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
	"os"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/internal/debug"
	"golang.org/x/sys/unix"

	log "github.com/golang/glog"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func init() {
	fwdport.Register(fwdpb.PortType_PORT_TYPE_TAP, tapBuilder{})
}

// tapPort is a ports that receives from and writes a linux network device.
// A TAP interface reads and writes from a fd instead of raw socket,
// this ensures kernel correctly responds to lower level protocols such as ARP, ICMP, etc.
type tapPort struct {
	fwdobject.Base
	input  fwdaction.Actions
	output fwdaction.Actions
	ctx    *fwdcontext.Context // Forwarding context containing the port
	fd     int
	file   *os.File
}

func (p *tapPort) String() string {
	desc := fmt.Sprintf("Type=%v;DeviceName=%v;<Input=%v>;<Output=%v>", fwdpb.PortType_PORT_TYPE_KERNEL, p.fd, p.input, p.output)
	if state, err := p.State(nil); err == nil {
		desc += fmt.Sprintf("<State=%v>;", state)
	}
	return desc
}

func (p *tapPort) Cleanup() {
	p.input.Cleanup()
	p.output.Cleanup()
	p.file.Close()
	p.input = nil
	p.output = nil
}

// Update updates the actions of the port.
func (p *tapPort) Update(upd *fwdpb.PortUpdateDesc) error {
	var err error
	defer func() {
		if err != nil {
			p.Cleanup()
		}
	}()
	// TODO: Should we have a common update type for input/output actions or an update type per port type?
	kernelUpd, ok := upd.Port.(*fwdpb.PortUpdateDesc_Kernel)
	if !ok {
		return fmt.Errorf("invalid type for port update")
	}

	// Acquire new actions before releasing the old ones.
	if p.input, err = fwdaction.NewActions(kernelUpd.Kernel.GetInputs(), p.ctx); err != nil {
		return fmt.Errorf("ports: input actions for port %v failed, err %v", p, err)
	}
	if p.output, err = fwdaction.NewActions(kernelUpd.Kernel.GetOutputs(), p.ctx); err != nil {
		return fmt.Errorf("ports: output actions for port %v failed, err %v", p, err)
	}
	return nil
}

func (p *tapPort) process() {
	go func() {
		buf := make([]byte, 1500) // TODO: MTU
		for {
			n, err := p.file.Read(buf)
			if err != nil {
				log.Warningf("failed to read packet: %v", err)
			}
			fwdPkt, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, buf[0:n])
			if err != nil {
				log.Warningf("failed to create new packet: %v", err)
				fwdport.Increment(p, n, fwdpb.CounterId_COUNTER_ID_RX_BAD_PACKETS, fwdpb.CounterId_COUNTER_ID_RX_BAD_OCTETS)
				continue
			}
			fwdPkt.Debug(debug.TAPPortPacketTrace)
			fwdport.Process(p, fwdPkt, fwdpb.PortAction_PORT_ACTION_INPUT, p.ctx, "TAP")
		}
	}()
}

// Write writes a packet out. If successful, the port returns
// fwdaction.CONSUME.
func (p *tapPort) Write(packet fwdpacket.Packet) (fwdaction.State, error) {
	if _, err := p.file.Write(packet.Frame()); err != nil {
		return fwdaction.DROP, fmt.Errorf("failed to write eth packet: %v", err)
	}
	return fwdaction.CONSUME, nil
}

// Actions returns the port actions of the specified type
func (p *tapPort) Actions(dir fwdpb.PortAction) fwdaction.Actions {
	switch dir {
	case fwdpb.PortAction_PORT_ACTION_INPUT:
		return p.input
	case fwdpb.PortAction_PORT_ACTION_OUTPUT:
		return p.output
	}
	return nil
}

// State return the state of the port (UP).
// TODO: handle port state correctly.
func (p *tapPort) State(op *fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	ready := fwdpb.PortStateReply{
		LocalPort: &fwdpb.PortInfo{
			Laser: fwdpb.PortLaserState_PORT_LASER_STATE_ENABLED,
		},
		Link: &fwdpb.LinkStateDesc{
			State: fwdpb.LinkState_LINK_STATE_UP,
			RemotePort: &fwdpb.PortInfo{
				Laser: fwdpb.PortLaserState_PORT_LASER_STATE_ENABLED,
			},
		},
	}
	return &ready, nil
}

type tapBuilder struct {
}

// Build creates a new port.
func (tapBuilder) Build(portDesc *fwdpb.PortDesc, ctx *fwdcontext.Context) (fwdport.Port, error) {
	kp, ok := portDesc.Port.(*fwdpb.PortDesc_Tap)
	if !ok {
		return nil, fmt.Errorf("invalid port type in proto")
	}
	err := unix.SetNonblock(int(kp.Tap.Fd), true)
	if err != nil {
		return nil, fmt.Errorf("failed to set tap in non-blocking mode: %v", err)
	}
	file := os.NewFile(uintptr(kp.Tap.Fd), "/dev/tun")
	p := &tapPort{
		ctx:  ctx,
		file: file,
		fd:   int(kp.Tap.Fd),
	}
	list := append(fwdport.CounterList, fwdaction.CounterList...)
	if err := p.InitCounters("", list...); err != nil {
		return nil, err
	}
	p.process()
	return p, nil
}
