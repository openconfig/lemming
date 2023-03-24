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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A testAction is a test action that has an id and a reference to a map.
// When the action processes a packet, it increments the count associated
// with the id.
type testAction struct {
	id       int         // id of this action
	countMap map[int]int // reference to a map of counters indexed by id
}

// Processes the packet by incrementing the counter associated with the id.
func (t *testAction) Process(fwdpacket.Packet, fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if c, ok := t.countMap[t.id]; ok {
		t.countMap[t.id] = c + 1
	} else {
		t.countMap[t.id] = 1
	}
	return nil, fwdaction.CONTINUE
}

// String returns the testAction as a string.
func (t *testAction) String() string {
	counter := 0
	if c, ok := t.countMap[t.id]; ok {
		counter = c
	}
	return fmt.Sprintf("id %v counter %v", t.id, counter)
}

// A testBuilder builds a test action. The builder has a map used to count the number
// of times an action with a specific id is called.
// Before an action is built, the builder
type testBuilder struct {
	countMap map[int]int
}

// newBuilder creates and registers a new builder for the test actions.
func newBuilder() *testBuilder {
	tb := testBuilder{
		countMap: make(map[int]int),
	}
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_TEST, &tb)
	return &tb
}

// Build creates a new test action.
func (t *testBuilder) Build(desc *fwdpb.ActionDesc, _ *fwdcontext.Context) (fwdaction.Action, error) {
	ta, ok := desc.Action.(*fwdpb.ActionDesc_Test)
	if !ok {
		return nil, fmt.Errorf("actions: Build for test action failed, missing extension")
	}
	return &testAction{
		id:       int(ta.Test.GetInt1()),
		countMap: t.countMap,
	}, nil
}

// genActionDesc creates a desc for a test action with the specified id.
func genActionDesc(id int) *fwdpb.ActionDesc {
	desc := &fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_ACTION_TYPE_TEST,
	}
	desc.Action = &fwdpb.ActionDesc_Test{
		Test: &fwdpb.TestActionDesc{
			Int1: uint32(id),
		},
	}
	return desc
}

// TestSelectActionList tests the select action list action and builder.
func TestSelectActionList(t *testing.T) {
	// Create a controller and forwarding context.
	ctx := fwdcontext.New("test", "fwd")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Maintain a map of "secondary" action-id indexed by a "primary"
	// action-id. The action lists are created to always execute a primary
	// and secondary action, and the association is used to validate that
	// the entire list is executed.
	actionMap := map[int]int{
		1: 2,
		3: 4,
		5: 6,
	}

	// Create a set of action lists. Each action list has two actions,
	// one with the "primary" and the other with the "secondary" id.
	var actionListDesc []*fwdpb.ActionList
	for p, s := range actionMap {
		var k fwdpb.ActionList
		k.Actions = append(k.Actions,
			genActionDesc(p))
		k.Actions = append(k.Actions,
			genActionDesc(s))
		actionListDesc = append(actionListDesc, &k)
	}

	// Packet fields used to compute the hash.
	fields := []*fwdpb.PacketFieldId{{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO,
		},
	}, {
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP,
		},
	}}

	// List of each type of select function.
	hashFn := []fwdpb.SelectActionListActionDesc_SelectAlgorithm{
		fwdpb.SelectActionListActionDesc_SELECT_ALGORITHM_CRC16,
		fwdpb.SelectActionListActionDesc_SELECT_ALGORITHM_CRC32,
		fwdpb.SelectActionListActionDesc_SELECT_ALGORITHM_RANDOM,
	}

	// For each type of select function, create a select action with the
	// specified actions lists.
	for _, hash := range hashFn {
		desc := fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_SELECT_ACTION_LIST,
		}
		ext := fwdpb.SelectActionListActionDesc{
			SelectAlgorithm: hash,
			FieldIds:        fields,
			ActionLists:     actionListDesc,
		}
		desc.Action = &fwdpb.ActionDesc_Select{
			Select: &ext,
		}

		b := newBuilder()
		actions, err := fwdaction.NewActions([]*fwdpb.ActionDesc{&desc}, ctx)
		if err != nil {
			t.Errorf("NewAction failed for desc %v, err %v.", &desc, err)
		}

		// Send multiple packets using all possible 1 byte keys. This ensures that a hash
		// function would select at least two distinct action lists. Use the map to verify
		// that the action lists are fully executed.
		for v := 0; v < 256; v++ {
			packet := mock_fwdpacket.NewMockPacket(ctrl)
			packet.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			packet.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			packet.EXPECT().Field(gomock.Any()).Return([]byte{uint8(v), 0, 0, 0, 0, 0, 0, 0}, nil).AnyTimes()

			s, err := fwdaction.ProcessPacket(packet, actions, nil)
			if err != nil {
				t.Errorf("Action %v failed to process packet, err %v.", actions, err)
			}
			if s != fwdaction.CONTINUE {
				t.Errorf("Action %v failed to process packet, got state %v, want state %v.", actions, s, fwdaction.CONTINUE)
			}
		}

		if len(b.countMap) < 2 {
			t.Errorf("Action %v failed to select multiple lists, map is %+v.", actions, b.countMap)
		}
		for p, s := range actionMap {
			pcount := b.countMap[p]
			scount := b.countMap[s]
			if pcount != scount {
				t.Errorf("Did not select unique* %+v %v, got %v, want %v", b.countMap, actions, scount, pcount)
			}
		}
	}
}
