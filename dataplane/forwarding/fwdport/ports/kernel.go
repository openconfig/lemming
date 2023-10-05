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
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/afpacket"
	"github.com/vishvananda/netlink"

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
	fwdport.Register(fwdpb.PortType_PORT_TYPE_KERNEL, kernelBuilder{})
}

// kernelPort is a ports that receives from and writes a linux network device.
type kernelPort struct {
	fwdobject.Base
	devName      string
	input        fwdaction.Actions
	output       fwdaction.Actions
	ctx          *fwdcontext.Context // Forwarding context containing the port
	handle       packetHandle
	doneCh       chan struct{}
	linkUpdateCh chan netlink.LinkUpdate
	ifaceMgr     kernel.Interfaces
}

type packetHandle interface {
	gopacket.PacketDataSource
	WritePacketData([]byte) error
	Close()
}

func (p *kernelPort) String() string {
	desc := fmt.Sprintf("Type=%v;DeviceName=%v;<Input=%v>;<Output=%v>", fwdpb.PortType_PORT_TYPE_KERNEL, p.devName, p.input, p.output)
	if state, err := p.State(nil); err == nil {
		desc += fmt.Sprintf("<State=%v>;", state)
	}
	return desc
}

func (p *kernelPort) Cleanup() {
	p.input.Cleanup()
	p.output.Cleanup()
	close(p.doneCh)
	p.input = nil
	p.output = nil
}

// Update updates the actions of the port.
func (p *kernelPort) Update(upd *fwdpb.PortUpdateDesc) error {
	var err error
	defer func() {
		if err != nil {
			p.Cleanup()
		}
	}()
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

func (p *kernelPort) process() {
	startStateWatch(p.linkUpdateCh, p.doneCh, p.devName, p, p.ctx)
	go func() {
		for {
			select {
			case <-p.doneCh:
				log.Warningf("src chan closed: %v", p.devName)
				p.handle.Close()
				return
			default:
				d, _, err := p.handle.ReadPacketData()
				if err == afpacket.ErrTimeout || err == afpacket.ErrPoll { // Don't log this error as it is very spammy.
					continue
				}
				if err != nil {
					log.Warningf("err reading packet data for %v: %v", p.devName, err)
					continue
				}
				fwdPkt, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, d)
				if err != nil {
					log.Warningf("failed to create new packet: %v", err)
					log.V(1).Info(d)
					fwdport.Increment(p, len(d), fwdpb.CounterId_COUNTER_ID_RX_BAD_PACKETS, fwdpb.CounterId_COUNTER_ID_RX_BAD_OCTETS)
					continue
				}
				fwdPkt.Debug(debug.ExternalPortPacketTrace)
				fwdPkt.Log().V(2).Info("input packet", "device", p.devName, "port", p.ID(), "frame", fwdpacket.IncludeFrameInLog)
				fwdport.Process(p, fwdPkt, fwdpb.PortAction_PORT_ACTION_INPUT, p.ctx, "Kernel")
			}
		}
	}()
}

// Write writes a packet out. If successful, the port returns
// fwdaction.CONSUME.
func (p *kernelPort) Write(packet fwdpacket.Packet) (fwdaction.State, error) {
	if err := p.handle.WritePacketData(packet.Frame()); err != nil {
		return fwdaction.DROP, fmt.Errorf("failed to write eth packet: %v", err)
	}
	return fwdaction.CONSUME, nil
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
	return getAndSetState(p.devName, &p.ifaceMgr, pi)
}

type kernelBuilder struct{}

// Build creates a new port.
func (kernelBuilder) Build(portDesc *fwdpb.PortDesc, ctx *fwdcontext.Context) (fwdport.Port, error) {
	kp, ok := portDesc.Port.(*fwdpb.PortDesc_Kernel)
	if !ok {
		return nil, fmt.Errorf("invalid port type in proto")
	}
	l, err := netlink.LinkByName(kp.Kernel.GetDeviceName())
	if err != nil {
		return nil, fmt.Errorf("failed to get interface: %v", err)
	}
	// Since the TAP port and external ports have different MACs and thus DST_MAC may be not correct,
	// read all packets for processing.
	if l.Attrs().Promisc == 0 {
		if err := netlink.SetPromiscOn(l); err != nil {
			return nil, fmt.Errorf("failed to set sec promisc on: %v", err)
		}
	}
	// Make port only reply to IPs it has.
	if err := os.WriteFile(fmt.Sprintf("/proc/sys/net/ipv4/conf/%s/arp_ignore", kp.Kernel.GetDeviceName()), []byte("2"), 0o600); err != nil {
		return nil, fmt.Errorf("failed to set arp_ignore to 2: %v", err)
	}

	// TODO: configure MTU
	handle, err := afpacket.NewTPacket(afpacket.OptInterface(kp.Kernel.GetDeviceName()), afpacket.OptPollTimeout(time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to create afpacket: %v", err)
	}
	p := &kernelPort{
		ctx:          ctx,
		handle:       handle,
		devName:      kp.Kernel.DeviceName,
		doneCh:       make(chan struct{}),
		linkUpdateCh: make(chan netlink.LinkUpdate),
	}
	list := append(fwdport.CounterList, fwdaction.CounterList...)
	if err := p.InitCounters("", list...); err != nil {
		return nil, err
	}
	if err := p.ifaceMgr.LinkSubscribe(p.linkUpdateCh, p.doneCh); err != nil {
		return nil, err
	}

	p.process()
	return p, nil
}
