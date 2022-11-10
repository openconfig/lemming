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

package fakedevice

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/gnmi/reconciler"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	api "github.com/osrg/gobgp/v3/api"
	"github.com/osrg/gobgp/v3/pkg/server"
)

const (
	listenPort = 1234
)

// NewGoBGPTask creates a new GoBGP task.
func NewGoBGPTask() *reconciler.BuiltReconciler {
	return reconciler.NewBuilder("gobgp").WithStart(startGoBGPFunc).Build()
}

// startGoBGPTask tries to establish a simple BGP session using GoBGP. It returns an error if initialization failed.
//
// TODO(wenbli): Break this function up.
//
// # How to achieve declarative configuration for GoBGP
//
// Goal:
//   - Declarative configuration: Get the system into the correct state regardless of what the diff of the intended config against the current applied config is.
//
// Requirements:
//  1. When there is a change in intended config, the system must arrive at an eventually-consistent state.
//  2. Make sure that all applied config and derived state are updated correctly whether there is a passive state change (i.e. through the event watcher), or active state change (i.e. when there is a change in intended config)
//
// Assumptions:
//   - All possible actions to take within this task (BGP in this case) can be put in DAG order.
//   - The number of possible actions is low such that maintaining the DAG
//     order of these actions is not burdensome. This is a poor assumption -- we
//     may very well have to use an alternative method such as an event loop for
//     managing the scheduling of actions such that we don't have to manually
//     maintain their order.
//
// Algorithm:
// Process all possible actions in DAG order so that if a previous one happens that unblocks later ones, those later ones will actually execute.
//
//  0. Maintain a view of the current intended and applied configs in memory.
//     For each new intended config update:
//  1. Identify the DAG order of library actions and the paths associated with each of them.
//  2. Process the actions in DAG order. The action to take depends on the current view of the config. The results of the action determines how we will update the current view of the applied config (which we then forward to the central DB).
//     e.g. for Global if the intended config matches the applied config we will simply skip this step. If not we will actually do the action (either start, stop, or stop-start (aka. update)), and if it succeeds, update the applied config both in the DB as well as the current view of the applied config.
//     When there is a dependency, the later actions will NOT directly depend on the results of the previous actions, but will just look at the current config view to determine the appropriate action.
//     e.g. for peers, if the global setting has been set up, then we can create the peers and update the applied config if it succeeds, but if not then we don't do anything.
//     ; however, if the global setting hasn't been set up, we actually need to erase the entirety of the applied config. This is because the watcher doesn't tell us this information.
func startGoBGPFunc(ctx context.Context, yclient *ygnmi.Client) error {
	b := &ocpath.Batch{}
	bgpPath := ocpath.Root().NetworkInstance(DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()
	b.AddPaths(
		bgpPath.Global().As().Config().PathStruct(),
		bgpPath.Global().RouterId().Config().PathStruct(),
		bgpPath.NeighborAny().PeerAs().Config().PathStruct(),
		bgpPath.NeighborAny().NeighborAddress().Config().PathStruct(),
	)

	s := server.NewBgpServer()
	go s.Serve()

	// The code below implements declarative configuration for setting up a basic BGP session.

	appliedRoot := &oc.Root{}
	// appliedBgp is the SoT for BGP applied configuration. It is maintained locally by the task.
	appliedBgp := appliedRoot.GetOrCreateNetworkInstance(DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").GetOrCreateBgp()
	appliedBgp.PopulateDefaults()
	var appliedBgpMu sync.Mutex

	// updateAppliedConfig computes the diff between a previous applied
	// configuration and the current SoT, and sends the updates to the
	// central DB.
	updateAppliedConfig := func(prevApplied *oc.NetworkInstance_Protocol_Bgp, grabLock bool) bool {
		if grabLock {
			appliedBgpMu.Lock()
			defer appliedBgpMu.Unlock()
		}
		no, err := ygot.Diff(prevApplied, appliedBgp)
		if err != nil {
			log.Errorf("goBgpTask: error while creating update notification for updating applied configuration: %v", err)
			return false
		}
		if len(no.GetUpdate())+len(no.GetDelete()) > 0 {
			_, err := gnmiclient.Replace(ctx, yclient, bgpPath.State(), appliedBgp)
			if err != nil {
				log.Errorf("goBgpTask: error while writing update to applied configuration: %v", err)
				return false
			}
		}
		return true
	}

	// monitor the change of the peer state
	if err := s.WatchEvent(context.Background(), &api.WatchEventRequest{Peer: &api.WatchEventRequest_Peer{}}, func(r *api.WatchEventResponse) {
		if p := r.GetPeer(); p != nil && p.Type == api.WatchEventResponse_PeerEvent_STATE {
			log.V(1).Info("Got peer event update:", p)
			ps := p.GetPeer().State
			appliedBgpMu.Lock()
			defer appliedBgpMu.Unlock()
			prevAppliedIntf, err := ygot.DeepCopy(appliedBgp)
			if err != nil {
				log.Fatalf("goBgpTask: Could not copy applied configuration: %v", err)
			}
			prevApplied := prevAppliedIntf.(*oc.NetworkInstance_Protocol_Bgp)
			if neigh, ok := appliedBgp.Neighbor[ps.NeighborAddress]; ok {
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
			} else {
				log.Warning("Received peer update from an unknown peer:", ps.NeighborAddress)
			}
			if success := updateAppliedConfig(prevApplied, false); !success {
				log.Errorf("goBgpTask: updating applied configuration failed")
			}
		}
	}); err != nil {
		return fmt.Errorf("goBgpTask failed to initialize due to error: %v", err)
	}

	var global api.Global
	global.ListenPort = listenPort

	bgpWatcher := ygnmi.Watch(
		context.Background(),
		yclient,
		b.Config(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}
			intended := rootVal.GetOrCreateNetworkInstance(DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").GetOrCreateBgp()

			processBgp := func() {
				log.V(1).Info("Processing BGP update")
				// Visit paths in action DAG order.
				// - The output at each stage are
				//   1. The GoBGP action to take.
				//   2. The paths to update.

				// Action 1: StartBGP / StopBGP / no-op
				//   - as-path, route-id
				intendedGlobal, appliedGlobal := intended.GetOrCreateGlobal(), appliedBgp.GetOrCreateGlobal()
				if !reflect.DeepEqual(intendedGlobal.As, appliedGlobal.As) || !reflect.DeepEqual(intendedGlobal.RouterId, appliedGlobal.RouterId) {
					switch {
					case intendedGlobal.As == nil || intendedGlobal.RouterId == nil:
						if hasGlobal := appliedGlobal.As != nil && appliedGlobal.RouterId != nil; hasGlobal {
							log.V(1).Info("Stopping BGP")
							if err := s.StopBgp(context.Background(), &api.StopBgpRequest{}); err != nil {
								log.Fatal(err)
							} else {
								appliedGlobal.As = nil
								appliedGlobal.RouterId = nil
							}
							for neighAddr := range appliedBgp.Neighbor {
								delete(appliedBgp.Neighbor, neighAddr)
							}
						}
					case appliedGlobal.As == nil:
						log.V(1).Info("Starting BGP")
						global.Asn = *intendedGlobal.As
						global.RouterId = *intendedGlobal.RouterId
						err := s.StartBgp(context.Background(), &api.StartBgpRequest{
							Global: &global,
						})
						if err != nil {
							log.Error(err)
						} else {
							appliedGlobal.As = ygot.Uint32(*intendedGlobal.As)
							appliedGlobal.RouterId = ygot.String(*intendedGlobal.RouterId)
						}
					default:
						log.V(1).Info("Restarting BGP due to changed global configuration")
						if err := s.StopBgp(context.Background(), &api.StopBgpRequest{}); err != nil {
							log.Fatal(err)
						}
						for neighAddr := range appliedBgp.Neighbor {
							delete(appliedBgp.Neighbor, neighAddr)
						}
						global.Asn = *intendedGlobal.As
						global.RouterId = *intendedGlobal.RouterId
						err := s.StartBgp(context.Background(), &api.StartBgpRequest{
							Global: &global,
						})
						if err != nil {
							log.Error(err)
						} else {
							appliedGlobal.As = ygot.Uint32(*intendedGlobal.As)
							appliedGlobal.RouterId = ygot.String(*intendedGlobal.RouterId)
						}
					}
				}
				// States used by later actions.
				hasGlobal := appliedGlobal.As != nil && appliedGlobal.RouterId != nil

				// Action 2: AddPeer / UpdatePeer / DeletePeer
				//   - peer-as, neighbor-addr

				// Delete non-existent neighbours.
				for neighAddr := range appliedBgp.Neighbor {
					if neigh, ok := intended.Neighbor[neighAddr]; !ok || !hasGlobal || ok && (neigh.PeerAs == nil) {
						log.V(1).Info("Deleting BGP peer: ", neighAddr)
						if err := s.DeletePeer(context.Background(), &api.DeletePeerRequest{
							Address: neighAddr,
						}); err != nil {
							log.Error(err)
						} else {
							delete(appliedBgp.Neighbor, neighAddr)
						}
					}
				}
				// Add/update neighbours.
				if hasGlobal {
					for neighAddr, neigh := range intended.Neighbor {
						if curNeigh, ok := appliedBgp.Neighbor[neighAddr]; ok {
							if reflect.DeepEqual(neigh.PeerAs, curNeigh.PeerAs) {
								continue
							}
							log.V(1).Info("Updating BGP peer: ", neighAddr)
							if resp, err := s.UpdatePeer(context.Background(), &api.UpdatePeerRequest{
								Peer: &api.Peer{
									Conf: &api.PeerConf{
										NeighborAddress: neighAddr,
										PeerAsn:         *neigh.PeerAs,
									},
									Transport: &api.Transport{
										RemotePort: listenPort,
									},
								},
							}); err != nil {
								log.Errorf("goBgpTask: %v", err)
								continue
							} else if resp.NeedsSoftResetIn {
								log.V(1).Info("Updating BGP peer with softResetIn: ", neighAddr)
								if _, err := s.UpdatePeer(context.Background(), &api.UpdatePeerRequest{
									Peer: &api.Peer{
										Conf: &api.PeerConf{
											NeighborAddress: neighAddr,
											PeerAsn:         *neigh.PeerAs,
										},
										Transport: &api.Transport{
											RemotePort: listenPort,
										},
									},
									DoSoftResetIn: true,
								}); err != nil {
									log.Errorf("goBgpTask: retry UpdatePeer with DoSoftResetIn: true failed: %v", err)
									continue
								}
							}
						} else {
							if neigh.PeerAs == nil {
								// Doesn't have enough information yet to create the peer.
								continue
							}
							log.V(1).Info("Adding BGP peer: ", neighAddr)
							if err := s.AddPeer(context.Background(), &api.AddPeerRequest{
								Peer: &api.Peer{
									Conf: &api.PeerConf{
										NeighborAddress: neighAddr,
										PeerAsn:         *neigh.PeerAs,
									},
									Transport: &api.Transport{
										RemotePort: listenPort,
									},
								},
							}); err != nil {
								log.Errorf("goBgpTask: %v", err)
								continue
							}
						}
						n := appliedBgp.GetOrCreateNeighbor(neighAddr)
						n.PeerAs = ygot.Uint32(*neigh.PeerAs)
					}
				}
			}

			appliedBgpMu.Lock()
			prevAppliedIntf, err := ygot.DeepCopy(appliedBgp)
			if err != nil {
				log.Fatalf("goBgpTask: Could not copy applied configuration: %v", err)
			}
			prevApplied := prevAppliedIntf.(*oc.NetworkInstance_Protocol_Bgp)
			processBgp()
			if success := updateAppliedConfig(prevApplied, false); !success {
				log.Errorf("goBgpTask: updating applied configuration failed")
			}
			appliedBgpMu.Unlock()

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
