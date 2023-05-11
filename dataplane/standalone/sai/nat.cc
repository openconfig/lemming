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

#include "dataplane/standalone/sai/nat.h"
#include <glog/logging.h>
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_nat_api_t l_nat = {
	.create_nat_entry = l_create_nat_entry,
	.remove_nat_entry = l_remove_nat_entry,
	.set_nat_entry_attribute = l_set_nat_entry_attribute,
	.get_nat_entry_attribute = l_get_nat_entry_attribute,
	.create_nat_entries = l_create_nat_entries,
	.remove_nat_entries = l_remove_nat_entries,
	.set_nat_entries_attribute = l_set_nat_entries_attribute,
	.get_nat_entries_attribute = l_get_nat_entries_attribute,
	.create_nat_zone_counter = l_create_nat_zone_counter,
	.remove_nat_zone_counter = l_remove_nat_zone_counter,
	.set_nat_zone_counter_attribute = l_set_nat_zone_counter_attribute,
	.get_nat_zone_counter_attribute = l_get_nat_zone_counter_attribute,
};


sai_status_t l_create_nat_entry(const sai_nat_entry_t *nat_entry, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__; common_entry_t entry = {.nat_entry = nat_entry}; 
	return translator->create(SAI_OBJECT_TYPE_NAT_ENTRY, entry, attr_count, attr_list);
}

sai_status_t l_remove_nat_entry(const sai_nat_entry_t *nat_entry) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__; common_entry_t entry = {.nat_entry = nat_entry}; 
	return translator->remove(SAI_OBJECT_TYPE_NAT_ENTRY, entry);
}

sai_status_t l_set_nat_entry_attribute(const sai_nat_entry_t *nat_entry, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__; common_entry_t entry = {.nat_entry = nat_entry}; 
	return translator->set_attribute(SAI_OBJECT_TYPE_NAT_ENTRY, entry, attr);
}

sai_status_t l_get_nat_entry_attribute(const sai_nat_entry_t *nat_entry, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__; common_entry_t entry = {.nat_entry = nat_entry}; 
	return translator->get_attribute(SAI_OBJECT_TYPE_NAT_ENTRY, entry, attr_count, attr_list);
}

sai_status_t l_create_nat_entries(uint32_t object_count, const sai_nat_entry_t *nat_entry, const uint32_t *attr_count, const sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__; common_entry_t entry = {.nat_entry = nat_entry}; 
	return translator->create_bulk(SAI_OBJECT_TYPE_NAT_ENTRY, object_count, entry, attr_count, attr_list, mode, object_statuses);
}

sai_status_t l_remove_nat_entries(uint32_t object_count, const sai_nat_entry_t *nat_entry, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__; common_entry_t entry = {.nat_entry = nat_entry}; 
	return translator->remove_bulk(SAI_OBJECT_TYPE_NAT_ENTRY, object_count, entry, mode, object_statuses);
}

sai_status_t l_set_nat_entries_attribute(uint32_t object_count, const sai_nat_entry_t *nat_entry, const sai_attribute_t *attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__; common_entry_t entry = {.nat_entry = nat_entry}; 
	return translator->set_attribute_bulk(SAI_OBJECT_TYPE_NAT_ENTRY, object_count, entry, attr_list, mode, object_statuses);
}

sai_status_t l_get_nat_entries_attribute(uint32_t object_count, const sai_nat_entry_t *nat_entry, const uint32_t *attr_count, sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__; common_entry_t entry = {.nat_entry = nat_entry}; 
	return translator->get_attribute_bulk(SAI_OBJECT_TYPE_NAT_ENTRY, object_count, entry, attr_count, attr_list, mode, object_statuses);
}

sai_status_t l_create_nat_zone_counter(sai_object_id_t *nat_zone_counter_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_NAT_ZONE_COUNTER, nat_zone_counter_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_nat_zone_counter(sai_object_id_t nat_zone_counter_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_NAT_ZONE_COUNTER, nat_zone_counter_id);
}

sai_status_t l_set_nat_zone_counter_attribute(sai_object_id_t nat_zone_counter_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_NAT_ZONE_COUNTER, nat_zone_counter_id, attr);
}

sai_status_t l_get_nat_zone_counter_attribute(sai_object_id_t nat_zone_counter_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_NAT_ZONE_COUNTER, nat_zone_counter_id, attr_count, attr_list);
}

