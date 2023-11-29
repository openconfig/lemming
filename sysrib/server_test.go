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
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/gopacket"
	"github.com/openconfig/ygnmi/schemaless"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	dpb "github.com/openconfig/lemming/proto/dataplane"
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

func programmedRoutesQuery(t *testing.T) ygnmi.WildcardQuery[*dpb.Route] {
	t.Helper()
	q, err := schemaless.NewWildcard[*dpb.Route]("/dataplane/routes/route[prefix=*][vrf=*]", gnmi.InternalOrigin)
	if err != nil {
		t.Fatal(err)
	}
	return q
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

	if _, err := gnmiclient.Update(context.Background(), yclient, ocpath.Root().Interface(intf.name).State(), ocintf); err != nil {
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
	return zero
}

func mapAddressTo6Bytes(v4Address [4]byte) [16]byte {
	ipv6Bytes := [16]byte{0x20, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	copy(ipv6Bytes[v4v6ConversionStartPos:], v4Address[:])
	return ipv6Bytes
}

func mapAddressSliceTo6BytesSlice(v4Address []byte) []byte {
	ipv6Bytes := []byte{0x20, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	copy(ipv6Bytes[v4v6ConversionStartPos:], v4Address)
	return ipv6Bytes
}

func mapAddressTo6BytesSlice(v4Address [4]byte) []byte {
	r := mapAddressTo6Bytes(v4Address)
	return r[:]
}

func mapPrefixLenTo6(pfxLen int) int {
	return pfxLen + v4v6ConversionStartPos*8
}

// mapAddressTo6 converts an input address to an IPv6 address that is *not* an
// IPv4-mapped IPv6 address. This allows running the same test cases except
// with IPv6 addresses.
func mapAddressTo6(t *testing.T, addrStr string) string {
	t.Helper()
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
		newAddrStr = netip.PrefixFrom(netip.AddrFrom16(ipv6Bytes), mapPrefixLenTo6(pfxLen)).String()
	} else {
		newAddrStr = netip.AddrFrom16(ipv6Bytes).String()
	}

	return newAddrStr
}

func mapResolvedRouteTo6(t *testing.T, route *dpb.Route) {
	route.Prefix.Cidr = mapAddressTo6(t, route.Prefix.GetCidr())
	for _, nh := range route.GetNextHops().GetHops() {
		if nh.GetNextHopIp() != "" {
			nh.NextHopIp = mapAddressTo6(t, nh.GetNextHopIp())
		}
		if nh.GetGue() != nil {
			nh.GetGue().DstIp = mapAddressSliceTo6BytesSlice(nh.GetGue().DstIp)
			nh.GetGue().SrcIp = mapAddressSliceTo6BytesSlice(nh.GetGue().SrcIp)
		}
	}
}

func selectGUEHeaders(t *testing.T, v4 bool, layers ...gopacket.SerializableLayer) []gopacket.SerializableLayer {
	layerN := len(layers)
	if layerN == 0 || layerN%2 != 0 {
		t.Fatalf("Input layers is not even and non-zero: %v", layerN)
	}
	if v4 {
		return layers[:layerN/2]
	}
	return layers[layerN/2:]
}

func mapPrefixTo6(t *testing.T, prefix *pb.Prefix) {
	if prefix.Family == pb.Prefix_FAMILY_IPV4 {
		prefix.Family = pb.Prefix_FAMILY_IPV6
		prefix.MaskLength = uint32(mapPrefixLenTo6(int(prefix.MaskLength)))
	}
	prefix.Address = mapAddressTo6(t, prefix.Address)
}

type SetAndProgramPair struct {
	SetRouteRequestAction *SetRouteRequestAction
	ResolvedRoutes        []*ResolvedRoute
}

func getConnectedIntfSetupVarsV4() ([]*AddIntfAction, []*dpb.Route) {
	return []*AddIntfAction{{
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
		}, {
			name:    "eth5",
			ifindex: 5,
			enabled: false,
			prefix:  "192.168.6.1/24",
			niName:  "DEFAULT",
		}}, []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "192.168.1.0/24",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						Interface: &dpb.OCInterface{
							Interface: "eth0",
						},
					}},
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "192.168.2.0/24",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						Interface: &dpb.OCInterface{
							Interface: "eth1",
						},
					}},
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "192.168.3.0/24",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						Interface: &dpb.OCInterface{
							Interface: "eth2",
						},
					}},
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "192.168.4.0/24",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						Interface: &dpb.OCInterface{
							Interface: "eth3",
						},
					}},
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "192.168.5.0/24",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						Interface: &dpb.OCInterface{
							Interface: "eth4",
						},
					}},
				},
			},
		}}
}

