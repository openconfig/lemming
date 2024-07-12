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
	"bytes"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	"github.com/openconfig/lemming/dataplane/forwarding/util/hash/csum16"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Size of various fields in an IP4 header.
const (
	ip4AddrBytes   = 4            // Number of bytes in the IP4 address
	ip4DescBytes   = 1            // Number of bytes in the header description
	ip4DescPos     = 0            // Offset in bytes of the header description
	ip4HeaderBits  = 4            // Number of bits in the IP header length
	ip4HeaderPos   = 0            // Offset in bits of the IP header length
	ip4TosBytes    = 1            // Number of bytes for the tos
	ip4TosPos      = 1            // Offset in bytes for the tos
	ip4LengthBytes = 2            // Number of bytes in the IP length
	ip4LengthPos   = 2            // Offset in bytes of the IP length
	ttlBytes       = 1            // Number of bytes in the TTL
	ttlPos         = 8            // Offset in bytes of the TTL
	ip4ProtoBytes  = 1            // Number of bytes in the IP protocol field
	ip4ProtoPos    = 9            // Offset in bytes of the IP protocol
	ip4CSumBytes   = 2            // Number of bytes in the checksum
	ip4CSumPos     = 10           // Offset in bytes of the checksum
	ip4SrcBytes    = ip4AddrBytes // Number of bytes in the source IP address
	ip4SrcPos      = 12           // Offset in bytes of the source IP address
	ip4DstBytes    = ip4AddrBytes // Number of bytes in the dest IP address
	ip4DstPos      = 16           // Offset in bytes of the dest IP address
	ip6to4Offset   = 2            // Offset in bytes of the encoded ip4 address
)

// Constants defined for various 6to4 tunnels
var ip6to4Prefix = []byte{0x20, 0x02} // Prefix of the computed 6to4 address

// An IP4 represents an IPv4 header including IP options in the packet.
// It can be query, update, add and remove the IP header.
type IP4 struct {
	header  frame.Header // IPv4 header.
	payload int64        // Length of the payload in bytes.
}

// Header returns the IPv4 header as a slice of bytes.
func (ip *IP4) Header() []byte {
	return ip.header
}

// ID returns the protocol header ID.
func (IP4) ID() fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4
}

