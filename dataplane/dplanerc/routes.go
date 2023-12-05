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

package dplanerc

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/netip"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/openconfig/ygnmi/schemaless"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/gnmi"

	log "github.com/golang/glog"

	saipb "github.com/openconfig/lemming/dataplane/proto"
	"github.com/openconfig/lemming/dataplane/saiserver"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// RouteQuery returns a ygnmi query for a route with the given prefix and vrf.
func RouteQuery(ni string, prefix string) ygnmi.ConfigQuery[*dpb.Route] {
	q, err := schemaless.NewConfig[*dpb.Route](fmt.Sprintf("/dataplane/routes/route[prefix=%s][vrf=%s]", prefix, ni), gnmi.InternalOrigin)
	if err != nil {
		log.Fatal(err)
	}
	return q
}

// MustWildcardQuery returns a wildcard card query for all routes.
func MustWildcardQuery() ygnmi.WildcardQuery[*dpb.Route] {
	q, err := schemaless.NewWildcard[*dpb.Route]("/dataplane/routes/route[prefix=*][vrf=*]", gnmi.InternalOrigin)
	if err != nil {
		log.Fatal(err)
	}
	return q
}

func (ni *Reconciler) StartRoute(ctx context.Context, client *ygnmi.Client) error {
	ctx, cancelFn := context.WithCancel(ctx)
	w := ygnmi.WatchAll(ctx, client, MustWildcardQuery(), func(v *ygnmi.Value[*dpb.Route]) error {
		route, present := v.Val()
		prefix, err := netip.ParsePrefix(v.Path.Elem[2].Key["prefix"])
		if err != nil {
			log.Warningf("failed to parse cidr: %v", err)
			return ygnmi.Continue
		}
		ipBytes := prefix.Masked().Addr().AsSlice()
		mask := net.CIDRMask(prefix.Bits(), len(ipBytes)*8)

		entry := &saipb.RouteEntry{
			SwitchId: ni.switchID,
			VrId:     0, // TODO: support vrf-ids other than 0.
			Destination: &saipb.IpPrefix{
				Addr: ipBytes,
				Mask: mask,
			},
		}

		if !present {
			log.Infof("removing route: %v", prefix)
			_, err := ni.routeClient.RemoveRouteEntry(ctx, &saipb.RemoveRouteEntryRequest{
				Entry: entry,
			})
			if err != nil {
				log.Warningf("failed to delete route: %v", err)
			}
			return ygnmi.Continue
		}
		rReq := saipb.CreateRouteEntryRequest{
			Entry:        entry,
			PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD.Enum(),
		}

		if route.GetInterface() != nil { // If next hop is a interface.
			// TODO: Add support for subinterfaces.
			data := ni.ocInterfaceData[ocInterface{name: route.GetInterface().GetInterface(), subintf: route.GetInterface().GetSubinterface()}]
			rReq.NextHopId = &data.rifID

			if _, err := ni.routeClient.CreateRouteEntry(ctx, &rReq); err != nil {
				log.Warningf("failed to create route: %v", err)
			}
			log.Infof("added connected route: %v", &rReq)
			return ygnmi.Continue
		}
		var hopID uint64
		if len(route.GetNextHops().GetHops()) == 1 {
			hopID, err = ni.createNextHop(ctx, route.GetNextHops().Hops[0])
			if err != nil {
				log.Warningf("failed to create next hop: %v", err)
				return ygnmi.Continue
			}
		} else {
			group, err := ni.nextHopGroupClient.CreateNextHopGroup(ctx, &saipb.CreateNextHopGroupRequest{
				Switch: ni.switchID,
				Type:   saipb.NextHopGroupType_NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP.Enum(),
			})
			hopID = group.Oid
			if err != nil {
				log.Warningf("failed to create next hop group: %v", err)
				return ygnmi.Continue
			}
			for i, nh := range route.GetNextHops().GetHops() {
				hopID, err = ni.createNextHop(ctx, nh)
				if err != nil {
					log.Warningf("failed to create next hop: %v", err)
					return ygnmi.Continue
				}
				_, err = ni.nextHopGroupClient.CreateNextHopGroupMember(ctx, &saipb.CreateNextHopGroupMemberRequest{
					Switch:         ni.switchID,
					NextHopGroupId: &group.Oid,
					NextHopId:      &hopID,
					Weight:         proto.Uint32(uint32(route.GetNextHops().Weights[i])),
				})
				if err != nil {
					log.Warningf("failed to create next group member: %v", err)
					return ygnmi.Continue
				}
			}
		}
		rReq.NextHopId = proto.Uint64(hopID)
		if _, err := ni.routeClient.CreateRouteEntry(ctx, &rReq); err != nil {
			log.Warningf("failed to create route: %v", err)
			return ygnmi.Continue
		}
		log.Infof("created route entry: %v", &rReq)

		return ygnmi.Continue
	})
	go func() {
		// TODO: handle error
		if _, err := w.Await(); err != nil {
			log.Warningf("routes watch err: %v", err)
		}
	}()
	ni.closers = append(ni.closers, cancelFn)
	return nil
}

