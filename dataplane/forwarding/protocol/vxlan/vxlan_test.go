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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
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

