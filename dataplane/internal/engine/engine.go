// Copyright 2022 Google LLC
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

package engine

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/openconfig/lemming/dataplane/forwarding"
	"github.com/openconfig/lemming/dataplane/forwarding/attributes"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"

	log "github.com/golang/glog"

	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	fibV4Table       = "fib-v4"
	fibV6Table       = "fib-v6"
	srcMACTable      = "port-mac"
	fibSelectorTable = "fib-selector"
	neighborTable    = "neighbor"
	nhgTable         = "nhg-table"
	nhTable          = "nh-table"
	layer2PuntTable  = "layer2-punt"
	layer3PuntTable  = "layer3-punt"
	arpPuntTable     = "arp-punt"
)

// Engine contains a routing context and methods to manage it.
type Engine struct {
	dpb.UnimplementedDataplaneServer
	*forwarding.Server
	id        string
	idToNIDMu sync.RWMutex
	// idToNID is map from RPC ID (proto), to internal object NID.
	idToNID   map[string]uint64
	nextNHGID atomic.Uint64
	nextNHID  atomic.Uint64
}

// New creates a new engine and sets up the forwarding tables.
func New(ctx context.Context) (*Engine, error) {
	e := &Engine{
		id:      "lucius",
		Server:  forwarding.New("engine"),
		idToNID: map[string]uint64{},
	}

	_, err := e.Server.ContextCreate(context.Background(), &fwdpb.ContextCreateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
	})
	if err != nil {
		return nil, err
	}

	v4FIB := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_PREFIX,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibV4Table}},
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
	if _, err := e.Server.TableCreate(ctx, v4FIB); err != nil {
		return nil, err
	}
	v6FIB := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_PREFIX,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibV6Table}},
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
	if _, err := e.Server.TableCreate(ctx, v6FIB); err != nil {
		return nil, err
	}
	portMAC := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: srcMACTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT,
						},
					}},
				},
			},
		},
	}
	if _, err := e.Server.TableCreate(ctx, portMAC); err != nil {
		return nil, err
	}
	neighbor := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: neighborTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP,
						},
					}},
				},
			},
		},
	}
	if _, err := e.Server.TableCreate(ctx, neighbor); err != nil {
		return nil, err
	}
	nh := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: nhTable}},
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
	if _, err := e.Server.TableCreate(ctx, nh); err != nil {
		return nil, err
	}
	nhg := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: nhgTable}},
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
	if _, err := e.Server.TableCreate(ctx, nhg); err != nil {
		return nil, err
	}
	if err := createFIBSelector(ctx, e.id, e.Server); err != nil {
		return nil, err
	}
	if err := createLayer2PuntTable(ctx, e.id, e.Server); err != nil {
		return nil, err
	}
	if err := createLayer3PuntTable(ctx, e.id, e.Server); err != nil {
		return nil, err
	}
	return e, nil
}

// ID returns the engine's forwarding context id.
func (e *Engine) ID() string {
	return e.id
}

func (e *Engine) CreatePort(ctx context.Context, req *dpb.CreatePortRequest) (*dpb.CreatePortResponse, error) {
	var err error
	switch req.Type {
	case fwdpb.PortType_PORT_TYPE_KERNEL:
		err = e.CreateExternalPort(ctx, req.GetId(), req.GetKernelDev())
	case fwdpb.PortType_PORT_TYPE_TAP:
		err = e.CreateLocalPort(ctx, req.GetId(), req.GetKernelDev(), req.GetExternalPort())
	}
	return &dpb.CreatePortResponse{}, err
}

// AddLayer3PuntRule adds rule to output packets to a corresponding port based on the destination IP and input port.
func (e *Engine) AddLayer3PuntRule(ctx context.Context, portName string, ip []byte) error {
	e.idToNIDMu.Lock()
	defer e.idToNIDMu.Unlock()
	portID := e.idToNID[portName]

	nidBytes := make([]byte, binary.Size(portID))
	binary.BigEndian.PutUint64(nidBytes, portID)

	log.Infof("adding layer3 punt rule: portName %s, id %d, ip %v", portName, portID, ip)

	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: layer3PuntTable,
			},
		},
		Entries: []*fwdpb.TableEntryAddRequest_Entry{{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes: nidBytes,
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT,
								},
							},
						}, {
							Bytes: ip,
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
								},
							},
						}},
					},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_SWAP_OUTPUT_INTERNAL_EXTERNAL,
			}, {
				ActionType: fwdpb.ActionType_ACTION_TYPE_OUTPUT,
			}},
		}},
	}
	if _, err := e.Server.TableEntryAdd(ctx, entries); err != nil {
		return err
	}
	return nil
}

