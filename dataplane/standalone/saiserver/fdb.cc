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

#include "dataplane/standalone/saiserver/fdb.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/fdb.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Fdb::CreateFdbEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateFdbEntryRequest* req,
    lemming::dataplane::sai::CreateFdbEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Fdb::RemoveFdbEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveFdbEntryRequest* req,
    lemming::dataplane::sai::RemoveFdbEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Fdb::SetFdbEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetFdbEntryAttributeRequest* req,
    lemming::dataplane::sai::SetFdbEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Fdb::GetFdbEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetFdbEntryAttributeRequest* req,
    lemming::dataplane::sai::GetFdbEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Fdb::CreateFdbEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateFdbEntriesRequest* req,
    lemming::dataplane::sai::CreateFdbEntriesResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Fdb::RemoveFdbEntries(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveFdbEntriesRequest* req,
    lemming::dataplane::sai::RemoveFdbEntriesResponse* resp) {
  return grpc::Status::OK;
}
