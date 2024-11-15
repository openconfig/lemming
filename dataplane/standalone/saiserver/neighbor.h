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

#ifndef DATAPLANE_STANDALONE_SAI_NEIGHBOR_H_
#define DATAPLANE_STANDALONE_SAI_NEIGHBOR_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/neighbor.grpc.pb.h"
#include "dataplane/proto/sai/neighbor.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Neighbor final : public lemming::dataplane::sai::Neighbor::Service {
 public:
  grpc::Status CreateNeighborEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNeighborEntryRequest* req,
      lemming::dataplane::sai::CreateNeighborEntryResponse* resp);

  grpc::Status RemoveNeighborEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNeighborEntryRequest* req,
      lemming::dataplane::sai::RemoveNeighborEntryResponse* resp);

  grpc::Status SetNeighborEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetNeighborEntryAttributeRequest* req,
      lemming::dataplane::sai::SetNeighborEntryAttributeResponse* resp);

  grpc::Status GetNeighborEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetNeighborEntryAttributeRequest* req,
      lemming::dataplane::sai::GetNeighborEntryAttributeResponse* resp);

  grpc::Status CreateNeighborEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNeighborEntriesRequest* req,
      lemming::dataplane::sai::CreateNeighborEntriesResponse* resp);

  grpc::Status RemoveNeighborEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNeighborEntriesRequest* req,
      lemming::dataplane::sai::RemoveNeighborEntriesResponse* resp);

  sai_neighbor_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_NEIGHBOR_H_
