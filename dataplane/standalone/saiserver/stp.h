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

#ifndef DATAPLANE_STANDALONE_SAI_STP_H_
#define DATAPLANE_STANDALONE_SAI_STP_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/stp.grpc.pb.h"
#include "dataplane/proto/sai/stp.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Stp final : public lemming::dataplane::sai::Stp::Service {
 public:
  grpc::Status CreateStp(grpc::ServerContext* context,
                         const lemming::dataplane::sai::CreateStpRequest* req,
                         lemming::dataplane::sai::CreateStpResponse* resp);

  grpc::Status RemoveStp(grpc::ServerContext* context,
                         const lemming::dataplane::sai::RemoveStpRequest* req,
                         lemming::dataplane::sai::RemoveStpResponse* resp);

  grpc::Status GetStpAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetStpAttributeRequest* req,
      lemming::dataplane::sai::GetStpAttributeResponse* resp);

  grpc::Status CreateStpPort(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateStpPortRequest* req,
      lemming::dataplane::sai::CreateStpPortResponse* resp);

  grpc::Status RemoveStpPort(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveStpPortRequest* req,
      lemming::dataplane::sai::RemoveStpPortResponse* resp);

  grpc::Status SetStpPortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetStpPortAttributeRequest* req,
      lemming::dataplane::sai::SetStpPortAttributeResponse* resp);

  grpc::Status GetStpPortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetStpPortAttributeRequest* req,
      lemming::dataplane::sai::GetStpPortAttributeResponse* resp);

  grpc::Status CreateStpPorts(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateStpPortsRequest* req,
      lemming::dataplane::sai::CreateStpPortsResponse* resp);

  grpc::Status RemoveStpPorts(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveStpPortsRequest* req,
      lemming::dataplane::sai::RemoveStpPortsResponse* resp);

  sai_stp_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_STP_H_
