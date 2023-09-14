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
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/google/go-cmp/cmp"

	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
)

func TestCreateSwitch(t *testing.T) {
	c, mgr, stopFn := newTestSwitch(t)
	defer stopFn()
	got, err := c.CreateSwitch(context.Background(), &saipb.CreateSwitchRequest{})
	if err != nil {
		t.Fatalf("CreateSwitch() unexpected error: %v", err)
	}
	want := &saipb.CreateSwitchResponse{
		Oid: 1,
	}
	if d := cmp.Diff(got, want, protocmp.Transform()); d != "" {
		t.Fatalf("CreateSwitch() failed: diff(-got,+want)\n:%s", d)
	}
	wantAttr := &saipb.SwitchAttribute{
		CpuPort:                          proto.Uint64(2),
		NumberOfActivePorts:              proto.Uint32(0),
		AclEntryMinimumPriority:          proto.Uint32(1),
		AclTableMaximumPriority:          proto.Uint32(100),
		MaxAclActionCount:                proto.Uint32(50),
		NumberOfEcmpGroups:               proto.Uint32(1024),
		DefaultVlanId:                    proto.Uint64(3),
		DefaultStpInstId:                 proto.Uint64(4),
		DefaultVirtualRouterId:           proto.Uint64(5),
		DefaultOverrideVirtualRouterId:   proto.Uint64(5),
		Default_1QBridgeId:               proto.Uint64(6),
		DefaultTrapGroup:                 proto.Uint64(7),
		IngressAcl:                       proto.Uint64(0),
		EgressAcl:                        proto.Uint64(0),
		QosMaxNumberOfTrafficClasses:     proto.Uint32(0),
		TotalBufferSize:                  proto.Uint64(1024 * 1024),
		AvailableIpv4RouteEntry:          proto.Uint32(1024),
		AvailableIpv6RouteEntry:          proto.Uint32(1024),
		AvailableIpv4NexthopEntry:        proto.Uint32(1024),
		AvailableIpv6NexthopEntry:        proto.Uint32(1024),
		AvailableIpv4NeighborEntry:       proto.Uint32(1024),
		AvailableIpv6NeighborEntry:       proto.Uint32(1024),
		AvailableNextHopGroupEntry:       proto.Uint32(1024),
		AvailableNextHopGroupMemberEntry: proto.Uint32(1024),
		AvailableFdbEntry:                proto.Uint32(1024),
		AvailableL2McEntry:               proto.Uint32(1024),
		AvailableIpmcEntry:               proto.Uint32(1024),
		AvailableSnatEntry:               proto.Uint32(1024),
		AvailableDnatEntry:               proto.Uint32(1024),
		EcmpHash:                         proto.Uint64(8),
		LagHash:                          proto.Uint64(8),
		RestartWarm:                      proto.Bool(false),
		WarmRecover:                      proto.Bool(false),
		LagDefaultHashAlgorithm:          saipb.HashAlgorithm_HASH_ALGORITHM_CRC.Enum(),
		LagDefaultHashSeed:               proto.Uint32(0),
		LagDefaultSymmetricHash:          proto.Bool(false),
		QosDefaultTc:                     proto.Uint32(0),
		QosDot1PToTcMap:                  proto.Uint64(0),
		QosDot1PToColorMap:               proto.Uint64(0),
		QosTcToQueueMap:                  proto.Uint64(0),
		QosTcAndColorToDot1PMap:          proto.Uint64(0),
		QosTcAndColorToDscpMap:           proto.Uint64(0),
		QosTcAndColorToMplsExpMap:        proto.Uint64(0),
		SwitchShellEnable:                proto.Bool(false),
		SwitchProfileId:                  proto.Uint32(0),
		NatZoneCounterObjectId:           proto.Uint64(0),
	}
	attr := &saipb.SwitchAttribute{}
	if err := mgr.PopulateAllAttributes("1", attr); err != nil {
		t.Fatal(err)
	}
	if d := cmp.Diff(attr, wantAttr, protocmp.Transform()); d != "" {
		t.Fatalf("CreateSwitch() failed: diff(-got,+want)\n:%s", d)
	}
}

func newTestSwitch(t testing.TB) (saipb.SwitchClient, *attrmgr.AttrMgr, func()) {
	t.Helper()
	mgr := attrmgr.New()
	lis, err := net.Listen("tcp", ("127.0.0.1:0"))
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()), grpc.ChainUnaryInterceptor(mgr.Interceptor))
	newSwitch(mgr, srv)
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve forwarding server: %v", err)
		}
	}()
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	return saipb.NewSwitchClient(conn), mgr, srv.Stop
}
