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
	"encoding/binary"
	"fmt"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	testStpInstId = uint64(101)
)

func TestCreateVlan(t *testing.T) {
	tests := []struct {
		desc            string
		hasExistingVlan bool
		req             *saipb.CreateVlanRequest
		want            *saipb.CreateVlanResponse
		wantAttr        *saipb.VlanAttribute
		wantErr         string
	}{{
		desc: "success",
		req: &saipb.CreateVlanRequest{
			Switch:       1,
			VlanId:       proto.Uint32(DefaultVlanId),
			LearnDisable: proto.Bool(true),
		},
		want: &saipb.CreateVlanResponse{
			Oid: 1,
		},
		wantAttr: &saipb.VlanAttribute{
			StpInstance:                        proto.Uint64(testStpInstId),
			VlanId:                             proto.Uint32(DefaultVlanId),
			LearnDisable:                       proto.Bool(true),
			UnknownNonIpMcastOutputGroupId:     proto.Uint64(0),
			UnknownIpv4McastOutputGroupId:      proto.Uint64(0),
			UnknownIpv6McastOutputGroupId:      proto.Uint64(0),
			UnknownLinklocalMcastOutputGroupId: proto.Uint64(0),
			IngressAcl:                         proto.Uint64(0),
			EgressAcl:                          proto.Uint64(0),
			UnknownUnicastFloodGroup:           proto.Uint64(0),
			UnknownMulticastFloodGroup:         proto.Uint64(0),
			BroadcastFloodGroup:                proto.Uint64(0),
		},
	}, {
		desc:            "existing VLAN",
		hasExistingVlan: true,
		req: &saipb.CreateVlanRequest{
			Switch:       1,
			VlanId:       proto.Uint32(DefaultVlanId),
			LearnDisable: proto.Bool(true),
		},
		wantErr: "found existing VLAN",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, mgr, stopFn := newTestVlan(t, dplane)
			defer stopFn()
			mgr.StoreAttributes(1, &saipb.SwitchAttribute{
				DefaultStpInstId: proto.Uint64(testStpInstId),
			})
			if tt.hasExistingVlan {
				if _, err := c.CreateVlan(context.TODO(), &saipb.CreateVlanRequest{
					Switch:       1,
					VlanId:       proto.Uint32(DefaultVlanId),
					LearnDisable: proto.Bool(true),
				}); err != nil {
					t.Fatalf("failed to create VLAN")
				}
			}
			got, gotErr := c.CreateVlan(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateVlan() unexpected err: %s w/ error %v", diff, gotErr)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateVlan() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.VlanAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateVlan() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

// TestVlanOperations is an end-to-end test.
func TestVlanOperations(t *testing.T) {
	testPort1Id := uint64(11)
	testPort2Id := uint64(12)
	testPort3Id := uint64(13)
	testPort4Id := uint64(14)
	testVlanId := uint32(10)
	dplane := &fakeSwitchDataplane{}
	c, mgr, stopFn := newTestVlan(t, dplane)
	defer stopFn()
	mgr.StoreAttributes(1, &saipb.SwitchAttribute{
		DefaultStpInstId: proto.Uint64(testStpInstId),
	})
	// Crete bridge port mapping in attrmgr for each port
	mgr.StoreAttributes(11, &saipb.BridgePortAttribute{PortId: proto.Uint64(11)})
	mgr.StoreAttributes(12, &saipb.BridgePortAttribute{PortId: proto.Uint64(12)})
	mgr.StoreAttributes(13, &saipb.BridgePortAttribute{PortId: proto.Uint64(13)})
	mgr.StoreAttributes(14, &saipb.BridgePortAttribute{PortId: proto.Uint64(14)})
	ctx := context.TODO()

	getVLANMembers := func(vlanOID uint64) ([]uint64, error) {
		vlanAttrReq := &saipb.GetVlanAttributeRequest{Oid: vlanOID, AttrType: []saipb.VlanAttr{saipb.VlanAttr_VLAN_ATTR_MEMBER_LIST}}
		vlanAttrResp := &saipb.GetVlanAttributeResponse{}
		if err := mgr.PopulateAttributes(vlanAttrReq, vlanAttrResp); err != nil {
			return nil, fmt.Errorf("failed to populate attributes of default VLAN.")
		}
		return vlanAttrResp.GetAttr().GetMemberList(), nil
	}

	// By default, the default VLAN has all the ports (port1 to port4).
	resp, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{
		Switch:       1,
		VlanId:       proto.Uint32(DefaultVlanId),
		LearnDisable: proto.Bool(true),
	})
	if err != nil {
		t.Fatalf("Failed to create VLAN: %v", err)
	}
	vlanOID := resp.GetOid()
	resps, err := c.CreateVlanMembers(ctx, &saipb.CreateVlanMembersRequest{
		Reqs: []*saipb.CreateVlanMemberRequest{
			{
				Switch:       1,
				VlanId:       &vlanOID,
				BridgePortId: &testPort1Id,
			}, {
				Switch:       1,
				VlanId:       &vlanOID,
				BridgePortId: &testPort2Id,
			}, {
				Switch:       1,
				VlanId:       &vlanOID,
				BridgePortId: &testPort3Id,
			}, {
				Switch:       1,
				VlanId:       &vlanOID,
				BridgePortId: &testPort4Id,
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to add ports to the default VLAN: %v", err)
	}
	defVLANMembers := []uint64{}
	for _, r := range resps.GetResps() {
		defVLANMembers = append(defVLANMembers, r.GetOid())
	}
	// Ensure default VLAN's member list is correct.
	members, err := getVLANMembers(vlanOID)
	if err != nil {
		t.Fatalf("Failed to populate attributes of the default VLAN: %v", err)
	}
	if diff := cmp.Diff(defVLANMembers, members); diff != "" {
		t.Fatalf("Failed to create VLAN members: %v", diff)
	}

	// Creates a non-default VLAN with two ports. These two ports are moved from
	// the default VLAN to the new VLAN.
	resp, err = c.CreateVlan(ctx, &saipb.CreateVlanRequest{
		Switch:       1,
		VlanId:       proto.Uint32(testVlanId),
		LearnDisable: proto.Bool(true),
	})
	if err != nil {
		t.Fatalf("Failed to create a non-default VLAN %d: %v", testVlanId, err)
	}
	testVlanOID := resp.GetOid()
	resps, err = c.CreateVlanMembers(ctx, &saipb.CreateVlanMembersRequest{
		Reqs: []*saipb.CreateVlanMemberRequest{
			{
				Switch:       1,
				VlanId:       &testVlanOID,
				BridgePortId: &testPort1Id,
			}, {
				Switch:       1,
				VlanId:       &testVlanOID,
				BridgePortId: &testPort2Id,
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to add ports to the non-default VLAN: %v", err)
	}
	testVLANMembers := []uint64{}
	for _, r := range resps.GetResps() {
		testVLANMembers = append(testVLANMembers, r.GetOid())
	}
	testm, err := getVLANMembers(testVlanOID)
	if err != nil {
		t.Fatalf("Failed to get members of the non-default VLAN: %v", err)
	}
	defaultm, err := getVLANMembers(vlanOID)
	if err != nil {
		t.Fatalf("Failed to get members of the default VLAN: %v", err)
	}
	if len(testm) != 2 || len(defaultm) != 2 {
		t.Fatalf("Failed to create a non-default VLAN: got testm=%v defaultm=%v", testm, defaultm)
	}

	// Remove one port from the non-default VLAN.
	_, err = c.RemoveVlanMember(ctx, &saipb.RemoveVlanMemberRequest{
		Oid: testVLANMembers[0],
	})
	if err != nil {
		t.Fatalf("Failed to remove VLAN member (OID=%d): %v", testVLANMembers[0], err)
	}
	testm, err = getVLANMembers(testVlanOID)
	if err != nil {
		t.Fatalf("Failed to get members of the non-default VLAN: %v", err)
	}
	defaultm, err = getVLANMembers(vlanOID)
	if err != nil {
		t.Fatalf("Failed to get members of default VLAN: %v", err)
	}
	if len(testm) != 1 || len(defaultm) != 2 {
		t.Errorf("Failed to remove a non-default VLAN member: got testm=%v defaultm=%v", testm, defaultm)
	}

	// Remove the non-default VLAN.
	_, err = c.RemoveVlan(ctx, &saipb.RemoveVlanRequest{
		Oid: testVlanOID,
	})
	if err != nil {
		t.Fatalf("Failed to remove the non-default VLAN %d: %v", testVlanId, err)
	}
	defaultm, err = getVLANMembers(vlanOID)
	if err != nil {
		t.Fatalf("Failed to get members of the default VLAN: %v", err)
	}
	// Ensure the default VLAN only has 2 members.
	if len(defaultm) != 2 {
		t.Errorf("Failed to remove the non-default VLAN")
	}
}

func newTestVlan(t testing.TB, api switchDataplaneAPI) (saipb.VlanClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newVlan(mgr, api, srv)
	})
	return saipb.NewVlanClient(conn), mgr, stopFn
}

func TestCreateVlanMemberLogicalMapping(t *testing.T) {
	dplane := &fakeSwitchDataplane{}
	c, mgr, stopFn := newTestVlan(t, dplane)
	defer stopFn()
	ctx := context.TODO()

	mgr.StoreAttributes(1, &saipb.SwitchAttribute{
		DefaultStpInstId: proto.Uint64(10),
	})

	// Create vlan.
	resp, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{
		Switch: 1,
		VlanId: proto.Uint32(10),
	})
	if err != nil {
		t.Fatalf("CreateVlan failed: %v", err)
	}
	vlanOid := resp.Oid

	// Create mock bridge port mapping in attrmgr for given port id.
	bridgePortOid := uint64(100)
	testPortId := uint64(10)
	mgr.StoreAttributes(bridgePortOid, &saipb.BridgePortAttribute{
		PortId: proto.Uint64(testPortId),
	})

	// Create vlan member using bridge port oid.
	_, err = c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
		Switch:       1,
		VlanId:       &vlanOid,
		BridgePortId: &bridgePortOid,
	})
	if err != nil {
		t.Fatalf("CreateVlanMember failed: %v", err)
	}

	// Verify that ObjectNID was queried with the port id.
	found := false
	for _, req := range dplane.gotObjectNIDReqs {
		if req.GetObjectId().GetId() == "10" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ObjectNID was not queried with port id 10, got reqs: %+v", dplane.gotObjectNIDReqs)
	}
}

func TestRemoveVlan(t *testing.T) {
	dplane := &fakeSwitchDataplane{}
	c, mgr, stopFn := newTestVlan(t, dplane)
	defer stopFn()
	ctx := context.TODO()

	mgr.StoreAttributes(1, &saipb.SwitchAttribute{
		DefaultStpInstId: proto.Uint64(testStpInstId),
	})

	// Create vlan.
	vlanResp, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{
		Switch: 1,
		VlanId: proto.Uint32(10),
	})
	if err != nil {
		t.Fatalf("CreateVlan failed: %v", err)
	}
	vlanOid := vlanResp.Oid

	// Create mapping for mock bridge ports.
	bridgePortOid := uint64(100)
	physicalPortId := uint64(10)
	mgr.StoreAttributes(bridgePortOid, &saipb.BridgePortAttribute{
		PortId: proto.Uint64(physicalPortId),
	})

	// Create vlan member with bridge port.
	memberResp, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
		Switch:       1,
		VlanId:       &vlanOid,
		BridgePortId: &bridgePortOid,
	})
	if err != nil {
		t.Fatalf("CreateVlanMember failed: %v", err)
	}

	// Remove vlan member.
	_, err = c.RemoveVlanMember(ctx, &saipb.RemoveVlanMemberRequest{
		Oid: memberResp.Oid,
	})
	if err != nil {
		t.Fatalf("RemoveVlanMember failed: %v", err)
	}

	// Remove vlan.
	_, err = c.RemoveVlan(ctx, &saipb.RemoveVlanRequest{
		Oid: vlanOid,
	})
	if err != nil {
		t.Fatalf("RemoveVlan failed: %v", err)
	}
}

func TestCreateVlanMemberTaggingMode(t *testing.T) {
	const (
		vlanID           = uint32(666)
		untaggedPortID   = uint64(10)
		taggedPortID     = uint64(20)
		untaggedBridgeID = uint64(100)
		taggedBridgeID   = uint64(200)
	)

	tests := []struct {
		desc        string
		bridgePort  uint64
		taggingMode saipb.VlanTaggingMode
		wantActions []*fwdpb.ActionDesc
	}{{
		desc:        "untagged member encapsulates 802.1Q header and sets VLAN tag",
		bridgePort:  untaggedBridgeID,
		taggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_UNTAGGED,
		wantActions: []*fwdpb.ActionDesc{
			fwdconfig.Action(fwdconfig.EncapAction(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN)).Build(),
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG).WithValue(binary.BigEndian.AppendUint16(nil, uint16(vlanID)))).Build(),
		},
	}, {
		desc:        "tagged member preserves existing 802.1Q header without re-encapsulating",
		bridgePort:  taggedBridgeID,
		taggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_TAGGED,
		wantActions: nil,
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, mgr, stopFn := newTestVlan(t, dplane)
			defer stopFn()
			ctx := context.Background()

			mgr.StoreAttributes(1, &saipb.SwitchAttribute{
				DefaultStpInstId: proto.Uint64(testStpInstId),
			})
			mgr.StoreAttributes(untaggedBridgeID, &saipb.BridgePortAttribute{
				PortId: proto.Uint64(untaggedPortID),
			})
			mgr.StoreAttributes(taggedBridgeID, &saipb.BridgePortAttribute{
				PortId: proto.Uint64(taggedPortID),
			})

			vlanResp, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{
				Switch: 1,
				VlanId: proto.Uint32(vlanID),
			})
			if err != nil {
				t.Fatalf("CreateVlan failed: %v", err)
			}

			_, err = c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
				Switch:          1,
				VlanId:          proto.Uint64(vlanResp.GetOid()),
				BridgePortId:    proto.Uint64(tt.bridgePort),
				VlanTaggingMode: tt.taggingMode.Enum(),
			})
			if err != nil {
				t.Fatalf("CreateVlanMember failed: %v", err)
			}

			if len(dplane.gotEntryAddReqs) < 1 {
				t.Fatalf("got %d TableEntryAdd requests, want at least 1", len(dplane.gotEntryAddReqs))
			}
			vlanEntryReq := dplane.gotEntryAddReqs[len(dplane.gotEntryAddReqs)-1]
			if vlanEntryReq.GetTableId().GetObjectId().GetId() != VlanTable {
				t.Fatalf("got TableId %q, want %q", vlanEntryReq.GetTableId().GetObjectId().GetId(), VlanTable)
			}
			gotEntries := vlanEntryReq.GetEntries()
			if len(gotEntries) != 1 {
				t.Fatalf("got %d table entries, want 1", len(gotEntries))
			}
			if d := cmp.Diff(gotEntries[0].GetActions(), tt.wantActions, protocmp.Transform()); d != "" {
				t.Errorf("CreateVlanMember() actions diff (-got +want):\n%s", d)
			}
		})
	}
}

