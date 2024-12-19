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

#include "dataplane/standalone/saiserver/udf.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/udf.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Udf::CreateUdf(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateUdfRequest* req,
    lemming::dataplane::sai::CreateUdfResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Udf::RemoveUdf(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveUdfRequest* req,
    lemming::dataplane::sai::RemoveUdfResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_udf(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Udf::SetUdfAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetUdfAttributeRequest* req,
    lemming::dataplane::sai::SetUdfAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Udf::GetUdfAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetUdfAttributeRequest* req,
    lemming::dataplane::sai::GetUdfAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Udf::CreateUdfMatch(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateUdfMatchRequest* req,
    lemming::dataplane::sai::CreateUdfMatchResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Udf::RemoveUdfMatch(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveUdfMatchRequest* req,
    lemming::dataplane::sai::RemoveUdfMatchResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_udf_match(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Udf::GetUdfMatchAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetUdfMatchAttributeRequest* req,
    lemming::dataplane::sai::GetUdfMatchAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Udf::CreateUdfGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateUdfGroupRequest* req,
    lemming::dataplane::sai::CreateUdfGroupResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Udf::RemoveUdfGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveUdfGroupRequest* req,
    lemming::dataplane::sai::RemoveUdfGroupResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_udf_group(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Udf::GetUdfGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetUdfGroupAttributeRequest* req,
    lemming::dataplane::sai::GetUdfGroupAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