// field returns the bytes as specified by id.
func (ip *IP4) field(id fwdpacket.FieldID) frame.Field {
	if id.IsUDF {
		return protocol.UDF(ip.header, id)
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC:
		return ip.header.Field(ip4SrcPos, ip4SrcBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST:
		return ip.header.Field(ip4DstPos, ip4DstBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP:
		return ip.header.Field(ttlPos, ttlBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO:
		return ip.header.Field(ip4ProtoPos, ip4ProtoBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS:
		return ip.header.Field(ip4TosPos, ip4TosBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION:
		return ipVersion(ip.header)

	default:
		return nil
	}
}

// Find returns a copy of the specified field.
func (ip *IP4) Find(id fwdpacket.FieldID) ([]byte, error) {
	if field := ip.field(id); field != nil {
		return field.Copy(), nil
	}
	return nil, fmt.Errorf("ip4: Find failed, field %v does not exist", id)
}

// Update updates the specified field.
func (ip *IP4) Update(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	field := ip.field(id)
	if field == nil {
		return false, fmt.Errorf("ip4: Update failed, field %v does not exist", id)
	}
	switch op {
	case fwdpacket.OpDec:
		// When decrementing the TTL, adjust the checksum and report
		// the header as clean.
		if !id.IsUDF && id.Num == fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP {
			pttl := uint16(field.Value())
			nttl := pttl - 1
			field.SetValue(uint(nttl))

			psum := uint16(ip.header.Field(ip4CSumPos, ip4CSumBytes).Value())
			nsum := ^(^psum + ^pttl + nttl)
			ip.header.Field(ip4CSumPos, ip4CSumBytes).SetValue(uint(nsum))
			return false, nil
		}

	case fwdpacket.OpSet:
		return true, field.Set(arg)
	}
	return false, fmt.Errorf("ip4: update failed, unsupported op %v for field %v", op, id)
}

// Payload gets the payload information.
func (ip *IP4) Payload() (fwdpb.PacketHeaderId, int64) {
	return ip.ID(), int64(len(ip.header)) + ip.payload
}

// SetPayload sets the payload.
func (ip *IP4) SetPayload(id fwdpb.PacketHeaderId, length int64) {
	proto, ok := headerProto[id]
	if !ok {
		proto = protoReserved
	}
	if id != fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE { // If the packet header is an unknown type, don't change it.
		ip.header.Field(ip4ProtoPos, ip4ProtoBytes).SetValue(uint(proto))
	}

	ip.header.Field(ip4LengthPos, ip4LengthBytes).SetValue(uint(length) + uint(len(ip.header)))
	ip.payload = length

	ip.header.Field(ip4CSumPos, ip4CSumBytes).SetValue(0)
	var sum csum16.Sum
	sum.Write(ip.header)
	ip.header.Field(ip4CSumPos, ip4CSumBytes).SetValue(uint(sum))
}

// newIP4 creates a new IPv4 header.
func newIP4() header {
	ip := &IP4{}
	ip.header = make(frame.Header, 20)
	ip.header.Field(ip4DescPos, ip4DescBytes).SetValue(0x45)
	return ip
}

// check6to4Tunnel checks the inner and outer IP headers to detect if it
// represents a valid 6to4 tunne.
func check6to4Tunnel(id fwdpb.PacketHeaderId, inner, outer header) error {
	if pid := inner.ID(); pid != fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6 {
		return fmt.Errorf("ip4: Incorrect 6to4 tunnel, inner header %v is not IP6", pid)
	}
	if pid := outer.ID(); pid != fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4 {
		return fmt.Errorf("ip4: Incorrect 6to4 tunnel, outer header %v is not IP4", pid)
	}

	check := func(inner, outer header, fid fwdpacket.FieldID) error {
		innerIP, e := inner.Find(fid)
		if e != nil {
			return e
		}
		if len(innerIP) != ip6AddrBytes {
			return fmt.Errorf("address %x has incorrect size %v", innerIP, len(innerIP))
		}
		outerIP, e := outer.Find(fid)
		if e != nil {
			return e
		}
		if len(outerIP) != ip4AddrBytes {
			return fmt.Errorf("address %x has incorrect size %v", outerIP, len(outerIP))
		}
		if !bytes.Equal(outerIP, innerIP[ip6to4Offset:ip6to4Offset+ip4AddrBytes]) {
			return fmt.Errorf("inner IP6 address %x does not contain the outer IP4 address %x", innerIP, outerIP)
		}
		if !bytes.Equal(ip6to4Prefix, innerIP[:len(ip6to4Prefix)]) {
			return fmt.Errorf("inner IP6 address %x does not contain the 6to4 prefix %x", innerIP, ip6to4Prefix)
		}
		return nil
	}

	switch id {
	case fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO:
		if e := check(inner, outer, fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, fwdpacket.FirstField)); e != nil {
			return fmt.Errorf("ip4: Incorrect 6to4 tunnel, bad dst address, err %v", e)
		}

	case fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE:
		if e := check(inner, outer, fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, fwdpacket.FirstField)); e != nil {
			return fmt.Errorf("ip4: Incorrect 6to4 tunnel, bad dst address, err %v", e)
		}
		if e := check(inner, outer, fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, fwdpacket.FirstField)); e != nil {
			return fmt.Errorf("ip4: Incorrect 6to4 tunnel, bad src address, err %v", e)
		}
	}
	return nil
}

// new6to4 creates a new IP4 header to carry IPv6 packets.
func new6to4Tunnel(id fwdpb.PacketHeaderId, inner header) (header, error) {
	if pid := inner.ID(); pid != fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6 {
		return nil, fmt.Errorf("ip4: Unable to add 6to4 tunnel for %v as payload", pid)
	}
	outer := newIP4()

	update := func(inner, outer header, fid fwdpacket.FieldID) error {
		a, e := inner.Find(fid)
		if e != nil {
			return e
		}
		if len(a) != ip6AddrBytes {
			return fmt.Errorf("ip4: Unable to add 6to4 tunnel, address %x has incorrect size %v", a, len(a))
		}

		b := a[ip6to4Offset : ip6to4Offset+ip4AddrBytes]
		outer.Update(fid, fwdpacket.OpSet, b)
		return nil
	}

	switch id {
	case fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_AUTO:
		if e := update(inner, outer, fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, fwdpacket.FirstField)); e != nil {
			return nil, fmt.Errorf("ip4: Unable to add 6to4 tunnel %v, %v", id, e)
		}

	case fwdpb.PacketHeaderId_PACKET_HEADER_ID_TUNNEL_6TO4_SECURE:
		if e := update(inner, outer, fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, fwdpacket.FirstField)); e != nil {
			return nil, fmt.Errorf("ip4: Unable to add 6to4 tunnel %v, %v", id, e)
		}
		if e := update(inner, outer, fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, fwdpacket.FirstField)); e != nil {
			return nil, fmt.Errorf("ip4: Unable to add 6to4 tunnel %v, %v", id, e)
		}

	default:
		return nil, fmt.Errorf("ip4: Unable to add 6to4 tunnel %v, unknown type", id)
	}
	return outer, nil
}

// makeIP4 creates a new IPv4 header from a frame.
func makeIP4(frame *frame.Frame) (header, fwdpb.PacketHeaderId, error) {
	desc, err := frame.Peek(ip4DescPos, ip4DescBytes)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("ip4: makeIP4 failed, err %v", err)
	}
	size := desc.BitField(ip4HeaderPos, ip4HeaderBits)
	header, err := frame.ReadHeader(int(size.Value() << 2))
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("ip4: makeIP4 failed, err %v", err)
	}
	ip := &IP4{
		header:  header,
		payload: int64(frame.Len()),
	}
	if next, ok := protoHeader[uint8(header.Field(ip4ProtoPos, ip4ProtoBytes).Value())]; ok {
		return ip, next, nil
	}
	return ip, fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE, nil
}
