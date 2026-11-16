// Copyright 2026 Google LLC
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

package ip

import (
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestMakeIP4HeaderLengthValidation(t *testing.T) {
	ethHeader := []byte{
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, // Dst MAC
		0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, // Src MAC
		0x08, 0x00, // EtherType: IPv4
	}

	tests := []struct {
		name       string
		ipBytes    []byte
		wantErr    bool
		wantNextID fwdpb.PacketHeaderId
		wantHdrLen int
	}{
		{
			name:    "empty_frame",
			ipBytes: []byte{},
			wantErr: true,
		},
		{
			name: "ihl_0_zero_length_header",
			ipBytes: []byte{
				0x40, 0x00, 0x00, 0x14,
				0x00, 0x00, 0x00, 0x00,
			},
			wantErr: true,
		},
		{
			name: "ihl_2_truncated_8_byte_header",
			ipBytes: []byte{
				0x42, 0x00, 0x00, 0x08,
				0x00, 0x00, 0x00, 0x00,
			},
			wantErr: true,
		},
		{
			name: "ihl_4_max_invalid_16_byte_header",
			ipBytes: []byte{
				0x44, 0x00, 0x00, 0x10,
				0x00, 0x00, 0x00, 0x00,
				0x40, 0x06, 0x00, 0x00,
				0x0a, 0x00, 0x00, 0x01,
			},
			wantErr: true,
		},
		{
			name: "ihl_5_valid_minimum_20_byte_header_tcp",
			ipBytes: []byte{
				0x45, 0x00, 0x00, 0x14,
				0x00, 0x01, 0x00, 0x00,
				0x40, 0x06, 0x00, 0x00, // TTL=64, Proto=6 (TCP)
				0x0a, 0x00, 0x00, 0x01, // Src IP: 10.0.0.1
				0x0a, 0x00, 0x00, 0x02, // Dst IP: 10.0.0.2
			},
			wantErr:    false,
			wantNextID: fwdpb.PacketHeaderId_PACKET_HEADER_ID_TCP,
			wantHdrLen: 20,
		},
		{
			name: "ihl_15_60_byte_header_truncated_buffer",
			ipBytes: []byte{
				0x4f, 0x00, 0x00, 0x3c, // Version=4, IHL=15 (60 bytes), but only 20 bytes in buffer
				0x00, 0x01, 0x00, 0x00,
				0x40, 0x11, 0x00, 0x00,
				0x0a, 0x00, 0x00, 0x01,
				0x0a, 0x00, 0x00, 0x02,
			},
			wantErr: true,
		},
		{
			name: "ihl_15_valid_60_byte_header_with_options",
			ipBytes: append([]byte{
				0x4f, 0x00, 0x00, 0x3c, // Version=4, IHL=15 (60 bytes)
				0x00, 0x01, 0x00, 0x00,
				0x40, 0x11, 0x00, 0x00, // TTL=64, Proto=17 (UDP)
				0x0a, 0x00, 0x00, 0x01,
				0x0a, 0x00, 0x00, 0x02,
			}, make([]byte, 40)...), // 40 bytes of IPv4 options
			wantErr:    false,
			wantNextID: fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP,
			wantHdrLen: 60,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pkt := append(append([]byte(nil), ethHeader...), tc.ipBytes...)
			f := frame.NewFrame(pkt)
			if _, err := f.ReadHeader(len(ethHeader)); err != nil {
				t.Fatalf("failed to strip ethernet header: %v", err)
			}

			hdr, nextID, err := makeIP4(f)
			if (err != nil) != tc.wantErr {
				t.Fatalf("makeIP4() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if nextID != tc.wantNextID {
				t.Errorf("makeIP4() nextID = %v, want %v", nextID, tc.wantNextID)
			}
			if gotLen := len(hdr.Header()); gotLen != tc.wantHdrLen {
				t.Errorf("len(hdr.Header()) = %d, want %d", gotLen, tc.wantHdrLen)
			}
		})
	}
}

func FuzzParseIP(f *testing.F) {
	// Seed corpus: empty, IHL < 5, valid IPv4 (IHL=5), IPv4 with options (IHL=15),
	// IPv6 (40B), IPv4-GRE-IPv4 tunnel (with GRE key & sequence), and GRE wrapping truncated IPv4.
	f.Add([]byte{})
	f.Add([]byte{0x42, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00})
	f.Add([]byte{
		0x45, 0x00, 0x00, 0x14, 0x00, 0x01, 0x00, 0x00,
		0x40, 0x06, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x01,
		0x0a, 0x00, 0x00, 0x02,
	})
	f.Add(append([]byte{
		0x4f, 0x00, 0x00, 0x3c, 0x00, 0x01, 0x00, 0x00,
		0x40, 0x04, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x01,
		0x0a, 0x00, 0x00, 0x02,
	}, make([]byte, 40)...))
	// IPv6 base header (40 bytes, Version=6, NextHeader=17 UDP).
	f.Add(append([]byte{
		0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x40,
	}, make([]byte, 32)...))
	// Outer IPv4 (proto=47 GRE) + GRE (Key+Seq flags=0x30, EtherType=0x0800 IPv4 + 8B Key/Seq) + Inner IPv4 (20B).
	f.Add([]byte{
		0x45, 0x00, 0x00, 0x34, 0x00, 0x01, 0x00, 0x00,
		0x40, 0x2f, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x01,
		0x0a, 0x00, 0x00, 0x02,
		0x30, 0x00, 0x08, 0x00, // GRE flags (Key + Seq), Proto 0x0800 (IPv4)
		0x00, 0x00, 0x00, 0x01, // GRE Key (4B)
		0x00, 0x00, 0x00, 0x02, // GRE Sequence (4B)
		0x45, 0x00, 0x00, 0x14, 0x00, 0x01, 0x00, 0x00,
		0x40, 0x06, 0x00, 0x00, 0xc0, 0xa8, 0x01, 0x01,
		0xc0, 0xa8, 0x01, 0x02,
	})
	// Outer IPv4 (proto=47 GRE) + GRE (4B) + malformed Inner IPv4 (IHL=2).
	f.Add([]byte{
		0x45, 0x00, 0x00, 0x20, 0x00, 0x01, 0x00, 0x00,
		0x40, 0x2f, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x01,
		0x0a, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x08, 0x00,
		0x42, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00,
	})

	f.Fuzz(func(t *testing.T, data []byte) {
		h, _, err := parseIP(frame.NewFrame(data), nil)
		if err != nil {
			return
		}
		if len(h.Header()) < ip4HeaderBytes {
			t.Fatalf("parsed IP header length %d is smaller than minimum %d", len(h.Header()), ip4HeaderBytes)
		}
	})
}
