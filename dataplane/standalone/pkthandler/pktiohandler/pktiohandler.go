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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"github.com/openconfig/lemming/dataplane/forwarding/util/queue"
	"github.com/openconfig/lemming/dataplane/kernel"
	"github.com/openconfig/lemming/dataplane/kernel/genetlink"

	log "github.com/golang/glog"

	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
)

// New returns a new PacketIOMgr
func New(portFile string) (*PacketIOMgr, error) {
	q, err := queue.NewUnbounded("send")
	if err != nil {
		return nil, err
	}
	q.Run()
	return &PacketIOMgr{
		hostifs:           map[uint64]*port{},
		dplanePortIfIndex: map[uint64]int{},
		sendQueue:         q,
		portFile:          portFile,
	}, nil
}

// PacketIOMgr creates and delete ports and reads and writes to them.
type PacketIOMgr struct {
	hostifs           map[uint64]*port
	dplanePortIfIndex map[uint64]int // For tap devices, maps the dataport port id to hostif if index.
	sendQueue         *queue.Queue
	portFile          string
}

type port struct {
	portIO
	cancelFn func()
	msg      *pktiopb.HostPortControlMessage
}

type portIO interface {
	Delete() error
	Write([]byte, *kernel.PacketMetadata) (int, error)
	Read([]byte) (int, error)
}

// StreamPackets sends and receives packets from a lucius CPU port.
func (m *PacketIOMgr) StreamPackets(c pktiopb.PacketIO_CPUPacketStreamClient) error {
	if err := c.Send(&pktiopb.PacketIn{Msg: &pktiopb.PacketIn_Init{}}); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-c.Context().Done():
				return
			case data := <-m.sendQueue.Receive():
				if err := c.Send(&pktiopb.PacketIn{Msg: &pktiopb.PacketIn_Packet{Packet: data.(*pktiopb.Packet)}}); err != nil {
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
				log.Warningf("skipping unknown port id: %v", out.GetPacket().GetHostPort())
				continue
			}

			if _, err := port.Write(out.GetPacket().GetFrame(), m.metadataFromPacket(out.GetPacket())); err != nil {
				log.Warningf("port write err: %v", err)
				continue
			}
		}
	}
}

func (m *PacketIOMgr) metadataFromPacket(p *pktiopb.Packet) *kernel.PacketMetadata {
	md := &kernel.PacketMetadata{
		SrcIfIndex: int16(m.dplanePortIfIndex[p.GetInputPort()]),
		DstIfIndex: int16(m.dplanePortIfIndex[p.GetOutputPort()]),
	}

	return md
}

func (m *PacketIOMgr) writePorts() error {
	if m.portFile == "" {
		return nil
	}
	msg := struct {
		Hostifs map[uint64]*pktiopb.HostPortControlMessage
		PortIds map[uint64]int
	}{
		Hostifs: make(map[uint64]*pktiopb.HostPortControlMessage),
		PortIds: m.dplanePortIfIndex,
	}

	for id, h := range m.hostifs {
		msg.Hostifs[id] = h.msg
	}
	contents, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return os.WriteFile(m.portFile, contents, 0666) //nolint:gosec
}

// ManagePorts handles HostPortControl message from a forwarding server.
func (m *PacketIOMgr) ManagePorts(c pktiopb.PacketIO_HostPortControlClient) error {
	if err := c.Send(&pktiopb.HostPortControlRequest{Msg: &pktiopb.HostPortControlRequest_Init{}}); err != nil {
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
			sendErr := c.Send(&pktiopb.HostPortControlRequest{Msg: &pktiopb.HostPortControlRequest_Status{
				Status: st,
			}})
			if sendErr != nil {
				return sendErr
			}
			if err := m.writePorts(); err != nil {
				log.Warningf("failed to write file: %v", err)
			}
		} else {
			p, ok := m.hostifs[resp.GetPortId()]
			if !ok {
				sendErr := c.Send(&pktiopb.HostPortControlRequest{Msg: &pktiopb.HostPortControlRequest_Status{
					Status: &status.Status{
						Code:    int32(codes.FailedPrecondition),
						Message: fmt.Sprintf("port %v doesn't exist", resp.GetPortId()),
					},
				}})
				if sendErr != nil {
					return sendErr
				}
				continue
			}

			m.hostifs[resp.GetPortId()].cancelFn()

			if err := p.Delete(); err != nil {
				sendErr := c.Send(&pktiopb.HostPortControlRequest{Msg: &pktiopb.HostPortControlRequest_Status{
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
			sendErr := c.Send(&pktiopb.HostPortControlRequest{Msg: &pktiopb.HostPortControlRequest_Status{
				Status: &status.Status{
					Code: int32(codes.OK),
				},
			}})
			if sendErr != nil {
				return sendErr
			}
			if err := m.writePorts(); err != nil {
				log.Warningf("failed to write file: %v", err)
			}
		}
	}
}

var createTAPFunc = kernel.NewTap

func (m *PacketIOMgr) createPort(msg *pktiopb.HostPortControlMessage) error {
	var p portIO
	switch msg.GetPort().(type) {
	case *pktiopb.HostPortControlMessage_Genetlink:
		portDesc := msg.GetGenetlink()
		var err error
		p, err = genetlink.NewGenetlinkPort(portDesc.Family, portDesc.Group)
		if err != nil {
			return err
		}
		log.Infof("add to new genetlink port: %v %v", portDesc.Family, portDesc.Group)
	case *pktiopb.HostPortControlMessage_Netdev:
		name := msg.GetNetdev().GetName()
		var err error
		kp, err := createTAPFunc(name)
		if err != nil {
			return err
		}
		p = kp
		m.dplanePortIfIndex[msg.GetDataplanePort()] = kp.IfIndex()
		log.Infof("add to new netdev port: %v", name)
	default:
		return fmt.Errorf("unsupported port type: %v", msg.GetPort())
	}

	doneCh := make(chan struct{})

	m.hostifs[msg.GetPortId()] = &port{
		portIO:   p,
		cancelFn: func() { close(doneCh) },
		msg:      msg,
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
					time.Sleep(time.Millisecond)
					continue
				}
				pkt := &pktiopb.Packet{
					HostPort: id,
					Frame:    buf[0:n],
				}
				m.sendQueue.Write(pkt)
				time.Sleep(time.Millisecond)
			}
		}
	}()
}
