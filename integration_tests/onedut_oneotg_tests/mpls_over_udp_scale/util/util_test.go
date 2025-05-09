// Copyright 2025 Google LLC
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

package util

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	gribipb "github.com/openconfig/gribi/v1/proto/service"
	"github.com/openconfig/gribigo/fluent"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/gnmi/fakedevice"
)

// entriesToOperationProtos accepts a slice of entries and returns next hop, next hop group, and prefix proto messages.
func entriesToOperationProtos(t *testing.T, entries []fluent.GRIBIEntry) ([]*gribipb.AFTOperation, []*gribipb.AFTOperation, []*gribipb.AFTOperation) {
	t.Helper()
	var nhs, nhgs, prefixes []*gribipb.AFTOperation
	for i, entry := range entries {
		op, err := entry.OpProto()
		if err != nil {
			t.Fatalf("Failed to build entry #%d: %v", i, err)
		}
		switch op.Entry.(type) {
		case *gribipb.AFTOperation_NextHop:
			nhs = append(nhs, op)
		case *gribipb.AFTOperation_NextHopGroup:
			nhgs = append(nhgs, op)
		case *gribipb.AFTOperation_Ipv4, *gribipb.AFTOperation_Ipv6:
			prefixes = append(prefixes, op)
		default:
			// Ignore other types for now
		}
	}
	return nhs, nhgs, prefixes
}

// TestGetFirstAddrFromPrefix tests the GetFirstAddrFromPrefix helper function.
func TestGetFirstAddrFromPrefix(t *testing.T) {
	tests := []struct {
		desc    string
		cidr    string
		want    string
		wantErr bool
	}{
		{
			desc: "valid IPv6 prefix",
			cidr: "2001:db8:abcd::/48",
			want: "2001:db8:abcd::",
		},
		{
			desc: "valid IPv6 host prefix",
			cidr: "2001:db8::1/128",
			want: "2001:db8::1",
		},
		{
			desc: "valid IPv4 prefix",
			cidr: "192.0.2.0/24",
			want: "192.0.2.0",
		},
		{
			desc: "valid IPv4 host prefix",
			cidr: "10.0.0.1/32",
			want: "10.0.0.1",
		},
		{
			desc:    "invalid CIDR",
			cidr:    "not-a-cidr",
			wantErr: true,
		},
		{
			desc:    "empty string",
			cidr:    "",
			wantErr: true,
		},
		{
			desc:    "prefix without length",
			cidr:    "2001:db8::",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := GetFirstAddrFromPrefix(tt.cidr)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFirstAddrFromPrefix(%q) got error: %v, wantErr: %v", tt.cidr, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetFirstAddrFromPrefix(%q) = %q, want %q", tt.cidr, got, tt.want)
			}
		})
	}
}

func TestGeneratePrefix(t *testing.T) {
	tests := []struct {
		name            string
		startPrefixCIDR string
		index           int
		wantPrefix      string
		wantErr         bool
	}{
		{
			name:            "ipv6 first prefix",
			startPrefixCIDR: "2001:db8:abcd::/48",
			index:           0,
			wantPrefix:      "2001:db8:abcd::/128",
			wantErr:         false,
		},
		{
			name:            "ipv6 second prefix",
			startPrefixCIDR: "2001:db8:abcd::/48",
			index:           1,
			wantPrefix:      "2001:db8:abcd::1/128",
			wantErr:         false,
		},
		{
			name:            "ipv6 large index",
			startPrefixCIDR: "2001:db8:1234::/48",
			index:           255,
			wantPrefix:      "2001:db8:1234::ff/128",
			wantErr:         false,
		},
		{
			name:            "ipv6 another large index",
			startPrefixCIDR: "2001:db8:1234::/48",
			index:           256,
			wantPrefix:      "2001:db8:1234::100/128",
			wantErr:         false,
		},
		{
			name:            "invalid start prefix",
			startPrefixCIDR: "invalid-cidr",
			index:           0,
			wantPrefix:      "",
			wantErr:         true,
		},
		// TODO: Add IPv4 tests when supported
		// {
		// 	name:            "ipv4 first prefix",
		// 	startPrefixCIDR: "192.168.1.0/24",
		// 	index:           0,
		// 	wantPrefix:      "192.168.1.0/32",
		// 	wantErr:         false,
		// },
		// {
		// 	name:            "ipv4 second prefix",
		// 	startPrefixCIDR: "192.168.1.0/24",
		// 	index:           1,
		// 	wantPrefix:      "192.168.1.1/32",
		// 	wantErr:         false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrefix, err := GeneratePrefix(tt.startPrefixCIDR, tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPrefix != tt.wantPrefix {
				t.Errorf("GeneratePrefix() = %v, want %v", gotPrefix, tt.wantPrefix)
			}
		})
	}
}

