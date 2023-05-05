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
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	cc "modernc.org/cc/v4"
)

type saiApi struct {
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

func getFuncAndTypes(ast *cc.AST) *saiApi {
	sa := saiApi{
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
#include "dataplane/standalone/log/log.h"

const {{ .APIType }} l_{{ .APIName }} = {
{{- range .Funcs }}
	.{{ .Name }} = l_{{ .Name }},
{{- end }}
};

{{ range .Funcs }}
{{ .ReturnType }} l_{{ .Name }}({{ .Args }}) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}

{{ end }}
`))
)

type templateFunc struct {
	ReturnType string
	Name       string
	Args       string
}

type templateData struct {
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

	for _, iface := range sai.ifaces {
		nameTrimmed := strings.TrimSuffix(strings.TrimPrefix(iface.name, "sai_"), "_api_t")
		data := templateData{
			IncludeGuard: fmt.Sprintf("DATAPLANE_STANDALONE_SAI_%s_H_", strings.ToUpper(nameTrimmed)),
			Header:       fmt.Sprintf("%s.h", nameTrimmed),
			APIType:      iface.name,
			APIName:      nameTrimmed,
		}
		for _, fn := range iface.funcs {
			var params []string
			for _, param := range sai.funcs[fn.typ].params {
				name := param.name
				typ := param.typ
				params = append(params, fmt.Sprintf("%s %s", typ, name))
			}
			data.Funcs = append(data.Funcs, templateFunc{
				ReturnType: sai.funcs[fn.typ].returnType,
				Name:       strings.TrimSuffix(strings.TrimPrefix(fn.name, "sai_"), "_fn"),
				Args:       strings.Join(params, ", "),
			})
		}
		header, err := os.Create(filepath.Join(outDir, data.Header))
		if err != nil {
			return err
		}
		impl, err := os.Create(filepath.Join(outDir, fmt.Sprintf("%s.cc", nameTrimmed)))
		if err != nil {
			return err
		}
		if err := headerTmpl.Execute(header, data); err != nil {
			return err
		}
		if err := ccTmpl.Execute(impl, data); err != nil {
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
