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
	"os"
	"strings"

	"google.golang.org/protobuf/encoding/prototext"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/configs"
	configpb "github.com/openconfig/lemming/proto/config"
)

const (
	VendorArista = "arista"
	VendorNokia  = "nokia"
)

var parseFromEmbeddedFn = parseFromEmbedded

// Load loads the lemming configuration with merging and environment variable support.
// The returned configuration is immutable and should not be modified after loading.
// User config is merged with defaults for any missing sections (Not fields).
func Load(configFile string) (*configpb.Config, error) {
	configPath, err := determineConfigPath(configFile)
	if err != nil {
		return nil, err
	}
	log.Infof("Loading Lemming config from: %s", configPath)

	// Try to load user config from the determined path
	var userConfig *configpb.Config
	if embeddedPath, ok := strings.CutPrefix(configPath, "embedded:"); ok {
		userConfig, err = parseFromEmbeddedFn(embeddedPath)
	} else {
		userConfig, err = parseFromFile(configPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load user config from %s: %v", configPath, err)
	}
	if userConfig == nil {
		log.Info("No config file specified or loaded, using complete defaults")
	}

	// Merge user config with defaults
	config := mergeWithDefaults(userConfig)
	if err := validate(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// determineConfigPath determines the config file path from various sources
func determineConfigPath(flagValue string) (string, error) {
	// Use flag value if provided
	if flagValue != "" {
		if _, err := os.Stat(flagValue); err != nil {
			return "", fmt.Errorf("config file specified via flag '%s' not found: %w", flagValue, err)
		}
		return flagValue, nil
	}

	// Check environment variable
	if envConfigFile := os.Getenv("LEMMING_CONFIG_FILE"); envConfigFile != "" {
		switch strings.ToLower(envConfigFile) {
		// TODO(mtr002): Add config files for juniper and cisco
		case VendorArista, VendorNokia:
			// For known vendor presets, use the embedded config
			return fmt.Sprintf("embedded:%s_default.textproto", envConfigFile), nil
		default:
			if _, err := os.Stat(envConfigFile); err == nil {
				return envConfigFile, nil
			}
			log.Warningf("Config file from LEMMING_CONFIG_FILE (%s) does not exist, falling back to default", envConfigFile)
		}
	}

	// Use default Lemming config
	return "embedded:lemming_default.textproto", nil
}

// mergeWithDefaults merges user config with defaults for any missing sections
func mergeWithDefaults(userConfig *configpb.Config) *configpb.Config {
	config := &configpb.Config{
		Vendor:            defaultVendor(),
		Components:        defaultComponents(),
		Processes:         defaultProcesses(),
		Timing:            defaultTiming(),
		NetworkSimulation: defaultNetworkSimulation(),
		Interfaces:        defaultInterfaces(),
		LinkQualification: defaultLinkQualification(),
	}

	if userConfig == nil {
		return config
	}

	if userConfig.Vendor != nil {
		config.Vendor = userConfig.Vendor
	}

	if userConfig.Components != nil {
		config.Components = userConfig.Components
	}

	if userConfig.Processes != nil {
		config.Processes = userConfig.Processes
	}

	if userConfig.Timing != nil {
		config.Timing = userConfig.Timing
	}

	if userConfig.NetworkSimulation != nil {
		config.NetworkSimulation = userConfig.NetworkSimulation
	}

	if userConfig.Interfaces != nil {
		config.Interfaces = userConfig.Interfaces
	}

	if userConfig.LinkQualification != nil {
		config.LinkQualification = userConfig.LinkQualification
	}

	return config
}

// defaultVendor returns default vendor configuration
func defaultVendor() *configpb.VendorConfig {
	return &configpb.VendorConfig{
		Name:      "OpenConfig",
		Model:     "Lemming",
		OsVersion: "1.0.0",
	}
}

// defaultComponents returns default component configuration
func defaultComponents() *configpb.ComponentConfig {
	return &configpb.ComponentConfig{
		Supervisor1Name: "Supervisor1",
		Supervisor2Name: "Supervisor2",
		ChassisName:     "chassis",
		LinecardPrefix:  "Linecard",
		FabricPrefix:    "Fabric",
		Linecard: &configpb.ComponentTypeConfig{
			Count:      8,
			StartIndex: 0,
			Step:       1,
		},
		Fabric: &configpb.ComponentTypeConfig{
			Count:      6,
			StartIndex: 0,
			Step:       1,
		},
	}
}

// defaultProcesses returns default process configurations
func defaultProcesses() *configpb.ProcessesConfig {
	return &configpb.ProcessesConfig{
		DefaultRestartOnKill: true,
		Process: []*configpb.ProcessConfig{
			{Name: "Octa", Pid: 1001, CpuUsageUser: 1000000, CpuUsageSystem: 500000, CpuUtilization: 1, MemoryUsage: 10485760, MemoryUtilization: 2},
			{Name: "Gribi", Pid: 1002, CpuUsageUser: 1000000, CpuUsageSystem: 500000, CpuUtilization: 1, MemoryUsage: 10485760, MemoryUtilization: 2},
			{Name: "emsd", Pid: 1003, CpuUsageUser: 1000000, CpuUsageSystem: 500000, CpuUtilization: 1, MemoryUsage: 10485760, MemoryUtilization: 2},
			{Name: "kim", Pid: 1004, CpuUsageUser: 1000000, CpuUsageSystem: 500000, CpuUtilization: 1, MemoryUsage: 10485760, MemoryUtilization: 2},
			{Name: "grpc_server", Pid: 1005, CpuUsageUser: 1000000, CpuUsageSystem: 500000, CpuUtilization: 1, MemoryUsage: 10485760, MemoryUtilization: 2},
			{Name: "fibd", Pid: 1006, CpuUsageUser: 1000000, CpuUsageSystem: 500000, CpuUtilization: 1, MemoryUsage: 10485760, MemoryUtilization: 2},
			{Name: "rpd", Pid: 1007, CpuUsageUser: 1000000, CpuUsageSystem: 500000, CpuUtilization: 1, MemoryUsage: 10485760, MemoryUtilization: 2},
		},
	}
}

// defaultTiming returns default timing configuration
func defaultTiming() *configpb.TimingConfig {
	return &configpb.TimingConfig{
		SwitchoverDurationMs:  2000,
		RebootDurationMs:      2000,
		ProcessRestartDelayMs: 2000,
	}
}

// defaultNetworkSim returns default network simulation configuration
func defaultNetworkSimulation() *configpb.NetworkSimConfig {
	return &configpb.NetworkSimConfig{
		BaseLatencyMs:   50,
		LatencyJitterMs: 20,
		PacketLossRate:  0.0,
		DefaultTtl:      64,
		PacketErrorRate: 0.0,
	}
}

// defaultInterfaces returns default interface configuration
func defaultInterfaces() *configpb.InterfaceConfig {
	return &configpb.InterfaceConfig{
		Interface: []*configpb.InterfaceSpec{
			{Name: "eth0", Description: "Standard Ethernet interface", IfIndex: 1},
			{Name: "eth1", Description: "Standard Ethernet interface", IfIndex: 2},
			{Name: "eth2", Description: "Line card port 1", IfIndex: 3},
			{Name: "eth3", Description: "Line card port 2", IfIndex: 4},
			{Name: "HundredGigE0/0/0/29", Description: "100G port", IfIndex: 5},
			{Name: "FourHundredGigE0/0/0/26", Description: "400G port", IfIndex: 6},
		},
	}
}

// defaultLinkQualification returns default link qualification configuration
func defaultLinkQualification() *configpb.LinkQualificationConfig {
	return &configpb.LinkQualificationConfig{
		MaxBps:                400000000000, // 400G
		MaxPps:                500000000,    // 500M PPS
		MinMtu:                64,
		MaxMtu:                9000,
		MinSetupDurationMs:    1000,  // 1 second
		MinTeardownDurationMs: 1000,  // 1 second
		MinSampleIntervalMs:   10000, // 10 seconds
		DefaultPacketRate:     138888,
		DefaultPacketSize:     8184,
		DefaultTestDurationMs: 5000, // 5 seconds
		MaxHistoricalResults:  10,
	}
}

// parseFromEmbedded parses configuration from an embedded file
func parseFromEmbedded(path string) (*configpb.Config, error) {
	data, err := configs.FS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded config file %s: %v", path, err)
	}
	config := &configpb.Config{}
	if err := prototext.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse protobuf text config from embedded file %s: %v", path, err)
	}
	return config, nil
}

// parseFromFile parses configuration from a file without validation
func parseFromFile(filename string) (*configpb.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", filename, err)
	}

	config := &configpb.Config{}

	// Check for supported extensions
	lowerFilename := strings.ToLower(filename)

	if strings.HasSuffix(lowerFilename, ".textproto") ||
		strings.HasSuffix(lowerFilename, ".pb.txt") ||
		strings.HasSuffix(lowerFilename, ".pbtxt") {
		if err := prototext.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse protobuf text config file %s: %v", filename, err)
		}
		return config, nil
	}

	return nil, fmt.Errorf("unsupported config file format for %s", filename)
}

