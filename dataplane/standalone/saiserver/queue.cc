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

#include "dataplane/standalone/saiserver/queue.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/queue.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Queue::CreateQueue(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateQueueRequest* req,
    lemming::dataplane::sai::CreateQueueResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Queue::RemoveQueue(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveQueueRequest* req,
    lemming::dataplane::sai::RemoveQueueResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Queue::SetQueueAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetQueueAttributeRequest* req,
    lemming::dataplane::sai::SetQueueAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Queue::GetQueueAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetQueueAttributeRequest* req,
    lemming::dataplane::sai::GetQueueAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Queue::GetQueueStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetQueueStatsRequest* req,
    lemming::dataplane::sai::GetQueueStatsResponse* resp) {
  return grpc::Status::OK;
}
