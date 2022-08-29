// Copyright 2022 Google LLC
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

// Package fwdattribute provides types and mechanisms used to manage attributes
// for forwarding objects. An attribute is a named property identified by a
// human readable name and a string value. The name and value of an attribute
// is opaque to the forwarding api. An attribute may be associated with a
// forwarding object, forwarding context or the global context.
package fwdattribute

import (
	"fmt"
	"strings"
)

// An ID is a human readable string that identifies an attribute.
type ID string

// A Set maps IDs to an opaque string. Note that a set is always passed by
// reference.
type Set map[ID]string

// List is a list of all possible attributes and the corresponding description.
var List = make(map[ID]string)

// Global is a global set of attributes.
var Global = NewSet()

// Register registers an attribute as used. All attributes are expected to be
// registered during package initialization.
func Register(id ID, help string) {
	List[id] = help
}

// NewSet creates a new set of attributes.
func NewSet() Set {
	return make(Set)
}

// Add adds (or updates) the value of an attribute.
func (a Set) Add(key ID, value string) {
	a[key] = value
}

// Delete deletes an attribute.
func (a Set) Delete(key ID) {
	delete(a, key)
}

// Get returns the value of the specified key if it exists.
func (a Set) Get(key ID) (string, bool) {
	value, ok := a[key]
	return value, ok
}

// Override overrides attributes in the current set with the corresponding
// attributes in the specified set.
func (a Set) Override(b Set) {
	for key, value := range b {
		a[key] = value
	}
}

// String formats the attributes into a string.
func (a Set) String() string {
	var buffer []string
	for key, value := range a {
		buffer = append(buffer, fmt.Sprintf("%v->%v", key, value))
	}
	return strings.Join(buffer, ".")
}
