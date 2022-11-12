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

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdflowcounter"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// TestFlowCounter tests the flowcounter action and flowcounterBuilder builder.
func TestFlowCounter(t *testing.T) {
	// Create a controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a context for the test.
	ctx := fwdcontext.New("test", "fwd")

	// Create a flow-counter.
	createReq := &fwdpb.FlowCounterCreateRequest{
		ContextId: &fwdpb.ContextId{Id: "test"},
		Id:        &fwdpb.FlowCounterId{ObjectId: &fwdpb.ObjectId{Id: "fc1"}},
	}
	fcobj, err := fwdflowcounter.New(ctx, createReq)
	if err != nil {
		t.Errorf("FlowCounter creation failed: %v", err)
	}

	// Create a flowcounter action using its builder.
	desc := fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_ACTION_TYPE_FLOW_COUNTER,
	}
	flowcounterdesc := fwdpb.FlowCounterActionDesc{
		CounterId: &fwdpb.FlowCounterId{ObjectId: &fwdpb.ObjectId{Id: "fc1"}},
	}
	desc.Action = &fwdpb.ActionDesc_Flow{
		Flow: &flowcounterdesc,
	}
	action, err := fwdaction.New(&desc, ctx)
	if err != nil {
		t.Errorf("NewAction failed for desc %v, err %v.", &desc, err)
	}

	// Create a packet of fixed length.
	const octets = 10
	const packets = 1
	packet := mock_fwdpacket.NewMockPacket(ctrl)
	packet.EXPECT().Length().Return(octets)

	// Execute the action to increment flow counters.
	next, state := action.Process(packet, nil)

	// Check results of flowcounter action.
	if next != nil {
		t.Errorf("%v processing returned bad actions. Got %v, expected nil", action, next)
	}
	if state != fwdaction.CONTINUE {
		t.Errorf("%v processing returned bad actions. Got %v, expected CONTINUE", action, state)
	}

	// Query the flow counter to get the octet and packet counts.
	fcval, qerr := fcobj.Query()
	if qerr != nil {
		t.Errorf("FlowCounter Query failed: %v", err)
	}

	// Check the counts against expected values.
	if fcval.Octets != octets {
		t.Errorf("FlowCounter octets mismatch: Expected %v, Got %v", octets, fcval.Octets)
	}
	if fcval.Packets != packets {
		t.Errorf("FlowCounter packetss mismatch: Expected %v, Got %v", packets, fcval.Packets)
	}
}
