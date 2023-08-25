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
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygot/ygot"

	valpb "github.com/openconfig/lemming/bgp/tests/proto/policyval"
)

const (
	acceptedCommunitySet = oc.UnionString("23456:23456")
	rejectedCommunitySet = oc.UnionString("12345:12345")
	lowerLocalPref       = 41
	higherLocalPref      = 42
	lowerMED             = oc.UnionUint32(10)
	higherMED            = oc.UnionUint32(11)
)

var (
	rejectedASPathSet = []string{"64502"}
)

// TestAttributes tests BGP attributes.
// - set-community and community set.
// - set-local-pref and local pref.
//
// DUT2's import policy from DUT1 sets the attribute values.
// DUT2's export policy to DUT3 filters the prefix with the attribute value.
func TestAttributes(t *testing.T) {
	routeList := []string{
		"10.1.0.0/16",
		"10.2.0.0/16",
		"10.10.0.0/16",
		"10.11.0.0/16",
		"10.12.0.0/16",
		"10.13.0.0/16",
		"10.14.0.0/16",
		"10.15.0.0/16",
	}
	testPolicy(t, PolicyTestCase{
		spec: &valpb.PolicyTestCase{
			Description: "Test that one NLRI gets accepted and the otheris rejected via various attribute values.",
			RouteTests: []*valpb.RouteTestCase{{
				Description: "Accepted route with no attributes",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[0],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}, {
				Description: "Accepted route with some attributes",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[1],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}, {
				Description: "Rejected route due to community set",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[2],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
			}, {
				Description: "Unpreferred route due to local-pref",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[3],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
			}, {
				Description: "Unpreferred route due to MED",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[4],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
			}, {
				Description: "Unpreferred route due to AS path prepend",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[6],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
			}, {
				Description: "Rejected route due to AS path match on prepended AS",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[7],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
			}},
			LongerPathRouteTests: []*valpb.RouteTestCase{{
				Description: "Accepted route due to higher local-pref",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[3],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}, {
				Description: "Rejected route due to AS path match",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[5],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
			}, {
				Description: "Accepted route due to shorter AS path after competing route's AS path prepend",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[6],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}},
			AlternatePathRouteTests: []*valpb.RouteTestCase{{
				Description: "Accepted route due to lower MED",
				Input: &valpb.TestRoute{
					ReachPrefix: routeList[4],
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}},
		},
		installPolicies: func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device) {
			if debug {
				fmt.Println("Installing test policies")
			}

			dut1PolicyName := "set-attributes-dut1"
			dut1Policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			dut5PolicyName := "set-attributes-dut5"
			dut5Policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			dut1RejectPolicyName := "dut1-import-policy"
			dut1RejectPolicy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			dut5RejectPolicyName := "dut5-import-policy"
			dut5RejectPolicy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}

			for i, route := range routeList {
				// Create prefix set
				prefixSetName := "accept-" + route
				prefixPath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName).Prefix(route, "exact").IpPrefix()
				Replace(t, dut1, prefixPath.Config(), route)
				Replace(t, dut2, prefixPath.Config(), route)
				Replace(t, dut5, prefixPath.Config(), route)

				var installDut1Stmt bool
				dut1Stmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-setattr-policy-dut1")}
				dut1Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut1Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)

				var installDut5Stmt bool
				dut5Stmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-setattr-policy-dut5")}
				dut5Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut5Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)

				var installDut1RejectStmt bool
				dut1RejectStmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-reject-policy-dut1")}
				dut1RejectStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut1RejectStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)

				var installDut5RejectStmt bool
				dut5RejectStmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-reject-policy-dut5")}
				dut5RejectStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut5RejectStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)

				// Create the corresponding set and filter policies for each test route.
				switch i {
				case 1:
					// Set communities
					installDut1Stmt = true
					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
						[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
							acceptedCommunitySet,
						},
					)

					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().SetSetLocalPref(higherLocalPref)
				case 2:
					// Set communities
					installDut1Stmt = true
					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
						[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
							rejectedCommunitySet,
						},
					)

					// Create community set
					installDut1RejectStmt = true
					rejectCommSetName := "reject-community-set"
					rejCommMemberPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(rejectCommSetName).CommunityMember()
					Replace(t, dut2, rejCommMemberPath.Config(), []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
						rejectedCommunitySet,
					})
					Replace(t, dut2, ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(rejectCommSetName).MatchSetOptions().Config(), oc.RoutingPolicy_MatchSetOptionsType_ANY)

					// Match on given list of community set members.
					dut1RejectStmt.GetOrCreateConditions().GetOrCreateBgpConditions().SetCommunitySet(rejectCommSetName)
				case 3:
					// Set local-pref
					installDut1Stmt = true
					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().SetSetLocalPref(lowerLocalPref)

					installDut5Stmt = true
					dut5Stmt.GetOrCreateActions().GetOrCreateBgpActions().SetSetLocalPref(higherLocalPref)
				case 4:
					// Set MED
					installDut1Stmt = true
					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().SetSetMed(higherMED)

					installDut5Stmt = true
					dut5Stmt.GetOrCreateActions().GetOrCreateBgpActions().SetSetMed(lowerMED)
				case 5:
					// AS-path-set
					installDut5RejectStmt = true
					rejectASPathSetName := "reject-as-path-set"
					rejASPathSetPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().AsPathSet(rejectASPathSetName)
					Replace(t, dut2, rejASPathSetPath.AsPathSetName().Config(), rejectASPathSetName)
					Replace(t, dut2, rejASPathSetPath.AsPathSetMember().Config(), rejectedASPathSet)

					// Match on given list of AS path set members.
					dut5RejectStmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchAsPathSet().SetAsPathSet(rejectASPathSetName)
					dut5RejectStmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchAsPathSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsType_ANY)
				case 6:
					// Set AS Path Prepend
					installDut1Stmt = true
					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetAsPathPrepend().SetAsn(64499)
					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetAsPathPrepend().SetRepeatN(2)
				case 7:
					// Set AS Path Prepend
					installDut1Stmt = true
					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetAsPathPrepend().SetAsn(64499)
					dut1Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetAsPathPrepend().SetRepeatN(2)

					// AS-path-set
					installDut1RejectStmt = true
					rejectASPathSetName := "reject-as-path-set2"
					rejASPathSetPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().AsPathSet(rejectASPathSetName)
					Replace(t, dut2, rejASPathSetPath.AsPathSetName().Config(), rejectASPathSetName)
					Replace(t, dut2, rejASPathSetPath.AsPathSetMember().Config(), []string{"64499"})

					// Match on given list of AS path set members.
					dut1RejectStmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchAsPathSet().SetAsPathSet(rejectASPathSetName)
					dut1RejectStmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchAsPathSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsType_ANY)
				}
				if installDut1Stmt {
					dut1Stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)
					if err := dut1Policy.Append(dut1Stmt); err != nil {
						t.Fatalf("Cannot append new BGP policy statement: %v", err)
					}
				}
				if installDut5Stmt {
					dut5Stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)
					if err := dut5Policy.Append(dut5Stmt); err != nil {
						t.Fatalf("Cannot append new BGP policy statement: %v", err)
					}
				}
				if installDut1RejectStmt {
					dut1RejectStmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)
					if err := dut1RejectPolicy.Append(dut1RejectStmt); err != nil {
						t.Fatalf("Cannot append new BGP policy statement: %v", err)
					}
				}
				if installDut5RejectStmt {
					dut5RejectStmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)
					if err := dut5RejectPolicy.Append(dut5RejectStmt); err != nil {
						t.Fatalf("Cannot append new BGP policy statement: %v", err)
					}
				}
			}
			// Install export policies
			Replace(t, dut1, ocpath.Root().RoutingPolicy().PolicyDefinition(dut1PolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut1PolicyName), Statement: dut1Policy})
			Replace(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().Config(), []string{dut1PolicyName})
			Replace(t, dut5, ocpath.Root().RoutingPolicy().PolicyDefinition(dut5PolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut5PolicyName), Statement: dut5Policy})
			Replace(t, dut5, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().Config(), []string{dut5PolicyName})
			// Install import policies
			Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(dut1RejectPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut1RejectPolicyName), Statement: dut1RejectPolicy})
			Replace(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{dut1RejectPolicyName})
			Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(dut5RejectPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut5RejectPolicyName), Statement: dut5RejectPolicy})
			Replace(t, dut2, bgp.BGPPath.Neighbor(dut5.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{dut5RejectPolicyName})
		},
	})
}
