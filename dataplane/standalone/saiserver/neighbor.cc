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

#include "dataplane/standalone/saiserver/neighbor.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/neighbor.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Neighbor::CreateNeighborEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNeighborEntryRequest* req,
    lemming::dataplane::sai::CreateNeighborEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Neighbor::RemoveNeighborEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNeighborEntryRequest* req,
    lemming::dataplane::sai::RemoveNeighborEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;

  const sai_neighbor_entry_t* neighbor_entry entry =
      convert_to_neighbor_entry(req);
  auto status = api->remove_neighbor_entry(entry);

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

grpc::Status Neighbor::SetNeighborEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetNeighborEntryAttributeRequest* req,
    lemming::dataplane::sai::SetNeighborEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Neighbor::GetNeighborEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetNeighborEntryAttributeRequest* req,
    lemming::dataplane::sai::GetNeighborEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Neighbor::CreateNeighborEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNeighborEntriesRequest* req,
    lemming::dataplane::sai::CreateNeighborEntriesResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Neighbor::RemoveNeighborEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNeighborEntriesRequest* req,
    lemming::dataplane::sai::RemoveNeighborEntriesResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
