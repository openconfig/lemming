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

#ifndef DATAPLANE_STANDALONE_SAI_ACL_H_
#define DATAPLANE_STANDALONE_SAI_ACL_H_

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

extern const sai_acl_api_t l_acl;


sai_status_t l_create_acl_table(sai_object_id_t *acl_table_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_acl_table(sai_object_id_t acl_table_id);

sai_status_t l_set_acl_table_attribute(sai_object_id_t acl_table_id, const sai_attribute_t *attr);

sai_status_t l_get_acl_table_attribute(sai_object_id_t acl_table_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_acl_entry(sai_object_id_t *acl_entry_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_acl_entry(sai_object_id_t acl_entry_id);

sai_status_t l_set_acl_entry_attribute(sai_object_id_t acl_entry_id, const sai_attribute_t *attr);

sai_status_t l_get_acl_entry_attribute(sai_object_id_t acl_entry_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_acl_counter(sai_object_id_t *acl_counter_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_acl_counter(sai_object_id_t acl_counter_id);

sai_status_t l_set_acl_counter_attribute(sai_object_id_t acl_counter_id, const sai_attribute_t *attr);

sai_status_t l_get_acl_counter_attribute(sai_object_id_t acl_counter_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_acl_range(sai_object_id_t *acl_range_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_acl_range(sai_object_id_t acl_range_id);

sai_status_t l_set_acl_range_attribute(sai_object_id_t acl_range_id, const sai_attribute_t *attr);

sai_status_t l_get_acl_range_attribute(sai_object_id_t acl_range_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_acl_table_group(sai_object_id_t *acl_table_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_acl_table_group(sai_object_id_t acl_table_group_id);

sai_status_t l_set_acl_table_group_attribute(sai_object_id_t acl_table_group_id, const sai_attribute_t *attr);

sai_status_t l_get_acl_table_group_attribute(sai_object_id_t acl_table_group_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_acl_table_group_member(sai_object_id_t *acl_table_group_member_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_acl_table_group_member(sai_object_id_t acl_table_group_member_id);

sai_status_t l_set_acl_table_group_member_attribute(sai_object_id_t acl_table_group_member_id, const sai_attribute_t *attr);

sai_status_t l_get_acl_table_group_member_attribute(sai_object_id_t acl_table_group_member_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_acl_table_chain_group(sai_object_id_t *acl_table_chain_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_acl_table_chain_group(sai_object_id_t acl_table_chain_group_id);

sai_status_t l_set_acl_table_chain_group_attribute(sai_object_id_t acl_table_chain_group_id, const sai_attribute_t *attr);

sai_status_t l_get_acl_table_chain_group_attribute(sai_object_id_t acl_table_chain_group_id, uint32_t attr_count, sai_attribute_t *attr_list);


#endif  // DATAPLANE_STANDALONE_SAI_ACL_H_
