package sysrib

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gribigo/afthelper"
	pb "github.com/openconfig/lemming/proto/sysrib"
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
	incomingRoutes []*ResolvedRoute
}

func (dp *FakeDataplane) ProgramRoute(r *ResolvedRoute) error {
	dp.incomingRoutes = append(dp.incomingRoutes, r)
	// Assume all routes are programmed successfully.
	return nil
}

func (dp *FakeDataplane) GetRoutesAndClearQueue() []*ResolvedRoute {
	rs := dp.incomingRoutes
	dp.incomingRoutes = []*ResolvedRoute{}
	return rs
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

func TestServer(t *testing.T) {
	// TODO(wenbli): This test should be refactored such that the wantResolvedRoutes is inlined with the inSetRouteRequests for easier reading.
	tests := []struct {
		desc                       string
		inInterfaces               []AddIntfAction
		wantInitialConnectedRoutes []*ResolvedRoute
		inSetRouteRequests         []SetRouteRequestAction
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
					Family:     pb.Prefix_IPv4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
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
					Family:     pb.Prefix_IPv4,
					Address:    "20.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
					Address: "10.10.10.10",
				}},
			},
		}, {
			Desc: "3rd level indirect route",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 10,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_IPv4,
					Address:    "30.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
					Address: "20.10.10.10",
				}},
			},
		}, {
			Desc: "secondary 1st level indirect route that has higher admin distance",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 20,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_IPv4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
					Address: "192.168.2.42",
				}},
			},
		}, {
			Desc: "secondary 1st level indirect route that has lower admin distance",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 5,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_IPv4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
					Address: "192.168.3.42",
				}},
			},
		}, {
			Desc: "secondary 1st level indirect route that has higher metric",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 5,
				Metric:        999,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_IPv4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
					Address: "192.168.4.42",
				}},
			},
		}, {
			Desc: "secondary 1st level indirect route that has lower metric",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 5,
				Metric:        5,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_IPv4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
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
					Family:     pb.Prefix_IPv4,
					Address:    "15.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
					Address: "11.10.10.10",
				}},
			},
		}, {
			Desc: "route that resolves over down interface",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 10,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_IPv4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
					Address: "192.168.1.42",
				}},
			},
		}, {
			Desc: "same route that resolves over up interface with higher admin distance",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 20,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_IPv4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
					Address: "192.168.2.42",
				}},
			},
		}, {
			Desc: "ECMP",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 20,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_IPv4,
					Address:    "10.0.0.0",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_IPV4,
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
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dp := NewFakeDataplane()
			s, err := NewServer(dp)
			if err != nil {
				t.Fatal(err)
			}

			for _, intf := range tt.inInterfaces {
				if err := s.addInterface(intf.name, intf.ifindex, intf.enabled, intf.prefix, intf.niName); err != nil {
					t.Fatal(err)
				}
			}

			if err := checkResolvedRoutesEqual(dp.GetRoutesAndClearQueue(), tt.wantInitialConnectedRoutes); err != nil {
				t.Fatalf("After initial interface operations: %v", err)
			}

			for i, routeReq := range tt.inSetRouteRequests {
				// TODO(wenbli): Test SetRouteResponse
				if _, err := s.SetRoute(context.Background(), routeReq.RouteReq); err != nil {
					t.Fatalf("%s: Got unexpected error during call to SetRoute: %v", routeReq.Desc, err)
				}
				if err := checkResolvedRoutesEqual(dp.GetRoutesAndClearQueue(), tt.wantResolvedRoutes[i]); err != nil {
					t.Fatalf("%s: %v", routeReq.Desc, err)
				}
			}
		})
	}
}
