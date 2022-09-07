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

// Package tcp implements the TCP header support in Lucius.
package tcp

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	"github.com/openconfig/lemming/dataplane/forwarding/util/hash/csum16"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	portBytes  = 2  // number of bytes in the TCP port
	srcOffset  = 0  // offset in bytes of the TCP source port
	dstOffset  = 2  // offset in bytes of the TCP destination port
	attrBytes  = 2  // number of bytes in the TCP attributes
	attrOffset = 12 // offset in bytes of the TCP attributes
	flagBits   = 9  // number of bits in the TCP flags
	flagPos    = 0  // offset in bits of the TCP flags
	dataBits   = 4  // number of bits in the TCP data offset
	dataPos    = 12 // offset in bits of the TCP data offset
	csumBytes  = 2  // number of bytes in the TCP checksum
	csumOffset = 16 // offset in bytes of the TCP checksum
	tcpBytes   = 20 // number of bytes in a TCP header

	protoTCP = 6 // TCP packet.
)

// A TCP represents a TCP header in the packet. It can parse, query and
// operate on the TCP header (including TCP extensions). However it cannot
// add a new TCP header to a packet and it cannot interpret, add or remove
// TCP extensions. All TCP payload is considered OPAQUE.
type TCP struct {
	header frame.Header   // TCP header.
	desc   *protocol.Desc // Protocol descriptor.
}

// Header returns the TCP header as a slice of bytes.
func (tcp *TCP) Header() []byte {
	return tcp.header
}

// Trailer returns the no trailing bytes.
func (TCP) Trailer() []byte {
	return nil
}

// ID returns the TCP protocol header ID.
func (TCP) ID(int) fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_PACKET_HEADER_ID_TCP
}

