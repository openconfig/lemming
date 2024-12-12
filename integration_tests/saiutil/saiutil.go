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

package saiutil

import (
	"context"
	"net"
	"net/netip"
	"strconv"
	"testing"

	"github.com/openconfig/ondatra"
	obind "github.com/openconfig/ondatra/binding"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/gnmi/oc"
)

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

type configer struct {
	cfg   func(t testing.TB, s *Suite, dut *ondatra.DUTDevice)
	uncfg func(t testing.TB, s *Suite, dut *ondatra.DUTDevice)
}

func (c *configer) Config(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {
	c.cfg(t, s, dut)
}

func (c *configer) UnConfig(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {
	c.uncfg(t, s, dut)
}

func CreateRIF(t testing.TB, dut *ondatra.DUTDevice, port *ondatra.Port, smac string) uint64 {
	t.Helper()
	conn := dataplaneConn(t, dut)
	ric := saipb.NewRouterInterfaceClient(conn)
	port1ID, err := strconv.ParseUint(port.Name(), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	mac, err := net.ParseMAC(smac)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := ric.CreateRouterInterface(context.Background(), &saipb.CreateRouterInterfaceRequest{
		Switch:        1,
		PortId:        proto.Uint64(port1ID),
		Type:          saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT.Enum(),
		SrcMacAddress: mac,
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.Oid
}

func CreateRoute(t testing.TB, dut *ondatra.DUTDevice, prefix string, nexthop uint64, vrId uint64) {
	t.Helper()
	conn := dataplaneConn(t, dut)
	rc := saipb.NewRouteClient(conn)
	pre, err := netip.ParsePrefix(prefix)
	if err != nil {
		t.Fatal(err)
	}
	ip := pre.Addr().AsSlice()
	mask := net.CIDRMask(pre.Bits(), pre.Addr().BitLen())

	_, err = rc.CreateRouteEntry(context.Background(), &saipb.CreateRouteEntryRequest{
		Entry: &saipb.RouteEntry{
			SwitchId: 1,
			VrId:     vrId,
			Destination: &saipb.IpPrefix{
				Addr: ip,
				Mask: mask,
			},
		},
		PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD.Enum(),
		NextHopId:    &nexthop,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func CreateNeighbor(t testing.TB, dut *ondatra.DUTDevice, ip string, dmac string, rif uint64) {
	t.Helper()
	conn := dataplaneConn(t, dut)
	nc := saipb.NewNeighborClient(conn)
	mac, err := net.ParseMAC(dmac)
	if err != nil {
		t.Fatal(err)
	}

	_, err = nc.CreateNeighborEntry(context.Background(), &saipb.CreateNeighborEntryRequest{
		Entry: &saipb.NeighborEntry{
			SwitchId:  1,
			RifId:     rif,
			IpAddress: net.ParseIP(ip),
		},
		DstMacAddress: mac,
	})
	if err != nil {
		t.Fatal(err)
	}
}

type InterfaceRef struct {
	Intf    string
	SubIntf uint32
}

func ConfigAft(vrf string, aft *oc.NetworkInstance_Afts) *configer {
	c := &configer{}

	c.cfg = func(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {
		conn := dataplaneConn(t, dut)
		nhc := saipb.NewNextHopClient(conn)
		nhgc := saipb.NewNextHopGroupClient(conn)

		for id, nh := range aft.NextHop {
			ref := InterfaceRef{
				Intf:    nh.GetInterfaceRef().GetInterface(),
				SubIntf: nh.GetInterfaceRef().GetSubinterface(),
			}
			req := &saipb.CreateNextHopRequest{
				Switch:            1,
				Type:              saipb.NextHopType_NEXT_HOP_TYPE_IP.Enum(),
				RouterInterfaceId: proto.Uint64(s.interfaceMap[ref]),
				Ip:                net.ParseIP(*nh.IpAddress),
			}
			resp, err := nhc.CreateNextHop(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to create next hop %d: %v", id, err)
			}
			s.oc2SAINextHop[id] = resp.GetOid()
		}

		for id, nhg := range aft.NextHopGroup {
			req := &saipb.CreateNextHopGroupRequest{
				Type: saipb.NextHopGroupType_NEXT_HOP_GROUP_TYPE_ECMP_WITH_MEMBERS.Enum(),
			}
			for id, nh := range nhg.NextHop {
				req.NextHopList = append(req.NextHopList, s.oc2SAINextHop[id])
				req.NextHopMemberWeightList = append(req.NextHopMemberWeightList, uint32(nh.GetWeight()))
			}
			resp, err := nhgc.CreateNextHopGroup(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to create next hop group %d: %v", id, err)
			}
			s.oc2SAINextHopGrp[id] = resp.GetOid()
		}

		for pre, entry := range aft.Ipv4Entry {
			CreateRoute(t, dut, pre, s.oc2SAINextHopGrp[entry.GetNextHopGroup()], s.oc2SAIVRF[vrf])
		}
		for pre, entry := range aft.Ipv6Entry {
			CreateRoute(t, dut, pre, s.oc2SAINextHopGrp[entry.GetNextHopGroup()], s.oc2SAIVRF[vrf])
		}
	}
	c.uncfg = func(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {}

	return c
}

func ConfigRIF(portID string, mac string, vlan string) *configer {
	c := &configer{}

	c.cfg = func(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {
		id := CreateRIF(t, dut, dut.Port(t, portID), mac)
		s.interfaceMap[InterfaceRef{Intf: portID}] = id
	}

	c.uncfg = func(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {}
	return c
}

func ConfigVLANSubIntf(portID string, index uint32, smac string, vlan uint16, vrf string) *configer {
	c := &configer{}

	c.cfg = func(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {
		conn := dataplaneConn(t, dut)
		ric := saipb.NewRouterInterfaceClient(conn)
		port1ID, err := strconv.ParseUint(dut.Port(t, portID).Name(), 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		mac, err := net.ParseMAC(smac)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := ric.CreateRouterInterface(context.Background(), &saipb.CreateRouterInterfaceRequest{
			Switch:        1,
			PortId:        proto.Uint64(port1ID),
			Type:          saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_SUB_PORT.Enum(),
			OuterVlanId:   proto.Uint32(uint32(vlan)),
			SrcMacAddress: mac,
		})
		if err != nil {
			t.Fatal(err)
		}
		s.interfaceMap[InterfaceRef{Intf: portID, SubIntf: index}] = resp.Oid
	}

	c.uncfg = func(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {}
	return c
}

func ConfigNeighbor(intf InterfaceRef, ip, mac string) *configer {
	c := &configer{}

	c.cfg = func(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {
		CreateNeighbor(t, dut, ip, mac, s.interfaceMap[intf])
	}

	c.uncfg = func(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {}
	return c
}

func ConfigVRF(name string) *configer {
	c := &configer{
		cfg: func(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {
			conn := dataplaneConn(t, dut)
			vrc := saipb.NewVirtualRouterClient(conn)
			vrc.CreateVirtualRouter(context.Background(), &saipb.CreateVirtualRouterRequest{
				Switch: 1,
			})
		},
	}
	c.uncfg = func(t testing.TB, s *Suite, dut *ondatra.DUTDevice) {}
	return c
}

func MustParseMac(t testing.TB, mac string) net.HardwareAddr {
	addr, err := net.ParseMAC(mac)
	if err != nil {
		t.Fatal(err)
	}
	return addr
}
