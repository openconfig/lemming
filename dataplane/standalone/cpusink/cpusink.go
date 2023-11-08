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

// The cpusink package subscribes to a forwarding context and writes packets to genetlink interfaces.
// It also configures the forwarding context with the IPs assigned to specific netdev interfaces.
package main

import "C"

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/vishvananda/netlink"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/internal/kernel"

	log "github.com/golang/glog"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// TODO: A better way to get this value.
const (
	contextID  = "lucius"
	addr       = "10.0.2.2:50000"
	ip2meTable = "ip2me"
)

//export StartSink
func StartSink() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	sink := &cpuSink{
		client:       fwdpb.NewForwardingClient(conn),
		ethDevToPort: make(map[string]string),
	}

	go func() {
		if err := sink.receivePackets(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
}

type cpuSink struct {
	client       fwdpb.ForwardingClient
	ethDevToPort map[string]string
}

func (cpuSink *cpuSink) receivePackets(ctx context.Context) error {
	subClient, err := cpuSink.client.PacketSinkSubscribe(ctx, &fwdpb.PacketSinkRequest{ContextId: &fwdpb.ContextId{Id: contextID}})
	if err != nil {
		return err
	}
	ports := map[string]*kernel.GenetlinkPort{}

	for {
		msg, err := subClient.Recv()
		if err != nil {
			return err
		}
		switch resp := msg.Resp.(type) {
		case *fwdpb.PacketSinkResponse_Port:
			switch desc := resp.Port.Port; desc.PortType {
			case fwdpb.PortType_PORT_TYPE_GENETLINK:
				portDesc := desc.GetGenetlink()
				p, err := kernel.NewGenetlinkPort(portDesc.FamilyName, portDesc.GroupName)
				if err != nil {
					log.Errorf("failed to create port: %v", err)
				}
				ports[resp.Port.Port.PortId.ObjectId.Id] = p
			case fwdpb.PortType_PORT_TYPE_KERNEL, fwdpb.PortType_PORT_TYPE_TAP:
				name := desc.GetKernel().GetDeviceName()
				if name == "" {
					name = desc.GetTap().GetDeviceName()
				}
				cpuSink.ethDevToPort[name] = desc.PortId.ObjectId.Id
			}

		case *fwdpb.PacketSinkResponse_Packet:
			p, ok := ports[resp.Packet.Egress.ObjectId.Id]
			if !ok {
				log.Infof("skipping port with id %v", resp.Packet.Egress.ObjectId.Id)
			}
			if err := p.Write(resp.Packet.Bytes); err != nil {
				log.Warningf("failed to write packet: %v", err)
			}
		}
	}
}

func (cpuSink cpuSink) handleIPUpdates(ctx context.Context, client fwdpb.ForwardingClient) error {
	updCh := make(chan netlink.AddrUpdate)
	doneCh := make(chan struct{})

	ipToDevName := map[string]string{}

	go func() {
		for {
			select {
			case upd := <-updCh:
				l, err := netlink.LinkByIndex(upd.LinkIndex)
				if err != nil {
					log.Warningf("failed to get link: %v", err)
					continue
				}
				ip := upd.LinkAddress.IP.To4()
				if ip == nil {
					ip = upd.LinkAddress.IP.To16()
				}
				entry := fwdconfig.EntryDesc(fwdconfig.ExtactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(ip)))

				if upd.NewAddr {
					log.Infof("added new ip %s to device %s", upd.LinkAddress.IP.String(), l.Attrs().Name)
					ipToDevName[upd.LinkAddress.IP.String()] = l.Attrs().Name

					_, err := client.TableEntryAdd(ctx, fwdconfig.TableEntryAddRequest(contextID, ip2meTable).
						AppendEntry(entry, fwdconfig.Action(fwdconfig.TransmitAction(cpuSink.ethDevToPort[l.Attrs().Name]))).Build())
					if err != nil {
						log.Warningf("failed to get link: %v", err)
						continue
					}
				} else {
					_, err := client.TableEntryRemove(ctx, &fwdpb.TableEntryRemoveRequest{
						ContextId: &fwdpb.ContextId{Id: contextID},
						TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: ip2meTable}},
						EntryDesc: entry.Build(),
					})
					if err != nil {
						log.Warningf("failed to get link: %v", err)
						continue
					}
					delete(ipToDevName, upd.LinkAddress.IP.String())
				}
			case <-ctx.Done():
				close(doneCh)
				return
			}
		}
	}()
	netlink.AddrSubscribe(updCh, doneCh)
	return nil
}

func main() {
}
