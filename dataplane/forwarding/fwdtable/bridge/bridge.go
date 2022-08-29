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

// Package bridge implements an exact match table that operates on the packet's
// mac address. It contains mac entries that are either provisioned or learned
// from the processed packets. It emulates a small L2 learning bridge.
package bridge

import (
	"encoding/binary"
	"fmt"

	log "github.com/golang/glog"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/exact"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/util/queue"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// learnRequest is a request to learn the specified portID for the mac address.
type learnRequest struct {
	mac     []byte
	portNID []byte
}

// String generates a debug string for a learn request.
func (req *learnRequest) String() string {
	return fmt.Sprintf("mac=%x, portNID=%v", req.mac, req.portNID)
}

// DebugString generates a debug string for a learn request after the port has been
// identified.
func (req *learnRequest) DebugString(port fwdport.Port) string {
	return fmt.Sprintf("mac=%x, portNID=%v, port=%v", req.mac, req.portNID, fwdport.GetID(port))
}

// Table is a learning bridge table that learns source mac and input ports from
// packets and processes packets using their destination mac address.
// The table can contain two types of entries:
// - Static entries that are added via provisioning.
// - Dynamic entries that are learned from the packet's source mac address.
// The dynamic entries are learned as "best-effort" i.e. a packet is not
// dropped if its mac address cannot be learned. Learned entries are timed
// out if they are not used.
//
// When processing packets, the table creates learn requests for the packet's
// source mac and input port and enqueues them to a channel. A goroutine
// monitors the channel and adds the corresponding entries. Before enqueuing
// a learn request, the process function checks if the mac address already
// exists. This prevents learning existing mac addresses continously. When a
// mac address is initially learned, it is possible to have more than one
// learn request for the same mac address. Hence the channel used for learning
// is buffered.
type Table struct {
	*exact.Table                     // exact table containing mac entries
	learn        *queue.Queue        // unbounded queue for learn requests
	ctx          *fwdcontext.Context // context for finding objects
	notify       chan bool           // if not nil, a notification is generated when an entry is learned (test only)
}

// Clear clears the table by deleting all its entries.
func (t *Table) Clear() {
	t.Table.Clear()
}

// Cleanup cleans up the exact match table and stops learning.
func (t *Table) Cleanup() {
	t.learn.Close()
	t.Table.Cleanup()
}

// String returns the table as a formatted string.
func (t *Table) String() string {
	return fmt.Sprintf("Type=BridgeTable;%s;<Queue=%v>;", t.Table.String(), t.learn)
}

// processLearn processes a learn request and adds an action in the exact match
// table to transmit packets with the specified dest. mac onto the specified
// port. Since learnRequest modifies the exact match table, it acquires a
// write lock on the table's context.
func (t *Table) processLearn(v interface{}) {
	req, ok := v.(*learnRequest)
	if !ok {
		log.Errorf("bridge: processLearn failed for %v", req)
		return
	}

	t.ctx.Lock()
	defer t.ctx.Unlock()

	if t.notify != nil {
		defer func() { t.notify <- true }()
	}

	nid := fwdobject.NID(binary.BigEndian.Uint64(req.portNID))
	p, err := t.ctx.Objects.FindNID(nid)
	if err != nil {
		log.Errorf("bridge: processLearn failed for %v: %v.", req, err)
		return
	}
	port, ok := p.(fwdport.Port)
	if !ok {
		log.Errorf("bridge: processLearn failed for %v: object is not a port.", req)
		return
	}
	ad := fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_TRANSMIT_ACTION.Enum(),
	}
	tac := fwdpb.TransmitActionDesc{
		Immediate: proto.Bool(true),
		PortId:    fwdport.GetID(port),
	}
	proto.SetExtension(&ad, fwdpb.E_TransmitActionDesc_Extension, &tac)

	desc := &fwdpb.EntryDesc{}
	exact := &fwdpb.ExactEntryDesc{
		Transient: proto.Bool(true),
		Fields: []*fwdpb.PacketFieldBytes{
			{
				FieldId: &fwdpb.PacketFieldId{
					Field: &fwdpb.PacketField{
						FieldNum: fwdpb.PacketFieldNum_ETHER_MAC_DST.Enum(),
					},
				},
				Bytes: req.mac,
			},
		},
	}
	proto.SetExtension(desc, fwdpb.E_ExactEntryDesc_Extension, exact)

	// It is not an error if an entry cannot be added. This may happen if
	// we try to learn a mac address that has a static entry.
	if err := t.AddEntry(desc, []*fwdpb.ActionDesc{&ad}); err != nil {
		log.Infof("bridge: Skipping learn for %v %v.", req.DebugString(port), err)
	}
}

