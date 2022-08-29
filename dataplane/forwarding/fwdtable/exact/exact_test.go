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

package exact

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/tabletestutil"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdaction/actions"
)

// exactDesc creates a desc for an exact match entry.
func exactDesc(value int, transient bool) *fwdpb.EntryDesc {
	desc := &fwdpb.EntryDesc{}
	exact := &fwdpb.ExactEntryDesc{
		Transient: proto.Bool(transient),
		Fields: []*fwdpb.PacketFieldBytes{
			{
				FieldId: &fwdpb.PacketFieldId{
					Field: &fwdpb.PacketField{
						FieldNum: fwdpb.PacketFieldNum_IP_VERSION.Enum(),
					},
				},
				Bytes: []byte{uint8(value)},
			},
		},
	}
	proto.SetExtension(desc, fwdpb.E_ExactEntryDesc_Extension, exact)
	return desc
}

// exactMatchTable creates an exact match table. A unique name is created for
// the table object using the index.
func exactMatchTable(ctx *fwdcontext.Context, index int) (fwdtable.Table, error) {
	// Exact match table descriptor.
	desc := &fwdpb.TableDesc{
		TableType: fwdpb.TableType_EXACT_TABLE.Enum(),
		Actions:   tabletestutil.ActionDesc(),
		TableId:   fwdtable.MakeID(fwdobject.NewID(fmt.Sprintf("TABLE=%v", index))),
	}

	// Setup the key specification.
	exact := &fwdpb.ExactTableDesc{
		FieldIds: []*fwdpb.PacketFieldId{
			{
				Field: &fwdpb.PacketField{
					FieldNum: fwdpb.PacketFieldNum_IP_VERSION.Enum(),
				},
			},
		},
	}
	proto.SetExtension(desc, fwdpb.E_ExactTableDesc_Extension, exact)
	return fwdtable.New(ctx, desc)
}

// TestExactTable tests various operations on an exact match table.
func TestExactTable(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Register a mock parser.
	const validSize = 10
	parser := mock_fwdpacket.NewMockParser(ctrl)
	parser.EXPECT().MaxSize(gomock.Any()).Return(validSize).AnyTimes()
	parser.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	fwdpacket.Register(parser)

	ctx := fwdcontext.New("test", "fwd")

	table, err := exactMatchTable(ctx, 1)
	if err != nil {
		t.Errorf("Exact match table create failed, err %v.", err)
	}
	if entries := table.Entries(); len(entries) != 0 {
		t.Errorf("Incorrect number of table entries. Got %v, want 0.", len(entries))
	}

	// Operate on *count* entries in the table.
	const count = 10

	// Add entries to the table.
	for index := 0; index < count; index++ {
		if err := table.AddEntry(exactDesc(index, false), tabletestutil.ActionDesc()); err != nil {
			t.Errorf("AddEntry failed, err %v.", err)
		}
	}
	if entries := table.Entries(); len(entries) != count {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	} else {
		t.Logf("List %v.", entries)
	}

	// Update an entry in the table.
	if err := table.AddEntry(exactDesc(1, false), tabletestutil.ActionDesc()); err != nil {
		t.Errorf("AddEntry failed, err %v.", err)
	}
	if entries := table.Entries(); len(entries) != count {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	} else {
		t.Logf("List %v.", entries)
	}

	// Remove an entry from the table.
	if err := table.RemoveEntry(exactDesc(1, false)); err != nil {
		t.Errorf("RemoveEntry failed, err %v.", err)
	}
	if entries := table.Entries(); len(entries) != (count - 1) {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	} else {
		t.Logf("List %v.", entries)
	}

	// Add an entry again.
	if err := table.AddEntry(exactDesc(1, false), tabletestutil.ActionDesc()); err != nil {
		t.Errorf("AddEntry failed, err %v.", err)
	}
	if entries := table.Entries(); len(entries) != count {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	} else {
		t.Logf("List %v.", entries)
	}

	// Remove a non existing entry.
	if err := table.RemoveEntry(exactDesc(count+1, false)); err != nil {
		t.Logf("RemoveEntry failed as expected, err %v.", err)
	} else {
		t.Error("Unexpected entry removed")
	}
	if entries := table.Entries(); len(entries) != count {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	} else {
		t.Logf("List %v.", entries)
	}

	// Clear the table and ensure that it is empty. Attempt to add another
	// entry.
	table.Clear()
	if entries := table.Entries(); len(entries) != 0 {
		t.Errorf("Incorrect number of table entries after clear. Got %v, want 0.", len(entries))
	}
	if err := table.AddEntry(exactDesc(1, false), tabletestutil.ActionDesc()); err != nil {
		t.Errorf("AddEntry failed, err %v.", err)
	}
	if entries := table.Entries(); len(entries) != 1 {
		t.Errorf("Incorrect number of table entries. Got %v, want 1.", len(entries))
	}
}

