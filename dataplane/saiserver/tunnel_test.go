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
	"math"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
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
	}, {
		desc: "vxlan tunnel",
		req: &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
			UnderlayInterface: proto.Uint64(10),
			VxlanUdpSport:     proto.Uint32(1234),
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
								FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP,
								Instance: 1,
							}},
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP}},
						},
					},
				}, {
					ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
					Action: &fwdpb.ActionDesc_Update{
						Update: &fwdpb.UpdateActionDesc{
							Type:    fwdpb.UpdateType_UPDATE_TYPE_SET,
							Field:   &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{}},
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST}},
							Value:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0xb5},
						},
					},
				}, {
					ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
					Action: &fwdpb.ActionDesc_Update{
						Update: &fwdpb.UpdateActionDesc{
							Type:    fwdpb.UpdateType_UPDATE_TYPE_SET,
							Field:   &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{}},
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC}},
							Value:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xd2},
						},
					},
				}},
			}},
		},
	}, {
		desc: "vxlan tunnel pipe ttl",
		req: &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
			EncapTtlMode:      saipb.TunnelTtlMode_TUNNEL_TTL_MODE_PIPE_MODEL.Enum(),
			EncapTtlVal:       proto.Uint32(32),
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
							Type:    fwdpb.UpdateType_UPDATE_TYPE_SET,
							Field:   &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{}},
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP}},
							Value:   []byte{32},
						},
					},
				}, {
					ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
					Action: &fwdpb.ActionDesc_Update{
						Update: &fwdpb.UpdateActionDesc{
							Type:    fwdpb.UpdateType_UPDATE_TYPE_SET,
							Field:   &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{}},
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST}},
							Value:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0xb5},
						},
					},
				}},
			}},
		},
	}, {
		desc: "vxlan tunnel pipe default ttl",
		req: &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
			EncapTtlMode:      saipb.TunnelTtlMode_TUNNEL_TTL_MODE_PIPE_MODEL.Enum(),
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
							Type:    fwdpb.UpdateType_UPDATE_TYPE_SET,
							Field:   &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{}},
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP}},
							Value:   []byte{defaultEncapTTL},
						},
					},
				}, {
					ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
					Action: &fwdpb.ActionDesc_Update{
						Update: &fwdpb.UpdateActionDesc{
							Type:    fwdpb.UpdateType_UPDATE_TYPE_SET,
							Field:   &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{}},
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST}},
							Value:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0xb5},
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

func TestRemoveTunnel(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.RemoveTunnelRequest
		wantErr string
	}{{
		desc:    "invalid request",
		req:     &saipb.RemoveTunnelRequest{},
		wantErr: "not found",
	}, {
		desc: "success",
		req: &saipb.RemoveTunnelRequest{
			Oid: 1,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, _, stopFn := newTestTunnel(t, dplane)
			defer stopFn()
			_, err := c.CreateTunnel(context.TODO(), &saipb.CreateTunnelRequest{
				Type:              saipb.TunnelType_TUNNEL_TYPE_IPINIP.Enum(),
				EncapEcnMode:      saipb.TunnelEncapEcnMode_TUNNEL_ENCAP_ECN_MODE_STANDARD.Enum(),
				EncapDscpMode:     saipb.TunnelDscpMode_TUNNEL_DSCP_MODE_UNIFORM_MODEL.Enum(),
				EncapTtlMode:      saipb.TunnelTtlMode_TUNNEL_TTL_MODE_UNIFORM_MODEL.Enum(),
				UnderlayInterface: proto.Uint64(10),
			})
			if err != nil {
				t.Fatalf("CreateTunnel() unexpected err: %s", err)
			}
			_, gotErr := c.RemoveTunnel(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("RemoveTunnel() unexpected err: %s", diff)
			}
		})
	}
}

