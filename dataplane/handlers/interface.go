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

// Package handlers contains gNMI task handlers.
package handlers

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/openconfig/lemming/dataplane/internal/engine"
	"github.com/openconfig/lemming/dataplane/internal/kernel"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/vishvananda/netlink"

	log "github.com/golang/glog"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Interface handles config updates to the /interfaces/... paths.
type Interface struct {
	c             *ygnmi.Client
	watchCancelFn context.CancelFunc
	linkDoneCh    chan struct{}
	addrDoneCh    chan struct{}
	fwd           fwdpb.ServiceClient
	stateMu       sync.RWMutex
	state         map[string]*oc.Interface
	idxToName     map[int]string
}

// NewInterface creates a new interface handler.
func NewInterface(yc *ygnmi.Client, fwd fwdpb.ServiceClient) *Interface {
	return &Interface{
		c:         yc,
		fwd:       fwd,
		idxToName: map[int]string{},
		state:     map[string]*oc.Interface{},
	}
}

// Start starts running the handler, watching the cache and the kernel interfaces.
func (ni *Interface) Start(ctx context.Context) error {
	log.Info("starting interface handler")
	b := &ocpath.Batch{}

	if err := ni.setupPorts(ctx); err != nil {
		return fmt.Errorf("failed to setup ports: %v", err)
	}

	b.AddPaths(
		ocpath.Root().InterfaceAny().Name().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Ethernet().MacAddress().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Enabled().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv4().AddressAny().Ip().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv4().AddressAny().PrefixLength().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv6().AddressAny().Ip().Config().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv6().AddressAny().PrefixLength().Config().PathStruct(),
	)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	ni.watchCancelFn = cancelFn

	watcher := ygnmi.Watch(cancelCtx, ni.c, b.Config(), func(val *ygnmi.Value[*oc.Root]) error {
		log.V(2).Info("reconciling interfaces")
		root, ok := val.Val()
		if !ok || root.Interface == nil {
			return ygnmi.Continue
		}
		for _, i := range root.Interface {
			ni.reconcile(i)
		}
		return ygnmi.Continue
	})

	ni.linkDoneCh = make(chan struct{})
	linkUpdateCh := make(chan netlink.LinkUpdate)
	ni.addrDoneCh = make(chan struct{})
	addrUpdateCh := make(chan netlink.AddrUpdate)

	if err := netlink.LinkSubscribe(linkUpdateCh, ni.linkDoneCh); err != nil {
		return fmt.Errorf("failed to sub to link: %v", err)
	}
	if err := netlink.AddrSubscribe(addrUpdateCh, ni.addrDoneCh); err != nil {
		return fmt.Errorf("failed to sub to addr: %v", err)
	}

	go func() {
		for {
			select {
			case up := <-linkUpdateCh:
				ni.handleLinkUpdate(ctx, &up)
			case up := <-addrUpdateCh:
				ni.handleAddrUpdate(ctx, &up)
			}
		}
	}()

	go func() {
		// TODO: handle error
		if _, err := watcher.Await(); err != nil {
			log.Warningf("interface watch err: %v", err)
		}
	}()

	return nil
}

// Stop stops all watchers.
func (ni *Interface) Stop() {
	// TODO: prevent stopping more than once.
	ni.watchCancelFn()
	close(ni.linkDoneCh)
	close(ni.addrDoneCh)
}

