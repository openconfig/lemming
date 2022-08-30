package kernel

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"

	dpb "github.com/openconfig/lemming/proto/dataplane"
)

func UpdateTAPInterface(i *dpb.UpdatePortRequest) error {
	tapName := 
	if err := createTAP(tapName); err != nil {
		return err
	}
	link, err := netlink.LinkByName(tapName)
	if err != nil {
		return fmt.Errorf("failed to get tap interface: %v", err)
	}
	if err := netlink.LinkSetHardwareAddr(link, i.hwaddr); err != nil {
		return fmt.Errorf("failed to get hwaddr of link: %v", err)
	}
	var allIPs []*net.IPNet
	allIPs = append(allIPs, i.ipv4...)
	allIPs = append(allIPs, i.ipv6...)
	for _, ip := range allIPs {
		if err := netlink.AddrReplace(link, &netlink.Addr{IPNet: ip}); err != nil {
			return fmt.Errorf("failed to add ip to link: %v", err)
		}
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("failed to turn link up: %v", err)
	}
	return nil
}

func createTAP(name string) error {
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
