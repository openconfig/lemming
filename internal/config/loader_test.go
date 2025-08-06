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
				parseFromEmbeddedFn = func(path string) (*configpb.Config, error) {
					config := &configpb.Config{}
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
		userConfig *configpb.Config
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
			userConfig: &configpb.Config{
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
			userConfig: &configpb.Config{
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
			if result.GetInterfaces() == nil {
				t.Error("Interfaces should not be nil")
			}
			if result.GetLinkQualification() == nil {
				t.Error("LinkQualification should not be nil")
			}

			// Verify new sections have expected defaults
			if result.GetNetworkSimulation().GetPacketErrorRate() != 0.0 {
				t.Errorf("Expected PacketErrorRate 0.0, got %f", result.GetNetworkSimulation().GetPacketErrorRate())
			}
			if len(result.GetInterfaces().GetInterface()) == 0 {
				t.Error("Expected default interfaces to be present")
			}
			if result.GetLinkQualification().GetMaxBps() == 0 {
				t.Error("Expected default link qualification settings to be present")
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *configpb.Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid config",
			config: &configpb.Config{
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
			config: &configpb.Config{
				Components: nil,
			},
			wantError: true,
			errorMsg:  "components configuration is required",
		},
		{
			name: "duplicate process PIDs",
			config: &configpb.Config{
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
			config: &configpb.Config{
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
				PacketErrorRate: 0.05,
				DefaultTtl:      64,
			},
			wantError: false,
		},
		{
			name: "invalid packet loss rate",
			config: &configpb.NetworkSimConfig{
				PacketLossRate: 1.5,
			},
			wantError: true,
			errorMsg:  "packet_loss_rate must be between 0.0 and 1.0",
		},
		{
			name: "invalid packet error rate too high",
			config: &configpb.NetworkSimConfig{
				PacketErrorRate: 1.5,
			},
			wantError: true,
			errorMsg:  "packet_error_rate must be between 0.0 and 1.0",
		},
		{
			name: "invalid packet error rate negative",
			config: &configpb.NetworkSimConfig{
				PacketErrorRate: -0.1,
			},
			wantError: true,
			errorMsg:  "packet_error_rate must be between 0.0 and 1.0",
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
				SwitchoverDurationMs:  2000,
				RebootDurationMs:      5000,
				ProcessRestartDelayMs: 2000,
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
			name: "negative process restart delay",
			config: &configpb.TimingConfig{
				ProcessRestartDelayMs: -500,
			},
			wantError: true,
			errorMsg:  "process_restart_delay_ms must be non-negative",
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
		{
			name: "excessive process restart delay",
			config: &configpb.TimingConfig{
				ProcessRestartDelayMs: 400000,
			},
			wantError: true,
			errorMsg:  "exceeds reasonable maximum 300000ms",
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

func TestValidateInterfaces(t *testing.T) {
	tests := []struct {
		name      string
		config    *configpb.InterfaceConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid interfaces",
			config: &configpb.InterfaceConfig{
				Interface: []*configpb.InterfaceSpec{
					{Name: "eth0", Description: "First interface", IfIndex: 1},
					{Name: "eth1", Description: "Second interface", IfIndex: 2},
				},
			},
			wantError: false,
		},
		{
			name: "empty interfaces",
			config: &configpb.InterfaceConfig{
				Interface: []*configpb.InterfaceSpec{},
			},
			wantError: true,
			errorMsg:  "interfaces section requires at least one interface",
		},
		{
			name: "nil interfaces",
			config: &configpb.InterfaceConfig{
				Interface: nil,
			},
			wantError: true,
			errorMsg:  "interfaces section requires at least one interface",
		},
		{
			name: "missing interface name",
			config: &configpb.InterfaceConfig{
				Interface: []*configpb.InterfaceSpec{
					{Name: "", Description: "Interface with no name", IfIndex: 1},
				},
			},
			wantError: true,
			errorMsg:  "name is required",
		},
		{
			name: "duplicate interface names",
			config: &configpb.InterfaceConfig{
				Interface: []*configpb.InterfaceSpec{
					{Name: "eth0", Description: "First", IfIndex: 1},
					{Name: "eth0", Description: "Duplicate", IfIndex: 2},
				},
			},
			wantError: true,
			errorMsg:  "duplicate interface name",
		},
		{
			name: "duplicate if_index",
			config: &configpb.InterfaceConfig{
				Interface: []*configpb.InterfaceSpec{
					{Name: "eth0", Description: "First", IfIndex: 1},
					{Name: "eth1", Description: "Duplicate index", IfIndex: 1},
				},
			},
			wantError: true,
			errorMsg:  "duplicate interface if_index",
		},
		{
			name: "interface name too long",
			config: &configpb.InterfaceConfig{
				Interface: []*configpb.InterfaceSpec{
					{Name: strings.Repeat("a", 65), Description: "Long name", IfIndex: 1},
				},
			},
			wantError: true,
			errorMsg:  "too long: 65 characters",
		},
		{
			name: "description too long",
			config: &configpb.InterfaceConfig{
				Interface: []*configpb.InterfaceSpec{
					{Name: "eth0", Description: strings.Repeat("a", 256), IfIndex: 1},
				},
			},
			wantError: true,
			errorMsg:  "description too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInterfaces(tt.config)

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

func TestValidateLinkQualification(t *testing.T) {
	tests := []struct {
		name      string
		config    *configpb.LinkQualificationConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid_link_qualification_config",
			config: &configpb.LinkQualificationConfig{
				MaxBps:                400000000000,
				MaxPps:                500000000,
				MinMtu:                64,
				MaxMtu:                9000,
				MinSetupDurationMs:    1000,   // proto field 5 - REQUIRED
				MinTeardownDurationMs: 1000,   // proto field 6 - REQUIRED
				MinSampleIntervalMs:   10000,  // proto field 7 - REQUIRED
				DefaultPacketRate:     138888, // proto field 8 - REQUIRED
				DefaultPacketSize:     8184,   // proto field 9 - REQUIRED
				DefaultTestDurationMs: 5000,
				MaxHistoricalResults:  10,
			},
			wantError: false,
		},
		// Test missing/zero required timing fields (proto fields 5, 6, 7)
		{
			name: "negative_min_setup_duration_ms",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    -100, // Negative value
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
			},
			wantError: true,
			errorMsg:  "min_setup_duration_ms must be non-negative",
		},
		{
			name: "negative_min_teardown_duration_ms",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: -100, // Negative value
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
			},
			wantError: true,
			errorMsg:  "min_teardown_duration_ms must be non-negative",
		},
		{
			name: "zero_min_sample_interval_ms",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   0, // Zero value (invalid)
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
			},
			wantError: true,
			errorMsg:  "min_sample_interval_ms must be greater than 0",
		},
		{
			name: "negative_min_sample_interval_ms",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   -100, // Negative value
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
			},
			wantError: true,
			errorMsg:  "min_sample_interval_ms must be greater than 0",
		},
		// Test missing/zero required simulation fields (proto fields 8, 9)
		{
			name: "zero_default_packet_rate",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     0, // Zero value (invalid)
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
			},
			wantError: true,
			errorMsg:  "default_packet_rate must be positive",
		},
		{
			name: "zero_default_packet_size",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     0, // Zero value (invalid)
				DefaultTestDurationMs: 5000,
			},
			wantError: true,
			errorMsg:  "default_packet_size must be positive",
		},
		// Test missing/zero default test duration
		{
			name: "zero_default_test_duration_ms",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 0, // Zero value (invalid)
			},
			wantError: true,
			errorMsg:  "default_test_duration_ms must be positive",
		},
		{
			name: "negative_default_test_duration_ms",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: -5000, // Negative value (invalid)
			},
			wantError: true,
			errorMsg:  "default_test_duration_ms must be positive",
		},
		// Test excessive values
		{
			name: "excessive_min_sample_interval_ms",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   20000, // exceeds 10000ms limit
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
			},
			wantError: true,
			errorMsg:  "exceeds reasonable maximum 10000ms",
		},
		{
			name: "excessive_default_packet_rate",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     20000000000, // exceeds 10B PPS
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
			},
			wantError: true,
			errorMsg:  "exceeds reasonable maximum 10B PPS",
		},
		// Test packet size bounds
		{
			name: "packet_size_too_small",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     32, // below 64 bytes minimum
				DefaultTestDurationMs: 5000,
			},
			wantError: true,
			errorMsg:  "should be between 64 and 9000 bytes",
		},
		{
			name: "packet_size_too_large",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     10000, // above 9000 bytes maximum
				DefaultTestDurationMs: 5000,
			},
			wantError: true,
			errorMsg:  "should be between 64 and 9000 bytes",
		},
		// Test optional field bounds (capability fields)
		{
			name: "excessive_max_bps",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
				MaxBps:                20000000000000, // 20Tbps exceeds 10Tbps
			},
			wantError: true,
			errorMsg:  "exceeds reasonable maximum 10Tbps",
		},
		{
			name: "excessive_max_pps",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
				MaxPps:                20000000000, // 20B PPS exceeds 10B PPS
			},
			wantError: true,
			errorMsg:  "exceeds reasonable maximum 10B PPS",
		},
		{
			name: "invalid_max_historical_results_low",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
				MaxHistoricalResults:  1, // must be > 1
			},
			wantError: true,
			errorMsg:  "must be between 2 and 20",
		},
		{
			name: "invalid_max_historical_results_high",
			config: &configpb.LinkQualificationConfig{
				MinSetupDurationMs:    1000,
				MinTeardownDurationMs: 1000,
				MinSampleIntervalMs:   10000,
				DefaultPacketRate:     138888,
				DefaultPacketSize:     8184,
				DefaultTestDurationMs: 5000,
				MaxHistoricalResults:  25, // exceeds 20
			},
			wantError: true,
			errorMsg:  "must be between 2 and 20",
		},
		// Test nil config
		{
			name:      "nil_config",
			config:    nil,
			wantError: true,
			errorMsg:  "link qualification config is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLinkQualification(tt.config)

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

func TestIsValidRPCMethod(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   bool
	}{
		{
			name:   "valid gnoi system method",
			method: "/gnoi.system.System/Reboot",
			want:   true,
		},
		{
			name:   "valid gnoi file method",
			method: "/gnoi.file.File/Get",
			want:   true,
		},
		{
			name:   "valid custom service",
			method: "/com.example.service.MyService/DoSomething",
			want:   true,
		},
		{
			name:   "missing leading slash",
			method: "gnoi.system.System/Reboot",
			want:   false,
		},
		{
			name:   "empty string",
			method: "",
			want:   false,
		},
		{
			name:   "only slash",
			method: "/",
			want:   false,
		},
		{
			name:   "missing service part",
			method: "/Reboot",
			want:   false,
		},
		{
			name:   "missing method part",
			method: "/gnoi.system.System/",
			want:   false,
		},
		{
			name:   "no package in service",
			method: "/System/Reboot",
			want:   false,
		},
		{
			name:   "empty service part",
			method: "//Reboot",
			want:   false,
		},
		{
			name:   "empty method part",
			method: "/gnoi.system.System/",
			want:   false,
		},
		{
			name:   "too many slashes",
			method: "/gnoi.system.System/Reboot/Extra",
			want:   false,
		},
		{
			name:   "single character parts",
			method: "/a.b/C",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidRPCMethod(tt.method)
			if got != tt.want {
				t.Errorf("isValidRPCMethod(%q) = %v, want %v", tt.method, got, tt.want)
			}
		})
	}
}
