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
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/golang/glog"
	log "github.com/golang/glog"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdattribute"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Tracking dirty headers in packets.
//
// Rationale:
// Some header fields are implicit; they are computed using the values of other
// fields in the same header or in different headers of the packet. Hence when
// the packet frame is rebuilt, these values must be recomputed. Some examples
// of these in the TCP stack are:
// 1. TCP over IPv4 and IPv6 has a checksum computed over the TCP header, Pseudo IP
//    header and payload.
// 2. UDP over IPv6 has a checksum computed over the UDP header, Pseudo IP header and
//    payload
// 3. ICMPv4 has a checksum computed over the ICMP message.
// 4. ICMPv6 has a checksum computed over the ICMP message and the Pseudo IP header.
// 5. The length field and next header field in every packet header.
// Updates such as [5] are very inexpensive, while [1] and [2] are expensive
// when considering messages like BGP, gRPC that have large payloads that
// do not change in flight (usually). Tracking dirty headers in the packet,
// allows a header to determine if it needs to recompute a field depending
// on whether itself or the headers it depends on are dirty. This mechanism is
// used to avoid costly recomputations such as [1]...[4]. Inexpensive updates
// such as [5] are not avoided to keep things simple.
//
// Mechanism:
// 1. When a field in a header is updated, the header reports back if the header
//    was dirtied. The infrastructure uses this flag to mark the header as dirty.
// 2. When a header is added or removed, the header is marked dirty as it needs
//    to compute implicit fields. Note that "removing a header" may just remove
//    a part of the header (such as a vlan tag) while leaving the header (such as
//    Ethernet) in place.
// 3. Headers are marked as clean only once the entire packet has been rebuilt.
//
// Note: The IP header is always modified on each hop as the TTL is decremented.
// The IP TTL decrement code always adjusts the IP checksum and reports the IP
// header as clean. This is essential to avoid recomputing TCP/UDP checksums
// on each hop.

// A Packet is a network packet that can be queried and manipulated.
// A packet is created by parsing a frame to create a chain of headers.
type Packet struct {
	headers    []*Desc              // Descriptors for each header in the packet
	debug      bool                 // true if the packet is being debugged
	desc       string               // describes the packet in human readable form
	attributes fwdattribute.Set     // set of attributes associated with the packet
	start      fwdpb.PacketHeaderId // Start header of the packet
	loggerMu   sync.Mutex
	logger     logr.Logger
	logSync    *packetLogger
}

// fieldDesc returns the Desc of the packet and the corresponding field id that
// contains the specified field.
//
// In case of non-UDF, the specified field id is returned unmodified. It may
// panic if a field unknown to the parser is specified.
//
// In case UDF, it searches for the field in the packet starting from the
// specifiedi desc. It then returns the desc and the corresponding field id. It may panic if
// a header unknown to the parser is specified.
func (p *Packet) fieldDesc(id fwdpacket.FieldID) (*Desc, fwdpacket.FieldID) {
	switch id.IsUDF {
	case true:
		// In case of UDF, we need to search successive headers that
		// may contain the specified field and return an updated field
		// id.
		h, ok := GroupAttr[id.Header]
		if !ok {
			panic(fmt.Sprintf("protocol: fieldDesc failed, field %+v contains unknown header", id))
		}
		for pos := h.Position; pos < len(Sequence); pos++ {
			hid := Sequence[pos]
			id.Header = hid // Use the specified header for the UDF
			if d := p.headers[hid]; d != nil {
				// If the field is found in the current header
				// use it. Otherwise update the field id.
				bl := uint8(len(d.handler.Header()))
				if id.Offset < bl {
					return d, id
				}
				id.Offset -= bl
			}
		}
		return nil, id

	case false:
		// In case of well known fields, we can find the header directly and use the
		// unmodified field id.
		if attr, ok := FieldAttr[id.Num]; ok {
			return p.headers[attr.Group], id
		}
		panic(fmt.Sprintf("protocol: fieldDesc failed, field %+v contains unknown field number", id))
	}
	return nil, id
}

// rebuildHeaders rebuilds the headers preceeding the specified header.
// The rebuilt headers are marked as clean.
func (p *Packet) rebuildHeaders(header *Desc) {
	for header != nil {
		header.handler.Rebuild()
		header = header.prev
	}
}

