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

#include "dataplane/standalone/sai/buffer.h"
#include <glog/logging.h>
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_buffer_api_t l_buffer = {
	.create_buffer_pool = l_create_buffer_pool,
	.remove_buffer_pool = l_remove_buffer_pool,
	.set_buffer_pool_attribute = l_set_buffer_pool_attribute,
	.get_buffer_pool_attribute = l_get_buffer_pool_attribute,
	.get_buffer_pool_stats = l_get_buffer_pool_stats,
	.get_buffer_pool_stats_ext = l_get_buffer_pool_stats_ext,
	.clear_buffer_pool_stats = l_clear_buffer_pool_stats,
	.create_ingress_priority_group = l_create_ingress_priority_group,
	.remove_ingress_priority_group = l_remove_ingress_priority_group,
	.set_ingress_priority_group_attribute = l_set_ingress_priority_group_attribute,
	.get_ingress_priority_group_attribute = l_get_ingress_priority_group_attribute,
	.get_ingress_priority_group_stats = l_get_ingress_priority_group_stats,
	.get_ingress_priority_group_stats_ext = l_get_ingress_priority_group_stats_ext,
	.clear_ingress_priority_group_stats = l_clear_ingress_priority_group_stats,
	.create_buffer_profile = l_create_buffer_profile,
	.remove_buffer_profile = l_remove_buffer_profile,
	.set_buffer_profile_attribute = l_set_buffer_profile_attribute,
	.get_buffer_profile_attribute = l_get_buffer_profile_attribute,
};


sai_status_t l_create_buffer_pool(sai_object_id_t *buffer_pool_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_buffer_pool(sai_object_id_t buffer_pool_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id);
}

sai_status_t l_set_buffer_pool_attribute(sai_object_id_t buffer_pool_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id, attr);
}

sai_status_t l_get_buffer_pool_attribute(sai_object_id_t buffer_pool_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id, attr_count, attr_list);
}

sai_status_t l_get_buffer_pool_stats(sai_object_id_t buffer_pool_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id, number_of_counters, counter_ids, counters);
}

sai_status_t l_get_buffer_pool_stats_ext(sai_object_id_t buffer_pool_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats_ext(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id, number_of_counters, counter_ids, mode, counters);
}

sai_status_t l_clear_buffer_pool_stats(sai_object_id_t buffer_pool_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->clear_stats(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id, number_of_counters, counter_ids);
}

sai_status_t l_create_ingress_priority_group(sai_object_id_t *ingress_priority_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP, ingress_priority_group_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_ingress_priority_group(sai_object_id_t ingress_priority_group_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP, ingress_priority_group_id);
}

sai_status_t l_set_ingress_priority_group_attribute(sai_object_id_t ingress_priority_group_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP, ingress_priority_group_id, attr);
}

sai_status_t l_get_ingress_priority_group_attribute(sai_object_id_t ingress_priority_group_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP, ingress_priority_group_id, attr_count, attr_list);
}

sai_status_t l_get_ingress_priority_group_stats(sai_object_id_t ingress_priority_group_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP, ingress_priority_group_id, number_of_counters, counter_ids, counters);
}

sai_status_t l_get_ingress_priority_group_stats_ext(sai_object_id_t ingress_priority_group_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats_ext(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP, ingress_priority_group_id, number_of_counters, counter_ids, mode, counters);
}

sai_status_t l_clear_ingress_priority_group_stats(sai_object_id_t ingress_priority_group_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->clear_stats(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP, ingress_priority_group_id, number_of_counters, counter_ids);
}

sai_status_t l_create_buffer_profile(sai_object_id_t *buffer_profile_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_BUFFER_PROFILE, buffer_profile_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_buffer_profile(sai_object_id_t buffer_profile_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_BUFFER_PROFILE, buffer_profile_id);
}

sai_status_t l_set_buffer_profile_attribute(sai_object_id_t buffer_profile_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_BUFFER_PROFILE, buffer_profile_id, attr);
}

sai_status_t l_get_buffer_profile_attribute(sai_object_id_t buffer_profile_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_BUFFER_PROFILE, buffer_profile_id, attr_count, attr_list);
}

