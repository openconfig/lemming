package config

import (
	"reflect"
	"testing"

	configpb "github.com/openconfig/lemming/proto/config"
)

func TestGetComponentName(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		index    int32
		expected string
	}{
		{
			name:     "with prefix",
			prefix:   "Linecard",
			index:    5,
			expected: "Linecard5",
		},
		{
			name:     "empty prefix",
			prefix:   "",
			index:    3,
			expected: "3",
		},
		{
			name:     "zero index",
			prefix:   "Fabric",
			index:    0,
			expected: "Fabric0",
		},
		{
			name:     "high index",
			prefix:   "LC",
			index:    42,
			expected: "LC42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetComponentName(tt.prefix, tt.index)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetProcessByPID(t *testing.T) {
	config := &configpb.Config{
		Processes: &configpb.ProcessesConfig{
			Process: []*configpb.ProcessConfig{
				{Name: "proc1", Pid: 1001},
				{Name: "proc2", Pid: 1002},
				{Name: "proc3", Pid: 1003},
			},
		},
	}

	tests := []struct {
		name     string
		pid      uint32
		expected *configpb.ProcessConfig
	}{
		{
			name:     "existing PID",
			pid:      1002,
			expected: &configpb.ProcessConfig{Name: "proc2", Pid: 1002},
		},
		{
			name:     "non-existing PID",
			pid:      9999,
			expected: nil,
		},
		{
			name:     "first PID",
			pid:      1001,
			expected: &configpb.ProcessConfig{Name: "proc1", Pid: 1001},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetProcessByPID(config, tt.pid)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestGetProcessByName(t *testing.T) {
	config := &configpb.Config{
		Processes: &configpb.ProcessesConfig{
			Process: []*configpb.ProcessConfig{
				{Name: "octa", Pid: 1001},
				{Name: "gribi", Pid: 1002},
				{Name: "fibd", Pid: 1003},
			},
		},
	}

	tests := []struct {
		name        string
		processName string
		expected    *configpb.ProcessConfig
	}{
		{
			name:        "existing process name",
			processName: "gribi",
			expected:    &configpb.ProcessConfig{Name: "gribi", Pid: 1002},
		},
		{
			name:        "non-existing process name",
			processName: "nonexistent",
			expected:    nil,
		},
		{
			name:        "case sensitive",
			processName: "OCTA",
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetProcessByName(config, tt.processName)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestGetAllLinecardNames(t *testing.T) {
	tests := []struct {
		name     string
		config   *configpb.Config
		expected []string
	}{
		{
			name: "normal linecard config",
			config: &configpb.Config{
				Components: &configpb.ComponentConfig{
					LinecardPrefix: "LC",
					Linecard: &configpb.ComponentTypeConfig{
						Count:      3,
						StartIndex: 1,
						Step:       2,
					},
				},
			},
			expected: []string{"LC1", "LC3", "LC5"},
		},
		{
			name: "empty prefix",
			config: &configpb.Config{
				Components: &configpb.ComponentConfig{
					LinecardPrefix: "",
					Linecard: &configpb.ComponentTypeConfig{
						Count:      2,
						StartIndex: 0,
						Step:       1,
					},
				},
			},
			expected: []string{"0", "1"},
		},
		{
			name: "nil components",
			config: &configpb.Config{
				Components: nil,
			},
			expected: nil,
		},
		{
			name: "nil linecard config",
			config: &configpb.Config{
				Components: &configpb.ComponentConfig{
					Linecard: nil,
				},
			},
			expected: nil,
		},
		{
			name: "zero count",
			config: &configpb.Config{
				Components: &configpb.ComponentConfig{
					LinecardPrefix: "LC",
					Linecard: &configpb.ComponentTypeConfig{
						Count:      0,
						StartIndex: 0,
						Step:       1,
					},
				},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAllLinecardNames(tt.config)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetAllFabricNames(t *testing.T) {
	tests := []struct {
		name     string
		config   *configpb.Config
		expected []string
	}{
		{
			name: "normal fabric config",
			config: &configpb.Config{
				Components: &configpb.ComponentConfig{
					FabricPrefix: "Fabric",
					Fabric: &configpb.ComponentTypeConfig{
						Count:      2,
						StartIndex: 0,
						Step:       1,
					},
				},
			},
			expected: []string{"Fabric0", "Fabric1"},
		},
		{
			name: "fabric with step",
			config: &configpb.Config{
				Components: &configpb.ComponentConfig{
					FabricPrefix: "FC",
					Fabric: &configpb.ComponentTypeConfig{
						Count:      3,
						StartIndex: 2,
						Step:       3,
					},
				},
			},
			expected: []string{"FC2", "FC5", "FC8"},
		},
		{
			name: "nil components",
			config: &configpb.Config{
				Components: nil,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAllFabricNames(tt.config)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsValidComponentName(t *testing.T) {
	config := &configpb.Config{
		Components: &configpb.ComponentConfig{
			Supervisor1Name: "Supervisor1",
			Supervisor2Name: "Supervisor2",
			ChassisName:     "chassis",
			LinecardPrefix:  "Linecard",
			FabricPrefix:    "Fabric",
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
	}

	tests := []struct {
		name          string
		componentName string
		expected      bool
	}{
		{
			name:          "valid supervisor1",
			componentName: "Supervisor1",
			expected:      true,
		},
		{
			name:          "valid supervisor2",
			componentName: "Supervisor2",
			expected:      true,
		},
		{
			name:          "valid chassis",
			componentName: "chassis",
			expected:      true,
		},
		{
			name:          "valid linecard",
			componentName: "Linecard0",
			expected:      true,
		},
		{
			name:          "valid linecard second",
			componentName: "Linecard1",
			expected:      true,
		},
		{
			name:          "valid fabric",
			componentName: "Fabric0",
			expected:      true,
		},
		{
			name:          "invalid linecard",
			componentName: "Linecard2",
			expected:      false,
		},
		{
			name:          "invalid component",
			componentName: "InvalidComponent",
			expected:      false,
		},
		{
			name:          "empty string",
			componentName: "",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidComponentName(config, tt.componentName)
			if result != tt.expected {
				t.Errorf("Expected %v for component %q, got %v", tt.expected, tt.componentName, result)
			}
		})
	}
}

func TestIsValidComponentNameNilConfig(t *testing.T) {
	tests := []struct {
		name          string
		config        *configpb.Config
		componentName string
		expected      bool
	}{
		{
			name:          "nil config",
			config:        nil,
			componentName: "anything",
			expected:      false,
		},
		{
			name: "nil components",
			config: &configpb.Config{
				Components: nil,
			},
			componentName: "anything",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidComponentName(tt.config, tt.componentName)
			if result != tt.expected {
				t.Errorf("Expected %v for component %q, got %v", tt.expected, tt.componentName, result)
			}
		})
	}
}

func TestGetInterfaceByName(t *testing.T) {
	config := &configpb.Config{
		Interfaces: &configpb.InterfaceConfig{
			Interface: []*configpb.InterfaceSpec{
				{Name: "eth0", Description: "First interface", IfIndex: 1},
				{Name: "eth1", Description: "Second interface", IfIndex: 2},
				{Name: "Ethernet1/1", Description: "Third interface", IfIndex: 3},
			},
		},
	}

	tests := []struct {
		name          string
		interfaceName string
		expected      *configpb.InterfaceSpec
	}{
		{
			name:          "existing interface eth0",
			interfaceName: "eth0",
			expected:      &configpb.InterfaceSpec{Name: "eth0", Description: "First interface", IfIndex: 1},
		},
		{
			name:          "existing interface Ethernet1/1",
			interfaceName: "Ethernet1/1",
			expected:      &configpb.InterfaceSpec{Name: "Ethernet1/1", Description: "Third interface", IfIndex: 3},
		},
		{
			name:          "non-existing interface",
			interfaceName: "eth99",
			expected:      nil,
		},
		{
			name:          "empty string",
			interfaceName: "",
			expected:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetInterfaceByName(config, tt.interfaceName)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestGetInterfaceByIndex(t *testing.T) {
	config := &configpb.Config{
		Interfaces: &configpb.InterfaceConfig{
			Interface: []*configpb.InterfaceSpec{
				{Name: "eth0", Description: "First interface", IfIndex: 1},
				{Name: "eth1", Description: "Second interface", IfIndex: 2},
				{Name: "Ethernet1/1", Description: "Third interface", IfIndex: 10},
			},
		},
	}

	tests := []struct {
		name     string
		ifIndex  uint32
		expected *configpb.InterfaceSpec
	}{
		{
			name:     "existing index 1",
			ifIndex:  1,
			expected: &configpb.InterfaceSpec{Name: "eth0", Description: "First interface", IfIndex: 1},
		},
		{
			name:     "existing index 10",
			ifIndex:  10,
			expected: &configpb.InterfaceSpec{Name: "Ethernet1/1", Description: "Third interface", IfIndex: 10},
		},
		{
			name:     "non-existing index",
			ifIndex:  99,
			expected: nil,
		},
		{
			name:     "zero index",
			ifIndex:  0,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetInterfaceByIndex(config, tt.ifIndex)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestGetAllInterfaceNames(t *testing.T) {
	tests := []struct {
		name     string
		config   *configpb.Config
		expected []string
	}{
		{
			name: "normal interface config",
			config: &configpb.Config{
				Interfaces: &configpb.InterfaceConfig{
					Interface: []*configpb.InterfaceSpec{
						{Name: "eth0", IfIndex: 1},
						{Name: "eth1", IfIndex: 2},
						{Name: "Ethernet1/1", IfIndex: 3},
					},
				},
			},
			expected: []string{"eth0", "eth1", "Ethernet1/1"},
		},
		{
			name: "empty interfaces",
			config: &configpb.Config{
				Interfaces: &configpb.InterfaceConfig{
					Interface: []*configpb.InterfaceSpec{},
				},
			},
			expected: []string{},
		},
		{
			name: "nil interfaces",
			config: &configpb.Config{
				Interfaces: nil,
			},
			expected: nil,
		},
		{
			name: "nil interface list",
			config: &configpb.Config{
				Interfaces: &configpb.InterfaceConfig{
					Interface: nil,
				},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAllInterfaceNames(tt.config)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsValidInterfaceName(t *testing.T) {
	config := &configpb.Config{
		Interfaces: &configpb.InterfaceConfig{
			Interface: []*configpb.InterfaceSpec{
				{Name: "eth0", IfIndex: 1},
				{Name: "eth1", IfIndex: 2},
				{Name: "Ethernet1/1", IfIndex: 3},
			},
		},
	}

	tests := []struct {
		name          string
		interfaceName string
		expected      bool
	}{
		{
			name:          "valid interface eth0",
			interfaceName: "eth0",
			expected:      true,
		},
		{
			name:          "valid interface Ethernet1/1",
			interfaceName: "Ethernet1/1",
			expected:      true,
		},
		{
			name:          "invalid interface",
			interfaceName: "eth99",
			expected:      false,
		},
		{
			name:          "empty string",
			interfaceName: "",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidInterfaceName(config, tt.interfaceName)
			if result != tt.expected {
				t.Errorf("Expected %v for interface %q, got %v", tt.expected, tt.interfaceName, result)
			}
		})
	}
}

func TestInterfaceHelpersNilConfig(t *testing.T) {
	tests := []struct {
		name          string
		config        *configpb.Config
		interfaceName string
		ifIndex       uint32
	}{
		{
			name:          "nil config",
			config:        nil,
			interfaceName: "eth0",
			ifIndex:       1,
		},
		{
			name: "nil interfaces",
			config: &configpb.Config{
				Interfaces: nil,
			},
			interfaceName: "eth0",
			ifIndex:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// All functions should handle nil configs gracefully
			if result := GetInterfaceByName(tt.config, tt.interfaceName); result != nil {
				t.Errorf("GetInterfaceByName should return nil for nil config, got %+v", result)
			}
			if result := GetInterfaceByIndex(tt.config, tt.ifIndex); result != nil {
				t.Errorf("GetInterfaceByIndex should return nil for nil config, got %+v", result)
			}
			if result := GetAllInterfaceNames(tt.config); result != nil {
				t.Errorf("GetAllInterfaceNames should return nil for nil config, got %+v", result)
			}
			if result := IsValidInterfaceName(tt.config, tt.interfaceName); result {
				t.Errorf("IsValidInterfaceName should return false for nil config, got %v", result)
			}
		})
	}
}
