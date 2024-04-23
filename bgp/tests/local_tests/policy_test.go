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
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/lemming/bgp"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/testing/protocmp"

	valpb "github.com/openconfig/lemming/proto/policyval"
)

const (
	debug         = false
	rejectTimeout = 20 * time.Second
)

// PolicyTestCase contains the specifications for a single policy test.
//
// Topology:
//
//	DUT1 (AS 64500) -> DUT2 (AS 64500) -> DUT3 (AS 64501)
//	                    ^
//	                    |
//	DUT4 (AS 64502) -> DUT5 (AS 64500)
//
// Currently by convention, all policies are installed on DUT1 (export), DUT5
// (export), and DUT2 (import). This is because GoBGP only withdraws routes on
// import policy change after a soft reset:
// https://github.com/osrg/gobgp/blob/master/docs/sources/policy.md#policy-and-soft-reset
type PolicyTestCase struct {
	spec            *valpb.PolicyTestCase
	installPolicies func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device)
}

// testPolicy is the helper policy integration tests can call to instantiate
// policy tests.
func testPolicy(t *testing.T, testspec PolicyTestCase) {
	t.Helper()
	dut1, stop1 := newLemming(t, 1, 64500, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "192.0.2.1/31",
		niName:  "DEFAULT",
	}})
	defer stop1()
	dut2, stop2 := newLemming(t, 2, 64500, nil)
	defer stop2()
	dut3, stop3 := newLemming(t, 3, 64501, nil)
	defer stop3()
	dut4, stop4 := newLemming(t, 4, 64502, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "192.0.2.1/30",
		niName:  "DEFAULT",
	}})
	defer stop4()
	dut5, stop5 := newLemming(t, 5, 64500, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "193.0.2.1/30",
		niName:  "DEFAULT",
	}})
	defer stop5()

	testPolicyAux(t, testspec, dut1, dut2, dut3, dut4, dut5)
}