func getConnectedIntfSetupVars(t *testing.T) ([]*AddIntfAction, []*dpb.Route) {
	inInterfaces, wantConnectedRoutes := getConnectedIntfSetupVarsV4()
	inInterface6s, wantConnectedRoute6s := getConnectedIntfSetupVarsV4()
	for _, intf := range inInterface6s {
		intf.prefix = mapAddressTo6(t, intf.prefix)
	}
	for _, route := range wantConnectedRoute6s {
		mapResolvedRouteTo6(t, route)
	}
	inInterfaces = append(inInterfaces, inInterface6s...)
	wantConnectedRoutes = append(wantConnectedRoutes, wantConnectedRoute6s...)
	return inInterfaces, wantConnectedRoutes
}

func TestServer(t *testing.T) {
	routesQuery := programmedRoutesQuery(t)
	inInterfaces, wantConnectedRoutes := getConnectedIntfSetupVars(t)

	tests := []struct {
		desc string
		// noV6 indicates whether the test should be run again after converting all
		// IPv4 addresses to IPv6.
		noV6               bool
		inSetRouteRequests []*SetRouteRequestAction
		wantRoutes         []*dpb.Route
		wantErr            bool
	}{{
		desc: "Route Additions",
		inSetRouteRequests: []*SetRouteRequestAction{{
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
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "10.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						NextHopIp: "192.168.5.42",
						Interface: &dpb.OCInterface{
							Interface: "eth4",
						},
					}},
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "20.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						NextHopIp: "192.168.5.42",
						Interface: &dpb.OCInterface{
							Interface: "eth4",
						},
					}},
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "30.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						NextHopIp: "192.168.5.42",
						Interface: &dpb.OCInterface{
							Interface: "eth4",
						},
					}},
				},
			},
		}},
	}, {
		desc: "Route Deletions",
		inSetRouteRequests: []*SetRouteRequestAction{{
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
					Weight:  1,
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
			Desc: "delete 2nd-level indirect route",
			RouteReq: &pb.SetRouteRequest{
				Delete:        true,
				AdminDistance: 10,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "20.0.0.0",
					MaskLength: 8,
				},
			},
		}, {
			Desc: "delete 3rd-level indirect route",
			RouteReq: &pb.SetRouteRequest{
				Delete:        true,
				AdminDistance: 10,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "30.0.0.0",
					MaskLength: 8,
				},
			},
		}, {
			Desc: "delete 3rd-level indirect route again",
			RouteReq: &pb.SetRouteRequest{
				Delete:        true,
				AdminDistance: 10,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "30.0.0.0",
					MaskLength: 8,
				},
			},
		}},
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "10.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
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
		desc:    "Invalid Route",
		noV6:    true,
		wantErr: true,
		inSetRouteRequests: []*SetRouteRequestAction{{
			Desc: "invalid route",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 20,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV4,
					Address:    "10.0.0.0.256",
					MaskLength: 8,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "192.168.2.42",
				}},
			},
		}},
	}, {
		desc: "Unresolvable and ECMP",
		inSetRouteRequests: []*SetRouteRequestAction{{
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
					Address: "192.168.6.42",
				}},
			},
		}, {
			Desc: "ECMP route that resolves over up interface with higher admin distance",
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
				}, {
					Type:    pb.Nexthop_TYPE_IPV4,
					Address: "192.168.3.42",
				}},
			},
		}},
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "10.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0, 0},
					Hops: []*dpb.NextHop{{
						NextHopIp: "192.168.2.42",
						Interface: &dpb.OCInterface{
							Interface: "eth1",
						},
					}, {
						NextHopIp: "192.168.3.42",
						Interface: &dpb.OCInterface{
							Interface: "eth2",
						},
					}},
				},
			},
		}},
	}, {
		desc: "IPv4-mapped IPv6",
		noV6: true,
		inSetRouteRequests: []*SetRouteRequestAction{{
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
		}},
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "10.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
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
				Cidr:            "20.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
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
				Cidr:            "2002::/49",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
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
		desc: "IPv4-mapped IPv6 in IPv4 format",
		noV6: true,
		inSetRouteRequests: []*SetRouteRequestAction{{
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
		}, {
			Desc: "1st level indirect route ipv4-mapped ipv6 using IPv4-formatted address",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 10,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV6,
					Address:    "2003::aaaa",
					MaskLength: 49,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV6,
					Address: "10.10.10.10",
				}},
			},
		}},
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "10.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
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
				Cidr:            "10.10.0.0/16",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						NextHopIp: "192.168.2.42",
						Interface: &dpb.OCInterface{
							Interface: "eth1",
						},
					}},
				},
			},
		}, {
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "2003::/49",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						NextHopIp: "192.168.2.42",
						Interface: &dpb.OCInterface{
							Interface: "eth1",
						},
					}},
				},
			},
		}},
	}, {
		desc: "IPv4-mapped IPv6 route deleted due to specific route not resolvable",
		noV6: true,
		inSetRouteRequests: []*SetRouteRequestAction{{
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
			Desc: "secondary 1st level indirect route that is more specific but higher admin distance and unresolvable",
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
					Address: "192.168.6.42",
				}},
			},
		}, {
			Desc: "1st level indirect route ipv4-mapped ipv6 using IPv4-formatted address",
			RouteReq: &pb.SetRouteRequest{
				AdminDistance: 10,
				Metric:        10,
				Prefix: &pb.Prefix{
					Family:     pb.Prefix_FAMILY_IPV6,
					Address:    "2003::aaaa",
					MaskLength: 49,
				},
				Nexthops: []*pb.Nexthop{{
					Type:    pb.Nexthop_TYPE_IPV6,
					Address: "10.10.10.10",
				}},
			},
		}},
		wantRoutes: []*dpb.Route{{
			Prefix: &dpb.RoutePrefix{
				NetworkInstance: "DEFAULT",
				Cidr:            "10.0.0.0/8",
			},
			Hop: &dpb.Route_NextHops{
				NextHops: &dpb.NextHopList{
					Weight: []uint64{0},
					Hops: []*dpb.NextHop{{
						NextHopIp: "192.168.1.42",
						Interface: &dpb.OCInterface{
							Interface: "eth0",
						},
					}},
				},
			},
		}},
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

	c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("cannot create ygnmi client: %v", err)
	}

	for _, intf := range inInterfaces {
		configureInterface(t, intf, c)
	}

	// Wait for Sysrib to pick up the connected prefixes.
	for i := 0; i != maxGNMIWaitQuanta; i++ {
		var routes []*dpb.Route
		routes, err = ygnmi.GetAll(context.Background(), c, routesQuery)
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
			for _, v4 := range []bool{true, false} {
				desc := "v4"
				if !v4 {
					if tt.noV6 {
						continue
					}
					desc = "v6"
					// Convert all v4 addresses to v6.
					for _, req := range tt.inSetRouteRequests {
						mapPrefixTo6(t, req.RouteReq.Prefix)
						for _, nh := range req.RouteReq.Nexthops {
							if nh.Type == pb.Nexthop_TYPE_IPV4 {
								nh.Type = pb.Nexthop_TYPE_IPV6
							}
							nh.Address = mapAddressTo6(t, nh.Address)
						}
					}
					for _, route := range tt.wantRoutes {
						mapResolvedRouteTo6(t, route)
					}
				}

				t.Run(desc, func(t *testing.T) {
					for _, req := range tt.inSetRouteRequests {
						// TODO(wenbli): Test SetRouteResponse
						_, err := s.SetRoute(context.Background(), req.RouteReq)
						hasErr := err != nil
						if hasErr != tt.wantErr {
							t.Fatalf("%s: got error during call to SetRoute: %v, wantErr: %v", req.Desc, err, tt.wantErr)
						}
					}

					routes, err := ygnmi.GetAll(context.Background(), c, routesQuery)
					if err != nil {
						t.Fatal(err)
					}
					if diff := cmp.Diff(append(append([]*dpb.Route{}, wantConnectedRoutes...), tt.wantRoutes...), routes, protocmp.Transform(), protocmp.SortRepeatedFields(new(dpb.NextHopList), "hops"), cmpopts.SortSlices(func(a, b *dpb.Route) bool {
						return a.GetPrefix().GetCidr() < b.GetPrefix().GetCidr()
					})); diff != "" {
						t.Errorf("routes not equal to wantRoutes (-want, +got):\n%s", diff)
					}

					// Clean-up
					for _, req := range tt.inSetRouteRequests {
						isDelete := req.RouteReq.Delete
						req.RouteReq.Delete = true
						_, err := s.SetRoute(context.Background(), req.RouteReq)
						req.RouteReq.Delete = isDelete
						hasErr := err != nil
						if hasErr != tt.wantErr {
							t.Fatalf("%s: got error during call to SetRoute: %v, wantErr: %v", req.Desc, err, tt.wantErr)
						}
					}

					if routes, err = ygnmi.GetAll(context.Background(), c, routesQuery); err != nil {
						t.Fatal(err)
					}
					if diff := cmp.Diff(wantConnectedRoutes, routes, protocmp.Transform(), protocmp.SortRepeatedFields(new(dpb.NextHopList), "hops")); diff != "" {
						t.Errorf("routes not equal to wantConnectedRoutes (-want, +got):\n%s", diff)
					}
				})
			}
		})
	}
}

