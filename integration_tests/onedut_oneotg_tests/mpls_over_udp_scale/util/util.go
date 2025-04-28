package mplsoverudpscale

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/netip"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/openconfig/gribigo/fluent"
)

// ScaleProfileConfig defines the parameters for generating gRIBI entries for a specific scale profile.
type ScaleProfileConfig struct {
	// Common parameters
	AddrFamily          string // "ipv6"
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
	DstIPStart string
	NumDstIP   int
	DSCP       uint64
	IPTTL      uint64
}

// populateNextHops generates NextHop entries and returns the updated slice.
func populateNextHops(cfg *ScaleProfileConfig) ([]fluent.GRIBIEntry, error) {
	totalNextHops := cfg.NumNexthopPerNHG * cfg.NumNexthopGroup
	entries := []fluent.GRIBIEntry{}
	if totalNextHops <= 0 {
		return entries, nil
	}

	startNHIP, err := netip.ParseAddr(cfg.NexthopIPStart)
	if err != nil {
		// This should ideally not happen due to prior validation in GenerateScaleProfileEntries
		return nil, fmt.Errorf("internal error: failed to parse NexthopIPStart %q: %w", cfg.NexthopIPStart, err)
	}

	// Generate Destination IPs for Encap Header
	dstIPs := []string{}
	if cfg.NumDstIP > 0 {
		startDstIP, err := netip.ParseAddr(cfg.DstIPStart)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DstIPStart %q: %w", cfg.DstIPStart, err)
		}
		currentDstIP := startDstIP
		for i := 0; i < cfg.NumDstIP; i++ {
			if !currentDstIP.IsValid() {
				return nil, fmt.Errorf("ran out of valid destination IP addresses starting from %s after %d IPs", cfg.DstIPStart, i)
			}
			dstIPs = append(dstIPs, currentDstIP.String())
			currentDstIP = currentDstIP.Next()
		}
	} else {
		return nil, errors.New("NumDstIP must be positive to generate destination IPs")
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
				// TODO: fluent currently doesn't support UDPV4 encap header builder. Support fluent.UDPV4EncapHeader to allow ipv4 support here.
				fluent.UDPV6EncapHeader().
					WithDstUDPPort(cfg.UDPDstPort).
					WithSrcUDPPort(cfg.UDPSrcPort).
					WithSrcIP(cfg.SrcIP).
					WithDstIP(dstIPs[(i-1)%cfg.NumDstIP]).
					WithDSCP(cfg.DSCP).
					WithIPTTL(cfg.IPTTL),
			)

		entries = append(entries, nhEntry)

		currentNHIP = currentNHIP.Next()
		if !currentNHIP.IsValid() {
			// This logic might need refinement depending on expected IP range behavior
			return nil, fmt.Errorf("ran out of valid next hop IP addresses starting from %s", cfg.NexthopIPStart)
		}
	}
	return entries, nil // Return the modified slice
}

