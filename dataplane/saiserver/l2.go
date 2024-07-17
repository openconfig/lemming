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
	"encoding/binary"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

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
	groups    map[uint64]map[uint64]*l2mcGroupMember // map[group id][member oid]
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

// updateGroupMember updates the member of a L2MC group.
// If m is nil, remove mid from the group(key: group id), otherwise add m to groups with mid as the key.
func (mg *l2mcGroup) updateGroupMember(ctx context.Context, gid, mid uint64, m *l2mcGroupMember) error {
	if mg.groups[gid] == nil {
		return status.Errorf(codes.FailedPrecondition, "L2MC group %d does not exist", gid)
	}
	if m == nil {
		// Remove the member.
		delete(mg.groups[gid], mid)
		if len(mg.groups[gid]) == 0 {
			delete(mg.groups, gid)
		}
		gReq := &saipb.GetL2McGroupAttributeRequest{Oid: gid, AttrType: []saipb.L2McGroupAttr{saipb.L2McGroupAttr_L2MC_GROUP_ATTR_L2MC_MEMBER_LIST, saipb.L2McGroupAttr_L2MC_GROUP_ATTR_L2MC_OUTPUT_COUNT}}
		gResp := &saipb.GetL2McGroupAttributeResponse{}
		if err := mg.mgr.PopulateAttributes(gReq, gResp); err != nil {
			return err
		}
		gAttrs := gResp.GetAttr()
		newMemList := []uint64{}
		for _, i := range gAttrs.GetL2McMemberList() {
			if i != mid {
				newMemList = append(newMemList, i)
			}
		}
		gAttrs.L2McMemberList = newMemList
		*gAttrs.L2McOutputCount -= 1
		mg.mgr.StoreAttributes(gid, gAttrs)
	} else {
		// Add a member.
		mg.groups[gid][mid] = m
		attr := &saipb.L2McGroupMemberAttribute{
			L2McGroupId:  &m.groupId,
			L2McOutputId: &m.outputId,
		}
		mg.mgr.StoreAttributes(mid, attr)
		// Update L2MC Group Attributes.
		gReq := &saipb.GetL2McGroupAttributeRequest{Oid: gid, AttrType: []saipb.L2McGroupAttr{saipb.L2McGroupAttr_L2MC_GROUP_ATTR_L2MC_MEMBER_LIST, saipb.L2McGroupAttr_L2MC_GROUP_ATTR_L2MC_OUTPUT_COUNT}}
		gResp := &saipb.GetL2McGroupAttributeResponse{}
		if err := mg.mgr.PopulateAttributes(gReq, gResp); err != nil {
			return err
		}
		gAttrs := gResp.GetAttr()
		gAttrs.L2McMemberList = append(gAttrs.GetL2McMemberList(), mid)
		*gAttrs.L2McOutputCount += 1
		mg.mgr.StoreAttributes(gid, gAttrs)
	}
	var actions []*fwdpb.ActionDesc
	for _, member := range mg.groups[gid] {
		actions = append(actions, &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_MIRROR,
			Action: &fwdpb.ActionDesc_Mirror{
				Mirror: &fwdpb.MirrorActionDesc{
					PortId: &fwdpb.PortId{
						ObjectId: &fwdpb.ObjectId{
							Id: fmt.Sprint(member),
						},
					},
				},
			},
		})
	}
	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: mg.dataplane.ID()},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: L2MCGroupTable,
			},
		},
		Entries: []*fwdpb.TableEntryAddRequest_Entry{{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes: binary.BigEndian.AppendUint64(nil, gid),
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID,
								},
							},
						}},
					},
				},
			},
			Actions: actions,
		}},
	}
	_, err := mg.dataplane.TableEntryAdd(ctx, entries)
	return err
}

func (mg *l2mcGroup) CreateL2McGroupMember(ctx context.Context, req *saipb.CreateL2McGroupMemberRequest) (*saipb.CreateL2McGroupMemberResponse, error) {
	id := mg.mgr.NextID()
	if err := mg.updateGroupMember(ctx, req.GetL2McGroupId(), id, &l2mcGroupMember{
		oid:      id,
		outputId: req.GetL2McOutputId(),
		groupId:  req.GetL2McGroupId(),
	}); err != nil {
		return nil, err
	}
	return &saipb.CreateL2McGroupMemberResponse{Oid: id}, nil
}

func (mg *l2mcGroup) RemoveL2McGroup(ctx context.Context, req *saipb.RemoveL2McGroupRequest) (*saipb.RemoveL2McGroupResponse, error) {
	if mg.groups[req.GetOid()] == nil {
		return nil, fmt.Errorf("cannot find L2MC group %q", req.GetOid())
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
	return &saipb.RemoveL2McGroupResponse{}, nil
}

func (mg *l2mcGroup) RemoveL2McGroupMember(ctx context.Context, req *saipb.RemoveL2McGroupMemberRequest) (*saipb.RemoveL2McGroupMemberResponse, error) {
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
	if err := mg.updateGroupMember(ctx, m.groupId, m.oid, nil); err != nil {
		return nil, err
	}
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
