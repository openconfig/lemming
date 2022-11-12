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

package actions

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// The update operations SET / INC / DEC are convered by the
// corresponding protocol handler tests, as the actual operation occurs in
// the corresponding protocol handler. The copy and bit operations are unit
// tested here, as the operation happens within the action's process function.

// TestCopy tests the copy update action.
func TestCopy(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a context with a test port.
	ctx := fwdcontext.New("test", "fwd")

	// Source and destination fields.
	srcField := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC,
		},
	}

	dstField :=
		&fwdpb.PacketFieldId{
			Field: &fwdpb.PacketField{
				FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
			},
		}

	// Create a copy update action for a packet field of the same size.
	// This is expected to succeed.
	desc := fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
	}
	update := fwdpb.UpdateActionDesc{
		FieldId: dstField,
		Type:    fwdpb.UpdateType_UPDATE_TYPE_COPY,
		Field:   srcField,
	}
	desc.Action = &fwdpb.ActionDesc_Update{
		Update: &update,
	}
	action, err := fwdaction.New(&desc, ctx)
	if err != nil {
		t.Errorf("NewAction failed, desc %v failed, err %v.", &desc, err)
	}

	// Verify the action by processing a packet where the source field exists.
	fieldBytes := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04}
	packet1 := mock_fwdpacket.NewMockPacket(ctrl)
	packet1.EXPECT().Logf(gomock.Any(), gomock.Any()).AnyTimes()
	packet1.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	packet1.EXPECT().Field(fwdpacket.NewFieldID(srcField)).Return(fieldBytes, nil)
	packet1.EXPECT().Update(fwdpacket.NewFieldID(dstField), fwdpacket.OpSet, fieldBytes)

	if _, state := action.Process(packet1, nil); state != fwdaction.CONTINUE {
		t.Errorf("%v processing returned bad result. Got %v want %v.", action, state, fwdaction.CONTINUE)
	}

	// Verify the action by processing a packet where the source field does not
	// exist.
	packet2 := mock_fwdpacket.NewMockPacket(ctrl)
	packet2.EXPECT().Logf(gomock.Any(), gomock.Any()).AnyTimes()
	packet2.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	packet2.EXPECT().Field(fwdpacket.NewFieldID(srcField)).Return(nil, errors.New("No field"))
	packet2.EXPECT().String().Return("").AnyTimes()

	if _, state := action.Process(packet2, nil); state != fwdaction.DROP {
		t.Errorf("%v processing returned bad result. Got %v want %v.", action, state, fwdaction.DROP)
	}
}

// TestBitWrite tests the bit update action.
func TestBitWrite(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a context with a test port.
	ctx := fwdcontext.New("test", "fwd")

	// Packet field
	field := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC,
		},
	}

	tests := []struct {
		original   []byte // original bytes of the field.
		bitCount   uint32 // number of bits to modify
		bitOffset  uint32 // offset in bits to modify
		value      []byte // bits to modify
		final      []byte // final bytes of the field
		buildErr   bool   // true if a build error is expected
		processErr bool   // true if a packet processing error is expected
	}{
		// Test a build failure as the bit count is higher than the value.
		{
			bitCount: 1024,
			value:    []byte{0x01, 0x02, 0x03},
			buildErr: true,
		},
		// Test a processing failure as the bit count is higher than the original packet field size.
		{
			original:   []byte{0x01, 0x02, 0x03, 0x04},
			bitCount:   33,
			bitOffset:  0,
			value:      []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			processErr: true,
		},
		// Test a processing failure as the bit offset is higher than the original packet field size.
		{
			original:   []byte{0x01, 0x02, 0x03, 0x04},
			bitCount:   1,
			bitOffset:  64,
			value:      []byte{0x01, 0x02, 0x03},
			processErr: true,
		},
		// Test a processing failure as the bit offset + bit count is higher than the original packet field size.
		{
			original:   []byte{0x01, 0x02, 0x03, 0x04},
			bitCount:   2,
			bitOffset:  31,
			value:      []byte{0x01},
			processErr: true,
		},
		// Test a successful case with a single bit.
		{
			original:  []byte{0x40, 0x00, 0x00, 0x01},
			bitCount:  1,
			bitOffset: 1,
			value:     []byte{0x01},
			final:     []byte{0x40, 0x00, 0x00, 0x03},
		},
		// Test a successful case with a single bit.
		{
			original:  []byte{0x80, 0x00, 0x00, 0x01},
			bitCount:  1,
			bitOffset: 30,
			value:     []byte{0x01},
			final:     []byte{0xC0, 0x00, 0x00, 0x01},
		},
		// Test a successful case with a single bit.
		{
			original:  []byte{0x40, 0x00, 0x00, 0x01},
			bitCount:  2,
			bitOffset: 30,
			value:     []byte{0x03},
			final:     []byte{0xC0, 0x00, 0x00, 0x01},
		},
	}

	for tid, test := range tests {
		t.Logf("%d: Running test %+v", tid, test)

		// Create an update action and check for build errors.
		desc := fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
		}
		update := fwdpb.UpdateActionDesc{
			FieldId:   field,
			Type:      fwdpb.UpdateType_UPDATE_TYPE_BIT_WRITE,
			BitCount:  test.bitCount,
			BitOffset: test.bitOffset,
			Value:     test.value,
		}
		desc.Action = &fwdpb.ActionDesc_Update{
			Update: &update,
		}
		action, err := fwdaction.New(&desc, ctx)
		switch {
		case !test.buildErr && err != nil:
			t.Fatalf("%d: NewAction failed, desc %v failed, err %v.", tid, &desc, err)
		case test.buildErr && err == nil:
			t.Fatalf("%d: NewAction failed, desc %v did not fail.", tid, &desc)
		}
		if test.buildErr {
			continue
		}

		// If packet processing is expected to fail.
		if test.processErr {
			packet := mock_fwdpacket.NewMockPacket(ctrl)
			packet.EXPECT().Logf(gomock.Any(), gomock.Any()).AnyTimes()
			packet.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			packet.EXPECT().Field(fwdpacket.NewFieldID(field)).Return(test.original, nil)
			packet.EXPECT().String().Return("").AnyTimes()
			if _, state := action.Process(packet, nil); state != fwdaction.DROP {
				t.Errorf("%d: Packet processing did not fail for value %x, desc %+v", tid, test.original, &desc)
			}
			continue
		}

		// If packet processing is expected to suceed.
		packet := mock_fwdpacket.NewMockPacket(ctrl)
		packet.EXPECT().Logf(gomock.Any(), gomock.Any()).AnyTimes()
		packet.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		packet.EXPECT().Field(fwdpacket.NewFieldID(field)).Return(test.original, nil)
		packet.EXPECT().Update(fwdpacket.NewFieldID(field), fwdpacket.OpSet, test.final)
		if _, state := action.Process(packet, nil); state != fwdaction.CONTINUE {
			t.Errorf("%d: Packet processing failed for value %x, desc %+v", tid, test.original, &desc)
		}
	}
}

