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

#include "dataplane/standalone/saiserver/switch.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/switch.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Switch::CreateSwitch(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateSwitchRequest* req,
    lemming::dataplane::sai::CreateSwitchResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Switch::RemoveSwitch(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveSwitchRequest* req,
    lemming::dataplane::sai::RemoveSwitchResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Switch::SetSwitchAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetSwitchAttributeRequest* req,
    lemming::dataplane::sai::SetSwitchAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Switch::GetSwitchAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetSwitchAttributeRequest* req,
    lemming::dataplane::sai::GetSwitchAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Switch::GetSwitchStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetSwitchStatsRequest* req,
    lemming::dataplane::sai::GetSwitchStatsResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Switch::CreateSwitchTunnel(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateSwitchTunnelRequest* req,
    lemming::dataplane::sai::CreateSwitchTunnelResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Switch::RemoveSwitchTunnel(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveSwitchTunnelRequest* req,
    lemming::dataplane::sai::RemoveSwitchTunnelResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Switch::SetSwitchTunnelAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetSwitchTunnelAttributeRequest* req,
    lemming::dataplane::sai::SetSwitchTunnelAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Switch::GetSwitchTunnelAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetSwitchTunnelAttributeRequest* req,
    lemming::dataplane::sai::GetSwitchTunnelAttributeResponse* resp) {
  return grpc::Status::OK;
}
