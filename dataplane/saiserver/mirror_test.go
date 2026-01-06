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

func TestMirrorSession(t *testing.T) {
	tests := []struct {
		desc     string
		req      *saipb.CreateMirrorSessionRequest
		want     *saipb.CreateMirrorSessionResponse
		wantAttr *saipb.MirrorSessionAttribute
		wantErr  string
	}{{
		desc: "success",
		req: &saipb.CreateMirrorSessionRequest{
			Switch:          1,
			Type:            saipb.MirrorSessionType_MIRROR_SESSION_TYPE_ENHANCED_REMOTE.Enum(),
			MonitorPort:     proto.Uint64(10),
			VlanId:          proto.Uint32(100),
			SrcIpAddress:    []byte{192, 168, 1, 1},
			DstIpAddress:    []byte{192, 168, 1, 2},
			SrcMacAddress:   []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			DstMacAddress:   []byte{0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb},
			GreProtocolType: proto.Uint32(0x8800),
		},
		want: &saipb.CreateMirrorSessionResponse{
			Oid: 1,
		},
		wantAttr: &saipb.MirrorSessionAttribute{
			Type:            saipb.MirrorSessionType_MIRROR_SESSION_TYPE_ENHANCED_REMOTE.Enum(),
			MonitorPort:     proto.Uint64(10),
			VlanId:          proto.Uint32(100),
			SrcIpAddress:    []byte{192, 168, 1, 1},
			DstIpAddress:    []byte{192, 168, 1, 2},
			SrcMacAddress:   []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			DstMacAddress:   []byte{0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb},
			GreProtocolType: proto.Uint32(0x8800),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			conn, mgr, stopFn := newTestMirror(t)
			defer stopFn()
			c := saipb.NewMirrorClient(conn)

			got, gotErr := c.CreateMirrorSession(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateMirrorSession() unexpected err: %s w/ error %v", diff, gotErr)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateMirrorSession() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.MirrorSessionAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateMirrorSession() failed: diff(-got,+want)\n:%s", d)
			}

			// Test Set Attribute
			_, err := c.SetMirrorSessionAttribute(context.TODO(), &saipb.SetMirrorSessionAttributeRequest{
				Oid: got.Oid,
				Tc:  proto.Uint32(5),
			})
			if err != nil {
				t.Fatalf("SetMirrorSessionAttribute() failed: %v", err)
			}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if attr.GetTc() != 5 {
				t.Errorf("SetMirrorSessionAttribute() failed: got %v, want 5", attr.GetTc())
			}

			// Test Get Attribute
			resp, err := c.GetMirrorSessionAttribute(context.TODO(), &saipb.GetMirrorSessionAttributeRequest{
				Oid:      got.Oid,
				AttrType: []saipb.MirrorSessionAttr{saipb.MirrorSessionAttr_MIRROR_SESSION_ATTR_TC},
			})
			if err != nil {
				t.Fatalf("GetMirrorSessionAttribute() failed: %v", err)
			}
			if resp.GetAttr().GetTc() != 5 {
				t.Errorf("GetMirrorSessionAttribute() failed: got %v, want 5", resp.GetAttr().GetTc())
			}
			// Test Remove
			if _, err := c.RemoveMirrorSession(context.TODO(), &saipb.RemoveMirrorSessionRequest{Oid: got.Oid}); err != nil {
				t.Fatalf("RemoveMirrorSession() failed: %v", err)
			}
		})
	}
}

func newTestMirror(t testing.TB) (grpc.ClientConnInterface, *attrmgr.AttrMgr, func()) {
	// We cannot use newTestServer because it initializes all services in saiserver.go,
	// but the test server in saiserver_test.go only registers the service passed in the callback.
	// However, saiserver.go RegisterEntrypointServer uses the struct Server which we can use
	// but we need to populate it.
	// To simplify, let's use the pattern from saiserver_test.go if available, or just manually set up.
	// Looking at vlan_test.go, it uses newTestServer.
	// Let's modify the newTestServer callback to register Mirror server.
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		// We need to register the mirror server.
		// In saiserver.go New() creates everything.
		// We can just create the mirror struct directly.
		m := &mirror{mgr: mgr}
		saipb.RegisterMirrorServer(srv, m)
	})
	return conn, mgr, stopFn
}

// Helper getter to access connection from AttrMgr for testing if needed, or we just return it.
// In vlan_test.go: return saipb.NewVlanClient(conn), mgr, stopFn
// We can do similar.
