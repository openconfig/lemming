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

#include "dataplane/standalone/saiserver/nat.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/nat.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Nat::CreateNatEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNatEntryRequest* req,
    lemming::dataplane::sai::CreateNatEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Nat::RemoveNatEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNatEntryRequest* req,
    lemming::dataplane::sai::RemoveNatEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Nat::SetNatEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetNatEntryAttributeRequest* req,
    lemming::dataplane::sai::SetNatEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Nat::GetNatEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetNatEntryAttributeRequest* req,
    lemming::dataplane::sai::GetNatEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Nat::CreateNatEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNatEntriesRequest* req,
    lemming::dataplane::sai::CreateNatEntriesResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Nat::RemoveNatEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNatEntriesRequest* req,
    lemming::dataplane::sai::RemoveNatEntriesResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Nat::CreateNatZoneCounter(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNatZoneCounterRequest* req,
    lemming::dataplane::sai::CreateNatZoneCounterResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Nat::RemoveNatZoneCounter(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNatZoneCounterRequest* req,
    lemming::dataplane::sai::RemoveNatZoneCounterResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Nat::SetNatZoneCounterAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetNatZoneCounterAttributeRequest* req,
    lemming::dataplane::sai::SetNatZoneCounterAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Nat::GetNatZoneCounterAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetNatZoneCounterAttributeRequest* req,
    lemming::dataplane::sai::GetNatZoneCounterAttributeResponse* resp) {
  return grpc::Status::OK;
}
