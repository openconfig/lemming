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

#ifndef DATAPLANE_STANDALONE_SAI_TUNNEL_H_
#define DATAPLANE_STANDALONE_SAI_TUNNEL_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/tunnel.grpc.pb.h"
#include "dataplane/proto/sai/tunnel.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Tunnel final : public lemming::dataplane::sai::Tunnel::Service {
 public:
  grpc::Status CreateTunnelMap(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTunnelMapRequest* req,
      lemming::dataplane::sai::CreateTunnelMapResponse* resp);

  grpc::Status RemoveTunnelMap(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTunnelMapRequest* req,
      lemming::dataplane::sai::RemoveTunnelMapResponse* resp);

  grpc::Status GetTunnelMapAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTunnelMapAttributeRequest* req,
      lemming::dataplane::sai::GetTunnelMapAttributeResponse* resp);

  grpc::Status CreateTunnel(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTunnelRequest* req,
      lemming::dataplane::sai::CreateTunnelResponse* resp);

  grpc::Status RemoveTunnel(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTunnelRequest* req,
      lemming::dataplane::sai::RemoveTunnelResponse* resp);

  grpc::Status SetTunnelAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTunnelAttributeRequest* req,
      lemming::dataplane::sai::SetTunnelAttributeResponse* resp);

  grpc::Status GetTunnelAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTunnelAttributeRequest* req,
      lemming::dataplane::sai::GetTunnelAttributeResponse* resp);

  grpc::Status GetTunnelStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTunnelStatsRequest* req,
      lemming::dataplane::sai::GetTunnelStatsResponse* resp);

  grpc::Status CreateTunnelTermTableEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTunnelTermTableEntryRequest* req,
      lemming::dataplane::sai::CreateTunnelTermTableEntryResponse* resp);

  grpc::Status RemoveTunnelTermTableEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTunnelTermTableEntryRequest* req,
      lemming::dataplane::sai::RemoveTunnelTermTableEntryResponse* resp);

  grpc::Status SetTunnelTermTableEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTunnelTermTableEntryAttributeRequest*
          req,
      lemming::dataplane::sai::SetTunnelTermTableEntryAttributeResponse* resp);

  grpc::Status GetTunnelTermTableEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTunnelTermTableEntryAttributeRequest*
          req,
      lemming::dataplane::sai::GetTunnelTermTableEntryAttributeResponse* resp);

  grpc::Status CreateTunnelMapEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTunnelMapEntryRequest* req,
      lemming::dataplane::sai::CreateTunnelMapEntryResponse* resp);

  grpc::Status RemoveTunnelMapEntry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTunnelMapEntryRequest* req,
      lemming::dataplane::sai::RemoveTunnelMapEntryResponse* resp);

  grpc::Status GetTunnelMapEntryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTunnelMapEntryAttributeRequest* req,
      lemming::dataplane::sai::GetTunnelMapEntryAttributeResponse* resp);

  grpc::Status CreateTunnels(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTunnelsRequest* req,
      lemming::dataplane::sai::CreateTunnelsResponse* resp);

  grpc::Status RemoveTunnels(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTunnelsRequest* req,
      lemming::dataplane::sai::RemoveTunnelsResponse* resp);

  sai_tunnel_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_TUNNEL_H_
