// Copyright 2023 Google LLC
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

package saiserver

import (
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
	saipb "github.com/openconfig/lemming/dataplane/proto/sai"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestCreateHostif(t *testing.T) {
	tests := []struct {
		desc     string
		req      *saipb.CreateHostifRequest
		want     *saipb.CreateHostifResponse
		wantAttr *saipb.HostifAttribute
		wantErr  string
	}{{
		desc: "unknown type",
		req:  &saipb.CreateHostifRequest{},
		want: &saipb.CreateHostifResponse{
			Oid: 1,
		},
		wantErr: "unknown type",
	}, {
		desc: "success netdev",
		req: &saipb.CreateHostifRequest{
			Type:  saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId: proto.Uint64(2),
		},
		want: &saipb.CreateHostifResponse{
			Oid: 3,
		},
		wantAttr: &saipb.HostifAttribute{
			Type:  saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId: proto.Uint64(2),
		},
	}, {
		desc: "success cpu port",
		req: &saipb.CreateHostifRequest{
			Type:  saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId: proto.Uint64(10),
		},
		want: &saipb.CreateHostifResponse{
			Oid: 3,
		},
		wantAttr: &saipb.HostifAttribute{
			Type:  saipb.HostifType_HOSTIF_TYPE_NETDEV.Enum(),
			ObjId: proto.Uint64(10),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{
				ctx: fwdcontext.New("foo", "foo"),
			}
			dplane.ctx.SetCPUPortSink(func(po *pktiopb.PacketOut) error { return nil }, func() {})
			c, mgr, stopFn := newTestHostif(t, dplane)
			// Create switch and ports
			mgr.StoreAttributes(mgr.NextID(), &saipb.SwitchAttribute{
				CpuPort: proto.Uint64(10),
			})
			mgr.StoreAttributes(10, &saipb.PortAttribute{
				OperStatus: saipb.PortOperStatus_PORT_OPER_STATUS_UP.Enum(),
			})
			mgr.StoreAttributes(mgr.NextID(), &saipb.PortAttribute{
				OperStatus: saipb.PortOperStatus_PORT_OPER_STATUS_DOWN.Enum(),
			})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			msgCh := make(chan *pktiopb.HostPortControlMessage, 2)
			pc, err := c.HostPortControl(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if err := pc.Send(&pktiopb.HostPortControlRequest{Msg: &pktiopb.HostPortControlRequest_Init{}}); err != nil {
				t.Fatal(err)
			}
			time.Sleep(time.Millisecond)
			processRequest := func() {
				msg, _ := pc.Recv()
				msgCh <- msg
				pc.Send(&pktiopb.HostPortControlRequest{
					Msg: &pktiopb.HostPortControlRequest_Status{
						Status: &status.Status{
							Code:    int32(codes.OK),
							Message: "",
						},
					},
				})
			}

			go processRequest()
			go processRequest()

			defer stopFn()
			got, gotErr := c.CreateHostif(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateHostif() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("CreateHostif() failed: diff(-got,+want)\n:%s", d)
			}

			_ = <-msgCh
			_ = <-msgCh

			attr := &saipb.HostifAttribute{}
			if err := mgr.PopulateAllAttributes("3", attr); err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(attr, tt.wantAttr, protocmp.Transform()); d != "" {
				t.Errorf("CreateHostif() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestRemoveHostif(t *testing.T) {
	tests := []struct {
		desc    string
		req     *saipb.RemoveHostifRequest
		want    *pktiopb.HostPortControlMessage
		wantErr string
	}{{
		desc: "sucess",
		req: &saipb.RemoveHostifRequest{
			Oid: 1,
		},
		want: &pktiopb.HostPortControlMessage{
			PortId: 1,
			Op:     pktiopb.PortOperation_PORT_OPERATION_DELETE,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{
				portIDToNID: map[string]uint64{
					"1": 10,
				},
			}
			c, mgr, stopFn := newTestHostif(t, dplane)
			defer stopFn()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mgr.StoreAttributes(1, &saipb.CreateHostifRequest{Name: []byte("eth1")})

			msgCh := make(chan *pktiopb.HostPortControlMessage, 1)
			pc, err := c.HostPortControl(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if err := pc.Send(&pktiopb.HostPortControlRequest{Msg: &pktiopb.HostPortControlRequest_Init{}}); err != nil {
				t.Fatal(err)
			}
			time.Sleep(time.Millisecond)
			go func() {
				msg, _ := pc.Recv()
				msgCh <- msg
				pc.Send(&pktiopb.HostPortControlRequest{
					Msg: &pktiopb.HostPortControlRequest_Status{
						Status: &status.Status{
							Code:    int32(codes.OK),
							Message: "",
						},
					},
				})
			}()

			_, gotErr := c.RemoveHostif(context.TODO(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("RemoveHostif() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			got := <-msgCh
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("RemoveHostif() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestCPUPacketStream(t *testing.T) {
	dplane := &fakeSwitchDataplane{
		ctx: fwdcontext.New("test", "test"),
	}
	p, err := fwdport.New(&fwdpb.PortDesc{
		PortType: fwdpb.PortType_PORT_TYPE_CPU_PORT,
		PortId:   &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "2"}},
		Port:     &fwdpb.PortDesc_Cpu{},
	}, dplane.ctx)
	if err != nil {
		t.Fatal(err)
	}

	c, mgr, stopFn := newTestHostif(t, dplane)
	mgr.StoreAttributes(1, &saipb.SwitchAttribute{
		CpuPort: proto.Uint64(2),
	})

	defer stopFn()

	s, err := c.PacketIOClient.CPUPacketStream(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Send(&pktiopb.PacketIn{}); err != nil {
		t.Fatal(err)
	}

	t.Run("send", func(t *testing.T) {
		if err := s.Send(&pktiopb.PacketIn{
			Msg: &pktiopb.PacketIn_Packet{
				Packet: &pktiopb.Packet{
					Frame: []byte("hello"),
				},
			},
		}); err != nil {
			t.Fatal(err)
		}
		time.Sleep(10 * time.Millisecond)
		want := [][]byte{[]byte("hello")}
		if d := cmp.Diff(dplane.gotPackets, want); d != "" {
			t.Errorf("PacketStream() failed: diff(-got,+want)\n:%s", d)
		}
	})
	t.Run("recv", func(t *testing.T) {
		if _, err := p.Write(createPacket(t, uint64(p.NID()))); err != nil {
			t.Fatal(err)
		}
		got, err := s.Recv()
		if err != nil {
			t.Fatal(err)
		}
		want := &pktiopb.PacketOut{
			Packet: &pktiopb.Packet{
				HostPort:  1,
				InputPort: 2,
			},
		}
		if d := cmp.Diff(got, want, protocmp.Transform(), protocmp.IgnoreFields(&pktiopb.Packet{}, protoreflect.Name("frame"))); d != "" {
			t.Errorf("PacketStream() failed: diff(-got,+want)\n:%s", d)
		}
	})
}

func createPacket(t testing.TB, nid uint64) fwdpacket.Packet {
	t.Helper()
	eth := &layers.Ethernet{
		SrcMAC:       parseMac(t, "00:00:00:00:00:01"),
		DstMAC:       parseMac(t, "00:00:00:00:00:02"),
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

	p, err := fwdpacket.New(fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET, buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	err = p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_HOST_PORT_ID, 0), fwdpacket.OpSet, binary.BigEndian.AppendUint64(nil, 1))
	if err != nil {
		t.Fatal(err)
	}
	err = p.Update(fwdpacket.NewFieldIDFromNum(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT, 0), fwdpacket.OpSet, binary.BigEndian.AppendUint64(nil, nid))
	if err != nil {
		t.Fatal(err)
	}

	return p
}

func parseMac(t testing.TB, mac string) net.HardwareAddr {
	addr, err := net.ParseMAC(mac)
	if err != nil {
		t.Fatal(err)
	}
	return addr
}

type hostifClient struct {
	saipb.HostifClient
	pktiopb.PacketIOClient
}

func newTestHostif(t testing.TB, api switchDataplaneAPI) (*hostifClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newHostif(mgr, api, srv, &dplaneopts.Options{
			HostifNetDevType: fwdpb.PortType_PORT_TYPE_KERNEL,
		})
	})
	return &hostifClient{
		HostifClient:   saipb.NewHostifClient(conn),
		PacketIOClient: pktiopb.NewPacketIOClient(conn),
	}, mgr, stopFn
}
