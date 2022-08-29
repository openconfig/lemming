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

// Package ethernet implements the Ethernet header.
package ethernet

import (
	"encoding/binary"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Size of various fields in an ethernet header.
const (
	macBytes      = 6                     // Number of bytes in a mac address.
	nextBytes     = 2                     // Number of bytes in the ethertype.
	vlanBytes     = 2                     // Number of bytes used to encode a vlan tag and qos bits.
	tagBitSize    = 12                    // Number of bits used to encode a vlan tag.
	tagBitPos     = 0                     // Offset in bits of the vlan tag.
	qosBitSize    = 4                     // Number of bits used to encode the vlan qos.
	qosBitPos     = 12                    // Offset in bits of the vlan qos.
	ethernetBytes = 14                    // Number of bytes in an ethernet header.
	headerBytes   = vlanBytes + nextBytes // Number of bytes used to encode a vlan header.
)

// Offset in bytes of various fields from the start of the ethernet header.
const (
	dstMACOffset = 0 // Destination mac address.
	srcMACOffset = 6 // Source mac address.
)

// nextOffsets is an array of offsets in bytes where the "next" field can occur.
// The array is indexed by the instance of the "next" field starting with 0.
// The maximum number of next fields for an ethernet frame with 2 vlan tags
// is 3.
var nextOffsets = []int{12, 16, 20}

// tagOffsets is an array of offset in bytes where the vlan tags can occur.
// The array is indexed by the instance of the vlan tag starting with 0.
// The maximum number of vlan tags for an ethernet frame with 2 vlan tags
// is 2.
var tagOffsets = []int{14, 18}

// Constants defined for various ether types in network byte order.
const (
	nextSingle   = 0x8100       // Single tagged frame.
	nextDouble   = 0x9100       // Double tagged frame.
	nextARP      = 0x0806       // ARP payload.
	nextIP4      = 0x0800       // IP4 payload.
	nextIP6      = 0x86DD       // IP6 payload.
	nextReserved = 0xFFFF       // Opaque data.
	Reserved     = nextReserved // Opaque data.
)

// NextHeader maps ether-types to packet header.
var NextHeader = map[uint16]fwdpb.PacketHeaderId{}

// HeaderNext maps packet headers to ether-types.
var HeaderNext = map[fwdpb.PacketHeaderId]uint16{
	fwdpb.PacketHeaderId_IP4: nextIP4,
	fwdpb.PacketHeaderId_IP6: nextIP6,
	fwdpb.PacketHeaderId_ARP: nextARP,
}

// An Ethernet presents a ethernet header. It can parse, query, add, remove
// and operate on Ethernet headers in the packet with zero, one and two
// vlan tags.
type Ethernet struct {
	next   int                  // Offset of the ethertype.
	header frame.Header         // Ethernet header.
	desc   *protocol.Desc       // Protocol descriptor.
	id     fwdpb.PacketHeaderId // Header ID.
	vlans  []int                // Offset of each vlan tag.
}

// Header returns the ethernet header as a slice of bytes.
func (eth *Ethernet) Header() []byte {
	return eth.header
}

// Trailer returns the no trailing bytes.
func (Ethernet) Trailer() []byte {
	return nil
}

// ID returns the ethernet header's ID.
func (eth *Ethernet) ID(int) fwdpb.PacketHeaderId {
	return eth.id
}

// field returns the bytes as specified by id.
func (eth *Ethernet) field(id fwdpacket.FieldID) frame.Field {
	if id.IsUDF {
		return protocol.UDF(eth.header, id)
	}

	if id.Instance == fwdpacket.LastField {
		id.Instance = uint8(len(eth.vlans) - 1)
	}

	switch id.Num {
	case fwdpb.PacketFieldNum_ETHER_MAC_SRC:
		return eth.header.Field(srcMACOffset, macBytes)

	case fwdpb.PacketFieldNum_ETHER_MAC_DST:
		return eth.header.Field(dstMACOffset, macBytes)

	case fwdpb.PacketFieldNum_ETHER_TYPE:
		return eth.header.Field(eth.next, nextBytes)

	case fwdpb.PacketFieldNum_VLAN_TAG, fwdpb.PacketFieldNum_VLAN_PRIORITY:
		if int(id.Instance) >= len(eth.vlans) {
			// TOOD(neeleshb): We currently ignore the vlan priority due to b/31199367.
			if id.Num == fwdpb.PacketFieldNum_VLAN_PRIORITY {
				return make([]byte, vlanBytes)
			}
			return nil
		}
		return eth.header.Field(eth.vlans[id.Instance], vlanBytes)

	default:
		return nil
	}
}

// Field copies the specified bytes from the ethernet header.
func (eth *Ethernet) Field(id fwdpacket.FieldID) ([]byte, error) {
	field := eth.field(id)
	if field == nil {
		return nil, fmt.Errorf("ethernet: Field failed, field %v does not exist", id)
	}

	if id.IsUDF {
		return field.Copy(), nil
	}

	switch id.Num {
	case fwdpb.PacketFieldNum_VLAN_TAG:
		return field.BitField(tagBitPos, tagBitSize), nil

	case fwdpb.PacketFieldNum_VLAN_PRIORITY:
		return field.BitField(qosBitPos, qosBitSize), nil

	default:
		return field.Copy(), nil
	}
}

// UpdateField sets bytes within the ethernet header.
func (eth *Ethernet) UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	field := eth.field(id)
	if field == nil {
		return false, fmt.Errorf("ethernet: UpdateField failed, field %v does not exist", id)
	}

	if op != fwdpacket.OpSet {
		return false, fmt.Errorf("ethernet: UpdateField failed, unsupported op %v for field %v", op, id)
	}

	if id.IsUDF {
		return true, field.Set(arg)
	}
	switch id.Num {

	case fwdpb.PacketFieldNum_VLAN_TAG:
		field.SetBits(tagBitPos, tagBitSize, uint64(binary.BigEndian.Uint16(arg)))
		return true, nil

	case fwdpb.PacketFieldNum_VLAN_PRIORITY:
		field.SetBits(qosBitPos, qosBitSize, uint64(binary.BigEndian.Uint16(arg)))
		return true, nil

	default:
		return true, field.Set(arg)
	}
}

