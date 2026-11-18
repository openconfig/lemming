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

package saiserver

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"slices"
	"strings"
	"sync"

	"github.com/openconfig/gnmi/errlist"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol/vxlan"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	defaultEncapTTL = 255 // Default of SAI_TUNNEL_ATTR_ENCAP_TTL_VAL.
	maxVXLANVNI     = 0xFFFFFF
	ipVersion4      = 4
	ipVersion6      = 6
	ipProtoIP4InIP  = 4
	ipProtoIP6InIP  = 41
	ipProtoGRE      = 47

	// Lucius FlowEntryDesc matches numerically smaller Priority values first, whereas SAI
	// SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_PRIORITY treats higher values as higher precedence.
	maxTermSAIPriority = (math.MaxUint32 - 4) / 4
	termRankP2P        = 0
	termRankP2MP       = 1
	termRankMP2P       = 1
	termRankMP2MP      = 2
)

type tunnel struct {
	saipb.UnimplementedTunnelServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI

	mu         sync.Mutex
	tunnelMaps map[uint64]saipb.TunnelMapType
	mapEntries map[uint64]map[uint64]*saipb.CreateTunnelMapEntryRequest // mapOID -> entryOID -> req
	entryToMap map[uint64]uint64                                         // entryOID -> mapOID
	mapTunnels map[uint64]map[uint64]struct{}                            // mapOID -> tunnelOID -> struct{}{}
	termKeys   map[string]uint64                                         // canonical flow key -> termOID
}

func vniToBytes(vni uint32) []byte {
	return []byte{byte(vni >> 16), byte(vni >> 8), byte(vni)}
}

