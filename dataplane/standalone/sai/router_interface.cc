

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

#include "dataplane/standalone/sai/router_interface.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/router_interface.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_router_interface_api_t l_router_interface = {
    .create_router_interface = l_create_router_interface,
    .remove_router_interface = l_remove_router_interface,
    .set_router_interface_attribute = l_set_router_interface_attribute,
    .get_router_interface_attribute = l_get_router_interface_attribute,
    .get_router_interface_stats = l_get_router_interface_stats,
    .get_router_interface_stats_ext = l_get_router_interface_stats_ext,
    .clear_router_interface_stats = l_clear_router_interface_stats,
};

lemming::dataplane::sai::CreateRouterInterfaceRequest
convert_create_router_interface(sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateRouterInterfaceRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ROUTER_INTERFACE_ATTR_VIRTUAL_ROUTER_ID:
        msg.set_virtual_router_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_TYPE:
        msg.set_type(static_cast<lemming::dataplane::sai::RouterInterfaceType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_ROUTER_INTERFACE_ATTR_PORT_ID:
        msg.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_VLAN_ID:
        msg.set_vlan_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_OUTER_VLAN_ID:
        msg.set_outer_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_INNER_VLAN_ID:
        msg.set_inner_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_BRIDGE_ID:
        msg.set_bridge_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS:
        msg.set_src_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_ROUTER_INTERFACE_ATTR_ADMIN_V4_STATE:
        msg.set_admin_v4_state(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_ADMIN_V6_STATE:
        msg.set_admin_v6_state(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_MTU:
        msg.set_mtu(attr_list[i].value.u32);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_INGRESS_ACL:
        msg.set_ingress_acl(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_EGRESS_ACL:
        msg.set_egress_acl(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_NEIGHBOR_MISS_PACKET_ACTION:
        msg.set_neighbor_miss_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_ROUTER_INTERFACE_ATTR_V4_MCAST_ENABLE:
        msg.set_v4_mcast_enable(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_V6_MCAST_ENABLE:
        msg.set_v6_mcast_enable(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_LOOPBACK_PACKET_ACTION:
        msg.set_loopback_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_ROUTER_INTERFACE_ATTR_IS_VIRTUAL:
        msg.set_is_virtual(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_NAT_ZONE_ID:
        msg.set_nat_zone_id(attr_list[i].value.u8);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_DISABLE_DECREMENT_TTL:
        msg.set_disable_decrement_ttl(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_ADMIN_MPLS_STATE:
        msg.set_admin_mpls_state(attr_list[i].value.booldata);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_router_interface(sai_object_id_t *router_interface_id,
                                       sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateRouterInterfaceRequest req =
      convert_create_router_interface(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateRouterInterfaceResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status =
      router_interface->CreateRouterInterface(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *router_interface_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_router_interface(sai_object_id_t router_interface_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveRouterInterfaceRequest req;
  lemming::dataplane::sai::RemoveRouterInterfaceResponse resp;
  grpc::ClientContext context;
  req.set_oid(router_interface_id);

  grpc::Status status =
      router_interface->RemoveRouterInterface(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_router_interface_attribute(
    sai_object_id_t router_interface_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetRouterInterfaceAttributeRequest req;
  lemming::dataplane::sai::SetRouterInterfaceAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(router_interface_id);

  switch (attr->id) {
    case SAI_ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS:
      req.set_src_mac_address(attr->value.mac, sizeof(attr->value.mac));
      break;
    case SAI_ROUTER_INTERFACE_ATTR_ADMIN_V4_STATE:
      req.set_admin_v4_state(attr->value.booldata);
      break;
    case SAI_ROUTER_INTERFACE_ATTR_ADMIN_V6_STATE:
      req.set_admin_v6_state(attr->value.booldata);
      break;
    case SAI_ROUTER_INTERFACE_ATTR_MTU:
      req.set_mtu(attr->value.u32);
      break;
    case SAI_ROUTER_INTERFACE_ATTR_INGRESS_ACL:
      req.set_ingress_acl(attr->value.oid);
      break;
    case SAI_ROUTER_INTERFACE_ATTR_EGRESS_ACL:
      req.set_egress_acl(attr->value.oid);
      break;
    case SAI_ROUTER_INTERFACE_ATTR_NEIGHBOR_MISS_PACKET_ACTION:
      req.set_neighbor_miss_packet_action(
          static_cast<lemming::dataplane::sai::PacketAction>(attr->value.s32 +
                                                             1));
      break;
    case SAI_ROUTER_INTERFACE_ATTR_V4_MCAST_ENABLE:
      req.set_v4_mcast_enable(attr->value.booldata);
      break;
    case SAI_ROUTER_INTERFACE_ATTR_V6_MCAST_ENABLE:
      req.set_v6_mcast_enable(attr->value.booldata);
      break;
    case SAI_ROUTER_INTERFACE_ATTR_LOOPBACK_PACKET_ACTION:
      req.set_loopback_packet_action(
          static_cast<lemming::dataplane::sai::PacketAction>(attr->value.s32 +
                                                             1));
      break;
    case SAI_ROUTER_INTERFACE_ATTR_NAT_ZONE_ID:
      req.set_nat_zone_id(attr->value.u8);
      break;
    case SAI_ROUTER_INTERFACE_ATTR_DISABLE_DECREMENT_TTL:
      req.set_disable_decrement_ttl(attr->value.booldata);
      break;
    case SAI_ROUTER_INTERFACE_ATTR_ADMIN_MPLS_STATE:
      req.set_admin_mpls_state(attr->value.booldata);
      break;
  }

  grpc::Status status =
      router_interface->SetRouterInterfaceAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_router_interface_attribute(
    sai_object_id_t router_interface_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetRouterInterfaceAttributeRequest req;
  lemming::dataplane::sai::GetRouterInterfaceAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(router_interface_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::RouterInterfaceAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status =
      router_interface->GetRouterInterfaceAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ROUTER_INTERFACE_ATTR_VIRTUAL_ROUTER_ID:
        attr_list[i].value.oid = resp.attr().virtual_router_id();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_PORT_ID:
        attr_list[i].value.oid = resp.attr().port_id();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_VLAN_ID:
        attr_list[i].value.oid = resp.attr().vlan_id();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_OUTER_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().outer_vlan_id();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_INNER_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().inner_vlan_id();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_BRIDGE_ID:
        attr_list[i].value.oid = resp.attr().bridge_id();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS:
        memcpy(attr_list[i].value.mac, resp.attr().src_mac_address().data(),
               sizeof(sai_mac_t));
        break;
      case SAI_ROUTER_INTERFACE_ATTR_ADMIN_V4_STATE:
        attr_list[i].value.booldata = resp.attr().admin_v4_state();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_ADMIN_V6_STATE:
        attr_list[i].value.booldata = resp.attr().admin_v6_state();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_MTU:
        attr_list[i].value.u32 = resp.attr().mtu();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_INGRESS_ACL:
        attr_list[i].value.oid = resp.attr().ingress_acl();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_EGRESS_ACL:
        attr_list[i].value.oid = resp.attr().egress_acl();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_NEIGHBOR_MISS_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().neighbor_miss_packet_action() - 1);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_V4_MCAST_ENABLE:
        attr_list[i].value.booldata = resp.attr().v4_mcast_enable();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_V6_MCAST_ENABLE:
        attr_list[i].value.booldata = resp.attr().v6_mcast_enable();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_LOOPBACK_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().loopback_packet_action() - 1);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_IS_VIRTUAL:
        attr_list[i].value.booldata = resp.attr().is_virtual();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_NAT_ZONE_ID:
        attr_list[i].value.u8 = resp.attr().nat_zone_id();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_DISABLE_DECREMENT_TTL:
        attr_list[i].value.booldata = resp.attr().disable_decrement_ttl();
        break;
      case SAI_ROUTER_INTERFACE_ATTR_ADMIN_MPLS_STATE:
        attr_list[i].value.booldata = resp.attr().admin_mpls_state();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_router_interface_stats(sai_object_id_t router_interface_id,
                                          uint32_t number_of_counters,
                                          const sai_stat_id_t *counter_ids,
                                          uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetRouterInterfaceStatsRequest req;
  lemming::dataplane::sai::GetRouterInterfaceStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(router_interface_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(
        static_cast<lemming::dataplane::sai::RouterInterfaceStat>(
            counter_ids[i] + 1));
  }
  grpc::Status status =
      router_interface->GetRouterInterfaceStats(&context, req, &resp);
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

sai_status_t l_get_router_interface_stats_ext(
    sai_object_id_t router_interface_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids, sai_stats_mode_t mode,
    uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_router_interface_stats(sai_object_id_t router_interface_id,
                                            uint32_t number_of_counters,
                                            const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
