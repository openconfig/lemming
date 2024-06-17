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

#include "dataplane/proto/sai/buffer.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/standalone/sai/common.h"

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

lemming::dataplane::sai::CreateBufferPoolRequest convert_create_buffer_pool(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateBufferPoolRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BUFFER_POOL_ATTR_TYPE:
        msg.set_type(static_cast<lemming::dataplane::sai::BufferPoolType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_BUFFER_POOL_ATTR_SIZE:
        msg.set_size(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_POOL_ATTR_THRESHOLD_MODE:
        msg.set_threshold_mode(
            static_cast<lemming::dataplane::sai::BufferPoolThresholdMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BUFFER_POOL_ATTR_TAM:
        msg.mutable_tam()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_BUFFER_POOL_ATTR_XOFF_SIZE:
        msg.set_xoff_size(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_POOL_ATTR_WRED_PROFILE_ID:
        msg.set_wred_profile_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateIngressPriorityGroupRequest
convert_create_ingress_priority_group(sai_object_id_t switch_id,
                                      uint32_t attr_count,
                                      const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateIngressPriorityGroupRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_BUFFER_PROFILE:
        msg.set_buffer_profile(attr_list[i].value.oid);
        break;
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_PORT:
        msg.set_port(attr_list[i].value.oid);
        break;
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_TAM:
        msg.mutable_tam()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_INDEX:
        msg.set_index(attr_list[i].value.u8);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateBufferProfileRequest
convert_create_buffer_profile(sai_object_id_t switch_id, uint32_t attr_count,
                              const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateBufferProfileRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BUFFER_PROFILE_ATTR_POOL_ID:
        msg.set_pool_id(attr_list[i].value.oid);
        break;
      case SAI_BUFFER_PROFILE_ATTR_RESERVED_BUFFER_SIZE:
        msg.set_reserved_buffer_size(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_PROFILE_ATTR_THRESHOLD_MODE:
        msg.set_threshold_mode(
            static_cast<lemming::dataplane::sai::BufferProfileThresholdMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BUFFER_PROFILE_ATTR_SHARED_DYNAMIC_TH:
        msg.set_shared_dynamic_th(attr_list[i].value.s8);
        break;
      case SAI_BUFFER_PROFILE_ATTR_SHARED_STATIC_TH:
        msg.set_shared_static_th(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_PROFILE_ATTR_XOFF_TH:
        msg.set_xoff_th(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_PROFILE_ATTR_XON_TH:
        msg.set_xon_th(attr_list[i].value.u64);
        break;
      case SAI_BUFFER_PROFILE_ATTR_XON_OFFSET_TH:
        msg.set_xon_offset_th(attr_list[i].value.u64);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_buffer_pool(sai_object_id_t *buffer_pool_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateBufferPoolRequest req =
      convert_create_buffer_pool(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateBufferPoolResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = buffer->CreateBufferPool(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *buffer_pool_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_buffer_pool(sai_object_id_t buffer_pool_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveBufferPoolRequest req;
  lemming::dataplane::sai::RemoveBufferPoolResponse resp;
  grpc::ClientContext context;
  req.set_oid(buffer_pool_id);

  grpc::Status status = buffer->RemoveBufferPool(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_buffer_pool_attribute(sai_object_id_t buffer_pool_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetBufferPoolAttributeRequest req;
  lemming::dataplane::sai::SetBufferPoolAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(buffer_pool_id);

  switch (attr->id) {
    case SAI_BUFFER_POOL_ATTR_SIZE:
      req.set_size(attr->value.u64);
      break;
    case SAI_BUFFER_POOL_ATTR_TAM:
      req.mutable_tam()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_BUFFER_POOL_ATTR_XOFF_SIZE:
      req.set_xoff_size(attr->value.u64);
      break;
    case SAI_BUFFER_POOL_ATTR_WRED_PROFILE_ID:
      req.set_wred_profile_id(attr->value.oid);
      break;
  }

  grpc::Status status = buffer->SetBufferPoolAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_buffer_pool_attribute(sai_object_id_t buffer_pool_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetBufferPoolAttributeRequest req;
  lemming::dataplane::sai::GetBufferPoolAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(buffer_pool_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::BufferPoolAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = buffer->GetBufferPoolAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BUFFER_POOL_ATTR_SHARED_SIZE:
        attr_list[i].value.u64 = resp.attr().shared_size();
        break;
      case SAI_BUFFER_POOL_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_BUFFER_POOL_ATTR_SIZE:
        attr_list[i].value.u64 = resp.attr().size();
        break;
      case SAI_BUFFER_POOL_ATTR_THRESHOLD_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().threshold_mode() - 1);
        break;
      case SAI_BUFFER_POOL_ATTR_TAM:
        copy_list(attr_list[i].value.objlist.list, resp.attr().tam(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_BUFFER_POOL_ATTR_XOFF_SIZE:
        attr_list[i].value.u64 = resp.attr().xoff_size();
        break;
      case SAI_BUFFER_POOL_ATTR_WRED_PROFILE_ID:
        attr_list[i].value.oid = resp.attr().wred_profile_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_buffer_pool_stats(sai_object_id_t buffer_pool_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetBufferPoolStatsRequest req;
  lemming::dataplane::sai::GetBufferPoolStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(buffer_pool_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(static_cast<lemming::dataplane::sai::BufferPoolStat>(
        counter_ids[i] + 1));
  }
  grpc::Status status = buffer->GetBufferPoolStats(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_buffer_pool_stats_ext(sai_object_id_t buffer_pool_id,
                                         uint32_t number_of_counters,
                                         const sai_stat_id_t *counter_ids,
                                         sai_stats_mode_t mode,
                                         uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_buffer_pool_stats(sai_object_id_t buffer_pool_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_ingress_priority_group(
    sai_object_id_t *ingress_priority_group_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIngressPriorityGroupRequest req =
      convert_create_ingress_priority_group(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateIngressPriorityGroupResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status =
      buffer->CreateIngressPriorityGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *ingress_priority_group_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_ingress_priority_group(
    sai_object_id_t ingress_priority_group_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveIngressPriorityGroupRequest req;
  lemming::dataplane::sai::RemoveIngressPriorityGroupResponse resp;
  grpc::ClientContext context;
  req.set_oid(ingress_priority_group_id);

  grpc::Status status =
      buffer->RemoveIngressPriorityGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_ingress_priority_group_attribute(
    sai_object_id_t ingress_priority_group_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetIngressPriorityGroupAttributeRequest req;
  lemming::dataplane::sai::SetIngressPriorityGroupAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(ingress_priority_group_id);

  switch (attr->id) {
    case SAI_INGRESS_PRIORITY_GROUP_ATTR_BUFFER_PROFILE:
      req.set_buffer_profile(attr->value.oid);
      break;
    case SAI_INGRESS_PRIORITY_GROUP_ATTR_TAM:
      req.mutable_tam()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
  }

  grpc::Status status =
      buffer->SetIngressPriorityGroupAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_ingress_priority_group_attribute(
    sai_object_id_t ingress_priority_group_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIngressPriorityGroupAttributeRequest req;
  lemming::dataplane::sai::GetIngressPriorityGroupAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(ingress_priority_group_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::IngressPriorityGroupAttr>(
            attr_list[i].id + 1));
  }
  grpc::Status status =
      buffer->GetIngressPriorityGroupAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_BUFFER_PROFILE:
        attr_list[i].value.oid = resp.attr().buffer_profile();
        break;
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_PORT:
        attr_list[i].value.oid = resp.attr().port();
        break;
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_TAM:
        copy_list(attr_list[i].value.objlist.list, resp.attr().tam(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_INGRESS_PRIORITY_GROUP_ATTR_INDEX:
        attr_list[i].value.u8 = resp.attr().index();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_ingress_priority_group_stats(
    sai_object_id_t ingress_priority_group_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIngressPriorityGroupStatsRequest req;
  lemming::dataplane::sai::GetIngressPriorityGroupStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(ingress_priority_group_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(
        static_cast<lemming::dataplane::sai::IngressPriorityGroupStat>(
            counter_ids[i] + 1));
  }
  grpc::Status status =
      buffer->GetIngressPriorityGroupStats(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_ingress_priority_group_stats_ext(
    sai_object_id_t ingress_priority_group_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids, sai_stats_mode_t mode,
    uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_ingress_priority_group_stats(
    sai_object_id_t ingress_priority_group_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_buffer_profile(sai_object_id_t *buffer_profile_id,
                                     sai_object_id_t switch_id,
                                     uint32_t attr_count,
                                     const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateBufferProfileRequest req =
      convert_create_buffer_profile(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateBufferProfileResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = buffer->CreateBufferProfile(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *buffer_profile_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_buffer_profile(sai_object_id_t buffer_profile_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveBufferProfileRequest req;
  lemming::dataplane::sai::RemoveBufferProfileResponse resp;
  grpc::ClientContext context;
  req.set_oid(buffer_profile_id);

  grpc::Status status = buffer->RemoveBufferProfile(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_buffer_profile_attribute(sai_object_id_t buffer_profile_id,
                                            const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetBufferProfileAttributeRequest req;
  lemming::dataplane::sai::SetBufferProfileAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(buffer_profile_id);

  switch (attr->id) {
    case SAI_BUFFER_PROFILE_ATTR_RESERVED_BUFFER_SIZE:
      req.set_reserved_buffer_size(attr->value.u64);
      break;
    case SAI_BUFFER_PROFILE_ATTR_SHARED_DYNAMIC_TH:
      req.set_shared_dynamic_th(attr->value.s8);
      break;
    case SAI_BUFFER_PROFILE_ATTR_SHARED_STATIC_TH:
      req.set_shared_static_th(attr->value.u64);
      break;
    case SAI_BUFFER_PROFILE_ATTR_XOFF_TH:
      req.set_xoff_th(attr->value.u64);
      break;
    case SAI_BUFFER_PROFILE_ATTR_XON_TH:
      req.set_xon_th(attr->value.u64);
      break;
    case SAI_BUFFER_PROFILE_ATTR_XON_OFFSET_TH:
      req.set_xon_offset_th(attr->value.u64);
      break;
  }

  grpc::Status status = buffer->SetBufferProfileAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_buffer_profile_attribute(sai_object_id_t buffer_profile_id,
                                            uint32_t attr_count,
                                            sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetBufferProfileAttributeRequest req;
  lemming::dataplane::sai::GetBufferProfileAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(buffer_profile_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::BufferProfileAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = buffer->GetBufferProfileAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BUFFER_PROFILE_ATTR_POOL_ID:
        attr_list[i].value.oid = resp.attr().pool_id();
        break;
      case SAI_BUFFER_PROFILE_ATTR_RESERVED_BUFFER_SIZE:
        attr_list[i].value.u64 = resp.attr().reserved_buffer_size();
        break;
      case SAI_BUFFER_PROFILE_ATTR_THRESHOLD_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().threshold_mode() - 1);
        break;
      case SAI_BUFFER_PROFILE_ATTR_SHARED_DYNAMIC_TH:
        attr_list[i].value.s8 = resp.attr().shared_dynamic_th();
        break;
      case SAI_BUFFER_PROFILE_ATTR_SHARED_STATIC_TH:
        attr_list[i].value.u64 = resp.attr().shared_static_th();
        break;
      case SAI_BUFFER_PROFILE_ATTR_XOFF_TH:
        attr_list[i].value.u64 = resp.attr().xoff_th();
        break;
      case SAI_BUFFER_PROFILE_ATTR_XON_TH:
        attr_list[i].value.u64 = resp.attr().xon_th();
        break;
      case SAI_BUFFER_PROFILE_ATTR_XON_OFFSET_TH:
        attr_list[i].value.u64 = resp.attr().xon_offset_th();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
