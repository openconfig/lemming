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

#ifndef DATAPLANE_STANDALONE_SAI_SYSTEM_PORT_H_
#define DATAPLANE_STANDALONE_SAI_SYSTEM_PORT_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/system_port.grpc.pb.h"
#include "dataplane/proto/sai/system_port.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class SystemPort final : public lemming::dataplane::sai::SystemPort::Service {
 public:
  grpc::Status CreateSystemPort(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateSystemPortRequest* req,
      lemming::dataplane::sai::CreateSystemPortResponse* resp);

  grpc::Status RemoveSystemPort(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveSystemPortRequest* req,
      lemming::dataplane::sai::RemoveSystemPortResponse* resp);

  grpc::Status SetSystemPortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetSystemPortAttributeRequest* req,
      lemming::dataplane::sai::SetSystemPortAttributeResponse* resp);

  grpc::Status GetSystemPortAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetSystemPortAttributeRequest* req,
      lemming::dataplane::sai::GetSystemPortAttributeResponse* resp);

  sai_system_port_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_SYSTEM_PORT_H_
