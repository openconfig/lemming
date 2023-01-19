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
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/gribigo/chk"
	"github.com/openconfig/gribigo/client"
	"github.com/openconfig/gribigo/constants"
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/integration_tests/binding"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	otgtelemetry "github.com/openconfig/ondatra/gnmi/otg"
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
	dut2AS = 64500
	ateAS  = 64502
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

// configureOTG configures ports and other configurations on the OTG device.
func configureOTG(t *testing.T, otg *otg.OTG) gosnappi.Config {
	t.Helper()
	config := otg.NewConfig(t)

	// Configure port1
	config.Ports().Add().SetName(atePort1.Name)
	i1 := config.Devices().Add().SetName(atePort1.Name)
	eth1 := i1.Ethernets().Add().SetName(atePort1.Name + ".Eth").
		SetPortName(i1.Name()).SetMac(atePort1.MAC)
	eth1.Ipv4Addresses().Add().SetName(i1.Name() + ".IPv4").
		SetAddress(atePort1.IPv4).SetGateway(dutPort1.IPv4).
		SetPrefix(int32(atePort1.IPv4Len))
	// Configure capture format.
	config.Captures().Add().SetName("ca1").SetPortNames([]string{atePort1.Name}).SetFormat(gosnappi.CaptureFormat.PCAP)

	// Configure port2
	config.Ports().Add().SetName(atePort2.Name)
	i2 := config.Devices().Add().SetName(atePort2.Name)
	eth2 := i2.Ethernets().Add().SetName(atePort2.Name + ".Eth").
		SetPortName(i2.Name()).SetMac(atePort2.MAC)
	eth2.Ipv4Addresses().Add().SetName(i2.Name() + ".IPv4").
		SetAddress(atePort2.IPv4).SetGateway(dutPort2.IPv4).
		SetPrefix(int32(atePort2.IPv4Len))
	// Configure capture format.
	config.Captures().Add().SetName("ca2").SetPortNames([]string{atePort2.Name}).SetFormat(gosnappi.CaptureFormat.PCAP)

	// Configure port3
	config.Ports().Add().SetName(atePort3.Name)
	i3 := config.Devices().Add().SetName(atePort3.Name)
	eth3 := i3.Ethernets().Add().SetName(atePort3.Name + ".Eth").
		SetPortName(i3.Name()).SetMac(atePort3.MAC)
	eth3.Ipv4Addresses().Add().SetName(i3.Name() + ".IPv4").
		SetAddress(atePort3.IPv4).SetGateway(dut2Port1.IPv4).
		SetPrefix(int32(atePort3.IPv4Len))

	// Configure BGP neighbour
	// This causes the route ateIndirectNHCIDR -> atePort2's IP to be
	// exchanged from OTG to DUT on DUT port 1.
	bgp4ObjectMap := make(map[string]gosnappi.BgpV4Peer)
	ipv4ObjectMap := make(map[string]gosnappi.DeviceIpv4)

	ateName := atePort2.Name
	ateNeighbor := dutPort1
	devName := ateName + ".dev"
	bgpNexthop := atePort2.IPv4

	bgp := i1.Bgp().SetRouterId(atePort2.IPv4)

	ipv4 := eth1.Ipv4Addresses().Add().SetName(devName + ".IPv4").SetAddress(atePort1.IPv4).SetGateway(ateNeighbor.IPv4).SetPrefix(ipv4PrefixLen)
	bgp4Name := devName + ".BGP4.peer"
	bgp4Peer := bgp.Ipv4Interfaces().Add().SetIpv4Name(ipv4.Name()).Peers().Add().SetName(bgp4Name).SetPeerAddress(ipv4.Gateway()).SetAsNumber(int32(ateAS)).SetAsType(gosnappi.BgpV4PeerAsType.EBGP)

	bgp4Peer.Capability().SetIpv4UnicastAddPath(true).SetIpv6UnicastAddPath(true)
	bgp4Peer.LearnedInformationFilter().SetUnicastIpv4Prefix(true).SetUnicastIpv6Prefix(true)

	bgp4ObjectMap[bgp4Name] = bgp4Peer
	ipv4ObjectMap[devName+".IPv4"] = ipv4

	bgpName := ateName + ".dev.BGP4.peer"
	bgpPeer := bgp4ObjectMap[bgpName]
	firstAdvAddr := strings.Split(ateIndirectNHCIDR, "/")[0]
	firstAdvPrefix, _ := strconv.Atoi(strings.Split(ateIndirectNHCIDR, "/")[1])
	bgp4PeerRoutes := bgpPeer.V4Routes().Add().SetName(bgpName + ".rr4").SetNextHopIpv4Address(bgpNexthop).SetNextHopAddressType(gosnappi.BgpV4RouteRangeNextHopAddressType.IPV4).SetNextHopMode(gosnappi.BgpV4RouteRangeNextHopMode.MANUAL)
	bgp4PeerRoutes.Addresses().Add().SetAddress(firstAdvAddr).SetPrefix(int32(firstAdvPrefix)).SetCount(1)
	bgp4PeerRoutes.AddPath().SetPathId(1)

	return config
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

// NewOCInterface returns a new *oc.Interface configured with these attributes.
func (a *Attributes) NewOCInterface(name string) *oc.Interface {
	return a.ConfigOCInterface(&oc.Interface{Name: ygot.String(name)})
}

// ConfigOCInterface configures an OpenConfig interface with these attributes.
func (a *Attributes) ConfigOCInterface(intf *oc.Interface) *oc.Interface {
	if a.Desc != "" {
		intf.Description = ygot.String(a.Desc)
	}
	intf.Type = oc.IETFInterfaces_InterfaceType_ethernetCsmacd
	if a.MTU > 0 {
		intf.Mtu = ygot.Uint16(a.MTU + 14)
	}
	e := intf.GetOrCreateEthernet()
	if a.MAC != "" {
		e.MacAddress = ygot.String(a.MAC)
	}

	s := intf.GetOrCreateSubinterface(0)
	if a.IPv4 != "" {
		s4 := s.GetOrCreateIpv4()
		if a.MTU > 0 {
			s4.Mtu = ygot.Uint16(a.MTU)
		}
		a4 := s4.GetOrCreateAddress(a.IPv4)
		if a.IPv4Len > 0 {
			a4.PrefixLength = ygot.Uint8(a.IPv4Len)
		}
	}

	if a.IPv6 != "" {
		s6 := s.GetOrCreateIpv6()
		if a.MTU > 0 {
			s6.Mtu = ygot.Uint32(uint32(a.MTU))
		}
		a6 := s6.GetOrCreateAddress(a.IPv6)
		if a.IPv6Len > 0 {
			a6.PrefixLength = ygot.Uint8(a.IPv6Len)
		}
	}
	return intf
}

// configureDUT1 configures ports on DUT1.
func configureDUT1(t *testing.T, dut *ondatra.DUTDevice) {
	p1 := dut.Port(t, "port1")
	gnmi.Replace(t, dut, ocpath.Root().Interface(p1.Name()).Config(), dutPort1.NewOCInterface(p1.Name()))

	p2 := dut.Port(t, "port2")
	gnmi.Replace(t, dut, ocpath.Root().Interface(p2.Name()).Config(), dutPort2.NewOCInterface(p2.Name()))

	p3 := dut.Port(t, "port3")
	gnmi.Replace(t, dut, ocpath.Root().Interface(p3.Name()).Config(), dutPort3.NewOCInterface(p3.Name()))

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv4().Address(dutPort1.IPv4).Ip().State(), time.Minute, dutPort1.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv4().Address(dutPort2.IPv4).Ip().State(), time.Minute, dutPort2.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port3").Name()).Subinterface(0).Ipv4().Address(dutPort3.IPv4).Ip().State(), time.Minute, dutPort3.IPv4)

	// Start a new BGP session that should exchange the necessary gRIBI
	// route that recursively resolves and thus enables traffic to flow.
	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()

	// Remove any existing BGP config
	gnmi.Delete(t, dut, bgpPath.Config())

	dutConf := bgpWithNbr(dutAS, dutPort3.IPv4, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dut2AS),
		NeighborAddress: ygot.String(dut2Port2.IPv4),
	}, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(ateAS),
		NeighborAddress: ygot.String(atePort1.IPv4),
	})
	gnmi.Replace(t, dut, bgpPath.Config(), dutConf)
}