// Rebuild updates the ethernet header's next field.
func (eth *Ethernet) Rebuild() error {
	if next, ok := HeaderNext[eth.desc.PayloadID()]; ok {
		eth.header.Field(eth.next, nextBytes).SetValue(uint(next))
	}
	return nil
}

// Remove removes the ethernet header. Depending on the header id, it can remove
// just the vlan tags or it can remove the entire header.
func (eth *Ethernet) Remove(id fwdpb.PacketHeaderId) error {
	switch id {
	case fwdpb.PacketHeaderId_ETHERNET_VLAN, fwdpb.PacketHeaderId_ETHERNET_1Q:
		if eth.id != id {
			return fmt.Errorf("ethernet: Remove failed, cannot remove header %v from header %v", id, eth.id)
		}
		next := eth.header.Field(eth.next, nextBytes).Value()
		t, length := newEthernet(uint16(next), eth.desc)
		t.header = eth.header[:length]
		t.header.Field(t.next, nextBytes).SetValue(next)
		*eth = *t
		return nil

	case fwdpb.PacketHeaderId_ETHERNET:
		if eth.header == nil {
			return fmt.Errorf("ethernet: Remove failed, cannot remove %v from non-ethernet frame", id)
		}
		eth.header = nil
		eth.vlans = nil
		eth.id = fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE
		return nil

	default:
		return fmt.Errorf("ethernet: Remove failed, unsupported header %v", id)
	}
}

// Modify modifies an existing ethernet header by adding vlan tags.
func (eth *Ethernet) Modify(id fwdpb.PacketHeaderId) error {
	if eth.id != fwdpb.PacketHeaderId_ETHERNET {
		return fmt.Errorf("ethernet: Modify failed, cannot add header %v to header %v", id, eth.id)
	}

	next := eth.header.Field(eth.next, nextBytes).Value()
	src := eth.header.Field(srcMACOffset, macBytes)
	dst := eth.header.Field(dstMACOffset, macBytes)

	var length int
	switch id {
	case fwdpb.PacketHeaderId_ETHERNET_VLAN:
		var t *Ethernet
		t, length = newEthernet(nextSingle, eth.desc)
		*eth = *t
		eth.header = make(frame.Header, length)
		eth.header.Field(nextOffsets[0], nextBytes).SetValue(nextSingle)

	case fwdpb.PacketHeaderId_ETHERNET_1Q:
		var t *Ethernet
		t, length = newEthernet(nextDouble, eth.desc)
		*eth = *t
		eth.header = make(frame.Header, length)
		eth.header.Field(nextOffsets[0], nextBytes).SetValue(nextDouble)
		eth.header.Field(nextOffsets[1], nextBytes).SetValue(nextSingle)

	default:
		return fmt.Errorf("ethernet: Modify failed, cannot add header %v to header %v", id, eth.id)
	}
	eth.header.Field(eth.next, nextBytes).SetValue(next)
	eth.header.Field(srcMACOffset, macBytes).Set(src)
	eth.header.Field(dstMACOffset, macBytes).Set(dst)
	return nil
}

