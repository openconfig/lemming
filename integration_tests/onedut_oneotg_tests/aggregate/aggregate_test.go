// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aggregate_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"strconv"
	"strings"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/gnmi/oc"
	otgtelemetry "github.com/openconfig/ondatra/gnmi/otg"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/lemming/internal/binding"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.KNE("."))
}

// Settings for configuring the aggregate testbed with the test
// topology.  IxNetwork flow requires both source and destination
// networks be configured on the ATE.  It is not possible to send
// packets to the ether.
//
// The testbed consists of ate:port1 -> dut:port1 and dut:port{2-9} ->
// ate:port{2-9}.  The first pair is called the "source" pair, and the
// second aggregate link the "destination" pair.
//
//   - Source: ate:port1 -> dut:port1 subnet 192.0.2.0/30 2001:db8::0/126
//   - Destination: dut:port{2-9} -> ate:port{2-9}
//     subnet 192.0.2.4/30 2001:db8::4/126
//
// Note that the first (.0, .4) and last (.3, .7) IPv4 addresses are
// reserved from the subnet for broadcast, so a /30 leaves exactly 2
// usable addresses.  This does not apply to IPv6 which allows /127
// for point to point links, but we use /126 so the numbering is
// consistent with IPv4.
//
// A traffic flow is configured from ate:port1 as source and ate:port{2-9}
// as destination.
const (
	plen4 = 30
	plen6 = 126
)

var (
	dutSrc = attrs.Attributes{
		Desc:    "dutsrc",
		IPv4:    "192.0.2.1",
		IPv6:    "2001:db8::1",
		IPv4Len: plen4,
		IPv6Len: plen6,
	}

	ateSrc = attrs.Attributes{
		Name:    "atesrc",
		MAC:     "02:11:01:00:00:01",
		IPv4:    "192.0.2.2",
		IPv6:    "2001:db8::2",
		IPv4Len: plen4,
		IPv6Len: plen6,
	}

	dutDst = attrs.Attributes{
		Desc:    "dutdst",
		IPv4:    "192.0.2.5",
		IPv6:    "2001:db8::5",
		IPv4Len: plen4,
		IPv6Len: plen6,
	}

	ateDst = attrs.Attributes{
		Name:    "atedst",
		MAC:     "02:12:01:00:00:01",
		IPv4:    "192.0.2.6",
		IPv6:    "2001:db8::6",
		IPv4Len: plen4,
		IPv6Len: plen6,
	}
)

const (
	lagTypeLACP   = oc.IfAggregate_AggregationType_LACP
	lagTypeSTATIC = oc.IfAggregate_AggregationType_STATIC
)

type testCase struct {
	lagType oc.E_IfAggregate_AggregationType

	dut *ondatra.DUTDevice
	ate *ondatra.ATEDevice
	top gosnappi.Config

	dutPorts []*ondatra.Port
	atePorts []*ondatra.Port
	aggID    string
}

func (tc *testCase) configSrcDUT(i *oc.Interface, a *attrs.Attributes) {
	i.Description = ygot.String(a.Desc)
	i.Enabled = ygot.Bool(true)

	s := i.GetOrCreateSubinterface(0)
	s4 := s.GetOrCreateIpv4()
	s4.Enabled = ygot.Bool(true)

	a4 := s4.GetOrCreateAddress(a.IPv4)
	a4.PrefixLength = ygot.Uint8(plen4)

	s6 := s.GetOrCreateIpv6()
	s6.Enabled = ygot.Bool(true)

	s6.GetOrCreateAddress(a.IPv6).PrefixLength = ygot.Uint8(plen6)
}

func (tc *testCase) configDstAggregateDUT(i *oc.Interface, a *attrs.Attributes) {
	tc.configSrcDUT(i, a)
	i.Type = ieee8023adLag
	g := i.GetOrCreateAggregation()
	g.LagType = tc.lagType
}

func (tc *testCase) configDstMemberDUT(i *oc.Interface, p *ondatra.Port) {
	i.Description = ygot.String(p.String())
	i.Type = ethernetCsmacd
	i.Enabled = ygot.Bool(true)

	e := i.GetOrCreateEthernet()
	e.AggregateId = ygot.String(tc.aggID)
}

