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

#ifndef DATAPLANE_STANDALONE_SAI_HASH_H_
#define DATAPLANE_STANDALONE_SAI_HASH_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/hash.grpc.pb.h"
#include "dataplane/proto/sai/hash.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Hash final : public lemming::dataplane::sai::Hash::Service {
 public:
  grpc::Status CreateHash(grpc::ServerContext* context,
                          const lemming::dataplane::sai::CreateHashRequest* req,
                          lemming::dataplane::sai::CreateHashResponse* resp);

  grpc::Status RemoveHash(grpc::ServerContext* context,
                          const lemming::dataplane::sai::RemoveHashRequest* req,
                          lemming::dataplane::sai::RemoveHashResponse* resp);

  grpc::Status SetHashAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetHashAttributeRequest* req,
      lemming::dataplane::sai::SetHashAttributeResponse* resp);

  grpc::Status GetHashAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetHashAttributeRequest* req,
      lemming::dataplane::sai::GetHashAttributeResponse* resp);

  grpc::Status CreateFineGrainedHashField(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateFineGrainedHashFieldRequest* req,
      lemming::dataplane::sai::CreateFineGrainedHashFieldResponse* resp);

  grpc::Status RemoveFineGrainedHashField(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveFineGrainedHashFieldRequest* req,
      lemming::dataplane::sai::RemoveFineGrainedHashFieldResponse* resp);

  grpc::Status GetFineGrainedHashFieldAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetFineGrainedHashFieldAttributeRequest*
          req,
      lemming::dataplane::sai::GetFineGrainedHashFieldAttributeResponse* resp);

  sai_hash_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_HASH_H_
