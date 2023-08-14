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

package main

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/openconfig/lemming/dataplane/standalone/apigen/saiast"
)

func createCCData(meta *saiast.FuncMetadata, sai *saiast.SAIAPI, fn *saiast.TypeDecl) *templateFunc {
	name := meta.Name
	tf := &templateFunc{
		ReturnType: sai.Funcs[fn.Typ].ReturnType,
		Name:       name,
		Operation:  meta.Operation,
		TypeName:   meta.TypeName,
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

	tf.UseCommonAPI = supportedOperation[tf.Operation]

	// Function or types that don't follow standard naming.
	if strings.Contains(tf.TypeName, "PORT_ALL") || strings.Contains(tf.TypeName, "ALL_NEIGHBOR") {
		tf.UseCommonAPI = false
	}
	return tf
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