func gueHeader(t *testing.T, layers ...gopacket.SerializableLayer) []byte {
	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{}, layers...); err != nil {
		t.Fatalf("failed to serialize GUE headers: %v", err)
	}
	return buf.Bytes()
}

func addrToBytes(t testing.TB, isV4 bool, v4 string) []byte {
	t.Helper()
	addr, err := netip.ParseAddr(v4)
	if err != nil {
		t.Fatal(err)
	}
	if isV4 {
		return addr.AsSlice()
	}
	return append(addr.AsSlice(), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
}

func TestBGPGUEPolicy(t *testing.T) {
	routesQuery := programmedRoutesQuery(t)
	inInterfaces, wantConnectedRoutes := getConnectedIntfSetupVars(t)

	// Note: This is a sequential test -- each test case depends on the previous one.
	tests := []struct {
		desc string
		// Skip this step in V4 testing.
		skipV4 bool
		// Skip this step in V6 testing.
		skipV6             bool
		inSetRouteRequests []*pb.SetRouteRequest
		inAddPolicies      map[string]GUEPolicy
		inDeletePolicies   []string
		wantRoutes         func(v4 bool) []*dpb.Route
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
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
						}},
					},
				},
			}}
		},
	}, {
		desc: "Add Policy",
		inAddPolicies: map[string]GUEPolicy{
			"192.168.0.0/16": {
				dstPortv4: 8,
				dstPortv6: 16,
				srcIP4:    [4]byte{42, 42, 42, 42},
				srcIP6:    [16]byte{42, 42, 42, 42, 42},
			},
		},
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
							Encap: &dpb.NextHop_Gue{
								Gue: &dpb.GUE{
									SrcIp:   []byte{42, 42, 42, 42},
									DstIp:   []byte{192, 168, 1, 42},
									DstPort: 8,
								},
							},
						}},
					},
				},
			}}
		},
	}, {
		desc:             "Remove Policy",
		inDeletePolicies: []string{"192.168.0.0/16"},
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
						}},
					},
				},
			}}
		},
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
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "30.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
						}},
					},
				},
			}}
		},
	}, {
		desc: "Add Policy for second BGP route",
		inAddPolicies: map[string]GUEPolicy{
			"10.10.0.0/16": {
				dstPortv4: 9,
				dstPortv6: 18,
				srcIP4:    [4]byte{43, 43, 43, 43},
				srcIP6:    [16]byte{43, 43, 43, 43, 43},
			},
		},
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "30.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
							Encap: &dpb.NextHop_Gue{
								Gue: &dpb.GUE{
									SrcIp:   []byte{43, 43, 43, 43},
									DstIp:   []byte{10, 10, 10, 10},
									DstPort: 9,
								},
							},
						}},
					},
				},
			}}
		},
	}, {
		desc:             "Remove Policy for second BGP route",
		inDeletePolicies: []string{"10.10.0.0/16"},
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "30.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
						}},
					},
				},
			}}
		},
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
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "30.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "40.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
						}},
					},
				},
			}}
		},
	}, {
		desc: "Add a policy that applies to two BGP routes",
		inAddPolicies: map[string]GUEPolicy{
			"10.0.0.0/8": {
				dstPortv4: 8,
				dstPortv6: 16,
				srcIP4:    [4]byte{8, 8, 8, 8},
				srcIP6:    [16]byte{8, 8, 8, 8, 8},
			},
		},
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "30.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
							Encap: &dpb.NextHop_Gue{
								Gue: &dpb.GUE{
									SrcIp:   []byte{8, 8, 8, 8},
									DstIp:   []byte{10, 10, 10, 10},
									DstPort: 8,
								},
							},
						}},
					},
				},
			}, {
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "40.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
							Encap: &dpb.NextHop_Gue{
								Gue: &dpb.GUE{
									SrcIp:   []byte{8, 8, 8, 8},
									DstIp:   []byte{10, 10, 20, 20},
									DstPort: 8,
								},
							},
						}},
					},
				},
			}}
		},
	}, {
		desc: "Add a more specific policy that applies to a BGP route",
		inAddPolicies: map[string]GUEPolicy{
			"10.10.20.0/24": {
				dstPortv4: 16,
				dstPortv6: 32,
				srcIP4:    [4]byte{16, 16, 16, 16},
				srcIP6:    [16]byte{16, 16, 16, 16, 16},
			},
		},
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "30.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
							Encap: &dpb.NextHop_Gue{
								Gue: &dpb.GUE{
									SrcIp:   []byte{8, 8, 8, 8},
									DstIp:   []byte{10, 10, 10, 10},
									DstPort: 8,
								},
							},
						}},
					},
				},
			}, {
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "40.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
							Encap: &dpb.NextHop_Gue{
								Gue: &dpb.GUE{
									SrcIp:   []byte{16, 16, 16, 16},
									DstIp:   []byte{10, 10, 20, 20},
									DstPort: 16,
								},
							},
						}},
					},
				},
			}}
		},
	}, {
		desc:             "Remove the less-specific policy",
		inDeletePolicies: []string{"10.0.0.0/8"},
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "30.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "40.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
							Encap: &dpb.NextHop_Gue{
								Gue: &dpb.GUE{
									SrcIp:   []byte{16, 16, 16, 16},
									DstIp:   []byte{10, 10, 20, 20},
									DstPort: 16,
								},
							},
						}},
					},
				},
			}}
		},
	}, {
		desc:   "Add an IPv4-mapped IPv6 BGP route over the more specific policy",
		skipV6: true,
		inSetRouteRequests: []*pb.SetRouteRequest{{
			AdminDistance: 20, // EBGP
			Metric:        10,
			Prefix: &pb.Prefix{
				Family:     pb.Prefix_FAMILY_IPV6,
				Address:    "4242::",
				MaskLength: 42,
			},
			Nexthops: []*pb.Nexthop{{
				Type:    pb.Nexthop_TYPE_IPV6,
				Address: "::ffff:10.10.20.30",
			}},
		}},
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "30.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "40.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
							Encap: &dpb.NextHop_Gue{
								Gue: &dpb.GUE{
									SrcIp:   []byte{16, 16, 16, 16},
									DstIp:   []byte{10, 10, 20, 20},
									DstPort: 16,
								},
							},
						}},
					},
				},
			}, {
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "4242::/42",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
							Encap: &dpb.NextHop_Gue{
								Gue: &dpb.GUE{
									SrcIp:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									DstIp:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									DstPort: 32,
								},
							},
						}},
					},
				},
			}}
		},
	}, {
		desc:             "Remove the more-specific policy",
		skipV4:           true,
		inDeletePolicies: []string{"10.10.20.0/24"},
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "30.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "40.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
						}},
					},
				},
			}}
		},
	}, {
		desc:             "Remove the more-specific policy",
		skipV6:           true,
		inDeletePolicies: []string{"10.10.20.0/24"},
		wantRoutes: func(v4 bool) []*dpb.Route {
			return []*dpb.Route{{
				Prefix: &dpb.RoutePrefix{
					NetworkInstance: "DEFAULT",
					Cidr:            "10.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "20.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "30.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "40.0.0.0/8",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
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
					Cidr:            "4242::/42",
				},
				Hop: &dpb.Route_NextHops{
					NextHops: &dpb.NextHopList{
						Weight: []uint64{0},
						Hops: []*dpb.NextHop{{
							NextHopIp: "192.168.1.42",
							Interface: &dpb.OCInterface{
								Interface: "eth0",
							},
						}},
					},
				},
			}}
		},
	}}

	for _, v4 := range []bool{true, false} {
		desc := "v4"
		if !v4 {
			desc = "v6"
		}

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

		c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
		if err != nil {
			t.Fatalf("cannot create ygnmi client: %v", err)
		}

		for _, intf := range inInterfaces {
			configureInterface(t, intf, c)
		}

		// Wait for Sysrib to pick up the connected prefixes.
		for i := 0; i != maxGNMIWaitQuanta; i++ {
			var routes []*dpb.Route
			routes, err = ygnmi.GetAll(context.Background(), c, routesQuery)
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

		t.Run(desc, func(t *testing.T) {
			for _, tt := range tests {
				if v4 && tt.skipV4 {
					continue
				}
				if !v4 && tt.skipV6 {
					continue
				}

				wantRoutes := tt.wantRoutes(v4)
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
					inAddPolicies := map[string]GUEPolicy{}
					for prefix, gueHeaders := range tt.inAddPolicies {
						inAddPolicies[mapAddressTo6(t, prefix)] = mapPolicyTo6(gueHeaders)
					}
					tt.inAddPolicies = inAddPolicies
					for i := range tt.inDeletePolicies {
						tt.inDeletePolicies[i] = mapAddressTo6(t, tt.inDeletePolicies[i])
					}
					if !v4 {
						for _, route := range wantRoutes {
							mapResolvedRouteTo6(t, route)
						}
					}
				}

				t.Run(tt.desc, func(t *testing.T) {
					for _, routeReq := range tt.inSetRouteRequests {
						if _, err := s.SetRoute(context.Background(), routeReq); err != nil {
							t.Fatalf("Got unexpected error during call to SetRoute: %v", err)
						}
					}
					for prefix, policy := range tt.inAddPolicies {
						s.setGUEPolicy(context.Background(), prefix, policy)
					}
					for _, prefix := range tt.inDeletePolicies {
						s.deleteGUEPolicy(context.Background(), prefix)
					}

					routes, err := ygnmi.GetAll(context.Background(), c, routesQuery)
					if err != nil {
						t.Fatal(err)
					}
					if diff := cmp.Diff(append(append([]*dpb.Route{}, wantConnectedRoutes...), wantRoutes...), routes, protocmp.Transform(), protocmp.SortRepeatedFields(new(dpb.NextHopList), "hops"), cmpopts.SortSlices(func(a, b *dpb.Route) bool {
						return a.GetPrefix().GetCidr() < b.GetPrefix().GetCidr()
					})); diff != "" {
						t.Errorf("routes not equal to wantRoutes (-want, +got):\n%s", diff)
					}
				})
			}
		})
	}
}
