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

#ifndef DATAPLANE_STANDALONE_SAI_MIRROR_H_
#define DATAPLANE_STANDALONE_SAI_MIRROR_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/mirror.grpc.pb.h"
#include "dataplane/proto/sai/mirror.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Mirror final : public lemming::dataplane::sai::Mirror::Service {
 public:
  grpc::Status CreateMirrorSession(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateMirrorSessionRequest* req,
      lemming::dataplane::sai::CreateMirrorSessionResponse* resp);

  grpc::Status RemoveMirrorSession(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveMirrorSessionRequest* req,
      lemming::dataplane::sai::RemoveMirrorSessionResponse* resp);

  grpc::Status SetMirrorSessionAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetMirrorSessionAttributeRequest* req,
      lemming::dataplane::sai::SetMirrorSessionAttributeResponse* resp);

  grpc::Status GetMirrorSessionAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMirrorSessionAttributeRequest* req,
      lemming::dataplane::sai::GetMirrorSessionAttributeResponse* resp);

  sai_mirror_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_MIRROR_H_
