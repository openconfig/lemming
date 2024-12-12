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

package subintf

import (
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/openconfig/ondatra"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/integration_tests/saiutil"
	"github.com/openconfig/lemming/internal/attrs"
	"github.com/openconfig/lemming/internal/binding"
)

var pm = &binding.PortMgr{}

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.Local(".", binding.WithOverridePortManager(pm)))
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

func TestVLANSubIntfMatch(t *testing.T) {
	s := saiutil.NewSuite()
	s.BaseConfig = []saiutil.ConfigOp{
		saiutil.ConfigVRF("DEFAULT"),
		saiutil.ConfigVRF("NON_DEFAULT"),
		saiutil.ConfigRIF("port1", dutPort1.MAC, "DEFAULT"),
		saiutil.ConfigRIF("port2", dutPort2.MAC, "DEFAULT"),
		saiutil.ConfigVLANSubIntf("port1", 1, "10:10:10:10:10:12", 100, "NON_DEFAULT"),
	}
	s.Case = []*saiutil.Case{{
		Config: []saiutil.ConfigOp{
			saiutil.ConfigAft("NON_DEFAULT", &oc.NetworkInstance_Afts{
				Ipv4Entry: map[string]*oc.NetworkInstance_Afts_Ipv4Entry{
					"192.0.1.1/32": {
						NextHopGroup: proto.Uint64(1),
					},
				},
				NextHopGroup: map[uint64]*oc.NetworkInstance_Afts_NextHopGroup{
					1: {
						NextHop: map[uint64]*oc.NetworkInstance_Afts_NextHopGroup_NextHop{
							1: {Weight: proto.Uint64(1), Index: proto.Uint64(1)},
						},
					},
				},
				NextHop: map[uint64]*oc.NetworkInstance_Afts_NextHop{
					1: {
						IpAddress: proto.String("192.0.2.2"),
						InterfaceRef: &oc.NetworkInstance_Afts_NextHop_InterfaceRef{
							Interface:    proto.String("port2"),
							Subinterface: proto.Uint32(0),
						},
					},
				},
			}),
			saiutil.ConfigNeighbor(saiutil.InterfaceRef{Intf: "port2"}, "192.0.2.2", dutPort2.MAC),
		},
		In: &saiutil.Packet{
			Port: "port1",
			Layers: []gopacket.SerializableLayer{
				&layers.Ethernet{
					SrcMAC:       saiutil.MustParseMac(t, "10:10:10:10:10:09"),
					DstMAC:       saiutil.MustParseMac(t, "10:10:10:10:10:10"),
					EthernetType: layers.EthernetTypeDot1Q,
				},
				&layers.Dot1Q{
					Type:           layers.EthernetTypeIPv4,
					VLANIdentifier: 100,
				},
				&layers.IPv4{
					SrcIP:    net.ParseIP("192.0.2.1"),
					DstIP:    net.ParseIP("192.0.1.1"),
					TTL:      10,
					Version:  4,
					Protocol: layers.IPProtocolNoNextHeader,
				},
				gopacket.Payload{},
			},
		},
		Out: &saiutil.Packet{
			Port: "port2",
			Layers: []gopacket.SerializableLayer{
				&layers.Ethernet{
					SrcMAC:       saiutil.MustParseMac(t, "10:10:10:10:10:11"),
					DstMAC:       saiutil.MustParseMac(t, "10:10:10:10:10:11"),
					EthernetType: layers.EthernetTypeIPv4,
				},
				&layers.IPv4{
					SrcIP:    net.ParseIP("192.0.2.1"),
					DstIP:    net.ParseIP("192.0.1.1"),
					TTL:      9,
					Version:  4,
					Protocol: layers.IPProtocolNoNextHeader,
					Length:   42,
					Checksum: 11927,
					IHL:      5,
				},
				gopacket.Payload{},
			},
		},
	}}

	s.Run(t, pm)
}
