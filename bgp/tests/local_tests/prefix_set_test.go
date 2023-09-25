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

	valpb "github.com/openconfig/lemming/proto/policyval"
)

func TestPrefixSetMode(t *testing.T) {
	dut1, stop1 := newLemming(t, 1, 64500, []*AddIntfAction{{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "192.0.2.1/31",
		niName:  "DEFAULT",
	}})
	defer stop1()

	prefix1 := "10.33.0.0/16"
	prefix2 := "10.34.0.0/16"
	prefix2v6 := "10::34/16"

	// Create prefix set
	prefixSetName := "reject-" + prefix1
	prefixSetPath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName)
	Replace(t, dut1, prefixSetPath.Mode().Config(), oc.PrefixSet_Mode_IPV4)
	Replace(t, dut1, prefixSetPath.Prefix(prefix1, "exact").IpPrefix().Config(), prefix1)
	ReplaceExpectFail(t, dut1, prefixSetPath.Prefix(prefix2v6, "exact").IpPrefix().Config(), prefix2v6)
	Replace(t, dut1, prefixSetPath.Prefix(prefix2, "exact").IpPrefix().Config(), prefix2)
}

func TestPrefixSet(t *testing.T) {
	installPolicies := func(t *testing.T, dut1, dut2, _, _, _ *Device, invert bool) {
		if debug {
			fmt.Println("Installing test policies")
		}
		prefix1 := "10.33.0.0/16"
		prefix2 := "10.34.0.0/16"

		// Policy to reject routes with the given prefix set
		policyName := "def1"

		// Create prefix set
		prefixSetName := "reject-" + prefix1
		prefixSetPath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName)
		Replace(t, dut2, prefixSetPath.Mode().Config(), oc.PrefixSet_Mode_IPV4)
		prefix1Path := prefixSetPath.Prefix(prefix1, "exact").IpPrefix()
		Replace(t, dut2, prefix1Path.Config(), prefix1)
		prefix2Path := prefixSetPath.Prefix(prefix2, "16..23").IpPrefix()
		Replace(t, dut2, prefix2Path.Config(), prefix2)

		policy := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
		stmt, err := policy.AppendNew("stmt1")
		if err != nil {
			t.Fatalf("Cannot append new BGP policy statement: %v", err)
		}
		// Match on prefix set & reject route
		stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
		if invert {
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_INVERT)
		} else {
			stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.RoutingPolicy_MatchSetOptionsRestrictedType_ANY)
		}
		stmt.GetOrCreateActions().SetPolicyResult(oc.RoutingPolicy_PolicyResultType_REJECT_ROUTE)
		// Install policy
		Replace(t, dut2, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config(), &oc.RoutingPolicy_PolicyDefinition{Statement: policy})
		Replace(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).ApplyPolicy().ImportPolicy().Config(), []string{policyName})
	}

	invertResult := func(result valpb.RouteTestResult, invert bool) valpb.RouteTestResult {
		if invert {
			switch result {
			case valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT:
				return valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD
			case valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD:
				return valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT
			default:
			}
		}
		return result
	}

	getspec := func(invert bool) *valpb.PolicyTestCase {
		spec := &valpb.PolicyTestCase{
			Description: "Test that one prefix gets accepted and the other rejected via an ANY prefix-set.",
			RouteTests: []*valpb.RouteTestCase{{
				Description: "Exact match",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.33.0.0/16",
				},
				ExpectedResult: invertResult(valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD, invert),
			}, {
				Description: "Not exact match",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.33.0.0/17",
				},
				ExpectedResult: invertResult(valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT, invert),
			}, {
				Description: "No match with any prefix",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.3.0.0/16",
				},
				ExpectedResult: invertResult(valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT, invert),
			}, {
				Description: "mask length too short",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.34.0.0/15",
				},
				ExpectedResult: invertResult(valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT, invert),
			}, {
				Description: "Lower end of mask length",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.34.0.0/16",
				},
				ExpectedResult: invertResult(valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD, invert),
			}, {
				Description: "Middle of mask length",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.34.0.0/20",
				},
				ExpectedResult: invertResult(valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD, invert),
			}, {
				Description: "Upper end of mask length",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.34.0.0/23",
				},
				ExpectedResult: invertResult(valpb.RouteTestResult_ROUTE_TEST_RESULT_DISCARD, invert),
			}, {
				Description: "mask length too long",
				Input: &valpb.TestRoute{
					ReachPrefix: "10.34.0.0/24",
				},
				ExpectedResult: invertResult(valpb.RouteTestResult_ROUTE_TEST_RESULT_ACCEPT, invert),
			}},
		}
		return spec
	}

	t.Run("ANY", func(t *testing.T) {
		testPolicy(t, PolicyTestCase{
			spec: getspec(false),
			installPolicies: func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device) {
				installPolicies(t, dut1, dut2, dut3, dut4, dut5, false)
			},
		})
	})
	t.Run("INVERT", func(t *testing.T) {
		testPolicy(t, PolicyTestCase{
			spec: getspec(true),
			installPolicies: func(t *testing.T, dut1, dut2, dut3, dut4, dut5 *Device) {
				installPolicies(t, dut1, dut2, dut3, dut4, dut5, true)
			},
		})
	})
}
