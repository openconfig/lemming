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
	"time"

	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"

	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/lemming/internal/binding"
)

const (
	ipv4PrefixLen          = 30
	ateDstNetCIDR          = "198.51.100.0/24"
	ateIndirectNH          = "203.0.113.1"
	ateIndirectNHCIDR      = ateIndirectNH + "/32"
	nhIndex                = 1
	nhgIndex               = 42
	nhIndex2               = 2
	nhgIndex2              = 52
	nhIndex3               = 3
	nhgIndex3              = 62
	defaultNetworkInstance = "DEFAULT"
)

var (
	dutPort1 = attrs.Attributes{
		Desc:    "dutPort1",
		IPv4:    "192.0.2.1",
		IPv4Len: ipv4PrefixLen,
	}

	atePort1 = attrs.Attributes{
		Name:    "port1",
		MAC:     "02:00:01:01:01:01",
		IPv4:    "192.0.2.2",
		IPv4Len: ipv4PrefixLen,
	}

	dutPort2 = attrs.Attributes{
		Desc:    "dutPort2",
		IPv4:    "192.0.2.5",
		IPv4Len: ipv4PrefixLen,
	}

	atePort2 = attrs.Attributes{
		Name:    "port2",
		MAC:     "02:00:02:01:01:01",
		IPv4:    "192.0.2.6",
		IPv4Len: ipv4PrefixLen,
	}
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.Local("."))
}

func TestTraffic(t *testing.T) {
	ondatra.DUT(t, "dut")

	ate := ondatra.ATE(t, "ate")
	ateTop := configureATE(t, ate)
	ate.OTG().PushConfig(t, ateTop)
	ate.OTG().StartProtocols(t)

	// Send some traffic to make sure neighbor cache is warmed up on the dut.
	loss := testTraffic(t, ate, ateTop, atePort1, atePort2, 10*time.Second)
	if loss > 1 {
		t.Errorf("loss %f, greater than 1", loss)
	}
}

// configureATE configures port1 and port2 on the ATE.
func configureATE(t *testing.T, ate *ondatra.ATEDevice) gosnappi.Config {
	top := gosnappi.NewConfig()

	p1 := ate.Port(t, "port1")
	p2 := ate.Port(t, "port2")

	atePort1.AddToOTG(top, p1, &dutPort1)
	atePort2.AddToOTG(top, p2, &dutPort2)
	return top
}

// testTraffic generates traffic flow from source network to
// destination network via srcEndPoint to dstEndPoint and checks for
// packet loss and returns loss percentage as float.
func testTraffic(t *testing.T, ate *ondatra.ATEDevice, top gosnappi.Config, srcEndPoint, dstEndPoint attrs.Attributes, dur time.Duration) float32 {
	otg := ate.OTG()
	top.Flows().Clear().Items()

	ipFLow := top.Flows().Add().SetName("Flow")
	ipFLow.Metrics().SetEnable(true)
	ipFLow.TxRx().Port().SetTxName(srcEndPoint.Name).SetRxNames([]string{dstEndPoint.Name})

	ipFLow.Rate().SetChoice("pps").SetPps(10)

	// OTG specifies that the order of the <flow>.Packet().Add() calls determines
	// the stack of headers that are to be used, starting at the outer-most and
	// ending with the inner-most.

	// Set up ethernet layer.
	eth := ipFLow.Packet().Add().Ethernet()
	eth.Src().SetChoice("value").SetValue(srcEndPoint.MAC)
	eth.Dst().SetChoice("value").SetValue(dstEndPoint.MAC)

	ip4 := ipFLow.Packet().Add().Ipv4()
	ip4.Src().SetChoice("value").SetValue(srcEndPoint.IPv4)
	ip4.Dst().SetChoice("value").SetValue(dstEndPoint.IPv4)
	ip4.Version().SetChoice("value").SetValue(4)

	otg.PushConfig(t, top)

	otg.StartTraffic(t)
	time.Sleep(dur)
	t.Logf("Stop traffic")
	otg.StopTraffic(t)

	time.Sleep(5 * time.Second)

	txPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().OutPkts().State())
	rxPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().InPkts().State())
	lossPct := (txPkts - rxPkts) * 100 / txPkts
	return float32(lossPct)
}
