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

#ifndef DATAPLANE_STANDALONE_SAI_BFD_H_
#define DATAPLANE_STANDALONE_SAI_BFD_H_

#include "dataplane/proto/sai/bfd.grpc.pb.h"
#include "dataplane/proto/sai/bfd.pb.h"
#include "dataplane/proto/sai/common.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Bfd final : public lemming::dataplane::sai::Bfd::Service {
 public:
  grpc::Status CreateBfdSession(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateBfdSessionRequest* req,
      lemming::dataplane::sai::CreateBfdSessionResponse* resp);

  grpc::Status RemoveBfdSession(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveBfdSessionRequest* req,
      lemming::dataplane::sai::RemoveBfdSessionResponse* resp);

  grpc::Status SetBfdSessionAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetBfdSessionAttributeRequest* req,
      lemming::dataplane::sai::SetBfdSessionAttributeResponse* resp);

  grpc::Status GetBfdSessionAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetBfdSessionAttributeRequest* req,
      lemming::dataplane::sai::GetBfdSessionAttributeResponse* resp);

  grpc::Status GetBfdSessionStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetBfdSessionStatsRequest* req,
      lemming::dataplane::sai::GetBfdSessionStatsResponse* resp);

  sai_bfd_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_BFD_H_
