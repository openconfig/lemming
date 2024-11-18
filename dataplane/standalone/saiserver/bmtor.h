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

#ifndef DATAPLANE_STANDALONE_SAI_BMTOR_H_
#define DATAPLANE_STANDALONE_SAI_BMTOR_H_

#include "dataplane/proto/sai/bmtor.grpc.pb.h"
#include "dataplane/proto/sai/bmtor.pb.h"
#include "dataplane/proto/sai/common.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Bmtor final : public lemming::dataplane::sai::Bmtor::Service {
 public:
  grpc::Status CreateTableBitmapClassificationEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::
          CreateTableBitmapClassificationEntryRequest* req,
      lemming::dataplane::sai::CreateTableBitmapClassificationEntryResponse*
          resp);

  grpc::Status RemoveTableBitmapClassificationEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::
          RemoveTableBitmapClassificationEntryRequest* req,
      lemming::dataplane::sai::RemoveTableBitmapClassificationEntryResponse*
          resp);

  grpc::Status GetTableBitmapClassificationEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::
          GetTableBitmapClassificationEntryAttributeRequest* req,
      lemming::dataplane::sai::
          GetTableBitmapClassificationEntryAttributeResponse* resp);

  grpc::Status GetTableBitmapClassificationEntryStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::
          GetTableBitmapClassificationEntryStatsRequest* req,
      lemming::dataplane::sai::GetTableBitmapClassificationEntryStatsResponse*
          resp);

  grpc::Status CreateTableBitmapRouterEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTableBitmapRouterEntryRequest* req,
      lemming::dataplane::sai::CreateTableBitmapRouterEntryResponse* resp);

  grpc::Status RemoveTableBitmapRouterEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTableBitmapRouterEntryRequest* req,
      lemming::dataplane::sai::RemoveTableBitmapRouterEntryResponse* resp);

  grpc::Status GetTableBitmapRouterEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTableBitmapRouterEntryAttributeRequest*
          req,
      lemming::dataplane::sai::GetTableBitmapRouterEntryAttributeResponse*
          resp);

  grpc::Status GetTableBitmapRouterEntryStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTableBitmapRouterEntryStatsRequest* req,
      lemming::dataplane::sai::GetTableBitmapRouterEntryStatsResponse* resp);

  grpc::Status CreateTableMetaTunnelEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTableMetaTunnelEntryRequest* req,
      lemming::dataplane::sai::CreateTableMetaTunnelEntryResponse* resp);

  grpc::Status RemoveTableMetaTunnelEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTableMetaTunnelEntryRequest* req,
      lemming::dataplane::sai::RemoveTableMetaTunnelEntryResponse* resp);

  grpc::Status GetTableMetaTunnelEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTableMetaTunnelEntryAttributeRequest*
          req,
      lemming::dataplane::sai::GetTableMetaTunnelEntryAttributeResponse* resp);

  grpc::Status GetTableMetaTunnelEntryStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTableMetaTunnelEntryStatsRequest* req,
      lemming::dataplane::sai::GetTableMetaTunnelEntryStatsResponse* resp);

  sai_bmtor_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_BMTOR_H_
