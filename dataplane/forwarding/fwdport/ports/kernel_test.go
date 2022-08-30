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
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type fakePacketHandle struct {
	writeErr error
}

func (f *fakePacketHandle) WritePacketData(pkt []byte) error {
	return f.writeErr
}

func (f *fakePacketHandle) ReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	return nil, gopacket.CaptureInfo{}, nil
}

func (f *fakePacketHandle) Close() {}

func createEthPacket(t testing.TB) fwdpacket.Packet {
	t.Helper()
	buffer := gopacket.NewSerializeBuffer()
	src, err := net.ParseMAC("11:11:11:11:11:11")
	if err != nil {
		t.Fatal(err)
	}
	dst, err := net.ParseMAC("11:11:11:11:11:11")
	if err != nil {
		t.Fatal(err)
	}
	err = gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{},
		&layers.Ethernet{
			SrcMAC: src,
			DstMAC: dst,
		},
		gopacket.Payload([]byte("hi")),
	)
	if err != nil {
		t.Fatal(err)
	}
	p, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, buffer.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestKernelWrite(t *testing.T) {
	tests := []struct {
		desc       string
		inPacket   fwdpacket.Packet
		inWriteErr error
		wantData   []byte
		wantState  fwdaction.State
		wantErr    bool
	}{{
		desc:       "write err",
		inPacket:   createEthPacket(t),
		inWriteErr: fmt.Errorf("write err"),
		wantState:  fwdaction.DROP,
		wantErr:    true,
	}, {
		desc:      "success",
		inPacket:  createEthPacket(t),
		wantState: fwdaction.CONSUME,
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			port := kernelPort{
				handle: &fakePacketHandle{
					writeErr: tt.inWriteErr,
				},
			}
			state, err := port.Write(tt.inPacket)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Write() unexpected error: got %v, want %v", err, tt.wantErr)
			}
			if state != tt.wantState {
				t.Fatalf("Write() unexpected result: got %v, want %v", state, tt.wantState)
			}
		})
	}
}
