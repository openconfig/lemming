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
	"encoding/binary"
	"fmt"
	"log/slog"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/gnmi/errlist"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
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
	entry := fwdconfig.TableEntryAddRequest(n.dataplane.ID(), NeighborTable).AppendEntry(fwdconfig.EntryDesc(fwdconfig.ExactEntry(
		fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE).WithUint64(req.GetEntry().GetRifId()),
		fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP).WithBytes(req.GetEntry().GetIpAddress()),
	)), fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).WithValue(req.GetDstMacAddress())),
	).Build()

	// After updated the dst_mac, continue to egress.
	entry.Entries[0].Actions = append(entry.Entries[0].Actions, getEgressPipeline()...)

	if _, err := n.dataplane.TableEntryAdd(ctx, entry); err != nil {
		return nil, err
	}
	return &saipb.CreateNeighborEntryResponse{}, nil
}

// CreateNeighborEntry adds a neighbor to the neighbor table.
func (n *neighbor) RemoveNeighborEntry(ctx context.Context, req *saipb.RemoveNeighborEntryRequest) (*saipb.RemoveNeighborEntryResponse, error) {
	entry := fwdconfig.EntryDesc(fwdconfig.ExactEntry(
		fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE).WithUint64(req.GetEntry().GetRifId()),
		fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP).WithBytes(req.GetEntry().GetIpAddress()),
	)).Build()

	if _, err := n.dataplane.TableEntryRemove(ctx, &fwdpb.TableEntryRemoveRequest{
		ContextId: &fwdpb.ContextId{Id: n.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: NeighborTable}},
		EntryDesc: entry,
	}); err != nil {
		return nil, err
	}
	return &saipb.RemoveNeighborEntryResponse{}, nil
}

// RemoveNeighborEntries removes multiple neighbors to the neighbor table.
func (n *neighbor) RemoveNeighborEntries(ctx context.Context, re *saipb.RemoveNeighborEntriesRequest) (*saipb.RemoveNeighborEntriesResponse, error) {
	resp := &saipb.RemoveNeighborEntriesResponse{}
	for _, req := range re.GetReqs() {
		res, err := attrmgr.InvokeAndSave(ctx, n.mgr, n.RemoveNeighborEntry, req)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, res)
	}
	return resp, nil
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

type groupMember struct {
	nextHop uint64 // ID of the next hop
	weight  uint32
}

type nextHopGroup struct {
	saipb.UnimplementedNextHopGroupServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
	groups    map[uint64]map[uint64]*groupMember // groups is map of next hop groups to a map of next hops
	groupIsV4 map[uint64]bool                    // map from group id to IP protocol version
}

func newNextHopGroup(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *nextHopGroup {
	n := &nextHopGroup{
		mgr:       mgr,
		dataplane: dataplane,
		groups:    map[uint64]map[uint64]*groupMember{},
		groupIsV4: map[uint64]bool{},
	}
	saipb.RegisterNextHopGroupServer(s, n)
	return n
}

// CreateNextHopGroup creates a next hop group.
func (nhg *nextHopGroup) CreateNextHopGroup(ctx context.Context, req *saipb.CreateNextHopGroupRequest) (*saipb.CreateNextHopGroupResponse, error) {
	id := nhg.mgr.NextID()

	switch req.GetType() {
	case saipb.NextHopGroupType_NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP:
		nhg.groups[id] = map[uint64]*groupMember{}
		return &saipb.CreateNextHopGroupResponse{
			Oid: id,
		}, nil
	case saipb.NextHopGroupType_NEXT_HOP_GROUP_TYPE_ECMP_WITH_MEMBERS:
		nhg.groups[id] = map[uint64]*groupMember{}

		memReq := &saipb.CreateNextHopGroupMembersRequest{}
		for i, nh := range req.GetNextHopList() {
			memReq.Reqs = append(memReq.Reqs, &saipb.CreateNextHopGroupMemberRequest{
				Switch:         switchID,
				NextHopGroupId: proto.Uint64(id),
				NextHopId:      proto.Uint64(nh),
				Weight:         proto.Uint32(req.GetNextHopMemberWeightList()[i]),
			})
		}

		_, err := attrmgr.InvokeAndSave(ctx, nhg.mgr, nhg.CreateNextHopGroupMembers, memReq)
		if err != nil {
			return nil, err
		}
		return &saipb.CreateNextHopGroupResponse{
			Oid: id,
		}, nil

	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported req type: %v", req.GetType())
	}
}

// CreateNextHopGroups creates multiple next hop groups.
func (nhg *nextHopGroup) CreateNextHopGroups(ctx context.Context, r *saipb.CreateNextHopGroupsRequest) (*saipb.CreateNextHopGroupsResponse, error) {
	resp := &saipb.CreateNextHopGroupsResponse{}
	for _, req := range r.GetReqs() {
		res, err := attrmgr.InvokeAndSave(ctx, nhg.mgr, nhg.CreateNextHopGroup, req)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, res)
	}
	return resp, nil
}

