// Copyright 2024 Google LLC
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
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/policytest"
	"github.com/openconfig/ygot/ygot"
)

func TestRouteType(t *testing.T) {
	routesUnderTest := []string{
		// For debugging: just comment out the ones you don't want to run.
		0: "10.0.0.0/10",
		1: "10.0.0.0/11",
		2: "10.0.0.0/12",
		3: "10.0.0.0/13",
		4: "10.0.0.0/14",
		5: "10.0.0.0/15",
	}

	installRejectPolicy := func(t *testing.T, dut1, dut2, _, _, dut5 *Device) {
		if debug {
			fmt.Println("Installing test policies")
		}

		policyName := "route-type"
		policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}

		for i, route := range routesUnderTest {
			// Create prefix set
			prefixSetName := "for-" + route
			prefixSetPath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName)
			Replace(t, dut2, prefixSetPath.Mode().Config(), oc.PrefixSet_Mode_IPV4)
			Replace(t, dut2, prefixSetPath.Prefix(route, "exact").IpPrefix().Config(), route)

			stmt, err := policy.AppendNew(fmt.Sprintf("stmt%d", i))
			if err != nil {
				t.Fatalf("Cannot append new BGP policy statement: %v", err)
			}
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)

			switch i {
			case 0, 3:
			case 1, 4:
				stmt.GetOrCreateConditions().GetOrCreateBgpConditions().SetRouteType(oc.BgpConditions_RouteType_INTERNAL)
			case 2, 5:
				stmt.GetOrCreateConditions().GetOrCreateBgpConditions().SetRouteType(oc.BgpConditions_RouteType_EXTERNAL)
			default:
				t.Fatalf("BGP set policy not specified for test case %d", i)
			}
			stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)
		}
		// Install policy
		Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(policyName), Statement: policy})
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{policyName})
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut5.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{policyName})
		Await(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().ImportPolicy().State(), []string{policyName})
		Await(t, dut2, bgp.BGPPath.Neighbor(dut5.RouterID).ApplyPolicy().ImportPolicy().State(), []string{policyName})
	}

	testPolicy(t, &PolicyTestCase{
		description:         "route-type",
		skipValidateAttrSet: true,
		dut1IsEBGP:          true,
		routeTests: []*policytest.RouteTestCase{{
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[0],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:    "not set",
				ExpectedResult: policytest.RouteDiscarded,
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[1],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:    "internal",
				ExpectedResult: policytest.RouteAccepted,
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[2],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:    "external",
				ExpectedResult: policytest.RouteDiscarded,
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[3],
			},
			AlternatePathRouteTest: &policytest.RoutePathTestCase{
				Description:    "not set",
				ExpectedResult: policytest.RouteDiscarded,
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[4],
			},
			AlternatePathRouteTest: &policytest.RoutePathTestCase{
				Description:    "internal",
				ExpectedResult: policytest.RouteDiscarded,
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[5],
			},
			AlternatePathRouteTest: &policytest.RoutePathTestCase{
				Description:    "external",
				ExpectedResult: policytest.RouteAccepted,
			},
		}},
		installPolicies: func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device) {
			installRejectPolicy(t, dut1, dut2, dut3, dut4, dut5)
		},
	})
}
