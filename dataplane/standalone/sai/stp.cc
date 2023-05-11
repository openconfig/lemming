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

#include "dataplane/standalone/sai/stp.h"

#include <glog/logging.h>

#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_stp_api_t l_stp = {
    .create_stp = l_create_stp,
    .remove_stp = l_remove_stp,
    .set_stp_attribute = l_set_stp_attribute,
    .get_stp_attribute = l_get_stp_attribute,
    .create_stp_port = l_create_stp_port,
    .remove_stp_port = l_remove_stp_port,
    .set_stp_port_attribute = l_set_stp_port_attribute,
    .get_stp_port_attribute = l_get_stp_port_attribute,
    .create_stp_ports = l_create_stp_ports,
    .remove_stp_ports = l_remove_stp_ports,
};

sai_status_t l_create_stp(sai_object_id_t *stp_id, sai_object_id_t switch_id,
                          uint32_t attr_count,
                          const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_STP, stp_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_stp(sai_object_id_t stp_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_STP, stp_id);
}

sai_status_t l_set_stp_attribute(sai_object_id_t stp_id,
                                 const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_STP, stp_id, attr);
}

sai_status_t l_get_stp_attribute(sai_object_id_t stp_id, uint32_t attr_count,
                                 sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_STP, stp_id, attr_count,
                                   attr_list);
}

sai_status_t l_create_stp_port(sai_object_id_t *stp_port_id,
                               sai_object_id_t switch_id, uint32_t attr_count,
                               const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_STP_PORT, stp_port_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_stp_port(sai_object_id_t stp_port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_STP_PORT, stp_port_id);
}

sai_status_t l_set_stp_port_attribute(sai_object_id_t stp_port_id,
                                      const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_STP_PORT, stp_port_id, attr);
}

sai_status_t l_get_stp_port_attribute(sai_object_id_t stp_port_id,
                                      uint32_t attr_count,
                                      sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_STP_PORT, stp_port_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_stp_ports(sai_object_id_t switch_id,
                                uint32_t object_count,
                                const uint32_t *attr_count,
                                const sai_attribute_t **attr_list,
                                sai_bulk_op_error_mode_t mode,
                                sai_object_id_t *object_id,
                                sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create_bulk(SAI_OBJECT_TYPE_STP_PORT, switch_id,
                                 object_count, attr_count, attr_list, mode,
                                 object_id, object_statuses);
}

sai_status_t l_remove_stp_ports(uint32_t object_count,
                                const sai_object_id_t *object_id,
                                sai_bulk_op_error_mode_t mode,
                                sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove_bulk(SAI_OBJECT_TYPE_STP_PORT, object_count,
                                 object_id, mode, object_statuses);
}
