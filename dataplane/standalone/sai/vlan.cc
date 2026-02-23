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

#include "dataplane/standalone/sai/vlan.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/vlan.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_vlan_api_t l_vlan = {
    .create_vlan = l_create_vlan,
    .remove_vlan = l_remove_vlan,
    .set_vlan_attribute = l_set_vlan_attribute,
    .get_vlan_attribute = l_get_vlan_attribute,
    .create_vlan_member = l_create_vlan_member,
    .remove_vlan_member = l_remove_vlan_member,
    .set_vlan_member_attribute = l_set_vlan_member_attribute,
    .get_vlan_member_attribute = l_get_vlan_member_attribute,
    .create_vlan_members = l_create_vlan_members,
    .remove_vlan_members = l_remove_vlan_members,
    .get_vlan_stats = l_get_vlan_stats,
    .get_vlan_stats_ext = l_get_vlan_stats_ext,
    .clear_vlan_stats = l_clear_vlan_stats,
};

lemming::dataplane::sai::CreateVlanRequest convert_create_vlan(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t* attr_list) {
  lemming::dataplane::sai::CreateVlanRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_VLAN_ATTR_VLAN_ID:
        msg.set_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_VLAN_ATTR_MAX_LEARNED_ADDRESSES:
        msg.set_max_learned_addresses(attr_list[i].value.u32);
        break;
      case SAI_VLAN_ATTR_STP_INSTANCE:
        msg.set_stp_instance(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_LEARN_DISABLE:
        msg.set_learn_disable(attr_list[i].value.booldata);
        break;
      case SAI_VLAN_ATTR_IPV4_MCAST_LOOKUP_KEY_TYPE:
        msg.set_ipv4_mcast_lookup_key_type(
            convert_sai_vlan_mcast_lookup_key_type_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_VLAN_ATTR_IPV6_MCAST_LOOKUP_KEY_TYPE:
        msg.set_ipv6_mcast_lookup_key_type(
            convert_sai_vlan_mcast_lookup_key_type_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_VLAN_ATTR_UNKNOWN_NON_IP_MCAST_OUTPUT_GROUP_ID:
        msg.set_unknown_non_ip_mcast_output_group_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_UNKNOWN_IPV4_MCAST_OUTPUT_GROUP_ID:
        msg.set_unknown_ipv4_mcast_output_group_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_UNKNOWN_IPV6_MCAST_OUTPUT_GROUP_ID:
        msg.set_unknown_ipv6_mcast_output_group_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_UNKNOWN_LINKLOCAL_MCAST_OUTPUT_GROUP_ID:
        msg.set_unknown_linklocal_mcast_output_group_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_INGRESS_ACL:
        msg.set_ingress_acl(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_EGRESS_ACL:
        msg.set_egress_acl(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_META_DATA:
        msg.set_meta_data(attr_list[i].value.u32);
        break;
      case SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
        msg.set_unknown_unicast_flood_control_type(
            convert_sai_vlan_flood_control_type_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
        msg.set_unknown_unicast_flood_group(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
        msg.set_unknown_multicast_flood_control_type(
            convert_sai_vlan_flood_control_type_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
        msg.set_unknown_multicast_flood_group(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
        msg.set_broadcast_flood_control_type(
            convert_sai_vlan_flood_control_type_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_VLAN_ATTR_BROADCAST_FLOOD_GROUP:
        msg.set_broadcast_flood_group(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_CUSTOM_IGMP_SNOOPING_ENABLE:
        msg.set_custom_igmp_snooping_enable(attr_list[i].value.booldata);
        break;
      case SAI_VLAN_ATTR_TAM_OBJECT:
        msg.mutable_tam_object()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_VLAN_ATTR_STATS_COUNT_MODE:
        msg.set_stats_count_mode(
            convert_sai_stats_count_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_VLAN_ATTR_SELECTIVE_COUNTER_LIST:
        msg.mutable_selective_counter_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateVlanMemberRequest convert_create_vlan_member(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t* attr_list) {
  lemming::dataplane::sai::CreateVlanMemberRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_VLAN_MEMBER_ATTR_VLAN_ID:
        msg.set_vlan_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_MEMBER_ATTR_BRIDGE_PORT_ID:
        msg.set_bridge_port_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_MEMBER_ATTR_VLAN_TAGGING_MODE:
        msg.set_vlan_tagging_mode(
            convert_sai_vlan_tagging_mode_t_to_proto(attr_list[i].value.s32));
        break;
    }
  }
  return msg;
}

sai_status_t l_create_vlan(sai_object_id_t* vlan_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateVlanRequest req =
      convert_create_vlan(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateVlanResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = vlan->CreateVlan(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (vlan_id) {
    *vlan_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_vlan(sai_object_id_t vlan_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveVlanRequest req;
  lemming::dataplane::sai::RemoveVlanResponse resp;
  grpc::ClientContext context;
  req.set_oid(vlan_id);

  grpc::Status status = vlan->RemoveVlan(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_vlan_attribute(sai_object_id_t vlan_id,
                                  const sai_attribute_t* attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetVlanAttributeRequest req;
  lemming::dataplane::sai::SetVlanAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(vlan_id);

  switch (attr->id) {
    case SAI_VLAN_ATTR_MAX_LEARNED_ADDRESSES:
      req.set_max_learned_addresses(attr->value.u32);
      break;
    case SAI_VLAN_ATTR_STP_INSTANCE:
      req.set_stp_instance(attr->value.oid);
      break;
    case SAI_VLAN_ATTR_LEARN_DISABLE:
      req.set_learn_disable(attr->value.booldata);
      break;
    case SAI_VLAN_ATTR_IPV4_MCAST_LOOKUP_KEY_TYPE:
      req.set_ipv4_mcast_lookup_key_type(
          convert_sai_vlan_mcast_lookup_key_type_t_to_proto(attr->value.s32));
      break;
    case SAI_VLAN_ATTR_IPV6_MCAST_LOOKUP_KEY_TYPE:
      req.set_ipv6_mcast_lookup_key_type(
          convert_sai_vlan_mcast_lookup_key_type_t_to_proto(attr->value.s32));
      break;
    case SAI_VLAN_ATTR_UNKNOWN_NON_IP_MCAST_OUTPUT_GROUP_ID:
      req.set_unknown_non_ip_mcast_output_group_id(attr->value.oid);
      break;
    case SAI_VLAN_ATTR_UNKNOWN_IPV4_MCAST_OUTPUT_GROUP_ID:
      req.set_unknown_ipv4_mcast_output_group_id(attr->value.oid);
      break;
    case SAI_VLAN_ATTR_UNKNOWN_IPV6_MCAST_OUTPUT_GROUP_ID:
      req.set_unknown_ipv6_mcast_output_group_id(attr->value.oid);
      break;
    case SAI_VLAN_ATTR_UNKNOWN_LINKLOCAL_MCAST_OUTPUT_GROUP_ID:
      req.set_unknown_linklocal_mcast_output_group_id(attr->value.oid);
      break;
    case SAI_VLAN_ATTR_INGRESS_ACL:
      req.set_ingress_acl(attr->value.oid);
      break;
    case SAI_VLAN_ATTR_EGRESS_ACL:
      req.set_egress_acl(attr->value.oid);
      break;
    case SAI_VLAN_ATTR_META_DATA:
      req.set_meta_data(attr->value.u32);
      break;
    case SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
      req.set_unknown_unicast_flood_control_type(
          convert_sai_vlan_flood_control_type_t_to_proto(attr->value.s32));
      break;
    case SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
      req.set_unknown_unicast_flood_group(attr->value.oid);
      break;
    case SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
      req.set_unknown_multicast_flood_control_type(
          convert_sai_vlan_flood_control_type_t_to_proto(attr->value.s32));
      break;
    case SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
      req.set_unknown_multicast_flood_group(attr->value.oid);
      break;
    case SAI_VLAN_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
      req.set_broadcast_flood_control_type(
          convert_sai_vlan_flood_control_type_t_to_proto(attr->value.s32));
      break;
    case SAI_VLAN_ATTR_BROADCAST_FLOOD_GROUP:
      req.set_broadcast_flood_group(attr->value.oid);
      break;
    case SAI_VLAN_ATTR_CUSTOM_IGMP_SNOOPING_ENABLE:
      req.set_custom_igmp_snooping_enable(attr->value.booldata);
      break;
    case SAI_VLAN_ATTR_TAM_OBJECT:
      req.mutable_tam_object()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_VLAN_ATTR_STATS_COUNT_MODE:
      req.set_stats_count_mode(
          convert_sai_stats_count_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_VLAN_ATTR_SELECTIVE_COUNTER_LIST:
      req.mutable_selective_counter_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
  }

  grpc::Status status = vlan->SetVlanAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_vlan_attribute(sai_object_id_t vlan_id, uint32_t attr_count,
                                  sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetVlanAttributeRequest req;
  lemming::dataplane::sai::GetVlanAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(vlan_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_vlan_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = vlan->GetVlanAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_VLAN_ATTR_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().vlan_id();
        break;
      case SAI_VLAN_ATTR_MEMBER_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().member_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_VLAN_ATTR_MAX_LEARNED_ADDRESSES:
        attr_list[i].value.u32 = resp.attr().max_learned_addresses();
        break;
      case SAI_VLAN_ATTR_STP_INSTANCE:
        attr_list[i].value.oid = resp.attr().stp_instance();
        break;
      case SAI_VLAN_ATTR_LEARN_DISABLE:
        attr_list[i].value.booldata = resp.attr().learn_disable();
        break;
      case SAI_VLAN_ATTR_IPV4_MCAST_LOOKUP_KEY_TYPE:
        attr_list[i].value.s32 =
            convert_sai_vlan_mcast_lookup_key_type_t_to_sai(
                resp.attr().ipv4_mcast_lookup_key_type());
        break;
      case SAI_VLAN_ATTR_IPV6_MCAST_LOOKUP_KEY_TYPE:
        attr_list[i].value.s32 =
            convert_sai_vlan_mcast_lookup_key_type_t_to_sai(
                resp.attr().ipv6_mcast_lookup_key_type());
        break;
      case SAI_VLAN_ATTR_UNKNOWN_NON_IP_MCAST_OUTPUT_GROUP_ID:
        attr_list[i].value.oid =
            resp.attr().unknown_non_ip_mcast_output_group_id();
        break;
      case SAI_VLAN_ATTR_UNKNOWN_IPV4_MCAST_OUTPUT_GROUP_ID:
        attr_list[i].value.oid =
            resp.attr().unknown_ipv4_mcast_output_group_id();
        break;
      case SAI_VLAN_ATTR_UNKNOWN_IPV6_MCAST_OUTPUT_GROUP_ID:
        attr_list[i].value.oid =
            resp.attr().unknown_ipv6_mcast_output_group_id();
        break;
      case SAI_VLAN_ATTR_UNKNOWN_LINKLOCAL_MCAST_OUTPUT_GROUP_ID:
        attr_list[i].value.oid =
            resp.attr().unknown_linklocal_mcast_output_group_id();
        break;
      case SAI_VLAN_ATTR_INGRESS_ACL:
        attr_list[i].value.oid = resp.attr().ingress_acl();
        break;
      case SAI_VLAN_ATTR_EGRESS_ACL:
        attr_list[i].value.oid = resp.attr().egress_acl();
        break;
      case SAI_VLAN_ATTR_META_DATA:
        attr_list[i].value.u32 = resp.attr().meta_data();
        break;
      case SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
        attr_list[i].value.s32 = convert_sai_vlan_flood_control_type_t_to_sai(
            resp.attr().unknown_unicast_flood_control_type());
        break;
      case SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
        attr_list[i].value.oid = resp.attr().unknown_unicast_flood_group();
        break;
      case SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
        attr_list[i].value.s32 = convert_sai_vlan_flood_control_type_t_to_sai(
            resp.attr().unknown_multicast_flood_control_type());
        break;
      case SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
        attr_list[i].value.oid = resp.attr().unknown_multicast_flood_group();
        break;
      case SAI_VLAN_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
        attr_list[i].value.s32 = convert_sai_vlan_flood_control_type_t_to_sai(
            resp.attr().broadcast_flood_control_type());
        break;
      case SAI_VLAN_ATTR_BROADCAST_FLOOD_GROUP:
        attr_list[i].value.oid = resp.attr().broadcast_flood_group();
        break;
      case SAI_VLAN_ATTR_CUSTOM_IGMP_SNOOPING_ENABLE:
        attr_list[i].value.booldata = resp.attr().custom_igmp_snooping_enable();
        break;
      case SAI_VLAN_ATTR_TAM_OBJECT:
        copy_list(attr_list[i].value.objlist.list, resp.attr().tam_object(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_VLAN_ATTR_STATS_COUNT_MODE:
        attr_list[i].value.s32 = convert_sai_stats_count_mode_t_to_sai(
            resp.attr().stats_count_mode());
        break;
      case SAI_VLAN_ATTR_SELECTIVE_COUNTER_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().selective_counter_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_vlan_member(sai_object_id_t* vlan_member_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateVlanMemberRequest req =
      convert_create_vlan_member(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateVlanMemberResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = vlan->CreateVlanMember(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (vlan_member_id) {
    *vlan_member_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_vlan_member(sai_object_id_t vlan_member_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveVlanMemberRequest req;
  lemming::dataplane::sai::RemoveVlanMemberResponse resp;
  grpc::ClientContext context;
  req.set_oid(vlan_member_id);

  grpc::Status status = vlan->RemoveVlanMember(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_vlan_member_attribute(sai_object_id_t vlan_member_id,
                                         const sai_attribute_t* attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetVlanMemberAttributeRequest req;
  lemming::dataplane::sai::SetVlanMemberAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(vlan_member_id);

  switch (attr->id) {
    case SAI_VLAN_MEMBER_ATTR_VLAN_TAGGING_MODE:
      req.set_vlan_tagging_mode(
          convert_sai_vlan_tagging_mode_t_to_proto(attr->value.s32));
      break;
  }

  grpc::Status status = vlan->SetVlanMemberAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_vlan_member_attribute(sai_object_id_t vlan_member_id,
                                         uint32_t attr_count,
                                         sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetVlanMemberAttributeRequest req;
  lemming::dataplane::sai::GetVlanMemberAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(vlan_member_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_vlan_member_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = vlan->GetVlanMemberAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_VLAN_MEMBER_ATTR_VLAN_ID:
        attr_list[i].value.oid = resp.attr().vlan_id();
        break;
      case SAI_VLAN_MEMBER_ATTR_BRIDGE_PORT_ID:
        attr_list[i].value.oid = resp.attr().bridge_port_id();
        break;
      case SAI_VLAN_MEMBER_ATTR_VLAN_TAGGING_MODE:
        attr_list[i].value.s32 = convert_sai_vlan_tagging_mode_t_to_sai(
            resp.attr().vlan_tagging_mode());
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_vlan_members(sai_object_id_t switch_id,
                                   uint32_t object_count,
                                   const uint32_t* attr_count,
                                   const sai_attribute_t** attr_list,
                                   sai_bulk_op_error_mode_t mode,
                                   sai_object_id_t* object_id,
                                   sai_status_t* object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateVlanMembersRequest req;
  lemming::dataplane::sai::CreateVlanMembersResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    auto r = convert_create_vlan_member(switch_id, attr_count[i], attr_list[i]);
    *req.add_reqs() = r;
  }

  grpc::Status status = vlan->CreateVlanMembers(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (object_count != resp.resps().size()) {
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < object_count; i++) {
    object_id[i] = resp.resps(i).oid();
    object_statuses[i] = SAI_STATUS_SUCCESS;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_vlan_members(uint32_t object_count,
                                   const sai_object_id_t* object_id,
                                   sai_bulk_op_error_mode_t mode,
                                   sai_status_t* object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveVlanMembersRequest req;
  lemming::dataplane::sai::RemoveVlanMembersResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    req.add_reqs()->set_oid(object_id[i]);
  }

  grpc::Status status = vlan->RemoveVlanMembers(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (object_count != resp.resps().size()) {
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < object_count; i++) {
    object_statuses[i] = SAI_STATUS_SUCCESS;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_vlan_stats(sai_object_id_t vlan_id,
                              uint32_t number_of_counters,
                              const sai_stat_id_t* counter_ids,
                              uint64_t* counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetVlanStatsRequest req;
  lemming::dataplane::sai::GetVlanStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(vlan_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(convert_sai_vlan_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status = vlan->GetVlanStats(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_vlan_stats_ext(sai_object_id_t vlan_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t* counter_ids,
                                  sai_stats_mode_t mode, uint64_t* counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_vlan_stats(sai_object_id_t vlan_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t* counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
