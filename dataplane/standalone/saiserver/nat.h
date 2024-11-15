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

#ifndef DATAPLANE_STANDALONE_SAI_NAT_H_
#define DATAPLANE_STANDALONE_SAI_NAT_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/nat.grpc.pb.h"
#include "dataplane/proto/sai/nat.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Nat final : public lemming::dataplane::sai::Nat::Service {
 public:
  grpc::Status CreateNatEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNatEntryRequest* req,
      lemming::dataplane::sai::CreateNatEntryResponse* resp);

  grpc::Status RemoveNatEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNatEntryRequest* req,
      lemming::dataplane::sai::RemoveNatEntryResponse* resp);

  grpc::Status SetNatEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetNatEntryAttributeRequest* req,
      lemming::dataplane::sai::SetNatEntryAttributeResponse* resp);

  grpc::Status GetNatEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetNatEntryAttributeRequest* req,
      lemming::dataplane::sai::GetNatEntryAttributeResponse* resp);

  grpc::Status CreateNatEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNatEntriesRequest* req,
      lemming::dataplane::sai::CreateNatEntriesResponse* resp);

  grpc::Status RemoveNatEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNatEntriesRequest* req,
      lemming::dataplane::sai::RemoveNatEntriesResponse* resp);

  grpc::Status CreateNatZoneCounter(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNatZoneCounterRequest* req,
      lemming::dataplane::sai::CreateNatZoneCounterResponse* resp);

  grpc::Status RemoveNatZoneCounter(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNatZoneCounterRequest* req,
      lemming::dataplane::sai::RemoveNatZoneCounterResponse* resp);

  grpc::Status SetNatZoneCounterAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetNatZoneCounterAttributeRequest* req,
      lemming::dataplane::sai::SetNatZoneCounterAttributeResponse* resp);

  grpc::Status GetNatZoneCounterAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetNatZoneCounterAttributeRequest* req,
      lemming::dataplane::sai::GetNatZoneCounterAttributeResponse* resp);

  sai_nat_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_NAT_H_
