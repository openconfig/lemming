// Copyright 2021 Google LLC
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

// Package sysrib implements a system-level RIB that is populated initially using
// an OpenConfig configuration.
package sysrib

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"sync"

	"github.com/kentik/patricia"
	"github.com/kentik/patricia/generics_tree"
	"github.com/openconfig/gribigo/afthelper"
	oc "github.com/openconfig/gribigo/ocrt"
	"github.com/openconfig/ygot/ytypes"
)

// SysRIB is a RIB data structure that can be used to resolve routing entries to their egress interfaces.
// Currently it supports only IPV4 entries.
type SysRIB struct {
	// mu protects the map of network instance RIBs.
	mu sync.RWMutex
	// NI is the list of network instances (aka VRFs)
	NI        map[string]*NIRIB
	defaultNI string
}

// NIRIB is the RIB for a single network instance.
type NIRIB struct {
	// IPV4 is the IPv4 RIB
	IPV4 *generics_tree.TreeV4[*Route]
}

type RoutePreference struct {
	// AdminDistance is the admin distance of the protocol that added this
	// route.
	AdminDistance uint8 `json:"admin-distance"`
	// Metric is the metric of the route. It is comparable only within
	// routes of the same protocol, and therefore the same admin distance.
	Metric uint32 `json:"metric"`
}

// Route is used to store a route in the radix tree.
type Route struct {
	// Prefix is a prefix that was being stored.
	Prefix string `json:"prefix"`
	// Connected indicates that the route is directly connected.
	Connected *Interface `json:"connected"`
	// NextHops is the set of IP nexthops that the route uses if
	// it is not a connected route.
	NextHops  []*afthelper.NextHopSummary `json:"nexthops"`
	RoutePref RoutePreference
}

// NewSysRIB returns a SysRIB from an input parsed OpenConfig configuration.
func NewSysRIB(cfg *oc.Device) (*SysRIB, error) {
	sr := &SysRIB{
		NI: map[string]*NIRIB{},
	}

	if cfg != nil {
		cr, err := connectedRoutesFromConfig(cfg)
		if err != nil {
			return nil, err
		}

		for ni, niR := range cr {
			sr.NI[ni] = &NIRIB{
				IPV4: generics_tree.NewTreeV4[*Route](),
			}
			if niR.T == oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE {
				sr.defaultNI = ni
			}
			for _, r := range niR.Rts {
				if err := sr.AddRoute(ni, r); err != nil {
					return nil, err
				}
			}
		}
	} else {
		sr.defaultNI = "DEFAULT"
		sr.NI[sr.defaultNI] = &NIRIB{
			IPV4: generics_tree.NewTreeV4[*Route](),
		}
	}

	return sr, nil
}

// AddRoute adds a route, r, to the network instance, ni, in the sysRIB.
// It returns an error if it cannot be added.
func (sr *SysRIB) AddRoute(ni string, r *Route) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	if _, ok := sr.NI[ni]; !ok {
		return fmt.Errorf("cannot find network instance %s", ni)
	}
	addr, _, err := patricia.ParseIPFromString(r.Prefix)
	if err != nil {
		return fmt.Errorf("cannot create prefix for %s, %v", r.Prefix, err)
	}
	if added, _ := sr.NI[ni].IPV4.Add(*addr, r, nil); !added {
		return fmt.Errorf("cannot insert route in network instance %s %s", ni, r.Prefix)
	}
	return nil
}

// NewRoute returns a new route for the specified prefix.
// Note - today this doesn't actually result in a viable
// forwarding entry unless its a connected route :-)
func NewRouteViaIF(pfx string, intf *Interface) *Route {
	return &Route{Prefix: pfx, Connected: intf}
}

// NewSysRIBFromJSON returns a new SysRIB from an RFC7951 marshalled JSON OpenConfig configuration.
func NewSysRIBFromJSON(jsonCfg []byte) (*SysRIB, error) {
	cfg := &oc.Device{}
	if err := oc.Unmarshal(jsonCfg, cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal JSON configuration, %v", err)
	}
	return NewSysRIB(cfg)
}

// Interface describes an interface of a device.
type Interface struct {
	Name         string `json:"name"`
	Index        int32  `json:"index"`
	Subinterface uint32 `json:"subinterface"`
}

// entryForCIDR returns the RIB entry for the IP address specified by ip within
// the specified network instance. It returns a bool indicating whether the
// entry was found, a slice of strings which contains its tags, and an optional
// error.
func (sr *SysRIB) entryForCIDR(ni string, ip *net.IPNet) (bool, []*Route, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	rib, ok := sr.NI[ni]
	if !ok {
		return false, nil, fmt.Errorf("cannot find a RIB for network instance %s", ni)
	}
	addr, _, err := patricia.ParseFromIPAddr(ip)
	if err != nil {
		return false, nil, fmt.Errorf("cannot parse IP to lookup, %s: %v", ip, err)
	}
	found, tags := rib.IPV4.FindDeepestTags(*addr)
	return found, tags, nil
}

