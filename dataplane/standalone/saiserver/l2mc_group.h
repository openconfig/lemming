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

#ifndef DATAPLANE_STANDALONE_SAI_L2MC_GROUP_H_
#define DATAPLANE_STANDALONE_SAI_L2MC_GROUP_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/l2mc_group.grpc.pb.h"
#include "dataplane/proto/sai/l2mc_group.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class L2mcGroup final : public lemming::dataplane::sai::L2mcGroup::Service {
 public:
  grpc::Status CreateL2mcGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateL2mcGroupRequest* req,
      lemming::dataplane::sai::CreateL2mcGroupResponse* resp);

  grpc::Status RemoveL2mcGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveL2mcGroupRequest* req,
      lemming::dataplane::sai::RemoveL2mcGroupResponse* resp);

  grpc::Status GetL2mcGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetL2mcGroupAttributeRequest* req,
      lemming::dataplane::sai::GetL2mcGroupAttributeResponse* resp);

  grpc::Status CreateL2mcGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateL2mcGroupMemberRequest* req,
      lemming::dataplane::sai::CreateL2mcGroupMemberResponse* resp);

  grpc::Status RemoveL2mcGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveL2mcGroupMemberRequest* req,
      lemming::dataplane::sai::RemoveL2mcGroupMemberResponse* resp);

  grpc::Status GetL2mcGroupMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetL2mcGroupMemberAttributeRequest* req,
      lemming::dataplane::sai::GetL2mcGroupMemberAttributeResponse* resp);

  sai_l2mc_group_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_L2MC_GROUP_H_
