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

//go:build linux

// Package dplanerc contains gNMI task handlers.
package dplanerc

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/kernel"
	"github.com/openconfig/lemming/dataplane/protocol/lldp"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	log "github.com/golang/glog"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type ocInterface struct {
	name    string
	subintf uint32
}

type interfaceData struct {
	portID          uint64
	portNID         uint64
	hostifID        uint64
	hostifIfIndex   int
	hostifDevName   string
	rifID           uint64
	lagMembershipID uint64
	isAggregate     bool
	networkInstance string
}

type interfaceMap map[ocInterface]*interfaceData

func (d interfaceMap) findByIfIndex(ifIndex int) (ocInterface, *interfaceData) {
	for k, v := range d {
		if v.hostifIfIndex == ifIndex {
			return k, v
		}
	}
	return ocInterface{}, nil
}

func (d interfaceMap) findByPortID(portID uint64) (ocInterface, *interfaceData) {
	for k, v := range d {
		if v.portID == portID {
			return k, v
		}
	}
	return ocInterface{}, nil
}

type ocRoute struct {
	vrf    uint64
	prefix string
}

type routeData struct {
	nh    uint64 // OID of the NextHop
	isNHG bool
	nhg   map[uint64]map[uint64]uint64 // NHG_ID/NH_ID -> Member_ID
}

type routeMap map[ocRoute]*routeData

func (r routeMap) findRoute(ipPrefix string, vrfID uint64) *routeData {
	key := ocRoute{vrf: vrfID, prefix: ipPrefix}
	return r[key]
}

type protocolHanlder interface {
	Reconcile(context.Context, *oc.Root, *ygnmi.Client) error
}

// Reconciler handles config updates to the paths.
type Reconciler struct {
	c *ygnmi.Client
	// closers functions should all be invoked when the interface handler stops running.
	closers            []func()
	hostifClient       saipb.HostifClient
	portClient         saipb.PortClient
	switchClient       saipb.SwitchClient
	ifaceClient        saipb.RouterInterfaceClient
	neighborClient     saipb.NeighborClient
	routeClient        saipb.RouteClient
	nextHopClient      saipb.NextHopClient
	fwdClient          fwdpb.ForwardingClient
	nextHopGroupClient saipb.NextHopGroupClient
	lagClient          saipb.LagClient
	vrClient           saipb.VirtualRouterClient
	stateMu            sync.RWMutex
	lldp               protocolHanlder
	// state keeps track of the applied state of the device's interfaces so that we do not issue duplicate configuration commands to the device's interfaces.
	state           map[string]*oc.Interface
	switchID        uint64
	ifaceMgr        interfaceManager
	ocInterfaceData interfaceMap
	ocRouteData     routeMap
	cpuPortID       uint64
	contextID       string
	niDetail        map[string]*netInst
}

type netInst struct {
	vrOID uint64
}

type interfaceManager interface {
	SetHWAddr(name string, addr string) error
	SetState(name string, up bool) error
	ReplaceIP(name string, ip string, prefixLen int) error
	DeleteIP(name string, ip string, prefixLen int) error
	LinkSubscribe(ch chan<- netlink.LinkUpdate, done <-chan struct{}) error
	AddrSubscribe(ch chan<- netlink.AddrUpdate, done <-chan struct{}) error
	NeighSubscribe(ch chan<- netlink.NeighUpdate, done <-chan struct{}) error
	LinkList() ([]netlink.Link, error)
	LinkAdd(link netlink.Link) error
	LinkByName(name string) (netlink.Link, error)
	LinkByIndex(idx int) (netlink.Link, error)
	LinkSetDown(link netlink.Link) error
	LinkSetUp(link netlink.Link) error
	LinkSetMaster(member netlink.Link, link netlink.Link) error
	LinkSetNoMaster(link netlink.Link) error
	LinkModify(link netlink.Link) error
}

// New creates a new interface handler.
func New(conn grpc.ClientConnInterface, switchID, cpuPortID uint64, contextID string) *Reconciler {
	r := &Reconciler{
		state:              map[string]*oc.Interface{},
		ifaceMgr:           &kernel.Interfaces{},
		switchID:           switchID,
		cpuPortID:          cpuPortID,
		contextID:          contextID,
		ocInterfaceData:    interfaceMap{},
		ocRouteData:        routeMap{},
		hostifClient:       saipb.NewHostifClient(conn),
		portClient:         saipb.NewPortClient(conn),
		switchClient:       saipb.NewSwitchClient(conn),
		ifaceClient:        saipb.NewRouterInterfaceClient(conn),
		neighborClient:     saipb.NewNeighborClient(conn),
		routeClient:        saipb.NewRouteClient(conn),
		nextHopClient:      saipb.NewNextHopClient(conn),
		nextHopGroupClient: saipb.NewNextHopGroupClient(conn),
		fwdClient:          fwdpb.NewForwardingClient(conn),
		lagClient:          saipb.NewLagClient(conn),
		vrClient:           saipb.NewVirtualRouterClient(conn),
		lldp:               lldp.New(),
		niDetail:           map[string]*netInst{},
	}
	return r
}

