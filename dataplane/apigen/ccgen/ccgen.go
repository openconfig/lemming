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
	"strings"
	"text/template"
	"unicode"

	strcase "github.com/stoewer/go-strcase"

	"github.com/openconfig/lemming/dataplane/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/apigen/protogen"
	"github.com/openconfig/lemming/dataplane/apigen/saiast"
)

// Generate generates the C++ code for the SAI library.
func Generate(doc *docparser.SAIInfo, sai *saiast.SAIAPI) (map[string]string, error) {
	files := make(map[string]string)
	for _, iface := range sai.Ifaces {
		apiName := strings.TrimSuffix(strings.TrimPrefix(iface.Name, "sai_"), "_api_t")
		ccData := ccTemplateData{
			IncludeGuard: fmt.Sprintf("DATAPLANE_STANDALONE_SAI_%s_H_", strings.ToUpper(apiName)),
			Header:       fmt.Sprintf("%s.h", apiName),
			APIType:      iface.Name,
			APIName:      apiName,
			ProtoInclude: apiName + ".pb",
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
	}
	return files, nil
}

func sanitizeProtoName(inName string) string {
	name := strings.ReplaceAll(inName, "inline", "inline_") // inline is C++ keyword
	if unicode.IsDigit(rune(name[0])) {
		name = fmt.Sprintf("_%s", name)
	}
	return name
}

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

	// If the func has entry, then we don't use ids, instead pass the entry to the proto.
	if meta.Entry == "" {
		opFn.OidVar = sai.Funcs[fn.Typ].Params[0].Name
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
	case "create":
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
	case "get_attribute":
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
	case "set_attribute":
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
	case "get_stats":
		convertFn = nil
		opFn.AttrEnumType = strcase.UpperCamelCase(meta.TypeName + " stat")
	default:
		convertFn = nil
	}

	opFn.UseCommonAPI = supportedOperation[opFn.Operation]
	// Function or types that don't follow standard naming.
	if strings.Contains(opFn.TypeName, "PORT_ALL") || strings.Contains(opFn.TypeName, "ALL_NEIGHBOR") {
		opFn.UseCommonAPI = false
	}
	return opFn, convertFn
}

const protoNS = "lemming::dataplane::sai::"

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
		aType:           acl,
	},
	"sai_acl_action_data_t sai_object_id_t": {
		accessor:        "oid",
		convertFromFunc: "convert_from_acl_action_data",
		aType:           acl,
	},
	"sai_acl_action_data_t sai_packet_action_t": {
		accessor:        "s32",
		convertFromFunc: "convert_from_acl_action_data_action",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_acl_ip_type_t": {
		accessor:        "s32",
		convertFromFunc: "convert_from_acl_field_data_ip_type",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_uint8_t": {
		accessor:        "u8",
		convertFromFunc: "convert_from_acl_field_data",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_uint16_t": {
		accessor:        "u16",
		convertFromFunc: "convert_from_acl_field_data",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_ip6_t": {
		accessor:        "ip6",
		convertFromFunc: "convert_from_acl_field_data_ip6",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_mac_t": {
		accessor:        "mac",
		convertFromFunc: "convert_from_acl_field_data_mac",
		aType:           acl,
	},
	"sai_acl_field_data_t sai_object_id_t": {
		accessor:        "oid",
		convertFromFunc: "convert_from_acl_field_data",
		aType:           acl,
	},
}

func protoFieldSetter(saiType, protoVar, protoField, varName string, info *docparser.SAIInfo) (*AttrSwitchSmt, error) {
	smt := &AttrSwitchSmt{
		ProtoFunc: fmt.Sprintf("set_%s", protoField),
	}

	if _, ok := info.Enums[saiType]; ok {
		pType, _, err := protogen.SaiTypeToProtoType(saiType, info)
		if err != nil {
			return nil, err
		}
		smt.Args = fmt.Sprintf("static_cast<%s%s>(%s.s32 + 1)", protoNS, pType, varName)
		return smt, nil
	}

	ua, ok := typeToUnionAccessor[saiType]
	if !ok {
		return nil, fmt.Errorf("unknown sai type: %q", saiType)
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
		smt.ConvertFunc = fmt.Sprintf("static_cast<int>")
		smt.ConvertFuncArgs = " - 1"
		return smt, nil
	}
	ua, ok := typeToUnionAccessor[saiType]
	if !ok {
		return nil, fmt.Errorf("unknown sai type: %q", saiType)
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
	}
	return nil, fmt.Errorf("unknown accessor type %q", saiType)
}

var supportedOperation = map[string]bool{
	"create":             true,
	"remove":             true,
	"get_attribute":      true,
	"set_attribute":      true,
	"clear_stats":        true,
	"get_stats":          true,
	"get_stats_ext":      true,
	"create_bulk":        true,
	"remove_bulk":        false,
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
{{ define "getattr" }}
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
{{ define "setattr" }}
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
{{ end }}


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

#include "dataplane/standalone/sai/{{ .Header }}"
#include <glog/logging.h>
#include "dataplane/standalone/sai/common.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/{{ .ProtoInclude }}.h"

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
	{{ if .OidVar -}} {{ .OidVar }} = resp.oid(); {{ end }}
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
		*r.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); 
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
		req.add_attr_type(static_cast<lemming::dataplane::sai::{{ .AttrEnumType }}>(attr_list[i].id + 1));
	}
	grpc::Status status = {{ .Client }}->{{ .RPCMethod }}(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		{{ template "getattr" .AttrSwitch }}
	}
	{{ else if and (eq .Operation "set_attribute") (ne (len .AttrSwitch.Attrs) 0) }}
	lemming::dataplane::sai::{{ .ReqType }} req;
	lemming::dataplane::sai::{{ .RespType }} resp;
	grpc::ClientContext context;
	{{ if .OidVar -}} req.set_oid({{ .OidVar }}); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
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
	{{ else if eq .Operation "get_stats" }}
	lemming::dataplane::sai::{{ .ReqType }} req;
	lemming::dataplane::sai::{{ .RespType }} resp;
	grpc::ClientContext context;
	{{ if .OidVar -}} req.set_oid({{ .OidVar }}); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	for (uint32_t i = 0; i < number_of_counters; i++) {
		req.add_counter_ids(static_cast<lemming::dataplane::sai::{{ .AttrEnumType }}>(counter_ids[i] + 1));
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
)

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
	AttrType            string
	AttrEnumType        string
	SwitchScoped        bool
	EntryConversionFunc string
	EntryVar            string
	ConvertFunc         string
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
}
