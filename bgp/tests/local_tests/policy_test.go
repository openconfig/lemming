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

	"github.com/openconfig/lemming/bgp"
	"github.com/openconfig/lemming/bgp/tests/proto/policyval"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

	valpb "github.com/openconfig/lemming/bgp/tests/proto/policyval"
)

const (
	debug         = true
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

	t.Run("installPolicyBeforeRoutes", func(t *testing.T) {
		testPolicyAux(t, testspec, dut1, dut2, dut3, dut4, dut5)
	})
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
	case policyval.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT:
		t.Logf("Waiting for %s to be propagated", prefix)
		Await(t, currDUT, v4uni.Neighbor(prevDUT.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.LocRib().Route(prefix, oc.UnionString(prevDUT.RouterID), 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextDUT.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, nextDUT, v4uni.Neighbor(currDUT.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	case policyval.RouteTestResult_ROUTE_TEST_RESULT_DISCARD:
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
	case policyval.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED:
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
	}
	installDefaultPolicies()

	if testspec.installPolicies != nil {
		testspec.installPolicies(t, dut1, dut2, dut3, dut4, dut5)
	}

	establishSessionPairs(t, []DevicePair{{dut1, dut2}, {dut2, dut3}, {dut4, dut5}, {dut5, dut2}}...)

	staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	v := GetAll(t, dut1, staticp.StaticAny().Config())
	t.Logf("Installed static route on %v: %s", dut1, formatYgot(v))

	for _, routeTest := range testspec.spec.RouteTests {
		testPropagation(t, routeTest, dut1, dut2, dut3)
	}
	for _, routeTest := range testspec.spec.LongerPathRouteTests {
		testPropagation(t, routeTest, dut5, dut2, dut3)
	}
}
