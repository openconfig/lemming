// Copyright 2023 Google LLC
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

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func init() {
	fwdport.Register(fwdpb.PortType_PORT_TYPE_GENETLINK, genetlinkBuilder{})
}

// genetlinkPort is a ports that receives from and writes a linux network device.
// Note: This is a "virtual" port, it doesn't perform any IO, but represents a genetlink
// port running on a remote CPU port.
type genetlinkPort struct {
	fwdobject.Base
	input  fwdaction.Actions
	output fwdaction.Actions
	ctx    *fwdcontext.Context // Forwarding context containing the port
}

func (p *genetlinkPort) String() string {
	desc := fmt.Sprintf("Type=%v;<Input=%v>;<Output=%v>", fwdpb.PortType_PORT_TYPE_KERNEL, p.input, p.output)
	if state, err := p.State(nil); err == nil {
		desc += fmt.Sprintf("<State=%v>;", state)
	}
	return desc
}

func (p *genetlinkPort) Cleanup() {
	p.input.Cleanup()
	p.output.Cleanup()
	p.input = nil
	p.output = nil
}

// Update updates the actions of the port.
func (p *genetlinkPort) Update(upd *fwdpb.PortUpdateDesc) error {
	var err error
	defer func() {
		if err != nil {
			p.Cleanup()
		}
	}()
	genUpd, ok := upd.Port.(*fwdpb.PortUpdateDesc_Genetlink)
	if !ok {
		return fmt.Errorf("invalid type for port update")
	}

	// Acquire new actions before releasing the old ones.
	if p.input, err = fwdaction.NewActions(genUpd.Genetlink.GetInputs(), p.ctx); err != nil {
		return fmt.Errorf("ports: input actions for port %v failed, err %v", p, err)
	}
	if p.output, err = fwdaction.NewActions(genUpd.Genetlink.GetOutputs(), p.ctx); err != nil {
		return fmt.Errorf("ports: output actions for port %v failed, err %v", p, err)
	}
	return nil
}

// Write writes a packet out. If successful, the port returns
// fwdaction.CONSUME.
func (p *genetlinkPort) Write(packet fwdpacket.Packet) (fwdaction.State, error) {
	return fwdaction.CONSUME, nil
}

// Actions returns the port actions of the specified type
func (p *genetlinkPort) Actions(dir fwdpb.PortAction) fwdaction.Actions {
	switch dir {
	case fwdpb.PortAction_PORT_ACTION_INPUT:
		return p.input
	case fwdpb.PortAction_PORT_ACTION_OUTPUT:
		return p.output
	}
	return nil
}

// State returns the state of the port.
func (p *genetlinkPort) State(*fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	return &fwdpb.PortStateReply{
		Status: &fwdpb.PortInfo{
			OperStatus:  fwdpb.PortState_PORT_STATE_ENABLED_UP,
			AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
		},
	}, nil
}

type genetlinkBuilder struct{}

// Build creates a new port.
func (genetlinkBuilder) Build(portDesc *fwdpb.PortDesc, ctx *fwdcontext.Context) (fwdport.Port, error) {
	return &genetlinkPort{}, nil
}
