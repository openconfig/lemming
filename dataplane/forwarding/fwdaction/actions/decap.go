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

// decap is an action that removes a header from a packet.
type decap struct {
	id fwdpb.PacketHeaderId
}

// String formats the state of the action as a string.
func (d *decap) String() string {
	return fmt.Sprintf("Type=%s;HeaderId=%v;", fwdpb.ActionType_DECAP_ACTION, d.id)
}

// Process removes a header from the packet. Decap errors are explicitly ignored.
func (d *decap) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	_ = packet.Decap(d.id)
	return nil, fwdaction.CONTINUE
}

// decapBuilder builds decap actions.
type decapBuilder struct{}

func init() {
	// Register a builder for the decap action type.
	fwdaction.Register(fwdpb.ActionType_DECAP_ACTION, &decapBuilder{})
}

// Build creates a new decap action.
func (*decapBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	if !proto.HasExtension(desc, fwdpb.E_DecapActionDesc_Extension) {
		return nil, fmt.Errorf("actions: Build for decap action failed, missing extension %s", fwdpb.E_DecapActionDesc_Extension.Name)
	}
	decapExt := proto.GetExtension(desc, fwdpb.E_DecapActionDesc_Extension).(*fwdpb.DecapActionDesc)
	return &decap{id: decapExt.GetHeaderId()}, nil
}
