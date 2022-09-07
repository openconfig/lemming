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
	"encoding/binary"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdport"
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

// TestTransmit tests the transmit action and builder.
func TestTransmit(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a context with a test port.
	ctx := fwdcontext.New("test", "fwd")

	// Setup the mock port. Note that before insertion, the ID and NID of the
	// port must be set to the invalid values.
	id := fwdobject.NewID("Port1")
	pid := fwdport.MakeID(id)
	port := mock_fwdport.NewMockPort(ctrl)
	port.EXPECT().ID().Return(fwdobject.InvalidID).Times(1)
	port.EXPECT().NID().Return(fwdobject.InvalidNID).Times(1)
	port.EXPECT().Init(id, gomock.Any(), gomock.Any()).Times(1)
	port.EXPECT().Acquire().Return(nil).Times(1)
	port.EXPECT().Release(false /*forceCleanup*/).Return(nil).Times(1)
	port.EXPECT().Increment(gomock.Any(), gomock.Any()).AnyTimes()
	if err := ctx.Objects.Insert(port, pid.ObjectId); err != nil {
		t.Fatalf("Port insert failed, err %v.", err)
	}
	port.EXPECT().NID().Return(fwdobject.NID(1)).AnyTimes()
	port.EXPECT().ID().Return(id).AnyTimes()

	// Create a transmit action using its builder.
	desc := fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_ACTION_TYPE_TRANSMIT,
	}
	transmit := fwdpb.TransmitActionDesc{
		PortId: pid,
	}
	desc.Action = &fwdpb.ActionDesc_Transmit{
		Transmit: &transmit,
	}
	action, err := fwdaction.New(&desc, ctx)
	if err != nil {
		t.Errorf("NewAction failed, desc %v failed, err %v.", desc, err)
	}

	// Verify the action by processing a packet and verifying the counters
	// and results. Note that a transmit action always terminates
	// processing of the packet.
	const length = 10
	verify := func(want fwdaction.State) {
		var base fwdobject.Base
		if err := base.InitCounters("prefix", "desc", fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_TX_ERROR_OCTETS); err != nil {
			t.Fatalf("InitCounters failed, %v", err)
		}

		packet := mock_fwdpacket.NewMockPacket(ctrl)
		packet.EXPECT().Length().Return(length).AnyTimes()

		if want == fwdaction.CONTINUE {
			fid := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT, 0)
			pid := make([]byte, protocol.SizeUint64)
			binary.BigEndian.PutUint64(pid, uint64(port.NID()))
			packet.EXPECT().Update(fid, fwdpacket.OpSet, pid).Times(1)
		}

		next, state := action.Process(packet, &base)
		switch {
		case next != nil:
			t.Errorf("%v processing returned bad actions. Got %v want nil.", action, next)
		case state != want:
			t.Errorf("%v processing returned bad result. Got %v want %v.", action, state, want)
		}
	}

	// Verify a transmit action associated with a port.
	verify(fwdaction.CONTINUE)

	// Verify a transmit action drops the packet after it has been cleaned up.
	c := action.(fwdobject.Composite)
	c.Cleanup()
	verify(fwdaction.DROP)
}
