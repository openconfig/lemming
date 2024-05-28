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
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/dataplane/cpusink"
	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type saiSwitch struct {
	saipb.UnimplementedSwitchServer
	dataplane       switchDataplaneAPI
	acl             *acl
	port            *port
	vlan            *vlan
	stp             *stp
	vr              *virtualRouter
	bridge          *bridge
	hostif          *hostif
	hash            *hash
	isolationGroup  *isolationGroup
	myMac           *myMac
	neighbor        *neighbor
	nextHopGroup    *nextHopGroup
	nextHop         *nextHop
	policer         *policer
	route           *route
	lag             *lag
	tunnel          *tunnel
	routerInterface *routerInterface
	udf             *udf
	mgr             *attrmgr.AttrMgr
}

type switchDataplaneAPI interface {
	NotifySubscribe(sub *fwdpb.NotifySubscribeRequest, srv fwdpb.Forwarding_NotifySubscribeServer) error
	TableCreate(context.Context, *fwdpb.TableCreateRequest) (*fwdpb.TableCreateReply, error)
	TableEntryAdd(context.Context, *fwdpb.TableEntryAddRequest) (*fwdpb.TableEntryAddReply, error)
	TableEntryRemove(context.Context, *fwdpb.TableEntryRemoveRequest) (*fwdpb.TableEntryRemoveReply, error)
	ID() string
	PortState(ctx context.Context, req *fwdpb.PortStateRequest) (*fwdpb.PortStateReply, error)
	ObjectCounters(context.Context, *fwdpb.ObjectCountersRequest) (*fwdpb.ObjectCountersReply, error)
	FindContext(*fwdpb.ContextId) (*fwdcontext.Context, error)
	PortCreate(context.Context, *fwdpb.PortCreateRequest) (*fwdpb.PortCreateReply, error)
	PortUpdate(context.Context, *fwdpb.PortUpdateRequest) (*fwdpb.PortUpdateReply, error)
	AttributeUpdate(context.Context, *fwdpb.AttributeUpdateRequest) (*fwdpb.AttributeUpdateReply, error)
	ObjectNID(context.Context, *fwdpb.ObjectNIDRequest) (*fwdpb.ObjectNIDReply, error)
	InjectPacket(contextID *fwdpb.ContextId, id *fwdpb.PortId, hid fwdpb.PacketHeaderId, frame []byte, preActions []*fwdpb.ActionDesc, debug bool, dir fwdpb.PortAction) error
	ObjectDelete(context.Context, *fwdpb.ObjectDeleteRequest) (*fwdpb.ObjectDeleteReply, error)
	FlowCounterCreate(_ context.Context, request *fwdpb.FlowCounterCreateRequest) (*fwdpb.FlowCounterCreateReply, error)
	FlowCounterQuery(_ context.Context, request *fwdpb.FlowCounterQueryRequest) (*fwdpb.FlowCounterQueryReply, error)
}

const (
	inputIfaceTable       = "input-iface"
	outputIfaceTable      = "output-iface"
	IngressVRFTable       = "ingress-vrf"
	FIBV4Table            = "fib-v4"
	FIBV6Table            = "fib-v6"
	SRCMACTable           = "port-mac"
	FIBSelectorTable      = "fib-selector"
	NeighborTable         = "neighbor"
	NHGTable              = "nhg-table"
	NHTable               = "nh-table"
	layer2PuntTable       = "layer2-punt"
	layer3PuntTable       = "layer3-punt"
	arpPuntTable          = "arp-punt"
	PreIngressActionTable = "preingress-table"
	IngressActionTable    = "ingress-table"
	EgressActionTable     = "egress-action-table"
	NHActionTable         = "nh-action"
	TunnelEncap           = "tunnel-encap"
	MyMacTable            = "my-mac-table"
	hostifToPortTable     = "cpu-input"
	portToHostifTable     = "cpu-output"
	tunTermTable          = "tun-term"
	VlanTable             = "vlan"
	DefaultVlanId         = 4095 // An reserved VLAN ID used as the default VLAN ID for internal usage.
)

func newSwitch(mgr *attrmgr.AttrMgr, engine switchDataplaneAPI, s *grpc.Server, opts *dplaneopts.Options) (*saiSwitch, error) {
	port, err := newPort(mgr, engine, s, opts)
	if err != nil {
		return nil, err
	}
	sw := &saiSwitch{
		dataplane:       engine,
		acl:             newACL(mgr, engine, s),
		policer:         newPolicer(mgr, engine, s),
		port:            port,
		vlan:            newVlan(mgr, engine, s),
		stp:             &stp{},
		vr:              &virtualRouter{},
		bridge:          newBridge(mgr, engine, s),
		hostif:          newHostif(mgr, engine, s, opts),
		hash:            newHash(mgr, engine, s),
		isolationGroup:  newIsolationGroup(mgr, engine, s),
		myMac:           newMyMac(mgr, engine, s),
		neighbor:        newNeighbor(mgr, engine, s),
		nextHopGroup:    newNextHopGroup(mgr, engine, s),
		nextHop:         newNextHop(mgr, engine, s),
		route:           newRoute(mgr, engine, s),
		routerInterface: newRouterInterface(mgr, engine, s),
		lag:             newLAG(mgr, engine, s),
		tunnel:          newTunnel(mgr, engine, s),
		udf:             newUdf(mgr, engine, s),
		mgr:             mgr,
	}
	saipb.RegisterSwitchServer(s, sw)
	saipb.RegisterStpServer(s, sw.stp)
	saipb.RegisterVirtualRouterServer(s, sw.vr)
	return sw, nil
}

