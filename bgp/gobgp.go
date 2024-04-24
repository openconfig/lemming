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
	"fmt"
	"math"
	"net/netip"
	"reflect"
	"strconv"
	"strings"
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
	api "github.com/osrg/gobgp/v3/api"
	"github.com/osrg/gobgp/v3/pkg/config"
	gobgpoc "github.com/osrg/gobgp/v3/pkg/config/oc"
	"github.com/osrg/gobgp/v3/pkg/server"
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

// NewGoBGPTask creates a new GoBGP task implementing OpenConfig BGP functionalities.
func NewGoBGPTask(targetName, zapiURL string, listenPort uint16) *reconciler.BuiltReconciler {
	gobgpTask := newBgpTask(targetName, zapiURL, listenPort)
	return reconciler.NewBuilder("gobgp").WithStart(gobgpTask.start).WithStop(gobgpTask.stop).WithValidator(
		[]ygnmi.PathStruct{
			RoutingPolicyPath.DefinedSets().PrefixSetAny().Mode().Config().PathStruct(),
		}, validatePrefixSetMode).Build()
}

// validatePrefixSetMode check that all prefix sets have the correct mode.
func validatePrefixSetMode(root *oc.Root) error {
	definedSets := root.GetRoutingPolicy().GetDefinedSets()
	if definedSets == nil {
		return nil
	}
	for _, prefixSet := range definedSets.PrefixSet {
		if len(prefixSet.Prefix) == 0 {
			continue
		}
		if prefixSet.GetMode() == oc.PrefixSet_Mode_MIXED {
			// This is always valid.
			continue
		}
		var gotMode oc.E_PrefixSet_Mode
		for _, pfx := range prefixSet.Prefix {
			p, err := netip.ParsePrefix(pfx.GetIpPrefix())
			if err != nil {
				return fmt.Errorf("invalid prefix %q in prefix set %q", pfx.GetIpPrefix(), prefixSet.GetName())
			}
			switch gotMode {
			case oc.PrefixSet_Mode_UNSET:
				if p.Addr().Is4() {
					gotMode = oc.PrefixSet_Mode_IPV4
				} else if p.Addr().Is6() {
					gotMode = oc.PrefixSet_Mode_IPV6
				}
			case oc.PrefixSet_Mode_IPV4:
				if p.Addr().Is6() {
					gotMode = oc.PrefixSet_Mode_MIXED
				}
			case oc.PrefixSet_Mode_IPV6:
				if p.Addr().Is4() {
					gotMode = oc.PrefixSet_Mode_MIXED
				}
			}
		}

		if wantMode := prefixSet.GetMode(); gotMode != wantMode {
			return fmt.Errorf("prefix set %q has mode %s based on parsing given prefixes, but has configured mode %s", prefixSet.GetName(), gotMode, wantMode)
		}
	}
	return nil
}

// bgpTask can be used to create a reconciler-compatible BGP task.
type bgpTask struct {
	targetName    string
	zapiURL       string
	bgpServer     *server.BgpServer
	currentConfig *gobgpoc.BgpConfigSet
	listenPort    uint16

	bgpStarted bool

	yclient *ygnmi.Client

	commAttrTracker *ocRIBAttrIndicesTracker[string]
	attrSetTracker  *ocRIBAttrIndicesTracker[ribAttrSet]

	appliedStateMu       sync.Mutex
	appliedState         *oc.Root
	appliedBGP           *oc.NetworkInstance_Protocol_Bgp
	appliedRoutingPolicy *oc.RoutingPolicy
}

// newBgpTask creates a new bgpTask.
func newBgpTask(targetName, zapiURL string, listenPort uint16) *bgpTask {
	appliedState := &oc.Root{}
	// appliedBGP is the SoT for BGP applied configuration. It is maintained locally by the task.
	appliedBGP := appliedState.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).GetOrCreateBgp()
	appliedRoutingPolicy := appliedState.GetOrCreateRoutingPolicy()

	return &bgpTask{
		targetName: targetName,
		zapiURL:    zapiURL,
		listenPort: listenPort,

		commAttrTracker: newOCRIBAttrIndices[string](),
		attrSetTracker:  newOCRIBAttrIndices[ribAttrSet](),

		appliedState:         appliedState,
		appliedBGP:           appliedBGP,
		appliedRoutingPolicy: appliedRoutingPolicy,
	}
}

