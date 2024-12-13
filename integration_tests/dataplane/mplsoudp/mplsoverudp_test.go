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

package mplsoudp

import (
	"context"
	"encoding/binary"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/integration_tests/saiutil"
	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/lemming/internal/binding"

	obind "github.com/openconfig/ondatra/binding"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

var pm = &binding.PortMgr{}

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.Local(".", binding.WithOverridePortManager(pm)))
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

var (
	dutPort1 = attrs.Attributes{
		Desc: "dutPort1",
		MAC:  "10:10:10:10:10:10",
	}

	dutPort2 = attrs.Attributes{
		Desc: "dutPort2",
		MAC:  "10:10:10:10:10:11",
	}
)

const neighborMAC = "10:10:10:10:10:12"

func configureDUT(t testing.TB, dut *ondatra.DUTDevice, hop *oc.NetworkInstance_Afts_NextHop, routePrefix string) {
	t.Helper()
	conn := dataplaneConn(t, dut)
	// Allow all traffic to L3 processing.
	mmc := saipb.NewMyMacClient(conn)
	_, err := mmc.CreateMyMac(context.Background(), &saipb.CreateMyMacRequest{
		Switch:         1,
		Priority:       proto.Uint32(1),
		MacAddress:     []byte{0, 0, 0, 0, 0, 0},
		MacAddressMask: []byte{0, 0, 0, 0, 0, 0},
	})
	if err != nil {
		t.Fatal(err)
	}
	saiutil.CreateRIF(t, dut, dut.Port(t, "port1"), dutPort1.MAC, 0)
	outRIF := saiutil.CreateRIF(t, dut, dut.Port(t, "port2"), dutPort2.MAC, 0)

	nhc := saipb.NewNextHopClient(conn)

	nh, err := nhc.CreateNextHop(context.Background(), &saipb.CreateNextHopRequest{
		Switch:            1,
		Type:              saipb.NextHopType_NEXT_HOP_TYPE_IP.Enum(),
		RouterInterfaceId: &outRIF,
		Ip:                net.ParseIP(*hop.IpAddress),
	})
	if err != nil {
		t.Fatal(err)
	}

	ip := &layers.IPv6{
		Version:    6,
		NextHeader: layers.IPProtocolUDP,
		SrcIP:      net.ParseIP(*hop.EncapHeader[0].UdpV6.SrcIp),
		DstIP:      net.ParseIP(*hop.EncapHeader[0].UdpV6.DstIp),
		HopLimit:   *hop.EncapHeader[0].UdpV6.IpTtl,
	}
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(*hop.EncapHeader[0].UdpV6.SrcUdpPort),
		DstPort: layers.UDPPort(*hop.EncapHeader[0].UdpV6.DstUdpPort),
	}
	mpls := &layers.MPLS{
		Label: uint32(hop.EncapHeader[1].Mpls.MplsLabelStack[0].(oc.UnionUint32)),
	}

	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{}, ip, udp, mpls); err != nil {
		t.Fatalf("failed to serialize headers: %v", err)
	}

	acts := []*fwdpb.ActionDesc{{
		ActionType: fwdpb.ActionType_ACTION_TYPE_REPARSE,
		Action: &fwdpb.ActionDesc_Reparse{
			Reparse: &fwdpb.ReparseActionDesc{
				HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6,
				FieldIds: []*fwdpb.PacketFieldId{ // Copy all metadata fields.
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP}},
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT}},
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT}},
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_INPUT_IFACE}},
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE}},
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TRAP_ID}},
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_GROUP_ID}},
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID}},
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF}},
				},
				// After the UDP header, the rest of the packet (original packet) will be classified as payload.
				Prepend: buf.Bytes(),
			},
		},
	}}
	actReq := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: "lucius"},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: saiserver.NHActionTable}},
		EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Exact{
			Exact: &fwdpb.ExactEntryDesc{
				Fields: []*fwdpb.PacketFieldBytes{{
					FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID}},
					Bytes:   binary.BigEndian.AppendUint64(nil, nh.GetOid()),
				}},
			},
		}},
		Actions: acts,
	}
	fwd := fwdpb.NewForwardingClient(conn)
	if _, err := fwd.TableEntryAdd(context.Background(), actReq); err != nil {
		t.Fatal(err)
	}
	saiutil.CreateRoute(t, dut, routePrefix, nh.GetOid(), 0)
	saiutil.CreateNeighbor(t, dut, *hop.IpAddress, neighborMAC, outRIF)
}

