// Copyright 2024 Google LLC
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

package local_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/lemming/bgp"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/testing/protocmp"
)

// This test case creates a new topology in which all IP addresses are fixed. The topology consists
// of only two connected devices. The second device sets the next-hop to itself.
func TestNextHopAttr(t *testing.T) {

	prefix := "10.0.0.0/10"

	dut1, dut2, stop := setupShortTopology(t)
	defer stop()

	// Install regular test route into DUT1.
	route := &oc.NetworkInstance_Protocol_Static{
		Prefix: ygot.String(prefix),
		NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
			"single": {
				Index:   ygot.String("single"),
				NextHop: oc.UnionString("192.0.2.1"),
				Recurse: ygot.Bool(true),
			},
		},
	}
	installStaticRoute(t, dut1, route)

	testAttrsShortPath(t, prefix, dut1, dut2, testCaseShortPath{
		Dut1AdjRibOutPreNextHop:  "192.0.2.1",
		Dut1AdjRibOutPostNextHop: "192.0.2.1",
		Dut2AdjRibInPreNextHop:   "192.0.2.1",
		Dut2AdjRibInPostNextHop:  "192.0.2.1",
		Dut2LocalRibNextHop:      "192.0.2.1",
	})
}

func setupShortTopology(t *testing.T) (*Device, *Device, func()) {
	t.Helper()

	dut1, stop1 := newLemming(t, 1, 64500, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "192.0.2.1/30",
		niName:  "DEFAULT",
	}})
	dut2, stop2 := newLemming(t, 2, 64500, nil)

	// Remove any existing BGP config
	Delete(t, dut1, bgp.BGPPath.Config())
	Delete(t, dut2, bgp.BGPPath.Config())
	Delete(t, dut1, bgp.RoutingPolicyPath.Config())
	Delete(t, dut2, bgp.RoutingPolicyPath.Config())

	establishSessionPairs(t, DevicePair{dut1, dut2})

	// Clear the path for routes to be propagated.
	Replace(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
	Replace(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)

	// Wait until policies are installed.
	Await(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultExportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
	Await(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().DefaultImportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)

	return dut1, dut2, func() {
		stop1()
		stop2()
	}
}

type testCaseShortPath struct {
	Dut1AdjRibOutPreNextHop  string
	Dut1AdjRibOutPostNextHop string
	Dut2AdjRibInPreNextHop   string
	Dut2AdjRibInPostNextHop  string
	Dut2LocalRibNextHop      string
}

func testAttrsShortPath(t *testing.T, prefix string, dut1, dut2 *Device, routeTest testCaseShortPath) {
	dut1AttrSetMap := Lookup(t, dut1, bgp.BGPPath.Rib().AttrSetMap().State())
	dut1AttrMap, _ := dut1AttrSetMap.Val()
	dut2AttrSetMap := Lookup(t, dut2, bgp.BGPPath.Rib().AttrSetMap().State())
	dut2AttrMap, _ := dut2AttrSetMap.Val()
	updateAttrMaps := func() {
		dut1AttrSetMap = Lookup(t, dut1, bgp.BGPPath.Rib().AttrSetMap().State())
		dut1AttrMap, _ = dut1AttrSetMap.Val()
		dut2AttrSetMap = Lookup(t, dut2, bgp.BGPPath.Rib().AttrSetMap().State())
		dut2AttrMap, _ = dut2AttrSetMap.Val()
	}

	v4uni := bgp.BGPPath.Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()

	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, dut1, dut1AttrMap, v4uni.Neighbor(dut2.RouterID).AdjRibOutPre().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(routeTest.Dut1AdjRibOutPreNextHop, attrs.GetNextHop(), protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Fatalf("DUT 1 AdjRibOutPre attribute difference (prefix %s) (-want, +got):\n%s", prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, dut1, dut1AttrMap, v4uni.Neighbor(dut2.RouterID).AdjRibOutPost().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(routeTest.Dut1AdjRibOutPostNextHop, attrs.GetNextHop(), protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Fatalf("DUT 1 AdjRibOutPost attribute difference (prefix %s) (-want, +got):\n%s", prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, dut2, dut2AttrMap, v4uni.Neighbor(dut1.RouterID).AdjRibInPre().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(routeTest.Dut2AdjRibInPreNextHop, attrs.GetNextHop(), protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Fatalf("DUT 2 AdjRibInPre attribute difference (prefix %s) (-want, +got):\n%s", prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, dut2, dut2AttrMap, v4uni.Neighbor(dut1.RouterID).AdjRibInPost().Route(prefix, 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(routeTest.Dut2AdjRibInPostNextHop, attrs.GetNextHop(), protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Fatalf("DUT 2 AdjRibInPost attribute difference (prefix %s) (-want, +got):\n%s", prefix, diff)
	}
	if diff := awaitNoDiff(func() string {
		attrs, err := getAttrs(t, dut2, dut2AttrMap, v4uni.LocRib().Route(prefix, oc.UnionString(dut1.RouterID), 0).AttrIndex().State())
		if err != nil {
			return err.Error()
		}
		return cmp.Diff(routeTest.Dut2LocalRibNextHop, attrs.GetNextHop(), protocmp.Transform())
	}, updateAttrMaps); diff != "" {
		t.Fatalf("DUT 2 LocRib routeTest difference (prefix %s) (-want, +got):\n%s", prefix, diff)
	}
}
