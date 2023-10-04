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

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type portDataplaneAPI interface {
	ID() string
	CreatePort(ctx context.Context, req *dpb.CreatePortRequest) (*dpb.CreatePortResponse, error)
	PortState(ctx context.Context, req *fwdpb.PortStateRequest) (*fwdpb.PortStateReply, error)
}

func newPort(mgr *attrmgr.AttrMgr, dataplane portDataplaneAPI, s *grpc.Server) *port {
	p := &port{
		mgr:       mgr,
		dataplane: dataplane,
		portToEth: make(map[uint64]string),
		nextEth:   1, // Start at eth1
	}
	saipb.RegisterPortServer(s, p)
	return p
}

type port struct {
	saipb.UnimplementedPortServer
	mgr       *attrmgr.AttrMgr
	dataplane portDataplaneAPI
	nextEth   int
	portToEth map[uint64]string
}

// stub for testing
var getInterface = net.InterfaceByName

// CreatePort creates a new port, mapping the port to ethX, where X is assigned sequentially from 1 to n.
// Note: If more ports are created than eth devices, no error is returned, but the OperStatus is set to NOT_PRESENT.
func (port *port) CreatePort(ctx context.Context, _ *saipb.CreatePortRequest) (*saipb.CreatePortResponse, error) {
	id := port.mgr.NextID()

	dev := fmt.Sprintf("eth%v", port.nextEth)
	port.nextEth++

	attrs := &saipb.PortAttribute{
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
	}

	// For ports that don't exist, do not create dataplane ports.
	if _, err := getInterface(dev); err != nil {
		attrs.OperStatus = saipb.PortOperStatus_PORT_OPER_STATUS_NOT_PRESENT.Enum()
		port.mgr.StoreAttributes(id, attrs)
		return &saipb.CreatePortResponse{
			Oid: id,
		}, nil
	}
	port.portToEth[id] = dev

	_, err := port.dataplane.CreatePort(ctx, &dpb.CreatePortRequest{
		Id:   fmt.Sprint(id),
		Type: fwdpb.PortType_PORT_TYPE_KERNEL,
		Src: &dpb.CreatePortRequest_KernelDev{
			KernelDev: dev,
		},
		Location: dpb.PortLocation_PORT_LOCATION_EXTERNAL,
	})
	if err != nil {
		return nil, err
	}
	attrs.OperStatus = saipb.PortOperStatus_PORT_OPER_STATUS_UP.Enum()
	port.mgr.StoreAttributes(id, attrs)

	return &saipb.CreatePortResponse{
		Oid: id,
	}, nil
}

func (port *port) createCPUPort(ctx context.Context) (uint64, error) {
	id := port.mgr.NextID()

	_, err := port.dataplane.CreatePort(ctx, &dpb.CreatePortRequest{
		Id:       fmt.Sprint(id),
		Type:     fwdpb.PortType_PORT_TYPE_CPU_PORT,
		Location: dpb.PortLocation_PORT_LOCATION_CPU,
	})
	if err != nil {
		return 0, err
	}

	cpuPort := &saipb.PortAttribute{
		Type:                             saipb.PortType_PORT_TYPE_CPU.Enum(),
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
		SupportedSpeed:                   []uint32{1024},
		OperSpeed:                        proto.Uint32(1024),
		SupportedFecMode:                 []saipb.PortFecMode{saipb.PortFecMode_PORT_FEC_MODE_NONE},
		NumberOfIngressPriorityGroups:    proto.Uint32(0),
		QosMaximumHeadroomSize:           proto.Uint32(0),
		AdminState:                       proto.Bool(true),
		AutoNegMode:                      proto.Bool(false),
		Mtu:                              proto.Uint32(1514),
	}
	port.mgr.SetType(fmt.Sprint(id), saipb.ObjectType_OBJECT_TYPE_PORT)
	port.mgr.StoreAttributes(id, cpuPort)

	return id, nil
}

