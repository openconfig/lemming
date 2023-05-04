
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

#include "dataplane/standalone/sai/next_hop_group.h"

#include "dataplane/standalone/log/log.h"

const sai_next_hop_group_api_t l_next_hop_group = {
    .create_next_hop_group = l_create_next_hop_group,
    .remove_next_hop_group = l_remove_next_hop_group,
    .set_next_hop_group_attribute = l_set_next_hop_group_attribute,
    .get_next_hop_group_attribute = l_get_next_hop_group_attribute,
    .create_next_hop_group_member = l_create_next_hop_group_member,
    .remove_next_hop_group_member = l_remove_next_hop_group_member,
    .set_next_hop_group_member_attribute =
        l_set_next_hop_group_member_attribute,
    .get_next_hop_group_member_attribute =
        l_get_next_hop_group_member_attribute,
    .create_next_hop_group_members = l_create_next_hop_group_members,
    .remove_next_hop_group_members = l_remove_next_hop_group_members,
};

sai_status_t l_create_next_hop_group(sai_object_id_t *next_hop_group_id,
                                     sai_object_id_t switch_id,
                                     uint32_t attr_count,
                                     const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_next_hop_group(sai_object_id_t next_hop_group_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_next_hop_group_attribute(sai_object_id_t next_hop_group_id,
                                            const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_next_hop_group_attribute(sai_object_id_t next_hop_group_id,
                                            uint32_t attr_count,
                                            sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_next_hop_group_member(
    sai_object_id_t *next_hop_group_member_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_next_hop_group_member(
    sai_object_id_t next_hop_group_member_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_next_hop_group_member_attribute(
    sai_object_id_t next_hop_group_member_id, const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_next_hop_group_member_attribute(
    sai_object_id_t next_hop_group_member_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_next_hop_group_members(sai_object_id_t switch_id,
                                             uint32_t object_count,
                                             const uint32_t *attr_count,
                                             const sai_attribute_t **attr_list,
                                             sai_bulk_op_error_mode_t mode,
                                             sai_object_id_t *object_id,
                                             sai_status_t *object_statuses) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_next_hop_group_members(uint32_t object_count,
                                             const sai_object_id_t *object_id,
                                             sai_bulk_op_error_mode_t mode,
                                             sai_status_t *object_statuses) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}
