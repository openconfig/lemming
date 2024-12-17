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

#include "dataplane/standalone/saiserver/port.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/port.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Port::CreatePort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreatePortRequest* req,
    lemming::dataplane::sai::CreatePortResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::RemovePort(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemovePortRequest* req,
    lemming::dataplane::sai::RemovePortResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_port(req.get_oid());

  auto status = api->remove_port(entry);
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

grpc::Status Port::SetPortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetPortAttributeRequest* req,
    lemming::dataplane::sai::SetPortAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::GetPortAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetPortAttributeRequest* req,
    lemming::dataplane::sai::GetPortAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::GetPortStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetPortStatsRequest* req,
    lemming::dataplane::sai::GetPortStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::CreatePortPool(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreatePortPoolRequest* req,
    lemming::dataplane::sai::CreatePortPoolResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::RemovePortPool(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemovePortPoolRequest* req,
    lemming::dataplane::sai::RemovePortPoolResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_port_pool(req.get_oid());

  auto status = api->remove_port_pool(entry);
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

grpc::Status Port::SetPortPoolAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetPortPoolAttributeRequest* req,
    lemming::dataplane::sai::SetPortPoolAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::GetPortPoolAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetPortPoolAttributeRequest* req,
    lemming::dataplane::sai::GetPortPoolAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::GetPortPoolStats(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetPortPoolStatsRequest* req,
    lemming::dataplane::sai::GetPortPoolStatsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::CreatePortConnector(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreatePortConnectorRequest* req,
    lemming::dataplane::sai::CreatePortConnectorResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::RemovePortConnector(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemovePortConnectorRequest* req,
    lemming::dataplane::sai::RemovePortConnectorResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_port_connector(req.get_oid());

  auto status = api->remove_port_connector(entry);
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

grpc::Status Port::SetPortConnectorAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetPortConnectorAttributeRequest* req,
    lemming::dataplane::sai::SetPortConnectorAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::GetPortConnectorAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetPortConnectorAttributeRequest* req,
    lemming::dataplane::sai::GetPortConnectorAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::CreatePortSerdes(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreatePortSerdesRequest* req,
    lemming::dataplane::sai::CreatePortSerdesResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::RemovePortSerdes(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemovePortSerdesRequest* req,
    lemming::dataplane::sai::RemovePortSerdesResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_port_serdes(req.get_oid());

  auto status = api->remove_port_serdes(entry);
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

grpc::Status Port::GetPortSerdesAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetPortSerdesAttributeRequest* req,
    lemming::dataplane::sai::GetPortSerdesAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::CreatePorts(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreatePortsRequest* req,
    lemming::dataplane::sai::CreatePortsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Port::RemovePorts(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemovePortsRequest* req,
    lemming::dataplane::sai::RemovePortsResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
