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

package engine

import (
	"context"
	"encoding/binary"
	"encoding/hex"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

var (
	etherBroadcast     = mustParseHex("FFFFFFFFFFFF")
	etherBroadcastMask = mustParseHex("FFFFFFFFFFFF")
	etherMulticast     = mustParseHex("0180C2000000")
	etherMulticastMask = mustParseHex("FFFFFF000000")
	etherIPV6Multi     = mustParseHex("333300000000")
	etherIPV6MultiMask = mustParseHex("FFFF00000000")
)

func mustParseHex(hexStr string) []byte {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return b
}

// createFIBSelector creates a table that controls which forwarding table is used.
func createFIBSelector(ctx context.Context, id string, c fwdpb.ForwardingServer) error {
	fieldID := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION,
		},
	}

	ipVersion := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: id},
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
		ContextId: &fwdpb.ContextId{Id: id},
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

// createLayer2PuntTable creates a table to packets to punt at layer 2 (input port and mac dst).
func createLayer2PuntTable(ctx context.Context, id string, c fwdpb.ForwardingServer) error {
	arp := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: id},
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
		ContextId: &fwdpb.ContextId{Id: id},
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
				ActionType: fwdpb.ActionType_ACTION_TYPE_SWAP_OUTPUT_INTERNAL_EXTERNAL,
			}, {
				ActionType: fwdpb.ActionType_ACTION_TYPE_OUTPUT,
			}},
		}},
	}
	if _, err := c.TableEntryAdd(ctx, entries); err != nil {
		return err
	}
	layer2 := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: id},
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
func addLayer2PuntRule(ctx context.Context, ctxID string, c fwdpb.ForwardingServer, portNID uint64, mac, macMask []byte) error {
	nidBytes := make([]byte, binary.Size(portNID))
	binary.BigEndian.PutUint64(nidBytes, portNID)

	entries := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: ctxID},
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
				ActionType: fwdpb.ActionType_ACTION_TYPE_SWAP_OUTPUT_INTERNAL_EXTERNAL,
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
func createLayer3PuntTable(ctx context.Context, ctxID string, c fwdpb.ForwardingServer) error {
	multicast := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: ctxID},
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

// createPortAndEntries creates a port and sets up the punt rules
func createPortAndEntries(ctx context.Context, ctxID string, c fwdpb.ForwardingServer, t fwdpb.PortType, id, devName string) (uint64, error) {
	port := &fwdpb.PortCreateRequest{
		ContextId: &fwdpb.ContextId{Id: ctxID},
		Port: &fwdpb.PortDesc{
			PortType: t,
			PortId: &fwdpb.PortId{
				ObjectId: &fwdpb.ObjectId{Id: id},
			},
			Port: &fwdpb.PortDesc_Kernel{
				Kernel: &fwdpb.KernelPortDesc{
					DeviceName: devName,
				},
			},
		},
	}
	if t == fwdpb.PortType_PORT_TYPE_TAP {
		port.Port.Port = &fwdpb.PortDesc_Tap{
			Tap: &fwdpb.TAPPortDesc{
				DeviceName: devName,
			},
		}
	}
	portID, err := c.PortCreate(ctx, port)
	if err != nil {
		return 0, err
	}
	if err := addLayer2PuntRule(ctx, ctxID, c, portID.GetObjectIndex().GetIndex(), etherBroadcast, etherBroadcastMask); err != nil {
		return 0, err
	}
	if err := addLayer2PuntRule(ctx, ctxID, c, portID.GetObjectIndex().GetIndex(), etherMulticast, etherMulticastMask); err != nil {
		return 0, err
	}
	if err := addLayer2PuntRule(ctx, ctxID, c, portID.GetObjectIndex().GetIndex(), etherIPV6Multi, etherIPV6MultiMask); err != nil {
		return 0, err
	}
	return portID.GetObjectIndex().GetIndex(), nil
}
