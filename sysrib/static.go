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
	"context"

	log "github.com/golang/glog"
	"github.com/openconfig/gribigo/afthelper"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
)

// convertStaticRoute converts an OC static route to a sysrib Route
func convertStaticRoute(sroute *oc.NetworkInstance_Protocol_Static) *Route {
	var nexthops []*afthelper.NextHopSummary
	for _, snh := range sroute.NextHop {
		// TODO(wenbli): Implement recurse option.
		snh.SetRecurse(true)
		switch nh := snh.NextHop.(type) {
		case nil:
		case oc.UnionString:
			nexthops = append(nexthops, &afthelper.NextHopSummary{
				Weight:          1,
				Address:         string(nh),
				NetworkInstance: fakedevice.DefaultNetworkInstance,
			})
		default:
			log.Warningf("sysrib: Unhandled static route nexthop type (%T): %v", nh, nh)
		}
	}
	if len(nexthops) == 0 {
		return nil
	}
	return &Route{
		Prefix:   *sroute.Prefix,
		NextHops: nexthops,
		RoutePref: RoutePreference{
			AdminDistance: 1,
		},
	}
}

// monitorStaticRoutes starts a gothread to check for static route
// configuration changes.
// It returns an error if there is an error before monitoring can begin.
func (s *Server) monitorStaticRoutes(yclient *ygnmi.Client) error {
	b := &ocpath.Batch{}
	staticroot := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	staticpath := staticroot.StaticAny()
	b.AddPaths(
		staticpath.NextHopAny().NextHop().Config().PathStruct(),
		// TODO(wenbli): Handle these paths.
		// staticpath.NextHopAny().Preference().Config().PathStruct(),
		// staticpath.NextHopAny().Metric().Config().PathStruct(),
		staticpath.NextHopAny().Recurse().Config().PathStruct(),
		staticpath.Prefix().Config().PathStruct(),
	)

	staticRouteWatcher := ygnmi.Watch(
		context.Background(),
		yclient,
		b.Config(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}
			staticp := rootVal.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).GetOrCreateProtocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
			for _, sroute := range staticp.Static {
				if sroute == nil || sroute.Prefix == nil {
					continue
				}
				if route := convertStaticRoute(sroute); route != nil {
					if err := s.setRoute(fakedevice.DefaultNetworkInstance, route, false); err != nil {
						log.Warningf("Failed to add static route: %v", err)
					} else {
						gnmiclient.Replace(context.Background(), yclient, staticroot.Static(sroute.GetPrefix()).State(), sroute)
					}
				}
			}
			return ygnmi.Continue
		},
	)

	// TODO(wenbli): Support static route removal.
	go func() {
		if _, err := staticRouteWatcher.Await(); err != nil {
			log.Warningf("Static route watcher has stopped: %v", err)
		}
	}()
	return nil
}
