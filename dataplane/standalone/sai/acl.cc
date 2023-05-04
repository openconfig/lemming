
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

#include "dataplane/standalone/sai/acl.h"

#include "dataplane/standalone/log/log.h"

const sai_acl_api_t l_acl = {
    .create_acl_table = l_create_acl_table,
    .remove_acl_table = l_remove_acl_table,
    .set_acl_table_attribute = l_set_acl_table_attribute,
    .get_acl_table_attribute = l_get_acl_table_attribute,
    .create_acl_entry = l_create_acl_entry,
    .remove_acl_entry = l_remove_acl_entry,
    .set_acl_entry_attribute = l_set_acl_entry_attribute,
    .get_acl_entry_attribute = l_get_acl_entry_attribute,
    .create_acl_counter = l_create_acl_counter,
    .remove_acl_counter = l_remove_acl_counter,
    .set_acl_counter_attribute = l_set_acl_counter_attribute,
    .get_acl_counter_attribute = l_get_acl_counter_attribute,
    .create_acl_range = l_create_acl_range,
    .remove_acl_range = l_remove_acl_range,
    .set_acl_range_attribute = l_set_acl_range_attribute,
    .get_acl_range_attribute = l_get_acl_range_attribute,
    .create_acl_table_group = l_create_acl_table_group,
    .remove_acl_table_group = l_remove_acl_table_group,
    .set_acl_table_group_attribute = l_set_acl_table_group_attribute,
    .get_acl_table_group_attribute = l_get_acl_table_group_attribute,
    .create_acl_table_group_member = l_create_acl_table_group_member,
    .remove_acl_table_group_member = l_remove_acl_table_group_member,
    .set_acl_table_group_member_attribute =
        l_set_acl_table_group_member_attribute,
    .get_acl_table_group_member_attribute =
        l_get_acl_table_group_member_attribute,
};

sai_status_t l_create_acl_table(sai_object_id_t *acl_table_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_acl_table(sai_object_id_t acl_table_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_acl_table_attribute(sai_object_id_t acl_table_id,
                                       const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_acl_table_attribute(sai_object_id_t acl_table_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_acl_entry(sai_object_id_t *acl_entry_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_acl_entry(sai_object_id_t acl_entry_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_acl_entry_attribute(sai_object_id_t acl_entry_id,
                                       const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_acl_entry_attribute(sai_object_id_t acl_entry_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_acl_counter(sai_object_id_t *acl_counter_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_acl_counter(sai_object_id_t acl_counter_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_acl_counter_attribute(sai_object_id_t acl_counter_id,
                                         const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_acl_counter_attribute(sai_object_id_t acl_counter_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_acl_range(sai_object_id_t *acl_range_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_acl_range(sai_object_id_t acl_range_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_acl_range_attribute(sai_object_id_t acl_range_id,
                                       const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_acl_range_attribute(sai_object_id_t acl_range_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_acl_table_group(sai_object_id_t *acl_table_group_id,
                                      sai_object_id_t switch_id,
                                      uint32_t attr_count,
                                      const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_acl_table_group(sai_object_id_t acl_table_group_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_acl_table_group_attribute(sai_object_id_t acl_table_group_id,
                                             const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_acl_table_group_attribute(sai_object_id_t acl_table_group_id,
                                             uint32_t attr_count,
                                             sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_acl_table_group_member(
    sai_object_id_t *acl_table_group_member_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_acl_table_group_member(
    sai_object_id_t acl_table_group_member_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_acl_table_group_member_attribute(
    sai_object_id_t acl_table_group_member_id, const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_acl_table_group_member_attribute(
    sai_object_id_t acl_table_group_member_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}
