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
	"net/netip"
	"reflect"
	"sort"
	"sync"

	log "github.com/golang/glog"
	"github.com/kentik/patricia"
	"github.com/kentik/patricia/generics_tree"
	"github.com/openconfig/gribigo/afthelper"
	oc "github.com/openconfig/gribigo/ocrt"
	"github.com/openconfig/lemming/gnmi/fakedevice"
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

	// GUEPolicies are the configured BGP-triggered GUE policies.
	// Every update to this should trigger a re-computation of the resolved
	// routes.
	// TODO(wenbli): Support v6 GUE policies. When adding support, make
	// sure to error out when an IPv6-mapped IPv4 prefix is provided.
	GUEPoliciesV4 *generics_tree.TreeV4[GUEPolicy]
	GUEPoliciesV6 *generics_tree.TreeV6[GUEPolicy]
}

// NIRIB is the RIB for a single network instance.
type NIRIB struct {
	// IPV4 is the IPv4 RIB
	IPV4 *generics_tree.TreeV4[*Route]
	// IPV6 is the IPv6 RIB
	IPV6 *generics_tree.TreeV6[*Route]
}

// GUEPolicy represents the static values in the IP and UDP headers that are
// to encapsulate the packet.
//
// srcPort is a hash.
// dstIP is the nexthop of the BGP route.
type GUEPolicy struct {
	// dstPortv4 is the UDP port used when the packet payload is IPv4.
	dstPortv4 uint16
	// dstPortv6 is the UDP port used when the packet payload is IPv6.
	dstPortv6 uint16
	srcIP4    [4]byte
	srcIP6    [16]byte
	isV6      bool
}

// GUEHeaders represents the IP and UDP headers that are to encapsulate the
// packet.
type GUEHeaders struct {
	GUEPolicy
	dstIP4 [4]byte
	dstIP6 [16]byte
}

// getGUEHeader retrieves the GUEHeader for the given address if it matched a
// GUE policy. The boolean return indicates whether a GUE policy was
// matched.
func (sr *SysRIB) getGUEHeader(address string) (GUEHeaders, bool, error) {
	addr4, addr6, err := patricia.ParseIPFromString(address)
	if err != nil {
		return GUEHeaders{}, false, err
	}
	var ok bool
	var policy GUEPolicy
	switch {
	case addr4 != nil:
		ok, policy = sr.GUEPoliciesV4.FindDeepestTag(*addr4)
	case addr6 != nil:
		ok, policy = sr.GUEPoliciesV6.FindDeepestTag(*addr6)
	default:
		return GUEHeaders{}, false, fmt.Errorf("Invalid IP address for looking up GUE header")
	}
	if !ok {
		return GUEHeaders{}, false, nil
	}
	ip := net.ParseIP(address)
	if ip == nil {
		return GUEHeaders{}, false, fmt.Errorf("cannot parse IP address: %q", address)
	}
	if policy.isV6 {
		var dstIP6 [16]byte
		for i, octet := range ip.To16() {
			dstIP6[i] = octet
		}
		return GUEHeaders{
			GUEPolicy: policy,
			dstIP6:    dstIP6,
		}, true, nil
	}
	var dstIP4 [4]byte
	for i, octet := range ip.To4() {
		dstIP4[i] = octet
	}
	return GUEHeaders{
		GUEPolicy: policy,
		dstIP4:    dstIP4,
	}, true, nil
}

type RoutePreference struct {
	// AdminDistance is the admin distance of the protocol that added this
	// route.
	// See https://docs.frrouting.org/en/latest/zebra.html#administrative-distance
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

func (r *Route) String() string {
	readable := fmt.Sprintf("%s (%+v)", r.Prefix, r.RoutePref)
	switch {
	case r.Connected != nil:
		return fmt.Sprintf("%s: connected interface %+v", readable, *r.Connected)
	default:
		readable += ": ["
		for i, nh := range r.NextHops {
			if i != 0 {
				readable += ", "
			}
			readable += fmt.Sprintf("nexthop %d: ", i)
			if nh == nil {
				readable += "nil nexthop"
			}
			readable += fmt.Sprintf("%+v", nh)
		}
		readable += "]"
		return readable
	}
}

// NewSysRIB returns a SysRIB from an input parsed OpenConfig configuration.
func NewSysRIB(cfg *oc.Device) (*SysRIB, error) {
	sr := &SysRIB{
		NI:            map[string]*NIRIB{},
		GUEPoliciesV4: generics_tree.NewTreeV4[GUEPolicy](),
		GUEPoliciesV6: generics_tree.NewTreeV6[GUEPolicy](),
	}

	if cfg != nil {
		cr, err := connectedRoutesFromConfig(cfg)
		if err != nil {
			return nil, err
		}

		for ni, niR := range cr {
			sr.NI[ni] = &NIRIB{
				IPV4: generics_tree.NewTreeV4[*Route](),
				IPV6: generics_tree.NewTreeV6[*Route](),
			}
			if niR.T == oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE {
				sr.defaultNI = ni
			}
			for _, r := range niR.Rts {
				if _, err := sr.AddRoute(ni, r); err != nil {
					return nil, err
				}
			}
		}
	} else {
		sr.defaultNI = fakedevice.DefaultNetworkInstance
		sr.NI[sr.defaultNI] = &NIRIB{
			IPV4: generics_tree.NewTreeV4[*Route](),
			IPV6: generics_tree.NewTreeV6[*Route](),
		}
	}

	return sr, nil
}

// routeMatches return if two routes are equal.
func routeMatches(a *Route, b *Route) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return reflect.DeepEqual(*a, *b)
}

