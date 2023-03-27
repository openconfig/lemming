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

// package kerneltest contains fake implemetation of structs in kernel package
package kerneltest

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// Iface is a basic representation of a network interface.
type Iface struct {
	HWAddr net.HardwareAddr
	Idx    int
	IPs    map[string]struct{}
	up     bool
}

// New returns a new FakeInterfaces.
func New(initialIfaces map[string]*Iface) *FakeInterfaces {
	return &FakeInterfaces{
		Links: initialIfaces,
	}
}

// FakeInterfaces contains a fake implementation methods for modifying networking interfaces.
type FakeInterfaces struct {
	Links map[string]*Iface
	// Channels are set by calls LinkSubscribe and AddrSubscribe.
	linkCh chan<- netlink.LinkUpdate
	addrCh chan<- netlink.AddrUpdate
}

// SetHWAddr sets the MAC address of a network interface.
func (fi *FakeInterfaces) SetHWAddr(name string, addr string) error {
	l, ok := fi.Links[name]
	if !ok {
		return fmt.Errorf("link %s doesn't exist", name)
	}
	hwAddr, err := net.ParseMAC(addr)
	if err != nil {
		return err
	}
	l.HWAddr = hwAddr
	fi.sendLinkUpdate(name, l)

	return nil
}

// ReplaceIP sets the IP addresses of a network interface.
func (fi *FakeInterfaces) ReplaceIP(name string, ip string, prefixLen int) error {
	l, ok := fi.Links[name]
	if !ok {
		return fmt.Errorf("link %s doesn't exist", name)
	}
	cidrStr := fmt.Sprintf("%s/%d", ip, prefixLen)
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return err
	}
	if l.IPs == nil {
		l.IPs = make(map[string]struct{})
	}
	l.IPs[cidrStr] = struct{}{}
	fi.sendAddrUpdate(l, ipNet, true)

	return nil
}

// DeleteIP delete an IP addresses from a network interface.
func (fi *FakeInterfaces) DeleteIP(name string, ip string, prefixLen int) error {
	l, ok := fi.Links[name]
	if !ok {
		return fmt.Errorf("link %s doesn't exist", name)
	}
	cidrStr := fmt.Sprintf("%s/%d", ip, prefixLen)
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return err
	}
	if _, ok := l.IPs[cidrStr]; !ok {
		return fmt.Errorf("ip %s/%d not set in interface %s", ip, prefixLen, name)
	}
	delete(l.IPs, cidrStr)
	fi.sendAddrUpdate(l, ipNet, false)

	return nil
}

// SetState sets a links up or down.
func (fi *FakeInterfaces) SetState(name string, up bool) error {
	l, ok := fi.Links[name]
	if !ok {
		return fmt.Errorf("link %s doesn't exist", name)
	}
	l.up = up
	fi.sendLinkUpdate(name, l)
	return nil
}

// CreateTAP creates kernel TAP interface.
func (fi *FakeInterfaces) CreateTAP(name string) (int, error) {
	if _, ok := fi.Links[name]; ok {
		return 0, fmt.Errorf("link %s already exist", name)
	}
	fi.Links[name] = &Iface{
		HWAddr: []byte{1, 1, 0, 0, 1, 1},
		Idx:    len(fi.Links) + 1,
	}
	return 0, nil
}

// GetAll returns all interfaces.
func (fi *FakeInterfaces) GetAll() ([]net.Interface, error) {
	var ifaces []net.Interface
	for name, link := range fi.Links {
		var flags net.Flags
		if link.up {
			flags = net.FlagUp
		}
		ifaces = append(ifaces, net.Interface{
			Index:        link.Idx,
			Name:         name,
			HardwareAddr: link.HWAddr,
			Flags:        flags,
		})
	}
	return ifaces, nil
}

// GetByName returns all interfaces.
func (fi *FakeInterfaces) GetByName(name string) (*net.Interface, error) {
	l, ok := fi.Links[name]
	if !ok {
		return nil, fmt.Errorf("link %s doesn't exist", name)
	}
	var flags net.Flags
	if l.up {
		flags = net.FlagUp
	}
	return &net.Interface{
		Index:        l.Idx,
		Name:         name,
		HardwareAddr: l.HWAddr,
		Flags:        flags,
	}, nil
}

// LinkSubscribe subscribes to link status for all interfaces.
func (fi *FakeInterfaces) LinkSubscribe(ch chan<- netlink.LinkUpdate, done <-chan struct{}) error {
	fi.linkCh = ch
	go func() {
		<-done
		fi.linkCh = nil
		close(ch)
	}()
	return nil
}

// AddrSubscribe subscribes to address changes for all interfaces.
func (fi *FakeInterfaces) AddrSubscribe(ch chan<- netlink.AddrUpdate, done <-chan struct{}) error {
	fi.addrCh = ch
	go func() {
		<-done
		fi.addrCh = nil
		close(ch)
	}()
	return nil
}

// NeighSubscribe subscribes to neighbor table updates.
// TODO: not implemented.
func (fi *FakeInterfaces) NeighSubscribe(chan<- netlink.NeighUpdate, <-chan struct{}) error {
	return nil
}

func (fi *FakeInterfaces) sendLinkUpdate(name string, l *Iface) {
	if fi.linkCh == nil {
		return
	}
	var flag net.Flags
	var state netlink.LinkOperState = netlink.OperDown
	if l.up {
		flag = net.FlagUp
		state = netlink.OperUp
	}

	fi.linkCh <- netlink.LinkUpdate{
		Link: &netlink.Dummy{
			LinkAttrs: netlink.LinkAttrs{
				Name:         name,
				HardwareAddr: l.HWAddr,
				Flags:        flag,
				OperState:    state,
			},
		},
	}
}

func (fi *FakeInterfaces) sendAddrUpdate(l *Iface, ip *net.IPNet, new bool) {
	if fi.addrCh == nil {
		return
	}

	fi.addrCh <- netlink.AddrUpdate{
		LinkIndex:   l.Idx,
		LinkAddress: *ip,
		NewAddr:     new,
	}
}
