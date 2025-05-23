// Copyright 2023 Google LLC
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

// Package fwdconfig contains builders for varius proto types.
package fwdconfig

import (
	"encoding/binary"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type ActionDescBuilder interface {
	set(*fwdpb.ActionDesc)
	actionType() fwdpb.ActionType
}

var (
	// Compile-time checks that builders implement actionDescBuilder interface.
	_ ActionDescBuilder = &UpdateActionBuilder{}
	_ ActionDescBuilder = &TransmitActionBuilder{}
	_ ActionDescBuilder = &LookupActionBuilder{}
	_ ActionDescBuilder = &EncapActionBuilder{}
	_ ActionDescBuilder = &DecapActionBuilder{}
	_ ActionDescBuilder = &DropActionBuilder{}
)

// ActionBuilder is a builder for forward action types.
type ActionBuilder struct {
	adb ActionDescBuilder
}

// Action returns a new action builder.
func Action(adb ActionDescBuilder) *ActionBuilder {
	return &ActionBuilder{
		adb: adb,
	}
}

// WithActionDesc sets the action description
func (ab *ActionBuilder) WithActionDesc(adb ActionDescBuilder) *ActionBuilder {
	ab.adb = adb
	return ab
}

// Build returns a new action.
func (ab *ActionBuilder) Build() *fwdpb.ActionDesc {
	a := &fwdpb.ActionDesc{
		ActionType: ab.adb.actionType(),
	}
	ab.adb.set(a)
	return a
}

// UpdateActionBuilder is a builder for an update action.
type UpdateActionBuilder struct {
	fieldIDNum    fwdpb.PacketFieldNum
	updateType    fwdpb.UpdateType
	fieldSrc      fwdpb.PacketFieldNum
	instance      uint32
	srcInstance   uint32
	offset, count uint32 // Offset and count in bits for BIT_WRITE update.
	value         []byte
}

// UpdateAction returns a new update action builder.
func UpdateAction(t fwdpb.UpdateType, num fwdpb.PacketFieldNum) *UpdateActionBuilder {
	return &UpdateActionBuilder{
		fieldIDNum: num,
		updateType: t,
	}
}

// WithFieldIDNum sets the packet field id enum.
func (u *UpdateActionBuilder) WithFieldIDNum(num fwdpb.PacketFieldNum) *UpdateActionBuilder {
	u.fieldIDNum = num
	return u
}

// WithUpdateType sets the update type.
func (u *UpdateActionBuilder) WithUpdateType(t fwdpb.UpdateType) *UpdateActionBuilder {
	u.updateType = t
	return u
}

// WithFieldIDInstance sets the instance of the id num.
func (u *UpdateActionBuilder) WithFieldIDInstance(i uint32) *UpdateActionBuilder {
	u.instance = i
	return u
}

// WithFieldSrcInstance sets the instance of the src id num.
func (u *UpdateActionBuilder) WithFieldSrcInstance(i uint32) *UpdateActionBuilder {
	u.srcInstance = i
	return u
}

// WithUpdateType sets the update source field id.
func (u *UpdateActionBuilder) WithFieldSrc(num fwdpb.PacketFieldNum) *UpdateActionBuilder {
	u.fieldSrc = num
	return u
}

// WithUint64Value sets the value to a uint64 value.
func (u *UpdateActionBuilder) WithUint64Value(v uint64) *UpdateActionBuilder {
	u.value = binary.BigEndian.AppendUint64(nil, v)
	return u
}

// WithUint64Value sets the value to a byte slice.
func (u *UpdateActionBuilder) WithValue(v []byte) *UpdateActionBuilder {
	u.value = v
	return u
}

// WithBitOp sets the bit count and offset fields.
func (u *UpdateActionBuilder) WithBitOp(count, offset int) *UpdateActionBuilder {
	u.offset = uint32(offset)
	u.count = uint32(count)
	return u
}

func (u *UpdateActionBuilder) set(ad *fwdpb.ActionDesc) {
	upd := &fwdpb.ActionDesc_Update{
		Update: &fwdpb.UpdateActionDesc{
			FieldId: &fwdpb.PacketFieldId{
				Field: &fwdpb.PacketField{
					FieldNum: u.fieldIDNum,
					Instance: u.instance,
				},
			},
			Type:      u.updateType,
			Value:     u.value,
			BitOffset: u.offset,
			BitCount:  u.count,
			Field: &fwdpb.PacketFieldId{
				Field: &fwdpb.PacketField{
					FieldNum: u.fieldSrc,
					Instance: u.srcInstance,
				},
			},
		},
	}
	ad.Action = upd
}

func (u *UpdateActionBuilder) actionType() fwdpb.ActionType {
	return fwdpb.ActionType_ACTION_TYPE_UPDATE
}

// TransmitActionBuilder is a builder for a transmit action.
type TransmitActionBuilder struct {
	portID    string
	immediate bool
}

// TransmitAction returns a new update action builder.
func TransmitAction(portID string) *TransmitActionBuilder {
	return &TransmitActionBuilder{
		portID: portID,
	}
}

// WithFieldIDNum sets the port id value.
func (u *TransmitActionBuilder) WithPortID(id string) *TransmitActionBuilder {
	u.portID = id
	return u
}

// WithImmediate sets immediate option value.
func (u *TransmitActionBuilder) WithImmediate(immediate bool) *TransmitActionBuilder {
	u.immediate = immediate
	return u
}

func (u *TransmitActionBuilder) set(ad *fwdpb.ActionDesc) {
	upd := &fwdpb.ActionDesc_Transmit{
		Transmit: &fwdpb.TransmitActionDesc{
			PortId: &fwdpb.PortId{
				ObjectId: &fwdpb.ObjectId{
					Id: u.portID,
				},
			},
			Immediate: u.immediate,
		},
	}
	ad.Action = upd
}

func (u *TransmitActionBuilder) actionType() fwdpb.ActionType {
	return fwdpb.ActionType_ACTION_TYPE_TRANSMIT
}

// LookupActionBuilder is a builder for a lookup action.
type LookupActionBuilder struct {
	tableID string
}

// LookupAction returns a new lookup action builder.
func LookupAction(tableID string) *LookupActionBuilder {
	return &LookupActionBuilder{
		tableID: tableID,
	}
}

// WithTableID sets the table id.
func (u *LookupActionBuilder) WithTableID(id string) *LookupActionBuilder {
	u.tableID = id
	return u
}

func (u *LookupActionBuilder) set(ad *fwdpb.ActionDesc) {
	upd := &fwdpb.ActionDesc_Lookup{
		Lookup: &fwdpb.LookupActionDesc{
			TableId: &fwdpb.TableId{
				ObjectId: &fwdpb.ObjectId{
					Id: u.tableID,
				},
			},
		},
	}
	ad.Action = upd
}

func (u *LookupActionBuilder) actionType() fwdpb.ActionType {
	return fwdpb.ActionType_ACTION_TYPE_LOOKUP
}

// EncapActionBuilder is a builder for a lookup action.
type EncapActionBuilder struct {
	header fwdpb.PacketHeaderId
}

// LookupAction returns a new lookup action builder.
func EncapAction(header fwdpb.PacketHeaderId) *EncapActionBuilder {
	return &EncapActionBuilder{
		header: header,
	}
}

// WithHeaderID sets the header id.
func (u *EncapActionBuilder) WithHeaderID(header fwdpb.PacketHeaderId) *EncapActionBuilder {
	u.header = header
	return u
}

func (u *EncapActionBuilder) set(ad *fwdpb.ActionDesc) {
	upd := &fwdpb.ActionDesc_Encap{
		Encap: &fwdpb.EncapActionDesc{
			HeaderId: u.header,
		},
	}
	ad.Action = upd
}

func (u *EncapActionBuilder) actionType() fwdpb.ActionType {
	return fwdpb.ActionType_ACTION_TYPE_ENCAP
}

// DecapActionBuilder is a builder for a lookup action.
type DecapActionBuilder struct {
	header fwdpb.PacketHeaderId
}

// LookupAction returns a new lookup action builder.
func DecapAction(header fwdpb.PacketHeaderId) *DecapActionBuilder {
	return &DecapActionBuilder{
		header: header,
	}
}

// WithHeaderID sets the header id.
func (u *DecapActionBuilder) WithHeaderID(header fwdpb.PacketHeaderId) *DecapActionBuilder {
	u.header = header
	return u
}

func (u *DecapActionBuilder) set(ad *fwdpb.ActionDesc) {
	upd := &fwdpb.ActionDesc_Decap{
		Decap: &fwdpb.DecapActionDesc{
			HeaderId: u.header,
		},
	}
	ad.Action = upd
}

func (u *DecapActionBuilder) actionType() fwdpb.ActionType {
	return fwdpb.ActionType_ACTION_TYPE_DECAP
}

// DropActionBuilder is a builder for a drop action.
type DropActionBuilder struct{}

// LookupAction returns a new drop action builder.
func DropAction() *DropActionBuilder {
	return &DropActionBuilder{}
}

func (u *DropActionBuilder) set(*fwdpb.ActionDesc) {
}

func (u *DropActionBuilder) actionType() fwdpb.ActionType {
	return fwdpb.ActionType_ACTION_TYPE_DROP
}

// FlowActionBuilder is a builder for a flow action.
type FlowCounterActionBuilder struct {
	id string
}

// FlowAction returns a new flow action builder.
func FlowCounterAction(id string) *FlowCounterActionBuilder {
	return &FlowCounterActionBuilder{
		id: id,
	}
}

func (u *FlowCounterActionBuilder) set(a *fwdpb.ActionDesc) {
	a.Action = &fwdpb.ActionDesc_Flow{
		Flow: &fwdpb.FlowCounterActionDesc{
			CounterId: &fwdpb.FlowCounterId{
				ObjectId: &fwdpb.ObjectId{
					Id: u.id,
				},
			},
		},
	}
}

func (u *FlowCounterActionBuilder) actionType() fwdpb.ActionType {
	return fwdpb.ActionType_ACTION_TYPE_FLOW_COUNTER
}

// ContinueActionBuilder is a builder for a continue action.
type ContinueActionBuilder struct{}

// ContinueAction returns a new continue action builder.
func ContinueAction() *ContinueActionBuilder {
	return &ContinueActionBuilder{}
}

func (u *ContinueActionBuilder) set(*fwdpb.ActionDesc) {
}

func (u *ContinueActionBuilder) actionType() fwdpb.ActionType {
	return fwdpb.ActionType_ACTION_TYPE_CONTINUE
}

// MirrorActionBuilder is a builder for a mirror action.
type MirrorActionBuilder struct {
	portID  string
	portAct fwdpb.PortAction
	fields  []*PacketFieldIdBuilder
	act     []*ActionBuilder
}

// MirrorAction returns a new mirror action builder.
func MirrorAction() *MirrorActionBuilder {
	return &MirrorActionBuilder{}
}

func (m *MirrorActionBuilder) WithPort(id string, act fwdpb.PortAction) *MirrorActionBuilder {
	m.portID = id
	m.portAct = act
	return m
}

func (m *MirrorActionBuilder) WithFields(f ...*PacketFieldIdBuilder) *MirrorActionBuilder {
	m.fields = f
	return m
}

func (m *MirrorActionBuilder) WithActions(a ...*ActionBuilder) *MirrorActionBuilder {
	m.act = a
	return m
}

func (m *MirrorActionBuilder) set(a *fwdpb.ActionDesc) {
	fields := []*fwdpb.PacketFieldId{}
	for _, f := range m.fields {
		fields = append(fields, f.Build())
	}
	act := []*fwdpb.ActionDesc{}
	for _, a := range m.act {
		act = append(act, a.Build())
	}

	a.Action = &fwdpb.ActionDesc_Mirror{
		Mirror: &fwdpb.MirrorActionDesc{
			PortId:     &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: m.portID}},
			PortAction: m.portAct,
			FieldIds:   fields,
			Actions:    act,
		},
	}
}

func (m *MirrorActionBuilder) actionType() fwdpb.ActionType {
	return fwdpb.ActionType_ACTION_TYPE_MIRROR
}
