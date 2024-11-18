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

#ifndef DATAPLANE_STANDALONE_SAI_BUFFER_H_
#define DATAPLANE_STANDALONE_SAI_BUFFER_H_

#include "dataplane/proto/sai/buffer.grpc.pb.h"
#include "dataplane/proto/sai/buffer.pb.h"
#include "dataplane/proto/sai/common.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Buffer final : public lemming::dataplane::sai::Buffer::Service {
 public:
  grpc::Status CreateBufferPool(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateBufferPoolRequest* req,
      lemming::dataplane::sai::CreateBufferPoolResponse* resp);

  grpc::Status RemoveBufferPool(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveBufferPoolRequest* req,
      lemming::dataplane::sai::RemoveBufferPoolResponse* resp);

  grpc::Status SetBufferPoolAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetBufferPoolAttributeRequest* req,
      lemming::dataplane::sai::SetBufferPoolAttributeResponse* resp);

  grpc::Status GetBufferPoolAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetBufferPoolAttributeRequest* req,
      lemming::dataplane::sai::GetBufferPoolAttributeResponse* resp);

  grpc::Status GetBufferPoolStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetBufferPoolStatsRequest* req,
      lemming::dataplane::sai::GetBufferPoolStatsResponse* resp);

  grpc::Status CreateIngressPriorityGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateIngressPriorityGroupRequest* req,
      lemming::dataplane::sai::CreateIngressPriorityGroupResponse* resp);

  grpc::Status RemoveIngressPriorityGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveIngressPriorityGroupRequest* req,
      lemming::dataplane::sai::RemoveIngressPriorityGroupResponse* resp);

  grpc::Status SetIngressPriorityGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetIngressPriorityGroupAttributeRequest*
          req,
      lemming::dataplane::sai::SetIngressPriorityGroupAttributeResponse* resp);

  grpc::Status GetIngressPriorityGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIngressPriorityGroupAttributeRequest*
          req,
      lemming::dataplane::sai::GetIngressPriorityGroupAttributeResponse* resp);

  grpc::Status GetIngressPriorityGroupStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIngressPriorityGroupStatsRequest* req,
      lemming::dataplane::sai::GetIngressPriorityGroupStatsResponse* resp);

  grpc::Status CreateBufferProfile(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateBufferProfileRequest* req,
      lemming::dataplane::sai::CreateBufferProfileResponse* resp);

  grpc::Status RemoveBufferProfile(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveBufferProfileRequest* req,
      lemming::dataplane::sai::RemoveBufferProfileResponse* resp);

  grpc::Status SetBufferProfileAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetBufferProfileAttributeRequest* req,
      lemming::dataplane::sai::SetBufferProfileAttributeResponse* resp);

  grpc::Status GetBufferProfileAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetBufferProfileAttributeRequest* req,
      lemming::dataplane::sai::GetBufferProfileAttributeResponse* resp);

  sai_buffer_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_BUFFER_H_