// updateNextHopGroupMember updates the next hop group.
// If m is nil, remove mid from the group(key: nhgid), otherwise add m to group with mid as the key.
func (nhg *nextHopGroup) updateNextHopGroupMember(ctx context.Context, nhgid, mid uint64, m *groupMember) error {
	group := nhg.groups[nhgid]
	if group == nil {
		return status.Errorf(codes.FailedPrecondition, "group %d does not exist", nhgid)
	}
	if m != nil {
		if _, ok := nhg.groupIsV4[nhgid]; !ok { // Use the first member added to group to determine if the group is ipv4.
			nhAttr := &saipb.GetNextHopAttributeResponse{}
			err := nhg.mgr.PopulateAttributes(&saipb.GetNextHopAttributeRequest{
				Oid:      m.nextHop,
				AttrType: []saipb.NextHopAttr{saipb.NextHopAttr_NEXT_HOP_ATTR_IP},
			}, nhAttr)
			if err != nil {
				return fmt.Errorf("failed to retrieve next hop attr: %v", err)
			}
			nhg.groupIsV4[nhgid] = len(nhAttr.GetAttr().GetIp()) == 4
		}
		group[mid] = m
	} else {
		delete(group, mid)
	}
	var actLists []*fwdpb.ActionList
	for _, member := range group {
		action := fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID).WithUint64Value(member.nextHop))
		actLists = append(actLists, &fwdpb.ActionList{
			Weight:  uint64(member.weight),
			Actions: []*fwdpb.ActionDesc{action.Build()},
		})
	}

	swAttr := &saipb.GetSwitchAttributeResponse{}
	err := nhg.mgr.PopulateAttributes(&saipb.GetSwitchAttributeRequest{
		Oid:      switchID,
		AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_ECMP_HASH_IPV4, saipb.SwitchAttr_SWITCH_ATTR_ECMP_HASH_IPV6},
	}, swAttr)
	if err != nil {
		return fmt.Errorf("failed to retrieve hash id: %v", err)
	}
	hashID := swAttr.GetAttr().GetEcmpHashIpv6()
	if nhg.groupIsV4[nhgid] {
		hashID = swAttr.GetAttr().GetEcmpHashIpv4()
	}
	hashAttr := &saipb.HashAttribute{}
	err = nhg.mgr.PopulateAllAttributes(fmt.Sprint(hashID), hashAttr)
	if err != nil {
		return fmt.Errorf("failed to retrieve hash attrs: %v", err)
	}

	fieldsID, err := convertHashFields(hashAttr.GetNativeHashFieldList())
	if err != nil {
		return fmt.Errorf("failed to compute hash fields: %v", err)
	}

	actions := []*fwdpb.ActionDesc{{
		ActionType: fwdpb.ActionType_ACTION_TYPE_SELECT_ACTION_LIST,
		Action: &fwdpb.ActionDesc_Select{
			Select: &fwdpb.SelectActionListActionDesc{
				SelectAlgorithm: fwdpb.SelectActionListActionDesc_SELECT_ALGORITHM_CRC32, // TODO: should algo + hash be configurable?
				FieldIds:        fieldsID,
				ActionLists:     actLists,
			},
		},
	}, {
		ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
		Action: &fwdpb.ActionDesc_Lookup{
			Lookup: &fwdpb.LookupActionDesc{
				TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{
					Id: NHTable,
				}},
			},
		},
	}}

	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: nhg.dataplane.ID()},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: NHGTable,
			},
		},
		Entries: []*fwdpb.TableEntryAddRequest_Entry{{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes: binary.BigEndian.AppendUint64(nil, nhgid),
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_GROUP_ID,
								},
							},
						}},
					},
				},
			},
			Actions: actions,
		}},
	}
	_, err = nhg.dataplane.TableEntryAdd(ctx, entries)
	return err
}

// RemoveNextHopGroup removes the next hop group specified in the OID.
func (nhg *nextHopGroup) RemoveNextHopGroup(_ context.Context, req *saipb.RemoveNextHopGroupRequest) (*saipb.RemoveNextHopGroupResponse, error) {
	oid := req.GetOid()
	if _, ok := nhg.groups[oid]; !ok {
		return nil, status.Errorf(codes.FailedPrecondition, "group %d does not exist", oid)
	}
	delete(nhg.groups, oid)

	entry := fwdconfig.EntryDesc(fwdconfig.ExactEntry(
		fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_GROUP_ID).WithUint64(oid))).Build()
	nhgReq := &fwdpb.TableEntryRemoveRequest{
		ContextId: &fwdpb.ContextId{Id: nhg.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: NHGTable}},
		EntryDesc: entry,
	}

	if _, err := nhg.dataplane.TableEntryRemove(context.Background(), nhgReq); err != nil {
		return nil, err
	}
	return &saipb.RemoveNextHopGroupResponse{}, nil
}

// CreateNextHopGroupMember adds a next hop to a next hop group.
func (nhg *nextHopGroup) CreateNextHopGroupMember(ctx context.Context, req *saipb.CreateNextHopGroupMemberRequest) (*saipb.CreateNextHopGroupMemberResponse, error) {
	nhgid := req.GetNextHopGroupId()
	mid := nhg.mgr.NextID()
	m := &groupMember{
		nextHop: req.GetNextHopId(),
		weight:  req.GetWeight(),
	}
	if err := nhg.updateNextHopGroupMember(ctx, nhgid, mid, m); err != nil {
		return nil, err
	}
	return &saipb.CreateNextHopGroupMemberResponse{Oid: mid}, nil
}

// CreateNextHopGroups creates multiple next hop groups.
func (nhg *nextHopGroup) CreateNextHopGroupMembers(ctx context.Context, r *saipb.CreateNextHopGroupMembersRequest) (*saipb.CreateNextHopGroupMembersResponse, error) {
	resp := &saipb.CreateNextHopGroupMembersResponse{}
	for _, req := range r.GetReqs() {
		res, err := attrmgr.InvokeAndSave(ctx, nhg.mgr, nhg.CreateNextHopGroupMember, req)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, res)
	}
	return resp, nil
}

