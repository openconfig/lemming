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
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdset"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
	"google.golang.org/protobuf/proto"
)

// A PacketKey is a sequence of bytes formed by concatenating packet fields in
// a sequence specified by a FieldList. Note that each field is padded to its
// maximum size. If the packet does not contain a field in the FieldList, the
// value is set to zero.
type PacketKey []byte

// A PacketQualifier is a set of packet fields as specified by a FieldList.
// Note that each field is padded to its maximum size. If the packet does not
// contain a field in the FieldList, the field is not added to the qualifier.
type PacketQualifier map[fwdpacket.FieldID][]byte

// An EntryKey describes packet fields in a FlowDesc which can be matched to
// specific masked values. An EntryKey is formed by concatenating the value
// and mask of the packet fields in the sequence described by a FieldList.
// Note that each field is padded to its maximum size. If the packet does not
// contain a field in the FieldList, the value is set to zero.
type EntryKey struct {
	value []byte // sequence of bytes describing the packet fields
	mask  []byte // sequence of bytes describing the mask of valid bits
}

// String returns the entry key as a formatted string.
func (k *EntryKey) String() string {
	return fmt.Sprintf("<byte=%x>;<mask=%x>;", k.value, k.mask)
}

// Equal returns true if two flow keys are equal.
func (k *EntryKey) Equal(k2 *EntryKey) bool {
	if bytes.Equal(k.value, k2.value) && bytes.Equal(k.mask, k2.mask) {
		return true
	}
	return false
}

// Match returns true if a PacketKey matches the EntryKey.
// Note:
//   - If the length of the EntryKey is 0, then all packets match the key.
//   - The EntryKey and PacketKey are assumed to be formed from the same
//     FieldList, non-existing fields are added as zeroes. Hence the length
//     of the EntryKey must be equal to the length of the PacketKey.
func (k *EntryKey) Match(p PacketKey) bool {
	if len(p) != len(k.value) {
		return false
	}
	for pos := range p {
		if p[pos]&k.mask[pos] != k.value[pos]&k.mask[pos] {
			return false
		}
	}
	return true
}

// An EntryQualifier describes packet fields in a FlowDesc which can be matched
// to values in the specified forwarding set (fwdset.Set). Note that the
// qualifier holds references to the forwarding sets.
type EntryQualifier struct {
	values map[fwdpacket.FieldID]*fwdset.Set
}

// String returns the entry qualifier as a formatted string.
func (q *EntryQualifier) String() string {
	index := 0
	buffer := make([]string, len(q.values))
	for id, set := range q.values {
		buffer[index] = fmt.Sprintf("ID=%v,Set=%v", id, set.ID())
		index++
	}
	return strings.Join(buffer, ";")
}

// Cleanup cleans up any references in the qualifier.
func (q *EntryQualifier) Cleanup() {
	for id, set := range q.values {
		set.Release(false /*forceCleanup*/)
		delete(q.values, id)
	}
	q.values = nil
}

// Equal returns true if two qualifiers are equal. Two qualifiers are set to
// be equal if they contain the exact same set of fields, and each field is
// associated to the same set. Note that the contents of the sets are not
// matched.
func (q *EntryQualifier) Equal(q2 *EntryQualifier) bool {
	if len(q.values) != len(q2.values) {
		return false
	}
	for f1, v1 := range q.values {
		if v2, ok := q2.values[f1]; !ok || v1.ID() != v2.ID() {
			return false
		}
	}
	return true
}

// Match returns true if a PacketQualifier matches the EntryQualifier.
// Note:
//   - If the length of the EntryQualifier is 0, then all packets match it.
//   - When creating qualifiers from a FieldList, non-existing fields are not
//     added. Hence the length of the qualifiers may not match.
func (q *EntryQualifier) Match(p PacketQualifier) bool {
	for f, s := range q.values {
		if v, ok := p[f]; !ok || !s.Contains(v) {
			return false
		}
	}
	return true
}

// A FieldList is a sequence of field ids.
type FieldList struct {
	fields []fwdpacket.FieldID // sequence of fields
	size   int                 // number of bytes needed by the fields after padding
}

// NewFieldList creates a new FieldList from a set of fields.
func NewFieldList(fm map[fwdpacket.FieldID]bool) *FieldList {
	size := 0
	var fields []fwdpacket.FieldID
	for id := range fm {
		size += fwdpacket.MaxSize(id)
		fields = append(fields, id)
	}
	return &FieldList{
		fields: fields,
		size:   size,
	}
}

// String returns a FieldList formatted as a string.
func (l *FieldList) String() string {
	return fmt.Sprintf("%v", l.fields)
}

// Size returns the number of bytes required for all fields in the FieldList.
func (l *FieldList) Size() int {
	return l.size
}

// MakeEntryKey creates an EntryKey for a FlowDesc. The EntryKey is formed by
// concatenating the value and mask of the fieldKey in the FlowDesc in the
// sequence specified by the FieldList.
func (l *FieldList) MakeEntryKey(fd *Desc) *EntryKey {
	k := EntryKey{
		value: make([]byte, 0, l.size),
		mask:  make([]byte, 0, l.size),
	}
	for _, id := range l.fields {
		value := []byte{}
		mask := []byte{}
		if field, ok := fd.keys[id]; ok {
			value = field.value
			mask = field.mask
		}
		value = fwdpacket.Pad(id, value)
		mask = fwdpacket.Pad(id, mask)
		k.value = append(k.value, value...)
		k.mask = append(k.mask, mask...)
	}
	return &k
}

