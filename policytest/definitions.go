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

package policytest

import "github.com/openconfig/lemming/gnmi/oc"

type RouteTestCase struct {
	Input                   TestRoute
	RouteTest               *RoutePathTestCase
	AlternatePathRouteTest  *RoutePathTestCase
	LongerPathRouteTest     *RoutePathTestCase
	SkipDUT2RouteValidation bool
}

type RoutePathTestCase struct {
	Description    string
	ExpectedResult RouteTestResult

	PrevAdjRibOutPreCommunities []string
	// An export policy is applied here by convention.
	PrevAdjRibOutPostCommunities []string

	AdjRibInPreCommunities []string
	// An import policy is applied here by convention.
	AdjRibInPostCommunities  []string
	LocalRibCommunities      []string
	AdjRibOutPreCommunities  []string
	AdjRibOutPostCommunities []string

	NextAdjRibInPreCommunities []string
	NextLocalRibCommunities    []string

	PrevAdjRibOutPreAttrs  *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet
	PrevAdjRibOutPostAttrs *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet
	AdjRibInPreAttrs       *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet
	AdjRibInPostAttrs      *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet
	LocalRibAttrs          *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet
	AdjRibOutPreAttrs      *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet
	AdjRibOutPostAttrs     *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet
	NextAdjRibInPreAttrs   *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet
	NextAdjRibInPostAttrs  *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet
	NextLocalRibAttrs      *oc.NetworkInstance_Protocol_Bgp_Rib_AttrSet
}

// This message represents a single prefix and its associated BGP attributes.
type TestRoute struct {
	ReachPrefix string
}

// The expected result for a RouteTestCase
type RouteTestResult int

const (
	// RouteUnspecified should never be the value of an actual test case.
	RouteUnspecified RouteTestResult = iota
	// RouteAccepted means to accept the input TestRoute.
	RouteAccepted
	// RouteDiscarded means to discard the input TestRoute.
	RouteDiscarded
	// RouteNotPreferred means not selected by best path selection.
	RouteNotPreferred
)
