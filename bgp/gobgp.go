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
	"github.com/wenovus/gobgp/v3/pkg/config"
	"github.com/wenovus/gobgp/v3/pkg/server"
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

// NewGoBGPTask creates a new GoBGP task using the declarative configuration style.
func NewGoBGPTask(targetName, zapiURL string, listenPort uint16) *reconciler.BuiltReconciler {
	gobgpTask := newBgpTask(targetName, zapiURL, listenPort)
	return reconciler.NewBuilder("gobgp-decl").WithStart(gobgpTask.start).WithStop(gobgpTask.stop).Build()
}

// bgpTask can be used to create a reconciler-compatible BGP task.
type bgpTask struct {
	targetName    string
	zapiURL       string
	bgpServer     *server.BgpServer
	currentConfig *config.BgpConfigSet
	listenPort    uint16

	bgpStarted bool

	yclient *ygnmi.Client

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
	appliedBGP.PopulateDefaults()
	appliedRoutingPolicy := appliedState.GetOrCreateRoutingPolicy()
	appliedRoutingPolicy.PopulateDefaults()

	return &bgpTask{
		targetName: targetName,
		zapiURL:    zapiURL,
		listenPort: listenPort,

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
func (t *bgpTask) start(_ context.Context, yclient *ygnmi.Client) error {
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
		RoutingPolicyPath.DefinedSets().PrefixSetAny().PrefixAny().IpPrefix().Config().PathStruct(),
		RoutingPolicyPath.DefinedSets().PrefixSetAny().PrefixAny().MasklengthRange().Config().PathStruct(),
		// -- community sets
		ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySetAny().CommunityMember().Config().PathStruct(),
		ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().CommunitySetAny().MatchSetOptions().Config().PathStruct(),
		// -- AS path sets
		ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().AsPathSetAny().AsPathSetName().Config().PathStruct(),
		ocpath.Root().RoutingPolicy().DefinedSets().BgpDefinedSets().AsPathSetAny().AsPathSetMember().Config().PathStruct(),
	)

	if err := t.createNewGoBGPServer(context.Background()); err != nil {
		log.Errorf("Failed to start GoBGP server: %v", err)
	}

	// Initialize values required for reconile to be called.
	t.currentConfig = &config.BgpConfigSet{}

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

	return nil
}

// reconcile examines the difference between the intended and applied
// configuration, and makes GoBGP API calls accordingly to update the applied
// configuration in the direction of intended configuration.
func (t *bgpTask) reconcile(intended *oc.Root) error {
	intendedBGP := intended.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).GetOrCreateBgp()
	intendedPolicy := intended.GetOrCreateRoutingPolicy()
	newConfig := intendedToGoBGP(intendedBGP, intendedPolicy, t.zapiURL, t.listenPort)

	intendedGlobal := intendedBGP.GetOrCreateGlobal()
	bgpShouldStart := intendedGlobal.As != nil && intendedGlobal.RouterId != nil
	switch {
	case bgpShouldStart && !t.bgpStarted:
		log.V(1).Info("Starting BGP")
		var err error
		t.currentConfig, err = config.InitialConfig(context.Background(), t.bgpServer, newConfig, gracefulRestart)
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
		t.currentConfig = &config.BgpConfigSet{}
		*t.appliedBGP = oc.NetworkInstance_Protocol_Bgp{}
		t.appliedBGP.PopulateDefaults()
	case t.bgpStarted:
		log.V(1).Info("Updating BGP")
		var err error
		t.currentConfig, err = config.UpdateConfig(context.Background(), t.bgpServer, t.currentConfig, newConfig)
		if err != nil {
			return fmt.Errorf("Failed to update BGP service: %v", newConfig)
		}
	default:
		// Waiting for BGP to be startable.
		return nil
	}

	return nil
}

