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
	t.Cleanup(func() {
		if _, err := ric.RemoveRouterInterface(context.Background(), &saipb.RemoveRouterInterfaceRequest{Oid: resp.Oid}); err != nil {
			t.Errorf("RemoveRouterInterface failed: %v", err)
		}
	})
	return resp.Oid
}

func createVirtualRouter(t testing.TB, dut *ondatra.DUTDevice) uint64 {
	t.Helper()
	conn := dataplaneConn(t, dut)
	vrc := saipb.NewVirtualRouterClient(conn)
	resp, err := vrc.CreateVirtualRouter(context.Background(), &saipb.CreateVirtualRouterRequest{
		Switch:       1,
		AdminV4State: proto.Bool(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := vrc.RemoveVirtualRouter(context.Background(), &saipb.RemoveVirtualRouterRequest{Oid: resp.Oid}); err != nil {
			t.Errorf("RemoveVirtualRouter failed: %v", err)
		}
	})
	return resp.Oid
}

func createTunnelDecapEntry(t testing.TB, dut *ondatra.DUTDevice, underlayRIF, vrID uint64, tunType saipb.TunnelType, srcPrefix, dstPrefix string) {
	t.Helper()
	conn := dataplaneConn(t, dut)

	tun := saipb.NewTunnelClient(conn)
	tunResp, err := tun.CreateTunnel(context.Background(), &saipb.CreateTunnelRequest{
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
	entryType := saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P.Enum()
	if srcPrefix == "" {
		entryType = saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP.Enum()
	}
	req := &saipb.CreateTunnelTermTableEntryRequest{
		Switch:         1,
		Type:           entryType,
		TunnelType:     &tunType,
		VrId:           &vrID,
		ActionTunnelId: &tunResp.Oid,
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

	termResp, err := tun.CreateTunnelTermTableEntry(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := tun.RemoveTunnelTermTableEntry(context.Background(), &saipb.RemoveTunnelTermTableEntryRequest{Oid: termResp.Oid}); err != nil {
			t.Errorf("RemoveTunnelTermTableEntry failed: %v", err)
		}
		if _, err := tun.RemoveTunnel(context.Background(), &saipb.RemoveTunnelRequest{Oid: tunResp.Oid}); err != nil {
			t.Errorf("RemoveTunnel failed: %v", err)
		}
	})
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
	routeReq := &saipb.CreateRouteEntryRequest{
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
	}
	_, err = rc.CreateRouteEntry(context.Background(), routeReq)
	if err != nil {
		t.Fatal(err)
	}

	nc := saipb.NewNeighborClient(conn)
	neighReq := &saipb.CreateNeighborEntryRequest{
		Entry: &saipb.NeighborEntry{
			SwitchId:  1,
			RifId:     underlay,
			IpAddress: nextHop.AsSlice(),
		},
		DstMacAddress: []byte{0x2, 0x0, 0x2, 0x1, 0x1, 0x1},
	}
	_, err = nc.CreateNeighborEntry(context.Background(), neighReq)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := nc.RemoveNeighborEntry(context.Background(), &saipb.RemoveNeighborEntryRequest{Entry: neighReq.Entry}); err != nil {
			t.Errorf("RemoveNeighborEntry failed: %v", err)
		}
		if _, err := rc.RemoveRouteEntry(context.Background(), &saipb.RemoveRouteEntryRequest{Entry: routeReq.Entry}); err != nil {
			t.Errorf("RemoveRouteEntry failed: %v", err)
		}
		if _, err := nh.RemoveNextHop(context.Background(), &saipb.RemoveNextHopRequest{Oid: nhResp.Oid}); err != nil {
			t.Errorf("RemoveNextHop failed: %v", err)
		}
		if _, err := tun.RemoveTunnel(context.Background(), &saipb.RemoveTunnelRequest{Oid: tunResp.Oid}); err != nil {
			t.Errorf("RemoveTunnel failed: %v", err)
		}
	})
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
	routeReq := &saipb.CreateRouteEntryRequest{
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
	}
	_, err = rc.CreateRouteEntry(context.Background(), routeReq)
	if err != nil {
		t.Fatal(err)
	}

	nc := saipb.NewNeighborClient(conn)
	neighReq := &saipb.CreateNeighborEntryRequest{
		Entry: &saipb.NeighborEntry{
			SwitchId:  1,
			RifId:     rif,
			IpAddress: neighborIP.AsSlice(),
		},
		DstMacAddress: neighborMac,
	}
	_, err = nc.CreateNeighborEntry(context.Background(), neighReq)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := nc.RemoveNeighborEntry(context.Background(), &saipb.RemoveNeighborEntryRequest{Entry: neighReq.Entry}); err != nil {
			t.Errorf("RemoveNeighborEntry failed: %v", err)
		}
		if _, err := rc.RemoveRouteEntry(context.Background(), &saipb.RemoveRouteEntryRequest{Entry: routeReq.Entry}); err != nil {
			t.Errorf("RemoveRouteEntry failed: %v", err)
		}
	})
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

	testTraffic(t, ate, ateTop, atePort1, atePort2)
}