// EgressInterface looks up the IP destination address ip in the routes for network instance
// named inputNI. It returns a slice of the interfaces that the packet would be forwarded
// via.
//
// TODO(robjs): support determining the NI based solely on the input interface.
// TODO(robjs): support a better description of a packet using the formats that ONDATRA uses.
//
// TODO(robjs): support WCMP
//
// This is really a POC that we can emulate our FIB for basic IPV4 routes.
func (sr *SysRIB) EgressInterface(inputNI string, ip *net.IPNet) ([]*Interface, error) {
	// no RIB recursion currently
	if inputNI == "" {
		inputNI = sr.defaultNI
	}

	found, routes, err := sr.entryForCIDR(inputNI, ip)
	if err != nil {
		return nil, fmt.Errorf("cannot lookup IP %s", ip)
	}

	if !found {
		return nil, nil
	}

	egressIfs := []*Interface{}
	for _, cr := range routes {
		if cr.Connected != nil {
			egressIfs = append(egressIfs, cr.Connected)
			continue
		}

		// This isn't a connected route, check whether we can resolve the next-hops.
		for _, nh := range cr.NextHops {
			_, nhop, err := net.ParseCIDR(fmt.Sprintf("%s/32", nh.Address))
			if err != nil {
				return nil, fmt.Errorf("can't parse %s/32 into CIDR, %v", nh.Address, err)
			}
			recursiveNHIfs, err := sr.EgressInterface(nh.NetworkInstance, nhop)
			if err != nil {
				return nil, fmt.Errorf("for nexthop %s, can't resolve: %v", nh.Address, err)
			}
			egressIfs = append(egressIfs, recursiveNHIfs...)
		}
	}
	return egressIfs, nil
}

// EgressNexthops returns the resolved nexthops for the input IP prefix for
// network instance inputNI based on the device's interface state.
func (sr *SysRIB) EgressNexthops(inputNI string, ip *net.IPNet, interfaces map[Interface]bool) (map[ResolvedNexthop]bool, error) {
	// no RIB recursion currently
	if inputNI == "" {
		inputNI = sr.defaultNI
	}

	found, routes, err := sr.entryForCIDR(inputNI, ip)
	if err != nil {
		return nil, fmt.Errorf("cannot lookup IP %s", ip)
	}

	if !found {
		return nil, nil
	}

	// For each route entry for the prefix, recursively resolve their nexthops.
	// Then, select the set of resolved nexthops for the route entry according to the following preference:
	//   1. Has at least one enabled & connected nexthop after resolution.
	//   2. Lowest admin distance.
	//   3. Lowest metric.
	// When there is a tie, use regular ECMP/WCMP rules.
	//
	// TODO(wenbli): Support WCMP.
	allEgressNhs := map[RoutePreference]map[ResolvedNexthop]bool{}
	for _, cr := range routes {
		if allEgressNhs[cr.RoutePref] == nil {
			allEgressNhs[cr.RoutePref] = map[ResolvedNexthop]bool{}
		}
		egressNhs := allEgressNhs[cr.RoutePref]
		if cr.Connected != nil {
			if interfaces[*cr.Connected] {
				nh := ResolvedNexthop{
					Port: *cr.Connected,
					NextHopSummary: afthelper.NextHopSummary{
						NetworkInstance: inputNI,
					},
				}
				if length, _ := ip.Mask.Size(); length == 32 {
					nh.Address = ip.IP.String()
				}
				// TODO(wenbli): Implement WCMP: there could be a merger of two nexthops, in which case we add their weights.
				egressNhs[nh] = true
			}
			continue
		}

		// This isn't a connected route, check whether we can resolve the next-hops.
		for _, nh := range cr.NextHops {
			_, nhop, err := net.ParseCIDR(fmt.Sprintf("%s/32", nh.Address))
			if err != nil {
				return nil, fmt.Errorf("can't parse %s/32 into CIDR, %v", nh.Address, err)
			}
			recursiveNHs, err := sr.EgressNexthops(nh.NetworkInstance, nhop, interfaces)
			if err != nil {
				return nil, fmt.Errorf("for nexthop %s, can't resolve: %v", nh.Address, err)
			}
			for nh := range recursiveNHs {
				// TODO(wenbli): Implement WCMP: there could be a merger of two nexthops, in which case we add their weights.
				egressNhs[nh] = true
			}
		}
	}

	var allRoutePrefs []RoutePreference
	for rp := range allEgressNhs {
		allRoutePrefs = append(allRoutePrefs, rp)
	}

	sort.Slice(allRoutePrefs, func(i, j int) bool {
		return allRoutePrefs[i].AdminDistance < allRoutePrefs[j].AdminDistance || allRoutePrefs[i].AdminDistance == allRoutePrefs[j].AdminDistance && allRoutePrefs[i].Metric < allRoutePrefs[j].Metric
	})

	for _, rp := range allRoutePrefs {
		if len(allEgressNhs[rp]) != 0 {
			return allEgressNhs[rp], nil
		}
	}

	return nil, nil
}

