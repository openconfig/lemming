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
	rejectTimeout = 5 * time.Second
)

// PolicyTestCase contains the specifications for a single policy test.
//
// Topology:
//
//	DUT1 -> DUT2 -> DUT3
//	         ^
//	         |
//	DUT4 -> DUT5
//
// Currently, all policies are installed on DUT2, and by convention,
// the import policies set attributes, and export policy to DUT3 filters
// prefixes.
type PolicyTestCase struct {
	spec            *valpb.PolicyTestCase
	installPolicies func(t *testing.T, dut2 *ygnmi.Client)
}

// testPolicy is the helper policy integration tests can call to instantiate
// policy tests.
func testPolicy(t *testing.T, testspec PolicyTestCase) {
	t.Helper()
	t.Run("installPolicyBeforeRoutes", func(t *testing.T) {
		testPolicyAux(t, testspec, false)
	})

	t.Run("installPolicyAfterRoutes", func(t *testing.T) {
		testPolicyAux(t, testspec, true)
	})
}

func testPropagation(t *testing.T, routeTest *valpb.RouteTestCase, prevDUT, currDUT, nextDUT *ygnmi.Client, prevRouterID, currRouterID, nextRouterID string, filterPoliciesInstalled bool) {
	t.Helper()
	v4uni := bgp.BGPPath.Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()
	expectedResult := routeTest.GetExpectedResultBeforePolicy()
	if filterPoliciesInstalled {
		expectedResult = routeTest.GetExpectedResult()
	}

	prefix := routeTest.GetInput().GetReachPrefix()
	// Check propagation to AdjRibOutPre for all prefixes.
	Await(t, prevDUT, v4uni.Neighbor(currRouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
	Await(t, prevDUT, v4uni.Neighbor(currRouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
	Await(t, currDUT, v4uni.Neighbor(prevRouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	Await(t, currDUT, v4uni.Neighbor(prevRouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), prefix)
	switch expectedResult {
	case policyval.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT:
		Await(t, currDUT, v4uni.LocRib().Route(prefix, oc.UnionString(prevRouterID), 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextRouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		t.Logf("Waiting for %s to be propagated", prefix)
		Await(t, currDUT, v4uni.Neighbor(nextRouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, nextDUT, v4uni.Neighbor(currRouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	case policyval.RouteTestResult_ROUTE_TEST_RESULT_DISCARD:
		Await(t, currDUT, v4uni.LocRib().Route(prefix, oc.UnionString(prevRouterID), 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextRouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		w := Watch(t, currDUT, v4uni.Neighbor(nextRouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return ok
		})
		if _, ok := w.Await(t); ok {
			t.Errorf("prefix %q (%s) was not rejected from adj-rib-out-post of %v within timeout.", prefix, routeTest.GetDescription(), currDUT)
			break
		}
		t.Logf("prefix %q (%s) was successfully rejected from adj-rib-out-post of %v within timeout.", prefix, routeTest.GetDescription(), currDUT)

		w = Watch(t, nextDUT, v4uni.Neighbor(currRouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return ok
		})
		if _, ok := w.Await(t); ok {
			t.Errorf("prefix %q (%s) was not rejected from adj-rib-in-pre of %v within timeout.", prefix, routeTest.GetDescription(), nextDUT)
			break
		}
		t.Logf("prefix %q (%s) was successfully rejected from adj-rib-in-pre of %v within timeout.", prefix, routeTest.GetDescription(), nextDUT)
	case policyval.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED:
		w := Watch(t, currDUT, v4uni.LocRib().Route(prefix, oc.UnionString(prevRouterID), 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return ok
		})
		if _, ok := w.Await(t); ok {
			t.Errorf("prefix %q with origin %q (%s) was selected into loc-rib of %v.", prefix, prevRouterID, routeTest.GetDescription(), currDUT)
			break
		}
		t.Logf("prefix %q with origin %q (%s) was successfully not selected into loc-rib of %v within timeout.", prefix, prevRouterID, routeTest.GetDescription(), currDUT)

		Await(t, currDUT, v4uni.Neighbor(nextRouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, currDUT, v4uni.Neighbor(nextRouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, nextDUT, v4uni.Neighbor(currRouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
	default:
		t.Fatalf("Invalid or unhandled policy result: %v", expectedResult)
	}
}

func testPolicyAux(t *testing.T, testspec PolicyTestCase, installPolicyAfterRoutes bool) {
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
	dut4, stop4 := newLemming(t, dut4spec, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "192.0.2.1/30",
		niName:  "DEFAULT",
	}})
	defer stop4()
	dut5, stop5 := newLemming(t, dut5spec, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "193.0.2.1/30",
		niName:  "DEFAULT",
	}})
	defer stop5()

	installDefaultPolicies := func() {
		// Clear the path for routes to be propagated.
		// DUT1 -> DUT2 -> DUT3
		Replace(t, dut1, bgp.BGPPath.Neighbor(dut2spec.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut1spec.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut3spec.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut3, bgp.BGPPath.Neighbor(dut2spec.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)

		// This is an alternate source of routes towards DUT2 and thereby DUT3.
		// Note that this path is longer than the above path:
		// DUT4 -> DUT5 -> DUT2 (-> DUT3)
		Replace(t, dut4, bgp.BGPPath.Neighbor(dut5spec.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut5, bgp.BGPPath.Neighbor(dut4spec.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut5, bgp.BGPPath.Neighbor(dut2spec.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut5spec.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
	}
	installDefaultPolicies()

	if testspec.installPolicies != nil && !installPolicyAfterRoutes {
		testspec.installPolicies(t, dut2)
	}

	establishSessionPair(t, dut1, dut2, dut1spec, dut2spec)
	establishSessionPair(t, dut2, dut3, dut2spec, dut3spec)
	establishSessionPair(t, dut4, dut5, dut4spec, dut5spec)
	establishSessionPair(t, dut5, dut2, dut5spec, dut2spec)

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

	staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	v := GetAll(t, dut1, staticp.StaticAny().Config())
	t.Logf("Installed static route on %v: %s", dut1, formatYgot(v))

	for _, routeTest := range testspec.spec.RouteTests {
		testPropagation(t, routeTest, dut1, dut2, dut3, dut1spec.RouterID, dut2spec.RouterID, dut3spec.RouterID, !installPolicyAfterRoutes)
	}
	for _, routeTest := range testspec.spec.LongerPathRouteTests {
		testPropagation(t, routeTest, dut5, dut2, dut3, dut5spec.RouterID, dut2spec.RouterID, dut3spec.RouterID, !installPolicyAfterRoutes)
	}

	if installPolicyAfterRoutes {
		awaitNewSession := make(chan error)
		go func() {
			if _, err := AwaitWithErr(dut2, bgp.BGPPath.Neighbor(dut3spec.RouterID).SessionState().State(), oc.Bgp_Neighbor_SessionState_IDLE); err != nil {
				awaitNewSession <- err
			}
			if _, err := AwaitWithErr(dut3, bgp.BGPPath.Neighbor(dut2spec.RouterID).SessionState().State(), oc.Bgp_Neighbor_SessionState_IDLE); err != nil {
				awaitNewSession <- err
			}
			if _, err := AwaitWithErr(dut2, bgp.BGPPath.Neighbor(dut3spec.RouterID).SessionState().State(), oc.Bgp_Neighbor_SessionState_ESTABLISHED); err != nil {
				awaitNewSession <- err
			}
			if _, err := AwaitWithErr(dut3, bgp.BGPPath.Neighbor(dut2spec.RouterID).SessionState().State(), oc.Bgp_Neighbor_SessionState_ESTABLISHED); err != nil {
				awaitNewSession <- err
			}
			awaitNewSession <- nil
		}()
		testspec.installPolicies(t, dut2)
		// Changing policy resets the BGP session, which causes routes
		// to disappear from the AdjRIBs, so we need to wait for
		// re-establishment first.
		if err := <-awaitNewSession; err != nil {
			t.Fatal(err)
		}

		for _, routeTest := range testspec.spec.RouteTests {
			testPropagation(t, routeTest, dut1, dut2, dut3, dut1spec.RouterID, dut2spec.RouterID, dut3spec.RouterID, true)
		}
		for _, routeTest := range testspec.spec.LongerPathRouteTests {
			testPropagation(t, routeTest, dut5, dut2, dut3, dut5spec.RouterID, dut2spec.RouterID, dut3spec.RouterID, true)
		}
	}
}
