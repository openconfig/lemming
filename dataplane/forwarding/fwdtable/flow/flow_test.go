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

package flow

import (
	"errors"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/go-logr/logr/testr"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdset"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A testAction is a test action that encodes a flow and priority.
type testAction struct {
	flow     *Desc
	priority uint32
}

func (testAction) Process(fwdpacket.Packet, fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	return nil, fwdaction.DROP
}

func (testAction) String() string {
	return ""
}

// TestFlowMapManagement tests management of flows in the flow map.
func TestFlowMapManagement(t *testing.T) {
	// Create a controller and context.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := fwdcontext.New("test", "fwd")

	// Register a mock parser which uses a fixed fieldSize.
	const fieldSize = 4
	parser := mock_fwdpacket.NewMockParser(ctrl)
	parser.EXPECT().MaxSize(gomock.Any()).Return(fieldSize).AnyTimes()
	parser.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	fwdpacket.Register(parser)

	// Create a flow table and perform multiple add, remove and lookups
	id := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO,
		},
	}
	fid := fwdpacket.NewFieldID(id)
	fl := NewFieldList(map[fwdpacket.FieldID]bool{
		fid: true,
	})
	fm := NewMap(fl, fl, 0, 0)

	// Set of sets used to to test.
	testSets := []struct {
		id      string
		members [][]byte
	}{{
		id: "1",
		members: [][]byte{
			{0x01, 0x02},
			{0x01, 0x03},
		},
	}}

	// Set of flows used to test.
	testFlows := []struct {
		keys       []*fwdpb.PacketFieldMaskedBytes
		qualifiers []*fwdpb.PacketFieldSet
	}{{
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id,
			Bytes:   []byte{0x01, 0x02, 0x00, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
	}, {
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id,
			Bytes:   []byte{0x01, 0x00, 0x00, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
	}, {
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id,
			Bytes:   []byte{0x02, 0x00, 0x00, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
	}, {
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id,
			Bytes:   []byte{0x05, 0x06, 0x07, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
	}, {
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id,
			Bytes:   []byte{0x07, 0x07, 0x06, 0x08},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
	}, {
		qualifiers: []*fwdpb.PacketFieldSet{{
			FieldId: id,
			SetId: &fwdpb.SetId{
				ObjectId: &fwdpb.ObjectId{
					Id: "1",
				},
			},
		}},
	}, {
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id,
			Bytes:   []byte{0x07, 0x07, 0x06, 0x08},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
		qualifiers: []*fwdpb.PacketFieldSet{{
			FieldId: id,
			SetId: &fwdpb.SetId{
				ObjectId: &fwdpb.ObjectId{
					Id: "1",
				},
			},
		}},
	}}

	// Build the sets required for tests.
	for pos, ts := range testSets {
		set, err := fwdset.New(ctx, &fwdpb.SetId{
			ObjectId: &fwdpb.ObjectId{
				Id: ts.id,
			},
		})
		if err != nil {
			t.Errorf("%v: Unable to create test set, err %v", pos, err)
		}
		set.Update(ts.members)
	}

	// Build a list of flow desc. from the test.
	var flows []*Desc
	for pos, tf := range testFlows {
		fd, err := NewDesc(ctx, tf.keys, tf.qualifiers)
		if err == nil {
			flows = append(flows, fd)
		} else {
			t.Errorf("%v: Unable to create test flow, err %v", pos, err)
		}
	}

	// For each flow desc, perform a Lookup and Add.
	for _, f := range flows {
		if got := fm.Lookup(f); got != nil {
			t.Errorf("Lookup(%v)=%v before add.", f, got)
		}
		fm.Add(f, []*fwdaction.ActionAttr{fwdaction.NewActionAttr(&testAction{flow: f}, false)})
		got := fm.Lookup(f)
		if got == nil {
			t.Errorf("Lookup(%v)=nil after add.", f)
		}
		if len(got) != 1 {
			t.Errorf("Lookup(%v) has invalid length. Got %v, want 1.", f, got)
		}
		a := got[0].Action()
		if a == nil {
			t.Errorf("Lookup(%v)=nil after add.", f)
		}
		fa, ok := a.(*testAction)
		if !ok || !f.Equal(fa.flow) {
			t.Errorf("Lookup(%v) is invalid. Got %v, want %v.", f, got, f)
		}
	}

	// For each flow desc., perform a Lookup and Remove.
	for _, f := range flows {
		fm.Remove(f)
		if got := fm.Lookup(f); got != nil {
			t.Errorf("Lookup(%v) invalid after remove. Got %v, want nil.", f, got)
		}
	}
}

// TestFlowMapMatch tests matches in a flow map.
func TestFlowMapMatch(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := fwdcontext.New("test", "fwd")

	// Register a mock parser which uses a fixed fieldSize.
	const fieldSize = 4
	parser := mock_fwdpacket.NewMockParser(ctrl)
	parser.EXPECT().MaxSize(gomock.Any()).Return(fieldSize).AnyTimes()
	parser.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	fwdpacket.Register(parser)

	// Create a flow table and perform multiple add, remove and lookups.
	// The flow map has a single field. Hence for each flow, the field
	// mask is the flow key.
	id1 := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO,
		},
	}
	fid1 := fwdpacket.NewFieldID(id1)
	fl1 := NewFieldList(map[fwdpacket.FieldID]bool{
		fid1: true,
	})
	id2 := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP,
		},
	}
	fid2 := fwdpacket.NewFieldID(id2)
	fl2 := NewFieldList(map[fwdpacket.FieldID]bool{
		fid2: true,
	})
	fm := NewMap(fl1, fl2, 0, 0)

	// Set of sets used to to test.
	testSets := []struct {
		id      string
		members [][]byte
	}{{
		id: "1",
		members: [][]byte{
			{0x01, 0x02, 0x03, 0x04},
			{0x01, 0x03, 0x04, 0x05},
		},
	}}

	// List of flows used to test. Each flow is associated with a non-'id'
	noFlow := 0
	testFlows := []struct {
		id         int
		keys       []*fwdpb.PacketFieldMaskedBytes
		qualifiers []*fwdpb.PacketFieldSet
	}{{
		id: 1,
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id1,
			Bytes:   []byte{0x01, 0x02, 0x00, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
	}, {
		id: 2,
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id1,
			Bytes:   []byte{0x01, 0x00, 0x00, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
	}, {
		id: 3,
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id1,
			Bytes:   []byte{0x01, 0x00, 0x01, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFD, 0xFF},
		}},
	}, {
		id: 4,
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id1,
			Bytes:   []byte{0x02, 0x00, 0x00, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
	}, {
		id: 5,
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id1,
			Bytes:   []byte{0x20, 0x04, 0x01, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
	}, {
		id: 6,
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id1,
			Bytes:   []byte{0x20, 0x04, 0xF6, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
	}, {
		id: 7,
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id1,
			Bytes:   []byte{0x20, 0x04, 0xF6, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xF0, 0xFF},
		}},
	}, {
		id: 8,
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id1,
			Bytes:   []byte{0x20, 0x04, 0xF6, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xF0, 0xFF},
		}},
		qualifiers: []*fwdpb.PacketFieldSet{{
			FieldId: id2,
			SetId: &fwdpb.SetId{
				ObjectId: &fwdpb.ObjectId{
					Id: "1",
				},
			},
		}},
	}, {
		id: 9,
		keys: []*fwdpb.PacketFieldMaskedBytes{{
			FieldId: id1,
			Bytes:   []byte{0x20, 0x04, 0xF6, 0x00},
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF},
		}},
		qualifiers: []*fwdpb.PacketFieldSet{{
			FieldId: id2,
			SetId: &fwdpb.SetId{
				ObjectId: &fwdpb.ObjectId{
					Id: "1",
				},
			},
		}},
	}, {
		id: 10,
		qualifiers: []*fwdpb.PacketFieldSet{{
			FieldId: id2,
			SetId: &fwdpb.SetId{
				ObjectId: &fwdpb.ObjectId{
					Id: "1",
				},
			},
		}},
	}}
	matches := []struct {
		key       PacketKey
		qualifier PacketQualifier
		id        int
	}{{
		key: []byte{0x01, 0x02, 0x00, 0x00},
		id:  1,
	}, {
		key: []byte{0x01, 0x00, 0x03, 0x00},
		id:  3,
	}, {
		key: []byte{0x20, 0x04, 0xF6, 0x00},
		id:  7,
	}, {
		key: []byte{0x20, 0x04, 0xF7, 0x00},
		id:  7,
	}, {
		key: []byte{0x20, 0x04, 0xE2, 0x00},
		id:  noFlow,
	}, {
		key: []byte{0x20, 0x04, 0x00, 0x01},
		id:  noFlow,
	}, {
		key: []byte{0xAB, 0xCD, 0xEF, 0x00},
		qualifier: PacketQualifier{
			fid2: []byte{0x01, 0x03, 0x04, 0x05},
		},
		id: 10,
	}, {
		key: []byte{0xAB, 0xCD, 0xEF, 0x00},
		qualifier: PacketQualifier{
			fid2: []byte{0x01, 0x02, 0x03, 0x04},
		},
		id: 10,
	}, {
		key: []byte{0xAB, 0xCD, 0xEF, 0x00},
		qualifier: PacketQualifier{
			fid2: []byte{0x01, 0x02, 0x03, 0x05},
		},
		id: noFlow,
	}, {
		key: []byte{0x20, 0x04, 0xF6, 0x00},
		qualifier: PacketQualifier{
			fid2: []byte{0x01, 0x02, 0x03, 0x04},
		},
		id: 10,
	}}

	// Build the sets required for tests.
	for pos, ts := range testSets {
		set, err := fwdset.New(ctx, &fwdpb.SetId{
			ObjectId: &fwdpb.ObjectId{
				Id: ts.id,
			},
		})
		if err != nil {
			t.Errorf("%v: Unable to create test set, err %v", pos, err)
		}
		set.Update(ts.members)
	}

	// map of flows indexed by the id.
	mapFlows := make(map[int]*Desc)

	for pos, f := range testFlows {
		fd, err := NewDesc(ctx, f.keys, f.qualifiers)
		if err != nil {
			t.Fatalf("%v: Unable to create flow test, err %v", pos, err)
		}
		mapFlows[f.id] = fd
		fm.Add(fd, []*fwdaction.ActionAttr{fwdaction.NewActionAttr(&testAction{flow: fd}, false)})
	}

	for id, match := range matches {
		_, _, got := fm.Match(match.key, match.qualifier)
		want := mapFlows[match.id]
		switch {
		case want == nil && len(got) != 0:
			t.Errorf("%v: Match(%x, %v) is invalid. Got %v len %v", id, match.key, match.qualifier, got, len(got))

		case want != nil && len(got) == 0:
			t.Errorf("%v: Match(%x, %v) is invalid. Want %v", id, match.key, match.qualifier, want)

		case want != nil && len(got) != 1:
			t.Errorf("%v: Match(%x, %v) has invalid length. Got %v, want 1.", id, match.key, match.qualifier, got)

		case want != nil && len(got) == 1:
			a := got[0].Action()
			fa, ok := a.(*testAction)
			if !ok || !want.Equal(fa.flow) {
				t.Errorf("%v: Match(%x, %v) is invalid. Got %v, want %v.", id, match.key, match.qualifier, fa.flow, want)
			}
		}
	}
}

