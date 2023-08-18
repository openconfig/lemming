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

package bgp

import (
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/wenovus/gobgp/v3/pkg/bgpconfig"
)

// convertPolicyName converts from OC policy name to a neighbour-qualified
// policy name in order to put all the policies into a global list.
func convertPolicyName(neighAddr, ocPolicyName string) string {
	return neighAddr + "|" + ocPolicyName
}

func convertPolicyNames(neighAddr string, ocPolicyNames []string) []string {
	var convertedNames []string
	for _, n := range ocPolicyNames {
		convertedNames = append(convertedNames, convertPolicyName(neighAddr, n))
	}
	return convertedNames
}

// convertPolicyDefinition converts an OC policy definition to GoBGP policy definition.
//
// It adds neighbour set to disambiguate it from another instance of the policy
// for another neighbour. This is necessary since all policies will go into a
// single apply-policy list.
func convertPolicyDefinition(policy *oc.RoutingPolicy_PolicyDefinition, neighAddr string, occommset map[string]*oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet) bgpconfig.PolicyDefinition {
	var statements []bgpconfig.Statement
	for _, statement := range policy.Statement.Values() {
		var setCommunitiesList []string
		for _, comm := range statement.GetActions().BgpActions.GetSetCommunity().GetInline().GetCommunities() {
			setCommunitiesList = append(setCommunitiesList, convertCommunity(comm))
		}
		statements = append(statements, bgpconfig.Statement{
			Name: statement.GetName(),
			Conditions: bgpconfig.Conditions{
				MatchPrefixSet: bgpconfig.MatchPrefixSet{
					PrefixSet:       statement.GetConditions().GetMatchPrefixSet().GetPrefixSet(),
					MatchSetOptions: convertMatchSetOptionsRestrictedType(statement.GetConditions().GetMatchPrefixSet().GetMatchSetOptions()),
				},
				MatchNeighborSet: bgpconfig.MatchNeighborSet{
					// Name the neighbor set as the policy so that the policy only applies to referring neighbours.
					NeighborSet: neighAddr,
				},
				BgpConditions: bgpconfig.BgpConditions{
					MatchCommunitySet: bgpconfig.MatchCommunitySet{
						CommunitySet:    statement.Conditions.GetBgpConditions().GetCommunitySet(),
						MatchSetOptions: convertMatchSetOptionsType(occommset[statement.Conditions.GetBgpConditions().GetCommunitySet()].GetMatchSetOptions()),
					},
				},
			},
			Actions: bgpconfig.Actions{
				RouteDisposition: convertRouteDisposition(statement.GetActions().GetPolicyResult()),
				BgpActions: bgpconfig.BgpActions{
					SetCommunity: bgpconfig.SetCommunity{
						SetCommunityMethod: bgpconfig.SetCommunityMethod{
							CommunitiesList: setCommunitiesList,
						},
						Options: statement.GetActions().GetBgpActions().GetSetCommunity().GetOptions().String(),
					},
				},
			},
		})
	}

	return bgpconfig.PolicyDefinition{
		Name:       convertPolicyName(neighAddr, policy.GetName()),
		Statements: statements,
	}
}

func convertNeighborApplyPolicy(neigh *oc.NetworkInstance_Protocol_Bgp_Neighbor) bgpconfig.ApplyPolicy {
	return bgpconfig.ApplyPolicy{
		Config: bgpconfig.ApplyPolicyConfig{
			DefaultImportPolicy: convertDefaultPolicy(neigh.GetApplyPolicy().GetDefaultImportPolicy()),
			DefaultExportPolicy: convertDefaultPolicy(neigh.GetApplyPolicy().GetDefaultExportPolicy()),
			ImportPolicyList:    neigh.GetApplyPolicy().GetImportPolicy(),
			ExportPolicyList:    neigh.GetApplyPolicy().GetExportPolicy(),
		},
	}
}

func convertDefaultPolicy(ocpolicy oc.E_RoutingPolicy_DefaultPolicyType) bgpconfig.DefaultPolicyType {
	switch ocpolicy {
	case oc.RoutingPolicy_DefaultPolicyType_REJECT_ROUTE:
		return bgpconfig.DEFAULT_POLICY_TYPE_REJECT_ROUTE
	case oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE:
		return bgpconfig.DEFAULT_POLICY_TYPE_ACCEPT_ROUTE
	default:
		return bgpconfig.DEFAULT_POLICY_TYPE_REJECT_ROUTE
	}
}

