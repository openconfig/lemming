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
	"slices"
	"strings"
	"text/template"

	"github.com/openconfig/lemming/dataplane/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/apigen/saiast"
	"github.com/openconfig/lemming/dataplane/apigen/typeinfo"
)

var unsupportedEnum = map[string]struct{}{
	"FDB_FLUSH":     {},
	"HOSTIF_PACKET": {},
}

// GenClient generates the C++ code for the SAI library.
func GenClient(doc *docparser.SAIInfo, sai *saiast.SAIAPI, protoOutDir, ccOutDir string) (map[string]string, error) {
	files := make(map[string]string)
	enums := enumFile{
		ProtoOutDir: protoOutDir,
		CCOutDir:    ccOutDir,
	}
	d, err := typeinfo.Data(doc, sai, "", "", ccOutDir, protoOutDir, false)
	if err != nil {
		return nil, err
	}
	for apiName, iface := range d.APIs {
		var headerBuilder, implBuilder strings.Builder
		if err := clientHdrTmpl.Execute(&headerBuilder, iface); err != nil {
			return nil, err
		}
		if err := clientCCTmpl.Execute(&implBuilder, iface); err != nil {
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

var (
	clientHdrTmpl = template.Must(template.New("header").Parse(`// Copyright 2023 Google LLC
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

extern "C" {
#include "experimental/saiextensions.h"
}

extern const {{ .APIType }} l_{{ .APIName }};

{{ range .Funcs }}
{{- if .Name }}
{{ .ReturnType }} l_{{ .Name }}({{ .Args }});
{{- end }}
{{ end }}

#endif  // {{ .IncludeGuard }}
`))
	clientCCTmpl = template.Must(template.New("cc").Parse(`
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
{{- if .Name }}
	.{{ .Name }} = l_{{ .Name }},
{{- end }}
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
{{- if .Name }}
{{ .ReturnType }} l_{{ .Name }}({{ .Args }}) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	{{- if .UseCommonAPI }}
	{{ if eq .Operation "create" }}
	lemming::dataplane::sai::{{ .ProtoRequestType }} req = {{.ConvertFunc}}({{.Vars}});
	lemming::dataplane::sai::{{ .ProtoResponseType }} resp;
	grpc::ClientContext context;
	{{ if .SwitchScoped }} req.set_switch_(switch_id); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	grpc::Status status = {{ .Client }}->{{ .ProtoRPCName }}(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	{{ if .OidVar -}}
	if ({{.OidPointer }}) {
	{{ .OidVar }} = resp.oid(); 
  	}
	{{ end }}
	{{ else if eq .Operation "create_bulk" }}
	lemming::dataplane::sai::{{ .ProtoRequestType }} req;
	lemming::dataplane::sai::{{ .ProtoResponseType }} resp;
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

	grpc::Status status = {{ .Client }}->{{ .ProtoRPCName }}(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
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
	lemming::dataplane::sai::{{ .ProtoRequestType }} req;
	lemming::dataplane::sai::{{ .ProtoResponseType }} resp;
	grpc::ClientContext context;
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	{{ if .OidVar -}} req.set_oid({{ .OidVar }}); {{ end }}

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_{{ .AttrType }}_to_proto(attr_list[i].id));
	}
	grpc::Status status = {{ .Client }}->{{ .ProtoRPCName }}(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		{{ .AttrConvertInsert }}
		{{ template "getattr" .AttrSwitch }}
	}
	{{ else if and (eq .Operation "set_attribute") (ne (len .AttrSwitch.Attrs) 0) }}
	lemming::dataplane::sai::{{ .ProtoRequestType }} req;
	lemming::dataplane::sai::{{ .ProtoResponseType }} resp;
	grpc::ClientContext context;
	{{ if .OidVar -}} req.set_oid({{ .OidVar }}); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	{{ .AttrConvertInsert }}
	{{ template "setattr" .AttrSwitch }}
	grpc::Status status = {{ .Client }}->{{ .ProtoRPCName }}(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	{{ else if eq .Operation "remove" }}
	lemming::dataplane::sai::{{ .ProtoRequestType }} req;
	lemming::dataplane::sai::{{ .ProtoResponseType }} resp;
	grpc::ClientContext context;
	{{ if .OidVar -}} req.set_oid({{ .OidVar }}); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	grpc::Status status = {{ .Client }}->{{ .ProtoRPCName }}(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	{{ else if eq .Operation "remove_bulk" }}
	lemming::dataplane::sai::{{ .ProtoRequestType }} req;
	lemming::dataplane::sai::{{ .ProtoResponseType }} resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) {
		{{ if .OidVar -}} req.add_reqs()->set_oid(object_id[i]); {{ end }}
		{{ if .EntryVar }} *req.add_reqs()->mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}[i]); {{ end }}
	}

	grpc::Status status = {{ .Client }}->{{ .ProtoRPCName }}(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (object_count != resp.resps().size()) {
		return SAI_STATUS_FAILURE;
	}
	for (uint32_t i = 0; i < object_count; i++) {
		object_statuses[i] = SAI_STATUS_SUCCESS;
	}
	{{ else if eq .Operation "get_stats" }}
	lemming::dataplane::sai::{{ .ProtoRequestType }} req;
	lemming::dataplane::sai::{{ .ProtoResponseType }} resp;
	grpc::ClientContext context;
	{{ if .OidVar -}} req.set_oid({{ .OidVar }}); {{ end }}
	{{ if .EntryVar }} *req.mutable_entry() = {{ .EntryConversionFunc }}({{ .EntryVar }}); {{ end }}
	for (uint32_t i = 0; i < number_of_counters; i++) {
		req.add_counter_ids(convert_{{ .AttrType }}_to_proto(counter_ids[i]));
	}
	grpc::Status status = {{ .Client }}->{{ .ProtoRPCName }}(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
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
{{- end }}
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
}

extern "C" {
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