func testPropagation(t *testing.T, routeTest *valpb.RouteTestCase, prevDUT, currDUT, nextDUT *Device) {
	t.Helper()
	v4uni := bgp.BGPPath.Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()

	prefix := routeTest.GetInput().GetReachPrefix()
	// Check propagation to AdjRibOutPre for all prefixes.
	Await(t, prevDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
	Await(t, prevDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
	Await(t, currDUT, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	switch expectedResult := routeTest.GetExpectedResult(); expectedResult {
	case valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT:
		t.Logf("Waiting for %s to be propagated", prefix)
		Await(t, currDUT, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.LocRib().Route(prefix, oc.UnionString(prevDUT.RouterID), 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, nextDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	case valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD:
		w := Watch(t, currDUT, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return !ok
		})
		if _, ok := w.Await(t); !ok {
			t.Errorf("prefix %q (%s) was not rejected from adj-rib-in-post of %v (neighbour %v) within timeout.", prefix, routeTest.GetDescription(), currDUT, prevDUT.ID)
			break
		}
		t.Logf("prefix %q (%s) was successfully rejected from adj-rib-in-post of %v (neighbour %v) within timeout.", prefix, routeTest.GetDescription(), currDUT, prevDUT.ID)

		// Test withdrawal in the case of InstallPolicyAfterRoutes.
		w = Watch(t, nextDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return !ok
		})
		if _, ok := w.Await(t); !ok {
			t.Errorf("prefix %q (%s) was not rejected from adj-rib-in-pre of %v (neighbour %v) within timeout.", prefix, routeTest.GetDescription(), nextDUT, currDUT.ID)
			break
		}
		t.Logf("prefix %q (%s) was successfully rejected from adj-rib-in-pre of %v (neighbour %v) within timeout.", prefix, routeTest.GetDescription(), nextDUT, currDUT.ID)
	case valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED:
		Await(t, currDUT, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), prefix)
		w := Watch(t, currDUT, v4uni.LocRib().Route(prefix, oc.UnionString(prevDUT.RouterID), 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return !ok
		})
		if _, ok := w.Await(t); !ok {
			t.Errorf("prefix %q with origin %q (%s) was selected into loc-rib of %v.", prefix, prevDUT.ID, routeTest.GetDescription(), currDUT)
			break
		}
		t.Logf("prefix %q with origin %q (%s) was successfully not selected into loc-rib of %v within timeout.", prefix, prevDUT.ID, routeTest.GetDescription(), currDUT)

		Await(t, currDUT, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, nextDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	default:
		t.Fatalf("Invalid or unhandled policy result: %v", expectedResult)
	}
}

func testPolicyAux(t *testing.T, testspec PolicyTestCase, dut1, dut2, dut3, dut4, dut5 *Device) {
	// Remove any existing BGP config
	//
	// TODO(wenbli): Debug why sometimes this causes GoBGP to transiently
	// report "not found policy" error.
	Delete(t, dut1, bgp.BGPPath.Config())
	Delete(t, dut2, bgp.BGPPath.Config())
	Delete(t, dut3, bgp.BGPPath.Config())
	Delete(t, dut4, bgp.BGPPath.Config())
	Delete(t, dut5, bgp.BGPPath.Config())
	Delete(t, dut1, bgp.RoutingPolicyPath.Config())
	Delete(t, dut2, bgp.RoutingPolicyPath.Config())
	Delete(t, dut3, bgp.RoutingPolicyPath.Config())
	Delete(t, dut4, bgp.RoutingPolicyPath.Config())
	Delete(t, dut5, bgp.RoutingPolicyPath.Config())

	establishSessionPairs(t, []DevicePair{{dut1, dut2}, {dut2, dut3}, {dut4, dut5}, {dut5, dut2}}...)

	installDefaultPolicies := func() {
		// Clear the path for routes to be propagated.
		// DUT1 -> DUT2 -> DUT3
		Replace(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut3.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut3, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)

		// This is an alternate source of routes towards DUT2 and thereby DUT3.
		// Note that this path is longer than the above path:
		// DUT4 -> DUT5 -> DUT2 (-> DUT3)
		Replace(t, dut4, bgp.BGPPath.Neighbor(dut5.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut5, bgp.BGPPath.Neighbor(dut4.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut5, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut5.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)

		// Wait until policies are installed.
		Await(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultExportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().DefaultImportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut2, bgp.BGPPath.Neighbor(dut3.RouterID).ApplyPolicy().DefaultExportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut3, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultImportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)

		Await(t, dut4, bgp.BGPPath.Neighbor(dut5.RouterID).ApplyPolicy().DefaultExportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut5, bgp.BGPPath.Neighbor(dut4.RouterID).ApplyPolicy().DefaultImportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut5, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultExportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut2, bgp.BGPPath.Neighbor(dut5.RouterID).ApplyPolicy().DefaultImportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
	}
	installDefaultPolicies()

	if testspec.installPolicies != nil {
		testspec.installPolicies(t, dut1, dut2, dut3, dut4, dut5)
	}

	for _, routeTest := range testspec.spec.RouteTests {
		// Install all regular test routes into DUT1.
		route := &oc.NetworkInstance_Protocol_Static{
			Prefix: ygot.String(routeTest.GetInput().GetReachPrefix()),
			NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
				"single": {
					Index:   ygot.String("single"),
					NextHop: oc.UnionString("192.0.2.1"),
					Recurse: ygot.Bool(true),
				},
			},
		}
		installStaticRoute(t, dut1, route)
	}

	for _, routeTest := range testspec.spec.LongerPathRouteTests {
		// Install all longer-path test routes into DUT4.
		route := &oc.NetworkInstance_Protocol_Static{
			Prefix: ygot.String(routeTest.GetInput().GetReachPrefix()),
			NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
				"single": {
					Index:   ygot.String("single"),
					NextHop: oc.UnionString("192.0.2.1"),
					Recurse: ygot.Bool(true),
				},
			},
		}
		installStaticRoute(t, dut4, route)
	}

	for _, routeTest := range testspec.spec.AlternatePathRouteTests {
		// Install all alternate-path test routes into DUT5.
		route := &oc.NetworkInstance_Protocol_Static{
			Prefix: ygot.String(routeTest.GetInput().GetReachPrefix()),
			NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
				"single": {
					Index:   ygot.String("single"),
					NextHop: oc.UnionString("193.0.2.1"),
					Recurse: ygot.Bool(true),
				},
			},
		}
		installStaticRoute(t, dut5, route)
	}

	if debug {
		staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
		v := GetAll(t, dut1, staticp.StaticAny().Config())
		t.Logf("Installed static route on %v: %s", dut1, formatYgot(v))
	}

	for _, routeTest := range testspec.spec.RouteTests {
		testPropagation(t, routeTest, dut1, dut2, dut3)
		testCommunities(t, routeTest, dut1, dut2, dut3)
		testAttrs(t, routeTest, dut1, dut2, dut3)
	}
	for _, routeTest := range testspec.spec.LongerPathRouteTests {
		testPropagation(t, routeTest, dut5, dut2, dut3)
		testCommunities(t, routeTest, dut5, dut2, dut3)
		testAttrs(t, routeTest, dut5, dut2, dut3)
	}
}

func testCommunities(t *testing.T, routeTest *valpb.RouteTestCase, prevDUT, currDUT, nextDUT *Device) {
	prevCommunityMap := Lookup(t, prevDUT, bgp.BGPPath.Rib().CommunityMap().State())
	prevCommMap, _ := prevCommunityMap.Val()
	currCommunityMap := Lookup(t, currDUT, bgp.BGPPath.Rib().CommunityMap().State())
	currCommMap, _ := currCommunityMap.Val()
	nextCommunityMap := Lookup(t, nextDUT, bgp.BGPPath.Rib().CommunityMap().State())
	nextCommMap, _ := nextCommunityMap.Val()
	v4uni := bgp.BGPPath.Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()

	prefix := routeTest.GetInput().GetReachPrefix()

	if diff := cmp.Diff(routeTest.PrevAdjRibOutPreCommunities, getCommunities(t, prevDUT, prevCommMap, v4uni.Neighbor(currDUT.RouterID).AdjRibOutPre().Route(prefix, 0).CommunityIndex().State())); diff != "" {
		t.Errorf("DUT %v AdjRibOutPre communities difference (prefix %s) (-want, +got):\n%s", prevDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(routeTest.PrevAdjRibOutPostCommunities, getCommunities(t, prevDUT, prevCommMap, v4uni.Neighbor(currDUT.RouterID).AdjRibOutPost().Route(prefix, 0).CommunityIndex().State())); diff != "" {
		t.Errorf("DUT %v AdjRibOutPost communities difference (prefix %s) (-want, +got):\n%s", prevDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(routeTest.AdjRibInPreCommunities, getCommunities(t, currDUT, currCommMap, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPre().Route(prefix, 0).CommunityIndex().State())); diff != "" {
		t.Errorf("DUT %v AdjRibInPre communities difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(routeTest.AdjRibInPostCommunities, getCommunities(t, currDUT, currCommMap, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPost().Route(prefix, 0).CommunityIndex().State())); diff != "" {
		t.Errorf("DUT %v AdjRibInPost communities difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(routeTest.LocalRibCommunities, getCommunities(t, currDUT, currCommMap, v4uni.LocRib().Route(prefix, oc.UnionString(prevDUT.RouterID), 0).CommunityIndex().State())); diff != "" {
		t.Errorf("DUT %v LocRib communities difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(routeTest.AdjRibOutPreCommunities, getCommunities(t, currDUT, currCommMap, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPre().Route(prefix, 0).CommunityIndex().State())); diff != "" {
		t.Errorf("DUT %v AdjRibOutPre communities difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(routeTest.AdjRibOutPostCommunities, getCommunities(t, currDUT, currCommMap, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPost().Route(prefix, 0).CommunityIndex().State())); diff != "" {
		t.Errorf("DUT %v AdjRibOutPost communities difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(routeTest.NextAdjRibInPreCommunities, getCommunities(t, nextDUT, nextCommMap, v4uni.Neighbor(currDUT.RouterID).AdjRibInPre().Route(prefix, 0).CommunityIndex().State())); diff != "" {
		t.Errorf("DUT %v AdjRibInPre communities difference (prefix %s) (-want, +got):\n%s", nextDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(routeTest.NextLocalRibCommunities, getCommunities(t, nextDUT, nextCommMap, v4uni.LocRib().Route(prefix, oc.UnionString(currDUT.RouterID), 0).CommunityIndex().State())); diff != "" {
		t.Errorf("DUT %v LocRib communities difference (prefix %s) (-want, +got):\n%s", nextDUT.ID, prefix, diff)
	}
}

// getCommunities gets the communities of the given route query to a community index.
//
// If the community index doesn't exist (e.g. the route doesn't exist), nil is returned.
func getCommunities(t *testing.T, dut *Device, commMap map[uint64]*oc.NetworkInstance_Protocol_Bgp_Rib_Community, query ygnmi.SingletonQuery[uint64]) []string {
	commIndexVal := Lookup(t, dut, query)
	commIndex, ok := commIndexVal.Val()
	if !ok {
		return nil
	}
	comms, ok := commMap[commIndex]
	if !ok {
		t.Errorf("RIB communities does not have expected community index: %v", commIndex)
		return nil
	}
	var gotCommunities []string
	for _, comm := range comms.GetCommunity() {
		switch c := bgp.ConvertCommunity(comm); c {
		case "":
			t.Errorf("Unexpected community type: (%T, %v)", c, c)
		default:
			gotCommunities = append(gotCommunities, c)
		}
	}
	return gotCommunities
}

func testAttrs(t *testing.T, routeTest *valpb.RouteTestCase, prevDUT, currDUT, nextDUT *Device) {
	prevAttrSetMap := Lookup(t, prevDUT, bgp.BGPPath.Rib().AttrSetMap().State())
	prevAttrMap, _ := prevAttrSetMap.Val()
	currAttrSetMap := Lookup(t, currDUT, bgp.BGPPath.Rib().AttrSetMap().State())
	currAttrMap, _ := currAttrSetMap.Val()
	nextAttrSetMap := Lookup(t, nextDUT, bgp.BGPPath.Rib().AttrSetMap().State())
	nextAttrMap, _ := nextAttrSetMap.Val()
	v4uni := bgp.BGPPath.Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()

	prefix := routeTest.GetInput().GetReachPrefix()

	igpAttr := &valpb.RibAttributes{
		AttrSet: &valpb.AttrSet{
			Origin: "IGP",
		},
	}

	// NOTE: GoBGP doesn't seem to set origin properly -- it is always set to zero.
	//       So, for simplicity just always test if it's "IGP" instead of
	//       specifying confusing test cases where the origin should not be IGP.
	// If this is ever supported by GoBGP properly, OR if we decide to use
	// SetOrigin to artificially set this attribute, then remove this
	// function and associated logic.
	nonNilOrIGP := func(a *valpb.RibAttributes, rejected bool) *valpb.RibAttributes {
		if a == nil && !rejected {
			return igpAttr
		}
		return a
	}

	if diff := cmp.Diff(nonNilOrIGP(routeTest.PrevAdjRibOutPreAttrs, false), getAttrs(t, prevDUT, prevAttrMap, v4uni.Neighbor(currDUT.RouterID).AdjRibOutPre().Route(prefix, 0).AttrIndex().State()), protocmp.Transform()); diff != "" {
		t.Errorf("DUT %v AdjRibOutPre attribute difference (prefix %s) (-want, +got):\n%s", prevDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(nonNilOrIGP(routeTest.PrevAdjRibOutPostAttrs, false), getAttrs(t, prevDUT, prevAttrMap, v4uni.Neighbor(currDUT.RouterID).AdjRibOutPost().Route(prefix, 0).AttrIndex().State()), protocmp.Transform()); diff != "" {
		t.Errorf("DUT %v AdjRibOutPost attribute difference (prefix %s) (-want, +got):\n%s", prevDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(nonNilOrIGP(routeTest.AdjRibInPreAttrs, false), getAttrs(t, currDUT, currAttrMap, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPre().Route(prefix, 0).AttrIndex().State()), protocmp.Transform()); diff != "" {
		t.Errorf("DUT %v AdjRibInPre attribute difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(nonNilOrIGP(routeTest.AdjRibInPostAttrs, routeTest.ExpectedResult == valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD), getAttrs(t, currDUT, currAttrMap, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPost().Route(prefix, 0).AttrIndex().State()), protocmp.Transform()); diff != "" {
		t.Errorf("DUT %v AdjRibInPost attribute difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(nonNilOrIGP(routeTest.LocalRibAttrs, routeTest.ExpectedResult != valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT), getAttrs(t, currDUT, currAttrMap, v4uni.LocRib().Route(prefix, oc.UnionString(prevDUT.RouterID), 0).AttrIndex().State()), protocmp.Transform()); diff != "" {
		t.Errorf("DUT %v LocRib attrs difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(nonNilOrIGP(routeTest.AdjRibOutPreAttrs, routeTest.ExpectedResult == valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD), getAttrs(t, currDUT, currAttrMap, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPre().Route(prefix, 0).AttrIndex().State()), protocmp.Transform()); diff != "" {
		fmt.Println("<1>", prefix, routeTest.ExpectedResult, valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT)
		t.Errorf("DUT %v AdjRibOutPre attribute difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(nonNilOrIGP(routeTest.AdjRibOutPostAttrs, routeTest.ExpectedResult == valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD), getAttrs(t, currDUT, currAttrMap, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPost().Route(prefix, 0).AttrIndex().State()), protocmp.Transform()); diff != "" {
		t.Errorf("DUT %v AdjRibOutPost attribute difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(nonNilOrIGP(routeTest.NextAdjRibInPreAttrs, routeTest.ExpectedResult == valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD), getAttrs(t, nextDUT, nextAttrMap, v4uni.Neighbor(currDUT.RouterID).AdjRibInPre().Route(prefix, 0).AttrIndex().State()), protocmp.Transform()); diff != "" {
		t.Errorf("DUT %v AdjRibInPre attribute difference (prefix %s) (-want, +got):\n%s", nextDUT.ID, prefix, diff)
	}
	if diff := cmp.Diff(nonNilOrIGP(routeTest.NextLocalRibAttrs, routeTest.ExpectedResult == valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD), getAttrs(t, nextDUT, nextAttrMap, v4uni.LocRib().Route(prefix, oc.UnionString(currDUT.RouterID), 0).AttrIndex().State()), protocmp.Transform()); diff != "" {
		t.Errorf("DUT %v LocRib attribute difference (prefix %s) (-want, +got):\n%s", nextDUT.ID, prefix, diff)
	}
}

// getAttrs gets the attribute of the given route query to a attr-set index.
//
// If the attr-set index doesn't exist (e.g. the route doesn't exist), nil is returned.
func getAttrs(t *testing.T, dut *Device, attrSetMap map[uint64]*oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet, query ygnmi.SingletonQuery[uint64]) *valpb.RibAttributes {
	attrIndexVal := Lookup(t, dut, query)
	attrIndex, ok := attrIndexVal.Val()
	if !ok {
		return nil
	}
	attrs, ok := attrSetMap[attrIndex]
	if !ok {
		t.Errorf("RIB attributes does not have expected attribute index: %v", attrIndex)
		return nil
	}

	gotAttrs := &valpb.RibAttributes{AttrSet: &valpb.AttrSet{}}
	if origin := attrs.GetOrigin(); origin != oc.BgpTypes_BgpOriginAttrType_UNSET {
		gotAttrs.AttrSet.Origin = origin.String()
	}
	if attrs.Med != nil {
		gotAttrs.AttrSet.Med = attrs.GetMed()
	}
	return gotAttrs
}
