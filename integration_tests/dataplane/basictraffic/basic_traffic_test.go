// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO: Fill this in with a real test
package basictraffic

import (
	"testing"

	"github.com/openconfig/ondatra"

	"github.com/openconfig/lemming/internal/binding"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.Local("."))
}

func TestTraffic(t *testing.T) {
	ondatra.DUT(t, "dut")
	ondatra.DUT(t, "dut2")
}
