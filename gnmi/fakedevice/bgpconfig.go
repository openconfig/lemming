package fakedevice

import (
	"errors"
	"fmt"
	"sync"

	"golang.org/x/net/context"
	apb "google.golang.org/protobuf/types/known/anypb"

	"github.com/golang/glog"

	"github.com/openconfig/lemming/gnmi/oc"
	api "github.com/wenovus/gobgp/v3/api"
	"github.com/wenovus/gobgp/v3/pkg/apiutil"
	"github.com/wenovus/gobgp/v3/pkg/bgpconfig"
	"github.com/wenovus/gobgp/v3/pkg/log"
	"github.com/wenovus/gobgp/v3/pkg/packet/bgp"
	"github.com/wenovus/gobgp/v3/pkg/server"
	"github.com/wenovus/gobgp/v3/pkg/table"
)

// ReadConfigFile parses a config file into a BgpConfigSet which can be applied
// using InitialConfig and UpdateConfig.
func ReadConfigFile(configFile, configType string) (*bgpconfig.BgpConfigSet, error) {
	return bgpconfig.ReadConfigfile(configFile, configType)
}

func marshalRouteTargets(l []string) ([]*apb.Any, error) {
	rtList := make([]*apb.Any, 0, len(l))
	for _, rtString := range l {
		rt, err := bgp.ParseRouteTarget(rtString)
		if err != nil {
			return nil, err
		}
		a, err := apiutil.MarshalRT(rt)
		if err != nil {
			return nil, err
		}
		rtList = append(rtList, a)
	}
	return rtList, nil
}

func assignGlobalpolicy(ctx context.Context, bgpServer *server.BgpServer, a *bgpconfig.ApplyPolicyConfig) {
	toDefaultTable := func(r bgpconfig.DefaultPolicyType) table.RouteType {
		var def table.RouteType
		switch r {
		case bgpconfig.DEFAULT_POLICY_TYPE_ACCEPT_ROUTE:
			def = table.ROUTE_TYPE_ACCEPT
		case bgpconfig.DEFAULT_POLICY_TYPE_REJECT_ROUTE:
			def = table.ROUTE_TYPE_REJECT
		}
		return def
	}
	toPolicies := func(r []string) []*table.Policy {
		p := make([]*table.Policy, 0, len(r))
		for _, n := range r {
			p = append(p, &table.Policy{
				Name: n,
			})
		}
		return p
	}

	def := toDefaultTable(a.DefaultImportPolicy)
	ps := toPolicies(a.ImportPolicyList)
	bgpServer.SetPolicyAssignment(ctx, &api.SetPolicyAssignmentRequest{
		Assignment: table.NewAPIPolicyAssignmentFromTableStruct(&table.PolicyAssignment{
			Name:     table.GLOBAL_RIB_NAME,
			Type:     table.POLICY_DIRECTION_IMPORT,
			Policies: ps,
			Default:  def,
		}),
	})

	def = toDefaultTable(a.DefaultExportPolicy)
	ps = toPolicies(a.ExportPolicyList)
	bgpServer.SetPolicyAssignment(ctx, &api.SetPolicyAssignmentRequest{
		Assignment: table.NewAPIPolicyAssignmentFromTableStruct(&table.PolicyAssignment{
			Name:     table.GLOBAL_RIB_NAME,
			Type:     table.POLICY_DIRECTION_EXPORT,
			Policies: ps,
			Default:  def,
		}),
	})

}

func addPeerGroups(ctx context.Context, bgpServer *server.BgpServer, addedPg []bgpconfig.PeerGroup) {
	for _, pg := range addedPg {
		bgpServer.Log().Info("Add PeerGroup",
			log.Fields{
				"Topic": "config",
				"Key":   pg.Config.PeerGroupName,
			})

		if err := bgpServer.AddPeerGroup(ctx, &api.AddPeerGroupRequest{
			PeerGroup: bgpconfig.NewPeerGroupFromConfigStruct(&pg),
		}); err != nil {
			bgpServer.Log().Warn("Failed to add PeerGroup",
				log.Fields{
					"Topic": "config",
					"Key":   pg.Config.PeerGroupName,
					"Error": err})
		}
	}
}

func deletePeerGroups(ctx context.Context, bgpServer *server.BgpServer, deletedPg []bgpconfig.PeerGroup) {
	for _, pg := range deletedPg {
		bgpServer.Log().Info("delete PeerGroup",
			log.Fields{
				"Topic": "config",
				"Key":   pg.Config.PeerGroupName})
		if err := bgpServer.DeletePeerGroup(ctx, &api.DeletePeerGroupRequest{
			Name: pg.Config.PeerGroupName,
		}); err != nil {
			bgpServer.Log().Warn("Failed to delete PeerGroup",
				log.Fields{
					"Topic": "config",
					"Key":   pg.Config.PeerGroupName,
					"Error": err})
		}
	}
}

