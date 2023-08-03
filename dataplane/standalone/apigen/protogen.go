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

package main

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"

	strcase "github.com/stoewer/go-strcase"
)

// populateTmplDataFromFunc populatsd the protobuf template struct from a SAI function call.
func populateTmplDataFromFunc(protoData *protoApiData, funcName, entryType, operation, typeName, apiName string, isSwitchScoped bool) error {
	if _, ok := protoData.apis[apiName]; !ok {
		protoData.apis[apiName] = &protoAPITmplData{
			Enums:       make(map[string]protoEnum),
			ServiceName: trimSAIName(apiName, true, false),
		}
	}

	req := &protoTmplMessage{
		Name: strcase.UpperCamelCase(funcName + "_request"),
	}
	resp := &protoTmplMessage{
		Name: strcase.UpperCamelCase(funcName + "_response"),
	}

	idField := protoTmplField{
		Index:     1,
		ProtoType: "uint64",
		Name:      "oid",
	}
	if entryType != "" {
		idField = protoTmplField{
			Index:     1,
			ProtoType: entryType,
			Name:      "entry",
		}
	}

	// Handle proto generation
	switch operation {
	case "create":
		requestIdx := 1
		if isSwitchScoped {
			req.Fields = append(req.Fields, protoTmplField{
				Index:     requestIdx,
				ProtoType: "uint64",
				Name:      "switch",
			})
		} else if entryType != "" {
			req.Fields = append(req.Fields, idField)
		}
		requestIdx++
		attrs, err := createAttrs(requestIdx, protoData.docInfo, protoData.docInfo.attrs[typeName].createFields, false)
		if err != nil {
			return err
		}
		req.Attrs = attrs

		resp.Fields = append(resp.Fields, protoTmplField{
			Index:     1,
			ProtoType: "uint64",
			Name:      "oid",
		})
	case "set_attribute":
		// If there are no settable attributes, do nothing.
		if len(protoData.docInfo.attrs[typeName].setFields) == 0 {
			return nil
		}
		req.Fields = append(req.Fields, idField)
		req.AttrsWrapperStart = "oneof attr {"
		req.AttrsWrapperEnd = "}"
		attrs, err := createAttrs(2, protoData.docInfo, protoData.docInfo.attrs[typeName].setFields, true)
		if err != nil {
			return err
		}
		req.Attrs = attrs
	case "get_attribute":
		req.Fields = append(req.Fields, idField, protoTmplField{
			ProtoType: strcase.UpperCamelCase(typeName + " attr"),
			Index:     2,
			Name:      "attr_type",
		})

		// attrEnum is the special emun that describes the possible values can be set/get for the API.
		attrEnum := protoEnum{
			Name:   strcase.UpperCamelCase(typeName + "_ATTR"),
			Values: []protoEnumValues{{Index: 0, Name: typeName + "_ATTR_UNSPECIFIED"}},
		}

		// For the attributes, generate code for the type if needed.
		for i, attr := range protoData.docInfo.attrs[typeName].readFields {
			attrEnum.Values = append(attrEnum.Values, protoEnumValues{
				Index: i + 1,
				Name:  strings.TrimPrefix(attr.EnumName, "SAI_"),
			})
		}
		protoData.apis[apiName].Enums[attrEnum.Name] = attrEnum

		attrs, err := createAttrs(1, protoData.docInfo, protoData.docInfo.attrs[typeName].readFields, false)
		if err != nil {
			return err
		}
		resp.Attrs = attrs
	default:
		return nil
	}
	protoData.apis[apiName].Messages = append(protoData.apis[apiName].Messages, *req, *resp)
	protoData.apis[apiName].RPCs = append(protoData.apis[apiName].RPCs, protoRPC{
		RequestName:  req.Name,
		ResponseName: resp.Name,
		Name:         strcase.UpperCamelCase(funcName),
	})
	return nil
}

func createAttrs(startIdx int, xmlInfo *protoGenInfo, attrs []attrTypeName, inOneof bool) ([]protoTmplField, error) {
	fields := []protoTmplField{}
	for _, attr := range attrs {
		// Function pointers are attempted as streaming RPCs.
		if strings.Contains(attr.SaiType, "sai_pointer_t") {
			continue
		}
		// Proto field names can't beging with numbers, prepend _.
		name := attr.MemberName
		if unicode.IsDigit(rune(attr.MemberName[0])) {
			name = fmt.Sprintf("_%s", name)
		}
		field := protoTmplField{
			Index: startIdx,
			Name:  name,
		}
		typ, _, err := saiTypeToProtoType(attr.SaiType, xmlInfo, inOneof)
		if err != nil {
			return nil, err
		}
		field.ProtoType = typ
		fields = append(fields, field)
		startIdx++
	}
	return fields, nil
}

