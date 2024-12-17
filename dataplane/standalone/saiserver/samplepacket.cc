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

#include "dataplane/standalone/saiserver/samplepacket.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/samplepacket.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Samplepacket::CreateSamplepacket(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateSamplepacketRequest* req,
    lemming::dataplane::sai::CreateSamplepacketResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Samplepacket::RemoveSamplepacket(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveSamplepacketRequest* req,
    lemming::dataplane::sai::RemoveSamplepacketResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_samplepacket(req.get_oid());

  auto status = api->remove_samplepacket(entry);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return grpc::Status::INTERNAL;
  }

  return grpc::Status::OK;
}

grpc::Status Samplepacket::SetSamplepacketAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetSamplepacketAttributeRequest* req,
    lemming::dataplane::sai::SetSamplepacketAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Samplepacket::GetSamplepacketAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetSamplepacketAttributeRequest* req,
    lemming::dataplane::sai::GetSamplepacketAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
