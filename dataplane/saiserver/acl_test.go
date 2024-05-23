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
	"net/netip"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestCreateAclEntry(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.CreateAclEntryRequest
		wantErr string
		want    *fwdpb.TableEntryAddRequest
	}{{
		desc:    "table not member of a group",
		wantErr: "FailedPrecondition",
		req:     &saipb.CreateAclEntryRequest{},
	}, {
		desc:    "no fields",
		wantErr: "InvalidArgument",
		req: &saipb.CreateAclEntryRequest{
			TableId: proto.Uint64(1),
		},
	}, {
		desc: "all fields",
		req: &saipb.CreateAclEntryRequest{
			TableId: proto.Uint64(1),
			FieldDstIp: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataIp{
					DataIp: []byte{127, 0, 0, 1},
				},
				Mask: &saipb.AclFieldData_MaskIp{
					MaskIp: []byte{255, 255, 255, 0},
				},
			},
			FieldInPort: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataOid{DataOid: 1},
			},
			FieldDscp: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataUint{DataUint: 10},
				Mask: &saipb.AclFieldData_MaskUint{MaskUint: 0xff},
			},
			FieldDstIpv6Word3: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataIp{
					DataIp: netip.MustParseAddr("cafe:beef::").AsSlice(),
				},
				Mask: &saipb.AclFieldData_MaskIp{
					MaskIp: netip.MustParseAddr("ffff:ffff::").AsSlice(),
				},
			},
			FieldDstIpv6Word2: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataIp{
					DataIp: netip.MustParseAddr("0:0:cafe:beef::").AsSlice(),
				},
				Mask: &saipb.AclFieldData_MaskIp{
					MaskIp: netip.MustParseAddr("0:0:ffff:ffff::").AsSlice(),
				},
			},
			FieldDstMac: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataMac{
					DataMac: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
				},
				Mask: &saipb.AclFieldData_MaskMac{
					MaskMac: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
			},
			FieldEtherType: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataUint{DataUint: 0x800},
				Mask: &saipb.AclFieldData_MaskUint{MaskUint: 0xfff},
			},
			FieldIcmpv6Type: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataUint{DataUint: 0x01},
				Mask: &saipb.AclFieldData_MaskUint{MaskUint: 0xff},
			},
			FieldIpProtocol: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataUint{DataUint: 0x04},
				Mask: &saipb.AclFieldData_MaskUint{MaskUint: 0xFF},
			},
			FieldL4DstPort: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataUint{DataUint: 22},
				Mask: &saipb.AclFieldData_MaskUint{MaskUint: 0xffff},
			},
			FieldSrcMac: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataMac{
					DataMac: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
				},
				Mask: &saipb.AclFieldData_MaskMac{
					MaskMac: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
			},
			FieldTtl: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataUint{DataUint: 0x01},
				Mask: &saipb.AclFieldData_MaskUint{MaskUint: 0xff},
			},
		},
		want: &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id: 1,
						Fields: []*fwdpb.PacketFieldMaskedBytes{{
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
							Bytes:   []byte{127, 0, 0, 1},
							Masks:   []byte{255, 255, 255, 0},
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT}},
							Bytes:   binary.BigEndian.AppendUint64(nil, 1),
							Masks:   binary.BigEndian.AppendUint64(nil, math.MaxUint64),
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS}},
							Bytes:   []byte{10},
							Masks:   []byte{255},
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
							Bytes:   netip.MustParseAddr("cafe:beef::").AsSlice(),
							Masks:   netip.MustParseAddr("ffff:ffff::").AsSlice(),
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
							Bytes:   netip.MustParseAddr("0:0:cafe:beef::").AsSlice(),
							Masks:   netip.MustParseAddr("0:0:ffff:ffff::").AsSlice(),
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST}},
							Bytes:   []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
							Masks:   []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE}},
							Bytes:   []byte{0x8, 0x0},
							Masks:   []byte{0x0f, 0xff},
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ICMP_TYPE}},
							Bytes:   []byte{0x01},
							Masks:   []byte{0xff},
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO}},
							Bytes:   []byte{0x04},
							Masks:   []byte{0xff},
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST}},
							Bytes:   []byte{0, 22},
							Masks:   []byte{0xff, 0xff},
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC}},
							Bytes:   []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
							Masks:   []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
						}, {
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP}},
							Bytes:   []byte{0x01},
							Masks:   []byte{0xff},
						}},
					},
				},
			},
		},
	}, {
		desc: "vrf action",
		req: &saipb.CreateAclEntryRequest{
			TableId: proto.Uint64(1),
			FieldDstIp: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataIp{
					DataIp: []byte{127, 0, 0, 1},
				},
				Mask: &saipb.AclFieldData_MaskIp{
					MaskIp: []byte{255, 255, 255, 0},
				},
			},
			ActionSetVrf: &saipb.AclActionData{
				Parameter: &saipb.AclActionData_Oid{
					Oid: 1,
				},
			},
		},
		want: &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id: 1,
						Fields: []*fwdpb.PacketFieldMaskedBytes{{
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
							Bytes:   []byte{127, 0, 0, 1},
							Masks:   []byte{255, 255, 255, 0},
						}},
					},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
				Action: &fwdpb.ActionDesc_Update{
					Update: &fwdpb.UpdateActionDesc{
						FieldId: &fwdpb.PacketFieldId{
							Field: &fwdpb.PacketField{
								FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF,
							},
						},
						Field: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{}},
						Type:  fwdpb.UpdateType_UPDATE_TYPE_SET,
						Value: binary.BigEndian.AppendUint64(nil, 1),
					},
				},
			}},
		},
	}, {
		desc: "user trap action",
		req: &saipb.CreateAclEntryRequest{
			TableId: proto.Uint64(1),
			FieldDstIp: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataIp{
					DataIp: []byte{127, 0, 0, 1},
				},
				Mask: &saipb.AclFieldData_MaskIp{
					MaskIp: []byte{255, 255, 255, 0},
				},
			},
			ActionSetUserTrapId: &saipb.AclActionData{
				Parameter: &saipb.AclActionData_Oid{
					Oid: 1,
				},
			},
		},
		want: &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id: 1,
						Fields: []*fwdpb.PacketFieldMaskedBytes{{
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
							Bytes:   []byte{127, 0, 0, 1},
							Masks:   []byte{255, 255, 255, 0},
						}},
					},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
				Action: &fwdpb.ActionDesc_Update{
					Update: &fwdpb.UpdateActionDesc{
						FieldId: &fwdpb.PacketFieldId{
							Field: &fwdpb.PacketField{
								FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TRAP_ID,
							},
						},
						Field: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{}},
						Type:  fwdpb.UpdateType_UPDATE_TYPE_SET,
						Value: binary.BigEndian.AppendUint64(nil, 1),
					},
				},
			}},
		},
	}, {
		desc: "drop action",
		req: &saipb.CreateAclEntryRequest{
			TableId: proto.Uint64(1),
			FieldDstIp: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataIp{
					DataIp: []byte{127, 0, 0, 1},
				},
				Mask: &saipb.AclFieldData_MaskIp{
					MaskIp: []byte{255, 255, 255, 0},
				},
			},
			ActionPacketAction: &saipb.AclActionData{
				Parameter: &saipb.AclActionData_PacketAction{
					PacketAction: saipb.PacketAction_PACKET_ACTION_DROP,
				},
			},
		},
		want: &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id: 1,
						Fields: []*fwdpb.PacketFieldMaskedBytes{{
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
							Bytes:   []byte{127, 0, 0, 1},
							Masks:   []byte{255, 255, 255, 0},
						}},
					},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_DROP,
			}},
		},
	}, {
		desc: "forward action",
		req: &saipb.CreateAclEntryRequest{
			TableId: proto.Uint64(1),
			FieldDstIp: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataIp{
					DataIp: []byte{127, 0, 0, 1},
				},
				Mask: &saipb.AclFieldData_MaskIp{
					MaskIp: []byte{255, 255, 255, 0},
				},
			},
			ActionPacketAction: &saipb.AclActionData{
				Parameter: &saipb.AclActionData_PacketAction{
					PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD,
				},
			},
		},
		want: &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id: 1,
						Fields: []*fwdpb.PacketFieldMaskedBytes{{
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
							Bytes:   []byte{127, 0, 0, 1},
							Masks:   []byte{255, 255, 255, 0},
						}},
					},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE,
			}},
		},
	}, {
		desc: "counter action",
		req: &saipb.CreateAclEntryRequest{
			TableId: proto.Uint64(1),
			FieldDstIp: &saipb.AclFieldData{
				Data: &saipb.AclFieldData_DataIp{
					DataIp: []byte{127, 0, 0, 1},
				},
				Mask: &saipb.AclFieldData_MaskIp{
					MaskIp: []byte{255, 255, 255, 0},
				},
			},
			ActionCounter: &saipb.AclActionData{
				Parameter: &saipb.AclActionData_Oid{
					Oid: 10,
				},
			},
		},
		want: &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id: 1,
						Fields: []*fwdpb.PacketFieldMaskedBytes{{
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
							Bytes:   []byte{127, 0, 0, 1},
							Masks:   []byte{255, 255, 255, 0},
						}},
					},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_FLOW_COUNTER,
				Action: &fwdpb.ActionDesc_Flow{
					Flow: &fwdpb.FlowCounterActionDesc{
						CounterId: &fwdpb.FlowCounterId{
							ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(10)},
						},
					},
				},
			}},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{
				ctx: fwdcontext.New("foo", "foo"),
			}
			dplane.ctx.Objects.Insert(&fwdobject.Base{}, &fwdpb.ObjectId{Id: "1"})
			c, a, stopFn := newTestACL(t, dplane)
			a.tableToLocation[1] = tableLocation{
				groupID: "1",
				bank:    0,
			}
			defer stopFn()
			_, gotErr := c.CreateAclEntry(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateAclEntry() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotEntryAddReqs[0], tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateAclEntry() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestRemoveAclEntry(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.RemoveAclEntryRequest
		wantErr string
		want    *fwdpb.TableEntryRemoveRequest
	}{{
		desc: "not found",
		req: &saipb.RemoveAclEntryRequest{
			Oid: 2,
		},
		wantErr: "FailedPrecondition",
	}, {
		desc: "success",
		req: &saipb.RemoveAclEntryRequest{
			Oid: 1,
		},
		want: &fwdpb.TableEntryRemoveRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Flow{
					Flow: &fwdpb.FlowEntryDesc{
						Id: 1,
						Fields: []*fwdpb.PacketFieldMaskedBytes{{
							FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
							Bytes:   []byte{127, 0, 0, 1},
							Masks:   []byte{255, 255, 255, 0},
						}},
					},
				},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}

			c, a, stopFn := newTestACL(t, dplane)
			defer stopFn()
			a.tableToLocation[1] = tableLocation{
				groupID: "1",
				bank:    0,
			}
			_, err := c.CreateAclEntry(context.TODO(), &saipb.CreateAclEntryRequest{
				TableId: proto.Uint64(1),
				FieldDstIp: &saipb.AclFieldData{
					Data: &saipb.AclFieldData_DataIp{
						DataIp: []byte{127, 0, 0, 1},
					},
					Mask: &saipb.AclFieldData_MaskIp{
						MaskIp: []byte{255, 255, 255, 0},
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			_, gotErr := c.RemoveAclEntry(context.Background(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("RemoveAclCounter() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotEntryRemoveReqs[0], tt.want, protocmp.Transform()); d != "" {
				t.Errorf("RemoveAclCounter() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestCreateAclCounter(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.CreateAclCounterRequest
		wantErr string
		want    *fwdpb.FlowCounterCreateRequest
	}{{
		desc: "success",
		req:  &saipb.CreateAclCounterRequest{},
		want: &fwdpb.FlowCounterCreateRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			Id:        &fwdpb.FlowCounterId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, _, stopFn := newTestACL(t, dplane)
			defer stopFn()
			_, gotErr := c.CreateAclCounter(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateAclCounter() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotFlowCounterCreateReqs[0], tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateAclCounter() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestRemoveAclCounter(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.RemoveAclCounterRequest
		wantErr string
		want    *fwdpb.ObjectDeleteRequest
	}{{
		desc: "not found",
		req: &saipb.RemoveAclCounterRequest{
			Oid: 2,
		},
		wantErr: "not found",
	}, {
		desc: "success",
		req: &saipb.RemoveAclCounterRequest{
			Oid: 1,
		},
		want: &fwdpb.ObjectDeleteRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			ObjectId:  &fwdpb.ObjectId{Id: "1"},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{
				ctx: fwdcontext.New("foo", "foo"),
			}
			dplane.ctx.Objects.Insert(&fwdobject.Base{}, &fwdpb.ObjectId{Id: "1"})
			c, a, stopFn := newTestACL(t, dplane)
			a.mgr.StoreAttributes(1, &saipb.CreateAclCounterRequest{EnablePacketCount: proto.Bool(true)})
			defer stopFn()
			_, gotErr := c.RemoveAclCounter(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("RemoveAclCounter() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotObjectDeleteReqs[0], tt.want, protocmp.Transform()); d != "" {
				t.Errorf("RemoveAclCounter() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestGetAclCounterAttribute(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.GetAclCounterAttributeRequest
		wantErr string
		want    *saipb.GetAclCounterAttributeResponse
	}{{
		desc: "success",
		req: &saipb.GetAclCounterAttributeRequest{
			Oid:      1,
			AttrType: []saipb.AclCounterAttr{saipb.AclCounterAttr_ACL_COUNTER_ATTR_PACKETS},
		},
		want: &saipb.GetAclCounterAttributeResponse{
			Attr: &saipb.AclCounterAttribute{
				Packets: proto.Uint64(1),
				Bytes:   proto.Uint64(0),
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{
				flowQueryReplies: []*fwdpb.FlowCounterQueryReply{{
					Counters: []*fwdpb.FlowCounter{{
						Packets: 1,
					}},
				}},
			}
			c, _, stopFn := newTestACL(t, dplane)
			defer stopFn()
			got, gotErr := c.GetAclCounterAttribute(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("GetAclCounterAttribute() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("GetAclCounterAttribute() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func newTestACL(t testing.TB, api switchDataplaneAPI) (saipb.AclClient, *acl, func()) {
	var a *acl
	conn, _, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		a = newACL(mgr, api, srv)
	})
	return saipb.NewAclClient(conn), a, stopFn
}
