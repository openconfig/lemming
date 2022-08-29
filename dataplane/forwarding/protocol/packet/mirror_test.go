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

package packet_test

import (
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// TestMirror tests the mirroring of a packet. The original packet may or
// may not have the input and output port initialized. A mirror is created
// while copying the input and output port from the original packet. The test
// verifies that the input and output port of the mirrored packet are identical
// to the input and output port (if set) of the original packet. Additionally,
// the test updates the mirror and verifies that the original packet is
// unaffected.
func TestMirror(t *testing.T) {
	bytes := []byte{0x01, 0x02, 0x03, 0x04}
	tests := []struct {
		inputPort  fwdobject.NID // input port NID. InvalidNID implies that no input port is set
		outputPort fwdobject.NID // output port NID. InvalidNID implies that no output port is set
	}{
		{
			inputPort:  fwdobject.InvalidNID,
			outputPort: fwdobject.InvalidNID,
		},
		{
			inputPort:  10,
			outputPort: fwdobject.InvalidNID,
		},
		{
			inputPort:  10,
			outputPort: 25,
		},
		{
			inputPort:  fwdobject.InvalidNID,
			outputPort: 25,
		},
	}

	fields := []fwdpacket.FieldID{
		fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_INPUT, 0),
		fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_OUTPUT, 0),
	}
	for idx, test := range tests {
		t.Logf("%v: Running test %+v", idx, test)

		// Create the packet and set its input and output port as specified.
		original, err := protocol.NewPacket(fwdpb.PacketHeaderId_OPAQUE,
			frame.NewFrame(bytes))
		if err != nil {
			t.Fatalf("Unable to create original packet from %x, err %v", bytes, err)
		}
		if test.inputPort != fwdobject.InvalidNID {
			if err := fwdpacket.SetNID(original, test.inputPort, fwdpb.PacketFieldNum_PACKET_PORT_INPUT); err != nil {
				t.Fatalf("Unable to set input port on original, err %v", err)
			}
		}
		if test.outputPort != fwdobject.InvalidNID {
			if err := fwdpacket.SetNID(original, test.outputPort, fwdpb.PacketFieldNum_PACKET_PORT_OUTPUT); err != nil {
				t.Fatalf("Unable to set output port on original, err %v", err)
			}
		}

		// Mirror the packet.
		mirror, err := original.Mirror(fields)
		if err != nil {
			t.Fatalf("Unable to mirror packet %v, err %v", original, err)
		}

		// Query the mirror and ensure that the input and outputs ports are
		// copied from the original packet.
		if test.inputPort != fwdobject.InvalidNID {
			nb, err := mirror.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_INPUT, 0))
			if err != nil {
				t.Fatalf("Unable to query input port for mirror, err %v", err)
			}
			if nid := fwdobject.NID(frame.Field(nb).Value()); nid != test.inputPort {
				t.Fatalf("Unexpected input port for mirror, got %v, want %v", nid, test.inputPort)
			}
		}

		if test.outputPort != fwdobject.InvalidNID {
			nb, err := mirror.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_OUTPUT, 0))
			if err != nil {
				t.Fatalf("Unable to query output port for mirror, err %v", err)
			}
			if nid := fwdobject.NID(frame.Field(nb).Value()); nid != test.outputPort {
				t.Fatalf("Unexpected output port for mirror, got %v, want %v", nid, test.outputPort)
			}
		}

		// Change the mirror and ensure that the original is not changed.
		sampleNID := fwdobject.NID(40)
		if err := fwdpacket.SetNID(mirror, sampleNID, fwdpb.PacketFieldNum_PACKET_PORT_OUTPUT); err != nil {
			t.Fatalf("Unable to update output port on mirror, err %v", err)
		}
		nb, err := mirror.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_OUTPUT, 0))
		if err != nil {
			t.Fatalf("Unable to query updated output port from mirror, err %v", err)
		}
		if nid := fwdobject.NID(frame.Field(nb).Value()); nid != sampleNID {
			t.Fatalf("Unexpected output port from mirror, got %v, want %v", nid, sampleNID)
		}
		if nb, err = original.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_OUTPUT, 0)); err != nil {
			t.Fatalf("Unable to query output port from original, err %v", err)
		}
		if nid := fwdobject.NID(frame.Field(nb).Value()); nid == sampleNID {
			t.Fatalf("Unexpected output port for original, got %v, want != %v", nid, sampleNID)
		}
	}
}
