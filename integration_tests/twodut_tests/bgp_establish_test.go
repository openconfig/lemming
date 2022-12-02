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
	"testing"
	"time"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/integration_tests/binding"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ygot/ygot"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.Get("."))
}

// Settings for configuring the baseline testbed with the test
// topology.
//
// The testbed consists of dut:port1 -> dut2:port1
//
//   - dut:port1 -> dut2:port1 subnet 192.0.2.0/30
const (
	ipv4PrefixLen = 30
)

const (
	dutAS  = 64500
	dut2AS = 64501
)

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

	dut2Port1 = Attributes{
		Desc:    "dut2Port1",
		IPv4:    "192.0.2.2",
		IPv4Len: ipv4PrefixLen,
	}
)

// configInterfaceDUT configures the interface with the Addrs.
func configInterfaceDUT(i *oc.Interface, a *Attributes) *oc.Interface {
	i.Description = ygot.String(a.Desc)
	i.Type = oc.IETFInterfaces_InterfaceType_ethernetCsmacd

	s := i.GetOrCreateSubinterface(0)
	s.Enabled = ygot.Bool(true)
	s4 := s.GetOrCreateIpv4()
	s4a := s4.GetOrCreateAddress(a.IPv4)
	s4a.PrefixLength = ygot.Uint8(ipv4PrefixLen)

	return i
}

// configureDUT configures port1 on the DUT.
func configureDUT(t *testing.T, dut *ondatra.DUTDevice, attr Attributes) {
	p1 := dut.Port(t, "port1")
	i1 := &oc.Interface{Name: ygot.String(p1.Name())}
	gnmi.Replace(t, dut, ocpath.Root().Interface(p1.Name()).Config(), configInterfaceDUT(i1, &attr))

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv4().Address(attr.IPv4).Ip().State(), time.Minute, attr.IPv4)
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

func TestEstablish(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	configureDUT(t, dut, dutPort1)
	dut2 := ondatra.DUT(t, "dut2")
	configureDUT(t, dut2, dut2Port1)

	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()

	// Remove any existing BGP config
	gnmi.Delete(t, dut, bgpPath.Config())
	gnmi.Delete(t, dut2, bgpPath.Config())

	// Start a new session
	dutConf := bgpWithNbr(dutAS, dutPort1.IPv4, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dut2AS),
		NeighborAddress: ygot.String(dut2Port1.IPv4),
	})
	dut2Conf := bgpWithNbr(dut2AS, dut2Port1.IPv4, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dutAS),
		NeighborAddress: ygot.String(dutPort1.IPv4),
	})
	gnmi.Replace(t, dut, bgpPath.Config(), dutConf)
	gnmi.Replace(t, dut2, bgpPath.Config(), dut2Conf)

	nbrPath := bgpPath.Neighbor(dut2Port1.IPv4)
	gnmi.Await(t, dut, nbrPath.SessionState().State(), 60*time.Second, oc.Bgp_Neighbor_SessionState_ESTABLISHED)
}