// validate validates the configuration for comprehensive consistency and correctness
func validate(config *configpb.Config) error {
	if config.Components == nil {
		return fmt.Errorf("components configuration is required")
	}

	if err := validateComponents(config.Components); err != nil {
		return fmt.Errorf("components validation failed: %v", err)
	}
	if config.Processes != nil {
		if err := validateProcessesConfig(config.Processes); err != nil {
			return fmt.Errorf("processes validation failed: %v", err)
		}
	}

	if config.NetworkSimulation != nil {
		if err := validateNetworkSim(config.NetworkSimulation); err != nil {
			return fmt.Errorf("network simulation validation failed: %v", err)
		}
	}

	if config.Timing != nil {
		if err := validateTiming(config.Timing); err != nil {
			return fmt.Errorf("timing validation failed: %v", err)
		}
	}

	if config.Vendor != nil {
		if err := validateVendor(config.Vendor); err != nil {
			return fmt.Errorf("vendor validation failed: %v", err)
		}
	}

	if config.Interfaces != nil {
		if err := validateInterfaces(config.Interfaces); err != nil {
			return fmt.Errorf("interfaces validation failed: %v", err)
		}
	}

	if config.LinkQualification != nil {
		if err := validateLinkQualification(config.LinkQualification); err != nil {
			return fmt.Errorf("link qualification validation failed: %v", err)
		}
	}

	return nil
}