// newEthernet create a new ethernet handler based on the ethertype.
func newEthernet(next uint16, desc *protocol.Desc) (*Ethernet, int) {
	switch next {
	case nextSingle:
		// An ethernet frame with a single vlan has 1 vlan tag.
		return &Ethernet{
			next:  nextOffsets[1],
			vlans: []int{tagOffsets[0]},
			desc:  desc,
			id:    fwdpb.PacketHeaderId_ETHERNET_VLAN,
		}, ethernetBytes + 1*headerBytes

	case nextDouble:
		// An ethernet frame with a two vlans has 2 vlan tags.
		return &Ethernet{
			next:  nextOffsets[2],
			vlans: []int{tagOffsets[0], tagOffsets[1]},
			desc:  desc,
			id:    fwdpb.PacketHeaderId_ETHERNET_1Q,
		}, ethernetBytes + 2*headerBytes

	default:
		// An ethernet frame by default has no vlan tags.
		return &Ethernet{
			next: nextOffsets[0],
			desc: desc,
			id:   fwdpb.PacketHeaderId_ETHERNET,
		}, ethernetBytes
	}
}

// add adds an empty ethernet header.
func add(id fwdpb.PacketHeaderId, desc *protocol.Desc) (protocol.Handler, error) {
	switch id {
	case fwdpb.PacketHeaderId_ETHERNET:
		eth, length := newEthernet(nextReserved, desc)
		eth.header = make(frame.Header, length)
		return eth, nil

	case fwdpb.PacketHeaderId_ETHERNET_VLAN:
		// An ethernet frame with one vlan tag, sets the first "next" field
		// to indicate the presence of another tag.
		eth, length := newEthernet(nextSingle, desc)
		eth.header = make(frame.Header, length)
		eth.header.Field(nextOffsets[0], nextBytes).SetValue(nextSingle)
		return eth, nil

	case fwdpb.PacketHeaderId_ETHERNET_1Q:
		// An ethernet frame with two vlan tags, sets the first "next" field
		// to indicate the presence of a pair of tags, and sets the second
		// "next" field to indicate the presence of another tag.
		eth, length := newEthernet(nextDouble, desc)
		eth.header = make(frame.Header, length)
		eth.header.Field(nextOffsets[0], nextBytes).SetValue(nextDouble)
		eth.header.Field(nextOffsets[1], nextBytes).SetValue(nextSingle)
		return eth, nil

	default:
		return nil, fmt.Errorf("ethernet: add failed, unknown header %v", id)
	}
}

// parse parses an ethernet header.
func parse(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	next, err := frame.Peek(nextOffsets[0], nextBytes)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("ethernet: parse failed, err %v", err)
	}

	eth, length := newEthernet(uint16(next.Value()), desc)
	if eth.header, err = frame.ReadHeader(length); err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("ethernet: parse failed, err %v", err)
	}
	if id, ok := NextHeader[uint16(eth.header.Field(eth.next, nextBytes).Value())]; ok {
		return eth, id, nil
	}
	return eth, fwdpb.PacketHeaderId_OPAQUE, nil
}

func init() {
	// Compute reverse mapping tables.
	for id, next := range HeaderNext {
		NextHeader[next] = id
	}

	// Register the parse and add functions for the ETHERNET header and its variants.
	protocol.Register(fwdpb.PacketHeaderId_ETHERNET, parse, add)
	protocol.Register(fwdpb.PacketHeaderId_ETHERNET_VLAN, parse, add)
	protocol.Register(fwdpb.PacketHeaderId_ETHERNET_1Q, parse, add)
}
