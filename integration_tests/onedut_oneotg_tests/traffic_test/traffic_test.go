// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration_test

import (
	"context"
	"net/netip"
	"slices"
	"testing"
	"time"

	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/gribigo/chk"
	"github.com/openconfig/gribigo/client"
	"github.com/openconfig/gribigo/constants"
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/gnmi/otg/otgpath"
	"github.com/openconfig/ondatra/otg"
	"github.com/openconfig/ygnmi/ygnmi"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/lemming/internal/binding"

	gribipb "github.com/openconfig/gribi/v1/proto/service"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.KNE(".."))
}

// Settings for configuring the baseline testbed with the test
// topology.
//
// The testbed consists of ate:port1 -> dut:port1
// and dut:port2 -> ate:port2.
//
//   - ate:port1 -> dut:port1 subnet 192.0.2.0/30
//   - ate:port2 -> dut:port2 subnet 192.0.2.4/30
const (
	ipv4PrefixLen = 30
	ipv6PrefixLen = 99
)

var (
	dutPort1 = attrs.Attributes{
		Desc:    "dutPort1",
		IPv4:    "192.0.2.1",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001::aaaa:bbbb:aa",
		IPv6Len: ipv6PrefixLen,
	}

	atePort1 = attrs.Attributes{
		Name:    "port1",
		MAC:     "02:00:01:01:01:01",
		IPv4:    "192.0.2.2",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001::aaaa:bbbb:bb",
		IPv6Len: ipv6PrefixLen,
	}

	dutPort2 = attrs.Attributes{
		Desc:    "dutPort2",
		IPv4:    "192.0.2.5",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001::aaab:bbbb:aa",
		IPv6Len: ipv6PrefixLen,
	}

	atePort2 = attrs.Attributes{
		Name:    "port2",
		MAC:     "02:00:02:01:01:01",
		IPv4:    "192.0.2.6",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001::aaab:bbbb:bb",
		IPv6Len: ipv6PrefixLen,
	}
)

// configureOTG configures port1 and port2 on the ATE.
func configureOTG(t *testing.T, ate *ondatra.ATEDevice) gosnappi.Config {
	top := gosnappi.NewConfig()

	p1 := ate.Port(t, "port1")
	p2 := ate.Port(t, "port2")

	atePort1.AddToOTG(top, p1, &dutPort1)
	atePort2.AddToOTG(top, p2, &dutPort2)

	return top
}

var gatewayMap = map[attrs.Attributes]attrs.Attributes{
	atePort1: dutPort1,
	atePort2: dutPort2,
}

// configureDUT configures port1 and port2 on the DUT.
func configureDUT(t *testing.T, dut *ondatra.DUTDevice) {
	p1 := dut.Port(t, "port1")
	gnmi.Replace(t, dut, ocpath.Root().Interface(p1.Name()).Config(), dutPort1.NewOCInterface(p1.Name(), dut))

	p2 := dut.Port(t, "port2")
	gnmi.Replace(t, dut, ocpath.Root().Interface(p2.Name()).Config(), dutPort2.NewOCInterface(p2.Name(), dut))

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv4().Address(dutPort1.IPv4).Ip().State(), time.Minute, dutPort1.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv4().Address(dutPort2.IPv4).Ip().State(), time.Minute, dutPort2.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv6().Address(dutPort1.IPv6).Ip().State(), time.Minute, dutPort1.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv6().Address(dutPort2.IPv6).Ip().State(), time.Minute, dutPort2.IPv6)
}

func waitOTGARPEntry(t *testing.T) {
	ate := ondatra.ATE(t, "ate")

	val, ok := gnmi.WatchAll(t, ate.OTG(), otgpath.Root().InterfaceAny().Ipv4NeighborAny().LinkLayerAddress().State(), time.Minute, func(v *ygnmi.Value[string]) bool {
		return v.IsPresent()
	}).Await(t)
	if !ok {
		t.Fatal("failed to get neighbor")
	}
	lla, _ := val.Val()
	t.Logf("Neighbor %v", lla)

	val, ok = gnmi.WatchAll(t, ate.OTG(), otgpath.Root().InterfaceAny().Ipv6NeighborAny().LinkLayerAddress().State(), time.Minute, func(v *ygnmi.Value[string]) bool {
		return v.IsPresent()
	}).Await(t)
	if !ok {
		t.Fatal("failed to get neighbor")
	}
	lla, _ = val.Val()
	t.Logf("Neighbor %v", lla)
}

