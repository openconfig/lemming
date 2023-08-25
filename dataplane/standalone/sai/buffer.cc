

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

#include "dataplane/standalone/sai/buffer.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/buffer.pb.h"
#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_buffer_api_t l_buffer = {
    .create_buffer_pool = l_create_buffer_pool,
    .remove_buffer_pool = l_remove_buffer_pool,
    .set_buffer_pool_attribute = l_set_buffer_pool_attribute,
    .get_buffer_pool_attribute = l_get_buffer_pool_attribute,
    .get_buffer_pool_stats = l_get_buffer_pool_stats,
    .get_buffer_pool_stats_ext = l_get_buffer_pool_stats_ext,
    .clear_buffer_pool_stats = l_clear_buffer_pool_stats,
    .create_ingress_priority_group = l_create_ingress_priority_group,
    .remove_ingress_priority_group = l_remove_ingress_priority_group,
    .set_ingress_priority_group_attribute =
        l_set_ingress_priority_group_attribute,
    .get_ingress_priority_group_attribute =
        l_get_ingress_priority_group_attribute,
    .get_ingress_priority_group_stats = l_get_ingress_priority_group_stats,
    .get_ingress_priority_group_stats_ext =
        l_get_ingress_priority_group_stats_ext,
    .clear_ingress_priority_group_stats = l_clear_ingress_priority_group_stats,
    .create_buffer_profile = l_create_buffer_profile,
    .remove_buffer_profile = l_remove_buffer_profile,
    .set_buffer_profile_attribute = l_set_buffer_profile_attribute,
    .get_buffer_profile_attribute = l_get_buffer_profile_attribute,
};

sai_status_t l_create_buffer_pool(sai_object_id_t *buffer_pool_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateBufferPoolRequest req;
  lemming::dataplane::sai::CreateBufferPoolResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BUFFER_POOL_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::BufferPoolType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_BUFFER_POOL_ATTR_SIZE:
        req.set_size(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_POOL_ATTR_THRESHOLD_MODE:
        req.set_threshold_mode(
            static_cast<lemming::dataplane::sai::BufferPoolThresholdMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BUFFER_POOL_ATTR_TAM:
        req.mutable_tam()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_BUFFER_POOL_ATTR_XOFF_SIZE:
        req.set_xoff_size(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_POOL_ATTR_WRED_PROFILE_ID:
        req.set_wred_profile_id(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = buffer->CreateBufferPool(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *buffer_pool_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_buffer_pool(sai_object_id_t buffer_pool_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id);
}

sai_status_t l_set_buffer_pool_attribute(sai_object_id_t buffer_pool_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id,
                                   attr);
}

sai_status_t l_get_buffer_pool_attribute(sai_object_id_t buffer_pool_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_buffer_pool_stats(sai_object_id_t buffer_pool_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_buffer_pool_stats_ext(sai_object_id_t buffer_pool_id,
                                         uint32_t number_of_counters,
                                         const sai_stat_id_t *counter_ids,
                                         sai_stats_mode_t mode,
                                         uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_buffer_pool_stats(sai_object_id_t buffer_pool_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_BUFFER_POOL, buffer_pool_id,
                                 number_of_counters, counter_ids);
}

sai_status_t l_create_ingress_priority_group(
    sai_object_id_t *ingress_priority_group_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIngressPriorityGroupRequest req;
  lemming::dataplane::sai::CreateIngressPriorityGroupResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_BUFFER_PROFILE:
        req.set_buffer_profile(attr_list[i].value.oid);
        break;
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_PORT:
        req.set_port(attr_list[i].value.oid);
        break;
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_TAM:
        req.mutable_tam()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_INDEX:
        req.set_index(attr_list[i].value.u8);
        break;
    }
  }
  grpc::Status status =
      buffer->CreateIngressPriorityGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *ingress_priority_group_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP,
                            ingress_priority_group_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_ingress_priority_group(
    sai_object_id_t ingress_priority_group_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP,
                            ingress_priority_group_id);
}

sai_status_t l_set_ingress_priority_group_attribute(
    sai_object_id_t ingress_priority_group_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP,
                                   ingress_priority_group_id, attr);
}

sai_status_t l_get_ingress_priority_group_attribute(
    sai_object_id_t ingress_priority_group_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP,
                                   ingress_priority_group_id, attr_count,
                                   attr_list);
}

sai_status_t l_get_ingress_priority_group_stats(
    sai_object_id_t ingress_priority_group_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP,
                               ingress_priority_group_id, number_of_counters,
                               counter_ids, counters);
}

sai_status_t l_get_ingress_priority_group_stats_ext(
    sai_object_id_t ingress_priority_group_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids, sai_stats_mode_t mode,
    uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(
      SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP, ingress_priority_group_id,
      number_of_counters, counter_ids, mode, counters);
}

sai_status_t l_clear_ingress_priority_group_stats(
    sai_object_id_t ingress_priority_group_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP,
                                 ingress_priority_group_id, number_of_counters,
                                 counter_ids);
}

sai_status_t l_create_buffer_profile(sai_object_id_t *buffer_profile_id,
                                     sai_object_id_t switch_id,
                                     uint32_t attr_count,
                                     const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateBufferProfileRequest req;
  lemming::dataplane::sai::CreateBufferProfileResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BUFFER_PROFILE_ATTR_POOL_ID:
        req.set_pool_id(attr_list[i].value.oid);
        break;
      case SAI_BUFFER_PROFILE_ATTR_RESERVED_BUFFER_SIZE:
        req.set_reserved_buffer_size(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_PROFILE_ATTR_THRESHOLD_MODE:
        req.set_threshold_mode(
            static_cast<lemming::dataplane::sai::BufferProfileThresholdMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BUFFER_PROFILE_ATTR_SHARED_DYNAMIC_TH:
        req.set_shared_dynamic_th(attr_list[i].value.s8);
        break;
      case SAI_BUFFER_PROFILE_ATTR_SHARED_STATIC_TH:
        req.set_shared_static_th(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_PROFILE_ATTR_XOFF_TH:
        req.set_xoff_th(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_PROFILE_ATTR_XON_TH:
        req.set_xon_th(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_PROFILE_ATTR_XON_OFFSET_TH:
        req.set_xon_offset_th(attr_list[i].value.u64);
        break;
    }
  }
  grpc::Status status = buffer->CreateBufferProfile(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *buffer_profile_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_BUFFER_PROFILE, buffer_profile_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_buffer_profile(sai_object_id_t buffer_profile_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_BUFFER_PROFILE, buffer_profile_id);
}

sai_status_t l_set_buffer_profile_attribute(sai_object_id_t buffer_profile_id,
                                            const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_BUFFER_PROFILE,
                                   buffer_profile_id, attr);
}

sai_status_t l_get_buffer_profile_attribute(sai_object_id_t buffer_profile_id,
                                            uint32_t attr_count,
                                            sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_BUFFER_PROFILE,
                                   buffer_profile_id, attr_count, attr_list);
}
