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
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/tabletestutil"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdaction/actions"
)

// A lookupTest defines a pattern to match in the table and
// the corresponding expected result.
type lookupTest struct {
	pattern []byte
	result  *key
}

// An opTest defines an operation on the prefix table with the specified key
// and its expected result.
type opTest struct {
	key    *key
	hasErr bool
}

// A matchTest defines a table of prefixes and set of lookups to be performed
// on the table.
type matchTest struct {
	adds    []opTest
	deletes []opTest
	lookups []lookupTest
}

// A dataAction is an action that stores some data.
// It satisfies the interface fwdaction.Action.
type dataAction struct {
	data *key
}

func (*dataAction) Process(_ fwdpacket.Packet, _ fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	return nil, fwdaction.DROP
}

func (t *dataAction) String() string {
	return t.data.String()
}

// prefixActions creates Actions from the specified data.
func prefixActions(data *key) fwdaction.Actions {
	return []*fwdaction.ActionAttr{fwdaction.NewActionAttr(&dataAction{
		data: data,
	}, false)}
}

// prefixResult returns data encoded using prefixActions.
func prefixResult(actions fwdaction.Actions) *key {
	if len(actions) != 1 {
		return nil
	}
	a := actions[0].Action()
	if test, ok := a.(*dataAction); ok {
		return test.data
	}
	return nil
}

// match returns true if an error exists as expected.
func match(err error, expect bool) bool {
	return (err != nil) == (expect)
}

// TestPrefixMatch tests prefix matching in the prefix table.
func TestPrefixMatch(t *testing.T) {
	tests := []matchTest{
		// Table with no entries.
		{
			lookups: []lookupTest{
				{
					pattern: []byte{},
					result:  nil,
				},
				{
					pattern: []byte{0x01, 0x23},
					result:  nil,
				},
			},
		},
		// Table with disjoint first levels.
		{
			adds: []opTest{
				{
					key: newKey([]byte{0x18}, 5),
				},
				{
					key: newKey([]byte{0x12}, 6),
				},
			},
			lookups: []lookupTest{
				{
					pattern: []byte{0x18, 0x20},
					result:  newKey([]byte{0x18}, 5),
				},
				{
					pattern: []byte{0x12, 0x20},
					result:  newKey([]byte{0x12}, 6),
				},
				{
					pattern: []byte{0x22, 0x20},
					result:  nil,
				},
				{
					pattern: []byte{0x00},
					result:  nil,
				},
				{
					pattern: []byte{0x13},
					result:  newKey([]byte{0x12}, 6),
				},
				{
					pattern: []byte{0x14},
					result:  nil,
				},
			},
		},
		// Table with multiple levels.
		{
			adds: []opTest{
				{
					key: newKey([]byte{0x00, 0x00, 0x00}, 10),
				},
				{
					key: newKey([]byte{0x00, 0x10, 0x00}, 10),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x10}, 21),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0xE0}, 18),
				},
				{
					key: newKey([]byte{0x00}, 1),
				},
				{
					key: newKey([]byte{0x00}, 2),
				},
			},
			lookups: []lookupTest{
				{
					pattern: []byte{0x10},
					result:  newKey([]byte{0x00}, 2),
				},
				{
					pattern: []byte{0x70},
					result:  newKey([]byte{0x10}, 1),
				},
				{
					pattern: []byte{0x20},
					result:  newKey([]byte{0x20}, 2),
				},
				{
					pattern: []byte{0x23},
					result:  newKey([]byte{0x20}, 2),
				},
				{
					pattern: []byte{0x80},
					result:  nil,
				},
				{
					pattern: []byte{0x00, 0x00, 0x00, 0x01},
					result:  newKey([]byte{0x00, 0x00, 0x00}, 10),
				},
				{
					pattern: []byte{0x00, 0x10, 0x00, 0x01},
					result:  newKey([]byte{0x00, 0x10, 0x00}, 10),
				},
				{
					pattern: []byte{0x00, 0x00, 0xF0, 0x01},
					result:  newKey([]byte{0x00, 0x00, 0xE0}, 18),
				},
				{
					pattern: []byte{0xF0, 0x00, 0xF0, 0x01},
					result:  nil,
				},
			},
		},
	}

	for id, test := range tests {
		table := &Table{
			root: newLevel(newKey([]byte{}, 0), nil),
		}

		// Add all entries into the table.
		for _, entry := range test.adds {
			table.add(entry.key.Copy(), prefixActions(entry.key.Copy()))
		}
		t.Logf("#%d: Prefix table is %v:\n%v", id, table, strings.Join(table.Entries(), "\n"))

		// Perform all the matches.
		for index, lookup := range test.lookups {
			_, a := table.match(lookup.pattern)
			result := prefixResult(a)
			switch {
			case result == nil && lookup.result != nil:
				t.Errorf("#%d.%d Incorrect result found for pattern %x, got nil want %v.", id, index, lookup.pattern, lookup.result)

			case result != nil && lookup.result == nil:
				t.Errorf("#%d.%d Incorrect result found for pattern %x, got %v want nil.", id, index, lookup.pattern, result)

			case result != nil && lookup.result != nil && !result.IsEqual(lookup.result):
				t.Errorf("#%d.%d Incorrect result found for pattern %x, got %v want %v.", id, index, lookup.pattern, result, lookup.result)
			}
		}
	}
}

