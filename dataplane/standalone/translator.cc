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
  *attr = {
      iter->first,
      iter->second,
  };
  return SAI_STATUS_SUCCESS;
}

sai_status_t Translator::create(sai_object_type_t type, sai_object_id_t* id,
                                uint32_t attr_count,
                                const sai_attribute_t* attr_list) {
  return sai_status_t();
}

sai_status_t Translator::create(sai_object_type_t type, common_entry_t id,
                                uint32_t attr_count,
                                const sai_attribute_t* attr_list) {
  return sai_status_t();
}

sai_status_t Translator::create(sai_object_type_t type, sai_object_id_t* id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t* attr_list) {
  return sai_status_t();
}

sai_status_t Translator::remove(sai_object_type_t type, sai_object_id_t id) {
  return sai_status_t();
}

sai_status_t Translator::remove(sai_object_type_t type, common_entry_t id) {
  return sai_status_t();
}

sai_status_t Translator::set_attribute(sai_object_type_t type,
                                       sai_object_id_t id,
                                       const sai_attribute_t* attr) {
  return sai_status_t();
}

sai_status_t Translator::get_attribute(sai_object_type_t type,
                                       sai_object_id_t id, uint32_t attr_count,
                                       sai_attribute_t* attr_list) {
  return sai_status_t();
}

sai_status_t Translator::set_attribute(sai_object_type_t type,
                                       common_entry_t id,
                                       const sai_attribute_t* attr) {
  return sai_status_t();
}

sai_status_t Translator::get_attribute(sai_object_type_t type,
                                       common_entry_t id, uint32_t attr_count,
                                       sai_attribute_t* attr_list) {
  return sai_status_t();
}

sai_status_t Translator::get_stats(sai_object_type_t type, sai_object_id_t id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t* counter_ids,
                                   uint64_t* counters) {
  return sai_status_t();
}

sai_status_t Translator::get_stats_ext(sai_object_type_t type,
                                       sai_object_id_t bfd_session_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t* counter_ids,
                                       sai_stats_mode_t mode,
                                       uint64_t* counters) {
  return sai_status_t();
}

sai_status_t Translator::clear_stats(sai_object_type_t type,
                                     sai_object_id_t bfd_session_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t* counter_ids) {
  return sai_status_t();
}

sai_status_t Translator::create_bulk(
    sai_object_type_t type, sai_object_id_t switch_id, uint32_t object_count,
    const uint32_t* attr_count, const sai_attribute_t** attr_list,
    sai_bulk_op_error_mode_t mode, sai_object_id_t* object_id,
    sai_status_t* object_statuses) {
  return sai_status_t();
}

sai_status_t Translator::remove_bulk(sai_object_type_t type,
                                     uint32_t object_count,
                                     const sai_object_id_t* object_id,
                                     sai_bulk_op_error_mode_t mode,
                                     sai_status_t* object_statuses) {
  return sai_status_t();
}

sai_status_t Translator::create_bulk(
    sai_object_type_t type, uint32_t object_count, common_entry_t object_id,
    const uint32_t* attr_count, const sai_attribute_t** attr_list,
    sai_bulk_op_error_mode_t mode, sai_status_t* object_statuses) {
  return sai_status_t();
}

sai_status_t Translator::remove_bulk(sai_object_type_t type,
                                     uint32_t object_count,
                                     common_entry_t object_id,
                                     sai_bulk_op_error_mode_t mode,
                                     sai_status_t* object_statuses) {
  return sai_status_t();
}

sai_status_t Translator::set_attribute_bulk(sai_object_type_t type,
                                            uint32_t object_count,
                                            common_entry_t object_id,
                                            const sai_attribute_t* attr_list,
                                            sai_bulk_op_error_mode_t mode,
                                            sai_status_t* object_statuses) {
  return sai_status_t();
}

sai_status_t Translator::get_attribute_bulk(
    sai_object_type_t type, uint32_t object_count, common_entry_t object_id,
    const uint32_t* attr_count, sai_attribute_t** attr_list,
    sai_bulk_op_error_mode_t mode, sai_status_t* object_statuses) {
  return sai_status_t();
}
