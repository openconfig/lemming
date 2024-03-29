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

// Package debug contains flags for turning on or off debug messages.
package debug

// TODO: Move these flags to viper so that they're accessible via CLI.
const (
	// ExternalPortPacketTrace turns on packet tracing for lemming's ports
	// that interface with other devices.
	ExternalPortPacketTrace = true
	// TAPPortPacketTrace turns on packet tracing for lemming's TAP interface ports
	// that communicate between each external port and lemming's tasks.
	TAPPortPacketTrace = true
	// SysRIB controls whether debug messages for the main system RIB are printed.
	SysRIB = false
)
