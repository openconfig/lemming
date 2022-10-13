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
	api "github.com/wenovus/gobgp/v3/api"
	"github.com/wenovus/gobgp/v3/pkg/bgpconfig"
	"github.com/wenovus/gobgp/v3/pkg/server"
)

const (
	gracefulRestart = false
)

// NewGoBGPTaskDecl creates a new GoBGP task using the declarative configuration style.
func NewGoBGPTaskDecl() *reconciler.BuiltReconciler {
	return reconciler.NewBuilder("gobgp-decl").WithStart(startGoBGPFuncDecl).Build()
}

func setIfNotZero[T any](setter func(T), v T) {
	if !reflect.ValueOf(v).IsZero() {
		setter(v)
	}
}

func updateState(yclient *ygnmi.Client, appliedBgp *oc.NetworkInstance_Protocol_Bgp) {
	log.V(1).Infof("BGP task: updating state")
	if _, err := gnmiclient.Replace(context.Background(), yclient, ocpath.Root().NetworkInstance(DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp().State(), appliedBgp); err != nil {
		log.Errorf("BGP failed to update state: %v", err)
	}
}

func startGoBGPFuncDecl(ctx context.Context, yclient *ygnmi.Client) error {
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

	// TODO(wenbli): This should look cleaner.
	bgpHandler := bgpReconciler{
		bgpServer:     server.NewBgpServer(),
		currentConfig: &bgpconfig.BgpConfigSet{},
	}
	go bgpHandler.bgpServer.Serve()

	// monitor the change of the peer state
	if err := bgpHandler.bgpServer.WatchEvent(context.Background(), &api.WatchEventRequest{Peer: &api.WatchEventRequest_Peer{}}, func(r *api.WatchEventResponse) {
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
			if err := bgpHandler.reconcile(intendedBgp, appliedBgp, &appliedBgpMu); err != nil {
				log.Errorf("BGP failed to reconcile: %v", err)
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

	return nil
}
