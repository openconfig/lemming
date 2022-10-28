package sysrib

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gribigo/afthelper"
	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"

	dpb "github.com/openconfig/lemming/proto/dataplane"
	pb "github.com/openconfig/lemming/proto/sysrib"
)

const (
	// Each quantum is 100 ms
	maxGNMIWaitQuanta = 100
)

type AddIntfAction struct {
	name    string
	ifindex int32
	enabled bool
	prefix  string
	niName  string
}

type SetRouteRequestAction struct {
	Desc     string
	RouteReq *pb.SetRouteRequest
}

type FakeDataplane struct {
	dpb.HALClient

	mu             sync.Mutex
	incomingRoutes []*ResolvedRoute

	// failRoutes are routes that will fail to program.
	failRoutes []*ResolvedRoute
}

func (dp *FakeDataplane) ProgramRoute(r *ResolvedRoute) error {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	for _, failroute := range dp.failRoutes {
		if reflect.DeepEqual(r, failroute) {
			return fmt.Errorf("route failed to program: %v", r)
		}
	}
	dp.incomingRoutes = append(dp.incomingRoutes, r)
	return nil
}

// SetupFailRoutes sets routes that will fail to program.
func (dp *FakeDataplane) SetupFailRoutes(failRoutes []*ResolvedRoute) {
	dp.failRoutes = failRoutes
}

func (dp *FakeDataplane) GetRoutes() []*ResolvedRoute {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	return dp.incomingRoutes
}

func (dp *FakeDataplane) ClearQueue() {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.incomingRoutes = []*ResolvedRoute{}
}

func NewFakeDataplane() *FakeDataplane {
	return &FakeDataplane{}
}

// routeSliceToMap converts a slice of ResolvedRoute to a map keyed by their
// RouteKeys. It returns an error if any of the routes were nil or if there is
// a duplicate.
func routeSliceToMap(rs []*ResolvedRoute) (map[RouteKey]*ResolvedRoute, error) {
	ret := map[RouteKey]*ResolvedRoute{}
	for _, rr := range rs {
		if rr == nil {
			return nil, fmt.Errorf("Got nil route in ResolvedRoute slice")
		}
		if existing, ok := ret[rr.RouteKey]; ok {
			return nil, fmt.Errorf("Got duplicate route key:\nFirst: %+v\nDuplicate: %+v", existing, rr)
		}
		ret[rr.RouteKey] = rr
	}
	return ret, nil
}

func checkResolvedRoutesEqual(got, want []*ResolvedRoute) error {
	gotRoutes, err := routeSliceToMap(got)
	if err != nil {
		return err
	}
	wantRoutes, err := routeSliceToMap(want)
	if err != nil {
		return err
	}

	if diff := cmp.Diff(gotRoutes, wantRoutes); diff != "" {
		return fmt.Errorf("Resolved routes are not equal: (-got, +want):\n%s", diff)
	}
	return nil
}

func configureInterface(t *testing.T, intf AddIntfAction, yclient *ygnmi.Client) {
	t.Helper()

	ocintf := &oc.Interface{}
	ocintf.Name = ygot.String(intf.name)
	ocintf.Enabled = ygot.Bool(intf.enabled)
	ocintf.Ifindex = ygot.Uint32(uint32(intf.ifindex))
	ss := strings.Split(intf.prefix, "/")
	if len(ss) != 2 {
		t.Fatalf("Invalid prefix: %q", intf.prefix)
	}
	ocaddr := ocintf.GetOrCreateSubinterface(0).GetOrCreateIpv4().GetOrCreateAddress(ss[0])
	plen, err := strconv.Atoi(ss[1])
	if err != nil {
		t.Fatalf("Invalid prefix: %v", err)
	}
	ocaddr.PrefixLength = ygot.Uint8(uint8(plen))

	if _, err := gnmiclient.Replace(context.Background(), yclient, ocpath.Root().Interface(intf.name).State(), ocintf); err != nil {
		t.Fatalf("Cannot configure interface: %v", err)
	}
}

