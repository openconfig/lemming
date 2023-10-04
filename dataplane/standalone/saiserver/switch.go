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
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type saiSwitch struct {
	saipb.UnimplementedSwitchServer
	dataplane       switchDataplaneAPI
	port            *port
	vlan            *vlan
	stp             *stp
	vr              *virtualRouter
	bridge          *bridge
	hostif          *hostif
	hash            *hash
	neighbor        *neighbor
	nextHopGroup    *nextHopGroup
	nextHop         *nextHop
	route           *route
	routerInterface *routerInterface
	mgr             *attrmgr.AttrMgr
}

type switchDataplaneAPI interface {
	portDataplaneAPI
	routingDataplaneAPI
	NotifySubscribe(sub *fwdpb.NotifySubscribeRequest, srv fwdpb.Forwarding_NotifySubscribeServer) error
}

func newSwitch(mgr *attrmgr.AttrMgr, engine switchDataplaneAPI, s *grpc.Server) *saiSwitch {
	sw := &saiSwitch{
		dataplane:       engine,
		port:            newPort(mgr, engine, s),
		vlan:            newVlan(mgr, engine, s),
		stp:             &stp{},
		vr:              &virtualRouter{},
		bridge:          newBridge(mgr, engine, s),
		hostif:          newHostif(mgr, engine, s),
		hash:            &hash{},
		neighbor:        newNeighbor(mgr, engine, s),
		nextHopGroup:    newNextHopGroup(mgr, engine, s),
		nextHop:         newNextHop(mgr, engine, s),
		route:           newRoute(mgr, engine, s),
		routerInterface: newRouterInterface(mgr, engine, s),
		mgr:             mgr,
	}
	saipb.RegisterSwitchServer(s, sw)
	saipb.RegisterStpServer(s, sw.stp)
	saipb.RegisterVirtualRouterServer(s, sw.vr)
	saipb.RegisterHashServer(s, sw.hash)
	return sw
}

