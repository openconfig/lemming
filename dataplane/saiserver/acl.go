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
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
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

func (a *acl) createAclEntryFields(req *saipb.CreateAclEntryRequest, id uint64, gid string, bank int) (*fwdpb.TableEntryAddRequest, error) {
	aReq := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: a.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: gid}},
		EntryDesc: &fwdpb.EntryDesc{
			Entry: &fwdpb.EntryDesc_Flow{
				Flow: &fwdpb.FlowEntryDesc{
					Priority: math.MaxUint32 - req.GetPriority(), // TODO: SAI and Lucius have reversed definition of priority, this be cleaner if lucius supported both.
					Id:       uint32(id),
					Bank:     uint32(bank),
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
		fwdCtx, err := a.dataplane.FindContext(&fwdpb.ContextId{Id: a.dataplane.ID()})
		if err != nil {
			return nil, err
		}
		obj, err := fwdCtx.Objects.FindID(&fwdpb.ObjectId{Id: fmt.Sprint(req.FieldInPort.GetDataOid())})
		if err != nil {
			return nil, err
		}
		nid := obj.NID()
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT}},
			Bytes:   binary.BigEndian.AppendUint64(nil, uint64(nid)),
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
		case saipb.AclIpType_ACL_IP_TYPE_IP:
			fieldMask = &fwdpb.PacketFieldMaskedBytes{
				FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION}},
				Bytes:   []byte{0x4}, // IPv4 0100 IPv6 0110
				Masks:   []byte{0x4}, // Mask 0100
			}
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unsupported ACL_IP_TYPE: %v", t)
		}
		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, fieldMask)
	}
	if req.GetFieldDscp() != nil {
		// The QOS header in lucius corresponds to DSCP and ECN, so shift the bits left by 2.
		data := byte(req.GetFieldDscp().GetDataUint()) << 2
		mask := byte(req.GetFieldDscp().GetMaskUint()) << 2

		aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS}},
			Bytes:   []byte{data},
			Masks:   []byte{mask},
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
	if len(req.GetUserDefinedFieldGroupMin()) > 0 {
		table := &saipb.AclTableAttribute{}
		if err := a.mgr.PopulateAllAttributes(fmt.Sprint(req.GetTableId()), table); err != nil {
			return nil, err
		}
		for id, field := range req.GetUserDefinedFieldGroupMin() {
			udfGroup := &saipb.UdfGroupAttribute{}
			if err := a.mgr.PopulateAllAttributes(fmt.Sprint(table.UserDefinedFieldGroupMin[id]), udfGroup); err != nil {
				return nil, err
			}
			for _, udfID := range udfGroup.GetUdfList() {
				udf := &saipb.UdfAttribute{}
				if err := a.mgr.PopulateAllAttributes(fmt.Sprint(udfID), udf); err != nil {
					return nil, err
				}
				hg := fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L2
				switch udf.GetBase() {
				case saipb.UdfBase_UDF_BASE_L2:
					hg = fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L2
				case saipb.UdfBase_UDF_BASE_L3:
					hg = fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L3
				case saipb.UdfBase_UDF_BASE_L4:
					hg = fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L4
				}

				aReq.EntryDesc.GetFlow().Fields = append(aReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
					FieldId: &fwdpb.PacketFieldId{Bytes: &fwdpb.PacketBytes{
						HeaderGroup: hg,
						Instance:    0,
						Offset:      udf.GetOffset(),
						Size:        udfGroup.GetLength(),
					}},
					Bytes: field.GetDataU8List(),
					Masks: field.GetMaskU8List(),
				})
			}
		}
	}
	if len(aReq.EntryDesc.GetFlow().Fields) == 0 {
		return nil, status.Error(codes.InvalidArgument, "either no fields or not unsupports fields in entry req")
	}
	return aReq, nil
}

