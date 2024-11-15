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

#include "dataplane/standalone/saiserver/hostif.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/hostif.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Hostif::CreateHostif(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHostifRequest* req,
    lemming::dataplane::sai::CreateHostifResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::RemoveHostif(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHostifRequest* req,
    lemming::dataplane::sai::RemoveHostifResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::SetHostifAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetHostifAttributeRequest* req,
    lemming::dataplane::sai::SetHostifAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::GetHostifAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHostifAttributeRequest* req,
    lemming::dataplane::sai::GetHostifAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::CreateHostifTableEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHostifTableEntryRequest* req,
    lemming::dataplane::sai::CreateHostifTableEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::RemoveHostifTableEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHostifTableEntryRequest* req,
    lemming::dataplane::sai::RemoveHostifTableEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::GetHostifTableEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHostifTableEntryAttributeRequest* req,
    lemming::dataplane::sai::GetHostifTableEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::CreateHostifTrapGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHostifTrapGroupRequest* req,
    lemming::dataplane::sai::CreateHostifTrapGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::RemoveHostifTrapGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHostifTrapGroupRequest* req,
    lemming::dataplane::sai::RemoveHostifTrapGroupResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::SetHostifTrapGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetHostifTrapGroupAttributeRequest* req,
    lemming::dataplane::sai::SetHostifTrapGroupAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::GetHostifTrapGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHostifTrapGroupAttributeRequest* req,
    lemming::dataplane::sai::GetHostifTrapGroupAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::CreateHostifTrap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHostifTrapRequest* req,
    lemming::dataplane::sai::CreateHostifTrapResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::RemoveHostifTrap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHostifTrapRequest* req,
    lemming::dataplane::sai::RemoveHostifTrapResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::SetHostifTrapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetHostifTrapAttributeRequest* req,
    lemming::dataplane::sai::SetHostifTrapAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::GetHostifTrapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHostifTrapAttributeRequest* req,
    lemming::dataplane::sai::GetHostifTrapAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::CreateHostifUserDefinedTrap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHostifUserDefinedTrapRequest* req,
    lemming::dataplane::sai::CreateHostifUserDefinedTrapResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::RemoveHostifUserDefinedTrap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHostifUserDefinedTrapRequest* req,
    lemming::dataplane::sai::RemoveHostifUserDefinedTrapResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::SetHostifUserDefinedTrapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetHostifUserDefinedTrapAttributeRequest*
        req,
    lemming::dataplane::sai::SetHostifUserDefinedTrapAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Hostif::GetHostifUserDefinedTrapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHostifUserDefinedTrapAttributeRequest*
        req,
    lemming::dataplane::sai::GetHostifUserDefinedTrapAttributeResponse* resp) {
  return grpc::Status::OK;
}
