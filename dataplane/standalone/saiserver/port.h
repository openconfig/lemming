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

#ifndef DATAPLANE_STANDALONE_SAI_PORT_H_
#define DATAPLANE_STANDALONE_SAI_PORT_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/port.grpc.pb.h"
#include "dataplane/proto/sai/port.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Port final : public lemming::dataplane::sai::Port::Service {
 public:
  grpc::Status CreatePort(grpc::ServerContext* context,
                          const lemming::dataplane::sai::CreatePortRequest* req,
                          lemming::dataplane::sai::CreatePortResponse* resp);

  grpc::Status RemovePort(grpc::ServerContext* context,
                          const lemming::dataplane::sai::RemovePortRequest* req,
                          lemming::dataplane::sai::RemovePortResponse* resp);

  grpc::Status SetPortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetPortAttributeRequest* req,
      lemming::dataplane::sai::SetPortAttributeResponse* resp);

  grpc::Status GetPortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetPortAttributeRequest* req,
      lemming::dataplane::sai::GetPortAttributeResponse* resp);

  grpc::Status GetPortStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetPortStatsRequest* req,
      lemming::dataplane::sai::GetPortStatsResponse* resp);

  grpc::Status CreatePortPool(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreatePortPoolRequest* req,
      lemming::dataplane::sai::CreatePortPoolResponse* resp);

  grpc::Status RemovePortPool(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemovePortPoolRequest* req,
      lemming::dataplane::sai::RemovePortPoolResponse* resp);

  grpc::Status SetPortPoolAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetPortPoolAttributeRequest* req,
      lemming::dataplane::sai::SetPortPoolAttributeResponse* resp);

  grpc::Status GetPortPoolAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetPortPoolAttributeRequest* req,
      lemming::dataplane::sai::GetPortPoolAttributeResponse* resp);

  grpc::Status GetPortPoolStats(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetPortPoolStatsRequest* req,
      lemming::dataplane::sai::GetPortPoolStatsResponse* resp);

  grpc::Status CreatePortConnector(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreatePortConnectorRequest* req,
      lemming::dataplane::sai::CreatePortConnectorResponse* resp);

  grpc::Status RemovePortConnector(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemovePortConnectorRequest* req,
      lemming::dataplane::sai::RemovePortConnectorResponse* resp);

  grpc::Status SetPortConnectorAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetPortConnectorAttributeRequest* req,
      lemming::dataplane::sai::SetPortConnectorAttributeResponse* resp);

  grpc::Status GetPortConnectorAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetPortConnectorAttributeRequest* req,
      lemming::dataplane::sai::GetPortConnectorAttributeResponse* resp);

  grpc::Status CreatePortSerdes(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreatePortSerdesRequest* req,
      lemming::dataplane::sai::CreatePortSerdesResponse* resp);

  grpc::Status RemovePortSerdes(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemovePortSerdesRequest* req,
      lemming::dataplane::sai::RemovePortSerdesResponse* resp);

  grpc::Status GetPortSerdesAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetPortSerdesAttributeRequest* req,
      lemming::dataplane::sai::GetPortSerdesAttributeResponse* resp);

  grpc::Status CreatePorts(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreatePortsRequest* req,
      lemming::dataplane::sai::CreatePortsResponse* resp);

  grpc::Status RemovePorts(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemovePortsRequest* req,
      lemming::dataplane::sai::RemovePortsResponse* resp);

  sai_port_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_PORT_H_
