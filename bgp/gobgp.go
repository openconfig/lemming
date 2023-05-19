// Copyright 2022 Google LLC
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
	"context"
	"errors"
	"fmt"
	"net/netip"
	"reflect"
	"sync"
	"time"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/gnmi/reconciler"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	api "github.com/wenovus/gobgp/v3/api"
	"github.com/wenovus/gobgp/v3/pkg/bgpconfig"
	"github.com/wenovus/gobgp/v3/pkg/server"
	"github.com/wenovus/gobgp/v3/pkg/zebra"
)

const (
	gracefulRestart = false
)

var (
	BGPPath                = ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).Bgp()
	BGPStatePath           = BGPPath.State()
	RoutingPolicyPath      = ocpath.Root().RoutingPolicy()
	RoutingPolicyStatePath = ocpath.Root().RoutingPolicy().State()
)

// NewGoBGPTaskDecl creates a new GoBGP task using the declarative configuration style.
func NewGoBGPTaskDecl(zapiURL string, listenPort uint16) *reconciler.BuiltReconciler {
	gobgpTask := newBgpDeclTask(zapiURL, listenPort)
	return reconciler.NewBuilder("gobgp-decl").WithStart(gobgpTask.startGoBGPFuncDecl).WithStop(gobgpTask.stop).Build()
}

// bgpDeclTask can be used to create a reconciler-compatible BGP task.
type bgpDeclTask struct {
	zapiURL       string
	bgpServer     *server.BgpServer
	currentConfig *bgpconfig.BgpConfigSet
	listenPort    uint16

	bgpStarted bool

	yclient *ygnmi.Client

	appliedStateMu       sync.Mutex
	appliedState         *oc.Root
	appliedBGP           *oc.NetworkInstance_Protocol_Bgp
	appliedRoutingPolicy *oc.RoutingPolicy
}

// newBgpDeclTask creates a new bgpDeclTask.
func newBgpDeclTask(zapiURL string, listenPort uint16) *bgpDeclTask {
	appliedState := &oc.Root{}
	// appliedBGP is the SoT for BGP applied configuration. It is maintained locally by the task.
	appliedBGP := appliedState.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).GetOrCreateBgp()
	appliedBGP.PopulateDefaults()
	appliedRoutingPolicy := appliedState.GetOrCreateRoutingPolicy()
	appliedRoutingPolicy.PopulateDefaults()

	return &bgpDeclTask{
		zapiURL:    zapiURL,
		bgpServer:  server.NewBgpServer(),
		listenPort: listenPort,

		appliedState:         appliedState,
		appliedBGP:           appliedBGP,
		appliedRoutingPolicy: appliedRoutingPolicy,
	}
}

func updateAppliedStateHelper[T ygot.GoStruct](yclient *ygnmi.Client, path ygnmi.SingletonQuery[T], appliedState T) {
	if _, err := gnmiclient.Replace(context.Background(), yclient, path, appliedState); err != nil {
		log.Errorf("BGP failed to update state at path %v: %v", path, err)
	}
}

// updateAppliedState is the ONLY function that's called when updating the appliedState.
//
// The input function is expected to make modifications to the applied state,
// which then this function will use to update the central cache.
func (t *bgpDeclTask) updateAppliedState(f func() error) error {
	log.V(1).Infof("BGP task: updating state")
	t.appliedStateMu.Lock()
	defer t.appliedStateMu.Unlock()
	if err := f(); err != nil {
		return err
	}
	updateAppliedStateHelper(t.yclient, BGPStatePath, t.appliedBGP)
	updateAppliedStateHelper(t.yclient, RoutingPolicyStatePath, t.appliedRoutingPolicy)
	return nil
}

// stop stops the GoBGP server.
func (t *bgpDeclTask) stop(context.Context) error {
	t.bgpServer.Stop()
	return nil
}

