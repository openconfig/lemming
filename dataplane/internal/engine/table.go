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
	"encoding/hex"
	"net"

	log "github.com/golang/glog"

	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	DefaultContextID = "default"
	fibV4Table       = "fib-v4"
	fibV6Table       = "fib-v6"
	srcMACTable      = "port-mac"
	fibSelectorTable = "fib-selector"
	neighborTable    = "neighbor"
	layer2PuntTable  = "layer2-punt"
	layer3PuntTable  = "layer3-punt"
	arpPuntTable     = "arp-punt"
)

func mustParseHex(hexStr string) []byte {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return b
}

// SetupForwardingTables creates the forwarding tables.
func SetupForwardingTables(ctx context.Context, c fwdpb.ServiceClient) error {
	_, err := c.ContextCreate(context.Background(), &fwdpb.ContextCreateRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
	})
	if err != nil {
		return err
	}

	v4FIB := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_PREFIX,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibV4Table}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Prefix{
				Prefix: &fwdpb.PrefixTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
						},
					}},
				},
			},
		},
	}
	if _, err := c.TableCreate(ctx, v4FIB); err != nil {
		return err
	}
	v6FIB := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_PREFIX,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibV6Table}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_DROP}},
			Table: &fwdpb.TableDesc_Prefix{
				Prefix: &fwdpb.PrefixTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
						},
					}},
				},
			},
		},
	}
	if _, err := c.TableCreate(ctx, v6FIB); err != nil {
		return err
	}
	portMAC := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
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
		return err
	}
	neighbor := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
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
		return err
	}
	if err := createFIBSelector(ctx, c); err != nil {
		return err
	}
	if err := createLayer2PuntTable(ctx, c); err != nil {
		return err
	}
	if err := createLayer3PuntTable(ctx, c); err != nil {
		return err
	}
	return nil
}

// createFIBSelector creates a table that controls which forwarding table is used.
func createFIBSelector(ctx context.Context, c fwdpb.ServiceClient) error {
	fieldID := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION,
		},
	}

	ipVersion := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibSelectorTable}},
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
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: fibSelectorTable,
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
						TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibV4Table}},
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
						TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibV6Table}},
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

// createLayer2PuntTable creates a table to packets to punt at layer 2 (input port and mac dst).
func createLayer2PuntTable(ctx context.Context, c fwdpb.ServiceClient) error {
	arp := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: arpPuntTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE,
						},
					}},
				},
			},
		},
	}
	if _, err := c.TableCreate(ctx, arp); err != nil {
		return err
	}
	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: arpPuntTable,
			},
		},
		Entries: []*fwdpb.TableEntryAddRequest_Entry{{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes: mustParseHex("0806"),
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE,
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
	if _, err := c.TableEntryAdd(ctx, entries); err != nil {
		return err
	}
	layer2 := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_PREFIX,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: layer2PuntTable}},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
				Action: &fwdpb.ActionDesc_Lookup{
					Lookup: &fwdpb.LookupActionDesc{
						TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: arpPuntTable}},
					},
				},
			}},
			Table: &fwdpb.TableDesc_Prefix{
				Prefix: &fwdpb.PrefixTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT,
						},
					}, {
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST,
						},
					}},
				},
			},
		},
	}
	if _, err := c.TableCreate(ctx, layer2); err != nil {
		return err
	}
	return nil
}

// addLayer2PuntRule adds rule to output packets to a corresponding port based on the destination MAC and input port.
func addLayer2PuntRule(ctx context.Context, c fwdpb.ServiceClient, portID uint64, mac, macMask []byte) error {
	nidBytes := make([]byte, binary.Size(portID))
	binary.BigEndian.PutUint64(nidBytes, portID)

	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: layer2PuntTable,
			},
		},
		Entries: []*fwdpb.TableEntryAddRequest_Entry{{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Prefix{
					Prefix: &fwdpb.PrefixEntryDesc{
						Fields: []*fwdpb.PacketFieldMaskedBytes{{
							Bytes: nidBytes,
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT,
								},
							},
						}, {
							Bytes: mac,
							Masks: macMask,
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST,
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
	if _, err := c.TableEntryAdd(ctx, entries); err != nil {
		return err
	}
	return nil
}

// createLayer3PuntTable creates a table controlling whether packets to punt at layer 3 (input port and IP dst).
func createLayer3PuntTable(ctx context.Context, c fwdpb.ServiceClient) error {
	multicast := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: layer3PuntTable}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
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
		},
	}
	if _, err := c.TableCreate(ctx, multicast); err != nil {
		return err
	}
	return nil
}

// AddLayer3PuntRule adds rule to output packets to a corresponding port based on the destination IP and input port.
func AddLayer3PuntRule(ctx context.Context, c fwdpb.ServiceClient, portName string, ip []byte) error {
	nameToIDMu.Lock()
	defer nameToIDMu.Unlock()
	portID := nameToID[portName]

	nidBytes := make([]byte, binary.Size(portID))
	binary.BigEndian.PutUint64(nidBytes, portID)

	log.Infof("adding layer3 punt rule: portName %s, id %d, ip %v", portName, portID, ip)

	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
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
	if _, err := c.TableEntryAdd(ctx, entries); err != nil {
		return err
	}
	return nil
}

