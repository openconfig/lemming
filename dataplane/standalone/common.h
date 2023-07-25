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

#ifndef DATAPLANE_STANDALONE_COMMON_H_
#define DATAPLANE_STANDALONE_COMMON_H_

#include <memory>
#include <string>
#include <unordered_map>

#include "dataplane/standalone/sai/entry.h"
#include "proto/dataplane/dataplane.grpc.pb.h"
#include "proto/dataplane/dataplane.pb.h"
#include "proto/forwarding/forwarding_service.grpc.pb.h"
#include "proto/forwarding/forwarding_service.pb.h"

extern "C" {
#include "inc/sai.h"
}

const char contextID[] = "lucius";

// SaiObject is an object and its attributes.
class SaiObject {
 public:
  sai_object_type_t type;
  std::string switch_id;
  std::unordered_map<sai_attr_id_t, sai_attribute_value_t> attributes;
};

// AttributeManager tracks objects and their attributes values.
class AttributeManager {
 public:
  AttributeManager() {
    objects[std::string("0")] = {
        .type = SAI_OBJECT_TYPE_NULL};  // ID == 0, is invalid so skip.
  }
  sai_object_id_t create(sai_object_type_t type, std::string switch_id);
  std::string create(sai_object_type_t type, common_entry_t entry);
  sai_object_type_t get_type(std::string id);
  std::string get_switch_id(std::string id);

  void set_attribute(std::string id, const sai_attribute_t* attr);
  void set_attribute(std::string id, sai_attribute_t attr);

  sai_status_t get_attribute(std::string id, uint32_t attr_count,
                             sai_attribute_t* attr_list);

  std::string serialize_entry(sai_object_type_t type, common_entry_t id);
  std::string entry_to_switch_id(sai_object_type_t type, common_entry_t id);

 private:
  std::unordered_map<std::string, SaiObject> objects;
};

// APIBase is a base that all implementation of SAI APIs should inherit.
// TODO(dgrau): Verify no concurrent access or add mutex.
class APIBase {
 public:
  APIBase(std::string id, std::shared_ptr<AttributeManager> mgr,
          std::shared_ptr<forwarding::Forwarding::Stub> fwd,
          std::shared_ptr<lemming::dataplane::Dataplane::Stub> dplane)
      : id(id), attrMgr(mgr), fwd(fwd), dataplane(dplane) {}
  virtual ~APIBase() = default;
  virtual sai_status_t create(common_entry_t id, _In_ uint32_t attr_count,
                              _In_ const sai_attribute_t* attr_list);
  virtual sai_status_t create(_In_ uint32_t attr_count,
                              _In_ const sai_attribute_t* attr_list);
  virtual sai_status_t set_attribute(_In_ const sai_attribute_t* attr);

 protected:
  std::string id;
  std::shared_ptr<AttributeManager> attrMgr;
  std::shared_ptr<forwarding::Forwarding::Stub> fwd;
  std::shared_ptr<lemming::dataplane::Dataplane::Stub> dataplane;
};

#endif  // DATAPLANE_STANDALONE_COMMON_H_
