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

#ifndef DATAPLANE_STANDALONE_SAI_VLAN_H_
#define DATAPLANE_STANDALONE_SAI_VLAN_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/vlan.grpc.pb.h"
#include "dataplane/proto/sai/vlan.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Vlan final : public lemming::dataplane::sai::Vlan::Service {
 public:
  grpc::Status CreateVlan(grpc::ServerContext* context,
                          const lemming::dataplane::sai::CreateVlanRequest* req,
                          lemming::dataplane::sai::CreateVlanResponse* resp);

  grpc::Status RemoveVlan(grpc::ServerContext* context,
                          const lemming::dataplane::sai::RemoveVlanRequest* req,
                          lemming::dataplane::sai::RemoveVlanResponse* resp);

  grpc::Status SetVlanAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetVlanAttributeRequest* req,
      lemming::dataplane::sai::SetVlanAttributeResponse* resp);

  grpc::Status GetVlanAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetVlanAttributeRequest* req,
      lemming::dataplane::sai::GetVlanAttributeResponse* resp);

  grpc::Status CreateVlanMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateVlanMemberRequest* req,
      lemming::dataplane::sai::CreateVlanMemberResponse* resp);

  grpc::Status RemoveVlanMember(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveVlanMemberRequest* req,
      lemming::dataplane::sai::RemoveVlanMemberResponse* resp);

  grpc::Status SetVlanMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetVlanMemberAttributeRequest* req,
      lemming::dataplane::sai::SetVlanMemberAttributeResponse* resp);

  grpc::Status GetVlanMemberAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetVlanMemberAttributeRequest* req,
      lemming::dataplane::sai::GetVlanMemberAttributeResponse* resp);

  grpc::Status CreateVlanMembers(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateVlanMembersRequest* req,
      lemming::dataplane::sai::CreateVlanMembersResponse* resp);

  grpc::Status RemoveVlanMembers(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveVlanMembersRequest* req,
      lemming::dataplane::sai::RemoveVlanMembersResponse* resp);

  grpc::Status GetVlanStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetVlanStatsRequest* req,
      lemming::dataplane::sai::GetVlanStatsResponse* resp);

  sai_vlan_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_VLAN_H_
