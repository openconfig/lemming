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
	"net/netip"
	"slices"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/internal/lemmingutil"
	gobgpoc "github.com/osrg/gobgp/v3/pkg/config/oc"
	"github.com/osrg/gobgp/v3/pkg/zebra"
)

// intendedToGoBGP translates from OC to GoBGP intended config.
//
// GoBGP's notion of config vs. state does not conform to OpenConfig (see
// https://github.com/osrg/gobgp/issues/2584)
// Therefore, we need a compatibility layer between the two configs.
func intendedToGoBGP(bgpoc *oc.NetworkInstance_Protocol_Bgp, policyoc *oc.RoutingPolicy, zapiURL string, listenPort uint16) *gobgpoc.BgpConfigSet {
	bgpConfig := &gobgpoc.BgpConfigSet{}

	// Global config
	global := bgpoc.GetOrCreateGlobal()

	bgpConfig.Global.Config.As = global.GetAs()
	bgpConfig.Global.Config.RouterId = global.GetRouterId()
	bgpConfig.Global.Config.Port = int32(listenPort)

	if localAddr, err := netip.ParseAddr(global.GetRouterId()); err == nil && localAddr.IsLoopback() {
		// Have GoBGP listen only on local address instead of all
		// addresses when testing BGP server on localhost.
		bgpConfig.Global.Config.LocalAddressList = []string{localAddr.String()}
	}

	for neighAddr, neigh := range bgpoc.Neighbor {
		// Add neighbour config.
		bgpConfig.Neighbors = append(bgpConfig.Neighbors, gobgpoc.Neighbor{
			Config: gobgpoc.NeighborConfig{
				PeerAs:          neigh.GetPeerAs(),
				NeighborAddress: neighAddr,
			},
			// This is needed because GoBGP's configuration diffing
			// logic may check the state value instead of the
			// config value.
			State: gobgpoc.NeighborState{
				PeerAs:          neigh.GetPeerAs(),
				NeighborAddress: neighAddr,
			},
			Transport: gobgpoc.Transport{
				Config: gobgpoc.TransportConfig{
					LocalAddress: neigh.GetTransport().GetLocalAddress(),
					RemotePort:   neigh.GetNeighborPort(),
				},
			},
			AsPathOptions: gobgpoc.AsPathOptions{
				Config: gobgpoc.AsPathOptionsConfig{
					AllowOwnAs:    neigh.GetAsPathOptions().GetAllowOwnAs(),
					ReplacePeerAs: neigh.GetAsPathOptions().GetReplacePeerAs(),
				},
			},
		})
	}

	intendedToGoBGPPolicies(bgpoc, policyoc, bgpConfig)

	bgpConfig.Zebra.Config = gobgpoc.ZebraConfig{
		Enabled: true,
		Url:     zapiURL,
		// TODO(wenbli): This should actually be filled with the types
		// of routes it wants redistributed instead of getting all
		// routes.
		RedistributeRouteTypeList: []string{},
		Version:                   zebra.MaxZapiVer,
		NexthopTriggerEnable:      false,
		SoftwareName:              "frr8.2",
	}

	return bgpConfig
}