// niConnected is a description of a set of connected routes within a network instance.
type niConnected struct {
	// N is the network instance to which the route belongs.
	N string
	// t is the type of netowrk instance.
	T oc.E_NetworkInstanceTypes_NETWORK_INSTANCE_TYPE
	// rts is the set of connected routes within the network instance.
	Rts []*Route
}

// connectedRoutesFromConfig returns the set of 'connected' routes from the input configuration supplied.
// Connected routes are defined to be those that are directly configured as a subnet to which the
// system is attached.
//
// This function only returns connected IPV4 routes.
func connectedRoutesFromConfig(cfg *oc.Device) (map[string]*niConnected, error) {
	// TODO(robjs): figure out where the reference that is referencing policy
	// definitions is that has not yet been removed, improve ygot error message.
	if err := cfg.Validate(&ytypes.LeafrefOptions{
		IgnoreMissingData: true,
		Log:               true,
	}); err != nil {
		return nil, fmt.Errorf("invalid input configuration, %v", err)
	}

	matched := map[string]map[uint32]bool{}
	// intfRoute is a map, keyed by the name of a physical interface, of maps, keyed by the id
	// of a subinterface, that points to the set of connected routes that are configured on the
	// interface.
	intfRoute := map[string]map[uint32][]*Route{}
	for intName, intf := range cfg.Interface {
		intfRoute[intf.GetName()] = map[uint32][]*Route{}
		for subIntIdx, subintf := range intf.Subinterface {
			if subintf.GetIpv4() != nil {
				for _, a := range subintf.GetIpv4().Address {
					_, cidr, err := net.ParseCIDR(fmt.Sprintf("%s/%d", a.GetIp(), a.GetPrefixLength()))
					if err != nil {
						return nil, fmt.Errorf("invalid IPV4 prefix on interface %s, subinterface %d, %s/%d", intf.GetName(), subintf.GetIndex(), a.GetIp(), a.GetPrefixLength())
					}
					rt := &Route{
						Prefix: cidr.String(),
						Connected: &Interface{
							Name:         intf.GetName(),
							Subinterface: subintf.GetIndex(),
						},
					}
					intfRoute[intName][subIntIdx] = append(intfRoute[intName][subIntIdx], rt)
					if matched[intf.GetName()] == nil {
						matched[intf.GetName()] = map[uint32]bool{}
					}
					matched[intf.GetName()][subintf.GetIndex()] = false
				}
			}
		}
	}

	var (
		defName string
		ni      = make(map[string]*niConnected)
	)

	for _, n := range cfg.NetworkInstance {
		netInstRoutes := &niConnected{
			N: n.GetName(),
		}

		// We don't support L2 adjacencies.
		switch n.GetType() {
		case oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE:
			if defName != "" {
				return nil, fmt.Errorf("cannot have >1 default instance, got %s and %s", n.GetName(), defName)
			}
			defName = n.GetName()
			netInstRoutes.T = n.GetType()
		case oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_L2P2P, oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_L2VSI, oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_UNSET:
			return nil, fmt.Errorf("invalid network instance type specified for NI %s, %s", n.GetName(), n.GetType())
		default:
			netInstRoutes.T = n.GetType()
		}

		for _, i := range n.Interface {
			if i.Subinterface == nil {
				// an L3 adjacency can only be associated with a subinterface in openconfig.
				continue
			}
			if intfRoute[i.GetInterface()] != nil && intfRoute[i.GetInterface()][i.GetSubinterface()] != nil {
				netInstRoutes.Rts = append(netInstRoutes.Rts, intfRoute[i.GetInterface()][i.GetSubinterface()]...)
				matched[i.GetInterface()][i.GetSubinterface()] = true
			}
		}

		sort.Slice(netInstRoutes.Rts, func(i, j int) bool {
			return netInstRoutes.Rts[i].Prefix < netInstRoutes.Rts[j].Prefix
		})

		ni[n.GetName()] = netInstRoutes
	}

	if defName == "" {
		return nil, errors.New("no default network instance, invalid")
	}

	for intfName, i := range intfRoute {
		for subintIndex, routes := range i {
			if !matched[intfName][subintIndex] {
				// any unmatched interface is mapped to the default network instance.
				ni[defName].Rts = append(ni[defName].Rts, routes...)
			}
		}
	}

	return ni, nil
}