func TestTunnelTermTableEntry(t *testing.T) {
	tests := []struct {
		desc           string
		req            *saipb.CreateTunnelTermTableEntryRequest
		wantAddReqs    []*fwdpb.TableEntryAddRequest
		wantRemoveReqs []*fwdpb.TableEntryRemoveRequest
		wantErr        string
	}{
		{
			desc: "vxlan p2mp ipv4 with action tunnel id",
			req: &saipb.CreateTunnelTermTableEntryRequest{
				TunnelType:     saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
				Type:           saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP.Enum(),
				DstIp:          []byte{10, 0, 0, 1},
				VrId:           proto.Uint64(10),
				ActionTunnelId: proto.Uint64(100),
			},
			wantAddReqs: []*fwdpb.TableEntryAddRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 2, // (MaxUint32 - 3) - 0 + termRankP2MP
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{17}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST).WithBytes([]byte{0x12, 0xB5}, []byte{0xFF, 0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{0, 0, 0, 0}, []byte{0, 0, 0, 0}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{10, 0, 0, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
						},
					},
				}},
				Actions: []*fwdpb.ActionDesc{
					fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64Value(100)).Build(),
					fwdconfig.Action(fwdconfig.LookupAction(vniToVrfTable)).Build(),
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET}},
					},
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4}},
					},
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP}},
					},
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN}},
					},
				},
			}},
			wantRemoveReqs: []*fwdpb.TableEntryRemoveRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 2,
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{17}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST).WithBytes([]byte{0x12, 0xB5}, []byte{0xFF, 0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{0, 0, 0, 0}, []byte{0, 0, 0, 0}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{10, 0, 0, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
						},
					},
				}},
			}},
		},
		{
			desc: "vxlan p2p ipv6 with action tunnel id",
			req: &saipb.CreateTunnelTermTableEntryRequest{
				TunnelType:     saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
				Type:           saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
				SrcIp:          net.ParseIP("2001:db8::1"),
				DstIp:          net.ParseIP("2001:db8::2"),
				VrId:           proto.Uint64(10),
				ActionTunnelId: proto.Uint64(200),
			},
			wantAddReqs: []*fwdpb.TableEntryAddRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 3, // (MaxUint32 - 3) - 0 + termRankP2P
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{6}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{17}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST).WithBytes([]byte{0x12, 0xB5}, []byte{0xFF, 0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(net.ParseIP("2001:db8::1"), ipV6ExactMask).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(net.ParseIP("2001:db8::2"), ipV6ExactMask).Build(),
						},
					},
				}},
				Actions: []*fwdpb.ActionDesc{
					fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64Value(200)).Build(),
					fwdconfig.Action(fwdconfig.LookupAction(vniToVrfTable)).Build(),
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET}},
					},
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6}},
					},
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP}},
					},
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN}},
					},
				},
			}},
			wantRemoveReqs: []*fwdpb.TableEntryRemoveRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 3,
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{6}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{17}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST).WithBytes([]byte{0x12, 0xB5}, []byte{0xFF, 0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(net.ParseIP("2001:db8::1"), ipV6ExactMask).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(net.ParseIP("2001:db8::2"), ipV6ExactMask).Build(),
						},
					},
				}},
			}},
		},
		{
			desc: "vxlan without action tunnel id",
			req: &saipb.CreateTunnelTermTableEntryRequest{
				TunnelType: saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
				Type:       saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP.Enum(),
				DstIp:      []byte{10, 0, 0, 1},
			},
			wantErr: "action_tunnel_id is required",
		},
		{
			desc: "gre mp2p ipv4",
			req: &saipb.CreateTunnelTermTableEntryRequest{
				TunnelType: saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE.Enum(),
				Type:       saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2P.Enum(),
				SrcIp:      []byte{10, 0, 0, 1},
				DstIp:      []byte{20, 0, 0, 0},
				DstIpMask:  []byte{255, 255, 255, 0},
				VrId:       proto.Uint64(5),
			},
			wantAddReqs: []*fwdpb.TableEntryAddRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 2, // (MaxUint32 - 3) - 0 + termRankMP2P
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{47}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{10, 0, 0, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{20, 0, 0, 0}, []byte{255, 255, 255, 0}).Build(),
						},
					},
				}},
				Actions: []*fwdpb.ActionDesc{
					fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64Value(5)).Build(),
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4}},
					},
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE}},
					},
				},
			}},
			wantRemoveReqs: []*fwdpb.TableEntryRemoveRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 2,
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{47}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{10, 0, 0, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{20, 0, 0, 0}, []byte{255, 255, 255, 0}).Build(),
						},
					},
				}},
			}},
		},
		{
			desc: "gre p2p ipv4 nil vrid",
			req: &saipb.CreateTunnelTermTableEntryRequest{
				TunnelType: saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE.Enum(),
				Type:       saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
				SrcIp:      []byte{10, 0, 0, 1},
				DstIp:      []byte{20, 0, 0, 2},
			},
			wantAddReqs: []*fwdpb.TableEntryAddRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 3, // (MaxUint32 - 3) - 0 + termRankP2P
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{47}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{10, 0, 0, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{20, 0, 0, 2}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
						},
					},
				}},
				Actions: []*fwdpb.ActionDesc{
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4}},
					},
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE}},
					},
				},
			}},
			wantRemoveReqs: []*fwdpb.TableEntryRemoveRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 3,
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{47}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{10, 0, 0, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{20, 0, 0, 2}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
						},
					},
				}},
			}},
		},
		{
			desc: "ipinip p2p ipv4",
			req: &saipb.CreateTunnelTermTableEntryRequest{
				TunnelType: saipb.TunnelType_TUNNEL_TYPE_IPINIP.Enum(),
				Type:       saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
				SrcIp:      []byte{192, 168, 1, 1},
				DstIp:      []byte{192, 168, 1, 2},
				VrId:       proto.Uint64(0),
			},
			wantAddReqs: []*fwdpb.TableEntryAddRequest{
				{
					ContextId: &fwdpb.ContextId{Id: "foo"},
					TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
					EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
						Flow: &fwdpb.FlowEntryDesc{
							Id:       1,
							Priority: math.MaxUint32 - 3,
							Bank:     0,
							Fields: []*fwdpb.PacketFieldMaskedBytes{
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{4}, []byte{0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{192, 168, 1, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{192, 168, 1, 2}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
							},
						},
					}},
					Actions: []*fwdpb.ActionDesc{
						fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64Value(0)).Build(),
						{
							ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
							Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4}},
						},
					},
				},
				{
					ContextId: &fwdpb.ContextId{Id: "foo"},
					TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
					EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
						Flow: &fwdpb.FlowEntryDesc{
							Id:       1,
							Priority: math.MaxUint32 - 3,
							Bank:     0,
							Fields: []*fwdpb.PacketFieldMaskedBytes{
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{41}, []byte{0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{192, 168, 1, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{192, 168, 1, 2}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
							},
						},
					}},
					Actions: []*fwdpb.ActionDesc{
						fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64Value(0)).Build(),
						{
							ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
							Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4}},
						},
					},
				},
			},
			wantRemoveReqs: []*fwdpb.TableEntryRemoveRequest{
				{
					ContextId: &fwdpb.ContextId{Id: "foo"},
					TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
					EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
						Flow: &fwdpb.FlowEntryDesc{
							Id:       1,
							Priority: math.MaxUint32 - 3,
							Bank:     0,
							Fields: []*fwdpb.PacketFieldMaskedBytes{
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{4}, []byte{0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{192, 168, 1, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{192, 168, 1, 2}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
							},
						},
					}},
				},
				{
					ContextId: &fwdpb.ContextId{Id: "foo"},
					TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
					EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
						Flow: &fwdpb.FlowEntryDesc{
							Id:       1,
							Priority: math.MaxUint32 - 3,
							Bank:     0,
							Fields: []*fwdpb.PacketFieldMaskedBytes{
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{41}, []byte{0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{192, 168, 1, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
								fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{192, 168, 1, 2}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
							},
						},
					}},
				},
			},
		},
		{
			desc: "ipv4 16-byte representation normalization",
			req: &saipb.CreateTunnelTermTableEntryRequest{
				TunnelType: saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE.Enum(),
				Type:       saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
				SrcIp:      net.ParseIP("10.0.0.1"), // 16-byte slice without .To4()
				DstIp:      net.ParseIP("10.0.0.2"), // 16-byte slice without .To4()
			},
			wantAddReqs: []*fwdpb.TableEntryAddRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 3,
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{47}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{10, 0, 0, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{10, 0, 0, 2}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
						},
					},
				}},
				Actions: []*fwdpb.ActionDesc{
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4}},
					},
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE}},
					},
				},
			}},
			wantRemoveReqs: []*fwdpb.TableEntryRemoveRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 3,
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{47}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes([]byte{10, 0, 0, 1}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes([]byte{10, 0, 0, 2}, []byte{0xFF, 0xFF, 0xFF, 0xFF}).Build(),
						},
					},
				}},
			}},
		},
		{
			desc: "ipv6 dst with sai default 4-byte zero mask",
			req: &saipb.CreateTunnelTermTableEntryRequest{
				TunnelType: saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE.Enum(),
				Type:       saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP.Enum(),
				DstIp:      net.ParseIP("2001:db8::2"),
				SrcIpMask:  []byte{0, 0, 0, 0},
			},
			wantAddReqs: []*fwdpb.TableEntryAddRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 2, // (MaxUint32 - 3) - 0 + termRankP2MP
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{6}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{47}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(make([]byte, 16), make([]byte, 16)).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(net.ParseIP("2001:db8::2"), ipV6ExactMask).Build(),
						},
					},
				}},
				Actions: []*fwdpb.ActionDesc{
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6}},
					},
					{
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action:     &fwdpb.ActionDesc_Decap{Decap: &fwdpb.DecapActionDesc{HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE}},
					},
				},
			}},
			wantRemoveReqs: []*fwdpb.TableEntryRemoveRequest{{
				ContextId: &fwdpb.ContextId{Id: "foo"},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
				EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id:       1,
						Priority: math.MaxUint32 - 2,
						Bank:     0,
						Fields: []*fwdpb.PacketFieldMaskedBytes{
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{6}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).WithBytes([]byte{47}, []byte{0xFF}).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(make([]byte, 16), make([]byte, 16)).Build(),
							fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(net.ParseIP("2001:db8::2"), ipV6ExactMask).Build(),
						},
					},
				}},
			}},
		},
		{
			desc: "mixed ipv4 and ipv6",
			req: &saipb.CreateTunnelTermTableEntryRequest{
				TunnelType: saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE.Enum(),
				Type:       saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
				SrcIp:      []byte{10, 0, 0, 1},
				DstIp:      net.ParseIP("2001:db8::2"),
			},
			wantErr: "mismatched IP address families",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, _, stopFn := newTestTunnel(t, dplane)
			defer stopFn()

			resp, err := c.CreateTunnelTermTableEntry(context.TODO(), tt.req)
			if diff := errdiff.Check(err, tt.wantErr); diff != "" {
				t.Fatalf("CreateTunnelTermTableEntry() unexpected err: %s", diff)
			}
			if err != nil {
				return
			}
			if len(dplane.gotEntryAddReqs) != len(tt.wantAddReqs) {
				t.Fatalf("CreateTunnelTermTableEntry() got %d add requests, want %d", len(dplane.gotEntryAddReqs), len(tt.wantAddReqs))
			}
			for i := range tt.wantAddReqs {
				if diff := cmp.Diff(dplane.gotEntryAddReqs[i], tt.wantAddReqs[i], protocmp.Transform()); diff != "" {
					t.Errorf("CreateTunnelTermTableEntry() addReq[%d] diff (-got +want):\n%s", i, diff)
				}
			}

			_, err = c.RemoveTunnelTermTableEntry(context.TODO(), &saipb.RemoveTunnelTermTableEntryRequest{
				Oid: resp.GetOid(),
			})
			if err != nil {
				t.Fatalf("RemoveTunnelTermTableEntry() unexpected err: %v", err)
			}
			if len(dplane.gotEntryRemoveReqs) != len(tt.wantRemoveReqs) {
				t.Fatalf("RemoveTunnelTermTableEntry() got %d remove requests, want %d", len(dplane.gotEntryRemoveReqs), len(tt.wantRemoveReqs))
			}
			for i := range tt.wantRemoveReqs {
				if diff := cmp.Diff(dplane.gotEntryRemoveReqs[i], tt.wantRemoveReqs[i], protocmp.Transform()); diff != "" {
					t.Errorf("RemoveTunnelTermTableEntry() removeReq[%d] diff (-got +want):\n%s", i, diff)
				}
			}
		})
	}
}

