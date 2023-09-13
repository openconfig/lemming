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

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
)

type saiSwitch struct {
	saipb.UnimplementedSwitchServer
	port   *port
	vlan   *vlan
	stp    *stp
	vr     *virtualRouter
	bridge *bridge
	hostif *hostif
	hash   *hash
	mgr    *attrmgr.AttrMgr
}

func newSwitch(mgr *attrmgr.AttrMgr, s *grpc.Server) *saiSwitch {
	sw := &saiSwitch{
		port:   &port{},
		vlan:   &vlan{},
		stp:    &stp{},
		vr:     &virtualRouter{},
		bridge: &bridge{},
		hostif: &hostif{},
		hash:   &hash{},
		mgr:    mgr,
	}
	saipb.RegisterSwitchServer(s, sw)
	saipb.RegisterPortServer(s, sw.port)
	saipb.RegisterVlanServer(s, sw.vlan)
	saipb.RegisterStpServer(s, sw.stp)
	saipb.RegisterVirtualRouterServer(s, sw.vr)
	saipb.RegisterBridgeServer(s, sw.bridge)
	saipb.RegisterHostifServer(s, sw.hostif)
	saipb.RegisterHashServer(s, sw.hash)
	return sw
}

// CreateSwitch a creates a new switch and populates its default values.
func (sw *saiSwitch) CreateSwitch(ctx context.Context, _ *saipb.CreateSwitchRequest) (*saipb.CreateSwitchResponse, error) {
	swID := sw.mgr.NextID()

	// TODO: The port type is not a settable attribute, figure out a pattern for this.
	cpuPortID := sw.mgr.NextID()
	cpuPort := &saipb.PortAttribute{
		Type: saipb.PortType_PORT_TYPE_CPU.Enum(),
	}
	sw.mgr.SetType(fmt.Sprint(cpuPortID), saipb.ObjectType_OBJECT_TYPE_PORT)
	vlanResp, err := attrmgr.InvokeAndSave(ctx, sw.mgr, sw.vlan.CreateVlan, &saipb.CreateVlanRequest{
		Switch: swID,
	})
	if err != nil {
		return nil, err
	}
	stpResp, err := attrmgr.InvokeAndSave(ctx, sw.mgr, sw.stp.CreateStp, &saipb.CreateStpRequest{
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
	attrs := &saipb.SwitchAttribute{
		CpuPort:                          proto.Uint64(cpuPortID),
		NumberOfActivePorts:              proto.Uint32(0),
		AclEntryMinimumPriority:          proto.Uint32(1),
		AclTableMaximumPriority:          proto.Uint32(100),
		MaxAclActionCount:                proto.Uint32(50),
		NumberOfEcmpGroups:               proto.Uint32(1024),
		DefaultVlanId:                    &vlanResp.Oid,
		DefaultStpInstId:                 &stpResp.Oid,
		DefaultVirtualRouterId:           &vrResp.Oid,
		DefaultOverrideVirtualRouterId:   &vrResp.Oid,
		Default_1QBridgeId:               &brResp.Oid,
		DefaultTrapGroup:                 &trGroupResp.Oid,
		IngressAcl:                       proto.Uint64(0),
		EgressAcl:                        proto.Uint64(0),
		QosMaxNumberOfTrafficClasses:     proto.Uint32(0),
		TotalBufferSize:                  proto.Uint64(1024 * 1024),
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
		EcmpHash:                         &hashResp.Oid,
		LagHash:                          &hashResp.Oid,
		RestartWarm:                      proto.Bool(false),
		WarmRecover:                      proto.Bool(false),
		LagDefaultHashAlgorithm:          saipb.HashAlgorithm_HASH_ALGORITHM_CRC.Enum(),
		LagDefaultHashSeed:               proto.Uint32(0),
		LagDefaultSymmetricHash:          proto.Bool(false),
		QosDefaultTc:                     proto.Uint32(0),
		QosDot1PToTcMap:                  proto.Uint64(0),
		QosDot1PToColorMap:               proto.Uint64(0),
		QosTcToQueueMap:                  proto.Uint64(0),
		QosTcAndColorToDot1PMap:          proto.Uint64(0),
		QosTcAndColorToDscpMap:           proto.Uint64(0),
		QosTcAndColorToMplsExpMap:        proto.Uint64(0),
		SwitchShellEnable:                proto.Bool(false),
		SwitchProfileId:                  proto.Uint32(0),
		NatZoneCounterObjectId:           proto.Uint64(0),
	}
	sw.mgr.StoreAttributes(swID, attrs)
	sw.mgr.StoreAttributes(cpuPortID, cpuPort)

	return &saipb.CreateSwitchResponse{
		Oid: swID,
	}, nil
}
