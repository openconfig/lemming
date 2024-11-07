// Copyright 2024 Google LLC
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

package typeinfo

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/openconfig/lemming/dataplane/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/apigen/saiast"

	log "github.com/golang/glog"
	strcase "github.com/stoewer/go-strcase"
)

type TemplateData struct {
	APIs map[string]*APITemplate
}

type APITemplate struct {
	Messages       []ProtoTmplMessage
	RPCs           []ProtoRPC
	Enums          []ProtoEnum
	ServiceName    string
	ProtoPackage   string
	ProtoGoPackage string
	ProtoOutDir    string
	IncludeGuard   string
	Header         string
	ProtoInclude   string
	GRPCInclude    string
	APIType        string
	ProtoClass     string
	APIName        string
	Globals        []string
	Funcs          []*TemplateFunc
	ConvertFuncs   []*TemplateFunc
	CCOutDir       string
}

func Data(doc *docparser.SAIInfo, sai *saiast.SAIAPI, protoPackage, protoGoPackage, ccOutDir, protoOutDir string) (*TemplateData, error) {
	data := &TemplateData{
		APIs: map[string]*APITemplate{},
	}
	for _, iface := range sai.Ifaces {
		apiName := strings.TrimSuffix(strings.TrimPrefix(iface.Name, "sai_"), "_api_t")
		if _, ok := data.APIs[apiName]; !ok {
			data.APIs[apiName] = &APITemplate{
				ServiceName:    saiast.TrimSAIName(apiName, true, false),
				Enums:          []ProtoEnum{},
				ProtoPackage:   protoPackage,
				ProtoGoPackage: protoGoPackage,
				ProtoOutDir:    protoOutDir,
				IncludeGuard:   fmt.Sprintf("DATAPLANE_STANDALONE_SAI_%s_H_", strings.ToUpper(apiName)),
				Header:         fmt.Sprintf("%s.h", apiName),
				APIType:        iface.Name,
				APIName:        apiName,
				ProtoInclude:   apiName + ".pb",
				CCOutDir:       ccOutDir,
			}
		}
		switch apiName {
		case "switch":
			data.APIs[apiName].Globals = append(data.APIs[apiName].Globals, "std::unique_ptr<PortStateReactor> port_state;")
		}

		for _, fn := range iface.Funcs {
			meta := sai.GetFuncMeta(fn)
			if err := populateTmplDataFromFunc(data.APIs, doc, apiName, meta); err != nil {
				return nil, err
			}
			opFn, convertFn := createCCData(meta, apiName, sai, doc, fn)
			if opFn != nil {
				data.APIs[apiName].Funcs = append(data.APIs[apiName].Funcs, opFn)
			}
			if convertFn != nil {
				data.APIs[apiName].ConvertFuncs = append(data.APIs[apiName].ConvertFuncs, convertFn)
			}
		}
	}
	return data, nil
}