// configureDUT2 configures ports on DUT2.
func configureDUT2(t *testing.T, dut *ondatra.DUTDevice) {
	p1 := dut.Port(t, "port1")
	gnmi.Replace(t, dut, ocpath.Root().Interface(p1.Name()).Config(), dut2Port1.NewOCInterface(p1.Name()))

	p2 := dut.Port(t, "port2")
	gnmi.Replace(t, dut, ocpath.Root().Interface(p2.Name()).Config(), dut2Port2.NewOCInterface(p2.Name()))

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv4().Address(dut2Port1.IPv4).Ip().State(), time.Minute, dut2Port1.IPv4)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv4().Address(dut2Port2.IPv4).Ip().State(), time.Minute, dut2Port2.IPv4)

	// Start a new BGP session that should exchange the necessary gRIBI
	// route that recursively resolves and thus enables traffic to flow.
	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()

	// Remove any existing BGP config
	gnmi.Delete(t, dut, bgpPath.Config())

	dut2Conf := bgpWithNbr(dut2AS, dut2Port2.IPv4, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dutAS),
		NeighborAddress: ygot.String(dutPort3.IPv4),
	})
	gnmi.Replace(t, dut, bgpPath.Config(), dut2Conf)

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
func testTraffic(t *testing.T, otg *otg.OTG, srcEndPoint, dstEndPoint Attributes, startingIP string) float32 {
	waitOTGARPEntry(t)
	otgConfig := otg.FetchConfig(t)
	otgConfig.Flows().Clear().Items()
	flowipv4 := otgConfig.Flows().Add().SetName("Flow")
	flowipv4.Metrics().SetEnable(true)
	flowipv4.TxRx().Device().
		SetTxNames([]string{srcEndPoint.Name + ".IPv4"}).
		SetRxNames([]string{dstEndPoint.Name + ".IPv4"})
	flowipv4.Duration().SetChoice("continuous")
	flowipv4.Packet().Add().Ethernet()
	v4 := flowipv4.Packet().Add().Ipv4()
	v4.Src().SetValue(srcEndPoint.IPv4)
	v4.Dst().Increment().SetStart(startingIP).SetCount(250)
	otg.PushConfig(t, otgConfig)

	otg.StartTraffic(t)
	time.Sleep(15 * time.Second)
	t.Logf("Stop traffic")
	otg.StopTraffic(t)

	time.Sleep(5 * time.Second)

	txPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().OutPkts().State())
	rxPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().InPkts().State())
	lossPct := (txPkts - rxPkts) * 100 / txPkts
	return float32(lossPct)
}

