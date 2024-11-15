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

#ifndef DATAPLANE_STANDALONE_SAI_COUNTER_H_
#define DATAPLANE_STANDALONE_SAI_COUNTER_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/counter.grpc.pb.h"
#include "dataplane/proto/sai/counter.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Counter final : public lemming::dataplane::sai::Counter::Service {
 public:
  grpc::Status CreateCounter(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateCounterRequest* req,
      lemming::dataplane::sai::CreateCounterResponse* resp);

  grpc::Status RemoveCounter(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveCounterRequest* req,
      lemming::dataplane::sai::RemoveCounterResponse* resp);

  grpc::Status SetCounterAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetCounterAttributeRequest* req,
      lemming::dataplane::sai::SetCounterAttributeResponse* resp);

  grpc::Status GetCounterAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetCounterAttributeRequest* req,
      lemming::dataplane::sai::GetCounterAttributeResponse* resp);

  grpc::Status GetCounterStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetCounterStatsRequest* req,
      lemming::dataplane::sai::GetCounterStatsResponse* resp);

  sai_counter_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_COUNTER_H_
