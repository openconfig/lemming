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

package crc16

import (
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/util/hash/hash16"
)

func TestCRC16(t *testing.T) {
	tests := []struct {
		ansi uint16
		in   string
	}{
		{0x0, ""},
		{0xe8c1, "a"},
		{0x79a8, "ab"},
		{0x9738, "abc"},
		{0x3997, "abcd"},
		{0x85b8, "abcde"},
		{0x5805, "abcdef"},
		{0xe9d9, "abcdefg"},
		{0x7429, "abcdefgh"},
		{0xf075, "abcdefghi"},
		{0xc8b1, "abcdefghij"},
		{0x2ea0, "Discard medicine more than two years old."},
		{0x276b, "He who has a shady past knows that nice guys finish last."},
		{0x1abb, "I wouldn't marry him with a ten foot pole."},
		{0x9499, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
		{0xabfd, "The days of the digital watch are numbered.  -Tom Stoppard"},
		{0x4ee5, "Nepal premier won't resign."},
		{0x761c, "For every action there is an equal and opposite government program."},
		{0xb823, "His money is twice tainted: 'taint yours and 'taint mine."},
		{0xd283, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
		{0x364a, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
		{0x657f, "size:  a.out:  bad magic"},
		{0xe8ec, "The major problem is with sendmail.  -Mark Horton"},
		{0xcb79, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
		{0x3032, "If the enemy is within range, then so are you."},
		{0xc114, "It's well we cannot hear the screams/That we create in others' dreams."},
		{0x161f, "You remind me of a TV show, but that's all right: I watch it anyway."},
		{0x12c6, "C is as portable as Stonehedge!!"},
		{0xc633, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
		{0xf768, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
		{0xbcef, "How can you write a big system without C++?  -Paul Glick"},
	}
	check := func(b []byte, sum hash16.Hash16, want uint16, name string) {
		sum.Write(b)
		if sum.Sum16() != want {
			t.Errorf("%s: unexpected result for %s. Got %x, want %x.", name, b, sum.Sum16(), want)
		}
	}
	for _, test := range tests {
		b := []byte(test.in)
		check(b, NewANSI(), test.ansi, "CRC16-ANSI")
	}
}