func (tc *testCase) setupAggregateAtomically(t *testing.T) {
	d := &oc.Root{}

	if tc.lagType == lagTypeLACP {
		d.GetOrCreateLacp().GetOrCreateInterface(tc.aggID)
	}

	agg := d.GetOrCreateInterface(tc.aggID)
	agg.GetOrCreateAggregation().LagType = tc.lagType
	agg.Type = ieee8023adLag

	for _, port := range tc.dutPorts[1:] {
		i := d.GetOrCreateInterface(port.Name())
		i.GetOrCreateEthernet().AggregateId = ygot.String(tc.aggID)
		i.Type = ethernetCsmacd
		i.Enabled = ygot.Bool(true)
	}

	p := gnmi.OC()
	gnmi.Update(t, tc.dut, p.Config(), d)
}

func (tc *testCase) clearAggregate(t *testing.T) {
	// Clear the aggregate minlink.
	gnmi.Delete(t, tc.dut, gnmi.OC().Interface(tc.aggID).Aggregation().MinLinks().Config())

	// Clear the members of the aggregate.
	for _, port := range tc.dutPorts[1:] {
		gnmi.Delete(t, tc.dut, gnmi.OC().Interface(port.Name()).Ethernet().AggregateId().Config())
	}
}

func (tc *testCase) configureDUT(t *testing.T) {
	t.Logf("dut ports = %v", tc.dutPorts)
	if len(tc.dutPorts) < 2 {
		t.Fatalf("Testbed requires at least 2 ports, got %d", len(tc.dutPorts))
	}

	d := gnmi.OC()

	lacp := &oc.Lacp_Interface{Name: ygot.String(tc.aggID)}
	if tc.lagType == lagTypeLACP {
		lacp.LacpMode = oc.Lacp_LacpActivityType_ACTIVE
	} else {
		lacp.LacpMode = oc.Lacp_LacpActivityType_UNSET
	}
	lacpPath := d.Lacp().Interface(tc.aggID)
	if tc.lagType == lagTypeLACP {
		gnmi.Replace(t, tc.dut, lacpPath.Config(), lacp)
	}

	agg := &oc.Interface{Name: ygot.String(tc.aggID)}
	tc.configDstAggregateDUT(agg, &dutDst)
	aggPath := d.Interface(tc.aggID)
	gnmi.Replace(t, tc.dut, aggPath.Config(), agg)

	srcp := tc.dutPorts[0]
	srci := &oc.Interface{Name: ygot.String(srcp.Name())}
	tc.configSrcDUT(srci, &dutSrc)
	srci.Type = ethernetCsmacd
	srciPath := d.Interface(srcp.Name())
	gnmi.Replace(t, tc.dut, srciPath.Config(), srci)
	for _, port := range tc.dutPorts[1:] {
		i := &oc.Interface{Name: ygot.String(port.Name())}
		i.Type = ethernetCsmacd
		i.Enabled = ygot.Bool(true)

		tc.configDstMemberDUT(i, port)
		iPath := d.Interface(port.Name())
		gnmi.Replace(t, tc.dut, iPath.Config(), i)
	}
}

