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

#ifndef DATAPLANE_STANDALONE_SAI_VIRTUAL_ROUTER_H_
#define DATAPLANE_STANDALONE_SAI_VIRTUAL_ROUTER_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/virtual_router.grpc.pb.h"
#include "dataplane/proto/sai/virtual_router.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class VirtualRouter final
    : public lemming::dataplane::sai::VirtualRouter::Service {
 public:
  grpc::Status CreateVirtualRouter(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateVirtualRouterRequest* req,
      lemming::dataplane::sai::CreateVirtualRouterResponse* resp);

  grpc::Status RemoveVirtualRouter(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveVirtualRouterRequest* req,
      lemming::dataplane::sai::RemoveVirtualRouterResponse* resp);

  grpc::Status SetVirtualRouterAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetVirtualRouterAttributeRequest* req,
      lemming::dataplane::sai::SetVirtualRouterAttributeResponse* resp);

  grpc::Status GetVirtualRouterAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetVirtualRouterAttributeRequest* req,
      lemming::dataplane::sai::GetVirtualRouterAttributeResponse* resp);

  sai_virtual_router_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_VIRTUAL_ROUTER_H_
