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

#include "dataplane/standalone/sai/tunnel.h"

#include <glog/logging.h>

#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_tunnel_api_t l_tunnel = {
    .create_tunnel_map = l_create_tunnel_map,
    .remove_tunnel_map = l_remove_tunnel_map,
    .set_tunnel_map_attribute = l_set_tunnel_map_attribute,
    .get_tunnel_map_attribute = l_get_tunnel_map_attribute,
    .create_tunnel = l_create_tunnel,
    .remove_tunnel = l_remove_tunnel,
    .set_tunnel_attribute = l_set_tunnel_attribute,
    .get_tunnel_attribute = l_get_tunnel_attribute,
    .get_tunnel_stats = l_get_tunnel_stats,
    .get_tunnel_stats_ext = l_get_tunnel_stats_ext,
    .clear_tunnel_stats = l_clear_tunnel_stats,
    .create_tunnel_term_table_entry = l_create_tunnel_term_table_entry,
    .remove_tunnel_term_table_entry = l_remove_tunnel_term_table_entry,
    .set_tunnel_term_table_entry_attribute =
        l_set_tunnel_term_table_entry_attribute,
    .get_tunnel_term_table_entry_attribute =
        l_get_tunnel_term_table_entry_attribute,
    .create_tunnel_map_entry = l_create_tunnel_map_entry,
    .remove_tunnel_map_entry = l_remove_tunnel_map_entry,
    .set_tunnel_map_entry_attribute = l_set_tunnel_map_entry_attribute,
    .get_tunnel_map_entry_attribute = l_get_tunnel_map_entry_attribute,
};

sai_status_t l_create_tunnel_map(sai_object_id_t *tunnel_map_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_TUNNEL_MAP, tunnel_map_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tunnel_map(sai_object_id_t tunnel_map_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_TUNNEL_MAP, tunnel_map_id);
}

sai_status_t l_set_tunnel_map_attribute(sai_object_id_t tunnel_map_id,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_TUNNEL_MAP, tunnel_map_id,
                                   attr);
}

sai_status_t l_get_tunnel_map_attribute(sai_object_id_t tunnel_map_id,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_TUNNEL_MAP, tunnel_map_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_tunnel(sai_object_id_t *tunnel_id,
                             sai_object_id_t switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_TUNNEL, tunnel_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_tunnel(sai_object_id_t tunnel_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_TUNNEL, tunnel_id);
}

sai_status_t l_set_tunnel_attribute(sai_object_id_t tunnel_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_TUNNEL, tunnel_id, attr);
}

sai_status_t l_get_tunnel_attribute(sai_object_id_t tunnel_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_TUNNEL, tunnel_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_tunnel_stats(sai_object_id_t tunnel_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids,
                                uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_stats(SAI_OBJECT_TYPE_TUNNEL, tunnel_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_tunnel_stats_ext(sai_object_id_t tunnel_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t *counter_ids,
                                    sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_stats_ext(SAI_OBJECT_TYPE_TUNNEL, tunnel_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_tunnel_stats(sai_object_id_t tunnel_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->clear_stats(SAI_OBJECT_TYPE_TUNNEL, tunnel_id,
                                 number_of_counters, counter_ids);
}

sai_status_t l_create_tunnel_term_table_entry(
    sai_object_id_t *tunnel_term_table_entry_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY,
                            tunnel_term_table_entry_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_tunnel_term_table_entry(
    sai_object_id_t tunnel_term_table_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY,
                            tunnel_term_table_entry_id);
}

sai_status_t l_set_tunnel_term_table_entry_attribute(
    sai_object_id_t tunnel_term_table_entry_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY,
                                   tunnel_term_table_entry_id, attr);
}

sai_status_t l_get_tunnel_term_table_entry_attribute(
    sai_object_id_t tunnel_term_table_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY,
                                   tunnel_term_table_entry_id, attr_count,
                                   attr_list);
}

sai_status_t l_create_tunnel_map_entry(sai_object_id_t *tunnel_map_entry_id,
                                       sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_TUNNEL_MAP_ENTRY,
                            tunnel_map_entry_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_tunnel_map_entry(sai_object_id_t tunnel_map_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_TUNNEL_MAP_ENTRY,
                            tunnel_map_entry_id);
}

sai_status_t l_set_tunnel_map_entry_attribute(
    sai_object_id_t tunnel_map_entry_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_TUNNEL_MAP_ENTRY,
                                   tunnel_map_entry_id, attr);
}

sai_status_t l_get_tunnel_map_entry_attribute(
    sai_object_id_t tunnel_map_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_TUNNEL_MAP_ENTRY,
                                   tunnel_map_entry_id, attr_count, attr_list);
}