func TestServer(t *testing.T) {
	// TODO(wenbli): This test should be refactored such that the wantResolvedRoutes is inlined with the inSetRouteRequests for easier reading.
	tests := []struct {
		desc                       string
		inInterfaces               []AddIntfAction
		wantInitialConnectedRoutes []*ResolvedRoute
		inSetRouteRequests         []SetRouteRequestAction
		inFailRoutes               []*ResolvedRoute
		wantResolvedRoutes         [][]*ResolvedRoute
	}{{
		desc: "Route Additions", // TODO(wenbli): test route deletion in this test case once it's implemented.
		inInterfaces: []AddIntfAction{{
			name:    "eth0",
			ifindex: 0,
			enabled: true,
			prefix:  "192.168.1.1/24",
			niName:  "DEFAULT",
		}, {
			name:    "eth1",
			ifindex: 1,
			enabled: true,
			prefix:  "192.168.2.1/24",
			niName:  "DEFAULT",
		}, {
			name:    "eth2",
			ifindex: 2,
			enabled: true,
			prefix:  "192.168.3.1/24",
			niName:  "DEFAULT",
		}, {
			name:    "eth3",
			ifindex: 3,
			enabled: true,
			prefix:  "192.168.4.1/24",
			niName:  "DEFAULT",
		}, {
			name:    "eth4",
			ifindex: 4,
			enabled: true,
			prefix:  "192.168.5.1/24",
			niName:  "DEFAULT",
		}},
		wantInitialConnectedRoutes: []*ResolvedRoute{{
			RouteKey: RouteKey{
				Prefix: "192.168.1.0/24",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
					},
					Port: Interface{
						Name:  "eth0",
						Index: 0,
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "192.168.2.0/24",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
					},
					Port: Interface{
						Name:  "eth1",
						Index: 1,
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "192.168.3.0/24",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
					},
					Port: Interface{
						Name:  "eth2",
						Index: 2,
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "192.168.4.0/24",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
					},
					Port: Interface{
						Name:  "eth3",
						Index: 3,
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "192.168.5.0/24",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
					},
					Port: Interface{
						Name:  "eth4",
						Index: 4,
					},
				}: true,
			},
		}},
		inSetRouteRequests: []SetRouteRequestAction{{
			Desc: "1st level indirect route",
			RouteReq: &pb.SetRouteRequest{
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
					// TODO(wenbli): Implement WCMP, for all route requests in this test.
					Weight: 1,
				}},
			},
		}, {
			Desc: "2nd level indirect route",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 10,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "20.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "10.10.10.10",
				}},
			},
		}, {
			Desc: "3rd level indirect route",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 10,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "30.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "20.10.10.10",
				}},
			},
		}, {
			Desc: "secondary 1st level indirect route that has higher admin distance",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 20,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "192.168.2.42",
				}},
			},
		}, {
			Desc: "secondary 1st level indirect route that has lower admin distance",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 5,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "192.168.3.42",
				}},
			},
		}, {
			Desc: "secondary 1st level indirect route that has higher metric",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 5,
				Metric:        999,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "192.168.4.42",
				}},
			},
		}, {
			Desc: "secondary 1st level indirect route that has lower metric",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 5,
				Metric:        5,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "192.168.5.42",
				}},
			},
		}},
		wantResolvedRoutes: [][]*ResolvedRoute{{{
			RouteKey: RouteKey{
				Prefix: "10.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.1.42",
					},
					Port: Interface{
						Name:  "eth0",
						Index: 0,
					},
				}: true,
			},
		}}, {{
			RouteKey: RouteKey{
				Prefix: "20.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.1.42",
					},
					Port: Interface{
						Name:  "eth0",
						Index: 0,
					},
				}: true,
			},
		}}, {{
			RouteKey: RouteKey{
				Prefix: "30.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.1.42",
					},
					Port: Interface{
						Name:  "eth0",
						Index: 0,
					},
				}: true,
			},
		}}, {}, {{
			RouteKey: RouteKey{
				Prefix: "10.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.3.42",
					},
					Port: Interface{
						Name:  "eth2",
						Index: 2,
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "20.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.3.42",
					},
					Port: Interface{
						Name:  "eth2",
						Index: 2,
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "30.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.3.42",
					},
					Port: Interface{
						Name:  "eth2",
						Index: 2,
					},
				}: true,
			},
		}}, {}, {{
			RouteKey: RouteKey{
				Prefix: "10.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.5.42",
					},
					Port: Interface{
						Name:  "eth4",
						Index: 4,
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "20.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.5.42",
					},
					Port: Interface{
						Name:  "eth4",
						Index: 4,
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "30.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.5.42",
					},
					Port: Interface{
						Name:  "eth4",
						Index: 4,
					},
				}: true,
			},
		}}},
	}, {
		desc: "Unresolvable and ECMP",
		inInterfaces: []AddIntfAction{{
			name:    "eth0",
			ifindex: 0,
			enabled: false,
			prefix:  "192.168.1.1/24",
			niName:  "DEFAULT",
		}, {
			name:    "eth1",
			ifindex: 1,
			enabled: true,
			prefix:  "192.168.2.1/24",
			niName:  "DEFAULT",
		}, {
			name:    "eth2",
			ifindex: 2,
			enabled: true,
			prefix:  "192.168.3.1/24",
			niName:  "DEFAULT",
		}},
		wantInitialConnectedRoutes: []*ResolvedRoute{{
			RouteKey: RouteKey{
				Prefix: "192.168.2.0/24",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
					},
					Port: Interface{
						Name:  "eth1",
						Index: 1,
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "192.168.3.0/24",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
					},
					Port: Interface{
						Name:  "eth2",
						Index: 2,
					},
				}: true,
			},
		}},
		inSetRouteRequests: []SetRouteRequestAction{{
			Desc: "unresolvable route",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 10,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "15.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "11.10.10.10",
				}},
			},
		}, {
			Desc: "route that resolves over down interface",
			RouteReq: &pb.SetRouteRequest{
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
				}},
			},
		}, {
			Desc: "same route that resolves over up interface with higher admin distance",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 20,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "192.168.2.42",
				}},
			},
		}, {
			Desc: "ECMP",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 20,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "192.168.3.42",
				}},
			},
		}},
		wantResolvedRoutes: [][]*ResolvedRoute{{}, {}, {{
			RouteKey: RouteKey{
				Prefix: "10.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.2.42",
					},
					Port: Interface{
						Name:  "eth1",
						Index: 1,
					},
				}: true,
			},
		}}, {{
			RouteKey: RouteKey{
				Prefix: "10.0.0.0/8",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.2.42",
					},
					Port: Interface{
						Name:  "eth1",
						Index: 1,
					},
				}: true,
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
						Address:         "192.168.3.42",
					},
					Port: Interface{
						Name:  "eth2",
						Index: 2,
					},
				}: true,
			},
		}}},
	}, {
		desc: "test route program failures",
		inInterfaces: []AddIntfAction{{
			name:    "eth1",
			ifindex: 3,
			enabled: true,
			prefix:  "192.0.2.1/30",
			niName:  "DEFAULT",
		}, {
			name:    "eth2",
			ifindex: 4,
			enabled: true,
			prefix:  "192.0.2.5/30",
			niName:  "DEFAULT",
		}, {
			name:    "eth3",
			ifindex: 5,
			enabled: true,
			prefix:  "192.0.2.9/30",
			niName:  "DEFAULT",
		}},
		inFailRoutes: []*ResolvedRoute{{
			RouteKey: RouteKey{
				Prefix: "192.0.2.0/30",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
					},
					Port: Interface{
						Name:  "eth1",
						Index: 3,
					},
				}: true,
			},
		}},
		wantInitialConnectedRoutes: []*ResolvedRoute{{
			RouteKey: RouteKey{
				Prefix: "192.0.2.4/30",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
					},
					Port: Interface{
						Name:  "eth2",
						Index: 4,
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "192.0.2.8/30",
				NIName: "DEFAULT",
			},
			Nexthops: map[ResolvedNexthop]bool{
				{
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: "DEFAULT",
					},
					Port: Interface{
						Name:  "eth3",
						Index: 5,
					},
				}: true,
			},
		}},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// TODO(wenbli): Don't re-create gNMI server, simply erase it and then reconnect to it afterwards.
			grpcServer := grpc.NewServer()
			gnmiServer, err := gnmi.New(grpcServer, "local")
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

			dp := NewFakeDataplane()
			dp.SetupFailRoutes(tt.inFailRoutes)
			s, err := New(dp)
			if err != nil {
				t.Fatal(err)
			}

			// Update the interface configuration on the gNMI server.
			client := gnmiServer.LocalClient()
			if err := s.Start(client, "local"); err != nil {
				t.Fatalf("cannot start sysrib server, %v", err)
			}

			c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
			if err != nil {
				t.Fatalf("cannot create ygnmi client: %v", err)
			}
			for _, intf := range tt.inInterfaces {
				configureInterface(t, intf, c)
			}

			// Wait for Sysrib to pick up the connected prefixes.
			for i := 0; i != maxGNMIWaitQuanta; i++ {
				if err = checkResolvedRoutesEqual(dp.GetRoutes(), tt.wantInitialConnectedRoutes); err == nil {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			if err != nil {
				t.Fatalf("After initial interface operations: %v", err)
			}
			dp.ClearQueue()

			for i, routeReq := range tt.inSetRouteRequests {
				// TODO(wenbli): Test SetRouteResponse
				if _, err := s.SetRoute(context.Background(), routeReq.RouteReq); err != nil {
					t.Fatalf("%s: Got unexpected error during call to SetRoute: %v", routeReq.Desc, err)
				}
				if err := checkResolvedRoutesEqual(dp.GetRoutes(), tt.wantResolvedRoutes[i]); err != nil {
					t.Fatalf("%s: %v", routeReq.Desc, err)
				}
				dp.ClearQueue()
			}
		})
	}
}
