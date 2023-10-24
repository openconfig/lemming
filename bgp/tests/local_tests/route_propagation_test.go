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

package local_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/openconfig/gribigo/chk"
	"github.com/openconfig/gribigo/client"
	"github.com/openconfig/gribigo/constants"
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/lemming/bgp"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
)

func installStaticRoute(t *testing.T, dut *Device, route *oc.NetworkInstance_Protocol_Static) {
	staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).
		Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	Replace(t, dut, staticp.Static(*route.Prefix).Config(), route)
	Await(t, dut, staticp.Static(*route.Prefix).State(), route)
}

// awaitTimeout calls a fluent client Await, adding a timeout to the context.
func awaitTimeout(ctx context.Context, c *fluent.GRIBIClient, t testing.TB, timeout time.Duration) error {
	subctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Await(subctx, t)
}

func installGRIBIRoutes(ctx context.Context, t *testing.T, dut *Device, entries []fluent.GRIBIEntry, results []*client.OpResult, isDelete bool) {
	if err := awaitTimeout(ctx, dut.gribic, t, time.Minute); err != nil {
		t.Fatalf("Await got error during session negotiation: %v", err)
	}

	dut.gribic.Modify().AddEntry(t, entries...)
	if !isDelete {
	} else {
		dut.gribic.Modify().DeleteEntry(t, entries...)
	}
	if err := awaitTimeout(ctx, dut.gribic, t, time.Minute); err != nil {
		t.Fatalf("Await got error for entries: %v", err)
	}

	for _, wantResult := range results {
		chk.HasResult(t, dut.gribic.Results(t), wantResult, chk.IgnoreOperationID())
	}
}

func formatYgot(v any) string {
	str, err := ygot.Marshal7951(v, ygot.JSONIndent("  "))
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(str)
}

func awaitNotPresent[T any](t *testing.T, d *Device, q ygnmi.SingletonQuery[T]) {
	t.Helper()
	w := Watch(t, d, q, rejectTimeout, func(val *ygnmi.Value[T]) bool {
		_, ok := val.Val()
		return !ok
	})
	if _, ok := w.Await(t); !ok {
		t.Fatalf("prefix was not rejected within timeout at %v.", d)
	}
}

