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

#ifndef DATAPLANE_STANDALONE_SAI_IPSEC_H_
#define DATAPLANE_STANDALONE_SAI_IPSEC_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/ipsec.grpc.pb.h"
#include "dataplane/proto/sai/ipsec.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Ipsec final : public lemming::dataplane::sai::Ipsec::Service {
 public:
  grpc::Status CreateIpsec(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateIpsecRequest* req,
      lemming::dataplane::sai::CreateIpsecResponse* resp);

  grpc::Status RemoveIpsec(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveIpsecRequest* req,
      lemming::dataplane::sai::RemoveIpsecResponse* resp);

  grpc::Status SetIpsecAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetIpsecAttributeRequest* req,
      lemming::dataplane::sai::SetIpsecAttributeResponse* resp);

  grpc::Status GetIpsecAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIpsecAttributeRequest* req,
      lemming::dataplane::sai::GetIpsecAttributeResponse* resp);

  grpc::Status CreateIpsecPort(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateIpsecPortRequest* req,
      lemming::dataplane::sai::CreateIpsecPortResponse* resp);

  grpc::Status RemoveIpsecPort(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveIpsecPortRequest* req,
      lemming::dataplane::sai::RemoveIpsecPortResponse* resp);

  grpc::Status SetIpsecPortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetIpsecPortAttributeRequest* req,
      lemming::dataplane::sai::SetIpsecPortAttributeResponse* resp);

  grpc::Status GetIpsecPortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIpsecPortAttributeRequest* req,
      lemming::dataplane::sai::GetIpsecPortAttributeResponse* resp);

  grpc::Status GetIpsecPortStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIpsecPortStatsRequest* req,
      lemming::dataplane::sai::GetIpsecPortStatsResponse* resp);

  grpc::Status CreateIpsecSa(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateIpsecSaRequest* req,
      lemming::dataplane::sai::CreateIpsecSaResponse* resp);

  grpc::Status RemoveIpsecSa(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveIpsecSaRequest* req,
      lemming::dataplane::sai::RemoveIpsecSaResponse* resp);

  grpc::Status SetIpsecSaAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetIpsecSaAttributeRequest* req,
      lemming::dataplane::sai::SetIpsecSaAttributeResponse* resp);

  grpc::Status GetIpsecSaAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIpsecSaAttributeRequest* req,
      lemming::dataplane::sai::GetIpsecSaAttributeResponse* resp);

  grpc::Status GetIpsecSaStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetIpsecSaStatsRequest* req,
      lemming::dataplane::sai::GetIpsecSaStatsResponse* resp);

  sai_ipsec_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_IPSEC_H_
