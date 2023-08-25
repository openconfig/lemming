

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

#include "dataplane/standalone/sai/hostif.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/hostif.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_hostif_api_t l_hostif = {
    .create_hostif = l_create_hostif,
    .remove_hostif = l_remove_hostif,
    .set_hostif_attribute = l_set_hostif_attribute,
    .get_hostif_attribute = l_get_hostif_attribute,
    .create_hostif_table_entry = l_create_hostif_table_entry,
    .remove_hostif_table_entry = l_remove_hostif_table_entry,
    .set_hostif_table_entry_attribute = l_set_hostif_table_entry_attribute,
    .get_hostif_table_entry_attribute = l_get_hostif_table_entry_attribute,
    .create_hostif_trap_group = l_create_hostif_trap_group,
    .remove_hostif_trap_group = l_remove_hostif_trap_group,
    .set_hostif_trap_group_attribute = l_set_hostif_trap_group_attribute,
    .get_hostif_trap_group_attribute = l_get_hostif_trap_group_attribute,
    .create_hostif_trap = l_create_hostif_trap,
    .remove_hostif_trap = l_remove_hostif_trap,
    .set_hostif_trap_attribute = l_set_hostif_trap_attribute,
    .get_hostif_trap_attribute = l_get_hostif_trap_attribute,
    .create_hostif_user_defined_trap = l_create_hostif_user_defined_trap,
    .remove_hostif_user_defined_trap = l_remove_hostif_user_defined_trap,
    .set_hostif_user_defined_trap_attribute =
        l_set_hostif_user_defined_trap_attribute,
    .get_hostif_user_defined_trap_attribute =
        l_get_hostif_user_defined_trap_attribute,
    .recv_hostif_packet = l_recv_hostif_packet,
    .send_hostif_packet = l_send_hostif_packet,
    .allocate_hostif_packet = l_allocate_hostif_packet,
    .free_hostif_packet = l_free_hostif_packet,
};

sai_status_t l_create_hostif(sai_object_id_t *hostif_id,
                             sai_object_id_t switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHostifRequest req;
  lemming::dataplane::sai::CreateHostifResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::HostifType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_HOSTIF_ATTR_OBJ_ID:
        req.set_obj_id(attr_list[i].value.oid);
        break;
      case SAI_HOSTIF_ATTR_NAME:
        req.set_name(attr_list[i].value.chardata);
        break;
      case SAI_HOSTIF_ATTR_OPER_STATUS:
        req.set_oper_status(attr_list[i].value.booldata);
        break;
      case SAI_HOSTIF_ATTR_QUEUE:
        req.set_queue(attr_list[i].value.u32);
        break;
      case SAI_HOSTIF_ATTR_VLAN_TAG:
        req.set_vlan_tag(static_cast<lemming::dataplane::sai::HostifVlanTag>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_HOSTIF_ATTR_GENETLINK_MCGRP_NAME:
        req.set_genetlink_mcgrp_name(attr_list[i].value.chardata);
        break;
    }
  }
  grpc::Status status = hostif->CreateHostif(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *hostif_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_HOSTIF, hostif_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_hostif(sai_object_id_t hostif_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_HOSTIF, hostif_id);
}

sai_status_t l_set_hostif_attribute(sai_object_id_t hostif_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_HOSTIF, hostif_id, attr);
}

sai_status_t l_get_hostif_attribute(sai_object_id_t hostif_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_HOSTIF, hostif_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_hostif_table_entry(sai_object_id_t *hostif_table_entry_id,
                                         sai_object_id_t switch_id,
                                         uint32_t attr_count,
                                         const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHostifTableEntryRequest req;
  lemming::dataplane::sai::CreateHostifTableEntryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::HostifTableEntryType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_OBJ_ID:
        req.set_obj_id(attr_list[i].value.oid);
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_TRAP_ID:
        req.set_trap_id(attr_list[i].value.oid);
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_CHANNEL_TYPE:
        req.set_channel_type(
            static_cast<lemming::dataplane::sai::HostifTableEntryChannelType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_HOST_IF:
        req.set_host_if(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = hostif->CreateHostifTableEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *hostif_table_entry_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY,
                            hostif_table_entry_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_hostif_table_entry(
    sai_object_id_t hostif_table_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY,
                            hostif_table_entry_id);
}

sai_status_t l_set_hostif_table_entry_attribute(
    sai_object_id_t hostif_table_entry_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY,
                                   hostif_table_entry_id, attr);
}

sai_status_t l_get_hostif_table_entry_attribute(
    sai_object_id_t hostif_table_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY,
                                   hostif_table_entry_id, attr_count,
                                   attr_list);
}

sai_status_t l_create_hostif_trap_group(sai_object_id_t *hostif_trap_group_id,
                                        sai_object_id_t switch_id,
                                        uint32_t attr_count,
                                        const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHostifTrapGroupRequest req;
  lemming::dataplane::sai::CreateHostifTrapGroupResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_TRAP_GROUP_ATTR_ADMIN_STATE:
        req.set_admin_state(attr_list[i].value.booldata);
        break;
      case SAI_HOSTIF_TRAP_GROUP_ATTR_QUEUE:
        req.set_queue(attr_list[i].value.u32);
        break;
      case SAI_HOSTIF_TRAP_GROUP_ATTR_POLICER:
        req.set_policer(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = hostif->CreateHostifTrapGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *hostif_trap_group_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP,
                            hostif_trap_group_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_hostif_trap_group(sai_object_id_t hostif_trap_group_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP,
                            hostif_trap_group_id);
}

sai_status_t l_set_hostif_trap_group_attribute(
    sai_object_id_t hostif_trap_group_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP,
                                   hostif_trap_group_id, attr);
}

sai_status_t l_get_hostif_trap_group_attribute(
    sai_object_id_t hostif_trap_group_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP,
                                   hostif_trap_group_id, attr_count, attr_list);
}