// RemoveNextHopGroupMember remove the next hop group member specified in the OID.
// Only need to remove with the desc.
func (nhg *nextHopGroup) RemoveNextHopGroupMember(ctx context.Context, req *saipb.RemoveNextHopGroupMemberRequest) (*saipb.RemoveNextHopGroupMemberResponse, error) {
	locateMember := func(oid uint64) (uint64, uint64, error) {
		for nhgid, nhg := range nhg.groups {
			for mid := range nhg {
				if mid == oid {
					return nhgid, mid, nil
				}
			}
		}
		return 0, 0, fmt.Errorf("cannot find member with id=%d", oid)
	}
	nhgid, mid, err := locateMember(req.GetOid())
	if err != nil {
		return nil, err
	}

	if err := nhg.updateNextHopGroupMember(ctx, nhgid, mid, nil); err != nil {
		return nil, err
	}
	return &saipb.RemoveNextHopGroupMemberResponse{}, nil
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

	var actions []*fwdpb.ActionDesc

	switch req.GetType() {
	case saipb.NextHopType_NEXT_HOP_TYPE_IP:
		actions = []*fwdpb.ActionDesc{
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE).WithUint64Value(req.GetRouterInterfaceId())).Build(),
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP).WithValue(req.GetIp())).Build(),
			fwdconfig.Action(fwdconfig.LookupAction(NHActionTable)).Build(),
		}
	case saipb.NextHopType_NEXT_HOP_TYPE_TUNNEL_ENCAP:
		tunnel := &saipb.GetTunnelAttributeResponse{}
		err := nh.mgr.PopulateAttributes(&saipb.GetTunnelAttributeRequest{Oid: req.GetTunnelId(), AttrType: []saipb.TunnelAttr{saipb.TunnelAttr_TUNNEL_ATTR_TYPE}}, tunnel)
		if err != nil {
			return nil, err
		}

		headerID := fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6
		if len(req.Ip) == 4 {
			headerID = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4
		}
		actions = []*fwdpb.ActionDesc{}

		switch tunnel.GetAttr().GetType() {
		case saipb.TunnelType_TUNNEL_TYPE_IPINIP:
			actions = append(actions, &fwdpb.ActionDesc{
				ActionType: fwdpb.ActionType_ACTION_TYPE_ENCAP,
				Action: &fwdpb.ActionDesc_Encap{
					Encap: &fwdpb.EncapActionDesc{
						HeaderId: headerID,
					},
				},
			})
		case saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE:
			actions = append(actions, &fwdpb.ActionDesc{
				ActionType: fwdpb.ActionType_ACTION_TYPE_ENCAP,
				Action: &fwdpb.ActionDesc_Encap{
					Encap: &fwdpb.EncapActionDesc{
						HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
					},
				},
			}, &fwdpb.ActionDesc{
				ActionType: fwdpb.ActionType_ACTION_TYPE_ENCAP,
				Action: &fwdpb.ActionDesc_Encap{
					Encap: &fwdpb.EncapActionDesc{
						HeaderId: headerID,
					},
				},
			})
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unsupported tunnel type: %v", tunnel.GetAttr().GetType())
		}

		actions = append(actions,
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithValue(req.GetIp())).Build(),
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP).WithValue(req.GetIp())).Build(),
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64Value(req.GetTunnelId())).Build(),
			fwdconfig.Action(fwdconfig.LookupAction(NHActionTable)).Build(),
			fwdconfig.Action(fwdconfig.LookupAction(TunnelEncap)).Build(),
		)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported req type: %v", req.GetType())
	}

	nhReq := fwdconfig.TableEntryAddRequest(nh.dataplane.ID(), NHTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID).WithUint64(id))),
	).Build()
	nhReq.Entries[0].Actions = actions

	if _, err := nh.dataplane.TableEntryAdd(ctx, nhReq); err != nil {
		return nil, err
	}
	return &saipb.CreateNextHopResponse{
		Oid: id,
	}, nil
}

// RemoveNextHop removes the next hop with the OID specified.
func (nh *nextHop) RemoveNextHop(ctx context.Context, r *saipb.RemoveNextHopRequest) (*saipb.RemoveNextHopResponse, error) {
	entry := fwdconfig.EntryDesc(fwdconfig.ExactEntry(
		fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID).WithUint64(r.GetOid()))).Build()
	nhReq := &fwdpb.TableEntryRemoveRequest{
		ContextId: &fwdpb.ContextId{Id: nh.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: NHTable}},
		EntryDesc: entry,
	}

	if _, err := nh.dataplane.TableEntryRemove(ctx, nhReq); err != nil {
		return nil, err
	}
	return &saipb.RemoveNextHopResponse{}, nil
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