// TestPrefixDelete tests prefix deletion from the prefix table.
func TestPrefixDelete(t *testing.T) {
	tests := []matchTest{
		// Table with no entries.
		{
			deletes: []opTest{
				{
					key:    newKey([]byte{0x18}, 5),
					hasErr: true,
				},
				{
					key:    newKey([]byte{0x12}, 6),
					hasErr: true,
				},
			},
			lookups: []lookupTest{
				{
					pattern: []byte{},
					result:  nil,
				},
				{
					pattern: []byte{0x01, 0x23},
					result:  nil,
				},
			},
		},
		// Table with disjoint first levels.
		{
			adds: []opTest{
				{
					key: newKey([]byte{0x18}, 5),
				},
				{
					key: newKey([]byte{0x12}, 6),
				},
				{
					key: newKey([]byte{0xF8}, 1),
				},
				{
					key: newKey([]byte{0xF8}, 2),
				},
			},
			deletes: []opTest{
				{
					key:    newKey([]byte{0xF8}, 1),
					hasErr: false,
				},
				{
					key:    newKey([]byte{0x12}, 6),
					hasErr: false,
				},
			},
			lookups: []lookupTest{
				{
					pattern: []byte{0x18, 00},
					result:  newKey([]byte{0x18}, 5),
				},
				{
					pattern: []byte{0x80, 0x01},
					result:  nil,
				},
				{
					pattern: []byte{0x14},
					result:  nil,
				},
				{
					pattern: []byte{0xC1, 0xAB},
					result:  newKey([]byte{0xF8}, 2),
				},
			},
		},
		// Table with multiple levels being deleted.
		{
			adds: []opTest{
				{
					key: newKey([]byte{0x18}, 5),
				},
				{
					key: newKey([]byte{0x12}, 6),
				},
				{
					key: newKey([]byte{0x12}, 7),
				},
				{
					key: newKey([]byte{0x12}, 8),
				},
				{
					key: newKey([]byte{0xF8}, 1),
				},
				{
					key: newKey([]byte{0xF8}, 2),
				},
				{
					key: newKey([]byte{0xF8}, 4),
				},
				{
					key: newKey([]byte{0xF8}, 5),
				},
			},
			deletes: []opTest{
				{
					key:    newKey([]byte{0xF8}, 1),
					hasErr: false,
				},
				{
					key:    newKey([]byte{0xF8}, 5),
					hasErr: false,
				},
				{
					key:    newKey([]byte{0xF8}, 4),
					hasErr: false,
				},
				{
					key:    newKey([]byte{0xF8}, 2),
					hasErr: false,
				},
			},
			lookups: []lookupTest{
				{
					pattern: []byte{0xF8},
					result:  nil,
				},
				{
					pattern: []byte{0x18, 00},
					result:  newKey([]byte{0x18}, 5),
				},
				{
					pattern: []byte{0x14},
					result:  nil,
				},
				{
					pattern: []byte{0x12},
					result:  newKey([]byte{0x12}, 8),
				},
			},
		},
		// Table empty after deletion.
		{
			adds: []opTest{
				{
					key: newKey([]byte{0x18}, 5),
				},
				{
					key: newKey([]byte{0x12}, 6),
				},
				{
					key: newKey([]byte{0xF8}, 1),
				},
				{
					key: newKey([]byte{0xF8}, 2),
				},
			},
			deletes: []opTest{
				{
					key:    newKey([]byte{0x18}, 5),
					hasErr: false,
				},
				{
					key:    newKey([]byte{0x12}, 6),
					hasErr: false,
				},
				{
					key:    newKey([]byte{0xF8}, 1),
					hasErr: false,
				},
				{
					key:    newKey([]byte{0xF8}, 2),
					hasErr: false,
				},
			},
			lookups: []lookupTest{
				{
					pattern: []byte{0x18, 00},
					result:  nil,
				},
				{
					pattern: []byte{0xC1, 0xAB},
					result:  nil,
				},
			},
		},
		// Sequence of operations used to export delete bug.
		{
			adds: []opTest{
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7f, 0x00, 0x01, 0x01}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7f, 0x00, 0x01, 0x02}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7f, 0x7f, 0x7f, 0x7f}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7f, 0x7f, 0x7f, 0x7e}, 152),
				},

				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xc8}, 150),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xc8}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xc9}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xca}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xcc}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xe0}, 150),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xe0}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xe2}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xe4}, 150),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xe4}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xe6}, 152),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xe8}, 150),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xec}, 150),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xf0}, 150),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xf4}, 150),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xf8}, 150),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xfc}, 150),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x17, 0xc0, 0x00}, 145),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x17, 0xc0, 0x20}, 147),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x17, 0xc1}, 144),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x17, 0xc3}, 144),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x17, 0xc4}, 143),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x17, 0xc5}, 144),
				},
				{
					key: newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x17, 0xd0}, 140),
				},
			},
			deletes: []opTest{
				{
					key:    newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x17, 0xc3}, 144),
					hasErr: false,
				},
				{
					key:    newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x17, 0xc0, 0x020}, 147),
					hasErr: false,
				},
				{
					key:    newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x17, 0xc5}, 144),
					hasErr: false,
				},
				{
					key:    newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xca}, 152),
					hasErr: false,
				},
			},
			lookups: []lookupTest{
				{
					pattern: []byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xc8},
					result:  newKey([]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xac, 0x13, 0xca, 0xc8}, 152),
				},
			},
		},
	}

	for id, test := range tests {
		table := &Table{
			root: newLevel(newKey([]byte{}, 0), nil),
		}

		// Add all entries into the table.
		for _, entry := range test.adds {
			table.add(entry.key.Copy(), prefixActions(entry.key.Copy()))
		}
		// Delete specified entries into the table.
		for _, entry := range test.deletes {
			err := table.remove(entry.key)
			if !match(err, entry.hasErr) {
				t.Errorf("#%d Incorrect result found for remove %x, got error %v want any error %v.", id, entry.key, err, entry.hasErr)
			}
		}
		t.Logf("#%d: Prefix table is %v:\n%v", id, table, strings.Join(table.Entries(), "\n"))

		// Perform all the matches.
		for index, lookup := range test.lookups {
			_, a := table.match(lookup.pattern)
			result := prefixResult(a)
			switch {
			case result == nil && lookup.result != nil:
				t.Errorf("#%d.%d Incorrect result found for pattern %x, got nil want %v.", id, index, lookup.pattern, lookup.result)

			case result != nil && lookup.result == nil:
				t.Errorf("#%d.%d Incorrect result found for pattern %x, got %v want nil.", id, index, lookup.pattern, result)

			case result != nil && lookup.result != nil && !result.IsEqual(lookup.result):
				t.Errorf("#%d.%d Incorrect result found for pattern %x, got %v want %v.", id, index, lookup.pattern, result, lookup.result)
			}
		}
	}
}

