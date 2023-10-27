// Copyright 2022 Google LLC
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

package ports

import (
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// A RecordPacketSink records the last injected packet and generates a
// notification on its channel.
type RecordPacketSink struct {
	notify     chan bool
	lastPacket *fwdpb.PacketEjectResponse
}

// PacketSink records the inject packet request and generates a notification.
func (p *RecordPacketSink) PacketSink(resp *fwdpb.PacketEjectResponse) error {
	p.lastPacket = resp
	p.notify <- true
	return nil
}

// Tests the writes to the CPU port.
func TestCpuWrite(t *testing.T) {
	name := "TestPort"
	ctx := fwdcontext.New("test", "fwd")
	ps := &RecordPacketSink{
		notify: make(chan bool),
	}
	ctx.SetPacketSink(ps.PacketSink)

	// Create a CPU port that exports the ETHER_TYPE and IP_ADDR_DST
	desc := &fwdpb.PortDesc{
		PortType: fwdpb.PortType_PORT_TYPE_CPU_PORT,
		PortId:   fwdport.MakeID(fwdobject.NewID(name)),
	}
	ethertype := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_TYPE,
			Instance: 0,
		},
	}
	ipaddress := &fwdpb.PacketFieldId{
		Field: &fwdpb.PacketField{
			FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST,
			Instance: 0,
		},
	}
	cpu := &fwdpb.CPUPortDesc{
		QueueId:        name,
		ExportFieldIds: []*fwdpb.PacketFieldId{ethertype, ipaddress},
	}
	desc.Port = &fwdpb.PortDesc_Cpu{
		Cpu: cpu,
	}
	port, err := fwdport.New(desc, ctx)
	if err != nil {
		t.Fatalf("Port creation failed, err %v.", err)
	}

	// Create an ARP packet which has an ETHER_TYPE and does not have IP_ADDR_DST.
	arp := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x08, 0x06, 0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D,
		0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
		0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c,
	}
	packet, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, arp)
	if err != nil {
		t.Fatalf("Unable to create ARP packet, err %v.", err)
	}

	// Write the packet out of the cpu port and wait for it to be received
	// by the packet sink.
	if _, err = port.Write(packet); err != nil {
		t.Fatalf("Write failed, err %v.", err)
	}
	<-ps.notify
	t.Logf("Got request %+v", ps.lastPacket)

	// Verify that the packet was received and the parsed fields only have
	// the ETHER_TYPE set to ARP.
	list := ps.lastPacket.GetParsedFields()
	if len(list) != 1 {
		t.Fatalf("Write failed to get parsed fields, got %v, want 1.", len(list))
	}
	want := &fwdpb.PacketFieldBytes{
		FieldId: ethertype,
		Bytes:   []byte{0x08, 0x06},
	}
	if !proto.Equal(list[0], want) {
		t.Fatalf("Write failed to get parsed fields, Got %v, want %v", list[0], want)
	}
}
