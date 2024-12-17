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

#include "dataplane/standalone/saiserver/system_port.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/system_port.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status SystemPort::CreateSystemPort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateSystemPortRequest* req,
    lemming::dataplane::sai::CreateSystemPortResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status SystemPort::RemoveSystemPort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveSystemPortRequest* req,
    lemming::dataplane::sai::RemoveSystemPortResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_system_port(req.get_oid());

  auto status = api->remove_system_port(entry);
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

grpc::Status SystemPort::SetSystemPortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetSystemPortAttributeRequest* req,
    lemming::dataplane::sai::SetSystemPortAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status SystemPort::GetSystemPortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetSystemPortAttributeRequest* req,
    lemming::dataplane::sai::GetSystemPortAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
