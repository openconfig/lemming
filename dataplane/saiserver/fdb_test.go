// Copyright 2026 Google LLC
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
	"bytes"
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

func newTestFdb(t testing.TB, api switchDataplaneAPI) (saipb.FdbClient, *fdb, *attrmgr.AttrMgr, func()) {
	var fdbServer *fdb
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		var err error
		fdbServer, err = newFdb(mgr, api, srv)
		if err != nil {
			t.Fatalf("newFdb failed: %v", err)
		}
	})
	return saipb.NewFdbClient(conn), fdbServer, mgr, stopFn
}

func TestCreateAndRemoveFdbEntry(t *testing.T) {
	dplane := &fakeSwitchDataplane{}
	c, fdbSrv, mgr, stopFn := newTestFdb(t, dplane)
	defer stopFn()
	ctx := context.Background()

	// Subscribe to notifications.
	notifCh := make(chan *saipb.FdbEventNotificationResponse, 10)
	unsub := fdbSrv.subscribe(notifCh)
	defer unsub()

	// Set up bridge port attribute.
	bpOID := uint64(10)
	mgr.StoreAttributes(bpOID, &saipb.BridgePortAttribute{
		PortId: proto.Uint64(1),
	})

	mac := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	bvID := uint64(20)

	// Create FDB Entry.
	_, err := c.CreateFdbEntry(ctx, &saipb.CreateFdbEntryRequest{
		Entry: &saipb.FdbEntry{
			SwitchId:   1,
			MacAddress: mac,
			BvId:       bvID,
		},
		BridgePortId: proto.Uint64(bpOID),
		Type:         saipb.FdbEntryType_FDB_ENTRY_TYPE_STATIC.Enum(),
		PacketAction: saipb.PacketAction_PACKET_ACTION_FORWARD.Enum(),
	})
	if err != nil {
		t.Fatalf("CreateFdbEntry() failed: %v", err)
	}

	// Verify notification.
	select {
	case notif := <-notifCh:
		if len(notif.GetData()) != 1 {
			t.Fatalf("Expected 1 notification item, got %d", len(notif.GetData()))
		}
		data := notif.GetData()[0]
		if data.GetEventType() != saipb.FdbEvent_FDB_EVENT_LEARNED {
			t.Errorf("Expected LEARNED event, got %v", data.GetEventType())
		}
		if !bytes.Equal(data.GetFdbEntry().GetMacAddress(), mac) {
			t.Errorf("Expected MAC %x, got %x", mac, data.GetFdbEntry().GetMacAddress())
		}
	default:
		t.Fatal("Expected FDB notification not received")
	}

	// Remove FDB Entry.
	_, err = c.RemoveFdbEntry(ctx, &saipb.RemoveFdbEntryRequest{
		Entry: &saipb.FdbEntry{
			SwitchId:   1,
			MacAddress: mac,
			BvId:       bvID,
		},
	})
	if err != nil {
		t.Fatalf("RemoveFdbEntry() failed: %v", err)
	}

	// Verify Aged/Removed notification.
	select {
	case notif := <-notifCh:
		if len(notif.GetData()) != 1 {
			t.Fatalf("Expected 1 notification item, got %d", len(notif.GetData()))
		}
		data := notif.GetData()[0]
		if data.GetEventType() != saipb.FdbEvent_FDB_EVENT_AGED {
			t.Errorf("Expected AGED event, got %v", data.GetEventType())
		}
	default:
		t.Fatal("Expected FDB aged notification not received")
	}
}

