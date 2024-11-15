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

#ifndef DATAPLANE_STANDALONE_SAI_L2MC_H_
#define DATAPLANE_STANDALONE_SAI_L2MC_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/l2mc.grpc.pb.h"
#include "dataplane/proto/sai/l2mc.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class L2mc final : public lemming::dataplane::sai::L2mc::Service {
 public:
  grpc::Status CreateL2mcEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateL2mcEntryRequest* req,
      lemming::dataplane::sai::CreateL2mcEntryResponse* resp);

  grpc::Status RemoveL2mcEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveL2mcEntryRequest* req,
      lemming::dataplane::sai::RemoveL2mcEntryResponse* resp);

  grpc::Status SetL2mcEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetL2mcEntryAttributeRequest* req,
      lemming::dataplane::sai::SetL2mcEntryAttributeResponse* resp);

  grpc::Status GetL2mcEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetL2mcEntryAttributeRequest* req,
      lemming::dataplane::sai::GetL2mcEntryAttributeResponse* resp);

  sai_l2mc_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_L2MC_H_
