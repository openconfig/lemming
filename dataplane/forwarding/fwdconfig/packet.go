// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fwdconfig

import (
	"encoding/binary"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// PacketFieldBytesBuilder is a builder for PacketFieldBytes.
type PacketFieldBytesBuilder struct {
	bytes    []byte
	field    fwdpb.PacketFieldNum
	instance uint32
}

// PacketFieldBytes creates a new PacketFieldBytesBuilder
func PacketFieldBytes(field fwdpb.PacketFieldNum) *PacketFieldBytesBuilder {
	return &PacketFieldBytesBuilder{}
}

// WithBytes sets the bytes value.
func (pfb *PacketFieldBytesBuilder) WithBytes(b []byte) *PacketFieldBytesBuilder {
	pfb.bytes = b
	return pfb
}

// WithUint64 sets the bytes value with big endian encoded uint.
func (pfb *PacketFieldBytesBuilder) WithUint64(d uint64) *PacketFieldBytesBuilder {
	pfb.bytes = binary.BigEndian.AppendUint64(nil, d)
	return pfb
}

// WithUint32 sets the bytes value with big endian encoded uint.
func (pfb *PacketFieldBytesBuilder) WithUint32(d uint32) *PacketFieldBytesBuilder {
	pfb.bytes = binary.BigEndian.AppendUint32(nil, d)
	return pfb
}

// WithUint16 sets the bytes value with big endian encoded uint.
func (pfb *PacketFieldBytesBuilder) WithUint16(d uint16) *PacketFieldBytesBuilder {
	pfb.bytes = binary.BigEndian.AppendUint16(nil, d)
	return pfb
}

// Build returns a new PacketFieldBytes.
func (pfb *PacketFieldBytesBuilder) Build() *fwdpb.PacketFieldBytes {
	return &fwdpb.PacketFieldBytes{
		Bytes: pfb.bytes,
		FieldId: &fwdpb.PacketFieldId{
			Field: &fwdpb.PacketField{
				FieldNum: pfb.field,
				Instance: pfb.instance,
			},
		},
	}
}
