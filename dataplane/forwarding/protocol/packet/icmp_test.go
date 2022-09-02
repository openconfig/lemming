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
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/icmp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ip"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/packettestutil"
)

// IP4 header carrying ICMP4 echo request message.
var ip4icmp = []byte{0x45, 0x01, 0x00, 0x1e, 0x00, 0x00, 0x00, 0x00, 0xff, 0x01,
	0xa1, 0xc0, 0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b, 0x0c, 0x0d}

// ICMPv4 echo request message. The checksum is set to zero.
var icmp4 = []byte{0x01, 0x02, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}

// IP6 header carrying ICMP6 echo request message.
var ip6icmp = []byte{0x61, 0x00, 0x02, 0x00, 0x00, 0x0a, 0x01, 0x04, 0x01, 0x02,
	0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03,
	0x02, 0x01, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x01, 0x02,
	0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

// ICMP6 echo request packet message. The checksum is set to zero.
var icmp6 = []byte{0x01, 0x02, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}

// ICMP6 neighbour advertisement with a target link-layer address option.
var icmp6NATLL = []byte{
	0x6e, 0x00, 0x00, 0x00, 0x00, 0x20, 0x3a, 0xff, 0xfe, 0x80, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00, 0xff, 0x02,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x01, 0x88, 0x00, 0x9a, 0xbb, 0xa0, 0x00, 0x00, 0x00, 0xfe, 0x80, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00, 0x02,
	0x01, 0xc2, 0x00, 0x54, 0xf5, 0x00, 0x00}

// ICMP6 neighbour advertisement with no options.
var icmp6NA = []byte{
	0x6e, 0x00, 0x00, 0x00, 0x00, 0x18, 0x3a, 0xff, 0xfe, 0x80, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00, 0xff, 0x02,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x01, 0x88, 0x00, 0xb3, 0xba, 0xa0, 0x00, 0x00, 0x00, 0xfe, 0x80, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00}

// ICMP6 neighbour solicitation with a source link-layer address option.
var icmp6NSSLL = []byte{0x6e, 0x00, 0x00, 0x00, 0x00, 0x20, 0x3a, 0xff, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x01, 0xff, 0xf5, 0x00, 0x00, 0x87, 0x00, 0x4f, 0x3d, 0x00,
	0x00, 0x00, 0x00, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0,
	0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00, 0x01, 0x01, 0xc2, 0x00, 0x54,
	0xf5, 0x00, 0x00}

// ICMP6 neighbour solicitation with no options.
var icmp6NS = []byte{0x6e, 0x00, 0x00, 0x00, 0x00, 0x18, 0x3a, 0xff, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x01, 0xff, 0xf5, 0x00, 0x00, 0x87, 0x00, 0x67, 0x3c, 0x00, 0x00,
	0x00, 0x00, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00,
	0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00}

// ICMP6 router advertisement with a source link-layer address option.
var icmp6RASLL = []byte{
	0x6e, 0x00, 0x00, 0x00, 0x00, 0x40, 0x3a, 0xff, 0xfe, 0x80, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00,
	0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x01, 0x86, 0x00, 0xc4, 0xfe, 0x40, 0x00, 0x07, 0x08,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0xc2, 0x00,
	0x54, 0xf5, 0x00, 0x00, 0x05, 0x01, 0x00, 0x00, 0x00, 0x00, 0x05, 0xdc,
	0x03, 0x04, 0x40, 0xc0, 0x00, 0x27, 0x8d, 0x00, 0x00, 0x09, 0x3a, 0x80,
	0x00, 0x00, 0x00, 0x00, 0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x01,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

// ICMP6 router advertisement with no options.
var icmp6RA = []byte{
	0x6e, 0x00, 0x00, 0x00, 0x00, 0x10, 0x3a, 0xff, 0xfe, 0x80, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00,
	0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x01, 0x86, 0x00, 0x21, 0x32, 0x40, 0x00, 0x07, 0x08,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

// ICMP6 router solicitation with a source link-layer address option.
var icmp6RSSLL = []byte{
	0x60, 0x00, 0x00, 0x00, 0x00, 0x10, 0x3a, 0xff, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x02, 0x85, 0x00, 0x63, 0xb9, 0x00, 0x00, 0x00, 0x00,
	0x01, 0x01, 0xc2, 0x00, 0x54, 0xf5, 0x00, 0x00}

// ICMP6 router solicitation with no options.
var icmp6RS = []byte{
	0x60, 0x00, 0x00, 0x00, 0x00, 0x08, 0x3a, 0xff, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x02, 0x85, 0x00, 0x7b, 0xb8, 0x00, 0x00, 0x00, 0x00}

func TestICMPv4(t *testing.T) {
	queries := []packettestutil.FieldQuery{
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
			Result: []byte{0x01},
		},
		{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
			Result: []byte{0x02},
		},
	}
	updates := []packettestutil.FieldUpdate{
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
			Arg: []byte{0x10},
			Op:  fwdpacket.OpSet,
		},
		{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
			Arg: []byte{0x11},
			Op:  fwdpacket.OpSet,
		},
	}
	tests := []packettestutil.PacketFieldTest{
		// ICMPv4 echo request message. The checksum is initially set to
		// zero. The ICMP message changes with updates.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP4,
				ip4icmp,
				icmp4,
			},
			Queries: queries,
			Updates: updates,
			Final: [][]byte{
				ethernetIP4,
				ip4icmp,
				[]byte{0x10, 0x11, 0xe6, 0xe2, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
			},
		},
	}

	packettestutil.TestPacketFields("icmp", t, tests)
}