func newTunnel(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *tunnel {
	t := &tunnel{
		mgr:        mgr,
		dataplane:  dataplane,
		tunnelMaps: make(map[uint64]saipb.TunnelMapType),
		mapEntries: make(map[uint64]map[uint64]*saipb.CreateTunnelMapEntryRequest),
		entryToMap: make(map[uint64]uint64),
		mapTunnels: make(map[uint64]map[uint64]struct{}),
		termKeys:   make(map[string]uint64),
	}
	saipb.RegisterTunnelServer(s, t)
	return t
}

// CreateTunnel creates a new tunnel supporting IPinIP, GRE, and VXLAN tunnel types.
func (t *tunnel) CreateTunnel(ctx context.Context, req *saipb.CreateTunnelRequest) (*saipb.CreateTunnelResponse, error) {
	tunType := req.GetType()

	supportedTypes := []saipb.TunnelType{
		saipb.TunnelType_TUNNEL_TYPE_IPINIP,
		saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE,
		saipb.TunnelType_TUNNEL_TYPE_VXLAN,
	}
	if !slices.Contains(supportedTypes, tunType) {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported tunnel type: %v", tunType)
	}

	// TODO: Support SAI_TUNNEL_ATTR_DECAP_TTL_MODE, DECAP_DSCP_MODE, and RFC 6040 DECAP_ECN_MODE.

	actions := []*fwdpb.ActionDesc{}

	// TODO: Support parsing QOS into ECN and DSCP bits separately.
	ecnMode := req.GetEncapEcnMode()
	dscpMode := req.GetEncapDscpMode()
	if ecnMode == saipb.TunnelEncapEcnMode_TUNNEL_ENCAP_ECN_MODE_STANDARD && dscpMode == saipb.TunnelDscpMode_TUNNEL_DSCP_MODE_UNIFORM_MODEL { // Copy the QOS bits from the inner IP header.
		actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_COPY, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS).
			WithFieldSrc(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_QOS).WithFieldSrcInstance(1)).Build())
	}

	ttlMode := req.GetEncapTtlMode()
	switch ttlMode {
	case saipb.TunnelTtlMode_TUNNEL_TTL_MODE_UNSPECIFIED, saipb.TunnelTtlMode_TUNNEL_TTL_MODE_UNIFORM_MODEL:
		// Copy the TTL from the inner IP header.
		actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_COPY, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP).
			WithFieldSrc(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP).WithFieldSrcInstance(1)).Build())
	case saipb.TunnelTtlMode_TUNNEL_TTL_MODE_PIPE_MODEL:
		ttl := uint8(defaultEncapTTL)
		if req.EncapTtlVal != nil {
			ttl = uint8(req.GetEncapTtlVal())
		}
		actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_HOP).
			WithValue([]byte{ttl})).Build())
	}

	if req.GetEncapSrcIp() != nil {
		actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(
			fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithValue(req.EncapSrcIp)).Build())
	}
	if req.GetEncapDstIp() != nil {
		actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(
			fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithValue(req.EncapDstIp)).Build())
	}

	if tunType == saipb.TunnelType_TUNNEL_TYPE_VXLAN {
		actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST).WithUint64Value(vxlan.DefaultPort)).Build())
		if req.GetVxlanUdpSport() != 0 {
			actions = append(actions, fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_SRC).WithUint64Value(uint64(req.GetVxlanUdpSport()))).Build())
		}
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	uniqueMapIDs := slices.Clone(req.GetDecapMappers())
	slices.Sort(uniqueMapIDs)
	uniqueMapIDs = slices.Compact(uniqueMapIDs)

	for _, mapID := range uniqueMapIDs {
		if _, ok := t.tunnelMaps[mapID]; !ok {
			return nil, status.Errorf(codes.NotFound, "tunnel map %d not found", mapID)
		}
	}

	id := t.mgr.NextID()

	under := req.GetUnderlayInterface()
	entry := fwdconfig.TableEntryAddRequest(t.dataplane.ID(), TunnelEncap).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(id))),
		fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_OUTPUT_IFACE).WithUint64Value(under),
	).Build()
	entry.Entries[0].Actions = append(entry.Entries[0].Actions, actions...)
	if _, err := t.dataplane.TableEntryAdd(ctx, entry); err != nil {
		return nil, err
	}

	var addedEntries []struct {
		tunID uint64
		vni   uint32
	}
	for _, mapID := range uniqueMapIDs {
		if t.tunnelMaps[mapID] == saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID {
			entryOIDs := make([]uint64, 0, len(t.mapEntries[mapID]))
			for eOID := range t.mapEntries[mapID] {
				entryOIDs = append(entryOIDs, eOID)
			}
			slices.Sort(entryOIDs)
			for _, eOID := range entryOIDs {
				entryReq := t.mapEntries[mapID][eOID]
				vEntry := fwdconfig.TableEntryAddRequest(t.dataplane.ID(), vniToVrfTable).AppendEntry(
					fwdconfig.EntryDesc(fwdconfig.ExactEntry(
						fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(id),
						fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI).WithBytes(vniToBytes(entryReq.GetVniIdKey())),
					)),
					fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64Value(entryReq.GetVirtualRouterIdValue()),
				).Build()
				if _, err := t.dataplane.TableEntryAdd(ctx, vEntry); err != nil {
					// Roll back added vniToVrfTable entries in reverse order.
					for i := len(addedEntries) - 1; i >= 0; i-- {
						rollbackReq := &fwdpb.TableEntryRemoveRequest{
							ContextId: &fwdpb.ContextId{Id: t.dataplane.ID()},
							TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: vniToVrfTable}},
							EntryDesc: fwdconfig.EntryDesc(fwdconfig.ExactEntry(
								fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(addedEntries[i].tunID),
								fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI).WithBytes(vniToBytes(addedEntries[i].vni)),
							)).Build(),
						}
						_, _ = t.dataplane.TableEntryRemove(ctx, rollbackReq)
					}
					// Remove TunnelEncap entry.
					encapRemoveReq := &fwdpb.TableEntryRemoveRequest{
						ContextId: &fwdpb.ContextId{Id: t.dataplane.ID()},
						TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: TunnelEncap}},
						EntryDesc: fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(id))).Build(),
					}
					_, _ = t.dataplane.TableEntryRemove(ctx, encapRemoveReq)
					return nil, err
				}
				addedEntries = append(addedEntries, struct {
					tunID uint64
					vni   uint32
				}{tunID: id, vni: entryReq.GetVniIdKey()})
			}
		}
	}

	for _, mapID := range uniqueMapIDs {
		if t.mapTunnels[mapID] == nil {
			t.mapTunnels[mapID] = make(map[uint64]struct{})
		}
		t.mapTunnels[mapID][id] = struct{}{}
	}

	t.mgr.StoreAttributes(id, req)

	return &saipb.CreateTunnelResponse{
		Oid: id,
	}, nil
}

