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
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/go-logr/logr/testr"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport/porttestutil"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// createPortGroup creates a port group using the specified ports. For writing
// packets, the port group uses the specified hash function and a single byte
// packet field to select a port.
func createPortGroup(t *testing.T, ctx *fwdcontext.Context, ports []fwdport.Port, hashFn fwdpb.AggregateHashAlgorithm, index int) fwdport.Port {
	desc := &fwdpb.PortDesc{
		PortType: fwdpb.PortType_PORT_TYPE_AGGREGATE_PORT,
		PortId:   fwdport.MakeID(fwdobject.NewID(fmt.Sprintf("GROUP_PORT=%v", index))),
	}
	port, err := fwdport.New(desc, ctx)
	if err != nil {
		t.Fatalf("Port create failed: %v.", err)
	}
	var update fwdpb.PortUpdateDesc
	var portList []*fwdpb.PortId
	for _, port := range ports {
		portList = append(portList, fwdport.GetID(port))
	}
	group := &fwdpb.AggregatePortUpdateDesc{
		PortIds: portList,
		Hash:    hashFn,
		FieldIds: []*fwdpb.PacketFieldId{{
			Field: &fwdpb.PacketField{
				FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO,
				Instance: 10,
			},
		}},
	}
	update.Port = &fwdpb.PortUpdateDesc_Aggregate{
		Aggregate: group,
	}
	if err := port.Update(&update); err != nil {
		t.Fatalf("Port update failed: %v.", err)
	}
	return port
}

// TestPortGroupWrite tests write operations for a port group.
// It verifies that the same port is selected when the same packet
// is sent repeatedly.
func TestPortGroupWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	names := []string{"p1", "p2", "p3", "p4"}
	ctx := fwdcontext.New("test", "fwd")
	var ports []fwdport.Port
	for _, name := range names {
		ports = append(ports, porttestutil.CreateTestPort(t, ctx, name))
	}

	pg := createPortGroup(t, ctx, ports, fwdpb.AggregateHashAlgorithm_AGGREGATE_HASH_ALGORITHM_CRC32, 0)
	packet := mock_fwdpacket.NewMockPacket(ctrl)
	packet.EXPECT().Length().Return(10).AnyTimes()
	packet.EXPECT().Log().Return(testr.New(t)).AnyTimes()
	packet.EXPECT().Field(gomock.Any()).Return(make([]byte, 8), nil).AnyTimes()
	packet.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	packet.EXPECT().Frame().Return(nil).AnyTimes()

	pg.Write(packet)
	pm := porttestutil.PortMap(ports)
	if len(pm) != 1 {
		t.Fatalf("Port group selected incorrect number of ports. Got %v, want 1.", len(pm))
	}

	for index := 0; index < 10; index++ {
		pg.Write(packet)
		cur := porttestutil.PortMap(ports)
		if len(cur) != 1 {
			t.Fatalf("Port group selected incorrect number of ports. Got %v, want 1.", len(cur))
		}
		for p := range cur {
			if _, ok := pm[p]; !ok {
				t.Fatalf("Port group selected unexpected port. Got %v, want one of %v.", p, pm)
			}
		}
	}
	if obj, ok := pg.(fwdobject.Composite); ok {
		obj.Cleanup()
	}
}

