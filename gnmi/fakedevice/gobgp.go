package fakedevice

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/coalesce"
	"github.com/openconfig/gnmi/ctree"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/lemming/gnmi/gnmit"
	"github.com/openconfig/lemming/gnmi/internal/config"
	configpath "github.com/openconfig/lemming/gnmi/internal/config/device"
	"github.com/openconfig/lemming/gnmi/internal/telemetry"
	telemetrypath "github.com/openconfig/lemming/gnmi/internal/telemetry/device"
	"github.com/openconfig/ygot/ygot"
	api "github.com/osrg/gobgp/v3/api"
	"github.com/osrg/gobgp/v3/pkg/server"
	"google.golang.org/protobuf/encoding/prototext"
)

const (
	listenPort = 179
)

// goBgpTask tries to establish a simple BGP session using GoBGP. It returns an error if initialization failed.
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
func goBgpTask(getIntendedConfig func() *config.Device, q gnmit.Queue, update gnmit.UpdateFn, target string, remove func()) error {
	bgpStatePath, _, err := ygot.ResolvePath(telemetrypath.DeviceRoot("").NetworkInstance("default").Protocol(telemetry.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp())
	if err != nil {
		return fmt.Errorf("goBgpTask failed to initialize due to error: %v", err)
	}

	s := server.NewBgpServer()
	go s.Serve()

	// The code below implements declarative configuration for setting up a basic BGP session.

	appliedRoot := &telemetry.Device{}
	// appliedBgp is the SoT for BGP applied configuration. It is maintained locally by the task.
	appliedBgp := appliedRoot.GetOrCreateNetworkInstance("default").GetOrCreateProtocol(telemetry.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").GetOrCreateBgp()
	appliedBgp.PopulateDefaults()
	var appliedBgpMu sync.Mutex

	// updateAppliedConfig computes the diff between a previous applied
	// configuration and the current SoT, and sends the updates to the
	// central DB.
	updateAppliedConfig := func(prevApplied *telemetry.NetworkInstance_Protocol_Bgp, grabLock bool) bool {
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
			log.V(1).Info("Updating BGP applied configuration: ", prototext.Format(no))
			no.Timestamp = time.Now().UnixNano()
			no.Prefix = &gpb.Path{Origin: "openconfig", Target: target, Elem: bgpStatePath.Elem}

			if err := update(no); err != nil {
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
			prevApplied := prevAppliedIntf.(*telemetry.NetworkInstance_Protocol_Bgp)
			if neigh, ok := appliedBgp.Neighbor[ps.NeighborAddress]; ok {
				found := false
				if ps.SessionState.String() == "UNKNOWN" {
					neigh.SessionState = telemetry.Bgp_Neighbor_SessionState_UNSET
					found = true
				} else {
					for enumCode, v := range neigh.SessionState.ΛMap()[reflect.TypeOf(neigh.SessionState).Name()] {
						if v.Name == ps.SessionState.String() {
							newSessionState := telemetry.E_Bgp_Neighbor_SessionState(enumCode)
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

	bgpPath := configpath.DeviceRoot("").NetworkInstance("default").Protocol(config.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()
	asPaths, _, err := ygot.ResolvePath(bgpPath.Global().As())
	if err != nil {
		return fmt.Errorf("goBgpTask failed to initialize due to error: %v", err)
	}
	routeIDPaths, _, err := ygot.ResolvePath(bgpPath.Global().RouterId())
	if err != nil {
		return fmt.Errorf("goBgpTask failed to initialize due to error: %v", err)
	}
	peerAsPaths, _, err := ygot.ResolvePath(bgpPath.NeighborAny().PeerAs())
	if err != nil {
		return fmt.Errorf("goBgpTask failed to initialize due to error: %v", err)
	}
	neighAddrPaths, _, err := ygot.ResolvePath(bgpPath.NeighborAny().NeighborAddress())
	if err != nil {
		return fmt.Errorf("goBgpTask failed to initialize due to error: %v", err)
	}

	var global api.Global
	global.ListenPort = listenPort

	go func() {
		defer remove()
		for {
			item, _, err := q.Next(context.Background())
			if coalesce.IsClosedQueue(err) {
				return
			}
			n, ok := item.(*ctree.Leaf)
			if !ok || n == nil {
				log.Errorf("goBgpTask invalid cache node: %#v", item)
				return
			}
			v := n.Value()
			no, ok := v.(*gpb.Notification)
			if !ok || no == nil {
				log.Errorf("goBgpTask invalid cache node, expected non-nil *gpb.Notification type, got: %#v", v)
				return
			}

			// Update the view of the intended config.
			// Note: We're guaranteed that whatever update came from the collector is valid since we validate before storing.
			shouldProcess := false
			for _, u := range no.Update {
				switch {
				case matchingPath(u.Path, asPaths), matchingPath(u.Path, routeIDPaths), matchingPath(u.Path, neighAddrPaths), matchingPath(u.Path, peerAsPaths):
					log.V(1).Infof("Received update path: %s", prototext.Format(u))
					shouldProcess = true
				default:
					log.V(1).Infof("goBgpTask: update path received isn't matched by any handlers: %s", prototext.Format(u.Path))
				}
			}
			for _, u := range no.Delete {
				log.V(1).Infof("Received delete path: %s", prototext.Format(u))
				switch {
				case len(u.Elem) > 0:
				case len(u.Element) > 0: //nolint:staticcheck //lint:ignore SA1019 gnmi cache currently doesn't support PathElem for deletions.
					// Since gNMI still sends delete paths using the deprecated Element field, we need to translate it into path-elems first.
					// We also need to strip the first element for origin.
					//nolint:staticcheck //lint:ignore SA1019 gnmi cache currently doesn't support PathElem for deletions.
					elems, err := pathTranslator.PathElem(u.Element[1:])
					if err != nil {
						log.Errorf("goBgpTask: failed to translate delete path: %s", prototext.Format(u))
						return
					}
					u.Elem = elems
				default:
					log.Errorf("Unhandled: delete at root: %s", prototext.Format(u))
					return
				}
				switch {
				case matchingPath(u, asPaths), matchingPath(u, routeIDPaths), matchingPath(u, neighAddrPaths), matchingPath(u, peerAsPaths):
					shouldProcess = true
				default:
					log.V(1).Infof("goBgpTask: delete path received isn't matched by any handlers: %s", prototext.Format(u))
				}
			}

			if !shouldProcess {
				continue
			}

			intendedRoot := getIntendedConfig()
			intended := intendedRoot.GetNetworkInstance("default").GetProtocol(config.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").GetBgp()
			if intended == nil {
				intended = &config.NetworkInstance_Protocol_Bgp{}
				intended.PopulateDefaults()
			}

			processBgp := func() {
				// Visit paths in action DAG order.
				// - The output at each stage are
				//   1. The GoBGP action to take.
				//   2. The paths to update.

				// Action 1: StartBGP / StopBGP / no-op
				//   - as-path, route-id
				intendedGlobal, appliedGlobal := intended.GetGlobal(), appliedBgp.GetOrCreateGlobal()
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
			prevApplied := prevAppliedIntf.(*telemetry.NetworkInstance_Protocol_Bgp)
			processBgp()
			if success := updateAppliedConfig(prevApplied, false); !success {
				log.Errorf("goBgpTask: updating applied configuration failed")
			}
			appliedBgpMu.Unlock()
		}
	}()

	return nil
}
