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

type GenType struct {
	Name   string
	Option string
	Fields []protoTmplField
}

type GenFunc struct {
	ProtoRPCName        string
	ProtoRequestType    string
	ProtoResponseType   string
	ReturnType          string
	Name                string
	Args                string
	TypeName            string
	Operation           string
	Vars                string
	UseCommonAPI        bool
	AttrSwitch          *AttrSwitch
	Client              string
	OidVar              string
	OidPointer          string
	AttrType            string
	AttrEnumType        string
	SwitchScoped        bool
	EntryConversionFunc string
	EntryVar            string
	EntryConversionToFunc string
	ConvertFunc         string
	AttrConvertInsert   string
	IsStreaming         bool
}

type APITemplate struct {
	Types          []*GenType
	Funcs          []*GenFunc
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
	ConvertFuncs   []*GenFunc
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
				GRPCInclude:    apiName + ".grpc.pb",
				CCOutDir:       ccOutDir,
			}
		}
		switch apiName {
		case "switch":
			data.APIs[apiName].Globals = append(data.APIs[apiName].Globals, "std::unique_ptr<PortStateReactor> port_state;")
		}

		for _, fn := range iface.Funcs {
			meta := sai.GetFuncMeta(fn)
			gFunc := &GenFunc{}
			protoReqType, protoRespType, err := genProtoReqResp(doc, apiName, meta)
			if err != nil {
				return nil, err
			}
			if protoReqType != nil && protoRespType != nil {
				gFunc.ProtoRPCName = strcase.UpperCamelCase(meta.Name)
				gFunc.ProtoRequestType = protoReqType.Name
				gFunc.ProtoResponseType = protoRespType.Name
				data.APIs[apiName].Types = append(data.APIs[apiName].Types, protoReqType, protoRespType)
			}

			populateCCInfo(meta, apiName, sai, doc, fn, gFunc)

			if gFunc.Operation == getAttrOp {
				enum := genProtoEnum(doc, apiName, meta)
				if enum != nil {
					data.APIs[apiName].Enums = append(data.APIs[apiName].Enums, *enum)
				}
				fns, msgs := genStreamingRPC(doc, apiName, meta)
				data.APIs[apiName].Funcs = append(data.APIs[apiName].Funcs, fns...)
				data.APIs[apiName].Types = append(data.APIs[apiName].Types, msgs...)
			}
			if gFunc.Operation == createOp {
				convertFn := genConvertFunc(gFunc, meta, doc, sai, fn)
				data.APIs[apiName].ConvertFuncs = append(data.APIs[apiName].ConvertFuncs, convertFn)
			}
			data.APIs[apiName].Funcs = append(data.APIs[apiName].Funcs, gFunc)
		}
	}
	return data, nil
}

