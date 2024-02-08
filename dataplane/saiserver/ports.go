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
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/cpusink"
	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	log "github.com/golang/glog"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func newPort(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server, opts *dplaneopts.Options) (*port, error) {
	p := &port{
		mgr:       mgr,
		dataplane: dataplane,
		portToEth: make(map[uint64]string),
		nextEth:   1, // Start at eth1
		opts:      opts,
	}
	if opts.PortConfigFile != "" {
		data, err := os.ReadFile(opts.PortConfigFile)
		if err != nil {
			return nil, err
		}
		p.config = &dplaneopts.PortConfig{}
		if err := json.Unmarshal(data, p.config); err != nil {
			return nil, err
		}
	}

	saipb.RegisterPortServer(s, p)
	return p, nil
}

type port struct {
	saipb.UnimplementedPortServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
	nextEth   int
	portToEth map[uint64]string
	opts      *dplaneopts.Options
	config    *dplaneopts.PortConfig
}

// stub for testing
var getInterface = net.InterfaceByName

// CreatePort creates a new port, mapping the port to ethX, where X is assigned sequentially from 1 to n.
// Note: If more ports are created than eth devices, no error is returned, but the OperStatus is set to NOT_PRESENT.
func (port *port) CreatePort(ctx context.Context, req *saipb.CreatePortRequest) (*saipb.CreatePortResponse, error) {
	id := port.mgr.NextID()

	// By default, create port sequentially starting at eth1.
	dev := fmt.Sprintf("eth%v", port.nextEth)
	port.nextEth++

	// If a port config is set, then use the hardware lanes to find the interface names.
	if port.config != nil {
		if len(req.HwLaneList) == 0 {
			return nil, fmt.Errorf("port lanes are required got %v", req.HwLaneList)
		}
		var b strings.Builder
		b.WriteString(fmt.Sprint(req.HwLaneList[0]))
		for _, l := range req.HwLaneList[1:] {
			b.WriteString(",")
			b.WriteString(fmt.Sprint(l))
		}
		dev = ""
		lanes := b.String()
		for n, cfg := range port.config.Ports {
			if cfg.Lanes == lanes {
				dev = n
			}
		}
		if dev == "" {
			return nil, fmt.Errorf("could not find lanes %v in config", lanes)
		}
		if port.opts.PortMap != nil && len(port.opts.PortMap) != 0 {
			if portMapPort := port.opts.PortMap[dev]; portMapPort == "" {
				log.Warningf("port named %q doesn't exist in the port map", dev)
			} else {
				dev = portMapPort
			}
		}
	}

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

	fwdPort := &fwdpb.PortCreateRequest{
		ContextId: &fwdpb.ContextId{Id: port.dataplane.ID()},
		Port: &fwdpb.PortDesc{
			PortId: &fwdpb.PortId{
				ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(id)},
			},
		},
	}

	switch port.opts.PortType {
	case fwdpb.PortType_PORT_TYPE_KERNEL:
		fwdPort.Port.Port = &fwdpb.PortDesc_Kernel{
			Kernel: &fwdpb.KernelPortDesc{
				DeviceName: dev,
			},
		}
		// For ports that don't exist, do not create dataplane ports.
		if _, err := getInterface(dev); err != nil {
			attrs.OperStatus = saipb.PortOperStatus_PORT_OPER_STATUS_NOT_PRESENT.Enum()
			port.mgr.StoreAttributes(id, attrs)
			// TODO: This should be a real error, improve once we a correct config solution.
			// For now, create dummy port with no actions so we don't get a bunch error for a nonexistant port.
			fwdPort := &fwdpb.PortCreateRequest{
				ContextId: &fwdpb.ContextId{Id: port.dataplane.ID()},
				Port: &fwdpb.PortDesc{
					PortType: fwdpb.PortType_PORT_TYPE_GENETLINK,
					PortId: &fwdpb.PortId{
						ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(id)},
					},
				},
			}
			_, err := port.dataplane.PortCreate(ctx, fwdPort)
			if err != nil {
				return nil, err
			}
			return &saipb.CreatePortResponse{
				Oid: id,
			}, nil
		}
		port.portToEth[id] = dev

	case fwdpb.PortType_PORT_TYPE_FAKE:
		fwdPort.Port.Port = &fwdpb.PortDesc_Fake{
			Fake: &fwdpb.FakePortDesc{},
		}
	default:
		return nil, fmt.Errorf("unsupported port type: %v", port.opts.PortType)
	}
	fwdPort.Port.PortType = port.opts.PortType

	_, err := port.dataplane.PortCreate(ctx, fwdPort)
	if err != nil {
		return nil, err
	}
	update := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: port.dataplane.ID()},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(id)}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_Kernel{
				Kernel: &fwdpb.KernelPortUpdateDesc{
					Inputs: []*fwdpb.ActionDesc{
						fwdconfig.Action(fwdconfig.LookupAction(inputIfaceTable)).Build(),                               // Match packet to interface.
						fwdconfig.Action(fwdconfig.LookupAction(IngressVRFTable)).Build(),                               /// Match interface to VRF.
						fwdconfig.Action(fwdconfig.LookupAction(PreIngressActionTable)).Build(),                         // Run pre-ingress actions.
						fwdconfig.Action(fwdconfig.DecapAction(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET)).Build(), // Decap L2 header.
						fwdconfig.Action(fwdconfig.LookupAction(IngressActionTable)).Build(),                            // Run ingress action.
						fwdconfig.Action(fwdconfig.LookupAction(FIBSelectorTable)).Build(),                              // Lookup in FIB.
						fwdconfig.Action(fwdconfig.EncapAction(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET)).Build(), // Encap L2 header.
						fwdconfig.Action(fwdconfig.LookupAction(outputIfaceTable)).Build(),                              // Match interface to port
						fwdconfig.Action(fwdconfig.LookupAction(NeighborTable)).Build(),                                 // Lookup in the neighbor table.
						fwdconfig.Action(fwdconfig.LookupAction(EgressActionTable)).Build(),                             // Run egress actions
						fwdconfig.Action(fwdconfig.LookupAction(SRCMACTable)).Build(),                                   // Lookup interface's MAC addr.
						{
							ActionType: fwdpb.ActionType_ACTION_TYPE_OUTPUT,
						},
					},
					Outputs: []*fwdpb.ActionDesc{},
				},
			},
		},
	}
	if _, err := port.dataplane.PortUpdate(ctx, update); err != nil {
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

	_, err := port.dataplane.PortCreate(ctx, &fwdpb.PortCreateRequest{
		ContextId: &fwdpb.ContextId{Id: port.dataplane.ID()},
		Port: &fwdpb.PortDesc{
			PortType: fwdpb.PortType_PORT_TYPE_CPU_PORT,
			PortId:   &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(id)}},
			Port: &fwdpb.PortDesc_Cpu{
				Cpu: &fwdpb.CPUPortDesc{},
			},
		},
	})
	if err != nil {
		return 0, err
	}
	_, err = port.dataplane.PortUpdate(ctx, &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: port.dataplane.ID()},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(id)}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_Cpu{
				Cpu: &fwdpb.CPUPortUpdateDesc{
					Outputs: []*fwdpb.ActionDesc{
						fwdconfig.Action(fwdconfig.LookupAction(hostifTable)).Build(),
						fwdconfig.Action(fwdconfig.LookupAction(cpusink.IP2MeTable)).Build(),
					},
				},
			},
		},
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

