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

#ifndef DATAPLANE_STANDALONE_SAI_IPMC_H_
#define DATAPLANE_STANDALONE_SAI_IPMC_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/ipmc.grpc.pb.h"
#include "dataplane/proto/sai/ipmc.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Ipmc final : public lemming::dataplane::sai::Ipmc::Service {
 public:
  grpc::Status CreateIpmcEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateIpmcEntryRequest* req,
      lemming::dataplane::sai::CreateIpmcEntryResponse* resp);

  grpc::Status RemoveIpmcEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveIpmcEntryRequest* req,
      lemming::dataplane::sai::RemoveIpmcEntryResponse* resp);

  grpc::Status SetIpmcEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetIpmcEntryAttributeRequest* req,
      lemming::dataplane::sai::SetIpmcEntryAttributeResponse* resp);

  grpc::Status GetIpmcEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIpmcEntryAttributeRequest* req,
      lemming::dataplane::sai::GetIpmcEntryAttributeResponse* resp);

  sai_ipmc_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_IPMC_H_