// prefixToPrimitives returns the primitive types of the route prefix.
// ip addr bytes, ip mask bytes, is ipv4, vrf id, error.
func prefixToPrimitives(prefix *dpb.RoutePrefix) ([]byte, []byte, bool, uint64, error) {
	var ip []byte
	var mask []byte
	var isIPv4 bool
	vrf := prefix.GetVrfId()

	switch pre := prefix.GetPrefix().(type) {
	case *dpb.RoutePrefix_Cidr:
		_, ipNet, err := net.ParseCIDR(pre.Cidr)
		if err != nil {
			return ip, mask, isIPv4, vrf, fmt.Errorf("failed to parse ip prefix: %v", err)
		}
		ip = ipNet.IP.To4()
		mask = ipNet.Mask
		isIPv4 = true
		if ip == nil {
			ip = ipNet.IP.To16()
			mask = ipNet.Mask
			isIPv4 = false
		}
	case *dpb.RoutePrefix_Mask:
		ip = pre.Mask.Addr
		mask = pre.Mask.Mask
		switch len(ip) {
		case net.IPv4len:
			isIPv4 = true
		case net.IPv6len:
			isIPv4 = false
		default:
			return ip, mask, isIPv4, vrf, fmt.Errorf("invalid ip addr length")
		}
	default:
		return ip, mask, isIPv4, vrf, fmt.Errorf("invalid prefix type")
	}
	return ip, mask, isIPv4, vrf, nil
}

// addNextHopList creates all the next hops from the message, then create a next hop group if there are multiple next hops.
func (e *Engine) addNextHopList(ctx context.Context, nhg *dpb.NextHopList) ([]*fwdpb.ActionDesc, error) {
	if len(nhg.GetHops()) == 1 {
		nhID := e.nextNHID.Add(1)
		if err := e.addNextHop(ctx, nhID, nhg.Hops[0]); err != nil {
			return nil, err
		}
		return []*fwdpb.ActionDesc{
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID).WithUint64Value(nhID)).Build(),
			fwdconfig.Action(fwdconfig.LookupAction(nhTable)).Build(),
		}, nil
	}

	idList := &dpb.NextHopIDList{}
	for _, hop := range nhg.GetHops() {
		nhID := e.nextNHID.Add(1)
		if err := e.addNextHop(ctx, nhID, hop); err != nil {
			return nil, err
		}
		idList.Hops = append(idList.Hops, nhID)
		idList.Weights = append(idList.Weights, hop.Weight)
	}
	nhgID := e.nextNHGID.Add(1)
	if err := e.addNextHopGroupIDList(ctx, nhgID, idList); err != nil {
		return nil, err
	}
	return []*fwdpb.ActionDesc{
		fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_GROUP_ID).WithUint64Value(nhgID)).Build(),
		fwdconfig.Action(fwdconfig.LookupAction(nhgTable)).Build(),
	}, nil
}

// addNextHopGroupIDList adds an entry to the next hop group table.
func (e *Engine) addNextHopGroupIDList(ctx context.Context, id uint64, nhg *dpb.NextHopIDList) error {
	var actLists []*fwdpb.ActionList
	for i, nh := range nhg.GetHops() {
		actLists = append(actLists, &fwdpb.ActionList{
			Weight: nhg.Weights[i],
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
				Action: &fwdpb.ActionDesc_Update{
					Update: &fwdpb.UpdateActionDesc{
						FieldId: &fwdpb.PacketFieldId{
							Field: &fwdpb.PacketField{
								FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID,
							},
						},
						Type:  fwdpb.UpdateType_UPDATE_TYPE_SET,
						Value: binary.BigEndian.AppendUint64(nil, nh),
					},
				},
			}},
		})
	}
	// If there are multiple next-hops, configure the route to use ECMP or WCMP.
	actions := []*fwdpb.ActionDesc{{
		ActionType: fwdpb.ActionType_ACTION_TYPE_SELECT_ACTION_LIST,
		Action: &fwdpb.ActionDesc_Select{
			Select: &fwdpb.SelectActionListActionDesc{
				SelectAlgorithm: fwdpb.SelectActionListActionDesc_SELECT_ALGORITHM_CRC32, // TODO: should algo + hash be configurable?
				FieldIds: []*fwdpb.PacketFieldId{{Field: &fwdpb.PacketField{ // Hash the traffic flow, identified, IP protocol, L3 SRC, DST address, and L4 ports (if present).
					FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO,
				}}, {Field: &fwdpb.PacketField{
					FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC,
				}}, {Field: &fwdpb.PacketField{
					FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
				}}, {Field: &fwdpb.PacketField{
					FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC,
				}}, {Field: &fwdpb.PacketField{
					FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST,
				}}},
				ActionLists: actLists,
			},
		},
	}, {
		ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
		Action: &fwdpb.ActionDesc_Lookup{
			Lookup: &fwdpb.LookupActionDesc{
				TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{
					Id: nhTable,
				}},
			},
		},
	}}

	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: nhgTable,
			},
		},
		Entries: []*fwdpb.TableEntryAddRequest_Entry{{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes: binary.BigEndian.AppendUint64(nil, id),
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_GROUP_ID,
								},
							},
						}},
					},
				},
			},
			Actions: actions,
		}},
	}
	if _, err := e.Server.TableEntryAdd(ctx, entries); err != nil {
		return err
	}

	return nil
}