// guePolicyUnconditionalMatch always returns true. It is intended to be used
// for kentik/patricia's matchFunc argument for Delete.
func guePolicyUnconditionalMatch(_ GUEPolicy, _ GUEPolicy) bool {
	return true
}

// AddRoute adds a route, r, to the network instance, ni, in the sysRIB.
// It returns true if the route was added, and false if not. If the route
// already exists, it returns (false, nil)
func (sr *SysRIB) AddRoute(ni string, r *Route) (bool, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	if _, ok := sr.NI[ni]; !ok {
		return false, fmt.Errorf("cannot find network instance %s", ni)
	}
	prefix, err := canonicalPrefix(r.Prefix)
	if err != nil {
		return false, fmt.Errorf("sysrib: prefix cannot be parsed: %v", err)
	}
	addr4, addr6, err := patricia.ParseIPFromString(prefix.String())
	if err != nil {
		return false, fmt.Errorf("cannot create prefix for %s, %v", r.Prefix, err)
	}
	switch {
	case addr4 != nil:
		added, _ := sr.NI[ni].IPV4.Add(*addr4, r, routeMatches)
		log.V(1).Infof("AddRoute attempt: %v, %v, result: %v", *addr4, r, added)
		return added, nil
	case addr6 != nil:
		added, _ := sr.NI[ni].IPV6.Add(*addr6, r, routeMatches)
		log.V(1).Infof("AddRoute attempt: %v, %v, result: %v", *addr6, r, added)
		return added, nil
	default:
		return false, fmt.Errorf("route prefix is neither v4 or v6: %v", r.Prefix)
	}
}

// SetGUEPolicy sets a GUE Policy in the RIB.
func (sr *SysRIB) SetGUEPolicy(prefix string, policy GUEPolicy) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	addr4, addr6, err := patricia.ParseIPFromString(prefix)
	if err != nil {
		return fmt.Errorf("cannot create prefix for %s, %v", prefix, err)
	}
	switch {
	case addr4 != nil:
		sr.GUEPoliciesV4.Set(*addr4, policy)
	case addr6 != nil:
		sr.GUEPoliciesV6.Set(*addr6, policy)
	default:
		return fmt.Errorf("SetGUEPolicy: invalid prefix: %v", prefix)
	}
	return nil
}

// DeleteGUEPolicy sets a GUE Policy in the RIB.
// It returns true if a policy was deleted, and false if not.
func (sr *SysRIB) DeleteGUEPolicy(prefix string) (bool, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	addr4, addr6, err := patricia.ParseIPFromString(prefix)
	if err != nil {
		return false, fmt.Errorf("cannot create prefix for %s, %v", prefix, err)
	}
	var count int
	switch {
	case addr4 != nil:
		count = sr.GUEPoliciesV4.Delete(*addr4, guePolicyUnconditionalMatch, GUEPolicy{})
	case addr6 != nil:
		count = sr.GUEPoliciesV6.Delete(*addr6, guePolicyUnconditionalMatch, GUEPolicy{})
	default:
		return false, fmt.Errorf("DeleteGUEPolicy: invalid prefix: %v", prefix)
	}
	return count > 0, nil
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
	addr4, addr6, err := patricia.ParseFromIPAddr(ip)
	if err != nil {
		return false, nil, fmt.Errorf("cannot parse IP to lookup, %s: %v", ip, err)
	}
	switch {
	case addr4 != nil:
		found, tags := rib.IPV4.FindDeepestTags(*addr4)
		return found, tags, nil
	case addr6 != nil:
		found, tags := rib.IPV6.FindDeepestTags(*addr6)
		return found, tags, nil
	default:
		return false, nil, fmt.Errorf("route prefix is neither v4 or v6: %v", ip)
	}
}

