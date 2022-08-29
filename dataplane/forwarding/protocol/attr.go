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

package protocol

import (
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Sizes of various attribute types in bytes.
const (
	SizeUint64 = 8
	SizeUint32 = 4
	SizeUint24 = 3
	SizeUint16 = 2
	SizeUint8  = 1
	SizeMAC    = 6
	SizeIP4    = SizeUint32
	SizeIP6    = 16
)

// addFn adds a protocol header to the packet described by Desc and returns
// the corresponding handler.
type addFn func(fwdpb.PacketHeaderId, *Desc) (Handler, error)

// parseFn parses a protocol from a frame and returns the id of the next
// packet header.
type parseFn func(*frame.Frame, *Desc) (Handler, fwdpb.PacketHeaderId, error)

// HeaderAttr contains attributes for packet headers indexed by the header id.
//
// HeaderAttr is fully initialized by init routines and is logically constant
// after init.
var HeaderAttr = map[fwdpb.PacketHeaderId]struct {
	Parse parseFn                 // registered function to parse the header
	Add   addFn                   // registered function to add the header
	Group fwdpb.PacketHeaderGroup // computed group of the header
}{}

// FieldAttr contains attributes for each packet field indexed by the field
// number. Each field has a valid size described as a discrete set of sizes.
//
// FieldAttr is fully initialized by init routines and is logically constant
// after init.
var FieldAttr = map[fwdpb.PacketFieldNum]struct {
	Sizes       []int                   // discrete valid sizes
	DefaultSize int                     // default size computed from Sizes in init
	Group       fwdpb.PacketHeaderGroup // header group containing this field computed in init
}{
	fwdpb.PacketFieldNum_PACKET_LENGTH: {
		Sizes: []int{SizeUint64},
	},
	fwdpb.PacketFieldNum_PACKET_PORT_INPUT: {
		Sizes: []int{SizeUint64},
	},
	fwdpb.PacketFieldNum_PACKET_PORT_OUTPUT: {
		Sizes: []int{SizeUint64},
	},
	fwdpb.PacketFieldNum_PACKET_VRF: {
		Sizes: []int{SizeUint16},
	},
	fwdpb.PacketFieldNum_PACKET_ATTRIBUTE_32: {
		Sizes: []int{SizeUint32},
	},
	fwdpb.PacketFieldNum_PACKET_ATTRIBUTE_24: {
		Sizes: []int{SizeUint24},
	},
	fwdpb.PacketFieldNum_PACKET_ATTRIBUTE_16: {
		Sizes: []int{SizeUint16},
	},
	fwdpb.PacketFieldNum_PACKET_ATTRIBUTE_8: {
		Sizes: []int{SizeUint8},
	},
	fwdpb.PacketFieldNum_ETHER_MAC_SRC: {
		Sizes: []int{SizeMAC},
	},
	fwdpb.PacketFieldNum_ETHER_MAC_DST: {
		Sizes: []int{SizeMAC},
	},
	fwdpb.PacketFieldNum_VLAN_TAG: {
		Sizes: []int{SizeUint16},
	},
	fwdpb.PacketFieldNum_VLAN_PRIORITY: {
		Sizes: []int{SizeUint16},
	},
	fwdpb.PacketFieldNum_ETHER_TYPE: {
		Sizes: []int{SizeUint16},
	},
	fwdpb.PacketFieldNum_IP_VERSION: {
		Sizes: []int{SizeUint8},
	},
	fwdpb.PacketFieldNum_IP_ADDR_SRC: {
		Sizes: []int{SizeIP4, SizeIP6},
	},
	fwdpb.PacketFieldNum_IP_ADDR_DST: {
		Sizes: []int{SizeIP4, SizeIP6},
	},
	fwdpb.PacketFieldNum_IP_PROTO: {
		Sizes: []int{SizeUint8},
	},
	fwdpb.PacketFieldNum_IP_HOP: {
		Sizes: []int{SizeUint8},
	},
	fwdpb.PacketFieldNum_IP_QOS: {
		// Trigger all IP QOS fields to be padded to 4B.
		// IPv6 traffic class is 1B in a 4B field.
		Sizes: []int{SizeUint32},
	},
	fwdpb.PacketFieldNum_IP6_FLOW: {
		// Trigger all IP FLOW fields to be padded to 4B.
		// IPv6 flow label is 20 bits in a 4B field.
		Sizes: []int{SizeUint32},
	},
	fwdpb.PacketFieldNum_ICMP_TYPE: {
		Sizes: []int{SizeUint8},
	},
	fwdpb.PacketFieldNum_ICMP_CODE: {
		Sizes: []int{SizeUint8},
	},
	fwdpb.PacketFieldNum_ICMP6_ND_SLL: {
		Sizes: []int{SizeMAC},
	},
	fwdpb.PacketFieldNum_ICMP6_ND_TLL: {
		Sizes: []int{SizeMAC},
	},
	fwdpb.PacketFieldNum_ICMP6_ND_TARGET: {
		Sizes: []int{SizeIP6},
	},
	fwdpb.PacketFieldNum_L4_PORT_SRC: {
		Sizes: []int{SizeUint16},
	},
	fwdpb.PacketFieldNum_L4_PORT_DST: {
		Sizes: []int{SizeUint16},
	},
	fwdpb.PacketFieldNum_TCP_FLAGS: {
		Sizes: []int{SizeUint16},
	},
	fwdpb.PacketFieldNum_ARP_TPA: {
		Sizes: []int{SizeUint32},
	},
	fwdpb.PacketFieldNum_ARP_SPA: {
		Sizes: []int{SizeUint32},
	},
	fwdpb.PacketFieldNum_GRE_KEY: {
		Sizes: []int{SizeUint32},
	},
	fwdpb.PacketFieldNum_ARP_TMAC: {
		Sizes: []int{SizeMAC},
	},
	fwdpb.PacketFieldNum_ARP_SMAC: {
		Sizes: []int{SizeMAC},
	},
	fwdpb.PacketFieldNum_GRE_SEQUENCE: {
		Sizes: []int{SizeUint32},
	},
}

// GroupAttr contains attributes for each packet header group.
//
// Each packet header group is associated with a set of packet headers and
// a set of packet fields that may be present in the group.
//
// Each packet header group is also given a position relative to other
// packet header groups within the packet.
//
// GroupAttr is logically constant and should not be changed.
var GroupAttr = map[fwdpb.PacketHeaderGroup]struct {
	Position int                    // position of the group within a frame
	headers  []fwdpb.PacketHeaderId // headers in the group
	fields   []fwdpb.PacketFieldNum // packet fields that can be queried in this group
}{
	fwdpb.PacketHeaderGroup_PACKET: {
		Position: 0,
		headers: []fwdpb.PacketHeaderId{
			fwdpb.PacketHeaderId_METADATA,
		},
		fields: []fwdpb.PacketFieldNum{
			fwdpb.PacketFieldNum_PACKET_PORT_INPUT,
			fwdpb.PacketFieldNum_PACKET_PORT_OUTPUT,
			fwdpb.PacketFieldNum_PACKET_LENGTH,
			fwdpb.PacketFieldNum_PACKET_VRF,
			fwdpb.PacketFieldNum_PACKET_ATTRIBUTE_32,
			fwdpb.PacketFieldNum_PACKET_ATTRIBUTE_24,
			fwdpb.PacketFieldNum_PACKET_ATTRIBUTE_16,
			fwdpb.PacketFieldNum_PACKET_ATTRIBUTE_8,
		},
	},
	fwdpb.PacketHeaderGroup_L2: {
		Position: 1,
		headers: []fwdpb.PacketHeaderId{
			fwdpb.PacketHeaderId_ETHERNET,
			fwdpb.PacketHeaderId_ETHERNET_VLAN,
			fwdpb.PacketHeaderId_ETHERNET_1Q,
		},
		fields: []fwdpb.PacketFieldNum{
			fwdpb.PacketFieldNum_ETHER_MAC_SRC,
			fwdpb.PacketFieldNum_ETHER_MAC_DST,
			fwdpb.PacketFieldNum_ETHER_TYPE,
			fwdpb.PacketFieldNum_VLAN_TAG,
			fwdpb.PacketFieldNum_VLAN_PRIORITY,
		},
	},
	fwdpb.PacketHeaderGroup_L3: {
		Position: 2,
		headers: []fwdpb.PacketHeaderId{
			fwdpb.PacketHeaderId_IP4,
			fwdpb.PacketHeaderId_IP,
			fwdpb.PacketHeaderId_TUNNEL_6TO4_AUTO,
			fwdpb.PacketHeaderId_TUNNEL_6TO4_SECURE,
			fwdpb.PacketHeaderId_IP6,
			fwdpb.PacketHeaderId_GRE,
		},
		fields: []fwdpb.PacketFieldNum{
			fwdpb.PacketFieldNum_IP_VERSION,
			fwdpb.PacketFieldNum_IP_ADDR_SRC,
			fwdpb.PacketFieldNum_IP_ADDR_DST,
			fwdpb.PacketFieldNum_IP_HOP,
			fwdpb.PacketFieldNum_IP_PROTO,
			fwdpb.PacketFieldNum_IP_QOS,
			fwdpb.PacketFieldNum_IP6_FLOW,
			fwdpb.PacketFieldNum_GRE_KEY,
			fwdpb.PacketFieldNum_GRE_SEQUENCE,
		},
	},
	fwdpb.PacketHeaderGroup_L4: {
		Position: 3,
		headers: []fwdpb.PacketHeaderId{
			fwdpb.PacketHeaderId_TCP,
			fwdpb.PacketHeaderId_UDP,
		},
		fields: []fwdpb.PacketFieldNum{
			fwdpb.PacketFieldNum_L4_PORT_SRC,
			fwdpb.PacketFieldNum_L4_PORT_DST,
			fwdpb.PacketFieldNum_TCP_FLAGS,
		},
	},
	fwdpb.PacketHeaderGroup_PAYLOAD: {
		Position: 4,
		headers: []fwdpb.PacketHeaderId{
			fwdpb.PacketHeaderId_ICMP4,
			fwdpb.PacketHeaderId_ICMP6,
			fwdpb.PacketHeaderId_ARP,
			fwdpb.PacketHeaderId_OPAQUE,
		},
		fields: []fwdpb.PacketFieldNum{
			fwdpb.PacketFieldNum_ICMP_TYPE,
			fwdpb.PacketFieldNum_ICMP_CODE,
			fwdpb.PacketFieldNum_ICMP6_ND_TARGET,
			fwdpb.PacketFieldNum_ICMP6_ND_TLL,
			fwdpb.PacketFieldNum_ICMP6_ND_SLL,
			fwdpb.PacketFieldNum_ARP_TPA,
			fwdpb.PacketFieldNum_ARP_SPA,
			fwdpb.PacketFieldNum_ARP_TMAC,
			fwdpb.PacketFieldNum_ARP_SMAC,
		},
	},
}

// Sequence is the sequence of packet header groups used to reconstruct a frame.
// Each group can occur exactly once in the sequence.
//
// Sequence is computed during initialization and is logically constnant and
// should not be changed after init.
var Sequence []fwdpb.PacketHeaderGroup

func init() {
	Sequence = make([]fwdpb.PacketHeaderGroup, len(fwdpb.PacketHeaderGroup_name))

	// Compute the default size for each field by using the maximum size.
	for pos, attr := range FieldAttr {
		for _, size := range attr.Sizes {
			if size > attr.DefaultSize {
				attr.DefaultSize = size
			}
		}
		FieldAttr[pos] = attr
	}

	// Process group attributes to update header and field attributes
	// and determine the sequence.
	for group, attr := range GroupAttr {
		for _, id := range attr.headers {
			attr := HeaderAttr[id]
			attr.Group = group
			HeaderAttr[id] = attr
		}

		for _, field := range attr.fields {
			attr := FieldAttr[field]
			attr.Group = group
			FieldAttr[field] = attr
		}

		if attr.Position > len(Sequence) {
			panic(fmt.Sprintf("protocol: incorrect position %v specified for group %v", attr.Position, group))
		}
		Sequence[attr.Position] = group
	}
}

// Register registers handlers to add and parse a packet header.
func Register(id fwdpb.PacketHeaderId, parse parseFn, add addFn) {
	attr := HeaderAttr[id]
	attr.Parse = parse
	attr.Add = add
	HeaderAttr[id] = attr
}