// TestPortGroupHash tests hash operations for a port group.
// For different hash algorithms, its writes packets with varying bytes to the
// port group and ensures that more than one port is selected.
func TestPortGroupHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	names := []string{"p1", "p2", "p3", "p4"}
	ctx := fwdcontext.New("test", "fwd")

	hashFn := []fwdpb.AggregateHashAlgorithm{
		fwdpb.AggregateHashAlgorithm_AGGREGATE_HASH_ALGORITHM_CRC16,
		fwdpb.AggregateHashAlgorithm_AGGREGATE_HASH_ALGORITHM_CRC32,
	}

	for id, fn := range hashFn {
		var ports []fwdport.Port
		for _, name := range names {
			ports = append(ports, porttestutil.CreateTestPort(t, ctx, fmt.Sprintf("%v-%v", name, id)))
		}
		pg := createPortGroup(t, ctx, ports, fn, id)

		// Send multiple packets using all possible 1 byte keys.
		for v := 0; v < 256; v++ {
			packet := mock_fwdpacket.NewMockPacket(ctrl)
			packet.EXPECT().Length().Return(10).AnyTimes()
			packet.EXPECT().Log().Return(testr.New(t)).AnyTimes()
			packet.EXPECT().Field(gomock.Any()).Return([]byte{uint8(v), 0, 0, 0, 0, 0, 0, 0}, nil).AnyTimes()
			packet.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			packet.EXPECT().Frame().Return(nil).AnyTimes()
			pg.Write(packet)
		}
		pm := porttestutil.PortMap(ports)
		if len(pm) < 2 {
			t.Errorf("Port group hashFn %v selected too few ports. Got %v, want at-least 2.", fn, len(pm))
		}
		if c, ok := pg.(fwdobject.Composite); ok {
			c.Cleanup()
		}
		for _, p := range ports {
			if c, ok := p.(fwdobject.Composite); ok {
				c.Cleanup()
			}
		}
	}
}

