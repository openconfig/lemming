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

#include "dataplane/standalone/saiserver/l2mc_group.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/l2mc_group.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status L2mcGroup::CreateL2mcGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateL2mcGroupRequest* req,
    lemming::dataplane::sai::CreateL2mcGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status L2mcGroup::RemoveL2mcGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveL2mcGroupRequest* req,
    lemming::dataplane::sai::RemoveL2mcGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status L2mcGroup::GetL2mcGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetL2mcGroupAttributeRequest* req,
    lemming::dataplane::sai::GetL2mcGroupAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status L2mcGroup::CreateL2mcGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateL2mcGroupMemberRequest* req,
    lemming::dataplane::sai::CreateL2mcGroupMemberResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status L2mcGroup::RemoveL2mcGroupMember(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveL2mcGroupMemberRequest* req,
    lemming::dataplane::sai::RemoveL2mcGroupMemberResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status L2mcGroup::GetL2mcGroupMemberAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetL2mcGroupMemberAttributeRequest* req,
    lemming::dataplane::sai::GetL2mcGroupMemberAttributeResponse* resp) {
  return grpc::Status::OK;
}
