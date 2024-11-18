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

#ifndef DATAPLANE_STANDALONE_SAI_QOS_MAP_H_
#define DATAPLANE_STANDALONE_SAI_QOS_MAP_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/qos_map.grpc.pb.h"
#include "dataplane/proto/sai/qos_map.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class QosMap final : public lemming::dataplane::sai::QosMap::Service {
 public:
  grpc::Status CreateQosMap(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateQosMapRequest* req,
      lemming::dataplane::sai::CreateQosMapResponse* resp);

  grpc::Status RemoveQosMap(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveQosMapRequest* req,
      lemming::dataplane::sai::RemoveQosMapResponse* resp);

  grpc::Status SetQosMapAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetQosMapAttributeRequest* req,
      lemming::dataplane::sai::SetQosMapAttributeResponse* resp);

  grpc::Status GetQosMapAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetQosMapAttributeRequest* req,
      lemming::dataplane::sai::GetQosMapAttributeResponse* resp);

  sai_qos_map_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_QOS_MAP_H_
