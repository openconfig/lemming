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
	"io"
	"os"
	"time"

	log "github.com/golang/glog"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/internal/debug"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func init() {
	fwdport.Register(fwdpb.PortType_PORT_TYPE_FAKE, fakeBuilder{})
}

// fakePort is a ports that receives from and writes a linux network device.
type fakePort struct {
	fwdobject.Base
	port   fwdcontext.Port
	input  fwdaction.Actions
	output fwdaction.Actions
	ctx    *fwdcontext.Context // Forwarding context containing the port
}

func (p *fakePort) String() string {
	desc := fmt.Sprintf("Type=%v;<Input=%v>;<Output=%v>", fwdpb.PortType_PORT_TYPE_KERNEL, p.input, p.output)
	if state, err := p.State(nil); err == nil {
		desc += fmt.Sprintf("<State=%v>;", state)
	}
	return desc
}

func (p *fakePort) Type() fwdpb.PortType {
	return fwdpb.PortType_PORT_TYPE_FAKE
}

func (p *fakePort) Cleanup() {
	p.input.Cleanup()
	p.output.Cleanup()
	p.input = nil
	p.output = nil
}

// Update updates the actions of the port.
func (p *fakePort) Update(upd *fwdpb.PortUpdateDesc) error {
	var err error
	defer func() {
		if err != nil {
			p.Cleanup()
		}
	}()
	fakeUpd, ok := upd.Port.(*fwdpb.PortUpdateDesc_Kernel)
	if !ok {
		return fmt.Errorf("invalid type for port update")
	}

	// Acquire new actions before releasing the old ones.
	if p.input, err = fwdaction.NewActions(fakeUpd.Kernel.GetInputs(), p.ctx); err != nil {
		return fmt.Errorf("ports: input actions for port %v failed, err %v", p, err)
	}
	if p.output, err = fwdaction.NewActions(fakeUpd.Kernel.GetOutputs(), p.ctx); err != nil {
		return fmt.Errorf("ports: output actions for port %v failed, err %v", p, err)
	}
	return nil
}

func (p *fakePort) process() {
	src := gopacket.NewPacketSource(p.port, layers.LinkTypeEthernet)
	go func() {
		for {
			select {
			case pkt, ok := <-src.Packets():
				if !ok {
					log.Warning("src chan closed")
					return
				}
				fwdPkt, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, pkt.Data())
				if err != nil {
					log.Warningf("failed to create new packet: %v", err)
					log.V(1).Info(pkt.Dump())
					fwdport.Increment(p, len(pkt.Data()), fwdpb.CounterId_COUNTER_ID_RX_BAD_PACKETS, fwdpb.CounterId_COUNTER_ID_RX_BAD_OCTETS)
					continue
				}
				fwdPkt.Debug(debug.ExternalPortPacketTrace)
				fwdPkt.Log().V(2).Info("input packet", "port", p.ID(), "frame", fwdpacket.IncludeFrameInLog)
				fwdport.Process(p, fwdPkt, fwdpb.PortAction_PORT_ACTION_INPUT, p.ctx, "Kernel")
			}
		}
	}()
}

var (
	// Stubs for tests
	createFile = func(filename string) (io.Writer, error) {
		return os.Create(filename)
	}

	openFile = func(filename string) (io.Reader, error) {
		return os.Open(filename)
	}
	timeNow = time.Now
)

// Write writes a packet out. If successful, the port returns
// fwdaction.CONSUME.
func (p *fakePort) Write(packet fwdpacket.Packet) (fwdaction.State, error) {
	if err := p.port.WritePacketData(packet.Frame()); err != nil {
		return fwdaction.DROP, fmt.Errorf("failed to write eth packet: %v", err)
	}
	return fwdaction.CONSUME, nil
}

// Actions returns the port actions of the specified type
func (p *fakePort) Actions(dir fwdpb.PortAction) fwdaction.Actions {
	switch dir {
	case fwdpb.PortAction_PORT_ACTION_INPUT:
		return p.input
	case fwdpb.PortAction_PORT_ACTION_OUTPUT:
		return p.output
	}
	return nil
}

// State returns the state of the port.
func (p *fakePort) State(*fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	return &fwdpb.PortStateReply{
		Status: &fwdpb.PortInfo{
			OperStatus:  fwdpb.PortState_PORT_STATE_ENABLED_UP,
			AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
		},
	}, nil
}

type pcapManager struct {
	desc *fwdpb.FakePortDesc
}

func (mgr *pcapManager) CreatePort(string) (fwdcontext.Port, error) {
	inFile, err := openFile(mgr.desc.InFile)
	if err != nil {
		return nil, err
	}

	r, err := pcapgo.NewNgReader(inFile, pcapgo.DefaultNgReaderOptions)
	if err != nil {
		return nil, err
	}

	outFile, err := createFile(mgr.desc.OutFile)
	if err != nil {
		return nil, err
	}

	w, err := pcapgo.NewNgWriter(outFile, layers.LinkTypeEthernet)
	if err != nil {
		return nil, err
	}
	return &pcapPort{
		reader: r,
		writer: w,
	}, nil
}

type pcapPort struct {
	writer *pcapgo.NgWriter
	reader *pcapgo.NgReader
}

func (p *pcapPort) WritePacketData(data []byte) error {
	return p.writer.WritePacket(gopacket.CaptureInfo{
		Timestamp:     timeNow(),
		CaptureLength: len(data),
		Length:        len(data),
	}, data)
}

func (p *pcapPort) ReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	return p.reader.ReadPacketData()
}

type fakeBuilder struct{}

// Build creates a new port.
func (fakeBuilder) Build(portDesc *fwdpb.PortDesc, ctx *fwdcontext.Context) (fwdport.Port, error) {
	fp, ok := portDesc.Port.(*fwdpb.PortDesc_Fake)
	if !ok {
		return nil, fmt.Errorf("invalid port type in proto, got %T, expected *fwdpb.PortDesc_Fake", portDesc.Port)
	}

	var mgr fwdcontext.FakePortManager
	if ctx != nil {
		mgr = ctx.FakePortManager
	}
	if mgr == nil {
		mgr = &pcapManager{
			desc: fp.Fake,
		}
	}
	port, err := mgr.CreatePort(portDesc.GetPortId().GetObjectId().GetId())
	if err != nil {
		return nil, err
	}

	p := &fakePort{
		ctx:  ctx,
		port: port,
	}
	list := append(fwdport.CounterList, fwdaction.CounterList...)
	if err := p.InitCounters("", list...); err != nil {
		return nil, err
	}

	p.process()
	return p, nil
}