// String returns a string representation of the packet.
func (p *Packet) String() string {
	var buf bytes.Buffer
	for h := p.headers[Sequence[0]]; h != nil; h = h.next {
		fmt.Fprintf(&buf, "%v->%x/%x; ", h.group, h.handler.Header(), h.handler.Trailer())
	}
	return buf.String()
}

// Attributes returns the attributes associated with the packet.
func (p *Packet) Attributes() fwdattribute.Set {
	return p.attributes
}

// Length returns the number of bytes in the packet.
func (p *Packet) Length() int {
	length, err := p.Field(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_LENGTH, 0))
	if err != nil {
		log.Errorf("Unable to get length for packet %v, err %v.", p, err)
		return 0
	}
	return int(binary.BigEndian.Uint64(length))
}

// Frame returns the packet's frame as a slice of bytes. The full frame of the
// packet is the payload of the first header which is METADATA by default.
func (p *Packet) Frame() []byte {
	if d := p.headers[Sequence[0]]; d != nil {
		return d.Payload()
	}
	return []byte{}
}

// Field returns the bytes associated with a field ID.
func (p *Packet) Field(id fwdpacket.FieldID) ([]byte, error) {
	header, id := p.fieldDesc(id)
	if header == nil {
		return nil, fmt.Errorf("Field %v failed for packet %v, Header does not exist", id, p)
	}
	field, err := header.handler.Field(id)
	if err != nil {
		return nil, fmt.Errorf("Field %v failed for packet %v, err %v", id, p, err)
	}
	return field, nil
}

// Update updates a field in the packet.
func (p *Packet) Update(id fwdpacket.FieldID, op int, arg []byte) error {
	header, id := p.fieldDesc(id)
	if header == nil {
		return fmt.Errorf("Update operation %v on %v failed for packet %v, Header does not exist", op, id, p)
	}
	dirty, err := header.handler.UpdateField(id, op, arg)
	if err != nil {
		return fmt.Errorf("Update operation %v on %v failed for packet %v, err %v", op, id, p, err)
	}
	if dirty {
		header.MarkDirty(true)
	}
	return nil
}

// Debug controls debugging for the packet.
func (p *Packet) Debug(enable bool) {
	p.debug = enable
}

var _ logr.LogSink = &packetLogger{}

type packetLogger struct {
	funcr.Formatter
	msgs []string
}

// Init receives optional information about the logr library for LogSink
// implementations that need it.
func (pl packetLogger) Init(info logr.RuntimeInfo) {}

// Enabled tests whether this LogSink is enabled at the specified V-level.
// For example, commandline flags might be used to set the logging
// verbosity and disable some info logs.
func (pl packetLogger) Enabled(level int) bool {
	return bool(glog.V(glog.Level(level)))
}

// Info logs a non-error message with the given key/value pairs as context.
// The level argument is provided for optional logging.  This method will
// only be called when Enabled(level) is true. See Logger.Info for more
// details.
func (pl *packetLogger) Info(level int, msg string, keysAndValues ...interface{}) {
	prefix, arg := pl.FormatInfo(level, msg, keysAndValues)
	if prefix == "" {
		pl.msgs = append(pl.msgs, arg)
	} else {
		pl.msgs = append(pl.msgs, prefix+" "+arg)
	}
}

// Error logs an error, with the given message and key/value pairs as
// context.  See Logger.Error for more details.
func (pl *packetLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	prefix, arg := pl.FormatError(err, msg, keysAndValues)
	if prefix == "" {
		pl.msgs = append(pl.msgs, arg)
	} else {
		pl.msgs = append(pl.msgs, prefix+" "+arg)
	}
}

// WithValues returns a new LogSink with additional key/value pairs.  See
// Logger.WithValues for more details.
func (pl packetLogger) WithValues(keysAndValues ...interface{}) logr.LogSink {
	pl.Formatter.AddValues(keysAndValues)
	return &pl
}

// WithName returns a new LogSink with the specified name appended.  See
// Logger.WithName for more details.
func (pl packetLogger) WithName(name string) logr.LogSink {
	pl.Formatter.AddName(name)
	return &pl
}

