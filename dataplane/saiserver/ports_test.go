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
	"fmt"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestCreatePort(t *testing.T) {
	tests := []struct {
		desc            string
		req             *saipb.CreatePortRequest
		getInterfaceErr error
		want            *saipb.CreatePortResponse
		wantAttr        *saipb.PortAttribute
		wantErr         string
	}{{
		desc:            "non-existent interface",
		req:             &saipb.CreatePortRequest{},
		getInterfaceErr: fmt.Errorf("no interface"),
		want: &saipb.CreatePortResponse{
			Oid: 1,
		},
		wantAttr: &saipb.PortAttribute{
			OperStatus:                       saipb.PortOperStatus_PORT_OPER_STATUS_NOT_PRESENT.Enum(),
			QosNumberOfQueues:                proto.Uint32(0),
			QosQueueList:                     []uint64{},
			QosNumberOfSchedulerGroups:       proto.Uint32(0),
			QosSchedulerGroupList:            []uint64{},
			IngressPriorityGroupList:         []uint64{},
			FloodStormControlPolicerId:       proto.Uint64(0),
			BroadcastStormControlPolicerId:   proto.Uint64(0),
			MulticastStormControlPolicerId:   proto.Uint64(0),
			IngressAcl:                       proto.Uint64(0),
			EgressAcl:                        proto.Uint64(0),
			IngressMacsecAcl:                 proto.Uint64(0),
			EgressMacsecAcl:                  proto.Uint64(0),
			MacsecPortList:                   []uint64{},
			IngressMirrorSession:             []uint64{},
			EgressMirrorSession:              []uint64{},
			IngressSamplepacketEnable:        proto.Uint64(0),
			EgressSamplepacketEnable:         proto.Uint64(0),
			IngressSampleMirrorSession:       []uint64{},
			EgressSampleMirrorSession:        []uint64{},
			PolicerId:                        proto.Uint64(0),
			QosDot1PToTcMap:                  proto.Uint64(0),
			QosDot1PToColorMap:               proto.Uint64(0),
			QosDscpToTcMap:                   proto.Uint64(0),
			QosDscpToColorMap:                proto.Uint64(0),
			QosTcToQueueMap:                  proto.Uint64(0),
			QosTcAndColorToDot1PMap:          proto.Uint64(0),
			QosTcAndColorToDscpMap:           proto.Uint64(0),
			QosTcToPriorityGroupMap:          proto.Uint64(0),
			QosPfcPriorityToPriorityGroupMap: proto.Uint64(0),
			QosPfcPriorityToQueueMap:         proto.Uint64(0),
			QosSchedulerProfileId:            proto.Uint64(0),
			QosIngressBufferProfileList:      []uint64{},
			QosEgressBufferProfileList:       []uint64{},
			EgressBlockPortList:              []uint64{},
			PortPoolList:                     []uint64{},
			IsolationGroup:                   proto.Uint64(0),
			TamObject:                        []uint64{},
			PortSerdesId:                     proto.Uint64(0),
			QosMplsExpToTcMap:                proto.Uint64(0),
			QosMplsExpToColorMap:             proto.Uint64(0),
			QosTcAndColorToMplsExpMap:        proto.Uint64(0),
			SystemPort:                       proto.Uint64(0),
			QosDscpToForwardingClassMap:      proto.Uint64(0),
			QosMplsExpToForwardingClassMap:   proto.Uint64(0),
			IpsecPort:                        proto.Uint64(0),
			SupportedSpeed:                   []uint32{1000, 10000, 40000},
			OperSpeed:                        proto.Uint32(40000),
			SupportedFecMode:                 []saipb.PortFecMode{saipb.PortFecMode_PORT_FEC_MODE_NONE},
			NumberOfIngressPriorityGroups:    proto.Uint32(0),
			QosMaximumHeadroomSize:           proto.Uint32(0),
			AdminState:                       proto.Bool(true),
			AutoNegMode:                      proto.Bool(true),
			Mtu:                              proto.Uint32(1514),
		},
	}, {
		desc: "existing interface",
		req:  &saipb.CreatePortRequest{},
		want: &saipb.CreatePortResponse{
			Oid: 1,
		},
		wantAttr: &saipb.PortAttribute{
			OperStatus:                       saipb.PortOperStatus_PORT_OPER_STATUS_DOWN.Enum(),
			QosNumberOfQueues:                proto.Uint32(0),
			QosQueueList:                     []uint64{},
			QosNumberOfSchedulerGroups:       proto.Uint32(0),
			QosSchedulerGroupList:            []uint64{},
			IngressPriorityGroupList:         []uint64{},
			FloodStormControlPolicerId:       proto.Uint64(0),
			BroadcastStormControlPolicerId:   proto.Uint64(0),
			MulticastStormControlPolicerId:   proto.Uint64(0),
			IngressAcl:                       proto.Uint64(0),
			EgressAcl:                        proto.Uint64(0),
			IngressMacsecAcl:                 proto.Uint64(0),
			EgressMacsecAcl:                  proto.Uint64(0),
			MacsecPortList:                   []uint64{},
			IngressMirrorSession:             []uint64{},
			EgressMirrorSession:              []uint64{},
			IngressSamplepacketEnable:        proto.Uint64(0),
			EgressSamplepacketEnable:         proto.Uint64(0),
			IngressSampleMirrorSession:       []uint64{},
			EgressSampleMirrorSession:        []uint64{},
			PolicerId:                        proto.Uint64(0),
			QosDot1PToTcMap:                  proto.Uint64(0),
			QosDot1PToColorMap:               proto.Uint64(0),
			QosDscpToTcMap:                   proto.Uint64(0),
			QosDscpToColorMap:                proto.Uint64(0),
			QosTcToQueueMap:                  proto.Uint64(0),
			QosTcAndColorToDot1PMap:          proto.Uint64(0),
			QosTcAndColorToDscpMap:           proto.Uint64(0),
			QosTcToPriorityGroupMap:          proto.Uint64(0),
			QosPfcPriorityToPriorityGroupMap: proto.Uint64(0),
			QosPfcPriorityToQueueMap:         proto.Uint64(0),
			QosSchedulerProfileId:            proto.Uint64(0),
			QosIngressBufferProfileList:      []uint64{},
			QosEgressBufferProfileList:       []uint64{},
			EgressBlockPortList:              []uint64{},
			PortPoolList:                     []uint64{},
			IsolationGroup:                   proto.Uint64(0),
			TamObject:                        []uint64{},
			PortSerdesId:                     proto.Uint64(0),
			QosMplsExpToTcMap:                proto.Uint64(0),
			QosMplsExpToColorMap:             proto.Uint64(0),
			QosTcAndColorToMplsExpMap:        proto.Uint64(0),
			SystemPort:                       proto.Uint64(0),
			QosDscpToForwardingClassMap:      proto.Uint64(0),
			QosMplsExpToForwardingClassMap:   proto.Uint64(0),
			IpsecPort:                        proto.Uint64(0),
			SupportedSpeed:                   []uint32{1000, 10000, 40000},
			OperSpeed:                        proto.Uint32(40000),
			SupportedFecMode:                 []saipb.PortFecMode{saipb.PortFecMode_PORT_FEC_MODE_NONE},
			NumberOfIngressPriorityGroups:    proto.Uint32(0),
			QosMaximumHeadroomSize:           proto.Uint32(0),
			AdminState:                       proto.Bool(false),
			AutoNegMode:                      proto.Bool(true),
			Mtu:                              proto.Uint32(1514),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			getInterface = func(string) (*net.Interface, error) {
				return nil, tt.getInterfaceErr
			}
			dplane := &fakeSwitchDataplane{}
			c, mgr, stopFn := newTestPort(t, dplane)
			defer stopFn()
			got, gotErr := c.CreatePort(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreatePort() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreatePort() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.PortAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreatePort() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestCreatePorts(t *testing.T) {
	tests := []struct {
		desc            string
		req             *saipb.CreatePortsRequest
		getInterfaceErr error
		want            *saipb.CreatePortsResponse
		wantAttr        *saipb.PortAttribute
		wantErr         string
	}{{
		desc: "success",
		req: &saipb.CreatePortsRequest{
			Reqs: []*saipb.CreatePortRequest{{}},
		},
		want: &saipb.CreatePortsResponse{
			Resps: []*saipb.CreatePortResponse{{Oid: 1}},
		},
		wantAttr: &saipb.PortAttribute{
			OperStatus:                       saipb.PortOperStatus_PORT_OPER_STATUS_DOWN.Enum(),
			QosNumberOfQueues:                proto.Uint32(0),
			QosQueueList:                     []uint64{},
			QosNumberOfSchedulerGroups:       proto.Uint32(0),
			QosSchedulerGroupList:            []uint64{},
			IngressPriorityGroupList:         []uint64{},
			FloodStormControlPolicerId:       proto.Uint64(0),
			BroadcastStormControlPolicerId:   proto.Uint64(0),
			MulticastStormControlPolicerId:   proto.Uint64(0),
			IngressAcl:                       proto.Uint64(0),
			EgressAcl:                        proto.Uint64(0),
			IngressMacsecAcl:                 proto.Uint64(0),
			EgressMacsecAcl:                  proto.Uint64(0),
			MacsecPortList:                   []uint64{},
			IngressMirrorSession:             []uint64{},
			EgressMirrorSession:              []uint64{},
			IngressSamplepacketEnable:        proto.Uint64(0),
			EgressSamplepacketEnable:         proto.Uint64(0),
			IngressSampleMirrorSession:       []uint64{},
			EgressSampleMirrorSession:        []uint64{},
			PolicerId:                        proto.Uint64(0),
			QosDot1PToTcMap:                  proto.Uint64(0),
			QosDot1PToColorMap:               proto.Uint64(0),
			QosDscpToTcMap:                   proto.Uint64(0),
			QosDscpToColorMap:                proto.Uint64(0),
			QosTcToQueueMap:                  proto.Uint64(0),
			QosTcAndColorToDot1PMap:          proto.Uint64(0),
			QosTcAndColorToDscpMap:           proto.Uint64(0),
			QosTcToPriorityGroupMap:          proto.Uint64(0),
			QosPfcPriorityToPriorityGroupMap: proto.Uint64(0),
			QosPfcPriorityToQueueMap:         proto.Uint64(0),
			QosSchedulerProfileId:            proto.Uint64(0),
			QosIngressBufferProfileList:      []uint64{},
			QosEgressBufferProfileList:       []uint64{},
			EgressBlockPortList:              []uint64{},
			PortPoolList:                     []uint64{},
			IsolationGroup:                   proto.Uint64(0),
			TamObject:                        []uint64{},
			PortSerdesId:                     proto.Uint64(0),
			QosMplsExpToTcMap:                proto.Uint64(0),
			QosMplsExpToColorMap:             proto.Uint64(0),
			QosTcAndColorToMplsExpMap:        proto.Uint64(0),
			SystemPort:                       proto.Uint64(0),
			QosDscpToForwardingClassMap:      proto.Uint64(0),
			QosMplsExpToForwardingClassMap:   proto.Uint64(0),
			IpsecPort:                        proto.Uint64(0),
			SupportedSpeed:                   []uint32{1000, 10000, 40000},
			OperSpeed:                        proto.Uint32(40000),
			SupportedFecMode:                 []saipb.PortFecMode{saipb.PortFecMode_PORT_FEC_MODE_NONE},
			NumberOfIngressPriorityGroups:    proto.Uint32(0),
			QosMaximumHeadroomSize:           proto.Uint32(0),
			AdminState:                       proto.Bool(false),
			AutoNegMode:                      proto.Bool(true),
			Mtu:                              proto.Uint32(1514),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, mgr, stopFn := newTestPort(t, dplane)
			defer stopFn()
			got, gotErr := c.CreatePorts(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreatePort() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreatePorts() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.PortAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreatePorts() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestSetPortAttribute(t *testing.T) {
	tests := []struct {
		desc            string
		req             *saipb.SetPortAttributeRequest
		getInterfaceErr error
		wantAttr        *saipb.PortAttribute
		wantReq         *fwdpb.PortStateRequest
		wantErr         string
	}{{
		desc: "admin status",
		req: &saipb.SetPortAttributeRequest{
			Oid:        1,
			AdminState: proto.Bool(false),
		},
		wantReq: &fwdpb.PortStateRequest{
			PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
			Operation: &fwdpb.PortInfo{AdminStatus: fwdpb.PortState_PORT_STATE_DISABLED_DOWN},
			ContextId: &fwdpb.ContextId{Id: "foo"},
		},
		wantAttr: &saipb.PortAttribute{
			OperStatus: saipb.PortOperStatus_PORT_OPER_STATUS_DOWN.Enum(),
			AdminState: proto.Bool(false),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			getInterface = func(string) (*net.Interface, error) {
				return nil, tt.getInterfaceErr
			}
			dplane := &fakeSwitchDataplane{}
			c, mgr, stopFn := newTestPort(t, dplane)
			mgr.StoreAttributes(1, &saipb.PortAttribute{
				OperStatus: saipb.PortOperStatus_PORT_OPER_STATUS_DOWN.Enum(),
			})
			defer stopFn()
			_, gotErr := c.SetPortAttribute(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("SetPortAttribute() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotPortStateReq[0], tt.wantReq, protocmp.Transform()); d != "" {
				t.Errorf("SetPortAttribute() failed: diff(-got,+want)\n:%s", d)
			}
			attr := &saipb.PortAttribute{}
			if err := mgr.PopulateAllAttributes("1", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("SetPortAttribute() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestGetPortStats(t *testing.T) {
	tests := []struct {
		desc         string
		req          *saipb.GetPortStatsRequest
		counterReply *fwdpb.ObjectCountersReply
		want         *saipb.GetPortStatsResponse
		wantErr      string
	}{{
		desc: "all stats",
		req: &saipb.GetPortStatsRequest{
			Oid: 1,
			CounterIds: []saipb.PortStat{
				saipb.PortStat_PORT_STAT_IF_IN_UCAST_PKTS,
				saipb.PortStat_PORT_STAT_IF_IN_NON_UCAST_PKTS,
				saipb.PortStat_PORT_STAT_IF_IN_ERRORS,
				saipb.PortStat_PORT_STAT_IF_OUT_UCAST_PKTS,
				saipb.PortStat_PORT_STAT_IF_OUT_NON_UCAST_PKTS,
				saipb.PortStat_PORT_STAT_IF_OUT_ERRORS,
				saipb.PortStat_PORT_STAT_IF_IN_OCTETS,
				saipb.PortStat_PORT_STAT_IF_OUT_OCTETS,
				saipb.PortStat_PORT_STAT_IF_IN_DISCARDS,
				saipb.PortStat_PORT_STAT_IF_OUT_DISCARDS,
			},
		},
		counterReply: &fwdpb.ObjectCountersReply{
			Counters: []*fwdpb.Counter{{
				Id:    fwdpb.CounterId_COUNTER_ID_TX_DROP_PACKETS,
				Value: 1,
			}, {
				Id:    fwdpb.CounterId_COUNTER_ID_RX_DROP_PACKETS,
				Value: 2,
			}, {
				Id:    fwdpb.CounterId_COUNTER_ID_TX_OCTETS,
				Value: 3,
			}, {
				Id:    fwdpb.CounterId_COUNTER_ID_RX_OCTETS,
				Value: 4,
			}, {
				Id:    fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS,
				Value: 5,
			}, {
				Id:    fwdpb.CounterId_COUNTER_ID_TX_NON_UCAST_PACKETS,
				Value: 6,
			}, {
				Id:    fwdpb.CounterId_COUNTER_ID_TX_UCAST_PACKETS,
				Value: 7,
			}, {
				Id:    fwdpb.CounterId_COUNTER_ID_RX_ERROR_PACKETS,
				Value: 8,
			}, {
				Id:    fwdpb.CounterId_COUNTER_ID_RX_NON_UCAST_PACKETS,
				Value: 9,
			}, {
				Id:    fwdpb.CounterId_COUNTER_ID_RX_UCAST_PACKETS,
				Value: 10,
			}},
		},
		want: &saipb.GetPortStatsResponse{
			Values: []uint64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{
				counterReplies: []*fwdpb.ObjectCountersReply{tt.counterReply},
			}
			c, _, stopFn := newTestPort(t, dplane)
			defer stopFn()
			got, gotErr := c.GetPortStats(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("SetPortAttribute() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("SetPortAttribute() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func newTestPort(t testing.TB, api switchDataplaneAPI) (saipb.PortClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newPort(mgr, api, srv, &dplaneopts.Options{PortType: fwdpb.PortType_PORT_TYPE_KERNEL})
	})
	return saipb.NewPortClient(conn), mgr, stopFn
}

func newTestLAG(t testing.TB, api switchDataplaneAPI) (saipb.LagClient, func()) {
	conn, _, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newLAG(mgr, api, srv)
	})
	return saipb.NewLagClient(conn), stopFn
}

func TestCreateLag(t *testing.T) {
	tests := []struct {
		desc            string
		req             *saipb.CreateLagRequest
		getInterfaceErr error
		want            *fwdpb.PortCreateRequest
		wantErr         string
	}{{
		desc: "success",
		req:  &saipb.CreateLagRequest{},
		want: &fwdpb.PortCreateRequest{
			ContextId: &fwdpb.ContextId{Id: "foo"},
			Port: &fwdpb.PortDesc{
				PortType: fwdpb.PortType_PORT_TYPE_AGGREGATE_PORT,
				PortId:   &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, stopFn := newTestLAG(t, dplane)
			defer stopFn()
			_, gotErr := c.CreateLag(context.Background(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateLag() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(dplane.gotPortCreateReqs[0], tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateLag() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestLagMember(t *testing.T) {
	dplane := &fakeSwitchDataplane{}
	c, stopFn := newTestLAG(t, dplane)
	defer stopFn()

	createResp, err := c.CreateLagMember(context.Background(), &saipb.CreateLagMemberRequest{
		LagId:  proto.Uint64(1),
		PortId: proto.Uint64(2),
	})
	if err != nil {
		t.Fatalf("CreateLagMember() unexpected err: %v", err)
	}
	want := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: "foo"},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_AggregateAdd{
				AggregateAdd: &fwdpb.AggregatePortAddMemberUpdateDesc{
					InstanceCount: 1,
					PortId:        &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "2"}},
				},
			},
		},
	}
	if d := cmp.Diff(dplane.gotPortUpdateReqs[0], want, protocmp.Transform()); d != "" {
		t.Errorf("CreateLagMember() failed: diff(-got,+want)\n:%s", d)
	}

	_, err = c.RemoveLagMember(context.Background(), &saipb.RemoveLagMemberRequest{
		Oid: createResp.Oid,
	})
	if err != nil {
		t.Fatalf("RemoveLagMember() unexpected err: %v", err)
	}
	want = &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: "foo"},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "1"}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_AggregateDel{
				AggregateDel: &fwdpb.AggregatePortRemoveMemberUpdateDesc{
					PortId: &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "2"}},
				},
			},
		},
	}
	if d := cmp.Diff(dplane.gotPortUpdateReqs[1], want, protocmp.Transform()); d != "" {
		t.Errorf("RemoveLagMember() failed: diff(-got,+want)\n:%s", d)
	}
}
