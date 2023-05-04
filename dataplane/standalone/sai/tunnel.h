
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

#ifndef DATAPLANE_STANDALONE_SAI_TUNNEL_H_
#define DATAPLANE_STANDALONE_SAI_TUNNEL_H_

extern "C" {
	#include "inc/sai.h"
}

extern const sai_tunnel_api_t l_tunnel;


sai_status_t l_create_tunnel_map(sai_object_id_t *tunnel_map_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_tunnel_map(sai_object_id_t tunnel_map_id);

sai_status_t l_set_tunnel_map_attribute(sai_object_id_t tunnel_map_id, const sai_attribute_t *attr);

sai_status_t l_get_tunnel_map_attribute(sai_object_id_t tunnel_map_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_tunnel(sai_object_id_t *tunnel_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_tunnel(sai_object_id_t tunnel_id);

sai_status_t l_set_tunnel_attribute(sai_object_id_t tunnel_id, const sai_attribute_t *attr);

sai_status_t l_get_tunnel_attribute(sai_object_id_t tunnel_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_get_tunnel_stats(sai_object_id_t tunnel_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters);

sai_status_t l_get_tunnel_stats_ext(sai_object_id_t tunnel_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters);

sai_status_t l_clear_tunnel_stats(sai_object_id_t tunnel_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids);

sai_status_t l_create_tunnel_term_table_entry(sai_object_id_t *tunnel_term_table_entry_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_tunnel_term_table_entry(sai_object_id_t tunnel_term_table_entry_id);

sai_status_t l_set_tunnel_term_table_entry_attribute(sai_object_id_t tunnel_term_table_entry_id, const sai_attribute_t *attr);

sai_status_t l_get_tunnel_term_table_entry_attribute(sai_object_id_t tunnel_term_table_entry_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_tunnel_map_entry(sai_object_id_t *tunnel_map_entry_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_tunnel_map_entry(sai_object_id_t tunnel_map_entry_id);

sai_status_t l_set_tunnel_map_entry_attribute(sai_object_id_t tunnel_map_entry_id, const sai_attribute_t *attr);

sai_status_t l_get_tunnel_map_entry_attribute(sai_object_id_t tunnel_map_entry_id, uint32_t attr_count, sai_attribute_t *attr_list);


#endif  // DATAPLANE_STANDALONE_SAI_TUNNEL_H_
