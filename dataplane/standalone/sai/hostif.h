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

#ifndef DATAPLANE_STANDALONE_SAI_HOSTIF_H_
#define DATAPLANE_STANDALONE_SAI_HOSTIF_H_

extern "C" {
#include "inc/sai.h"
#include "experimental/saiextensions.h"
}

extern const sai_hostif_api_t l_hostif;

sai_status_t l_create_hostif(sai_object_id_t *hostif_id,
                             sai_object_id_t switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list);

sai_status_t l_remove_hostif(sai_object_id_t hostif_id);

sai_status_t l_set_hostif_attribute(sai_object_id_t hostif_id,
                                    const sai_attribute_t *attr);

sai_status_t l_get_hostif_attribute(sai_object_id_t hostif_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list);

sai_status_t l_create_hostif_table_entry(sai_object_id_t *hostif_table_entry_id,
                                         sai_object_id_t switch_id,
                                         uint32_t attr_count,
                                         const sai_attribute_t *attr_list);

sai_status_t l_remove_hostif_table_entry(sai_object_id_t hostif_table_entry_id);

sai_status_t l_set_hostif_table_entry_attribute(
    sai_object_id_t hostif_table_entry_id, const sai_attribute_t *attr);

sai_status_t l_get_hostif_table_entry_attribute(
    sai_object_id_t hostif_table_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list);

sai_status_t l_create_hostif_trap_group(sai_object_id_t *hostif_trap_group_id,
                                        sai_object_id_t switch_id,
                                        uint32_t attr_count,
                                        const sai_attribute_t *attr_list);

sai_status_t l_remove_hostif_trap_group(sai_object_id_t hostif_trap_group_id);

sai_status_t l_set_hostif_trap_group_attribute(
    sai_object_id_t hostif_trap_group_id, const sai_attribute_t *attr);

sai_status_t l_get_hostif_trap_group_attribute(
    sai_object_id_t hostif_trap_group_id, uint32_t attr_count,
    sai_attribute_t *attr_list);

sai_status_t l_create_hostif_trap(sai_object_id_t *hostif_trap_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list);

sai_status_t l_remove_hostif_trap(sai_object_id_t hostif_trap_id);

sai_status_t l_set_hostif_trap_attribute(sai_object_id_t hostif_trap_id,
                                         const sai_attribute_t *attr);

sai_status_t l_get_hostif_trap_attribute(sai_object_id_t hostif_trap_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list);

sai_status_t l_create_hostif_user_defined_trap(
    sai_object_id_t *hostif_user_defined_trap_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_hostif_user_defined_trap(
    sai_object_id_t hostif_user_defined_trap_id);

sai_status_t l_set_hostif_user_defined_trap_attribute(
    sai_object_id_t hostif_user_defined_trap_id, const sai_attribute_t *attr);

sai_status_t l_get_hostif_user_defined_trap_attribute(
    sai_object_id_t hostif_user_defined_trap_id, uint32_t attr_count,
    sai_attribute_t *attr_list);

sai_status_t l_recv_hostif_packet(sai_object_id_t hostif_id,
                                  sai_size_t *buffer_size, void *buffer,
                                  uint32_t *attr_count,
                                  sai_attribute_t *attr_list);

sai_status_t l_send_hostif_packet(sai_object_id_t hostif_id,
                                  sai_size_t buffer_size, const void *buffer,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list);

sai_status_t l_allocate_hostif_packet(sai_object_id_t hostif_id,
                                      sai_size_t buffer_size, void **buffer,
                                      uint32_t attr_count,
                                      const sai_attribute_t *attr_list);

sai_status_t l_free_hostif_packet(sai_object_id_t hostif_id, void *buffer);

#endif  // DATAPLANE_STANDALONE_SAI_HOSTIF_H_
