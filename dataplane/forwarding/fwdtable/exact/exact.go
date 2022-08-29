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

// Package exact implements a Lucius table that performs packet matches using
// an exact match match and satisfies the interface fwdtable.Table.
package exact

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/tableutil"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// An Entry is an entry in an exact-match table. It maps a key to Actions.
// An entry is a part of a hash table collision list and an optional
// stale entry list.
type Entry struct {
	key                  tableutil.Key     // key of the entry
	actions              fwdaction.Actions // actions associated with the key
	hashNext, hashPrev   *Entry            // links within the hash table
	staleTime            time.Time         // time when the entry will be considered stale
	staleNext, stalePrev *Entry            // links within the stale list
	transient            bool              // indicates if the entry is transient
}

// String returns the entry as a formatted string.
func (e *Entry) String() string {
	return fmt.Sprintf("<key=%x>;<actions=%v>;<timeout=%v>;", e.key, e.actions, e.staleTime)
}

// isStale returns true if the entry is now stale.
func (e *Entry) isStale(now time.Time) bool {
	return e != nil && e.staleTime.Before(now)
}

// A staleList is a list of entries that will be removed when they become stale.
// An entry is stale if it is not used for a specified amount of time. All
// entries in the list have the same timeout for stale detection.
type staleList struct {
	head, tail *Entry           // the head and tail of the stale list
	timeout    time.Duration    // duration after which unused entries are considered stale
	stop       chan bool        // channel used to stop the goroutine monitoring the stale list
	now        func() time.Time // function used to get current time
}

// newStaleList creates a new stale list.
func newStaleList(timeout time.Duration, now func() time.Time) *staleList {
	return &staleList{
		now:     now,
		stop:    make(chan bool),
		timeout: timeout,
	}
}

// add adds an entry to the stale list and sets its timeout.
func (l *staleList) add(e *Entry) {
	e.staleTime = l.now().Add(l.timeout)
	e.stalePrev = l.tail
	e.staleNext = nil
	if l.head == nil {
		l.head = e
	} else {
		l.tail.staleNext = e
	}
	l.tail = e
}

// remove removes the entry from the stale list.
func (l *staleList) remove(e *Entry) {
	if e.staleNext != nil {
		e.staleNext.stalePrev = e.stalePrev
	}
	if e.stalePrev != nil {
		e.stalePrev.staleNext = e.staleNext
	}
	if l.tail == e {
		l.tail = e.stalePrev
	}
	if l.head == e {
		l.head = e.staleNext
	}
}

// use resets the entry's position in the stale list.
func (l *staleList) use(e *Entry) {
	l.remove(e)
	l.add(e)
}

// process scans the stale list and removes stale entries. It also returns
// the time till the next stale check. If there are no entries being
// monitored, the next stale check is requested after 1 minute.
func (l *staleList) process(t *Table) time.Duration {
	now := l.now()
	for l.head.isStale(now) {
		e := l.head
		l.remove(e)
		t.remove(e)
	}
	sleep := 1 * time.Minute
	if l.head != nil {
		sleep = l.head.staleTime.Sub(l.now())
	}
	return sleep
}

// A Table is a table that matches packets to entries using exact-match.
// Entries in the table are described by a keyDesc. Packets match an entry if
// the key made from the packet is equal to the key made from the entry.
// The table also has a set of actions which are used to process packets that
// do not match any entry in the table.
//
// The forwarding infrastructure calls the AddEntry, RemoveEntry and Entries
// methods while holding a write lock on the context, and calls the Process
// function while holding a read lock on the context. The stale list monitor
// goroutine calls the Remove method while holding the write lock on the
// context. Since the stale list can be manipulated simultaneously by the
// monitor and packet processing goroutines, the table has a mutex to protect
// stale list operations.
type Table struct {
	fwdobject.Base
	desc    tableutil.KeyDesc   // describes all entries in the table
	entries map[uint32]*Entry   // entries stored in a hash table
	ctx     *fwdcontext.Context // context for finding objects
	actions fwdaction.Actions   // default actions

	staleMu sync.Mutex // mutex to protect the staleList
	stale   *staleList // list of entries that are monitored for stale detection
}

