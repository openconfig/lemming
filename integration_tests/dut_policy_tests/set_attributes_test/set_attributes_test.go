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

package integration_test

import (
	"testing"

	"github.com/openconfig/lemming/internal/binding"
	"github.com/openconfig/lemming/policytest"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/gnmi/oc"
	"github.com/openconfig/ygot/ygot"

	valpb "github.com/openconfig/lemming/bgp/tests/proto/policyval"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.Get(".."))
}

const (
	acceptedCommunitySet = oc.UnionString("23456:23456")
	rejectedCommunitySet = oc.UnionString("12345:12345")
	lowerLocalPref       = 41
	higherLocalPref      = 42
	lowerMED             = oc.UnionUint32(10)
	higherMED            = oc.UnionUint32(11)

	rejectCommSetName    = "reject-community-set"
	rejectASPathSetName  = "reject-as-path-set"
	rejectASPathSetName2 = "reject-as-path-set2"
)

var (
	rejectedASPathSet = []string{"64502"}
)

func singletonPrefixSetName(route string) string {
	return "only-" + route
}

// TestSetAttributes tests setting BGP attributes.
func TestSetAttributes(t *testing.T) {
	routesUnderTest := []string{
		"10.0.0.0/16",
		"10.1.0.0/16",
		"10.2.0.0/16",
		"10.3.0.0/16",
		"10.4.0.0/16",
		"10.5.0.0/16",
		"10.6.0.0/16",
		"10.7.0.0/16",
	}

	installDefinedSets := func(t *testing.T, dut1, dut2, dut5 *policytest.Device) {
		for _, route := range routesUnderTest {
			// Create prefix set containing just the route under test.
			prefixPath := policytest.RoutingPolicyPath.DefinedSets().PrefixSet(singletonPrefixSetName(route)).Prefix(route, "exact").IpPrefix()
			gnmi.Replace(t, dut1, prefixPath.Config(), route)
			gnmi.Replace(t, dut2, prefixPath.Config(), route)
			gnmi.Replace(t, dut5, prefixPath.Config(), route)
		}

		// Create a community set
		rejCommMemberPath := policytest.RoutingPolicyPath.DefinedSets().BgpDefinedSets().CommunitySet(rejectCommSetName).CommunityMember()
		gnmi.Replace(t, dut2, rejCommMemberPath.Config(), []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
			rejectedCommunitySet,
		})
		gnmi.Replace(t, dut2, policytest.RoutingPolicyPath.DefinedSets().BgpDefinedSets().CommunitySet(rejectCommSetName).MatchSetOptions().Config(), oc.RoutingPolicy_MatchSetOptionsType_ANY)

		// Create AS path-sets
		rejASPathSetPath := policytest.RoutingPolicyPath.DefinedSets().BgpDefinedSets().AsPathSet(rejectASPathSetName)
		gnmi.Replace(t, dut2, rejASPathSetPath.AsPathSetName().Config(), rejectASPathSetName)
		gnmi.Replace(t, dut2, rejASPathSetPath.AsPathSetMember().Config(), rejectedASPathSet)

		rejASPathSet2Path := policytest.RoutingPolicyPath.DefinedSets().BgpDefinedSets().AsPathSet(rejectASPathSetName2)
		gnmi.Replace(t, dut2, rejASPathSet2Path.AsPathSetName().Config(), rejectASPathSetName2)
		gnmi.Replace(t, dut2, rejASPathSet2Path.AsPathSetMember().Config(), []string{"64499"})
	}

	policytest.TestPolicy(t, policytest.TestCase{
		Spec: &valpb.PolicyTestCase{
			Description: "Test that one NLRI gets accepted and the otheris rejected via various attribute values.",
			RouteTests: []*valpb.RouteTestCase{{
				Description: "Accepted route with no attributes",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[0],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}, {
				Description: "Accepted route with some attributes",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[1],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}, {
				Description: "Rejected route due to community set",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[2],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
			}, {
				Description: "Unpreferred route due to local-pref",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[3],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
			}, {
				Description: "Unpreferred route due to MED",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[4],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
			}, {
				Description: "Unpreferred route due to AS path prepend",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[6],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
			}, {
				Description: "Rejected route due to AS path match on prepended AS",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[7],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
			}},
			LongerPathRouteTests: []*valpb.RouteTestCase{{
				Description: "Accepted route due to higher local-pref",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[3],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}, {
				Description: "Rejected route due to AS path match",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[5],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
			}, {
				Description: "Accepted route due to shorter AS path after competing route's AS path prepend",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[6],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}},
			AlternatePathRouteTests: []*valpb.RouteTestCase{{
				Description: "Accepted route due to lower MED",
				Input: &valpb.TestRoute{
					ReachPrefix: routesUnderTest[4],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}},
		},
		InstallPolicies: func(t *testing.T, pair12, pair52, pair23 *policytest.DevicePair) {
			dut1, dut2, dut5 := pair12.First, pair12.Second, pair52.First
			port1, port21, port25, port52 := pair12.FirstPort, pair12.SecondPort, pair52.SecondPort, pair52.FirstPort
			installDefinedSets(t, dut1, dut2, dut5)

			// DUT1 -> DUT2 export policy
			dut1ExportPolicyName := "set-attributes-dut1"
			dut1ExportPolicy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			// DUT5 -> DUT2 export policy
			dut5ExportPolicyName := "set-attributes-dut5"
			dut5ExportPolicy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			// DUT1 -> DUT2 import policy
			dut1ImportPolicyName := "dut1-import-policy"
			dut1ImportPolicy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			// DUT5 -> DUT2 import policy
			dut5ImportPolicyName := "dut5-import-policy"
			dut5ImportPolicy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}

			for i, route := range routesUnderTest {
				// Restrict each policy to apply only to the route under test.
				prefixSetName := singletonPrefixSetName(route)

				// Initialize policy statements that apply only to the route under test.
				// If the boolean is set to true, then the statement will be installed.
				var installDut1ExportStmt bool
				dut1ExportStmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-setattr-policy-dut1")}
				dut1ExportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut1ExportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)

				var installDut5ExportStmt bool
				dut5ExportStmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-setattr-policy-dut5")}
				dut5ExportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut5ExportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)

				var installDut1ImportStmt bool
				dut1ImportStmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-reject-policy-dut1")}
				dut1ImportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut1ImportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)

				var installDut5ImportStmt bool
				dut5ImportStmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-reject-policy-dut5")}
				dut5ImportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut5ImportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)

				// Create the corresponding set and filter policies for each test route.
				switch i {
				case 1:
					// Set a attributes that are not filtered so that the route remains accepted.
					installDut1ExportStmt = true
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
						[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
							acceptedCommunitySet,
						},
					)

					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().SetSetLocalPref(higherLocalPref)
				case 2:
					// Set communities to a filtered value.
					installDut1ExportStmt = true
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
						[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
							rejectedCommunitySet,
						},
					)

					// Match on given list of community set members and reject it.
					installDut1ImportStmt = true
					dut1ImportStmt.GetOrCreateConditions().GetOrCreateBgpConditions().SetCommunitySet(rejectCommSetName)
				case 3:
					// Set local-pref such that the longer AS Path is preferred.
					installDut1ExportStmt = true
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().SetSetLocalPref(lowerLocalPref)

					installDut5ExportStmt = true
					dut5ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().SetSetLocalPref(higherLocalPref)
				case 4:
					// Set MED and make higher RouterID's route be preferred.
					installDut1ExportStmt = true
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().SetSetMed(higherMED)

					installDut5ExportStmt = true
					dut5ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().SetSetMed(lowerMED)
				case 5:
					// Reject route with matching AS path set members.
					installDut5ImportStmt = true
					dut5ImportStmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchAsPathSet().SetAsPathSet(rejectASPathSetName)
					dut5ImportStmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchAsPathSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsType_ANY)
				case 6:
					// Set AS Path Prepend to lengthen AS path list such that route is no longer preferred.
					installDut1ExportStmt = true
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetAsPathPrepend().SetAsn(64499)
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetAsPathPrepend().SetRepeatN(2)
				case 7:
					// Set AS Path Prepend and then match on it for rejection.
					installDut1ExportStmt = true
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetAsPathPrepend().SetAsn(64499)
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetAsPathPrepend().SetRepeatN(2)

					installDut1ImportStmt = true
					dut1ImportStmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchAsPathSet().SetAsPathSet(rejectASPathSetName2)
					dut1ImportStmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchAsPathSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsType_ANY)
				}
				if installDut1ExportStmt {
					dut1ExportStmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)
					if err := dut1ExportPolicy.Append(dut1ExportStmt); err != nil {
						t.Fatalf("Cannot append new BGP policy statement: %v", err)
					}
				}
				if installDut5ExportStmt {
					dut5ExportStmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)
					if err := dut5ExportPolicy.Append(dut5ExportStmt); err != nil {
						t.Fatalf("Cannot append new BGP policy statement: %v", err)
					}
				}
				if installDut1ImportStmt {
					dut1ImportStmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)
					if err := dut1ImportPolicy.Append(dut1ImportStmt); err != nil {
						t.Fatalf("Cannot append new BGP policy statement: %v", err)
					}
				}
				if installDut5ImportStmt {
					dut5ImportStmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)
					if err := dut5ImportPolicy.Append(dut5ImportStmt); err != nil {
						t.Fatalf("Cannot append new BGP policy statement: %v", err)
					}
				}
			}
			// Install export policies
			gnmi.Replace(t, dut1, policytest.RoutingPolicyPath.PolicyDefinition(dut1ExportPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut1ExportPolicyName), Statement: dut1ExportPolicy})
			gnmi.Replace(t, dut1, policytest.BGPPath.Neighbor(port21.IPv4).ApplyPolicy().ExportPolicy().Config(), []string{dut1ExportPolicyName})
			gnmi.Replace(t, dut5, policytest.RoutingPolicyPath.PolicyDefinition(dut5ExportPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut5ExportPolicyName), Statement: dut5ExportPolicy})
			gnmi.Replace(t, dut5, policytest.BGPPath.Neighbor(port25.IPv4).ApplyPolicy().ExportPolicy().Config(), []string{dut5ExportPolicyName})
			// Install import policies
			gnmi.Replace(t, dut2, policytest.RoutingPolicyPath.PolicyDefinition(dut1ImportPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut1ImportPolicyName), Statement: dut1ImportPolicy})
			gnmi.Replace(t, dut2, policytest.BGPPath.Neighbor(port1.IPv4).ApplyPolicy().ImportPolicy().Config(), []string{dut1ImportPolicyName})
			gnmi.Replace(t, dut2, policytest.RoutingPolicyPath.PolicyDefinition(dut5ImportPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut5ImportPolicyName), Statement: dut5ImportPolicy})
			gnmi.Replace(t, dut2, policytest.BGPPath.Neighbor(port52.IPv4).ApplyPolicy().ImportPolicy().Config(), []string{dut5ImportPolicyName})
		},
	})
}