func (t *tunnel) RemoveTunnel(ctx context.Context, req *saipb.RemoveTunnelRequest) (*saipb.RemoveTunnelResponse, error) {
	var errs errlist.List

	rReq := &fwdpb.TableEntryRemoveRequest{
		ContextId: &fwdpb.ContextId{Id: t.dataplane.ID()},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: TunnelEncap}},
		EntryDesc: fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(req.GetOid()))).Build(),
	}
	if _, err := t.dataplane.TableEntryRemove(ctx, rReq); err != nil {
		errs.Add(err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	tunID := req.GetOid()
	mapIDs := make([]uint64, 0, len(t.mapTunnels))
	for mID := range t.mapTunnels {
		mapIDs = append(mapIDs, mID)
	}
	slices.Sort(mapIDs)

	for _, mapID := range mapIDs {
		tunnels := t.mapTunnels[mapID]
		if _, ok := tunnels[tunID]; ok {
			if t.tunnelMaps[mapID] == saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID {
				entryOIDs := make([]uint64, 0, len(t.mapEntries[mapID]))
				for eOID := range t.mapEntries[mapID] {
					entryOIDs = append(entryOIDs, eOID)
				}
				slices.Sort(entryOIDs)
				for _, eOID := range entryOIDs {
					entryReq := t.mapEntries[mapID][eOID]
					vReq := &fwdpb.TableEntryRemoveRequest{
						ContextId: &fwdpb.ContextId{Id: t.dataplane.ID()},
						TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: vniToVrfTable}},
						EntryDesc: fwdconfig.EntryDesc(fwdconfig.ExactEntry(
							fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(tunID),
							fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI).WithBytes(vniToBytes(entryReq.GetVniIdKey())),
						)).Build(),
					}
					if _, err := t.dataplane.TableEntryRemove(ctx, vReq); err != nil {
						errs.Add(err)
					}
				}
			}
			delete(t.mapTunnels[mapID], tunID)
		}
	}

	if err := errs.Err(); err != nil {
		return nil, err
	}
	return &saipb.RemoveTunnelResponse{}, nil
}

func (t *tunnel) CreateTunnels(ctx context.Context, req *saipb.CreateTunnelsRequest) (*saipb.CreateTunnelsResponse, error) {
	resp := &saipb.CreateTunnelsResponse{}
	for _, r := range req.GetReqs() {
		res, err := t.CreateTunnel(ctx, r)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, res)
	}
	return resp, nil
}

func (t *tunnel) RemoveTunnels(ctx context.Context, req *saipb.RemoveTunnelsRequest) (*saipb.RemoveTunnelsResponse, error) {
	resp := &saipb.RemoveTunnelsResponse{}
	for _, r := range req.GetReqs() {
		res, err := t.RemoveTunnel(ctx, r)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, res)
	}
	return resp, nil
}

var (
	ipV4ExactMask = []byte{0xFF, 0xFF, 0xFF, 0xFF}
	ipV4AnyMask   = make([]byte, 4)
	ipV6ExactMask = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	ipV6AnyMask   = make([]byte, 16)
)

func maskBytes(ip, mask []byte) []byte {
	if len(ip) != len(mask) {
		return ip
	}
	out := make([]byte, len(ip))
	for i := 0; i < len(ip); i++ {
		out[i] = ip[i] & mask[i]
	}
	return out
}

func normalizeIPAddr(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return nil, nil
	}
	if len(b) == 4 {
		return slices.Clone(b), nil
	}
	if len(b) == 16 {
		if v4 := net.IP(b).To4(); v4 != nil {
			return slices.Clone(v4), nil
		}
		return slices.Clone(b), nil
	}
	return nil, status.Errorf(codes.InvalidArgument, "invalid IP address length %d", len(b))
}

func normalizeIPMask(mask []byte, addrLen int) ([]byte, error) {
	if len(mask) == 0 {
		if addrLen == 4 {
			return slices.Clone(ipV4AnyMask), nil
		}
		return slices.Clone(ipV6AnyMask), nil
	}
	if addrLen == 16 && len(mask) == 4 && slices.Equal(mask, ipV4AnyMask) {
		return slices.Clone(ipV6AnyMask), nil
	}
	if len(mask) != addrLen {
		return nil, status.Errorf(codes.InvalidArgument, "invalid IP mask length %d for address length %d", len(mask), addrLen)
	}
	return slices.Clone(mask), nil
}

