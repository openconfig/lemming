

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

#include "dataplane/standalone/sai/port.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/port.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_port_api_t l_port = {
    .create_port = l_create_port,
    .remove_port = l_remove_port,
    .set_port_attribute = l_set_port_attribute,
    .get_port_attribute = l_get_port_attribute,
    .get_port_stats = l_get_port_stats,
    .get_port_stats_ext = l_get_port_stats_ext,
    .clear_port_stats = l_clear_port_stats,
    .clear_port_all_stats = l_clear_port_all_stats,
    .create_port_pool = l_create_port_pool,
    .remove_port_pool = l_remove_port_pool,
    .set_port_pool_attribute = l_set_port_pool_attribute,
    .get_port_pool_attribute = l_get_port_pool_attribute,
    .get_port_pool_stats = l_get_port_pool_stats,
    .get_port_pool_stats_ext = l_get_port_pool_stats_ext,
    .clear_port_pool_stats = l_clear_port_pool_stats,
    .create_port_connector = l_create_port_connector,
    .remove_port_connector = l_remove_port_connector,
    .set_port_connector_attribute = l_set_port_connector_attribute,
    .get_port_connector_attribute = l_get_port_connector_attribute,
    .create_port_serdes = l_create_port_serdes,
    .remove_port_serdes = l_remove_port_serdes,
    .set_port_serdes_attribute = l_set_port_serdes_attribute,
    .get_port_serdes_attribute = l_get_port_serdes_attribute,
};

