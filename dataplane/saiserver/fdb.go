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

type fdbEntryKey struct {
	bvID uint64
	mac  string
}

type fdbEntryRecord struct {
	mac          []byte
	bvID         uint64
	bridgePortID uint64
	entryType    saipb.FdbEntryType
	switchID     uint64
}

type fdb struct {
	saipb.UnimplementedFdbServer
	mgr         *attrmgr.AttrMgr
	dataplane   switchDataplaneAPI
	mu          sync.RWMutex
	subscribers map[chan *saipb.FdbEventNotificationResponse]struct{}
	entries     map[fdbEntryKey]*fdbEntryRecord
}

func newFdb(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) (*fdb, error) {
	f := &fdb{
		mgr:         mgr,
		dataplane:   dataplane,
		subscribers: make(map[chan *saipb.FdbEventNotificationResponse]struct{}),
		entries:     make(map[fdbEntryKey]*fdbEntryRecord),
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
	f.mu.Lock()
	if f.entries == nil {
		f.entries = make(map[fdbEntryKey]*fdbEntryRecord)
	}
	key := fdbEntryKey{
		bvID: data.GetFdbEntry().GetBvId(),
		mac:  string(data.GetFdbEntry().GetMacAddress()),
	}
	switch data.GetEventType() {
	case saipb.FdbEvent_FDB_EVENT_LEARNED:
		rec := &fdbEntryRecord{
			mac:       data.GetFdbEntry().GetMacAddress(),
			bvID:      data.GetFdbEntry().GetBvId(),
			switchID:  data.GetFdbEntry().GetSwitchId(),
			entryType: saipb.FdbEntryType_FDB_ENTRY_TYPE_DYNAMIC,
		}
		for _, attr := range data.GetAttrs() {
			if attr.BridgePortId != nil {
				rec.bridgePortID = attr.GetBridgePortId()
			}
			if attr.Type != nil {
				rec.entryType = attr.GetType()
			}
		}
		f.entries[key] = rec
	case saipb.FdbEvent_FDB_EVENT_AGED, saipb.FdbEvent_FDB_EVENT_FLUSHED:
		delete(f.entries, key)
	}
	f.mu.Unlock()

	f.broadcastNotification(data)
}

func (f *fdb) broadcastNotification(data *saipb.FdbEventNotificationData) {
	resp := &saipb.FdbEventNotificationResponse{
		Data: []*saipb.FdbEventNotificationData{data},
	}
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
	slog.InfoContext(ctx, "FlushFdbEntries called", "bridgePortId", req.GetBridgePortId(), "bvId", req.GetBvId(), "entryType", req.GetEntryType())

	var brTable *fwdbridge.Table
	if fwdCtx, err := f.dataplane.FindContext(&fwdpb.ContextId{Id: f.dataplane.ID()}); err == nil && fwdCtx != nil {
		if obj, err := fwdCtx.Objects.FindID(&fwdpb.ObjectId{Id: FDBTable}); err == nil {
			brTable, _ = obj.(*fwdbridge.Table)
		}
	}

	f.mu.Lock()
	var toFlush []*fdbEntryRecord
	for key, entry := range f.entries {
		if req.GetBridgePortId() != 0 && entry.bridgePortID != req.GetBridgePortId() {
			continue
		}
		if req.GetBvId() != 0 && entry.bvID != req.GetBvId() {
			continue
		}
		switch req.GetEntryType() {
		case saipb.FdbFlushEntryType_FDB_FLUSH_ENTRY_TYPE_DYNAMIC:
			if entry.entryType != saipb.FdbEntryType_FDB_ENTRY_TYPE_DYNAMIC {
				continue
			}
		case saipb.FdbFlushEntryType_FDB_FLUSH_ENTRY_TYPE_STATIC:
			if entry.entryType != saipb.FdbEntryType_FDB_ENTRY_TYPE_STATIC {
				continue
			}
		}
		toFlush = append(toFlush, entry)
		delete(f.entries, key)
	}
	f.mu.Unlock()

	for _, entry := range toFlush {
		if brTable != nil {
			_ = brTable.Remove(entry.mac)
		}

		delReq := fwdconfig.TableEntryRemoveRequest(f.dataplane.ID(), FDBTable).AppendEntry(
			fwdconfig.EntryDesc(fwdconfig.ExactEntry(
				fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).WithBytes(entry.mac),
			)),
		).Build()

		if _, err := f.dataplane.TableEntryRemove(ctx, delReq); err != nil {
			slog.WarnContext(ctx, "FlushFdbEntries: failed to remove entry from dataplane", "mac", fmt.Sprintf("%x", entry.mac), "err", err)
		}

		f.broadcastNotification(&saipb.FdbEventNotificationData{
			EventType: saipb.FdbEvent_FDB_EVENT_AGED,
			FdbEntry: &saipb.FdbEntry{
				SwitchId:   entry.switchID,
				MacAddress: entry.mac,
				BvId:       entry.bvID,
			},
			Attrs: []*saipb.FdbEntryAttribute{
				{
					BridgePortId: proto.Uint64(entry.bridgePortID),
				},
				{
					Type: entry.entryType.Enum(),
				},
			},
		})
	}

	// If flushing all entries (or all dynamic entries across the whole switch), also ensure bridge table is cleared.
	if req.GetBridgePortId() == 0 && req.GetBvId() == 0 && (req.GetEntryType() == saipb.FdbFlushEntryType_FDB_FLUSH_ENTRY_TYPE_ALL || req.GetEntryType() == saipb.FdbFlushEntryType_FDB_FLUSH_ENTRY_TYPE_DYNAMIC || req.GetEntryType() == saipb.FdbFlushEntryType_FDB_FLUSH_ENTRY_TYPE_UNSPECIFIED) {
		if brTable != nil {
			brTable.Clear()
		}
	}

	return &saipb.FlushFdbEntriesResponse{}, nil
}