// addNextHop adds an entry to the next hop table.
// TODO: Remove workaround that nexthop IP is not specified that the packet is treated as directly connected.
func (e *Engine) addNextHop(ctx context.Context, id uint64, nh *dpb.NextHop) error {
	var nextHopIP []byte
	if nhIPStr := nh.GetIp(); nhIPStr != "" {
		nextHop := net.ParseIP(nhIPStr)
		nextHopIP = nextHop.To4()
		if nextHopIP == nil {
			nextHopIP = nextHop.To16()
		}
	}
	// Set the next hop IP in the packet's metadata.
	nextHopAct := fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP).WithValue(nextHopIP)).Build()
	if nextHopIP == nil {
		nextHopAct = fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_COPY, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP).WithFieldSrc(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST)).Build()
	}
	// Set the output port of the packet.
	transmitAct := fwdconfig.Action(fwdconfig.TransmitAction(nh.GetPort())).Build()

	acts := append([]*fwdpb.ActionDesc{nextHopAct, transmitAct}, nh.GetPreTransmitActions()...)
	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: nhTable,
			},
		},
		Entries: []*fwdpb.TableEntryAddRequest_Entry{{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes: binary.BigEndian.AppendUint64(nil, id),
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID,
								},
							},
						}},
					},
				},
			},
			Actions: acts,
		}},
	}
	if _, err := e.Server.TableEntryAdd(ctx, entries); err != nil {
		return err
	}

	return nil
}

func (e *Engine) actionsFromRoute(ctx context.Context, route *dpb.Route) ([]*fwdpb.ActionDesc, error) {
	// If action is DROP, then skip handling next hops.
	if route.GetAction() == dpb.PacketAction_PACKET_ACTION_DROP {
		return []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}}, nil
	}

	var actions []*fwdpb.ActionDesc
	switch hop := route.GetHop().(type) {
	case *dpb.Route_PortId:
		actions = []*fwdpb.ActionDesc{
			// Set the next hop IP in the packet's metadata.
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_COPY, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP).WithFieldSrc(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST)).Build(),
			// Set the output port.
			fwdconfig.Action(fwdconfig.TransmitAction(hop.PortId)).Build(),
		}
	case *dpb.Route_NextHopId:
		actions = []*fwdpb.ActionDesc{ // Set the next hop ID in the packet's metadata.
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID).WithUint64Value(hop.NextHopId)).Build(),
			fwdconfig.Action(fwdconfig.LookupAction(nhTable)).Build(),
		}
	case *dpb.Route_NextHopGroupId:
		actions = []*fwdpb.ActionDesc{ // Set the next hop group ID in the packet's metadata.
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_GROUP_ID).WithUint64Value(hop.NextHopGroupId)).Build(),
			fwdconfig.Action(fwdconfig.LookupAction(nhgTable)).Build(),
		}
	case *dpb.Route_NextHops:
		var err error
		actions, err = e.addNextHopList(ctx, hop.NextHops)
		if err != nil {
			return nil, err
		}
	}
	return actions, nil
}

