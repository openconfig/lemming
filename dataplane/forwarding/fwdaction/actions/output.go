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
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// An output is an action that causes the packet to be transmitted if the packet
// has an output port.
type output struct {
	ctx *fwdcontext.Context
}

// String formats the state of the action as a string.
func (output) String() string {
	return fmt.Sprintf("Type=%v;", fwdpb.ActionType_ACTION_TYPE_OUTPUT)
}

// Process evaluates if the packet has an output port and transmits it.
func (o *output) Process(pkt fwdpacket.Packet, _ fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if p, err := fwdport.OutputPort(pkt, o.ctx); err == nil && p != nil {
		return nil, fwdaction.OUTPUT
	}
	return nil, fwdaction.CONTINUE
}

// An outputBuilder builds output actions.
type outputBuilder struct{}

// init registers a builder for the output action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_OUTPUT, &outputBuilder{})
}

// Build creates a new output action.
func (*outputBuilder) Build(_ *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	return &output{
		ctx: ctx,
	}, nil
}
