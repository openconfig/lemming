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

package actions

import (
	"bytes"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/go-logr/logr/testr"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// A recordPort records the last packet written to it.
type recordPort struct {
	fwdobject.Base
	last fwdpacket.Packet
}

// Update is ignored.
func (recordPort) Update(*fwdpb.PortUpdateDesc) error { return nil }

// Write records the last packet written out of the port.
func (r *recordPort) Write(packet fwdpacket.Packet) (fwdaction.State, error) {
	r.last = packet
	return fwdaction.CONSUME, nil
}

func (recordPort) Desc() *fwdpb.PortDesc { return nil }

// String returns an empty string.
func (recordPort) String() string { return "" }

// Type returns the type.
func (recordPort) Type() fwdpb.PortType { return fwdpb.PortType_PORT_TYPE_UNSPECIFIED }

// Actions returns the port actions as nil.
func (recordPort) Actions(fwdpb.PortAction) fwdaction.Actions { return nil }

// State is ignored.
func (recordPort) State(*fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	return &fwdpb.PortStateReply{}, nil
}

// Ensures recordPort implements fwdport.Port interface.
var _ fwdport.Port = &recordPort{}

// makeUpdateDstMacAction returns an action desc to update the mac address.
func makeUpdateDstMacAction(mac []byte) *fwdpb.ActionDesc {
	desc := fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
	}
	update := fwdpb.UpdateActionDesc{
		Type: fwdpb.UpdateType_UPDATE_TYPE_SET,
		FieldId: &fwdpb.PacketFieldId{
			Field: &fwdpb.PacketField{
				FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST,
			},
		},
		Value: mac,
	}
	desc.Action = &fwdpb.ActionDesc_Update{
		Update: &update,
	}
	return &desc
}

// TestMirror tests the mirror action and builder.
func TestMirror(t *testing.T) {
	// Packet data used for generating the test.
	// The original packet.
	orgFrame := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c}

	// The packet with the updated destination mac address.
	updFrame := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c}

	// The mac address used for the update.
	dmac := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xfe}

	// Create a controller and forwarding context
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := fwdcontext.New("test", "fwd")

	tests := []struct {
		hasActions bool // actions applied to the mirrored packet
		hasPort    bool // true if the mirror action has a port
	}{
		// Mirror action never returns a DROP for the orignal packet.
		{
			hasActions: false,
			hasPort:    false,
		},
		// Mirror action with only a port.
		{
			hasActions: false,
			hasPort:    true,
		},
		// Mirror action with only actions and no port.
		{
			hasActions: true,
			hasPort:    false,
		},
		// Mirror action with actions and a port.
		{
			hasActions: true,
			hasPort:    true,
		},
	}

	for idx, test := range tests {
		t.Logf("%d: Running test %+v", idx, test)

		// Create a port that can be used by the test. If the test uses
		// a port, setup the expectations for the various operations and
		// insert it in the object table.
		id := fwdobject.NewID(fmt.Sprintf("Port-%v", idx))
		pid := fwdport.MakeID(id)
		port := &recordPort{}
		if test.hasPort {
			if err := ctx.Objects.Insert(port, pid.ObjectId); err != nil {
				t.Fatalf("%d: Port insert failed, err %v.", idx, err)
			}
		}

		// Create a mirror action using its builder.
		desc := fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_MIRROR,
		}
		mirror := fwdpb.MirrorActionDesc{}

		if test.hasActions {
			ad := makeUpdateDstMacAction(dmac)
			mirror.Actions = []*fwdpb.ActionDesc{ad}
		}
		if test.hasPort {
			mirror.PortId = pid
			mirror.PortAction = fwdpb.PortAction_PORT_ACTION_OUTPUT
		}
		desc.Action = &fwdpb.ActionDesc_Mirror{
			Mirror: &mirror,
		}
		action, err := fwdaction.New(&desc, ctx)
		if err != nil {
			t.Fatalf("%v: NewAction failed, desc %v failed, err %v.", idx, &desc, err)
		}

		// Verify the action by processing a packet and verifying the counters
		// and results.
		var base fwdobject.Base
		if err := base.InitCounters("desc", fwdpb.CounterId_COUNTER_ID_MIRROR_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_MIRROR_ERROR_OCTETS, fwdpb.CounterId_COUNTER_ID_MIRROR_PACKETS, fwdpb.CounterId_COUNTER_ID_MIRROR_OCTETS); err != nil {
			t.Fatalf("%v: InitCounters failed, %v", idx, err)
		}

		// List of fields that are implicitly mirrored
		opFID := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT, 0)
		inFID := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT, 0)

		fields := []fwdpacket.FieldID{
			inFID, opFID,
		}

		// The mock mirrored packet. Note that the expectations on
		// the mirrored packet verifies that the specified updates
		// are performed. Note that the mirrored packet always has
		// an input and output port setup because the original packet
		// always returns an input and output port.
		mirrored := mock_fwdpacket.NewMockPacket(ctrl)
		mirrored.EXPECT().Length().Return(len(orgFrame)).AnyTimes()
		if test.hasActions {
			mirrored.EXPECT().Frame().Return(updFrame).AnyTimes()
		} else {
			mirrored.EXPECT().Frame().Return(orgFrame).AnyTimes()
		}
		mirrored.EXPECT().Attributes().Return(nil).AnyTimes()
		mirrored.EXPECT().Log().Return(testr.New(t)).AnyTimes()
		mirrored.EXPECT().LogMsgs().Return(nil).AnyTimes()
		mirrored.EXPECT().Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0)).Return([]byte{0}, nil).AnyTimes()
		mirrored.EXPECT().Update(opFID, fwdpacket.OpSet, gomock.Any()).Return(nil).AnyTimes()
		mirrored.EXPECT().Update(inFID, fwdpacket.OpSet, gomock.Any()).Return(nil).AnyTimes()
		mirrored.EXPECT().Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST, 0),
			gomock.Any(), gomock.Any()).AnyTimes()

		// The mock original packet. Note that there are no expectations on the
		// original packet that update its fields.
		original := mock_fwdpacket.NewMockPacket(ctrl)
		original.EXPECT().Length().Return(len(orgFrame)).AnyTimes()
		original.EXPECT().Frame().Return(orgFrame).AnyTimes()
		original.EXPECT().Mirror(fields).Return(mirrored, nil).AnyTimes()
		original.EXPECT().Attributes().Return(nil).AnyTimes()
		original.EXPECT().Log().Return(testr.New(t)).AnyTimes()
		original.EXPECT().Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC, 0)).Return([]byte{0}, nil).AnyTimes()
		original.EXPECT().Field(opFID).Return(make([]byte, protocol.SizeUint64), nil).AnyTimes()
		original.EXPECT().Field(inFID).Return(make([]byte, protocol.SizeUint64), nil).AnyTimes()

		// Verify that the action always continues packet processing.
		next, state := action.Process(original, &base)
		switch {
		case next != nil:
			t.Errorf("%v: %v processing returned bad actions. Got %v want nil.", idx, action, next)
		case state != fwdaction.CONTINUE:
			t.Errorf("%v: %v processing returned bad result. Got %v want %v.", idx, action, state, fwdaction.CONTINUE)
		}

		// Verify the packet processed by the port is the mirrored packet.
		if test.hasPort {
			got := port.last.Frame()
			want := orgFrame
			if test.hasActions {
				want = updFrame
			}
			if !bytes.Equal(got, want) {
				t.Errorf("%v: port processed unexpected frame. Got %x, want %x", idx, got, want)
			}
		}
	}
}
