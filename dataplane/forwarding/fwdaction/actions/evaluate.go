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

// An evaluate is an action that causes all actions marked as "on-evaluate" to be applied to the packet
// and transmit it if the packet has an output port.
type evaluate struct {
	ctx *fwdcontext.Context
}

// String formats the state of the action as a string.
func (evaluate) String() string {
	return fmt.Sprintf("Type=%v;", fwdpb.ActionType_ACTION_TYPE_EVALUATE)
}

// Process evaluates if the packet has an evaluate port and transmits it.
func (o *evaluate) Process(fwdpacket.Packet, fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	return nil, fwdaction.EVALUATE
}

// An evaluateBuilder builds evaluate actions.
type evaluateBuilder struct{}

// init registers a builder for the evaluate action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_EVALUATE, &evaluateBuilder{})
}

// Build creates a new evaluate action.
func (*evaluateBuilder) Build(_ *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	return &evaluate{
		ctx: ctx,
	}, nil
}
