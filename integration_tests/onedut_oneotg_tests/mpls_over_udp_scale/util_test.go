package mplsoverudpscale

import (
	"context"
	"strings"
	"testing"
)

func TestGenerateScaleProfileAEntries_Validation(t *testing.T) {
	validCfg := &ScaleProfileConfig{
		AddrFamily:         "ipv6",
		NumNetworkInstance: 1,
		NumPrefixes:        10,
		NumNexthopGroup:    10,
		NumNexthopPerNHG:   1,
		PrefixStart:        "2001:db8::/64",
		NexthopIPStart:     "2001:db8:1::1",
		UseSameMPLSLabel:   true,
		MPLSLabelStart:     100,
	}

	tests := []struct {
		desc          string
		cfg           *ScaleProfileConfig
		wantSubErrStr string
	}{
		{
			desc: "valid config",
			cfg:  validCfg,
		},
		{
			desc: "invalid NumNetworkInstance",
			cfg: &ScaleProfileConfig{
				AddrFamily:         "ipv6",
				NumNetworkInstance: 0, // Invalid
				NumPrefixes:        10,
				NumNexthopGroup:    10,
				NumNexthopPerNHG:   1,
				PrefixStart:        "2001:db8::/64",
				NexthopIPStart:     "2001:db8:1::1",
			},
			wantSubErrStr: "NumNetworkInstance",
		},
		{
			desc: "invalid NumPrefixes",
			cfg: &ScaleProfileConfig{
				AddrFamily:         "ipv6",
				NumNetworkInstance: 1,
				NumPrefixes:        0, // Invalid
				NumNexthopGroup:    10,
				NumNexthopPerNHG:   1,
				PrefixStart:        "2001:db8::/64",
				NexthopIPStart:     "2001:db8:1::1",
			},
			wantSubErrStr: "NumPrefixes",
		},
		{
			desc: "invalid NumNexthopGroup",
			cfg: &ScaleProfileConfig{
				AddrFamily:         "ipv6",
				NumNetworkInstance: 1,
				NumPrefixes:        10,
				NumNexthopGroup:    -1, // Invalid
				NumNexthopPerNHG:   1,
				PrefixStart:        "2001:db8::/64",
				NexthopIPStart:     "2001:db8:1::1",
			},
			wantSubErrStr: "NumNexthopGroup",
		},
		{
			desc: "invalid NumNexthopPerNHG",
			cfg: &ScaleProfileConfig{
				AddrFamily:         "ipv6",
				NumNetworkInstance: 1,
				NumPrefixes:        10,
				NumNexthopGroup:    10,
				NumNexthopPerNHG:   0, // Invalid
				PrefixStart:        "2001:db8::/64",
				NexthopIPStart:     "2001:db8:1::1",
			},
			wantSubErrStr: "NumNexthopPerNHG",
		},
		{
			desc: "invalid AddrFamily",
			cfg: &ScaleProfileConfig{
				AddrFamily:         "ipvx", // Invalid
				NumNetworkInstance: 1,
				NumPrefixes:        10,
				NumNexthopGroup:    10,
				NumNexthopPerNHG:   1,
				PrefixStart:        "2001:db8::/64",
				NexthopIPStart:     "2001:db8:1::1",
			},
			wantSubErrStr: "invalid AddrFamily",
		},
		{
			desc: "invalid PrefixStart",
			cfg: &ScaleProfileConfig{
				AddrFamily:         "ipv4",
				NumNetworkInstance: 1,
				NumPrefixes:        10,
				NumNexthopGroup:    10,
				NumNexthopPerNHG:   1,
				PrefixStart:        "192.168.1.300/24", // Invalid IP
				NexthopIPStart:     "10.0.0.1",
			},
			wantSubErrStr: "invalid PrefixStart",
		},
		{
			desc: "invalid NexthopIPStart",
			cfg: &ScaleProfileConfig{
				AddrFamily:         "ipv6",
				NumNetworkInstance: 1,
				NumPrefixes:        10,
				NumNexthopGroup:    10,
				NumNexthopPerNHG:   1,
				PrefixStart:        "2001:db8::/64",
				NexthopIPStart:     "not-an-ip", // Invalid
			},
			wantSubErrStr: "invalid NexthopIPStart",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := GenerateScaleProfileAEntries(context.Background(), tt.cfg)

			if (err != nil) != (tt.wantSubErrStr != "") {
				t.Fatalf("Got inconsistent error: %v, want error?: %v", err, tt.wantSubErrStr == "")
			}

			if err != nil && !strings.Contains(err.Error(), tt.wantSubErrStr) {
				t.Errorf("Got error %v, want substring %s", err.Error(), tt.wantSubErrStr)
			}
		})
	}
}
