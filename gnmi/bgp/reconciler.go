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
	"errors"
	"fmt"
	"sync"

	"github.com/golang/glog"
	"github.com/openconfig/lemming/gnmi/oc"
	api "github.com/wenovus/gobgp/v3/api"
	"github.com/wenovus/gobgp/v3/pkg/bgpconfig"
	"github.com/wenovus/gobgp/v3/pkg/server"
	"golang.org/x/net/context"
)

const (
	gracefulRestart = false
)

type bgpReconciler struct {
	bgpServer     *server.BgpServer
	currentConfig *bgpconfig.BgpConfigSet

	bgpStarted bool
}

// newBgpReconciler creates a new BGP reconciler that reconciles between
// intended and applied configuration.
//
// It assumes that the bgpServer has been started WITHOUT any initial calls
// made to it.
func newBgpReconciler(bgpServer *server.BgpServer) *bgpReconciler {
	return &bgpReconciler{
		bgpServer:     bgpServer,
		currentConfig: &bgpconfig.BgpConfigSet{},
	}
}

// reconcile examines the difference between the intended and applied
// configuration, and makes GoBGP API calls accordingly to update the applied
// configuration in the direction of intended configuration.
func (r *bgpReconciler) reconcile(intended, applied *oc.NetworkInstance_Protocol_Bgp, appliedMu *sync.Mutex) error {
	appliedMu.Lock()
	defer appliedMu.Unlock()

	intendedGlobal := intended.GetOrCreateGlobal()
	newConfig := intendedToGoBGP(intended)

	bgpShouldStart := intendedGlobal.As != nil && intendedGlobal.RouterId != nil
	switch {
	case bgpShouldStart && !r.bgpStarted:
		glog.V(1).Info("Starting BGP")
		var err error
		r.currentConfig, err = InitialConfig(context.Background(), applied, r.bgpServer, newConfig, gracefulRestart)
		if err != nil {
			return fmt.Errorf("Failed to apply initial BGP configuration %v", newConfig)
		} else {
			r.bgpStarted = true
		}
	case !bgpShouldStart && r.bgpStarted:
		glog.V(1).Info("Stopping BGP")
		if err := r.bgpServer.StopBgp(context.Background(), &api.StopBgpRequest{}); err != nil {
			return errors.New("Failed to stop BGP service")
		} else {
			r.bgpStarted = false
		}
		r.currentConfig = &bgpconfig.BgpConfigSet{}
		*applied = oc.NetworkInstance_Protocol_Bgp{}
		applied.PopulateDefaults()
	case r.bgpStarted:
		glog.V(1).Info("Updating BGP")
		var err error
		r.currentConfig, err = UpdateConfig(context.Background(), applied, r.bgpServer, r.currentConfig, newConfig)
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
func intendedToGoBGP(bgpoc *oc.NetworkInstance_Protocol_Bgp) *bgpconfig.BgpConfigSet {
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
			Transport: bgpconfig.Transport{
				Config: bgpconfig.TransportConfig{
					RemotePort: listenPort,
				},
			},
		})
	}

	return bgpConfig
}
