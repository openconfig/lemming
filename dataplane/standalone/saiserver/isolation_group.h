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

#ifndef DATAPLANE_STANDALONE_SAI_ISOLATION_GROUP_H_
#define DATAPLANE_STANDALONE_SAI_ISOLATION_GROUP_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/isolation_group.grpc.pb.h"
#include "dataplane/proto/sai/isolation_group.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class IsolationGroup final
    : public lemming::dataplane::sai::IsolationGroup::Service {
 public:
  grpc::Status CreateIsolationGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateIsolationGroupRequest* req,
      lemming::dataplane::sai::CreateIsolationGroupResponse* resp);

  grpc::Status RemoveIsolationGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveIsolationGroupRequest* req,
      lemming::dataplane::sai::RemoveIsolationGroupResponse* resp);

  grpc::Status GetIsolationGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIsolationGroupAttributeRequest* req,
      lemming::dataplane::sai::GetIsolationGroupAttributeResponse* resp);

  grpc::Status CreateIsolationGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateIsolationGroupMemberRequest* req,
      lemming::dataplane::sai::CreateIsolationGroupMemberResponse* resp);

  grpc::Status RemoveIsolationGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveIsolationGroupMemberRequest* req,
      lemming::dataplane::sai::RemoveIsolationGroupMemberResponse* resp);

  grpc::Status GetIsolationGroupMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIsolationGroupMemberAttributeRequest*
          req,
      lemming::dataplane::sai::GetIsolationGroupMemberAttributeResponse* resp);

  sai_isolation_group_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_ISOLATION_GROUP_H_
