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
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// TestLookup tests the lookup action and builder.
func TestLookup(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a context with a test table.
	ctx := fwdcontext.New("test", "fwd")

	// Default results for the table.
	next := fwdaction.Actions([]*fwdaction.ActionAttr{fwdaction.NewActionAttr(&drop{}, false)})
	result := fwdaction.CONTINUE

	// Setup the mock port. Note that before insertion, the ID and NID of the
	// port must be set to the invalid values.
	id := fwdobject.NewID("Table1")
	tid := fwdtable.MakeID(id)
	table := mock_fwdtable.NewMockTable(ctrl)
	table.EXPECT().Process(gomock.Any(), gomock.Any()).Return(next, result).Times(1)
	table.EXPECT().ID().Return(fwdobject.InvalidID).Times(1)
	table.EXPECT().NID().Return(fwdobject.InvalidNID).Times(1)
	table.EXPECT().Init(id, gomock.Any(), gomock.Any()).Times(1)
	table.EXPECT().Acquire().Return(nil).Times(1)
	table.EXPECT().Release(false /*forceCleanup*/).Return(nil).Times(1)
	if err := ctx.Objects.Insert(table, tid.ObjectId); err != nil {
		t.Fatalf("Port insert failed, err %v.", err)
	}
	table.EXPECT().ID().Return(id).AnyTimes()
	table.EXPECT().NID().Return(fwdobject.NID(1)).AnyTimes()

	// Create a lookup action using its builder.
	desc := fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
	}
	lookup := fwdpb.LookupActionDesc{
		TableId: tid,
	}
	desc.Action = &fwdpb.ActionDesc_Lookup{
		Lookup: &lookup,
	}
	action, err := fwdaction.New(&desc, ctx)
	if err != nil {
		t.Errorf("NewAction failed, desc %v failed, err %v.", &desc, err)
	}

	// Verify the action by processing a packet and verifying the counters
	// and results. Note that a lookup action always returns the table's result.
	const length = 10
	verify := func(errorPackets, errorBytes uint64, n fwdaction.Actions, s fwdaction.State) {
		var base fwdobject.Base
		if err := base.InitCounters("prefix", "desc", fwdpb.CounterId_COUNTER_ID_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_ERROR_OCTETS); err != nil {
			t.Fatalf("InitCounters failed, %v", err)
		}

		packet := mock_fwdpacket.NewMockPacket(ctrl)
		packet.EXPECT().Length().Return(length).AnyTimes()

		next, state := action.Process(packet, &base)
		switch {
		case n == nil && next != nil:
			t.Errorf("%v processing returned bad actions. Got %v want %v.", action, next, n)

		case n != nil && len(n) != len(next) && n[0] != next[0]:
			t.Errorf("%v processing returned bad actions. Got %v want %v.", action, next, n)
		case state != s:
			t.Errorf("%v processing returned bad result. Got %v want %v.", action, state, s)
		}
		counters := base.Counters()
		if counter, ok := counters[fwdpb.CounterId_COUNTER_ID_ERROR_PACKETS]; !ok || counter.Value != errorPackets {
			t.Errorf("Invalid counter %v on drop. Got %v, want %v.", counter, counter.Value, errorPackets)
		}
		if counter, ok := counters[fwdpb.CounterId_COUNTER_ID_ERROR_OCTETS]; !ok || counter.Value != errorBytes {
			t.Errorf("Invalid counter %v on drop.  Got %v, want %v.", counter, counter.Value, errorBytes)
		}
	}

	// Verify lookup processing when the action has a valid table.
	verify(0, 0, next, result)

	// Verify a lookup action drops the packet after it has been cleaned up.
	c := action.(fwdobject.Composite)
	c.Cleanup()
	verify(1, length, nil, fwdaction.DROP)
}