// AddIPRoute adds a route to the FIB. It operates in two modes:
// 1. Client-managed IDs: each next hop and next hop group must be created before adding to a route with user provided ids.
// 2. Server-managed IDs: each next hop and next hop group must be specified with route. The server implicitly creates ids.
// TODO: Enforce that only one mode can be used.
func (e *Engine) AddIPRoute(ctx context.Context, req *dpb.AddIPRouteRequest) (*dpb.AddIPRouteResponse, error) {
	ip, mask, isIPv4, vrf, err := prefixToPrimitives(req.GetRoute().GetPrefix())
	if err != nil {
		return nil, err
	}
	fib := fibV6Table
	if isIPv4 {
		fib = fibV4Table
	}

	entry := &fwdpb.TableEntryAddRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fib}},
		ContextId: &fwdpb.ContextId{Id: e.id},
		EntryDesc: &fwdpb.EntryDesc{
			Entry: &fwdpb.EntryDesc_Prefix{
				Prefix: &fwdpb.PrefixEntryDesc{
					Fields: []*fwdpb.PacketFieldMaskedBytes{{
						FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF}},
						Bytes:   binary.BigEndian.AppendUint64(nil, vrf),
					}, {
						FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
						Bytes:   ip,
						Masks:   mask,
					}},
				},
			},
		},
	}

	entry.Actions, err = e.actionsFromRoute(ctx, req.GetRoute())
	if err != nil {
		return nil, err
	}

	if _, err := e.Server.TableEntryAdd(ctx, entry); err != nil {
		return nil, err
	}

	return &dpb.AddIPRouteResponse{}, nil
}

// Remove deletes a route from the FIB.
// TODO: Clean up orphaned next-hop and next-hop-groups for server managed ids.
func (e *Engine) RemoveIPRoute(ctx context.Context, req *dpb.RemoveIPRouteRequest) (*dpb.RemoveIPRouteResponse, error) {
	ip, mask, isIPv4, vrf, err := prefixToPrimitives(req.GetPrefix())
	if err != nil {
		return nil, err
	}
	fib := fibV6Table
	if isIPv4 {
		fib = fibV4Table
	}
	entry := &fwdpb.TableEntryRemoveRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fib}},
		ContextId: &fwdpb.ContextId{Id: e.id},
		EntryDesc: &fwdpb.EntryDesc{
			Entry: &fwdpb.EntryDesc_Prefix{
				Prefix: &fwdpb.PrefixEntryDesc{
					Fields: []*fwdpb.PacketFieldMaskedBytes{{
						FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF}},
						Bytes:   binary.BigEndian.AppendUint64(nil, vrf),
					}, {
						FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
						Bytes:   ip,
						Masks:   mask,
					}},
				},
			},
		},
	}
	if _, err := e.Server.TableEntryRemove(ctx, entry); err != nil {
		return nil, err
	}
	return &dpb.RemoveIPRouteResponse{}, nil
}

// AddNeighbor adds a neighbor to the neighbor table.
func (e *Engine) AddNeighbor(ctx context.Context, ip, mac []byte) error {
	entry := &fwdpb.TableEntryAddRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: neighborTable}},
		ContextId: &fwdpb.ContextId{Id: e.id},
		EntryDesc: &fwdpb.EntryDesc{
			Entry: &fwdpb.EntryDesc_Exact{
				Exact: &fwdpb.ExactEntryDesc{
					Fields: []*fwdpb.PacketFieldBytes{{
						FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP}},
						Bytes:   ip,
					}},
				},
			},
		},
		Actions: []*fwdpb.ActionDesc{{ // Set the dst MAC.
			ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
			Action: &fwdpb.ActionDesc_Update{
				Update: &fwdpb.UpdateActionDesc{
					FieldId: &fwdpb.PacketFieldId{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST,
						},
					},
					Type:  fwdpb.UpdateType_UPDATE_TYPE_SET,
					Value: mac,
				},
			},
		}},
	}
	if _, err := e.Server.TableEntryAdd(ctx, entry); err != nil {
		return err
	}

	return nil
}

// RemoveNeighbor removes a neighbor from the neighbor table.
func (e *Engine) RemoveNeighbor(ctx context.Context, ip []byte) error {
	entry := &fwdpb.TableEntryRemoveRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: neighborTable}},
		ContextId: &fwdpb.ContextId{Id: e.id},
		EntryDesc: &fwdpb.EntryDesc{
			Entry: &fwdpb.EntryDesc_Exact{
				Exact: &fwdpb.ExactEntryDesc{
					Fields: []*fwdpb.PacketFieldBytes{{
						FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP}},
						Bytes:   ip,
					}},
				},
			},
		},
	}
	if _, err := e.Server.TableEntryRemove(ctx, entry); err != nil {
		return err
	}

	return nil
}