func TestRoutePropagation(t *testing.T) {
	dut1, stop1 := newLemming(t, 1, 64500, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "127.0.0.0/30",
		niName:  "DEFAULT",
	}, {
		name:    "eth1",
		ifindex: 1,
		enabled: true,
		prefix:  "192.0.2.1/31",
		niName:  "DEFAULT",
	}})
	defer stop1()
	dut2, stop2 := newLemming(t, 2, 64501, nil)
	defer stop2()
	dut3, stop3 := newLemming(t, 3, 64502, nil)
	defer stop3()
	dut4, stop4 := newLemming(t, 4, 64503, nil)
	defer stop4()

	installDefaultPolicies := func() {
		fmt.Println("Installing default policies")
		// Clear the path for routes to be propagated.
		Replace(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut3.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut3, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut3, bgp.BGPPath.Neighbor(dut4.RouterID).ApplyPolicy().DefaultExportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Replace(t, dut4, bgp.BGPPath.Neighbor(dut3.RouterID).ApplyPolicy().DefaultImportPolicy().Config(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
	}
	installDefaultPolicies()

	establishSessionPairs(t, []DevicePair{{dut1, dut2}, {dut2, dut3}, {dut3, dut4}}...)

	awaitDefaultPolicies := func() {
		Await(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultExportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().DefaultImportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut2, bgp.BGPPath.Neighbor(dut3.RouterID).ApplyPolicy().DefaultExportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut3, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().DefaultImportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut3, bgp.BGPPath.Neighbor(dut4.RouterID).ApplyPolicy().DefaultExportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
		Await(t, dut4, bgp.BGPPath.Neighbor(dut3.RouterID).ApplyPolicy().DefaultImportPolicy().State(), oc.RoutingPolicy_DefaultPolicyType_ACCEPT_ROUTE)
	}
	awaitDefaultPolicies()

	staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)

	checkpath := func(t *testing.T, prefix string, static bool) {
		if static {
			installStaticRoute(t, dut1, &oc.NetworkInstance_Protocol_Static{
				Prefix: ygot.String(prefix),
				NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
					"single": {
						Index:   ygot.String("single"),
						NextHop: oc.UnionString("192.0.2.1"),
						Recurse: ygot.Bool(true),
					},
				},
			})

			v := GetAll(t, dut1, staticp.StaticAny().Config())
			t.Logf("Installed static route: %s", formatYgot(v))
		} else {
			installGRIBIRoutes(context.Background(), t, dut1,
				[]fluent.GRIBIEntry{
					fluent.NextHopEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
						WithIndex(42).WithIPAddress("192.0.2.1"),
					fluent.NextHopGroupEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
						WithID(10).AddNextHop(42, 1),
					fluent.IPv4Entry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
						WithPrefix(prefix).WithNextHopGroup(10),
				}, []*client.OpResult{
					fluent.OperationResult().
						WithNextHopOperation(42).
						WithProgrammingResult(fluent.InstalledInFIB).
						WithOperationType(constants.Add).
						AsResult(),
					fluent.OperationResult().
						WithNextHopGroupOperation(10).
						WithProgrammingResult(fluent.InstalledInFIB).
						WithOperationType(constants.Add).
						AsResult(),
					fluent.OperationResult().
						WithIPv4Operation(prefix).
						WithProgrammingResult(fluent.InstalledInFIB).
						WithOperationType(constants.Add).
						AsResult(),
				},
				false,
			)
		}

		v4uni := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).Bgp().Rib().AfiSafi(oc.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST).Ipv4Unicast()
		// Check route is in Adj-In of dut2.
		Await(t, dut2, v4uni.Neighbor(dut1.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, dut2, v4uni.Neighbor(dut1.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State(), prefix)
		// Check route is in Loc-RIB of dut2.
		Await(t, dut2, v4uni.LocRib().Route(prefix, oc.UnionString(dut1.RouterID), 0).Prefix().State(), prefix)
		// Check route is in Adj-Out of dut2.
		Await(t, dut2, v4uni.Neighbor(dut3.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State(), prefix)
		Await(t, dut2, v4uni.Neighbor(dut3.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State(), prefix)

		// Check route is in Adj-In of dut3.
		Await(t, dut3, v4uni.Neighbor(dut2.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)

		// Check route is in Adj-In of dut4.
		Await(t, dut4, v4uni.Neighbor(dut3.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State(), prefix)

		// -------------- Test deletion ----------------
		if static {
			Delete[*oc.NetworkInstance_Protocol_Static](t, dut1, staticp.Static(prefix).Config())
		} else {
			installGRIBIRoutes(context.Background(), t, dut1,
				[]fluent.GRIBIEntry{
					fluent.IPv4Entry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
						WithPrefix(prefix).WithNextHopGroup(10),
					fluent.NextHopGroupEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
						WithID(10).AddNextHop(42, 1),
					fluent.NextHopEntry().WithNetworkInstance(fakedevice.DefaultNetworkInstance).
						WithIndex(42).WithIPAddress("192.0.2.1"),
				}, []*client.OpResult{
					fluent.OperationResult().
						WithNextHopOperation(42).
						WithProgrammingResult(fluent.InstalledInFIB).
						WithOperationType(constants.Delete).
						AsResult(),
					fluent.OperationResult().
						WithNextHopGroupOperation(10).
						WithProgrammingResult(fluent.InstalledInFIB).
						WithOperationType(constants.Delete).
						AsResult(),
					fluent.OperationResult().
						WithIPv4Operation(prefix).
						WithProgrammingResult(fluent.InstalledInFIB).
						WithOperationType(constants.Delete).
						AsResult(),
				},
				true,
			)
		}

		// Check route is not in Adj-In of dut2.
		awaitNotPresent(t, dut2, v4uni.Neighbor(dut1.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State())
		awaitNotPresent(t, dut2, v4uni.Neighbor(dut1.RouterID).AdjRibInPost().Route(prefix, 0).Prefix().State())
		// Check route is not in Loc-RIB of dut2.
		awaitNotPresent(t, dut2, v4uni.LocRib().Route(prefix, oc.UnionString(dut1.RouterID), 0).Prefix().State())
		// Check route is not in Adj-Out of dut2.
		awaitNotPresent(t, dut2, v4uni.Neighbor(dut3.RouterID).AdjRibOutPre().Route(prefix, 0).Prefix().State())
		awaitNotPresent(t, dut2, v4uni.Neighbor(dut3.RouterID).AdjRibOutPost().Route(prefix, 0).Prefix().State())

		// Check route is not in Adj-In of dut3.
		awaitNotPresent(t, dut3, v4uni.Neighbor(dut2.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State())

		// Check route is not in Adj-In of dut4.
		awaitNotPresent(t, dut4, v4uni.Neighbor(dut3.RouterID).AdjRibInPre().Route(prefix, 0).Prefix().State())
	}

	t.Run("static", func(t *testing.T) {
		checkpath(t, "10.10.10.0/24", true)
	})
	t.Run("gribi", func(t *testing.T) {
		checkpath(t, "20.20.20.0/24", false)
	})
}
