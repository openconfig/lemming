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

// The apigen command generates C++ and protobuf code for the SAI API.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	strcase "github.com/stoewer/go-strcase"
	cc "modernc.org/cc/v4"
)

type saiAPI struct {
	ifaces []*saiInterface
	funcs  map[string]*saiFunc
}

type saiInterface struct {
	name  string
	funcs []typeDecl
}

type saiFunc struct {
	name       string
	returnType string
	params     []typeDecl
}

type typeDecl struct {
	name string
	typ  string
}

func handleFunc(name string, decl *cc.Declaration, directDecl *cc.DirectDeclarator) *saiFunc {
	sf := &saiFunc{
		name:       name,
		returnType: decl.DeclarationSpecifiers.Type().String(),
	}
	if directDecl.ParameterTypeList == nil {
		return nil
	}
	for paramList := directDecl.ParameterTypeList.ParameterList; paramList != nil; paramList = paramList.ParameterList {
		if paramList.ParameterDeclaration.Declarator == nil { // When the parameter is void.
			return nil
		}
		pd := typeDecl{
			name: paramList.ParameterDeclaration.Declarator.Name(),
			typ:  paramList.ParameterDeclaration.Declarator.Type().String(),
		}
		if ptr, ok := paramList.ParameterDeclaration.Declarator.Type().(*cc.PointerType); ok {
			pd.typ = ptr.Elem().String()
			pd.name = fmt.Sprintf("*%s", pd.name)
			if ptr2, ok := ptr.Elem().(*cc.PointerType); ok {
				pd.typ = ptr2.Elem().String()
				pd.name = fmt.Sprintf("*%s", pd.name)
			}
		}

		if paramList.ParameterDeclaration.DeclarationSpecifiers.TypeQualifier != nil && paramList.ParameterDeclaration.DeclarationSpecifiers.TypeQualifier.Case == cc.TypeQualifierConst {
			pd.typ = fmt.Sprintf("const %s", pd.typ)
		}
		sf.params = append(sf.params, pd)
	}
	return sf
}

func handleIfaces(name string, decl *cc.Declaration) *saiInterface {
	ts := decl.DeclarationSpecifiers.DeclarationSpecifiers.TypeSpecifier
	if ts.StructOrUnionSpecifier == nil || ts.StructOrUnionSpecifier.StructDeclarationList == nil {
		return nil
	}
	if !strings.Contains(name, "api_t") {
		return nil
	}
	si := &saiInterface{
		name: name,
	}

	structSpec := ts.StructOrUnionSpecifier.StructDeclarationList
	for sd := structSpec; sd != nil; sd = sd.StructDeclarationList {
		si.funcs = append(si.funcs, typeDecl{
			name: sd.StructDeclaration.StructDeclaratorList.StructDeclarator.Declarator.Name(),
			typ:  sd.StructDeclaration.StructDeclaratorList.StructDeclarator.Declarator.Type().String(),
		})
	}
	return si
}

func getFuncAndTypes(ast *cc.AST) *saiAPI {
	sa := saiAPI{
		funcs: map[string]*saiFunc{},
	}
	for unit := ast.TranslationUnit; unit != nil; unit = unit.TranslationUnit {
		decl := unit.ExternalDeclaration.Declaration
		if decl == nil {
			continue
		}
		if decl.InitDeclaratorList == nil {
			continue
		}

		name := decl.InitDeclaratorList.InitDeclarator.Declarator.Name()
		if !strings.Contains(name, "sai") {
			continue
		}

		directDecl := decl.InitDeclaratorList.InitDeclarator.Declarator.DirectDeclarator
		if directDecl != nil { // Possible func declaration.
			if fn := handleFunc(name, decl, directDecl); fn != nil {
				sa.funcs[fn.name] = fn
			}
		}
		// Possible struct type declaration.
		if decl.DeclarationSpecifiers != nil && decl.DeclarationSpecifiers.DeclarationSpecifiers != nil && decl.DeclarationSpecifiers.DeclarationSpecifiers.TypeSpecifier != nil {
			if si := handleIfaces(name, decl); si != nil {
				sa.ifaces = append(sa.ifaces, si)
			}
		}
	}
	return &sa
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
	ccTmpl = template.Must(template.New("cc").Parse(`// Copyright 2023 Google LLC
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

const {{ .APIType }} l_{{ .APIName }} = {
{{- range .Funcs }}
	.{{ .Name }} = l_{{ .Name }},
{{- end }}
};

{{ range .Funcs }}
{{ .ReturnType }} l_{{ .Name }}({{ .Args }}) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	{{- if .UseCommonAPI }}
	{{- if .Entry }} {{ .Entry }} {{ end }}
	return translator->{{ .Operation }}(SAI_OBJECT_TYPE_{{ .TypeName }}, {{ .Vars }});
	{{- else }}
	return SAI_STATUS_NOT_IMPLEMENTED;
	{{- end }}
}
{{ end }}
`))
)