// createCCData returns a two structs with the template data for the given function.
// The first is the implementation of the API: CreateFoo.
// The second is the a conversion func from attribute list to the proto message. covert_create_foo.
func createCCData(meta *saiast.FuncMetadata, apiName string, sai *saiast.SAIAPI, info *docparser.SAIInfo, fn *saiast.TypeDecl) (*TemplateFunc, *TemplateFunc) {
	if info.Attrs[meta.TypeName] == nil {
		fmt.Printf("no doc info for type: %v\n", meta.TypeName)
		return nil, nil
	}
	opFn := &TemplateFunc{
		ReturnType: sai.Funcs[fn.Typ].ReturnType,
		Name:       meta.Name,
		Operation:  meta.Operation,
		TypeName:   meta.TypeName,
		ReqType:    strcase.UpperCamelCase(meta.Name + "_request"),
		RespType:   strcase.UpperCamelCase(meta.Name + "_response"),
	}
	convertFn := &TemplateFunc{
		Name:      "convert_" + meta.Name,
		Operation: meta.Operation,
		TypeName:  meta.TypeName,
		ReqType:   strcase.UpperCamelCase(meta.Name + "_request"),
		RespType:  strcase.UpperCamelCase(meta.Name + "_response"),
	}

	var paramDefs []string
	var paramVars []string
	for _, param := range sai.Funcs[fn.Typ].Params {
		paramDefs = append(paramDefs, fmt.Sprintf("%s %s", param.Typ, param.Name))
		name := strings.ReplaceAll(param.Name, "*", "")
		// Functions that operator on entries take some entry type instead of an object id as argument.
		// Generate a entry union with the pointer to entry instead.
		if strings.Contains(param.Typ, "entry") {
			opFn.Entry = fmt.Sprintf("common_entry_t entry = {.%s = %s};", name, name)
			name = "entry"
		}
		paramVars = append(paramVars, name)
	}
	opFn.Args = strings.Join(paramDefs, ", ")
	opFn.Vars = strings.Join(paramVars[1:], ", ")
	convertFn.Args = strings.Join(paramDefs[1:], ", ")
	convertFn.Vars = strings.Join(paramVars[1:], ", ")
	opFn.Client = strcase.SnakeCase(apiName)
	if opFn.Client == "switch" { // switch is C++ keyword.
		opFn.Client = "switch_"
	}
	opFn.RPCMethod = strcase.UpperCamelCase(meta.Name)
	opFn.SwitchScoped = meta.IsSwitchScoped
	opFn.AttrEnumType = strcase.UpperCamelCase(meta.TypeName + " attr")
	opFn.AttrType = strcase.SnakeCase("sai_" + meta.TypeName + "_attr_t")

	// If the func has entry, then we don't use ids, instead pass the entry to the proto.
	if meta.Entry == "" {
		opFn.OidVar = sai.Funcs[fn.Typ].Params[0].Name
		opFn.OidPointer = strings.TrimPrefix(opFn.OidVar, "*")
	} else {
		i := 0
		if strings.Contains(opFn.Operation, "bulk") {
			i = 1
		}
		entryType := strings.TrimPrefix(sai.Funcs[fn.Typ].Params[i].Typ, "const ")
		if ua, ok := typeToUnionAccessor[entryType]; ok {
			opFn.EntryConversionFunc = ua.convertFromFunc
			opFn.EntryVar = sai.Funcs[fn.Typ].Params[i].Name
		}
	}

	switch opFn.Operation {
	case createOp:
		convertFn.AttrSwitch = &AttrSwitch{
			Var:      "attr_list[i].id",
			ProtoVar: "msg",
		}
		opFn.ConvertFunc = strcase.SnakeCase("convert_create " + meta.TypeName)
		convertFn.ReturnType = opFn.ReqType
		for _, attr := range info.Attrs[meta.TypeName].CreateFields {
			name := sanitizeProtoName(attr.MemberName)
			smt, err := protoFieldSetter(attr.SaiType, convertFn.AttrSwitch.ProtoVar, name, "attr_list[i].value", info)
			if err != nil {
				fmt.Println("skipping due to error: ", err)
				continue
			}
			smt.EnumValue = attr.EnumName
			convertFn.AttrSwitch.Attrs = append(convertFn.AttrSwitch.Attrs, smt)
		}
	case getAttrOp:
		opFn.AttrSwitch = &AttrSwitch{
			Var:      "attr_list[i].id",
			ProtoVar: "resp.attr()",
		}
		convertFn = nil
		for _, attr := range info.Attrs[meta.TypeName].ReadFields {
			name := sanitizeProtoName(attr.MemberName)
			smt, err := protoFieldGetter(attr.SaiType, name, "attr_list[i].value", info)
			if err != nil {
				fmt.Println("skipping due to error: ", err)
				continue
			}
			smt.EnumValue = attr.EnumName
			opFn.AttrSwitch.Attrs = append(opFn.AttrSwitch.Attrs, smt)
		}
	case setAttrOp:
		convertFn = nil
		opFn.AttrSwitch = &AttrSwitch{
			Var:      "attr->id",
			ProtoVar: "req",
		}
		for _, attr := range info.Attrs[meta.TypeName].SetFields {
			name := sanitizeProtoName(attr.MemberName)
			smt, err := protoFieldSetter(attr.SaiType, opFn.AttrSwitch.ProtoVar, name, "attr->value", info)
			if err != nil {
				fmt.Println("skipping due to error: ", err)
				continue
			}
			smt.EnumValue = attr.EnumName
			opFn.AttrSwitch.Attrs = append(opFn.AttrSwitch.Attrs, smt)
		}
	case "create_bulk":
		convertFn = nil
		opFn.EntryVar = strings.TrimPrefix(opFn.EntryVar, "*") // Usual entry is pointer, but for remove_bulk it's an array.
		opFn.ConvertFunc = strcase.SnakeCase("convert_create " + meta.TypeName)
		for _, attr := range info.Attrs[meta.TypeName].CreateFields {
			name := sanitizeProtoName(attr.MemberName)
			smt, err := protoFieldSetter(attr.SaiType, "", name, "attr_list[i].value", info)
			if err != nil {
				fmt.Println("skipping due to error: ", err)
				continue
			}
			smt.EnumValue = attr.EnumName
		}
	case "remove_bulk":
		convertFn = nil
		opFn.EntryVar = strings.TrimPrefix(opFn.EntryVar, "*") // Usual entry is pointer, but for remove_bulk it's an array.
		for _, attr := range info.Attrs[meta.TypeName].CreateFields {
			name := sanitizeProtoName(attr.MemberName)
			smt, err := protoFieldSetter(attr.SaiType, "", name, "attr_list[i].value", info)
			if err != nil {
				fmt.Println("skipping due to error: ", err)
				continue
			}
			smt.EnumValue = attr.EnumName
		}
	case "get_stats":
		convertFn = nil
		opFn.AttrType = strcase.SnakeCase("sai_" + meta.TypeName + "_stat_t")
		opFn.AttrEnumType = strcase.UpperCamelCase(meta.TypeName + " stat")
	default:
		convertFn = nil
	}

	// Patches for non-standard APIS
	if meta.TypeName == "ACL_TABLE" {
		switch meta.Operation {
		case createOp:
			convertFn.AttrConvertInsert = `
if (attr_list[i].id >= SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
  (*msg.mutable_user_defined_field_group_min())[attr_list[i].id - SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN] = attr_list[i].value.oid;
}`
		case getAttrOp:
			opFn.AttrConvertInsert = `
if (attr_list[i].id >= SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
  attr_list[i].value.oid = resp.attr().user_defined_field_group_min().at(attr_list[i].id - SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN);
}`
		}
	}
	if meta.TypeName == "ACL_ENTRY" {
		switch meta.Operation {
		case createOp:
			convertFn.AttrConvertInsert = `
if (attr_list[i].id >= SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
    *(*msg.mutable_user_defined_field_group_min())[attr_list[i].id  - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN].mutable_data_u8list() = std::string(attr_list[i].value.aclfield.data.u8list.list, attr_list[i].value.aclfield.data.u8list.list + attr_list[i].value.aclfield.data.u8list.count);
    *(*msg.mutable_user_defined_field_group_min())[attr_list[i].id  - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN].mutable_mask_u8list() = std::string(attr_list[i].value.aclfield.mask.u8list.list, attr_list[i].value.aclfield.mask.u8list.list + attr_list[i].value.aclfield.mask.u8list.count);
}`
		case setAttrOp:
			opFn.AttrConvertInsert = `
if (attr->id >= SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr->id < SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
*(*req.mutable_user_defined_field_group_min())[attr->id  - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN].mutable_data_u8list() = std::string(attr->value.aclfield.data.u8list.list, attr->value.aclfield.data.u8list.list + attr->value.aclfield.data.u8list.count);
*(*req.mutable_user_defined_field_group_min())[attr->id  - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN].mutable_mask_u8list() = std::string(attr->value.aclfield.mask.u8list.list, attr->value.aclfield.mask.u8list.list + attr->value.aclfield.mask.u8list.count);
}`
		case getAttrOp:
			opFn.AttrConvertInsert = `
if (attr_list[i].id >= SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
	auto acl_attr = resp.attr().user_defined_field_group_min().at(attr_list[i].id - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN);
	memcpy(attr_list[i].value.aclfield.data.u8list.list, acl_attr.data_u8list().data(), acl_attr.data_u8list().size());
	memcpy(attr_list[i].value.aclfield.mask.u8list.list, acl_attr.mask_u8list().data(), acl_attr.mask_u8list().size());
    attr_list[i].value.aclfield.data.u8list.count = acl_attr.data_u8list().size();
	attr_list[i].value.aclfield.mask.u8list.count = acl_attr.mask_u8list().size();
}`
		}
	}

	opFn.UseCommonAPI = supportedOperation[opFn.Operation]
	// Function or types that don't follow standard naming.
	if strings.Contains(opFn.TypeName, "PORT_ALL") || strings.Contains(opFn.TypeName, "ALL_NEIGHBOR") {
		opFn.UseCommonAPI = false
	}
	return opFn, convertFn
}

