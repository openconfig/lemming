
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

#include "dataplane/standalone/sai/vlan.h"
#include "dataplane/standalone/log/log.h"

const sai_vlan_api_t l_vlan = {
	.create_vlan = l_create_vlan,
	.remove_vlan = l_remove_vlan,
	.set_vlan_attribute = l_set_vlan_attribute,
	.get_vlan_attribute = l_get_vlan_attribute,
	.create_vlan_member = l_create_vlan_member,
	.remove_vlan_member = l_remove_vlan_member,
	.set_vlan_member_attribute = l_set_vlan_member_attribute,
	.get_vlan_member_attribute = l_get_vlan_member_attribute,
	.create_vlan_members = l_create_vlan_members,
	.remove_vlan_members = l_remove_vlan_members,
	.get_vlan_stats = l_get_vlan_stats,
	.get_vlan_stats_ext = l_get_vlan_stats_ext,
	.clear_vlan_stats = l_clear_vlan_stats,
};


sai_status_t l_create_vlan(sai_object_id_t *vlan_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_remove_vlan(sai_object_id_t vlan_id) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_set_vlan_attribute(sai_object_id_t vlan_id, const sai_attribute_t *attr) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_get_vlan_attribute(sai_object_id_t vlan_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_create_vlan_member(sai_object_id_t *vlan_member_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_remove_vlan_member(sai_object_id_t vlan_member_id) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_set_vlan_member_attribute(sai_object_id_t vlan_member_id, const sai_attribute_t *attr) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_get_vlan_member_attribute(sai_object_id_t vlan_member_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_create_vlan_members(sai_object_id_t switch_id, uint32_t object_count, const uint32_t *attr_count, const sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_object_id_t *object_id, sai_status_t *object_statuses) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_remove_vlan_members(uint32_t object_count, const sai_object_id_t *object_id, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_get_vlan_stats(sai_object_id_t vlan_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_get_vlan_stats_ext(sai_object_id_t vlan_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_clear_vlan_stats(sai_object_id_t vlan_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


