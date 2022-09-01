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

// Package kernel contains funcs that interact with the kernel (sycalls, netlink).
package kernel

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// SetInterfaceHWAddr sets the MAC address of a network interface.
func SetInterfaceHWAddr(name string, addr net.HardwareAddr) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("failed to get interface: %w", err)
	}
	if err := netlink.LinkSetHardwareAddr(link, addr); err != nil {
		return fmt.Errorf("failed to get hwaddr of link: %w", err)
	}
	return nil
}

// SetInterfaceIPs sets the IP addresses of a network interface.
// Any existing IPs on the interface are deleted.
func SetInterfaceIPs(name string, ips []*net.IPNet) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("failed to get interface: %w", err)
	}
	currIPs, err := netlink.AddrList(link, unix.AF_UNSPEC)
	if err != nil {
		return fmt.Errorf("failed to get existing ips: %w", err)
	}
	for i := range currIPs {
		if err := netlink.AddrDel(link, &currIPs[i]); err != nil {
			return fmt.Errorf("failed to add delete existing IP: %w", err)
		}
	}
	for _, ip := range ips {
		if err := netlink.AddrReplace(link, &netlink.Addr{IPNet: ip}); err != nil {
			return fmt.Errorf("failed to add ip to link: %w", err)
		}
	}
	return nil
}

// SetInterfaceState sets a links up or down.
func SetInterfaceState(name string, up bool) error {
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
func CreateTAP(name string) error {
	fd, err := unix.Open("/dev/net/tun", unix.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open tun file: %w", err)
	}
	req, err := unix.NewIfreq(name)
	if err != nil {
		return fmt.Errorf("failed to create interface req: %w", err)
	}
	req.SetUint16(unix.IFF_TAP | unix.IFF_NO_PI)
	if err := unix.IoctlIfreq(fd, unix.TUNSETIFF, req); err != nil {
		return fmt.Errorf("failed to do ioctl: %v", err)
	}
	return nil
}
