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

#include "dataplane/standalone/saiserver/bmtor.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/bmtor.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Bmtor::CreateTableBitmapClassificationEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTableBitmapClassificationEntryRequest*
        req,
    lemming::dataplane::sai::CreateTableBitmapClassificationEntryResponse*
        resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bmtor::RemoveTableBitmapClassificationEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTableBitmapClassificationEntryRequest*
        req,
    lemming::dataplane::sai::RemoveTableBitmapClassificationEntryResponse*
        resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_table_bitmap_classification_entry(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Bmtor::GetTableBitmapClassificationEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::
        GetTableBitmapClassificationEntryAttributeRequest* req,
    lemming::dataplane::sai::GetTableBitmapClassificationEntryAttributeResponse*
        resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bmtor::GetTableBitmapClassificationEntryStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::
        GetTableBitmapClassificationEntryStatsRequest* req,
    lemming::dataplane::sai::GetTableBitmapClassificationEntryStatsResponse*
        resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bmtor::CreateTableBitmapRouterEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTableBitmapRouterEntryRequest* req,
    lemming::dataplane::sai::CreateTableBitmapRouterEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bmtor::RemoveTableBitmapRouterEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTableBitmapRouterEntryRequest* req,
    lemming::dataplane::sai::RemoveTableBitmapRouterEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_table_bitmap_router_entry(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Bmtor::GetTableBitmapRouterEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTableBitmapRouterEntryAttributeRequest*
        req,
    lemming::dataplane::sai::GetTableBitmapRouterEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bmtor::GetTableBitmapRouterEntryStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTableBitmapRouterEntryStatsRequest* req,
    lemming::dataplane::sai::GetTableBitmapRouterEntryStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bmtor::CreateTableMetaTunnelEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTableMetaTunnelEntryRequest* req,
    lemming::dataplane::sai::CreateTableMetaTunnelEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bmtor::RemoveTableMetaTunnelEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTableMetaTunnelEntryRequest* req,
    lemming::dataplane::sai::RemoveTableMetaTunnelEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_table_meta_tunnel_entry(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Bmtor::GetTableMetaTunnelEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTableMetaTunnelEntryAttributeRequest* req,
    lemming::dataplane::sai::GetTableMetaTunnelEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Bmtor::GetTableMetaTunnelEntryStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTableMetaTunnelEntryStatsRequest* req,
    lemming::dataplane::sai::GetTableMetaTunnelEntryStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
