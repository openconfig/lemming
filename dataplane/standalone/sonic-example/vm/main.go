package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/vishvananda/netlink"
)

func main() {
	args := []string{
		"-display", "none", // No display driver
		"-accel", "kvm", // Enable accel
		"-m", "32768", // 32Gi RAM
		"-smp", "12", // 12 CPUs
		"-nographic",                          // Don't launch any windows
		"-drive", "file=/vm.img,format=qcow2", // OS disk
		// "-device", "pci-bridge,chassis_nr=1,id=pci.1",
		"-netdev", "bridge,id=mgmt",
		"-device", "e1000,netdev=mgmt",
		"-netdev", "user,id=mgmt2,restrict=y", // include a dummy which might be needed.
		"-device", "e1000,netdev=mgmt2",
	}
	if err := setUpMgmt(); err != nil {
		log.Fatal(err)
	}
	tapsArgs, err := createTaps()
	if err != nil {
		log.Fatal(err)
	}
	args = append(args, tapsArgs...)
	fmt.Println(args)
	cmd := exec.Command("qemu-system-x86_64", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func createTaps() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces: %v", err)
	}
	var ifaceArgs []string
	ifaceArgs = append(ifaceArgs, "-device", "pci-bridge,chassis_nr=1,id=pci.1")
	for i, iface := range ifaces {
		var ifx int
		if count, err := fmt.Sscanf(iface.Name, "eth%d", &ifx); err != nil || count != 1 {
			continue
		}
		if ifx == 0 {
			continue
		}
		ifaceArgs = append(ifaceArgs,
			"-netdev",
			fmt.Sprintf("tap,id=dataport%d,ifname=tap%d,script=/tc-tap.sh,downscript=no", ifx, ifx),
			"-device",
			fmt.Sprintf("e1000,netdev=dataport%d,bus=pci.1,addr=0x%d", ifx, i),
		)
	}
	return ifaceArgs, nil
}

func setUpMgmt() error {
	if out, err := exec.Command("brctl", "addbr", "br0").CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create bridge: %v: Output:\n%s", err, out)
	}
	if err := os.MkdirAll("/etc/qemu", 0755); err != nil {
		return err
	}
	if err := os.WriteFile("/etc/qemu/bridge.conf", []byte("allow br0\n"), 0644); err != nil {
		return err
	}

	eth0Link, err := netlink.LinkByName("eth0")
	if err != nil {
		return fmt.Errorf("failed to get eth0: %v", err)
	}
	addrs, err := netlink.AddrList(eth0Link, netlink.FAMILY_V4)
	if err != nil {
		return fmt.Errorf("failed to get eth0 addresses: %v", err)
	}
	fmt.Println("len addrs")

	if err := netlink.AddrDel(eth0Link, &addrs[0]); err != nil {
		return fmt.Errorf("failed to get eth0 addresses: %v", err)
	}
	if out, err := exec.Command("brctl", "addif", "br0", "eth0").CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add eth0 to bridge: %v: Output:\n%s", err, out)
	}
	// if out, err := exec.Command("ip", "addr", "add", "dev", "br0", addrs[0].IPNet.String()).CombinedOutput(); err != nil {
	// 	return fmt.Errorf("failed to add tap0 to bridge: %v: Output:\n%s", err, out)
	// }
	if out, err := exec.Command("ip", "link", "set", "dev", "br0", "up").CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set br0 up: %v: Output:\n%s", err, out)
	}
	// if out, err := exec.Command("dnsmasq", "--bind-interfaces", "--interface=br0", fmt.Sprintf("--dhcp-range=%s,%s", addrs[0].IP.String(), addrs[0].IP.String())).CombinedOutput(); err != nil {
	// 	return fmt.Errorf("failed to set br0 up: %v: Output:\n%s", err, out)
	// }
	return nil
}
