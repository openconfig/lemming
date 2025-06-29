package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/testing/protocmp"

	configpb "github.com/openconfig/lemming/proto/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		configFile     string
		envConfigFile  string
		setupFiles     map[string]string
		mockEmbed      string
		wantComponents *configpb.ComponentConfig
		wantError      bool
	}{
		{
			name:       "no config file uses defaults",
			configFile: "",
			mockEmbed: `
components {
  supervisor1_name: "Supervisor1"
  supervisor2_name: "Supervisor2"
  chassis_name: "chassis"
  linecard_prefix: "Linecard"
  fabric_prefix: "Fabric"
  linecard { count: 8, start_index: 0, step: 1 }
  fabric { count: 6, start_index: 0, step: 1 }
}`,
			wantComponents: &configpb.ComponentConfig{
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
			},
			wantError: false,
		},
		{
			name:       "nonexistent config file returns error",
			configFile: "nonexistent.textproto",
			wantError:  true,
		},
		{
			name:       "valid config file",
			configFile: "test_config.textproto",
			setupFiles: map[string]string{
				"test_config.textproto": `
vendor {
  name: "TestVendor"
  model: "TestModel"
  os_version: "1.0.0"
}
components {
  supervisor1_name: "TestSup1"
  supervisor2_name: "TestSup2"
  chassis_name: "TestChassis"
  linecard_prefix: "LC"
  fabric_prefix: "FC"
  linecard {
    count: 4
    start_index: 1
    step: 2
  }
  fabric {
    count: 2
    start_index: 0
    step: 1
  }
}`,
			},
			wantComponents: &configpb.ComponentConfig{
				Supervisor1Name: "TestSup1",
				Supervisor2Name: "TestSup2",
				ChassisName:     "TestChassis",
				LinecardPrefix:  "LC",
				FabricPrefix:    "FC",
				Linecard: &configpb.ComponentTypeConfig{
					Count:      4,
					StartIndex: 1,
					Step:       2,
				},
				Fabric: &configpb.ComponentTypeConfig{
					Count:      2,
					StartIndex: 0,
					Step:       1,
				},
			},
			wantError: false,
		},
		{
			name:       "invalid config file is a fatal error",
			configFile: "invalid_config.textproto",
			setupFiles: map[string]string{
				"invalid_config.textproto": `invalid protobuf content`,
			},
			wantError: true,
		},
		{
			name:          "environment variable is a direct path",
			configFile:    "",
			envConfigFile: "env_config.textproto",
			setupFiles: map[string]string{
				"env_config.textproto": `
components {
  supervisor1_name: "EnvSup1"
  supervisor2_name: "EnvSup2"
  chassis_name: "EnvChassis"
  linecard_prefix: "EnvLC"
  fabric_prefix: "EnvFC"
  linecard {
    count: 2
    start_index: 0
    step: 1
  }
  fabric {
    count: 1
    start_index: 0
    step: 1
  }
}`,
			},
			wantComponents: &configpb.ComponentConfig{
				Supervisor1Name: "EnvSup1",
				Supervisor2Name: "EnvSup2",
				ChassisName:     "EnvChassis",
				LinecardPrefix:  "EnvLC",
				FabricPrefix:    "EnvFC",
				Linecard: &configpb.ComponentTypeConfig{
					Count:      2,
					StartIndex: 0,
					Step:       1,
				},
				Fabric: &configpb.ComponentTypeConfig{
					Count:      1,
					StartIndex: 0,
					Step:       1,
				},
			},
			wantError: false,
		},
		{
			name:          "environment variable is known vendor arista",
			configFile:    "",
			envConfigFile: "arista",
			mockEmbed: `
components {
  supervisor1_name: "AristaSup1"
  supervisor2_name: "AristaSup2"
  chassis_name: "AristaChassis"
  linecard_prefix: "AristaLC"
  fabric_prefix: "AristaFC"
  linecard {
    count: 4
    start_index: 1
    step: 1
  }
  fabric {
    count: 2
    start_index: 0
    step: 1
  }
}`,
			wantComponents: &configpb.ComponentConfig{
				Supervisor1Name: "AristaSup1",
				Supervisor2Name: "AristaSup2",
				ChassisName:     "AristaChassis",
				LinecardPrefix:  "AristaLC",
				FabricPrefix:    "AristaFC",
				Linecard: &configpb.ComponentTypeConfig{
					Count:      4,
					StartIndex: 1,
					Step:       1,
				},
				Fabric: &configpb.ComponentTypeConfig{
					Count:      2,
					StartIndex: 0,
					Step:       1,
				},
			},
			wantError: false,
		},
		{
			name:          "environment variable known vendor cisco",
			configFile:    "",
			envConfigFile: "cisco",
			mockEmbed: `
components {
  supervisor1_name: "CiscoSup1"
  supervisor2_name: "CiscoSup2"
  chassis_name: "CiscoChassis"
  linecard_prefix: "CiscoLC"
  fabric_prefix: "CiscoFC"
  linecard {
    count: 6
    start_index: 0
    step: 1
  }
  fabric {
    count: 3
    start_index: 0
    step: 1
  }
}`,
			wantComponents: &configpb.ComponentConfig{
				Supervisor1Name: "CiscoSup1",
				Supervisor2Name: "CiscoSup2",
				ChassisName:     "CiscoChassis",
				LinecardPrefix:  "CiscoLC",
				FabricPrefix:    "CiscoFC",
				Linecard: &configpb.ComponentTypeConfig{
					Count:      6,
					StartIndex: 0,
					Step:       1,
				},
				Fabric: &configpb.ComponentTypeConfig{
					Count:      3,
					StartIndex: 0,
					Step:       1,
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup temporary directory
			tmpDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			os.Chdir(tmpDir)

			// Create test files
			for filename, content := range tt.setupFiles {
				// Create directory if needed
				dir := filepath.Dir(filename)
				if dir != "." {
					os.MkdirAll(dir, 0755)
				}
				err := os.WriteFile(filename, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", filename, err)
				}
			}

			// Set environment variable if needed
			if tt.envConfigFile != "" {
				os.Setenv("LEMMING_CONFIG_FILE", tt.envConfigFile)
				defer os.Unsetenv("LEMMING_CONFIG_FILE")
			}

			// Mock embedded file parser
			if tt.mockEmbed != "" {
				originalParse := parseFromEmbeddedFn
				parseFromEmbeddedFn = func(path string) (*configpb.LemmingConfig, error) {
					config := &configpb.LemmingConfig{}
					if err := prototext.Unmarshal([]byte(tt.mockEmbed), config); err != nil {
						return nil, err
					}
					return config, nil
				}
				defer func() { parseFromEmbeddedFn = originalParse }()
			}

			// Test Load function
			config, err := Load(tt.configFile)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantError {
				if diff := cmp.Diff(tt.wantComponents, config.GetComponents(), protocmp.Transform()); diff != "" {
					t.Errorf("Components mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestDetermineConfigPath(t *testing.T) {
	tests := []struct {
		name        string
		flagValue   string
		envValue    string
		defaultFile bool
		setupFiles  []string
		expected    string
		wantError   bool
	}{
		{
			name:       "flag value exists",
			flagValue:  "flag_config.txt",
			setupFiles: []string{"flag_config.txt"},
			expected:   "flag_config.txt",
			wantError:  false,
		},
		{
			name:       "flag value doesn't exist, returns error",
			flagValue:  "nonexistent.txt",
			envValue:   "env_config.txt",
			setupFiles: []string{"env_config.txt"},
			expected:   "",
			wantError:  true,
		},
		{
			name:      "env value is known vendor arista",
			flagValue: "",
			envValue:  "arista",
			expected:  "embedded:arista_default.textproto",
			wantError: false,
		},
		{
			name:       "env value is not a known vendor, treated as a path",
			flagValue:  "",
			envValue:   "my_special_config.txt",
			setupFiles: []string{"my_special_config.txt"},
			expected:   "my_special_config.txt",
			wantError:  false,
		},
		{
			name:        "env value is not a known vendor and path does not exist, falls back to default",
			flagValue:   "",
			envValue:    "nonexistent_vendor",
			defaultFile: true,
			expected:    "embedded:lemming_default.textproto",
			wantError:   false,
		},
		{
			name:        "neither flag nor env exist, default exists",
			flagValue:   "",
			defaultFile: true,
			expected:    "embedded:lemming_default.textproto",
			wantError:   false,
		},
		{
			name:      "no files exist",
			flagValue: "",
			expected:  "embedded:lemming_default.textproto",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup temporary directory
			tmpDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			os.Chdir(tmpDir)

			// Create test files
			for _, filename := range tt.setupFiles {
				dir := filepath.Dir(filename)
				if dir != "." {
					os.MkdirAll(dir, 0755)
				}
				err := os.WriteFile(filename, []byte("test content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", filename, err)
				}
			}

			// Set environment variable if needed
			if tt.envValue != "" {
				os.Setenv("LEMMING_CONFIG_FILE", tt.envValue)
				defer os.Unsetenv("LEMMING_CONFIG_FILE")
			}

			// Test determineConfigPath
			result, err := determineConfigPath(tt.flagValue)

			if tt.wantError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return // Test is done if error was expected
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMergeWithDefaults(t *testing.T) {
	tests := []struct {
		name       string
		userConfig *configpb.LemmingConfig
		wantVendor *configpb.VendorConfig
	}{
		{
			name:       "nil config uses all defaults",
			userConfig: nil,
			wantVendor: &configpb.VendorConfig{
				Name:      "OpenConfig",
				Model:     "Lemming",
				OsVersion: "1.0.0",
			},
		},
		{
			name: "partial config merges with defaults",
			userConfig: &configpb.LemmingConfig{
				Vendor: &configpb.VendorConfig{
					Name: "CustomVendor",
				},
			},
			wantVendor: &configpb.VendorConfig{
				Name: "CustomVendor",
			},
		},
		{
			name: "empty sections use defaults",
			userConfig: &configpb.LemmingConfig{
				Components: nil,
			},
			wantVendor: &configpb.VendorConfig{
				Name:      "OpenConfig",
				Model:     "Lemming",
				OsVersion: "1.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeWithDefaults(tt.userConfig)

			if diff := cmp.Diff(tt.wantVendor, result.GetVendor(), protocmp.Transform()); diff != "" {
				t.Errorf("Vendor mismatch (-want +got):\n%s", diff)
			}

			// Verify other sections are present
			if result.GetComponents() == nil {
				t.Error("Components should not be nil")
			}
			if result.GetProcesses() == nil || len(result.GetProcesses().GetProcess()) == 0 {
				t.Error("Processes should not be empty")
			}
			if result.GetTiming() == nil {
				t.Error("Timing should not be nil")
			}
			if result.GetNetworkSimulation() == nil {
				t.Error("NetworkSimulation should not be nil")
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *configpb.LemmingConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid config",
			config: &configpb.LemmingConfig{
				Components: &configpb.ComponentConfig{
					Supervisor1Name: "Sup1",
					Supervisor2Name: "Sup2",
					ChassisName:     "Chassis",
					Linecard: &configpb.ComponentTypeConfig{
						Count:      1,
						StartIndex: 0,
						Step:       1,
					},
					Fabric: &configpb.ComponentTypeConfig{
						Count:      1,
						StartIndex: 0,
						Step:       1,
					},
				},
			},
			wantError: false,
		},
		{
			name: "missing components",
			config: &configpb.LemmingConfig{
				Components: nil,
			},
			wantError: true,
			errorMsg:  "components configuration is required",
		},
		{
			name: "duplicate process PIDs",
			config: &configpb.LemmingConfig{
				Components: &configpb.ComponentConfig{
					Supervisor1Name: "Sup1",
					Supervisor2Name: "Sup2",
					ChassisName:     "Chassis",
					Linecard: &configpb.ComponentTypeConfig{
						Count:      1,
						StartIndex: 0,
						Step:       1,
					},
					Fabric: &configpb.ComponentTypeConfig{
						Count:      1,
						StartIndex: 0,
						Step:       1,
					},
				},
				Processes: &configpb.ProcessesConfig{
					Process: []*configpb.ProcessConfig{
						{Name: "proc1", Pid: 1001},
						{Name: "proc2", Pid: 1001}, // Duplicate PID
					},
				},
			},
			wantError: true,
			errorMsg:  "duplicate PID 1001 found",
		},
		{
			name: "invalid packet loss rate",
			config: &configpb.LemmingConfig{
				Components: &configpb.ComponentConfig{
					Supervisor1Name: "Sup1",
					Supervisor2Name: "Sup2",
					ChassisName:     "Chassis",
					Linecard: &configpb.ComponentTypeConfig{
						Count:      1,
						StartIndex: 0,
						Step:       1,
					},
					Fabric: &configpb.ComponentTypeConfig{
						Count:      1,
						StartIndex: 0,
						Step:       1,
					},
				},
				NetworkSimulation: &configpb.NetworkSimConfig{
					PacketLossRate: 1.5, // Invalid rate > 1.0
				},
			},
			wantError: true,
			errorMsg:  "packet_loss_rate must be between 0.0 and 1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.config)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && tt.errorMsg != "" {
				if !containsSubstring(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestValidateComponents(t *testing.T) {
	tests := []struct {
		name      string
		config    *configpb.ComponentConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid components",
			config: &configpb.ComponentConfig{
				Supervisor1Name: "Sup1",
				Supervisor2Name: "Sup2",
				ChassisName:     "Chassis",
				Linecard: &configpb.ComponentTypeConfig{
					Count:      2,
					StartIndex: 0,
					Step:       1,
				},
				Fabric: &configpb.ComponentTypeConfig{
					Count:      1,
					StartIndex: 0,
					Step:       1,
				},
			},
			wantError: false,
		},
		{
			name: "missing chassis name",
			config: &configpb.ComponentConfig{
				Supervisor1Name: "Sup1",
				Supervisor2Name: "Sup2",
				ChassisName:     "",
			},
			wantError: true,
			errorMsg:  "chassis name is required",
		},
		{
			name: "duplicate supervisor names",
			config: &configpb.ComponentConfig{
				Supervisor1Name: "Sup1",
				Supervisor2Name: "Sup1", // Same as supervisor1
				ChassisName:     "Chassis",
			},
			wantError: true,
			errorMsg:  "supervisor1_name and supervisor2_name must be different",
		},
		{
			name: "missing linecard config",
			config: &configpb.ComponentConfig{
				Supervisor1Name: "Sup1",
				Supervisor2Name: "Sup2",
				ChassisName:     "Chassis",
				Linecard:        nil,
			},
			wantError: true,
			errorMsg:  "components section requires linecard configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateComponents(tt.config)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && tt.errorMsg != "" {
				if !containsSubstring(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestParseFromFile(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		content      string
		wantError    bool
		expectVendor string
	}{
		{
			name:     "valid textproto file",
			filename: "test.textproto",
			content: `vendor {
  name: "TestVendor"
  model: "TestModel"
}`,
			wantError:    false,
			expectVendor: "TestVendor",
		},
		{
			name:     "valid pb.txt file",
			filename: "test.pb.txt",
			content: `vendor {
  name: "PbTxtVendor"
}`,
			wantError:    false,
			expectVendor: "PbTxtVendor",
		},
		{
			name:      "unsupported file format",
			filename:  "test.json",
			content:   `{"vendor": {"name": "JsonVendor"}}`,
			wantError: true,
		},
		{
			name:      "invalid protobuf content",
			filename:  "test.textproto",
			content:   `invalid protobuf syntax`,
			wantError: true,
		},
		{
			name:      "nonexistent file",
			filename:  "nonexistent.textproto",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			os.Chdir(tmpDir)

			// Create test file if content is provided
			if tt.content != "" {
				err := os.WriteFile(tt.filename, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			config, err := parseFromFile(tt.filename)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantError && config != nil && tt.expectVendor != "" {
				if config.GetVendor().GetName() != tt.expectVendor {
					t.Errorf("Expected vendor name %q, got %q", tt.expectVendor, config.GetVendor().GetName())
				}
			}
		})
	}
}

func TestValidateProcesses(t *testing.T) {
	tests := []struct {
		name      string
		processes []*configpb.ProcessConfig
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty processes list",
			processes: []*configpb.ProcessConfig{},
			wantError: false,
		},
		{
			name: "valid processes",
			processes: []*configpb.ProcessConfig{
				{Name: "proc1", Pid: 1001},
				{Name: "proc2", Pid: 1002},
			},
			wantError: false,
		},
		{
			name: "missing process name",
			processes: []*configpb.ProcessConfig{
				{Name: "", Pid: 1001},
			},
			wantError: true,
			errorMsg:  "name is required",
		},
		{
			name: "zero PID",
			processes: []*configpb.ProcessConfig{
				{Name: "proc1", Pid: 0},
			},
			wantError: true,
			errorMsg:  "invalid PID 0",
		},
		{
			name: "duplicate process names",
			processes: []*configpb.ProcessConfig{
				{Name: "proc1", Pid: 1001},
				{Name: "proc1", Pid: 1002}, // Duplicate name
			},
			wantError: true,
			errorMsg:  "duplicate process name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProcesses(tt.processes)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && tt.errorMsg != "" {
				if !containsSubstring(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestValidateComponentType(t *testing.T) {
	tests := []struct {
		name      string
		typeName  string
		config    *configpb.ComponentTypeConfig
		wantError bool
		errorMsg  string
	}{
		{
			name:     "valid config",
			typeName: "linecard",
			config: &configpb.ComponentTypeConfig{
				Count:      4,
				StartIndex: 0,
				Step:       1,
			},
			wantError: false,
		},
		{
			name:     "zero count",
			typeName: "fabric",
			config: &configpb.ComponentTypeConfig{
				Count:      0,
				StartIndex: 0,
				Step:       1,
			},
			wantError: true,
			errorMsg:  "count must be positive",
		},
		{
			name:     "zero step",
			typeName: "linecard",
			config: &configpb.ComponentTypeConfig{
				Count:      2,
				StartIndex: 0,
				Step:       0,
			},
			wantError: true,
			errorMsg:  "step must be positive",
		},
		{
			name:     "count exceeds maximum",
			typeName: "fabric",
			config: &configpb.ComponentTypeConfig{
				Count:      100,
				StartIndex: 0,
				Step:       1,
			},
			wantError: true,
			errorMsg:  "exceeds reasonable maximum 64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateComponentType(tt.typeName, tt.config)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && tt.errorMsg != "" {
				if !containsSubstring(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestValidateNetworkSim(t *testing.T) {
	tests := []struct {
		name      string
		config    *configpb.NetworkSimConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid config",
			config: &configpb.NetworkSimConfig{
				BaseLatencyMs:   50,
				LatencyJitterMs: 20,
				PacketLossRate:  0.1,
				DefaultTtl:      64,
			},
			wantError: false,
		},
		{
			name: "negative base latency",
			config: &configpb.NetworkSimConfig{
				BaseLatencyMs: -10,
			},
			wantError: true,
			errorMsg:  "base_latency_ms must be non-negative",
		},
		{
			name: "negative jitter",
			config: &configpb.NetworkSimConfig{
				LatencyJitterMs: -5,
			},
			wantError: true,
			errorMsg:  "latency_jitter_ms must be non-negative",
		},
		{
			name: "TTL too high",
			config: &configpb.NetworkSimConfig{
				DefaultTtl: 300,
			},
			wantError: true,
			errorMsg:  "default_ttl must be between 0 and 255",
		},
		{
			name: "excessive base latency",
			config: &configpb.NetworkSimConfig{
				BaseLatencyMs: 15000,
			},
			wantError: true,
			errorMsg:  "exceeds reasonable maximum 10000ms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNetworkSim(tt.config)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && tt.errorMsg != "" {
				if !containsSubstring(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestValidateTiming(t *testing.T) {
	tests := []struct {
		name      string
		config    *configpb.TimingConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid config",
			config: &configpb.TimingConfig{
				SwitchoverDurationMs: 2000,
				RebootDurationMs:     5000,
			},
			wantError: false,
		},
		{
			name: "negative switchover duration",
			config: &configpb.TimingConfig{
				SwitchoverDurationMs: -100,
			},
			wantError: true,
			errorMsg:  "switchover_duration_ms must be non-negative",
		},
		{
			name: "excessive switchover duration",
			config: &configpb.TimingConfig{
				SwitchoverDurationMs: 700000,
			},
			wantError: true,
			errorMsg:  "exceeds reasonable maximum 600000ms",
		},
		{
			name: "excessive reboot duration",
			config: &configpb.TimingConfig{
				RebootDurationMs: 800000,
			},
			wantError: true,
			errorMsg:  "exceeds reasonable maximum 600000ms",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTiming(tt.config)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantError && err != nil && tt.errorMsg != "" {
				if !containsSubstring(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return strings.Contains(s, substr)
}
