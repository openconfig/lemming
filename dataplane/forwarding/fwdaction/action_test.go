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

package fwdaction

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A testAction is an action that tracks its state.
type testAction struct {
	id        int
	cleaned   bool
	allocated bool
}

// Cleanup releases all held references (satisfies interface Composite).
func (action *testAction) Cleanup() {
	action.cleaned = true
}

// Process ensures that testAction satisfies interface Action.
func (action *testAction) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (Actions, State) {
	return nil, DROP
}

// String formats the state of the action as a string.
func (action *testAction) String() string {
	return fmt.Sprintf("Id=%v, cleaned=%v, allocated=%v, ", action.id, action.cleaned, action.allocated)
}

// testBuilder builds test actions using a prebuilt action.
type testBuilder struct {
	action *testAction
}

// Build uses the prebuilt action.
func (builder *testBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (Action, error) {
	builder.action.allocated = true
	return builder.action, nil
}

// newTestBuilder creates a new test builder and registers it.
func newTestBuilder(id int, actionType fwdpb.ActionType) *testBuilder {
	builder := &testBuilder{
		action: &testAction{
			id:        id,
			allocated: false,
			cleaned:   false,
		},
	}
	Register(actionType, builder)
	return builder
}

// unregister unregisters a builder for the specified action type.
func unregister(atype fwdpb.ActionType) {
	delete(builders, atype)
}

// validateCleanup verifies that the actions were cleanedup.
func validateCleanup(t *testing.T, builder *testBuilder) {
	if builder.action.allocated && !builder.action.cleaned {
		t.Errorf("Allocated action not cleaned, action %v.", builder.action)
	}
}

// validateAllocate verifies that the actions were allocated.
func validateAllocate(t *testing.T, builder *testBuilder) {
	if !builder.action.allocated {
		t.Errorf("Action not allocated as expected, action %v.", builder.action)
	}
}

// TestActions tests the creation of various action sets.
func TestActions(t *testing.T) {
	var tests = []struct {
		count   int
		actions []fwdpb.ActionType
	}{
		{0, []fwdpb.ActionType{fwdpb.ActionType_DROP_ACTION, fwdpb.ActionType_TRANSMIT_ACTION}},
		{1, []fwdpb.ActionType{fwdpb.ActionType_DROP_ACTION, fwdpb.ActionType_TRANSMIT_ACTION}},
		{2, []fwdpb.ActionType{fwdpb.ActionType_DROP_ACTION, fwdpb.ActionType_TRANSMIT_ACTION}},
	}

	for id, test := range tests {
		t.Logf("#%d: Test with %d actions from %v.", id, test.count, test.actions)

		// For each action type, unregister an existing builder and create a
		// desc. For the first *count* types, register a new builder.
		index := 0
		var descs []*fwdpb.ActionDesc
		var builders []*testBuilder

		for _, actionType := range test.actions {
			unregister(actionType)
			current := actionType
			descs = append(descs, &fwdpb.ActionDesc{ActionType: &current})
			if index < test.count {
				builders = append(builders, newTestBuilder(index, actionType))
				index++
			}
		}

		// Build the action set.
		list, err := NewActions(descs, nil)

		// If we had a builder per action type, we should be successful.
		if len(test.actions) == test.count {
			if err != nil {
				t.Errorf("#%d: NewActions failed, err %s.", id, err)
			}

			// validate allocation in each builder.
			for _, builder := range builders {
				validateAllocate(t, builder)
			}
			return
		}

		// We expect a failure and builders that are cleaned.
		if err == nil {
			t.Errorf("#%d: Unexpected action set created, list %v.", id, list)
		} else {
			t.Logf("#%d: Expected error %v.", id, err)
		}

		// validate cleanup on each builder.
		for _, builder := range builders {
			validateCleanup(t, builder)
		}
	}
}

// recordAction is an action whose processing can be specified and it appends
// its id to a record on each execution.
type recordAction struct {
	id     int
	record *[]int
	result State
	next   Actions
}

// String returns the state of the object as a formatted string. Note that we
// cannot add the "next" to the string because we construct infinit
// action chains in the test.
func (action *recordAction) String() string {
	return fmt.Sprintf("record=%v result=%v", action.id, action.result)
}

// Process ensures that testAction satisfies the interface Action.
func (action *recordAction) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (Actions, State) {
	*action.record = append(*action.record, action.id)
	return action.next, action.result
}

// TestPacketProcessing tests various sequence of operations on a packet.
func TestPacketProcessing(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var record []int

	// Action used to construct an infinite test.
	infinite := &recordAction{
		id:     0,
		record: &record,
		result: CONTINUE,
	}
	infinite.next = []*ActionAttr{NewActionAttr(infinite, false)}

	var processingTests = []struct {
		count   int
		err     bool
		result  State
		actions Actions
	}{
		// Test a packet with no action.
		{0, false, CONTINUE, nil},

		// Test a packet with one action.
		{1, false, CONTINUE,
			[]*ActionAttr{
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: CONTINUE,
				}, false),
			},
		},

		// Test a packet with multiple actions (last terminating).
		{3, false, DROP,
			[]*ActionAttr{
				NewActionAttr(
					&recordAction{
						id:     0,
						record: &record,
						result: CONTINUE,
					}, false),
				NewActionAttr(&recordAction{
					id:     1,
					record: &record,
					result: CONTINUE,
				}, false),
				NewActionAttr(&recordAction{
					id:     2,
					record: &record,
					result: DROP,
				}, false),
			},
		},

		// Test a packet with multiple actions (last terminating).
		{3, false, OUTPUT,
			[]*ActionAttr{
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: CONTINUE,
				}, false),
				NewActionAttr(&recordAction{
					id:     1,
					record: &record,
					result: CONTINUE,
				}, false),
				NewActionAttr(&recordAction{
					id:     2,
					record: &record,
					result: OUTPUT,
				}, false),
			},
		},

		// Test a packet with multiple actions (middle terminating).
		{2, false, CONSUME,
			[]*ActionAttr{
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: CONTINUE,
				}, false),
				NewActionAttr(&recordAction{
					id:     1,
					record: &record,
					result: CONSUME,
				}, false),
				NewActionAttr(&recordAction{
					id:     2,
					record: &record,
					result: DROP,
				}, false),
			},
		},

		// Test a packet with multiple actions (middle terminating).
		{2, false, OUTPUT,
			[]*ActionAttr{
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: CONTINUE,
				}, false),
				NewActionAttr(&recordAction{
					id:     1,
					record: &record,
					result: OUTPUT,
				}, false),
				NewActionAttr(&recordAction{
					id:     2,
					record: &record,
					result: DROP,
				}, false),
			},
		},

		// Test a packet with multiple actions (complex).
		{7, false, DROP,
			[]*ActionAttr{
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: CONTINUE,
				}, false),
				NewActionAttr(&recordAction{
					id:     1,
					record: &record,
					result: CONTINUE,
				}, false),
				NewActionAttr(&recordAction{
					id:     2,
					record: &record,
					result: CONTINUE,
					next: []*ActionAttr{
						NewActionAttr(&recordAction{
							id:     3,
							record: &record,
							result: CONTINUE,
						}, false),
						NewActionAttr(&recordAction{
							id:     4,
							record: &record,
							result: CONTINUE,
							next: []*ActionAttr{
								NewActionAttr(&recordAction{
									id:     5,
									record: &record,
									result: CONTINUE,
								}, false),
								NewActionAttr(&recordAction{
									id:     6,
									record: &record,
									result: DROP,
								}, false),
								NewActionAttr(&recordAction{
									id:     7,
									record: &record,
									result: DROP,
								}, false),
							},
						}, false),
						NewActionAttr(&recordAction{
							id:     10,
							record: &record,
							result: DROP,
						}, false),
					},
				}, false),
			},
		},

		// Test a packet with multiple actions (complex).
		{9, false, CONTINUE,
			[]*ActionAttr{
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: CONTINUE,
				}, false),
				NewActionAttr(&recordAction{
					id:     1,
					record: &record,
					result: CONTINUE,
				}, false),
				NewActionAttr(&recordAction{
					id:     2,
					record: &record,
					result: CONTINUE,
					next: []*ActionAttr{
						NewActionAttr(&recordAction{
							id:     3,
							record: &record,
							result: CONTINUE,
						}, false),
						NewActionAttr(&recordAction{
							id:     4,
							record: &record,
							result: CONTINUE,
							next: []*ActionAttr{
								NewActionAttr(&recordAction{
									id:     5,
									record: &record,
									result: CONTINUE,
								}, false),
								NewActionAttr(&recordAction{
									id:     6,
									record: &record,
									result: CONTINUE,
								}, false),
								NewActionAttr(&recordAction{
									id:     7,
									record: &record,
									result: CONTINUE,
								}, false),
							},
						}, false),
						NewActionAttr(&recordAction{
							id:     8,
							record: &record,
							result: CONTINUE,
						}, false),
					},
				}, false),
			},
		},

		// Test a packet with infinite actions.
		{50, true, DROP, []*ActionAttr{NewActionAttr(infinite, false)}},

		// Test a packet with one onEvaluate action.
		{0, false, CONTINUE,
			[]*ActionAttr{
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: CONTINUE,
				}, true),
			},
		},

		// Test a packet with one Evaluate action.
		{1, false, CONTINUE,
			[]*ActionAttr{
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: EVALUATE,
				}, false),
			},
		},

		// Test a packet with multiple actions, two onevaluate and no evaluate.
		{1, false, DROP,
			[]*ActionAttr{
				NewActionAttr(
					&recordAction{
						id:     1,
						record: &record,
						result: CONTINUE,
					}, true),
				NewActionAttr(&recordAction{
					id:     2,
					record: &record,
					result: CONTINUE,
				}, true),
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: DROP,
				}, false),
			},
		},

		// Test a packet with multiple actions, two onevaluate and one drop
		// before evaluate.
		{1, false, DROP,
			[]*ActionAttr{
				NewActionAttr(
					&recordAction{
						id:     2,
						record: &record,
						result: CONTINUE,
					}, true),
				NewActionAttr(&recordAction{
					id:     3,
					record: &record,
					result: CONTINUE,
				}, true),
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: DROP,
				}, false),
				NewActionAttr(&recordAction{
					id:     1,
					record: &record,
					result: EVALUATE,
				}, false),
			},
		},

		// Test a packet with multiple actions, two onevaluate and no drop
		// before evaluate.
		{4, false, CONTINUE,
			[]*ActionAttr{
				NewActionAttr(
					&recordAction{
						id:     2,
						record: &record,
						result: CONTINUE,
					}, true),
				NewActionAttr(&recordAction{
					id:     3,
					record: &record,
					result: CONTINUE,
				}, true),
				NewActionAttr(&recordAction{
					id:     0,
					record: &record,
					result: CONTINUE,
				}, false),
				NewActionAttr(&recordAction{
					id:     1,
					record: &record,
					result: EVALUATE,
				}, false),
			},
		},
	}

	// The id of the actions should be in the increasing order of their expected execution.
	for id, test := range processingTests {
		record = nil

		packet := mock_fwdpacket.NewMockPacket(ctrl)
		packet.EXPECT().Length().Return(100).AnyTimes()
		packet.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		packet.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		state, err := ProcessPacket(packet, test.actions, nil)
		switch {
		case err != nil && test.err:
			continue
		case err == nil && test.err:
			t.Errorf("#%d: Unexpected sucessful packet processing.", id)
			continue
		case err != nil && !test.err:
			t.Errorf("#%d: Packet processing failed %v.", id, err)

		}
		if state != test.result {
			t.Errorf("#%d: Incorrect packet state, got state %v want %v.", id, state, test.result)
		}

		t.Logf("#%d: Sequence of execution %v", id, record)
		if len(record) != test.count {
			t.Errorf("#%d: Unexpected number of actions executed, got %d want %d.", id, len(record), test.count)
		}
		for pos, value := range record {
			if pos != value {
				t.Errorf("#%d: Unexpected sequence of actions executed, got %d want %d.", id, value, pos)
			}
		}
	}
}

