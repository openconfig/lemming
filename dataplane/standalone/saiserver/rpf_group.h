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

#ifndef DATAPLANE_STANDALONE_SAI_RPF_GROUP_H_
#define DATAPLANE_STANDALONE_SAI_RPF_GROUP_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/rpf_group.grpc.pb.h"
#include "dataplane/proto/sai/rpf_group.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class RpfGroup final : public lemming::dataplane::sai::RpfGroup::Service {
 public:
  grpc::Status CreateRpfGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateRpfGroupRequest* req,
      lemming::dataplane::sai::CreateRpfGroupResponse* resp);

  grpc::Status RemoveRpfGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveRpfGroupRequest* req,
      lemming::dataplane::sai::RemoveRpfGroupResponse* resp);

  grpc::Status GetRpfGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetRpfGroupAttributeRequest* req,
      lemming::dataplane::sai::GetRpfGroupAttributeResponse* resp);

  grpc::Status CreateRpfGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateRpfGroupMemberRequest* req,
      lemming::dataplane::sai::CreateRpfGroupMemberResponse* resp);

  grpc::Status RemoveRpfGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveRpfGroupMemberRequest* req,
      lemming::dataplane::sai::RemoveRpfGroupMemberResponse* resp);

  grpc::Status GetRpfGroupMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetRpfGroupMemberAttributeRequest* req,
      lemming::dataplane::sai::GetRpfGroupMemberAttributeResponse* resp);

  sai_rpf_group_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_RPF_GROUP_H_
