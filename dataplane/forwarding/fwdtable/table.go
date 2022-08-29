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

// Package fwdtable contains routines and types to manage forwarding tables.
//
// A forwarding table is a mapping of entries to a set of actions. Packets
// are matched to entries within the table using an algorithm specific to
// the table. Lucius has a variety of forwarding tables (implemented by various
// types). All ports are created by provisioning.
//
// This package defines the following mechanisms to manage tables
//  1. It allows different types of tables to register builders during
//     package initialization.
//  2. Provisioning can create tables using the registered builders.
//  3. It defines an interface that can be used to operate on a table.
package fwdtable

import (
	"errors"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// CounterList is a set of counters incremented by tables.
var CounterList = []fwdpb.CounterId{
	fwdpb.CounterId_COUNTER_ID_ERROR_PACKETS,
	fwdpb.CounterId_COUNTER_ID_ERROR_OCTETS,
}

// A Table is a mapping of entries to a set of actions. It is used to match
// entries to packets to determine the actions to be performed on the packet.
//
// Tables are always created by provisioning.
type Table interface {
	fwdobject.Object

	// Process processes the packet using the specified counters.
	// Process returns an updated processing state for the packet.
	// It may also return new Actions that are used to process the packet.
	// Process assumes that counters are always valid (not nil).
	Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State)

	// Add adds or updates a table entry.
	AddEntry(entryDesc *fwdpb.EntryDesc, descs []*fwdpb.ActionDesc) error

	// Remove removes a table entry.
	RemoveEntry(entryDesc *fwdpb.EntryDesc) error

	// Entries lists all entries in a table.
	Entries() []string

	// Clear removes all entries in the table.
	Clear()
}

// A Builder can build a Table of the specified type.
type Builder interface {
	// Build builds a table.
	Build(ctx *fwdcontext.Context, desc *fwdpb.TableDesc) (Table, error)
}

// builders is a map of builders for various types of tables.
var builders = make(map[fwdpb.TableType]Builder)

// Register registers a builder for a table type. Note that builders are
// expected to be registered during package initialization.
func Register(tableType fwdpb.TableType, builder Builder) {
	builders[tableType] = builder
}

// New creates a new table.
func New(ctx *fwdcontext.Context, desc *fwdpb.TableDesc) (Table, error) {
	if desc == nil {
		return nil, errors.New("fwdtable: new failed, missing description")
	}
	builder, ok := builders[desc.GetTableType()]
	if !ok {
		return nil, fmt.Errorf("fwdtable: new failed, no builder for table %s", desc)
	}
	table, err := builder.Build(ctx, desc)
	if err != nil {
		return nil, err
	}
	tid := desc.GetTableId()
	if tid == nil {
		return nil, errors.New("fwdtable: new failed, missing id")
	}
	if err = ctx.Objects.Insert(table, tid.GetObjectId()); err != nil {
		return nil, err
	}
	return table, nil
}

// Find finds a table.
func Find(ctx *fwdcontext.Context, id *fwdpb.TableId) (Table, error) {
	if id == nil {
		return nil, errors.New("fwdtable: find failed, no table specified")
	}
	object, err := ctx.Objects.FindID(id.GetObjectId())
	if err != nil {
		return nil, err
	}
	if table, ok := object.(Table); ok {
		return table, nil
	}
	return nil, fmt.Errorf("fwdtable: find failed, %v is not a table", id)
}

// Acquire acquires a reference to a table.
func Acquire(ctx *fwdcontext.Context, id *fwdpb.TableId) (Table, error) {
	if id == nil {
		return nil, errors.New("fwdtable: acquire failed, no table specified")
	}
	object, err := ctx.Objects.Acquire(id.GetObjectId())
	if err != nil {
		return nil, err
	}
	if table, ok := object.(Table); ok {
		return table, nil
	}

	// Release the object if there was an error.
	_ = object.Release(false /*forceCleanup*/)
	return nil, fmt.Errorf("fwdtable: acquire failed, %v is not a table", id)
}

// Release releases a reference to a table.
func Release(table Table) error {
	if table == nil {
		return errors.New("fwdtable: release failed, no table specified")
	}
	return table.Release(false /*forceCleanup*/)
}

// MakeID makes a TableID corresponding to the specified object ID.
func MakeID(id fwdobject.ID) *fwdpb.TableId {
	return &fwdpb.TableId{
		ObjectId: fwdobject.MakeID(id),
	}
}

// GetID returns the TableID for the given port.
func GetID(table Table) *fwdpb.TableId {
	return MakeID(table.ID())
}
