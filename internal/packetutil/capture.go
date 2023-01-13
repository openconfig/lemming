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
	"io"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

// DisplayCaptureFile displays the first n packets.
func DisplayCaptureFile(pcapfile string, n int) error {
	f, err := os.Open(pcapfile)
	if err != nil {
		return fmt.Errorf("could not open pcap file %s: %v", f.Name(), err)
	}
	defer f.Close()

	return displayCapture(f, n)
}

// DisplayCapture displays the first n packets from a PCAP capture.
func DisplayCapture(captureBytes []byte, n int) error {
	if len(captureBytes) == 0 {
		return fmt.Errorf("packetutil: zero captured bytes")
	}

	return displayCapture(bytes.NewReader(captureBytes), n)
}

func displayCapture(reader io.Reader, n int) error {
	handleRead, err := pcapgo.NewReader(reader)
	if err != nil {
		return fmt.Errorf("packetutil: could not create reader")
	}
	ps := gopacket.NewPacketSource(handleRead, layers.LinkTypeEthernet)

	for i := 0; i != n; i++ {
		pkt, err := ps.NextPacket()
		if err != nil {
			return fmt.Errorf("packetutil: error reading next packet: %v", err)
		}
		fmt.Println(pkt.Dump())
	}
	return nil
}