// MakeEntryQualifier creates an EntryQualifier for a FlowDesc. The
// EntryQualifier is formed by referencing the Set specified in the
// FlowDesc for each field in FieldList.
func (l *FieldList) MakeEntryQualifier(fd *Desc) (*EntryQualifier, error) {
	q := EntryQualifier{
		values: make(map[fwdpacket.FieldID]*fwdset.Set),
	}

	for _, id := range l.fields {
		sid, ok := fd.qualifiers[id]
		if !ok {
			continue
		}
		set, err := fwdset.Find(fd.ctx, sid)
		if err != nil {
			return nil, fmt.Errorf("flow: NewEntryQualifier failed, cannot find set, %v", err)
		}
		if err := set.Acquire(); err != nil {
			return nil, fmt.Errorf("flow: NewEntryQualifier failed, cannot acquire set, %v", err)
		}
		q.values[id] = set
	}
	return &q, nil
}

// MakePacketKey makes a PacketKey for the specified packet.
func (l *FieldList) MakePacketKey(packet fwdpacket.Packet) PacketKey {
	var k PacketKey
	for _, id := range l.fields {
		field, err := packet.Field(id)
		if err != nil {
			field = []byte{}
		}
		field = fwdpacket.Pad(id, field)
		k = append(k, field...)
	}
	return k
}

// MakePacketQualifier makes a PacketQualifier for the specified packet.
func (l *FieldList) MakePacketQualifier(packet fwdpacket.Packet) PacketQualifier {
	q := make(PacketQualifier)
	for _, id := range l.fields {
		if field, err := packet.Field(id); err == nil {
			q[id] = fwdpacket.Pad(id, field)
		}
	}
	return q
}

// A fieldKey describes a single field as a value and mask.
type fieldKey struct {
	value []byte
	mask  []byte
}

// A Desc describes a packet flow by specifying keys and qualifiers.
// keys specify values of packet fields by describing the bits in them.
// qualifiers specify values of packet fields as a value in a fwdset.Set.
type Desc struct {
	keys       map[fwdpacket.FieldID]fieldKey
	qualifiers map[fwdpacket.FieldID]*fwdpb.SetId
	hash       uint32
	ctx        *fwdcontext.Context
}

// NewDesc creates a new Desc corresponding to the specified keys and qualifiers.
func NewDesc(ctx *fwdcontext.Context, keys []*fwdpb.PacketFieldMaskedBytes, qualifiers []*fwdpb.PacketFieldSet) (*Desc, error) {
	hash := fnv.New32()
	fd := &Desc{
		keys:       make(map[fwdpacket.FieldID]fieldKey),
		qualifiers: make(map[fwdpacket.FieldID]*fwdpb.SetId),
		ctx:        ctx,
	}
	for _, d := range keys {
		b := d.GetBytes()
		m := d.GetMasks()

		// Build the field descriptor with the values and mask. Copy
		// the bytes to all the desc to be garbage collected. If a mask
		// is specified, ensure it is the same size as the value. If a
		// mask is not specified, use all the bits of the field.
		field := fieldKey{
			value: make([]byte, len(b)),
			mask:  make([]byte, len(b)),
		}
		copy(field.value, b)
		if len(m) != 0 {
			if len(b) != len(m) {
				return nil, fmt.Errorf("flow: NewFlowDesc failed, len(bytes)=%v != len(masks)=%v", len(b), len(m))
			}
			copy(field.mask, m)
		} else {
			for pos := range b {
				field.mask[pos] = 0xFF
			}
		}
		hash.Write(field.value)
		hash.Write(field.mask)
		fd.keys[fwdpacket.NewFieldID(d.GetFieldId())] = field
	}

	for _, d := range qualifiers {
		setid := d.GetSetId()
		fd.qualifiers[fwdpacket.NewFieldID(d.GetFieldId())] = setid
	}
	fd.hash = hash.Sum32()
	return fd, nil
}

// Fields returns the list of fields specified in the flow desc's keys and
// qualifiers.
func (d *Desc) Fields() (keys []fwdpacket.FieldID, qualifiers []fwdpacket.FieldID) {
	for id := range d.keys {
		keys = append(keys, id)
	}
	for id := range d.qualifiers {
		qualifiers = append(qualifiers, id)
	}
	return keys, qualifiers
}

// String formats a flow descriptor as a string.
func (d *Desc) String() string {
	index := 0
	buffer := make([]string, len(d.keys)+len(d.qualifiers))
	for pos, v := range d.keys {
		buffer[index] = fmt.Sprintf("ID=%v,Value=%x,Mask=%x", pos, v.value, v.mask)
		index++
	}
	for pos, set := range d.qualifiers {
		buffer[index] = fmt.Sprintf("ID=%v,Set=%v", pos, set)
		index++
	}
	return strings.Join(buffer, ";")
}

// Equal returns true if two flow desc are equal.
func (d *Desc) Equal(d2 *Desc) bool {
	if (len(d.keys) != len(d2.keys)) || (len(d.qualifiers) != len(d2.qualifiers)) {
		return false
	}
	for f1, v1 := range d.keys {
		if v2, ok := d2.keys[f1]; !ok || !bytes.Equal(v1.value, v2.value) || !bytes.Equal(v1.mask, v2.mask) {
			return false
		}
	}
	for f1, q1 := range d.qualifiers {
		if q2, ok := d2.qualifiers[f1]; !ok || !proto.Equal(q1, q2) {
			return false
		}
	}
	return true
}

// Hash returns a hash value for the flow desc.
func (d *Desc) Hash() uint32 {
	return d.hash
}
