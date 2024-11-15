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

#include "dataplane/standalone/saiserver/scheduler.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/scheduler.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Scheduler::CreateScheduler(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateSchedulerRequest* req,
    lemming::dataplane::sai::CreateSchedulerResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Scheduler::RemoveScheduler(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveSchedulerRequest* req,
    lemming::dataplane::sai::RemoveSchedulerResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Scheduler::SetSchedulerAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetSchedulerAttributeRequest* req,
    lemming::dataplane::sai::SetSchedulerAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Scheduler::GetSchedulerAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetSchedulerAttributeRequest* req,
    lemming::dataplane::sai::GetSchedulerAttributeResponse* resp) {
  return grpc::Status::OK;
}
