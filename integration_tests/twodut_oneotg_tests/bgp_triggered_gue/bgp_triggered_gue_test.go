/*
 Copyright 2022 Google LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/gribigo/chk"
	"github.com/openconfig/gribigo/client"
	"github.com/openconfig/gribigo/constants"
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/gnmi/otg/otgpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

	gribipb "github.com/openconfig/gribi/v1/proto/service"

	kinit "github.com/openconfig/ondatra/knebind/init"
)

// Settings for configuring the baseline testbed with the test
// topology.
//
// The testbed consists of,
//
//   - ate:port1 -> dut:port1 subnet 192.0.2.0/30
//   - ate:port2 -> dut:port2 subnet 192.0.2.4/30
//   - ate:port3 -> dut2:port1 subnet 203.0.113.0/30
//   - dut:port3 -> dut2:port2 subnet 192.1.2.8/30
const (
	ipv4PrefixLen          = 30
	ateDstNetCIDR1         = "198.51.0.0/24"
	ateDstNetCIDR2         = "198.51.2.0/24"
	ateDstNetCIDR3         = "198.51.4.0/24"
	ateIndirectNH1         = "203.0.113.1"
	ateIndirectNH2         = "203.0.113.5"
	ateIndirectNH3         = "203.0.113.9"
	ateIndirectNHCIDR      = "203.0.113.0/24"
	nhIndex1               = 1
	nhIndex2               = 2
	nhIndex3               = 3
	nhgIndex1              = 42
	nhgIndex2              = 43
	nhgIndex3              = 44
	defaultNetworkInstance = "DEFAULT"

	dutAS  = 64500
	dut2AS = 64501
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, kinit.Init)
}

// Attributes bundles some common attributes for devices and/or interfaces.
// It provides helpers to generate appropriate configuration for OpenConfig
// and for an ATETopology.  All fields are optional; only those that are
// non-empty will be set when configuring an interface.
type Attributes struct {
	IPv4    string
	IPv6    string
	MAC     string
	Name    string // Interface name, only applied to ATE ports.
	Desc    string // Description, only applied to DUT interfaces.
	IPv4Len uint8  // Prefix length for IPv4.
	IPv6Len uint8  // Prefix length for IPv6.
	MTU     uint16
}

var (
	dutPort1 = Attributes{
		Desc:    "dutPort1",
		IPv4:    "192.0.2.1",
		IPv4Len: ipv4PrefixLen,
	}

	atePort1 = Attributes{
		Name:    "port1",
		MAC:     "02:00:01:01:01:01",
		IPv4:    "192.0.2.2",
		IPv4Len: ipv4PrefixLen,
	}

	dutPort2 = Attributes{
		Desc:    "dutPort2",
		IPv4:    "192.0.2.5",
		IPv4Len: ipv4PrefixLen,
	}

	atePort2 = Attributes{
		Name:    "port2",
		MAC:     "02:00:02:01:01:01",
		IPv4:    "192.0.2.6",
		IPv4Len: ipv4PrefixLen,
	}

	dutPort3 = Attributes{
		Desc:    "dutPort3",
		IPv4:    "192.0.2.9",
		IPv4Len: ipv4PrefixLen,
	}

	atePort3 = Attributes{
		Name:    "port3",
		MAC:     "02:00:03:01:01:01",
		IPv4:    ateIndirectNH3,
		IPv4Len: 28,
	}

	dut2Port1 = Attributes{
		Desc: "dut2Port1",
		IPv4: "203.0.113.2",
		// Make sure the ATE indirect prefixes that are to be exchanged
		// to DUT1 are actually resolvable at DUT2.
		IPv4Len: 28,
	}

	dut2Port2 = Attributes{
		Desc:    "dut2Port2",
		IPv4:    "192.0.2.10",
		IPv4Len: ipv4PrefixLen,
	}
)

// configureATE configures port1 and port2 on the ATE.
func configureATE(t *testing.T, ate *ondatra.ATEDevice) gosnappi.Config {
	t.Helper()
	otg := ate.OTG()
	top := otg.NewConfig(t)

	top.Ports().Add().SetName(atePort1.Name)
	i1 := top.Devices().Add().SetName(atePort1.Name)
	eth1 := i1.Ethernets().Add().SetName(atePort1.Name + ".Eth").
		SetPortName(i1.Name()).SetMac(atePort1.MAC)
	eth1.Ipv4Addresses().Add().SetName(i1.Name() + ".IPv4").
		SetAddress(atePort1.IPv4).SetGateway(dutPort1.IPv4).
		SetPrefix(int32(atePort1.IPv4Len))

	top.Ports().Add().SetName(atePort2.Name)
	i2 := top.Devices().Add().SetName(atePort2.Name)
	eth2 := i2.Ethernets().Add().SetName(atePort2.Name + ".Eth").
		SetPortName(i2.Name()).SetMac(atePort2.MAC)
	eth2.Ipv4Addresses().Add().SetName(i2.Name() + ".IPv4").
		SetAddress(atePort2.IPv4).SetGateway(dutPort2.IPv4).
		SetPrefix(int32(atePort2.IPv4Len))
	// Configure capture format.
	top.Captures().Add().SetName("ca").SetPortNames([]string{atePort2.Name}).SetFormat(gosnappi.CaptureFormat.PCAP)

	top.Ports().Add().SetName(atePort3.Name)
	i3 := top.Devices().Add().SetName(atePort3.Name)
	eth3 := i3.Ethernets().Add().SetName(atePort3.Name + ".Eth").
		SetPortName(i3.Name()).SetMac(atePort3.MAC)
	eth3.Ipv4Addresses().Add().SetName(i3.Name() + ".IPv4").
		SetAddress(atePort3.IPv4).SetGateway(dut2Port1.IPv4).
		SetPrefix(int32(atePort3.IPv4Len))
	return top
}

var gatewayMap = map[Attributes]Attributes{
	atePort1: dutPort1,
	atePort2: dutPort2,
	atePort3: dut2Port1,
}

// configInterfaceDUT configures the interface with the Addrs.
func configInterfaceDUT(i *oc.Interface, a *Attributes) *oc.Interface {
	i.Description = ygot.String(a.Desc)
	i.Type = oc.IETFInterfaces_InterfaceType_ethernetCsmacd

	s := i.GetOrCreateSubinterface(0)
	s.Enabled = ygot.Bool(true)
	s4 := s.GetOrCreateIpv4()
	s4a := s4.GetOrCreateAddress(a.IPv4)
	s4a.PrefixLength = ygot.Uint8(a.IPv4Len)

	return i
}

// configureDUT1 configures ports on DUT1.
func configureDUT1(t *testing.T, dut *ondatra.DUTDevice) {
	p1 := dut.Port(t, "port1")
	i1 := &oc.Interface{Name: ygot.String(p1.Name())}
	gnmi.Replace(t, dut, ocpath.Root().Interface(p1.Name()).Config(), configInterfaceDUT(i1, &dutPort1))

	p2 := dut.Port(t, "port2")
	i2 := &oc.Interface{Name: ygot.String(p2.Name())}
	gnmi.Replace(t, dut, ocpath.Root().Interface(p2.Name()).Config(), configInterfaceDUT(i2, &dutPort2))

	p3 := dut.Port(t, "port3")
	i3 := &oc.Interface{Name: ygot.String(p3.Name())}
	gnmi.Replace(t, dut, ocpath.Root().Interface(p3.Name()).Config(), configInterfaceDUT(i3, &dutPort3))

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv4().Address(dutPort1.IPv4).Ip().State(), time.Minute, dutPort1.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv4().Address(dutPort2.IPv4).Ip().State(), time.Minute, dutPort2.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port3").Name()).Subinterface(0).Ipv4().Address(dutPort3.IPv4).Ip().State(), time.Minute, dutPort3.IPv4)
}

// configureDUT2 configures ports on DUT2.
func configureDUT2(t *testing.T, dut *ondatra.DUTDevice) {
	p1 := dut.Port(t, "port1")
	i1 := &oc.Interface{Name: ygot.String(p1.Name())}
	gnmi.Replace(t, dut, ocpath.Root().Interface(p1.Name()).Config(), configInterfaceDUT(i1, &dut2Port1))

	p2 := dut.Port(t, "port2")
	i2 := &oc.Interface{Name: ygot.String(p2.Name())}
	gnmi.Replace(t, dut, ocpath.Root().Interface(p2.Name()).Config(), configInterfaceDUT(i2, &dut2Port2))

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv4().Address(dut2Port1.IPv4).Ip().State(), time.Minute, dut2Port1.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv4().Address(dut2Port2.IPv4).Ip().State(), time.Minute, dut2Port2.IPv4)
}

func waitOTGARPEntry(t *testing.T) {
	ate := ondatra.ATE(t, "ate")

	ate.OTG().Telemetry().InterfaceAny().Ipv4NeighborAny().LinkLayerAddress()
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
func testTraffic(t *testing.T, ate *ondatra.ATEDevice, top gosnappi.Config, srcEndPoint, dstEndPoint Attributes, startingIP string) float32 {
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
	v4.Dst().Increment().SetStart(startingIP).SetCount(250)
	otg.PushConfig(t, top)

	otg.StartTraffic(t)
	time.Sleep(15 * time.Second)
	t.Logf("Stop traffic")
	otg.StopTraffic(t)

	time.Sleep(5 * time.Second)

	txPkts := otg.Telemetry().Flow("Flow").Counters().OutPkts().Get(t)
	rxPkts := otg.Telemetry().Flow("Flow").Counters().InPkts().Get(t)
	lossPct := (txPkts - rxPkts) * 100 / txPkts
	return float32(lossPct)
}

// awaitTimeout calls a fluent client Await, adding a timeout to the context.
func awaitTimeout(ctx context.Context, c *fluent.GRIBIClient, t testing.TB, timeout time.Duration) error {
	subctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Await(subctx, t)
}

func bgpWithNbr(as uint32, routerID string, nbr *oc.NetworkInstance_Protocol_Bgp_Neighbor) *oc.NetworkInstance_Protocol_Bgp {
	bgp := &oc.NetworkInstance_Protocol_Bgp{}
	bgp.GetOrCreateGlobal().As = ygot.Uint32(as)
	if routerID != "" {
		bgp.Global.RouterId = ygot.String(routerID)
	}
	bgp.AppendNeighbor(nbr)
	return bgp
}

func configureGRIBIEntry(t *testing.T, dut *ondatra.DUTDevice, entries []fluent.GRIBIEntry) *fluent.GRIBIClient {
	t.Helper()
	gribic := dut.RawAPIs().GRIBI().Default(t)
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

	c.Modify().AddEntry(t, entries...)
	if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
		t.Fatalf("Await got error for entries: %v", err)
	}
	return c
}

func TestBGPTriggeredGUE(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	configureDUT1(t, dut)
	dut2 := ondatra.DUT(t, "dut2")
	configureDUT2(t, dut2)

	ate := ondatra.ATE(t, "ate")
	ateTop := configureATE(t, ate)
	ate.OTG().PushConfig(t, ateTop)

	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()

	// Remove any existing BGP config
	gnmi.Delete(t, dut, bgpPath.Config())
	gnmi.Delete(t, dut2, bgpPath.Config())

	dutEntries := []fluent.GRIBIEntry{
		// Add an IPv4Entry for 203.0.113.0/24 pointing to 192.0.2.6.
		fluent.NextHopEntry().WithNetworkInstance(defaultNetworkInstance).
			WithIndex(nhIndex1).WithIPAddress(atePort2.IPv4),
		fluent.NextHopGroupEntry().WithNetworkInstance(defaultNetworkInstance).
			WithID(nhgIndex1).AddNextHop(nhIndex1, 1),
		fluent.IPv4Entry().WithNetworkInstance(defaultNetworkInstance).
			WithPrefix(ateIndirectNHCIDR).WithNextHopGroup(nhgIndex1),
	}
	c := configureGRIBIEntry(t, dut, dutEntries)

	dut2Entries := []fluent.GRIBIEntry{
		// Add an IPv4Entry for 198.51.0.0/24 pointing to 203.0.113.1.
		fluent.NextHopEntry().WithNetworkInstance(defaultNetworkInstance).
			WithIndex(nhIndex1).WithIPAddress(ateIndirectNH1),
		fluent.NextHopGroupEntry().WithNetworkInstance(defaultNetworkInstance).
			WithID(nhgIndex1).AddNextHop(nhIndex1, 1),
		fluent.IPv4Entry().WithNetworkInstance(defaultNetworkInstance).
			WithPrefix(ateDstNetCIDR1).WithNextHopGroup(nhgIndex1),
		// Add an IPv4Entry for 198.51.2.0/24 pointing to 203.0.113.5.
		fluent.NextHopEntry().WithNetworkInstance(defaultNetworkInstance).
			WithIndex(nhIndex2).WithIPAddress(ateIndirectNH2),
		fluent.NextHopGroupEntry().WithNetworkInstance(defaultNetworkInstance).
			WithID(nhgIndex2).AddNextHop(nhIndex2, 1),
		fluent.IPv4Entry().WithNetworkInstance(defaultNetworkInstance).
			WithPrefix(ateDstNetCIDR2).WithNextHopGroup(nhgIndex2),
		// Add an IPv4Entry for 198.51.4.0/24 pointing to 203.0.113.9.
		fluent.NextHopEntry().WithNetworkInstance(defaultNetworkInstance).
			WithIndex(nhIndex3).WithIPAddress(ateIndirectNH3),
		fluent.NextHopGroupEntry().WithNetworkInstance(defaultNetworkInstance).
			WithID(nhgIndex3).AddNextHop(nhIndex3, 1),
		fluent.IPv4Entry().WithNetworkInstance(defaultNetworkInstance).
			WithPrefix(ateDstNetCIDR3).WithNextHopGroup(nhgIndex3),
	}
	c2 := configureGRIBIEntry(t, dut2, dut2Entries)

	wantOperationResultsDUT1 := []*client.OpResult{
		fluent.OperationResult().
			WithNextHopOperation(nhIndex1).
			WithProgrammingResult(fluent.InstalledInFIB).
			WithOperationType(constants.Add).
			AsResult(),
		fluent.OperationResult().
			WithNextHopGroupOperation(nhgIndex1).
			WithProgrammingResult(fluent.InstalledInFIB).
			WithOperationType(constants.Add).
			AsResult(),
		fluent.OperationResult().
			WithIPv4Operation(ateIndirectNHCIDR).
			WithProgrammingResult(fluent.InstalledInFIB).
			WithOperationType(constants.Add).
			AsResult(),
	}

	wantOperationResultsDUT2 := []*client.OpResult{
		fluent.OperationResult().
			WithNextHopOperation(nhIndex1).
			WithProgrammingResult(fluent.InstalledInFIB).
			WithOperationType(constants.Add).
			AsResult(),
		fluent.OperationResult().
			WithNextHopGroupOperation(nhgIndex1).
			WithProgrammingResult(fluent.InstalledInFIB).
			WithOperationType(constants.Add).
			AsResult(),
		fluent.OperationResult().
			WithIPv4Operation(ateDstNetCIDR1).
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
			WithIPv4Operation(ateDstNetCIDR2).
			WithProgrammingResult(fluent.InstalledInFIB).
			WithOperationType(constants.Add).
			AsResult(),
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
			WithIPv4Operation(ateDstNetCIDR3).
			WithProgrammingResult(fluent.InstalledInFIB).
			WithOperationType(constants.Add).
			AsResult(),
	}

	for _, wantResult := range wantOperationResultsDUT1 {
		chk.HasResult(t, c.Results(t), wantResult, chk.IgnoreOperationID())
	}
	for _, wantResult := range wantOperationResultsDUT2 {
		chk.HasResult(t, c2.Results(t), wantResult, chk.IgnoreOperationID())
	}

	policy1Pfx := "203.0.113.0/30"
	gnmi.Replace(t, dut, ocpath.Root().BgpGuePolicy(policy1Pfx).Config(), &oc.BgpGuePolicy{
		DstPort: ygot.Uint16(42),
		Prefix:  ygot.String(policy1Pfx),
		SrcIp:   ygot.String("42.42.42.42"),
	})

	policy2Pfx := "203.0.113.0/29"
	gnmi.Replace(t, dut, ocpath.Root().BgpGuePolicy(policy2Pfx).Config(), &oc.BgpGuePolicy{
		DstPort: ygot.Uint16(84),
		Prefix:  ygot.String(policy2Pfx),
		SrcIp:   ygot.String("84.84.84.84"),
	})

	// Start a new BGP session that should exchange the necessary gRIBI
	// route that recursively resolves and thus enables traffic to flow.
	dutConf := bgpWithNbr(dutAS, dutPort3.IPv4, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dut2AS),
		NeighborAddress: ygot.String(dut2Port2.IPv4),
	})
	dut2Conf := bgpWithNbr(dut2AS, dut2Port2.IPv4, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dutAS),
		NeighborAddress: ygot.String(dutPort3.IPv4),
	})
	gnmi.Replace(t, dut, bgpPath.Config(), dutConf)
	gnmi.Replace(t, dut2, bgpPath.Config(), dut2Conf)

	nbrPath := bgpPath.Neighbor(dut2Port2.IPv4)
	gnmi.Await(t, dut, nbrPath.SessionState().State(), 60*time.Second, oc.Bgp_Neighbor_SessionState_ESTABLISHED)

	testTrafficAndPrintCapture := func(title string, startingIP string) {
		fmt.Println(title, startingIP)
		otg := ate.OTG()
		otg.StartCapture(t, atePort2.Name)

		if loss := testTraffic(t, ate, ateTop, atePort1, atePort2, startingIP); loss > 1 {
			t.Errorf("Loss: got %g, want <= 1", loss)
		}

		otg.StopCapture(t, atePort2.Name)

		captureBytes := ate.OTG().FetchCapture(t, atePort2.Name)

		f, err := os.CreateTemp(".", "pcap")
		if err != nil {
			t.Fatalf("ERROR: Could not create temporary pcap file: %v\n", err)
		}
		defer os.Remove(f.Name())

		if _, err := f.Write(captureBytes); err != nil {
			t.Fatalf("ERROR: Could not write bytes to pcap file: %v\n", err)
		}
		f.Close()

		handleRead, err := pcap.OpenOffline(f.Name())
		if err != nil {
			t.Fatalf("ERROR: Could not open pcap file %s: %v\n", f.Name(), err)
		}
		ps := gopacket.NewPacketSource(handleRead, layers.LinkTypeEthernet)
		for i := 0; i != 10; i++ {
			pkt, err := ps.NextPacket()

			if err != nil {
				t.Fatalf("error reading next packet: %v", err)
			}
			fmt.Printf("received packet %d: %v\ndump: %s\n", i, pkt, pkt.Dump())
			for i, layer := range pkt.Layers() {
				fmt.Printf("layer %d: %v\n%v, bytes: %d\n", i, layer.LayerType().String(), layer.LayerContents(), len(layer.LayerContents()))
			}
		}
	}

	//testTrafficAndPrintCapture("without policy")

	testTrafficAndPrintCapture("with policy", "198.51.0.0")
	testTrafficAndPrintCapture("with policy", "198.51.2.0")
	testTrafficAndPrintCapture("with policy", "198.51.4.0")

	dut.RawAPIs().GRIBI().Default(t).Flush(context.Background(), &gribipb.FlushRequest{
		NetworkInstance: &gribipb.FlushRequest_All{All: &gribipb.Empty{}},
	})
	dut2.RawAPIs().GRIBI().Default(t).Flush(context.Background(), &gribipb.FlushRequest{
		NetworkInstance: &gribipb.FlushRequest_All{All: &gribipb.Empty{}},
	})
	// TODO: Test that entries are deleted and that there is no more traffic.
}
