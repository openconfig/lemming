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

#ifndef DATAPLANE_STANDALONE_SAI_SRV6_H_
#define DATAPLANE_STANDALONE_SAI_SRV6_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/srv6.grpc.pb.h"
#include "dataplane/proto/sai/srv6.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Srv6 final : public lemming::dataplane::sai::Srv6::Service {
 public:
  grpc::Status CreateSrv6Sidlist(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateSrv6SidlistRequest* req,
      lemming::dataplane::sai::CreateSrv6SidlistResponse* resp);

  grpc::Status RemoveSrv6Sidlist(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveSrv6SidlistRequest* req,
      lemming::dataplane::sai::RemoveSrv6SidlistResponse* resp);

  grpc::Status SetSrv6SidlistAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetSrv6SidlistAttributeRequest* req,
      lemming::dataplane::sai::SetSrv6SidlistAttributeResponse* resp);

  grpc::Status GetSrv6SidlistAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetSrv6SidlistAttributeRequest* req,
      lemming::dataplane::sai::GetSrv6SidlistAttributeResponse* resp);

  grpc::Status CreateSrv6Sidlists(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateSrv6SidlistsRequest* req,
      lemming::dataplane::sai::CreateSrv6SidlistsResponse* resp);

  grpc::Status RemoveSrv6Sidlists(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveSrv6SidlistsRequest* req,
      lemming::dataplane::sai::RemoveSrv6SidlistsResponse* resp);

  grpc::Status CreateMySidEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateMySidEntryRequest* req,
      lemming::dataplane::sai::CreateMySidEntryResponse* resp);

  grpc::Status RemoveMySidEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveMySidEntryRequest* req,
      lemming::dataplane::sai::RemoveMySidEntryResponse* resp);

  grpc::Status SetMySidEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetMySidEntryAttributeRequest* req,
      lemming::dataplane::sai::SetMySidEntryAttributeResponse* resp);

  grpc::Status GetMySidEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMySidEntryAttributeRequest* req,
      lemming::dataplane::sai::GetMySidEntryAttributeResponse* resp);

  grpc::Status CreateMySidEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateMySidEntriesRequest* req,
      lemming::dataplane::sai::CreateMySidEntriesResponse* resp);

  grpc::Status RemoveMySidEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveMySidEntriesRequest* req,
      lemming::dataplane::sai::RemoveMySidEntriesResponse* resp);

  sai_srv6_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_SRV6_H_
