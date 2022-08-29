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

package exact

import (
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/tableutil"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// EntryDesc describes a table entry as a set of pairs of packet-field-ids and
// packet-field-values.
type EntryDesc []*fwdpb.PacketFieldBytes

// newExactKey makes a exact match key from an EntryDesc.
// Note that each field is padded to the maxiumum size when forming the key.
func newExactKey(kd tableutil.KeyDesc, ed EntryDesc) (tableutil.Key, error) {
	// Build a map of fields and verify the byte lengths and uniqueness.
	set := make(map[fwdpacket.FieldID]tableutil.Key)
	for _, desc := range ed {
		id := fwdpacket.NewFieldID(desc.GetFieldId())
		bytes := desc.GetBytes()
		if err := fwdpacket.Validate(id, len(bytes)); err != nil {
			return nil, fmt.Errorf("exact: newExactKey failed, err %v", err)
		}
		set[id] = fwdpacket.Pad(id, bytes)
	}

	// Ensure that there are no duplicate fields.
	if len(set) != len(ed) {
		return nil, fmt.Errorf("exact: newExactKey failed, field-ids %v contain duplicate field-id", ed)
	}

	// Build the key using the bytes in the sequence defined by the desc.
	var k tableutil.Key
	for _, id := range kd {
		bytes, ok := set[id]
		if !ok {
			return nil, fmt.Errorf("exact: newExactKey failed, missing field %v in %v", id, ed)
		}
		k = append(k, bytes...)
	}
	return k, nil
}
