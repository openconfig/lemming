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

#include "dataplane/standalone/sai/port.h"

#include <glog/logging.h>

#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_port_api_t l_port = {
    .create_port = l_create_port,
    .remove_port = l_remove_port,
    .set_port_attribute = l_set_port_attribute,
    .get_port_attribute = l_get_port_attribute,
    .get_port_stats = l_get_port_stats,
    .get_port_stats_ext = l_get_port_stats_ext,
    .clear_port_stats = l_clear_port_stats,
    .clear_port_all_stats = l_clear_port_all_stats,
    .create_port_pool = l_create_port_pool,
    .remove_port_pool = l_remove_port_pool,
    .set_port_pool_attribute = l_set_port_pool_attribute,
    .get_port_pool_attribute = l_get_port_pool_attribute,
    .get_port_pool_stats = l_get_port_pool_stats,
    .get_port_pool_stats_ext = l_get_port_pool_stats_ext,
    .clear_port_pool_stats = l_clear_port_pool_stats,
    .create_port_connector = l_create_port_connector,
    .remove_port_connector = l_remove_port_connector,
    .set_port_connector_attribute = l_set_port_connector_attribute,
    .get_port_connector_attribute = l_get_port_connector_attribute,
    .create_port_serdes = l_create_port_serdes,
    .remove_port_serdes = l_remove_port_serdes,
    .set_port_serdes_attribute = l_set_port_serdes_attribute,
    .get_port_serdes_attribute = l_get_port_serdes_attribute,
};

sai_status_t l_create_port(sai_object_id_t *port_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_PORT, port_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_port(sai_object_id_t port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_PORT, port_id);
}

sai_status_t l_set_port_attribute(sai_object_id_t port_id,
                                  const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_PORT, port_id, attr);
}

sai_status_t l_get_port_attribute(sai_object_id_t port_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_PORT, port_id, attr_count,
                                   attr_list);
}

sai_status_t l_get_port_stats(sai_object_id_t port_id,
                              uint32_t number_of_counters,
                              const sai_stat_id_t *counter_ids,
                              uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_stats(SAI_OBJECT_TYPE_PORT, port_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_port_stats_ext(sai_object_id_t port_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids,
                                  sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_stats_ext(SAI_OBJECT_TYPE_PORT, port_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_port_stats(sai_object_id_t port_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->clear_stats(SAI_OBJECT_TYPE_PORT, port_id,
                                 number_of_counters, counter_ids);
}

sai_status_t l_clear_port_all_stats(sai_object_id_t port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_port_pool(sai_object_id_t *port_pool_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_port_pool(sai_object_id_t port_pool_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id);
}

sai_status_t l_set_port_pool_attribute(sai_object_id_t port_pool_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id,
                                   attr);
}

sai_status_t l_get_port_pool_attribute(sai_object_id_t port_pool_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_port_pool_stats(sai_object_id_t port_pool_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids,
                                   uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_stats(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_port_pool_stats_ext(sai_object_id_t port_pool_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids,
                                       sai_stats_mode_t mode,
                                       uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_stats_ext(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_port_pool_stats(sai_object_id_t port_pool_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->clear_stats(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id,
                                 number_of_counters, counter_ids);
}

sai_status_t l_create_port_connector(sai_object_id_t *port_connector_id,
                                     sai_object_id_t switch_id,
                                     uint32_t attr_count,
                                     const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_PORT_CONNECTOR, port_connector_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_port_connector(sai_object_id_t port_connector_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_PORT_CONNECTOR, port_connector_id);
}

sai_status_t l_set_port_connector_attribute(sai_object_id_t port_connector_id,
                                            const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_PORT_CONNECTOR,
                                   port_connector_id, attr);
}

sai_status_t l_get_port_connector_attribute(sai_object_id_t port_connector_id,
                                            uint32_t attr_count,
                                            sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_PORT_CONNECTOR,
                                   port_connector_id, attr_count, attr_list);
}

sai_status_t l_create_port_serdes(sai_object_id_t *port_serdes_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_PORT_SERDES, port_serdes_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_port_serdes(sai_object_id_t port_serdes_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_PORT_SERDES, port_serdes_id);
}

sai_status_t l_set_port_serdes_attribute(sai_object_id_t port_serdes_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_PORT_SERDES, port_serdes_id,
                                   attr);
}

sai_status_t l_get_port_serdes_attribute(sai_object_id_t port_serdes_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_PORT_SERDES, port_serdes_id,
                                   attr_count, attr_list);
}
