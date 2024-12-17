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
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::RemoveTunnelMap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTunnelMapRequest* req,
    lemming::dataplane::sai::RemoveTunnelMapResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tunnel_map(req.get_oid());

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

grpc::Status Tunnel::GetTunnelMapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTunnelMapAttributeRequest* req,
    lemming::dataplane::sai::GetTunnelMapAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::CreateTunnel(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTunnelRequest* req,
    lemming::dataplane::sai::CreateTunnelResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::RemoveTunnel(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTunnelRequest* req,
    lemming::dataplane::sai::RemoveTunnelResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tunnel(req.get_oid());

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

grpc::Status Tunnel::SetTunnelAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTunnelAttributeRequest* req,
    lemming::dataplane::sai::SetTunnelAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::GetTunnelAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTunnelAttributeRequest* req,
    lemming::dataplane::sai::GetTunnelAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::GetTunnelStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTunnelStatsRequest* req,
    lemming::dataplane::sai::GetTunnelStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::CreateTunnelTermTableEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTunnelTermTableEntryRequest* req,
    lemming::dataplane::sai::CreateTunnelTermTableEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::RemoveTunnelTermTableEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTunnelTermTableEntryRequest* req,
    lemming::dataplane::sai::RemoveTunnelTermTableEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tunnel_term_table_entry(req.get_oid());

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

grpc::Status Tunnel::SetTunnelTermTableEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTunnelTermTableEntryAttributeRequest* req,
    lemming::dataplane::sai::SetTunnelTermTableEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::GetTunnelTermTableEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTunnelTermTableEntryAttributeRequest* req,
    lemming::dataplane::sai::GetTunnelTermTableEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::CreateTunnelMapEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTunnelMapEntryRequest* req,
    lemming::dataplane::sai::CreateTunnelMapEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::RemoveTunnelMapEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTunnelMapEntryRequest* req,
    lemming::dataplane::sai::RemoveTunnelMapEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tunnel_map_entry(req.get_oid());

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

grpc::Status Tunnel::GetTunnelMapEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTunnelMapEntryAttributeRequest* req,
    lemming::dataplane::sai::GetTunnelMapEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::CreateTunnels(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTunnelsRequest* req,
    lemming::dataplane::sai::CreateTunnelsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tunnel::RemoveTunnels(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTunnelsRequest* req,
    lemming::dataplane::sai::RemoveTunnelsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