func sanitizeProtoName(inName string) string {
	name := strings.ReplaceAll(inName, "inline", "inline_") // inline is C++ keyword
	if unicode.IsDigit(rune(name[0])) {
		name = fmt.Sprintf("_%s", name)
	}
	return name
}

const (
	createOp  = "create"
	getAttrOp = "get_attribute"
	setAttrOp = "set_attribute"
)

type accessorType int

const (
	scalar accessorType = iota
	fixedSizedArray
	variableSizedArray
	convertFunc
	callbackRPC
	acl
)

type unionAccessor struct {
	accessor        string
	protoAccessor   string
	pointerOf       bool
	aType           accessorType
	convertFromFunc string
	convertToFunc   string
	convertToCopy   bool // If there is preallocated list, need to copy elems into it.
	assignmentVar   string
}

var typeToUnionAccessor = map[string]*unionAccessor{
	"sai_object_list_t": {
		accessor: "objlist",
		aType:    variableSizedArray,
	},
	"sai_s32_list_t": {
		accessor: "s32list",
		aType:    variableSizedArray,
	},
	"sai_u32_list_t": {
		accessor: "u32list",
		aType:    variableSizedArray,
	},
	"sai_u8_list_t": {
		accessor: "u8list",
		aType:    variableSizedArray,
	},
	"sai_s8_list_t": {
		accessor: "s8list",
		aType:    variableSizedArray,
	},
	"sai_mac_t": {
		accessor: "mac",
		aType:    fixedSizedArray,
	},
	"sai_ip4_t": {
		accessor:  "ip4",
		aType:     fixedSizedArray,
		pointerOf: true,
	},
	"sai_ip6_t": {
		accessor: "ip6",
		aType:    fixedSizedArray,
	},
	"sai_object_id_t": {
		accessor: "oid",
		aType:    scalar,
	},
	"sai_uint64_t": {
		accessor: "u64",
		aType:    scalar,
	},
	"sai_uint32_t": {
		accessor: "u32",
		aType:    scalar,
	},
	"sai_uint16_t": {
		accessor: "u16",
		aType:    scalar,
	},
	"sai_uint8_t": {
		accessor: "u8",
		aType:    scalar,
	},
	"sai_int8_t": {
		accessor: "s8",
		aType:    scalar,
	},
	"bool": {
		accessor: "booldata",
		aType:    scalar,
	},
	"char": {
		accessor: "chardata",
		aType:    scalar,
	},
	"sai_ip_address_t": {
		accessor:        "ipaddr",
		convertFromFunc: "convert_from_ip_address",
		convertToFunc:   "convert_to_ip_address",
		aType:           convertFunc,
	},
	"sai_route_entry_t": {
		convertFromFunc: "convert_from_route_entry",
		convertToFunc:   "convert_to_route_entry",
		aType:           convertFunc,
	},
	"sai_neighbor_entry_t": {
		convertFromFunc: "convert_from_neighbor_entry",
		convertToFunc:   "convert_to_neighbor_entry",
		aType:           convertFunc,
	},
	"sai_pointer_t sai_port_state_change_notification_fn": {
		aType:           callbackRPC,
		assignmentVar:   "port_state",
		convertFromFunc: "std::make_unique<PortStateReactor>",
	},
	"sai_acl_capability_t": {
		accessor:        "aclcapability",
		aType:           convertFunc,
		convertToCopy:   true,
		convertFromFunc: "convert_from_acl_capability",
		convertToFunc:   "convert_to_acl_capability",
	},
	"sai_acl_field_data_t sai_ip4_t": {
		accessor:        "ip4",
		convertFromFunc: "convert_from_acl_field_data_ip4",
		convertToFunc:   "convert_to_acl_field_data",
		protoAccessor:   "ip",
		aType:           acl,
	},
	"sai_acl_action_data_t sai_object_id_t": {
		accessor:        "oid",
		convertFromFunc: "convert_from_acl_action_data",
		convertToFunc:   "convert_to_acl_action_data",
		protoAccessor:   "oid",
		aType:           acl,
	},
	"sai_acl_action_data_t sai_packet_action_t": {
		accessor:        "s32",
		convertFromFunc: "convert_from_acl_action_data_action",
		convertToFunc:   "convert_to_acl_action_data_action",
		protoAccessor:   "packet_action",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_acl_ip_type_t": {
		accessor:        "s32",
		convertFromFunc: "convert_from_acl_field_data_ip_type",
		convertToFunc:   "convert_to_acl_field_data_ip_type",
		protoAccessor:   "ip_type",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_uint8_t": {
		accessor:        "u8",
		convertFromFunc: "convert_from_acl_field_data",
		convertToFunc:   "convert_to_acl_field_data_u8",
		protoAccessor:   "uint",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_uint16_t": {
		accessor:        "u16",
		convertFromFunc: "convert_from_acl_field_data",
		convertToFunc:   "convert_to_acl_field_data_u16",
		protoAccessor:   "uint",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_uint32_t": {
		accessor:        "u32",
		convertFromFunc: "convert_from_acl_field_data",
		convertToFunc:   "convert_to_acl_field_data_u32",
		protoAccessor:   "uint",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_ip6_t": {
		accessor:        "ip6",
		convertFromFunc: "convert_from_acl_field_data_ip6",
		convertToFunc:   "convert_to_acl_field_data_ip6",
		protoAccessor:   "ip",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_mac_t": {
		accessor:        "mac",
		convertFromFunc: "convert_from_acl_field_data_mac",
		convertToFunc:   "convert_to_acl_field_data_mac",
		protoAccessor:   "mac",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_object_id_t": {
		accessor:        "oid",
		convertFromFunc: "convert_from_acl_field_data",
		convertToFunc:   "convert_to_acl_field_data",
		protoAccessor:   "oid",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_u8_list_t": {
		accessor:        "u8list",
		convertFromFunc: "convert_from_acl_field_data",
		convertToFunc:   "convert_to_acl_field_data_u8",
		protoAccessor:   "sai_u8_list_t",
		aType:           acl,
	},
}