// GetPortStats returns the stats for a port.
func (port *port) GetPortStats(ctx context.Context, req *saipb.GetPortStatsRequest) (*saipb.GetPortStatsResponse, error) {
	resp := &saipb.GetPortStatsResponse{}
	counters, err := port.dataplane.ObjectCounters(ctx, &fwdpb.ObjectCountersRequest{
		ContextId: &fwdpb.ContextId{Id: port.dataplane.ID()},
		ObjectId:  &fwdpb.ObjectId{Id: fmt.Sprint(req.GetOid())},
	})
	if err != nil {
		return nil, err
	}
	counterMap := map[fwdpb.CounterId]uint64{}
	for _, c := range counters.GetCounters() {
		counterMap[c.GetId()] = c.GetValue()
	}

	for _, id := range req.GetCounterIds() {
		switch id {
		case saipb.PortStat_PORT_STAT_IF_IN_UCAST_PKTS:
			resp.Values = append(resp.Values, counterMap[fwdpb.CounterId_COUNTER_ID_RX_UCAST_PACKETS])
		case saipb.PortStat_PORT_STAT_IF_IN_NON_UCAST_PKTS:
			resp.Values = append(resp.Values, counterMap[fwdpb.CounterId_COUNTER_ID_RX_NON_UCAST_PACKETS])
		case saipb.PortStat_PORT_STAT_IF_IN_ERRORS:
			resp.Values = append(resp.Values, counterMap[fwdpb.CounterId_COUNTER_ID_RX_ERROR_PACKETS])
		case saipb.PortStat_PORT_STAT_IF_OUT_UCAST_PKTS:
			resp.Values = append(resp.Values, counterMap[fwdpb.CounterId_COUNTER_ID_TX_UCAST_PACKETS])
		case saipb.PortStat_PORT_STAT_IF_OUT_NON_UCAST_PKTS:
			resp.Values = append(resp.Values, counterMap[fwdpb.CounterId_COUNTER_ID_TX_NON_UCAST_PACKETS])
		case saipb.PortStat_PORT_STAT_IF_OUT_ERRORS:
			resp.Values = append(resp.Values, counterMap[fwdpb.CounterId_COUNTER_ID_TX_ERROR_PACKETS])
		case saipb.PortStat_PORT_STAT_IF_IN_OCTETS:
			resp.Values = append(resp.Values, counterMap[fwdpb.CounterId_COUNTER_ID_RX_OCTETS])
		case saipb.PortStat_PORT_STAT_IF_OUT_OCTETS:
			resp.Values = append(resp.Values, counterMap[fwdpb.CounterId_COUNTER_ID_TX_OCTETS])
		case saipb.PortStat_PORT_STAT_IF_IN_DISCARDS:
			resp.Values = append(resp.Values, counterMap[fwdpb.CounterId_COUNTER_ID_RX_DROP_PACKETS])
		case saipb.PortStat_PORT_STAT_IF_OUT_DISCARDS:
			resp.Values = append(resp.Values, counterMap[fwdpb.CounterId_COUNTER_ID_TX_DROP_PACKETS])
		default:
			resp.Values = append(resp.Values, 0)
			log.Infof("unknown port stat: %v", id)
		}
	}
	return resp, nil
}

