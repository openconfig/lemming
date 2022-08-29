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

package protocol

import (
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func init() {
	// Register a packet parser.
	fwdpacket.Register(parser{})
}

// A parser can create a packet from a slice of bytes.
type parser struct {
}

// New creates a new packet interpreting a slice of bytes to begin with the
// specified packet header.
func (parser) New(l2 fwdpb.PacketHeaderId, f []byte) (fwdpacket.Packet, error) {
	pkt, err := NewPacket(l2, frame.NewFrame(f))
	if err != nil {
		return nil, err
	}
	return pkt, nil
}

// Validate returns an error if the size of the specified field in invalid.
//
// Valid lengths for a field may be specified in multiple ways
//  1. In case of UDF, the valid size is encoded in the FieldID.
//  2. In case of field numbers, valid sizes are either specified as a range
//     or set of discrete sizes.
func (parser) Validate(id fwdpacket.FieldID, size int) error {
	if id.IsUDF {
		if int(id.Size) != size {
			return fmt.Errorf("Validate failed for field %v, got %v want %v", id, size, id.Size)
		}
		return nil
	}

	attr, ok := FieldAttr[id.Num]
	if !ok {
		return fmt.Errorf("Validate failed for field %v, unknown field", id)
	}
	for _, s := range attr.Sizes {
		if size == s {
			return nil
		}
	}
	return fmt.Errorf("Validate failed for field %v, %v does not satisfy %+v", id, size, attr)
}

// MaxSize returns the maximum size of a field. In case of UDF, the maximum
// size is encoded within the FieldID. For other cases, the default size is
// the maximum size of the field.
func (parser) MaxSize(id fwdpacket.FieldID) int {
	if id.IsUDF {
		return int(id.Size)
	}
	if attr, ok := FieldAttr[id.Num]; ok {
		return attr.DefaultSize
	}
	return 0
}
