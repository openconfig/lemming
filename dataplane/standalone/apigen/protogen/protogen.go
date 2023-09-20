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

// Package protogen implements a generator for SAI to protobuf.
package protogen

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/openconfig/lemming/dataplane/standalone/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/standalone/apigen/saiast"

	log "github.com/golang/glog"
	strcase "github.com/stoewer/go-strcase"
)

// Generates returns a map of files containing the generated code code.
func Generate(doc *docparser.SAIInfo, sai *saiast.SAIAPI) (map[string]string, error) {
	files := map[string]string{}
	common, err := generateCommonTypes(doc)
	if err != nil {
		return nil, err
	}
	files["common.proto"] = common

	apis := map[string]*protoAPITmplData{}
	for _, iface := range sai.Ifaces {
		apiName := strings.TrimSuffix(strings.TrimPrefix(iface.Name, "sai_"), "_api_t")
		for _, fn := range iface.Funcs {
			meta := sai.GetFuncMeta(fn)
			if err := populateTmplDataFromFunc(apis, doc, apiName, meta); err != nil {
				return nil, err
			}
		}
		var builder strings.Builder
		if err := protoTmpl.Execute(&builder, apis[apiName]); err != nil {
			return nil, err
		}
		files[apiName+".proto"] = builder.String()
	}
	return files, nil
}

func rangeInOrder[T any](m map[string]T, pred func(key string, val T) error) error {
	keys := maps.Keys(m)
	slices.Sort(keys)
	for _, key := range keys {
		if err := pred(key, m[key]); err != nil {
			return err
		}
	}
	return nil
}

// generateCommonTypes returns all contents of the common proto.
// These all reside in the common.proto file to simplify handling imports.
func generateCommonTypes(docInfo *docparser.SAIInfo) (string, error) {
	common := &protoCommonTmplData{}

	// Generate the hand-crafted messages.
	rangeInOrder(saiTypeToProto, func(key string, typeInfo saiTypeInfo) error {
		if typeInfo.MessageDef != "" {
			common.Messages = append(common.Messages, typeInfo.MessageDef)
		}
		return nil
	})

	seenEnums := map[string]bool{}
	// Generate non-attribute enums.
	rangeInOrder(docInfo.Enums, func(name string, vals []string) error {
		protoName := saiast.TrimSAIName(name, true, false)
		unspecifiedName := saiast.TrimSAIName(name, false, true) + "_UNSPECIFIED"
		enum := &protoEnum{
			Name:   protoName,
			Values: []protoEnumValues{{Index: 0, Name: unspecifiedName}},
		}
		for i, val := range vals {
			name := strings.TrimPrefix(val, "SAI_")
			// If the SAI name conflicts with unspecified proto name, then add SAI prefix,
			// that way the proto enum value is always 1 greater than the c enum.
			if name == unspecifiedName {
				name = strings.TrimSuffix(saiast.TrimSAIName(name, false, true), "_UNSPECIFIED") + "_SAI_UNSPECIFIED"
			}
			enum.Values = append(enum.Values, protoEnumValues{
				Index: i + 1,
				Name:  name,
			})
		}
		if _, ok := seenEnums[protoName]; !ok {
			common.Enums = append(common.Enums, enum)
			seenEnums[protoName] = true
		}
		return nil
	})

	err := rangeInOrder(docInfo.Attrs, func(n string, attr *docparser.Attr) error {
		attrFields, err := createAttrs(1, n, docInfo, attr.ReadFields)
		if err != nil {
			return err
		}
		common.Lists = append(common.Lists, &protoTmplMessage{
			Fields: attrFields,
			Name:   saiast.TrimSAIName(n, true, false) + "Attribute",
		})
		return nil
	})
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	if err := protoCommonTmpl.Execute(&builder, common); err != nil {
		return "", err
	}
	return builder.String(), nil
}

