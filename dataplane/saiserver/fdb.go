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
	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

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
	f.mu.RLock()
	defer f.mu.RUnlock()
	resp := &saipb.FdbEventNotificationResponse{
		Data: []*saipb.FdbEventNotificationData{data},
	}
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

	f.sendNotification(&saipb.FdbEventNotificationData{
		EventType: saipb.FdbEvent_FDB_EVENT_LEARNED,
		FdbEntry: &saipb.FdbEntry{
			SwitchId:   entry.GetSwitchId(),
			MacAddress: mac,
			BvId:       entry.GetBvId(),
		},
		Attrs: []*saipb.FdbEntryAttribute{
			{
				AttributeId: proto.Int32(int32(saipb.FdbEntryAttr_FDB_ENTRY_ATTR_BRIDGE_PORT_ID)),
				Value: &saipb.AttributeValue{
					Value: &saipb.AttributeValue_Oid{
						Oid: portOID,
					},
				},
			},
		},
	})

	return &saipb.CreateFdbEntryResponse{}, nil
}

func (f *fdb) RemoveFdbEntry(ctx context.Context, req *saipb.RemoveFdbEntryRequest) (*saipb.RemoveFdbEntryResponse, error) {
	entry := req.GetEntry()
	if entry == nil {
		return nil, status.Errorf(codes.InvalidArgument, "FDB entry is required")
	}

	mac := entry.GetMacAddress()
	if len(mac) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "MAC address is required")
	}

	delReq := fwdconfig.TableEntryRemoveRequest(f.dataplane.ID(), FDBTable).AppendEntry(
		fwdconfig.EntryDesc(fwdconfig.ExactEntry(
			fwdconfig.PacketFieldBytes(fwdpb.PacketFieldNum_PACKET_FIELD_NUM_ETHER_MAC_DST).WithBytes(mac),
		)),
	).Build()

	if _, err := f.dataplane.TableEntryRemove(ctx, delReq); err != nil {
		return nil, fmt.Errorf("failed to remove FDB entry from dataplane: %v", err)
	}

	f.sendNotification(&saipb.FdbEventNotificationData{
		EventType: saipb.FdbEvent_FDB_EVENT_AGED,
		FdbEntry: &saipb.FdbEntry{
			SwitchId:   entry.GetSwitchId(),
			MacAddress: mac,
			BvId:       entry.GetBvId(),
		},
	})

	return &saipb.RemoveFdbEntryResponse{}, nil
}
