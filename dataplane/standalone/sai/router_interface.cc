

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

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/router_interface.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_router_interface_api_t l_router_interface = {
    .create_router_interface = l_create_router_interface,
    .remove_router_interface = l_remove_router_interface,
    .set_router_interface_attribute = l_set_router_interface_attribute,
    .get_router_interface_attribute = l_get_router_interface_attribute,
    .get_router_interface_stats = l_get_router_interface_stats,
    .get_router_interface_stats_ext = l_get_router_interface_stats_ext,
    .clear_router_interface_stats = l_clear_router_interface_stats,
};

sai_status_t l_create_router_interface(sai_object_id_t *router_interface_id,
                                       sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateRouterInterfaceRequest req;
  lemming::dataplane::sai::CreateRouterInterfaceResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ROUTER_INTERFACE_ATTR_VIRTUAL_ROUTER_ID:
        req.set_virtual_router_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::RouterInterfaceType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_ROUTER_INTERFACE_ATTR_PORT_ID:
        req.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_VLAN_ID:
        req.set_vlan_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_OUTER_VLAN_ID:
        req.set_outer_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_INNER_VLAN_ID:
        req.set_inner_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_BRIDGE_ID:
        req.set_bridge_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS:
        req.set_src_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_ROUTER_INTERFACE_ATTR_ADMIN_V4_STATE:
        req.set_admin_v4_state(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_ADMIN_V6_STATE:
        req.set_admin_v6_state(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_MTU:
        req.set_mtu(attr_list[i].value.u32);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_INGRESS_ACL:
        req.set_ingress_acl(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_EGRESS_ACL:
        req.set_egress_acl(attr_list[i].value.oid);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_NEIGHBOR_MISS_PACKET_ACTION:
        req.set_neighbor_miss_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_ROUTER_INTERFACE_ATTR_V4_MCAST_ENABLE:
        req.set_v4_mcast_enable(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_V6_MCAST_ENABLE:
        req.set_v6_mcast_enable(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_LOOPBACK_PACKET_ACTION:
        req.set_loopback_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_ROUTER_INTERFACE_ATTR_IS_VIRTUAL:
        req.set_is_virtual(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_NAT_ZONE_ID:
        req.set_nat_zone_id(attr_list[i].value.u8);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_DISABLE_DECREMENT_TTL:
        req.set_disable_decrement_ttl(attr_list[i].value.booldata);
        break;
      case SAI_ROUTER_INTERFACE_ATTR_ADMIN_MPLS_STATE:
        req.set_admin_mpls_state(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status =
      router_interface->CreateRouterInterface(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *router_interface_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_ROUTER_INTERFACE,
                            router_interface_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_router_interface(sai_object_id_t router_interface_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_ROUTER_INTERFACE,
                            router_interface_id);
}

sai_status_t l_set_router_interface_attribute(
    sai_object_id_t router_interface_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_ROUTER_INTERFACE,
                                   router_interface_id, attr);
}

sai_status_t l_get_router_interface_attribute(
    sai_object_id_t router_interface_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_ROUTER_INTERFACE,
                                   router_interface_id, attr_count, attr_list);
}

sai_status_t l_get_router_interface_stats(sai_object_id_t router_interface_id,
                                          uint32_t number_of_counters,
                                          const sai_stat_id_t *counter_ids,
                                          uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_ROUTER_INTERFACE,
                               router_interface_id, number_of_counters,
                               counter_ids, counters);
}

sai_status_t l_get_router_interface_stats_ext(
    sai_object_id_t router_interface_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids, sai_stats_mode_t mode,
    uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_ROUTER_INTERFACE,
                                   router_interface_id, number_of_counters,
                                   counter_ids, mode, counters);
}

sai_status_t l_clear_router_interface_stats(sai_object_id_t router_interface_id,
                                            uint32_t number_of_counters,
                                            const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_ROUTER_INTERFACE,
                                 router_interface_id, number_of_counters,
                                 counter_ids);
}
