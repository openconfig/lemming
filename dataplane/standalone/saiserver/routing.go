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
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	log "github.com/golang/glog"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	"github.com/openconfig/lemming/proto/forwarding"
)

type neighbor struct {
	saipb.UnimplementedNeighborServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newNeighbor(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *neighbor {
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

// CreateNeighborEntries adds multiple neighbors to the neighbor table.
func (n *neighbor) CreateNeighborEntries(ctx context.Context, re *saipb.CreateNeighborEntriesRequest) (*saipb.CreateNeighborEntriesResponse, error) {
	resp := &saipb.CreateNeighborEntriesResponse{}
	for _, req := range re.GetReqs() {
		res, err := attrmgr.InvokeAndSave(ctx, n.mgr, n.CreateNeighborEntry, req)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, res)
	}
	return resp, nil
}

type nextHopGroup struct {
	saipb.UnimplementedNextHopGroupServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newNextHopGroup(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *nextHopGroup {
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
	dataplane switchDataplaneAPI
}

func newNextHop(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *nextHop {
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

func (nh *nextHop) CreateNextHops(ctx context.Context, r *saipb.CreateNextHopsRequest) (*saipb.CreateNextHopsResponse, error) {
	resp := &saipb.CreateNextHopsResponse{}
	for _, req := range r.GetReqs() {
		res, err := attrmgr.InvokeAndSave(ctx, nh.mgr, nh.CreateNextHop, req)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, res)
	}
	return resp, nil
}

type route struct {
	saipb.UnimplementedRouteServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newRoute(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *route {
	r := &route{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterRouteServer(s, r)
	return r
}

const ip2meTableID = "ip2metable"

// CreateRouteEntry creates a new route entry.
func (r *route) CreateRouteEntry(ctx context.Context, req *saipb.CreateRouteEntryRequest) (*saipb.CreateRouteEntryResponse, error) {
	rReq := &dpb.AddIPRouteRequest{
		Route: &dpb.Route{
			Prefix: &dpb.RoutePrefix{
				Prefix: &dpb.RoutePrefix_Mask{
					Mask: &dpb.IpMask{
						Addr: req.GetEntry().GetDestination().GetAddr(),
						Mask: req.GetEntry().GetDestination().GetMask(),
					},
				},
				VrfId: req.GetEntry().GetVrId(),
			},
		},
	}

	// TODO(dgrau): Implement CPU actions.
	switch req.GetPacketAction() {
	case saipb.PacketAction_PACKET_ACTION_DROP,
		saipb.PacketAction_PACKET_ACTION_TRAP, // COPY and DROP
		saipb.PacketAction_PACKET_ACTION_DENY: // COPY_CANCEL and DROP
		rReq.Route.Action = dpb.PacketAction_PACKET_ACTION_DROP
	case saipb.PacketAction_PACKET_ACTION_FORWARD,
		saipb.PacketAction_PACKET_ACTION_LOG,     // COPY and FORWARD
		saipb.PacketAction_PACKET_ACTION_TRANSIT: // COPY_CANCEL and FORWARD
		rReq.Route.Action = dpb.PacketAction_PACKET_ACTION_FORWARD
	case saipb.PacketAction_PACKET_ACTION_UNSPECIFIED: // Default action is forward.
		rReq.Route.Action = dpb.PacketAction_PACKET_ACTION_FORWARD
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown action type: %v", req.GetPacketAction())
	}
	nextType := r.mgr.GetType(fmt.Sprint(req.GetNextHopId()))

	swReq := &saipb.GetSwitchAttributeRequest{
		Oid:      req.GetEntry().GetSwitchId(),
		AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT},
	}
	swAttr := &saipb.GetSwitchAttributeResponse{}
	if err := r.mgr.PopulateAttributes(swReq, swAttr); err != nil {
		return nil, err
	}

	// Handle "IP2ME" routes specially.
	if nextType == saipb.ObjectType_OBJECT_TYPE_PORT && req.GetNextHopId() == swAttr.GetAttr().GetCpuPort() {
		ipReq := fwdconfig.TableEntryAddRequest(r.dataplane.ID(), ip2meTableID).AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ExtactEntry(fwdconfig.PacketFieldBytes(forwarding.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).
				WithBytes(req.GetEntry().Destination.GetAddr()))),
			fwdconfig.Action(f))
	}

	// If the packet action is drop, then next hop is optional.
	if rReq.Route.Action == dpb.PacketAction_PACKET_ACTION_FORWARD {
		switch nextType {
		case saipb.ObjectType_OBJECT_TYPE_NEXT_HOP:
			rReq.Route.Hop = &dpb.Route_NextHopId{NextHopId: req.GetNextHopId()}
		case saipb.ObjectType_OBJECT_TYPE_NEXT_HOP_GROUP:
			rReq.Route.Hop = &dpb.Route_NextHopGroupId{NextHopGroupId: req.GetNextHopId()}
		case saipb.ObjectType_OBJECT_TYPE_ROUTER_INTERFACE:
			rReq.Route.Hop = &dpb.Route_InterfaceId{InterfaceId: fmt.Sprint(req.GetNextHopId())}
		case saipb.ObjectType_OBJECT_TYPE_PORT:
			rReq.Route.Hop = &dpb.Route_PortId{PortId: fmt.Sprint(req.GetNextHopId())}
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unknown next hop type: %v", nextType)
		}
	}

	_, err := r.dataplane.AddIPRoute(ctx, rReq)
	if err != nil {
		return nil, err
	}
	return &saipb.CreateRouteEntryResponse{}, nil
}

func (r *route) CreateRouteEntries(ctx context.Context, re *saipb.CreateRouteEntriesRequest) (*saipb.CreateRouteEntriesResponse, error) {
	resp := &saipb.CreateRouteEntriesResponse{}
	for _, req := range re.GetReqs() {
		res, err := attrmgr.InvokeAndSave(ctx, r.mgr, r.CreateRouteEntry, req)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, res)
	}
	return resp, nil
}

type routerInterface struct {
	saipb.UnimplementedRouterInterfaceServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newRouterInterface(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *routerInterface {
	r := &routerInterface{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterRouterInterfaceServer(s, r)
	return r
}

// CreateRouterInterfaces creates a new router interface.
func (ri *routerInterface) CreateRouterInterface(ctx context.Context, req *saipb.CreateRouterInterfaceRequest) (*saipb.CreateRouterInterfaceResponse, error) {
	id := ri.mgr.NextID()
	iReq := &dpb.AddInterfaceRequest{
		Id:      fmt.Sprint(id),
		VrfId:   uint32(req.GetVirtualRouterId()),
		Mtu:     uint64(req.GetMtu()),
		PortIds: []string{fmt.Sprint(req.GetPortId())},
		Mac:     req.GetSrcMacAddress(),
	}
	switch req.GetType() {
	case saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT:
		iReq.Type = dpb.InterfaceType_INTERFACE_TYPE_PORT

	case saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_LOOPBACK: // TODO: Support loopback interfaces
		log.Warning("loopback interfaces not supported")
		return &saipb.CreateRouterInterfaceResponse{Oid: id}, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown interface type: %v", req.GetType())
	}
	if _, err := ri.dataplane.AddInterface(ctx, iReq); err != nil {
		return nil, err
	}

	return &saipb.CreateRouterInterfaceResponse{Oid: id}, nil
}

type vlan struct {
	saipb.UnimplementedVlanServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newVlan(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *vlan {
	v := &vlan{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterVlanServer(s, v)
	return v
}

func (vlan *vlan) CreateVlan(context.Context, *saipb.CreateVlanRequest) (*saipb.CreateVlanResponse, error) {
	id := vlan.mgr.NextID()

	req := &saipb.GetSwitchAttributeRequest{Oid: 1, AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_DEFAULT_STP_INST_ID}}
	resp := &saipb.GetSwitchAttributeResponse{}

	if err := vlan.mgr.PopulateAttributes(req, resp); err != nil {
		return nil, err
	}

	attrs := &saipb.VlanAttribute{
		MemberList:                         []uint64{},
		StpInstance:                        resp.Attr.DefaultStpInstId,
		UnknownNonIpMcastOutputGroupId:     proto.Uint64(0),
		UnknownIpv4McastOutputGroupId:      proto.Uint64(0),
		UnknownIpv6McastOutputGroupId:      proto.Uint64(0),
		UnknownLinklocalMcastOutputGroupId: proto.Uint64(0),
		IngressAcl:                         proto.Uint64(0),
		EgressAcl:                          proto.Uint64(0),
		UnknownUnicastFloodGroup:           proto.Uint64(0),
		UnknownMulticastFloodGroup:         proto.Uint64(0),
		BroadcastFloodGroup:                proto.Uint64(0),
		TamObject:                          []uint64{},
	}
	vlan.mgr.StoreAttributes(id, attrs)
	return &saipb.CreateVlanResponse{
		Oid: id,
	}, nil
}

type bridge struct {
	saipb.UnimplementedBridgeServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newBridge(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *bridge {
	b := &bridge{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterBridgeServer(s, b)
	return b
}

func (br *bridge) CreateBridge(context.Context, *saipb.CreateBridgeRequest) (*saipb.CreateBridgeResponse, error) {
	id := br.mgr.NextID()
	attrs := &saipb.BridgeAttribute{
		PortList:                   []uint64{},
		UnknownUnicastFloodGroup:   proto.Uint64(0),
		UnknownMulticastFloodGroup: proto.Uint64(0),
		BroadcastFloodGroup:        proto.Uint64(0),
	}
	br.mgr.StoreAttributes(id, attrs)
	return &saipb.CreateBridgeResponse{
		Oid: id,
	}, nil
}
