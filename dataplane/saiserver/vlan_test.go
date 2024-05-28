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

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"

	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

func TestCreateVlan(t *testing.T) {
	defaultStpInstId := uint64(101)
	defaultVlanId := uint32(4095)
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
			VlanId:       proto.Uint32(defaultVlanId),
			LearnDisable: proto.Bool(true),
		},
		want: &saipb.CreateVlanResponse{
			Oid: 1,
		},
		wantAttr: &saipb.VlanAttribute{
			StpInstance:                        proto.Uint64(defaultStpInstId),
			VlanId:                             proto.Uint32(defaultVlanId),
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
			VlanId:       proto.Uint32(defaultVlanId),
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
				DefaultStpInstId: proto.Uint64(defaultStpInstId),
			})
			if tt.hasExistingVlan {
				if _, err := c.CreateVlan(context.TODO(), &saipb.CreateVlanRequest{
					Switch:       1,
					VlanId:       proto.Uint32(defaultVlanId),
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

func newTestVlan(t testing.TB, api switchDataplaneAPI) (saipb.VlanClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newVlan(mgr, api, srv)
	})
	return saipb.NewVlanClient(conn), mgr, stopFn
}
