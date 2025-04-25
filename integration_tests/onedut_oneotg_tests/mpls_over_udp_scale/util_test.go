// google3/path/to/util_test.go
package mplsoverudpscale

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	gribipb "github.com/openconfig/gribi/v1/proto/service"
	"github.com/openconfig/gribigo/fluent"
	"google.golang.org/protobuf/testing/protocmp"
)

// TestGenerateScaleProfileEntries tests both validation and correct generation.
func TestGenerateScaleProfileEntries(t *testing.T) {
	// Define a valid config with specific encap details for the 'want' case
	validCfg := &ScaleProfileConfig{
		AddrFamily:          "ipv6",
		NetworkInstanceName: "DEFAULT",
		NumPrefixes:         2,
		NumNexthopGroup:     1,
		NumNexthopPerNHG:    2,
		PrefixStart:         "2001:db8::/64",
		NexthopIPStart:      "2001:db8:1::1",
		UseSameMPLSLabel:    true,
		MPLSLabelStart:      100,
		UDPSrcPort:          5000,
		UDPDstPort:          6000,
		SrcIP:               "2001:db8:f::1",
		DstIP:               "2001:db8:f::2",
		DSCP:                46,
		IPTTL:               64,
	}

	wantValidEntries := []fluent.GRIBIEntry{
		fluent.NextHopEntry().
			WithNetworkInstance("DEFAULT").
			WithIndex(1).
			WithIPAddress("2001:db8:1::1").
			AddEncapHeader(
				fluent.MPLSEncapHeader().WithLabels(100),
				fluent.UDPV6EncapHeader().
					WithDstUDPPort(6000).WithSrcUDPPort(5000).
					WithSrcIP("2001:db8:f::1").WithDstIP("2001:db8:f::2").
					WithDSCP(46).WithIPTTL(64),
			),
		fluent.NextHopEntry().
			WithNetworkInstance("DEFAULT").
			WithIndex(2).
			WithIPAddress("2001:db8:1::2").
			AddEncapHeader(
				fluent.MPLSEncapHeader().WithLabels(100),
				fluent.UDPV6EncapHeader().
					WithDstUDPPort(6000).WithSrcUDPPort(5000).
					WithSrcIP("2001:db8:f::1").WithDstIP("2001:db8:f::2").
					WithDSCP(46).WithIPTTL(64),
			),
		// TODO: Add expected NHG and AFT entries once implemented
	}

	tests := []struct {
		desc          string
		cfg           *ScaleProfileConfig
		want          []fluent.GRIBIEntry
		wantSubErrStr string
	}{
		{
			desc: "valid config",
			cfg:  validCfg,
			want: wantValidEntries,
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
				DstIP:               "2001:db8:f::2",
			},
			wantSubErrStr: "NetworkInstanceName",
		},
		{
			desc: "invalid NumPrefixes",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: "DEFAULT",
				NumPrefixes:         0, // Invalid
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "2001:db8:1::1",
				SrcIP:               "2001:db8:f::1",
				DstIP:               "2001:db8:f::2",
			},
			wantSubErrStr: "NumPrefixes",
		},
		// ... other validation error test cases ...
		{
			desc: "invalid NexthopIPStart",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: "DEFAULT",
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "not-an-ip", // Invalid
				SrcIP:               "2001:db8:f::1",
				DstIP:               "2001:db8:f::2",
			},
			wantSubErrStr: "invalid NexthopIPStart",
		},
		{
			desc: "missing SrcIP",
			cfg: &ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: "DEFAULT",
				NumPrefixes:         10,
				NumNexthopGroup:     10,
				NumNexthopPerNHG:    1,
				PrefixStart:         "2001:db8::/64",
				NexthopIPStart:      "2001:db8:1::1",
				DstIP:               "2001:db8:f::2",
			},
			wantSubErrStr: "SrcIP and DstIP",
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
			// Build protos from fluent entries for comparison
			var wantProtos []*gribipb.AFTOperation
			for i, entry := range tt.want {
				op, buildErr := entry.OpProto()
				if buildErr != nil {
					t.Fatalf("Failed to build wantEntry #%d: %v", i, buildErr)
				}
				wantProtos = append(wantProtos, op)
			}

			var gotProtos []*gribipb.AFTOperation
			for i, entry := range got {
				op, buildErr := entry.OpProto()
				if buildErr != nil {
					t.Fatalf("Failed to build gotEntry #%d: %v", i, buildErr)
				}
				gotProtos = append(gotProtos, op)
			}

			if diff := cmp.Diff(gotProtos, wantProtos, protocmp.Transform()); diff != "" {
				t.Logf("Got entries (fluent):\n%v", got)
				t.Errorf("GenerateScaleProfileEntries() returned unexpected proto diff (-got +want):\n%s", diff)
			}
			t.Logf("Got %v", got)
		})
	}
}
