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

#ifndef DATAPLANE_STANDALONE_SAI_NEXT_HOP_GROUP_H_
#define DATAPLANE_STANDALONE_SAI_NEXT_HOP_GROUP_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/next_hop_group.grpc.pb.h"
#include "dataplane/proto/sai/next_hop_group.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class NextHopGroup final
    : public lemming::dataplane::sai::NextHopGroup::Service {
 public:
  grpc::Status CreateNextHopGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNextHopGroupRequest* req,
      lemming::dataplane::sai::CreateNextHopGroupResponse* resp);

  grpc::Status RemoveNextHopGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNextHopGroupRequest* req,
      lemming::dataplane::sai::RemoveNextHopGroupResponse* resp);

  grpc::Status SetNextHopGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetNextHopGroupAttributeRequest* req,
      lemming::dataplane::sai::SetNextHopGroupAttributeResponse* resp);

  grpc::Status GetNextHopGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetNextHopGroupAttributeRequest* req,
      lemming::dataplane::sai::GetNextHopGroupAttributeResponse* resp);

  grpc::Status CreateNextHopGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNextHopGroupMemberRequest* req,
      lemming::dataplane::sai::CreateNextHopGroupMemberResponse* resp);

  grpc::Status RemoveNextHopGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNextHopGroupMemberRequest* req,
      lemming::dataplane::sai::RemoveNextHopGroupMemberResponse* resp);

  grpc::Status SetNextHopGroupMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetNextHopGroupMemberAttributeRequest* req,
      lemming::dataplane::sai::SetNextHopGroupMemberAttributeResponse* resp);

  grpc::Status GetNextHopGroupMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetNextHopGroupMemberAttributeRequest* req,
      lemming::dataplane::sai::GetNextHopGroupMemberAttributeResponse* resp);

  grpc::Status CreateNextHopGroupMembers(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNextHopGroupMembersRequest* req,
      lemming::dataplane::sai::CreateNextHopGroupMembersResponse* resp);

  grpc::Status RemoveNextHopGroupMembers(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNextHopGroupMembersRequest* req,
      lemming::dataplane::sai::RemoveNextHopGroupMembersResponse* resp);

  grpc::Status CreateNextHopGroupMap(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNextHopGroupMapRequest* req,
      lemming::dataplane::sai::CreateNextHopGroupMapResponse* resp);

  grpc::Status RemoveNextHopGroupMap(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNextHopGroupMapRequest* req,
      lemming::dataplane::sai::RemoveNextHopGroupMapResponse* resp);

  grpc::Status SetNextHopGroupMapAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetNextHopGroupMapAttributeRequest* req,
      lemming::dataplane::sai::SetNextHopGroupMapAttributeResponse* resp);

  grpc::Status GetNextHopGroupMapAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetNextHopGroupMapAttributeRequest* req,
      lemming::dataplane::sai::GetNextHopGroupMapAttributeResponse* resp);

  grpc::Status CreateNextHopGroups(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNextHopGroupsRequest* req,
      lemming::dataplane::sai::CreateNextHopGroupsResponse* resp);

  grpc::Status RemoveNextHopGroups(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNextHopGroupsRequest* req,
      lemming::dataplane::sai::RemoveNextHopGroupsResponse* resp);

  sai_next_hop_group_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_NEXT_HOP_GROUP_H_
