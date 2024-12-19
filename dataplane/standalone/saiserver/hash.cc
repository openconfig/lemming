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

#include "dataplane/standalone/saiserver/hash.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/hash.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Hash::CreateHash(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHashRequest* req,
    lemming::dataplane::sai::CreateHashResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hash::RemoveHash(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHashRequest* req,
    lemming::dataplane::sai::RemoveHashResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_hash(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Hash::SetHashAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetHashAttributeRequest* req,
    lemming::dataplane::sai::SetHashAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hash::GetHashAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHashAttributeRequest* req,
    lemming::dataplane::sai::GetHashAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hash::CreateFineGrainedHashField(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateFineGrainedHashFieldRequest* req,
    lemming::dataplane::sai::CreateFineGrainedHashFieldResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hash::RemoveFineGrainedHashField(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveFineGrainedHashFieldRequest* req,
    lemming::dataplane::sai::RemoveFineGrainedHashFieldResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_fine_grained_hash_field(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Hash::GetFineGrainedHashFieldAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetFineGrainedHashFieldAttributeRequest* req,
    lemming::dataplane::sai::GetFineGrainedHashFieldAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
