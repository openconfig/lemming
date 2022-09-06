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

package ip

import (
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Size of various fields in a GRE header.
const (
	greDescBytes   = 1 // Number of bytes in the header description.
	greDescPos     = 0 // Offset in bytes of the header description.
	greSeqBits     = 1 // Number of bits for the sequence flag.
	greSeqPos      = 4 // Offset in bits for the sequence flag.
	greKeyBits     = 1 // Number of bits for the key flag.
	greKeyPos      = 5 // Offset in bits for the key flag.
	greSeqBytes    = 4 // Number of bytes for the gre sequence number.
	greKeyBytes    = 4 // Number of bytes for the gre key.
	greProtoBytes  = 2 // Number of bytes in the IP protocol field.
	greProtoPos    = 2 // Offset in bytes of the IP protocol.
	greHeaderBytes = 4 // Number of bytes in the fixed sized GRE header.
)

// A GRE represents a GRE header in the packet. It can be queried, updated
// added and removed. It supports GRE headers with the optional key and
// sequence-id attributes.
type GRE struct {
	header  frame.Header // Mandatory GRE header.
	key     frame.Field  // Optional GRE key field.
	seq     frame.Field  // Optional GRE sequence field.
	payload int64        // Length of payload.
}

// Header returns the GRE header as a slice of bytes.
func (gre *GRE) Header() []byte {
	b := append(gre.header, gre.key...)
	b = append(b, gre.seq...)
	return b
}

// ID returns the protocol header ID.
func (GRE) ID() fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE
}

// Payload gets the payload information.
func (gre *GRE) Payload() (fwdpb.PacketHeaderId, int64) {
	return gre.ID(), int64(len(gre.header)+len(gre.key)+len(gre.seq)) + gre.payload
}

// SetPayload sets the payload.
func (gre *GRE) SetPayload(id fwdpb.PacketHeaderId, length int64) {
	gre.payload = length
	next, ok := ethernet.HeaderNext[id]
	if !ok {
		next = ethernet.Reserved
	}
	gre.header.Field(greProtoPos, greProtoBytes).SetValue(uint(next))
}

// Find returns a copy of the field specified by id.
func (gre *GRE) Find(id fwdpacket.FieldID) ([]byte, error) {
	if id.IsUDF {
		return nil, fmt.Errorf("gre: Find failed, field %v does not exist", id)
	}
	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_GRE_KEY:
		if gre.key != nil {
			return gre.key.Copy(), nil
		}

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_GRE_SEQUENCE:
		if gre.seq != nil {
			return gre.seq.Copy(), nil
		}
	}
	return nil, fmt.Errorf("gre: Find failed, field %v does not exist", id)
}

// Update updates a slice of bytes identified by id.
func (gre *GRE) Update(id fwdpacket.FieldID, oper int, arg []byte) (bool, error) {
	if id.IsUDF {
		return false, fmt.Errorf("gre: Update failed, field %v is not supported", id)
	}
	if oper != fwdpacket.OpSet {
		return false, fmt.Errorf("gre: Update failed, operation %v is not supported", oper)
	}

	switch id.Num {
	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_GRE_KEY:
		if gre.key == nil {
			gre.key = make(frame.Field, greKeyBytes)
			gre.header.Field(greDescPos, greDescBytes).SetBits(greKeyPos, greKeyBits, 0x1)
		}
		gre.key.Set(arg)
		return true, nil

	case fwdpb.PacketFieldNum_PACKET_FIELD_NUM_GRE_SEQUENCE:
		if gre.seq == nil {
			gre.seq = make(frame.Field, greSeqBytes)
			gre.header.Field(greDescPos, greDescBytes).SetBits(greSeqPos, greSeqBits, 0x1)
		}
		gre.seq.Set(arg)
		return true, nil
	default:
		return false, fmt.Errorf("gre: Update failed, field %v is not supported", id)
	}
}

// newGRE creates an empty GRE header.
func newGRE() header {
	return &GRE{
		header: make(frame.Header, greHeaderBytes),
	}
}

// makeGRE parses a GRE header.
func makeGRE(f *frame.Frame) (header, fwdpb.PacketHeaderId, error) {
	desc, err := f.Peek(greDescPos, greDescBytes)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("gre: makeGRE failed, err %v", err)
	}

	keyPresent := false // Indicates if the gre key is present.
	seqPresent := false // Indicates if the gre sequence is present.
	if desc.BitField(greKeyPos, greKeyBits).Value() == 1 {
		keyPresent = true
	}
	if desc.BitField(greSeqPos, greSeqBits).Value() == 1 {
		seqPresent = true
	}

	// Read the GRE header followed by the optional gre key and sequence.
	header, err := f.ReadHeader(greHeaderBytes)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("gre: makeGRE failed to read header, err %v", err)
	}

	optional := func(size int, present bool, name string) (frame.Field, error) {
		if !present {
			return nil, nil
		}

		h, err := f.ReadHeader(size)
		if err != nil {
			return nil, fmt.Errorf("gre: makeGRE failed to read %v, err %v", name, err)
		}
		return h.Field(0, size), nil
	}

	key, err := optional(greKeyBytes, keyPresent, "key")
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, err
	}

	seq, err := optional(greSeqBytes, seqPresent, "sequence")
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, err
	}

	gre := &GRE{
		header:  header,
		key:     key,
		seq:     seq,
		payload: int64(f.Len()),
	}
	if next, ok := ethernet.NextHeader[uint16(header.Field(greProtoPos, greProtoBytes).Value())]; ok {
		return gre, next, nil
	}
	return gre, fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE, nil
}
