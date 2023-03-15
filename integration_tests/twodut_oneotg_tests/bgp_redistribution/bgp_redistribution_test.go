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
	"testing"
	"time"

	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/gribigo/chk"
	"github.com/openconfig/gribigo/client"
	"github.com/openconfig/gribigo/constants"
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/internal/binding"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/gnmi/otg/otgpath"
	"github.com/openconfig/ondatra/otg"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

	gribipb "github.com/openconfig/gribi/v1/proto/service"
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
	// IPv4
	ipv4PrefixLen     = 30
	ateDstNetCIDR     = "198.51.100.0/24"
	ateIndirectNH     = "203.0.113.1"
	ateIndirectNHCIDR = ateIndirectNH + "/24"
	nhIndex           = 1
	nhgIndex          = 42
	// IPv6
	ipv6PrefixLen       = 99
	ateDstNetCIDRv6     = "2003::/49"
	ateIndirectNHv6     = "2002::"
	ateIndirectNHCIDRv6 = ateIndirectNHv6 + "/49"

	dutAS  = 64500
	dut2AS = 64501
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.Get(".."))
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
		IPv6:    "2001::aaaa:bbbb:aa",
		IPv6Len: ipv6PrefixLen,
	}

	atePort1 = Attributes{
		Name:    "port1",
		MAC:     "02:00:01:01:01:01",
		IPv4:    "192.0.2.2",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001::aaaa:bbbb:bb",
		IPv6Len: ipv6PrefixLen,
	}

	dutPort2 = Attributes{
		Desc:    "dutPort2",
		IPv4:    "192.0.2.5",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001::aaab:bbbb:aa",
		IPv6Len: ipv6PrefixLen,
	}

	atePort2 = Attributes{
		Name:    "port2",
		MAC:     "02:00:02:01:01:01",
		IPv4:    "192.0.2.6",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001::aaab:bbbb:bb",
		IPv6Len: ipv6PrefixLen,
	}

	dutPort3 = Attributes{
		Desc:    "dutPort3",
		IPv4:    "192.0.2.9",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001::aaac:bbbb:aa",
		IPv6Len: ipv6PrefixLen,
	}

	atePort3 = Attributes{
		Name:    "port3",
		MAC:     "02:00:03:01:01:01",
		IPv4:    ateIndirectNH,
		IPv4Len: ipv4PrefixLen,
		IPv6:    ateIndirectNHv6,
		IPv6Len: ipv6PrefixLen,
	}

	dut2Port1 = Attributes{
		Desc: "dut2Port1",
		IPv4: "203.0.113.2",
		// Need to make this more specific than ateIndirectNHCIDR so
		// that routes at DUT2 will continue to resolve over this after
		// BGP route exchange.
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2002::cc",
		// Need to make this more specific than ateIndirectNHCIDRv6 so
		// that routes at DUT2 will continue to resolve over this after
		// BGP route exchange.
		IPv6Len: ipv6PrefixLen,
	}

	dut2Port2 = Attributes{
		Desc:    "dut2Port2",
		IPv4:    "192.0.2.10",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001::aaac:bbbb:bb",
		IPv6Len: ipv6PrefixLen,
	}
)

