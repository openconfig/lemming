// Copyright 2023 Google LLC
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

#ifndef TRANSLATOR_H
#define TRANSLATOR_H

#include "proto/forwarding/forwarding_service.pb.h"
#include "proto/forwarding/forwarding_service.grpc.pb.h"
#include <grpcpp/channel.h>
#include <grpcpp/security/credentials.h>
#include "dataplane/standalone/lucius/lucius_clib.h"

extern "C"
{
#include "sai.h"
}

class Translator
{
public:
    Translator(std::shared_ptr<grpc::Channel> chan) { client = forwarding::Forwarding::NewStub(chan); };
    sai_status_t create_switch(_Out_ sai_object_id_t *switch_id, _In_ uint32_t attr_count, _In_ const sai_attribute_t *attr_list);

private:
    std::unique_ptr<forwarding::Forwarding::Stub> client;
};

#endif