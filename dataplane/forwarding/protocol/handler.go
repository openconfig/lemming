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

package protocol

import (
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A Handler encapsulates a protocol header, providing methods to query
// and mutate it. Note that some protocol headers are complex i.e they
// can have multiple instances of themselves (IP-IP tunnels) or can be
// composed of multiple headers (ethernet with vlans).
type Handler interface {
	// Header returns bytes of the current protocol header.
	Header() []byte

	// Trailer returns the trailing bytes of the current protocol header.
	// This is typically used for padding or trailing CRC.
	Trailer() []byte

	// ID returns the protocol header ID of the specified instance.
	ID(instance int) fwdpb.PacketHeaderId

	// Field finds the field specified by id.
	Field(id fwdpacket.FieldID) ([]byte, error)

	// UpdateField updates the field specified by id. The type of update
	// is determined by op using the constants defined in fwdpacket.
	// It returns true if the update dirties the header.
	UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error)

	// Remove removes the protocol header specified by id from the
	// current header.
	Remove(id fwdpb.PacketHeaderId) error

	// Modify extends the current header with the header specified
	// by id.
	Modify(id fwdpb.PacketHeaderId) error

	// Rebuild rebuilds the current header
	Rebuild() error
}

// A Desc is a set of attributes that describe a network protocol in a packet.
// The desc in the packet are maintained in a doubly linked list in the
// order in which they occur within the packet. By default, all protocol
// headers are marked as clean.
type Desc struct {
	group   fwdpb.PacketHeaderGroup // Group of the header.
	next    *Desc                   // Next header in the packet.
	prev    *Desc                   // Previous header in the packet.
	handler Handler                 // Handler for the current packet.
	Packet  *Packet                 // Reference to the packet.
	dirty   bool                    // true if the network protocol is dirty.
}

// PayloadLength returns the length of the payload of a protocol in the packet.
// The payload of a desc is represented by the set of all subsequent desc
// in the packet (header and trailer).
func (d *Desc) PayloadLength() int {
	var l int
	for desc := d.next; desc != nil; desc = desc.next {
		l += len(desc.handler.Header()) + len(desc.handler.Trailer())
	}
	return l
}

// PayloadID returns the payload's header id. The payload of a desc is the
// header of the first instance of the succeeding desc.
func (d *Desc) PayloadID() fwdpb.PacketHeaderId {
	if d.next == nil {
		return fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE
	}
	return d.next.handler.ID(fwdpacket.FirstField)
}

// PayloadDesc returns the payload's desc.
func (d *Desc) PayloadDesc() *Desc {
	return d.next
}

// Payload returns the payload. The payload of a desc is represented by the
// set of all subsequent desc in the packet. handlers are rebuilt and the
// header and trailer bytes of each protocol handler is combined to form the
// frame. All headers are marked as clean after being rebuilt.
func (d *Desc) Payload() []byte {
	var headers [][]byte
	var trailers [][]byte

	// setup appends a copy of entry to a list if the entry is
	// non-nil. Otherwise it creates an empty entry at the same position
	// in the list.
	setup := func(list [][]byte, entry []byte) [][]byte {
		if len(entry) == 0 {
			return append(list, nil)
		}
		t := make([]byte, len(entry))
		copy(t, entry)
		return append(list, t)
	}

	// Due to the setup headers and trailers are arrays of equal length.
	// This ensures that the headers/trailers of the same protocol handler
	// are present at the same position in the corresponding arrays.
	for h := d.next; h != nil; h = h.next {
		h.handler.Rebuild()
		headers = setup(headers, h.handler.Header())
		trailers = setup(trailers, h.handler.Trailer())
	}

	var payload []byte
	for pos := len(headers) - 1; pos >= 0; pos-- {
		if h := headers[pos]; h != nil {
			payload = append(h, payload...)
		}
		if t := trailers[pos]; t != nil {
			payload = append(payload, t...)
		}
	}

	for h := d.next; h != nil; h = h.next {
		h.MarkDirty(false)
	}

	return payload
}

// EnvelopeID returns the encapsulator's header id. The envelope of a desc is
// the header of the last instance of the preceding desc.
func (d *Desc) EnvelopeID() fwdpb.PacketHeaderId {
	if d.prev == nil {
		return fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE
	}
	return d.prev.handler.ID(fwdpacket.LastField)
}

// EnvelopeDesc returns the envelope's desc.
func (d *Desc) EnvelopeDesc() *Desc {
	return d.prev
}

// Dirty returns true if the header is dirty.
func (d *Desc) Dirty() bool {
	return d.dirty
}

// MarkDirty marks the header as dirty or clean.
func (d *Desc) MarkDirty(dirty bool) {
	d.dirty = dirty
}
