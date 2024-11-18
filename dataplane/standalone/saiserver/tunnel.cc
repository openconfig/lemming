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

#include "dataplane/standalone/saiserver/tunnel.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/tunnel.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Tunnel::CreateTunnelMap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTunnelMapRequest* req,
    lemming::dataplane::sai::CreateTunnelMapResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::RemoveTunnelMap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTunnelMapRequest* req,
    lemming::dataplane::sai::RemoveTunnelMapResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::GetTunnelMapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTunnelMapAttributeRequest* req,
    lemming::dataplane::sai::GetTunnelMapAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::CreateTunnel(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTunnelRequest* req,
    lemming::dataplane::sai::CreateTunnelResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::RemoveTunnel(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTunnelRequest* req,
    lemming::dataplane::sai::RemoveTunnelResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::SetTunnelAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTunnelAttributeRequest* req,
    lemming::dataplane::sai::SetTunnelAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::GetTunnelAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTunnelAttributeRequest* req,
    lemming::dataplane::sai::GetTunnelAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::GetTunnelStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTunnelStatsRequest* req,
    lemming::dataplane::sai::GetTunnelStatsResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::CreateTunnelTermTableEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTunnelTermTableEntryRequest* req,
    lemming::dataplane::sai::CreateTunnelTermTableEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::RemoveTunnelTermTableEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTunnelTermTableEntryRequest* req,
    lemming::dataplane::sai::RemoveTunnelTermTableEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::SetTunnelTermTableEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTunnelTermTableEntryAttributeRequest* req,
    lemming::dataplane::sai::SetTunnelTermTableEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::GetTunnelTermTableEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTunnelTermTableEntryAttributeRequest* req,
    lemming::dataplane::sai::GetTunnelTermTableEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::CreateTunnelMapEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTunnelMapEntryRequest* req,
    lemming::dataplane::sai::CreateTunnelMapEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::RemoveTunnelMapEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTunnelMapEntryRequest* req,
    lemming::dataplane::sai::RemoveTunnelMapEntryResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::GetTunnelMapEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTunnelMapEntryAttributeRequest* req,
    lemming::dataplane::sai::GetTunnelMapEntryAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::CreateTunnels(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTunnelsRequest* req,
    lemming::dataplane::sai::CreateTunnelsResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Tunnel::RemoveTunnels(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTunnelsRequest* req,
    lemming::dataplane::sai::RemoveTunnelsResponse* resp) {
  return grpc::Status::OK;
}
