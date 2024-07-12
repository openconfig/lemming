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

package ip

import (
	"encoding/binary"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Size of various fields in an ip6 header.
const (
	ip6AddrBytes   = 16           // Number of bytes in the IP address.
	ip6DescBytes   = 4            // Number of bytes in the header description.
	ip6DescPos     = 0            // Offset in bytes of the header description.
	ip6TosBits     = 8            // Number of bits for the tos.
	ip6TosPos      = 20           // Offset in bits for the tos.
	ip6FlowBits    = 20           // Number of bits in the IP flow.
	ip6FlowPos     = 0            // Offset in bits for the IP flow.
	ip6VersionBits = 4            // Number of bits in the IP version.
	ip6VersionPos  = 28           // Offset in bits for the IP version.
	ip6LengthBytes = 2            // Number of bytes in the IP length.
	ip6LengthPos   = 4            // Offset in bytes of the IP length.
	hopBytes       = 1            // Number of bytes in the Hop.
	hopPos         = 7            // Offset in bytes of the Hop.
	ip6ProtoBytes  = 1            // Number of bytes in the IP protocol field.
	ip6ProtoPos    = 6            // Offset in bytes of the IP protocol.
	ip6SrcBytes    = ip6AddrBytes // Number of bytes in the source IP address.
	ip6SrcPos      = 8            // Offset in bytes of the source IP address.
	ip6DstBytes    = ip6AddrBytes // Number of bytes in the dest IP address.
	ip6DstPos      = 24           // Offset in bytes of the dest IP address.
	ip6HeaderBytes = 40           // Number of bytes in the fixed sized IP6 header.
)

// An IP6 represents an IPv6 header in the packet. It can query,
// update, add and remove the IPv6 header. Note that it does not support
// IPv6 extensions.
type IP6 struct {
	header  frame.Header // IPv6 header.
	payload int64        // Length of the payload in bytes.
}

// Header returns the header as a slice of bytes.
func (ip *IP6) Header() []byte {
	return ip.header
}

// ID returns the protocol header ID.
func (IP6) ID() fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6
}

// field returns the bytes as specified by id.
func (ip *IP6) field(id fwdpacket.FieldID) frame.Field {
	if id.IsUDF {
		return protocol.UDF(ip.header, id)
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC:
		return ip.header.Field(ip6SrcPos, ip6SrcBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST:
		return ip.header.Field(ip6DstPos, ip6DstBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP:
		return ip.header.Field(hopPos, hopBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO:
		return ip.header.Field(ip6ProtoPos, ip6ProtoBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP6_FLOW:
		return ip.header.Field(ip6DescPos, ip6DescBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION:
		return ipVersion(ip.header)

	default:
		return nil
	}
}

// Find returns a copy of the specified field.
func (ip *IP6) Find(id fwdpacket.FieldID) ([]byte, error) {
	field := ip.field(id)
	if field == nil {
		return nil, fmt.Errorf("ip6: find failed, field %v does not exist", id)
	}

	if id.IsUDF {
		return field.Copy(), nil
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS:
		return field.BitField(ip6TosPos, ip6TosBits), nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP6_FLOW:
		return field.BitField(ip6FlowPos, ip6FlowBits), nil

	default:
		return field.Copy(), nil
	}
}

// Update updates the specified field.
func (ip *IP6) Update(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	field := ip.field(id)
	if field == nil {
		return false, fmt.Errorf("ip6: update failed, field %v does not exist", id)
	}

	switch op {
	case fwdpacket.OpDec:
		if id.IsUDF || id.Num != fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP {
			return false, fmt.Errorf("ip6: update failed, unsupported op %v for field %v", op, id)
		}
		ttl := field.Value()
		field.SetValue(ttl - 1)
		return false, nil

	case fwdpacket.OpSet:
		switch id.Num {
		case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS:
			field.SetBits(ip6TosPos, ip6TosBits, uint64(binary.BigEndian.Uint32(arg)))
			return true, nil

		case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP6_FLOW:
			field.SetBits(ip6FlowPos, ip6FlowBits, uint64(binary.BigEndian.Uint32(arg)))
			return true, nil

		default:
			return true, field.Set(arg)
		}

	default:
		return false, fmt.Errorf("ip6: update failed, unsupported op %v for field %v", op, id)
	}
}

// Payload gets the payload information.
func (ip *IP6) Payload() (fwdpb.PacketHeaderId, int64) {
	return ip.ID(), int64(len(ip.header)) + ip.payload
}

// SetPayload sets the payload.
func (ip *IP6) SetPayload(id fwdpb.PacketHeaderId, length int64) {
	proto, ok := headerProto[id]
	if !ok {
		proto = protoReserved
	}
	if id != fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE { // If the packet header is an unknown type, don't change it.
		ip.header.Field(ip6ProtoPos, ip6ProtoBytes).SetValue(uint(proto))
	}

	ip.header.Field(ip6LengthPos, ip6LengthBytes).SetValue(uint(length))
	ip.payload = length
}

func newIP6() header {
	ip := &IP6{}
	ip.header = make(frame.Header, ip6HeaderBytes)
	ip.header.Field(ip6DescPos, ip6DescBytes).SetBits(ip6VersionPos, ip6VersionBits, 0x6)
	return ip
}

// makeIP6 parses an IP6 header.
func makeIP6(frame *frame.Frame) (header, fwdpb.PacketHeaderId, error) {
	header, err := frame.ReadHeader(ip6HeaderBytes)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("ip6: makeIP6 failed, err %v", err)
	}

	ip := &IP6{
		header:  header,
		payload: int64(frame.Len()),
	}
	if next, ok := protoHeader[uint8(header.Field(ip6ProtoPos, ip6ProtoBytes).Value())]; ok {
		return ip, next, nil
	}
	return ip, fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE, nil
}