// prefixDesc creates a desc for a prefix entry.
func prefixDesc(versionBytes, versionMask, vrfBytes, vrfMask []byte) *fwdpb.EntryDesc {
	desc := &fwdpb.EntryDesc{}
	prefix := &fwdpb.PrefixEntryDesc{
		Fields: []*fwdpb.PacketFieldMaskedBytes{
			{
				FieldId: &fwdpb.PacketFieldId{
					Field: &fwdpb.PacketField{
						FieldNum: fwdpb.PacketFieldNum_IP_VERSION.Enum(),
					},
				},
				Bytes: versionBytes,
				Masks: versionMask,
			},
			{
				FieldId: &fwdpb.PacketFieldId{
					Field: &fwdpb.PacketField{
						FieldNum: fwdpb.PacketFieldNum_PACKET_VRF.Enum(),
					},
				},
				Bytes: vrfBytes,
				Masks: vrfMask,
			},
		},
	}
	proto.SetExtension(desc, fwdpb.E_PrefixEntryDesc_Extension, prefix)
	return desc
}

// prefixMatchTable creates a prefix table.
func prefixMatchTable(ctx *fwdcontext.Context) (fwdtable.Table, error) {
	// Prefix match table descriptor.
	desc := &fwdpb.TableDesc{
		TableType: fwdpb.TableType_PREFIX_TABLE.Enum(),
		Actions:   tabletestutil.ActionDesc(),
		TableId:   fwdtable.MakeID(fwdobject.NewID("prefixtable")),
	}

	// Setup the key specification.
	prefix := &fwdpb.PrefixTableDesc{
		FieldIds: []*fwdpb.PacketFieldId{
			{
				Field: &fwdpb.PacketField{
					FieldNum: fwdpb.PacketFieldNum_IP_VERSION.Enum(),
				},
			},
			{
				Field: &fwdpb.PacketField{
					FieldNum: fwdpb.PacketFieldNum_PACKET_VRF.Enum(),
				},
			},
		},
	}
	proto.SetExtension(desc, fwdpb.E_PrefixTableDesc_Extension, prefix)
	return fwdtable.New(ctx, desc)
}

