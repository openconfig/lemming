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

#ifndef DATAPLANE_STANDALONE_SAI_MACSEC_H_
#define DATAPLANE_STANDALONE_SAI_MACSEC_H_

extern "C" {
#include "inc/sai.h"
#include "experimental/saiextensions.h"
}

extern const sai_macsec_api_t l_macsec;

sai_status_t l_create_macsec(sai_object_id_t *macsec_id,
                             sai_object_id_t switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list);

sai_status_t l_remove_macsec(sai_object_id_t macsec_id);

sai_status_t l_set_macsec_attribute(sai_object_id_t macsec_id,
                                    const sai_attribute_t *attr);

sai_status_t l_get_macsec_attribute(sai_object_id_t macsec_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list);

sai_status_t l_create_macsec_port(sai_object_id_t *macsec_port_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list);

sai_status_t l_remove_macsec_port(sai_object_id_t macsec_port_id);

sai_status_t l_set_macsec_port_attribute(sai_object_id_t macsec_port_id,
                                         const sai_attribute_t *attr);

sai_status_t l_get_macsec_port_attribute(sai_object_id_t macsec_port_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list);

sai_status_t l_get_macsec_port_stats(sai_object_id_t macsec_port_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     uint64_t *counters);

sai_status_t l_get_macsec_port_stats_ext(sai_object_id_t macsec_port_id,
                                         uint32_t number_of_counters,
                                         const sai_stat_id_t *counter_ids,
                                         sai_stats_mode_t mode,
                                         uint64_t *counters);

sai_status_t l_clear_macsec_port_stats(sai_object_id_t macsec_port_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids);

sai_status_t l_create_macsec_flow(sai_object_id_t *macsec_flow_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list);

sai_status_t l_remove_macsec_flow(sai_object_id_t macsec_flow_id);

sai_status_t l_set_macsec_flow_attribute(sai_object_id_t macsec_flow_id,
                                         const sai_attribute_t *attr);

sai_status_t l_get_macsec_flow_attribute(sai_object_id_t macsec_flow_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list);

sai_status_t l_get_macsec_flow_stats(sai_object_id_t macsec_flow_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     uint64_t *counters);

sai_status_t l_get_macsec_flow_stats_ext(sai_object_id_t macsec_flow_id,
                                         uint32_t number_of_counters,
                                         const sai_stat_id_t *counter_ids,
                                         sai_stats_mode_t mode,
                                         uint64_t *counters);

sai_status_t l_clear_macsec_flow_stats(sai_object_id_t macsec_flow_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids);

sai_status_t l_create_macsec_sc(sai_object_id_t *macsec_sc_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list);

sai_status_t l_remove_macsec_sc(sai_object_id_t macsec_sc_id);

sai_status_t l_set_macsec_sc_attribute(sai_object_id_t macsec_sc_id,
                                       const sai_attribute_t *attr);

sai_status_t l_get_macsec_sc_attribute(sai_object_id_t macsec_sc_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list);

sai_status_t l_get_macsec_sc_stats(sai_object_id_t macsec_sc_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids,
                                   uint64_t *counters);

sai_status_t l_get_macsec_sc_stats_ext(sai_object_id_t macsec_sc_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids,
                                       sai_stats_mode_t mode,
                                       uint64_t *counters);

sai_status_t l_clear_macsec_sc_stats(sai_object_id_t macsec_sc_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids);

sai_status_t l_create_macsec_sa(sai_object_id_t *macsec_sa_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list);

sai_status_t l_remove_macsec_sa(sai_object_id_t macsec_sa_id);

sai_status_t l_set_macsec_sa_attribute(sai_object_id_t macsec_sa_id,
                                       const sai_attribute_t *attr);

sai_status_t l_get_macsec_sa_attribute(sai_object_id_t macsec_sa_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list);

sai_status_t l_get_macsec_sa_stats(sai_object_id_t macsec_sa_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids,
                                   uint64_t *counters);

sai_status_t l_get_macsec_sa_stats_ext(sai_object_id_t macsec_sa_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids,
                                       sai_stats_mode_t mode,
                                       uint64_t *counters);

sai_status_t l_clear_macsec_sa_stats(sai_object_id_t macsec_sa_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids);

#endif  // DATAPLANE_STANDALONE_SAI_MACSEC_H_