var supportedOperation = map[string]bool{
	createOp:             true,
	"remove":             true,
	getAttrOp:            true,
	setAttrOp:            true,
	"clear_stats":        true,
	"get_stats":          true,
	"get_stats_ext":      true,
	"create_bulk":        true,
	"remove_bulk":        true,
	"set_attribute_bulk": false,
	"get_attribute_bulk": false,
}

func protoFieldSetter(saiType, protoVar, protoField, varName string, info *docparser.SAIInfo) (*AttrSwitchSmt, error) {
	smt := &AttrSwitchSmt{
		ProtoFunc: fmt.Sprintf("set_%s", protoField),
	}

	if _, ok := info.Enums[saiType]; ok {
		smt.Args = fmt.Sprintf("convert_%s_to_proto(%s.s32)", saiType, varName)
		return smt, nil
	}

	ua, ok := typeToUnionAccessor[saiType]
	if !ok {
		split := strings.Split(saiType, " ")
		if len(split) != 2 {
			return nil, fmt.Errorf("unknown sai type: %q", saiType)
		}
		if _, ok := info.Enums[split[1]]; split[0] == "sai_s32_list_t" && ok {
			smt.ProtoFunc = fmt.Sprintf("mutable_%s()->CopyFrom", protoField)
			smt.Args = fmt.Sprintf("convert_list_%s_to_proto(%s.%s)", split[1], varName, "s32list")
			return smt, nil
		} else {
			return nil, fmt.Errorf("unknown sai type: %q", saiType)
		}
	}

	switch ua.aType {
	case scalar:
		smt.Args = fmt.Sprintf("%s.%s", varName, ua.accessor)
	case convertFunc:
		smt.Args = fmt.Sprintf("%s(%s.%s)", ua.convertFromFunc, varName, ua.accessor)
	case fixedSizedArray:
		smt.Args = fmt.Sprintf("%s.%s, sizeof(%s.%s)", varName, ua.accessor, varName, ua.accessor)
		if ua.pointerOf {
			smt.Args = fmt.Sprintf("&%s.%s, sizeof(%s.%s)", varName, ua.accessor, varName, ua.accessor)
		}
	case variableSizedArray:
		smt.ProtoFunc = fmt.Sprintf("mutable_%s()->Add", protoField)
		smt.Args = fmt.Sprintf("%s.%s.list, %s.%s.list + %s.%s.count", varName, ua.accessor, varName, ua.accessor, varName, ua.accessor)
	case callbackRPC:
		smt.Var = ua.assignmentVar
		smt.ConvertFunc = ua.convertFromFunc
		fnType := strings.Split(saiType, " ")[1]
		smt.Args = fmt.Sprintf("switch_, reinterpret_cast<%s>(%s.ptr)", fnType, varName)
	case acl:
		smt.Var = fmt.Sprintf("*%s.mutable_%s()", protoVar, protoField)
		smt.ConvertFunc = ua.convertFromFunc
		access := "aclaction"
		smt.Args = fmt.Sprintf("%s.%s, %s.%s.parameter.%s", varName, access, varName, access, ua.accessor)
		if strings.Contains(saiType, "sai_acl_field_data_t") {
			access = "aclfield"
			smt.Args = fmt.Sprintf("%s.%s, %s.%s.data.%s, %s.%s.mask.%s", varName, access, varName, access, ua.accessor, varName, access, ua.accessor)
			if strings.Contains(saiType, "sai_object_id_t") {
				smt.Args = fmt.Sprintf("%s.%s, %s.%s.data.%s", varName, access, varName, access, ua.accessor)
			}
		}
	default:
		return nil, fmt.Errorf("unknown accessor type %q", ua.aType)
	}
	return smt, nil
}

