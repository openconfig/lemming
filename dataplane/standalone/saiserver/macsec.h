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

#ifndef DATAPLANE_STANDALONE_SAI_MACSEC_H_
#define DATAPLANE_STANDALONE_SAI_MACSEC_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/macsec.grpc.pb.h"
#include "dataplane/proto/sai/macsec.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Macsec final : public lemming::dataplane::sai::Macsec::Service {
 public:
  grpc::Status CreateMacsec(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateMacsecRequest* req,
      lemming::dataplane::sai::CreateMacsecResponse* resp);

  grpc::Status RemoveMacsec(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveMacsecRequest* req,
      lemming::dataplane::sai::RemoveMacsecResponse* resp);

  grpc::Status SetMacsecAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetMacsecAttributeRequest* req,
      lemming::dataplane::sai::SetMacsecAttributeResponse* resp);

  grpc::Status GetMacsecAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMacsecAttributeRequest* req,
      lemming::dataplane::sai::GetMacsecAttributeResponse* resp);

  grpc::Status CreateMacsecPort(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateMacsecPortRequest* req,
      lemming::dataplane::sai::CreateMacsecPortResponse* resp);

  grpc::Status RemoveMacsecPort(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveMacsecPortRequest* req,
      lemming::dataplane::sai::RemoveMacsecPortResponse* resp);

  grpc::Status SetMacsecPortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetMacsecPortAttributeRequest* req,
      lemming::dataplane::sai::SetMacsecPortAttributeResponse* resp);

  grpc::Status GetMacsecPortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMacsecPortAttributeRequest* req,
      lemming::dataplane::sai::GetMacsecPortAttributeResponse* resp);

  grpc::Status GetMacsecPortStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMacsecPortStatsRequest* req,
      lemming::dataplane::sai::GetMacsecPortStatsResponse* resp);

  grpc::Status CreateMacsecFlow(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateMacsecFlowRequest* req,
      lemming::dataplane::sai::CreateMacsecFlowResponse* resp);

  grpc::Status RemoveMacsecFlow(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveMacsecFlowRequest* req,
      lemming::dataplane::sai::RemoveMacsecFlowResponse* resp);

  grpc::Status GetMacsecFlowAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMacsecFlowAttributeRequest* req,
      lemming::dataplane::sai::GetMacsecFlowAttributeResponse* resp);

  grpc::Status GetMacsecFlowStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMacsecFlowStatsRequest* req,
      lemming::dataplane::sai::GetMacsecFlowStatsResponse* resp);

  grpc::Status CreateMacsecSc(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateMacsecScRequest* req,
      lemming::dataplane::sai::CreateMacsecScResponse* resp);

  grpc::Status RemoveMacsecSc(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveMacsecScRequest* req,
      lemming::dataplane::sai::RemoveMacsecScResponse* resp);

  grpc::Status SetMacsecScAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetMacsecScAttributeRequest* req,
      lemming::dataplane::sai::SetMacsecScAttributeResponse* resp);

  grpc::Status GetMacsecScAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMacsecScAttributeRequest* req,
      lemming::dataplane::sai::GetMacsecScAttributeResponse* resp);

  grpc::Status GetMacsecScStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMacsecScStatsRequest* req,
      lemming::dataplane::sai::GetMacsecScStatsResponse* resp);

  grpc::Status CreateMacsecSa(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateMacsecSaRequest* req,
      lemming::dataplane::sai::CreateMacsecSaResponse* resp);

  grpc::Status RemoveMacsecSa(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveMacsecSaRequest* req,
      lemming::dataplane::sai::RemoveMacsecSaResponse* resp);

  grpc::Status SetMacsecSaAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetMacsecSaAttributeRequest* req,
      lemming::dataplane::sai::SetMacsecSaAttributeResponse* resp);

  grpc::Status GetMacsecSaAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMacsecSaAttributeRequest* req,
      lemming::dataplane::sai::GetMacsecSaAttributeResponse* resp);

  grpc::Status GetMacsecSaStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMacsecSaStatsRequest* req,
      lemming::dataplane::sai::GetMacsecSaStatsResponse* resp);

  sai_macsec_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_MACSEC_H_