func TestCreateVlanMemberHybridPort(t *testing.T) {
	const (
		vlan42ID  = uint32(42)
		vlan100ID = uint32(100)
		vlan666ID = uint32(666)
		portID    = uint64(55)
		bridgeID  = uint64(550)
	)

	getVLANMembers := func(t *testing.T, mgr *attrmgr.AttrMgr, vlanOID uint64) []uint64 {
		t.Helper()
		req := &saipb.GetVlanAttributeRequest{Oid: vlanOID, AttrType: []saipb.VlanAttr{saipb.VlanAttr_VLAN_ATTR_MEMBER_LIST}}
		resp := &saipb.GetVlanAttributeResponse{}
		if err := mgr.PopulateAttributes(req, resp); err != nil {
			t.Fatalf("PopulateAttributes(%d) failed: %v", vlanOID, err)
		}
		return resp.GetAttr().GetMemberList()
	}

	wantUntaggedActions := func(vid uint32) []*fwdpb.ActionDesc {
		return []*fwdpb.ActionDesc{
			fwdconfig.Action(fwdconfig.EncapAction(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET_VLAN)).Build(),
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG).WithValue(binary.BigEndian.AppendUint16(nil, uint16(vid)))).Build(),
		}
	}

	findTableAddEntry := func(dplane *fakeSwitchDataplane, tableID string, vlanTag uint16) *fwdpb.TableEntryAddRequest_Entry {
		tagBytes := binary.BigEndian.AppendUint16(nil, vlanTag)
		for i := len(dplane.gotEntryAddReqs) - 1; i >= 0; i-- {
			req := dplane.gotEntryAddReqs[i]
			if req.GetTableId().GetObjectId().GetId() != tableID {
				continue
			}
			for _, entry := range req.GetEntries() {
				for _, f := range entry.GetEntryDesc().GetExact().GetFields() {
					if f.GetFieldId().GetField().GetFieldNum() == fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG &&
						cmp.Equal(f.GetBytes(), tagBytes) {
						return entry
					}
				}
			}
		}
		return nil
	}

	t.Run("untagged first then tagged (hybrid port)", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, mgr, stopFn := newTestVlan(t, dplane)
		defer stopFn()
		ctx := context.Background()

		mgr.StoreAttributes(1, &saipb.SwitchAttribute{DefaultStpInstId: proto.Uint64(testStpInstId)})
		mgr.StoreAttributes(bridgeID, &saipb.BridgePortAttribute{PortId: proto.Uint64(portID)})

		v42, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(vlan42ID)})
		if err != nil {
			t.Fatalf("CreateVlan(42) failed: %v", err)
		}
		v666, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(vlan666ID)})
		if err != nil {
			t.Fatalf("CreateVlan(666) failed: %v", err)
		}

		m42, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v42.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_UNTAGGED.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateVlanMember(42, UNTAGGED) failed: %v", err)
		}
		m666, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v666.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_TAGGED.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateVlanMember(666, TAGGED) failed: %v", err)
		}

		if got := getVLANMembers(t, mgr, v42.GetOid()); !cmp.Equal(got, []uint64{m42.GetOid()}) {
			t.Errorf("VLAN 42 members = %v, want [%d]", got, m42.GetOid())
		}
		if got := getVLANMembers(t, mgr, v666.GetOid()); !cmp.Equal(got, []uint64{m666.GetOid()}) {
			t.Errorf("VLAN 666 members = %v, want [%d]", got, m666.GetOid())
		}

		untaggedEntry := findTableAddEntry(dplane, VlanTable, 0)
		if untaggedEntry == nil {
			t.Fatalf("missing VlanTable entry for VLAN_TAG=0")
		}
		if d := cmp.Diff(untaggedEntry.GetActions(), wantUntaggedActions(vlan42ID), protocmp.Transform()); d != "" {
			t.Errorf("VlanTable(VLAN_TAG=0) actions diff (-got +want):\n%s", d)
		}
		taggedEntry := findTableAddEntry(dplane, VlanTable, uint16(vlan666ID))
		if taggedEntry == nil {
			t.Fatalf("missing VlanTable entry for VLAN_TAG=%d", vlan666ID)
		}
		if len(taggedEntry.GetActions()) != 0 {
			t.Errorf("VlanTable(VLAN_TAG=%d) actions = %v, want empty (Continue)", vlan666ID, taggedEntry.GetActions())
		}
		egress42 := findTableAddEntry(dplane, EgressVlanTable, uint16(vlan42ID))
		if egress42 == nil {
			t.Fatalf("missing EgressVlanTable entry for VLAN_TAG=%d", vlan42ID)
		}
		egress666 := findTableAddEntry(dplane, EgressVlanTable, uint16(vlan666ID))
		if egress666 != nil {
			t.Errorf("unexpected EgressVlanTable entry for tagged VLAN %d: %v", vlan666ID, egress666)
		}
	})

	t.Run("tagged first then untagged (hybrid port)", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, mgr, stopFn := newTestVlan(t, dplane)
		defer stopFn()
		ctx := context.Background()

		mgr.StoreAttributes(1, &saipb.SwitchAttribute{DefaultStpInstId: proto.Uint64(testStpInstId)})
		mgr.StoreAttributes(bridgeID, &saipb.BridgePortAttribute{PortId: proto.Uint64(portID)})

		v42, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(vlan42ID)})
		if err != nil {
			t.Fatalf("CreateVlan(42) failed: %v", err)
		}
		v666, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(vlan666ID)})
		if err != nil {
			t.Fatalf("CreateVlan(666) failed: %v", err)
		}

		m666, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v666.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_TAGGED.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateVlanMember(666, TAGGED) failed: %v", err)
		}
		m42, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v42.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_UNTAGGED.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateVlanMember(42, UNTAGGED) failed: %v", err)
		}

		if got := getVLANMembers(t, mgr, v42.GetOid()); !cmp.Equal(got, []uint64{m42.GetOid()}) {
			t.Errorf("VLAN 42 members = %v, want [%d]", got, m42.GetOid())
		}
		if got := getVLANMembers(t, mgr, v666.GetOid()); !cmp.Equal(got, []uint64{m666.GetOid()}) {
			t.Errorf("VLAN 666 members = %v, want [%d]", got, m666.GetOid())
		}

		untaggedEntry := findTableAddEntry(dplane, VlanTable, 0)
		if untaggedEntry == nil {
			t.Fatalf("missing VlanTable entry for VLAN_TAG=0")
		}
		if d := cmp.Diff(untaggedEntry.GetActions(), wantUntaggedActions(vlan42ID), protocmp.Transform()); d != "" {
			t.Errorf("VlanTable(VLAN_TAG=0) actions diff (-got +want):\n%s", d)
		}
		taggedEntry := findTableAddEntry(dplane, VlanTable, uint16(vlan666ID))
		if taggedEntry == nil {
			t.Fatalf("missing VlanTable entry for VLAN_TAG=%d", vlan666ID)
		}
		if len(taggedEntry.GetActions()) != 0 {
			t.Errorf("VlanTable(VLAN_TAG=%d) actions = %v, want empty (Continue)", vlan666ID, taggedEntry.GetActions())
		}
	})

	t.Run("tagged then tagged (multi-VLAN trunk port)", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, mgr, stopFn := newTestVlan(t, dplane)
		defer stopFn()
		ctx := context.Background()

		mgr.StoreAttributes(1, &saipb.SwitchAttribute{DefaultStpInstId: proto.Uint64(testStpInstId)})
		mgr.StoreAttributes(bridgeID, &saipb.BridgePortAttribute{PortId: proto.Uint64(portID)})

		v100, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(vlan100ID)})
		if err != nil {
			t.Fatalf("CreateVlan(100) failed: %v", err)
		}
		v666, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(vlan666ID)})
		if err != nil {
			t.Fatalf("CreateVlan(666) failed: %v", err)
		}

		m100, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v100.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_TAGGED.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateVlanMember(100, TAGGED) failed: %v", err)
		}
		m666, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v666.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_TAGGED.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateVlanMember(666, TAGGED) failed: %v", err)
		}

		if got := getVLANMembers(t, mgr, v100.GetOid()); !cmp.Equal(got, []uint64{m100.GetOid()}) {
			t.Errorf("VLAN 100 members = %v, want [%d]", got, m100.GetOid())
		}
		if got := getVLANMembers(t, mgr, v666.GetOid()); !cmp.Equal(got, []uint64{m666.GetOid()}) {
			t.Errorf("VLAN 666 members = %v, want [%d]", got, m666.GetOid())
		}
	})

	t.Run("replace untagged VLAN on hybrid port while preserving tagged VLAN", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, mgr, stopFn := newTestVlan(t, dplane)
		defer stopFn()
		ctx := context.Background()

		mgr.StoreAttributes(1, &saipb.SwitchAttribute{DefaultStpInstId: proto.Uint64(testStpInstId)})
		mgr.StoreAttributes(bridgeID, &saipb.BridgePortAttribute{PortId: proto.Uint64(portID)})

		v1, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(DefaultVlanId)})
		if err != nil {
			t.Fatalf("CreateVlan(1) failed: %v", err)
		}
		v42, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(vlan42ID)})
		if err != nil {
			t.Fatalf("CreateVlan(42) failed: %v", err)
		}
		v666, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(vlan666ID)})
		if err != nil {
			t.Fatalf("CreateVlan(666) failed: %v", err)
		}

		if _, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v1.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_UNTAGGED.Enum(),
		}); err != nil {
			t.Fatalf("CreateVlanMember(1, UNTAGGED) failed: %v", err)
		}
		m666, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v666.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_TAGGED.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateVlanMember(666, TAGGED) failed: %v", err)
		}
		m42, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v42.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_UNTAGGED.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateVlanMember(42, UNTAGGED) failed: %v", err)
		}

		if got := getVLANMembers(t, mgr, v1.GetOid()); len(got) != 0 {
			t.Errorf("VLAN 1 members = %v, want empty after untagged replacement", got)
		}
		if got := getVLANMembers(t, mgr, v42.GetOid()); !cmp.Equal(got, []uint64{m42.GetOid()}) {
			t.Errorf("VLAN 42 members = %v, want [%d]", got, m42.GetOid())
		}
		if got := getVLANMembers(t, mgr, v666.GetOid()); !cmp.Equal(got, []uint64{m666.GetOid()}) {
			t.Errorf("VLAN 666 members = %v, want [%d]", got, m666.GetOid())
		}
	})

	t.Run("remove tagged member on hybrid port preserves untagged member", func(t *testing.T) {
		dplane := &fakeSwitchDataplane{}
		c, mgr, stopFn := newTestVlan(t, dplane)
		defer stopFn()
		ctx := context.Background()

		mgr.StoreAttributes(1, &saipb.SwitchAttribute{DefaultStpInstId: proto.Uint64(testStpInstId)})
		mgr.StoreAttributes(bridgeID, &saipb.BridgePortAttribute{PortId: proto.Uint64(portID)})

		v42, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(vlan42ID)})
		if err != nil {
			t.Fatalf("CreateVlan(42) failed: %v", err)
		}
		v666, err := c.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: 1, VlanId: proto.Uint32(vlan666ID)})
		if err != nil {
			t.Fatalf("CreateVlan(666) failed: %v", err)
		}

		m42, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v42.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_UNTAGGED.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateVlanMember(42, UNTAGGED) failed: %v", err)
		}
		m666, err := c.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
			Switch:          1,
			VlanId:          proto.Uint64(v666.GetOid()),
			BridgePortId:    proto.Uint64(bridgeID),
			VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_TAGGED.Enum(),
		})
		if err != nil {
			t.Fatalf("CreateVlanMember(666, TAGGED) failed: %v", err)
		}

		if _, err := c.RemoveVlanMember(ctx, &saipb.RemoveVlanMemberRequest{Oid: m666.GetOid()}); err != nil {
			t.Fatalf("RemoveVlanMember(666) failed: %v", err)
		}
		if got := getVLANMembers(t, mgr, v666.GetOid()); len(got) != 0 {
			t.Errorf("VLAN 666 members = %v, want empty after RemoveVlanMember", got)
		}
		if got := getVLANMembers(t, mgr, v42.GetOid()); !cmp.Equal(got, []uint64{m42.GetOid()}) {
			t.Errorf("VLAN 42 members = %v, want [%d]", got, m42.GetOid())
		}
	})
}

