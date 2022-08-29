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

package fwdflowcounter

import (
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func createFlowCounterCreateRequest(fc, ctx string) *fwdpb.FlowCounterCreateRequest {
	return &fwdpb.FlowCounterCreateRequest{
		ContextId: &fwdpb.ContextId{Id: ctx},
		Id:        &fwdpb.FlowCounterId{ObjectId: &fwdpb.ObjectId{Id: fc}},
	}
}

// TestFlowCountersCreate tests the creation of FlowCounters.
func TestFlowCountersCreate(t *testing.T) {
	ctxstr := "test"
	ctx := fwdcontext.New(ctxstr, "fwd")

	// List of flow counters created in the test.
	creates := []struct {
		name    string // name of the flow counter
		success bool   // whether creation should be successful
	}{
		{
			name:    "fc1",
			success: true,
		},
		{
			name:    "fc2",
			success: true,
		},
		{
			name:    "fc1",
			success: false,
		},
	}

	// Create all flow counters.
	for _, create := range creates {
		fcreq := createFlowCounterCreateRequest(create.name, ctxstr)
		_, err := New(ctx, fcreq)
		switch {
		case err == nil && create.success == false:
			t.Errorf("New(%v) succeeded, expected error", create.name)

		case err != nil && create.success == true:
			t.Errorf("New(%v) failed with error %v, expected success", create.name, err)
		}
	}
}

// TestFlowCountersQuery tests the reading of packets and bytes counts.
func TestFlowCountersQuery(t *testing.T) {
	ctxstr := "test"
	ctx := fwdcontext.New(ctxstr, "fwd")

	// Create some counters for querying first.
	fc1req := createFlowCounterCreateRequest("fc1", ctxstr)
	fc2req := createFlowCounterCreateRequest("fc2", ctxstr)
	fc1, _ := New(ctx, fc1req)
	fc2, _ := New(ctx, fc2req)

	// Now query them for expected values.
	expects := []struct {
		name    *FlowCounter // name of the flow counter
		octets  uint64       // count of octets to expect
		packets uint64       // count of packets to expect
	}{
		{
			name:    fc1,
			octets:  0,
			packets: 0,
		},
		{
			name:    fc2,
			octets:  0,
			packets: 0,
		},
		{
			name:    fc1,
			octets:  0,
			packets: 0,
		},
	}
	for _, expect := range expects {
		fwfc, err := expect.name.Query()
		if (fwfc.Octets) != expect.octets {
			t.Errorf("Query(%v) Fail: Octets found (%v), Expected (%v)",
				fwfc.Id, (fwfc.Octets), expect.octets)
		}
		if (fwfc.Packets) != expect.packets {
			t.Errorf("Query(%v) Fail: Packets found (%v), Expected (%v)",
				fwfc.Id, (fwfc.Packets), expect.packets)
		}
		if err != nil {
			t.Errorf("Query(%v) Fail: Unexpected err (%v)", fwfc.Id, err)
		}
	}
}
