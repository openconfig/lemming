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

#ifndef DATAPLANE_STANDALONE_PORT_H_
#define DATAPLANE_STANDALONE_PORT_H_

#include <memory>
#include <unordered_map>

#include "dataplane/standalone/translator.h"
#include "proto/forwarding/forwarding_service.grpc.pb.h"
#include "proto/forwarding/forwarding_service.pb.h"

extern "C" {
#include "inc/sai.h"
}

class Translator;

class Port {
 public:
  Port(Translator* translator, std::shared_ptr<forwarding::Forwarding::Stub> c)
      : translator(translator), client(c) {}
  sai_status_t create_port(_Out_ sai_object_id_t* port_id,
                           _In_ sai_object_id_t switch_id,
                           _In_ uint32_t attr_count,
                           _In_ const sai_attribute_t* attr_list);
  sai_status_t set_port_attribute(_In_ sai_object_id_t port_id,
                                  _In_ const sai_attribute_t* attr);
  sai_status_t get_port_attribute(_In_ sai_object_id_t port_id,
                                  _In_ uint32_t attr_count,
                                  _Inout_ sai_attribute_t* attr_list);

 private:
  std::shared_ptr<Translator> translator;
  std::shared_ptr<forwarding::Forwarding::Stub> client;
  std::unordered_map<sai_attr_id_t, sai_attribute_value_t> attributes;
};

#endif  // DATAPLANE_STANDALONE_PORT_H_