func termFieldSetsFromReq(req *saipb.CreateTunnelTermTableEntryRequest) ([][]*fwdpb.PacketFieldMaskedBytes, fwdpb.PacketHeaderId, error) {
	srcIP, err := normalizeIPAddr(req.GetSrcIp())
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, err
	}
	dstIP, err := normalizeIPAddr(req.GetDstIp())
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, err
	}

	switch req.GetType() {
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P:
		if len(srcIP) == 0 || len(dstIP) == 0 {
			return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "P2P tunnel term entry requires both SrcIp and DstIp")
		}
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP:
		if len(dstIP) == 0 {
			return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "P2MP tunnel term entry requires DstIp")
		}
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2P:
		if len(srcIP) == 0 {
			return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "MP2P tunnel term entry requires SrcIp")
		}
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2MP:
		if len(srcIP) == 0 && len(dstIP) == 0 && len(req.GetSrcIpMask()) == 0 && len(req.GetDstIpMask()) == 0 {
			return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "MP2MP tunnel term entry requires at least one IP or mask")
		}
	default:
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "invalid tunnel term table entry type: %v", req.GetType())
	}

	var addrLen int
	if len(srcIP) != 0 && len(dstIP) != 0 {
		if len(srcIP) != len(dstIP) {
			return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "mismatched IP address families: src len %d, dst len %d", len(srcIP), len(dstIP))
		}
		addrLen = len(srcIP)
	} else if len(srcIP) != 0 {
		addrLen = len(srcIP)
	} else if len(dstIP) != 0 {
		addrLen = len(dstIP)
	} else {
		srcMaskLen := len(req.GetSrcIpMask())
		dstMaskLen := len(req.GetDstIpMask())
		if srcMaskLen != 0 && dstMaskLen != 0 {
			if srcMaskLen != dstMaskLen || (srcMaskLen != 4 && srcMaskLen != 16) {
				return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "invalid or mismatched IP mask lengths: src %d, dst %d", srcMaskLen, dstMaskLen)
			}
			addrLen = srcMaskLen
		} else if srcMaskLen == 4 || srcMaskLen == 16 {
			addrLen = srcMaskLen
		} else if dstMaskLen == 4 || dstMaskLen == 16 {
			addrLen = dstMaskLen
		} else {
			return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "cannot determine IP address family from empty IPs and masks")
		}
	}

	srcMask, err := normalizeIPMask(req.GetSrcIpMask(), addrLen)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, err
	}
	dstMask, err := normalizeIPMask(req.GetDstIpMask(), addrLen)
	if err != nil {
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, err
	}

	var zeroIP, exactMask []byte
	var headerID fwdpb.PacketHeaderId
	var ipVer byte
	if addrLen == 4 {
		zeroIP = slices.Clone(ipV4AnyMask)
		exactMask = slices.Clone(ipV4ExactMask)
		headerID = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4
		ipVer = ipVersion4
	} else {
		zeroIP = slices.Clone(ipV6AnyMask)
		exactMask = slices.Clone(ipV6ExactMask)
		headerID = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6
		ipVer = ipVersion6
	}

	if len(srcIP) == 0 {
		srcIP = slices.Clone(zeroIP)
	}
	if len(dstIP) == 0 {
		dstIP = slices.Clone(zeroIP)
	}

	var ipFields []*fwdpb.PacketFieldMaskedBytes
	switch req.GetType() {
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P:
		ipFields = []*fwdpb.PacketFieldMaskedBytes{
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(maskBytes(srcIP, exactMask), slices.Clone(exactMask)).Build(),
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(maskBytes(dstIP, exactMask), slices.Clone(exactMask)).Build(),
		}
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP:
		ipFields = []*fwdpb.PacketFieldMaskedBytes{
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(maskBytes(srcIP, srcMask), srcMask).Build(),
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(maskBytes(dstIP, exactMask), slices.Clone(exactMask)).Build(),
		}
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2P:
		ipFields = []*fwdpb.PacketFieldMaskedBytes{
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(maskBytes(srcIP, exactMask), slices.Clone(exactMask)).Build(),
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(maskBytes(dstIP, dstMask), dstMask).Build(),
		}
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2MP:
		ipFields = []*fwdpb.PacketFieldMaskedBytes{
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_SRC).WithBytes(maskBytes(srcIP, srcMask), srcMask).Build(),
			fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_ADDR_DST).WithBytes(maskBytes(dstIP, dstMask), dstMask).Build(),
		}
	}

	verField := fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_VERSION).
		WithBytes([]byte{ipVer}, []byte{0xFF}).Build()

	var fieldSets [][]*fwdpb.PacketFieldMaskedBytes
	switch req.GetTunnelType() {
	case saipb.TunnelType_TUNNEL_TYPE_IPINIP:
		// Return 2 field sets: one for IPv4-in-IP (proto 4) and one for IPv6-in-IP (proto 41).
		proto4 := fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).
			WithBytes([]byte{ipProtoIP4InIP}, []byte{0xFF}).Build()
		set1 := append([]*fwdpb.PacketFieldMaskedBytes{verField, proto4}, ipFields...)

		proto41 := fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).
			WithBytes([]byte{ipProtoIP6InIP}, []byte{0xFF}).Build()
		set2 := append([]*fwdpb.PacketFieldMaskedBytes{verField, proto41}, ipFields...)

		fieldSets = [][]*fwdpb.PacketFieldMaskedBytes{set1, set2}

	case saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE:
		protoGRE := fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).
			WithBytes([]byte{ipProtoGRE}, []byte{0xFF}).Build()
		set := append([]*fwdpb.PacketFieldMaskedBytes{verField, protoGRE}, ipFields...)
		fieldSets = [][]*fwdpb.PacketFieldMaskedBytes{set}

	case saipb.TunnelType_TUNNEL_TYPE_VXLAN:
		protoUDP := fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_IP_PROTO).
			WithBytes([]byte{ipProtoUDP}, []byte{0xFF}).Build()
		dstPort := fwdconfig.PacketFieldMaskedBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_L4_PORT_DST).
			WithBytes(binary.BigEndian.AppendUint16(nil, vxlan.DefaultPort), []byte{0xFF, 0xFF}).Build()
		set := append([]*fwdpb.PacketFieldMaskedBytes{verField, protoUDP, dstPort}, ipFields...)
		fieldSets = [][]*fwdpb.PacketFieldMaskedBytes{set}

	default:
		return nil, fwdpb.PacketHeaderId_PACKET_HEADER_ID_UNSPECIFIED, status.Errorf(codes.InvalidArgument, "invalid tunnel type: %v", req.GetTunnelType())
	}

	return fieldSets, headerID, nil
}