type templateFunc struct {
	ReturnType   string
	Name         string
	Args         string
	TypeName     string
	Operation    string
	Vars         string
	UseCommonAPI bool
	Entry        string
}

type ccTemplateData struct {
	IncludeGuard string
	Header       string
	APIType      string
	APIName      string
	Funcs        []templateFunc
}

func parse(header string, includePaths ...string) (*cc.AST, error) {
	cfg, err := cc.NewConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, err
	}
	for _, p := range includePaths {
		cfg.SysIncludePaths = append(cfg.SysIncludePaths, p)
	}

	sources := []cc.Source{{Name: "<predefined>", Value: cfg.Predefined}, {Name: "<builtin>", Value: cc.Builtin}, {Name: header}}
	ast, err := cc.Translate(cfg, sources)
	if err != nil {
		return nil, err
	}
	return ast, nil
}

const (
	saiPath = "bazel-lemming/external/com_github_opencomputeproject_sai"
	outDir  = "dataplane/standalone/sai"
)

var (
	supportedOperation = map[string]bool{
		"create":        true,
		"remove":        true,
		"get_attribute": true,
		"set_attribute": true,
		"clear_stats":   true,
		"get_stats":     true,
		"get_stats_ext": true,
	}
	funcExpr = regexp.MustCompile(`^([a-z]*_)(\w*)_(attribute|stats_ext|stats)|([a-z]*)_(\w*)$`)
)

// TODO: Enable generation.
var _ = template.Must(template.New("cc").Parse(`
syntax = "proto3";

package lemming.dataplane.sai;

option go_package = "github.com/openconfig/lemming/proto/dataplane/sai";

{{ range .Messages }}
message {{ .RequestName }} {
	{{ .RequestFieldsWrapperStart -}}
	{{- range .RequestFields }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }};
	{{- end }}
	{{ .RequestFieldsWrapperEnd }}
}

message {{ .ResponseName }} {
	{{ .ResponseFieldsWrapperStart -}}
	{{- range .ResponseFields }}
	{{ .ProtoType }} {{ .Name }} = {{ .Index }};
	{{- end }}
	{{ .ResponseFieldsWrapperEnd }}
}

{{ end }}


service {{ .ServiceName }} {
	{{- range .RPCs }}
	rpc {{ .Name }} ({{ .RequestName }}) returns ({{ .ResponseName }}) {}
	{{- end }}
}
`))

type protoTmplData struct {
	Messages    []protoTmplMessage
	RPCs        []protoRPC
	ServiceName string
}

type protoTmplMessage struct {
	RequestName                string
	ResponseName               string
	RequestFieldsWrapperStart  string
	RequestFieldsWrapperEnd    string
	RequestFields              []protoTmplField
	ResponseFieldsWrapperStart string
	ResponseFieldsWrapperEnd   string
	ResponseFields             []protoTmplField
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
	Repeated  bool
	ProtoType string
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
		ProtoType: "sai_port_err_status_t",
	},
	"sai_vlan_list_t": {
		Repeated:  true,
		ProtoType: "uint32",
	},
	"sai_u32_range_t": {
		ProtoType: "Uint32Range",
	},
	"sai_ip_address_t": {
		ProtoType: "IpAddress",
	},
	"sai_map_list_t": {
		Repeated:  true,
		ProtoType: "Uint32Range",
	},
	"sai_tlv_list_t": {
		Repeated:  true,
		ProtoType: "TLV",
	},
	"sai_qos_map_list_t": {
		Repeated:  true,
		ProtoType: "QOSMap",
	},
	"sai_system_port_config_t": {
		Repeated:  true,
		ProtoType: "SystemPortConfig",
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
	},
	"sai_prbs_rx_state_t": {
		ProtoType: "PRBS_RXState",
	},
	"sai_fabric_port_reachability_t": {
		ProtoType: "FabricPortReachability",
	},
	"sai_acl_resource_list_t": {
		Repeated:  true,
		ProtoType: "ACLResource",
	},
	"sai_acl_capability_t": {
		Repeated:  true,
		ProtoType: "ACLCapability",
	},
}

