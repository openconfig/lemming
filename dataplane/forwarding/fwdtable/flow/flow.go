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

// Package flow implements a Lucius table that performs packet matches using
// an flow matching and satisfies the interface fwdtable.Table.
package flow

import (
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A level is a set of flows at the specified priority level.
type level struct {
	priority   uint32
	next, prev *level
	flows      *Map
}

// A flowBank is a collection for flow maps associated with a priority.
// The flow maps are maintained in a map indexed by priority and chained in
// the order of decreasing priority.
type flowBank struct {
	levels map[uint32]*level
	head   *level
}

// addLevel adds a new FlowMap to the bank at the specified priority level.
func (b *flowBank) addLevel(fm *Map, priority uint32) {
	curr := &level{
		flows:    fm,
		priority: priority,
	}
	b.levels[priority] = curr
	if b.head == nil {
		b.head = curr
		return
	}

	if priority < b.head.priority {
		curr.next = b.head
		b.head.prev = curr
		b.head = curr
		return
	}

	pos := b.head
	//nolint:revive // Empty block is ok here.
	for ; pos.next != nil && pos.next.priority < priority; pos = pos.next {
	}
	curr.next = pos.next
	curr.prev = pos
	pos.next = curr
	if curr.next != nil {
		curr.next.prev = curr
	}
}

// removeLevel removes the specified priority level
func (b *flowBank) removeLevel(priority uint32) {
	curr := b.levels[priority]
	if b.head == curr {
		b.head = curr.next
	}
	if curr.next != nil {
		curr.next.prev = curr.prev
	}
	if curr.prev != nil {
		curr.prev.next = curr.next
	}
	delete(b.levels, priority)
}

// A FieldDesc describes a set of fields.
type FieldDesc struct {
	fields map[fwdpacket.FieldID]bool // map of fields used to describe flows in the table
	list   *FieldList                 // FieldList precomputed when the fields in the FieldDesc are updated
}

// NewFieldDesc creates a FieldDesc with no fields.
func NewFieldDesc() *FieldDesc {
	return &FieldDesc{
		fields: make(map[fwdpacket.FieldID]bool),
		list:   NewFieldList(nil),
	}
}

// Update updates a FieldDesc with a set of fields. If new fields are added
// to the FieldDesc, update recomputes the field list and returns true.
func (f *FieldDesc) Update(fields []fwdpacket.FieldID) bool {
	updated := false
	for _, id := range fields {
		if _, ok := f.fields[id]; !ok {
			updated = true
			f.fields[id] = true
		}
	}

	if updated {
		f.list = NewFieldList(f.fields)
	}
	return updated
}

// String returns the FieldDesc formatted as a string.
func (f *FieldDesc) String() string {
	return f.list.String()
}

// A Table matches packet flows to actions. The table provisioning is modeled
// as a TCAM; it consists of multiple "banks" each of which contains a
// priority ordered set of flows. When a packet is processed by the table it is
// matched against flows in priority order within each bank.
//
// There are two types of fields in the flow desc; keys and qualifiers.
//   - keys are fields whose values are described as a set of bits. The bits are
//     expressed as the value and mask.
//   - qualifiers are fields whose values are described as one of the values in a
//     precreated forwarding set.
//
// The entries within the flow table are described using FieldDesc, formed by
// a union of all fields used in various flows in the table. The fields in the
// desc are learned as flow entries are added into the table. If the FieldDesc
// does not contain a field in a newly added flow entry, it is expanded and the
// flow table is rebuilt. Note that fields are never removed from the key desc.
// It is assumed that most flows within the table will have
// similar fields.
//
// The flow table optimizes packet matching by computing FieldDesc for keys and
// qualifiers. For each packet, the qualifier and key is computed for all fields
// just once and then these are matched against all flows sequentially. This
// reduces the overhead of query the same field of the packet for each field.
type Table struct {
	fwdobject.Base
	ctx           *fwdcontext.Context // context for finding objects
	actions       fwdaction.Actions   // default actions
	banks         []*flowBank         // map of banks indexed by the bankID
	keyDesc       *FieldDesc          // Describes fields that are matched using a value and mask
	qualifierDesc *FieldDesc          // Describes fields that are matched using set of valid values
}

// Clear removes all entries in the table.
func (t *Table) Clear() {
	for _, b := range t.banks {
		for pos, level := range b.levels {
			level.flows.Cleanup()
			delete(b.levels, pos)
		}
		b.head = nil
	}
}

// Cleanup releases all references held by the table and its entries.
func (t *Table) Cleanup() {
	t.Clear()
	t.actions.Cleanup()
	t.actions = nil
	t.banks = nil
}

// String returns the table as a formatted string.
func (t *Table) String() string {
	return fmt.Sprintf("Type=FlowTable;Name=%v;<Key=%v>;<Qualifier=%v><Default=%v>;%v", t.ID(), t.keyDesc, t.qualifierDesc, t.actions, t.BaseInfo())
}

// AddEntry adds or updates the actions associated with the specified flow
// entry.
func (t *Table) AddEntry(ed *fwdpb.EntryDesc, ad []*fwdpb.ActionDesc) error {
	fl, ok := ed.Entry.(*fwdpb.EntryDesc_Flow)
	if !ok {
		return fmt.Errorf("flow: AddEntry failed, missing desc")
	}
	bankID := fl.Flow.GetBank()
	if int(bankID) >= len(t.banks) {
		for i := len(t.banks) - 1; i < int(bankID); i++ {
			t.banks = append(t.banks, &flowBank{
				levels: make(map[uint32]*level),
			})
		}
	}
	bank := t.banks[bankID]
	desc, err := NewDesc(t.ctx, fl.Flow.GetFields(), fl.Flow.GetQualifiers())
	if err != nil {
		return err
	}

	keys, qualifiers := desc.Fields()

	// If the current key descriptor does not contain all fields of the new
	// flow descriptor, expand the key descriptor and rebuild all the
	// flow maps.
	rebuild := t.keyDesc.Update(keys)
	if rebuild {
		for _, b := range t.banks {
			for _, level := range b.levels {
				level.flows.rebuild(t.keyDesc.list)
			}
		}
	}
	t.qualifierDesc.Update(qualifiers)

	// Add the flow to the flow map.
	priority := fl.Flow.GetPriority()
	actions, err := fwdaction.NewActions(ad, t.ctx)
	if err != nil {
		return fmt.Errorf("flow: AddEntry failed, err %v", err)
	}
	if level, ok := bank.levels[priority]; ok {
		level.flows.Add(desc, actions)
		return nil
	}
	flows := NewMap(t.keyDesc.list, t.qualifierDesc.list, bankID, priority)
	bank.addLevel(flows, priority)
	flows.Add(desc, actions)
	return nil
}

// RemoveEntry removes the specified flow entry.
func (t *Table) RemoveEntry(ed *fwdpb.EntryDesc) error {
	fl, ok := ed.Entry.(*fwdpb.EntryDesc_Flow)
	if !ok {
		return fmt.Errorf("flow: RemoveEntry failed, missing extension")
	}
	bankID := int(fl.Flow.GetBank())
	if bankID >= len(t.banks) {
		return fmt.Errorf("flow: RemoveEntry failed, invalid bankID %v", fl.Flow.GetBank())
	}
	bank := t.banks[fl.Flow.GetBank()]
	priority := fl.Flow.GetPriority()
	level, ok := bank.levels[priority]
	if !ok {
		return fmt.Errorf("flow: RemoveEntry failed, invalid priority %v", priority)
	}
	desc, err := NewDesc(t.ctx, fl.Flow.GetFields(), fl.Flow.GetQualifiers())
	if err != nil {
		return err
	}
	if err := level.flows.Remove(desc); err != nil {
		return err
	}
	if level.flows.count == 0 {
		bank.removeLevel(priority)
	}
	return nil
}

// Entries lists all table entries in each bank in decreasing order of priority.
func (t *Table) Entries() []string {
	var list []string
	for _, b := range t.banks {
		for priority := b.head; priority != nil; priority = priority.next {
			l := priority.flows.Entries()
			list = append(list, l...)
		}
	}
	return list
}

// Process matches the packet to the entries within the table to determine the
// actions to be performed. If the packet does not match any entries, the
// default actions are used. In case of errors, the packet is dropped.
func (t *Table) Process(packet fwdpacket.Packet, _ fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	key := t.keyDesc.list.MakePacketKey(packet)
	qualifier := t.qualifierDesc.list.MakePacketQualifier(packet)
	var actions fwdaction.Actions
	match := false
	for bid, bank := range t.banks {
		m, f, a, p := func(b *flowBank) (bool, *Desc, fwdaction.Actions, uint32) {
			for priority := b.head; priority != nil; priority = priority.next {
				if match, f, a := priority.flows.Match(key, qualifier); match {
					return true, f, a, priority.priority
				}
			}
			return false, nil, nil, 0
		}(bank)
		if m {
			packet.Log().V(3).Info("flow table matched action", "table", t.ID(), "bank", bid, "priority", p, "entry", f, "actions", a)
			actions = append(actions, a...)
			match = true
			break
		}
	}
	if match {
		return actions, fwdaction.CONTINUE
	}
	packet.Log().V(3).Info("flow table default actions", "table", t.ID(), "actions", t.actions)
	return t.actions, fwdaction.CONTINUE
}

// A flowBuilder builds a flow-match table.
type flowBuilder struct{}

// init registers a builder for flow-match flow
func init() {
	fwdtable.Register(fwdpb.TableType_TABLE_TYPE_FLOW, &flowBuilder{})
}

// Build creates a new flow-match table.
func (flowBuilder) Build(ctx *fwdcontext.Context, td *fwdpb.TableDesc) (fwdtable.Table, error) {
	fl, ok := td.Table.(*fwdpb.TableDesc_Flow)
	if !ok {
		return nil, fmt.Errorf("flow: Build for flow table failed, missing desc")
	}

	banks := make([]*flowBank, fl.Flow.GetBankCount())
	for bid := int(fl.Flow.GetBankCount()) - 1; bid >= 0; bid-- {
		banks[bid] = &flowBank{
			levels: make(map[uint32]*level),
		}
	}
	t := &Table{
		banks:         banks,
		ctx:           ctx,
		keyDesc:       NewFieldDesc(),
		qualifierDesc: NewFieldDesc(),
	}
	var err error
	if t.actions, err = fwdaction.NewActions(td.GetActions(), ctx); err != nil {
		return nil, fmt.Errorf("flow: Build for flow table failed, err %v", err)
	}
	if err := t.InitCounters("", fwdtable.CounterList...); err != nil {
		return nil, fmt.Errorf("flow: Build for flow table failed, counter init error, %v", err)
	}
	return t, nil
}
