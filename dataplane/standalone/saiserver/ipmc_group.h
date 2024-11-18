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

#ifndef DATAPLANE_STANDALONE_SAI_IPMC_GROUP_H_
#define DATAPLANE_STANDALONE_SAI_IPMC_GROUP_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/ipmc_group.grpc.pb.h"
#include "dataplane/proto/sai/ipmc_group.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class IpmcGroup final : public lemming::dataplane::sai::IpmcGroup::Service {
 public:
  grpc::Status CreateIpmcGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateIpmcGroupRequest* req,
      lemming::dataplane::sai::CreateIpmcGroupResponse* resp);

  grpc::Status RemoveIpmcGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveIpmcGroupRequest* req,
      lemming::dataplane::sai::RemoveIpmcGroupResponse* resp);

  grpc::Status GetIpmcGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIpmcGroupAttributeRequest* req,
      lemming::dataplane::sai::GetIpmcGroupAttributeResponse* resp);

  grpc::Status CreateIpmcGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateIpmcGroupMemberRequest* req,
      lemming::dataplane::sai::CreateIpmcGroupMemberResponse* resp);

  grpc::Status RemoveIpmcGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveIpmcGroupMemberRequest* req,
      lemming::dataplane::sai::RemoveIpmcGroupMemberResponse* resp);

  grpc::Status GetIpmcGroupMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIpmcGroupMemberAttributeRequest* req,
      lemming::dataplane::sai::GetIpmcGroupMemberAttributeResponse* resp);

  sai_ipmc_group_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_IPMC_GROUP_H_
