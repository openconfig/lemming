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

func CreateRoute(t testing.TB, dut *ondatra.DUTDevice, prefix string, nexthop uint64) {
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
			VrId:     0,
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
