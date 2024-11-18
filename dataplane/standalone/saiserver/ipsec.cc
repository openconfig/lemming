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

#include "dataplane/standalone/saiserver/ipsec.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/ipsec.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Ipsec::CreateIpsec(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateIpsecRequest* req,
    lemming::dataplane::sai::CreateIpsecResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::RemoveIpsec(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveIpsecRequest* req,
    lemming::dataplane::sai::RemoveIpsecResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::SetIpsecAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetIpsecAttributeRequest* req,
    lemming::dataplane::sai::SetIpsecAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::GetIpsecAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIpsecAttributeRequest* req,
    lemming::dataplane::sai::GetIpsecAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::CreateIpsecPort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateIpsecPortRequest* req,
    lemming::dataplane::sai::CreateIpsecPortResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::RemoveIpsecPort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveIpsecPortRequest* req,
    lemming::dataplane::sai::RemoveIpsecPortResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::SetIpsecPortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetIpsecPortAttributeRequest* req,
    lemming::dataplane::sai::SetIpsecPortAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::GetIpsecPortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIpsecPortAttributeRequest* req,
    lemming::dataplane::sai::GetIpsecPortAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::GetIpsecPortStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIpsecPortStatsRequest* req,
    lemming::dataplane::sai::GetIpsecPortStatsResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::CreateIpsecSa(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateIpsecSaRequest* req,
    lemming::dataplane::sai::CreateIpsecSaResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::RemoveIpsecSa(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveIpsecSaRequest* req,
    lemming::dataplane::sai::RemoveIpsecSaResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::SetIpsecSaAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetIpsecSaAttributeRequest* req,
    lemming::dataplane::sai::SetIpsecSaAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::GetIpsecSaAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIpsecSaAttributeRequest* req,
    lemming::dataplane::sai::GetIpsecSaAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipsec::GetIpsecSaStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIpsecSaStatsRequest* req,
    lemming::dataplane::sai::GetIpsecSaStatsResponse* resp) {
  return grpc::Status::OK;
}