// A testKey describes a key formed by a single packet field. The testKey
// descibes the field's id  and bytes. The field contains size number of
// bytes with value encoded at index 0.
type testKey struct {
	id    fwdpb.PacketFieldNum // id identifies the field in the key
	value byte                 // value of the byte at index 0
	size  int                  // number of bytes in the field
}

// field returns the FieldID of a test key.
func (k *testKey) field() fwdpacket.FieldID {
	return fwdpacket.NewFieldIDFromNum(k.id, 0)
}

// bytes returns the bytes of a test key.
func (k *testKey) bytes() []byte {
	b := make([]byte, k.size)
	b[0] = k.value
	return b
}

// mask returns the mask of a test key.
func (k *testKey) mask() []byte {
	m := make([]byte, k.size)
	m[0] = byte(0xFF)
	return m
}

// A testQualifier describes a qualifier formed by a single packet field. The
// testKey descibes the field's id, set-id and packet bytes. The packet bytes
// are a slice of bytes contains size number of bytes with value encoded at
// index 0. The packet field bytes must be a member of the forwarding set
// referenced by the qualifier.
type testQualifier struct {
	id    fwdpb.PacketFieldNum // id identifies the field in the key
	value byte                 // value of the byte at index 0
	size  int                  // number of bytes in the field
	set   string               // id of the forwarding set
}