// Start starts running the handler, watching the cache and the kernel interfaces.
func (ni *Reconciler) StartInterface(ctx context.Context, client *ygnmi.Client) error {
	log.Info("starting interface handler")
	b := &ocpath.Batch{}
	ni.c = client

	vrID, err := ni.switchClient.GetSwitchAttribute(ctx, &saipb.GetSwitchAttributeRequest{
		Oid:      ni.switchID,
		AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_DEFAULT_VIRTUAL_ROUTER_ID},
	})
	if err != nil {
		return err
	}

	ni.niDetail[fakedevice.DefaultNetworkInstance] = &netInst{
		vrOID: vrID.GetAttr().GetDefaultVirtualRouterId(),
	}

	if err := ni.setupPorts(ctx); err != nil {
		return fmt.Errorf("failed to setup ports: %v", err)
	}

	b.AddPaths(
		ocpath.Root().InterfaceAny().Name().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Ethernet().MacAddress().Config().PathStruct(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Enabled().Config().PathStruct(), // TODO: Support the parent interface config/enabled controling the subinterface state.
		ocpath.Root().InterfaceAny().SubinterfaceAny().Ipv4().AddressAny().Ip().Config().PathStruct(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Ipv4().AddressAny().PrefixLength().Config().PathStruct(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Ipv6().AddressAny().Ip().Config().PathStruct(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Ipv6().AddressAny().PrefixLength().Config().PathStruct(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Vlan().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Aggregation().LagType().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Ethernet().AggregateId().Config().PathStruct(),
		ocpath.Root().Lldp().Enabled().Config().PathStruct(),
		ocpath.Root().Lldp().InterfaceAny().Config().PathStruct(),
		ocpath.Root().NetworkInstanceAny().InterfaceAny().Config().PathStruct(),
		ocpath.Root().NetworkInstanceAny().Name().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Type().Config().PathStruct(),
	)
	cancelCtx, cancelFn := context.WithCancel(ctx)

	watcher := ygnmi.Watch(cancelCtx, ni.c, b.Config(), func(val *ygnmi.Value[*oc.Root]) error {
		log.V(2).Info("reconciling interfaces")
		root, ok := val.Val()
		if !ok {
			return ygnmi.Continue
		}
		for _, netInst := range root.NetworkInstance {
			ni.reconcileNI(netInst)
		}
		if root.Interface != nil {
			for _, i := range root.Interface {
				ni.reconcile(cancelCtx, i)
			}
		}
		if root.Lldp.Interface != nil {
			ni.reconcileLldp(cancelCtx, root)
		}

		return ygnmi.Continue
	})
	linkDoneCh := make(chan struct{})
	linkUpdateCh := make(chan netlink.LinkUpdate)
	addrDoneCh := make(chan struct{})
	addrUpdateCh := make(chan netlink.AddrUpdate)
	neighDoneCh := make(chan struct{})
	neighUpdateCh := make(chan netlink.NeighUpdate)
	ni.closers = append(ni.closers, func() {
		close(linkDoneCh)
		close(addrDoneCh)
		close(neighDoneCh)
	}, cancelFn)

	if err := ni.ifaceMgr.LinkSubscribe(linkUpdateCh, linkDoneCh); err != nil {
		return fmt.Errorf("failed to sub to link: %v", err)
	}
	if err := ni.ifaceMgr.AddrSubscribe(addrUpdateCh, addrDoneCh); err != nil {
		return fmt.Errorf("failed to sub to addr: %v", err)
	}
	if err := ni.ifaceMgr.NeighSubscribe(neighUpdateCh, addrDoneCh); err != nil {
		return fmt.Errorf("failed to sub to neighbor: %v", err)
	}
	notifClient, err := ni.switchClient.PortStateChangeNotification(cancelCtx, &saipb.PortStateChangeNotificationRequest{})
	if err != nil {
		return err
	}
	go func() {
		for {
			n, err := notifClient.Recv()
			if err != nil {
				return
			}
			ni.handleDataplaneEvent(ctx, n)
		}
	}()

	go func() {
		for {
			select {
			case up := <-linkUpdateCh:
				ni.handleLinkUpdate(ctx, &up)
			case up := <-addrUpdateCh:
				ni.handleAddrUpdate(ctx, &up)
			case up := <-neighUpdateCh:
				ni.handleNeighborUpdate(ctx, &up)
			}
		}
	}()

	go func() {
		// TODO: handle error
		if _, err := watcher.Await(); err != nil {
			log.Warningf("interface watch err: %v", err)
		}
	}()

	ni.startCounterUpdates(ctx)

	return nil
}

// Stop stops all watchers.
func (ni *Reconciler) Stop(context.Context) error {
	// TODO: prevent stopping more than once.
	for _, closeFn := range ni.closers {
		closeFn()
	}
	return nil
}

// startCounterUpdates starts a goroutine for updating counters for configured
// interfaces.
func (ni *Reconciler) startCounterUpdates(ctx context.Context) {
	tick := time.NewTicker(time.Second)
	ni.closers = append(ni.closers, tick.Stop)
	go func() {
		// Design comment:
		// This polling can be eliminated if either the forwarding
		// service supported streaming the counters, or if somehow the
		// gnmi cache were able to forward queries to prompt the data
		// producer to populate the leaf.
		//
		// However, given counters are likely frequently-updated values
		// anyways, it may be fine for counter values to be polled.
		for range tick.C {
			ni.stateMu.RLock()
			var intfNames []ocInterface
			for intfName := range ni.state {
				// TODO(wenbli): Support interface state deletion when interface is deleted.
				intfNames = append(intfNames, ocInterface{name: intfName})
			}
			ni.stateMu.RUnlock()
			for _, intfName := range intfNames {
				intf := ni.ocInterfaceData[intfName]
				if intf == nil {
					continue
				}
				stats, err := ni.portClient.GetPortStats(ctx, &saipb.GetPortStatsRequest{
					Oid: intf.portID,
					CounterIds: []saipb.PortStat{
						saipb.PortStat_PORT_STAT_IF_IN_UCAST_PKTS,      // 0
						saipb.PortStat_PORT_STAT_IF_IN_NON_UCAST_PKTS,  // 1
						saipb.PortStat_PORT_STAT_IF_OUT_UCAST_PKTS,     // 2
						saipb.PortStat_PORT_STAT_IF_OUT_NON_UCAST_PKTS, // 3
						saipb.PortStat_PORT_STAT_IF_IN_OCTETS,          // 4
						saipb.PortStat_PORT_STAT_IF_OUT_OCTETS,         // 5
						saipb.PortStat_PORT_STAT_IF_OUT_MULTICAST_PKTS, // 6
						saipb.PortStat_PORT_STAT_IF_OUT_BROADCAST_PKTS, // 7
					},
				})
				log.V(2).Infof("querying counters for interface %q, got %v", intfName, stats)
				if err != nil {
					log.Errorf("interface handler: could not retrieve counter for interface %q", intfName)
					continue
				}
				sb := &ygnmi.SetBatch{}
				gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intfName.name).Counters().InPkts().State(), stats.Values[0]+stats.Values[1])
				gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intfName.name).Counters().OutPkts().State(), stats.Values[2]+stats.Values[2])
				gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intfName.name).Counters().InUnicastPkts().State(), stats.Values[0])
				gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intfName.name).Counters().OutUnicastPkts().State(), stats.Values[2])
				gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intfName.name).Counters().InOctets().State(), stats.Values[4])
				gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intfName.name).Counters().OutOctets().State(), stats.Values[5])

				if _, err := sb.Set(ctx, ni.c); err != nil {
					log.Errorf("interface handler: %v", err)
				}
			}
		}
	}()
}

func intfRefToDevName(intf ocInterface) string {
	if intf.subintf == 0 {
		return intf.name
	}
	return fmt.Sprintf("%s.%d", intf.name, intf.subintf)
}

