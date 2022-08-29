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

// Package fwdobject contains routines and interfaces used to implement various
// forwarding objects.
//
// Lucius implements a variety of forwarding behaviors with the help of
// forwarding objects explictly created by provisioning. These objects are
// implemented as types satisfying the interface Object.
//
// When a forwarding object is created in a context, the client specifies an
// object ID. The ID is a human readable opaque string used by provisioning for
// managing and updating this object or other objects. Additionally Lucius also
// assigns a NID for each object and uses it for its internal operations.
//
// When provisioning removes an object, its ID is freed and the object is removed
// from the Table. After this point, provisioning cannot find the object using the
// ID. When the reference count of the object drops to zero, its NID is released
// to the pool and its entry is removed. In other words, the object cannot be
// found by ID after it is removed, and cannot be found by NID after all
// references drop to zero.
//
// Provisioning can make and break links between objects. To track the
// correctness of these links, each object maintains a reference count.
// The reference count is incremented and decremented by the Aquire and
// Release methods respectively.
//
// Note that when an object is created, the object's reference count is set
// to one. When the reference count drops to zero, the object is deleted.
// In a correctly behaving system, the reference count drops to zero only when
// provisioning explicitly removes an object.
//
// An object may contain references to other objects. An object may also
// contain instances of types that are not directly managed by provisioning
// i.e. they are implicitly created when the object is created. These
// implicit types may in turn be composed of references to other objects.
// To ensure correct cleanup, all types (objects or
// non-objects) must implement the interface Composite.
package fwdobject

import (
	"fmt"
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdattribute"
	"github.com/openconfig/lemming/dataplane/forwarding/util/stats"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Composite is an interface that must be implemented by objects that
// reference other objects.
type Composite interface {
	// Cleanup releases all references held by the struct.
	Cleanup()
}

// Counter is a counter id and its corresponding counter value.
type Counter struct {
	ID    fwdpb.CounterId
	Value uint64
}

// String formats the counter.
func (c *Counter) String() string {
	return fmt.Sprintf("%v = %v", c.ID, c.Value)
}

// Counters is an interface to maintain counters that can be queried by
// provisioning and incremented while processing packets.
type Counters interface {
	// Counters retrieves all available counters.
	Counters() map[fwdpb.CounterId]Counter

	// Increment increments the specified counter.
	Increment(id fwdpb.CounterId, delta uint32)
}

// ID is an opaque string used in the public API to identify a forwarding object
// within a context. The ID is an opaque string and is expected to be human
// readable.
type ID string

// InvalidID represents a non existing ID.
const InvalidID = ID("")

// NewID creates an ID from the specified string.
func NewID(s string) ID {
	return ID(strings.TrimSpace(s))
}

// MakeID makes an ObjectID to the specified ID.
func MakeID(id ID) *fwdpb.ObjectId {
	return &fwdpb.ObjectId{
		Id: string(id),
	}
}

// Object is an interface for visible objects. A visible entity is an object
// that can be referenced by other objects or provisioning.
type Object interface {
	Counters

	// ID returns the object's ID.
	ID() ID

	// NID returns the object's NID.
	NID() NID

	// String formats the state of the object as a string.
	String() string

	// Acquire acquires a reference on the object.
	Acquire() error

	// Release releases a reference on the object.
	// If forceCleanup is set, the object is cleaned up regardless of references.
	Release(forceCleanup bool) error

	// Init initalizes the object's ID, NID, reference count and relevant
	// interfaces. It also sets a cleanup function called when the
	// object's reference count drops to zero
	Init(id ID, nid NID, cleanup func())

	// Attributes returns the set of attributes associated with the object.
	Attributes() fwdattribute.Set
}

// A Base is a partial implementation of Object used to facilitate
// the implementation of a forwarding object.
// A forwarding object can be implemented as follows:
// 1. Embed a Base.
// 2. Implement String using a call to BaseInfo.
// 3. Implement Composite if it contains references to other Composites.
// A Base is not safe for simultaneous use by multiple goroutines.
type Base struct {
	id         ID
	nid        NID
	refCount   uint64
	attributes fwdattribute.Set
	cleanup    func()
	s          *stats.Stats
}

// ID returns b's id.
func (b *Base) ID() ID {
	return b.id
}

// NID returns b's nid.
func (b *Base) NID() NID {
	return b.nid
}

// Init initalizes b's id and reference count.
func (b *Base) Init(id ID, nid NID, cleanup func()) {
	b.id = id
	b.nid = nid
	b.refCount = 1
	b.attributes = fwdattribute.NewSet()
	b.cleanup = cleanup
}

// Attributes returns the set of attributes associated with the object.
func (b *Base) Attributes() fwdattribute.Set {
	return b.attributes
}

// BaseInfo formats the state of the object as a string.
func (b *Base) BaseInfo() (str string) {
	return fmt.Sprintf("ID=%v, NID=%v, refCnt=%v", b.id, b.nid, b.refCount)
}

// Acquire acquires a reference on the object.
func (b *Base) Acquire() error {
	if b.refCount == 0 {
		if object, ok := interface{}(b).(Object); ok {
			return fmt.Errorf("fwdobject: acquire failed, refcount = 0: %s", object)
		}
		return fmt.Errorf("fwdobject: acquire failed, refcount = 0: %s", b.BaseInfo())
	}
	b.refCount++
	return nil
}

// Release releases a reference on the object.
// If forceCleanup is set, object will be cleaned up regardless of refcount.
func (b *Base) Release(forceCleanup bool) error {
	if b.refCount == 0 {
		if object, ok := interface{}(b).(Object); ok {
			return fmt.Errorf("fwdobject: release failed, refcount = 0: %s", object)
		}
		return fmt.Errorf("fwdobject: release failed, refcount = 0: %s", b.BaseInfo())
	}
	b.refCount--
	if forceCleanup || b.refCount == 0 {
		b.cleanup()
		b.cleanup = func() {} // Don't clean up an object twice.
	}
	return nil
}

// InitCounters initializes all the counters that are maintained on the object.
// Note that this must be called before any call to Increment.
func (b *Base) InitCounters(prefix, desc string, ids ...fwdpb.CounterId) error {
	var entries []stats.EntryDesc
	for _, id := range ids {
		// We use CounterId as id of the stat and id.String(), which
		// returns CounterId_COUNTER_ID_name, as the name of the stat.
		entries = append(entries, stats.EntryDesc{int(id), id.String()})
	}
	ss, err := stats.New(prefix, desc, entries...)
	if err != nil {
		return fmt.Errorf("failed to make stats, %v", err)
	}
	b.s = ss
	return nil
}

// Counters returns all available counters.
func (b *Base) Counters() map[fwdpb.CounterId]Counter {
	if b.s == nil {
		return nil
	}
	r := make(map[fwdpb.CounterId]Counter)
	// stats.GetAll() returns the id and value as int and int64. Therefore,
	// we need to convert accordingly.
	for id, c := range b.s.GetAll() {
		r[fwdpb.CounterId(id)] = Counter{ID: fwdpb.CounterId(id), Value: uint64(c)}
	}
	return r
}

// Increment increments the specified counter if it exists.
func (b *Base) Increment(id fwdpb.CounterId, delta uint32) {
	if b.s == nil {
		log.Errorf("fwdobject: counter is not initialized for %s", b.ID())
		return
	}
	if err := b.s.Add(int(id), int64(delta)); err != nil {
		log.Errorf("%v: missing counter-id %s %v", b, id, int(id))
	}
}
