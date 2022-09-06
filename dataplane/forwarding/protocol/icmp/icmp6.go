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

// Package icmp implements the ICMP header support in Lucius.
package icmp

import (
	"encoding/binary"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	"github.com/openconfig/lemming/dataplane/forwarding/util/hash/csum16"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	lengthBytes    = 1                // Number of bytes in the ICMP Option length
	lengthOffset   = 1                // Offset in bytes of the ICMP Option length
	ndTargetBytes  = protocol.SizeIP6 // Number of bytes in the target address
	ndTargetOffset = 8                // Offset in bytes of the target address
	ndSLLBytes     = protocol.SizeMAC // Number of bytes in the source link address
	ndTLLBytes     = protocol.SizeMAC // Number of bytes in the target link address
)

// Set of supported ICMP6 message types.
const (
	ICMP6RouterSolicitation    = uint8(133)
	ICMP6RouterAdvertisement   = uint8(134)
	ICMP6NeighborSolicitation  = uint8(135)
	ICMP6NeighborAdvertisement = uint8(136)
	ICMP6NeighborRedirect      = uint8(137)

	protoICMP6 = 58
)

// Set of supported ICMPv6 options
const (
	ICMP6SLL = uint8(1)
	ICMP6TLL = uint8(2)
)

// Offset of the option indexed by the ICMPv6 message type.
var optionOffset = map[uint8]int{
	ICMP6RouterSolicitation:    8,
	ICMP6RouterAdvertisement:   16,
	ICMP6NeighborSolicitation:  24,
	ICMP6NeighborAdvertisement: 24,
	ICMP6NeighborRedirect:      40,
}

// An ICMP6 represents an ICMP header in the packet. It can update the ICMP
// header but it cannot add or remove an ICMP header. The ICMP message is
// assumed to contain the ICMP header and all bytes following it.
//
// Note that the ICMP tracks only one instance of each optional field.
type ICMP6 struct {
	ICMP
	id      uint8                                // message type of the ICMPv6
	options map[fwdpb.PacketFieldNum]frame.Field // map of optional values indexed by their well known field numbers
	desc    *protocol.Desc
}

// Parse the header to create the map of options. We read the options until
// the header is exhausted. Note that the value of an option can be nil. If
// the header has no options, an empty map is returned.
func parseOptions(id uint8, header frame.Header) (map[fwdpb.PacketFieldNum]frame.Field, error) {
	options := make(map[fwdpb.PacketFieldNum]frame.Field)

	// find the offset for options based on the header type.
	offset, ok := optionOffset[id]
	if !ok {
		return options, nil
	}
	for {
		var f frame.Field

		// read the type of the option. The read is expected to fail if the buffer is exhausted.
		if f = header.Field(typeOffset+offset, typeBytes); f == nil {
			break
		}
		id := f.Value()

		// read the type of the option. The read is not expected.
		if f = header.Field(lengthOffset+offset, lengthBytes); f == nil {
			return nil, fmt.Errorf("Unable to read length for option %v", id)
		}
		length := int(8 * f.Value())

		// length in the TLV includes the bytes containing the type and length.
		// Use it to compute the offet to the next option. If the option is
		// interesting, determine the length of the value.
		valueLen := length - (typeBytes + lengthBytes)
		if valueLen < 0 {
			return nil, fmt.Errorf("Incorrect length %v for option %v", length, id)
		}

		switch uint8(id) {
		case ICMP6SLL:
			options[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ICMP6_ND_SLL] = header.Field(lengthOffset+offset+lengthBytes, valueLen)

		case ICMP6TLL:
			options[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ICMP6_ND_TLL] = header.Field(lengthOffset+offset+lengthBytes, valueLen)
		}
		offset += length
	}
	return options, nil
}

// ID returns the ICMP protocol header ID.
func (ICMP6) ID(int) fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_PACKET_HEADER_ID_ICMP6
}

// field returns bytes within the ICMP header as identified by id.
func (i *ICMP6) field(id fwdpacket.FieldID) frame.Field {
	if f := i.commonField(id); f != nil {
		return f
	}

	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ICMP6_ND_TARGET:
		switch i.id {
		case ICMP6NeighborSolicitation, ICMP6NeighborAdvertisement, ICMP6NeighborRedirect:
			return i.header.Field(ndTargetOffset, ndTargetBytes)
		default:
			return nil
		}

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ICMP6_ND_SLL, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ICMP6_ND_TLL:
		if option, ok := i.options[id.Num]; ok {
			return option
		}
	}
	return nil
}

// Field returns a copy of the bytes for the specfied field in the ICMP header.
func (i *ICMP6) Field(id fwdpacket.FieldID) ([]byte, error) {
	if field := i.field(id); field != nil {
		return field.Copy(), nil
	}
	return nil, fmt.Errorf("icmp6: Field failed, field %v does not exist", id)
}

// UpdateField sets bytes within the ICMP header.
func (i *ICMP6) UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	if field := i.field(id); field != nil && op == fwdpacket.OpSet {
		return true, field.Set(arg)
	}
	return false, fmt.Errorf("icmp6: UpdateField failed, unsupported op %v for field %v", op, id)
}

// Rebuild updates the ICMP checksum. The checksum is computed over the ICMP
// message and the IPv6 pseudo-header. The checksum is recomputed only if the
// envelope and/or ICMP header are dirty.
func (i *ICMP6) Rebuild() error {
	e := i.desc.EnvelopeDesc()
	if e != nil && !e.Dirty() && !i.desc.Dirty() {
		return nil
	}

	i.header.Field(csumOffset, csumBytes).SetValue(0)
	var sum csum16.Sum

	var err error
	var f []byte
	if f, err = i.desc.Packet.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, fwdpacket.LastField)); err != nil {
		return fmt.Errorf("icmp6: Rebuild failed: %v", err)
	}
	sum.Write(f)
	if f, err = i.desc.Packet.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, fwdpacket.LastField)); err != nil {
		return fmt.Errorf("icmp6: Rebuild failed: %v", err)
	}
	sum.Write(f)

	length := len(i.header)
	f = make([]byte, protocol.SizeUint32)
	binary.BigEndian.PutUint32(f, uint32(length))
	sum.Write(f)

	sum.Write([]byte{0, 0, 0, protoICMP6})
	sum.Write(i.header)
	i.header.Field(csumOffset, csumBytes).SetValue(uint(sum))
	return nil
}

// parseICMP6 parses an ICMPv6 header in the packet. All data following the ICMP
// header is assumed to be a part of the ICMP message.
func parseICMP6(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	if frame.Len() < icmpBytes {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("icmp6: parse failed, invalid frame length %v, expected size is greater than %v bytes", frame.Len(), icmpBytes)
	}
	header, err := frame.ReadHeader(frame.Len())
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("icmp6: parse failed, err %v", err)
	}
	field := header.Field(typeOffset, typeBytes)
	id := uint8(field.Value())

	options, err := parseOptions(id, header)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("icmp6: parse failed, err %v", err)
	}
	i := &ICMP6{
		id:      id,
		options: options,
		desc:    desc,
	}
	i.ICMP.header = header
	return i, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, nil
}

func init() {
	// ICMP header cannot be added to a packet.
	protocol.Register(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ICMP6, parseICMP6, nil)
}
