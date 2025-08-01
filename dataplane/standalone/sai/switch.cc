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
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/switch.pb.h"
#include <glog/logging.h>

const sai_switch_api_t l_switch = {
	.create_switch = l_create_switch,
	.remove_switch = l_remove_switch,
	.set_switch_attribute = l_set_switch_attribute,
	.get_switch_attribute = l_get_switch_attribute,
	.get_switch_stats = l_get_switch_stats,
	.get_switch_stats_ext = l_get_switch_stats_ext,
	.clear_switch_stats = l_clear_switch_stats,
	.create_switch_tunnel = l_create_switch_tunnel,
	.remove_switch_tunnel = l_remove_switch_tunnel,
	.set_switch_tunnel_attribute = l_set_switch_tunnel_attribute,
	.get_switch_tunnel_attribute = l_get_switch_tunnel_attribute,
};

std::unique_ptr<PortStateReactor> port_state;

lemming::dataplane::sai::CreateSwitchRequest convert_create_switch(uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateSwitchRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_SWITCH_ATTR_INGRESS_ACL:
	msg.set_ingress_acl(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_EGRESS_ACL:
	msg.set_egress_acl(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_RESTART_WARM:
	msg.set_restart_warm(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_WARM_RECOVER:
	msg.set_warm_recover(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_SWITCHING_MODE:
	msg.set_switching_mode(convert_sai_switch_switching_mode_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_BCAST_CPU_FLOOD_ENABLE:
	msg.set_bcast_cpu_flood_enable(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_MCAST_CPU_FLOOD_ENABLE:
	msg.set_mcast_cpu_flood_enable(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_SRC_MAC_ADDRESS:
	msg.set_src_mac_address(attr_list[i].value.mac, sizeof(attr_list[i].value.mac));
	break;
  case SAI_SWITCH_ATTR_MAX_LEARNED_ADDRESSES:
	msg.set_max_learned_addresses(attr_list[i].value.u32);
	break;
  case SAI_SWITCH_ATTR_FDB_AGING_TIME:
	msg.set_fdb_aging_time(attr_list[i].value.u32);
	break;
  case SAI_SWITCH_ATTR_FDB_UNICAST_MISS_PACKET_ACTION:
	msg.set_fdb_unicast_miss_packet_action(convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_FDB_BROADCAST_MISS_PACKET_ACTION:
	msg.set_fdb_broadcast_miss_packet_action(convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_FDB_MULTICAST_MISS_PACKET_ACTION:
	msg.set_fdb_multicast_miss_packet_action(convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_ALGORITHM:
	msg.set_ecmp_default_hash_algorithm(convert_sai_hash_algorithm_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_SEED:
	msg.set_ecmp_default_hash_seed(attr_list[i].value.u32);
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_OFFSET:
	msg.set_ecmp_default_hash_offset(attr_list[i].value.u8);
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_SYMMETRIC_HASH:
	msg.set_ecmp_default_symmetric_hash(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_ECMP_HASH_IPV4:
	msg.set_ecmp_hash_ipv4(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_ECMP_HASH_IPV4_IN_IPV4:
	msg.set_ecmp_hash_ipv4_in_ipv4(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_ECMP_HASH_IPV6:
	msg.set_ecmp_hash_ipv6(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM:
	msg.set_lag_default_hash_algorithm(convert_sai_hash_algorithm_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_SEED:
	msg.set_lag_default_hash_seed(attr_list[i].value.u32);
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_OFFSET:
	msg.set_lag_default_hash_offset(attr_list[i].value.u8);
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH:
	msg.set_lag_default_symmetric_hash(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_LAG_HASH_IPV4:
	msg.set_lag_hash_ipv4(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_LAG_HASH_IPV4_IN_IPV4:
	msg.set_lag_hash_ipv4_in_ipv4(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_LAG_HASH_IPV6:
	msg.set_lag_hash_ipv6(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_COUNTER_REFRESH_INTERVAL:
	msg.set_counter_refresh_interval(attr_list[i].value.u32);
	break;
  case SAI_SWITCH_ATTR_QOS_DEFAULT_TC:
	msg.set_qos_default_tc(attr_list[i].value.u8);
	break;
  case SAI_SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP:
	msg.set_qos_dot1p_to_tc_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP:
	msg.set_qos_dot1p_to_color_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_DSCP_TO_TC_MAP:
	msg.set_qos_dscp_to_tc_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_DSCP_TO_COLOR_MAP:
	msg.set_qos_dscp_to_color_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP:
	msg.set_qos_tc_to_queue_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
	msg.set_qos_tc_and_color_to_dot1p_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
	msg.set_qos_tc_and_color_to_dscp_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_SWITCH_SHELL_ENABLE:
	msg.set_switch_shell_enable(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_SWITCH_PROFILE_ID:
	msg.set_switch_profile_id(attr_list[i].value.u32);
	break;
  case SAI_SWITCH_ATTR_SWITCH_HARDWARE_INFO:
	msg.mutable_switch_hardware_info()->Add(attr_list[i].value.s8list.list, attr_list[i].value.s8list.list + attr_list[i].value.s8list.count);
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_PATH_NAME:
	msg.mutable_firmware_path_name()->Add(attr_list[i].value.s8list.list, attr_list[i].value.s8list.list + attr_list[i].value.s8list.count);
	break;
  case SAI_SWITCH_ATTR_INIT_SWITCH:
	msg.set_init_switch(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_PORT_STATE_CHANGE_NOTIFY:
	port_state = std::make_unique<PortStateReactor>( switch_, reinterpret_cast<sai_port_state_change_notification_fn>(attr_list[i].value.ptr) );
	break;
  case SAI_SWITCH_ATTR_FAST_API_ENABLE:
	msg.set_fast_api_enable(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_MIRROR_TC:
	msg.set_mirror_tc(attr_list[i].value.u8);
	break;
  case SAI_SWITCH_ATTR_PFC_DLR_PACKET_ACTION:
	msg.set_pfc_dlr_packet_action(convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_TPID_OUTER_VLAN:
	msg.set_tpid_outer_vlan(attr_list[i].value.u16);
	break;
  case SAI_SWITCH_ATTR_TPID_INNER_VLAN:
	msg.set_tpid_inner_vlan(attr_list[i].value.u16);
	break;
  case SAI_SWITCH_ATTR_CRC_CHECK_ENABLE:
	msg.set_crc_check_enable(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_CRC_RECALCULATION_ENABLE:
	msg.set_crc_recalculation_enable(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_ECN_ECT_THRESHOLD_ENABLE:
	msg.set_ecn_ect_threshold_enable(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_VXLAN_DEFAULT_ROUTER_MAC:
	msg.set_vxlan_default_router_mac(attr_list[i].value.mac, sizeof(attr_list[i].value.mac));
	break;
  case SAI_SWITCH_ATTR_VXLAN_DEFAULT_PORT:
	msg.set_vxlan_default_port(attr_list[i].value.u16);
	break;
  case SAI_SWITCH_ATTR_UNINIT_DATA_PLANE_ON_REMOVAL:
	msg.set_uninit_data_plane_on_removal(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_TAM_OBJECT_ID:
	msg.mutable_tam_object_id()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_PRE_SHUTDOWN:
	msg.set_pre_shutdown(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID:
	msg.set_nat_zone_counter_object_id(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_NAT_ENABLE:
	msg.set_nat_enable(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_HARDWARE_ACCESS_BUS:
	msg.set_hardware_access_bus(convert_sai_switch_hardware_access_bus_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_PLATFROM_CONTEXT:
	msg.set_platfrom_context(attr_list[i].value.u64);
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_BROADCAST:
	msg.set_firmware_download_broadcast(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_LOAD_METHOD:
	msg.set_firmware_load_method(convert_sai_switch_firmware_load_method_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_LOAD_TYPE:
	msg.set_firmware_load_type(convert_sai_switch_firmware_load_type_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_EXECUTE:
	msg.set_firmware_download_execute(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_BROADCAST_STOP:
	msg.set_firmware_broadcast_stop(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_VERIFY_AND_INIT_SWITCH:
	msg.set_firmware_verify_and_init_switch(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_TYPE:
	msg.set_type(convert_sai_switch_type_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_MACSEC_OBJECT_LIST:
	msg.mutable_macsec_object_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
	msg.set_qos_mpls_exp_to_tc_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
	msg.set_qos_mpls_exp_to_color_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
	msg.set_qos_tc_and_color_to_mpls_exp_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_SWITCH_ID:
	msg.set_switch_id(attr_list[i].value.u32);
	break;
  case SAI_SWITCH_ATTR_MAX_SYSTEM_CORES:
	msg.set_max_system_cores(attr_list[i].value.u32);
	break;
  case SAI_SWITCH_ATTR_FAILOVER_CONFIG_MODE:
	msg.set_failover_config_mode(convert_sai_switch_failover_config_mode_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_TUNNEL_OBJECTS_LIST:
	msg.mutable_tunnel_objects_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_PRE_INGRESS_ACL:
	msg.set_pre_ingress_acl(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_SLAVE_MDIO_ADDR_LIST:
	msg.mutable_slave_mdio_addr_list()->Add(attr_list[i].value.u8list.list, attr_list[i].value.u8list.list + attr_list[i].value.u8list.count);
	break;
  case SAI_SWITCH_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
	msg.set_qos_dscp_to_forwarding_class_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
	msg.set_qos_mpls_exp_to_forwarding_class_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_IPSEC_OBJECT_ID:
	msg.set_ipsec_object_id(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_IPSEC_SA_TAG_TPID:
	msg.set_ipsec_sa_tag_tpid(attr_list[i].value.u16);
	break;
  case SAI_SWITCH_ATTR_ECMP_MEMBER_COUNT:
	msg.set_ecmp_member_count(attr_list[i].value.u32);
	break;
  case SAI_SWITCH_ATTR_ARS_PROFILE:
	msg.set_ars_profile(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_ATTR_CREDIT_WD:
	msg.set_credit_wd(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_CREDIT_WD_TIMER:
	msg.set_credit_wd_timer(attr_list[i].value.u32);
	break;
  case SAI_SWITCH_ATTR_SWITCH_ISOLATE:
	msg.set_switch_isolate(attr_list[i].value.booldata);
	break;
  case SAI_SWITCH_ATTR_HOSTIF_OPER_STATUS_UPDATE_MODE:
	msg.set_hostif_oper_status_update_mode(convert_sai_switch_hostif_oper_status_update_mode_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_ATTR_REG_FATAL_SWITCH_ASIC_SDK_HEALTH_CATEGORY:
	msg.mutable_reg_fatal_switch_asic_sdk_health_category()->CopyFrom(convert_list_sai_switch_asic_sdk_health_category_t_to_proto(attr_list[i].value.s32list));
	break;
  case SAI_SWITCH_ATTR_REG_WARNING_SWITCH_ASIC_SDK_HEALTH_CATEGORY:
	msg.mutable_reg_warning_switch_asic_sdk_health_category()->CopyFrom(convert_list_sai_switch_asic_sdk_health_category_t_to_proto(attr_list[i].value.s32list));
	break;
  case SAI_SWITCH_ATTR_REG_NOTICE_SWITCH_ASIC_SDK_HEALTH_CATEGORY:
	msg.mutable_reg_notice_switch_asic_sdk_health_category()->CopyFrom(convert_list_sai_switch_asic_sdk_health_category_t_to_proto(attr_list[i].value.s32list));
	break;
  case SAI_SWITCH_ATTR_POE_DEVICE_LIST:
	msg.mutable_poe_device_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateSwitchTunnelRequest convert_create_switch_tunnel(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateSwitchTunnelRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_TYPE:
	msg.set_tunnel_type(convert_sai_tunnel_type_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
	msg.set_loopback_packet_action(convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_ENCAP_ECN_MODE:
	msg.set_tunnel_encap_ecn_mode(convert_sai_tunnel_encap_ecn_mode_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_TUNNEL_ATTR_ENCAP_MAPPERS:
	msg.mutable_encap_mappers()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_DECAP_ECN_MODE:
	msg.set_tunnel_decap_ecn_mode(convert_sai_tunnel_decap_ecn_mode_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_TUNNEL_ATTR_DECAP_MAPPERS:
	msg.mutable_decap_mappers()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_VXLAN_UDP_SPORT_MODE:
	msg.set_tunnel_vxlan_udp_sport_mode(convert_sai_tunnel_vxlan_udp_sport_mode_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT:
	msg.set_vxlan_udp_sport(attr_list[i].value.u16);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
	msg.set_vxlan_udp_sport_mask(attr_list[i].value.u8);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
	msg.set_encap_qos_tc_and_color_to_dscp_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
	msg.set_encap_qos_tc_to_queue_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
	msg.set_decap_qos_dscp_to_tc_map(attr_list[i].value.oid);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
	msg.set_decap_qos_tc_to_priority_group_map(attr_list[i].value.oid);
	break;
}

}
return msg;
}

sai_status_t l_create_switch(sai_object_id_t *switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateSwitchRequest req = convert_create_switch(attr_count, attr_list);
	lemming::dataplane::sai::CreateSwitchResponse resp;
	grpc::ClientContext context;
	
	
	grpc::Status status = switch_->CreateSwitch(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (switch_id) {
	*switch_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_switch(sai_object_id_t switch_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveSwitchRequest req;
	lemming::dataplane::sai::RemoveSwitchResponse resp;
	grpc::ClientContext context;
	req.set_oid(switch_id); 
	
	grpc::Status status = switch_->RemoveSwitch(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_switch_attribute(sai_object_id_t switch_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetSwitchAttributeRequest req;
	lemming::dataplane::sai::SetSwitchAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(switch_id); 
	
	
	

switch (attr->id) {
  
  case SAI_SWITCH_ATTR_INGRESS_ACL:
	req.set_ingress_acl(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_EGRESS_ACL:
	req.set_egress_acl(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_RESTART_WARM:
	req.set_restart_warm(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_WARM_RECOVER:
	req.set_warm_recover(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_SWITCHING_MODE:
	req.set_switching_mode(convert_sai_switch_switching_mode_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_ATTR_BCAST_CPU_FLOOD_ENABLE:
	req.set_bcast_cpu_flood_enable(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_MCAST_CPU_FLOOD_ENABLE:
	req.set_mcast_cpu_flood_enable(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_SRC_MAC_ADDRESS:
	req.set_src_mac_address(attr->value.mac, sizeof(attr->value.mac));
	break;
  case SAI_SWITCH_ATTR_MAX_LEARNED_ADDRESSES:
	req.set_max_learned_addresses(attr->value.u32);
	break;
  case SAI_SWITCH_ATTR_FDB_AGING_TIME:
	req.set_fdb_aging_time(attr->value.u32);
	break;
  case SAI_SWITCH_ATTR_FDB_UNICAST_MISS_PACKET_ACTION:
	req.set_fdb_unicast_miss_packet_action(convert_sai_packet_action_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_ATTR_FDB_BROADCAST_MISS_PACKET_ACTION:
	req.set_fdb_broadcast_miss_packet_action(convert_sai_packet_action_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_ATTR_FDB_MULTICAST_MISS_PACKET_ACTION:
	req.set_fdb_multicast_miss_packet_action(convert_sai_packet_action_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_ALGORITHM:
	req.set_ecmp_default_hash_algorithm(convert_sai_hash_algorithm_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_SEED:
	req.set_ecmp_default_hash_seed(attr->value.u32);
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_OFFSET:
	req.set_ecmp_default_hash_offset(attr->value.u8);
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_SYMMETRIC_HASH:
	req.set_ecmp_default_symmetric_hash(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_ECMP_HASH_IPV4:
	req.set_ecmp_hash_ipv4(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_ECMP_HASH_IPV4_IN_IPV4:
	req.set_ecmp_hash_ipv4_in_ipv4(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_ECMP_HASH_IPV6:
	req.set_ecmp_hash_ipv6(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM:
	req.set_lag_default_hash_algorithm(convert_sai_hash_algorithm_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_SEED:
	req.set_lag_default_hash_seed(attr->value.u32);
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_OFFSET:
	req.set_lag_default_hash_offset(attr->value.u8);
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH:
	req.set_lag_default_symmetric_hash(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_LAG_HASH_IPV4:
	req.set_lag_hash_ipv4(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_LAG_HASH_IPV4_IN_IPV4:
	req.set_lag_hash_ipv4_in_ipv4(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_LAG_HASH_IPV6:
	req.set_lag_hash_ipv6(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_COUNTER_REFRESH_INTERVAL:
	req.set_counter_refresh_interval(attr->value.u32);
	break;
  case SAI_SWITCH_ATTR_QOS_DEFAULT_TC:
	req.set_qos_default_tc(attr->value.u8);
	break;
  case SAI_SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP:
	req.set_qos_dot1p_to_tc_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP:
	req.set_qos_dot1p_to_color_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_DSCP_TO_TC_MAP:
	req.set_qos_dscp_to_tc_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_DSCP_TO_COLOR_MAP:
	req.set_qos_dscp_to_color_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP:
	req.set_qos_tc_to_queue_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
	req.set_qos_tc_and_color_to_dot1p_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
	req.set_qos_tc_and_color_to_dscp_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_SWITCH_SHELL_ENABLE:
	req.set_switch_shell_enable(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_PORT_STATE_CHANGE_NOTIFY:
	port_state = std::make_unique<PortStateReactor>( switch_, reinterpret_cast<sai_port_state_change_notification_fn>(attr->value.ptr) );
	break;
  case SAI_SWITCH_ATTR_FAST_API_ENABLE:
	req.set_fast_api_enable(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_MIRROR_TC:
	req.set_mirror_tc(attr->value.u8);
	break;
  case SAI_SWITCH_ATTR_PFC_DLR_PACKET_ACTION:
	req.set_pfc_dlr_packet_action(convert_sai_packet_action_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_ATTR_TPID_OUTER_VLAN:
	req.set_tpid_outer_vlan(attr->value.u16);
	break;
  case SAI_SWITCH_ATTR_TPID_INNER_VLAN:
	req.set_tpid_inner_vlan(attr->value.u16);
	break;
  case SAI_SWITCH_ATTR_CRC_CHECK_ENABLE:
	req.set_crc_check_enable(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_CRC_RECALCULATION_ENABLE:
	req.set_crc_recalculation_enable(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_ECN_ECT_THRESHOLD_ENABLE:
	req.set_ecn_ect_threshold_enable(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_VXLAN_DEFAULT_ROUTER_MAC:
	req.set_vxlan_default_router_mac(attr->value.mac, sizeof(attr->value.mac));
	break;
  case SAI_SWITCH_ATTR_VXLAN_DEFAULT_PORT:
	req.set_vxlan_default_port(attr->value.u16);
	break;
  case SAI_SWITCH_ATTR_UNINIT_DATA_PLANE_ON_REMOVAL:
	req.set_uninit_data_plane_on_removal(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_TAM_OBJECT_ID:
	req.mutable_tam_object_id()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_PRE_SHUTDOWN:
	req.set_pre_shutdown(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID:
	req.set_nat_zone_counter_object_id(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_NAT_ENABLE:
	req.set_nat_enable(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_EXECUTE:
	req.set_firmware_download_execute(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_BROADCAST_STOP:
	req.set_firmware_broadcast_stop(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_VERIFY_AND_INIT_SWITCH:
	req.set_firmware_verify_and_init_switch(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_MACSEC_OBJECT_LIST:
	req.mutable_macsec_object_list()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
	req.set_qos_mpls_exp_to_tc_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
	req.set_qos_mpls_exp_to_color_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
	req.set_qos_tc_and_color_to_mpls_exp_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_FAILOVER_CONFIG_MODE:
	req.set_failover_config_mode(convert_sai_switch_failover_config_mode_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_ATTR_TUNNEL_OBJECTS_LIST:
	req.mutable_tunnel_objects_list()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_PRE_INGRESS_ACL:
	req.set_pre_ingress_acl(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
	req.set_qos_dscp_to_forwarding_class_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
	req.set_qos_mpls_exp_to_forwarding_class_map(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_IPSEC_OBJECT_ID:
	req.set_ipsec_object_id(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_IPSEC_SA_TAG_TPID:
	req.set_ipsec_sa_tag_tpid(attr->value.u16);
	break;
  case SAI_SWITCH_ATTR_ECMP_MEMBER_COUNT:
	req.set_ecmp_member_count(attr->value.u32);
	break;
  case SAI_SWITCH_ATTR_ARS_PROFILE:
	req.set_ars_profile(attr->value.oid);
	break;
  case SAI_SWITCH_ATTR_CREDIT_WD:
	req.set_credit_wd(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_CREDIT_WD_TIMER:
	req.set_credit_wd_timer(attr->value.u32);
	break;
  case SAI_SWITCH_ATTR_SWITCH_ISOLATE:
	req.set_switch_isolate(attr->value.booldata);
	break;
  case SAI_SWITCH_ATTR_HOSTIF_OPER_STATUS_UPDATE_MODE:
	req.set_hostif_oper_status_update_mode(convert_sai_switch_hostif_oper_status_update_mode_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_ATTR_REG_FATAL_SWITCH_ASIC_SDK_HEALTH_CATEGORY:
	req.mutable_reg_fatal_switch_asic_sdk_health_category()->CopyFrom(convert_list_sai_switch_asic_sdk_health_category_t_to_proto(attr->value.s32list));
	break;
  case SAI_SWITCH_ATTR_REG_WARNING_SWITCH_ASIC_SDK_HEALTH_CATEGORY:
	req.mutable_reg_warning_switch_asic_sdk_health_category()->CopyFrom(convert_list_sai_switch_asic_sdk_health_category_t_to_proto(attr->value.s32list));
	break;
  case SAI_SWITCH_ATTR_REG_NOTICE_SWITCH_ASIC_SDK_HEALTH_CATEGORY:
	req.mutable_reg_notice_switch_asic_sdk_health_category()->CopyFrom(convert_list_sai_switch_asic_sdk_health_category_t_to_proto(attr->value.s32list));
	break;
  case SAI_SWITCH_ATTR_POE_DEVICE_LIST:
	req.mutable_poe_device_list()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
}

	grpc::Status status = switch_->SetSwitchAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}










sai_status_t l_get_switch_attribute(sai_object_id_t switch_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetSwitchAttributeRequest req;
	lemming::dataplane::sai::GetSwitchAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(switch_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_switch_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = switch_->GetSwitchAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		
		

switch (attr_list[i].id) {
  
  case SAI_SWITCH_ATTR_NUMBER_OF_ACTIVE_PORTS:
	 attr_list[i].value.u32 =   resp.attr().number_of_active_ports();
	break;
  case SAI_SWITCH_ATTR_MAX_NUMBER_OF_SUPPORTED_PORTS:
	 attr_list[i].value.u32 =   resp.attr().max_number_of_supported_ports();
	break;
  case SAI_SWITCH_ATTR_PORT_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().port_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_PORT_MAX_MTU:
	 attr_list[i].value.u32 =   resp.attr().port_max_mtu();
	break;
  case SAI_SWITCH_ATTR_CPU_PORT:
	 attr_list[i].value.oid =   resp.attr().cpu_port();
	break;
  case SAI_SWITCH_ATTR_MAX_VIRTUAL_ROUTERS:
	 attr_list[i].value.u32 =   resp.attr().max_virtual_routers();
	break;
  case SAI_SWITCH_ATTR_FDB_TABLE_SIZE:
	 attr_list[i].value.u32 =   resp.attr().fdb_table_size();
	break;
  case SAI_SWITCH_ATTR_L3_NEIGHBOR_TABLE_SIZE:
	 attr_list[i].value.u32 =   resp.attr().l3_neighbor_table_size();
	break;
  case SAI_SWITCH_ATTR_L3_ROUTE_TABLE_SIZE:
	 attr_list[i].value.u32 =   resp.attr().l3_route_table_size();
	break;
  case SAI_SWITCH_ATTR_LAG_MEMBERS:
	 attr_list[i].value.u32 =   resp.attr().lag_members();
	break;
  case SAI_SWITCH_ATTR_NUMBER_OF_LAGS:
	 attr_list[i].value.u32 =   resp.attr().number_of_lags();
	break;
  case SAI_SWITCH_ATTR_ECMP_MEMBERS:
	 attr_list[i].value.u32 =   resp.attr().ecmp_members();
	break;
  case SAI_SWITCH_ATTR_NUMBER_OF_ECMP_GROUPS:
	 attr_list[i].value.u32 =   resp.attr().number_of_ecmp_groups();
	break;
  case SAI_SWITCH_ATTR_NUMBER_OF_UNICAST_QUEUES:
	 attr_list[i].value.u32 =   resp.attr().number_of_unicast_queues();
	break;
  case SAI_SWITCH_ATTR_NUMBER_OF_MULTICAST_QUEUES:
	 attr_list[i].value.u32 =   resp.attr().number_of_multicast_queues();
	break;
  case SAI_SWITCH_ATTR_NUMBER_OF_QUEUES:
	 attr_list[i].value.u32 =   resp.attr().number_of_queues();
	break;
  case SAI_SWITCH_ATTR_NUMBER_OF_CPU_QUEUES:
	 attr_list[i].value.u32 =   resp.attr().number_of_cpu_queues();
	break;
  case SAI_SWITCH_ATTR_ON_LINK_ROUTE_SUPPORTED:
	 attr_list[i].value.booldata =   resp.attr().on_link_route_supported();
	break;
  case SAI_SWITCH_ATTR_OPER_STATUS:
	 attr_list[i].value.s32 =  convert_sai_switch_oper_status_t_to_sai(resp.attr().oper_status());
	break;
  case SAI_SWITCH_ATTR_MAX_NUMBER_OF_TEMP_SENSORS:
	 attr_list[i].value.u8 =   resp.attr().max_number_of_temp_sensors();
	break;
  case SAI_SWITCH_ATTR_TEMP_LIST:
	copy_list(attr_list[i].value.s32list.list, resp.attr().temp_list(), &attr_list[i].value.s32list.count);
	break;
  case SAI_SWITCH_ATTR_ACL_TABLE_MINIMUM_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().acl_table_minimum_priority();
	break;
  case SAI_SWITCH_ATTR_ACL_TABLE_MAXIMUM_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().acl_table_maximum_priority();
	break;
  case SAI_SWITCH_ATTR_ACL_ENTRY_MINIMUM_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().acl_entry_minimum_priority();
	break;
  case SAI_SWITCH_ATTR_ACL_ENTRY_MAXIMUM_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().acl_entry_maximum_priority();
	break;
  case SAI_SWITCH_ATTR_ACL_TABLE_GROUP_MINIMUM_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().acl_table_group_minimum_priority();
	break;
  case SAI_SWITCH_ATTR_ACL_TABLE_GROUP_MAXIMUM_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().acl_table_group_maximum_priority();
	break;
  case SAI_SWITCH_ATTR_DEFAULT_VLAN_ID:
	 attr_list[i].value.oid =   resp.attr().default_vlan_id();
	break;
  case SAI_SWITCH_ATTR_DEFAULT_STP_INST_ID:
	 attr_list[i].value.oid =   resp.attr().default_stp_inst_id();
	break;
  case SAI_SWITCH_ATTR_MAX_STP_INSTANCE:
	 attr_list[i].value.u32 =   resp.attr().max_stp_instance();
	break;
  case SAI_SWITCH_ATTR_DEFAULT_VIRTUAL_ROUTER_ID:
	 attr_list[i].value.oid =   resp.attr().default_virtual_router_id();
	break;
  case SAI_SWITCH_ATTR_DEFAULT_OVERRIDE_VIRTUAL_ROUTER_ID:
	 attr_list[i].value.oid =   resp.attr().default_override_virtual_router_id();
	break;
  case SAI_SWITCH_ATTR_DEFAULT_1Q_BRIDGE_ID:
	 attr_list[i].value.oid =   resp.attr().default_1q_bridge_id();
	break;
  case SAI_SWITCH_ATTR_INGRESS_ACL:
	 attr_list[i].value.oid =   resp.attr().ingress_acl();
	break;
  case SAI_SWITCH_ATTR_EGRESS_ACL:
	 attr_list[i].value.oid =   resp.attr().egress_acl();
	break;
  case SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_TRAFFIC_CLASSES:
	 attr_list[i].value.u8 =   resp.attr().qos_max_number_of_traffic_classes();
	break;
  case SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_SCHEDULER_GROUP_HIERARCHY_LEVELS:
	 attr_list[i].value.u32 =   resp.attr().qos_max_number_of_scheduler_group_hierarchy_levels();
	break;
  case SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_SCHEDULER_GROUPS_PER_HIERARCHY_LEVEL:
	copy_list(attr_list[i].value.u32list.list, resp.attr().qos_max_number_of_scheduler_groups_per_hierarchy_level(), &attr_list[i].value.u32list.count);
	break;
  case SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_CHILDS_PER_SCHEDULER_GROUP:
	 attr_list[i].value.u32 =   resp.attr().qos_max_number_of_childs_per_scheduler_group();
	break;
  case SAI_SWITCH_ATTR_TOTAL_BUFFER_SIZE:
	 attr_list[i].value.u64 =   resp.attr().total_buffer_size();
	break;
  case SAI_SWITCH_ATTR_INGRESS_BUFFER_POOL_NUM:
	 attr_list[i].value.u32 =   resp.attr().ingress_buffer_pool_num();
	break;
  case SAI_SWITCH_ATTR_EGRESS_BUFFER_POOL_NUM:
	 attr_list[i].value.u32 =   resp.attr().egress_buffer_pool_num();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_IPV4_ROUTE_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_ipv4_route_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_IPV6_ROUTE_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_ipv6_route_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_IPV4_NEXTHOP_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_ipv4_nexthop_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_IPV6_NEXTHOP_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_ipv6_nexthop_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_IPV4_NEIGHBOR_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_ipv4_neighbor_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_IPV6_NEIGHBOR_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_ipv6_neighbor_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_next_hop_group_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_MEMBER_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_next_hop_group_member_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_FDB_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_fdb_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_L2MC_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_l2mc_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_IPMC_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_ipmc_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_SNAT_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_snat_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_DNAT_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_dnat_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_DOUBLE_NAT_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_double_nat_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_MY_SID_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_my_sid_entry();
	break;
  case SAI_SWITCH_ATTR_DEFAULT_TRAP_GROUP:
	 attr_list[i].value.oid =   resp.attr().default_trap_group();
	break;
  case SAI_SWITCH_ATTR_ECMP_HASH:
	 attr_list[i].value.oid =   resp.attr().ecmp_hash();
	break;
  case SAI_SWITCH_ATTR_LAG_HASH:
	 attr_list[i].value.oid =   resp.attr().lag_hash();
	break;
  case SAI_SWITCH_ATTR_RESTART_WARM:
	 attr_list[i].value.booldata =   resp.attr().restart_warm();
	break;
  case SAI_SWITCH_ATTR_WARM_RECOVER:
	 attr_list[i].value.booldata =   resp.attr().warm_recover();
	break;
  case SAI_SWITCH_ATTR_RESTART_TYPE:
	 attr_list[i].value.s32 =  convert_sai_switch_restart_type_t_to_sai(resp.attr().restart_type());
	break;
  case SAI_SWITCH_ATTR_MIN_PLANNED_RESTART_INTERVAL:
	 attr_list[i].value.u32 =   resp.attr().min_planned_restart_interval();
	break;
  case SAI_SWITCH_ATTR_NV_STORAGE_SIZE:
	 attr_list[i].value.u64 =   resp.attr().nv_storage_size();
	break;
  case SAI_SWITCH_ATTR_MAX_ACL_ACTION_COUNT:
	 attr_list[i].value.u32 =   resp.attr().max_acl_action_count();
	break;
  case SAI_SWITCH_ATTR_MAX_ACL_RANGE_COUNT:
	 attr_list[i].value.u32 =   resp.attr().max_acl_range_count();
	break;
  case SAI_SWITCH_ATTR_ACL_CAPABILITY:
	convert_to_acl_capability(attr_list[i].value.aclcapability, resp.attr().acl_capability());
	break;
  case SAI_SWITCH_ATTR_MCAST_SNOOPING_CAPABILITY:
	 attr_list[i].value.s32 =  convert_sai_switch_mcast_snooping_capability_t_to_sai(resp.attr().mcast_snooping_capability());
	break;
  case SAI_SWITCH_ATTR_SWITCHING_MODE:
	 attr_list[i].value.s32 =  convert_sai_switch_switching_mode_t_to_sai(resp.attr().switching_mode());
	break;
  case SAI_SWITCH_ATTR_BCAST_CPU_FLOOD_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().bcast_cpu_flood_enable();
	break;
  case SAI_SWITCH_ATTR_MCAST_CPU_FLOOD_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().mcast_cpu_flood_enable();
	break;
  case SAI_SWITCH_ATTR_SRC_MAC_ADDRESS:
	memcpy(attr_list[i].value.mac, resp.attr().src_mac_address().data(), sizeof(sai_mac_t));
	break;
  case SAI_SWITCH_ATTR_MAX_LEARNED_ADDRESSES:
	 attr_list[i].value.u32 =   resp.attr().max_learned_addresses();
	break;
  case SAI_SWITCH_ATTR_FDB_AGING_TIME:
	 attr_list[i].value.u32 =   resp.attr().fdb_aging_time();
	break;
  case SAI_SWITCH_ATTR_FDB_UNICAST_MISS_PACKET_ACTION:
	 attr_list[i].value.s32 =  convert_sai_packet_action_t_to_sai(resp.attr().fdb_unicast_miss_packet_action());
	break;
  case SAI_SWITCH_ATTR_FDB_BROADCAST_MISS_PACKET_ACTION:
	 attr_list[i].value.s32 =  convert_sai_packet_action_t_to_sai(resp.attr().fdb_broadcast_miss_packet_action());
	break;
  case SAI_SWITCH_ATTR_FDB_MULTICAST_MISS_PACKET_ACTION:
	 attr_list[i].value.s32 =  convert_sai_packet_action_t_to_sai(resp.attr().fdb_multicast_miss_packet_action());
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_ALGORITHM:
	 attr_list[i].value.s32 =  convert_sai_hash_algorithm_t_to_sai(resp.attr().ecmp_default_hash_algorithm());
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_SEED:
	 attr_list[i].value.u32 =   resp.attr().ecmp_default_hash_seed();
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_OFFSET:
	 attr_list[i].value.u8 =   resp.attr().ecmp_default_hash_offset();
	break;
  case SAI_SWITCH_ATTR_ECMP_DEFAULT_SYMMETRIC_HASH:
	 attr_list[i].value.booldata =   resp.attr().ecmp_default_symmetric_hash();
	break;
  case SAI_SWITCH_ATTR_ECMP_HASH_IPV4:
	 attr_list[i].value.oid =   resp.attr().ecmp_hash_ipv4();
	break;
  case SAI_SWITCH_ATTR_ECMP_HASH_IPV4_IN_IPV4:
	 attr_list[i].value.oid =   resp.attr().ecmp_hash_ipv4_in_ipv4();
	break;
  case SAI_SWITCH_ATTR_ECMP_HASH_IPV6:
	 attr_list[i].value.oid =   resp.attr().ecmp_hash_ipv6();
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM:
	 attr_list[i].value.s32 =  convert_sai_hash_algorithm_t_to_sai(resp.attr().lag_default_hash_algorithm());
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_SEED:
	 attr_list[i].value.u32 =   resp.attr().lag_default_hash_seed();
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_OFFSET:
	 attr_list[i].value.u8 =   resp.attr().lag_default_hash_offset();
	break;
  case SAI_SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH:
	 attr_list[i].value.booldata =   resp.attr().lag_default_symmetric_hash();
	break;
  case SAI_SWITCH_ATTR_LAG_HASH_IPV4:
	 attr_list[i].value.oid =   resp.attr().lag_hash_ipv4();
	break;
  case SAI_SWITCH_ATTR_LAG_HASH_IPV4_IN_IPV4:
	 attr_list[i].value.oid =   resp.attr().lag_hash_ipv4_in_ipv4();
	break;
  case SAI_SWITCH_ATTR_LAG_HASH_IPV6:
	 attr_list[i].value.oid =   resp.attr().lag_hash_ipv6();
	break;
  case SAI_SWITCH_ATTR_COUNTER_REFRESH_INTERVAL:
	 attr_list[i].value.u32 =   resp.attr().counter_refresh_interval();
	break;
  case SAI_SWITCH_ATTR_QOS_DEFAULT_TC:
	 attr_list[i].value.u8 =   resp.attr().qos_default_tc();
	break;
  case SAI_SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_dot1p_to_tc_map();
	break;
  case SAI_SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_dot1p_to_color_map();
	break;
  case SAI_SWITCH_ATTR_QOS_DSCP_TO_TC_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_dscp_to_tc_map();
	break;
  case SAI_SWITCH_ATTR_QOS_DSCP_TO_COLOR_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_dscp_to_color_map();
	break;
  case SAI_SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_tc_to_queue_map();
	break;
  case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_tc_and_color_to_dot1p_map();
	break;
  case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_tc_and_color_to_dscp_map();
	break;
  case SAI_SWITCH_ATTR_SWITCH_SHELL_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().switch_shell_enable();
	break;
  case SAI_SWITCH_ATTR_SWITCH_PROFILE_ID:
	 attr_list[i].value.u32 =   resp.attr().switch_profile_id();
	break;
  case SAI_SWITCH_ATTR_SWITCH_HARDWARE_INFO:
	copy_list(attr_list[i].value.s8list.list, resp.attr().switch_hardware_info(), &attr_list[i].value.s8list.count);
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_PATH_NAME:
	copy_list(attr_list[i].value.s8list.list, resp.attr().firmware_path_name(), &attr_list[i].value.s8list.count);
	break;
  case SAI_SWITCH_ATTR_INIT_SWITCH:
	 attr_list[i].value.booldata =   resp.attr().init_switch();
	break;
  case SAI_SWITCH_ATTR_FAST_API_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().fast_api_enable();
	break;
  case SAI_SWITCH_ATTR_MIRROR_TC:
	 attr_list[i].value.u8 =   resp.attr().mirror_tc();
	break;
  case SAI_SWITCH_ATTR_ACL_STAGE_INGRESS:
	convert_to_acl_capability(attr_list[i].value.aclcapability, resp.attr().acl_stage_ingress());
	break;
  case SAI_SWITCH_ATTR_ACL_STAGE_EGRESS:
	convert_to_acl_capability(attr_list[i].value.aclcapability, resp.attr().acl_stage_egress());
	break;
  case SAI_SWITCH_ATTR_SRV6_MAX_SID_DEPTH:
	 attr_list[i].value.u32 =   resp.attr().srv6_max_sid_depth();
	break;
  case SAI_SWITCH_ATTR_SRV6_TLV_TYPE:
	convert_list_sai_tlv_type_t_to_sai(attr_list[i].value.s32list.list, resp.attr().srv6_tlv_type(), &attr_list[i].value.s32list.count);
	break;
  case SAI_SWITCH_ATTR_QOS_NUM_LOSSLESS_QUEUES:
	 attr_list[i].value.u32 =   resp.attr().qos_num_lossless_queues();
	break;
  case SAI_SWITCH_ATTR_PFC_DLR_PACKET_ACTION:
	 attr_list[i].value.s32 =  convert_sai_packet_action_t_to_sai(resp.attr().pfc_dlr_packet_action());
	break;
  case SAI_SWITCH_ATTR_SUPPORTED_PROTECTED_OBJECT_TYPE:
	convert_list_sai_object_type_t_to_sai(attr_list[i].value.s32list.list, resp.attr().supported_protected_object_type(), &attr_list[i].value.s32list.count);
	break;
  case SAI_SWITCH_ATTR_TPID_OUTER_VLAN:
	 attr_list[i].value.u16 =   resp.attr().tpid_outer_vlan();
	break;
  case SAI_SWITCH_ATTR_TPID_INNER_VLAN:
	 attr_list[i].value.u16 =   resp.attr().tpid_inner_vlan();
	break;
  case SAI_SWITCH_ATTR_CRC_CHECK_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().crc_check_enable();
	break;
  case SAI_SWITCH_ATTR_CRC_RECALCULATION_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().crc_recalculation_enable();
	break;
  case SAI_SWITCH_ATTR_NUMBER_OF_BFD_SESSION:
	 attr_list[i].value.u32 =   resp.attr().number_of_bfd_session();
	break;
  case SAI_SWITCH_ATTR_MAX_BFD_SESSION:
	 attr_list[i].value.u32 =   resp.attr().max_bfd_session();
	break;
  case SAI_SWITCH_ATTR_SUPPORTED_IPV4_BFD_SESSION_OFFLOAD_TYPE:
	convert_list_sai_bfd_session_offload_type_t_to_sai(attr_list[i].value.s32list.list, resp.attr().supported_ipv4_bfd_session_offload_type(), &attr_list[i].value.s32list.count);
	break;
  case SAI_SWITCH_ATTR_SUPPORTED_IPV6_BFD_SESSION_OFFLOAD_TYPE:
	convert_list_sai_bfd_session_offload_type_t_to_sai(attr_list[i].value.s32list.list, resp.attr().supported_ipv6_bfd_session_offload_type(), &attr_list[i].value.s32list.count);
	break;
  case SAI_SWITCH_ATTR_MIN_BFD_RX:
	 attr_list[i].value.u32 =   resp.attr().min_bfd_rx();
	break;
  case SAI_SWITCH_ATTR_MIN_BFD_TX:
	 attr_list[i].value.u32 =   resp.attr().min_bfd_tx();
	break;
  case SAI_SWITCH_ATTR_ECN_ECT_THRESHOLD_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().ecn_ect_threshold_enable();
	break;
  case SAI_SWITCH_ATTR_VXLAN_DEFAULT_ROUTER_MAC:
	memcpy(attr_list[i].value.mac, resp.attr().vxlan_default_router_mac().data(), sizeof(sai_mac_t));
	break;
  case SAI_SWITCH_ATTR_VXLAN_DEFAULT_PORT:
	 attr_list[i].value.u16 =   resp.attr().vxlan_default_port();
	break;
  case SAI_SWITCH_ATTR_MAX_MIRROR_SESSION:
	 attr_list[i].value.u32 =   resp.attr().max_mirror_session();
	break;
  case SAI_SWITCH_ATTR_MAX_SAMPLED_MIRROR_SESSION:
	 attr_list[i].value.u32 =   resp.attr().max_sampled_mirror_session();
	break;
  case SAI_SWITCH_ATTR_SUPPORTED_EXTENDED_STATS_MODE:
	convert_list_sai_stats_mode_t_to_sai(attr_list[i].value.s32list.list, resp.attr().supported_extended_stats_mode(), &attr_list[i].value.s32list.count);
	break;
  case SAI_SWITCH_ATTR_UNINIT_DATA_PLANE_ON_REMOVAL:
	 attr_list[i].value.booldata =   resp.attr().uninit_data_plane_on_removal();
	break;
  case SAI_SWITCH_ATTR_TAM_OBJECT_ID:
	copy_list(attr_list[i].value.objlist.list, resp.attr().tam_object_id(), &attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_SUPPORTED_OBJECT_TYPE_LIST:
	convert_list_sai_object_type_t_to_sai(attr_list[i].value.s32list.list, resp.attr().supported_object_type_list(), &attr_list[i].value.s32list.count);
	break;
  case SAI_SWITCH_ATTR_PRE_SHUTDOWN:
	 attr_list[i].value.booldata =   resp.attr().pre_shutdown();
	break;
  case SAI_SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID:
	 attr_list[i].value.oid =   resp.attr().nat_zone_counter_object_id();
	break;
  case SAI_SWITCH_ATTR_NAT_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().nat_enable();
	break;
  case SAI_SWITCH_ATTR_HARDWARE_ACCESS_BUS:
	 attr_list[i].value.s32 =  convert_sai_switch_hardware_access_bus_t_to_sai(resp.attr().hardware_access_bus());
	break;
  case SAI_SWITCH_ATTR_PLATFROM_CONTEXT:
	 attr_list[i].value.u64 =   resp.attr().platfrom_context();
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_BROADCAST:
	 attr_list[i].value.booldata =   resp.attr().firmware_download_broadcast();
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_LOAD_METHOD:
	 attr_list[i].value.s32 =  convert_sai_switch_firmware_load_method_t_to_sai(resp.attr().firmware_load_method());
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_LOAD_TYPE:
	 attr_list[i].value.s32 =  convert_sai_switch_firmware_load_type_t_to_sai(resp.attr().firmware_load_type());
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_EXECUTE:
	 attr_list[i].value.booldata =   resp.attr().firmware_download_execute();
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_BROADCAST_STOP:
	 attr_list[i].value.booldata =   resp.attr().firmware_broadcast_stop();
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_VERIFY_AND_INIT_SWITCH:
	 attr_list[i].value.booldata =   resp.attr().firmware_verify_and_init_switch();
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_STATUS:
	 attr_list[i].value.booldata =   resp.attr().firmware_status();
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_MAJOR_VERSION:
	 attr_list[i].value.u32 =   resp.attr().firmware_major_version();
	break;
  case SAI_SWITCH_ATTR_FIRMWARE_MINOR_VERSION:
	 attr_list[i].value.u32 =   resp.attr().firmware_minor_version();
	break;
  case SAI_SWITCH_ATTR_PORT_CONNECTOR_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().port_connector_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_PROPOGATE_PORT_STATE_FROM_LINE_TO_SYSTEM_PORT_SUPPORT:
	 attr_list[i].value.booldata =   resp.attr().propogate_port_state_from_line_to_system_port_support();
	break;
  case SAI_SWITCH_ATTR_TYPE:
	 attr_list[i].value.s32 =  convert_sai_switch_type_t_to_sai(resp.attr().type());
	break;
  case SAI_SWITCH_ATTR_MACSEC_OBJECT_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().macsec_object_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_mpls_exp_to_tc_map();
	break;
  case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_mpls_exp_to_color_map();
	break;
  case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_tc_and_color_to_mpls_exp_map();
	break;
  case SAI_SWITCH_ATTR_SWITCH_ID:
	 attr_list[i].value.u32 =   resp.attr().switch_id();
	break;
  case SAI_SWITCH_ATTR_MAX_SYSTEM_CORES:
	 attr_list[i].value.u32 =   resp.attr().max_system_cores();
	break;
  case SAI_SWITCH_ATTR_NUMBER_OF_SYSTEM_PORTS:
	 attr_list[i].value.u32 =   resp.attr().number_of_system_ports();
	break;
  case SAI_SWITCH_ATTR_SYSTEM_PORT_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().system_port_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_NUMBER_OF_FABRIC_PORTS:
	 attr_list[i].value.u32 =   resp.attr().number_of_fabric_ports();
	break;
  case SAI_SWITCH_ATTR_FABRIC_PORT_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().fabric_port_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_PACKET_DMA_MEMORY_POOL_SIZE:
	 attr_list[i].value.u32 =   resp.attr().packet_dma_memory_pool_size();
	break;
  case SAI_SWITCH_ATTR_FAILOVER_CONFIG_MODE:
	 attr_list[i].value.s32 =  convert_sai_switch_failover_config_mode_t_to_sai(resp.attr().failover_config_mode());
	break;
  case SAI_SWITCH_ATTR_SUPPORTED_FAILOVER_MODE:
	 attr_list[i].value.booldata =   resp.attr().supported_failover_mode();
	break;
  case SAI_SWITCH_ATTR_TUNNEL_OBJECTS_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().tunnel_objects_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_PACKET_AVAILABLE_DMA_MEMORY_POOL_SIZE:
	 attr_list[i].value.u32 =   resp.attr().packet_available_dma_memory_pool_size();
	break;
  case SAI_SWITCH_ATTR_PRE_INGRESS_ACL:
	 attr_list[i].value.oid =   resp.attr().pre_ingress_acl();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_SNAPT_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_snapt_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_DNAPT_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_dnapt_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_DOUBLE_NAPT_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_double_napt_entry();
	break;
  case SAI_SWITCH_ATTR_SLAVE_MDIO_ADDR_LIST:
	copy_list(attr_list[i].value.u8list.list, resp.attr().slave_mdio_addr_list(), &attr_list[i].value.u8list.count);
	break;
  case SAI_SWITCH_ATTR_MY_MAC_TABLE_MINIMUM_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().my_mac_table_minimum_priority();
	break;
  case SAI_SWITCH_ATTR_MY_MAC_TABLE_MAXIMUM_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().my_mac_table_maximum_priority();
	break;
  case SAI_SWITCH_ATTR_MY_MAC_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().my_mac_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_ATTR_INSTALLED_MY_MAC_ENTRIES:
	 attr_list[i].value.u32 =   resp.attr().installed_my_mac_entries();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_MY_MAC_ENTRIES:
	 attr_list[i].value.u32 =   resp.attr().available_my_mac_entries();
	break;
  case SAI_SWITCH_ATTR_MAX_NUMBER_OF_FORWARDING_CLASSES:
	 attr_list[i].value.u8 =   resp.attr().max_number_of_forwarding_classes();
	break;
  case SAI_SWITCH_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_dscp_to_forwarding_class_map();
	break;
  case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
	 attr_list[i].value.oid =   resp.attr().qos_mpls_exp_to_forwarding_class_map();
	break;
  case SAI_SWITCH_ATTR_IPSEC_OBJECT_ID:
	 attr_list[i].value.oid =   resp.attr().ipsec_object_id();
	break;
  case SAI_SWITCH_ATTR_IPSEC_SA_TAG_TPID:
	 attr_list[i].value.u16 =   resp.attr().ipsec_sa_tag_tpid();
	break;
  case SAI_SWITCH_ATTR_MAX_ECMP_MEMBER_COUNT:
	 attr_list[i].value.u32 =   resp.attr().max_ecmp_member_count();
	break;
  case SAI_SWITCH_ATTR_ECMP_MEMBER_COUNT:
	 attr_list[i].value.u32 =   resp.attr().ecmp_member_count();
	break;
  case SAI_SWITCH_ATTR_ARS_PROFILE:
	 attr_list[i].value.oid =   resp.attr().ars_profile();
	break;
  case SAI_SWITCH_ATTR_ACL_STAGE_POST_INGRESS:
	convert_to_acl_capability(attr_list[i].value.aclcapability, resp.attr().acl_stage_post_ingress());
	break;
  case SAI_SWITCH_ATTR_CREDIT_WD:
	 attr_list[i].value.booldata =   resp.attr().credit_wd();
	break;
  case SAI_SWITCH_ATTR_CREDIT_WD_TIMER:
	 attr_list[i].value.u32 =   resp.attr().credit_wd_timer();
	break;
  case SAI_SWITCH_ATTR_SWITCH_ISOLATE:
	 attr_list[i].value.booldata =   resp.attr().switch_isolate();
	break;
  case SAI_SWITCH_ATTR_HOSTIF_OPER_STATUS_UPDATE_MODE:
	 attr_list[i].value.s32 =  convert_sai_switch_hostif_oper_status_update_mode_t_to_sai(resp.attr().hostif_oper_status_update_mode());
	break;
  case SAI_SWITCH_ATTR_REG_FATAL_SWITCH_ASIC_SDK_HEALTH_CATEGORY:
	convert_list_sai_switch_asic_sdk_health_category_t_to_sai(attr_list[i].value.s32list.list, resp.attr().reg_fatal_switch_asic_sdk_health_category(), &attr_list[i].value.s32list.count);
	break;
  case SAI_SWITCH_ATTR_REG_WARNING_SWITCH_ASIC_SDK_HEALTH_CATEGORY:
	convert_list_sai_switch_asic_sdk_health_category_t_to_sai(attr_list[i].value.s32list.list, resp.attr().reg_warning_switch_asic_sdk_health_category(), &attr_list[i].value.s32list.count);
	break;
  case SAI_SWITCH_ATTR_REG_NOTICE_SWITCH_ASIC_SDK_HEALTH_CATEGORY:
	convert_list_sai_switch_asic_sdk_health_category_t_to_sai(attr_list[i].value.s32list.list, resp.attr().reg_notice_switch_asic_sdk_health_category(), &attr_list[i].value.s32list.count);
	break;
  case SAI_SWITCH_ATTR_TUNNEL_TERM_TABLE_ENTRY_MINIMUM_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().tunnel_term_table_entry_minimum_priority();
	break;
  case SAI_SWITCH_ATTR_TUNNEL_TERM_TABLE_ENTRY_MAXIMUM_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().tunnel_term_table_entry_maximum_priority();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_TWAMP_SESSION:
	 attr_list[i].value.u32 =   resp.attr().available_twamp_session();
	break;
  case SAI_SWITCH_ATTR_MAX_TWAMP_SESSION:
	 attr_list[i].value.u32 =   resp.attr().max_twamp_session();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_IP_NEXT_HOP_GROUP_MEMBER_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_ip_next_hop_group_member_entry();
	break;
  case SAI_SWITCH_ATTR_AVAILABLE_SYSTEM_VOQS:
	 attr_list[i].value.u32 =   resp.attr().available_system_voqs();
	break;
  case SAI_SWITCH_ATTR_ACL_STAGE_PRE_INGRESS:
	convert_to_acl_capability(attr_list[i].value.aclcapability, resp.attr().acl_stage_pre_ingress());
	break;
  case SAI_SWITCH_ATTR_POE_DEVICE_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().poe_device_list(), &attr_list[i].value.objlist.count);
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_switch_stats(sai_object_id_t switch_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetSwitchStatsRequest req;
	lemming::dataplane::sai::GetSwitchStatsResponse resp;
	grpc::ClientContext context;
	req.set_oid(switch_id); 
	
	for (uint32_t i = 0; i < number_of_counters; i++) {
		req.add_counter_ids(convert_sai_switch_stat_t_to_proto(counter_ids[i]));
	}
	grpc::Status status = switch_->GetSwitchStats(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < number_of_counters && i < uint32_t(resp.values_size()); i++ ) {
		counters[i] = resp.values(i);
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_switch_stats_ext(sai_object_id_t switch_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_switch_stats(sai_object_id_t switch_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}



sai_status_t l_create_switch_tunnel(sai_object_id_t *switch_tunnel_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateSwitchTunnelRequest req = convert_create_switch_tunnel(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateSwitchTunnelResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = switch_->CreateSwitchTunnel(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (switch_tunnel_id) {
	*switch_tunnel_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_switch_tunnel(sai_object_id_t switch_tunnel_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveSwitchTunnelRequest req;
	lemming::dataplane::sai::RemoveSwitchTunnelResponse resp;
	grpc::ClientContext context;
	req.set_oid(switch_tunnel_id); 
	
	grpc::Status status = switch_->RemoveSwitchTunnel(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_switch_tunnel_attribute(sai_object_id_t switch_tunnel_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetSwitchTunnelAttributeRequest req;
	lemming::dataplane::sai::SetSwitchTunnelAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(switch_tunnel_id); 
	
	
	

switch (attr->id) {
  
  case SAI_SWITCH_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
	req.set_loopback_packet_action(convert_sai_packet_action_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_VXLAN_UDP_SPORT_MODE:
	req.set_tunnel_vxlan_udp_sport_mode(convert_sai_tunnel_vxlan_udp_sport_mode_t_to_proto(attr->value.s32));
	break;
  case SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT:
	req.set_vxlan_udp_sport(attr->value.u16);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
	req.set_vxlan_udp_sport_mask(attr->value.u8);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
	req.set_encap_qos_tc_and_color_to_dscp_map(attr->value.oid);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
	req.set_encap_qos_tc_to_queue_map(attr->value.oid);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
	req.set_decap_qos_dscp_to_tc_map(attr->value.oid);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
	req.set_decap_qos_tc_to_priority_group_map(attr->value.oid);
	break;
}

	grpc::Status status = switch_->SetSwitchTunnelAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_switch_tunnel_attribute(sai_object_id_t switch_tunnel_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetSwitchTunnelAttributeRequest req;
	lemming::dataplane::sai::GetSwitchTunnelAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(switch_tunnel_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_switch_tunnel_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = switch_->GetSwitchTunnelAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		
		

switch (attr_list[i].id) {
  
  case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_TYPE:
	 attr_list[i].value.s32 =  convert_sai_tunnel_type_t_to_sai(resp.attr().tunnel_type());
	break;
  case SAI_SWITCH_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
	 attr_list[i].value.s32 =  convert_sai_packet_action_t_to_sai(resp.attr().loopback_packet_action());
	break;
  case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_ENCAP_ECN_MODE:
	 attr_list[i].value.s32 =  convert_sai_tunnel_encap_ecn_mode_t_to_sai(resp.attr().tunnel_encap_ecn_mode());
	break;
  case SAI_SWITCH_TUNNEL_ATTR_ENCAP_MAPPERS:
	copy_list(attr_list[i].value.objlist.list, resp.attr().encap_mappers(), &attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_DECAP_ECN_MODE:
	 attr_list[i].value.s32 =  convert_sai_tunnel_decap_ecn_mode_t_to_sai(resp.attr().tunnel_decap_ecn_mode());
	break;
  case SAI_SWITCH_TUNNEL_ATTR_DECAP_MAPPERS:
	copy_list(attr_list[i].value.objlist.list, resp.attr().decap_mappers(), &attr_list[i].value.objlist.count);
	break;
  case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_VXLAN_UDP_SPORT_MODE:
	 attr_list[i].value.s32 =  convert_sai_tunnel_vxlan_udp_sport_mode_t_to_sai(resp.attr().tunnel_vxlan_udp_sport_mode());
	break;
  case SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT:
	 attr_list[i].value.u16 =   resp.attr().vxlan_udp_sport();
	break;
  case SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
	 attr_list[i].value.u8 =   resp.attr().vxlan_udp_sport_mask();
	break;
  case SAI_SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
	 attr_list[i].value.oid =   resp.attr().encap_qos_tc_and_color_to_dscp_map();
	break;
  case SAI_SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
	 attr_list[i].value.oid =   resp.attr().encap_qos_tc_to_queue_map();
	break;
  case SAI_SWITCH_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
	 attr_list[i].value.oid =   resp.attr().decap_qos_dscp_to_tc_map();
	break;
  case SAI_SWITCH_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
	 attr_list[i].value.oid =   resp.attr().decap_qos_tc_to_priority_group_map();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}



