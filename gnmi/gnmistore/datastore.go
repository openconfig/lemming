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

// Package datastore contains a gNMI cache implementation that can be used as
// an on-device/fake device implementation.
package gnmistore

// GNMIMode indicates the mode in which the gNMI service operates.
type GNMIMode string

const (
	GNMIModeMetadataKey = "gnmi-mode"
	// ConfigMode indicates that the gNMI service will allow updates to
	// intended configuration, but not operational state values.
	ConfigMode GNMIMode = "config"
	// StateMode indicates that the gNMI service will allow updates to
	// operational state, but not intended configuration values.
	StateMode GNMIMode = "state"
)
