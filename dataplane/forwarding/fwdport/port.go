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

// Package fwdport contains routines and types to manage forwarding ports.
//
// A Port is an entry or exit point for packets within the forwarding plane.
// Lucius has different types of ports (implemented by various types). All ports
// are created by provisioning.
//
// This package defines the following mechanisms to manage ports
//  1. It allows different types of ports to register builders during
//     package initialization.
//  2. Provisioning can create ports using the registered builders.
//  3. It defines an interface that can be used to operate on a port.
package fwdport

import (
	"encoding/binary"
	"errors"
	"fmt"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdattribute"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// CounterList is a set of counters incremented by ports.
var CounterList = []fwdpb.CounterId{
	fwdpb.CounterId_COUNTER_ID_RX_PACKETS,
	fwdpb.CounterId_COUNTER_ID_RX_OCTETS,
	fwdpb.CounterId_COUNTER_ID_RX_BAD_PACKETS,
	fwdpb.CounterId_COUNTER_ID_RX_BAD_OCTETS,
	fwdpb.CounterId_COUNTER_ID_RX_ADMIN_DROP_PACKETS,
	fwdpb.CounterId_COUNTER_ID_RX_ADMIN_DROP_OCTETS,
	fwdpb.CounterId_COUNTER_ID_RX_ERROR_PACKETS,
	fwdpb.CounterId_COUNTER_ID_RX_ERROR_OCTETS,
	fwdpb.CounterId_COUNTER_ID_RX_DROP_PACKETS,
	fwdpb.CounterId_COUNTER_ID_RX_DROP_OCTETS,
	fwdpb.CounterId_COUNTER_ID_RX_DEBUG_PACKETS,
	fwdpb.CounterId_COUNTER_ID_RX_DEBUG_OCTETS,
	fwdpb.CounterId_COUNTER_ID_TX_PACKETS,
	fwdpb.CounterId_COUNTER_ID_TX_OCTETS,
	fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS,
	fwdpb.CounterId_COUNTER_ID_TX_ERROR_OCTETS,
	fwdpb.CounterId_COUNTER_ID_TX_DROP_PACKETS,
	fwdpb.CounterId_COUNTER_ID_TX_DROP_OCTETS,
	fwdpb.CounterId_COUNTER_ID_TX_ADMIN_DROP_PACKETS,
	fwdpb.CounterId_COUNTER_ID_TX_ADMIN_DROP_OCTETS,
	fwdpb.CounterId_COUNTER_ID_TX_UCAST_PACKETS,
	fwdpb.CounterId_COUNTER_ID_TX_NON_UCAST_PACKETS,
	fwdpb.CounterId_COUNTER_ID_RX_UCAST_PACKETS,
	fwdpb.CounterId_COUNTER_ID_RX_NON_UCAST_PACKETS,
}

// A Port is an entry or exit point within the forwarding plane. Each port
// has a set of actions to apply for incoming and outgoing packets.
//
// Ports are always created by provisioning.
type Port interface {
	fwdobject.Object

	// Update updates a port.
	Update(upd *fwdpb.PortUpdateDesc) error

	// Write writes a packet out. If successful, the port returns
	// fwdaction.CONSUME.
	Write(packet fwdpacket.Packet) (fwdaction.State, error)

	// Actions returns the port actions of the specified type
	Actions(dir fwdpb.PortAction) fwdaction.Actions

	// State manages the state of the port.
	State(op *fwdpb.PortInfo) (*fwdpb.PortStateReply, error)
}

// A Builder can build Ports of the specified type.
type Builder interface {
	// Build builds a port.
	Build(portDesc *fwdpb.PortDesc, ctx *fwdcontext.Context) (Port, error)
}

// builders is a map of builders for various types of ports.
var builders = make(map[fwdpb.PortType]Builder)

// Register registers a builder for a port type. Note that builders are
// expected to be registered during package initialization.
func Register(portType fwdpb.PortType, builder Builder) {
	builders[portType] = builder
}

// New creates a new port.
func New(desc *fwdpb.PortDesc, ctx *fwdcontext.Context) (Port, error) {
	if desc == nil {
		return nil, errors.New("fwdport: new failed, missing description")
	}
	builder, ok := builders[desc.GetPortType()]
	if !ok {
		return nil, fmt.Errorf("fwdport: new failed, no builder for port %s", desc)
	}

	port, err := builder.Build(desc, ctx)
	if err != nil {
		return nil, err
	}
	pid := desc.GetPortId()
	if pid == nil {
		return nil, errors.New("fwdport: new failed, missing id")
	}
	if err = ctx.Objects.Insert(port, pid.GetObjectId()); err != nil {
		return nil, err
	}
	return port, nil
}

// Find finds a port.
func Find(id *fwdpb.PortId, ctx *fwdcontext.Context) (Port, error) {
	if id == nil {
		return nil, errors.New("fwdport: find failed, no port specified")
	}
	object, err := ctx.Objects.FindID(id.GetObjectId())
	if err != nil {
		return nil, err
	}
	if port, ok := object.(Port); ok {
		return port, nil
	}
	return nil, fmt.Errorf("fwdport: find failed, %v is not a port", id)
}

// Acquire acquires a reference to a port.
func Acquire(id *fwdpb.PortId, ctx *fwdcontext.Context) (Port, error) {
	if id == nil {
		return nil, errors.New("fwdport: acquire failed, no port specified")
	}
	object, err := ctx.Objects.Acquire(id.GetObjectId())
	if err != nil {
		return nil, err
	}
	if port, ok := object.(Port); ok {
		return port, nil
	}

	// Release the object if there was an error.
	_ = object.Release(false /*forceCleanup*/)
	return nil, fmt.Errorf("fwdport: acquire failed, %v is not a port", id)
}

// Release releases a reference to a port.
func Release(port Port) error {
	if port == nil {
		return errors.New("fwdport: release failed, no port specified")
	}
	return port.Release(false /*forceCleanup*/)
}

// MakeID makes a PortID corresponding to the specified object ID.
func MakeID(id fwdobject.ID) *fwdpb.PortId {
	return &fwdpb.PortId{
		ObjectId: fwdobject.MakeID(id),
	}
}

// GetID returns the PortID for the given port.
func GetID(port Port) *fwdpb.PortId {
	return MakeID(port.ID())
}

// findPort returns the port associated with the specifed field of the packet.
func findPort(packet fwdpacket.Packet, ctx *fwdcontext.Context, num fwdpb.PacketFieldNum) (Port, error) {
	field, err := packet.Field(fwdpacket.NewFieldIDFromNum(num, 0))
	if err != nil {
		return nil, err
	}
	if len(field) != protocol.SizeUint64 {
		return nil, fmt.Errorf("fwdport: field %v has bad size %v", num, len(field))
	}
	nid := fwdobject.NID(binary.BigEndian.Uint64(field))
	obj, err := ctx.Objects.FindNID(nid)
	if err != nil {
		return nil, fmt.Errorf("fwdport: field %v error, %v", num, err)
	}
	if port, ok := obj.(Port); ok {
		return port, nil
	}
	return nil, fmt.Errorf("fwdport: field %v has non-port object %v", num, obj.ID())
}

// InputPort returns the input port of the packet.
func InputPort(packet fwdpacket.Packet, ctx *fwdcontext.Context) (Port, error) {
	return findPort(packet, ctx, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT)
}

// OutputPort returns the output port of the packet.
func OutputPort(packet fwdpacket.Packet, ctx *fwdcontext.Context) (Port, error) {
	return findPort(packet, ctx, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT)
}

// SetInputPort initializes the input port of the packet.
func SetInputPort(packet fwdpacket.Packet, port Port) {
	fwdpacket.SetNID(packet, port.NID(), fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT)
}

// SetOutputPort initializes the output port of the packet.
func SetOutputPort(packet fwdpacket.Packet, port Port) {
	fwdpacket.SetNID(packet, port.NID(), fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT)
}

// Increment increments a packet and octet counters on the port.
func Increment(port Port, octets int, packetID, octetID fwdpb.CounterId) {
	port.Increment(packetID, 1)
	port.Increment(octetID, uint32(octets))
}

// Input processes an incoming packet. The specified port actions are applied
// to the packet, and if the packet has an output port, output actions are
// performed on the packet. All appropriate counters are incremented.
func Input(port Port, packet fwdpacket.Packet, dir fwdpb.PortAction, ctx *fwdcontext.Context) (err error) {
	defer func() {
		if err != nil {
			packet.Log().Error(err, "input processing failed", "frame", fwdpacket.IncludeFrameInLog)
		}
	}()

	Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_RX_PACKETS, fwdpb.CounterId_COUNTER_ID_RX_OCTETS)
	SetInputPort(packet, port)
	mac, err := packet.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0))
	if err != nil {
		Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_RX_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_RX_ERROR_OCTETS)
		return err
	}
	if mac[0]%2 == 0 { // Unicast address is when is least significant bit of the 1st octet is 0.
		port.Increment(fwdpb.CounterId_COUNTER_ID_RX_UCAST_PACKETS, 1)
	} else {
		port.Increment(fwdpb.CounterId_COUNTER_ID_RX_NON_UCAST_PACKETS, 1)
	}

	packet.Log().V(3).Info("input packet", "port", port.ID(), "frame", fwdpacket.IncludeFrameInLog)
	state, err := fwdaction.ProcessPacket(packet, port.Actions(dir), port)
	if err != nil {
		Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_RX_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_RX_ERROR_OCTETS)
		return err
	}

	switch state {
	case fwdaction.DROP:
		packet.Log().V(1).Info("input dropped frame", "port", port.ID(), "frame", fwdpacket.IncludeFrameInLog)
		Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_RX_DROP_PACKETS, fwdpb.CounterId_COUNTER_ID_RX_DROP_OCTETS)
		return nil

	case fwdaction.CONSUME:
		return nil

	case fwdaction.OUTPUT, fwdaction.CONTINUE:
		// If we don't have an output port, count it as a drop.
		out, err := OutputPort(packet, ctx)
		if err != nil {
			Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_RX_DROP_PACKETS, fwdpb.CounterId_COUNTER_ID_RX_DROP_OCTETS)
			return nil
		}
		packet.Log().V(3).Info("transmitting packet", "port", out.ID())
		packet.Log().WithValues("context", ctx.ID, "port", out.ID())
		Output(out, packet, fwdpb.PortAction_PORT_ACTION_OUTPUT, ctx)
		return nil
	}
	return fmt.Errorf("fwdport: unknown state %v", state)
}