// combinationKey generates a unique, sorted string key for a slice of NH indices.
func combinationKey(indices []uint64) string {
	sort.Slice(indices, func(i, j int) bool { return indices[i] < indices[j] })

	var sb strings.Builder
	for i, idx := range indices {
		sb.WriteString(strconv.FormatUint(idx, 10))
		if i < len(indices)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

// populateNextHopGroups generates NextHopGroup entries, ensuring unique combinations of NHs per NHG.
func populateNextHopGroups(cfg *ScaleProfileConfig) ([]fluent.GRIBIEntry, error) {
	totalNHsAvailable := cfg.NumNexthopPerNHG * cfg.NumNexthopGroup
	k := cfg.NumNexthopPerNHG
	entries := []fluent.GRIBIEntry{}

	// 1. Create a slice of all available NH indices
	allNHIndices := make([]uint64, totalNHsAvailable)
	for i := 0; i < totalNHsAvailable; i++ {
		allNHIndices[i] = uint64(i + 1)
	}

	// 2. Shuffle the indices randomly
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(allNHIndices), func(i, j int) {
		allNHIndices[i], allNHIndices[j] = allNHIndices[j], allNHIndices[i]
	})

	// 3. Assign indices ensuring unique combinations
	usedCombinations := make(map[string]bool)
	currentIndex := 0 // Starting index for the sliding window in allNHIndices

	for nhgIdx := uint64(1); nhgIdx <= uint64(cfg.NumNexthopGroup); nhgIdx++ {
		foundCombination := false
		for currentIndex+k <= totalNHsAvailable {
			indSlice := allNHIndices[currentIndex : currentIndex+k]
			key := combinationKey(indSlice)

			if !usedCombinations[key] {
				usedCombinations[key] = true
				nhgEntry := fluent.NextHopGroupEntry().
					WithNetworkInstance(cfg.NetworkInstanceName).
					WithID(nhgIdx)
				finalIndSlice := make([]uint64, k)
				copy(finalIndSlice, indSlice)
				for _, nhIndex := range finalIndSlice {
					nhgEntry.AddNextHop(nhIndex, 1) // Weight is 1 for now
				}
				entries = append(entries, nhgEntry)

				currentIndex += 1
				foundCombination = true
				break
			}
			currentIndex++
		}

		if !foundCombination {
			return nil, fmt.Errorf("failed to find unique NH combination for NHG ID %d after checking all possibilities. (n=%d, k=%d, groups=%d)", nhgIdx, totalNHsAvailable, k, cfg.NumNexthopGroup)
		}
	}

	return entries, nil
}

// populatePrefixes generates IPv4 or IPv6 entries based on the configuration.
func populatePrefixes(cfg *ScaleProfileConfig) ([]fluent.GRIBIEntry, error) {
	entries := []fluent.GRIBIEntry{}

	startPrefix, err := netip.ParsePrefix(cfg.PrefixStart)
	if err != nil {
		// This should ideally not happen due to prior validation
		return nil, fmt.Errorf("internal error: failed to parse PrefixStart %q: %w", cfg.PrefixStart, err)
	}

	currentPrefix := startPrefix
	originalBits := startPrefix.Bits()

	for i := 0; i < cfg.NumPrefixes; i++ {
		prefixStr := currentPrefix.String()
		// Assign NHG ID in a round-robin fashion
		nhgID := uint64(i%cfg.NumNexthopGroup + 1)

		var entry fluent.GRIBIEntry
		switch cfg.AddrFamily {
		case "ipv6":
			entry = fluent.IPv6Entry().
				WithNetworkInstance(cfg.NetworkInstanceName).
				WithPrefix(prefixStr).
				WithNextHopGroup(nhgID)
		case "ipv4":
			// TODO: Add fluent.UDPV4EncapHeader() support first.
			return nil, fmt.Errorf("ipv4 entry generation not yet implemented")
		default:
			// Should not happen due to validation
			return nil, fmt.Errorf("internal error: unsupported AddrFamily %q", cfg.AddrFamily)
		}
		entries = append(entries, entry)

		if i < cfg.NumPrefixes-1 {
			addr := currentPrefix.Addr()
			nextAddr := addr.Next()

			if !nextAddr.IsValid() {
				return nil, fmt.Errorf("ran out of valid IP addresses after generating %d prefixes, starting from %s", i+1, cfg.PrefixStart)
			}
			if nextAddr == startPrefix.Addr() {
				return nil, fmt.Errorf("prefix generation wrapped around back to the start address %s after %d prefixes; address space likely too small for %d prefixes", startPrefix.Addr(), i+1, cfg.NumPrefixes)
			}

			// Create the next prefix using the next sequential address and the original prefix length.
			currentPrefix = netip.PrefixFrom(nextAddr, originalBits)
		}
	}

	return entries, nil
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
	if cfg.AddrFamily != "ipv6" {
		// TODO: Support ipv4 when fluent.UDPV4EncapHeader is available
		return nil, fmt.Errorf("invalid AddrFamily: %q, must be 'ipv6'", cfg.AddrFamily)
	}
	prefixAddr, err := netip.ParsePrefix(cfg.PrefixStart)
	if err != nil {
		return nil, fmt.Errorf("invalid PrefixStart %q: %w", cfg.PrefixStart, err)
	}
	if cfg.AddrFamily == "ipv6" && !prefixAddr.Addr().Is6() {
		return nil, fmt.Errorf("AddrFamily %q does not match PrefixStart %q", cfg.AddrFamily, cfg.PrefixStart)
	}
	_, err = netip.ParseAddr(cfg.NexthopIPStart)
	if err != nil {
		return nil, fmt.Errorf("invalid NexthopIPStart %q: %w", cfg.NexthopIPStart, err)
	}
	if cfg.SrcIP == "" {
		return nil, errors.New("SrcIP for encapsulation must be provided")
	}
	if cfg.DstIPStart == "" {
		return nil, errors.New("DstIPStart for encapsulation must be provided")
	}
	if cfg.NumDstIP <= 0 {
		return nil, errors.New("NumDstIP must be positive")
	}
	_, err = netip.ParseAddr(cfg.DstIPStart)
	if err != nil {
		return nil, fmt.Errorf("invalid DstIPStart %q: %w", cfg.DstIPStart, err)
	}

	// 2. Initialize basic information.
	entries := []fluent.GRIBIEntry{}

	// 3. Generates next hops.
	nhs, err := populateNextHops(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to populate next hops: %w", err)
	}
	entries = append(entries, nhs...)

	// 4. Generates next hop groups.
	nhgs, err := populateNextHopGroups(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to populate next hop groups: %w", err)
	}
	entries = append(entries, nhgs...)

	// 5. Generates prefixes (AFT entries).
	prefixEntries, err := populatePrefixes(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to populate prefixes: %w", err)
	}
	entries = append(entries, prefixEntries...)
	return entries, nil
}
