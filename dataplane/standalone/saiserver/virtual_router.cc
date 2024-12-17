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

#include "dataplane/standalone/saiserver/virtual_router.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/virtual_router.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status VirtualRouter::CreateVirtualRouter(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateVirtualRouterRequest* req,
    lemming::dataplane::sai::CreateVirtualRouterResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status VirtualRouter::RemoveVirtualRouter(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveVirtualRouterRequest* req,
    lemming::dataplane::sai::RemoveVirtualRouterResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_virtual_router(req.get_oid());

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

grpc::Status VirtualRouter::SetVirtualRouterAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetVirtualRouterAttributeRequest* req,
    lemming::dataplane::sai::SetVirtualRouterAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status VirtualRouter::GetVirtualRouterAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetVirtualRouterAttributeRequest* req,
    lemming::dataplane::sai::GetVirtualRouterAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
