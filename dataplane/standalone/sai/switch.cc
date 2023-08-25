

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

#include "dataplane/standalone/sai/switch.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/switch.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_switch_api_t l_switch = {
    .create_switch = l_create_switch,
    .remove_switch = l_remove_switch,
    .set_switch_attribute = l_set_switch_attribute,
    .get_switch_attribute = l_get_switch_attribute,
    .get_switch_stats = l_get_switch_stats,
    .get_switch_stats_ext = l_get_switch_stats_ext,
    .clear_switch_stats = l_clear_switch_stats,
    .switch_mdio_read = l_switch_mdio_read,
    .switch_mdio_write = l_switch_mdio_write,
    .create_switch_tunnel = l_create_switch_tunnel,
    .remove_switch_tunnel = l_remove_switch_tunnel,
    .set_switch_tunnel_attribute = l_set_switch_tunnel_attribute,
    .get_switch_tunnel_attribute = l_get_switch_tunnel_attribute,
};

sai_status_t l_create_switch(sai_object_id_t *switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateSwitchRequest req;
  lemming::dataplane::sai::CreateSwitchResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SWITCH_ATTR_INGRESS_ACL:
        req.set_ingress_acl(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_EGRESS_ACL:
        req.set_egress_acl(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_RESTART_WARM:
        req.set_restart_warm(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_WARM_RECOVER:
        req.set_warm_recover(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_SWITCHING_MODE:
        req.set_switching_mode(
            static_cast<lemming::dataplane::sai::SwitchSwitchingMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_BCAST_CPU_FLOOD_ENABLE:
        req.set_bcast_cpu_flood_enable(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_MCAST_CPU_FLOOD_ENABLE:
        req.set_mcast_cpu_flood_enable(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_SRC_MAC_ADDRESS:
        req.set_src_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_SWITCH_ATTR_MAX_LEARNED_ADDRESSES:
        req.set_max_learned_addresses(attr_list[i].value.u32);
        break;
      case SAI_SWITCH_ATTR_FDB_AGING_TIME:
        req.set_fdb_aging_time(attr_list[i].value.u32);
        break;
      case SAI_SWITCH_ATTR_FDB_UNICAST_MISS_PACKET_ACTION:
        req.set_fdb_unicast_miss_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_FDB_BROADCAST_MISS_PACKET_ACTION:
        req.set_fdb_broadcast_miss_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_FDB_MULTICAST_MISS_PACKET_ACTION:
        req.set_fdb_multicast_miss_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_ALGORITHM:
        req.set_ecmp_default_hash_algorithm(
            static_cast<lemming::dataplane::sai::HashAlgorithm>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_SEED:
        req.set_ecmp_default_hash_seed(attr_list[i].value.u32);
        break;
      case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_OFFSET:
        req.set_ecmp_default_hash_offset(attr_list[i].value.u8);
        break;
      case SAI_SWITCH_ATTR_ECMP_DEFAULT_SYMMETRIC_HASH:
        req.set_ecmp_default_symmetric_hash(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_ECMP_HASH_IPV4:
        req.set_ecmp_hash_ipv4(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_ECMP_HASH_IPV4_IN_IPV4:
        req.set_ecmp_hash_ipv4_in_ipv4(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_ECMP_HASH_IPV6:
        req.set_ecmp_hash_ipv6(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM:
        req.set_lag_default_hash_algorithm(
            static_cast<lemming::dataplane::sai::HashAlgorithm>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_SEED:
        req.set_lag_default_hash_seed(attr_list[i].value.u32);
        break;
      case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_OFFSET:
        req.set_lag_default_hash_offset(attr_list[i].value.u8);
        break;
      case SAI_SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH:
        req.set_lag_default_symmetric_hash(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_LAG_HASH_IPV4:
        req.set_lag_hash_ipv4(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_LAG_HASH_IPV4_IN_IPV4:
        req.set_lag_hash_ipv4_in_ipv4(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_LAG_HASH_IPV6:
        req.set_lag_hash_ipv6(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_COUNTER_REFRESH_INTERVAL:
        req.set_counter_refresh_interval(attr_list[i].value.u32);
        break;
      case SAI_SWITCH_ATTR_QOS_DEFAULT_TC:
        req.set_qos_default_tc(attr_list[i].value.u8);
        break;
      case SAI_SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP:
        req.set_qos_dot1p_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP:
        req.set_qos_dot1p_to_color_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_QOS_DSCP_TO_TC_MAP:
        req.set_qos_dscp_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_QOS_DSCP_TO_COLOR_MAP:
        req.set_qos_dscp_to_color_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP:
        req.set_qos_tc_to_queue_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
        req.set_qos_tc_and_color_to_dot1p_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
        req.set_qos_tc_and_color_to_dscp_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_SWITCH_SHELL_ENABLE:
        req.set_switch_shell_enable(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_SWITCH_PROFILE_ID:
        req.set_switch_profile_id(attr_list[i].value.u32);
        break;
      case SAI_SWITCH_ATTR_SWITCH_HARDWARE_INFO:
        req.mutable_switch_hardware_info()->Add(
            attr_list[i].value.s8list.list,
            attr_list[i].value.s8list.list + attr_list[i].value.s8list.count);
        break;
      case SAI_SWITCH_ATTR_FIRMWARE_PATH_NAME:
        req.mutable_firmware_path_name()->Add(
            attr_list[i].value.s8list.list,
            attr_list[i].value.s8list.list + attr_list[i].value.s8list.count);
        break;
      case SAI_SWITCH_ATTR_INIT_SWITCH:
        req.set_init_switch(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_FAST_API_ENABLE:
        req.set_fast_api_enable(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_MIRROR_TC:
        req.set_mirror_tc(attr_list[i].value.u8);
        break;
      case SAI_SWITCH_ATTR_PFC_DLR_PACKET_ACTION:
        req.set_pfc_dlr_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_TPID_OUTER_VLAN:
        req.set_tpid_outer_vlan(attr_list[i].value.u16);
        break;
      case SAI_SWITCH_ATTR_TPID_INNER_VLAN:
        req.set_tpid_inner_vlan(attr_list[i].value.u16);
        break;
      case SAI_SWITCH_ATTR_CRC_CHECK_ENABLE:
        req.set_crc_check_enable(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_CRC_RECALCULATION_ENABLE:
        req.set_crc_recalculation_enable(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_ECN_ECT_THRESHOLD_ENABLE:
        req.set_ecn_ect_threshold_enable(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_VXLAN_DEFAULT_ROUTER_MAC:
        req.set_vxlan_default_router_mac(attr_list[i].value.mac,
                                         sizeof(attr_list[i].value.mac));
        break;
      case SAI_SWITCH_ATTR_VXLAN_DEFAULT_PORT:
        req.set_vxlan_default_port(attr_list[i].value.u16);
        break;
      case SAI_SWITCH_ATTR_UNINIT_DATA_PLANE_ON_REMOVAL:
        req.set_uninit_data_plane_on_removal(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_TAM_OBJECT_ID:
        req.mutable_tam_object_id()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_SWITCH_ATTR_PRE_SHUTDOWN:
        req.set_pre_shutdown(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID:
        req.set_nat_zone_counter_object_id(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_NAT_ENABLE:
        req.set_nat_enable(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_HARDWARE_ACCESS_BUS:
        req.set_hardware_access_bus(
            static_cast<lemming::dataplane::sai::SwitchHardwareAccessBus>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_PLATFROM_CONTEXT:
        req.set_platfrom_context(attr_list[i].value.u64);
        break;
      case SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_BROADCAST:
        req.set_firmware_download_broadcast(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_FIRMWARE_LOAD_METHOD:
        req.set_firmware_load_method(
            static_cast<lemming::dataplane::sai::SwitchFirmwareLoadMethod>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_FIRMWARE_LOAD_TYPE:
        req.set_firmware_load_type(
            static_cast<lemming::dataplane::sai::SwitchFirmwareLoadType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_EXECUTE:
        req.set_firmware_download_execute(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_FIRMWARE_BROADCAST_STOP:
        req.set_firmware_broadcast_stop(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_FIRMWARE_VERIFY_AND_INIT_SWITCH:
        req.set_firmware_verify_and_init_switch(attr_list[i].value.booldata);
        break;
      case SAI_SWITCH_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::SwitchType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_MACSEC_OBJECT_LIST:
        req.mutable_macsec_object_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
        req.set_qos_mpls_exp_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
        req.set_qos_mpls_exp_to_color_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
        req.set_qos_tc_and_color_to_mpls_exp_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_SWITCH_ID:
        req.set_switch_id(attr_list[i].value.u32);
        break;
      case SAI_SWITCH_ATTR_MAX_SYSTEM_CORES:
        req.set_max_system_cores(attr_list[i].value.u32);
        break;
      case SAI_SWITCH_ATTR_FAILOVER_CONFIG_MODE:
        req.set_failover_config_mode(
            static_cast<lemming::dataplane::sai::SwitchFailoverConfigMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_ATTR_TUNNEL_OBJECTS_LIST:
        req.mutable_tunnel_objects_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_SWITCH_ATTR_PRE_INGRESS_ACL:
        req.set_pre_ingress_acl(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_SLAVE_MDIO_ADDR_LIST:
        req.mutable_slave_mdio_addr_list()->Add(
            attr_list[i].value.u8list.list,
            attr_list[i].value.u8list.list + attr_list[i].value.u8list.count);
        break;
      case SAI_SWITCH_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
        req.set_qos_dscp_to_forwarding_class_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
        req.set_qos_mpls_exp_to_forwarding_class_map(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_IPSEC_OBJECT_ID:
        req.set_ipsec_object_id(attr_list[i].value.oid);
        break;
      case SAI_SWITCH_ATTR_IPSEC_SA_TAG_TPID:
        req.set_ipsec_sa_tag_tpid(attr_list[i].value.u16);
        break;
    }
  }
  grpc::Status status = switch_->CreateSwitch(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *switch_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_SWITCH, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_switch(sai_object_id_t switch_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_SWITCH, switch_id);
}

sai_status_t l_set_switch_attribute(sai_object_id_t switch_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_SWITCH, switch_id, attr);
}

sai_status_t l_get_switch_attribute(sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_SWITCH, switch_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_switch_stats(sai_object_id_t switch_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids,
                                uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_SWITCH, switch_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_switch_stats_ext(sai_object_id_t switch_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t *counter_ids,
                                    sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_SWITCH, switch_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_switch_stats(sai_object_id_t switch_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_SWITCH, switch_id,
                                 number_of_counters, counter_ids);
}

sai_status_t l_switch_mdio_read(sai_object_id_t switch_id, uint32_t device_addr,
                                uint32_t start_reg_addr,
                                uint32_t number_of_registers,
                                uint32_t *reg_val) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_switch_mdio_write(sai_object_id_t switch_id,
                                 uint32_t device_addr, uint32_t start_reg_addr,
                                 uint32_t number_of_registers,
                                 const uint32_t *reg_val) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_switch_tunnel(sai_object_id_t *switch_tunnel_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateSwitchTunnelRequest req;
  lemming::dataplane::sai::CreateSwitchTunnelResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_TYPE:
        req.set_tunnel_type(static_cast<lemming::dataplane::sai::TunnelType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
        req.set_loopback_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_ENCAP_ECN_MODE:
        req.set_tunnel_encap_ecn_mode(
            static_cast<lemming::dataplane::sai::TunnelEncapEcnMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_TUNNEL_ATTR_ENCAP_MAPPERS:
        req.mutable_encap_mappers()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_DECAP_ECN_MODE:
        req.set_tunnel_decap_ecn_mode(
            static_cast<lemming::dataplane::sai::TunnelDecapEcnMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_TUNNEL_ATTR_DECAP_MAPPERS:
        req.mutable_decap_mappers()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_VXLAN_UDP_SPORT_MODE:
        req.set_tunnel_vxlan_udp_sport_mode(
            static_cast<lemming::dataplane::sai::TunnelVxlanUdpSportMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT:
        req.set_vxlan_udp_sport(attr_list[i].value.u16);
        break;
      case SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
        req.set_vxlan_udp_sport_mask(attr_list[i].value.u8);
        break;
    }
  }
  grpc::Status status = switch_->CreateSwitchTunnel(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *switch_tunnel_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_SWITCH_TUNNEL, switch_tunnel_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_switch_tunnel(sai_object_id_t switch_tunnel_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_SWITCH_TUNNEL, switch_tunnel_id);
}

sai_status_t l_set_switch_tunnel_attribute(sai_object_id_t switch_tunnel_id,
                                           const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_SWITCH_TUNNEL,
                                   switch_tunnel_id, attr);
}

sai_status_t l_get_switch_tunnel_attribute(sai_object_id_t switch_tunnel_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_SWITCH_TUNNEL,
                                   switch_tunnel_id, attr_count, attr_list);
}
