// Copyright 2021 Google LLC
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
	"fmt"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openconfig/gribigo/afthelper"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/ygot/ygot"
)

func mustCIDR(s string) net.IPNet {
	_, cidr, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return *cidr
}

func TestRoutesFromConfig(t *testing.T) {
	tests := []struct {
		desc       string
		inConfig   *oc.Root
		wantRoutes map[string]*niConnected
		wantErr    bool
	}{{
		desc: "only default NI with explicit interfaces",
		inConfig: func() *oc.Root {
			d := &oc.Root{}

			dni := d.GetOrCreateNetworkInstance("DEFAULT")
			dni.Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			ei := dni.GetOrCreateInterface("eth0-assoc")
			ei.Interface = ygot.String("eth0")
			ei.Subinterface = ygot.Uint32(0)

			d.GetOrCreateInterface("eth0").GetOrCreateSubinterface(0).GetOrCreateIpv4().GetOrCreateAddress("1.1.1.1").PrefixLength = ygot.Uint8(24)
			d.GetOrCreateInterface("eth0").GetOrCreateSubinterface(0).GetOrCreateIpv6().GetOrCreateAddress("2001::ffff").PrefixLength = ygot.Uint8(42)
			return d
		}(),
		wantRoutes: map[string]*niConnected{
			"DEFAULT": {
				N: "DEFAULT",
				T: oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE,
				Rts: []*Route{{
					Prefix: "1.1.1.0/24",
					Connected: &Interface{
						Name:         "eth0",
						Subinterface: 0,
					},
				}, {
					Prefix: "2001::/42",
					Connected: &Interface{
						Name:         "eth0",
						Subinterface: 0,
					},
				}},
			},
		},
	}, {
		desc: "only default NI with implicit interfaces",
		inConfig: func() *oc.Root {
			d := &oc.Root{}
			dni := d.GetOrCreateNetworkInstance("DEFAULT")
			dni.Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			d.GetOrCreateInterface("eth0").GetOrCreateSubinterface(0).GetOrCreateIpv4().GetOrCreateAddress("192.0.2.1").PrefixLength = ygot.Uint8(24)
			d.GetOrCreateInterface("eth0").GetOrCreateSubinterface(0).GetOrCreateIpv6().GetOrCreateAddress("1999::eeee").PrefixLength = ygot.Uint8(42)
			return d
		}(),
		wantRoutes: map[string]*niConnected{
			"DEFAULT": {
				N: "DEFAULT",
				T: oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE,
				Rts: []*Route{{
					Prefix: "192.0.2.0/24",
					Connected: &Interface{
						Name:         "eth0",
						Subinterface: 0,
					},
				}, {
					Prefix: "1999::/42",
					Connected: &Interface{
						Name:         "eth0",
						Subinterface: 0,
					},
				}},
			},
		},
	}, {
		desc: "non-default NI, with multiple routes",
		inConfig: func() *oc.Root {
			d := &oc.Root{}
			dni := d.GetOrCreateNetworkInstance("DEFAULT")
			dni.Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			d.GetOrCreateInterface("eth0").GetOrCreateSubinterface(0).GetOrCreateIpv4().GetOrCreateAddress("192.0.2.1").PrefixLength = ygot.Uint8(24)
			d.GetOrCreateInterface("eth0").GetOrCreateSubinterface(0).GetOrCreateIpv6().GetOrCreateAddress("1999::eeee").PrefixLength = ygot.Uint8(42)

			vrf := d.GetOrCreateNetworkInstance("VRF-1")
			vrf.Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_L3VRF
			for i, spec := range []struct {
				i string
				s uint32
			}{{"eth0", 1}, {"eth24", 42}, {"eth42", 84}} {
				if err := vrf.AppendInterface(&oc.NetworkInstance_Interface{
					Id:           ygot.String(fmt.Sprintf("vrf-if%d", i)),
					Interface:    ygot.String(spec.i),
					Subinterface: ygot.Uint32(spec.s),
				}); err != nil {
					panic(fmt.Errorf("cannot add vrf interface, %v", err))
				}
			}
			d.GetOrCreateInterface("eth0").GetOrCreateSubinterface(1).GetOrCreateIpv4().GetOrCreateAddress("10.0.0.1").PrefixLength = ygot.Uint8(24)
			d.GetOrCreateInterface("eth24").GetOrCreateSubinterface(42).GetOrCreateIpv6().GetOrCreateAddress("4242::4242").PrefixLength = ygot.Uint8(42)
			d.GetOrCreateInterface("eth42").GetOrCreateSubinterface(84).GetOrCreateIpv4().GetOrCreateAddress("10.0.2.0").PrefixLength = ygot.Uint8(24)
			return d
		}(),
		wantRoutes: map[string]*niConnected{
			"DEFAULT": {
				N: "DEFAULT",
				T: oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE,
				Rts: []*Route{{
					Prefix: "192.0.2.0/24",
					Connected: &Interface{
						Name:         "eth0",
						Subinterface: 0,
					},
				}, {
					Prefix: "1999::/42",
					Connected: &Interface{
						Name:         "eth0",
						Subinterface: 0,
					},
				}},
			},
			"VRF-1": {
				N: "VRF-1",
				T: oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_L3VRF,
				Rts: []*Route{{
					Prefix: "10.0.0.0/24",
					Connected: &Interface{
						Name:         "eth0",
						Subinterface: 1,
					},
				}, {
					Prefix: "10.0.2.0/24",
					Connected: &Interface{
						Name:         "eth42",
						Subinterface: 84,
					},
				}, {
					Prefix: "4242::/42",
					Connected: &Interface{
						Name:         "eth24",
						Subinterface: 42,
					},
				}},
			},
		},
	}, {
		desc: "multiple default instances",
		inConfig: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateNetworkInstance("one").Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			d.GetOrCreateNetworkInstance("two").Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		wantErr: true,
	}, {
		desc: "invalid NI type",
		inConfig: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateNetworkInstance("one").Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_L2P2P
			return d
		}(),
		wantErr: true,
	}, {
		desc: "skip if with no subif",
		inConfig: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateNetworkInstance("D").Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			d.GetOrCreateInterface("eth0")
			return d
		}(),
		wantRoutes: map[string]*niConnected{
			"D": {
				N: "D",
				T: oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE,
			},
		},
	}, {
		desc: "no default instance",
		inConfig: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0")
			return d
		}(),
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := connectedRoutesFromConfig(tt.inConfig)
			if (err != nil) != tt.wantErr {
				t.Fatalf("did not get expected error, got: %v wantErr? %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(got, tt.wantRoutes); diff != "" {
				t.Fatalf("did not get expected routes, diff(-got,+want):\n%s", diff)
			}
		})
	}
}

