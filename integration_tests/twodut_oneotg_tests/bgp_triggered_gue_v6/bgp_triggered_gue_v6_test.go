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
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/internal/binding"
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
//   - ate:port1 -> dut:port1 subnet 2001::aaaa:bbbb:cccc/99
//   - ate:port2 -> dut:port2 subnet 2001::aaab:bbbb:cccc/99
//   - ate:port3 -> dut2:port1 subnet 2002::/49
//   - dut:port3 -> dut2:port2 subnet 2001::aaac:bbbb:aa/99
const (
	ipv4PrefixLen          = 30
	ipv6PrefixLen          = 99
	ateDstNetCIDR1         = "2003:aaaa::/49"
	ateDstNetCIDR2         = "2003:bbbb::/49"
	ateDstNetCIDR3         = "2003:cccc::/49"
	ateIndirectNH1         = "2002:0:0:0:10::"
	ateIndirectNH2         = "2002:0:0:10::"
	ateIndirectNH3         = "2002:0:10::"
	ateIndirectNHCIDR      = "2002::/32"
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
		IPv6:    "2001::aaaa:bbbb:aa",
		IPv6Len: ipv6PrefixLen,
	}

	atePort1 = Attributes{
		Name:    "port1",
		MAC:     "02:00:01:01:01:01",
		IPv6:    "2001::aaaa:bbbb:bb",
		IPv6Len: ipv6PrefixLen,
	}

	dutPort2 = Attributes{
		Desc:    "dutPort2",
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
		IPv6:    ateIndirectNH3,
		IPv6Len: 32,
	}

	dut2Port1 = Attributes{
		Desc: "dut2Port1",
		IPv6: "2002::cc",
		// Make sure the ATE indirect prefixes that are to be exchanged
		// to DUT1 are actually resolvable at DUT2.
		IPv6Len: 32,
	}

	dut2Port2 = Attributes{
		Desc:    "dut2Port2",
		IPv4:    "192.0.2.10",
		IPv4Len: ipv4PrefixLen,
		IPv6:    "2001::aaac:bbbb:bb",
		IPv6Len: ipv6PrefixLen,
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
	eth1.Ipv6Addresses().Add().SetName(i1.Name() + ".IPv6").
		SetAddress(atePort1.IPv6).SetGateway(dutPort1.IPv6).
		SetPrefix(int32(atePort1.IPv6Len))
	// Configure capture format.
	config.Captures().Add().SetName("ca1").SetPortNames([]string{atePort1.Name}).SetFormat(gosnappi.CaptureFormat.PCAP)

	// Configure port2
	config.Ports().Add().SetName(atePort2.Name)
	i2 := config.Devices().Add().SetName(atePort2.Name)
	eth2 := i2.Ethernets().Add().SetName(atePort2.Name + ".Eth").
		SetPortName(i2.Name()).SetMac(atePort2.MAC)
	eth2.Ipv6Addresses().Add().SetName(i2.Name() + ".IPv6").
		SetAddress(atePort2.IPv6).SetGateway(dutPort2.IPv6).
		SetPrefix(int32(atePort2.IPv6Len))
	// Configure capture format.
	config.Captures().Add().SetName("ca2").SetPortNames([]string{atePort2.Name}).SetFormat(gosnappi.CaptureFormat.PCAP)

	// Configure port3
	config.Ports().Add().SetName(atePort3.Name)
	i3 := config.Devices().Add().SetName(atePort3.Name)
	eth3 := i3.Ethernets().Add().SetName(atePort3.Name + ".Eth").
		SetPortName(i3.Name()).SetMac(atePort3.MAC)
	eth3.Ipv6Addresses().Add().SetName(i3.Name() + ".IPv6").
		SetAddress(atePort3.IPv6).SetGateway(dut2Port1.IPv6).
		SetPrefix(int32(atePort3.IPv6Len))

	// Configure BGP neighbour
	// This causes the route ateIndirectNHCIDR -> atePort2's IP to be
	// exchanged from OTG to DUT on DUT port 1.
	bgp6ObjectMap := make(map[string]gosnappi.BgpV6Peer)
	ipv6ObjectMap := make(map[string]gosnappi.DeviceIpv6)

	ateName := atePort2.Name
	ateNeighbor := dutPort1
	devName := ateName + ".dev"
	bgpNexthop := atePort2.IPv6

	bgp := i1.Bgp().SetRouterId(atePort2.IPv4)

	ipv6 := eth1.Ipv6Addresses().Add().SetName(devName + ".IPv6").SetAddress(atePort1.IPv6).SetGateway(ateNeighbor.IPv6).SetPrefix(ipv6PrefixLen)
	bgp6Name := devName + ".BGP6.peer"
	bgp6Peer := bgp.Ipv6Interfaces().Add().SetIpv6Name(ipv6.Name()).Peers().Add().SetName(bgp6Name).SetPeerAddress(ipv6.Gateway()).SetAsNumber(int32(ateAS)).SetAsType(gosnappi.BgpV6PeerAsType.EBGP)

	bgp6Peer.Capability().SetIpv6UnicastAddPath(true).SetIpv6UnicastAddPath(true)
	bgp6Peer.LearnedInformationFilter().SetUnicastIpv6Prefix(true).SetUnicastIpv6Prefix(true)

	bgp6ObjectMap[bgp6Name] = bgp6Peer
	ipv6ObjectMap[devName+".IPv6"] = ipv6

	bgpName := ateName + ".dev.BGP6.peer"
	bgpPeer := bgp6ObjectMap[bgpName]
	firstAdvAddr := strings.Split(ateIndirectNHCIDR, "/")[0]
	firstAdvPrefix, _ := strconv.Atoi(strings.Split(ateIndirectNHCIDR, "/")[1])
	bgp6PeerRoutes := bgpPeer.V6Routes().Add().SetName(bgpName + ".rr6").SetNextHopIpv6Address(bgpNexthop).SetNextHopAddressType(gosnappi.BgpV6RouteRangeNextHopAddressType.IPV6).SetNextHopMode(gosnappi.BgpV6RouteRangeNextHopMode.MANUAL)
	bgp6PeerRoutes.Addresses().Add().SetAddress(firstAdvAddr).SetPrefix(int32(firstAdvPrefix)).SetCount(1)
	bgp6PeerRoutes.AddPath().SetPathId(1)

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
	s6 := s.GetOrCreateIpv6()
	s6a := s6.GetOrCreateAddress(a.IPv6)
	s6a.PrefixLength = ygot.Uint8(a.IPv6Len)

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

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv6().Address(dutPort1.IPv6).Ip().State(), time.Minute, dutPort1.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv6().Address(dutPort2.IPv6).Ip().State(), time.Minute, dutPort2.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port3").Name()).Subinterface(0).Ipv6().Address(dutPort3.IPv6).Ip().State(), time.Minute, dutPort3.IPv6)

	// Start a new BGP session that should exchange the necessary gRIBI
	// route that recursively resolves and thus enables traffic to flow.
	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()

	// Remove any existing BGP config
	gnmi.Delete(t, dut, bgpPath.Config())

	dutConf := bgpWithNbr(dutAS, dutPort3.IPv4, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dut2AS),
		NeighborAddress: ygot.String(dut2Port2.IPv6),
	}, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(ateAS),
		NeighborAddress: ygot.String(atePort1.IPv6),
	})
	gnmi.Replace(t, dut, bgpPath.Config(), dutConf)
}

