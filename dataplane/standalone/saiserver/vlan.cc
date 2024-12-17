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

#include "dataplane/standalone/saiserver/vlan.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/vlan.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Vlan::CreateVlan(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateVlanRequest* req,
    lemming::dataplane::sai::CreateVlanResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Vlan::RemoveVlan(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveVlanRequest* req,
    lemming::dataplane::sai::RemoveVlanResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_vlan(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Vlan::SetVlanAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetVlanAttributeRequest* req,
    lemming::dataplane::sai::SetVlanAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Vlan::GetVlanAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetVlanAttributeRequest* req,
    lemming::dataplane::sai::GetVlanAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Vlan::CreateVlanMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateVlanMemberRequest* req,
    lemming::dataplane::sai::CreateVlanMemberResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Vlan::RemoveVlanMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveVlanMemberRequest* req,
    lemming::dataplane::sai::RemoveVlanMemberResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  auto status = api->remove_vlan_member(req->oid());

  if (status != SAI_STATUS_SUCCESS) {
    context->AddTrailingMetadata("status-code", "500");
    context->AddTrailingMetadata("message", "Internal server error");
    return grpc::Status(grpc::StatusCode::INTERNAL, "Internal error occurred");
  }

  return grpc::Status::OK;
}

grpc::Status Vlan::SetVlanMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetVlanMemberAttributeRequest* req,
    lemming::dataplane::sai::SetVlanMemberAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Vlan::GetVlanMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetVlanMemberAttributeRequest* req,
    lemming::dataplane::sai::GetVlanMemberAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Vlan::CreateVlanMembers(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateVlanMembersRequest* req,
    lemming::dataplane::sai::CreateVlanMembersResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Vlan::RemoveVlanMembers(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveVlanMembersRequest* req,
    lemming::dataplane::sai::RemoveVlanMembersResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Vlan::GetVlanStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetVlanStatsRequest* req,
    lemming::dataplane::sai::GetVlanStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
