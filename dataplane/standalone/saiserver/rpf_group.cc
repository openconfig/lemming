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

#include "dataplane/standalone/saiserver/rpf_group.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/rpf_group.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status RpfGroup::CreateRpfGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateRpfGroupRequest* req,
    lemming::dataplane::sai::CreateRpfGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status RpfGroup::RemoveRpfGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveRpfGroupRequest* req,
    lemming::dataplane::sai::RemoveRpfGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status RpfGroup::GetRpfGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetRpfGroupAttributeRequest* req,
    lemming::dataplane::sai::GetRpfGroupAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status RpfGroup::CreateRpfGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateRpfGroupMemberRequest* req,
    lemming::dataplane::sai::CreateRpfGroupMemberResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status RpfGroup::RemoveRpfGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveRpfGroupMemberRequest* req,
    lemming::dataplane::sai::RemoveRpfGroupMemberResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status RpfGroup::GetRpfGroupMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetRpfGroupMemberAttributeRequest* req,
    lemming::dataplane::sai::GetRpfGroupMemberAttributeResponse* resp) {
  return grpc::Status::OK;
}