// TestExactTableStale tests stale monitoring in an exact match table.
func TestExactTableStale(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const fieldSize = 4
	parser := mock_fwdpacket.NewMockParser(ctrl)
	parser.EXPECT().MaxSize(gomock.Any()).Return(fieldSize).AnyTimes()
	parser.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	fwdpacket.Register(parser)
	ctx := fwdcontext.New("test", "fwd")

	const timeout = 10 * time.Second

	// Create an exact match table with stale timer that uses a fake clock.
	// This is done by first creating a table with no configured timer and
	// then creating a stale list that is run manually.
	tableCreate := func(now *time.Time, index int) fwdtable.Table {
		table, err := exactMatchTable(ctx, index)
		if err != nil {
			t.Fatalf("Exact match table create failed, err %v.", err)
		}

		et := table.(*Table)
		et.stale = newStaleList(timeout, func() time.Time {
			return *now
		})
		return table
	}

	// Each test simulates a series of events on the table. After each event
	// the test verifies the number of entries in the table and the time till
	// the next event according to the stale list processing.
	type testEvent struct {
		elapse    time.Duration // time elapsed since previous event
		transient []int         // transient entries to add
		static    []int         // static entries to add
		reset     []int         // entries to reset
		want      int           // number of entries expected
		next      time.Duration // time pending for next expiration
	}
	tests := [][]testEvent{
		// Table with no transient entries.
		[]testEvent{
			{
				elapse: 0 * time.Second,
				static: []int{1, 2, 3},
				want:   3,
				next:   1 * time.Minute,
			},
		},
		// Table containing static entries and unused transient entries.
		[]testEvent{
			{
				elapse:    0 * time.Second,
				static:    []int{1, 2, 3},
				transient: []int{4, 5},
				want:      5,
				next:      10 * time.Second,
			},
			{
				elapse: 4 * time.Second,
				want:   5,
				next:   6 * time.Second,
			},
			{
				elapse: 7 * time.Second,
				want:   3,
				next:   1 * time.Minute,
			},
		},
		// Table containing static entries, used transient entries and
		// unused transient entries.
		[]testEvent{
			{
				elapse:    0 * time.Second,
				static:    []int{1, 2, 3},
				transient: []int{4, 5},
				want:      5,
				next:      10 * time.Second,
			},
			{
				elapse: 4 * time.Second,
				reset:  []int{4},
				want:   5,
				next:   6 * time.Second,
			},
			{
				elapse: 7 * time.Second,
				want:   4,
				next:   3 * time.Second,
			},
			{
				elapse: 4 * time.Second,
				want:   3,
				next:   1 * time.Minute,
			},
		},
		// Table containing static and transient entries. A transient
		// entry is updated by another transient entry.
		[]testEvent{
			{
				elapse:    0 * time.Second,
				static:    []int{1, 2, 3},
				transient: []int{4, 5},
				want:      5,
				next:      10 * time.Second,
			},
			{
				elapse:    4 * time.Second,
				transient: []int{4},
				want:      5,
				next:      6 * time.Second,
			},
			{
				elapse: 7 * time.Second,
				want:   4,
				next:   3 * time.Second,
			},
			{
				elapse: 4 * time.Second,
				want:   3,
				next:   1 * time.Minute,
			},
		},
		// Table containing static and transient entries. A transient
		// entry is updated by a static entry.
		[]testEvent{
			{
				elapse:    0 * time.Second,
				static:    []int{1, 2, 3},
				transient: []int{4, 5},
				want:      5,
				next:      10 * time.Second,
			},
			{
				elapse: 4 * time.Second,
				static: []int{4},
				want:   5,
				next:   6 * time.Second,
			},
			{
				elapse: 7 * time.Second,
				want:   4,
				next:   1 * time.Minute,
			},
			{
				elapse: 4 * time.Second,
				want:   4,
				next:   1 * time.Minute,
			},
		},
	}
	for tid, test := range tests {
		var now time.Time
		table := tableCreate(&now, tid)
		et := table.(*Table)
		for eid, e := range test {
			now = now.Add(e.elapse)
			for _, id := range e.static {
				if err := table.AddEntry(exactDesc(id, false), tabletestutil.ActionDesc()); err != nil {
					t.Errorf("test %d, event %d: add failed for static entry %d: %v", tid, eid, id, err)
				}
			}
			for _, id := range e.transient {
				if err := table.AddEntry(exactDesc(id, true), tabletestutil.ActionDesc()); err != nil {
					t.Errorf("test %d, event %d: add failed for transient entry %d: %v", tid, eid, id, err)
				}
			}
			for _, reset := range e.reset {
				packet := mock_fwdpacket.NewMockPacket(ctrl)
				packet.EXPECT().Field(gomock.Any()).Return([]byte{byte(reset)}, nil).AnyTimes()
				packet.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				table.Process(packet, nil)
			}
			if next := et.stale.process(et); next != e.next {
				t.Errorf("test %d, event %d: Unexpected duration for next event. Got %v, want %v.", tid, eid, next, e.next)
			}
			entries := table.Entries()
			if len(entries) != e.want {
				t.Errorf("test %d, event %d: Incorrect number of table entries. Got %v, want %v.", tid, eid, len(entries), e.want)
			}
		}
	}
}
