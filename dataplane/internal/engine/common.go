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

// Package engine contains funcs for interacting with the forwarding engine.
package engine

import "strings"

// IntfNameToTapName returns the connected tap interface of an interface.
func IntfNameToTapName(name string) string {
	if IsTap(name) {
		return name
	}
	return name + "-tap"
}

// TapNameToIntfName returns connected external interface from a tap interface.
func TapNameToIntfName(name string) string {
	return strings.TrimSuffix(name, "-tap")
}

// IsTap returns wether the interface is a tap interface.
func IsTap(name string) bool {
	return strings.HasSuffix(name, "-tap")
}
