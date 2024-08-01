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

package ccgen

import (
	"fmt"
	"slices"
	"strings"
	"text/template"
	"unicode"

	strcase "github.com/stoewer/go-strcase"

	"github.com/openconfig/lemming/dataplane/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/apigen/saiast"
)

// Generate generates the C++ code for the SAI library.
func Generate(doc *docparser.SAIInfo, sai *saiast.SAIAPI, protoOutDir, ccOutDir string) (map[string]string, error) {
	files := make(map[string]string)
	enums := enumFile{
		ProtoOutDir: protoOutDir,
		CCOutDir:    ccOutDir,
	}
	for _, iface := range sai.Ifaces {
		apiName := strings.TrimSuffix(strings.TrimPrefix(iface.Name, "sai_"), "_api_t")
		ccData := ccTemplateData{
			IncludeGuard: fmt.Sprintf("DATAPLANE_STANDALONE_SAI_%s_H_", strings.ToUpper(apiName)),
			Header:       fmt.Sprintf("%s.h", apiName),
			APIType:      iface.Name,
			APIName:      apiName,
			ProtoInclude: apiName + ".pb",
			ProtoOutDir:  protoOutDir,
			CCOutDir:     ccOutDir,
		}
		switch apiName {
		case "switch":
			ccData.Globals = append(ccData.Globals, "std::unique_ptr<PortStateReactor> port_state;")
		}
		for _, fn := range iface.Funcs {
			meta := sai.GetFuncMeta(fn)
			opFn, convertFn := createCCData(meta, apiName, sai, doc, fn)
			if opFn != nil {
				ccData.Funcs = append(ccData.Funcs, opFn)
			}
			if convertFn != nil {
				ccData.ConvertFuncs = append(ccData.ConvertFuncs, convertFn)
			}
		}
		var headerBuilder, implBuilder strings.Builder
		if err := headerTmpl.Execute(&headerBuilder, ccData); err != nil {
			return nil, err
		}
		if err := ccTmpl.Execute(&implBuilder, ccData); err != nil {
			return nil, err
		}
		files[apiName+".h"] = headerBuilder.String()
		files[apiName+".cc"] = implBuilder.String()
		enums.ProtoIncludes = append(enums.ProtoIncludes, apiName+".pb")
	}
	for name, enumDesc := range doc.Enums {
		e := enum{
			SAIName:         name,
			ProtoName:       saiast.TrimSAIName(name, true, false),
			UnspecifiedName: saiast.TrimSAIName(name, false, true) + "_UNSPECIFIED",
		}
		if len(enumDesc) == 0 {
			continue
		}
		seenVal := map[int]struct{}{}
		for _, enumElemDesc := range enumDesc {
			if _, ok := seenVal[enumElemDesc.Value]; ok {
				continue
			}
			seenVal[enumElemDesc.Value] = struct{}{}
			name := strings.TrimPrefix(enumElemDesc.Name, "SAI_")
			// If the SAI name conflicts with unspecified proto name, then add SAI prefix,
			// that way the proto enum value is always 1 greater than the c enum.
			if name == e.UnspecifiedName {
				name = strings.TrimSuffix(saiast.TrimSAIName(name, false, true), "_UNSPECIFIED") + "_SAI_UNSPECIFIED"
			}
			e.Elems = append(e.Elems, enumElem{
				SAIName:   enumElemDesc.Name,
				ProtoName: name,
			})
		}
		e.DefaultReturn = e.Elems[0].SAIName
		enums.Enums = append(enums.Enums, e)
	}
	for name, attr := range doc.Attrs {
		if _, ok := unsupportedEnum[name]; ok {
			continue
		}
		e := enum{
			SAIName:         "sai_" + strings.ToLower(name) + "_attr_t",
			ProtoName:       saiast.TrimSAIName(name+"_ATTR", true, false),
			UnspecifiedName: saiast.TrimSAIName(name+"_ATTR", false, true) + "_UNSPECIFIED",
		}
		for _, val := range attr.ReadFields {
			name := strings.TrimPrefix(val.EnumName, "SAI_")
			// If the SAI name conflicts with unspecified proto name, then add SAI prefix,
			// that way the proto enum value is always 1 greater than the c enum.
			if name == e.UnspecifiedName {
				name = strings.TrimSuffix(saiast.TrimSAIName(name+"_ATTR", false, true), "_UNSPECIFIED") + "_SAI_UNSPECIFIED"
			}
			e.Elems = append(e.Elems, enumElem{
				SAIName:   val.EnumName,
				ProtoName: name,
			})
		}
		e.DefaultReturn = e.Elems[0].SAIName
		enums.Enums = append(enums.Enums, e)
	}
	slices.SortFunc(enums.Enums, func(a, b enum) int { return strings.Compare(a.SAIName, b.SAIName) })
	var headerBuilder, implBuilder strings.Builder
	if err := enumHeaderTmpl.Execute(&headerBuilder, enums); err != nil {
		return nil, err
	}
	if err := enumCCTmpl.Execute(&implBuilder, enums); err != nil {
		return nil, err
	}
	files["enum.h"] = headerBuilder.String()
	files["enum.cc"] = implBuilder.String()

	return files, nil
}

