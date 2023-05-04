
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

#ifndef DATAPLANE_STANDALONE_SAI_NAT_H_
#define DATAPLANE_STANDALONE_SAI_NAT_H_

extern "C" {
	#include "inc/sai.h"
}

extern const sai_nat_api_t l_nat;


sai_status_t l_create_nat_entry(const sai_nat_entry_t *nat_entry, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_nat_entry(const sai_nat_entry_t *nat_entry);

sai_status_t l_set_nat_entry_attribute(const sai_nat_entry_t *nat_entry, const sai_attribute_t *attr);

sai_status_t l_get_nat_entry_attribute(const sai_nat_entry_t *nat_entry, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_nat_entries(uint32_t object_count, const sai_nat_entry_t *nat_entry, const uint32_t *attr_count, const sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses);

sai_status_t l_remove_nat_entries(uint32_t object_count, const sai_nat_entry_t *nat_entry, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses);

sai_status_t l_set_nat_entries_attribute(uint32_t object_count, const sai_nat_entry_t *nat_entry, const sai_attribute_t *attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses);

sai_status_t l_get_nat_entries_attribute(uint32_t object_count, const sai_nat_entry_t *nat_entry, const uint32_t *attr_count, sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses);

sai_status_t l_create_nat_zone_counter(sai_object_id_t *nat_zone_counter_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_nat_zone_counter(sai_object_id_t nat_zone_counter_id);

sai_status_t l_set_nat_zone_counter_attribute(sai_object_id_t nat_zone_counter_id, const sai_attribute_t *attr);

sai_status_t l_get_nat_zone_counter_attribute(sai_object_id_t nat_zone_counter_id, uint32_t attr_count, sai_attribute_t *attr_list);


#endif  // DATAPLANE_STANDALONE_SAI_NAT_H_
