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

#ifndef DATAPLANE_STANDALONE_SAI_IPSEC_H_
#define DATAPLANE_STANDALONE_SAI_IPSEC_H_

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

extern const sai_ipsec_api_t l_ipsec;

sai_status_t l_create_ipsec(sai_object_id_t* ipsec_id,
                            sai_object_id_t switch_id, uint32_t attr_count,
                            const sai_attribute_t* attr_list);

sai_status_t l_remove_ipsec(sai_object_id_t ipsec_id);

sai_status_t l_set_ipsec_attribute(sai_object_id_t ipsec_id,
                                   const sai_attribute_t* attr);

sai_status_t l_get_ipsec_attribute(sai_object_id_t ipsec_id,
                                   uint32_t attr_count,
                                   sai_attribute_t* attr_list);

sai_status_t l_create_ipsec_port(sai_object_id_t* ipsec_port_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t* attr_list);

sai_status_t l_remove_ipsec_port(sai_object_id_t ipsec_port_id);

sai_status_t l_set_ipsec_port_attribute(sai_object_id_t ipsec_port_id,
                                        const sai_attribute_t* attr);

sai_status_t l_get_ipsec_port_attribute(sai_object_id_t ipsec_port_id,
                                        uint32_t attr_count,
                                        sai_attribute_t* attr_list);

sai_status_t l_get_ipsec_port_stats(sai_object_id_t ipsec_port_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t* counter_ids,
                                    uint64_t* counters);

sai_status_t l_get_ipsec_port_stats_ext(sai_object_id_t ipsec_port_id,
                                        uint32_t number_of_counters,
                                        const sai_stat_id_t* counter_ids,
                                        sai_stats_mode_t mode,
                                        uint64_t* counters);

sai_status_t l_clear_ipsec_port_stats(sai_object_id_t ipsec_port_id,
                                      uint32_t number_of_counters,
                                      const sai_stat_id_t* counter_ids);

sai_status_t l_create_ipsec_sa(sai_object_id_t* ipsec_sa_id,
                               sai_object_id_t switch_id, uint32_t attr_count,
                               const sai_attribute_t* attr_list);

sai_status_t l_remove_ipsec_sa(sai_object_id_t ipsec_sa_id);

sai_status_t l_set_ipsec_sa_attribute(sai_object_id_t ipsec_sa_id,
                                      const sai_attribute_t* attr);

sai_status_t l_get_ipsec_sa_attribute(sai_object_id_t ipsec_sa_id,
                                      uint32_t attr_count,
                                      sai_attribute_t* attr_list);

sai_status_t l_get_ipsec_sa_stats(sai_object_id_t ipsec_sa_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t* counter_ids,
                                  uint64_t* counters);

sai_status_t l_get_ipsec_sa_stats_ext(sai_object_id_t ipsec_sa_id,
                                      uint32_t number_of_counters,
                                      const sai_stat_id_t* counter_ids,
                                      sai_stats_mode_t mode,
                                      uint64_t* counters);

sai_status_t l_clear_ipsec_sa_stats(sai_object_id_t ipsec_sa_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t* counter_ids);

#endif  // DATAPLANE_STANDALONE_SAI_IPSEC_H_
