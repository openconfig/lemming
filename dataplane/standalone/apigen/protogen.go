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

	strcase "github.com/stoewer/go-strcase"
)

// populateTmplDataFromFunc populatsd the protobuf template struct from a SAI function call.
func populateTmplDataFromFunc(protoTmplData *protoTmplData, funcName, entryType, operation, typeName, apiName string, isSwitchScoped bool, xmlInfo *protoGenInfo) error {
	msg := &protoTmplMessage{
		RequestName:  strcase.UpperCamelCase(funcName + "_request"),
		ResponseName: strcase.UpperCamelCase(funcName + "_response"),
	}
	nameTrimmed := strings.TrimSuffix(strings.TrimPrefix(apiName, "sai_"), "_api_t")
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
	// TODO: Enable proto generation.
	switch operation {
	case "create":
		requestIdx := 1
		if isSwitchScoped {
			msg.RequestFields = append(msg.RequestFields, protoTmplField{
				Index:     requestIdx,
				ProtoType: "uint64",
				Name:      "switch_id",
			})
		} else if entryType != "" {
			msg.RequestFields = append(msg.RequestFields, idField)
		}
		requestIdx++
		attrs, err := createAttrs(requestIdx, xmlInfo, xmlInfo.attrs[typeName].createFields, false)
		if err != nil {
			return err
		}
		msg.RequestAttrs = attrs

		msg.ResponseFields = append(msg.ResponseFields, protoTmplField{
			Index:     1,
			ProtoType: "uint64",
			Name:      "oid",
		})
	case "set_attribute":
		// If there are no settable attributes, do nothing.
		if len(xmlInfo.attrs[typeName].setFields) == 0 {
			return nil
		}
		msg.RequestFields = append(msg.RequestFields, idField)
		msg.RequestAttrsWrapperStart = "oneof attr {"
		msg.RequestAttrsWrapperEnd = "}"
		attrs, err := createAttrs(2, xmlInfo, xmlInfo.attrs[typeName].setFields, true)
		if err != nil {
			return err
		}
		msg.RequestAttrs = attrs
	case "get_attribute":
		msg.RequestFields = append(msg.RequestFields, idField, protoTmplField{
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
		for i, attr := range xmlInfo.attrs[typeName].readFields {
			attrEnum.Values = append(attrEnum.Values, protoEnumValues{
				Index: i + 1,
				Name:  strings.TrimPrefix(attr.EnumName, "SAI_"),
			})

			// For types that begin sai_api_name, generate them in the package.
			// TODO: Otherwise mark them for generation in a common package.
			if strings.HasPrefix(attr.SaiType, "sai_"+nameTrimmed) {
				protoName := trimSAIName(attr.SaiType, true, false)
				// TODO: Generated code for non-enum types.
				if vals, ok := xmlInfo.enums[attr.SaiType]; ok {
					enum := protoEnum{
						Name:   protoName,
						Values: []protoEnumValues{{Index: 0, Name: trimSAIName(attr.SaiType, false, true) + "_UNSPECIFIED"}},
					}
					for i, val := range vals {
						enum.Values = append(enum.Values, protoEnumValues{
							Index: i + 1,
							Name:  strings.TrimPrefix(val, "SAI_"),
						})
					}
					protoTmplData.Enums[protoName] = enum
				}
			}
		}
		protoTmplData.Enums[attrEnum.Name] = attrEnum

		attrs, err := createAttrs(1, xmlInfo, xmlInfo.attrs[typeName].readFields, false)
		if err != nil {
			return err
		}
		msg.ResponseAttrs = attrs
	default:
		return nil
	}
	protoTmplData.Messages = append(protoTmplData.Messages, *msg)
	protoTmplData.RPCs = append(protoTmplData.RPCs, protoRPC{
		RequestName:  msg.RequestName,
		ResponseName: msg.ResponseName,
		Name:         strcase.UpperCamelCase(funcName),
	})
	return nil
}

func createAttrs(startIdx int, xmlInfo *protoGenInfo, attrs []attrTypeName, inOneof bool) ([]protoTmplField, error) {
	fields := []protoTmplField{}
	for _, attr := range attrs {
		field := protoTmplField{
			Index: startIdx,
			Name:  attr.MemberName,
		}
		typ, err := saiTypeToProtoType(attr.SaiType, xmlInfo, inOneof)
		if err != nil {
			return nil, err
		}
		field.ProtoType = typ
		fields = append(fields, field)
		startIdx++
	}
	return fields, nil
}

// TODO: Enable generation.
var protoTmpl = template.Must(template.New("cc").Parse(`
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
message {{ .RequestName }} {
	{{- range .RequestFields }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }};
	{{- end }}
	{{ .RequestAttrsWrapperStart -}}
	{{- range .RequestAttrs }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }};
	{{- end }}
	{{ .RequestAttrsWrapperEnd }}
}

message {{ .ResponseName }} {
	{{- range .ResponseFields }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }};
	{{- end }}
	{{ .ResponseAttrsWrapperStart -}}
	{{- range .ResponseAttrs }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }};
	{{- end }}
	{{ .ResponseAttrsWrapperEnd }}
}

{{ end }}


service {{ .ServiceName }} {
	{{- range .RPCs }}
	rpc {{ .Name }} ({{ .RequestName }}) returns ({{ .ResponseName }}) {}
	{{- end }}
}
`))

// protoTmplData contains the formated information needed to render the protobuf template.
type protoTmplData struct {
	Messages    []protoTmplMessage
	RPCs        []protoRPC
	Enums       map[string]protoEnum
	ServiceName string
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
	RequestName               string
	ResponseName              string
	RequestAttrsWrapperStart  string
	RequestAttrsWrapperEnd    string
	RequestFields             []protoTmplField
	RequestAttrs              []protoTmplField
	ResponseAttrsWrapperStart string
	ResponseAttrsWrapperEnd   string
	ResponseFields            []protoTmplField
	ResponseAttrs             []protoTmplField
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
	"sai_map_list_t": {
		Repeated:  true,
		ProtoType: "map<uint32, uint32>",
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
	PacketColor color 7;
	uint32 mpls_exp;
	uint32 fc;
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
	uint32 num_voq 6;
}`,
	},
	"sai_system_port_config_list_t": {
		Repeated:  true,
		ProtoType: "SystemPortConfig",
	},
	"sai_ip_address_list_t": {
		Repeated:  true,
		ProtoType: "IpAddress",
	},
	"sai_port_eye_values_list_t": {
		Repeated:  true,
		ProtoType: "PortEyeValues",
		MessageDef: `message PortLaneEyeValues {
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
	uint32 error_count 2;
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
	BindPointType bind_point = 2;
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
		uint64 uint = 2;
		uint64 int = 3;
		bytes mac = 4;
		bytes ip = 5;
		Uint64List list = 6;
	};
	oneof data {
		bool booldata = 7;
		uint64 uint = 8;
		int64 int = 9;
		bytes mac = 10;
		bytes ip = 11;
		Uint64List list = 12;
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
		IpAddress ipaddr = 8;
	};
}`,
	},
}

// saiTypeToProtoTypeCompound handles compound sai types (eg list of enums).
// The map key contains the base type (eg list) and func accepts the subtype (eg an enum type)
// and returns the full type string (eg repeated sample_enum).
var saiTypeToProtoTypeCompound = map[string]func(subType string, xmlInfo *protoGenInfo, inOneof bool) string{
	"sai_s32_list_t": func(subType string, xmlInfo *protoGenInfo, inOneof bool) string {
		if _, ok := xmlInfo.enums[subType]; !ok {
			return ""
		}
		if inOneof {
			return "List" + trimSAIName(subType, true, false)
		}
		return "repeated " + trimSAIName(subType, true, false)
	},
	// TODO: Support these types
	"sai_acl_field_data_t":  func(next string, xmlInfo *protoGenInfo, inOneof bool) string { return "AclFieldData" },
	"sai_acl_action_data_t": func(next string, xmlInfo *protoGenInfo, inOneof bool) string { return "AclActionData" },
	"sai_pointer_t":         func(next string, xmlInfo *protoGenInfo, inOneof bool) string { return "-" },
}

// saiTypeToProtoType returns the protobuf type string for a SAI type.
// example: sai_u8_list_t -> repeated uint32
func saiTypeToProtoType(saiType string, xmlInfo *protoGenInfo, inOneof bool) (string, error) {
	if protoType, ok := saiTypeToProto[saiType]; ok {
		if protoType.Repeated {
			return "repeated " + protoType.ProtoType, nil
		}
		return protoType.ProtoType, nil
	}
	if _, ok := xmlInfo.enums[saiType]; ok {
		return trimSAIName(saiType, true, false), nil
	}

	if splits := strings.Split(saiType, " "); len(splits) == 2 {
		fn, ok := saiTypeToProtoTypeCompound[splits[0]]
		if !ok {
			return "", fmt.Errorf("unknown sai type: %v", saiType)
		}
		return fn(splits[1], xmlInfo, inOneof), nil
	}

	return "", fmt.Errorf("unknown sai type: %v", saiType)
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
