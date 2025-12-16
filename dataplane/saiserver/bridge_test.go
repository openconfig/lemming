// Copyright 2025 Google LLC
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

func TestBridgePort(t *testing.T) {
	tests := []struct {
		desc     string
		req      *saipb.CreateBridgePortRequest
		want     *saipb.CreateBridgePortResponse
		wantAttr *saipb.BridgePortAttribute
		wantErr  string
	}{{
		desc: "success",
		req: &saipb.CreateBridgePortRequest{
			Switch:          1,
			Type:            saipb.BridgePortType_BRIDGE_PORT_TYPE_PORT.Enum(),
			PortId:          proto.Uint64(10),
			FdbLearningMode: saipb.BridgePortFdbLearningMode_BRIDGE_PORT_FDB_LEARNING_MODE_DISABLE.Enum(),
			AdminState:      proto.Bool(true),
		},
		want: &saipb.CreateBridgePortResponse{
			Oid: 1,
		},
		wantAttr: &saipb.BridgePortAttribute{
			Type:            saipb.BridgePortType_BRIDGE_PORT_TYPE_PORT.Enum(),
			PortId:          proto.Uint64(10),
			FdbLearningMode: saipb.BridgePortFdbLearningMode_BRIDGE_PORT_FDB_LEARNING_MODE_DISABLE.Enum(),
			AdminState:      proto.Bool(true),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			conn, mgr, stopFn := newTestBridge(t)
			defer stopFn()
			c := saipb.NewBridgeClient(conn)

			got, gotErr := c.CreateBridgePort(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateBridgePort() unexpected err: %s w/ error %v", diff, gotErr)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateBridgePort() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.BridgePortAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateBridgePort() failed: diff(-got,+want)\n:%s", d)
			}

			// Test Set Attribute
			_, err := c.SetBridgePortAttribute(context.TODO(), &saipb.SetBridgePortAttributeRequest{
				Oid:        got.Oid,
				AdminState: proto.Bool(false),
			})
			if err != nil {
				t.Fatalf("SetBridgePortAttribute() failed: %v", err)
			}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if attr.GetAdminState() != false {
				t.Errorf("SetBridgePortAttribute() failed: got %v, want false", attr.GetAdminState())
			}

			// Test Get Attribute
			resp, err := c.GetBridgePortAttribute(context.TODO(), &saipb.GetBridgePortAttributeRequest{
				Oid:      got.Oid,
				AttrType: []saipb.BridgePortAttr{saipb.BridgePortAttr_BRIDGE_PORT_ATTR_ADMIN_STATE},
			})
			if err != nil {
				t.Fatalf("GetBridgePortAttribute() failed: %v", err)
			}
			if resp.GetAttr().GetAdminState() != false {
				t.Errorf("GetBridgePortAttribute() failed: got %v, want false", resp.GetAttr().GetAdminState())
			}
			// Test Remove
			if _, err := c.RemoveBridgePort(context.TODO(), &saipb.RemoveBridgePortRequest{Oid: got.Oid}); err != nil {
				t.Fatalf("RemoveBridgePort() failed: %v", err)
			}
		})
	}
}

func newTestBridge(t testing.TB) (grpc.ClientConnInterface, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newBridge(mgr, nil, srv)
	})
	return conn, mgr, stopFn
}
