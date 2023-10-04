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

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"
	dpb "github.com/openconfig/lemming/proto/dataplane"
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
			OperStatus:                       saipb.PortOperStatus_PORT_OPER_STATUS_UP.Enum(),
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
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			getInterface = func(name string) (*net.Interface, error) {
				return nil, tt.getInterfaceErr
			}
			dplane := &fakePortDataplaneAPI{}
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
			ContextId: &fwdpb.ContextId{},
		},
		wantAttr: &saipb.PortAttribute{
			OperStatus: saipb.PortOperStatus_PORT_OPER_STATUS_DOWN.Enum(),
			AdminState: proto.Bool(false),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			getInterface = func(name string) (*net.Interface, error) {
				return nil, tt.getInterfaceErr
			}
			dplane := &fakePortDataplaneAPI{}
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
			dplane := &fakePortDataplaneAPI{}
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
			ContextId: &fwdpb.ContextId{},
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
			dplane := &fakePortDataplaneAPI{}
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

func newTestPort(t testing.TB, api portDataplaneAPI) (saipb.PortClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newPort(mgr, api, srv)
	})
	return saipb.NewPortClient(conn), mgr, stopFn
}

func newTestHostif(t testing.TB, api portDataplaneAPI) (saipb.HostifClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newHostif(mgr, api, srv)
	})
	return saipb.NewHostifClient(conn), mgr, stopFn
}

type fakePortDataplaneAPI struct {
	gotPortStateReq  []*fwdpb.PortStateRequest
	gotPortCreateReq []*dpb.CreatePortRequest
}

func (f *fakePortDataplaneAPI) ID() string {
	return ""
}

func (f *fakePortDataplaneAPI) CreatePort(_ context.Context, req *dpb.CreatePortRequest) (*dpb.CreatePortResponse, error) {
	f.gotPortCreateReq = append(f.gotPortCreateReq, req)
	return nil, nil
}

func (f *fakePortDataplaneAPI) PortState(_ context.Context, req *fwdpb.PortStateRequest) (*fwdpb.PortStateReply, error) {
	f.gotPortStateReq = append(f.gotPortStateReq, req)
	return nil, nil
}