func (nh *nextHop) RemoveNextHops(ctx context.Context, r *saipb.RemoveNextHopsRequest) (*saipb.RemoveNextHopsResponse, error) {
	resp := &saipb.RemoveNextHopsResponse{}
	for _, req := range r.GetReqs() {
		res, err := attrmgr.InvokeAndSave(ctx, nh.mgr, nh.RemoveNextHop, req)
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

// CreateRouteEntry creates a new route entry.
func (r *route) CreateRouteEntry(ctx context.Context, req *saipb.CreateRouteEntryRequest) (*saipb.CreateRouteEntryResponse, error) {
	fib := FIBV6Table
	if len(req.GetEntry().GetDestination().GetAddr()) == 4 {
		fib = FIBV4Table
	}

	entry := fwdconfig.TableEntryAddRequest(r.dataplane.ID(), fib).AppendEntry(fwdconfig.EntryDesc(
		fwdconfig.PrefixEntry(
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64(req.GetEntry().GetVrId()),
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(
				req.GetEntry().GetDestination().GetAddr(),
				req.GetEntry().GetDestination().GetMask(),
			),
		),
	))

	forward := true
	if req.GetPacketAction() == saipb.PacketAction_PACKET_ACTION_DROP ||
		req.GetPacketAction() == saipb.PacketAction_PACKET_ACTION_TRAP ||
		req.GetPacketAction() == saipb.PacketAction_PACKET_ACTION_DENY {
		forward = false
	}
	nextType := r.mgr.GetType(fmt.Sprint(req.GetNextHopId()))

	// If the packet action is drop, then next hop is optional.
	if forward {
		switch nextType {
		case saipb.ObjectType_OBJECT_TYPE_NEXT_HOP:
			entry.AppendActions(fwdconfig.Action(
				fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID).WithUint64Value(req.GetNextHopId())),
				fwdconfig.Action(fwdconfig.LookupAction(NHTable)),
			)
		case saipb.ObjectType_OBJECT_TYPE_NEXT_HOP_GROUP:
			entry.AppendActions(fwdconfig.Action(
				fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_GROUP_ID).WithUint64Value(req.GetNextHopId())),
				fwdconfig.Action(fwdconfig.LookupAction(NHGTable)),
			)
		case saipb.ObjectType_OBJECT_TYPE_ROUTER_INTERFACE:
			entry.AppendActions(
				// Set the next hop IP in the packet's metadata.
				fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_COPY, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP).WithFieldSrc(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST)),
				// Set the output iface.
				fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE).WithUint64Value(req.GetNextHopId())),
			)
		case saipb.ObjectType_OBJECT_TYPE_PORT:
			attrReq := &saipb.GetSwitchAttributeRequest{
				Oid:      req.GetEntry().GetSwitchId(),
				AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT},
			}
			resp := &saipb.GetSwitchAttributeResponse{}
			if err := r.mgr.PopulateAttributes(attrReq, resp); err != nil {
				return nil, err
			}
			if req.GetNextHopId() == *resp.Attr.CpuPort {
				_, err := r.dataplane.TableEntryAdd(ctx, fwdconfig.TableEntryAddRequest(r.dataplane.ID(), trapTableID).
					AppendEntry(
						fwdconfig.EntryDesc(fwdconfig.FlowEntry(
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(
								req.GetEntry().GetDestination().GetAddr(),
								req.GetEntry().GetDestination().GetMask()),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64(req.GetEntry().GetVrId()))),
						fwdconfig.Action(fwdconfig.TransmitAction(fmt.Sprint(req.GetNextHopId())).WithImmediate(true))).
					Build())
				if err != nil {
					return nil, status.Errorf(codes.Internal, "failed to add next IP2ME route: %v", nextType)
				}
				return &saipb.CreateRouteEntryResponse{}, nil
			}
			entry.AppendActions(
				// Set the next hop IP in the packet's metadata.
				fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_COPY, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP).WithFieldSrc(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST)),
				// Set the output port.
				fwdconfig.Action(fwdconfig.TransmitAction(fmt.Sprint(req.GetNextHopId()))),
			)
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unknown next hop type: %v", nextType)
		}
	} else {
		entry.AppendActions(fwdconfig.Action(fwdconfig.DropAction()))
	}

	route := entry.Build()
	_, err := r.dataplane.TableEntryAdd(ctx, route)
	if err != nil {
		return nil, err
	}
	return &saipb.CreateRouteEntryResponse{}, nil
}

func (r *route) CreateRouteEntries(ctx context.Context, re *saipb.CreateRouteEntriesRequest) (*saipb.CreateRouteEntriesResponse, error) {
	var errs errlist.List
	resp := &saipb.CreateRouteEntriesResponse{}
	for _, req := range re.GetReqs() {
		res, err := attrmgr.InvokeAndSave(ctx, r.mgr, r.CreateRouteEntry, req)
		errs.Add(err)
		resp.Resps = append(resp.Resps, res)
	}
	return resp, errs.Err()
}

func (r *route) RemoveRouteEntry(ctx context.Context, req *saipb.RemoveRouteEntryRequest) (*saipb.RemoveRouteEntryResponse, error) {
	fib := FIBV6Table
	if len(req.GetEntry().GetDestination().GetAddr()) == 4 {
		fib = FIBV4Table
	}

	_, err := r.dataplane.TableEntryRemove(ctx, &fwdpb.TableEntryRemoveRequest{
		ContextId: &fwdpb.ContextId{Id: r.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fib}},
		EntryDesc: fwdconfig.EntryDesc(
			fwdconfig.PrefixEntry(
				fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64(req.GetEntry().GetVrId()),
				fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(
					req.GetEntry().GetDestination().GetAddr(),
					req.GetEntry().GetDestination().GetMask(),
				),
			),
		).Build(),
	})
	return &saipb.RemoveRouteEntryResponse{}, err
}