// addressToPrefix returns a prefix of /32 or /128 of the input v4 or v6 address.
//
// It returns an error if the input address cannot be parsed.
//
// e.g. 1.1.1.1 -> 1.1.1.1/32
// e.g. 2001:: -> 2001::/128
func addressToPrefix(address string) (*net.IPNet, error) {
	addr, err := netip.ParseAddr(address)
	if err != nil {
		return nil, fmt.Errorf("sysrib.addressToPrefix: cannot parse address: %v", address)
	}
	mask := 32
	if addr.Is6() {
		mask = 128
	}
	_, nhop, err := net.ParseCIDR(fmt.Sprintf("%s/%d", address, mask))
	if err != nil {
		return nil, fmt.Errorf("can't parse %s/%d into CIDR, %v", address, mask, err)
	}
	return nhop, nil
}

// EgressInterface looks up the IP destination address ip in the routes for network instance
// named inputNI. It returns a slice of the interfaces that the packet would be forwarded
// via.
//
// TODO(robjs): support determining the NI based solely on the input interface.
// TODO(robjs): support a better description of a packet using the formats that ONDATRA uses.
//
// TODO(robjs): support WCMP
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
			nhop, err := addressToPrefix(nh.Address)
			if err != nil {
				return nil, err
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

// EgressNexthops returns the resolved nexthops for the input IP prefix. It
// also returns the top-level route (at this level) that was successfully
// resolved (if any). This is useful for determining the properties of the
// route that was ultimately resolved, for example its route preference and
// first-level nexthops.
//
// - inputNI is the network instance of the input prefix.
// - interfaces is the set of known interface states on the device.
func (sr *SysRIB) EgressNexthops(inputNI string, ip *net.IPNet, interfaces map[Interface]bool) (map[ResolvedNexthop]bool, *Route, error) {
	// no RIB recursion currently
	if inputNI == "" {
		inputNI = sr.defaultNI
	}

	found, routes, err := sr.entryForCIDR(inputNI, ip)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot lookup IP %s", ip)
	}

	if !found {
		log.V(1).Infof("Prefix not found in RIB: %v", ip)
		return nil, nil, nil
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
	resolvedRoutes := map[RoutePreference]*Route{}
	for _, cr := range routes {
		if allEgressNhs[cr.RoutePref] == nil {
			allEgressNhs[cr.RoutePref] = map[ResolvedNexthop]bool{}
			resolvedRoutes[cr.RoutePref] = cr
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
				if length, _ := ip.Mask.Size(); (ip.IP.To4() != nil && length == 32) || (ip.IP.To16() != nil && length == 128) {
					nh.Address = ip.IP.String()
				}
				// TODO(wenbli): Implement WCMP: there could be a merger of two nexthops, in which case we add their weights.
				egressNhs[nh] = true
			}
			continue
		}

		// This isn't a connected route, check whether we can resolve the next-hops.
		for _, nh := range cr.NextHops {
			nhop, err := addressToPrefix(nh.Address)
			if err != nil {
				return nil, nil, err
			}
			recursiveNHs, _, err := sr.EgressNexthops(nh.NetworkInstance, nhop, interfaces)
			if err != nil {
				return nil, nil, fmt.Errorf("for nexthop %s, can't resolve: %v", nh.Address, err)
			}

			// pseudocode for BGP-triggered GUE:
			// if route is BGP, then
			// - For each of its nexthops,
			//   - if the nexthop falls within the policy prefix, and
			//   - the same nexthop doesn't also already have a BGP-triggered GUE action
			//   then, add the route with an encap action for the nexthop.
			var encapHeaders GUEHeaders
			if cr.RoutePref.AdminDistance == 20 { // EBGP
				gueHeaders, ok, err := sr.getGUEHeader(nh.Address)
				if ok {
					encapHeaders = gueHeaders
				}
				if err != nil {
					return nil, nil, fmt.Errorf("Error during GUE policy look-up: %v", err)
				}
			}
			for rnh := range recursiveNHs {
				switch {
				case rnh.HasGUE():
					return nil, nil, fmt.Errorf("route %v resolves over another route that has a BGP-triggered GUE action, the behaviour is undefined, nexthop: %v, recursive nexthop: %v", cr, nh, rnh)
				}
				rnh.GUEHeaders = encapHeaders
				// TODO(wenbli): Implement WCMP: there could be a merger of two nexthops, in which case we add their weights.
				egressNhs[rnh] = true
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
			return allEgressNhs[rp], resolvedRoutes[rp], nil
		}
	}

	log.V(1).Infof("Route with prefix %v cannot be resolved in RIB.", ip)
	return nil, nil, nil
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
			if subintf.GetIpv6() != nil {
				for _, a := range subintf.GetIpv6().Address {
					_, cidr, err := net.ParseCIDR(fmt.Sprintf("%s/%d", a.GetIp(), a.GetPrefixLength()))
					if err != nil {
						return nil, fmt.Errorf("invalid IPV6 prefix on interface %s, subinterface %d, %s/%d", intf.GetName(), subintf.GetIndex(), a.GetIp(), a.GetPrefixLength())
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

		// Make deterministic
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
