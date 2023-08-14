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

// Package docparser implements a parser for Doxygen xml.
package docparser

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Info contains all the info parsed from the doxygen.
type Info struct {
	// attrs is a map from sai type (sai_port_t) to its attributes.
	Attrs map[string]*Attr
	// attrs is a map from enum name (sai_port_media_type_t) to the values of the enum.
	Enums map[string][]string
}

// attrInfo holds values and types for an attribute enum.
type Attr struct {
	CreateFields []*AttrTypeName
	SetFields    []*AttrTypeName
	ReadFields   []*AttrTypeName
}

// attrTypeName contains the type and name of the attribute.
type AttrTypeName struct {
	MemberName string
	SaiType    string
	EnumName   string
	Comment    string
}

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
	Name                string      `xml:"name"`
	Initializer         string      `xml:"initializer"`
	DetailedDescription Description `xml:"detaileddescription"`
	BriefDescription    Description `xml:"briefdescription"`
}

// Description contains extra information about an enum value.
type Description struct {
	Paragraph Paragraph `xml:"para"`
}

// Paragraph is a generic paragraph.
type Paragraph struct {
	InlineText string       `xml:",chardata"` // For BriefDescription, the paragraph contains raw text.
	SimpleSect []SimpleSect `xml:"simplesect"`
}

// SimpleSect contains a description of an element.
type SimpleSect struct {
	Para string `xml:"para"`
}

const xmlPath = "dataplane/standalone/apigen/xml"

// ParseSAIXMLDir parses all the SAI Doxygen XML files in a directory.
func ParseSAIXMLDir() (*Info, error) {
	i := &Info{
		Attrs: make(map[string]*Attr),
		Enums: make(map[string][]string),
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

// memberToAttrInfo converts the MemberDef into attrInfo extracting the enum values,
// their types, and if they createable, readable, and/or writable.
func memberToAttrInfo(enum MemberDef) *Attr {
	info := &Attr{}
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
		atn := &AttrTypeName{
			EnumName:   value.Name,
			MemberName: strings.TrimPrefix(strings.ToLower(value.Name), trimStr),
			SaiType:    saiType,
			Comment:    strings.TrimSpace(value.BriefDescription.Paragraph.InlineText),
		}
		if canCreate {
			info.CreateFields = append(info.CreateFields, atn)
		}
		if canRead {
			info.ReadFields = append(info.ReadFields, atn)
		}
		if canSet {
			info.SetFields = append(info.SetFields, atn)
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
func parseXMLFile(file string, xmlInfo *Info) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	dox := &Doxygen{}
	if err := xml.Unmarshal(b, dox); err != nil {
		return err
	}

	for _, section := range dox.CompoundDef.SectionDef {
		for _, enum := range section.MemberDef {
			if enum.Kind != "enum" {
				fmt.Printf("skipping kind %s\n", enum.Kind)
				continue
			}
			if strings.Contains(enum.Name, "attr_t") {
				matches := typeNameExpr.FindStringSubmatch(enum.Name)
				if len(matches) != 2 {
					return fmt.Errorf("unexpected number of matches: got %v", matches)
				}
				info := memberToAttrInfo(enum)
				xmlInfo.Attrs[strings.ToUpper(matches[1])] = info
			} else {
				xmlInfo.Enums[strings.TrimPrefix(enum.Name, "_")] = memberToEnumValueStrings(enum)
			}
		}
	}

	return nil
}