// saiTypeToProtoTypeCompound handles compound sai types (eg list of enums).
// The map key contains the base type (eg list) and func accepts the subtype (eg an enum type)
// and returns the full type string (eg repeated sample_enum).
var saiTypeToProtoTypeCompound = map[string]func(subType string, xmlInfo *xmlInfo) string{
	"sai_s32_list_t": func(subType string, xmlInfo *xmlInfo) string {
		if _, ok := xmlInfo.enums[subType]; !ok {
			return ""
		}
		return "repeated " + subType
	},
	// TODO: Support these types
	"sai_acl_field_data_t":  func(next string, xmlInfo *xmlInfo) string { return "-" },
	"sai_acl_action_data_t": func(next string, xmlInfo *xmlInfo) string { return "-" },
	"sai_pointer_t":         func(next string, xmlInfo *xmlInfo) string { return "-" },
}

// saiTypeToProtoType returns the protobuf type string for a SAI type.
// example: sai_u8_list_t -> repeated uint32
func saiTypeToProtoType(saiType string, xmlInfo *xmlInfo) (string, error) {
	typ := ""

	if protoType, ok := saiTypeToProto[saiType]; ok {
		typ = protoType.ProtoType
		if protoType.Repeated {
			typ = "repeated " + typ
		}
	} else if _, ok := xmlInfo.enums[saiType]; ok {
		typ = saiType
	} else if splits := strings.Split(saiType, " "); len(splits) == 2 {
		if fn, ok := saiTypeToProtoTypeCompound[splits[0]]; ok {
			typ = fn(splits[1], xmlInfo)
		}
	} else {
		return "", fmt.Errorf("unknown sai type: %v", saiType)
	}
	return typ, nil
}

