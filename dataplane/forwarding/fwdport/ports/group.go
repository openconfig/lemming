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

package ports

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/util/hash/crc16"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// AsyncPacketQueueDepth is the number of packets buffered for the async write
// requests. Heuristic says 10 is sufficient.
const AsyncPacketQueueDepth = 10

// A member represents a member of a port group and the set of actions to be
// applied when the member is selected. It also indicates the number of
// instances the member should be tried for a packet.
type member struct {
	port      fwdport.Port
	actions   fwdaction.Actions
	parent    fwdport.Port
	instances int
	ctx       *fwdcontext.Context

	stop    chan bool             // Channel for stopping goroutines.
	packets chan fwdpacket.Packet // Buffered channel of packets
}

// String formats a member into a string.
func (m *member) String() string {
	return fmt.Sprintf("<Port=%v><Actions=%v><Count=%d>", m.port, m.actions, m.instances)
}

// Write writes a packet out of the member after applying the member actions.
func (m *member) Write(packet fwdpacket.Packet, caller string) error {
	packet.Logf(fwdpacket.LogDesc, fmt.Sprintf("%v: %v %v:", m.ctx.ID, caller, m.port.ID()))
	state, err := fwdaction.ProcessPacket(packet, m.actions, m.port)
	if err != nil {
		return fmt.Errorf("ports: %v unable to write out of member %v: %v ", m.parent, m.port, err)
	}
	if state != fwdaction.CONTINUE {
		return fmt.Errorf("ports: %v unable to write out of member %v: unexpected state %v", m.parent, m.port, state)
	}
	return fwdport.Output(m.port, packet, fwdpb.PortAction_PORT_ACTION_OUTPUT, m.ctx)
}

// AsyncWrite queues a packet to be written asynchronously out of a member.
func (m *member) AsyncWrite(packet fwdpacket.Packet) {
	m.packets <- packet
}

// Cleanup cleans up references held by a member.
func (m *member) Cleanup() {
	fwdport.Release(m.port)
	m.actions.Cleanup()
	close(m.stop)
	close(m.packets)
}

// Ready returns true if the member is ready to transmit packets.
func (m *member) Ready() bool {
	ps, err := m.port.State(nil)
	if err != nil {
		log.Warningf("ports: error querying port state (%v)", err)
		return false
	}
	return ps.GetLink().GetState() == fwdpb.LinkState_LINK_STATE_UP
}

// A portGroup is a port that writes packets to a group of ports. A port can
// appear multiple times in the group. This is used by clients to mimic a
// weighted group.
type portGroup struct {
	fwdobject.Base
	fields   []fwdpacket.FieldID                                    // packet fields used to create a packet hash
	hashFn   func(key []byte) int                                   // function used to hash a set of bytes
	hash     fwdpb.AggregateHashAlgorithm                           // hash algorithm used to select the port
	packetFn func(packet fwdpacket.Packet) (fwdaction.State, error) // function used to process packets
	ctx      *fwdcontext.Context
	// list of members used to hash for packets. If a member has two instances
	// in the group, it appears twice in this array.
	members []*member

	// map of members indexed by the port id.
	memberMap map[fwdobject.ID]*member
}

// String returns the port as a formatted string.
func (p *portGroup) String() string {
	desc := fmt.Sprintf("Type=%v;<Members=%v>;<Fields=%v>;<Hash=%v>;%v", fwdpb.PortType_PORT_TYPE_AGGREGATE_PORT, p.members, p.fields, p.hash, p.BaseInfo())
	if state, err := p.State(nil); err == nil {
		desc += fmt.Sprintf("<State=%v>;", state)
	}
	return desc
}

// Cleanup releases references held by the ports and its entries.
func (p *portGroup) Cleanup() {
	for _, m := range p.memberMap {
		m.Cleanup()
	}
	p.memberMap = nil
	p.members = nil
}

// recompute recomputes the members list of the port group from its
// current member map.
func (p *portGroup) recompute() {
	p.members = nil
	for _, m := range p.memberMap {
		for index := 0; index < m.instances; index++ {
			p.members = append(p.members, m)
		}
	}
}

