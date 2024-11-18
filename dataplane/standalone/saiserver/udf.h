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

#ifndef DATAPLANE_STANDALONE_SAI_UDF_H_
#define DATAPLANE_STANDALONE_SAI_UDF_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/udf.grpc.pb.h"
#include "dataplane/proto/sai/udf.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Udf final : public lemming::dataplane::sai::Udf::Service {
 public:
  grpc::Status CreateUdf(grpc::ServerContext* context,
                         const lemming::dataplane::sai::CreateUdfRequest* req,
                         lemming::dataplane::sai::CreateUdfResponse* resp);

  grpc::Status RemoveUdf(grpc::ServerContext* context,
                         const lemming::dataplane::sai::RemoveUdfRequest* req,
                         lemming::dataplane::sai::RemoveUdfResponse* resp);

  grpc::Status SetUdfAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetUdfAttributeRequest* req,
      lemming::dataplane::sai::SetUdfAttributeResponse* resp);

  grpc::Status GetUdfAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetUdfAttributeRequest* req,
      lemming::dataplane::sai::GetUdfAttributeResponse* resp);

  grpc::Status CreateUdfMatch(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateUdfMatchRequest* req,
      lemming::dataplane::sai::CreateUdfMatchResponse* resp);

  grpc::Status RemoveUdfMatch(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveUdfMatchRequest* req,
      lemming::dataplane::sai::RemoveUdfMatchResponse* resp);

  grpc::Status GetUdfMatchAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetUdfMatchAttributeRequest* req,
      lemming::dataplane::sai::GetUdfMatchAttributeResponse* resp);

  grpc::Status CreateUdfGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateUdfGroupRequest* req,
      lemming::dataplane::sai::CreateUdfGroupResponse* resp);

  grpc::Status RemoveUdfGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveUdfGroupRequest* req,
      lemming::dataplane::sai::RemoveUdfGroupResponse* resp);

  grpc::Status GetUdfGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetUdfGroupAttributeRequest* req,
      lemming::dataplane::sai::GetUdfGroupAttributeResponse* resp);

  sai_udf_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_UDF_H_
