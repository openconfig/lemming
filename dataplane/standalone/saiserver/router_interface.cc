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

#include "dataplane/standalone/saiserver/router_interface.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/router_interface.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status RouterInterface::CreateRouterInterface(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateRouterInterfaceRequest* req,
    lemming::dataplane::sai::CreateRouterInterfaceResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status RouterInterface::RemoveRouterInterface(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveRouterInterfaceRequest* req,
    lemming::dataplane::sai::RemoveRouterInterfaceResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_router_interface(req.get_oid());

  auto status = api->remove_router_interface(entry);
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

grpc::Status RouterInterface::SetRouterInterfaceAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetRouterInterfaceAttributeRequest* req,
    lemming::dataplane::sai::SetRouterInterfaceAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status RouterInterface::GetRouterInterfaceAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetRouterInterfaceAttributeRequest* req,
    lemming::dataplane::sai::GetRouterInterfaceAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status RouterInterface::GetRouterInterfaceStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetRouterInterfaceStatsRequest* req,
    lemming::dataplane::sai::GetRouterInterfaceStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