// startGoBGPFuncDecl starts a GoBGP server.
func (t *bgpDeclTask) startGoBGPFuncDecl(_ context.Context, yclient *ygnmi.Client) error {
	t.yclient = yclient

	b := &ocpath.Batch{}
	b.AddPaths(
		// Basic BGP paths for session establishment.
		BGPPath.Global().As().Config().PathStruct(),
		BGPPath.Global().RouterId().Config().PathStruct(),
		BGPPath.NeighborAny().PeerAs().Config().PathStruct(),
		BGPPath.NeighborAny().NeighborAddress().Config().PathStruct(),
		BGPPath.NeighborAny().NeighborPort().Config().PathStruct(),
		// BGP Policy paths
		RoutingPolicyPath.DefinedSets().PrefixSetAny().PrefixAny().IpPrefix().Config().PathStruct(),
		RoutingPolicyPath.DefinedSets().PrefixSetAny().PrefixAny().MasklengthRange().Config().PathStruct(),
		RoutingPolicyPath.PolicyDefinitionAny().StatementAny().Conditions().MatchPrefixSet().PrefixSet().Config().PathStruct(),
		RoutingPolicyPath.PolicyDefinitionAny().StatementAny().Conditions().MatchPrefixSet().MatchSetOptions().Config().PathStruct(),
		RoutingPolicyPath.PolicyDefinitionAny().StatementAny().Actions().PolicyResult().Config().PathStruct(),
		BGPPath.NeighborAny().ApplyPolicy().DefaultImportPolicy().Config().PathStruct(),
		BGPPath.NeighborAny().ApplyPolicy().DefaultExportPolicy().Config().PathStruct(),
		BGPPath.NeighborAny().ApplyPolicy().ImportPolicy().Config().PathStruct(),
		BGPPath.NeighborAny().ApplyPolicy().ExportPolicy().Config().PathStruct(),
	)

	if log.V(2) {
		if err := t.bgpServer.SetLogLevel(context.Background(), &api.SetLogLevelRequest{
			Level: api.SetLogLevelRequest_DEBUG,
		}); err != nil {
			log.Errorf("Error setting GoBGP log level: %v", err)
		}
	}
	go t.bgpServer.Serve()

	// monitor the change of the peer state
	if err := t.bgpServer.WatchEvent(context.Background(), &api.WatchEventRequest{Peer: &api.WatchEventRequest_Peer{}}, func(r *api.WatchEventResponse) {
		if p := r.GetPeer(); p != nil && p.Type == api.WatchEventResponse_PeerEvent_STATE {
			log.V(1).Info("Got peer event update:", p)
			ps := p.GetPeer().State

			t.updateAppliedState(func() error {
				neigh := t.appliedBGP.GetOrCreateNeighbor(ps.NeighborAddress)

				found := false
				if ps.SessionState.String() == "UNKNOWN" {
					neigh.SessionState = oc.Bgp_Neighbor_SessionState_UNSET
					found = true
				} else {
					for enumCode, v := range neigh.SessionState.ΛMap()[reflect.TypeOf(neigh.SessionState).Name()] {
						if v.Name == ps.SessionState.String() {
							newSessionState := oc.E_Bgp_Neighbor_SessionState(enumCode)
							if neigh.SessionState != newSessionState {
								log.V(1).Infof("Peer %s transitioned to session state %s", ps.NeighborAddress, v.Name)
								neigh.SessionState = newSessionState
							}
							found = true
							break
						}
					}
				}
				if !found {
					log.Warningf("Unknown neighbor session-state value received: %v", ps.SessionState)
				}
				return nil
			})
		}
	}); err != nil {
		return fmt.Errorf("goBgpTask failed to initialize due to error: %v", err)
	}

	// Initialize values required for reconile to be called.
	t.currentConfig = &bgpconfig.BgpConfigSet{}

	// Monitor changes to BGP intended config and apply them.
	bgpWatcher := ygnmi.Watch(
		context.Background(),
		yclient,
		b.Config(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}

			t.updateAppliedState(func() error {
				return t.reconcile(rootVal)
			})

			return ygnmi.Continue
		},
	)

	go func() {
		if _, err := bgpWatcher.Await(); err != nil {
			log.Warningf("GoBGP Task's watcher has stopped: %v", err)
		}
	}()

	// Periodically query the BGP table and update the RIBs.
	// TODO: Break this out into its own function.
	go func() {
		tick := time.NewTicker(5 * time.Second)
		for range tick.C {
			if err := t.bgpServer.ListPath(context.Background(), &api.ListPathRequest{
				TableType: api.TableType_GLOBAL,
				Family: &api.Family{
					Afi:  api.Family_AFI_IP,
					Safi: api.Family_SAFI_UNICAST,
				},
			}, func(d *api.Destination) {
				log.V(1).Infof("GoBGP global table path: %v", d)
			}); err != nil {
				log.Errorf("GoBGP ListPath call failed (global table): %v", err)
			} else {
				log.V(1).Info("GoBGP ListPath call completed (global table)")
			}

			t.updateAppliedState(func() error {
				v4uni := t.appliedBGP.GetOrCreateRib().GetOrCreateAfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).GetOrCreateIpv4Unicast()

				// TODO: Support IPv6
				t.queryTable("", "local", api.TableType_LOCAL, func(routes []*api.Destination) {
					v4uni.LocRib = nil
					locRib := v4uni.GetOrCreateLocRib()
					for _, route := range routes {
						for j, path := range route.Paths {
							var origin oc.NetworkInstance_Protocol_Bgp_Rib_AfiSafi_Ipv4Unicast_LocRib_Route_Origin_Union
							if path.SourceId == "" {
								// TODO: For locally-originated routes figure out how to get the originating protocol.
								origin = oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_UNSET
							} else {
								origin = oc.UnionString(path.SourceId)
							}
							// TODO: this ID should match the ID in adj-rib-in-post.
							locRib.GetOrCreateRoute(route.Prefix, origin, uint32(j))
						}
					}
				})

				for neigh := range t.appliedBGP.Neighbor {
					neighContainer := v4uni.GetOrCreateNeighbor(neigh)
					neighContainer.AdjRibInPre = nil
					neighContainer.AdjRibInPost = nil
					neighContainer.AdjRibOutPre = nil
					neighContainer.AdjRibOutPost = nil
					t.queryTable(neigh, "adj-rib-in", api.TableType_ADJ_IN, func(routes []*api.Destination) {
						for _, route := range routes {
							for j, path := range route.Paths {
								neighContainer.GetOrCreateAdjRibInPre().GetOrCreateRoute(route.Prefix, uint32(j))
								if !path.Filtered {
									neighContainer.GetOrCreateAdjRibInPost().GetOrCreateRoute(route.Prefix, uint32(j))
								}
							}
						}
					})

					t.queryTable(neigh, "adj-rib-out", api.TableType_ADJ_OUT, func(routes []*api.Destination) {
						for _, route := range routes {
							for j, path := range route.Paths {
								// Per OpenConfig the ID of this should be the ID assigned when exchanging add-path routes. However
								// GoBGP doesn't seem to support the add-path capability and so just going to use the first path
								// with 0 as the ID here. GoBGP does support AddPath as a gRPC call but when advertising the routes
								// the generated UUID isn't propagated.
								//
								// Note that path.NeighborIp is <nil> for some reason so have to use neigh.
								neighContainer.GetOrCreateAdjRibOutPre().GetOrCreateRoute(route.Prefix, uint32(j))
								if !path.Filtered {
									neighContainer.GetOrCreateAdjRibOutPost().GetOrCreateRoute(route.Prefix, uint32(j))
								}
							}
						}
					})
				}
				return nil
			})
		}
	}()

	return nil
}