func termEntryPriority(req *saipb.CreateTunnelTermTableEntryRequest) (uint32, error) {
	var rank uint32
	switch req.GetType() {
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P:
		rank = termRankP2P
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP:
		rank = termRankP2MP
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2P:
		rank = termRankMP2P
	case saipb.TunnelTermTableEntryType_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2MP:
		rank = termRankMP2MP
	default:
		return 0, status.Errorf(codes.InvalidArgument, "invalid tunnel term table entry type: %v", req.GetType())
	}
	saiPrio := min(req.GetPriority(), maxTermSAIPriority)
	return (math.MaxUint32 - 3) - saiPrio*4 + rank, nil
}

func canonicalFlowKey(priority uint32, fields []*fwdpb.PacketFieldMaskedBytes) string {
	var b strings.Builder
	fmt.Fprintf(&b, "prio:%d;", priority)
	for _, f := range fields {
		fmt.Fprintf(&b, "f:%d:%x/%x;", f.GetFieldId().GetField().GetFieldNum(), f.GetBytes(), f.GetMasks())
	}
	return b.String()
}

// overlayVRFForTunnel resolves the post-decap VRF from the tunnel's SAI_TUNNEL_ATTR_OVERLAY_INTERFACE.
// Returns (vrf, true, nil) when the tunnel has an overlay router interface configured,
// (0, false, nil) when actionTunnelID == 0 or OverlayInterface is unset, or InvalidArgument
// if a non-zero actionTunnelID or OverlayInterface does not reference a live SAI object.
func (t *tunnel) overlayVRFForTunnel(actionTunnelID uint64) (uint64, bool, error) {
	if actionTunnelID == 0 {
		return 0, false, nil
	}
	tunID := fmt.Sprint(actionTunnelID)
	if t.mgr.GetType(tunID) != saipb.ObjectType_OBJECT_TYPE_TUNNEL ||
		t.mgr.GetAttribute(tunID, int32(saipb.TunnelAttr_TUNNEL_ATTR_TYPE)) == nil {
		return 0, false, status.Errorf(codes.InvalidArgument, "action tunnel %d not found", actionTunnelID)
	}
	tunReq := &saipb.CreateTunnelRequest{}
	if err := t.mgr.PopulateAllAttributes(tunID, tunReq); err != nil {
		return 0, false, status.Errorf(codes.Internal, "failed to populate tunnel %d attributes: %v", actionTunnelID, err)
	}
	if tunReq.OverlayInterface == nil {
		return 0, false, nil
	}
	rifID := fmt.Sprint(tunReq.GetOverlayInterface())
	if t.mgr.GetType(rifID) != saipb.ObjectType_OBJECT_TYPE_ROUTER_INTERFACE ||
		t.mgr.GetAttribute(rifID, int32(saipb.RouterInterfaceAttr_ROUTER_INTERFACE_ATTR_TYPE)) == nil {
		return 0, false, status.Errorf(codes.InvalidArgument, "overlay interface %d not found", tunReq.GetOverlayInterface())
	}
	rifReq := &saipb.CreateRouterInterfaceRequest{}
	if err := t.mgr.PopulateAllAttributes(rifID, rifReq); err != nil {
		return 0, false, status.Errorf(codes.Internal, "failed to populate router interface %d attributes: %v", tunReq.GetOverlayInterface(), err)
	}
	return rifReq.GetVirtualRouterId(), true, nil
}

