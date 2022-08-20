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

package sysrib

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strconv"

	log "github.com/golang/glog"
	"github.com/openconfig/gribigo/afthelper"

	pb "github.com/openconfig/lemming/proto/sysrib"
)

const SockAddr = "/tmp/sysrib.api"

// Server is the implementation of the Sysrib API.
//
// API:
// - SetRoute
// - addInterface
// - setInterface
type Server struct {
	pb.UnimplementedSysribServer // For forward-compatibility

	// rib contains the current set of all raw routes from the routing
	// clients as well as the configured connected prefixes from the
	// dataplane.
	rib *SysRIB

	// interfaces contains the status of all existing interfaces as
	// indicated by the forwarding plane.
	interfaces map[Interface]bool

	// resolvedRibs contains resolvable routes for each network instance.
	//resolvedRibs map[string]*string_tree.TreeV4

	// resolvedRoutes contain a map of resolved routes with which to do
	// diff for sending to the dataplane for programming.
	resolvedRoutes map[RouteKey]*ResolvedRoute

	dataplane Dataplane
}

type Dataplane interface {
	ProgramRoute(*ResolvedRoute) error
}

// NewServer instantiates server to handle client queries.
func NewServer(dataplane Dataplane) (*Server, error) {
	rib, err := NewSysRIB(nil)
	if err != nil {
		return nil, err
	}

	return &Server{
		rib:            rib,
		interfaces:     map[Interface]bool{},
		resolvedRoutes: map[RouteKey]*ResolvedRoute{},
		dataplane:      dataplane,
	}, nil
}

type RouteKey struct {
	Prefix string
	NIName string
}

// ResolvedRoute represents a route that is ready to be programmed into the forwarding plane.
type ResolvedRoute struct {
	RouteKey

	// NOTE: The order of the nexthops should not matter when being programmed into the forwarding plane. As such, the forwarding plane should sort these nexthops before assigning the hash output for ECMP.
	Nexthops map[ResolvedNexthop]bool
	// TODO(wenbli): backup nexthops.
}

type ResolvedNexthop struct {
	afthelper.NextHopSummary

	Port Interface
}

func vrfIdToNiName(vrfId uint32) string {
	switch vrfId {
	case 0:
		return "DEFAULT"
	default:
		return strconv.Itoa(int(vrfId))
	}
}

func prefixString(prefix *pb.Prefix) (string, error) {
	switch fam := prefix.GetFamily(); fam {
	case pb.Prefix_IPv4:
		// TODO(wenbli): Handle invalid input values.
		return fmt.Sprintf("%s/%d", prefix.GetAddress(), prefix.GetMaskLength()), nil
	default:
		return "", fmt.Errorf("unrecognized prefix family: %v", fam)
	}
}

// programRoute programs the route in the dataplane, returning an error on failure.
func (s *Server) programRoute(r *ResolvedRoute) error {
	// TODO(wenbli): Interface with Daniel's dataplane.
	return s.dataplane.ProgramRoute(r)
}

// ResolveAndProgramDiff walks through the resolved RIBs, updates the forwarding plane.
// TODO(wenbli): handle route deletion.
func (s *Server) ResolveAndProgramDiff() error {
	log.Info("Recalculating resolved RIB")
	for niName, ni := range s.rib.NI {
		for it := ni.IPv4.Iterate(); it.Next(); {
			_, prefix, err := net.ParseCIDR(it.Address().String())
			if err != nil {
				return fmt.Errorf("sysrib: %v", err)
			}
			nhs, err := s.rib.EgressNexthops(niName, prefix, s.interfaces)
			if err != nil {
				return err
			}

			rr := &ResolvedRoute{
				RouteKey: RouteKey{
					// TODO(wenbli): Could it.Address() be different from prefix.String()?
					Prefix: prefix.String(),
					NIName: niName,
				},
				Nexthops: nhs,
			}

			currentRoute, ok := s.resolvedRoutes[rr.RouteKey]
			switch {
			case !ok && len(nhs) > 0, ok && !reflect.DeepEqual(currentRoute, rr):
				if err := s.programRoute(rr); err != nil {
					return fmt.Errorf("failed to program route %+v", rr)
				}
				s.resolvedRoutes[rr.RouteKey] = rr
			default:
				// No diff, so don't do anything.
			}
		}
	}
	return nil
}

