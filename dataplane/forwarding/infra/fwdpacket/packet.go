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

// Package fwdpacket contains routines and types for manipulating packets.
package fwdpacket

import (
	"encoding/binary"
	"fmt"
	"sync"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdattribute"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Operations that can be applied to a field.
const (
	OpInc = iota // Increment a field with the specified value.
	OpDec        // Decrement a field with the specified value.
	OpSet        // Set a field to the specified value.
)

// Enumeration of various operations that can be performed on the packet log.
const (
	LogDebugMessage = iota // log a debug message
	LogDebugFrame          // log a debug message with the frame
	LogErrorFrame          // log an error message with the frame
	LogDesc                // sets the packet's description (prefixed to all messages)
)

// AttrPacketDebug controls packet debugging.
var AttrPacketDebug = fwdattribute.ID("PacketDebug")

func init() {
	fwdattribute.Register(AttrPacketDebug, "Enables packet debugging if set to true")
}

// Packet is a network packet that can be queried and manipulated.
//
// Packet repesents a network packet that has been parsed and verfied.
// It supports queries on various fields. Additionally it can change the
// network packet by adding and removing headers or setting packet fields.
type Packet interface {
	// Field returns the bytes associated with a field ID.
	Field(id FieldID) ([]byte, error)

	// Update updates a field in the packet.
	Update(id FieldID, op int, arg []byte) error

	// Decap removes the outermost header of the specified type.
	Decap(id fwdpb.PacketHeaderId) error

	// Encap adds an outermost header of the specified type.
	Encap(id fwdpb.PacketHeaderId) error

	// Reparse reparses the current packet to start from the specified packet header
	// id. Note that the current packet is rebuilt before it is reparsed. Note that
	// by default reparsing creates a new packet, so metadata fields will be lost.
	// Additional fields specified during reparsing ensures that the field values
	// are copied from the old packet to the new packet. This can be used to retain
	// metadata field values across a packet reparse. It can also prepend the
	// rebuilt packet with the specified set of bytes before reparsing the packet.
	Reparse(id fwdpb.PacketHeaderId, fields []FieldID, prepend []byte) error

	// Mirror creates a new packet from the current packet. Note that the current
	// packet is rebuilt before it is mirrored. Note that by default the metadata
	// fields are lost. Additional fields specified during the mirror ensures that
	// the field values are copied from the old packet to the new packet. This can
	// be used to retain  metadata field values across a mirror.
	Mirror(fields []FieldID) (Packet, error)

	// String formats the packet into a string.
	String() string

	// Length returns the length of the packet in bytes.
	Length() int

	// Frame returns the packet as a slice of bytes.
	Frame() []byte

	// Debug control debugging for the packet.
	Debug(enable bool)

	// Logf controls the packet's message log.
	Logf(level int, fmt string, args ...interface{})

	// Log returns the contents of the packet's log.
	Log() []string

	// Attributes returns the attributes associated with the packet.
	Attributes() fwdattribute.Set

	// StartHeader returns the first header of the packet.
	StartHeader() fwdpb.PacketHeaderId
}

// Parser can create a Packet from a slice of bytes containing the packet.
type Parser interface {
	// New creates a new packet using the specified first header id.
	New(hid fwdpb.PacketHeaderId, bytes []byte) (Packet, error)

	// Validate validates the size of the specified field.
	Validate(id FieldID, size int) error

	// MaxSize returns the maximum size of the specified field.
	MaxSize(id FieldID) int
}

// parser is the Parser used in Lucius.
var parser Parser

// Register registers a parser for Lucius.
func Register(p Parser) {
	parser = p
}

// New builds a new packet from a stream of bytes.
func New(hid fwdpb.PacketHeaderId, bytes []byte) (Packet, error) {
	if parser == nil {
		return nil, fmt.Errorf("fwdpacket: Cannot create packet, no parser registered")
	}
	return parser.New(hid, bytes)
}

// NewNID creates a new packet from the specified frame and sets up its input port
// if specified.
func NewNID(hid fwdpb.PacketHeaderId, frame []byte, nid fwdobject.NID) (Packet, error) {
	pkt, err := New(hid, frame)
	if err != nil {
		return nil, err
	}
	if nid != fwdobject.InvalidNID {
		if err := SetNID(pkt, nid, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT); err != nil {
			return nil, err
		}
	}
	return pkt, nil
}

// Validate validates the size of the specified field.
func Validate(id FieldID, size int) error {
	if parser == nil {
		return fmt.Errorf("fwdpacket: Validate field %v failed, no parser registered", id)
	}
	return parser.Validate(id, size)
}

// MaxSize returns the maximum size of the specified field.
func MaxSize(id FieldID) int {
	if parser == nil {
		return 0
	}
	return parser.MaxSize(id)
}

// Pad pads the field to the maximum size.
func Pad(id FieldID, field []byte) []byte {
	curr := len(field)
	max := MaxSize(id)
	if curr > max {
		log.Errorf("fwdpacket: Pad for field %v failed, length %v is greater than maximum size %v.", id, curr, max)
		return field
	}
	return frame.Resize(field, max)
}

// Truncate truncates a padded field to the specified size.
func Truncate(field []byte, size int) []byte {
	curr := len(field)
	if curr < size {
		log.Errorf("fwdpacket: Truncate for field %v failed, length %v is smaller than required size %v.", field, curr, size)
		return field
	}
	return frame.Resize(field, size)
}

// SetNID sets a NID in the specified field.
func SetNID(packet Packet, nid fwdobject.NID, num fwdpb.PacketFieldNum) error {
	fid := NewFieldIDFromNum(num, 0)
	pid := make([]byte, MaxSize(fid))
	binary.BigEndian.PutUint64(pid, uint64(nid))
	return packet.Update(fid, OpSet, pid)
}

// global lock for printing the packet log.
var logMu sync.Mutex

// Log logs the packet's buffer if it is not empty.
func Log(packet Packet) {
	if m := packet.Log(); len(m) != 0 {
		logMu.Lock()
		log.V(1).Infof("Packet Trace:\n")
		for _, msg := range m {
			log.V(1).Infof("%v\n", msg)
		}
		logMu.Unlock()
	}
}
