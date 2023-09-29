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

package protogen

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"

	"github.com/openconfig/lemming/dataplane/standalone/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/standalone/apigen/saiast"
)

const (
	commonType = `
syntax = "proto3";

package lemming.dataplane.sai;
	
import "google/protobuf/timestamp.proto";
import "google/protobuf/descriptor.proto";

option go_package = "github.com/openconfig/lemming/dataplane/standalone/proto;sai";

extend google.protobuf.FieldOptions {
	optional int32 attr_enum_value = 50000;
}

extend google.protobuf.MessageOptions {
	optional ObjectType sai_type = 50001;
}

message AclActionData {
	bool enable = 1;
	oneof parameter {
		uint64 uint = 2;
		uint64 int = 3;
		bytes mac = 4;
		bytes ip = 5;
		uint64 oid = 6;
		Uint64List objlist = 7;
		bytes ipaddr = 8;
	};
}

message ACLCapability {
	bool is_action_list_mandatory = 1;
	repeated AclActionType action_list = 2;
}

message AclFieldData {
	bool enable = 1;
	oneof mask {
		uint64 mask_uint = 2;
		uint64 mask_int = 3;
		bytes mask_mac = 4;
		bytes mask_ip = 5;
		Uint64List mask_list = 6;
	};
	oneof data {
		bool data_bool = 7;
		uint64 data_uint = 8;
		int64 data_int = 9;
		bytes data_mac = 10;
		bytes data_ip = 11;
		Uint64List data_list = 12;
	};
}

message Uint64List {
	repeated uint64 list = 1;
}

message ACLResource {
	AclStage stage = 1;
	AclBindPointType bind_point = 2;
	uint32 avail_num = 3;
}

message BfdSessionStateChangeNotificationData {
	uint64 bfd_session_id = 1;
	BfdSessionState session_state = 2;
}

message FabricPortReachability {
	uint32 switch_id = 1;
	bool reachable = 2;
}

message FdbEntry {
	uint64 switch_id = 1;
	bytes mac_address = 2;
	uint64 bv_id = 3;
}


message FdbEventNotificationData {
    FdbEvent event_type = 1;
	FdbEntry fdb_entry = 2;
	repeated FdbEntryAttribute attrs = 3;
}

message InsegEntry {
	uint64 switch_id = 1;
	uint32 label = 2;
}

message IpPrefix {
	bytes addr = 1;
	bytes mask = 2;
}

message IpmcEntry {
	uint64 switch_id = 1;
	uint64 vr_id = 2;
	IpmcEntryType type = 3;
	bytes destination = 4;
	bytes source = 5;
}

message IpsecSaStatusNotificationData {
    uint64 ipsec_sa_id = 1;
	IpsecSaOctetCountStatus ipsec_sa_octet_count_status = 2;
	bool ipsec_egress_sn_at_max_limit = 3;
}

message L2mcEntry {
	uint64 switch_id = 1;
	uint64 bv_id = 2;
	L2mcEntryType type = 3;
	bytes destination = 4;
	bytes source = 5;
}

message LatchStatus {
	bool current_status = 1;
	bool changed = 2;
}

message UintMap {
	map<uint32, uint32> uintmap = 1;
}

message McastFdbEntry {
	uint64 switch_id = 1;
	bytes mac_address = 2;
	uint64 bv_id = 3;
}

message MySidEntry {
	uint64 switch_id = 1;
	uint64 vr_id = 2;
	uint32 locator_block_len = 3;
	uint32 locator_node_len = 4;
	uint32 function_len = 5;
	uint32 args_len = 6;
	bytes sid = 7;
}

message NatEntryData{
	oneof key {
		bytes key_src_ip = 2;
		bytes key_dst_ip = 3;
		uint32 key_proto = 4;
		uint32 key_l4_src_port = 5;
		uint32 key_l4_dst_port = 6;
	};
	oneof mask {
		bytes mask_src_ip = 7;
		bytes mask_dst_ip = 8;
		uint32 mask_proto = 9;
		uint32 mask_l4_src_port = 10;
		uint32 mask_l4_dst_port = 11;
	};
}

message NatEntry {
	uint64 switch_id = 1;
	uint64 vr_id = 2;
	NatType nat_type = 3;
	NatEntryData data = 4;
}

message NeighborEntry {
	uint64 switch_id = 1;
	uint64 rif_id = 2;
	bytes ip_address = 3;
}

message PortEyeValues {
	uint32 lane = 1;
	int32 left = 2;
	int32 right = 3;
	int32 up = 4;
	int32 down = 5;
}

message PortLaneLatchStatus {
	uint32 lane = 1;
	LatchStatus value = 2;
}

message PortOperStatusNotification {
	uint64 port_id = 1;
	PortOperStatus port_state = 2;
}

message PRBS_RXState {
	PortPrbsRxStatus rx_status = 1;
	uint32 error_count = 2;
}


message	QOSMapParams {
	uint32 tc = 1;
	uint32 dscp = 2;
	uint32 dot1p = 3;
	uint32 prio = 4;
	uint32 pg = 5;
	uint32 queue_index = 6;
	PacketColor color = 7;
	uint32 mpls_exp = 8;
	uint32 fc = 9;
}

message QOSMap {
	QOSMapParams key = 1;
	QOSMapParams value = 2;
}

message QueueDeadlockNotificationData {
	uint64 queue_id = 1;
	QueuePfcDeadlockEventType event= 2;
	bool app_managed_recovery = 3;
}

message RouteEntry {
	uint64 switch_id = 1;
	uint64 vr_id = 2;
	IpPrefix destination = 3;
}

message SystemPortConfig {
	uint32 port_id = 1;
	uint32 attached_switch_id = 2;
	uint32 attached_core_index = 3;
	uint32 attached_core_port_index = 4;
	uint32 speed = 5;
	uint32 num_voq = 6;
}

message HMAC {
	uint32 key_id = 1;
	repeated uint32 hmac = 2;
}

message TLVEntry {
	oneof entry {
		bytes ingress_node = 1; 
		bytes egress_node = 2;
		bytes opaque_container = 3;
		HMAC hmac = 4;
	}
}

message Uint32Range {
	uint64 min = 1;
	uint64 max = 2;
}

message ObjectTypeQueryRequest {
	uint64 object = 1;
}

message ObjectTypeQueryResponse {
	ObjectType type = 1;
}

service Entrypoint {
  rpc ObjectTypeQuery(ObjectTypeQueryRequest) returns (ObjectTypeQueryResponse) {}
}
`
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		desc    string
		inAst   *saiast.SAIAPI
		inInfo  *docparser.SAIInfo
		want    map[string]string
		wantErr string
	}{{
		desc:   "empty",
		inAst:  &saiast.SAIAPI{},
		inInfo: &docparser.SAIInfo{},
		want: map[string]string{
			"common.proto": commonType,
		},
	}, {
		desc:  "common enum",
		inAst: &saiast.SAIAPI{},
		inInfo: &docparser.SAIInfo{
			Enums: map[string][]*docparser.Enum{
				"sai_foo_t": {{Name: "SAI_FOO_ONE", Value: 0}, {Name: "SAI_FOO_TWO", Value: 1}},
			},
		},
		want: map[string]string{
			"common.proto": commonType + `
enum Foo {
	FOO_UNSPECIFIED = 0;
	FOO_ONE = 1;
	FOO_TWO = 2;
}
`,
		},
	}, {
		desc:  "common enum with unspecified value",
		inAst: &saiast.SAIAPI{},
		inInfo: &docparser.SAIInfo{
			Enums: map[string][]*docparser.Enum{
				"sai_foo_t": {{Name: "SAI_FOO_UNSPECIFIED", Value: 0}, {Name: "SAI_FOO_TWO", Value: 1}},
			},
		},
		want: map[string]string{
			"common.proto": commonType + `
enum Foo {
	FOO_UNSPECIFIED = 0;
	FOO_SAI_UNSPECIFIED = 1;
	FOO_TWO = 2;
}
`,
		},
	}, {
		desc:  "common list and attribute",
		inAst: &saiast.SAIAPI{},
		inInfo: &docparser.SAIInfo{
			Attrs: map[string]*docparser.Attr{
				"FOO": {
					ReadFields: []*docparser.AttrTypeName{{
						SaiType:    "sai_u32_list_t",
						MemberName: "sample_list",
					}, {
						SaiType:    "sai_uint8_t",
						MemberName: "sample_int",
					}},
				},
			},
		},
		want: map[string]string{
			"common.proto": commonType + `
message FooAttribute {
	repeated uint32 sample_list = 1 [(attr_enum_value) = 1];
	optional uint32 sample_int = 2 [(attr_enum_value) = 2];
}
`,
		},
	}, {
		desc: "simple API",
		inAst: &saiast.SAIAPI{
			Ifaces: []*saiast.SAIInterface{{
				Name: "sai_sample_api_t",
				Funcs: []*saiast.TypeDecl{{
					Name: "create_foo",
					Typ:  "sai_create_foo_t",
				}, {
					Name: "remove_foo",
					Typ:  "sai_remove_foo_t",
				}, {
					Name: "set_foo_attribute",
					Typ:  "sai_set_foo_attribute_t",
				}, {
					Name: "get_foo_attribute",
					Typ:  "sai_get_foo_attribute_t",
				}},
			}},
			Funcs: map[string]*saiast.SAIFunc{
				"sai_create_foo_t":        {},
				"sai_remove_foo_t":        {},
				"sai_set_foo_attribute_t": {},
				"sai_get_foo_attribute_t": {},
			},
		},
		inInfo: &docparser.SAIInfo{
			Attrs: map[string]*docparser.Attr{
				"FOO": {
					CreateFields: []*docparser.AttrTypeName{{
						SaiType:    "sai_uint8_t",
						MemberName: "sample_uint",
						EnumName:   "SAI_FOO_ATTR_UINT",
					}},
					SetFields: []*docparser.AttrTypeName{{
						SaiType:    "sai_int8_t",
						MemberName: "sample_int",
						EnumName:   "SAI_FOO_ATTR_INT",
					}},
					ReadFields: []*docparser.AttrTypeName{{
						SaiType:    "sai_int8_t",
						MemberName: "sample_int",
						EnumName:   "SAI_FOO_ATTR_UINT",
					}, {
						SaiType:    "sai_uint8_t",
						MemberName: "sample_uint",
						EnumName:   "SAI_FOO_ATTR_INT",
					}},
				},
			},
		},
		want: map[string]string{
			"common.proto": commonType + `
message FooAttribute {
	optional int32 sample_int = 1 [(attr_enum_value) = 1];
	optional uint32 sample_uint = 2 [(attr_enum_value) = 2];
}
`,
			"sample.proto": `
syntax = "proto3";

package lemming.dataplane.sai;

import "dataplane/standalone/proto/common.proto";

option go_package = "github.com/openconfig/lemming/dataplane/standalone/proto;sai";


enum FooAttr {
	FOO_ATTR_UNSPECIFIED = 0;
	FOO_ATTR_UINT = 1;
	FOO_ATTR_INT = 2;
}

message CreateFooRequest {
	option (sai_type) = OBJECT_TYPE_UNSPECIFIED;
	optional uint32 sample_uint = 1;
}

message CreateFooResponse {
	uint64 oid = 1;
}

message RemoveFooRequest {
	uint64 oid = 1;
}

message RemoveFooResponse {
}

message SetFooAttributeRequest {
	uint64 oid = 1;
	optional int32 sample_int = 2;
}

message SetFooAttributeResponse {
}

message GetFooAttributeRequest {
	uint64 oid = 1;
	repeated FooAttr attr_type = 2;
}

message GetFooAttributeResponse {
	FooAttribute attr = 1;
}


service Sample {
	rpc CreateFoo (CreateFooRequest) returns (CreateFooResponse) {}
	rpc RemoveFoo (RemoveFooRequest) returns (RemoveFooResponse) {}
	rpc SetFooAttribute (SetFooAttributeRequest) returns (SetFooAttributeResponse) {}
	rpc GetFooAttribute (GetFooAttributeRequest) returns (GetFooAttributeResponse) {}
}
`,
		},
	}, {
		desc: "simple switch scoped API API",
		inAst: &saiast.SAIAPI{
			Ifaces: []*saiast.SAIInterface{{
				Name: "sai_sample_api_t",
				Funcs: []*saiast.TypeDecl{{
					Name: "create_foo",
					Typ:  "sai_create_foo_t",
				}, {
					Name: "remove_foo",
					Typ:  "sai_remove_foo_t",
				}, {
					Name: "set_foo_attribute",
					Typ:  "sai_set_foo_attribute_t",
				}, {
					Name: "get_foo_attribute",
					Typ:  "sai_get_foo_attribute_t",
				}},
			}},
			Funcs: map[string]*saiast.SAIFunc{
				"sai_create_foo_t":        {Params: []saiast.TypeDecl{{Name: "object_id", Typ: "*sai_object_id_t"}, {Name: "switch_id", Typ: "sai_object_id_t"}}},
				"sai_remove_foo_t":        {},
				"sai_set_foo_attribute_t": {},
				"sai_get_foo_attribute_t": {},
			},
		},
		inInfo: &docparser.SAIInfo{
			Attrs: map[string]*docparser.Attr{
				"FOO": {
					CreateFields: []*docparser.AttrTypeName{{
						SaiType:    "sai_uint8_t",
						MemberName: "sample_uint",
						EnumName:   "SAI_FOO_ATTR_UINT",
					}},
					ReadFields: []*docparser.AttrTypeName{{
						SaiType:    "sai_int8_t",
						MemberName: "sample_int",
						EnumName:   "SAI_FOO_ATTR_UINT",
					}, {
						SaiType:    "sai_uint8_t",
						MemberName: "sample_uint",
						EnumName:   "SAI_FOO_ATTR_INT",
					}},
				},
			},
		},
		want: map[string]string{
			"common.proto": commonType + `
message FooAttribute {
	optional int32 sample_int = 1 [(attr_enum_value) = 1];
	optional uint32 sample_uint = 2 [(attr_enum_value) = 2];
}
`,
			"sample.proto": `
syntax = "proto3";

package lemming.dataplane.sai;

import "dataplane/standalone/proto/common.proto";

option go_package = "github.com/openconfig/lemming/dataplane/standalone/proto;sai";


enum FooAttr {
	FOO_ATTR_UNSPECIFIED = 0;
	FOO_ATTR_UINT = 1;
	FOO_ATTR_INT = 2;
}

message CreateFooRequest {
	option (sai_type) = OBJECT_TYPE_UNSPECIFIED;
	uint64 switch = 1;
	optional uint32 sample_uint = 2;
}

message CreateFooResponse {
	uint64 oid = 1;
}

message RemoveFooRequest {
	uint64 oid = 1;
}

message RemoveFooResponse {
}

message GetFooAttributeRequest {
	uint64 oid = 1;
	repeated FooAttr attr_type = 2;
}

message GetFooAttributeResponse {
	FooAttribute attr = 1;
}


service Sample {
	rpc CreateFoo (CreateFooRequest) returns (CreateFooResponse) {}
	rpc RemoveFoo (RemoveFooRequest) returns (RemoveFooResponse) {}
	rpc GetFooAttribute (GetFooAttributeRequest) returns (GetFooAttributeResponse) {}
}
`,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, gotErr := Generate(tt.inInfo, tt.inAst)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("Generate() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want); d != "" {
				t.Fatalf("Generate() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}