// stop stops the GoBGP server.
func (t *bgpTask) stop(context.Context) error {
	t.bgpServer.Stop()
	return nil
}

// start starts a GoBGP server.
func (t *bgpTask) start(ctx context.Context, yclient *ygnmi.Client) error {
	t.yclient = yclient

	b := &ocpath.Batch{}
	b.AddPaths(
		// Basic BGP paths for session establishment.
		BGPPath.Global().As().Config().PathStruct(),
		BGPPath.Global().RouterId().Config().PathStruct(),
		BGPPath.NeighborAny().PeerAs().Config().PathStruct(),
		BGPPath.NeighborAny().NeighborAddress().Config().PathStruct(),
		BGPPath.NeighborAny().NeighborPort().Config().PathStruct(),
		BGPPath.NeighborAny().Transport().LocalAddress().Config().PathStruct(),
		// BGP Policy statements
		RoutingPolicyPath.PolicyDefinitionAny().Name().Config().PathStruct(),
		RoutingPolicyPath.PolicyDefinitionAny().StatementMap().Config().PathStruct(),
		BGPPath.NeighborAny().ApplyPolicy().DefaultImportPolicy().Config().PathStruct(),
		BGPPath.NeighborAny().ApplyPolicy().DefaultExportPolicy().Config().PathStruct(),
		BGPPath.NeighborAny().ApplyPolicy().ImportPolicy().Config().PathStruct(),
		BGPPath.NeighborAny().ApplyPolicy().ExportPolicy().Config().PathStruct(),
		// BGP defined sets
		// -- prefix sets
		RoutingPolicyPath.DefinedSets().PrefixSetAny().Name().Config().PathStruct(),
		RoutingPolicyPath.DefinedSets().PrefixSetAny().Mode().Config().PathStruct(),
		RoutingPolicyPath.DefinedSets().PrefixSetAny().PrefixAny().IpPrefix().Config().PathStruct(),
		RoutingPolicyPath.DefinedSets().PrefixSetAny().PrefixAny().MasklengthRange().Config().PathStruct(),
		// -- community sets
		ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySetAny().CommunitySetName().Config().PathStruct(),
		ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySetAny().CommunityMember().Config().PathStruct(),
		ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySetAny().MatchSetOptions().Config().PathStruct(),
		// -- AS path sets
		ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().AsPathSetAny().AsPathSetName().Config().PathStruct(),
		ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().AsPathSetAny().AsPathSetMember().Config().PathStruct(),
	)

	if err := t.createNewGoBGPServer(ctx); err != nil {
		log.Errorf("Failed to start GoBGP server: %v", err)
	}

	// Initialize values required for reconile to be called.
	t.currentConfig = &gobgpoc.BgpConfigSet{}

	// Monitor changes to BGP intended config and apply them.
	bgpWatcher := ygnmi.Watch(
		ctx,
		yclient,
		b.Config(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}

			t.updateAppliedState(ctx, func() error {
				return t.reconcile(ctx, rootVal)
			})

			return ygnmi.Continue
		},
	)

	go func() {
		if _, err := bgpWatcher.Await(); err != nil {
			log.Warningf("GoBGP Task's watcher has stopped: %v", err)
		}
	}()

	return nil
}