// populateTmplDataFromFunc populatsd the protobuf template struct from a SAI function call.
func populateTmplDataFromFunc(apis map[string]*protoAPITmplData, docInfo *docparser.SAIInfo, apiName string, meta *saiast.FuncMetadata) error {
	if docInfo.Attrs[meta.TypeName] == nil {
		fmt.Printf("no doc info for type: %v\n", meta.TypeName)
		return nil
	}

	if _, ok := apis[apiName]; !ok {
		apis[apiName] = &protoAPITmplData{
			ServiceName: saiast.TrimSAIName(apiName, true, false),
		}
	}

	req := &protoTmplMessage{
		Name: strcase.UpperCamelCase(meta.Name + "_request"),
	}
	resp := &protoTmplMessage{
		Name: strcase.UpperCamelCase(meta.Name + "_response"),
	}

	idField := protoTmplField{
		Index:     1,
		ProtoType: "uint64",
		Name:      "oid",
	}
	if meta.Entry != "" {
		idField = protoTmplField{
			Index:     1,
			ProtoType: meta.Entry,
			Name:      "entry",
		}
	}

	// Handle proto generation
	switch meta.Operation {
	case "create":
		requestIdx := 1
		if meta.IsSwitchScoped {
			req.Fields = append(req.Fields, protoTmplField{
				Index:     requestIdx,
				ProtoType: "uint64",
				Name:      "switch",
			})
			requestIdx++
		} else if meta.Entry != "" {
			req.Fields = append(req.Fields, idField)
			requestIdx++
		}
		for _, v := range docInfo.Enums["sai_object_type_t"] {
			if v == fmt.Sprintf("SAI_OBJECT_TYPE_%s", meta.TypeName) {
				req.Option = fmt.Sprintf("option (sai_type) = OBJECT_TYPE_%s", meta.TypeName)
			}
		}
		if req.Option == "" {
			req.Option = "option (sai_type) = OBJECT_TYPE_UNSPECIFIED"
		}

		attrs, err := createAttrs(requestIdx, meta.TypeName, docInfo, docInfo.Attrs[meta.TypeName].CreateFields)
		if err != nil {
			return err
		}
		req.Fields = append(req.Fields, attrs...)
		if meta.Entry == "" { // Entries don't have id.
			resp.Fields = append(resp.Fields, idField)
		}
	case "set_attribute":
		// If there are no settable attributes, do nothing.
		if len(docInfo.Attrs[meta.TypeName].SetFields) == 0 {
			return nil
		}
		req.Fields = append(req.Fields, idField)
		attrs, err := createAttrs(2, meta.TypeName, docInfo, docInfo.Attrs[meta.TypeName].SetFields)
		if err != nil {
			return err
		}
		req.Fields = append(req.Fields, attrs...)
	case "get_attribute":
		req.Fields = append(req.Fields, idField, protoTmplField{
			ProtoType: "repeated " + strcase.UpperCamelCase(meta.TypeName+" attr"),
			Index:     2,
			Name:      "attr_type",
		})
		resp.Fields = append(resp.Fields, protoTmplField{
			Index:     1,
			Name:      "attr",
			ProtoType: strcase.UpperCamelCase(meta.TypeName + "Attribute"),
		})

		// attrEnum is the special emun that describes the possible values can be set/get for the API.
		attrEnum := protoEnum{
			Name:   strcase.UpperCamelCase(meta.TypeName + "_ATTR"),
			Values: []protoEnumValues{{Index: 0, Name: meta.TypeName + "_ATTR_UNSPECIFIED"}},
		}

		// For the attributes, generate code for the type if needed.
		for i, attr := range docInfo.Attrs[meta.TypeName].ReadFields {
			attrEnum.Values = append(attrEnum.Values, protoEnumValues{
				Index: i + 1,
				Name:  strings.TrimPrefix(attr.EnumName, "SAI_"),
			})
			// Handle function pointers as streaming RPCs.
			if strings.Contains(attr.SaiType, "sai_pointer_t") {
				funcName := strings.Split(attr.SaiType, " ")[1]
				name := saiast.TrimSAIName(strings.TrimSuffix(funcName, "_fn"), true, false)
				req := &protoTmplMessage{
					Name: strcase.UpperCamelCase(name + "_request"),
				}
				resp, ok := funcToStreamResp[funcName]
				if !ok {
					// TODO: There are 2 function pointers that don't follow this pattern, support them.
					log.Warningf("skipping unknown func type %q\n", funcName)
					continue
				}
				apis[apiName].Messages = append(apis[apiName].Messages, *req, resp)
				apis[apiName].RPCs = append(apis[apiName].RPCs, protoRPC{
					RequestName:  req.Name,
					ResponseName: "stream " + resp.Name,
					Name:         strcase.UpperCamelCase(name),
				})
			}
		}
		apis[apiName].Enums = append(apis[apiName].Enums, attrEnum)
	case "remove":
		req.Fields = append(req.Fields, idField)
	default:
		return nil
	}
	apis[apiName].Messages = append(apis[apiName].Messages, *req, *resp)
	apis[apiName].RPCs = append(apis[apiName].RPCs, protoRPC{
		RequestName:  req.Name,
		ResponseName: resp.Name,
		Name:         strcase.UpperCamelCase(meta.Name),
	})
	return nil
}

