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

package tunnel

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"testing"
	"time"

	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/lemming/internal/binding"

	obind "github.com/openconfig/ondatra/binding"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

const (
	ipv4PrefixLen = 30
)

var (
	dut1Port1 = attrs.Attributes{
		Desc:    "dutPort1",
		IPv4:    "192.0.2.1",
		MAC:     "02:00:01:01:01:01",
		IPv4Len: ipv4PrefixLen,
	}
	atePort1 = attrs.Attributes{
		Name:    "port1",
		MAC:     "02:00:02:01:01:01",
		IPv4:    "192.0.2.2",
		IPv4Len: ipv4PrefixLen,
	}
	dut1Port2 = attrs.Attributes{
		Desc:    "dutPort1",
		IPv4:    "192.0.2.5",
		MAC:     "02:00:03:01:01:01",
		IPv4Len: ipv4PrefixLen,
	}
	dut2Port1 = attrs.Attributes{
		Desc:    "dutPort1",
		IPv4:    "192.0.2.6",
		MAC:     "02:00:04:01:01:01",
		IPv4Len: ipv4PrefixLen,
	}
	dut2Port2 = attrs.Attributes{
		Desc:    "dutPort1",
		IPv4:    "192.0.2.9",
		MAC:     "02:00:05:01:01:01",
		IPv4Len: ipv4PrefixLen,
	}
	atePort2 = attrs.Attributes{
		Name:    "port2",
		IPv4:    "192.0.2.10",
		MAC:     "02:00:06:01:01:01",
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
	var kneDUT interface {
		DialGRPC(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
	}
	if err := obind.DUTAs(dut.RawAPIs().BindingDUT(), &lemming); err != nil {
		if err := obind.DUTAs(dut.RawAPIs().BindingDUT(), &kneDUT); err != nil {
			t.Fatalf("failed to get kne dut dut: %v", err)
		}
		conn, err := kneDUT.DialGRPC(context.Background(), "dataplane", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatal(err)
		}
		return conn
	}
	conn, err := lemming.DataplaneConn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func createRIF(t testing.TB, dut *ondatra.DUTDevice, portID string, mac string, vrID uint64) uint64 {
	t.Helper()
	conn := dataplaneConn(t, dut)

	smac, err := net.ParseMAC(mac)
	if err != nil {
		t.Fatal(err)
	}

	port1ID, err := strconv.ParseUint(dut.Port(t, portID).Name(), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	ric := saipb.NewRouterInterfaceClient(conn)

	resp, err := ric.CreateRouterInterface(context.Background(), &saipb.CreateRouterInterfaceRequest{
		Switch:          1,
		PortId:          proto.Uint64(port1ID),
		Type:            saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT.Enum(),
		SrcMacAddress:   smac,
		VirtualRouterId: &vrID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.Oid
}

func createTunnelDecapEntry(t testing.TB, dut *ondatra.DUTDevice, underlayRIF, vrID uint64, tunType saipb.TunnelType, srcPrefix, dstPrefix string) {
	t.Helper()
	conn := dataplaneConn(t, dut)

	tun := saipb.NewTunnelClient(conn)
	_, err := tun.CreateTunnel(context.Background(), &saipb.CreateTunnelRequest{
		Switch:            1,
		Type:              &tunType,
		UnderlayInterface: &underlayRIF,
		EncapEcnMode:      saipb.TunnelEncapEcnMode_TUNNEL_ENCAP_ECN_MODE_STANDARD.Enum(),
		EncapDscpMode:     saipb.TunnelDscpMode_TUNNEL_DSCP_MODE_UNIFORM_MODEL.Enum(),
		EncapTtlMode:      saipb.TunnelTtlMode_TUNNEL_TTL_MODE_UNIFORM_MODEL.Enum(),
	})
	if err != nil {
		t.Fatal(err)
	}
	req := &saipb.CreateTunnelTermTableEntryRequest{
		Switch:     1,
		Type:       saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum(),
		TunnelType: &tunType,
		VrId:       &vrID,
	}
	if srcPrefix != "" {
		src, err := netip.ParsePrefix(srcPrefix)
		if err != nil {
			t.Fatal(err)
		}
		req.SrcIp = src.Addr().AsSlice()
		req.SrcIpMask = net.CIDRMask(src.Bits(), src.Addr().BitLen())
	}
	if dstPrefix != "" {
		dst, err := netip.ParsePrefix(dstPrefix)
		if err != nil {
			t.Fatal(err)
		}
		req.DstIp = dst.Addr().AsSlice()
		req.DstIpMask = net.CIDRMask(dst.Bits(), dst.Addr().BitLen())
	}

	_, err = tun.CreateTunnelTermTableEntry(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
}

func createEncapRoute(t testing.TB, dut *ondatra.DUTDevice, underlay, vrID uint64, tunnelType saipb.TunnelType, prefix, nhIP string) {
	t.Helper()
	conn := dataplaneConn(t, dut)

	nextHop, err := netip.ParseAddr(nhIP)
	if err != nil {
		t.Fatal(err)
	}

	routePrefix, err := netip.ParsePrefix(prefix)
	if err != nil {
		t.Fatal(err)
	}

	tun := saipb.NewTunnelClient(conn)
	tunResp, err := tun.CreateTunnel(context.Background(), &saipb.CreateTunnelRequest{
		Switch:            1,
		Type:              &tunnelType,
		UnderlayInterface: &underlay,
		EncapEcnMode:      saipb.TunnelEncapEcnMode_TUNNEL_ENCAP_ECN_MODE_STANDARD.Enum(),
		EncapDscpMode:     saipb.TunnelDscpMode_TUNNEL_DSCP_MODE_UNIFORM_MODEL.Enum(),
		EncapTtlMode:      saipb.TunnelTtlMode_TUNNEL_TTL_MODE_UNIFORM_MODEL.Enum(),
	})
	if err != nil {
		t.Fatal(err)
	}

	nh := saipb.NewNextHopClient(conn)
	nhResp, err := nh.CreateNextHop(context.Background(), &saipb.CreateNextHopRequest{
		Switch:   1,
		Type:     saipb.NextHopType_NEXT_HOP_TYPE_TUNNEL_ENCAP.Enum(),
		TunnelId: &tunResp.Oid,
		Ip:       nextHop.AsSlice(),
	})
	if err != nil {
		t.Fatal(err)
	}

	rc := saipb.NewRouteClient(conn)
	_, err = rc.CreateRouteEntry(context.Background(), &saipb.CreateRouteEntryRequest{
		Entry: &saipb.RouteEntry{
			SwitchId: 1,
			VrId:     vrID,
			Destination: &saipb.IpPrefix{
				Addr: routePrefix.Addr().AsSlice(),
				Mask: net.CIDRMask(routePrefix.Bits(), routePrefix.Addr().BitLen()),
			},
		},
		PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD.Enum(),
		NextHopId:    &nhResp.Oid,
	})
	if err != nil {
		t.Fatal(err)
	}

	nc := saipb.NewNeighborClient(conn)
	_, err = nc.CreateNeighborEntry(context.Background(), &saipb.CreateNeighborEntryRequest{
		Entry: &saipb.NeighborEntry{
			SwitchId:  1,
			RifId:     underlay,
			IpAddress: nextHop.AsSlice(),
		},
		DstMacAddress: []byte{0x2, 0x0, 0x2, 0x1, 0x1, 0x1},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createConnectedRouteAndNeighbor(t testing.TB, dut *ondatra.DUTDevice, rif, vrID uint64, prefix, neighIP string, neighMac string) {
	t.Helper()
	conn := dataplaneConn(t, dut)

	routePrefix, err := netip.ParsePrefix(prefix)
	if err != nil {
		t.Fatal(err)
	}
	neighborIP, err := netip.ParseAddr(neighIP)
	if err != nil {
		t.Fatal(err)
	}
	neighborMac, err := net.ParseMAC(neighMac)
	if err != nil {
		t.Fatal(err)
	}

	rc := saipb.NewRouteClient(conn)
	_, err = rc.CreateRouteEntry(context.Background(), &saipb.CreateRouteEntryRequest{
		Entry: &saipb.RouteEntry{
			SwitchId: 1,
			VrId:     vrID,
			Destination: &saipb.IpPrefix{
				Addr: routePrefix.Addr().AsSlice(),
				Mask: net.CIDRMask(routePrefix.Bits(), routePrefix.Addr().BitLen()),
			},
		},
		PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD.Enum(),
		NextHopId:    &rif,
	})
	if err != nil {
		t.Fatal(err)
	}

	nc := saipb.NewNeighborClient(conn)
	_, err = nc.CreateNeighborEntry(context.Background(), &saipb.CreateNeighborEntryRequest{
		Entry: &saipb.NeighborEntry{
			SwitchId:  1,
			RifId:     rif,
			IpAddress: neighborIP.AsSlice(),
		},
		DstMacAddress: neighborMac,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTunnelEncapDecap(t *testing.T) {
	ate := ondatra.ATE(t, "ate")
	ateTop := configureATE(t, ate)
	ate.OTG().PushConfig(t, ateTop)
	ate.OTG().StartProtocols(t)

	dut1 := ondatra.DUT(t, "dut1")
	dut2 := ondatra.DUT(t, "dut2")

	createRIF(t, dut1, "port1", dut1Port1.MAC, 0)
	dut1Rif2 := createRIF(t, dut1, "port2", dut1Port2.MAC, 0)
	createRIF(t, dut2, "port1", dut2Port1.MAC, 0)
	dut2Rif2 := createRIF(t, dut2, "port2", dut2Port2.MAC, 0)

	createEncapRoute(t, dut1, dut1Rif2, 0, saipb.TunnelType_TUNNEL_TYPE_IPINIP, fmt.Sprint(atePort2.IPv4, "/32"), "20.20.20.20")
	createTunnelDecapEntry(t, dut2, dut2Rif2, 0, saipb.TunnelType_TUNNEL_TYPE_IPINIP, "", "20.20.20.20/32")
	createConnectedRouteAndNeighbor(t, dut2, dut2Rif2, 0, fmt.Sprint(dut2Port2.IPv4, "/", dut2Port2.IPv4Len), atePort2.IPv4, atePort2.MAC)

	loss := testTraffic(t, ate, ateTop, atePort1, atePort2)
	if loss > 1 {
		t.Errorf("loss %f, greater than 1", loss)
	}
}

// configureATE configures port1 and port2 on the ATE.
func configureATE(t *testing.T, ate *ondatra.ATEDevice) gosnappi.Config {
	top := gosnappi.NewConfig()

	p1 := ate.Port(t, "port1")
	p2 := ate.Port(t, "port2")

	atePort1.AddToOTG(top, p1, &dut1Port1)
	atePort2.AddToOTG(top, p2, &dut2Port2)
	return top
}

// testTraffic generates traffic flow from source network to
// destination network via srcEndPoint to dstEndPoint and checks for
// packet loss and returns loss percentage as float.
func testTraffic(t *testing.T, ate *ondatra.ATEDevice, top gosnappi.Config, srcEndPoint, dstEndPoint attrs.Attributes) float32 {
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
	eth.Dst().SetValue(dstEndPoint.MAC)

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
