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
// The test has a desciption of packets that are used to generate unit tests.
// For each described packet, it does the following:
// 1. Generate the packet frame by stringing the ethernet frame with all the
// specified headers.
// 2. Generate a query test for all fields present in the packet.
// 3. Generate an update test for all fields that can be updated.
func TestTunnelParsing(t *testing.T) {
	descs := []struct {
		text     string // Desciption of the test.
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
// The test has a desciption of packets that are used to generate unit tests.
// For each described packet, it does the following:
// 1. Generate the packet frame by stringing the ethernet frame with all the
// specified headers.
// 2. Generate a decap test for the specified headers.
// 3. Verify the packet frame after all decap are completed.
func TestTunnelDecap(t *testing.T) {
	descs := []struct {
		text          string // Description of the packet.
		headers       []int  // List of headers in the original packet (in sequence).
		origEthernet  []byte // Orignal ethernet header.
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
// The test has a desciption of packets that are used to generate unit tests.
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
