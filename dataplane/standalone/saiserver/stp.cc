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

#include "dataplane/standalone/saiserver/stp.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/stp.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Stp::CreateStp(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateStpRequest* req,
    lemming::dataplane::sai::CreateStpResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Stp::RemoveStp(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveStpRequest* req,
    lemming::dataplane::sai::RemoveStpResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_stp(req.get_oid());

  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return grpc::Status::INTERNAL;
  }

  return grpc::Status::OK;
}

grpc::Status Stp::GetStpAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetStpAttributeRequest* req,
    lemming::dataplane::sai::GetStpAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Stp::CreateStpPort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateStpPortRequest* req,
    lemming::dataplane::sai::CreateStpPortResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Stp::RemoveStpPort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveStpPortRequest* req,
    lemming::dataplane::sai::RemoveStpPortResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_stp_port(req.get_oid());

  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return grpc::Status::INTERNAL;
  }

  return grpc::Status::OK;
}

grpc::Status Stp::SetStpPortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetStpPortAttributeRequest* req,
    lemming::dataplane::sai::SetStpPortAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Stp::GetStpPortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetStpPortAttributeRequest* req,
    lemming::dataplane::sai::GetStpPortAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Stp::CreateStpPorts(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateStpPortsRequest* req,
    lemming::dataplane::sai::CreateStpPortsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Stp::RemoveStpPorts(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveStpPortsRequest* req,
    lemming::dataplane::sai::RemoveStpPortsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
