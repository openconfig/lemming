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
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/lemming/internal/binding"

	obind "github.com/openconfig/ondatra/binding"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
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
		MAC:     "10:10:10:10:10:10",
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
		MAC:     "10:10:10:10:10:11",
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

func dataplaneConn(t testing.TB, dut *ondatra.DUTDevice) *grpc.ClientConn {
	t.Helper()
	var lemming interface {
		DataplaneConn(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error)
	}
	if err := obind.DUTAs(dut.RawAPIs().BindingDUT(), &lemming); err != nil {
		t.Fatalf("failed to get lemming dut: %v", err)
	}
	conn, err := lemming.DataplaneConn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func configureDUT(t testing.TB, dut *ondatra.DUTDevice) {
	t.Helper()
	conn := dataplaneConn(t, dut)
	ric := saipb.NewRouterInterfaceClient(conn)
	port1ID, err := strconv.ParseUint(dut.Port(t, "port1").Name(), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	port2ID, err := strconv.ParseUint(dut.Port(t, "port2").Name(), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	// Allow all traffic to L3 processing.
	mmc := saipb.NewMyMacClient(conn)
	_, err = mmc.CreateMyMac(context.Background(), &saipb.CreateMyMacRequest{
		Switch:         1,
		Priority:       proto.Uint32(1),
		MacAddress:     []byte{0, 0, 0, 0, 0, 0},
		MacAddressMask: []byte{0, 0, 0, 0, 0, 0},
	})
	if err != nil {
		t.Fatal(err)
	}
	mac1, err := net.ParseMAC(dutPort1.MAC)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ric.CreateRouterInterface(context.Background(), &saipb.CreateRouterInterfaceRequest{
		Switch:        1,
		PortId:        proto.Uint64(port1ID),
		Type:          saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT.Enum(),
		SrcMacAddress: mac1,
	})
	if err != nil {
		t.Fatal(err)
	}
	mac2, err := net.ParseMAC(dutPort2.MAC)
	if err != nil {
		t.Fatal(err)
	}
	rif2Resp, err := ric.CreateRouterInterface(context.Background(), &saipb.CreateRouterInterfaceRequest{
		Switch:        1,
		PortId:        proto.Uint64(port2ID),
		Type:          saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT.Enum(),
		SrcMacAddress: mac2,
	})
	if err != nil {
		t.Fatal(err)
	}
	rc := saipb.NewRouteClient(conn)
	_, err = rc.CreateRouteEntry(context.Background(), &saipb.CreateRouteEntryRequest{
		Entry: &saipb.RouteEntry{
			SwitchId: 1,
			VrId:     0,
			Destination: &saipb.IpPrefix{
				Addr: []byte{192, 0, 2, 6},
				Mask: []byte{255, 255, 255, 255},
			},
		},
		PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD.Enum(),
		NextHopId:    &rif2Resp.Oid,
	})
	if err != nil {
		t.Fatal(err)
	}

	nc := saipb.NewNeighborClient(conn)
	_, err = nc.CreateNeighborEntry(context.Background(), &saipb.CreateNeighborEntryRequest{
		Entry: &saipb.NeighborEntry{
			SwitchId:  1,
			RifId:     rif2Resp.Oid,
			IpAddress: []byte{192, 0, 2, 6},
		},
		DstMacAddress: []byte{0o2, 0o0, 0o2, 0o1, 0o1, 0o1},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTraffic(t *testing.T) {
	ate := ondatra.ATE(t, "ate")
	ateTop := configureATE(t, ate)
	ate.OTG().PushConfig(t, ateTop)
	ate.OTG().StartProtocols(t)

	dut := ondatra.DUT(t, "dut")
	configureDUT(t, dut)

	loss := testTraffic(t, ate, ateTop, atePort1, dutPort1, atePort2)
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
func testTraffic(t *testing.T, ate *ondatra.ATEDevice, top gosnappi.Config, srcEndPoint, srcPeerEndpoint, dstEndPoint attrs.Attributes) float32 {
	otg := ate.OTG()
	top.Flows().Clear().Items()

	ipFlow := top.Flows().Add().SetName("Flow")
	ipFlow.Metrics().SetEnable(true)
	ipFlow.TxRx().Port().SetTxName(srcEndPoint.Name).SetRxNames([]string{dstEndPoint.Name})

	txPkts := uint64(100)
	ipFlow.Rate().SetPps(100)
	ipFlow.Duration().FixedPackets().SetPackets(uint32(txPkts))

	// Set up ethernet layer.
	eth := ipFlow.Packet().Add().Ethernet()
	eth.Src().SetValue(srcEndPoint.MAC)
	eth.Dst().SetValue(srcPeerEndpoint.MAC)

	ip4 := ipFlow.Packet().Add().Ipv4()
	ip4.Src().SetValue(srcEndPoint.IPv4)
	ip4.Dst().SetValue(dstEndPoint.IPv4)
	ip4.Version().SetValue(4)

	otg.PushConfig(t, top)
	otg.StartTraffic(t)

	gnmi.Await(t, otg, gnmi.OTG().Flow("Flow").Counters().OutPkts().State(), 5*time.Second, txPkts)
	rxPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().InPkts().State())
	lossPct := (txPkts - rxPkts) * 100 / txPkts
	return float32(lossPct)
}
