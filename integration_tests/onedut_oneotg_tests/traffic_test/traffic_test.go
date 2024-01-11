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
	"github.com/openconfig/ygnmi/ygnmi"

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
	ipv4PrefixLen          = 30
	ateDstNetCIDR          = "198.51.100.0/24"
	ateIndirectNH          = "203.0.113.1"
	ateIndirectNHCIDR      = ateIndirectNH + "/32"
	nhIndex                = 1
	nhgIndex               = 42
	nhIndex2               = 2
	nhgIndex2              = 52
	nhIndex3               = 3
	nhgIndex3              = 62
	defaultNetworkInstance = "DEFAULT"
)

var (
	dutPort1 = attrs.Attributes{
		Desc:    "dutPort1",
		IPv4:    "192.0.2.1",
		IPv4Len: ipv4PrefixLen,
	}

	atePort1 = attrs.Attributes{
		Name:    "port1",
		MAC:     "02:00:01:01:01:01",
		IPv4:    "192.0.2.2",
		IPv4Len: ipv4PrefixLen,
	}

	dutPort2 = attrs.Attributes{
		Desc:    "dutPort2",
		IPv4:    "192.0.2.5",
		IPv4Len: ipv4PrefixLen,
	}

	atePort2 = attrs.Attributes{
		Name:    "port2",
		MAC:     "02:00:02:01:01:01",
		IPv4:    "192.0.2.6",
		IPv4Len: ipv4PrefixLen,
	}
)

// configureATE configures port1 and port2 on the ATE.
func configureATE(t *testing.T, ate *ondatra.ATEDevice) gosnappi.Config {
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
}

