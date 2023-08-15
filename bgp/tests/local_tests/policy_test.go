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
	rejectTimeout = 10 * time.Second
)

// PolicyTestCase contains the specifications for a single policy test.
//
// Limitations:
// * Does not check path attributes.
// * Only checks export policies.
type PolicyTestCase struct {
	spec          *valpb.PolicyTestCase
	installPolicy func(*testing.T, *ygnmi.Client)
}

func testPolicy(t *testing.T, testspec PolicyTestCase, installPolicyAfterRoutes bool) {
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

	installDefaultPolicies := func() {
		// Clear the path for routes to be propagated.
		Replace(t, dut1, bgp.BGPPath.Neighbor(dut2spec.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut1spec.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut3spec.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut3, bgp.BGPPath.Neighbor(dut2spec.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
	}
	installDefaultPolicies()

	if !installPolicyAfterRoutes {
		testspec.installPolicy(t, dut2)
	}

	establishSessionPair(t, dut1, dut2, dut1spec, dut2spec)
	establishSessionPair(t, dut2, dut3, dut2spec, dut3spec)

	for _, routeTest := range testspec.spec.RouteTests {
		// Install all test routes into DUT1.
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

	staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	v := GetAll(t, dut1, staticp.StaticAny().Config())
	t.Logf("Installed static route on %v: %s", dut1, formatYgot(v))

	v4uni := bgp.BGPPath.Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()

	for _, routeTest := range testspec.spec.RouteTests {
		prefix := routeTest.GetInput().GetReachPrefix()
		// Check propagation to AdjRibOutPre for all prefixes.
		Await(t, dut1, v4uni.Neighbor(dut2spec.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, dut1, v4uni.Neighbor(dut2spec.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, dut2, v4uni.LocRib().Route(prefix, oc.UnionString(dut1spec.RouterID), 0).Prefix().State(), prefix)
		Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		if routeTest.GetExpectedResult() == policyval.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT || installPolicyAfterRoutes {
			Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
			Await(t, dut3, v4uni.Neighbor(dut2spec.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
		} else {
			w := Watch(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
				_, ok := val.Val()
				return ok
			})
			if _, ok := w.Await(t); ok {
				t.Errorf("prefix %q was not rejected.", prefix)
			}
		}
	}

	if installPolicyAfterRoutes {
		testspec.installPolicy(t, dut2)

		for _, routeTest := range testspec.spec.RouteTests {
			prefix := routeTest.GetInput().GetReachPrefix()

			Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
			Await(t, dut2, v4uni.Neighbor(dut1spec.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), prefix)
			Await(t, dut2, v4uni.LocRib().Route(prefix, oc.UnionString(dut1spec.RouterID), 0).Prefix().State(), prefix)
			Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
			if routeTest.GetExpectedResult() == policyval.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT {
				Await(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)
				Await(t, dut3, v4uni.Neighbor(dut2spec.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
			} else {
				// Check rejected prefix has been withdrawn from AdjRib-out of DUT2 and AdjRib-in of DUT3.
				w := Watch(t, dut2, v4uni.Neighbor(dut3spec.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
					_, ok := val.Val()
					return !ok
				})
				if _, ok := w.Await(t); !ok {
					t.Errorf("prefix %q was not rejected within timeout.", prefix)
				}
				w = Watch(t, dut3, v4uni.Neighbor(dut2spec.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
					_, ok := val.Val()
					return !ok
				})
				if _, ok := w.Await(t); !ok {
					t.Errorf("prefix %q was not withdrawn from DUT3 within timeout.", prefix)
				}
			}
		}
	}
}
