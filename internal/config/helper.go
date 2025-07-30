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

package config

import (
	"fmt"
	"slices"

	configpb "github.com/openconfig/lemming/proto/config"
)

// GetComponentName generates component names based on configuration
func GetComponentName(prefix string, index int32) string {
	if prefix == "" {
		return fmt.Sprintf("%d", index)
	}
	return fmt.Sprintf("%s%d", prefix, index)
}

// GetProcessByPID finds a process configuration by PID
func GetProcessByPID(config *configpb.Config, pid uint32) *configpb.ProcessConfig {
	if config.GetProcesses() == nil {
		return nil
	}
	for _, proc := range config.GetProcesses().GetProcess() {
		if proc.GetPid() == pid {
			return proc
		}
	}
	return nil
}

// GetProcessByName finds a process configuration by name
func GetProcessByName(config *configpb.Config, name string) *configpb.ProcessConfig {
	if config.GetProcesses() == nil {
		return nil
	}
	for _, proc := range config.GetProcesses().GetProcess() {
		if proc.GetName() == name {
			return proc
		}
	}
	return nil
}

// GetAllLinecardNames returns all linecard component names based on configuration
func GetAllLinecardNames(config *configpb.Config) []string {
	comp := config.GetComponents()
	if comp == nil || comp.GetLinecard() == nil {
		return nil
	}

	var names []string
	lc := comp.GetLinecard()
	for i := int32(0); i < lc.GetCount(); i++ {
		index := lc.GetStartIndex() + (i * lc.GetStep())
		names = append(names, GetComponentName(comp.GetLinecardPrefix(), index))
	}
	if names == nil {
		names = []string{}
	}
	return names
}

// GetAllFabricNames returns all fabric component names based on configuration
func GetAllFabricNames(config *configpb.Config) []string {
	comp := config.GetComponents()
	if comp == nil || comp.GetFabric() == nil {
		return nil
	}

	var names []string
	fab := comp.GetFabric()
	for i := int32(0); i < fab.GetCount(); i++ {
		index := fab.GetStartIndex() + (i * fab.GetStep())
		names = append(names, GetComponentName(comp.GetFabricPrefix(), index))
	}
	if names == nil {
		names = []string{}
	}
	return names
}

// IsValidComponentName checks if a given name matches any component in the configuration
func IsValidComponentName(config *configpb.Config, name string) bool {
	comp := config.GetComponents()
	if comp == nil {
		return false
	}

	// Check supervisors
	if name == comp.GetSupervisor1Name() || name == comp.GetSupervisor2Name() {
		return true
	}

	// Check chassis
	if name == comp.GetChassisName() {
		return true
	}

	// Check linecards
	if slices.Contains(GetAllLinecardNames(config), name) {
		return true
	}

	// Check fabrics
	if slices.Contains(GetAllFabricNames(config), name) {
		return true
	}

	return false
}

// GetInterfaceByName finds an interface configuration by name
func GetInterfaceByName(config *configpb.Config, name string) *configpb.InterfaceSpec {
	if config.GetInterfaces() == nil {
		return nil
	}
	for _, iface := range config.GetInterfaces().GetInterface() {
		if iface.GetName() == name {
			return iface
		}
	}
	return nil
}

// GetInterfaceByIndex finds an interface configuration by if_index
func GetInterfaceByIndex(config *configpb.Config, ifIndex uint32) *configpb.InterfaceSpec {
	if config.GetInterfaces() == nil {
		return nil
	}
	for _, iface := range config.GetInterfaces().GetInterface() {
		if iface.GetIfIndex() == ifIndex {
			return iface
		}
	}
	return nil
}

// GetAllInterfaceNames returns all interface names from the configuration
func GetAllInterfaceNames(config *configpb.Config) []string {
	if config.GetInterfaces() == nil {
		return nil
	}

	var names []string
	for _, iface := range config.GetInterfaces().GetInterface() {
		names = append(names, iface.GetName())
	}
	if names == nil {
		names = []string{}
	}
	return names
}

// IsValidInterfaceName checks if a given name matches any interface in the configuration
func IsValidInterfaceName(config *configpb.Config, name string) bool {
	return GetInterfaceByName(config, name) != nil
}
