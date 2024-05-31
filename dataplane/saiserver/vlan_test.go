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

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"

	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
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
	if len(testm) != 2 && len(defaultm) != 2 {
		t.Fatalf("Failed to create a non-default VLAN")
	}

	// Remove one port from the non-default VLAN, and ensure it is moved to the default VLAN.
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
	if len(testm) != 1 && len(defaultm) != 3 {
		t.Errorf("Failed to remove a non-default VLAN member")
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
	// Ensure all of the 4 ports are now under the default VLAN.
	if len(defaultm) != 4 {
		t.Errorf("Failed to remove the non-default VLAN")
	}
}

func newTestVlan(t testing.TB, api switchDataplaneAPI) (saipb.VlanClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newVlan(mgr, api, srv)
	})
	return saipb.NewVlanClient(conn), mgr, stopFn
}
