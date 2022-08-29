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

// Package opaque implements the opaque packet header.
package opaque

import (
	"errors"
	"fmt"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// An Opaque represents payload that cannot be parsed by Lucius.
// Note opaque payload cannot be added or removed.
type Opaque struct {
	header frame.Header
}

// Header returns the opaque frame.
func (op *Opaque) Header() []byte {
	return op.header
}

// Trailer returns the no trailing bytes.
func (Opaque) Trailer() []byte {
	return nil
}

// ID returns the protocol header ID.
func (Opaque) ID(int) fwdpb.PacketHeaderId {
	return fwdpb.PacketHeaderId_OPAQUE
}

// Field returns a slice of bytes as identified by id.
func (op *Opaque) Field(id fwdpacket.FieldID) ([]byte, error) {
	if id.IsUDF {
		if field := protocol.UDF(op.header, id); field != nil {
			return field.Copy(), nil
		}
	}
	return nil, fmt.Errorf("opaque: Field failed, field %v does not exist", id)
}

// UpdateField updates a slice of bytes identified by id.
func (op *Opaque) UpdateField(id fwdpacket.FieldID, oper int, arg []byte) (bool, error) {
	if id.IsUDF && oper == fwdpacket.OpSet {
		if field := protocol.UDF(op.header, id); field != nil {
			return true, field.Set(arg)
		}
	}
	return false, fmt.Errorf("opaque: UpdateField failed, field %v does not exist", id)
}

// Rebuild does not perform any updates.
func (Opaque) Rebuild() error {
	return nil
}

// Remove returns an error as the opaque header cannot be removed.
func (Opaque) Remove(fwdpb.PacketHeaderId) error {
	return errors.New("opaque: Remove is unsupported")
}

// Modify returns an error as the opaque header's type cannot be modified.
func (Opaque) Modify(fwdpb.PacketHeaderId) error {
	return errors.New("opaque: Modify is unsupported")
}

// parse reads an Opaque header from the frame.
func parse(frame *frame.Frame, _ *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	header, err := frame.ReadHeader(frame.Len())
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, fmt.Errorf("opaque: parse failed, err %v", err)
	}
	return &Opaque{
		header: header,
	}, fwdpb.PacketHeaderId_PACKET_HEADER_ID_NONE, nil
}

func init() {
	// Register the parse function for the OPAQUE headers.
	// Note that opaque cannot be added explicitly.
	protocol.Register(fwdpb.PacketHeaderId_OPAQUE, parse, nil)
}