func (r *route) RemoveRouteEntries(ctx context.Context, re *saipb.RemoveRouteEntriesRequest) (*saipb.RemoveRouteEntriesResponse, error) {
	resp := &saipb.RemoveRouteEntriesResponse{}
	for _, req := range re.GetReqs() {
		res, err := attrmgr.InvokeAndSave(ctx, r.mgr, r.RemoveRouteEntry, req)
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

func ifaceCounterID(oid uint64, input bool) string {
	if input {
		return fmt.Sprintf("%d-in-counter", oid)
	}
	return fmt.Sprintf("%d-out-counter", oid)
}

// CreateRouterInterfaces creates a new router interface.
func (ri *routerInterface) CreateRouterInterface(ctx context.Context, req *saipb.CreateRouterInterfaceRequest) (*saipb.CreateRouterInterfaceResponse, error) {
	id := ri.mgr.NextID()
	switch req.GetType() {
	case saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT:
	case saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_LOOPBACK: // TODO: Support loopback interfaces
		slog.WarnContext(ctx, "loopback interfaces not supported")
		return &saipb.CreateRouterInterfaceResponse{Oid: id}, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown interface type: %v", req.GetType())
	}
	fwdCtx, err := ri.dataplane.FindContext(&fwdpb.ContextId{Id: ri.dataplane.ID()})
	if err != nil {
		return nil, err
	}
	obj, err := fwdCtx.Objects.FindID(&fwdpb.ObjectId{Id: fmt.Sprint(req.GetPortId())})
	if err != nil {
		return nil, err
	}

	inCounter := ifaceCounterID(id, true)
	outCounter := ifaceCounterID(id, false)

	_, err = ri.dataplane.FlowCounterCreate(ctx, &fwdpb.FlowCounterCreateRequest{
		ContextId: &fwdpb.ContextId{Id: ri.dataplane.ID()},
		Id:        &fwdpb.FlowCounterId{ObjectId: &fwdpb.ObjectId{Id: inCounter}},
	})
	if err != nil {
		return nil, err
	}
	_, err = ri.dataplane.FlowCounterCreate(ctx, &fwdpb.FlowCounterCreateRequest{
		ContextId: &fwdpb.ContextId{Id: ri.dataplane.ID()},
		Id:        &fwdpb.FlowCounterId{ObjectId: &fwdpb.ObjectId{Id: outCounter}},
	})
	if err != nil {
		return nil, err
	}

	// Link the port to the interface.
	_, err = ri.dataplane.TableEntryAdd(ctx, fwdconfig.TableEntryAddRequest(ri.dataplane.ID(), inputIfaceTable).
		AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT).WithUint64(uint64(obj.NID())))),
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_INPUT_IFACE).WithUint64Value(id)),
			fwdconfig.Action(fwdconfig.FlowCounterAction(inCounter)),
		).Build())
	if err != nil {
		return nil, err
	}

	_, err = ri.dataplane.TableEntryAdd(ctx, fwdconfig.TableEntryAddRequest(ri.dataplane.ID(), outputIfaceTable).
		AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE).WithUint64(id))),
			fwdconfig.Action(fwdconfig.TransmitAction(fmt.Sprint(req.GetPortId()))),
			fwdconfig.Action(fwdconfig.FlowCounterAction(outCounter)),
		).Build())
	if err != nil {
		return nil, err
	}

	// Link the interface to a VRF.
	_, err = ri.dataplane.TableEntryAdd(ctx, fwdconfig.TableEntryAddRequest(ri.dataplane.ID(), IngressVRFTable).
		AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_INPUT_IFACE).WithUint64(id))),
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64Value(req.GetVirtualRouterId())),
		).Build())
	if err != nil {
		return nil, err
	}

	// Give the interface a SMAC.
	_, err = ri.dataplane.TableEntryAdd(ctx, fwdconfig.TableEntryAddRequest(ri.dataplane.ID(), SRCMACTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE).WithUint64(id))),
		fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC).WithValue(req.GetSrcMacAddress())),
	).Build())
	if err != nil {
		return nil, err
	}

	return &saipb.CreateRouterInterfaceResponse{Oid: id}, nil
}

func (ri *routerInterface) RemoveRouterInterface(ctx context.Context, req *saipb.RemoveRouterInterfaceRequest) (*saipb.RemoveRouterInterfaceResponse, error) {
	resp := &saipb.GetRouterInterfaceAttributeResponse{}
	err := ri.mgr.PopulateAttributes(&saipb.GetRouterInterfaceAttributeRequest{
		Oid:      req.GetOid(),
		AttrType: []saipb.RouterInterfaceAttr{saipb.RouterInterfaceAttr_ROUTER_INTERFACE_ATTR_PORT_ID},
	}, resp)
	if err != nil {
		return nil, err
	}

	nid, err := ri.dataplane.ObjectNID(ctx, &fwdpb.ObjectNIDRequest{
		ContextId: &fwdpb.ContextId{Id: ri.dataplane.ID()},
		ObjectId:  &fwdpb.ObjectId{Id: fmt.Sprint(resp.GetAttr().GetPortId())},
	})
	if err != nil {
		return nil, err
	}

	_, err = ri.dataplane.TableEntryRemove(ctx, fwdconfig.TableEntryRemoveRequest(ri.dataplane.ID(), inputIfaceTable).
		AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT).WithUint64(nid.GetNid()))),
		).Build())
	if err != nil {
		return nil, err
	}

	_, err = ri.dataplane.TableEntryRemove(ctx, fwdconfig.TableEntryRemoveRequest(ri.dataplane.ID(), outputIfaceTable).
		AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE).WithUint64(req.GetOid()))),
		).Build())
	if err != nil {
		return nil, err
	}

	// Link the interface to a VRF.
	_, err = ri.dataplane.TableEntryRemove(ctx, fwdconfig.TableEntryRemoveRequest(ri.dataplane.ID(), IngressVRFTable).
		AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_INPUT_IFACE).WithUint64(req.GetOid()))),
		).Build())
	if err != nil {
		return nil, err
	}

	// Give the interface a SMAC.
	_, err = ri.dataplane.TableEntryRemove(ctx, fwdconfig.TableEntryRemoveRequest(ri.dataplane.ID(), SRCMACTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE).WithUint64(req.GetOid()))),
	).Build())
	if err != nil {
		return nil, err
	}

	return &saipb.RemoveRouterInterfaceResponse{}, nil
}

