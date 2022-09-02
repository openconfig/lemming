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

package fwdpacket

import (
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Constants used to indicate relative instances.
const (
	FirstField = 0            // The first (outermost) instance of a field.
	LastField  = (1 << 8) - 1 // The last (innermost) instance of a field.
)

// A FieldID identifies a field in a packet.
// It supports fields using field numbers and user defined fields (UDF).
type FieldID struct {
	IsUDF    bool                    // true if the field is a UDF.
	Instance uint8                   // Instance of the field in a packet.
	Num      fwdpb.PacketFieldNum    // Field number.
	Header   fwdpb.PacketHeaderGroup // Header type of UDF.
	Offset   uint8                   // UDF's offset within the header.
	Size     uint8                   // UDF size in bytes.
}

// NewFieldIDFromBytes creates a new field id from a set of bytes.
func NewFieldIDFromBytes(header fwdpb.PacketHeaderGroup, offset, size, instance uint32) FieldID {
	return FieldID{
		IsUDF:    true,
		Size:     uint8(size),
		Header:   header,
		Offset:   uint8(offset),
		Instance: uint8(instance),
	}
}

// NewFieldIDFromNum creates a new field id from a field number.
func NewFieldIDFromNum(num fwdpb.PacketFieldNum, instance uint32) FieldID {
	return FieldID{
		Num:      num,
		Instance: uint8(instance),
	}
}

// NewFieldID creates a new packet field id.
func NewFieldID(field *fwdpb.PacketFieldId) FieldID {
	if b := field.GetBytes(); b != nil {
		return NewFieldIDFromBytes(b.GetHeaderGroup(), b.GetOffset(), b.GetSize(), b.GetInstance())
	}
	if f := field.GetField(); f != nil {
		return NewFieldIDFromNum(f.GetFieldNum(), f.GetInstance())
	}
	return FieldID{}
}
