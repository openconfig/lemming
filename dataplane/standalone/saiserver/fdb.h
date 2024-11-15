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

#ifndef DATAPLANE_STANDALONE_SAI_FDB_H_
#define DATAPLANE_STANDALONE_SAI_FDB_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/fdb.grpc.pb.h"
#include "dataplane/proto/sai/fdb.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Fdb final : public lemming::dataplane::sai::Fdb::Service {
 public:
  grpc::Status CreateFdbEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateFdbEntryRequest* req,
      lemming::dataplane::sai::CreateFdbEntryResponse* resp);

  grpc::Status RemoveFdbEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveFdbEntryRequest* req,
      lemming::dataplane::sai::RemoveFdbEntryResponse* resp);

  grpc::Status SetFdbEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetFdbEntryAttributeRequest* req,
      lemming::dataplane::sai::SetFdbEntryAttributeResponse* resp);

  grpc::Status GetFdbEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetFdbEntryAttributeRequest* req,
      lemming::dataplane::sai::GetFdbEntryAttributeResponse* resp);

  grpc::Status CreateFdbEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateFdbEntriesRequest* req,
      lemming::dataplane::sai::CreateFdbEntriesResponse* resp);

  grpc::Status RemoveFdbEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveFdbEntriesRequest* req,
      lemming::dataplane::sai::RemoveFdbEntriesResponse* resp);

  sai_fdb_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_FDB_H_