func (tc *testCase) configureATE(t *testing.T) {
	if len(tc.atePorts) < 2 {
		t.Fatalf("Testbed requires at least 2 ports, got: %v", tc.atePorts)
	}

	p0 := tc.atePorts[0]
	tc.top.Ports().Add().SetName(p0.ID())
	srcDev := tc.top.Devices().Add().SetName(ateSrc.Name)
	srcEth := srcDev.Ethernets().Add().SetName(ateSrc.Name + ".Eth").SetMac(ateSrc.MAC)
	srcEth.Connection().SetChoice(gosnappi.EthernetConnectionChoice.PORT_NAME).SetPortName(p0.ID())
	srcEth.Ipv4Addresses().Add().SetName(ateSrc.Name + ".IPv4").SetAddress(ateSrc.IPv4).SetGateway(dutSrc.IPv4).SetPrefix(uint32(ateSrc.IPv4Len))
	srcEth.Ipv6Addresses().Add().SetName(ateSrc.Name + ".IPv6").SetAddress(ateSrc.IPv6).SetGateway(dutSrc.IPv6).SetPrefix(uint32(ateSrc.IPv6Len))

	// Adding the rest of the ports to the configuration and to the LAG
	agg := tc.top.Lags().Add().SetName(ateDst.Name)
	if tc.lagType == lagTypeSTATIC {
		lagId, _ := strconv.Atoi(tc.aggID)
		agg.Protocol().SetChoice("static").Static().SetLagId(uint32(lagId))
		for i, p := range tc.atePorts[1:] {
			port := tc.top.Ports().Add().SetName(p.ID())
			newMac, err := incrementMAC(ateDst.MAC, i+1)
			if err != nil {
				t.Fatal(err)
			}
			agg.Ports().Add().SetPortName(port.Name()).Ethernet().SetMac(newMac).SetName("LAGRx-" + strconv.Itoa(i))
		}
	} else {
		agg.Protocol().SetChoice("lacp")
		agg.Protocol().Lacp().SetActorKey(1).SetActorSystemPriority(1).SetActorSystemId(ateDst.MAC)
		for i, p := range tc.atePorts[1:] {
			port := tc.top.Ports().Add().SetName(p.ID())
			newMac, err := incrementMAC(ateDst.MAC, i+1)
			if err != nil {
				t.Fatal(err)
			}
			lagPort := agg.Ports().Add().SetPortName(port.Name())
			lagPort.Ethernet().SetMac(newMac).SetName("LAGRx-" + strconv.Itoa(i))
			lagPort.Lacp().SetActorActivity("active").SetActorPortNumber(uint32(i) + 1).SetActorPortPriority(1).SetLacpduTimeout(0)
		}
	}

	// Disable FEC for 100G-FR ports because Novus does not support it.
	p100gbasefr := []string{}
	for _, p := range tc.atePorts {
		if p.PMD() == ondatra.PMD100GBASEFR {
			p100gbasefr = append(p100gbasefr, p.ID())
		}
	}

	if len(p100gbasefr) > 0 {
		l1Settings := tc.top.Layer1().Add().SetName("L1").SetPortNames(p100gbasefr)
		l1Settings.SetAutoNegotiate(true).SetIeeeMediaDefaults(false).SetSpeed("speed_100_gbps")
		autoNegotiate := l1Settings.AutoNegotiation()
		autoNegotiate.SetRsFec(false)
	}

	dstDev := tc.top.Devices().Add().SetName(agg.Name() + ".dev")
	dstEth := dstDev.Ethernets().Add().SetName(ateDst.Name + ".Eth").SetMac(ateDst.MAC)
	dstEth.Connection().SetChoice(gosnappi.EthernetConnectionChoice.LAG_NAME).SetLagName(agg.Name())
	dstEth.Ipv4Addresses().Add().SetName(ateDst.Name + ".IPv4").SetAddress(ateDst.IPv4).SetGateway(dutDst.IPv4).SetPrefix(uint32(ateDst.IPv4Len))
	dstEth.Ipv6Addresses().Add().SetName(ateDst.Name + ".IPv6").SetAddress(ateDst.IPv6).SetGateway(dutDst.IPv6).SetPrefix(uint32(ateDst.IPv6Len))

	// Fail early if the topology is bad.
	tc.ate.OTG().PushConfig(t, tc.top)
	tc.ate.OTG().StartProtocols(t)
}

const (
	ethernetCsmacd = oc.IETFInterfaces_InterfaceType_ethernetCsmacd
	ieee8023adLag  = oc.IETFInterfaces_InterfaceType_ieee8023adLag
	adminUp        = oc.Interface_AdminStatus_UP
	opUp           = oc.Interface_OperStatus_UP
	opDown         = oc.Interface_OperStatus_DOWN
	full           = oc.Ethernet_DuplexMode_FULL
	dynamic        = oc.IfIp_NeighborOrigin_DYNAMIC
)

func (tc *testCase) verifyAggID(t *testing.T, dp *ondatra.Port) {
	dip := gnmi.OC().Interface(dp.Name())
	di := gnmi.Get(t, tc.dut, dip.State())
	if lagID := di.GetEthernet().GetAggregateId(); lagID != tc.aggID {
		t.Errorf("%s LagID got %v, want %v", dp, lagID, tc.aggID)
	}
}

func (tc *testCase) verifyInterfaceDUT(t *testing.T, dp *ondatra.Port) {
	dip := gnmi.OC().Interface(dp.Name())
	di := gnmi.Get(t, tc.dut, dip.State())

	if got := di.GetAdminStatus(); got != adminUp {
		t.Errorf("%s admin-status got %v, want %v", dp, got, adminUp)
	}

	// LAG members may fall behind, so wait for them to be up.
	gnmi.Await(t, tc.dut, dip.OperStatus().State(), time.Minute, opUp)
}

