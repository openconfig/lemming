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

#include "dataplane/standalone/saiserver/scheduler_group.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/scheduler_group.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status SchedulerGroup::CreateSchedulerGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateSchedulerGroupRequest* req,
    lemming::dataplane::sai::CreateSchedulerGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status SchedulerGroup::RemoveSchedulerGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveSchedulerGroupRequest* req,
    lemming::dataplane::sai::RemoveSchedulerGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status SchedulerGroup::SetSchedulerGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetSchedulerGroupAttributeRequest* req,
    lemming::dataplane::sai::SetSchedulerGroupAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status SchedulerGroup::GetSchedulerGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetSchedulerGroupAttributeRequest* req,
    lemming::dataplane::sai::GetSchedulerGroupAttributeResponse* resp) {
  return grpc::Status::OK;
}
