// Copyright 2024 Google LLC
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

package pktiohandler

import (
	"errors"
	"fmt"
	"io"
	"time"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"github.com/openconfig/lemming/dataplane/forwarding/util/queue"
	"github.com/openconfig/lemming/dataplane/internal/kernel"

	log "github.com/golang/glog"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// New returns a new PacketIOMgr
func New() (*PacketIOMgr, error) {
	q, err := queue.NewUnbounded("send")
	if err != nil {
		return nil, err
	}
	q.Run()
	return &PacketIOMgr{
		hostifs:           map[uint64]*port{},
		dplanePortIfIndex: map[string]int{},
		sendQueue:         q,
	}, nil
}

// PacketIOMgr creates and delete ports and reads and writes to them.
type PacketIOMgr struct {
	hostifs           map[uint64]*port
	dplanePortIfIndex map[string]int // For tap devices, maps the dataport port id to hostif if index.
	sendQueue         *queue.Queue
}

const contextID = "lucius"

type port struct {
	portIO
	cancelFn func()
}

type portIO interface {
	Delete() error
	Write([]byte, *kernel.PacketMetadata) (int, error)
	Read([]byte) (int, error)
}

// StreamPackets sends and receives packets from a lucius CPU port.
func (m *PacketIOMgr) StreamPackets(c fwdpb.Forwarding_CPUPacketStreamClient) error {
	if err := c.Send(&fwdpb.PacketIn{Msg: &fwdpb.PacketIn_ContextId{ContextId: &fwdpb.ContextId{Id: contextID}}}); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-c.Context().Done():
				return
			case data := <-m.sendQueue.Receive():
				if err := c.Send(&fwdpb.PacketIn{Msg: &fwdpb.PacketIn_Packet{Packet: data.(*fwdpb.Packet)}}); err != nil {
					continue
				}
			}
		}
	}()

	for {
		select {
		case <-c.Context().Done():
			return c.Context().Err()
		default:
			out, err := c.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Warning("received EOF from server, exiting")
					return nil
				}
				log.Warningf("received err from server: %v", err)
				continue
			}
			port, ok := m.hostifs[out.GetPacket().GetHostPort()]
			if !ok {
				continue
			}

			if _, err := port.Write(out.GetPacket().GetFrame(), m.metadataFromPacket(out.GetPacket())); err != nil {
				continue
			}
		}
	}
}

func (m *PacketIOMgr) metadataFromPacket(p *fwdpb.Packet) *kernel.PacketMetadata {
	md := &kernel.PacketMetadata{
		SrcIfIndex: m.dplanePortIfIndex[p.GetInputPort().GetObjectId().GetId()],
		DstIfIndex: m.dplanePortIfIndex[p.GetOutputPort().GetObjectId().GetId()],
	}

	return md
}

// ManagePorts handles HostPortControl message from a forwarding server.
func (m *PacketIOMgr) ManagePorts(c fwdpb.Forwarding_HostPortControlClient) error {
	if err := c.Send(&fwdpb.HostPortControlRequest{Msg: &fwdpb.HostPortControlRequest_ContextId{ContextId: &fwdpb.ContextId{Id: contextID}}}); err != nil {
		return err
	}
	for {
		resp, err := c.Recv()
		if err != nil {
			return err
		}
		log.Infof("received port control message: %+v", resp)
		if resp.Create {
			st := &status.Status{
				Code: int32(codes.OK),
			}
			if err := m.createPort(resp); err != nil {
				st = &status.Status{
					Code:    int32(codes.Internal),
					Message: err.Error(),
				}
			}
			sendErr := c.Send(&fwdpb.HostPortControlRequest{Msg: &fwdpb.HostPortControlRequest_Status{
				Status: st,
			}})
			if sendErr != nil {
				return sendErr
			}
		} else {
			p, ok := m.hostifs[resp.GetPortId()]
			if !ok {
				sendErr := c.Send(&fwdpb.HostPortControlRequest{Msg: &fwdpb.HostPortControlRequest_Status{
					Status: &status.Status{
						Code:    int32(codes.FailedPrecondition),
						Message: fmt.Sprintf("port %v doesn't exist", resp.GetPort().GetPortId().GetObjectId().GetId()),
					},
				}})
				if sendErr != nil {
					return sendErr
				}
				continue
			}

			m.hostifs[resp.GetPortId()].cancelFn()

			if err := p.Delete(); err != nil {
				sendErr := c.Send(&fwdpb.HostPortControlRequest{Msg: &fwdpb.HostPortControlRequest_Status{
					Status: &status.Status{
						Code:    int32(codes.Internal),
						Message: err.Error(),
					},
				}})
				if sendErr != nil {
					return sendErr
				}
			}

			delete(m.hostifs, resp.GetPortId())
			sendErr := c.Send(&fwdpb.HostPortControlRequest{Msg: &fwdpb.HostPortControlRequest_Status{
				Status: &status.Status{
					Code: int32(codes.OK),
				},
			}})
			if sendErr != nil {
				return sendErr
			}
		}
	}
}

var createTAPFunc = kernel.NewTap

func (m *PacketIOMgr) createPort(msg *fwdpb.HostPortControlMessage) error {
	var p portIO
	switch msg.GetPort().GetPortType() {
	case fwdpb.PortType_PORT_TYPE_GENETLINK:
		portDesc := msg.GetPort().GetGenetlink()
		var err error
		p, err = kernel.NewGenetlinkPort(portDesc.FamilyName, portDesc.GroupName)
		if err != nil {
			return err
		}
		log.Infof("add to new genetlink port: %v %v", portDesc.FamilyName, portDesc.GroupName)
	case fwdpb.PortType_PORT_TYPE_TAP:
		name := msg.GetPort().GetTap().GetDeviceName()
		var err error
		kp, err := createTAPFunc(name)
		if err != nil {
			return err
		}
		p = kp
		m.dplanePortIfIndex[msg.GetDataplanePort().GetObjectId().GetId()] = kp.IfIndex()
		log.Infof("add to new netdev port: %v", name)
	default:
		return fmt.Errorf("unsupported port type: %q", msg.GetPort().GetPortType())
	}

	doneCh := make(chan struct{})

	m.hostifs[msg.GetPortId()] = &port{
		portIO:   p,
		cancelFn: func() { close(doneCh) },
	}

	m.queueRead(msg.GetPortId(), doneCh)

	return nil
}

func (m *PacketIOMgr) queueRead(id uint64, done chan struct{}) {
	p := m.hostifs[id]
	go func() {
		buf := make([]byte, 9100) // TODO: Configurable MTU.
		for {
			select {
			case <-done:
				return
			default:
				n, err := p.Read(buf)
				if err != nil || n == 0 {
					continue
				}
				pkt := &fwdpb.Packet{
					HostPort: id,
					Frame:    buf[0:n],
				}
				m.sendQueue.Write(pkt)
				time.Sleep(time.Microsecond)
			}
		}
	}()
}
