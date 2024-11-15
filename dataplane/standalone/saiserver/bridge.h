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

#ifndef DATAPLANE_STANDALONE_SAI_BRIDGE_H_
#define DATAPLANE_STANDALONE_SAI_BRIDGE_H_

#include "dataplane/proto/sai/bridge.grpc.pb.h"
#include "dataplane/proto/sai/bridge.pb.h"
#include "dataplane/proto/sai/common.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Bridge final : public lemming::dataplane::sai::Bridge::Service {
 public:
  grpc::Status CreateBridge(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateBridgeRequest* req,
      lemming::dataplane::sai::CreateBridgeResponse* resp);

  grpc::Status RemoveBridge(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveBridgeRequest* req,
      lemming::dataplane::sai::RemoveBridgeResponse* resp);

  grpc::Status SetBridgeAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetBridgeAttributeRequest* req,
      lemming::dataplane::sai::SetBridgeAttributeResponse* resp);

  grpc::Status GetBridgeAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetBridgeAttributeRequest* req,
      lemming::dataplane::sai::GetBridgeAttributeResponse* resp);

  grpc::Status GetBridgeStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetBridgeStatsRequest* req,
      lemming::dataplane::sai::GetBridgeStatsResponse* resp);

  grpc::Status CreateBridgePort(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateBridgePortRequest* req,
      lemming::dataplane::sai::CreateBridgePortResponse* resp);

  grpc::Status RemoveBridgePort(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveBridgePortRequest* req,
      lemming::dataplane::sai::RemoveBridgePortResponse* resp);

  grpc::Status SetBridgePortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetBridgePortAttributeRequest* req,
      lemming::dataplane::sai::SetBridgePortAttributeResponse* resp);

  grpc::Status GetBridgePortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetBridgePortAttributeRequest* req,
      lemming::dataplane::sai::GetBridgePortAttributeResponse* resp);

  grpc::Status GetBridgePortStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetBridgePortStatsRequest* req,
      lemming::dataplane::sai::GetBridgePortStatsResponse* resp);

  sai_bridge_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_BRIDGE_H_
