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

#include "dataplane/standalone/saiserver/bridge.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/bridge.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Bridge::CreateBridge(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateBridgeRequest* req,
    lemming::dataplane::sai::CreateBridgeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bridge::RemoveBridge(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveBridgeRequest* req,
    lemming::dataplane::sai::RemoveBridgeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_bridge(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Bridge::SetBridgeAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetBridgeAttributeRequest* req,
    lemming::dataplane::sai::SetBridgeAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bridge::GetBridgeAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetBridgeAttributeRequest* req,
    lemming::dataplane::sai::GetBridgeAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bridge::GetBridgeStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetBridgeStatsRequest* req,
    lemming::dataplane::sai::GetBridgeStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bridge::CreateBridgePort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateBridgePortRequest* req,
    lemming::dataplane::sai::CreateBridgePortResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bridge::RemoveBridgePort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveBridgePortRequest* req,
    lemming::dataplane::sai::RemoveBridgePortResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_bridge_port(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Bridge::SetBridgePortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetBridgePortAttributeRequest* req,
    lemming::dataplane::sai::SetBridgePortAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bridge::GetBridgePortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetBridgePortAttributeRequest* req,
    lemming::dataplane::sai::GetBridgePortAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bridge::GetBridgePortStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetBridgePortStatsRequest* req,
    lemming::dataplane::sai::GetBridgePortStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