// field returns the FieldID of a test qualifier.
func (q *testQualifier) field() fwdpacket.FieldID {
	return fwdpacket.NewFieldIDFromNum(q.id, 0)
}

// bytes returns the bytes of a test qualifier.
func (q *testQualifier) bytes() []byte {
	b := make([]byte, q.size)
	b[0] = q.value
	return b
}

// A testBuilder builds a test action. Before an action is built, the builder
// is initialized with a set of flow descriptors indexed by priority. The
// ActionDesc encodes the priority of the flow. To build an action, the builder
// looks up the flow descriptor corresponding to the priority encoded in the
// ActionDesc, and encodes it as a part of the action.
type testBuilder struct {
	flows map[int]*Desc
}

// newBuilder creates and registers a new builder for the test actions.
func newBuilder() *testBuilder {
	tb := testBuilder{
		flows: make(map[int]*Desc),
	}
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_TEST, &tb)
	return &tb
}

// Build creates a new test action.
func (t *testBuilder) Build(desc *fwdpb.ActionDesc, _ *fwdcontext.Context) (fwdaction.Action, error) {
	ta, ok := desc.Action.(*fwdpb.ActionDesc_Test)
	if !ok {
		return nil, fmt.Errorf("flow: Build for test action failed, missing desc")
	}
	fd, _ := t.flows[int(ta.Test.GetInt1())]
	return &testAction{
		flow:     fd,
		priority: ta.Test.GetInt1(),
	}, nil
}