func (port *port) Reset() {
	port.portToEth = make(map[uint64]string)
	port.nextEth = 1
}

type lagMember struct {
	lagID  uint64
	portID uint64
}

type lag struct {
	saipb.UnimplementedLagServer
	mgr         *attrmgr.AttrMgr
	dataplane   switchDataplaneAPI
	memberships map[uint64]*lagMember
}

func newLAG(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *lag {
	l := &lag{
		mgr:         mgr,
		dataplane:   dataplane,
		memberships: map[uint64]*lagMember{},
	}
	saipb.RegisterLagServer(s, l)
	return l
}

func (l *lag) Reset() {
	l.memberships = make(map[uint64]*lagMember)
}

func (l *lag) CreateLag(ctx context.Context, _ *saipb.CreateLagRequest) (*saipb.CreateLagResponse, error) {
	id := l.mgr.NextID()

	pReq := &fwdpb.PortCreateRequest{
		ContextId: &fwdpb.ContextId{Id: l.dataplane.ID()},
		Port: &fwdpb.PortDesc{
			PortType: fwdpb.PortType_PORT_TYPE_AGGREGATE_PORT,
			PortId:   &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(id)}},
		},
	}
	_, err := l.dataplane.PortCreate(ctx, pReq)
	if err != nil {
		return nil, err
	}

	upd := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: l.dataplane.ID()},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(id)}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_AggregateAlgo{
				AggregateAlgo: &fwdpb.AggregatePortAlgorithmUpdateDesc{
					Hash: fwdpb.AggregateHashAlgorithm_AGGREGATE_HASH_ALGORITHM_CRC32,
					FieldIds: []*fwdpb.PacketFieldId{
						{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC}},
						{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST}},
						{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC}},
						{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
						{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC}},
						{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST}},
					},
				},
			},
		},
	}
	_, err = l.dataplane.PortUpdate(ctx, upd)
	if err != nil {
		return nil, err
	}

	return &saipb.CreateLagResponse{Oid: id}, err
}

func (l *lag) CreateLagMember(ctx context.Context, req *saipb.CreateLagMemberRequest) (*saipb.CreateLagMemberResponse, error) {
	id := l.mgr.NextID()

	pReq := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: l.dataplane.ID()},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(req.GetLagId())}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_AggregateAdd{
				AggregateAdd: &fwdpb.AggregatePortAddMemberUpdateDesc{
					InstanceCount: 1,
					PortId:        &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(req.GetPortId())}},
				},
			},
		},
	}
	_, err := l.dataplane.PortUpdate(ctx, pReq)
	if err != nil {
		return nil, err
	}
	l.memberships[id] = &lagMember{lagID: req.GetLagId(), portID: req.GetPortId()}
	return &saipb.CreateLagMemberResponse{Oid: id}, err
}

func (l *lag) RemoveLagMember(ctx context.Context, req *saipb.RemoveLagMemberRequest) (*saipb.RemoveLagMemberResponse, error) {
	member := l.memberships[req.GetOid()]

	pReq := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: l.dataplane.ID()},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(member.lagID)}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_AggregateDel{
				AggregateDel: &fwdpb.AggregatePortRemoveMemberUpdateDesc{
					PortId: &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(member.portID)}},
				},
			},
		},
	}
	_, err := l.dataplane.PortUpdate(ctx, pReq)
	return &saipb.RemoveLagMemberResponse{}, err
}
