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

#include "dataplane/standalone/saiserver/ipmc.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/ipmc.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Ipmc::CreateIpmcEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateIpmcEntryRequest* req,
    lemming::dataplane::sai::CreateIpmcEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipmc::RemoveIpmcEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveIpmcEntryRequest* req,
    lemming::dataplane::sai::RemoveIpmcEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipmc::SetIpmcEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetIpmcEntryAttributeRequest* req,
    lemming::dataplane::sai::SetIpmcEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Ipmc::GetIpmcEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetIpmcEntryAttributeRequest* req,
    lemming::dataplane::sai::GetIpmcEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}