type AttrSwitchSmt struct {
	// EnumValue is the name of enum value string. (eg SAI_HOSTIF_ATTR_OPER_STATUS).
	EnumValue string
	// ProtoFunc is the name of the protobuf getter or setter (eg set_obj_id() or obj_id()).
	ProtoFunc string
	// Args are the arguments to pass the ProtoFunc as comma seperated values.
	Args string
	// Var is the name of the variable that is used for assignment.
	Var string
	// ConvertFunc is a name the function that is called before assignment to var.
	ConvertFunc string
	// ConvertFuncArgs are the extra arguments to pass the convert func, in addition to ProtoFunc.
	ConvertFuncArgs string
	// CopyConvertFunc is the name of func to use. If set, then this function is used instead assignment.
	// For example, to copy a string, set this strncpy, since char arrays can't assigned be directly.
	CopyConvertFunc string
	// ConvertFuncArgs are the extra arguments to pass the copy convert func.
	CopyConvertFuncArgs string
	// CustomText is any text to include for in the case.
	CustomText string
}

type AttrSwitch struct {
	Var      string
	Attrs    []*AttrSwitchSmt
	ProtoVar string
}

type TemplateFunc struct {
	ReturnType          string
	Name                string
	Args                string
	TypeName            string
	Operation           string
	Vars                string
	UseCommonAPI        bool
	Entry               string
	AttrSwitch          *AttrSwitch
	ReqType             string
	RespType            string
	Client              string
	RPCMethod           string
	OidVar              string
	OidPointer          string
	AttrType            string
	AttrEnumType        string
	SwitchScoped        bool
	EntryConversionFunc string
	EntryVar            string
	ConvertFunc         string
	AttrConvertInsert   string
}

