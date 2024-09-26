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

// Package kernel contains funcs that interact with the kernel (sycalls, netlink).
package kernel

import (
	"fmt"
	"net"
	"sort"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// Interfaces contains methods for modifying networking interfaces.
type Interfaces struct{}

// SetHWAddr sets the MAC address of a network interface.
func (k *Interfaces) SetHWAddr(name string, addr string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("failed to get interface: %w", err)
	}
	addrBytes, err := net.ParseMAC(addr)
	if err != nil {
		return fmt.Errorf("failed to set parse addres: %v", err)
	}
	if err := netlink.LinkSetHardwareAddr(link, addrBytes); err != nil {
		return fmt.Errorf("failed to get hwaddr of link: %w", err)
	}
	return nil
}

// ReplaceIP sets the IP addresses of a network interface.
func (k *Interfaces) ReplaceIP(name string, ip string, prefixLen int) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("failed to get interface: %w", err)
	}
	ipNet := &net.IPNet{}
	ipNet.IP = net.ParseIP(ip)
	if ipNet.IP == nil {
		return fmt.Errorf("failed to parse ip")
	}
	ipNet.Mask = net.CIDRMask(prefixLen, 128)
	if ipNet.IP.To4() != nil { // If ip is IPv4.
		ipNet.Mask = net.CIDRMask(prefixLen, 32)
	}
	if err := netlink.AddrReplace(link, &netlink.Addr{IPNet: ipNet}); err != nil {
		return fmt.Errorf("failed to add ip to link: %w", err)
	}
	return nil
}

// DeleteIP delete an IP addresses from a network interface.
func (k *Interfaces) DeleteIP(name string, ip string, prefixLen int) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("failed to get interface: %w", err)
	}
	ipNet := &net.IPNet{}
	ipNet.IP = net.ParseIP(ip)
	if ipNet.IP == nil {
		return fmt.Errorf("failed to parse ip")
	}
	ipNet.Mask = net.CIDRMask(prefixLen, 128)
	if ipNet.IP.To4() != nil { // If ip is IPv4.
		ipNet.Mask = net.CIDRMask(prefixLen, 32)
	}
	if err := netlink.AddrDel(link, &netlink.Addr{IPNet: ipNet}); err != nil {
		return fmt.Errorf("failed to add ip to link: %w", err)
	}
	return nil
}

// SetState sets a links up or down.
func (k *Interfaces) SetState(name string, up bool) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("failed to get interface: %w", err)
	}
	if up {
		return netlink.LinkSetUp(link)
	}
	return netlink.LinkSetDown(link)
}

// CreateTAP creates kernel TAP interface.
func (k *Interfaces) CreateTAP(name string) (int, error) {
	fd, err := unix.Open("/dev/net/tun", unix.O_RDWR, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to open tun file: %w", err)
	}
	req, err := unix.NewIfreq(name)
	if err != nil {
		return 0, fmt.Errorf("failed to create interface req: %w", err)
	}
	req.SetUint16(unix.IFF_TAP | unix.IFF_NO_PI)
	if err := unix.IoctlIfreq(fd, unix.TUNSETIFF, req); err != nil {
		return 0, fmt.Errorf("failed to do ioctl: %v", err)
	}
	return fd, nil
}

// LinkByName returns a network interface by name.
func (k *Interfaces) LinkByName(name string) (netlink.Link, error) {
	return netlink.LinkByName(name)
}

// LinkSubscribe subscribes to link status for all interfaces.
func (k *Interfaces) LinkSubscribe(ch chan<- netlink.LinkUpdate, done <-chan struct{}) error {
	return netlink.LinkSubscribe(ch, done)
}

// AddrSubscribe subscribes to address changes for all interfaces.
func (k *Interfaces) AddrSubscribe(ch chan<- netlink.AddrUpdate, done <-chan struct{}) error {
	return netlink.AddrSubscribe(ch, done)
}

// NeighSubscribe subscribes to neighbor table updates.
func (k *Interfaces) NeighSubscribe(ch chan<- netlink.NeighUpdate, done <-chan struct{}) error {
	return netlink.NeighSubscribe(ch, done)
}

// LinkList lists all Linux network interfaces.
func (k *Interfaces) LinkList() ([]netlink.Link, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	sort.Slice(links, func(i, j int) bool {
		return links[i].Attrs().Name < links[j].Attrs().Name
	})
	return links, err
}

// LinkAdd adds a new network interface.
func (k *Interfaces) LinkAdd(link netlink.Link) error {
	return netlink.LinkAdd(link)
}

// LinkByIndex finds a link by its index.
func (k *Interfaces) LinkByIndex(idx int) (netlink.Link, error) {
	return netlink.LinkByIndex(idx)
}

// LinkSetDown sets the admin state to down.
func (k *Interfaces) LinkSetDown(link netlink.Link) error {
	return netlink.LinkSetDown(link)
}

// LinkSetDown sets the admin state to up.
func (k *Interfaces) LinkSetUp(link netlink.Link) error {
	return netlink.LinkSetUp(link)
}

// LinkSetMaster sets the member link's master to the other link.
func (k *Interfaces) LinkSetMaster(member netlink.Link, link netlink.Link) error {
	return netlink.LinkSetMaster(member, link)
}

// LinkSetNoMaster removes the master from the link.
func (k *Interfaces) LinkSetNoMaster(link netlink.Link) error {
	return netlink.LinkSetNoMaster(link)
}

// LinkModify modifies the link.
func (k *Interfaces) LinkModify(link netlink.Link) error {
	return netlink.LinkModify(link)
}

type PacketMetadata struct {
	SrcIfIndex int16
	DstIfIndex int16
	Context    uint32 // Context is extra value that can be set by the forwarding pipeline.
}
