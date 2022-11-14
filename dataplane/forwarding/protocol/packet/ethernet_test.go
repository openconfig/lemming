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
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/packettestutil"
)

// Ethernet headers carrying IPv4 packets.
var ethernetIP4 = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x08, 0x00}
var ethernetVLANIP4 = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x81, 0x00, 0x71, 0x23, 0x08, 0x00}
var ethernet1QIP4 = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x91, 0x00, 0x24, 0x56, 0x81, 0x00, 0x71, 0x23, 0x08, 0x00}

// Ethernet headers carrying IPv6 packets.
var ethernetIP6 = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x86, 0xDD}
var ethernetVLANIP6 = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x81, 0x00, 0x71, 0x23, 0x86, 0xDD}
var ethernet1QIP6 = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x91, 0x00, 0x24, 0x56, 0x81, 0x00, 0x71, 0x23, 0x86, 0xDD}

// Ethernet headers carrying ARP payload.
var ethernetARP = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x08, 0x06}
var ethernetVLANARP = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x81, 0x00, 0x71, 0x23, 0x08, 0x06}
var ethernet1QARP = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x91, 0x00, 0x24, 0x56, 0x81, 0x00, 0x71, 0x23, 0x08, 0x06}

// Ethernet header carrying ISIS
var payloadISIS = []byte{0x00,
	0x67, 0xfe, 0xfe, 0x03, 0x83, 0x1b, 0x01, 0x00, 0x14, 0x01, 0x00, 0x00, 0x00, 0x64, 0x04, 0xaf,
	0x44, 0x44, 0x44, 0x44, 0x44, 0x44, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0xf2, 0x52, 0x03, 0x01,
	0x04, 0x03, 0x49, 0x00, 0x14, 0x81, 0x01, 0xcc, 0x89, 0x02, 0x52, 0x34, 0x84, 0x04, 0x0a, 0x00,
	0x14, 0x01, 0x80, 0x0c, 0x0a, 0x80, 0x80, 0x80, 0x0a, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfc,
	0x02, 0x0c, 0x00, 0x0a, 0x80, 0x80, 0x80, 0x44, 0x44, 0x44, 0x44, 0x44, 0x44, 0x01, 0x80, 0x18,
	0x0a, 0x80, 0x80, 0x80, 0x0a, 0x00, 0x14, 0x00, 0xff, 0xff, 0xff, 0xfc, 0x14, 0x80, 0x80, 0x80,
	0xc0, 0xa8, 0x14, 0x00, 0xff, 0xff, 0xff, 0x00}

// TestEthernetFields parses various type of ethernet frames, and
// performs queries and updates on various ethernet header fields.
func TestEthernetFields(t *testing.T) {
	tests := []packettestutil.PacketFieldTest{
		// Parse an ethernet frame with no payload.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig:        [][]byte{{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x08, 0x23}},
			Queries: []packettestutil.FieldQuery{{
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				Err: "failed",
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
				Result: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
				Result: []byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0),
				Result: []byte{0x08, 0x23},
			}, {
				ID:     fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L2, 1, 2, 0),
				Result: []byte{0x01, 0x02},
			}, {
				ID:  fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L2, 10, 20, 0),
				Err: "failed",
			}},
			Updates: []packettestutil.FieldUpdate{{
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				Op:  fwdpacket.OpSet,
				Err: "failed",
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
				Arg: []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
				Arg: []byte{0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0),
				Arg: []byte{0x1C, 0x1D},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L2, 1, 2, 0),
				Arg: []byte{0x21, 0x22},
				Op:  fwdpacket.OpSet,
			}},
			Final: [][]byte{{0x10, 0x21, 0x22, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D}},
		}, { // Parse an ethernet frame with one vlan tag.
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig:        [][]byte{{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x81, 0x00, 0x71, 0x23, 0x08, 0x23}},
			Queries: []packettestutil.FieldQuery{{
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				Err: "failed",
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
				Result: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
				Result: []byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0),
				Result: []byte{0x08, 0x23},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 0),
				Result: []byte{0x01, 0x23},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_PRIORITY, 0),
				Result: []byte{0x00, 0x07},
			}, {
				ID:     fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L2, 15, 1, 0),
				Result: []byte{0x23},
			}},
			Updates: []packettestutil.FieldUpdate{{
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				Err: "failed",
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
				Arg: []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
				Arg: []byte{0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0),
				Arg: []byte{0x1C, 0x1D},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 0),
				Arg: []byte{0x02, 0x34},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_PRIORITY, 0),
				Arg: []byte{0x00, 0x04},
				Op:  fwdpacket.OpSet,
			}},
			Final: [][]byte{{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x81, 0x00, 0x42, 0x34, 0x1C, 0x1D}},
		}, { // Parse an ethernet frame with two vlan tags.
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig:        [][]byte{{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x91, 0x00, 0x24, 0x56, 0x81, 0x00, 0x71, 0x23, 0x08, 0x23}},
			Queries: []packettestutil.FieldQuery{{
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				Err: "failed",
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
				Result: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
				Result: []byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0),
				Result: []byte{0x08, 0x23},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 0),
				Result: []byte{0x04, 0x56},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_PRIORITY, 0),
				Result: []byte{0x00, 0x02},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 1),
				Result: []byte{0x01, 0x23},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_PRIORITY, 1),
				Result: []byte{0x00, 0x07},
			}},
			Updates: []packettestutil.FieldUpdate{{
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				Err: "failed",
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
				Arg: []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
				Arg: []byte{0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0),
				Arg: []byte{0x1C, 0x1D},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 0),
				Arg: []byte{0x05, 0x67},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_PRIORITY, 0),
				Arg: []byte{0x00, 0x01},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 1),
				Arg: []byte{0x02, 0x34},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_PRIORITY, 1),
				Arg: []byte{0x00, 0x04},
				Op:  fwdpacket.OpSet,
			},
			},
			Final: [][]byte{{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x91, 0x00, 0x15, 0x67, 0x81, 0x00, 0x42, 0x34, 0x1C, 0x1D}},
		}, { // Parse an ethernet frame with an ISIS payload.
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				{0x01, 0x80, 0xc2, 0x00, 0x00, 0x15, 0xc2, 0x03, 0x29, 0xa9, 0x00, 0x00},
				payloadISIS,
			},
			Queries: []packettestutil.FieldQuery{{
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				Err: "failed",
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
				Result: []byte{0x01, 0x80, 0xc2, 0x00, 0x00, 0x15},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
				Result: []byte{0xc2, 0x03, 0x29, 0xa9, 0x00, 0x00},
			}, {
				ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0),
				Result: []byte{0x00, 0x67},
			}},
			Updates: []packettestutil.FieldUpdate{{
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
				Arg: []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15},
				Op:  fwdpacket.OpSet,
			}, {
				ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
				Arg: []byte{0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b},
				Op:  fwdpacket.OpSet,
			}},
			Final: [][]byte{
				{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b},
				payloadISIS,
			},
		},
	}

	packettestutil.TestPacketFields("ethernet", t, tests)
}