func (tc *testCase) verifyDUT(t *testing.T) {
	// Wait for LAG negotiation and verify LAG type for the aggregate interface.
	gnmi.Await(t, tc.dut, gnmi.OC().Interface(tc.aggID).Type().State(), time.Minute, ieee8023adLag)

	for n, port := range tc.dutPorts {
		if n < 1 {
			// We designate port 0 as the source link, not part of LAG.
			t.Run(fmt.Sprintf("%s [source]", port.ID()), func(t *testing.T) {
				tc.verifyInterfaceDUT(t, port)
			})
			continue
		}
		t.Run(fmt.Sprintf("%s [member]", port.ID()), func(t *testing.T) {
			tc.verifyInterfaceDUT(t, port)
			tc.verifyAggID(t, port)
		})
	}
}

// verifyATE checks the telemetry against the parameters set by
// configureDUT().
func (tc *testCase) verifyATE(t *testing.T) {
	ap := tc.atePorts[0]
	// State for the interface.
	time.Sleep(3 * time.Second)

	portMetrics := gnmi.Get(t, tc.ate.OTG(), gnmi.OTG().Port(ap.ID()).State())
	if portMetrics.GetLink() != otgtelemetry.Port_Link_UP {
		t.Errorf("%s oper-status got %v, want %v", ap.ID(), portMetrics.GetLink(), otgtelemetry.Port_Link_UP)
	}
	t.Logf("Checking if LAG is up on OTG")
	gnmi.Watch(t, tc.ate.OTG(), gnmi.OTG().Lag(ateDst.Name).OperStatus().State(), time.Minute, func(val *ygnmi.Value[otgtelemetry.E_Lag_OperStatus]) bool {
		state, present := val.Val()
		return present && state.String() == "UP"
	}).Await(t)
}

// setDutInterfaceWithState sets the admin state to a member of the lag
func (tc *testCase) setDutInterfaceWithState(t testing.TB, p *ondatra.Port, state bool) {
	dc := gnmi.OC()
	i := &oc.Interface{}
	i.Enabled = ygot.Bool(state)
	i.Type = ethernetCsmacd
	i.Name = ygot.String(p.Name())
	gnmi.Update(t, tc.dut, dc.Interface(p.Name()).Config(), i)
}

// sortPorts sorts the ports by the testbed port ID.
func sortPorts(ports []*ondatra.Port) []*ondatra.Port {
	sort.SliceStable(ports, func(i, j int) bool {
		return ports[i].ID() < ports[j].ID()
	})
	return ports
}

// incrementMAC increments the MAC by i. Returns error if the mac cannot be parsed or overflows the mac address space
func incrementMAC(mac string, i int) (string, error) {
	macAddr, err := net.ParseMAC(mac)
	if err != nil {
		return "", err
	}
	convMac := binary.BigEndian.Uint64(append([]byte{0, 0}, macAddr...))
	convMac = convMac + uint64(i)
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, convMac)
	if err != nil {
		return "", err
	}
	newMac := net.HardwareAddr(buf.Bytes()[2:8])
	return newMac.String(), nil
}

func (tc *testCase) verifyMinLinks(t *testing.T) {
	totalPorts := len(tc.dutPorts)
	numLagPorts := totalPorts - 1
	minLinks := uint16(numLagPorts - 1)
	gnmi.Replace(t, tc.dut, gnmi.OC().Interface(tc.aggID).Aggregation().MinLinks().Config(), minLinks)

	tests := []struct {
		desc      string
		downCount int
		want      []oc.E_Interface_OperStatus
	}{
		{
			desc:      "MinLink + 1",
			downCount: 0,
			want:      []oc.E_Interface_OperStatus{opUp},
		},
		{
			desc:      "MinLink",
			downCount: 1,
			want:      []oc.E_Interface_OperStatus{opUp},
		},
		{
			desc:      "MinLink - 1",
			downCount: 2,
			want:      []oc.E_Interface_OperStatus{oc.Interface_OperStatus_LOWER_LAYER_DOWN, opDown},
		},
	}
	for _, tf := range tests {
		t.Run(tf.desc, func(t *testing.T) {
			for _, port := range tc.atePorts[1 : 1+tf.downCount] {
				dp := tc.dut.Port(t, port.ID())
				tc.setDutInterfaceWithState(t, dp, false)
				defer tc.setDutInterfaceWithState(t, dp, true)

				if tc.lagType == oc.IfAggregate_AggregationType_LACP {
					time.Sleep(3 * time.Second)

					t.Logf("Awaiting LAG DUT port: %v to stop collecting", dp)
					gnmi.WatchAll(t, tc.dut, gnmi.OC().Lacp().InterfaceAny().Member(dp.Name()).Collecting().State(), time.Minute, func(val *ygnmi.Value[bool]) bool {
						col, present := val.Val()
						return present && !col
					}).Await(t)
					t.Logf("Awaiting LAG DUT port: %v to stop distributing", dp)
					gnmi.WatchAll(t, tc.dut, gnmi.OC().Lacp().InterfaceAny().Member(dp.Name()).Distributing().State(), time.Minute, func(val *ygnmi.Value[bool]) bool {
						dist, present := val.Val()
						return present && !dist
					}).Await(t)

				}
			}
			opStatus, statusCheckResult := gnmi.Watch(t, tc.dut, gnmi.OC().Interface(tc.aggID).OperStatus().State(), 1*time.Minute, func(y *ygnmi.Value[oc.E_Interface_OperStatus]) bool {
				opStatus, ok := y.Val()
				if !ok {
					return false
				}
				for _, expectedStatus := range tf.want {
					if opStatus == expectedStatus {
						return true
					}
				}
				return false
			}).Await(t)
			if !statusCheckResult {
				val, _ := opStatus.Val()
				t.Errorf("Check of OperStatus for Interface %s is failed, want: %v, got: %s", tc.aggID, tf.want, val.String())
			}
		})
	}
}

