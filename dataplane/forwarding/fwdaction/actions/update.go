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
	"fmt"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	byteBitCount = 8 // number of bits stored in a byte
	byteLogCount = 3 // log to the base 2 of the byteBitCount
)

// getBit returns the bit at the specified position within a slice of bytes.
// It is assumed that bitpos (0 indexed) exists within the byte slice and
// bytes has at least one byte.
func getBit(bytes []byte, bitpos uint) byte {
	idx := len(bytes) - 1 - int(bitpos>>byteLogCount)
	pos := (bitpos & (byteBitCount - 1))
	return byte((bytes[idx] >> uint(pos))) & 0x1
}

// writeBit sets the bit at the specified position.
// It is assumed that bitpos (0 indexed) exists within the byte slice and
// bytes has atleast one byte.
func writeBit(bytes []byte, bitpos uint, v byte) {
	idx := len(bytes) - 1 - int(bitpos>>byteLogCount)
	pos := (bitpos & (byteBitCount - 1))
	mask := byte(1 << uint(pos))
	if v == 0x1 {
		bytes[idx] |= mask
	} else {
		bytes[idx] &= ^mask
	}
}

// update is an action that writes a packet to a port.
type update struct {
	fieldID   fwdpacket.FieldID // Field being update.
	bytesArg  []byte            // Operand specified as bytes.
	fieldArg  fwdpacket.FieldID // Operand specified as a field.
	op        fwdpb.UpdateType  // Operation to be performed on the field.
	bitCount  uint              // number of bits to set.
	bitOffset uint              // offset in the packet field to set bits.
}

// String formats the state of the action as a string.
func (u *update) String() string {
	return fmt.Sprintf("Type=%v;Field=%v;Op=%v;ByteArg=%x;FieldArg=%v;BitCount=%v;BitOffset=%v", fwdpb.ActionType_ACTION_TYPE_UPDATE, u.fieldID, u.op, u.bytesArg, u.fieldArg, u.bitCount, u.bitOffset)
}

// Process updates packet by applying an operation on a field.
// INC/DEC: These are rarely used operations only supported for select fields
//
//	like IP TTL/HOP. They are implemented within the  protocol handlers
//	as they are accompanied by optimized checksum adjustments.
//
// SET:  This is a primitive operation which is implemented for all fields in the
//
//	protocol handlers.
//
// COPY: This operation is implemented using the Get/Set of the packet field.
// BIT:  This operation is implemented using the Get/Set of the packet field.
func (u *update) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	var e error
	defer func() {
		if e != nil {
			log.Errorf("actions: Update %v on packet %v failed, err %v", u, packet, e)
		}
	}()
	switch u.op {
	case fwdpb.UpdateType_UPDATE_TYPE_INC:
		e = packet.Update(u.fieldID, fwdpacket.OpInc, u.bytesArg)

	case fwdpb.UpdateType_UPDATE_TYPE_DEC:
		e = packet.Update(u.fieldID, fwdpacket.OpDec, u.bytesArg)

	case fwdpb.UpdateType_UPDATE_TYPE_SET:
		e = packet.Update(u.fieldID, fwdpacket.OpSet, u.bytesArg)

	case fwdpb.UpdateType_UPDATE_TYPE_COPY:
		arg, err := packet.Field(u.fieldArg)
		if err != nil {
			e = err
			return nil, fwdaction.DROP
		}
		e = packet.Update(u.fieldID, fwdpacket.OpSet, fwdpacket.Pad(u.fieldID, arg))

	case fwdpb.UpdateType_UPDATE_TYPE_BIT_WRITE:
		arg, err := packet.Field(u.fieldID)
		if err != nil {
			e = err
			return nil, fwdaction.DROP
		}
		if (byteBitCount * uint(len(arg))) < (u.bitCount + u.bitOffset) {
			e = fmt.Errorf("actions: Update %v on packet %v failed, field of length %v does not have %d bits at offset %d", u, packet, len(arg), u.bitCount, u.bitOffset)
			return nil, fwdaction.DROP
		}
		for i := uint(0); i < u.bitCount; i++ {
			writeBit(arg, i+u.bitOffset, getBit(u.bytesArg, i))
		}
		e = packet.Update(u.fieldID, fwdpacket.OpSet, arg)

	case fwdpb.UpdateType_UPDATE_TYPE_BIT_AND:
		arg, err := packet.Field(u.fieldID)
		if err != nil {
			e = err
			return nil, fwdaction.DROP
		}
		if len(arg) < len(u.bytesArg) {
			e = fmt.Errorf("actions: Update %v on packet %v failed, field %x of length %v is shorter than %x", u, packet, arg, len(arg), u.bytesArg)
			return nil, fwdaction.DROP
		}
		for i := 0; i < len(u.bytesArg); i++ {
			arg[i] = arg[i] & u.bytesArg[i]
		}
		e = packet.Update(u.fieldID, fwdpacket.OpSet, arg)

	case fwdpb.UpdateType_UPDATE_TYPE_BIT_OR:
		arg, err := packet.Field(u.fieldID)
		if err != nil {
			e = err
			return nil, fwdaction.DROP
		}
		if len(arg) < len(u.bytesArg) {
			e = fmt.Errorf("actions: Update %v on packet %v failed, field %x of length %v is shorter than %x", u, packet, arg, len(arg), u.bytesArg)
			return nil, fwdaction.DROP
		}
		for i := 0; i < len(u.bytesArg); i++ {
			arg[i] = arg[i] | u.bytesArg[i]
		}
		e = packet.Update(u.fieldID, fwdpacket.OpSet, arg)

	}
	if e != nil {
		return nil, fwdaction.DROP
	}
	return nil, fwdaction.CONTINUE
}

