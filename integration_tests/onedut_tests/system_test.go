/*
 Copyright 2022 Google LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package integration_test

import (
	"testing"
	"time"

	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/integration_tests/binding"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.Get("."))
}

// TestHostname verifies that the hostname configuration paths can be read,
// updated, and deleted.
//
// config_path:/system/config/hostname
// telemetry_path:/system/state/hostname
func TestHostname(t *testing.T) {
	testCases := []struct {
		description string
		hostname    string
	}{
		{"15 Letters", "abcdefghijkmnop"},
		{"15 Numbers", "123456789012345"},
		{"Single Character", "x"},
		{"Periods", "test.name.example"},
		{"63 Characters", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
	}

	dut := ondatra.DUT(t, "dut")

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			config := ocpath.Root().System().Hostname().Config()
			state := ocpath.Root().System().Hostname().State()

			gnmi.Replace(t, dut, config, testCase.hostname)

			t.Run("Get Hostname Config", func(t *testing.T) {
				configGot := gnmi.Get[string](t, dut, config)
				if configGot != testCase.hostname {
					t.Errorf("Config hostname: got %s, want %s", configGot, testCase.hostname)
				}
			})

			t.Run("Get Hostname State", func(t *testing.T) {
				stateGot := gnmi.Await(t, dut, state, 5*time.Second, testCase.hostname)
				if got, ok := stateGot.Val(); !ok || (ok && got != testCase.hostname) {
					t.Errorf("State hostname: got %v, want %s", stateGot, testCase.hostname)
				}
			})

			t.Run("Delete Hostname", func(t *testing.T) {
				gnmi.Delete(t, dut, config)
				if qs := gnmi.Lookup[string](t, dut, config); qs.IsPresent() == true {
					t.Errorf("Delete hostname fail: got %v", qs)
				}
			})
		})
	}
}

// TestDomainName verifies that the domainname configuration paths can be read,
// updated, and deleted.
//
// config_path:/system/config/domain-name
// telemetry_path:/system/state/domain-name
func TestDomainName(t *testing.T) {
	testCases := []struct {
		description string
		domainname  string
	}{
		{"15 Letters", "abcdefghijkmnop"},
		{"15 Numbers", "123456789012345"},
		{"Single Character", "x"},
		{"Dash and Underscore", "foo_bar-baz"},
		{"Periods", "test.name.example"},
		{"63 Characters", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
	}

	dut := ondatra.DUT(t, "dut")

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			config := ocpath.Root().System().DomainName().Config()
			state := ocpath.Root().System().DomainName().State()

			gnmi.Replace(t, dut, config, testCase.domainname)

			t.Run("Get Domainname Config", func(t *testing.T) {
				configGot := gnmi.Get[string](t, dut, config)
				if configGot != testCase.domainname {
					t.Errorf("Config domainname: got %s, want %s", configGot, testCase.domainname)
				}
			})

			t.Run("Get Domainname State", func(t *testing.T) {
				stateGot := gnmi.Await(t, dut, state, 5*time.Second, testCase.domainname)
				if got, ok := stateGot.Val(); !ok || (ok && got != testCase.domainname) {
					t.Errorf("State domainname: got %v, want %s", stateGot, testCase.domainname)
				}
			})

			t.Run("Delete Domainname", func(t *testing.T) {
				gnmi.Delete(t, dut, config)
				if qs := gnmi.Lookup[string](t, dut, config); qs.IsPresent() == true {
					t.Errorf("Delete domainname fail: got %v", qs)
				}
			})
		})
	}
}