// validateComponents validates component configuration structure and names
func validateComponents(comp *configpb.ComponentConfig) error {
	if comp.ChassisName == "" {
		return fmt.Errorf("chassis name is required")
	}
	if comp.Supervisor1Name == "" {
		return fmt.Errorf("components section requires supervisor1_name to be specified")
	}
	if comp.Supervisor2Name == "" {
		return fmt.Errorf("components section requires supervisor2_name to be specified")
	}

	// Note: linecard_prefix and fabric_prefix are optional
	// Empty string means component names will be just indices (e.g., "0", "1", "2")

	if comp.Supervisor1Name == comp.Supervisor2Name {
		return fmt.Errorf("supervisor1_name and supervisor2_name must be different")
	}

	if comp.Linecard == nil {
		return fmt.Errorf("components section requires linecard configuration")
	}
	if err := validateComponentType("linecard", comp.Linecard); err != nil {
		return err
	}

	if comp.Fabric == nil {
		return fmt.Errorf("components section requires fabric configuration")
	}
	if err := validateComponentType("fabric", comp.Fabric); err != nil {
		return err
	}
	return nil
}

// validateComponentType validates component type configuration parameters
func validateComponentType(typeName string, config *configpb.ComponentTypeConfig) error {
	if config.Count <= 0 {
		return fmt.Errorf("%s count must be positive, got %d", typeName, config.Count)
	}
	if config.Step <= 0 {
		return fmt.Errorf("%s step must be positive, got %d", typeName, config.Step)
	}

	// Reasonable upper bounds for component counts
	if config.Count > 64 {
		return fmt.Errorf("%s count %d exceeds reasonable maximum 64", typeName, config.Count)
	}

	return nil
}

// validateProcessesConfig validates the complete processes configuration structure
func validateProcessesConfig(processesConfig *configpb.ProcessesConfig) error {
	// Validate individual processes
	return validateProcesses(processesConfig.GetProcess())
}

// validateProcesses validates process configuration including PIDs and names
func validateProcesses(processes []*configpb.ProcessConfig) error {
	pidSet := make(map[uint32]bool)
	nameSet := make(map[string]bool)

	for i, proc := range processes {
		if proc.Name == "" {
			return fmt.Errorf("process[%d] name is required", i)
		}

		if proc.Pid == 0 {
			return fmt.Errorf("process[%d] '%s' has invalid PID 0", i, proc.Name)
		}

		// Check for duplicate PIDs
		if pidSet[proc.Pid] {
			return fmt.Errorf("duplicate PID %d found in process configuration", proc.Pid)
		}
		pidSet[proc.Pid] = true

		// Check for duplicate process names
		if nameSet[proc.Name] {
			return fmt.Errorf("duplicate process name '%s' found in configuration", proc.Name)
		}
		nameSet[proc.Name] = true
	}

	return nil
}

