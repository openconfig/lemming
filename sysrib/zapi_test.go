// Copyright 2022 Google LLC
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

// Copyright 2016, 2017 zebra project.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"syscall"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openconfig/lemming/gnmi"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	pb "github.com/openconfig/lemming/proto/sysrib"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/wenovus/gobgp/v3/pkg/zebra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestZServer(t *testing.T) {
	t.Run("hello", testHello)
	t.Run("RouteAdd", testRouteAdd)
	t.Run("RouteRedistribution-routeReadyBeforeDial", func(t *testing.T) { testRouteRedistribution(t, true) })
	t.Run("RouteRedistribution-routeNotReadyBeforeDial", func(t *testing.T) { testRouteRedistribution(t, false) })
}

// SendMessage sends a zebra message to the given connection.
func SendMessage(t *testing.T, conn net.Conn, msg *zebra.Message) error {
	s, err := msg.Serialize(zebra.MaxSoftware)
	if err != nil {
		return err
	}
	n, err := conn.Write(s)
	t.Logf("SendMessage: %d bytes written", n)
	return err
}

func Dial() (net.Conn, error) {
	conn, err := net.DialTimeout("unix", "/tmp/zserv.api", 10*time.Second)
	if err != nil {
		return nil, err
	}
	return conn, err
}

func ZAPIServerStart(t *testing.T) *ZServer {
	t.Helper()
	sysribServer, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}

	s, err := StartZServer(context.Background(), "unix:/tmp/zserv.api", 0, sysribServer)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func testHello(t *testing.T) {
	s := ZAPIServerStart(t)
	defer s.Stop()

	conn, err := Dial()
	if err != nil {
		t.Errorf("Dial failed %v\n", err)
		return
	}
	defer conn.Close()

	serverVersion := zebra.MaxZapiVer
	serverSoftware := zebra.MaxSoftware
	msg := zebra.Message{
		Header: zebra.Header{
			Len:     zebra.HeaderSize(serverVersion),
			Marker:  zebra.HeaderMarker(serverVersion),
			Version: serverVersion,
			Command: zebra.Hello.ToEach(serverVersion, serverSoftware),
		},
		Body: &zebra.HelloBody{},
	}
	if err := SendMessage(t, conn, &msg); err != nil {
		t.Error(err)
	}
}

func testRouteAdd(t *testing.T) {
	tests := []struct {
		desc   string
		inBody *zebra.IPRouteBody
	}{{
		desc: "IPv4",
		inBody: &zebra.IPRouteBody{Prefix: zebra.Prefix{
			Family:    syscall.AF_INET,
			PrefixLen: 24,
			Prefix:    net.IPv4(10, 0, 0, 0),
		}},
	}, {
		desc: "IPv6",
		inBody: &zebra.IPRouteBody{Prefix: zebra.Prefix{
			Family:    syscall.AF_INET6,
			PrefixLen: 49,
			Prefix:    net.ParseIP("2001:aaaa:bbbb::"),
		}},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			s := ZAPIServerStart(t)
			defer s.Stop()

			conn, err := Dial()
			if err != nil {
				t.Errorf("Dial failed %v\n", err)
				return
			}
			defer conn.Close()

			serverVersion := zebra.MaxZapiVer
			serverSoftware := zebra.MaxSoftware
			helloMsg := zebra.Message{
				Header: zebra.Header{
					Len:     zebra.HeaderSize(serverVersion),
					Marker:  zebra.HeaderMarker(serverVersion),
					Version: serverVersion,
					Command: zebra.Hello.ToEach(serverVersion, serverSoftware),
				},
				Body: &zebra.HelloBody{},
			}
			if err := SendMessage(t, conn, &helloMsg); err != nil {
				t.Error(err)
			}

			msg := zebra.Message{
				Header: zebra.Header{
					Len:     zebra.HeaderSize(serverVersion),
					Marker:  zebra.HeaderMarker(serverVersion),
					Version: serverVersion,
					Command: zebra.RouteAdd.ToEach(serverVersion, serverSoftware),
				},
				Body: tt.inBody,
			}
			if err := SendMessage(t, conn, &msg); err != nil {
				t.Error(err)
			}
		})
	}
}

