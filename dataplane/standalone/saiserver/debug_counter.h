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

#ifndef DATAPLANE_STANDALONE_SAI_DEBUG_COUNTER_H_
#define DATAPLANE_STANDALONE_SAI_DEBUG_COUNTER_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/debug_counter.grpc.pb.h"
#include "dataplane/proto/sai/debug_counter.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class DebugCounter final
    : public lemming::dataplane::sai::DebugCounter::Service {
 public:
  grpc::Status CreateDebugCounter(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateDebugCounterRequest* req,
      lemming::dataplane::sai::CreateDebugCounterResponse* resp);

  grpc::Status RemoveDebugCounter(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveDebugCounterRequest* req,
      lemming::dataplane::sai::RemoveDebugCounterResponse* resp);

  grpc::Status SetDebugCounterAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetDebugCounterAttributeRequest* req,
      lemming::dataplane::sai::SetDebugCounterAttributeResponse* resp);

  grpc::Status GetDebugCounterAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetDebugCounterAttributeRequest* req,
      lemming::dataplane::sai::GetDebugCounterAttributeResponse* resp);

  sai_debug_counter_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_DEBUG_COUNTER_H_