// testTraffic generates traffic flow from source network to
// destination network via srcEndPoint to dstEndPoint and checks for
// packet loss and returns loss percentage as float.
func testTraffic(t *testing.T, otg *otg.OTG, srcEndPoint, dstEndPoint attrs.Attributes, startAddress string, dur time.Duration) float32 {
	waitOTGARPEntry(t)
	top := otg.FetchConfig(t)
	top.Flows().Clear().Items()
	flowipv4 := top.Flows().Add().SetName("Flow")
	flowipv4.Metrics().SetEnable(true)
	flowipv4.TxRx().Device().
		SetTxNames([]string{srcEndPoint.Name + ".IPv4"}).
		SetRxNames([]string{dstEndPoint.Name + ".IPv4"})
	flowipv4.Duration().Continuous()
	flowipv4.Packet().Add().Ethernet()
	v4 := flowipv4.Packet().Add().Ipv4()
	v4.Src().SetValue(srcEndPoint.IPv4)
	v4.Dst().Increment().SetStart(startAddress).SetCount(24)
	otg.PushConfig(t, top)

	otg.StartTraffic(t)
	time.Sleep(dur)
	t.Logf("Stop traffic")
	otg.StopTraffic(t)

	time.Sleep(5 * time.Second)

	txPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().OutPkts().State())
	rxPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().InPkts().State())
	lossPct := (txPkts - rxPkts) * 100 / txPkts
	return float32(lossPct)
}

// testTrafficv6 generates traffic flow from source network to
// destination network via srcEndPoint to dstEndPoint and checks for
// packet loss and returns loss percentage as float.
func testTrafficv6(t *testing.T, otg *otg.OTG, srcEndPoint, dstEndPoint attrs.Attributes, startAddress string, dur time.Duration) float32 {
	waitOTGARPEntry(t)
	top := otg.FetchConfig(t)
	top.Flows().Clear().Items()
	flowipv6 := top.Flows().Add().SetName("Flow2")
	flowipv6.Metrics().SetEnable(true)
	flowipv6.TxRx().Device().
		SetTxNames([]string{srcEndPoint.Name + ".IPv6"}).
		SetRxNames([]string{dstEndPoint.Name + ".IPv6"})
	flowipv6.Duration().Continuous()
	flowipv6.Packet().Add().Ethernet()
	v6 := flowipv6.Packet().Add().Ipv6()
	v6.Src().SetValue(srcEndPoint.IPv6)
	v6.Dst().Increment().SetStart(startAddress).SetCount(24)
	otg.PushConfig(t, top)

	otg.StartTraffic(t)
	time.Sleep(dur)
	t.Logf("Stop traffic")
	otg.StopTraffic(t)

	time.Sleep(5 * time.Second)

	txPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow2").Counters().OutPkts().State())
	rxPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow2").Counters().InPkts().State())
	lossPct := (txPkts - rxPkts) * 100 / txPkts
	return float32(lossPct)
}