func (tc *testCase) getCounters(t *testing.T, when string) map[string]*oc.Interface_Counters {
	results := make(map[string]*oc.Interface_Counters)
	b := &strings.Builder{}
	w := tabwriter.NewWriter(b, 0, 0, 1, ' ', 0)

	fmt.Fprint(w, "Raw Interface Counters\n\n")
	fmt.Fprint(w, "Name\tInUnicastPkts\tInOctets\tOutUnicastPkts\tOutOctets\n")
	for _, port := range tc.dutPorts[1:] {
		counters := gnmi.Get(t, tc.dut, gnmi.OC().Interface(port.Name()).Counters().State())
		results[port.Name()] = counters
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\n",
			port.Name(),
			counters.GetInUnicastPkts(), counters.GetInOctets(),
			counters.GetOutUnicastPkts(), counters.GetOutOctets())
	}
	w.Flush()

	t.Log(b)

	return results
}

// generates a list of random tcp ports values
func generateRandomPortList(count uint) []uint32 {
	a := make([]uint32, count)
	for index := range a {
		a[index] = uint32(rand.Intn(65536-1) + 1)
	}
	return a
}

// normalize normalizes the input values so that the output values sum
// to 1.0 but reflect the proportions of the input.  For example,
// input [1, 2, 3, 4] is normalized to [0.1, 0.2, 0.3, 0.4].
func normalize(xs []uint64) (ys []float64, sum uint64) {
	for _, x := range xs {
		sum += x
	}
	ys = make([]float64, len(xs))
	for i, x := range xs {
		ys[i] = float64(x) / float64(sum)
	}
	return ys, sum
}

var approxOpt = cmpopts.EquateApprox(0 /* frac */, 0.01 /* absolute */)

// portWants converts the nextHop wanted weights to per-port wanted
// weights listed in the same order as atePorts.
func (tc *testCase) portWants() []float64 {
	numPorts := len(tc.dutPorts[1:])
	weights := []float64{}
	for i := 0; i < numPorts; i++ {
		weights = append(weights, 1/float64(numPorts))
	}
	return weights
}

func (tc *testCase) verifyCounterDiff(t *testing.T, before, after map[string]*oc.Interface_Counters) {
	b := &strings.Builder{}
	w := tabwriter.NewWriter(b, 0, 0, 1, ' ', 0)

	fmt.Fprint(w, "Interface Counter Deltas\n\n")
	fmt.Fprint(w, "Name\tInPkts\tInOctets\tOutPkts\tOutOctets\n")
	allInPkts := []uint64{}
	allOutPkts := []uint64{}

	for port := range before {
		inPkts := after[port].GetInUnicastPkts() - before[port].GetInUnicastPkts()
		allInPkts = append(allInPkts, inPkts)
		inOctets := after[port].GetInOctets() - before[port].GetInOctets()
		outPkts := after[port].GetOutUnicastPkts() - before[port].GetOutUnicastPkts()
		allOutPkts = append(allOutPkts, outPkts)
		outOctets := after[port].GetOutOctets() - before[port].GetOutOctets()

		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\n",
			port,
			inPkts, inOctets,
			outPkts, outOctets)
	}
	got, outSum := normalize(allOutPkts)
	want := tc.portWants()
	t.Logf("outPkts normalized got: %v", got)
	t.Logf("want: %v", want)
	t.Run("Ratio", func(t *testing.T) {
		if diff := cmp.Diff(want, got, approxOpt); diff != "" {
			t.Errorf("Packet distribution ratios -want,+got:\n%s", diff)
		}
	})
	t.Run("Loss", func(t *testing.T) {
		if allInPkts[0] > outSum {
			t.Errorf("Traffic flow received %d packets, sent only %d",
				allOutPkts[0], outSum)
		}
	})
	w.Flush()

	t.Log(b)
}