sai_status_t l_create_hostif_trap(sai_object_id_t *hostif_trap_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHostifTrapRequest req;
  lemming::dataplane::sai::CreateHostifTrapResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_TRAP_ATTR_TRAP_TYPE:
        req.set_trap_type(static_cast<lemming::dataplane::sai::HostifTrapType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_HOSTIF_TRAP_ATTR_PACKET_ACTION:
        req.set_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_HOSTIF_TRAP_ATTR_TRAP_PRIORITY:
        req.set_trap_priority(attr_list[i].value.u32);
        break;
      case SAI_HOSTIF_TRAP_ATTR_EXCLUDE_PORT_LIST:
        req.mutable_exclude_port_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_HOSTIF_TRAP_ATTR_TRAP_GROUP:
        req.set_trap_group(attr_list[i].value.oid);
        break;
      case SAI_HOSTIF_TRAP_ATTR_MIRROR_SESSION:
        req.mutable_mirror_session()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_HOSTIF_TRAP_ATTR_COUNTER_ID:
        req.set_counter_id(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = hostif->CreateHostifTrap(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *hostif_trap_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_HOSTIF_TRAP, hostif_trap_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_hostif_trap(sai_object_id_t hostif_trap_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_HOSTIF_TRAP, hostif_trap_id);
}

sai_status_t l_set_hostif_trap_attribute(sai_object_id_t hostif_trap_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_HOSTIF_TRAP, hostif_trap_id,
                                   attr);
}

sai_status_t l_get_hostif_trap_attribute(sai_object_id_t hostif_trap_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_HOSTIF_TRAP, hostif_trap_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_hostif_user_defined_trap(
    sai_object_id_t *hostif_user_defined_trap_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHostifUserDefinedTrapRequest req;
  lemming::dataplane::sai::CreateHostifUserDefinedTrapResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TYPE:
        req.set_type(
            static_cast<lemming::dataplane::sai::HostifUserDefinedTrapType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_PRIORITY:
        req.set_trap_priority(attr_list[i].value.u32);
        break;
      case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_GROUP:
        req.set_trap_group(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status =
      hostif->CreateHostifUserDefinedTrap(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *hostif_user_defined_trap_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP,
                            hostif_user_defined_trap_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_hostif_user_defined_trap(
    sai_object_id_t hostif_user_defined_trap_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP,
                            hostif_user_defined_trap_id);
}

sai_status_t l_set_hostif_user_defined_trap_attribute(
    sai_object_id_t hostif_user_defined_trap_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP,
                                   hostif_user_defined_trap_id, attr);
}

sai_status_t l_get_hostif_user_defined_trap_attribute(
    sai_object_id_t hostif_user_defined_trap_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP,
                                   hostif_user_defined_trap_id, attr_count,
                                   attr_list);
}

sai_status_t l_recv_hostif_packet(sai_object_id_t hostif_id,
                                  sai_size_t *buffer_size, void *buffer,
                                  uint32_t *attr_count,
                                  sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_send_hostif_packet(sai_object_id_t hostif_id,
                                  sai_size_t buffer_size, const void *buffer,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_allocate_hostif_packet(sai_object_id_t hostif_id,
                                      sai_size_t buffer_size, void **buffer,
                                      uint32_t attr_count,
                                      const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_free_hostif_packet(sai_object_id_t hostif_id, void *buffer) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}