// TestEthernetHeaderDecap performs decap operations using various
// combinations of ethernet headers and frames.
func TestEthernetHeaderDecap(t *testing.T) {
	tests := []packettestutil.PacketHeaderTest{
		// Strip out tags from a 1q ethernet frame.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernet1QARP,
				arp,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_1Q,
					Result: [][]byte{
						ethernetARP,
						arp,
					},
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_1Q,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Result: [][]byte{
						arp,
					},
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Err:   "failed",
				},
			},
		},
		// Strip out tags from a vlan ethernet frame.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetVLANARP,
				arp,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_1Q,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN,
					Result: [][]byte{
						ethernetARP,
						arp,
					},
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_1Q,
					Err:   "failed",
				},
				{
					Encap:  false,
					ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Result: [][]byte{arp},
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Err:   "failed",
				},
			},
		},
		// Strip out an ethernet header.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetARP,
				arp,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_1Q,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN,
					Err:   "failed",
				},
				{
					Encap:  false,
					ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Result: [][]byte{arp},
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Err:   "failed",
				},
			},
		},
		// Strip out a 1q ethernet header.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig:        [][]byte{ethernet1QARP, arp},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap:  false,
					ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Result: [][]byte{arp},
				},
			},
		},
		// Strip out a vlan ethernet header.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig:        [][]byte{ethernetVLANARP, arp},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap:  false,
					ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Result: [][]byte{arp},
				},
			},
		},
	}

	packettestutil.TestPacketHeaders("ethernet", t, tests)
}

// TestEthernetHeaderEncap performs encap operations using various
// combinations of ethernet headers and frames.
func TestEthernetHeaderEncap(t *testing.T) {
	tests := []packettestutil.PacketHeaderTest{
		// Remove and add a 1q ethernet header.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig:        [][]byte{ethernetVLANARP, arp},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap:  false,
					ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Result: [][]byte{arp},
				},
				{
					Encap:  true,
					ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_1Q,
					Result: [][]byte{ethernet1QARP, arp},
					Updates: []packettestutil.FieldUpdate{
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
							Arg: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
							Arg: []byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 0),
							Arg: []byte{0x04, 0x56},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_PRIORITY, 0),
							Arg: []byte{0x00, 0x02},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 1),
							Arg: []byte{0x01, 0x23},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_PRIORITY, 1),
							Arg: []byte{0x00, 0x07},
							Op:  fwdpacket.OpSet,
						},
					},
				},
			},
		},
		// Remove and add a vlan ethernet header.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig:        [][]byte{ethernetVLANARP, arp},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap:  false,
					ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Result: [][]byte{arp},
				},
				{
					Encap:  true,
					ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN,
					Result: [][]byte{ethernetVLANARP, arp},
					Updates: []packettestutil.FieldUpdate{
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
							Arg: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
							Arg: []byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 0),
							Arg: []byte{0x01, 0x23},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_PRIORITY, 0),
							Arg: []byte{0x00, 0x07},
							Op:  fwdpacket.OpSet,
						},
					},
				},
				{
					Encap: true,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_1Q,
					Err:   "failed",
				},
				{
					Encap: true,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN,
					Err:   "failed",
				},
				{
					Encap: true,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Err:   "failed",
				},
			},
		},
		// Add a vlan ethernet header to an IPv6 frame.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig:        [][]byte{ethernetIP6, ip6},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap:  true,
					ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN,
					Result: [][]byte{ethernetVLANIP6, ip6},
					Updates: []packettestutil.FieldUpdate{
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
							Arg: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
							Arg: []byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, 0),
							Arg: []byte{0x01, 0x23},
							Op:  fwdpacket.OpSet,
						},
						{
							ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_PRIORITY, 0),
							Arg: []byte{0x00, 0x07},
							Op:  fwdpacket.OpSet,
						},
					},
				},
				{
					Encap: true,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN,
					Err:   "failed",
				},
				{
					Encap: true,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					Err:   "failed",
				},
			},
		},
	}

	packettestutil.TestPacketHeaders("ethernet", t, tests)
}