func (ni *Reconciler) createLAG(ctx context.Context, intf ocInterface, lagType oc.E_IfAggregate_AggregationType) error {
	bond := netlink.NewLinkBond(netlink.NewLinkAttrs())
	bond.Name = intfRefToDevName(intf)
	bond.Mode = netlink.BOND_MODE_BALANCE_XOR
	if err := ni.ifaceMgr.LinkAdd(bond); err != nil {
		return fmt.Errorf("failed to create kernel lag interface: %v", err)
	}
	lagResp, err := ni.lagClient.CreateLag(ctx, &saipb.CreateLagRequest{
		Switch: ni.switchID,
	})
	if err != nil {
		return fmt.Errorf("failed to create router interface %q: %v", intf.name, err)
	}
	l, err := ni.ifaceMgr.LinkByName(intfRefToDevName(intf))
	if err != nil {
		return fmt.Errorf("failed to get bond intf %q: %v", intf.name, err)
	}
	rifResp, err := ni.ifaceClient.CreateRouterInterface(ctx, &saipb.CreateRouterInterfaceRequest{
		Switch:          ni.switchID,
		Type:            saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT.Enum(),
		PortId:          &lagResp.Oid,
		VirtualRouterId: proto.Uint64(0),
		SrcMacAddress:   l.Attrs().HardwareAddr,
	})
	if err != nil {
		return fmt.Errorf("failed to create router interface %q: %v", intfRefToDevName(intf), err)
	}
	aggData := &interfaceData{
		portID:          lagResp.Oid,
		rifID:           rifResp.Oid,
		hostifIfIndex:   bond.Index,
		hostifDevName:   intfRefToDevName(intf),
		isAggregate:     true,
		networkInstance: fakedevice.DefaultNetworkInstance,
	}

	ni.getOrCreateInterface(intf.name).GetOrCreateEthernet().SetHwMacAddress(l.Attrs().HardwareAddr.String())
	ni.getOrCreateInterface(intf.name).GetOrCreateEthernet().SetMacAddress(l.Attrs().HardwareAddr.String())
	ni.getOrCreateInterface(intf.name).GetOrCreateAggregation().SetLagType(lagType)

	sb := &ygnmi.SetBatch{}
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Ethernet().HwMacAddress().State(), l.Attrs().HardwareAddr.String())
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Ethernet().MacAddress().State(), l.Attrs().HardwareAddr.String())
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Aggregation().LagType().State(), lagType)
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Enabled().State(), true)
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).AdminStatus().State(), oc.Interface_AdminStatus_UP)
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).OperStatus().State(), oc.Interface_OperStatus_UP)
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Type().State(), oc.IETFInterfaces_InterfaceType_ieee8023adLag)
	if _, err := sb.Set(ctx, ni.c); err != nil {
		return fmt.Errorf("failed to update agg state: %v", err)
	}

	ni.ocInterfaceData[intf] = aggData
	return nil
}

const ifaceDownRetries = 3

func (ni *Reconciler) addLAGMember(ctx context.Context, intf ocInterface, memberData *interfaceData, aggID string) error {
	agg, ok := ni.ocInterfaceData[ocInterface{name: aggID}]
	if !ok {
		return fmt.Errorf("unknown aggregate id %q", aggID)
	}
	bondLink, err := ni.ifaceMgr.LinkByIndex(agg.hostifIfIndex)
	if err != nil {
		return fmt.Errorf("failed to find bond link: %v", err)
	}
	memberLink, err := ni.ifaceMgr.LinkByIndex(memberData.hostifIfIndex)
	if err != nil {
		return fmt.Errorf("failed to find member link: %v", err)
	}

	// Can only add links to a bond interface when it's down.
	if memberLink.Attrs().OperState != netlink.OperDown {
		log.Infof("aggregate link %v oper status %v, setting to down", intf.name, bondLink.Attrs().OperState)
		if err := ni.ifaceMgr.LinkSetDown(memberLink); err != nil {
			log.Warningf("failed to set link %v down: %v", intf.name, err)
		}
		for i := 0; i < ifaceDownRetries; i++ { // Ensure link is down before continuing.
			memberLink, err := ni.ifaceMgr.LinkByIndex(memberData.hostifIfIndex)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if memberLink.Attrs().OperState == netlink.OperDown {
				break
			}
			time.Sleep(time.Second)
		}

		defer func() {
			if err := ni.ifaceMgr.LinkSetUp(memberLink); err != nil {
				log.Warningf("failed to set link %v up: %v", intf.name, err)
			}
		}()
	}

	if err := ni.ifaceMgr.LinkSetMaster(memberLink, bondLink); err != nil {
		return fmt.Errorf("failed to add bond member: %v", err)
	}
	resp, err := ni.lagClient.CreateLagMember(ctx, &saipb.CreateLagMemberRequest{
		Switch: ni.switchID,
		LagId:  proto.Uint64(agg.portID),
		PortId: proto.Uint64(memberData.portID),
	})
	if err != nil {
		return fmt.Errorf("failed to create lag member: %v", err)
	}
	ni.getOrCreateInterface(intf.name).GetOrCreateEthernet().AggregateId = &aggID
	ni.getOrCreateInterface(aggID).GetAggregation().Member = append(ni.getOrCreateInterface(aggID).GetAggregation().Member, intf.name)
	sb := &ygnmi.SetBatch{}
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Ethernet().AggregateId().State(), aggID)
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(aggID).Aggregation().Member().State(), ni.getOrCreateInterface(aggID).GetAggregation().Member)
	if _, err := sb.Set(ctx, ni.c); err != nil {
		return fmt.Errorf("failed to update agg state: %v", err)
	}

	memberData.lagMembershipID = resp.Oid
	return nil
}

func (ni *Reconciler) removeLAGMember(ctx context.Context, intf ocInterface, memberData *interfaceData) error {
	memberLink, err := ni.ifaceMgr.LinkByIndex(memberData.hostifIfIndex)
	if err != nil {
		return fmt.Errorf("failed to find member link: %v", err)
	}
	if err := ni.ifaceMgr.LinkSetNoMaster(memberLink); err != nil {
		return fmt.Errorf("failed to remove bond: %v", err)
	}
	_, err = ni.lagClient.RemoveLagMember(ctx, &saipb.RemoveLagMemberRequest{
		Oid: memberData.lagMembershipID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove lag member: %v", err)
	}
	sb := &ygnmi.SetBatch{}
	aggID := ni.getOrCreateInterface(intf.name).GetOrCreateEthernet().GetAggregateId()
	ni.getOrCreateInterface(intf.name).GetOrCreateEthernet().AggregateId = nil
	idx := -1
	for i, member := range ni.getOrCreateInterface(aggID).GetOrCreateAggregation().Member {
		if member == intf.name {
			idx = i
			break
		}
	}
	if isAggregateMember := idx != -1; isAggregateMember {
		ni.getOrCreateInterface(aggID).GetOrCreateAggregation().Member = append(ni.getOrCreateInterface(aggID).GetOrCreateAggregation().Member[:idx],
			ni.getOrCreateInterface(aggID).GetOrCreateAggregation().Member[idx+1:]...)
	}
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(aggID).Aggregation().Member().State(), ni.getOrCreateInterface(aggID).GetOrCreateAggregation().Member)
	gnmiclient.BatchDelete(sb, ocpath.Root().Interface(intf.name).Ethernet().AggregateId().State())
	if _, err := sb.Set(ctx, ni.c); err != nil {
		return fmt.Errorf("failed to update agg state: %v", err)
	}
	memberData.lagMembershipID = 0
	return nil
}

