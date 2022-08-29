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

package frame

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// errTest returns a non-empty string if err and str do not match.
func ErrTest(err error, str string) string {
	if err == nil {
		if str == "" {
			return ""
		}
		return `did not get expected error "` + str + `"`
	}
	if str == "" {
		return fmt.Sprintf("unexpected error %q", err.Error())
	}
	if !strings.Contains(err.Error(), str) {
		return fmt.Sprintf("got error %q, want %q", err.Error(), str)
	}
	return ""
}

// TestFrame tests operations on a frame.
func TestFrame(t *testing.T) {
	frame := NewFrame([]byte{0, 1, 2, 3, 4, 5, 6, 7})

	// tests lists the successive read / peeks performed on the frame and
	// their expected results.
	tests := []struct {
		length int
		result []byte
		err    string
	}{
		{
			length: 1,
			result: []byte{0},
		},
		{
			length: 2,
			result: []byte{1, 2},
		},
		{
			length: 3,
			result: []byte{3, 4, 5},
		},
		{
			length: 5,
			err:    "failed",
		},
		{
			length: 2,
			result: []byte{6, 7},
		},
	}

	// Peek into the frame at various offsets.
	offset := 0
	for _, test := range tests {
		f, err := frame.Peek(offset, test.length)
		if str := ErrTest(err, test.err); str != "" {
			t.Fatalf("Peek failed at offset %v, err %v.", offset, str)
		}
		if err == nil && !bytes.Equal(f, test.result) {
			t.Fatalf("Peek failed at offset %v, got %v want %v.", offset, f, test.result)
		}
		if err == nil {
			offset += test.length
		}
	}

	// Read the frame.
	for _, test := range tests {
		h, err := frame.ReadHeader(test.length)
		if str := ErrTest(err, test.err); str != "" {
			t.Fatalf("ReadHeader failed at offset, err %v.", str)
		}
		if err == nil && !bytes.Equal(h, test.result) {
			t.Fatalf("ReadHeader failed at offset, got %v want %v.", h, test.result)
		}
	}
}

// TestFieldValue tests operations on numeric field values.
func TestFieldValue(t *testing.T) {
	tests := []struct {
		field Field
		value uint
		set   uint
	}{
		{
			field: []byte{0x01},
			value: 0x01,
			set:   0x02,
		},
		{
			field: []byte{0x01, 0x02},
			value: 0x0102,
			set:   0x0203,
		},
		{
			field: []byte{0x01, 0x02, 0x03, 0x04},
			value: 0x01020304,
			set:   0x02030405,
		},
		{
			field: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			value: 0x0102030405060708,
			set:   0x0203040507070809,
		},
	}

	// Read the value of a fields and compare.
	for _, test := range tests {
		if test.field.Value() != test.value {
			t.Errorf("Value failed for field %v, got %v want %v.", test.field, test.field.Value(), test.value)
		}

		test.field.SetValue(test.set)
		if test.field.Value() != test.set {
			t.Errorf("Set failed for field %v, got %v want %v.", test.field, test.field.Value(), test.set)
		}
	}
}

// TestBitField tests bit operations on numeric field values.
func TestBitField(t *testing.T) {
	tests := []struct {
		field    Field
		bitpos   uint8
		bitcount uint8
		value    uint
		set      uint64
	}{
		{
			field:    []byte{0x12},
			bitpos:   4,
			bitcount: 2,
			value:    1,
			set:      3,
		},
		{
			field:    []byte{0x01, 0x02},
			bitpos:   8,
			bitcount: 3,
			value:    1,
			set:      7,
		},
		{
			field:    []byte{0x01, 0x02, 0x03, 0x04},
			bitpos:   16,
			bitcount: 3,
			value:    2,
			set:      7,
		},
		{
			field:    []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			bitpos:   8,
			bitcount: 3,
			value:    7,
			set:      2,
		},
	}

	// Read the value of a fields and compare.
	for _, test := range tests {
		val := test.field.BitField(test.bitpos, test.bitcount)
		if val.Value() != test.value {
			t.Errorf("Value failed for bitpos %v bitcount %v from %v, got %v want %v.", test.bitpos, test.bitcount, test.field, val.Value(), test.value)
		}

		test.field.SetBits(test.bitpos, test.bitcount, test.set)
		set := test.field.BitField(test.bitpos, test.bitcount)
		if set.Value() != uint(test.set) {
			t.Errorf("Set failed for bitpos %v bitcount %v from %v, got %v want %v.", test.bitpos, test.bitcount, test.field, set.Value(), test.set)
		}
	}
}

