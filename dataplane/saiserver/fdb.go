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

type fdb struct {
	saipb.UnimplementedFdbServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newFdb(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) (*fdb, error) {
	f := &fdb{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterFdbServer(s, f)
	return f, nil
}

func (f *fdb) CreateFdbFlush(ctx context.Context, req *saipb.CreateFdbFlushRequest) (*saipb.CreateFdbFlushResponse, error) {
	return &saipb.CreateFdbFlushResponse{}, nil
}
