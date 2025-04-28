package mplsoverudpscale

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
		NexthopIPStart:      "2001:db8:1::1",
		UseSameMPLSLabel:    true,
		MPLSLabelStart:      100,
		UDPSrcPort:          5000,
		UDPDstPort:          6000,
		SrcIP:               "2001:db8:f::1",
		DstIPStart:          "2001:db8:d::1",
		NumDstIP:            2,
		DSCP:                46,
		IPTTL:               64,
	}

	tests := []struct {
		desc              string
		cfg               *ScaleProfileConfig
		wantNHs           []fluent.GRIBIEntry
		wantNHGs          []fluent.GRIBIEntry
		wantV6Prefixes    []fluent.GRIBIEntry
		wantTotalNHsAvail int
		wantSubErrStr     string
	}{
		{
			desc: "valid config",
			cfg:  validCfg,
			wantNHs: []fluent.GRIBIEntry{
				fluent.NextHopEntry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(1).
					WithIPAddress("2001:db8:1::1").
					AddEncapHeader(
						fluent.MPLSEncapHeader().WithLabels(100),
						fluent.UDPV6EncapHeader().
							WithDstUDPPort(6000).WithSrcUDPPort(5000).
							WithSrcIP("2001:db8:f::1").WithDstIP("2001:db8:d::1").
							WithDSCP(46).WithIPTTL(64),
					),
				fluent.NextHopEntry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(2).
					WithIPAddress("2001:db8:1::2").
					AddEncapHeader(
						fluent.MPLSEncapHeader().WithLabels(100),
						fluent.UDPV6EncapHeader().
							WithDstUDPPort(6000).WithSrcUDPPort(5000).
							WithSrcIP("2001:db8:f::1").WithDstIP("2001:db8:d::2").
							WithDSCP(46).WithIPTTL(64),
					),
				fluent.NextHopEntry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(3).
					WithIPAddress("2001:db8:1::3").
					AddEncapHeader(
						fluent.MPLSEncapHeader().WithLabels(100),
						fluent.UDPV6EncapHeader().
							WithDstUDPPort(6000).WithSrcUDPPort(5000).
							WithSrcIP("2001:db8:f::1").WithDstIP("2001:db8:d::1").
							WithDSCP(46).WithIPTTL(64),
					),
				fluent.NextHopEntry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(4).
					WithIPAddress("2001:db8:1::4").
					AddEncapHeader(
						fluent.MPLSEncapHeader().WithLabels(100),
						fluent.UDPV6EncapHeader().
							WithDstUDPPort(6000).WithSrcUDPPort(5000).
							WithSrcIP("2001:db8:f::1").WithDstIP("2001:db8:d::2").
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
					WithPrefix("2001:db8::/64").
					WithNextHopGroup(1),
				fluent.IPv6Entry().
					WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithPrefix("2001:db8::1/64").
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
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "2001:db8:1::1",
				SrcIP:               "2001:db8:f::1",
				DstIPStart:          "2001:db8:d::1",
				NumDstIP:            1,
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
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "2001:db8:1::1",
				SrcIP:               "2001:db8:f::1",
				DstIPStart:          "2001:db8:d::1",
				NumDstIP:            1,
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
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "2001:db8:1::1",
				SrcIP:               "2001:db8:f::1",
				DstIPStart:          "2001:db8:d::1",
				NumDstIP:            1,
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
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "2001:db8:1::1",
				SrcIP:               "2001:db8:f::1",
				DstIPStart:          "2001:db8:d::1",
				NumDstIP:            1,
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
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "2001:db8:1::1",
				SrcIP:               "2001:db8:f::1",
				DstIPStart:          "2001:db8:d::1",
				NumDstIP:            1,
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
				NexthopIPStart:      "2001:db8:1::1",
				SrcIP:               "2001:db8:f::1",
				DstIPStart:          "2001:db8:d::1",
				NumDstIP:            1,
			},
			wantSubErrStr: "invalid PrefixStart",
		},
		{
			desc: "addr family mismatch",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         10, NumNexthopGroup: 10, NumNexthopPerNHG: 1,
				PrefixStart:    "192.0.2.0/24",
				NexthopIPStart: "2001:db8:1::1",
				UDPSrcPort:     5000,
				UDPDstPort:     6000,
				SrcIP:          "::1",
				DstIPStart:     "2001:db8:d::1",
				NumDstIP:       1,
			},
			wantSubErrStr: "AddrFamily \"ipv6\" does not match PrefixStart",
		},
		{
			desc: "invalid NexthopIPStart",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "not-an-ip", // Invalid
				SrcIP:               "2001:db8:f::1",
				DstIPStart:          "2001:db8:d::1",
				NumDstIP:            1,
			},
			wantSubErrStr: "invalid NexthopIPStart",
		},
		{
			desc: "missing SrcIP",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: fakedevice.DefaultNetworkInstance,
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "2001:db8:1::1",
				DstIPStart:          "2001:db8:d::1",
				NumDstIP:            1,
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
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "2001:db8:1::1",
				SrcIP:               "2001:db8:f::1",
				NumDstIP:            1,
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
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "2001:db8:1::1",
				SrcIP:               "2001:db8:f::1",
				DstIPStart:          "2001:db8:d::1",
			},
			wantSubErrStr: "NumDstIP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := GenerateScaleProfileEntries(context.Background(), tt.cfg)

			if (err != nil) != (tt.wantSubErrStr != "") {
				t.Fatalf("Got inconsistent error: %v, want error?: %v", err, tt.wantSubErrStr == "")
			}

			if err != nil && !strings.Contains(err.Error(), tt.wantSubErrStr) {
				t.Errorf("Got error %v, want substring %s", err.Error(), tt.wantSubErrStr)
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
				// if nhgOp.GetId() != wantnhgs[i].GetId() {
				// 	t.Errorf("Got ID %v, want %v", gotnhgs[i].GetId(), wantnhgs[i].GetId())
				// }
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
