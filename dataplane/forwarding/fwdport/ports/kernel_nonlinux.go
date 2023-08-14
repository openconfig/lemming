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

//go:build !linux

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
	fwdport.Register(fwdpb.PortType_PORT_TYPE_KERNEL, kernelBuilder{})
}

// kernelPort is a ports that receives from and writes a linux network device.
type kernelPort struct {
	fwdobject.Base
}

func (p *kernelPort) String() string {
	desc := fmt.Sprintf("Type=%v", fwdpb.PortType_PORT_TYPE_KERNEL)
	if state, err := p.State(nil); err == nil {
		desc += fmt.Sprintf("<State=%v>;", state)
	}
	return desc
}

func (p *kernelPort) Cleanup() {
}

// Update updates the actions of the port.
func (p *kernelPort) Update(upd *fwdpb.PortUpdateDesc) error {
	return fmt.Errorf("kernel port only supported on linux")
}

// Write writes a packet out. If successful, the port returns
// fwdaction.CONSUME.
func (p *kernelPort) Write(packet fwdpacket.Packet) (fwdaction.State, error) {
	return fwdaction.CONSUME, fmt.Errorf("kernel port only supported on linux")
}

// Actions returns the port actions of the specified type
func (p *kernelPort) Actions(dir fwdpb.PortAction) fwdaction.Actions {
	switch dir {
	case fwdpb.PortAction_PORT_ACTION_INPUT:
		return p.input
	case fwdpb.PortAction_PORT_ACTION_OUTPUT:
		return p.output
	}
	return nil
}

// State returns the state of the port.
func (p *kernelPort) State(pi *fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	return nil, fmt.Errorf("kernel port only supported on linux")
}

type kernelBuilder struct{}

// Build creates a new port.
func (kernelBuilder) Build(portDesc *fwdpb.PortDesc, ctx *fwdcontext.Context) (fwdport.Port, error) {
	return nil, fmt.Errorf("kernel port only supported on linux")
}