// Output processes an outgoing packet. The specified port actions are applied
// to the packet, and if allowed the packet is written out of the port. All
// appropriate counters are incremented.
func Output(port Port, packet fwdpacket.Packet, dir fwdpb.PortAction, _ *fwdcontext.Context) (err error) {
	defer func() {
		if err != nil {
			packet.Log().Error(err, "output processing failed")
		}
	}()
	Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_TX_PACKETS, fwdpb.CounterId_COUNTER_ID_TX_OCTETS)
	SetOutputPort(packet, port)
	mac, err := packet.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0))
	if err != nil {
		Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_TX_ERROR_OCTETS)
		return err
	}
	if mac[0]%2 == 0 { // Unicast address is when is least significant bit of the 1st octet is 0.
		port.Increment(fwdpb.CounterId_COUNTER_ID_TX_UCAST_PACKETS, 1)
	} else {
		port.Increment(fwdpb.CounterId_COUNTER_ID_TX_NON_UCAST_PACKETS, 1)
	}

	packet.Log().V(3).Info("output packet", "frame", fwdpacket.IncludeFrameInLog)
	state, err := fwdaction.ProcessPacket(packet, port.Actions(dir), port)
	if err == nil && state == fwdaction.CONTINUE {
		state, err = port.Write(packet)
	}
	if err != nil {
		Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_TX_ERROR_OCTETS)
		return err
	}
	switch state {
	case fwdaction.DROP:
		packet.Log().V(1).Info("output dropped frame", "port", port.ID(), "frame", fwdpacket.IncludeFrameInLog)
		Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_TX_DROP_PACKETS, fwdpb.CounterId_COUNTER_ID_TX_DROP_OCTETS)
		return nil
	case fwdaction.CONSUME:
		packet.Log().V(1).Info("consumed frame", "port", port.ID(), "frame", fwdpacket.IncludeFrameInLog)
		return nil
	case fwdaction.CONTINUE:
		Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_TX_ERROR_OCTETS)
		return errors.New("fwdport: output processing results in continued processing")
	}
	return fmt.Errorf("fwdport: unknown state %v", state)
}

