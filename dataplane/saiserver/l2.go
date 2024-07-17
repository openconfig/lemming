// Copyright 2024 Google LLC
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
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

type l2mcGroupMember struct {
	oid      uint64
	groupId  uint64
	outputId uint64
}

type l2mcGroup struct {
	saipb.UnimplementedL2McGroupServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
	groups    map[uint64]map[uint64]*l2mcGroupMember // OID -> map of Port
}

func newL2mcGroup(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *l2mcGroup {
	mg := &l2mcGroup{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterL2McGroupServer(s, mg)
	return mg
}

func (mg *l2mcGroup) CreateL2McGroup(ctx context.Context, req *saipb.CreateL2McGroupRequest) (*saipb.CreateL2McGroupResponse, error) {
	id := mg.mgr.NextID()
	// Update internal data
	mg.groups[id] = map[uint64]*l2mcGroupMember{}
	// Update attributes
	l2mcAttrs := &saipb.L2McGroupAttribute{
		L2McOutputCount: proto.Uint32(0),
		L2McMemberList:  []uint64{},
	}
	mg.mgr.StoreAttributes(id, l2mcAttrs)
	return &saipb.CreateL2McGroupResponse{Oid: id}, nil
}

func (mg *l2mcGroup) CreateL2McGroupMember(ctx context.Context, req *saipb.CreateL2McGroupMemberRequest) (*saipb.CreateL2McGroupMemberResponse, error) {
	id := mg.mgr.NextID()
	// Update table entry
	r := fwdconfig.TableEntryAddRequest(mg.dataplane.ID(), L2MCGroupTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID).WithUint64(req.GetL2McGroupId())))).Build()
	for _, p := range mg.groups {
		r.Entries[0].Actions = append(r.Entries[0].Actions, &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_MIRROR,
			Action: &fwdpb.ActionDesc_Mirror{
				Mirror: &fwdpb.MirrorActionDesc{
					PortId: &fwdpb.PortId{
						ObjectId: &fwdpb.ObjectId{
							Id: fmt.Sprint(p),
						},
					},
				},
			},
		})
	}
	if _, err := mg.dataplane.TableEntryAdd(ctx, r); err != nil {
		return nil, err
	}
	// Update internal data
	if mg.groups[req.GetL2McGroupId()] == nil {
		return nil, fmt.Errorf("L2MC group id %q not found", req.GetL2McGroupId())
	}
	if _, ok := mg.groups[req.GetL2McGroupId()][req.GetL2McOutputId()]; ok {
		return nil, fmt.Errorf("found existing member %q", req.GetL2McOutputId())
	}
	mg.groups[req.GetL2McGroupId()][req.GetL2McOutputId()] = &l2mcGroupMember{
		oid:      id,
		outputId: req.GetL2McOutputId(),
		groupId:  req.GetL2McGroupId(),
	}
	// Update L2MC Group member attributes
	attr := &saipb.L2McGroupMemberAttribute{
		L2McGroupId:  req.L2McGroupId,
		L2McOutputId: req.L2McOutputId,
	}
	mg.mgr.StoreAttributes(id, attr)
	// Update L2MC Group Attributes.
	gReq := &saipb.GetL2McGroupAttributeRequest{Oid: req.GetL2McGroupId(), AttrType: []saipb.L2McGroupAttr{saipb.L2McGroupAttr_L2MC_GROUP_ATTR_L2MC_MEMBER_LIST, saipb.L2McGroupAttr_L2MC_GROUP_ATTR_L2MC_OUTPUT_COUNT}}
	gResp := &saipb.GetL2McGroupAttributeResponse{}
	if err := mg.mgr.PopulateAttributes(gReq, gResp); err != nil {
		return nil, err
	}
	gAttrs := gResp.GetAttr()
	gAttrs.L2McMemberList = append(gAttrs.GetL2McMemberList(), id)
	*gAttrs.L2McOutputCount += 1
	mg.mgr.StoreAttributes(req.GetL2McGroupId(), gAttrs)
	// Update Switch Attributes.
	swReq := &saipb.GetSwitchAttributeRequest{Oid: req.GetSwitch(), AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_AVAILABLE_L2MC_ENTRY}}
	swResp := &saipb.GetSwitchAttributeResponse{}
	if err := mg.mgr.PopulateAttributes(swReq, swResp); err != nil {
		return nil, err
	}
	attrs := swResp.GetAttr()
	*attrs.AvailableL2McEntry = attrs.GetAvailableL2McEntry() - 1
	mg.mgr.StoreAttributes(req.GetSwitch(), attrs)
	return &saipb.CreateL2McGroupMemberResponse{Oid: id}, nil
}