// CreateTunnelTermTableEntry creates a tunnel termination table entry supporting IPinIP, GRE, and VXLAN tunnel types.
func (t *tunnel) CreateTunnelTermTableEntry(ctx context.Context, req *saipb.CreateTunnelTermTableEntryRequest) (*saipb.CreateTunnelTermTableEntryResponse, error) {
	if req.GetTunnelType() == saipb.TunnelType_TUNNEL_TYPE_VXLAN && req.GetActionTunnelId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "action_tunnel_id is required for VXLAN tunnel termination")
	}

	fieldSets, headerID, err := termFieldSetsFromReq(req)
	if err != nil {
		return nil, err
	}

	prio, err := termEntryPriority(req)
	if err != nil {
		return nil, err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	for _, fs := range fieldSets {
		k := canonicalFlowKey(prio, fs)
		if _, exists := t.termKeys[k]; exists {
			return nil, status.Errorf(codes.AlreadyExists, "duplicate tunnel term table entry")
		}
	}

	var (
		overlayVRF    uint64
		hasOverlayVRF bool
	)
	if req.GetTunnelType() == saipb.TunnelType_TUNNEL_TYPE_IPINIP || req.GetTunnelType() == saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE {
		overlayVRF, hasOverlayVRF, err = t.overlayVRFForTunnel(req.GetActionTunnelId())
		if err != nil {
			return nil, err
		}
	}

	id := t.mgr.NextID()

	var actions []*fwdpb.ActionDesc
	if req.GetActionTunnelId() != 0 {
		actions = append(actions,
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64Value(req.GetActionTunnelId())).Build(),
		)
	}

	// TODO: Match outer packet PACKET_FIELD_NUM_PACKET_VRF against req.VrId in tunTermTable.
	if hasOverlayVRF {
		actions = append(actions,
			fwdconfig.Action(fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64Value(overlayVRF)).Build(),
		)
	}

	switch req.GetTunnelType() {
	case saipb.TunnelType_TUNNEL_TYPE_IPINIP:
		actions = append(actions, &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
			Action: &fwdpb.ActionDesc_Decap{
				Decap: &fwdpb.DecapActionDesc{
					HeaderId: headerID,
				},
			},
		})
	case saipb.TunnelType_TUNNEL_TYPE_IPINIP_GRE:
		actions = append(actions, &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
			Action: &fwdpb.ActionDesc_Decap{
				Decap: &fwdpb.DecapActionDesc{
					HeaderId: headerID,
				},
			},
		}, &fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
			Action: &fwdpb.ActionDesc_Decap{
				Decap: &fwdpb.DecapActionDesc{
					HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_GRE,
				},
			},
		})
	case saipb.TunnelType_TUNNEL_TYPE_VXLAN:
		// Lookup VNI-to-VRF map before stripping the VXLAN header so PACKET_FIELD_NUM_VXLAN_VNI is still present.
		// Note: PACKET_VRF is set by the vniToVrfTable lookup; req.VrId is NOT applied directly.
		actions = append(actions, fwdconfig.Action(fwdconfig.LookupAction(vniToVrfTable)).Build())

		// Strip outer Ethernet + IP + UDP + VXLAN so the inner Ethernet frame becomes the outermost L2 header.
		actions = append(actions,
			&fwdpb.ActionDesc{
				ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
				Action: &fwdpb.ActionDesc_Decap{
					Decap: &fwdpb.DecapActionDesc{
						HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
					},
				},
			},
			&fwdpb.ActionDesc{
				ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
				Action: &fwdpb.ActionDesc_Decap{
					Decap: &fwdpb.DecapActionDesc{
						HeaderId: headerID,
					},
				},
			},
			&fwdpb.ActionDesc{
				ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
				Action: &fwdpb.ActionDesc_Decap{
					Decap: &fwdpb.DecapActionDesc{
						HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_UDP,
					},
				},
			},
			&fwdpb.ActionDesc{
				ActionType: fwdpb.ActionType_ACTION_TYPE_DECAP,
				Action: &fwdpb.ActionDesc_Decap{
					Decap: &fwdpb.DecapActionDesc{
						HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_VXLAN,
					},
				},
			},
		)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid tunnel type: %v", req.GetTunnelType())
	}

	var addedReqs []*fwdpb.TableEntryRemoveRequest
	for _, fs := range fieldSets {
		tReq := &fwdpb.TableEntryAddRequest{
			ContextId: &fwdpb.ContextId{Id: t.dataplane.ID()},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
			EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
				Flow: &fwdpb.FlowEntryDesc{
					Id:       uint32(id),
					Priority: prio,
					Bank:     0,
					Fields:   fs,
				},
			}},
			Actions: actions,
		}
		if _, err := t.dataplane.TableEntryAdd(ctx, tReq); err != nil {
			for i := len(addedReqs) - 1; i >= 0; i-- {
				_, _ = t.dataplane.TableEntryRemove(ctx, addedReqs[i])
			}
			return nil, err
		}
		addedReqs = append(addedReqs, &fwdpb.TableEntryRemoveRequest{
			ContextId: &fwdpb.ContextId{Id: t.dataplane.ID()},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
			EntryDesc: tReq.EntryDesc,
		})
	}

	for _, fs := range fieldSets {
		k := canonicalFlowKey(prio, fs)
		t.termKeys[k] = id
	}

	t.mgr.StoreAttributes(id, req)

	return &saipb.CreateTunnelTermTableEntryResponse{
		Oid: id,
	}, nil
}

