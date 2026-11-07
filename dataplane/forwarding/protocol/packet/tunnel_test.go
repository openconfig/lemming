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

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ip"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/packettestutil"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/udp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/vxlan"
)

// IP tunnels are formed via various combinations of IP4, IP6, GRE tunnels.
// All tunnels are tested using recursive tunnels i.e the packet has an outer
// tunnel that transports an inner tunnel. This implies that the packet has
// 3 IP headers; the outer tunnel header, the inner tunnel header and the
// inner tunnel payload.
const (
	outerIP4IP4       = iota // IP4 header with an IP4 payload (outer tunnel)
	outerIP4IP6              // IP4 header with an IP6 payload (outer tunnel)
	outerIP4GRE              // IP4 header with a GRE payload (outer tunnel)
	outerIP4GREKeySeq        // IP4 header with a GRE payload with Key and Sequence (outer tunnel)
	outerIP6GRE              // IP6 header with a GRE payload (outer tunnel)
	GRE4                     // GRE header with an IP4 payload
	GRE6                     // GRE header with an IP6 payload
	midIP4IP4                // IP4 header with an IP4 payload (inner tunnel)
	midIP6IP6                // IP6 header with an IP6 payload (inner tunnel)
	innerIP4                 // IP4 header of the payload
	innerIP6                 // IP6 header of the payload
	AutoIP6                  // IP6 header in a 6to4 auto tunnel
	AutoIP4                  // IP4 header in a 6to4 auto tunnel
	SecureIP6                // IP6 header in a 6to4 secure tunnel
	SecureIP4                // IP4 header in a 6to4 secure tunnel
	GRE4Key                  // GRE header with an IP4 payload, and GRE key
	GRE4Seq                  // GRE header with an IP4 payload, and GRE sequence number
	GRE4KeySeq               // GRE header with an IP4 payload, and GRE key and sequence number
	outerIP4UDP              // IP4 header with UDP payload
	udpHeader                // UDP header
	vxlanHeader              // VXLAN header
	innerEthernet            // Inner Ethernet header
)

// A field describes a field in the header used to create a test packet.
// The description includes the ID of the field and the value of the field
// in the header.
type field struct {
	id    fwdpacket.FieldID // ID of the field.
	value []byte            // Value of the field.
}

// A header describes a header used to create test packets.
type header struct {
	id     fwdpb.PacketHeaderId // ID of the header.
	orig   []byte               // Bytes in the header.
	fields []field              // Fields in the header.
}

