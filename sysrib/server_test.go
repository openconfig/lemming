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

package sysrib

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"reflect"
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

	pb "github.com/openconfig/lemming/proto/sysrib"
)

const (
	// Each quantum is 100 ms
	maxGNMIWaitQuanta = 200 // 20s
	// v4v6ConversionStartPos is the position in the IPv6 address byte
	// slice into which to start copying the 4 bytes of the IPv4 address
	// for conversion.
	v4v6ConversionStartPos = 8
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
	mu             sync.Mutex
	incomingRoutes []*ResolvedRoute

	// failRoutes are routes that the fake dataplane will choose to fail to
	// program.
	failRoutes []*ResolvedRoute
}

func (dp *FakeDataplane) ProgramRoute(r *ResolvedRoute) error {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	// Intentionally fail to program failRoutes.
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

func configureInterface(t *testing.T, intf *AddIntfAction, yclient *ygnmi.Client) {
	t.Helper()

	ocintf := &oc.Interface{}
	ocintf.Name = ygot.String(intf.name)
	ocintf.Enabled = ygot.Bool(intf.enabled)
	ocintf.Ifindex = ygot.Uint32(uint32(intf.ifindex))
	prefix, err := netip.ParsePrefix(intf.prefix)
	if err != nil {
		t.Fatalf("Invalid prefix: %q", intf.prefix)
	}
	switch {
	case prefix.Addr().Is4():
		ocaddr := ocintf.GetOrCreateSubinterface(0).GetOrCreateIpv4().GetOrCreateAddress(prefix.Addr().String())
		ocaddr.PrefixLength = ygot.Uint8(uint8(prefix.Bits()))
	case prefix.Addr().Is6():
		ocaddr := ocintf.GetOrCreateSubinterface(0).GetOrCreateIpv6().GetOrCreateAddress(prefix.Addr().String())
		ocaddr.PrefixLength = ygot.Uint8(uint8(prefix.Bits()))
	default:
		t.Fatalf("Prefix is neither IPv4 nor IPv6: %q", intf.prefix)
	}

	if _, err := gnmiclient.Replace(context.Background(), yclient, ocpath.Root().Interface(intf.name).State(), ocintf); err != nil {
		t.Fatalf("Cannot configure interface: %v", err)
	}
}

func mapPolicyTo6(h GUEPolicy) GUEPolicy {
	zero := GUEPolicy{}
	if h == zero {
		return h
	}
	zero.dstPortv6 = h.dstPortv4
	zero.srcIP6 = mapAddressTo6Bytes(h.srcIP4)
	zero.isV6 = true
	return zero
}

func mapPolicyHeadersTo6(h GUEHeaders) GUEHeaders {
	zero := GUEHeaders{}
	if h == zero {
		return h
	}
	zero.dstPortv6 = h.dstPortv4
	zero.srcIP6 = mapAddressTo6Bytes(h.srcIP4)
	zero.dstIP6 = mapAddressTo6Bytes(h.dstIP4)
	zero.isV6 = true
	return zero
}

func mapAddressTo6Bytes(v4Address [4]byte) [16]byte {
	ipv6Bytes := [16]byte{0x20, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	copy(ipv6Bytes[v4v6ConversionStartPos:], v4Address[:])
	return ipv6Bytes
}

// mapAddressTo6 converts an input address to an IPv6 address that is *not* an
// IPv4-mapped IPv6 address. This allows running the same test cases except
// with IPv6 addresses.
func mapAddressTo6(t *testing.T, addrStr string) string {
	if addrStr == "" {
		return ""
	}
	var addr netip.Addr
	var pfxLen int
	isPrefix := strings.Contains(addrStr, "/")
	if isPrefix {
		pfx, err := netip.ParsePrefix(addrStr)
		if err != nil {
			t.Fatalf("not a valid prefix: %q", addrStr)
		}
		addr = pfx.Addr()
		pfxLen = pfx.Bits()
	} else {
		var err error
		if addr, err = netip.ParseAddr(addrStr); err != nil {
			t.Fatalf("not a valid address: %q", addrStr)
		}
	}

	// Don't convert IPv6 addresses.
	if !addr.Is4() {
		return addrStr
	}

	ipv6Bytes := mapAddressTo6Bytes(*(*[4]byte)(addr.AsSlice()))

	var newAddrStr string
	if isPrefix {
		newAddrStr = netip.PrefixFrom(netip.AddrFrom16(ipv6Bytes), pfxLen+v4v6ConversionStartPos*8).String()
	} else {
		newAddrStr = netip.AddrFrom16(ipv6Bytes).String()
	}

	return newAddrStr
}

func mapResolvedRouteTo6(t *testing.T, route *ResolvedRoute) {
	route.Prefix = mapAddressTo6(t, route.Prefix)
	// Since this is not a pointer need to overwrite.
	nexthops := map[ResolvedNexthop]bool{}
	for nh, v := range route.Nexthops {
		nh.Address = mapAddressTo6(t, nh.Address)
		nh.GUEHeaders = mapPolicyHeadersTo6(nh.GUEHeaders)
		nexthops[nh] = v
	}
	route.Nexthops = nexthops
}

func mapPrefixTo6(t *testing.T, prefix *pb.Prefix) {
	if prefix.Family == pb.Prefix_FAMILY_IPV4 {
		prefix.Family = pb.Prefix_FAMILY_IPV6
		prefix.MaskLength += v4v6ConversionStartPos * 8
	}
	prefix.Address = mapAddressTo6(t, prefix.Address)
}

type SetAndProgramPair struct {
	SetRouteRequestAction *SetRouteRequestAction
	ResolvedRoutes        []*ResolvedRoute
}

func TestServer(t *testing.T) {
	tests := []struct {
		desc string
		// noV6 indicates whether the test should be run again after converting all
		// IPv4 addresses to IPv6.
		noV6                       bool
		inInterfaces               []*AddIntfAction
		wantInitialConnectedRoutes []*ResolvedRoute
		inSetAndProgramPairs       []*SetAndProgramPair
		inFailRoutes               []*ResolvedRoute
	}{{
		desc: "Route Additions", // TODO(wenbli): test route deletion in this test case once it's implemented.
		inInterfaces: []*AddIntfAction{{
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
		inSetAndProgramPairs: []*SetAndProgramPair{{
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{{
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
			}},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{{
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
			}},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{{
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
			}},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{{
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
			}},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{{
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
			}},
		}},
	}, {
		desc: "Unresolvable and ECMP",
		inInterfaces: []*AddIntfAction{{
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
		inSetAndProgramPairs: []*SetAndProgramPair{{
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{{
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
			}},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{{
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
			}},
		}},
	}, {
		desc: "test route program failures",
		inInterfaces: []*AddIntfAction{{
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
	}, {
		desc: "IPv4-mapped IPv6",
		noV6: true,
		inInterfaces: []*AddIntfAction{{
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
		}},
		// TODO(wenbli): Test when a more specific route is not resolvable --> route is deleted.
		inSetAndProgramPairs: []*SetAndProgramPair{{
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{{
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
			}},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
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
			},
			ResolvedRoutes: []*ResolvedRoute{{
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
			}},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
				Desc: "3rd level indirect route ipv4-mapped ipv6",
				RouteReq: &pb.SetRouteRequest{
					AdminDistance: 10,
					Metric:        10,
					Prefix: &pb.Prefix{
						Family:     pb.Prefix_FAMILY_IPV6,
						Address:    "2002::aaaa",
						MaskLength: 49,
					},
					Nexthops: []*pb.Nexthop{{
						Type:    pb.Nexthop_TYPE_IPV6,
						Address: "::ffff:20.10.10.10",
					}},
				},
			},
			ResolvedRoutes: []*ResolvedRoute{{
				RouteKey: RouteKey{
					Prefix: "2002::/49",
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
			}},
		}, {
			SetRouteRequestAction: &SetRouteRequestAction{
				Desc: "secondary 1st level indirect route that is more specific but higher admin distance",
				RouteReq: &pb.SetRouteRequest{
					AdminDistance: 20,
					Metric:        10,
					Prefix: &pb.Prefix{
						Family:     pb.Prefix_FAMILY_IPV4,
						Address:    "10.10.0.0",
						MaskLength: 16,
					},
					Nexthops: []*pb.Nexthop{{
						Type:    pb.Nexthop_TYPE_IPV4,
						Address: "192.168.2.42",
					}},
				},
			},
			ResolvedRoutes: []*ResolvedRoute{{
				RouteKey: RouteKey{
					Prefix: "10.10.0.0/16",
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
			}, {
				RouteKey: RouteKey{
					Prefix: "20.0.0.0/8",
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
			}, {
				RouteKey: RouteKey{
					Prefix: "2002::/49",
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
			}},
		}},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			for _, v4 := range []bool{true, false} {
				desc := "v4"
				if !v4 {
					if tt.noV6 {
						continue
					}
					desc = "v6"
					// Convert all v4 addresses to v6.
					for _, intf := range tt.inInterfaces {
						intf.prefix = mapAddressTo6(t, intf.prefix)
					}
					for _, pair := range tt.inSetAndProgramPairs {
						mapPrefixTo6(t, pair.SetRouteRequestAction.RouteReq.Prefix)
						for _, nh := range pair.SetRouteRequestAction.RouteReq.Nexthops {
							if nh.Type == pb.Nexthop_TYPE_IPV4 {
								nh.Type = pb.Nexthop_TYPE_IPV6
							}
							nh.Address = mapAddressTo6(t, nh.Address)
						}
						for _, route := range pair.ResolvedRoutes {
							mapResolvedRouteTo6(t, route)
						}
					}
					for _, route := range tt.wantInitialConnectedRoutes {
						mapResolvedRouteTo6(t, route)
					}
					for _, route := range tt.inFailRoutes {
						mapResolvedRouteTo6(t, route)
					}
				}

				t.Run(desc, func(t *testing.T) {
					// TODO(wenbli): Don't re-create gNMI server, simply erase it and then reconnect to it afterwards.
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

					dp := NewFakeDataplane()
					dp.SetupFailRoutes(tt.inFailRoutes)
					s, err := New()
					if err != nil {
						t.Fatal(err)
					}

					// Update the interface configuration on the gNMI server.
					client := gnmiServer.LocalClient()
					if err := s.Start(client, "local", ""); err != nil {
						t.Fatalf("cannot start sysrib server, %v", err)
					}
					s.dataplane = dp
					defer s.Stop()

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

					for _, pair := range tt.inSetAndProgramPairs {
						// TODO(wenbli): Test SetRouteResponse
						if _, err := s.SetRoute(context.Background(), pair.SetRouteRequestAction.RouteReq); err != nil {
							t.Fatalf("%s: Got unexpected error during call to SetRoute: %v", pair.SetRouteRequestAction.Desc, err)
						}
						if err := checkResolvedRoutesEqual(dp.GetRoutes(), pair.ResolvedRoutes); err != nil {
							t.Fatalf("%s: %v", pair.SetRouteRequestAction.Desc, err)
						}
						dp.ClearQueue()
					}
				})
			}
		})
	}
}

func TestBGPGUEPolicy(t *testing.T) {
	intfs := []*AddIntfAction{{
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
	}}

	wantConnectedRoutes := []*ResolvedRoute{{
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
	}}

	// Note: This is a sequential test -- each test case depends on the previous one.
	tests := []struct {
		desc               string
		inSetRouteRequests []*pb.SetRouteRequest
		inAddPolicies      map[string]GUEPolicy
		inDeletePolicies   []string
		wantResolvedRoutes []*ResolvedRoute
	}{{
		desc: "Add static and BGP routes",
		inSetRouteRequests: []*pb.SetRouteRequest{{
			AdminDistance: 10, // not BGP
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
		}, {
			AdminDistance: 20, // EBGP
			Metric:        10,
			Prefix: &pb.Prefix{
				Family:     pb.Prefix_FAMILY_IPV4,
				Address:    "20.0.0.0",
				MaskLength: 8,
			},
			Nexthops: []*pb.Nexthop{{
				Type:    pb.Nexthop_TYPE_IPV4,
				Address: "192.168.1.42",
			}},
		}},
		wantResolvedRoutes: []*ResolvedRoute{{
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
		}, {
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
		}},
	}, {
		desc: "Add Policy",
		inAddPolicies: map[string]GUEPolicy{
			"192.168.0.0/16": {
				dstPortv4: 8,
				srcIP4:    [4]byte{42, 42, 42, 42},
				isV6:      false,
			},
		},
		wantResolvedRoutes: []*ResolvedRoute{{
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
					GUEHeaders: GUEHeaders{
						GUEPolicy: GUEPolicy{
							dstPortv4: 8,
							srcIP4:    [4]byte{42, 42, 42, 42},
							isV6:      false,
						},
						dstIP4: [4]byte{192, 168, 1, 42},
					},
				}: true,
			},
		}},
	}, {
		desc:             "Remove Policy",
		inDeletePolicies: []string{"192.168.0.0/16"},
		wantResolvedRoutes: []*ResolvedRoute{{
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
		}},
	}, {
		desc: "Add BGP route that resolves over the static route",
		inSetRouteRequests: []*pb.SetRouteRequest{{
			AdminDistance: 20, // EBGP
			Metric:        10,
			Prefix: &pb.Prefix{
				Family:     pb.Prefix_FAMILY_IPV4,
				Address:    "30.0.0.0",
				MaskLength: 8,
			},
			Nexthops: []*pb.Nexthop{{
				Type:    pb.Nexthop_TYPE_IPV4,
				Address: "10.10.10.10",
			}},
		}},
		wantResolvedRoutes: []*ResolvedRoute{{
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
		}},
	}, {
		desc: "Add Policy for second BGP route",
		inAddPolicies: map[string]GUEPolicy{
			"10.10.0.0/16": {
				dstPortv4: 9,
				srcIP4:    [4]byte{43, 43, 43, 43},
				isV6:      false,
			},
		},
		wantResolvedRoutes: []*ResolvedRoute{{
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
					GUEHeaders: GUEHeaders{
						GUEPolicy: GUEPolicy{
							dstPortv4: 9,
							srcIP4:    [4]byte{43, 43, 43, 43},
							isV6:      false,
						},
						dstIP4: [4]byte{10, 10, 10, 10},
					},
				}: true,
			},
		}},
	}, {
		desc:             "Remove Policy for second BGP route",
		inDeletePolicies: []string{"10.10.0.0/16"},
		wantResolvedRoutes: []*ResolvedRoute{{
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
		}},
	}, {
		desc: "Add another BGP route that resolves over the static route",
		inSetRouteRequests: []*pb.SetRouteRequest{{
			AdminDistance: 20, // EBGP
			Metric:        10,
			Prefix: &pb.Prefix{
				Family:     pb.Prefix_FAMILY_IPV4,
				Address:    "40.0.0.0",
				MaskLength: 8,
			},
			Nexthops: []*pb.Nexthop{{
				Type:    pb.Nexthop_TYPE_IPV4,
				Address: "10.10.20.20",
			}},
		}},
		wantResolvedRoutes: []*ResolvedRoute{{
			RouteKey: RouteKey{
				Prefix: "40.0.0.0/8",
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
		}},
	}, {
		desc: "Add a policy that applies to two BGP routes",
		inAddPolicies: map[string]GUEPolicy{
			"10.0.0.0/8": {
				dstPortv4: 8,
				srcIP4:    [4]byte{8, 8, 8, 8},
				isV6:      false,
			},
		},
		wantResolvedRoutes: []*ResolvedRoute{{
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
					GUEHeaders: GUEHeaders{
						GUEPolicy: GUEPolicy{
							dstPortv4: 8,
							srcIP4:    [4]byte{8, 8, 8, 8},
							isV6:      false,
						},
						dstIP4: [4]byte{10, 10, 10, 10},
					},
				}: true,
			},
		}, {
			RouteKey: RouteKey{
				Prefix: "40.0.0.0/8",
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
					GUEHeaders: GUEHeaders{
						GUEPolicy: GUEPolicy{
							dstPortv4: 8,
							srcIP4:    [4]byte{8, 8, 8, 8},
							isV6:      false,
						},
						dstIP4: [4]byte{10, 10, 20, 20},
					},
				}: true,
			},
		}},
	}, {
		desc: "Add a more specific policy that applies to a BGP route",
		inAddPolicies: map[string]GUEPolicy{
			"10.10.20.0/24": {
				dstPortv4: 16,
				srcIP4:    [4]byte{16, 16, 16, 16},
				isV6:      false,
			},
		},
		wantResolvedRoutes: []*ResolvedRoute{{
			RouteKey: RouteKey{
				Prefix: "40.0.0.0/8",
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
					GUEHeaders: GUEHeaders{
						GUEPolicy: GUEPolicy{
							dstPortv4: 16,
							srcIP4:    [4]byte{16, 16, 16, 16},
							isV6:      false,
						},
						dstIP4: [4]byte{10, 10, 20, 20},
					},
				}: true,
			},
		}},
	}, {
		desc:             "Remove the less-specific policy",
		inDeletePolicies: []string{"10.0.0.0/8"},
		wantResolvedRoutes: []*ResolvedRoute{{
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
		}},
	}, {
		desc:             "Remove the more-specific policy",
		inDeletePolicies: []string{"10.10.20.0/24"},
		wantResolvedRoutes: []*ResolvedRoute{{
			RouteKey: RouteKey{
				Prefix: "40.0.0.0/8",
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
		}},
	}}

	for _, v4 := range []bool{true, false} {
		desc := "v4"
		if !v4 {
			desc = "v6"
			// Convert all v4 addresses to v6.
			for _, intf := range intfs {
				intf.prefix = mapAddressTo6(t, intf.prefix)
			}
			for _, route := range wantConnectedRoutes {
				mapResolvedRouteTo6(t, route)
			}
		}

		t.Run(desc, func(t *testing.T) {
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

			dp := NewFakeDataplane()
			s, err := New()
			if err != nil {
				t.Fatal(err)
			}

			// Update the interface configuration on the gNMI server.
			client := gnmiServer.LocalClient()
			if err := s.Start(client, "local", ""); err != nil {
				t.Fatalf("cannot start sysrib server, %v", err)
			}
			s.dataplane = dp
			defer s.Stop()

			c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
			if err != nil {
				t.Fatalf("cannot create ygnmi client: %v", err)
			}

			for _, intf := range intfs {
				configureInterface(t, intf, c)
			}
			// Wait for Sysrib to pick up the connected prefixes.
			for i := 0; i != maxGNMIWaitQuanta; i++ {
				if err = checkResolvedRoutesEqual(dp.GetRoutes(), wantConnectedRoutes); err == nil {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			if err != nil {
				t.Fatalf("After initial interface operations: %v", err)
			}
			dp.ClearQueue()

			for _, tt := range tests {
				if !v4 { // Convert v4 to v6.
					for _, req := range tt.inSetRouteRequests {
						mapPrefixTo6(t, req.Prefix)
						for _, nh := range req.Nexthops {
							if nh.Type == pb.Nexthop_TYPE_IPV4 {
								nh.Type = pb.Nexthop_TYPE_IPV6
							}
							nh.Address = mapAddressTo6(t, nh.Address)
						}
					}
					for _, route := range tt.wantResolvedRoutes {
						mapResolvedRouteTo6(t, route)
					}
					inAddPolicies := map[string]GUEPolicy{}
					for prefix, gueHeaders := range tt.inAddPolicies {
						inAddPolicies[mapAddressTo6(t, prefix)] = mapPolicyTo6(gueHeaders)
					}
					tt.inAddPolicies = inAddPolicies
					for i := range tt.inDeletePolicies {
						tt.inDeletePolicies[i] = mapAddressTo6(t, tt.inDeletePolicies[i])
					}
				}

				t.Run(tt.desc, func(t *testing.T) {
					for _, routeReq := range tt.inSetRouteRequests {
						if _, err := s.SetRoute(context.Background(), routeReq); err != nil {
							t.Fatalf("Got unexpected error during call to SetRoute: %v", err)
						}
					}
					for prefix, policy := range tt.inAddPolicies {
						s.setGUEPolicy(prefix, policy)
					}
					for _, prefix := range tt.inDeletePolicies {
						s.deleteGUEPolicy(prefix)
					}
					if err := checkResolvedRoutesEqual(dp.GetRoutes(), tt.wantResolvedRoutes); err != nil {
						fmt.Println("debug", tt.desc)
						s.rib.PrintRIB()
						t.Fatalf("%v", err)
					}
					dp.ClearQueue()
				})
			}
		})
	}
}
