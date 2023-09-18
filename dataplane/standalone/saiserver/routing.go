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

package saiserver

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	dpb "github.com/openconfig/lemming/proto/dataplane"
)

type routingDataplaneAPI interface {
	AddNeighbor(ctx context.Context, req *dpb.AddNeighborRequest) (*dpb.AddNeighborResponse, error)
	AddNextHopGroup(ctx context.Context, req *dpb.AddNextHopGroupRequest) (*dpb.AddNextHopGroupResponse, error)
	AddNextHop(ctx context.Context, req *dpb.AddNextHopRequest) (*dpb.AddNextHopResponse, error)
	AddIPRoute(ctx context.Context, req *dpb.AddIPRouteRequest) (*dpb.AddIPRouteResponse, error)
	AddInterface(ctx context.Context, req *dpb.AddInterfaceRequest) (*dpb.AddInterfaceResponse, error)
}

type neighbor struct {
	saipb.UnimplementedNeighborServer
	mgr       *attrmgr.AttrMgr
	dataplane routingDataplaneAPI
}

func newNeighbor(mgr *attrmgr.AttrMgr, dataplane routingDataplaneAPI, s *grpc.Server) *neighbor {
	n := &neighbor{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterNeighborServer(s, n)
	return n
}

// CreateNeighborEntry adds a neighbor to the neighbor table.
func (n *neighbor) CreateNeighborEntry(ctx context.Context, req *saipb.CreateNeighborEntryRequest) (*saipb.CreateNeighborEntryResponse, error) {
	_, err := n.dataplane.AddNeighbor(ctx, &dpb.AddNeighborRequest{
		Dev: &dpb.AddNeighborRequest_InterfaceId{
			InterfaceId: fmt.Sprint(req.GetEntry().GetRifId()),
		},
		Mac: req.GetDstMacAddress(),
		Ip: &dpb.AddNeighborRequest_IpBytes{
			IpBytes: req.GetEntry().GetIpAddress(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &saipb.CreateNeighborEntryResponse{}, nil
}

type nextHopGroup struct {
	saipb.UnimplementedNextHopGroupServer
	mgr       *attrmgr.AttrMgr
	dataplane routingDataplaneAPI
}

func newNextHopGroup(mgr *attrmgr.AttrMgr, dataplane routingDataplaneAPI, s *grpc.Server) *nextHopGroup {
	n := &nextHopGroup{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterNextHopGroupServer(s, n)
	return n
}

// CreateNextHopGroupMember creates a next hop group.
func (nhg *nextHopGroup) CreateNextHopGroup(ctx context.Context, req *saipb.CreateNextHopGroupRequest) (*saipb.CreateNextHopGroupResponse, error) {
	id := nhg.mgr.NextID()

	if req.GetType() != saipb.NextHopGroupType_NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported req type: %v", req.GetType())
	}

	_, err := nhg.dataplane.AddNextHopGroup(ctx, &dpb.AddNextHopGroupRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &saipb.CreateNextHopGroupResponse{
		Oid: id,
	}, nil
}

// CreateNextHopGroupMember adds a next hop to a next hop group.
func (nhg *nextHopGroup) CreateNextHopGroupMember(ctx context.Context, req *saipb.CreateNextHopGroupMemberRequest) (*saipb.CreateNextHopGroupMemberResponse, error) {
	_, err := nhg.dataplane.AddNextHopGroup(ctx, &dpb.AddNextHopGroupRequest{
		Id: req.GetNextHopGroupId(),
		List: &dpb.NextHopIDList{
			Hops:    []uint64{req.GetNextHopId()},
			Weights: []uint64{uint64(req.GetWeight())},
		},
		Mode: dpb.GroupUpdateMode_GROUP_UPDATE_MODE_APPEND,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type nextHop struct {
	saipb.UnimplementedNextHopServer
	mgr       *attrmgr.AttrMgr
	dataplane routingDataplaneAPI
}

func newNextHop(mgr *attrmgr.AttrMgr, dataplane routingDataplaneAPI, s *grpc.Server) *nextHop {
	n := &nextHop{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterNextHopServer(s, n)
	return n
}

// CreateNextHop creates a new next hop.
func (nh *nextHop) CreateNextHop(ctx context.Context, req *saipb.CreateNextHopRequest) (*saipb.CreateNextHopResponse, error) {
	id := nh.mgr.NextID()

	if req.GetType() != saipb.NextHopType_NEXT_HOP_TYPE_IP {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported req type: %v", req.GetType())
	}

	_, err := nh.dataplane.AddNextHop(ctx, &dpb.AddNextHopRequest{
		Id: id,
		NextHop: &dpb.NextHop{
			Dev: &dpb.NextHop_Interface{
				Interface: fmt.Sprint(req.GetRouterInterfaceId()),
			},
			Ip: &dpb.NextHop_IpBytes{
				IpBytes: req.GetIp(),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &saipb.CreateNextHopResponse{
		Oid: id,
	}, nil
}

type route struct {
	saipb.UnimplementedRouteServer
	mgr       *attrmgr.AttrMgr
	dataplane routingDataplaneAPI
}

func newRoute(mgr *attrmgr.AttrMgr, dataplane routingDataplaneAPI, s *grpc.Server) *route {
	r := &route{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterRouteServer(s, r)
	return r
}

type routerInterface struct {
	saipb.UnimplementedRouterInterfaceServer
	mgr       *attrmgr.AttrMgr
	dataplane routingDataplaneAPI
}

func newRouterInterface(mgr *attrmgr.AttrMgr, dataplane routingDataplaneAPI, s *grpc.Server) *routerInterface {
	r := &routerInterface{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterRouterInterfaceServer(s, r)
	return r
}
