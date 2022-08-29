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
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/testutil"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

var arp = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c}

func TestARP(t *testing.T) {
	tests := []testutil.PacketFieldTest{
		// Parse an ethernet frame with an ARP message.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetARP,
				arp,
			},
			Queries: []testutil.FieldQuery{
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_VERSION, 0),
					Err: "failed",
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_TPA, 0),
					Result: []byte{0x19, 0x1a, 0x1b, 0x1c},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_TMAC, 0),
					Result: []byte{0x13, 0x14, 0x15, 0x16, 0x17, 0x18},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_SPA, 0),
					Result: []byte{0x0F, 0x10, 0x11, 0x12},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_SMAC, 0),
					Result: []byte{0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E},
				},
				{
					ID:     fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_L2, 38, 4, 0),
					Result: []byte{0x19, 0x1a, 0x1b, 0x1c},
				},
				{
					ID:     fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_L3, 24, 4, 0),
					Result: []byte{0x19, 0x1a, 0x1b, 0x1c},
				},
			},
			Updates: []testutil.FieldUpdate{
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_TPA, 0),
					Arg: []byte{0x29, 0x2a, 0x2b, 0x2c},
					Op:  fwdpacket.OpSet,
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_TMAC, 0),
					Arg: []byte{0x23, 0x24, 0x25, 0x26, 0x27, 0x28},
					Op:  fwdpacket.OpSet,
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_SPA, 0),
					Arg: []byte{0x3F, 0x30, 0x31, 0x32},
					Op:  fwdpacket.OpSet,
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_SMAC, 0),
					Arg: []byte{0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E},
					Op:  fwdpacket.OpSet,
				},
			},
			Final: [][]byte{
				ethernetARP,
				[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x30, 0x31, 0x32, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c},
			},
		},
		// Parse a vlan frame with an ARP message.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetVLANARP,
				arp,
			},
			Queries: []testutil.FieldQuery{
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_VERSION, 0),
					Err: "failed",
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ETHER_MAC_DST, 0),
					Result: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_TPA, 0),
					Result: []byte{0x19, 0x1a, 0x1b, 0x1c},
				},
				{
					ID:     fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_L2, 42, 4, 0),
					Result: []byte{0x19, 0x1a, 0x1b, 0x1c},
				},
				{
					ID:     fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_L3, 24, 4, 0),
					Result: []byte{0x19, 0x1a, 0x1b, 0x1c},
				},
			},
			Updates: []testutil.FieldUpdate{
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_TPA, 0),
					Arg: []byte{0x29, 0x2a, 0x2b, 0x2c},
					Op:  fwdpacket.OpSet,
				},
			},
			Final: [][]byte{
				ethernetVLANARP,
				[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x29, 0x2a, 0x2b, 0x2c},
			},
		},
		// Parse a 1q frame with an ARP message.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernet1QARP,
				arp,
			},
			Queries: []testutil.FieldQuery{
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_VERSION, 0),
					Err: "failed",
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ETHER_MAC_DST, 0),
					Result: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_TPA, 0),
					Result: []byte{0x19, 0x1a, 0x1b, 0x1c},
				},
				{
					ID:     fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_L2, 46, 4, 0),
					Result: []byte{0x19, 0x1a, 0x1b, 0x1c},
				},
				{
					ID:     fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_L3, 24, 4, 0),
					Result: []byte{0x19, 0x1a, 0x1b, 0x1c},
				},
				{
					ID:     fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_PAYLOAD, 24, 4, 0),
					Result: []byte{0x19, 0x1a, 0x1b, 0x1c},
				},
				{
					ID:  fwdpacket.NewFieldIDFromBytes(fwdpb.PacketHeaderGroup_L3, 38, 4, 0),
					Err: "failed",
				},
			},
			Updates: []testutil.FieldUpdate{
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ARP_TPA, 0),
					Arg: []byte{0x29, 0x2a, 0x2b, 0x2c},
					Op:  fwdpacket.OpSet,
				},
			},
			Final: [][]byte{
				ethernet1QARP,
				[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x29, 0x2a, 0x2b, 0x2c},
			},
		},
	}

	testutil.TestPacketFields("arp", t, tests)
}
