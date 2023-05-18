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

package ports

import (
	"fmt"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/deadlock"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/util/queue"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A cpuPort is a port that receives from and transmits to the controller.
type cpuPort struct {
	fwdobject.Base
	queueID string                 // CPU queue id
	queue   *queue.Queue           // Queue of packets
	input   fwdaction.Actions      // Actions used to process received packets
	output  fwdaction.Actions      // Actions used to process transmitted packets
	ctx     *fwdcontext.Context    // Forwarding context containing the port
	export  []*fwdpb.PacketFieldId // List of fields to export when writing the packet
}

// String returns the port as a formatted string.
func (p *cpuPort) String() string {
	desc := fmt.Sprintf("Type=%v;CPU=%v;%v;<Queue=%v><Input=%v>;<Output=%v>;<Export=%v>", fwdpb.PortType_PORT_TYPE_CPU_PORT, p.queueID, p.BaseInfo(), p.queue, p.input, p.output, p.export)
	if state, err := p.State(nil); err == nil {
		desc += fmt.Sprintf("<State=%v>;", state)
	}
	return desc
}

// Cleanup releases references held by the table and its entries.
func (p *cpuPort) Cleanup() {
	p.input.Cleanup()
	p.output.Cleanup()
	p.input = nil
	p.output = nil
	p.export = nil
}

// Update updates the actions for the port.
func (p *cpuPort) Update(upd *fwdpb.PortUpdateDesc) error {
	// Release any interim actions in case of errors.
	var err error
	defer func() {
		if err != nil {
			p.Cleanup()
		}
	}()
	cpu, ok := upd.Port.(*fwdpb.PortUpdateDesc_Cpu)
	if !ok {
		return fmt.Errorf("ports: missing desc")
	}

	// Acquire new actions before releasing the old ones.
	if p.input, err = fwdaction.NewActions(cpu.Cpu.GetInputs(), p.ctx); err != nil {
		return fmt.Errorf("ports: input actions for port %v failed, err %v", p, err)
	}
	if p.output, err = fwdaction.NewActions(cpu.Cpu.GetOutputs(), p.ctx); err != nil {
		return fmt.Errorf("ports: output actions for port %v failed, err %v", p, err)
	}
	return nil
}

// Write applies output actions and writes a packet to the cable.
func (p *cpuPort) Write(packet fwdpacket.Packet) (fwdaction.State, error) {
	if err := p.queue.Write(packet); err != nil {
		return fwdaction.DROP, err
	}
	return fwdaction.CONSUME, nil
}

// punt sends a packet to the packet sink. It implements queue.Handler.
// Note that the queue handler runs in its own goroutine, and hence it must
// relock the context. We also do not want to hold the lock when performing
// the gRPC request.
func (p *cpuPort) punt(v interface{}) {
	packet, ok := v.(fwdpacket.Packet)
	if !ok {
		fwdport.Increment(p, 1, fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_TX_ERROR_OCTETS)
		return
	}

	p.ctx.RLock()
	var ingressPID *fwdpb.PortId
	if port, err := fwdport.InputPort(packet, p.ctx); err == nil {
		ingressPID = fwdport.GetID(port)
	}
	egressPID := fwdport.GetID(p)
	if port, err := fwdport.OutputPort(packet, p.ctx); err == nil {
		egressPID = fwdport.GetID(port)
	}
	var parsed []*fwdpb.PacketFieldBytes
	for _, f := range p.export {
		value, err := packet.Field(fwdpacket.NewFieldID(f))
		if err != nil {
			continue
		}
		parsed = append(parsed, &fwdpb.PacketFieldBytes{
			FieldId: f,
			Bytes:   value,
		})
	}
	request := &fwdpb.PacketInjectRequest{
		ContextId:    &fwdpb.ContextId{Id: p.ctx.ID},
		PortId:       fwdport.GetID(p),
		Egress:       egressPID,
		Ingress:      ingressPID,
		Bytes:        packet.Frame(),
		Action:       fwdpb.PortAction_PORT_ACTION_OUTPUT,
		ParsedFields: parsed,
	}

	ps := p.ctx.PacketSink()
	p.ctx.RUnlock()
	if ps != nil {
		timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Punting packet from port %v", p))
		_, err := ps(request)
		timer.Stop()
		if err == nil {
			return
		}
		log.Errorf("ports: Unable to punt packet, request %+v, err %v.", request, err)
	}
	fwdport.Increment(p, packet.Length(), fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS, fwdpb.CounterId_COUNTER_ID_TX_ERROR_OCTETS)
}

// Actions returns the port actions of the specified type.
func (p *cpuPort) Actions(dir fwdpb.PortAction) fwdaction.Actions {
	switch dir {
	case fwdpb.PortAction_PORT_ACTION_INPUT:
		return p.input
	case fwdpb.PortAction_PORT_ACTION_OUTPUT:
		return p.output
	}
	return nil
}

// State implements the port interface. The CPU port state cannot be controlled
// (it is always enabled). It is considered to be connected if a packet sink
// is present in the port's context.
func (cpuPort) State(*fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	ready := &fwdpb.PortStateReply{
		Status: &fwdpb.PortInfo{
			OperStatus:  fwdpb.PortState_PORT_STATE_ENABLED_UP,
			AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
		},
	}
	return ready, nil
}

// A cpuBuilder is used to build cpu ports.
type cpuBuilder struct{}

// init registers a builder for cpu ports.
func init() {
	fwdport.Register(fwdpb.PortType_PORT_TYPE_CPU_PORT, &cpuBuilder{})
}

// Build creates a new CPU port.
func (*cpuBuilder) Build(pd *fwdpb.PortDesc, ctx *fwdcontext.Context) (fwdport.Port, error) {
	cpu, ok := pd.Port.(*fwdpb.PortDesc_Cpu)
	if !ok {
		return nil, fmt.Errorf("ports: Unable to create cpu port")
	}

	p := cpuPort{
		ctx:     ctx,
		queueID: cpu.Cpu.GetQueueId(),
		export:  cpu.Cpu.GetExportFieldIds(),
	}
	var err error
	if l := cpu.Cpu.GetQueueLength(); l != 0 {
		p.queue, err = queue.NewBounded("punt", int(l))
	} else {
		p.queue, err = queue.NewUnbounded("punt")
	}
	if err != nil {
		return nil, fmt.Errorf("ports: Unable to create cpu port %v with length %v, err %v", p.queueID, cpu.Cpu.GetQueueLength(), err)
	}
	p.queue.Run()
	ch := make(chan bool)
	go func() {
		// Unblock the caller after goroutine is scheduled.
		log.Infof("Goroutine for queue %v scheduled", p.queueID)
		ch <- true
		for {
			v, ok := <-p.queue.Receive()
			if !ok {
				return
			}
			p.punt(v)
		}
	}()
	// Block until goroutine to drain queue is scheduled.
	<-ch

	// Store counters for all ports and actions.
	list := append(fwdport.CounterList, fwdaction.CounterList...)
	if err := p.InitCounters("", list...); err != nil {
		return nil, fmt.Errorf("cpu: Unable to initialize counters, %v", err)
	}
	return &p, nil
}