func (t *tunnel) RemoveTunnelTermTableEntry(ctx context.Context, req *saipb.RemoveTunnelTermTableEntryRequest) (*saipb.RemoveTunnelTermTableEntryResponse, error) {
	cReq := &saipb.CreateTunnelTermTableEntryRequest{}
	if err := t.mgr.PopulateAllAttributes(fmt.Sprint(req.GetOid()), cReq); err != nil {
		return nil, err
	}
	prio, err := termEntryPriority(cReq)
	if err != nil {
		return nil, err
	}
	fieldSets, _, err := termFieldSetsFromReq(cReq)
	if err != nil {
		return nil, err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	var errs errlist.List
	for _, fs := range fieldSets {
		tReq := &fwdpb.TableEntryRemoveRequest{
			ContextId: &fwdpb.ContextId{Id: t.dataplane.ID()},
			TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: tunTermTable}},
			EntryDesc: &fwdpb.EntryDesc{Entry: &fwdpb.EntryDesc_Flow{
				Flow: &fwdpb.FlowEntryDesc{
					Id:       uint32(req.GetOid()),
					Priority: prio,
					Bank:     0,
					Fields:   fs,
				},
			}},
		}
		if _, err := t.dataplane.TableEntryRemove(ctx, tReq); err != nil {
			errs.Add(err)
		}
		k := canonicalFlowKey(prio, fs)
		delete(t.termKeys, k)
	}

	if err := errs.Err(); err != nil {
		return nil, err
	}
	return &saipb.RemoveTunnelTermTableEntryResponse{}, nil
}

// CreateTunnelMap creates a tunnel map.
// Note: TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID is programmed into vniToVrfTable when
// entries are added or tunnels are bound. L2 VNI map types (e.g. VNI_TO_VLAN_ID,
// VNI_TO_BRIDGE_IF) are stored for SAI/orchagent compatibility (L2 VNI dataplane bridging
// is not yet supported).
func (t *tunnel) CreateTunnelMap(ctx context.Context, req *saipb.CreateTunnelMapRequest) (*saipb.CreateTunnelMapResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	id := t.mgr.NextID()
	t.tunnelMaps[id] = req.GetType()
	t.mapEntries[id] = make(map[uint64]*saipb.CreateTunnelMapEntryRequest)
	t.mapTunnels[id] = make(map[uint64]struct{})
	t.mgr.StoreAttributes(id, req)
	return &saipb.CreateTunnelMapResponse{Oid: id}, nil
}

func (t *tunnel) RemoveTunnelMap(ctx context.Context, req *saipb.RemoveTunnelMapRequest) (*saipb.RemoveTunnelMapResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	id := req.GetOid()
	if _, ok := t.tunnelMaps[id]; !ok {
		return nil, status.Errorf(codes.NotFound, "tunnel map %d not found", id)
	}
	if len(t.mapTunnels[id]) > 0 || len(t.mapEntries[id]) > 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "tunnel map %d is in use", id)
	}

	delete(t.tunnelMaps, id)
	delete(t.mapEntries, id)
	delete(t.mapTunnels, id)
	return &saipb.RemoveTunnelMapResponse{}, nil
}

