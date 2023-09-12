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

#ifndef DATAPLANE_STANDALONE_SAI_PORT_H_
#define DATAPLANE_STANDALONE_SAI_PORT_H_

extern "C" {
#include "inc/sai.h"
#include "experimental/saiextensions.h"
}

extern const sai_port_api_t l_port;

sai_status_t l_create_port(sai_object_id_t *port_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list);

sai_status_t l_remove_port(sai_object_id_t port_id);

sai_status_t l_set_port_attribute(sai_object_id_t port_id,
                                  const sai_attribute_t *attr);

sai_status_t l_get_port_attribute(sai_object_id_t port_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list);

sai_status_t l_get_port_stats(sai_object_id_t port_id,
                              uint32_t number_of_counters,
                              const sai_stat_id_t *counter_ids,
                              uint64_t *counters);

sai_status_t l_get_port_stats_ext(sai_object_id_t port_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids,
                                  sai_stats_mode_t mode, uint64_t *counters);

sai_status_t l_clear_port_stats(sai_object_id_t port_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids);

sai_status_t l_clear_port_all_stats(sai_object_id_t port_id);

sai_status_t l_create_port_pool(sai_object_id_t *port_pool_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list);

sai_status_t l_remove_port_pool(sai_object_id_t port_pool_id);

sai_status_t l_set_port_pool_attribute(sai_object_id_t port_pool_id,
                                       const sai_attribute_t *attr);

sai_status_t l_get_port_pool_attribute(sai_object_id_t port_pool_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list);

sai_status_t l_get_port_pool_stats(sai_object_id_t port_pool_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids,
                                   uint64_t *counters);

sai_status_t l_get_port_pool_stats_ext(sai_object_id_t port_pool_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids,
                                       sai_stats_mode_t mode,
                                       uint64_t *counters);

sai_status_t l_clear_port_pool_stats(sai_object_id_t port_pool_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids);

sai_status_t l_create_port_connector(sai_object_id_t *port_connector_id,
                                     sai_object_id_t switch_id,
                                     uint32_t attr_count,
                                     const sai_attribute_t *attr_list);

sai_status_t l_remove_port_connector(sai_object_id_t port_connector_id);

sai_status_t l_set_port_connector_attribute(sai_object_id_t port_connector_id,
                                            const sai_attribute_t *attr);

sai_status_t l_get_port_connector_attribute(sai_object_id_t port_connector_id,
                                            uint32_t attr_count,
                                            sai_attribute_t *attr_list);

sai_status_t l_create_port_serdes(sai_object_id_t *port_serdes_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list);

sai_status_t l_remove_port_serdes(sai_object_id_t port_serdes_id);

sai_status_t l_set_port_serdes_attribute(sai_object_id_t port_serdes_id,
                                         const sai_attribute_t *attr);

sai_status_t l_get_port_serdes_attribute(sai_object_id_t port_serdes_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list);

#endif  // DATAPLANE_STANDALONE_SAI_PORT_H_
