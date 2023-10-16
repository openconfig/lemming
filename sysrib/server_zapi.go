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

package sysrib

import (
	"fmt"
	"net"

	log "github.com/golang/glog"
	"github.com/openconfig/gribigo/afthelper"
	"github.com/wenovus/gobgp/v3/pkg/zebra"
)

func distributeRoute(s *ZServer, rr *ResolvedRoute, route *Route) {
	// TODO(wenbli): RedistributeRouteDel
	zrouteBody, err := convertToZAPIRoute(rr.RouteKey, route, rr)
	if err != nil {
		log.Warningf("failed to convert resolved route to zebra BGP route: %v", err)
	}
	if zrouteBody != nil && s != nil {
		log.V(1).Info("Sending new route to ZAPI clients: ", zrouteBody)
		s.ClientMutex.RLock()
		for conn := range s.ClientMap {
			serverSendMessage(conn, zebra.RedistributeRouteAdd, zrouteBody)
		}
		s.ClientMutex.RUnlock()
	}
}

// convertToZAPIRoute converts a route to a ZAPI route for redistributing to
// other protocols (e.g. BGP).
func convertToZAPIRoute(routeKey RouteKey, route *Route, rr *ResolvedRoute) (*zebra.IPRouteBody, error) {
	if route.Connected != nil {
		// TODO(wenbli): Connected route redistribution not supported.
		// This is not needed right now since only need to redistribute
		// non-connected routes. It also breaks some of the integration tests.
		return nil, nil
	}
	vrfID, err := niNameToVrfID(routeKey.NIName)
	if err != nil {
		return nil, err
	}

	_, ipnet, err := net.ParseCIDR(rr.Prefix)
	if err != nil {
		return nil, fmt.Errorf("gribigo/zapi: %v", err)
	}
	prefixLen, _ := ipnet.Mask.Size()

	var nexthops []zebra.Nexthop
	for nh := range rr.Nexthops {
		nexthops = append(nexthops, zebra.Nexthop{
			VrfID:  vrfID,
			Gate:   net.ParseIP(nh.Address),
			Weight: uint32(nh.Weight),
		})
	}

	return &zebra.IPRouteBody{
		Flags:   zebra.FlagAllowRecursion,
		Type:    zebra.RouteStatic,
		Safi:    zebra.SafiUnicast,
		Message: zebra.MessageNexthop,
		Prefix: zebra.Prefix{
			Prefix:    ipnet.IP,
			PrefixLen: uint8(prefixLen),
		},
		Nexthops: nexthops,
		Distance: route.RoutePref.AdminDistance,
	}, nil
}

// setZebraRoute calls setRoute after reformatting a zebra-formatted input route.
func (s *Server) setZebraRoute(niName string, zroute *zebra.IPRouteBody) error {
	if s == nil {
		return fmt.Errorf("cannot add route to nil sysrib server")
	}
	log.V(1).Infof("setZebraRoute: %+v", *zroute)
	route := convertZebraRoute(niName, zroute)
	return s.setRoute(niName, route, false)
}

// convertZebraRoute converts a zebra route to a Sysrib route.
func convertZebraRoute(niName string, zroute *zebra.IPRouteBody) *Route {
	var nexthops []*afthelper.NextHopSummary
	for _, znh := range zroute.Nexthops {
		nexthops = append(nexthops, &afthelper.NextHopSummary{
			Weight:          1,
			Address:         znh.Gate.String(),
			NetworkInstance: niName,
		})
	}
	var routePref RoutePreference
	switch zroute.Type {
	case zebra.RouteBGP:
		routePref.AdminDistance = AdminDistanceBGP
	}
	routePref.Metric = zroute.Metric
	return &Route{
		Prefix: fmt.Sprintf("%s/%d", zroute.Prefix.Prefix.String(), zroute.Prefix.PrefixLen),
		// NextHops is the set of IP nexthops that the route uses if
		// it is not a connected route.
		NextHops:  nexthops,
		RoutePref: routePref,
	}
}
