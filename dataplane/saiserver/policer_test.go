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

func TestCreatePolicer(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.CreatePolicerRequest
		wantErr string
		want    *fwdpb.TableEntryAddRequest
	}{{
		desc: "trap action",
		req: &saipb.CreatePolicerRequest{
			GreenPacketAction: saipb.PacketAction_PACKET_ACTION_TRAP.Enum(),
		},
		want: &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: policerTabler}},
			Entries: []*fwdpb.TableEntryAddRequest_Entry{{
				EntryDesc: &fwdpb.EntryDesc{
					Entry: &fwdpb.EntryDesc_Exact{
						Exact: &fwdpb.ExactEntryDesc{
							Fields: []*fwdpb.PacketFieldBytes{{
								FieldId: &fwdpb.PacketFieldId{
									Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_POLICER_ID},
								},
								Bytes: binary.BigEndian.AppendUint64(nil, 2),
							}},
						},
					},
				},
				Actions: []*fwdpb.ActionDesc{{
					ActionType: fwdpb.ActionType_ACTION_TYPE_TRANSMIT,
					Action: &fwdpb.ActionDesc_Transmit{
						Transmit: &fwdpb.TransmitActionDesc{
							PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "10"}},
							Immediate: true,
						},
					},
				}},
			}},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, p, stopFn := newTestPolicer(t, dplane)
			p.mgr.StoreAttributes(p.mgr.NextID(), &saipb.SwitchAttribute{
				CpuPort: proto.Uint64(10),
			})
			defer stopFn()
			_, gotErr := c.CreatePolicer(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreatePolicer() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotEntryAddReqs[0], tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreatePolicer() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestRemovePolicer(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.RemovePolicerRequest
		wantErr string
		want    *fwdpb.TableEntryRemoveRequest
	}{{
		desc: "not found",
		req: &saipb.RemovePolicerRequest{
			Oid: 3,
		},
		wantErr: "not found",
	}, {
		desc: "success",
		req: &saipb.RemovePolicerRequest{
			Oid: 2,
		},
		want: &fwdpb.TableEntryRemoveRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: policerTabler}},
			Entries: []*fwdpb.EntryDesc{{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_POLICER_ID},
							},
							Bytes: binary.BigEndian.AppendUint64(nil, 2),
						}},
					},
				},
			}},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, a, stopFn := newTestPolicer(t, dplane)
			a.mgr.StoreAttributes(a.mgr.NextID(), &saipb.SwitchAttribute{
				CpuPort: proto.Uint64(10),
			})
			defer stopFn()
			_, err := c.CreatePolicer(context.TODO(), &saipb.CreatePolicerRequest{
				GreenPacketAction: saipb.PacketAction_PACKET_ACTION_TRAP.Enum(),
			})
			if err != nil {
				t.Fatal(err)
			}
			_, gotErr := c.RemovePolicer(context.Background(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("RemovePolicer() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotEntryRemoveReqs[0], tt.want, protocmp.Transform()); d != "" {
				t.Errorf("RemovePolicer() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func newTestPolicer(t testing.TB, api switchDataplaneAPI) (saipb.PolicerClient, *policer, func()) {
	var p *policer
	conn, _, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		p = newPolicer(mgr, api, srv)
	})
	return saipb.NewPolicerClient(conn), p, stopFn
}
