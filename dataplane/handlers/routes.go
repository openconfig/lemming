// Copyright 2022 Google LLC
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

package handlers

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/openconfig/lemming/dataplane/internal/engine"
	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/reconciler"
	"github.com/openconfig/ygnmi/schemaless"
	"github.com/openconfig/ygnmi/ygnmi"

	log "github.com/golang/glog"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type route struct {
	w   *ygnmi.Watcher[*dpb.InsertRouteRequest]
	fwd fwdpb.ServiceClient
}

// RouteQuery returns a ygnmi query for a route with the given prefix and vrf.
func RouteQuery(vrf uint64, prefix string) ygnmi.ConfigQuery[*dpb.InsertRouteRequest] {
	q, err := schemaless.NewConfig[*dpb.InsertRouteRequest](fmt.Sprintf("/dataplane/routes/route[prefix=%s][vrf=%d]", prefix, vrf), gnmi.InternalOrigin)
	if err != nil {
		log.Fatal(err)
	}
	return q
}

var (
	routesQuery ygnmi.WildcardQuery[*dpb.InsertRouteRequest]
)

// NewRoute returns a new route reconciler.
func NewRoute(fwd fwdpb.ServiceClient) *reconciler.BuiltReconciler {
	r := &route{
		fwd: fwd,
	}
	return reconciler.NewBuilder("dataplane-routes").WithStart(r.start).Build()
}

func (r *route) start(ctx context.Context, client *ygnmi.Client) error {
	r.w = ygnmi.WatchAll(ctx, client, routesQuery, func(v *ygnmi.Value[*dpb.InsertRouteRequest]) error {
		route, present := v.Val()
		prefix := v.Path.Elem[2].Key["prefix"]
		vrf, err := strconv.ParseUint(v.Path.Elem[2].Key["vrf"], 10, 64)
		if err != nil {
			log.Warningf("non-int vrf set in path: %v", err)
			return ygnmi.Continue
		}
		if vrf != 0 {
			log.Warningf("non-zero vrf")
			return ygnmi.Continue
		}

		_, ipNet, err := net.ParseCIDR(prefix)
		if err != nil {
			log.Warningf("failed to parse prefix: %v", err)
			return ygnmi.Continue
		}
		ip := ipNet.IP.To4()
		isIPv4 := true
		if ip == nil {
			ip = ipNet.IP.To16()
			isIPv4 = false
		}

		if !present {
			if err := engine.DeleteIPRoute(ctx, r.fwd, isIPv4, ipNet.IP, ipNet.Mask); err != nil {
				log.Warningf("failed to delete route: %v", err)
				return ygnmi.Continue
			}
			return ygnmi.Continue
		}
		if err := engine.AddIPRoute(ctx, r.fwd, isIPv4, ip, ipNet.Mask, route.GetNextHops()); err != nil {
			log.Warningf("failed to delete route: %v", err)
		}

		return ygnmi.Continue
	})
	return nil
}

func init() {
	q, err := schemaless.NewWildcard[*dpb.InsertRouteRequest]("/dataplane/routes/route[prefix=*][vrf=*]", gnmi.InternalOrigin)
	if err != nil {
		log.Fatal(err)
	}
	routesQuery = q
}