// Learn learns the source mac and input port of the packet.
//
// Before creating a learn request for the packet, it looks up the key in the
// exact match table. This prevents the learn channel from being filled up
// by learn requests for pre-existing mac addresses. It is assumed that the
// caller already holds the context's read lock i.e. it is called from an
// action or table.
//
// Since the fields learned from the packet occur concurrent to packet
// processing, it is important that the learnRequest does not contain
// references into the packet.
func (t *Table) Learn(packet fwdpacket.Packet) error {
	var err error
	var lr learnRequest

	macField := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ETHER_MAC_SRC, 0)
	if lr.mac, err = packet.Field(macField); err != nil {
		return fmt.Errorf("bridge: Unable to find source mac, %v", err)
	}
	if e := t.Find(lr.mac); e != nil {
		return nil
	}

	portField := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_INPUT, 0)
	if lr.portNID, err = packet.Field(portField); err != nil {
		return fmt.Errorf("bridge: Unable to find input port, %v", err)
	}
	return t.learn.Write(&lr)
}

// A builder builds a bridge table.
type builder struct{}

// init registers a builder for bridge tables.
func init() {
	fwdtable.Register(fwdpb.TableType_BRIDGE_TABLE, builder{})
}

// Build creates a new bridge table that consists of an exact match table that
// processes packets using the destination mac address. It also starts a
// goroutine to monitor and learn from the table's learn channel. The goroutine
// can safely acquire a write lock on the table's context before learning the
// entry.
func (builder) Build(ctx *fwdcontext.Context, td *fwdpb.TableDesc) (fwdtable.Table, error) {
	if !proto.HasExtension(td, fwdpb.E_BridgeTableDesc_Extension) {
		return nil, fmt.Errorf("bridge: Build for bridge table failed, missing extension %s", fwdpb.E_BridgeTableDesc_Extension.Name)
	}
	md := proto.GetExtension(td, fwdpb.E_BridgeTableDesc_Extension).(*fwdpb.BridgeTableDesc)

	desc := &fwdpb.TableDesc{
		TableType: fwdpb.TableType_EXACT_TABLE.Enum(),
		Actions:   td.GetActions(),
	}

	ed := &fwdpb.ExactTableDesc{
		TransientTimeout: proto.Uint32(md.GetTransientTimeout()),
		FieldIds: []*fwdpb.PacketFieldId{
			{
				Field: &fwdpb.PacketField{
					FieldNum: fwdpb.PacketFieldNum_ETHER_MAC_DST.Enum(),
				},
			},
		},
	}
	proto.SetExtension(desc, fwdpb.E_ExactTableDesc_Extension, ed)
	table, err := exact.New(ctx, desc)
	if err != nil {
		return nil, fmt.Errorf("bridge: Build for bridge table failed: %v", err)
	}

	t := &Table{
		Table: table,
		ctx:   ctx,
	}
	if t.learn, err = queue.NewUnbounded("learn"); err != nil {
		return nil, err
	}
	t.learn.Run()
	go func() {
		for {
			v, ok := <-t.learn.Receive()
			if !ok {
				return
			}
			t.processLearn(v)
		}
	}()
	return t, nil
}
