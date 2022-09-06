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

package bridge

import (
	"fmt"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A learn is an action that learns the packet in a bridge.
type learn struct {
	table *Table
}

// String formats the state of the action as a string.
func (l *learn) String() string {
	if l.table == nil {
		return fmt.Sprintf("Type=%v;<Table=nil>", fwdpb.ActionType_ACTION_TYPE_BRIDGE_LEARN)
	}
	return fmt.Sprintf("Type=%v;<Table=%v>", fwdpb.ActionType_ACTION_TYPE_BRIDGE_LEARN, l.table.ID())
}

// Cleanup releases the table.
func (l *learn) Cleanup() {
	if err := fwdtable.Release(l.table); err != nil {
		log.Errorf("actions: Cleanup failed for action %v, err %s.", l, err)
	}
	l.table = nil
}

// Process processes the packet by learning it. It does not drop the packet on
// error. The errors are logged and the packet processing continues.
func (l *learn) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if l.table == nil {
		counters.Increment(fwdpb.CounterId_COUNTER_ID_ERROR_PACKETS, 1)
		counters.Increment(fwdpb.CounterId_COUNTER_ID_ERROR_OCTETS, uint32(packet.Length()))
		return nil, fwdaction.DROP
	}
	if err := l.table.Learn(packet); err != nil {
		log.Warningf("bridge: Error during learn, err %v, action %v.", err, l)
	}
	return nil, fwdaction.CONTINUE
}

// A learnBuilder builds learn actions.
type learnBuilder struct{}

// init registers a builder for the learn action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_BRIDGE_LEARN, &learnBuilder{})
}

// Build creates a new learn action.
func (*learnBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	bl, ok := desc.Action.(*fwdpb.ActionDesc_Bridge)
	if !ok {
		return nil, fmt.Errorf("actions: Build for learn action failed, missing extension")
	}
	table, err := fwdtable.Acquire(ctx, bl.Bridge.GetTableId())
	if err != nil {
		return nil, fmt.Errorf("actions: Build for learn action failed, err %v", err)
	}
	b, ok := table.(*Table)
	if !ok {
		fwdtable.Release(table)
		return nil, fmt.Errorf("actions: Build for learn action failed, table %v is not a bridge", bl.Bridge.GetTableId())
	}
	return &learn{table: b}, nil
}
