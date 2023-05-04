
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

#include "dataplane/standalone/sai/policer.h"

#include "dataplane/standalone/log/log.h"

const sai_policer_api_t l_policer = {
    .create_policer = l_create_policer,
    .remove_policer = l_remove_policer,
    .set_policer_attribute = l_set_policer_attribute,
    .get_policer_attribute = l_get_policer_attribute,
    .get_policer_stats = l_get_policer_stats,
    .get_policer_stats_ext = l_get_policer_stats_ext,
    .clear_policer_stats = l_clear_policer_stats,
};

sai_status_t l_create_policer(sai_object_id_t *policer_id,
                              sai_object_id_t switch_id, uint32_t attr_count,
                              const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_policer(sai_object_id_t policer_id) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_policer_attribute(sai_object_id_t policer_id,
                                     const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_policer_attribute(sai_object_id_t policer_id,
                                     uint32_t attr_count,
                                     sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_policer_stats(sai_object_id_t policer_id,
                                 uint32_t number_of_counters,
                                 const sai_stat_id_t *counter_ids,
                                 uint64_t *counters) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_policer_stats_ext(sai_object_id_t policer_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     sai_stats_mode_t mode,
                                     uint64_t *counters) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_clear_policer_stats(sai_object_id_t policer_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}
