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
	"encoding/hex"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	contextID        = "default"
	fibV4Table       = "fib-v4"
	fibV6Table       = "fib-v6"
	srcMACTable      = "port-mac"
	fibSelectorTable = "fib-selector"
	neighborTable    = "neighbor"
)

var (
	etherTypeIPV4 = mustParseHex("0x0800")
	etherTypeIPV6 = mustParseHex("0x86DD")
	etherTypeARP  = mustParseHex("0x0806")
	ipProtoICMP   = mustParseHex("0x01")
	ipProtoICMPV6 = mustParseHex("0x3A")
)

func mustParseHex(hexStr string) []byte {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return b
}

// SetupTables creates the forwarding tables.
func SetupTables(ctx context.Context, c fwdpb.ServiceClient) error {
	_, err := c.ContextCreate(context.Background(), &fwdpb.ContextCreateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
	})
	if err != nil {
		return err
	}

	v4FIB := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
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
		ContextId: &fwdpb.ContextId{Id: contextID},
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
		ContextId: &fwdpb.ContextId{Id: contextID},
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
		ContextId: &fwdpb.ContextId{Id: contextID},
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
	return nil
}

// createFIBSelector creates a table that controls which forwarding table is used.
func createFIBSelector(ctx context.Context, c fwdpb.ServiceClient) error {
	fieldID := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE,
		},
	}

	etherType := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
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
	if _, err := c.TableCreate(ctx, etherType); err != nil {
		return err
	}
	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
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
							Bytes:   etherTypeIPV4,
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
							Bytes:   etherTypeIPV6,
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

// AddIPRoute adds a route to the FIB.
func AddIPRoute(ctx context.Context, c fwdpb.ServiceClient, v4 bool, ip, mask, nextHopIP []byte, port string) error {
	fib := fibV6Table
	if v4 {
		fib = fibV4Table
	}

	entry := &fwdpb.TableEntryAddRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fib}},
		ContextId: &fwdpb.ContextId{Id: contextID},
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
		Actions: []*fwdpb.ActionDesc{{ // Set the output port.
			ActionType: fwdpb.ActionType_ACTION_TYPE_TRANSMIT,
			Action: &fwdpb.ActionDesc_Transmit{
				Transmit: &fwdpb.TransmitActionDesc{
					PortId: &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: port}},
				},
			},
		}, { // Set the next hop IP in the packet's metadata.
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
		}, { // Lookup in the neighbor table.
			ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
			Action: &fwdpb.ActionDesc_Lookup{
				Lookup: &fwdpb.LookupActionDesc{
					TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: neighborTable}},
				},
			},
		}, { // Finally output the packet.
			ActionType: fwdpb.ActionType_ACTION_TYPE_OUTPUT,
		}},
	}
	if _, err := c.TableEntryAdd(ctx, entry); err != nil {
		return err
	}

	return nil
}

// AddNeighbor adds a neighbor to the neighbor table.
func AddNeighbor(ctx context.Context, c fwdpb.ServiceClient, ip, mac []byte) error {
	entry := &fwdpb.TableEntryAddRequest{
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: neighborTable}},
		ContextId: &fwdpb.ContextId{Id: contextID},
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