// nextHop returns the forwarding actions for a nexthop.
func nextHopToActions(nh *dpb.NextHop) []*fwdpb.ActionDesc {
	var nextHopIP []byte
	if nhIPStr := nh.GetIp(); nhIPStr != "" {
		nextHop := net.ParseIP(nhIPStr)
		nextHopIP = nextHop.To4()
		if nextHopIP == nil {
			nextHopIP = nextHop.To16()
		}
	}
	nextHopAct := &fwdpb.ActionDesc{ // Set the next hop IP in the packet's metadata.
		ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
		Action: &fwdpb.ActionDesc_Update{
			Update: &fwdpb.UpdateActionDesc{
				FieldId: &fwdpb.PacketFieldId{
					Field: &fwdpb.PacketField{
						FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP,
					},
				},
				Type:  fwdpb.UpdateType_UPDATE_TYPE_SET,
				Value: nextHopIP,
			},
		},
	}
	if nextHopIP == nil {
		nextHopAct = &fwdpb.ActionDesc{ // Set the next hop IP in the packet's metadata.
			ActionType: fwdpb.ActionType_ACTION_TYPE_UPDATE,
			Action: &fwdpb.ActionDesc_Update{
				Update: &fwdpb.UpdateActionDesc{
					FieldId: &fwdpb.PacketFieldId{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP,
						},
					},
					Type: fwdpb.UpdateType_UPDATE_TYPE_COPY,
					Field: &fwdpb.PacketFieldId{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
						},
					},
				},
			},
		}
	}
	return append([]*fwdpb.ActionDesc{{ // Set the output port.
		ActionType: fwdpb.ActionType_ACTION_TYPE_TRANSMIT,
		Action: &fwdpb.ActionDesc_Transmit{
			Transmit: &fwdpb.TransmitActionDesc{
				PortId: &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: nh.GetPort()}},
			},
		},
	}, nextHopAct}, nh.GetPreTransmitActions()...)
}

// AddIPRoute adds a route to the FIB with the input next hops.
func AddIPRoute(ctx context.Context, c fwdpb.ServiceClient, v4 bool, ip, mask []byte, nextHops []*dpb.NextHop) error {
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

		actions = []*fwdpb.ActionDesc{{
			ActionType: fwdpb.ActionType_ACTION_TYPE_SELECT_ACTION_LIST,
			Action: &fwdpb.ActionDesc_Select{
				Select: &fwdpb.SelectActionListActionDesc{
					SelectAlgorithm: fwdpb.SelectActionListActionDesc_SELECT_ALGORITHM_CRC32, // TODO: should algo + hash be configurable?
					FieldIds: []*fwdpb.PacketFieldId{{Field: &fwdpb.PacketField{
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
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		EntryDesc: &fwdpb.EntryDesc{
			Entry: &fwdpb.EntryDesc_Prefix{
				Prefix: &fwdpb.PrefixEntryDesc{
					Fields: []*fwdpb.PacketFieldMaskedBytes{{
						FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
						Bytes:   ip,
						Masks:   mask,
					}},
				},
			},
		},
		Actions: actions,
	}
	if _, err := c.TableEntryAdd(ctx, entry); err != nil {
		return err
	}

	return nil
}

// DeleteIPRoute deletes a route from the FIB.
func DeleteIPRoute(ctx context.Context, c fwdpb.ServiceClient, v4 bool, ip, mask []byte) error {
	fib := fibV6Table
	if v4 {
		fib = fibV4Table
	}
	entry := &fwdpb.TableEntryRemoveRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fib}},
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
		EntryDesc: &fwdpb.EntryDesc{
			Entry: &fwdpb.EntryDesc_Prefix{
				Prefix: &fwdpb.PrefixEntryDesc{
					Fields: []*fwdpb.PacketFieldMaskedBytes{{
						FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST}},
						Bytes:   ip,
						Masks:   mask,
					}},
				},
			},
		},
	}
	if _, err := c.TableEntryRemove(ctx, entry); err != nil {
		return err
	}

	return nil
}

// AddNeighbor adds a neighbor to the neighbor table.
func AddNeighbor(ctx context.Context, c fwdpb.ServiceClient, ip, mac []byte) error {
	entry := &fwdpb.TableEntryAddRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: neighborTable}},
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
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
	if _, err := c.TableEntryAdd(ctx, entry); err != nil {
		return err
	}

	return nil
}

// RemoveNeighbor removes a neighbor from the neighbor table.
func RemoveNeighbor(ctx context.Context, c fwdpb.ServiceClient, ip []byte) error {
	entry := &fwdpb.TableEntryRemoveRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: neighborTable}},
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
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
	if _, err := c.TableEntryRemove(ctx, entry); err != nil {
		return err
	}

	return nil
}

// UpdatePortSrcMAC updates a port's source mac address.
func UpdatePortSrcMAC(ctx context.Context, c fwdpb.ServiceClient, portName string, mac []byte) error {
	nameToIDMu.RLock()
	defer nameToIDMu.RUnlock()
	idBytes := make([]byte, binary.Size(nameToID[portName]))
	binary.BigEndian.PutUint64(idBytes, nameToID[portName])

	entry := &fwdpb.TableEntryAddRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: srcMACTable}},
		ContextId: &fwdpb.ContextId{Id: DefaultContextID},
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
	if _, err := c.TableEntryAdd(ctx, entry); err != nil {
		return err
	}

	return nil
}
