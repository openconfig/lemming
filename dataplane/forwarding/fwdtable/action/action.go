// Copyright 2023 Google LLC
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

package action

import (
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type entry struct {
	start      int
	numActions int
}

type Table struct {
	fwdobject.Base
	ctx            *fwdcontext.Context
	actions        fwdaction.Actions
	entries        map[string]*entry
	defaultActions fwdaction.Actions
}

// AddEntry adds or updates the actions associated with the specified key.
func (t *Table) AddEntry(ed *fwdpb.EntryDesc, ad []*fwdpb.ActionDesc) error {
	act, ok := ed.Entry.(*fwdpb.EntryDesc_Action)
	if !ok {
		return fmt.Errorf("action: AddEntry failed, missing desc")
	}
	a, err := fwdaction.NewActions(ad, t.ctx)
	if err != nil {
		return err
	}
	if _, ok := t.entries[act.Action.Id]; ok {
		if err := t.RemoveEntry(ed); err != nil {
			return err
		}
	}
	e := &entry{
		numActions: len(ad),
	}
	switch act.Action.InsertType {
	case fwdpb.ActionEntryDesc_TYPE_PREPEND:
		for _, entry := range t.entries {
			entry.start += len(ad)
		}
		t.actions = append(a, t.actions...)
	case fwdpb.ActionEntryDesc_TYPE_APPEND:
		e.start = len(ad)
		t.actions = append(t.actions, a...)
	default:
		return fmt.Errorf("unknown insert type: %v", act.Action.InsertType)
	}
	t.entries[act.Action.GetId()] = e
	return nil
}

// RemoveEntry removes an entry.
func (t *Table) RemoveEntry(ed *fwdpb.EntryDesc) error {
	act, ok := ed.Entry.(*fwdpb.EntryDesc_Action)
	if !ok {
		return fmt.Errorf("action: RemoveEntry failed, missing desc")
	}
	e, ok := t.entries[act.Action.GetId()]
	if !ok {
		return fmt.Errorf("action: RemoveEntry failed, unknown entry")
	}
	t.actions = append(t.actions[0:e.start], t.actions[e.start+e.numActions:]...)
	for _, otherEntry := range t.entries {
		if otherEntry.start > e.start {
			otherEntry.start -= e.numActions
		}
	}
	delete(t.entries, act.Action.GetId())
	return nil
}

// Entries lists all entries in a table.
func (t *Table) Entries() []string {
	return []string{}
}

// Clear removes all entries in the table.
func (t *Table) Clear() {
	t.actions.Cleanup()
	t.entries = make(map[string]*entry)
}

func (t *Table) Process(fwdpacket.Packet, fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if len(t.actions) == 0 {
		return t.defaultActions, fwdaction.CONTINUE
	}
	return t.actions, fwdaction.CONTINUE
}

// String returns the table as a formatted string.
func (t *Table) String() string {
	return fmt.Sprintf("Type=ActionTable;Name=%v,%v", t.ID(), t.BaseInfo())
}

func (t *Table) Cleanup() {
	t.actions.Cleanup()
	t.defaultActions.Cleanup()
}

// A builder builds a action table.
type builder struct{}

// init registers a builder for action tables.
func init() {
	fwdtable.Register(fwdpb.TableType_TABLE_TYPE_ACTION, builder{})
}

// Build builds a new action table.
func (builder) Build(ctx *fwdcontext.Context, td *fwdpb.TableDesc) (fwdtable.Table, error) {
	a, err := fwdaction.NewActions(td.Actions, ctx)
	if err != nil {
		return nil, err
	}
	return &Table{
		entries:        map[string]*entry{},
		actions:        fwdaction.Actions{},
		defaultActions: a,
		ctx:            ctx,
	}, nil
}
