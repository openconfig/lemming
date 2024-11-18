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

#include "dataplane/standalone/saiserver/mcast_fdb.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/mcast_fdb.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status McastFdb::CreateMcastFdbEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateMcastFdbEntryRequest* req,
    lemming::dataplane::sai::CreateMcastFdbEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status McastFdb::RemoveMcastFdbEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveMcastFdbEntryRequest* req,
    lemming::dataplane::sai::RemoveMcastFdbEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status McastFdb::SetMcastFdbEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetMcastFdbEntryAttributeRequest* req,
    lemming::dataplane::sai::SetMcastFdbEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status McastFdb::GetMcastFdbEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMcastFdbEntryAttributeRequest* req,
    lemming::dataplane::sai::GetMcastFdbEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}
