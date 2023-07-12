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

package prefix

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/mock_fwdpacket"
	tables "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/tabletestutil"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdaction/actions"
)

// ErrTest returns a non-empty string if err and str do not match.
func ErrTest(err error, str string) string {
	if err == nil {
		if str == "" {
			return ""
		}
		return `did not get expected error "` + str + `"`
	}
	if str == "" {
		return fmt.Sprintf("unexpected error %q", err.Error())
	}
	if !strings.Contains(err.Error(), str) {
		return fmt.Sprintf("got error %q, want %q", err.Error(), str)
	}
	return ""
}

// TestPacketField verifies that newPrefixKey accepts and rejects packet
// fields appropriately.
func TestPacketFields(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Register a mock parser.
	const validSize = 10
	const invalidSize = 11
	parser := mock_fwdpacket.NewMockParser(ctrl)
	parser.EXPECT().MaxSize(gomock.Any()).Return(validSize).AnyTimes()
	parser.EXPECT().Validate(gomock.Any(), invalidSize).Return(errors.New("fwdpacket: validation error")).AnyTimes()
	parser.EXPECT().Validate(gomock.Any(), validSize).Return(nil).AnyTimes()
	fwdpacket.Register(parser)

	tests := []struct {
		fields     []fwdpb.PacketFieldNum
		size       int
		countFudge int
		err        string
	}{
		{
			fields: []fwdpb.PacketFieldNum{
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC,
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE,
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION,
			},
			countFudge: -1,
			err:        "missing field",
			size:       validSize,
		},
		{
			fields: []fwdpb.PacketFieldNum{
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC,
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE,
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION,
			},
			countFudge: +1,
			err:        "duplicate field",
			size:       validSize,
		},
		{
			fields: []fwdpb.PacketFieldNum{
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC,
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE,
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION,
			},
			err:  "validation error",
			size: invalidSize,
		},
		{
			fields: []fwdpb.PacketFieldNum{
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_SRC,
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE,
				fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION,
			},
			size: validSize,
		},
	}

	for pos, test := range tests {
		kd := tables.KeyDescFields(test.fields)
		ed := EntryDescFields(test.fields, len(test.fields)+test.countFudge, test.size)
		_, err := newPrefixKey(kd, ed)
		if str := ErrTest(err, test.err); str != "" {
			t.Errorf("%d: newPrefixKey failed, err %v.", pos, str)
		}
	}
}

// EntryDescBytes builds a test entry descriptor from a set of packet bytes.
// The argument count can inject duplicate fields within the desc. The argument
// size determines the size of the field. The values of count and size
// determine if the entry descriptor can successfully generate a key.
func EntryDescBytes(fields []fwdpb.PacketBytes, count, size int) EntryDesc {
	var ed EntryDesc
	for i := 0; i < len(fields); i++ {
		if count == 0 {
			break
		}
		copied := proto.Clone(&fields[i]).(*fwdpb.PacketBytes)
		ed = append(ed, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{
				Bytes: copied,
			},
			Bytes: make([]byte, size),
		})
		count--
	}

	// Add additional fields as duplicates.
	for ; count > 0; count-- {
		field := proto.Clone(&fields[0]).(*fwdpb.PacketBytes)
		ed = append(ed, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{
				Bytes: field,
			},
			Bytes: make([]byte, size),
		})
	}
	return ed
}

// EntryDescFields builds a test entry descriptor from a set of fields.
// The argument count can inject duplicate fields within the desc. The argument
// size determines the size of the field. The values of count and size
// determine if the entry descriptor can successfully generate a key.
func EntryDescFields(fields []fwdpb.PacketFieldNum, count int, size int) EntryDesc {
	var ed EntryDesc
	for _, id := range fields {
		if count == 0 {
			break
		}
		id := id
		ed = append(ed, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{
				Field: &fwdpb.PacketField{
					FieldNum: id,
				},
			},
			Bytes: make([]byte, size),
		})
		count--
	}

	// Add additional fields as duplicates.
	for ; count > 0; count-- {
		id := fields[0]
		ed = append(ed, &fwdpb.PacketFieldMaskedBytes{
			FieldId: &fwdpb.PacketFieldId{
				Field: &fwdpb.PacketField{
					FieldNum: id,
				},
			},
			Bytes: make([]byte, size),
		})
	}
	return ed
}

// TestPacketBytes verifies that newPrefixKey accepts and rejects packet
// bytes appropriately.
func TestPacketBytes(t *testing.T) {
	// Create a controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Register a mock parser.
	const validSize = 10
	const invalidSize = 11
	parser := mock_fwdpacket.NewMockParser(ctrl)
	parser.EXPECT().MaxSize(gomock.Any()).Return(validSize).AnyTimes()
	parser.EXPECT().Validate(gomock.Any(), invalidSize).Return(errors.New("fwdpacket: validation error")).AnyTimes()
	parser.EXPECT().Validate(gomock.Any(), validSize).Return(nil).AnyTimes()
	fwdpacket.Register(parser)

	tests := []struct {
		fields     []fwdpb.PacketBytes
		size       int
		countFudge int
		err        string
	}{
		{
			fields: []fwdpb.PacketBytes{
				{
					HeaderGroup: fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L2,
					Offset:      20,
					Size:        4,
				},
				{
					HeaderGroup: fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L3,
					Offset:      20,
					Size:        4,
				},
			},
			countFudge: -1,
			err:        "missing field",
			size:       validSize,
		},
		{
			fields: []fwdpb.PacketBytes{
				{
					HeaderGroup: fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L2,
					Offset:      20,
					Size:        4,
				},
				{
					HeaderGroup: fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L3,
					Offset:      20,
					Size:        4,
				},
			},
			countFudge: 1,
			err:        "duplicate field",
			size:       validSize,
		},
		{
			fields: []fwdpb.PacketBytes{
				{
					HeaderGroup: fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L2,
					Offset:      20,
					Size:        4,
				},
				{
					HeaderGroup: fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L3,
					Offset:      20,
					Size:        4,
				},
			},
			err:  "validation error",
			size: invalidSize,
		},
		{
			fields: []fwdpb.PacketBytes{
				{
					HeaderGroup: fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L2,
					Offset:      20,
					Size:        4,
				},
				{
					HeaderGroup: fwdpb.PacketHeaderGroup_PACKET_HEADER_GROUP_L3,
					Offset:      20,
					Size:        4,
				},
			},
			size: validSize,
		},
	}

	for pos, test := range tests {
		kd := tables.KeyDescBytes(test.fields)
		ed := EntryDescBytes(test.fields, len(test.fields)+test.countFudge, test.size)
		_, err := newPrefixKey(kd, ed)
		if str := ErrTest(err, test.err); str != "" {
			t.Errorf("%d: newPrefixKey failed, err %v.", pos, str)
		}
	}
}