// fields is a set of flags for known fields used to create tests.
//
// The update flag indicates that the field can be updated. This is used to
// generate tests to update the field. Note that this can also added extensions
// to the header such as the GRE key and sequence extension.
//
// A field can optionally specify the argument used in an update.
var fields = map[fwdpb.PacketFieldNum]struct {
	update bool   // Indicates if the field can be updated.
	arg    []byte // Optional argument used in an update.
}{
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION: {
		update: false,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC: {
		update: true,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST: {
		update: true,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS: {
		update: true,
		arg:    []byte{0x00, 0x00, 0x00, 0x11},
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO: {
		update: false,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP6_FLOW: {
		update: true,
		arg:    []byte{0x00, 0x00, 0x02, 0x24},
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_GRE_KEY: {
		update: true,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_GRE_SEQUENCE: {
		update: true,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC: {
		update: true,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST: {
		update: true,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI: {
		update: true,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC: {
		update: true,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST: {
		update: true,
	},
	fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE: {
		update: false,
	},
}

// headers is a set of hand crafted packet headers used to create test packets.
var headers = map[int]header{
	outerIP4IP4: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
		orig: []byte{0x45, 0x01, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xa0, 0x9e, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				value: []byte{0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0),
				value: []byte{0x01, 0x02, 0x03, 0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0),
				value: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 0),
				value: []byte{0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 0),
				value: []byte{0x04},
			},
		},
	},
	midIP4IP4: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
		orig: []byte{0x45, 0x01, 0x00, 0x2a, 0x00, 0x00, 0x00, 0x00, 0x04, 0x04, 0x9c, 0xb2, 0x0a, 0x0b, 0x0c, 0x0d, 0x01, 0x02, 0x03, 0x04},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 1),
				value: []byte{0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 1),
				value: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 1),
				value: []byte{0x01, 0x02, 0x03, 0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 1),
				value: []byte{0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 1),
				value: []byte{0x04},
			},
		},
	},
	innerIP4: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
		orig: []byte{0x45, 0x01, 0x00, 0x16, 0x00, 0x00, 0x00, 0x00, 0x08, 0xff, 0x97, 0xcb, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d, 0x00, 0x00},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 2),
				value: []byte{0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 2),
				value: []byte{0x01, 0x02, 0x03, 0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 2),
				value: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 2),
				value: []byte{0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 2),
				value: []byte{0xff},
			},
		},
	},
	outerIP4IP6: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
		orig: []byte{0x45, 0x01, 0x00, 0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x29, 0xa0, 0x51, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				value: []byte{0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0),
				value: []byte{0x01, 0x02, 0x03, 0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0),
				value: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 0),
				value: []byte{0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 0),
				value: []byte{0x29},
			},
		},
	},
	midIP6IP6: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
		orig: []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x2a, 0x29, 0x04, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 1),
				value: []byte{0x06},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 1),
				value: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 1),
				value: []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 1),
				value: []byte{0x00, 0x00, 0x00, 0x10},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 1),
				value: []byte{0x29},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP6_FLOW, 1),
				value: []byte{0x00, 0x00, 0x02, 0x00},
			},
		},
	},
	innerIP6: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
		orig: []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x02, 0xff, 0x04, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x01, 0x02},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 2),
				value: []byte{0x06},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 2),
				value: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 2),
				value: []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 2),
				value: []byte{0x00, 0x00, 0x00, 0x10},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 2),
				value: []byte{0xff},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP6_FLOW, 2),
				value: []byte{0x00, 0x00, 0x02, 0x00},
			},
		},
	},
	outerIP4GRE: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
		orig: []byte{0x45, 0x01, 0x00, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2f, 0xa0, 0x6f, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				value: []byte{0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0),
				value: []byte{0x01, 0x02, 0x03, 0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0),
				value: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 0),
				value: []byte{0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 0),
				value: []byte{0x2F},
			},
		},
	},
	outerIP4GREKeySeq: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
		orig: []byte{0x45, 0x01, 0x00, 0x4a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2f, 0xa0, 0x67, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				value: []byte{0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0),
				value: []byte{0x01, 0x02, 0x03, 0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0),
				value: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 0),
				value: []byte{0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 0),
				value: []byte{0x2F},
			},
		},
	},
	// GRE header with IPv4 payload.
	GRE4: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
		orig: []byte{0x00, 0x00, 0x08, 0x00},
	},
	// GRE header with sequence number.
	GRE4Seq: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
		orig: []byte{0x10, 0x00, 0x08, 0x00, 0x04, 0x03, 0x02, 0x01},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_GRE_SEQUENCE, 0),
				value: []byte{0x04, 0x03, 0x02, 0x01},
			},
		},
	},
	// GRE header with key number.
	GRE4Key: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
		orig: []byte{0x20, 0x00, 0x08, 0x00, 0x04, 0x03, 0x02, 0x01},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_GRE_KEY, 0),
				value: []byte{0x04, 0x03, 0x02, 0x01},
			},
		},
	},
	// GRE header with key and sequence number.
	GRE4KeySeq: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
		orig: []byte{0x30, 0x00, 0x08, 0x00, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_GRE_KEY, 0),
				value: []byte{0x04, 0x03, 0x02, 0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_GRE_SEQUENCE, 0),
				value: []byte{0x01, 0x02, 0x03, 0x04},
			},
		},
	},
	outerIP6GRE: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
		orig: []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x56, 0x2f, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				value: []byte{0x06},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0),
				value: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0),
				value: []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 0),
				value: []byte{0x00, 0x00, 0x00, 0x10},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 0),
				value: []byte{0x2F},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP6_FLOW, 0),
				value: []byte{0x00, 0x00, 0x02, 0x00},
			},
		},
	},
	GRE6: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
		orig: []byte{0x00, 0x00, 0x86, 0xDD},
	},
	AutoIP6: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
		orig: []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x02, 0xff, 0x04, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x20, 0x02, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 2),
				value: []byte{0x06},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 2),
				value: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 2),
				value: []byte{0x20, 0x02, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 2),
				value: []byte{0x00, 0x00, 0x00, 0x10},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 2),
				value: []byte{0xff},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP6_FLOW, 2),
				value: []byte{0x00, 0x00, 0x02, 0x00},
			},
		},
	},
	// The destination address is encoded in the frame for comparison, but not
	// specified in the fields as it is auto-computed.
	AutoIP4: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO,
		orig: []byte{0x45, 0x01, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x29, 0xb0, 0x8d, 0x01, 0x02, 0x03, 0x04, 0x04, 0x03, 0x02, 0x01},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0),
				value: []byte{0x01, 0x02, 0x03, 0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 0),
				value: []byte{0x01},
			},
		},
	},
	SecureIP6: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
		orig: []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x02, 0xff, 0x04, 0x20, 0x02, 0x01, 0x02, 0x03, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x02, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 2),
				value: []byte{0x06},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 2),
				value: []byte{0x20, 0x02, 0x01, 0x02, 0x03, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 2),
				value: []byte{0x20, 0x02, 0x04, 0x03, 0x02, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 2),
				value: []byte{0x00, 0x00, 0x00, 0x10},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 2),
				value: []byte{0xff},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP6_FLOW, 2),
				value: []byte{0x00, 0x00, 0x02, 0x00},
			},
		},
	},
	// The source and destination address is encoded in the frame for comparison,
	// but not specified in the fields as they are auto-computed.
	SecureIP4: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE,
		orig: []byte{0x45, 0x01, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x29, 0xb0, 0x8d, 0x01, 0x02, 0x03, 0x04, 0x04, 0x03, 0x02, 0x01},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 0),
				value: []byte{0x01},
			},
		},
	},
	outerIP4UDP: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
		orig: []byte{0x45, 0x01, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0xa0, 0x8e, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
				value: []byte{0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0),
				value: []byte{0x01, 0x02, 0x03, 0x04},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0),
				value: []byte{0x0a, 0x0b, 0x0c, 0x0d},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, 0),
				value: []byte{0x01},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 0),
				value: []byte{0x11},
			},
		},
	},
	udpHeader: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP,
		orig: []byte{0x01, 0x02, 0x12, 0xb5, 0x00, 0x2a, 0x00, 0x00},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC, 0),
				value: []byte{0x01, 0x02},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST, 0),
				value: []byte{0x12, 0xb5},
			},
		},
	},
	vxlanHeader: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN,
		orig: []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI, 0),
				value: []byte{0x00, 0x00, 0x01},
			},
		},
	},
	innerEthernet: {
		id:   fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
		orig: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x08, 0x00},
		fields: []field{
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
				value: []byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
				value: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
			},
			{
				id:    fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0),
				value: []byte{0x08, 0x00},
			},
		},
	},
}

