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

// Package metadata implements the metadata packet header.
package metadata

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Metadata is a manadatory packet header that stores derived state of a packet.
// Since metadata is not a real packet header, it does not contribute to the
// packet's frame.
type Metadata struct {
	length           uint64           // Packet length.
	inputPort        []byte           // Input port identifier.
	outputPort       []byte           // Output port identifier.
	vrf              []byte           // VRF identifier.
	attribute32      map[uint8][]byte // Map of 32-bit attributes indexed by instance.
	attribute24      map[uint8][]byte // Map of 24-bit attributes indexed by instance.
	attribute16      map[uint8][]byte // Map of 16-bit attributes indexed by instance.
	attribute8       map[uint8][]byte // Map of 8-bit attributes indexed by instance.
	nextHopIP        []byte
	nextHopID        []byte         // ID of the next hop.
	nextHopGroupID   []byte         // ID of the next hop group.
	trapID           []byte         // ID of the trap rule that was applies to this packet.
	inputIface       []byte         // L3 input interface id.
	outputIface      []byte         // L3 output interface id.
	tunnelID         []byte         // Tunnel ID
	hostPortID       []byte         // Host port id
	l2mcGroupID      []byte         // L2MC Group ID
	policer          []byte         // Policer ID
	targetEgressPort []byte         // Target egress port
	packetAction     []byte         // Action that packet is taking.
	desc             *protocol.Desc // Protocol descriptor.
}

// Header returns nil as it does not contribute to the packet's frame.
func (Metadata) Header() []byte {
	return nil
}

// Trailer returns the no trailing bytes.
func (Metadata) Trailer() []byte {
	return nil
}

// ID returns the protocol header ID.
func (Metadata) ID(int) fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_PACKET_HEADER_ID_METADATA
}

// Field returns the values of the queried packet fields.
func (m *Metadata) Field(id fwdpacket.FieldID) ([]byte, error) {
	if id.IsUDF {
		return nil, fmt.Errorf("metadata: Field %v failed, unsupported field", id)
	}

	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT:
		return m.inputPort, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT:
		return m.outputPort, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF:
		vrf := make([]byte, protocol.FieldAttr[id.Num].DefaultSize)
		copy(vrf, m.vrf)
		return vrf, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_LENGTH:
		length := make([]byte, protocol.FieldAttr[id.Num].DefaultSize)
		binary.BigEndian.PutUint64(length, m.length)
		return length, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_32:
		if a, ok := m.attribute32[id.Instance]; ok {
			return a, nil
		}
		return make([]byte, 4), nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_24:
		if a, ok := m.attribute24[id.Instance]; ok {
			return a, nil
		}
		return make([]byte, 3), nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_16:
		if a, ok := m.attribute16[id.Instance]; ok {
			return a, nil
		}
		return make([]byte, 2), nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_8:
		if a, ok := m.attribute8[id.Instance]; ok {
			return a, nil
		}
		return make([]byte, 1), nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP:
		return m.nextHopIP, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_GROUP_ID:
		return m.nextHopGroupID, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID:
		return m.nextHopID, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_INPUT_IFACE:
		return m.inputIface, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE:
		return m.outputIface, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TRAP_ID:
		return m.trapID, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID:
		return m.tunnelID, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_HOST_PORT_ID:
		return m.hostPortID, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID:
		return m.l2mcGroupID, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_POLICER_ID:
		return m.policer, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TARGET_EGRESS_PORT:
		return m.targetEgressPort, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION:
		return m.packetAction, nil
	default:
		return nil, fmt.Errorf("metadata: Field %v failed, unsupported field", id)
	}
}

// updateNumeric updates a returns a numeric field after
// incrementing/decrementing it. Note that if a field is nil, an empty field
// is made of the appropriate length. It is assumed that the inputs are valid.
func updateNumeric(field, arg frame.Field, length, op int) frame.Field {
	if len(field) == 0 {
		field = frame.Field(make([]byte, length))
	}
	a := arg.Value()
	n := field.Value()
	switch op {
	case fwdpacket.OpDec:
		n -= a
	case fwdpacket.OpInc:
		n += a
	}
	field.SetValue(n)
	return field
}