func updatePeerGroups(ctx context.Context, bgpServer *server.BgpServer, updatedPg []bgpconfig.PeerGroup) bool {
	for _, pg := range updatedPg {
		bgpServer.Log().Info("update PeerGroup",
			log.Fields{
				"Topic": "config",
				"Key":   pg.Config.PeerGroupName})
		if u, err := bgpServer.UpdatePeerGroup(ctx, &api.UpdatePeerGroupRequest{
			PeerGroup: bgpconfig.NewPeerGroupFromConfigStruct(&pg),
		}); err != nil {
			bgpServer.Log().Warn("Failed to update PeerGroup",
				log.Fields{
					"Topic": "config",
					"Key":   pg.Config.PeerGroupName,
					"Error": err})
		} else {
			return u.NeedsSoftResetIn
		}
	}
	return false
}

func addDynamicNeighbors(ctx context.Context, bgpServer *server.BgpServer, dynamicNeighbors []bgpconfig.DynamicNeighbor) {
	for _, dn := range dynamicNeighbors {
		bgpServer.Log().Info("Add Dynamic Neighbor to PeerGroup",
			log.Fields{
				"Topic":  "config",
				"Key":    dn.Config.PeerGroup,
				"Prefix": dn.Config.Prefix})
		if err := bgpServer.AddDynamicNeighbor(ctx, &api.AddDynamicNeighborRequest{
			DynamicNeighbor: &api.DynamicNeighbor{
				Prefix:    dn.Config.Prefix,
				PeerGroup: dn.Config.PeerGroup,
			},
		}); err != nil {
			bgpServer.Log().Warn("Failed to add Dynamic Neighbor to PeerGroup",
				log.Fields{
					"Topic":  "config",
					"Key":    dn.Config.PeerGroup,
					"Prefix": dn.Config.Prefix,
					"Error":  err})
		}
	}
}

func addNeighbors(ctx context.Context, applied *oc.NetworkInstance_Protocol_Bgp, bgpServer *server.BgpServer, added []bgpconfig.Neighbor) {
	for _, p := range added {
		bgpServer.Log().Info("Add Peer",
			log.Fields{
				"Topic": "config",
				"Key":   p.State.NeighborAddress})
		if err := bgpServer.AddPeer(ctx, &api.AddPeerRequest{
			Peer: bgpconfig.NewPeerFromConfigStruct(&p),
		}); err != nil {
			bgpServer.Log().Warn("Failed to add Peer",
				log.Fields{
					"Topic": "config",
					"Key":   p.State.NeighborAddress,
					"Error": err})
		} else {
			addr, err := p.ExtractNeighborAddress()
			if err != nil {
				glog.Errorf("BGP task: unexpected internal error: add peer succeeded on invalid address: %v", err)
			}
			neighoc, err := applied.NewNeighbor(addr)
			if err != nil {
				glog.Errorf("Internal error: BGP Task: %v", err)
				continue
			}
			setIfNotZero(neighoc.SetPeerAs, p.Config.PeerAs)
		}
	}
}

func deleteNeighbors(ctx context.Context, bgpServer *server.BgpServer, deleted []bgpconfig.Neighbor) {
	for _, p := range deleted {
		bgpServer.Log().Info("Delete Peer",
			log.Fields{
				"Topic": "config",
				"Key":   p.State.NeighborAddress})
		if err := bgpServer.DeletePeer(ctx, &api.DeletePeerRequest{
			Address: p.State.NeighborAddress,
		}); err != nil {
			bgpServer.Log().Warn("Failed to delete Peer",
				log.Fields{
					"Topic": "config",
					"Key":   p.State.NeighborAddress,
					"Error": err})
		}
	}
}

func updateNeighbors(ctx context.Context, bgpServer *server.BgpServer, updated []bgpconfig.Neighbor) bool {
	for _, p := range updated {
		bgpServer.Log().Info("Update Peer",
			log.Fields{
				"Topic": "config",
				"Key":   p.State.NeighborAddress})
		if u, err := bgpServer.UpdatePeer(ctx, &api.UpdatePeerRequest{
			Peer: bgpconfig.NewPeerFromConfigStruct(&p),
		}); err != nil {
			bgpServer.Log().Warn("Failed to update Peer",
				log.Fields{
					"Topic": "config",
					"Key":   p.State.NeighborAddress,
					"Error": err})
		} else {
			return u.NeedsSoftResetIn
		}
	}
	return false
}

type bgpReconciler struct {
	bgpServer     *server.BgpServer
	bgpStarted    bool
	currentConfig *bgpconfig.BgpConfigSet
}

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

