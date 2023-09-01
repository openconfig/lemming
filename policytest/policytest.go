// Copyright 2023 Google LLC
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

package policytest

import (
	"testing"
	"time"

	"github.com/openconfig/lemming/bgp/tests/proto/policyval"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/gnmi/oc"
	"github.com/openconfig/ondatra/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

	valpb "github.com/openconfig/lemming/bgp/tests/proto/policyval"
)

// Settings for configuring the baseline testbed with the test
// topology.
//
// The testbed consists of dut:port1 -> dut2:port1
//
//   - dut:port1 -> dut2:port1 subnet 192.0.2.0/30
const (
	ipv4PrefixLen = 24

	awaitTimeout  = 30 * time.Second
	rejectTimeout = 20 * time.Second
)

var (
	BGPPath           = ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).Bgp()
	RoutingPolicyPath = ocpath.Root().RoutingPolicy()
)

type Device struct {
	*ondatra.DUTDevice
	AS       uint32
	RouterID string
}

// TestCase contains the specifications for a single policy test.
//
// Topology:
//
//	DUT1 (AS 64500) -> DUT2 (AS 64500) -> DUT3 (AS 64501)
//	                    ^
//	                    |
//	DUT4 (AS 64502) -> DUT5 (AS 64500)
//
//	Additionally, DUT0 is present as a neighbour for DUT1, DUT4, and DUT5
//	to allow a static route to be resolvable.
//
// Currently by convention, all policies are installed on DUT1 (export), DUT5
// (export), and DUT2 (import). This is because GoBGP only withdraws routes on
// import policy change after a soft reset:
// https://github.com/osrg/gobgp/blob/master/docs/sources/policy.md#policy-and-soft-reset
type TestCase struct {
	Spec            *valpb.PolicyTestCase
	InstallPolicies func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device)
}

// TestPolicy is the helper policy integration tests can call to instantiate
// policy tests.
func TestPolicy(t *testing.T, testspec TestCase) {
	t.Helper()
	t.Run("installPolicyBeforeRoutes", func(t *testing.T) {
		testPolicyAux(t, testspec, false)
	})

	t.Run("installPolicyAfterRoutes", func(t *testing.T) {
		// TODO(wenbli): Figure out how to reset and unskip this.
		testPolicyAux(t, testspec, true)
	})
}

var (
	dut0Ports = map[string]*attrs.Attributes{
		"port1": {
			Desc:    "dut0Port1",
			IPv4:    "192.0.1.0",
			IPv4Len: ipv4PrefixLen,
		},
		"port2": {
			Desc:    "dut0Port2",
			IPv4:    "192.0.4.0",
			IPv4Len: ipv4PrefixLen,
		},
		"port3": {
			Desc:    "dut0Port3",
			IPv4:    "192.0.5.0",
			IPv4Len: ipv4PrefixLen,
		},
	}

	dut1Ports = map[string]*attrs.Attributes{
		"port0": {
			Desc:    "dut1Port1",
			IPv4:    "192.0.1.1",
			IPv4Len: ipv4PrefixLen,
		},
		"port1": {
			Desc:    "dut1Port1",
			IPv4:    "192.1.0.1",
			IPv4Len: ipv4PrefixLen,
		},
	}

	dut2Ports = map[string]*attrs.Attributes{
		"port1": {
			Desc:    "dut2Port1",
			IPv4:    "192.1.0.2",
			IPv4Len: ipv4PrefixLen,
		},
		"port2": {
			Desc:    "dut2Port2",
			IPv4:    "192.2.0.2",
			IPv4Len: ipv4PrefixLen,
		},
		"port3": {
			Desc:    "dut2Port3",
			IPv4:    "192.3.0.2",
			IPv4Len: ipv4PrefixLen,
		},
	}

	dut3Ports = map[string]*attrs.Attributes{
		"port1": {
			Desc:    "dut3Port1",
			IPv4:    "192.2.0.3",
			IPv4Len: ipv4PrefixLen,
		},
	}

	dut4Ports = map[string]*attrs.Attributes{
		"port0": {
			Desc:    "dut1Port1",
			IPv4:    "192.0.4.4",
			IPv4Len: ipv4PrefixLen,
		},
		"port1": {
			Desc:    "dut4Port1",
			IPv4:    "192.4.0.4",
			IPv4Len: ipv4PrefixLen,
		},
	}

	dut5Ports = map[string]*attrs.Attributes{
		"port0": {
			Desc:    "dut1Port1",
			IPv4:    "192.0.5.5",
			IPv4Len: ipv4PrefixLen,
		},
		"port1": {
			Desc:    "dut5Port1",
			IPv4:    "192.4.0.5",
			IPv4Len: ipv4PrefixLen,
		},
		"port2": {
			Desc:    "dut5Port2",
			IPv4:    "192.3.0.5",
			IPv4Len: ipv4PrefixLen,
		},
	}
)

