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

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/hostif.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

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

lemming::dataplane::sai::CreateHostifRequest convert_create_hostif(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateHostifRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_ATTR_TYPE:
        msg.set_type(
            convert_sai_hostif_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_HOSTIF_ATTR_OBJ_ID:
        msg.set_obj_id(attr_list[i].value.oid);
        break;
      case SAI_HOSTIF_ATTR_NAME:
        msg.set_name(attr_list[i].value.chardata);
        break;
      case SAI_HOSTIF_ATTR_OPER_STATUS:
        msg.set_oper_status(attr_list[i].value.booldata);
        break;
      case SAI_HOSTIF_ATTR_QUEUE:
        msg.set_queue(attr_list[i].value.u32);
        break;
      case SAI_HOSTIF_ATTR_VLAN_TAG:
        msg.set_vlan_tag(
            convert_sai_hostif_vlan_tag_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_HOSTIF_ATTR_GENETLINK_MCGRP_NAME:
        msg.set_genetlink_mcgrp_name(attr_list[i].value.chardata);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateHostifTableEntryRequest
convert_create_hostif_table_entry(sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateHostifTableEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_TYPE:
        msg.set_type(convert_sai_hostif_table_entry_type_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_OBJ_ID:
        msg.set_obj_id(attr_list[i].value.oid);
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_TRAP_ID:
        msg.set_trap_id(attr_list[i].value.oid);
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_CHANNEL_TYPE:
        msg.set_channel_type(
            convert_sai_hostif_table_entry_channel_type_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_HOST_IF:
        msg.set_host_if(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateHostifTrapGroupRequest
convert_create_hostif_trap_group(sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateHostifTrapGroupRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_TRAP_GROUP_ATTR_ADMIN_STATE:
        msg.set_admin_state(attr_list[i].value.booldata);
        break;
      case SAI_HOSTIF_TRAP_GROUP_ATTR_QUEUE:
        msg.set_queue(attr_list[i].value.u32);
        break;
      case SAI_HOSTIF_TRAP_GROUP_ATTR_POLICER:
        msg.set_policer(attr_list[i].value.oid);
        break;
      case SAI_HOSTIF_TRAP_GROUP_ATTR_OBJECT_STAGE:
        msg.set_object_stage(
            convert_sai_object_stage_t_to_proto(attr_list[i].value.s32));
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateHostifTrapRequest convert_create_hostif_trap(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateHostifTrapRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_TRAP_ATTR_TRAP_TYPE:
        msg.set_trap_type(
            convert_sai_hostif_trap_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_HOSTIF_TRAP_ATTR_PACKET_ACTION:
        msg.set_packet_action(
            convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_HOSTIF_TRAP_ATTR_TRAP_PRIORITY:
        msg.set_trap_priority(attr_list[i].value.u32);
        break;
      case SAI_HOSTIF_TRAP_ATTR_EXCLUDE_PORT_LIST:
        msg.mutable_exclude_port_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_HOSTIF_TRAP_ATTR_TRAP_GROUP:
        msg.set_trap_group(attr_list[i].value.oid);
        break;
      case SAI_HOSTIF_TRAP_ATTR_MIRROR_SESSION:
        msg.mutable_mirror_session()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_HOSTIF_TRAP_ATTR_COUNTER_ID:
        msg.set_counter_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateHostifUserDefinedTrapRequest
convert_create_hostif_user_defined_trap(sai_object_id_t switch_id,
                                        uint32_t attr_count,
                                        const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateHostifUserDefinedTrapRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TYPE:
        msg.set_type(convert_sai_hostif_user_defined_trap_type_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_PRIORITY:
        msg.set_trap_priority(attr_list[i].value.u32);
        break;
      case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_GROUP:
        msg.set_trap_group(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_hostif(sai_object_id_t *hostif_id,
                             sai_object_id_t switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHostifRequest req =
      convert_create_hostif(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateHostifResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = hostif->CreateHostif(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (hostif_id) {
    *hostif_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_hostif(sai_object_id_t hostif_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveHostifRequest req;
  lemming::dataplane::sai::RemoveHostifResponse resp;
  grpc::ClientContext context;
  req.set_oid(hostif_id);

  grpc::Status status = hostif->RemoveHostif(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_hostif_attribute(sai_object_id_t hostif_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetHostifAttributeRequest req;
  lemming::dataplane::sai::SetHostifAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(hostif_id);

  switch (attr->id) {
    case SAI_HOSTIF_ATTR_OPER_STATUS:
      req.set_oper_status(attr->value.booldata);
      break;
    case SAI_HOSTIF_ATTR_QUEUE:
      req.set_queue(attr->value.u32);
      break;
    case SAI_HOSTIF_ATTR_VLAN_TAG:
      req.set_vlan_tag(convert_sai_hostif_vlan_tag_t_to_proto(attr->value.s32));
      break;
  }

  grpc::Status status = hostif->SetHostifAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_hostif_attribute(sai_object_id_t hostif_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetHostifAttributeRequest req;
  lemming::dataplane::sai::GetHostifAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(hostif_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_hostif_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = hostif->GetHostifAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_ATTR_TYPE:
        attr_list[i].value.s32 =
            convert_sai_hostif_type_t_to_sai(resp.attr().type());
        break;
      case SAI_HOSTIF_ATTR_OBJ_ID:
        attr_list[i].value.oid = resp.attr().obj_id();
        break;
      case SAI_HOSTIF_ATTR_NAME:
        strncpy(attr_list[i].value.chardata, resp.attr().name().data(), 32);
        break;
      case SAI_HOSTIF_ATTR_OPER_STATUS:
        attr_list[i].value.booldata = resp.attr().oper_status();
        break;
      case SAI_HOSTIF_ATTR_QUEUE:
        attr_list[i].value.u32 = resp.attr().queue();
        break;
      case SAI_HOSTIF_ATTR_VLAN_TAG:
        attr_list[i].value.s32 =
            convert_sai_hostif_vlan_tag_t_to_sai(resp.attr().vlan_tag());
        break;
      case SAI_HOSTIF_ATTR_GENETLINK_MCGRP_NAME:
        strncpy(attr_list[i].value.chardata,
                resp.attr().genetlink_mcgrp_name().data(), 32);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_hostif_table_entry(sai_object_id_t *hostif_table_entry_id,
                                         sai_object_id_t switch_id,
                                         uint32_t attr_count,
                                         const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHostifTableEntryRequest req =
      convert_create_hostif_table_entry(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateHostifTableEntryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = hostif->CreateHostifTableEntry(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (hostif_table_entry_id) {
    *hostif_table_entry_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_hostif_table_entry(
    sai_object_id_t hostif_table_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveHostifTableEntryRequest req;
  lemming::dataplane::sai::RemoveHostifTableEntryResponse resp;
  grpc::ClientContext context;
  req.set_oid(hostif_table_entry_id);

  grpc::Status status = hostif->RemoveHostifTableEntry(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_hostif_table_entry_attribute(
    sai_object_id_t hostif_table_entry_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_hostif_table_entry_attribute(
    sai_object_id_t hostif_table_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetHostifTableEntryAttributeRequest req;
  lemming::dataplane::sai::GetHostifTableEntryAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(hostif_table_entry_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_hostif_table_entry_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status =
      hostif->GetHostifTableEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_TYPE:
        attr_list[i].value.s32 =
            convert_sai_hostif_table_entry_type_t_to_sai(resp.attr().type());
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_OBJ_ID:
        attr_list[i].value.oid = resp.attr().obj_id();
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_TRAP_ID:
        attr_list[i].value.oid = resp.attr().trap_id();
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_CHANNEL_TYPE:
        attr_list[i].value.s32 =
            convert_sai_hostif_table_entry_channel_type_t_to_sai(
                resp.attr().channel_type());
        break;
      case SAI_HOSTIF_TABLE_ENTRY_ATTR_HOST_IF:
        attr_list[i].value.oid = resp.attr().host_if();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_hostif_trap_group(sai_object_id_t *hostif_trap_group_id,
                                        sai_object_id_t switch_id,
                                        uint32_t attr_count,
                                        const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHostifTrapGroupRequest req =
      convert_create_hostif_trap_group(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateHostifTrapGroupResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = hostif->CreateHostifTrapGroup(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (hostif_trap_group_id) {
    *hostif_trap_group_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_hostif_trap_group(sai_object_id_t hostif_trap_group_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveHostifTrapGroupRequest req;
  lemming::dataplane::sai::RemoveHostifTrapGroupResponse resp;
  grpc::ClientContext context;
  req.set_oid(hostif_trap_group_id);

  grpc::Status status = hostif->RemoveHostifTrapGroup(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_hostif_trap_group_attribute(
    sai_object_id_t hostif_trap_group_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetHostifTrapGroupAttributeRequest req;
  lemming::dataplane::sai::SetHostifTrapGroupAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(hostif_trap_group_id);

  switch (attr->id) {
    case SAI_HOSTIF_TRAP_GROUP_ATTR_ADMIN_STATE:
      req.set_admin_state(attr->value.booldata);
      break;
    case SAI_HOSTIF_TRAP_GROUP_ATTR_QUEUE:
      req.set_queue(attr->value.u32);
      break;
    case SAI_HOSTIF_TRAP_GROUP_ATTR_POLICER:
      req.set_policer(attr->value.oid);
      break;
  }

  grpc::Status status =
      hostif->SetHostifTrapGroupAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_hostif_trap_group_attribute(
    sai_object_id_t hostif_trap_group_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetHostifTrapGroupAttributeRequest req;
  lemming::dataplane::sai::GetHostifTrapGroupAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(hostif_trap_group_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_hostif_trap_group_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status =
      hostif->GetHostifTrapGroupAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_TRAP_GROUP_ATTR_ADMIN_STATE:
        attr_list[i].value.booldata = resp.attr().admin_state();
        break;
      case SAI_HOSTIF_TRAP_GROUP_ATTR_QUEUE:
        attr_list[i].value.u32 = resp.attr().queue();
        break;
      case SAI_HOSTIF_TRAP_GROUP_ATTR_POLICER:
        attr_list[i].value.oid = resp.attr().policer();
        break;
      case SAI_HOSTIF_TRAP_GROUP_ATTR_OBJECT_STAGE:
        attr_list[i].value.s32 =
            convert_sai_object_stage_t_to_sai(resp.attr().object_stage());
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_hostif_trap(sai_object_id_t *hostif_trap_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHostifTrapRequest req =
      convert_create_hostif_trap(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateHostifTrapResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = hostif->CreateHostifTrap(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (hostif_trap_id) {
    *hostif_trap_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_hostif_trap(sai_object_id_t hostif_trap_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveHostifTrapRequest req;
  lemming::dataplane::sai::RemoveHostifTrapResponse resp;
  grpc::ClientContext context;
  req.set_oid(hostif_trap_id);

  grpc::Status status = hostif->RemoveHostifTrap(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_hostif_trap_attribute(sai_object_id_t hostif_trap_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetHostifTrapAttributeRequest req;
  lemming::dataplane::sai::SetHostifTrapAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(hostif_trap_id);

  switch (attr->id) {
    case SAI_HOSTIF_TRAP_ATTR_PACKET_ACTION:
      req.set_packet_action(
          convert_sai_packet_action_t_to_proto(attr->value.s32));
      break;
    case SAI_HOSTIF_TRAP_ATTR_TRAP_PRIORITY:
      req.set_trap_priority(attr->value.u32);
      break;
    case SAI_HOSTIF_TRAP_ATTR_EXCLUDE_PORT_LIST:
      req.mutable_exclude_port_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_HOSTIF_TRAP_ATTR_TRAP_GROUP:
      req.set_trap_group(attr->value.oid);
      break;
    case SAI_HOSTIF_TRAP_ATTR_MIRROR_SESSION:
      req.mutable_mirror_session()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_HOSTIF_TRAP_ATTR_COUNTER_ID:
      req.set_counter_id(attr->value.oid);
      break;
  }

  grpc::Status status = hostif->SetHostifTrapAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_hostif_trap_attribute(sai_object_id_t hostif_trap_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetHostifTrapAttributeRequest req;
  lemming::dataplane::sai::GetHostifTrapAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(hostif_trap_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_hostif_trap_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = hostif->GetHostifTrapAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_TRAP_ATTR_TRAP_TYPE:
        attr_list[i].value.s32 =
            convert_sai_hostif_trap_type_t_to_sai(resp.attr().trap_type());
        break;
      case SAI_HOSTIF_TRAP_ATTR_PACKET_ACTION:
        attr_list[i].value.s32 =
            convert_sai_packet_action_t_to_sai(resp.attr().packet_action());
        break;
      case SAI_HOSTIF_TRAP_ATTR_TRAP_PRIORITY:
        attr_list[i].value.u32 = resp.attr().trap_priority();
        break;
      case SAI_HOSTIF_TRAP_ATTR_EXCLUDE_PORT_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().exclude_port_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_HOSTIF_TRAP_ATTR_TRAP_GROUP:
        attr_list[i].value.oid = resp.attr().trap_group();
        break;
      case SAI_HOSTIF_TRAP_ATTR_MIRROR_SESSION:
        copy_list(attr_list[i].value.objlist.list, resp.attr().mirror_session(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_HOSTIF_TRAP_ATTR_COUNTER_ID:
        attr_list[i].value.oid = resp.attr().counter_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_hostif_user_defined_trap(
    sai_object_id_t *hostif_user_defined_trap_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHostifUserDefinedTrapRequest req =
      convert_create_hostif_user_defined_trap(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateHostifUserDefinedTrapResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status =
      hostif->CreateHostifUserDefinedTrap(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (hostif_user_defined_trap_id) {
    *hostif_user_defined_trap_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_hostif_user_defined_trap(
    sai_object_id_t hostif_user_defined_trap_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveHostifUserDefinedTrapRequest req;
  lemming::dataplane::sai::RemoveHostifUserDefinedTrapResponse resp;
  grpc::ClientContext context;
  req.set_oid(hostif_user_defined_trap_id);

  grpc::Status status =
      hostif->RemoveHostifUserDefinedTrap(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_hostif_user_defined_trap_attribute(
    sai_object_id_t hostif_user_defined_trap_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetHostifUserDefinedTrapAttributeRequest req;
  lemming::dataplane::sai::SetHostifUserDefinedTrapAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(hostif_user_defined_trap_id);

  switch (attr->id) {
    case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_PRIORITY:
      req.set_trap_priority(attr->value.u32);
      break;
    case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_GROUP:
      req.set_trap_group(attr->value.oid);
      break;
  }

  grpc::Status status =
      hostif->SetHostifUserDefinedTrapAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_hostif_user_defined_trap_attribute(
    sai_object_id_t hostif_user_defined_trap_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetHostifUserDefinedTrapAttributeRequest req;
  lemming::dataplane::sai::GetHostifUserDefinedTrapAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(hostif_user_defined_trap_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_hostif_user_defined_trap_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status =
      hostif->GetHostifUserDefinedTrapAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TYPE:
        attr_list[i].value.s32 =
            convert_sai_hostif_user_defined_trap_type_t_to_sai(
                resp.attr().type());
        break;
      case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_PRIORITY:
        attr_list[i].value.u32 = resp.attr().trap_priority();
        break;
      case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_GROUP:
        attr_list[i].value.oid = resp.attr().trap_group();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
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