var (
	protoTmpl = template.Must(template.New("cc").Parse(`
syntax = "proto3";

package lemming.dataplane.sai;

import "dataplane/standalone/proto/common.proto";

option go_package = "github.com/openconfig/lemming/proto/dataplane/sai";

{{ range .Enums }}
enum {{ .Name }} {
	{{- range .Values }}
	{{ .Name }} = {{ .Index }};
	{{- end}}
}
{{ end }}

{{ range .Messages }}
message {{ .Name }} {
	{{- range .Fields }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }};
	{{- end }}
	{{ .AttrsWrapperStart -}}
	{{- range .Attrs }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }};
	{{- end }}
	{{ .AttrsWrapperEnd }}
}
{{ end }}


service {{ .ServiceName }} {
	{{- range .RPCs }}
	rpc {{ .Name }} ({{ .RequestName }}) returns ({{ .ResponseName }}) {}
	{{- end }}
}
`))
	protoCommonTmpl = template.Must(template.New("common").Parse(`
syntax = "proto3";

package lemming.dataplane.sai;
	
option go_package = "github.com/openconfig/lemming/proto/dataplane/sai";

{{ range .Enums }}
enum {{ .Name }} {
	{{- range .Values }}
	{{ .Name }} = {{ .Index }};
	{{- end}}
}
{{ end }}

{{ range .Messages }}
{{ . }}
{{ end }}

{{ range .Lists }}
message {{ .Name }} {
	{{- range .Fields }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }};
	{{- end }}
}
{{ end }}
`))
)

// protoApiData contains the input and output for protobuf generation.
type protoApiData struct {
	docInfo *protoGenInfo
	apis    map[string]*protoAPITmplData
	common  *protoCommonTmplData
}

// protoAPITmplData contains the formated information needed to render the protobuf template.
type protoAPITmplData struct {
	Messages    []protoTmplMessage
	RPCs        []protoRPC
	Enums       map[string]protoEnum
	ServiceName string
}

// protoCommonTmplData contains the formated information needed to render the protobuf template.
type protoCommonTmplData struct {
	Messages []string
	Enums    map[string]*protoEnum
	Lists    map[string]*protoTmplMessage
}

type protoEnum struct {
	Name   string
	Values []protoEnumValues
}

type protoEnumValues struct {
	Index int
	Name  string
}

type protoTmplMessage struct {
	Name              string
	AttrsWrapperStart string
	AttrsWrapperEnd   string
	Fields            []protoTmplField
	Attrs             []protoTmplField
}

type protoTmplField struct {
	ProtoType string
	Name      string
	Index     int
}

type protoRPC struct {
	RequestName  string
	ResponseName string
	Name         string
}

type saiTypeInfo struct {
	Repeated   bool
	ProtoType  string
	MessageDef string
}

