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

#ifndef DATAPLANE_STANDALONE_SAI_ROUTER_INTERFACE_H_
#define DATAPLANE_STANDALONE_SAI_ROUTER_INTERFACE_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/router_interface.grpc.pb.h"
#include "dataplane/proto/sai/router_interface.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class RouterInterface final
    : public lemming::dataplane::sai::RouterInterface::Service {
 public:
  grpc::Status CreateRouterInterface(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateRouterInterfaceRequest* req,
      lemming::dataplane::sai::CreateRouterInterfaceResponse* resp);

  grpc::Status RemoveRouterInterface(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveRouterInterfaceRequest* req,
      lemming::dataplane::sai::RemoveRouterInterfaceResponse* resp);

  grpc::Status SetRouterInterfaceAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetRouterInterfaceAttributeRequest* req,
      lemming::dataplane::sai::SetRouterInterfaceAttributeResponse* resp);

  grpc::Status GetRouterInterfaceAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetRouterInterfaceAttributeRequest* req,
      lemming::dataplane::sai::GetRouterInterfaceAttributeResponse* resp);

  grpc::Status GetRouterInterfaceStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetRouterInterfaceStatsRequest* req,
      lemming::dataplane::sai::GetRouterInterfaceStatsResponse* resp);

  sai_router_interface_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_ROUTER_INTERFACE_H_
