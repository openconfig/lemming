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
	groups    map[uint64]map[uint64]*l2mcGroupMember // map[group id][member oid]
}

func newL2mcGroup(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *l2mcGroup {
	mg := &l2mcGroup{
		mgr:       mgr,
		dataplane: dataplane,
		groups:    map[uint64]map[uint64]*l2mcGroupMember{},
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

func (mg *l2mcGroup) portNidFromBrirdgeId(ctx context.Context, outputId uint64) (uint64, error) {
	req := &saipb.GetBridgePortAttributeRequest{Oid: outputId, AttrType: []saipb.BridgePortAttr{saipb.BridgePortAttr_BRIDGE_PORT_ATTR_PORT_ID}}
	resp := &saipb.GetBridgePortAttributeResponse{}
	if err := mg.mgr.PopulateAttributes(req, resp); err != nil {
		return 0, fmt.Errorf("failed to populate OutputId (oid=%d): %v", outputId, err)
	}
	if resp.GetAttr().PortId == nil {
		return 0, fmt.Errorf("cannot find portId for bridge port %q", outputId)
	}
	return resp.GetAttr().GetPortId(), nil
}

// updateGroupMember updates the member of a L2MC group.
// If m is nil, remove mid from the group(key: group id), otherwise add m to groups with mid as the key.
func (mg *l2mcGroup) updateGroupMember(ctx context.Context, gid, mid uint64, m *l2mcGroupMember) error {
	if m == nil {
		// Remove the member.
		delete(mg.groups[gid], mid)
	} else {
		// Add a member.
		mg.groups[gid][mid] = m
		attr := &saipb.L2McGroupMemberAttribute{
			L2McGroupId:  &m.groupId,
			L2McOutputId: &m.outputId,
		}
		mg.mgr.StoreAttributes(mid, attr)
	}
	mList := []uint64{}
	for _, m := range mg.groups[gid] {
		mList = append(mList, m.oid)
	}
	gAttr := &saipb.L2McGroupAttribute{}
	gAttr.L2McMemberList = mList
	cnt := uint32(len(mList))
	gAttr.L2McOutputCount = &cnt
	mg.mgr.StoreAttributes(gid, gAttr)

	var actions []*fwdpb.ActionDesc
	for _, member := range mg.groups[gid] {
		portId, err := mg.portNidFromBrirdgeId(ctx, member.outputId)
		if err != nil {
			return err
		}
		actions = append(actions, &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_MIRROR,
			Action: &fwdpb.ActionDesc_Mirror{
				Mirror: &fwdpb.MirrorActionDesc{
					PortId: &fwdpb.PortId{
						ObjectId: &fwdpb.ObjectId{
							Id: fmt.Sprint(portId),
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
	if mg.groups[req.GetL2McGroupId()] == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "cannot find L2MC group with group ID=%d", req.GetL2McGroupId())
	}
	if _, err := mg.portNidFromBrirdgeId(ctx, req.GetL2McOutputId()); err != nil {
		return nil, err
	}
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
		return nil, status.Errorf(codes.FailedPrecondition, "cannot find L2MC group with group ID=%d", req.GetOid())
	}
	delete(mg.groups, req.GetOid())
	if _, err := mg.dataplane.TableEntryRemove(ctx, fwdconfig.TableEntryRemoveRequest(mg.dataplane.ID(), L2MCGroupTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID).WithUint64(req.GetOid())))).Build()); err != nil {
		return nil, err
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
		return nil, status.Errorf(codes.FailedPrecondition, "cannot find L2MC group member with OID %d", req.GetOid())
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
