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

#include "dataplane/standalone/saiserver/srv6.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/srv6.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Srv6::CreateSrv6Sidlist(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateSrv6SidlistRequest* req,
    lemming::dataplane::sai::CreateSrv6SidlistResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::RemoveSrv6Sidlist(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveSrv6SidlistRequest* req,
    lemming::dataplane::sai::RemoveSrv6SidlistResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::SetSrv6SidlistAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetSrv6SidlistAttributeRequest* req,
    lemming::dataplane::sai::SetSrv6SidlistAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::GetSrv6SidlistAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetSrv6SidlistAttributeRequest* req,
    lemming::dataplane::sai::GetSrv6SidlistAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::CreateSrv6Sidlists(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateSrv6SidlistsRequest* req,
    lemming::dataplane::sai::CreateSrv6SidlistsResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::RemoveSrv6Sidlists(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveSrv6SidlistsRequest* req,
    lemming::dataplane::sai::RemoveSrv6SidlistsResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::CreateMySidEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateMySidEntryRequest* req,
    lemming::dataplane::sai::CreateMySidEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::RemoveMySidEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveMySidEntryRequest* req,
    lemming::dataplane::sai::RemoveMySidEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::SetMySidEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetMySidEntryAttributeRequest* req,
    lemming::dataplane::sai::SetMySidEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::GetMySidEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMySidEntryAttributeRequest* req,
    lemming::dataplane::sai::GetMySidEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::CreateMySidEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateMySidEntriesRequest* req,
    lemming::dataplane::sai::CreateMySidEntriesResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Srv6::RemoveMySidEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveMySidEntriesRequest* req,
    lemming::dataplane::sai::RemoveMySidEntriesResponse* resp) {
  return grpc::Status::OK;
}
