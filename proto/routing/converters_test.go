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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func TestToLayers(t *testing.T) {
	tests := []struct {
		desc string
		hdr  *Header
		want []gopacket.SerializableLayer
	}{{
		desc: "mpls",
		hdr: &Header{
			Type:   HeaderType_HEADER_TYPE_MPLS,
			Labels: []uint32{100},
		},
		want: []gopacket.SerializableLayer{&layers.MPLS{
			Label: 100,
		}},
	}, {
		desc: "UDP4",
		hdr: &Header{
			Type:    HeaderType_HEADER_TYPE_UDP4,
			SrcIp:   "127.0.0.1",
			DstIp:   "127.0.0.2",
			SrcPort: 1000,
			DstPort: 2000,
			IpTtl:   1,
		},
		want: []gopacket.SerializableLayer{&layers.IPv4{
			SrcIP:    net.ParseIP("127.0.0.1"),
			DstIP:    net.ParseIP("127.0.0.2"),
			Protocol: layers.IPProtocolUDP,
			TTL:      1,
			Version:  4,
		}, &layers.UDP{
			SrcPort: layers.UDPPort(1000),
			DstPort: layers.UDPPort(2000),
		}},
	}, {
		desc: "UDP6",
		hdr: &Header{
			Type:    HeaderType_HEADER_TYPE_UDP6,
			SrcIp:   "::1",
			DstIp:   "::2",
			SrcPort: 1000,
			DstPort: 2000,
			IpTtl:   1,
		},
		want: []gopacket.SerializableLayer{&layers.IPv6{
			SrcIP:      net.ParseIP("::1"),
			DstIP:      net.ParseIP("::2"),
			NextHeader: layers.IPProtocolUDP,
			HopLimit:   1,
			Version:    6,
		}, &layers.UDP{
			SrcPort: layers.UDPPort(1000),
			DstPort: layers.UDPPort(2000),
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := ToLayers(tt.hdr)
			if d := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(layers.IPv6{}), cmpopts.IgnoreUnexported(layers.UDP{})); d != "" {
				t.Errorf("ToLayers() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}
