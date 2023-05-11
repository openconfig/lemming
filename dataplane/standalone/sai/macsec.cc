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

#include "dataplane/standalone/sai/macsec.h"
#include <glog/logging.h>
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_macsec_api_t l_macsec = {
	.create_macsec = l_create_macsec,
	.remove_macsec = l_remove_macsec,
	.set_macsec_attribute = l_set_macsec_attribute,
	.get_macsec_attribute = l_get_macsec_attribute,
	.create_macsec_port = l_create_macsec_port,
	.remove_macsec_port = l_remove_macsec_port,
	.set_macsec_port_attribute = l_set_macsec_port_attribute,
	.get_macsec_port_attribute = l_get_macsec_port_attribute,
	.get_macsec_port_stats = l_get_macsec_port_stats,
	.get_macsec_port_stats_ext = l_get_macsec_port_stats_ext,
	.clear_macsec_port_stats = l_clear_macsec_port_stats,
	.create_macsec_flow = l_create_macsec_flow,
	.remove_macsec_flow = l_remove_macsec_flow,
	.set_macsec_flow_attribute = l_set_macsec_flow_attribute,
	.get_macsec_flow_attribute = l_get_macsec_flow_attribute,
	.get_macsec_flow_stats = l_get_macsec_flow_stats,
	.get_macsec_flow_stats_ext = l_get_macsec_flow_stats_ext,
	.clear_macsec_flow_stats = l_clear_macsec_flow_stats,
	.create_macsec_sc = l_create_macsec_sc,
	.remove_macsec_sc = l_remove_macsec_sc,
	.set_macsec_sc_attribute = l_set_macsec_sc_attribute,
	.get_macsec_sc_attribute = l_get_macsec_sc_attribute,
	.get_macsec_sc_stats = l_get_macsec_sc_stats,
	.get_macsec_sc_stats_ext = l_get_macsec_sc_stats_ext,
	.clear_macsec_sc_stats = l_clear_macsec_sc_stats,
	.create_macsec_sa = l_create_macsec_sa,
	.remove_macsec_sa = l_remove_macsec_sa,
	.set_macsec_sa_attribute = l_set_macsec_sa_attribute,
	.get_macsec_sa_attribute = l_get_macsec_sa_attribute,
	.get_macsec_sa_stats = l_get_macsec_sa_stats,
	.get_macsec_sa_stats_ext = l_get_macsec_sa_stats_ext,
	.clear_macsec_sa_stats = l_clear_macsec_sa_stats,
};


sai_status_t l_create_macsec(sai_object_id_t *macsec_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_MACSEC, macsec_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_macsec(sai_object_id_t macsec_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_MACSEC, macsec_id);
}

sai_status_t l_set_macsec_attribute(sai_object_id_t macsec_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_MACSEC, macsec_id, attr);
}

sai_status_t l_get_macsec_attribute(sai_object_id_t macsec_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_MACSEC, macsec_id, attr_count, attr_list);
}

sai_status_t l_create_macsec_port(sai_object_id_t *macsec_port_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_MACSEC_PORT, macsec_port_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_macsec_port(sai_object_id_t macsec_port_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_MACSEC_PORT, macsec_port_id);
}

sai_status_t l_set_macsec_port_attribute(sai_object_id_t macsec_port_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_MACSEC_PORT, macsec_port_id, attr);
}

sai_status_t l_get_macsec_port_attribute(sai_object_id_t macsec_port_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_MACSEC_PORT, macsec_port_id, attr_count, attr_list);
}

sai_status_t l_get_macsec_port_stats(sai_object_id_t macsec_port_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats(SAI_OBJECT_TYPE_MACSEC_PORT, macsec_port_id, number_of_counters, counter_ids, counters);
}

sai_status_t l_get_macsec_port_stats_ext(sai_object_id_t macsec_port_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats_ext(SAI_OBJECT_TYPE_MACSEC_PORT, macsec_port_id, number_of_counters, counter_ids, mode, counters);
}