// newValue creates a new value for a given field using its current value.
// If the field has a preconfigured arg, it used. Otherwise a new value
// is generated by incrementing each byte.
//
// Note that the newValue is always the same size as the current value. This
// allows us to generate values for fields like QoS which have different sizes
// depending on the packet type.
func newValue(arg, cur []byte) []byte {
	if arg == nil {
		var upd []byte
		for _, b := range cur {
			upd = append(upd, b+1)
		}
		return upd
	}
	if len(arg) >= len(cur) {
		return arg[0:len(cur)]
	}
	upd := make([]byte, len(cur))
	copy(upd, arg)
	return upd
}

// TestTunnelFields performs query and update tests for various fields
// in IP tunnel packets.
//
// The test has a description of packets that are used to generate unit tests.
// For each described packet, it does the following:
// 1. Generate the packet frame by stringing the ethernet frame with all the
// specified headers.
// 2. Generate a query test for all fields present in the packet.
// 3. Generate an update test for all fields that can be updated.
func TestTunnelParsing(t *testing.T) {
	descs := []struct {
		text     string // Description of the test.
		headers  []int  // List of headers in the packet (in sequence).
		ethernet []byte // Ethernet header for the packet.
	}{
		{
			text:     "IP4-IP4-IP4",
			headers:  []int{outerIP4IP4, midIP4IP4, innerIP4},
			ethernet: ethernetIP4,
		},
		{
			text:     "IP4-IP6-IP6",
			headers:  []int{outerIP4IP6, midIP6IP6, innerIP6},
			ethernet: ethernetIP4,
		},
		{
			text:     "IP4-GRE-IP4-IP4",
			headers:  []int{outerIP4GRE, GRE4, midIP4IP4, innerIP4},
			ethernet: ethernetIP4,
		},
		{
			text:     "IP4-GRE(Key)-IP4-IP4",
			headers:  []int{outerIP4GRE, GRE4Key, midIP4IP4, innerIP4},
			ethernet: ethernetIP4,
		},
		{
			text:     "IP4-GRE(Seq)-IP4-IP4",
			headers:  []int{outerIP4GRE, GRE4Seq, midIP4IP4, innerIP4},
			ethernet: ethernetIP4,
		},
		{
			text:     "IP4-GRE(KeySeq)-IP4-IP4",
			headers:  []int{outerIP4GRE, GRE4KeySeq, midIP4IP4, innerIP4},
			ethernet: ethernetIP4,
		},
		{
			text:     "IP6-GRE-IP6-IP6",
			headers:  []int{outerIP6GRE, GRE6, midIP6IP6, innerIP6},
			ethernet: ethernetIP6,
		},
		{
			text:     "IP6-GRE-IP6-IP6",
			headers:  []int{outerIP6GRE, GRE6, midIP6IP6, innerIP6},
			ethernet: nil, // parsing a tunnel as IP
		},
		{
			text:     "IP4-IP4-IP4",
			headers:  []int{outerIP4IP4, midIP4IP4, innerIP4},
			ethernet: nil, // parsing a tunnel as IP
		},
	}

	for _, desc := range descs {
		var queries []packettestutil.FieldQuery
		var updates []packettestutil.FieldUpdate
		frame := [][]byte{desc.ethernet}
		start := fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET
		if len(desc.ethernet) == 0 {
			start = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP
		}

		for _, header := range desc.headers {
			h := headers[header]

			frame = append(frame, h.orig)
			for _, field := range h.fields {
				attr := fields[field.id.Num]
				queries = append(queries, packettestutil.FieldQuery{
					ID:     field.id,
					Result: field.value,
				})
				if attr.update {
					updates = append(updates, packettestutil.FieldUpdate{
						ID:  field.id,
						Op:  fwdpacket.OpSet,
						Arg: newValue(attr.arg, field.value),
					})
				}
			}
		}
		test := packettestutil.PacketFieldTest{
			StartHeader: start,
			Orig:        frame,
			Updates:     updates,
			Queries:     queries,
		}
		packettestutil.TestPacketFields(desc.text, t, []packettestutil.PacketFieldTest{test})
	}
}

