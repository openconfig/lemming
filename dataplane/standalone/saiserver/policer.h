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

#ifndef DATAPLANE_STANDALONE_SAI_POLICER_H_
#define DATAPLANE_STANDALONE_SAI_POLICER_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/policer.grpc.pb.h"
#include "dataplane/proto/sai/policer.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Policer final : public lemming::dataplane::sai::Policer::Service {
 public:
  grpc::Status CreatePolicer(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreatePolicerRequest* req,
      lemming::dataplane::sai::CreatePolicerResponse* resp);

  grpc::Status RemovePolicer(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemovePolicerRequest* req,
      lemming::dataplane::sai::RemovePolicerResponse* resp);

  grpc::Status SetPolicerAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetPolicerAttributeRequest* req,
      lemming::dataplane::sai::SetPolicerAttributeResponse* resp);

  grpc::Status GetPolicerAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetPolicerAttributeRequest* req,
      lemming::dataplane::sai::GetPolicerAttributeResponse* resp);

  grpc::Status GetPolicerStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetPolicerStatsRequest* req,
      lemming::dataplane::sai::GetPolicerStatsResponse* resp);

  sai_policer_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_POLICER_H_
