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

package integration_test

import (
	"context"
	"encoding/binary"
	"net/netip"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/open-traffic-generator/snappi/gosnappi"
	gribipb "github.com/openconfig/gribi/v1/proto/service"
	"github.com/openconfig/gribigo/constants"
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/otg"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	scaleutil "github.com/openconfig/lemming/integration_tests/onedut_oneotg_tests/mpls_over_udp_scale/util"
	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/lemming/internal/binding"
)

const (
	ipv4PrefixLen = 30
	ipv6PrefixLen = 126 // Using /126 for point-to-point links as is common.

	// Scale parameters
	numPrefixes        = 200
	numVlans           = 1000
	numPolicers        = 1000
	numPolicerPolicies = 200
	flowQ              = 200  // flow_q
	flowR              = 1    // flow_r (updates per second)
	schedQ             = 1000 // sched_q
	schedR             = 60   // sched_r (seconds)

	// IP Prefixes (using only IPv6 for now)
	// Use a base range for generating inner destination prefixes. /48 gives plenty of space.
	innerIPv6DstRange = "2001:aa:bb::/48"
	innerIPv6DstLen   = 128 // Scale profiles use /128 destinations

	// Outer encap values
	outerIPv6Src     = "2001:f:a:1::0" // Source IP for the UDP tunnel
	outerIPv6DstA    = "2001:f:c:e::1" // Example destination IP for the UDP tunnel
	outerIPv6DstB    = "2001:f:c:e::2" // Example destination IP for the UDP tunnel
	outerIPv6DstDef  = "2001:1:1:1::0" // Default route target
	outerDstUDPPort  = uint64(5555)
	outerDSCP        = uint64(26)
	outerIPTTL       = uint64(64)
	defaultMPLSLabel = uint64(100) // Default starting label for same-label profiles

	// Default network instance
	defaultNI = fakedevice.DefaultNetworkInstance

	gribiBatchSize = 100 // Number of entries to send per ModifyRequest

	// Traffic Flow Parameters
	trafficFlowName   = "ScaleFlowIPv6"
	trafficFlowRate   = 100
	trafficFlowDur    = 10 * time.Second
	trafficLossTarget = 0.0 // Expect no loss for programmed routes
)

var (
	dutPort1 = attrs.Attributes{
		Desc:    "dutPort1",
		IPv4:    "192.0.2.1",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001:db8:1::1",
		IPv6Len: ipv6PrefixLen,
	}

	atePort1 = attrs.Attributes{
		Name:    "port1",
		MAC:     "02:00:01:01:01:01",
		IPv4:    "192.0.2.2",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001:db8:1::2",
		IPv6Len: ipv6PrefixLen,
	}

	dutPort2 = attrs.Attributes{
		Desc:    "dutPort2",
		IPv4:    "192.0.2.5",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001:db8:2::1",
		IPv6Len: ipv6PrefixLen,
	}

	atePort2 = attrs.Attributes{
		Name:    "port2",
		MAC:     "02:00:02:01:01:01",
		IPv4:    "192.0.2.6",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001:db8:2::2",
		IPv6Len: ipv6PrefixLen,
	}
)

// packetResult stores the expected values for validating captured packets.
type packetResult struct {
	mplsLabel  uint64
	udpSrcPort uint16
	udpDstPort uint16
	ipTTL      uint8 // Expected TTL in the *received* packet (usually DUT configured TTL - 1)
	srcIP      string
	dstIP      string
}

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.KNE(".."))
}

// configureDUT configures basic IP addresses on DUT ports.
func configureDUT(t *testing.T, dut *ondatra.DUTDevice) {
	t.Helper()
	p1 := dut.Port(t, "port1")
	p2 := dut.Port(t, "port2")

	// Basic interface configuration
	gnmi.Replace(t, dut, ocpath.Root().Interface(p1.Name()).Config(), dutPort1.NewOCInterface(p1.Name(), dut))
	gnmi.Replace(t, dut, ocpath.Root().Interface(p2.Name()).Config(), dutPort2.NewOCInterface(p2.Name(), dut))

	// Wait for interfaces to be up
	gnmi.Await(t, dut, ocpath.Root().Interface(p1.Name()).OperStatus().State(), 1*time.Minute, oc.Interface_OperStatus_UP)
	gnmi.Await(t, dut, ocpath.Root().Interface(p2.Name()).OperStatus().State(), 1*time.Minute, oc.Interface_OperStatus_UP)

	// Wait for IPv6 addresses to be assigned
	gnmi.Await(t, dut, ocpath.Root().Interface(p1.Name()).Subinterface(0).Ipv6().Address(dutPort1.IPv6).Ip().State(), 1*time.Minute, dutPort1.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(p2.Name()).Subinterface(0).Ipv6().Address(dutPort2.IPv6).Ip().State(), 1*time.Minute, dutPort2.IPv6)

	// TODO: Add VLAN configuration if needed based on specific scale profiles.
	// TODO: Add QoS configuration if needed.
}