// reconcile examines the difference between the intended and applied
// configuration, and makes GoBGP API calls accordingly to update the applied
// configuration in the direction of intended configuration.
func (t *bgpTask) reconcile(ctx context.Context, intended *oc.Root) error {
	intendedBGP := intended.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).GetOrCreateBgp()
	intendedPolicy := intended.GetOrCreateRoutingPolicy()
	newConfig := intendedToGoBGP(intendedBGP, intendedPolicy, t.zapiURL, t.listenPort)

	intendedGlobal := intendedBGP.GetOrCreateGlobal()
	bgpShouldStart := intendedGlobal.As != nil && intendedGlobal.RouterId != nil
	switch {
	case bgpShouldStart && !t.bgpStarted:
		log.V(1).Info("Starting BGP")
		var err error
		t.currentConfig, err = config.InitialConfig(ctx, t.bgpServer, newConfig, gracefulRestart)
		if err != nil {
			return fmt.Errorf("Failed to apply initial BGP configuration %v", newConfig)
		}
		t.bgpStarted = true
	case t.bgpStarted:
		log.V(1).Info("Updating BGP")
		var err error
		t.currentConfig, err = config.UpdateConfig(ctx, t.bgpServer, t.currentConfig, newConfig)
		if err != nil {
			return fmt.Errorf("Failed to update BGP service: %v", newConfig)
		}
	default:
		// Waiting for BGP to be startable.
		return nil
	}

	err := ygot.MergeStructInto(t.appliedBGP, intendedBGP, &ygot.MergeOverwriteExistingFields{})
	// TODO(wenbli): Since policy definitions is an atomic node,
	// unsupported policy leaves will be merged as well. Therefore omitting
	// them from the applied state until we find a way to to prune out
	// unsupported paths prior to merge.
	t.appliedRoutingPolicy.PolicyDefinition = nil
	return err
}

// updateAppliedState is the ONLY function that's called when updating the appliedState.
//
// The input function is expected to make modifications to the applied state,
// which then this function will use to update the central cache.
func (t *bgpTask) updateAppliedState(ctx context.Context, f func() error) error {
	log.V(1).Infof("BGP task: updating state")
	t.appliedStateMu.Lock()
	defer t.appliedStateMu.Unlock()
	if err := f(); err != nil {
		return err
	}
	updateAppliedStateHelper(ctx, t.yclient, BGPStatePath, t.appliedBGP)
	updateAppliedStateHelper(ctx, t.yclient, RoutingPolicyStatePath, t.appliedRoutingPolicy)
	return nil
}

func updateAppliedStateHelper[T ygot.GoStruct](ctx context.Context, yclient *ygnmi.Client, path ygnmi.SingletonQuery[T], appliedState T) {
	if _, err := gnmiclient.Replace(ctx, yclient, path, appliedState); err != nil {
		log.Errorf("BGP failed to update state at path %v: %v", path, err)
	}
}