// InitialConfig applies initial configuration to a pristine gobgp instance. It
// can only be called once for an instance. Subsequent changes to the
// configuration can be applied using UpdateConfig. The BgpConfigSet can be
// obtained by calling ReadConfigFile. If graceful restart behavior is desired,
// pass true for isGracefulRestart. Otherwise, pass false.
func InitialConfig(ctx context.Context, applied *oc.NetworkInstance_Protocol_Bgp, bgpServer *server.BgpServer, newConfig *bgpconfig.BgpConfigSet, isGracefulRestart bool) (*bgpconfig.BgpConfigSet, error) {
	if err := bgpServer.StartBgp(ctx, &api.StartBgpRequest{
		Global: bgpconfig.NewGlobalFromConfigStruct(&newConfig.Global),
	}); err != nil {
		bgpServer.Log().Fatal("failed to set global config",
			log.Fields{"Topic": "config", "Error": err})
	} else {
		applied.GetOrCreateGlobal().SetAs(newConfig.Global.Config.As)
		applied.GetOrCreateGlobal().SetRouterId(newConfig.Global.Config.RouterId)
	}

	if newConfig.Zebra.Config.Enabled {
		tps := newConfig.Zebra.Config.RedistributeRouteTypeList
		l := make([]string, 0, len(tps))
		l = append(l, tps...)
		if err := bgpServer.EnableZebra(ctx, &api.EnableZebraRequest{
			Url:                  newConfig.Zebra.Config.Url,
			RouteTypes:           l,
			Version:              uint32(newConfig.Zebra.Config.Version),
			NexthopTriggerEnable: newConfig.Zebra.Config.NexthopTriggerEnable,
			NexthopTriggerDelay:  uint32(newConfig.Zebra.Config.NexthopTriggerDelay),
			MplsLabelRangeSize:   uint32(newConfig.Zebra.Config.MplsLabelRangeSize),
			SoftwareName:         newConfig.Zebra.Config.SoftwareName,
		}); err != nil {
			bgpServer.Log().Fatal("failed to set zebra config",
				log.Fields{"Topic": "config", "Error": err})
		}
	}

	if len(newConfig.Collector.Config.Url) > 0 {
		bgpServer.Log().Fatal("collector feature is not supported",
			log.Fields{"Topic": "config"})
	}

	for _, c := range newConfig.RpkiServers {
		if err := bgpServer.AddRpki(ctx, &api.AddRpkiRequest{
			Address:  c.Config.Address,
			Port:     c.Config.Port,
			Lifetime: c.Config.RecordLifetime,
		}); err != nil {
			bgpServer.Log().Fatal("failed to set rpki config",
				log.Fields{"Topic": "config", "Error": err})
		}
	}
	for _, c := range newConfig.BmpServers {
		if err := bgpServer.AddBmp(ctx, &api.AddBmpRequest{
			Address:           c.Config.Address,
			Port:              c.Config.Port,
			SysName:           c.Config.SysName,
			SysDescr:          c.Config.SysDescr,
			Policy:            api.AddBmpRequest_MonitoringPolicy(c.Config.RouteMonitoringPolicy.ToInt()),
			StatisticsTimeout: int32(c.Config.StatisticsTimeout),
		}); err != nil {
			bgpServer.Log().Fatal("failed to set bmp config",
				log.Fields{"Topic": "config", "Error": err})
		}
	}
	for _, vrf := range newConfig.Vrfs {
		rd, err := bgp.ParseRouteDistinguisher(vrf.Config.Rd)
		if err != nil {
			bgpServer.Log().Fatal("failed to load vrf rd config",
				log.Fields{"Topic": "config", "Error": err})
		}

		importRtList, err := marshalRouteTargets(vrf.Config.ImportRtList)
		if err != nil {
			bgpServer.Log().Fatal("failed to load vrf import rt config",
				log.Fields{"Topic": "config", "Error": err})
		}
		exportRtList, err := marshalRouteTargets(vrf.Config.ExportRtList)
		if err != nil {
			bgpServer.Log().Fatal("failed to load vrf export rt config",
				log.Fields{"Topic": "config", "Error": err})
		}

		a, err := apiutil.MarshalRD(rd)
		if err != nil {
			bgpServer.Log().Fatal("failed to set vrf config",
				log.Fields{"Topic": "config", "Error": err})
		}
		if err := bgpServer.AddVrf(ctx, &api.AddVrfRequest{
			Vrf: &api.Vrf{
				Name:     vrf.Config.Name,
				Rd:       a,
				Id:       uint32(vrf.Config.Id),
				ImportRt: importRtList,
				ExportRt: exportRtList,
			},
		}); err != nil {
			bgpServer.Log().Fatal("failed to set vrf config",
				log.Fields{"Topic": "config", "Error": err})
		}
	}
	for _, c := range newConfig.MrtDump {
		if len(c.Config.FileName) == 0 {
			continue
		}
		if err := bgpServer.EnableMrt(ctx, &api.EnableMrtRequest{
			Type:             api.EnableMrtRequest_DumpType(c.Config.DumpType.ToInt()),
			Filename:         c.Config.FileName,
			DumpInterval:     c.Config.DumpInterval,
			RotationInterval: c.Config.RotationInterval,
		}); err != nil {
			bgpServer.Log().Fatal("failed to set mrt config",
				log.Fields{"Topic": "config", "Error": err})
		}
	}
	p := bgpconfig.ConfigSetToRoutingPolicy(newConfig)
	rp, err := table.NewAPIRoutingPolicyFromConfigStruct(p)
	if err != nil {
		bgpServer.Log().Fatal("failed to update policy config",
			log.Fields{"Topic": "config", "Error": err})
	} else {
		bgpServer.SetPolicies(ctx, &api.SetPoliciesRequest{
			DefinedSets: rp.DefinedSets,
			Policies:    rp.Policies,
		})
	}

	assignGlobalpolicy(ctx, bgpServer, &newConfig.Global.ApplyPolicy.Config)

	added := newConfig.Neighbors
	addedPg := newConfig.PeerGroups
	if isGracefulRestart {
		for i, n := range added {
			if n.GracefulRestart.Config.Enabled {
				added[i].GracefulRestart.State.LocalRestarting = true
			}
		}
	}

	addPeerGroups(ctx, bgpServer, addedPg)
	addDynamicNeighbors(ctx, bgpServer, newConfig.DynamicNeighbors)
	addNeighbors(ctx, applied, bgpServer, added)
	return newConfig, nil
}

