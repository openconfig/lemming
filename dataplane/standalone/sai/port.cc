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

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/port.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_port_api_t l_port = {
    .create_port = l_create_port,
    .remove_port = l_remove_port,
    .set_port_attribute = l_set_port_attribute,
    .get_port_attribute = l_get_port_attribute,
    .get_port_stats = l_get_port_stats,
    .get_port_stats_ext = l_get_port_stats_ext,
    .clear_port_stats = l_clear_port_stats,
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
    .create_ports = l_create_ports,
    .remove_ports = l_remove_ports,
    .set_ports_attribute = l_set_ports_attribute,
    .get_ports_attribute = l_get_ports_attribute,
};

lemming::dataplane::sai::CreatePortRequest convert_create_port(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreatePortRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_ATTR_HW_LANE_LIST:
        msg.mutable_hw_lane_list()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SPEED:
        msg.set_speed(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_FULL_DUPLEX_MODE:
        msg.set_full_duplex_mode(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_AUTO_NEG_MODE:
        msg.set_auto_neg_mode(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_ADMIN_STATE:
        msg.set_admin_state(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_MEDIA_TYPE:
        msg.set_media_type(
            convert_sai_port_media_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_ADVERTISED_SPEED:
        msg.mutable_advertised_speed()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_ADVERTISED_HALF_DUPLEX_SPEED:
        msg.mutable_advertised_half_duplex_speed()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_ADVERTISED_AUTO_NEG_MODE:
        msg.set_advertised_auto_neg_mode(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_ADVERTISED_FLOW_CONTROL_MODE:
        msg.set_advertised_flow_control_mode(
            convert_sai_port_flow_control_mode_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_ADVERTISED_ASYMMETRIC_PAUSE_MODE:
        msg.set_advertised_asymmetric_pause_mode(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_ADVERTISED_MEDIA_TYPE:
        msg.set_advertised_media_type(
            convert_sai_port_media_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_ADVERTISED_OUI_CODE:
        msg.set_advertised_oui_code(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_PORT_VLAN_ID:
        msg.set_port_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_PORT_ATTR_DEFAULT_VLAN_PRIORITY:
        msg.set_default_vlan_priority(attr_list[i].value.u8);
        break;
      case SAI_PORT_ATTR_DROP_UNTAGGED:
        msg.set_drop_untagged(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_DROP_TAGGED:
        msg.set_drop_tagged(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_INTERNAL_LOOPBACK_MODE:
        msg.set_internal_loopback_mode(
            convert_sai_port_internal_loopback_mode_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_USE_EXTENDED_FEC:
        msg.set_use_extended_fec(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_FEC_MODE:
        msg.set_fec_mode(
            convert_sai_port_fec_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_FEC_MODE_EXTENDED:
        msg.set_fec_mode_extended(convert_sai_port_fec_mode_extended_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_UPDATE_DSCP:
        msg.set_update_dscp(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_MTU:
        msg.set_mtu(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_FLOOD_STORM_CONTROL_POLICER_ID:
        msg.set_flood_storm_control_policer_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_BROADCAST_STORM_CONTROL_POLICER_ID:
        msg.set_broadcast_storm_control_policer_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_MULTICAST_STORM_CONTROL_POLICER_ID:
        msg.set_multicast_storm_control_policer_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_MODE:
        msg.set_global_flow_control_mode(
            convert_sai_port_flow_control_mode_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_INGRESS_ACL:
        msg.set_ingress_acl(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_EGRESS_ACL:
        msg.set_egress_acl(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_INGRESS_MACSEC_ACL:
        msg.set_ingress_macsec_acl(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_EGRESS_MACSEC_ACL:
        msg.set_egress_macsec_acl(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_INGRESS_MIRROR_SESSION:
        msg.mutable_ingress_mirror_session()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_EGRESS_MIRROR_SESSION:
        msg.mutable_egress_mirror_session()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
        msg.set_ingress_samplepacket_enable(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_EGRESS_SAMPLEPACKET_ENABLE:
        msg.set_egress_samplepacket_enable(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_INGRESS_SAMPLE_MIRROR_SESSION:
        msg.mutable_ingress_sample_mirror_session()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_EGRESS_SAMPLE_MIRROR_SESSION:
        msg.mutable_egress_sample_mirror_session()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_POLICER_ID:
        msg.set_policer_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_DEFAULT_TC:
        msg.set_qos_default_tc(attr_list[i].value.u8);
        break;
      case SAI_PORT_ATTR_QOS_DOT1P_TO_TC_MAP:
        msg.set_qos_dot1p_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_DOT1P_TO_COLOR_MAP:
        msg.set_qos_dot1p_to_color_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_DSCP_TO_TC_MAP:
        msg.set_qos_dscp_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_DSCP_TO_COLOR_MAP:
        msg.set_qos_dscp_to_color_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
        msg.set_qos_tc_to_queue_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
        msg.set_qos_tc_and_color_to_dot1p_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
        msg.set_qos_tc_and_color_to_dscp_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_TC_TO_PRIORITY_GROUP_MAP:
        msg.set_qos_tc_to_priority_group_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_PRIORITY_GROUP_MAP:
        msg.set_qos_pfc_priority_to_priority_group_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_QUEUE_MAP:
        msg.set_qos_pfc_priority_to_queue_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_SCHEDULER_PROFILE_ID:
        msg.set_qos_scheduler_profile_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_INGRESS_BUFFER_PROFILE_LIST:
        msg.mutable_qos_ingress_buffer_profile_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_QOS_EGRESS_BUFFER_PROFILE_LIST:
        msg.mutable_qos_egress_buffer_profile_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_MODE:
        msg.set_priority_flow_control_mode(
            convert_sai_port_priority_flow_control_mode_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL:
        msg.set_priority_flow_control(attr_list[i].value.u8);
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_RX:
        msg.set_priority_flow_control_rx(attr_list[i].value.u8);
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_TX:
        msg.set_priority_flow_control_tx(attr_list[i].value.u8);
        break;
      case SAI_PORT_ATTR_META_DATA:
        msg.set_meta_data(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_EGRESS_BLOCK_PORT_LIST:
        msg.mutable_egress_block_port_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_HW_PROFILE_ID:
        msg.set_hw_profile_id(attr_list[i].value.u64);
        break;
      case SAI_PORT_ATTR_EEE_ENABLE:
        msg.set_eee_enable(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_EEE_IDLE_TIME:
        msg.set_eee_idle_time(attr_list[i].value.u16);
        break;
      case SAI_PORT_ATTR_EEE_WAKE_TIME:
        msg.set_eee_wake_time(attr_list[i].value.u16);
        break;
      case SAI_PORT_ATTR_ISOLATION_GROUP:
        msg.set_isolation_group(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_PKT_TX_ENABLE:
        msg.set_pkt_tx_enable(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_TAM_OBJECT:
        msg.mutable_tam_object()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_SERDES_PREEMPHASIS:
        msg.mutable_serdes_preemphasis()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SERDES_IDRIVER:
        msg.mutable_serdes_idriver()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SERDES_IPREDRIVER:
        msg.mutable_serdes_ipredriver()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_LINK_TRAINING_ENABLE:
        msg.set_link_training_enable(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_PTP_MODE:
        msg.set_ptp_mode(
            convert_sai_port_ptp_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_INTERFACE_TYPE:
        msg.set_interface_type(
            convert_sai_port_interface_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_REFERENCE_CLOCK:
        msg.set_reference_clock(attr_list[i].value.u64);
        break;
      case SAI_PORT_ATTR_PRBS_POLYNOMIAL:
        msg.set_prbs_polynomial(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_PRBS_CONFIG:
        msg.set_prbs_config(
            convert_sai_port_prbs_config_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_DISABLE_DECREMENT_TTL:
        msg.set_disable_decrement_ttl(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
        msg.set_qos_mpls_exp_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
        msg.set_qos_mpls_exp_to_color_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
        msg.set_qos_tc_and_color_to_mpls_exp_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_TPID:
        msg.set_tpid(attr_list[i].value.u16);
        break;
      case SAI_PORT_ATTR_AUTO_NEG_FEC_MODE_OVERRIDE:
        msg.set_auto_neg_fec_mode_override(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_LOOPBACK_MODE:
        msg.set_loopback_mode(
            convert_sai_port_loopback_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_MDIX_MODE_CONFIG:
        msg.set_mdix_mode_config(convert_sai_port_mdix_mode_config_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_AUTO_NEG_CONFIG_MODE:
        msg.set_auto_neg_config_mode(
            convert_sai_port_auto_neg_config_mode_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_1000X_SGMII_SLAVE_AUTODETECT:
        msg.set__1000x_sgmii_slave_autodetect(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_MODULE_TYPE:
        msg.set_module_type(
            convert_sai_port_module_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_DUAL_MEDIA:
        msg.set_dual_media(
            convert_sai_port_dual_media_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_PORT_ATTR_IPG:
        msg.set_ipg(attr_list[i].value.u32);
        break;
      case SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_FORWARD:
        msg.set_global_flow_control_forward(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_FORWARD:
        msg.set_priority_flow_control_forward(attr_list[i].value.booldata);
        break;
      case SAI_PORT_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
        msg.set_qos_dscp_to_forwarding_class_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
        msg.set_qos_mpls_exp_to_forwarding_class_map(attr_list[i].value.oid);
        break;
      case SAI_PORT_ATTR_FABRIC_ISOLATE:
        msg.set_fabric_isolate(attr_list[i].value.booldata);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreatePortPoolRequest convert_create_port_pool(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreatePortPoolRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_POOL_ATTR_PORT_ID:
        msg.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_POOL_ATTR_BUFFER_POOL_ID:
        msg.set_buffer_pool_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_POOL_ATTR_QOS_WRED_PROFILE_ID:
        msg.set_qos_wred_profile_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreatePortConnectorRequest
convert_create_port_connector(sai_object_id_t switch_id, uint32_t attr_count,
                              const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreatePortConnectorRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_PORT_ID:
        msg.set_system_side_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_CONNECTOR_ATTR_LINE_SIDE_PORT_ID:
        msg.set_line_side_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_FAILOVER_PORT_ID:
        msg.set_system_side_failover_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_CONNECTOR_ATTR_LINE_SIDE_FAILOVER_PORT_ID:
        msg.set_line_side_failover_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_CONNECTOR_ATTR_FAILOVER_MODE:
        msg.set_failover_mode(
            convert_sai_port_connector_failover_mode_t_to_proto(
                attr_list[i].value.s32));
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreatePortSerdesRequest convert_create_port_serdes(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreatePortSerdesRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_SERDES_ATTR_PORT_ID:
        msg.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_PORT_SERDES_ATTR_PREEMPHASIS:
        msg.mutable_preemphasis()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_IDRIVER:
        msg.mutable_idriver()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_IPREDRIVER:
        msg.mutable_ipredriver()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_PRE1:
        msg.mutable_tx_fir_pre1()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_PRE2:
        msg.mutable_tx_fir_pre2()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_PRE3:
        msg.mutable_tx_fir_pre3()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_MAIN:
        msg.mutable_tx_fir_main()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_POST1:
        msg.mutable_tx_fir_post1()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_POST2:
        msg.mutable_tx_fir_post2()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_POST3:
        msg.mutable_tx_fir_post3()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_ATTN:
        msg.mutable_tx_fir_attn()->Add(
            attr_list[i].value.s32list.list,
            attr_list[i].value.s32list.list + attr_list[i].value.s32list.count);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_port(sai_object_id_t *port_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePortRequest req =
      convert_create_port(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreatePortResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = port->CreatePort(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *port_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_port(sai_object_id_t port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemovePortRequest req;
  lemming::dataplane::sai::RemovePortResponse resp;
  grpc::ClientContext context;
  req.set_oid(port_id);

  grpc::Status status = port->RemovePort(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_port_attribute(sai_object_id_t port_id,
                                  const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetPortAttributeRequest req;
  lemming::dataplane::sai::SetPortAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(port_id);

  switch (attr->id) {
    case SAI_PORT_ATTR_SPEED:
      req.set_speed(attr->value.u32);
      break;
    case SAI_PORT_ATTR_AUTO_NEG_MODE:
      req.set_auto_neg_mode(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_ADMIN_STATE:
      req.set_admin_state(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_MEDIA_TYPE:
      req.set_media_type(
          convert_sai_port_media_type_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_ADVERTISED_SPEED:
      req.mutable_advertised_speed()->Add(
          attr->value.u32list.list,
          attr->value.u32list.list + attr->value.u32list.count);
      break;
    case SAI_PORT_ATTR_ADVERTISED_HALF_DUPLEX_SPEED:
      req.mutable_advertised_half_duplex_speed()->Add(
          attr->value.u32list.list,
          attr->value.u32list.list + attr->value.u32list.count);
      break;
    case SAI_PORT_ATTR_ADVERTISED_AUTO_NEG_MODE:
      req.set_advertised_auto_neg_mode(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_ADVERTISED_FLOW_CONTROL_MODE:
      req.set_advertised_flow_control_mode(
          convert_sai_port_flow_control_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_ADVERTISED_ASYMMETRIC_PAUSE_MODE:
      req.set_advertised_asymmetric_pause_mode(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_ADVERTISED_MEDIA_TYPE:
      req.set_advertised_media_type(
          convert_sai_port_media_type_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_ADVERTISED_OUI_CODE:
      req.set_advertised_oui_code(attr->value.u32);
      break;
    case SAI_PORT_ATTR_PORT_VLAN_ID:
      req.set_port_vlan_id(attr->value.u16);
      break;
    case SAI_PORT_ATTR_DEFAULT_VLAN_PRIORITY:
      req.set_default_vlan_priority(attr->value.u8);
      break;
    case SAI_PORT_ATTR_DROP_UNTAGGED:
      req.set_drop_untagged(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_DROP_TAGGED:
      req.set_drop_tagged(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_INTERNAL_LOOPBACK_MODE:
      req.set_internal_loopback_mode(
          convert_sai_port_internal_loopback_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_USE_EXTENDED_FEC:
      req.set_use_extended_fec(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_FEC_MODE:
      req.set_fec_mode(convert_sai_port_fec_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_FEC_MODE_EXTENDED:
      req.set_fec_mode_extended(
          convert_sai_port_fec_mode_extended_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_UPDATE_DSCP:
      req.set_update_dscp(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_MTU:
      req.set_mtu(attr->value.u32);
      break;
    case SAI_PORT_ATTR_FLOOD_STORM_CONTROL_POLICER_ID:
      req.set_flood_storm_control_policer_id(attr->value.oid);
      break;
    case SAI_PORT_ATTR_BROADCAST_STORM_CONTROL_POLICER_ID:
      req.set_broadcast_storm_control_policer_id(attr->value.oid);
      break;
    case SAI_PORT_ATTR_MULTICAST_STORM_CONTROL_POLICER_ID:
      req.set_multicast_storm_control_policer_id(attr->value.oid);
      break;
    case SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_MODE:
      req.set_global_flow_control_mode(
          convert_sai_port_flow_control_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_INGRESS_ACL:
      req.set_ingress_acl(attr->value.oid);
      break;
    case SAI_PORT_ATTR_EGRESS_ACL:
      req.set_egress_acl(attr->value.oid);
      break;
    case SAI_PORT_ATTR_INGRESS_MACSEC_ACL:
      req.set_ingress_macsec_acl(attr->value.oid);
      break;
    case SAI_PORT_ATTR_EGRESS_MACSEC_ACL:
      req.set_egress_macsec_acl(attr->value.oid);
      break;
    case SAI_PORT_ATTR_INGRESS_MIRROR_SESSION:
      req.mutable_ingress_mirror_session()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_PORT_ATTR_EGRESS_MIRROR_SESSION:
      req.mutable_egress_mirror_session()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_PORT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
      req.set_ingress_samplepacket_enable(attr->value.oid);
      break;
    case SAI_PORT_ATTR_EGRESS_SAMPLEPACKET_ENABLE:
      req.set_egress_samplepacket_enable(attr->value.oid);
      break;
    case SAI_PORT_ATTR_INGRESS_SAMPLE_MIRROR_SESSION:
      req.mutable_ingress_sample_mirror_session()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_PORT_ATTR_EGRESS_SAMPLE_MIRROR_SESSION:
      req.mutable_egress_sample_mirror_session()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_PORT_ATTR_POLICER_ID:
      req.set_policer_id(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_DEFAULT_TC:
      req.set_qos_default_tc(attr->value.u8);
      break;
    case SAI_PORT_ATTR_QOS_DOT1P_TO_TC_MAP:
      req.set_qos_dot1p_to_tc_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_DOT1P_TO_COLOR_MAP:
      req.set_qos_dot1p_to_color_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_DSCP_TO_TC_MAP:
      req.set_qos_dscp_to_tc_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_DSCP_TO_COLOR_MAP:
      req.set_qos_dscp_to_color_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
      req.set_qos_tc_to_queue_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
      req.set_qos_tc_and_color_to_dot1p_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
      req.set_qos_tc_and_color_to_dscp_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_TC_TO_PRIORITY_GROUP_MAP:
      req.set_qos_tc_to_priority_group_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_PRIORITY_GROUP_MAP:
      req.set_qos_pfc_priority_to_priority_group_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_QUEUE_MAP:
      req.set_qos_pfc_priority_to_queue_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_SCHEDULER_PROFILE_ID:
      req.set_qos_scheduler_profile_id(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_INGRESS_BUFFER_PROFILE_LIST:
      req.mutable_qos_ingress_buffer_profile_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_PORT_ATTR_QOS_EGRESS_BUFFER_PROFILE_LIST:
      req.mutable_qos_egress_buffer_profile_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_MODE:
      req.set_priority_flow_control_mode(
          convert_sai_port_priority_flow_control_mode_t_to_proto(
              attr->value.s32));
      break;
    case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL:
      req.set_priority_flow_control(attr->value.u8);
      break;
    case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_RX:
      req.set_priority_flow_control_rx(attr->value.u8);
      break;
    case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_TX:
      req.set_priority_flow_control_tx(attr->value.u8);
      break;
    case SAI_PORT_ATTR_META_DATA:
      req.set_meta_data(attr->value.u32);
      break;
    case SAI_PORT_ATTR_EGRESS_BLOCK_PORT_LIST:
      req.mutable_egress_block_port_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_PORT_ATTR_HW_PROFILE_ID:
      req.set_hw_profile_id(attr->value.u64);
      break;
    case SAI_PORT_ATTR_EEE_ENABLE:
      req.set_eee_enable(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_EEE_IDLE_TIME:
      req.set_eee_idle_time(attr->value.u16);
      break;
    case SAI_PORT_ATTR_EEE_WAKE_TIME:
      req.set_eee_wake_time(attr->value.u16);
      break;
    case SAI_PORT_ATTR_ISOLATION_GROUP:
      req.set_isolation_group(attr->value.oid);
      break;
    case SAI_PORT_ATTR_PKT_TX_ENABLE:
      req.set_pkt_tx_enable(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_TAM_OBJECT:
      req.mutable_tam_object()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_PORT_ATTR_SERDES_PREEMPHASIS:
      req.mutable_serdes_preemphasis()->Add(
          attr->value.u32list.list,
          attr->value.u32list.list + attr->value.u32list.count);
      break;
    case SAI_PORT_ATTR_SERDES_IDRIVER:
      req.mutable_serdes_idriver()->Add(
          attr->value.u32list.list,
          attr->value.u32list.list + attr->value.u32list.count);
      break;
    case SAI_PORT_ATTR_SERDES_IPREDRIVER:
      req.mutable_serdes_ipredriver()->Add(
          attr->value.u32list.list,
          attr->value.u32list.list + attr->value.u32list.count);
      break;
    case SAI_PORT_ATTR_LINK_TRAINING_ENABLE:
      req.set_link_training_enable(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_PTP_MODE:
      req.set_ptp_mode(convert_sai_port_ptp_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_INTERFACE_TYPE:
      req.set_interface_type(
          convert_sai_port_interface_type_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_PRBS_POLYNOMIAL:
      req.set_prbs_polynomial(attr->value.u32);
      break;
    case SAI_PORT_ATTR_PRBS_CONFIG:
      req.set_prbs_config(
          convert_sai_port_prbs_config_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_DISABLE_DECREMENT_TTL:
      req.set_disable_decrement_ttl(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
      req.set_qos_mpls_exp_to_tc_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
      req.set_qos_mpls_exp_to_color_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
      req.set_qos_tc_and_color_to_mpls_exp_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_TPID:
      req.set_tpid(attr->value.u16);
      break;
    case SAI_PORT_ATTR_AUTO_NEG_FEC_MODE_OVERRIDE:
      req.set_auto_neg_fec_mode_override(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_LOOPBACK_MODE:
      req.set_loopback_mode(
          convert_sai_port_loopback_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_MDIX_MODE_CONFIG:
      req.set_mdix_mode_config(
          convert_sai_port_mdix_mode_config_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_AUTO_NEG_CONFIG_MODE:
      req.set_auto_neg_config_mode(
          convert_sai_port_auto_neg_config_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_1000X_SGMII_SLAVE_AUTODETECT:
      req.set__1000x_sgmii_slave_autodetect(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_MODULE_TYPE:
      req.set_module_type(
          convert_sai_port_module_type_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_DUAL_MEDIA:
      req.set_dual_media(
          convert_sai_port_dual_media_t_to_proto(attr->value.s32));
      break;
    case SAI_PORT_ATTR_IPG:
      req.set_ipg(attr->value.u32);
      break;
    case SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_FORWARD:
      req.set_global_flow_control_forward(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_FORWARD:
      req.set_priority_flow_control_forward(attr->value.booldata);
      break;
    case SAI_PORT_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
      req.set_qos_dscp_to_forwarding_class_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
      req.set_qos_mpls_exp_to_forwarding_class_map(attr->value.oid);
      break;
    case SAI_PORT_ATTR_FABRIC_ISOLATE:
      req.set_fabric_isolate(attr->value.booldata);
      break;
  }

  grpc::Status status = port->SetPortAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_port_attribute(sai_object_id_t port_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetPortAttributeRequest req;
  lemming::dataplane::sai::GetPortAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(port_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_port_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = port->GetPortAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_ATTR_TYPE:
        attr_list[i].value.s32 =
            convert_sai_port_type_t_to_sai(resp.attr().type());
        break;
      case SAI_PORT_ATTR_OPER_STATUS:
        attr_list[i].value.s32 =
            convert_sai_port_oper_status_t_to_sai(resp.attr().oper_status());
        break;
      case SAI_PORT_ATTR_CURRENT_BREAKOUT_MODE_TYPE:
        attr_list[i].value.s32 = convert_sai_port_breakout_mode_type_t_to_sai(
            resp.attr().current_breakout_mode_type());
        break;
      case SAI_PORT_ATTR_QOS_NUMBER_OF_QUEUES:
        attr_list[i].value.u32 = resp.attr().qos_number_of_queues();
        break;
      case SAI_PORT_ATTR_QOS_QUEUE_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().qos_queue_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_QOS_NUMBER_OF_SCHEDULER_GROUPS:
        attr_list[i].value.u32 = resp.attr().qos_number_of_scheduler_groups();
        break;
      case SAI_PORT_ATTR_QOS_SCHEDULER_GROUP_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().qos_scheduler_group_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_QOS_MAXIMUM_HEADROOM_SIZE:
        attr_list[i].value.u32 = resp.attr().qos_maximum_headroom_size();
        break;
      case SAI_PORT_ATTR_SUPPORTED_SPEED:
        copy_list(attr_list[i].value.u32list.list,
                  resp.attr().supported_speed(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SUPPORTED_HALF_DUPLEX_SPEED:
        copy_list(attr_list[i].value.u32list.list,
                  resp.attr().supported_half_duplex_speed(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SUPPORTED_AUTO_NEG_MODE:
        attr_list[i].value.booldata = resp.attr().supported_auto_neg_mode();
        break;
      case SAI_PORT_ATTR_SUPPORTED_FLOW_CONTROL_MODE:
        attr_list[i].value.s32 = convert_sai_port_flow_control_mode_t_to_sai(
            resp.attr().supported_flow_control_mode());
        break;
      case SAI_PORT_ATTR_SUPPORTED_ASYMMETRIC_PAUSE_MODE:
        attr_list[i].value.booldata =
            resp.attr().supported_asymmetric_pause_mode();
        break;
      case SAI_PORT_ATTR_SUPPORTED_MEDIA_TYPE:
        attr_list[i].value.s32 = convert_sai_port_media_type_t_to_sai(
            resp.attr().supported_media_type());
        break;
      case SAI_PORT_ATTR_REMOTE_ADVERTISED_SPEED:
        copy_list(attr_list[i].value.u32list.list,
                  resp.attr().remote_advertised_speed(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_REMOTE_ADVERTISED_HALF_DUPLEX_SPEED:
        copy_list(attr_list[i].value.u32list.list,
                  resp.attr().remote_advertised_half_duplex_speed(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_REMOTE_ADVERTISED_AUTO_NEG_MODE:
        attr_list[i].value.booldata =
            resp.attr().remote_advertised_auto_neg_mode();
        break;
      case SAI_PORT_ATTR_REMOTE_ADVERTISED_FLOW_CONTROL_MODE:
        attr_list[i].value.s32 = convert_sai_port_flow_control_mode_t_to_sai(
            resp.attr().remote_advertised_flow_control_mode());
        break;
      case SAI_PORT_ATTR_REMOTE_ADVERTISED_ASYMMETRIC_PAUSE_MODE:
        attr_list[i].value.booldata =
            resp.attr().remote_advertised_asymmetric_pause_mode();
        break;
      case SAI_PORT_ATTR_REMOTE_ADVERTISED_MEDIA_TYPE:
        attr_list[i].value.s32 = convert_sai_port_media_type_t_to_sai(
            resp.attr().remote_advertised_media_type());
        break;
      case SAI_PORT_ATTR_REMOTE_ADVERTISED_OUI_CODE:
        attr_list[i].value.u32 = resp.attr().remote_advertised_oui_code();
        break;
      case SAI_PORT_ATTR_NUMBER_OF_INGRESS_PRIORITY_GROUPS:
        attr_list[i].value.u32 =
            resp.attr().number_of_ingress_priority_groups();
        break;
      case SAI_PORT_ATTR_INGRESS_PRIORITY_GROUP_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().ingress_priority_group_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_OPER_SPEED:
        attr_list[i].value.u32 = resp.attr().oper_speed();
        break;
      case SAI_PORT_ATTR_HW_LANE_LIST:
        copy_list(attr_list[i].value.u32list.list, resp.attr().hw_lane_list(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SPEED:
        attr_list[i].value.u32 = resp.attr().speed();
        break;
      case SAI_PORT_ATTR_FULL_DUPLEX_MODE:
        attr_list[i].value.booldata = resp.attr().full_duplex_mode();
        break;
      case SAI_PORT_ATTR_AUTO_NEG_MODE:
        attr_list[i].value.booldata = resp.attr().auto_neg_mode();
        break;
      case SAI_PORT_ATTR_ADMIN_STATE:
        attr_list[i].value.booldata = resp.attr().admin_state();
        break;
      case SAI_PORT_ATTR_MEDIA_TYPE:
        attr_list[i].value.s32 =
            convert_sai_port_media_type_t_to_sai(resp.attr().media_type());
        break;
      case SAI_PORT_ATTR_ADVERTISED_SPEED:
        copy_list(attr_list[i].value.u32list.list,
                  resp.attr().advertised_speed(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_ADVERTISED_HALF_DUPLEX_SPEED:
        copy_list(attr_list[i].value.u32list.list,
                  resp.attr().advertised_half_duplex_speed(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_ADVERTISED_AUTO_NEG_MODE:
        attr_list[i].value.booldata = resp.attr().advertised_auto_neg_mode();
        break;
      case SAI_PORT_ATTR_ADVERTISED_FLOW_CONTROL_MODE:
        attr_list[i].value.s32 = convert_sai_port_flow_control_mode_t_to_sai(
            resp.attr().advertised_flow_control_mode());
        break;
      case SAI_PORT_ATTR_ADVERTISED_ASYMMETRIC_PAUSE_MODE:
        attr_list[i].value.booldata =
            resp.attr().advertised_asymmetric_pause_mode();
        break;
      case SAI_PORT_ATTR_ADVERTISED_MEDIA_TYPE:
        attr_list[i].value.s32 = convert_sai_port_media_type_t_to_sai(
            resp.attr().advertised_media_type());
        break;
      case SAI_PORT_ATTR_ADVERTISED_OUI_CODE:
        attr_list[i].value.u32 = resp.attr().advertised_oui_code();
        break;
      case SAI_PORT_ATTR_PORT_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().port_vlan_id();
        break;
      case SAI_PORT_ATTR_DEFAULT_VLAN_PRIORITY:
        attr_list[i].value.u8 = resp.attr().default_vlan_priority();
        break;
      case SAI_PORT_ATTR_DROP_UNTAGGED:
        attr_list[i].value.booldata = resp.attr().drop_untagged();
        break;
      case SAI_PORT_ATTR_DROP_TAGGED:
        attr_list[i].value.booldata = resp.attr().drop_tagged();
        break;
      case SAI_PORT_ATTR_INTERNAL_LOOPBACK_MODE:
        attr_list[i].value.s32 =
            convert_sai_port_internal_loopback_mode_t_to_sai(
                resp.attr().internal_loopback_mode());
        break;
      case SAI_PORT_ATTR_USE_EXTENDED_FEC:
        attr_list[i].value.booldata = resp.attr().use_extended_fec();
        break;
      case SAI_PORT_ATTR_FEC_MODE:
        attr_list[i].value.s32 =
            convert_sai_port_fec_mode_t_to_sai(resp.attr().fec_mode());
        break;
      case SAI_PORT_ATTR_FEC_MODE_EXTENDED:
        attr_list[i].value.s32 = convert_sai_port_fec_mode_extended_t_to_sai(
            resp.attr().fec_mode_extended());
        break;
      case SAI_PORT_ATTR_UPDATE_DSCP:
        attr_list[i].value.booldata = resp.attr().update_dscp();
        break;
      case SAI_PORT_ATTR_MTU:
        attr_list[i].value.u32 = resp.attr().mtu();
        break;
      case SAI_PORT_ATTR_FLOOD_STORM_CONTROL_POLICER_ID:
        attr_list[i].value.oid = resp.attr().flood_storm_control_policer_id();
        break;
      case SAI_PORT_ATTR_BROADCAST_STORM_CONTROL_POLICER_ID:
        attr_list[i].value.oid =
            resp.attr().broadcast_storm_control_policer_id();
        break;
      case SAI_PORT_ATTR_MULTICAST_STORM_CONTROL_POLICER_ID:
        attr_list[i].value.oid =
            resp.attr().multicast_storm_control_policer_id();
        break;
      case SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_MODE:
        attr_list[i].value.s32 = convert_sai_port_flow_control_mode_t_to_sai(
            resp.attr().global_flow_control_mode());
        break;
      case SAI_PORT_ATTR_INGRESS_ACL:
        attr_list[i].value.oid = resp.attr().ingress_acl();
        break;
      case SAI_PORT_ATTR_EGRESS_ACL:
        attr_list[i].value.oid = resp.attr().egress_acl();
        break;
      case SAI_PORT_ATTR_INGRESS_MACSEC_ACL:
        attr_list[i].value.oid = resp.attr().ingress_macsec_acl();
        break;
      case SAI_PORT_ATTR_EGRESS_MACSEC_ACL:
        attr_list[i].value.oid = resp.attr().egress_macsec_acl();
        break;
      case SAI_PORT_ATTR_MACSEC_PORT_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().macsec_port_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_INGRESS_MIRROR_SESSION:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().ingress_mirror_session(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_EGRESS_MIRROR_SESSION:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().egress_mirror_session(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
        attr_list[i].value.oid = resp.attr().ingress_samplepacket_enable();
        break;
      case SAI_PORT_ATTR_EGRESS_SAMPLEPACKET_ENABLE:
        attr_list[i].value.oid = resp.attr().egress_samplepacket_enable();
        break;
      case SAI_PORT_ATTR_INGRESS_SAMPLE_MIRROR_SESSION:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().ingress_sample_mirror_session(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_EGRESS_SAMPLE_MIRROR_SESSION:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().egress_sample_mirror_session(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_POLICER_ID:
        attr_list[i].value.oid = resp.attr().policer_id();
        break;
      case SAI_PORT_ATTR_QOS_DEFAULT_TC:
        attr_list[i].value.u8 = resp.attr().qos_default_tc();
        break;
      case SAI_PORT_ATTR_QOS_DOT1P_TO_TC_MAP:
        attr_list[i].value.oid = resp.attr().qos_dot1p_to_tc_map();
        break;
      case SAI_PORT_ATTR_QOS_DOT1P_TO_COLOR_MAP:
        attr_list[i].value.oid = resp.attr().qos_dot1p_to_color_map();
        break;
      case SAI_PORT_ATTR_QOS_DSCP_TO_TC_MAP:
        attr_list[i].value.oid = resp.attr().qos_dscp_to_tc_map();
        break;
      case SAI_PORT_ATTR_QOS_DSCP_TO_COLOR_MAP:
        attr_list[i].value.oid = resp.attr().qos_dscp_to_color_map();
        break;
      case SAI_PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
        attr_list[i].value.oid = resp.attr().qos_tc_to_queue_map();
        break;
      case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
        attr_list[i].value.oid = resp.attr().qos_tc_and_color_to_dot1p_map();
        break;
      case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
        attr_list[i].value.oid = resp.attr().qos_tc_and_color_to_dscp_map();
        break;
      case SAI_PORT_ATTR_QOS_TC_TO_PRIORITY_GROUP_MAP:
        attr_list[i].value.oid = resp.attr().qos_tc_to_priority_group_map();
        break;
      case SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_PRIORITY_GROUP_MAP:
        attr_list[i].value.oid =
            resp.attr().qos_pfc_priority_to_priority_group_map();
        break;
      case SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_QUEUE_MAP:
        attr_list[i].value.oid = resp.attr().qos_pfc_priority_to_queue_map();
        break;
      case SAI_PORT_ATTR_QOS_SCHEDULER_PROFILE_ID:
        attr_list[i].value.oid = resp.attr().qos_scheduler_profile_id();
        break;
      case SAI_PORT_ATTR_QOS_INGRESS_BUFFER_PROFILE_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().qos_ingress_buffer_profile_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_QOS_EGRESS_BUFFER_PROFILE_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().qos_egress_buffer_profile_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_MODE:
        attr_list[i].value.s32 =
            convert_sai_port_priority_flow_control_mode_t_to_sai(
                resp.attr().priority_flow_control_mode());
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL:
        attr_list[i].value.u8 = resp.attr().priority_flow_control();
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_RX:
        attr_list[i].value.u8 = resp.attr().priority_flow_control_rx();
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_TX:
        attr_list[i].value.u8 = resp.attr().priority_flow_control_tx();
        break;
      case SAI_PORT_ATTR_META_DATA:
        attr_list[i].value.u32 = resp.attr().meta_data();
        break;
      case SAI_PORT_ATTR_EGRESS_BLOCK_PORT_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().egress_block_port_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_HW_PROFILE_ID:
        attr_list[i].value.u64 = resp.attr().hw_profile_id();
        break;
      case SAI_PORT_ATTR_EEE_ENABLE:
        attr_list[i].value.booldata = resp.attr().eee_enable();
        break;
      case SAI_PORT_ATTR_EEE_IDLE_TIME:
        attr_list[i].value.u16 = resp.attr().eee_idle_time();
        break;
      case SAI_PORT_ATTR_EEE_WAKE_TIME:
        attr_list[i].value.u16 = resp.attr().eee_wake_time();
        break;
      case SAI_PORT_ATTR_PORT_POOL_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().port_pool_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_ISOLATION_GROUP:
        attr_list[i].value.oid = resp.attr().isolation_group();
        break;
      case SAI_PORT_ATTR_PKT_TX_ENABLE:
        attr_list[i].value.booldata = resp.attr().pkt_tx_enable();
        break;
      case SAI_PORT_ATTR_TAM_OBJECT:
        copy_list(attr_list[i].value.objlist.list, resp.attr().tam_object(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_PORT_ATTR_SERDES_PREEMPHASIS:
        copy_list(attr_list[i].value.u32list.list,
                  resp.attr().serdes_preemphasis(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SERDES_IDRIVER:
        copy_list(attr_list[i].value.u32list.list, resp.attr().serdes_idriver(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_SERDES_IPREDRIVER:
        copy_list(attr_list[i].value.u32list.list,
                  resp.attr().serdes_ipredriver(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_PORT_ATTR_LINK_TRAINING_ENABLE:
        attr_list[i].value.booldata = resp.attr().link_training_enable();
        break;
      case SAI_PORT_ATTR_PTP_MODE:
        attr_list[i].value.s32 =
            convert_sai_port_ptp_mode_t_to_sai(resp.attr().ptp_mode());
        break;
      case SAI_PORT_ATTR_INTERFACE_TYPE:
        attr_list[i].value.s32 = convert_sai_port_interface_type_t_to_sai(
            resp.attr().interface_type());
        break;
      case SAI_PORT_ATTR_REFERENCE_CLOCK:
        attr_list[i].value.u64 = resp.attr().reference_clock();
        break;
      case SAI_PORT_ATTR_PRBS_POLYNOMIAL:
        attr_list[i].value.u32 = resp.attr().prbs_polynomial();
        break;
      case SAI_PORT_ATTR_PORT_SERDES_ID:
        attr_list[i].value.oid = resp.attr().port_serdes_id();
        break;
      case SAI_PORT_ATTR_LINK_TRAINING_FAILURE_STATUS:
        attr_list[i].value.s32 =
            convert_sai_port_link_training_failure_status_t_to_sai(
                resp.attr().link_training_failure_status());
        break;
      case SAI_PORT_ATTR_LINK_TRAINING_RX_STATUS:
        attr_list[i].value.s32 =
            convert_sai_port_link_training_rx_status_t_to_sai(
                resp.attr().link_training_rx_status());
        break;
      case SAI_PORT_ATTR_PRBS_CONFIG:
        attr_list[i].value.s32 =
            convert_sai_port_prbs_config_t_to_sai(resp.attr().prbs_config());
        break;
      case SAI_PORT_ATTR_PRBS_LOCK_STATUS:
        attr_list[i].value.booldata = resp.attr().prbs_lock_status();
        break;
      case SAI_PORT_ATTR_PRBS_LOCK_LOSS_STATUS:
        attr_list[i].value.booldata = resp.attr().prbs_lock_loss_status();
        break;
      case SAI_PORT_ATTR_PRBS_RX_STATUS:
        attr_list[i].value.s32 = convert_sai_port_prbs_rx_status_t_to_sai(
            resp.attr().prbs_rx_status());
        break;
      case SAI_PORT_ATTR_AUTO_NEG_STATUS:
        attr_list[i].value.booldata = resp.attr().auto_neg_status();
        break;
      case SAI_PORT_ATTR_DISABLE_DECREMENT_TTL:
        attr_list[i].value.booldata = resp.attr().disable_decrement_ttl();
        break;
      case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
        attr_list[i].value.oid = resp.attr().qos_mpls_exp_to_tc_map();
        break;
      case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
        attr_list[i].value.oid = resp.attr().qos_mpls_exp_to_color_map();
        break;
      case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
        attr_list[i].value.oid = resp.attr().qos_tc_and_color_to_mpls_exp_map();
        break;
      case SAI_PORT_ATTR_TPID:
        attr_list[i].value.u16 = resp.attr().tpid();
        break;
      case SAI_PORT_ATTR_FABRIC_ATTACHED:
        attr_list[i].value.booldata = resp.attr().fabric_attached();
        break;
      case SAI_PORT_ATTR_FABRIC_ATTACHED_SWITCH_TYPE:
        attr_list[i].value.s32 = convert_sai_switch_type_t_to_sai(
            resp.attr().fabric_attached_switch_type());
        break;
      case SAI_PORT_ATTR_FABRIC_ATTACHED_SWITCH_ID:
        attr_list[i].value.u32 = resp.attr().fabric_attached_switch_id();
        break;
      case SAI_PORT_ATTR_FABRIC_ATTACHED_PORT_INDEX:
        attr_list[i].value.u32 = resp.attr().fabric_attached_port_index();
        break;
      case SAI_PORT_ATTR_SYSTEM_PORT:
        attr_list[i].value.oid = resp.attr().system_port();
        break;
      case SAI_PORT_ATTR_AUTO_NEG_FEC_MODE_OVERRIDE:
        attr_list[i].value.booldata = resp.attr().auto_neg_fec_mode_override();
        break;
      case SAI_PORT_ATTR_LOOPBACK_MODE:
        attr_list[i].value.s32 = convert_sai_port_loopback_mode_t_to_sai(
            resp.attr().loopback_mode());
        break;
      case SAI_PORT_ATTR_MDIX_MODE_STATUS:
        attr_list[i].value.s32 = convert_sai_port_mdix_mode_status_t_to_sai(
            resp.attr().mdix_mode_status());
        break;
      case SAI_PORT_ATTR_MDIX_MODE_CONFIG:
        attr_list[i].value.s32 = convert_sai_port_mdix_mode_config_t_to_sai(
            resp.attr().mdix_mode_config());
        break;
      case SAI_PORT_ATTR_AUTO_NEG_CONFIG_MODE:
        attr_list[i].value.s32 = convert_sai_port_auto_neg_config_mode_t_to_sai(
            resp.attr().auto_neg_config_mode());
        break;
      case SAI_PORT_ATTR_1000X_SGMII_SLAVE_AUTODETECT:
        attr_list[i].value.booldata =
            resp.attr()._1000x_sgmii_slave_autodetect();
        break;
      case SAI_PORT_ATTR_MODULE_TYPE:
        attr_list[i].value.s32 =
            convert_sai_port_module_type_t_to_sai(resp.attr().module_type());
        break;
      case SAI_PORT_ATTR_DUAL_MEDIA:
        attr_list[i].value.s32 =
            convert_sai_port_dual_media_t_to_sai(resp.attr().dual_media());
        break;
      case SAI_PORT_ATTR_AUTO_NEG_FEC_MODE_EXTENDED:
        attr_list[i].value.s32 = convert_sai_port_fec_mode_extended_t_to_sai(
            resp.attr().auto_neg_fec_mode_extended());
        break;
      case SAI_PORT_ATTR_IPG:
        attr_list[i].value.u32 = resp.attr().ipg();
        break;
      case SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_FORWARD:
        attr_list[i].value.booldata = resp.attr().global_flow_control_forward();
        break;
      case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_FORWARD:
        attr_list[i].value.booldata =
            resp.attr().priority_flow_control_forward();
        break;
      case SAI_PORT_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
        attr_list[i].value.oid = resp.attr().qos_dscp_to_forwarding_class_map();
        break;
      case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
        attr_list[i].value.oid =
            resp.attr().qos_mpls_exp_to_forwarding_class_map();
        break;
      case SAI_PORT_ATTR_IPSEC_PORT:
        attr_list[i].value.oid = resp.attr().ipsec_port();
        break;
      case SAI_PORT_ATTR_SUPPORTED_LINK_TRAINING_MODE:
        attr_list[i].value.booldata =
            resp.attr().supported_link_training_mode();
        break;
      case SAI_PORT_ATTR_FABRIC_ISOLATE:
        attr_list[i].value.booldata = resp.attr().fabric_isolate();
        break;
      case SAI_PORT_ATTR_MAX_FEC_SYMBOL_ERRORS_DETECTABLE:
        attr_list[i].value.u32 = resp.attr().max_fec_symbol_errors_detectable();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_port_stats(sai_object_id_t port_id,
                              uint32_t number_of_counters,
                              const sai_stat_id_t *counter_ids,
                              uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetPortStatsRequest req;
  lemming::dataplane::sai::GetPortStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(port_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(convert_sai_port_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status = port->GetPortStats(&context, req, &resp);
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

sai_status_t l_get_port_stats_ext(sai_object_id_t port_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids,
                                  sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_port_stats(sai_object_id_t port_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_port_pool(sai_object_id_t *port_pool_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePortPoolRequest req =
      convert_create_port_pool(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreatePortPoolResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = port->CreatePortPool(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *port_pool_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_port_pool(sai_object_id_t port_pool_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemovePortPoolRequest req;
  lemming::dataplane::sai::RemovePortPoolResponse resp;
  grpc::ClientContext context;
  req.set_oid(port_pool_id);

  grpc::Status status = port->RemovePortPool(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_port_pool_attribute(sai_object_id_t port_pool_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetPortPoolAttributeRequest req;
  lemming::dataplane::sai::SetPortPoolAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(port_pool_id);

  switch (attr->id) {
    case SAI_PORT_POOL_ATTR_QOS_WRED_PROFILE_ID:
      req.set_qos_wred_profile_id(attr->value.oid);
      break;
  }

  grpc::Status status = port->SetPortPoolAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_port_pool_attribute(sai_object_id_t port_pool_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetPortPoolAttributeRequest req;
  lemming::dataplane::sai::GetPortPoolAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(port_pool_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_port_pool_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = port->GetPortPoolAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_POOL_ATTR_PORT_ID:
        attr_list[i].value.oid = resp.attr().port_id();
        break;
      case SAI_PORT_POOL_ATTR_BUFFER_POOL_ID:
        attr_list[i].value.oid = resp.attr().buffer_pool_id();
        break;
      case SAI_PORT_POOL_ATTR_QOS_WRED_PROFILE_ID:
        attr_list[i].value.oid = resp.attr().qos_wred_profile_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_port_pool_stats(sai_object_id_t port_pool_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids,
                                   uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetPortPoolStatsRequest req;
  lemming::dataplane::sai::GetPortPoolStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(port_pool_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(convert_sai_port_pool_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status = port->GetPortPoolStats(&context, req, &resp);
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

sai_status_t l_get_port_pool_stats_ext(sai_object_id_t port_pool_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids,
                                       sai_stats_mode_t mode,
                                       uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_port_pool_stats(sai_object_id_t port_pool_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_port_connector(sai_object_id_t *port_connector_id,
                                     sai_object_id_t switch_id,
                                     uint32_t attr_count,
                                     const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePortConnectorRequest req =
      convert_create_port_connector(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreatePortConnectorResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = port->CreatePortConnector(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *port_connector_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_port_connector(sai_object_id_t port_connector_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemovePortConnectorRequest req;
  lemming::dataplane::sai::RemovePortConnectorResponse resp;
  grpc::ClientContext context;
  req.set_oid(port_connector_id);

  grpc::Status status = port->RemovePortConnector(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_port_connector_attribute(sai_object_id_t port_connector_id,
                                            const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetPortConnectorAttributeRequest req;
  lemming::dataplane::sai::SetPortConnectorAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(port_connector_id);

  switch (attr->id) {
    case SAI_PORT_CONNECTOR_ATTR_FAILOVER_MODE:
      req.set_failover_mode(
          convert_sai_port_connector_failover_mode_t_to_proto(attr->value.s32));
      break;
  }

  grpc::Status status = port->SetPortConnectorAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_port_connector_attribute(sai_object_id_t port_connector_id,
                                            uint32_t attr_count,
                                            sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetPortConnectorAttributeRequest req;
  lemming::dataplane::sai::GetPortConnectorAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(port_connector_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_port_connector_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = port->GetPortConnectorAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_PORT_ID:
        attr_list[i].value.oid = resp.attr().system_side_port_id();
        break;
      case SAI_PORT_CONNECTOR_ATTR_LINE_SIDE_PORT_ID:
        attr_list[i].value.oid = resp.attr().line_side_port_id();
        break;
      case SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_FAILOVER_PORT_ID:
        attr_list[i].value.oid = resp.attr().system_side_failover_port_id();
        break;
      case SAI_PORT_CONNECTOR_ATTR_LINE_SIDE_FAILOVER_PORT_ID:
        attr_list[i].value.oid = resp.attr().line_side_failover_port_id();
        break;
      case SAI_PORT_CONNECTOR_ATTR_FAILOVER_MODE:
        attr_list[i].value.s32 =
            convert_sai_port_connector_failover_mode_t_to_sai(
                resp.attr().failover_mode());
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_port_serdes(sai_object_id_t *port_serdes_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePortSerdesRequest req =
      convert_create_port_serdes(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreatePortSerdesResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = port->CreatePortSerdes(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *port_serdes_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_port_serdes(sai_object_id_t port_serdes_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemovePortSerdesRequest req;
  lemming::dataplane::sai::RemovePortSerdesResponse resp;
  grpc::ClientContext context;
  req.set_oid(port_serdes_id);

  grpc::Status status = port->RemovePortSerdes(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_port_serdes_attribute(sai_object_id_t port_serdes_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_port_serdes_attribute(sai_object_id_t port_serdes_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetPortSerdesAttributeRequest req;
  lemming::dataplane::sai::GetPortSerdesAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(port_serdes_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_port_serdes_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = port->GetPortSerdesAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_PORT_SERDES_ATTR_PORT_ID:
        attr_list[i].value.oid = resp.attr().port_id();
        break;
      case SAI_PORT_SERDES_ATTR_PREEMPHASIS:
        copy_list(attr_list[i].value.s32list.list, resp.attr().preemphasis(),
                  &attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_IDRIVER:
        copy_list(attr_list[i].value.s32list.list, resp.attr().idriver(),
                  &attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_IPREDRIVER:
        copy_list(attr_list[i].value.s32list.list, resp.attr().ipredriver(),
                  &attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_PRE1:
        copy_list(attr_list[i].value.s32list.list, resp.attr().tx_fir_pre1(),
                  &attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_PRE2:
        copy_list(attr_list[i].value.s32list.list, resp.attr().tx_fir_pre2(),
                  &attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_PRE3:
        copy_list(attr_list[i].value.s32list.list, resp.attr().tx_fir_pre3(),
                  &attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_MAIN:
        copy_list(attr_list[i].value.s32list.list, resp.attr().tx_fir_main(),
                  &attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_POST1:
        copy_list(attr_list[i].value.s32list.list, resp.attr().tx_fir_post1(),
                  &attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_POST2:
        copy_list(attr_list[i].value.s32list.list, resp.attr().tx_fir_post2(),
                  &attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_POST3:
        copy_list(attr_list[i].value.s32list.list, resp.attr().tx_fir_post3(),
                  &attr_list[i].value.s32list.count);
        break;
      case SAI_PORT_SERDES_ATTR_TX_FIR_ATTN:
        copy_list(attr_list[i].value.s32list.list, resp.attr().tx_fir_attn(),
                  &attr_list[i].value.s32list.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_ports(sai_object_id_t switch_id, uint32_t object_count,
                            const uint32_t *attr_count,
                            const sai_attribute_t **attr_list,
                            sai_bulk_op_error_mode_t mode,
                            sai_object_id_t *object_id,
                            sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePortsRequest req;
  lemming::dataplane::sai::CreatePortsResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    auto r = convert_create_port(switch_id, attr_count[i], attr_list[i]);
    *req.add_reqs() = r;
  }

  grpc::Status status = port->CreatePorts(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
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

sai_status_t l_remove_ports(uint32_t object_count,
                            const sai_object_id_t *object_id,
                            sai_bulk_op_error_mode_t mode,
                            sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemovePortsRequest req;
  lemming::dataplane::sai::RemovePortsResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    req.add_reqs()->set_oid(object_id[i]);
  }

  grpc::Status status = port->RemovePorts(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
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

sai_status_t l_set_ports_attribute(uint32_t object_count,
                                   const sai_object_id_t *object_id,
                                   const sai_attribute_t *attr_list,
                                   sai_bulk_op_error_mode_t mode,
                                   sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_ports_attribute(uint32_t object_count,
                                   const sai_object_id_t *object_id,
                                   const uint32_t *attr_count,
                                   sai_attribute_t **attr_list,
                                   sai_bulk_op_error_mode_t mode,
                                   sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}
