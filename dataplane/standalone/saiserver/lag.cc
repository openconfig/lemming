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

#include "dataplane/standalone/saiserver/lag.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/lag.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Lag::CreateLag(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateLagRequest* req,
    lemming::dataplane::sai::CreateLagResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Lag::RemoveLag(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveLagRequest* req,
    lemming::dataplane::sai::RemoveLagResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_lag(req.get_oid());

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

grpc::Status Lag::SetLagAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetLagAttributeRequest* req,
    lemming::dataplane::sai::SetLagAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Lag::GetLagAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetLagAttributeRequest* req,
    lemming::dataplane::sai::GetLagAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Lag::CreateLagMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateLagMemberRequest* req,
    lemming::dataplane::sai::CreateLagMemberResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Lag::RemoveLagMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveLagMemberRequest* req,
    lemming::dataplane::sai::RemoveLagMemberResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_lag_member(req.get_oid());

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

grpc::Status Lag::SetLagMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetLagMemberAttributeRequest* req,
    lemming::dataplane::sai::SetLagMemberAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Lag::GetLagMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetLagMemberAttributeRequest* req,
    lemming::dataplane::sai::GetLagMemberAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Lag::CreateLagMembers(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateLagMembersRequest* req,
    lemming::dataplane::sai::CreateLagMembersResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Lag::RemoveLagMembers(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveLagMembersRequest* req,
    lemming::dataplane::sai::RemoveLagMembersResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
