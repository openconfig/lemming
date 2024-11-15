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

#ifndef DATAPLANE_STANDALONE_SAI_MY_MAC_H_
#define DATAPLANE_STANDALONE_SAI_MY_MAC_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/my_mac.grpc.pb.h"
#include "dataplane/proto/sai/my_mac.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class MyMac final : public lemming::dataplane::sai::MyMac::Service {
 public:
  grpc::Status CreateMyMac(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateMyMacRequest* req,
      lemming::dataplane::sai::CreateMyMacResponse* resp);

  grpc::Status RemoveMyMac(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveMyMacRequest* req,
      lemming::dataplane::sai::RemoveMyMacResponse* resp);

  grpc::Status SetMyMacAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetMyMacAttributeRequest* req,
      lemming::dataplane::sai::SetMyMacAttributeResponse* resp);

  grpc::Status GetMyMacAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetMyMacAttributeRequest* req,
      lemming::dataplane::sai::GetMyMacAttributeResponse* resp);

  sai_my_mac_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_MY_MAC_H_