// CreateSwitch a creates a new switch and populates its default values.
func (sw *saiSwitch) CreateSwitch(ctx context.Context, _ *saipb.CreateSwitchRequest) (*saipb.CreateSwitchResponse, error) {
	swID := sw.mgr.NextID()

	// Setup forwarding tables.
	ingressVRF := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: IngressVRFTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}}, // TODO: Should this be drop?
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_INPUT_IFACE,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, ingressVRF); err != nil {
		return nil, err
	}

	v4FIB := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_PREFIX,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: FIBV4Table}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Prefix{
				Prefix: &fwdpb.PrefixTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF,
						},
					}, {
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, v4FIB); err != nil {
		return nil, err
	}
	v6FIB := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_PREFIX,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: FIBV6Table}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Prefix{
				Prefix: &fwdpb.PrefixTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF,
						},
					}, {
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, v6FIB); err != nil {
		return nil, err
	}
	portMAC := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: SRCMACTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, portMAC); err != nil {
		return nil, err
	}
	vlanReq := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: VlanTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, vlanReq); err != nil {
		return nil, err
	}
	myMAC := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_FLOW,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: MyMacTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Flow{
				Flow: &fwdpb.FlowTableDesc{
					BankCount: 1,
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, myMAC); err != nil {
		return nil, err
	}
	neighbor := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: NeighborTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE,
						},
					}, {
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, neighbor); err != nil {
		return nil, err
	}
	nh := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: NHTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, nh); err != nil {
		return nil, err
	}
	nhg := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: NHGTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_GROUP_ID,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, nhg); err != nil {
		return nil, err
	}
	if err := createFIBSelector(ctx, sw.dataplane.ID(), sw.dataplane); err != nil {
		return nil, err
	}

	action := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_ACTION,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: PreIngressActionTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
			Table: &fwdpb.TableDesc_Action{
				Action: &fwdpb.ActionTableDesc{},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, action); err != nil {
		return nil, err
	}
	action.Desc.TableId.ObjectId.Id = IngressActionTable
	if _, err := sw.dataplane.TableCreate(ctx, action); err != nil {
		return nil, err
	}
	action.Desc.TableId.ObjectId.Id = EgressActionTable
	if _, err := sw.dataplane.TableCreate(ctx, action); err != nil {
		return nil, err
	}
	nexthopAction := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: NHActionTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, nexthopAction); err != nil {
		return nil, err
	}

	// Setup the packet io tables. A packet is punted by setting the output port to the CPU port.
	// There a two places where packets can be punted:
	//   1. pre-fib: the trap table contains rules that may punt the packets.
	// Once a packet is sent to the CPU port, it must be matched to a hostif:
	//   1. ip2me: a table maps IP DST to hostif port. (populated by the CPU port).
	//   2. hostif table: a table the maps TRAP IP to the hostif. (trap id is set by the ACL actions).
	//   3. default/wildcard: each hostif is created with a corresponding port, use that mapping to determine correct hostif.
	// Once the output port is determined, based on the hostif type:
	//   1. For genetlink: send the packets using the CPU port gRPC connection.
	//   2. For netdev (lucius kernel/tap): write the packets directly to the hostif.

	// Create the trap table and add it to the end of ingress stage.
	_, err := sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
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
	// Add at the end of the pre-ingress table, since this need to happen before the layer 2 header is decapped.
	_, err = sw.dataplane.TableEntryAdd(ctx, fwdconfig.TableEntryAddRequest(sw.dataplane.ID(), PreIngressActionTable).
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
	_, err = sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: cpusink.IP2MeTable}},
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT,
						},
					}, {
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

	trapToHostifTableReq := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: trapIDToHostifTable}},
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
		},
	}

	if !sw.port.opts.RemoteCPUPort {
		trapToHostifTableReq.Desc.Actions = append(trapToHostifTableReq.Desc.Actions, &fwdpb.ActionDesc{ActionType: fwdpb.ActionType_ACTION_TYPE_SWAP_OUTPUT_INTERNAL_EXTERNAL})
	}

	if _, err = sw.dataplane.TableCreate(ctx, trapToHostifTableReq); err != nil {
		return nil, err
	}

	_, err = sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: inputIfaceTable}},
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT, // TODO: Figure out all the ways ports can be mapped to interfaces.
						},
					}},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: outputIfaceTable}},
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE,
						},
					}},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	tunnel := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: TunnelEncap}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, tunnel); err != nil {
		return nil, err
	}
	_, err = sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: hostifToPortTable}},
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_HOST_PORT_ID,
						},
					}},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	_, err = sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: portToHostifTable}},
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT,
						},
					}},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	_, err = sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
			TableType: fwdpb.TableType_TABLE_TYPE_FLOW,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
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
		Switch:       swID,
		VlanId:       proto.Uint32(DefaultVlanId), // Create the default VLAN.
		LearnDisable: proto.Bool(true),            // TODO: figure out what does this do?
	})
	if err != nil {
		return nil, err
	}
	sw.mgr.StoreAttributes(swID, &saipb.SwitchAttribute{DefaultVlanId: &vlanResp.Oid})
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

