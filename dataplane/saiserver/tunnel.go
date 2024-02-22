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