// Clear removes all entries in the table by walking all entries in the table and deleting them.
func (t *Table) Clear() {
	for pos, head := range t.entries {
		for entry := head; entry != nil; entry = entry.hashNext {
			entry.actions.Cleanup()
		}
		delete(t.entries, pos)
	}
}

// Cleanup releases all references held by the table and its entries.
func (t *Table) Cleanup() {
	t.Clear()

	// Release the default actions and stale list.
	t.actions.Cleanup()
	t.actions = nil
	if t.stale != nil {
		close(t.stale.stop)
		t.stale = nil
	}
}

// String returns the table as a formatted string.
func (t *Table) String() string {
	return fmt.Sprintf("Type=ExactTable;Name=%v;<Desc=%v>;<Default=%v>;%v", t.ID(), t.desc, t.actions, t.BaseInfo())
}

// bucket derives a hash bucket using FNV.
func (*Table) bucket(key tableutil.Key) uint32 {
	hash := fnv.New32()
	hash.Write(key)
	return hash.Sum32()
}

// Find looks up a key within the table and returns the entry if it is found.
func (t *Table) Find(key tableutil.Key) *Entry {
	for entry := t.entries[t.bucket(key)]; entry != nil; entry = entry.hashNext {
		if bytes.Equal(entry.key, key) {
			return entry
		}
	}
	return nil
}

// insert inserts an entry into the table.
// It assumes that the key does not exist.
func (t *Table) insert(key tableutil.Key, actions fwdaction.Actions) *Entry {
	bucket := t.bucket(key)
	entry := &Entry{
		key:     key,
		actions: actions,
	}
	if head := t.entries[bucket]; head != nil {
		entry.hashNext = head
		head.hashPrev = entry
	}
	t.entries[bucket] = entry
	return entry
}

// remove removes the specified key from the table.
// It returns an error if the key is not found.
func (t *Table) remove(entry *Entry) {
	if entry.hashNext != nil {
		entry.hashNext.hashPrev = entry.hashPrev
	}
	if entry.hashPrev != nil {
		entry.hashPrev.hashNext = entry.hashNext
	} else {
		bucket := t.bucket(entry.key)
		if entry.hashNext == nil {
			delete(t.entries, bucket)
		} else {
			t.entries[bucket] = entry.hashNext
		}
	}
	entry.actions.Cleanup()
}

// AddEntry adds or updates the actions associated with the specified key.
func (t *Table) AddEntry(ed *fwdpb.EntryDesc, ad []*fwdpb.ActionDesc) error {
	if !proto.HasExtension(ed, fwdpb.E_ExactEntryDesc_Extension) {
		return fmt.Errorf("exact: AddEntry failed, missing extension %s", fwdpb.E_ExactEntryDesc_Extension.Name)
	}
	desc := proto.GetExtension(ed, fwdpb.E_ExactEntryDesc_Extension).(*fwdpb.ExactEntryDesc)
	key, err := newExactKey(t.desc, desc.GetFields())
	if err != nil {
		return fmt.Errorf("exact: AddEntry failed, err %v", err)
	}
	actions, err := fwdaction.NewActions(ad, t.ctx)
	if err != nil {
		return fmt.Errorf("exact: AddEntry failed, err %v", err)
	}
	transient := desc.GetTransient()
	entry := t.Find(key)
	if entry != nil {
		if transient && !entry.transient {
			defer actions.Cleanup()

			// If the actions are identical, it is safe to not return an error.
			if entry.actions.IsEqual(actions) {
				return nil
			}
			return fmt.Errorf("exact: AddEntry failed, cannot override static entry %v", entry.String())
		}
		entry.actions.Cleanup()
		entry.actions = actions
	} else {
		entry = t.insert(key, actions)
	}
	if t.stale != nil {
		t.staleMu.Lock()
		defer t.staleMu.Unlock()

		if entry.transient {
			t.stale.remove(entry)
		}
		if transient {
			entry.transient = transient
			t.stale.add(entry)
		}
	}
	return nil
}