func TestFlushFdbEntries(t *testing.T) {
	dplane := &fakeSwitchDataplane{}
	c, fdbSrv, mgr, stopFn := newTestFdb(t, dplane)
	defer stopFn()
	ctx := context.Background()

	notifCh := make(chan *saipb.FdbEventNotificationResponse, 20)
	unsub := fdbSrv.subscribe(notifCh)
	defer unsub()

	mgr.StoreAttributes(10, &saipb.BridgePortAttribute{PortId: proto.Uint64(1)})
	mgr.StoreAttributes(20, &saipb.BridgePortAttribute{PortId: proto.Uint64(2)})
	mgr.StoreAttributes(30, &saipb.BridgePortAttribute{PortId: proto.Uint64(3)})

	mac1 := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	mac2 := []byte{0x00, 0x22, 0x33, 0x44, 0x55, 0x66}
	mac3 := []byte{0x00, 0x33, 0x44, 0x55, 0x66, 0x77}

	// Simulate learned dynamic entry 1 on VLAN 20, port 10.
	fdbSrv.sendNotification(&saipb.FdbEventNotificationData{
		EventType: saipb.FdbEvent_FDB_EVENT_LEARNED,
		FdbEntry: &saipb.FdbEntry{
			SwitchId:   1,
			MacAddress: mac1,
			BvId:       20,
		},
		Attrs: []*saipb.FdbEntryAttribute{
			{BridgePortId: proto.Uint64(10)},
			{Type: saipb.FdbEntryType_FDB_ENTRY_TYPE_DYNAMIC.Enum()},
		},
	})
	// Simulate learned dynamic entry 2 on VLAN 20, port 20.
	fdbSrv.sendNotification(&saipb.FdbEventNotificationData{
		EventType: saipb.FdbEvent_FDB_EVENT_LEARNED,
		FdbEntry: &saipb.FdbEntry{
			SwitchId:   1,
			MacAddress: mac2,
			BvId:       20,
		},
		Attrs: []*saipb.FdbEntryAttribute{
			{BridgePortId: proto.Uint64(20)},
			{Type: saipb.FdbEntryType_FDB_ENTRY_TYPE_DYNAMIC.Enum()},
		},
	})
	// Create static entry 3 on VLAN 30, port 30.
	_, err := c.CreateFdbEntry(ctx, &saipb.CreateFdbEntryRequest{
		Entry: &saipb.FdbEntry{
			SwitchId:   1,
			MacAddress: mac3,
			BvId:       30,
		},
		BridgePortId: proto.Uint64(30),
		Type:         saipb.FdbEntryType_FDB_ENTRY_TYPE_STATIC.Enum(),
	})
	if err != nil {
		t.Fatalf("CreateFdbEntry failed: %v", err)
	}

	// Drain setup notifications.
	for len(notifCh) > 0 {
		<-notifCh
	}

	// Test 1: Flush dynamic entries on VLAN 20.
	_, err = c.FlushFdbEntries(ctx, &saipb.FlushFdbEntriesRequest{
		Switch:    1,
		BvId:      proto.Uint64(20),
		EntryType: saipb.FdbFlushEntryType_FDB_FLUSH_ENTRY_TYPE_DYNAMIC.Enum(),
	})
	if err != nil {
		t.Fatalf("FlushFdbEntries failed: %v", err)
	}

	// Should receive 2 AGED notifications for mac1 and mac2.
	flushedMacs := make(map[string]bool)
	for i := 0; i < 2; i++ {
		select {
		case notif := <-notifCh:
			for _, d := range notif.GetData() {
				if d.GetEventType() == saipb.FdbEvent_FDB_EVENT_AGED {
					flushedMacs[string(d.GetFdbEntry().GetMacAddress())] = true
				}
			}
		default:
			t.Fatalf("Expected 2 flush notifications, got %d", len(flushedMacs))
		}
	}
	if !flushedMacs[string(mac1)] || !flushedMacs[string(mac2)] {
		t.Errorf("Expected mac1 and mac2 to be flushed, got %v", flushedMacs)
	}

	// Test 2: Flush all remaining entries.
	_, err = c.FlushFdbEntries(ctx, &saipb.FlushFdbEntriesRequest{
		Switch:    1,
		EntryType: saipb.FdbFlushEntryType_FDB_FLUSH_ENTRY_TYPE_ALL.Enum(),
	})
	if err != nil {
		t.Fatalf("FlushFdbEntries(ALL) failed: %v", err)
	}

	// Should receive 1 AGED notification for static mac3.
	select {
	case notif := <-notifCh:
		if len(notif.GetData()) != 1 || notif.GetData()[0].GetEventType() != saipb.FdbEvent_FDB_EVENT_AGED {
			t.Errorf("Expected mac3 flush notification, got %v", notif)
		}
		if !bytes.Equal(notif.GetData()[0].GetFdbEntry().GetMacAddress(), mac3) {
			t.Errorf("Expected mac3 %x, got %x", mac3, notif.GetData()[0].GetFdbEntry().GetMacAddress())
		}
	default:
		t.Fatal("Expected mac3 flush notification not received")
	}
}
