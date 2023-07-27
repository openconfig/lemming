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
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Doxygen is the root of the generated xml struct.
type Doxygen struct {
	CompoundDef CompoundDef `xml:"compounddef"`
}

// CompoundDef contains a list sections in the xml.
type CompoundDef struct {
	Title      string       `xml:"title"`
	SectionDef []SectionDef `xml:"sectiondef"`
}

// SectionDef contains a list of members.
type SectionDef struct {
	MemberDef []MemberDef `xml:"memberdef"`
}

// MemberDef is the definition of a single type.
type MemberDef struct {
	Name       string      `xml:"name"`
	EnumValues []EnumValue `xml:"enumvalue"`
	Kind       string      `xml:"kind,attr"`
}

// EnumValue is a single values in a enum.
type EnumValue struct {
	Name                string              `xml:"name"`
	Initializer         string              `xml:"initializer"`
	DetailedDescription DetailedDescription `xml:"detaileddescription"`
}

// DetailedDescription contains extra information about an enum value.
type DetailedDescription struct {
	Para Para `xml:"para"`
}

// Para is a generic paragraph.
type Para struct {
	SinpleSect []SinpleSect `xml:"simplesect"`
}

// SinpleSect contains a description of an element.
type SinpleSect struct {
	Para string `xml:"para"`
}

const xmlPath = "dataplane/standalone/apigen/xml"

func parseXml() (*xmlInfo, error) {
	i := &xmlInfo{
		attrs: make(map[string]attrInfo),
		enums: make(map[string][]string),
	}
	files, err := os.ReadDir(xmlPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if err := parseXmlFile(filepath.Join(xmlPath, file.Name()), i); err != nil {
			return nil, err
		}
	}
	return i, nil
}

var typeNameExpr = regexp.MustCompile("sai_(.*)_attr.*")

func handleEnumAttr(enum MemberDef) attrInfo {
	info := attrInfo{}
	trimStr := strings.TrimSuffix(strings.TrimPrefix(enum.Name, "_"), "_t") + "_"
	for _, value := range enum.EnumValues {
		var canCreate, canRead, canSet bool
		var saiType string
		for _, details := range value.DetailedDescription.Para.SinpleSect {
			annotation := strings.TrimSpace(details.Para)
			switch {
			case strings.HasPrefix(annotation, "@@type"):
				saiType = strings.TrimSpace(strings.TrimPrefix(annotation, "@@type"))
			case strings.HasPrefix(annotation, "@@flags"):
				switch {
				case strings.Contains(annotation, "CREATE_ONLY"):
					canCreate = true
					canRead = true
				case strings.Contains(annotation, "CREATE_AND_SET"):
					canCreate = true
					canRead = true
					canCreate = true
				case strings.Contains(annotation, "READ_ONLY"):
					canRead = true
				}
			}
		}
		if !canCreate && !canRead && !canSet {
			continue
		}
		atn := attrTypeName{
			EnumName:   value.Name,
			MemberName: strings.TrimPrefix(strings.ToLower(value.Name), trimStr),
			SaiType:    saiType,
		}
		if canCreate {
			info.createFields = append(info.createFields, atn)
		}
		if canRead {
			info.setFields = append(info.setFields, atn)
		}
		if canSet {
			info.readFields = append(info.readFields, atn)
		}
	}
	return info
}

func handleEnum(enum MemberDef) []string {
	res := []string{}
	for _, value := range enum.EnumValues {
		res = append(res, value.Name)
	}
	return res
}

func parseXmlFile(file string, xmlInfo *xmlInfo) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	dox := &Doxygen{}
	if err := xml.Unmarshal(b, dox); err != nil {
		return err
	}

	for _, enum := range dox.CompoundDef.SectionDef[0].MemberDef {
		if enum.Kind != "enum" {
			fmt.Printf("skipping kind %s\n", enum.Kind)
			continue
		}

		if strings.Contains(enum.Name, "attr_t") {
			matches := typeNameExpr.FindStringSubmatch(enum.Name)
			if len(matches) != 2 {
				return fmt.Errorf("unexpected number of matches: got %v", matches)
			}
			info := handleEnumAttr(enum)
			xmlInfo.attrs[strings.ToUpper(matches[1])] = info
		} else {
			xmlInfo.enums[strings.TrimPrefix(enum.Name, "_")] = handleEnum(enum)
		}

	}
	return nil
}

// attrInfo holds values and types for an attribute enum.
type attrInfo struct {
	createFields []attrTypeName
	setFields    []attrTypeName
	readFields   []attrTypeName
}

// attrTypeName contains the type and name of the attribute.
type attrTypeName struct {
	MemberName string
	SaiType    string
	EnumName   string
}

// xmlInfo contains all the info parsed from the doxygen.
type xmlInfo struct {
	attrs map[string]attrInfo
	enums map[string][]string
}
