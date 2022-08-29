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

	log "github.com/golang/glog"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdflowcounter"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// flowcounter is an action that increments the per flow counters based on the packet.
type flowcounter struct {
	counter *fwdflowcounter.FlowCounter
}

// String formats the type of the action as a string.
func (f *flowcounter) String() string {
	if f.counter == nil {
		return fmt.Sprintf("Type=%v;<FlowCounter=nil>", fwdpb.ActionType_FLOW_COUNTER_ACTION)
	}
	return fmt.Sprintf("Type=%v;<FlowCounter=%v>", fwdpb.ActionType_FLOW_COUNTER_ACTION, f.counter.ID())
}

// Cleanup releases the flowCounter.
func (f *flowcounter) Cleanup() {
	if err := fwdflowcounter.Release(f.counter); err != nil {
		log.Errorf("actions: Cleanup failed for action flowcounter, err %s", err)
	}
	f.counter = nil
}

// Process increments the octet and packet counter fields, based on the packet.
func (f *flowcounter) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if f.counter != nil {
		octetCount := uint32(packet.Length())
		const packetCount = 1
		f.counter.Process(octetCount, packetCount)
	}
	return nil, fwdaction.CONTINUE
}

// flowcounterBuilder builds flowcounter actions.
type flowcounterBuilder struct{}

// init registers a builder for the flowcounter action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_FLOW_COUNTER_ACTION, &flowcounterBuilder{})
}

// Build creates a new flowcounter action.
func (*flowcounterBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	if !proto.HasExtension(desc, fwdpb.E_FlowCounterActionDesc_Extension) {
		return nil, fmt.Errorf("actions: Build for flowcounter action failed, missing extension %s", fwdpb.E_FlowCounterActionDesc_Extension.Name)
	}
	fcExt := proto.GetExtension(desc, fwdpb.E_FlowCounterActionDesc_Extension).(*fwdpb.FlowCounterActionDesc)
	fctr, err := fwdflowcounter.Acquire(ctx, fcExt.GetCounterId())
	if err != nil {
		return nil, fmt.Errorf("actions: Build for flowcounter action failed, err %v", err)
	}
	return &flowcounter{counter: fctr}, nil
}