func convertMatchSetOptionsType(ocMatchSetOpts oc.E_RoutingPolicy_MatchSetOptionsType) bgpconfig.MatchSetOptionsType {
	switch ocMatchSetOpts {
	case oc.RoutingPolicy_MatchSetOptionsType_INVERT:
		return bgpconfig.MATCH_SET_OPTIONS_TYPE_INVERT
	case oc.RoutingPolicy_MatchSetOptionsType_ANY:
		return bgpconfig.MATCH_SET_OPTIONS_TYPE_ANY
	default:
		return bgpconfig.MATCH_SET_OPTIONS_TYPE_ANY
	}
}

func convertMatchSetOptionsRestrictedType(ocrestrictedMatchSetOpts oc.E_RoutingPolicy_MatchSetOptionsRestrictedType) bgpconfig.MatchSetOptionsRestrictedType {
	switch ocrestrictedMatchSetOpts {
	case oc.RoutingPolicy_MatchSetOptionsRestrictedType_INVERT:
		return bgpconfig.MATCH_SET_OPTIONS_RESTRICTED_TYPE_INVERT
	case oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY:
		return bgpconfig.MATCH_SET_OPTIONS_RESTRICTED_TYPE_ANY
	default:
		return bgpconfig.MATCH_SET_OPTIONS_RESTRICTED_TYPE_ANY
	}
}

func convertRouteDisposition(ocpolicyresult oc.E_RoutingPolicy_PolicyResultType) bgpconfig.RouteDisposition {
	switch ocpolicyresult {
	case oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE:
		return bgpconfig.ROUTE_DISPOSITION_REJECT_ROUTE
	case oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE:
		return bgpconfig.ROUTE_DISPOSITION_ACCEPT_ROUTE
	default:
		return bgpconfig.ROUTE_DISPOSITION_NONE
	}
}

func defaultPolicyToRouteDisp(gobgpdefaultpolicy bgpconfig.DefaultPolicyType) bgpconfig.RouteDisposition {
	switch gobgpdefaultpolicy {
	case bgpconfig.DEFAULT_POLICY_TYPE_REJECT_ROUTE:
		return bgpconfig.ROUTE_DISPOSITION_REJECT_ROUTE
	case bgpconfig.DEFAULT_POLICY_TYPE_ACCEPT_ROUTE:
		return bgpconfig.ROUTE_DISPOSITION_ACCEPT_ROUTE
	default:
		return bgpconfig.ROUTE_DISPOSITION_REJECT_ROUTE
	}
}

func convertCommunity(community any) string {
	switch c := community.(type) {
	case oc.UnionString:
		return string(c)
	case oc.UnionUint32:
		uint32ToCommunityString(uint32(c))
	case oc.E_BgpTypes_BGP_WELL_KNOWN_STD_COMMUNITY:
		switch c {
		case oc.BgpTypes_BGP_WELL_KNOWN_STD_COMMUNITY_NO_EXPORT:
			return "65535:65281"
		case oc.BgpTypes_BGP_WELL_KNOWN_STD_COMMUNITY_NO_ADVERTISE:
			return "65535:65282"
		case oc.BgpTypes_BGP_WELL_KNOWN_STD_COMMUNITY_NO_EXPORT_SUBCONFED:
			return "65535:65283"
		case oc.BgpTypes_BGP_WELL_KNOWN_STD_COMMUNITY_NOPEER:
			return "65535:65284"
		}
	}
	return ""
}

func convertCommunitySet(occommset map[string]*oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet) []bgpconfig.CommunitySet {
	var commsets []bgpconfig.CommunitySet
	for communitySetName, communitySet := range occommset {
		var communityList []string
		for _, community := range communitySet.CommunityMember {
			communityList = append(communityList, convertCommunity(community))
		}

		commsets = append(commsets, bgpconfig.CommunitySet{
			CommunitySetName: communitySetName,
			CommunityList:    communityList,
		})
	}
	return commsets
}
