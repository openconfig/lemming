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

package mpls

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type mplsLabel struct {
	hdr frame.Header
}

type mpls struct {
	protocol.Handler
	desc   *protocol.Desc
	labels []*mplsLabel
}

func (m *mpls) Header() []byte {
	b := make([][]byte, len(m.labels))
	for p, h := range m.labels {
		b[p] = h.hdr
	}
	return bytes.Join(b, []byte{})
}

func (m *mpls) Trailer() []byte {
	return nil
}

func (m *mpls) ID(instance int) fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS
}

const (
	labelOffset    = 0
	labelByteLen   = 4  // To make dealing with label easier, encode it as uint32.
	labelBitOffset = 12 // The label is the leftmost 20 bits of the header, so shift 12 bits.
	labelBitLen    = 20
	tcOffset       = 2
	tcByteLen      = 1
	tcBitOffset    = 1
	tcBitLen       = 3
	ttlOffset      = 3
	ttlByteLen     = 1
	botOffset      = 2
	botByteLen     = 1
	botBitOffset   = 0
	botBitLen      = 1
)

func (m *mpls) field(id fwdpacket.FieldID) frame.Field {
	l := m.labels[id.Instance]
	if id.IsUDF {
		return protocol.UDF(l.hdr, id)
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_LABEL:
		return l.hdr.Field(labelOffset, labelByteLen)
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TC:
		return l.hdr.Field(tcOffset, tcByteLen)
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TTL:
		return l.hdr.Field(ttlOffset, ttlByteLen)
	default:
		return nil
	}
}

func (m *mpls) Field(id fwdpacket.FieldID) ([]byte, error) {
	field := m.field(id)
	if field == nil {
		return nil, fmt.Errorf("ethernet: Field failed, field %v does not exist", id)
	}
	if id.IsUDF {
		return field.Copy(), nil
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_LABEL:
		return field.BitField(labelBitOffset, labelBitLen), nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TC:
		return field.BitField(tcBitOffset, tcBitLen), nil
	default:
		return field.Copy(), nil
	}
}

func (m *mpls) UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	field := m.field(id)
	if field == nil {
		return false, fmt.Errorf("mpls: UpdateField failed, field %v does not exist", id)
	}
	if op == fwdpacket.OpDec && id.Num == fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TTL {
		pttl := uint8(field.Value())
		nttl := pttl - 1
		field.SetValue(uint(nttl))
		return false, nil
	}
	if op != fwdpacket.OpSet {
		return false, fmt.Errorf("only SET is supported")
	}

	if id.IsUDF && op == fwdpacket.OpSet {
		return true, field.Set(arg)
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_LABEL:
		field.SetBits(labelBitOffset, labelBitLen, uint64(binary.BigEndian.Uint32(arg)))
		return true, nil
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TC:
		field.SetBits(tcBitOffset, tcBitLen, uint64(arg[0]))
		return true, nil
	default:
		return true, field.Set(arg)
	}
}

func (m *mpls) Remove(id fwdpb.PacketHeaderId) error {
	if len(m.labels) > 1 {
		m.labels = m.labels[1:]
		return nil
	}
	m.labels = nil
	return nil
}

func (m *mpls) Modify(id fwdpb.PacketHeaderId) error {
	m.labels = append([]*mplsLabel{{hdr: make(frame.Header, 4)}}, m.labels...)
	return nil
}

func (m *mpls) Rebuild() error {
	if len(m.labels) == 0 {
		return nil
	}
	// Set the bottom of stack bit 0 on all but the last label. Set that it to 1.
	for i := 0; i < len(m.labels)-1; i++ {
		m.labels[i].hdr.Field(botOffset, botByteLen).SetBits(botBitOffset, botBitLen, 0)
	}
	m.labels[len(m.labels)-1].hdr.Field(botOffset, botByteLen).SetBits(botBitOffset, botBitLen, 1)
	return nil
}

const (
	ipv4ExplicitNull = 0
	ipv6ExplicitNull = 3
)

func parseMPLS(f *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	m := &mpls{
		desc: desc,
	}

	var label frame.Field
	for {
		hdr, err := f.ReadHeader(4)
		if err != nil {
			return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("failed to parse MPLS: %v", err)
		}
		botStack := hdr.Field(botOffset, botByteLen).BitField(botBitOffset, botBitLen) // bottom of stack is the last bit of the 3 byte.

		m.labels = append(m.labels, &mplsLabel{
			hdr: hdr,
		})
		label = hdr.Field(labelOffset, labelByteLen).BitField(labelBitOffset, labelBitLen)
		if botStack.Value() == 1 {
			break
		}
	}
	// If the last label is explicit null, then use that to parse the next header, otherwise treat the rest of the packet as opaque.
	next := fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE
	switch label.Value() {
	case ipv4ExplicitNull:
		next = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4
	case ipv6ExplicitNull:
		next = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6
	}

	return m, next, nil
}

func add(id fwdpb.PacketHeaderId, desc *protocol.Desc) (protocol.Handler, error) {
	if id != fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS {
		return nil, fmt.Errorf("unsupport type: %v", id)
	}
	return &mpls{
		desc:   desc,
		labels: []*mplsLabel{{hdr: make(frame.Header, 4)}},
	}, nil
}

func init() {
	protocol.Register(fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS, parseMPLS, add)
}