// UpdateConfig updates the configuration of a running gobgp instance.
// InitialConfig must have been called once before this can be called for
// subsequent changes to bgpconfig. The differences are that this call 1) does not
// hangle graceful restart and 2) requires a BgpConfigSet for the previous
// configuration so that it can compute the delta between it and the new
// bgpconfig. The new BgpConfigSet can be obtained using ReadConfigFile.
func UpdateConfig(ctx context.Context, applied *oc.NetworkInstance_Protocol_Bgp, bgpServer *server.BgpServer, c, newConfig *bgpconfig.BgpConfigSet) (*bgpconfig.BgpConfigSet, error) {
	addedPg, deletedPg, updatedPg := bgpconfig.UpdatePeerGroupConfig(bgpServer.Log(), c, newConfig)
	added, deleted, updated := bgpconfig.UpdateNeighborConfig(bgpServer.Log(), c, newConfig)
	updatePolicy := bgpconfig.CheckPolicyDifference(bgpServer.Log(), bgpconfig.ConfigSetToRoutingPolicy(c), bgpconfig.ConfigSetToRoutingPolicy(newConfig))

	if updatePolicy {
		bgpServer.Log().Info("policy config is update", log.Fields{"Topic": "config"})
		p := bgpconfig.ConfigSetToRoutingPolicy(newConfig)
		rp, err := table.NewAPIRoutingPolicyFromConfigStruct(p)
		if err != nil {
			bgpServer.Log().Warn("failed to update policy config",
				log.Fields{
					"Topic": "config",
					"Error": err})
		} else {
			bgpServer.SetPolicies(ctx, &api.SetPoliciesRequest{
				DefinedSets: rp.DefinedSets,
				Policies:    rp.Policies,
			})
		}
	}
	// global policy update
	if !newConfig.Global.ApplyPolicy.Config.Equal(&c.Global.ApplyPolicy.Config) {
		assignGlobalpolicy(ctx, bgpServer, &newConfig.Global.ApplyPolicy.Config)
		updatePolicy = true
	}

	addPeerGroups(ctx, bgpServer, addedPg)
	deletePeerGroups(ctx, bgpServer, deletedPg)
	needsSoftResetIn := updatePeerGroups(ctx, bgpServer, updatedPg)
	updatePolicy = updatePolicy || needsSoftResetIn
	addDynamicNeighbors(ctx, bgpServer, newConfig.DynamicNeighbors)
	addNeighbors(ctx, applied, bgpServer, added)
	deleteNeighbors(ctx, bgpServer, deleted)
	needsSoftResetIn = updateNeighbors(ctx, bgpServer, updated)
	updatePolicy = updatePolicy || needsSoftResetIn

	if updatePolicy {
		if err := bgpServer.ResetPeer(ctx, &api.ResetPeerRequest{
			Address:   "",
			Direction: api.ResetPeerRequest_IN,
			Soft:      true,
		}); err != nil {
			bgpServer.Log().Fatal("failed to update policy config",
				log.Fields{"Topic": "config", "Error": err})
		}
	}
	return newConfig, nil
}
