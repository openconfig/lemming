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

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ip"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// IP6 header carrying opaque data.
var ip6 = []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x02, 0xff, 0x04, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x01, 0x02}

func TestIP6(t *testing.T) {
	queries := []testutil.FieldQuery{
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_VERSION, 0),
			Result: []byte{0x06},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_ADDR_SRC, 0),
			Result: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_ADDR_DST, 0),
			Result: []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_QOS, 0),
			Result: []byte{0x00, 0x00, 0x00, 0x10},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_HOP, 0),
			Result: []byte{0x04},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_PROTO, 0),
			Result: []byte{0xff},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP6_FLOW, 0),
			Result: []byte{0x00, 0x00, 0x02, 0x00},
		},
	}

	updates := []testutil.FieldUpdate{
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_ADDR_SRC, 0),
			Arg: []byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x18, 0x17, 0x16, 0x15, 0x14, 0x13, 0x12, 0x11},
			Op:  fwdpacket.OpSet,
		},
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_ADDR_DST, 0),
			Arg: []byte{0x18, 0x17, 0x16, 0x15, 0x14, 0x13, 0x12, 0x11, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18},
			Op:  fwdpacket.OpSet,
		},
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_QOS, 0),
			Arg: []byte{0x00, 0x00, 0x00, 0x20},
			Op:  fwdpacket.OpSet,
		},
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_HOP, 0),
			Arg: []byte{0x01},
			Op:  fwdpacket.OpDec,
		},
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP6_FLOW, 0),
			Arg: []byte{0x00, 0x00, 0x00, 0x04},
			Op:  fwdpacket.OpSet,
		},
	}

	ipFinal := []byte{0x62, 0x00, 0x00, 0x04, 0x00, 0x02, 0xff, 0x03, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x18, 0x17, 0x16, 0x15, 0x14, 0x13, 0x12, 0x11, 0x18, 0x17, 0x16, 0x15, 0x14, 0x13, 0x12, 0x11, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x01, 0x02}

	tests := []testutil.PacketFieldTest{
		// IP.
		{
			StartHeader: fwdpb.PacketHeaderId_IP,
			Orig: [][]byte{
				ip6,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ipFinal,
			},
		},
		// IP6.
		{
			StartHeader: fwdpb.PacketHeaderId_IP6,
			Orig: [][]byte{
				ip6,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ipFinal,
			},
		},
		// IP6 over Ethernet frame.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				ip6,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernetIP6,
				ipFinal,
			},
		},
		// IP6 over VLAN frame.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetVLANIP6,
				ip6,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernetVLANIP6,
				ipFinal,
			},
		},
		// IP6 over 1Q frame.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernet1QIP6,
				ip6,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernet1QIP6,
				ipFinal,
			},
		},
	}
	testutil.TestPacketFields("ipv6", t, tests)
}

func TestIP6TTL(t *testing.T) {
	ethernetIP6 := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x86, 0xDD}

	// Initial IP6 header carrying TCP data.
	// IP6 header carrying TCP data.
	ip6tcpInitial := []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x18, 0x06, 0x04, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	// TCP header.
	tcpSegment := []byte{0x01, 0x02, 0x03, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x51, 0x34, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x0b, 0x0c, 0x0d}

	updates := []testutil.FieldUpdate{
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_IP_HOP, 0),
			Arg: []byte{0x01},
			Op:  fwdpacket.OpDec,
		},
	}

	// Final IP6 header carrying TCP data.
	ip6tcpFinal := []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x18, 0x06, 0x03, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	tests := []testutil.PacketFieldTest{
		// IP6 over Ethernet frame.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				ip6tcpInitial,
				tcpSegment,
			},
			Updates: updates,
			Final: [][]byte{
				ethernetIP6,
				ip6tcpFinal,
				tcpSegment,
			},
		},
	}
	testutil.TestPacketFields("ipv6", t, tests)
}
