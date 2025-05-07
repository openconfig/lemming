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
	"log/slog"
	"net"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

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
	opts            *dplaneopts.Options
	acl             *acl
	buffer          *buffer
	port            *port
	vlan            *vlan
	stp             *stp
	bridge          *bridge
	hostif          *hostif
	hash            *hash
	isolationGroup  *isolationGroup
	l2mc            *l2mc
	l2mcGroup       *l2mcGroup
	myMac           *myMac
	neighbor        *neighbor
	nextHopGroup    *nextHopGroup
	nextHop         *nextHop
	policer         *policer
	route           *route
	lag             *lag
	tunnel          *tunnel
	queue           *queue
	sg              *schedulerGroup
	routerInterface *routerInterface
	virtualRouter   *virtualRouter
	udf             *udf
	scheduler       *scheduler
	qosMap          *qosMap
	rpf             *rpfGroup
	wred            *wred
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

type luciusTrace struct {
	switchDataplaneAPI
	tracer trace.Tracer
}

func (l *luciusTrace) TableCreate(ctx context.Context, req *fwdpb.TableCreateRequest) (*fwdpb.TableCreateReply, error) {
	ctx, span := l.tracer.Start(ctx, "TableCreate")
	defer span.End()
	return l.switchDataplaneAPI.TableCreate(ctx, req)
}

func (l *luciusTrace) TableEntryAdd(ctx context.Context, req *fwdpb.TableEntryAddRequest) (*fwdpb.TableEntryAddReply, error) {
	ctx, span := l.tracer.Start(ctx, "TableEntryAdd")
	defer span.End()
	return l.switchDataplaneAPI.TableEntryAdd(ctx, req)
}

func (l *luciusTrace) TableEntryRemove(ctx context.Context, req *fwdpb.TableEntryRemoveRequest) (*fwdpb.TableEntryRemoveReply, error) {
	ctx, span := l.tracer.Start(ctx, "TableEntryRemove")
	defer span.End()
	return l.switchDataplaneAPI.TableEntryRemove(ctx, req)
}

func (l *luciusTrace) PortState(ctx context.Context, req *fwdpb.PortStateRequest) (*fwdpb.PortStateReply, error) {
	ctx, span := l.tracer.Start(ctx, "PortState")
	defer span.End()
	return l.switchDataplaneAPI.PortState(ctx, req)
}

func (l *luciusTrace) ObjectCounters(ctx context.Context, req *fwdpb.ObjectCountersRequest) (*fwdpb.ObjectCountersReply, error) {
	ctx, span := l.tracer.Start(ctx, "ObjectCounters")
	defer span.End()
	return l.switchDataplaneAPI.ObjectCounters(ctx, req)
}

func (l *luciusTrace) PortCreate(ctx context.Context, req *fwdpb.PortCreateRequest) (*fwdpb.PortCreateReply, error) {
	ctx, span := l.tracer.Start(ctx, "PortCreate")
	defer span.End()
	return l.switchDataplaneAPI.PortCreate(ctx, req)
}

func (l *luciusTrace) PortUpdate(ctx context.Context, req *fwdpb.PortUpdateRequest) (*fwdpb.PortUpdateReply, error) {
	ctx, span := l.tracer.Start(ctx, "PortUpdate")
	defer span.End()
	return l.switchDataplaneAPI.PortUpdate(ctx, req)
}

func (l *luciusTrace) AttributeUpdate(ctx context.Context, req *fwdpb.AttributeUpdateRequest) (*fwdpb.AttributeUpdateReply, error) {
	ctx, span := l.tracer.Start(ctx, "AttributeUpdate")
	defer span.End()
	return l.switchDataplaneAPI.AttributeUpdate(ctx, req)
}

func (l *luciusTrace) ObjectNID(ctx context.Context, req *fwdpb.ObjectNIDRequest) (*fwdpb.ObjectNIDReply, error) {
	ctx, span := l.tracer.Start(ctx, "ObjectNID")
	defer span.End()
	return l.switchDataplaneAPI.ObjectNID(ctx, req)
}

func (l *luciusTrace) ObjectDelete(ctx context.Context, req *fwdpb.ObjectDeleteRequest) (*fwdpb.ObjectDeleteReply, error) {
	ctx, span := l.tracer.Start(ctx, "ObjectDelete")
	defer span.End()
	return l.switchDataplaneAPI.ObjectDelete(ctx, req)
}

func (l *luciusTrace) FlowCounterCreate(ctx context.Context, req *fwdpb.FlowCounterCreateRequest) (*fwdpb.FlowCounterCreateReply, error) {
	ctx, span := l.tracer.Start(ctx, "FlowCounterCreate")
	defer span.End()
	return l.switchDataplaneAPI.FlowCounterCreate(ctx, req)
}