// TestTunnelDecap tests decap for various IP tunnel packets.
//
// The test has a description of packets that are used to generate unit tests.
// For each described packet, it does the following:
// 1. Generate the packet frame by stringing the ethernet frame with all the
// specified headers.
// 2. Generate a decap test for the specified headers.
// 3. Verify the packet frame after all decap are completed.
func TestTunnelDecap(t *testing.T) {
	descs := []struct {
		text          string // Description of the packet.
		headers       []int  // List of headers in the original packet (in sequence).
		origEthernet  []byte // Original ethernet header.
		finalEthernet []byte // Final ethernet header.
		depth         int    // Number of headers to decap.
	}{
		{
			text:          "IP4-IP4-IP4",
			headers:       []int{outerIP4IP4, midIP4IP4, innerIP4},
			origEthernet:  ethernetIP4,
			finalEthernet: ethernetIP4,
			depth:         1,
		},
		{
			text:          "IP4-IP6-IP6",
			headers:       []int{outerIP4IP6, midIP6IP6, innerIP6},
			origEthernet:  ethernetIP4,
			finalEthernet: ethernetIP6,
			depth:         1,
		},
		{
			text:          "IP4-GRE-IP4-IP4",
			headers:       []int{outerIP4GRE, GRE4, midIP4IP4, innerIP4},
			origEthernet:  ethernetIP4,
			finalEthernet: ethernetIP4,
			depth:         2,
		},
		{
			text:          "IP4-GRE(KeySeq)-IP4-IP4",
			headers:       []int{outerIP4GREKeySeq, GRE4KeySeq, midIP4IP4, innerIP4},
			origEthernet:  ethernetIP4,
			finalEthernet: ethernetIP4,
			depth:         2,
		},
		{
			text:          "IP6-GRE-IP6-IP6",
			headers:       []int{outerIP6GRE, GRE6, midIP6IP6, innerIP6},
			origEthernet:  ethernetIP6,
			finalEthernet: ethernetIP6,
			depth:         2,
		},
		{
			text:          "6to4-AUTO",
			headers:       []int{AutoIP4, AutoIP6},
			origEthernet:  ethernetIP4,
			finalEthernet: ethernetIP6,
			depth:         1,
		},
		{
			text:          "6to4-SECURE",
			headers:       []int{SecureIP4, SecureIP6},
			origEthernet:  ethernetIP4,
			finalEthernet: ethernetIP6,
			depth:         1,
		},
	}

	for _, desc := range descs {
		var ids []fwdpb.PacketHeaderId        // List of headers to decap.
		orig := [][]byte{desc.origEthernet}   // Original frame (with all headers).
		final := [][]byte{desc.finalEthernet} // Final frame.

		for pos, header := range desc.headers {
			h := headers[header]
			orig = append(orig, h.orig)
			if pos < desc.depth {
				ids = append(ids, h.id)
			} else {
				final = append(final, h.orig)
			}
		}

		var updates []packettestutil.HeaderUpdate
		for _, id := range ids {
			updates = append(updates, packettestutil.HeaderUpdate{
				ID:    id,
				Encap: false,
			})
		}
		updates[len(updates)-1].Result = final
		test := packettestutil.PacketHeaderTest{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig:        orig,
			Updates:     updates,
		}
		packettestutil.TestPacketHeaders(desc.text, t, []packettestutil.PacketHeaderTest{test})
	}
}

