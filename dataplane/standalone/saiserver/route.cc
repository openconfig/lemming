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

#include "dataplane/standalone/saiserver/route.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/route.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Route::CreateRouteEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateRouteEntryRequest* req,
    lemming::dataplane::sai::CreateRouteEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Route::RemoveRouteEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveRouteEntryRequest* req,
    lemming::dataplane::sai::RemoveRouteEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto entry = convert_to_route_entry(req->entry());
  auto status = api->remove_route_entry(&entry);

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Route::SetRouteEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetRouteEntryAttributeRequest* req,
    lemming::dataplane::sai::SetRouteEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Route::GetRouteEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetRouteEntryAttributeRequest* req,
    lemming::dataplane::sai::GetRouteEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Route::CreateRouteEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateRouteEntriesRequest* req,
    lemming::dataplane::sai::CreateRouteEntriesResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Route::RemoveRouteEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveRouteEntriesRequest* req,
    lemming::dataplane::sai::RemoveRouteEntriesResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