func TestVlanLagAndBridgeLearn(t *testing.T) {
	dplane := &fakeSwitchDataplane{}
	var sw *saiSwitch
	var fdbSrv *fdb
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		var err error
		sw, err = newSwitch(mgr, dplane, srv, &dplaneopts.Options{})
		if err != nil {
			t.Fatalf("newSwitch failed: %v", err)
		}
		fdbSrv, err = newFdb(mgr, dplane, srv)
		if err != nil {
			t.Fatalf("newFdb failed: %v", err)
		}
		sw.fdb = fdbSrv
	})
	defer stopFn()
	ctx := context.Background()

	swClient := saipb.NewSwitchClient(conn)
	vlanClient := saipb.NewVlanClient(conn)
	lagClient := saipb.NewLagClient(conn)

	swResp, err := swClient.CreateSwitch(ctx, &saipb.CreateSwitchRequest{})
	if err != nil {
		t.Fatalf("CreateSwitch failed: %v", err)
	}

	// Create LAG (PortChannel1) and add member port 10 before CreateVlanMember,
	// and member port 20 after CreateVlanMember.
	lagResp, err := lagClient.CreateLag(ctx, &saipb.CreateLagRequest{Switch: swResp.GetOid()})
	if err != nil {
		t.Fatalf("CreateLag failed: %v", err)
	}
	if gotType := mgr.GetType(fmt.Sprint(lagResp.GetOid())); gotType != saipb.ObjectType_OBJECT_TYPE_LAG {
		t.Fatalf("LAG object type = %v, want OBJECT_TYPE_LAG", gotType)
	}

	const (
		physPort1   = uint64(10)
		physPort2   = uint64(20)
		hybridPort  = uint64(55)
		lagBridgeID = uint64(500)
		hybBridgeID = uint64(550)
	)
	if _, err := lagClient.CreateLagMember(ctx, &saipb.CreateLagMemberRequest{
		Switch: swResp.GetOid(),
		LagId:  proto.Uint64(lagResp.GetOid()),
		PortId: proto.Uint64(physPort1),
	}); err != nil {
		t.Fatalf("CreateLagMember(physPort1) failed: %v", err)
	}

	mgr.StoreAttributes(lagBridgeID, &saipb.BridgePortAttribute{PortId: proto.Uint64(lagResp.GetOid())})
	mgr.StoreAttributes(hybBridgeID, &saipb.BridgePortAttribute{PortId: proto.Uint64(hybridPort)})

	v42, err := vlanClient.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: swResp.GetOid(), VlanId: proto.Uint32(42)})
	if err != nil {
		t.Fatalf("CreateVlan(42) failed: %v", err)
	}
	v666, err := vlanClient.CreateVlan(ctx, &saipb.CreateVlanRequest{Switch: swResp.GetOid(), VlanId: proto.Uint32(666)})
	if err != nil {
		t.Fatalf("CreateVlan(666) failed: %v", err)
	}

	// Add LAG to VLAN 42 (untagged). This should program physPort1.
	if _, err := vlanClient.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
		Switch:          swResp.GetOid(),
		VlanId:          proto.Uint64(v42.GetOid()),
		BridgePortId:    proto.Uint64(lagBridgeID),
		VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_UNTAGGED.Enum(),
	}); err != nil {
		t.Fatalf("CreateVlanMember(LAG, 42) failed: %v", err)
	}

	// Add physPort2 to LAG after CreateVlanMember. This should dynamically program physPort2 into VLAN 42.
	if _, err := lagClient.CreateLagMember(ctx, &saipb.CreateLagMemberRequest{
		Switch: swResp.GetOid(),
		LagId:  proto.Uint64(lagResp.GetOid()),
		PortId: proto.Uint64(physPort2),
	}); err != nil {
		t.Fatalf("CreateLagMember(physPort2) failed: %v", err)
	}

	// Configure hybridPort in VLAN 42 (untagged) and VLAN 666 (tagged).
	if _, err := vlanClient.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
		Switch:          swResp.GetOid(),
		VlanId:          proto.Uint64(v42.GetOid()),
		BridgePortId:    proto.Uint64(hybBridgeID),
		VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_UNTAGGED.Enum(),
	}); err != nil {
		t.Fatalf("CreateVlanMember(hybrid, 42) failed: %v", err)
	}
	if _, err := vlanClient.CreateVlanMember(ctx, &saipb.CreateVlanMemberRequest{
		Switch:          swResp.GetOid(),
		VlanId:          proto.Uint64(v666.GetOid()),
		BridgePortId:    proto.Uint64(hybBridgeID),
		VlanTaggingMode: saipb.VlanTaggingMode_VLAN_TAGGING_MODE_TAGGED.Enum(),
	}); err != nil {
		t.Fatalf("CreateVlanMember(hybrid, 666) failed: %v", err)
	}

	// Verify onBridgeLearn resolves the exact VLAN OID for both untagged (42) and tagged (666)
	// packets on hybridPort, as well as physical member ports of the LAG.
	notifCh := make(chan *saipb.FdbEventNotificationResponse, 10)
	unsub := fdbSrv.subscribe(notifCh)
	defer unsub()

	sw.onBridgeLearn([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x42}, fmt.Sprint(hybridPort), 42)
	notif42 := <-notifCh
	if gotBv := notif42.GetData()[0].GetFdbEntry().GetBvId(); gotBv != v42.GetOid() {
		t.Errorf("onBridgeLearn(hybrid, vid=42) BvId = %d, want %d (VLAN 42 OID)", gotBv, v42.GetOid())
	}

	sw.onBridgeLearn([]byte{0x00, 0x01, 0x02, 0x03, 0x06, 0x66}, fmt.Sprint(hybridPort), 666)
	notif666 := <-notifCh
	if gotBv := notif666.GetData()[0].GetFdbEntry().GetBvId(); gotBv != v666.GetOid() {
		t.Errorf("onBridgeLearn(hybrid, vid=666) BvId = %d, want %d (VLAN 666 OID)", gotBv, v666.GetOid())
	}

	sw.onBridgeLearn([]byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x20}, fmt.Sprint(physPort2), 42)
	notifLag := <-notifCh
	if gotBv := notifLag.GetData()[0].GetFdbEntry().GetBvId(); gotBv != v42.GetOid() {
		t.Errorf("onBridgeLearn(physPort2 in LAG, vid=42) BvId = %d, want %d", gotBv, v42.GetOid())
	}
	if gotBp := notifLag.GetData()[0].GetAttrs()[0].GetBridgePortId(); gotBp != lagBridgeID {
		t.Errorf("onBridgeLearn(physPort2 in LAG, vid=42) BridgePortId = %d, want %d (LAG bridge port)", gotBp, lagBridgeID)
	}
}
