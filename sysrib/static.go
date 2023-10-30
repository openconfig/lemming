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
func convertStaticRoute(prefix string, sroute *oc.NetworkInstance_Protocol_Static) *Route {
	var nexthops []*afthelper.NextHopSummary
	if sroute != nil {
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
	}
	return &Route{
		Prefix:   prefix,
		NextHops: nexthops,
		RoutePref: RoutePreference{
			AdminDistance: 1,
		},
	}
}

// monitorStaticRoutes starts a gothread to check for static route
// configuration changes.
// It returns an error if there is an error before monitoring can begin.
func (s *Server) monitorStaticRoutes(ctx context.Context, yclient *ygnmi.Client) error {
	staticroot := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
	staticpath := staticroot.StaticAny()
	staticpathMap := staticroot.StaticMap()

	b := ygnmi.NewBatch[map[string]*oc.NetworkInstance_Protocol_Static](staticpathMap.Config())
	b.AddPaths(
		staticpath.NextHopAny().NextHop().Config(),
		// TODO(wenbli): Handle these paths.
		// staticpath.NextHopAny().Preference().Config(),
		// staticpath.NextHopAny().Metric().Config(),
		staticpath.NextHopAny().Recurse().Config(),
		staticpath.Prefix().Config(),
	)

	prevIntfs := map[string]struct{}{}

	staticRouteWatcher := ygnmi.Watch(
		ctx,
		yclient,
		b.Query(),
		func(static *ygnmi.Value[map[string]*oc.NetworkInstance_Protocol_Static]) error {
			staticMap, ok := static.Val()
			currentIntfs := map[string]struct{}{}
			if ok {
				for prefix, sroute := range staticMap {
					if sroute == nil || sroute.Prefix == nil {
						continue
					}
					if route := convertStaticRoute(prefix, sroute); route != nil {
						currentIntfs[prefix] = struct{}{}
						if err := s.setRoute(ctx, fakedevice.DefaultNetworkInstance, route, false); err != nil {
							log.Warningf("Failed to add static route: %v", err)
						} else {
							gnmiclient.Replace(ctx, yclient, staticroot.Static(sroute.GetPrefix()).State(), sroute)
						}
					}
				}
			}
			for prefix := range prevIntfs {
				if _, ok := currentIntfs[prefix]; !ok {
					if route := convertStaticRoute(prefix, nil); route != nil {
						if err := s.setRoute(ctx, fakedevice.DefaultNetworkInstance, route, true); err != nil {
							log.Warningf("Failed to delete static route: %v", err)
						} else {
							gnmiclient.Delete(ctx, yclient, staticroot.Static(prefix).State())
						}
					}
				}
			}
			prevIntfs = currentIntfs
			return ygnmi.Continue
		},
	)

	go func() {
		if _, err := staticRouteWatcher.Await(); err != nil {
			log.Warningf("Static route watcher has stopped: %v", err)
		}
	}()
	return nil
}
