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

// Package ip handles the IP L3 portion of the packet. Lucius supports IP
// packets such as IPv4, IPv6 and IP tunnels such as GRE and IP-IP tunnels
// with IPv4, IPv6 payload and IPv4, IPv6 transport.
package ip

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Size of various fields in an IP header.
const (
	versionBits = 4 // Number of bits in the IP version.
	versionPos  = 4 // Offset in bits of the IP version.

	versionByteOffset = 0 // Offset in bytes of the IP version.
	versionByteSize   = 1 // Number of bytes encoding the IP version
)

// Constants defined for various protocol types.
const (
	protoICMP4    = 1   // ICMP4 packet.
	protoICMP6    = 58  // ICMP6 packet.
	protoTCP      = 6   // TCP packet.
	protoUDP      = 17  // UDP packet.
	protoIP4IP4   = 4   // IPv4 over IPv4 tunnel.
	protoIP6IP4   = 41  // IPv6 over IPv4 tunnel.
	protoGRE      = 47  // GRE tunnel.
	protoReserved = 255 // Reserved payload.
)

// protoHeader maps protocols to packet headers.
var protoHeader = map[uint8]fwdpb.PacketHeaderId{}

// headerProto maps packet headers to protocols.
var headerProto = map[fwdpb.PacketHeaderId]uint8{
	fwdpb.PacketHeaderId_IP4:    protoIP4IP4,
	fwdpb.PacketHeaderId_IP6:    protoIP6IP4,
	fwdpb.PacketHeaderId_GRE:    protoGRE,
	fwdpb.PacketHeaderId_TCP:    protoTCP,
	fwdpb.PacketHeaderId_UDP:    protoUDP,
	fwdpb.PacketHeaderId_ICMP4:  protoICMP4,
	fwdpb.PacketHeaderId_ICMP6:  protoICMP6,
	fwdpb.PacketHeaderId_OPAQUE: protoReserved,
}

// ipVersion extracts the IP version from a byte.
func ipVersion(b []byte) frame.Field {
	if version := frame.Header(b).Field(versionByteOffset, versionByteSize); version != nil {
		return version.BitField(versionPos, versionBits)
	}
	return nil
}

// header is a set of functions to manipulate an individual IP header.
type header interface {
	// Header returns the protocol header.
	Header() []byte

	// ID returns the protocol header ID.
	ID() fwdpb.PacketHeaderId

	// Find returns a copy of the specified field.
	Find(id fwdpacket.FieldID) ([]byte, error)

	// Update updates the specified field and returns true if the
	// header is dirty.
	Update(id fwdpacket.FieldID, op int, arg []byte) (bool, error)

	// Payload gets the payload information.
	Payload() (fwdpb.PacketHeaderId, int64)

	// SetPayload sets the payload.
	SetPayload(id fwdpb.PacketHeaderId, length int64)
}

// An IP is a group of IP headers in the L3 portion of the packet.
// The L3/IP portion of a packet can contain a series IPv4, IPv6 and GRE
// headers in-case of nested tunnels.
//
// Note that the GRE implementation does not support UDF.
type IP struct {
	headers []header       // Sequence of IP headers.
	desc    *protocol.Desc // Descriptor for L3/IP headers.
}

// Header returns the IP/L3 as a slice of bytes.
func (ip *IP) Header() []byte {
	var b = make([][]byte, len(ip.headers))
	for p, h := range ip.headers {
		b[p] = h.Header()
	}
	return bytes.Join(b, []byte{})
}

// Trailer returns the no trailing bytes.
func (IP) Trailer() []byte {
	return nil
}

// ID returns the protocol header ID of the outermost ip header.
func (ip *IP) ID(instance int) fwdpb.PacketHeaderId {
	index := instance
	if instance == fwdpacket.LastField {
		index = len(ip.headers) - 1
	}
	if index >= len(ip.headers) {
		return fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE
	}
	return ip.headers[index].ID()
}

// lookup finds the IP header that contains the field.
//
// A couple of implementation notes:
// The GRE header does not support UDF.
// IPv4 and IPv6 contain the same fields (for most part).
func (ip *IP) lookup(id fwdpacket.FieldID) (header, error) {
	greField := false // Indicates if we need to find a GRE header.
	if id.Num == fwdpb.PacketFieldNum_GRE_KEY || id.Num == fwdpb.PacketFieldNum_GRE_SEQUENCE {
		greField = true
	}

	var found header
	instance := uint8(0)
	for _, header := range ip.headers {
		if hid := header.ID(); greField == (hid == fwdpb.PacketHeaderId_GRE) {
			found = header
			if instance == id.Instance {
				break
			}
			instance++
		}
	}
	if found != nil {
		return found, nil
	}
	return nil, fmt.Errorf("ip: lookup failed, header for field %v does not exist", id)
}

// Field returns a copy of the specified field.
func (ip *IP) Field(id fwdpacket.FieldID) ([]byte, error) {
	header, err := ip.lookup(id)
	if err != nil {
		return nil, fmt.Errorf("ip: Field failed, err %v", err)
	}
	return header.Find(id)
}

// UpdateField updates the specified field.
func (ip *IP) UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	header, err := ip.lookup(id)
	if err != nil {
		return false, fmt.Errorf("ip: UpdateField failed, err %v", err)
	}
	return header.Update(id, op, arg)
}

