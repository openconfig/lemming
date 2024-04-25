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

	"github.com/openconfig/ygot/ygot"

	"github.com/openconfig/lemming/bgp"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/policytest"
)

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

var rejectedASPathSet = []string{"64502"}

func singletonPrefixSetName(route string) string {
	return "only-" + route
}

// TestSetAttributes tests setting BGP attributes.
func TestSetAttributes(t *testing.T) {
	routesUnderTest := []string{
		"10.1.0.0/16",
		"10.2.0.0/16",
		"10.10.0.0/16",
		"10.11.0.0/16",
		"10.12.0.0/16",
		"10.13.0.0/16",
		"10.14.0.0/16",
		"10.15.0.0/16",
	}

	installDefinedSets := func(t *testing.T, dut1, dut2, dut5 *Device) {
		for _, route := range routesUnderTest {
			// Create prefix set containing just the route under test.
			prefixModePath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(singletonPrefixSetName(route)).Mode()
			prefixPath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(singletonPrefixSetName(route)).Prefix(route, "exact").IpPrefix()
			Replace(t, dut1, prefixModePath.Config(), oc.PrefixSet_Mode_IPV4)
			Replace(t, dut1, prefixPath.Config(), route)
			Replace(t, dut2, prefixModePath.Config(), oc.PrefixSet_Mode_IPV4)
			Replace(t, dut2, prefixPath.Config(), route)
			Replace(t, dut5, prefixModePath.Config(), oc.PrefixSet_Mode_IPV4)
			Replace(t, dut5, prefixPath.Config(), route)
		}

		// Create a community set
		rejCommMemberPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(rejectCommSetName).CommunityMember()
		Replace(t, dut2, rejCommMemberPath.Config(), []oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
			rejectedCommunitySet,
		})
		Replace(t, dut2, ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySet(rejectCommSetName).MatchSetOptions().Config(), oc.PolicyTypes_MatchSetOptionsType_ANY)

		// Create AS path-sets
		rejASPathSetPath := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().AsPathSet(rejectASPathSetName)
		Replace(t, dut2, rejASPathSetPath.AsPathSetName().Config(), rejectASPathSetName)
		Replace(t, dut2, rejASPathSetPath.AsPathSetMember().Config(), rejectedASPathSet)

		rejASPathSet2Path := ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().AsPathSet(rejectASPathSetName2)
		Replace(t, dut2, rejASPathSet2Path.AsPathSetName().Config(), rejectASPathSetName2)
		Replace(t, dut2, rejASPathSet2Path.AsPathSetMember().Config(), []string{"64499"})
	}

	testPolicy(t, &PolicyTestCase{
		description: "Test that one NLRI gets accepted and the other is rejected via various attribute values.",
		routeTests: []*policytest.RouteTestCase{{
			Description: "Accepted route with no attributes",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[0],
			},
			ExpectedResult: policytest.RouteAccepted,
		}, {
			Description: "Accepted route with some attributes",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[1],
			},
			ExpectedResult:               policytest.RouteAccepted,
			PrevAdjRibOutPreCommunities:  nil,
			PrevAdjRibOutPostCommunities: []string{"23456:23456"},
			AdjRibInPreCommunities:       []string{"23456:23456"},
			AdjRibInPostCommunities:      []string{"23456:23456"},
			LocalRibCommunities:          []string{"23456:23456"},
			AdjRibOutPreCommunities:      []string{"23456:23456"},
			AdjRibOutPostCommunities:     []string{"23456:23456"},
			NextAdjRibInPreCommunities:   []string{"23456:23456"},
			NextLocalRibCommunities:      []string{"23456:23456"},
			PrevAdjRibOutPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
			},
			PrevAdjRibOutPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
			AdjRibInPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
			AdjRibInPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
			LocalRibAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
			AdjRibOutPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
			AdjRibOutPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
				// local-pref doesn't propagate to EBGP neighbour.
			},
			NextAdjRibInPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
			},
			NextAdjRibInPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
			},
			NextLocalRibAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
			},
		}, {
			Description: "Rejected route due to community set",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[2],
			},
			ExpectedResult:               policytest.RouteDiscarded,
			PrevAdjRibOutPreCommunities:  nil,
			PrevAdjRibOutPostCommunities: []string{"12345:12345"},
			AdjRibInPreCommunities:       []string{"12345:12345"},
			AdjRibInPostCommunities:      nil,
			LocalRibCommunities:          nil,
			AdjRibOutPreCommunities:      nil,
			AdjRibOutPostCommunities:     nil,
			NextAdjRibInPreCommunities:   nil,
			NextLocalRibCommunities:      nil,
		}, {
			Description: "Unpreferred route due to local-pref",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[3],
			},
			ExpectedResult: policytest.RouteNotPreferred,
			PrevAdjRibOutPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
			},
			PrevAdjRibOutPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(lowerLocalPref),
			},
			AdjRibInPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(lowerLocalPref),
			},
			AdjRibInPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(lowerLocalPref),
			},
			AdjRibOutPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
		}, {
			Description: "Unpreferred route due to MED",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[4],
			},
			ExpectedResult: policytest.RouteNotPreferred,
			PrevAdjRibOutPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
			},
			PrevAdjRibOutPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
				Med:    ygot.Uint32(uint32(higherMED)),
			},
			AdjRibInPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
				Med:    ygot.Uint32(uint32(higherMED)),
			},
			AdjRibInPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
				Med:    ygot.Uint32(uint32(higherMED)),
			},
		}, {
			Description: "Unpreferred route due to AS path prepend",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[6],
			},
			ExpectedResult: policytest.RouteNotPreferred,
		}, {
			Description: "Rejected route due to AS path match on prepended AS",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[7],
			},
			ExpectedResult: policytest.RouteDiscarded,
		}},
		longerPathRouteTests: []*policytest.RouteTestCase{{
			Description: "Rejected route due to longer AS Path",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[0],
			},
			ExpectedResult: policytest.RouteNotPreferred,
		}, {
			Description: "Accepted route due to higher local-pref",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[3],
			},
			ExpectedResult: policytest.RouteAccepted,
			PrevAdjRibOutPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
			},
			PrevAdjRibOutPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
			AdjRibInPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
			AdjRibInPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
			LocalRibAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
			AdjRibOutPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin:    oc.BgpTypes_BgpOriginAttrType_IGP,
				LocalPref: ygot.Uint32(higherLocalPref),
			},
			AdjRibOutPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
				// local-pref doesn't propagate to EBGP neighbour.
			},
			NextAdjRibInPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
			},
			NextAdjRibInPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
			},
			NextLocalRibAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
			},
		}, {
			Description: "Rejected route due to AS path match",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[5],
			},
			ExpectedResult: policytest.RouteDiscarded,
		}, {
			Description: "Accepted route due to shorter AS path after competing route's AS path prepend",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[6],
			},
			ExpectedResult: policytest.RouteAccepted,
		}},
		alternatePathRouteTests: []*policytest.RouteTestCase{{
			Description: "Accepted route due to lower MED",
			Input: policytest.TestRoute{
				ReachPrefix: routesUnderTest[4],
			},
			ExpectedResult: policytest.RouteAccepted,
			PrevAdjRibOutPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
				Med:    ygot.Uint32(uint32(lowerMED)),
			},
			PrevAdjRibOutPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
				Med:    ygot.Uint32(uint32(lowerMED)),
			},
			AdjRibInPreAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
				Med:    ygot.Uint32(uint32(lowerMED)),
			},
			AdjRibInPostAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
				Med:    ygot.Uint32(uint32(lowerMED)),
			},
			LocalRibAttrs: &oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet{
				Origin: oc.BgpTypes_BgpOriginAttrType_IGP,
				Med:    ygot.Uint32(uint32(lowerMED)),
			},
		}},
		installPolicies: func(t *testing.T, dut1, dut2, _, _, dut5 *Device) {
			if debug {
				fmt.Println("Installing test policies")
			}
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
				dut1ExportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)

				var installDut5ExportStmt bool
				dut5ExportStmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-setattr-policy-dut5")}
				dut5ExportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut5ExportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)

				var installDut1ImportStmt bool
				dut1ImportStmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-reject-policy-dut1")}
				dut1ImportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut1ImportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)

				var installDut5ImportStmt bool
				dut5ImportStmt := &oc.RoutingPolicy_PolicyDefinition_Statement{Name: ygot.String(route + "-reject-policy-dut5")}
				dut5ImportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
				dut5ImportStmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)

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
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_INLINE)

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
					dut1ExportStmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_INLINE)

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
					dut5ImportStmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchAsPathSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsType_ANY)
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
					dut1ImportStmt.GetOrCreateConditions().GetOrCreateBgpConditions().GetOrCreateMatchAsPathSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsType_ANY)
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
			Replace(t, dut1, ocpath.Root().RoutingPolicy().PolicyDefinition(dut1ExportPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut1ExportPolicyName), Statement: dut1ExportPolicy})
			Replace(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().Config(), []string{dut1ExportPolicyName})
			Await(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().State(), []string{dut1ExportPolicyName})
			Replace(t, dut5, ocpath.Root().RoutingPolicy().PolicyDefinition(dut5ExportPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut5ExportPolicyName), Statement: dut5ExportPolicy})
			Replace(t, dut5, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().Config(), []string{dut5ExportPolicyName})
			Await(t, dut5, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().State(), []string{dut5ExportPolicyName})
			// Install import policies
			Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(dut1ImportPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut1ImportPolicyName), Statement: dut1ImportPolicy})
			Replace(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{dut1ImportPolicyName})
			Await(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().ImportPolicy().State(), []string{dut1ImportPolicyName})
			Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(dut5ImportPolicyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(dut5ImportPolicyName), Statement: dut5ImportPolicy})
			Replace(t, dut2, bgp.BGPPath.Neighbor(dut5.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{dut5ImportPolicyName})
			Await(t, dut2, bgp.BGPPath.Neighbor(dut5.RouterID).ApplyPolicy().ImportPolicy().State(), []string{dut5ImportPolicyName})
		},
	})
}
