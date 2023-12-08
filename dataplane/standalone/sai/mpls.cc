

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

#include "dataplane/standalone/sai/mpls.h"

#include <glog/logging.h>

#include "dataplane/proto/common.pb.h"
#include "dataplane/proto/mpls.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_mpls_api_t l_mpls = {
    .create_inseg_entry = l_create_inseg_entry,
    .remove_inseg_entry = l_remove_inseg_entry,
    .set_inseg_entry_attribute = l_set_inseg_entry_attribute,
    .get_inseg_entry_attribute = l_get_inseg_entry_attribute,
    .create_inseg_entries = l_create_inseg_entries,
    .remove_inseg_entries = l_remove_inseg_entries,
    .set_inseg_entries_attribute = l_set_inseg_entries_attribute,
    .get_inseg_entries_attribute = l_get_inseg_entries_attribute,
};

lemming::dataplane::sai::CreateInsegEntryRequest convert_create_inseg_entry(
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateInsegEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_INSEG_ENTRY_ATTR_NUM_OF_POP:
        msg.set_num_of_pop(attr_list[i].value.u8);
        break;
      case SAI_INSEG_ENTRY_ATTR_PACKET_ACTION:
        msg.set_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_INSEG_ENTRY_ATTR_TRAP_PRIORITY:
        msg.set_trap_priority(attr_list[i].value.u8);
        break;
      case SAI_INSEG_ENTRY_ATTR_NEXT_HOP_ID:
        msg.set_next_hop_id(attr_list[i].value.oid);
        break;
      case SAI_INSEG_ENTRY_ATTR_PSC_TYPE:
        msg.set_psc_type(
            static_cast<lemming::dataplane::sai::InsegEntryPscType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_INSEG_ENTRY_ATTR_QOS_TC:
        msg.set_qos_tc(attr_list[i].value.u8);
        break;
      case SAI_INSEG_ENTRY_ATTR_MPLS_EXP_TO_TC_MAP:
        msg.set_mpls_exp_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_INSEG_ENTRY_ATTR_MPLS_EXP_TO_COLOR_MAP:
        msg.set_mpls_exp_to_color_map(attr_list[i].value.oid);
        break;
      case SAI_INSEG_ENTRY_ATTR_POP_TTL_MODE:
        msg.set_pop_ttl_mode(
            static_cast<lemming::dataplane::sai::InsegEntryPopTtlMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_INSEG_ENTRY_ATTR_POP_QOS_MODE:
        msg.set_pop_qos_mode(
            static_cast<lemming::dataplane::sai::InsegEntryPopQosMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_INSEG_ENTRY_ATTR_COUNTER_ID:
        msg.set_counter_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_inseg_entry(const sai_inseg_entry_t *inseg_entry,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateInsegEntryRequest req =
      convert_create_inseg_entry(attr_count, attr_list);
  lemming::dataplane::sai::CreateInsegEntryResponse resp;
  grpc::ClientContext context;

  grpc::Status status = mpls->CreateInsegEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_inseg_entry(const sai_inseg_entry_t *inseg_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveInsegEntryRequest req;
  lemming::dataplane::sai::RemoveInsegEntryResponse resp;
  grpc::ClientContext context;

  grpc::Status status = mpls->RemoveInsegEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_inseg_entry_attribute(const sai_inseg_entry_t *inseg_entry,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetInsegEntryAttributeRequest req;
  lemming::dataplane::sai::SetInsegEntryAttributeResponse resp;
  grpc::ClientContext context;

  switch (attr->id) {
    case SAI_INSEG_ENTRY_ATTR_NUM_OF_POP:
      req.set_num_of_pop(attr->value.u8);
      break;
    case SAI_INSEG_ENTRY_ATTR_PACKET_ACTION:
      req.set_packet_action(static_cast<lemming::dataplane::sai::PacketAction>(
          attr->value.s32 + 1));
      break;
    case SAI_INSEG_ENTRY_ATTR_TRAP_PRIORITY:
      req.set_trap_priority(attr->value.u8);
      break;
    case SAI_INSEG_ENTRY_ATTR_NEXT_HOP_ID:
      req.set_next_hop_id(attr->value.oid);
      break;
    case SAI_INSEG_ENTRY_ATTR_PSC_TYPE:
      req.set_psc_type(static_cast<lemming::dataplane::sai::InsegEntryPscType>(
          attr->value.s32 + 1));
      break;
    case SAI_INSEG_ENTRY_ATTR_QOS_TC:
      req.set_qos_tc(attr->value.u8);
      break;
    case SAI_INSEG_ENTRY_ATTR_MPLS_EXP_TO_TC_MAP:
      req.set_mpls_exp_to_tc_map(attr->value.oid);
      break;
    case SAI_INSEG_ENTRY_ATTR_MPLS_EXP_TO_COLOR_MAP:
      req.set_mpls_exp_to_color_map(attr->value.oid);
      break;
    case SAI_INSEG_ENTRY_ATTR_POP_TTL_MODE:
      req.set_pop_ttl_mode(
          static_cast<lemming::dataplane::sai::InsegEntryPopTtlMode>(
              attr->value.s32 + 1));
      break;
    case SAI_INSEG_ENTRY_ATTR_POP_QOS_MODE:
      req.set_pop_qos_mode(
          static_cast<lemming::dataplane::sai::InsegEntryPopQosMode>(
              attr->value.s32 + 1));
      break;
    case SAI_INSEG_ENTRY_ATTR_COUNTER_ID:
      req.set_counter_id(attr->value.oid);
      break;
  }

  grpc::Status status = mpls->SetInsegEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_inseg_entry_attribute(const sai_inseg_entry_t *inseg_entry,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetInsegEntryAttributeRequest req;
  lemming::dataplane::sai::GetInsegEntryAttributeResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::InsegEntryAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = mpls->GetInsegEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_INSEG_ENTRY_ATTR_NUM_OF_POP:
        attr_list[i].value.u8 = resp.attr().num_of_pop();
        break;
      case SAI_INSEG_ENTRY_ATTR_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().packet_action() - 1);
        break;
      case SAI_INSEG_ENTRY_ATTR_TRAP_PRIORITY:
        attr_list[i].value.u8 = resp.attr().trap_priority();
        break;
      case SAI_INSEG_ENTRY_ATTR_NEXT_HOP_ID:
        attr_list[i].value.oid = resp.attr().next_hop_id();
        break;
      case SAI_INSEG_ENTRY_ATTR_PSC_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().psc_type() - 1);
        break;
      case SAI_INSEG_ENTRY_ATTR_QOS_TC:
        attr_list[i].value.u8 = resp.attr().qos_tc();
        break;
      case SAI_INSEG_ENTRY_ATTR_MPLS_EXP_TO_TC_MAP:
        attr_list[i].value.oid = resp.attr().mpls_exp_to_tc_map();
        break;
      case SAI_INSEG_ENTRY_ATTR_MPLS_EXP_TO_COLOR_MAP:
        attr_list[i].value.oid = resp.attr().mpls_exp_to_color_map();
        break;
      case SAI_INSEG_ENTRY_ATTR_POP_TTL_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().pop_ttl_mode() - 1);
        break;
      case SAI_INSEG_ENTRY_ATTR_POP_QOS_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().pop_qos_mode() - 1);
        break;
      case SAI_INSEG_ENTRY_ATTR_COUNTER_ID:
        attr_list[i].value.oid = resp.attr().counter_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_inseg_entries(uint32_t object_count,
                                    const sai_inseg_entry_t *inseg_entry,
                                    const uint32_t *attr_count,
                                    const sai_attribute_t **attr_list,
                                    sai_bulk_op_error_mode_t mode,
                                    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateInsegEntriesRequest req;
  lemming::dataplane::sai::CreateInsegEntriesResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    auto r = convert_create_inseg_entry(attr_count[i], attr_list[i]);
    *req.add_reqs() = r;
  }

  grpc::Status status = mpls->CreateInsegEntries(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < object_count; i++) {
    object_statuses[i] = SAI_STATUS_SUCCESS;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_inseg_entries(uint32_t object_count,
                                    const sai_inseg_entry_t *inseg_entry,
                                    sai_bulk_op_error_mode_t mode,
                                    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_inseg_entries_attribute(uint32_t object_count,
                                           const sai_inseg_entry_t *inseg_entry,
                                           const sai_attribute_t *attr_list,
                                           sai_bulk_op_error_mode_t mode,
                                           sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_inseg_entries_attribute(uint32_t object_count,
                                           const sai_inseg_entry_t *inseg_entry,
                                           const uint32_t *attr_count,
                                           sai_attribute_t **attr_list,
                                           sai_bulk_op_error_mode_t mode,
                                           sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}
