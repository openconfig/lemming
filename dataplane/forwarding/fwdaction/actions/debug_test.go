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
	"testing"

	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// TestDebug tests the debug action and builder.
func TestDebug(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a debug action using its builder.
	desc := fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_DEBUG_ACTION.Enum(),
	}
	proto.SetExtension(&desc, fwdpb.E_DebugActionDesc_Extension, &fwdpb.DebugActionDesc{})
	action, err := fwdaction.New(&desc, nil)
	if err != nil {
		t.Errorf("NewAction failed for desc %v, err %v.", desc, err)
	}

	// Verify the action by processing a packet and verifying the counters
	// and results.
	var counters fwdobject.Base
	if err := counters.InitCounters("prefix", "desc", fwdpb.CounterId_RX_DEBUG_PACKETS, fwdpb.CounterId_RX_DEBUG_OCTETS); err != nil {
		t.Fatalf("InitCounter failed, %v", err)
	}

	// Packet of fixed length.
	const length = 10
	packet := mock_fwdpacket.NewMockPacket(ctrl)
	packet.EXPECT().Length().Return(length)
	packet.EXPECT().Debug(true).Times(1)
	next, state := action.Process(packet, &counters)
	switch {
	case next != nil:
		t.Errorf("%v processing returned bad actions. Got %v want nil.", action, next)
	case state != fwdaction.CONTINUE:
		t.Errorf("%v processing returned bad result. Got %v want CONTINUE.", action, state)
	}
	for _, counter := range counters.Counters() {
		switch counter.ID {
		case fwdpb.CounterId_RX_DEBUG_PACKETS:
			if counter.Value != 1 {
				t.Errorf("Invalid counter %v. Got %v, want 1.", counter, counter.Value)
			}
		case fwdpb.CounterId_RX_DEBUG_OCTETS:
			if counter.Value != length {
				t.Errorf("Invalid counter %v. Got %v, want %v.", counter, counter.Value, length)
			}
		}
	}
}
