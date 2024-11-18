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

#include "dataplane/standalone/saiserver/macsec.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/macsec.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Macsec::CreateMacsec(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateMacsecRequest* req,
    lemming::dataplane::sai::CreateMacsecResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::RemoveMacsec(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveMacsecRequest* req,
    lemming::dataplane::sai::RemoveMacsecResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::SetMacsecAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetMacsecAttributeRequest* req,
    lemming::dataplane::sai::SetMacsecAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::GetMacsecAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMacsecAttributeRequest* req,
    lemming::dataplane::sai::GetMacsecAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::CreateMacsecPort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateMacsecPortRequest* req,
    lemming::dataplane::sai::CreateMacsecPortResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::RemoveMacsecPort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveMacsecPortRequest* req,
    lemming::dataplane::sai::RemoveMacsecPortResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::SetMacsecPortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetMacsecPortAttributeRequest* req,
    lemming::dataplane::sai::SetMacsecPortAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::GetMacsecPortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMacsecPortAttributeRequest* req,
    lemming::dataplane::sai::GetMacsecPortAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::GetMacsecPortStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMacsecPortStatsRequest* req,
    lemming::dataplane::sai::GetMacsecPortStatsResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::CreateMacsecFlow(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateMacsecFlowRequest* req,
    lemming::dataplane::sai::CreateMacsecFlowResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::RemoveMacsecFlow(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveMacsecFlowRequest* req,
    lemming::dataplane::sai::RemoveMacsecFlowResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::GetMacsecFlowAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMacsecFlowAttributeRequest* req,
    lemming::dataplane::sai::GetMacsecFlowAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::GetMacsecFlowStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMacsecFlowStatsRequest* req,
    lemming::dataplane::sai::GetMacsecFlowStatsResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::CreateMacsecSc(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateMacsecScRequest* req,
    lemming::dataplane::sai::CreateMacsecScResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::RemoveMacsecSc(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveMacsecScRequest* req,
    lemming::dataplane::sai::RemoveMacsecScResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::SetMacsecScAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetMacsecScAttributeRequest* req,
    lemming::dataplane::sai::SetMacsecScAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::GetMacsecScAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMacsecScAttributeRequest* req,
    lemming::dataplane::sai::GetMacsecScAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::GetMacsecScStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMacsecScStatsRequest* req,
    lemming::dataplane::sai::GetMacsecScStatsResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::CreateMacsecSa(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateMacsecSaRequest* req,
    lemming::dataplane::sai::CreateMacsecSaResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::RemoveMacsecSa(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveMacsecSaRequest* req,
    lemming::dataplane::sai::RemoveMacsecSaResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::SetMacsecSaAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetMacsecSaAttributeRequest* req,
    lemming::dataplane::sai::SetMacsecSaAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::GetMacsecSaAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMacsecSaAttributeRequest* req,
    lemming::dataplane::sai::GetMacsecSaAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Macsec::GetMacsecSaStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMacsecSaStatsRequest* req,
    lemming::dataplane::sai::GetMacsecSaStatsResponse* resp) {
  return grpc::Status::OK;
}
