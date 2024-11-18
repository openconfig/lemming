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

#include "dataplane/standalone/saiserver/acl.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/acl.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Acl::CreateAclTable(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateAclTableRequest* req,
    lemming::dataplane::sai::CreateAclTableResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::RemoveAclTable(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveAclTableRequest* req,
    lemming::dataplane::sai::RemoveAclTableResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::GetAclTableAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetAclTableAttributeRequest* req,
    lemming::dataplane::sai::GetAclTableAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::CreateAclEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateAclEntryRequest* req,
    lemming::dataplane::sai::CreateAclEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::RemoveAclEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveAclEntryRequest* req,
    lemming::dataplane::sai::RemoveAclEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::SetAclEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetAclEntryAttributeRequest* req,
    lemming::dataplane::sai::SetAclEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::GetAclEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetAclEntryAttributeRequest* req,
    lemming::dataplane::sai::GetAclEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::CreateAclCounter(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateAclCounterRequest* req,
    lemming::dataplane::sai::CreateAclCounterResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::RemoveAclCounter(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveAclCounterRequest* req,
    lemming::dataplane::sai::RemoveAclCounterResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::SetAclCounterAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetAclCounterAttributeRequest* req,
    lemming::dataplane::sai::SetAclCounterAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::GetAclCounterAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetAclCounterAttributeRequest* req,
    lemming::dataplane::sai::GetAclCounterAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::CreateAclRange(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateAclRangeRequest* req,
    lemming::dataplane::sai::CreateAclRangeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::RemoveAclRange(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveAclRangeRequest* req,
    lemming::dataplane::sai::RemoveAclRangeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::GetAclRangeAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetAclRangeAttributeRequest* req,
    lemming::dataplane::sai::GetAclRangeAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::CreateAclTableGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateAclTableGroupRequest* req,
    lemming::dataplane::sai::CreateAclTableGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::RemoveAclTableGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveAclTableGroupRequest* req,
    lemming::dataplane::sai::RemoveAclTableGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::GetAclTableGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetAclTableGroupAttributeRequest* req,
    lemming::dataplane::sai::GetAclTableGroupAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::CreateAclTableGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateAclTableGroupMemberRequest* req,
    lemming::dataplane::sai::CreateAclTableGroupMemberResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::RemoveAclTableGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveAclTableGroupMemberRequest* req,
    lemming::dataplane::sai::RemoveAclTableGroupMemberResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Acl::GetAclTableGroupMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetAclTableGroupMemberAttributeRequest* req,
    lemming::dataplane::sai::GetAclTableGroupMemberAttributeResponse* resp) {
  return grpc::Status::OK;
}
