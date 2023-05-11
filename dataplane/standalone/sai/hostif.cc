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

#include "dataplane/standalone/sai/hostif.h"
#include <glog/logging.h>
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_hostif_api_t l_hostif = {
	.create_hostif = l_create_hostif,
	.remove_hostif = l_remove_hostif,
	.set_hostif_attribute = l_set_hostif_attribute,
	.get_hostif_attribute = l_get_hostif_attribute,
	.create_hostif_table_entry = l_create_hostif_table_entry,
	.remove_hostif_table_entry = l_remove_hostif_table_entry,
	.set_hostif_table_entry_attribute = l_set_hostif_table_entry_attribute,
	.get_hostif_table_entry_attribute = l_get_hostif_table_entry_attribute,
	.create_hostif_trap_group = l_create_hostif_trap_group,
	.remove_hostif_trap_group = l_remove_hostif_trap_group,
	.set_hostif_trap_group_attribute = l_set_hostif_trap_group_attribute,
	.get_hostif_trap_group_attribute = l_get_hostif_trap_group_attribute,
	.create_hostif_trap = l_create_hostif_trap,
	.remove_hostif_trap = l_remove_hostif_trap,
	.set_hostif_trap_attribute = l_set_hostif_trap_attribute,
	.get_hostif_trap_attribute = l_get_hostif_trap_attribute,
	.create_hostif_user_defined_trap = l_create_hostif_user_defined_trap,
	.remove_hostif_user_defined_trap = l_remove_hostif_user_defined_trap,
	.set_hostif_user_defined_trap_attribute = l_set_hostif_user_defined_trap_attribute,
	.get_hostif_user_defined_trap_attribute = l_get_hostif_user_defined_trap_attribute,
	.recv_hostif_packet = l_recv_hostif_packet,
	.send_hostif_packet = l_send_hostif_packet,
	.allocate_hostif_packet = l_allocate_hostif_packet,
	.free_hostif_packet = l_free_hostif_packet,
};


sai_status_t l_create_hostif(sai_object_id_t *hostif_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_HOSTIF, hostif_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_hostif(sai_object_id_t hostif_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_HOSTIF, hostif_id);
}

sai_status_t l_set_hostif_attribute(sai_object_id_t hostif_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_HOSTIF, hostif_id, attr);
}

sai_status_t l_get_hostif_attribute(sai_object_id_t hostif_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_HOSTIF, hostif_id, attr_count, attr_list);
}

sai_status_t l_create_hostif_table_entry(sai_object_id_t *hostif_table_entry_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY, hostif_table_entry_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_hostif_table_entry(sai_object_id_t hostif_table_entry_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY, hostif_table_entry_id);
}

sai_status_t l_set_hostif_table_entry_attribute(sai_object_id_t hostif_table_entry_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY, hostif_table_entry_id, attr);
}

sai_status_t l_get_hostif_table_entry_attribute(sai_object_id_t hostif_table_entry_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY, hostif_table_entry_id, attr_count, attr_list);
}

sai_status_t l_create_hostif_trap_group(sai_object_id_t *hostif_trap_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP, hostif_trap_group_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_hostif_trap_group(sai_object_id_t hostif_trap_group_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP, hostif_trap_group_id);
}

sai_status_t l_set_hostif_trap_group_attribute(sai_object_id_t hostif_trap_group_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP, hostif_trap_group_id, attr);
}

sai_status_t l_get_hostif_trap_group_attribute(sai_object_id_t hostif_trap_group_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP, hostif_trap_group_id, attr_count, attr_list);
}

sai_status_t l_create_hostif_trap(sai_object_id_t *hostif_trap_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_HOSTIF_TRAP, hostif_trap_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_hostif_trap(sai_object_id_t hostif_trap_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_HOSTIF_TRAP, hostif_trap_id);
}

sai_status_t l_set_hostif_trap_attribute(sai_object_id_t hostif_trap_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_HOSTIF_TRAP, hostif_trap_id, attr);
}

sai_status_t l_get_hostif_trap_attribute(sai_object_id_t hostif_trap_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_HOSTIF_TRAP, hostif_trap_id, attr_count, attr_list);
}

sai_status_t l_create_hostif_user_defined_trap(sai_object_id_t *hostif_user_defined_trap_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP, hostif_user_defined_trap_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_hostif_user_defined_trap(sai_object_id_t hostif_user_defined_trap_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP, hostif_user_defined_trap_id);
}

sai_status_t l_set_hostif_user_defined_trap_attribute(sai_object_id_t hostif_user_defined_trap_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP, hostif_user_defined_trap_id, attr);
}

sai_status_t l_get_hostif_user_defined_trap_attribute(sai_object_id_t hostif_user_defined_trap_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP, hostif_user_defined_trap_id, attr_count, attr_list);
}

sai_status_t l_recv_hostif_packet(sai_object_id_t hostif_id, sai_size_t *buffer_size, void *buffer, uint32_t *attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_send_hostif_packet(sai_object_id_t hostif_id, sai_size_t buffer_size, const void *buffer, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_allocate_hostif_packet(sai_object_id_t hostif_id, sai_size_t buffer_size, void **buffer, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_free_hostif_packet(sai_object_id_t hostif_id, void *buffer) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return SAI_STATUS_NOT_IMPLEMENTED;
}

