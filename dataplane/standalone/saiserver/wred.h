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

#ifndef DATAPLANE_STANDALONE_SAI_WRED_H_
#define DATAPLANE_STANDALONE_SAI_WRED_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/wred.grpc.pb.h"
#include "dataplane/proto/sai/wred.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Wred final : public lemming::dataplane::sai::Wred::Service {
 public:
  grpc::Status CreateWred(grpc::ServerContext* context,
                          const lemming::dataplane::sai::CreateWredRequest* req,
                          lemming::dataplane::sai::CreateWredResponse* resp);

  grpc::Status RemoveWred(grpc::ServerContext* context,
                          const lemming::dataplane::sai::RemoveWredRequest* req,
                          lemming::dataplane::sai::RemoveWredResponse* resp);

  grpc::Status SetWredAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetWredAttributeRequest* req,
      lemming::dataplane::sai::SetWredAttributeResponse* resp);

  grpc::Status GetWredAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetWredAttributeRequest* req,
      lemming::dataplane::sai::GetWredAttributeResponse* resp);

  sai_wred_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_WRED_H_
