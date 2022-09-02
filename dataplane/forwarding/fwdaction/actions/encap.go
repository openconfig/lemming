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

package actions

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// encap is an action that adds a header to a packet.
type encap struct {
	id fwdpb.PacketHeaderId
}

// String formats the state of the action as a string.
func (e *encap) String() string {
	return fmt.Sprintf("Type=%v;HeaderId=%v;", fwdpb.ActionType_ENCAP_ACTION, e.id)
}

// Process adds a header to the packet.
func (e *encap) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if err := packet.Encap(e.id); err != nil {
		counters.Increment(fwdpb.CounterId_ENCAP_ERROR_PACKETS, 1)
		counters.Increment(fwdpb.CounterId_ENCAP_ERROR_OCTETS, uint32(packet.Length()))
		return nil, fwdaction.DROP
	}
	return nil, fwdaction.CONTINUE
}

// encapBuilder builds encap actions.
type encapBuilder struct{}

func init() {
	// Register a builder for the encap action type.
	fwdaction.Register(fwdpb.ActionType_ENCAP_ACTION, &encapBuilder{})
}

// Build creates a new encap action.
func (*encapBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	if !proto.HasExtension(desc, fwdpb.E_EncapActionDesc_Extension) {
		return nil, fmt.Errorf("actions: Build for encap action failed, missing extension %s", fwdpb.E_EncapActionDesc_Extension.Name)
	}
	encapExt := proto.GetExtension(desc, fwdpb.E_EncapActionDesc_Extension).(*fwdpb.EncapActionDesc)
	return &encap{id: encapExt.GetHeaderId()}, nil
}
