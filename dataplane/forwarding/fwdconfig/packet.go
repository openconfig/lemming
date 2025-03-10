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
	"math"

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
	return &PacketFieldBytesBuilder{
		field: field,
	}
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

// PacketFieldMaskedBytesBuilder is a builder for PacketFieldBytes.
type PacketFieldMaskedBytesBuilder struct {
	bytes    []byte
	mask     []byte
	field    fwdpb.PacketFieldNum
	instance uint32
}

// PacketFieldMaskedBytes creates a new PacketFieldBytesBuilder
func PacketFieldMaskedBytes(field fwdpb.PacketFieldNum) *PacketFieldMaskedBytesBuilder {
	return &PacketFieldMaskedBytesBuilder{
		field: field,
	}
}

// WithBytes sets the bytes and mask value.
func (b *PacketFieldMaskedBytesBuilder) WithBytes(bytes, mask []byte) *PacketFieldMaskedBytesBuilder {
	b.bytes = bytes
	b.mask = mask
	return b
}

// WithUint64 sets the bytes value with big endian encoded uint and the mask to an exact match.
func (b *PacketFieldMaskedBytesBuilder) WithUint64(d uint64) *PacketFieldMaskedBytesBuilder {
	b.bytes = binary.BigEndian.AppendUint64(nil, d)
	b.mask = binary.BigEndian.AppendUint64(nil, math.MaxUint64)
	return b
}

// WithUint32 sets the bytes value with big endian encoded uint and the mask to an exact match.
func (b *PacketFieldMaskedBytesBuilder) WithUint32(d uint32) *PacketFieldMaskedBytesBuilder {
	b.bytes = binary.BigEndian.AppendUint32(nil, d)
	b.mask = binary.BigEndian.AppendUint32(nil, math.MaxUint32)
	return b
}

// WithUint16 sets the bytes value with big endian encoded uint and the mask to an exact match.
func (b *PacketFieldMaskedBytesBuilder) WithUint16(d uint16) *PacketFieldMaskedBytesBuilder {
	b.bytes = binary.BigEndian.AppendUint16(nil, d)
	b.mask = binary.BigEndian.AppendUint16(nil, math.MaxUint16)
	return b
}

// Build returns a new PacketFieldMaskedBytes.
func (b *PacketFieldMaskedBytesBuilder) Build() *fwdpb.PacketFieldMaskedBytes {
	return &fwdpb.PacketFieldMaskedBytes{
		Bytes: b.bytes,
		Masks: b.mask,
		FieldId: &fwdpb.PacketFieldId{
			Field: &fwdpb.PacketField{
				FieldNum: b.field,
				Instance: b.instance,
			},
		},
	}
}

// PacketFieldIdBuilder is a builder for PacketField.
type PacketFieldIdBuilder struct {
	field        fwdpb.PacketFieldNum
	instance     uint32
	grp          fwdpb.PacketHeaderGroup
	offset, size uint32
}

func PacketFieldIDField(f fwdpb.PacketFieldNum, instance int) *PacketFieldIdBuilder {
	return &PacketFieldIdBuilder{
		field:    f,
		instance: uint32(instance),
	}
}

func PacketFieldIDBytes(grp fwdpb.PacketHeaderGroup, instance, offset, size int) *PacketFieldIdBuilder {
	return &PacketFieldIdBuilder{
		grp:      grp,
		instance: uint32(instance),
		offset:   uint32(offset),
		size:     uint32(size),
	}
}

func (b *PacketFieldIdBuilder) Build() *fwdpb.PacketFieldId {
	if b.grp != fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_UNSPECIFIED {
		return &fwdpb.PacketFieldId{Bytes: &fwdpb.PacketBytes{
			HeaderGroup: b.grp,
			Instance:    b.instance,
			Offset:      b.offset,
			Size:        b.size,
		}}
	}

	return &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: b.field,
			Instance: b.instance,
		},
	}
}
