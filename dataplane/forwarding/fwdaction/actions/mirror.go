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
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A mirror is an action that injects a replica of the packet to a port.
type mirror struct {
	actions fwdaction.Actions   // set of actions to apply to the mirrored packet
	port    fwdport.Port        // Port to which the mirrored packet is injected
	dir     fwdpb.PortAction    // Selects how the mirrored packet is injected
	ctx     *fwdcontext.Context // Context containing port
	fields  []fwdpacket.FieldID // Fields that are copied to the mirrored packet
}

// String formats the state of the action as a string.
func (m *mirror) String() string {
	pid := "nil"
	if m.port != nil {
		pid = string(m.port.ID())
	}
	return fmt.Sprintf("Type=%s;Dir=%v;<Port=%v>;<Actions=%v>;<Fields=%v>;", fwdpb.ActionType_MIRROR_ACTION, m.dir, pid, m.actions, m.fields)
}

// Cleanup releases the port.
func (m *mirror) Cleanup() {
	if m.port != nil {
		if err := fwdport.Release(m.port); err != nil {
			log.Errorf("Cleanup failed for action %v, err %s.", m, err)
		}
		m.port = nil
	}
	m.actions.Cleanup()
	m.actions = nil
}

// Process processes the packet by making a replica and injecting it into the
// specified port. Note that the mirror action never drops a packet. The
// mirrored packet is processed inlined as the mirror action already holds the
// right references.
//
// TODO: Evaluate if we should process the mirrored packet in a
// goroutine to improve latency for the original packet.
func (m *mirror) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	err := func() error {
		if m.port == nil && m.actions == nil {
			return fmt.Errorf("actions: mirror failed, no mirror")
		}

		// Mirror the packet and replicate the specified fields.
		cp, err := packet.Mirror(m.fields)
		if err != nil {
			return err
		}

		// Apply the mirror actions on the copied packet. If the actions do not
		// fully process the packet, inject it into the port if specified. Note
		// that irrespective of the fate of the mirrored packet, we continue
		// processing the original packet.
		//
		// The port processing the mirrored packet will log the trace of the
		// mirrored packet. If the mirrored packet is not being processed by a
		// port, we need to explicitly log the packet.
		state, err := fwdaction.ProcessPacket(cp, m.actions, counters)
		if err != nil {
			fwdpacket.Log(cp)
			return fmt.Errorf("actions: mirror actions failed, err %v", err)
		}
		switch state {
		case fwdaction.OUTPUT:
			// If we don't have an output port, drop the mirrored packet.
			out, err := fwdport.OutputPort(cp, m.ctx)
			if err != nil {
				return fmt.Errorf("actions: mirror actions failed to get output port, err %v", err)
			}
			packet.Logf(fwdpacket.LogDebugMessage, "transmitting packet to %q", out.ID())
			packet.Logf(fwdpacket.LogDesc, fmt.Sprintf("%v: Transmit %v", m.ctx.ID, out.ID()))
			fwdport.Output(out, cp, fwdpb.PortAction_PORT_ACTION_OUTPUT, m.ctx)

		case fwdaction.CONTINUE:
			if m.port != nil {
				fwdport.Process(m.port, cp, m.dir, m.ctx, "mirror")
			} else {
				fwdpacket.Log(cp)
			}
		default:
			fwdpacket.Log(cp)
		}
		return nil
	}()
	// Note that if the mirror fails, we increment counters and log the error.
	// However the packet processing is continued for the original packet.
	if err != nil {
		packet.Logf(fwdpacket.LogErrorFrame, "mirrored packet failed, %v", err)
		counters.Increment(fwdpb.CounterId_MIRROR_PACKETS, 1)
		counters.Increment(fwdpb.CounterId_MIRROR_OCTETS, uint32(packet.Length()))
	}
	return nil, fwdaction.CONTINUE
}

// A mirrorBuilder builds mirror actions.
type mirrorBuilder struct{}

// init registers a builder for the mirror action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_MIRROR_ACTION, &mirrorBuilder{})
}

// Build creates a new mirror action.
func (*mirrorBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	if !proto.HasExtension(desc, fwdpb.E_MirrorActionDesc_Extension) {
		return nil, fmt.Errorf("actions: Build for mirror action failed, missing extension %s", fwdpb.E_MirrorActionDesc_Extension.Name)
	}
	mirrorExt := proto.GetExtension(desc, fwdpb.E_MirrorActionDesc_Extension).(*fwdpb.MirrorActionDesc)
	actions, err := fwdaction.NewActions(mirrorExt.GetActions(), ctx)
	if err != nil {
		return nil, fmt.Errorf("actions: Unable to create actions %v, err %v", mirrorExt.GetActions(), err)
	}
	var port fwdport.Port
	if pid := mirrorExt.GetPortId(); pid != nil {
		if port, err = fwdport.Acquire(pid, ctx); err != nil {
			actions.Cleanup()
			return nil, fmt.Errorf("actions: Build for mirror action failed, err %v", err)
		}
	}
	var fields []fwdpacket.FieldID
	for _, f := range mirrorExt.GetFieldIds() {
		fields = append(fields, fwdpacket.NewFieldID(f))
	}

	// Append the attibute for the INPUT and OUTPUT port.
	fields = append(fields, fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_INPUT, 0))
	fields = append(fields, fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_PORT_OUTPUT, 0))

	return &mirror{port: port, dir: mirrorExt.GetPortAction(), ctx: ctx, actions: actions, fields: fields}, nil
}