func (f *fdb) CreateFdbEntry(ctx context.Context, req *saipb.CreateFdbEntryRequest) (*saipb.CreateFdbEntryResponse, error) {
	slog.InfoContext(ctx, "CreateFdbEntry called", "mac", req.GetEntry().GetMacAddress(), "vlan", req.GetEntry().GetBvId(), "bridge_port", req.GetBridgePortId())
	entry := req.GetEntry()
	if entry == nil {
		return nil, status.Errorf(codes.InvalidArgument, "FDB entry is required")
	}

	mac := entry.GetMacAddress()
	if len(mac) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "MAC address is required")
	}

	portOID := req.GetBridgePortId()
	if portOID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Bridge port ID is required")
	}

	bpReq := &saipb.GetBridgePortAttributeRequest{
		Oid:      portOID,
		AttrType: []saipb.BridgePortAttr{saipb.BridgePortAttr_BRIDGE_PORT_ATTR_PORT_ID},
	}
	bpResp := &saipb.GetBridgePortAttributeResponse{}
	if err := f.mgr.PopulateAttributes(bpReq, bpResp); err != nil {
		return nil, fmt.Errorf("failed to populate bridge port %d: %v", portOID, err)
	}

	portID := bpResp.GetAttr().GetPortId()
	if portID == 0 {
		return nil, fmt.Errorf("cannot find port ID for bridge port %d", portOID)
	}

	addReq := fwdconfig.TableEntryAddRequest(f.dataplane.ID(), FDBTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(
			fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).WithBytes(mac),
		)),
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
		FdbEntry: &saipb.FdbEntry{
			SwitchId:   entry.GetSwitchId(),
			MacAddress: mac,
			BvId:       entry.GetBvId(),
		},
		Attrs: []*saipb.FdbEntryAttribute{
			{
				BridgePortId: proto.Uint64(portOID),
			},
			{
				Type: entryType.Enum(),
			},
			{
				PacketAction: req.GetPacketAction().Enum(),
			},
		},
	})

	return &saipb.CreateFdbEntryResponse{}, nil
}

func (f *fdb) RemoveFdbEntry(ctx context.Context, req *saipb.RemoveFdbEntryRequest) (*saipb.RemoveFdbEntryResponse, error) {
	slog.InfoContext(ctx, "RemoveFdbEntry called", "mac", req.GetEntry().GetMacAddress(), "vlan", req.GetEntry().GetBvId())
	entry := req.GetEntry()
	if entry == nil {
		return nil, status.Errorf(codes.InvalidArgument, "FDB entry is required")
	}

	mac := entry.GetMacAddress()
	if len(mac) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "MAC address is required")
	}

	if fwdCtx, err := f.dataplane.FindContext(&fwdpb.ContextId{Id: f.dataplane.ID()}); err == nil && fwdCtx != nil {
		if obj, err := fwdCtx.Objects.FindID(&fwdpb.ObjectId{Id: FDBTable}); err == nil {
			if brTable, ok := obj.(*fwdbridge.Table); ok {
				_ = brTable.Remove(mac)
			}
		}
	}

	delReq := fwdconfig.TableEntryRemoveRequest(f.dataplane.ID(), FDBTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(
			fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).WithBytes(mac),
		)),
	).Build()

	if _, err := f.dataplane.TableEntryRemove(ctx, delReq); err != nil {
		slog.WarnContext(ctx, "failed to remove FDB entry from dataplane", "err", err)
	}

	key := fdbEntryKey{
		bvID: entry.GetBvId(),
		mac:  string(mac),
	}
	f.mu.Lock()
	rec := f.entries[key]
	delete(f.entries, key)
	f.mu.Unlock()

	var bpID uint64
	var entryType *saipb.FdbEntryType
	if rec != nil {
		bpID = rec.bridgePortID
		entryType = rec.entryType.Enum()
	}

	f.sendNotification(&saipb.FdbEventNotificationData{
		EventType: saipb.FdbEvent_FDB_EVENT_AGED,
		FdbEntry: &saipb.FdbEntry{
			SwitchId:   entry.GetSwitchId(),
			MacAddress: mac,
			BvId:       entry.GetBvId(),
		},
		Attrs: []*saipb.FdbEntryAttribute{
			{
				BridgePortId: proto.Uint64(bpID),
			},
			{
				Type: entryType,
			},
		},
	})

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

func (f *fdb) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.entries = make(map[fdbEntryKey]*fdbEntryRecord)
}