func (ni *Reconciler) setMinLinks(intf ocInterface, data *interfaceData, minLink uint16) error {
	link, err := ni.ifaceMgr.LinkByIndex(data.hostifIfIndex)
	if err != nil {
		return fmt.Errorf("failed to find link: %v", err)
	}
	bond, ok := link.(*netlink.Bond)
	if !ok {
		return fmt.Errorf("link %s is not a bond interface", intf.name)
	}
	bond.MinLinks = int(minLink)
	if err := ni.ifaceMgr.LinkModify(bond); err != nil {
		return fmt.Errorf("failed to modify link %s: %v", intf.name, err)
	}
	ni.getOrCreateInterface(intf.name).GetOrCreateAggregation().SetMinLinks(minLink)
	return nil
}

// reconcileLldp compares the LLDP config with state and modifies state to match config.
func (ni *Reconciler) reconcileLldp(ctx context.Context, intent *oc.Root) {
	if err := ni.lldp.Reconcile(ctx, intent, ni.c); err != nil {
		log.Warningf("error found LLDP reconciliation: %v", err)
	}
}

// reconcile compares the interface config with state and modifies state to match config.
func (ni *Reconciler) reconcile(ctx context.Context, config *oc.Interface) {
	ni.stateMu.RLock()
	defer ni.stateMu.RUnlock()

	intf := ocInterface{name: config.GetName(), subintf: 0}
	state := ni.getOrCreateInterface(config.GetName())
	data := ni.ocInterfaceData[intf]

	ni.reconcileEthernet(ctx, config, state, intf, data)

	sb := &ygnmi.SetBatch{}
	if config.GetEnabled() != state.GetEnabled() {
		log.Infof("reconciling config enabled on intf: %v", intf)
		adminStatus := oc.Interface_AdminStatus_DOWN
		if config.GetEnabled() {
			adminStatus = oc.Interface_AdminStatus_UP
		}
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Enabled().State(), config.GetEnabled())
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).AdminStatus().State(), adminStatus)
	}
	if _, err := sb.Set(ctx, ni.c); err != nil {
		log.Warningf("failed to set link status: %v", err)
	}

	if config.GetOrCreateAggregation().GetMinLinks() != state.GetOrCreateAggregation().GetMinLinks() {
		if err := ni.setMinLinks(intf, data, config.GetOrCreateAggregation().GetMinLinks()); err != nil {
			log.Warningf("failed to set min links: %v", err)
		}
	}
	ni.reconcileSubIntf(ctx, config, state)
	ni.reconcileIPs(config, state)
}

func (ni *Reconciler) reconcileEthernet(ctx context.Context, config, state *oc.Interface, intf ocInterface, data *interfaceData) {
	if data == nil {
		return
	}
	if config.GetOrCreateEthernet().GetAggregateId() != state.GetOrCreateEthernet().GetAggregateId() {
		log.Infof("reconciling lag member intf %v: config agg id %v, state agg id %v", intf.name, config.GetEthernet().GetAggregateId(), state.GetEthernet().GetAggregateId())
		if data.lagMembershipID != 0 {
			log.Infof("intf %v has existing lab membership lag membership id %d", intf.name, data.lagMembershipID)
			if err := ni.removeLAGMember(ctx, intf, data); err != nil {
				log.Warningf("intf %v failed to remove lag member: %v", intf.name, err)
			}
		}
		if config.GetEthernet().GetAggregateId() != "" {
			log.Infof("intf %v adding to agg id %v ", intf.name, config.GetEthernet().GetAggregateId())
			if err := ni.addLAGMember(ctx, intf, data, config.GetEthernet().GetAggregateId()); err != nil {
				log.Warningf("intf %v failed to add lag member %v: %v", intf.name, config.GetEthernet().GetAggregateId(), err)
			}
		}
	}

	if config.GetOrCreateEthernet().MacAddress != nil {
		if config.GetEthernet().GetMacAddress() != state.GetEthernet().GetMacAddress() {
			log.V(1).Infof("setting interface %s hw-addr %q", data.hostifDevName, config.GetEthernet().GetMacAddress())
			if err := ni.ifaceMgr.SetHWAddr(config.GetName(), config.GetEthernet().GetMacAddress()); err != nil {
				log.Warningf("Failed to set mac address of port: %v", err)
			}
		}
	} else {
		// Deleting the configured MAC address means it should be the system-assigned MAC address, as detailed in the OpenConfig schema.
		// https://openconfig.net/projects/models/schemadocs/yangdoc/openconfig-interfaces.html#interfaces-interface-ethernet-state-mac-address
		if state.GetEthernet().GetHwMacAddress() != state.GetEthernet().GetMacAddress() {
			log.V(1).Infof("resetting interface %s hw-addr %q", data.hostifDevName, state.GetEthernet().GetHwMacAddress())
			if err := ni.ifaceMgr.SetHWAddr(config.GetName(), state.GetEthernet().GetHwMacAddress()); err != nil {
				log.Warningf("Failed to set mac address of port: %v", err)
			}
		}
	}
}

func (ni *Reconciler) reconcileSubIntf(ctx context.Context, config, state *oc.Interface) {
	for idx, subintf := range config.Subinterface {
		intfRef := ocInterface{
			name:    config.GetName(),
			subintf: idx,
		}
		data, ok := ni.ocInterfaceData[intfRef]
		// We have a new interface.
		if !ok {
			switch {
			case config.GetType() == oc.IETFInterfaces_InterfaceType_ieee8023adLag:
				log.Infof("creating new lag interface: %v", intfRefToDevName(intfRef))
				if err := ni.createLAG(ctx, intfRef, config.GetAggregation().GetLagType()); err != nil {
					log.Warningf("failed to create lag: %v", err)
				}
			case idx != 0 && subintf.Vlan != nil: // TODO: add support for vlan on the subintf 0.
				log.Infof("creating new vlan intf: %v", intfRefToDevName(intfRef))
				ni.createVLANSubIntf(ctx, intfRef, config)
			default:
				log.Warningf("new interface %+v, can't be created", intfRef)
				return
			}
			data = ni.ocInterfaceData[intfRef]
		}

		enabled := config.GetSubinterface(idx).GetEnabled() && config.GetEnabled()
		if data.hostifID != 0 {
			_, err := ni.hostifClient.SetHostifAttribute(ctx, &saipb.SetHostifAttributeRequest{
				Oid:        data.hostifID,
				OperStatus: proto.Bool(enabled),
			})
			if err != nil {
				log.Warningf("Failed to set oper status of hostif: %v", err)
			}
		}
		if !data.isAggregate {
			_, err := ni.portClient.SetPortAttribute(ctx, &saipb.SetPortAttributeRequest{
				Oid:        data.portID,
				AdminState: proto.Bool(enabled),
			})
			if err != nil {
				log.Warningf("Failed to set admin state of port: %v", err)
			}
		}
		if err := ni.ifaceMgr.SetState(data.hostifDevName, config.GetSubinterface(intfRef.subintf).GetEnabled()); err != nil {
			log.Warningf("Failed to set admin state of hostif: %v", err)
		}
		sb := &ygnmi.SetBatch{}

		adminStatus := oc.Interface_AdminStatus_DOWN
		if enabled {
			adminStatus = oc.Interface_AdminStatus_UP
		}
		// TODO: Right now treating subinterface 0 and interface as the same.
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intfRef.name).Subinterface(intfRef.subintf).Enabled().State(), enabled)
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intfRef.name).Subinterface(intfRef.subintf).AdminStatus().State(), adminStatus)
		if _, err := sb.Set(ctx, ni.c); err != nil {
			log.Warningf("failed to set link status: %v", err)
		}
	}
}

