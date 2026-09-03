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
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/openconfig/lemming/dataplane/apigen/ccgen"
	"github.com/openconfig/lemming/dataplane/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/apigen/protogen"
	"github.com/openconfig/lemming/dataplane/apigen/saiast"

	cc "modernc.org/cc/v4"
)

func parse(headers []string, includePaths ...string) (*cc.AST, error) {
	cfg, err := cc.NewConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, err
	}
	for _, p := range includePaths {
		cfg.SysIncludePaths = append(cfg.SysIncludePaths, p)
	}

	sources := []cc.Source{
		{Name: "<predefined>", Value: cfg.Predefined},
		{Name: "<builtin>", Value: cc.Builtin},
		{Name: "stdbool.h", Value: "#define bool _Bool\n#define true 1\n#define false 0"},
		{Name: "stddef.h", Value: "typedef unsigned long size_t; typedef long ptrdiff_t; typedef int wchar_t;"},
		{Name: "stdint.h", Value: "typedef signed char int8_t; typedef short int int16_t; typedef int int32_t; typedef long int int64_t; typedef unsigned char uint8_t; typedef unsigned short int uint16_t; typedef unsigned int uint32_t; typedef unsigned long int uint64_t;"},
	}

	for _, hdr := range headers {
		sources = append(sources, cc.Source{Name: hdr})
	}
	ast, err := cc.Translate(cfg, sources)
	if err != nil {
		return nil, err
	}
	return ast, nil
}

var (
	saiPath        = flag.String("sai_path", "bazel-lemming/external/com_github_opencomputeproject_sai", "Path to SAI repo")
	clientOutDir   = flag.String("client_out_dir", "dataplane/standalone/sai", "Output directory for C++ client code")
	serverOutDir   = flag.String("server_out_dir", "dataplane/standalone/saiserver", "Output directory for C++ server code")
	protoOutDir    = flag.String("proto_out_dir", "dataplane/proto/sai", "Output dirrectory for proto code")
	protoPackage   = flag.String("proto_package", "lemming.dataplane.sai", "Package for generated proto code")
	protoGoPackage = flag.String("proto_go_package", "github.com/openconfig/lemming/dataplane/proto/sai", "Go package option in proto code")
	xmlPath        = flag.String("xml_path", "dataplane/apigen/xml", "Path to generate SAI doxygen XML files.")
)

func generate() error {
	saiHeaderFile, err := filepath.Abs(filepath.Join(*saiPath, "inc/sai.h"))
	if err != nil {
		return err
	}
	expHeaderFile, err := filepath.Abs(filepath.Join(*saiPath, "experimental/saiextensions.h"))
	if err != nil {
		return err
	}
	incDir, err := filepath.Abs(filepath.Join(*saiPath, "inc"))
	if err != nil {
		return err
	}
	experiDir, err := filepath.Abs(filepath.Join(*saiPath, "experimental"))
	if err != nil {
		return err
	}
	ast, err := parse([]string{saiHeaderFile, expHeaderFile}, incDir, experiDir)
	if err != nil {
		return err
	}
	sai := saiast.Get(ast)
	xmlInfo, err := docparser.ParseSAIXMLDir(*xmlPath)
	if err != nil {
		return err
	}

	protos, err := protogen.Generate(xmlInfo, sai, *protoPackage, *protoGoPackage, *protoOutDir)
	if err != nil {
		return err
	}
	clientFiles, err := ccgen.GenClient(xmlInfo, sai, *protoOutDir, *clientOutDir)
	if err != nil {
		return err
	}
	serverFiles, err := ccgen.GenServer(xmlInfo, sai, *protoOutDir, *serverOutDir)
	if err != nil {
		return err
	}

	for file, content := range protos {
		if err := os.WriteFile(filepath.Join(*protoOutDir, file), []byte(content), 0o600); err != nil {
			return err
		}
	}
	for file, content := range clientFiles {
		if err := os.WriteFile(filepath.Join(*clientOutDir, file), []byte(content), 0o600); err != nil {
			return err
		}
	}
	for file, content := range serverFiles {
		if err := os.WriteFile(filepath.Join(*serverOutDir, file), []byte(content), 0o600); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()
	if err := generate(); err != nil {
		log.Fatal(err)
	}
}
