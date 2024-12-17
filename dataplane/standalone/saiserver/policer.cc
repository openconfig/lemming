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

#include "dataplane/standalone/saiserver/policer.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/policer.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Policer::CreatePolicer(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreatePolicerRequest* req,
    lemming::dataplane::sai::CreatePolicerResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Policer::RemovePolicer(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemovePolicerRequest* req,
    lemming::dataplane::sai::RemovePolicerResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_policer(req.get_oid());

  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return grpc::Status::INTERNAL;
  }

  return grpc::Status::OK;
}

grpc::Status Policer::SetPolicerAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetPolicerAttributeRequest* req,
    lemming::dataplane::sai::SetPolicerAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Policer::GetPolicerAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetPolicerAttributeRequest* req,
    lemming::dataplane::sai::GetPolicerAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Policer::GetPolicerStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetPolicerStatsRequest* req,
    lemming::dataplane::sai::GetPolicerStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