func protoFieldGetter(saiType, protoField, varName string, info *docparser.SAIInfo) (*AttrSwitchSmt, error) {
	smt := &AttrSwitchSmt{
		ProtoFunc: protoField,
	}
	if _, ok := info.Enums[saiType]; ok {
		smt.Var = varName + ".s32"
		smt.ConvertFunc = fmt.Sprintf("convert_%s_to_sai", saiType)
		smt.ConvertFuncArgs = ""
		return smt, nil
	}
	ua, ok := typeToUnionAccessor[saiType]
	if !ok {
		split := strings.Split(saiType, " ")
		if len(split) != 2 {
			return nil, fmt.Errorf("unknown sai type: %q", saiType)
		}
		if _, ok := info.Enums[split[1]]; split[0] == "sai_s32_list_t" && ok {
			smt.Var = fmt.Sprintf("%s.%s", varName, "s32list")
			smt.CopyConvertFuncArgs = fmt.Sprintf(", &%s.count", smt.Var)
			smt.CopyConvertFunc = fmt.Sprintf("convert_list_%s_to_sai", split[1])
			smt.Var += ".list"
			return smt, nil
		} else {
			return nil, fmt.Errorf("unknown sai type: %q", saiType)
		}
	}
	smt.Var = fmt.Sprintf("%s.%s", varName, ua.accessor)

	switch ua.aType {
	case scalar:
		if saiType == "char" {
			smt.CopyConvertFunc = "strncpy"
			smt.CopyConvertFuncArgs = ".data(), 32"
		}
		return smt, nil
	case convertFunc:
		if ua.convertToCopy {
			smt.CopyConvertFunc = ua.convertToFunc
		} else {
			smt.ConvertFunc = ua.convertToFunc
		}
		return smt, nil
	case fixedSizedArray:
		if saiType == "sai_ip4_t" {
			smt.Var = fmt.Sprintf("&%s.%s", varName, ua.accessor)
		}
		smt.CopyConvertFunc = "memcpy"
		smt.CopyConvertFuncArgs = fmt.Sprintf(".data(), sizeof(%s)", saiType)
		return smt, nil
	case variableSizedArray:
		smt.CopyConvertFunc = "copy_list"
		smt.CopyConvertFuncArgs = fmt.Sprintf(", &%s.count", smt.Var)
		smt.Var += ".list"
		return smt, nil
	case acl:
		smt.ConvertFunc = ua.convertToFunc
		access := "aclaction"
		smt.ConvertFuncArgs = fmt.Sprintf(", resp.attr().%s().%s()", smt.ProtoFunc, ua.protoAccessor)
		if strings.Contains(saiType, "sai_acl_field_data_t") {
			access = "aclfield"
			smt.ConvertFuncArgs = fmt.Sprintf(", resp.attr().%s().data_%s(), resp.attr().%s().mask_%s()", smt.ProtoFunc, ua.protoAccessor, smt.ProtoFunc, ua.protoAccessor)
			if strings.Contains(saiType, "sai_object_id_t") || strings.Contains(saiType, "sai_acl_ip_type_t") {
				smt.ConvertFuncArgs = fmt.Sprintf(", resp.attr().%s().data_%s()", smt.ProtoFunc, ua.protoAccessor)
			}
		}
		smt.Var = varName + "." + access
		return smt, nil
	}
	return nil, fmt.Errorf("unknown accessor type %q", saiType)
}

