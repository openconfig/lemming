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

// Package saiast contains information parsed from AST of the SAI header.
package saiast

import (
	"fmt"
	"regexp"
	"strings"

	"modernc.org/cc/v4"

	"github.com/stoewer/go-strcase"
)

// SAIAPI contains the information retreived from the AST.
type SAIAPI struct {
	Ifaces []*SAIInterface
	Funcs  map[string]*SAIFunc
}

// SAIInterface contains the functions contained in the interface.
type SAIInterface struct {
	Name  string
	Funcs []*TypeDecl
}

// SAIFunc is a function definition.
type SAIFunc struct {
	Name       string
	ReturnType string
	Params     []TypeDecl
}

// TypeDecl stores the name and type of a declation.
type TypeDecl struct {
	Name string
	Typ  string
}

// Get returns the information parses from the AST.
func Get(ast *cc.AST) *SAIAPI {
	sa := SAIAPI{
		Funcs: map[string]*SAIFunc{},
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
				sa.Funcs[fn.Name] = fn
			}
		}
		// Possible struct type declaration.
		if decl.DeclarationSpecifiers != nil && decl.DeclarationSpecifiers.DeclarationSpecifiers != nil && decl.DeclarationSpecifiers.DeclarationSpecifiers.TypeSpecifier != nil {
			if si := handleIfaces(name, decl); si != nil {
				sa.Ifaces = append(sa.Ifaces, si)
			}
		}
	}
	return &sa
}

// TrimSAIName trims sai_ prefix and _t from the string
func TrimSAIName(name string, makeCamel, makeUpper bool) string {
	str := strings.TrimSuffix(strings.TrimPrefix(name, "sai_"), "_t")
	if makeCamel {
		str = strcase.UpperCamelCase(str)
	}
	if makeUpper {
		str = strings.ToUpper(str)
	}

	return str
}

var funcExpr = regexp.MustCompile(`^([a-z]*_)(\w*)_(attribute|stats_ext|stats)|([a-z]*)_(\w*)$`)

// FuncMetadata contains additional information derived from a func definition.
type FuncMetadata struct {
	// Name is the name of the function: create_acl_table
	Name string
	// Operation is the action to take: create
	Operation string
	// TypeName is the name of the type: acl_table
	TypeName string
	// Entry is the name of the entry using the instead of id.
	Entry string
	// IsSwitchScoped is whether the API take a switch id as a parameter.
	IsSwitchScoped bool
}

// GetFuncMeta returns the metadata for a SAI func.
func (sai *SAIAPI) GetFuncMeta(fn *TypeDecl) *FuncMetadata {
	meta := &FuncMetadata{}
	meta.Name = strings.TrimSuffix(strings.TrimPrefix(fn.Name, "sai_"), "_fn")
	matches := funcExpr.FindStringSubmatch(meta.Name)
	meta.Operation = matches[1] + matches[4] + matches[3]
	meta.TypeName = strings.ToUpper(matches[2]) + strings.ToUpper(matches[5])

	for i, param := range sai.Funcs[fn.Typ].Params {
		if strings.Contains(param.Typ, "entry") {
			meta.Entry = TrimSAIName(strings.TrimPrefix(param.Typ, "const "), true, false)
		}
		// The switch_id is either the first or second arguments to function.
		if (i == 0 || i == 1) && param.Name == "switch_id" {
			meta.IsSwitchScoped = true
		}
	}
	// Handle plural types using the bulk API.
	if strings.HasSuffix(meta.TypeName, "PORTS") || strings.HasSuffix(meta.TypeName, "ENTRIES") || strings.HasSuffix(meta.TypeName, "MEMBERS") ||
		strings.HasSuffix(meta.TypeName, "LISTS") || strings.HasSuffix(meta.TypeName, "GROUPS") || strings.HasSuffix(meta.TypeName, "HOPS") ||
		strings.HasSuffix(meta.TypeName, "TUNNELS") || strings.HasSuffix(meta.TypeName, "INTERFACES") {
		meta.Operation += "_bulk"
		meta.TypeName = strings.TrimSuffix(meta.TypeName, "S")
		if strings.HasSuffix(meta.TypeName, "IE") {
			meta.TypeName = strings.TrimSuffix(meta.TypeName, "IE")
			meta.TypeName += "Y"
		}
	}

	return meta
}

func handleFunc(name string, decl *cc.Declaration, directDecl *cc.DirectDeclarator) *SAIFunc {
	sf := &SAIFunc{
		Name:       name,
		ReturnType: decl.DeclarationSpecifiers.Type().String(),
	}
	if directDecl.ParameterTypeList == nil {
		return nil
	}
	for paramList := directDecl.ParameterTypeList.ParameterList; paramList != nil; paramList = paramList.ParameterList {
		if paramList.ParameterDeclaration.Declarator == nil { // When the parameter is void.
			return nil
		}
		pd := TypeDecl{
			Name: paramList.ParameterDeclaration.Declarator.Name(),
			Typ:  paramList.ParameterDeclaration.Declarator.Type().String(),
		}
		if ptr, ok := paramList.ParameterDeclaration.Declarator.Type().(*cc.PointerType); ok {
			pd.Typ = ptr.Elem().String()
			pd.Name = fmt.Sprintf("*%s", pd.Name)
			if ptr2, ok := ptr.Elem().(*cc.PointerType); ok {
				pd.Typ = ptr2.Elem().String()
				pd.Name = fmt.Sprintf("*%s", pd.Name)
			}
		}

		if paramList.ParameterDeclaration.DeclarationSpecifiers.TypeQualifier != nil && paramList.ParameterDeclaration.DeclarationSpecifiers.TypeQualifier.Case == cc.TypeQualifierConst {
			pd.Typ = fmt.Sprintf("const %s", pd.Typ)
		}
		sf.Params = append(sf.Params, pd)
	}
	return sf
}

func handleIfaces(name string, decl *cc.Declaration) *SAIInterface {
	ts := decl.DeclarationSpecifiers.DeclarationSpecifiers.TypeSpecifier
	if ts.StructOrUnionSpecifier == nil || ts.StructOrUnionSpecifier.StructDeclarationList == nil {
		return nil
	}
	if !strings.Contains(name, "api_t") {
		return nil
	}
	si := &SAIInterface{
		Name: name,
	}

	structSpec := ts.StructOrUnionSpecifier.StructDeclarationList
	for sd := structSpec; sd != nil; sd = sd.StructDeclarationList {
		si.Funcs = append(si.Funcs, &TypeDecl{
			Name: sd.StructDeclaration.StructDeclaratorList.StructDeclarator.Declarator.Name(),
			Typ:  sd.StructDeclaration.StructDeclaratorList.StructDeclarator.Declarator.Type().String(),
		})
	}
	return si
}
