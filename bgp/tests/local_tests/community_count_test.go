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

func TestCommunityCount(t *testing.T) {
	installRejectPolicy := func(t *testing.T, dut1, dut2, _, _, _ *Device) {
		policyName := "community-count"
		policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}

		// Create prefix set for prefixes under test.
		prefixSetName := "tenRoutes"
		tenRoutes := "10.0.0.0/8"
		prefixSetPath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName)
		Replace(t, dut2, prefixSetPath.Mode().Config(), oc.PrefixSet_Mode_IPV4)
		Replace(t, dut2, prefixSetPath.Prefix(tenRoutes, "8..32").IpPrefix().Config(), tenRoutes)

		// Reject size 3 communities.
		stmt, err := policy.AppendNew("reject-communities-sized-3")
		if err != nil {
			t.Fatalf("Cannot append new BGP policy statement: %v", err)
		}
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateCommunityCount().SetOperator(oc.PolicyTypes_ATTRIBUTE_COMPARISON_ATTRIBUTE_EQ)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateCommunityCount().SetValue(3)
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)

		stmt, err = policy.AppendNew("reject-communities-sized-le-1")
		if err != nil {
			t.Fatalf("Cannot append new BGP policy statement: %v", err)
		}
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateCommunityCount().SetOperator(oc.PolicyTypes_ATTRIBUTE_COMPARISON_ATTRIBUTE_LE)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateCommunityCount().SetValue(1)
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)

		stmt, err = policy.AppendNew("reject-communities-sized-ge-5")
		if err != nil {
			t.Fatalf("Cannot append new BGP policy statement: %v", err)
		}
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateCommunityCount().SetOperator(oc.PolicyTypes_ATTRIBUTE_COMPARISON_ATTRIBUTE_GE)
		stmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateCommunityCount().SetValue(5)
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)

		// Install policy
		Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Statement: policy})
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{policyName})

		Await(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().ImportPolicy().State(), []string{policyName})
	}

	// 10.* are subject to rejection due to community-count policy.
	// 20.* must propagate in all cases due to prefix-match rule.
	routesUnderTest := map[int]string{
		// For debugging: just comment out the ones you don't want to run.
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
		11: "20.0.0.0/21",
		12: "20.0.0.0/22",
		13: "20.0.0.0/23",
	}

	installSetPolicy := func(t *testing.T, dut1, dut2, _, _, _ *Device) {
		if debug {
			fmt.Println("Installing test policies")
		}

		policyName := "set-communities"
		policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}

		for i, route := range routesUnderTest {
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

			comms := []oc.UnionString{
				oc.UnionString("11111:11111"),
				oc.UnionString("22222:22222"),
				oc.UnionString("33333:33333"),
				oc.UnionString("44444:44444"),
				oc.UnionString("55555:55555"),
				oc.UnionString("60000:60000"),
				oc.UnionString("7777:7777"),
			}

			switch i {
			case 0:
			case 1, 2, 3, 4, 5, 6, 7:
				configureSetCommunityPolicy(t, i, dut1, stmt, true, comms[:i]...)
			case 8:
			case 9:
				configureSetCommunityPolicy(t, i, dut1, stmt, true, comms[:1]...)
			case 10:
				configureSetCommunityPolicy(t, i, dut1, stmt, true, comms[:2]...)
			case 11:
				configureSetCommunityPolicy(t, i, dut1, stmt, true, comms[:3]...)
			case 12:
				configureSetCommunityPolicy(t, i, dut1, stmt, true, comms[:4]...)
			case 13:
				configureSetCommunityPolicy(t, i, dut1, stmt, true, comms[:5]...)
			default:
				t.Fatalf("BGP set policy not specified for test case %d", i)
			}
			stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)
		}
		// Install policy
		Replace(t, dut1, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(policyName), Statement: policy})
		Replace(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().Config(), []string{policyName})
		Await(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().State(), []string{policyName})
	}

	testPolicy(t, &PolicyTestCase{
		description:         "community-count",
		skipValidateAttrSet: true,
		routeTests: []*policytest.RouteTestCase{{
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[0],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "0",
				ExpectedResult:               policytest.RouteDiscarded,
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
				ReachPrefix: routesUnderTest[1],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "1",
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
				ReachPrefix: routesUnderTest[2],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "2",
				ExpectedResult:               policytest.RouteAccepted,
				PrevAdjRibOutPreCommunities:  nil,
				PrevAdjRibOutPostCommunities: []string{"11111:11111", "22222:22222"},
				AdjRibInPreCommunities:       []string{"11111:11111", "22222:22222"},
				AdjRibInPostCommunities:      []string{"11111:11111", "22222:22222"},
				LocalRibCommunities:          []string{"11111:11111", "22222:22222"},
				AdjRibOutPreCommunities:      []string{"11111:11111", "22222:22222"},
				AdjRibOutPostCommunities:     []string{"11111:11111", "22222:22222"},
				NextAdjRibInPreCommunities:   []string{"11111:11111", "22222:22222"},
				NextLocalRibCommunities:      []string{"11111:11111", "22222:22222"},
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[3],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "3",
				ExpectedResult:               policytest.RouteDiscarded,
				PrevAdjRibOutPreCommunities:  nil,
				PrevAdjRibOutPostCommunities: []string{"11111:11111", "22222:22222", "33333:33333"},
				AdjRibInPreCommunities:       []string{"11111:11111", "22222:22222", "33333:33333"},
				AdjRibInPostCommunities:      nil,
				LocalRibCommunities:          nil,
				AdjRibOutPreCommunities:      nil,
				AdjRibOutPostCommunities:     nil,
				NextAdjRibInPreCommunities:   nil,
				NextLocalRibCommunities:      nil,
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[4],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "4",
				ExpectedResult:               policytest.RouteAccepted,
				PrevAdjRibOutPreCommunities:  nil,
				PrevAdjRibOutPostCommunities: []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				AdjRibInPreCommunities:       []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				AdjRibInPostCommunities:      []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				LocalRibCommunities:          []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				AdjRibOutPreCommunities:      []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				AdjRibOutPostCommunities:     []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				NextAdjRibInPreCommunities:   []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				NextLocalRibCommunities:      []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[5],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "5",
				ExpectedResult:               policytest.RouteDiscarded,
				PrevAdjRibOutPreCommunities:  nil,
				PrevAdjRibOutPostCommunities: []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555"},
				AdjRibInPreCommunities:       []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555"},
				AdjRibInPostCommunities:      nil,
				LocalRibCommunities:          nil,
				AdjRibOutPreCommunities:      nil,
				AdjRibOutPostCommunities:     nil,
				NextAdjRibInPreCommunities:   nil,
				NextLocalRibCommunities:      nil,
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[6],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "6",
				ExpectedResult:               policytest.RouteDiscarded,
				PrevAdjRibOutPreCommunities:  nil,
				PrevAdjRibOutPostCommunities: []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555", "60000:60000"},
				AdjRibInPreCommunities:       []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555", "60000:60000"},
				AdjRibInPostCommunities:      nil,
				LocalRibCommunities:          nil,
				AdjRibOutPreCommunities:      nil,
				AdjRibOutPostCommunities:     nil,
				NextAdjRibInPreCommunities:   nil,
				NextLocalRibCommunities:      nil,
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[7],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "7",
				ExpectedResult:               policytest.RouteDiscarded,
				PrevAdjRibOutPreCommunities:  nil,
				PrevAdjRibOutPostCommunities: []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555", "60000:60000", "7777:7777"},
				AdjRibInPreCommunities:       []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555", "60000:60000", "7777:7777"},
				AdjRibInPostCommunities:      nil,
				LocalRibCommunities:          nil,
				AdjRibOutPreCommunities:      nil,
				AdjRibOutPostCommunities:     nil,
				NextAdjRibInPreCommunities:   nil,
				NextLocalRibCommunities:      nil,
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[8],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "0-different-prefix",
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
				ReachPrefix: routesUnderTest[9],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "1-different-prefix",
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
				ReachPrefix: routesUnderTest[10],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "2-different-prefix",
				ExpectedResult:               policytest.RouteAccepted,
				PrevAdjRibOutPreCommunities:  nil,
				PrevAdjRibOutPostCommunities: []string{"11111:11111", "22222:22222"},
				AdjRibInPreCommunities:       []string{"11111:11111", "22222:22222"},
				AdjRibInPostCommunities:      []string{"11111:11111", "22222:22222"},
				LocalRibCommunities:          []string{"11111:11111", "22222:22222"},
				AdjRibOutPreCommunities:      []string{"11111:11111", "22222:22222"},
				AdjRibOutPostCommunities:     []string{"11111:11111", "22222:22222"},
				NextAdjRibInPreCommunities:   []string{"11111:11111", "22222:22222"},
				NextLocalRibCommunities:      []string{"11111:11111", "22222:22222"},
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[11],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "3-different-prefix",
				ExpectedResult:               policytest.RouteAccepted,
				PrevAdjRibOutPreCommunities:  nil,
				PrevAdjRibOutPostCommunities: []string{"11111:11111", "22222:22222", "33333:33333"},
				AdjRibInPreCommunities:       []string{"11111:11111", "22222:22222", "33333:33333"},
				AdjRibInPostCommunities:      []string{"11111:11111", "22222:22222", "33333:33333"},
				LocalRibCommunities:          []string{"11111:11111", "22222:22222", "33333:33333"},
				AdjRibOutPreCommunities:      []string{"11111:11111", "22222:22222", "33333:33333"},
				AdjRibOutPostCommunities:     []string{"11111:11111", "22222:22222", "33333:33333"},
				NextAdjRibInPreCommunities:   []string{"11111:11111", "22222:22222", "33333:33333"},
				NextLocalRibCommunities:      []string{"11111:11111", "22222:22222", "33333:33333"},
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[12],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "4-different-prefix",
				ExpectedResult:               policytest.RouteAccepted,
				PrevAdjRibOutPreCommunities:  nil,
				PrevAdjRibOutPostCommunities: []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				AdjRibInPreCommunities:       []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				AdjRibInPostCommunities:      []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				LocalRibCommunities:          []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				AdjRibOutPreCommunities:      []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				AdjRibOutPostCommunities:     []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				NextAdjRibInPreCommunities:   []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
				NextLocalRibCommunities:      []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444"},
			},
		}, {
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[13],
			},
			RouteTest: &policytest.RoutePathTestCase{
				Description:                  "5-different-prefix",
				ExpectedResult:               policytest.RouteAccepted,
				PrevAdjRibOutPreCommunities:  nil,
				PrevAdjRibOutPostCommunities: []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555"},
				AdjRibInPreCommunities:       []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555"},
				AdjRibInPostCommunities:      []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555"},
				LocalRibCommunities:          []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555"},
				AdjRibOutPreCommunities:      []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555"},
				AdjRibOutPostCommunities:     []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555"},
				NextAdjRibInPreCommunities:   []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555"},
				NextLocalRibCommunities:      []string{"11111:11111", "22222:22222", "33333:33333", "44444:44444", "55555:55555"},
			},
		}},
		installPolicies: func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device) {
			installSetPolicy(t, dut1, dut2, dut3, dut4, dut5)
			installRejectPolicy(t, dut1, dut2, dut3, dut4, dut5)
		},
	})
}
