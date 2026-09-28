// Copyright 2024 Google LLC
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

// Package vxlan implements the VXLAN header support in Lucius.
package vxlan

import (
	"errors"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	"github.com/openconfig/lemming/dataplane/forwarding/util/hash/crc16"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	// DefaultPort is the default UDP port for VXLAN (4789).
	DefaultPort = 4789
	vxlanBytes  = 8
	flagsOffset = 0
	flagsBytes  = 1
	vniOffset   = 4
	vniBytes    = 3 // SizeUint24
	iFlag       = 0x08

	// EphemeralPortBase is the IANA dynamic/private port range base (49152) per RFC 7348 §5.
	EphemeralPortBase = 0xC000
	// EphemeralPortMask is the 14-bit mask for the dynamic port range (16383).
	EphemeralPortMask = 0x3FFF
)

// VXLAN represents a VXLAN header in the packet.
type VXLAN struct {
	header frame.Header
	desc   *protocol.Desc
}

// Header returns the VXLAN header.
func (v *VXLAN) Header() []byte {
	return v.header
}

// Trailer returns nil.
func (VXLAN) Trailer() []byte {
	return nil
}

// ID returns the VXLAN protocol header ID.
func (VXLAN) ID(int) fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN
}

func (v *VXLAN) field(id fwdpacket.FieldID) frame.Field {
	if id.IsUDF {
		return protocol.UDF(v.header, id)
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI:
		return v.header.Field(vniOffset, vniBytes)
	default:
		return nil
	}
}

// Field finds bytes within the VXLAN header.
func (v *VXLAN) Field(id fwdpacket.FieldID) ([]byte, error) {
	if field := v.field(id); field != nil {
		return field.Copy(), nil
	}
	return nil, fmt.Errorf("vxlan: Field failed, field %v does not exist", id)
}

// UpdateField sets bytes within the VXLAN header.
func (v *VXLAN) UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	field := v.field(id)
	if field == nil {
		return false, fmt.Errorf("vxlan: UpdateField failed, field %v does not exist", id)
	}
	if op != fwdpacket.OpSet {
		return false, fmt.Errorf("vxlan: UpdateField failed, unsupported op %v for field %v", op, id)
	}
	if id.IsUDF {
		return true, field.Set(arg)
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI:
		return true, field.Set(arg)
	default:
		return false, fmt.Errorf("vxlan: UpdateField failed, field %v not supported", id)
	}
}

// Remove removes the VXLAN header.
func (v *VXLAN) Remove(id fwdpb.PacketHeaderId) error {
	if id != fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN {
		return fmt.Errorf("vxlan: Remove header %v failed, outermost header is %v", id, fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN)
	}
	v.header = nil
	return nil
}

// Modify returns an error as VXLAN does not support extensions.
func (VXLAN) Modify(_ fwdpb.PacketHeaderId) error {
	return errors.New("vxlan: Modify is unsupported")
}

// Rebuild rebuilds the VXLAN header.
func (v *VXLAN) Rebuild() error {
	return nil
}

// InnerFlowKey extracts inner header invariant fields (5-tuple and MACs) per RFC 7348.
// It queries Instance 1 (standard inner frame of a VXLAN-encapsulated packet).
// Only invariant header fields are included;
// (such as IP ID or TCP sequence numbers) are excluded to ensure deterministic ECMP flow stability.
func InnerFlowKey(desc *protocol.Desc) []byte {
	if desc == nil || desc.Packet == nil {
		return nil
	}
	// Inner fields: MACs (12B) + IPv6/IPv4 IPs (32B) + Proto (1B) + L4 Ports (4B)
	key := make([]byte, 0, 49)
	fields := []fwdpb.PacketFieldNum{
		fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC,
		fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST,
		fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC,
		fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
		fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO,
		fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC,
		fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST,
	}
	for _, num := range fields {
		if f, err := desc.Packet.Field(fwdpacket.NewFieldIDFromNum(num, 1)); err == nil && len(f) > 0 {
			key = append(key, f...)
		}
	}
	return key
}

// ComputeEntropySourcePort calculates a 16-bit UDP source port within the dynamic range [49152, 65535]
// based on the inner packet headers per RFC 7348 §5. Returns 0 if no inner header fields are found.
func ComputeEntropySourcePort(desc *protocol.Desc) uint16 {
	key := InnerFlowKey(desc)
	if len(key) == 0 {
		return 0
	}
	hashVal := crc16.ChecksumANSI(key)
	return uint16(EphemeralPortBase | (hashVal & EphemeralPortMask))
}

// Parse parses a VXLAN header in the packet.
func Parse(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	return parse(frame, desc)
}

// parse parses a VXLAN header in the packet.
func parse(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	if frame.Len() < vxlanBytes {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("vxlan: parse failed, frame length %v too small to contain a VXLAN header", frame.Len())
	}
	header, err := frame.ReadHeader(vxlanBytes)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("vxlan: unable read header: %v", err)
	}
	return &VXLAN{
		header: header,
		desc:   desc,
	}, fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, nil
}

// add adds a VXLAN header to the packet.
func add(_ fwdpb.PacketHeaderId, desc *protocol.Desc) (protocol.Handler, error) {
	header := make(frame.Header, vxlanBytes)
	header.Field(flagsOffset, flagsBytes).SetValue(iFlag)
	return &VXLAN{
		header: header,
		desc:   desc,
	}, nil
}

func init() {
	protocol.Register(fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN, parse, add)
}