func TestEgressInterface(t *testing.T) {
	tests := []struct {
		desc  string
		inCfg *oc.Root
		// keyed by network-instance
		inAddRoutes   map[string][]*Route
		inNI          string
		inIP          net.IPNet
		wantInterface []*Interface
		wantErr       bool
	}{{
		desc: "v4 single connected route",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.0").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inNI: "DEFAULT",
		inIP: mustCIDR("192.0.2.1/32"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
		},
	}, {
		desc: "v6 single connected route",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001::ab:cd:ef").
				PrefixLength = ygot.Uint8(42)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inNI: "DEFAULT",
		inIP: mustCIDR("2001::ab:cd:ef/128"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
		},
	}, {
		desc: "v4 connected and less specific",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.0").
				PrefixLength = ygot.Uint8(30)
			d.GetOrCreateInterface("eth1").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.0").
				PrefixLength = ygot.Uint8(28)
			d.GetOrCreateInterface("eth2").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.0").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inNI: "DEFAULT",
		inIP: mustCIDR("192.0.2.1/32"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
		},
	}, {
		desc: "v6 connected and less specific",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001::ab:cd:ef").
				PrefixLength = ygot.Uint8(62)
			d.GetOrCreateInterface("eth1").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001::ab:cd:ef").
				PrefixLength = ygot.Uint8(52)
			d.GetOrCreateInterface("eth2").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001::ab:cd:ef").
				PrefixLength = ygot.Uint8(42)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inNI: "DEFAULT",
		inIP: mustCIDR("2001::ab:cd:ef/128"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
		},
	}, {
		desc: "v4 ecmp",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.1").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateInterface("eth1").
				GetOrCreateSubinterface(1).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.2").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateInterface("eth2").
				GetOrCreateSubinterface(2).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.3").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inNI: "DEFAULT",
		inIP: mustCIDR("192.0.2.1/32"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
			{Name: "eth2", Subinterface: 2},
			{Name: "eth1", Subinterface: 1},
		},
	}, {
		desc: "v6 ecmp",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001::aaaa").
				PrefixLength = ygot.Uint8(42)
			d.GetOrCreateInterface("eth1").
				GetOrCreateSubinterface(1).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001::bbbb").
				PrefixLength = ygot.Uint8(42)
			d.GetOrCreateInterface("eth2").
				GetOrCreateSubinterface(2).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001::cccc").
				PrefixLength = ygot.Uint8(42)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inNI: "DEFAULT",
		inIP: mustCIDR("2001::dddd/42"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
			{Name: "eth2", Subinterface: 2},
			{Name: "eth1", Subinterface: 1},
		},
	}, {
		desc: "v4 recursive route onto connected route",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.1").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "8.8.8.8/32",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "192.0.2.1",
					NetworkInstance: "DEFAULT",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("8.8.8.8/32"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
		},
	}, {
		desc: "v6 recursive route onto connected route",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001::aaaa").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "2023::2025/128",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "2001::abcd",
					NetworkInstance: "DEFAULT",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("2023::2025/128"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
		},
	}, {
		desc: "v6 recursive route onto connected v4 route in IPv4 format",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.1").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "2023::2025/128",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "192.0.2.1",
					NetworkInstance: "DEFAULT",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("2023::2025/128"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
		},
	}, {
		desc: "v6 recursive route onto connected v4 route in IPv4-mapped IPv6 format",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.1").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "2023::2025/128",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "::ffff:192.0.2.1",
					NetworkInstance: "DEFAULT",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("2023::2025/128"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
		},
	}, {
		desc: "v4 route that recurses onto itself - not connected",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "1.1.1.1/32",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "1.1.1.1",
					NetworkInstance: "DEFAULT",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("1.1.1.1/32"),
	}, {
		desc: "v4 route that recurses onto a prefix that contains this addr",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "1.0.0.0/8",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "1.1.1.1",
					NetworkInstance: "DEFAULT",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("1.1.1.1/32"),
	}, {
		desc: "v4 route that recurses onto itself - via an intermediate route",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "1.0.0.0/8",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "2.2.2.2",
					NetworkInstance: "DEFAULT",
				}},
			}, {
				Prefix: "2.0.0.0/8",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "1.1.1.1",
					NetworkInstance: "DEFAULT",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("1.1.1.1/32"),
	}, {
		desc: "v4 route that recurses onto itself on one ECMP path",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.1").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateInterface("eth1").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.3.1").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "1.1.1.0/24",
				NextHops: []*afthelper.NextHopSummary{{
					Address: "192.0.2.1",
				}, {
					Address: "2.2.2.2",
				}},
			}, {
				Prefix: "2.0.0.0/8",
				NextHops: []*afthelper.NextHopSummary{{
					Address: "1.1.1.1",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("1.1.1.1/32"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
		},
	}, {
		desc: "v6 route that recurses onto itself on one ECMP path",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001:db8:1::").
				PrefixLength = ygot.Uint8(48)
			d.GetOrCreateInterface("eth1").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001:db8:2::").
				PrefixLength = ygot.Uint8(48)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "2001:db8:42::/48",
				NextHops: []*afthelper.NextHopSummary{{
					Address: "2001:db8:1::1",
				}, {
					Address: "2001:db8:0b::1",
				}},
			}, {
				Prefix: "2001:db8:0b::1/128",
				NextHops: []*afthelper.NextHopSummary{{
					Address: "2001:db8:42::1",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("2001:db8:42::42/128"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
		},
	}, {
		desc: "v4 recursive route onto two connected route",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.1").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateInterface("eth1").
				GetOrCreateSubinterface(42).
				GetOrCreateIpv4().
				GetOrCreateAddress("172.16.12.1").
				PrefixLength = ygot.Uint8(16)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "8.8.8.8/32",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "192.0.2.1",
					NetworkInstance: "DEFAULT",
				}, {
					Address:         "172.16.12.4",
					NetworkInstance: "DEFAULT",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("8.8.8.8/32"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
			{Name: "eth1", Subinterface: 42},
		},
	}, {
		desc: "v6 recursive route onto two connected route",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv6().
				GetOrCreateAddress("2001::aaaa").
				PrefixLength = ygot.Uint8(42)
			d.GetOrCreateInterface("eth1").
				GetOrCreateSubinterface(42).
				GetOrCreateIpv6().
				GetOrCreateAddress("4002::bbbb").
				PrefixLength = ygot.Uint8(40)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "2023::2025/128",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "2001::abcd",
					NetworkInstance: "DEFAULT",
				}, {
					Address:         "4002::def0",
					NetworkInstance: "DEFAULT",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("2023::2025/128"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
			{Name: "eth1", Subinterface: 42},
		},
	}, {
		desc: "v6 recursive route onto v4 and v6 connected routes",
		inCfg: func() *oc.Root {
			d := &oc.Root{}
			d.GetOrCreateInterface("eth0").
				GetOrCreateSubinterface(0).
				GetOrCreateIpv4().
				GetOrCreateAddress("192.0.2.1").
				PrefixLength = ygot.Uint8(24)
			d.GetOrCreateInterface("eth1").
				GetOrCreateSubinterface(42).
				GetOrCreateIpv6().
				GetOrCreateAddress("4002::bbbb").
				PrefixLength = ygot.Uint8(40)
			d.GetOrCreateNetworkInstance("DEFAULT").
				Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
			return d
		}(),
		inAddRoutes: map[string][]*Route{
			"DEFAULT": {{
				Prefix: "2023::2025/128",
				NextHops: []*afthelper.NextHopSummary{{
					Address:         "::ffff:192.0.2.1",
					NetworkInstance: "DEFAULT",
				}, {
					Address:         "4002::def0",
					NetworkInstance: "DEFAULT",
				}},
			}},
		},
		inNI: "DEFAULT",
		inIP: mustCIDR("2023::2025/128"),
		wantInterface: []*Interface{
			{Name: "eth0", Subinterface: 0},
			{Name: "eth1", Subinterface: 42},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			r, err := NewSysRIB(tt.inCfg)
			if err != nil {
				t.Fatalf("could not build RIB, %v", err)
			}
			for ni, routes := range tt.inAddRoutes {
				for _, rt := range routes {
					if err := r.AddRoute(ni, rt); err != nil {
						t.Fatalf("cannot add route %s to NI %s, err: %v", rt.Prefix, ni, err)
					}
				}
			}

			got, err := r.EgressInterface(tt.inNI, &tt.inIP)
			if (err != nil) != tt.wantErr {
				t.Fatalf("did not get expected error, got: %v, wantErr? %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(got, tt.wantInterface, cmpopts.SortSlices(func(a, b *Interface) bool {
				return a.Name < b.Name
			})); diff != "" {
				t.Fatalf("did not get expected interface set, diff(-got,+want):\n%s", diff)
			}
		})
	}
}

func baseCfg() *oc.Root {
	d := &oc.Root{}
	d.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
	return d
}

func TestAddAndDeleteRouteSequences(t *testing.T) {
	tests := [][]struct {
		desc         string
		inRoute      *Route
		inDelete     bool
		wantV4Routes []*Route
		wantV6Routes []*Route
	}{{
		{
			desc: "v4",
			inRoute: &Route{
				Prefix: "8.8.8.8/32",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			wantV4Routes: []*Route{{
				Prefix: "8.8.8.8/32",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			}},
		}, {
			desc: "v4-idempotent",
			inRoute: &Route{
				Prefix: "8.8.8.8/32",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			wantV4Routes: []*Route{{
				Prefix: "8.8.8.8/32",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			}},
		}, {
			desc: "v4-delete",
			inRoute: &Route{
				Prefix: "8.8.8.8/32",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			inDelete: true,
		}, {
			desc: "v4-delete-idempotent",
			inRoute: &Route{
				Prefix: "8.8.8.8/32",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			inDelete: true,
		},
	}, {
		{
			desc: "v4-noncanonical",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			wantV4Routes: []*Route{{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			}},
		}, {
			desc: "v4-non-canonical-delete",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			inDelete: true,
		},
	}, {
		{
			desc: "v4-noncanonical-2",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			wantV4Routes: []*Route{{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			}},
		}, {
			desc: "v4-non-canonical-2-delete-canonical",
			inRoute: &Route{
				Prefix: "8.8.8.0/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			inDelete: true,
		},
	}, {
		{
			desc: "v4-metric-test",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 0,
					Metric:        0,
				},
			},
			wantV4Routes: []*Route{{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 0,
					Metric:        0,
				},
			}},
		}, {
			desc: "v4-metric-test-add-diff-metric",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 0,
					Metric:        100,
				},
			},
			wantV4Routes: []*Route{{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 0,
					Metric:        100,
				},
			}},
		}, {
			desc: "v4-metric-test-delete-same-metric",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 0,
					Metric:        0,
				},
			},
			inDelete: true,
		},
	}, {
		{
			desc: "add-v4-admin-test",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 0,
					Metric:        0,
				},
			},
			wantV4Routes: []*Route{{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 0,
					Metric:        0,
				},
			}},
		}, {
			desc: "add-v4-diff-admin",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 100,
					Metric:        0,
				},
			},
			wantV4Routes: []*Route{{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 0,
					Metric:        0,
				},
			}, {
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 100,
					Metric:        0,
				},
			}},
		}, {
			desc: "v4-delete-admin",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 0,
					Metric:        0,
				},
			},
			inDelete: true,
			wantV4Routes: []*Route{{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 100,
					Metric:        0,
				},
			}},
		}, {
			desc: "v4-delete-admin-again",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 0,
					Metric:        0,
				},
			},
			inDelete: true,
			wantV4Routes: []*Route{{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 100,
					Metric:        0,
				},
			}},
		}, {
			desc: "v4-delete-diff-admin",
			inRoute: &Route{
				Prefix: "8.8.8.8/24",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
				RoutePref: RoutePreference{
					AdminDistance: 100,
					Metric:        0,
				},
			},
			inDelete: true,
		},
	}, {
		{
			desc: "v6",
			inRoute: &Route{
				Prefix: "2001::ab:cd:ef/42",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			wantV6Routes: []*Route{{
				Prefix: "2001::ab:cd:ef/42",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			}},
		}, {
			desc: "v6-2",
			inRoute: &Route{
				Prefix: "2023::ab:cd:ef/42",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			wantV6Routes: []*Route{{
				Prefix: "2001::ab:cd:ef/42",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			}, {
				Prefix: "2023::ab:cd:ef/42",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			}},
		},
	}, {
		{
			desc: "v6-noncanonical",
			inRoute: &Route{
				Prefix: "2001::ab:cd:ef/42",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			wantV6Routes: []*Route{{
				Prefix: "2001::ab:cd:ef/42",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			}},
		}, {
			desc: "v6-non-canonical-delete",
			inRoute: &Route{
				Prefix: "2001::ab:0/42",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			inDelete: true,
		}, {
			desc: "v6-non-canonical-delete-idempotent",
			inRoute: &Route{
				Prefix: "2001::ab:cd:ef/42",
				Connected: &Interface{
					Name:         "eth0",
					Subinterface: 0,
				},
			},
			inDelete: true,
		},
	}}

	for _, ttt := range tests {
		s, err := NewSysRIB(baseCfg())
		if err != nil {
			t.Fatalf("cannot create new system RIB, got err: %v", err)
		}
		for _, tt := range ttt {
			t.Run(tt.desc, func(t *testing.T) {
				if tt.inDelete {
					// Delete first route
					if err := s.DeleteRoute(fakedevice.DefaultNetworkInstance, tt.inRoute); err != nil {
						t.Error(err)
					}
				} else {
					// Add a new route
					if err := s.AddRoute(fakedevice.DefaultNetworkInstance, tt.inRoute); err != nil {
						t.Error(err)
					}
				}

				var gotV4Routes, gotV6Routes []*Route
				for _, ni := range s.NI {
					for it := ni.IPV4.Iterate(); it.Next(); {
						gotV4Routes = append(gotV4Routes, it.Tags()...)
					}
					for it := ni.IPV6.Iterate(); it.Next(); {
						gotV6Routes = append(gotV6Routes, it.Tags()...)
					}
				}
				if diff := cmp.Diff(tt.wantV4Routes, gotV4Routes, cmpopts.SortSlices(func(a, b *Route) bool {
					if a.Prefix < b.Prefix {
						return true
					} else if a.Prefix > b.Prefix {
						return false
					}
					return a.RoutePref.AdminDistance < b.RoutePref.AdminDistance
				})); diff != "" {
					t.Errorf("v4Routes (-want, +got):\n%s", diff)
				}
				if diff := cmp.Diff(tt.wantV6Routes, gotV6Routes, cmpopts.SortSlices(func(a, b *Route) bool {
					if a.Prefix < b.Prefix {
						return true
					} else if a.Prefix > b.Prefix {
						return false
					}
					return a.RoutePref.AdminDistance < b.RoutePref.AdminDistance
				})); diff != "" {
					t.Errorf("v6Routes (-want, +got):\n%s", diff)
				}
			})
		}
	}
}
