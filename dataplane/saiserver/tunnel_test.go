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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestCreateTunnel(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.CreateTunnelRequest
		wantReq *fwdpb.TableEntryAddRequest
		types   map[string]saipb.ObjectType
		wantErr string
	}{{
		desc:    "invalid request",
		req:     &saipb.CreateTunnelRequest{},
		wantErr: "InvalidArgument",
	}, {
		desc: "ipinip tunnel",
		req: &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_IPINIP.Enum(),
			EncapEcnMode:      saipb.TunnelEncapEcnMode_TUNNEL_ENCAP_ECN_MODE_STANDARD.Enum(),
			EncapDscpMode:     saipb.TunnelDscpMode_TUNNEL_DSCP_MODE_UNIFORM_MODEL.Enum(),
			EncapTtlMode:      saipb.TunnelTtlMode_TUNNEL_TTL_MODE_UNIFORM_MODEL.Enum(),
			UnderlayInterface: proto.Uint64(10),
		},
		wantReq: &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: TunnelEncap}},
			Entries: []*fwdpb.TableEntryAddRequest_Entry{{
				EntryDesc: &fwdpb.EntryDesc{
					Entry: &fwdpb.EntryDesc_Exact{
						Exact: &fwdpb.ExactEntryDesc{
							Fields: []*fwdpb.PacketFieldBytes{{
								FieldId: &fwdpb.PacketFieldId{
									Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID},
								},
								Bytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
							}},
						},
					},
				},
				Actions: []*fwdpb.ActionDesc{{
					ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
					Action: &fwdpb.ActionDesc_Update{
						Update: &fwdpb.UpdateActionDesc{
							Type:    fwdpb.UpdateType_UPDATE_TYPE_SET,
							Field:   &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{}},
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE}},
							Value:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a},
						},
					},
				}, {
					ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
					Action: &fwdpb.ActionDesc_Update{
						Update: &fwdpb.UpdateActionDesc{
							Type: fwdpb.UpdateType_UPDATE_TYPE_COPY,
							Field: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{
								FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS,
								Instance: 1,
							}},
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS}},
						},
					},
				}, {
					ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
					Action: &fwdpb.ActionDesc_Update{
						Update: &fwdpb.UpdateActionDesc{
							Type: fwdpb.UpdateType_UPDATE_TYPE_COPY,
							Field: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{
								FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP,
								Instance: 1,
							}},
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP}},
						},
					},
				}},
			}},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, mgr, stopFn := newTestTunnel(t, dplane)
			defer stopFn()
			for k, v := range tt.types {
				mgr.SetType(k, v)
			}
			mgr.StoreAttributes(1, &saipb.SwitchAttribute{
				CpuPort: proto.Uint64(10),
			})
			_, gotErr := c.CreateTunnel(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateTunnel() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotEntryAddReqs[0], tt.wantReq, protocmp.Transform()); d != "" {
				t.Errorf("CreateTunnel() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func newTestTunnel(t testing.TB, api switchDataplaneAPI) (saipb.TunnelClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newTunnel(mgr, api, srv)
	})
	return saipb.NewTunnelClient(conn), mgr, stopFn
}