func createAttrs(startIdx int, typeName string, xmlInfo *docparser.SAIInfo, attrs []*docparser.AttrTypeName) ([]protoTmplField, error) {
	fields := []protoTmplField{}
	for _, attr := range attrs {
		// Function pointers are implemented as streaming RPCs instead of settable attributes.
		// TODO: Implement these.
		if strings.Contains(attr.SaiType, "sai_pointer_t") {
			continue
		}
		// Proto field names can't begin with numbers, prepend _.
		name := attr.MemberName
		if unicode.IsDigit(rune(attr.MemberName[0])) {
			name = fmt.Sprintf("_%s", name)
		}
		field := protoTmplField{
			Index: startIdx,
			Name:  name,
		}
		typ, repeated, err := SaiTypeToProtoType(attr.SaiType, xmlInfo)
		if err != nil {
			return nil, err
		}
		for i, val := range xmlInfo.Attrs[typeName].ReadFields {
			if val == attr {
				field.Option = fmt.Sprintf("[(attr_enum_value) = %d]", i+1)
			}
		}
		field.ProtoType = typ
		if !repeated {
			field.ProtoType = "optional " + typ
		}
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

option go_package = "github.com/openconfig/lemming/dataplane/standalone/proto";

{{ range .Enums }}
enum {{ .Name }} {
	{{- range .Values }}
	{{ .Name }} = {{ .Index }};
	{{- end}}
}
{{ end -}}

{{ range .Messages }}
message {{ .Name }} {
	{{- if .Option }}
	{{ .Option }};
	{{- end -}}
	{{- range .Fields }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }} {{- if .Option -}} {{ .Option }} {{- end -}};
	{{- end }}
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
	
import "google/protobuf/timestamp.proto";
import "google/protobuf/descriptor.proto";

option go_package = "github.com/openconfig/lemming/dataplane/standalone/proto";

extend google.protobuf.FieldOptions {
	optional int32 attr_enum_value = 50000;
}

extend google.protobuf.MessageOptions {
	optional ObjectType sai_type = 50001;
}
{{ range .Messages }}
{{ . }}
{{ end }}
message ObjectTypeQueryRequest {
	uint64 object = 1;
}

message ObjectTypeQueryResponse {
	ObjectType type = 1;
}

service Entrypoint {
  rpc ObjectTypeQuery(ObjectTypeQueryRequest) returns (ObjectTypeQueryResponse) {}
}
{{ range .Enums }}
enum {{ .Name }} {
	{{- range .Values }}
	{{ .Name }} = {{ .Index }};
	{{- end}}
}
{{ end -}}

{{- range .Lists }}
message {{ .Name }} {
	{{- range .Fields }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }} {{ if .Option -}} {{ .Option }} {{- end -}};
	{{- end }}
}
{{ end -}}
`))
)

// protoAPITmplData contains the formated information needed to render the protobuf template.
type protoAPITmplData struct {
	Messages    []protoTmplMessage
	RPCs        []protoRPC
	Enums       []protoEnum
	ServiceName string
}

// protoCommonTmplData contains the formated information needed to render the protobuf template.
type protoCommonTmplData struct {
	Messages []string
	Enums    []*protoEnum
	Lists    []*protoTmplMessage
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
	Name   string
	Option string
	Fields []protoTmplField
}

type protoTmplField struct {
	ProtoType string
	Name      string
	Index     int
	Option    string
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

var (
	saiTypeToProto = map[string]saiTypeInfo{
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
		"sai_json_t": {
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
		},
		"sai_timespec_t": {
			ProtoType: "google.protobuf.Timestamp",
		},
		// The non-scalar types could be autogenerated, but that aren't that many so create messages by hand.
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
		"sai_latch_status_t": {
			ProtoType: "LatchStatus",
			MessageDef: `message LatchStatus {
	bool current_status = 1;
	bool changed = 2;
}`,
		},
		"sai_port_lane_latch_status_list_t": {
			Repeated:  true,
			ProtoType: "PortLaneLatchStatus",
			MessageDef: `message PortLaneLatchStatus {
	uint32 lane = 1;
	LatchStatus value = 2;
}`,
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
}

message Uint64List {
	repeated uint64 list = 1;
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
		"sai_fdb_event_notification_data_t": {
			ProtoType: "FdbEventNotificationData",
			MessageDef: `
message FdbEventNotificationData {
    FdbEvent event_type = 1;
	FdbEntry fdb_entry = 2;
	repeated FdbEntryAttribute attrs = 3;
}`,
		},
		"sai_port_oper_status_notification_t": {
			ProtoType: "PortOperStatusNotification",
			MessageDef: `message PortOperStatusNotification {
	uint64 port_id = 1;
	PortOperStatus port_state = 2;
}`,
		},
		"sai_queue_deadlock_notification_data_t": {
			ProtoType: "QueueDeadlockNotificationData",
			MessageDef: `message QueueDeadlockNotificationData {
	uint64 queue_id = 1;
	QueuePfcDeadlockEventType event= 2;
	bool app_managed_recovery = 3;
}`,
		},
		"sai_bfd_session_state_notification_t": {
			ProtoType: "BfdSessionStateChangeNotificationData",
			MessageDef: `message BfdSessionStateChangeNotificationData {
	uint64 bfd_session_id = 1;
	BfdSessionState session_state = 2;
}`,
		},
		"sai_ipsec_sa_status_notification_t": {
			ProtoType: "IpsecSaStatusNotificationData",
			MessageDef: `message IpsecSaStatusNotificationData {
    uint64 ipsec_sa_id = 1;
	IpsecSaOctetCountStatus ipsec_sa_octet_count_status = 2;
	bool ipsec_egress_sn_at_max_limit = 3;
}`,
		},
	}
	// The notification function types are implemented as streaming RPCs.
	funcToStreamResp = map[string]protoTmplMessage{
		"sai_switch_state_change_notification_fn": {
			Name: "SwitchStateChangeNotificationResponse",
			Fields: []protoTmplField{{
				Index:     1,
				ProtoType: "uint64",
				Name:      "switch_id",
			}, {
				Index:     2,
				ProtoType: "SwitchOperStatus",
				Name:      "switch_oper_status",
			}},
		},
		"sai_switch_shutdown_request_notification_fn": {
			Name: "SwitchShutdownRequestNotificationResponse",
			Fields: []protoTmplField{{
				Index:     1,
				ProtoType: "uint64",
				Name:      "switch_id",
			}},
		},
		"sai_fdb_event_notification_fn": {
			Name: "FdbEventNotificationResponse",
			Fields: []protoTmplField{{
				Index:     1,
				ProtoType: "repeated FdbEventNotificationData",
				Name:      "data",
			}},
		},
		"sai_port_state_change_notification_fn": {
			Name: "PortStateChangeNotificationResponse",
			Fields: []protoTmplField{{
				Index:     1,
				ProtoType: "repeated PortOperStatusNotification",
				Name:      "data",
			}},
		},
		"sai_packet_event_notification_fn": {
			Name: "PacketEventNotificationResponse",
			Fields: []protoTmplField{{
				Index:     1,
				ProtoType: "uint64",
				Name:      "switch_id",
			}, {
				Index:     2,
				ProtoType: "bytes",
				Name:      "buffer",
			}, {
				Index:     3,
				ProtoType: "repeated HostifPacketAttribute",
				Name:      "attrs",
			}},
		},
		"sai_queue_pfc_deadlock_notification_fn": {
			Name: "QueuePfcDeadlockNotificationResponse",
			Fields: []protoTmplField{{
				Index:     1,
				ProtoType: "repeated QueueDeadlockNotificationData",
				Name:      "data",
			}},
		},
		"sai_bfd_session_state_change_notification_fn": {
			Name: "BfdSessionStateChangeNotificationResponse",
			Fields: []protoTmplField{{
				Index:     1,
				ProtoType: "repeated BfdSessionStateChangeNotificationData",
				Name:      "data",
			}},
		},
		"sai_tam_event_notification_fn": {
			Name: "TamEventNotificationResponse",
			Fields: []protoTmplField{{
				Index:     1,
				ProtoType: "uint64",
				Name:      "tam_event_id",
			}, {
				Index:     2,
				ProtoType: "bytes",
				Name:      "buffer",
			}, {
				Index:     3,
				ProtoType: "repeated TamEventActionAttribute",
				Name:      "attrs",
			}},
		},
		"sai_ipsec_sa_status_change_notification_fn": {
			Name: "IpsecSaStatusNotificationDataResponse",
			Fields: []protoTmplField{{
				Index:     1,
				ProtoType: "repeated IpsecSaStatusNotificationData",
				Name:      "data",
			}},
		},
	}
)

// saiTypeToProtoTypeCompound handles compound sai types (eg list of enums).
// The map key contains the base type (eg list) and func accepts the subtype (eg an enum type)
// and returns the full type string (eg repeated sample_enum).
var saiTypeToProtoTypeCompound = map[string]func(subType string, xmlInfo *docparser.SAIInfo) (string, bool){
	"sai_s32_list_t": func(subType string, xmlInfo *docparser.SAIInfo) (string, bool) {
		if _, ok := xmlInfo.Enums[subType]; !ok {
			return "", false
		}
		return "repeated " + saiast.TrimSAIName(subType, true, false), true
	},
	"sai_acl_field_data_t": func(next string, xmlInfo *docparser.SAIInfo) (string, bool) {
		return "AclFieldData", false
	},
	"sai_acl_action_data_t": func(next string, xmlInfo *docparser.SAIInfo) (string, bool) {
		return "AclActionData", false
	},
	"sai_pointer_t": func(next string, xmlInfo *docparser.SAIInfo) (string, bool) { return "-", false }, // Noop, these are special cases.
}

// saiTypeToProtoType returns the protobuf type string for a SAI type.
// example: sai_u8_list_t -> repeated uint32
func SaiTypeToProtoType(saiType string, xmlInfo *docparser.SAIInfo) (string, bool, error) {
	saiType = strings.TrimPrefix(saiType, "const ")

	if pt, ok := saiTypeToProto[saiType]; ok {
		if pt.Repeated {
			return "repeated " + pt.ProtoType, true, nil
		}
		return pt.ProtoType, false, nil
	}
	if _, ok := xmlInfo.Enums[saiType]; ok {
		return saiast.TrimSAIName(saiType, true, false), false, nil
	}

	if splits := strings.Split(saiType, " "); len(splits) == 2 {
		fn, ok := saiTypeToProtoTypeCompound[splits[0]]
		if !ok {
			return "", false, fmt.Errorf("unknown sai type: %v", saiType)
		}
		name, isRepeated := fn(splits[1], xmlInfo)
		return name, isRepeated, nil
	}

	return "", false, fmt.Errorf("unknown sai type: %v", saiType)
}
