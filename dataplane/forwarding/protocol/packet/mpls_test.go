// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package packet_test

import (
	"encoding/binary"
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ip"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/mpls"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/packettestutil"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestMPLSFields(t *testing.T) {
	tests := []packettestutil.PacketFieldTest{{
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS,
		Orig:        [][]byte{genMPLS(t, 50, 7, true, 1)},
		Queries: []packettestutil.FieldQuery{{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_LABEL, 0),
			Result: binary.BigEndian.AppendUint32([]byte{}, 50),
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TC, 0),
			Result: []byte{7},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TTL, 0),
			Result: []byte{1},
		}},
	}, {
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS,
		Orig:        [][]byte{genMPLS(t, 50, 7, false, 1), genMPLS(t, 100, 3, true, 50)},
		Queries: []packettestutil.FieldQuery{{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_LABEL, 0),
			Result: binary.BigEndian.AppendUint32([]byte{}, 50),
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TC, 0),
			Result: []byte{7},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TTL, 0),
			Result: []byte{1},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_LABEL, 1),
			Result: binary.BigEndian.AppendUint32([]byte{}, 100),
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TC, 1),
			Result: []byte{3},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TTL, 1),
			Result: []byte{50},
		}},
	}, {
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
		Orig:        [][]byte{genEth(t, "00:00:00:00:00:00", "00:00:00:00:00:00", layers.EthernetTypeMPLSUnicast), genMPLS(t, 0, 7, true, 1), genIP(t, true)},
		Queries: []packettestutil.FieldQuery{{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_LABEL, 0),
			Result: binary.BigEndian.AppendUint32([]byte{}, 0),
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TC, 0),
			Result: []byte{7},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TTL, 0),
			Result: []byte{1},
		}},
	}, {
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
		Orig:        [][]byte{genEth(t, "00:00:00:00:00:00", "00:00:00:00:00:00", layers.EthernetTypeMPLSUnicast), genMPLS(t, 3, 7, true, 1), genIP(t, false)},
		Queries: []packettestutil.FieldQuery{{
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_LABEL, 0),
			Result: binary.BigEndian.AppendUint32([]byte{}, 3),
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TC, 0),
			Result: []byte{7},
		}, {
			ID:     fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_MPLS_TTL, 0),
			Result: []byte{1},
		}},
	}}
	packettestutil.TestPacketFields("mpls", t, tests)
}

func TestMPLSEncap(t *testing.T) {
	tests := []packettestutil.PacketHeaderTest{{
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP,
		Orig:        [][]byte{genIP(t, true)},
		Updates: []packettestutil.HeaderUpdate{{
			ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS,
			Encap:  true,
			Result: [][]byte{genMPLS(t, 0, 0, true, 0), genIP(t, true)},
		}},
	}, {
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP,
		Orig:        [][]byte{genIP(t, true)},
		Updates: []packettestutil.HeaderUpdate{{
			ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS,
			Encap:  true,
			Result: [][]byte{genMPLS(t, 0, 0, true, 0), genIP(t, true)},
		}, {
			ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS,
			Encap:  true,
			Result: [][]byte{genMPLS(t, 0, 0, false, 0), genMPLS(t, 0, 0, true, 0), genIP(t, true)},
		}},
	}, {
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP,
		Orig:        [][]byte{genIP(t, true)},
		Updates: []packettestutil.HeaderUpdate{{
			ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS,
			Encap:  true,
			Result: [][]byte{genMPLS(t, 0, 0, true, 0), genIP(t, true)},
		}, {
			ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
			Encap:  true,
			Result: [][]byte{genEth(t, "00:00:00:00:00:00", "00:00:00:00:00:00", layers.EthernetTypeMPLSUnicast), genMPLS(t, 0, 0, true, 0), genIP(t, true)},
		}},
	}}
	packettestutil.TestPacketHeaders("mpls", t, tests)
}

func TestMPLSDecap(t *testing.T) {
	tests := []packettestutil.PacketHeaderTest{{
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS,
		Orig:        [][]byte{genMPLS(t, 10, 0, false, 0), genMPLS(t, 15, 0, true, 0), genIP(t, true)},
		Updates: []packettestutil.HeaderUpdate{{
			ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS,
			Encap:  false,
			Result: [][]byte{genMPLS(t, 15, 0, true, 0), genIP(t, true)},
		}},
	}, {
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS,
		Orig:        [][]byte{genMPLS(t, 15, 0, true, 0), genIP(t, true)},
		Updates: []packettestutil.HeaderUpdate{{
			ID:     fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS,
			Encap:  false,
			Result: [][]byte{genIP(t, true)},
		}},
	}}
	packettestutil.TestPacketHeaders("mpls", t, tests)
}

func genEth(t testing.TB, srcMAC, dstMAC string, eType layers.EthernetType) []byte {
	src, err := net.ParseMAC(srcMAC)
	if err != nil {
		t.Fatal(err)
	}
	dst, err := net.ParseMAC(dstMAC)
	if err != nil {
		t.Fatal(err)
	}

	b := gopacket.NewSerializeBuffer()
	l := layers.Ethernet{
		SrcMAC:       src,
		DstMAC:       dst,
		EthernetType: eType,
	}
	if err := l.SerializeTo(b, gopacket.SerializeOptions{}); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()[0:14] // gopacket pads the packet, which we don't want
}

func genIP(t testing.TB, v4 bool) []byte {
	b := gopacket.NewSerializeBuffer()
	ip := net.ParseIP("::1")
	l := layers.IPv6{
		Version: 6,
		SrcIP:   ip,
		DstIP:   ip,
	}
	var err error
	err = l.SerializeTo(b, gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true})
	if v4 {
		ip = net.ParseIP("127.0.0.1")
		l := layers.IPv4{
			Version: 4,
			SrcIP:   ip,
			DstIP:   ip,
		}
		err = l.SerializeTo(b, gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true})
	}
	if err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func genMPLS(t testing.TB, label uint32, tc uint8, bos bool, ttl uint8) []byte {
	b := gopacket.NewSerializeBuffer()
	l := layers.MPLS{
		Label:        label,
		TrafficClass: tc,
		StackBottom:  bos,
		TTL:          ttl,
	}
	err := l.SerializeTo(b, gopacket.SerializeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