func (l *luciusTrace) FlowCounterQuery(ctx context.Context, req *fwdpb.FlowCounterQueryRequest) (*fwdpb.FlowCounterQueryReply, error) {
	ctx, span := l.tracer.Start(ctx, "FlowCounterQuery")
	defer span.End()
	return l.switchDataplaneAPI.FlowCounterQuery(ctx, req)
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
	L2MCGroupTable        = "l2mcg"
	policerTabler         = "policerTable"
	invalidIngress        = "invalid-ingress"
	invalidIngressV4Table = "invalid-ingress-v4"
	invalidIngressV6Table = "invalid-ingress-v6"
	outputTable           = "output-table"
	DefaultVlanId         = 1
)

func newSwitch(mgr *attrmgr.AttrMgr, engine switchDataplaneAPI, s *grpc.Server, opts *dplaneopts.Options) (*saiSwitch, error) {
	dplane := &luciusTrace{switchDataplaneAPI: engine, tracer: otel.Tracer("lucius")}

	vlan := newVlan(mgr, dplane, s)
	q := newQueue(mgr, dplane, s)
	sg := newSchedulerGroup(mgr, dplane, s)
	port, err := newPort(mgr, dplane, s, vlan, q, sg, opts)
	if err != nil {
		return nil, err
	}

	sw := &saiSwitch{
		dataplane:       dplane,
		opts:            opts,
		acl:             newACL(mgr, dplane, s),
		policer:         newPolicer(mgr, dplane, s),
		port:            port,
		vlan:            vlan,
		stp:             &stp{},
		bridge:          newBridge(mgr, engine, s),
		hostif:          newHostif(mgr, engine, s, opts),
		hash:            newHash(mgr, engine, s),
		isolationGroup:  newIsolationGroup(mgr, engine, s),
		l2mc:            newL2mc(mgr, engine, s),
		l2mcGroup:       newL2mcGroup(mgr, engine, s),
		myMac:           newMyMac(mgr, engine, s, opts),
		neighbor:        newNeighbor(mgr, engine, s),
		nextHopGroup:    newNextHopGroup(mgr, engine, s),
		nextHop:         newNextHop(mgr, engine, s),
		route:           newRoute(mgr, engine, s),
		routerInterface: newRouterInterface(mgr, engine, s),
		lag:             newLAG(mgr, engine, s),
		tunnel:          newTunnel(mgr, engine, s),
		udf:             newUdf(mgr, engine, s),
		scheduler:       newScheduler(mgr, engine, s),
		qosMap:          newQOSMap(mgr, engine, s),
		virtualRouter:   newVirtualRouter(mgr, engine, s),
		rpf:             newRpfGroup(mgr, engine, s),
		buffer:          newBuffer(mgr, engine, s),
		wred:            newWRED(mgr, engine, s),
		queue:           q,
		sg:              sg,
		mgr:             mgr,
	}
	saipb.RegisterSwitchServer(s, sw)
	saipb.RegisterStpServer(s, sw.stp)
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
			Actions:   []*fwdpb.ActionDesc{fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_BIT_WRITE, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBitOp(1, 0).WithValue([]byte{0})).Build()},
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
			Actions:   []*fwdpb.ActionDesc{fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_BIT_WRITE, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBitOp(1, 0).WithValue([]byte{0})).Build()},
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
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
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
	l2mcGroupReq := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: L2MCGroupTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L2MC_GROUP_ID,
						},
					}},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, l2mcGroupReq); err != nil {
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

	if !sw.opts.SkipIPValidation {
		if err := sw.createInvalidPacketFilter(ctx); err != nil {
			return nil, err
		}
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
			Actions:   []*fwdpb.ActionDesc{fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_BIT_WRITE, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBitOp(1, 0).WithValue([]byte{0})).Build()},
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
			Actions:   []*fwdpb.ActionDesc{fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_BIT_WRITE, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBitOp(1, 0).WithValue([]byte{0})).Build()},
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
	if err := sw.createFIBSelector(ctx); err != nil {
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
	_, err := sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
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
	_, err = sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_FLOW,
			// Actions:   []*fwdpb.ActionDesc{fwdconfig.Action(fwdconfig.LookupAction(portToHostifTable)).Build()},
			TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: trapTableID}},
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
			fwdconfig.LookupAction(trapTableID)).
		Build(),
	)
	if err != nil {
		return nil, err
	}

	trapToHostifTableReq := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: trapIDToHostifTable}},
			Actions:   []*fwdpb.ActionDesc{fwdconfig.Action(fwdconfig.LookupAction(portToHostifTable)).Build()},
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
					}, {
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VLAN_TAG, // TODO: Figure out all the ways ports can be mapped to interfaces.
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

	_, err = sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: policerTabler}},
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_POLICER_ID,
						},
					}},
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

	err = sw.createOutputTable(ctx, fmt.Sprint(cpuPortID))
	if err != nil {
		return nil, err
	}

	myMAC := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_FLOW,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: MyMacTable}},
			Actions:   getL2Pipeline(),
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
	vrResp, err := attrmgr.InvokeAndSave(ctx, sw.mgr, sw.virtualRouter.CreateVirtualRouter, &saipb.CreateVirtualRouterRequest{
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
		NativeHashFieldList: []saipb.NativeHashField{
			saipb.NativeHashField_NATIVE_HASH_FIELD_SRC_IP,
			saipb.NativeHashField_NATIVE_HASH_FIELD_DST_IP,
			saipb.NativeHashField_NATIVE_HASH_FIELD_L4_SRC_PORT,
			saipb.NativeHashField_NATIVE_HASH_FIELD_L4_DST_PORT,
		},
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
		AclEntryMaximumPriority:          proto.Uint32(10000),
		AclTableMinimumPriority:          proto.Uint32(1),
		AclTableMaximumPriority:          proto.Uint32(10000),
		MaxAclActionCount:                proto.Uint32(1000),
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

// Set up rules to drop packets that contain invalid IP or ttl == 0.
// https://www.rfc-editor.org/rfc/rfc1812#section-5.3.7
func (sw *saiSwitch) createInvalidPacketFilter(ctx context.Context) error {
	ips := map[string]map[fwdpb.PacketFieldNum][]string{
		invalidIngressV4Table: { /*                               LOOPBACK      BROADCAST             MULTICAST */
			fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC: {"127.0.0.0/8", "255.255.255.255/32", "224.0.0.0/4"},
			fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST: {"127.0.0.0/8", "255.255.255.255/32"},
		},
		invalidIngressV6Table: { /*                              LOOPBACK  MULTICAST*/
			fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC: {"::1/128", "ff00::/8"},
			fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST: {"::1/128"},
		},
	}
	// Packets can't have multicast, or loopback IP as the source IP.
	// Only unicast MAC address are processed at this stage, so multicast IPs are invalid
	for table, ipsByField := range ips {
		_, err := sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
			ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
			Desc: &fwdpb.TableDesc{
				Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
				TableType: fwdpb.TableType_TABLE_TYPE_FLOW,
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: table}},
				Table: &fwdpb.TableDesc_Flow{
					Flow: &fwdpb.FlowTableDesc{
						BankCount: 1,
					},
				},
			},
		})
		if err != nil {
			return err
		}

		for field, ips := range ipsByField {
			for _, ip := range ips {
				_, prefix, err := net.ParseCIDR(ip)
				if err != nil {
					return err
				}
				req := fwdconfig.TableEntryAddRequest(sw.dataplane.ID(), table).
					AppendEntry(
						fwdconfig.EntryDesc(fwdconfig.FlowEntry(fwdconfig.PacketFieldMaskedBytes(field).WithBytes(prefix.IP, prefix.Mask))),
						fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_BIT_WRITE, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBitOp(1, 0).WithValue([]byte{0}),
					).Build()
				if _, err := sw.dataplane.TableEntryAdd(ctx, req); err != nil {
					return err
				}
			}
		}
		// Before the TTL is decremented and after the packets may be punted, drop packet with TTL == 1 or TTL == 0.
		req := fwdconfig.TableEntryAddRequest(sw.dataplane.ID(), table).
			AppendEntry(
				fwdconfig.EntryDesc(fwdconfig.FlowEntry(fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP).WithBytes([]byte{0x00}, []byte{0xFF}))),
				fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_BIT_WRITE, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBitOp(1, 0).WithValue([]byte{0}),
			).Build()
		if _, err := sw.dataplane.TableEntryAdd(ctx, req); err != nil {
			return err
		}
		req = fwdconfig.TableEntryAddRequest(sw.dataplane.ID(), table).
			AppendEntry(
				fwdconfig.EntryDesc(fwdconfig.FlowEntry(fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP).WithBytes([]byte{0x01}, []byte{0xFF}))),
				fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_BIT_WRITE, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBitOp(1, 0).WithValue([]byte{0}),
			).Build()
		if _, err := sw.dataplane.TableEntryAdd(ctx, req); err != nil {
			return err
		}
	}

	_, err := sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: invalidIngress}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{Field: &fwdpb.PacketField{
						FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION,
					}}},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	verReq := fwdconfig.TableEntryAddRequest(sw.dataplane.ID(), invalidIngress).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{4}))),
		fwdconfig.LookupAction(invalidIngressV4Table),
	).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).WithBytes([]byte{6}))),
		fwdconfig.LookupAction(invalidIngressV6Table),
	).Build()

	if _, err := sw.dataplane.TableEntryAdd(ctx, verReq); err != nil {
		return err
	}

	return nil
}

