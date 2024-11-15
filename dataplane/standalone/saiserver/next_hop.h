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

#ifndef DATAPLANE_STANDALONE_SAI_NEXT_HOP_H_
#define DATAPLANE_STANDALONE_SAI_NEXT_HOP_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/next_hop.grpc.pb.h"
#include "dataplane/proto/sai/next_hop.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class NextHop final : public lemming::dataplane::sai::NextHop::Service {
 public:
  grpc::Status CreateNextHop(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNextHopRequest* req,
      lemming::dataplane::sai::CreateNextHopResponse* resp);

  grpc::Status RemoveNextHop(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNextHopRequest* req,
      lemming::dataplane::sai::RemoveNextHopResponse* resp);

  grpc::Status SetNextHopAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetNextHopAttributeRequest* req,
      lemming::dataplane::sai::SetNextHopAttributeResponse* resp);

  grpc::Status GetNextHopAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetNextHopAttributeRequest* req,
      lemming::dataplane::sai::GetNextHopAttributeResponse* resp);

  grpc::Status CreateNextHops(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateNextHopsRequest* req,
      lemming::dataplane::sai::CreateNextHopsResponse* resp);

  grpc::Status RemoveNextHops(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveNextHopsRequest* req,
      lemming::dataplane::sai::RemoveNextHopsResponse* resp);

  sai_next_hop_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_NEXT_HOP_H_
