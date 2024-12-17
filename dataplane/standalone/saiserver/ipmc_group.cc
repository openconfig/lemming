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

#include "dataplane/standalone/saiserver/ipmc_group.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/ipmc_group.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status IpmcGroup::CreateIpmcGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateIpmcGroupRequest* req,
    lemming::dataplane::sai::CreateIpmcGroupResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status IpmcGroup::RemoveIpmcGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveIpmcGroupRequest* req,
    lemming::dataplane::sai::RemoveIpmcGroupResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_ipmc_group(req.get_oid());

  auto status = api->remove_ipmc_group(entry);
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

grpc::Status IpmcGroup::GetIpmcGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIpmcGroupAttributeRequest* req,
    lemming::dataplane::sai::GetIpmcGroupAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status IpmcGroup::CreateIpmcGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateIpmcGroupMemberRequest* req,
    lemming::dataplane::sai::CreateIpmcGroupMemberResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status IpmcGroup::RemoveIpmcGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveIpmcGroupMemberRequest* req,
    lemming::dataplane::sai::RemoveIpmcGroupMemberResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_ipmc_group_member(req.get_oid());

  auto status = api->remove_ipmc_group_member(entry);
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

grpc::Status IpmcGroup::GetIpmcGroupMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIpmcGroupMemberAttributeRequest* req,
    lemming::dataplane::sai::GetIpmcGroupMemberAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