// CreateAclEntry adds an entry in the a bank.
func (a *acl) CreateAclEntry(ctx context.Context, req *saipb.CreateAclEntryRequest) (*saipb.CreateAclEntryResponse, error) {
	id := a.mgr.NextID()
	gb, ok := a.tableToLocation[req.GetTableId()]
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition, "table is not member of a group")
	}
	aReq, err := a.createAclEntryFields(req, id, gb.groupID, gb.bank)
	if err != nil {
		return nil, err
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
	if req.ActionCounter != nil {
		aReq.Actions = append(aReq.Actions, &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_FLOW_COUNTER,
			Action: &fwdpb.ActionDesc_Flow{
				Flow: &fwdpb.FlowCounterActionDesc{
					CounterId: &fwdpb.FlowCounterId{
						ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(req.GetActionCounter().GetOid())},
					},
				},
			},
		})
	}
	if req.ActionRedirect != nil {
		switch typ := a.mgr.GetType(fmt.Sprint(req.GetActionRedirect().GetOid())); typ {
		case saipb.ObjectType_OBJECT_TYPE_L2MC_GROUP:
			aReq.Actions = append(aReq.Actions, []*fwdpb.ActionDesc{
				fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET,
					fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID).WithUint64Value(req.GetActionRedirect().GetOid())).Build(),
				fwdconfig.Action(fwdconfig.LookupAction(L2MCGroupTable)).Build(), // Check L2MC group.
			}...)
		default:
			return nil, status.Errorf(codes.InvalidArgument, "type %q is not supported; only support L2MC Group for ACL Redirect for now", typ.String())
		}
	}
	if req.ActionSetPolicer != nil {
		aReq.Actions = append(aReq.Actions,
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_POLICER_ID).
				WithUint64Value(req.GetActionSetPolicer().GetOid())).Build(),
			fwdconfig.Action(fwdconfig.LookupAction(policerTabler)).Build(),
		)
	}

	cpuPortReq := &saipb.GetSwitchAttributeRequest{Oid: switchID, AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT}}
	resp := &saipb.GetSwitchAttributeResponse{}
	if err := a.mgr.PopulateAttributes(cpuPortReq, resp); err != nil {
		return nil, err
	}

	if req.ActionPacketAction != nil {
		switch req.GetActionPacketAction().GetPacketAction() {
		case saipb.PacketAction_PACKET_ACTION_DROP, saipb.PacketAction_PACKET_ACTION_DENY: // COPY_CANCEL and DROP
			aReq.Actions = append(aReq.Actions, &fwdpb.ActionDesc{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP})
		case saipb.PacketAction_PACKET_ACTION_TRAP: // COPY and DROP
			aReq.Actions = append(aReq.Actions, fwdconfig.Action(fwdconfig.TransmitAction(fmt.Sprint(resp.GetAttr().GetCpuPort())).WithImmediate(true)).Build())
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

func (a *acl) RemoveAclEntry(ctx context.Context, req *saipb.RemoveAclEntryRequest) (*saipb.RemoveAclEntryResponse, error) {
	cReq := &saipb.CreateAclEntryRequest{}
	if err := a.mgr.PopulateAllAttributes(fmt.Sprint(req.GetOid()), cReq); err != nil {
		return nil, err
	}
	gb, ok := a.tableToLocation[cReq.GetTableId()]
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition, "table is not member of a group")
	}

	aReq, err := a.createAclEntryFields(cReq, req.GetOid(), gb.groupID, gb.bank)
	if err != nil {
		return nil, err
	}

	if _, err := a.dataplane.TableEntryRemove(ctx, &fwdpb.TableEntryRemoveRequest{
		ContextId: aReq.ContextId,
		TableId:   aReq.TableId,
		EntryDesc: aReq.EntryDesc,
	}); err != nil {
		return nil, err
	}

	return &saipb.RemoveAclEntryResponse{}, nil
}