// TestTunnelEncap tests encap for various IP tunnel packets.
//
// The test has a description of packets that are used to generate unit tests.
// For each described packet, it does the following:
// 1. Generate the packet frame by stringing the ethernet frame with all the
// specified headers.
// 2. Generate a encap test for the specified headers.
// 3. Add fields that are present in the newly added headers.
// 4. Verify the packet frame after the encap and updates are done.
func TestTunnelEncap(t *testing.T) {
	descs := []struct {
		text          string // Description of the packet
		headers       []int  // List of headers in the final frame (in sequence)
		origEthernet  []byte // Original ethernet header
		finalEthernet []byte // Final ethernet header
		depth         int    // Number of headers to encap
	}{
		{
			text:          "IP4-IP4-IP4",
			headers:       []int{outerIP4IP4, midIP4IP4, innerIP4},
			origEthernet:  ethernetIP4,
			finalEthernet: ethernetIP4,
			depth:         1,
		},
		{
			text:          "IP4-IP6-IP6",
			headers:       []int{outerIP4IP6, midIP6IP6, innerIP6},
			origEthernet:  ethernetIP6,
			finalEthernet: ethernetIP4,
			depth:         1,
		},
		{
			text:          "IP4-GRE-IP4-IP4",
			headers:       []int{outerIP4GRE, GRE4, midIP4IP4, innerIP4},
			origEthernet:  ethernetIP4,
			finalEthernet: ethernetIP4,
			depth:         2,
		},
		{
			text:          "IP4-GRE(KeySeq)-IP4-IP4",
			headers:       []int{outerIP4GREKeySeq, GRE4KeySeq, midIP4IP4, innerIP4},
			origEthernet:  ethernetIP4,
			finalEthernet: ethernetIP4,
			depth:         2,
		},
		{
			text:          "IP6-GRE-IP6-IP6",
			headers:       []int{outerIP6GRE, GRE6, midIP6IP6, innerIP6},
			origEthernet:  ethernetIP6,
			finalEthernet: ethernetIP6,
			depth:         2,
		},
		{
			text:          "6to4-AUTO",
			headers:       []int{AutoIP4, AutoIP6},
			origEthernet:  ethernetIP6,
			finalEthernet: ethernetIP4,
			depth:         1,
		},
		{
			text:          "6to4-SECURE",
			headers:       []int{SecureIP4, SecureIP6},
			origEthernet:  ethernetIP6,
			finalEthernet: ethernetIP4,
			depth:         1,
		},
	}

	for _, desc := range descs {
		var fupd []packettestutil.FieldUpdate // Field updates for the new headers.
		var ids []fwdpb.PacketHeaderId        // List of headers to add (reverse order).
		orig := [][]byte{desc.origEthernet}   // Original frame.
		final := [][]byte{desc.finalEthernet} // Final frame.

		for pos, header := range desc.headers {
			h := headers[header]
			final = append(final, h.orig)
			if pos >= desc.depth {
				orig = append(orig, h.orig)
			} else {
				ids = append([]fwdpb.PacketHeaderId{h.id}, ids...)
				for _, f := range h.fields {
					fupd = append(fupd, packettestutil.FieldUpdate{
						ID:  f.id,
						Op:  fwdpacket.OpSet,
						Arg: f.value,
					})
				}
			}
		}

		var updates []packettestutil.HeaderUpdate
		for _, id := range ids {
			updates = append(updates, packettestutil.HeaderUpdate{
				ID:    id,
				Encap: true,
			})
		}
		updates[len(updates)-1].Result = final
		updates[len(updates)-1].Updates = fupd
		test := packettestutil.PacketHeaderTest{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig:        orig,
			Updates:     updates,
		}
		packettestutil.TestPacketHeaders(desc.text, t, []packettestutil.PacketHeaderTest{test})
	}
}

// TestTunnelDecapErrors tests tunnel decap
func TestTunnelDecapErrors(t *testing.T) {
	tests := []packettestutil.PacketHeaderTest{
		// Decap IP4 frame.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[innerIP4].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE,
					Err:   "failed",
				},
			},
		},
		// Decap IP4-IP4 tunnel.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[outerIP4IP4].orig,
				headers[midIP4IP4].orig,
				headers[innerIP4].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
					Result: [][]byte{
						ethernetIP4,
						headers[midIP4IP4].orig,
						headers[innerIP4].orig,
					},
				},
			},
		},
		// Decap IP4-IP6 tunnel.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[outerIP4IP6].orig,
				headers[midIP6IP6].orig,
				headers[innerIP6].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
					Result: [][]byte{
						ethernetIP6,
						headers[midIP6IP6].orig,
						headers[innerIP6].orig,
					},
				},
			},
		},
		// Decap IP4-GRE-IP4 tunnel.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[outerIP4GRE].orig,
				headers[GRE4].orig,
				headers[innerIP4].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE,
					Err:   "failed",
				},
				{
					// The IPv4 header must be decapped before GRE.
					// The result is not verified because this will
					// be an invalid ethernet frame until the GRE
					// header is stripped as well.
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
					Result: [][]byte{
						ethernetIP4,
						headers[innerIP4].orig,
					},
				},
			},
		},
		// Decap IP4-GRE(KeySeq)-IP4 tunnel.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[outerIP4GREKeySeq].orig,
				headers[GRE4KeySeq].orig,
				headers[innerIP4].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE,
					Err:   "failed",
				},
				{
					// The IPv4 header must be decapped before GRE.
					// The result is not verified because this will
					// be an invalid ethernet frame until the GRE
					// header is stripped as well.
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
					Result: [][]byte{
						ethernetIP4,
						headers[innerIP4].orig,
					},
				},
			},
		},
		// Decap 6to4 auto tunnel.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[AutoIP4].orig,
				headers[AutoIP6].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO,
					Result: [][]byte{
						ethernetIP6,
						headers[AutoIP6].orig,
					},
				},
			},
		},
		// Decap 6to4 auto tunnel using IP4
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[AutoIP4].orig,
				headers[AutoIP6].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
					Result: [][]byte{
						ethernetIP6,
						headers[AutoIP6].orig,
					},
				},
			},
		},
		// Decap 6to4 secure tunnel.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[SecureIP4].orig,
				headers[SecureIP6].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
					Err:   "failed",
				},
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE,
					Result: [][]byte{
						ethernetIP6,
						headers[SecureIP6].orig,
					},
				},
			},
		},
		// Decap 6to4 secure tunnel using 6to4 auto.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[SecureIP4].orig,
				headers[SecureIP6].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO,
					Result: [][]byte{
						ethernetIP6,
						headers[SecureIP6].orig,
					},
				},
			},
		},
		// Decap 6to4 secure tunnel using IP4
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[SecureIP4].orig,
				headers[SecureIP6].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: false,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
					Result: [][]byte{
						ethernetIP6,
						headers[SecureIP6].orig,
					},
				},
			},
		},
	}
	packettestutil.TestPacketHeaders("tunnel-decap-errors", t, tests)
}

