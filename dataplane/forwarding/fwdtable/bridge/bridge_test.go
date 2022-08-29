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

package bridge

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport/porttestutil"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdattribute"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdaction/actions"
	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdport/ports"
)

// addEntry adds a table entry associating the mac address to an action that
// transmits on the specified port.
func addEntry(table fwdtable.Table, mac []byte, portID fwdobject.ID) error {
	ad := fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_TRANSMIT_ACTION.Enum(),
	}
	tac := fwdpb.TransmitActionDesc{
		PortId: fwdport.MakeID(portID),
	}
	proto.SetExtension(&ad, fwdpb.E_TransmitActionDesc_Extension, &tac)

	desc := &fwdpb.EntryDesc{}
	exact := &fwdpb.ExactEntryDesc{
		Fields: []*fwdpb.PacketFieldBytes{
			{
				FieldId: &fwdpb.PacketFieldId{
					Field: &fwdpb.PacketField{
						FieldNum: fwdpb.PacketFieldNum_ETHER_MAC_DST.Enum(),
					},
				},
				Bytes: mac,
			},
		},
	}
	proto.SetExtension(desc, fwdpb.E_ExactEntryDesc_Extension, exact)
	return table.AddEntry(desc, []*fwdpb.ActionDesc{&ad})
}

// createBridge create a bridge table. For test purposes, the bridge table
// does not age its entries.
func createBridge(ctx *fwdcontext.Context, id *fwdpb.TableId) (fwdtable.Table, error) {
	desc := &fwdpb.TableDesc{
		TableType: fwdpb.TableType_BRIDGE_TABLE.Enum(),
		TableId:   id,
	}
	proto.SetExtension(desc, fwdpb.E_BridgeTableDesc_Extension, &fwdpb.BridgeTableDesc{})
	return fwdtable.New(ctx, desc)
}

// createLearn creates a learning action for the specified bridge table.
func createLearn(ctx *fwdcontext.Context, id *fwdpb.TableId) (fwdaction.Action, error) {
	desc := &fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_BRIDGE_LEARN_ACTION.Enum(),
	}
	ext := &fwdpb.BridgeLearnActionDesc{
		TableId: id,
	}
	proto.SetExtension(desc, fwdpb.E_BridgeLearnActionDesc_Extension, ext)
	return fwdaction.New(desc, ctx)
}

// network is a hypothetical network connected to the bridge with exactly one
// point-to-point link. Thus only one mac address can be reached from the port.
type network struct {
	name string       // name of the port
	mac  []byte       // mac address reachable on the network
	port fwdport.Port // port connected to the network
}

// packet is a network packet that records all updates to its fields and
// returns the new field value when queried.
type packet struct {
	fields map[fwdpacket.FieldID][]byte
}

// Field returns the bytes associated with a field ID.
func (p *packet) Field(id fwdpacket.FieldID) ([]byte, error) {
	if b, ok := p.fields[id]; ok {
		return b, nil
	}
	return nil, nil
}

// Update updates a field in the packet.
func (p *packet) Update(id fwdpacket.FieldID, _ int, arg []byte) error {
	p.fields[id] = arg
	return nil
}

// Decap removes the outermost header of the specified type.
func (packet) Decap(fwdpb.PacketHeaderId) error { return nil }

// Encap adds an outermost header of the specified type.
func (packet) Encap(fwdpb.PacketHeaderId) error { return nil }

// Reparse reparses the packet with the specified frame.
func (packet) Reparse(fwdpb.PacketHeaderId, []fwdpacket.FieldID, []byte) error { return nil }

// Mirror mirrors the packet
func (packet) Mirror(fields []fwdpacket.FieldID) (fwdpacket.Packet, error) { return nil, nil }

// String formats the packet into a string.
func (packet) String() string { return "" }

// Length returns the length of the packet in bytes.
func (packet) Length() int { return 0 }

// Frame returns the packet as a slice of bytes.
func (packet) Frame() []byte { return nil }

// Debug control debugging for the packet.
func (packet) Debug(bool) {}

// Logf writes a message to the packet's log.
func (packet) Logf(int, string, ...interface{}) {}

// Log returns the contents of the packet's log.
func (packet) Log() []string { return nil }

// Attributes returns the attributes associated with the packet.
func (packet) Attributes() fwdattribute.Set { return nil }

// StartHeader returns the first header of the packet.
func (packet) StartHeader() fwdpb.PacketHeaderId { return fwdpb.PacketHeaderId_ETHERNET }

// sendPacket sends a packet from the source network to the destination network
// through the bridge and returns the port used for transmission.
func sendPacket(ctrl *gomock.Controller, table fwdtable.Table, action fwdaction.Action, networks []network, src, dst int, ctx *fwdcontext.Context) (fwdobject.ID, error) {
	p := &packet{
		fields: make(map[fwdpacket.FieldID][]byte),
	}
	fwdport.SetInputPort(p, networks[src].port)
	p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ETHER_MAC_SRC, 0), fwdpacket.OpSet, networks[src].mac)
	p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ETHER_MAC_DST, 0), fwdpacket.OpSet, networks[dst].mac)

	if _, state := action.Process(p, nil); state != fwdaction.CONTINUE {
		return "", errors.New("bridge: learn failed")
	}
	actions, state := table.Process(p, nil)
	if state != fwdaction.CONTINUE {
		return "", fmt.Errorf("bridge: packet processing stopped. Got %v, want %v", state, fwdaction.CONTINUE)
	}
	fwdaction.ProcessPacket(p, actions, nil)
	port, err := fwdport.OutputPort(p, ctx)
	if err != nil {
		return "", errors.New("bridge: packet was not sent")
	}
	return port.ID(), nil
}

