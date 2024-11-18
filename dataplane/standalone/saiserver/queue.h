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

#ifndef DATAPLANE_STANDALONE_SAI_QUEUE_H_
#define DATAPLANE_STANDALONE_SAI_QUEUE_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/queue.grpc.pb.h"
#include "dataplane/proto/sai/queue.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Queue final : public lemming::dataplane::sai::Queue::Service {
 public:
  grpc::Status CreateQueue(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateQueueRequest* req,
      lemming::dataplane::sai::CreateQueueResponse* resp);

  grpc::Status RemoveQueue(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveQueueRequest* req,
      lemming::dataplane::sai::RemoveQueueResponse* resp);

  grpc::Status SetQueueAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetQueueAttributeRequest* req,
      lemming::dataplane::sai::SetQueueAttributeResponse* resp);

  grpc::Status GetQueueAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetQueueAttributeRequest* req,
      lemming::dataplane::sai::GetQueueAttributeResponse* resp);

  grpc::Status GetQueueStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetQueueStatsRequest* req,
      lemming::dataplane::sai::GetQueueStatsResponse* resp);

  sai_queue_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_QUEUE_H_
