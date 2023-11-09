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

	valpb "github.com/openconfig/lemming/proto/policyval"
)

func TestCommunitySet(t *testing.T) {
	installRejectPolicy := func(t *testing.T, dut1, dut2, dut3, _, _ *Device) {
		// Policy to reject routes with the given community set conditions
		policyName := "community-sets"
		policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}

		// Create prefix set for ANY/ALL policy
		prefixSetName := "tenRoutes"
		tenRoutes := "10.0.0.0/8"
		prefixSetPath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName)
		Replace(t, dut2, prefixSetPath.Mode().Config(), oc.PrefixSet_Mode_IPV4)
		Replace(t, dut2, prefixSetPath.Prefix(tenRoutes, "8..32").IpPrefix().Config(), tenRoutes)

		// Create ANY community set
		anyCommSetName := "ANY-community-set"
		anyCommMemberPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(anyCommSetName).CommunityMember()
		Replace(t, dut2, anyCommMemberPath.Config(), []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
			oc.UnionString("11111:11111"),
			oc.UnionString("22222:22222"),
		})
		Replace(t, dut2, ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(anyCommSetName).MatchSetOptions().Config(), oc.PolicyTypes_MatchSetOptionsType_ANY)

		// Create ALL community set
		allCommSetName := "ALL-community-set"
		allCommMemberPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(allCommSetName).CommunityMember()
		Replace(t, dut2, allCommMemberPath.Config(), []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
			oc.UnionString("33333:33333"),
			oc.UnionString("44444:44444"),
		})
		Replace(t, dut2, ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(allCommSetName).MatchSetOptions().Config(), oc.PolicyTypes_MatchSetOptionsType_ALL)

		// Create INVERT community set
		invertCommSetName := "INVERT-community-set"
		invertCommMemberPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(invertCommSetName).CommunityMember()
		Replace(t, dut2, invertCommMemberPath.Config(), []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
			oc.UnionString("11111:11111"),
			oc.UnionString("22222:22222"),
		})
		Replace(t, dut2, ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(invertCommSetName).MatchSetOptions().Config(), oc.PolicyTypes_MatchSetOptionsType_INVERT)

		// Match on given list of community set members and reject them.
		stmt, err := policy.AppendNew("reject-any-community-sets")
		if err != nil {
			t.Fatalf("Cannot append new BGP policy statement: %v", err)
		}
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().SetCommunitySet(anyCommSetName)
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)

		stmt, err = policy.AppendNew("reject-all-community-sets")
		if err != nil {
			t.Fatalf("Cannot append new BGP policy statement: %v", err)
		}
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().SetCommunitySet(allCommSetName)
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)

		// Create prefix set for INVERT policy
		prefixSetName = "twentyRoutes"
		twentyRoutes := "20.0.0.0/8"
		prefixSetPath = ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName)
		Replace(t, dut2, prefixSetPath.Mode().Config(), oc.PrefixSet_Mode_IPV4)
		Replace(t, dut2, prefixSetPath.Prefix(twentyRoutes, "8..32").IpPrefix().Config(), twentyRoutes)

		stmt, err = policy.AppendNew("reject-invert-community-sets")
		if err != nil {
			t.Fatalf("Cannot append new BGP policy statement: %v", err)
		}
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().SetCommunitySet(invertCommSetName)
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)

		// Policy to reject routes with the given community set conditions
		uselessPolicyName := "useless"
		uselessPolicy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
		stmt, err = uselessPolicy.AppendNew("reject-any-community-sets")
		if err != nil {
			t.Fatalf("Cannot append new BGP policy statement: %v", err)
		}
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().SetCommunitySet(anyCommSetName)
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)

		// Install policy
		Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Statement: policy})
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{policyName})
		// This apply-policy is a no-op because everything is rejected
		// in the reverse direction anyways: its purpose is to check
		// that statements across different policies can have the same
		// name (GoBGP requires all statement names to be unique).
		Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(uselessPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Statement: uselessPolicy})
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut3.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{uselessPolicyName})

		Await(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().ImportPolicy().State(), []string{policyName})
		Await(t, dut2, bgp.BGPPath.Neighbor(dut3.RouterID).ApplyPolicy().ImportPolicy().State(), []string{uselessPolicyName})
	}

	routeUnderTestList := []string{
		"10.0.0.0/10",
		"10.0.0.0/11",
		"10.0.0.0/12",
		"10.0.0.0/13",
		"10.0.0.0/14",
		"10.0.0.0/15",
		"10.0.0.0/16",
		"10.0.0.0/17",
		"20.0.0.0/18",
		"20.0.0.0/19",
		"20.0.0.0/20",
	}

	installSetPolicy := func(t *testing.T, dut1, dut2, _, _, _ *Device, testRef bool) {
		if debug {
			fmt.Println("Installing test policies")
		}

		policyName := "set-communities"
		policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}

		for i, route := range routeUnderTestList {
			// Create prefix set
			prefixSetName := "accept-" + route
			prefixSetPath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName)
			Replace(t, dut1, prefixSetPath.Mode().Config(), oc.PrefixSet_Mode_IPV4)
			Replace(t, dut1, prefixSetPath.Prefix(route, "exact").IpPrefix().Config(), route)

			stmt, err := policy.AppendNew(fmt.Sprintf("stmt%d", i))
			if err != nil {
				t.Fatalf("Cannot append new BGP policy statement: %v", err)
			}
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)

			installSetCommunityPolicy := func(comms ...oc.UnionString) {
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_REPLACE)
				if testRef {
					commSetName := fmt.Sprintf("ref-set-%d", i)
					commPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(commSetName)
					var commUnions []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union
					for _, c := range comms {
						commUnions = append(commUnions, c)
					}
					Replace(t, dut1, commPath.CommunitySetName().Config(), commSetName)
					Replace(t, dut1, commPath.CommunityMember().Config(), commUnions)
					stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateReference().SetCommunitySetRefs([]string{commSetName})
					stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_REFERENCE)
				} else {
					var commUnions []oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union
					for _, c := range comms {
						commUnions = append(commUnions, c)
					}
					stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(commUnions)
					stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_INLINE)
				}
			}

			switch i {
			case 0:
				installSetCommunityPolicy(oc.UnionString("10000:10000"))
			case 1:
				installSetCommunityPolicy(oc.UnionString("11111:11111"))
			case 2:
				installSetCommunityPolicy(oc.UnionString("33333:33333"))
			case 3:
				installSetCommunityPolicy(oc.UnionString("33333:33333"), oc.UnionString("44444:44444"))
			case 4:
				installSetCommunityPolicy(
					oc.UnionString("55555:55555"),
					oc.UnionString("44444:44444"),
					oc.UnionString("33333:33333"),
				)
			case 5:
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_REPLACE)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
					[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
						oc.UnionString("33333:33333"),
					},
				)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_INLINE)

				stmt, err = policy.AppendNew(fmt.Sprintf("stmt%d-2", i))
				if err != nil {
					t.Fatalf("Cannot append new BGP policy statement: %v", err)
				}
				stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
					[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
						oc.UnionString("44444:44444"),
					},
				)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_INLINE)
			case 6:
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
					[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
						oc.UnionString("11111:11111"),
						oc.UnionString("22222:22222"),
					},
				)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_INLINE)

				stmt, err = policy.AppendNew(fmt.Sprintf("stmt%d-2", i))
				if err != nil {
					t.Fatalf("Cannot append new BGP policy statement: %v", err)
				}
				stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)

				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_REMOVE)
				commSetName := fmt.Sprintf("ref-set-%d", i)
				commPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(commSetName)
				commUnions := []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
					oc.UnionString("11111:11111"),
				}
				Replace(t, dut1, commPath.CommunitySetName().Config(), commSetName)
				Replace(t, dut1, commPath.CommunityMember().Config(), commUnions)
				commSetName2 := fmt.Sprintf("ref-set-%d-2", i)
				commPath2 := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(commSetName2)
				commUnions2 := []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
					oc.UnionString("22222:22222"),
					oc.UnionString("33333:33333"),
				}
				Replace(t, dut1, commPath2.CommunitySetName().Config(), commSetName2)
				Replace(t, dut1, commPath2.CommunityMember().Config(), commUnions2)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateReference().SetCommunitySetRefs([]string{commSetName, commSetName2})
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_REFERENCE)
			case 7:
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
					[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
						oc.UnionString("11111:11111"),
						oc.UnionString("22222:22222"),
					},
				)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_INLINE)

				stmt, err = policy.AppendNew(fmt.Sprintf("stmt%d-2", i))
				if err != nil {
					t.Fatalf("Cannot append new BGP policy statement: %v", err)
				}
				stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)

				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_REMOVE)
				commSetName := fmt.Sprintf("ref-set-%d", i)
				commPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(commSetName)
				commUnions := []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
					oc.UnionString("[0-9]+:[0-9]+"),
				}
				Replace(t, dut1, commPath.CommunitySetName().Config(), commSetName)
				Replace(t, dut1, commPath.CommunityMember().Config(), commUnions)
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateReference().SetCommunitySetRefs([]string{commSetName})
				stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_REFERENCE)
			case 8:
				installSetCommunityPolicy(oc.UnionString("10000:10000"))
			case 9:
				installSetCommunityPolicy(oc.UnionString("11111:11111"))
			case 10:
				installSetCommunityPolicy(oc.UnionString("22222:22222"))
			}
			stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)
		}
		// Install policy
		Replace(t, dut1, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(policyName), Statement: policy})
		Replace(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().Config(), []string{policyName})
		Await(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().State(), []string{policyName})
	}

	test := func(testRef bool) {
		testName := "test-set-community-inline"
		if testRef {
			testName = "test-set-community-ref"
		}
		t.Run(testName, func(t *testing.T) {
			testPolicy(t, PolicyTestCase{
				spec: &valpb.PolicyTestCase{
					Description: "Test community set ANY and ALL",
					RouteTests: []*valpb.RouteTestCase{{
						Description: "No match",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[0],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
					}, {
						Description: "Matches ANY",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[1],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
					}, {
						Description: "Partially matches ALL",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[2],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
					}, {
						Description: "Matches ALL",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[3],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
					}, {
						Description: "Matches ALL reversed and with extra community",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[4],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
					}, {
						Description: "Matches ALL after ADD",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[5],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
					}, {
						Description: "REMOVE",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[6],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
					}, {
						Description: "REMOVE-ALL",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[7],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
					}, {
						Description: "matches-invert",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[8],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
					}, {
						Description: "no-match-invert",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[9],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
					}, {
						Description: "no-match-invert-2",
						Input: &valpb.TestRoute{
							ReachPrefix: routeUnderTestList[10],
						},
						ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
					}},
				},
				installPolicies: func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device) {
					installSetPolicy(t, dut1, dut2, dut3, dut4, dut5, testRef)
					installRejectPolicy(t, dut1, dut2, dut3, dut4, dut5)
				},
			})
		})
	}

	test(false)
	test(true)
}
