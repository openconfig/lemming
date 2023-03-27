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
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// reparse is an action that reparses the current packet with the specified
// start header.
type reparse struct {
	id      fwdpb.PacketHeaderId
	fields  []fwdpacket.FieldID
	prepend []byte
}

// String formats the state of the action as a string.
func (r *reparse) String() string {
	return fmt.Sprintf("Type=%v;HeaderId=%v;Fields=%+v;Prepend=%x", fwdpb.ActionType_ACTION_TYPE_REPARSE, r.id, r.fields, r.prepend)
}

// Process reparses the packet.
func (r *reparse) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if err := packet.Reparse(r.id, r.fields, r.prepend); err != nil {
		counters.Increment(fwdpb.CounterId_COUNTER_ID_ERROR_PACKETS, 1)
		counters.Increment(fwdpb.CounterId_COUNTER_ID_ERROR_OCTETS, uint32(packet.Length()))
		log.Errorf("actions: Failed to reparse packet, err %v", err)
		return nil, fwdaction.DROP
	}
	return nil, fwdaction.CONTINUE
}

// reparseBuilder builds reparse actions.
type reparseBuilder struct{}

func init() {
	// Register a builder for the reparse action type.
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_REPARSE, &reparseBuilder{})
}

// Build creates a new reparse action.
func (*reparseBuilder) Build(desc *fwdpb.ActionDesc, _ *fwdcontext.Context) (fwdaction.Action, error) {
	r, ok := desc.Action.(*fwdpb.ActionDesc_Reparse)
	if !ok {
		return nil, fmt.Errorf("actions: Build for reparse action failed, missing desc")
	}
	var fields []fwdpacket.FieldID
	for _, f := range r.Reparse.GetFieldIds() {
		fields = append(fields, fwdpacket.NewFieldID(f))
	}
	return &reparse{id: r.Reparse.GetHeaderId(), fields: fields, prepend: r.Reparse.GetPrepend()}, nil
}
