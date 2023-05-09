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

package local_test

import (
	"fmt"
	"testing"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
)

func installStaticRoute(t *testing.T, dut *ygnmi.Client, route *oc.NetworkInstance_Protocol_Static) {
	staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).
		Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	Replace(t, dut, staticp.Static(*route.Prefix).Config(), route)
	Await(t, dut, staticp.Static(*route.Prefix).State(), route)
}

func formatYgot(v any) string {
	str, err := ygot.Marshal7951(v, ygot.JSONIndent("  "))
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(str)
}

func TestRoutePropagation(t *testing.T) {
	dut1, dut2, stop := establishSession(t, 1111, 1112, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "192.0.2.1/31",
		niName:  "DEFAULT",
	}}, nil)
	defer stop()

	prefix := "10.10.10.0/24"
	installStaticRoute(t, dut1, &oc.NetworkInstance_Protocol_Static{
		Prefix: ygot.String(prefix),
		NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
			"single": {
				Index:   ygot.String("single"),
				NextHop: oc.UnionString("192.0.2.1"),
				Recurse: ygot.Bool(true),
			},
		},
	})

	staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	v := GetAll(t, dut1, staticp.StaticAny().Config())
	t.Logf("Installed static route: %s", formatYgot(v))

	// Check route is in Loc-RIB of dut2.
	Await(t, dut2, ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).Bgp().Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast().LocRib().Route(prefix, oc.UnionString(dut1BGP.RouterID), 0).Prefix().State(), prefix)
}
