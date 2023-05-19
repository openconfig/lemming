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
	"fmt"
	"testing"

	"github.com/openconfig/lemming/bgp"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"

	valpb "github.com/openconfig/lemming/bgp/tests/proto/policyval"
)

func TestPrefixSet(t *testing.T) {
	testspec := PolicyTestCase{
		spec: &valpb.PolicyTestCase{
			Description: "Test that one prefix gets accepted and the other rejected via a prefix-set.",
			RouteTest: []*valpb.RouteTestCase{{
				Description: "Rejected route",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.33.0.0/16",
				},
				ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD,
			}, {
				Description: "Accepted route",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.3.0.0/16",
				},
				ExpectedResult: valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT,
			}},
		},
		installPolicy: func(t *testing.T, dut2 *ygnmi.Client) {
			if debug {
				fmt.Println("Installing test policies")
			}
			prefix1 := "10.33.0.0/16"

			// Custom policy
			prefixSetName := "reject-" + prefix1
			prefix1Path := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName).Prefix(prefix1, "exact").IpPrefix()
			Replace(t, dut2, prefix1Path.Config(), prefix1)

			policyName := "def1"
			Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Statement("stmt1").Conditions().MatchPrefixSet().PrefixSet().Config(), prefixSetName)
			Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Statement("stmt1").Conditions().MatchPrefixSet().MatchSetOptions().Config(), oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)
			Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Statement("stmt1").Actions().PolicyResult().Config(), oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)
			Replace(t, dut2, bgp.BGPPath.Neighbor(dut3spec.RouterID).ApplyPolicy().ExportPolicy().Config(), []string{policyName})
		},
	}

	testPolicy(t, testspec, false)
	testPolicy(t, testspec, true)
}
