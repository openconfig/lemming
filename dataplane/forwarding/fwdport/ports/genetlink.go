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
	desc       *fwdpb.PortDesc
	devName    string
	adminState fwdpb.PortState
	ctx        *fwdcontext.Context // Forwarding context containing the port
}

// Desc returns the port description proto.
func (p *genetlinkPort) Desc() *fwdpb.PortDesc {
	return p.desc
}

func (p *genetlinkPort) String() string {
	desc := fmt.Sprintf("Type=%v;DeviceName=%v;", fwdpb.PortType_PORT_TYPE_GENETLINK, p.devName)
	if state, err := p.State(nil); err == nil {
		desc += fmt.Sprintf("<State=%v>;", state)
	}
	return desc
}

func (p *genetlinkPort) Type() fwdpb.PortType {
	return fwdpb.PortType_PORT_TYPE_GENETLINK
}

func (p *genetlinkPort) Cleanup() {
}

// Update updates the actions of the port.
func (p *genetlinkPort) Update(upd *fwdpb.PortUpdateDesc) error {
	defer func() {
		p.Cleanup()
	}()
	return nil
}

func (p *genetlinkPort) Write(fwdpacket.Packet) (fwdaction.State, error) {
	return fwdaction.CONSUME, nil
}

// Actions returns the port actions of the specified type
func (p *genetlinkPort) Actions(dir fwdpb.PortAction) fwdaction.Actions {
	return nil
}

// State returns the state of the port.
func (p *genetlinkPort) State(pi *fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
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

type genetlinkBuilder struct{}

// Build creates a new port.
func (genetlinkBuilder) Build(portDesc *fwdpb.PortDesc, ctx *fwdcontext.Context) (fwdport.Port, error) {
	kp, ok := portDesc.Port.(*fwdpb.PortDesc_Genetlink)
	if !ok {
		return nil, fmt.Errorf("invalid port type in proto")
	}
	var adminState fwdpb.PortState
	if kp.Genetlink.GetAdminState() {
		adminState = fwdpb.PortState_PORT_STATE_ENABLED_UP
	} else {
		adminState = fwdpb.PortState_PORT_STATE_DISABLED_DOWN
	}
	p := genetlinkPort{
		ctx:        ctx,
		desc:       portDesc,
		adminState: adminState,
	}
	list := append(fwdport.CounterList, fwdaction.CounterList...)
	if err := p.InitCounters("", list...); err != nil {
		return nil, err
	}
	return &p, nil
}
