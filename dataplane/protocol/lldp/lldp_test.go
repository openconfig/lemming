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

package lldp

import (
	"encoding/hex"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"

	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
)

func TestIsLldp(t *testing.T) {
	tests := []struct {
		desc string
		pkt  *pktiopb.PacketOut
		want bool
	}{{
		desc: "Not LLDP.",
		pkt: &pktiopb.PacketOut{
			Packet: &pktiopb.Packet{
				HostPort: 1,
				Frame:    []byte("hello"),
			},
		},
		want: false,
	}, {
		desc: "A LLDP frame.",
		pkt: &pktiopb.PacketOut{
			Packet: &pktiopb.Packet{
				HostPort: 1,
				Frame:    []byte{0x02, 0x00, 0x02, 0x01, 0x01, 0x01, 0x10, 0x10, 0x10, 0x10, 0x10, 0x11, 0x88, 0xCC},
			},
		},
		want: false,
	}}

	for _, tt := range tests {
		got := IsLldp(tt.pkt)
		if got != tt.want {
			t.Errorf("IsLldp = %v, want %v", got, tt.want)
		}
	}
}

func TestLldpFrame(t *testing.T) {
	tests := []struct {
		desc    string
		content lldpInfo
		want    []byte
		wantErr string
	}{{
		desc: "legit lldp frame",
		content: lldpInfo{
			HostIfId:     "Ethernet1",
			SysName:      "System1",
			SysDesc:      "Sysem Desciption1",
			PortName:     "Ethernet1",
			PortDesc:     "A NIC",
			HardwareAddr: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
			Interval:     1,
		},
		want: []byte{
			0x01, 0x80, 0xc2, 0x00, 0x00, 0x0e, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x88, 0xcc, 0x02, 0x07, 0x04, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x04, 0x0a, 0x05, 0x45, 0x74, 0x68, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x31, 0x06, 0x02, 0x00, 0x02, 0x08, 0x05, 0x41, 0x20, 0x4e, 0x49, 0x43, 0x0a, 0x07, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x31, 0x0c, 0x11, 0x53, 0x79, 0x73, 0x65, 0x6d, 0x20, 0x44, 0x65, 0x73, 0x63, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x31, 0x00, 0x00,
		},
	}}

	for _, tt := range tests {
		got, gotErr := tt.content.Frame()
		if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
			t.Errorf("Want error: %+v, got: %v", tt.wantErr, gotErr)
		}
		if gotErr != nil {
			return
		}
		if d := cmp.Diff(got, tt.want); d != "" {
			t.Errorf("Wanted: %v, got: %v", tt.want, hex.EncodeToString(got))
		}
	}
}