// updateIncDec performs an Inc or Dec on a numeric field.
func (m *Metadata) updateIncDec(id fwdpacket.FieldID, arg frame.Field, op int) (bool, error) {
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_32:
		m.attribute32[id.Instance] = updateNumeric(m.attribute32[id.Instance], arg, 4, op)
		return true, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_16:
		m.attribute16[id.Instance] = updateNumeric(m.attribute16[id.Instance], arg, 2, op)
		return true, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_8:
		m.attribute8[id.Instance] = updateNumeric(m.attribute8[id.Instance], arg, 1, op)
		return true, nil
	default:
		return false, fmt.Errorf("metadata: UpdateField failed, unsupported inc/dec op %v for field %v", op, id)
	}
}

// updateSet sets a metadata field to the specified value.
func (m *Metadata) updateSet(id fwdpacket.FieldID, arg []byte) (bool, error) {
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT:
		m.inputPort = arg
		return true, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT:
		m.outputPort = arg
		return true, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF:
		copy(m.vrf, arg)
		return true, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_32:
		a := make([]byte, 4)
		copy(a, arg)
		m.attribute32[id.Instance] = a
		return true, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_24:
		a := make([]byte, 3)
		copy(a, arg)
		m.attribute24[id.Instance] = a
		return true, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_16:
		a := make([]byte, 2)
		copy(a, arg)
		m.attribute16[id.Instance] = a
		return true, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_8:
		a := make([]byte, 1)
		copy(a, arg)
		m.attribute8[id.Instance] = a
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP:
		m.nextHopIP = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID:
		m.nextHopID = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_GROUP_ID:
		m.nextHopGroupID = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TRAP_ID:
		m.trapID = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_INPUT_IFACE:
		m.inputIface = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE:
		m.outputIface = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID:
		m.tunnelID = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_HOST_PORT_ID:
		m.hostPortID = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID:
		m.l2mcGroupID = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_POLICER_ID:
		m.policer = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TARGET_EGRESS_PORT:
		m.targetEgressPort = arg
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION:
		m.packetAction = arg
		return true, nil
	default:
		return false, fmt.Errorf("metadata: UpdateField failed, set unsupported for field %v", id)
	}
}

// UpdateField can update the input and output port of the packet.
func (m *Metadata) UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	if id.IsUDF {
		return false, fmt.Errorf("metadata: UpdateField failed, unsupported op %v for UDF %v", op, id)
	}

	switch op {
	case fwdpacket.OpSet:
		return m.updateSet(id, arg)

	case fwdpacket.OpInc, fwdpacket.OpDec:
		return m.updateIncDec(id, frame.Field(arg), op)

	default:
		return false, fmt.Errorf("metadata: UpdateField failed, unsupported op %v for field %v", op, id)
	}
}

// Rebuild updates the packet length. The length is zero if the packet had
// zero bytes.
func (m *Metadata) Rebuild() error {
	m.length = uint64(m.desc.PayloadLength())
	return nil
}

// Remove returns an error as the metadata header cannot be removed.
func (Metadata) Remove(fwdpb.PacketHeaderId) error {
	return errors.New("metadata: Remove is unsupported")
}

// Modify returns an error as the metadata header cannot be changed.
func (Metadata) Modify(fwdpb.PacketHeaderId) error {
	return errors.New("metadata: Modify is unsupported")
}

// parse accepts the frame and initializes metadata state. It always returns
// the next header as "NONE" because the first header is decided by the parser
// (and the port).
func parse(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	return &Metadata{
		desc:             desc,
		length:           uint64(frame.Len()),
		vrf:              make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF].DefaultSize),
		inputPort:        make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT].DefaultSize),
		outputPort:       make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT].DefaultSize),
		nextHopIP:        make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP].DefaultSize),
		inputIface:       make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_INPUT_IFACE].DefaultSize),
		outputIface:      make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE].DefaultSize),
		tunnelID:         make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID].DefaultSize),
		hostPortID:       make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_HOST_PORT_ID].DefaultSize),
		l2mcGroupID:      make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID].DefaultSize),
		policer:          make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_POLICER_ID].DefaultSize),
		targetEgressPort: make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TARGET_EGRESS_PORT].DefaultSize),
		packetAction:     make([]byte, protocol.FieldAttr[fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION].DefaultSize),
		attribute32:      make(map[uint8][]byte),
		attribute24:      make(map[uint8][]byte),
		attribute16:      make(map[uint8][]byte),
		attribute8:       make(map[uint8][]byte),
	}, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, nil
}

func init() {
	// Register the parse function for the METADATA headers.
	// Note that metadata cannot be added explicitly.
	protocol.Register(fwdpb.PacketHeaderId_PACKET_HEADER_ID_METADATA, parse, nil)
}
