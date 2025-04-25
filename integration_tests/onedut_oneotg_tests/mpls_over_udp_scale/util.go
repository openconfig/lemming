package mplsoverudpscale

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	"github.com/openconfig/gribigo/fluent"
)

// ScaleProfileConfig defines the parameters for generating gRIBI entries for a specific scale profile.
type ScaleProfileConfig struct {
	// Common parameters
	AddrFamily          string // "ipv4" or "ipv6"
	NetworkInstanceName string // Network instance name. Required.
	NumPrefixes         int    // Total number of IP prefixes (AFT entries) to generate
	NumNexthopGroup     int    // Total number of Next Hop Groups (NHGs) to generate
	NumNexthopPerNHG    int    // Number of Next Hops (NHs) per NHG

	// Prefix generation details
	PrefixStart string // Starting IP prefix (e.g., "10.5.1.1/32" or "2001:aa:bb::1/128")

	// Nexthop generation details
	NexthopIPStart string // Starting IP address for the next hops

	// MPLS Label details
	UseSameMPLSLabel bool   // If true, all NHs use MPLSLabelStart. If false, labels increment.
	MPLSLabelStart   uint64 // Starting MPLS label value

	// Encapheader Configs
	UDPSrcPort uint64
	UDPDstPort uint64
	SrcIP      string
	DstIP      string
	DSCP       uint64
	IPTTL      uint64
}

// populateNextHops generates NextHop entries and returns the updated slice.
func populateNextHops(entries []fluent.GRIBIEntry, cfg *ScaleProfileConfig) ([]fluent.GRIBIEntry, error) {
	totalNextHops := cfg.NumNexthopPerNHG * cfg.NumNexthopGroup
	if totalNextHops <= 0 {
		return entries, nil
	}

	startNHIP, err := netip.ParseAddr(cfg.NexthopIPStart)
	if err != nil {
		// This should ideally not happen due to prior validation in GenerateScaleProfileEntries
		return nil, fmt.Errorf("internal error: failed to parse NexthopIPStart %q: %w", cfg.NexthopIPStart, err)
	}

	currentNHIP := startNHIP
	for i := 1; i <= totalNextHops; i++ {
		nhIPStr := currentNHIP.String() // IP of the immediate next hop device
		mplsLabel := cfg.MPLSLabelStart
		if !cfg.UseSameMPLSLabel {
			mplsLabel = cfg.MPLSLabelStart + uint64(i-1)
		}

		nhEntry := fluent.NextHopEntry().
			WithNetworkInstance(cfg.NetworkInstanceName).
			WithIndex(uint64(i)).
			WithIPAddress(nhIPStr).
			AddEncapHeader(
				fluent.MPLSEncapHeader().WithLabels(mplsLabel),
				// NOTE: fluent currently doesn't support UDPV4 encap header builder.
				fluent.UDPV6EncapHeader().
					WithDstUDPPort(cfg.UDPDstPort).
					WithSrcUDPPort(cfg.UDPSrcPort).
					WithSrcIP(cfg.SrcIP).
					WithDstIP(cfg.DstIP).
					WithDSCP(cfg.DSCP).
					WithIPTTL(cfg.IPTTL),
			)

		entries = append(entries, nhEntry)

		currentNHIP = currentNHIP.Next()
		if !currentNHIP.IsValid() {
			// This logic might need refinement depending on expected IP range behavior
			return nil, fmt.Errorf("ran out of valid IP addresses starting from %s", cfg.NexthopIPStart)
		}
	}
	return entries, nil // Return the modified slice
}

// GenerateScaleProfileEntries randomly generates fluent gRIBI entries based on ScaleProfileConfig in one network instance.
func GenerateScaleProfileEntries(ctx context.Context, cfg *ScaleProfileConfig) ([]fluent.GRIBIEntry, error) {
	// 1. Validation of the parameters.
	if cfg.NetworkInstanceName == "" {
		return nil, errors.New("NetworkInstanceName must be given")
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
	prefixAddr, err := netip.ParsePrefix(cfg.PrefixStart)
	if err != nil {
		return nil, fmt.Errorf("invalid PrefixStart %q: %w", cfg.PrefixStart, err)
	}
	if (cfg.AddrFamily == "ipv4" && !prefixAddr.Addr().Is4()) || (cfg.AddrFamily == "ipv6" && !prefixAddr.Addr().Is6()) {
		return nil, fmt.Errorf("AddrFamily %q does not match PrefixStart %q", cfg.AddrFamily, cfg.PrefixStart)
	}
	_, err = netip.ParseAddr(cfg.NexthopIPStart)
	if err != nil {
		return nil, fmt.Errorf("invalid NexthopIPStart %q: %w", cfg.NexthopIPStart, err)
	}
	if cfg.SrcIP == "" || cfg.DstIP == "" {
		return nil, errors.New("SrcIP and DstIP for encapsulation must be provided")
	}

	// 2. Initialize basic information.
	// Initialize the slice. Give it some capacity based on expected size.
	entries := []fluent.GRIBIEntry{}

	// 3. Generates next hops.
	entries, err = populateNextHops(entries, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to populate next hops: %w", err)
	}

	return entries, nil
}
