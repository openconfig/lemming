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

//go:build linux

package ports

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/sys/unix"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/kernel"
	"github.com/openconfig/lemming/internal/debug"

	log "github.com/golang/glog"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const packetBufferSize = 16384 // 16KB buffer

func init() {
	fwdport.Register(fwdpb.PortType_PORT_TYPE_TAP, tapBuilder{})
}

// tapPort is a ports that receives from and writes a linux network device.
// A TAP interface reads and writes from a fd instead of raw socket,
// this ensures kernel correctly responds to lower level protocols such as ARP, ICMP, etc.
type tapPort struct {
	fwdobject.Base
	input    fwdaction.Actions
	output   fwdaction.Actions
	desc     *fwdpb.PortDesc
	ctx      *fwdcontext.Context // Forwarding context containing the port
	fd       int
	devName  string
	doneCh   chan struct{}
	file     *os.File
	ifaceMgr kernel.Interfaces
}

// Desc returns the port description proto.
func (p *tapPort) Desc() *fwdpb.PortDesc {
	return p.desc
}

func (p *tapPort) String() string {
	desc := fmt.Sprintf("Type=%v;DeviceName=%v;<Input=%v>;<Output=%v>", fwdpb.PortType_PORT_TYPE_KERNEL, p.fd, p.input, p.output)
	if state, err := p.State(nil); err == nil {
		desc += fmt.Sprintf("<State=%v>;", state)
	}
	return desc
}

func (p *tapPort) Type() fwdpb.PortType {
	return fwdpb.PortType_PORT_TYPE_TAP
}

func (p *tapPort) Cleanup() {
	p.input.Cleanup()
	p.output.Cleanup()
	defaultLinkWatcher.unwatch(p.devName)
	close(p.doneCh)
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
		buf := make([]byte, packetBufferSize)
		for {
			select {
			case <-p.doneCh:
				log.Warningf("stopping reading tap packet to read packet: %v", p.devName)
				p.file.Close()
				return
			default:
				p.file.SetReadDeadline(time.Now().Add(time.Second))
				n, err := p.file.Read(buf)
				if errors.Is(err, os.ErrDeadlineExceeded) { // Don't log this error as it is very spammy.
					continue
				}
				if err != nil {
					log.Warningf("failed to read packet: %v", err)
					time.Sleep(2 * time.Second)
					continue
				}
				fwdPkt, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, buf[0:n])
				if err != nil {
					log.Warningf("failed to create new packet: %v", err)
					fwdport.Increment(p, n, fwdpb.CounterId_COUNTER_ID_RX_BAD_PACKETS, fwdpb.CounterId_COUNTER_ID_RX_BAD_OCTETS)
					continue
				}
				fwdPkt.Debug(debug.TAPPortPacketTrace)
				fwdPkt.Log().V(2).Info("input packet", "device", p.devName, "port", p.ID(), "frame", fwdpacket.IncludeFrameInLog)
				fwdport.Process(p, fwdPkt, fwdpb.PortAction_PORT_ACTION_INPUT, p.ctx, "TAP")
			}
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
	l, err := tp.ifaceMgr.LinkByName(kp.Tap.GetDeviceName())
	if err != nil {
		return nil, fmt.Errorf("failed to get tap interface %q: %w", kp.Tap.GetDeviceName(), err)
	}
	isUp := l.Attrs().Flags&net.FlagUp != 0
	if kp.Tap.GetAdminState() && !isUp {
		if err := tp.ifaceMgr.LinkSetUp(l); err != nil {
			return nil, fmt.Errorf("failed to set tap interface %q up: %w", kp.Tap.GetDeviceName(), err)
		}
	} else if !kp.Tap.GetAdminState() && isUp {
		if err := tp.ifaceMgr.LinkSetDown(l); err != nil {
			return nil, fmt.Errorf("failed to set tap interface %q down: %w", kp.Tap.GetDeviceName(), err)
		}
	}
	file := os.NewFile(uintptr(fd), "/dev/tun")
	p := &tapPort{
		ctx:     ctx,
		file:    file,
		fd:      fd,
		devName: kp.Tap.GetDeviceName(),
		desc:    portDesc,
		doneCh:  make(chan struct{}),
	}
	list := append(fwdport.CounterList, fwdaction.CounterList...)
	if err := p.InitCounters("", list...); err != nil {
		return nil, err
	}
	if err := defaultLinkWatcher.watch(p.devName, p, p.ctx); err != nil {
		return nil, err
	}

	p.process()
	return p, nil
}
