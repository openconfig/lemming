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

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type debugCounter struct {
	saipb.UnimplementedDebugCounterServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newDebugCounter(mgr *attrmgr.AttrMgr, engine switchDataplaneAPI, s *grpc.Server) *debugCounter {
	d := &debugCounter{
		mgr:       mgr,
		dataplane: engine,
	}
	saipb.RegisterDebugCounterServer(s, d)
	return d
}

func (d *debugCounter) CreateDebugCounter(ctx context.Context, req *saipb.CreateDebugCounterRequest) (*saipb.CreateDebugCounterResponse, error) {
	oid := d.mgr.NextID()
	index := uint32(0)
	for _, reason := range req.InDropReasonList {
		if reason == saipb.InDropReason_IN_DROP_REASON_LPM4_MISS {
			index = 0
			break
		}
		if reason == saipb.InDropReason_IN_DROP_REASON_LPM6_MISS {
			index = 1
			break
		}
	}

	d.mgr.StoreAttributes(oid, &saipb.DebugCounterAttribute{
		Index: proto.Uint32(index),
	})
	return &saipb.CreateDebugCounterResponse{Oid: oid}, nil
}

func (d *debugCounter) RemoveDebugCounter(ctx context.Context, req *saipb.RemoveDebugCounterRequest) (*saipb.RemoveDebugCounterResponse, error) {
	return &saipb.RemoveDebugCounterResponse{}, nil
}

func (d *debugCounter) SetDebugCounterAttribute(ctx context.Context, req *saipb.SetDebugCounterAttributeRequest) (*saipb.SetDebugCounterAttributeResponse, error) {
	return &saipb.SetDebugCounterAttributeResponse{}, nil
}

