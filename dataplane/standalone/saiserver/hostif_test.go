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
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestCreateHostif(t *testing.T) {
	tests := []struct {
		desc            string
		req             *saipb.CreateHostifRequest
		getInterfaceErr error
		want            *saipb.CreateHostifResponse
		wantAttr        *saipb.HostifAttribute
		wantErr         string
	}{{
		desc: "unknown type",
		req:  &saipb.CreateHostifRequest{},
		want: &saipb.CreateHostifResponse{
			Oid: 1,
		},
		wantErr: "unknown type",
	}, {
		desc: "success netdev",
		req: &saipb.CreateHostifRequest{
			Type:  saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId: proto.Uint64(2),
		},
		want: &saipb.CreateHostifResponse{
			Oid: 1,
		},
		wantAttr: &saipb.HostifAttribute{
			Type:       saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId:      proto.Uint64(2),
			OperStatus: proto.Bool(true),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			getInterface = func(name string) (*net.Interface, error) {
				return nil, tt.getInterfaceErr
			}
			dplane := &fakeSwitchDataplane{}
			c, mgr, stopFn := newTestHostif(t, dplane)
			mgr.StoreAttributes(2, &saipb.PortAttribute{
				OperStatus: saipb.PortOperStatus_PORT_OPER_STATUS_DOWN.Enum(),
			})
			defer stopFn()
			got, gotErr := c.CreateHostif(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateHostif() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateHostif() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.HostifAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateHostif() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestSetHostifAttribute(t *testing.T) {
	tests := []struct {
		desc            string
		req             *saipb.SetHostifAttributeRequest
		getInterfaceErr error
		wantAttr        *saipb.HostifAttribute
		wantReq         *fwdpb.PortStateRequest
		wantErr         string
	}{{
		desc: "oper status",
		req: &saipb.SetHostifAttributeRequest{
			Oid:        1,
			OperStatus: proto.Bool(false),
		},
		wantReq: &fwdpb.PortStateRequest{
			PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
			Operation: &fwdpb.PortInfo{AdminStatus: fwdpb.PortState_PORT_STATE_DISABLED_DOWN},
			ContextId: &fwdpb.ContextId{Id: "foo"},
		},
		wantAttr: &saipb.HostifAttribute{
			OperStatus: proto.Bool(false),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			getInterface = func(name string) (*net.Interface, error) {
				return nil, tt.getInterfaceErr
			}
			dplane := &fakeSwitchDataplane{}
			c, mgr, stopFn := newTestHostif(t, dplane)
			defer stopFn()
			_, gotErr := c.SetHostifAttribute(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("SetHostifAttribute() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotPortStateReq[0], tt.wantReq, protocmp.Transform()); d != "" {
				t.Errorf("SetHostifAttribute() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.HostifAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("SetHostifAttribute() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func newTestHostif(t testing.TB, api switchDataplaneAPI) (saipb.HostifClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newHostif(mgr, api, srv)
	})
	return saipb.NewHostifClient(conn), mgr, stopFn
}
