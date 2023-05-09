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

#ifndef DATAPLANE_STANDALONE_TRANSLATOR_H_
#define DATAPLANE_STANDALONE_TRANSLATOR_H_

#include <grpcpp/channel.h>
#include <grpcpp/security/credentials.h>

#include <memory>
#include <unordered_map>
#include <utility>
#include <vector>

#include "dataplane/standalone/port.h"
#include "dataplane/standalone/switch.h"
#include "proto/forwarding/forwarding_service.grpc.pb.h"
#include "proto/forwarding/forwarding_service.pb.h"

extern "C" {
#include "inc/sai.h"
}

class Switch;
class Port;

class SaiObject {
 public:
  sai_object_type_t type;
  std::unordered_map<sai_attr_id_t, sai_attribute_value_t> attributes;
};

class Translator {
 public:
  explicit Translator(std::shared_ptr<grpc::Channel> chan) {
    client = std::shared_ptr<forwarding::Forwarding::Stub>(
        forwarding::Forwarding::NewStub(chan));
    objects[0] = {.type =
                      SAI_OBJECT_TYPE_NULL};  // ID == 0, is invalid so skip.
    sw = std::make_unique<Switch>(this, client);
    port = std::make_unique<Port>(this, client);
  }
  sai_object_type_t getObjectType(sai_object_id_t id);
  sai_object_id_t createObject(sai_object_type_t type);
  void setAttribute(sai_object_id_t id, sai_attribute_t attr);
  sai_status_t getAttribute(sai_object_id_t id, sai_attribute_t *attr);

  std::unique_ptr<Switch> sw;
  std::unique_ptr<Port> port;

 private:
  std::shared_ptr<forwarding::Forwarding::Stub> client;
  std::unordered_map<sai_object_id_t, SaiObject> objects;
};

#endif  // DATAPLANE_STANDALONE_TRANSLATOR_H_
