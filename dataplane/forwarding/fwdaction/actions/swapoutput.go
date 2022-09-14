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
	"github.com/openconfig/lemming/dataplane/internal/engine"

	log "github.com/golang/glog"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// swapOutput is a action used to set a packet's output port to the input port related port.
// If the a packet's input port is an external port, than its the output the corresponding TAP port or vice versa.
type swapOutput struct {
	ctx *fwdcontext.Context
}

// String formats the state of the action as a string.
func (swapOutput) String() string {
	return fmt.Sprintf("Type=%v;", fwdpb.ActionType_ACTION_TYPE_SWAP_OUTPUT)
}

// Process processes the packet by updating the output port.
func (s *swapOutput) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	p, err := fwdport.InputPort(packet, s.ctx)
	if err != nil {
		log.Warningf("failed to get input port: %v", err)
		return nil, fwdaction.DROP
	}
	var outPort string

	if engine.IsTap(string(p.ID())) {
		outPort = engine.TapNameToIntfName(string(p.ID()))
	} else {
		outPort = engine.IntfNameToTapName(string(p.ID()))
	}
	port, err := fwdport.Acquire(&fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: outPort}}, s.ctx)
	if err != nil {
		log.Warningf("failed to get output port: %v", err)
		return nil, fwdaction.DROP
	}
	fwdport.SetOutputPort(packet, port)
	return nil, fwdaction.CONTINUE
}

// A swapOutputBuilder builds swap output actions.
type swapOutputBuilder struct{}

// init registers a builder for the swap output action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_SWAP_OUTPUT, &swapOutputBuilder{})
}

// Build creates a new swap output action.
func (*swapOutputBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	return &swapOutput{
		ctx: ctx,
	}, nil
}
