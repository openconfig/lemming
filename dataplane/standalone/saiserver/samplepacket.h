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

#ifndef DATAPLANE_STANDALONE_SAI_SAMPLEPACKET_H_
#define DATAPLANE_STANDALONE_SAI_SAMPLEPACKET_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/samplepacket.grpc.pb.h"
#include "dataplane/proto/sai/samplepacket.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Samplepacket final
    : public lemming::dataplane::sai::Samplepacket::Service {
 public:
  grpc::Status CreateSamplepacket(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateSamplepacketRequest* req,
      lemming::dataplane::sai::CreateSamplepacketResponse* resp);

  grpc::Status RemoveSamplepacket(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveSamplepacketRequest* req,
      lemming::dataplane::sai::RemoveSamplepacketResponse* resp);

  grpc::Status SetSamplepacketAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetSamplepacketAttributeRequest* req,
      lemming::dataplane::sai::SetSamplepacketAttributeResponse* resp);

  grpc::Status GetSamplepacketAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetSamplepacketAttributeRequest* req,
      lemming::dataplane::sai::GetSamplepacketAttributeResponse* resp);

  sai_samplepacket_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_SAMPLEPACKET_H_