func (a *acl) CreateAclCounter(ctx context.Context, req *saipb.CreateAclCounterRequest) (*saipb.CreateAclCounterResponse, error) {
	id := a.mgr.NextID()

	_, err := a.dataplane.FlowCounterCreate(ctx, &fwdpb.FlowCounterCreateRequest{
		ContextId: &fwdpb.ContextId{Id: a.dataplane.ID()},
		Id:        &fwdpb.FlowCounterId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(id)}},
	})
	if err != nil {
		return nil, err
	}
	return &saipb.CreateAclCounterResponse{Oid: id}, nil
}

func (a *acl) RemoveAclCounter(ctx context.Context, req *saipb.RemoveAclCounterRequest) (*saipb.RemoveAclCounterResponse, error) {
	_, err := a.dataplane.ObjectDelete(ctx, &fwdpb.ObjectDeleteRequest{
		ContextId: &fwdpb.ContextId{Id: a.dataplane.ID()},
		ObjectId:  &fwdpb.ObjectId{Id: fmt.Sprint(req.GetOid())},
	})
	if err != nil {
		return nil, err
	}
	return &saipb.RemoveAclCounterResponse{}, nil
}

func (a *acl) GetAclCounterAttribute(ctx context.Context, req *saipb.GetAclCounterAttributeRequest) (*saipb.GetAclCounterAttributeResponse, error) {
	fetchStats := false
	for _, attr := range req.GetAttrType() {
		switch attr {
		case saipb.AclCounterAttr_ACL_COUNTER_ATTR_PACKETS, saipb.AclCounterAttr_ACL_COUNTER_ATTR_BYTES:
			fetchStats = true
		}
	}
	if !fetchStats {
		return &saipb.GetAclCounterAttributeResponse{}, nil
	}
	count, err := a.dataplane.FlowCounterQuery(ctx, &fwdpb.FlowCounterQueryRequest{
		ContextId: &fwdpb.ContextId{Id: a.dataplane.ID()},
		Ids:       []*fwdpb.FlowCounterId{{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(req.Oid)}}},
	})
	if err != nil {
		return nil, err
	}
	resp := &saipb.GetAclCounterAttributeResponse{
		Attr: &saipb.AclCounterAttribute{
			Packets: &count.GetCounters()[0].Packets,
			Bytes:   &count.GetCounters()[0].Octets,
		},
	}
	a.mgr.StoreAttributes(req.GetOid(), resp.GetAttr())
	return resp, nil
}

// entryDescFromReq returns the EntryDesc based on req.
func entryDescFromReq(m *myMac, req *saipb.CreateMyMacRequest) (*fwdpb.EntryDesc, error) {
	if req.Priority == nil {
		return nil, fmt.Errorf("priority needs to be specified")
	}
	fields := []*fwdconfig.PacketFieldMaskedBytesBuilder{
		fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).
			WithBytes(req.GetMacAddress(), req.GetMacAddressMask()),
	}

	if req.VlanId != nil {
		fields = append(fields,
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG).
				WithBytes(binary.BigEndian.AppendUint16(nil, uint16(req.GetVlanId())), binary.BigEndian.AppendUint16(nil, 0x0FFF)))
	}

	if req.PortId != nil {
		fwdCtx, err := m.dataplane.FindContext(&fwdpb.ContextId{Id: m.dataplane.ID()})
		if err != nil {
			return nil, err
		}
		obj, err := fwdCtx.Objects.FindID(&fwdpb.ObjectId{Id: fmt.Sprint(req.GetPortId())})
		if err != nil {
			return nil, err
		}
		fields = append(fields,
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT).
				WithBytes(binary.BigEndian.AppendUint64(nil, uint64(obj.NID())), binary.BigEndian.AppendUint64(nil, math.MaxUint64)))
	}

	ed := &fwdpb.EntryDesc{
		Entry: &fwdpb.EntryDesc_Flow{
			Flow: &fwdpb.FlowEntryDesc{
				Priority: req.GetPriority(),
				Bank:     1,
			},
		},
	}
	for _, f := range fields {
		ed.GetFlow().Fields = append(ed.GetFlow().Fields, f.Build())
	}
	return ed, nil
}

