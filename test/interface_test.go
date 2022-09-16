package traffic_test

import (
	"testing"
	"time"

	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/gnmi/oc"
	"github.com/openconfig/ondatra/gnmi/oc/ocpath"
	"github.com/openconfig/ygot/ygot"

	kinit "github.com/openconfig/ondatra/knebind/init"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, kinit.Init)
}

func TestConfigure(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	i := &oc.Interface{}
	i.GetOrCreateEthernet().MacAddress = ygot.String("02:1a:c0:00:02:02")
	i.GetOrCreateSubinterface(0).GetOrCreateIpv4().GetOrCreateAddress("192.0.2.2").PrefixLength = ygot.Uint8(30)
	i.GetOrCreateSubinterface(0).GetOrCreateIpv6().GetOrCreateAddress("2001:db8::2").PrefixLength = ygot.Uint8(126)
	i.GetOrCreateSubinterface(0).Enabled = ygot.Bool(true)

	gnmi.Replace(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Config(), i)
	time.Sleep(time.Second)
	lastVal := gnmi.Get(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).State())

	t.Log(lastVal)
}
