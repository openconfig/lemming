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

// Package packet takes a packet hex string and prints the decoded packet.
package packet

import (
	"encoding/hex"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/spf13/cobra"

	"github.com/openconfig/lemming/internal/packetutil"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packet",
		Short: "Packet printing and decoding utilities.",
	}
	decode := &cobra.Command{
		Use:   "decode",
		Short: "Pretty prints a hex-encoded packet.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bytes, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}
			packet := gopacket.NewPacket(bytes, layers.LayerTypeEthernet, gopacket.Default)
			cmd.Println(packet.Dump())
			return nil
		},
	}

	var packetCount int

	pcap := &cobra.Command{
		Use:   "printpcap",
		Short: "Prints the contents of a pcap file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packetutil.DisplayCaptureFile(args[0], packetCount)
			return nil
		},
	}
	pcap.Flags().IntVar(&packetCount, "packet_count", 10, "Number of packets to packet_count")

	cmd.AddCommand(decode, pcap)

	return cmd
}
