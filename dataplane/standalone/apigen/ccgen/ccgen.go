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

	"github.com/openconfig/lemming/dataplane/standalone/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/standalone/apigen/protogen"
	"github.com/openconfig/lemming/dataplane/standalone/apigen/saiast"
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
		for _, fn := range iface.Funcs {
			meta := sai.GetFuncMeta(fn)
			tf := createCCData(meta, apiName, sai, doc, fn)
			ccData.Funcs = append(ccData.Funcs, *tf)
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

// createCCData returns a struct with the template data for the given function.
func createCCData(meta *saiast.FuncMetadata, apiName string, sai *saiast.SAIAPI, info *docparser.SAIInfo, fn *saiast.TypeDecl) *templateFunc {
	tf := &templateFunc{
		ReturnType: sai.Funcs[fn.Typ].ReturnType,
		Name:       meta.Name,
		Operation:  meta.Operation,
		TypeName:   meta.TypeName,
		ReqType:    strcase.UpperCamelCase(meta.Name + "_request"),
		RespType:   strcase.UpperCamelCase(meta.Name + "_response"),
	}

	var paramDefs []string
	var paramVars []string
	for _, param := range sai.Funcs[fn.Typ].Params {
		paramDefs = append(paramDefs, fmt.Sprintf("%s %s", param.Typ, param.Name))
		name := strings.ReplaceAll(param.Name, "*", "")
		// Functions that operator on entries take some entry type instead of an object id as argument.
		// Generate a entry union with the pointer to entry instead.
		if strings.Contains(param.Typ, "entry") {
			tf.Entry = fmt.Sprintf("common_entry_t entry = {.%s = %s};", name, name)
			name = "entry"
		}
		paramVars = append(paramVars, name)
	}
	tf.Args = strings.Join(paramDefs, ", ")
	tf.Vars = strings.Join(paramVars, ", ")
	tf.Client = strcase.SnakeCase(apiName)
	if tf.Client == "switch" { // switch is C++ keyword.
		tf.Client = "switch_"
	}
	tf.RPCMethod = strcase.UpperCamelCase(meta.Name)
	tf.SwitchScoped = meta.IsSwitchScoped
	tf.AttrEnumType = strcase.UpperCamelCase(meta.TypeName + " attr")

	// If the func has entry, then we don't use ids, instead pass the entry to the proto.
	if meta.Entry == "" {
		tf.OidVar = sai.Funcs[fn.Typ].Params[0].Name
	} else {
		entryType := strings.TrimPrefix(sai.Funcs[fn.Typ].Params[0].Typ, "const ")
		if ua, ok := typeToUnionAccessor[entryType]; ok {
			tf.EntryConversionFunc = ua.convertFromFunc
			tf.EntryVar = sai.Funcs[fn.Typ].Params[0].Name
		}
	}

	switch tf.Operation {
	case "create":
		tf.AttrSwitch = &AttrSwitch{
			Var:      "attr_list[i].id",
			ProtoVar: "req",
		}
		for _, attr := range info.Attrs[meta.TypeName].CreateFields {
			name := sanitizeProtoName(attr.MemberName)
			pFunc, arg, err := protoFieldSetter(attr.SaiType, name, "attr_list[i].value", info)
			if err != nil {
				fmt.Println("skipping due to error: ", err)
				continue
			}
			tf.AttrSwitch.Attrs = append(tf.AttrSwitch.Attrs, &AttrSwitchSmt{
				EnumValue: attr.EnumName,
				ProtoFunc: pFunc,
				Args:      arg,
			})
		}
	case "get_attribute":
		tf.AttrSwitch = &AttrSwitch{
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
			tf.AttrSwitch.Attrs = append(tf.AttrSwitch.Attrs, smt)
		}
	case "set_attribute":
		tf.AttrSwitch = &AttrSwitch{
			Var:      "attr->id",
			ProtoVar: "req",
		}
		for _, attr := range info.Attrs[meta.TypeName].SetFields {
			name := sanitizeProtoName(attr.MemberName)
			pFunc, arg, err := protoFieldSetter(attr.SaiType, name, "attr->value", info)
			if err != nil {
				fmt.Println("skipping due to error: ", err)
				continue
			}
			tf.AttrSwitch.Attrs = append(tf.AttrSwitch.Attrs, &AttrSwitchSmt{
				EnumValue: attr.EnumName,
				ProtoFunc: pFunc,
				Args:      arg,
			})
		}
	}

	tf.UseCommonAPI = supportedOperation[tf.Operation]
	// Function or types that don't follow standard naming.
	if strings.Contains(tf.TypeName, "PORT_ALL") || strings.Contains(tf.TypeName, "ALL_NEIGHBOR") {
		tf.UseCommonAPI = false
	}
	return tf
}

const protoNS = "lemming::dataplane::sai::"

type accessorType int

const (
	scalar accessorType = iota
	fixedSizedArray
	variableSizedArray
	convertFunc
)

type unionAccessor struct {
	accessor        string
	pointerOf       bool
	aType           accessorType
	convertFromFunc string
	convertToFunc   string
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
}

func protoFieldSetter(saiType, protoField, varName string, info *docparser.SAIInfo) (string, string, error) {
	setFn := fmt.Sprintf("set_%s", protoField)
	if _, ok := info.Enums[saiType]; ok {
		pType, _, err := protogen.SaiTypeToProtoType(saiType, info)
		if err != nil {
			return "", "", err
		}
		return setFn, fmt.Sprintf("static_cast<%s%s>(%s.s32 + 1)", protoNS, pType, varName), nil
	}

	ua, ok := typeToUnionAccessor[saiType]
	if !ok {
		return "", "", fmt.Errorf("unknown sai type: %q", saiType)
	}
	switch ua.aType {
	case scalar:
		return setFn, fmt.Sprintf("%s.%s", varName, ua.accessor), nil
	case convertFunc:
		return setFn, fmt.Sprintf("%s(%s.%s)", ua.convertFromFunc, varName, ua.accessor), nil
	case fixedSizedArray:
		if ua.pointerOf {
			return setFn, fmt.Sprintf("&%s.%s, sizeof(%s.%s)", varName, ua.accessor, varName, ua.accessor), nil
		}
		return setFn, fmt.Sprintf("%s.%s, sizeof(%s.%s)", varName, ua.accessor, varName, ua.accessor), nil
	case variableSizedArray:
		setFn = fmt.Sprintf("mutable_%s()->Add", protoField)
		return setFn, fmt.Sprintf("%s.%s.list, %s.%s.list + %s.%s.count", varName, ua.accessor, varName, ua.accessor, varName, ua.accessor), nil
	}
	return "", "", fmt.Errorf("unknown accessor type %q", ua.aType)
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
		smt.ConvertFunc = ua.convertToFunc
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
		smt.CopyConvertFuncArgs = fmt.Sprintf(", %s.count", smt.Var)
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
	"remove_bulk":        true,
	"set_attribute_bulk": true,
	"get_attribute_bulk": true,
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
	{{ $parent.ProtoVar }}.{{ .ProtoFunc }}({{ .Args }});
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
#include "dataplane/standalone/sai/entry.h"
#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/{{ .ProtoInclude }}.h"

const {{ .APIType }} l_{{ .APIName }} = {
{{- range .Funcs }}
	.{{ .Name }} = l_{{ .Name }},
{{- end }}
};

{{ range .Funcs }}
{{ .ReturnType }} l_{{ .Name }}({{ .Args }}) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	{{- if .UseCommonAPI }}
	{{ if eq .Operation "create" }}
	lemming::dataplane::sai::{{ .ReqType }} req;
	lemming::dataplane::sai::{{ .RespType }} resp;
	grpc::ClientContext context;
	{{ if .SwitchScoped }} req.set_switch_(switch_id); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		{{ template "setattr" .AttrSwitch }}
	}
	grpc::Status status = {{ .Client }}->{{ .RPCMethod }}(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	{{ if .OidVar -}} {{ .OidVar }} = resp.oid(); {{ end }}

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
}

type ccTemplateData struct {
	IncludeGuard string
	Header       string
	ProtoInclude string
	APIType      string
	APIName      string
	Funcs        []templateFunc
}