// Remove removes the outermost IP header from the packet after verifying its type.
func (ip *IP) Remove(id fwdpb.PacketHeaderId) error {
	switch id {
	case fwdpb.PacketHeaderId_TUNNEL_6TO4_AUTO, fwdpb.PacketHeaderId_TUNNEL_6TO4_SECURE:
		// There must be at-least two IP headers for an IP tunnel.
		if len(ip.headers) < 2 {
			return fmt.Errorf("ip: Remove header %v failed, insufficient header length %v", id, len(ip.headers))
		}
		if err := check6to4Tunnel(id, ip.headers[1], ip.headers[0]); err != nil {
			return fmt.Errorf("ip: Remove header %v failed, err %v", id, err)
		}

	default:
		if hid := ip.headers[0].ID(); hid != id {
			return fmt.Errorf("ip: Remove header %v failed, outermost header is %v", id, hid)
		}
	}
	ip.headers = ip.headers[1:]
	return nil
}

// Modify adds an addition IP header effectively tunneling the packet.
func (ip *IP) Modify(id fwdpb.PacketHeaderId) error {
	var h header
	switch id {
	case fwdpb.PacketHeaderId_GRE:
		if hid := ip.headers[0].ID(); hid == fwdpb.PacketHeaderId_GRE {
			return errors.New("ip: Modify header failed, GRE cannot encapsulate a GRE header")
		}
		h = newGRE()

	case fwdpb.PacketHeaderId_IP4:
		h = newIP4()

	case fwdpb.PacketHeaderId_IP6:
		h = newIP6()

	case fwdpb.PacketHeaderId_TUNNEL_6TO4_AUTO, fwdpb.PacketHeaderId_TUNNEL_6TO4_SECURE:
		if len(ip.headers) == 0 {
			return fmt.Errorf("ip: Modify header failed for header %v, no ip headers", id)
		}
		var err error
		if h, err = new6to4Tunnel(id, ip.headers[0]); err != nil {
			return err
		}

	default:
		return fmt.Errorf("ip: Modify header failed for header %v", id)
	}
	curr := ip.headers
	ip.headers = []header{h}
	ip.headers = append(ip.headers, curr...)
	return nil
}

// Rebuild updates all the IP headers starting from the innermost IP header.
func (ip *IP) Rebuild() error {
	length := int64(ip.desc.PayloadLength())
	id := ip.desc.PayloadID()
	for index := len(ip.headers) - 1; index >= 0; index-- {
		header := ip.headers[index]
		header.SetPayload(id, length)
		id, length = header.Payload()
	}
	return nil
}

// add adds an IP header to an empty L3.
// Note that the L3 cannot contain a stand-alone GRE header.
func add(id fwdpb.PacketHeaderId, desc *protocol.Desc) (protocol.Handler, error) {
	var h header
	switch id {
	case fwdpb.PacketHeaderId_IP4:
		h = newIP4()

	case fwdpb.PacketHeaderId_IP6:
		h = newIP6()

	default:
		return nil, fmt.Errorf("ip: add failed, unable to add header %v", id)
	}
	return &IP{
		headers: []header{h},
	}, nil
}

// parseTunnel processes a frame to parse all IP tunnel headers.
func (ip *IP) parseTunnel(frame *frame.Frame, next fwdpb.PacketHeaderId) (fwdpb.PacketHeaderId, error) {
	var err error
	var header header

	for frame.Len() != 0 {
		switch next {
		case fwdpb.PacketHeaderId_GRE:
			header, next, err = makeGRE(frame)

		case fwdpb.PacketHeaderId_IP4:
			header, next, err = makeIP4(frame)

		case fwdpb.PacketHeaderId_IP6:
			header, next, err = makeIP6(frame)

		default:
			return next, nil
		}
		if err != nil {
			return fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, err
		}
		ip.headers = append(ip.headers, header)
	}
	return fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, nil
}

// parseIP6 parses an L3 starting with an IP6 header.
func parseIP6(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	h, next, err := makeIP6(frame)
	if err != nil {
		return nil, next, err
	}
	ip := &IP{
		desc:    desc,
		headers: []header{h},
	}
	if next, err = ip.parseTunnel(frame, next); err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, err
	}
	return ip, next, err
}

// parseIP4 parses an L3 starting with an IP4 header.
func parseIP4(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	h, next, err := makeIP4(frame)
	if err != nil {
		return nil, next, err
	}
	ip := &IP{
		desc:    desc,
		headers: []header{h},
	}
	if next, err = ip.parseTunnel(frame, next); err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, err
	}
	return ip, next, nil
}

// parseIP parses an L3 as an IP header. It peeks into the IP version and parses
// the frames as IPv4 or IPv6.
func parseIP(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	vb, err := frame.Peek(versionByteOffset, versionByteSize)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("ip: parseIP failed, err %v", err)
	}
	switch version := ipVersion(vb); version.Value() {
	case 4: // IP4
		return parseIP4(frame, desc)

	case 6: // IP6
		return parseIP6(frame, desc)

	default:
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("ip: parseIP failed, unknown version %x", version)
	}
}

func init() {
	// Compute reverse mapping tables.
	for id, proto := range headerProto {
		protoHeader[proto] = id
	}

	// Register the parse and add functions for the IP4 and IP6 header and its variants.
	protocol.Register(fwdpb.PacketHeaderId_IP4, parseIP4, add)
	protocol.Register(fwdpb.PacketHeaderId_IP6, parseIP6, add)

	// Register the parse function for the generic IP header. The generic IP
	// header can be used to parse the frame. However one cannot add a generic
	// IP header.
	protocol.Register(fwdpb.PacketHeaderId_IP, parseIP, nil)
}