// genActionDesc creates a desc for a test action encoding the key and its
// priority.
func genActionDesc(keys []testKey, priority uint32) []*fwdpb.ActionDesc {
	var b []byte
	for _, k := range keys {
		t := k.bytes()
		b = append(b, t...)
	}
	desc := &fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_ACTION_TYPE_TEST,
	}
	desc.Action = &fwdpb.ActionDesc_Test{
		Test: &fwdpb.TestActionDesc{
			Int1:   priority,
			Bytes1: b,
		},
	}
	return []*fwdpb.ActionDesc{desc}
}

// genFlowDesc generates a flow desc for a test flow using the specified bank
// and priority.
func genFlowDesc(t *testing.T, tb *testBuilder, bank, priority uint32, keys []testKey, qualifiers []testQualifier) *fwdpb.EntryDesc {
	var f []*fwdpb.PacketFieldMaskedBytes
	for _, k := range keys {
		f = append(f, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{
				Field: &fwdpb.PacketField{FieldNum: k.id},
			},
			Bytes: k.bytes(),
			Masks: k.mask(),
		})
	}
	var q []*fwdpb.PacketFieldSet
	for _, k := range qualifiers {
		q = append(q, &fwdpb.PacketFieldSet{
			FieldId: &fwdpb.PacketFieldId{
				Field: &fwdpb.PacketField{FieldNum: k.id},
			},
			SetId: &fwdpb.SetId{
				ObjectId: &fwdpb.ObjectId{
					Id: k.set,
				},
			},
		})
	}
	desc := &fwdpb.EntryDesc{}
	flow := &fwdpb.FlowEntryDesc{
		Priority:   priority,
		Bank:       bank,
		Fields:     f,
		Qualifiers: q,
	}
	fd, err := NewDesc(nil, nil, nil)
	if err != nil {
		t.Fatalf("Unable to initialize flow for test, err %v", err)
		return nil
	}
	tb.flows[int(priority)] = fd
	desc.Entry = &fwdpb.EntryDesc_Flow{
		Flow: flow,
	}
	return desc
}

// flowTable creates a flow table with the specified number of banks.
func flowTable(ctx *fwdcontext.Context, bank uint32, index int) (fwdtable.Table, error) {
	desc := &fwdpb.TableDesc{
		TableType: fwdpb.TableType_TABLE_TYPE_FLOW,
		TableId:   fwdtable.MakeID(fwdobject.NewID(fmt.Sprintf("TABLE=%v", index))),
	}
	flow := &fwdpb.FlowTableDesc{
		BankCount: bank,
	}
	desc.Table = &fwdpb.TableDesc_Flow{
		Flow: flow,
	}
	return fwdtable.New(ctx, desc)
}

