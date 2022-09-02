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

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ip"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/packettestutil"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/udp"
)

// IP4 header carrying UDP data.
var ip4udp = []byte{0x45, 0x01, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0xff, 0x11, 0xa1, 0xae, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d}

// IP6 header carrying UDP data.
var ip6udp = []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x0c, 0x11, 0x04, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

// UDP header.
var udp = []byte{0x01, 0x02, 0x03, 0x04, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x0b, 0x0c, 0x0d}

func TestUDP(t *testing.T) {
	queries := []packettestutil.FieldQuery{
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_L4_PORT_SRC, 0),
			Result: []byte{0x01, 0x02},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_L4_PORT_DST, 0),
			Result: []byte{0x03, 0x04},
		},
	}
	updates := []packettestutil.FieldUpdate{
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_L4_PORT_SRC, 0),
			Arg: []byte{0x11, 0x12},
			Op:  fwdpacket.OpSet,
		},
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_L4_PORT_DST, 0),
			Arg: []byte{0x13, 0x14},
			Op:  fwdpacket.OpSet,
		},
	}
	tests := []packettestutil.PacketFieldTest{
		// UDP over IP4.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				ip4udp,
				udp,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernetIP4,
				ip4udp,
				[]byte{0x11, 0x12, 0x13, 0x14, 0x00, 0x0c, 0x00, 0x00, 0x0a, 0x0b, 0x0c, 0x0d},
			},
		},
		// UDP over IP6.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				ip6udp,
				udp,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernetIP6,
				ip6udp,
				[]byte{0x11, 0x12, 0x13, 0x14, 0x00, 0x0c, 0x7d, 0x50, 0x0a, 0x0b, 0x0c, 0x0d},
			},
		},
	}

	packettestutil.TestPacketFields("udp", t, tests)
}
