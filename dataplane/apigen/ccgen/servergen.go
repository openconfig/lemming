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

package ccgen

import (
	"slices"
	"strings"
	"text/template"

	"github.com/openconfig/lemming/dataplane/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/apigen/saiast"
	"github.com/openconfig/lemming/dataplane/apigen/typeinfo"
)

func GenServer(doc *docparser.SAIInfo, sai *saiast.SAIAPI, protoOutDir, ccOutDir string) (map[string]string, error) {
	files := make(map[string]string)
	enums := enumFile{
		ProtoOutDir: protoOutDir,
		CCOutDir:    ccOutDir,
	}
	d, err := typeinfo.Data(doc, sai, "", "", ccOutDir, protoOutDir)
	if err != nil {
		return nil, err
	}
	for apiName, iface := range d.APIs {
		var headerBuilder, implBuilder strings.Builder
		if err := serverHdrTmpl.Execute(&headerBuilder, iface); err != nil {
			return nil, err
		}
		if err := serverCCTmpl.Execute(&implBuilder, iface); err != nil {
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
	serverHdrTmpl = template.Must(template.New("header").Parse(`// Copyright 2024 Google LLC
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

#include "{{ .ProtoOutDir }}/common.pb.h"
#include "{{ .ProtoOutDir }}/{{ .ProtoInclude }}.h"
#include "{{ .ProtoOutDir }}/{{ .GRPCInclude }}.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class {{ .ServiceName }} final : public lemming::dataplane::sai::{{ .ServiceName }}::Service {
    public:
        {{ range .Funcs }}
		{{- if and .ProtoRPCName (not .IsStreaming) }}
        grpc::Status {{ .ProtoRPCName }}(grpc::ServerContext* context, const lemming::dataplane::sai::{{ .ProtoRequestType }}* req, lemming::dataplane::sai::{{ .ProtoResponseType }}* resp);
        {{- end }}
		{{ end }}
		{{ .APIType }} *api;
};

#endif  // {{ .IncludeGuard }}
`))
	serverCCTmpl = template.Must(template.New("cc").Parse(`// Copyright 2024 Google LLC
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

{{ range .Funcs }}
{{- if and .ProtoRPCName (not .IsStreaming) }}
grpc::Status {{ $.ServiceName }}::{{ .ProtoRPCName }}(grpc::ServerContext* context, const lemming::dataplane::sai::{{ .ProtoRequestType }}* req, lemming::dataplane::sai::{{ .ProtoResponseType }}* resp) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	{{- if .UseCommonAPI }}
	{{ if eq .Operation "remove" }}
	grpc::ClientContext context;
	{{ if .OidVar -}} auto status = api->{{ .Name }}(req.get_oid()); {{ end }}
  	{{ if .EntryVar }} 
	{{ .Args }} entry = {{ .EntryConversionToFunc }}(req); {{ end }}
	auto status = api->{{.Name}}(entry);
	if(!status.ok()) {
	auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return grpc::Status::INTERNAL;
	}
	{{end}}
	{{- end}}
	return grpc::Status::OK;
}

{{- end }}
{{ end }}
`))
)
