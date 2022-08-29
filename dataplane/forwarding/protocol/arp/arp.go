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

// Package arp implements the ARP header.
package arp

import (
	"errors"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	tpaBytes   = 4                // Number of bytes in the target protocol address.
	tpaOffset  = 24               // Offset in bytes of the target protocol address.
	spaBytes   = 4                // Number of bytes in the sender protocol address.
	spaOffset  = 14               // Offset in bytes of the sender protocol address.
	tmacBytes  = protocol.SizeMAC // Number of bytes in the target mac address.
	tmacOffset = 18               // Offset in bytes of the target mac address.
	smacBytes  = protocol.SizeMAC // Number of bytes in the source mac address.
	smacOffset = 8                // Offset in bytes of the source mac address.
	arpBytes   = 28               // Number of bytes in an arp header.
)

// An ARP is an ARP header. It can parse, query and operate on an ARP
// header in the packet. However it cannot add a new ARP header to a packet.
type ARP struct {
	header frame.Header   // ARP header.
	desc   *protocol.Desc // Protocol descriptor.
}

// Header returns the ARP header.
func (a *ARP) Header() []byte {
	return a.header
}

// Trailer returns the no trailing bytes.
func (ARP) Trailer() []byte {
	return nil
}

// ID returns the ARP protocol header ID.
func (ARP) ID(int) fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_ARP
}

// field returns the bytes within the ARP header as specified by id.
func (a *ARP) field(id fwdpacket.FieldID) frame.Field {
	switch {
	case id.IsUDF:
		return protocol.UDF(a.header, id)

	case id.Num == fwdpb.PacketFieldNum_ARP_TPA:
		return a.header.Field(tpaOffset, tpaBytes)

	case id.Num == fwdpb.PacketFieldNum_ARP_SPA:
		return a.header.Field(spaOffset, spaBytes)

	case id.Num == fwdpb.PacketFieldNum_ARP_TMAC:
		return a.header.Field(tmacOffset, tmacBytes)

	case id.Num == fwdpb.PacketFieldNum_ARP_SMAC:
		return a.header.Field(smacOffset, smacBytes)

	default:
		return nil
	}
}

// Field returns bytes within the ARP header as specified by id.
func (a *ARP) Field(id fwdpacket.FieldID) ([]byte, error) {
	if field := a.field(id); field != nil {
		return field.Copy(), nil
	}
	return nil, fmt.Errorf("arp: Field failed, field %v does not exist", id)
}

// UpdateField sets bytes within the ARP header.
func (a *ARP) UpdateField(id fwdpacket.FieldID, op int, arg []byte) (bool, error) {
	if field := a.field(id); field != nil && op == fwdpacket.OpSet {
		return true, field.Set(arg)
	}
	return false, fmt.Errorf("arp: UpdateField failed, unsupported op %v for field %v", op, id)
}

// Rebuild succeeds by default as the ARP header does not need updates.
func (ARP) Rebuild() error {
	return nil
}

// Remove returns an error as the ARP header cannot be removed.
func (ARP) Remove(fwdpb.PacketHeaderId) error {
	return errors.New("arp: Remove is unsupported")
}

// Modify returns an error as the ARP header does not support extensions.
func (ARP) Modify(fwdpb.PacketHeaderId) error {
	return errors.New("arp: Modify is unsupported")
}

// parse parses an ARP header.
func parse(frame *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	if frame.Len() != arpBytes {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("arp: parse failed, invalid frame length %v, expected %v bytes", frame.Len(), arpBytes)
	}
	header, err := frame.ReadHeader(arpBytes)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("arp: parse failed, err %v", err)
	}
	return &ARP{
		header: header,
		desc:   desc,
	}, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, nil
}

func init() {
	// ARP header cannot be added to a packet.
	protocol.Register(fwdpb.PacketHeaderId_ARP, parse, nil)
}