// TestFrameStripTrailingPeek tests stripping trailer bytes in a frame after
// a peek operation. Verify that when we strip trailer bytes, only unread bytes
// are considered irrespective of prior peeks.
func TestFrameStripTrailingPeek(t *testing.T) {
	tests := []struct {
		frame       *Frame //  Initial frame used in the test
		peekLength  int    // Number of bytes peeked before the frame is stripped
		stripLength int    // Number of bytes to be stripped
		result      []byte // Remaining bytes in the frame after it is stripped
		err         string // Error if the strip operation has an error
	}{
		{
			frame:       NewFrame([]byte{0, 1, 2, 3, 4, 5, 6, 7}),
			peekLength:  0,
			stripLength: 0,
			result:      []byte{0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			frame:       NewFrame([]byte{0, 1, 2, 3, 4, 5, 6, 7}),
			peekLength:  5,
			stripLength: 5,
			result:      []byte{0, 1, 2},
		},
		{
			frame:       NewFrame([]byte{0, 1, 2, 3, 4, 5, 6, 7}),
			peekLength:  5,
			stripLength: 8,
			result:      []byte{},
		},
		{
			// Verify that the frame is unchanged by a failed strip.
			frame:       NewFrame([]byte{0, 1, 2, 3, 4, 5, 6, 7}),
			peekLength:  5,
			stripLength: 9,
			err:         "failed",
			result:      []byte{0, 1, 2, 3, 4, 5, 6, 7},
		},
	}

	for _, test := range tests {
		if _, err := test.frame.Peek(0, test.peekLength); err != nil {
			t.Fatalf("Peek for length %v failed before stripping trailer %v bytes, err %v.", test.peekLength, test.stripLength, err)
		}

		err := test.frame.StripTrailing(test.stripLength)
		if str := ErrTest(err, test.err); str != "" {
			t.Fatalf("Stripping trailer of length %v returned unexpected result, err %v.", test.stripLength, str)
		}

		f, err := test.frame.Peek(0, test.frame.Len())
		if err != nil {
			t.Fatalf("Peek for length %v failed after stripping trailer %v bytes, err %v.", test.peekLength, test.stripLength, err)
		}
		if err == nil && !bytes.Equal(f, test.result) {
			t.Fatalf("Peek for length %v failed after stripping trailer %v bytes, got %v, want %v.", test.peekLength, test.stripLength, f, test.result)
		}
	}
}

// TestFrameStripTrailingRead tests stripping trailer bytes in a frame after
// a read operation. Verify that when we strip trailer bytes we consider the
// unread bytes.
func TestFrameStripTrailingRead(t *testing.T) {
	tests := []struct {
		frame       *Frame //  Initial frame used in the test
		readLength  int    // Number of bytes read before the frame is stripped
		stripLength int    // Number of bytes to be stripped
		result      []byte // Remaining bytes in the frame after it is stripped
		err         string // Error if the strip operation has an error
	}{
		{
			frame:       NewFrame([]byte{0, 1, 2, 3, 4, 5, 6, 7}),
			readLength:  0,
			stripLength: 0,
			result:      []byte{0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			// Verify that the unread frame is unchanged by a failed strip.
			frame:       NewFrame([]byte{0, 1, 2, 3, 4, 5, 6, 7}),
			readLength:  5,
			stripLength: 5,
			err:         "failed",
			result:      []byte{5, 6, 7},
		},
		{
			frame:       NewFrame([]byte{0, 1, 2, 3, 4, 5, 6, 7}),
			readLength:  1,
			stripLength: 7,
			result:      []byte{},
		},
	}

	// Peek into the frame at various offsets and then strip the trailer.
	for _, test := range tests {
		if _, err := test.frame.ReadHeader(test.readLength); err != nil {
			t.Fatalf("Read for length %v failed before stripping trailer %v bytes, err %v.", test.readLength, test.stripLength, err)
		}
		err := test.frame.StripTrailing(test.stripLength)
		if str := ErrTest(err, test.err); str != "" {
			t.Fatalf("Stripping trailer of length %v returned unexpected result, err %v.", test.stripLength, str)
		}

		f, err := test.frame.ReadHeader(test.frame.Len())
		if err != nil {
			t.Fatalf("Read for length %v failed after stripping trailer %v bytes, err %v.", test.readLength, test.stripLength, err)
		}
		if err == nil && !bytes.Equal(f, test.result) {
			t.Fatalf("Read for length %v failed after stripping trailer %v bytes, got %v, want %v.", test.readLength, test.stripLength, f, test.result)
		}
	}
}

// TestResize tests resizing fields of various sizes.
func TestResize(t *testing.T) {
	tests := []struct {
		input     []byte // Input byte slice
		wantSize  int    // wanted size of the byte slice
		wantBytes []byte // wanted byte slice
	}{
		{
			input:     []byte{},
			wantSize:  4,
			wantBytes: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			input:     []byte{0x01},
			wantSize:  4,
			wantBytes: []byte{0x00, 0x00, 0x00, 0x01},
		},
		{
			input:     []byte{0x01, 0x02},
			wantSize:  4,
			wantBytes: []byte{0x00, 0x00, 0x01, 0x02},
		},
		{
			input:     []byte{0x00, 0x00, 0x01, 0x02},
			wantSize:  2,
			wantBytes: []byte{0x01, 0x02},
		},
		{
			input:     []byte{0x00, 0x00, 0x01, 0x02},
			wantSize:  1,
			wantBytes: []byte{0x02},
		},
		{
			input:     []byte{0x00, 0x00, 0x01, 0x02},
			wantSize:  0,
			wantBytes: []byte{},
		},
	}
	for idx, test := range tests {
		got := Resize(test.input, test.wantSize)
		if len(got) != test.wantSize {
			t.Errorf("%v: Incorrect number of bytes on resize of %x, got %v bytes (buffer %x), want %v bytes",
				idx, test.input, len(got), got, test.wantSize)
		}
		if !bytes.Equal(got, test.wantBytes) {
			t.Errorf("%v: Incorrect number of bytes on resize of %x, got buffer %x, want buffer %x",
				idx, test.input, got, test.wantBytes)
		}
	}
}
