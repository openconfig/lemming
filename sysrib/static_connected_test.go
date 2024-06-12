// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sysrib

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	gpb "github.com/openconfig/gnmi/proto/gnmi"

	dpb "github.com/openconfig/lemming/proto/dataplane"
)

var staticp = ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).
	Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)

func installStaticRoute(t *testing.T, c *ygnmi.Client, route *oc.NetworkInstance_Protocol_Static) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(50*time.Second))
	defer cancel()

	if _, err := ygnmi.Replace(ctx, c, staticp.Static(*route.Prefix).Config(), route); err != nil {
		t.Fatal(err)
	}
	if _, err := ygnmi.Await(ctx, c, staticp.Static(*route.Prefix).State(), route); err != nil {
		r, getErr := ygnmi.Get(ctx, c, staticp.Static(*route.Prefix).State())
		if getErr != nil {
			t.Error(getErr)
		}
		t.Fatalf("Did not get static route (%v) state within deadline: %v", r, err)
	}
}

// TestStaticRouteAndIntfs tests static and interface route addition/deletion.
func TestStaticRouteAndIntfs(t *testing.T) {
	routesQuery := programmedRoutesQuery(t)

	// Note: This is a sequential test -- each test case depends on the previous one.
	tests := []struct {
		desc               string
		inStaticRouteOp    func(t *testing.T, c *ygnmi.Client)
		inConnectedRouteOp func(t *testing.T, c *ygnmi.Client)
		wantRoutes         []*dpb.Route
	}{{
		desc: "Add static route v4",
		inStaticRouteOp: func(t *testing.T, c *ygnmi.Client) {
			installStaticRoute(t, c, &oc.NetworkInstance_Protocol_Static{
				Prefix: ygot.String("10.0.0.0/8"),
				NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
					"single": {
						Index:   ygot.String("single"),
						NextHop: oc.UnionString("192.168.1.42"),
						Recurse: ygot.Bool(true),
					},
				},
			})
		},
		inConnectedRouteOp: func(t *testing.T, c *ygnmi.Client) {
			if _, err := gnmiclient.Replace(context.Background(), c, ocpath.Root().Interface("eth0").State(), &oc.Interface{
				Name:    ygot.String("eth0"),
				Enabled: ygot.Bool(true),
				Ifindex: ygot.Uint32(1),
				Subinterface: map[uint32]*oc.Interface_Subinterface{
					0: {
						Index: ygot.Uint32(0),
						Ipv4: &oc.Interface_Subinterface_Ipv4{
							Address: map[string]*oc.Interface_Subinterface_Ipv4_Address{
								"192.168.1.1": {
									Ip:           ygot.String("192.168.1.1"),
									PrefixLength: ygot.Uint8(24),
								},
							},
						},
					},
				},
			}); err != nil {
				t.Fatalf("Cannot configure interface: %v", err)
			}
		},
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "192.168.1.0/24",
			},
			Hop: &dpb.Route_Interface{
				Interface: &dpb.OCInterface{
					Interface: "eth0",
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "10.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weights: []uint64{0},
					Hops: []*dpb.NextHop{{
						NextHopIp: "192.168.1.42",
						Interface: &dpb.OCInterface{
							Interface: "eth0",
						},
					}},
				},
			},
		}},
	}, {
		desc: "Add static route v6",
		inStaticRouteOp: func(t *testing.T, c *ygnmi.Client) {
			installStaticRoute(t, c, &oc.NetworkInstance_Protocol_Static{
				Prefix: ygot.String(mapAddressTo6(t, "10.0.0.0/8")),
				NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
					"single": {
						Index:   ygot.String("single"),
						NextHop: oc.UnionString(mapAddressTo6(t, "192.168.1.42")),
						Recurse: ygot.Bool(true),
					},
				},
			})
		},
		inConnectedRouteOp: func(t *testing.T, c *ygnmi.Client) {
			if _, err := gnmiclient.Update(context.Background(), c, ocpath.Root().Interface("eth0").State(), &oc.Interface{
				Name:    ygot.String("eth0"),
				Enabled: ygot.Bool(true),
				Ifindex: ygot.Uint32(1),
				Subinterface: map[uint32]*oc.Interface_Subinterface{
					0: {
						Index: ygot.Uint32(0),
						Ipv6: &oc.Interface_Subinterface_Ipv6{
							Address: map[string]*oc.Interface_Subinterface_Ipv6_Address{
								mapAddressTo6(t, "192.168.1.1"): {
									Ip:           ygot.String(mapAddressTo6(t, "192.168.1.1")),
									PrefixLength: ygot.Uint8(uint8(mapPrefixLenTo6(24))),
								},
							},
						},
					},
				},
			}); err != nil {
				t.Fatalf("Cannot configure interface: %v", err)
			}
		},
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "192.168.1.0/24",
			},
			Hop: &dpb.Route_Interface{
				Interface: &dpb.OCInterface{
					Interface: "eth0",
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "10.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weights: []uint64{0},
					Hops: []*dpb.NextHop{{
						NextHopIp: "192.168.1.42",
						Interface: &dpb.OCInterface{
							Interface: "eth0",
						},
					}},
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            mapAddressTo6(t, "192.168.1.0/24"),
			},
			Hop: &dpb.Route_Interface{
				Interface: &dpb.OCInterface{
					Interface: "eth0",
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            mapAddressTo6(t, "10.0.0.0/8"),
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weights: []uint64{0},
					Hops: []*dpb.NextHop{{
						NextHopIp: mapAddressTo6(t, "192.168.1.42"),
						Interface: &dpb.OCInterface{
							Interface: "eth0",
						},
					}},
				},
			},
		}},
	}, {
		desc: "Delete static route v4",
		inStaticRouteOp: func(t *testing.T, c *ygnmi.Client) {
			if _, err := ygnmi.Delete[*oc.NetworkInstance_Protocol_Static](context.Background(), c, staticp.Static("10.0.0.0/8").Config()); err != nil {
				t.Fatal(err)
			}
		},
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "192.168.1.0/24",
			},
			Hop: &dpb.Route_Interface{
				Interface: &dpb.OCInterface{
					Interface: "eth0",
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            mapAddressTo6(t, "192.168.1.0/24"),
			},
			Hop: &dpb.Route_Interface{
				Interface: &dpb.OCInterface{
					Interface: "eth0",
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            mapAddressTo6(t, "10.0.0.0/8"),
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weights: []uint64{0},
					Hops: []*dpb.NextHop{{
						NextHopIp: mapAddressTo6(t, "192.168.1.42"),
						Interface: &dpb.OCInterface{
							Interface: "eth0",
						},
					}},
				},
			},
		}},
	}, {
		desc: "Delete static route v6",
		inStaticRouteOp: func(t *testing.T, c *ygnmi.Client) {
			if _, err := ygnmi.Delete[*oc.NetworkInstance_Protocol_Static](context.Background(), c, staticp.Static(mapAddressTo6(t, "10.0.0.0/8")).Config()); err != nil {
				t.Fatal(err)
			}
		},
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "192.168.1.0/24",
			},
			Hop: &dpb.Route_Interface{
				Interface: &dpb.OCInterface{
					Interface: "eth0",
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            mapAddressTo6(t, "192.168.1.0/24"),
			},
			Hop: &dpb.Route_Interface{
				Interface: &dpb.OCInterface{
					Interface: "eth0",
				},
			},
		}},
	}, {
		desc: "Delete connected interface",
		inStaticRouteOp: func(t *testing.T, c *ygnmi.Client) {
			if _, err := ygnmi.Delete[*oc.NetworkInstance_Protocol_Static](context.Background(), c, staticp.Static(mapAddressTo6(t, "10.0.0.0/8")).Config()); err != nil {
				t.Fatal(err)
			}
		},
		inConnectedRouteOp: func(t *testing.T, c *ygnmi.Client) {
			if _, err := gnmiclient.Delete(context.Background(), c, ocpath.Root().Interface("eth0").State()); err != nil {
				t.Fatal(err)
			}
		},
	}}

	grpcServer := grpc.NewServer()
	gnmiServer, err := gnmi.New(grpcServer, "local", nil)
	if err != nil {
		t.Fatal(err)
	}
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	go func() {
		grpcServer.Serve(lis)
	}()

	s, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Update the interface configuration on the gNMI server.
	client := gnmiServer.LocalClient()
	if err := s.Start(context.Background(), client, "local", "", "/tmp/sysrib.api"); err != nil {
		t.Fatalf("cannot start sysrib server, %v", err)
	}
	defer s.Stop()

	stateC, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("cannot create ygnmi client: %v", err)
	}

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		t.Fatalf("cannot dial gNMI server, %v", err)
	}
	configC, err := ygnmi.NewClient(gpb.NewGNMIClient(conn), ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("cannot create ygnmi client: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.inStaticRouteOp != nil {
				tt.inStaticRouteOp(t, configC)
			}
			if tt.inConnectedRouteOp != nil {
				tt.inConnectedRouteOp(t, stateC)
			}

			var err error
			for i := 0; i != maxGNMIWaitQuanta; i++ {
				var routesVal []*ygnmi.Value[*dpb.Route]
				if routesVal, err = ygnmi.LookupAll(context.Background(), configC, routesQuery); err == nil {
					var routes []*dpb.Route
					for _, v := range routesVal {
						r, ok := v.Val()
						if ok {
							routes = append(routes, r)
						}
					}
					if diff := cmp.Diff(tt.wantRoutes, routes, protocmp.Transform(), protocmp.SortRepeatedFields(new(dpb.NextHopList), "hops"), cmpopts.SortSlices(func(a, b *dpb.Route) bool {
						return a.GetPrefix().GetCidr() < b.GetPrefix().GetCidr()
					})); diff != "" {
						err = fmt.Errorf("routes not equal to wantRoutes (-want, +got):\n%s", diff)
					} else {
						break
					}
				}
				time.Sleep(100 * time.Millisecond)
			}
			if err != nil {
				t.Fatalf("After static route operation: %v", err)
			}
		})
	}
}
