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
	"net"
	"os"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/internal/kernel"
	"github.com/openconfig/lemming/internal/debug"

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
	input        fwdaction.Actions
	output       fwdaction.Actions
	ctx          *fwdcontext.Context // Forwarding context containing the port
	fd           int
	devName      string
	linkDoneCh   chan struct{}
	linkUpdateCh chan netlink.LinkUpdate
	file         *os.File
	ifaceMgr     kernel.Interfaces
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
	close(p.linkDoneCh)
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
	startStateWatch(p.linkUpdateCh, p.devName, p, p.ctx)
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
			fwdPkt.Log().V(1).Info("input packet", "device", p.devName, "port", p.ID(), "frame", fwdpacket.IncludeFrameInLog)
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

func getPortState(name string, ifMgr *kernel.Interfaces) (fwdpb.PortState, fwdpb.PortState, error) {
	iface, err := ifMgr.LinkByName(name)
	if err != nil {
		return fwdpb.PortState_PORT_STATE_UNSPECIFIED, fwdpb.PortState_PORT_STATE_UNSPECIFIED, err
	}
	admin, oper := stateFromAttrs(iface.Attrs())
	return admin, oper, nil
}

func stateFromAttrs(attrs *netlink.LinkAttrs) (fwdpb.PortState, fwdpb.PortState) {
	adminState := fwdpb.PortState_PORT_STATE_DISABLED_DOWN
	if attrs.Flags&net.FlagUp != 0 {
		adminState = fwdpb.PortState_PORT_STATE_ENABLED_UP
	}

	var state fwdpb.PortState
	switch attrs.OperState {
	case netlink.OperDown:
		state = fwdpb.PortState_PORT_STATE_DISABLED_DOWN
	case netlink.OperUp, netlink.OperUnknown: // TAP interface may be unknown state because the dataplane doesn't bind to its fd, so treat unknown as up.
		state = fwdpb.PortState_PORT_STATE_ENABLED_UP
	}
	return adminState, state
}

// getAndSetState returns the state of the port and optionally sets the state to the new value.
// Note: the reply doesn't contain the updated oper-status (if applicable).
func getAndSetState(name string, ifMgr *kernel.Interfaces, pi *fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	adminState, operState, err := getPortState(name, ifMgr)
	if err != nil {
		return nil, err
	}
	reply := &fwdpb.PortStateReply{
		Status: &fwdpb.PortInfo{
			OperStatus:  operState,
			AdminStatus: adminState,
		},
	}
	if pi == nil {
		log.V(1).Infof("dataplane read port %q: admin state %v, oper state %v", name, adminState.String(), operState.String())
		return reply, nil
	}
	log.V(1).Infof("dataplane write port %q: admin state %v, oper state %v, setting to %v", name, adminState.String(), operState.String(), pi.AdminStatus.String())
	if adminState == pi.GetAdminStatus() {
		return reply, nil
	}
	if pi.AdminStatus == fwdpb.PortState_PORT_STATE_DISABLED_DOWN {
		err = ifMgr.SetState(name, false)
	} else if pi.AdminStatus == fwdpb.PortState_PORT_STATE_ENABLED_UP {
		err = ifMgr.SetState(name, true)
	}
	if err != nil {
		reply.Status.AdminStatus = pi.AdminStatus
	}
	return reply, err
}

func startStateWatch(updCh chan netlink.LinkUpdate, devName string, port fwdport.Port, ctx *fwdcontext.Context) {
	go func() {
		for {
			upd, ok := <-updCh
			if !ok {
				return
			}
			if upd.Attrs().Name != devName {
				continue
			}
			admin, oper := stateFromAttrs(upd.Attrs())
			log.V(1).Infof("dataplane receive link update: port %q, admin %v, oper %v", devName, admin.String(), oper.String())
			ctx.Notify(&fwdpb.EventDesc{
				Event: fwdpb.Event_EVENT_PORT,
				Desc: &fwdpb.EventDesc_Port{
					Port: &fwdpb.PortEventDesc{
						Context: &fwdpb.ContextId{Id: ctx.ID},
						PortId:  &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: string(port.ID())}},
						PortInfo: &fwdpb.PortInfo{
							AdminStatus: admin,
							OperStatus:  oper,
						},
					},
				},
			})
		}
	}()
}

// State returns the oper state of the port.
func (p *tapPort) State(pi *fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	return getAndSetState(p.devName, &p.ifaceMgr, pi)
}

type tapBuilder struct {
	ifaceMgr kernel.Interfaces
}

// Build creates a new port.
func (tp tapBuilder) Build(portDesc *fwdpb.PortDesc, ctx *fwdcontext.Context) (fwdport.Port, error) {
	kp, ok := portDesc.Port.(*fwdpb.PortDesc_Tap)
	if !ok {
		return nil, fmt.Errorf("invalid port type in proto")
	}
	fd, err := tp.ifaceMgr.CreateTAP(kp.Tap.GetDeviceName())
	if err != nil {
		return nil, fmt.Errorf("failed to create tap port %q: %w", kp.Tap.GetDeviceName(), err)
	}
	if err := unix.SetNonblock(fd, true); err != nil {
		return nil, fmt.Errorf("failed to set tap in non-blocking mode: %v", err)
	}
	file := os.NewFile(uintptr(fd), "/dev/tun")
	p := &tapPort{
		ctx:          ctx,
		file:         file,
		fd:           fd,
		devName:      kp.Tap.GetDeviceName(),
		linkDoneCh:   make(chan struct{}),
		linkUpdateCh: make(chan netlink.LinkUpdate),
	}
	list := append(fwdport.CounterList, fwdaction.CounterList...)
	if err := p.InitCounters("", list...); err != nil {
		return nil, err
	}
	if err := tp.ifaceMgr.LinkSubscribe(p.linkUpdateCh, p.linkDoneCh); err != nil {
		return nil, err
	}

	p.process()
	return p, nil
}