func (ni *Reconciler) createVLANSubIntf(ctx context.Context, intfRef ocInterface, config *oc.Interface) {
	rootPort := ni.ocInterfaceData[ocInterface{name: intfRef.name, subintf: 0}]

	vlanIntf := &netlink.Vlan{
		LinkAttrs: netlink.LinkAttrs{
			Name:        intfRefToDevName(intfRef),
			ParentIndex: rootPort.hostifIfIndex,
		},
		VlanId:       int(config.GetSubinterface(intfRef.subintf).GetVlan().GetMatch().GetSingleTagged().GetVlanId()),
		VlanProtocol: netlink.VLAN_PROTOCOL_8021Q,
	}
	if err := ni.ifaceMgr.LinkAdd(vlanIntf); err != nil {
		log.Warningf("failed to add vlan intf: %v", err)
		return
	}
	rootPortAttr, err := ni.ifaceClient.GetRouterInterfaceAttribute(ctx, &saipb.GetRouterInterfaceAttributeRequest{
		Oid:      rootPort.rifID,
		AttrType: []saipb.RouterInterfaceAttr{saipb.RouterInterfaceAttr_ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS},
	})
	if err != nil {
		log.Warningf("failed to get root port mac, %v", err)
		return
	}

	rifResp, err := ni.ifaceClient.CreateRouterInterface(ctx, &saipb.CreateRouterInterfaceRequest{
		Switch:          ni.switchID,
		Type:            saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_SUB_PORT.Enum(),
		PortId:          &rootPort.portID,
		OuterVlanId:     proto.Uint32(uint32(config.GetSubinterface(intfRef.subintf).GetVlan().GetMatch().GetSingleTagged().GetVlanId())),
		VirtualRouterId: proto.Uint64(0),
		SrcMacAddress:   rootPortAttr.GetAttr().SrcMacAddress,
	})
	if err != nil {
		log.Warningf("failed to add vlan intf: %v", err)
		return
	}

	ni.ocInterfaceData[intfRef] = &interfaceData{
		rifID:           rifResp.GetOid(),
		hostifIfIndex:   vlanIntf.Index,
		hostifDevName:   intfRefToDevName(intfRef),
		networkInstance: fakedevice.DefaultNetworkInstance,
	}
	sb := &ygnmi.SetBatch{}
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intfRef.name).Subinterface(intfRef.subintf).Vlan().Match().SingleTagged().VlanId().State(), uint16(vlanIntf.VlanId))
	if _, err := sb.Set(ctx, ni.c); err != nil {
		log.Warningf("failed to set link status: %v", err)
	}
}

func (ni *Reconciler) reconcileIPs(config, state *oc.Interface) {
	type prefixPair struct {
		cfgIP, stateIP *string
		cfgPL, statePL *uint8
	}

	for idx := range config.Subinterface {
		intfRef := ocInterface{
			name:    config.GetName(),
			subintf: idx,
		}
		data, ok := ni.ocInterfaceData[intfRef]
		if !ok {
			log.Infof("skipping ip reconcilation for %+v", intfRef)
			continue
		}

		// Get all state IPs and their corresponding config IPs (if they exist).
		var interfacePairs []*prefixPair
		for _, addr := range state.GetOrCreateSubinterface(idx).GetOrCreateIpv4().Address {
			pair := &prefixPair{
				stateIP: addr.Ip,
				statePL: addr.PrefixLength,
			}
			if pairAddr := config.GetSubinterface(idx).GetIpv4().GetAddress(addr.GetIp()); pairAddr != nil {
				pair.cfgIP = pairAddr.Ip
				pair.cfgPL = pairAddr.PrefixLength
			}
			interfacePairs = append(interfacePairs, pair)
		}
		for _, addr := range state.GetOrCreateSubinterface(idx).GetOrCreateIpv6().Address {
			pair := &prefixPair{
				stateIP: addr.Ip,
				statePL: addr.PrefixLength,
			}
			if pairAddr := config.GetSubinterface(idx).GetIpv6().GetAddress(addr.GetIp()); pairAddr != nil {
				pair.cfgIP = pairAddr.Ip
				pair.cfgPL = pairAddr.PrefixLength
			}
			interfacePairs = append(interfacePairs, pair)
		}

		for _, addr := range config.GetOrCreateSubinterface(idx).GetOrCreateIpv4().Address {
			pair := &prefixPair{
				cfgIP: addr.Ip,
				cfgPL: addr.PrefixLength,
			}
			if pairAddr := state.GetSubinterface(idx).GetIpv4().GetAddress(addr.GetIp()); pairAddr != nil {
				pair.stateIP = pairAddr.Ip
				pair.statePL = pairAddr.PrefixLength
			}
			interfacePairs = append(interfacePairs, pair)
		}
		for _, addr := range config.GetOrCreateSubinterface(idx).GetOrCreateIpv6().Address {
			pair := &prefixPair{
				cfgIP: addr.Ip,
				cfgPL: addr.PrefixLength,
			}
			if pairAddr := state.GetSubinterface(idx).GetIpv6().GetAddress(addr.GetIp()); pairAddr != nil {
				pair.stateIP = pairAddr.Ip
				pair.statePL = pairAddr.PrefixLength
			}
			interfacePairs = append(interfacePairs, pair)
		}

		for _, pair := range interfacePairs {
			// If an IP exists in state, but not in config, remove the IP.
			if (pair.stateIP != nil && pair.statePL != nil) && (pair.cfgIP == nil && pair.cfgPL == nil) {
				log.V(1).Infof("Delete Config IP: %v, Config PL: %v. State IP: %v, State PL: %v", pair.cfgIP, pair.cfgPL, *pair.stateIP, *pair.statePL)
				log.V(2).Infof("deleting interface %s ip %s/%d", data.hostifDevName, *pair.stateIP, *pair.statePL)
				if err := ni.ifaceMgr.DeleteIP(data.hostifDevName, *pair.stateIP, int(*pair.statePL)); err != nil {
					log.Warningf("Failed to set ip address of port: %v", err)
				}
			}
			// If an IP exists in config, but not in state (or state is different) add the IP.
			if (pair.cfgIP != nil && pair.cfgPL != nil) && (pair.stateIP == nil || *pair.statePL != *pair.cfgPL) {
				log.V(1).Infof("Set Config IP: %v, Config PL: %v. State IP: %v, State PL: %v", *pair.cfgIP, *pair.cfgPL, pair.stateIP, pair.statePL)
				log.V(2).Infof("setting interface %s ip %s/%d", data.hostifDevName, *pair.cfgIP, *pair.cfgPL)
				if err := ni.ifaceMgr.ReplaceIP(data.hostifDevName, *pair.cfgIP, int(*pair.cfgPL)); err != nil {
					log.Warningf("Failed to set ip address of port: %v", err)
				}
			}
		}
	}
}