// reconcile compares the interface config with state and modifies state to match config.
func (ni *Interface) reconcile(config *oc.Interface) {
	ni.stateMu.RLock()
	defer ni.stateMu.RUnlock()

	tapName := engine.IntfNameToTapName(config.GetName())
	state := ni.state[config.GetName()]

	// TODO: handle deleting interface.
	if config.GetOrCreateEthernet().MacAddress != nil {
		if config.GetEthernet().GetMacAddress() != state.GetEthernet().GetMacAddress() {
			log.V(1).Infof("setting interface %s hw-addr %q", engine.IntfNameToTapName(config.GetName()), config.GetEthernet().GetMacAddress())
			if err := kernel.SetInterfaceHWAddr(engine.IntfNameToTapName(config.GetName()), config.GetEthernet().GetMacAddress()); err != nil {
				log.Warningf("Failed to set mac address of port: %v", err)
				return
			}
		}
	}
	if config.GetOrCreateSubinterface(0).Enabled != nil {
		if state.GetSubinterface(0).Enabled == nil || config.GetSubinterface(0).GetEnabled() != state.GetSubinterface(0).GetEnabled() {
			log.V(1).Infof("setting interface %s enabled %t", engine.IntfNameToTapName(config.GetName()), config.GetSubinterface(0).GetEnabled())
			if err := kernel.SetInterfaceState(engine.IntfNameToTapName(config.GetName()), config.GetSubinterface(0).GetEnabled()); err != nil {
				log.Warningf("Failed to set state address of port: %v", err)
				return
			}
		}
	}
	// TODO: refactor this.
	for _, addr := range config.GetOrCreateSubinterface(0).GetOrCreateIpv4().Address {
		configIP := addr.Ip
		configPL := addr.PrefixLength
		var stateIP *string
		var statePL *uint8
		stateAddr := state.GetSubinterface(0).GetIpv4().GetAddress(addr.GetIp())
		if stateAddr != nil {
			stateIP = stateAddr.Ip
			statePL = stateAddr.PrefixLength
		}

		if configIP != nil && configPL != nil && (stateIP == nil || *statePL != *configPL) {
			log.V(1).Infof("Config IP: %v, Config PL: %v. State IP: %v, State PL: %v", addr.GetIp(), addr.GetPrefixLength(), stateAddr.GetIp(), stateAddr.GetPrefixLength())
			log.V(2).Infof("setting interface %s ip %s/%d", tapName, *configIP, *configPL)
			if err := kernel.SetInterfaceIP(engine.IntfNameToTapName(config.GetName()), *configIP, int(*configPL)); err != nil {
				log.Warningf("Failed to set ip address of port: %v", err)
				return
			}
		}
	}
	for _, addr := range config.GetOrCreateSubinterface(0).GetOrCreateIpv6().Address {
		configIP := addr.Ip
		configPL := addr.PrefixLength
		var stateIP *string
		var statePL *uint8
		stateAddr := state.GetSubinterface(0).GetIpv6().GetAddress(addr.GetIp())
		if stateAddr != nil {
			stateIP = stateAddr.Ip
			statePL = stateAddr.PrefixLength
		}

		if configIP != nil && configPL != nil && (stateIP == nil || *statePL != *configPL) {
			log.V(1).Infof("Config IP: %v, Config PL: %v. State IP: %v, State PL: %v", addr.GetIp(), addr.GetPrefixLength(), stateAddr.GetIp(), stateAddr.GetPrefixLength())
			log.V(2).Infof("setting interface %s ip %s/%d", tapName, *configIP, *configPL)
			if err := kernel.SetInterfaceIP(engine.IntfNameToTapName(config.GetName()), *configIP, int(*configPL)); err != nil {
				log.Warningf("Failed to set ip address of port: %v", err)
				return
			}
		}
	}
	// TODO: delete IPs
}

// getOrCreateInterface returns the state interface from the cache.
func (ni *Interface) getOrCreateInterface(iface string) *oc.Interface {
	if _, ok := ni.state[iface]; !ok {
		ni.state[iface] = &oc.Interface{
			Name: &iface,
		}
	}
	return ni.state[iface]
}

// handleLinkUpdate modifies the state based on changes to link state.
func (ni *Interface) handleLinkUpdate(ctx context.Context, lu *netlink.LinkUpdate) {
	ni.stateMu.Lock()
	defer ni.stateMu.Unlock()
	if !engine.IsTap(lu.Attrs().Name) {
		return
	}
	log.V(1).Infof("handling link update for %s", lu.Attrs().Name)

	modelName := engine.TapNameToIntfName(lu.Attrs().Name)
	iface := ni.getOrCreateInterface(modelName)
	iface.GetOrCreateEthernet().MacAddress = ygot.String(lu.Attrs().HardwareAddr.String())
	iface.Ifindex = ygot.Uint32(uint32(lu.Attrs().Index))
	iface.Enabled = ygot.Bool(lu.Attrs().Flags&net.FlagUp != 0)
	iface.AdminStatus = oc.Interface_AdminStatus_DOWN
	if *iface.Enabled {
		iface.AdminStatus = oc.Interface_AdminStatus_UP
	}
	// TODO: handle other states.
	var operStatus oc.E_Interface_OperStatus
	switch lu.Attrs().OperState {
	case netlink.OperDown:
		operStatus = oc.Interface_OperStatus_DOWN
	case netlink.OperUp, netlink.OperUnknown: // TAP interface may be unknown state because the dataplane doesn't bind to its fd, so treat unknown as up.
		operStatus = oc.Interface_OperStatus_UP
	}
	iface.OperStatus = operStatus

	sb := &ygnmi.SetBatch{}
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(modelName).Ifindex().State(), *iface.Ifindex)
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(modelName).Enabled().State(), *iface.Enabled)
	gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(modelName).OperStatus().State(), iface.OperStatus)
	if _, err := sb.Set(ctx, ni.c); err != nil {
		log.Warningf("failed to set link status: %v", err)
	}
}

