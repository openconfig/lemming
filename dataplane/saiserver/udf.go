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

	"google.golang.org/grpc"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
)

type udf struct {
	saipb.UnimplementedUdfServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newUdf(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *udf {
	udf := &udf{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterUdfServer(s, udf)
	return udf
}

func (u *udf) CreateUdf(ctx context.Context, req *saipb.CreateUdfRequest) (*saipb.CreateUdfResponse, error) {
	id := u.mgr.NextID()

	udfGroup := &saipb.UdfGroupAttribute{}
	if err := u.mgr.PopulateAllAttributes(fmt.Sprint(req.GetGroupId()), udfGroup); err != nil {
		return nil, err
	}
	udfGroup.UdfList = append(udfGroup.UdfList, id)
	u.mgr.StoreAttributes(id, udfGroup)

	return &saipb.CreateUdfResponse{
		Oid: id,
	}, nil
}

func (u *udf) CreateUdfGroup(context.Context, *saipb.CreateUdfGroupRequest) (*saipb.CreateUdfGroupResponse, error) {
	return &saipb.CreateUdfGroupResponse{
		Oid: u.mgr.NextID(),
	}, nil
}

func (u *udf) CreateUdfMatch(context.Context, *saipb.CreateUdfMatchRequest) (*saipb.CreateUdfMatchResponse, error) {
	return &saipb.CreateUdfMatchResponse{
		Oid: u.mgr.NextID(),
	}, nil
}
