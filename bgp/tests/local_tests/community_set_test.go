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
	"github.com/openconfig/lemming/policytest"
	"github.com/openconfig/ygot/ygot"
)

func TestCommunitySet(t *testing.T) {
	installRejectPolicy := func(t *testing.T, dut1, dut2, dut3, _, _ *Device, testRef bool) {
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
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchCommunitySet().SetCommunitySet(anyCommSetName)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchCommunitySet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsType_ANY)
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)

		stmt, err = policy.AppendNew("reject-all-community-sets")
		if err != nil {
			t.Fatalf("Cannot append new BGP policy statement: %v", err)
		}
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchCommunitySet().SetCommunitySet(allCommSetName)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchCommunitySet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsType_ALL)
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
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchCommunitySet().SetCommunitySet(invertCommSetName)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchCommunitySet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsType_INVERT)
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
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchCommunitySet().SetCommunitySet(anyCommSetName)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchCommunitySet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsType_ANY)
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)

		// Create prefix set for setting a new community
		prefixSetName = "thirtyRoutes"
		thirtyRoutes := "30.0.0.0/8"
		prefixSetPath = ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName)
		Replace(t, dut2, prefixSetPath.Mode().Config(), oc.PrefixSet_Mode_IPV4)
		Replace(t, dut2, prefixSetPath.Prefix(thirtyRoutes, "8..32").IpPrefix().Config(), thirtyRoutes)

		stmt, err = policy.AppendNew("set-new-community-sets")
		if err != nil {
			t.Fatalf("Cannot append new BGP policy statement: %v", err)
		}
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
		configureSetCommunityPolicy(t, 11, dut2, stmt, testRef, "22222:22222")
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)

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

	routeUnderTestList := map[int]string{
		0:  "10.0.0.0/10",
		1:  "10.0.0.0/11",
		2:  "10.0.0.0/12",
		3:  "10.0.0.0/13",
		4:  "10.0.0.0/14",
		5:  "10.0.0.0/15",
		6:  "10.0.0.0/16",
		7:  "10.0.0.0/17",
		8:  "20.0.0.0/18",
		9:  "20.0.0.0/19",
		10: "20.0.0.0/20",
		11: "30.0.0.0/21",
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

			switch i {
			case 0:
				configureSetCommunityPolicy(t, i, dut1, stmt, testRef, oc.UnionString("10000:10000"))
			case 1:
				configureSetCommunityPolicy(t, i, dut1, stmt, testRef, oc.UnionString("11111:11111"))
			case 2:
				configureSetCommunityPolicy(t, i, dut1, stmt, testRef, oc.UnionString("33333:33333"))
			case 3:
				configureSetCommunityPolicy(t, i, dut1, stmt, testRef, oc.UnionString("33333:33333"), oc.UnionString("44444:44444"))
			case 4:
				configureSetCommunityPolicy(t, i, dut1, stmt, testRef,
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
				configureSetCommunityPolicy(t, i, dut1, stmt, testRef, oc.UnionString("10000:10000"))
			case 9:
				configureSetCommunityPolicy(t, i, dut1, stmt, testRef, oc.UnionString("11111:11111"))
			case 10:
				configureSetCommunityPolicy(t, i, dut1, stmt, testRef, oc.UnionString("22222:22222"))
			case 11:
				configureSetCommunityPolicy(t, i, dut1, stmt, testRef, oc.UnionString("11111:11111"))
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
			testPolicy(t, &PolicyTestCase{
				description:         "Test community set ANY and ALL",
				skipValidateAttrSet: true,
				routeTests: []*policytest.RouteTestCase{{
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[0],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "No match",
						ExpectedResult:               policytest.RouteAccepted,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: []string{"10000:10000"},
						AdjRibInPreCommunities:       []string{"10000:10000"},
						AdjRibInPostCommunities:      []string{"10000:10000"},
						LocalRibCommunities:          []string{"10000:10000"},
						AdjRibOutPreCommunities:      []string{"10000:10000"},
						AdjRibOutPostCommunities:     []string{"10000:10000"},
						NextAdjRibInPreCommunities:   []string{"10000:10000"},
						NextLocalRibCommunities:      []string{"10000:10000"},
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[1],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "Matches ANY",
						ExpectedResult:               policytest.RouteDiscarded,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: []string{"11111:11111"},
						AdjRibInPreCommunities:       []string{"11111:11111"},
						AdjRibInPostCommunities:      nil,
						LocalRibCommunities:          nil,
						AdjRibOutPreCommunities:      nil,
						AdjRibOutPostCommunities:     nil,
						NextAdjRibInPreCommunities:   nil,
						NextLocalRibCommunities:      nil,
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[2],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "Partially matches ALL",
						ExpectedResult:               policytest.RouteAccepted,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: []string{"33333:33333"},
						AdjRibInPreCommunities:       []string{"33333:33333"},
						AdjRibInPostCommunities:      []string{"33333:33333"},
						LocalRibCommunities:          []string{"33333:33333"},
						AdjRibOutPreCommunities:      []string{"33333:33333"},
						AdjRibOutPostCommunities:     []string{"33333:33333"},
						NextAdjRibInPreCommunities:   []string{"33333:33333"},
						NextLocalRibCommunities:      []string{"33333:33333"},
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[3],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "Matches ALL",
						ExpectedResult:               policytest.RouteDiscarded,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: []string{"33333:33333", "44444:44444"},
						AdjRibInPreCommunities:       []string{"33333:33333", "44444:44444"},
						AdjRibInPostCommunities:      nil,
						LocalRibCommunities:          nil,
						AdjRibOutPreCommunities:      nil,
						AdjRibOutPostCommunities:     nil,
						NextAdjRibInPreCommunities:   nil,
						NextLocalRibCommunities:      nil,
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[4],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "Matches ALL reversed and with extra community",
						ExpectedResult:               policytest.RouteDiscarded,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: []string{"55555:55555", "44444:44444", "33333:33333"},
						AdjRibInPreCommunities:       []string{"55555:55555", "44444:44444", "33333:33333"},
						AdjRibInPostCommunities:      nil,
						LocalRibCommunities:          nil,
						AdjRibOutPreCommunities:      nil,
						AdjRibOutPostCommunities:     nil,
						NextAdjRibInPreCommunities:   nil,
						NextLocalRibCommunities:      nil,
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[5],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "Matches ALL after ADD",
						ExpectedResult:               policytest.RouteDiscarded,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: []string{"33333:33333", "44444:44444"},
						AdjRibInPreCommunities:       []string{"33333:33333", "44444:44444"},
						AdjRibInPostCommunities:      nil,
						LocalRibCommunities:          nil,
						AdjRibOutPreCommunities:      nil,
						AdjRibOutPostCommunities:     nil,
						NextAdjRibInPreCommunities:   nil,
						NextLocalRibCommunities:      nil,
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[6],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "REMOVE",
						ExpectedResult:               policytest.RouteAccepted,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: nil,
						AdjRibInPreCommunities:       nil,
						AdjRibInPostCommunities:      nil,
						LocalRibCommunities:          nil,
						AdjRibOutPreCommunities:      nil,
						AdjRibOutPostCommunities:     nil,
						NextAdjRibInPreCommunities:   nil,
						NextLocalRibCommunities:      nil,
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[7],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "REMOVE-ALL",
						ExpectedResult:               policytest.RouteAccepted,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: nil,
						AdjRibInPreCommunities:       nil,
						AdjRibInPostCommunities:      nil,
						LocalRibCommunities:          nil,
						AdjRibOutPreCommunities:      nil,
						AdjRibOutPostCommunities:     nil,
						NextAdjRibInPreCommunities:   nil,
						NextLocalRibCommunities:      nil,
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[8],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "matches-invert",
						ExpectedResult:               policytest.RouteDiscarded,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: []string{"10000:10000"},
						AdjRibInPreCommunities:       []string{"10000:10000"},
						AdjRibInPostCommunities:      nil,
						LocalRibCommunities:          nil,
						AdjRibOutPreCommunities:      nil,
						AdjRibOutPostCommunities:     nil,
						NextAdjRibInPreCommunities:   nil,
						NextLocalRibCommunities:      nil,
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[9],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "no-match-invert",
						ExpectedResult:               policytest.RouteAccepted,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: []string{"11111:11111"},
						AdjRibInPreCommunities:       []string{"11111:11111"},
						AdjRibInPostCommunities:      []string{"11111:11111"},
						LocalRibCommunities:          []string{"11111:11111"},
						AdjRibOutPreCommunities:      []string{"11111:11111"},
						AdjRibOutPostCommunities:     []string{"11111:11111"},
						NextAdjRibInPreCommunities:   []string{"11111:11111"},
						NextLocalRibCommunities:      []string{"11111:11111"},
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[10],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "no-match-invert-2",
						ExpectedResult:               policytest.RouteAccepted,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: []string{"22222:22222"},
						AdjRibInPreCommunities:       []string{"22222:22222"},
						AdjRibInPostCommunities:      []string{"22222:22222"},
						LocalRibCommunities:          []string{"22222:22222"},
						AdjRibOutPreCommunities:      []string{"22222:22222"},
						AdjRibOutPostCommunities:     []string{"22222:22222"},
						NextAdjRibInPreCommunities:   []string{"22222:22222"},
						NextLocalRibCommunities:      []string{"22222:22222"},
					},
				}, {
					Input: policytest.TestRoute{
						ReachPrefix: routeUnderTestList[11],
					},
					RouteTest: &policytest.RoutePathTestCase{
						Description:                  "community-modified-at-import",
						ExpectedResult:               policytest.RouteAccepted,
						PrevAdjRibOutPreCommunities:  nil,
						PrevAdjRibOutPostCommunities: []string{"11111:11111"},
						AdjRibInPreCommunities:       []string{"11111:11111"},
						AdjRibInPostCommunities:      []string{"22222:22222"},
						LocalRibCommunities:          []string{"22222:22222"},
						AdjRibOutPreCommunities:      []string{"22222:22222"},
						AdjRibOutPostCommunities:     []string{"22222:22222"},
						NextAdjRibInPreCommunities:   []string{"22222:22222"},
						NextLocalRibCommunities:      []string{"22222:22222"},
					},
				}},
				installPolicies: func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device) {
					installSetPolicy(t, dut1, dut2, dut3, dut4, dut5, testRef)
					installRejectPolicy(t, dut1, dut2, dut3, dut4, dut5, testRef)
				},
			})
		})
	}

	test(false)
	test(true)
}

// configureSetCommunityPolicy adds the community set to the given device, and
// configures the input statement with a community replace action.
func configureSetCommunityPolicy(t *testing.T, i int, dut *Device, stmt *oc.RoutingPolicy_PolicyDefinition_Statement, testRef bool, comms ...oc.UnionString) {
	stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_REPLACE)
	if testRef {
		commSetName := fmt.Sprintf("ref-set-%d", i)
		commPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(commSetName)
		var commUnions []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union
		for _, c := range comms {
			commUnions = append(commUnions, c)
		}
		Replace(t, dut, commPath.CommunitySetName().Config(), commSetName)
		Replace(t, dut, commPath.CommunityMember().Config(), commUnions)
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