// removeGroupMember removes the specified port id from the port group. Note
// that all instances of the member are removed.
func (p *portGroup) removeGroupMember(u *fwdpb.AggregatePortRemoveMemberUpdateDesc) error {
	pid := fwdobject.ID(u.GetPortId().GetObjectId().GetId())
	m, ok := p.memberMap[pid]
	if !ok {
		return fmt.Errorf("ports: Unable to find port %v in %v", pid, p)
	}
	delete(p.memberMap, pid)
	p.recompute()
	m.Cleanup()
	return nil
}

// newGroupMember creates a new member.
func newGroupMember(pid *fwdpb.PortId, actions fwdaction.Actions, ctx *fwdcontext.Context, parent fwdport.Port, instances int) (*member, error) {
	port, err := fwdport.Acquire(pid, ctx)
	if err != nil {
		return nil, fmt.Errorf("ports: Unable to create port group, missing member: %v", err)
	}

	m := &member{
		port:      port,
		parent:    parent,
		actions:   actions,
		instances: instances,
		ctx:       ctx,
		stop:      make(chan bool),
		packets:   make(chan fwdpacket.Packet, AsyncPacketQueueDepth),
	}
	go func(m *member) {
		for {
			select {
			case <-m.stop:
				log.Infof("ports: stopping async write goroutine for member %v", m)
				return
			case packet := <-m.packets:
				if packet == nil {
					continue
				}
				if err := m.Write(packet, "async"); err != nil {
					packet.Logf(fwdpacket.LogErrorFrame, "Async write to %v failed, err %v", m, err)
				}
			}
		}
	}(m)
	return m, nil
}

// addGroupMember adds the specified port id from the port group with the
// specified number of instances. Note that it is an error if the port
// already exists. It also starts a goroutine to service async write requests.
// Note that we still do not need to lock the member as the write operations
// do not change the state of the struct.
func (p *portGroup) addGroupMember(u *fwdpb.AggregatePortAddMemberUpdateDesc) error {
	pid := fwdobject.ID(u.GetPortId().GetObjectId().GetId())
	if _, ok := p.memberMap[pid]; ok {
		return fmt.Errorf("ports: Unable to overwrite an existing port %v in %v", pid, p)
	}
	desc := u.GetSelectActions()
	actions, err := fwdaction.NewActions(desc, p.ctx)
	if err != nil {
		return fmt.Errorf("ports: Unable to create select actions %v for member %v: %v", desc, pid, err)
	}

	m, err := newGroupMember(u.GetPortId(), actions, p.ctx, p, int(u.GetInstanceCount()))
	if err != nil {
		return fmt.Errorf("ports: Unable to create member: %v", err)
	}
	p.memberMap[pid] = m
	p.recompute()
	return nil
}

// updateAlgorithm updates how the port group selects its contituents to
// process a packet.
func (p *portGroup) updateAlgorithm(fields []*fwdpb.PacketFieldId, hash fwdpb.AggregateHashAlgorithm) error {
	// Setup the fields for the packet hash.
	p.fields = make([]fwdpacket.FieldID, 0, len(fields))
	for _, field := range fields {
		p.fields = append(p.fields, fwdpacket.NewFieldID(field))
	}

	// Setup the packet hash function.
	p.hash = hash
	switch p.hash {
	case fwdpb.AggregateHashAlgorithm_AGGREGATE_HASH_ALGORITHM_CRC32:
		p.hashFn = hashCRC32
		p.packetFn = p.selectLink
	case fwdpb.AggregateHashAlgorithm_AGGREGATE_HASH_ALGORITHM_CRC16:
		p.hashFn = hashCRC16
		p.packetFn = p.selectLink
	case fwdpb.AggregateHashAlgorithm_AGGREGATE_HASH_ALGORITHM_FLOOD:
		p.packetFn = p.floodLink
	default:
		return fmt.Errorf("ports: Unable to find hash function %v", hash)
	}
	return nil
}

