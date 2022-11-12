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

// Package prefix implements a Lucius table that performs packet matches using
// the longest prefix match and satisfies the interface fwdtable.Table.
// All entries in the prefix tree are maintained in network order.
package prefix

import (
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/tableutil"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A level is a level within the prefix tree. It represents a sub-string
// within a prefix entry. It is associated with an optional result and
// set of child levels.
type level struct {
	prefix  *key // prefix represented by this level.
	actions fwdaction.Actions
	result  bool
	child   [2]*level // child levels.
	parent  *level    // parent level
}

// A Table is a table which is looked up using a longest prefix match.
// Entries in the table are described by a keyDesc. Packets match an entry if
// the key made from the packet is equal to the key made from the entry.
// The table also has a set of actions which are used to process packets that
// do not match any entry in the table.
type Table struct {
	fwdobject.Base
	desc    tableutil.KeyDesc   // Describes all entries in the table.
	ctx     *fwdcontext.Context // Context for finding objects.
	actions fwdaction.Actions   // Default actions.
	root    *level              // Root for the prefix tree.
}

// newLevel creates a new level in the prefix tree.
func newLevel(prefix *key, parent *level) *level {
	l := &level{
		prefix: prefix.Copy(),
		child:  [2]*level{},
		parent: parent,
	}
	if parent != nil {
		parent.child[l.prefix.Bit(0)] = l
	}
	return l
}

// setResult sets the actions for a level.
func (l *level) setResult(actions fwdaction.Actions) {
	if l.result {
		l.actions.Cleanup()
	}
	l.result = true
	l.actions = actions
}

// clearResult clears the actions for a level.
func (l *level) clearResult() {
	if l.result {
		l.actions.Cleanup()
		l.result = false
		l.actions = nil
	}
}

// getResult returns the actions and true if the level has a result.
func (l *level) getResult() (fwdaction.Actions, bool) {
	return l.actions, l.result
}

// entries formats an entry and its childern as a list of strings.
func (l *level) entries(parent *key) []string {
	prefix := combine(parent, l.prefix)

	var list []string
	if actions, ok := l.getResult(); ok {
		list = append(list, fmt.Sprintf("<prefix=%v>;<actions=%v>;", prefix, actions))
	}
	for _, child := range l.child {
		if child != nil {
			l := child.entries(prefix)
			list = append(list, l...)
		}
	}
	return list
}

// String recurisively prints a table to a string.
func (t *Table) String() string {
	return fmt.Sprintf("Type=PrefixTable;Name=%v;<Desc=%v>;<Default=%v>;%v", t.ID(), t.desc, t.actions, t.BaseInfo())
}

// match matches an sequence of bytes to an entry in the table.
func (t *Table) match(in []byte) (*key, fwdaction.Actions) {
	// Start the iteration with the table's root.
	var actions fwdaction.Actions
	key := newKey(in, Calculate(len(in)))
	curr := t.root
	record := curr.prefix

	// At the start of each iteration, curr is known to be a prefix of key.
	for {
		// store the result of the longest matched prefix.
		if a, ok := curr.getResult(); ok {
			record = curr.prefix
			actions = a
		}

		//  Return the result if the key matched the prefix exactly.
		if key.IsEqual(curr.prefix) {
			break
		}

		// Strip out the prefix and check if the 0 or 1 child is a prefix for
		// the new key.
		key.TrimPrefix(curr.prefix)
		curr = curr.child[key.Bit(0)]
		if curr == nil {
			break
		}
		if !key.HasPrefix(curr.prefix) {
			break
		}
	}
	return record, actions
}

// add adds or updates an entry in the prefix table.
func (t *Table) add(pre *key, actions fwdaction.Actions) {
	// Start the iteration with the table's root.
	curr := t.root
	key := pre

	// At the start of each iteration, curr is known to be a prefix of key.
	for {
		// If the key match exactly, use the current level for the result.
		if key.IsEqual(curr.prefix) {
			break
		}

		// Strip out the prefix and determine the next level based on the 0th bit.
		key.TrimPrefix(curr.prefix)
		child := curr.child[key.Bit(0)]

		// If the child does not exist, add the new level for the result.
		if child == nil {
			curr = newLevel(key, curr)
			break
		}

		// If the complete child is a prefix of the key, continue
		// searching using the child.
		if key.HasPrefix(child.prefix) {
			curr = child
			continue
		}

		// Find the common prefix between the child and the key and add a
		// level for the shared prefix. Since the child is not a prefix of
		// the key, we can strip the common prefix from the child and make the
		// suffix a child of the common prefix. Continue the loop with the
		// common prefix.
		p := prefixKey(child.prefix, key)
		curr = newLevel(p, curr)
		child.prefix.TrimPrefix(p)
		curr.child[child.prefix.Bit(0)] = child
		child.parent = curr
	}
	curr.setResult(actions)
}

// remove removes an entry from the prefix table.
func (t *Table) remove(prefix *key) error {
	// Start the iteration with the table's root.
	curr := t.root
	key := prefix

	// At the start of each iteration, curr is known to be a prefix of key.
	for {
		//  Return the result if the key matched the prefix exactly.
		if key.IsEqual(curr.prefix) {
			break
		}

		// Strip out the prefix and check if the 0 or 1 child is a prefix for
		// the new key.
		key.TrimPrefix(curr.prefix)
		if curr = curr.child[key.Bit(0)]; curr == nil {
			return fmt.Errorf("prefix: Unable to find entry for prefix %v", key)
		}
		if !key.HasPrefix(curr.prefix) {
			return fmt.Errorf("prefix: Unable to find entry for prefix %v from %v", key, curr)
		}
	}

	// Clear the result node.
	curr.clearResult()

	// Try compressing the tree starting from this level towards the root.
	for curr.parent != t.root {
		// If the current level has a result, stop compression.
		if _, ok := curr.getResult(); ok {
			break
		}

		child1 := curr.child[0]
		child2 := curr.child[1]

		// If there are no children, we can just remove the current level.
		if child1 == nil && child2 == nil {
			curr.parent.child[curr.prefix.Bit(0)] = nil
			curr = curr.parent
			continue
		}

		// If there are two children, stop compression.
		if child1 != nil && child2 != nil {
			break
		}

		// If there is exactly one child, we can replace the current level with it.
		child := child1
		if child2 != nil {
			child = child2
		}

		child.prefix = combine(curr.prefix, child.prefix)
		curr.parent.child[curr.prefix.Bit(0)] = child
		child.parent = curr.parent
		curr = curr.parent
	}
	return nil
}

// Clear removes all entries in the table and reallocating an empty root.
func (t *Table) Clear() {
	list := []*level{t.root}
	for len(list) != 0 {
		curr := list[0]
		list = list[1:]
		if curr.child[0] != nil {
			list = append(list, curr.child[0])
			curr.child[0] = nil
		}
		if curr.child[1] != nil {
			list = append(list, curr.child[1])
			curr.child[1] = nil
		}
		curr.clearResult()
	}
	t.root = newLevel(newKey([]byte{}, 0), nil)
}

// Cleanup releases all references held by the table and its entries.
func (t *Table) Cleanup() {
	t.Clear()

	// Release the default actions.
	t.actions.Cleanup()
	t.actions = nil
}

// AddEntry adds or updates the actions associated with the specified key.
func (t *Table) AddEntry(ed *fwdpb.EntryDesc, ad []*fwdpb.ActionDesc) error {
	prefix, ok := ed.Entry.(*fwdpb.EntryDesc_Prefix)
	if !ok {
		return fmt.Errorf("prefix: AddEntry failed, missing desc")
	}
	key, err := newPrefixKey(t.desc, prefix.Prefix.GetFields())
	if err != nil {
		return fmt.Errorf("prefix: AddEntry failed, err %v", err)
	}
	actions, err := fwdaction.NewActions(ad, t.ctx)
	if err != nil {
		return fmt.Errorf("prefix: AddEntry failed, err %v", err)
	}
	t.add(key, actions)
	return nil
}

// RemoveEntry removes an entry.
func (t *Table) RemoveEntry(ed *fwdpb.EntryDesc) error {
	prefix, ok := ed.Entry.(*fwdpb.EntryDesc_Prefix)
	if !ok {
		return fmt.Errorf("prefix: RemoveEntry failed, missing desc")
	}
	key, err := newPrefixKey(t.desc, prefix.Prefix.GetFields())
	if err != nil {
		return fmt.Errorf("prefix: RemoveEntry failed, err %v", err)
	}
	return t.remove(key)
}

// Entries lists all entries in a table.
func (t *Table) Entries() []string {
	return t.root.entries(t.root.prefix)
}

// Process matches the packet to the entries within the table to determine the
// actions to be performed. If the packet does not match any entries, the
// default actions are used. In case of errors, the packet is dropped.
func (t *Table) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	key := t.desc.MakePacketKey(packet)
	if record, actions := t.match(key); actions != nil {
		packet.Logf(fwdpacket.LogDebugMessage, "%v: matched entry %v, actions %v", t.ID(), record, actions)
		return actions, fwdaction.CONTINUE
	}
	packet.Logf(fwdpacket.LogDebugMessage, "%v: default actions %v", t.ID(), t.actions)
	return t.actions, fwdaction.CONTINUE
}

// A builder builds a prefix table.
type builder struct{}

// init registers a builder for prefix tables
func init() {
	fwdtable.Register(fwdpb.TableType_TABLE_TYPE_PREFIX, &builder{})
}

// Build creates a new prefix table.
func (builder) Build(ctx *fwdcontext.Context, td *fwdpb.TableDesc) (fwdtable.Table, error) {
	prefix, ok := td.Table.(*fwdpb.TableDesc_Prefix)
	if !ok {
		return nil, fmt.Errorf("prefix: Build for table failed, missing desc")
	}
	table := &Table{
		desc: tableutil.MakeKeyDesc(prefix.Prefix.GetFieldIds()),
		ctx:  ctx,
		root: newLevel(newKey([]byte{}, 0), nil),
	}
	var err error
	if table.actions, err = fwdaction.NewActions(td.GetActions(), ctx); err != nil {
		return nil, fmt.Errorf("prefix: Build for table failed, err %v", err)
	}
	if err := table.InitCounters("", "", fwdtable.CounterList...); err != nil {
		return nil, fmt.Errorf("prefix: Build for table failed, counter init error, %v", err)
	}
	return table, nil
}
