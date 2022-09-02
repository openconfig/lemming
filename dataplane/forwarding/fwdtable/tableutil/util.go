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

// Package tableutil contains utilites used to implement tables in Lucius like
// prefix, flow and exact match.
package tableutil

import (
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A Key is a sequence of bytes used to identify entries in a table.
type Key []byte

// A KeyDesc describes the structure of a key as a sequence of packet field IDs.
// It can make keys from packets or entry descriptors.
type KeyDesc []fwdpacket.FieldID

// MakeKeyDesc makes a keyDesc from a series of packet field id.
func MakeKeyDesc(fields []*fwdpb.PacketFieldId) KeyDesc {
	var kd KeyDesc
	for _, field := range fields {
		kd = append(kd, fwdpacket.NewFieldID(field))
	}
	return kd
}

// MakePacketKey makes a key from a packet. The key desc determines the
// sequence of fields used to create the key. If the key desc contains a
// field absent in the desc, then the corresponding bytes are set to
// zero. Each field is padded to the maximum size when forming the key.
func (kd KeyDesc) MakePacketKey(packet fwdpacket.Packet) Key {
	var k Key
	for _, id := range kd {
		field, err := packet.Field(id)
		if err != nil {
			field = []byte{}
		}
		field = fwdpacket.Pad(id, field)
		k = append(k, field...)
	}
	return k
}

// Size returns the size of the key in bytes built using the key desc.
func (kd KeyDesc) Size() int {
	size := 0
	for _, id := range kd {
		size += fwdpacket.MaxSize(id)
	}
	return size
}