func (ri *routerInterface) SetRouterInterfaceAttribute(context.Context, *saipb.SetRouterInterfaceAttributeRequest) (*saipb.SetRouterInterfaceAttributeResponse, error) {
	return &saipb.SetRouterInterfaceAttributeResponse{}, nil
}

func (ri *routerInterface) GetRouterInterfaceStats(ctx context.Context, req *saipb.GetRouterInterfaceStatsRequest) (*saipb.GetRouterInterfaceStatsResponse, error) {
	counters, err := ri.dataplane.FlowCounterQuery(ctx, &fwdpb.FlowCounterQueryRequest{
		ContextId: &fwdpb.ContextId{Id: ri.dataplane.ID()},
		Ids: []*fwdpb.FlowCounterId{{
			ObjectId: &fwdpb.ObjectId{Id: ifaceCounterID(req.GetOid(), true)},
		}, {
			ObjectId: &fwdpb.ObjectId{Id: ifaceCounterID(req.GetOid(), false)},
		}},
	})
	if err != nil {
		return nil, err
	}

	vals := []uint64{}

	for _, counter := range req.GetCounterIds() {
		switch counter {
		case saipb.RouterInterfaceStat_ROUTER_INTERFACE_STAT_IN_OCTETS:
			vals = append(vals, counters.GetCounters()[0].Octets)
		case saipb.RouterInterfaceStat_ROUTER_INTERFACE_STAT_IN_PACKETS:
			vals = append(vals, counters.GetCounters()[0].Packets)
		case saipb.RouterInterfaceStat_ROUTER_INTERFACE_STAT_OUT_OCTETS:
			vals = append(vals, counters.GetCounters()[1].Octets)
		case saipb.RouterInterfaceStat_ROUTER_INTERFACE_STAT_OUT_PACKETS:
			vals = append(vals, counters.GetCounters()[1].Packets)
		default: // TODO: Support more stats.
			vals = append(vals, 0)
		}
	}

	return &saipb.GetRouterInterfaceStatsResponse{Values: vals}, nil
}

// vlanMember contains the info of a VLAN member.
type vlanMember struct {
	Oid    uint64
	PortID uint64
	Vid    uint32
	Mode   saipb.VlanTaggingMode
}

type vlan struct {
	mu sync.Mutex
	saipb.UnimplementedVlanServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
	oidByVId  map[uint32]uint64                 // VID -> VLAN_OID
	vlans     map[uint64]map[uint64]*vlanMember // VLAN_OID -> VLAN_Member_OID (port)
}

