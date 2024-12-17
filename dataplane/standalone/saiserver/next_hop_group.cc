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

#include "dataplane/standalone/saiserver/next_hop_group.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/next_hop_group.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status NextHopGroup::CreateNextHopGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNextHopGroupRequest* req,
    lemming::dataplane::sai::CreateNextHopGroupResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::RemoveNextHopGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNextHopGroupRequest* req,
    lemming::dataplane::sai::RemoveNextHopGroupResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_next_hop_group(req.get_oid());

  auto status = api->remove_next_hop_group(entry);
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

grpc::Status NextHopGroup::SetNextHopGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetNextHopGroupAttributeRequest* req,
    lemming::dataplane::sai::SetNextHopGroupAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::GetNextHopGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetNextHopGroupAttributeRequest* req,
    lemming::dataplane::sai::GetNextHopGroupAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::CreateNextHopGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNextHopGroupMemberRequest* req,
    lemming::dataplane::sai::CreateNextHopGroupMemberResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::RemoveNextHopGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNextHopGroupMemberRequest* req,
    lemming::dataplane::sai::RemoveNextHopGroupMemberResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_next_hop_group_member(req.get_oid());

  auto status = api->remove_next_hop_group_member(entry);
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

grpc::Status NextHopGroup::SetNextHopGroupMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetNextHopGroupMemberAttributeRequest* req,
    lemming::dataplane::sai::SetNextHopGroupMemberAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::GetNextHopGroupMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetNextHopGroupMemberAttributeRequest* req,
    lemming::dataplane::sai::GetNextHopGroupMemberAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::CreateNextHopGroupMembers(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNextHopGroupMembersRequest* req,
    lemming::dataplane::sai::CreateNextHopGroupMembersResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::RemoveNextHopGroupMembers(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNextHopGroupMembersRequest* req,
    lemming::dataplane::sai::RemoveNextHopGroupMembersResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::CreateNextHopGroupMap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNextHopGroupMapRequest* req,
    lemming::dataplane::sai::CreateNextHopGroupMapResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::RemoveNextHopGroupMap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNextHopGroupMapRequest* req,
    lemming::dataplane::sai::RemoveNextHopGroupMapResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_next_hop_group_map(req.get_oid());

  auto status = api->remove_next_hop_group_map(entry);
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

grpc::Status NextHopGroup::SetNextHopGroupMapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetNextHopGroupMapAttributeRequest* req,
    lemming::dataplane::sai::SetNextHopGroupMapAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::GetNextHopGroupMapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetNextHopGroupMapAttributeRequest* req,
    lemming::dataplane::sai::GetNextHopGroupMapAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::CreateNextHopGroups(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNextHopGroupsRequest* req,
    lemming::dataplane::sai::CreateNextHopGroupsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status NextHopGroup::RemoveNextHopGroups(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNextHopGroupsRequest* req,
    lemming::dataplane::sai::RemoveNextHopGroupsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
