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

package prefix

import (
	"testing"

	log "github.com/golang/glog"
)

// TestKeyEqual tests key equality.
func TestKeyEqual(t *testing.T) {
	// Create a set of keys identified by the index.
	specs := []struct {
		index int
		bytes []byte
		count int
	}{
		{
			index: 0,
			bytes: []byte{0x01, 0x20},
			count: 10,
		},
		{
			index: 1,
			bytes: []byte{0x01, 0x30},
			count: 10,
		},
		{
			index: 2,
			bytes: []byte{0x11, 0x30},
			count: 10,
		},
	}
	keys := make(map[int]*key)
	for _, s := range specs {
		k := newKey(s.bytes, s.count)
		keys[s.index] = k
		log.Infof("%d: Create key %v from %x/%v", s.index, k, s.bytes, s.count)
	}

	// Test if the specified keys are equal.
	tests := []struct {
		index1 int  // index of the first key in the comparison
		index2 int  // index of the second key in the comparison
		equal  bool // true if the comparision is expected to return true
	}{
		{
			// Test keys that have identical bits and bytes.
			index1: 0,
			index2: 0,
			equal:  true,
		},
		{
			// Test keys that have identical bits but not bytes.
			index1: 0,
			index2: 1,
			equal:  true,
		},
		{
			// Test keys that have different bits and bytes.
			index1: 0,
			index2: 2,
			equal:  false,
		},
	}
	for pos, test := range tests {
		s1 := keys[test.index1]
		s2 := keys[test.index2]
		if got := s1.IsEqual(s2); got != test.equal {
			t.Errorf("%d: Equality failed for %v and %v. Got %v, want %v.", pos, s1, s2, got, test.equal)
		}
	}
}

// TestKeyAccess test the set and get for each bit in the key.
func TestKeyAccess(t *testing.T) {
	key := newKey([]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}, 64)
	log.Infof("Key used is %v.", key)
	tests := []struct {
		pos int  // bit position to check
		set byte // expected value at the specified position.
	}{
		{
			pos: 4,
			set: 0x0,
		},
		{
			pos: 5,
			set: 0x0,
		},
		{
			pos: 6,
			set: 0x0,
		},
		{
			pos: 7,
			set: 0x1,
		},
		{
			pos: 8,
			set: 0x0,
		},
		{
			pos: 56,
			set: 0x1,
		},
		{
			pos: 57,
			set: 0x1,
		},
		{
			pos: 58,
			set: 0x1,
		},
		{
			pos: 59,
			set: 0x0,
		},
	}
	for pos, test := range tests {
		k := key.Copy()
		if result := k.Bit(test.pos); result != test.set {
			t.Errorf("%d: Bit get failed at index %v. Got %v, want %v.", pos, test.pos, result, test.set)
		}

		// flip the bit and verify.
		set := byte(0x0)
		if test.set == 0x0 {
			set = 0x1
		}
		k.Set(test.pos, set)
		if result := k.Bit(test.pos); result != set {
			t.Errorf("%d: Bit flip set at index %v. Got %v, want %v.", pos, test.pos, result, set)
		}
	}
}

// TestKeyPrefix tests HasPrefix and TrimPrefix operations.
func TestKeyPrefix(t *testing.T) {
	// Create a set of keys identified by the index.
	specs := []struct {
		index int
		bytes []byte
		count int
	}{
		{
			index: 0,
			bytes: []byte{0x12, 0x34},
			count: 11,
		},
		{
			index: 1,
			bytes: []byte{0x13, 0x00},
			count: 7,
		},
		{
			index: 2,
			bytes: []byte{0x13, 0x00},
			count: 8,
		},
	}
	keys := make(map[int]*key)
	for _, s := range specs {
		k := newKey(s.bytes, s.count)
		keys[s.index] = k
		log.Infof("%d: Create key %v from %x/%v", s.index, k, s.bytes, s.count)
	}

	// Test if the specified keys are equal.
	tests := []struct {
		index1 int  // index of the first key in the comparison
		index2 int  // index of the second key in the comparison
		prefix bool // true if the key at index2 is a prefix of the key at index1
	}{
		{
			// Test keys with identical bits.
			index1: 0,
			index2: 0,
			prefix: true,
		},
		{
			// Test keys that are a prefix.
			index1: 0,
			index2: 1,
			prefix: true,
		},
		{
			// Test keys that are not a prefix.
			index1: 0,
			index2: 2,
			prefix: false,
		},
	}
	for pos, test := range tests {
		s1 := keys[test.index1]
		s2 := keys[test.index2]
		if got := s1.HasPrefix(s2); got != test.prefix {
			t.Errorf("%d: HasPrefix failed for %v and %v. Got %v, want %v.", pos, s1, s2, got, test.prefix)
		}

		if test.prefix {
			k := s1.Copy()
			k.TrimPrefix(s2)
			if k.HasPrefix(s2) {
				t.Errorf("%d: HasPrefix succeed for %v and %v after trimming prefix.", pos, s1, s2)
			}
		}
	}
}

// TestKeyFindPrefix tests finding the prefix between pairs of keys
func TestKeyFindPrefix(t *testing.T) {
	tests := []struct {
		key1   *key // first key
		key2   *key // second key
		prefix *key // expected shared prefix
	}{
		{
			// First key is a complete prefix of the second.
			key1:   newKey([]byte{0x81, 0x7A}, 5),
			key2:   newKey([]byte{0x87, 0x9B}, 10),
			prefix: newKey([]byte{0x80}, 5),
		},
		{
			// Second key is a complete prefix of the first.
			key1:   newKey([]byte{0x87, 0x9B}, 10),
			key2:   newKey([]byte{0x81, 0x7A}, 5),
			prefix: newKey([]byte{0x80}, 5),
		},
		{
			// Keys share a common prefix.
			key1:   newKey([]byte{0x87, 0x9B}, 10),
			key2:   newKey([]byte{0x81, 0x7A}, 9),
			prefix: newKey([]byte{0x80}, 5),
		},
		{
			// Keys have no common prefix.
			key1:   newKey([]byte{0x77, 0x9B}, 10),
			key2:   newKey([]byte{0x81, 0x7A}, 9),
			prefix: newKey([]byte{}, 0),
		},
	}
	for pos, test := range tests {
		s1 := prefixKey(test.key1, test.key2)
		if !s1.IsEqual(test.prefix) {
			t.Errorf("%d: Split did not return the expected prefix. Got %v want %v, (%v,%v). ", pos, s1, test.prefix, test.key1, test.key2)
		}
	}
}