// UpdatePortSrcMAC updates a port's source mac address.
func (e *Engine) UpdatePortSrcMAC(ctx context.Context, portID string, mac []byte) error {
	e.idToNIDMu.RLock()
	defer e.idToNIDMu.RUnlock()
	idBytes := make([]byte, binary.Size(e.idToNID[portID]))
	binary.BigEndian.PutUint64(idBytes, e.idToNID[portID])

	entry := &fwdpb.TableEntryAddRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: srcMACTable}},
		ContextId: &fwdpb.ContextId{Id: e.id},
		EntryDesc: &fwdpb.EntryDesc{
			Entry: &fwdpb.EntryDesc_Exact{
				Exact: &fwdpb.ExactEntryDesc{
					Fields: []*fwdpb.PacketFieldBytes{{
						FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT}},
						Bytes:   idBytes,
					}},
				},
			},
		},
		Actions: []*fwdpb.ActionDesc{{
			ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
			Action: &fwdpb.ActionDesc_Update{
				Update: &fwdpb.UpdateActionDesc{
					FieldId: &fwdpb.PacketFieldId{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC,
						},
					},
					Type:  fwdpb.UpdateType_UPDATE_TYPE_SET,
					Value: mac,
				},
			},
		}},
	}
	if _, err := e.Server.TableEntryAdd(ctx, entry); err != nil {
		return err
	}

	return nil
}

// CreateExternalPort creates an external port (connected to other devices).
func (e *Engine) CreateExternalPort(ctx context.Context, id, devName string) error {
	nid, err := createKernelPort(ctx, e.id, e.Server, id, devName)
	if err != nil {
		return err
	}
	e.idToNIDMu.Lock()
	e.idToNID[devName] = nid
	e.idToNIDMu.Unlock()

	update := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: id}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_Kernel{
				Kernel: &fwdpb.KernelPortUpdateDesc{
					Inputs: []*fwdpb.ActionDesc{{ // Lookup in layer 2 table.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: layer2PuntTable}},
							},
						},
					}, { // Lookup in layer 3 table.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: layer3PuntTable}},
							},
						},
					}, { // Decap L2 header.
						ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
						Action: &fwdpb.ActionDesc_Decap{
							Decap: &fwdpb.DecapActionDesc{
								HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
							},
						},
					}, { // Lookup in FIB.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibSelectorTable}},
							},
						},
					}, { // Encap a L2 header.
						ActionType: fwdpb.ActionType_ACTION_TYPE_ENCAP,
						Action: &fwdpb.ActionDesc_Encap{
							Encap: &fwdpb.EncapActionDesc{
								HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
							},
						},
					}, { // Lookup in the neighbor table.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: neighborTable}},
							},
						},
					}, {
						ActionType: fwdpb.ActionType_ACTION_TYPE_OUTPUT,
					}},
					Outputs: []*fwdpb.ActionDesc{{ // update the src mac address with the configured port's mac address.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: srcMACTable}},
							},
						},
					}},
				},
			},
		},
	}
	if _, err := e.Server.PortUpdate(ctx, update); err != nil {
		return err
	}
	return nil
}

// CreateLocalPort creates an local (ie TAP) port for the given linux device name.
func (e *Engine) CreateLocalPort(ctx context.Context, id, devName, externalID string) error {
	nid, err := createTapPort(ctx, e.id, e.Server, id, devName)
	if err != nil {
		return err
	}

	e.idToNIDMu.Lock()
	e.idToNID[id] = nid
	e.idToNIDMu.Unlock()

	update := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: id}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_Kernel{
				Kernel: &fwdpb.KernelPortUpdateDesc{
					Inputs: []*fwdpb.ActionDesc{{ // Lookup in layer 2 table.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: layer2PuntTable}},
							},
						},
					}, { // Lookup in FIB.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibSelectorTable}},
							},
						},
					}, { // Lookup in the neighbor table.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: neighborTable}},
							},
						},
					}, {
						ActionType: fwdpb.ActionType_ACTION_TYPE_OUTPUT,
					}},
				},
			},
		},
	}
	if _, err := e.Server.PortUpdate(ctx, update); err != nil {
		return err
	}
	if externalID == "" {
		return nil
	}
	_, err = e.Server.AttributeUpdate(ctx, &fwdpb.AttributeUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		ObjectId:  &fwdpb.ObjectId{Id: id},
		AttrId:    attributes.SwapActionRelatedPort,
		AttrValue: externalID,
	})
	if err != nil {
		return err
	}
	_, err = e.Server.AttributeUpdate(ctx, &fwdpb.AttributeUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		ObjectId:  &fwdpb.ObjectId{Id: externalID},
		AttrId:    attributes.SwapActionRelatedPort,
		AttrValue: id,
	})
	if err != nil {
		return err
	}

	return nil
}

// CreateLocalPort returns the counters for the object by name.
func (e *Engine) GetCounters(ctx context.Context, name string) (*fwdpb.ObjectCountersReply, error) {
	return e.Server.ObjectCounters(ctx, &fwdpb.ObjectCountersRequest{
		ObjectId:  &fwdpb.ObjectId{Id: name},
		ContextId: &fwdpb.ContextId{Id: e.id},
	})
}
