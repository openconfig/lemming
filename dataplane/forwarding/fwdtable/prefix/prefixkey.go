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
)

const (
	byteBitCount = 8 // number of bits stored in a byte
	byteLogCount = 3 // log to the base 2 of the byteBitCount
)

// byteCount returns the smallest number of bytes needed to store the
// specified number of bits.
func byteCount(bitCount int) int {
	if bitCount <= byteBitCount {
		return 1
	}
	return (bitCount + byteBitCount - 1) >> byteLogCount
}

// wordset is a slice of bytes used to store a set of bits. It is assumed that
// the wordset is only used for the key and all operations are valid.
type wordset []byte

// bit returns the bit at the specified position.
func (ws wordset) bit(bitpos int) byte {
	idx := bitpos >> byteLogCount
	pos := byteBitCount - 1 - (bitpos & (byteBitCount - 1))
	return (ws[idx] >> uint(pos)) & 0x1
}

// set sets the bit at the specified position.
func (ws wordset) set(bitpos int, v byte) {
	idx := bitpos >> byteLogCount
	pos := byteBitCount - 1 - (bitpos & (byteBitCount - 1))
	mask := byte(1 << uint(pos))
	if v == 0x1 {
		ws[idx] |= mask
	} else {
		ws[idx] &= ^mask
	}
}

// ones returns the number of ones in a wordset.
func (ws wordset) ones() int {
	cmap := map[byte]int{
		0x0: 0,
		0x1: 1,
		0x2: 1,
		0x3: 2,
		0x4: 1,
		0x5: 2,
		0x6: 2,
		0x7: 3,
		0x8: 1,
		0x9: 2,
		0xa: 2,
		0xb: 3,
		0xc: 2,
		0xd: 3,
		0xe: 3,
		0xf: 4,
	}
	c := 0
	for _, w := range ws {
		c += cmap[w&0x0F] + cmap[(w>>4)&0x0F]
	}
	return c
}

// A key is a prefix table key specified as a set of bits stored in a slice of
// bytes. It is assumed the key is used only for the prefix table, and all
// operations are valid.
type key struct {
	bytes    wordset // bytes containing the bits
	bitCount int     // number of bits in the set
	bitStart int     // starting bit position
}

// newKey creates a new key from a slice of bytes containing the specified
// number of bits.
func newKey(value []byte, bitCount int) *key {
	byteCount := byteCount(bitCount)
	if byteCount < len(value) {
		byteCount = len(value)
	}
	bytes := make([]byte, byteCount)
	copy(bytes, value)
	return &key{
		bytes:    bytes,
		bitCount: bitCount,
		bitStart: 0,
	}
}

// HasPrefix returns true if 'q' is a prefix of 's'.
func (s *key) HasPrefix(q *key) bool {
	if s.bitCount < q.bitCount {
		return false
	}
	for pos := 0; pos < q.bitCount; pos++ {
		if s.Bit(pos) != q.Bit(pos) {
			return false
		}
	}
	return true
}

// Pack packs the bit set by removing any unused bits.
func (s *key) Pack() {
	bytes := wordset(make([]byte, byteCount(s.bitCount)))
	for pos := 0; pos < s.bitCount; pos++ {
		bytes.set(pos, s.Bit(pos))
	}
	s.bitStart = 0
	s.bytes = bytes
}

// String returns the set formatted as a string after it is packed.
func (s *key) String() string {
	s.Pack()
	return fmt.Sprintf("%x/%d", s.bytes, s.bitCount)
}

// IsEqual returns true if the two keys are equal.
func (s *key) IsEqual(q *key) bool {
	if s.bitCount != q.bitCount {
		return false
	}
	for pos := 0; pos < s.bitCount; pos++ {
		if s.Bit(pos) != q.Bit(pos) {
			return false
		}
	}
	return true
}

// Bit gets the bit at the specified position.
func (s *key) Bit(bitpos int) byte {
	return s.bytes.bit(s.bitStart + bitpos)
}

// Set sets the bit at the specified index to true or false.
func (s *key) Set(bitpos int, v byte) {
	s.bytes.set(s.bitStart+bitpos, v)
}

// TrimPrefix removes a key which is known to be a prefix and returns a new key.
func (s *key) TrimPrefix(q *key) *key {
	return &key{
		bytes:    s.bytes,
		bitCount: s.bitCount - q.bitCount,
		bitStart: s.bitStart + q.bitCount,
	}
}

// BitCount returns the number of bits in the key.
func (s *key) BitCount() int {
	return s.bitCount
}

// Copy creates a copy of the key.
func (s *key) Copy() *key {
	ws := wordset(make([]byte, len(s.bytes)))
	copy(ws, s.bytes)
	return &key{
		bytes:    ws,
		bitCount: s.bitCount,
		bitStart: s.bitStart,
	}
}

// combine creates a new key by appending s2 to s1.
func combine(s1, s2 *key) *key {
	bitCount := s1.bitCount + s2.bitCount
	s := &key{
		bytes:    make([]byte, byteCount(bitCount)),
		bitCount: bitCount,
		bitStart: 0,
	}
	for pos := 0; pos < s1.bitCount; pos++ {
		s.Set(pos, s1.Bit(pos))
	}
	for pos := 0; pos < s2.bitCount; pos++ {
		s.Set(s1.bitCount+pos, s2.Bit(pos))
	}
	return s
}

// prefixKey finds the common prefixKey between two keys.
func prefixKey(s1, s2 *key) *key {
	max := s1.bitCount
	if s2.bitCount < max {
		max = s2.bitCount
	}
	bytes := wordset(make([]byte, byteCount(max)))
	var count int
	for count = 0; count < max; count++ {
		v := s1.Bit(count)
		if v != s2.Bit(count) {
			break
		}
		bytes.set(count, v)
	}
	prefix := &key{
		bytes:    bytes,
		bitCount: count,
		bitStart: 0,
	}
	return prefix
}

// Calculate calculates the number of bits in the specified number of bytes.
func Calculate(byteCount int) int {
	return byteBitCount * byteCount
}

// Count counts the number of bits set in a slice of bytes.
func Count(bs []byte) int {
	ws := wordset(bs)
	return ws.ones()
}