// createNewGoBGPServer creates and starts a new GoBGP Server.
func (t *bgpTask) createNewGoBGPServer(ctx context.Context) error {
	t.bgpServer = server.NewBgpServer()

	if log.V(2) {
		if err := t.bgpServer.SetLogLevel(ctx, &api.SetLogLevelRequest{
			Level: api.SetLogLevelRequest_DEBUG,
		}); err != nil {
			log.Errorf("Error setting GoBGP log level: %v", err)
		}
	}
	go t.bgpServer.Serve()

	// monitor the change of the peer state
	if err := t.bgpServer.WatchEvent(ctx, &api.WatchEventRequest{Peer: &api.WatchEventRequest_Peer{}}, func(r *api.WatchEventResponse) {
		if p := r.GetPeer(); p != nil && p.Type == api.WatchEventResponse_PeerEvent_STATE {
			log.V(1).Info("Got peer event update:", p)
			ps := p.GetPeer().State

			t.updateAppliedState(ctx, func() error {
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

	go func() {
		tick := time.NewTicker(5 * time.Second)
		for range tick.C {
			if err := t.updateRIBs(ctx); err != nil {
				log.Warning("Error while updating BGP RIB data: %v", err)
			}
		}
	}()

	return nil
}

// updateRIBs updates the BGP RIBs.
func (t *bgpTask) updateRIBs(ctx context.Context) error {
	// Log global tables
	t.queryTable(ctx, "", false, api.TableType_GLOBAL, api.Family_AFI_IP, nil)
	t.queryTable(ctx, "", false, api.TableType_GLOBAL, api.Family_AFI_IP6, nil)

	return t.updateAppliedState(ctx, func() error {
		t.beginAttrPopulation()
		defer t.completeAttrPopulation()

		bgpRIB := t.appliedBGP.GetOrCreateRib()
		v4uni := t.appliedBGP.GetOrCreateRib().GetOrCreateAfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).GetOrCreateIpv4Unicast()
		v6uni := t.appliedBGP.GetOrCreateRib().GetOrCreateAfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV6_UNICAST).GetOrCreateIpv6Unicast()

		t.queryTable(ctx, "", false, api.TableType_LOCAL, api.Family_AFI_IP, func(routes []*api.Destination) {
			v4uni.LocRib = nil
			locRib := v4uni.GetOrCreateLocRib()
			for _, route := range routes {
				for i, path := range route.Paths {
					var origin oc.NetworkInstance_Protocol_Bgp_Rib_AfiSafi_Ipv4Unicast_LocRib_Route_Origin_Union
					if path.SourceId == "" {
						// TODO: For locally-originated routes figure out how to get the originating protocol.
						origin = oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_UNSET
					} else {
						origin = oc.UnionString(path.NeighborIp)
					}
					// TODO: this ID should match the ID in adj-rib-in-post.
					t.populateRIBAttrs(path, bgpRIB, locRib.GetOrCreateRoute(route.Prefix, origin, uint32(i)))
				}
			}
		})
		t.queryTable(ctx, "", false, api.TableType_LOCAL, api.Family_AFI_IP6, func(routes []*api.Destination) {
			v6uni.LocRib = nil
			locRib := v6uni.GetOrCreateLocRib()
			for _, route := range routes {
				for i, path := range route.Paths {
					var origin oc.NetworkInstance_Protocol_Bgp_Rib_AfiSafi_Ipv6Unicast_LocRib_Route_Origin_Union
					if path.SourceId == "" {
						// TODO: For locally-originated routes figure out how to get the originating protocol.
						origin = oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_UNSET
					} else {
						origin = oc.UnionString(path.NeighborIp)
					}
					// TODO: this ID should match the ID in adj-rib-in-post.
					t.populateRIBAttrs(path, bgpRIB, locRib.GetOrCreateRoute(route.Prefix, origin, uint32(i)))
				}
			}
		})

		for neigh := range t.appliedBGP.Neighbor {
			neighContainer := v4uni.GetOrCreateNeighbor(neigh)
			neighContainer.AdjRibInPre = nil
			neighContainer.AdjRibInPost = nil
			neighContainer.AdjRibOutPre = nil
			neighContainer.AdjRibOutPost = nil
			t.queryTable(ctx, neigh, false, api.TableType_ADJ_IN, api.Family_AFI_IP, func(routes []*api.Destination) {
				for _, route := range routes {
					for i, path := range route.Paths {
						// TODO: this ID should be retrieved from the update message.
						t.populateRIBAttrs(path, bgpRIB, neighContainer.GetOrCreateAdjRibInPre().GetOrCreateRoute(route.Prefix, uint32(i)))
					}
				}
			})
			t.queryTable(ctx, neigh, true, api.TableType_ADJ_IN, api.Family_AFI_IP, func(routes []*api.Destination) {
				for _, route := range routes {
					for i, path := range route.Paths {
						// TODO: this ID should be retrieved from the update message.
						if !path.Filtered {
							t.populateRIBAttrs(path, bgpRIB, neighContainer.GetOrCreateAdjRibInPost().GetOrCreateRoute(route.Prefix, uint32(i)))
						}
					}
				}
			})
			t.queryTable(ctx, neigh, false, api.TableType_ADJ_OUT, api.Family_AFI_IP, func(routes []*api.Destination) {
				for _, route := range routes {
					for i, path := range route.Paths {
						// Per OpenConfig the ID of this should be the ID assigned when exchanging add-path routes. However
						// GoBGP doesn't seem to support the add-path capability and so just going to use the first path
						// with 0 as the ID here. GoBGP does support AddPath as a gRPC call but when advertising the routes
						// the generated UUID isn't propagated.
						//
						// Note that path.NeighborIp is <nil> for some reason so have to use neigh.
						t.populateRIBAttrs(path, bgpRIB, neighContainer.GetOrCreateAdjRibOutPre().GetOrCreateRoute(route.Prefix, uint32(i)))
					}
				}
			})
			t.queryTable(ctx, neigh, true, api.TableType_ADJ_OUT, api.Family_AFI_IP, func(routes []*api.Destination) {
				for _, route := range routes {
					for i, path := range route.Paths {
						// Per OpenConfig the ID of this should be the ID assigned when exchanging add-path routes. However
						// GoBGP doesn't seem to support the add-path capability and so just going to use the first path
						// with 0 as the ID here. GoBGP does support AddPath as a gRPC call but when advertising the routes
						// the generated UUID isn't propagated.
						//
						// Note that path.NeighborIp is <nil> for some reason so have to use neigh.
						if !path.Filtered {
							t.populateRIBAttrs(path, bgpRIB, neighContainer.GetOrCreateAdjRibOutPost().GetOrCreateRoute(route.Prefix, uint32(i)))
						}
					}
				}
			})
		}

		for neigh := range t.appliedBGP.Neighbor {
			neighContainer := v6uni.GetOrCreateNeighbor(neigh)
			neighContainer.AdjRibInPre = nil
			neighContainer.AdjRibInPost = nil
			neighContainer.AdjRibOutPre = nil
			neighContainer.AdjRibOutPost = nil
			t.queryTable(ctx, neigh, false, api.TableType_ADJ_IN, api.Family_AFI_IP6, func(routes []*api.Destination) {
				for _, route := range routes {
					for i, path := range route.Paths {
						// TODO: this ID should be retrieved from the update message.
						t.populateRIBAttrs(path, bgpRIB, neighContainer.GetOrCreateAdjRibInPre().GetOrCreateRoute(route.Prefix, uint32(i)))
					}
				}
			})
			t.queryTable(ctx, neigh, true, api.TableType_ADJ_IN, api.Family_AFI_IP6, func(routes []*api.Destination) {
				for _, route := range routes {
					for i, path := range route.Paths {
						// TODO: this ID should be retrieved from the update message.
						if !path.Filtered {
							t.populateRIBAttrs(path, bgpRIB, neighContainer.GetOrCreateAdjRibInPost().GetOrCreateRoute(route.Prefix, uint32(i)))
						}
					}
				}
			})
			t.queryTable(ctx, neigh, false, api.TableType_ADJ_OUT, api.Family_AFI_IP6, func(routes []*api.Destination) {
				for _, route := range routes {
					for i, path := range route.Paths {
						// Per OpenConfig the ID of this should be the ID assigned when exchanging add-path routes. However
						// GoBGP doesn't seem to support the add-path capability and so just going to use the first path
						// with 0 as the ID here. GoBGP does support AddPath as a gRPC call but when advertising the routes
						// the generated UUID isn't propagated.
						//
						// Note that path.NeighborIp is <nil> for some reason so have to use neigh.
						t.populateRIBAttrs(path, bgpRIB, neighContainer.GetOrCreateAdjRibOutPre().GetOrCreateRoute(route.Prefix, uint32(i)))
					}
				}
			})
			t.queryTable(ctx, neigh, true, api.TableType_ADJ_OUT, api.Family_AFI_IP6, func(routes []*api.Destination) {
				for _, route := range routes {
					for i, path := range route.Paths {
						// Per OpenConfig the ID of this should be the ID assigned when exchanging add-path routes. However
						// GoBGP doesn't seem to support the add-path capability and so just going to use the first path
						// with 0 as the ID here. GoBGP does support AddPath as a gRPC call but when advertising the routes
						// the generated UUID isn't propagated.
						//
						// Note that path.NeighborIp is <nil> for some reason so have to use neigh.
						if !path.Filtered {
							t.populateRIBAttrs(path, bgpRIB, neighContainer.GetOrCreateAdjRibOutPost().GetOrCreateRoute(route.Prefix, uint32(i)))
						}
					}
				}
			})
		}
		return nil
	})
}

// queryTable queries for all routes stored in the specified table, applying f
// to the routes that are queried if the query was successful or logging an
// error otherwise.
func (t *bgpTask) queryTable(ctx context.Context, neighbor string, postPolicy bool, tableType api.TableType, afi api.Family_Afi, f func(route []*api.Destination)) {
	tableName := fmt.Sprint(tableType)
	var enableFiltered bool
	switch tableType {
	case api.TableType_ADJ_IN:
		enableFiltered = postPolicy
		if postPolicy {
			tableName += "-post"
		} else {
			tableName += "-pre"
		}
	case api.TableType_ADJ_OUT:
		// NOTE: This doesn't intuitively make sense since by meaning, filtered == postPolicy.
		// However, to avoid a breaking change this is used.
		// For background see https://github.com/osrg/gobgp/issues/2765
		enableFiltered = !postPolicy
		if postPolicy {
			tableName += "-post"
		} else {
			tableName += "-pre"
		}
	}
	tableName += fmt.Sprintf("-%v", afi)

	var routes []*api.Destination
	if err := t.bgpServer.ListPath(ctx, &api.ListPathRequest{
		Name:      neighbor,
		TableType: tableType,
		Family: &api.Family{
			Afi:  afi,
			Safi: api.Family_SAFI_UNICAST,
		},
		EnableFiltered: enableFiltered,
	}, func(d *api.Destination) {
		if f != nil {
			routes = append(routes, d)
		}
		log.V(0).Infof("%s: GoBGP %s table path (neighbor if applicable: %q): %v", t.targetName, tableName, neighbor, d)
	}); err != nil {
		if err.Error() != "bgp server hasn't started yet" {
			log.Errorf("GoBGP ListPath call failed (%s, %s, %s table): %v", tableName, tableType, afi, err)
		}
	} else {
		log.V(1).Infof("GoBGP ListPath call completed (%s, %s, %s table)", tableName, tableType, afi)
		if f != nil {
			f(routes)
		}
	}
}

// beginAttrPopulation clears the RIB's attributes and starts to assign or
// re-assign attribute indices.
//
// This MUST be called prior to calling .populateRIBAttrs() over all routes
// in the RIB.
func (t *bgpTask) beginAttrPopulation() {
	// Clear RIB attributes for fresh population.
	t.appliedBGP.GetOrCreateRib().Community = nil
	t.appliedBGP.GetOrCreateRib().AttrSet = nil
	t.commAttrTracker.beginAllocation()
	t.attrSetTracker.beginAllocation()
}

// completeAttrPopulation garbage collects indices and cleans up the state for
// th next round of attribute population.
//
// This MUST be called after calling .populateRIBAttrs() over all routes in
// the RIB.
func (t *bgpTask) completeAttrPopulation() {
	// Clear RIB attributes for fresh population.
	t.commAttrTracker.completeAllocation()
	t.attrSetTracker.completeAllocation()
}

// populateRIBAttrs populates path attributes of routes in the RIB.
//
// - path is the GoBGP path attributes.
// - rib is the BGP RIB applied state to be populated.
// - route is the route in the RIB whose attribute reference needs to be updated.
func (t *bgpTask) populateRIBAttrs(path *api.Path, rib *oc.NetworkInstance_Protocol_Bgp_Rib, route ocRIBRoute) {
	commsToString := func(comms []uint32) string {
		var b strings.Builder
		for i, comm := range comms {
			if i != 0 {
				b.WriteRune(' ')
			}
			b.WriteString(strconv.FormatUint(uint64(comm), 10))
		}
		return b.String()
	}

	asSegmentsToString := func(segs []*api.AsSegment) string {
		var b strings.Builder
		for i, s := range segs {
			if i != 0 {
				b.WriteRune(' ')
			}
			// Here we trust that the proto String() implementation
			// is a function (in the math sense), and will not map
			// two different segments to the same string
			// representation.
			b.WriteString(s.String())
			b.WriteRune('\n')
		}
		return b.String()
	}

	var (
		hasCommunity       bool
		commIndex          uint64
		hasOrigin          bool
		hasMED             bool
		hasLocalPref       bool
		hasASPathAttribute bool
		asSegments         []*api.AsSegment
		attrSet            ribAttrSet
	)

	for _, attr := range path.GetPattrs() {
		m, err := attr.UnmarshalNew()
		if err != nil {
			log.Errorf("BGP: Unable to unmarshal a GoBGP path attribute")
			continue
		}
		switch m := m.(type) {
		case *api.CommunitiesAttribute:
			if comms := m.GetCommunities(); len(comms) > 0 {
				hasCommunity = true
				commIndex = t.commAttrTracker.getOrAllocIndex(commsToString(comms))
				rib.GetOrCreateCommunity(commIndex).SetCommunity(communitiesToOC(comms))
			}
		case *api.OriginAttribute:
			hasOrigin = true
			switch origin := m.GetOrigin(); origin {
			case 0:
				attrSet.origin = oc.BgpTypes_BgpOriginAttrType_IGP
			case 1:
				attrSet.origin = oc.BgpTypes_BgpOriginAttrType_EGP
			case 2:
				attrSet.origin = oc.BgpTypes_BgpOriginAttrType_INCOMPLETE
			default:
				log.Errorf("BGP: Unrecognized origin attribute value: %v", origin)
			}
		case *api.MultiExitDiscAttribute:
			hasMED = true
			attrSet.med = m.GetMed()
		case *api.LocalPrefAttribute:
			hasLocalPref = true
			attrSet.localPref = m.GetLocalPref()
		case *api.AsPathAttribute:
			hasASPathAttribute = true
			asSegments = m.GetSegments()
			attrSet.asPath = asSegmentsToString(asSegments)
		}
	}
	if hasCommunity {
		route.SetCommunityIndex(commIndex)
	}
	if hasOrigin || hasMED || hasLocalPref || hasASPathAttribute {
		attrSetIndex := t.attrSetTracker.getOrAllocIndex(attrSet)
		route.SetAttrIndex(attrSetIndex)
		attrSetOC := rib.GetOrCreateAttrSet(attrSetIndex)
		if hasOrigin {
			attrSetOC.SetOrigin(attrSet.origin)
		}
		if hasMED {
			attrSetOC.SetMed(attrSet.med)
		}
		if hasLocalPref {
			attrSetOC.SetLocalPref(attrSet.localPref)
		}
		if hasASPathAttribute {
			for i, s := range asSegments {
				segmentOC := attrSetOC.GetOrCreateAsSegment(uint32(i))
				segmentOC.SetType(convertSegmentTypeToOC(s.GetType()))
				segmentOC.SetMember(s.GetNumbers())
			}
		}
	}
}

type ocRIBRoute interface {
	SetCommunityIndex(uint64)
	SetAttrIndex(uint64)
}

type ribAttrSet struct {
	origin    oc.E_BgpTypes_BgpOriginAttrType
	med       uint32
	localPref uint32
	asPath    string
}

// ocRIBAttrIndicesTracker is used to track and populate BGP RIB attribute
// information.
type ocRIBAttrIndicesTracker[T comparable] struct {
	lastIndex         uint64
	allocs            map[T]uint64
	indiciesInUse     map[uint64]T
	indiciesInUseNext map[uint64]T
}

func newOCRIBAttrIndices[T comparable]() *ocRIBAttrIndicesTracker[T] {
	return &ocRIBAttrIndicesTracker[T]{
		lastIndex:         0,
		allocs:            map[T]uint64{},
		indiciesInUse:     map[uint64]T{},
		indiciesInUseNext: map[uint64]T{},
	}
}

// beginAllocation indicates a fresh round of attribute index allocation.
func (r *ocRIBAttrIndicesTracker[T]) beginAllocation() {
	r.indiciesInUseNext = map[uint64]T{}
}

// completeAllocation indicates to finish allocation and delete unused indices.
func (r *ocRIBAttrIndicesTracker[T]) completeAllocation() {
	for i, key := range r.indiciesInUse {
		if _, ok := r.indiciesInUseNext[i]; !ok {
			delete(r.allocs, key)
		}
	}

	r.indiciesInUse = r.indiciesInUseNext
	r.indiciesInUseNext = map[uint64]T{}
}

func (r *ocRIBAttrIndicesTracker[T]) getOrAllocIndex(key T) uint64 {
	i, ok := r.allocs[key]
	if ok {
		r.indiciesInUseNext[i] = key
		return i
	}

	if len(r.indiciesInUse) == math.MaxInt {
		log.Fatal("Way too many unities")
	}
	r.lastIndex++
	//revive:disable:empty-block while loop usage.
	for _, ok := r.indiciesInUse[r.lastIndex]; ok; r.lastIndex++ {
	}
	//revive:enable:empty-block
	i = r.lastIndex

	r.allocs[key] = i
	r.indiciesInUse[i] = key
	r.indiciesInUseNext[i] = key
	return i
}
