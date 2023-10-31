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
	"math"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func newHostif(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *hostif {
	p := &hostif{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterHostifServer(s, p)
	return p
}

type hostif struct {
	saipb.UnimplementedHostifServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

// CreateHostif creates a hostif interface (usually a tap interface).
func (hostif *hostif) CreateHostif(ctx context.Context, req *saipb.CreateHostifRequest) (*saipb.CreateHostifResponse, error) {
	id := hostif.mgr.NextID()

	switch req.GetType() {
	case saipb.HostifType_HOSTIF_TYPE_GENETLINK: // TODO: support this type
		return &saipb.CreateHostifResponse{Oid: id}, nil
	case saipb.HostifType_HOSTIF_TYPE_NETDEV:
		portReq := &dpb.CreatePortRequest{
			Id:   fmt.Sprint(id),
			Type: fwdpb.PortType_PORT_TYPE_KERNEL,
			Src: &dpb.CreatePortRequest_KernelDev{
				KernelDev: string(req.GetName()),
			},
			Location: dpb.PortLocation_PORT_LOCATION_INTERNAL,
		}
		attrReq := &saipb.GetPortAttributeRequest{Oid: req.GetObjId(), AttrType: []saipb.PortAttr{saipb.PortAttr_PORT_ATTR_OPER_STATUS}}
		p := &saipb.GetPortAttributeResponse{}
		if err := hostif.mgr.PopulateAttributes(attrReq, p); err != nil {
			return nil, err
		}
		if p.GetAttr().GetOperStatus() != saipb.PortOperStatus_PORT_OPER_STATUS_NOT_PRESENT {
			portReq.ExternalPort = fmt.Sprint(req.GetObjId())
		}
		if _, err := hostif.dataplane.CreatePort(ctx, portReq); err != nil {
			if err != nil {
				return nil, err
			}
		}
		attr := &saipb.HostifAttribute{
			OperStatus: proto.Bool(true),
		}
		hostif.mgr.StoreAttributes(id, attr)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown type %v", req.GetType())
	}
	return &saipb.CreateHostifResponse{Oid: id}, nil
}

// SetHostifAttribute sets the attributes in the request.
func (hostif *hostif) SetHostifAttribute(ctx context.Context, req *saipb.SetHostifAttributeRequest) (*saipb.SetHostifAttributeResponse, error) {
	if req.OperStatus != nil {
		stateReq := &fwdpb.PortStateRequest{
			ContextId: &fwdpb.ContextId{Id: hostif.dataplane.ID()},
			PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(req.GetOid())}},
		}
		stateReq.Operation = &fwdpb.PortInfo{
			AdminStatus: fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
		}
		if req.GetOperStatus() {
			stateReq.Operation.AdminStatus = fwdpb.PortState_PORT_STATE_ENABLED_UP
		}
		_, err := hostif.dataplane.PortState(ctx, stateReq)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

var (
	etherTypeARP  = []byte{0x08, 0x06}
	bgpPort       = binary.BigEndian.AppendUint16(nil, 179)
	udldDstMAC    = []byte{0x01, 0x00, 0x0C, 0xCC, 0xCC, 0xCC}
	etherTypeLLDP = []byte{0x88, 0xcc}
	ndDstMAC      = []byte{0x33, 0x33, 0x00, 0x00, 0x00, 0x00} // ND is generic IPv6 multicast MAC.
	ndDstMACMask  = []byte{0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00}
	lacpDstMAC    = []byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x02}
)

func (hostif *hostif) CreateHostifTrap(_ context.Context, req *saipb.CreateHostifTrapRequest) (*saipb.CreateHostifTrapResponse, error) {
	id := hostif.mgr.NextID()

	fwdReq := fwdpb.TableEntryAddRequest{
		EntryDesc: &fwdpb.EntryDesc{
			Entry: &fwdpb.EntryDesc_Flow{
				Flow: &fwdpb.FlowEntryDesc{
					Priority: req.GetTrapPriority(),
					Id:       uint32(id),
				},
			},
		},
	}
	switch tType := req.GetTrapType(); tType {
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_ARP_REQUEST, saipb.HostifTrapType_HOSTIF_TRAP_TYPE_ARP_RESPONSE:
		fwdReq.EntryDesc.GetFlow().Fields = append(fwdReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE}},
			Bytes:   etherTypeARP,
			Masks:   []byte{0xFF, 0xFF},
		})
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_UDLD:
		fwdReq.EntryDesc.GetFlow().Fields = append(fwdReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST}},
			Bytes:   udldDstMAC,
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		})
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_LLDP:
		fwdReq.EntryDesc.GetFlow().Fields = append(fwdReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE}},
			Bytes:   etherTypeLLDP,
			Masks:   []byte{0xFF, 0xFF},
		})
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_DISCOVERY:
		fwdReq.EntryDesc.GetFlow().Fields = append(fwdReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST}},
			Bytes:   ndDstMAC,
			Masks:   ndDstMACMask,
		})
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_LACP:
		fwdReq.EntryDesc.GetFlow().Fields = append(fwdReq.EntryDesc.GetFlow().Fields, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST}},
			Bytes:   lacpDstMAC,
			Masks:   []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		})
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_IP2ME:
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_BGP, saipb.HostifTrapType_HOSTIF_TRAP_TYPE_BGPV6:
		fwdReq.EntryDesc = nil
		fwdReq.Entries = append(fwdReq.Entries, &fwdpb.TableEntryAddRequest_Entry{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Priority: req.GetTrapPriority(),
						Id:       uint32(id),
						Fields: []*fwdpb.PacketFieldMaskedBytes{{
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC}},
							Bytes:   bgpPort,
							Masks:   binary.BigEndian.AppendUint16(nil, math.MaxUint16),
						}},
					},
				},
			},
		}, &fwdpb.TableEntryAddRequest_Entry{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Priority: req.GetTrapPriority(),
						Id:       uint32(id),
						Fields: []*fwdpb.PacketFieldMaskedBytes{{
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST}},
							Bytes:   bgpPort,
							Masks:   binary.BigEndian.AppendUint16(nil, math.MaxUint16),
						}},
					},
				},
			},
		})
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown trap type: %v", tType)
	}
	switch act := req.GetPacketAction(); act {
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown action type: %v", act)
	}
	// group := req.GetTrapGroup()
	return &saipb.CreateHostifTrapResponse{
		Oid: id,
	}, nil
}

func (hostif *hostif) CreateHostifTrapGroup(context.Context, *saipb.CreateHostifTrapGroupRequest) (*saipb.CreateHostifTrapGroupResponse, error) {
	return nil, nil
}

func (hostif *hostif) CreateHostifUserDefinedTrap(context.Context, *saipb.CreateHostifUserDefinedTrapRequest) (*saipb.CreateHostifUserDefinedTrapResponse, error) {
	return nil, nil
}

func (hostif *hostif) CreateHostifTableEntry(context.Context, *saipb.CreateHostifTableEntryRequest) (*saipb.CreateHostifTableEntryResponse, error) {
	return nil, nil
}
