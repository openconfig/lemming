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

// Package udp implements the UDP header support in Lucius.
package udp

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/vxlan"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	"github.com/openconfig/lemming/dataplane/forwarding/util/hash/csum16"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	portBytes  = 2 // number of bytes in the UDP port
	srcOffset  = 0 // offset in bytes of the UDP source port
	dstOffset  = 2 // offset in bytes of the UDP destination port
	csumBytes  = 2 // number of bytes in the UDP checksum
	csumOffset = 6 // offset in bytes of the UDP checksum
	lenBytes   = 2 // number of bytes in the UDP length
	lenOffset  = 4 // offset in bytes of the UDP length
	udpBytes   = 8 // number of bytes in an udp header

	protoUDP = 17 // UDP packet
)

// An UDP represents a UDP header in the packet. It can add, remove and update
// the UDP header. The UDP payload is treated as OPAQUE.
type UDP struct {
	header frame.Header
	vxlan  protocol.Handler
	desc   *protocol.Desc
}

// Header returns the UDP header (including optional embedded VXLAN).
func (udp *UDP) Header() []byte {
	header := append([]byte(nil), udp.header...)
	if udp.vxlan != nil {
		header = append(header, udp.vxlan.Header()...)
	}
	return header
}

// Trailer returns the no trailing bytes.
func (UDP) Trailer() []byte {
	return nil
}

// ID returns the UDP or VXLAN protocol header ID.
func (udp *UDP) ID(instance int) fwdpb.PacketHeaderId {
	if instance > 0 && udp.vxlan != nil {
		return udp.vxlan.ID(instance - 1)
	}
	return fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP
}

// field returns bytes within the UDP header as identified by id.
func (udp *UDP) field(id fwdpacket.FieldID) frame.Field {
	if id.IsUDF {
		return protocol.UDF(udp.header, id)
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC:
		return udp.header.Field(srcOffset, portBytes)

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST:
		return udp.header.Field(dstOffset, portBytes)

	default:
		return nil
	}
}

// Field finds bytes within the UDP or VXLAN header.
func (udp *UDP) Field(id fwdpacket.FieldID) ([]byte, error) {
	if id.Num == fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI {
		if udp.vxlan != nil {
			return udp.vxlan.Field(id)
		}
		return nil, fmt.Errorf("udp: Field failed, VXLAN header does not exist")
	}
	if field := udp.field(id); field != nil {
		return field.Copy(), nil
	}
	return nil, fmt.Errorf("udp: Field failed, field %v does not exist", id)
}

// UpdateField sets bytes within the UDP or VXLAN header.
func (udp *UDP) UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	if id.Num == fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI {
		if udp.vxlan != nil {
			return udp.vxlan.UpdateField(id, op, arg)
		}
		return false, fmt.Errorf("udp: UpdateField failed, VXLAN header does not exist")
	}
	if field := udp.field(id); field != nil && op == fwdpacket.OpSet {
		return true, field.Set(arg)
	}
	return false, fmt.Errorf("udp: UpdateField failed, unsupported op %v for field %v", op, id)
}

// Remove removes the UDP or VXLAN header.
func (udp *UDP) Remove(id fwdpb.PacketHeaderId) error {
	switch id {
	case fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN:
		if udp.vxlan != nil {
			err := udp.vxlan.Remove(id)
			udp.vxlan = nil
			return err
		}
		return fmt.Errorf("udp: Remove header %v failed, VXLAN header does not exist", id)
	case fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP:
		udp.header = nil
		return nil
	default:
		return fmt.Errorf("udp: Remove header %v failed, outermost header is %v", id, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP)
	}
}

// Modify returns an error as the UDP header has no extensions.
func (UDP) Modify(_ fwdpb.PacketHeaderId) error {
	return errors.New("udp: Modify is unsupported")
}

// Rebuild updates the UDP header length and csum (if over IPv6) and rebuilds VXLAN if present.
func (udp *UDP) Rebuild() error {
	if udp.vxlan != nil {
		if err := udp.vxlan.Rebuild(); err != nil {
			return err
		}
	}

	// If the envelope, payload and udp header are unmodified, skip updates.
	e := udp.desc.EnvelopeDesc()
	p := udp.desc.PayloadDesc()
	if e != nil && !e.Dirty() && p != nil && !p.Dirty() && !udp.desc.Dirty() {
		return nil
	}

	// Update the length and reset the checksum.
	length := len(udp.Header()) + udp.desc.PayloadLength()
	udp.header.Field(lenOffset, lenBytes).SetValue(uint(length))
	udp.header.Field(csumOffset, csumBytes).SetValue(0)

	// If UDP is over IP4, we keep the checksum as zero.
	if udp.desc.EnvelopeID() != fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6 {
		return nil
	}

	// Checksum is computed over the UDP header, IP source and destination
	// address, the packet length and protocol.
	var err error
	var f []byte
	var sum csum16.Sum
	_, _ = sum.Write(udp.Header())
	_, _ = sum.Write(udp.desc.Payload())
	if f, err = udp.desc.Packet.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC, fwdpacket.LastField)); err != nil {
		return fmt.Errorf("udp: Rebuild failed: %v", err)
	}
	_, _ = sum.Write(f)
	if f, err = udp.desc.Packet.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST, fwdpacket.LastField)); err != nil {
		return fmt.Errorf("udp: Rebuild failed: %v", err)
	}
	_, _ = sum.Write(f)

	f = make([]byte, protocol.SizeUint32)
	binary.BigEndian.PutUint32(f, uint32(length))
	_, _ = sum.Write(f)
	_, _ = sum.Write([]byte{0, 0, 0, protoUDP})
	udp.header.Field(csumOffset, csumBytes).SetValue(uint(sum))
	return nil
}

// add adds a UDP header to the packet.
func add(_ fwdpb.PacketHeaderId, desc *protocol.Desc) (protocol.Handler, error) {
	return &UDP{
		header: make([]byte, udpBytes),
		desc:   desc,
	}, nil
}

// parse parses a UDP header in the packet. VXLAN traffic on port 4789 embeds
// a VXLAN header while other UDP payloads default to PACKET_HEADER_ID_OPAQUE.
func parse(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	if frame.Len() < udpBytes {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("udp: parse failed, frame length %v too small to contain a UDP header", frame.Len())
	}
	header, err := frame.ReadHeader(udpBytes)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("udp: unable read header: %v", err)
	}
	dstPort := header.Field(dstOffset, portBytes).Value()
	var vxlanHandler protocol.Handler
	next := fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE
	if dstPort == vxlan.DefaultPort {
		var err error
		vxlanHandler, next, err = vxlan.Parse(frame, desc)
		if err != nil {
			return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, err
		}
	}
	return &UDP{
		header: header,
		vxlan:  vxlanHandler,
		desc:   desc,
	}, next, nil
}

func init() {
	protocol.Register(fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP, parse, add)
}
