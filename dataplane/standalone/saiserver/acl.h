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

#ifndef DATAPLANE_STANDALONE_SAI_ACL_H_
#define DATAPLANE_STANDALONE_SAI_ACL_H_

#include "dataplane/proto/sai/acl.grpc.pb.h"
#include "dataplane/proto/sai/acl.pb.h"
#include "dataplane/proto/sai/common.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Acl final : public lemming::dataplane::sai::Acl::Service {
 public:
  grpc::Status CreateAclTable(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateAclTableRequest* req,
      lemming::dataplane::sai::CreateAclTableResponse* resp);

  grpc::Status RemoveAclTable(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveAclTableRequest* req,
      lemming::dataplane::sai::RemoveAclTableResponse* resp);

  grpc::Status GetAclTableAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetAclTableAttributeRequest* req,
      lemming::dataplane::sai::GetAclTableAttributeResponse* resp);

  grpc::Status CreateAclEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateAclEntryRequest* req,
      lemming::dataplane::sai::CreateAclEntryResponse* resp);

  grpc::Status RemoveAclEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveAclEntryRequest* req,
      lemming::dataplane::sai::RemoveAclEntryResponse* resp);

  grpc::Status SetAclEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetAclEntryAttributeRequest* req,
      lemming::dataplane::sai::SetAclEntryAttributeResponse* resp);

  grpc::Status GetAclEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetAclEntryAttributeRequest* req,
      lemming::dataplane::sai::GetAclEntryAttributeResponse* resp);

  grpc::Status CreateAclCounter(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateAclCounterRequest* req,
      lemming::dataplane::sai::CreateAclCounterResponse* resp);

  grpc::Status RemoveAclCounter(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveAclCounterRequest* req,
      lemming::dataplane::sai::RemoveAclCounterResponse* resp);

  grpc::Status SetAclCounterAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetAclCounterAttributeRequest* req,
      lemming::dataplane::sai::SetAclCounterAttributeResponse* resp);

  grpc::Status GetAclCounterAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetAclCounterAttributeRequest* req,
      lemming::dataplane::sai::GetAclCounterAttributeResponse* resp);

  grpc::Status CreateAclRange(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateAclRangeRequest* req,
      lemming::dataplane::sai::CreateAclRangeResponse* resp);

  grpc::Status RemoveAclRange(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveAclRangeRequest* req,
      lemming::dataplane::sai::RemoveAclRangeResponse* resp);

  grpc::Status GetAclRangeAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetAclRangeAttributeRequest* req,
      lemming::dataplane::sai::GetAclRangeAttributeResponse* resp);

  grpc::Status CreateAclTableGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateAclTableGroupRequest* req,
      lemming::dataplane::sai::CreateAclTableGroupResponse* resp);

  grpc::Status RemoveAclTableGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveAclTableGroupRequest* req,
      lemming::dataplane::sai::RemoveAclTableGroupResponse* resp);

  grpc::Status GetAclTableGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetAclTableGroupAttributeRequest* req,
      lemming::dataplane::sai::GetAclTableGroupAttributeResponse* resp);

  grpc::Status CreateAclTableGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateAclTableGroupMemberRequest* req,
      lemming::dataplane::sai::CreateAclTableGroupMemberResponse* resp);

  grpc::Status RemoveAclTableGroupMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveAclTableGroupMemberRequest* req,
      lemming::dataplane::sai::RemoveAclTableGroupMemberResponse* resp);

  grpc::Status GetAclTableGroupMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetAclTableGroupMemberAttributeRequest*
          req,
      lemming::dataplane::sai::GetAclTableGroupMemberAttributeResponse* resp);

  sai_acl_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_ACL_H_