// SetRoute implements ROUTE_ADD and ROUTE_DELETE
func (s *Server) SetRoute(_ context.Context, req *pb.SetRouteRequest) (*pb.SetRouteResponse, error) {
	// TODO(wenbli): Handle route deletion.
	if req.Delete {
		return nil, fmt.Errorf("route delete not yet supported")
	}

	pfx, err := prefixString(req.Prefix)
	if err != nil {
		return nil, err
	}

	nexthops := []*afthelper.NextHopSummary{}
	for _, nh := range req.GetNexthops() {
		if nh.GetType() != pb.Nexthop_IPV4 {
			return nil, fmt.Errorf("non-IPv4 nexthop not supported")
		}
		nexthops = append(nexthops, &afthelper.NextHopSummary{
			Weight:          nh.GetWeight(),
			Address:         nh.GetAddress(),
			NetworkInstance: vrfIdToNiName(nh.GetVrfId()),
		})
	}

	// TODO(wenbli): Check if recursive resolution is an infinite recursion. This happens if there is a cycle.

	niName := vrfIdToNiName(req.GetVrfId())
	if err := s.rib.AddRoute(niName, &Route{
		Prefix:   pfx,
		NextHops: nexthops,
		RoutePref: RoutePreference{
			AdminDistance: uint8(req.GetAdminDistance()),
			Metric:        req.GetMetric(),
		},
	}); err != nil {
		return nil, err
	}

	if err := s.ResolveAndProgramDiff(); err != nil {
		return nil, err
	}

	// There could be operations carried out by ResolveAndProgramDiff() other than the input route, so we look up our particular prefix.
	status := pb.SetRouteResponse_FAIL
	if _, ok := s.resolvedRoutes[RouteKey{Prefix: pfx, NIName: niName}]; ok {
		status = pb.SetRouteResponse_SUCCESS
	}
	return &pb.SetRouteResponse{
		Status: status,
	}, nil
}

// AddInterface responds to INTERFACE_ADD messages from the dataplane.
// TODO(wenbli): This is only provided for convenience. It should not be public.
func (s *Server) AddInterface(name string, ifindex int32, enabled bool, prefix string, niName string) error {
	return s.addInterface(name, ifindex, enabled, prefix, niName)
}

// addInterface responds to INTERFACE_ADD messages from the dataplane.
func (s *Server) addInterface(name string, ifindex int32, enabled bool, prefix string, niName string) error {
	intf := Interface{
		Name:  name,
		Index: ifindex,
	}
	s.interfaces[intf] = enabled

	_, pfx, err := net.ParseCIDR(prefix)
	if err != nil {
		return fmt.Errorf("cannot parse connected interface's prefix: %v", err)
	}
	_, routes, err := s.rib.entryForCIDR(niName, pfx)
	if err != nil {
		return err
	}

	connectedRoute := &Route{
		Prefix:    prefix,
		Connected: &intf,
		RoutePref: RoutePreference{
			// Connected routes have admin-distance of 0.
			AdminDistance: 0,
		},
	}

	for i, route := range routes {
		if route.Prefix == prefix && route.Connected != nil {
			// Update the connected route in the tree.
			// FIXME(wenbli): synchronization required?
			routes[i] = connectedRoute
			return s.ResolveAndProgramDiff()
		}
	}

	s.rib.AddRoute(niName, connectedRoute)
	return s.ResolveAndProgramDiff()
}

// TODO(wenbli): Do we need to handle interface deletion?

// setInterface responds to INTERFACE_UP/INTERFACE_DOWN messages from the dataplane.
func (s *Server) setInterface(name string, ifindex int32, enabled bool, prefix string, niName string) error {
	s.interfaces[Interface{
		Name:  name,
		Index: ifindex,
	}] = enabled

	return s.ResolveAndProgramDiff()
}