// TestActionsEquality tests the equality checking between two sets of actions.
func TestActionsEquality(t *testing.T) {
	// Register a builder for the various types of actions used in the test
	// with different ids.
	newTestBuilder(0, fwdpb.ActionType_DROP_ACTION)
	newTestBuilder(1, fwdpb.ActionType_OUTPUT_ACTION)

	// Run tests of various inputs.
	var tests = []struct {
		first     []*fwdpb.ActionDesc
		second    []*fwdpb.ActionDesc
		wantEqual bool // true if we expect them to be equal
	}{
		// Empty lists are trivially true.
		{
			wantEqual: true,
		},
		// Lists are not equal if the first has actions and the second does not.
		{
			first: []*fwdpb.ActionDesc{
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_DROP_ACTION.Enum(),
				},
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_OUTPUT_ACTION.Enum(),
				},
			},
			wantEqual: false,
		},
		// Lists are not equal if the second has actions and the first does not.
		{
			second: []*fwdpb.ActionDesc{
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_DROP_ACTION.Enum(),
				},
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_OUTPUT_ACTION.Enum(),
				},
			},
			wantEqual: false,
		},
		// Lists are not equal if the order is different.
		{
			first: []*fwdpb.ActionDesc{
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_OUTPUT_ACTION.Enum(),
				},
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_DROP_ACTION.Enum(),
				},
			},
			second: []*fwdpb.ActionDesc{
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_DROP_ACTION.Enum(),
				},
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_OUTPUT_ACTION.Enum(),
				},
			},
			wantEqual: false,
		},
		// Lists are not equal if the second is a subset of the first.
		{
			first: []*fwdpb.ActionDesc{
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_OUTPUT_ACTION.Enum(),
				},
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_DROP_ACTION.Enum(),
				},
			},
			second: []*fwdpb.ActionDesc{
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_DROP_ACTION.Enum(),
				},
			},
			wantEqual: false,
		},
		// Lists are not equal if the first is a subset of the second.
		{
			first: []*fwdpb.ActionDesc{
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_DROP_ACTION.Enum(),
				},
			},
			second: []*fwdpb.ActionDesc{
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_DROP_ACTION.Enum(),
				},
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_OUTPUT_ACTION.Enum(),
				},
			},
			wantEqual: false,
		},
		// Lists are equal if the first is the same as the second.
		{
			first: []*fwdpb.ActionDesc{
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_DROP_ACTION.Enum(),
				},
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_OUTPUT_ACTION.Enum(),
				},
			},
			second: []*fwdpb.ActionDesc{
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_DROP_ACTION.Enum(),
				},
				&fwdpb.ActionDesc{
					ActionType: fwdpb.ActionType_OUTPUT_ACTION.Enum(),
				},
			},
			wantEqual: true,
		},
	}

	for id, test := range tests {
		t.Logf("#%d: Test with first: %v, second %v.", id, test.first, test.second)

		// Build the actions.
		first, err := NewActions(test.first, nil)
		if err != nil {
			t.Errorf("#%d: Unable to builds actions from %v, err %v", id, test.first, err)
			continue
		}

		second, err := NewActions(test.second, nil)
		if err != nil {
			t.Errorf("#%d: Unable to builds actions from %v, err %v", id, test.second, err)
			continue
		}

		if gotEqual := first.IsEqual(second); gotEqual != test.wantEqual {
			t.Errorf("#%d: Equality check between %v and %v failed, got %v, want %v", id, test.first, test.second, gotEqual, test.wantEqual)
			continue
		}
	}
}
