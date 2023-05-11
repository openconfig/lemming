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

#ifndef DATAPLANE_STANDALONE_SAI_BRIDGE_H_
#define DATAPLANE_STANDALONE_SAI_BRIDGE_H_

extern "C" {
	#include "inc/sai.h"
}

extern const sai_bridge_api_t l_bridge;


sai_status_t l_create_bridge(sai_object_id_t *bridge_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_bridge(sai_object_id_t bridge_id);

sai_status_t l_set_bridge_attribute(sai_object_id_t bridge_id, const sai_attribute_t *attr);

sai_status_t l_get_bridge_attribute(sai_object_id_t bridge_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_get_bridge_stats(sai_object_id_t bridge_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters);

sai_status_t l_get_bridge_stats_ext(sai_object_id_t bridge_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters);

sai_status_t l_clear_bridge_stats(sai_object_id_t bridge_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids);

sai_status_t l_create_bridge_port(sai_object_id_t *bridge_port_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_bridge_port(sai_object_id_t bridge_port_id);

sai_status_t l_set_bridge_port_attribute(sai_object_id_t bridge_port_id, const sai_attribute_t *attr);

sai_status_t l_get_bridge_port_attribute(sai_object_id_t bridge_port_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_get_bridge_port_stats(sai_object_id_t bridge_port_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters);

sai_status_t l_get_bridge_port_stats_ext(sai_object_id_t bridge_port_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters);

sai_status_t l_clear_bridge_port_stats(sai_object_id_t bridge_port_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids);


#endif  // DATAPLANE_STANDALONE_SAI_BRIDGE_H_
