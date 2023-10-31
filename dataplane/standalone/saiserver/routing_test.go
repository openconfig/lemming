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
			dplane := &fakeSwitchDataplane{}
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
			dplane := &fakeSwitchDataplane{}
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
			dplane := &fakeSwitchDataplane{}
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
			dplane := &fakeSwitchDataplane{}
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

func TestCreateRouteEntry(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.CreateRouteEntryRequest
		wantReq *dpb.AddIPRouteRequest
		types   map[string]saipb.ObjectType
		wantErr string
	}{{
		desc:    "unknown action",
		req:     &saipb.CreateRouteEntryRequest{},
		wantErr: "InvalidArgument",
	}, {
		desc: "drop action",
		req: &saipb.CreateRouteEntryRequest{
			Entry: &saipb.RouteEntry{
				Destination: &saipb.IpPrefix{
					Addr: []byte{127, 0, 0, 1},
					Mask: []byte{255, 255, 255, 255},
				},
			},
			PacketAction: saipb.PacketAction_PACKET_ACTION_DROP.Enum(),
		},
		wantReq: &dpb.AddIPRouteRequest{
			Route: &dpb.Route{
				Prefix: &dpb.RoutePrefix{
					Prefix: &dpb.RoutePrefix_Mask{
						Mask: &dpb.IpMask{
							Addr: []byte{127, 0, 0, 1},
							Mask: []byte{255, 255, 255, 255},
						},
					},
				},
				Action: dpb.PacketAction_PACKET_ACTION_DROP,
			},
		},
	}, {
		desc:  "forward port action",
		types: map[string]saipb.ObjectType{"100": saipb.ObjectType_OBJECT_TYPE_PORT},
		req: &saipb.CreateRouteEntryRequest{
			Entry: &saipb.RouteEntry{
				Destination: &saipb.IpPrefix{
					Addr: []byte{127, 0, 0, 1},
					Mask: []byte{255, 255, 255, 255},
				},
			},
			PacketAction: saipb.PacketAction_PACKET_ACTION_TRANSIT.Enum(),
			NextHopId:    proto.Uint64(100),
		},
		wantReq: &dpb.AddIPRouteRequest{
			Route: &dpb.Route{
				Prefix: &dpb.RoutePrefix{
					Prefix: &dpb.RoutePrefix_Mask{
						Mask: &dpb.IpMask{
							Addr: []byte{127, 0, 0, 1},
							Mask: []byte{255, 255, 255, 255},
						},
					},
				},
				Hop:    &dpb.Route_PortId{PortId: "100"},
				Action: dpb.PacketAction_PACKET_ACTION_FORWARD,
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, mgr, stopFn := newTestRoute(t, dplane)
			defer stopFn()
			for k, v := range tt.types {
				mgr.SetType(k, v)
			}
			_, gotErr := c.CreateRouteEntry(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateRouteEntry() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotAddIPRouteReq[0], tt.wantReq, protocmp.Transform()); d != "" {
				t.Errorf("CreateRouteEntry() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestCreateRouterInterface(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.CreateRouterInterfaceRequest
		wantReq *dpb.AddInterfaceRequest
		wantErr string
	}{{
		desc:    "unknown type",
		req:     &saipb.CreateRouterInterfaceRequest{},
		wantErr: "InvalidArgument",
	}, {
		desc: "success port",
		req: &saipb.CreateRouterInterfaceRequest{
			PortId: proto.Uint64(10),
			Type:   saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT.Enum(),
		},
		wantReq: &dpb.AddInterfaceRequest{
			Id:      "1",
			PortIds: []string{"10"},
			Type:    dpb.InterfaceType_INTERFACE_TYPE_PORT,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, _, stopFn := newTestRouterInterface(t, dplane)
			defer stopFn()
			_, gotErr := c.CreateRouterInterface(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateRouterInterface() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotAddInterfaceReq[0], tt.wantReq, protocmp.Transform()); d != "" {
				t.Errorf("CreateRouterInterface() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func newTestNeighbor(t testing.TB, api switchDataplaneAPI) (saipb.NeighborClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newNeighbor(mgr, api, srv)
	})
	return saipb.NewNeighborClient(conn), mgr, stopFn
}

func newTestNextHopGroup(t testing.TB, api switchDataplaneAPI) (saipb.NextHopGroupClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newNextHopGroup(mgr, api, srv)
	})
	return saipb.NewNextHopGroupClient(conn), mgr, stopFn
}

func newTestNextHop(t testing.TB, api switchDataplaneAPI) (saipb.NextHopClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newNextHop(mgr, api, srv)
	})
	return saipb.NewNextHopClient(conn), mgr, stopFn
}

func newTestRoute(t testing.TB, api switchDataplaneAPI) (saipb.RouteClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newRoute(mgr, api, srv)
	})
	return saipb.NewRouteClient(conn), mgr, stopFn
}

func newTestRouterInterface(t testing.TB, api switchDataplaneAPI) (saipb.RouterInterfaceClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newRouterInterface(mgr, api, srv)
	})
	return saipb.NewRouterInterfaceClient(conn), mgr, stopFn
}