func (rec *Reconciler) reconcileNI(config *oc.NetworkInstance) {
	if _, ok := rec.niDetail[config.GetName()]; !ok {
		resp, err := rec.vrClient.CreateVirtualRouter(context.Background(), &saipb.CreateVirtualRouterRequest{
			Switch: rec.switchID,
		})
		if err != nil {
			log.Warningf("failed to create virtual router for %q: %v", config.GetName(), err)
		} else {
			log.Infof("created virtual router for %q: %v", config.GetName(), resp.GetOid())
			rec.niDetail[config.GetName()] = &netInst{vrOID: resp.GetOid()}
		}
	}

	for _, intf := range config.Interface {
		if intf.Interface == nil || intf.Subinterface == nil {
			continue
		}
		intfRef := ocInterface{
			name:    intf.GetInterface(),
			subintf: intf.GetSubinterface(),
		}
		if data, ok := rec.ocInterfaceData[intfRef]; ok {
			if data.networkInstance != config.GetName() {
				log.Infof("rif %s/%d moving network instance from %q to %q", intfRef.name, intfRef.subintf, data.networkInstance, config.GetName())

				attr, err := rec.ifaceClient.GetRouterInterfaceAttribute(context.Background(), &saipb.GetRouterInterfaceAttributeRequest{
					Oid: data.rifID,
					AttrType: []saipb.RouterInterfaceAttr{
						saipb.RouterInterfaceAttr_ROUTER_INTERFACE_ATTR_TYPE,
						saipb.RouterInterfaceAttr_ROUTER_INTERFACE_ATTR_PORT_ID,
						saipb.RouterInterfaceAttr_ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS,
					},
				})
				if err != nil {
					log.Warningf("failed to get rif attrs %s/%d %v", intfRef.name, intfRef.subintf, err)
					continue
				}
				if attr.GetAttr().GetType() == saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_SUB_PORT {
					attr, err = rec.ifaceClient.GetRouterInterfaceAttribute(context.Background(), &saipb.GetRouterInterfaceAttributeRequest{
						Oid: data.rifID,
						AttrType: []saipb.RouterInterfaceAttr{
							saipb.RouterInterfaceAttr_ROUTER_INTERFACE_ATTR_TYPE,
							saipb.RouterInterfaceAttr_ROUTER_INTERFACE_ATTR_PORT_ID,
							saipb.RouterInterfaceAttr_ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS,
							saipb.RouterInterfaceAttr_ROUTER_INTERFACE_ATTR_OUTER_VLAN_ID,
						},
					})
					if err != nil {
						log.Warningf("failed to get rif attrs %s/%d %v", intfRef.name, intfRef.subintf, err)
						continue
					}
				}

				_, err = rec.ifaceClient.RemoveRouterInterface(context.Background(), &saipb.RemoveRouterInterfaceRequest{
					Oid: data.rifID,
				})
				if err != nil {
					log.Warningf("failed to remove rif %s/%d %v", intfRef.name, intfRef.subintf, err)
				}
				rifResp, err := rec.ifaceClient.CreateRouterInterface(context.Background(), &saipb.CreateRouterInterfaceRequest{
					Type:            attr.GetAttr().Type,
					PortId:          attr.GetAttr().PortId,
					SrcMacAddress:   attr.GetAttr().SrcMacAddress,
					OuterVlanId:     attr.GetAttr().OuterVlanId,
					VirtualRouterId: &rec.niDetail[config.GetName()].vrOID,
				})
				if err != nil {
					log.Warningf("failed to create rif %s/%d %v", intfRef.name, intfRef.subintf, err)
				}
				data.rifID = rifResp.Oid
			}
		}
	}
}

// getOrCreateInterface returns the state interface from the cache.
func (ni *Reconciler) getOrCreateInterface(iface string) *oc.Interface {
	if _, ok := ni.state[iface]; !ok {
		ni.state[iface] = &oc.Interface{
			Name: &iface,
		}
	}
	return ni.state[iface]
}

func (ni *Reconciler) handleDataplaneEvent(ctx context.Context, resp *saipb.PortStateChangeNotificationResponse) {
	for _, event := range resp.Data {
		log.V(1).Infof("handling dataplane update on: %q", event.String())
		intf, data := ni.ocInterfaceData.findByPortID(event.GetPortId())
		if data == nil {
			return
		}
		operStatus := oc.Interface_OperStatus_UNKNOWN
		switch event.PortState {
		case saipb.PortOperStatus_PORT_OPER_STATUS_DOWN:
			operStatus = oc.Interface_OperStatus_DOWN
		case saipb.PortOperStatus_PORT_OPER_STATUS_UP:
			operStatus = oc.Interface_OperStatus_UP
		}

		sb := &ygnmi.SetBatch{}
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).OperStatus().State(), operStatus)

		if _, err := sb.Set(ctx, ni.c); err != nil {
			log.Warningf("failed to set link status: %v", err)
		}
	}
}

