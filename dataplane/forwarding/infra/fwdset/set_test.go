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

package fwdset

import (
	"fmt"
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// TestSetContains tests the contains operation on a set.
func TestSetContains(t *testing.T) {
	ctx := fwdcontext.New("test", "fwd")

	// check describes a member whose membership is checked in the set.
	type check struct {
		member   []byte // member to check in the set
		contains bool   // true if the member should be found in the set
	}

	tests := []struct {
		members [][]byte // members in the test set
		checks  []check  // checkes performed on the set
	}{
		// Set with no member.
		{
			checks: []check{
				{
					member:   []byte{0x00},
					contains: false,
				},
				{
					member:   []byte{0x01},
					contains: false,
				},
				{
					member:   []byte{},
					contains: false,
				},
			},
		},
		// Set with one member.
		{
			members: [][]byte{
				[]byte{0x00},
			},
			checks: []check{
				{
					member:   []byte{0x00},
					contains: true,
				},
				{
					member:   []byte{0x01},
					contains: false,
				},
			},
		},
		// Set with one member, the empty string.
		{
			members: [][]byte{
				[]byte{},
			},
			checks: []check{
				{
					member:   []byte{0x00},
					contains: false,
				},
				{
					member:   []byte{},
					contains: true,
				},
			},
		},
		// Set with multiple members.
		{
			members: [][]byte{
				[]byte{0x00},
				[]byte{0x00, 0x01},
				[]byte{0x01},
				[]byte{0x01, 0x02, 0x03},
			},
			checks: []check{
				{
					member:   []byte{0x00},
					contains: true,
				},
				{
					member:   []byte{},
					contains: false,
				},
				{
					member:   []byte{0x01},
					contains: true,
				},
				{
					member:   []byte{0x02},
					contains: false,
				},
				{
					member:   []byte{0x01, 0x03, 0x02},
					contains: false,
				},
				{
					member:   []byte{0x01, 0x02, 0x03},
					contains: true,
				},
			},
		},
	}

	// For each test, create the specified set and perform the specified checks.
	for id, test := range tests {
		name := fmt.Sprintf("Set %v", id)

		c, err := New(ctx, &fwdpb.SetId{ObjectId: &fwdpb.ObjectId{Id: name}})
		if err != nil {
			t.Errorf("%d: Failed to create set %v, err %v.", id, name, err)
			continue
		}
		c.Update(test.members)
		t.Logf("Created set %v with %v", name, c)
		for pos, check := range test.checks {
			if got := c.Contains(check.member); got != check.contains {
				t.Errorf("%d.%d: Unexpected result for Find(%v). Got %v, want %v.", id, pos, check.member, got, check.contains)
			}
		}
	}
}

// TestSetOperation tests the operations on sets.
func TestSetOperations(t *testing.T) {
	ctx := fwdcontext.New("test", "fwd")

	// List of sets created in the test.
	creates := []struct {
		name    string // name of the set
		success bool   // true if the set creation should be successful
	}{
		{
			name:    "c1",
			success: true,
		},
		{
			name:    "c2",
			success: true,
		},
		{
			name:    "c1",
			success: false,
		},
	}

	// List of sets searched in the test.
	searchs := []struct {
		name    string // name of the set
		success bool   // true if the set creation should be successful
	}{
		{
			name:    "c1",
			success: true,
		},
		{
			name:    "c2",
			success: true,
		},
		{
			name:    "c3",
			success: false,
		},
	}

	// Create all sets.
	for _, create := range creates {
		_, err := New(ctx, &fwdpb.SetId{ObjectId: &fwdpb.ObjectId{Id: create.name}})
		switch {
		case err == nil && create.success == false:
			t.Errorf("Create(%v) succeeded. Want error", create.name)

		case err != nil && create.success == true:
			t.Errorf("Create(%v) failed. Got error %v", create.name, err)
		}
	}

	// Search the specified sets.
	for _, search := range searchs {
		cid := &fwdpb.SetId{
			ObjectId: &fwdpb.ObjectId{
				Id: search.name,
			},
		}
		_, err := Find(ctx, cid)
		switch {
		case err == nil && search.success == false:
			t.Errorf("Find(%v) succeeded. Want error", search.name)

		case err != nil && search.success == true:
			t.Errorf("Find(%v) failed. Got error %v", search.name, err)

		}
	}
}