// updateGroup updates all attributes of the port group.
func (p *portGroup) updateGroup(u *fwdpb.AggregatePortUpdateDesc) error {
	if err := p.updateAlgorithm(u.GetFieldIds(), u.GetHash()); err != nil {
		return err
	}

	// Create a map of actions to be applied by port id.
	actions := make(map[fwdobject.ID]fwdaction.Actions)
	for _, sd := range u.GetSelectActions() {
		desc := sd.GetActions()
		pid := sd.GetPortId()
		a, err := fwdaction.NewActions(desc, p.ctx)
		if err != nil {
			return fmt.Errorf("ports: Unable to create select actions %v for member %v: %v", desc, pid, err)
		}
		actions[fwdobject.ID(pid.GetObjectId().GetId())] = a
	}

	// Create a map of members indexed by the port-id. Track the number of
	// instances of each member.
	memberMap := make(map[fwdobject.ID]*member)
	for _, id := range u.GetPortIds() {
		pid := fwdobject.ID(id.GetObjectId().GetId())
		m, ok := memberMap[pid]
		if ok {
			m.instances++
			continue
		}
		port, err := fwdport.Acquire(id, p.ctx)
		if err != nil {
			return fmt.Errorf("ports: Unable to create port group, missing member: %v", err)
		}
		a, ok := actions[port.ID()]
		if ok {
			delete(actions, port.ID())
		}
		m, err = newGroupMember(id, a, p.ctx, p, 1)
		if err != nil {
			return fmt.Errorf("ports: Unable to create port group member: %v", err)
		}
		memberMap[pid] = m
	}

	// Store the curr map and rebuild the member list.
	curr := p.memberMap
	p.members = nil
	p.memberMap = memberMap
	for _, m := range p.memberMap {
		for index := 0; index < m.instances; index++ {
			p.members = append(p.members, m)
		}
	}

	// Cleanup the old members and any unused actions.
	for _, m := range curr {
		m.Cleanup()
	}
	for _, a := range actions {
		a.Cleanup()
	}
	return nil
}

// Update updates the port group as defined by the update extension.
// Note that only one extension can be valid at a time.
func (p *portGroup) Update(upd *fwdpb.PortUpdateDesc) error {
	switch agg := upd.Port.(type) {
	case *fwdpb.PortUpdateDesc_Aggregate:
		return p.updateGroup(agg.Aggregate)
	case *fwdpb.PortUpdateDesc_AggregateAdd:
		return p.addGroupMember(agg.AggregateAdd)
	case *fwdpb.PortUpdateDesc_AggregateDel:
		return p.removeGroupMember(agg.AggregateDel)
	case *fwdpb.PortUpdateDesc_AggregateAlgo:
		return p.updateAlgorithm(agg.AggregateAlgo.GetFieldIds(), agg.AggregateAlgo.GetHash())
	}
	return errors.New("ports: no extension specified")
}

// Actions returns nil as a port group does not have actions.
func (portGroup) Actions(fwdpb.PortAction) fwdaction.Actions {
	return nil
}

// State implements the port interface. The port group state cannot be
// externally controlled. The group is considered ready to transmit
// at-least one constituent is ready to transmit.
func (p *portGroup) State(op *fwdpb.PortInfo) (fwdpb.PortStateReply, error) {
	for _, m := range p.members {
		if m.Ready() {
			ready := fwdpb.PortStateReply{
				LocalPort: &fwdpb.PortInfo{
					Laser: fwdpb.PortLaserState_PORT_LASER_STATE_ENABLED,
				},
				Link: &fwdpb.LinkStateDesc{
					State: fwdpb.LinkState_LINK_STATE_UP,
					RemotePort: &fwdpb.PortInfo{
						Laser: fwdpb.PortLaserState_PORT_LASER_STATE_ENABLED,
					},
				},
			}
			return ready, nil
		}
	}
	down := fwdpb.PortStateReply{
		LocalPort: &fwdpb.PortInfo{Laser: fwdpb.PortLaserState_PORT_LASER_STATE_DISABLED},
		Link:      &fwdpb.LinkStateDesc{State: fwdpb.LinkState_LINK_STATE_DOWN},
	}
	return down, nil
}

