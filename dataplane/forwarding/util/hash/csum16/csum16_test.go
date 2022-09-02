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

package csum16

import (
	"testing"
)

// TestChecksumFixed tests the checksum computation for some known set of bytes.
func TestChecksumFixed(t *testing.T) {
	tests := []struct {
		buffer []byte
		result uint16
	}{
		{
			nil,
			0x0,
		},
		{
			[]byte{0x01},
			0xfeff,
		},
		{
			[]byte{0x45, 0x00, 0x00, 0x1c, 0x03, 0xde, 0x00, 0x00, 0x40, 0x01, 0x00, 0x00, 0x7f, 0x00, 0x00, 0x01, 0x7f, 0x00, 0x00, 0x01},
			0x7901,
		},
		{
			[]byte{0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x06, 0x00, 0x00, 0xac, 0x10, 0x0a, 0x63, 0xac, 0x10, 0x0a, 0x0c},
			0xb1e6,
		},
		{
			[]byte{0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x06, 0x00, 0x00, 0xac, 0x10, 0x0a, 0x63, 0xac, 0x10, 0x0a, 0x0c, 0x10},
			0xa1e6,
		},
		{
			[]byte{0x45, 0x00, 0x00, 0x40, 0xf9, 0xa6, 0x40, 0x00, 0x40, 0x06, 0x00, 0x00, 0x0a, 0xc5, 0x0a, 0x04, 0x0a, 0xc5, 0x0a, 0x05},
			0x177f,
		},
	}
	for _, test := range tests {
		var sum Sum
		sum.Write(test.buffer)
		if csum := sum.Sum16(); csum != test.result {
			t.Errorf("Checksum failed for buffer %v, got %x want %x", test.buffer, csum, test.result)
		}
	}
}

// TestChecksumBuffer tests the checksum computation for complete buffers of
// varying sizes. The test validates the checksum as described by RFC-1624.
//
// For convenience, each buffer is assumed to have the first two bytes reserved
// for the checksum. Initially the checksum bytes in the buffer are set to zero
// and the checksum is computed. The computed checksum is then set in the buffer
// and the checksum is recomputed. If the recomputed checksum is zero, it
// indicates that the computed checksum was correct.
func TestChecksumBuffer(t *testing.T) {
	// Arbitrary slice of bytes.
	source := []byte{0x00, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74}
	for size := 2; size < len(source); size++ {
		buffer := source[:size]

		buffer[0] = 0
		buffer[1] = 0
		var sum Sum
		sum.Write(buffer)
		csum1 := sum.Sum16()

		buffer[0] = byte(csum1 >> 8)
		buffer[1] = byte(csum1 & 0xff)
		sum.Reset()
		sum.Write(buffer)
		if csum := sum.Sum16(); csum != 0 {
			t.Errorf("Checksum failed for buffer %v, got %x want 0 (checksum %x)", buffer, csum, csum1)
		}
	}
}

// TestChecksumBufferSplit tests the checksum computation for split buffers of
// varying sizes. The test validates the checksum as described by RFC-1624 and
// RFC-1071.
//
// For convenience, each buffer is assumed to have a minimum size of 4 and the
// first two bytes reserved for the checksum. Initially the checksum bytes in
// the buffer are set to zero and the checksum is computed over two halves.
// The computed checksum is then set in the buffer and the checksum is
// recomputed over the same two halves. If the recomputed checksum is zero, it
// indicates that the computed checksum was correct.
func TestChecksumBufferSplit(t *testing.T) {
	// Arbitrary slice of bytes.
	source := []byte{0x00, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74}
	for size := 4; size < len(source); size++ {
		buffer := source[:size]
		step := size / 2

		buffer[0] = 0
		buffer[1] = 0
		var sum Sum
		sum.Write(buffer[0:step])
		sum.Write(buffer[step:])
		csum1 := sum.Sum16()

		buffer[0] = byte(csum1 >> 8)
		buffer[1] = byte(csum1 & 0xff)
		sum.Reset()
		sum.Write(buffer[0:step])
		sum.Write(buffer[step:])
		if csum := sum.Sum16(); csum != 0 {
			t.Errorf("Checksum failed for buffer %v, got %x want 0 (checksum %x)", buffer, csum, csum1)
		}
	}
}