// TestFlowTable tests various operations on a flow table table.
func TestFlowTable(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Register a mock parser which uses a fixed fieldSize.
	const fieldSize = 4
	parser := mock_fwdpacket.NewMockParser(ctrl)
	parser.EXPECT().MaxSize(gomock.Any()).Return(fieldSize).AnyTimes()
	parser.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	fwdpacket.Register(parser)

	ctx := fwdcontext.New("test", "fwd")

	// Operations that can be performed on the table.
	const (
		ADD = iota
		REMOVE
		MATCH
		CLEAR
	)

	// Set of sets used to to test.
	testSets := []struct {
		id      string
		members [][]byte
	}{{
		id: "1",
		members: [][]byte{
			{0x0a, 0x00, 0x00, 0x00},
			{0x0b, 0x00, 0x00, 0x00},
		},
	}, {
		id: "2",
		members: [][]byte{
			{0x0c, 0x00, 0x00, 0x00},
			{0x0d, 0x00, 0x00, 0x00},
		},
	}}

	// Build the sets required for tests.
	for pos, ts := range testSets {
		set, err := fwdset.New(ctx, &fwdpb.SetId{
			ObjectId: &fwdpb.ObjectId{
				Id: ts.id,
			},
		})
		if err != nil {
			t.Errorf("%v: Unable to create test set, err %v", pos, err)
		}
		set.Update(ts.members)
	}

	// testOp describes an operation on the table.
	//
	// In case of ADD (or REMOVE), the priority and keys describe the
	// flow being added (or removed). In case of MATCH, keys describes
	// the packet being processed by the table and priority describes
	// the priority of the expected flow matching the packet.
	type testOp struct {
		event      int             // type of operation
		priority   uint32          // priority of flow
		keys       []testKey       // keys describing the flow keys
		qualifiers []testQualifier // keys describing the flow qualifiers
		count      int             // number of entries in the table after the operation
	}

	// Canned keys used to describe tests.
	field1 := testKey{
		id:    fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF,
		value: 10,
		size:  fieldSize,
	}
	field2 := testKey{
		id:    fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_LENGTH,
		value: 20,
		size:  fieldSize,
	}

	// Canned qualfiers used to describe tests. Note that the value must
	// match the forwarding sets created for the test.
	qual1 := testQualifier{
		id:    fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC,
		set:   "1",
		value: 10,
		size:  fieldSize,
	}
	qual2 := testQualifier{
		id:    fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
		set:   "2",
		value: 12,
		size:  fieldSize,
	}
	tests := [][]testOp{{{ // Flow table with a single entry.
		count:    1,
		event:    ADD,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    MATCH,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    0,
		event:    REMOVE,
		priority: 10,
		keys:     []testKey{field1},
	}}, {{ // Flow table with adds / removes of the same flow with different priorities.
		count:    1,
		event:    ADD,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    MATCH,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    2,
		event:    ADD,
		priority: 5,
		keys:     []testKey{field1},
	}, {
		count:    2,
		event:    MATCH,
		priority: 5,
		keys:     []testKey{field1},
	}, {
		count:    3,
		event:    ADD,
		priority: 15,
		keys:     []testKey{field1},
	}, {
		count:    3,
		event:    MATCH,
		priority: 5,
		keys:     []testKey{field1},
	}, {
		count:    2,
		event:    REMOVE,
		priority: 5,
		keys:     []testKey{field1},
	}, {
		count:    2,
		event:    MATCH,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    REMOVE,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    MATCH,
		priority: 15,
		keys:     []testKey{field1},
	}}, {{ // Flow table with adds / removes of flows with different keys
		count:    1,
		event:    ADD,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    MATCH,
		priority: 10,
		keys:     []testKey{field1, field2},
	}, {
		count:    2,
		event:    ADD,
		priority: 5,
		keys:     []testKey{field1, field2},
	}, {
		count:    2,
		event:    MATCH,
		priority: 5,
		keys:     []testKey{field1, field2},
	}, {
		count:    1,
		event:    REMOVE,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    MATCH,
		priority: 5,
		keys:     []testKey{field1, field2},
	}}, {{ // Flow table with a add and match of three entries of the same flow.
		count:    1,
		event:    ADD,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    2,
		event:    ADD,
		priority: 20,
		keys:     []testKey{field1},
	}, {
		count:    3,
		event:    ADD,
		priority: 15,
		keys:     []testKey{field1},
	}, {
		count:    3,
		event:    MATCH,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    2,
		event:    REMOVE,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    2,
		event:    MATCH,
		priority: 15,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    REMOVE,
		priority: 15,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    MATCH,
		priority: 20,
		keys:     []testKey{field1},
	}}, {{ // Flow table with a single entry (with keys and qualifiers).
		count:      1,
		event:      ADD,
		priority:   10,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}, {
		count:      1,
		event:      MATCH,
		priority:   10,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}, {
		count:      0,
		event:      REMOVE,
		priority:   10,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}}, {{ // Flow table with adds / removes of the same flow with different priorities (keys and qualifiers).
		count:      1,
		event:      ADD,
		priority:   10,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}, {
		count:      2,
		event:      ADD,
		priority:   15,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}, {
		count:      2,
		event:      MATCH,
		priority:   10,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}, {
		count:      3,
		event:      ADD,
		priority:   5,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}, {
		count:      3,
		event:      MATCH,
		priority:   5,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}, {
		count:      4,
		event:      ADD,
		priority:   6,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1, qual2},
	}, {
		count:      4,
		event:      MATCH,
		priority:   5,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual2, qual1},
	}, {
		count:      4,
		event:      MATCH,
		priority:   5,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}, {
		count:      3,
		event:      REMOVE,
		priority:   5,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}, {
		count:      3,
		event:      MATCH,
		priority:   10,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual1},
	}, {
		count:      3,
		event:      MATCH,
		priority:   6,
		keys:       []testKey{field1},
		qualifiers: []testQualifier{qual2, qual1},
	}}, {{ // Flow table with a single entry (no keys).
		count:      1,
		event:      ADD,
		priority:   10,
		qualifiers: []testQualifier{qual1},
	}, {
		count:      1,
		event:      MATCH,
		priority:   10,
		qualifiers: []testQualifier{qual1},
	}, {
		count:      0,
		event:      REMOVE,
		priority:   10,
		qualifiers: []testQualifier{qual1},
	}}, {{ // Flow table with a single entry that can be matched. Repeated after a clear.
		count:    1,
		event:    ADD,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    MATCH,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count: 0,
		event: CLEAR,
	}, {
		count:    1,
		event:    ADD,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    MATCH,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    0,
		event:    REMOVE,
		priority: 10,
		keys:     []testKey{field1},
	}}, {{ // Flow table is cleared when it is empty.
		count: 0,
		event: CLEAR,
	}, {
		count:    1,
		event:    ADD,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    MATCH,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count: 0,
		event: CLEAR,
	}, {
		count:    1,
		event:    ADD,
		priority: 10,
		keys:     []testKey{field1},
	}, {
		count:    1,
		event:    MATCH,
		priority: 10,
		keys:     []testKey{field1},
	}}}
