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

// Package packetutil contains utilities for packets.
package packetutil

import (
	"bytes"
	"fmt"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

func DisplayCaptureFile(pcapfile string) error {
	f, err := os.Open(pcapfile)
	if err != nil {
		return fmt.Errorf("could not open pcap file %s: %v", f.Name(), err)
	}
	defer f.Close()

	handleRead, err := pcapgo.NewReader(f)
	if err != nil {
		return fmt.Errorf("could not create reader on pcap file %s: %v", f.Name(), err)
	}
	ps := gopacket.NewPacketSource(handleRead, layers.LinkTypeEthernet)

	for i := 0; i != 10; i++ {
		pkt, err := ps.NextPacket()
		if err != nil {
			return fmt.Errorf("error reading next packet: %v", err)
		}
		fmt.Println(pkt.Dump())
	}
	return nil
}

func DisplayCapture(captureBytes []byte) error {
	if len(captureBytes) == 0 {
		return fmt.Errorf("packetutil: zero captured bytes")
	}

	handleRead, err := pcapgo.NewReader(bytes.NewReader(captureBytes))
	if err != nil {
		return fmt.Errorf("could not create reader on captured bytes")
	}
	ps := gopacket.NewPacketSource(handleRead, layers.LinkTypeEthernet)

	for i := 0; i != 10; i++ {
		pkt, err := ps.NextPacket()
		if err != nil {
			return fmt.Errorf("error reading next packet: %v", err)
		}
		fmt.Println(pkt.Dump())
	}
	return nil
}