func TestTunnelTermDuplicateAndPriority(t *testing.T) {
	t.Run("duplicate entry rejection", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, _, stopFn := newTestTunnel(t, dplane)
		defer stopFn()

		req := &saipb.CreateTunnelTermTableEntryRequest{
			TunnelType: saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE.Enum(),
			Type:       saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
			SrcIp:      []byte{10, 0, 0, 1},
			DstIp:      []byte{20, 0, 0, 2},
		}
		_, err := c.CreateTunnelTermTableEntry(context.TODO(), req)
		if err != nil {
			t.Fatalf("CreateTunnelTermTableEntry() first call failed: %v", err)
		}
		_, err = c.CreateTunnelTermTableEntry(context.TODO(), req)
		if s, ok := status.FromError(err); !ok || s.Code() != codes.AlreadyExists {
			t.Fatalf("CreateTunnelTermTableEntry() second call expected AlreadyExists, got: %v", err)
		}
	})

	t.Run("term specificity and sai priority ordering", func(t *testing.T) {
		p2pReq := &saipb.CreateTunnelTermTableEntryRequest{
			Type: saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
		}
		p2mpReq := &saipb.CreateTunnelTermTableEntryRequest{
			Type: saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP.Enum(),
		}
		mp2mpReq := &saipb.CreateTunnelTermTableEntryRequest{
			Type: saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2MP.Enum(),
		}

		p2pPrio, _ := termEntryPriority(p2pReq)
		p2mpPrio, _ := termEntryPriority(p2mpReq)
		mp2mpPrio, _ := termEntryPriority(mp2mpReq)

		if p2pPrio >= p2mpPrio || p2mpPrio >= mp2mpPrio {
			t.Errorf("term priority rank order invalid: p2p=%d, p2mp=%d, mp2mp=%d", p2pPrio, p2mpPrio, mp2mpPrio)
		}

		highSaiPrioReq := &saipb.CreateTunnelTermTableEntryRequest{
			Type:     saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
			Priority: proto.Uint32(100),
		}
		lowSaiPrioReq := &saipb.CreateTunnelTermTableEntryRequest{
			Type:     saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
			Priority: proto.Uint32(10),
		}
		highPrio, _ := termEntryPriority(highSaiPrioReq)
		lowPrio, _ := termEntryPriority(lowSaiPrioReq)
		if highPrio >= lowPrio {
			t.Errorf("higher SAI priority must yield lower numeric Lucius priority: highSAI=%d, lowSAI=%d", highPrio, lowPrio)
		}
	})
}

