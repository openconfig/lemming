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

package vxlan

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ip"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/tcp"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestVXLANParseAndBuild(t *testing.T) {
	// VXLAN header: 8 bytes. Flags (1 byte, iFlag = 0x08), Reserved (2 bytes), VNI (3 bytes = 0x000123), Reserved (2 bytes).
	raw := []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00}
	f := frame.NewFrame(raw)
	handler, next, err := parse(f, &protocol.Desc{})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if next != fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET {
		t.Errorf("got next header %v, want ETHERNET", next)
	}

	// Test Field get for VNI
	vniFieldID := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI, 0)
	val, err := handler.Field(vniFieldID)
	if err != nil {
		t.Fatalf("Field VNI failed: %v", err)
	}
	wantVNI := []byte{0x00, 0x12, 0x34}
	if diff := cmp.Diff(val, wantVNI); diff != "" {
		t.Errorf("VNI field diff (-got +want):\n%s", diff)
	}

	// Test UpdateField (set VNI)
	newVNI := []byte{0x00, 0xab, 0xcd}
	updated, err := handler.UpdateField(vniFieldID, fwdpacket.OpSet, newVNI)
	if err != nil || !updated {
		t.Fatalf("UpdateField VNI failed: %v, updated=%v", err, updated)
	}

	val, err = handler.Field(vniFieldID)
	if err != nil {
		t.Fatalf("Field VNI after update failed: %v", err)
	}
	if diff := cmp.Diff(val, newVNI); diff != "" {
		t.Errorf("Updated VNI field diff (-got +want):\n%s", diff)
	}

	// Test build / add
	h2, err := add(fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN, &protocol.Desc{})
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if h2.Header() == nil {
		t.Errorf("added header is nil")
	}
}

func TestVXLANParseTruncated(t *testing.T) {
	for _, size := range []int{0, 4, 7} {
		raw := make([]byte, size)
		f := frame.NewFrame(raw)
		_, _, err := parse(f, &protocol.Desc{})
		if err == nil {
			t.Errorf("parse() with size %d expected error, got nil", size)
		}
	}
}

func TestVXLANUnsupportedField(t *testing.T) {
	raw := []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00}
	f := frame.NewFrame(raw)
	handler, _, err := parse(f, &protocol.Desc{})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	badFieldID := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_UNSPECIFIED, 0)
	_, err = handler.Field(badFieldID)
	if err == nil {
		t.Errorf("Field() with unsupported field expected error, got nil")
	}
}

func TestVXLANUpdateFieldErrors(t *testing.T) {
	raw := []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00}
	f := frame.NewFrame(raw)
	handler, _, err := parse(f, &protocol.Desc{})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	vniFieldID := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI, 0)
	badFieldID := fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_UNSPECIFIED, 0)

	// Test unsupported op
	_, err = handler.UpdateField(vniFieldID, fwdpacket.OpInc, []byte{1})
	if err == nil {
		t.Errorf("UpdateField() with unsupported op expected error, got nil")
	}

	// Test unsupported field
	_, err = handler.UpdateField(badFieldID, fwdpacket.OpSet, []byte{1})
	if err == nil {
		t.Errorf("UpdateField() with unsupported field expected error, got nil")
	}
}

func TestVXLANModify(t *testing.T) {
	raw := []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00}
	f := frame.NewFrame(raw)
	handler, _, err := parse(f, &protocol.Desc{})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	err = handler.Modify(fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN)
	if err == nil {
		t.Errorf("Modify() expected error, got nil")
	}
}

func TestVXLANRemove(t *testing.T) {
	raw := []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00}
	f := frame.NewFrame(raw)
	handler, _, err := parse(f, &protocol.Desc{})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	// Non-matching header ID
	err = handler.Remove(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET)
	if err == nil {
		t.Errorf("Remove() with non-matching header ID expected error, got nil")
	}

	// Matching header ID
	err = handler.Remove(fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN)
	if err != nil {
		t.Errorf("Remove() with matching header ID failed: %v", err)
	}
}

func TestVXLANEntropyCalculation(t *testing.T) {
	vxlanHdr := []byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00}
	innerEth1 := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x08, 0x00}
	innerEth2 := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x08, 0x00}

	innerIPAndTCP := []byte{
		0x45, 0x00, 0x00, 0x50, 0x00, 0x00, 0x00, 0x00, 0x40, 0x06, 0x00, 0x00,
		10, 0, 0, 1,
		10, 0, 0, 2,
		0x1f, 0x90, 0x00, 0x50, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x50, 0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	data1 := []byte("HTTP payload 1")
	data2 := []byte("HTTP payload 2 different")

	flow1Pkt1 := slices.Concat(vxlanHdr, innerEth1, innerIPAndTCP, data1)
	flow1Pkt2 := slices.Concat(vxlanHdr, innerEth1, innerIPAndTCP, data2)
	flow2Pkt := slices.Concat(vxlanHdr, innerEth2, innerIPAndTCP, data1)

	pkt1, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN, flow1Pkt1)
	if err != nil {
		t.Fatalf("fwdpacket.New failed: %v", err)
	}
	pkt2, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN, flow1Pkt2)
	if err != nil {
		t.Fatalf("fwdpacket.New failed: %v", err)
	}
	pkt3, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN, flow2Pkt)
	if err != nil {
		t.Fatalf("fwdpacket.New failed: %v", err)
	}

	port1 := ComputeEntropySourcePort(&protocol.Desc{Packet: pkt1.(*protocol.Packet)})
	port2 := ComputeEntropySourcePort(&protocol.Desc{Packet: pkt2.(*protocol.Packet)})
	port3 := ComputeEntropySourcePort(&protocol.Desc{Packet: pkt3.(*protocol.Packet)})

	if port1 < EphemeralPortBase {
		t.Errorf("port1 = %d, want >= %d", port1, EphemeralPortBase)
	}
	if port1 != port2 {
		t.Errorf("expected identical port for same flow with different payload: got %d vs %d", port1, port2)
	}
	if port1 == port3 {
		t.Errorf("expected different ports for different flows: got %d for both", port1)
	}
	if got := ComputeEntropySourcePort(nil); got != 0 {
		t.Errorf("ComputeEntropySourcePort(nil) = %d, want 0", got)
	}
}

