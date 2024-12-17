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

#include "dataplane/standalone/saiserver/my_mac.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/my_mac.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status MyMac::CreateMyMac(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateMyMacRequest* req,
    lemming::dataplane::sai::CreateMyMacResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status MyMac::RemoveMyMac(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveMyMacRequest* req,
    lemming::dataplane::sai::RemoveMyMacResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_my_mac(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status MyMac::SetMyMacAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetMyMacAttributeRequest* req,
    lemming::dataplane::sai::SetMyMacAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status MyMac::GetMyMacAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMyMacAttributeRequest* req,
    lemming::dataplane::sai::GetMyMacAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