// validateNetworkSim validates network simulation parameters for realistic values
func validateNetworkSim(netSim *configpb.NetworkSimConfig) error {
	if netSim.PacketLossRate < 0 || netSim.PacketLossRate > 1 {
		return fmt.Errorf("packet_loss_rate must be between 0.0 and 1.0, got %f", netSim.PacketLossRate)
	}
	if netSim.PacketErrorRate < 0 || netSim.PacketErrorRate > 1 {
		return fmt.Errorf("packet_error_rate must be between 0.0 and 1.0, got %f", netSim.PacketErrorRate)
	}
	if netSim.BaseLatencyMs < 0 {
		return fmt.Errorf("base_latency_ms must be non-negative, got %d", netSim.BaseLatencyMs)
	}
	if netSim.LatencyJitterMs < 0 {
		return fmt.Errorf("latency_jitter_ms must be non-negative, got %d", netSim.LatencyJitterMs)
	}
	if netSim.DefaultTtl < 0 || netSim.DefaultTtl > 255 {
		return fmt.Errorf("default_ttl must be between 0 and 255, got %d", netSim.DefaultTtl)
	}
	// Validate reasonable upper bounds for network simulation
	if netSim.BaseLatencyMs > 10000 {
		return fmt.Errorf("base_latency_ms %d exceeds reasonable maximum 10000ms", netSim.BaseLatencyMs)
	}
	if netSim.LatencyJitterMs > 5000 {
		return fmt.Errorf("latency_jitter_ms %d exceeds reasonable maximum 5000ms", netSim.LatencyJitterMs)
	}
	return nil
}

// validateTiming validates timing configuration for system operations
func validateTiming(timing *configpb.TimingConfig) error {
	if timing.SwitchoverDurationMs < 0 {
		return fmt.Errorf("switchover_duration_ms must be non-negative, got %d", timing.SwitchoverDurationMs)
	}
	if timing.RebootDurationMs < 0 {
		return fmt.Errorf("reboot_duration_ms must be non-negative, got %d", timing.RebootDurationMs)
	}
	if timing.ProcessRestartDelayMs < 0 {
		return fmt.Errorf("process_restart_delay_ms must be non-negative, got %d", timing.ProcessRestartDelayMs)
	}
	if timing.SwitchoverDurationMs > 600000 {
		return fmt.Errorf("switchover_duration_ms %d exceeds reasonable maximum 600000ms", timing.SwitchoverDurationMs)
	}
	if timing.RebootDurationMs > 600000 {
		return fmt.Errorf("reboot_duration_ms %d exceeds reasonable maximum 600000ms", timing.RebootDurationMs)
	}
	if timing.ProcessRestartDelayMs > 300000 {
		return fmt.Errorf("process_restart_delay_ms %d exceeds reasonable maximum 300000ms (5 minutes)", timing.ProcessRestartDelayMs)
	}
	return nil
}

// validateVendor validates vendor configuration for reasonable values
func validateVendor(vendor *configpb.VendorConfig) error {
	// Vendor fields are optional, but if provided should not be excessively long
	if len(vendor.Name) > 64 {
		return fmt.Errorf("vendor name too long: %d characters (max 64)", len(vendor.Name))
	}
	if len(vendor.Model) > 64 {
		return fmt.Errorf("vendor model too long: %d characters (max 64)", len(vendor.Model))
	}
	if len(vendor.OsVersion) > 32 {
		return fmt.Errorf("vendor os_version too long: %d characters (max 32)", len(vendor.OsVersion))
	}
	return nil
}

// validateInterfaces validates interface configuration structure and names
func validateInterfaces(interfaces *configpb.InterfaceConfig) error {
	if len(interfaces.Interface) == 0 {
		return fmt.Errorf("interfaces section requires at least one interface")
	}

	nameSet := make(map[string]bool)
	ifIndexSet := make(map[uint32]bool)

	for i, iface := range interfaces.Interface {
		if iface.Name == "" {
			return fmt.Errorf("interface[%d] name is required", i)
		}

		// Check for duplicate interface names
		if nameSet[iface.Name] {
			return fmt.Errorf("duplicate interface name '%s' found in configuration", iface.Name)
		}
		nameSet[iface.Name] = true

		// Check for duplicate interface indices (only if if_index is specified)
		if iface.IfIndex != 0 {
			if ifIndexSet[iface.IfIndex] {
				return fmt.Errorf("duplicate interface if_index %d found for interface '%s'", iface.IfIndex, iface.Name)
			}
			ifIndexSet[iface.IfIndex] = true
		}

		// Validate reasonable interface name length
		if len(iface.Name) > 64 {
			return fmt.Errorf("interface[%d] name '%s' too long: %d characters (max 64)", i, iface.Name, len(iface.Name))
		}

		// Validate description length if provided
		if len(iface.Description) > 255 {
			return fmt.Errorf("interface[%d] '%s' description too long: %d characters (max 255)", i, iface.Name, len(iface.Description))
		}
	}
	return nil
}