// configureOTG configures port1 and port2 on the ATE.
//
// TODO: Add helpers, or see if featureprofiles already has them, to reduce
// replication.
func configureOTG(t *testing.T, otg *otg.OTG) gosnappi.Config {
	t.Helper()
	top := otg.NewConfig(t)

	top.Ports().Add().SetName(atePort1.Name)
	i1 := top.Devices().Add().SetName(atePort1.Name)
	eth1 := i1.Ethernets().Add().SetName(atePort1.Name + ".Eth").
		SetPortName(i1.Name()).SetMac(atePort1.MAC)
	eth1.Ipv4Addresses().Add().SetName(i1.Name() + ".IPv4").
		SetAddress(atePort1.IPv4).SetGateway(dutPort1.IPv4).
		SetPrefix(int32(atePort1.IPv4Len))
	eth1.Ipv6Addresses().Add().SetName(i1.Name() + ".IPv6").
		SetAddress(atePort1.IPv6).SetGateway(dutPort1.IPv6).
		SetPrefix(int32(atePort1.IPv6Len))

	top.Ports().Add().SetName(atePort2.Name)
	i2 := top.Devices().Add().SetName(atePort2.Name)
	eth2 := i2.Ethernets().Add().SetName(atePort2.Name + ".Eth").
		SetPortName(i2.Name()).SetMac(atePort2.MAC)
	eth2.Ipv4Addresses().Add().SetName(i2.Name() + ".IPv4").
		SetAddress(atePort2.IPv4).SetGateway(dutPort2.IPv4).
		SetPrefix(int32(atePort2.IPv4Len))
	eth2.Ipv6Addresses().Add().SetName(i2.Name() + ".IPv6").
		SetAddress(atePort2.IPv6).SetGateway(dutPort2.IPv6).
		SetPrefix(int32(atePort2.IPv6Len))

	top.Ports().Add().SetName(atePort3.Name)
	i3 := top.Devices().Add().SetName(atePort3.Name)
	eth3 := i3.Ethernets().Add().SetName(atePort3.Name + ".Eth").
		SetPortName(i3.Name()).SetMac(atePort3.MAC)
	eth3.Ipv4Addresses().Add().SetName(i3.Name() + ".IPv4").
		SetAddress(atePort3.IPv4).SetGateway(dut2Port1.IPv4).
		SetPrefix(int32(atePort3.IPv4Len))
	eth3.Ipv6Addresses().Add().SetName(i3.Name() + ".IPv6").
		SetAddress(atePort3.IPv6).SetGateway(dut2Port1.IPv6).
		SetPrefix(int32(atePort3.IPv6Len))
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
	s4a.PrefixLength = ygot.Uint8(ipv4PrefixLen)
	s6 := s.GetOrCreateIpv6()
	s6a := s6.GetOrCreateAddress(a.IPv6)
	s6a.PrefixLength = ygot.Uint8(ipv6PrefixLen)

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
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv6().Address(dutPort1.IPv6).Ip().State(), time.Minute, dutPort1.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv6().Address(dutPort2.IPv6).Ip().State(), time.Minute, dutPort2.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port3").Name()).Subinterface(0).Ipv6().Address(dutPort3.IPv6).Ip().State(), time.Minute, dutPort3.IPv6)
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
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv6().Address(dut2Port1.IPv6).Ip().State(), time.Minute, dut2Port1.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv6().Address(dut2Port2.IPv6).Ip().State(), time.Minute, dut2Port2.IPv6)
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
func testTraffic(t *testing.T, otg *otg.OTG, srcEndPoint, dstEndPoint Attributes, startAddress string, dur time.Duration) float32 {
	waitOTGARPEntry(t)
	top := otg.FetchConfig(t)
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
	v4.Dst().Increment().SetStart(startAddress).SetCount(250)
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
func testTrafficv6(t *testing.T, otg *otg.OTG, srcEndPoint, dstEndPoint Attributes, startAddress string, dur time.Duration) float32 {
	waitOTGARPEntry(t)
	top := otg.FetchConfig(t)
	top.Flows().Clear().Items()
	flowipv6 := top.Flows().Add().SetName("Flow2")
	flowipv6.Metrics().SetEnable(true)
	flowipv6.TxRx().Device().
		SetTxNames([]string{srcEndPoint.Name + ".IPv6"}).
		SetRxNames([]string{dstEndPoint.Name + ".IPv6"})
	flowipv6.Duration().SetChoice("continuous")
	flowipv6.Packet().Add().Ethernet()
	v6 := flowipv6.Packet().Add().Ipv6()
	v6.Src().SetValue(srcEndPoint.IPv6)
	v6.Dst().Increment().SetStart(startAddress).SetCount(250)
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

func bgpWithNbr(as uint32, routerID string, nbrs ...*oc.NetworkInstance_Protocol_Bgp_Neighbor) *oc.NetworkInstance_Protocol_Bgp {
	bgp := &oc.NetworkInstance_Protocol_Bgp{}
	bgp.GetOrCreateGlobal().As = ygot.Uint32(as)
	if routerID != "" {
		bgp.Global.RouterId = ygot.String(routerID)
	}
	for _, nbr := range nbrs {
		bgp.AppendNeighbor(nbr)
	}
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

func installStaticRoute(t *testing.T, dut *ondatra.DUTDevice, route *oc.NetworkInstance_Protocol_Static) {
	staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).
		Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	gnmi.Replace(t, dut, staticp.Static(*route.Prefix).Config(), route)
	gnmi.Await(t, dut, staticp.Static(*route.Prefix).State(), 30*time.Second, route)
}

func TestBGPRouteAdvertisement(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	configureDUT1(t, dut)
	dut2 := ondatra.DUT(t, "dut2")
	configureDUT2(t, dut2)

	ate := ondatra.ATE(t, "ate")
	otg := ate.OTG()
	otgConfig := configureOTG(t, otg)
	otg.PushConfig(t, otgConfig)
	otg.StartProtocols(t)

	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()

	// Remove any existing BGP config
	gnmi.Delete(t, dut, bgpPath.Config())
	gnmi.Delete(t, dut2, bgpPath.Config())

	dutEntries := []fluent.GRIBIEntry{
		// Add an IPv4Entry for 203.0.113.1/24 pointing to 192.0.2.6.
		fluent.NextHopEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
			WithIndex(nhIndex).WithIPAddress(atePort2.IPv4),
		fluent.NextHopGroupEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
			WithID(nhgIndex).AddNextHop(nhIndex, 1),
		fluent.IPv4Entry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
			WithPrefix(ateIndirectNHCIDR).WithNextHopGroup(nhgIndex),
	}
	c := configureGRIBIEntry(t, dut, dutEntries)

	dut2Entries := []fluent.GRIBIEntry{
		// Add an IPv4Entry for 198.51.100.0/24 pointing to 203.0.113.1
		fluent.NextHopEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
			WithIndex(nhIndex).WithIPAddress(ateIndirectNH),
		fluent.NextHopGroupEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
			WithID(nhgIndex).AddNextHop(nhIndex, 1),
		fluent.IPv4Entry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
			WithPrefix(ateDstNetCIDR).WithNextHopGroup(nhgIndex),
	}
	c2 := configureGRIBIEntry(t, dut2, dut2Entries)

	wantOperationResults := []*client.OpResult{
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
			WithIPv4Operation(ateIndirectNHCIDR).
			WithProgrammingResult(fluent.InstalledInFIB).
			WithOperationType(constants.Add).
			AsResult(),
	}

	wantOperationResults2 := []*client.OpResult{
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
	}

	for _, wantResult := range wantOperationResults {
		chk.HasResult(t, c.Results(t), wantResult, chk.IgnoreOperationID())
	}
	for _, wantResult := range wantOperationResults2 {
		chk.HasResult(t, c2.Results(t), wantResult, chk.IgnoreOperationID())
	}

	// Install static route 2002::/49 pointing to 2001::aaab:bbbb:bb.
	installStaticRoute(t, dut, &oc.NetworkInstance_Protocol_Static{
		Prefix: ygot.String(ateIndirectNHCIDRv6),
		NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
			"single": {
				Index:   ygot.String("single"),
				NextHop: oc.UnionString(atePort2.IPv6),
				Recurse: ygot.Bool(true),
			},
		},
	})

	// Install static route 2003::/49 pointing to 2002::.
	installStaticRoute(t, dut2, &oc.NetworkInstance_Protocol_Static{
		Prefix: ygot.String(ateDstNetCIDRv6),
		NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
			"single": {
				Index:   ygot.String("single"),
				NextHop: oc.UnionString(ateIndirectNHv6),
				Recurse: ygot.Bool(true),
			},
		},
	})

	loss := testTraffic(t, otg, atePort1, atePort2, "198.51.100.0", 3*time.Second)
	if loss != 100 {
		t.Errorf("Loss: got %g, want 100", loss)
	}
	loss = testTrafficv6(t, otg, atePort1, atePort2, "2003::", 3*time.Second)
	if loss != 100 {
		t.Errorf("Loss: got %g, want 100", loss)
	}

	// Start a new BGP session that should exchange the necessary gRIBI
	// route that recursively resolves and thus enables traffic to flow.
	dutConf := bgpWithNbr(dutAS, dutPort3.IPv4, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dut2AS),
		NeighborAddress: ygot.String(dut2Port2.IPv4),
	}, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dut2AS),
		NeighborAddress: ygot.String(dut2Port2.IPv6),
	})
	dut2Conf := bgpWithNbr(dut2AS, dut2Port2.IPv4, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dutAS),
		NeighborAddress: ygot.String(dutPort3.IPv4),
	}, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dutAS),
		NeighborAddress: ygot.String(dutPort3.IPv6),
	})
	gnmi.Replace(t, dut, bgpPath.Config(), dutConf)
	gnmi.Replace(t, dut2, bgpPath.Config(), dut2Conf)

	nbrPath := bgpPath.Neighbor(dut2Port2.IPv4)
	gnmi.Await(t, dut, nbrPath.SessionState().State(), 60*time.Second, oc.Bgp_Neighbor_SessionState_ESTABLISHED)
	nbrPathv6 := bgpPath.Neighbor(dut2Port2.IPv6)
	gnmi.Await(t, dut, nbrPathv6.SessionState().State(), 60*time.Second, oc.Bgp_Neighbor_SessionState_ESTABLISHED)

	loss = testTraffic(t, otg, atePort1, atePort2, "198.51.100.0", 15*time.Second)
	if loss > 1 {
		t.Errorf("Loss: got %g, want <= 1", loss)
	}
	loss = testTrafficv6(t, otg, atePort1, atePort2, "2003::", 15*time.Second)
	if loss > 1 {
		t.Errorf("Loss: got %g, want <= 1", loss)
	}

	dut.RawAPIs().GRIBI().Default(t).Flush(context.Background(), &gribipb.FlushRequest{
		NetworkInstance: &gribipb.FlushRequest_All{All: &gribipb.Empty{}},
	})
	dut2.RawAPIs().GRIBI().Default(t).Flush(context.Background(), &gribipb.FlushRequest{
		NetworkInstance: &gribipb.FlushRequest_All{All: &gribipb.Empty{}},
	})
	// TODO: Test that entries are deleted and that there is no more traffic.
}
