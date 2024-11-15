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

#include "dataplane/standalone/saiserver/mpls.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/mpls.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Mpls::CreateInsegEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateInsegEntryRequest* req,
    lemming::dataplane::sai::CreateInsegEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Mpls::RemoveInsegEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveInsegEntryRequest* req,
    lemming::dataplane::sai::RemoveInsegEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Mpls::SetInsegEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetInsegEntryAttributeRequest* req,
    lemming::dataplane::sai::SetInsegEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Mpls::GetInsegEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetInsegEntryAttributeRequest* req,
    lemming::dataplane::sai::GetInsegEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Mpls::CreateInsegEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateInsegEntriesRequest* req,
    lemming::dataplane::sai::CreateInsegEntriesResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Mpls::RemoveInsegEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveInsegEntriesRequest* req,
    lemming::dataplane::sai::RemoveInsegEntriesResponse* resp) {
  return grpc::Status::OK;
}