var saiTypeToProto = map[string]saiTypeInfo{
	"bool": {
		ProtoType: "bool",
	},
	"char": {
		ProtoType: "bytes",
	},
	"sai_uint8_t": {
		ProtoType: "uint32",
	},
	"sai_int8_t": {
		ProtoType: "int32",
	},
	"sai_uint16_t": {
		ProtoType: "uint32",
	},
	"sai_int16_t": {
		ProtoType: "int32",
	},
	"sai_uint32_t": {
		ProtoType: "uint32",
	},
	"sai_int32_t": {
		ProtoType: "uint32",
	},
	"sai_uint64_t": {
		ProtoType: "uint64",
	},
	"sai_int64_t": {
		ProtoType: "int64",
	},
	"sai_mac_t": {
		ProtoType: "bytes",
	},
	"sai_ip4_t": {
		ProtoType: "bytes",
	},
	"sai_ip6_t": {
		ProtoType: "bytes",
	},
	"sai_s32_list_t": {
		Repeated:  true,
		ProtoType: "int32",
	},
	"sai_object_id_t": {
		ProtoType: "uint64",
	},
	"sai_object_list_t": {
		Repeated:  true,
		ProtoType: "uint64",
	},
	"sai_encrypt_key_t": {
		ProtoType: "bytes",
	},
	"sai_auth_key_t": {
		ProtoType: "bytes",
	},
	"sai_macsec_sak_t": {
		ProtoType: "bytes",
	},
	"sai_macsec_auth_key_t": {
		ProtoType: "bytes",
	},
	"sai_macsec_salt_t": {
		ProtoType: "bytes",
	},
	"sai_u32_list_t": {
		Repeated:  true,
		ProtoType: "uint32",
	},
	"sai_segment_list_t": {
		Repeated:  true,
		ProtoType: "bytes",
	},
	"sai_s8_list_t": {
		Repeated:  true,
		ProtoType: "int32",
	},
	"sai_u8_list_t": {
		Repeated:  true,
		ProtoType: "uint32",
	},
	"sai_port_err_status_list_t": {
		Repeated:  true,
		ProtoType: "PortErrStatus",
	},
	"sai_vlan_list_t": {
		Repeated:  true,
		ProtoType: "uint32",
	}, // The non-scalar types could be autogenerated, but that aren't that many so create messages by hand.
	"sai_u32_range_t": {
		ProtoType: "Uint32Range",
		MessageDef: `message Uint32Range {
uint64 min = 1;
uint64 max = 2;
}`,
	},
	"sai_ip_address_t": {
		ProtoType: "bytes",
	},
	"sai_map_list_t": { // Wrap the map in a message because maps can't be repeated.
		Repeated:  true,
		ProtoType: "UintMap",
		MessageDef: `message UintMap {
	map<uint32, uint32> uintmap = 1;
}`,
	},
	"sai_tlv_list_t": {
		Repeated:  true,
		ProtoType: "TLVEntry",
		MessageDef: `message HMAC {
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
}`,
	},
	"sai_qos_map_list_t": {
		Repeated:  true,
		ProtoType: "QOSMap",
		MessageDef: `
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
}`,
	},
	"sai_system_port_config_t": {
		ProtoType: "SystemPortConfig",
		MessageDef: `message SystemPortConfig {
	uint32 port_id = 1;
	uint32 attached_switch_id = 2;
	uint32 attached_core_index = 3;
	uint32 attached_core_port_index = 4;
	uint32 speed = 5;
	uint32 num_voq = 6;
}`,
	},
	"sai_system_port_config_list_t": {
		Repeated:  true,
		ProtoType: "SystemPortConfig",
	},
	"sai_ip_address_list_t": {
		Repeated:  true,
		ProtoType: "bytes",
	},
	"sai_port_eye_values_list_t": {
		Repeated:  true,
		ProtoType: "PortEyeValues",
		MessageDef: `message PortEyeValues {
	uint32 lane = 1;
	int32 left = 2;
	int32 right = 3;
	int32 up = 4;
	int32 down = 5;
}`,
	},
	"sai_prbs_rx_state_t": {
		ProtoType: "PRBS_RXState",
		MessageDef: `message PRBS_RXState {
	PortPrbsRxStatus rx_status = 1;
	uint32 error_count = 2;
}`,
	},
	"sai_fabric_port_reachability_t": {
		ProtoType: "FabricPortReachability",
		MessageDef: `message FabricPortReachability {
	uint32 switch_id = 1;
	bool reachable = 2;
}`,
	},
	"sai_acl_resource_list_t": {
		Repeated:  true,
		ProtoType: "ACLResource",
		MessageDef: `message ACLResource {
	AclStage stage = 1;
	AclBindPointType bind_point = 2;
	uint32 avail_num = 3;
}`,
	},
	"sai_acl_capability_t": {
		ProtoType: "ACLCapability",
		MessageDef: `message ACLCapability {
	bool is_action_list_mandatory = 1;
	repeated int32 action_list = 2;
}`,
	},
	"sai_acl_field_data_t": {
		ProtoType: "AclFieldData",
		MessageDef: `message AclFieldData {
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
}`,
	},
	"sai_acl_action_data_t": {
		ProtoType: "AclActionData",
		MessageDef: `message AclActionData {
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
}`,
	},
	"sai_fdb_entry_t": {
		ProtoType: "FdbEntry",
		MessageDef: `message FdbEntry {
	uint64 switch_id = 1;
	bytes mac_address = 2;
	uint64 bv_id = 3;
}`,
	},
	"sai_ipmc_entry_t": {
		ProtoType: "IpmcEntry",
		MessageDef: `message IpmcEntry {
	uint64 switch_id = 1;
	uint64 vr_id = 2;
	IpmcEntryType type = 3;
	bytes destination = 4;
	bytes source = 5;
}`,
	},
	"sai_l2mc_entry_t": {
		ProtoType: "L2mcEntry",
		MessageDef: `message L2mcEntry {
	uint64 switch_id = 1;
	uint64 bv_id = 2;
	L2mcEntryType type = 3;
	bytes destination = 4;
	bytes source = 5;
}`,
	},
	"sai_mcast_fdb_entry_t": {
		ProtoType: "McastFdbEntry",
		MessageDef: `message McastFdbEntry {
	uint64 switch_id = 1;
	bytes mac_address = 2;
	uint64 bv_id = 3;
}`,
	},
	"sai_inseg_entry_t": {
		ProtoType: "InsegEntry",
		MessageDef: `message InsegEntry {
	uint64 switch_id = 1;
	uint32 label = 2;
}`,
	},
	"sai_nat_entry_data_t": {
		ProtoType: "NatEntryData",
		MessageDef: `message NatEntryData{
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
}`,
	},
	"sai_nat_entry_t": {
		ProtoType: "NatEntry",
		MessageDef: `message NatEntry {
	uint64 switch_id = 1;
	uint64 vr_id = 2;
	NatType nat_type = 3;
	NatEntryData data = 4;
}`,
	},
	"sai_neighbor_entry_t": {
		ProtoType: "NeighborEntry",
		MessageDef: `message NeighborEntry {
	uint64 switch_id = 1;
	uint64 rif_id = 2;
	bytes ip_address = 3;
}`,
	},
	"sai_ip_prefix_t": {
		ProtoType: "IpPrefix",
		MessageDef: `message IpPrefix {
	bytes addr = 1;
	bytes mask = 2;
}`,
	},
	"sai_route_entry_t": {
		ProtoType: "RouteEntry",
		MessageDef: `message RouteEntry {
	uint64 switch_id = 1;
	uint64 vr_id = 2;
	IpPrefix destination = 3;
}`,
	},
	"sai_my_sid_entry_t": {
		ProtoType: "MySidEntry",
		MessageDef: `message MySidEntry {
	uint64 switch_id = 1;
	uint64 vr_id = 2;
	uint32 locator_block_len = 3;
	uint32 locator_node_len = 4;
	uint32 function_len = 5;
	uint32 args_len = 6;
	bytes sid = 7;
}`,
	},
}

