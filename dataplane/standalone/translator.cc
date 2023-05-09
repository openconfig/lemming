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

#include "dataplane/standalone/translator.h"

#include "dataplane/standalone/log/log.h"
#include "dataplane/standalone/switch.h"

extern "C" {
#include "inc/sai.h"
}

sai_object_type_t Translator::getObjectType(sai_object_id_t id) {
  auto iter = this->objects.find(id);
  if (iter == this->objects.end()) {
    return SAI_OBJECT_TYPE_NULL;
  }
  return iter->second.type;
}

sai_object_id_t Translator::createObject(sai_object_type_t type) {
  auto id = this->objects.size() + 1;
  this->objects[id] = {
      .type = type,
      .attributes = std::unordered_map<sai_attr_id_t, sai_attribute_value_t>(),
  };
  return id;
}

void Translator::setAttribute(sai_object_id_t id, sai_attribute_t attr) {
  this->objects[id].attributes[attr.id] = attr.value;
}

sai_status_t Translator::getAttribute(sai_object_id_t id,
                                      sai_attribute_t* attr) {
  auto iter = this->objects[id].attributes.find(attr->id);
  if (iter == this->objects[id].attributes.end()) {
    return SAI_STATUS_ITEM_NOT_FOUND;
  }
  attr->value = iter->second;
  return SAI_STATUS_SUCCESS;
}