func newVlan(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *vlan {
	v := &vlan{
		mgr:       mgr,
		dataplane: dataplane,
		oidByVId:  map[uint32]uint64{},
		vlans:     map[uint64]map[uint64]*vlanMember{},
	}
	saipb.RegisterVlanServer(s, v)
	return v
}

func (vlan *vlan) vidByOid(oid uint64) (uint32, error) {
	for v, o := range vlan.oidByVId {
		if o == oid {
			return v, nil
		}
	}
	return 0, fmt.Errorf("cannot find VLAN with OID %d", oid)
}

func (vlan *vlan) memberByOid(oid uint64) *vlanMember {
	for _, v := range vlan.vlans {
		for mOid, member := range v {
			if mOid == oid {
				return member
			}
		}
	}
	return nil
}

func (vlan *vlan) memberByPortId(oid uint64) *vlanMember {
	for _, v := range vlan.vlans {
		for _, member := range v {
			if member.PortID == oid {
				return member
			}
		}
	}
	return nil
}

func (vlan *vlan) CreateVlan(ctx context.Context, r *saipb.CreateVlanRequest) (*saipb.CreateVlanResponse, error) {
	if _, ok := vlan.oidByVId[r.GetVlanId()]; ok {
		return nil, fmt.Errorf("found existing VLAN %d", r.GetVlanId())
	}
	id := vlan.mgr.NextID()
	req := &saipb.GetSwitchAttributeRequest{Oid: 1, AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_DEFAULT_STP_INST_ID}}
	resp := &saipb.GetSwitchAttributeResponse{}

	if err := vlan.mgr.PopulateAttributes(req, resp); err != nil {
		return nil, err
	}

	attrs := &saipb.VlanAttribute{
		MemberList:                         []uint64{},
		StpInstance:                        resp.Attr.DefaultStpInstId,
		VlanId:                             r.VlanId,
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
	vlan.vlans[id] = map[uint64]*vlanMember{}
	vlan.oidByVId[r.GetVlanId()] = id
	return &saipb.CreateVlanResponse{
		Oid: id,
	}, nil
}

func (vlan *vlan) RemoveVlan(ctx context.Context, r *saipb.RemoveVlanRequest) (*saipb.RemoveVlanResponse, error) {
	vId, err := vlan.vidByOid(r.GetOid())
	if err != nil {
		return nil, fmt.Errorf("cannot find VLAN ID for OID %d", r.GetOid())
	}
	if vId == DefaultVlanId {
		return nil, fmt.Errorf("cannot remove default VLAN")
	}
	for _, v := range vlan.vlans[r.GetOid()] {
		_, err := attrmgr.InvokeAndSave(ctx, vlan.mgr, vlan.RemoveVlanMember, &saipb.RemoveVlanMemberRequest{
			Oid: v.Oid,
		})
		if err != nil {
			return nil, err
		}
	}
	// Update the internal map.
	vlan.mu.Lock()
	delete(vlan.vlans, r.GetOid())
	vlan.mu.Unlock()
	return &saipb.RemoveVlanResponse{}, nil
}

func (vlan *vlan) CreateVlanMember(ctx context.Context, r *saipb.CreateVlanMemberRequest) (*saipb.CreateVlanMemberResponse, error) {
	vOid := r.GetVlanId()
	vId, err := vlan.vidByOid(vOid)
	if err != nil {
		return nil, err
	}
	member := vlan.memberByPortId(r.GetBridgePortId()) // Keep the vlan member if this port was assigned to any VLAN.
	mOid := vlan.mgr.NextID()
	nid, err := vlan.dataplane.ObjectNID(ctx, &fwdpb.ObjectNIDRequest{
		ContextId: &fwdpb.ContextId{Id: vlan.dataplane.ID()},
		ObjectId:  &fwdpb.ObjectId{Id: fmt.Sprint(r.GetBridgePortId())},
	})
	if err != nil {
		slog.InfoContext(ctx, "Failed to find NID for port", "bridge_port", r.GetBridgePortId(), "err", err)
		return nil, err
	}
	vlanReq := fwdconfig.TableEntryAddRequest(vlan.dataplane.ID(), VlanTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT).WithUint64(nid.GetNid())))).Build()
	vlanReq.Entries[0].Actions = []*fwdpb.ActionDesc{
		fwdconfig.Action(fwdconfig.EncapAction(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN)).Build(),
		fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG).WithUint64Value(uint64(vId))).Build(),
	}
	if _, err := vlan.dataplane.TableEntryAdd(ctx, vlanReq); err != nil {
		return nil, err
	}
	// Update the attributes and intenal data.
	vlanAttrReq := &saipb.GetVlanAttributeRequest{Oid: vOid, AttrType: []saipb.VlanAttr{saipb.VlanAttr_VLAN_ATTR_MEMBER_LIST}}
	vlanAttrResp := &saipb.GetVlanAttributeResponse{}
	if err := vlan.mgr.PopulateAttributes(vlanAttrReq, vlanAttrResp); err != nil {
		return nil, err
	}
	vlanAttrResp.GetAttr().MemberList = append(vlanAttrResp.GetAttr().MemberList, mOid)
	vlan.mgr.StoreAttributes(vOid, vlanAttrResp.GetAttr())
	vlan.mu.Lock()
	vlan.vlans[vOid][mOid] = &vlanMember{Oid: mOid, PortID: r.GetBridgePortId(), Vid: vId, Mode: r.GetVlanTaggingMode()}
	vlan.mu.Unlock()
	if member != nil {
		preVlanOid := vlan.oidByVId[member.Vid]
		vlanAttrReq = &saipb.GetVlanAttributeRequest{Oid: preVlanOid, AttrType: []saipb.VlanAttr{saipb.VlanAttr_VLAN_ATTR_MEMBER_LIST}}
		vlanAttrResp = &saipb.GetVlanAttributeResponse{}
		if err := vlan.mgr.PopulateAttributes(vlanAttrReq, vlanAttrResp); err != nil {
			return nil, err
		}
		newMemList := []uint64{}
		for _, i := range vlanAttrResp.Attr.GetMemberList() {
			if i != member.Oid {
				newMemList = append(newMemList, i)
			}
		}
		vlanAttrResp.GetAttr().MemberList = newMemList
		vlan.mgr.StoreAttributes(preVlanOid, vlanAttrResp.GetAttr())
		// Update internal map.
		vlan.mu.Lock()
		delete(vlan.vlans[preVlanOid], member.Oid)
		vlan.mu.Unlock()
	}
	return &saipb.CreateVlanMemberResponse{Oid: mOid}, nil
}

func (vlan *vlan) RemoveVlanMember(ctx context.Context, r *saipb.RemoveVlanMemberRequest) (*saipb.RemoveVlanMemberResponse, error) {
	member := vlan.memberByOid(r.GetOid())
	if member == nil {
		return nil, fmt.Errorf("cannot find member with OID %d", r.GetOid())
	}
	nid, err := vlan.dataplane.ObjectNID(ctx, &fwdpb.ObjectNIDRequest{
		ContextId: &fwdpb.ContextId{Id: vlan.dataplane.ID()},
		ObjectId:  &fwdpb.ObjectId{Id: fmt.Sprint(member.PortID)},
	})
	if err != nil {
		slog.InfoContext(ctx, "Failed to find NID for port", "bridge_port", member.PortID, "err", err)
		return nil, err
	}
	if _, err := vlan.dataplane.TableEntryRemove(ctx, fwdconfig.TableEntryRemoveRequest(vlan.dataplane.ID(), VlanTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT).WithUint64(nid.GetNid())))).Build()); err != nil {
		return nil, err
	}
	return &saipb.RemoveVlanMemberResponse{}, nil
}

