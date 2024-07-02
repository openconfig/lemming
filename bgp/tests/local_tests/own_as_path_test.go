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
	"fmt"
	"testing"

	"github.com/openconfig/lemming/bgp"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/policytest"
	"github.com/openconfig/ygot/ygot"
)

//
// dut1 RIB --> export policy (Set the AS Path) --> BGP session --> AS path checker --> Import Policy (default accept) --> dut2 RIB
//

func TestOwnASPath(t *testing.T) {
	// 10.* are subject to rejection due to community-count policy.
	// 20.* must propagate in all cases due to prefix-match rule.
	routeUnderTestList := map[int]string{
		0: "10.0.0.0/10",
		1: "11.0.0.0/10",
		2: "12.0.0.0/10",
		3: "13.0.0.0/10",
		4: "14.0.0.0/10",
		5: "15.0.0.0/10",
	}

	installASPathPrepend := func(t *testing.T, dut1, dut2, _, _, _ *Device, num uint8) {
		Update(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).AsPathOptions().AllowOwnAs().Config(), num)
		awaitSessionEstablished(t, dut2, dut1)
	}

	installSetPolicy := func(t *testing.T, dut1, dut2, _, _, _ *Device) {
		if debug {
			fmt.Println("Installing test policies")
		}

		policyName := "set-as-path"
		policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}

		for i, route := range routeUnderTestList {
			// Create prefix set
			prefixSetName := "accept-" + route
			prefixSetPath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName)
			Replace(t, dut1, prefixSetPath.Mode().Config(), oc.PrefixSet_Mode_IPV4)
			Replace(t, dut1, prefixSetPath.Prefix(route, "exact").IpPrefix().Config(), route)

			stmt, err := policy.AppendNew(fmt.Sprintf("stmt%d", i))
			if err != nil {
				t.Fatalf("Cannot append new BGP policy statement: %v", err)
			}
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)

			if i > 0 {
				asPathPrepend := stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetAsPathPrepend()
				asPathPrepend.SetRepeatN(uint8(i))
				asPathPrepend.SetAsn(dut2.AS)
			}
			stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE)
		}
		// Install policy
		Replace(t, dut1, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(policyName), Statement: policy})
		Replace(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().Config(), []string{policyName})
		Await(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).ApplyPolicy().ExportPolicy().State(), []string{policyName})
	}

	for i := 0; i <= 5; i++ {

		var routeTests []*policytest.RouteTestCase
		for routeId, prefix := range routeUnderTestList {
			result := policytest.RouteAccepted
			if routeId > i {
				result = policytest.RouteDiscarded
			}
			routeTests = append(routeTests, &policytest.RouteTestCase{
				Input: policytest.TestRoute{
					ReachPrefix: prefix,
				},
				RouteTest: &policytest.RoutePathTestCase{
					Description:    fmt.Sprintf("%d", routeId),
					ExpectedResult: result,
				},
				SkipDUT2RouteValidation: true,
			})
		}

		testPolicy(t, &PolicyTestCase{
			description:         fmt.Sprintf("own-as-path-%d", i),
			skipValidateAttrSet: true,
			routeTests:          routeTests,
			dut1IsEBGP:          true,
			installPolicies: func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device) {
				installSetPolicy(t, dut1, dut2, dut3, dut4, dut5)
				installASPathPrepend(t, dut1, dut2, dut3, dut4, dut5, uint8(i))
			},
		})
	}
}