func (t *tunnel) CreateTunnelMapEntry(ctx context.Context, req *saipb.CreateTunnelMapEntryRequest) (*saipb.CreateTunnelMapEntryResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	mapID := req.GetTunnelMap()
	mapType, ok := t.tunnelMaps[mapID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "tunnel map %d not found", mapID)
	}
	if req.TunnelMapType != nil && req.GetTunnelMapType() != mapType {
		return nil, status.Errorf(codes.InvalidArgument, "tunnel map entry type %v does not match map type %v", req.GetTunnelMapType(), mapType)
	}

	if mapType == saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID {
		if req.GetVniIdKey() > maxVXLANVNI {
			return nil, status.Errorf(codes.InvalidArgument, "VNI %d exceeds 24-bit max (%d)", req.GetVniIdKey(), maxVXLANVNI)
		}
		for _, existing := range t.mapEntries[mapID] {
			if existing.GetVniIdKey() == req.GetVniIdKey() {
				return nil, status.Errorf(codes.AlreadyExists, "tunnel map entry with VNI %d already exists in map %d", req.GetVniIdKey(), mapID)
			}
		}

		tunIDs := make([]uint64, 0, len(t.mapTunnels[mapID]))
		for tunID := range t.mapTunnels[mapID] {
			tunIDs = append(tunIDs, tunID)
		}
		slices.Sort(tunIDs)
		var addedTuns []uint64
		for _, tunID := range tunIDs {
			entry := fwdconfig.TableEntryAddRequest(t.dataplane.ID(), vniToVrfTable).AppendEntry(
				fwdconfig.EntryDesc(fwdconfig.ExactEntry(
					fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(tunID),
					fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI).WithBytes(vniToBytes(req.GetVniIdKey())),
				)),
				fwdconfig.UpdateAction(fwdpb.UpdateType_UPDATE_TYPE_SET, fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_VRF).WithUint64Value(req.GetVirtualRouterIdValue()),
			).Build()
			if _, err := t.dataplane.TableEntryAdd(ctx, entry); err != nil {
				// Roll back added entries in reverse order.
				for i := len(addedTuns) - 1; i >= 0; i-- {
					vReq := &fwdpb.TableEntryRemoveRequest{
						ContextId: &fwdpb.ContextId{Id: t.dataplane.ID()},
						TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: vniToVrfTable}},
						EntryDesc: fwdconfig.EntryDesc(fwdconfig.ExactEntry(
							fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(addedTuns[i]),
							fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI).WithBytes(vniToBytes(req.GetVniIdKey())),
						)).Build(),
					}
					_, _ = t.dataplane.TableEntryRemove(ctx, vReq)
				}
				return nil, err
			}
			addedTuns = append(addedTuns, tunID)
		}
	}

	id := t.mgr.NextID()
	if t.mapEntries[mapID] == nil {
		t.mapEntries[mapID] = make(map[uint64]*saipb.CreateTunnelMapEntryRequest)
	}
	t.mapEntries[mapID][id] = proto.Clone(req).(*saipb.CreateTunnelMapEntryRequest)
	t.entryToMap[id] = mapID
	t.mgr.StoreAttributes(id, req)

	return &saipb.CreateTunnelMapEntryResponse{Oid: id}, nil
}

func (t *tunnel) RemoveTunnelMapEntry(ctx context.Context, req *saipb.RemoveTunnelMapEntryRequest) (*saipb.RemoveTunnelMapEntryResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	entryOID := req.GetOid()
	mapID, ok := t.entryToMap[entryOID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "tunnel map entry %d not found", entryOID)
	}

	entryReq := t.mapEntries[mapID][entryOID]
	var errs errlist.List
	if t.tunnelMaps[mapID] == saipb.TunnelMapType_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID {
		tunIDs := make([]uint64, 0, len(t.mapTunnels[mapID]))
		for tunID := range t.mapTunnels[mapID] {
			tunIDs = append(tunIDs, tunID)
		}
		slices.Sort(tunIDs)
		for _, tunID := range tunIDs {
			vReq := &fwdpb.TableEntryRemoveRequest{
				ContextId: &fwdpb.ContextId{Id: t.dataplane.ID()},
				TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: vniToVrfTable}},
				EntryDesc: fwdconfig.EntryDesc(fwdconfig.ExactEntry(
					fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_TUNNEL_ID).WithUint64(tunID),
					fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_VXLAN_VNI).WithBytes(vniToBytes(entryReq.GetVniIdKey())),
				)).Build(),
			}
			if _, err := t.dataplane.TableEntryRemove(ctx, vReq); err != nil {
				errs.Add(err)
			}
		}
	}
	delete(t.mapEntries[mapID], entryOID)
	delete(t.entryToMap, entryOID)
	if err := errs.Err(); err != nil {
		return nil, err
	}
	return &saipb.RemoveTunnelMapEntryResponse{}, nil
}
