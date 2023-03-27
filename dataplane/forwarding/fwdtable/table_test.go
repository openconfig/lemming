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

package fwdtable

import (
	"fmt"
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A testTable is a Table that tracks its state.
type testTable struct {
	fwdobject.Base
	id        int
	cleaned   bool
	allocated bool
}

// Name returns the name of the table.
func (testTable) Name() string { return "" }

// String returns the state of the port as a formatted string.
func (table *testTable) String() string {
	return fmt.Sprintf("ID=%v, cleaned=%v, allocated=%v, ", table.id, table.cleaned, table.allocated)
}

// Cleanup releases all held references (satisfies interface Composite).
func (table *testTable) Cleanup() {
	table.cleaned = true
}

// Process ensures that testAction satisfies interface Action.
func (table *testTable) Process(fwdpacket.Packet, fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	return nil, fwdaction.DROP
}

// Add adds or updates a table entry (satisfies interface Table).
func (table *testTable) AddEntry(*fwdpb.EntryDesc, []*fwdpb.ActionDesc) error {
	return nil
}

// Remove removes a table entry (satisfies interface Table).
func (table *testTable) RemoveEntry(*fwdpb.EntryDesc) error {
	return nil
}

// Entries lists all entries in a table (satisfies interface Table).
func (table *testTable) Entries() []string {
	return nil
}

// Clear removes all entries in the table.
func (table *testTable) Clear() {}

// testBuilder builds test tables using a prebuilt table.
type testBuilder struct {
	table *testTable
}

// Build uses the prebuilt table.
func (builder *testBuilder) Build(*fwdcontext.Context, *fwdpb.TableDesc) (Table, error) {
	builder.table.allocated = true
	return builder.table, nil
}

// unregister unregisters a builder for the specified table type.
func unregister(tableType fwdpb.TableType) {
	delete(builders, tableType)
}

// newTestBuilder creates a new test builder and registers it.
func newTestBuilder(id int, tableType fwdpb.TableType) *testBuilder {
	builder := &testBuilder{
		table: &testTable{
			id:        id,
			allocated: false,
			cleaned:   false,
		},
	}
	Register(tableType, builder)
	return builder
}

// TestTable tests various table operations.
func TestTable(t *testing.T) {
	tableType := fwdpb.TableType_TABLE_TYPE_EXACT
	unregister(tableType)

	ctx := fwdcontext.New("test", "fwd")

	// Create a table, no builder registered.
	table, err := New(ctx, &fwdpb.TableDesc{TableType: tableType})
	if err != nil {
		t.Logf("Got expected error %s.", err)
	} else {
		t.Errorf("Unexpected object created, got %v.", table)
	}

	// Create a table, builder registered.
	builder := newTestBuilder(10, tableType)
	table, err = New(ctx, &fwdpb.TableDesc{
		TableType: tableType,
		TableId:   MakeID(fwdobject.NewID("TestTable")),
	})
	if err != nil {
		t.Errorf("Table create failed, err %v.", err)
	}
	if !builder.table.allocated {
		t.Errorf("Table not allocated, table %v.", builder.table)
	}
	id := string(table.ID())

	// Find the table using an invalid object id.
	invalid := id + "1"
	if table, err = Find(ctx, &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: invalid}}); err != nil {
		t.Logf("Got expected error %v.", err)
	} else {
		t.Errorf("Found unexpected table %v.", table)
	}

	// Find the table using a valid object id.
	if _, err = Find(ctx, &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: id}}); err != nil {
		t.Errorf("Table find failed, err %v.", err)
	}
}
