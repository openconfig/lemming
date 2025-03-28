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

#ifndef DATAPLANE_STANDALONE_SAI_SCHEDULER_H_
#define DATAPLANE_STANDALONE_SAI_SCHEDULER_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/scheduler.grpc.pb.h"
#include "dataplane/proto/sai/scheduler.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Scheduler final : public lemming::dataplane::sai::Scheduler::Service {
 public:
  grpc::Status CreateScheduler(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateSchedulerRequest* req,
      lemming::dataplane::sai::CreateSchedulerResponse* resp);

  grpc::Status RemoveScheduler(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveSchedulerRequest* req,
      lemming::dataplane::sai::RemoveSchedulerResponse* resp);

  grpc::Status SetSchedulerAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetSchedulerAttributeRequest* req,
      lemming::dataplane::sai::SetSchedulerAttributeResponse* resp);

  grpc::Status GetSchedulerAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetSchedulerAttributeRequest* req,
      lemming::dataplane::sai::GetSchedulerAttributeResponse* resp);

  sai_scheduler_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_SCHEDULER_H_