// validateLinkQualification validates link qualification configuration parameters
func validateLinkQualification(lq *configpb.LinkQualificationConfig) error {
	if lq == nil {
		return fmt.Errorf("link qualification config is nil")
	}

	// Validate required timing fields
	if err := validateTimingFields(lq); err != nil {
		return err
	}

	// Validate required simulation fields
	if err := validateSimulationFields(lq); err != nil {
		return err
	}

	// Validate optional fields for reasonableness (if set)
	if err := validateCapabilityFields(lq); err != nil {
		return err
	}

	return nil
}

// validateTimingFields validates the required timing proto fields
func validateTimingFields(lq *configpb.LinkQualificationConfig) error {
	if lq.MinSetupDurationMs < 0 {
		return fmt.Errorf("min_setup_duration_ms must be non-negative, got %d", lq.MinSetupDurationMs)
	}
	if lq.MinTeardownDurationMs < 0 {
		return fmt.Errorf("min_teardown_duration_ms must be non-negative, got %d", lq.MinTeardownDurationMs)
	}

	if lq.MinSampleIntervalMs <= 0 {
		return fmt.Errorf("min_sample_interval_ms must be greater than 0, got %d", lq.MinSampleIntervalMs)
	}
	if lq.MinSampleIntervalMs > 10000 {
		return fmt.Errorf("min_sample_interval_ms %d exceeds reasonable maximum 10000ms", lq.MinSampleIntervalMs)
	}

	if lq.DefaultTestDurationMs <= 0 {
		return fmt.Errorf("default_test_duration_ms must be positive, got %d", lq.DefaultTestDurationMs)
	}
	const maxSafeMs = 600000 // 10 minutes
	if lq.MinSetupDurationMs > maxSafeMs {
		return fmt.Errorf("min_setup_duration_ms %d exceeds maximum safe value %d (time.Duration overflow)", lq.MinSetupDurationMs, maxSafeMs)
	}
	if lq.MinTeardownDurationMs > maxSafeMs {
		return fmt.Errorf("min_teardown_duration_ms %d exceeds maximum safe value %d (time.Duration overflow)", lq.MinTeardownDurationMs, maxSafeMs)
	}
	if lq.DefaultTestDurationMs > maxSafeMs {
		return fmt.Errorf("default_test_duration_ms %d exceeds maximum safe value %d (time.Duration overflow)", lq.DefaultTestDurationMs, maxSafeMs)
	}

	return nil
}

// validateSimulationFields validates the required simulation proto fields
func validateSimulationFields(lq *configpb.LinkQualificationConfig) error {
	if lq.DefaultPacketRate == 0 {
		return fmt.Errorf("default_packet_rate must be positive, got %d", lq.DefaultPacketRate)
	}
	if lq.DefaultPacketRate > 10000000000 { // 10B PPS reasonable limit
		return fmt.Errorf("default_packet_rate %d exceeds reasonable maximum 10B PPS", lq.DefaultPacketRate)
	}

	if lq.DefaultPacketSize == 0 {
		return fmt.Errorf("default_packet_size must be positive, got %d", lq.DefaultPacketSize)
	}
	if lq.DefaultPacketSize < 64 || lq.DefaultPacketSize > 9000 {
		return fmt.Errorf("default_packet_size %d should be between 64 and 9000 bytes", lq.DefaultPacketSize)
	}

	return nil
}

// validateCapabilityFields validates optional fields for reasonableness when they are set
func validateCapabilityFields(lq *configpb.LinkQualificationConfig) error {
	// Device capability limits
	if lq.MaxBps > 0 {
		if lq.MaxBps > 10000000000000 { // 10Tbps
			return fmt.Errorf("max_bps %d exceeds reasonable maximum 10Tbps", lq.MaxBps)
		}
	}
	if lq.MaxPps > 0 {
		if lq.MaxPps > 10000000000 { // 10B PPS
			return fmt.Errorf("max_pps %d exceeds reasonable maximum 10B PPS", lq.MaxPps)
		}
	}
	if lq.MaxHistoricalResults <= 1 || lq.MaxHistoricalResults > 20 {
		return fmt.Errorf("max_historical_results must be between 2 and 20, got %d", lq.MaxHistoricalResults)
	}

	return nil
}
