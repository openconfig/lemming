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

#include "dataplane/standalone/saiserver/counter.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/counter.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Counter::CreateCounter(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateCounterRequest* req,
    lemming::dataplane::sai::CreateCounterResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Counter::RemoveCounter(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveCounterRequest* req,
    lemming::dataplane::sai::RemoveCounterResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Counter::SetCounterAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetCounterAttributeRequest* req,
    lemming::dataplane::sai::SetCounterAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Counter::GetCounterAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetCounterAttributeRequest* req,
    lemming::dataplane::sai::GetCounterAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Counter::GetCounterStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetCounterStatsRequest* req,
    lemming::dataplane::sai::GetCounterStatsResponse* resp) {
  return grpc::Status::OK;
}