func parseMac(t testing.TB, mac string) net.HardwareAddr {
	addr, err := net.ParseMAC(mac)
	if err != nil {
		t.Fatal(err)
	}
	return addr
}

func TestMPLSoverUDP(t *testing.T) {
	hop := &oc.NetworkInstance_Afts_NextHop{
		IpAddress: ygot.String("2003::3"),
		EncapHeader: map[uint8]*oc.NetworkInstance_Afts_NextHop_EncapHeader{
			0: {
				UdpV6: &oc.NetworkInstance_Afts_NextHop_EncapHeader_UdpV6{
					SrcIp:      ygot.String("2003::1"),
					DstIp:      ygot.String("2003::2"),
					SrcUdpPort: ygot.Uint16(60000),
					DstUdpPort: ygot.Uint16(60001),
					IpTtl:      ygot.Uint8(10),
				},
			},
			1: {
				Mpls: &oc.NetworkInstance_Afts_NextHop_EncapHeader_Mpls{
					MplsLabelStack: []oc.NetworkInstance_Afts_NextHop_EncapHeader_Mpls_MplsLabelStack_Union{oc.UnionUint32(100)},
				},
			},
		},
	}

	configureDUT(t, ondatra.DUT(t, "dut"), hop, "2003::10/128")

	// Create test packet
	eth := &layers.Ethernet{
		SrcMAC:       parseMac(t, "00:00:00:00:00:01"),
		DstMAC:       parseMac(t, dutPort1.MAC),
		EthernetType: layers.EthernetTypeIPv6,
	}

	ip := &layers.IPv6{
		Version:  6,
		SrcIP:    net.ParseIP("2003::9"),
		DstIP:    net.ParseIP("2003::10"),
		HopLimit: 255,
	}
	payload := gopacket.Payload([]byte("hello world"))

	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{FixLengths: true}, eth, ip, payload); err != nil {
		t.Fatalf("failed to serialize headers: %v", err)
	}

	// Send test packet to port1.
	p1 := pm.GetPort(ondatra.DUT(t, "dut").Port(t, "port1"))
	p1.RXQueue.Write(buf.Bytes())

	// Received forwarded packet from port2.
	p2 := pm.GetPort(ondatra.DUT(t, "dut").Port(t, "port2"))
	packet := (<-p2.TXQueue.Receive()).([]byte)
	p := gopacket.NewPacket(packet, layers.LayerTypeEthernet, gopacket.Default)
	t.Log(p.Dump())

	wantEth := &layers.Ethernet{
		SrcMAC:       parseMac(t, dutPort2.MAC),
		DstMAC:       parseMac(t, neighborMAC),
		EthernetType: layers.EthernetTypeIPv6,
	}

	wantIP := &layers.IPv6{
		Version:    6,
		SrcIP:      net.ParseIP(*hop.EncapHeader[0].UdpV6.SrcIp),
		DstIP:      net.ParseIP(*hop.EncapHeader[0].UdpV6.DstIp),
		NextHeader: layers.IPProtocolUDP,
		HopLimit:   9,
		Length:     63,
	}
	wantUDP := &layers.UDP{
		SrcPort: layers.UDPPort(*hop.EncapHeader[0].UdpV6.SrcUdpPort),
		DstPort: layers.UDPPort(*hop.EncapHeader[0].UdpV6.DstUdpPort),
	}

	wantMPLS := &layers.MPLS{
		Label: 100,
	}
	buf = gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{FixLengths: true}, wantMPLS, ip, payload); err != nil {
		t.Fatalf("failed to serialize GUE headers: %v", err)
	}
	wantPayload := gopacket.Payload(buf.Bytes())

	if d := cmp.Diff(p.Layers(), []gopacket.Layer{wantEth, wantIP, wantUDP, &wantPayload}, cmpopts.IgnoreUnexported(layers.IPv6{}, layers.UDP{}, layers.MPLS{}), cmpopts.IgnoreFields(layers.UDP{}, "BaseLayer"), cmpopts.IgnoreFields(layers.Ethernet{}, "BaseLayer"), cmpopts.IgnoreFields(layers.IPv6{}, "BaseLayer")); d != "" {
		t.Fatal(d)
	}
}
