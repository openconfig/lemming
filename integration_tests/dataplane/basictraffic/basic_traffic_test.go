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