// intendedToGoBGPPolicies populates bgpConfig's policies from the OC configuration.
func intendedToGoBGPPolicies(bgpoc *oc.NetworkInstance_Protocol_Bgp, policyoc *oc.RoutingPolicy, bgpConfig *gobgpoc.BgpConfigSet) {
	var communitySetIndexMap map[string]int
	// community sets
	bgpConfig.DefinedSets.BgpDefinedSets.CommunitySets, communitySetIndexMap = convertCommunitySet(policyoc.GetOrCreateDefinedSets().GetOrCreateBgpDefinedSets().CommunitySet)
	// Prefix sets
	bgpConfig.DefinedSets.PrefixSets = convertPrefixSets(policyoc.GetOrCreateDefinedSets().PrefixSet)
	// AS Path Sets
	bgpConfig.DefinedSets.BgpDefinedSets.AsPathSets = convertASPathSets(policyoc.GetOrCreateDefinedSets().GetOrCreateBgpDefinedSets().AsPathSet)

	neighAddrs := lemmingutil.Mapkeys(bgpoc.Neighbor)
	slices.Sort(neighAddrs)

	// Neighbours, global policy definitions, and global apply policy list.
	for _, neighAddr := range neighAddrs {
		// Ideally a simple conversion of apply-policy is sufficient, but due to GoBGP using
		// a global set of apply-policy instead of per-neighbour policies, we need to create
		// neighbour sets and modify input policy statements so that we retain the same
		// per-neighbour behaviour while only using a single set of global policies.
		//
		// To do this, we create a neighbour set for each neighbour containing just the
		// single neighbour address, then duplicate the policies to make a copy for each
		// neighbour that uses it, and then concatenate the ApplyPolicy lists of every
		// neighbour's ApplyPolicy into the global ApplyPolicy list.
		//
		// The resulting policies is of the following form:
		// Neighbour sets: [neigh1, neigh2, neigh3, ...]
		// PolicyDefinitions: [neigh1polA, neigh1polB, ..., neigh1default-import, neigh1default-export,
		//                     neigh2polA, neigh2polB, ..., ...
		//                     ...]
		// Global ApplyPolicy list: [same as policy-definitions]
		bgpConfig.DefinedSets.NeighborSets = append(bgpConfig.DefinedSets.NeighborSets, gobgpoc.NeighborSet{
			NeighborSetName:  neighAddr,
			NeighborInfoList: []string{neighAddr},
		})

		applyPolicy := convertNeighborApplyPolicy(bgpoc.Neighbor[neighAddr])

		// populatePolicies populates the global policy definitions and the ApplyPolicy
		// list, and returns the list of converted policies' names.
		policies := map[string]bool{}
		populatePolicies := func(policyList []string) []string {
			var applyPolicyList []string
			for _, policyName := range policyList {
				convertedPolicyName := convertPolicyName(neighAddr, policyName)
				if policies[policyName] {
					// Already processed
					applyPolicyList = append(applyPolicyList, convertedPolicyName)
					continue
				}
				policies[policyName] = true
				policy, ok := policyoc.PolicyDefinition[policyName]
				if !ok {
					log.Errorf("Neighbour policy doesn't exist in policy definitions: %q", policyName)
					continue
				}
				convertedPolicy := convertPolicyDefinition(policy, neighAddr, policyoc.GetOrCreateDefinedSets().GetOrCreateBgpDefinedSets().CommunitySet, bgpConfig.DefinedSets.BgpDefinedSets.CommunitySets, communitySetIndexMap)
				bgpConfig.PolicyDefinitions = append(bgpConfig.PolicyDefinitions, convertedPolicy)
				applyPolicyList = append(applyPolicyList, convertedPolicyName)
			}
			return applyPolicyList
		}
		bgpConfig.Global.ApplyPolicy.Config.ImportPolicyList = append(bgpConfig.Global.ApplyPolicy.Config.ImportPolicyList, populatePolicies(applyPolicy.Config.ImportPolicyList)...)
		bgpConfig.Global.ApplyPolicy.Config.ExportPolicyList = append(bgpConfig.Global.ApplyPolicy.Config.ExportPolicyList, populatePolicies(applyPolicy.Config.ExportPolicyList)...)

		// Create per-neighbour default policies.
		defaultImportPolicyName := "default-import|" + neighAddr
		defaultExportPolicyName := "default-export|" + neighAddr
		bgpConfig.PolicyDefinitions = append(bgpConfig.PolicyDefinitions, gobgpoc.PolicyDefinition{
			Name: defaultImportPolicyName,
			Statements: []gobgpoc.Statement{{
				// Use a customized name for the default policies.
				Name: defaultImportPolicyName,
				Conditions: gobgpoc.Conditions{
					MatchNeighborSet: gobgpoc.MatchNeighborSet{
						NeighborSet: neighAddr,
					},
				},
				Actions: gobgpoc.Actions{
					RouteDisposition: defaultPolicyToRouteDisp(applyPolicy.Config.DefaultImportPolicy),
				},
			}},
		}, gobgpoc.PolicyDefinition{
			Name: defaultExportPolicyName,
			Statements: []gobgpoc.Statement{{
				// Use a customized name for the default policies.
				Name: defaultExportPolicyName,
				Conditions: gobgpoc.Conditions{
					MatchNeighborSet: gobgpoc.MatchNeighborSet{
						NeighborSet: neighAddr,
					},
				},
				Actions: gobgpoc.Actions{
					RouteDisposition: defaultPolicyToRouteDisp(applyPolicy.Config.DefaultExportPolicy),
				},
			}},
		})
		bgpConfig.Global.ApplyPolicy.Config.ImportPolicyList = append(bgpConfig.Global.ApplyPolicy.Config.ImportPolicyList, defaultImportPolicyName)
		bgpConfig.Global.ApplyPolicy.Config.ExportPolicyList = append(bgpConfig.Global.ApplyPolicy.Config.ExportPolicyList, defaultExportPolicyName)
	}
}
