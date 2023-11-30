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

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/internal/kernel"
	"github.com/openconfig/lemming/dataplane/standalone/packetio/cpusink"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	log "github.com/golang/glog"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type ocInterface struct {
	name    string
	subintf uint32
}

type interfaceData struct {
	portID        uint64
	hostifID      uint64
	hostifIfIndex int
	hostifDevName string
	rifID         uint64
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
	stateMu            sync.RWMutex
	// state keeps track of the applied state of the device's interfaces so that we do not issue duplicate configuration commands to the device's interfaces.
	state           map[string]*oc.Interface
	switchID        uint64
	ifaceMgr        interfaceManager
	ocInterfaceData interfaceMap
	cpuPortID       uint64
	contextID       string
}

type interfaceManager interface {
	SetHWAddr(name string, addr string) error
	SetState(name string, up bool) error
	ReplaceIP(name string, ip string, prefixLen int) error
	DeleteIP(name string, ip string, prefixLen int) error
	GetAll() ([]net.Interface, error)
	GetByName(name string) (*net.Interface, error)
	CreateTAP(name string) (int, error)
	LinkSubscribe(ch chan<- netlink.LinkUpdate, done <-chan struct{}) error
	AddrSubscribe(ch chan<- netlink.AddrUpdate, done <-chan struct{}) error
	NeighSubscribe(ch chan<- netlink.NeighUpdate, done <-chan struct{}) error
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
		hostifClient:       saipb.NewHostifClient(conn),
		portClient:         saipb.NewPortClient(conn),
		switchClient:       saipb.NewSwitchClient(conn),
		ifaceClient:        saipb.NewRouterInterfaceClient(conn),
		neighborClient:     saipb.NewNeighborClient(conn),
		routeClient:        saipb.NewRouteClient(conn),
		nextHopClient:      saipb.NewNextHopClient(conn),
		nextHopGroupClient: saipb.NewNextHopGroupClient(conn),
		fwdClient:          fwdpb.NewForwardingClient(conn),
	}
	return r
}

// Start starts running the handler, watching the cache and the kernel interfaces.
func (ni *Reconciler) StartInterface(ctx context.Context, client *ygnmi.Client) error {
	log.Info("starting interface handler")
	b := &ocpath.Batch{}
	ni.c = client

	if err := ni.setupPorts(ctx); err != nil {
		return fmt.Errorf("failed to setup ports: %v", err)
	}

	b.AddPaths(
		ocpath.Root().InterfaceAny().Name().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Ethernet().MacAddress().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Enabled().Config().PathStruct(), // TODO: Support the parent interface config/enabled controling the subinterface state.
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv4().AddressAny().Ip().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv4().AddressAny().PrefixLength().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv6().AddressAny().Ip().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv6().AddressAny().PrefixLength().Config().PathStruct(),
	)
	cancelCtx, cancelFn := context.WithCancel(ctx)

	watcher := ygnmi.Watch(cancelCtx, ni.c, b.Config(), func(val *ygnmi.Value[*oc.Root]) error {
		log.V(2).Info("reconciling interfaces")
		root, ok := val.Val()
		if !ok || root.Interface == nil {
			return ygnmi.Continue
		}
		for _, i := range root.Interface {
			ni.reconcile(cancelCtx, i)
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
				stats, err := ni.portClient.GetPortStats(ctx, &saipb.GetPortStatsRequest{
					Oid: ni.ocInterfaceData[intfName].portID,
					CounterIds: []saipb.PortStat{
						saipb.PortStat_PORT_STAT_IF_IN_UCAST_PKTS,
						saipb.PortStat_PORT_STAT_IF_IN_NON_UCAST_PKTS,
						saipb.PortStat_PORT_STAT_IF_OUT_UCAST_PKTS,
						saipb.PortStat_PORT_STAT_IF_OUT_NON_UCAST_PKTS,
					},
				})
				log.V(2).Infof("querying counters for interface %q, got %v", intfName, stats)
				if err != nil {
					log.Errorf("interface handler: could not retrieve counter for interface %q", intfName)
					continue
				}
				if _, err := gnmiclient.Replace(ctx, ni.c, ocpath.Root().Interface(intfName.name).Counters().InPkts().State(), stats.Values[0]+stats.Values[1]); err != nil {
					log.Errorf("interface handler: %v", err)
				}
				if _, err := gnmiclient.Replace(ctx, ni.c, ocpath.Root().Interface(intfName.name).Counters().OutPkts().State(), stats.Values[2]+stats.Values[2]); err != nil {
					log.Errorf("interface handler: %v", err)
				}
			}
		}
	}()
}