// updateAppliedState is the ONLY function that's called when updating the appliedState.
//
// The input function is expected to make modifications to the applied state,
// which then this function will use to update the central cache.
func (t *bgpTask) updateAppliedState(f func() error) error {
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

func updateAppliedStateHelper[T ygot.GoStruct](yclient *ygnmi.Client, path ygnmi.SingletonQuery[T], appliedState T) {
	if _, err := gnmiclient.Replace(context.Background(), yclient, path, appliedState); err != nil {
		log.Errorf("BGP failed to update state at path %v: %v", path, err)
	}
}

// createNewGoBGPServer creates and starts a new GoBGP Server.
func (t *bgpTask) createNewGoBGPServer(context.Context) error {
	t.bgpServer = server.NewBgpServer()

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

	go func() {
		tick := time.NewTicker(5 * time.Second)
		for range tick.C {
			t.updateRIBs()
		}
	}()

	return nil
}

// updateRIBs updates the BGP RIBs.
func (t *bgpTask) updateRIBs() {
	if err := t.bgpServer.ListPath(context.Background(), &api.ListPathRequest{
		TableType: api.TableType_GLOBAL,
		Family: &api.Family{
			Afi:  api.Family_AFI_IP,
			Safi: api.Family_SAFI_UNICAST,
		},
	}, func(d *api.Destination) {
		log.V(1).Infof("%s: GoBGP global table path: %v", t.targetName, d)
	}); err != nil {
		if err.Error() != "bgp server hasn't started yet" {
			log.Errorf("GoBGP ListPath call failed (global table): %v", err)
		}
	} else {
		log.V(1).Info("GoBGP ListPath call completed (global table)")
	}

	t.updateAppliedState(func() error {
		v4uni := t.appliedBGP.GetOrCreateRib().GetOrCreateAfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).GetOrCreateIpv4Unicast()
		rib := t.appliedBGP.GetOrCreateRib()
		ribattrs := &ocRIBAttrIndices{}

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
						origin = oc.UnionString(path.NeighborIp)
					}
					// TODO: this ID should match the ID in adj-rib-in-post.
					ribattrs.populateRIBAttrs(path, rib, locRib.GetOrCreateRoute(route.Prefix, origin, uint32(j)))
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
						// TODO: this ID should be retrieved from the update message.
						ribattrs.populateRIBAttrs(path, rib, neighContainer.GetOrCreateAdjRibInPre().GetOrCreateRoute(route.Prefix, uint32(j)))
						if !path.Filtered {
							ribattrs.populateRIBAttrs(path, rib, neighContainer.GetOrCreateAdjRibInPost().GetOrCreateRoute(route.Prefix, uint32(j)))
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
						ribattrs.populateRIBAttrs(path, rib, neighContainer.GetOrCreateAdjRibOutPre().GetOrCreateRoute(route.Prefix, uint32(j)))
						if !path.Filtered {
							ribattrs.populateRIBAttrs(path, rib, neighContainer.GetOrCreateAdjRibOutPost().GetOrCreateRoute(route.Prefix, uint32(j)))
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
func (t *bgpTask) queryTable(neighbor, tableName string, tableType api.TableType, f func(route []*api.Destination)) {
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
		log.V(0).Infof("%s: GoBGP %s table path (neighbor if applicable: %q): %v", t.targetName, tableName, neighbor, d)
	}); err != nil {
		if err.Error() != "bgp server hasn't started yet" {
			log.Errorf("GoBGP ListPath call failed (%s table): %v", tableType, err)
		}
	} else {
		log.V(1).Info("GoBGP ListPath call completed (%s table)", tableName)
		f(routes)
	}
}

type ocRIBRoute interface {
	SetCommunityIndex(uint64)
}

type ocRIBAttrIndices struct {
	commIndex uint64
}

// populateRIBAttrs populates path attributes of routes in the RIB.
//
// TODO(wenbli): Keep a cache and keep indices stable rather than changing.
func (ribattrs *ocRIBAttrIndices) populateRIBAttrs(path *api.Path, rib *oc.NetworkInstance_Protocol_Bgp_Rib, r ocRIBRoute) {
	for _, attr := range path.GetPattrs() {
		m, err := attr.UnmarshalNew()
		if err != nil {
			log.Errorf("BGP: Unable to unmarshal a GoBGP path attribute")
		}
		switch m := m.(type) {
		case *api.CommunitiesAttribute:
			if comms := m.GetCommunities(); len(comms) > 0 {
				rib.GetOrCreateCommunity(ribattrs.commIndex).SetCommunity(communitiesToOC(comms))
				ribattrs.commIndex++
			}
		}
	}
	r.SetCommunityIndex(ribattrs.commIndex)
}