func generate() error {
	headerFile, err := filepath.Abs(filepath.Join(saiPath, "inc/sai.h"))
	if err != nil {
		return err
	}
	incDir, err := filepath.Abs(filepath.Join(saiPath, "inc"))
	if err != nil {
		return err
	}
	experiDir, err := filepath.Abs(filepath.Join(saiPath, "experimental"))
	if err != nil {
		return err
	}
	ast, err := parse(headerFile, incDir, experiDir)
	if err != nil {
		return err
	}
	sai := getFuncAndTypes(ast)
	xmlInfo, err := parseXML()
	if err != nil {
		return err
	}

	for _, iface := range sai.ifaces {
		nameTrimmed := strings.TrimSuffix(strings.TrimPrefix(iface.name, "sai_"), "_api_t")
		ccData := ccTemplateData{
			IncludeGuard: fmt.Sprintf("DATAPLANE_STANDALONE_SAI_%s_H_", strings.ToUpper(nameTrimmed)),
			Header:       fmt.Sprintf("%s.h", nameTrimmed),
			APIType:      iface.name,
			APIName:      nameTrimmed,
		}
		protoData := protoTmplData{
			ServiceName: nameTrimmed, // TODO: prettier name
		}
		for _, fn := range iface.funcs {
			name := strings.TrimSuffix(strings.TrimPrefix(fn.name, "sai_"), "_fn")
			tf := templateFunc{
				ReturnType: sai.funcs[fn.typ].returnType,
				Name:       name,
			}

			var paramDefs []string
			var paramVars []string
			for _, param := range sai.funcs[fn.typ].params {
				paramDefs = append(paramDefs, fmt.Sprintf("%s %s", param.typ, param.name))
				name := strings.ReplaceAll(param.name, "*", "")
				// Functions that operator on entries take some entry type instead of an object id as argument.
				// Generate a entry union with the pointer to entry instead.
				if strings.Contains(param.typ, "entry") {
					tf.Entry = fmt.Sprintf("common_entry_t entry = {.%s = %s};", name, name)
					name = "entry"
				}
				paramVars = append(paramVars, name)
			}
			tf.Args = strings.Join(paramDefs, ", ")
			tf.Vars = strings.Join(paramVars, ", ")

			matches := funcExpr.FindStringSubmatch(name)
			tf.Operation = matches[1] + matches[4] + matches[3]

			tf.UseCommonAPI = supportedOperation[tf.Operation]
			tf.TypeName = strings.ToUpper(matches[2]) + strings.ToUpper(matches[5])

			// Handle plural types using the bulk API.
			if strings.HasSuffix(tf.TypeName, "PORTS") || strings.HasSuffix(tf.TypeName, "ENTRIES") || strings.HasSuffix(tf.TypeName, "MEMBERS") || strings.HasSuffix(tf.TypeName, "LISTS") {
				tf.Operation += "_bulk"
				tf.TypeName = strings.TrimSuffix(tf.TypeName, "S")
				if strings.HasSuffix(tf.TypeName, "IE") {
					tf.TypeName = strings.TrimSuffix(tf.TypeName, "IE")
					tf.TypeName += "Y"
				}
			}

			// Function or types that don't follow standard naming.
			if strings.Contains(tf.TypeName, "PORT_ALL") || strings.Contains(tf.TypeName, "ALL_NEIGHBOR") {
				tf.UseCommonAPI = false
			}
			ccData.Funcs = append(ccData.Funcs, tf)

			msg := protoTmplMessage{
				RequestName:  strcase.UpperCamelCase(tf.Name + "_request"),
				ResponseName: strcase.UpperCamelCase(tf.Name + "_response"),
			}

			// Handle proto generation
			// TODO: Enable proto generation and handle other funcs.
			switch tf.Operation {
			case "create":
				for i, attr := range xmlInfo.attrs[tf.TypeName].createFields {
					field := protoTmplField{
						Index: i + 1,
						Name:  attr.MemberName,
					}
					typ, err := saiTypeToProtoType(attr.SaiType, xmlInfo)
					if err != nil {
						return err
					}
					field.ProtoType = typ
					msg.RequestFields = append(msg.RequestFields, field)
				}
				protoData.Messages = append(protoData.Messages, msg)
				protoData.RPCs = append(protoData.RPCs, protoRPC{
					RequestName:  msg.RequestName,
					ResponseName: msg.ResponseName,
					Name:         strcase.UpperCamelCase(tf.Name),
				})
			case "set_attribute":
				for i, attr := range xmlInfo.attrs[tf.TypeName].setFields {
					field := protoTmplField{
						Index: i + 1,
						Name:  attr.MemberName,
					}
					msg.RequestFieldsWrapperStart = "oneof attr {"
					msg.RequestFieldsWrapperEnd = "}"
					typ, err := saiTypeToProtoType(attr.SaiType, xmlInfo)
					if err != nil {
						return fmt.Errorf("failed to get proto type for attr %s: %v", attr.MemberName, err)
					}
					field.ProtoType = typ
					msg.RequestFields = append(msg.RequestFields, field)
				}
				protoData.Messages = append(protoData.Messages, msg)
				protoData.RPCs = append(protoData.RPCs, protoRPC{
					RequestName:  msg.RequestName,
					ResponseName: msg.ResponseName,
					Name:         strcase.UpperCamelCase(tf.Name),
				})
			case "get_attribute":
				for i, attr := range xmlInfo.attrs[tf.TypeName].readFields {
					field := protoTmplField{
						Index: i + 1,
						Name:  attr.MemberName,
					}
					msg.ResponseFieldsWrapperStart = "oneof attr {"
					msg.ResponseFieldsWrapperEnd = "}"
					typ, err := saiTypeToProtoType(attr.SaiType, xmlInfo)
					if err != nil {
						return fmt.Errorf("failed to get proto type for attr %s: %v", attr.MemberName, err)
					}
					field.ProtoType = typ
					msg.ResponseFields = append(msg.ResponseFields, field)
				}
				protoData.Messages = append(protoData.Messages, msg)
				protoData.RPCs = append(protoData.RPCs, protoRPC{
					RequestName:  msg.RequestName,
					ResponseName: msg.ResponseName,
					Name:         strcase.UpperCamelCase(tf.Name),
				})
			}
		}
		header, err := os.Create(filepath.Join(outDir, ccData.Header))
		if err != nil {
			return err
		}
		impl, err := os.Create(filepath.Join(outDir, fmt.Sprintf("%s.cc", nameTrimmed)))
		if err != nil {
			return err
		}
		if err := headerTmpl.Execute(header, ccData); err != nil {
			return err
		}
		if err := ccTmpl.Execute(impl, ccData); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := generate(); err != nil {
		log.Fatal(err)
	}
}