func testPropagation(t *testing.T, routeTest *valpb.RouteTestCase, pair1, pair2 devicePair, filterPoliciesInstalled bool) {
	t.Helper()
	v4uni := BGPPath.Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()
	expectedResult := routeTest.GetExpectedResultBeforePolicy()
	if filterPoliciesInstalled {
		expectedResult = routeTest.GetExpectedResult()
	}

	prevDUT, currDUT, nextDUT := pair1.first, pair1.second, pair2.second
	port1, port21, port23, port3 := pair1.firstPort, pair1.secondPort, pair2.firstPort, pair2.secondPort

	prefix := routeTest.GetInput().GetReachPrefix()
	// Check propagation to AdjRibOutPre for all prefixes.
	gnmi.Await(t, prevDUT, v4uni.Neighbor(port21.IPv4).AdjRibOutPre().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
	gnmi.Await(t, prevDUT, v4uni.Neighbor(port21.IPv4).AdjRibOutPost().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
	gnmi.Await(t, currDUT, v4uni.Neighbor(port1.IPv4).AdjRibInPre().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
	switch expectedResult {
	case policyval.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT:
		t.Logf("Waiting for %s to be propagated", prefix)
		gnmi.Await(t, currDUT, v4uni.Neighbor(port1.IPv4).AdjRibInPost().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
		gnmi.Await(t, currDUT, v4uni.LocRib().Route(prefix, oc.UnionString(port1.IPv4), 0).Prefix().State(), awaitTimeout, prefix)
		gnmi.Await(t, currDUT, v4uni.Neighbor(port3.IPv4).AdjRibOutPre().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
		gnmi.Await(t, currDUT, v4uni.Neighbor(port3.IPv4).AdjRibOutPost().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
		gnmi.Await(t, nextDUT, v4uni.Neighbor(port23.IPv4).AdjRibInPre().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
	case policyval.RouteTestResult_ROUTE_TEST_RESULT_DISCARD:
		w := gnmi.Watch(t, currDUT, v4uni.Neighbor(port1.IPv4).AdjRibInPost().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return !ok
		})
		if _, ok := w.Await(t); !ok {
			t.Errorf("prefix %q (%s) was not rejected from adj-rib-in-post of %v (neighbour %v) within timeout.", prefix, routeTest.GetDescription(), currDUT, prevDUT)
			break
		}
		t.Logf("prefix %q (%s) was successfully rejected from adj-rib-in-post of %v (neighbour %v) within timeout.", prefix, routeTest.GetDescription(), currDUT, prevDUT)

		// Test withdrawal in the case of InstallPolicyAfterRoutes.
		w = gnmi.Watch(t, nextDUT, v4uni.Neighbor(port23.IPv4).AdjRibInPre().Route(prefix, 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return !ok
		})
		if _, ok := w.Await(t); !ok {
			t.Errorf("prefix %q (%s) was not rejected from adj-rib-in-pre of %v (neighbour %v) within timeout.", prefix, routeTest.GetDescription(), nextDUT, currDUT)
			break
		}
		t.Logf("prefix %q (%s) was successfully rejected from adj-rib-in-pre of %v (neighbour %v) within timeout.", prefix, routeTest.GetDescription(), nextDUT, currDUT)
	case policyval.RouteTestResult_ROUTE_TEST_RESULT_NOT_PREFERRED:
		gnmi.Await(t, currDUT, v4uni.Neighbor(port1.IPv4).AdjRibInPost().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
		w := gnmi.Watch(t, currDUT, v4uni.LocRib().Route(prefix, oc.UnionString(port1.IPv4), 0).Prefix().State(), rejectTimeout, func(val *ygnmi.Value[string]) bool {
			_, ok := val.Val()
			return !ok
		})
		if _, ok := w.Await(t); !ok {
			t.Errorf("prefix %q with origin %q (%s) was selected into loc-rib of %v.", prefix, prevDUT, routeTest.GetDescription(), currDUT)
			break
		}
		t.Logf("prefix %q with origin %q (%s) was successfully not selected into loc-rib of %v within timeout.", prefix, prevDUT, routeTest.GetDescription(), currDUT)

		gnmi.Await(t, currDUT, v4uni.Neighbor(port3.IPv4).AdjRibOutPre().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
		gnmi.Await(t, currDUT, v4uni.Neighbor(port3.IPv4).AdjRibOutPost().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
		gnmi.Await(t, nextDUT, v4uni.Neighbor(port23.IPv4).AdjRibInPre().Route(prefix, 0).Prefix().State(), awaitTimeout, prefix)
	default:
		t.Fatalf("Invalid or unhandled policy result: %v", expectedResult)
	}
}

// configureDUT configures the ports on the DUT.
func configureDUT(t *testing.T, dut *ondatra.DUTDevice, ports map[string]*attrs.Attributes) {
	for portName, attr := range ports {
		p := dut.Port(t, portName)
		gnmi.Replace(t, dut, gnmi.OC().Interface(p.Name()).Config(), attr.NewOCInterface(p.Name(), dut))
		gnmi.Await(t, dut, gnmi.OC().Interface(p.Name()).Subinterface(0).Ipv4().Address(attr.IPv4).Ip().State(), awaitTimeout, attr.IPv4)
	}
}

type devicePair struct {
	first      *Device
	second     *Device
	firstPort  *attrs.Attributes
	secondPort *attrs.Attributes
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

func installStaticRoute(t *testing.T, dut *Device, route *oc.NetworkInstance_Protocol_Static) {
	staticp := gnmi.OC().NetworkInstance(fakedevice.DefaultNetworkInstance).
		Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	gnmi.Replace(t, dut, staticp.Static(*route.Prefix).Config(), route)
	gnmi.Await(t, dut, staticp.Static(*route.Prefix).State(), 30*time.Second, route)
}

func awaitSessionEstablished(t *testing.T, dutPair devicePair) {
	t.Helper()
	dut1, dut2 := dutPair.first, dutPair.second
	port1, port2 := dutPair.firstPort, dutPair.secondPort
	gnmi.Await(t, dut1, BGPPath.Neighbor(port2.IPv4).SessionState().State(), awaitTimeout, oc.Bgp_Neighbor_SessionState_ESTABLISHED)
	gnmi.Await(t, dut2, BGPPath.Neighbor(port1.IPv4).SessionState().State(), awaitTimeout, oc.Bgp_Neighbor_SessionState_ESTABLISHED)
}

func establishSessionPairs(t *testing.T, dutPairs ...devicePair) {
	t.Helper()
	for _, pair := range dutPairs {
		dut1, dut2 := pair.first, pair.second
		port1, port2 := pair.firstPort, pair.secondPort
		dutConf := bgpWithNbr(dut1.AS, dut1.RouterID, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
			PeerAs:          ygot.Uint32(dut2.AS),
			NeighborAddress: ygot.String(port2.IPv4),
			Transport: &oc.NetworkInstance_Protocol_Bgp_Neighbor_Transport{
				LocalAddress: ygot.String(port1.IPv4),
			},
		})
		dut2Conf := bgpWithNbr(dut2.AS, dut2.RouterID, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
			PeerAs:          ygot.Uint32(dut1.AS),
			NeighborAddress: ygot.String(port1.IPv4),
			Transport: &oc.NetworkInstance_Protocol_Bgp_Neighbor_Transport{
				LocalAddress: ygot.String(port2.IPv4),
			},
		})
		gnmi.Update(t, dut1, BGPPath.Config(), dutConf)
		gnmi.Update(t, dut2, BGPPath.Config(), dut2Conf)
	}

	for _, pair := range dutPairs {
		awaitSessionEstablished(t, pair)
	}
}

func installDefaultAllowPolicies(t *testing.T, dutPair devicePair) {
	t.Helper()
	dut1, dut2 := dutPair.first, dutPair.second
	port1, port2 := dutPair.firstPort, dutPair.secondPort
	gnmi.Replace(t, dut1, BGPPath.Neighbor(port2.IPv4).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
	gnmi.Replace(t, dut2, BGPPath.Neighbor(port1.IPv4).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
}

func testPolicyAux(t *testing.T, testspec TestCase, installPolicyAfterRoutes bool) {
	dut0 := ondatra.DUT(t, "dut0")
	dut1 := &Device{
		DUTDevice: ondatra.DUT(t, "dut1"),
		AS:        64500,
		RouterID:  dut1Ports["port1"].IPv4,
	}
	dut2 := &Device{
		DUTDevice: ondatra.DUT(t, "dut2"),
		AS:        64500,
		RouterID:  dut2Ports["port1"].IPv4,
	}
	dut3 := &Device{
		DUTDevice: ondatra.DUT(t, "dut3"),
		AS:        64501,
		RouterID:  dut3Ports["port1"].IPv4,
	}
	dut4 := &Device{
		DUTDevice: ondatra.DUT(t, "dut4"),
		AS:        64502,
		RouterID:  dut4Ports["port1"].IPv4,
	}
	dut5 := &Device{
		DUTDevice: ondatra.DUT(t, "dut5"),
		AS:        64500,
		RouterID:  dut5Ports["port1"].IPv4,
	}
	configureDUT(t, dut0, dut0Ports)
	configureDUT(t, dut1.DUTDevice, dut1Ports)
	configureDUT(t, dut2.DUTDevice, dut2Ports)
	configureDUT(t, dut3.DUTDevice, dut3Ports)
	configureDUT(t, dut4.DUTDevice, dut4Ports)
	configureDUT(t, dut5.DUTDevice, dut5Ports)

	// Remove any existing BGP config
	gnmi.Delete(t, dut1, BGPPath.Config())
	gnmi.Delete(t, dut2, BGPPath.Config())
	gnmi.Delete(t, dut3, BGPPath.Config())
	gnmi.Delete(t, dut4, BGPPath.Config())
	gnmi.Delete(t, dut5, BGPPath.Config())
	gnmi.Delete(t, dut1, RoutingPolicyPath.Config())
	gnmi.Delete(t, dut2, RoutingPolicyPath.Config())
	gnmi.Delete(t, dut3, RoutingPolicyPath.Config())
	gnmi.Delete(t, dut4, RoutingPolicyPath.Config())
	gnmi.Delete(t, dut5, RoutingPolicyPath.Config())

	pair12 := devicePair{dut1, dut2, dut1Ports["port1"], dut2Ports["port1"]}
	pair23 := devicePair{dut2, dut3, dut2Ports["port2"], dut3Ports["port1"]}
	pair45 := devicePair{dut4, dut5, dut4Ports["port1"], dut5Ports["port1"]}
	pair52 := devicePair{dut5, dut2, dut5Ports["port2"], dut2Ports["port3"]}

	// Clear the path for routes to be propagated.
	// DUT1 -> DUT2 -> DUT3
	installDefaultAllowPolicies(t, pair12)
	installDefaultAllowPolicies(t, pair23)
	// This is an alternate source of routes towards DUT2 and thereby DUT3.
	// Note that this path is longer than the above path:
	// DUT4 -> DUT5 -> DUT2 (-> DUT3)
	installDefaultAllowPolicies(t, pair45)
	installDefaultAllowPolicies(t, pair52)

	if testspec.InstallPolicies != nil && !installPolicyAfterRoutes {
		testspec.InstallPolicies(t, dut1, dut2, dut3, dut4, dut5)
	}

	establishSessionPairs(t, pair12, pair23, pair45, pair52)

	for _, routeTest := range testspec.Spec.RouteTests {
		// Install all regular test routes into DUT1.
		route := &oc.NetworkInstance_Protocol_Static{
			Prefix: ygot.String(routeTest.GetInput().GetReachPrefix()),
			NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
				"single": {
					Index:   ygot.String("single"),
					NextHop: oc.UnionString(dut1Ports["port0"].IPv4),
					Recurse: ygot.Bool(true),
				},
			},
		}
		installStaticRoute(t, dut1, route)
	}

	for _, routeTest := range testspec.Spec.LongerPathRouteTests {
		// Install all longer-path test routes into DUT4.
		route := &oc.NetworkInstance_Protocol_Static{
			Prefix: ygot.String(routeTest.GetInput().GetReachPrefix()),
			NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
				"single": {
					Index:   ygot.String("single"),
					NextHop: oc.UnionString(dut4Ports["port0"].IPv4),
					Recurse: ygot.Bool(true),
				},
			},
		}
		installStaticRoute(t, dut4, route)
	}

	for _, routeTest := range testspec.Spec.AlternatePathRouteTests {
		// Install all alternate-path test routes into DUT5.
		route := &oc.NetworkInstance_Protocol_Static{
			Prefix: ygot.String(routeTest.GetInput().GetReachPrefix()),
			NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
				"single": {
					Index:   ygot.String("single"),
					NextHop: oc.UnionString(dut5Ports["port0"].IPv4),
					Recurse: ygot.Bool(true),
				},
			},
		}
		installStaticRoute(t, dut5, route)
	}

	for _, routeTest := range testspec.Spec.RouteTests {
		testPropagation(t, routeTest, pair12, pair23, !installPolicyAfterRoutes)
	}
	for _, routeTest := range testspec.Spec.LongerPathRouteTests {
		testPropagation(t, routeTest, pair52, pair23, !installPolicyAfterRoutes)
	}

	if installPolicyAfterRoutes {
		if testspec.InstallPolicies != nil {
			testspec.InstallPolicies(t, dut1, dut2, dut3, dut4, dut5)
		}
		// Changing policy currently causes a hard reset of the BGP
		// session, which causes routes to disappear from the AdjRIBs,
		// so we need to wait for re-establishment first. To do this,
		// we sleep for a reasonable amount of time for the sessions to
		// be teared down.
		time.Sleep(10 * time.Second)
		awaitSessionEstablished(t, pair12)
		awaitSessionEstablished(t, pair52)

		for _, routeTest := range testspec.Spec.RouteTests {
			testPropagation(t, routeTest, pair12, pair23, true)
		}
		for _, routeTest := range testspec.Spec.LongerPathRouteTests {
			testPropagation(t, routeTest, pair52, pair23, true)
		}
	}
}