// reconcile compares the interface config with state and modifies state to match config.
func (ni *Reconciler) reconcile(ctx context.Context, config *oc.Interface) {
	ni.stateMu.RLock()
	defer ni.stateMu.RUnlock()

	intf := ocInterface{name: config.GetName(), subintf: 0}
	data := ni.ocInterfaceData[intf]
	if data == nil {
		return
	}
	state := ni.getOrCreateInterface(config.GetName())

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

	if config.GetOrCreateSubinterface(intf.subintf).Enabled != nil {
		if state.GetOrCreateSubinterface(intf.subintf).Enabled == nil || config.GetSubinterface(intf.subintf).GetEnabled() != state.GetSubinterface(intf.subintf).GetEnabled() {
			log.V(1).Infof("setting interface %s enabled %t", data.hostifDevName, config.GetSubinterface(intf.subintf).GetEnabled())
			_, err := ni.hostifClient.SetHostifAttribute(ctx, &saipb.SetHostifAttributeRequest{
				Oid:        data.hostifID,
				OperStatus: proto.Bool(config.GetSubinterface(0).GetEnabled()),
			})
			if err != nil {
				log.Warningf("Failed to set state address of hostif: %v", err)
			}
			_, err = ni.portClient.SetPortAttribute(ctx, &saipb.SetPortAttributeRequest{
				Oid:        data.portID,
				AdminState: proto.Bool(config.GetSubinterface(0).GetEnabled()),
			})
			if err != nil {
				log.Warningf("Failed to set state address of port: %v", err)
			}
			sb := &ygnmi.SetBatch{}
			enabled := config.GetSubinterface(intf.subintf).GetEnabled() && config.GetEnabled()
			adminStatus := oc.Interface_AdminStatus_DOWN
			if enabled {
				adminStatus = oc.Interface_AdminStatus_UP
			}
			// TODO: Right now treating subinterface 0 and interface as the same.
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Enabled().State(), enabled)
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).AdminStatus().State(), adminStatus)
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).Enabled().State(), enabled)
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(intf.name).Subinterface(intf.subintf).AdminStatus().State(), adminStatus)
			if _, err := sb.Set(ctx, ni.c); err != nil {
				log.Warningf("failed to set link status: %v", err)
			}
		}
	}

	type prefixPair struct {
		cfgIP, stateIP *string
		cfgPL, statePL *uint8
	}

	// Get all state IPs and their corresponding config IPs (if they exist).
	var interfacePairs []*prefixPair
	for _, addr := range state.GetOrCreateSubinterface(0).GetOrCreateIpv4().Address {
		pair := &prefixPair{
			stateIP: addr.Ip,
			statePL: addr.PrefixLength,
		}
		if pairAddr := config.GetSubinterface(0).GetIpv4().GetAddress(addr.GetIp()); pairAddr != nil {
			pair.cfgIP = pairAddr.Ip
			pair.cfgPL = pairAddr.PrefixLength
		}
		interfacePairs = append(interfacePairs, pair)
	}
	for _, addr := range state.GetOrCreateSubinterface(0).GetOrCreateIpv6().Address {
		pair := &prefixPair{
			stateIP: addr.Ip,
			statePL: addr.PrefixLength,
		}
		if pairAddr := config.GetSubinterface(0).GetIpv6().GetAddress(addr.GetIp()); pairAddr != nil {
			pair.cfgIP = pairAddr.Ip
			pair.cfgPL = pairAddr.PrefixLength
		}
		interfacePairs = append(interfacePairs, pair)
	}

	// Get all config IPs and their corresponding state IPs (if they exist).
	for _, addr := range config.GetOrCreateSubinterface(0).GetOrCreateIpv4().Address {
		pair := &prefixPair{
			cfgIP: addr.Ip,
			cfgPL: addr.PrefixLength,
		}
		if pairAddr := state.GetSubinterface(0).GetIpv4().GetAddress(addr.GetIp()); pairAddr != nil {
			pair.stateIP = pairAddr.Ip
			pair.statePL = pairAddr.PrefixLength
		}
		interfacePairs = append(interfacePairs, pair)
	}
	for _, addr := range config.GetOrCreateSubinterface(0).GetOrCreateIpv6().Address {
		pair := &prefixPair{
			cfgIP: addr.Ip,
			cfgPL: addr.PrefixLength,
		}
		if pairAddr := state.GetSubinterface(0).GetIpv6().GetAddress(addr.GetIp()); pairAddr != nil {
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

	entry := fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(ipBytes)))

	log.V(1).Infof("handling addr update for %s ip %v pl %v", data.hostifDevName, ip, pl)
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
		_, err := ni.fwdClient.TableEntryAdd(ctx, fwdconfig.TableEntryAddRequest(ni.contextID, cpusink.IP2MeTable).
			AppendEntry(entry, fwdconfig.Action(fwdconfig.TransmitAction(fmt.Sprint(data.hostifID)))).Build())
		if err != nil {
			log.Warningf("failed to add route: %v", err)
			return
		}
		_, err = ni.routeClient.CreateRouteEntry(ctx, &saipb.CreateRouteEntryRequest{
			Entry: &saipb.RouteEntry{
				SwitchId:    ni.switchID,
				VrId:        0,
				Destination: &saipb.IpPrefix{Addr: ipBytes, Mask: mask},
			},
			NextHopId:    proto.Uint64(ni.cpuPortID),
			PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD.Enum(),
		})
		if err != nil {
			log.Warningf("failed to add route: %v", err)
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
		_, err := ni.fwdClient.TableEntryRemove(ctx, &fwdpb.TableEntryRemoveRequest{
			ContextId: &fwdpb.ContextId{Id: ni.contextID},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: cpusink.IP2MeTable}},
			EntryDesc: entry.Build(),
		})
		if err != nil {
			log.Warningf("failed to remove route: %v", err)
			return
		}
		_, err = ni.routeClient.RemoveRouteEntry(ctx, &saipb.RemoveRouteEntryRequest{
			Entry: &saipb.RouteEntry{
				SwitchId:    ni.switchID,
				VrId:        0,
				Destination: &saipb.IpPrefix{Addr: ipBytes, Mask: mask},
			},
		})
		if err != nil {
			log.Warningf("failed to remove route: %v", err)
			return
		}
	}
	if _, err := sb.Set(ctx, ni.c); err != nil {
		log.Warningf("failed to set link status: %v", err)
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
	ifs, err := ni.ifaceMgr.GetAll()
	if err != nil {
		return err
	}

	for _, i := range ifs {
		// Skip loopback, k8s pod interface, and tap interfaces.
		if i.Name == "lo" || i.Name == "eth0" || strings.HasSuffix(i.Name, internalSuffix) {
			continue
		}
		log.Info("creating interfaces for %v", i.Name)
		ocIntf := ocInterface{
			name:    i.Name,
			subintf: 0,
		}
		data := &interfaceData{}

		portResp, err := ni.portClient.CreatePort(ctx, &saipb.CreatePortRequest{
			Switch: ni.switchID,
		})
		if err != nil {
			return fmt.Errorf("failed to create port %q: %w", i.Name, err)
		}
		data.portID = portResp.Oid

		hostifName := i.Name + internalSuffix
		hostifResp, err := ni.hostifClient.CreateHostif(ctx, &saipb.CreateHostifRequest{
			Switch:     ni.switchID,
			Type:       saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId:      &portResp.Oid,
			Name:       []byte(hostifName),
			OperStatus: proto.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("failed to create host interface %q: %w", hostifName, err)
		}
		data.hostifID = hostifResp.Oid
		data.hostifDevName = hostifName

		tap, err := ni.ifaceMgr.GetByName(hostifName)
		if err != nil {
			return fmt.Errorf("failed to find tap interface %q: %w", hostifName, err)
		}
		data.hostifIfIndex = tap.Index

		rifResp, err := ni.ifaceClient.CreateRouterInterface(ctx, &saipb.CreateRouterInterfaceRequest{
			Switch:          ni.switchID,
			Type:            saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT.Enum(),
			PortId:          &portResp.Oid,
			SrcMacAddress:   tap.HardwareAddr,
			VirtualRouterId: proto.Uint64(0),
		})
		if err != nil {
			return fmt.Errorf("failed to update MAC address for interface %q: %w", i.Name, err)
		}
		data.rifID = rifResp.Oid
		ni.getOrCreateInterface(i.Name).GetOrCreateEthernet().SetHwMacAddress(tap.HardwareAddr.String())
		ni.getOrCreateInterface(i.Name).GetOrCreateEthernet().SetMacAddress(tap.HardwareAddr.String())
		if _, err := gnmiclient.Update(ctx, ni.c, ocpath.Root().Interface(i.Name).Ethernet().HwMacAddress().State(), tap.HardwareAddr.String()); err != nil {
			return fmt.Errorf("failed to set hw addr of interface %q: %v", tap.Name, err)
		}
		if _, err := gnmiclient.Update(ctx, ni.c, ocpath.Root().Interface(i.Name).Ethernet().MacAddress().State(), tap.HardwareAddr.String()); err != nil {
			return fmt.Errorf("failed to set hw addr of interface %q: %v", tap.Name, err)
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
