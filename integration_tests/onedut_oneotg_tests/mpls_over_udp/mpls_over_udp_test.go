// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License fbuildor the specific language governing permissions and
// limitations under the License.

package integration_test

import (
	"context"
	"net/netip"
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/gribigo/chk"
	"github.com/openconfig/gribigo/client"
	"github.com/openconfig/gribigo/constants"
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/otg"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

	gribipb "github.com/openconfig/gribi/v1/proto/service"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	ni "github.com/openconfig/lemming/gnmi/oc/networkinstance"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/lemming/internal/binding"
)

const (
	ipv4PrefixLen = 30
	ipv6PrefixLen = 99
)

const (
	// IPv6
	ateDstNetCIDRv6     = "2003::/48"
	ateIndirectNHv6     = "2002::"
	ateIndirectNHCIDRv6 = ateIndirectNHv6 + "/48"
	// Common attributes
	nhIndex        = 1
	nhgIndex       = 42
	mplsLabel      = uint64(100)     // Example MPLS label
	udpPort        = uint16(6635)    // Example UDP port
	outerIPv6Src   = "2001:f:a:1::0" // Example outer IPv6 src, adjust as needed
	outerIPv6Dst   = "2001:f:c:e::2" // Example outer IPv6 dst, adjust as needed
	ipTTL          = uint8(1)
	startAddressV6 = "2003::6464"
	flowNameV6     = "FlowV6"
	dscp           = 46
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

var destIP = atePort2.IPv6

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.KNE(".."))
}

// configureDUT configures port1 and port2 on the DUT.
func configureDUT(t *testing.T, dut *ondatra.DUTDevice) {
	p1 := dut.Port(t, "port1")
	p2 := dut.Port(t, "port2")
	portList := []*ondatra.Port{p1, p2}

	// Added for loop
	for idx, a := range []attrs.Attributes{dutPort1, dutPort2} {
		p := portList[idx]
		intf := a.NewOCInterface(p.Name(), dut)
		if p.PMD() == ondatra.PMD100GBASEFR && dut.Vendor() != ondatra.CISCO && dut.Vendor() != ondatra.JUNIPER {
			e := intf.GetOrCreateEthernet()
			e.AutoNegotiate = ygot.Bool(false)
			e.DuplexMode = oc.Ethernet_DuplexMode_FULL
			e.PortSpeed = oc.IfEthernet_ETHERNET_SPEED_SPEED_100GB
		}
		gnmi.Replace(t, dut, ocpath.Root().Interface(p.Name()).Config(), intf) // Modified
	}

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv4().Address(dutPort1.IPv4).Ip().State(), time.Minute, dutPort1.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv4().Address(dutPort2.IPv4).Ip().State(), time.Minute, dutPort2.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv6().Address(dutPort1.IPv6).Ip().State(), time.Minute, dutPort1.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv6().Address(dutPort2.IPv6).Ip().State(), time.Minute, dutPort2.IPv6)
}

// WaitForARP waits for ARP to resolve on all OTG interfaces for a given ipType, which is
// either "IPv4" or "IPv6".
func WaitForARP(t *testing.T, otg *otg.OTG, c gosnappi.Config, ipType string) {
	intfs := []string{}
	for _, d := range c.Devices().Items() {
		Eth := d.Ethernets().Items()[0]
		intfs = append(intfs, Eth.Name())
	}

	for _, intf := range intfs {
		switch ipType {
		case "IPv4":
			got, ok := gnmi.WatchAll(t, otg, gnmi.OTG().Interface(intf).Ipv4NeighborAny().LinkLayerAddress().State(), 2*time.Minute, func(val *ygnmi.Value[string]) bool {
				return val.IsPresent()
			}).Await(t)
			if !ok {
				t.Fatalf("Did not receive OTG Neighbor entry for interface %s, last got: %v", intf, got)
			}
		case "IPv6":
			got, ok := gnmi.WatchAll(t, otg, gnmi.OTG().Interface(intf).Ipv6NeighborAny().LinkLayerAddress().State(), 2*time.Minute, func(val *ygnmi.Value[string]) bool {
				return val.IsPresent()
			}).Await(t)
			if !ok {
				t.Fatalf("Did not receive OTG Neighbor entry for interface %s, last got: %v", intf, got)
			}
		}
	}
}