// configureOTG configures basic IP addresses on ATE ports.
func configureOTG(t *testing.T, ate *ondatra.ATEDevice) gosnappi.Config {
	t.Helper()
	top := gosnappi.NewConfig()

	p1 := ate.Port(t, "port1")
	p2 := ate.Port(t, "port2")

	// Configure ATE interfaces
	atePort1.AddToOTG(top, p1, &dutPort1)
	atePort2.AddToOTG(top, p2, &dutPort2)

	return top
}

// WaitForARP waits for ARP/ND to resolve on all OTG interfaces for a given ipType ("IPv4" or "IPv6").
func WaitForARP(t *testing.T, otg *otg.OTG, c gosnappi.Config, ipType string) {
	t.Helper()
	intfs := []string{}
	for _, d := range c.Devices().Items() {
		Eth := d.Ethernets().Items()[0]
		intfs = append(intfs, Eth.Name())
	}

	for _, intf := range intfs {
		switch ipType {
		case "IPv4":
			statePath := gnmi.OTG().Interface(intf).Ipv4NeighborAny().LinkLayerAddress().State()
			_, ok := gnmi.WatchAll(t, otg, statePath, 2*time.Minute, func(val *ygnmi.Value[string]) bool {
				return val.IsPresent()
			}).Await(t)
			if !ok {
				t.Fatalf("ND entries not populated for interface %s on OTG.", intf)
			}
		case "IPv6":
			statePath := gnmi.OTG().Interface(intf).Ipv6NeighborAny().LinkLayerAddress().State()
			_, ok := gnmi.WatchAll(t, otg, statePath, 2*time.Minute, func(val *ygnmi.Value[string]) bool {
				return val.IsPresent()
			}).Await(t)
			if !ok {
				t.Fatalf("ND entries not populated for interface %s on OTG.", intf)
			}
		default:
			t.Fatalf("Invalid IP type specified: %s", ipType)
		}
		t.Logf("ARP/ND resolved for OTG interface %s (%s)", intf, ipType)
	}
}

// awaitTimeout calls a fluent client Await, adding a timeout to the context.
func awaitTimeout(ctx context.Context, c *fluent.GRIBIClient, t testing.TB, timeout time.Duration) error {
	t.Helper()
	subctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Await(subctx, t)
}

// testTrafficv6 generates traffic flow from source network to
// destination network via srcEndPoint to dstEndPoint and checks for
// packet loss and returns loss percentage as float.
func testTrafficv6(t *testing.T, otg *otg.OTG, cfg *scaleutil.ScaleProfileConfig, srcEndPoint, dstEndPoint attrs.Attributes, startAddress string, dur time.Duration) float32 {
	WaitForARP(t, otg, otg.FetchConfig(t), "IPv6")
	top := otg.FetchConfig(t)
	top.Flows().Clear().Items()
	flowipv6 := top.Flows().Add().SetName(trafficFlowName)
	flowipv6.Metrics().SetEnable(true)
	flowipv6.TxRx().Device().
		SetTxNames([]string{srcEndPoint.Name + ".IPv6"}).
		SetRxNames([]string{dstEndPoint.Name + ".IPv6"})
	flowipv6.Duration().Continuous()
	flowipv6.Rate().SetPps(uint64(trafficFlowRate))
	flowipv6.Packet().Add().Ethernet()
	v6 := flowipv6.Packet().Add().Ipv6()
	v6.Src().SetValue(srcEndPoint.IPv6)
	v6.Dst().Increment().SetStart(startAddress).SetCount(uint32(cfg.NumPrefixes))
	otg.PushConfig(t, top)

	otg.StartTraffic(t)
	time.Sleep(dur)
	t.Logf("Stop traffic")
	otg.StopTraffic(t)

	time.Sleep(5 * time.Second)

	txPkts := gnmi.Get(t, otg, gnmi.OTG().Flow(trafficFlowName).Counters().OutPkts().State())
	rxPkts := gnmi.Get(t, otg, gnmi.OTG().Flow(trafficFlowName).Counters().InPkts().State())
	lossPct := (txPkts - rxPkts) * 100 / txPkts
	return float32(lossPct)
}

