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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/wenovus/gobgp/v3/pkg/config/gobgp"
)

func TestIntendedToGoBGPPolicies(t *testing.T) {
	tests := []struct {
		desc          string
		inOC          *oc.Root
		wantBGPConfig *gobgp.BgpConfigSet
	}{{
		desc: "big-test",
		inOC: func() *oc.Root {
			root := &oc.Root{}
			bgpoc := root.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).GetOrCreateBgp()
			policyoc := root.GetOrCreateRoutingPolicy()

			// DEFINED SETS
			prefixSet1Name := "V4-1"
			prefix1 := "10.10.10.0/27"
			prefixSetPath := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet(prefixSet1Name)
			prefixSetPath.SetMode(oc.PrefixSet_Mode_IPV4)
			prefixSetPath.GetOrCreatePrefix(prefix1, "exact").SetIpPrefix(prefix1)

			prefixSet2Name := "V4-2"
			prefix2 := "10.20.0.0/16"
			prefixSetPath = root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet(prefixSet2Name)
			prefixSetPath.SetMode(oc.PrefixSet_Mode_IPV4)
			prefixSetPath.GetOrCreatePrefix(prefix2, "29..29").SetIpPrefix(prefix2)

			commsetName := "COMM1"
			root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreateBgpDefinedSets().GetOrCreateCommunitySet(commsetName).SetCommunityMember(
				[]oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
					oc.UnionString("12345:54321"),
				},
			)

			// POLICY
			policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			policyName := "foo"
			policyoc.GetOrCreatePolicyDefinition(policyName).Statement = policy

			v4Stmt, err := policy.AppendNew(policyName + "-1")
			if err != nil {
				t.Fatal(err)
			}
			v4Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSet1Name)
			v4Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)
			v4Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateReference().SetCommunitySetRef(commsetName)
			v4Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_REFERENCE)
			v4Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
			v4Stmt, err = policy.AppendNew(policyName + "-2")
			if err != nil {
				t.Fatal(err)
			}
			v4Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSet2Name)
			v4Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)
			v4Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_REFERENCE)
			v4Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_REPLACE)
			v4Stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)

			bgpoc.GetOrCreateNeighbor("1.1.1.1").GetOrCreateApplyPolicy().SetExportPolicy([]string{policyName})
			bgpoc.GetOrCreateNeighbor("1.1.1.1").GetOrCreateApplyPolicy().SetImportPolicy([]string{policyName})
			bgpoc.GetOrCreateNeighbor("1.1.1.1").GetOrCreateApplyPolicy().SetDefaultImportPolicy(oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)

			bgpoc.GetOrCreateNeighbor("2.2.2.2").GetOrCreateApplyPolicy().SetExportPolicy([]string{policyName})
			return root
		}(),
		wantBGPConfig: &gobgp.BgpConfigSet{
			Global: gobgp.Global{
				ApplyPolicy: gobgp.ApplyPolicy{
					Config: gobgp.ApplyPolicyConfig{
						ImportPolicyList: []string{
							"1.1.1.1|foo",
							"default-import|1.1.1.1",
							"default-import|2.2.2.2",
						},
						DefaultImportPolicy: "",
						ExportPolicyList: []string{
							"1.1.1.1|foo",
							"default-export|1.1.1.1",
							"2.2.2.2|foo",
							"default-export|2.2.2.2",
						},
						DefaultExportPolicy: "",
						InPolicyList:        []string(nil),
						DefaultInPolicy:     "",
					},
				},
			},
			DefinedSets: gobgp.DefinedSets{
				PrefixSets: []gobgp.PrefixSet{{
					PrefixSetName: "V4-1",
					PrefixList: []gobgp.Prefix{{
						IpPrefix:        "10.10.10.0/27",
						MasklengthRange: "",
					}},
				}, {
					PrefixSetName: "V4-2",
					PrefixList: []gobgp.Prefix{{
						IpPrefix:        "10.20.0.0/16",
						MasklengthRange: "29..29",
					}},
				}},
				NeighborSets: []gobgp.NeighborSet{{
					NeighborSetName:  "1.1.1.1",
					NeighborInfoList: []string{"1.1.1.1"},
				}, {
					NeighborSetName:  "2.2.2.2",
					NeighborInfoList: []string{"2.2.2.2"},
				}},
				BgpDefinedSets: gobgp.BgpDefinedSets{
					CommunitySets: []gobgp.CommunitySet{{
						CommunitySetName: "COMM1",
						CommunityList:    []string{"12345:54321"},
					}},
				},
			},
			PolicyDefinitions: []gobgp.PolicyDefinition{{
				Name: "1.1.1.1|foo",
				Statements: []gobgp.Statement{{
					Name: "foo-1",
					Conditions: gobgp.Conditions{
						CallPolicy:       "",
						MatchPrefixSet:   gobgp.MatchPrefixSet{PrefixSet: "V4-1", MatchSetOptions: "any"},
						MatchNeighborSet: gobgp.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""}, MatchTagSet: gobgp.MatchTagSet{TagSet: "", MatchSetOptions: ""},
						InstallProtocolEq: "", IgpConditions: gobgp.IgpConditions{}, BgpConditions: gobgp.BgpConditions{
							MatchCommunitySet: gobgp.MatchCommunitySet{
								CommunitySet:    "",
								MatchSetOptions: "any",
							},
							MatchExtCommunitySet: gobgp.MatchExtCommunitySet{
								ExtCommunitySet: "",
								MatchSetOptions: "",
							},
							MatchAsPathSet: gobgp.MatchAsPathSet{
								AsPathSet: "", MatchSetOptions: "any",
							},
							MedEq:                0x0,
							OriginEq:             "",
							NextHopInList:        []string(nil),
							AfiSafiInList:        []gobgp.AfiSafiType(nil),
							LocalPrefEq:          0x0,
							CommunityCount:       gobgp.CommunityCount{Operator: "", Value: 0x0},
							AsPathLength:         gobgp.AsPathLength{Operator: "", Value: 0x0},
							RouteType:            "",
							RpkiValidationResult: "",
							MatchLargeCommunitySet: gobgp.MatchLargeCommunitySet{
								LargeCommunitySet: "",
								MatchSetOptions:   "",
							},
						},
					},
					Actions: gobgp.Actions{
						RouteDisposition: "none",
						IgpActions:       gobgp.IgpActions{SetTag: ""},
						BgpActions: gobgp.BgpActions{SetAsPathPrepend: gobgp.SetAsPathPrepend{RepeatN: 0x0, As: "0"},
							SetCommunity:      gobgp.SetCommunity{SetCommunityMethod: gobgp.SetCommunityMethod{CommunitiesList: []string{"12345:54321"}, CommunitySetRef: ""}, Options: "ADD"},
							SetExtCommunity:   gobgp.SetExtCommunity{SetExtCommunityMethod: gobgp.SetExtCommunityMethod{CommunitiesList: []string(nil), ExtCommunitySetRef: ""}, Options: ""},
							SetRouteOrigin:    "",
							SetLocalPref:      0x0,
							SetNextHop:        "",
							SetMed:            "",
							SetLargeCommunity: gobgp.SetLargeCommunity{SetLargeCommunityMethod: gobgp.SetLargeCommunityMethod{CommunitiesList: []string(nil)}, Options: ""},
						},
					},
				}, {
					Name: "foo-2",
					Conditions: gobgp.Conditions{
						CallPolicy:       "",
						MatchPrefixSet:   gobgp.MatchPrefixSet{PrefixSet: "V4-2", MatchSetOptions: "any"},
						MatchNeighborSet: gobgp.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""}, MatchTagSet: gobgp.MatchTagSet{TagSet: "", MatchSetOptions: ""},
						InstallProtocolEq: "", IgpConditions: gobgp.IgpConditions{}, BgpConditions: gobgp.BgpConditions{
							MatchCommunitySet: gobgp.MatchCommunitySet{
								CommunitySet:    "",
								MatchSetOptions: "any",
							},
							MatchExtCommunitySet: gobgp.MatchExtCommunitySet{
								ExtCommunitySet: "",
								MatchSetOptions: "",
							},
							MatchAsPathSet: gobgp.MatchAsPathSet{
								AsPathSet: "", MatchSetOptions: "any",
							},
							MedEq:                0x0,
							OriginEq:             "",
							NextHopInList:        []string(nil),
							AfiSafiInList:        []gobgp.AfiSafiType(nil),
							LocalPrefEq:          0x0,
							CommunityCount:       gobgp.CommunityCount{Operator: "", Value: 0x0},
							AsPathLength:         gobgp.AsPathLength{Operator: "", Value: 0x0},
							RouteType:            "",
							RpkiValidationResult: "",
							MatchLargeCommunitySet: gobgp.MatchLargeCommunitySet{
								LargeCommunitySet: "",
								MatchSetOptions:   "",
							},
						},
					},
					Actions: gobgp.Actions{
						RouteDisposition: "accept-route",
						IgpActions:       gobgp.IgpActions{SetTag: ""},
						BgpActions: gobgp.BgpActions{SetAsPathPrepend: gobgp.SetAsPathPrepend{RepeatN: 0x0, As: "0"},
							SetCommunity:      gobgp.SetCommunity{SetCommunityMethod: gobgp.SetCommunityMethod{CommunitiesList: []string(nil), CommunitySetRef: ""}, Options: "REPLACE"},
							SetExtCommunity:   gobgp.SetExtCommunity{SetExtCommunityMethod: gobgp.SetExtCommunityMethod{CommunitiesList: []string(nil), ExtCommunitySetRef: ""}, Options: ""},
							SetRouteOrigin:    "",
							SetLocalPref:      0x0,
							SetNextHop:        "",
							SetMed:            "",
							SetLargeCommunity: gobgp.SetLargeCommunity{SetLargeCommunityMethod: gobgp.SetLargeCommunityMethod{CommunitiesList: []string(nil)}, Options: ""},
						},
					},
				}},
			}, {
				Name: "default-import|1.1.1.1",
				Statements: []gobgp.Statement{{
					Name: "default-import|1.1.1.1",
					Conditions: gobgp.Conditions{
						MatchNeighborSet: gobgp.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""},
					},
					Actions: gobgp.Actions{
						RouteDisposition: "accept-route",
					},
				}},
			}, {
				Name: "default-export|1.1.1.1",
				Statements: []gobgp.Statement{{
					Name: "default-export|1.1.1.1",
					Conditions: gobgp.Conditions{
						MatchNeighborSet: gobgp.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""},
					},
					Actions: gobgp.Actions{
						RouteDisposition: "reject-route",
					},
				}},
			}, {
				Name: "2.2.2.2|foo",
				Statements: []gobgp.Statement{{
					Name: "foo-1",
					Conditions: gobgp.Conditions{
						CallPolicy:       "",
						MatchPrefixSet:   gobgp.MatchPrefixSet{PrefixSet: "V4-1", MatchSetOptions: "any"},
						MatchNeighborSet: gobgp.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""}, MatchTagSet: gobgp.MatchTagSet{TagSet: "", MatchSetOptions: ""},
						InstallProtocolEq: "", IgpConditions: gobgp.IgpConditions{}, BgpConditions: gobgp.BgpConditions{
							MatchCommunitySet: gobgp.MatchCommunitySet{
								CommunitySet:    "",
								MatchSetOptions: "any",
							},
							MatchExtCommunitySet: gobgp.MatchExtCommunitySet{
								ExtCommunitySet: "",
								MatchSetOptions: "",
							},
							MatchAsPathSet: gobgp.MatchAsPathSet{
								AsPathSet: "", MatchSetOptions: "any",
							},
							MedEq:                0x0,
							OriginEq:             "",
							NextHopInList:        []string(nil),
							AfiSafiInList:        []gobgp.AfiSafiType(nil),
							LocalPrefEq:          0x0,
							CommunityCount:       gobgp.CommunityCount{Operator: "", Value: 0x0},
							AsPathLength:         gobgp.AsPathLength{Operator: "", Value: 0x0},
							RouteType:            "",
							RpkiValidationResult: "",
							MatchLargeCommunitySet: gobgp.MatchLargeCommunitySet{
								LargeCommunitySet: "",
								MatchSetOptions:   "",
							},
						},
					},
					Actions: gobgp.Actions{
						RouteDisposition: "none",
						IgpActions:       gobgp.IgpActions{SetTag: ""},
						BgpActions: gobgp.BgpActions{SetAsPathPrepend: gobgp.SetAsPathPrepend{RepeatN: 0x0, As: "0"},
							SetCommunity:      gobgp.SetCommunity{SetCommunityMethod: gobgp.SetCommunityMethod{CommunitiesList: []string{"12345:54321"}, CommunitySetRef: ""}, Options: "ADD"},
							SetExtCommunity:   gobgp.SetExtCommunity{SetExtCommunityMethod: gobgp.SetExtCommunityMethod{CommunitiesList: []string(nil), ExtCommunitySetRef: ""}, Options: ""},
							SetRouteOrigin:    "",
							SetLocalPref:      0x0,
							SetNextHop:        "",
							SetMed:            "",
							SetLargeCommunity: gobgp.SetLargeCommunity{SetLargeCommunityMethod: gobgp.SetLargeCommunityMethod{CommunitiesList: []string(nil)}, Options: ""},
						},
					},
				}, {
					Name: "foo-2",
					Conditions: gobgp.Conditions{
						CallPolicy:       "",
						MatchPrefixSet:   gobgp.MatchPrefixSet{PrefixSet: "V4-2", MatchSetOptions: "any"},
						MatchNeighborSet: gobgp.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""}, MatchTagSet: gobgp.MatchTagSet{TagSet: "", MatchSetOptions: ""},
						InstallProtocolEq: "", IgpConditions: gobgp.IgpConditions{}, BgpConditions: gobgp.BgpConditions{
							MatchCommunitySet: gobgp.MatchCommunitySet{
								CommunitySet:    "",
								MatchSetOptions: "any",
							},
							MatchExtCommunitySet: gobgp.MatchExtCommunitySet{
								ExtCommunitySet: "",
								MatchSetOptions: "",
							},
							MatchAsPathSet: gobgp.MatchAsPathSet{
								AsPathSet: "", MatchSetOptions: "any",
							},
							MedEq:                0x0,
							OriginEq:             "",
							NextHopInList:        []string(nil),
							AfiSafiInList:        []gobgp.AfiSafiType(nil),
							LocalPrefEq:          0x0,
							CommunityCount:       gobgp.CommunityCount{Operator: "", Value: 0x0},
							AsPathLength:         gobgp.AsPathLength{Operator: "", Value: 0x0},
							RouteType:            "",
							RpkiValidationResult: "",
							MatchLargeCommunitySet: gobgp.MatchLargeCommunitySet{
								LargeCommunitySet: "",
								MatchSetOptions:   "",
							},
						},
					},
					Actions: gobgp.Actions{
						RouteDisposition: "accept-route",
						IgpActions:       gobgp.IgpActions{SetTag: ""},
						BgpActions: gobgp.BgpActions{SetAsPathPrepend: gobgp.SetAsPathPrepend{RepeatN: 0x0, As: "0"},
							SetCommunity:      gobgp.SetCommunity{SetCommunityMethod: gobgp.SetCommunityMethod{CommunitiesList: []string(nil), CommunitySetRef: ""}, Options: "REPLACE"},
							SetExtCommunity:   gobgp.SetExtCommunity{SetExtCommunityMethod: gobgp.SetExtCommunityMethod{CommunitiesList: []string(nil), ExtCommunitySetRef: ""}, Options: ""},
							SetRouteOrigin:    "",
							SetLocalPref:      0x0,
							SetNextHop:        "",
							SetMed:            "",
							SetLargeCommunity: gobgp.SetLargeCommunity{SetLargeCommunityMethod: gobgp.SetLargeCommunityMethod{CommunitiesList: []string(nil)}, Options: ""},
						},
					},
				}},
			}, {
				Name: "default-import|2.2.2.2",
				Statements: []gobgp.Statement{{
					Name: "default-import|2.2.2.2",
					Conditions: gobgp.Conditions{
						MatchNeighborSet: gobgp.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""},
					},
					Actions: gobgp.Actions{
						RouteDisposition: "reject-route",
					},
				}},
			}, {
				Name: "default-export|2.2.2.2",
				Statements: []gobgp.Statement{{
					Name: "default-export|2.2.2.2",
					Conditions: gobgp.Conditions{
						MatchNeighborSet: gobgp.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""},
					},
					Actions: gobgp.Actions{
						RouteDisposition: "reject-route",
					},
				}},
			}},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gotBGPConfig := &gobgp.BgpConfigSet{}
			intendedToGoBGPPolicies(
				tt.inOC.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).GetOrCreateBgp(),
				tt.inOC.GetOrCreateRoutingPolicy(),
				gotBGPConfig,
			)
			if diff := cmp.Diff(tt.wantBGPConfig, gotBGPConfig); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}
		})
	}
}