// TestTunnelEncapErrors tests errors during tunnel encap
func TestTunnelEncapErrors(t *testing.T) {
	tests := []packettestutil.PacketHeaderTest{
		// Encap tunnel headers on an arp frame.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				arp,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: true,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
					Err:   "failed",
				},
				{
					Encap: true,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO,
					Err:   "failed",
				},
				{
					Encap: true,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE,
					Err:   "failed",
				},
			},
		},
		// Encap 6TO4 tunnels to a IP4 packet.
		{
			StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				headers[innerIP4].orig,
			},
			Updates: []packettestutil.HeaderUpdate{
				{
					Encap: true,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO,
					Err:   "failed",
				},
				{
					Encap: true,
					ID:    fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE,
					Err:   "failed",
				},
			},
		},
	}
	packettestutil.TestPacketHeaders("tunnel-encap-errors", t, tests)
}

func TestVXLANPacketParsing(t *testing.T) {
	outerEthernet := []byte{
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, // dst mac
		0x00, 0x11, 0x22, 0x33, 0x44, 0x66, // src mac
		0x08, 0x00, // ethertype ipv4
	}
	outerIP := []byte{
		0x45, 0x01, 0x00, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x40, 0x11, 0xa0, 0x8e,
		0x0a, 0x00, 0x01, 0x01, // src ip: 10.0.1.1
		0x0a, 0x00, 0x01, 0x02, // dst ip: 10.0.1.2
	}
	udpHeader := []byte{
		0x01, 0x02, // src port
		0x12, 0xb5, // dst port: 4789
		0x00, 0x2a, // length
		0x00, 0x00, // csum
	}
	vxlanHeader := []byte{
		0x08, 0x00, 0x00, 0x00, // flags
		0x00, 0x4e, 0x21, // vni: 20001
		0x00, // reserved
	}
	innerEthernet := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, // dst mac
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xfe, // src mac
		0x08, 0x00, // ethertype ipv4
	}
	innerIP := []byte{
		0x45, 0x01, 0x00, 0x16, 0x00, 0x00, 0x00, 0x00, 0x40, 0xff, 0x97, 0xcb,
		0xc0, 0xa8, 0x01, 0x64, // src ip: 192.168.1.100
		0xc0, 0xa8, 0x01, 0xc8, // dst ip: 192.168.1.200
		0x00, 0x00,
	}

	rawBytes := append([]byte(nil), outerEthernet...)
	rawBytes = append(rawBytes, outerIP...)
	rawBytes = append(rawBytes, udpHeader...)
	rawBytes = append(rawBytes, vxlanHeader...)
	rawBytes = append(rawBytes, innerEthernet...)
	rawBytes = append(rawBytes, innerIP...)

	packet, err := protocol.NewPacket(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, frame.NewFrame(rawBytes))
	if err != nil {
		t.Fatalf("NewPacket failed: %v", err)
	}

	tests := []struct {
		name    string
		fieldID fwdpacket.FieldID
		want    []byte
	}{
		{
			name:    "Outer Ethernet MAC Dst",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
			want:    []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		},
		{
			name:    "Outer Ethernet MAC Src",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0),
			want:    []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x66},
		},
		{
			name:    "Outer IP Src",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0),
			want:    []byte{0x0a, 0x00, 0x01, 0x01},
		},
		{
			name:    "Outer IP Dst",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0),
			want:    []byte{0x0a, 0x00, 0x01, 0x02},
		},
		{
			name:    "Outer UDP Port Src",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC, 0),
			want:    []byte{0x01, 0x02},
		},
		{
			name:    "Outer UDP Port Dst (4789)",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST, 0),
			want:    []byte{0x12, 0xb5},
		},
		{
			name:    "VXLAN VNI (20001)",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI, 0),
			want:    []byte{0x00, 0x4e, 0x21},
		},
		{
			name:    "Inner Ethernet MAC Dst",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 1),
			want:    []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		},
		{
			name:    "Inner Ethernet MAC Src",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 1),
			want:    []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xfe},
		},
		{
			name:    "Inner IP Src",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 1),
			want:    []byte{0xc0, 0xa8, 0x01, 0x64},
		},
		{
			name:    "Inner IP Dst",
			fieldID: fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 1),
			want:    []byte{0xc0, 0xa8, 0x01, 0xc8},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := packet.Field(tc.fieldID)
			if err != nil {
				t.Fatalf("Field(%v) failed: %v", tc.fieldID, err)
			}
			if diff := cmp.Diff(b, tc.want); diff != "" {
				t.Errorf("Field(%v) diff (-got +want):\n%s", tc.fieldID, diff)
			}
		})
	}
}

