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
	"context"
	"io"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/dataplane/internal/kernel"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestStreamPackets(t *testing.T) {
	tests := []struct {
		desc                string
		sendPackets         []*fwdpb.Packet
		recvPackets         []*fwdpb.PacketOut
		wantSentPacket      []*fwdpb.PacketIn
		wantPortWrittenData []*portWriteData
	}{{
		desc: "send packets",
		sendPackets: []*fwdpb.Packet{{
			HostPort: 1,
			Frame:    []byte("hello"),
		}},
		wantSentPacket: []*fwdpb.PacketIn{{
			Msg: &fwdpb.PacketIn_ContextId{ContextId: &fwdpb.ContextId{Id: contextID}},
		}, {
			Msg: &fwdpb.PacketIn_Packet{
				Packet: &fwdpb.Packet{
					HostPort: 1,
					Frame:    []byte("hello"),
				},
			},
		}},
	}, {
		desc: "recv packets",
		recvPackets: []*fwdpb.PacketOut{{
			Packet: &fwdpb.Packet{
				HostPort: 1,
				Frame:    []byte("hello"),
			},
		}},
		wantSentPacket: []*fwdpb.PacketIn{{
			Msg: &fwdpb.PacketIn_ContextId{ContextId: &fwdpb.ContextId{Id: contextID}},
		}},
		wantPortWrittenData: []*portWriteData{{
			Frame: []byte("hello"),
			MD:    &kernel.PacketMetadata{},
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mgr, err := New()
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
		msgs    []*fwdpb.HostPortControlMessage
		wantErr string
		want    codes.Code
	}{{
		desc: "create",
		msgs: []*fwdpb.HostPortControlMessage{{
			Port: &fwdpb.PortDesc{
				PortType: fwdpb.PortType_PORT_TYPE_TAP,
				Port: &fwdpb.PortDesc_Tap{
					Tap: &fwdpb.TAPPortDesc{
						DeviceName: "foo1",
					},
				},
			},
			PortId:        1,
			DataplanePort: &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "2"}},
			Create:        true,
		}},
		want: codes.OK,
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mgr, err := New()
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
	fwdpb.Forwarding_CPUPacketStreamClient
	sendPackets    []*fwdpb.PacketIn
	ctx            context.Context
	recvPackets    []*fwdpb.PacketOut
	recvPacketsIdx int
}

func (f *fakePacketStream) Context() context.Context {
	if f.ctx == nil {
		return context.Background()
	}
	return f.ctx
}

func (f *fakePacketStream) Send(p *fwdpb.PacketIn) error {
	f.sendPackets = append(f.sendPackets, p)
	return nil
}

func (f *fakePacketStream) Recv() (*fwdpb.PacketOut, error) {
	if f.recvPacketsIdx >= len(f.recvPackets) {
		return nil, io.EOF
	}
	pkt := f.recvPackets[f.recvPacketsIdx]
	f.recvPacketsIdx++

	return pkt, nil
}

type fakeHostPortControl struct {
	fwdpb.Forwarding_HostPortControlClient
	gotReqs []*fwdpb.HostPortControlRequest
	msg     []*fwdpb.HostPortControlMessage
	msgIdx  int
}

func (f *fakeHostPortControl) Send(req *fwdpb.HostPortControlRequest) error {
	f.gotReqs = append(f.gotReqs, req)
	return nil
}

func (f *fakeHostPortControl) Recv() (*fwdpb.HostPortControlMessage, error) {
	if f.msgIdx >= len(f.msg) {
		return nil, io.EOF
	}
	req := f.msg[f.msgIdx]
	f.msgIdx++

	return req, nil
}
