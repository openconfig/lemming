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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	dpb "github.com/openconfig/lemming/proto/dataplane"
)

type fakeRoutingDataplaneAPI struct {
	gotAddNeighborReq     []*dpb.AddNeighborRequest
	gotAddNextHopGroupReq []*dpb.AddNextHopGroupRequest
	gotAddNextHopReq      []*dpb.AddNextHopRequest
}

func (f *fakeRoutingDataplaneAPI) AddNeighbor(_ context.Context, req *dpb.AddNeighborRequest) (*dpb.AddNeighborResponse, error) {
	f.gotAddNeighborReq = append(f.gotAddNeighborReq, req)
	return nil, nil
}

func (f *fakeRoutingDataplaneAPI) AddNextHopGroup(_ context.Context, req *dpb.AddNextHopGroupRequest) (*dpb.AddNextHopGroupResponse, error) {
	f.gotAddNextHopGroupReq = append(f.gotAddNextHopGroupReq, req)
	return nil, nil
}

func (f *fakeRoutingDataplaneAPI) AddNextHop(_ context.Context, req *dpb.AddNextHopRequest) (*dpb.AddNextHopResponse, error) {
	f.gotAddNextHopReq = append(f.gotAddNextHopReq, req)
	return nil, nil
}

func (f *fakeRoutingDataplaneAPI) AddIPRoute(context.Context, *dpb.AddIPRouteRequest) (*dpb.AddIPRouteResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (f *fakeRoutingDataplaneAPI) AddInterface(context.Context, *dpb.AddInterfaceRequest) (*dpb.AddInterfaceResponse, error) {
	panic("not implemented") // TODO: Implement
}

func TestCreateNeighborEntry(t *testing.T) {
	tests := []struct {
		desc     string
		req      *saipb.CreateNeighborEntryRequest
		want     *saipb.CreateNeighborEntryResponse
		wantAttr *saipb.NeighborEntryAttribute
		wantErr  string
	}{{
		desc:     "existing interface",
		req:      &saipb.CreateNeighborEntryRequest{},
		want:     &saipb.CreateNeighborEntryResponse{},
		wantAttr: &saipb.NeighborEntryAttribute{},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeRoutingDataplaneAPI{}
			c, mgr, stopFn := newTestNeighbor(t, dplane)
			defer stopFn()
			got, gotErr := c.CreateNeighborEntry(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateNeighborEntry() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateNeighborEntry() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.NeighborEntryAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateNeighborEntry() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestCreateNextHopGroup(t *testing.T) {
	tests := []struct {
		desc     string
		req      *saipb.CreateNextHopGroupRequest
		wantAttr *saipb.NextHopGroupAttribute
		wantReq  *dpb.AddNextHopGroupRequest
		wantErr  string
	}{{
		desc:    "unspeficied type",
		req:     &saipb.CreateNextHopGroupRequest{},
		wantErr: "InvalidArgument",
	}, {
		desc: "success",
		req: &saipb.CreateNextHopGroupRequest{
			Type: saipb.NextHopGroupType_NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP.Enum(),
		},
		wantAttr: &saipb.NextHopGroupAttribute{
			Type: saipb.NextHopGroupType_NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP.Enum(),
		},
		wantReq: &dpb.AddNextHopGroupRequest{
			Id: 1,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeRoutingDataplaneAPI{}
			c, mgr, stopFn := newTestNextHopGroup(t, dplane)
			defer stopFn()
			_, gotErr := c.CreateNextHopGroup(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateNextHopGroup() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotAddNextHopGroupReq[0], tt.wantReq, protocmp.Transform()); d != "" {
				t.Errorf("CreateNextHopGroup() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.NextHopGroupAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateNextHopGroup() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestCreateNextHopGroupMember(t *testing.T) {
	tests := []struct {
		desc     string
		req      *saipb.CreateNextHopGroupMemberRequest
		wantAttr *saipb.NextHopGroupMemberAttribute
		wantReq  *dpb.AddNextHopGroupRequest
		wantErr  string
	}{{
		desc: "success",
		req: &saipb.CreateNextHopGroupMemberRequest{
			NextHopGroupId: proto.Uint64(1),
			NextHopId:      proto.Uint64(2),
			Weight:         proto.Uint32(3),
		},
		wantReq: &dpb.AddNextHopGroupRequest{
			Id: 1,
			List: &dpb.NextHopIDList{
				Hops:    []uint64{2},
				Weights: []uint64{3},
			},
			Mode: dpb.GroupUpdateMode_GROUP_UPDATE_MODE_APPEND,
		},
		wantAttr: &saipb.NextHopGroupMemberAttribute{
			NextHopGroupId: proto.Uint64(1),
			NextHopId:      proto.Uint64(2),
			Weight:         proto.Uint32(3),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeRoutingDataplaneAPI{}
			c, mgr, stopFn := newTestNextHopGroup(t, dplane)
			defer stopFn()
			_, gotErr := c.CreateNextHopGroupMember(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateNextHopGroupMember() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotAddNextHopGroupReq[0], tt.wantReq, protocmp.Transform()); d != "" {
				t.Errorf("CreateNextHopGroupMember() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.NextHopGroupMemberAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateNextHopGroupMember() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestCreateNextHop(t *testing.T) {
	tests := []struct {
		desc     string
		req      *saipb.CreateNextHopRequest
		wantAttr *saipb.NextHopAttribute
		wantReq  *dpb.AddNextHopRequest
		wantErr  string
	}{{
		desc:    "unknown type",
		req:     &saipb.CreateNextHopRequest{},
		wantErr: "InvalidArgument",
	}, {
		desc: "success",
		req: &saipb.CreateNextHopRequest{
			Type:              saipb.NextHopType_NEXT_HOP_TYPE_IP.Enum(),
			RouterInterfaceId: proto.Uint64(10),
			Ip:                []byte{127, 0, 0, 1},
		},
		wantAttr: &saipb.NextHopAttribute{
			Type:              saipb.NextHopType_NEXT_HOP_TYPE_IP.Enum(),
			RouterInterfaceId: proto.Uint64(10),
			Ip:                []byte{127, 0, 0, 1},
		},
		wantReq: &dpb.AddNextHopRequest{
			Id: 1,
			NextHop: &dpb.NextHop{
				Dev: &dpb.NextHop_Interface{
					Interface: "10",
				},
				Ip: &dpb.NextHop_IpBytes{
					IpBytes: []byte{127, 0, 0, 1},
				},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeRoutingDataplaneAPI{}
			c, mgr, stopFn := newTestNextHop(t, dplane)
			defer stopFn()
			_, gotErr := c.CreateNextHop(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateNextHop() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotAddNextHopReq[0], tt.wantReq, protocmp.Transform()); d != "" {
				t.Errorf("CreateNextHop() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.NextHopAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateNextHop() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func newTestNeighbor(t testing.TB, api routingDataplaneAPI) (saipb.NeighborClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newNeighbor(mgr, api, srv)
	})
	return saipb.NewNeighborClient(conn), mgr, stopFn
}

func newTestNextHopGroup(t testing.TB, api routingDataplaneAPI) (saipb.NextHopGroupClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newNextHopGroup(mgr, api, srv)
	})
	return saipb.NewNextHopGroupClient(conn), mgr, stopFn
}

func newTestNextHop(t testing.TB, api routingDataplaneAPI) (saipb.NextHopClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newNextHop(mgr, api, srv)
	})
	return saipb.NewNextHopClient(conn), mgr, stopFn
}