// queryTable queries for all routes stored in the specified table, applying f
// to the routes that are queried if the query was successful or logging an
// error otherwise.
func (t *bgpDeclTask) queryTable(neighbor, tableName string, tableType api.TableType, f func(route []*api.Destination)) {
	var routes []*api.Destination
	if err := t.bgpServer.ListPath(context.Background(), &api.ListPathRequest{
		Name:      neighbor,
		TableType: tableType,
		Family: &api.Family{
			Afi:  api.Family_AFI_IP,
			Safi: api.Family_SAFI_UNICAST,
		},
		// This is always set to true since GoBGP doesn't actually
		// filter the paths out, only mark them as filtered out by the
		// IMPORT or EXPORT policy.
		EnableFiltered: true,
	}, func(d *api.Destination) {
		routes = append(routes, d)
		log.V(0).Infof("GoBGP %s table path (neighbor if applicable: %q): %v", tableName, neighbor, d)
	}); err != nil {
		log.Errorf("GoBGP ListPath call failed (%s table): %v", tableType, err)
	} else {
		log.V(1).Info("GoBGP ListPath call completed (%s table)", tableName)
		f(routes)
	}
}

// reconcile examines the difference between the intended and applied
// configuration, and makes GoBGP API calls accordingly to update the applied
// configuration in the direction of intended configuration.
func (t *bgpDeclTask) reconcile(intended *oc.Root) error {
	intendedBGP := intended.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).GetOrCreateBgp()
	intendedPolicy := intended.GetOrCreateRoutingPolicy()
	newConfig := intendedToGoBGP(intendedBGP, intendedPolicy, t.zapiURL, t.listenPort)

	intendedGlobal := intendedBGP.GetOrCreateGlobal()
	bgpShouldStart := intendedGlobal.As != nil && intendedGlobal.RouterId != nil
	switch {
	case bgpShouldStart && !t.bgpStarted:
		log.V(1).Info("Starting BGP")
		var err error
		t.currentConfig, err = InitialConfig(context.Background(), t.appliedBGP, t.bgpServer, newConfig, gracefulRestart)
		if err != nil {
			return fmt.Errorf("Failed to apply initial BGP configuration %v", newConfig)
		}
		t.bgpStarted = true
	case !bgpShouldStart && t.bgpStarted:
		log.V(1).Info("Stopping BGP")
		if err := t.bgpServer.StopBgp(context.Background(), &api.StopBgpRequest{}); err != nil {
			return errors.New("Failed to stop BGP service")
		}
		t.bgpStarted = false
		t.currentConfig = &bgpconfig.BgpConfigSet{}
		*t.appliedBGP = oc.NetworkInstance_Protocol_Bgp{}
		t.appliedBGP.PopulateDefaults()
	case t.bgpStarted:
		log.V(1).Info("Updating BGP")
		var err error
		t.currentConfig, err = UpdateConfig(context.Background(), t.appliedBGP, t.bgpServer, t.currentConfig, newConfig)
		if err != nil {
			return fmt.Errorf("Failed to update BGP service: %v", newConfig)
		}
	default:
		// Waiting for BGP to be startable.
		return nil
	}

	return nil
}