// RemoveEntry removes an entry.
func (t *Table) RemoveEntry(ed *fwdpb.EntryDesc) error {
	if !proto.HasExtension(ed, fwdpb.E_ExactEntryDesc_Extension) {
		return fmt.Errorf("exact: RemoveEntry failed, missing extension %s", fwdpb.E_ExactEntryDesc_Extension.Name)
	}
	desc := proto.GetExtension(ed, fwdpb.E_ExactEntryDesc_Extension).(*fwdpb.ExactEntryDesc)
	key, err := newExactKey(t.desc, desc.GetFields())
	if err != nil {
		return fmt.Errorf("exact: RemoveEntry failed, err %v", err)
	}
	entry := t.Find(key)
	if entry == nil {
		return fmt.Errorf("exact: RemoveEntry failed, cannot find key %v", key)
	}
	if t.stale != nil && entry.transient {
		t.staleMu.Lock()
		t.stale.remove(entry)
		t.staleMu.Unlock()
	}
	t.remove(entry)
	return nil
}

// Entries lists all entries in a table. Note that the order of entries is
// non-deterministic.
func (t *Table) Entries() []string {
	var list []string
	for _, head := range t.entries {
		for entry := head; entry != nil; entry = entry.hashNext {
			list = append(list, entry.String())
		}
	}
	return list
}

// staleMonitor periodically processes the stale list. Since the stale list
// processing can remove entries, it done while holding a write lock
// on the table context. This prevents races with other provisioning and
// can be safely done from the monitor goroutine.
func (t *Table) staleMonitor() {
	go func() {
		for {
			t.ctx.Lock()
			t.staleMu.Lock()
			sleep := t.stale.process(t)
			t.staleMu.Unlock()
			t.ctx.Unlock()
			select {
			case <-t.stale.stop:
				return
			case <-time.After(sleep):
			}
		}
	}()
}

// Process matches the packet to the entries within the table to determine the
// actions to be performed. If the packet does not match any entries, the
// default actions are used. In case of errors, the packet is dropped.
func (t *Table) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	key := t.desc.MakePacketKey(packet)
	if entry := t.Find(key); entry != nil {
		if t.stale != nil && entry.transient {
			t.staleMu.Lock()
			t.stale.use(entry)
			t.staleMu.Unlock()
		}
		packet.Logf(fwdpacket.LogDebugMessage, "%v: matched %v", t.ID(), entry)
		return entry.actions, fwdaction.CONTINUE
	}
	packet.Logf(fwdpacket.LogDebugMessage, "%v: default actions %v", t.ID(), t.actions)
	return t.actions, fwdaction.CONTINUE
}

// An builder builds an exact-match table.
type builder struct{}

// init registers a builder for exact-match tables
func init() {
	fwdtable.Register(fwdpb.TableType_EXACT_TABLE, builder{})
}

// New creates a new exact-match table. The table also has a stale list if a
// timeout is specified for transient entries.
func New(ctx *fwdcontext.Context, td *fwdpb.TableDesc) (*Table, error) {
	if !proto.HasExtension(td, fwdpb.E_ExactTableDesc_Extension) {
		return nil, fmt.Errorf("exact: Build for exact table failed, missing extension %s", fwdpb.E_ExactTableDesc_Extension.Name)
	}
	ed := proto.GetExtension(td, fwdpb.E_ExactTableDesc_Extension).(*fwdpb.ExactTableDesc)
	t := &Table{
		desc:    tableutil.MakeKeyDesc(ed.GetFieldIds()),
		entries: make(map[uint32]*Entry),
		ctx:     ctx,
	}
	if ed.GetTransientTimeout() != 0 {
		t.stale = newStaleList(time.Duration(ed.GetTransientTimeout())*time.Second, time.Now)
		t.staleMonitor()
	}
	var err error
	if t.actions, err = fwdaction.NewActions(td.GetActions(), ctx); err != nil {
		return nil, fmt.Errorf("exact: Build for extact table failed, err %v", err)
	}
	if err := t.InitCounters("", "", fwdtable.CounterList...); err != nil {
		return nil, fmt.Errorf("exact: Build for extact table failed, counter init error, %v", err)
	}
	return t, nil
}

// Build creates a new exact-match table.
func (builder) Build(ctx *fwdcontext.Context, td *fwdpb.TableDesc) (fwdtable.Table, error) {
	return New(ctx, td)
}
