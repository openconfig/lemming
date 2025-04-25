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
	DstIP      string
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

func generateRandomNHInd(n, nhRange int, r *rand.Rand) []uint64 {
	indSlice := []uint64{}
	for i := 0; i < n; i++ {
		randNum := uint64(r.Intn(int(nhRange))) + 1
		indSlice = append(indSlice, randNum)
	}
	return indSlice
}

// populateNextHopGroups generates NextHopGroup entries, assigning NHs using bootstrapping
// and avoiding duplicate NH combinations across NHGs.
func populateNextHopGroups(cfg *ScaleProfileConfig) ([]fluent.GRIBIEntry, error) {
	totalNHsAvailable := cfg.NumNexthopPerNHG * cfg.NumNexthopGroup

	usedCombinations := make(map[string]bool)
	maxRetries := max(cfg.NumNexthopPerNHG, 20)

	entries := []fluent.GRIBIEntry{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for nhgIdx := 1; nhgIdx <= cfg.NumNexthopGroup; nhgIdx++ {
		var key string
		var indSlice []uint64
		foundUnique := false

		for retry := 0; retry < maxRetries; retry++ {
			indSlice = generateRandomNHInd(cfg.NumNexthopPerNHG, totalNHsAvailable, r)

			key = combinationKey(indSlice)
			if !usedCombinations[key] {
				usedCombinations[key] = true
				foundUnique = true
				break // Found a unique combination
			}
		}

		if !foundUnique {
			return nil, fmt.Errorf("failed to generate a unique NH combination for NHG ID %d after %d retries", nhgIdx, maxRetries)
		}

		nhgEntry := fluent.NextHopGroupEntry().
			WithNetworkInstance(cfg.NetworkInstanceName).
			WithID(uint64(nhgIdx))
		for _, nhIndex := range indSlice {
			nhgEntry.AddNextHop(nhIndex, 1)
		}

		entries = append(entries, nhgEntry)
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
	if cfg.SrcIP == "" || cfg.DstIP == "" {
		return nil, errors.New("SrcIP and DstIP for encapsulation must be provided")
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
