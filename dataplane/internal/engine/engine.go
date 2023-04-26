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
	"sync"

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
	layer2PuntTable  = "layer2-punt"
	layer3PuntTable  = "layer3-punt"
	arpPuntTable     = "arp-punt"
)

// Engine contains a routing context and methods to manage it.
type Engine struct {
	client     fwdpb.ForwardingClient
	id         string
	nameToIDMu sync.RWMutex
	nameToID   map[string]uint64
}

// New creates a new engine and sets up the forwarding tables.
func New(ctx context.Context, id string, c fwdpb.ForwardingClient) (*Engine, error) {
	e := &Engine{
		id:       id,
		client:   c,
		nameToID: map[string]uint64{},
	}

	_, err := c.ContextCreate(context.Background(), &fwdpb.ContextCreateRequest{
		ContextId: &fwdpb.ContextId{Id: id},
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
	if _, err := c.TableCreate(ctx, v4FIB); err != nil {
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
	if _, err := c.TableCreate(ctx, v6FIB); err != nil {
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
	if _, err := c.TableCreate(ctx, portMAC); err != nil {
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
	if _, err := c.TableCreate(ctx, neighbor); err != nil {
		return nil, err
	}
	if err := createFIBSelector(ctx, e.id, c); err != nil {
		return nil, err
	}
	if err := createLayer2PuntTable(ctx, e.id, c); err != nil {
		return nil, err
	}
	if err := createLayer3PuntTable(ctx, e.id, c); err != nil {
		return nil, err
	}
	return e, nil
}

// AddLayer3PuntRule adds rule to output packets to a corresponding port based on the destination IP and input port.
func (e *Engine) AddLayer3PuntRule(ctx context.Context, portName string, ip []byte) error {
	e.nameToIDMu.Lock()
	defer e.nameToIDMu.Unlock()
	portID := e.nameToID[portName]

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
				ActionType: fwdpb.ActionType_ACTION_TYPE_SWAP_OUTPUT_TAP_EXTERNAL,
			}, {
				ActionType: fwdpb.ActionType_ACTION_TYPE_OUTPUT,
			}},
		}},
	}
	if _, err := e.client.TableEntryAdd(ctx, entries); err != nil {
		return err
	}
	return nil
}

// AddIPRoute adds a route to the FIB with the input next hops.
func (e *Engine) AddIPRoute(ctx context.Context, v4 bool, ip, mask []byte, vrf uint64, nextHops []*dpb.NextHop) error {
	fib := fibV6Table
	if v4 {
		fib = fibV4Table
	}

	actions := nextHopToActions(nextHops[0])

	if len(nextHops) > 1 {
		var actLists []*fwdpb.ActionList
		for _, nh := range nextHops {
			actLists = append(actLists, &fwdpb.ActionList{
				Weight:  nh.GetWeight(),
				Actions: nextHopToActions(nh),
			})
		}

		// If there are multiple next-hops, configure the route to use ECMP or WCMP.
		actions = []*fwdpb.ActionDesc{{
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
		}}
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
		Actions: actions,
	}
	if _, err := e.client.TableEntryAdd(ctx, entry); err != nil {
		return err
	}

	return nil
}

// DeleteIPRoute deletes a route from the FIB.
func (e *Engine) DeleteIPRoute(ctx context.Context, v4 bool, ip, mask []byte, vrf uint64) error {
	fib := fibV6Table
	if v4 {
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
	if _, err := e.client.TableEntryRemove(ctx, entry); err != nil {
		return err
	}

	return nil
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
	if _, err := e.client.TableEntryAdd(ctx, entry); err != nil {
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
	if _, err := e.client.TableEntryRemove(ctx, entry); err != nil {
		return err
	}

	return nil
}

// UpdatePortSrcMAC updates a port's source mac address.
func (e *Engine) UpdatePortSrcMAC(ctx context.Context, portName string, mac []byte) error {
	e.nameToIDMu.RLock()
	defer e.nameToIDMu.RUnlock()
	idBytes := make([]byte, binary.Size(e.nameToID[portName]))
	binary.BigEndian.PutUint64(idBytes, e.nameToID[portName])

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
	if _, err := e.client.TableEntryAdd(ctx, entry); err != nil {
		return err
	}

	return nil
}

// CreateExternalPort creates an external port (connected to other devices).
func (e *Engine) CreateExternalPort(ctx context.Context, name string) error {
	id, err := createKernelPort(ctx, e.id, e.client, name)
	if err != nil {
		return err
	}
	e.nameToIDMu.Lock()
	e.nameToID[name] = id
	e.nameToIDMu.Unlock()

	update := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: name}},
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
	if _, err := e.client.PortUpdate(ctx, update); err != nil {
		return err
	}
	return nil
}

// CreateLocalPort creates an local (ie TAP) port for the given linux device name.
func (e *Engine) CreateLocalPort(ctx context.Context, name string, fd int) error {
	id, err := createTapPort(ctx, e.id, e.client, name, fd)
	if err != nil {
		return err
	}
	e.nameToIDMu.Lock()
	e.nameToID[name] = id
	e.nameToIDMu.Unlock()

	update := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: name}},
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
	if _, err := e.client.PortUpdate(ctx, update); err != nil {
		return err
	}
	return nil
}

// CreateLocalPort returns the counters for the object by name.
func (e *Engine) GetCounters(ctx context.Context, name string) (*fwdpb.ObjectCountersReply, error) {
	return e.client.ObjectCounters(ctx, &fwdpb.ObjectCountersRequest{
		ObjectId:  &fwdpb.ObjectId{Id: name},
		ContextId: &fwdpb.ContextId{Id: e.id},
	})
}
