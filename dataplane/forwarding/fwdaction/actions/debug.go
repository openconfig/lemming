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

// A debug is an action that sets the debug flag on a packet.
type debug struct{}

// String formats the state of the action as a string.
func (debug) String() string {
	return fmt.Sprintf("Type=%v;", fwdpb.ActionType_ACTION_TYPE_DEBUG)
}

// Process processes the packet by setting it's debug flag.
func (*debug) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if counters != nil {
		counters.Increment(fwdpb.CounterId_COUNTER_ID_RX_DEBUG_PACKETS, 1)
		counters.Increment(fwdpb.CounterId_COUNTER_ID_RX_DEBUG_OCTETS, uint32(packet.Length()))
	}
	packet.Debug(true)
	return nil, fwdaction.CONTINUE
}

// A debugBuilder builds debug actions.
type debugBuilder struct{}

// init registers a builder for the debug action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_DEBUG, &debugBuilder{})
}

// Build creates a new debug action.
func (*debugBuilder) Build(*fwdpb.ActionDesc, *fwdcontext.Context) (fwdaction.Action, error) {
	return &debug{}, nil
}