// configureDUT2 configures ports on DUT2.
func configureDUT2(t *testing.T, dut *ondatra.DUTDevice) {
	p1 := dut.Port(t, "port1")
	gnmi.Replace(t, dut, ocpath.Root().Interface(p1.Name()).Config(), dut2Port1.NewOCInterface(p1.Name()))

	p2 := dut.Port(t, "port2")
	gnmi.Replace(t, dut, ocpath.Root().Interface(p2.Name()).Config(), dut2Port2.NewOCInterface(p2.Name()))

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv6().Address(dut2Port1.IPv6).Ip().State(), time.Minute, dut2Port1.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv6().Address(dut2Port2.IPv6).Ip().State(), time.Minute, dut2Port2.IPv6)

	// Start a new BGP session that should exchange the necessary gRIBI
	// route that recursively resolves and thus enables traffic to flow.
	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()

	// Remove any existing BGP config
	gnmi.Delete(t, dut, bgpPath.Config())

	dut2Conf := bgpWithNbr(dut2AS, dut2Port2.IPv4, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dutAS),
		NeighborAddress: ygot.String(dutPort3.IPv6),
	})
	gnmi.Replace(t, dut, bgpPath.Config(), dut2Conf)

}

func waitOTGARPEntry(t *testing.T) {
	ate := ondatra.ATE(t, "ate")

	val, ok := gnmi.WatchAll(t, ate.OTG(), otgpath.Root().InterfaceAny().Ipv6NeighborAny().LinkLayerAddress().State(), time.Minute, func(v *ygnmi.Value[string]) bool {
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
	flowipv6 := otgConfig.Flows().Add().SetName("Flow")
	flowipv6.Metrics().SetEnable(true)
	flowipv6.TxRx().Device().
		SetTxNames([]string{srcEndPoint.Name + ".IPv6"}).
		SetRxNames([]string{dstEndPoint.Name + ".IPv6"})
	flowipv6.Duration().SetChoice("continuous")
	flowipv6.Packet().Add().Ethernet()
	v6 := flowipv6.Packet().Add().Ipv6()
	v6.Src().SetValue(srcEndPoint.IPv6)
	v6.Dst().Increment().SetStart(startingIP).SetCount(250)
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
			wantNL := layers.IPv6{
				Version: 6,
				Length:  24,
				SrcIP:   net.IP{0x20, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0xaa, 0xaa, 0xbb, 0xbb, 0, 0xbb}, // ATE port 1.
			}
			var gotNL layers.IPv6
			nl := pkt.NetworkLayer()
			if nl == nil {
				t.Errorf("packet doesn't have network layer: %s", pkt.Dump())
				continue
			}
			if err := gotNL.DecodeFromBytes(nl.LayerContents(), gopacket.NilDecodeFeedback); err != nil {
				t.Errorf("cannot decode network layer header: %v", err)
				continue
			}
			// TODO(wenbli): What should NextHeader be?
			if diff := cmp.Diff(wantNL, gotNL, cmpopts.IgnoreUnexported(layers.IPv6{}), cmpopts.IgnoreFields(layers.IPv6{}, "BaseLayer", "TrafficClass", "FlowLabel", "HopLimit", "DstIP", "NextHeader")); diff != "" {
				t.Errorf("network layer (-want, +got):\n%s", diff)
			}
		} else {
			wantNL := layers.IPv6{
				Version:    6,
				Length:     72,
				NextHeader: layers.IPProtocolUDP,
				SrcIP:      encapFields.srcIP,
				DstIP:      encapFields.dstIP,
			}
			var gotNL layers.IPv6
			nl := pkt.NetworkLayer()
			if nl == nil {
				t.Errorf("packet doesn't have network layer: %s", pkt.Dump())
				continue
			}
			if err := gotNL.DecodeFromBytes(nl.LayerContents(), gopacket.NilDecodeFeedback); err != nil {
				t.Errorf("cannot decode network layer header: %v", err)
				continue
			}
			if diff := cmp.Diff(wantNL, gotNL, cmpopts.IgnoreUnexported(layers.IPv6{}), cmpopts.IgnoreFields(layers.IPv6{}, "BaseLayer", "TrafficClass", "FlowLabel", "HopLimit", "DstIP")); diff != "" {
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
				t.Errorf("transport layer (-want, +got):\n%s", diff)
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

func installStaticRoute(t *testing.T, dut *ondatra.DUTDevice, route *oc.NetworkInstance_Protocol_Static) {
	staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).
		Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	gnmi.Replace(t, dut, staticp.Static(*route.Prefix).Config(), route)
	gnmi.Await(t, dut, staticp.Static(*route.Prefix).State(), 30*time.Second, route)
}

func TestBGPTriggeredGUEV6(t *testing.T) {
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
	t.Logf("Verify DUT's DUT-DUT BGP session up")
	gnmi.Await(t, dut, bgpPath.Neighbor(dut2Port2.IPv6).SessionState().State(), 120*time.Second, oc.Bgp_Neighbor_SessionState_ESTABLISHED)
	t.Logf("Verify DUT's DUT-OTG BGP session up")
	gnmi.Await(t, dut, bgpPath.Neighbor(atePort1.IPv6).SessionState().State(), 120*time.Second, oc.Bgp_Neighbor_SessionState_ESTABLISHED)
	t.Logf("Verify OTG's OTG-DUT BGP session up")
	verifyOTGBGPTelemetry(t, otg, otgConfig, "ESTABLISHED")

	// Install entries to be propagated to DUT1.
	installStaticRoute(t, dut2, &oc.NetworkInstance_Protocol_Static{
		Prefix: ygot.String(ateDstNetCIDR1),
		NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
			"single": {
				Index:   ygot.String("single"),
				NextHop: oc.UnionString(ateIndirectNH1),
				Recurse: ygot.Bool(true),
			},
		},
	})

	installStaticRoute(t, dut2, &oc.NetworkInstance_Protocol_Static{
		Prefix: ygot.String(ateDstNetCIDR2),
		NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
			"single": {
				Index:   ygot.String("single"),
				NextHop: oc.UnionString(ateIndirectNH2),
				Recurse: ygot.Bool(true),
			},
		},
	})

	installStaticRoute(t, dut2, &oc.NetworkInstance_Protocol_Static{
		Prefix: ygot.String(ateDstNetCIDR3),
		NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
			"single": {
				Index:   ygot.String("single"),
				NextHop: oc.UnionString(ateIndirectNH3),
				Recurse: ygot.Bool(true),
			},
		},
	})

	tests := []struct {
		desc             string
		gnmiOp           func()
		wantEncapFields1 *EncapFields
		wantEncapFields2 *EncapFields
		wantEncapFields3 *EncapFields
		skip             bool
	}{{
		desc:   "without policy",
		gnmiOp: func() {},
	}, {
		desc: "with single policy",
		gnmiOp: func() {
			policy2Pfx := "2002::/48"
			gnmi.Replace(t, dut, ocpath.Root().BgpGueIpv6Policy(policy2Pfx).Config(), &oc.BgpGueIpv6Policy{
				DstPortIpv6: ygot.Uint16(168),
				Prefix:      ygot.String(policy2Pfx),
				SrcIp:       ygot.String("8484:8484::"),
			})
		},
		wantEncapFields1: &EncapFields{
			srcIP:   net.IP{0x84, 0x84, 0x84, 0x84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			dstIP:   net.IP{0x20, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			dstPort: 168,
		},
		wantEncapFields2: &EncapFields{
			srcIP:   net.IP{0x84, 0x84, 0x84, 0x84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			dstIP:   net.IP{0x20, 0x02, 0, 0, 0x10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			dstPort: 168,
		},
	}, {
		desc: "with two overlapping policies",
		gnmiOp: func() {
			policy1Pfx := "2002::/64"
			gnmi.Replace(t, dut, ocpath.Root().BgpGueIpv6Policy(policy1Pfx).Config(), &oc.BgpGueIpv6Policy{
				DstPortIpv6: ygot.Uint16(84),
				Prefix:      ygot.String(policy1Pfx),
				SrcIp:       ygot.String("4242:4242::"),
			})
		},
		wantEncapFields1: &EncapFields{
			srcIP:   net.IP{0x42, 0x42, 0x42, 0x42, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			dstIP:   net.IP{0x20, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			dstPort: 84,
		},
		wantEncapFields2: &EncapFields{
			srcIP:   net.IP{0x84, 0x84, 0x84, 0x84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			dstIP:   net.IP{0x20, 0x02, 0, 0, 0x10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			dstPort: 168,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}

			tests := []struct {
				startingIP      string
				wantEncapFields *EncapFields
			}{{
				startingIP:      "2003:aaaa::",
				wantEncapFields: tt.wantEncapFields1,
			}, {
				startingIP:      "2003:bbbb::",
				wantEncapFields: tt.wantEncapFields2,
			}, {
				startingIP:      "2003:cccc::",
				wantEncapFields: tt.wantEncapFields3,
			}}

			tt.gnmiOp()
			for _, tt := range tests {
				t.Run(tt.startingIP, func(t *testing.T) {
					testTrafficAndEncap(t, otg, tt.startingIP, tt.wantEncapFields)
				})
			}
		})
	}

	dut2.RawAPIs().GRIBI().Default(t).Flush(context.Background(), &gribipb.FlushRequest{
		NetworkInstance: &gribipb.FlushRequest_All{All: &gribipb.Empty{}},
	})
	// TODO: Test that entries are deleted and that there is no more traffic.
}
