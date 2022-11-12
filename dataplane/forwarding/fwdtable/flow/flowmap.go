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

package flow

import (
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
)

// A Entry associates a flow to a set of actions.
type Entry struct {
	desc               *Desc             // description about the flow
	actions            fwdaction.Actions // actions associated with the flow
	key                *EntryKey         // key for the flow precomputed using the flow map's key field list
	qualifier          *EntryQualifier   // qualifier for the flow precomputed using the flow map's qualifier field list
	hashNext, hashPrev *Entry            // flow entries are maintained in a hash table for management
	seqNext, seqPrev   *Entry            // flow entries are maintained in a list for matching sequentially
	bankID             uint32            // bank containing the flow entry (used for display)
	priority           uint32            // priority of the flow entry (used for display)
}

// String returns the entry as a formatted string.
func (e *Entry) String() string {
	return fmt.Sprintf("<flow=%v>;<key=%v>;<qualifier=%v>;<actions=%v>;<bank=%v priority=%v>", e.desc, e.key, e.qualifier, e.actions, e.bankID, e.priority)
}

// A Map maps flows to actions. It supports two methods for finding actions.
// - Lookup: Actions are found by looking up the flow in the hash table.
// - Match:  Actions are found by sequentially matching flows to a specified key.
// A flow is said to match a key if it is equal to the key masked with the flow.
//
// All flows in the map have an equal length.
type Map struct {
	hash             map[uint32]*Entry // entries stored in a hash table for lookup
	keyFields        *FieldList        // fields used for the entry key
	qualifierFields  *FieldList        // fields used for the entry qualifier
	seqHead, seqTail *Entry            // list of flows in the order of addition
	count            int               // number of entries in the flow map
	bankID           uint32            // bank containing the flow map (used for display)
	priority         uint32            // priority of the flow map (used for display)
}

// NewMap creates a new flow map using the specified bank, priority and
// fields for the keys and qualifiers.
func NewMap(keyFields, qualifierFields *FieldList, bankID, priority uint32) *Map {
	if qualifierFields == nil {
		qualifierFields = NewFieldList(nil)
	}
	if keyFields == nil {
		keyFields = NewFieldList(nil)
	}
	return &Map{
		keyFields:       keyFields,
		qualifierFields: qualifierFields,
		hash:            make(map[uint32]*Entry),
		bankID:          bankID,
		priority:        priority,
	}
}

// Cleanup removes all entries in the keyMap and releases the actions.
func (m *Map) Cleanup() {
	for entry := m.seqHead; entry != nil; entry = entry.seqNext {
		entry.qualifier.Cleanup()
		entry.actions.Cleanup()
	}
	m.hash = nil
}

// Find finds the flow and returns the entry if it is found.
func (m *Map) Find(fd *Desc) *Entry {
	for entry := m.hash[fd.Hash()]; entry != nil; entry = entry.hashNext {
		if fd.Equal(entry.desc) {
			return entry
		}
	}
	return nil
}

// Lookup returns the actions whose flow is equal to the specified flow.
func (m *Map) Lookup(fd *Desc) fwdaction.Actions {
	if entry := m.Find(fd); entry != nil {
		return entry.actions
	}
	return nil
}

// Match returns the actions whose flow matches the specified packet.
func (m *Map) Match(key PacketKey, qualifier PacketQualifier) (bool, *Desc, fwdaction.Actions) {
	for entry := m.seqHead; entry != nil; entry = entry.seqNext {
		if entry.key.Match(key) && entry.qualifier.Match(qualifier) {
			return true, entry.desc, entry.actions
		}
	}
	return false, nil, nil
}

// Add associates a set of actions with the specified flow.
func (m *Map) Add(fd *Desc, actions fwdaction.Actions) error {
	b := fd.Hash()

	// update an existing flow.
	for entry := m.hash[b]; entry != nil; entry = entry.hashNext {
		if fd.Equal(entry.desc) {
			a := entry.actions
			entry.actions = actions
			a.Cleanup()
			return nil
		}
	}
	qualifier, err := m.qualifierFields.MakeEntryQualifier(fd)
	if err != nil {
		return err
	}
	entry := &Entry{
		desc:      fd,
		key:       m.keyFields.MakeEntryKey(fd),
		qualifier: qualifier,
		actions:   actions,
		bankID:    m.bankID,
		priority:  m.priority,
	}

	// add the flow to the hash table.
	if head := m.hash[b]; head != nil {
		entry.hashNext = head
		head.hashPrev = entry
	}
	m.hash[b] = entry

	// add the flow to the sequential list.
	if m.seqHead != nil {
		m.seqHead.seqPrev = entry
	}
	entry.seqNext = m.seqHead
	m.seqHead = entry
	if m.seqTail == nil {
		m.seqTail = entry
	}
	m.count++
	return nil
}

// Remove removes the specified flow from the .
func (m *Map) Remove(fd *Desc) error {
	entry := m.Find(fd)
	if entry == nil {
		return fmt.Errorf("FlowMap: remove failed, cannot find %v", fd)
	}

	// remove the flow from the hash table.
	if entry.hashNext != nil {
		entry.hashNext.hashPrev = entry.hashPrev
	}
	if entry.hashPrev != nil {
		entry.hashPrev.hashNext = entry.hashNext
	} else {
		b := fd.Hash()
		if entry.hashNext == nil {
			delete(m.hash, b)
		} else {
			m.hash[b] = entry.hashNext
		}
	}

	// remove the flow from the sequential list.
	if entry.seqPrev != nil {
		entry.seqPrev.seqNext = entry.seqNext
	}
	if entry.seqNext != nil {
		entry.seqNext.seqPrev = entry.seqPrev
	}
	if m.seqHead == entry {
		m.seqHead = entry.seqNext
	}
	if m.seqTail == entry {
		m.seqTail = entry.seqPrev
	}
	entry.qualifier.Cleanup()
	entry.actions.Cleanup()
	m.count--
	return nil
}

// rebuild rebuilds the flow map using the new table descriptor. It recomputes
// the flow keys of each flow and rebuilds the hash table. The sequence of
// flows used for matching keys is unaltered.
func (m *Map) rebuild(kd *FieldList) {
	m.keyFields = kd

	for entry := m.seqTail; entry != nil; entry = entry.seqPrev {
		entry.key = m.keyFields.MakeEntryKey(entry.desc)
	}
}

// Entries lists all entries in the FlowMap.
func (m *Map) Entries() []string {
	var list []string
	for entry := m.seqHead; entry != nil; entry = entry.seqNext {
		list = append(list, entry.String())
	}
	return list
}
