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

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

var (
	etherBroadcast = mustParseHex("FFFFFFFFFFFF")
	etherMulticast = mustParseHex("010000000000")
	etherIPV6Multi = mustParseHex("333300000000")
)

// CreateExternalPort creates an external port (connected to other devices).
func CreateExternalPort(ctx context.Context, c fwdpb.ServiceClient, name string) error {
	return createKernelPort(ctx, c, name)
}

// CreateLocalPort creates an local (ie TAP) port for the given linux device name.
func CreateLocalPort(ctx context.Context, c fwdpb.ServiceClient, name string) error {
	return createKernelPort(ctx, c, name)
}

func createKernelPort(ctx context.Context, c fwdpb.ServiceClient, name string) error {
	port := &fwdpb.PortCreateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
		Port: &fwdpb.PortDesc{
			PortType: fwdpb.PortType_PORT_TYPE_KERNEL,
			PortId: &fwdpb.PortId{
				ObjectId: &fwdpb.ObjectId{Id: name},
			},
			Port: &fwdpb.PortDesc_Kernel{
				Kernel: &fwdpb.KernelPortDesc{
					DeviceName: name,
				},
			},
		},
	}
	portID, err := c.PortCreate(ctx, port)
	if err != nil {
		return err
	}
	if err := AddLayer2PuntRule(ctx, c, portID.GetObjectIndex().GetIndex(), etherBroadcast, mustParseHex("FFFFFFFFFFFF")); err != nil {
		return err
	}
	if err := AddLayer2PuntRule(ctx, c, portID.GetObjectIndex().GetIndex(), etherMulticast, mustParseHex("010000000000")); err != nil {
		return err
	}
	if err := AddLayer2PuntRule(ctx, c, portID.GetObjectIndex().GetIndex(), etherIPV6Multi, mustParseHex("FFFF00000000")); err != nil {
		return err
	}

	update := &fwdpb.PortUpdateRequest{
		ContextId: &fwdpb.ContextId{Id: contextID},
		PortId:    &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: name}},
		Update: &fwdpb.PortUpdateDesc{
			Port: &fwdpb.PortUpdateDesc_Kernel{
				Kernel: &fwdpb.KernelPortUpdateDesc{
					Inputs: []*fwdpb.ActionDesc{{ // Turn on packet tracing. TODO: put this behind a flag.
						ActionType: fwdpb.ActionType_ACTION_TYPE_DEBUG,
					}, { // Lookup in layer 2 table.
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
