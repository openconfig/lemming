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

package packet_test

import (
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/packettestutil"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ip"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// IP4 header carrying opaque data.
var ip4 = []byte{0x45, 0x01, 0x00, 0x16, 0x00, 0x00, 0x00, 0x00, 0x08, 0xff, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d, 0x00, 0x00}

func TestIP4(t *testing.T) {
	queries := []packettestutil.FieldQuery{
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
			Result: []byte{0x04},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0),
			Result: []byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0),
			Result: []byte{0x0a, 0x0b, 0x0c, 0x0d},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 0),
			Result: []byte{0x01},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP, 0),
			Result: []byte{0x08},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 0),
			Result: []byte{0xff},
		},
	}

	updates := []packettestutil.FieldUpdate{
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0),
			Arg: []byte{0x11, 0x12, 0x13, 0x14},
			Op:  fwdpacket.OpSet,
		},
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0),
			Arg: []byte{0x1a, 0x1b, 0x1c, 0x1d},
			Op:  fwdpacket.OpSet,
		},
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 0),
			Arg: []byte{0x0d},
			Op:  fwdpacket.OpSet,
		},
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP, 0),
			Arg: []byte{0x01},
			Op:  fwdpacket.OpDec,
		},
	}

	ipFinal := []byte{0x45, 0x0d, 0x00, 0x16, 0x00, 0x00, 0x00, 0x00, 0x07, 0xff, 0x58, 0x7f, 0x11, 0x12, 0x13, 0x14, 0x1a, 0x1b, 0x1c, 0x1d, 0x00, 0x00}

	tests := []packettestutil.PacketFieldTest{
		// IP packet.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP,
			Orig: [][]byte{
				ip4,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ipFinal,
			},
		},
		// IP4 packet.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
			Orig: [][]byte{
				ip4,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ipFinal,
			},
		},
		// IP4 over Ethernet frame.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				ip4,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernetIP4,
				ipFinal,
			},
		},
		// IP4 over VLAN frame.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetVLANIP4,
				ip4,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernetVLANIP4,
				ipFinal,
			},
		},
		// IP4 over 1Q frame.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernet1QIP4,
				ip4,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernet1QIP4,
				ipFinal,
			},
		},
	}
	packettestutil.TestPacketFields("ipv4", t, tests)
}

func TestIP4TTL(t *testing.T) {
	// Initial IP4 header carrying TCP data.
	ip4tcpInitial := []byte{0x45, 0x01, 0x00, 0x2c, 0x00, 0x00, 0x00, 0x00, 0xff, 0x06, 0xa1, 0xad, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d}

	// TCP header.
	tcpSegment := []byte{0x01, 0x02, 0x03, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x51, 0x34, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x0b, 0x0c, 0x0d}

	updates := []packettestutil.FieldUpdate{
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP, 0),
			Arg: []byte{0x01},
			Op:  fwdpacket.OpDec,
		},
	}

	// Final IP4 header carrying TCP data.
	ip4tcpFinal := []byte{0x45, 0x01, 0x00, 0x2c, 0x00, 0x00, 0x00, 0x00, 0xfe, 0x06, 0xa2, 0xad, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d}

	tests := []packettestutil.PacketFieldTest{
		// IP4 over Ethernet frame.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				ip4tcpInitial,
				tcpSegment,
			},
			Updates: updates,
			Final: [][]byte{
				ethernetIP4,
				ip4tcpFinal,
				tcpSegment,
			},
		},
	}
	packettestutil.TestPacketFields("ipv4", t, tests)
}

// TestReparse adds an opaque ethernet header to an ethernet/IP header and reparses it.
func TestReparse(t *testing.T) {
	// Create the original ethernet packet
	originalBytes := append(ethernetIP4, ip4...)
	p, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, originalBytes)
	if err != nil {
		t.Fatalf("Unable to create ethernet packet from %x", originalBytes)
	}

	// Set the VRF of the packet.
	vrf := []byte{0xff, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	updates := []packettestutil.FieldUpdate{
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF, 0),
			Arg: vrf,
			Op:  fwdpacket.OpSet,
		},
	}
	packettestutil.FieldUpdates(t, "reparse/current ethernet", 0, p, updates)

	// Query some fields of the ethernet header and the VRF.
	queries1 := []packettestutil.FieldQuery{
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0),
			Result: []byte{0x08, 0x00},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF, 0),
			Result: vrf,
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
			Result: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
			Result: []byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b},
		},
	}
	packettestutil.FieldQueries(t, "reparse/current ethernet", 0, p, queries1)

	// Reparse the packet as an ethernet header after prepending a new
	// ethernet header as an opaque set of bytes.
	prepend := []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x81, 0x00, 0x10, 0x01, 0x08, 0xCC}
	if err := p.Reparse(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, []fwdpacket.FieldID{
		fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF, 0),
	}, prepend); err != nil {
		t.Fatalf("reparse: Reparse failed with err %v", err)
	}
	t.Logf("reparsed packet is %+v", p)

	// Query fields of the ethernet header and the VRF. This effectively
	// encapsulates the ether/IP packet within another ethernet/vlan header.
	queries2 := []packettestutil.FieldQuery{
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0),
			Result: []byte{0x08, 0xCC},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF, 0),
			Result: vrf,
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
			Result: []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
			Result: []byte{0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 0),
			Result: []byte{0x00, 0x01},
		},
	}
	packettestutil.FieldQueries(t, "reparse/new ethernet", 0, p, queries2)
}