func TestICMPv6(t *testing.T) {
	tests := []packettestutil.PacketFieldTest{
		// ICMPv6 echo request message. The checksum is initially set to
		// zero. The ICMP message changes with updates.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				ip6icmp,
				icmp6,
			},
			Queries: []packettestutil.FieldQuery{
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
					Result: []byte{0x01},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
					Result: []byte{0x02},
				},
			},
			Updates: []packettestutil.FieldUpdate{
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
					Arg: []byte{0x10},
					Op:  fwdpacket.OpSet,
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
					Arg: []byte{0x11},
					Op:  fwdpacket.OpSet,
				},
			},
			Final: [][]byte{
				ethernetIP6,
				ip6icmp,
				[]byte{0x10, 0x11, 0xe6, 0xe2, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
			},
		},
		// ICMP6 neighbour advertisement with a target link-layer address option.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				icmp6NATLL,
			},
			Queries: []packettestutil.FieldQuery{
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
					Result: []byte{136},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
					Result: []byte{0},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TARGET, 0),
					Result: []byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00},
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_SLL, 0),
					Err: "does not exist",
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TLL, 0),
					Result: []byte{0xc2, 0x00, 0x54, 0xf5, 0x00, 0x00},
				},
			},
			Final: [][]byte{
				ethernetIP6,
				icmp6NATLL,
			},
		},
		// ICMP6 neighbour advertisement with no options.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				icmp6NA,
			},
			Queries: []packettestutil.FieldQuery{
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
					Result: []byte{136},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
					Result: []byte{0},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TARGET, 0),
					Result: []byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00},
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_SLL, 0),
					Err: "does not exist",
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TLL, 0),
					Err: "does not exist",
				},
			},
			Final: [][]byte{
				ethernetIP6,
				icmp6NA,
			},
		},
		// ICMP6 neighbour solicitation with a source link-layer address option.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				icmp6NSSLL,
			},
			Queries: []packettestutil.FieldQuery{
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
					Result: []byte{135},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
					Result: []byte{0},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TARGET, 0),
					Result: []byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_SLL, 0),
					Result: []byte{0xc2, 0x00, 0x54, 0xf5, 0x00, 0x00},
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TLL, 0),
					Err: "does not exist",
				},
			},
			Final: [][]byte{
				ethernetIP6,
				icmp6NSSLL,
			},
		},
		// ICMP6 neighbour solicitation with no options.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				icmp6NS,
			},
			Queries: []packettestutil.FieldQuery{
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
					Result: []byte{135},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
					Result: []byte{0},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TARGET, 0),
					Result: []byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x00},
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_SLL, 0),
					Err: "does not exist",
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TLL, 0),
					Err: "does not exist",
				},
			},
			Final: [][]byte{
				ethernetIP6,
				icmp6NS,
			},
		},
		// ICMP6 router advertisement with a source link-layer address option.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				icmp6RASLL,
			},
			Queries: []packettestutil.FieldQuery{
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
					Result: []byte{134},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
					Result: []byte{0},
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TARGET, 0),
					Err: "does not exist",
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_SLL, 0),
					Result: []byte{0xc2, 0x00, 0x54, 0xf5, 0x00, 0x00},
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TLL, 0),
					Err: "does not exist",
				},
			},
			Final: [][]byte{
				ethernetIP6,
				icmp6RASLL,
			},
		},
		// ICMP6 router advertisement with no options.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				icmp6RA,
			},
			Queries: []packettestutil.FieldQuery{
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
					Result: []byte{134},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
					Result: []byte{0},
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TARGET, 0),
					Err: "does not exist",
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_SLL, 0),
					Err: "does not exist",
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TLL, 0),
					Err: "does not exist",
				},
			},
			Final: [][]byte{
				ethernetIP6,
				icmp6RA,
			},
		},
		// ICMP6 router solicitation with no options.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				icmp6RS,
			},
			Queries: []packettestutil.FieldQuery{
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
					Result: []byte{133},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
					Result: []byte{0},
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TARGET, 0),
					Err: "does not exist",
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_SLL, 0),
					Err: "does not exist",
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TLL, 0),
					Err: "does not exist",
				},
			},
			Final: [][]byte{
				ethernetIP6,
				icmp6RS,
			},
		},
		// ICMP6 router solicitation with a source link-layer address option.
		{
			StartHeader: fwdpb.PacketHeaderId_ETHERNET,
			Orig: [][]byte{
				ethernetIP6,
				icmp6RSSLL,
			},
			Queries: []packettestutil.FieldQuery{
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_TYPE, 0),
					Result: []byte{133},
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP_CODE, 0),
					Result: []byte{0},
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TARGET, 0),
					Err: "does not exist",
				},
				{
					ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_SLL, 0),
					Result: []byte{0xc2, 0x00, 0x54, 0xf5, 0x00, 0x00},
				},
				{
					ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ICMP6_ND_TLL, 0),
					Err: "does not exist",
				},
			},
			Final: [][]byte{
				ethernetIP6,
				icmp6RSSLL,
			},
		},
	}

	packettestutil.TestPacketFields("icmp", t, tests)
}