// TestFormatMPLSHeader tests the FormatMPLSHeader helper function.
func TestFormatMPLSHeader(t *testing.T) {
	tests := []struct {
		desc string
		data []byte
		want string
	}{
		{
			desc: "valid 4 bytes - label 100, exp 0, bos 1, ttl 255",
			// 100 << 12 | 0 << 9 | 1 << 8 | 255 = 0x000641FF
			data: []byte{0x00, 0x06, 0x41, 0xFF},
			want: "MPLS Label: 100, EXP: 0, BoS: true, TTL: 255",
		},
		{
			desc: "valid 4 bytes - label 200, exp 5, bos 0, ttl 64",
			// 200 << 12 | 5 << 9 | 0 << 8 | 64 = 0x000C8A40
			data: []byte{0x00, 0x0C, 0x8A, 0x40},
			want: "MPLS Label: 200, EXP: 5, BoS: false, TTL: 64",
		},
		{
			desc: "valid with payload",
			data: []byte{0x00, 0x06, 0x41, 0xFF, 0xDE, 0xAD, 0xBE, 0xEF},
			want: "MPLS Label: 100, EXP: 0, BoS: true, TTL: 255, Payload: DE AD BE EF",
		},
		{
			desc: "less than 4 bytes",
			data: []byte{0x00, 0x06, 0x41},
			want: "Invalid MPLS data (len < 4): 00 06 41",
		},
		{
			desc: "nil data",
			data: nil,
			want: "Invalid MPLS data (len < 4): ", // Note: % X with nil results in empty string
		},
		{
			desc: "empty data",
			data: []byte{},
			want: "Invalid MPLS data (len < 4): ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := FormatMPLSHeader(tt.data)
			if got != tt.want {
				t.Errorf("FormatMPLSHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestMPLSLabelToPacketBytes tests the MPLSLabelToPacketBytes helper function.
func TestMPLSLabelToPacketBytes(t *testing.T) {
	tests := []struct {
		desc  string
		label uint32
		want  []byte
	}{
		{
			desc:  "label 100",
			label: 100,
			// 100 << 12 | 1 << 8 | 255 = 0x000641FF
			want: []byte{0x00, 0x06, 0x41, 0xFF},
		},
		{
			desc:  "label 0",
			label: 0,
			// 0 << 12 | 1 << 8 | 255 = 0x000001FF
			want: []byte{0x00, 0x00, 0x01, 0xFF},
		},
		{
			desc:  "max label (20 bits)",
			label: 0xFFFFF, // 1048575
			// 0xFFFFF << 12 | 1 << 8 | 255 = 0xFFFFF1FF
			want: []byte{0xFF, 0xFF, 0xF1, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := MPLSLabelToPacketBytes(tt.label)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("MPLSLabelToPacketBytes(%d) mismatch (-want +got):\n%s", tt.label, diff)
				t.Errorf("Got bytes: % X, Want bytes: % X", got, tt.want)
			}
		})
	}
}

// TestGenerateScaleProfileEntries tests both validation and correct generation.
func TestGenerateScaleProfileEntries(t *testing.T) {
	// Define a valid config with specific encap details for the 'want' case
	validCfg := &ScaleProfileConfig{
		AddrFamily:          "ipv6",
		NetworkInstanceName: fakedevice.DefaultNetworkInstance,
		NumPrefixes:         2,
		NumNexthopGroup:     2,
		NumNexthopPerNHG:    2,
		PrefixStart:         "2001:db8::/64",
		UseSameMPLSLabel:    true,
		MPLSLabelStart:      100,
		UDPSrcPort:          5000,
		UDPDstPort:          6000,
		SrcIP:               "2001:db8:f::1",
		DstIPStart:          "2001:db8:d::1",
		NumDstIP:            2,
		DSCP:                46,
		IPTTL:               64,
		EgressATEIPv6:       "2001:db8:a::2",
	}

	tests := []struct {
		desc           string
		cfg            *ScaleProfileConfig
		wantNHs        []fluent.GRIBIEntry
		wantNHGs       []fluent.GRIBIEntry
		wantV6Prefixes []fluent.GRIBIEntry
		wantSubErrStr  string
	}{
		{
			desc: "valid config",
			cfg:  validCfg,
			wantNHs: []fluent.GRIBIEntry{
				fluent.NextHopEntry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(1).
					WithIPAddress(validCfg.EgressATEIPv6).
					AddEncapHeader(
						fluent.MPLSEncapHeader().WithLabels(100),
						fluent.UDPV6EncapHeader().
							WithDstUDPPort(6000).WithSrcUDPPort(5000).
							WithSrcIP(validCfg.SrcIP).WithDstIP("2001:db8:d::1").
							WithDSCP(46).WithIPTTL(64),
					),
				fluent.NextHopEntry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(2).
					WithIPAddress(validCfg.EgressATEIPv6).
					AddEncapHeader(
						fluent.MPLSEncapHeader().WithLabels(100),
						fluent.UDPV6EncapHeader().
							WithDstUDPPort(6000).WithSrcUDPPort(5000).
							WithSrcIP(validCfg.SrcIP).WithDstIP("2001:db8:d::2").
							WithDSCP(46).WithIPTTL(64),
					),
				fluent.NextHopEntry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(3).
					WithIPAddress(validCfg.EgressATEIPv6).
					AddEncapHeader(
						fluent.MPLSEncapHeader().WithLabels(100),
						fluent.UDPV6EncapHeader().
							WithDstUDPPort(6000).WithSrcUDPPort(5000).
							WithSrcIP(validCfg.SrcIP).WithDstIP("2001:db8:d::1").
							WithDSCP(46).WithIPTTL(64),
					),
				fluent.NextHopEntry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(4).
					WithIPAddress(validCfg.EgressATEIPv6).
					AddEncapHeader(
						fluent.MPLSEncapHeader().WithLabels(100),
						fluent.UDPV6EncapHeader().
							WithDstUDPPort(6000).WithSrcUDPPort(5000).
							WithSrcIP(validCfg.SrcIP).WithDstIP("2001:db8:d::2").
							WithDSCP(46).WithIPTTL(64),
					),
			},
			wantNHGs: []fluent.GRIBIEntry{
				fluent.NextHopGroupEntry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithID(1).
					AddNextHop(1, 1).AddNextHop(2, 2),
				fluent.NextHopGroupEntry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithID(2).
					AddNextHop(3, 1).AddNextHop(4, 2),
			},
			wantV6Prefixes: []fluent.GRIBIEntry{
				fluent.IPv6Entry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithPrefix("2001:db8::/128").
					WithNextHopGroup(1),
				fluent.IPv6Entry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithPrefix("2001:db8::1/128").
					WithNextHopGroup(2),
			},
		},
		{
			desc: "invalid NetworkInstanceName",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: "", // Invalid
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         validCfg.PrefixStart,
				SrcIP:               validCfg.SrcIP,
				DstIPStart:          validCfg.DstIPStart,
				NumDstIP:            1,
				EgressATEIPv6:       validCfg.EgressATEIPv6,
			},
			wantSubErrStr: "NetworkInstanceName",
		},
		{
			desc: "invalid NumPrefixes",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         0, // Invalid
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         validCfg.PrefixStart,
				SrcIP:               validCfg.SrcIP,
				DstIPStart:          validCfg.DstIPStart,
				NumDstIP:            1,
				EgressATEIPv6:       validCfg.EgressATEIPv6,
			},
			wantSubErrStr: "NumPrefixes",
		},
		{
			desc: "invalid NumNexthopGroup",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         1,
				NumNexthopGroup:     0, // Invalid
				NumNexthopPerNHG:    1,
				PrefixStart:         validCfg.PrefixStart,
				SrcIP:               validCfg.SrcIP,
				DstIPStart:          validCfg.DstIPStart,
				NumDstIP:            1,
				EgressATEIPv6:       validCfg.EgressATEIPv6,
			},
			wantSubErrStr: "NumNexthopGroup",
		},
		{
			desc: "invalid NumNexthopPerNHG",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         1,
				NumNexthopGroup:     1,
				NumNexthopPerNHG:    0, // Invalid
				PrefixStart:         validCfg.PrefixStart,
				SrcIP:               validCfg.SrcIP,
				DstIPStart:          validCfg.DstIPStart,
				NumDstIP:            1,
				EgressATEIPv6:       validCfg.EgressATEIPv6,
			},
			wantSubErrStr: "NumNexthopPerNHG",
		},
		{
			desc: "invalid AddrFamily",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv4", // Invalid
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         1,
				NumNexthopGroup:     1,
				NumNexthopPerNHG:    1,
				PrefixStart:         validCfg.PrefixStart,
				SrcIP:               validCfg.SrcIP,
				DstIPStart:          validCfg.DstIPStart,
				NumDstIP:            1,
				EgressATEIPv6:       validCfg.EgressATEIPv6,
			},
			wantSubErrStr: "AddrFamily",
		},
		{
			desc: "invalid PrefixStart",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         "2001:db8:::::/64", // Invalid
				SrcIP:               validCfg.SrcIP,
				DstIPStart:          validCfg.DstIPStart,
				NumDstIP:            1,
				EgressATEIPv6:       validCfg.EgressATEIPv6,
			},
			wantSubErrStr: "invalid PrefixStart",
		},
		{
			desc: "addr family mismatch",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         10, NumNexthopGroup: 10, NumNexthopPerNHG: 1,
				PrefixStart:   "192.0.2.0/24", // IPv4
				UDPSrcPort:    5000,
				UDPDstPort:    6000,
				SrcIP:         "::1",
				DstIPStart:    validCfg.DstIPStart,
				NumDstIP:      1,
				EgressATEIPv6: validCfg.EgressATEIPv6,
			},
			wantSubErrStr: "AddrFamily \"ipv6\" does not match PrefixStart",
		},
		// Removed "invalid NexthopIPStart" test case
		{
			desc: "missing SrcIP",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         validCfg.PrefixStart,
				DstIPStart:          validCfg.DstIPStart,
				NumDstIP:            1,
				EgressATEIPv6:       validCfg.EgressATEIPv6,
			},
			wantSubErrStr: "SrcIP",
		},
		{
			desc: "missing DstIPStart",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         validCfg.PrefixStart,
				SrcIP:               validCfg.SrcIP,
				NumDstIP:            1,
				EgressATEIPv6:       validCfg.EgressATEIPv6,
			},
			wantSubErrStr: "DstIPStart",
		},
		{
			desc: "missing NumDstIP",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         validCfg.PrefixStart,
				SrcIP:               validCfg.SrcIP,
				DstIPStart:          validCfg.DstIPStart,
				EgressATEIPv6:       validCfg.EgressATEIPv6,
			},
			wantSubErrStr: "NumDstIP",
		},
		{ // Added test case
			desc: "missing EgressATEIPv6",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         validCfg.PrefixStart,
				SrcIP:               validCfg.SrcIP,
				DstIPStart:          validCfg.DstIPStart,
				NumDstIP:            1,
				// EgressATEIPv6:    "", // Missing
			},
			wantSubErrStr: "EgressATEIPv6 must be provided",
		},
		{ // Added test case
			desc: "invalid EgressATEIPv6",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         validCfg.PrefixStart,
				SrcIP:               validCfg.SrcIP,
				DstIPStart:          validCfg.DstIPStart,
				NumDstIP:            1,
				EgressATEIPv6:       "not-an-ip", // Invalid
			},
			wantSubErrStr: "invalid EgressATEIPv6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := GenerateScaleProfileEntries(context.Background(), tt.cfg)

			if (err != nil) != (tt.wantSubErrStr != "") {
				t.Fatalf("Got inconsistent error: %v, want error?: %v", err, tt.wantSubErrStr == "")
			}

			if err != nil {
				if !strings.Contains(err.Error(), tt.wantSubErrStr) {
					t.Errorf("Got error %v, want substring %s", err.Error(), tt.wantSubErrStr)
				}
				return
			}

			gotnhs, gotnhgs, gotPrefixes := entriesToOperationProtos(t, got)
			wantnhs, _, _ := entriesToOperationProtos(t, tt.wantNHs)
			_, wantnhgs, _ := entriesToOperationProtos(t, tt.wantNHGs)
			var wantPrefixes []*gribipb.AFTOperation
			switch tt.cfg.AddrFamily {
			case "ipv6":
				_, _, wantPrefixes = entriesToOperationProtos(t, tt.wantV6Prefixes)
			case "ipv4":
				// TODO: add support for fluent first.
			}

			if diff := cmp.Diff(gotnhs, wantnhs, protocmp.Transform()); diff != "" {
				t.Errorf("GenerateScaleProfileEntries() returned unexpected NH proto diff (-got +want):\n%s", diff)
			}

			if len(gotnhgs) != len(tt.wantNHGs) {
				t.Errorf("Got %d NHG entries, want %d", len(gotnhgs), len(tt.wantNHGs))
			}
			gotUniqueNHMap := make(map[string]bool)
			wantUniqueNHMap := make(map[string]bool)
			var gotNHInds, wantNHInds []uint64
			for i, nhgOp := range gotnhgs {
				nhg := nhgOp.GetNextHopGroup().GetNextHopGroup()
				wantnhg := wantnhgs[i].GetNextHopGroup().GetNextHopGroup()
				wantnhs := wantnhg.GetNextHop()
				if len(nhg.GetNextHop()) != len(wantnhg.GetNextHop()) {
					t.Errorf("Got %d next hops , want %d", len(nhg.GetNextHop()), len(wantnhg.GetNextHop()))
				}
				for i, nh := range nhg.GetNextHop() {
					gotNHInds = append(gotNHInds, nh.GetIndex())
					wantNHInds = append(wantNHInds, wantnhs[i].GetIndex())
					gotKey := combinationKey(gotNHInds)
					wantKey := combinationKey(wantNHInds)
					gotUniqueNHMap[gotKey] = true
					wantUniqueNHMap[wantKey] = true
				}
				if len(gotUniqueNHMap) != len(wantUniqueNHMap) {
					t.Errorf("Got length of unique NH map: %d, want %d", len(gotUniqueNHMap), len(wantUniqueNHMap))
				}

				nhg.NextHop = nil
				wantnhg.NextHop = nil
			}
			if diff := cmp.Diff(gotnhgs, wantnhgs, protocmp.Transform()); diff != "" {
				t.Errorf("GenerateScaleProfileEntries() returned unexpected NHG proto diff (-got +want):\n%s", diff)
			}

			if diff := cmp.Diff(gotPrefixes, wantPrefixes, protocmp.Transform()); diff != "" {
				t.Errorf("GenerateScaleProfileEntries() returned unexpected Prefixes proto diff (-got +want):\n%s", diff)
			}
		})
	}
}
