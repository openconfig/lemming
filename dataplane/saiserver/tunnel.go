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

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type tunnel struct {
	saipb.UnimplementedTunnelServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newTunnel(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *tunnel {
	t := &tunnel{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterTunnelServer(s, t)
	return t
}

func (t *tunnel) CreateTunnel(ctx context.Context, req *saipb.CreateTunnelRequest) (*saipb.CreateTunnelResponse, error) {
	id := t.mgr.NextID()

	tunType := req.GetType()
	switch tunType {
	case saipb.TunnelType_TUNNEL_TYPE_IPINIP:
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported tunnel type: %v", tunType)
	}

	actions := []*fwdpb.ActionDesc{}

	// TODO: Support parsing QOS into ECN and DSCP bits seperately.
	ecnMode := req.GetEncapEcnMode()
	dscpMode := req.GetEncapDscpMode()
	if ecnMode == saipb.TunnelEncapEcnMode_TUNNEL_ENCAP_ECN_MODE_STANDARD && dscpMode == saipb.TunnelDscpMode_TUNNEL_DSCP_MODE_UNIFORM_MODEL { // Copy the QOS bits from the inner IP header.
		actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_COPY, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS).
			WithFieldSrc(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS).WithFieldSrcInstance(1)).Build())
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported encap ecn mode %v and dscp mode: %v", ecnMode, dscpMode)
	}

	ttlMode := req.GetEncapTtlMode()
	switch ttlMode {
	case saipb.TunnelTtlMode_TUNNEL_TTL_MODE_UNIFORM_MODEL: // Copy the TTL from the inner IP header.
		actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_COPY, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP).
			WithFieldSrc(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP).WithFieldSrcInstance(1)).Build())
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported ttl mode: %v", ttlMode)
	}

	if req.GetEncapSrcIp() != nil {
		actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(
			fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithValue(req.EncapSrcIp)).Build())
	}
	if req.GetEncapDstIp() != nil {
		actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(
			fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithValue(req.EncapDstIp)).Build())
	}

	under := req.GetUnderlayInterface()
	entry := fwdconfig.TableEntryAddRequest(t.dataplane.ID(), TunnelEncap).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(id))),
		fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE).WithUint64Value(under)),
	).Build()
	entry.Entries[0].Actions = append(entry.Entries[0].Actions, actions...)
	if _, err := t.dataplane.TableEntryAdd(ctx, entry); err != nil {
		return nil, err
	}

	return &saipb.CreateTunnelResponse{
		Oid: id,
	}, nil
}

var (
	ipV4ExactMask = []byte{0xFF, 0xFF, 0xFF, 0xFF}
	ipV4AnyMask   = make([]byte, 4)
	ipV6ExactMask = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	ipV6AnyMask   = make([]byte, 16)
)

func (t *tunnel) CreateTunnelTermTableEntry(ctx context.Context, req *saipb.CreateTunnelTermTableEntryRequest) (*saipb.CreateTunnelTermTableEntryResponse, error) {
	id := t.mgr.NextID()

	fields := []*fwdpb.PacketFieldMaskedBytes{}

	// It is valid for some request fields to be omited, so initialize slices to the correct length for the given IP protocol.
	isV4 := len(req.SrcIp) == 4 || len(req.DstIp) == 4
	zeroIP := ipV6AnyMask
	exactMask := ipV6ExactMask
	headerID := fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6
	if isV4 {
		zeroIP = ipV4AnyMask
		exactMask = ipV4ExactMask
		headerID = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4
	}

	srcIP := req.GetSrcIp()
	if req.SrcIp == nil {
		srcIP = zeroIP
	}

	dstIP := req.GetDstIp()
	if req.DstIp == nil {
		dstIP = zeroIP
	}

	srcIPMask := req.GetSrcIpMask()
	if req.SrcIpMask == nil {
		srcIPMask = zeroIP
	}

	dstIPMask := req.GetDstIpMask()
	if req.DstIpMask == nil {
		dstIPMask = zeroIP
	}

	switch req.GetType() {
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P: // src IP, dst IP
		fields = append(fields,
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(srcIP, exactMask).Build(),
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(dstIP, exactMask).Build(),
		)
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP: // src IP, dst IP & mask
		fields = append(fields,
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(srcIP, exactMask).Build(),
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(dstIP, dstIPMask).Build(),
		)
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2P: // src IP & mask, dst IP
		fields = append(fields,
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(srcIP, srcIPMask).Build(),
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(dstIP, exactMask).Build(),
		)
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2MP: // src IP & mask, dst IP &mask
		fields = append(fields,
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(srcIP, srcIPMask).Build(),
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(srcIP, dstIPMask).Build(),
		)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid tunnel type: %v", req.GetType())
	}

	var actions []*fwdpb.ActionDesc
	switch req.GetTunnelType() { // TODO: The tunnel type should be used to match packets as well.
	case saipb.TunnelType_TUNNEL_TYPE_IPINIP:
		actions = append(actions, &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
			Action: &fwdpb.ActionDesc_Decap{
				Decap: &fwdpb.DecapActionDesc{
					HeaderId: headerID,
				},
			},
		})
	case saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE:
		actions = append(actions, &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
			Action: &fwdpb.ActionDesc_Decap{
				Decap: &fwdpb.DecapActionDesc{
					HeaderId: headerID,
				},
			},
		}, &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
			Action: &fwdpb.ActionDesc_Decap{
				Decap: &fwdpb.DecapActionDesc{
					HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
				},
			},
		})
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid tunnel type: %v", req.GetType())
	}
	actions = append(actions,
		fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64Value(req.GetVrId())).Build(),
	)

	tReq := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: t.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: turnTermTable}},
		EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
			Flow: &fwdpb.FlowEntryDesc{
				Id:       uint32(id),
				Priority: 1,
				Bank:     0,
				Fields:   fields,
			},
		}},
		Actions: actions,
	}
	if _, err := t.dataplane.TableEntryAdd(ctx, tReq); err != nil {
		return nil, err
	}

	return &saipb.CreateTunnelTermTableEntryResponse{
		Oid: id,
	}, nil
}
