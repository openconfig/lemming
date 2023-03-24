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
	api "github.com/wenovus/gobgp/v3/api"
	"github.com/wenovus/gobgp/v3/pkg/bgpconfig"
	"github.com/wenovus/gobgp/v3/pkg/server"
	"github.com/wenovus/gobgp/v3/pkg/zebra"
)

const (
	gracefulRestart = false
)

// NewGoBGPTaskDecl creates a new GoBGP task using the declarative configuration style.
func NewGoBGPTaskDecl(zapiURL string) *reconciler.BuiltReconciler {
	return reconciler.NewBuilder("gobgp-decl").WithStart(newBgpDeclTask(zapiURL).startGoBGPFuncDecl).Build()
}

func updateState(yclient *ygnmi.Client, appliedBgp *oc.NetworkInstance_Protocol_Bgp) {
	log.V(1).Infof("BGP task: updating state")
	if _, err := gnmiclient.Replace(context.Background(), yclient, ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp().State(), appliedBgp); err != nil {
		log.Errorf("BGP failed to update state: %v", err)
	}
}

// bgpDeclTask can be used to create a reconciler-compatible BGP task.
type bgpDeclTask struct {
	zapiURL       string
	bgpServer     *server.BgpServer
	currentConfig *bgpconfig.BgpConfigSet

	bgpStarted bool
}

// newBgpDeclTask creates a new bgpDeclTask.
func newBgpDeclTask(zapiURL string) *bgpDeclTask {
	return &bgpDeclTask{zapiURL: zapiURL}
}

// startGoBGPFuncDecl starts a GoBGP server.
func (t *bgpDeclTask) startGoBGPFuncDecl(_ context.Context, yclient *ygnmi.Client) error {
	b := &ocpath.Batch{}
	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()
	b.AddPaths(
		bgpPath.Global().As().Config().PathStruct(),
		bgpPath.Global().RouterId().Config().PathStruct(),
		bgpPath.NeighborAny().PeerAs().Config().PathStruct(),
		bgpPath.NeighborAny().NeighborAddress().Config().PathStruct(),
	)

	appliedRoot := &oc.Root{}
	// appliedBgp is the SoT for BGP applied configuration. It is maintained locally by the task.
	appliedBgp := appliedRoot.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").GetOrCreateBgp()
	appliedBgp.PopulateDefaults()
	var appliedBgpMu sync.Mutex

	bgpServer := server.NewBgpServer()
	if log.V(2) {
		if err := bgpServer.SetLogLevel(context.Background(), &api.SetLogLevelRequest{
			Level: api.SetLogLevelRequest_DEBUG,
		}); err != nil {
			log.Errorf("Error setting GoBGP log level: %v", err)
		}
	}
	go bgpServer.Serve()

	// monitor the change of the peer state
	if err := bgpServer.WatchEvent(context.Background(), &api.WatchEventRequest{Peer: &api.WatchEventRequest_Peer{}}, func(r *api.WatchEventResponse) {
		appliedBgpMu.Lock()
		defer appliedBgpMu.Unlock()
		if p := r.GetPeer(); p != nil && p.Type == api.WatchEventResponse_PeerEvent_STATE {
			log.V(1).Info("Got peer event update:", p)
			ps := p.GetPeer().State

			neigh := appliedBgp.GetOrCreateNeighbor(ps.NeighborAddress)

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

			updateState(yclient, appliedBgp)
		}
	}); err != nil {
		return fmt.Errorf("goBgpTask failed to initialize due to error: %v", err)
	}

	// Initialize values required for reconile to be called.
	t.bgpServer = bgpServer
	t.currentConfig = &bgpconfig.BgpConfigSet{}

	bgpWatcher := ygnmi.Watch(
		context.Background(),
		yclient,
		b.Config(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}

			intendedBgp := rootVal.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").GetOrCreateBgp()
			if err := t.reconcile(intendedBgp, appliedBgp, &appliedBgpMu); err != nil {
				log.Errorf("GoBGP failed to reconcile: %v", err)
				// TODO(wenbli): Instead of stopping BGP, we should simply keep trying.
				return err
			}

			appliedBgpMu.Lock()
			updateState(yclient, appliedBgp)
			appliedBgpMu.Unlock()

			return ygnmi.Continue
		},
	)

	go func() {
		if _, err := bgpWatcher.Await(); err != nil {
			log.Warningf("GoBGP Task's watcher has stopped: %v", err)
		}
	}()

	if log.V(1) {
		// Periodically print the BGP table.
		// TODO(wenbli): Put this in the BGP RIB schema.
		go func() {
			tick := time.NewTicker(5 * time.Second)
			for range tick.C {
				if err := bgpServer.ListPath(context.Background(), &api.ListPathRequest{
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

				if err := bgpServer.ListPath(context.Background(), &api.ListPathRequest{
					TableType: api.TableType_LOCAL,
					Family: &api.Family{
						Afi:  api.Family_AFI_IP,
						Safi: api.Family_SAFI_UNICAST,
					},
				}, func(d *api.Destination) {
					log.V(1).Infof("GoBGP local table path: %v", d)
				}); err != nil {
					log.Errorf("GoBGP ListPath call failed (local table): %v", err)
				} else {
					log.V(1).Info("GoBGP ListPath call completed (local table)")
				}
			}
		}()
	}

	return nil
}

// reconcile examines the difference between the intended and applied
// configuration, and makes GoBGP API calls accordingly to update the applied
// configuration in the direction of intended configuration.
func (t *bgpDeclTask) reconcile(intended, applied *oc.NetworkInstance_Protocol_Bgp, appliedMu *sync.Mutex) error {
	appliedMu.Lock()
	defer appliedMu.Unlock()

	intendedGlobal := intended.GetOrCreateGlobal()
	newConfig := intendedToGoBGP(intended, t.zapiURL)

	bgpShouldStart := intendedGlobal.As != nil && intendedGlobal.RouterId != nil
	switch {
	case bgpShouldStart && !t.bgpStarted:
		log.V(1).Info("Starting BGP")
		var err error
		t.currentConfig, err = InitialConfig(context.Background(), applied, t.bgpServer, newConfig, gracefulRestart)
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
		*applied = oc.NetworkInstance_Protocol_Bgp{}
		applied.PopulateDefaults()
	case t.bgpStarted:
		log.V(1).Info("Updating BGP")
		var err error
		t.currentConfig, err = UpdateConfig(context.Background(), applied, t.bgpServer, t.currentConfig, newConfig)
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
func intendedToGoBGP(bgpoc *oc.NetworkInstance_Protocol_Bgp, zapiURL string) *bgpconfig.BgpConfigSet {
	bgpConfig := &bgpconfig.BgpConfigSet{}
	global := bgpoc.GetOrCreateGlobal()

	bgpConfig.Global.Config.As = global.GetAs()
	bgpConfig.Global.Config.RouterId = global.GetRouterId()
	bgpConfig.Global.Config.Port = listenPort

	bgpConfig.Neighbors = []bgpconfig.Neighbor{}
	for neighAddr, neigh := range bgpoc.Neighbor {
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
					RemotePort: listenPort,
				},
			},
		})
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

	return bgpConfig
}
