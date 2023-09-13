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

	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type portDataplaneAPI interface {
	ContextID() string
	CreatePort(ctx context.Context, req *dpb.CreatePortRequest) (*dpb.CreatePortResponse, error)
	PortState(ctx context.Context, req *fwdpb.PortStateRequest) (*fwdpb.PortStateReply, error)
}

type port struct {
	saipb.UnimplementedPortServer
	mgr       *attrmgr.AttrMgr
	dataplane portDataplaneAPI
	nextEth   int
}

func (port *port) CreatePort(ctx context.Context, _ *saipb.CreatePortRequest) (*saipb.CreatePortResponse, error) {
	id := port.mgr.NextID()
	port.nextEth += 1
	dev := fmt.Sprintf("eth%v", port.nextEth)

	attrs := &saipb.PortAttribute{
		OperSpeed:                     proto.Uint32(1024),
		NumberOfIngressPriorityGroups: proto.Uint32(0),
		QosNumberOfQueues:             proto.Uint32(0),
		QosMaximumHeadroomSize:        proto.Uint32(0),
		AdminState:                    proto.Bool(true),
		AutoNegMode:                   proto.Bool(false),
		Mtu:                           proto.Uint32(1514),
	}

	// For ports that don't exist, do not create dataplane ports.
	if _, err := net.InterfaceByName(dev); err != nil {
		attrs.OperStatus = saipb.PortOperStatus_PORT_OPER_STATUS_NOT_PRESENT.Enum()
		port.mgr.StoreAttributes(id, attrs)
		return &saipb.CreatePortResponse{
			Oid: id,
		}, nil
	}

	_, err := port.dataplane.CreatePort(ctx, &dpb.CreatePortRequest{
		Id:   fmt.Sprint(id),
		Type: fwdpb.PortType_PORT_TYPE_KERNEL,
		Src: &dpb.CreatePortRequest_KernelDev{
			KernelDev: dev,
		},
	})
	if err != nil {
		return nil, err
	}
	attrs.OperStatus = saipb.PortOperStatus_PORT_OPER_STATUS_UP.Enum().Enum()
	port.mgr.StoreAttributes(id, attrs)

	return &saipb.CreatePortResponse{
		Oid: id,
	}, nil
}

func (port *port) SetPortAttribute(ctx context.Context, req *saipb.SetPortAttributeRequest) (*saipb.SetPortAttributeResponse, error) {
	if req.AdminState != nil {
		stateReq := &fwdpb.PortStateRequest{
			ContextId: &fwdpb.ContextId{Id: port.dataplane.ContextID()},
			PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(req.GetOid())}},
		}
		stateReq.Operation.AdminStatus = fwdpb.PortState_PORT_STATE_DISABLED_DOWN
		if req.GetAdminState() {
			stateReq.Operation.AdminStatus = fwdpb.PortState_PORT_STATE_ENABLED_UP
		}
		_, err := port.dataplane.PortState(ctx, stateReq)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