// selectLink selects a port after applying a hash on the packet and writes the packet out on the cable.
func (p *portGroup) selectLink(packet fwdpacket.Packet) (fwdaction.State, error) {
	if p.hashFn == nil {
		return fwdaction.DROP, fmt.Errorf("ports: write to group %v failed, no hash", p)
	}

	var key []byte
	for _, id := range p.fields {
		if f, err := packet.Field(id); err == nil {
			key = append(key, f...)
		}
	}
	index := p.hashFn(key)
	if index >= len(p.members) {
		index = index - len(p.members)*(index/len(p.members))
	}
	m := p.members[index]
	packet.Logf(fwdpacket.LogDebugMessage, "hash selected %v", m.port.ID())
	return fwdaction.CONSUME, m.Write(packet, "Hash")
}

// floodLink floods the packet onto all constituent ports. It does not flood
// the packet onto a port if the port is the input port or if the port is not
// ready. If the write to a constituent fails, the packet is still written
// to all other constituents.
func (p *portGroup) floodLink(packet fwdpacket.Packet) (fwdaction.State, error) {
	pid := fwdobject.InvalidNID
	if field, err := packet.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT, 0)); err == nil {
		pid = fwdobject.NID(binary.BigEndian.Uint64(field))
	}

	frame := packet.Frame()
	for _, m := range p.members {
		if !m.Ready() {
			packet.Logf(fwdpacket.LogDebugMessage, "flood to %v skipped, not ready", m)
			continue
		}

		// Note that this check works even if the packet does not an input port.
		if pid == m.port.NID() {
			packet.Logf(fwdpacket.LogDebugMessage, "flood to %v skipped, do not flood on input", m)
			continue
		}

		// make a copy of the original frame before retransmitting it.
		c := make([]byte, len(frame))
		copy(c, frame)
		pkt, err := fwdpacket.NewNID(packet.StartHeader(), c, pid)
		if err != nil {
			packet.Logf(fwdpacket.LogErrorFrame, "flood to %v failed, err %v", m, err)
			continue
		}

		m.AsyncWrite(pkt)
	}
	return fwdaction.CONSUME, nil
}

// Write transmits the packet out.
func (p *portGroup) Write(packet fwdpacket.Packet) (fwdaction.State, error) {
	if len(p.members) == 0 {
		return fwdaction.DROP, fmt.Errorf("ports: write to group %v failed, no ports", p)
	}
	if p.packetFn == nil {
		return fwdaction.DROP, fmt.Errorf("ports: write to group %v failed, no packet function", p)
	}
	return p.packetFn(packet)
}

// hashCRC32 computes the CRC32 checksum of the key.
func hashCRC32(key []byte) int {
	return int(crc32.ChecksumIEEE(key))
}

// hashCRC16 computes the CRC16 checksum of the key.
func hashCRC16(key []byte) int {
	return int(crc16.ChecksumANSI(key))
}

// A groupBuilder is used to build port groups.
type groupBuilder struct{}

// init registers a builder for port groups.
func init() {
	fwdport.Register(fwdpb.PortType_PORT_TYPE_AGGREGATE_PORT, groupBuilder{})
}

// Build creates a new port group.
func (groupBuilder) Build(pd *fwdpb.PortDesc, ctx *fwdcontext.Context) (fwdport.Port, error) {
	p := portGroup{
		ctx:       ctx,
		memberMap: make(map[fwdobject.ID]*member),
	}

	// Store counters for all ports and actions.
	list := append(fwdport.CounterList, fwdaction.CounterList...)
	if err := p.InitCounters("", "", list...); err != nil {
		return nil, fmt.Errorf("group: Unable to initialize counters, %v", err)
	}
	return &p, nil
}
