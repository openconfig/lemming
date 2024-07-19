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

func TestCreateL2McGroup(t *testing.T) {
	tests := []struct {
		desc     string
		req      *saipb.CreateL2McGroupRequest
		wantAttr *saipb.L2McGroupAttribute
		wantErr  string
	}{
		{
			desc: "success",
			wantAttr: &saipb.L2McGroupAttribute{
				L2McOutputCount: proto.Uint32(0),
				L2McMemberList:  []uint64{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, mg, stopFn := newTestL2McGroup(t, dplane)
			defer stopFn()
			resp, gotErr := c.CreateL2McGroup(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateL2McGroup() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			attr := &saipb.L2McGroupAttribute{}
			if err := mg.mgr.PopulateAllAttributes(fmt.Sprint(resp.GetOid()), attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateL2McGroup() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestCreateL2McGroupMember(t *testing.T) {
	bridgePort1Id := uint64(101)
	port1Id := uint64(1001)

	tests := []struct {
		desc           string
		withBridgePort bool
		req            *saipb.CreateL2McGroupMemberRequest
		want           *fwdpb.TableEntryAddRequest
		wantAttr       *saipb.L2McGroupMemberAttribute
		wantErr        string
	}{{
		desc: "no bridge port",
		req: &saipb.CreateL2McGroupMemberRequest{
			Switch:       1,
			L2McOutputId: proto.Uint64(111),
		},
		wantErr: "failed to populate OutputId",
	}, {
		desc:           "no group",
		withBridgePort: true,
		req: &saipb.CreateL2McGroupMemberRequest{
			Switch:       1,
			L2McGroupId:  proto.Uint64(999),
			L2McOutputId: proto.Uint64(100),
		},
		wantErr: "cannot find L2MC group with group ID",
	}, {
		desc:           "success",
		withBridgePort: true,
		req: &saipb.CreateL2McGroupMemberRequest{
			Switch:       1,
			L2McOutputId: proto.Uint64(bridgePort1Id),
		},
		want: &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: L2MCGroupTable}},
			Entries: []*fwdpb.TableEntryAddRequest_Entry{
				{
					EntryDesc: &fwdpb.EntryDesc{
						Entry: &fwdpb.EntryDesc_Exact{
							Exact: &fwdpb.ExactEntryDesc{
								Fields: []*fwdpb.PacketFieldBytes{{
									FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID}},
									Bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
								}},
							},
						},
					},
					Actions: []*fwdpb.ActionDesc{
						{
							ActionType: fwdpb.ActionType_ACTION_TYPE_MIRROR,
							Action: &fwdpb.ActionDesc_Mirror{
								Mirror: &fwdpb.MirrorActionDesc{PortId: &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(port1Id)}}},
							},
						},
					},
				},
			},
		},
		wantAttr: &saipb.L2McGroupMemberAttribute{
			L2McGroupId:  proto.Uint64(1),
			L2McOutputId: proto.Uint64(bridgePort1Id),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, mg, stopFn := newTestL2McGroup(t, dplane)
			if tt.withBridgePort {
				mg.mgr.StoreAttributes(bridgePort1Id, &saipb.BridgePortAttribute{PortId: proto.Uint64(port1Id)})
			}
			defer stopFn()
			ctx := context.TODO()
			resp, err := c.CreateL2McGroup(ctx, &saipb.CreateL2McGroupRequest{
				Switch: *proto.Uint64(1),
			})
			if err != nil {
				t.Fatalf("failed to create L2MC group: %v", err)
			}
			if tt.req.L2McGroupId == nil {
				tt.req.L2McGroupId = &resp.Oid
			}
			mResp, gotErr := c.CreateL2McGroupMember(ctx, tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateL2McGroupMember() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotEntryAddReqs[0], tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateL2McGroupMember() request check failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.L2McGroupMemberAttribute{}
			if err := mg.mgr.PopulateAllAttributes(fmt.Sprint(mResp.GetOid()), attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateL2McGroupMember() group attribute check failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestRemoveL2McGroupMember(t *testing.T) {
	bridgePort1Id := uint64(101)
	bridgePort2Id := uint64(102)
	port1Id := uint64(1001)
	port2Id := uint64(1002)

	tests := []struct {
		desc          string
		req           *saipb.RemoveL2McGroupMemberRequest
		want          *fwdpb.TableEntryAddRequest
		wantAttr      *saipb.L2McGroupMemberAttribute
		wantGroupAttr *saipb.L2McGroupAttribute
		wantErr       string
	}{{
		desc: "non-existing member",
		req: &saipb.RemoveL2McGroupMemberRequest{
			Oid: 123,
		},
		wantErr: "cannot find L2MC group member with OID",
	}, {
		desc: "success",
		req: &saipb.RemoveL2McGroupMemberRequest{
			Oid: 2, // first member.
		},
		want: &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: L2MCGroupTable}},
			Entries: []*fwdpb.TableEntryAddRequest_Entry{
				{
					EntryDesc: &fwdpb.EntryDesc{
						Entry: &fwdpb.EntryDesc_Exact{
							Exact: &fwdpb.ExactEntryDesc{
								Fields: []*fwdpb.PacketFieldBytes{{
									FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID}},
									Bytes:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
								}},
							},
						},
					},
					Actions: []*fwdpb.ActionDesc{
						{
							ActionType: fwdpb.ActionType_ACTION_TYPE_MIRROR,
							Action: &fwdpb.ActionDesc_Mirror{
								Mirror: &fwdpb.MirrorActionDesc{PortId: &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(port1Id)}}},
							},
						},
					},
				},
			},
		},
		wantGroupAttr: &saipb.L2McGroupAttribute{
			L2McOutputCount: proto.Uint32(1),
			L2McMemberList:  []uint64{3},
		},
		wantAttr: &saipb.L2McGroupMemberAttribute{
			L2McGroupId:  proto.Uint64(1),
			L2McOutputId: proto.Uint64(bridgePort2Id),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, mg, stopFn := newTestL2McGroup(t, dplane)
			mg.mgr.StoreAttributes(bridgePort1Id, &saipb.BridgePortAttribute{PortId: proto.Uint64(port1Id)})
			mg.mgr.StoreAttributes(bridgePort2Id, &saipb.BridgePortAttribute{PortId: proto.Uint64(port2Id)})
			defer stopFn()
			ctx := context.TODO()
			// Create one group with 2 members (OID 2 and 3 repectively).
			resp, err := c.CreateL2McGroup(ctx, &saipb.CreateL2McGroupRequest{
				Switch: *proto.Uint64(1),
			})
			if err != nil {
				t.Fatalf("failed to create L2MC group: %v", err)
			}
			groupId := resp.GetOid()
			_, err = c.CreateL2McGroupMember(ctx, &saipb.CreateL2McGroupMemberRequest{
				L2McGroupId:  &groupId,
				L2McOutputId: proto.Uint64(bridgePort1Id),
			})
			if err != nil {
				t.Fatal(err)
			}
			mResp2, err := c.CreateL2McGroupMember(ctx, &saipb.CreateL2McGroupMemberRequest{
				L2McGroupId:  &groupId,
				L2McOutputId: proto.Uint64(bridgePort2Id),
			})
			if err != nil {
				t.Fatal(err)
			}
			// Remove the 1st member if the test gets here.
			_, gotErr := c.RemoveL2McGroupMember(ctx, tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("RemoveL2McGroupMember() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}

			if d := cmp.Diff(dplane.gotEntryAddReqs[0], tt.want, protocmp.Transform()); d != "" {
				t.Errorf("RemoveL2McGroupMember() request check failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.L2McGroupMemberAttribute{}
			if err := mg.mgr.PopulateAllAttributes(fmt.Sprint(mResp2.GetOid()), attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("RemoveL2McGroupMember() member attribute check failed: diff(-got,+want)\n:%s", d)
			}
			gAttr := &saipb.L2McGroupAttribute{}
			if err := mg.mgr.PopulateAllAttributes(fmt.Sprint(groupId), gAttr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(gAttr, tt.wantGroupAttr, protocmp.Transform()); d != "" {
				t.Errorf("RemoveL2McGroupMember() group attribute check failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func newTestL2McGroup(t testing.TB, api switchDataplaneAPI) (saipb.L2McGroupClient, *l2mcGroup, func()) {
	var p *l2mcGroup
	conn, _, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		p = newL2mcGroup(mgr, api, srv)
	})
	return saipb.NewL2McGroupClient(conn), p, stopFn
}