// field returns bytes within the TCP header as identified by id.
func (tcp *TCP) field(id fwdpacket.FieldID) frame.Field {
	if id.IsUDF {
		return protocol.UDF(tcp.header, id)
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC:
		return tcp.header.Field(srcOffset, portBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST:
		return tcp.header.Field(dstOffset, portBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TCP_FLAGS:
		return tcp.header.Field(attrOffset, attrBytes)

	default:
		return nil
	}
}

// Field finds bytes within the TCP header.
func (tcp *TCP) Field(id fwdpacket.FieldID) ([]byte, error) {
	field := tcp.field(id)
	if field == nil {
		return nil, fmt.Errorf("tcp: Field failed, field %v does not exist", id)
	}

	if !id.IsUDF && id.Num == fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TCP_FLAGS {
		return field.BitField(flagPos, flagBits), nil
	}
	return field.Copy(), nil
}

// UpdateField sets bytes within the TCP header.
func (tcp *TCP) UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	if op != fwdpacket.OpSet {
		return false, fmt.Errorf("tcp: UpdateField failed, unsupported op %v for field %v", op, id)
	}

	field := tcp.field(id)
	if field == nil {
		return false, fmt.Errorf("tcp: UpdateField failed, field %v does not exist", id)
	}

	if !id.IsUDF && id.Num == fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TCP_FLAGS {
		field.SetBits(flagPos, flagBits, uint64(binary.BigEndian.Uint16(arg)))
		return true, nil
	}
	return true, field.Set(arg)
}

// Remove removes the TCP header.
func (tcp *TCP) Remove(id fwdpb.PacketHeaderId) error {
	if id != fwdpb.PacketHeaderId_PACKET_HEADER_ID_TCP {
		return fmt.Errorf("tcp: Remove header %v failed, outermost header is %v", id, fwdpb.PacketHeaderId_PACKET_HEADER_ID_TCP)
	}
	tcp.header = nil
	return nil
}

// Modify returns an error as the TCP header does not support adding extensions.
func (TCP) Modify(fwdpb.PacketHeaderId) error {
	return errors.New("tcp: Modify failed, unsupported operation")
}

// ChecksumIPv4 computes the TCP header checksum for IP4 packets.
// The checksum is computed over the TCP header, TCP payload and the source and
// destination IP addresses. A couple of notes:
// 1. It is assumed that the checksum in the TCP header is already zeroed.
// 2. The IP addresses can be padded. As the IPv4 addresses are padded, we need to truncate them.
func ChecksumIPv4(tcpHeader, tcpPayload, ipSrc, ipDst []byte) uint {
	var sum csum16.Sum
	sum.Write(fwdpacket.Truncate(ipSrc, protocol.SizeIP4))
	sum.Write(fwdpacket.Truncate(ipDst, protocol.SizeIP4))
	sum.Write([]byte{0, protoTCP})
	f := make([]byte, protocol.SizeUint16)
	binary.BigEndian.PutUint16(f, uint16(len(tcpHeader)+len(tcpPayload)))
	sum.Write(f)
	sum.Write(tcpHeader)
	sum.Write(tcpPayload)
	return uint(sum)
}

// ChecksumIPv6 computes the TCP header checksum for IP6 packets.
// The checksum is computed over the TCP header, TCP payload and the source and
// destination IP addresses. A couple of notes:
// 1. It is assumed that the checksum in the TCP header is already zeroed.
// 2. It is assumed that the IPv6 addresses are the correct padded, we need to truncate them.
func ChecksumIPv6(tcpHeader, tcpPayload, ipSrc, ipDst []byte) uint {
	var sum csum16.Sum
	sum.Write(fwdpacket.Truncate(ipSrc, protocol.SizeIP6))
	sum.Write(fwdpacket.Truncate(ipDst, protocol.SizeIP6))
	f := make([]byte, protocol.SizeUint32)
	binary.BigEndian.PutUint32(f, uint32(len(tcpHeader)+len(tcpPayload)))
	sum.Write([]byte{0, 0, 0, protoTCP})
	sum.Write(f)
	sum.Write(tcpHeader)
	sum.Write(tcpPayload)
	return uint(sum)
}

// Rebuild update the TCP header checksum using the update L3.
// The checksum is computed over the TCP header, the Opaque payload and various
// IP fields corresponding to the innermost IP header.
func (tcp *TCP) Rebuild() error {
	envelopeID := tcp.desc.EnvelopeID()

	// Update the checksum only if the TCP is over a valid IP header
	switch envelopeID {
	case fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4, fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6:
		break
	default:
		return nil
	}

	// Skip the checksum update if there are no changes.
	e := tcp.desc.EnvelopeDesc()
	p := tcp.desc.PayloadDesc()
	if e != nil && !e.Dirty() && p != nil && !p.Dirty() && !tcp.desc.Dirty() {
		return nil
	}

	// Reset the checksum in the TCP header.
	tcp.header.Field(csumOffset, csumBytes).SetValue(0)
	ipSrc, err := tcp.desc.Packet.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, fwdpacket.LastField))
	if err != nil {
		return fmt.Errorf("tcp: Rebuild failed, err %v", err)
	}

	ipDst, err := tcp.desc.Packet.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, fwdpacket.LastField))
	if err != nil {
		return fmt.Errorf("tcp: Rebuild failed, err %v", err)
	}

	var sum uint
	switch envelopeID {
	case fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4:
		sum = ChecksumIPv4(tcp.header, tcp.desc.Payload(), ipSrc, ipDst)

	case fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6:
		sum = ChecksumIPv6(tcp.header, tcp.desc.Payload(), ipSrc, ipDst)
	}
	tcp.header.Field(csumOffset, csumBytes).SetValue(sum)
	return nil
}

// parse parses a TCP header in the packet.
// It peeks into the TCP header to determine TCP header size (with extensions).
// The payload of TCP is handled as an OPAQUE header.
func parse(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	attr, err := frame.Peek(attrOffset, attrBytes)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("tcp: parsing failed: %v", err)
	}
	size := 4 * attr.BitField(dataPos, dataBits).Value()
	if size < tcpBytes {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("tcp: parse failed, data-offset has small size %v", size)
	}
	header, err := frame.ReadHeader(int(size))
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("tcp: unable read header: %v", err)
	}
	return &TCP{
		header: header,
		desc:   desc,
	}, fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE, nil
}

func init() {
	// TCP header cannot be added to a packet.
	protocol.Register(fwdpb.PacketHeaderId_PACKET_HEADER_ID_TCP, parse, nil)
}
