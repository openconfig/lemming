// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package pktiohandler

import (
	"context"
	"google3/third_party/golang/cmp/cmp"
	"google3/third_party/golang/grpc/codes/codes"
	"google3/third_party/golang/protobuf/v2/testing/protocmp/protocmp"
	"google3/third_party/openconfig/lemming/dataplane/kernel/kernel"
	pktiogrpc "google3/third_party/openconfig/lemming/dataplane/proto/packetio/packetio_go_grpc"
	pktiopb "google3/third_party/openconfig/lemming/dataplane/proto/packetio/packetio_go_proto"
	"io"
	"sync"
	"testing"
	"time"
)

func TestStreamPackets(t *testing.T) {
	tests := []struct {
		desc                string
		sendPackets         []*pktiopb.Packet
		recvPackets         []*pktiopb.PacketOut
		wantSentPacket      []*pktiopb.PacketIn
		wantPortWrittenData []*portWriteData
	}{{
		desc: "send packets",
		sendPackets: []*pktiopb.Packet{{
			HostPort: 1,
			Frame:    []byte("hello"),
		}},
		wantSentPacket: []*pktiopb.PacketIn{{
			Msg: &pktiopb.PacketIn_Init{},
		}, {
			Msg: &pktiopb.PacketIn_Packet{
				Packet: &pktiopb.Packet{
					HostPort: 1,
					Frame:    []byte("hello"),
				},
			},
		}},
	}, {
		desc: "recv packets",
		recvPackets: []*pktiopb.PacketOut{{
			Packet: &pktiopb.Packet{
				HostPort: 1,
				Frame:    []byte("hello"),
			},
		}},
		wantSentPacket: []*pktiopb.PacketIn{{
			Msg: &pktiopb.PacketIn_Init{},
		}},
		wantPortWrittenData: []*portWriteData{{
			Frame: []byte("hello"),
			MD:    &kernel.PacketMetadata{},
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mgr, err := New("")
			if err != nil {
				t.Fatalf("unexpected error on New(): %v", err)
			}
			fp := &fakePort{}
			mgr.hostifs[1] = &port{
				portIO:   fp,
				cancelFn: func() {},
			}
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()
			ps := &fakePacketStream{
				ctx:         ctx,
				recvPackets: tt.recvPackets,
			}
			for _, pkt := range tt.sendPackets {
				mgr.sendQueue.Write(pkt)
			}
			mgr.StreamPackets(ps)
			time.Sleep(time.Millisecond) // Sleep long enough to drain the send queue.
			if d := cmp.Diff(ps.sendPackets, tt.wantSentPacket, protocmp.Transform()); d != "" {
				t.Errorf("StreamPackets() failed: sent packet diff(-got,+want)\n:%s", d)
			}
			if d := cmp.Diff(fp.writtenData, tt.wantPortWrittenData); d != "" {
				t.Errorf("StreamPackets() failed: port packet diff(-got,+want)\n:%s", d)
			}
		})
	}
}
func TestManagePorts(t *testing.T) {
	tests := []struct {
		desc    string
		msgs    []*pktiopb.HostPortControlMessage
		wantErr string
		want    codes.Code
	}{{
		desc: "create",
		msgs: []*pktiopb.HostPortControlMessage{{
			Port: &pktiopb.HostPortControlMessage_Netdev{
				Netdev: &pktiopb.NetdevPort{
					Name: "foo1",
				},
			},
			PortId:        1,
			DataplanePort: 2,
			Create:        true,
		}},
		want: codes.OK,
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mgr, err := New("")
			if err != nil {
				t.Fatalf("unexpected error on New(): %v", err)
			}
			createTAPFunc = func(string) (*kernel.TapInterface, error) {
				return &kernel.TapInterface{}, nil
			}
			hpc := &fakeHostPortControl{
				msg: tt.msgs,
			}
			if err := mgr.ManagePorts(hpc); err != nil && err != io.EOF {
				t.Fatalf("ManagePorts() unexpected error: %v", err)
			}
			if got := codes.Code(hpc.gotReqs[1].GetStatus().GetCode()); got != tt.want {
				t.Fatalf("ManagePorts() unexpected result: got %v, want %v", got, tt.want)
			}
		})
	}
}

type portWriteData struct {
	Frame []byte
	MD    *kernel.PacketMetadata
}
type fakePort struct {
	portIO
	writtenData []*portWriteData
}

func (p *fakePort) Write(frame []byte, md *kernel.PacketMetadata) (int, error) {
	p.writtenData = append(p.writtenData, &portWriteData{Frame: frame, MD: md})
	return len(frame), nil
}

type fakePacketStream struct {
	pktiogrpc.PacketIO_CPUPacketStreamClient
	ctx            context.Context
	sMu            sync.Mutex
	sendPackets    []*pktiopb.PacketIn
	rMu            sync.Mutex
	recvPackets    []*pktiopb.PacketOut
	recvPacketsIdx int
}

func (f *fakePacketStream) Context() context.Context {
	if f.ctx == nil {
		return context.Background()
	}
	return f.ctx
}
func (f *fakePacketStream) Send(p *pktiopb.PacketIn) error {
	f.sMu.Lock()
	defer f.sMu.Unlock()
	f.sendPackets = append(f.sendPackets, p)
	return nil
}
func (f *fakePacketStream) Recv() (*pktiopb.PacketOut, error) {
	f.rMu.Lock()
	defer f.rMu.Unlock()
	if f.recvPacketsIdx >= len(f.recvPackets) {
		return nil, io.EOF
	}
	pkt := f.recvPackets[f.recvPacketsIdx]
	f.recvPacketsIdx++
	return pkt, nil
}

type fakeHostPortControl struct {
	pktiogrpc.PacketIO_HostPortControlClient
	gMu     sync.Mutex
	gotReqs []*pktiopb.HostPortControlRequest
	mMu     sync.Mutex
	msg     []*pktiopb.HostPortControlMessage
	msgIdx  int
}

func (f *fakeHostPortControl) Send(req *pktiopb.HostPortControlRequest) error {
	f.gMu.Lock()
	defer f.gMu.Unlock()
	f.gotReqs = append(f.gotReqs, req)
	return nil
}
func (f *fakeHostPortControl) Recv() (*pktiopb.HostPortControlMessage, error) {
	f.mMu.Lock()
	defer f.mMu.Unlock()
	if f.msgIdx >= len(f.msg) {
		return nil, io.EOF
	}
	req := f.msg[f.msgIdx]
	f.msgIdx++
	return req, nil
}