func TestTunnelMapAndEntry(t *testing.T) {
	t.Run("vlan id to vni map and entry create and remove", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, _, stopFn := newTestTunnel(t, dplane)
		defer stopFn()

		mapResp, err := c.CreateTunnelMap(context.TODO(), &saipb.CreateTunnelMapRequest{
			Type: saipb.TunnelMapType_TUNNEL_MAP_TYPE_VLAN_ID_TO_VNI.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMap() unexpected err: %v", err)
		}
		if mapResp.GetOid() == 0 {
			t.Errorf("CreateTunnelMap() returned 0 OID")
		}

		entryResp, err := c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap: proto.Uint64(mapResp.GetOid()),
			VlanIdKey: proto.Uint32(100),
			VniIdKey:  proto.Uint32(1000),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMapEntry() unexpected err: %v", err)
		}
		if entryResp.GetOid() == 0 {
			t.Errorf("CreateTunnelMapEntry() returned 0 OID")
		}

		// Non-VNI_TO_VRF map entries shouldn't program vni-to-vrf table.
		if len(dplane.gotEntryAddReqs) != 0 {
			t.Errorf("CreateTunnelMapEntry() added %d entries, want 0", len(dplane.gotEntryAddReqs))
		}

		_, err = c.RemoveTunnelMapEntry(context.TODO(), &saipb.RemoveTunnelMapEntryRequest{
			Oid: entryResp.GetOid(),
		})
		if err != nil {
			t.Fatalf("RemoveTunnelMapEntry() unexpected err: %v", err)
		}

		_, err = c.RemoveTunnelMap(context.TODO(), &saipb.RemoveTunnelMapRequest{
			Oid: mapResp.GetOid(),
		})
		if err != nil {
			t.Fatalf("RemoveTunnelMap() unexpected err: %v", err)
		}
	})

	t.Run("vni to vrf: tunnel then map entry (orchagent order)", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, _, stopFn := newTestTunnel(t, dplane)
		defer stopFn()

		mapResp, err := c.CreateTunnelMap(context.TODO(), &saipb.CreateTunnelMapRequest{
			Type: saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMap() unexpected err: %v", err)
		}

		tunResp, err := c.CreateTunnel(context.TODO(), &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
			UnderlayInterface: proto.Uint64(10),
			DecapMappers:      []uint64{mapResp.GetOid()},
		})
		if err != nil {
			t.Fatalf("CreateTunnel() unexpected err: %v", err)
		}

		// Only TunnelEncap entry added so far.
		if len(dplane.gotEntryAddReqs) != 1 {
			t.Fatalf("CreateTunnel() got %d add requests, want 1", len(dplane.gotEntryAddReqs))
		}

		entryResp, err := c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap:            proto.Uint64(mapResp.GetOid()),
			VniIdKey:             proto.Uint32(5000),
			VirtualRouterIdValue: proto.Uint64(42),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMapEntry() unexpected err: %v", err)
		}

		if len(dplane.gotEntryAddReqs) != 2 {
			t.Fatalf("CreateTunnelMapEntry() got %d add requests, want 2", len(dplane.gotEntryAddReqs))
		}

		wantVniEntry := fwdconfig.TableEntryAddRequest("foo", vniToVrfTable).AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ExactEntry(
				fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(tunResp.GetOid()),
				fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI).WithBytes(vniToBytes(5000)),
			)),
			fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64Value(42),
		).Build()

		if diff := cmp.Diff(dplane.gotEntryAddReqs[1], wantVniEntry, protocmp.Transform()); diff != "" {
			t.Errorf("vniToVrfTable entry diff (-got +want):\n%s", diff)
		}

		// Remove entry: verifies table entry removal.
		_, err = c.RemoveTunnelMapEntry(context.TODO(), &saipb.RemoveTunnelMapEntryRequest{
			Oid: entryResp.GetOid(),
		})
		if err != nil {
			t.Fatalf("RemoveTunnelMapEntry() unexpected err: %v", err)
		}

		wantVniRemove := &fwdpb.TableEntryRemoveRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: vniToVrfTable}},
			EntryDesc: fwdconfig.EntryDesc(fwdconfig.ExactEntry(
				fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(tunResp.GetOid()),
				fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI).WithBytes(vniToBytes(5000)),
			)).Build(),
		}
		if len(dplane.gotEntryRemoveReqs) != 1 {
			t.Fatalf("RemoveTunnelMapEntry() got %d remove requests, want 1", len(dplane.gotEntryRemoveReqs))
		}
		if diff := cmp.Diff(dplane.gotEntryRemoveReqs[0], wantVniRemove, protocmp.Transform()); diff != "" {
			t.Errorf("vniToVrfTable remove diff (-got +want):\n%s", diff)
		}

		_, err = c.RemoveTunnel(context.TODO(), &saipb.RemoveTunnelRequest{
			Oid: tunResp.GetOid(),
		})
		if err != nil {
			t.Fatalf("RemoveTunnel() unexpected err: %v", err)
		}

		_, err = c.RemoveTunnelMap(context.TODO(), &saipb.RemoveTunnelMapRequest{
			Oid: mapResp.GetOid(),
		})
		if err != nil {
			t.Fatalf("RemoveTunnelMap() unexpected err: %v", err)
		}
	})

	t.Run("vni to vrf: map entry then tunnel", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, _, stopFn := newTestTunnel(t, dplane)
		defer stopFn()

		mapResp, err := c.CreateTunnelMap(context.TODO(), &saipb.CreateTunnelMapRequest{
			Type: saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMap() unexpected err: %v", err)
		}

		entryResp, err := c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap:            proto.Uint64(mapResp.GetOid()),
			VniIdKey:             proto.Uint32(6000),
			VirtualRouterIdValue: proto.Uint64(55),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMapEntry() unexpected err: %v", err)
		}

		// No tunnel attached yet, so no table entries added.
		if len(dplane.gotEntryAddReqs) != 0 {
			t.Fatalf("CreateTunnelMapEntry() got %d add requests, want 0", len(dplane.gotEntryAddReqs))
		}

		tunResp, err := c.CreateTunnel(context.TODO(), &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
			UnderlayInterface: proto.Uint64(10),
			DecapMappers:      []uint64{mapResp.GetOid()},
		})
		if err != nil {
			t.Fatalf("CreateTunnel() unexpected err: %v", err)
		}

		// TableEntryAdd #0 is TunnelEncap, #1 is vniToVrfTable.
		if len(dplane.gotEntryAddReqs) != 2 {
			t.Fatalf("CreateTunnel() got %d add requests, want 2", len(dplane.gotEntryAddReqs))
		}

		wantVniEntry := fwdconfig.TableEntryAddRequest("foo", vniToVrfTable).AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ExactEntry(
				fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(tunResp.GetOid()),
				fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI).WithBytes(vniToBytes(6000)),
			)),
			fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64Value(55),
		).Build()

		if diff := cmp.Diff(dplane.gotEntryAddReqs[1], wantVniEntry, protocmp.Transform()); diff != "" {
			t.Errorf("vniToVrfTable entry diff (-got +want):\n%s", diff)
		}

		// RemoveTunnel first: should remove both TunnelEncap (#0) and vniToVrfTable (#1).
		_, err = c.RemoveTunnel(context.TODO(), &saipb.RemoveTunnelRequest{
			Oid: tunResp.GetOid(),
		})
		if err != nil {
			t.Fatalf("RemoveTunnel() unexpected err: %v", err)
		}

		if len(dplane.gotEntryRemoveReqs) != 2 {
			t.Fatalf("RemoveTunnel() got %d remove requests, want 2", len(dplane.gotEntryRemoveReqs))
		}

		wantVniRemove := &fwdpb.TableEntryRemoveRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: vniToVrfTable}},
			EntryDesc: fwdconfig.EntryDesc(fwdconfig.ExactEntry(
				fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(tunResp.GetOid()),
				fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI).WithBytes(vniToBytes(6000)),
			)).Build(),
		}
		if diff := cmp.Diff(dplane.gotEntryRemoveReqs[1], wantVniRemove, protocmp.Transform()); diff != "" {
			t.Errorf("vniToVrfTable remove diff (-got +want):\n%s", diff)
		}

		_, err = c.RemoveTunnelMapEntry(context.TODO(), &saipb.RemoveTunnelMapEntryRequest{
			Oid: entryResp.GetOid(),
		})
		if err != nil {
			t.Fatalf("RemoveTunnelMapEntry() unexpected err: %v", err)
		}

		_, err = c.RemoveTunnelMap(context.TODO(), &saipb.RemoveTunnelMapRequest{
			Oid: mapResp.GetOid(),
		})
		if err != nil {
			t.Fatalf("RemoveTunnelMap() unexpected err: %v", err)
		}
	})

	t.Run("map entry validation and error codes", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, _, stopFn := newTestTunnel(t, dplane)
		defer stopFn()

		// Unknown tunnel map -> NotFound.
		_, err := c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap: proto.Uint64(9999),
			VniIdKey:  proto.Uint32(100),
		})
		if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound {
			t.Errorf("CreateTunnelMapEntry with unknown map expected NotFound, got %v", err)
		}

		mapResp, err := c.CreateTunnelMap(context.TODO(), &saipb.CreateTunnelMapRequest{
			Type: saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMap() failed: %v", err)
		}

		// Entry type mismatch -> InvalidArgument.
		_, err = c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap:     proto.Uint64(mapResp.GetOid()),
			TunnelMapType: saipb.TunnelMapType_TUNNEL_MAP_TYPE_VLAN_ID_TO_VNI.Enum(),
		})
		if s, ok := status.FromError(err); !ok || s.Code() != codes.InvalidArgument {
			t.Errorf("CreateTunnelMapEntry with mismatched type expected InvalidArgument, got %v", err)
		}

		// VNI > 0xFFFFFF -> InvalidArgument.
		_, err = c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap: proto.Uint64(mapResp.GetOid()),
			VniIdKey:  proto.Uint32(0x1000000),
		})
		if s, ok := status.FromError(err); !ok || s.Code() != codes.InvalidArgument {
			t.Errorf("CreateTunnelMapEntry with VNI > 24-bit expected InvalidArgument, got %v", err)
		}

		// Add valid entry.
		entryResp, err := c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap: proto.Uint64(mapResp.GetOid()),
			VniIdKey:  proto.Uint32(5000),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMapEntry failed: %v", err)
		}

		// Duplicate VNI in same map -> AlreadyExists.
		_, err = c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap: proto.Uint64(mapResp.GetOid()),
			VniIdKey:  proto.Uint32(5000),
		})
		if s, ok := status.FromError(err); !ok || s.Code() != codes.AlreadyExists {
			t.Errorf("CreateTunnelMapEntry with duplicate VNI expected AlreadyExists, got %v", err)
		}

		// Remove entry.
		if _, err := c.RemoveTunnelMapEntry(context.TODO(), &saipb.RemoveTunnelMapEntryRequest{Oid: entryResp.GetOid()}); err != nil {
			t.Fatalf("RemoveTunnelMapEntry failed: %v", err)
		}
		// Removing already-removed entry -> NotFound.
		_, err = c.RemoveTunnelMapEntry(context.TODO(), &saipb.RemoveTunnelMapEntryRequest{Oid: entryResp.GetOid()})
		if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound {
			t.Errorf("RemoveTunnelMapEntry already removed expected NotFound, got %v", err)
		}
	})

	t.Run("map in use and deletion restrictions", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, _, stopFn := newTestTunnel(t, dplane)
		defer stopFn()

		// Unknown map -> NotFound.
		_, err := c.RemoveTunnelMap(context.TODO(), &saipb.RemoveTunnelMapRequest{Oid: 9999})
		if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound {
			t.Errorf("RemoveTunnelMap with unknown OID expected NotFound, got %v", err)
		}

		mapResp, err := c.CreateTunnelMap(context.TODO(), &saipb.CreateTunnelMapRequest{
			Type: saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMap() failed: %v", err)
		}

		tunResp, err := c.CreateTunnel(context.TODO(), &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
			UnderlayInterface: proto.Uint64(10),
			DecapMappers:      []uint64{mapResp.GetOid()},
		})
		if err != nil {
			t.Fatalf("CreateTunnel() failed: %v", err)
		}

		// Cannot remove map referenced by tunnel -> FailedPrecondition.
		_, err = c.RemoveTunnelMap(context.TODO(), &saipb.RemoveTunnelMapRequest{Oid: mapResp.GetOid()})
		if s, ok := status.FromError(err); !ok || s.Code() != codes.FailedPrecondition {
			t.Errorf("RemoveTunnelMap referenced by tunnel expected FailedPrecondition, got %v", err)
		}

		if _, err := c.RemoveTunnel(context.TODO(), &saipb.RemoveTunnelRequest{Oid: tunResp.GetOid()}); err != nil {
			t.Fatalf("RemoveTunnel() failed: %v", err)
		}

		// Add entry to map.
		entryResp, err := c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap: proto.Uint64(mapResp.GetOid()),
			VniIdKey:  proto.Uint32(2000),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMapEntry() failed: %v", err)
		}

		// Cannot remove map with entries -> FailedPrecondition.
		_, err = c.RemoveTunnelMap(context.TODO(), &saipb.RemoveTunnelMapRequest{Oid: mapResp.GetOid()})
		if s, ok := status.FromError(err); !ok || s.Code() != codes.FailedPrecondition {
			t.Errorf("RemoveTunnelMap with entries expected FailedPrecondition, got %v", err)
		}

		if _, err := c.RemoveTunnelMapEntry(context.TODO(), &saipb.RemoveTunnelMapEntryRequest{Oid: entryResp.GetOid()}); err != nil {
			t.Fatalf("RemoveTunnelMapEntry() failed: %v", err)
		}

		// Now removable.
		if _, err := c.RemoveTunnelMap(context.TODO(), &saipb.RemoveTunnelMapRequest{Oid: mapResp.GetOid()}); err != nil {
			t.Fatalf("RemoveTunnelMap() failed: %v", err)
		}
	})

	t.Run("shared tunnel map across two tunnels", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, _, stopFn := newTestTunnel(t, dplane)
		defer stopFn()

		mapResp, err := c.CreateTunnelMap(context.TODO(), &saipb.CreateTunnelMapRequest{
			Type: saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMap() failed: %v", err)
		}

		tun1, err := c.CreateTunnel(context.TODO(), &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
			UnderlayInterface: proto.Uint64(1),
			DecapMappers:      []uint64{mapResp.GetOid()},
		})
		if err != nil {
			t.Fatalf("CreateTunnel 1 failed: %v", err)
		}

		tun2, err := c.CreateTunnel(context.TODO(), &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
			UnderlayInterface: proto.Uint64(2),
			DecapMappers:      []uint64{mapResp.GetOid()},
		})
		if err != nil {
			t.Fatalf("CreateTunnel 2 failed: %v", err)
		}

		// 2 TunnelEncap entries added so far.
		if len(dplane.gotEntryAddReqs) != 2 {
			t.Fatalf("got %d add requests, want 2", len(dplane.gotEntryAddReqs))
		}

		entryResp, err := c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap:            proto.Uint64(mapResp.GetOid()),
			VniIdKey:             proto.Uint32(7000),
			VirtualRouterIdValue: proto.Uint64(99),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMapEntry failed: %v", err)
		}

		// Now 4 entries total (2 tunnels + 2 vniToVrfTable entries).
		if len(dplane.gotEntryAddReqs) != 4 {
			t.Fatalf("got %d add requests, want 4", len(dplane.gotEntryAddReqs))
		}

		// Remove entry: removes 2 entries from vniToVrfTable.
		if _, err := c.RemoveTunnelMapEntry(context.TODO(), &saipb.RemoveTunnelMapEntryRequest{Oid: entryResp.GetOid()}); err != nil {
			t.Fatalf("RemoveTunnelMapEntry failed: %v", err)
		}
		if len(dplane.gotEntryRemoveReqs) != 2 {
			t.Fatalf("got %d remove requests, want 2", len(dplane.gotEntryRemoveReqs))
		}

		if _, err := c.RemoveTunnel(context.TODO(), &saipb.RemoveTunnelRequest{Oid: tun1.GetOid()}); err != nil {
			t.Fatalf("RemoveTunnel 1 failed: %v", err)
		}
		if _, err := c.RemoveTunnel(context.TODO(), &saipb.RemoveTunnelRequest{Oid: tun2.GetOid()}); err != nil {
			t.Fatalf("RemoveTunnel 2 failed: %v", err)
		}
		if _, err := c.RemoveTunnelMap(context.TODO(), &saipb.RemoveTunnelMapRequest{Oid: mapResp.GetOid()}); err != nil {
			t.Fatalf("RemoveTunnelMap failed: %v", err)
		}
	})

	t.Run("dataplane error rollback on CreateTunnel", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, _, stopFn := newTestTunnel(t, dplane)
		defer stopFn()

		mapResp, err := c.CreateTunnelMap(context.TODO(), &saipb.CreateTunnelMapRequest{
			Type: saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMap failed: %v", err)
		}

		_, err = c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap: proto.Uint64(mapResp.GetOid()),
			VniIdKey:  proto.Uint32(8000),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMapEntry failed: %v", err)
		}

		// Call 1 will be TunnelEncap (succeeds), Call 2 will be vniToVrfTable (fails).
		dplane.addErrOnCall = 2

		_, err = c.CreateTunnel(context.TODO(), &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
			UnderlayInterface: proto.Uint64(10),
			DecapMappers:      []uint64{mapResp.GetOid()},
		})
		if err == nil {
			t.Fatalf("CreateTunnel expected error from injected dataplane failure, got nil")
		}

		// Verify TunnelEncap was rolled back via TableEntryRemove.
		if len(dplane.gotEntryRemoveReqs) != 1 {
			t.Fatalf("got %d remove requests during rollback, want 1", len(dplane.gotEntryRemoveReqs))
		}
	})

	t.Run("dataplane error rollback on CreateTunnelMapEntry", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, _, stopFn := newTestTunnel(t, dplane)
		defer stopFn()

		mapResp, err := c.CreateTunnelMap(context.TODO(), &saipb.CreateTunnelMapRequest{
			Type: saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateTunnelMap failed: %v", err)
		}

		_, err = c.CreateTunnel(context.TODO(), &saipb.CreateTunnelRequest{
			Type:              saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
			UnderlayInterface: proto.Uint64(10),
			DecapMappers:      []uint64{mapResp.GetOid()},
		})
		if err != nil {
			t.Fatalf("CreateTunnel failed: %v", err)
		}

		// Call 1 was TunnelEncap. Call 2 will be vniToVrfTable TableEntryAdd from CreateTunnelMapEntry.
		dplane.addErrOnCall = 2

		_, err = c.CreateTunnelMapEntry(context.TODO(), &saipb.CreateTunnelMapEntryRequest{
			TunnelMap: proto.Uint64(mapResp.GetOid()),
			VniIdKey:  proto.Uint32(9000),
		})
		if err == nil {
			t.Fatalf("CreateTunnelMapEntry expected error from injected dataplane failure, got nil")
		}

		// Verify no phantom entry was stored in entryToMap by trying to remove OID 3 (mgr next ID).
		_, err = c.RemoveTunnelMapEntry(context.TODO(), &saipb.RemoveTunnelMapEntryRequest{Oid: 3})
		if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound {
			t.Errorf("RemoveTunnelMapEntry for failed entry expected NotFound, got %v", err)
		}
	})
}

func newTestTunnel(t testing.TB, api switchDataplaneAPI) (saipb.TunnelClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newTunnel(mgr, api, srv)
	})
	return saipb.NewTunnelClient(conn), mgr, stopFn
}
