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

// Package frame contains utilities to implement various network protocols
// and describe their relationship within Lucius. Each supported network
// protocol implements the interface protocol.Handler and registers a parse
// and (an optional) add handler.
package frame

import (
	"encoding/binary"
	"fmt"
)

// ByteBitCount is the number of bits in a byte.
const ByteBitCount = 8

// Uint32BitCount is the number of bits in a uint32.
const Uint32BitCount = 4 * ByteBitCount

// Uint64BitCount is the number of bits in a uint64.
const Uint64BitCount = 8 * ByteBitCount

// A Field is a slice of bytes within a header. Fields of
// certain sizes may contain numeric data in BigEndian format.
type Field []byte

// Resize returns a slice of bytes either padded or truncated to the specified
// number of bytes. The intent of resize is to be used with the frame
// utilities which operate on bytes that represent data in network order.
// For padding data, zeroes are added at the MSB and when truncating data, the
// MSB is stripped.
func Resize(field []byte, size int) []byte {
	switch curr := len(field); {
	case curr < size:
		return append(make([]byte, size-curr), field...)

	case curr > size:
		return field[curr-size : curr]

	default:
		return field
	}
}

// SetValue sets a numeric field's value.
func (f Field) SetValue(value uint) {
	switch len(f) {
	case 1:
		f[0] = uint8(value)
	case 2:
		binary.BigEndian.PutUint16(f, uint16(value))
	case 4:
		binary.BigEndian.PutUint32(f, uint32(value))
	case 8:
		binary.BigEndian.PutUint64(f, uint64(value))
	default:
		panic(fmt.Sprintf("protocol: SetValue failed, unsupported size of field %v", f))
	}
}

// Value gets a numeric field's value.
func (f Field) Value() uint {
	switch len(f) {
	case 1:
		return uint(f[0])
	case 2:
		return uint(binary.BigEndian.Uint16(f))
	case 4:
		return uint(binary.BigEndian.Uint32(f))
	case 8:
		return uint(binary.BigEndian.Uint64(f))
	default:
		panic(fmt.Sprintf("protocol: Value failed, unsupported size of field %v", f))
	}
}

// SetBits sets bitcount bits at position bitpos in a numeric field.
func (f Field) SetBits(bitpos, bitcount uint8, bitmask uint64) {
	if (int(bitpos+bitcount-1) >= (len(f) * ByteBitCount)) || (int(bitcount) > (Uint32BitCount * ByteBitCount)) {
		panic(fmt.Sprintf("protocol: SetBits failed, bitcount %v, bitpos %v is incorrect for field %x", bitpos, bitcount, f))
	}
	clr := ^(((1 << bitcount) - 1) << bitpos)
	set := bitmask << bitpos
	switch len(f) {
	case 1:
		f[0] = (uint8(f[0]) & uint8(clr)) | uint8(set)
	case 2:
		value := binary.BigEndian.Uint16(f)
		value = (value & uint16(clr)) | uint16(set)
		binary.BigEndian.PutUint16(f, value)
	case 4:
		value := binary.BigEndian.Uint32(f)
		value = (value & uint32(clr)) | uint32(set)
		binary.BigEndian.PutUint32(f, value)
	case 8:
		value := binary.BigEndian.Uint64(f)
		value = (value & uint64(clr)) | set
		binary.BigEndian.PutUint64(f, value)

	default:
		panic(fmt.Sprintf("protocol: SetBits failed, unsupported size of field %v", f))
	}
}

// BitField creates a new field using bitcount bits at position bitpos
// in a numeric field.
func (f Field) BitField(bitpos, bitcount uint8) Field {
	if bitpos+bitcount != 0 && int(bitpos+bitcount-1) >= (len(f)*ByteBitCount) {
		panic(fmt.Sprintf("protocol: BitField failed, bitcount %v, bitpos %v is incorrect for field %x", bitpos, bitcount, f))
	}
	bit := make(Field, len(f))
	mask := ((1 << bitcount) - 1)
	switch len(f) {
	case 1:
		bit[0] = (uint8(f[0]) >> bitpos) & uint8(mask)
		return bit
	case 2:
		value := binary.BigEndian.Uint16(f)
		value = (value >> bitpos) & uint16(mask)
		binary.BigEndian.PutUint16(bit, value)
		return bit
	case 4:
		value := binary.BigEndian.Uint32(f)
		value = (value >> bitpos) & uint32(mask)
		binary.BigEndian.PutUint32(bit, value)
		return bit
	case 8:
		value := binary.BigEndian.Uint64(f)
		value = (value >> bitpos) & uint64(mask)
		binary.BigEndian.PutUint64(bit, value)
		return bit
	default:
		panic(fmt.Sprintf("protocol: BitField failed, unsupported size of field %v", f))
	}
}

// Set sets all the bytes of the field.
func (f Field) Set(value []byte) error {
	value = Resize(value, len(f))
	copy(f, value)
	return nil
}

// Copy returns a copy of the field.
func (f Field) Copy() Field {
	result := make(Field, len(f))
	copy(result, f)
	return result
}

// A Header is a sequence of fields of a protocol within a frame.
type Header []byte

// Field returns a field of n bytes at offset off within the header.
func (h Header) Field(off, n int) Field {
	if len(h) < int(off+n) {
		return nil
	}
	return Field(h[off : off+n])
}

// A Frame is a set of bytes that can be read sequentially (or parsed) as a
// series of packet headers.
type Frame struct {
	buffer  []byte // Underlying slice of bytes.
	readPos int    // Position in buffer for the next read.
}

// NewFrame creates a new frame from a slice of bytes.
func NewFrame(buffer []byte) *Frame {
	return &Frame{
		buffer: buffer,
	}
}

// Len returns the number of unread bytes in the frame.
func (f *Frame) Len() int {
	return len(f.buffer) - f.readPos
}

// Peek creates a field of n bytes at offset off from the read position.
// If the specified bytes do not exist, an error is returned.
func (f *Frame) Peek(off, n int) (Field, error) {
	if f.Len() < off+n {
		return nil, fmt.Errorf("protocol: Peek failed, cannot read %v bytes at offset %v in frame of length %v", n, off, f.Len())
	}
	return Field(f.buffer[f.readPos+off : f.readPos+off+n]), nil
}

// ReadHeader reads a header of n bytes from the read position.
// If the specified bytes do not exist, an error is returned.
func (f *Frame) ReadHeader(n int) (Header, error) {
	if f.Len() < n {
		return nil, fmt.Errorf("protocol: ReadHeader failed, want %v bytes, got %v bytes", n, f.Len())
	}
	buffer := f.buffer[f.readPos : f.readPos+n]
	f.readPos += n
	return buffer, nil
}

// StripTrailing strips the trailing 'n' bytes from the buffer. Note that there
// must be at-least 'n' unread bytes in the frame.
func (f *Frame) StripTrailing(n int) error {
	if f.Len() < n {
		return fmt.Errorf("protocol: StripTrailing failed, want to strip %v bytes from %v unread bytes", n, f.Len())
	}
	f.buffer = f.buffer[:len(f.buffer)-n]
	return nil
}

// String formats unread bytes in a frame as a string.
func (f *Frame) String() string {
	return fmt.Sprintf("%v", f.buffer[f.readPos:])
}
