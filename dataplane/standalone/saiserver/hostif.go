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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/internal/engine"
	"github.com/openconfig/lemming/dataplane/standalone/packetio/cpusink"
	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func newHostif(ctx context.Context, mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) (*hostif, error) {
	hostif := &hostif{
		mgr:              mgr,
		dataplane:        dataplane,
		trapIDToHostifID: map[uint64]uint64{},
		groupIDToQueue:   map[uint64]uint32{},
	}

	// Setup the packet io tables. A packet is punted by setting the output port to the CPU port.
	// There a two places where packets can be punted:
	//   1. pre-fib: the trap table contains rules that may punt the packets.
	//   2. fib: "IP2ME" routes, the fib contains routes for the IPs assigned to the hostif, these routes have the next hop as the CPU port.
	// Once a packet is sent to the CPU port, it must be matched to a hostif:
	//   1. ip2me: a table maps IP DST to hostif port. (populated by the CPU port).
	//   2. hostif table: a table the maps TRAP IP to the hostif. (trap id is set by the ACL actions).
	//   3. default/wildcard: each hostif is created with a corresponding port, use that mapping to determine correct hostif.
	// Once the output port is determined, based on the hostif type:
	//   1. For genetlink: send the packets using the CPU port gRPC connection.
	//   2. For netdev (lucius kernel/tap): write the packets directly to the hostif.

	// Create the trap table and add it to the end of ingress stage.
	_, err := hostif.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: hostif.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_FLOW,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: trapTableID}},
			Table: &fwdpb.TableDesc_Flow{
				Flow: &fwdpb.FlowTableDesc{
					BankCount: 1,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	_, err = hostif.dataplane.TableEntryAdd(ctx, fwdconfig.TableEntryAddRequest(hostif.dataplane.ID(), engine.IngressActionTable).
		AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ActionEntry("trap", fwdpb.ActionEntryDesc_INSERT_METHOD_APPEND)),
			fwdconfig.Action(fwdconfig.LookupAction(trapTableID))).
		Build(),
	)
	if err != nil {
		return nil, err
	}

	// Create the IP2MeTable and hostif tables, these map the packet to real hostif port.
	// These tables are set as output actions of the CPU port.
	_, err = hostif.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: hostif.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: cpusink.IP2MeTable}},
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
						},
					}},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE,
			}},
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = hostif.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: hostif.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: hostifTable}},
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TRAP_ID,
						},
					}},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_SWAP_OUTPUT_INTERNAL_EXTERNAL,
			}},
		},
	})
	if err != nil {
		return nil, err
	}
	saipb.RegisterHostifServer(s, hostif)
	return hostif, nil
}

type hostif struct {
	saipb.UnimplementedHostifServer
	mgr              *attrmgr.AttrMgr
	dataplane        switchDataplaneAPI
	trapIDToHostifID map[uint64]uint64
	groupIDToQueue   map[uint64]uint32
}

// CreateHostif creates a hostif interface (usually a tap interface).
func (hostif *hostif) CreateHostif(ctx context.Context, req *saipb.CreateHostifRequest) (*saipb.CreateHostifResponse, error) {
	id := hostif.mgr.NextID()

	switch req.GetType() {
	case saipb.HostifType_HOSTIF_TYPE_GENETLINK: // For genetlink device, pass the port description to the cpu sink.
		// First, create the port in the dataplane: this port is "virtual", it doesn't do any packet io.
		// It ensures that the packet metadata is correct.
		portReq := &fwdpb.PortCreateRequest{
			ContextId: &fwdpb.ContextId{Id: hostif.dataplane.ID()},
			Port: &fwdpb.PortDesc{
				PortId:   &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: fmt.Sprint(id)}},
				PortType: fwdpb.PortType_PORT_TYPE_GENETLINK,
				Port: &fwdpb.PortDesc_Genetlink{
					Genetlink: &fwdpb.GenetlinkPortDesc{
						FamilyName: string(req.GetName()),
						GroupName:  string(req.GetGenetlinkMcgrpName()),
					},
				},
			},
		}
		if _, err := hostif.dataplane.PortCreate(ctx, portReq); err != nil {
			return nil, err
		}
		// Notify the cpu sink about these port types.
		fwdCtx, err := hostif.dataplane.Context()
		if err != nil {
			return nil, err
		}
		fwdCtx.RLock()
		ps := fwdCtx.PacketSink()
		fwdCtx.RUnlock()
		ps(&fwdpb.PacketSinkResponse{
			Resp: &fwdpb.PacketSinkResponse_Port{
				Port: &fwdpb.PacketSinkPortInfo{
					Port: portReq.Port,
				},
			},
		})

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

var (
	etherTypeARP  = []byte{0x08, 0x06}
	udldDstMAC    = []byte{0x01, 0x00, 0x0C, 0xCC, 0xCC, 0xCC}
	etherTypeLLDP = []byte{0x88, 0xcc}
	ndDstMAC      = []byte{0x33, 0x33, 0x00, 0x00, 0x00, 0x00} // ND is generic IPv6 multicast MAC.
	ndDstMACMask  = []byte{0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00}
	lacpDstMAC    = []byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x02}
)

const (
	bgpPort        = 179
	trapTableID    = "trap-table"
	wildcardPortID = 0
)