// awaitTimeout calls a fluent client Await, adding a timeout to the context.
func awaitTimeout(ctx context.Context, c *fluent.GRIBIClient, t testing.TB, timeout time.Duration) error {
	subctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Await(subctx, t)
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

type testCase struct {
	desc                    string
	entries                 []fluent.GRIBIEntry
	wantAddOperationResults []*client.OpResult
	wantDelOperationResults []*client.OpResult
	startAddress            string
	v6Traffic               bool
}

type RouteType int

const (
	IPv4 RouteType = iota
	IPv6
	IPv4MappedIPv6
)

func newTestCase(t *testing.T, desc string, startAddress string, routeType RouteType, recursive bool) testCase {
	const (
		// IPv4
		ateDstNetCIDRv4     = "198.51.100.0/24"
		ateIndirectNHv4     = "203.0.113.1"
		ateIndirectNHCIDRv4 = ateIndirectNHv4 + "/32"
		// IPv6
		ateDstNetCIDRv6     = "2003::/48"
		ateIndirectNHv6     = "2002::"
		ateIndirectNHCIDRv6 = ateIndirectNHv6 + "/48"
		// Common attributes
		nhIndex   = 1
		nhgIndex  = 42
		nhIndex2  = 2
		nhgIndex2 = 52
		nhIndex3  = 3
		nhgIndex3 = 62
	)

	var (
		ateIndirectNH     = ateIndirectNHv4
		ateDstNetCIDR     = ateDstNetCIDRv4
		ateIndirectNHCIDR = ateIndirectNHCIDRv4
		destIP            = atePort2.IPv4
	)
	switch routeType {
	case IPv4:
	case IPv6:
		destIP = atePort2.IPv6
		ateIndirectNH = ateIndirectNHv6
		ateDstNetCIDR = ateDstNetCIDRv6
		ateIndirectNHCIDR = ateIndirectNHCIDRv6
	case IPv4MappedIPv6:
		if recursive {
			ateIndirectNHCIDR = ateIndirectNHCIDR
			destIP = atePort2.IPv4

			ateDstNetCIDR = ateDstNetCIDRv6
			mappedAddr, err := netip.ParseAddr("::ffff:" + ateIndirectNH)
			if err != nil {
				t.Fatal(err)
			}
			ateIndirectNH = mappedAddr.StringExpanded()
		} else {
			ateDstNetCIDR = ateDstNetCIDRv6
			mappedDestIP, err := netip.ParseAddr("::ffff:" + atePort2.IPv4)
			if err != nil {
				t.Fatal(err)
			}
			destIP = mappedDestIP.StringExpanded()
		}
	}

	fluentEntry := func(prefix string, nhgIndex uint64) fluent.GRIBIEntry {
		pfx, err := netip.ParsePrefix(prefix)
		if err != nil {
			t.Fatal(err)
		}
		if pfx.Addr().Is4() {
			return fluent.IPv4Entry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
				WithPrefix(prefix).WithNextHopGroup(nhgIndex)
		} else {
			return fluent.IPv6Entry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
				WithPrefix(prefix).WithNextHopGroup(nhgIndex)
		}
	}

	routeInstallResult := func(prefix string, c constants.OpType) *client.OpResult {
		pfx, err := netip.ParsePrefix(prefix)
		if err != nil {
			t.Fatal(err)
		}
		if pfx.Addr().Is4() {
			return fluent.OperationResult().
				WithIPv4Operation(prefix).
				WithProgrammingResult(fluent.InstalledInFIB).
				WithOperationType(c).
				AsResult()
		} else {
			return fluent.OperationResult().
				WithIPv6Operation(prefix).
				WithProgrammingResult(fluent.InstalledInFIB).
				WithOperationType(c).
				AsResult()
		}
	}

	if !recursive {
		return testCase{
			desc: desc,
			entries: []fluent.GRIBIEntry{
				fluent.NextHopEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(nhIndex).WithIPAddress(destIP),
				fluent.NextHopGroupEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithID(nhgIndex).AddNextHop(nhIndex, 1),
				fluentEntry(ateDstNetCIDR, nhgIndex),
			},
			wantAddOperationResults: []*client.OpResult{
				fluent.OperationResult().
					WithNextHopOperation(nhIndex).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Add).
					AsResult(),
				fluent.OperationResult().
					WithNextHopGroupOperation(nhgIndex).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Add).
					AsResult(),
				routeInstallResult(ateDstNetCIDR, constants.Add),
			},
			wantDelOperationResults: []*client.OpResult{
				routeInstallResult(ateDstNetCIDR, constants.Delete),
				fluent.OperationResult().
					WithNextHopGroupOperation(nhgIndex).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Delete).
					AsResult(),
				fluent.OperationResult().
					WithNextHopOperation(nhIndex).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Delete).
					AsResult(),
			},
			startAddress: startAddress,
			v6Traffic:    routeType != IPv4,
		}
	} else {
		return testCase{
			desc: desc,
			entries: []fluent.GRIBIEntry{
				fluent.NextHopEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(nhIndex3).WithIPAddress(ateIndirectNH),
				fluent.NextHopGroupEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithID(nhgIndex3).AddNextHop(nhIndex3, 1),
				fluentEntry(ateDstNetCIDR, nhgIndex3),
				fluent.NextHopEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(nhIndex2).WithIPAddress(destIP),
				fluent.NextHopGroupEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithID(nhgIndex2).AddNextHop(nhIndex2, 1),
				fluentEntry(ateIndirectNHCIDR, nhgIndex2),
			},
			wantAddOperationResults: []*client.OpResult{
				fluent.OperationResult().
					WithNextHopOperation(nhIndex3).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Add).
					AsResult(),
				fluent.OperationResult().
					WithNextHopGroupOperation(nhgIndex3).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Add).
					AsResult(),
				routeInstallResult(ateDstNetCIDR, constants.Add),
				fluent.OperationResult().
					WithNextHopOperation(nhIndex2).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Add).
					AsResult(),
				fluent.OperationResult().
					WithNextHopGroupOperation(nhgIndex2).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Add).
					AsResult(),
				routeInstallResult(ateIndirectNHCIDR, constants.Add),
			},
			wantDelOperationResults: []*client.OpResult{
				fluent.OperationResult().
					WithNextHopOperation(nhIndex3).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Delete).
					AsResult(),
				fluent.OperationResult().
					WithNextHopGroupOperation(nhgIndex3).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Delete).
					AsResult(),
				routeInstallResult(ateDstNetCIDR, constants.Delete),
				fluent.OperationResult().
					WithNextHopOperation(nhIndex2).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Delete).
					AsResult(),
				fluent.OperationResult().
					WithNextHopGroupOperation(nhgIndex2).
					WithProgrammingResult(fluent.InstalledInFIB).
					WithOperationType(constants.Delete).
					AsResult(),
				routeInstallResult(ateIndirectNHCIDR, constants.Delete),
			},
			startAddress: startAddress,
			v6Traffic:    routeType != IPv4,
		}
	}
}