func (sw *saiSwitch) createOutputTable(ctx context.Context, cpuPortID string) error {
	_, err := sw.dataplane.TableCreate(ctx, &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: outputTable}},
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION,
						},
					}},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	req := fwdconfig.TableEntryAddRequest(sw.dataplane.ID(), outputTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBytes([]byte{0}))), // DROP
		fwdconfig.DropAction(),
	).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBytes([]byte{1}))), // FORWARD
		fwdconfig.DecapAction(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET),                                                                           // Decap L2 header.
		fwdconfig.LookupAction(NHActionTable), // Apply additional encap actions
		fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_DEC, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP).WithValue([]byte{0x1}), // Decrement TTL.
		fwdconfig.EncapAction(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET),                                                         // Encap L2 header.
		fwdconfig.LookupAction(NeighborTable), // Lookup in the neighbor table.
		fwdconfig.LookupAction(SRCMACTable),   // Update source mac
	).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBytes([]byte{2}))), // COPY AND DROP
		fwdconfig.TransmitAction(cpuPortID),
	).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ACTION).WithBytes([]byte{3}))), // COPY AND FORWARD
		fwdconfig.MirrorAction().WithPort(cpuPortID, fwdpb.PortAction_PORT_ACTION_OUTPUT).WithFields(fwdconfig.PacketFieldIDField(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TRAP_ID, 0), fwdconfig.PacketFieldIDField(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TARGET_EGRESS_PORT, 0)),
		fwdconfig.DecapAction(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET),                                                         // Decap L2 header.
		fwdconfig.LookupAction(NHActionTable),                                                                                         // Apply additional encap actions
		fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_DEC, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP).WithValue([]byte{0x1}), // Decrement TTL.
		fwdconfig.EncapAction(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET),                                                         // Encap L2 header.
		fwdconfig.LookupAction(NeighborTable),                                                                                         // Lookup in the neighbor table.
		fwdconfig.LookupAction(SRCMACTable),                                                                                           // Update source mac
	)
	if _, err := sw.dataplane.TableEntryAdd(ctx, req.Build()); err != nil {
		return err
	}

	return nil
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
				slog.WarnContext(srv.Context(), "couldn't get numeric port id", "err", err)
				continue
			}
			oType := sw.mgr.GetType(ed.GetPort().GetPortId().GetObjectId().GetId())
			switch oType {
			case saipb.ObjectType_OBJECT_TYPE_PORT:
			case saipb.ObjectType_OBJECT_TYPE_BRIDGE_PORT:
			case saipb.ObjectType_OBJECT_TYPE_LAG:
			default:
				slog.InfoContext(srv.Context(), "skipping port state event", "type", oType)
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
			slog.InfoContext(srv.Context(), "send port event", "event", resp)
			err = srv.Send(resp)
			if err != nil {
				return err
			}
		}
	}
}

func (sw *saiSwitch) Reset() {
	sw.vlan.Reset()
	sw.port.Reset()
	sw.hostif.Reset()
}

// createFIBSelector creates a table that controls which forwarding table is used.
func (sw *saiSwitch) createFIBSelector(ctx context.Context) error {
	fieldID := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION,
		},
	}

	ipVersion := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: FIBSelectorTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{fieldID},
				},
			},
		},
	}
	if _, err := sw.dataplane.TableCreate(ctx, ipVersion); err != nil {
		return err
	}
	v4Acts := []*fwdpb.ActionDesc{{
		ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
		Action: &fwdpb.ActionDesc_Lookup{
			Lookup: &fwdpb.LookupActionDesc{
				TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: FIBV4Table}},
			},
		},
	}}
	v6Acts := []*fwdpb.ActionDesc{{
		ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
		Action: &fwdpb.ActionDesc_Lookup{
			Lookup: &fwdpb.LookupActionDesc{
				TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: FIBV6Table}},
			},
		},
	}}

	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: sw.dataplane.ID()},
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
			Actions: v4Acts,
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
			Actions: v6Acts,
		}},
	}
	if _, err := sw.dataplane.TableEntryAdd(ctx, entries); err != nil {
		return err
	}
	return nil
}
