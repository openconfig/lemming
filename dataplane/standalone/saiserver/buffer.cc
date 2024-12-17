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

#include "dataplane/standalone/saiserver/buffer.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/buffer.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Buffer::CreateBufferPool(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateBufferPoolRequest* req,
    lemming::dataplane::sai::CreateBufferPoolResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Buffer::RemoveBufferPool(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveBufferPoolRequest* req,
    lemming::dataplane::sai::RemoveBufferPoolResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_buffer_pool(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Buffer::SetBufferPoolAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetBufferPoolAttributeRequest* req,
    lemming::dataplane::sai::SetBufferPoolAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Buffer::GetBufferPoolAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetBufferPoolAttributeRequest* req,
    lemming::dataplane::sai::GetBufferPoolAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Buffer::GetBufferPoolStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetBufferPoolStatsRequest* req,
    lemming::dataplane::sai::GetBufferPoolStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Buffer::CreateIngressPriorityGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateIngressPriorityGroupRequest* req,
    lemming::dataplane::sai::CreateIngressPriorityGroupResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Buffer::RemoveIngressPriorityGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveIngressPriorityGroupRequest* req,
    lemming::dataplane::sai::RemoveIngressPriorityGroupResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_ingress_priority_group(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Buffer::SetIngressPriorityGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetIngressPriorityGroupAttributeRequest* req,
    lemming::dataplane::sai::SetIngressPriorityGroupAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Buffer::GetIngressPriorityGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIngressPriorityGroupAttributeRequest* req,
    lemming::dataplane::sai::GetIngressPriorityGroupAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Buffer::GetIngressPriorityGroupStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIngressPriorityGroupStatsRequest* req,
    lemming::dataplane::sai::GetIngressPriorityGroupStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Buffer::CreateBufferProfile(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateBufferProfileRequest* req,
    lemming::dataplane::sai::CreateBufferProfileResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Buffer::RemoveBufferProfile(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveBufferProfileRequest* req,
    lemming::dataplane::sai::RemoveBufferProfileResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_buffer_profile(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Buffer::SetBufferProfileAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetBufferProfileAttributeRequest* req,
    lemming::dataplane::sai::SetBufferProfileAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Buffer::GetBufferProfileAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetBufferProfileAttributeRequest* req,
    lemming::dataplane::sai::GetBufferProfileAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
