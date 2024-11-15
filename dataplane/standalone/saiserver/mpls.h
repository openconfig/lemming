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

#ifndef DATAPLANE_STANDALONE_SAI_MPLS_H_
#define DATAPLANE_STANDALONE_SAI_MPLS_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/mpls.grpc.pb.h"
#include "dataplane/proto/sai/mpls.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Mpls final : public lemming::dataplane::sai::Mpls::Service {
 public:
  grpc::Status CreateInsegEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateInsegEntryRequest* req,
      lemming::dataplane::sai::CreateInsegEntryResponse* resp);

  grpc::Status RemoveInsegEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveInsegEntryRequest* req,
      lemming::dataplane::sai::RemoveInsegEntryResponse* resp);

  grpc::Status SetInsegEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetInsegEntryAttributeRequest* req,
      lemming::dataplane::sai::SetInsegEntryAttributeResponse* resp);

  grpc::Status GetInsegEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetInsegEntryAttributeRequest* req,
      lemming::dataplane::sai::GetInsegEntryAttributeResponse* resp);

  grpc::Status CreateInsegEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateInsegEntriesRequest* req,
      lemming::dataplane::sai::CreateInsegEntriesResponse* resp);

  grpc::Status RemoveInsegEntries(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveInsegEntriesRequest* req,
      lemming::dataplane::sai::RemoveInsegEntriesResponse* resp);

  sai_mpls_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_MPLS_H_
