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
	"github.com/openconfig/lemming/dataplane"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/grpc/status"

	dpb "github.com/openconfig/lemming/proto/dataplane"
	pb "github.com/openconfig/lemming/proto/sysrib"
)

const (
	SockAddr  = "/tmp/sysrib.api"
	defaultNI = "DEFAULT"
)

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

	// resolvedRoutes contain a map of resolved routes with which to do
	// diff for sending to the dataplane for programming.
	resolvedRoutes map[RouteKey]*ResolvedRoute

	dataplane dataplaneAPI
}

type dataplaneAPI interface {
	ProgramRoute(*ResolvedRoute) error
}

// Dataplane is a wrapper around dpb.HALClient to enable testing before
// resolved route translation.
//
// TODO(wenbli): This is a temporary workaround due to the instability of the
// API. Once the dataplane API is stable, then we'll want to test at the API
// layer instead.
type Dataplane struct {
	dpb.HALClient
}

// programRoute programs the route in the dataplane, returning an error on failure.
func (d *Dataplane) ProgramRoute(r *ResolvedRoute) error {
	log.V(1).Infof("sysrib: programming resolved route: %+v", r)
	rr, err := resolvedRouteToRouteRequest(r)
	if err != nil {
		return err
	}
	_, err = d.InsertRoute(context.Background(), rr)
	return err
}

// NewServer instantiates server to handle client queries.
//
// If dp is nil, then a connection attempt is made.
func NewServer(dp dataplaneAPI) (*Server, error) {
	rib, err := NewSysRIB(nil)
	if err != nil {
		return nil, err
	}

	if dp == nil {
		opts := []grpc.DialOption{grpc.WithTransportCredentials(local.NewCredentials())}
		dpconn, err := grpc.Dial(fmt.Sprintf("localhost:%d", dataplane.Port), opts...)
		if err != nil {
			return nil, fmt.Errorf("cannot dial to HAL service, %v", err)
		}
		dp = &Dataplane{dpb.NewHALClient(dpconn)}
	}

	return &Server{
		rib:            rib,
		interfaces:     map[Interface]bool{},
		resolvedRoutes: map[RouteKey]*ResolvedRoute{},
		dataplane:      dp,
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

func vrfIDToNiName(vrfID uint32) string {
	switch vrfID {
	case 0:
		return defaultNI
	default:
		return strconv.Itoa(int(vrfID))
	}
}

func niNameToVrfID(niName string) (uint32, error) {
	switch niName {
	case defaultNI:
		return 0, nil
	default:
		// TODO(wenbli): This mapping should probably be stored in a map.
		return 1, fmt.Errorf("sysrib: only %s VRF is recognized", defaultNI)
	}
}

func prefixString(prefix *pb.Prefix) (string, error) {
	switch fam := prefix.GetFamily(); fam {
	case pb.Prefix_FAMILY_IPV4:
		// TODO(wenbli): Handle invalid input values.
		return fmt.Sprintf("%s/%d", prefix.GetAddress(), prefix.GetMaskLength()), nil
	default:
		return "", fmt.Errorf("unrecognized prefix family: %v", fam)
	}
}

func resolvedRouteToRouteRequest(r *ResolvedRoute) (*dpb.InsertRouteRequest, error) {
	vrfID, err := niNameToVrfID(r.NIName)
	if err != nil {
		return nil, err
	}

	var nexthops []*dpb.NextHop
	for nh := range r.Nexthops {
		nexthops = append(nexthops, &dpb.NextHop{
			Port:   nh.Port.Name,
			Ip:     nh.Address,
			Weight: nh.Weight,
		})
	}

	return &dpb.InsertRouteRequest{
		Vrf:      uint64(vrfID),
		Prefix:   r.Prefix,
		NextHops: nexthops,
	}, nil
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
		for it := ni.IPV4.Iterate(); it.Next(); {
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
		return nil, status.Errorf(codes.Unimplemented, "route delete not yet supported")
	}

	pfx, err := prefixString(req.Prefix)
	if err != nil {
		return nil, err
	}

	nexthops := []*afthelper.NextHopSummary{}
	for _, nh := range req.GetNexthops() {
		if nh.GetType() != pb.Nexthop_TYPE_IPV4 {
			return nil, status.Errorf(codes.Unimplemented, "non-IPV4 nexthop not supported")
		}
		nexthops = append(nexthops, &afthelper.NextHopSummary{
			Weight:          nh.GetWeight(),
			Address:         nh.GetAddress(),
			NetworkInstance: vrfIDToNiName(nh.GetVrfId()),
		})
	}

	// TODO(wenbli): Check if recursive resolution is an infinite recursion. This happens if there is a cycle.

	niName := vrfIDToNiName(req.GetVrfId())
	if err := s.rib.AddRoute(niName, &Route{
		Prefix:   pfx,
		NextHops: nexthops,
		RoutePref: RoutePreference{
			AdminDistance: uint8(req.GetAdminDistance()),
			Metric:        req.GetMetric(),
		},
	}); err != nil {
		return nil, status.Errorf(codes.Aborted, "error while adding route to sysrib: %v", err)
	}

	if err := s.ResolveAndProgramDiff(); err != nil {
		return nil, status.Errorf(codes.Aborted, "error while resolving sysrib: %v", err)
	}

	// There could be operations carried out by ResolveAndProgramDiff() other than the input route, so we look up our particular prefix.
	status := pb.SetRouteResponse_STATUS_FAIL
	if _, ok := s.resolvedRoutes[RouteKey{Prefix: pfx, NIName: niName}]; ok {
		status = pb.SetRouteResponse_STATUS_SUCCESS
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
	s.setInterface(name, ifindex, enabled)

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
			// FIXME(wenbli): synchronization is required if concurrent calls are possible.
			routes[i] = connectedRoute
			return s.ResolveAndProgramDiff()
		}
	}

	s.rib.AddRoute(niName, connectedRoute)
	return s.ResolveAndProgramDiff()
}

// TODO(wenbli): Do we need to handle interface deletion?
// This is not required in the MVP since basic tests will just need to enable/disable interfaces.

// setInterface responds to INTERFACE_UP/INTERFACE_DOWN messages from the dataplane.
func (s *Server) setInterface(name string, ifindex int32, enabled bool) error {
	s.interfaces[Interface{
		Name:  name,
		Index: ifindex,
	}] = enabled

	return s.ResolveAndProgramDiff()
}
