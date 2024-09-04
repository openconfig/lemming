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
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestCreateHostif(t *testing.T) {
	tests := []struct {
		desc     string
		req      *saipb.CreateHostifRequest
		want     *saipb.CreateHostifResponse
		wantAttr *saipb.HostifAttribute
		wantErr  string
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
			Oid: 3,
		},
		wantAttr: &saipb.HostifAttribute{
			Type:       saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId:      proto.Uint64(2),
			OperStatus: proto.Bool(true),
		},
	}, {
		desc: "success cpu port",
		req: &saipb.CreateHostifRequest{
			Type:  saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId: proto.Uint64(10),
		},
		want: &saipb.CreateHostifResponse{
			Oid: 3,
		},
		wantAttr: &saipb.HostifAttribute{
			Type:       saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId:      proto.Uint64(10),
			OperStatus: proto.Bool(true),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{
				ctx: fwdcontext.New("foo", "foo"),
			}
			dplane.ctx.SetPacketSink(func(*fwdpb.PacketSinkResponse) error { return nil })
			c, mgr, stopFn := newTestHostif(t, dplane)
			// Create switch and ports
			mgr.StoreAttributes(mgr.NextID(), &saipb.SwitchAttribute{
				CpuPort: proto.Uint64(10),
			})
			mgr.StoreAttributes(10, &saipb.PortAttribute{
				OperStatus: saipb.PortOperStatus_PORT_OPER_STATUS_UP.Enum(),
			})
			mgr.StoreAttributes(mgr.NextID(), &saipb.PortAttribute{
				OperStatus: saipb.PortOperStatus_PORT_OPER_STATUS_DOWN.Enum(),
			})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			msgCh := make(chan *pktiopb.HostPortControlMessage, 1)
			pc, err := c.HostPortControl(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if err := pc.Send(&pktiopb.HostPortControlRequest{Msg: &pktiopb.HostPortControlRequest_Init{}}); err != nil {
				t.Fatal(err)
			}
			time.Sleep(time.Millisecond)
			go func() {
				msg, _ := pc.Recv()
				msgCh <- msg
				pc.Send(&pktiopb.HostPortControlRequest{
					Msg: &pktiopb.HostPortControlRequest_Status{
						Status: &status.Status{
							Code:    int32(codes.OK),
							Message: "",
						},
					},
				})
			}()

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
			_ = <-msgCh
			attr := &saipb.HostifAttribute{}
			if err := mgr.PopulateAllAttributes("3", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateHostif() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestRemoveHostif(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.RemoveHostifRequest
		want    *pktiopb.HostPortControlMessage
		wantErr string
	}{{
		desc: "sucess",
		req: &saipb.RemoveHostifRequest{
			Oid: 1,
		},
		want: &pktiopb.HostPortControlMessage{
			PortId: 1,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{
				portIDToNID: map[string]uint64{
					"1": 10,
				},
			}
			c, mgr, stopFn := newTestHostif(t, dplane)
			defer stopFn()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mgr.StoreAttributes(1, &saipb.CreateHostifRequest{Name: []byte("eth1")})

			msgCh := make(chan *pktiopb.HostPortControlMessage, 1)
			pc, err := c.HostPortControl(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if err := pc.Send(&pktiopb.HostPortControlRequest{Msg: &pktiopb.HostPortControlRequest_Init{}}); err != nil {
				t.Fatal(err)
			}
			time.Sleep(time.Millisecond)
			go func() {
				msg, _ := pc.Recv()
				msgCh <- msg
				pc.Send(&pktiopb.HostPortControlRequest{
					Msg: &pktiopb.HostPortControlRequest_Status{
						Status: &status.Status{
							Code:    int32(codes.OK),
							Message: "",
						},
					},
				})
			}()

			_, gotErr := c.RemoveHostif(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("RemoveHostif() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			got := <-msgCh
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("RemoveHostif() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

type hostifClient struct {
	saipb.HostifClient
	pktiopb.PacketIOClient
}

func newTestHostif(t testing.TB, api switchDataplaneAPI) (*hostifClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newHostif(mgr, api, srv, &dplaneopts.Options{
			HostifNetDevType: fwdpb.PortType_PORT_TYPE_KERNEL,
		})
	})
	return &hostifClient{
		HostifClient:   saipb.NewHostifClient(conn),
		PacketIOClient: pktiopb.NewPacketIOClient(conn),
	}, mgr, stopFn
}
