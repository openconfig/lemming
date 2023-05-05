
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

#include "dataplane/standalone/sai/hash.h"

#include "dataplane/standalone/log/log.h"

const sai_hash_api_t l_hash = {
    .create_hash = l_create_hash,
    .remove_hash = l_remove_hash,
    .set_hash_attribute = l_set_hash_attribute,
    .get_hash_attribute = l_get_hash_attribute,
    .create_fine_grained_hash_field = l_create_fine_grained_hash_field,
    .remove_fine_grained_hash_field = l_remove_fine_grained_hash_field,
    .set_fine_grained_hash_field_attribute =
        l_set_fine_grained_hash_field_attribute,
    .get_fine_grained_hash_field_attribute =
        l_get_fine_grained_hash_field_attribute,
};

sai_status_t l_create_hash(sai_object_id_t *hash_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_hash(sai_object_id_t hash_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_hash_attribute(sai_object_id_t hash_id,
                                  const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_hash_attribute(sai_object_id_t hash_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_fine_grained_hash_field(
    sai_object_id_t *fine_grained_hash_field_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_fine_grained_hash_field(
    sai_object_id_t fine_grained_hash_field_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_fine_grained_hash_field_attribute(
    sai_object_id_t fine_grained_hash_field_id, const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_fine_grained_hash_field_attribute(
    sai_object_id_t fine_grained_hash_field_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}