next:
	// For simplicity, all flows are added to a single bank in the flow table.
	for tid, test := range tests {
		tb := newBuilder()
		bankCount := uint32(4)
		bankID := uint32(1)
		table, err := flowTable(ctx, bankCount, tid)
		if err != nil {
			t.Errorf("%d: Flow match table create failed, err %v.", tid, err)
		}
		if entries := table.Entries(); len(entries) != 0 {
			t.Errorf("%d: Incorrect number of table entries. Got %v, want 0.", tid, len(entries))
		}

		for opid, op := range test {
			switch op.event {
			case ADD:
				if err := table.AddEntry(genFlowDesc(t, tb, bankID, op.priority, op.keys, op.qualifiers), genActionDesc(op.keys, op.priority)); err != nil {
					t.Errorf("%d, opid %d: AddEntry failed, err %v.", tid, opid, err)
					continue next
				}
			case REMOVE:
				if err := table.RemoveEntry(genFlowDesc(t, tb, bankID, op.priority, op.keys, op.qualifiers)); err != nil {
					t.Errorf("%d, opid %d: RemoveEntry failed, err %v.", tid, opid, err)
					continue next
				}
			case MATCH:
				packet := mock_fwdpacket.NewMockPacket(ctrl)
				for _, k := range op.keys {
					packet.EXPECT().Field(k.field()).Return(k.bytes(), nil).AnyTimes()
				}
				for _, q := range op.qualifiers {
					packet.EXPECT().Field(q.field()).Return(q.bytes(), nil).AnyTimes()
				}
				packet.EXPECT().Field(gomock.Any()).Return(nil, errors.New("no field")).AnyTimes()
				packet.EXPECT().Log().Return(testr.New(t)).AnyTimes()

				actions, state := table.Process(packet, nil)
				if state != fwdaction.CONTINUE {
					t.Errorf("%d, opid %d: Incorrect state after packet processing. Got %v, want %v.", tid, opid, state, fwdaction.CONTINUE)
					continue next
				}
				if len(actions) != 1 {
					t.Errorf("%d, opid %d: Incorrect action count after packet processing. Got %v, want 1.", tid, opid, len(actions))
					continue next
				}
				fa := actions[0].Action().(*testAction)
				if fa.priority != op.priority {
					t.Errorf("%d, opid %d: Incorrect priority matched packet. Got %v, want %v.", tid, opid, fa.priority, op.priority)
					continue next
				}
			case CLEAR:
				table.Clear()
			}
			if entries := table.Entries(); len(entries) != op.count {
				t.Errorf("%d, opid %d: Incorrect number of table entries. Got %v, want %v.", tid, opid, len(entries), op.count)
			}
		}
	}
}
