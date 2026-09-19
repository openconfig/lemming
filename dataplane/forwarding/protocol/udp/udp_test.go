// Copyright 2024 Google LLC
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

package udp

import (
	"encoding/binary"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/vxlan"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestUDPParse(t *testing.T) {
	// Standard UDP packet: sport = 1234 (0x04d2), dport = 80 (0x0050), len = 12, csum = 0
	rawStd := []byte{0x04, 0xd2, 0x00, 0x50, 0x00, 0x0c, 0x00, 0x00, 0xde, 0xad, 0xbe, 0xef}
	fStd := frame.NewFrame(rawStd)
	handler, next, err := parse(fStd, &protocol.Desc{})
	if err != nil {
		t.Fatalf("parse standard UDP failed: %v", err)
	}
	if next != fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE {
		t.Errorf("got next header %v, want OPAQUE", next)
	}
	if len(handler.Header()) != udpBytes {
		t.Errorf("got header len %d, want %d", len(handler.Header()), udpBytes)
	}
	if fStd.Len() != 4 {
		t.Errorf("got frame len %d, want 4 payload bytes", fStd.Len())
	}

	// VXLAN UDP packet: sport = 1234 (0x04d2), dport = 4789 (vxlan.DefaultPort), len = 20, csum = 0
	dportBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dportBytes, vxlan.DefaultPort)
	rawVxlan := append([]byte{0x04, 0xd2}, dportBytes...)
	rawVxlan = append(rawVxlan, []byte{0x00, 0x14, 0x00, 0x00}...)
	rawVxlan = append(rawVxlan, []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00, 0xde, 0xad, 0xbe, 0xef}...)

	fVxlan := frame.NewFrame(rawVxlan)
	handlerVxlan, nextVxlan, err := parse(fVxlan, &protocol.Desc{})
	if err != nil {
		t.Fatalf("parse VXLAN UDP failed: %v", err)
	}
	if nextVxlan != fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET {
		t.Errorf("got next header %v, want ETHERNET", nextVxlan)
	}
	if len(handlerVxlan.Header()) != udpBytes+8 {
		t.Errorf("got header len %d, want %d", len(handlerVxlan.Header()), udpBytes+8)
	}
	if fVxlan.Len() != 4 {
		t.Errorf("got frame len %d, want 4 payload bytes", fVxlan.Len())
	}
}

func TestUDPParseTruncated(t *testing.T) {
	for _, size := range []int{0, 2, 7} {
		raw := make([]byte, size)
		f := frame.NewFrame(raw)
		_, _, err := parse(f, &protocol.Desc{})
		if err == nil {
			t.Errorf("parse() with size %d expected error, got nil", size)
		}
	}
}

func TestUDPFieldsAndUpdates(t *testing.T) {
	raw := []byte{0x04, 0xd2, 0x00, 0x50, 0x00, 0x08, 0x00, 0x00}
	f := frame.NewFrame(raw)
	handler, _, err := parse(f, &protocol.Desc{})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	srcFieldID := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC, 0)
	dstFieldID := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST, 0)

	// Test Field get
	srcVal, err := handler.Field(srcFieldID)
	if err != nil {
		t.Fatalf("Field SRC failed: %v", err)
	}
	if diff := cmp.Diff(srcVal, []byte{0x04, 0xd2}); diff != "" {
		t.Errorf("SRC field diff (-got +want):\n%s", diff)
	}

	dstVal, err := handler.Field(dstFieldID)
	if err != nil {
		t.Fatalf("Field DST failed: %v", err)
	}
	if diff := cmp.Diff(dstVal, []byte{0x00, 0x50}); diff != "" {
		t.Errorf("DST field diff (-got +want):\n%s", diff)
	}

	// Test UpdateField
	updated, err := handler.UpdateField(srcFieldID, fwdpacket.OpSet, []byte{0x11, 0x22})
	if err != nil || !updated {
		t.Fatalf("UpdateField SRC failed: %v, updated=%v", err, updated)
	}
	srcVal, err = handler.Field(srcFieldID)
	if err != nil {
		t.Fatalf("Field SRC after update failed: %v", err)
	}
	if diff := cmp.Diff(srcVal, []byte{0x11, 0x22}); diff != "" {
		t.Errorf("Updated SRC field diff (-got +want):\n%s", diff)
	}

	// Test unsupported operation or field
	_, err = handler.UpdateField(srcFieldID, fwdpacket.OpInc, []byte{1})
	if err == nil {
		t.Errorf("UpdateField with unsupported op expected error, got nil")
	}

	badFieldID := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_UNSPECIFIED, 0)
	_, err = handler.Field(badFieldID)
	if err == nil {
		t.Errorf("Field with unsupported field expected error, got nil")
	}
	_, err = handler.UpdateField(badFieldID, fwdpacket.OpSet, []byte{1})
	if err == nil {
		t.Errorf("UpdateField with unsupported field expected error, got nil")
	}
}

func TestUDPAdd(t *testing.T) {
	h, err := add(fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP, &protocol.Desc{})
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if len(h.Header()) != udpBytes {
		t.Errorf("added header len %d, want %d", len(h.Header()), udpBytes)
	}
}

func TestUDPMethods(t *testing.T) {
	h, _ := add(fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP, &protocol.Desc{})
	if h.Trailer() != nil {
		t.Errorf("Trailer() != nil")
	}
	if h.ID(0) != fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP {
		t.Errorf("ID() = %v, want UDP", h.ID(0))
	}
	if udpUDP, ok := h.(*UDP); ok {
		if err := udpUDP.Remove(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET); err == nil {
			t.Errorf("Remove with non-matching header expected error, got nil")
		}
		if err := udpUDP.Remove(fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP); err != nil {
			t.Errorf("Remove failed: %v", err)
		}
		if err := udpUDP.Modify(fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP); err == nil {
			t.Errorf("Modify expected error, got nil")
		}
	}
}