// SetPortAttributes sets the attributes in the request.
func (port *port) SetPortAttribute(ctx context.Context, req *saipb.SetPortAttributeRequest) (*saipb.SetPortAttributeResponse, error) {
	if req.AdminState != nil {
		// Skip ports that don't exsit.
		attrReq := &saipb.GetPortAttributeRequest{Oid: req.GetOid(), AttrType: []saipb.PortAttr{saipb.PortAttr_PORT_ATTR_OPER_STATUS}}
		p := &saipb.GetPortAttributeResponse{}
		if err := port.mgr.PopulateAttributes(attrReq, p); err != nil {
			return nil, err
		}
		if p.GetAttr().GetOperStatus() == saipb.PortOperStatus_PORT_OPER_STATUS_NOT_PRESENT {
			return nil, nil
		}

		stateReq := &fwdpb.PortStateRequest{
			ContextId: &fwdpb.ContextId{Id: port.dataplane.ID()},
			PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(req.GetOid())}},
		}
		stateReq.Operation = &fwdpb.PortInfo{
			AdminStatus: fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
		}
		if req.GetAdminState() {
			stateReq.Operation.AdminStatus = fwdpb.PortState_PORT_STATE_ENABLED_UP
		}
		_, err := port.dataplane.PortState(ctx, stateReq)
		if err != nil {
			return nil, err
		}
	}
	return &saipb.SetPortAttributeResponse{}, nil
}

func (port *port) Reset() {
	port.portToEth = make(map[uint64]string)
	port.nextEth = 1
}

func newHostif(mgr *attrmgr.AttrMgr, dataplane portDataplaneAPI, s *grpc.Server) *hostif {
	p := &hostif{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterHostifServer(s, p)
	return p
}

type hostif struct {
	saipb.UnimplementedHostifServer
	mgr       *attrmgr.AttrMgr
	dataplane portDataplaneAPI
}

// CreateHostif creates a hostif interface (usually a tap interface).
func (hostif *hostif) CreateHostif(ctx context.Context, req *saipb.CreateHostifRequest) (*saipb.CreateHostifResponse, error) {
	id := hostif.mgr.NextID()

	switch req.GetType() {
	case saipb.HostifType_HOSTIF_TYPE_GENETLINK: // TODO: support this type
		return &saipb.CreateHostifResponse{Oid: id}, nil
	case saipb.HostifType_HOSTIF_TYPE_NETDEV:
		portReq := &dpb.CreatePortRequest{
			Id:   fmt.Sprint(id),
			Type: fwdpb.PortType_PORT_TYPE_KERNEL,
			Src: &dpb.CreatePortRequest_KernelDev{
				KernelDev: string(req.GetName()),
			},
			Location: dpb.PortLocation_PORT_LOCATION_INTERNAL,
		}
		attrReq := &saipb.GetPortAttributeRequest{Oid: req.GetObjId(), AttrType: []saipb.PortAttr{saipb.PortAttr_PORT_ATTR_OPER_STATUS}}
		p := &saipb.GetPortAttributeResponse{}
		if err := hostif.mgr.PopulateAttributes(attrReq, p); err != nil {
			return nil, err
		}
		if p.GetAttr().GetOperStatus() != saipb.PortOperStatus_PORT_OPER_STATUS_NOT_PRESENT {
			portReq.ExternalPort = fmt.Sprint(req.GetObjId())
		}
		if _, err := hostif.dataplane.CreatePort(ctx, portReq); err != nil {
			if err != nil {
				return nil, err
			}
		}
		attr := &saipb.HostifAttribute{
			OperStatus: proto.Bool(true),
		}
		hostif.mgr.StoreAttributes(id, attr)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown type %v", req.GetType())
	}
	return &saipb.CreateHostifResponse{Oid: id}, nil
}

// SetHostifAttribute sets the attributes in the request.
func (hostif *hostif) SetHostifAttribute(ctx context.Context, req *saipb.SetHostifAttributeRequest) (*saipb.SetHostifAttributeResponse, error) {
	if req.OperStatus != nil {
		stateReq := &fwdpb.PortStateRequest{
			ContextId: &fwdpb.ContextId{Id: hostif.dataplane.ID()},
			PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(req.GetOid())}},
		}
		stateReq.Operation = &fwdpb.PortInfo{
			AdminStatus: fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
		}
		if req.GetOperStatus() {
			stateReq.Operation.AdminStatus = fwdpb.PortState_PORT_STATE_ENABLED_UP
		}
		_, err := hostif.dataplane.PortState(ctx, stateReq)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