// TestPrefixTable tests various operations on an prefix match table.
func TestPrefixTable(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Register a mock parser.
	const validSize = 10
	parser := mock_fwdpacket.NewMockParser(ctrl)
	parser.EXPECT().MaxSize(gomock.Any()).Return(validSize).AnyTimes()
	parser.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	fwdpacket.Register(parser)

	ctx := fwdcontext.New("test", "fwd")

	table, err := prefixMatchTable(ctx)
	if err != nil {
		t.Errorf("Prefix match table create failed, err %v.", err)
	}
	if entries := table.Entries(); len(entries) != 0 {
		t.Errorf("Incorrect number of table entries. Got %v, want 0.", len(entries))
	}

	// Operate on *count* entries in the table.
	const count = 10

	// Add entries to the table without specifying any masks.
	for index := 0; index < count; index++ {
		if err := table.AddEntry(prefixDesc([]byte{uint8(index)}, nil, []byte{uint8(index)}, nil), tabletestutil.ActionDesc()); err != nil {
			t.Errorf("AddEntry failed, err %v.", err)
		}
	}
	entries := table.Entries()
	t.Logf("List %v.", entries)
	if len(entries) != count {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	}

	// Add entries to the table specifying a masks for all bits. The
	// entry count should not change.
	for index := 0; index < count; index++ {
		if err := table.AddEntry(prefixDesc([]byte{uint8(index)}, []byte{0xFF}, []byte{uint8(index)}, []byte{0xFF}), tabletestutil.ActionDesc()); err != nil {
			t.Errorf("AddEntry failed, err %v.", err)
		}
	}
	entries = table.Entries()
	t.Logf("List %v.", entries)
	if len(entries) != count {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	}

	// Update an entry in the table.
	if err := table.AddEntry(prefixDesc([]byte{uint8(1)}, nil, []byte{uint8(1)}, nil), tabletestutil.ActionDesc()); err != nil {
		t.Errorf("AddEntry failed, err %v.", err)
	}
	entries = table.Entries()
	t.Logf("List %v.", entries)
	if len(entries) != count {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	}

	// Remove an entry from the table.
	if err := table.RemoveEntry(prefixDesc([]byte{uint8(1)}, nil, []byte{uint8(1)}, nil)); err != nil {
		t.Errorf("RemoveEntry failed, err %v.", err)
	}
	entries = table.Entries()
	t.Logf("List %v.", entries)
	if len(entries) != (count - 1) {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	}

	// Add an entry again.
	if err := table.AddEntry(prefixDesc([]byte{uint8(1)}, nil, []byte{uint8(1)}, nil), tabletestutil.ActionDesc()); err != nil {
		t.Errorf("AddEntry failed, err %v.", err)
	}
	entries = table.Entries()
	t.Logf("List %v.", entries)
	if len(entries) != count {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	}

	// Remove a non existing entry.
	if err := table.RemoveEntry(prefixDesc([]byte{uint8(count + 1)}, nil, []byte{uint8(1)}, nil)); err != nil {
		t.Logf("RemoveEntry failed as expected, err %v.", err)
	} else {
		t.Error("Unexpected entry removed")
	}
	entries = table.Entries()
	t.Logf("List %v.", entries)
	if len(entries) != count {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), count)
	}

	// Add entries to the table specifying a masks with a prefix. The
	// entry count should double.
	for index := 0; index < count; index++ {
		if err := table.AddEntry(prefixDesc([]byte{uint8(index)}, []byte{0xFF}, []byte{uint8(index)}, []byte{0xFE}), tabletestutil.ActionDesc()); err != nil {
			t.Errorf("AddEntry failed, err %v.", err)
		}
	}
	entries = table.Entries()
	t.Logf("List %v.", entries)
	if len(entries) != 2*count {
		t.Errorf("Incorrect number of table entries. Got %v, want %v.", len(entries), 2*count)
	}

	// Clear the table and ensure that there are no entres.
	table.Clear()
	if entries := table.Entries(); len(entries) != 0 {
		t.Errorf("Incorrect number of table entries. Got %v, want 0.", len(entries))
	}

	// Add a new entry into the table.
	if err := table.AddEntry(prefixDesc([]byte{uint8(1)}, nil, []byte{uint8(1)}, nil), tabletestutil.ActionDesc()); err != nil {
		t.Errorf("AddEntry failed, err %v.", err)
	}
	if entries := table.Entries(); len(entries) != 1 {
		t.Errorf("Incorrect number of table entries. Got %v, want 1.", len(entries))
	}
}
