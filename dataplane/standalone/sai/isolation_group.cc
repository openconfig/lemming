
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

#include "dataplane/standalone/sai/isolation_group.h"
#include "dataplane/standalone/log/log.h"

const sai_isolation_group_api_t l_isolation_group = {
	.create_isolation_group = l_create_isolation_group,
	.remove_isolation_group = l_remove_isolation_group,
	.set_isolation_group_attribute = l_set_isolation_group_attribute,
	.get_isolation_group_attribute = l_get_isolation_group_attribute,
	.create_isolation_group_member = l_create_isolation_group_member,
	.remove_isolation_group_member = l_remove_isolation_group_member,
	.set_isolation_group_member_attribute = l_set_isolation_group_member_attribute,
	.get_isolation_group_member_attribute = l_get_isolation_group_member_attribute,
};


sai_status_t l_create_isolation_group(sai_object_id_t *isolation_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_remove_isolation_group(sai_object_id_t isolation_group_id) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_set_isolation_group_attribute(sai_object_id_t isolation_group_id, const sai_attribute_t *attr) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_get_isolation_group_attribute(sai_object_id_t isolation_group_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_create_isolation_group_member(sai_object_id_t *isolation_group_member_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_remove_isolation_group_member(sai_object_id_t isolation_group_member_id) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_set_isolation_group_member_attribute(sai_object_id_t isolation_group_member_id, const sai_attribute_t *attr) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_get_isolation_group_member_attribute(sai_object_id_t isolation_group_member_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}

