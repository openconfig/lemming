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

#ifndef DATAPLANE_STANDALONE_SAI_TAM_H_
#define DATAPLANE_STANDALONE_SAI_TAM_H_

extern "C" {
#include "inc/sai.h"
#include "experimental/saiextensions.h"
}

extern const sai_tam_api_t l_tam;

sai_status_t l_create_tam(sai_object_id_t *tam_id, sai_object_id_t switch_id,
                          uint32_t attr_count,
                          const sai_attribute_t *attr_list);

sai_status_t l_remove_tam(sai_object_id_t tam_id);

sai_status_t l_set_tam_attribute(sai_object_id_t tam_id,
                                 const sai_attribute_t *attr);

sai_status_t l_get_tam_attribute(sai_object_id_t tam_id, uint32_t attr_count,
                                 sai_attribute_t *attr_list);

sai_status_t l_create_tam_math_func(sai_object_id_t *tam_math_func_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list);

sai_status_t l_remove_tam_math_func(sai_object_id_t tam_math_func_id);

sai_status_t l_set_tam_math_func_attribute(sai_object_id_t tam_math_func_id,
                                           const sai_attribute_t *attr);

sai_status_t l_get_tam_math_func_attribute(sai_object_id_t tam_math_func_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list);

sai_status_t l_create_tam_report(sai_object_id_t *tam_report_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list);

sai_status_t l_remove_tam_report(sai_object_id_t tam_report_id);

sai_status_t l_set_tam_report_attribute(sai_object_id_t tam_report_id,
                                        const sai_attribute_t *attr);

sai_status_t l_get_tam_report_attribute(sai_object_id_t tam_report_id,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list);

sai_status_t l_create_tam_event_threshold(
    sai_object_id_t *tam_event_threshold_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_tam_event_threshold(
    sai_object_id_t tam_event_threshold_id);

sai_status_t l_set_tam_event_threshold_attribute(
    sai_object_id_t tam_event_threshold_id, const sai_attribute_t *attr);

sai_status_t l_get_tam_event_threshold_attribute(
    sai_object_id_t tam_event_threshold_id, uint32_t attr_count,
    sai_attribute_t *attr_list);

sai_status_t l_create_tam_int(sai_object_id_t *tam_int_id,
                              sai_object_id_t switch_id, uint32_t attr_count,
                              const sai_attribute_t *attr_list);

sai_status_t l_remove_tam_int(sai_object_id_t tam_int_id);

sai_status_t l_set_tam_int_attribute(sai_object_id_t tam_int_id,
                                     const sai_attribute_t *attr);

sai_status_t l_get_tam_int_attribute(sai_object_id_t tam_int_id,
                                     uint32_t attr_count,
                                     sai_attribute_t *attr_list);

sai_status_t l_create_tam_tel_type(sai_object_id_t *tam_tel_type_id,
                                   sai_object_id_t switch_id,
                                   uint32_t attr_count,
                                   const sai_attribute_t *attr_list);

sai_status_t l_remove_tam_tel_type(sai_object_id_t tam_tel_type_id);

sai_status_t l_set_tam_tel_type_attribute(sai_object_id_t tam_tel_type_id,
                                          const sai_attribute_t *attr);

sai_status_t l_get_tam_tel_type_attribute(sai_object_id_t tam_tel_type_id,
                                          uint32_t attr_count,
                                          sai_attribute_t *attr_list);

sai_status_t l_create_tam_transport(sai_object_id_t *tam_transport_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list);

sai_status_t l_remove_tam_transport(sai_object_id_t tam_transport_id);

sai_status_t l_set_tam_transport_attribute(sai_object_id_t tam_transport_id,
                                           const sai_attribute_t *attr);

sai_status_t l_get_tam_transport_attribute(sai_object_id_t tam_transport_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list);

sai_status_t l_create_tam_telemetry(sai_object_id_t *tam_telemetry_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list);

sai_status_t l_remove_tam_telemetry(sai_object_id_t tam_telemetry_id);

sai_status_t l_set_tam_telemetry_attribute(sai_object_id_t tam_telemetry_id,
                                           const sai_attribute_t *attr);

sai_status_t l_get_tam_telemetry_attribute(sai_object_id_t tam_telemetry_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list);

sai_status_t l_create_tam_collector(sai_object_id_t *tam_collector_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list);

sai_status_t l_remove_tam_collector(sai_object_id_t tam_collector_id);

sai_status_t l_set_tam_collector_attribute(sai_object_id_t tam_collector_id,
                                           const sai_attribute_t *attr);

sai_status_t l_get_tam_collector_attribute(sai_object_id_t tam_collector_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list);

sai_status_t l_create_tam_event_action(sai_object_id_t *tam_event_action_id,
                                       sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list);

sai_status_t l_remove_tam_event_action(sai_object_id_t tam_event_action_id);

sai_status_t l_set_tam_event_action_attribute(
    sai_object_id_t tam_event_action_id, const sai_attribute_t *attr);

sai_status_t l_get_tam_event_action_attribute(
    sai_object_id_t tam_event_action_id, uint32_t attr_count,
    sai_attribute_t *attr_list);

sai_status_t l_create_tam_event(sai_object_id_t *tam_event_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list);

sai_status_t l_remove_tam_event(sai_object_id_t tam_event_id);

sai_status_t l_set_tam_event_attribute(sai_object_id_t tam_event_id,
                                       const sai_attribute_t *attr);

sai_status_t l_get_tam_event_attribute(sai_object_id_t tam_event_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list);

#endif  // DATAPLANE_STANDALONE_SAI_TAM_H_