// populateTmplDataFromFunc populatsd the protobuf template struct from a SAI function call.
func populateTmplDataFromFunc(apis map[string]*APITemplate, docInfo *docparser.SAIInfo, apiName string, meta *saiast.FuncMetadata) error {
	if docInfo.Attrs[meta.TypeName] == nil {
		fmt.Printf("no doc info for type: %v\n", meta.TypeName)
		return nil
	}

	req := &ProtoTmplMessage{
		Name: strcase.UpperCamelCase(meta.Name + "_request"),
	}
	resp := &ProtoTmplMessage{
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
			if v.Name == fmt.Sprintf("SAI_OBJECT_TYPE_%s", meta.TypeName) {
				req.Option = fmt.Sprintf("option (sai_type) = OBJECT_TYPE_%s", meta.TypeName)
			}
		}
		if req.Option == "" {
			req.Option = "option (sai_type) = OBJECT_TYPE_UNSPECIFIED"
		}

		attrs, err := CreateAttrs(requestIdx, meta.TypeName, docInfo, docInfo.Attrs[meta.TypeName].CreateFields)
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
		attrs, err := CreateAttrs(2, meta.TypeName, docInfo, docInfo.Attrs[meta.TypeName].SetFields)
		if err != nil {
			return err
		}
		req.Fields = append(req.Fields, attrs...)
	case "get_attribute":
		req.Fields = append(req.Fields, idField, protoTmplField{
			ProtoType: repeatedType + strcase.UpperCamelCase(meta.TypeName+" attr"),
			Index:     2,
			Name:      "attr_type",
		})
		resp.Fields = append(resp.Fields, protoTmplField{
			Index:     1,
			Name:      "attr",
			ProtoType: strcase.UpperCamelCase(meta.TypeName + "Attribute"),
		})

		// attrEnum is the special emun that describes the possible values can be set/get for the API.
		attrEnum := ProtoEnum{
			Name:   strcase.UpperCamelCase(meta.TypeName + "_ATTR"),
			Values: []ProtoEnumValues{{Index: 0, Name: meta.TypeName + "_ATTR_UNSPECIFIED"}},
		}

		// For the attributes, generate code for the type if needed.
		for i, attr := range docInfo.Attrs[meta.TypeName].ReadFields {
			attrEnum.Values = append(attrEnum.Values, ProtoEnumValues{
				Index: i + 1,
				Name:  strings.TrimPrefix(attr.EnumName, "SAI_"),
			})
			// Handle function pointers as streaming RPCs.
			if strings.Contains(attr.SaiType, "sai_pointer_t") {
				funcName := strings.Split(attr.SaiType, " ")[1]
				name := saiast.TrimSAIName(strings.TrimSuffix(funcName, "_fn"), true, false)
				req := &ProtoTmplMessage{
					Name: strcase.UpperCamelCase(name + "_request"),
				}
				resp, ok := funcToStreamResp[funcName]
				if !ok {
					// TODO: There are 2 function pointers that don't follow this pattern, support them.
					log.Warningf("skipping unknown func type %q\n", funcName)
					continue
				}
				apis[apiName].Messages = append(apis[apiName].Messages, *req, resp)
				apis[apiName].RPCs = append(apis[apiName].RPCs, ProtoRPC{
					RequestName:  req.Name,
					ResponseName: "stream " + resp.Name,
					Name:         strcase.UpperCamelCase(name),
				})
			}
		}
		apis[apiName].Enums = append(apis[apiName].Enums, attrEnum)
	case "remove":
		req.Fields = append(req.Fields, idField)
	case "create_bulk":
		req.Fields = append(req.Fields, protoTmplField{
			Name:      "reqs",
			ProtoType: repeatedType + strcase.UpperCamelCase("Create "+meta.TypeName+"Request"),
			Index:     1,
		})
		resp.Fields = append(resp.Fields, protoTmplField{
			Name:      "resps",
			ProtoType: repeatedType + strcase.UpperCamelCase("Create "+meta.TypeName+"Response"),
			Index:     1,
		})
	case "remove_bulk":
		req.Fields = append(req.Fields, protoTmplField{
			Name:      "reqs",
			ProtoType: repeatedType + strcase.UpperCamelCase("Remove "+meta.TypeName+"Request"),
			Index:     1,
		})
		resp.Fields = append(resp.Fields, protoTmplField{
			Name:      "resps",
			ProtoType: repeatedType + strcase.UpperCamelCase("Remove "+meta.TypeName+"Response"),
			Index:     1,
		})
	case "get_stats":
		req.Fields = append(req.Fields, idField, protoTmplField{
			Name:      "counter_ids",
			ProtoType: fmt.Sprintf("repeated %s", strcase.UpperCamelCase(meta.TypeName+" stat")),
			Index:     2,
		})
		resp.Fields = append(resp.Fields, protoTmplField{
			Name:      "values",
			ProtoType: "repeated uint64",
			Index:     1,
		})
	default:
		return nil
	}
	apis[apiName].Messages = append(apis[apiName].Messages, *req, *resp)
	apis[apiName].RPCs = append(apis[apiName].RPCs, ProtoRPC{
		RequestName:  req.Name,
		ResponseName: resp.Name,
		Name:         strcase.UpperCamelCase(meta.Name),
	})
	return nil
}

