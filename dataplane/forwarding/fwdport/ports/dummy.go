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
	fwdport.Register(fwdpb.PortType_PORT_TYPE_DUMMY, dummyBuilder{})
}

// dummyPort is a ports that drops all packets written to it. It is used to
// represent a port that is not present in the forwarding plane.
type dummyPort struct {
	fwdobject.Base
	desc       *fwdpb.PortDesc
	devName    string
	adminState fwdpb.PortState
	ctx        *fwdcontext.Context // Forwarding context containing the port
}

// Desc returns the port description proto.
func (p *dummyPort) Desc() *fwdpb.PortDesc {
	return p.desc
}

func (p *dummyPort) String() string {
	desc := fmt.Sprintf("Type=%v;DeviceName=%v;", fwdpb.PortType_PORT_TYPE_DUMMY, p.devName)
	if state, err := p.State(nil); err == nil {
		desc += fmt.Sprintf("<State=%v>;", state)
	}
	return desc
}

func (p *dummyPort) Type() fwdpb.PortType {
	return fwdpb.PortType_PORT_TYPE_DUMMY
}

func (p *dummyPort) Cleanup() {
}

// Update updates the actions of the port.
func (p *dummyPort) Update(upd *fwdpb.PortUpdateDesc) error {
	defer func() {
		p.Cleanup()
	}()
	return nil
}

func (p *dummyPort) Write(fwdpacket.Packet) (fwdaction.State, error) {
	return fwdaction.CONSUME, nil
}

// Actions returns the port actions of the specified type
func (p *dummyPort) Actions(dir fwdpb.PortAction) fwdaction.Actions {
	return nil
}

// State returns the state of the port.
func (p *dummyPort) State(pi *fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	if pi != nil {
		p.adminState = pi.AdminStatus
	}
	return &fwdpb.PortStateReply{
		Status: &fwdpb.PortInfo{
			OperStatus:  fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
			AdminStatus: p.adminState,
		},
	}, nil
}

type dummyBuilder struct{}

// Build creates a new port.
func (dummyBuilder) Build(portDesc *fwdpb.PortDesc, ctx *fwdcontext.Context) (fwdport.Port, error) {
	kp, ok := portDesc.Port.(*fwdpb.PortDesc_Dummy)
	if !ok {
		return nil, fmt.Errorf("invalid port type in proto")
	}
	var adminState fwdpb.PortState
	if kp.Dummy.GetAdminState() {
		adminState = fwdpb.PortState_PORT_STATE_ENABLED_UP
	} else {
		adminState = fwdpb.PortState_PORT_STATE_DISABLED_DOWN
	}
	p := dummyPort{
		ctx:        ctx,
		desc:       portDesc,
		devName:    kp.Dummy.GetDeviceName(),
		adminState: adminState,
	}
	list := append(fwdport.CounterList, fwdaction.CounterList...)
	if err := p.InitCounters("", list...); err != nil {
		return nil, err
	}
	return &p, nil
}
