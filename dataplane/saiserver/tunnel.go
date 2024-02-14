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

	"google.golang.org/grpc"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
)

type tunnel struct {
	saipb.UnimplementedTunnelServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newTunnel(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *tunnel {
	t := &tunnel{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterTunnelServer(s, t)
	return t
}

func (t *tunnel) CreateTunnel(ctx context.Context, req *saipb.CreateTunnelRequest) (*saipb.CreateTunnelResponse, error) {
	over := req.GetOverlayInterface()
	under := req.GetUnderlayInterface()
	tunType := req.GetType()
	ecnMode := req.GetDecapEcnMode()
	ttlMode := req.GetDecapTtlMode()
	dscpMode := req.GetDecapDscpMode()

	id := t.mgr.NextID()
	return &saipb.CreateTunnelResponse{
		Oid: id,
	}, nil
}

// func (t *tunnel) CreateTunnelTermTableEntry(ctx context.Context, req *saipb.CreateTunnelTermTableEntryRequest) (*saipb.CreateTunnelTermTableEntryResponse, error) {
// 	vr := req.GetVrId()
// 	entryType := req.GetType()
// 	tunnelType := req.GetTunnelType()
// 	tunnel := req.GetActionTunnelId()
// 	dstIP := req.GetDstIp()

// 	id := t.mgr.NextID()
// 	return &saipb.CreateTunnelTermTableEntryResponse{
// 		Oid: id,
// 	}, nil
// }
