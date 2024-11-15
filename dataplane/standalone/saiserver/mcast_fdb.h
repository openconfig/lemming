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

#ifndef DATAPLANE_STANDALONE_SAI_MCAST_FDB_H_
#define DATAPLANE_STANDALONE_SAI_MCAST_FDB_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/mcast_fdb.grpc.pb.h"
#include "dataplane/proto/sai/mcast_fdb.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class McastFdb final : public lemming::dataplane::sai::McastFdb::Service {
 public:
  grpc::Status CreateMcastFdbEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateMcastFdbEntryRequest* req,
      lemming::dataplane::sai::CreateMcastFdbEntryResponse* resp);

  grpc::Status RemoveMcastFdbEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveMcastFdbEntryRequest* req,
      lemming::dataplane::sai::RemoveMcastFdbEntryResponse* resp);

  grpc::Status SetMcastFdbEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetMcastFdbEntryAttributeRequest* req,
      lemming::dataplane::sai::SetMcastFdbEntryAttributeResponse* resp);

  grpc::Status GetMcastFdbEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMcastFdbEntryAttributeRequest* req,
      lemming::dataplane::sai::GetMcastFdbEntryAttributeResponse* resp);

  sai_mcast_fdb_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_MCAST_FDB_H_
