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

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type policer struct {
	saipb.UnimplementedPolicerServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newPolicer(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *policer {
	p := &policer{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterPolicerServer(s, p)
	return p
}

// CreatePolicer creates a new policer, QOS is not actually supported. the GREEN action is always taken.
func (p *policer) CreatePolicer(ctx context.Context, req *saipb.CreatePolicerRequest) (*saipb.CreatePolicerResponse, error) {
	id := p.mgr.NextID()

	cpuPortReq := &saipb.GetSwitchAttributeRequest{Oid: switchID, AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT}}
	resp := &saipb.GetSwitchAttributeResponse{}
	if err := p.mgr.PopulateAttributes(cpuPortReq, resp); err != nil {
		return nil, err
	}

	var action *fwdpb.ActionDesc

	switch req.GetGreenPacketAction() {
	case saipb.PacketAction_PACKET_ACTION_TRAP:
		action = fwdconfig.Action(fwdconfig.TransmitAction(fmt.Sprint(resp.GetAttr().GetCpuPort())).WithImmediate(true)).Build()
	case saipb.PacketAction_PACKET_ACTION_COPY:
		action = &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_MIRROR,
			Action: &fwdpb.ActionDesc_Mirror{Mirror: &fwdpb.MirrorActionDesc{
				PortId:     &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(resp.GetAttr().GetCpuPort())}},
				PortAction: fwdpb.PortAction_PORT_ACTION_OUTPUT,
			}},
		}
	case saipb.PacketAction_PACKET_ACTION_FORWARD, saipb.PacketAction_PACKET_ACTION_UNSPECIFIED: // If unset, the default action is FORWARD.
		action = fwdconfig.Action(fwdconfig.ContinueAction()).Build()
	default:
		return nil, fmt.Errorf("unsupport policer action: %v", req.GetGreenPacketAction())
	}

	tReq := fwdconfig.TableEntryAddRequest(p.dataplane.ID(), policerTabler).
		AppendEntry(fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_POLICER_ID).WithUint64(id)))).Build()

	tReq.Entries[0].Actions = append(tReq.Entries[0].Actions, action)

	if _, err := p.dataplane.TableEntryAdd(ctx, tReq); err != nil {
		return nil, err
	}

	return &saipb.CreatePolicerResponse{
		Oid: id,
	}, nil
}

// RemovePolicer removes the entry from the table.
func (p *policer) RemovePolicer(ctx context.Context, req *saipb.RemovePolicerRequest) (*saipb.RemovePolicerResponse, error) {
	tReq := fwdconfig.TableEntryRemoveRequest(p.dataplane.ID(), policerTabler).
		AppendEntry(fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_POLICER_ID).WithUint64(req.GetOid()))))

	if _, err := p.dataplane.TableEntryRemove(ctx, tReq.Build()); err != nil {
		return nil, err
	}

	return &saipb.RemovePolicerResponse{}, nil
}

func (p *policer) SetPolicerAttribute(context.Context, *saipb.SetPolicerAttributeRequest) (*saipb.SetPolicerAttributeResponse, error) {
	return &saipb.SetPolicerAttributeResponse{}, nil
}