// validateAFTState checks the AFT state for a sample of NHs, NHGs, and Prefixes.
// If expectPresent is true, it validates that entries exist and match the config.
// If expectPresent is false, it validates that entries do NOT exist.
func validateAFTState(t *testing.T, dut *ondatra.DUTDevice, cfg *scaleutil.ScaleProfileConfig, expectPresent bool) {
	t.Helper()
	expectation := "present"
	if !expectPresent {
		expectation = "absent"
	}
	t.Logf("Validating AFT state for profile: %s (expecting entries to be %s)", cfg.NetworkInstanceName, expectation)

	// --- Validate a sample of NextHops ---
	nhIndicesToCheck := []uint64{1}
	totalNHs := uint64(cfg.NumNexthopGroup * cfg.NumNexthopPerNHG)
	if totalNHs > 1 {
		nhIndicesToCheck = append(nhIndicesToCheck, totalNHs)
	}

	for _, nhIndex := range nhIndicesToCheck {
		t.Logf("Checking NH index: %d (expect %s)", nhIndex, expectation)
		nhPath := ocpath.Root().NetworkInstance(cfg.NetworkInstanceName).Afts().NextHop(nhIndex)
		nhStateVal := gnmi.Lookup(t, dut, nhPath.State())
		nhState, found := nhStateVal.Val()

		if expectPresent {
			if !found {
				t.Errorf("AFT NextHop %d not found in NI %s, but expected it to be present", nhIndex, cfg.NetworkInstanceName)
				continue
			}

			// Construct expected Encap Headers based on config
			expectedMPLSLabel := cfg.MPLSLabelStart
			if !cfg.UseSameMPLSLabel {
				expectedMPLSLabel = cfg.MPLSLabelStart + nhIndex - 1
			}
			expectedDstIP := cfg.DstIPStart // Assumes NumDstIP = 1 for Profile A

			wantEncapHeaders := map[uint8]*oc.NetworkInstance_Afts_NextHop_EncapHeader{
				1: { // MPLS Header
					Index: ygot.Uint8(1),
					Type:  oc.AftTypes_EncapsulationHeaderType_MPLS,
					Mpls: &oc.NetworkInstance_Afts_NextHop_EncapHeader_Mpls{
						MplsLabelStack: []oc.NetworkInstance_Afts_NextHop_EncapHeader_Mpls_MplsLabelStack_Union{
							oc.UnionUint32(expectedMPLSLabel),
						},
					},
				},
				2: { // UDPv6 Header
					Index: ygot.Uint8(2),
					Type:  oc.AftTypes_EncapsulationHeaderType_UDPV6,
					UdpV6: &oc.NetworkInstance_Afts_NextHop_EncapHeader_UdpV6{
						Dscp:       ygot.Uint8(uint8(cfg.DSCP)),
						DstIp:      ygot.String(expectedDstIP),
						DstUdpPort: ygot.Uint16(uint16(cfg.UDPDstPort)),
						IpTtl:      ygot.Uint8(uint8(cfg.IPTTL)),
						SrcIp:      ygot.String(cfg.SrcIP),
						SrcUdpPort: ygot.Uint16(uint16(cfg.UDPSrcPort)),
					},
				},
			}

			// Compare actual encap headers with expected
			gotEncapHeaders := make(map[uint8]*oc.NetworkInstance_Afts_NextHop_EncapHeader)
			for _, eh := range nhState.EncapHeader {
				gotEncapHeaders[eh.GetIndex()] = eh
			}

			if diff := cmp.Diff(wantEncapHeaders, gotEncapHeaders); diff != "" {
				t.Errorf("NH index %d EncapHeader mismatch (-want +got):\n%s", nhIndex, diff)
			}

			// Check IP address associated with NH
			if nhState.GetIpAddress() != cfg.EgressATEIPv6 {
				t.Errorf("NH index %d IP address mismatch: got %q, want %q", nhIndex, nhState.GetIpAddress(), cfg.EgressATEIPv6)
			}
		} else { // expectPresent == false
			if found {
				t.Errorf("AFT NextHop %d FOUND in NI %s, but expected it to be absent", nhIndex, cfg.NetworkInstanceName)
			}
		}
	}

	// --- Validate a sample of NextHopGroups ---
	nhgIndicesToCheck := []uint64{1}
	if uint64(cfg.NumNexthopGroup) > 1 {
		nhgIndicesToCheck = append(nhgIndicesToCheck, uint64(cfg.NumNexthopGroup))
	}

	for _, nhgIndex := range nhgIndicesToCheck {
		t.Logf("Checking NHG index: %d (expect %s)", nhgIndex, expectation)
		nhgPath := ocpath.Root().NetworkInstance(cfg.NetworkInstanceName).Afts().NextHopGroup(nhgIndex)
		nhgStateVal := gnmi.Lookup(t, dut, nhgPath.State())
		nhgState, found := nhgStateVal.Val()

		if expectPresent {
			if !found {
				t.Errorf("AFT NextHopGroup %d not found in NI %s, but expected it to be present", nhgIndex, cfg.NetworkInstanceName)
				continue
			}

			// Validate number of NHs in the group
			if len(nhgState.NextHop) != cfg.NumNexthopPerNHG {
				t.Errorf("NHG index %d has %d NHs, want %d", nhgIndex, len(nhgState.NextHop), cfg.NumNexthopPerNHG)
				continue
			}

			// Validate the NH index points correctly (simplified check for k=1)
			if cfg.NumNexthopPerNHG == 1 {
				var containedNHIndex uint64
				for nhIdxInGroup := range nhgState.NextHop {
					containedNHIndex = nhIdxInGroup
					break
				}
				maxExpectedNHIndex := uint64(cfg.NumNexthopGroup * cfg.NumNexthopPerNHG)
				if containedNHIndex < 1 || containedNHIndex > maxExpectedNHIndex {
					t.Errorf("NHG index %d contains invalid NH index %d (expected range [1, %d])", nhgIndex, containedNHIndex, maxExpectedNHIndex)
				} else {
					t.Logf("NHG index %d correctly contains NH index %d (within range [1, %d])", nhgIndex, containedNHIndex, maxExpectedNHIndex)
				}
			} else {
				t.Logf("Skipping specific NH index validation for NHG %d because NumNexthopPerNHG > 1", nhgIndex)
			}
		} else { // expectPresent == false
			if found {
				t.Errorf("AFT NextHopGroup %d FOUND in NI %s, but expected it to be absent", nhgIndex, cfg.NetworkInstanceName)
			}
		}
	}

	// --- Validate a sample of Prefixes ---
	firstPrefixStr, err := scaleutil.GeneratePrefix(cfg.PrefixStart, 0)
	if err != nil {
		t.Errorf("Failed to generate first prefix for validation: %v", err)
		return // Can't proceed if prefix generation fails
	}
	prefixesToCheck := []string{firstPrefixStr}

	if cfg.NumPrefixes > 1 {
		lastPrefixStr, err := scaleutil.GeneratePrefix(cfg.PrefixStart, cfg.NumPrefixes-1)
		if err != nil {
			t.Errorf("Failed to generate last prefix for validation: %v", err)
		} else {
			prefixesToCheck = append(prefixesToCheck, lastPrefixStr)
		}
	}

	for i, prefixStr := range prefixesToCheck {
		t.Logf("Checking Prefix: %s (expect %s)", prefixStr, expectation)
		var prefixStateVal *ygnmi.Value[*oc.NetworkInstance_Afts_Ipv6Entry] // Adjust type for IPv4 if needed
		var found bool

		if cfg.AddrFamily == "ipv6" {
			ipv6Path := ocpath.Root().NetworkInstance(cfg.NetworkInstanceName).Afts().Ipv6Entry(prefixStr)
			prefixStateVal = gnmi.Lookup(t, dut, ipv6Path.State())
			_, found = prefixStateVal.Val()
		} else {
			t.Errorf("IPv4 prefix validation not implemented yet.")
			continue
		}

		if expectPresent {
			if !found {
				t.Errorf("AFT Prefix %s not found in NI %s, but expected it to be present", prefixStr, cfg.NetworkInstanceName)
				continue
			}

			prefixState, _ := prefixStateVal.Val() // We know it's present here
			if prefixState == nil {
				t.Errorf("Internal error: prefixState is nil for prefix %s despite being found", prefixStr)
				continue
			}

			// Validate the NHG it points to (assuming round-robin assignment in util)
			expectedNHGIndex := uint64(i*(cfg.NumPrefixes-1))%uint64(cfg.NumNexthopGroup) + 1
			if i == 0 {
				expectedNHGIndex = 1
			} else if cfg.NumPrefixes > 0 {
				expectedNHGIndex = uint64((cfg.NumPrefixes-1)%cfg.NumNexthopGroup + 1)
			}

			if prefixState.GetNextHopGroup() != expectedNHGIndex {
				t.Errorf("Prefix %s points to NHG %d, want %d", prefixStr, prefixState.GetNextHopGroup(), expectedNHGIndex)
			}
		} else { // expectPresent == false
			if found {
				t.Errorf("AFT Prefix %s FOUND in NI %s, but expected it to be absent", prefixStr, cfg.NetworkInstanceName)
			}
		}
	}
	t.Logf("AFT state validation sample completed (expected %s).", expectation)
}

