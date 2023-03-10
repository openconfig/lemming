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

package packet_test

import (
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/packettestutil"
)

var data = []byte{0x01, 0x02, 0x03, 0x04}

func TestMetadata(t *testing.T) {
	tests := []packettestutil.PacketFieldTest{{
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE,
		Orig: [][]byte{
			data,
		},
		Queries: []packettestutil.FieldQuery{{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION, 0),
			Err: "failed",
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF, 0),
			Result: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_LENGTH, 0),
			Result: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_32, 3),
			Result: []byte{0x00, 0x00, 0x00, 0x00},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_32, 2),
			Result: []byte{0x00, 0x00, 0x00, 0x00},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_32, 0),
			Result: []byte{0x00, 0x00, 0x00, 0x00},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_24, 0),
			Result: []byte{0x00, 0x00, 0x00},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_16, 0),
			Result: []byte{0x00, 0x00},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_8, 0),
			Result: []byte{0x00},
		},
		},
		Updates: []packettestutil.FieldUpdate{{
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF, 0),
			Arg: []byte{0xab, 0xcd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			Op:  fwdpacket.OpSet,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_32, 3),
			Arg: []byte{0xab, 0xcd, 0x00, 0x00},
			Op:  fwdpacket.OpSet,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_32, 0),
			Arg: []byte{0xab, 0xcd, 0x12, 0x34},
			Op:  fwdpacket.OpSet,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_32, 0),
			Arg: []byte{0x01},
			Op:  fwdpacket.OpDec,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_32, 0),
			Arg: []byte{0x01},
			Op:  fwdpacket.OpInc,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_24, 0),
			Arg: []byte{0xab, 0xcd, 0x12},
			Op:  fwdpacket.OpSet,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_16, 0),
			Arg: []byte{0xab, 0xcd},
			Op:  fwdpacket.OpSet,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_16, 0),
			Arg: []byte{0x01},
			Op:  fwdpacket.OpDec,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_16, 0),
			Arg: []byte{0x01},
			Op:  fwdpacket.OpInc,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_8, 0),
			Arg: []byte{0xab},
			Op:  fwdpacket.OpSet,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_8, 0),
			Arg: []byte{0x01},
			Op:  fwdpacket.OpDec,
		}, {
			ID:  fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_ATTRIBUTE_8, 0),
			Arg: []byte{0x01},
			Op:  fwdpacket.OpInc,
		}},
		Final: [][]byte{
			data,
		},
	}}

	packettestutil.TestPacketFields("metadata", t, tests)
}
