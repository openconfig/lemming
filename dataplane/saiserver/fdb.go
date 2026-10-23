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
	"fmt"
	"log/slog"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdconfig"
	fwdbridge "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/bridge"
	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const fdbTable = "fdb"

type fdb struct {
	saipb.UnimplementedFdbServer
	mgr         *attrmgr.AttrMgr
	dataplane   switchDataplaneAPI
	mu          sync.RWMutex
	subscribers map[chan *saipb.FdbEventNotificationResponse]struct{}
}

func newFdb(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) (*fdb, error) {
	f := &fdb{
		mgr:         mgr,
		dataplane:   dataplane,
		subscribers: make(map[chan *saipb.FdbEventNotificationResponse]struct{}),
	}
	saipb.RegisterFdbServer(s, f)
	return f, nil
}

func (f *fdb) subscribe(ch chan *saipb.FdbEventNotificationResponse) func() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.subscribers[ch] = struct{}{}
	return func() {
		f.mu.Lock()
		defer f.mu.Unlock()
		delete(f.subscribers, ch)
	}
}

func (f *fdb) sendNotification(data *saipb.FdbEventNotificationData) {
	if data == nil || data.GetFdbEntry() == nil {
		return
	}
	resp := &saipb.FdbEventNotificationResponse{Data: []*saipb.FdbEventNotificationData{data}}
	f.mu.RLock()
	defer f.mu.RUnlock()
	for ch := range f.subscribers {
		select {
		case ch <- resp:
		default:
			slog.Warn("fdb notification channel full, dropping event")
		}
	}
}

func (f *fdb) FlushFdbEntries(ctx context.Context, req *saipb.FlushFdbEntriesRequest) (*saipb.FlushFdbEntriesResponse, error) {
	return &saipb.FlushFdbEntriesResponse{}, nil
}

func (f *fdb) CreateFdbEntry(ctx context.Context, req *saipb.CreateFdbEntryRequest) (*saipb.CreateFdbEntryResponse, error) {
	entry := req.GetEntry()
	if entry == nil || len(entry.GetMacAddress()) == 0 || req.GetBridgePortId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "FDB entry, MAC address, and bridge port ID are required")
	}
	mac, portOID := entry.GetMacAddress(), req.GetBridgePortId()
	bpReq := &saipb.GetBridgePortAttributeRequest{Oid: portOID, AttrType: []saipb.BridgePortAttr{saipb.BridgePortAttr_BRIDGE_PORT_ATTR_PORT_ID}}
	bpResp := &saipb.GetBridgePortAttributeResponse{}
	if err := f.mgr.PopulateAttributes(bpReq, bpResp); err != nil || bpResp.GetAttr().GetPortId() == 0 {
		return nil, fmt.Errorf("failed to resolve port ID for bridge port %d: %v", portOID, err)
	}
	portID := bpResp.GetAttr().GetPortId()
	addReq := fwdconfig.TableEntryAddRequest(f.dataplane.ID(), fdbTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).WithBytes(mac))),
		fwdconfig.TransmitAction(fmt.Sprint(portID)),
	).Build()
	if _, err := f.dataplane.TableEntryAdd(ctx, addReq); err != nil {
		return nil, fmt.Errorf("failed to add FDB entry to dataplane: %v", err)
	}
	entryType := req.GetType()
	if entryType == saipb.FdbEntryType_FDB_ENTRY_TYPE_UNSPECIFIED {
		entryType = saipb.FdbEntryType_FDB_ENTRY_TYPE_STATIC
	}
	f.sendNotification(&saipb.FdbEventNotificationData{
		EventType: saipb.FdbEvent_FDB_EVENT_LEARNED,
		FdbEntry:  &saipb.FdbEntry{SwitchId: entry.GetSwitchId(), MacAddress: mac, BvId: entry.GetBvId()},
		Attrs:     []*saipb.FdbEntryAttribute{{BridgePortId: proto.Uint64(portOID)}, {Type: entryType.Enum()}, {PacketAction: req.GetPacketAction().Enum()}},
	})
	return &saipb.CreateFdbEntryResponse{}, nil
}

func (f *fdb) RemoveFdbEntry(ctx context.Context, req *saipb.RemoveFdbEntryRequest) (*saipb.RemoveFdbEntryResponse, error) {
	entry := req.GetEntry()
	if entry == nil || len(entry.GetMacAddress()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "FDB entry and MAC address are required")
	}
	mac := entry.GetMacAddress()
	if fwdCtx, err := f.dataplane.FindContext(&fwdpb.ContextId{Id: f.dataplane.ID()}); err == nil && fwdCtx != nil {
		obj, _ := fwdCtx.Objects.FindID(&fwdpb.ObjectId{Id: fdbTable})
		if brTable, ok := obj.(*fwdbridge.Table); ok {
			_ = brTable.Remove(mac)
		}
	}
	delReq := fwdconfig.TableEntryRemoveRequest(f.dataplane.ID(), fdbTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).WithBytes(mac))),
	).Build()
	if _, err := f.dataplane.TableEntryRemove(ctx, delReq); err != nil {
		slog.WarnContext(ctx, "failed to remove FDB entry from dataplane", "err", err)
	}
	f.sendNotification(&saipb.FdbEventNotificationData{EventType: saipb.FdbEvent_FDB_EVENT_AGED, FdbEntry: &saipb.FdbEntry{SwitchId: entry.GetSwitchId(), MacAddress: mac, BvId: entry.GetBvId()}})
	return &saipb.RemoveFdbEntryResponse{}, nil
}

func (f *fdb) CreateFdbEntries(ctx context.Context, req *saipb.CreateFdbEntriesRequest) (*saipb.CreateFdbEntriesResponse, error) {
	resp := &saipb.CreateFdbEntriesResponse{}
	for _, r := range req.GetReqs() {
		entryResp, err := f.CreateFdbEntry(ctx, r)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, entryResp)
	}
	return resp, nil
}

func (f *fdb) RemoveFdbEntries(ctx context.Context, req *saipb.RemoveFdbEntriesRequest) (*saipb.RemoveFdbEntriesResponse, error) {
	resp := &saipb.RemoveFdbEntriesResponse{}
	for _, r := range req.GetReqs() {
		entryResp, err := f.RemoveFdbEntry(ctx, r)
		if err != nil {
			return nil, err
		}
		resp.Resps = append(resp.Resps, entryResp)
	}
	return resp, nil
}
