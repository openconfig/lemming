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
	"io"
	"log"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

func TestCreateSwitch(t *testing.T) {
	c, mgr, stopFn := newTestSwitch(t, &fakeSwitchDataplane{})
	defer stopFn()
	got, err := c.CreateSwitch(context.Background(), &saipb.CreateSwitchRequest{})
	if err != nil {
		t.Fatalf("CreateSwitch() unexpected error: %v", err)
	}
	want := &saipb.CreateSwitchResponse{
		Oid: 1,
	}
	if d := cmp.Diff(got, want, protocmp.Transform()); d != "" {
		t.Fatalf("CreateSwitch() failed: diff(-got,+want)\n:%s", d)
	}
	wantAttr := &saipb.SwitchAttribute{
		CpuPort:                          proto.Uint64(2),
		NumberOfActivePorts:              proto.Uint32(0),
		AclEntryMinimumPriority:          proto.Uint32(1),
		AclEntryMaximumPriority:          proto.Uint32(10000),
		AclTableMinimumPriority:          proto.Uint32(1),
		AclTableMaximumPriority:          proto.Uint32(10000),
		MaxAclActionCount:                proto.Uint32(1000),
		NumberOfEcmpGroups:               proto.Uint32(1024),
		PortList:                         []uint64{2},
		SwitchHardwareInfo:               []int32{},
		DefaultStpInstId:                 proto.Uint64(3),
		DefaultVlanId:                    proto.Uint64(4),
		DefaultVirtualRouterId:           proto.Uint64(5),
		DefaultOverrideVirtualRouterId:   proto.Uint64(5),
		Default_1QBridgeId:               proto.Uint64(6),
		DefaultTrapGroup:                 proto.Uint64(7),
		IngressAcl:                       proto.Uint64(0),
		EgressAcl:                        proto.Uint64(0),
		PreIngressAcl:                    proto.Uint64(0),
		AvailableIpv4RouteEntry:          proto.Uint32(1024),
		AvailableIpv6RouteEntry:          proto.Uint32(1024),
		AvailableIpv4NexthopEntry:        proto.Uint32(1024),
		AvailableIpv6NexthopEntry:        proto.Uint32(1024),
		AvailableIpv4NeighborEntry:       proto.Uint32(1024),
		AvailableIpv6NeighborEntry:       proto.Uint32(1024),
		AvailableNextHopGroupEntry:       proto.Uint32(1024),
		AvailableNextHopGroupMemberEntry: proto.Uint32(1024),
		AvailableFdbEntry:                proto.Uint32(1024),
		AvailableL2McEntry:               proto.Uint32(1024),
		AvailableIpmcEntry:               proto.Uint32(1024),
		AvailableSnatEntry:               proto.Uint32(1024),
		AvailableDnatEntry:               proto.Uint32(1024),
		MaxAclRangeCount:                 proto.Uint32(10),
		AclStageIngress: &saipb.ACLCapability{
			IsActionListMandatory: false,
			ActionList:            []saipb.AclActionType{saipb.AclActionType_ACL_ACTION_TYPE_PACKET_ACTION, saipb.AclActionType_ACL_ACTION_TYPE_MIRROR_INGRESS, saipb.AclActionType_ACL_ACTION_TYPE_NO_NAT},
		},
		AclStageEgress: &saipb.ACLCapability{
			IsActionListMandatory: false,
			ActionList:            []saipb.AclActionType{saipb.AclActionType_ACL_ACTION_TYPE_PACKET_ACTION},
		},
		EcmpHash:                       proto.Uint64(8),
		LagHash:                        proto.Uint64(8),
		EcmpHashIpv4:                   proto.Uint64(8),
		EcmpHashIpv4InIpv4:             proto.Uint64(8),
		EcmpHashIpv6:                   proto.Uint64(8),
		LagHashIpv4:                    proto.Uint64(8),
		LagHashIpv4InIpv4:              proto.Uint64(8),
		LagHashIpv6:                    proto.Uint64(8),
		RestartWarm:                    proto.Bool(false),
		WarmRecover:                    proto.Bool(false),
		LagDefaultHashAlgorithm:        saipb.HashAlgorithm_HASH_ALGORITHM_CRC.Enum(),
		LagDefaultHashSeed:             proto.Uint32(0),
		LagDefaultSymmetricHash:        proto.Bool(false),
		QosDefaultTc:                   proto.Uint32(0),
		QosDot1PToTcMap:                proto.Uint64(0),
		QosDot1PToColorMap:             proto.Uint64(0),
		QosTcToQueueMap:                proto.Uint64(0),
		QosTcAndColorToDot1PMap:        proto.Uint64(0),
		QosTcAndColorToDscpMap:         proto.Uint64(0),
		QosTcAndColorToMplsExpMap:      proto.Uint64(0),
		QosDscpToTcMap:                 proto.Uint64(0),
		QosDscpToColorMap:              proto.Uint64(0),
		QosMplsExpToTcMap:              proto.Uint64(0),
		QosMplsExpToColorMap:           proto.Uint64(0),
		QosDscpToForwardingClassMap:    proto.Uint64(0),
		QosMplsExpToForwardingClassMap: proto.Uint64(0),
		IpsecObjectId:                  proto.Uint64(0),
		TamObjectId:                    []uint64{},
		PortConnectorList:              []uint64{},
		MacsecObjectList:               []uint64{},
		SystemPortList:                 []uint64{},
		FabricPortList:                 []uint64{},
		TunnelObjectsList:              []uint64{},
		MyMacList:                      []uint64{},
		Type:                           saipb.SwitchType_SWITCH_TYPE_NPU.Enum(),
		NumberOfSystemPorts:            proto.Uint32(0),
		SwitchShellEnable:              proto.Bool(false),
		SwitchProfileId:                proto.Uint32(0),
		NatZoneCounterObjectId:         proto.Uint64(0),
	}
	attr := &saipb.SwitchAttribute{}
	if err := mgr.PopulateAllAttributes("1", attr); err != nil {
		t.Fatal(err)
	}
	if d := cmp.Diff(attr, wantAttr, protocmp.Transform()); d != "" {
		t.Fatalf("CreateSwitch() failed: diff(-got,+want)\n:%s", d)
	}
}

