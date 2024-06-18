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
	"github.com/openconfig/lemming/policytest"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/testing/protocmp"
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
	description         string
	routeTests          []*policytest.RouteTestCase
	skipValidateAttrSet bool // whether attr-sets are validated
	dut1IsEBGP          bool // whether DUT1 and DUT2 are in different ASes
	installPolicies     func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device)
}

// testPolicy is the helper policy integration tests can call to instantiate
// policy tests.
func testPolicy(t *testing.T, testspec *PolicyTestCase) {
	t.Helper()
	dut1AS := uint32(64500)
	if testspec.dut1IsEBGP {
		dut1AS = 64599
	}
	dut1, stop1 := newLemming(t, 1, dut1AS, []*AddIntfAction{{
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

func testPropagation(t *testing.T, route policytest.TestRoute, routeTest *policytest.RoutePathTestCase, prevDUT, currDUT, nextDUT *Device) {
	t.Helper()
	v4uni := bgp.BGPPath.Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()

	prefix := route.ReachPrefix
	// Check propagation to AdjRibOutPre for all prefixes.
	Await(t, prevDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
	Await(t, prevDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
	Await(t, currDUT, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	switch expectedResult := routeTest.ExpectedResult; expectedResult {
	case policytest.RouteAccepted:
		t.Logf("Waiting for %q (%s) to be propagated", prefix, routeTest.Description)
		Await(t, currDUT, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.LocRib().Route(prefix, oc.UnionString(prevDUT.RouterID), 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, nextDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	case policytest.RouteDiscarded:
		w := Watch(t, currDUT, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return !ok
		})
		if _, ok := w.Await(t); !ok {
			t.Errorf("prefix %q (%s) was not rejected from adj-rib-in-post of %v (neighbour %v) within timeout.", prefix, routeTest.Description, currDUT, prevDUT.ID)
			break
		}
		t.Logf("prefix %q (%s) was successfully rejected from adj-rib-in-post of %v (neighbour %v) within timeout.", prefix, routeTest.Description, currDUT, prevDUT.ID)

		// Test withdrawal in the case of InstallPolicyAfterRoutes.
		w = Watch(t, nextDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return !ok
		})
		if _, ok := w.Await(t); !ok {
			t.Errorf("prefix %q (%s) was not rejected from adj-rib-in-pre of %v (neighbour %v) within timeout.", prefix, routeTest.Description, nextDUT, currDUT.ID)
			break
		}
		t.Logf("prefix %q (%s) was successfully rejected from adj-rib-in-pre of %v (neighbour %v) within timeout.", prefix, routeTest.Description, nextDUT, currDUT.ID)
	case policytest.RouteNotPreferred:
		Await(t, currDUT, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), prefix)
		w := Watch(t, currDUT, v4uni.LocRib().Route(prefix, oc.UnionString(prevDUT.RouterID), 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return !ok
		})
		if _, ok := w.Await(t); !ok {
			t.Errorf("prefix %q with origin %q (%s) was selected into loc-rib of %v.", prefix, prevDUT.ID, routeTest.Description, currDUT)
			break
		}
		t.Logf("prefix %q with origin %q (%s) was successfully not selected into loc-rib of %v within timeout.", prefix, prevDUT.ID, routeTest.Description, currDUT)

		Await(t, currDUT, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, nextDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	default:
		t.Fatalf("Invalid or unhandled policy result: %v", expectedResult)
	}
}

func testPolicyAux(t *testing.T, testspec *PolicyTestCase, dut1, dut2, dut3, dut4, dut5 *Device) {
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

	for _, routeTest := range testspec.routeTests {
		if routeTest.Input.ReachPrefix == "" {
			continue
		}
		if routeTest.RouteTest != nil {
			// Install regular test route into DUT1.
			route := &oc.NetworkInstance_Protocol_Static{
				Prefix: ygot.String(routeTest.Input.ReachPrefix),
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
		if routeTest.LongerPathRouteTest != nil {
			// Install longer-path test route into DUT4.
			route := &oc.NetworkInstance_Protocol_Static{
				Prefix: ygot.String(routeTest.Input.ReachPrefix),
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
		if routeTest.AlternatePathRouteTest != nil {
			// Install alternate-path test route into DUT5.
			route := &oc.NetworkInstance_Protocol_Static{
				Prefix: ygot.String(routeTest.Input.ReachPrefix),
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
	}

	if debug {
		staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
		v := GetAll(t, dut1, staticp.StaticAny().Config())
		t.Logf("Installed static route on %v: %s", dut1, formatYgot(v))
	}

	for _, routeTest := range testspec.routeTests {
		if routeTest.Input.ReachPrefix == "" {
			continue
		}
		if routeTest.RouteTest != nil {
			testPropagation(t, routeTest.Input, routeTest.RouteTest, dut1, dut2, dut3)
			testCommunities(t, routeTest.Input, routeTest.RouteTest, dut1, dut2, dut3)
			if !testspec.skipValidateAttrSet {
				testAttrs(t, routeTest.Input, routeTest.RouteTest, dut1, dut2, dut3)
			}
		}

		if routeTest.LongerPathRouteTest != nil {
			testPropagation(t, routeTest.Input, routeTest.LongerPathRouteTest, dut5, dut2, dut3)
			testCommunities(t, routeTest.Input, routeTest.LongerPathRouteTest, dut5, dut2, dut3)
			if !testspec.skipValidateAttrSet {
				testAttrs(t, routeTest.Input, routeTest.LongerPathRouteTest, dut5, dut2, dut3)
			}
		}
	}
}

func testCommunities(t *testing.T, route policytest.TestRoute, routeTest *policytest.RoutePathTestCase, prevDUT, currDUT, nextDUT *Device) {
	prevCommunityMap := Lookup(t, prevDUT, bgp.BGPPath.Rib().CommunityMap().State())
	prevCommMap, _ := prevCommunityMap.Val()
	currCommunityMap := Lookup(t, currDUT, bgp.BGPPath.Rib().CommunityMap().State())
	currCommMap, _ := currCommunityMap.Val()
	nextCommunityMap := Lookup(t, nextDUT, bgp.BGPPath.Rib().CommunityMap().State())
	nextCommMap, _ := nextCommunityMap.Val()
	v4uni := bgp.BGPPath.Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()

	prefix := route.ReachPrefix

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

// awaitNoDiff periodically checks diffFunc's output until empty or timeout.
//
//   - updateRIBFunc is called periodically to update any RIB data used in the
//     validation.
func awaitNoDiff(diffFunc func() string, updateRIBFunc func()) string {
	timeout := time.After(awaitTimeLimit)
	updateTicker := time.NewTicker(5 * time.Second)
	defer updateTicker.Stop()

	var diff string

	for {
		select {
		case <-timeout:
			return diff
		case <-updateTicker.C:
			updateRIBFunc()
		default:
			diff = diffFunc()
			if diff == "" {
				return ""
			}
			time.Sleep(time.Second)
		}
	}
}

// setDefaultAttrs populates the default attribute fields of the provided AttrSet.
//
// If the route is rejected, the provided AttrSet is returned unchanged.
func setDefaultAttrs(a *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet, rejected bool) *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet {
	if rejected {
		return a
	}
	if a == nil {
		a = &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{}
	}
	// Set default values
	if a.Origin == oc.BgpTypes_BgpOriginAttrType_UNSET {
		a.Origin = oc.BgpTypes_BgpOriginAttrType_IGP
	}
	if a.Med == nil {
		a.Med = ygot.Uint32(0)
	}
	if a.LocalPref == nil {
		a.LocalPref = ygot.Uint32(100)
	}
	return a
}

func testAttrs(t *testing.T, route policytest.TestRoute, routeTest *policytest.RoutePathTestCase, prevDUT, currDUT, nextDUT *Device) {
	prevAttrSetMap := Lookup(t, prevDUT, bgp.BGPPath.Rib().AttrSetMap().State())
	prevAttrMap, _ := prevAttrSetMap.Val()
	currAttrSetMap := Lookup(t, currDUT, bgp.BGPPath.Rib().AttrSetMap().State())
	currAttrMap, _ := currAttrSetMap.Val()
	nextAttrSetMap := Lookup(t, nextDUT, bgp.BGPPath.Rib().AttrSetMap().State())
	nextAttrMap, _ := nextAttrSetMap.Val()
	updateAttrMaps := func() {
		prevAttrSetMap = Lookup(t, prevDUT, bgp.BGPPath.Rib().AttrSetMap().State())
		prevAttrMap, _ = prevAttrSetMap.Val()
		currAttrSetMap = Lookup(t, currDUT, bgp.BGPPath.Rib().AttrSetMap().State())
		currAttrMap, _ = currAttrSetMap.Val()
		nextAttrSetMap = Lookup(t, nextDUT, bgp.BGPPath.Rib().AttrSetMap().State())
		nextAttrMap, _ = nextAttrSetMap.Val()
	}

	v4uni := bgp.BGPPath.Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()

	prefix := route.ReachPrefix

	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, prevDUT, prevAttrMap, v4uni.Neighbor(currDUT.RouterID).AdjRibOutPre().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(setDefaultAttrs(routeTest.PrevAdjRibOutPreAttrs, false), attrs, protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Errorf("DUT %v AdjRibOutPre attribute difference (prefix %s) (-want, +got):\n%s", prevDUT.ID, prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, prevDUT, prevAttrMap, v4uni.Neighbor(currDUT.RouterID).AdjRibOutPost().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(setDefaultAttrs(routeTest.PrevAdjRibOutPostAttrs, false), attrs, protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Errorf("DUT %v AdjRibOutPost attribute difference (prefix %s) (-want, +got):\n%s", prevDUT.ID, prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, currDUT, currAttrMap, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPre().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(setDefaultAttrs(routeTest.AdjRibInPreAttrs, false), attrs, protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Errorf("DUT %v AdjRibInPre attribute difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, currDUT, currAttrMap, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPost().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(setDefaultAttrs(routeTest.AdjRibInPostAttrs, routeTest.ExpectedResult == policytest.RouteDiscarded), attrs, protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Errorf("DUT %v AdjRibInPost attribute difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, currDUT, currAttrMap, v4uni.LocRib().Route(prefix, oc.UnionString(prevDUT.RouterID), 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(setDefaultAttrs(routeTest.LocalRibAttrs, routeTest.ExpectedResult != policytest.RouteAccepted), attrs, protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Errorf("DUT %v LocRib routeTest difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, currDUT, currAttrMap, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPre().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(setDefaultAttrs(routeTest.AdjRibOutPreAttrs, routeTest.ExpectedResult == policytest.RouteDiscarded), attrs, protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Errorf("DUT %v AdjRibOutPre attribute difference (prefix %s) (-want, +got):\n%s\n%+v", currDUT.ID, prefix, diff, routeTest.AdjRibOutPreAttrs)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, currDUT, currAttrMap, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPost().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(setDefaultAttrs(routeTest.AdjRibOutPostAttrs, routeTest.ExpectedResult == policytest.RouteDiscarded), attrs, protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Errorf("DUT %v AdjRibOutPost attribute difference (prefix %s) (-want, +got):\n%s", currDUT.ID, prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, nextDUT, nextAttrMap, v4uni.Neighbor(currDUT.RouterID).AdjRibInPre().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(setDefaultAttrs(routeTest.NextAdjRibInPreAttrs, routeTest.ExpectedResult == policytest.RouteDiscarded), attrs, protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Errorf("DUT %v AdjRibInPre attribute difference (prefix %s) (-want, +got):\n%s", nextDUT.ID, prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, nextDUT, nextAttrMap, v4uni.LocRib().Route(prefix, oc.UnionString(currDUT.RouterID), 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(setDefaultAttrs(routeTest.NextLocalRibAttrs, routeTest.ExpectedResult == policytest.RouteDiscarded), attrs, protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Errorf("DUT %v LocRib attribute difference (prefix %s) (-want, +got):\n%s", nextDUT.ID, prefix, diff)
	}
}

// getAttrs gets the attribute of the given route query to a attr-set index.
//
// For optional attributes that have defaults (e.g. local-pref and med),
// they're automatically populated if not available, since GoBGP inconsistently
// populates them.
//
// If the attr-set index doesn't exist (e.g. the route doesn't exist), nil is returned.
func getAttrs(t *testing.T, dut *Device, attrSetMap map[uint64]*oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet, query ygnmi.SingletonQuery[uint64]) (*oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet, error) {
	attrIndexVal := Lookup(t, dut, query)
	attrIndex, ok := attrIndexVal.Val()
	if !ok {
		return nil, nil
	}
	attrs, ok := attrSetMap[attrIndex]
	if !ok {
		return nil, fmt.Errorf("RIB attributes does not have expected attribute index: %v", attrIndex)
	}
	attrs.Index = nil
	// Set default values.
	if attrs.Med == nil {
		attrs.Med = ygot.Uint32(0)
	}
	if attrs.LocalPref == nil {
		attrs.LocalPref = ygot.Uint32(100)
	}
	return attrs, nil
}
