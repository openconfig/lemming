package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	args := []string{
		"-display", "none", // No display driver
		"-accel", "kvm", // Enable accel
		"-m", "32768", // 32Gi RAM
		"-smp", "12", // 12 CPUs
		"-nographic",                          // Don't launch any windows
		"-drive", "file=/vm.img,format=qcow2", // OS disk
		"-netdev", "user,hostfwd=tcp::22-:22,id=mgmt",
		"-device", "e1000,netdev=mgmt",
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
	var ifaceArgs []string
	ifaceArgs = append(ifaceArgs, "-device", "pci-bridge,chassis_nr=1,id=pci.1")
	for i := 1; i <= 32; i++ {
		ifaceArgs = append(ifaceArgs,
			"-netdev",
			fmt.Sprintf("tap,id=dataport%d,ifname=tap%d,script=no,downscript=no", i, i),
			"-device",
			fmt.Sprintf("e1000,netdev=dataport%d,bus=pci.1,addr=0x%x", i, i-1),
		)
	}
	return ifaceArgs, nil
}