// validate flood write validates that the number of specified ports wrote the
// packet out. Note that this check is performed by reading the counts multiple
// times since the broadcast is performed with asynchronous writes.
func validateFloodWrite(ports []fwdport.Port, want int, t *testing.T) {
	for attempt := 0; ; attempt++ {
		pm := porttestutil.PortMap(ports)
		if len(pm) == want {
			break
		}
		if attempt == 20 {
			t.Fatalf("Port group selected incorrect number of ports. Got %v, want %v.", len(pm), want)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// TestFloodWrite tests write operations for a flood group. Note that this
// test uses a real ethernet/arp frame instead of a mock to ensure that the
// packet copy succeeds. It also ensures that the packet is not flooded onto
// the input port.
func TestFloodWrite(t *testing.T) {
	names := []string{"p1", "p2", "p3", "p4"}
	ctx := fwdcontext.New("test", "fwd")
	var ports []fwdport.Port
	for _, name := range names {
		ports = append(ports, porttestutil.CreateTestPort(t, ctx, name))
	}

	pg := createPortGroup(t, ctx, ports, fwdpb.AggregateHashAlgorithm_AGGREGATE_HASH_ALGORITHM_FLOOD, 0)

	// Create an ARP packet and set its input port.
	arp := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x08, 0x06, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c}
	packet, err := fwdpacket.NewNID(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, arp, ports[0].NID())
	if err != nil {
		t.Fatalf("Flood failed, err %v.", err)
	}

	state, err := pg.Write(packet)
	if err != nil {
		t.Fatalf("Flood failed, err %v.", err)
	}
	if state != fwdaction.CONSUME {
		t.Fatalf("Flood succeeded with incorrect state. Got %v, want %v.", state, fwdaction.CONSUME)
	}

	// Note that the packet should not be flooded to the input port.
	validateFloodWrite(ports, len(names)-1, t)

	if obj, ok := pg.(fwdobject.Composite); ok {
		obj.Cleanup()
	}
}

// addMember adds a member to a specified port group.
func addMember(t *testing.T, _ *fwdcontext.Context, pg fwdport.Port, name string, count int) {
	var update fwdpb.PortUpdateDesc
	add := &fwdpb.AggregatePortAddMemberUpdateDesc{
		PortId: &fwdpb.PortId{
			ObjectId: &fwdpb.ObjectId{
				Id: name,
			},
		},
		InstanceCount: uint32(count),
	}
	update.Port = &fwdpb.PortUpdateDesc_AggregateAdd{
		AggregateAdd: add,
	}
	if err := pg.Update(&update); err != nil {
		t.Fatalf("Port update failed: %v.", err)
	}
}

// removeMember removes a member from the specified port group.
func removeMember(t *testing.T, _ *fwdcontext.Context, pg fwdport.Port, name string) {
	var update fwdpb.PortUpdateDesc
	remove := &fwdpb.AggregatePortRemoveMemberUpdateDesc{
		PortId: &fwdpb.PortId{
			ObjectId: &fwdpb.ObjectId{
				Id: name,
			},
		},
	}
	update.Port = &fwdpb.PortUpdateDesc_AggregateDel{
		AggregateDel: remove,
	}
	if err := pg.Update(&update); err != nil {
		t.Fatalf("Port update failed: %v.", err)
	}
}

// createFloodGroup creates an empty port group that floods packets to all
// constituents. This exercises the incremental algorithm update for a port
// group.
func createFloodGroup(t *testing.T, ctx *fwdcontext.Context, index int) fwdport.Port {
	desc := &fwdpb.PortDesc{
		PortType: fwdpb.PortType_PORT_TYPE_AGGREGATE_PORT,
		PortId:   fwdport.MakeID(fwdobject.NewID(fmt.Sprintf("GROUP_PORT=%v", index))),
	}
	port, err := fwdport.New(desc, ctx)
	if err != nil {
		t.Fatalf("Port create failed: %v.", err)
	}
	var update fwdpb.PortUpdateDesc
	group := &fwdpb.AggregatePortAlgorithmUpdateDesc{
		Hash: fwdpb.AggregateHashAlgorithm_AGGREGATE_HASH_ALGORITHM_FLOOD,
	}
	update.Port = &fwdpb.PortUpdateDesc_AggregateAlgo{
		AggregateAlgo: group,
	}
	if err := port.Update(&update); err != nil {
		t.Fatalf("Port update failed: %v.", err)
	}
	return port
}

// TestIncremental tests incremental add and remove operations on a port group.
// The operation is verified by flooding the packet and checking the number of
// copies used. Note that this test uses a real ethernet/arp frame.
func TestIncremental(t *testing.T) {
	ctx := fwdcontext.New("test", "fwd")

	// Create a set of ports used for the test.
	names := []string{"p1", "p2", "p3", "p4"}
	var ports []fwdport.Port
	for _, name := range names {
		ports = append(ports, porttestutil.CreateTestPort(t, ctx, name))
	}

	// A test adds the specified set of ports into the group, and then
	// removes the specified set of ports. It then floods the packet
	// into the port group and verifies that the specified number of
	// ports processed the packet.
	type addOp struct {
		name  string
		count int
	}

	type testDesc struct {
		// ports to add into the group
		add []addOp

		// ports to remove from the group
		remove []string

		// number of ports that are expected to process the packet
		want int
	}

	tests := []testDesc{
		// Group with only the incoming port.
		{
			add: []addOp{
				{
					"p1",
					1,
				},
			},
			want: 0,
		},
		// Group with two ports.
		{
			add: []addOp{
				{
					"p1",
					1,
				},
				{
					"p2",
					1,
				},
			},
			want: 1,
		},
		// Group with two ports.
		{
			add: []addOp{
				{
					"p1",
					1,
				},
				{
					"p2",
					2,
				},
			},
			want: 1,
		},
		// Group with three ports.
		{
			add: []addOp{
				{
					"p1",
					1,
				},
				{
					"p2",
					2,
				},
				{
					"p3",
					1,
				},
			},
			want: 2,
		},
		// Group with two ports (3 added and 1 removed).
		{
			add: []addOp{
				{
					"p1",
					1,
				},
				{
					"p2",
					2,
				},
				{
					"p3",
					1,
				},
			},
			remove: []string{"p2"},
			want:   2,
		},
	}

	for index, test := range tests {
		// Create an ARP packet and set its input port.
		arp := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x08, 0x06, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c}
		packet, err := fwdpacket.NewNID(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, arp, ports[0].NID())
		if err != nil {
			t.Fatalf("Flood failed, err %v.", err)
		}

		// Create a port group with no members.
		pg := createFloodGroup(t, ctx, index)
		for _, add := range test.add {
			addMember(t, ctx, pg, add.name, add.count)
		}
		for _, remove := range test.remove {
			removeMember(t, ctx, pg, remove)
		}

		state, err := pg.Write(packet)
		if err != nil {
			t.Fatalf("Flood failed, err %v.", err)
		}
		if state != fwdaction.CONSUME {
			t.Fatalf("Flood succeeded with incorrect state. Got %v, want %v.", state, fwdaction.CONSUME)
		}

		validateFloodWrite(ports, test.want, t)
		if obj, ok := pg.(fwdobject.Composite); ok {
			obj.Cleanup()
		}
	}
}