func (ni *Reconciler) createNextHop(ctx context.Context, hop *dpb.NextHop) (uint64, error) {
	ip, err := netip.ParseAddr(hop.GetNextHopIp())
	if err != nil {
		return 0, err
	}
	data := ni.ocInterfaceData[ocInterface{name: hop.GetInterface().GetInterface(), subintf: hop.GetInterface().GetSubinterface()}]
	hopReq := saipb.CreateNextHopRequest{
		Switch:            ni.switchID,
		Type:              saipb.NextHopType_NEXT_HOP_TYPE_IP.Enum(),
		Ip:                ip.AsSlice(),
		RouterInterfaceId: proto.Uint64(data.rifID),
	}
	resp, err := ni.nextHopClient.CreateNextHop(ctx, &hopReq)
	if err != nil {
		return 0, err
	}
	log.Infof("created next hop: %v", &hopReq)
	if hop.GetGue() != nil {
		acts, err := gueActions(hop.GetGue())
		if err != nil {
			return 0, err
		}
		actReq := &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: ni.contextID},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: saiserver.NHActionTable}},
			EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Exact{
				Exact: &fwdpb.ExactEntryDesc{
					Fields: []*fwdpb.PacketFieldBytes{{
						FieldId: &fwdpb.PacketFieldId{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_ID}},
						Bytes:   binary.BigEndian.AppendUint64(nil, resp.Oid),
					}},
				},
			}},
			Actions: acts,
		}
		// TODO: Ideally, this would use the SAI tunnel, but it's not currently supported.
		_, err = ni.fwdClient.TableEntryAdd(ctx, actReq)
		if err != nil {
			return 0, err
		}
		log.Infof("created gue actions: %v", actReq)
	}
	return resp.Oid, nil
}

func gueActions(gueHeaders *dpb.GUE) ([]*fwdpb.ActionDesc, error) {
	var ip gopacket.SerializableLayer
	var headerID fwdpb.PacketHeaderId
	if !gueHeaders.IsV6 {
		ip = &layers.IPv4{
			Version:  4,
			IHL:      5,
			Protocol: layers.IPProtocolUDP,
			SrcIP:    gueHeaders.SrcIp,
			DstIP:    gueHeaders.DstIp,
		}
		headerID = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4
	} else {
		ip = &layers.IPv6{
			Version:    6,
			NextHeader: layers.IPProtocolUDP,
			SrcIP:      gueHeaders.SrcIp,
			DstIP:      gueHeaders.DstIp,
		}
		headerID = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6
	}

	udp := &layers.UDP{
		SrcPort: 0,  // TODO(wenbli): Implement hashing for srcPort.
		Length:  34, // TODO(wenbli): Figure out how to not make this hardcoded.
	}
	udp.DstPort = layers.UDPPort(gueHeaders.DstPort)
	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{}, ip, udp); err != nil {
		return nil, fmt.Errorf("failed to serialize GUE headers: %v", err)
	}

	return []*fwdpb.ActionDesc{{
		ActionType: fwdpb.ActionType_ACTION_TYPE_REPARSE,
		Action: &fwdpb.ActionDesc_Reparse{
			Reparse: &fwdpb.ReparseActionDesc{
				HeaderId: headerID,
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
	}}, nil
}
