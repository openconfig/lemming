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

#include "dataplane/standalone/sai/switch.h"

#include <glog/logging.h>

#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_switch_api_t l_switch = {
    .create_switch = l_create_switch,
    .remove_switch = l_remove_switch,
    .set_switch_attribute = l_set_switch_attribute,
    .get_switch_attribute = l_get_switch_attribute,
    .get_switch_stats = l_get_switch_stats,
    .get_switch_stats_ext = l_get_switch_stats_ext,
    .clear_switch_stats = l_clear_switch_stats,
    .switch_mdio_read = l_switch_mdio_read,
    .switch_mdio_write = l_switch_mdio_write,
    .create_switch_tunnel = l_create_switch_tunnel,
    .remove_switch_tunnel = l_remove_switch_tunnel,
    .set_switch_tunnel_attribute = l_set_switch_tunnel_attribute,
    .get_switch_tunnel_attribute = l_get_switch_tunnel_attribute,
};

sai_status_t l_create_switch(sai_object_id_t *switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_SWITCH, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_switch(sai_object_id_t switch_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_SWITCH, switch_id);
}

sai_status_t l_set_switch_attribute(sai_object_id_t switch_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_SWITCH, switch_id, attr);
}

sai_status_t l_get_switch_attribute(sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_SWITCH, switch_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_switch_stats(sai_object_id_t switch_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids,
                                uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_stats(SAI_OBJECT_TYPE_SWITCH, switch_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_switch_stats_ext(sai_object_id_t switch_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t *counter_ids,
                                    sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_stats_ext(SAI_OBJECT_TYPE_SWITCH, switch_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_switch_stats(sai_object_id_t switch_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->clear_stats(SAI_OBJECT_TYPE_SWITCH, switch_id,
                                 number_of_counters, counter_ids);
}

sai_status_t l_switch_mdio_read(sai_object_id_t switch_id, uint32_t device_addr,
                                uint32_t start_reg_addr,
                                uint32_t number_of_registers,
                                uint32_t *reg_val) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_switch_mdio_write(sai_object_id_t switch_id,
                                 uint32_t device_addr, uint32_t start_reg_addr,
                                 uint32_t number_of_registers,
                                 const uint32_t *reg_val) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_switch_tunnel(sai_object_id_t *switch_tunnel_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_SWITCH_TUNNEL, switch_tunnel_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_switch_tunnel(sai_object_id_t switch_tunnel_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_SWITCH_TUNNEL, switch_tunnel_id);
}

sai_status_t l_set_switch_tunnel_attribute(sai_object_id_t switch_tunnel_id,
                                           const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_SWITCH_TUNNEL,
                                   switch_tunnel_id, attr);
}

sai_status_t l_get_switch_tunnel_attribute(sai_object_id_t switch_tunnel_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_SWITCH_TUNNEL,
                                   switch_tunnel_id, attr_count, attr_list);
}