func (vlan *vlan) CreateVlanMembers(ctx context.Context, r *saipb.CreateVlanMembersRequest) (*saipb.CreateVlanMembersResponse, error) {
	resp := &saipb.CreateVlanMembersResponse{}
	for _, req := range r.GetReqs() {
		res, err := attrmgr.InvokeAndSave(ctx, vlan.mgr, vlan.CreateVlanMember, req)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, res)
	}
	return resp, nil
}

func (vlan *vlan) Reset() {
	slog.Info("resetting vlan")
	vlan.oidByVId = map[uint32]uint64{}
	vlan.vlans = map[uint64]map[uint64]*vlanMember{}
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

type hash struct {
	saipb.UnimplementedHashServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newHash(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *hash {
	m := &hash{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterHashServer(s, m)
	return m
}

func convertHashFields(list []saipb.NativeHashField) ([]*fwdpb.PacketFieldId, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("got 0 hash fields")
	}

	fields := []*fwdpb.PacketFieldId{}
	for _, field := range list {
		switch field {
		case saipb.NativeHashField_NATIVE_HASH_FIELD_SRC_IP, saipb.NativeHashField_NATIVE_HASH_FIELD_SRC_IPV4, saipb.NativeHashField_NATIVE_HASH_FIELD_SRC_IPV6:
			fields = append(fields, &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC}})
		case saipb.NativeHashField_NATIVE_HASH_FIELD_DST_IP, saipb.NativeHashField_NATIVE_HASH_FIELD_DST_IPV4, saipb.NativeHashField_NATIVE_HASH_FIELD_DST_IPV6:
			fields = append(fields, &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}})
		case saipb.NativeHashField_NATIVE_HASH_FIELD_L4_SRC_PORT:
			fields = append(fields, &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC}})
		case saipb.NativeHashField_NATIVE_HASH_FIELD_L4_DST_PORT:
			fields = append(fields, &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST}})
		case saipb.NativeHashField_NATIVE_HASH_FIELD_IPV6_FLOW_LABEL:
			fields = append(fields, &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP6_FLOW}})
		default:
			return nil, fmt.Errorf("unsupported hash field: %v", field)
		}
	}
	return fields, nil
}

func (h *hash) CreateHash(_ context.Context, req *saipb.CreateHashRequest) (*saipb.CreateHashResponse, error) {
	id := h.mgr.NextID()

	// Creating a hash doesn't affect the forwarding pipeline, just validate the arguments.
	_, err := convertHashFields(req.GetNativeHashFieldList())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &saipb.CreateHashResponse{Oid: id}, nil
}

type virtualRouter struct {
	saipb.UnimplementedVirtualRouterServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newVirtualRouter(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *virtualRouter {
	vr := &virtualRouter{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterVirtualRouterServer(s, vr)
	return vr
}

func (vr *virtualRouter) CreateVirtualRouter(context.Context, *saipb.CreateVirtualRouterRequest) (*saipb.CreateVirtualRouterResponse, error) {
	id := vr.mgr.NextID()

	return &saipb.CreateVirtualRouterResponse{Oid: id}, nil
}

func (vr *virtualRouter) RemoveVirtualRouter(context.Context, *saipb.RemoveVirtualRouterRequest) (*saipb.RemoveVirtualRouterResponse, error) {
	return &saipb.RemoveVirtualRouterResponse{}, nil
}

func (vr *virtualRouter) SetVirtualRouterAttribute(context.Context, *saipb.SetVirtualRouterAttributeRequest) (*saipb.SetVirtualRouterAttributeResponse, error) {
	return &saipb.SetVirtualRouterAttributeResponse{}, nil
}

type rpfGroup struct {
	saipb.UnimplementedRpfGroupServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newRpfGroup(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *rpfGroup {
	rpf := &rpfGroup{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterRpfGroupServer(s, rpf)
	return rpf
}

func (rpf *rpfGroup) CreateRpfGroup(context.Context, *saipb.CreateRpfGroupRequest) (*saipb.CreateRpfGroupResponse, error) {
	id := rpf.mgr.NextID()

	return &saipb.CreateRpfGroupResponse{Oid: id}, nil
}

func (rpf *rpfGroup) CreateRpfGroupMember(context.Context, *saipb.CreateRpfGroupMemberRequest) (*saipb.CreateRpfGroupMemberResponse, error) {
	id := rpf.mgr.NextID()

	return &saipb.CreateRpfGroupMemberResponse{Oid: id}, nil
}

func (rpf *rpfGroup) RemoveRpfGroup(context.Context, *saipb.RemoveRpfGroupRequest) (*saipb.RemoveRpfGroupResponse, error) {
	return &saipb.RemoveRpfGroupResponse{}, nil
}

func (rpf *rpfGroup) RemoveRpfGroupMember(context.Context, *saipb.RemoveRpfGroupMemberRequest) (*saipb.RemoveRpfGroupMemberResponse, error) {
	return &saipb.RemoveRpfGroupMemberResponse{}, nil
}