func TestVXLANPacketEncapDecap(t *testing.T) {
	innerEthernet := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, // dst mac
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xfe, // src mac
		0x08, 0x00, // ethertype ipv4
	}
	innerIP := []byte{
		0x45, 0x01, 0x00, 0x16, 0x00, 0x00, 0x00, 0x00, 0x40, 0xff, 0xf5, 0x6b,
		0xc0, 0xa8, 0x01, 0x64, // src ip: 192.168.1.100
		0xc0, 0xa8, 0x01, 0xc8, // dst ip: 192.168.1.200
		0x00, 0x00,
	}

	origFrame := append(append([]byte(nil), innerEthernet...), innerIP...)
	p, err := protocol.NewPacket(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, frame.NewFrame(append([]byte(nil), origFrame...)))
	if err != nil {
		t.Fatalf("NewPacket failed: %v", err)
	}

	// Perform VXLAN encap sequence matching SAI NextHop + TunnelEncap pipeline:
	// 1. Encap VXLAN, UDP, outer IP4
	for _, hid := range []fwdpb.PacketHeaderId{
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN,
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP,
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
	} {
		if err := p.Encap(hid); err != nil {
			t.Fatalf("Encap(%v) failed: %v", hid, err)
		}
	}

	// 2. Set VNI, inner ETHER_MAC_SRC (from inner ETHER_MAC_DST), inner ETHER_MAC_DST (TunnelMac)
	if err := p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI, 0), fwdpacket.OpSet, []byte{0x00, 0x4e, 0x21}); err != nil {
		t.Fatalf("Update(VXLAN_VNI) failed: %v", err)
	}
	innerDstMAC, err := p.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0))
	if err != nil {
		t.Fatalf("Field(ETHER_MAC_DST) failed: %v", err)
	}
	if err := p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0), fwdpacket.OpSet, innerDstMAC); err != nil {
		t.Fatalf("Update(inner ETHER_MAC_SRC) failed: %v", err)
	}
	if err := p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0), fwdpacket.OpSet, []byte{0x00, 0xaa, 0xbb, 0xcc, 0xdd, 0xee}); err != nil {
		t.Fatalf("Update(inner ETHER_MAC_DST) failed: %v", err)
	}

	// 3. Simulate outputTable: Decap(ETHERNET) when no outer L2 header exists
	// must fail without removing Inner Ethernet, and Encap(ETHERNET) then adds
	// the outer Ethernet header.
	if err := p.Decap(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET); err == nil {
		t.Fatalf("Decap(ETHERNET) without outer L2 header succeeded unexpectedly")
	}
	if err := p.Encap(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET); err != nil {
		t.Fatalf("Encap(ETHERNET) failed: %v", err)
	}

	// 4. Populate outer IP, UDP, and Ethernet fields
	if err := p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0), fwdpacket.OpSet, []byte{10, 0, 1, 1}); err != nil {
		t.Fatalf("Update(outer IP_ADDR_SRC) failed: %v", err)
	}
	if err := p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0), fwdpacket.OpSet, []byte{10, 0, 1, 2}); err != nil {
		t.Fatalf("Update(outer IP_ADDR_DST) failed: %v", err)
	}
	innerTTL, err := p.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP, 1))
	if err != nil {
		t.Fatalf("Field(inner IP_HOP) failed: %v", err)
	}
	if err := p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP, 0), fwdpacket.OpSet, innerTTL); err != nil {
		t.Fatalf("Update(outer IP_HOP) failed: %v", err)
	}
	if err := p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST, 0), fwdpacket.OpSet, []byte{0x12, 0xb5}); err != nil {
		t.Fatalf("Update(L4_PORT_DST) failed: %v", err)
	}
	if err := p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0), fwdpacket.OpSet, []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x66}); err != nil {
		t.Fatalf("Update(outer ETHER_MAC_SRC) failed: %v", err)
	}
	if err := p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0), fwdpacket.OpSet, []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}); err != nil {
		t.Fatalf("Update(outer ETHER_MAC_DST) failed: %v", err)
	}

	// Query UDP source port before Frame() to verify on-demand entropy computation.
	sportBeforeRebuild, err := p.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC, 0))
	if err != nil {
		t.Fatalf("Field(L4_PORT_SRC) before Frame() failed: %v", err)
	}
	sportVal := int(sportBeforeRebuild[0])<<8 | int(sportBeforeRebuild[1])
	if sportVal < 49152 || sportVal > 65535 {
		t.Errorf("entropy UDP source port before Frame() = %d, want in [49152, 65535]", sportVal)
	}

	// Rebuild frame and reparse to verify wire format.
	encapBytes := p.Frame()
	reparsed, err := protocol.NewPacket(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, frame.NewFrame(append([]byte(nil), encapBytes...)))
	if err != nil {
		t.Fatalf("NewPacket(encapBytes) failed: %v", err)
	}

	wantFields := []struct {
		name string
		id   fwdpacket.FieldID
		want []byte
	}{
		{"Outer MAC Dst", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0), []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}},
		{"Outer MAC Src", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0), []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x66}},
		{"Outer EtherType", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE, 0), []byte{0x08, 0x00}},
		{"Outer IP Src", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0), []byte{10, 0, 1, 1}},
		{"Outer IP Dst", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0), []byte{10, 0, 1, 2}},
		{"Outer IP TTL", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP, 0), []byte{0x40}},
		{"Outer IP Proto", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO, 0), []byte{17}},
		{"Outer UDP Sport", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC, 0), sportBeforeRebuild},
		{"Outer UDP Dport", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST, 0), []byte{0x12, 0xb5}},
		{"VXLAN VNI", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI, 0), []byte{0x00, 0x4e, 0x21}},
		{"Inner MAC Dst", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 1), []byte{0x00, 0xaa, 0xbb, 0xcc, 0xdd, 0xee}},
		{"Inner MAC Src", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 1), []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
		{"Inner IP Src", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 1), []byte{192, 168, 1, 100}},
		{"Inner IP Dst", fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 1), []byte{192, 168, 1, 200}},
	}
	for _, tc := range wantFields {
		got, err := reparsed.Field(tc.id)
		if err != nil {
			t.Fatalf("%s: Field(%v) failed: %v", tc.name, tc.id, err)
		}
		if diff := cmp.Diff(got, tc.want); diff != "" {
			t.Errorf("%s: diff (-got +want):\n%s", tc.name, diff)
		}
	}

	// Now test VXLAN Decap on reparsed packet (matching SAI termination + L2 rewrite).
	for _, hid := range []fwdpb.PacketHeaderId{
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP,
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN,
	} {
		if err := reparsed.Decap(hid); err != nil {
			t.Fatalf("Decap(%v) failed: %v", hid, err)
		}
	}

	wantDecap := append([]byte(nil),
		0x00, 0xaa, 0xbb, 0xcc, 0xdd, 0xee,
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
		0x08, 0x00,
	)
	wantDecap = append(wantDecap, innerIP...)
	if diff := cmp.Diff(reparsed.Frame(), wantDecap); diff != "" {
		t.Errorf("Decap frame diff (-got +want):\n%s", diff)
	}

	// Also verify IPv6 VXLAN outer encapsulation computes UDP checksum over outer IPv6 addresses.
	p6, err := protocol.NewPacket(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, frame.NewFrame(append([]byte(nil), origFrame...)))
	if err != nil {
		t.Fatalf("NewPacket(p6) failed: %v", err)
	}
	for _, hid := range []fwdpb.PacketHeaderId{
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN,
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP,
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
		fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
	} {
		if err := p6.Encap(hid); err != nil {
			t.Fatalf("p6.Encap(%v) failed: %v", hid, err)
		}
	}
	if err := p6.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, 0), fwdpacket.OpSet, []byte{0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}); err != nil {
		t.Fatalf("p6.Update(IP_ADDR_SRC) failed: %v", err)
	}
	if err := p6.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, 0), fwdpacket.OpSet, []byte{0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}); err != nil {
		t.Fatalf("p6.Update(IP_ADDR_DST) failed: %v", err)
	}
	if err := p6.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST, 0), fwdpacket.OpSet, []byte{0x12, 0xb5}); err != nil {
		t.Fatalf("p6.Update(L4_PORT_DST) failed: %v", err)
	}
	f6 := p6.Frame()
	// Outer Ethernet (14B) + Outer IPv6 (40B) + UDP (8B: checksum at offset 54+6=60..62)
	if len(f6) < 62 || (f6[60] == 0 && f6[61] == 0) {
		t.Errorf("IPv6 VXLAN outer UDP checksum was not computed (got %x)", f6[60:62])
	}
}

