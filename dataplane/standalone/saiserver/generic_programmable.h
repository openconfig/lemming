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

#ifndef DATAPLANE_STANDALONE_SAI_GENERIC_PROGRAMMABLE_H_
#define DATAPLANE_STANDALONE_SAI_GENERIC_PROGRAMMABLE_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/generic_programmable.grpc.pb.h"
#include "dataplane/proto/sai/generic_programmable.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class GenericProgrammable final
    : public lemming::dataplane::sai::GenericProgrammable::Service {
 public:
  grpc::Status CreateGenericProgrammable(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateGenericProgrammableRequest* req,
      lemming::dataplane::sai::CreateGenericProgrammableResponse* resp);

  grpc::Status RemoveGenericProgrammable(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveGenericProgrammableRequest* req,
      lemming::dataplane::sai::RemoveGenericProgrammableResponse* resp);

  grpc::Status SetGenericProgrammableAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetGenericProgrammableAttributeRequest*
          req,
      lemming::dataplane::sai::SetGenericProgrammableAttributeResponse* resp);

  grpc::Status GetGenericProgrammableAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetGenericProgrammableAttributeRequest*
          req,
      lemming::dataplane::sai::GetGenericProgrammableAttributeResponse* resp);

  sai_generic_programmable_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_GENERIC_PROGRAMMABLE_H_