sai_status_t l_create_port(sai_object_id_t *port_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePortRequest req;
  lemming::dataplane::sai::CreatePortResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_ATTR_HW_LANE_LIST:
        req.mutable_hw_lane_list()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SPEED:
        req.set_speed(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_FULL_DUPLEX_MODE:
        req.set_full_duplex_mode(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_AUTO_NEG_MODE:
        req.set_auto_neg_mode(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_ADMIN_STATE:
        req.set_admin_state(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_MEDIA_TYPE:
        req.set_media_type(static_cast<lemming::dataplane::sai::PortMediaType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_ADVERTISED_SPEED:
        req.mutable_advertised_speed()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_ADVERTISED_HALF_DUPLEX_SPEED:
        req.mutable_advertised_half_duplex_speed()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_ADVERTISED_AUTO_NEG_MODE:
        req.set_advertised_auto_neg_mode(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_ADVERTISED_FLOW_CONTROL_MODE:
        req.set_advertised_flow_control_mode(
            static_cast<lemming::dataplane::sai::PortFlowControlMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_ADVERTISED_ASYMMETRIC_PAUSE_MODE:
        req.set_advertised_asymmetric_pause_mode(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_ADVERTISED_MEDIA_TYPE:
        req.set_advertised_media_type(
            static_cast<lemming::dataplane::sai::PortMediaType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_ADVERTISED_OUI_CODE:
        req.set_advertised_oui_code(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_PORT_VLAN_ID:
        req.set_port_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_PORT_ATTR_DEFAULT_VLAN_PRIORITY:
        req.set_default_vlan_priority(attr_list[i].value.u8);
        break;
      case SAI_PORT_ATTR_DROP_UNTAGGED:
        req.set_drop_untagged(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_DROP_TAGGED:
        req.set_drop_tagged(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_INTERNAL_LOOPBACK_MODE:
        req.set_internal_loopback_mode(
            static_cast<lemming::dataplane::sai::PortInternalLoopbackMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_USE_EXTENDED_FEC:
        req.set_use_extended_fec(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_FEC_MODE:
        req.set_fec_mode(static_cast<lemming::dataplane::sai::PortFecMode>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_FEC_MODE_EXTENDED:
        req.set_fec_mode_extended(
            static_cast<lemming::dataplane::sai::PortFecModeExtended>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_UPDATE_DSCP:
        req.set_update_dscp(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_MTU:
        req.set_mtu(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_FLOOD_STORM_CONTROL_POLICER_ID:
        req.set_flood_storm_control_policer_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_BROADCAST_STORM_CONTROL_POLICER_ID:
        req.set_broadcast_storm_control_policer_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_MULTICAST_STORM_CONTROL_POLICER_ID:
        req.set_multicast_storm_control_policer_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_MODE:
        req.set_global_flow_control_mode(
            static_cast<lemming::dataplane::sai::PortFlowControlMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_INGRESS_ACL:
        req.set_ingress_acl(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_EGRESS_ACL:
        req.set_egress_acl(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_INGRESS_MACSEC_ACL:
        req.set_ingress_macsec_acl(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_EGRESS_MACSEC_ACL:
        req.set_egress_macsec_acl(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_INGRESS_MIRROR_SESSION:
        req.mutable_ingress_mirror_session()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_EGRESS_MIRROR_SESSION:
        req.mutable_egress_mirror_session()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
        req.set_ingress_samplepacket_enable(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_EGRESS_SAMPLEPACKET_ENABLE:
        req.set_egress_samplepacket_enable(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_INGRESS_SAMPLE_MIRROR_SESSION:
        req.mutable_ingress_sample_mirror_session()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_EGRESS_SAMPLE_MIRROR_SESSION:
        req.mutable_egress_sample_mirror_session()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_POLICER_ID:
        req.set_policer_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_DEFAULT_TC:
        req.set_qos_default_tc(attr_list[i].value.u8);
        break;
      case SAI_PORT_ATTR_QOS_DOT1P_TO_TC_MAP:
        req.set_qos_dot1p_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_DOT1P_TO_COLOR_MAP:
        req.set_qos_dot1p_to_color_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_DSCP_TO_TC_MAP:
        req.set_qos_dscp_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_DSCP_TO_COLOR_MAP:
        req.set_qos_dscp_to_color_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
        req.set_qos_tc_to_queue_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
        req.set_qos_tc_and_color_to_dot1p_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
        req.set_qos_tc_and_color_to_dscp_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_TC_TO_PRIORITY_GROUP_MAP:
        req.set_qos_tc_to_priority_group_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_PRIORITY_GROUP_MAP:
        req.set_qos_pfc_priority_to_priority_group_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_QUEUE_MAP:
        req.set_qos_pfc_priority_to_queue_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_SCHEDULER_PROFILE_ID:
        req.set_qos_scheduler_profile_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_INGRESS_BUFFER_PROFILE_LIST:
        req.mutable_qos_ingress_buffer_profile_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_QOS_EGRESS_BUFFER_PROFILE_LIST:
        req.mutable_qos_egress_buffer_profile_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_MODE:
        req.set_priority_flow_control_mode(
            static_cast<lemming::dataplane::sai::PortPriorityFlowControlMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL:
        req.set_priority_flow_control(attr_list[i].value.u8);
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_RX:
        req.set_priority_flow_control_rx(attr_list[i].value.u8);
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_TX:
        req.set_priority_flow_control_tx(attr_list[i].value.u8);
        break;
      case SAI_PORT_ATTR_META_DATA:
        req.set_meta_data(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_EGRESS_BLOCK_PORT_LIST:
        req.mutable_egress_block_port_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_HW_PROFILE_ID:
        req.set_hw_profile_id(attr_list[i].value.u64);
        break;
      case SAI_PORT_ATTR_EEE_ENABLE:
        req.set_eee_enable(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_EEE_IDLE_TIME:
        req.set_eee_idle_time(attr_list[i].value.u16);
        break;
      case SAI_PORT_ATTR_EEE_WAKE_TIME:
        req.set_eee_wake_time(attr_list[i].value.u16);
        break;
      case SAI_PORT_ATTR_ISOLATION_GROUP:
        req.set_isolation_group(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_PKT_TX_ENABLE:
        req.set_pkt_tx_enable(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_TAM_OBJECT:
        req.mutable_tam_object()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_SERDES_PREEMPHASIS:
        req.mutable_serdes_preemphasis()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SERDES_IDRIVER:
        req.mutable_serdes_idriver()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SERDES_IPREDRIVER:
        req.mutable_serdes_ipredriver()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_LINK_TRAINING_ENABLE:
        req.set_link_training_enable(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_PTP_MODE:
        req.set_ptp_mode(static_cast<lemming::dataplane::sai::PortPtpMode>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_INTERFACE_TYPE:
        req.set_interface_type(
            static_cast<lemming::dataplane::sai::PortInterfaceType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_REFERENCE_CLOCK:
        req.set_reference_clock(attr_list[i].value.u64);
        break;
      case SAI_PORT_ATTR_PRBS_POLYNOMIAL:
        req.set_prbs_polynomial(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_PRBS_CONFIG:
        req.set_prbs_config(
            static_cast<lemming::dataplane::sai::PortPrbsConfig>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_DISABLE_DECREMENT_TTL:
        req.set_disable_decrement_ttl(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
        req.set_qos_mpls_exp_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
        req.set_qos_mpls_exp_to_color_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
        req.set_qos_tc_and_color_to_mpls_exp_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_TPID:
        req.set_tpid(attr_list[i].value.u16);
        break;
      case SAI_PORT_ATTR_AUTO_NEG_FEC_MODE_OVERRIDE:
        req.set_auto_neg_fec_mode_override(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_LOOPBACK_MODE:
        req.set_loopback_mode(
            static_cast<lemming::dataplane::sai::PortLoopbackMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_MDIX_MODE_CONFIG:
        req.set_mdix_mode_config(
            static_cast<lemming::dataplane::sai::PortMdixModeConfig>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_AUTO_NEG_CONFIG_MODE:
        req.set_auto_neg_config_mode(
            static_cast<lemming::dataplane::sai::PortAutoNegConfigMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_1000X_SGMII_SLAVE_AUTODETECT:
        req.set__1000x_sgmii_slave_autodetect(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_MODULE_TYPE:
        req.set_module_type(
            static_cast<lemming::dataplane::sai::PortModuleType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_DUAL_MEDIA:
        req.set_dual_media(static_cast<lemming::dataplane::sai::PortDualMedia>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_PORT_ATTR_IPG:
        req.set_ipg(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_FORWARD:
        req.set_global_flow_control_forward(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_FORWARD:
        req.set_priority_flow_control_forward(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
        req.set_qos_dscp_to_forwarding_class_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
        req.set_qos_mpls_exp_to_forwarding_class_map(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = port->CreatePort(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *port_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_PORT, port_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_port(sai_object_id_t port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_PORT, port_id);
}

sai_status_t l_set_port_attribute(sai_object_id_t port_id,
                                  const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_PORT, port_id, attr);
}

sai_status_t l_get_port_attribute(sai_object_id_t port_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_PORT, port_id, attr_count,
                                   attr_list);
}

sai_status_t l_get_port_stats(sai_object_id_t port_id,
                              uint32_t number_of_counters,
                              const sai_stat_id_t *counter_ids,
                              uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_PORT, port_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_port_stats_ext(sai_object_id_t port_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids,
                                  sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_PORT, port_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_port_stats(sai_object_id_t port_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_PORT, port_id,
                                 number_of_counters, counter_ids);
}

sai_status_t l_clear_port_all_stats(sai_object_id_t port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_port_pool(sai_object_id_t *port_pool_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePortPoolRequest req;
  lemming::dataplane::sai::CreatePortPoolResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_POOL_ATTR_PORT_ID:
        req.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_POOL_ATTR_BUFFER_POOL_ID:
        req.set_buffer_pool_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_POOL_ATTR_QOS_WRED_PROFILE_ID:
        req.set_qos_wred_profile_id(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = port->CreatePortPool(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *port_pool_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_port_pool(sai_object_id_t port_pool_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id);
}

sai_status_t l_set_port_pool_attribute(sai_object_id_t port_pool_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id,
                                   attr);
}

sai_status_t l_get_port_pool_attribute(sai_object_id_t port_pool_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_port_pool_stats(sai_object_id_t port_pool_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids,
                                   uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_port_pool_stats_ext(sai_object_id_t port_pool_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids,
                                       sai_stats_mode_t mode,
                                       uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_port_pool_stats(sai_object_id_t port_pool_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_PORT_POOL, port_pool_id,
                                 number_of_counters, counter_ids);
}

sai_status_t l_create_port_connector(sai_object_id_t *port_connector_id,
                                     sai_object_id_t switch_id,
                                     uint32_t attr_count,
                                     const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePortConnectorRequest req;
  lemming::dataplane::sai::CreatePortConnectorResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_PORT_ID:
        req.set_system_side_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_CONNECTOR_ATTR_LINE_SIDE_PORT_ID:
        req.set_line_side_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_FAILOVER_PORT_ID:
        req.set_system_side_failover_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_CONNECTOR_ATTR_LINE_SIDE_FAILOVER_PORT_ID:
        req.set_line_side_failover_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_CONNECTOR_ATTR_FAILOVER_MODE:
        req.set_failover_mode(
            static_cast<lemming::dataplane::sai::PortConnectorFailoverMode>(
                attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = port->CreatePortConnector(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *port_connector_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_PORT_CONNECTOR, port_connector_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_port_connector(sai_object_id_t port_connector_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_PORT_CONNECTOR, port_connector_id);
}

sai_status_t l_set_port_connector_attribute(sai_object_id_t port_connector_id,
                                            const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_PORT_CONNECTOR,
                                   port_connector_id, attr);
}

sai_status_t l_get_port_connector_attribute(sai_object_id_t port_connector_id,
                                            uint32_t attr_count,
                                            sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_PORT_CONNECTOR,
                                   port_connector_id, attr_count, attr_list);
}

sai_status_t l_create_port_serdes(sai_object_id_t *port_serdes_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePortSerdesRequest req;
  lemming::dataplane::sai::CreatePortSerdesResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_SERDES_ATTR_PORT_ID:
        req.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_SERDES_ATTR_PREEMPHASIS:
        req.mutable_preemphasis()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_IDRIVER:
        req.mutable_idriver()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_IPREDRIVER:
        req.mutable_ipredriver()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_PRE1:
        req.mutable_tx_fir_pre1()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_PRE2:
        req.mutable_tx_fir_pre2()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_PRE3:
        req.mutable_tx_fir_pre3()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_MAIN:
        req.mutable_tx_fir_main()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_POST1:
        req.mutable_tx_fir_post1()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_POST2:
        req.mutable_tx_fir_post2()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_POST3:
        req.mutable_tx_fir_post3()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_ATTN:
        req.mutable_tx_fir_attn()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
    }
  }
  grpc::Status status = port->CreatePortSerdes(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *port_serdes_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_PORT_SERDES, port_serdes_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_port_serdes(sai_object_id_t port_serdes_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_PORT_SERDES, port_serdes_id);
}

sai_status_t l_set_port_serdes_attribute(sai_object_id_t port_serdes_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_PORT_SERDES, port_serdes_id,
                                   attr);
}

sai_status_t l_get_port_serdes_attribute(sai_object_id_t port_serdes_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_PORT_SERDES, port_serdes_id,
                                   attr_count, attr_list);
}
