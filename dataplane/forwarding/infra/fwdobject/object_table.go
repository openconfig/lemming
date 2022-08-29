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

package fwdobject

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A Table is set of visible forwarding objects. Each object is indexed by its
// ID and NID.
type Table struct {
	mu   sync.Mutex
	pool *NIDPool
	nids map[NID]Object
	ids  map[ID]Object
}

// NewTable returns a new empty object table.
func NewTable() *Table {
	return &Table{
		nids: make(map[NID]Object),
		ids:  make(map[ID]Object),
		pool: NewNIDPool(math.MaxUint64),
	}
}

// String returns all the objects in a ';' delimited string format.
func (table *Table) String() string {
	table.mu.Lock()
	defer table.mu.Unlock()

	list := make([]string, 0, 1+len(table.ids))
	for _, object := range table.ids {
		list = append(list, object.String())
	}
	return strings.Join(append(list, ""), ";")
}

// IDs returns a slice of object-ids.
func (table *Table) IDs() []ID {
	table.mu.Lock()
	defer table.mu.Unlock()

	var ids []ID
	for _, object := range table.ids {
		ids = append(ids, object.ID())
	}
	return ids
}

// Insert adds an object to the table and associates it with the specified
// object id. The object is also setup with a cleanup function that releases
// the NID and references held by the object, when the object's refcount drops
// to zero.
func (table *Table) Insert(object Object, oid *fwdpb.ObjectId) error {
	table.mu.Lock()
	defer table.mu.Unlock()

	if oid == nil {
		return errors.New("fwdobject: insert failed, no object id")
	}
	id := NewID(oid.GetId())

	composite, _ := object.(Composite)

	// Ensure that the object has not be initialized.
	if cid := object.ID(); cid != InvalidID {
		return fmt.Errorf("fwdobject: insert failed, object %v has ID %v", object, cid)
	}
	if cid := object.NID(); cid != InvalidNID {
		return fmt.Errorf("fwdobject: insert failed, object %v has NID %v", object, cid)
	}

	// Ensure the specified ID is unique.
	if current, ok := table.ids[id]; ok {
		return fmt.Errorf("fwdobject: insert failed, reinserting object %v over object %v", object, current)
	}

	// Assign a NID for the object.
	var err error
	var nid NID
	if nid, err = table.pool.Acquire(); err != nil {
		if composite != nil {
			composite.Cleanup()
		}
		return fmt.Errorf("fwdobject: insert failed, %s", err)
	}
	cleanup := func() {
		if composite != nil {
			composite.Cleanup()
		}
		delete(table.nids, nid)
		table.pool.Release(nid)
	}
	object.Init(id, nid, cleanup)
	table.ids[id] = object
	table.nids[nid] = object
	return nil
}

// Remove finds the object with the specific id and removes it from the table.
// It also releases the reference taken during the initial insert. Note that
// the object is findable by the NID until its refcount drops to zero.
// If forceCleanup is set, objects will be cleaned up regardless of refcount.
// forceCleanup should only be set during context deletion!
func (table *Table) Remove(oid *fwdpb.ObjectId, forceCleanup bool) error {
	table.mu.Lock()
	defer table.mu.Unlock()

	if oid == nil {
		return errors.New("fwdobject: remove failed, no object id")
	}
	id := NewID(oid.GetId())
	object, ok := table.ids[id]
	if !ok {
		return fmt.Errorf("fwdobject: remove failed, no object at id %v", id)
	}

	delete(table.ids, id)
	return object.Release(forceCleanup)
}

// Acquire finds the object with the specific id and acquires a reference.
func (table *Table) Acquire(oid *fwdpb.ObjectId) (Object, error) {
	table.mu.Lock()
	defer table.mu.Unlock()

	if oid == nil {
		return nil, errors.New("fwdobject: reference failed, no object id")
	}
	id := NewID(oid.GetId())
	if object, ok := table.ids[id]; ok {
		return object, object.Acquire()
	}
	return nil, fmt.Errorf("fwdobject: reference failed, no object at id %v", id)
}

// FindID finds the object with the specific ID without acquiring a reference.
func (table *Table) FindID(oid *fwdpb.ObjectId) (Object, error) {
	table.mu.Lock()
	defer table.mu.Unlock()

	if oid == nil {
		return nil, errors.New("fwdobject: remove failed, no object id")
	}
	id := NewID(oid.GetId())
	if object, ok := table.ids[id]; ok {
		return object, nil
	}
	return nil, fmt.Errorf("fwdobject: find failed, no object at id %v", id)
}

// FindNID finds the object with the specific NID without acquiring a reference.
func (table *Table) FindNID(nid NID) (Object, error) {
	table.mu.Lock()
	defer table.mu.Unlock()

	if object, ok := table.nids[nid]; ok {
		return object, nil
	}
	return nil, fmt.Errorf("fwdobject: find failed, no object at id %v", nid)
}
