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

//nolint:revive // Acl is incorrect, but it comes from generated code.
package saiserver

import (
	"context"
	"fmt"
	"sync"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type aclDataplaneAPI interface {
	ID() string
	TableCreate(context.Context, *fwdpb.TableCreateRequest) (*fwdpb.TableCreateReply, error)
	TableEntryAdd(context.Context, *fwdpb.TableEntryAddRequest) (*fwdpb.TableEntryAddReply, error)
}

// tableLocation indentifies the location of an acl table by the group, bank, member id.
type tableLocation struct {
	groupID  string
	bank     int
	memberID uint64
}

type acl struct {
	saipb.UnimplementedAclServer
	mgr       *attrmgr.AttrMgr
	dataplane aclDataplaneAPI
	// tableToLocation maps the acl table id to the lucius flow table and bank.
	tableToLocation     map[uint64]tableLocation
	groupNextFreeBankMu sync.Mutex
	// groupNextFreeBank contains the next free bank for a group.
	groupNextFreeBank map[uint64]int
}

func newACL(mgr *attrmgr.AttrMgr, dataplane aclDataplaneAPI, s *grpc.Server) *acl {
	a := &acl{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterAclServer(s, a)
	return a
}

// CreateAclTableGroup creates a lucius flow table, where the group members are banks in the flow table.
func (a *acl) CreateAclTableGroup(ctx context.Context, req *saipb.CreateAclTableGroupRequest) (*saipb.CreateAclTableGroupResponse, error) {
	id := a.mgr.NextID()

	stage := req.GetAclStage()
	typ := req.GetType()
	bind := req.GetAclBindPointTypeList()
	if stage != saipb.AclStage_ACL_STAGE_EGRESS && stage != saipb.AclStage_ACL_STAGE_PRE_INGRESS && stage != saipb.AclStage_ACL_STAGE_INGRESS {
		return nil, status.Errorf(codes.InvalidArgument, "invalid stage type: %v", stage)
	}
	if typ != saipb.AclTableGroupType_ACL_TABLE_GROUP_TYPE_PARALLEL {
		return nil, status.Errorf(codes.InvalidArgument, "invalid type: %v", typ)
	}
	if len(bind) != 1 && bind[0] == saipb.AclBindPointType_ACL_BIND_POINT_TYPE_SWITCH {
		return nil, status.Errorf(codes.InvalidArgument, "invalid bind configuration: %v", typ)
	}

	tReq := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: a.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_FLOW,
			TableId: &fwdpb.TableId{
				ObjectId: &fwdpb.ObjectId{
					Id: fmt.Sprint(id),
				},
			},
			Table: &fwdpb.TableDesc_Flow{
				Flow: &fwdpb.FlowTableDesc{
					BankCount: 1,
				},
			},
		},
	}
	if _, err := a.dataplane.TableCreate(ctx, tReq); err != nil {
		return nil, err
	}

	return &saipb.CreateAclTableGroupResponse{Oid: 1}, nil
}

// CreateAclTableGroupMember stores the acl table id to its corresponding lucius flow table and bank.
func (a *acl) CreateAclTableGroupMember(_ context.Context, req *saipb.CreateAclTableGroupMemberRequest) (*saipb.CreateAclTableGroupMemberResponse, error) {
	groupID := req.GetAclTableGroupId()
	tableID := req.GetAclTableId()
	memberID := a.mgr.NextID()

	a.groupNextFreeBankMu.Lock()
	bank := a.groupNextFreeBank[groupID]
	a.groupNextFreeBank[groupID] = bank + 1
	a.groupNextFreeBankMu.Unlock()
	a.tableToLocation[tableID] = tableLocation{
		groupID:  fmt.Sprint(groupID),
		bank:     bank,
		memberID: memberID,
	}

	return &saipb.CreateAclTableGroupMemberResponse{Oid: memberID}, nil
}

// CreateAclTable is noop as the table is already created in the group.
func (a *acl) CreateAclTable(context.Context, *saipb.CreateAclTableRequest) (*saipb.CreateAclTableResponse, error) {
	id := a.mgr.NextID()
	return &saipb.CreateAclTableResponse{Oid: id}, nil
}

// CreateAclEntry adds an entry in the a bank.
func (a *acl) CreateAclEntry(ctx context.Context, req *saipb.CreateAclEntryRequest) (*saipb.CreateAclEntryResponse, error) {
	id := a.mgr.NextID()
	gb, ok := a.tableToLocation[req.GetTableId()]
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition, "table is not member of a group")
	}

	aReq := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: a.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: gb.groupID}},
		EntryDesc: &fwdpb.EntryDesc{
			Entry: &fwdpb.EntryDesc_Flow{
				Flow: &fwdpb.FlowEntryDesc{
					Priority: req.GetPriority(),
					Id:       uint32(id),
					Bank:     uint32(gb.bank),
				},
			},
		},
	}
	switch {
	case req.GetFieldDstIp() != nil:
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
			Bytes:   req.GetFieldDstIp().GetDataIp(),
			Masks:   req.GetFieldDstIp().GetMaskIp(),
		})
	}
	switch {
	case req.ActionSetVrf != nil:
		aReq.Actions = append(aReq.Actions,
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).
				WithUint64Value(req.GetActionSetVrf().GetOid())).Build())
	}
	if _, err := a.dataplane.TableEntryAdd(ctx, aReq); err != nil {
		return nil, err
	}

	return &saipb.CreateAclEntryResponse{Oid: id}, nil
}