func sanitizeProtoName(inName string) string {
	name := strings.ReplaceAll(inName, "inline", "inline_") // inline is C++ keyword
	if unicode.IsDigit(rune(name[0])) {
		name = fmt.Sprintf("_%s", name)
	}
	return name
}

var unsupportedEnum = map[string]struct{}{
	"FDB_FLUSH":     {},
	"HOSTIF_PACKET": {},
}

const (
	createOp  = "create"
	getAttrOp = "get_attribute"
	setAttrOp = "set_attribute"
)

// createCCData returns a two structs with the template data for the given function.
// The first is the implementation of the API: CreateFoo.
// The second is the a conversion func from attribute list to the proto message. covert_create_foo.
func createCCData(meta *saiast.FuncMetadata, apiName string, sai *saiast.SAIAPI, info *docparser.SAIInfo, fn *saiast.TypeDecl) (*templateFunc, *templateFunc) {
	if info.Attrs[meta.TypeName] == nil {
		fmt.Printf("no doc info for type: %v\n", meta.TypeName)
		return nil, nil
	}
	opFn := &templateFunc{
		ReturnType: sai.Funcs[fn.Typ].ReturnType,
		Name:       meta.Name,
		Operation:  meta.Operation,
		TypeName:   meta.TypeName,
		ReqType:    strcase.UpperCamelCase(meta.Name + "_request"),
		RespType:   strcase.UpperCamelCase(meta.Name + "_response"),
	}
	convertFn := &templateFunc{
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
		convertFromFunc: "convert_from_acl_field_data",
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

var (
	headerTmpl = template.Must(template.New("header").Parse(`// Copyright 2023 Google LLC
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

#ifndef {{ .IncludeGuard }}
#define {{ .IncludeGuard }}

extern "C" {
	#include "inc/sai.h"
	#include "experimental/saiextensions.h"
}

extern const {{ .APIType }} l_{{ .APIName }};

{{ range .Funcs }}
{{ .ReturnType }} l_{{ .Name }}({{ .Args }});
{{ end }}

#endif  // {{ .IncludeGuard }}
`))
	ccTmpl = template.Must(template.New("cc").Parse(`
{{- define "getattr" }}
{{ $parent := . }}
switch ({{ .Var }}) {
  {{ range .Attrs }}
  case {{ .EnumValue }}:
	{{- if .CopyConvertFunc }}
	{{- if .ConvertFunc }}
	{{ .CopyConvertFunc}}({{.Var}}, {{ .ConvertFunc }}({{ $parent.ProtoVar }}.{{ .ProtoFunc }}({{ .Args }}){{ .ConvertFuncArgs }}));
	{{- else }}
	{{ .CopyConvertFunc}}({{.Var}}, {{ $parent.ProtoVar }}.{{ .ProtoFunc }}({{ .Args }}){{ .CopyConvertFuncArgs }});
	{{ end -}}
	{{ else }}
	{{- if .ConvertFunc }}
	{{ if .Var }} {{ .Var }} = {{ end }} {{ .ConvertFunc }}({{ $parent.ProtoVar }}.{{ .ProtoFunc }}({{ .Args }}){{ .ConvertFuncArgs }});
	{{ else }}
	{{ if .Var }} {{ .Var }} = {{ end }}  {{ $parent.ProtoVar }}.{{ .ProtoFunc }}({{ .Args }});
	{{ end -}}
	{{ end -}}
	break;
 {{- end }}
}
{{ end }}
{{- define "setattr" }}
{{ $parent := . }}
switch ({{ .Var }}) {
  {{ range .Attrs }}
  case {{ .EnumValue }}:
	{{- if .CustomText }}
	{{ .CustomText }}
	{{- end }}
	{{- if .Var }}
	{{ .Var }} = {{ .ConvertFunc }}( {{.Args}} );
	{{- else }}
	{{ $parent.ProtoVar }}.{{ .ProtoFunc }}({{ .Args }});
	{{- end }}
	break;
 {{- end }}
}
{{ end -}}
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

#include "{{ .CCOutDir }}/{{ .Header }}"
#include "{{ .CCOutDir }}/common.h"
#include "{{ .CCOutDir }}/enum.h"
#include "{{ .ProtoOutDir }}/common.pb.h"
#include "{{ .ProtoOutDir }}/{{ .ProtoInclude }}.h"
#include <glog/logging.h>

const {{ .APIType }} l_{{ .APIName }} = {
{{- range .Funcs }}
	.{{ .Name }} = l_{{ .Name }},
{{- end }}
};

{{ range .Globals }}
	{{- . }}
{{ end }}

{{- range .ConvertFuncs }}
lemming::dataplane::sai::{{ .ReturnType }} {{ .Name }}({{ .Args }}) {
{{ if eq .Operation "create" }}
lemming::dataplane::sai::{{ .ReturnType }} msg;
{{ if .SwitchScoped }} msg.set_switch_(switch_id); {{ end }}
{{ if .EntryVar }} *msg.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
 for(uint32_t i = 0; i < attr_count; i++ ) {
	{{ .AttrConvertInsert }}
	{{ template "setattr" .AttrSwitch }}
}
{{- else if eq .Operation "get_attribute" }}
{{ .ReturnType }} msg;
for(uint32_t i = 0; i < attr_count; i++ ) {
	{{ template "getattr" .AttrSwitch }}
}
{{- end }}
return msg;
}
{{ end }}

{{- range .Funcs }}
{{ .ReturnType }} l_{{ .Name }}({{ .Args }}) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	{{- if .UseCommonAPI }}
	{{ if eq .Operation "create" }}
	lemming::dataplane::sai::{{ .ReqType }} req = {{.ConvertFunc}}({{.Vars}});
	lemming::dataplane::sai::{{ .RespType }} resp;
	grpc::ClientContext context;
	{{ if .SwitchScoped }} req.set_switch_(switch_id); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	grpc::Status status = {{ .Client }}->{{ .RPCMethod }}(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	{{ if .OidVar -}}
	if ({{.OidPointer }}) {
	{{ .OidVar }} = resp.oid(); 
  	}
	{{ end }}
	{{ else if eq .Operation "create_bulk" }}
	lemming::dataplane::sai::{{ .ReqType }} req;
	lemming::dataplane::sai::{{ .RespType }} resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) {
		{{- if .SwitchScoped }}
		auto r = {{.ConvertFunc}}(switch_id, attr_count[i],attr_list[i]);
		{{ else }} 
		auto r = {{.ConvertFunc}}(attr_count[i], attr_list[i]);
		{{ end -}}
		{{- if .EntryVar }}
		*r.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}[i]); 
		{{ end -}}
		*req.add_reqs() = r;
	}

	grpc::Status status = {{ .Client }}->{{ .RPCMethod }}(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	if (object_count != resp.resps().size()) {
		return SAI_STATUS_FAILURE;
	}
	for (uint32_t i = 0; i < object_count; i++) {
		{{ if .OidVar -}} object_id[i] = resp.resps(i).oid(); {{ end }}
		object_statuses[i] = SAI_STATUS_SUCCESS;
	}

	{{ else if eq .Operation "get_attribute" }}
	lemming::dataplane::sai::{{ .ReqType }} req;
	lemming::dataplane::sai::{{ .RespType }} resp;
	grpc::ClientContext context;
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	{{ if .OidVar -}} req.set_oid({{ .OidVar }}); {{ end }}

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_{{ .AttrType }}_to_proto(attr_list[i].id));
	}
	grpc::Status status = {{ .Client }}->{{ .RPCMethod }}(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		{{ .AttrConvertInsert }}
		{{ template "getattr" .AttrSwitch }}
	}
	{{ else if and (eq .Operation "set_attribute") (ne (len .AttrSwitch.Attrs) 0) }}
	lemming::dataplane::sai::{{ .ReqType }} req;
	lemming::dataplane::sai::{{ .RespType }} resp;
	grpc::ClientContext context;
	{{ if .OidVar -}} req.set_oid({{ .OidVar }}); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	{{ .AttrConvertInsert }}
	{{ template "setattr" .AttrSwitch }}
	grpc::Status status = {{ .Client }}->{{ .RPCMethod }}(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	{{ else if eq .Operation "remove" }}
	lemming::dataplane::sai::{{ .ReqType }} req;
	lemming::dataplane::sai::{{ .RespType }} resp;
	grpc::ClientContext context;
	{{ if .OidVar -}} req.set_oid({{ .OidVar }}); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	grpc::Status status = {{ .Client }}->{{ .RPCMethod }}(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	{{ else if eq .Operation "remove_bulk" }}
	lemming::dataplane::sai::{{ .ReqType }} req;
	lemming::dataplane::sai::{{ .RespType }} resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) {
		{{ if .OidVar -}} req.add_reqs()->set_oid(object_id[i]); {{ end }}
		{{ if .EntryVar }} *req.add_reqs()->mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}[i]); {{ end }}
	}

	grpc::Status status = {{ .Client }}->{{ .RPCMethod }}(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	if (object_count != resp.resps().size()) {
		return SAI_STATUS_FAILURE;
	}
	for (uint32_t i = 0; i < object_count; i++) {
		object_statuses[i] = SAI_STATUS_SUCCESS;
	}
	{{ else if eq .Operation "get_stats" }}
	lemming::dataplane::sai::{{ .ReqType }} req;
	lemming::dataplane::sai::{{ .RespType }} resp;
	grpc::ClientContext context;
	{{ if .OidVar -}} req.set_oid({{ .OidVar }}); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	for (uint32_t i = 0; i < number_of_counters; i++) {
		req.add_counter_ids(convert_{{ .AttrType }}_to_proto(counter_ids[i]));
	}
	grpc::Status status = {{ .Client }}->{{ .RPCMethod }}(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < number_of_counters && i < uint32_t(resp.values_size()); i++ ) {
		counters[i] = resp.values(i);
	}
	{{ end }}
	return SAI_STATUS_SUCCESS;
	{{- else }}
	return SAI_STATUS_NOT_IMPLEMENTED;
	{{- end }}
}
{{ end }}
`))
	enumHeaderTmpl = template.Must(template.New("enum").Parse(`
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

#ifndef DATAPLANE_STANDALONE_SAI_ENUM_H_
#define DATAPLANE_STANDALONE_SAI_ENUM_H_

#include "{{ .ProtoOutDir }}/common.pb.h"
{{$file := .}} {{ range .ProtoIncludes }}
#include "{{ $file.ProtoOutDir }}/{{ . }}.h"
{{ end }}

extern "C" {
	#include "inc/sai.h"
	#include "experimental/saiextensions.h"
}

{{ range .Enums }}
lemming::dataplane::sai::{{ .ProtoName }} convert_{{ .SAIName }}_to_proto(const sai_int32_t val);
{{ .SAIName }} convert_{{ .SAIName }}_to_sai(lemming::dataplane::sai::{{ .ProtoName }} val);
google::protobuf::RepeatedField<int> convert_list_{{ .SAIName }}_to_proto(const sai_s32_list_t &list);
void convert_list_{{ .SAIName }}_to_sai(int32_t *list, const google::protobuf::RepeatedField<int> &proto_list, uint32_t *count);
{{ end }}

#endif  // DATAPLANE_STANDALONE_SAI_ENUM_H_
`))
	enumCCTmpl = template.Must(template.New("enum").Parse(`
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

#include "{{ .CCOutDir }}/enum.h"

{{ range .Enums }}
{{$enum := .}}
lemming::dataplane::sai::{{ .ProtoName }} convert_{{ .SAIName }}_to_proto(const sai_int32_t val) {
    switch (val) {
    {{ range .Elems }}
        case {{ .SAIName}}: return lemming::dataplane::sai::{{ .ProtoName}};
    {{ end }}
        default: return lemming::dataplane::sai::{{ $enum.UnspecifiedName}};
    }
}
{{ .SAIName }} convert_{{ .SAIName }}_to_sai(lemming::dataplane::sai::{{ .ProtoName }} val) {
    switch (val) {
    {{ range .Elems }}
        case lemming::dataplane::sai::{{ .ProtoName}}: return {{ .SAIName}};
    {{ end }}
        default: return {{ $enum.DefaultReturn }};
    }
}

google::protobuf::RepeatedField<int> convert_list_{{ .SAIName }}_to_proto(const sai_s32_list_t &list) {
	google::protobuf::RepeatedField<int> proto_list;
	for (int i = 0; i < list.count; i++) {
		proto_list.Add(convert_{{ .SAIName }}_to_proto(list.list[i]));
	}
	return proto_list;
}
void convert_list_{{ .SAIName }}_to_sai(int32_t *list, const google::protobuf::RepeatedField<int> &proto_list, uint32_t *count) {
	for (int i = 0; i < proto_list.size(); i++) {
		list[i] = convert_{{ .SAIName }}_to_sai(static_cast<lemming::dataplane::sai::{{.ProtoName}}>(proto_list[i]));
	}
	*count = proto_list.size();
}

{{ end }}
`))
)

type enumFile struct {
	ProtoIncludes []string
	Enums         []enum
	ProtoOutDir   string
	CCOutDir      string
}

type enum struct {
	SAIName         string
	ProtoName       string
	UnspecifiedName string
	Elems           []enumElem
	DefaultReturn   string
}

type enumElem struct {
	SAIName   string
	ProtoName string
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

type templateFunc struct {
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

type ccTemplateData struct {
	IncludeGuard string
	Header       string
	ProtoInclude string
	APIType      string
	APIName      string
	Globals      []string
	Funcs        []*templateFunc
	ConvertFuncs []*templateFunc
	ProtoOutDir  string
	CCOutDir     string
}
