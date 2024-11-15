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

#ifndef DATAPLANE_STANDALONE_SAI_ROUTE_H_
#define DATAPLANE_STANDALONE_SAI_ROUTE_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/route.grpc.pb.h"
#include "dataplane/proto/sai/route.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Route final : public lemming::dataplane::sai::Route::Service {
 public:
  grpc::Status CreateRouteEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateRouteEntryRequest* req,
      lemming::dataplane::sai::CreateRouteEntryResponse* resp);

  grpc::Status RemoveRouteEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveRouteEntryRequest* req,
      lemming::dataplane::sai::RemoveRouteEntryResponse* resp);

  grpc::Status SetRouteEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetRouteEntryAttributeRequest* req,
      lemming::dataplane::sai::SetRouteEntryAttributeResponse* resp);

  grpc::Status GetRouteEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetRouteEntryAttributeRequest* req,
      lemming::dataplane::sai::GetRouteEntryAttributeResponse* resp);

  grpc::Status CreateRouteEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateRouteEntriesRequest* req,
      lemming::dataplane::sai::CreateRouteEntriesResponse* resp);

  grpc::Status RemoveRouteEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveRouteEntriesRequest* req,
      lemming::dataplane::sai::RemoveRouteEntriesResponse* resp);

  sai_route_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_ROUTE_H_
