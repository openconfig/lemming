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

package registry

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/gnmi/errdiff"

	"github.com/openconfig/lemming/dataplane/proto/packetio"
	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
)

// TestRegistration tests the protocol registration and deregistration.
func TestRegistration(t *testing.T) {
	tests := []struct {
		desc         string
		action       string
		protocolName string
		wantErr      string
	}{{
		desc:         "Register new protocol pass",
		action:       "register",
		protocolName: "new",
	}, {
		desc:         "Register existing protocol fail",
		action:       "register",
		protocolName: "M3",
		wantErr:      "is existing",
	}, {
		desc:         "Deregister existing protocol pass",
		action:       "deregister",
		protocolName: "M3",
	}, {
		desc:         "Deregister nonexisting protocol fail",
		action:       "deregister",
		protocolName: "non-existing",
		wantErr:      "not found",
	}}
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	packet := &fakePacketStream{ctx: ctx}
	reg, err := NewRegistry(packet)
	if err != nil {
		t.Fatalf("Failed to create protocol register: %v", err)
	}

	op := newMul3Protocol()
	err = reg.Register("M3", op)
	if err != nil {
		t.Fatalf("Failed to register mul3 protocol: %v.", err)
	}
	ep := newMul5Protocol()
	err = reg.Register("M5", ep)
	if err != nil {
		t.Fatalf("Failed to register mul5 protocol: %v.", err)
	}

	for _, tt := range tests {
		switch tt.action {
		case "register":
			gotErr := reg.Register(tt.protocolName, nil)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("Unexpected registration err: %s w/ error %v", diff, gotErr)
			}
		case "deregister":
			gotErr := reg.Deregister(tt.protocolName)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("Unexpected deregistration err: %s w/ error %v", diff, gotErr)
			}
		}
	}
}

// TestMatching tests the interaction between the registry and procotol handlers.
func TestMatching(t *testing.T) {
	m3Pkt1 := &pktiopb.PacketOut{
		Packet: &pktiopb.Packet{
			HostPort: 3,
		},
	}
	m3Pkt2 := &pktiopb.PacketOut{
		Packet: &pktiopb.Packet{
			HostPort: 6,
		},
	}
	m5Pkt := &pktiopb.PacketOut{
		Packet: &pktiopb.Packet{
			HostPort: 5,
		},
	}
	m7Pkt := &pktiopb.PacketOut{
		Packet: &pktiopb.Packet{
			HostPort: 7,
		},
	}
	tests := []struct {
		desc    string
		packets []*pktiopb.PacketOut
		wantM3  []*pktiopb.PacketOut
		wantM5  []*pktiopb.PacketOut
		wantErr string
	}{{
		desc:    "Passing case",
		packets: []*pktiopb.PacketOut{m7Pkt, m3Pkt1, m3Pkt2, m5Pkt},
		wantM3:  []*pktiopb.PacketOut{m3Pkt1, m3Pkt2},
		wantM5:  []*pktiopb.PacketOut{m5Pkt},
	}}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	for _, tt := range tests {
		ps := &fakePacketStream{
			ctx:         ctx,
			recvPackets: tt.packets,
		}
		reg, err := NewRegistry(ps)
		if err != nil {
			t.Fatalf("Failed to create protocol register: %v", err)
		}
		m3 := newMul3Protocol()
		err = reg.Register("m3", m3)
		if err != nil {
			t.Errorf("Failed to register m3 protocol.")
		}
		m5 := newMul5Protocol()
		err = reg.Register("m5", m5)
		if err != nil {
			t.Errorf("Failed to register m5 protocol.")
		}
		reg.Start()
		time.Sleep(100 * time.Millisecond)

		if d := cmp.Diff(m3.recvPackets, tt.wantM3, protocmp.Transform()); d != "" {
			t.Errorf("Unexpected received packet for protocol m3: diff(-got,+want)\n:%s", d)
		}
		if d := cmp.Diff(m5.recvPackets, tt.wantM5, protocmp.Transform()); d != "" {
			t.Errorf("Unexpected received packet for protocol m5: diff(-got,+want)\n:%s", d)
		}
	}
}

// mul3Protocol only processes packet with port number that is multiple of 3.
type mul3Protocol struct {
	recvPackets []*pktiopb.PacketOut
	grpc.ClientStream
}

func newMul3Protocol() *mul3Protocol {
	return &mul3Protocol{}
}

func (d *mul3Protocol) Matched(p *packetio.PacketOut) bool {
	if p.Packet.GetHostPort()%3 == 0 {
		return true
	}
	return false
}

func (d *mul3Protocol) Process(p *pktiopb.PacketOut) error {
	d.recvPackets = append(d.recvPackets, p)
	return nil
}

func (d *mul3Protocol) Start() {
}

func (d *mul3Protocol) Stop() {
}

// mul5Protocol only processes packet with port number that is multiple of 5.
type mul5Protocol struct {
	mul3Protocol
}

func newMul5Protocol() *mul5Protocol {
	return &mul5Protocol{}
}

func (d *mul5Protocol) Matched(p *packetio.PacketOut) bool {
	if p.Packet.GetHostPort()%5 == 0 {
		return true
	}
	return false
}

type fakePacketStream struct {
	pktiopb.PacketIO_CPUPacketStreamClient
	ctx            context.Context
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
	return nil
}

func (f *fakePacketStream) Recv() (*pktiopb.PacketOut, error) {
	if f.recvPacketsIdx >= len(f.recvPackets) {
		return nil, io.EOF
	}
	pkt := f.recvPackets[f.recvPacketsIdx]
	f.recvPacketsIdx++
	return pkt, nil
}
