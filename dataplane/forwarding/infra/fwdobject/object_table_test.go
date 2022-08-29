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
	"fmt"
	"testing"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// testObject is an object used for testing.
type testObject struct {
	Base
	cleaned bool
}

// Cleanup sets a flag indicating that the object was cleaned.
func (t *testObject) Cleanup() {
	t.cleaned = true
}

// String returns the state of the object as a formatted string.
func (t *testObject) String() string {
	return fmt.Sprintf("%s, object cleaned=%v", t.BaseInfo(), t.cleaned)
}

// TestEmpty tests an empty object table.
func TestEmpty(t *testing.T) {
	table := NewTable()

	t.Logf("Table contents %s", table)
	if ids := table.IDs(); len(ids) != 0 {
		t.Fatalf("Unexpected number of entries. len(%v) = %v, want 0.", ids, len(ids))
	}

	id := "10"
	oid := &fwdpb.ObjectId{Id: &id}
	if object, err := table.FindID(oid); err == nil {
		t.Errorf("Found unexpected object %s at id %v.", object, id)
	} else {
		t.Logf("Expected error %s.", err)
	}

	if object, err := table.Acquire(oid); err == nil {
		t.Errorf("Acquired unexpected object %s at id %v.", object, id)
	} else {
		t.Logf("Expected error %s.", err)
	}

	if err := table.Remove(oid, false /*forceCleanup*/); err == nil {
		t.Errorf("Removed unexpected object at id %v.", id)
	} else {
		t.Logf("Expected error %s.", err)
	}

	var object testObject
	if err := object.Release(false /*forceCleanup*/); err == nil {
		t.Errorf("Released unexpected object at id %v.", id)
	} else {
		t.Logf("Expected error %s.", err)
	}
}

// TestObjects tests a table with multiple objects.
func TestObjects(t *testing.T) {
	ids := []string{"id1", "id2", "id3", "id4", "id5"}
	var testobjects []*testObject
	var objects []Object

	table := NewTable()
	testIndex := 1
	testForceIndex := 2

	// Add all objects into the table.
	for _, id := range ids {
		var object testObject
		testobjects = append(testobjects, &object)
		objects = append(objects, &object)
		if err := table.Insert(&object, &fwdpb.ObjectId{Id: &id}); err != nil {
			t.Fatalf("Insert failed, err %s.", err)
		}
	}

	// Check the content of the tables
	t.Logf("Table contents %s", table)
	if ids := table.IDs(); len(ids) != len(objects) {
		t.Errorf("Unexpected number of entries. len(%v) = %v, want %v.", ids, len(ids), len(objects))
	}
	for _, object := range objects {
		id := string(object.ID())
		if _, err := table.FindID(&fwdpb.ObjectId{Id: &id}); err != nil {
			t.Errorf("Find failed, err %s.", err)
		}
	}

	// Re-insert an object and expect it to fail.
	object := objects[testIndex]
	id := string(object.ID())
	oid := &fwdpb.ObjectId{Id: &id}

	if err := table.Insert(object, &fwdpb.ObjectId{Id: &id}); err == nil {
		t.Errorf("Inserted unexpected object %v.", object)
	} else {
		t.Logf("Expected error %s.", err)
	}

	// Take an additional reference on the object.
	if _, err := table.Acquire(oid); err != nil {
		t.Errorf("Acquire failed, %s.", err)
	}
	t.Logf("Table contents after acquire %s.", table)

	// Remove the object and ensure it is remove from the table but not cleaned.
	if err := table.Remove(oid, false /*forceCleanup*/); err != nil {
		t.Errorf("Remove failed, err %s.", err)
	}
	if _, err := table.FindID(oid); err == nil {
		t.Errorf("Found unexpected object %s at id %v.", object, object.ID())
	} else {
		t.Logf("Expected error, err %s.", err)
	}
	if testobjects[testIndex].cleaned {
		t.Errorf("Object is unexpectedly cleaned, object %s.", object)
	}
	t.Logf("Table contents after remove %s.", table)

	// Deference the object and ensure it is successful and the object is cleaned.
	if err := object.Release(false /*forceCleanup*/); err != nil {
		t.Errorf("Release failed, err %s.", err)
	}
	if !testobjects[testIndex].cleaned {
		t.Errorf("Object is not cleaned, object %s.", object)
	}
	t.Logf("Table contents after release %s.", table)

	// Remove a different object and force it to be cleaned up.
	object = objects[testForceIndex]
	id = string(object.ID())
	oid = &fwdpb.ObjectId{Id: &id}

	if err := table.Remove(oid, true /*forceCleanup*/); err != nil {
		t.Errorf("Remove failed, err %s.", err)
	}
	if _, err := table.FindID(oid); err == nil {
		t.Errorf("Found unexpected object %s at id %v.", object, object.ID())
	} else {
		t.Logf("Expected error, err %s.", err)
	}
	if !testobjects[testForceIndex].cleaned {
		t.Errorf("Object is unexpectedly not cleaned, object %s.", object)
	}
	t.Logf("Table contents after remove %s.", table)
}

// TestCounters tests the counters on an object.
func TestCounters(t *testing.T) {
	validate := func(counters Counters, length int) error {
		valid := 0
		list := counters.Counters()
		if len(list) != length {
			return fmt.Errorf("Unexpected number of entries, len(%v) = %d, want %d", list, len(list), length)
		}
		for _, counter := range list {
			if uint64(counter.ID) == counter.Value {
				valid++
			}
		}
		if valid != length {
			return fmt.Errorf("Unexpected counter set. Wanted %v valid counters, got %v valid counters, list %v", length, valid, list)
		}
		return nil
	}

	var base Base

	// Verify an empty counter set.
	list := base.Counters()
	if len(list) != 0 {
		t.Errorf("Unexpected number of entries. Wanted 0, got %v, list %v.", len(list), list)
	}

	// Add 4 counters to the counter set.
	testLength := 4
	if err := base.InitCounters("prefix", "desc", fwdpb.CounterId_RX_DROP_PACKETS, fwdpb.CounterId_RX_DROP_OCTETS, fwdpb.CounterId_RX_PACKETS, fwdpb.CounterId_RX_OCTETS); err != nil {
		t.Fatalf("InitCounters failed, %v", err)
	}
	list = base.Counters()
	if len(list) != testLength {
		t.Errorf("Unexpected number of entries. Wanted %v, got %v, list %v.", testLength, len(list), list)
	}

	// Increment all initialized counters and check their values.
	for _, counter := range list {
		base.Increment(counter.ID, uint32(counter.ID))
	}
	if err := validate(&base, testLength); err != nil {
		t.Error(err)
	}

	// Increment a non existing counter.
	base.Increment(fwdpb.CounterId_TX_OCTETS, 10)
	if err := validate(&base, testLength); err != nil {
		t.Error(err)
	}
}
