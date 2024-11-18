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

#include "dataplane/standalone/saiserver/isolation_group.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/isolation_group.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status IsolationGroup::CreateIsolationGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateIsolationGroupRequest* req,
    lemming::dataplane::sai::CreateIsolationGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status IsolationGroup::RemoveIsolationGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveIsolationGroupRequest* req,
    lemming::dataplane::sai::RemoveIsolationGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status IsolationGroup::GetIsolationGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIsolationGroupAttributeRequest* req,
    lemming::dataplane::sai::GetIsolationGroupAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status IsolationGroup::CreateIsolationGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateIsolationGroupMemberRequest* req,
    lemming::dataplane::sai::CreateIsolationGroupMemberResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status IsolationGroup::RemoveIsolationGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveIsolationGroupMemberRequest* req,
    lemming::dataplane::sai::RemoveIsolationGroupMemberResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status IsolationGroup::GetIsolationGroupMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIsolationGroupMemberAttributeRequest* req,
    lemming::dataplane::sai::GetIsolationGroupMemberAttributeResponse* resp) {
  return grpc::Status::OK;
}
