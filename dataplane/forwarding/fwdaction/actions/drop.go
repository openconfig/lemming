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

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A drop is an action that drops all packets.
type drop struct{}

// String formats the state of the action as a string.
func (drop) String() string {
	return fmt.Sprintf("Type=%v;", fwdpb.ActionType_DROP_ACTION)
}

// Process processes the packet by dropping it.
func (*drop) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if counters != nil {
		counters.Increment(fwdpb.CounterId_DROP_PACKETS, 1)
		counters.Increment(fwdpb.CounterId_DROP_OCTETS, uint32(packet.Length()))
	}
	return nil, fwdaction.DROP
}

// A dropBuilder builds drop actions.
type dropBuilder struct{}

// init registers a builder for the drop action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_DROP_ACTION, &dropBuilder{})
}

// Build creates a new drop action.
func (*dropBuilder) Build(desc *fwdpb.ActionDesc, _ *fwdcontext.Context) (fwdaction.Action, error) {
	return &drop{}, nil
}
