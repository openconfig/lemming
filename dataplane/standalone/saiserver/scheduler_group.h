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

#ifndef DATAPLANE_STANDALONE_SAI_SCHEDULER_GROUP_H_
#define DATAPLANE_STANDALONE_SAI_SCHEDULER_GROUP_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/scheduler_group.grpc.pb.h"
#include "dataplane/proto/sai/scheduler_group.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class SchedulerGroup final
    : public lemming::dataplane::sai::SchedulerGroup::Service {
 public:
  grpc::Status CreateSchedulerGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateSchedulerGroupRequest* req,
      lemming::dataplane::sai::CreateSchedulerGroupResponse* resp);

  grpc::Status RemoveSchedulerGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveSchedulerGroupRequest* req,
      lemming::dataplane::sai::RemoveSchedulerGroupResponse* resp);

  grpc::Status SetSchedulerGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetSchedulerGroupAttributeRequest* req,
      lemming::dataplane::sai::SetSchedulerGroupAttributeResponse* resp);

  grpc::Status GetSchedulerGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetSchedulerGroupAttributeRequest* req,
      lemming::dataplane::sai::GetSchedulerGroupAttributeResponse* resp);

  sai_scheduler_group_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_SCHEDULER_GROUP_H_