// testRouteRedistribution tests that a route redistribution is sent by the
// ZAPI server when a new route is added to sysrib, or when a new client ZAPI
// connection is added where there already exists routes in the sysrib.
//
// - routeReadyBeforeDial specifies whether to make the route ready before the
// client dials to the ZAPI server.
func testRouteRedistribution(t *testing.T, routeReadyBeforeDial bool) {
	routesQuery := programmedRoutesQuery(t)
	tests := []struct {
		desc              string
		inAddIntfAction   *AddIntfAction
		inSetRouteRequest *pb.SetRouteRequest
		wantRoutes        []*dpb.Route
	}{{
		desc: "IPv4",
		inAddIntfAction: &AddIntfAction{
			name:    "eth0",
			ifindex: 0,
			enabled: true,
			prefix:  "192.168.1.1/24",
			niName:  "DEFAULT",
		},
		inSetRouteRequest: &pb.SetRouteRequest{
			AdminDistance: 10,
			Metric:        10,
			Prefix: &pb.Prefix{
				Family:     pb.Prefix_FAMILY_IPV4,
				Address:    "10.0.0.0",
				MaskLength: 8,
			},
			Nexthops: []*pb.Nexthop{{
				Type:    pb.Nexthop_TYPE_IPV4,
				Address: "192.168.1.42",
				Weight:  1,
			}},
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
					Cidr: "192.168.1.0/24",
				},
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Hops: []*dpb.NextHop{{
						Dev: &dpb.NextHop_Port{
							Port: "eth0",
						},
						Weight: 0,
					}},
				},
			},
		}},
	}, {
		desc: "IPv6",
		inAddIntfAction: &AddIntfAction{
			name:    "eth0",
			ifindex: 0,
			enabled: true,
			prefix:  "2001::aaaa/42",
			niName:  "DEFAULT",
		},
		inSetRouteRequest: &pb.SetRouteRequest{
			AdminDistance: 10,
			Metric:        10,
			Prefix: &pb.Prefix{
				Family:     pb.Prefix_FAMILY_IPV6,
				Address:    "4242::4242",
				MaskLength: 42,
			},
			Nexthops: []*pb.Nexthop{{
				Type:    pb.Nexthop_TYPE_IPV6,
				Address: "2001::ffff",
				Weight:  1,
			}},
		},
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				VrfId: uint64(0),
				Prefix: &dpb.RoutePrefix_Cidr{
					Cidr: "4242::/42",
				},
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Hops: []*dpb.NextHop{{
						Ip: &dpb.NextHop_IpStr{IpStr: "2001::ffff"},
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
					Cidr: "2001::/42",
				},
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Hops: []*dpb.NextHop{{
						Dev: &dpb.NextHop_Port{
							Port: "eth0",
						},
						Weight: 0,
					}},
				},
			},
		}},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			s, err := New(nil)
			if err != nil {
				t.Fatal(err)
			}

			grpcServer := grpc.NewServer()
			gnmiServer, err := gnmi.New(grpcServer, "local", nil)
			if err != nil {
				t.Fatal(err)
			}
			client := gnmiServer.LocalClient()
			if err := s.Start(context.Background(), client, "local", "unix:/tmp/zserv.api"); err != nil {
				t.Fatalf("cannot start sysrib server, %v", err)
			}
			defer s.Stop()

			c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
			if err != nil {
				t.Fatalf("cannot create ygnmi client: %v", err)
			}

			configureInterface(t, tt.inAddIntfAction, c)

			if routeReadyBeforeDial {
				if _, err := s.SetRoute(context.Background(), tt.inSetRouteRequest); err != nil {
					t.Fatalf("Got unexpected error during call to SetRoute: %v", err)
				}

				var err error
				for i := 0; i != maxGNMIWaitQuanta; i++ {
					var routes []*dpb.Route
					if routes, err = ygnmi.GetAll(context.Background(), c, routesQuery); err == nil {
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
					t.Fatalf("Routes not resolved in time limit: %v", err)
				}
			}

			conn, err := Dial()
			if err != nil {
				t.Errorf("Dial failed %v\n", err)
				return
			}
			defer conn.Close()

			version := zebra.MaxZapiVer
			software := zebra.MaxSoftware

			SendMessage(t, conn, &zebra.Message{
				Header: zebra.Header{
					Len:     zebra.HeaderSize(version),
					Marker:  zebra.HeaderMarker(version),
					Version: version,
					Command: zebra.Hello.ToEach(version, software),
				},
				Body: &zebra.HelloBody{},
			})

			// The first message is expected to be a capabilities message which is
			// discarded since no client uses it.
			if _, err := zebra.ReceiveSingleMsg(topicLogger, conn, version, software, "test-client"); err != nil {
				t.Fatalf("Got error during call to first ReceiveSingleMsg: %v", err)
			}

			if !routeReadyBeforeDial {
				if _, err := s.SetRoute(context.Background(), tt.inSetRouteRequest); err != nil {
					t.Fatalf("Got unexpected error during call to SetRoute: %v", err)
				}
			}

			m, err := zebra.ReceiveSingleMsg(topicLogger, conn, version, software, "test-client")
			if err != nil {
				t.Fatal(err)
			} else if m == nil {
				t.Fatal("got empty message")
			}
			if got, want := m.Header.Command.ToCommon(version, software), zebra.RedistributeRouteAdd; got != want {
				t.Errorf("Got message %s, want %s", got, want)
			}
		})
	}
}
