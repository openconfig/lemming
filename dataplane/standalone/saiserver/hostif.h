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

#ifndef DATAPLANE_STANDALONE_SAI_HOSTIF_H_
#define DATAPLANE_STANDALONE_SAI_HOSTIF_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/hostif.grpc.pb.h"
#include "dataplane/proto/sai/hostif.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Hostif final : public lemming::dataplane::sai::Hostif::Service {
 public:
  grpc::Status CreateHostif(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateHostifRequest* req,
      lemming::dataplane::sai::CreateHostifResponse* resp);

  grpc::Status RemoveHostif(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveHostifRequest* req,
      lemming::dataplane::sai::RemoveHostifResponse* resp);

  grpc::Status SetHostifAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetHostifAttributeRequest* req,
      lemming::dataplane::sai::SetHostifAttributeResponse* resp);

  grpc::Status GetHostifAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetHostifAttributeRequest* req,
      lemming::dataplane::sai::GetHostifAttributeResponse* resp);

  grpc::Status CreateHostifTableEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateHostifTableEntryRequest* req,
      lemming::dataplane::sai::CreateHostifTableEntryResponse* resp);

  grpc::Status RemoveHostifTableEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveHostifTableEntryRequest* req,
      lemming::dataplane::sai::RemoveHostifTableEntryResponse* resp);

  grpc::Status GetHostifTableEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetHostifTableEntryAttributeRequest* req,
      lemming::dataplane::sai::GetHostifTableEntryAttributeResponse* resp);

  grpc::Status CreateHostifTrapGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateHostifTrapGroupRequest* req,
      lemming::dataplane::sai::CreateHostifTrapGroupResponse* resp);

  grpc::Status RemoveHostifTrapGroup(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveHostifTrapGroupRequest* req,
      lemming::dataplane::sai::RemoveHostifTrapGroupResponse* resp);

  grpc::Status SetHostifTrapGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetHostifTrapGroupAttributeRequest* req,
      lemming::dataplane::sai::SetHostifTrapGroupAttributeResponse* resp);

  grpc::Status GetHostifTrapGroupAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetHostifTrapGroupAttributeRequest* req,
      lemming::dataplane::sai::GetHostifTrapGroupAttributeResponse* resp);

  grpc::Status CreateHostifTrap(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateHostifTrapRequest* req,
      lemming::dataplane::sai::CreateHostifTrapResponse* resp);

  grpc::Status RemoveHostifTrap(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveHostifTrapRequest* req,
      lemming::dataplane::sai::RemoveHostifTrapResponse* resp);

  grpc::Status SetHostifTrapAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetHostifTrapAttributeRequest* req,
      lemming::dataplane::sai::SetHostifTrapAttributeResponse* resp);

  grpc::Status GetHostifTrapAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetHostifTrapAttributeRequest* req,
      lemming::dataplane::sai::GetHostifTrapAttributeResponse* resp);

  grpc::Status CreateHostifUserDefinedTrap(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateHostifUserDefinedTrapRequest* req,
      lemming::dataplane::sai::CreateHostifUserDefinedTrapResponse* resp);

  grpc::Status RemoveHostifUserDefinedTrap(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveHostifUserDefinedTrapRequest* req,
      lemming::dataplane::sai::RemoveHostifUserDefinedTrapResponse* resp);

  grpc::Status SetHostifUserDefinedTrapAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetHostifUserDefinedTrapAttributeRequest*
          req,
      lemming::dataplane::sai::SetHostifUserDefinedTrapAttributeResponse* resp);

  grpc::Status GetHostifUserDefinedTrapAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetHostifUserDefinedTrapAttributeRequest*
          req,
      lemming::dataplane::sai::GetHostifUserDefinedTrapAttributeResponse* resp);

  sai_hostif_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_HOSTIF_H_
