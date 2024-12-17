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

#include "dataplane/standalone/saiserver/l2mc.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/l2mc.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status L2mc::CreateL2mcEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateL2mcEntryRequest* req,
    lemming::dataplane::sai::CreateL2mcEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status L2mc::RemoveL2mcEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveL2mcEntryRequest* req,
    lemming::dataplane::sai::RemoveL2mcEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status L2mc::SetL2mcEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetL2mcEntryAttributeRequest* req,
    lemming::dataplane::sai::SetL2mcEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status L2mc::GetL2mcEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetL2mcEntryAttributeRequest* req,
    lemming::dataplane::sai::GetL2mcEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
