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

	// The IP address of the ATE interface on the egress link.
	EgressATEIPv6 string // e.g., atePort2.IPv6
}

// GetFirstAddrFromPrefix takes a CIDR string (e.g., "2001:aa:bb::/48")
// and returns the first usable address as a string (e.g., "2001:aa:bb::").
func GetFirstAddrFromPrefix(cidr string) (string, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return "", fmt.Errorf("failed to parse prefix %q: %w", cidr, err)
	}
	// For a prefix, Addr() returns the network address itself, which is the first address.
	return prefix.Addr().String(), nil
}

// GeneratePrefix calculates the specific prefix string based on a start prefix and an index offset.
func GeneratePrefix(startPrefixCIDR string, index int) (string, error) {
	startPrefix, err := netip.ParsePrefix(startPrefixCIDR)
	if err != nil {
		return "", fmt.Errorf("failed to parse start prefix %q: %w", startPrefixCIDR, err)
	}
	startAddr := startPrefix.Addr()

	currentAddr := startAddr
	for i := 0; i < index; i++ {
		nextAddr := currentAddr.Next()
		if !nextAddr.IsValid() || !startPrefix.Contains(nextAddr) {
			return "", fmt.Errorf("ran out of valid IP addresses after generating %d prefixes, starting from %s", i+1, startAddr)
		}
		currentAddr = nextAddr
	}

	// Assume generated prefixes are host routes
	var generatedPrefixLen int
	if currentAddr.Is6() {
		generatedPrefixLen = 128
	} else if currentAddr.Is4() {
		generatedPrefixLen = 32
	} else {
		return "", fmt.Errorf("invalid start address type: %s", startAddr)
	}

	prefix := netip.PrefixFrom(currentAddr, generatedPrefixLen)
	return prefix.String(), nil
}

// populateNextHops generates NextHop entries and returns the updated slice.
func populateNextHops(cfg *ScaleProfileConfig) ([]fluent.GRIBIEntry, error) {
	totalNextHops := cfg.NumNexthopPerNHG * cfg.NumNexthopGroup
	entries := []fluent.GRIBIEntry{}
	if totalNextHops <= 0 {
		return entries, nil
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

	for i := 1; i <= totalNextHops; i++ {
		mplsLabel := cfg.MPLSLabelStart
		if !cfg.UseSameMPLSLabel {
			mplsLabel = cfg.MPLSLabelStart + uint64(i-1)
		}

		nhEntry := fluent.NextHopEntry().
			WithNetworkInstance(cfg.NetworkInstanceName).
			WithIndex(uint64(i)).
			WithIPAddress(cfg.EgressATEIPv6).
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
	}
	return entries, nil
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

	for i := 0; i < cfg.NumPrefixes; i++ {
		prefixStr, err := GeneratePrefix(cfg.PrefixStart, i) // Use the helper function
		if err != nil {
			return nil, fmt.Errorf("failed to generate prefix at index %d: %w", i, err)
		}

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
	if cfg.EgressATEIPv6 == "" {
		return nil, errors.New("EgressATEIPv6 must be provided in ScaleProfileConfig")
	}
	if _, err := netip.ParseAddr(cfg.EgressATEIPv6); err != nil {
		return nil, fmt.Errorf("invalid EgressATEIPv6 %q: %w", cfg.EgressATEIPv6, err)
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
