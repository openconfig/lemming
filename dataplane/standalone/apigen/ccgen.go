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
)

func createCCData(sai *saiAPI, fn typeDecl) (*templateFunc, bool, string) {
	name := strings.TrimSuffix(strings.TrimPrefix(fn.name, "sai_"), "_fn")
	tf := &templateFunc{
		ReturnType: sai.funcs[fn.typ].returnType,
		Name:       name,
	}

	isSwitchScoped := false
	entryType := ""
	var paramDefs []string
	var paramVars []string
	for i, param := range sai.funcs[fn.typ].params {
		paramDefs = append(paramDefs, fmt.Sprintf("%s %s", param.typ, param.name))
		name := strings.ReplaceAll(param.name, "*", "")
		// Functions that operator on entries take some entry type instead of an object id as argument.
		// Generate a entry union with the pointer to entry instead.
		if strings.Contains(param.typ, "entry") {
			tf.Entry = fmt.Sprintf("common_entry_t entry = {.%s = %s};", name, name)
			name = "entry"
			entryType = trimSAIName(strings.TrimPrefix(param.typ, "const"), true, false)
		}
		if i == 1 && param.name == "switch_id" {
			isSwitchScoped = true
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
	return tf, isSwitchScoped, entryType
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
