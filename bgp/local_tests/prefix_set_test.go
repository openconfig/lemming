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
	"testing"
	"time"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
)

func TestPrefixSet(t *testing.T) {
	dut1, stop1 := newLemming(t, dut1spec, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "192.0.2.1/31",
		niName:  "DEFAULT",
	}})
	defer stop1()
	dut2, stop2 := newLemming(t, dut2spec, nil)
	defer stop2()
	dut3, stop3 := newLemming(t, dut3spec, nil)
	defer stop3()

	establishSessionPair(t, dut1, dut2, dut1spec, dut2spec)
	establishSessionPair(t, dut2, dut3, dut2spec, dut3spec)

	// Install both prefixes into DUT2.
	// TODO: Write some constant static routes for testing and add a helper
	// for installing them.
	prefix1 := "10.33.0.0/16"
	route1 := &oc.NetworkInstance_Protocol_Static{
		Prefix: ygot.String(prefix1),
		NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
			"single": {
				Index:   ygot.String("single"),
				NextHop: oc.UnionString("192.0.2.1"),
				Recurse: ygot.Bool(true),
			},
		},
	}
	installStaticRoute(t, dut1, route1)

	prefix2 := "10.3.0.0/16"
	route2 := &oc.NetworkInstance_Protocol_Static{
		Prefix: ygot.String(prefix2),
		NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
			"single": {
				Index:   ygot.String("single"),
				NextHop: oc.UnionString("192.0.2.1"),
				Recurse: ygot.Bool(true),
			},
		},
	}
	installStaticRoute(t, dut1, route2)

	staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	v := GetAll(t, dut1, staticp.StaticAny().Config())
	t.Logf("Installed static route on %v: %s", dut1, formatYgot(v))

	v4uni := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).Bgp().Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()
	// Check both prefixes are propagated from DUT1 to Loc-RIB of DUT2.
	Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPre().Route(prefix1, 0).Prefix().State(), prefix1)
	Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPost().Route(prefix1, 0).Prefix().State(), prefix1)
	Await(t, dut2, v4uni.LocRib().Route(prefix1, oc.UnionString(dut1spec.RouterID), 0).Prefix().State(), prefix1)
	Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPre().Route(prefix2, 0).Prefix().State(), prefix2)
	Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPost().Route(prefix2, 0).Prefix().State(), prefix2)
	Await(t, dut2, v4uni.LocRib().Route(prefix2, oc.UnionString(dut1spec.RouterID), 0).Prefix().State(), prefix2)

	// Check both prefixes are propagated to AdjRib-out of DUT2.
	Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPre().Route(prefix1, 0).Prefix().State(), prefix1)
	Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPost().Route(prefix1, 0).Prefix().State(), prefix1)
	Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPre().Route(prefix2, 0).Prefix().State(), prefix2)
	Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPost().Route(prefix2, 0).Prefix().State(), prefix2)

	// Now install policy
	prefixSetName := "reject-" + prefix1
	prefix1Path := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName).Prefix(prefix1, "21..24").IpPrefix()
	Replace(t, dut2, prefix1Path.Config(), prefix1)

	policyName := "def1"
	Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Statement("stmt1").Conditions().MatchPrefixSet().PrefixSet().Config(), prefixSetName)
	Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Statement("stmt1").Conditions().MatchPrefixSet().MatchSetOptions().Config(), oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)
	Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Statement("stmt1").Actions().PolicyResult().Config(), oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)
	Replace(t, dut2, ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).Bgp().Neighbor(dut3spec.RouterID).AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).ApplyPolicy().ExportPolicy().Config(), []string{policyName})

	// Check both prefixes are propagated from DUT1 to Loc-RIB of DUT2.
	Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPre().Route(prefix1, 0).Prefix().State(), prefix1)
	Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPost().Route(prefix1, 0).Prefix().State(), prefix1)
	Await(t, dut2, v4uni.LocRib().Route(prefix1, oc.UnionString(dut1spec.RouterID), 0).Prefix().State(), prefix1)
	Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPre().Route(prefix2, 0).Prefix().State(), prefix2)
	Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPost().Route(prefix2, 0).Prefix().State(), prefix2)
	Await(t, dut2, v4uni.LocRib().Route(prefix2, oc.UnionString(dut1spec.RouterID), 0).Prefix().State(), prefix2)

	// Check only prefix2 is still remains in AdjRib-out of DUT2.
	Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPre().Route(prefix2, 0).Prefix().State(), prefix2)
	Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPost().Route(prefix2, 0).Prefix().State(), prefix2)
	// Check only prefix1 has been withdrawn from AdjRib-out of DUT2.
	const rejectTimeout = 5 * time.Second
	w := Watch(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPre().Route(prefix1, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
		_, ok := val.Val()
		return !ok
	})
	if _, ok := w.Await(t); !ok {
		t.Errorf("prefix %q was not rejected within timeout.", prefix1)
	}

	Watch(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPost().Route(prefix1, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
		_, ok := val.Val()
		return !ok
	})
	if _, ok := w.Await(t); !ok {
		t.Errorf("prefix %q was not rejected within timeout.", prefix1)
	}
}