type headerType string

const (
	ipv4Header headerType = "ipv4"
	ipv6Header headerType = "ipv6"
)

func (tc *testCase) testFlow(t *testing.T, l3header headerType) {
	i1 := ateSrc.Name
	i2 := ateDst.Name

	tc.top.Flows().Clear().Items()
	flow := tc.top.Flows().Add().SetName(string(l3header))
	flow.Metrics().SetEnable(true)
	flow.Size().SetFixed(128)
	flow.Packet().Add().Ethernet().Src().SetValue(ateSrc.MAC)

	switch l3header {
	case ipv4Header:
		flow.TxRx().Device().SetTxNames([]string{i1 + ".IPv4"}).SetRxNames([]string{i2 + ".IPv4"})
		v4 := flow.Packet().Add().Ipv4()
		v4.Src().SetValue(ateSrc.IPv4)
		v4.Dst().SetValue(ateDst.IPv4)
	case ipv6Header:
		flow.TxRx().Device().SetTxNames([]string{i1 + ".IPv6"}).SetRxNames([]string{i2 + ".IPv6"})
		v6 := flow.Packet().Add().Ipv6()
		v6.Src().SetValue(ateSrc.IPv6)
		v6.Dst().SetValue(ateDst.IPv6)
	}

	tcp := flow.Packet().Add().Tcp()
	tcp.SrcPort().SetValues(generateRandomPortList(65534))
	tcp.DstPort().SetValues(generateRandomPortList(65534))
	tc.ate.OTG().PushConfig(t, tc.top)
	tc.ate.OTG().StartProtocols(t)

	tc.verifyDUT(t)
	tc.verifyATE(t)

	beforeTrafficCounters := tc.getCounters(t, "before")

	tc.ate.OTG().StartTraffic(t)
	time.Sleep(15 * time.Second)
	tc.ate.OTG().StopTraffic(t)

	recvMetric := gnmi.Get(t, tc.ate.OTG(), gnmi.OTG().Flow(flow.Name()).State())
	pkts := recvMetric.GetCounters().GetOutPkts()

	if pkts == 0 {
		t.Errorf("Flow sent packets: got %v, want non zero", pkts)
	}
	afterTrafficCounters := tc.getCounters(t, "after")
	tc.verifyCounterDiff(t, beforeTrafficCounters, afterTrafficCounters)
}

func TestNegotiation(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	ate := ondatra.ATE(t, "ate")
	aggID := "agg0"

	lagTypes := []oc.E_IfAggregate_AggregationType{lagTypeSTATIC}

	for _, lagType := range lagTypes {
		top := gosnappi.NewConfig()
		// Clean otg with an empty config
		ate.OTG().PushConfig(t, top)

		tc := &testCase{
			dut:     dut,
			ate:     ate,
			top:     top,
			lagType: lagType,

			dutPorts: sortPorts(dut.Ports()),
			atePorts: sortPorts(ate.Ports()),
			aggID:    aggID,
		}
		t.Run(fmt.Sprintf("LagType=%s", lagType), func(t *testing.T) {
			tc.configureDUT(t)
			time.Sleep(5 * time.Second)
			t.Run("VerifyDUT", tc.verifyDUT)

			tc.configureATE(t)
			t.Run("VerifyATE", tc.verifyATE)

			for _, flow := range []headerType{ipv4Header} {
				t.Run(fmt.Sprint("TestFlow ", flow), func(t *testing.T) {
					tc.testFlow(t, flow)
				})
			}

			if lagType == lagTypeLACP { // The Linux kernel bond driver only supports min_links for LACP.
				t.Run("MinLinks", tc.verifyMinLinks)
			}
		})
	}
}