// handleLinkUpdate modifies the state based on changes to link state.
// This is the callback from netlink.
func (ni *Reconciler) handleLinkUpdate(ctx context.Context, lu *netlink.LinkUpdate) {
	ni.stateMu.Lock()
	defer ni.stateMu.Unlock()

	log.V(1).Infof("handling link update for %s", lu.Attrs().Name)

	intf, data := ni.ocInterfaceData.findByIfIndex(lu.Attrs().Index)
	if data == nil {
		return
	}

	iface := ni.getOrCreateInterface(intf.name)
	_, err := ni.ifaceClient.SetRouterInterfaceAttribute(ctx, &saipb.SetRouterInterfaceAttributeRequest{
		Oid:           data.rifID,
		SrcMacAddress: lu.Attrs().HardwareAddr,
	})
	if err != nil {
		log.Warningf("failed to update src mac: %v", err)
	}
	iface.GetOrCreateEthernet().MacAddress = ygot.String(lu.Attrs().HardwareAddr.String())

	iface.Ifindex = ygot.Uint32(uint32(lu.Attrs().Index))
	sb := &ygnmi.SetBatch{}

	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Ethernet().MacAddress().State(), *iface.Ethernet.MacAddress)
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Ifindex().State(), *iface.Ifindex)
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ifindex().State(), *iface.Ifindex)
	if _, err := sb.Set(ctx, ni.c); err != nil {
		log.Warningf("failed to set link status: %v", err)
	}
}

// handleAddrUpdate modifies the state based on changes to addresses.
func (ni *Reconciler) handleAddrUpdate(ctx context.Context, au *netlink.AddrUpdate) {
	ni.stateMu.Lock()
	defer ni.stateMu.Unlock()

	intf, data := ni.ocInterfaceData.findByIfIndex(au.LinkIndex)
	if data == nil {
		log.Infof("skipping address reconcilion for %s/%d, no interface data", intf.name, intf.subintf)
		return
	}

	sb := &ygnmi.SetBatch{}
	sub := ni.getOrCreateInterface(intf.name).GetOrCreateSubinterface(intf.subintf)

	ip := au.LinkAddress.IP.String()
	ipBytes := au.LinkAddress.IP.To4()
	mask := net.CIDRMask(32, 32)
	if ipBytes == nil {
		ipBytes = au.LinkAddress.IP.To16()
		mask = net.CIDRMask(128, 128)
	}
	pl, _ := au.LinkAddress.Mask.Size()
	isV4 := au.LinkAddress.IP.To4() != nil

	log.V(1).Infof("handling addr update for %s ip %v pl %v", data.hostifDevName, ip, pl)

	if ni.niDetail[data.networkInstance] == nil {
		log.Infof("skipping address reconcilion for %s/%d, unknown VRF %q", intf.name, intf.subintf, data.networkInstance)
		return
	}

	// The dataplane does not monitor the local interface's IP addr, they must set externally.
	if au.NewAddr {
		if isV4 {
			sub.GetOrCreateIpv4().GetOrCreateAddress(ip).PrefixLength = ygot.Uint8(uint8(pl))
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv4().Address(ip).Ip().State(), au.LinkAddress.IP.String())
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv4().Address(ip).PrefixLength().State(), uint8(pl))
		} else {
			sub.GetOrCreateIpv6().GetOrCreateAddress(ip).PrefixLength = ygot.Uint8(uint8(pl))
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv6().Address(ip).Ip().State(), au.LinkAddress.IP.String())
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv6().Address(ip).PrefixLength().State(), uint8(pl))
		}
		_, err := ni.routeClient.CreateRouteEntry(ctx, &saipb.CreateRouteEntryRequest{
			Entry: &saipb.RouteEntry{
				SwitchId:    ni.switchID,
				VrId:        ni.niDetail[data.networkInstance].vrOID,
				Destination: &saipb.IpPrefix{Addr: ipBytes, Mask: mask},
			},
			NextHopId:    proto.Uint64(ni.cpuPortID),
			PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD.Enum(),
		})
		if err != nil {
			log.Warningf("failed to add connected on intf %v route: %v", intf.name, err)
			return
		}
	} else {
		if isV4 {
			sub.GetOrCreateIpv4().DeleteAddress(ip)
			gnmiclient.BatchDelete(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv4().Address(ip).State())
		} else {
			sub.GetOrCreateIpv6().DeleteAddress(ip)
			gnmiclient.BatchDelete(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv6().Address(ip).State())
		}

		_, err := ni.routeClient.RemoveRouteEntry(ctx, &saipb.RemoveRouteEntryRequest{
			Entry: &saipb.RouteEntry{
				SwitchId:    ni.switchID,
				VrId:        ni.niDetail[data.networkInstance].vrOID,
				Destination: &saipb.IpPrefix{Addr: ipBytes, Mask: mask},
			},
		})
		if err != nil {
			log.Warningf("failed to remove connected route on intf %v: %v", intf.name, err)
			return
		}
	}
	if _, err := sb.Set(ctx, ni.c); err != nil {
		log.Warningf("failed to set link ip address: %v", err)
	}
}

// handleNeighborUpdate modifies the state based on changes to the neighbor.
func (ni *Reconciler) handleNeighborUpdate(ctx context.Context, nu *netlink.NeighUpdate) {
	ni.stateMu.Lock()
	defer ni.stateMu.Unlock()
	log.V(1).Infof("handling neighbor update for %s on %d", nu.IP.String(), nu.LinkIndex)

	intf, data := ni.ocInterfaceData.findByIfIndex(nu.LinkIndex)
	if data == nil {
		return
	}

	sb := &ygnmi.SetBatch{}
	sub := ni.getOrCreateInterface(intf.name).GetOrCreateSubinterface(intf.subintf)

	switch nu.Type {
	case unix.RTM_DELNEIGH:
		_, err := ni.neighborClient.RemoveNeighborEntry(ctx, &saipb.RemoveNeighborEntryRequest{
			Entry: &saipb.NeighborEntry{
				SwitchId:  ni.switchID,
				RifId:     data.rifID,
				IpAddress: ipToBytes(nu.IP),
			},
		})
		if err != nil {
			log.Warningf("failed to remove neighbor to dataplane: %v", err)
			return
		}
		if nu.Family == unix.AF_INET6 {
			sub.GetOrCreateIpv6().DeleteNeighbor(nu.IP.String())
			gnmiclient.BatchDelete(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv6().Neighbor(nu.IP.String()).State())
		} else {
			sub.GetOrCreateIpv4().DeleteNeighbor(nu.IP.String())
			gnmiclient.BatchDelete(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv4().Neighbor(nu.IP.String()).State())
		}
	case unix.RTM_NEWNEIGH:
		if len(nu.HardwareAddr) == 0 {
			log.Info("skipping neighbor update with no hwaddr")
			return
		}
		_, err := ni.neighborClient.CreateNeighborEntry(ctx, &saipb.CreateNeighborEntryRequest{
			Entry: &saipb.NeighborEntry{
				SwitchId:  ni.switchID,
				RifId:     data.rifID,
				IpAddress: ipToBytes(nu.IP),
			},
			DstMacAddress: nu.HardwareAddr,
		})
		if err != nil {
			log.Warningf("failed to create neighbor entry: %v", err)
		}
		if nu.Family == unix.AF_INET6 {
			neigh := sub.GetOrCreateIpv6().GetOrCreateNeighbor(nu.IP.String())
			neigh.LinkLayerAddress = ygot.String(nu.HardwareAddr.String())
			if nu.Flags&unix.NUD_PERMANENT != 0 {
				neigh.Origin = oc.IfIp_NeighborOrigin_STATIC
			} else {
				neigh.Origin = oc.IfIp_NeighborOrigin_DYNAMIC
			}
			gnmiclient.BatchReplace(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv6().Neighbor(nu.IP.String()).Ip().State(), neigh.GetIp())
			gnmiclient.BatchReplace(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv6().Neighbor(nu.IP.String()).LinkLayerAddress().State(), neigh.GetLinkLayerAddress())
			gnmiclient.BatchReplace(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv6().Neighbor(nu.IP.String()).Origin().State(), neigh.GetOrigin())
		} else {
			neigh := sub.GetOrCreateIpv4().GetOrCreateNeighbor(nu.IP.String())
			neigh.LinkLayerAddress = ygot.String(nu.HardwareAddr.String())
			if nu.Flags&unix.NUD_PERMANENT != 0 {
				neigh.Origin = oc.IfIp_NeighborOrigin_STATIC
			} else {
				neigh.Origin = oc.IfIp_NeighborOrigin_DYNAMIC
			}
			gnmiclient.BatchReplace(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv4().Neighbor(nu.IP.String()).Ip().State(), neigh.GetIp())
			gnmiclient.BatchReplace(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv4().Neighbor(nu.IP.String()).LinkLayerAddress().State(), neigh.GetLinkLayerAddress())
			gnmiclient.BatchReplace(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Ipv4().Neighbor(nu.IP.String()).Origin().State(), neigh.GetOrigin())
		}
	default:
		log.Warningf("unknown neigh update type: %v", nu.Type)
	}

	if _, err := sb.Set(ctx, ni.c); err != nil {
		log.Warningf("failed to set link status: %v", err)
	}
}

