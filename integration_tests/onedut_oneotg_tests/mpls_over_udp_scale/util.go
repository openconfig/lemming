package mplsoverudpscale

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	gribipb "github.com/openconfig/gribi/v1/proto/service"
)

// ScaleProfileConfig defines the parameters for generating gRIBI entries for a specific scale profile.
type ScaleProfileConfig struct {
	// Common parameters
	AddrFamily         string // "ipv4" or "ipv6"
	NumNetworkInstance int    // Number of network instances (VRFs), MUST be > 1. If 1, "DEFAULT" will be applied.
	NumPrefixes        int    // Total number of IP prefixes (AFT entries) to generate
	NumNexthopGroup    int    // Total number of Next Hop Groups (NHGs) to generate
	NumNexthopPerNHG   int    // Number of Next Hops (NHs) per NHG

	// Prefix generation details
	PrefixStart string // Starting IP prefix (e.g., "10.5.1.1/32" or "2001:aa:bb::1/128")

	// Nexthop generation details
	NexthopIPStart string // Starting IP address for the next hops

	// MPLS Label details
	UseSameMPLSLabel bool   // If true, all NHs use MPLSLabelStart. If false, labels increment.
	MPLSLabelStart   uint32 // Starting MPLS label value

	// Network Instance details
	BaseNetworkInstanceName string // Base name for network instances (e.g., "vrf-") if NumNetworkInstance > 1. Default VRF often implied if 1.
}

// GenerateScaleProfileAEntries generates gRIBI ModifyRequest operations based on ScaleProfileConfig.
func GenerateScaleProfileAEntries(ctx context.Context, cfg *ScaleProfileConfig) ([]*gribipb.ModifyRequest, error) {
	// 1. Validation of the parameters.
	if cfg.NumNetworkInstance < 1 {
		return nil, errors.New("NumNetworkInstance must be at least 1")
	}
	if cfg.NumPrefixes <= 0 {
		return nil, errors.New("NumPrefixes must be positive")
	}
	if cfg.NumNexthopGroup <= 0 {
		return nil, errors.New("NumNexthopGroup must be positive")
	}
	if cfg.NumNexthopPerNHG <= 0 {
		return nil, errors.New("NumNexthopPerNHG must be positive")
	}
	if cfg.AddrFamily != "ipv4" && cfg.AddrFamily != "ipv6" {
		return nil, fmt.Errorf("invalid AddrFamily: %q, must be 'ipv4' or 'ipv6'", cfg.AddrFamily)
	}
	if _, err := netip.ParsePrefix(cfg.PrefixStart); err != nil {
		return nil, fmt.Errorf("invalid PrefixStart %q: %w", cfg.PrefixStart, err)
	}
	if _, err := netip.ParseAddr(cfg.NexthopIPStart); err != nil {
		return nil, fmt.Errorf("invalid NexthopIPStart %q: %w", cfg.NexthopIPStart, err)
	}

	// 2. Initialize gRIBI operations slice.
	var modifyReqs []*gribipb.ModifyRequest

	// 3. Loop NumPrefixes (or NumNexthopGroup) times:
	//    (Implementation pending)

	return modifyReqs, nil
}