func (mg *l2mcGroup) RemoveL2McGroup(ctx context.Context, req *saipb.RemoveL2McGroupRequest) (*saipb.RemoveL2McGroupResponse, error) {
	if mg.groups[req.GetOid()] == nil {
		return nil, fmt.Errorf("cannot find group %q", req.GetOid())
	}
	// Remove all members in the group
	for _, p := range mg.groups[req.GetOid()] {
		_, err := attrmgr.InvokeAndSave(ctx, mg.mgr, mg.RemoveL2McGroupMember, &saipb.RemoveL2McGroupMemberRequest{
			Oid: p.oid,
		})
		if err != nil {
			return nil, err
		}
	}
	// Update internal data
	delete(mg.groups, req.GetOid())
	return &saipb.RemoveL2McGroupResponse{}, nil
}

func (mg *l2mcGroup) RemoveL2McGroupMember(ctx context.Context, req *saipb.RemoveL2McGroupMemberRequest) (*saipb.RemoveL2McGroupMemberResponse, error) {
	// Remove table entry.
	r := fwdconfig.TableEntryRemoveRequest(mg.dataplane.ID(), L2MCGroupTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID).WithUint64(req.GetOid())))).Build()
	if _, err := mg.dataplane.TableEntryRemove(ctx, r); err != nil {
		return nil, err
	}
	// Update L2MC group attributes
	locateMember := func(oid uint64) *l2mcGroupMember {
		for _, members := range mg.groups {
			for k, v := range members {
				if k == oid {
					return v
				}
			}
		}
		return nil
	}
	m := locateMember(req.GetOid())
	if m == nil {
		return nil, fmt.Errorf("cannot find member with OID %d", req.GetOid())
	}
	gReq := &saipb.GetL2McGroupAttributeRequest{Oid: m.groupId, AttrType: []saipb.L2McGroupAttr{saipb.L2McGroupAttr_L2MC_GROUP_ATTR_L2MC_MEMBER_LIST, saipb.L2McGroupAttr_L2MC_GROUP_ATTR_L2MC_OUTPUT_COUNT}}
	gResp := &saipb.GetL2McGroupAttributeResponse{}
	if err := mg.mgr.PopulateAttributes(gReq, gResp); err != nil {
		return nil, err
	}
	gAttrs := gResp.GetAttr()
	newMemList := []uint64{}
	for _, i := range gAttrs.GetL2McMemberList() {
		if i != req.GetOid() {
			newMemList = append(newMemList, i)
		}
	}
	gAttrs.L2McMemberList = newMemList
	*gAttrs.L2McOutputCount -= 1
	mg.mgr.StoreAttributes(m.groupId, gAttrs)
	// Update Switch Attributes.
	swReq := &saipb.GetSwitchAttributeRequest{Oid: 1, AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_AVAILABLE_L2MC_ENTRY}}
	swResp := &saipb.GetSwitchAttributeResponse{}
	if err := mg.mgr.PopulateAttributes(swReq, swResp); err != nil {
		return nil, err
	}
	attrs := swResp.GetAttr()
	*attrs.AvailableL2McEntry = attrs.GetAvailableL2McEntry() + 1
	mg.mgr.StoreAttributes(1, attrs)
	// Update internal data
	delete(mg.groups[m.groupId], req.GetOid())
	return &saipb.RemoveL2McGroupMemberResponse{}, nil
}

type l2mc struct {
	saipb.UnimplementedL2McServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

// TODO: Implement this.
func newL2mc(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *l2mc {
	m := &l2mc{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterL2McServer(s, m)
	return m
}

// TODO: Implement this.
func (m *l2mc) CreateL2McEntry(context.Context, *saipb.CreateL2McEntryRequest) (*saipb.CreateL2McEntryResponse, error) {
	return &saipb.CreateL2McEntryResponse{}, nil
}

// TODO: Implement this.
func (m *l2mc) RemoveL2McEntry(context.Context, *saipb.RemoveL2McEntryRequest) (*saipb.RemoveL2McEntryResponse, error) {
	return &saipb.RemoveL2McEntryResponse{}, nil
}

// TODO: Implement this.
func (m *l2mc) SetL2McEntryAttribute(context.Context, *saipb.SetL2McEntryAttributeRequest) (*saipb.SetL2McEntryAttributeResponse, error) {
	return &saipb.SetL2McEntryAttributeResponse{}, nil
}
