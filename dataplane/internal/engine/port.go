package engine

import (
	"context"
	"fmt"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// CreateExternalPort creates an external (connected to other devices).
func CreateExternalPort(ctx context.Context, c fwdpb.ServiceClient, name string) error {
	port := &fwdpb.PortCreateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
		Port: &fwdpb.PortDesc{
			PortType: fwdpb.PortType_PORT_TYPE_KERNEL,
			Port: &fwdpb.PortDesc_Kernel{
				Kernel: &fwdpb.KernelPortDesc{
					DeviceName: name,
				},
			},
		},
	}
	if _, err := c.PortCreate(ctx, port); err != nil {
		return err
	}

	if err := setupPuntRules(ctx, c, name); err != nil {
		return err
	}

	update := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: name}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_Kernel{
				Kernel: &fwdpb.KernelPortUpdateDesc{
					Inputs: []*fwdpb.ActionDesc{{ // Check EtherType punt rules.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: puntEtherTypeTable(name)}},
							},
						},
					}, { // Check IP Proto punt rules.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: puntIPProtoType(name)}},
							},
						},
					}, { // Lookup in FIB.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibSelectorTable}},
							},
						},
					}},
					Outputs: []*fwdpb.ActionDesc{{ // Lookup in the to port to src MAC table.
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
	if _, err := c.PortUpdate(ctx, update); err != nil {
		return err
	}
	return nil
}

// CreateExternalPort creates an external (connected to other devices).
func CreateInternalPort(ctx context.Context, c fwdpb.ServiceClient, name string) error {
	port := &fwdpb.PortCreateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
		Port: &fwdpb.PortDesc{
			PortType: fwdpb.PortType_PORT_TYPE_KERNEL,
			Port: &fwdpb.PortDesc_Kernel{
				Kernel: &fwdpb.KernelPortDesc{
					DeviceName: name,
				},
			},
		},
	}
	if _, err := c.PortCreate(ctx, port); err != nil {
		return err
	}

	update := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: name}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_Kernel{
				Kernel: &fwdpb.KernelPortUpdateDesc{
					Inputs: []*fwdpb.ActionDesc{{ // Lookup in FIB.
						ActionType: fwdpb.ActionType_ACTION_TYPE_LOOKUP,
						Action: &fwdpb.ActionDesc_Lookup{
							Lookup: &fwdpb.LookupActionDesc{
								TableId: &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: fibSelectorTable}},
							},
						},
					}},
				},
			},
		},
	}
	if _, err := c.PortUpdate(ctx, update); err != nil {
		return err
	}
	return nil
}

func puntEtherTypeTable(port string) string {
	return fmt.Sprintf("%s-punt-etherType", port)
}

func puntIPProtoType(port string) string {
	return fmt.Sprintf("%s-punt-ipproto", port)
}

func setupPuntRules(ctx context.Context, c fwdpb.ServiceClient, name string) error {
	// Add rule to write ARP packets to tap interface.
	etherTypePunt := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: puntEtherTypeTable(name)}},
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
	if _, err := c.TableCreate(ctx, etherTypePunt); err != nil {
		return err
	}
	arp := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: puntEtherTypeTable(name),
			},
		},
		Entries: []*fwdpb.TableEntryAddRequest_Entry{{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes: etherTypeARP,
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
				ActionType: fwdpb.ActionType_ACTION_TYPE_TRANSMIT,
				Action: &fwdpb.ActionDesc_Transmit{
					Transmit: &fwdpb.TransmitActionDesc{
						PortId: &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: NameToTap(name)}},
					},
				},
			}},
		}},
	}
	if _, err := c.TableEntryAdd(ctx, arp); err != nil {
		return err
	}
	// Add rule to write ICMP and ICMPv6 packets to tap interface.
	ipProtoPunt := &fwdpb.TableCreateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
		Desc: &fwdpb.TableDesc{
			TableType: fwdpb.TableType_TABLE_TYPE_EXACT,
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: puntIPProtoType(name)}},
			Actions:   []*fwdpb.ActionDesc{{ActionType: fwdpb.ActionType_ACTION_TYPE_CONTINUE}},
			Table: &fwdpb.TableDesc_Exact{
				Exact: &fwdpb.ExactTableDesc{
					FieldIds: []*fwdpb.PacketFieldId{{
						Field: &fwdpb.PacketField{
							FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO,
						},
					}},
				},
			},
		},
	}
	if _, err := c.TableCreate(ctx, ipProtoPunt); err != nil {
		return err
	}

	icmp := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
		TableId: &fwdpb.TableId{
			ObjectId: &fwdpb.ObjectId{
				Id: puntIPProtoType(name),
			},
		},
		Entries: []*fwdpb.TableEntryAddRequest_Entry{{
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes: ipProtoICMP,
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO,
								},
							},
						}},
					},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_TRANSMIT,
				Action: &fwdpb.ActionDesc_Transmit{
					Transmit: &fwdpb.TransmitActionDesc{
						PortId: &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: NameToTap(name)}},
					},
				},
			}},
		}, {
			EntryDesc: &fwdpb.EntryDesc{
				Entry: &fwdpb.EntryDesc_Exact{
					Exact: &fwdpb.ExactEntryDesc{
						Fields: []*fwdpb.PacketFieldBytes{{
							Bytes: ipProtoICMPV6,
							FieldId: &fwdpb.PacketFieldId{
								Field: &fwdpb.PacketField{
									FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO,
								},
							},
						}},
					},
				},
			},
			Actions: []*fwdpb.ActionDesc{{
				ActionType: fwdpb.ActionType_ACTION_TYPE_TRANSMIT,
				Action: &fwdpb.ActionDesc_Transmit{
					Transmit: &fwdpb.TransmitActionDesc{
						PortId: &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: NameToTap(name)}},
					},
				},
			}},
		}},
	}
	if _, err := c.TableEntryAdd(ctx, icmp); err != nil {
		return err
	}

	return nil
}
