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

// Package packettestutil contains a set of routines used to test the processing
// of packet headers and fields.
package packettestutil

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// ErrTest returns a non-empty string if err and str do not match.
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

// frame builds a packet frame from a series of headers, each described as
// a slice of bytes.
func makeFrame(headers [][]byte) []byte {
	var f []byte
	for _, h := range headers {
		f = append(f, h...)
	}
	return f
}

// A FieldUpdate describes an update to a packet field.
type FieldUpdate struct {
	ID  fwdpacket.FieldID // Field identifier
	Arg []byte            // Value to use for the update
	Op  int               // The type of update described by fwdpacket.OpXyz
	Err string            // Expected error
}

// FieldUpdates applies a series of updates and checks its succeess or failure.
func FieldUpdates(t *testing.T, name string, id int, packet fwdpacket.Packet, updates []FieldUpdate) {
	for upos, update := range updates {
		var initial uint
		switch update.Op {
		case fwdpacket.OpInc, fwdpacket.OpDec:
			if result, err := packet.Field(update.ID); err == nil {
				initial = frame.Field(result).Value()
			}
		}

		err := packet.Update(update.ID, update.Op, update.Arg)
		if str := ErrTest(err, update.Err); str != "" {
			t.Errorf("%v %d.%d: field %+v update failed: %v", name, id, upos, update.ID, str)
			return
		}
		if err != nil {
			continue
		}
		result, err := packet.Field(update.ID)
		if err != nil {
			t.Errorf("%v %d.%d: field %+v query after update failed: %v", name, id, upos, update.ID, err)
			continue
		}
		switch update.Op {
		case fwdpacket.OpSet:
			if !bytes.Equal(result, update.Arg) {
				t.Errorf("%v %d.%d: field %+v failed after set: got %x, want %x.", name, id, upos, update.ID, result, update.Arg)
			}
		case fwdpacket.OpInc:
			arg := frame.Field(update.Arg).Value()
			value := frame.Field(result).Value()
			if value != (initial + arg) {
				t.Errorf("%v %d.%d: field %+v failed after inc: got %x, want %x.", name, id, upos, update.ID, value, (initial + arg))
			}
		case fwdpacket.OpDec:
			arg := frame.Field(update.Arg).Value()
			value := frame.Field(result).Value()
			if value != (initial - arg) {
				t.Errorf("%v %d.%d: field %+v failed after dec: got %x, want %x.", name, id, upos, update.ID, value, (initial - arg))
			}
		}
	}
}

// A FieldQuery describes a query of a packet's field and its expected result.
type FieldQuery struct {
	ID     fwdpacket.FieldID // Field identifier
	Result []byte            // Expected value of the field
	Err    string            // Expected error
}

// FieldQueries performs a series of queries to check their succeess or failure.
func FieldQueries(t *testing.T, name string, id int, packet fwdpacket.Packet, queries []FieldQuery) {
	for qpos, query := range queries {
		result, err := packet.Field(query.ID)
		if str := ErrTest(err, query.Err); str != "" {
			t.Errorf("%v %d.%d: field %+v query failed: %s", name, id, qpos, query.ID, str)
			continue
		}
		if !bytes.Equal(result, query.Result) {
			t.Errorf("%v %d.%ds: field %+v: got %x, want %x.", name, id, qpos, query.ID, result, query.Result)
		}
	}
}

// A PacketFieldTest describes a set of tests on a packet's fields.
// Each test describes the following:
// - An original packet frame as a series of byte slices.
// - A set of queries performed on the original frame.
// - A set of updates performed on the original frame.
// - A resultant frame as a series of byte slices.
type PacketFieldTest struct {
	StartHeader fwdpb.PacketHeaderId // First packet header in the frame
	Orig        [][]byte             // Set of headers in the original frame
	Final       [][]byte             // Set of headers in the final frame
	Queries     []FieldQuery
	Updates     []FieldUpdate
}