sai_status_t l_clear_macsec_port_stats(sai_object_id_t macsec_port_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->clear_stats(SAI_OBJECT_TYPE_MACSEC_PORT, macsec_port_id, number_of_counters, counter_ids);
}

sai_status_t l_create_macsec_flow(sai_object_id_t *macsec_flow_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_MACSEC_FLOW, macsec_flow_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_macsec_flow(sai_object_id_t macsec_flow_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_MACSEC_FLOW, macsec_flow_id);
}

sai_status_t l_set_macsec_flow_attribute(sai_object_id_t macsec_flow_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_MACSEC_FLOW, macsec_flow_id, attr);
}

sai_status_t l_get_macsec_flow_attribute(sai_object_id_t macsec_flow_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_MACSEC_FLOW, macsec_flow_id, attr_count, attr_list);
}

sai_status_t l_get_macsec_flow_stats(sai_object_id_t macsec_flow_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats(SAI_OBJECT_TYPE_MACSEC_FLOW, macsec_flow_id, number_of_counters, counter_ids, counters);
}

sai_status_t l_get_macsec_flow_stats_ext(sai_object_id_t macsec_flow_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats_ext(SAI_OBJECT_TYPE_MACSEC_FLOW, macsec_flow_id, number_of_counters, counter_ids, mode, counters);
}

sai_status_t l_clear_macsec_flow_stats(sai_object_id_t macsec_flow_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->clear_stats(SAI_OBJECT_TYPE_MACSEC_FLOW, macsec_flow_id, number_of_counters, counter_ids);
}

sai_status_t l_create_macsec_sc(sai_object_id_t *macsec_sc_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_MACSEC_SC, macsec_sc_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_macsec_sc(sai_object_id_t macsec_sc_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_MACSEC_SC, macsec_sc_id);
}

sai_status_t l_set_macsec_sc_attribute(sai_object_id_t macsec_sc_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_MACSEC_SC, macsec_sc_id, attr);
}

sai_status_t l_get_macsec_sc_attribute(sai_object_id_t macsec_sc_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_MACSEC_SC, macsec_sc_id, attr_count, attr_list);
}

sai_status_t l_get_macsec_sc_stats(sai_object_id_t macsec_sc_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats(SAI_OBJECT_TYPE_MACSEC_SC, macsec_sc_id, number_of_counters, counter_ids, counters);
}

sai_status_t l_get_macsec_sc_stats_ext(sai_object_id_t macsec_sc_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats_ext(SAI_OBJECT_TYPE_MACSEC_SC, macsec_sc_id, number_of_counters, counter_ids, mode, counters);
}

sai_status_t l_clear_macsec_sc_stats(sai_object_id_t macsec_sc_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->clear_stats(SAI_OBJECT_TYPE_MACSEC_SC, macsec_sc_id, number_of_counters, counter_ids);
}

sai_status_t l_create_macsec_sa(sai_object_id_t *macsec_sa_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_MACSEC_SA, macsec_sa_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_macsec_sa(sai_object_id_t macsec_sa_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_MACSEC_SA, macsec_sa_id);
}

sai_status_t l_set_macsec_sa_attribute(sai_object_id_t macsec_sa_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_MACSEC_SA, macsec_sa_id, attr);
}

sai_status_t l_get_macsec_sa_attribute(sai_object_id_t macsec_sa_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_MACSEC_SA, macsec_sa_id, attr_count, attr_list);
}

sai_status_t l_get_macsec_sa_stats(sai_object_id_t macsec_sa_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats(SAI_OBJECT_TYPE_MACSEC_SA, macsec_sa_id, number_of_counters, counter_ids, counters);
}

sai_status_t l_get_macsec_sa_stats_ext(sai_object_id_t macsec_sa_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_stats_ext(SAI_OBJECT_TYPE_MACSEC_SA, macsec_sa_id, number_of_counters, counter_ids, mode, counters);
}

sai_status_t l_clear_macsec_sa_stats(sai_object_id_t macsec_sa_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->clear_stats(SAI_OBJECT_TYPE_MACSEC_SA, macsec_sa_id, number_of_counters, counter_ids);
}