// TestBitAnd tests the bit AND/OR update action.
func TestBitAndOr(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a context with a test port.
	ctx := fwdcontext.New("test", "fwd")

	// Packet field
	field := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC,
		},
	}

	tests := []struct {
		original   []byte            // original bytes of the field
		value      []byte            // bits used for the operation
		op         *fwdpb.UpdateType // type of update
		final      []byte            // final bytes of the field after the operation
		processErr bool              // true if a packet processing error is expected
	}{
		// Test a processing failure as bit-mask is longer than the original packet field size.
		{
			original:   []byte{0x10, 0x01},
			value:      []byte{0xF1, 0x11, 0xFF},
			op:         fwdpb.UpdateType_UPDATE_TYPE_BIT_AND.Enum(),
			processErr: true,
		},
		// Test a processing failure as bit-mask is longer than the original packet field size.
		{
			original:   []byte{0x10, 0x01},
			value:      []byte{0xF1, 0x11, 0xFF},
			op:         fwdpb.UpdateType_UPDATE_TYPE_BIT_OR.Enum(),
			processErr: true,
		},
		// Test a successful AND case.
		{
			original: []byte{0x10, 0x01},
			value:    []byte{0xF1, 0x10},
			op:       fwdpb.UpdateType_UPDATE_TYPE_BIT_AND.Enum(),
			final:    []byte{0x10, 0x00},
		},
		// Test a successful OR case.
		{
			original: []byte{0x10, 0x01},
			value:    []byte{0xF1, 0x10},
			op:       fwdpb.UpdateType_UPDATE_TYPE_BIT_OR.Enum(),
			final:    []byte{0xF1, 0x11},
		},
	}

	for tid, test := range tests {
		t.Logf("%d: Running test %+v", tid, test)

		// Create an update action and check for build errors.
		desc := fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
		}
		update := fwdpb.UpdateActionDesc{
			FieldId: field,
			Type:    *test.op,
			Value:   test.value,
		}
		desc.Action = &fwdpb.ActionDesc_Update{
			Update: &update,
		}
		action, err := fwdaction.New(&desc, ctx)
		if err != nil {
			t.Fatalf("%d: NewAction failed, desc %v failed, err %v.", tid, &desc, err)
		}

		// If packet processing is expected to fail.
		if test.processErr {
			packet := mock_fwdpacket.NewMockPacket(ctrl)
			packet.EXPECT().Logf(gomock.Any(), gomock.Any()).AnyTimes()
			packet.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			packet.EXPECT().Field(fwdpacket.NewFieldID(field)).Return(test.original, nil)
			packet.EXPECT().String().Return("").AnyTimes()
			if _, state := action.Process(packet, nil); state != fwdaction.DROP {
				t.Errorf("%d: Packet processing did not fail for value %x, desc %+v", tid, test.original, &desc)
			}
			continue
		}

		// If packet processing is expected to suceed.
		packet := mock_fwdpacket.NewMockPacket(ctrl)
		packet.EXPECT().Logf(gomock.Any(), gomock.Any()).AnyTimes()
		packet.EXPECT().Logf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		packet.EXPECT().Field(fwdpacket.NewFieldID(field)).Return(test.original, nil)
		packet.EXPECT().Update(fwdpacket.NewFieldID(field), fwdpacket.OpSet, test.final)
		if _, state := action.Process(packet, nil); state != fwdaction.CONTINUE {
			t.Errorf("%d: Packet processing failed for value %x, desc %+v", tid, test.original, &desc)
		}
	}
}
