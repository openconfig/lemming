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
	dplane := &fakeSwitchDataplane{}
	c, _, stopFn := newTestTunnel(t, dplane)
	defer stopFn()

	req := &saipb.CreateTunnelTermTableEntryRequest{
		TunnelType: saipb.TunnelType_TUNNEL_TYPE_VXLAN.Enum(),
		Type:       saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
		SrcIp:      []byte{1, 2, 3, 4},
		DstIp:      []byte{5, 6, 7, 8},
		VrId:       proto.Uint64(10),
	}
	resp, err := c.CreateTunnelTermTableEntry(context.TODO(), req)
	if err != nil {
		t.Fatalf("CreateTunnelTermTableEntry() unexpected err: %v", err)
	}
	if resp.GetOid() == 0 {
		t.Errorf("CreateTunnelTermTableEntry() returned 0 OID")
	}

	_, err = c.RemoveTunnelTermTableEntry(context.TODO(), &saipb.RemoveTunnelTermTableEntryRequest{
		Oid: resp.GetOid(),
	})
	if err != nil {
		t.Fatalf("RemoveTunnelTermTableEntry() unexpected err: %v", err)
	}
}

func TestTunnelMapAndEntry(t *testing.T) {
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
}

func newTestTunnel(t testing.TB, api switchDataplaneAPI) (saipb.TunnelClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newTunnel(mgr, api, srv)
	})
	return saipb.NewTunnelClient(conn), mgr, stopFn
}