// configureOTG configures port1 and port2 on the ATE.
func configureOTG(t *testing.T, ate *ondatra.ATEDevice) gosnappi.Config {
	top := gosnappi.NewConfig()

	p1 := ate.Port(t, "port1")
	p2 := ate.Port(t, "port2")

	atePort1.AddToOTG(top, p1, &dutPort1)
	atePort2.AddToOTG(top, p2, &dutPort2)

	pmd100GFRPorts := []string{}
	for _, p := range top.Ports().Items() {
		port := ate.Port(t, p.Name())
		if port.PMD() == ondatra.PMD100GBASEFR {
			pmd100GFRPorts = append(pmd100GFRPorts, port.ID())
		}
	}
	// Disable FEC for 100G-FR ports because Novus does not support it.
	if len(pmd100GFRPorts) > 0 {
		l1Settings := top.Layer1().Add().SetName("L1").SetPortNames(pmd100GFRPorts)
		l1Settings.SetAutoNegotiate(true).SetIeeeMediaDefaults(false).SetSpeed("speed_100_gbps")
		autoNegotiate := l1Settings.AutoNegotiation()
		autoNegotiate.SetRsFec(false)
	}

	return top
}

// testTrafficv6 generates traffic flow from source network to
// destination network via srcEndPoint to dstEndPoint and checks for
// packet loss and returns loss percentage as float.
func testTrafficv6(t *testing.T, otg *otg.OTG, srcEndPoint, dstEndPoint attrs.Attributes, startAddress string, dur time.Duration) float32 {
	WaitForARP(t, otg, otg.FetchConfig(t), "IPv6")
	top := otg.FetchConfig(t)
	top.Flows().Clear().Items()
	flowipv6 := top.Flows().Add().SetName(flowNameV6)
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

	txPkts := gnmi.Get(t, otg, gnmi.OTG().Flow(flowNameV6).Counters().OutPkts().State())
	rxPkts := gnmi.Get(t, otg, gnmi.OTG().Flow(flowNameV6).Counters().InPkts().State())
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

func routeInstallResult(t *testing.T, prefix string, c constants.OpType) *client.OpResult {
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

func checkEncapHeaders(t *testing.T, dut *ondatra.DUTDevice, nhgPaths []*ni.NetworkInstance_Afts_NextHopGroupPath, wantEncapHeaders map[uint8]*oc.NetworkInstance_Afts_NextHop_EncapHeader) {
	for _, p := range nhgPaths {
		nhg, present := gnmi.Lookup(t, dut, p.State()).Val()
		if !present {
			return
		}
		nhs := nhg.NextHop
		for ind := range nhs {
			nhp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Afts().NextHop(ind)
			nh, present := gnmi.Lookup(t, dut, nhp.State()).Val()
			if !present {
				continue
			}
			ehs := nh.EncapHeader
			for i, eh := range ehs {
				if diff := cmp.Diff(eh, wantEncapHeaders[i]); diff != "" {
					t.Errorf("Diff (-got +want): %v", diff)
				}
			}
		}
	}
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

func TestMPLSOverUDPIPv6(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	configureDUT(t, dut)

	ate := ondatra.ATE(t, "ate")
	otg := ate.OTG()
	otgConfig := configureOTG(t, ate)
	t.Logf("Pushing config to ATE and starting protocols...")
	otg.PushConfig(t, otgConfig)
	t.Logf("starting protocols...")
	otg.StartProtocols(t)

	tests := []struct {
		desc                    string
		entries                 []fluent.GRIBIEntry
		nextHopGroupPaths       []*ni.NetworkInstance_Afts_NextHopGroupPath
		wantAddOperationResults []*client.OpResult
		wantAddEncapHeaders     map[uint8]*oc.NetworkInstance_Afts_NextHop_EncapHeader
		wantDelOperationResults []*client.OpResult
		wantDelEncapHeaders     map[uint8]*oc.NetworkInstance_Afts_NextHop_EncapHeader
		capturePort             string
		wantMPLSLabel           uint64
		wantOuterIP             string
	}{
		{
			desc: "mplsoudpv6",
			entries: []fluent.GRIBIEntry{
				fluent.NextHopEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithIndex(nhIndex).WithIPAddress(destIP).AddEncapHeader(
					fluent.MPLSEncapHeader().WithLabels(mplsLabel),
					fluent.UDPV6EncapHeader().WithDstUDPPort(uint64(udpPort)).WithSrcUDPPort(uint64(udpPort)).WithSrcIP(atePort1.IPv6).WithDstIP(atePort2.IPv6).WithDSCP(dscp).WithIPTTL(uint64(ipTTL)),
				),
				fluent.NextHopGroupEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithID(nhgIndex).AddNextHop(nhIndex, 1),
				fluent.IPv6Entry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
					WithPrefix(ateDstNetCIDRv6).WithNextHopGroup(nhgIndex),
			},
			nextHopGroupPaths: []*ni.NetworkInstance_Afts_NextHopGroupPath{
				(*ni.NetworkInstance_Afts_NextHopGroupPath)(ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Afts().NextHopGroup(nhgIndex)),
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
				routeInstallResult(t, ateDstNetCIDRv6, constants.Add),
			},
			wantAddEncapHeaders: map[uint8]*oc.NetworkInstance_Afts_NextHop_EncapHeader{
				1: {
					Index: ygot.Uint8(1),
					Type:  oc.AftTypes_EncapsulationHeaderType_MPLS,
					Mpls: &oc.NetworkInstance_Afts_NextHop_EncapHeader_Mpls{
						MplsLabelStack: []oc.NetworkInstance_Afts_NextHop_EncapHeader_Mpls_MplsLabelStack_Union{
							oc.UnionUint32(mplsLabel),
						},
					},
				},
				2: {
					Index: ygot.Uint8(2),
					Type:  oc.AftTypes_EncapsulationHeaderType_UDPV6,
					UdpV6: &oc.NetworkInstance_Afts_NextHop_EncapHeader_UdpV6{
						Dscp:       ygot.Uint8(dscp),
						DstIp:      ygot.String(atePort2.IPv6),
						DstUdpPort: ygot.Uint16(udpPort),
						IpTtl:      ygot.Uint8(ipTTL),
						SrcIp:      ygot.String(atePort1.IPv6),
						SrcUdpPort: ygot.Uint16(udpPort),
					},
				},
			},
			wantDelOperationResults: []*client.OpResult{
				routeInstallResult(t, ateDstNetCIDRv6, constants.Delete),
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
			wantDelEncapHeaders: map[uint8]*oc.NetworkInstance_Afts_NextHop_EncapHeader{},
			capturePort:         atePort2.Name,
			wantMPLSLabel:       mplsLabel,
			wantOuterIP:         atePort2.IPv6,
		},
	}
	for _, tc := range tests {
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

			t.Log("Sending ADD Modify request")
			c.Modify().AddEntry(t, tc.entries...)
			if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
				t.Fatalf("Await got error for entries: %v", err)
			}

			for _, wantResult := range tc.wantAddOperationResults {
				chk.HasResult(t, c.Results(t), wantResult, chk.IgnoreOperationID())
			}

			enableCapture(t, otg, otgConfig, tc.capturePort)
			startCapture(t, ate)
			if loss := testTrafficv6(t, otg, atePort1, atePort2, startAddressV6, 5*time.Second); loss > 1 {
				t.Errorf("Loss: got %g, want <= 1", loss)
			}
			stopCapture(t, ate)

			// TODO(tengyi): we need to wait LayerTypeAGUEVar0 updates in https://github.com/google/gopacket/blob/master/layers/layertypes.go
			// to include the packet validations and checks as required by TE-18.1.1 under featureprofiles.
			checkEncapHeaders(t, dut, tc.nextHopGroupPaths, tc.wantAddEncapHeaders)

			var txPkts, rxPkts uint64
			flowName := flowNameV6
			// counters are not erased, so have to accumulate the packets from previous subtests.
			txPkts += gnmi.Get(t, otg, gnmi.OTG().Flow(flowName).Counters().OutPkts().State())
			rxPkts += gnmi.Get(t, otg, gnmi.OTG().Flow(flowName).Counters().InPkts().State())
			testCounters(t, dut, txPkts, rxPkts)

			t.Log("Sending DELETE Modify request")
			slices.Reverse(tc.entries)
			c.Modify().DeleteEntry(t, tc.entries...)
			if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
				t.Fatalf("Await got error for entries: %v", err)
			}

			for _, wantResult := range tc.wantDelOperationResults {
				chk.HasResult(t, c.Results(t), wantResult, chk.IgnoreOperationID())
			}

			if loss := testTrafficv6(t, otg, atePort1, atePort2, startAddressV6, 5*time.Second); loss != 100 {
				t.Errorf("Loss: got %g, want 100", loss)
			}

			gribic.Flush(ctx, &gribipb.FlushRequest{
				NetworkInstance: &gribipb.FlushRequest_All{All: &gribipb.Empty{}},
			})
			checkEncapHeaders(t, dut, tc.nextHopGroupPaths, tc.wantDelEncapHeaders)
		})
	}
}
