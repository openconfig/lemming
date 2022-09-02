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

// Package icmp implements the ICMP header support in Lucius.
package icmp

import (
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	"github.com/openconfig/lemming/dataplane/forwarding/util/hash/csum16"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// An ICMP4 represents an ICMP4 header in the packet. It can update the ICMP
// header but it cannot add or remove an ICMP header. The ICMP message is
// assumed to contain the ICMP header and all bytes following it.
type ICMP4 struct {
	ICMP
	desc *protocol.Desc // Descriptor for L3/IP headers.
}

// ID returns the ICMP protocol header ID.
func (ICMP4) ID(int) fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_ICMP4
}

// Field returns a copy of the bytes for the specfied field in the ICMP header.
func (i *ICMP4) Field(id fwdpacket.FieldID) ([]byte, error) {
	if field := i.commonField(id); field != nil {
		return field.Copy(), nil
	}
	return nil, fmt.Errorf("icmp4: Field failed, field %v does not exist", id)
}

// UpdateField sets bytes within the ICMP header.
func (i *ICMP4) UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	if field := i.commonField(id); field != nil && op == fwdpacket.OpSet {
		return true, field.Set(arg)
	}
	return false, fmt.Errorf("icmp4: UpdateField failed, unsupported op %v for field %v", op, id)
}

// Rebuild updates the ICMP checksum if it is dirty.
func (i *ICMP4) Rebuild() error {
	if !i.desc.Dirty() {
		return nil
	}
	i.header.Field(csumOffset, csumBytes).SetValue(0)
	var sum csum16.Sum
	sum.Write(i.header)
	i.header.Field(csumOffset, csumBytes).SetValue(uint(sum))
	return nil
}

// parseICMP4 parses an ICMP4 header in the packet. All data following the ICMP
// header is assumed to be a part of the ICMP message.
func parseICMP4(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	if frame.Len() < icmpBytes {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("icmp: parse failed, invalid frame length %v, expected size is greater than %v bytes", frame.Len(), icmpBytes)
	}
	header, err := frame.ReadHeader(frame.Len())
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("icmp: parse failed, err %v", err)
	}
	i := &ICMP4{
		desc: desc,
	}
	i.ICMP.header = header
	return i, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, nil
}

func init() {
	// ICMP header cannot be added to a packet.
	protocol.Register(fwdpb.PacketHeaderId_ICMP4, parseICMP4, nil)
}
