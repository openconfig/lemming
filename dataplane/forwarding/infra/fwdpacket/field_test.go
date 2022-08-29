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

package fwdpacket

import (
	"testing"

	"google.golang.org/protobuf/proto"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// TestFieldID tests the packet FieldID
func TestFieldID(t *testing.T) {
	tests := []struct {
		fieldID *fwdpb.PacketFieldId
		want    FieldID
	}{
		// FieldID using a packet field number.
		{
			fieldID: &fwdpb.PacketFieldId{
				Field: &fwdpb.PacketField{
					FieldNum: fwdpb.PacketFieldNum_IP_PROTO.Enum(),
					Instance: proto.Uint32(10),
				}},
			want: FieldID{
				Instance: 10,
				Num:      fwdpb.PacketFieldNum_IP_PROTO,
			},
		},
		// FieldID using a UDF with a non-zero offset.
		{
			fieldID: &fwdpb.PacketFieldId{
				Bytes: &fwdpb.PacketBytes{
					HeaderGroup: fwdpb.PacketHeaderGroup_L2.Enum(),
					Instance:    proto.Uint32(10),
					Offset:      proto.Uint32(5),
					Size:        proto.Uint32(1),
				}},
			want: FieldID{
				IsUDF:    true,
				Instance: 10,
				Header:   fwdpb.PacketHeaderGroup_L2,
				Size:     1,
				Offset:   5,
			},
		},
		// FieldID using a UDF with a zero offset.
		{
			fieldID: &fwdpb.PacketFieldId{
				Bytes: &fwdpb.PacketBytes{
					HeaderGroup: fwdpb.PacketHeaderGroup_L2.Enum(),
					Instance:    proto.Uint32(10),
					Offset:      proto.Uint32(0),
					Size:        proto.Uint32(1),
				}},
			want: FieldID{
				IsUDF:    true,
				Instance: 10,
				Header:   fwdpb.PacketHeaderGroup_L2,
				Size:     1,
			},
		},
	}

	for _, test := range tests {
		if got := NewFieldID(test.fieldID); got != test.want {
			t.Errorf("NewFieldID(%v) = %+v; want %+v.", test.fieldID, got, test.want)
		}
	}
}
