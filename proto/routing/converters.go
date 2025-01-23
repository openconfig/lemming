// Copyright 2025 Google LLC
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

package routing

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ToLayers(hdr *Header) []gopacket.SerializableLayer {
	layer := []gopacket.SerializableLayer{}
	switch hdr.Type {
	case HeaderType_HEADER_TYPE_MPLS:
		for _, label := range hdr.Labels {
			layer = append(layer, &layers.MPLS{Label: label})
		}
	case HeaderType_HEADER_TYPE_UDP4:
		ip := &layers.IPv4{
			Version:  4,
			Protocol: layers.IPProtocolUDP,
			SrcIP:    net.ParseIP(hdr.SrcIp),
			DstIP:    net.ParseIP(hdr.DstIp),
			TTL:      uint8(hdr.IpTtl),
		}
		udp := &layers.UDP{
			SrcPort: layers.UDPPort(hdr.SrcPort),
			DstPort: layers.UDPPort(hdr.DstPort),
		}
		layer = append(layer, ip, udp)
	case HeaderType_HEADER_TYPE_UDP6:
		ip := &layers.IPv6{
			Version:    6,
			NextHeader: layers.IPProtocolUDP,
			SrcIP:      net.ParseIP(hdr.SrcIp),
			DstIP:      net.ParseIP(hdr.DstIp),
			HopLimit:   uint8(hdr.IpTtl),
		}
		udp := &layers.UDP{
			SrcPort: layers.UDPPort(hdr.SrcPort),
			DstPort: layers.UDPPort(hdr.DstPort),
		}
		layer = append(layer, ip, udp)
	}
	return layer
}