// Log returns a logger.
func (p *Packet) Log() logr.Logger {
	return p.logger
}

// LogMsgs returns the log messages for the packet.
func (p *Packet) LogMsgs() []string {
	return p.logSync.msgs
}

// NewPacket parses a frame into a Packet and returns it.
func NewPacket(start fwdpb.PacketHeaderId, frame *frame.Frame) (*Packet, error) {
	p := &Packet{
		headers:    make([]*Desc, len(Sequence)),
		desc:       "none",
		attributes: fwdattribute.NewSet(),
		start:      start,
	}

	sync := &packetLogger{
		Formatter: funcr.NewFormatter(funcr.Options{
			RenderArgsHook: func(kvList []interface{}) []interface{} {
				for i, kv := range kvList {
					if _, ok := kv.(fwdpacket.LogFrameValue); ok {
						kvList[i] = fmt.Sprintf("%x", p.Frame())
					}
				}
				return kvList
			},
		}),
	}
	p.logSync = sync
	p.logger = logr.New(p.logSync)

	// Start parsing the packet using the first Header. The metadata is a
	// lucius only special header. It is always followed by the specified
	// first header in the frame.
	var prev *Desc
	currID := fwdpb.PacketHeaderId_PACKET_HEADER_ID_METADATA
	for currID != fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE {
		// Parse the current frame .
		attr, ok := HeaderAttr[currID]
		if !ok || attr.Parse == nil {
			return nil, fmt.Errorf("Parse failed for packet %v, bad Header %v", frame, currID)
		}
		header := &Desc{
			group:  attr.Group,
			Packet: p,
			prev:   prev,
		}
		handler, nextID, err := attr.Parse(frame, header)
		if err != nil {
			return nil, fmt.Errorf("Parse failed for packet %v, err %v", frame, err)
		}

		// Adjust the returned next header if the current header is METADATA.
		if currID == fwdpb.PacketHeaderId_PACKET_HEADER_ID_METADATA {
			nextID = p.start
		}

		// Add the desc and update the chain of Headers.
		header.handler = handler
		p.headers[attr.Group] = header
		if prev != nil {
			prev.next = header
		}
		prev = header
		currID = nextID
	}
	return p, nil
}

// Decap removes the header specified by id. The header is marked as dirty.
func (p *Packet) Decap(id fwdpb.PacketHeaderId) error {
	attr, ok := HeaderAttr[id]
	if !ok {
		return fmt.Errorf("Decap %v failed for packet %v, unknown header", id, p)
	}
	header := p.headers[attr.Group]
	if header == nil {
		return fmt.Errorf("Decap %v failed for packet %v, non existing header %v", id, p, attr.Group)
	}
	if err := header.handler.Remove(id); err != nil {
		return fmt.Errorf("Decap %v failed for packet %v, err %v", id, p, err)
	}
	header.MarkDirty(true)

	// If the header is now empty, adjust the chain of Headers.
	if len(header.handler.Header()) == 0 {
		if header.next != nil {
			header.next.prev = header.prev
		}
		if header.prev != nil {
			header.prev.next = header.next
		}
		p.headers[attr.Group] = nil
	}
	p.rebuildHeaders(header.prev)
	return nil
}

// Encap adds a Header specified by id. If a header already exists in the
// header-group, then the existing header is modified. The header being added
// or updated is always marked as dirty.
func (p *Packet) Encap(id fwdpb.PacketHeaderId) error {
	attr, ok := HeaderAttr[id]
	if !ok {
		return fmt.Errorf("Encap header %v failed for packet %v, unknown header", id, p)
	}

	// Add the encap to an existing Header in the group.
	if header := p.headers[attr.Group]; header != nil {
		if err := header.handler.Modify(id); err != nil {
			return fmt.Errorf("Encap header %v failed for packet %v, err %v", id, p, err)
		}
		header.MarkDirty(true)
		p.rebuildHeaders(header.prev)
		return nil
	}

	// Check if the header can be added
	if attr.Add == nil {
		return fmt.Errorf("Encap header %v failed for packet %v, cannot add header", id, p)
	}

	// Find the position of the new Header in the frame and add it.
	group := GroupAttr[attr.Group]
	if group.Position <= 0 {
		return fmt.Errorf("Encap %v failed for packet %v, cannot add header at position %v", id, p, group.Position)
	}
	var prev *Desc
	for pos := group.Position - 1; pos >= 0 && prev == nil; pos-- {
		prev = p.headers[Sequence[pos]]
	}
	if prev == nil {
		return fmt.Errorf("Encap %v failed for packet %v, malformed packet", id, p)
	}

	header := &Desc{
		group:  attr.Group,
		Packet: p,
		prev:   prev,
		next:   prev.next,
	}
	var err error
	if header.handler, err = attr.Add(id, header); err != nil {
		return fmt.Errorf("Encap %v failed for packet %v, err %v", id, p, err)
	}
	if prev.next != nil {
		prev.next.prev = header
	}
	prev.next = header
	p.headers[attr.Group] = header
	header.MarkDirty(true)
	p.rebuildHeaders(header)
	return nil
}

