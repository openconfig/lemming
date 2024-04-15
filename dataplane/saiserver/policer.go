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

	"google.golang.org/grpc"

	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

type policer struct {
	saipb.UnimplementedPolicerServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newPolicer(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *policer {
	p := &policer{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterPolicerServer(s, p)
	return p
}

// CreatePolicer return policer.
func (p *policer) CreatePolicer(context.Context, *saipb.CreatePolicerRequest) (*saipb.CreatePolicerResponse, error) {
	// TODO: provide implementation.
	return &saipb.CreatePolicerResponse{}, nil
}
