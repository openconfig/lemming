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

	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

type l2mcGroup struct {
	saipb.UnimplementedL2McGroupServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

func newL2mcGroup(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *l2mcGroup {
	mg := &l2mcGroup{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterL2McGroupServer(s, mg)
	return mg
}

// TODO: Implement this.
func (mg *l2mcGroup) CreateL2McGroup(context.Context, *saipb.CreateL2McGroupRequest) (*saipb.CreateL2McGroupResponse, error) {
	id := mg.mgr.NextID()
	return &saipb.CreateL2McGroupResponse{Oid: id}, nil
}

// TODO: Implement this.
func (mg *l2mcGroup) CreateL2McGroupMember(context.Context, *saipb.CreateL2McGroupMemberRequest) (*saipb.CreateL2McGroupMemberResponse, error) {
	id := mg.mgr.NextID()
	return &saipb.CreateL2McGroupMemberResponse{Oid: id}, nil
}

// TODO: Implement this.
func (mg *l2mcGroup) RemoveL2McGroup(context.Context, *saipb.RemoveL2McGroupRequest) (*saipb.RemoveL2McGroupResponse, error) {
	return &saipb.RemoveL2McGroupResponse{}, nil
}

// TODO: Implement this.
func (mg *l2mcGroup) RemoveL2McGroupMember(context.Context, *saipb.RemoveL2McGroupMemberRequest) (*saipb.RemoveL2McGroupMemberResponse, error) {
	return &saipb.RemoveL2McGroupMemberResponse{}, nil
}

type l2mc struct {
	saipb.UnimplementedL2McServer
	mgr       *attrmgr.AttrMgr
	dataplane switchDataplaneAPI
}

// TODO: Implement this.
func newL2mc(mgr *attrmgr.AttrMgr, dataplane switchDataplaneAPI, s *grpc.Server) *l2mc {
	m := &l2mc{
		mgr:       mgr,
		dataplane: dataplane,
	}
	saipb.RegisterL2McServer(s, m)
	return m
}

// TODO: Implement this.
func (m *l2mc) CreateL2McEntry(context.Context, *saipb.CreateL2McEntryRequest) (*saipb.CreateL2McEntryResponse, error) {
	return &saipb.CreateL2McEntryResponse{}, nil
}

// TODO: Implement this.
func (m *l2mc) RemoveL2McEntry(context.Context, *saipb.RemoveL2McEntryRequest) (*saipb.RemoveL2McEntryResponse, error) {
	return &saipb.RemoveL2McEntryResponse{}, nil
}

// TODO: Implement this.
func (m *l2mc) SetL2McEntryAttribute(context.Context, *saipb.SetL2McEntryAttributeRequest) (*saipb.SetL2McEntryAttributeResponse, error) {
	return &saipb.SetL2McEntryAttributeResponse{}, nil
}
