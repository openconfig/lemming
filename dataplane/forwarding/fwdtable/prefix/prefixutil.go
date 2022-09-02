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
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/tableutil"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A EntryDesc describes a table entry as a set of packet field ids and the corresponding
// mask and value.
type EntryDesc []*fwdpb.PacketFieldMaskedBytes

// newPrefixKey creates and returns a bitset representing the prefix described
// by EntryDesc. Note that each field is padded to the maximum size and only
// the last field in the key can have a mask.
func newPrefixKey(kd tableutil.KeyDesc, ed EntryDesc) (*key, error) {
	// Build a map of fields.
	set := make(map[fwdpacket.FieldID]*fwdpb.PacketFieldMaskedBytes)
	for _, desc := range ed {
		id := fwdpacket.NewFieldID(desc.GetFieldId())
		set[id] = desc
	}

	// Ensure that there are no duplicate fields.
	if len(set) != len(ed) {
		return nil, fmt.Errorf("prefix: makePrefixKey failed, field-ids %v contain duplicate field-id", ed)
	}

	// Build the bitset using the bytes in the sequence defined by the desc.
	// For each field, do the following:
	// 1. Validate the length of the value.
	// 2. Other than the last field, all fields either have no mask or have
	//    a mask specifying all bits.
	// 3. Pad the value (and the mask) appropriately.
	bytes := []byte{}
	bits := 0
	for pos, id := range kd {
		desc, ok := set[id]
		if !ok {
			return nil, fmt.Errorf("prefix: makePrefixKey failed, missing field %v in %v", id, ed)
		}
		value := desc.GetBytes()
		if err := fwdpacket.Validate(id, len(value)); err != nil {
			return nil, fmt.Errorf("prefix: makePrefixKey failed, err %v", err)
		}
		valueBits := Calculate(len(value))
		maskBits := valueBits

		// If a mask is specified, count the number of continguous bits.
		// If this is the last field in the key, use the mask to
		// determine the prefix length. Otherwise validate that the
		// mask contains the same number of bits as the value.
		if m := desc.GetMasks(); m != nil {
			maskBits = Count(m)
			if pos != len(kd)-1 && maskBits != valueBits {
				return nil, fmt.Errorf("prefix: makePrefixKey failed, field %v in entry %v at index %v has bad mask length. Got %v, want %v", id, ed, pos, maskBits, valueBits)
			}
		}

		value = fwdpacket.Pad(id, value)
		maskBits += Calculate(len(value)) - valueBits
		bytes = append(bytes, value...)
		bits += maskBits
	}
	return newKey(bytes, bits), nil
}