func genConvertFunc(genFunc *GenFunc, meta *saiast.FuncMetadata, info *docparser.SAIInfo, sai *saiast.SAIAPI, fn *saiast.TypeDecl) *GenFunc {
	convertFn := &GenFunc{
		Name:              "convert_" + meta.Name,
		Operation:         genFunc.Operation,
		TypeName:          genFunc.TypeName,
		ProtoRequestType:  genFunc.ProtoRequestType,
		ProtoResponseType: genFunc.ProtoResponseType,
	}
	paramDefs, paramVars := getParamDefs(sai.Funcs[fn.Typ].Params)
	convertFn.Args = strings.Join(paramDefs[1:], ", ")
	convertFn.Vars = strings.Join(paramVars[1:], ", ")
	convertFn.AttrSwitch = &AttrSwitch{
		Var:      "attr_list[i].id",
		ProtoVar: "msg",
	}
	convertFn.ReturnType = genFunc.ProtoRequestType

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

	if meta.TypeName == "ACL_TABLE" {
		switch meta.Operation {
		case createOp:
			convertFn.AttrConvertInsert = `
if (attr_list[i].id >= SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
  (*msg.mutable_user_defined_field_group_min())[attr_list[i].id - SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN] = attr_list[i].value.oid;
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
		}
	}

	return convertFn
}

func getParamDefs(params []saiast.TypeDecl) ([]string, []string) {
	var paramDefs []string
	var paramVars []string
	for _, param := range params {
		paramDefs = append(paramDefs, fmt.Sprintf("%s %s", param.Typ, param.Name))
		name := strings.ReplaceAll(param.Name, "*", "")
		paramVars = append(paramVars, name)
	}
	return paramDefs, paramVars
}

// populateCCInfo returns a two structs with the template data for the given function.
// The first is the implementation of the API: CreateFoo.
// The second is the a conversion func from attribute list to the proto message. covert_create_foo.
func populateCCInfo(meta *saiast.FuncMetadata, apiName string, sai *saiast.SAIAPI, info *docparser.SAIInfo, fn *saiast.TypeDecl, genFunc *GenFunc) {
	if info.Attrs[meta.TypeName] == nil {
		fmt.Printf("no doc info for type: %v\n", meta.TypeName)
		return
	}
	genFunc.ReturnType = sai.Funcs[fn.Typ].ReturnType
	genFunc.Name = meta.Name
	genFunc.Operation = meta.Operation
	genFunc.TypeName = meta.TypeName

	paramDefs, paramVars := getParamDefs(sai.Funcs[fn.Typ].Params)
	genFunc.Args = strings.Join(paramDefs, ", ")
	genFunc.Vars = strings.Join(paramVars[1:], ", ")

	genFunc.Client = strcase.SnakeCase(apiName)
	if genFunc.Client == "switch" { // switch is C++ keyword.
		genFunc.Client = "switch_"
	}
	genFunc.SwitchScoped = meta.IsSwitchScoped
	genFunc.AttrEnumType = strcase.UpperCamelCase(meta.TypeName + " attr")
	genFunc.AttrType = strcase.SnakeCase("sai_" + meta.TypeName + "_attr_t")

	// If the func has entry, then we don't use ids, instead pass the entry to the proto.
	if meta.Entry == "" {
		genFunc.OidVar = sai.Funcs[fn.Typ].Params[0].Name
		genFunc.OidPointer = strings.TrimPrefix(genFunc.OidVar, "*")
	} else {
		i := 0
		if strings.Contains(genFunc.Operation, "bulk") {
			i = 1
		}
		entryType := strings.TrimPrefix(sai.Funcs[fn.Typ].Params[i].Typ, "const ")
		if ua, ok := typeToUnionAccessor[entryType]; ok {
			genFunc.EntryConversionFunc = ua.convertFromFunc
			genFunc.EntryConversionToFunc = ua.convertToFunc
			genFunc.EntryVar = sai.Funcs[fn.Typ].Params[i].Name
		}
	}

	switch genFunc.Operation {
	case createOp:
		genFunc.ConvertFunc = strcase.SnakeCase("convert_create " + meta.TypeName)
	case getAttrOp:
		genFunc.AttrSwitch = &AttrSwitch{
			Var:      "attr_list[i].id",
			ProtoVar: "resp.attr()",
		}
		for _, attr := range info.Attrs[meta.TypeName].ReadFields {
			name := sanitizeProtoName(attr.MemberName)
			smt, err := protoFieldGetter(attr.SaiType, name, "attr_list[i].value", info)
			if err != nil {
				fmt.Println("skipping due to error: ", err)
				continue
			}
			smt.EnumValue = attr.EnumName
			genFunc.AttrSwitch.Attrs = append(genFunc.AttrSwitch.Attrs, smt)
		}
	case setAttrOp:
		genFunc.AttrSwitch = &AttrSwitch{
			Var:      "attr->id",
			ProtoVar: "req",
		}
		for _, attr := range info.Attrs[meta.TypeName].SetFields {
			name := sanitizeProtoName(attr.MemberName)
			smt, err := protoFieldSetter(attr.SaiType, genFunc.AttrSwitch.ProtoVar, name, "attr->value", info)
			if err != nil {
				fmt.Println("skipping due to error: ", err)
				continue
			}
			smt.EnumValue = attr.EnumName
			genFunc.AttrSwitch.Attrs = append(genFunc.AttrSwitch.Attrs, smt)
		}
	case "create_bulk":
		genFunc.EntryVar = strings.TrimPrefix(genFunc.EntryVar, "*") // Usual entry is pointer, but for remove_bulk it's an array.
		genFunc.ConvertFunc = strcase.SnakeCase("convert_create " + meta.TypeName)
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
		genFunc.EntryVar = strings.TrimPrefix(genFunc.EntryVar, "*") // Usual entry is pointer, but for remove_bulk it's an array.
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
		genFunc.AttrType = strcase.SnakeCase("sai_" + meta.TypeName + "_stat_t")
		genFunc.AttrEnumType = strcase.UpperCamelCase(meta.TypeName + " stat")
	default:
	}

	// Patches for non-standard APIS
	if meta.TypeName == "ACL_TABLE" {
		switch meta.Operation {
		case getAttrOp:
			genFunc.AttrConvertInsert = `
if (attr_list[i].id >= SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
  attr_list[i].value.oid = resp.attr().user_defined_field_group_min().at(attr_list[i].id - SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN);
}`
		}
	}
	if meta.TypeName == "ACL_ENTRY" {
		switch meta.Operation {
		case setAttrOp:
			genFunc.AttrConvertInsert = `
if (attr->id >= SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr->id < SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
*(*req.mutable_user_defined_field_group_min())[attr->id  - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN].mutable_data_u8list() = std::string(attr->value.aclfield.data.u8list.list, attr->value.aclfield.data.u8list.list + attr->value.aclfield.data.u8list.count);
*(*req.mutable_user_defined_field_group_min())[attr->id  - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN].mutable_mask_u8list() = std::string(attr->value.aclfield.mask.u8list.list, attr->value.aclfield.mask.u8list.list + attr->value.aclfield.mask.u8list.count);
}`
		case getAttrOp:
			genFunc.AttrConvertInsert = `
if (attr_list[i].id >= SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
	auto acl_attr = resp.attr().user_defined_field_group_min().at(attr_list[i].id - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN);
	memcpy(attr_list[i].value.aclfield.data.u8list.list, acl_attr.data_u8list().data(), acl_attr.data_u8list().size());
	memcpy(attr_list[i].value.aclfield.mask.u8list.list, acl_attr.mask_u8list().data(), acl_attr.mask_u8list().size());
    attr_list[i].value.aclfield.data.u8list.count = acl_attr.data_u8list().size();
	attr_list[i].value.aclfield.mask.u8list.count = acl_attr.mask_u8list().size();
}`
		}
	}

	genFunc.UseCommonAPI = supportedOperation[genFunc.Operation]
	// Function or types that don't follow standard naming.
	if strings.Contains(genFunc.TypeName, "PORT_ALL") || strings.Contains(genFunc.TypeName, "ALL_NEIGHBOR") {
		genFunc.UseCommonAPI = false
	}
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

func genStreamingRPC(docInfo *docparser.SAIInfo, apiName string, meta *saiast.FuncMetadata) ([]*GenFunc, []*GenType) {
	types := []*GenType{}
	rpcs := []*GenFunc{}
	if docInfo.Attrs[meta.TypeName] == nil {
		return nil, nil
	}

	for _, attr := range docInfo.Attrs[meta.TypeName].ReadFields {
		if strings.Contains(attr.SaiType, "sai_pointer_t") {
			funcName := strings.Split(attr.SaiType, " ")[1]
			name := saiast.TrimSAIName(strings.TrimSuffix(funcName, "_fn"), true, false)
			req := &GenType{
				Name: strcase.UpperCamelCase(name + "_request"),
			}
			resp, ok := funcToStreamResp[funcName]
			if !ok {
				// TODO: There are 2 function pointers that don't follow this pattern, support them.
				log.Warningf("skipping unknown func type %q\n", funcName)
				continue
			}
			types = append(types, req, resp)
			rpcs = append(rpcs, &GenFunc{
				ProtoRequestType:  req.Name,
				ProtoResponseType: "stream " + resp.Name,
				ProtoRPCName:      strcase.UpperCamelCase(name),
				IsStreaming:       true,
			})
		}
	}
	return rpcs, types
}

func genProtoEnum(docInfo *docparser.SAIInfo, apiName string, meta *saiast.FuncMetadata) *ProtoEnum {
	// attrEnum is the special emun that describes the possible values can be set/get for the API.
	if docInfo.Attrs[meta.TypeName] == nil {
		return nil
	}

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
	}
	return &attrEnum
}

func genProtoReqResp(docInfo *docparser.SAIInfo, apiName string, meta *saiast.FuncMetadata) (*GenType, *GenType, error) {
	req := &GenType{
		Name: strcase.UpperCamelCase(meta.Name + "_request"),
	}
	resp := &GenType{
		Name: strcase.UpperCamelCase(meta.Name + "_response"),
	}
	if docInfo.Attrs[meta.TypeName] == nil {
		fmt.Printf("no doc info for type: %v\n", meta.TypeName)
		return nil, nil, nil
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
	case createOp:
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
			return nil, nil, err
		}
		req.Fields = append(req.Fields, attrs...)
		if meta.Entry == "" { // Entries don't have id.
			resp.Fields = append(resp.Fields, idField)
		}
	case setAttrOp:
		// If there are no settable attributes, do nothing.
		if len(docInfo.Attrs[meta.TypeName].SetFields) == 0 {
			return nil, nil, nil
		}
		req.Fields = append(req.Fields, idField)
		attrs, err := CreateAttrs(2, meta.TypeName, docInfo, docInfo.Attrs[meta.TypeName].SetFields)
		if err != nil {
			return nil, nil, err
		}
		req.Fields = append(req.Fields, attrs...)
	case getAttrOp:
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
		return nil, nil, nil
	}
	return req, resp, nil
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

const repeatedType = "repeated "
