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
	"slices"
	"strings"
	"text/template"

	"github.com/openconfig/lemming/dataplane/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/apigen/saiast"
	"github.com/openconfig/lemming/dataplane/apigen/typeinfo"
	"github.com/openconfig/lemming/internal/lemmingutil"
)

// Generates returns a map of files containing the generated code code.
func Generate(doc *docparser.SAIInfo, sai *saiast.SAIAPI, protoPackage, protoGoPackage, protoOutDir string) (map[string]string, error) {
	files := map[string]string{}
	common, err := generateCommonTypes(doc, protoPackage, protoGoPackage)
	if err != nil {
		return nil, err
	}
	files["common.proto"] = common
	d, err := typeinfo.Data(doc, sai, protoPackage, protoGoPackage, "", protoOutDir)
	if err != nil {
		return nil, err
	}
	for apiName, iface := range d.APIs {
		var builder strings.Builder
		if err := protoTmpl.Execute(&builder, iface); err != nil {
			return nil, err
		}
		files[apiName+".proto"] = builder.String()
	}
	return files, nil
}

func rangeInOrder[T any](m map[string]T, pred func(key string, val T) error) error {
	keys := lemmingutil.Mapkeys(m)
	slices.Sort(keys)
	for _, key := range keys {
		if err := pred(key, m[key]); err != nil {
			return err
		}
	}
	return nil
}

// protoCommonTmplData contains the formated information needed to render the protobuf template.
type protoCommonTmplData struct {
	Messages       []string
	Enums          []*typeinfo.ProtoEnum
	Lists          []*typeinfo.ProtoTmplMessage
	ProtoPackage   string
	ProtoGoPackage string
}

// generateCommonTypes returns all contents of the common proto.
// These all reside in the common.proto file to simplify handling imports.
func generateCommonTypes(docInfo *docparser.SAIInfo, protoPackage, protoGoPackage string) (string, error) {
	common := &protoCommonTmplData{
		ProtoPackage:   protoPackage,
		ProtoGoPackage: protoGoPackage,
	}

	// Generate the hand-crafted messages.
	rangeInOrder(typeinfo.SAITypeToProto, func(_ string, typeInfo typeinfo.SAITypeInfo) error {
		if typeInfo.MessageDef != "" {
			common.Messages = append(common.Messages, typeInfo.MessageDef)
		}
		return nil
	})

	seenEnums := map[string]bool{}
	// Generate non-attribute enums.
	rangeInOrder(docInfo.Enums, func(name string, vals []*docparser.Enum) error {
		protoName := saiast.TrimSAIName(name, true, false)
		unspecifiedName := saiast.TrimSAIName(name, false, true) + "_UNSPECIFIED"
		enum := &typeinfo.ProtoEnum{
			Name:   protoName,
			Values: []typeinfo.ProtoEnumValues{{Index: 0, Name: unspecifiedName}},
		}
		seenValues := map[int]struct{}{}
		for _, val := range vals {
			name := strings.TrimPrefix(val.Name, "SAI_")
			// If the SAI name conflicts with unspecified proto name, then add SAI prefix,
			// that way the proto enum value is always 1 greater than the c enum.
			if name == unspecifiedName {
				name = strings.TrimSuffix(saiast.TrimSAIName(name, false, true), "_UNSPECIFIED") + "_SAI_UNSPECIFIED"
			}

			if _, ok := seenValues[val.Value]; ok {
				enum.Alias = true
			}
			seenValues[val.Value] = struct{}{}
			enum.Values = append(enum.Values, typeinfo.ProtoEnumValues{
				Index: val.Value + 1,
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
		attrFields, err := typeinfo.CreateAttrs(1, n, docInfo, attr.ReadFields)
		if err != nil {
			return err
		}
		common.Lists = append(common.Lists, &typeinfo.ProtoTmplMessage{
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

var (
	protoTmpl = template.Must(template.New("proto").Parse(`
syntax = "proto3";

package {{ .ProtoPackage }};

import "{{ .ProtoOutDir }}/common.proto";

option go_package = "{{ .ProtoGoPackage }}";

{{ range .Enums }}
enum {{ .Name }} {
	{{- range .Values }}
	{{ .Name }} = {{ .Index }};
	{{- end}}
}
{{ end -}}

{{ range .Types }}
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
	{{- range .Funcs }}
	{{- if .ProtoRPCName }}
	rpc {{ .ProtoRPCName }} ({{ .ProtoRequestType }}) returns ({{ .ProtoResponseType }}) {}
	{{- end }}
	{{- end }}
}
`))
	protoCommonTmpl = template.Must(template.New("common").Parse(`
syntax = "proto3";

package {{ .ProtoPackage }};
	
import "google/protobuf/timestamp.proto";
import "google/protobuf/descriptor.proto";

option go_package = "{{ .ProtoGoPackage }}";

extend google.protobuf.FieldOptions {
	optional int32 attr_enum_value = 515153358;
}

extend google.protobuf.MessageOptions {
	optional ObjectType sai_type = 515146388;
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

message InitializeRequest {
}

message InitializeResponse {
}

message UninitializeRequest {
}

message UninitializeResponse {
}

service Entrypoint {
  rpc ObjectTypeQuery(ObjectTypeQueryRequest) returns (ObjectTypeQueryResponse) {}
  rpc Initialize(InitializeRequest) returns (InitializeResponse) {}
  rpc Uninitialize(UninitializeRequest) returns (UninitializeResponse) {}
}
{{ range .Enums }}
enum {{ .Name }} {
	{{- if .Alias }}
	option allow_alias = true;
	{{- end }}
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
