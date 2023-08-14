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
	"runtime"
	"strings"

	"github.com/openconfig/lemming/dataplane/standalone/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/standalone/apigen/protogen"
	"github.com/openconfig/lemming/dataplane/standalone/apigen/saiast"

	cc "modernc.org/cc/v4"
)

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
	saiPath     = "bazel-lemming/external/com_github_opencomputeproject_sai"
	ccOutDir    = "dataplane/standalone/sai"
	protoOutDir = "dataplane/standalone/proto"
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
	sai := saiast.Get(ast)
	xmlInfo, err := docparser.ParseSAIXMLDir()
	if err != nil {
		return err
	}

	protos, err := protogen.Generate(xmlInfo, sai)
	if err != nil {
		return err
	}

	for _, iface := range sai.Ifaces {
		nameTrimmed := strings.TrimSuffix(strings.TrimPrefix(iface.Name, "sai_"), "_api_t")
		ccData := ccTemplateData{
			IncludeGuard: fmt.Sprintf("DATAPLANE_STANDALONE_SAI_%s_H_", strings.ToUpper(nameTrimmed)),
			Header:       fmt.Sprintf("%s.h", nameTrimmed),
			APIType:      iface.Name,
			APIName:      nameTrimmed,
		}
		for _, fn := range iface.Funcs {
			meta := saiast.GetFuncMeta(fn, sai)
			tf := createCCData(meta, sai, fn)
			ccData.Funcs = append(ccData.Funcs, *tf)
		}
		header, err := os.Create(filepath.Join(ccOutDir, ccData.Header))
		if err != nil {
			return err
		}
		impl, err := os.Create(filepath.Join(ccOutDir, fmt.Sprintf("%s.cc", nameTrimmed)))
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
	for file, content := range protos {
		if err := os.WriteFile(filepath.Join(protoOutDir, file), []byte(content), 0666); err != nil {
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
