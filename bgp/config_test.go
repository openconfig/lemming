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
	gobgpoc "github.com/wenovus/gobgp/v3/pkg/config/oc"
)

func TestIntendedToGoBGPPolicies(t *testing.T) {
	tests := []struct {
		desc          string
		inOC          *oc.Root
		wantBGPConfig *gobgpoc.BgpConfigSet
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
					oc.UnionString("[0-9]+:[0-9]+"),
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
			v4Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
			v4Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateReference().SetCommunitySetRef(commsetName)
			v4Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_REFERENCE)
			v4Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
			v4Stmt, err = policy.AppendNew(policyName + "-2")
			if err != nil {
				t.Fatal(err)
			}
			v4Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSet2Name)
			v4Stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
			v4Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_REFERENCE)
			v4Stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_REPLACE)
			v4Stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)

			bgpoc.GetOrCreateNeighbor("1.1.1.1").GetOrCreateApplyPolicy().SetExportPolicy([]string{policyName})
			bgpoc.GetOrCreateNeighbor("1.1.1.1").GetOrCreateApplyPolicy().SetImportPolicy([]string{policyName})
			bgpoc.GetOrCreateNeighbor("1.1.1.1").GetOrCreateApplyPolicy().SetDefaultImportPolicy(oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)

			bgpoc.GetOrCreateNeighbor("2.2.2.2").GetOrCreateApplyPolicy().SetExportPolicy([]string{policyName})
			return root
		}(),
		wantBGPConfig: &gobgpoc.BgpConfigSet{
			Global: gobgpoc.Global{
				ApplyPolicy: gobgpoc.ApplyPolicy{
					Config: gobgpoc.ApplyPolicyConfig{
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
			DefinedSets: gobgpoc.DefinedSets{
				PrefixSets: []gobgpoc.PrefixSet{{
					PrefixSetName: "V4-1",
					PrefixList: []gobgpoc.Prefix{{
						IpPrefix:        "10.10.10.0/27",
						MasklengthRange: "",
					}},
				}, {
					PrefixSetName: "V4-2",
					PrefixList: []gobgpoc.Prefix{{
						IpPrefix:        "10.20.0.0/16",
						MasklengthRange: "29..29",
					}},
				}},
				NeighborSets: []gobgpoc.NeighborSet{{
					NeighborSetName:  "1.1.1.1",
					NeighborInfoList: []string{"1.1.1.1"},
				}, {
					NeighborSetName:  "2.2.2.2",
					NeighborInfoList: []string{"2.2.2.2"},
				}},
				BgpDefinedSets: gobgpoc.BgpDefinedSets{
					CommunitySets: []gobgpoc.CommunitySet{{
						CommunitySetName: "COMM1",
						CommunityList:    []string{"[0-9]+:[0-9]+"},
					}},
				},
			},
			PolicyDefinitions: []gobgpoc.PolicyDefinition{{
				Name: "1.1.1.1|foo",
				Statements: []gobgpoc.Statement{{
					Name: "1.1.1.1|foo:foo-1",
					Conditions: gobgpoc.Conditions{
						CallPolicy:       "",
						MatchPrefixSet:   gobgpoc.MatchPrefixSet{PrefixSet: "V4-1", MatchSetOptions: "any"},
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""}, MatchTagSet: gobgpoc.MatchTagSet{TagSet: "", MatchSetOptions: ""},
						InstallProtocolEq: "", IgpConditions: gobgpoc.IgpConditions{}, BgpConditions: gobgpoc.BgpConditions{
							MatchCommunitySet: gobgpoc.MatchCommunitySet{
								CommunitySet:    "",
								MatchSetOptions: "any",
							},
							MatchExtCommunitySet: gobgpoc.MatchExtCommunitySet{
								ExtCommunitySet: "",
								MatchSetOptions: "",
							},
							MatchAsPathSet: gobgpoc.MatchAsPathSet{
								AsPathSet: "", MatchSetOptions: "any",
							},
							MedEq:                0x0,
							OriginEq:             "",
							NextHopInList:        []string(nil),
							AfiSafiInList:        []gobgpoc.AfiSafiType(nil),
							LocalPrefEq:          0x0,
							CommunityCount:       gobgpoc.CommunityCount{Operator: "", Value: 0x0},
							AsPathLength:         gobgpoc.AsPathLength{Operator: "", Value: 0x0},
							RouteType:            "",
							RpkiValidationResult: "",
							MatchLargeCommunitySet: gobgpoc.MatchLargeCommunitySet{
								LargeCommunitySet: "",
								MatchSetOptions:   "",
							},
						},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "none",
						IgpActions:       gobgpoc.IgpActions{SetTag: ""},
						BgpActions: gobgpoc.BgpActions{SetAsPathPrepend: gobgpoc.SetAsPathPrepend{RepeatN: 0x0, As: "0"},
							SetCommunity:      gobgpoc.SetCommunity{SetCommunityMethod: gobgpoc.SetCommunityMethod{CommunitiesList: []string{"[0-9]+:[0-9]+"}, CommunitySetRef: ""}, Options: "add"},
							SetExtCommunity:   gobgpoc.SetExtCommunity{SetExtCommunityMethod: gobgpoc.SetExtCommunityMethod{CommunitiesList: []string(nil), ExtCommunitySetRef: ""}, Options: ""},
							SetRouteOrigin:    "",
							SetLocalPref:      0x0,
							SetNextHop:        "",
							SetMed:            "",
							SetLargeCommunity: gobgpoc.SetLargeCommunity{SetLargeCommunityMethod: gobgpoc.SetLargeCommunityMethod{CommunitiesList: []string(nil)}, Options: ""},
						},
					},
				}, {
					Name: "1.1.1.1|foo:foo-2",
					Conditions: gobgpoc.Conditions{
						CallPolicy:       "",
						MatchPrefixSet:   gobgpoc.MatchPrefixSet{PrefixSet: "V4-2", MatchSetOptions: "any"},
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""}, MatchTagSet: gobgpoc.MatchTagSet{TagSet: "", MatchSetOptions: ""},
						InstallProtocolEq: "", IgpConditions: gobgpoc.IgpConditions{}, BgpConditions: gobgpoc.BgpConditions{
							MatchCommunitySet: gobgpoc.MatchCommunitySet{
								CommunitySet:    "",
								MatchSetOptions: "any",
							},
							MatchExtCommunitySet: gobgpoc.MatchExtCommunitySet{
								ExtCommunitySet: "",
								MatchSetOptions: "",
							},
							MatchAsPathSet: gobgpoc.MatchAsPathSet{
								AsPathSet: "", MatchSetOptions: "any",
							},
							MedEq:                0x0,
							OriginEq:             "",
							NextHopInList:        []string(nil),
							AfiSafiInList:        []gobgpoc.AfiSafiType(nil),
							LocalPrefEq:          0x0,
							CommunityCount:       gobgpoc.CommunityCount{Operator: "", Value: 0x0},
							AsPathLength:         gobgpoc.AsPathLength{Operator: "", Value: 0x0},
							RouteType:            "",
							RpkiValidationResult: "",
							MatchLargeCommunitySet: gobgpoc.MatchLargeCommunitySet{
								LargeCommunitySet: "",
								MatchSetOptions:   "",
							},
						},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "accept-route",
						IgpActions:       gobgpoc.IgpActions{SetTag: ""},
						BgpActions: gobgpoc.BgpActions{SetAsPathPrepend: gobgpoc.SetAsPathPrepend{RepeatN: 0x0, As: "0"},
							SetCommunity:      gobgpoc.SetCommunity{SetCommunityMethod: gobgpoc.SetCommunityMethod{CommunitiesList: []string(nil), CommunitySetRef: ""}, Options: "replace"},
							SetExtCommunity:   gobgpoc.SetExtCommunity{SetExtCommunityMethod: gobgpoc.SetExtCommunityMethod{CommunitiesList: []string(nil), ExtCommunitySetRef: ""}, Options: ""},
							SetRouteOrigin:    "",
							SetLocalPref:      0x0,
							SetNextHop:        "",
							SetMed:            "",
							SetLargeCommunity: gobgpoc.SetLargeCommunity{SetLargeCommunityMethod: gobgpoc.SetLargeCommunityMethod{CommunitiesList: []string(nil)}, Options: ""},
						},
					},
				}},
			}, {
				Name: "default-import|1.1.1.1",
				Statements: []gobgpoc.Statement{{
					Name: "default-import|1.1.1.1",
					Conditions: gobgpoc.Conditions{
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "accept-route",
					},
				}},
			}, {
				Name: "default-export|1.1.1.1",
				Statements: []gobgpoc.Statement{{
					Name: "default-export|1.1.1.1",
					Conditions: gobgpoc.Conditions{
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "reject-route",
					},
				}},
			}, {
				Name: "2.2.2.2|foo",
				Statements: []gobgpoc.Statement{{
					Name: "2.2.2.2|foo:foo-1",
					Conditions: gobgpoc.Conditions{
						CallPolicy:       "",
						MatchPrefixSet:   gobgpoc.MatchPrefixSet{PrefixSet: "V4-1", MatchSetOptions: "any"},
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""}, MatchTagSet: gobgpoc.MatchTagSet{TagSet: "", MatchSetOptions: ""},
						InstallProtocolEq: "", IgpConditions: gobgpoc.IgpConditions{}, BgpConditions: gobgpoc.BgpConditions{
							MatchCommunitySet: gobgpoc.MatchCommunitySet{
								CommunitySet:    "",
								MatchSetOptions: "any",
							},
							MatchExtCommunitySet: gobgpoc.MatchExtCommunitySet{
								ExtCommunitySet: "",
								MatchSetOptions: "",
							},
							MatchAsPathSet: gobgpoc.MatchAsPathSet{
								AsPathSet: "", MatchSetOptions: "any",
							},
							MedEq:                0x0,
							OriginEq:             "",
							NextHopInList:        []string(nil),
							AfiSafiInList:        []gobgpoc.AfiSafiType(nil),
							LocalPrefEq:          0x0,
							CommunityCount:       gobgpoc.CommunityCount{Operator: "", Value: 0x0},
							AsPathLength:         gobgpoc.AsPathLength{Operator: "", Value: 0x0},
							RouteType:            "",
							RpkiValidationResult: "",
							MatchLargeCommunitySet: gobgpoc.MatchLargeCommunitySet{
								LargeCommunitySet: "",
								MatchSetOptions:   "",
							},
						},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "none",
						IgpActions:       gobgpoc.IgpActions{SetTag: ""},
						BgpActions: gobgpoc.BgpActions{SetAsPathPrepend: gobgpoc.SetAsPathPrepend{RepeatN: 0x0, As: "0"},
							SetCommunity:      gobgpoc.SetCommunity{SetCommunityMethod: gobgpoc.SetCommunityMethod{CommunitiesList: []string{"[0-9]+:[0-9]+"}, CommunitySetRef: ""}, Options: "add"},
							SetExtCommunity:   gobgpoc.SetExtCommunity{SetExtCommunityMethod: gobgpoc.SetExtCommunityMethod{CommunitiesList: []string(nil), ExtCommunitySetRef: ""}, Options: ""},
							SetRouteOrigin:    "",
							SetLocalPref:      0x0,
							SetNextHop:        "",
							SetMed:            "",
							SetLargeCommunity: gobgpoc.SetLargeCommunity{SetLargeCommunityMethod: gobgpoc.SetLargeCommunityMethod{CommunitiesList: []string(nil)}, Options: ""},
						},
					},
				}, {
					Name: "2.2.2.2|foo:foo-2",
					Conditions: gobgpoc.Conditions{
						CallPolicy:       "",
						MatchPrefixSet:   gobgpoc.MatchPrefixSet{PrefixSet: "V4-2", MatchSetOptions: "any"},
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""}, MatchTagSet: gobgpoc.MatchTagSet{TagSet: "", MatchSetOptions: ""},
						InstallProtocolEq: "", IgpConditions: gobgpoc.IgpConditions{}, BgpConditions: gobgpoc.BgpConditions{
							MatchCommunitySet: gobgpoc.MatchCommunitySet{
								CommunitySet:    "",
								MatchSetOptions: "any",
							},
							MatchExtCommunitySet: gobgpoc.MatchExtCommunitySet{
								ExtCommunitySet: "",
								MatchSetOptions: "",
							},
							MatchAsPathSet: gobgpoc.MatchAsPathSet{
								AsPathSet: "", MatchSetOptions: "any",
							},
							MedEq:                0x0,
							OriginEq:             "",
							NextHopInList:        []string(nil),
							AfiSafiInList:        []gobgpoc.AfiSafiType(nil),
							LocalPrefEq:          0x0,
							CommunityCount:       gobgpoc.CommunityCount{Operator: "", Value: 0x0},
							AsPathLength:         gobgpoc.AsPathLength{Operator: "", Value: 0x0},
							RouteType:            "",
							RpkiValidationResult: "",
							MatchLargeCommunitySet: gobgpoc.MatchLargeCommunitySet{
								LargeCommunitySet: "",
								MatchSetOptions:   "",
							},
						},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "accept-route",
						IgpActions:       gobgpoc.IgpActions{SetTag: ""},
						BgpActions: gobgpoc.BgpActions{SetAsPathPrepend: gobgpoc.SetAsPathPrepend{RepeatN: 0x0, As: "0"},
							SetCommunity:      gobgpoc.SetCommunity{SetCommunityMethod: gobgpoc.SetCommunityMethod{CommunitiesList: []string(nil), CommunitySetRef: ""}, Options: "replace"},
							SetExtCommunity:   gobgpoc.SetExtCommunity{SetExtCommunityMethod: gobgpoc.SetExtCommunityMethod{CommunitiesList: []string(nil), ExtCommunitySetRef: ""}, Options: ""},
							SetRouteOrigin:    "",
							SetLocalPref:      0x0,
							SetNextHop:        "",
							SetMed:            "",
							SetLargeCommunity: gobgpoc.SetLargeCommunity{SetLargeCommunityMethod: gobgpoc.SetLargeCommunityMethod{CommunitiesList: []string(nil)}, Options: ""},
						},
					},
				}},
			}, {
				Name: "default-import|2.2.2.2",
				Statements: []gobgpoc.Statement{{
					Name: "default-import|2.2.2.2",
					Conditions: gobgpoc.Conditions{
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "reject-route",
					},
				}},
			}, {
				Name: "default-export|2.2.2.2",
				Statements: []gobgpoc.Statement{{
					Name: "default-export|2.2.2.2",
					Conditions: gobgpoc.Conditions{
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "reject-route",
					},
				}},
			}},
		},
	}, {
		desc: "remove-community-set-two-statements",
		inOC: func() *oc.Root {
			root := &oc.Root{}
			bgpoc := root.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).GetOrCreateBgp()
			policyoc := root.GetOrCreateRoutingPolicy()

			// Create prefix set
			prefixSetName := "prefixset-foo"
			prefixSet := root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreatePrefixSet(prefixSetName)
			prefixSet.SetMode(oc.PrefixSet_Mode_IPV4)
			prefixSet.GetOrCreatePrefix("10.0.0.0/10", "8..32")

			// DEFINED SETS
			commsetName := "COMM1"
			root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreateBgpDefinedSets().GetOrCreateCommunitySet(commsetName).SetCommunityMember(
				[]oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
					oc.UnionString("54321:54321"),
				},
			)

			commsetName2 := "COMM2"
			root.GetOrCreateRoutingPolicy().GetOrCreateDefinedSets().GetOrCreateBgpDefinedSets().GetOrCreateCommunitySet(commsetName2).SetCommunityMember(
				[]oc.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
					oc.UnionString("12345:12345"),
				},
			)

			// POLICY
			stmtName := "stmt"
			policy1 := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			policy1Name := "foo"
			policyoc.GetOrCreatePolicyDefinition(policy1Name).Statement = policy1

			stmt, err := policy1.AppendNew(stmtName)
			if err != nil {
				t.Fatalf("Cannot append new BGP policy statement: %v", err)
			}
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)

			stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
			stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
				[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
					oc.UnionString("11111:11111"),
					oc.UnionString("22222:22222"),
				},
			)
			stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_INLINE)

			policy2 := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
			policy2Name := "bar"
			policyoc.GetOrCreatePolicyDefinition(policy2Name).Statement = policy2

			stmt, err = policy2.AppendNew(stmtName)
			if err != nil {
				t.Fatalf("Cannot append new BGP policy statement: %v", err)
			}
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
			stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_REMOVE)
			stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateReference().SetCommunitySetRefs([]string{commsetName, commsetName2})
			stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetMethod(oc.SetCommunity_Method_REFERENCE)
			stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)

			bgpoc.GetOrCreateNeighbor("1.1.1.1").GetOrCreateApplyPolicy().SetExportPolicy([]string{policy1Name})
			bgpoc.GetOrCreateNeighbor("2.2.2.2").GetOrCreateApplyPolicy().SetExportPolicy([]string{policy2Name})
			return root
		}(),
		wantBGPConfig: &gobgpoc.BgpConfigSet{
			Global: gobgpoc.Global{
				ApplyPolicy: gobgpoc.ApplyPolicy{
					Config: gobgpoc.ApplyPolicyConfig{
						ImportPolicyList: []string{
							"default-import|1.1.1.1",
							"default-import|2.2.2.2",
						},
						DefaultImportPolicy: "",
						ExportPolicyList: []string{
							"1.1.1.1|foo",
							"default-export|1.1.1.1",
							"2.2.2.2|bar",
							"default-export|2.2.2.2",
						},
						DefaultExportPolicy: "",
						InPolicyList:        []string(nil),
						DefaultInPolicy:     "",
					},
				},
			},
			DefinedSets: gobgpoc.DefinedSets{
				PrefixSets: []gobgpoc.PrefixSet{{
					PrefixSetName: "prefixset-foo",
					PrefixList: []gobgpoc.Prefix{{
						IpPrefix:        "10.0.0.0/10",
						MasklengthRange: "8..32",
					}},
				}},
				NeighborSets: []gobgpoc.NeighborSet{{
					NeighborSetName:  "1.1.1.1",
					NeighborInfoList: []string{"1.1.1.1"},
				}, {
					NeighborSetName:  "2.2.2.2",
					NeighborInfoList: []string{"2.2.2.2"},
				}},
				BgpDefinedSets: gobgpoc.BgpDefinedSets{
					CommunitySets: []gobgpoc.CommunitySet{{
						CommunitySetName: "COMM1",
						CommunityList:    []string{"54321:54321"},
					}, {
						CommunitySetName: "COMM2",
						CommunityList:    []string{"12345:12345"},
					}},
				},
			},
			PolicyDefinitions: []gobgpoc.PolicyDefinition{{
				Name: "1.1.1.1|foo",
				Statements: []gobgpoc.Statement{{
					Name: "1.1.1.1|foo:stmt",
					Conditions: gobgpoc.Conditions{
						CallPolicy:       "",
						MatchPrefixSet:   gobgpoc.MatchPrefixSet{PrefixSet: "prefixset-foo", MatchSetOptions: "any"},
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""}, MatchTagSet: gobgpoc.MatchTagSet{TagSet: "", MatchSetOptions: ""},
						InstallProtocolEq: "", IgpConditions: gobgpoc.IgpConditions{}, BgpConditions: gobgpoc.BgpConditions{
							MatchCommunitySet: gobgpoc.MatchCommunitySet{
								CommunitySet:    "",
								MatchSetOptions: "any",
							},
							MatchExtCommunitySet: gobgpoc.MatchExtCommunitySet{
								ExtCommunitySet: "",
								MatchSetOptions: "",
							},
							MatchAsPathSet: gobgpoc.MatchAsPathSet{
								AsPathSet: "", MatchSetOptions: "any",
							},
							MedEq:                0x0,
							OriginEq:             "",
							NextHopInList:        []string(nil),
							AfiSafiInList:        []gobgpoc.AfiSafiType(nil),
							LocalPrefEq:          0x0,
							CommunityCount:       gobgpoc.CommunityCount{Operator: "", Value: 0x0},
							AsPathLength:         gobgpoc.AsPathLength{Operator: "", Value: 0x0},
							RouteType:            "",
							RpkiValidationResult: "",
							MatchLargeCommunitySet: gobgpoc.MatchLargeCommunitySet{
								LargeCommunitySet: "",
								MatchSetOptions:   "",
							},
						},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "none",
						IgpActions:       gobgpoc.IgpActions{SetTag: ""},
						BgpActions: gobgpoc.BgpActions{SetAsPathPrepend: gobgpoc.SetAsPathPrepend{RepeatN: 0x0, As: "0"},
							SetCommunity:      gobgpoc.SetCommunity{SetCommunityMethod: gobgpoc.SetCommunityMethod{CommunitiesList: []string{"11111:11111", "22222:22222"}, CommunitySetRef: ""}, Options: "add"},
							SetExtCommunity:   gobgpoc.SetExtCommunity{SetExtCommunityMethod: gobgpoc.SetExtCommunityMethod{CommunitiesList: []string(nil), ExtCommunitySetRef: ""}, Options: ""},
							SetRouteOrigin:    "",
							SetLocalPref:      0x0,
							SetNextHop:        "",
							SetMed:            "",
							SetLargeCommunity: gobgpoc.SetLargeCommunity{SetLargeCommunityMethod: gobgpoc.SetLargeCommunityMethod{CommunitiesList: []string(nil)}, Options: ""},
						},
					},
				}},
			}, {
				Name: "default-import|1.1.1.1",
				Statements: []gobgpoc.Statement{{
					Name: "default-import|1.1.1.1",
					Conditions: gobgpoc.Conditions{
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "reject-route",
					},
				}},
			}, {
				Name: "default-export|1.1.1.1",
				Statements: []gobgpoc.Statement{{
					Name: "default-export|1.1.1.1",
					Conditions: gobgpoc.Conditions{
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "1.1.1.1", MatchSetOptions: ""},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "reject-route",
					},
				}},
			}, {
				Name: "2.2.2.2|bar",
				Statements: []gobgpoc.Statement{{
					Name: "2.2.2.2|bar:stmt",
					Conditions: gobgpoc.Conditions{
						CallPolicy:       "",
						MatchPrefixSet:   gobgpoc.MatchPrefixSet{PrefixSet: "prefixset-foo", MatchSetOptions: "any"},
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""}, MatchTagSet: gobgpoc.MatchTagSet{TagSet: "", MatchSetOptions: ""},
						InstallProtocolEq: "", IgpConditions: gobgpoc.IgpConditions{}, BgpConditions: gobgpoc.BgpConditions{
							MatchCommunitySet: gobgpoc.MatchCommunitySet{
								CommunitySet:    "",
								MatchSetOptions: "any",
							},
							MatchExtCommunitySet: gobgpoc.MatchExtCommunitySet{
								ExtCommunitySet: "",
								MatchSetOptions: "",
							},
							MatchAsPathSet: gobgpoc.MatchAsPathSet{
								AsPathSet: "", MatchSetOptions: "any",
							},
							MedEq:                0x0,
							OriginEq:             "",
							NextHopInList:        []string(nil),
							AfiSafiInList:        []gobgpoc.AfiSafiType(nil),
							LocalPrefEq:          0x0,
							CommunityCount:       gobgpoc.CommunityCount{Operator: "", Value: 0x0},
							AsPathLength:         gobgpoc.AsPathLength{Operator: "", Value: 0x0},
							RouteType:            "",
							RpkiValidationResult: "",
							MatchLargeCommunitySet: gobgpoc.MatchLargeCommunitySet{
								LargeCommunitySet: "",
								MatchSetOptions:   "",
							},
						},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "accept-route",
						IgpActions:       gobgpoc.IgpActions{SetTag: ""},
						BgpActions: gobgpoc.BgpActions{SetAsPathPrepend: gobgpoc.SetAsPathPrepend{RepeatN: 0x0, As: "0"},
							SetCommunity:      gobgpoc.SetCommunity{SetCommunityMethod: gobgpoc.SetCommunityMethod{CommunitiesList: []string{"54321:54321", "12345:12345"}, CommunitySetRef: ""}, Options: "remove"},
							SetExtCommunity:   gobgpoc.SetExtCommunity{SetExtCommunityMethod: gobgpoc.SetExtCommunityMethod{CommunitiesList: []string(nil), ExtCommunitySetRef: ""}, Options: ""},
							SetRouteOrigin:    "",
							SetLocalPref:      0x0,
							SetNextHop:        "",
							SetMed:            "",
							SetLargeCommunity: gobgpoc.SetLargeCommunity{SetLargeCommunityMethod: gobgpoc.SetLargeCommunityMethod{CommunitiesList: []string(nil)}, Options: ""},
						},
					},
				}},
			}, {
				Name: "default-import|2.2.2.2",
				Statements: []gobgpoc.Statement{{
					Name: "default-import|2.2.2.2",
					Conditions: gobgpoc.Conditions{
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "reject-route",
					},
				}},
			}, {
				Name: "default-export|2.2.2.2",
				Statements: []gobgpoc.Statement{{
					Name: "default-export|2.2.2.2",
					Conditions: gobgpoc.Conditions{
						MatchNeighborSet: gobgpoc.MatchNeighborSet{NeighborSet: "2.2.2.2", MatchSetOptions: ""},
					},
					Actions: gobgpoc.Actions{
						RouteDisposition: "reject-route",
					},
				}},
			}},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gotBGPConfig := &gobgpoc.BgpConfigSet{}
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
