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
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A lookup is an action that looks up the packet in a table.
type lookup struct {
	table fwdtable.Table
}

// String formats the state of the action as a string.
func (l *lookup) String() string {
	if l.table == nil {
		return fmt.Sprintf("Type=%v;<Table=nil>", fwdpb.ActionType_ACTION_TYPE_LOOKUP)
	}
	return fmt.Sprintf("Type=%v;<Table=%v>", fwdpb.ActionType_ACTION_TYPE_LOOKUP, l.table.ID())
}

// Cleanup releases the table.
func (l *lookup) Cleanup() {
	if err := fwdtable.Release(l.table); err != nil {
		log.Errorf("actions: Cleanup failed for action %v, err %s.", l, err)
	}
	l.table = nil
}

// Process processes the packet by looking up a table.
// A lookup action returns the result of the lookup if the table exists.
// If the lookup is not associated with a table, the packet is dropped.
func (l *lookup) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if l.table == nil {
		counters.Increment(fwdpb.CounterId_COUNTER_ID_ERROR_PACKETS, 1)
		counters.Increment(fwdpb.CounterId_COUNTER_ID_ERROR_OCTETS, uint32(packet.Length()))
		return nil, fwdaction.DROP
	}
	return l.table.Process(packet, counters)
}

// A lookupBuilder builds lookup actions.
type lookupBuilder struct{}

// init registers a builder for the lookup action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_LOOKUP, &lookupBuilder{})
}

// Build creates a new lookup action.
func (*lookupBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	look, ok := desc.Action.(*fwdpb.ActionDesc_Lookup)
	if !ok {
		return nil, fmt.Errorf("actions: Build for lookup action failed, missing desc")
	}
	table, err := fwdtable.Acquire(ctx, look.Lookup.GetTableId())
	if err != nil {
		return nil, fmt.Errorf("actions: Build for lookup action failed, err %v", err)
	}
	return &lookup{table: table}, nil
}