// handleLinkUpdate modifies the state based on changes to addresses.
func (ni *Interface) handleAddrUpdate(ctx context.Context, au *netlink.AddrUpdate) {
	ni.stateMu.Lock()
	defer ni.stateMu.Unlock()
	name := ni.idxToName[au.LinkIndex]
	if !engine.IsTap(name) {
		return
	}

	sb := &ygnmi.SetBatch{}
	modelName := engine.TapNameToIntfName(name)
	sub := ni.getOrCreateInterface(modelName).GetOrCreateSubinterface(0)

	ip := au.LinkAddress.IP.String()
	pl, _ := au.LinkAddress.Mask.Size()
	isV4 := au.LinkAddress.IP.To4() != nil
	log.V(1).Infof("handling addr update for %s ip %v pl %v", name, ip, pl)
	if au.NewAddr {
		if isV4 {
			sub.GetOrCreateIpv4().GetOrCreateAddress(ip).PrefixLength = ygot.Uint8(uint8(pl))
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(modelName).Subinterface(0).Ipv4().Address(ip).Ip().State(), au.LinkAddress.IP.String())
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(modelName).Subinterface(0).Ipv4().Address(ip).PrefixLength().State(), uint8(pl))
		} else {
			sub.GetOrCreateIpv6().GetOrCreateAddress(ip).PrefixLength = ygot.Uint8(uint8(pl))
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(modelName).Subinterface(0).Ipv6().Address(ip).Ip().State(), au.LinkAddress.IP.String())
			gnmiclient.BatchUpdate(sb, ocpath.Root().Interface(modelName).Subinterface(0).Ipv6().Address(ip).PrefixLength().State(), uint8(pl))
		}
	} else {
		if isV4 {
			sub.GetOrCreateIpv4().DeleteAddress(ip)
			gnmiclient.BatchDelete(sb, ocpath.Root().Interface(modelName).Subinterface(0).Ipv4().Address(ip).State())
		} else {
			sub.GetOrCreateIpv6().DeleteAddress(ip)
			gnmiclient.BatchDelete(sb, ocpath.Root().Interface(modelName).Subinterface(0).Ipv6().Address(ip).State())
		}
	}
	if _, err := sb.Set(ctx, ni.c); err != nil {
		log.Warningf("failed to set link status: %v", err)
	}
}

// setupPorts creates the dataplane ports and TAP interfaces for all interfaces on the device.
func (ni *Interface) setupPorts(ctx context.Context) error {
	ifs, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range ifs {
		// Skip loopback, k8s pod interface, and tap interfaces.
		if i.Name == "lo" || i.Name == "eth0" || engine.IsTap(i.Name) {
			continue
		}
		if err := kernel.CreateTAP(engine.IntfNameToTapName(i.Name)); err != nil {
			return fmt.Errorf("failed to create tap port %q: %w", engine.IntfNameToTapName(i.Name), err)
		}
		tap, err := net.InterfaceByName(engine.IntfNameToTapName(i.Name))
		if err != nil {
			return fmt.Errorf("failed to find tap interface %q: %w", engine.IntfNameToTapName(i.Name), err)
		}
		ni.idxToName[tap.Index] = tap.Name
		if err := engine.CreateLocalPort(ctx, ni.fwd, tap.Name); err != nil {
			return fmt.Errorf("failed to create internal port %q: %w", tap.Name, err)
		}
		if err := engine.CreateExternalPort(ctx, ni.fwd, i.Name); err != nil {
			return fmt.Errorf("failed to create external port %q: %w", i.Name, err)
		}
	}
	return nil
}