// enableCapture enables packet capture on specified list of ports on OTG
func enableCapture(t *testing.T, otg *otg.OTG, topo gosnappi.Config, otgPortName string) {
	t.Log("Enabling capture on ", otgPortName)
	topo.Captures().Add().SetName(otgPortName).SetPortNames([]string{otgPortName}).SetFormat(gosnappi.CaptureFormat.PCAP)
	pb, _ := topo.Marshal().ToProto()
	t.Log(pb.GetCaptures())
	otg.PushConfig(t, topo)
}

// startCapture starts the capture on the otg ports
func startCapture(t *testing.T, ate *ondatra.ATEDevice) {
	otg := ate.OTG()
	cs := gosnappi.NewControlState()
	cs.Port().Capture().SetState(gosnappi.StatePortCaptureState.START)
	otg.SetControlState(t, cs)
}

// stopCapture starts the capture on the otg ports
func stopCapture(t *testing.T, ate *ondatra.ATEDevice) {
	otg := ate.OTG()
	cs := gosnappi.NewControlState()
	cs.Port().Capture().SetState(gosnappi.StatePortCaptureState.STOP)
	otg.SetControlState(t, cs)
}

// validatePacketCapture reads capture files and checks the encapped packet for desired protocol, dscp and ttl
func validatePacketCapture(t *testing.T, ate *ondatra.ATEDevice, otgPortName string, pr *packetResult) {
	t.Helper()
	otg := ate.OTG()
	t.Logf("Validating packet capture from port %s...", otgPortName)

	// 1. Get captured bytes
	captureBytes := otg.GetCapture(t, gosnappi.NewCaptureRequest().SetPortName(otgPortName))

	// 2. Write to temporary pcap file
	f, err := os.CreateTemp("", "mpls-scale-capture-*.pcap")
	if err != nil {
		t.Fatalf("ERROR: Could not create temporary pcap file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(captureBytes); err != nil {
		t.Fatalf("ERROR: Could not write packetBytes to pcap file %q: %v", f.Name(), err)
	}
	f.Close()
	t.Logf("Wrote capture to %s", f.Name())

	// 3. Open pcap file
	handle, err := pcap.OpenOffline(f.Name())
	if err != nil {
		t.Fatalf("ERROR: pcap.OpenOffline(%s) failed: %v", f.Name(), err)
	}
	defer handle.Close()

	// 4. Process packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetCount := 0
	checkedPacketCount := 0
	mismatchedPacketFound := false

	// Parse expected IPs once
	wantSrcIP, err := netip.ParseAddr(pr.srcIP)
	if err != nil {
		t.Fatalf("ERROR: Could not parse expected source IP %q: %v", pr.srcIP, err)
	}
	wantDstIP, err := netip.ParseAddr(pr.dstIP)
	if err != nil {
		t.Fatalf("ERROR: Could not parse expected destination IP %q: %v", pr.dstIP, err)
	}

	for packet := range packetSource.Packets() {
		packetCount++
		ipv6Layer := packet.Layer(layers.LayerTypeIPv6)
		udpLayer := packet.Layer(layers.LayerTypeUDP)

		// Skip packets that don't have the outer IPv6/UDP structure
		if ipv6Layer == nil || udpLayer == nil {
			continue
		}
		checkedPacketCount++
		v6Packet, _ := ipv6Layer.(*layers.IPv6)
		udpPacket, _ := udpLayer.(*layers.UDP)

		// --- Perform Validations ---
		currentPacketValid := true

		// Validate Outer IPv6 Header
		gotSrcIP, ok := netip.AddrFromSlice(v6Packet.SrcIP)
		if !ok || gotSrcIP != wantSrcIP {
			t.Errorf("Packet %d: Outer IPv6 SrcIP mismatch: got %s, want %s", packetCount, v6Packet.SrcIP, pr.srcIP)
			currentPacketValid = false
		}
		gotDstIP, ok := netip.AddrFromSlice(v6Packet.DstIP)
		if !ok || gotDstIP != wantDstIP {
			t.Errorf("Packet %d: Outer IPv6 DstIP mismatch: got %s, want %s", packetCount, v6Packet.DstIP, pr.dstIP)
			currentPacketValid = false
		}

		// Check HopLimit (TTL). It should be decremented by the DUT.
		if v6Packet.HopLimit != pr.ipTTL {
			t.Errorf("Packet %d: Outer IPv6 HopLimit (TTL) mismatch: got %d, want %d", packetCount, v6Packet.HopLimit, pr.ipTTL)
			currentPacketValid = false
		}
		// TODO: Add DSCP validation if needed: v6Packet.TrafficClass

		// Validate Outer UDP Header
		if udpPacket.SrcPort != layers.UDPPort(pr.udpSrcPort) {
			t.Errorf("Packet %d: Outer UDP SrcPort mismatch: got %d, want %d", packetCount, udpPacket.SrcPort, pr.udpSrcPort)
			currentPacketValid = false
		}
		if udpPacket.DstPort != layers.UDPPort(pr.udpDstPort) {
			t.Errorf("Packet %d: Outer UDP DstPort mismatch: got %d, want %d", packetCount, udpPacket.DstPort, pr.udpDstPort)
			currentPacketValid = false
		}

		// Validate UDP Payload (Extract and check MPLS Label only)
		payload := udpPacket.LayerPayload()
		if len(payload) < 4 {
			t.Errorf("Packet %d: UDP Payload too short for MPLS Header (len %d)", packetCount, len(payload))
			currentPacketValid = false
		} else {
			// Extract label: Label (20 bits) = (headerValue >> 12) & 0xFFFFF
			headerValue := binary.BigEndian.Uint32(payload[:4])
			gotLabel := uint64((headerValue >> 12) & 0xFFFFF)
			if gotLabel != pr.mplsLabel {
				t.Errorf("Packet %d: MPLS Label mismatch: got %d, want %d (Full header: %s)",
					packetCount, gotLabel, pr.mplsLabel, scaleutil.FormatMPLSHeader(payload))
				currentPacketValid = false
			}
		}

		if !currentPacketValid {
			mismatchedPacketFound = true
		}
	}

	if packetCount == 0 {
		t.Errorf("No packets found in capture file %s", f.Name())
	} else if checkedPacketCount == 0 {
		t.Errorf("Found %d packets, but none had the expected outer IPv6/UDP structure.", packetCount)
	} else if mismatchedPacketFound {
		t.Errorf("Found %d packets with expected outer IPv6/UDP structure, but at least one had encapsulation mismatches (see errors above).", checkedPacketCount)
	} else {
		t.Logf("Successfully validated %d packets with expected encapsulation.", checkedPacketCount)
	}
}

// testCounters test packet counters and should be called after testTraffic
func testCounters(t *testing.T, dut *ondatra.DUTDevice, wantTxPkts, wantRxPkts uint64) {
	got := gnmi.Get(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Counters().InPkts().State())
	t.Logf("DUT port 1 in-pkts: %d", got)
	if got < wantTxPkts {
		t.Errorf("DUT got less packets (%d) than OTG sent (%d)", got, wantTxPkts)
	}

	got = gnmi.Get(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Counters().OutPkts().State())
	t.Logf("DUT port 2 out-pkts: %d", got)
	if got < wantRxPkts {
		t.Errorf("DUT got sent less packets (%d) than OTG received (%d)", got, wantRxPkts)
	}
}

// TestMPLSOverUDPScale sets up the basic DUT and ATE environment and runs scale profile tests.
func TestMPLSOverUDPScale(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	configureDUT(t, dut)

	ate := ondatra.ATE(t, "ate")
	otg := ate.OTG()
	otgConfig := configureOTG(t, ate)
	t.Log("Pushing ATE config and starting protocols...")
	otg.PushConfig(t, otgConfig)
	otg.StartProtocols(t)

	t.Log("Waiting for IPv6 ND resolution...")
	WaitForARP(t, otg, otgConfig, "IPv6")

	t.Log("Environment setup complete.")

	// Define scale profile test cases
	tests := []struct {
		desc        string
		config      *scaleutil.ScaleProfileConfig
		capturePort string
	}{
		{
			desc: "Scale Profile A - 1 NI, 20k NHG, 20k Prefixes, 1 NH/NHG, Same MPLS Label",
			config: &scaleutil.ScaleProfileConfig{
				AddrFamily:          "ipv6",
				NetworkInstanceName: defaultNI,
				NumPrefixes:         numPrefixes,
				NumNexthopGroup:     numPrefixes,
				NumNexthopPerNHG:    1,
				PrefixStart:         innerIPv6DstRange,
				UseSameMPLSLabel:    true,
				MPLSLabelStart:      defaultMPLSLabel,
				UDPSrcPort:          outerDstUDPPort,
				UDPDstPort:          outerDstUDPPort,
				SrcIP:               outerIPv6Src,
				DstIPStart:          outerIPv6DstA,
				NumDstIP:            1,
				DSCP:                outerDSCP,
				IPTTL:               outerIPTTL,
				EgressATEIPv6:       atePort2.IPv6,
			},
			capturePort: atePort2.Name,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			ctx := context.Background()

			// 1. Generate gRIBI entries
			entries, err := scaleutil.GenerateScaleProfileEntries(ctx, tc.config)
			if err != nil {
				t.Fatalf("Failed to generate gRIBI entries: %v", err)
			}
			t.Logf("Generated %d total gRIBI entries.", len(entries))
			expectedAddCount := tc.config.NumNexthopGroup*tc.config.NumNexthopPerNHG + tc.config.NumNexthopGroup + tc.config.NumPrefixes
			if len(entries) != expectedAddCount {
				t.Fatalf("Generated entry count mismatch: got %d, want %d", len(entries), expectedAddCount)
			}

			// 2. Establish gRIBI connection
			gribic := dut.RawAPIs().GRIBI(t)
			c := fluent.NewClient()
			c.Connection().WithStub(gribic).
				WithRedundancyMode(fluent.ElectedPrimaryClient).
				WithPersistence().
				WithFIBACK().
				WithInitialElectionID(1, 0)

			c.Start(ctx, t)
			defer c.Stop(t)
			c.StartSending(ctx, t)
			if err := awaitTimeout(ctx, c, t, 2*time.Minute); err != nil {
				t.Fatalf("Await got error during session negotiation: %v", err)
			}

			// 3. Add entries via gRIBI.Modify
			t.Logf("Sending %d ADD operations in batches of %d...", len(entries), gribiBatchSize)
			totalSent := 0
			for i := 0; i < len(entries); i += gribiBatchSize {
				end := i + gribiBatchSize
				if end > len(entries) {
					end = len(entries)
				}
				batch := entries[i:end]
				totalSent += len(batch)
				t.Logf("Sending batch %d/%d (%d entries, total sent: %d)", (i/gribiBatchSize)+1, (len(entries)+gribiBatchSize-1)/gribiBatchSize, len(batch), totalSent)
				c.Modify().AddEntry(t, batch...)
				batchTimeout := 3 * time.Minute
				if err := awaitTimeout(ctx, c, t, batchTimeout); err != nil {
					t.Fatalf("Await got error for ADD batch %d: %v", (i/gribiBatchSize)+1, err)
				}
			}
			t.Logf("Finished sending all %d entries.", totalSent)

			// 4. Validate length of the results.
			results := c.Results(t)
			gotInstalledCount := 0
			for _, res := range results {
				if res.ProgrammingResult == gribipb.AFTResult_FIB_PROGRAMMED {
					gotInstalledCount++
				}
			}
			t.Logf("Received total of %d results, found %d FIB_PROGRAMMED results.", len(results), gotInstalledCount)
			if gotInstalledCount != expectedAddCount {
				t.Errorf("Got %d results, want %d", gotInstalledCount, expectedAddCount)
			}

			// 5. AFT validations
			validateAFTState(t, dut, tc.config, true)

			// 6. Validate Traffic Flow
			enableCapture(t, otg, otgConfig, tc.capturePort)
			time.Sleep(1 * time.Second)
			startCapture(t, ate)
			time.Sleep(1 * time.Second)
			startAddrV6, err := scaleutil.GetFirstAddrFromPrefix(tc.config.PrefixStart)
			if err != nil {
				t.Fatalf("Failed to get start address from prefix %q", tc.config.PrefixStart)
			}
			if loss := testTrafficv6(t, otg, tc.config, atePort1, atePort2, startAddrV6, 5*time.Second); loss > 1 {
				t.Errorf("Loss: got %g, want <= 1", loss)
			}
			time.Sleep(10 * time.Second)
			stopCapture(t, ate)

			var txPkts, rxPkts uint64
			// counters are not erased, so have to accumulate the packets from previous subtests.
			txPkts += gnmi.Get(t, otg, gnmi.OTG().Flow(trafficFlowName).Counters().OutPkts().State())
			rxPkts += gnmi.Get(t, otg, gnmi.OTG().Flow(trafficFlowName).Counters().InPkts().State())
			testCounters(t, dut, txPkts, rxPkts)

			validatePacketCapture(t, ate, tc.capturePort,
				&packetResult{
					mplsLabel:  tc.config.MPLSLabelStart, // For Profile A, label is the same
					udpSrcPort: uint16(tc.config.UDPSrcPort),
					udpDstPort: uint16(tc.config.UDPDstPort),
					ipTTL:      uint8(tc.config.IPTTL - 1), // Expect TTL decremented by DUT
					srcIP:      tc.config.SrcIP,
					dstIP:      tc.config.DstIPStart, // For Profile A, only one DstIP
				})

			// 7. Delete Entries
			t.Logf("Sending DELETE operations for %d entries...", len(entries))
			// Reverse the order for deletion: Prefixes -> NHGs -> NHs
			// The 'entries' slice currently has NHs, then NHGs, then Prefixes.
			prefixStartIdx := tc.config.NumNexthopGroup*tc.config.NumNexthopPerNHG + tc.config.NumNexthopGroup
			nhgStartIdx := tc.config.NumNexthopGroup * tc.config.NumNexthopPerNHG

			deleteEntries := []fluent.GRIBIEntry{}
			deleteEntries = append(deleteEntries, entries[prefixStartIdx:]...)            // Prefixes
			deleteEntries = append(deleteEntries, entries[nhgStartIdx:prefixStartIdx]...) // NHGs
			deleteEntries = append(deleteEntries, entries[:nhgStartIdx]...)               // NHs

			totalSent = 0
			for i := 0; i < len(deleteEntries); i += gribiBatchSize {
				end := i + gribiBatchSize
				if end > len(deleteEntries) {
					end = len(deleteEntries)
				}
				batch := deleteEntries[i:end]
				totalSent += len(batch)
				t.Logf("Sending DELETE batch %d/%d (%d entries, total sent: %d)", (i/gribiBatchSize)+1, (len(deleteEntries)+gribiBatchSize-1)/gribiBatchSize, len(batch), totalSent)
				c.Modify().DeleteEntry(t, batch...) // Use DeleteEntry with the correct entries
				batchTimeout := 3 * time.Minute
				if err := awaitTimeout(ctx, c, t, batchTimeout); err != nil {
					t.Errorf("Await got error for DELETE batch %d: %v", (i/gribiBatchSize)+1, err)
				}
			}
			t.Logf("Finished sending all %d DELETE entries.", totalSent)

			// 8. Validate Deletion Results
			delResults := c.Results(t)
			gotDeletedCount := 0
			for _, res := range delResults {
				// Check if it's an AFT result and corresponds to a Delete operation
				if res.Details != nil && res.Details.Type == constants.Delete && res.ProgrammingResult == gribipb.AFTResult_FIB_PROGRAMMED {
					gotDeletedCount++
				}
			}
			if gotDeletedCount != expectedAddCount {
				t.Errorf("Got %d successful delete results, want %d", gotDeletedCount, expectedAddCount)
			}

			// 9. Verify Traffic Loss After Deletion
			t.Log("Verifying 100% loss after deletion...")
			if loss := testTrafficv6(t, otg, tc.config, atePort1, atePort2, startAddrV6, 5*time.Second); loss != 100 {
				t.Errorf("Loss after deletion: got %g, want 100", loss)
			}

			// 10. Flush gRIBI Entries (Optional but good practice)
			t.Log("Flushing all gRIBI entries...")
			// Use FlushAll for simplicity, or specify the network instance
			gribic.Flush(ctx, &gribipb.FlushRequest{
				NetworkInstance: &gribipb.FlushRequest_All{All: &gribipb.Empty{}},
			})

			// 11. Validate AFT State After Flush (Optional)
			// TODO: Call validateAFTState again, modifying it to expect *not found* results.
			validateAFTState(t, dut, tc.config, false)
		})
	}
}
