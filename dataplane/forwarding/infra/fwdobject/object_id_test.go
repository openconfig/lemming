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
	"testing"
)

// TestPool tests acquire / release of object-ids from an object id pool.
func TestPool(t *testing.T) {
	var err error
	poolSize := 5
	pool := NewNIDPool(uint64(poolSize))

	// Exhaust the pool.
	id := make([]NID, poolSize)
	for count := 0; count < poolSize; count++ {
		if id[count], err = pool.Acquire(); err != nil {
			t.Errorf("Acquire attempt %v failed. Error %s.", count, err)
		} else {
			t.Logf("Allocated id %v.", id[count])
		}
	}

	// Allocate one more id from the exhausted pool.
	if test, err := pool.Acquire(); err == nil {
		t.Errorf("Incorrectly allocated id %v.", test)
	} else {
		t.Logf("Expected error, err %s", err)
	}

	// Release an id to the pool and re-acquire.
	pool.Release(id[0])
	test, err := pool.Acquire()
	if err != nil {
		t.Errorf("Acquire attempt failed, err %s", err)
	}

	if test != id[0] {
		t.Errorf("Obtained unexpected id. Got %v want %v.", test, id[0])
	}
}