// testTraffic generates traffic flow from source network to
// destination network via srcEndPoint to dstEndPoint and checks for
// packet loss and returns loss percentage as float.
func testTraffic(t *testing.T, ate *ondatra.ATEDevice, top gosnappi.Config, srcEndPoint, dstEndPoint attrs.Attributes, dur time.Duration) float32 {
	otg := ate.OTG()
	waitOTGARPEntry(t)
	top.Flows().Clear().Items()
	flowipv4 := top.Flows().Add().SetName("Flow")
	flowipv4.Metrics().SetEnable(true)
	flowipv4.TxRx().Device().
		SetTxNames([]string{srcEndPoint.Name + ".IPv4"}).
		SetRxNames([]string{dstEndPoint.Name + ".IPv4"})
	flowipv4.Duration().SetChoice("continuous")
	flowipv4.Packet().Add().Ethernet()
	v4 := flowipv4.Packet().Add().Ipv4()
	v4.Src().SetValue(srcEndPoint.IPv4)
	v4.Dst().Increment().SetStart("198.51.100.0").SetCount(250)
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

// TestIPv4Entry tests a single IPv4Entry forwarding entry.
func TestIPv4Entry(t *testing.T) {
	ctx := context.Background()

	dut := ondatra.DUT(t, "dut")
	configureDUT(t, dut)

	ate := ondatra.ATE(t, "ate")
	ateTop := configureATE(t, ate)
	ate.OTG().PushConfig(t, ateTop)

	cases := []struct {
		desc                    string
		entries                 []fluent.GRIBIEntry
		wantAddOperationResults []*client.OpResult
		wantDelOperationResults []*client.OpResult
	}{{
		desc: "single-next-hop",
		entries: []fluent.GRIBIEntry{
			fluent.NextHopEntry().WithNetworkInstance(defaultNetworkInstance).
				WithIndex(nhIndex).WithIPAddress(atePort2.IPv4),
			fluent.NextHopGroupEntry().WithNetworkInstance(defaultNetworkInstance).
				WithID(nhgIndex).AddNextHop(nhIndex, 1),
			fluent.IPv4Entry().WithNetworkInstance(defaultNetworkInstance).
				WithPrefix(ateDstNetCIDR).WithNextHopGroup(nhgIndex),
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
			fluent.OperationResult().
				WithIPv4Operation(ateDstNetCIDR).
				WithProgrammingResult(fluent.InstalledInFIB).
				WithOperationType(constants.Add).
				AsResult(),
		},
		wantDelOperationResults: []*client.OpResult{
			fluent.OperationResult().
				WithIPv4Operation(ateDstNetCIDR).
				WithProgrammingResult(fluent.InstalledInFIB).
				WithOperationType(constants.Delete).
				AsResult(),
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
	}, {
		desc: "recursive-next-hop",
		entries: []fluent.GRIBIEntry{
			// Add an IPv4Entry for 198.51.100.0/24 pointing to 203.0.113.1/32.
			fluent.NextHopEntry().WithNetworkInstance(defaultNetworkInstance).
				WithIndex(nhIndex3).WithIPAddress(ateIndirectNH),
			fluent.NextHopGroupEntry().WithNetworkInstance(defaultNetworkInstance).
				WithID(nhgIndex3).AddNextHop(nhIndex3, 1),
			fluent.IPv4Entry().WithNetworkInstance(defaultNetworkInstance).
				WithPrefix(ateDstNetCIDR).WithNextHopGroup(nhgIndex3),
			// Add an IPv4Entry for 203.0.113.1/32 pointing to 192.0.2.6.
			fluent.NextHopEntry().WithNetworkInstance(defaultNetworkInstance).
				WithIndex(nhIndex2).WithIPAddress(atePort2.IPv4),
			fluent.NextHopGroupEntry().WithNetworkInstance(defaultNetworkInstance).
				WithID(nhgIndex2).AddNextHop(nhIndex2, 1),
			fluent.IPv4Entry().WithNetworkInstance(defaultNetworkInstance).
				WithPrefix(ateIndirectNHCIDR).WithNextHopGroup(nhgIndex2),
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
			fluent.OperationResult().
				WithIPv4Operation(ateDstNetCIDR).
				WithProgrammingResult(fluent.InstalledInFIB).
				WithOperationType(constants.Add).
				AsResult(),
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
			fluent.OperationResult().
				WithIPv4Operation(ateIndirectNHCIDR).
				WithProgrammingResult(fluent.InstalledInFIB).
				WithOperationType(constants.Add).
				AsResult(),
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
			fluent.OperationResult().
				WithIPv4Operation(ateDstNetCIDR).
				WithProgrammingResult(fluent.InstalledInFIB).
				WithOperationType(constants.Delete).
				AsResult(),
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
			fluent.OperationResult().
				WithIPv4Operation(ateIndirectNHCIDR).
				WithProgrammingResult(fluent.InstalledInFIB).
				WithOperationType(constants.Delete).
				AsResult(),
		},
	}}
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

			// Send some traffic to make sure neighbor cache is warmed up on the dut.
			testTraffic(t, ate, ateTop, atePort1, atePort2, 1*time.Second)

			if loss := testTraffic(t, ate, ateTop, atePort1, atePort2, 5*time.Second); loss > 1 {
				t.Errorf("Loss: got %g, want <= 1", loss)
			}

			otg := ate.OTG()
			// counters are not erased, so have to accumulate the packets from previous subtests.
			txPkts += gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().OutPkts().State())
			rxPkts += gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().InPkts().State())
			testCounters(t, dut, txPkts, rxPkts)

			slices.Reverse(tc.entries)
			c.Modify().DeleteEntry(t, tc.entries...)
			if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
				t.Fatalf("Await got error for entries: %v", err)
			}

			// Send some traffic to make sure neighbor cache is warmed up on the dut.
			testTraffic(t, ate, ateTop, atePort1, atePort2, 1*time.Second)

			for _, wantResult := range tc.wantDelOperationResults {
				chk.HasResult(t, c.Results(t), wantResult, chk.IgnoreOperationID())
			}

			if loss := testTraffic(t, ate, ateTop, atePort1, atePort2, 5*time.Second); loss != 100 {
				t.Errorf("Loss: got %g, want 100", loss)
			}

			// TODO: Test flush once it's implemented.
			gribic.Flush(context.Background(), &gribipb.FlushRequest{
				NetworkInstance: &gribipb.FlushRequest_All{All: &gribipb.Empty{}},
			})
		})
	}
}
