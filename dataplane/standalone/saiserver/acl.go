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
	"encoding/binary"
	"fmt"
	"math"
	"sync"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// tableLocation indentifies the location of an acl table by the group, bank, member id.
type tableLocation struct {
	groupID  string
	bank     int
	memberID uint64
}

type acl struct {
	saipb.UnimplementedAclServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
	// tableToLocation maps the acl table id to the lucius flow table and bank.
	tableToLocation     map[uint64]tableLocation
	groupNextFreeBankMu sync.Mutex
	// groupNextFreeBank contains the next free bank for a group.
	groupNextFreeBank map[uint64]int
}

func newACL(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *acl {
	a := &acl{
		mgr:               mgr,
		dataplane:         dataplane,
		tableToLocation:   make(map[uint64]tableLocation),
		groupNextFreeBank: make(map[uint64]int),
	}
	saipb.RegisterAclServer(s, a)
	return a
}

// CreateAclTableGroup creates a lucius flow table, where the group members are banks in the flow table.
func (a *acl) CreateAclTableGroup(ctx context.Context, req *saipb.CreateAclTableGroupRequest) (*saipb.CreateAclTableGroupResponse, error) {
	id := a.mgr.NextID()

	stage := req.GetAclStage()
	typ := req.GetType()
	if stage != saipb.AclStage_ACL_STAGE_EGRESS && stage != saipb.AclStage_ACL_STAGE_PRE_INGRESS && stage != saipb.AclStage_ACL_STAGE_INGRESS {
		return nil, status.Errorf(codes.InvalidArgument, "invalid stage type: %v", stage)
	}
	if typ != saipb.AclTableGroupType_ACL_TABLE_GROUP_TYPE_PARALLEL {
		return nil, status.Errorf(codes.InvalidArgument, "invalid type: %v", typ)
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

	return &saipb.CreateAclTableGroupResponse{Oid: id}, nil
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
	if req.GetFieldDstIp() != nil {
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
			Bytes:   req.GetFieldDstIp().GetDataIp(),
			Masks:   req.GetFieldDstIp().GetMaskIp(),
		})
	}
	if req.GetFieldInPort() != nil {
		nid, ok := a.dataplane.PortIDToNID(fmt.Sprint(req.FieldInPort.GetDataOid()))
		if !ok {
			return nil, fmt.Errorf("unknown port with id: %v", req.FieldInPort.GetDataOid())
		}
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT}},
			Bytes:   binary.BigEndian.AppendUint64(nil, nid),
			Masks:   binary.BigEndian.AppendUint64(nil, math.MaxUint64),
		})
	}
	if req.GetFieldAclIpType() != nil { // Use the EtherType header to match against specific protocols.
		fieldMask := &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE}},
			Bytes:   binary.BigEndian.AppendUint16(nil, 0x0000),
			Masks:   binary.BigEndian.AppendUint16(nil, 0xFFFF),
		}
		switch t := req.GetFieldAclIpType().GetDataIpType(); t {
		case saipb.AclIpType_ACL_IP_TYPE_ANY:
			fieldMask.Masks = binary.BigEndian.AppendUint16(nil, 0x0000) // Match any EtherType.
		case saipb.AclIpType_ACL_IP_TYPE_IPV4ANY:
			fieldMask.Bytes = binary.BigEndian.AppendUint16(nil, 0x0800) // Match IPv4.
		case saipb.AclIpType_ACL_IP_TYPE_IPV6ANY:
			fieldMask.Bytes = binary.BigEndian.AppendUint16(nil, 0x86DD) // Match IPv6.
		case saipb.AclIpType_ACL_IP_TYPE_ARP:
			fieldMask.Bytes = binary.BigEndian.AppendUint16(nil, 0x0806) // Match ARP.
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unspporrted ACL_IP_TYPE: %v", t)
		}
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, fieldMask)
	}
	if req.GetFieldDscp() != nil {
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS}},
			Bytes:   []byte{byte(req.GetFieldDscp().GetDataUint())},
			Masks:   []byte{byte(req.GetFieldDscp().GetMaskUint())},
		})
	}
	if req.GetFieldDstIpv6Word3() != nil { // Word3 is supposed to match the 127:96 bits of the IP, assume the caller is masking this correctly put the whole IP in the table.
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
			Bytes:   req.GetFieldDstIpv6Word3().GetDataIp(),
			Masks:   req.GetFieldDstIpv6Word3().GetMaskIp(),
		})
	}
	if req.GetFieldDstIpv6Word2() != nil { // Word2 is supposed to match the  95:64 bits of the IP, assume the caller is masking this correctly put the whole IP in the table.
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
			Bytes:   req.GetFieldDstIpv6Word2().GetDataIp(),
			Masks:   req.GetFieldDstIpv6Word2().GetMaskIp(),
		})
	}
	if req.GetFieldDstMac() != nil {
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST}},
			Bytes:   req.GetFieldDstMac().GetDataMac(),
			Masks:   req.GetFieldDstMac().GetMaskMac(),
		})
	}
	if req.GetFieldEtherType() != nil {
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE}},
			Bytes:   binary.BigEndian.AppendUint16(nil, uint16(req.GetFieldEtherType().GetDataUint())),
			Masks:   binary.BigEndian.AppendUint16(nil, uint16(req.GetFieldEtherType().GetMaskUint())),
		})
	}
	if req.GetFieldIcmpv6Type() != nil {
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ICMP_TYPE}},
			Bytes:   []byte{byte(req.GetFieldIcmpv6Type().GetDataUint())},
			Masks:   []byte{byte(req.GetFieldIcmpv6Type().GetMaskUint())},
		})
	}
	if req.GetFieldIpProtocol() != nil {
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO}},
			Bytes:   []byte{byte(req.GetFieldIpProtocol().GetDataUint())},
			Masks:   []byte{byte(req.GetFieldIpProtocol().GetMaskUint())},
		})
	}
	if req.GetFieldL4DstPort() != nil {
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST}},
			Bytes:   binary.BigEndian.AppendUint16(nil, uint16(req.GetFieldL4DstPort().GetDataUint())),
			Masks:   binary.BigEndian.AppendUint16(nil, uint16(req.GetFieldL4DstPort().GetMaskUint())),
		})
	}
	if req.GetFieldSrcMac() != nil {
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC}},
			Bytes:   req.GetFieldSrcMac().GetDataMac(),
			Masks:   req.GetFieldSrcMac().GetMaskMac(),
		})
	}
	if req.GetFieldTtl() != nil {
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP}},
			Bytes:   []byte{byte(req.GetFieldTtl().GetDataUint())},
			Masks:   []byte{byte(req.GetFieldTtl().GetMaskUint())},
		})
	}
	if len(aReq.EntryDesc.GetFlow().Fields) == 0 {
		return nil, status.Error(codes.InvalidArgument, "either no fields or not unsupports fields in entry req")
	}

	if req.ActionSetVrf != nil {
		aReq.Actions = append(aReq.Actions,
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).
				WithUint64Value(req.GetActionSetVrf().GetOid())).Build())
	}
	if req.ActionSetUserTrapId != nil {
		aReq.Actions = append(aReq.Actions,
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TRAP_ID).
				WithUint64Value(req.GetActionSetUserTrapId().GetOid())).Build())
	}
	if req.ActionPacketAction != nil {
		switch req.GetActionPacketAction().GetPacketAction() {
		case saipb.PacketAction_PACKET_ACTION_DROP,
			saipb.PacketAction_PACKET_ACTION_TRAP, // COPY and DROP
			saipb.PacketAction_PACKET_ACTION_DENY: // COPY_CANCEL and DROP
			aReq.Actions = append(aReq.Actions, &fwdpb.ActionDesc{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP})
		case saipb.PacketAction_PACKET_ACTION_FORWARD,
			saipb.PacketAction_PACKET_ACTION_LOG,     // COPY and FORWARD
			saipb.PacketAction_PACKET_ACTION_TRANSIT: // COPY_CANCEL and FORWARD
			aReq.Actions = append(aReq.Actions, &fwdpb.ActionDesc{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}) // Packets are forwarded by default so continue.
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unknown packet action type: %v", req.GetActionPacketAction().GetPacketAction())
		}
	}

	if _, err := a.dataplane.TableEntryAdd(ctx, aReq); err != nil {
		return nil, err
	}

	return &saipb.CreateAclEntryResponse{Oid: id}, nil
}
