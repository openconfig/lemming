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
	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/protobuf/testing/protocmp"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	dpb "github.com/openconfig/lemming/proto/dataplane"
)

var (
	staticp = ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).
		Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
)

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

func TestStaticRoute(t *testing.T) {
	routesQuery := programmedRoutesQuery(t)
	inInterfaces, wantConnectedRoutes := getConnectedIntfSetupVars(t)

	// Note: This is a sequential test -- each test case depends on the previous one.
	tests := []struct {
		desc            string
		inStaticRouteOp func(t *testing.T, c *ygnmi.Client)
		wantRoutes      []*dpb.Route
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
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				VrfId: uint64(0),
				Prefix: &dpb.RoutePrefix_Cidr{
					Cidr: "10.0.0.0/8",
				},
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Hops: []*dpb.NextHop{{
						Ip: &dpb.NextHop_IpStr{IpStr: "192.168.1.42"},
						Dev: &dpb.NextHop_Port{
							Port: "eth0",
						},
						Weight: 0,
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
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				VrfId: uint64(0),
				Prefix: &dpb.RoutePrefix_Cidr{
					Cidr: "10.0.0.0/8",
				},
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Hops: []*dpb.NextHop{{
						Ip: &dpb.NextHop_IpStr{IpStr: "192.168.1.42"},
						Dev: &dpb.NextHop_Port{
							Port: "eth0",
						},
						Weight: 0,
					}},
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				VrfId: uint64(0),
				Prefix: &dpb.RoutePrefix_Cidr{
					Cidr: mapAddressTo6(t, "10.0.0.0/8"),
				},
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Hops: []*dpb.NextHop{{
						Ip: &dpb.NextHop_IpStr{IpStr: mapAddressTo6(t, "192.168.1.42")},
						Dev: &dpb.NextHop_Port{
							Port: "eth0",
						},
						Weight: 0,
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
				VrfId: uint64(0),
				Prefix: &dpb.RoutePrefix_Cidr{
					Cidr: mapAddressTo6(t, "10.0.0.0/8"),
				},
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Hops: []*dpb.NextHop{{
						Ip: &dpb.NextHop_IpStr{IpStr: mapAddressTo6(t, "192.168.1.42")},
						Dev: &dpb.NextHop_Port{
							Port: "eth0",
						},
						Weight: 0,
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
	if err := s.Start(context.Background(), client, "local", ""); err != nil {
		t.Fatalf("cannot start sysrib server, %v", err)
	}
	defer s.Stop()

	stateC, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("cannot create ygnmi client: %v", err)
	}

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		t.Fatalf("cannot dial gNMI server, %v", err)
	}
	configC, err := ygnmi.NewClient(gpb.NewGNMIClient(conn), ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("cannot create ygnmi client: %v", err)
	}

	for _, intf := range inInterfaces {
		configureInterface(t, intf, stateC)
	}

	// Wait for Sysrib to pick up the connected prefixes.
	for i := 0; i != maxGNMIWaitQuanta; i++ {
		var routes []*dpb.Route
		routes, err = ygnmi.GetAll(context.Background(), configC, routesQuery)
		if err == nil {
			if diff := cmp.Diff(wantConnectedRoutes, routes, protocmp.Transform(), protocmp.SortRepeatedFields(new(dpb.NextHopList), "hops")); diff != "" {
				err = fmt.Errorf("routes not equal to wantConnectedRoutes (-want, +got):\n%s", diff)
			} else {
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("After initial interface operations: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.inStaticRouteOp(t, configC)

			var err error
			for i := 0; i != maxGNMIWaitQuanta; i++ {
				var routes []*dpb.Route
				routes, err = ygnmi.GetAll(context.Background(), configC, routesQuery)
				if err != nil {
					t.Fatal(err)
				}
				if diff := cmp.Diff(append(append([]*dpb.Route{}, wantConnectedRoutes...), tt.wantRoutes...), routes, protocmp.Transform(), protocmp.SortRepeatedFields(new(dpb.NextHopList), "hops"), cmpopts.SortSlices(func(a, b *dpb.Route) bool {
					return a.GetPrefix().GetCidr() < b.GetPrefix().GetCidr()
				})); diff != "" {
					err = fmt.Errorf("routes not equal to wantRoutes (-want, +got):\n%s", diff)
				}
			}
			if err != nil {
				t.Fatalf("After static route operation: %v", err)
			}
		})
	}
}
