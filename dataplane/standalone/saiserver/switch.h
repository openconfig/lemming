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

#ifndef DATAPLANE_STANDALONE_SAI_SWITCH_H_
#define DATAPLANE_STANDALONE_SAI_SWITCH_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/switch.grpc.pb.h"
#include "dataplane/proto/sai/switch.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Switch final : public lemming::dataplane::sai::Switch::Service {
 public:
  grpc::Status CreateSwitch(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateSwitchRequest* req,
      lemming::dataplane::sai::CreateSwitchResponse* resp);

  grpc::Status RemoveSwitch(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveSwitchRequest* req,
      lemming::dataplane::sai::RemoveSwitchResponse* resp);

  grpc::Status SetSwitchAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetSwitchAttributeRequest* req,
      lemming::dataplane::sai::SetSwitchAttributeResponse* resp);

  grpc::Status GetSwitchAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetSwitchAttributeRequest* req,
      lemming::dataplane::sai::GetSwitchAttributeResponse* resp);

  grpc::Status GetSwitchStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetSwitchStatsRequest* req,
      lemming::dataplane::sai::GetSwitchStatsResponse* resp);

  grpc::Status CreateSwitchTunnel(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateSwitchTunnelRequest* req,
      lemming::dataplane::sai::CreateSwitchTunnelResponse* resp);

  grpc::Status RemoveSwitchTunnel(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveSwitchTunnelRequest* req,
      lemming::dataplane::sai::RemoveSwitchTunnelResponse* resp);

  grpc::Status SetSwitchTunnelAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetSwitchTunnelAttributeRequest* req,
      lemming::dataplane::sai::SetSwitchTunnelAttributeResponse* resp);

  grpc::Status GetSwitchTunnelAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetSwitchTunnelAttributeRequest* req,
      lemming::dataplane::sai::GetSwitchTunnelAttributeResponse* resp);

  sai_switch_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_SWITCH_H_