func CreateAttrs(startIdx int, typeName string, xmlInfo *docparser.SAIInfo, attrs []*docparser.AttrTypeName) ([]protoTmplField, error) {
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

type SAITypeInfo struct {
	Repeated   bool
	ProtoType  string
	MessageDef string
	Required   bool
}

var (
	SAITypeToProto = map[string]SAITypeInfo{
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
		"map_sai_object_id_t": {
			ProtoType: "map<uint64, uint64>",
			Required:  true,
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
	repeated AclActionType action_list = 2;
}`,
		},
		"sai_acl_field_data_t": {
			ProtoType: "AclFieldData",
			MessageDef: `message AclFieldData {
	bool enable = 1;
	oneof mask {
		uint64 mask_uint = 2;
		int64 mask_int = 3;
		bytes mask_mac = 4;
		bytes mask_ip = 5;
		Uint64List mask_list = 6;
		bytes mask_u8list = 15;
	};
	oneof data {
		bool data_bool = 7;
		uint64 data_uint = 8;
		int64 data_int = 9;
		bytes data_mac = 10;
		bytes data_ip = 11;
		Uint64List data_list = 12;
		AclIpType data_ip_type = 13;
		uint64 data_oid = 14;
		bytes data_u8list = 16;
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
		PacketAction packet_action = 9;
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
	funcToStreamResp = map[string]ProtoTmplMessage{
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
	"sai_acl_field_data_t": func(_ string, _ *docparser.SAIInfo) (string, bool) {
		return "AclFieldData", false
	},
	"map_sai_acl_field_data_t": func(_ string, _ *docparser.SAIInfo) (string, bool) {
		return "map<uint64, AclFieldData>", true
	},
	"sai_acl_action_data_t": func(_ string, _ *docparser.SAIInfo) (string, bool) {
		return "AclActionData", false
	},
	"sai_pointer_t": func(_ string, _ *docparser.SAIInfo) (string, bool) { return "-", false }, // Noop, these are special cases.
}

// saiTypeToProtoType returns the protobuf type string for a SAI type.
// example: sai_u8_list_t -> repeated uint32
func SaiTypeToProtoType(saiType string, xmlInfo *docparser.SAIInfo) (string, bool, error) {
	saiType = strings.TrimPrefix(saiType, "const ")

	if pt, ok := SAITypeToProto[saiType]; ok {
		if pt.Repeated {
			return "repeated " + pt.ProtoType, true, nil
		}
		return pt.ProtoType, pt.Required, nil
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

type ProtoEnum struct {
	Name   string
	Values []ProtoEnumValues
	Alias  bool
}

type ProtoEnumValues struct {
	Index int
	Name  string
}

type ProtoTmplMessage struct {
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

type ProtoRPC struct {
	RequestName  string
	ResponseName string
	Name         string
}

const repeatedType = "repeated "