// Write writes out a packet through a port without changing it. No actions are applied.
func Write(port Port, packet fwdpacket.Packet) {
	packet.Log().V(1).Info("write packet", "port", port.ID(), "frame", fwdpacket.IncludeFrameInLog)
	Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_TX_PACKETS, fwdpb.CounterId_COUNTER_ID_TX_OCTETS)
	if _, err := port.Write(packet); err != nil {
		packet.Log().Error(err, "write failed")
		Increment(port, packet.Length(), fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_TX_ERROR_OCTETS)
	}
}

// Process processes a packet on the specified port, action direction and context.
// If the attribute "PacketDebug" is set, then the packet debugging is enabled.
func Process(port Port, packet fwdpacket.Packet, dir fwdpb.PortAction, ctx *fwdcontext.Context, prefix string) {
	// Update the prefix if the port and context are not nil
	if port != nil && ctx != nil {
		prefix = fmt.Sprintf("%v: %v %v", ctx.ID, prefix, port.ID())
	}

	// Initialize packet attributes using the context attributes overridden
	// by the port attributes.
	a := packet.Attributes()
	if a != nil {
		a.Override(fwdattribute.Global)
		a.Override(ctx.Attributes)
		a.Override(port.Attributes())
	}

	// Setup the packet level debugs.
	if value, ok := a.Get(fwdpacket.AttrPacketDebug); ok && value == "true" {
		packet.Debug(true)
	}

	packet.Log().WithName(prefix)
	switch dir {
	case fwdpb.PortAction_PORT_ACTION_INPUT:
		Input(port, packet, dir, ctx)

	case fwdpb.PortAction_PORT_ACTION_OUTPUT:
		Output(port, packet, dir, ctx)

	case fwdpb.PortAction_PORT_ACTION_WRITE:
		Write(port, packet)

	default:
		log.Errorf("%v: unknown action %v.", prefix, dir)
		return
	}
	fwdpacket.Log(packet)
}
