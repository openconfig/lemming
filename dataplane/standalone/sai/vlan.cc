

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

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/vlan.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

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

sai_status_t l_create_vlan(sai_object_id_t *vlan_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateVlanRequest req;
  lemming::dataplane::sai::CreateVlanResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_VLAN_ATTR_VLAN_ID:
        req.set_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_VLAN_ATTR_MAX_LEARNED_ADDRESSES:
        req.set_max_learned_addresses(attr_list[i].value.u32);
        break;
      case SAI_VLAN_ATTR_STP_INSTANCE:
        req.set_stp_instance(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_LEARN_DISABLE:
        req.set_learn_disable(attr_list[i].value.booldata);
        break;
      case SAI_VLAN_ATTR_IPV4_MCAST_LOOKUP_KEY_TYPE:
        req.set_ipv4_mcast_lookup_key_type(
            static_cast<lemming::dataplane::sai::VlanMcastLookupKeyType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_VLAN_ATTR_IPV6_MCAST_LOOKUP_KEY_TYPE:
        req.set_ipv6_mcast_lookup_key_type(
            static_cast<lemming::dataplane::sai::VlanMcastLookupKeyType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_VLAN_ATTR_UNKNOWN_NON_IP_MCAST_OUTPUT_GROUP_ID:
        req.set_unknown_non_ip_mcast_output_group_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_UNKNOWN_IPV4_MCAST_OUTPUT_GROUP_ID:
        req.set_unknown_ipv4_mcast_output_group_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_UNKNOWN_IPV6_MCAST_OUTPUT_GROUP_ID:
        req.set_unknown_ipv6_mcast_output_group_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_UNKNOWN_LINKLOCAL_MCAST_OUTPUT_GROUP_ID:
        req.set_unknown_linklocal_mcast_output_group_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_INGRESS_ACL:
        req.set_ingress_acl(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_EGRESS_ACL:
        req.set_egress_acl(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_META_DATA:
        req.set_meta_data(attr_list[i].value.u32);
        break;
      case SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
        req.set_unknown_unicast_flood_control_type(
            static_cast<lemming::dataplane::sai::VlanFloodControlType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
        req.set_unknown_unicast_flood_group(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
        req.set_unknown_multicast_flood_control_type(
            static_cast<lemming::dataplane::sai::VlanFloodControlType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
        req.set_unknown_multicast_flood_group(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
        req.set_broadcast_flood_control_type(
            static_cast<lemming::dataplane::sai::VlanFloodControlType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_VLAN_ATTR_BROADCAST_FLOOD_GROUP:
        req.set_broadcast_flood_group(attr_list[i].value.oid);
        break;
      case SAI_VLAN_ATTR_CUSTOM_IGMP_SNOOPING_ENABLE:
        req.set_custom_igmp_snooping_enable(attr_list[i].value.booldata);
        break;
      case SAI_VLAN_ATTR_TAM_OBJECT:
        req.mutable_tam_object()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
    }
  }
  grpc::Status status = vlan->CreateVlan(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *vlan_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_VLAN, vlan_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_vlan(sai_object_id_t vlan_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_VLAN, vlan_id);
}

sai_status_t l_set_vlan_attribute(sai_object_id_t vlan_id,
                                  const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_VLAN, vlan_id, attr);
}

sai_status_t l_get_vlan_attribute(sai_object_id_t vlan_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_VLAN, vlan_id, attr_count,
                                   attr_list);
}

sai_status_t l_create_vlan_member(sai_object_id_t *vlan_member_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateVlanMemberRequest req;
  lemming::dataplane::sai::CreateVlanMemberResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_VLAN_MEMBER_ATTR_VLAN_ID:
        req.set_vlan_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_MEMBER_ATTR_BRIDGE_PORT_ID:
        req.set_bridge_port_id(attr_list[i].value.oid);
        break;
      case SAI_VLAN_MEMBER_ATTR_VLAN_TAGGING_MODE:
        req.set_vlan_tagging_mode(
            static_cast<lemming::dataplane::sai::VlanTaggingMode>(
                attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = vlan->CreateVlanMember(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *vlan_member_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_VLAN_MEMBER, vlan_member_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_vlan_member(sai_object_id_t vlan_member_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_VLAN_MEMBER, vlan_member_id);
}

sai_status_t l_set_vlan_member_attribute(sai_object_id_t vlan_member_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_VLAN_MEMBER, vlan_member_id,
                                   attr);
}

sai_status_t l_get_vlan_member_attribute(sai_object_id_t vlan_member_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_VLAN_MEMBER, vlan_member_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_vlan_members(sai_object_id_t switch_id,
                                   uint32_t object_count,
                                   const uint32_t *attr_count,
                                   const sai_attribute_t **attr_list,
                                   sai_bulk_op_error_mode_t mode,
                                   sai_object_id_t *object_id,
                                   sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->create_bulk(SAI_OBJECT_TYPE_VLAN_MEMBER, switch_id,
                                 object_count, attr_count, attr_list, mode,
                                 object_id, object_statuses);
}

sai_status_t l_remove_vlan_members(uint32_t object_count,
                                   const sai_object_id_t *object_id,
                                   sai_bulk_op_error_mode_t mode,
                                   sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove_bulk(SAI_OBJECT_TYPE_VLAN_MEMBER, object_count,
                                 object_id, mode, object_statuses);
}

sai_status_t l_get_vlan_stats(sai_object_id_t vlan_id,
                              uint32_t number_of_counters,
                              const sai_stat_id_t *counter_ids,
                              uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_VLAN, vlan_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_vlan_stats_ext(sai_object_id_t vlan_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids,
                                  sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_VLAN, vlan_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_vlan_stats(sai_object_id_t vlan_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_VLAN, vlan_id,
                                 number_of_counters, counter_ids);
}
