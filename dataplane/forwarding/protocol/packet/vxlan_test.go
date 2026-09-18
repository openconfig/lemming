// Copyright 2024 Google LLC
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

package packet_test

import (
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ip"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/packettestutil"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/udp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/vxlan"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// IP4 header carrying UDP/VXLAN data (proto 17, total length 54).
var ip4vxlan = []byte{0x45, 0x01, 0x00, 0x36, 0x00, 0x00, 0x00, 0x00, 0xff, 0x11, 0xa1, 0x98, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d}

// IP6 header carrying UDP/VXLAN data (proto 17, payload length 34).
var ip6vxlan = []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x22, 0x11, 0x04, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

// UDP header targeting VXLAN port 4789 (0x12b5), length 34.
var udpVxlan = []byte{0x01, 0x02, 0x12, 0xb5, 0x00, 0x22, 0x00, 0x00}

// VXLAN header with VNI 0x001234 (Flags 0x08).
var vxlanHdr = []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00}

// Inner Ethernet header (carrying opaque payload 0x0823).
var innerEth = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x08, 0x23}

// Inner payload.
var innerPayload = []byte{0x0a, 0x0b, 0x0c, 0x0d}

func TestVXLANFields(t *testing.T) {
	queries := []packettestutil.FieldQuery{
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC, 0),
			Result: []byte{0x01, 0x02},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST, 0),
			Result: []byte{0x12, 0xb5},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI, 0),
			Result: []byte{0x00, 0x12, 0x34},
		},
	}
	updates := []packettestutil.FieldUpdate{
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI, 0),
			Arg: []byte{0x00, 0x56, 0x78},
			Op:  fwdpacket.OpSet,
		},
	}
	tests := []packettestutil.PacketFieldTest{
		// VXLAN over IPv4
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				ip4vxlan,
				udpVxlan,
				vxlanHdr,
				innerEth,
				innerPayload,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernetIP4,
				ip4vxlan,
				udpVxlan,
				{0x08, 0x00, 0x00, 0x00, 0x00, 0x56, 0x78, 0x00},
				innerEth,
				innerPayload,
			},
		},
		// VXLAN over IPv6
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				ip6vxlan,
				udpVxlan,
				vxlanHdr,
				innerEth,
				innerPayload,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernetIP6,
				ip6vxlan,
				{0x01, 0x02, 0x12, 0xb5, 0x00, 0x22, 0xe6, 0xf5},
				{0x08, 0x00, 0x00, 0x00, 0x00, 0x56, 0x78, 0x00},
				innerEth,
				innerPayload,
			},
		},
	}
	packettestutil.TestPacketFields("vxlan", t, tests)
}
