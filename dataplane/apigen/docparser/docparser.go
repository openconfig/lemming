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
	"strconv"
	"strings"
)

// SAIInfo contains all the info parsed from the doxygen.
type SAIInfo struct {
	// Attrs is a map from sai type (sai_port_t) to its attributes.
	Attrs map[string]*Attr
	// Enums is a map from enum name (sai_port_media_type_t) to the values of the enum.
	Enums map[string][]*Enum
}

type Enum struct {
	Name  string
	Value int
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
	Value      int
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

// ParseSAIXMLDir parses all the SAI Doxygen XML files in a directory.
func ParseSAIXMLDir(xmlPath string) (*SAIInfo, error) {
	i := &SAIInfo{
		Attrs: make(map[string]*Attr),
		Enums: make(map[string][]*Enum),
	}
	files, err := os.ReadDir(xmlPath)
	if err != nil {
		return nil, err
	}
	members := []MemberDef{}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), "xml") {
			continue
		}
		if strings.Contains(file.Name(), "saimetadata") {
			fmt.Println("skipping file ", file.Name())
			continue
		}
		file := filepath.Join(xmlPath, file.Name())
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		dox := &Doxygen{}
		if err := xml.Unmarshal(b, dox); err != nil {
			return nil, err
		}
		for _, section := range dox.CompoundDef.SectionDef {
			members = append(members, section.MemberDef...)
		}
	}
	if err := parseXMLEntries(members, i); err != nil {
		return nil, err
	}

	return i, nil
}

var typeNameExpr = regexp.MustCompile("sai_(.*)_attr.*")

// memberToAttrInfo converts the MemberDef into attrInfo extracting the enum values,
// their types, and if they createable, readable, and/or writable.
func memberToAttrInfo(enum MemberDef) (*Attr, error) {
	info := &Attr{}
	trimStr := strings.TrimSuffix(strings.TrimPrefix(enum.Name, "_"), "_t") + "_"

	for i, value := range enum.EnumValues {
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
		val, err := parseInitializer(i, enum.EnumValues)
		if err != nil {
			return nil, err
		}
		if !canCreate && !canRead && !canSet {
			continue
		}
		atn := &AttrTypeName{
			EnumName:   value.Name,
			MemberName: strings.TrimPrefix(strings.ToLower(value.Name), trimStr),
			SaiType:    saiType,
			Comment:    strings.TrimSpace(value.BriefDescription.Paragraph.InlineText),
			Value:      val,
		}
		if strings.HasSuffix(value.Name, "_MIN") {
			atn.SaiType = fmt.Sprintf("map_%s", atn.SaiType)
		}
		if strings.HasSuffix(value.Name, "_MAX") {
			continue
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
	return info, nil
}

var saiConsts = map[string]int{
	"SAI_ACL_USER_DEFINED_FIELD_ATTR_ID_RANGE": 0xff,
}

func parseInitializer(i int, vals []EnumValue) (val int, rerr error) {
	defer func() {
		if rerr == nil {
			saiConsts[vals[i].Name] = val
		}
	}()

	findIndex := func(val string, vals []EnumValue) int {
		for i, eV := range vals {
			if eV.Name == val {
				return i
			}
		}
		return -1
	}

	init := strings.TrimSpace(strings.TrimPrefix(vals[i].Initializer, "= "))
	if init == "" {
		if i == 0 {
			return 0, nil
		}
		prev, err := parseInitializer(i-1, vals)
		return prev + 1, err
	}
	if strings.Contains(init, "+") {
		splits := strings.Split(init, "+")
		if len(splits) != 2 {
			return 0, fmt.Errorf("invalid init format expected A + B: %q", init)
		}
		lhs := strings.TrimSpace(splits[0])
		rhs := strings.TrimSpace(splits[1])
		rv, err := strconv.ParseUint(rhs, 0, 64)
		if err != nil {
			saiConst, ok := saiConsts[rhs]
			if !ok {
				return 0, err
			}
			rv = uint64(saiConst)
		}
		lv, err := parseInitializer(findIndex(lhs, vals), vals)
		if err != nil {
			return 0, err
		}
		return lv + int(rv), nil
	}
	if strings.Contains(init, "<<") {
		splits := strings.Split(init, "<<")
		if len(splits) != 2 {
			return 0, fmt.Errorf("invalid init format expected A << B: %q", init)
		}
		lhs := strings.TrimSpace(splits[0])
		rhs := strings.TrimSpace(splits[1])
		lv, err := strconv.ParseUint(lhs, 0, 64)
		if err != nil {
			return 0, err
		}
		rv, err := strconv.ParseUint(rhs, 0, 64)
		if err != nil {
			return 0, err
		}
		return int(lv << rv), nil
	}

	if v, err := strconv.ParseUint(init, 0, 64); err == nil {
		return int(v), nil
	}
	if v, ok := saiConsts[init]; ok {
		return v, nil
	}
	if idx := findIndex(init, vals); idx != -1 {
		return parseInitializer(idx, vals)
	}
	return 0, fmt.Errorf("failed to parse intialiazer for enum %v", vals[i].Name)
}

func memberToEnumValueStrings(enum MemberDef) ([]*Enum, error) {
	res := []*Enum{}
	for i, value := range enum.EnumValues {
		val, err := parseInitializer(i, enum.EnumValues)
		if err != nil {
			return nil, err
		}
		res = append(res, &Enum{
			Name:  value.Name,
			Value: val,
		})
	}
	return res, nil
}

// parseXMLEntries parses a single XML and appends the values into xmlInfo.
func parseXMLEntries(members []MemberDef, xmlInfo *SAIInfo) error {
	retryEntry := []MemberDef{}

	processEntry := func(enum MemberDef) error {
		if enum.Kind != "enum" {
			return nil
		}
		if strings.Contains(enum.Name, "attr_t") {
			matches := typeNameExpr.FindStringSubmatch(enum.Name)
			if len(matches) != 2 {
				return fmt.Errorf("unexpected number of matches: got %v", matches)
			}
			info, err := memberToAttrInfo(enum)
			if err != nil {
				return err
			}
			xmlInfo.Attrs[strings.ToUpper(matches[1])] = info
		} else {
			enums, err := memberToEnumValueStrings(enum)
			if err != nil {
				return err
			}
			xmlInfo.Enums[strings.TrimPrefix(enum.Name, "_")] = enums
		}
		return nil
	}

	for _, enum := range members {
		err := processEntry(enum)
		if err != nil {
			retryEntry = append(retryEntry, enum)
		}
	}
	// Retry processing enum once, for enums that refer to other enums.
	for _, enum := range retryEntry {
		err := processEntry(enum)
		if err != nil {
			return err
		}
	}

	return nil
}
