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

#include "dataplane/standalone/saiserver/debug_counter.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/debug_counter.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status DebugCounter::CreateDebugCounter(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateDebugCounterRequest* req,
    lemming::dataplane::sai::CreateDebugCounterResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status DebugCounter::RemoveDebugCounter(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveDebugCounterRequest* req,
    lemming::dataplane::sai::RemoveDebugCounterResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status DebugCounter::SetDebugCounterAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetDebugCounterAttributeRequest* req,
    lemming::dataplane::sai::SetDebugCounterAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status DebugCounter::GetDebugCounterAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetDebugCounterAttributeRequest* req,
    lemming::dataplane::sai::GetDebugCounterAttributeResponse* resp) {
  return grpc::Status::OK;
}