// CreateSwitch a creates a new switch and populates its default values.
func (sw *saiSwitch) CreateSwitch(ctx context.Context, _ *saipb.CreateSwitchRequest) (*saipb.CreateSwitchResponse, error) {
	swID := sw.mgr.NextID()

	cpuPortID, err := sw.port.createCPUPort(ctx)
	if err != nil {
		return nil, err
	}

	stpResp, err := attrmgr.InvokeAndSave(ctx, sw.mgr, sw.stp.CreateStp, &saipb.CreateStpRequest{
		Switch: swID,
	})
	if err != nil {
		return nil, err
	}
	sw.mgr.StoreAttributes(swID, &saipb.SwitchAttribute{DefaultStpInstId: &stpResp.Oid})

	vlanResp, err := attrmgr.InvokeAndSave(ctx, sw.mgr, sw.vlan.CreateVlan, &saipb.CreateVlanRequest{
		Switch: swID,
	})
	if err != nil {
		return nil, err
	}

	vrResp, err := attrmgr.InvokeAndSave(ctx, sw.mgr, sw.vr.CreateVirtualRouter, &saipb.CreateVirtualRouterRequest{
		Switch: swID,
	})
	if err != nil {
		return nil, err
	}
	brResp, err := attrmgr.InvokeAndSave(ctx, sw.mgr, sw.bridge.CreateBridge, &saipb.CreateBridgeRequest{
		Switch: swID,
	})
	if err != nil {
		return nil, err
	}
	trGroupResp, err := attrmgr.InvokeAndSave(ctx, sw.mgr, sw.hostif.CreateHostifTrapGroup, &saipb.CreateHostifTrapGroupRequest{
		Switch: swID,
	})
	if err != nil {
		return nil, err
	}
	hashResp, err := attrmgr.InvokeAndSave(ctx, sw.mgr, sw.hash.CreateHash, &saipb.CreateHashRequest{
		Switch: swID,
	})
	if err != nil {
		return nil, err
	}

	// These values are mostly meaningless, but clients expect these to be set.
	// The values either the default value for the attribute (https://github.com/opencomputeproject/SAI/blob/master/inc/saiswitch.h)
	// or for unsupported features a zero value.
	attrs := &saipb.SwitchAttribute{
		CpuPort:                          proto.Uint64(cpuPortID),
		NumberOfActivePorts:              proto.Uint32(0),
		AclEntryMinimumPriority:          proto.Uint32(1),
		AclEntryMaximumPriority:          proto.Uint32(100),
		AclTableMinimumPriority:          proto.Uint32(1),
		AclTableMaximumPriority:          proto.Uint32(100),
		MaxAclActionCount:                proto.Uint32(50),
		NumberOfEcmpGroups:               proto.Uint32(1024),
		PortList:                         []uint64{cpuPortID},
		SwitchHardwareInfo:               []int32{},
		DefaultVlanId:                    &vlanResp.Oid,
		DefaultVirtualRouterId:           &vrResp.Oid,
		DefaultOverrideVirtualRouterId:   &vrResp.Oid,
		Default_1QBridgeId:               &brResp.Oid,
		DefaultTrapGroup:                 &trGroupResp.Oid,
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
		EcmpHash:                       &hashResp.Oid,
		LagHash:                        &hashResp.Oid,
		EcmpHashIpv4:                   &hashResp.Oid,
		EcmpHashIpv4InIpv4:             &hashResp.Oid,
		EcmpHashIpv6:                   &hashResp.Oid,
		LagHashIpv4:                    &hashResp.Oid,
		LagHashIpv4InIpv4:              &hashResp.Oid,
		LagHashIpv6:                    &hashResp.Oid,
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
	sw.mgr.StoreAttributes(swID, attrs)

	return &saipb.CreateSwitchResponse{
		Oid: swID,
	}, nil
}

type fwdNotifServer struct {
	fwdpb.Forwarding_NotifySubscribeServer
	ch chan *fwdpb.EventDesc
}

func (s *fwdNotifServer) Send(ed *fwdpb.EventDesc) error {
	s.ch <- ed
	return nil
}

func (sw *saiSwitch) PortStateChangeNotification(_ *saipb.PortStateChangeNotificationRequest, srv saipb.Switch_PortStateChangeNotificationServer) error {
	req := &fwdpb.NotifySubscribeRequest{
		Context: &fwdpb.ContextId{
			Id: sw.dataplane.ID(),
		},
	}
	fwdSrv := &fwdNotifServer{
		ch: make(chan *fwdpb.EventDesc, 1),
	}
	errCh := make(chan error)
	go func() {
		errCh <- sw.dataplane.NotifySubscribe(req, fwdSrv)
	}()
	for {
		select {
		case err := <-errCh:
			return err
		case ed := <-fwdSrv.ch:
			num, err := strconv.Atoi(ed.GetPort().GetPortId().GetObjectId().GetId())
			if err != nil {
				log.Warningf("couldn't get numeric port id: %v", err)
				continue
			}
			oType := sw.mgr.GetType(ed.GetPort().GetPortId().GetObjectId().GetId())
			if oType != saipb.ObjectType_OBJECT_TYPE_PORT {
				log.Infof("skipping port state event for type %v", oType)
				continue
			}
			status := saipb.PortOperStatus_PORT_OPER_STATUS_UNKNOWN
			if ed.GetPort().PortInfo.OperStatus == fwdpb.PortState_PORT_STATE_ENABLED_UP {
				status = saipb.PortOperStatus_PORT_OPER_STATUS_UP
			} else if ed.GetPort().PortInfo.OperStatus == fwdpb.PortState_PORT_STATE_DISABLED_DOWN {
				status = saipb.PortOperStatus_PORT_OPER_STATUS_DOWN
			}

			err = srv.Send(&saipb.PortStateChangeNotificationResponse{
				Data: []*saipb.PortOperStatusNotification{{
					PortId:    uint64(num),
					PortState: status,
				}},
			})
			if err != nil {
				log.Warningf("failed to send port event %v", err)
			}
		}
	}
}

func (sw saiSwitch) Reset() {
	sw.port.Reset()
}
