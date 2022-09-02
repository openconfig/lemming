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

package actions

import (
	"testing"

	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// TestReparse tests the reparse action and builder.
func TestReparse(t *testing.T) {
	// Create a controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a context for the test.
	ctx := fwdcontext.New("test", "fwd")

	// List of packet field ids (protobuf and Lucius type) used to save
	// across the reparse and the header for the reparse.
	header := fwdpb.PacketHeaderId_OPAQUE
	field1 := []*fwdpb.PacketFieldId{
		{
			Field: &fwdpb.PacketField{
				FieldNum: fwdpb.PacketFieldNum_IP_ADDR_DST.Enum(),
			},
		},
		{
			Field: &fwdpb.PacketField{
				FieldNum: fwdpb.PacketFieldNum_IP_ADDR_SRC.Enum(),
			},
		},
	}
	var field2 []fwdpacket.FieldID
	for _, f := range field1 {
		field2 = append(field2, fwdpacket.NewFieldID(f))
	}

	// Create a reparse action using its builder. Prepend an ethernet header
	// as an opaque set of bytes before the reparse.
	desc := fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_REPARSE_ACTION.Enum(),
	}
	prepend := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x08, 0x00}
	reparseDesc := fwdpb.ReparseActionDesc{
		HeaderId: &header,
		FieldIds: field1,
		Prepend:  prepend,
	}
	proto.SetExtension(&desc, fwdpb.E_ReparseActionDesc_Extension, &reparseDesc)
	action, err := fwdaction.New(&desc, ctx)
	if err != nil {
		t.Errorf("NewAction failed for desc %v, err %v.", desc, err)
	}

	packet := mock_fwdpacket.NewMockPacket(ctrl)
	packet.EXPECT().Reparse(header, field2, prepend).Return(nil)
	action.Process(packet, nil)
}
