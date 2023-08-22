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
	"github.com/openconfig/ygnmi/ygnmi"

	valpb "github.com/openconfig/lemming/bgp/tests/proto/policyval"
)

const (
	rejectedCommunitySet = oc.UnionString("12345:12345")
)

// TestCommunitySet tests set-community and community set.
//
// DUT2's import policy from DUT1 sets the community value.
// DUT2's export policy to DUT3 filters the prefix with the community value.
func TestCommunitySet(t *testing.T) {
	testPolicy(t, PolicyTestCase{
		spec: &valpb.PolicyTestCase{
			Description: "Test that one NLRI gets accepted and the other rejected via a community-set.",
			RouteTests: []*valpb.RouteTestCase{{
				Description: "Rejected route",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.33.0.0/16",
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
			}, {
				Description: "Accepted route",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.3.0.0/16",
				},
				ExpectedResultBeforePolicy: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
				ExpectedResult:             valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}},
		},
		installImportSetPolicies: func(t *testing.T, dut2 *ygnmi.Client) {
			if debug {
				fmt.Println("Installing test policies")
			}
			prefix1 := "10.33.0.0/16"

			// Policy to set community on dut2
			policyName := "def1"

			// Create prefix set
			prefixSetName := "accept-" + prefix1
			prefix1Path := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName).Prefix(prefix1, "exact").IpPrefix()
			Replace(t, dut2, prefix1Path.Config(), prefix1)

			policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			stmt, err := policy.AppendNew("stmt1")
			if err != nil {
				t.Fatalf("Cannot append new BGP policy statement: %v", err)
			}
			// Match on prefix set
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)
			// Add community set
			stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
			stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
				[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
					rejectedCommunitySet,
				},
			)
			// Accept the route so that it may be advertised.
			stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)
			// Install policy
			Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Statement: policy})
			Replace(t, dut2, bgp.BGPPath.Neighbor(dut1spec.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{policyName})
		},
		installExportFilterPolicies: func(t *testing.T, dut2 *ygnmi.Client) {
			if debug {
				fmt.Println("Installing test policies")
			}

			// Policy to reject routes with the given community set
			policyName := "def2"

			// Create community set
			rejectCommSetName := "reject-community-set"
			rejCommMemberPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(rejectCommSetName).CommunityMember()
			Replace(t, dut2, rejCommMemberPath.Config(), []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
				rejectedCommunitySet,
			})
			Replace(t, dut2, ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(rejectCommSetName).MatchSetOptions().Config(), oc.RoutingPolicy_MatchSetOptionsType_ANY)

			policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			stmt, err := policy.AppendNew("stmt2")
			if err != nil {
				t.Fatalf("Cannot append new BGP policy statement: %v", err)
			}
			// Match on given list of community set members.
			stmt.GetOrCreateConditions().GetOrCreateBgpConditions().SetCommunitySet(rejectCommSetName)
			stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)
			// Install policy
			Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Statement: policy})
			Replace(t, dut2, bgp.BGPPath.Neighbor(dut3spec.RouterID).ApplyPolicy().ExportPolicy().Config(), []string{policyName})
		},
	})
}