// updateBuilder builds update actions.
type updateBuilder struct{}

func init() {
	// Register a builder for the update action type.
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_UPDATE, &updateBuilder{})
}

// Build creates a new update action.
func (*updateBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	upd, ok := desc.Action.(*fwdpb.ActionDesc_Update)
	if !ok {
		return nil, fmt.Errorf("actions: Build for update action failed, missing extension")
	}
	// validate that the specified arguments have enough bits.
	switch upd.Update.GetType() {
	case fwdpb.UpdateType_UPDATE_TYPE_INC, fwdpb.UpdateType_UPDATE_TYPE_DEC, fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.UpdateType_UPDATE_TYPE_BIT_OR, fwdpb.UpdateType_UPDATE_TYPE_BIT_AND:
		if len(upd.Update.GetValue()) == 0 {
			return nil, fmt.Errorf("actions: Build for update action failed, update %+v is missing the bytes argument", upd.Update)
		}
	case fwdpb.UpdateType_UPDATE_TYPE_BIT_WRITE:
		if len(upd.Update.GetValue()) == 0 {
			return nil, fmt.Errorf("actions: Build for update action failed, update %+v is missing the bytes argument", upd.Update)
		}
		if uint(upd.Update.GetBitCount()) > byteBitCount*uint(len(upd.Update.GetValue())) {
			return nil, fmt.Errorf("actions: Build for update action failed, bit update value %v has fewer bytes than the bit count %v", upd.Update.GetValue(), upd.Update.GetBitCount())
		}
	case fwdpb.UpdateType_UPDATE_TYPE_COPY:
		if upd.Update.Field == nil {
			return nil, fmt.Errorf("actions: Build for update action failed, update %+v is missing the field argument", upd.Update)
		}
	}
	return &update{
		fieldID:   fwdpacket.NewFieldID(upd.Update.GetFieldId()),
		op:        upd.Update.GetType(),
		bytesArg:  upd.Update.GetValue(),
		fieldArg:  fwdpacket.NewFieldID(upd.Update.GetField()),
		bitCount:  uint(upd.Update.GetBitCount()),
		bitOffset: uint(upd.Update.GetBitOffset()),
	}, nil
}
