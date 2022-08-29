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

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A transmit is an action that writes a packet to a port.
type transmit struct {
	port      fwdport.Port
	immediate bool
}

// String formats the state of the action as a string.
func (t *transmit) String() string {
	if t.port == nil {
		return fmt.Sprintf("Type=%s;Immediate=%v;<Port=nil>;", fwdpb.ActionType_ACTION_TYPE_TRANSMIT, t.immediate)
	}
	return fmt.Sprintf("Type=%s;Immediate=%v;<Port=%v>;", fwdpb.ActionType_ACTION_TYPE_TRANSMIT, t.immediate, t.port.ID())
}

// Cleanup releases the port.
func (t *transmit) Cleanup() {
	if err := fwdport.Release(t.port); err != nil {
		log.Errorf("Cleanup failed for action %v, err %s.", t, err)
	}
	t.port = nil
}

// Process processes the packet by setting up its output port.
// If transmit does not have a port, the packet is dropped.
func (t *transmit) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if t.port == nil {
		counters.Increment(fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS, 1)
		counters.Increment(fwdpb.CounterId_COUNTER_ID_TX_ERROR_OCTETS, uint32(packet.Length()))
		return nil, fwdaction.DROP
	}
	fwdport.SetOutputPort(packet, t.port)
	if t.immediate {
		return nil, fwdaction.OUTPUT
	}
	return nil, fwdaction.CONTINUE
}

// A transmitBuilder builds transmit actions.
type transmitBuilder struct{}

// init registers a builder for the transmit action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_TRANSMIT, &transmitBuilder{})
}

// Build creates a new transmit action.
func (*transmitBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	tr, ok := desc.Action.(*fwdpb.ActionDesc_Transmit)
	if !ok {
		return nil, fmt.Errorf("actions: Build for lookup action failed, missing extension")
	}
	port, err := fwdport.Acquire(tr.Transmit.GetPortId(), ctx)
	if err != nil {
		return nil, fmt.Errorf("actions: Build for transmit action failed, err %v", err)
	}
	return &transmit{port: port, immediate: tr.Transmit.GetImmediate()}, nil
}
