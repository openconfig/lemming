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

// Package crc16 computes the 16 bit checksum over a series of bytes.
package crc16

import "github.com/openconfig/lemming/dataplane/forwarding/util/hash/hash16"

// Size of a CRC-16 checksum in bytes.
const Size = 2

// Predefined polynomials.
const (
	// ANSI is the polynomial x16+x15+x2+1.
	ANSI = 0xA001
)

// Table is a 256-word table representing the polynomial for efficient processing.
type Table [256]uint16

// update returns the result of adding the bytes in p to the crc.
func (t *Table) update(crc uint16, p []byte) uint16 {
	for _, v := range p {
		crc = (crc >> 8) ^ t[(crc^uint16(v))&0xff]
	}
	return crc
}

// ANSITable is the table for the ANSI polynomial.
var ANSITable = makeTable(ANSI)

// makeTable constructs a table for the specified CRC16 polynomial.
func makeTable(poly uint16) *Table {
	t := new(Table)
	for i := 0; i < 256; i++ {
		crc := uint16(i)
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ poly
			} else {
				crc >>= 1
			}
		}
		t[i] = crc
	}
	return t
}

// MakeTable returns the Table constructed from the specified polynomial.
func MakeTable(poly uint16) *Table {
	switch poly {
	case ANSI:
		return ANSITable
	default:
		return makeTable(poly)
	}
}

// digest represents the partial evaluation of a checksum.
type digest struct {
	crc uint16
	tab *Table
}

// NewANSI creates a new hash.Hash16 that uses the CRC16-ANSI polynomial.
func NewANSI() hash16.Hash16 { return &digest{0, ANSITable} }

// Size returns the number of bytes Sum will return.
func (d *digest) Size() int { return Size }

// BlockSize returns the hash's underlying block size in bytes.
func (d *digest) BlockSize() int { return 1 }

// Reset resets the hash to its initial state.
func (d *digest) Reset() { d.crc = 0 }

// Write adds more data to the running checksum.
func (d *digest) Write(p []byte) (n int, err error) {
	d.crc = d.tab.update(d.crc, p)
	return len(p), nil
}

// Sum16 returns the current value of the checksum.
func (d *digest) Sum16() uint16 { return d.crc }

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (d *digest) Sum(in []byte) []byte {
	s := d.Sum16()
	return append(in, byte(s>>8), byte(s))
}

// ChecksumANSI returns the checksum of data using the CRC16-ANSI polynomial.
func ChecksumANSI(data []byte) uint16 { return ANSITable.update(0, data) }