// intendedToGoBGP translates from OC to GoBGP intended config.
//
// GoBGP's notion of config vs. state does not conform to OpenConfig (see
// https://github.com/osrg/gobgp/issues/2584)
// Therefore, we need a compatibility layer between the two configs.
func intendedToGoBGP(bgpoc *oc.NetworkInstance_Protocol_Bgp, policyoc *oc.RoutingPolicy, zapiURL string, listenPort uint16) *bgpconfig.BgpConfigSet {
	bgpConfig := &bgpconfig.BgpConfigSet{}

	// Global config
	global := bgpoc.GetOrCreateGlobal()

	bgpConfig.Global.Config.As = global.GetAs()
	bgpConfig.Global.Config.RouterId = global.GetRouterId()
	bgpConfig.Global.Config.Port = int32(listenPort)

	localAddress := ""
	if localAddr, err := netip.ParseAddr(global.GetRouterId()); err == nil && localAddr.IsLoopback() {
		localAddress = localAddr.String()
	}

	bgpConfig.Zebra.Config = bgpconfig.ZebraConfig{
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

	// Neighbours & apply policy
	bgpConfig.Neighbors = []bgpconfig.Neighbor{}
	neighborSets := map[string]bgpconfig.NeighborSet{}
	var neighborDefaultPolicies []bgpconfig.PolicyDefinition
	for neighAddr, neigh := range bgpoc.Neighbor {
		// Ideally A simple conversion of apply-policy is sufficient, but due to GoBGP using
		// a global set of apply-policy instead of per-neighbour policies, we need to create
		// neighbour sets and modify input policy statements so that we retain the same
		// per-neighbour behaviour while only using a single set of global policies.
		//
		// To do this, we first make the assumption that each policy is only be used for a
		// single neighbour across import and export policies (TODO: this is a simplifying
		// assumption that will be changed). Then, we add a neighbor-set for each policy based on the
		// relevant neighbour, and lastly simply concatenate the set of all apply-policy for
		// each neighbour into the global apply-policy.
		applyPolicy := bgpconfig.ApplyPolicy{
			Config: bgpconfig.ApplyPolicyConfig{
				// TODO(wenbli): Update applied state.
				DefaultImportPolicy: convertDefaultPolicy(neigh.GetApplyPolicy().GetDefaultImportPolicy()),
				DefaultExportPolicy: convertDefaultPolicy(neigh.GetApplyPolicy().GetDefaultExportPolicy()),
				ImportPolicyList:    neigh.GetApplyPolicy().GetImportPolicy(),
				ExportPolicyList:    neigh.GetApplyPolicy().GetExportPolicy(),
			},
		}

		neighborsetName := neighAddr + "default"
		neighborSets[neighborsetName] = bgpconfig.NeighborSet{
			NeighborSetName:  neighborsetName,
			NeighborInfoList: []string{neighAddr},
		}
		defaultImportPolicyName := neighborsetName + "import"
		defaultExportPolicyName := neighborsetName + "export"
		neighborDefaultPolicies = append(neighborDefaultPolicies, bgpconfig.PolicyDefinition{
			Name: defaultImportPolicyName,
			Statements: []bgpconfig.Statement{{
				// Use a customized name for the default policies.
				Name: defaultImportPolicyName,
				Conditions: bgpconfig.Conditions{
					MatchNeighborSet: bgpconfig.MatchNeighborSet{
						NeighborSet: neighborsetName,
					},
				},
				Actions: bgpconfig.Actions{
					RouteDisposition: defaultPolicyToRouteDisp(applyPolicy.Config.DefaultImportPolicy),
				},
			}},
		}, bgpconfig.PolicyDefinition{
			Name: defaultExportPolicyName,
			Statements: []bgpconfig.Statement{{
				// Use a customized name for the default policies.
				Name: defaultExportPolicyName,
				Conditions: bgpconfig.Conditions{
					MatchNeighborSet: bgpconfig.MatchNeighborSet{
						NeighborSet: neighborsetName,
					},
				},
				Actions: bgpconfig.Actions{
					RouteDisposition: defaultPolicyToRouteDisp(applyPolicy.Config.DefaultExportPolicy),
				},
			}},
		})
		for _, importPolicyName := range applyPolicy.Config.ImportPolicyList {
			neighborSet := neighborSets[importPolicyName]
			neighborSet.NeighborSetName = importPolicyName
			neighborSet.NeighborInfoList = append(neighborSet.NeighborInfoList, neighAddr)
			neighborSets[importPolicyName] = neighborSet
		}
		for _, exportPolicyName := range applyPolicy.Config.ExportPolicyList {
			neighborSet := neighborSets[exportPolicyName]
			neighborSet.NeighborSetName = exportPolicyName
			neighborSet.NeighborInfoList = append(neighborSet.NeighborInfoList, neighAddr)
			neighborSets[exportPolicyName] = neighborSet
		}
		applyPolicy.Config.ImportPolicyList = append(applyPolicy.Config.ImportPolicyList, defaultImportPolicyName)
		applyPolicy.Config.ExportPolicyList = append(applyPolicy.Config.ExportPolicyList, defaultExportPolicyName)

		// GoBGP in non-route server mode doesn't have per-neighbor
		// policies -- instead it uses a global policy for everything.
		// So, concatenate policies from all
		// neighbours to create one giant apply-policy list.
		bgpConfig.Global.ApplyPolicy.Config.ImportPolicyList = append(bgpConfig.Global.ApplyPolicy.Config.ImportPolicyList, applyPolicy.Config.ImportPolicyList...)
		bgpConfig.Global.ApplyPolicy.Config.ExportPolicyList = append(bgpConfig.Global.ApplyPolicy.Config.ExportPolicyList, applyPolicy.Config.ExportPolicyList...)

		// Add neighbour config.
		bgpConfig.Neighbors = append(bgpConfig.Neighbors, bgpconfig.Neighbor{
			Config: bgpconfig.NeighborConfig{
				PeerAs:          neigh.GetPeerAs(),
				NeighborAddress: neighAddr,
			},
			// This is needed because GoBGP's configuration diffing
			// logic may check the state value instead of the
			// config value.
			State: bgpconfig.NeighborState{
				PeerAs:          neigh.GetPeerAs(),
				NeighborAddress: neighAddr,
			},
			Transport: bgpconfig.Transport{
				Config: bgpconfig.TransportConfig{
					LocalAddress: localAddress,
					RemotePort:   neigh.GetNeighborPort(),
				},
			},
			// NOTE: This line only matters when the neighbour is a
			// route server client, since otherwise GoBGP will
			// simply use the global RIB and global policies.
			ApplyPolicy: applyPolicy,
		})
	}

	for prefixSetName, prefixSet := range policyoc.GetOrCreateDefinedSets().PrefixSet {
		var prefixList []bgpconfig.Prefix
		for _, prefix := range prefixSet.Prefix {
			r := prefix.GetMasklengthRange()
			if r == "exact" {
				// GoBGP recognizes "" instead of "exact"
				r = ""
			}
			prefixList = append(prefixList, bgpconfig.Prefix{
				IpPrefix:        prefix.GetIpPrefix(),
				MasklengthRange: r,
			})
		}

		bgpConfig.DefinedSets.PrefixSets = append(bgpConfig.DefinedSets.PrefixSets, bgpconfig.PrefixSet{
			PrefixSetName: prefixSetName,
			PrefixList:    prefixList,
		})
	}

	// Policy definitions
	//
	// TODO(wenbli): Duplicate policy definitions so that we allow the same
	// policy to be used by different apply-policy references.
	for _, policy := range policyoc.PolicyDefinition {
		var statements []bgpconfig.Statement
		for _, statement := range policy.Statement {
			statements = append(statements, bgpconfig.Statement{
				Name: statement.GetName(),
				Conditions: bgpconfig.Conditions{
					MatchPrefixSet: bgpconfig.MatchPrefixSet{
						PrefixSet:       statement.GetConditions().GetMatchPrefixSet().GetPrefixSet(),
						MatchSetOptions: convertMatchSetOptionsRestrictedType(statement.GetConditions().GetMatchPrefixSet().GetMatchSetOptions()),
					},
					MatchNeighborSet: bgpconfig.MatchNeighborSet{
						// Name the neighbor set as the policy so that the policy only applies to referring neighbours.
						NeighborSet: policy.GetName(),
					},
				},
				Actions: bgpconfig.Actions{
					RouteDisposition: convertRouteDisposition(statement.GetActions().GetPolicyResult()),
				},
			})
		}

		bgpConfig.PolicyDefinitions = append(bgpConfig.PolicyDefinitions, bgpconfig.PolicyDefinition{
			Name:       policy.GetName(),
			Statements: statements,
		})

		if _, ok := neighborSets[policy.GetName()]; !ok {
			neighborSets[policy.GetName()] = bgpconfig.NeighborSet{NeighborSetName: policy.GetName()}
		}
	}
	bgpConfig.PolicyDefinitions = append(bgpConfig.PolicyDefinitions, neighborDefaultPolicies...)

	// Convert defined sets.
	for _, neighborSet := range neighborSets {
		bgpConfig.DefinedSets.NeighborSets = append(bgpConfig.DefinedSets.NeighborSets, neighborSet)
	}

	return bgpConfig
}