func TestVXLANTunnelEncapDecap(t *testing.T) {
	ate := ondatra.ATE(t, "ate")
	ateTop := configureATE(t, ate)
	ate.OTG().PushConfig(t, ateTop)
	ate.OTG().StartProtocols(t)

	dut1 := ondatra.DUT(t, "dut1")
	dut2 := ondatra.DUT(t, "dut2")

	conn1 := dataplaneConn(t, dut1)
	conn2 := dataplaneConn(t, dut2)

	dut1OverlayVRF := createVirtualRouter(t, dut1)
	dut2OverlayVRF := createVirtualRouter(t, dut2)

	// dut1 interfaces: port1 in overlay VRF, port2 in underlay VRF 0.
	createRIF(t, dut1, "port1", dut1Port1.MAC, dut1OverlayVRF)
	dut1Rif2 := createRIF(t, dut1, "port2", dut1Port2.MAC, 0)

	// dut2 interfaces: port1 in underlay VRF 0, port2 in overlay VRF.
	dut2Rif1 := createRIF(t, dut2, "port1", dut2Port1.MAC, 0)
	dut2Rif2 := createRIF(t, dut2, "port2", dut2Port2.MAC, dut2OverlayVRF)

	// Configure dut1 VXLAN tunnel, NextHop, overlay route, and underlay neighbor.
	tunClient1 := saipb.NewTunnelClient(conn1)
	tunType := saipb.TunnelType_TUNNEL_TYPE_VXLAN
	tunResp1, err := tunClient1.CreateTunnel(context.Background(), &saipb.CreateTunnelRequest{
		Switch:            1,
		Type:              &tunType,
		UnderlayInterface: proto.Uint64(dut1Rif2),
		EncapSrcIp:        net.ParseIP("10.10.10.1").To4(),
		EncapTtlMode:      saipb.TunnelTtlMode_TUNNEL_TTL_MODE_PIPE_MODEL.Enum(),
		EncapTtlVal:       proto.Uint32(64),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := tunClient1.RemoveTunnel(context.Background(), &saipb.RemoveTunnelRequest{Oid: tunResp1.Oid}); err != nil {
			t.Errorf("RemoveTunnel failed: %v", err)
		}
	})

	tunnelMac, err := net.ParseMAC("02:00:aa:bb:cc:dd")
	if err != nil {
		t.Fatal(err)
	}

	nhClient1 := saipb.NewNextHopClient(conn1)
	nhResp1, err := nhClient1.CreateNextHop(context.Background(), &saipb.CreateNextHopRequest{
		Switch:    1,
		Type:      saipb.NextHopType_NEXT_HOP_TYPE_TUNNEL_ENCAP.Enum(),
		TunnelId:  proto.Uint64(tunResp1.Oid),
		Ip:        net.ParseIP("20.20.20.2").To4(),
		TunnelVni: proto.Uint32(5000),
		TunnelMac: tunnelMac,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := nhClient1.RemoveNextHop(context.Background(), &saipb.RemoveNextHopRequest{Oid: nhResp1.Oid}); err != nil {
			t.Errorf("RemoveNextHop failed: %v", err)
		}
	})

	rc1 := saipb.NewRouteClient(conn1)
	routeReq1 := &saipb.CreateRouteEntryRequest{
		Entry: &saipb.RouteEntry{
			SwitchId: 1,
			VrId:     dut1OverlayVRF,
			Destination: &saipb.IpPrefix{
				Addr: net.ParseIP(atePort2.IPv4).To4(),
				Mask: []byte{0xFF, 0xFF, 0xFF, 0xFF},
			},
		},
		PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD.Enum(),
		NextHopId:    &nhResp1.Oid,
	}
	if _, err := rc1.CreateRouteEntry(context.Background(), routeReq1); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := rc1.RemoveRouteEntry(context.Background(), &saipb.RemoveRouteEntryRequest{Entry: routeReq1.Entry}); err != nil {
			t.Errorf("RemoveRouteEntry failed: %v", err)
		}
	})

	nc1 := saipb.NewNeighborClient(conn1)
	dut2Port1MAC, err := net.ParseMAC(dut2Port1.MAC)
	if err != nil {
		t.Fatal(err)
	}
	neighReq1 := &saipb.CreateNeighborEntryRequest{
		Entry: &saipb.NeighborEntry{
			SwitchId:  1,
			RifId:     dut1Rif2,
			IpAddress: net.ParseIP("20.20.20.2").To4(),
		},
		DstMacAddress: dut2Port1MAC,
	}
	if _, err := nc1.CreateNeighborEntry(context.Background(), neighReq1); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := nc1.RemoveNeighborEntry(context.Background(), &saipb.RemoveNeighborEntryRequest{Entry: neighReq1.Entry}); err != nil {
			t.Errorf("RemoveNeighborEntry failed: %v", err)
		}
	})

	// Configure dut2: TunnelMap, VXLAN Tunnel with DecapMappers, P2MP TunnelTermTableEntry.
	tunClient2 := saipb.NewTunnelClient(conn2)
	mapResp, err := tunClient2.CreateTunnelMap(context.Background(), &saipb.CreateTunnelMapRequest{
		Switch: 1,
		Type:   saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID.Enum(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := tunClient2.RemoveTunnelMap(context.Background(), &saipb.RemoveTunnelMapRequest{Oid: mapResp.Oid}); err != nil {
			t.Errorf("RemoveTunnelMap failed: %v", err)
		}
	})

	tunResp2, err := tunClient2.CreateTunnel(context.Background(), &saipb.CreateTunnelRequest{
		Switch:            1,
		Type:              &tunType,
		UnderlayInterface: proto.Uint64(dut2Rif1),
		DecapMappers:      []uint64{mapResp.Oid},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := tunClient2.RemoveTunnel(context.Background(), &saipb.RemoveTunnelRequest{Oid: tunResp2.Oid}); err != nil {
			t.Errorf("RemoveTunnel failed: %v", err)
		}
	})

	termResp2, err := tunClient2.CreateTunnelTermTableEntry(context.Background(), &saipb.CreateTunnelTermTableEntryRequest{
		Switch:         1,
		Type:           saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP.Enum(),
		TunnelType:     &tunType,
		DstIp:          net.ParseIP("20.20.20.2").To4(),
		ActionTunnelId: proto.Uint64(tunResp2.Oid),
		VrId:           proto.Uint64(0),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := tunClient2.RemoveTunnelTermTableEntry(context.Background(), &saipb.RemoveTunnelTermTableEntryRequest{Oid: termResp2.Oid}); err != nil {
			t.Errorf("RemoveTunnelTermTableEntry failed: %v", err)
		}
	})

	// Add TunnelMapEntry *after* CreateTunnel (matching orchagent order).
	mapEntryResp, err := tunClient2.CreateTunnelMapEntry(context.Background(), &saipb.CreateTunnelMapEntryRequest{
		Switch:               1,
		TunnelMap:            proto.Uint64(mapResp.Oid),
		VniIdKey:             proto.Uint32(5000),
		VirtualRouterIdValue: proto.Uint64(dut2OverlayVRF),
	})
	if err != nil {
		t.Fatal(err)
	}
	removedFirst := false
	t.Cleanup(func() {
		if !removedFirst {
			if _, err := tunClient2.RemoveTunnelMapEntry(context.Background(), &saipb.RemoveTunnelMapEntryRequest{Oid: mapEntryResp.Oid}); err != nil {
				t.Errorf("RemoveTunnelMapEntry failed: %v", err)
			}
		}
	})

	// Connected route + neighbor for atePort2 in dut2OverlayVRF (NOT VRF 0).
	createConnectedRouteAndNeighbor(t, dut2, dut2Rif2, dut2OverlayVRF, fmt.Sprint(dut2Port2.IPv4, "/", dut2Port2.IPv4Len), atePort2.IPv4, atePort2.MAC)

	// Step 4: Verify 100% packet delivery for traffic.
	testTraffic(t, ate, ateTop, atePort1, atePort2)

	// Step 5: Verify negative VNI isolation by removing the map entry.
	// Traffic without the VNI map entry should NOT reach dut2OverlayVRF.
	if _, err := tunClient2.RemoveTunnelMapEntry(context.Background(), &saipb.RemoveTunnelMapEntryRequest{Oid: mapEntryResp.Oid}); err != nil {
		t.Fatal(err)
	}
	removedFirst = true

	// Strengthen negative VNI isolation test: install a route for atePort2 in underlay VRF 0 on dut2.
	// Even though underlay VRF 0 has a valid route to atePort2, packets with unmapped VNI must be dropped
	// by vniToVrfTable's default DROP action, never leaking into underlay VRF 0.
	rc2 := saipb.NewRouteClient(conn2)
	underlayRouteReq := &saipb.CreateRouteEntryRequest{
		Entry: &saipb.RouteEntry{
			SwitchId: 1,
			VrId:     0,
			Destination: &saipb.IpPrefix{
				Addr: net.ParseIP(atePort2.IPv4).To4(),
				Mask: []byte{0xFF, 0xFF, 0xFF, 0xFF},
			},
		},
		PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD.Enum(),
		NextHopId:    &dut2Rif2,
	}
	if _, err := rc2.CreateRouteEntry(context.Background(), underlayRouteReq); err != nil {
		t.Fatal(err)
	}
	removedUnderlayRoute := false
	cleanupUnderlayRoute := func() {
		if !removedUnderlayRoute {
			if _, err := rc2.RemoveRouteEntry(context.Background(), &saipb.RemoveRouteEntryRequest{Entry: underlayRouteReq.Entry}); err != nil {
				t.Errorf("RemoveRouteEntry failed: %v", err)
			}
			removedUnderlayRoute = true
		}
	}
	t.Cleanup(cleanupUnderlayRoute)

	testTrafficDropped(t, ate, ateTop, atePort1, atePort2)
	cleanupUnderlayRoute()

	// Re-add TunnelMapEntry and verify traffic resumes 100% delivery.
	reAddResp, err := tunClient2.CreateTunnelMapEntry(context.Background(), &saipb.CreateTunnelMapEntryRequest{
		Switch:               1,
		TunnelMap:            proto.Uint64(mapResp.Oid),
		VniIdKey:             proto.Uint32(5000),
		VirtualRouterIdValue: proto.Uint64(dut2OverlayVRF),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := tunClient2.RemoveTunnelMapEntry(context.Background(), &saipb.RemoveTunnelMapEntryRequest{Oid: reAddResp.Oid}); err != nil {
			t.Errorf("RemoveTunnelMapEntry failed: %v", err)
		}
	})
	testTraffic(t, ate, ateTop, atePort1, atePort2)
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
// destination network via srcEndPoint to dstEndPoint and checks that
// all packets are received.
func testTraffic(t *testing.T, ate *ondatra.ATEDevice, top gosnappi.Config, srcEndPoint, dstEndPoint attrs.Attributes) {
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
	defer otg.StopTraffic(t)

	gnmi.Await(t, otg, gnmi.OTG().Flow("Flow").Counters().OutPkts().State(), 5*time.Second, txPkts)
	gnmi.Await(t, otg, gnmi.OTG().Flow("Flow").Counters().InPkts().State(), 5*time.Second, txPkts)
}

func testTrafficDropped(t *testing.T, ate *ondatra.ATEDevice, top gosnappi.Config, srcEndPoint, dstEndPoint attrs.Attributes) {
	otg := ate.OTG()
	top.Flows().Clear().Items()

	flow := top.Flows().Add().SetName("Flow")
	flow.Metrics().SetEnable(true)
	flow.TxRx().Port().SetTxName(srcEndPoint.Name).SetRxNames([]string{dstEndPoint.Name})

	txPkts := uint64(50)
	flow.Rate().SetPps(100)
	flow.Duration().FixedPackets().SetPackets(uint32(txPkts))

	eth := flow.Packet().Add().Ethernet()
	eth.Src().SetValue(srcEndPoint.MAC)
	eth.Dst().SetValue(dstEndPoint.MAC)

	ip4 := flow.Packet().Add().Ipv4()
	ip4.Src().SetValue(srcEndPoint.IPv4)
	ip4.Dst().SetValue(dstEndPoint.IPv4)
	ip4.Version().SetValue(4)

	otg.PushConfig(t, top)
	otg.StartTraffic(t)
	defer otg.StopTraffic(t)

	gnmi.Await(t, otg, gnmi.OTG().Flow("Flow").Counters().OutPkts().State(), 5*time.Second, txPkts)
	time.Sleep(500 * time.Millisecond)
	if val, ok := gnmi.Lookup(t, otg, gnmi.OTG().Flow("Flow").Counters().InPkts().State()).Val(); ok && val != 0 {
		t.Fatalf("testTrafficDropped() got %d rx packets, want 0", val)
	}
}