// clone clones a packet with the specified prepended bytes and copies the
// specified fields from the original packet to the clone. Note that by
// default fields not encoded in the frame (such as metadata fields) are lost.
// The clone is expected to have the specified start packet header.
// If the replicate flag is true, we force the clone to have its own byte
// array. This allows the clone and original to be updated independently and
// incurs the cost of memory allocation. In some cases (such as reparse), the
// original packet is expected to be discarded, and hence we can let Go
// optimize the slice/array management. Likewise if the packet is being
// replicated, it has its own packet log. If not, the packet log is carried
// over from the original packet to the new packet.
func (p *Packet) clone(replicate bool, prepend []byte, id fwdpb.PacketHeaderId, fields []fwdpacket.FieldID) (*Packet, error) {
	// Get the packet's frame. This also causes the dirty headers to be
	// rebuilt.
	of := p.Frame()

	// Copy over the fields from the original to the cloned packet. This is
	// done after the call to Frame() on the original packet which ensures
	// that all headers are rebuilt before we make a copy of their fields.
	saved := make(map[fwdpacket.FieldID]frame.Field)
	for _, k := range fields {
		v, err := p.Field(k)
		if err != nil {
			return nil, fmt.Errorf("clone failed to query field %v, err %v", k, err)
		}
		saved[k] = v
	}
	if replicate {
		cp := make([]byte, len(of))
		copy(cp, of)
		of = cp
	}
	of = append(prepend, of...)
	np, err := NewPacket(id, frame.NewFrame(of))
	if err != nil {
		return nil, fmt.Errorf("clone failed to parse frame %x (start %v), err %v", of, id, err)
	}
	np.debug = p.debug
	np.desc = p.desc
	np.attributes = p.attributes
	np.logSync = p.logSync
	np.logger = logr.New(np.logSync)
	if !replicate {
		np.logSync.msgs = p.logSync.msgs
	}

	// Restore the saved values into the cloned packet.
	for k, v := range saved {
		if err := np.Update(k, fwdpacket.OpSet, v); err != nil {
			return nil, fmt.Errorf("clone with header %v failed to restore field %v to value %x, err %v", id, k, v, err)
		}
	}
	return np, nil
}

// Mirror creates a new packet from the current packet. By default the metadata
// fields are lost. Additional fields specified during the mirror ensures that
// the field values are copied from the old packet to the new packet. This can
// be used to retain  metadata field values across a mirror.
func (p *Packet) Mirror(fields []fwdpacket.FieldID) (fwdpacket.Packet, error) {
	return p.clone(true, nil, p.start, fields)
}

// Reparse reparses the current packet to start from the specified packet header
// id. Note that by default reparsing creates a new packet, so metadata fields
// will be lost. Additional fields specified during reparsing ensures that the
// field values are copied from the old packet to the new packet. This can be
// used to retain metadata field values across a packet reparse. It can also
// prepend a slice of bytes to the rebuilt packet before reparsing.
func (p *Packet) Reparse(id fwdpb.PacketHeaderId, fields []fwdpacket.FieldID, prepend []byte) error {
	np, err := p.clone(false, prepend, id, fields)
	if err != nil {
		return err
	}

	// Copy over the packet.
	*p = *np
	return nil
}

// StartHeader returns the start header of the packet.
func (p *Packet) StartHeader() fwdpb.PacketHeaderId {
	return p.start
}