// testTrafficAndEncap checks that traffic can reach from ATE port1 to ATE
// port2 with the correct GUE encap fields (if any).
func testTrafficAndEncap(t *testing.T, otg *otg.OTG, startingIP string, encapFields *EncapFields) {
	t.Helper()
	otg.StartCapture(t, atePort2.Name)

	if loss := testTraffic(t, otg, atePort1, atePort2, startingIP); loss > 1 {
		t.Errorf("Loss: got %g, want <= 1", loss)
	}

	otg.StopCapture(t, atePort2.Name)

	captureBytes := otg.FetchCapture(t, atePort2.Name)

	f, err := os.CreateTemp(".", "pcap")
	if err != nil {
		t.Fatalf("ERROR: Could not create temporary pcap file: %v\n", err)
	}
	defer os.Remove(f.Name())

	if _, err := f.Write(captureBytes); err != nil {
		t.Fatalf("ERROR: Could not write bytes to pcap file: %v\n", err)
	}
	f.Close()

	f, err = os.Open(f.Name())
	if err != nil {
		t.Fatalf("ERROR: Could not open pcap file %s: %v\n", f.Name(), err)
	}
	defer f.Close()

	handleRead, err := pcapgo.NewReader(f)
	if err != nil {
		t.Fatalf("ERROR: Could not create reader on pcap file %s: %v\n", f.Name(), err)
	}
	ps := gopacket.NewPacketSource(handleRead, layers.LinkTypeEthernet)

	for i := 0; i != 10; i++ {
		pkt, err := ps.NextPacket()
		if err != nil {
			t.Fatalf("error reading next packet: %v", err)
		}
		if encapFields == nil {
			wantNL := layers.IPv4{
				Version: 4,
				IHL:     5,
				Length:  46,
				SrcIP:   net.IP{192, 0, 2, 2}, // ATE port 1.
			}
			var gotNL layers.IPv4
			nl := pkt.NetworkLayer()
			if nl == nil {
				t.Errorf("packet doesn't have network layer: %s", pkt.Dump())
				continue
			}
			if err := gotNL.DecodeFromBytes(nl.LayerContents(), gopacket.NilDecodeFeedback); err != nil {
				t.Errorf("cannot decode network layer header: %v", err)
				continue
			}
			if diff := cmp.Diff(wantNL, gotNL, cmpopts.IgnoreUnexported(layers.IPv4{}), cmpopts.IgnoreFields(layers.IPv4{}, "BaseLayer", "Checksum", "TTL", "Protocol", "DstIP")); diff != "" {
				t.Errorf("network layer (-want, +got):\n%s", diff)
			}
		} else {
			wantNL := layers.IPv4{
				Version:  4,
				IHL:      5,
				Length:   74,
				Protocol: layers.IPProtocolUDP,
				SrcIP:    encapFields.srcIP,
				DstIP:    encapFields.dstIP,
			}
			var gotNL layers.IPv4
			nl := pkt.NetworkLayer()
			if nl == nil {
				t.Errorf("packet doesn't have network layer: %s", pkt.Dump())
				continue
			}
			if err := gotNL.DecodeFromBytes(nl.LayerContents(), gopacket.NilDecodeFeedback); err != nil {
				t.Errorf("cannot decode network layer header: %v", err)
				continue
			}
			if diff := cmp.Diff(wantNL, gotNL, cmpopts.IgnoreUnexported(layers.IPv4{}), cmpopts.IgnoreFields(layers.IPv4{}, "BaseLayer", "Checksum")); diff != "" {
				t.Errorf("network layer (-want, +got):\n%s", diff)
			}

			wantTL := layers.UDP{
				SrcPort: 0, // TODO(wenbli): Implement and test hashing for srcPort.
				DstPort: layers.UDPPort(encapFields.dstPort),
				Length:  34,
			}
			var gotTL layers.UDP
			tl := pkt.TransportLayer()
			if tl == nil {
				t.Errorf("packet doesn't have transport layer: %s", pkt.Dump())
				continue
			}
			if err := gotTL.DecodeFromBytes(tl.LayerContents(), gopacket.NilDecodeFeedback); err != nil {
				t.Errorf("cannot decode transport layer header: %v", err)
				continue
			}
			if diff := cmp.Diff(wantTL, gotTL, cmpopts.IgnoreUnexported(layers.UDP{}), cmpopts.IgnoreFields(layers.UDP{}, "BaseLayer")); diff != "" {
				t.Errorf("network layer (-want, +got):\n%s", diff)
			}
			// TODO: Check that lower layers is the original packet.
		}
	}
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

func verifyOTGBGPTelemetry(t *testing.T, otg *otg.OTG, c gosnappi.Config, state string) {
	for _, d := range c.Devices().Items() {
		for _, ip := range d.Bgp().Ipv4Interfaces().Items() {
			for _, configPeer := range ip.Peers().Items() {
				nbrPath := gnmi.OTG().BgpPeer(configPeer.Name())
				_, ok := gnmi.Watch(t, otg, nbrPath.SessionState().State(), time.Minute, func(val *ygnmi.Value[otgtelemetry.E_BgpPeer_SessionState]) bool {
					currState, ok := val.Val()
					return ok && currState.String() == state
				}).Await(t)
				if !ok {
					t.Log("BGP reported state", nbrPath.State(), gnmi.Get(t, otg, nbrPath.State()))
					t.Errorf("No BGP neighbor formed for peer %s", configPeer.Name())
				}
			}
		}
		for _, ip := range d.Bgp().Ipv6Interfaces().Items() {
			for _, configPeer := range ip.Peers().Items() {
				nbrPath := gnmi.OTG().BgpPeer(configPeer.Name())
				_, ok := gnmi.Watch(t, otg, nbrPath.SessionState().State(), time.Minute, func(val *ygnmi.Value[otgtelemetry.E_BgpPeer_SessionState]) bool {
					currState, ok := val.Val()
					return ok && currState.String() == state
				}).Await(t)
				if !ok {
					t.Log("BGP reported state", nbrPath.State(), gnmi.Get(t, otg, nbrPath.State()))
					t.Errorf("No BGP neighbor formed for peer %s", configPeer.Name())
				}
			}
		}
	}
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

// EncapFields are the expected encap fields for a GUE-encapped packet.
type EncapFields struct {
	srcIP   net.IP
	dstIP   net.IP
	dstPort uint16
}

func installGRIBIEntries(t *testing.T, dut2 *ondatra.DUTDevice) {
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

	for _, wantResult := range wantOperationResultsDUT2 {
		chk.HasResult(t, c2.Results(t), wantResult, chk.IgnoreOperationID())
	}
}

func TestBGPTriggeredGUE(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	configureDUT1(t, dut)
	dut2 := ondatra.DUT(t, "dut2")
	configureDUT2(t, dut2)

	ate := ondatra.ATE(t, "ate")
	otg := ate.OTG()
	otgConfig := configureOTG(t, otg)
	t.Logf("Pushing config to ATE and starting protocols...")
	otg.PushConfig(t, otgConfig)
	otg.StartProtocols(t)

	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()
	nbrPath := bgpPath.Neighbor(dut2Port2.IPv4)
	gnmi.Await(t, dut, nbrPath.SessionState().State(), 120*time.Second, oc.Bgp_Neighbor_SessionState_ESTABLISHED)

	t.Logf("Verify DUT's DUT-OTG BGP sessions up")
	gnmi.Await(t, dut, bgpPath.Neighbor(atePort1.IPv4).SessionState().State(), 120*time.Second, oc.Bgp_Neighbor_SessionState_ESTABLISHED)
	t.Logf("Verify OTG's DUT-OTG BGP sessions up")
	verifyOTGBGPTelemetry(t, otg, otgConfig, "ESTABLISHED")

	// Install entries to be propagated to DUT1.
	installGRIBIEntries(t, dut2)

	tests := []struct {
		desc         string
		gnmiOp       func()
		encapFields1 *EncapFields
		encapFields2 *EncapFields
		encapFields3 *EncapFields
		skip         bool
	}{{
		desc:   "without policy",
		gnmiOp: func() {},
	}, {
		desc: "with single policy",
		gnmiOp: func() {
			policy2Pfx := "203.0.113.0/29"
			gnmi.Replace(t, dut, ocpath.Root().BgpGueIpv4Policy(policy2Pfx).Config(), &oc.BgpGueIpv4Policy{
				// TODO(wenbli): Support IPv4-mapped IPv6 traffic that would use a different dstPort.
				DstPortIpv4: ygot.Uint16(84),
				DstPortIpv6: ygot.Uint16(168),
				Prefix:      ygot.String(policy2Pfx),
				SrcIp:       ygot.String("84.84.84.84"),
			})
		},
		encapFields1: &EncapFields{
			srcIP:   net.IP{84, 84, 84, 84},
			dstIP:   net.IP{203, 0, 113, 1},
			dstPort: 84,
		},
		encapFields2: &EncapFields{
			srcIP:   net.IP{84, 84, 84, 84},
			dstIP:   net.IP{203, 0, 113, 5},
			dstPort: 84,
		},
	}, {
		desc: "with two overlapping policies",
		gnmiOp: func() {
			policy1Pfx := "203.0.113.0/30"
			gnmi.Replace(t, dut, ocpath.Root().BgpGueIpv4Policy(policy1Pfx).Config(), &oc.BgpGueIpv4Policy{
				DstPortIpv4: ygot.Uint16(42),
				DstPortIpv6: ygot.Uint16(168),
				Prefix:      ygot.String(policy1Pfx),
				SrcIp:       ygot.String("42.42.42.42"),
			})
		},
		encapFields1: &EncapFields{
			srcIP:   net.IP{42, 42, 42, 42},
			dstIP:   net.IP{203, 0, 113, 1},
			dstPort: 42,
		},
		encapFields2: &EncapFields{
			srcIP:   net.IP{84, 84, 84, 84},
			dstIP:   net.IP{203, 0, 113, 5},
			dstPort: 84,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}

			tests := []struct {
				startingIP  string
				encapFields *EncapFields
			}{{
				startingIP:  "198.51.0.0",
				encapFields: tt.encapFields1,
			}, {
				startingIP:  "198.51.2.0",
				encapFields: tt.encapFields2,
			}, {
				startingIP:  "198.51.4.0",
				encapFields: tt.encapFields3,
			}}

			tt.gnmiOp()
			for _, tt := range tests {
				t.Run(tt.startingIP, func(t *testing.T) {
					testTrafficAndEncap(t, otg, tt.startingIP, tt.encapFields)
				})
			}
		})
	}

	dut2.RawAPIs().GRIBI().Default(t).Flush(context.Background(), &gribipb.FlushRequest{
		NetworkInstance: &gribipb.FlushRequest_All{All: &gribipb.Empty{}},
	})
	// TODO: Test that entries are deleted and that there is no more traffic.
}
