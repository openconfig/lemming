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
	Paragraph Paragraph `xml:"para"`
}

// Paragraph is a generic paragraph.
type Paragraph struct {
	SimpleSect []SimpleSect `xml:"simplesect"`
}

// SimpleSect contains a description of an element.
type SimpleSect struct {
	Para string `xml:"para"`
}

const xmlPath = "dataplane/standalone/apigen/xml"

// parseSAIXMLDir parses all the SAI Doxygen XML files in a directory.
func parseSAIXMLDir() (*xmlInfo, error) {
	i := &xmlInfo{
		attrs: make(map[string]attrInfo),
		enums: make(map[string][]string),
	}
	files, err := os.ReadDir(xmlPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if err := parseXMLFile(filepath.Join(xmlPath, file.Name()), i); err != nil {
			return nil, err
		}
	}
	return i, nil
}

var typeNameExpr = regexp.MustCompile("sai_(.*)_attr.*")

// handleEnumAttr converts the MemberDef into attrInfo extracting the enum values,
// their types, and if they createable, readable, and/or writable.
func handleEnumAttr(enum MemberDef) attrInfo {
	info := attrInfo{}
	trimStr := strings.TrimSuffix(strings.TrimPrefix(enum.Name, "_"), "_t") + "_"
	for _, value := range enum.EnumValues {
		var canCreate, canRead, canSet bool
		var saiType string
		for _, details := range value.DetailedDescription.Paragraph.SimpleSect {
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
					canSet = true
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
			info.readFields = append(info.readFields, atn)
		}
		if canSet {
			info.setFields = append(info.setFields, atn)
		}
	}
	return info
}

func memberToEnumValueStrings(enum MemberDef) []string {
	res := []string{}
	for _, value := range enum.EnumValues {
		res = append(res, value.Name)
	}
	return res
}

// parseXMLFile parses a single XML and appends the values into xmlInfo.
func parseXMLFile(file string, xmlInfo *xmlInfo) error {
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
			xmlInfo.enums[strings.TrimPrefix(enum.Name, "_")] = memberToEnumValueStrings(enum)
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
	// attrs is a map from sai type (sai_port_t) to its attributes.
	attrs map[string]attrInfo
	// attrs is a map from enum name (sai_port_media_type_t) to the values of the enum.
	enums map[string][]string
}