func (sw *saiSwitch) SetSwitchAttribute(ctx context.Context, req *saipb.SetSwitchAttributeRequest) (*saipb.SetSwitchAttributeResponse, error) {
	switch {
	case req.PreIngressAcl != nil:
		if err := sw.bindACLTable(ctx, fmt.Sprint(req.GetPreIngressAcl()), PreIngressActionTable); err != nil {
			return nil, err
		}
	case req.IngressAcl != nil:
		if err := sw.bindACLTable(ctx, fmt.Sprint(req.GetIngressAcl()), IngressActionTable); err != nil {
			return nil, err
		}
	case req.EgressAcl != nil:
		if err := sw.bindACLTable(ctx, fmt.Sprint(req.GetEgressAcl()), EgressActionTable); err != nil {
			return nil, err
		}
	}
	return &saipb.SetSwitchAttributeResponse{}, nil
}

func (sw *saiSwitch) bindACLTable(ctx context.Context, aclTableID, stageID string) error {
	_, err := sw.dataplane.TableEntryAdd(ctx, &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: stageID}},
		Actions:   []*fwdpb.ActionDesc{fwdconfig.Action(fwdconfig.LookupAction(aclTableID)).Build()},
		EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Action{Action: &fwdpb.ActionEntryDesc{Id: "acl", InsertMethod: fwdpb.ActionEntryDesc_INSERT_METHOD_PREPEND}}},
	})
	return err
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
			switch oType {
			case saipb.ObjectType_OBJECT_TYPE_PORT:
			case saipb.ObjectType_OBJECT_TYPE_BRIDGE_PORT:
			case saipb.ObjectType_OBJECT_TYPE_LAG:
			default:
				log.Infof("skipping port state event for type %v", oType)
				continue
			}
			status := saipb.PortOperStatus_PORT_OPER_STATUS_UNKNOWN
			if ed.GetPort().PortInfo.OperStatus == fwdpb.PortState_PORT_STATE_ENABLED_UP {
				status = saipb.PortOperStatus_PORT_OPER_STATUS_UP
			} else if ed.GetPort().PortInfo.OperStatus == fwdpb.PortState_PORT_STATE_DISABLED_DOWN {
				status = saipb.PortOperStatus_PORT_OPER_STATUS_DOWN
			}
			resp := &saipb.PortStateChangeNotificationResponse{
				Data: []*saipb.PortOperStatusNotification{{
					PortId:    uint64(num),
					PortState: status,
				}},
			}
			log.Infof("send port event: %+v", resp)
			err = srv.Send(resp)
			if err != nil {
				return err
			}
		}
	}
}

func (sw *saiSwitch) Reset() {
	sw.port.Reset()
	sw.hostif.Reset()
}

// createFIBSelector creates a table that controls which forwarding table is used.
func createFIBSelector(ctx context.Context, id string, c switchDataplaneAPI) error {
	fieldID := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION,
		},
	}

	ipVersion := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: id},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: FIBSelectorTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{fieldID},
				},
			},
		},
	}
	if _, err := c.TableCreate(ctx, ipVersion); err != nil {
		return err
	}
	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: id},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: FIBSelectorTable,
			},
		},
		Entries: []*fwdpb.TableEntryAddRequest_Entry{{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes:   []byte{0x4},
							FieldId: fieldID,
						}},
					},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
				Action: &fwdpb.ActionDesc_Lookup{
					Lookup: &fwdpb.LookupActionDesc{
						TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: FIBV4Table}},
					},
				},
			}},
		}, {
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes:   []byte{0x6},
							FieldId: fieldID,
						}},
					},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
				Action: &fwdpb.ActionDesc_Lookup{
					Lookup: &fwdpb.LookupActionDesc{
						TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: FIBV6Table}},
					},
				},
			}},
		}},
	}
	if _, err := c.TableEntryAdd(ctx, entries); err != nil {
		return err
	}
	return nil
}