type myMac struct {
	saipb.UnimplementedMyMacServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newMyMac(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *myMac {
	m := &myMac{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterMyMacServer(s, m)
	return m
}

func (m *myMac) CreateMyMac(ctx context.Context, req *saipb.CreateMyMacRequest) (*saipb.CreateMyMacResponse, error) {
	if req.MacAddress == nil || req.MacAddressMask == nil {
		return nil, status.Errorf(codes.InvalidArgument, "MAC address and MAC address mask cannot be empty")
	}
	id := m.mgr.NextID()
	ed, err := entryDescFromReq(m, req)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "failed to create entry descriptor: %v", err)
	}
	mReq := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: m.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: MyMacTable}},
		EntryDesc: ed,
		Actions:   getL3Pipeline(),
	}
	if _, err := m.dataplane.TableEntryAdd(ctx, mReq); err != nil {
		return nil, err
	}
	// Populate the switch attribute.
	saReq := &saipb.GetSwitchAttributeRequest{
		Oid: switchID,
		AttrType: []saipb.SwitchAttr{
			saipb.SwitchAttr_SWITCH_ATTR_MY_MAC_LIST,
		},
	}
	saResp := &saipb.GetSwitchAttributeResponse{}
	if err := m.mgr.PopulateAttributes(saReq, saResp); err != nil {
		return nil, fmt.Errorf("Failed to populate switch attributes: %v", err)
	}
	mml := append(saResp.GetAttr().MyMacList, id)
	m.mgr.StoreAttributes(switchID, &saipb.SwitchAttribute{
		MyMacList: mml,
	})
	return &saipb.CreateMyMacResponse{Oid: id}, nil
}

func (m *myMac) RemoveMyMac(ctx context.Context, req *saipb.RemoveMyMacRequest) (*saipb.RemoveMyMacResponse, error) {
	cReq := &saipb.CreateMyMacRequest{}
	if err := m.mgr.PopulateAllAttributes(fmt.Sprint(req.GetOid()), cReq); err != nil {
		return nil, err
	}
	ed, err := entryDescFromReq(m, cReq)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "failed to create entry descriptor: %v", err)
	}

	mReq := &fwdpb.TableEntryRemoveRequest{
		ContextId: &fwdpb.ContextId{Id: m.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: MyMacTable}},
		EntryDesc: ed,
	}
	if _, err := m.dataplane.TableEntryRemove(ctx, mReq); err != nil {
		return nil, err
	}

	// Populate the switch attribute.
	saReq := &saipb.GetSwitchAttributeRequest{
		Oid: switchID,
		AttrType: []saipb.SwitchAttr{
			saipb.SwitchAttr_SWITCH_ATTR_MY_MAC_LIST,
		},
	}
	saResp := &saipb.GetSwitchAttributeResponse{}
	if err := m.mgr.PopulateAttributes(saReq, saResp); err != nil {
		return nil, fmt.Errorf("Failed to populate switch attributes: %v", err)
	}
	locate := func(uint64) int {
		for i := range saResp.GetAttr().MyMacList {
			if saResp.GetAttr().MyMacList[i] == req.GetOid() {
				return i
			}
		}
		return -1
	}
	if loc := locate(req.GetOid()); loc != -1 {
		m.mgr.StoreAttributes(switchID, &saipb.SwitchAttribute{
			MyMacList: append(saResp.GetAttr().MyMacList[:loc], saResp.GetAttr().MyMacList[loc+1:]...),
		})
	} else {
		return nil, fmt.Errorf("Failed to store switch attributes: %v", err)
	}

	return &saipb.RemoveMyMacResponse{}, nil
}
