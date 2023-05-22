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

	"github.com/openconfig/lemming/bgp"
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
	dut1, stop1 := newLemming(t, dut1spec, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "127.0.0.0/30",
		niName:  "DEFAULT",
	}, {
		name:    "eth1",
		ifindex: 1,
		enabled: true,
		prefix:  "192.0.2.1/31",
		niName:  "DEFAULT",
	}})
	defer stop1()
	dut2, stop2 := newLemming(t, dut2spec, nil)
	defer stop2()
	dut3, stop3 := newLemming(t, dut3spec, nil)
	defer stop3()
	dut4, stop4 := newLemming(t, dut4spec, nil)
	defer stop4()

	installDefaultPolicies := func() {
		// TODO: clean-up
		// This print statement is to help with debugging a non-deterministic "update is stale" error in the gNMI cache.
		// See https://pipelines.actions.githubusercontent.com/serviceHosts/8897a2bd-3f38-443c-8744-4fc2e0f58aae/_apis/pipelines/1/runs/2930/signedlogcontent/4?urlExpires=2023-05-22T21%3A31%3A09.1265307Z&urlSigningMethod=HMACV1&urlSignature=lYO37gwTl6Vtt%2FCc2Z%2Faeffh32zXn1cszkosj69uGSM%3D
		// route_propagation_test.go:68: Replace(t) on ygnmi client (target: "dut1") at /network-instances/network-instance[name=DEFAULT]/protocols/protocol[identifier=BGP][name=BGP]/bgp/neighbors/neighbor[neighbor-address=127.0.0.2]/apply-policy/config/default-export-policy: Replace(t) at path origin:"openconfig" elem:{name:"network-instances"} elem:{name:"network-instance" key:{key:"name" value:"DEFAULT"}} elem:{name:"protocols"} elem:{name:"protocol" key:{key:"identifier" value:"BGP"} key:{key:"name" value:"BGP"}} elem:{name:"bgp"} elem:{name:"neighbors"} elem:{name:"neighbor" key:{key:"neighbor-address" value:"127.0.0.2"}} elem:{name:"apply-policy"} elem:{name:"config"} elem:{name:"default-export-policy"}: rpc error: code = Internal desc = update is stale, update is stale, update is stale
		fmt.Println("Installing default policies")
		// Clear the path for routes to be propagated.
		Replace(t, dut1, bgp.BGPPath.Neighbor(dut2spec.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut1spec.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut3spec.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut3, bgp.BGPPath.Neighbor(dut2spec.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut3, bgp.BGPPath.Neighbor(dut4spec.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut4, bgp.BGPPath.Neighbor(dut3spec.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
	}
	installDefaultPolicies()

	establishSessionPair(t, dut1, dut2, dut1spec, dut2spec)
	establishSessionPair(t, dut2, dut3, dut2spec, dut3spec)
	establishSessionPair(t, dut3, dut4, dut3spec, dut4spec)

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

	v4uni := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).Bgp().Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()
	// Check route is in Adj-In of dut2.
	Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), prefix)
	// Check route is in Loc-RIB of dut2.
	Await(t, dut2, v4uni.LocRib().Route(prefix, oc.UnionString(dut1spec.RouterID), 0).Prefix().State(), prefix)
	// Check route is in Adj-Out of dut2.
	Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
	Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)

	// Check route is in Adj-In of dut3.
	Await(t, dut3, v4uni.Neighbor(dut2spec.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)

	// Check route is in Adj-In of dut4.
	Await(t, dut4, v4uni.Neighbor(dut3spec.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
}