func (hostif *hostif) CreateHostifTrap(ctx context.Context, req *saipb.CreateHostifTrapRequest) (*saipb.CreateHostifTrapResponse, error) {
	id := hostif.mgr.NextID()
	fwdReq := fwdconfig.TableEntryAddRequest(hostif.dataplane.ID(), trapTableID)

	swReq := &saipb.GetSwitchAttributeRequest{
		Oid:      req.GetSwitch(),
		AttrType: []saipb.SwitchAttr{saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT},
	}
	swAttr := &saipb.GetSwitchAttributeResponse{}
	if err := hostif.mgr.PopulateAttributes(swReq, swAttr); err != nil {
		return nil, err
	}
	entriesAdded := 1
	switch tType := req.GetTrapType(); tType {
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_ARP_REQUEST, saipb.HostifTrapType_HOSTIF_TRAP_TYPE_ARP_RESPONSE:
		fwdReq.AppendEntry(fwdconfig.EntryDesc(fwdconfig.FlowEntry(
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE).
				WithBytes(etherTypeARP, []byte{0xFF, 0xFF}))))
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_UDLD:
		fwdReq.AppendEntry(fwdconfig.EntryDesc(fwdconfig.FlowEntry(
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).
				WithBytes(udldDstMAC, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}))))
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_LLDP:
		fwdReq.AppendEntry(fwdconfig.EntryDesc(fwdconfig.FlowEntry(
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE).
				WithBytes(etherTypeLLDP, []byte{0xFF, 0xFF}))))
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_DISCOVERY:
		fwdReq.AppendEntry(fwdconfig.EntryDesc(fwdconfig.FlowEntry(
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).
				WithBytes(ndDstMAC, ndDstMACMask))))
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_LACP:
		fwdReq.AppendEntry(fwdconfig.EntryDesc(fwdconfig.FlowEntry(
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).
				WithBytes(lacpDstMAC, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}))))
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_IP2ME:
		// IP2ME routes are added to the FIB, do nothing here.
		return &saipb.CreateHostifTrapResponse{
			Oid: id,
		}, nil
	case saipb.HostifTrapType_HOSTIF_TRAP_TYPE_BGP, saipb.HostifTrapType_HOSTIF_TRAP_TYPE_BGPV6:
		// TODO: This should only match for packets destined to the management IP.
		fwdReq.AppendEntry(fwdconfig.EntryDesc(fwdconfig.FlowEntry(
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC).
				WithUint16(bgpPort))),
		)
		fwdReq.AppendEntry(fwdconfig.EntryDesc(fwdconfig.FlowEntry(
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST).
				WithUint16(bgpPort))),
		)
		entriesAdded = 2
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown trap type: %v", tType)
	}

	switch act := req.GetPacketAction(); act {
	case saipb.PacketAction_PACKET_ACTION_TRAP: // TRAP means COPY to CPU and DROP, just transmit immediately, which interrupts any pending actions.
		for i := 0; i < entriesAdded; i++ {
			fwdReq.AppendActions(fwdconfig.Action(fwdconfig.TransmitAction(fmt.Sprint(swAttr.GetAttr().GetCpuPort())).WithImmediate(true)))
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown action type: %v", act)
	}
	if _, err := hostif.dataplane.TableEntryAdd(ctx, fwdReq.Build()); err != nil {
		return nil, err
	}
	// TODO: Support multiple queues, by using the group ID.
	return &saipb.CreateHostifTrapResponse{
		Oid: id,
	}, nil
}

func (hostif *hostif) CreateHostifTrapGroup(_ context.Context, req *saipb.CreateHostifTrapGroupRequest) (*saipb.CreateHostifTrapGroupResponse, error) {
	id := hostif.mgr.NextID()
	hostif.groupIDToQueue[id] = req.GetQueue()
	return &saipb.CreateHostifTrapGroupResponse{Oid: id}, nil
}

func (hostif *hostif) CreateHostifUserDefinedTrap(_ context.Context, req *saipb.CreateHostifUserDefinedTrapRequest) (*saipb.CreateHostifUserDefinedTrapResponse, error) {
	if req.GetType() != saipb.HostifUserDefinedTrapType_HOSTIF_USER_DEFINED_TRAP_TYPE_ACL {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported trap type: %v", req.GetType())
	}
	return &saipb.CreateHostifUserDefinedTrapResponse{Oid: hostif.mgr.NextID()}, nil
}

const (
	hostifTable = "hostiftable"
)

func (hostif *hostif) CreateHostifTableEntry(ctx context.Context, req *saipb.CreateHostifTableEntryRequest) (*saipb.CreateHostifTableEntryResponse, error) {
	switch entryType := req.GetType(); entryType {
	case saipb.HostifTableEntryType_HOSTIF_TABLE_ENTRY_TYPE_TRAP_ID:
		hostif.trapIDToHostifID[req.GetTrapId()] = req.GetHostIf()
		_, err := hostif.dataplane.TableEntryAdd(ctx, fwdconfig.TableEntryAddRequest(hostif.dataplane.ID(), hostifTable).
			AppendEntry(
				fwdconfig.EntryDesc(fwdconfig.ExactEntry(
					fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TRAP_ID).WithUint64(req.GetTrapId()))),
				fwdconfig.Action(fwdconfig.TransmitAction(fmt.Sprint(req.GetHostIf())))).
			Build())
		if err != nil {
			return nil, err
		}
	case saipb.HostifTableEntryType_HOSTIF_TABLE_ENTRY_TYPE_WILDCARD:
		hostif.trapIDToHostifID[req.GetTrapId()] = wildcardPortID
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported entry type: %v", entryType)
	}

	return nil, nil
}