// saiTypeToProtoTypeCompound handles compound sai types (eg list of enums).
// The map key contains the base type (eg list) and func accepts the subtype (eg an enum type)
// and returns the full type string (eg repeated sample_enum).
var saiTypeToProtoTypeCompound = map[string]func(subType string, xmlInfo *protoGenInfo, inOneof bool) (string, bool){
	"sai_s32_list_t": func(subType string, xmlInfo *protoGenInfo, inOneof bool) (string, bool) {
		if _, ok := xmlInfo.enums[subType]; !ok {
			return "", false
		}
		if inOneof {
			return trimSAIName(subType, true, false) + "List", true
		}
		return "repeated " + trimSAIName(subType, true, false), true
	},
	// TODO: Support these types
	"sai_acl_field_data_t":  func(next string, xmlInfo *protoGenInfo, inOneof bool) (string, bool) { return "AclFieldData", false },
	"sai_acl_action_data_t": func(next string, xmlInfo *protoGenInfo, inOneof bool) (string, bool) { return "AclActionData", false },
	"sai_pointer_t":         func(next string, xmlInfo *protoGenInfo, inOneof bool) (string, bool) { return "-", false },
}

// saiTypeToProtoType returns the protobuf type string for a SAI type.
// example: sai_u8_list_t -> repeated uint32
func saiTypeToProtoType(saiType string, xmlInfo *protoGenInfo, inOneof bool) (string, bool, error) {
	if pt, ok := saiTypeToProto[saiType]; ok {
		if pt.Repeated {
			if inOneof {
				return strcase.UpperCamelCase(pt.ProtoType + " List"), true, nil
			}
			return "repeated " + pt.ProtoType, true, nil
		}
		return pt.ProtoType, false, nil
	}
	if _, ok := xmlInfo.enums[saiType]; ok {
		return trimSAIName(saiType, true, false), false, nil
	}

	if splits := strings.Split(saiType, " "); len(splits) == 2 {
		fn, ok := saiTypeToProtoTypeCompound[splits[0]]
		if !ok {
			return "", false, fmt.Errorf("unknown sai type: %v", saiType)
		}
		name, isRepeated := fn(splits[1], xmlInfo, inOneof)
		return name, isRepeated, nil
	}

	return "", false, fmt.Errorf("unknown sai type: %v", saiType)
}

// trimSAIName trims sai_ prefix and _t from the string
func trimSAIName(name string, makeCamel, makeUpper bool) string {
	str := strings.TrimSuffix(strings.TrimPrefix(name, "sai_"), "_t")
	if makeCamel {
		str = strcase.UpperCamelCase(str)
	}
	if makeUpper {
		str = strings.ToUpper(str)
	}

	return str
}
