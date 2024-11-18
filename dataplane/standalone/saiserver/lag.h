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

#ifndef DATAPLANE_STANDALONE_SAI_LAG_H_
#define DATAPLANE_STANDALONE_SAI_LAG_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/lag.grpc.pb.h"
#include "dataplane/proto/sai/lag.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Lag final : public lemming::dataplane::sai::Lag::Service {
 public:
  grpc::Status CreateLag(grpc::ServerContext* context,
                         const lemming::dataplane::sai::CreateLagRequest* req,
                         lemming::dataplane::sai::CreateLagResponse* resp);

  grpc::Status RemoveLag(grpc::ServerContext* context,
                         const lemming::dataplane::sai::RemoveLagRequest* req,
                         lemming::dataplane::sai::RemoveLagResponse* resp);

  grpc::Status SetLagAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetLagAttributeRequest* req,
      lemming::dataplane::sai::SetLagAttributeResponse* resp);

  grpc::Status GetLagAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetLagAttributeRequest* req,
      lemming::dataplane::sai::GetLagAttributeResponse* resp);

  grpc::Status CreateLagMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateLagMemberRequest* req,
      lemming::dataplane::sai::CreateLagMemberResponse* resp);

  grpc::Status RemoveLagMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveLagMemberRequest* req,
      lemming::dataplane::sai::RemoveLagMemberResponse* resp);

  grpc::Status SetLagMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetLagMemberAttributeRequest* req,
      lemming::dataplane::sai::SetLagMemberAttributeResponse* resp);

  grpc::Status GetLagMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetLagMemberAttributeRequest* req,
      lemming::dataplane::sai::GetLagMemberAttributeResponse* resp);

  grpc::Status CreateLagMembers(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateLagMembersRequest* req,
      lemming::dataplane::sai::CreateLagMembersResponse* resp);

  grpc::Status RemoveLagMembers(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveLagMembersRequest* req,
      lemming::dataplane::sai::RemoveLagMembersResponse* resp);

  sai_lag_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_LAG_H_
