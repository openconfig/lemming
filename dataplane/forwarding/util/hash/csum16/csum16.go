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

// Package csum16 computes the 16 bit Internet checksum.
package csum16

// The size of a checksum in bytes.
const size = 2

// Sum stores the partially computed checksum.
type Sum uint16

// Size returns the number of bytes Sum will return.
func (Sum) Size() int { return size }

// BlockSize returns the hash's underlying block size.
func (Sum) BlockSize() int { return 1 }

// Reset resets the hash to its initial state.
func (s *Sum) Reset() { *s = 0 }

// Sum16 returns the current value of the checksum.
func (s Sum) Sum16() uint16 { return uint16(s) }

// Sum appends the current hash to b and returns the resulting slice.
func (s Sum) Sum(b []byte) []byte {
	v := s.Sum16()
	return append(b, byte(v>>8), byte(v))
}

// Write computes the checksum over a series of bytes.
func (s *Sum) Write(p []byte) (n int, err error) {
	// Accumulate the bytes in a uint32 to accumulate the carry.
	sum := uint32(*s ^ 0xffff)
	for pos, b := range p {
		curr := uint32(b)
		if pos&0x1 == 0 {
			curr <<= 8
		}
		sum += curr
	}

	// Fold the sum and return the one's complement.
	for sum&0xffff0000 != 0 {
		sum = sum&0xffff + sum>>16
	}
	*s = Sum(sum ^ 0xffff)
	return len(p), nil
}