// TestPacketFields performs a series of packet field tests.
func TestPacketFields(name string, t *testing.T, tests []PacketFieldTest) {
	for pos, test := range tests {
		orig := makeFrame(test.Orig)
		packet, err := fwdpacket.New(test.StartHeader, orig)
		if err != nil {
			t.Fatalf("%v %d: Unable to create packet from frame %x, err %v.", name, pos, orig, err)
		}
		t.Logf("%v %d: Created packet %v from frame %x", name, pos, packet, orig)
		if packet.Length() != len(orig) {
			t.Errorf("%v %d: packet length: got %v, want %v.", name, pos, packet.Length(), len(orig))
		}
		FieldQueries(t, name, pos, packet, test.Queries)
		f := packet.Frame()
		if !bytes.Equal(orig, f) {
			t.Errorf("%v %d: rebuilt frame: got %x, want %x.", name, pos, f, orig)
		}
		FieldUpdates(t, name, pos, packet, test.Updates)
		if final := makeFrame(test.Final); final != nil {
			f = packet.Frame()
			if !bytes.Equal(f, final) {
				t.Errorf("%v %d: updated rebuilt frame: got %x, want %x.", name, pos, f, final)
			}
		}
	}
}

// A HeaderUpdate describes an update to the packet header via
// an encap or decap operation.
type HeaderUpdate struct {
	ID      fwdpb.PacketHeaderId // Header being operated
	Encap   bool                 // Encap if true, decap otherwise
	Err     string               // Expected error
	Updates []FieldUpdate        // Series of updates applied after an encap
	Result  [][]byte             // Expected frame after the update
}

// A PacketHeaderTest describes a set of tests on a packet's headers.
// Each test describes the following:
// - An original packet frame as a series of byte slices.
// - A set of header operations performed on the original frame.
type PacketHeaderTest struct {
	StartHeader fwdpb.PacketHeaderId // First packet header in the frame
	Orig        [][]byte             // Original frame
	Updates     []HeaderUpdate       // Series of operations on the header
}

// TestPacketHeaders performs a series of packet header tests.
func TestPacketHeaders(name string, t *testing.T, tests []PacketHeaderTest) {

	for pos, test := range tests {
		orig := makeFrame(test.Orig)
		packet, err := fwdpacket.New(test.StartHeader, orig)
		if err != nil {
			t.Fatalf("%v %d: Unable to create packet from frame %x, err %v.", name, pos, orig, err)
		}
		t.Logf("%v %d: Created packet %v from frame %x", name, pos, packet, orig)
		if packet.Length() != len(orig) {
			t.Errorf("%v %d: packet length: got %v, want %v.", name, pos, packet.Length(), len(orig))
		}
		for upos, update := range test.Updates {
			t.Logf("%v %d.%d: Current packet %v frame %x", name, pos, upos, packet, packet.Frame())
			var err error
			var desc string
			if update.Encap {
				desc = "encap"
				err = packet.Encap(update.ID)
			} else {
				desc = "decap"
				err = packet.Decap(update.ID)
			}
			if str := ErrTest(err, update.Err); str != "" {
				t.Errorf("%v %d.%d: %v failed for header %+v, failed %v", name, pos, upos, desc, update.ID, str)
				continue
			}
			FieldUpdates(t, name, pos, packet, update.Updates)
			if result := makeFrame(update.Result); result != nil {
				f := packet.Frame()
				if err == nil && !bytes.Equal(f, result) {
					t.Errorf("%v %d.%d: %v has incorrect result for header %+v, got %x want %x.", name, pos, upos, desc, update.ID, f, result)
				}
			}
			t.Logf("%v %d.%d: New packet %v frame %x", name, pos, upos, packet, packet.Frame())
		}
		t.Logf("%v %d: Final packet %v, frame %x", name, pos, packet, packet.Frame())
	}
}
