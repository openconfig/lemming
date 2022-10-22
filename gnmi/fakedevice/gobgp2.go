package fakedevice

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/gnmi/reconciler"
	"github.com/openconfig/ygnmi/ygnmi"
	api "github.com/wenovus/gobgp/v3/api"
	"github.com/wenovus/gobgp/v3/pkg/bgpconfig"
	"github.com/wenovus/gobgp/v3/pkg/server"
)

// NewGoBGPTaskDecl creates a new GoBGP task using the declarative configuration style.
func NewGoBGPTaskDecl(zapiURL string) *reconciler.BuiltReconciler {
	return reconciler.NewBuilder("gobgp-decl").WithStart(newBgpDeclTask(zapiURL).startGoBGPFuncDecl).Build()
}

func updateState(yclient *ygnmi.Client, appliedBgp *oc.NetworkInstance_Protocol_Bgp) {
	log.V(1).Infof("BGP task: updating state")
	if _, err := gnmiclient.Replace(context.Background(), yclient, ocpath.Root().NetworkInstance(DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp().State(), appliedBgp); err != nil {
		log.Errorf("BGP failed to update state: %v", err)
	}
}

// bgpDeclTask can be used to create a reconciler-compatible BGP task.
type bgpDeclTask struct {
	zapiURL string
}

// newBgpDeclTask creates a new bgpDeclTask.
func newBgpDeclTask(zapiURL string) *bgpDeclTask {
	return &bgpDeclTask{zapiURL: zapiURL}
}

// startGoBGPFuncDecl starts a GoBGP server.
func (t *bgpDeclTask) startGoBGPFuncDecl(ctx context.Context, yclient *ygnmi.Client) error {
	b := &ocpath.Batch{}
	bgpPath := ocpath.Root().NetworkInstance(DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()
	b.AddPaths(
		bgpPath.Global().As().Config().PathStruct(),
		bgpPath.Global().RouterId().Config().PathStruct(),
		bgpPath.NeighborAny().PeerAs().Config().PathStruct(),
		bgpPath.NeighborAny().NeighborAddress().Config().PathStruct(),
	)

	appliedRoot := &oc.Root{}
	// appliedBgp is the SoT for BGP applied configuration. It is maintained locally by the task.
	appliedBgp := appliedRoot.GetOrCreateNetworkInstance(DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").GetOrCreateBgp()
	appliedBgp.PopulateDefaults()
	var appliedBgpMu sync.Mutex

	bgpServer := server.NewBgpServer()
	if err := bgpServer.SetLogLevel(context.Background(), &api.SetLogLevelRequest{
		Level: api.SetLogLevelRequest_DEBUG,
	}); err != nil {
		log.Errorf("Error setting GoBGP log level: %v", err)
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

	recon := newBgpReconciler(
		bgpServer,
		&bgpconfig.BgpConfigSet{},
		t.zapiURL,
	)

	bgpWatcher := ygnmi.Watch(
		context.Background(),
		yclient,
		b.Config(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}

			intendedBgp := rootVal.GetOrCreateNetworkInstance(DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").GetOrCreateBgp()
			if err := recon.reconcile(intendedBgp, appliedBgp, &appliedBgpMu); err != nil {
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