func TestSwitchPortStateChangeNotification(t *testing.T) {
	tests := []struct {
		desc    string
		want    []*saipb.PortStateChangeNotificationResponse
		notifs  []*fwdpb.EventDesc
		wantErr string
	}{{
		desc: "port state up",
		notifs: []*fwdpb.EventDesc{{
			Event: fwdpb.Event_EVENT_PORT,
			Desc: &fwdpb.EventDesc_Port{
				Port: &fwdpb.PortEventDesc{
					PortId:   &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
					PortInfo: &fwdpb.PortInfo{OperStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP},
				},
			},
		}},
		want: []*saipb.PortStateChangeNotificationResponse{{
			Data: []*saipb.PortOperStatusNotification{{
				PortId:    1,
				PortState: saipb.PortOperStatus_PORT_OPER_STATUS_UP,
			}},
		}},
	}, {
		desc: "port state down",
		notifs: []*fwdpb.EventDesc{{
			Event: fwdpb.Event_EVENT_PORT,
			Desc: &fwdpb.EventDesc_Port{
				Port: &fwdpb.PortEventDesc{
					PortId:   &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
					PortInfo: &fwdpb.PortInfo{OperStatus: fwdpb.PortState_PORT_STATE_DISABLED_DOWN},
				},
			},
		}},
		want: []*saipb.PortStateChangeNotificationResponse{{
			Data: []*saipb.PortOperStatusNotification{{
				PortId:    1,
				PortState: saipb.PortOperStatus_PORT_OPER_STATUS_DOWN,
			}},
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{
				events: tt.notifs,
			}
			c, mgr, stopFn := newTestSwitch(t, dplane)
			mgr.SetType("1", saipb.ObjectType_OBJECT_TYPE_PORT)
			defer stopFn()
			notifs, gotErr := c.PortStateChangeNotification(context.TODO(), &saipb.PortStateChangeNotificationRequest{})
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateNeighborEntry() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			got := []*saipb.PortStateChangeNotificationResponse{}
			for {
				r, err := notifs.Recv()
				if err != nil {
					break
				}
				got = append(got, r)
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateNeighborEntry() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

type fakeSwitchDataplane struct {
	events                   []*fwdpb.EventDesc
	gotEntryAddReqs          []*fwdpb.TableEntryAddRequest
	gotPortStateReq          []*fwdpb.PortStateRequest
	counterReplies           []*fwdpb.ObjectCountersReply
	gotPortCreateReqs        []*fwdpb.PortCreateRequest
	gotPortUpdateReqs        []*fwdpb.PortUpdateRequest
	gotObjectDeleteReqs      []*fwdpb.ObjectDeleteRequest
	gotFlowCounterCreateReqs []*fwdpb.FlowCounterCreateRequest
	gotFlowCounterQueryReqs  []*fwdpb.FlowCounterQueryRequest
	gotEntryRemoveReqs       []*fwdpb.TableEntryRemoveRequest
	gotPackets               [][]byte
	portIDToNID              map[string]uint64
	counterRepliesIdx        int
	flowQueryReplies         []*fwdpb.FlowCounterQueryReply
	flowQueryRepliesIdx      int
	ctx                      *fwdcontext.Context
}

func (f *fakeSwitchDataplane) NotifySubscribe(_ *fwdpb.NotifySubscribeRequest, srv fwdpb.Forwarding_NotifySubscribeServer) error {
	for _, e := range f.events {
		srv.Send(e)
	}
	return io.EOF
}

func (f *fakeSwitchDataplane) TableCreate(context.Context, *fwdpb.TableCreateRequest) (*fwdpb.TableCreateReply, error) {
	return nil, nil
}

func (f *fakeSwitchDataplane) TableEntryAdd(_ context.Context, req *fwdpb.TableEntryAddRequest) (*fwdpb.TableEntryAddReply, error) {
	f.gotEntryAddReqs = append(f.gotEntryAddReqs, req)
	return nil, nil
}

func (f *fakeSwitchDataplane) TableEntryRemove(_ context.Context, req *fwdpb.TableEntryRemoveRequest) (*fwdpb.TableEntryRemoveReply, error) {
	f.gotEntryRemoveReqs = append(f.gotEntryRemoveReqs, req)
	return nil, nil
}

func (f *fakeSwitchDataplane) PortIDToNID(id string) (uint64, bool) {
	nid, ok := f.portIDToNID[id]
	return nid, ok
}

func (f *fakeSwitchDataplane) PortState(_ context.Context, req *fwdpb.PortStateRequest) (*fwdpb.PortStateReply, error) {
	f.gotPortStateReq = append(f.gotPortStateReq, req)
	return nil, nil
}

func (f *fakeSwitchDataplane) ObjectCounters(context.Context, *fwdpb.ObjectCountersRequest) (*fwdpb.ObjectCountersReply, error) {
	if f.counterRepliesIdx > len(f.counterReplies) {
		return nil, io.EOF
	}
	r := f.counterReplies[f.counterRepliesIdx]
	f.counterRepliesIdx++
	return r, nil
}

func (f *fakeSwitchDataplane) ID() string {
	return "foo"
}

func (f *fakeSwitchDataplane) FindContext(*fwdpb.ContextId) (*fwdcontext.Context, error) {
	return f.ctx, nil
}

func (f *fakeSwitchDataplane) PortCreate(_ context.Context, req *fwdpb.PortCreateRequest) (*fwdpb.PortCreateReply, error) {
	f.gotPortCreateReqs = append(f.gotPortCreateReqs, req)
	return nil, nil
}

func (f *fakeSwitchDataplane) PortUpdate(_ context.Context, req *fwdpb.PortUpdateRequest) (*fwdpb.PortUpdateReply, error) {
	f.gotPortUpdateReqs = append(f.gotPortUpdateReqs, req)
	return nil, nil
}

func (f *fakeSwitchDataplane) AttributeUpdate(context.Context, *fwdpb.AttributeUpdateRequest) (*fwdpb.AttributeUpdateReply, error) {
	return nil, nil
}

func (f *fakeSwitchDataplane) ObjectNID(context.Context, *fwdpb.ObjectNIDRequest) (*fwdpb.ObjectNIDReply, error) {
	return nil, nil
}

func (f *fakeSwitchDataplane) InjectPacket(_ *fwdpb.ContextId, _ *fwdpb.PortId, _ fwdpb.PacketHeaderId, pkt []byte, _ []*fwdpb.ActionDesc, _ bool, _ fwdpb.PortAction) error {
	f.gotPackets = append(f.gotPackets, pkt)
	return nil
}

func (f *fakeSwitchDataplane) ObjectDelete(_ context.Context, req *fwdpb.ObjectDeleteRequest) (*fwdpb.ObjectDeleteReply, error) {
	f.gotObjectDeleteReqs = append(f.gotObjectDeleteReqs, req)
	return nil, nil
}

func (f *fakeSwitchDataplane) FlowCounterCreate(_ context.Context, req *fwdpb.FlowCounterCreateRequest) (*fwdpb.FlowCounterCreateReply, error) {
	f.gotFlowCounterCreateReqs = append(f.gotFlowCounterCreateReqs, req)
	return nil, nil
}

func (f *fakeSwitchDataplane) FlowCounterQuery(_ context.Context, req *fwdpb.FlowCounterQueryRequest) (*fwdpb.FlowCounterQueryReply, error) {
	f.gotFlowCounterQueryReqs = append(f.gotFlowCounterQueryReqs, req)
	if f.flowQueryRepliesIdx > len(f.flowQueryReplies) {
		return nil, io.EOF
	}
	r := f.flowQueryReplies[f.flowQueryRepliesIdx]
	f.flowQueryRepliesIdx++
	return r, nil
}

func newTestServer(t testing.TB, newSrvFn func(mgr *attrmgr.AttrMgr, srv *grpc.Server)) (grpc.ClientConnInterface, *attrmgr.AttrMgr, func()) {
	t.Helper()
	mgr := attrmgr.New()
	lis, err := net.Listen("tcp", ("127.0.0.1:0"))
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()), grpc.ChainUnaryInterceptor(mgr.Interceptor))
	if newSrvFn != nil {
		newSrvFn(mgr, srv)
	}
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve forwarding server: %v", err)
		}
	}()
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	return conn, mgr, srv.Stop
}

func newTestSwitch(t testing.TB, dplane switchDataplaneAPI) (saipb.SwitchClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newSwitch(mgr, dplane, srv, &dplaneopts.Options{})
	})
	return saipb.NewSwitchClient(conn), mgr, stopFn
}
