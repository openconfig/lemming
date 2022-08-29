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
	"errors"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	typeBytes  = 1 // Number of bytes in the ICMP type
	typeOffset = 0 // Offset in bytes of the ICMP type
	codeBytes  = 1 // Number of bytes in the ICMP code
	codeOffset = 1 // Offset in bytes of the ICMP code
	csumBytes  = 2 // Number of bytes in the ICMP checksum
	csumOffset = 2 // Offset in bytes of the ICMP checksum
	icmpBytes  = 8 // Number of bytes in an ICMP header
)

// An ICMP represents the ICMP header in the packet. Its state and functions
// handle parts of ICMP that are common to IPv4 and IPv6.
type ICMP struct {
	header frame.Header // ICMP header.
}

// Header returns the ICMP header.
func (i *ICMP) Header() []byte {
	return i.header
}

// Trailer returns the no trailing bytes.
func (ICMP) Trailer() []byte {
	return nil
}

// commonField returns bytes within the ICMP header as identified by id.
func (i *ICMP) commonField(id fwdpacket.FieldID) frame.Field {
	switch {
	case id.IsUDF:
		return protocol.UDF(i.header, id)

	case id.Num == fwdpb.PacketFieldNum_ICMP_TYPE:
		return i.header.Field(typeOffset, typeBytes)

	case id.Num == fwdpb.PacketFieldNum_ICMP_CODE:
		return i.header.Field(codeOffset, codeBytes)

	default:
		return nil
	}
}

// Remove returns an error as the ICMP header cannot be removed.
func (ICMP) Remove(fwdpb.PacketHeaderId) error {
	return errors.New("icmp: Remove is unsupported")
}

// Modify returns an error as the ICMP header does not support adding extensions.
func (ICMP) Modify(fwdpb.PacketHeaderId) error {
	return errors.New("icmp: Modify failed, unsupported operation")
}