const (
	internalSuffix = "-internal"
)

// setupPorts creates the dataplane ports and TAP interfaces for all interfaces on the device.
func (ni *Reconciler) setupPorts(ctx context.Context) error {
	ifs, err := ni.ifaceMgr.LinkList()
	if err != nil {
		return err
	}

	for idx, i := range ifs {
		// Skip loopback, k8s pod interface, and tap interfaces.
		if i.Attrs().Name == "lo" || i.Attrs().Name == "eth0" || strings.HasSuffix(i.Attrs().Name, internalSuffix) {
			continue
		}
		log.Info("creating interfaces for %v", i.Attrs().Name)
		ocIntf := ocInterface{
			name:    i.Attrs().Name,
			subintf: 0,
		}
		data := &interfaceData{
			networkInstance: fakedevice.DefaultNetworkInstance,
		}

		portResp, err := ni.portClient.CreatePort(ctx, &saipb.CreatePortRequest{
			Switch:     ni.switchID,
			HwLaneList: []uint32{uint32(idx)},
			AdminState: proto.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("failed to create port %q: %w", i.Attrs().Name, err)
		}
		data.portID = portResp.Oid
		nid, err := ni.fwdClient.ObjectNID(ctx, &fwdpb.ObjectNIDRequest{
			ContextId: &fwdpb.ContextId{Id: ni.contextID},
			ObjectId:  &fwdpb.ObjectId{Id: fmt.Sprint(portResp.Oid)},
		})
		if err != nil {
			return fmt.Errorf("failed to get port %q nid: %w", i.Attrs().Name, err)
		}
		data.portNID = nid.Nid

		hostifName := i.Attrs().Name + internalSuffix
		hostifResp, err := ni.hostifClient.CreateHostif(ctx, &saipb.CreateHostifRequest{
			Switch:     ni.switchID,
			Type:       saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId:      &portResp.Oid,
			Name:       []byte(hostifName),
			OperStatus: proto.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("failed to create host interface %q: %w", hostifName, err)
		}
		data.hostifID = hostifResp.Oid
		data.hostifDevName = hostifName

		tap, err := ni.ifaceMgr.LinkByName(hostifName)
		if err != nil {
			return fmt.Errorf("failed to find tap interface %q: %w", hostifName, err)
		}
		data.hostifIfIndex = tap.Attrs().Index

		log.Infof("creating router interface dev: %v, port id: %v, mac: %s, vr id: %d", intfRefToDevName(ocIntf), portResp.GetOid(), tap.Attrs().HardwareAddr.String(), ni.niDetail[fakedevice.DefaultNetworkInstance].vrOID)

		rifResp, err := ni.ifaceClient.CreateRouterInterface(ctx, &saipb.CreateRouterInterfaceRequest{
			Switch:          ni.switchID,
			Type:            saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT.Enum(),
			PortId:          &portResp.Oid,
			SrcMacAddress:   tap.Attrs().HardwareAddr,
			VirtualRouterId: proto.Uint64(ni.niDetail[fakedevice.DefaultNetworkInstance].vrOID), // Implicitly add the RIF to DEFAULT network instance.
		})
		if err != nil {
			return fmt.Errorf("failed to update MAC address for interface %q: %w", i.Attrs().Name, err)
		}
		data.rifID = rifResp.Oid

		ni.getOrCreateInterface(i.Attrs().Name).GetOrCreateEthernet().SetHwMacAddress(tap.Attrs().HardwareAddr.String())
		ni.getOrCreateInterface(i.Attrs().Name).GetOrCreateEthernet().SetMacAddress(tap.Attrs().HardwareAddr.String())

		sb := &ygnmi.SetBatch{}
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(i.Attrs().Name).Ethernet().HwMacAddress().State(), tap.Attrs().HardwareAddr.String())
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(i.Attrs().Name).Ethernet().MacAddress().State(), tap.Attrs().HardwareAddr.String())
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(i.Attrs().Name).OperStatus().State(), oc.Interface_OperStatus_UP)
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(i.Attrs().Name).AdminStatus().State(), oc.Interface_AdminStatus_UP)
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(i.Attrs().Name).Enabled().State(), true)
		gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(i.Attrs().Name).Subinterface(0).AdminStatus().State(), oc.Interface_AdminStatus_UP)
		if _, err := sb.Set(ctx, ni.c); err != nil {
			log.Warningf("failed to set link status: %v", err)
		}
		ni.ocInterfaceData[ocIntf] = data
	}
	return nil
}

// ipToBytes converts a net.IP to a slice of bytes of the correct length (4 for IPv4, 16 for IPv6).
func ipToBytes(ip net.IP) []byte {
	if ip.To4() != nil {
		return ip.To4()
	}
	return ip.To16()
}
