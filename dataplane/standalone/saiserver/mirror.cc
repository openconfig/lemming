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

#include "dataplane/standalone/saiserver/mirror.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/mirror.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Mirror::CreateMirrorSession(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateMirrorSessionRequest* req,
    lemming::dataplane::sai::CreateMirrorSessionResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Mirror::RemoveMirrorSession(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveMirrorSessionRequest* req,
    lemming::dataplane::sai::RemoveMirrorSessionResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Mirror::SetMirrorSessionAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetMirrorSessionAttributeRequest* req,
    lemming::dataplane::sai::SetMirrorSessionAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Mirror::GetMirrorSessionAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetMirrorSessionAttributeRequest* req,
    lemming::dataplane::sai::GetMirrorSessionAttributeResponse* resp) {
  return grpc::Status::OK;
}
