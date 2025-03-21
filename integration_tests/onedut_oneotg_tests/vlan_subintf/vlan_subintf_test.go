// Copyright 2024 Google LLC
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

package vlan_subintf

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
	"github.com/openconfig/ondatra/otg"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/lemming/internal/binding"
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

const (
	TestNI = "TEST"
)

// configureDUT configures port1 and port2 on the DUT.
func configureDUT(t *testing.T, dut *ondatra.DUTDevice) {
	p1 := dut.Port(t, "port1")

	p1OC := dutPort1.NewOCInterface(p1.Name(), dut)
	p1OC.GetOrCreateSubinterface(1).GetOrCreateVlan().GetOrCreateMatch().GetOrCreateSingleTagged().SetVlanId(uint16(1))
	gnmi.Replace(t, dut, ocpath.Root().Interface(p1.Name()).Config(), p1OC)

	p2 := dut.Port(t, "port2")
	gnmi.Replace(t, dut, ocpath.Root().Interface(p2.Name()).Config(), dutPort2.NewOCInterface(p2.Name(), dut))

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv4().Address(dutPort1.IPv4).Ip().State(), time.Minute, dutPort1.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv4().Address(dutPort2.IPv4).Ip().State(), time.Minute, dutPort2.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv6().Address(dutPort1.IPv6).Ip().State(), time.Minute, dutPort1.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv6().Address(dutPort2.IPv6).Ip().State(), time.Minute, dutPort2.IPv6)

	gnmi.Update(t, dut, ocpath.Root().NetworkInstance(TestNI).Config(), &oc.NetworkInstance{
		Name: ygot.String(TestNI),
		Interface: map[string]*oc.NetworkInstance_Interface{
			p1.Name(): {
				Id:           ygot.String(p1.Name()),
				Interface:    ygot.String(p1.Name()),
				Subinterface: ygot.Uint32(1),
			},
		},
	})
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
func testTraffic(t *testing.T, otg *otg.OTG, srcEndPoint, dstEndPoint attrs.Attributes, dstIP string, dur time.Duration, vlanID *uint32) float32 {
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
	if vlanID != nil {
		flowipv4.Packet().Add().Vlan().Id().SetValue(*vlanID)
	}
	v4 := flowipv4.Packet().Add().Ipv4()
	v4.Src().SetValue(srcEndPoint.IPv4)
	v4.Dst().SetValue(dstIP)
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

// TestGRIBIEntry tests single IPv4 and IPv6 gRIBI forwarding entries.
func TestVLANVRF(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	configureDUT(t, dut)

	ate := ondatra.ATE(t, "ate")
	otg := ate.OTG()
	otgConfig := configureOTG(t, ate)
	otg.PushConfig(t, otgConfig)
	otg.StartProtocols(t)

	nhIndex := uint64(100)
	nhgIndex := uint64(10)

	tests := []struct {
		desc                    string
		entries                 []fluent.GRIBIEntry
		wantAddOperationResults []*client.OpResult
		flowDstAddr             string
		flowVlanID              uint32
		untaggedLoss            bool
		taggedLoss              bool
	}{{
		desc:       "tagged and untagged in different vrf",
		flowVlanID: 1,
		entries: []fluent.GRIBIEntry{
			fluent.NextHopEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
				WithIndex(nhIndex).WithIPAddress(atePort2.IPv4),
			fluent.NextHopGroupEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
				WithID(10).AddNextHop(nhIndex, 1),
			fluent.IPv4Entry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
				WithPrefix("198.51.100.1/32").WithNextHopGroup(nhgIndex),
		},
		flowDstAddr: "198.51.100.1",
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
				WithIPv4Operation("198.51.100.1/32").
				WithProgrammingResult(fluent.InstalledInFIB).
				WithOperationType(constants.Add).
				AsResult(),
		},
		taggedLoss: true,
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
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

			c.Modify().AddEntry(t, tt.entries...)
			if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
				t.Fatalf("Await got error for entries: %v", err)
			}

			for _, wantResult := range tt.wantAddOperationResults {
				chk.HasResult(t, c.Results(t), wantResult, chk.IgnoreOperationID())
			}

			t.Log("starting untagged flow")
			// First send traffic without VLAN tag.
			loss := testTraffic(t, otg, atePort1, atePort2, tt.flowDstAddr, 5*time.Second, nil)
			if tt.untaggedLoss && loss < 0.99 {
				t.Errorf("untaggedLoss untagged to flow to fail, got loss %f", loss)
			} else if !tt.taggedLoss && loss > 0.01 {
				t.Errorf("unexpected untagged to flow to pass, got loss %f", loss)
			}
			// Next send traffic with VLAN tag.
			t.Log("starting tagged flow")
			loss = testTraffic(t, otg, atePort1, atePort2, tt.flowDstAddr, 5*time.Second, &tt.flowVlanID)
			if tt.taggedLoss && loss < 0.99 {
				t.Errorf("unexpected tagged to flow to fail, got loss %f", loss)
			} else if !tt.taggedLoss && loss > 0.01 {
				t.Errorf("unexpected tagged to flow to pass, got loss %f", loss)
			}

			slices.Reverse(tt.entries)
			c.Modify().DeleteEntry(t, tt.entries...)
			if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
				t.Fatalf("Await got error for entries: %v", err)
			}
		})
	}
}
