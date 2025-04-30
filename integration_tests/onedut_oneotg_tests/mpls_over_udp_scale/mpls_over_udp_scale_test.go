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
	"testing"
	"time"

	"github.com/open-traffic-generator/snappi/gosnappi"
	gribipb "github.com/openconfig/gribi/v1/proto/service"
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/otg"
	"github.com/openconfig/ygnmi/ygnmi"

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
		desc   string
		config *scaleutil.ScaleProfileConfig
		// TODO: Add fields for validation later (e.g., expected results)
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

			// 5. Validate Traffic Flow
			startAddrV6, err := scaleutil.GetFirstAddrFromPrefix(tc.config.PrefixStart)
			if err != nil {
				t.Fatalf("Failed to get start address from prefix %q", tc.config.PrefixStart)
			}
			if loss := testTrafficv6(t, otg, tc.config, atePort1, atePort2, startAddrV6, 5*time.Second); loss > 1 {
				t.Errorf("Loss: got %g, want <= 1", loss)
			}
		})
	}
}