// TestBridge tests the bridge table.
func TestBridge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a context with two ports and two cables.
	ctx := fwdcontext.New("test", "fwd")

	nw := []network{
		{
			name: "p1",
			mac:  []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		},
		{
			name: "p2",
			mac:  []byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16},
		},
		{
			name: "p3",
			mac:  []byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26},
		},
	}

	for pos, n := range nw {
		nw[pos].port = porttestutil.CreateTestPort(t, ctx, fmt.Sprintf("%v-%v", n.name, pos))
	}

	// Setup the parser with fields needed for the bridge.
	parser := mock_fwdpacket.NewMockParser(ctrl)
	parser.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	fid := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ETHER_MAC_SRC, 0)
	parser.EXPECT().MaxSize(fid).Return(6).AnyTimes()

	fid = fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_ETHER_MAC_DST, 0)
	parser.EXPECT().MaxSize(fid).Return(6).AnyTimes()

	fid = fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_INPUT, 0)
	parser.EXPECT().MaxSize(fid).Return(protocol.SizeUint64).AnyTimes()

	fid = fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_OUTPUT, 0)
	parser.EXPECT().MaxSize(fid).Return(protocol.SizeUint64).AnyTimes()
	fwdpacket.Register(parser)

	// packet defines a packet flow from the source network to the
	// destination network.
	type packet struct {
		src, dst int  // index of the source and destination network
		err      bool // true if an error is expected
		learn    bool // true if the packet should trigger a learn
	}

	// Each test sets up a bridge and generates packet flows. The test
	// validates the port selected by the bridge and the number of entries
	// in the bridge.
	tests := []struct {
		static     []int    // prepopulated networks in the bridge
		packets    []packet // packets to transmit
		entryCount int      // expected number of entries in the bridge after the test is run
	}{
		{
			// Transmit packets between two known networks.
			static: []int{0, 1},
			packets: []packet{
				{
					src: 0,
					dst: 1,
				},
				{
					src: 1,
					dst: 0,
				},
			},
			entryCount: 2,
		},
		{
			// Transmit packets between a known network and an unknown network.
			// The first packet learns the unknown network.
			static: []int{0},
			packets: []packet{
				{
					src:   1,
					dst:   0,
					learn: true,
				},
				{
					src: 0,
					dst: 1,
				},
			},
			entryCount: 2,
		},
		{
			// Transmit packets between a known network and an unknown network.
			// The first packet is destined to the unknown network.
			static: []int{0},
			packets: []packet{
				{
					src: 0,
					dst: 1,
					err: true,
				},
				{
					src:   1,
					dst:   0,
					learn: true,
				},
				{
					src: 0,
					dst: 1,
				},
			},
			entryCount: 2,
		},
	}

	for tid, test := range tests {
		bid := fwdtable.MakeID(fwdobject.NewID(fmt.Sprintf("TABLE=%v", tid)))

		// Create a bridge table.
		table, err := createBridge(ctx, bid)
		if err != nil {
			t.Fatalf("%d: Unable to create bridge: %v.", tid, err)
		}

		action, err := createLearn(ctx, bid)
		if err != nil {
			t.Fatalf("%d: Unable to create bridge learn action: %v.", tid, err)
		}

		for _, e := range test.static {
			if err := addEntry(table, nw[e].mac, nw[e].port.ID()); err != nil {
				t.Fatalf("%d: Unable to add network %d to bridge: %v.", tid, e, err)
			}
		}

		// Set the bridge's notification channel.
		bt := table.(*Table)
		bt.notify = make(chan bool)

	next:
		for id, packet := range test.packets {
			got, err := sendPacket(ctrl, table, action, nw, packet.src, packet.dst, ctx)
			switch packet.err {
			case true:
				if err == nil {
					t.Fatalf("%d packet %d: Incorrectly sent packet from %v to %v on port %v.", tid, id, nw[packet.src], nw[packet.dst], got)
				}
				continue next
			case false:
				if err != nil {
					t.Fatalf("%d packet %d: Unable to send packet from %v to %v: %v.", tid, id, nw[packet.src], nw[packet.dst], err)
				}
			}
			want := nw[packet.dst].port.ID()
			if got != want {
				t.Fatalf("%d packet %d: Packet sent on incorrect port. Got %v, want %v.", tid, id, got, want)
			}

			// Drain the learn channel manually.
			if packet.learn {
				select {
				case <-bt.notify:
				case <-time.After(1 * time.Second):
					t.Fatalf("%d packet %d: learn processing timeout.", tid, id)
				}
			}
		}

		entries := table.Entries()
		if test.entryCount != len(entries) {
			t.Fatalf("%d: Bridge has incorrect number of entries. Got %v, want %v.", tid, len(entries), test.entryCount)
		}
	}
}