// TestGRIBIEntry tests single IPv4 and IPv6 gRIBI forwarding entries.
func TestGRIBIEntry(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	configureDUT(t, dut)

	ate := ondatra.ATE(t, "ate")
	otg := ate.OTG()
	otgConfig := configureOTG(t, ate)
	otg.PushConfig(t, otgConfig)
	otg.StartProtocols(t)

	cases := []testCase{
		newTestCase(t, "single-next-hop-IPv4", "198.51.100.0", IPv4, false),
		newTestCase(t, "recursive-next-hop-IPv4", "198.51.100.64", IPv4, true),
		newTestCase(t, "single-next-hop-IPv6", "2003::6464", IPv6, false),
		newTestCase(t, "recursive-next-hop-IPv6", "2003::", IPv6, true),
		newTestCase(t, "single-next-hop-IPv4-mapped-IPv6", "2003::3232", IPv4MappedIPv6, false),
		newTestCase(t, "recursive-next-hop-IPv4-mapped-IPv6", "2003::2424", IPv4MappedIPv6, true),
	}
	for _, tc := range cases {
		var txPkts, rxPkts uint64
		t.Run(tc.desc, func(t *testing.T) {
			gribic := dut.RawAPIs().GRIBI(t)
			c := fluent.NewClient()
			c.Connection().WithStub(gribic).
				WithRedundancyMode(fluent.ElectedPrimaryClient).
				WithPersistence().
				WithFIBACK().
				WithInitialElectionID(1, 0)
			ctx := context.Background()
			c.Start(ctx, t)
			defer c.Stop(t)
			c.StartSending(ctx, t)
			if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
				t.Fatalf("Await got error during session negotiation: %v", err)
			}

			c.Modify().AddEntry(t, tc.entries...)
			if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
				t.Fatalf("Await got error for entries: %v", err)
			}

			for _, wantResult := range tc.wantAddOperationResults {
				chk.HasResult(t, c.Results(t), wantResult, chk.IgnoreOperationID())
			}

			testTrafficFn := testTraffic
			flowName := "Flow"
			if tc.v6Traffic {
				testTrafficFn = testTrafficv6
				flowName = "Flow2"
			}

			if loss := testTrafficFn(t, otg, atePort1, atePort2, tc.startAddress, 5*time.Second); loss > 1 {
				t.Errorf("Loss: got %g, want <= 1", loss)
			}

			// counters are not erased, so have to accumulate the packets from previous subtests.
			txPkts += gnmi.Get(t, otg, gnmi.OTG().Flow(flowName).Counters().OutPkts().State())
			rxPkts += gnmi.Get(t, otg, gnmi.OTG().Flow(flowName).Counters().InPkts().State())
			testCounters(t, dut, txPkts, rxPkts)

			slices.Reverse(tc.entries)
			c.Modify().DeleteEntry(t, tc.entries...)
			if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
				t.Fatalf("Await got error for entries: %v", err)
			}

			for _, wantResult := range tc.wantDelOperationResults {
				chk.HasResult(t, c.Results(t), wantResult, chk.IgnoreOperationID())
			}

			if loss := testTrafficFn(t, otg, atePort1, atePort2, tc.startAddress, 5*time.Second); loss != 100 {
				t.Errorf("Loss: got %g, want 100", loss)
			}

			// TODO: Test flush once it's implemented.
			gribic.Flush(context.Background(), &gribipb.FlushRequest{
				NetworkInstance: &gribipb.FlushRequest_All{All: &gribipb.Empty{}},
			})
		})
	}
}
