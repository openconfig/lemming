
// Copyright 2024 Google LLC
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

#ifndef DATAPLANE_STANDALONE_SAI_ENUM_H_
#define DATAPLANE_STANDALONE_SAI_ENUM_H_

#include "dataplane/proto/sai/acl.pb.h"
#include "dataplane/proto/sai/ars.pb.h"
#include "dataplane/proto/sai/ars_profile.pb.h"
#include "dataplane/proto/sai/bfd.pb.h"
#include "dataplane/proto/sai/bmtor.pb.h"
#include "dataplane/proto/sai/bridge.pb.h"
#include "dataplane/proto/sai/buffer.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/counter.pb.h"
#include "dataplane/proto/sai/dash_acl.pb.h"
#include "dataplane/proto/sai/dash_appliance.pb.h"
#include "dataplane/proto/sai/dash_direction_lookup.pb.h"
#include "dataplane/proto/sai/dash_eni.pb.h"
#include "dataplane/proto/sai/dash_flow.pb.h"
#include "dataplane/proto/sai/dash_ha.pb.h"
#include "dataplane/proto/sai/dash_inbound_routing.pb.h"
#include "dataplane/proto/sai/dash_meter.pb.h"
#include "dataplane/proto/sai/dash_outbound_ca_to_pa.pb.h"
#include "dataplane/proto/sai/dash_outbound_routing.pb.h"
#include "dataplane/proto/sai/dash_pa_validation.pb.h"
#include "dataplane/proto/sai/dash_tunnel.pb.h"
#include "dataplane/proto/sai/dash_vip.pb.h"
#include "dataplane/proto/sai/dash_vnet.pb.h"
#include "dataplane/proto/sai/debug_counter.pb.h"
#include "dataplane/proto/sai/dtel.pb.h"
#include "dataplane/proto/sai/fdb.pb.h"
#include "dataplane/proto/sai/generic_programmable.pb.h"
#include "dataplane/proto/sai/hash.pb.h"
#include "dataplane/proto/sai/hostif.pb.h"
#include "dataplane/proto/sai/icmp_echo.pb.h"
#include "dataplane/proto/sai/ipmc.pb.h"
#include "dataplane/proto/sai/ipmc_group.pb.h"
#include "dataplane/proto/sai/ipsec.pb.h"
#include "dataplane/proto/sai/isolation_group.pb.h"
#include "dataplane/proto/sai/l2mc.pb.h"
#include "dataplane/proto/sai/l2mc_group.pb.h"
#include "dataplane/proto/sai/lag.pb.h"
#include "dataplane/proto/sai/macsec.pb.h"
#include "dataplane/proto/sai/mcast_fdb.pb.h"
#include "dataplane/proto/sai/mirror.pb.h"
#include "dataplane/proto/sai/mpls.pb.h"
#include "dataplane/proto/sai/my_mac.pb.h"
#include "dataplane/proto/sai/nat.pb.h"
#include "dataplane/proto/sai/neighbor.pb.h"
#include "dataplane/proto/sai/next_hop.pb.h"
#include "dataplane/proto/sai/next_hop_group.pb.h"
#include "dataplane/proto/sai/poe.pb.h"
#include "dataplane/proto/sai/policer.pb.h"
#include "dataplane/proto/sai/port.pb.h"
#include "dataplane/proto/sai/qos_map.pb.h"
#include "dataplane/proto/sai/queue.pb.h"
#include "dataplane/proto/sai/route.pb.h"
#include "dataplane/proto/sai/router_interface.pb.h"
#include "dataplane/proto/sai/rpf_group.pb.h"
#include "dataplane/proto/sai/samplepacket.pb.h"
#include "dataplane/proto/sai/scheduler.pb.h"
#include "dataplane/proto/sai/scheduler_group.pb.h"
#include "dataplane/proto/sai/srv6.pb.h"
#include "dataplane/proto/sai/stp.pb.h"
#include "dataplane/proto/sai/switch.pb.h"
#include "dataplane/proto/sai/system_port.pb.h"
#include "dataplane/proto/sai/tam.pb.h"
#include "dataplane/proto/sai/tunnel.pb.h"
#include "dataplane/proto/sai/twamp.pb.h"
#include "dataplane/proto/sai/udf.pb.h"
#include "dataplane/proto/sai/virtual_router.pb.h"
#include "dataplane/proto/sai/vlan.pb.h"
#include "dataplane/proto/sai/wred.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

lemming::dataplane::sai::AclActionType convert_sai_acl_action_type_t_to_proto(
    const sai_int32_t val);
sai_acl_action_type_t convert_sai_acl_action_type_t_to_sai(
    lemming::dataplane::sai::AclActionType val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_action_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_acl_action_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclBindPointType
convert_sai_acl_bind_point_type_t_to_proto(const sai_int32_t val);
sai_acl_bind_point_type_t convert_sai_acl_bind_point_type_t_to_sai(
    lemming::dataplane::sai::AclBindPointType val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_bind_point_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_acl_bind_point_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclCounterAttr convert_sai_acl_counter_attr_t_to_proto(
    const sai_int32_t val);
sai_acl_counter_attr_t convert_sai_acl_counter_attr_t_to_sai(
    lemming::dataplane::sai::AclCounterAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_counter_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_acl_counter_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclDtelFlowOp convert_sai_acl_dtel_flow_op_t_to_proto(
    const sai_int32_t val);
sai_acl_dtel_flow_op_t convert_sai_acl_dtel_flow_op_t_to_sai(
    lemming::dataplane::sai::AclDtelFlowOp val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_dtel_flow_op_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_acl_dtel_flow_op_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclEntryAttr convert_sai_acl_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_acl_entry_attr_t convert_sai_acl_entry_attr_t_to_sai(
    lemming::dataplane::sai::AclEntryAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_acl_entry_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclIpFrag convert_sai_acl_ip_frag_t_to_proto(
    const sai_int32_t val);
sai_acl_ip_frag_t convert_sai_acl_ip_frag_t_to_sai(
    lemming::dataplane::sai::AclIpFrag val);
google::protobuf::RepeatedField<int> convert_list_sai_acl_ip_frag_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_ip_frag_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclIpType convert_sai_acl_ip_type_t_to_proto(
    const sai_int32_t val);
sai_acl_ip_type_t convert_sai_acl_ip_type_t_to_sai(
    lemming::dataplane::sai::AclIpType val);
google::protobuf::RepeatedField<int> convert_list_sai_acl_ip_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_ip_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclRangeAttr convert_sai_acl_range_attr_t_to_proto(
    const sai_int32_t val);
sai_acl_range_attr_t convert_sai_acl_range_attr_t_to_sai(
    lemming::dataplane::sai::AclRangeAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_acl_range_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_range_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclRangeType convert_sai_acl_range_type_t_to_proto(
    const sai_int32_t val);
sai_acl_range_type_t convert_sai_acl_range_type_t_to_sai(
    lemming::dataplane::sai::AclRangeType val);
google::protobuf::RepeatedField<int> convert_list_sai_acl_range_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_range_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclStage convert_sai_acl_stage_t_to_proto(
    const sai_int32_t val);
sai_acl_stage_t convert_sai_acl_stage_t_to_sai(
    lemming::dataplane::sai::AclStage val);
google::protobuf::RepeatedField<int> convert_list_sai_acl_stage_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_stage_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclTableAttr convert_sai_acl_table_attr_t_to_proto(
    const sai_int32_t val);
sai_acl_table_attr_t convert_sai_acl_table_attr_t_to_sai(
    lemming::dataplane::sai::AclTableAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_acl_table_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_table_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclTableChainGroupAttr
convert_sai_acl_table_chain_group_attr_t_to_proto(const sai_int32_t val);
sai_acl_table_chain_group_attr_t
convert_sai_acl_table_chain_group_attr_t_to_sai(
    lemming::dataplane::sai::AclTableChainGroupAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_chain_group_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_table_chain_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclTableChainGroupStage
convert_sai_acl_table_chain_group_stage_t_to_proto(const sai_int32_t val);
sai_acl_table_chain_group_stage_t
convert_sai_acl_table_chain_group_stage_t_to_sai(
    lemming::dataplane::sai::AclTableChainGroupStage val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_chain_group_stage_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_table_chain_group_stage_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclTableChainGroupType
convert_sai_acl_table_chain_group_type_t_to_proto(const sai_int32_t val);
sai_acl_table_chain_group_type_t
convert_sai_acl_table_chain_group_type_t_to_sai(
    lemming::dataplane::sai::AclTableChainGroupType val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_chain_group_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_table_chain_group_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclTableGroupAttr
convert_sai_acl_table_group_attr_t_to_proto(const sai_int32_t val);
sai_acl_table_group_attr_t convert_sai_acl_table_group_attr_t_to_sai(
    lemming::dataplane::sai::AclTableGroupAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_group_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_acl_table_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclTableGroupMemberAttr
convert_sai_acl_table_group_member_attr_t_to_proto(const sai_int32_t val);
sai_acl_table_group_member_attr_t
convert_sai_acl_table_group_member_attr_t_to_sai(
    lemming::dataplane::sai::AclTableGroupMemberAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_group_member_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_table_group_member_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclTableGroupType
convert_sai_acl_table_group_type_t_to_proto(const sai_int32_t val);
sai_acl_table_group_type_t convert_sai_acl_table_group_type_t_to_sai(
    lemming::dataplane::sai::AclTableGroupType val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_group_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_acl_table_group_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclTableMatchType
convert_sai_acl_table_match_type_t_to_proto(const sai_int32_t val);
sai_acl_table_match_type_t convert_sai_acl_table_match_type_t_to_sai(
    lemming::dataplane::sai::AclTableMatchType val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_match_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_acl_table_match_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::AclTableSupportedMatchType
convert_sai_acl_table_supported_match_type_t_to_proto(const sai_int32_t val);
sai_acl_table_supported_match_type_t
convert_sai_acl_table_supported_match_type_t_to_sai(
    lemming::dataplane::sai::AclTableSupportedMatchType val);
google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_supported_match_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_acl_table_supported_match_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::ApiExtensions convert_sai_api_extensions_t_to_proto(
    const sai_int32_t val);
sai_api_extensions_t convert_sai_api_extensions_t_to_sai(
    lemming::dataplane::sai::ApiExtensions val);
google::protobuf::RepeatedField<int> convert_list_sai_api_extensions_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_api_extensions_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::Api convert_sai_api_t_to_proto(const sai_int32_t val);
sai_api_t convert_sai_api_t_to_sai(lemming::dataplane::sai::Api val);
google::protobuf::RepeatedField<int> convert_list_sai_api_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_api_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BfdEncapsulationType
convert_sai_bfd_encapsulation_type_t_to_proto(const sai_int32_t val);
sai_bfd_encapsulation_type_t convert_sai_bfd_encapsulation_type_t_to_sai(
    lemming::dataplane::sai::BfdEncapsulationType val);
google::protobuf::RepeatedField<int>
convert_list_sai_bfd_encapsulation_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_bfd_encapsulation_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BfdSessionAttr convert_sai_bfd_session_attr_t_to_proto(
    const sai_int32_t val);
sai_bfd_session_attr_t convert_sai_bfd_session_attr_t_to_sai(
    lemming::dataplane::sai::BfdSessionAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_bfd_session_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_bfd_session_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BfdSessionOffloadType
convert_sai_bfd_session_offload_type_t_to_proto(const sai_int32_t val);
sai_bfd_session_offload_type_t convert_sai_bfd_session_offload_type_t_to_sai(
    lemming::dataplane::sai::BfdSessionOffloadType val);
google::protobuf::RepeatedField<int>
convert_list_sai_bfd_session_offload_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_bfd_session_offload_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BfdSessionStat convert_sai_bfd_session_stat_t_to_proto(
    const sai_int32_t val);
sai_bfd_session_stat_t convert_sai_bfd_session_stat_t_to_sai(
    lemming::dataplane::sai::BfdSessionStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_bfd_session_stat_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_bfd_session_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BfdSessionState
convert_sai_bfd_session_state_t_to_proto(const sai_int32_t val);
sai_bfd_session_state_t convert_sai_bfd_session_state_t_to_sai(
    lemming::dataplane::sai::BfdSessionState val);
google::protobuf::RepeatedField<int>
convert_list_sai_bfd_session_state_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_bfd_session_state_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BfdSessionType convert_sai_bfd_session_type_t_to_proto(
    const sai_int32_t val);
sai_bfd_session_type_t convert_sai_bfd_session_type_t_to_sai(
    lemming::dataplane::sai::BfdSessionType val);
google::protobuf::RepeatedField<int>
convert_list_sai_bfd_session_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_bfd_session_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BridgeAttr convert_sai_bridge_attr_t_to_proto(
    const sai_int32_t val);
sai_bridge_attr_t convert_sai_bridge_attr_t_to_sai(
    lemming::dataplane::sai::BridgeAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_bridge_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_bridge_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BridgeFloodControlType
convert_sai_bridge_flood_control_type_t_to_proto(const sai_int32_t val);
sai_bridge_flood_control_type_t convert_sai_bridge_flood_control_type_t_to_sai(
    lemming::dataplane::sai::BridgeFloodControlType val);
google::protobuf::RepeatedField<int>
convert_list_sai_bridge_flood_control_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_bridge_flood_control_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BridgePortAttr convert_sai_bridge_port_attr_t_to_proto(
    const sai_int32_t val);
sai_bridge_port_attr_t convert_sai_bridge_port_attr_t_to_sai(
    lemming::dataplane::sai::BridgePortAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_bridge_port_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_bridge_port_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BridgePortFdbLearningMode
convert_sai_bridge_port_fdb_learning_mode_t_to_proto(const sai_int32_t val);
sai_bridge_port_fdb_learning_mode_t
convert_sai_bridge_port_fdb_learning_mode_t_to_sai(
    lemming::dataplane::sai::BridgePortFdbLearningMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_bridge_port_fdb_learning_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_bridge_port_fdb_learning_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BridgePortStat convert_sai_bridge_port_stat_t_to_proto(
    const sai_int32_t val);
sai_bridge_port_stat_t convert_sai_bridge_port_stat_t_to_sai(
    lemming::dataplane::sai::BridgePortStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_bridge_port_stat_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_bridge_port_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BridgePortTaggingMode
convert_sai_bridge_port_tagging_mode_t_to_proto(const sai_int32_t val);
sai_bridge_port_tagging_mode_t convert_sai_bridge_port_tagging_mode_t_to_sai(
    lemming::dataplane::sai::BridgePortTaggingMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_bridge_port_tagging_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_bridge_port_tagging_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BridgePortType convert_sai_bridge_port_type_t_to_proto(
    const sai_int32_t val);
sai_bridge_port_type_t convert_sai_bridge_port_type_t_to_sai(
    lemming::dataplane::sai::BridgePortType val);
google::protobuf::RepeatedField<int>
convert_list_sai_bridge_port_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_bridge_port_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BridgeStat convert_sai_bridge_stat_t_to_proto(
    const sai_int32_t val);
sai_bridge_stat_t convert_sai_bridge_stat_t_to_sai(
    lemming::dataplane::sai::BridgeStat val);
google::protobuf::RepeatedField<int> convert_list_sai_bridge_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_bridge_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BridgeType convert_sai_bridge_type_t_to_proto(
    const sai_int32_t val);
sai_bridge_type_t convert_sai_bridge_type_t_to_sai(
    lemming::dataplane::sai::BridgeType val);
google::protobuf::RepeatedField<int> convert_list_sai_bridge_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_bridge_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BufferPoolAttr convert_sai_buffer_pool_attr_t_to_proto(
    const sai_int32_t val);
sai_buffer_pool_attr_t convert_sai_buffer_pool_attr_t_to_sai(
    lemming::dataplane::sai::BufferPoolAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_buffer_pool_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_buffer_pool_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BufferPoolStat convert_sai_buffer_pool_stat_t_to_proto(
    const sai_int32_t val);
sai_buffer_pool_stat_t convert_sai_buffer_pool_stat_t_to_sai(
    lemming::dataplane::sai::BufferPoolStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_buffer_pool_stat_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_buffer_pool_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BufferPoolThresholdMode
convert_sai_buffer_pool_threshold_mode_t_to_proto(const sai_int32_t val);
sai_buffer_pool_threshold_mode_t
convert_sai_buffer_pool_threshold_mode_t_to_sai(
    lemming::dataplane::sai::BufferPoolThresholdMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_buffer_pool_threshold_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_buffer_pool_threshold_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BufferPoolType convert_sai_buffer_pool_type_t_to_proto(
    const sai_int32_t val);
sai_buffer_pool_type_t convert_sai_buffer_pool_type_t_to_sai(
    lemming::dataplane::sai::BufferPoolType val);
google::protobuf::RepeatedField<int>
convert_list_sai_buffer_pool_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_buffer_pool_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BufferProfileAttr
convert_sai_buffer_profile_attr_t_to_proto(const sai_int32_t val);
sai_buffer_profile_attr_t convert_sai_buffer_profile_attr_t_to_sai(
    lemming::dataplane::sai::BufferProfileAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_buffer_profile_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_buffer_profile_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BufferProfileThresholdMode
convert_sai_buffer_profile_threshold_mode_t_to_proto(const sai_int32_t val);
sai_buffer_profile_threshold_mode_t
convert_sai_buffer_profile_threshold_mode_t_to_sai(
    lemming::dataplane::sai::BufferProfileThresholdMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_buffer_profile_threshold_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_buffer_profile_threshold_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::BulkOpErrorMode
convert_sai_bulk_op_error_mode_t_to_proto(const sai_int32_t val);
sai_bulk_op_error_mode_t convert_sai_bulk_op_error_mode_t_to_sai(
    lemming::dataplane::sai::BulkOpErrorMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_bulk_op_error_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_bulk_op_error_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::CommonApi convert_sai_common_api_t_to_proto(
    const sai_int32_t val);
sai_common_api_t convert_sai_common_api_t_to_sai(
    lemming::dataplane::sai::CommonApi val);
google::protobuf::RepeatedField<int> convert_list_sai_common_api_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_common_api_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::CounterAttr convert_sai_counter_attr_t_to_proto(
    const sai_int32_t val);
sai_counter_attr_t convert_sai_counter_attr_t_to_sai(
    lemming::dataplane::sai::CounterAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_counter_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_counter_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::CounterStat convert_sai_counter_stat_t_to_proto(
    const sai_int32_t val);
sai_counter_stat_t convert_sai_counter_stat_t_to_sai(
    lemming::dataplane::sai::CounterStat val);
google::protobuf::RepeatedField<int> convert_list_sai_counter_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_counter_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::CounterType convert_sai_counter_type_t_to_proto(
    const sai_int32_t val);
sai_counter_type_t convert_sai_counter_type_t_to_sai(
    lemming::dataplane::sai::CounterType val);
google::protobuf::RepeatedField<int> convert_list_sai_counter_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_counter_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashCapsHaScopeLevel
convert_sai_dash_caps_ha_scope_level_t_to_proto(const sai_int32_t val);
sai_dash_caps_ha_scope_level_t convert_sai_dash_caps_ha_scope_level_t_to_sai(
    lemming::dataplane::sai::DashCapsHaScopeLevel val);
google::protobuf::RepeatedField<int>
convert_list_sai_dash_caps_ha_scope_level_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_dash_caps_ha_scope_level_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashDirection convert_sai_dash_direction_t_to_proto(
    const sai_int32_t val);
sai_dash_direction_t convert_sai_dash_direction_t_to_sai(
    lemming::dataplane::sai::DashDirection val);
google::protobuf::RepeatedField<int> convert_list_sai_dash_direction_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_dash_direction_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashEncapsulation
convert_sai_dash_encapsulation_t_to_proto(const sai_int32_t val);
sai_dash_encapsulation_t convert_sai_dash_encapsulation_t_to_sai(
    lemming::dataplane::sai::DashEncapsulation val);
google::protobuf::RepeatedField<int>
convert_list_sai_dash_encapsulation_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_dash_encapsulation_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashEniMacOverrideType
convert_sai_dash_eni_mac_override_type_t_to_proto(const sai_int32_t val);
sai_dash_eni_mac_override_type_t
convert_sai_dash_eni_mac_override_type_t_to_sai(
    lemming::dataplane::sai::DashEniMacOverrideType val);
google::protobuf::RepeatedField<int>
convert_list_sai_dash_eni_mac_override_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_dash_eni_mac_override_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashFlowAction convert_sai_dash_flow_action_t_to_proto(
    const sai_int32_t val);
sai_dash_flow_action_t convert_sai_dash_flow_action_t_to_sai(
    lemming::dataplane::sai::DashFlowAction val);
google::protobuf::RepeatedField<int>
convert_list_sai_dash_flow_action_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_dash_flow_action_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashFlowEnabledKey
convert_sai_dash_flow_enabled_key_t_to_proto(const sai_int32_t val);
sai_dash_flow_enabled_key_t convert_sai_dash_flow_enabled_key_t_to_sai(
    lemming::dataplane::sai::DashFlowEnabledKey val);
google::protobuf::RepeatedField<int>
convert_list_sai_dash_flow_enabled_key_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_dash_flow_enabled_key_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashFlowEntryBulkGetSessionFilterKey
convert_sai_dash_flow_entry_bulk_get_session_filter_key_t_to_proto(
    const sai_int32_t val);
sai_dash_flow_entry_bulk_get_session_filter_key_t
convert_sai_dash_flow_entry_bulk_get_session_filter_key_t_to_sai(
    lemming::dataplane::sai::DashFlowEntryBulkGetSessionFilterKey val);
google::protobuf::RepeatedField<int>
convert_list_sai_dash_flow_entry_bulk_get_session_filter_key_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_dash_flow_entry_bulk_get_session_filter_key_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashFlowEntryBulkGetSessionMode
convert_sai_dash_flow_entry_bulk_get_session_mode_t_to_proto(
    const sai_int32_t val);
sai_dash_flow_entry_bulk_get_session_mode_t
convert_sai_dash_flow_entry_bulk_get_session_mode_t_to_sai(
    lemming::dataplane::sai::DashFlowEntryBulkGetSessionMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_dash_flow_entry_bulk_get_session_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_dash_flow_entry_bulk_get_session_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashFlowEntryBulkGetSessionOpKey
convert_sai_dash_flow_entry_bulk_get_session_op_key_t_to_proto(
    const sai_int32_t val);
sai_dash_flow_entry_bulk_get_session_op_key_t
convert_sai_dash_flow_entry_bulk_get_session_op_key_t_to_sai(
    lemming::dataplane::sai::DashFlowEntryBulkGetSessionOpKey val);
google::protobuf::RepeatedField<int>
convert_list_sai_dash_flow_entry_bulk_get_session_op_key_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_dash_flow_entry_bulk_get_session_op_key_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashHaRole convert_sai_dash_ha_role_t_to_proto(
    const sai_int32_t val);
sai_dash_ha_role_t convert_sai_dash_ha_role_t_to_sai(
    lemming::dataplane::sai::DashHaRole val);
google::protobuf::RepeatedField<int> convert_list_sai_dash_ha_role_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_dash_ha_role_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashHaState convert_sai_dash_ha_state_t_to_proto(
    const sai_int32_t val);
sai_dash_ha_state_t convert_sai_dash_ha_state_t_to_sai(
    lemming::dataplane::sai::DashHaState val);
google::protobuf::RepeatedField<int> convert_list_sai_dash_ha_state_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_dash_ha_state_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashRoutingActions
convert_sai_dash_routing_actions_t_to_proto(const sai_int32_t val);
sai_dash_routing_actions_t convert_sai_dash_routing_actions_t_to_sai(
    lemming::dataplane::sai::DashRoutingActions val);
google::protobuf::RepeatedField<int>
convert_list_sai_dash_routing_actions_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_dash_routing_actions_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DashTunnelDscpMode
convert_sai_dash_tunnel_dscp_mode_t_to_proto(const sai_int32_t val);
sai_dash_tunnel_dscp_mode_t convert_sai_dash_tunnel_dscp_mode_t_to_sai(
    lemming::dataplane::sai::DashTunnelDscpMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_dash_tunnel_dscp_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_dash_tunnel_dscp_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DebugCounterAttr
convert_sai_debug_counter_attr_t_to_proto(const sai_int32_t val);
sai_debug_counter_attr_t convert_sai_debug_counter_attr_t_to_sai(
    lemming::dataplane::sai::DebugCounterAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_debug_counter_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_debug_counter_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DebugCounterBindMethod
convert_sai_debug_counter_bind_method_t_to_proto(const sai_int32_t val);
sai_debug_counter_bind_method_t convert_sai_debug_counter_bind_method_t_to_sai(
    lemming::dataplane::sai::DebugCounterBindMethod val);
google::protobuf::RepeatedField<int>
convert_list_sai_debug_counter_bind_method_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_debug_counter_bind_method_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DebugCounterType
convert_sai_debug_counter_type_t_to_proto(const sai_int32_t val);
sai_debug_counter_type_t convert_sai_debug_counter_type_t_to_sai(
    lemming::dataplane::sai::DebugCounterType val);
google::protobuf::RepeatedField<int>
convert_list_sai_debug_counter_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_debug_counter_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DtelAttr convert_sai_dtel_attr_t_to_proto(
    const sai_int32_t val);
sai_dtel_attr_t convert_sai_dtel_attr_t_to_sai(
    lemming::dataplane::sai::DtelAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_dtel_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_dtel_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DtelEventAttr convert_sai_dtel_event_attr_t_to_proto(
    const sai_int32_t val);
sai_dtel_event_attr_t convert_sai_dtel_event_attr_t_to_sai(
    lemming::dataplane::sai::DtelEventAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_dtel_event_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_dtel_event_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DtelEventType convert_sai_dtel_event_type_t_to_proto(
    const sai_int32_t val);
sai_dtel_event_type_t convert_sai_dtel_event_type_t_to_sai(
    lemming::dataplane::sai::DtelEventType val);
google::protobuf::RepeatedField<int>
convert_list_sai_dtel_event_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_dtel_event_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DtelIntSessionAttr
convert_sai_dtel_int_session_attr_t_to_proto(const sai_int32_t val);
sai_dtel_int_session_attr_t convert_sai_dtel_int_session_attr_t_to_sai(
    lemming::dataplane::sai::DtelIntSessionAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_dtel_int_session_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_dtel_int_session_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DtelQueueReportAttr
convert_sai_dtel_queue_report_attr_t_to_proto(const sai_int32_t val);
sai_dtel_queue_report_attr_t convert_sai_dtel_queue_report_attr_t_to_sai(
    lemming::dataplane::sai::DtelQueueReportAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_dtel_queue_report_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_dtel_queue_report_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::DtelReportSessionAttr
convert_sai_dtel_report_session_attr_t_to_proto(const sai_int32_t val);
sai_dtel_report_session_attr_t convert_sai_dtel_report_session_attr_t_to_sai(
    lemming::dataplane::sai::DtelReportSessionAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_dtel_report_session_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_dtel_report_session_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::EcnMarkMode convert_sai_ecn_mark_mode_t_to_proto(
    const sai_int32_t val);
sai_ecn_mark_mode_t convert_sai_ecn_mark_mode_t_to_sai(
    lemming::dataplane::sai::EcnMarkMode val);
google::protobuf::RepeatedField<int> convert_list_sai_ecn_mark_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ecn_mark_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::ErspanEncapsulationType
convert_sai_erspan_encapsulation_type_t_to_proto(const sai_int32_t val);
sai_erspan_encapsulation_type_t convert_sai_erspan_encapsulation_type_t_to_sai(
    lemming::dataplane::sai::ErspanEncapsulationType val);
google::protobuf::RepeatedField<int>
convert_list_sai_erspan_encapsulation_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_erspan_encapsulation_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::FdbEntryAttr convert_sai_fdb_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_fdb_entry_attr_t convert_sai_fdb_entry_attr_t_to_sai(
    lemming::dataplane::sai::FdbEntryAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_fdb_entry_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_fdb_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::FdbEntryType convert_sai_fdb_entry_type_t_to_proto(
    const sai_int32_t val);
sai_fdb_entry_type_t convert_sai_fdb_entry_type_t_to_sai(
    lemming::dataplane::sai::FdbEntryType val);
google::protobuf::RepeatedField<int> convert_list_sai_fdb_entry_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_fdb_entry_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::FdbEvent convert_sai_fdb_event_t_to_proto(
    const sai_int32_t val);
sai_fdb_event_t convert_sai_fdb_event_t_to_sai(
    lemming::dataplane::sai::FdbEvent val);
google::protobuf::RepeatedField<int> convert_list_sai_fdb_event_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_fdb_event_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::FdbFlushAttr convert_sai_fdb_flush_attr_t_to_proto(
    const sai_int32_t val);
sai_fdb_flush_attr_t convert_sai_fdb_flush_attr_t_to_sai(
    lemming::dataplane::sai::FdbFlushAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_fdb_flush_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_fdb_flush_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::FdbFlushEntryType
convert_sai_fdb_flush_entry_type_t_to_proto(const sai_int32_t val);
sai_fdb_flush_entry_type_t convert_sai_fdb_flush_entry_type_t_to_sai(
    lemming::dataplane::sai::FdbFlushEntryType val);
google::protobuf::RepeatedField<int>
convert_list_sai_fdb_flush_entry_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_fdb_flush_entry_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::FineGrainedHashFieldAttr
convert_sai_fine_grained_hash_field_attr_t_to_proto(const sai_int32_t val);
sai_fine_grained_hash_field_attr_t
convert_sai_fine_grained_hash_field_attr_t_to_sai(
    lemming::dataplane::sai::FineGrainedHashFieldAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_fine_grained_hash_field_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_fine_grained_hash_field_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::GenericProgrammableAttr
convert_sai_generic_programmable_attr_t_to_proto(const sai_int32_t val);
sai_generic_programmable_attr_t convert_sai_generic_programmable_attr_t_to_sai(
    lemming::dataplane::sai::GenericProgrammableAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_generic_programmable_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_generic_programmable_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HaScopeEvent convert_sai_ha_scope_event_t_to_proto(
    const sai_int32_t val);
sai_ha_scope_event_t convert_sai_ha_scope_event_t_to_sai(
    lemming::dataplane::sai::HaScopeEvent val);
google::protobuf::RepeatedField<int> convert_list_sai_ha_scope_event_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ha_scope_event_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HaSetEvent convert_sai_ha_set_event_t_to_proto(
    const sai_int32_t val);
sai_ha_set_event_t convert_sai_ha_set_event_t_to_sai(
    lemming::dataplane::sai::HaSetEvent val);
google::protobuf::RepeatedField<int> convert_list_sai_ha_set_event_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ha_set_event_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HashAlgorithm convert_sai_hash_algorithm_t_to_proto(
    const sai_int32_t val);
sai_hash_algorithm_t convert_sai_hash_algorithm_t_to_sai(
    lemming::dataplane::sai::HashAlgorithm val);
google::protobuf::RepeatedField<int> convert_list_sai_hash_algorithm_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_hash_algorithm_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HashAttr convert_sai_hash_attr_t_to_proto(
    const sai_int32_t val);
sai_hash_attr_t convert_sai_hash_attr_t_to_sai(
    lemming::dataplane::sai::HashAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_hash_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_hash_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HealthDataType convert_sai_health_data_type_t_to_proto(
    const sai_int32_t val);
sai_health_data_type_t convert_sai_health_data_type_t_to_sai(
    lemming::dataplane::sai::HealthDataType val);
google::protobuf::RepeatedField<int>
convert_list_sai_health_data_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_health_data_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifAttr convert_sai_hostif_attr_t_to_proto(
    const sai_int32_t val);
sai_hostif_attr_t convert_sai_hostif_attr_t_to_sai(
    lemming::dataplane::sai::HostifAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_hostif_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_hostif_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifTableEntryAttr
convert_sai_hostif_table_entry_attr_t_to_proto(const sai_int32_t val);
sai_hostif_table_entry_attr_t convert_sai_hostif_table_entry_attr_t_to_sai(
    lemming::dataplane::sai::HostifTableEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_hostif_table_entry_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_hostif_table_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifTableEntryChannelType
convert_sai_hostif_table_entry_channel_type_t_to_proto(const sai_int32_t val);
sai_hostif_table_entry_channel_type_t
convert_sai_hostif_table_entry_channel_type_t_to_sai(
    lemming::dataplane::sai::HostifTableEntryChannelType val);
google::protobuf::RepeatedField<int>
convert_list_sai_hostif_table_entry_channel_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_hostif_table_entry_channel_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifTableEntryType
convert_sai_hostif_table_entry_type_t_to_proto(const sai_int32_t val);
sai_hostif_table_entry_type_t convert_sai_hostif_table_entry_type_t_to_sai(
    lemming::dataplane::sai::HostifTableEntryType val);
google::protobuf::RepeatedField<int>
convert_list_sai_hostif_table_entry_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_hostif_table_entry_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifTrapAttr convert_sai_hostif_trap_attr_t_to_proto(
    const sai_int32_t val);
sai_hostif_trap_attr_t convert_sai_hostif_trap_attr_t_to_sai(
    lemming::dataplane::sai::HostifTrapAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_hostif_trap_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_hostif_trap_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifTrapGroupAttr
convert_sai_hostif_trap_group_attr_t_to_proto(const sai_int32_t val);
sai_hostif_trap_group_attr_t convert_sai_hostif_trap_group_attr_t_to_sai(
    lemming::dataplane::sai::HostifTrapGroupAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_hostif_trap_group_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_hostif_trap_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifTrapType convert_sai_hostif_trap_type_t_to_proto(
    const sai_int32_t val);
sai_hostif_trap_type_t convert_sai_hostif_trap_type_t_to_sai(
    lemming::dataplane::sai::HostifTrapType val);
google::protobuf::RepeatedField<int>
convert_list_sai_hostif_trap_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_hostif_trap_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifTxType convert_sai_hostif_tx_type_t_to_proto(
    const sai_int32_t val);
sai_hostif_tx_type_t convert_sai_hostif_tx_type_t_to_sai(
    lemming::dataplane::sai::HostifTxType val);
google::protobuf::RepeatedField<int> convert_list_sai_hostif_tx_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_hostif_tx_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifType convert_sai_hostif_type_t_to_proto(
    const sai_int32_t val);
sai_hostif_type_t convert_sai_hostif_type_t_to_sai(
    lemming::dataplane::sai::HostifType val);
google::protobuf::RepeatedField<int> convert_list_sai_hostif_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_hostif_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifUserDefinedTrapAttr
convert_sai_hostif_user_defined_trap_attr_t_to_proto(const sai_int32_t val);
sai_hostif_user_defined_trap_attr_t
convert_sai_hostif_user_defined_trap_attr_t_to_sai(
    lemming::dataplane::sai::HostifUserDefinedTrapAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_hostif_user_defined_trap_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_hostif_user_defined_trap_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifUserDefinedTrapType
convert_sai_hostif_user_defined_trap_type_t_to_proto(const sai_int32_t val);
sai_hostif_user_defined_trap_type_t
convert_sai_hostif_user_defined_trap_type_t_to_sai(
    lemming::dataplane::sai::HostifUserDefinedTrapType val);
google::protobuf::RepeatedField<int>
convert_list_sai_hostif_user_defined_trap_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_hostif_user_defined_trap_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::HostifVlanTag convert_sai_hostif_vlan_tag_t_to_proto(
    const sai_int32_t val);
sai_hostif_vlan_tag_t convert_sai_hostif_vlan_tag_t_to_sai(
    lemming::dataplane::sai::HostifVlanTag val);
google::protobuf::RepeatedField<int>
convert_list_sai_hostif_vlan_tag_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_hostif_vlan_tag_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IcmpEchoSessionAttr
convert_sai_icmp_echo_session_attr_t_to_proto(const sai_int32_t val);
sai_icmp_echo_session_attr_t convert_sai_icmp_echo_session_attr_t_to_sai(
    lemming::dataplane::sai::IcmpEchoSessionAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_icmp_echo_session_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_icmp_echo_session_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IcmpEchoSessionStat
convert_sai_icmp_echo_session_stat_t_to_proto(const sai_int32_t val);
sai_icmp_echo_session_stat_t convert_sai_icmp_echo_session_stat_t_to_sai(
    lemming::dataplane::sai::IcmpEchoSessionStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_icmp_echo_session_stat_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_icmp_echo_session_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IcmpEchoSessionState
convert_sai_icmp_echo_session_state_t_to_proto(const sai_int32_t val);
sai_icmp_echo_session_state_t convert_sai_icmp_echo_session_state_t_to_sai(
    lemming::dataplane::sai::IcmpEchoSessionState val);
google::protobuf::RepeatedField<int>
convert_list_sai_icmp_echo_session_state_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_icmp_echo_session_state_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::InDropReason convert_sai_in_drop_reason_t_to_proto(
    const sai_int32_t val);
sai_in_drop_reason_t convert_sai_in_drop_reason_t_to_sai(
    lemming::dataplane::sai::InDropReason val);
google::protobuf::RepeatedField<int> convert_list_sai_in_drop_reason_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_in_drop_reason_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IngressPriorityGroupAttr
convert_sai_ingress_priority_group_attr_t_to_proto(const sai_int32_t val);
sai_ingress_priority_group_attr_t
convert_sai_ingress_priority_group_attr_t_to_sai(
    lemming::dataplane::sai::IngressPriorityGroupAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_ingress_priority_group_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ingress_priority_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IngressPriorityGroupStat
convert_sai_ingress_priority_group_stat_t_to_proto(const sai_int32_t val);
sai_ingress_priority_group_stat_t
convert_sai_ingress_priority_group_stat_t_to_sai(
    lemming::dataplane::sai::IngressPriorityGroupStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_ingress_priority_group_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ingress_priority_group_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::InsegEntryAttr convert_sai_inseg_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_inseg_entry_attr_t convert_sai_inseg_entry_attr_t_to_sai(
    lemming::dataplane::sai::InsegEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_inseg_entry_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_inseg_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::InsegEntryPopQosMode
convert_sai_inseg_entry_pop_qos_mode_t_to_proto(const sai_int32_t val);
sai_inseg_entry_pop_qos_mode_t convert_sai_inseg_entry_pop_qos_mode_t_to_sai(
    lemming::dataplane::sai::InsegEntryPopQosMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_inseg_entry_pop_qos_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_inseg_entry_pop_qos_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::InsegEntryPopTtlMode
convert_sai_inseg_entry_pop_ttl_mode_t_to_proto(const sai_int32_t val);
sai_inseg_entry_pop_ttl_mode_t convert_sai_inseg_entry_pop_ttl_mode_t_to_sai(
    lemming::dataplane::sai::InsegEntryPopTtlMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_inseg_entry_pop_ttl_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_inseg_entry_pop_ttl_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::InsegEntryPscType
convert_sai_inseg_entry_psc_type_t_to_proto(const sai_int32_t val);
sai_inseg_entry_psc_type_t convert_sai_inseg_entry_psc_type_t_to_sai(
    lemming::dataplane::sai::InsegEntryPscType val);
google::protobuf::RepeatedField<int>
convert_list_sai_inseg_entry_psc_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_inseg_entry_psc_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpAddrFamily convert_sai_ip_addr_family_t_to_proto(
    const sai_int32_t val);
sai_ip_addr_family_t convert_sai_ip_addr_family_t_to_sai(
    lemming::dataplane::sai::IpAddrFamily val);
google::protobuf::RepeatedField<int> convert_list_sai_ip_addr_family_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ip_addr_family_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpmcEntryAttr convert_sai_ipmc_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_ipmc_entry_attr_t convert_sai_ipmc_entry_attr_t_to_sai(
    lemming::dataplane::sai::IpmcEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_ipmc_entry_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_ipmc_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpmcEntryType convert_sai_ipmc_entry_type_t_to_proto(
    const sai_int32_t val);
sai_ipmc_entry_type_t convert_sai_ipmc_entry_type_t_to_sai(
    lemming::dataplane::sai::IpmcEntryType val);
google::protobuf::RepeatedField<int>
convert_list_sai_ipmc_entry_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_ipmc_entry_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpmcGroupAttr convert_sai_ipmc_group_attr_t_to_proto(
    const sai_int32_t val);
sai_ipmc_group_attr_t convert_sai_ipmc_group_attr_t_to_sai(
    lemming::dataplane::sai::IpmcGroupAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_ipmc_group_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_ipmc_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpmcGroupMemberAttr
convert_sai_ipmc_group_member_attr_t_to_proto(const sai_int32_t val);
sai_ipmc_group_member_attr_t convert_sai_ipmc_group_member_attr_t_to_sai(
    lemming::dataplane::sai::IpmcGroupMemberAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_ipmc_group_member_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_ipmc_group_member_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpsecAttr convert_sai_ipsec_attr_t_to_proto(
    const sai_int32_t val);
sai_ipsec_attr_t convert_sai_ipsec_attr_t_to_sai(
    lemming::dataplane::sai::IpsecAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_ipsec_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ipsec_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpsecCipher convert_sai_ipsec_cipher_t_to_proto(
    const sai_int32_t val);
sai_ipsec_cipher_t convert_sai_ipsec_cipher_t_to_sai(
    lemming::dataplane::sai::IpsecCipher val);
google::protobuf::RepeatedField<int> convert_list_sai_ipsec_cipher_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ipsec_cipher_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpsecDirection convert_sai_ipsec_direction_t_to_proto(
    const sai_int32_t val);
sai_ipsec_direction_t convert_sai_ipsec_direction_t_to_sai(
    lemming::dataplane::sai::IpsecDirection val);
google::protobuf::RepeatedField<int>
convert_list_sai_ipsec_direction_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_ipsec_direction_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpsecPortAttr convert_sai_ipsec_port_attr_t_to_proto(
    const sai_int32_t val);
sai_ipsec_port_attr_t convert_sai_ipsec_port_attr_t_to_sai(
    lemming::dataplane::sai::IpsecPortAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_ipsec_port_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_ipsec_port_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpsecPortStat convert_sai_ipsec_port_stat_t_to_proto(
    const sai_int32_t val);
sai_ipsec_port_stat_t convert_sai_ipsec_port_stat_t_to_sai(
    lemming::dataplane::sai::IpsecPortStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_ipsec_port_stat_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_ipsec_port_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpsecSaAttr convert_sai_ipsec_sa_attr_t_to_proto(
    const sai_int32_t val);
sai_ipsec_sa_attr_t convert_sai_ipsec_sa_attr_t_to_sai(
    lemming::dataplane::sai::IpsecSaAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_ipsec_sa_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ipsec_sa_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpsecSaOctetCountStatus
convert_sai_ipsec_sa_octet_count_status_t_to_proto(const sai_int32_t val);
sai_ipsec_sa_octet_count_status_t
convert_sai_ipsec_sa_octet_count_status_t_to_sai(
    lemming::dataplane::sai::IpsecSaOctetCountStatus val);
google::protobuf::RepeatedField<int>
convert_list_sai_ipsec_sa_octet_count_status_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ipsec_sa_octet_count_status_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IpsecSaStat convert_sai_ipsec_sa_stat_t_to_proto(
    const sai_int32_t val);
sai_ipsec_sa_stat_t convert_sai_ipsec_sa_stat_t_to_sai(
    lemming::dataplane::sai::IpsecSaStat val);
google::protobuf::RepeatedField<int> convert_list_sai_ipsec_sa_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ipsec_sa_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IsolationGroupAttr
convert_sai_isolation_group_attr_t_to_proto(const sai_int32_t val);
sai_isolation_group_attr_t convert_sai_isolation_group_attr_t_to_sai(
    lemming::dataplane::sai::IsolationGroupAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_isolation_group_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_isolation_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IsolationGroupMemberAttr
convert_sai_isolation_group_member_attr_t_to_proto(const sai_int32_t val);
sai_isolation_group_member_attr_t
convert_sai_isolation_group_member_attr_t_to_sai(
    lemming::dataplane::sai::IsolationGroupMemberAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_isolation_group_member_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_isolation_group_member_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::IsolationGroupType
convert_sai_isolation_group_type_t_to_proto(const sai_int32_t val);
sai_isolation_group_type_t convert_sai_isolation_group_type_t_to_sai(
    lemming::dataplane::sai::IsolationGroupType val);
google::protobuf::RepeatedField<int>
convert_list_sai_isolation_group_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_isolation_group_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::L2mcEntryAttr convert_sai_l2mc_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_l2mc_entry_attr_t convert_sai_l2mc_entry_attr_t_to_sai(
    lemming::dataplane::sai::L2mcEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_l2mc_entry_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_l2mc_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::L2mcEntryType convert_sai_l2mc_entry_type_t_to_proto(
    const sai_int32_t val);
sai_l2mc_entry_type_t convert_sai_l2mc_entry_type_t_to_sai(
    lemming::dataplane::sai::L2mcEntryType val);
google::protobuf::RepeatedField<int>
convert_list_sai_l2mc_entry_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_l2mc_entry_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::L2mcGroupAttr convert_sai_l2mc_group_attr_t_to_proto(
    const sai_int32_t val);
sai_l2mc_group_attr_t convert_sai_l2mc_group_attr_t_to_sai(
    lemming::dataplane::sai::L2mcGroupAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_l2mc_group_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_l2mc_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::L2mcGroupMemberAttr
convert_sai_l2mc_group_member_attr_t_to_proto(const sai_int32_t val);
sai_l2mc_group_member_attr_t convert_sai_l2mc_group_member_attr_t_to_sai(
    lemming::dataplane::sai::L2mcGroupMemberAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_l2mc_group_member_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_l2mc_group_member_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::LagAttr convert_sai_lag_attr_t_to_proto(
    const sai_int32_t val);
sai_lag_attr_t convert_sai_lag_attr_t_to_sai(
    lemming::dataplane::sai::LagAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_lag_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_lag_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::LagMemberAttr convert_sai_lag_member_attr_t_to_proto(
    const sai_int32_t val);
sai_lag_member_attr_t convert_sai_lag_member_attr_t_to_sai(
    lemming::dataplane::sai::LagMemberAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_lag_member_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_lag_member_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::LogLevel convert_sai_log_level_t_to_proto(
    const sai_int32_t val);
sai_log_level_t convert_sai_log_level_t_to_sai(
    lemming::dataplane::sai::LogLevel val);
google::protobuf::RepeatedField<int> convert_list_sai_log_level_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_log_level_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecAttr convert_sai_macsec_attr_t_to_proto(
    const sai_int32_t val);
sai_macsec_attr_t convert_sai_macsec_attr_t_to_sai(
    lemming::dataplane::sai::MacsecAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_macsec_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_macsec_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecCipherSuite
convert_sai_macsec_cipher_suite_t_to_proto(const sai_int32_t val);
sai_macsec_cipher_suite_t convert_sai_macsec_cipher_suite_t_to_sai(
    lemming::dataplane::sai::MacsecCipherSuite val);
google::protobuf::RepeatedField<int>
convert_list_sai_macsec_cipher_suite_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_macsec_cipher_suite_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecDirection
convert_sai_macsec_direction_t_to_proto(const sai_int32_t val);
sai_macsec_direction_t convert_sai_macsec_direction_t_to_sai(
    lemming::dataplane::sai::MacsecDirection val);
google::protobuf::RepeatedField<int>
convert_list_sai_macsec_direction_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_macsec_direction_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecFlowAttr convert_sai_macsec_flow_attr_t_to_proto(
    const sai_int32_t val);
sai_macsec_flow_attr_t convert_sai_macsec_flow_attr_t_to_sai(
    lemming::dataplane::sai::MacsecFlowAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_macsec_flow_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_macsec_flow_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecFlowStat convert_sai_macsec_flow_stat_t_to_proto(
    const sai_int32_t val);
sai_macsec_flow_stat_t convert_sai_macsec_flow_stat_t_to_sai(
    lemming::dataplane::sai::MacsecFlowStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_macsec_flow_stat_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_macsec_flow_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecMaxSecureAssociationsPerSc
convert_sai_macsec_max_secure_associations_per_sc_t_to_proto(
    const sai_int32_t val);
sai_macsec_max_secure_associations_per_sc_t
convert_sai_macsec_max_secure_associations_per_sc_t_to_sai(
    lemming::dataplane::sai::MacsecMaxSecureAssociationsPerSc val);
google::protobuf::RepeatedField<int>
convert_list_sai_macsec_max_secure_associations_per_sc_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_macsec_max_secure_associations_per_sc_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecPortAttr convert_sai_macsec_port_attr_t_to_proto(
    const sai_int32_t val);
sai_macsec_port_attr_t convert_sai_macsec_port_attr_t_to_sai(
    lemming::dataplane::sai::MacsecPortAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_macsec_port_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_macsec_port_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecPortStat convert_sai_macsec_port_stat_t_to_proto(
    const sai_int32_t val);
sai_macsec_port_stat_t convert_sai_macsec_port_stat_t_to_sai(
    lemming::dataplane::sai::MacsecPortStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_macsec_port_stat_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_macsec_port_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecSaAttr convert_sai_macsec_sa_attr_t_to_proto(
    const sai_int32_t val);
sai_macsec_sa_attr_t convert_sai_macsec_sa_attr_t_to_sai(
    lemming::dataplane::sai::MacsecSaAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_macsec_sa_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_macsec_sa_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecSaStat convert_sai_macsec_sa_stat_t_to_proto(
    const sai_int32_t val);
sai_macsec_sa_stat_t convert_sai_macsec_sa_stat_t_to_sai(
    lemming::dataplane::sai::MacsecSaStat val);
google::protobuf::RepeatedField<int> convert_list_sai_macsec_sa_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_macsec_sa_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecScAttr convert_sai_macsec_sc_attr_t_to_proto(
    const sai_int32_t val);
sai_macsec_sc_attr_t convert_sai_macsec_sc_attr_t_to_sai(
    lemming::dataplane::sai::MacsecScAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_macsec_sc_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_macsec_sc_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MacsecScStat convert_sai_macsec_sc_stat_t_to_proto(
    const sai_int32_t val);
sai_macsec_sc_stat_t convert_sai_macsec_sc_stat_t_to_sai(
    lemming::dataplane::sai::MacsecScStat val);
google::protobuf::RepeatedField<int> convert_list_sai_macsec_sc_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_macsec_sc_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::McastFdbEntryAttr
convert_sai_mcast_fdb_entry_attr_t_to_proto(const sai_int32_t val);
sai_mcast_fdb_entry_attr_t convert_sai_mcast_fdb_entry_attr_t_to_sai(
    lemming::dataplane::sai::McastFdbEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_mcast_fdb_entry_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_mcast_fdb_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MeterType convert_sai_meter_type_t_to_proto(
    const sai_int32_t val);
sai_meter_type_t convert_sai_meter_type_t_to_sai(
    lemming::dataplane::sai::MeterType val);
google::protobuf::RepeatedField<int> convert_list_sai_meter_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_meter_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MirrorSessionAttr
convert_sai_mirror_session_attr_t_to_proto(const sai_int32_t val);
sai_mirror_session_attr_t convert_sai_mirror_session_attr_t_to_sai(
    lemming::dataplane::sai::MirrorSessionAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_mirror_session_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_mirror_session_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MirrorSessionCongestionMode
convert_sai_mirror_session_congestion_mode_t_to_proto(const sai_int32_t val);
sai_mirror_session_congestion_mode_t
convert_sai_mirror_session_congestion_mode_t_to_sai(
    lemming::dataplane::sai::MirrorSessionCongestionMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_mirror_session_congestion_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_mirror_session_congestion_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MirrorSessionType
convert_sai_mirror_session_type_t_to_proto(const sai_int32_t val);
sai_mirror_session_type_t convert_sai_mirror_session_type_t_to_sai(
    lemming::dataplane::sai::MirrorSessionType val);
google::protobuf::RepeatedField<int>
convert_list_sai_mirror_session_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_mirror_session_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MyMacAttr convert_sai_my_mac_attr_t_to_proto(
    const sai_int32_t val);
sai_my_mac_attr_t convert_sai_my_mac_attr_t_to_sai(
    lemming::dataplane::sai::MyMacAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_my_mac_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_my_mac_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MySidEntryAttr
convert_sai_my_sid_entry_attr_t_to_proto(const sai_int32_t val);
sai_my_sid_entry_attr_t convert_sai_my_sid_entry_attr_t_to_sai(
    lemming::dataplane::sai::MySidEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_my_sid_entry_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_my_sid_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MySidEntryEndpointBehaviorFlavor
convert_sai_my_sid_entry_endpoint_behavior_flavor_t_to_proto(
    const sai_int32_t val);
sai_my_sid_entry_endpoint_behavior_flavor_t
convert_sai_my_sid_entry_endpoint_behavior_flavor_t_to_sai(
    lemming::dataplane::sai::MySidEntryEndpointBehaviorFlavor val);
google::protobuf::RepeatedField<int>
convert_list_sai_my_sid_entry_endpoint_behavior_flavor_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_my_sid_entry_endpoint_behavior_flavor_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::MySidEntryEndpointBehavior
convert_sai_my_sid_entry_endpoint_behavior_t_to_proto(const sai_int32_t val);
sai_my_sid_entry_endpoint_behavior_t
convert_sai_my_sid_entry_endpoint_behavior_t_to_sai(
    lemming::dataplane::sai::MySidEntryEndpointBehavior val);
google::protobuf::RepeatedField<int>
convert_list_sai_my_sid_entry_endpoint_behavior_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_my_sid_entry_endpoint_behavior_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NatEntryAttr convert_sai_nat_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_nat_entry_attr_t convert_sai_nat_entry_attr_t_to_sai(
    lemming::dataplane::sai::NatEntryAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_nat_entry_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_nat_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NatEvent convert_sai_nat_event_t_to_proto(
    const sai_int32_t val);
sai_nat_event_t convert_sai_nat_event_t_to_sai(
    lemming::dataplane::sai::NatEvent val);
google::protobuf::RepeatedField<int> convert_list_sai_nat_event_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_nat_event_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NatType convert_sai_nat_type_t_to_proto(
    const sai_int32_t val);
sai_nat_type_t convert_sai_nat_type_t_to_sai(
    lemming::dataplane::sai::NatType val);
google::protobuf::RepeatedField<int> convert_list_sai_nat_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_nat_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NatZoneCounterAttr
convert_sai_nat_zone_counter_attr_t_to_proto(const sai_int32_t val);
sai_nat_zone_counter_attr_t convert_sai_nat_zone_counter_attr_t_to_sai(
    lemming::dataplane::sai::NatZoneCounterAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_nat_zone_counter_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_nat_zone_counter_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NativeHashField
convert_sai_native_hash_field_t_to_proto(const sai_int32_t val);
sai_native_hash_field_t convert_sai_native_hash_field_t_to_sai(
    lemming::dataplane::sai::NativeHashField val);
google::protobuf::RepeatedField<int>
convert_list_sai_native_hash_field_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_native_hash_field_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NeighborEntryAttr
convert_sai_neighbor_entry_attr_t_to_proto(const sai_int32_t val);
sai_neighbor_entry_attr_t convert_sai_neighbor_entry_attr_t_to_sai(
    lemming::dataplane::sai::NeighborEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_neighbor_entry_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_neighbor_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NextHopAttr convert_sai_next_hop_attr_t_to_proto(
    const sai_int32_t val);
sai_next_hop_attr_t convert_sai_next_hop_attr_t_to_sai(
    lemming::dataplane::sai::NextHopAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_next_hop_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_next_hop_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NextHopGroupAttr
convert_sai_next_hop_group_attr_t_to_proto(const sai_int32_t val);
sai_next_hop_group_attr_t convert_sai_next_hop_group_attr_t_to_sai(
    lemming::dataplane::sai::NextHopGroupAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_next_hop_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NextHopGroupMapAttr
convert_sai_next_hop_group_map_attr_t_to_proto(const sai_int32_t val);
sai_next_hop_group_map_attr_t convert_sai_next_hop_group_map_attr_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMapAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_map_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_next_hop_group_map_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NextHopGroupMapType
convert_sai_next_hop_group_map_type_t_to_proto(const sai_int32_t val);
sai_next_hop_group_map_type_t convert_sai_next_hop_group_map_type_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMapType val);
google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_map_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_next_hop_group_map_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NextHopGroupMemberAttr
convert_sai_next_hop_group_member_attr_t_to_proto(const sai_int32_t val);
sai_next_hop_group_member_attr_t
convert_sai_next_hop_group_member_attr_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMemberAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_member_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_next_hop_group_member_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NextHopGroupMemberConfiguredRole
convert_sai_next_hop_group_member_configured_role_t_to_proto(
    const sai_int32_t val);
sai_next_hop_group_member_configured_role_t
convert_sai_next_hop_group_member_configured_role_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMemberConfiguredRole val);
google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_member_configured_role_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_next_hop_group_member_configured_role_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NextHopGroupMemberObservedRole
convert_sai_next_hop_group_member_observed_role_t_to_proto(
    const sai_int32_t val);
sai_next_hop_group_member_observed_role_t
convert_sai_next_hop_group_member_observed_role_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMemberObservedRole val);
google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_member_observed_role_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_next_hop_group_member_observed_role_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NextHopGroupType
convert_sai_next_hop_group_type_t_to_proto(const sai_int32_t val);
sai_next_hop_group_type_t convert_sai_next_hop_group_type_t_to_sai(
    lemming::dataplane::sai::NextHopGroupType val);
google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_next_hop_group_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::NextHopType convert_sai_next_hop_type_t_to_proto(
    const sai_int32_t val);
sai_next_hop_type_t convert_sai_next_hop_type_t_to_sai(
    lemming::dataplane::sai::NextHopType val);
google::protobuf::RepeatedField<int> convert_list_sai_next_hop_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_next_hop_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::ObjectStage convert_sai_object_stage_t_to_proto(
    const sai_int32_t val);
sai_object_stage_t convert_sai_object_stage_t_to_sai(
    lemming::dataplane::sai::ObjectStage val);
google::protobuf::RepeatedField<int> convert_list_sai_object_stage_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_object_stage_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::ObjectTypeExtensions
convert_sai_object_type_extensions_t_to_proto(const sai_int32_t val);
sai_object_type_extensions_t convert_sai_object_type_extensions_t_to_sai(
    lemming::dataplane::sai::ObjectTypeExtensions val);
google::protobuf::RepeatedField<int>
convert_list_sai_object_type_extensions_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_object_type_extensions_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::ObjectType convert_sai_object_type_t_to_proto(
    const sai_int32_t val);
sai_object_type_t convert_sai_object_type_t_to_sai(
    lemming::dataplane::sai::ObjectType val);
google::protobuf::RepeatedField<int> convert_list_sai_object_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_object_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::OutDropReason convert_sai_out_drop_reason_t_to_proto(
    const sai_int32_t val);
sai_out_drop_reason_t convert_sai_out_drop_reason_t_to_sai(
    lemming::dataplane::sai::OutDropReason val);
google::protobuf::RepeatedField<int>
convert_list_sai_out_drop_reason_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_out_drop_reason_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::OutsegExpMode convert_sai_outseg_exp_mode_t_to_proto(
    const sai_int32_t val);
sai_outseg_exp_mode_t convert_sai_outseg_exp_mode_t_to_sai(
    lemming::dataplane::sai::OutsegExpMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_outseg_exp_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_outseg_exp_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::OutsegTtlMode convert_sai_outseg_ttl_mode_t_to_proto(
    const sai_int32_t val);
sai_outseg_ttl_mode_t convert_sai_outseg_ttl_mode_t_to_sai(
    lemming::dataplane::sai::OutsegTtlMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_outseg_ttl_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_outseg_ttl_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::OutsegType convert_sai_outseg_type_t_to_proto(
    const sai_int32_t val);
sai_outseg_type_t convert_sai_outseg_type_t_to_sai(
    lemming::dataplane::sai::OutsegType val);
google::protobuf::RepeatedField<int> convert_list_sai_outseg_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_outseg_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PacketAction convert_sai_packet_action_t_to_proto(
    const sai_int32_t val);
sai_packet_action_t convert_sai_packet_action_t_to_sai(
    lemming::dataplane::sai::PacketAction val);
google::protobuf::RepeatedField<int> convert_list_sai_packet_action_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_packet_action_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PacketColor convert_sai_packet_color_t_to_proto(
    const sai_int32_t val);
sai_packet_color_t convert_sai_packet_color_t_to_sai(
    lemming::dataplane::sai::PacketColor val);
google::protobuf::RepeatedField<int> convert_list_sai_packet_color_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_packet_color_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PacketVlan convert_sai_packet_vlan_t_to_proto(
    const sai_int32_t val);
sai_packet_vlan_t convert_sai_packet_vlan_t_to_sai(
    lemming::dataplane::sai::PacketVlan val);
google::protobuf::RepeatedField<int> convert_list_sai_packet_vlan_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_packet_vlan_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PoePortActiveChannelType
convert_sai_poe_port_active_channel_type_t_to_proto(const sai_int32_t val);
sai_poe_port_active_channel_type_t
convert_sai_poe_port_active_channel_type_t_to_sai(
    lemming::dataplane::sai::PoePortActiveChannelType val);
google::protobuf::RepeatedField<int>
convert_list_sai_poe_port_active_channel_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_poe_port_active_channel_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PoePortClassMethodType
convert_sai_poe_port_class_method_type_t_to_proto(const sai_int32_t val);
sai_poe_port_class_method_type_t
convert_sai_poe_port_class_method_type_t_to_sai(
    lemming::dataplane::sai::PoePortClassMethodType val);
google::protobuf::RepeatedField<int>
convert_list_sai_poe_port_class_method_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_poe_port_class_method_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PoePortSignatureType
convert_sai_poe_port_signature_type_t_to_proto(const sai_int32_t val);
sai_poe_port_signature_type_t convert_sai_poe_port_signature_type_t_to_sai(
    lemming::dataplane::sai::PoePortSignatureType val);
google::protobuf::RepeatedField<int>
convert_list_sai_poe_port_signature_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_poe_port_signature_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PolicerAttr convert_sai_policer_attr_t_to_proto(
    const sai_int32_t val);
sai_policer_attr_t convert_sai_policer_attr_t_to_sai(
    lemming::dataplane::sai::PolicerAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_policer_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_policer_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PolicerColorSource
convert_sai_policer_color_source_t_to_proto(const sai_int32_t val);
sai_policer_color_source_t convert_sai_policer_color_source_t_to_sai(
    lemming::dataplane::sai::PolicerColorSource val);
google::protobuf::RepeatedField<int>
convert_list_sai_policer_color_source_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_policer_color_source_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PolicerMode convert_sai_policer_mode_t_to_proto(
    const sai_int32_t val);
sai_policer_mode_t convert_sai_policer_mode_t_to_sai(
    lemming::dataplane::sai::PolicerMode val);
google::protobuf::RepeatedField<int> convert_list_sai_policer_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_policer_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PolicerStat convert_sai_policer_stat_t_to_proto(
    const sai_int32_t val);
sai_policer_stat_t convert_sai_policer_stat_t_to_sai(
    lemming::dataplane::sai::PolicerStat val);
google::protobuf::RepeatedField<int> convert_list_sai_policer_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_policer_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortAttrExtensions
convert_sai_port_attr_extensions_t_to_proto(const sai_int32_t val);
sai_port_attr_extensions_t convert_sai_port_attr_extensions_t_to_sai(
    lemming::dataplane::sai::PortAttrExtensions val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_attr_extensions_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_attr_extensions_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortAttr convert_sai_port_attr_t_to_proto(
    const sai_int32_t val);
sai_port_attr_t convert_sai_port_attr_t_to_sai(
    lemming::dataplane::sai::PortAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_port_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortAutoNegConfigMode
convert_sai_port_auto_neg_config_mode_t_to_proto(const sai_int32_t val);
sai_port_auto_neg_config_mode_t convert_sai_port_auto_neg_config_mode_t_to_sai(
    lemming::dataplane::sai::PortAutoNegConfigMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_auto_neg_config_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_auto_neg_config_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortBreakoutModeType
convert_sai_port_breakout_mode_type_t_to_proto(const sai_int32_t val);
sai_port_breakout_mode_type_t convert_sai_port_breakout_mode_type_t_to_sai(
    lemming::dataplane::sai::PortBreakoutModeType val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_breakout_mode_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_breakout_mode_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortCablePairState
convert_sai_port_cable_pair_state_t_to_proto(const sai_int32_t val);
sai_port_cable_pair_state_t convert_sai_port_cable_pair_state_t_to_sai(
    lemming::dataplane::sai::PortCablePairState val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_cable_pair_state_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_cable_pair_state_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortCableType convert_sai_port_cable_type_t_to_proto(
    const sai_int32_t val);
sai_port_cable_type_t convert_sai_port_cable_type_t_to_sai(
    lemming::dataplane::sai::PortCableType val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_cable_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_cable_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortConnectorAttr
convert_sai_port_connector_attr_t_to_proto(const sai_int32_t val);
sai_port_connector_attr_t convert_sai_port_connector_attr_t_to_sai(
    lemming::dataplane::sai::PortConnectorAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_connector_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_connector_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortConnectorFailoverMode
convert_sai_port_connector_failover_mode_t_to_proto(const sai_int32_t val);
sai_port_connector_failover_mode_t
convert_sai_port_connector_failover_mode_t_to_sai(
    lemming::dataplane::sai::PortConnectorFailoverMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_connector_failover_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_connector_failover_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortDatapathEnable
convert_sai_port_datapath_enable_t_to_proto(const sai_int32_t val);
sai_port_datapath_enable_t convert_sai_port_datapath_enable_t_to_sai(
    lemming::dataplane::sai::PortDatapathEnable val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_datapath_enable_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_datapath_enable_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortDualMedia convert_sai_port_dual_media_t_to_proto(
    const sai_int32_t val);
sai_port_dual_media_t convert_sai_port_dual_media_t_to_sai(
    lemming::dataplane::sai::PortDualMedia val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_dual_media_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_dual_media_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortErrStatus convert_sai_port_err_status_t_to_proto(
    const sai_int32_t val);
sai_port_err_status_t convert_sai_port_err_status_t_to_sai(
    lemming::dataplane::sai::PortErrStatus val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_err_status_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_err_status_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortErrorStatus
convert_sai_port_error_status_t_to_proto(const sai_int32_t val);
sai_port_error_status_t convert_sai_port_error_status_t_to_sai(
    lemming::dataplane::sai::PortErrorStatus val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_error_status_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_error_status_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortFecModeExtended
convert_sai_port_fec_mode_extended_t_to_proto(const sai_int32_t val);
sai_port_fec_mode_extended_t convert_sai_port_fec_mode_extended_t_to_sai(
    lemming::dataplane::sai::PortFecModeExtended val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_fec_mode_extended_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_fec_mode_extended_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortFecMode convert_sai_port_fec_mode_t_to_proto(
    const sai_int32_t val);
sai_port_fec_mode_t convert_sai_port_fec_mode_t_to_sai(
    lemming::dataplane::sai::PortFecMode val);
google::protobuf::RepeatedField<int> convert_list_sai_port_fec_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_fec_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortFlowControlMode
convert_sai_port_flow_control_mode_t_to_proto(const sai_int32_t val);
sai_port_flow_control_mode_t convert_sai_port_flow_control_mode_t_to_sai(
    lemming::dataplane::sai::PortFlowControlMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_flow_control_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_flow_control_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortHostTxReadyStatus
convert_sai_port_host_tx_ready_status_t_to_proto(const sai_int32_t val);
sai_port_host_tx_ready_status_t convert_sai_port_host_tx_ready_status_t_to_sai(
    lemming::dataplane::sai::PortHostTxReadyStatus val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_host_tx_ready_status_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_host_tx_ready_status_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortInterfaceType
convert_sai_port_interface_type_t_to_proto(const sai_int32_t val);
sai_port_interface_type_t convert_sai_port_interface_type_t_to_sai(
    lemming::dataplane::sai::PortInterfaceType val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_interface_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_interface_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortInternalLoopbackMode
convert_sai_port_internal_loopback_mode_t_to_proto(const sai_int32_t val);
sai_port_internal_loopback_mode_t
convert_sai_port_internal_loopback_mode_t_to_sai(
    lemming::dataplane::sai::PortInternalLoopbackMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_internal_loopback_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_internal_loopback_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortLinkTrainingFailureStatus
convert_sai_port_link_training_failure_status_t_to_proto(const sai_int32_t val);
sai_port_link_training_failure_status_t
convert_sai_port_link_training_failure_status_t_to_sai(
    lemming::dataplane::sai::PortLinkTrainingFailureStatus val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_link_training_failure_status_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_link_training_failure_status_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortLinkTrainingRxStatus
convert_sai_port_link_training_rx_status_t_to_proto(const sai_int32_t val);
sai_port_link_training_rx_status_t
convert_sai_port_link_training_rx_status_t_to_sai(
    lemming::dataplane::sai::PortLinkTrainingRxStatus val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_link_training_rx_status_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_link_training_rx_status_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortLoopbackMode
convert_sai_port_loopback_mode_t_to_proto(const sai_int32_t val);
sai_port_loopback_mode_t convert_sai_port_loopback_mode_t_to_sai(
    lemming::dataplane::sai::PortLoopbackMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_loopback_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_loopback_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortMdixModeConfig
convert_sai_port_mdix_mode_config_t_to_proto(const sai_int32_t val);
sai_port_mdix_mode_config_t convert_sai_port_mdix_mode_config_t_to_sai(
    lemming::dataplane::sai::PortMdixModeConfig val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_mdix_mode_config_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_mdix_mode_config_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortMdixModeStatus
convert_sai_port_mdix_mode_status_t_to_proto(const sai_int32_t val);
sai_port_mdix_mode_status_t convert_sai_port_mdix_mode_status_t_to_sai(
    lemming::dataplane::sai::PortMdixModeStatus val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_mdix_mode_status_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_mdix_mode_status_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortMediaType convert_sai_port_media_type_t_to_proto(
    const sai_int32_t val);
sai_port_media_type_t convert_sai_port_media_type_t_to_sai(
    lemming::dataplane::sai::PortMediaType val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_media_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_media_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortModuleType convert_sai_port_module_type_t_to_proto(
    const sai_int32_t val);
sai_port_module_type_t convert_sai_port_module_type_t_to_sai(
    lemming::dataplane::sai::PortModuleType val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_module_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_module_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortOperStatus convert_sai_port_oper_status_t_to_proto(
    const sai_int32_t val);
sai_port_oper_status_t convert_sai_port_oper_status_t_to_sai(
    lemming::dataplane::sai::PortOperStatus val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_oper_status_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_oper_status_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortPathTracingTimestampType
convert_sai_port_path_tracing_timestamp_type_t_to_proto(const sai_int32_t val);
sai_port_path_tracing_timestamp_type_t
convert_sai_port_path_tracing_timestamp_type_t_to_sai(
    lemming::dataplane::sai::PortPathTracingTimestampType val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_path_tracing_timestamp_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_path_tracing_timestamp_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortPoolAttr convert_sai_port_pool_attr_t_to_proto(
    const sai_int32_t val);
sai_port_pool_attr_t convert_sai_port_pool_attr_t_to_sai(
    lemming::dataplane::sai::PortPoolAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_port_pool_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_pool_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortPoolStat convert_sai_port_pool_stat_t_to_proto(
    const sai_int32_t val);
sai_port_pool_stat_t convert_sai_port_pool_stat_t_to_sai(
    lemming::dataplane::sai::PortPoolStat val);
google::protobuf::RepeatedField<int> convert_list_sai_port_pool_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_pool_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortPrbsConfig convert_sai_port_prbs_config_t_to_proto(
    const sai_int32_t val);
sai_port_prbs_config_t convert_sai_port_prbs_config_t_to_sai(
    lemming::dataplane::sai::PortPrbsConfig val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_prbs_config_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_prbs_config_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortPrbsRxStatus
convert_sai_port_prbs_rx_status_t_to_proto(const sai_int32_t val);
sai_port_prbs_rx_status_t convert_sai_port_prbs_rx_status_t_to_sai(
    lemming::dataplane::sai::PortPrbsRxStatus val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_prbs_rx_status_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_prbs_rx_status_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortPriorityFlowControlMode
convert_sai_port_priority_flow_control_mode_t_to_proto(const sai_int32_t val);
sai_port_priority_flow_control_mode_t
convert_sai_port_priority_flow_control_mode_t_to_sai(
    lemming::dataplane::sai::PortPriorityFlowControlMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_priority_flow_control_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_priority_flow_control_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortPtpMode convert_sai_port_ptp_mode_t_to_proto(
    const sai_int32_t val);
sai_port_ptp_mode_t convert_sai_port_ptp_mode_t_to_sai(
    lemming::dataplane::sai::PortPtpMode val);
google::protobuf::RepeatedField<int> convert_list_sai_port_ptp_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_ptp_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortSerdesAttr convert_sai_port_serdes_attr_t_to_proto(
    const sai_int32_t val);
sai_port_serdes_attr_t convert_sai_port_serdes_attr_t_to_sai(
    lemming::dataplane::sai::PortSerdesAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_serdes_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_serdes_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortStatExtensions
convert_sai_port_stat_extensions_t_to_proto(const sai_int32_t val);
sai_port_stat_extensions_t convert_sai_port_stat_extensions_t_to_sai(
    lemming::dataplane::sai::PortStatExtensions val);
google::protobuf::RepeatedField<int>
convert_list_sai_port_stat_extensions_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_port_stat_extensions_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortStat convert_sai_port_stat_t_to_proto(
    const sai_int32_t val);
sai_port_stat_t convert_sai_port_stat_t_to_sai(
    lemming::dataplane::sai::PortStat val);
google::protobuf::RepeatedField<int> convert_list_sai_port_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::PortType convert_sai_port_type_t_to_proto(
    const sai_int32_t val);
sai_port_type_t convert_sai_port_type_t_to_sai(
    lemming::dataplane::sai::PortType val);
google::protobuf::RepeatedField<int> convert_list_sai_port_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_port_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::QosMapAttr convert_sai_qos_map_attr_t_to_proto(
    const sai_int32_t val);
sai_qos_map_attr_t convert_sai_qos_map_attr_t_to_sai(
    lemming::dataplane::sai::QosMapAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_qos_map_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_qos_map_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::QosMapType convert_sai_qos_map_type_t_to_proto(
    const sai_int32_t val);
sai_qos_map_type_t convert_sai_qos_map_type_t_to_sai(
    lemming::dataplane::sai::QosMapType val);
google::protobuf::RepeatedField<int> convert_list_sai_qos_map_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_qos_map_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::QueueAttr convert_sai_queue_attr_t_to_proto(
    const sai_int32_t val);
sai_queue_attr_t convert_sai_queue_attr_t_to_sai(
    lemming::dataplane::sai::QueueAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_queue_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_queue_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::QueuePfcContinuousDeadlockState
convert_sai_queue_pfc_continuous_deadlock_state_t_to_proto(
    const sai_int32_t val);
sai_queue_pfc_continuous_deadlock_state_t
convert_sai_queue_pfc_continuous_deadlock_state_t_to_sai(
    lemming::dataplane::sai::QueuePfcContinuousDeadlockState val);
google::protobuf::RepeatedField<int>
convert_list_sai_queue_pfc_continuous_deadlock_state_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_queue_pfc_continuous_deadlock_state_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::QueuePfcDeadlockEventType
convert_sai_queue_pfc_deadlock_event_type_t_to_proto(const sai_int32_t val);
sai_queue_pfc_deadlock_event_type_t
convert_sai_queue_pfc_deadlock_event_type_t_to_sai(
    lemming::dataplane::sai::QueuePfcDeadlockEventType val);
google::protobuf::RepeatedField<int>
convert_list_sai_queue_pfc_deadlock_event_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_queue_pfc_deadlock_event_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::QueueStat convert_sai_queue_stat_t_to_proto(
    const sai_int32_t val);
sai_queue_stat_t convert_sai_queue_stat_t_to_sai(
    lemming::dataplane::sai::QueueStat val);
google::protobuf::RepeatedField<int> convert_list_sai_queue_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_queue_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::QueueType convert_sai_queue_type_t_to_proto(
    const sai_int32_t val);
sai_queue_type_t convert_sai_queue_type_t_to_sai(
    lemming::dataplane::sai::QueueType val);
google::protobuf::RepeatedField<int> convert_list_sai_queue_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_queue_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::RouteEntryAttr convert_sai_route_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_route_entry_attr_t convert_sai_route_entry_attr_t_to_sai(
    lemming::dataplane::sai::RouteEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_route_entry_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_route_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::RouterInterfaceAttr
convert_sai_router_interface_attr_t_to_proto(const sai_int32_t val);
sai_router_interface_attr_t convert_sai_router_interface_attr_t_to_sai(
    lemming::dataplane::sai::RouterInterfaceAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_router_interface_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_router_interface_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::RouterInterfaceStat
convert_sai_router_interface_stat_t_to_proto(const sai_int32_t val);
sai_router_interface_stat_t convert_sai_router_interface_stat_t_to_sai(
    lemming::dataplane::sai::RouterInterfaceStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_router_interface_stat_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_router_interface_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::RouterInterfaceType
convert_sai_router_interface_type_t_to_proto(const sai_int32_t val);
sai_router_interface_type_t convert_sai_router_interface_type_t_to_sai(
    lemming::dataplane::sai::RouterInterfaceType val);
google::protobuf::RepeatedField<int>
convert_list_sai_router_interface_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_router_interface_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::RpfGroupAttr convert_sai_rpf_group_attr_t_to_proto(
    const sai_int32_t val);
sai_rpf_group_attr_t convert_sai_rpf_group_attr_t_to_sai(
    lemming::dataplane::sai::RpfGroupAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_rpf_group_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_rpf_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::RpfGroupMemberAttr
convert_sai_rpf_group_member_attr_t_to_proto(const sai_int32_t val);
sai_rpf_group_member_attr_t convert_sai_rpf_group_member_attr_t_to_sai(
    lemming::dataplane::sai::RpfGroupMemberAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_rpf_group_member_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_rpf_group_member_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SamplepacketAttr
convert_sai_samplepacket_attr_t_to_proto(const sai_int32_t val);
sai_samplepacket_attr_t convert_sai_samplepacket_attr_t_to_sai(
    lemming::dataplane::sai::SamplepacketAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_samplepacket_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_samplepacket_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SamplepacketMode
convert_sai_samplepacket_mode_t_to_proto(const sai_int32_t val);
sai_samplepacket_mode_t convert_sai_samplepacket_mode_t_to_sai(
    lemming::dataplane::sai::SamplepacketMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_samplepacket_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_samplepacket_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SamplepacketType
convert_sai_samplepacket_type_t_to_proto(const sai_int32_t val);
sai_samplepacket_type_t convert_sai_samplepacket_type_t_to_sai(
    lemming::dataplane::sai::SamplepacketType val);
google::protobuf::RepeatedField<int>
convert_list_sai_samplepacket_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_samplepacket_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SchedulerAttr convert_sai_scheduler_attr_t_to_proto(
    const sai_int32_t val);
sai_scheduler_attr_t convert_sai_scheduler_attr_t_to_sai(
    lemming::dataplane::sai::SchedulerAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_scheduler_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_scheduler_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SchedulerGroupAttr
convert_sai_scheduler_group_attr_t_to_proto(const sai_int32_t val);
sai_scheduler_group_attr_t convert_sai_scheduler_group_attr_t_to_sai(
    lemming::dataplane::sai::SchedulerGroupAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_scheduler_group_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_scheduler_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SchedulingType convert_sai_scheduling_type_t_to_proto(
    const sai_int32_t val);
sai_scheduling_type_t convert_sai_scheduling_type_t_to_sai(
    lemming::dataplane::sai::SchedulingType val);
google::protobuf::RepeatedField<int>
convert_list_sai_scheduling_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_scheduling_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SerCorrectionType
convert_sai_ser_correction_type_t_to_proto(const sai_int32_t val);
sai_ser_correction_type_t convert_sai_ser_correction_type_t_to_sai(
    lemming::dataplane::sai::SerCorrectionType val);
google::protobuf::RepeatedField<int>
convert_list_sai_ser_correction_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_ser_correction_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SerLogType convert_sai_ser_log_type_t_to_proto(
    const sai_int32_t val);
sai_ser_log_type_t convert_sai_ser_log_type_t_to_sai(
    lemming::dataplane::sai::SerLogType val);
google::protobuf::RepeatedField<int> convert_list_sai_ser_log_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ser_log_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SerType convert_sai_ser_type_t_to_proto(
    const sai_int32_t val);
sai_ser_type_t convert_sai_ser_type_t_to_sai(
    lemming::dataplane::sai::SerType val);
google::protobuf::RepeatedField<int> convert_list_sai_ser_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_ser_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::Srv6SidlistAttr
convert_sai_srv6_sidlist_attr_t_to_proto(const sai_int32_t val);
sai_srv6_sidlist_attr_t convert_sai_srv6_sidlist_attr_t_to_sai(
    lemming::dataplane::sai::Srv6SidlistAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_srv6_sidlist_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_srv6_sidlist_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::Srv6SidlistStat
convert_sai_srv6_sidlist_stat_t_to_proto(const sai_int32_t val);
sai_srv6_sidlist_stat_t convert_sai_srv6_sidlist_stat_t_to_sai(
    lemming::dataplane::sai::Srv6SidlistStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_srv6_sidlist_stat_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_srv6_sidlist_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::Srv6SidlistType
convert_sai_srv6_sidlist_type_t_to_proto(const sai_int32_t val);
sai_srv6_sidlist_type_t convert_sai_srv6_sidlist_type_t_to_sai(
    lemming::dataplane::sai::Srv6SidlistType val);
google::protobuf::RepeatedField<int>
convert_list_sai_srv6_sidlist_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_srv6_sidlist_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::StatsCountMode convert_sai_stats_count_mode_t_to_proto(
    const sai_int32_t val);
sai_stats_count_mode_t convert_sai_stats_count_mode_t_to_sai(
    lemming::dataplane::sai::StatsCountMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_stats_count_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_stats_count_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::StatsMode convert_sai_stats_mode_t_to_proto(
    const sai_int32_t val);
sai_stats_mode_t convert_sai_stats_mode_t_to_sai(
    lemming::dataplane::sai::StatsMode val);
google::protobuf::RepeatedField<int> convert_list_sai_stats_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_stats_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::StpAttr convert_sai_stp_attr_t_to_proto(
    const sai_int32_t val);
sai_stp_attr_t convert_sai_stp_attr_t_to_sai(
    lemming::dataplane::sai::StpAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_stp_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_stp_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::StpPortAttr convert_sai_stp_port_attr_t_to_proto(
    const sai_int32_t val);
sai_stp_port_attr_t convert_sai_stp_port_attr_t_to_sai(
    lemming::dataplane::sai::StpPortAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_stp_port_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_stp_port_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::StpPortState convert_sai_stp_port_state_t_to_proto(
    const sai_int32_t val);
sai_stp_port_state_t convert_sai_stp_port_state_t_to_sai(
    lemming::dataplane::sai::StpPortState val);
google::protobuf::RepeatedField<int> convert_list_sai_stp_port_state_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_stp_port_state_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchAsicSdkHealthCategory
convert_sai_switch_asic_sdk_health_category_t_to_proto(const sai_int32_t val);
sai_switch_asic_sdk_health_category_t
convert_sai_switch_asic_sdk_health_category_t_to_sai(
    lemming::dataplane::sai::SwitchAsicSdkHealthCategory val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_asic_sdk_health_category_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_asic_sdk_health_category_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchAsicSdkHealthSeverity
convert_sai_switch_asic_sdk_health_severity_t_to_proto(const sai_int32_t val);
sai_switch_asic_sdk_health_severity_t
convert_sai_switch_asic_sdk_health_severity_t_to_sai(
    lemming::dataplane::sai::SwitchAsicSdkHealthSeverity val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_asic_sdk_health_severity_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_asic_sdk_health_severity_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchAttrExtensions
convert_sai_switch_attr_extensions_t_to_proto(const sai_int32_t val);
sai_switch_attr_extensions_t convert_sai_switch_attr_extensions_t_to_sai(
    lemming::dataplane::sai::SwitchAttrExtensions val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_attr_extensions_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_switch_attr_extensions_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchAttr convert_sai_switch_attr_t_to_proto(
    const sai_int32_t val);
sai_switch_attr_t convert_sai_switch_attr_t_to_sai(
    lemming::dataplane::sai::SwitchAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_switch_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchFailoverConfigMode
convert_sai_switch_failover_config_mode_t_to_proto(const sai_int32_t val);
sai_switch_failover_config_mode_t
convert_sai_switch_failover_config_mode_t_to_sai(
    lemming::dataplane::sai::SwitchFailoverConfigMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_failover_config_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_failover_config_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchFirmwareLoadMethod
convert_sai_switch_firmware_load_method_t_to_proto(const sai_int32_t val);
sai_switch_firmware_load_method_t
convert_sai_switch_firmware_load_method_t_to_sai(
    lemming::dataplane::sai::SwitchFirmwareLoadMethod val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_firmware_load_method_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_firmware_load_method_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchFirmwareLoadType
convert_sai_switch_firmware_load_type_t_to_proto(const sai_int32_t val);
sai_switch_firmware_load_type_t convert_sai_switch_firmware_load_type_t_to_sai(
    lemming::dataplane::sai::SwitchFirmwareLoadType val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_firmware_load_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_firmware_load_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchHardwareAccessBus
convert_sai_switch_hardware_access_bus_t_to_proto(const sai_int32_t val);
sai_switch_hardware_access_bus_t
convert_sai_switch_hardware_access_bus_t_to_sai(
    lemming::dataplane::sai::SwitchHardwareAccessBus val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_hardware_access_bus_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_hardware_access_bus_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchHostifOperStatusUpdateMode
convert_sai_switch_hostif_oper_status_update_mode_t_to_proto(
    const sai_int32_t val);
sai_switch_hostif_oper_status_update_mode_t
convert_sai_switch_hostif_oper_status_update_mode_t_to_sai(
    lemming::dataplane::sai::SwitchHostifOperStatusUpdateMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_hostif_oper_status_update_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_hostif_oper_status_update_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchMcastSnoopingCapability
convert_sai_switch_mcast_snooping_capability_t_to_proto(const sai_int32_t val);
sai_switch_mcast_snooping_capability_t
convert_sai_switch_mcast_snooping_capability_t_to_sai(
    lemming::dataplane::sai::SwitchMcastSnoopingCapability val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_mcast_snooping_capability_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_mcast_snooping_capability_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchOperStatus
convert_sai_switch_oper_status_t_to_proto(const sai_int32_t val);
sai_switch_oper_status_t convert_sai_switch_oper_status_t_to_sai(
    lemming::dataplane::sai::SwitchOperStatus val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_oper_status_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_switch_oper_status_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchRestartType
convert_sai_switch_restart_type_t_to_proto(const sai_int32_t val);
sai_switch_restart_type_t convert_sai_switch_restart_type_t_to_sai(
    lemming::dataplane::sai::SwitchRestartType val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_restart_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_switch_restart_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchStat convert_sai_switch_stat_t_to_proto(
    const sai_int32_t val);
sai_switch_stat_t convert_sai_switch_stat_t_to_sai(
    lemming::dataplane::sai::SwitchStat val);
google::protobuf::RepeatedField<int> convert_list_sai_switch_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchSwitchingMode
convert_sai_switch_switching_mode_t_to_proto(const sai_int32_t val);
sai_switch_switching_mode_t convert_sai_switch_switching_mode_t_to_sai(
    lemming::dataplane::sai::SwitchSwitchingMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_switching_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_switch_switching_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchTunnelAttr
convert_sai_switch_tunnel_attr_t_to_proto(const sai_int32_t val);
sai_switch_tunnel_attr_t convert_sai_switch_tunnel_attr_t_to_sai(
    lemming::dataplane::sai::SwitchTunnelAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_switch_tunnel_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_switch_tunnel_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SwitchType convert_sai_switch_type_t_to_proto(
    const sai_int32_t val);
sai_switch_type_t convert_sai_switch_type_t_to_sai(
    lemming::dataplane::sai::SwitchType val);
google::protobuf::RepeatedField<int> convert_list_sai_switch_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_switch_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SystemPortAttr convert_sai_system_port_attr_t_to_proto(
    const sai_int32_t val);
sai_system_port_attr_t convert_sai_system_port_attr_t_to_sai(
    lemming::dataplane::sai::SystemPortAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_system_port_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_system_port_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::SystemPortType convert_sai_system_port_type_t_to_proto(
    const sai_int32_t val);
sai_system_port_type_t convert_sai_system_port_type_t_to_sai(
    lemming::dataplane::sai::SystemPortType val);
google::protobuf::RepeatedField<int>
convert_list_sai_system_port_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_system_port_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TableBitmapClassificationEntryAction
convert_sai_table_bitmap_classification_entry_action_t_to_proto(
    const sai_int32_t val);
sai_table_bitmap_classification_entry_action_t
convert_sai_table_bitmap_classification_entry_action_t_to_sai(
    lemming::dataplane::sai::TableBitmapClassificationEntryAction val);
google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_classification_entry_action_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_table_bitmap_classification_entry_action_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TableBitmapClassificationEntryAttr
convert_sai_table_bitmap_classification_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_table_bitmap_classification_entry_attr_t
convert_sai_table_bitmap_classification_entry_attr_t_to_sai(
    lemming::dataplane::sai::TableBitmapClassificationEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_classification_entry_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_table_bitmap_classification_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TableBitmapClassificationEntryStat
convert_sai_table_bitmap_classification_entry_stat_t_to_proto(
    const sai_int32_t val);
sai_table_bitmap_classification_entry_stat_t
convert_sai_table_bitmap_classification_entry_stat_t_to_sai(
    lemming::dataplane::sai::TableBitmapClassificationEntryStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_classification_entry_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_table_bitmap_classification_entry_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TableBitmapRouterEntryAction
convert_sai_table_bitmap_router_entry_action_t_to_proto(const sai_int32_t val);
sai_table_bitmap_router_entry_action_t
convert_sai_table_bitmap_router_entry_action_t_to_sai(
    lemming::dataplane::sai::TableBitmapRouterEntryAction val);
google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_router_entry_action_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_table_bitmap_router_entry_action_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TableBitmapRouterEntryAttr
convert_sai_table_bitmap_router_entry_attr_t_to_proto(const sai_int32_t val);
sai_table_bitmap_router_entry_attr_t
convert_sai_table_bitmap_router_entry_attr_t_to_sai(
    lemming::dataplane::sai::TableBitmapRouterEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_router_entry_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_table_bitmap_router_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TableBitmapRouterEntryStat
convert_sai_table_bitmap_router_entry_stat_t_to_proto(const sai_int32_t val);
sai_table_bitmap_router_entry_stat_t
convert_sai_table_bitmap_router_entry_stat_t_to_sai(
    lemming::dataplane::sai::TableBitmapRouterEntryStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_router_entry_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_table_bitmap_router_entry_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TableMetaTunnelEntryAction
convert_sai_table_meta_tunnel_entry_action_t_to_proto(const sai_int32_t val);
sai_table_meta_tunnel_entry_action_t
convert_sai_table_meta_tunnel_entry_action_t_to_sai(
    lemming::dataplane::sai::TableMetaTunnelEntryAction val);
google::protobuf::RepeatedField<int>
convert_list_sai_table_meta_tunnel_entry_action_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_table_meta_tunnel_entry_action_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TableMetaTunnelEntryAttr
convert_sai_table_meta_tunnel_entry_attr_t_to_proto(const sai_int32_t val);
sai_table_meta_tunnel_entry_attr_t
convert_sai_table_meta_tunnel_entry_attr_t_to_sai(
    lemming::dataplane::sai::TableMetaTunnelEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_table_meta_tunnel_entry_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_table_meta_tunnel_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TableMetaTunnelEntryStat
convert_sai_table_meta_tunnel_entry_stat_t_to_proto(const sai_int32_t val);
sai_table_meta_tunnel_entry_stat_t
convert_sai_table_meta_tunnel_entry_stat_t_to_sai(
    lemming::dataplane::sai::TableMetaTunnelEntryStat val);
google::protobuf::RepeatedField<int>
convert_list_sai_table_meta_tunnel_entry_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_table_meta_tunnel_entry_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamAttr convert_sai_tam_attr_t_to_proto(
    const sai_int32_t val);
sai_tam_attr_t convert_sai_tam_attr_t_to_sai(
    lemming::dataplane::sai::TamAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_tam_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tam_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamBindPointType
convert_sai_tam_bind_point_type_t_to_proto(const sai_int32_t val);
sai_tam_bind_point_type_t convert_sai_tam_bind_point_type_t_to_sai(
    lemming::dataplane::sai::TamBindPointType val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_bind_point_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_bind_point_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamCollectorAttr
convert_sai_tam_collector_attr_t_to_proto(const sai_int32_t val);
sai_tam_collector_attr_t convert_sai_tam_collector_attr_t_to_sai(
    lemming::dataplane::sai::TamCollectorAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_collector_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_collector_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamCounterSubscriptionAttr
convert_sai_tam_counter_subscription_attr_t_to_proto(const sai_int32_t val);
sai_tam_counter_subscription_attr_t
convert_sai_tam_counter_subscription_attr_t_to_sai(
    lemming::dataplane::sai::TamCounterSubscriptionAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_counter_subscription_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tam_counter_subscription_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamEventActionAttr
convert_sai_tam_event_action_attr_t_to_proto(const sai_int32_t val);
sai_tam_event_action_attr_t convert_sai_tam_event_action_attr_t_to_sai(
    lemming::dataplane::sai::TamEventActionAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_event_action_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_event_action_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamEventAttr convert_sai_tam_event_attr_t_to_proto(
    const sai_int32_t val);
sai_tam_event_attr_t convert_sai_tam_event_attr_t_to_sai(
    lemming::dataplane::sai::TamEventAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_tam_event_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tam_event_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamEventThresholdAttr
convert_sai_tam_event_threshold_attr_t_to_proto(const sai_int32_t val);
sai_tam_event_threshold_attr_t convert_sai_tam_event_threshold_attr_t_to_sai(
    lemming::dataplane::sai::TamEventThresholdAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_event_threshold_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tam_event_threshold_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamEventThresholdUnit
convert_sai_tam_event_threshold_unit_t_to_proto(const sai_int32_t val);
sai_tam_event_threshold_unit_t convert_sai_tam_event_threshold_unit_t_to_sai(
    lemming::dataplane::sai::TamEventThresholdUnit val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_event_threshold_unit_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tam_event_threshold_unit_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamEventType convert_sai_tam_event_type_t_to_proto(
    const sai_int32_t val);
sai_tam_event_type_t convert_sai_tam_event_type_t_to_sai(
    lemming::dataplane::sai::TamEventType val);
google::protobuf::RepeatedField<int> convert_list_sai_tam_event_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tam_event_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamIntAttr convert_sai_tam_int_attr_t_to_proto(
    const sai_int32_t val);
sai_tam_int_attr_t convert_sai_tam_int_attr_t_to_sai(
    lemming::dataplane::sai::TamIntAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_tam_int_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tam_int_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamIntPresenceType
convert_sai_tam_int_presence_type_t_to_proto(const sai_int32_t val);
sai_tam_int_presence_type_t convert_sai_tam_int_presence_type_t_to_sai(
    lemming::dataplane::sai::TamIntPresenceType val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_int_presence_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_int_presence_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamIntType convert_sai_tam_int_type_t_to_proto(
    const sai_int32_t val);
sai_tam_int_type_t convert_sai_tam_int_type_t_to_sai(
    lemming::dataplane::sai::TamIntType val);
google::protobuf::RepeatedField<int> convert_list_sai_tam_int_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tam_int_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamMathFuncAttr
convert_sai_tam_math_func_attr_t_to_proto(const sai_int32_t val);
sai_tam_math_func_attr_t convert_sai_tam_math_func_attr_t_to_sai(
    lemming::dataplane::sai::TamMathFuncAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_math_func_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_math_func_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamReportAttr convert_sai_tam_report_attr_t_to_proto(
    const sai_int32_t val);
sai_tam_report_attr_t convert_sai_tam_report_attr_t_to_sai(
    lemming::dataplane::sai::TamReportAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_report_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_report_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamReportIntervalUnit
convert_sai_tam_report_interval_unit_t_to_proto(const sai_int32_t val);
sai_tam_report_interval_unit_t convert_sai_tam_report_interval_unit_t_to_sai(
    lemming::dataplane::sai::TamReportIntervalUnit val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_report_interval_unit_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tam_report_interval_unit_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamReportMode convert_sai_tam_report_mode_t_to_proto(
    const sai_int32_t val);
sai_tam_report_mode_t convert_sai_tam_report_mode_t_to_sai(
    lemming::dataplane::sai::TamReportMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_report_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_report_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamReportType convert_sai_tam_report_type_t_to_proto(
    const sai_int32_t val);
sai_tam_report_type_t convert_sai_tam_report_type_t_to_sai(
    lemming::dataplane::sai::TamReportType val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_report_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_report_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamReportingUnit
convert_sai_tam_reporting_unit_t_to_proto(const sai_int32_t val);
sai_tam_reporting_unit_t convert_sai_tam_reporting_unit_t_to_sai(
    lemming::dataplane::sai::TamReportingUnit val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_reporting_unit_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_reporting_unit_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamTelMathFuncType
convert_sai_tam_tel_math_func_type_t_to_proto(const sai_int32_t val);
sai_tam_tel_math_func_type_t convert_sai_tam_tel_math_func_type_t_to_sai(
    lemming::dataplane::sai::TamTelMathFuncType val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_tel_math_func_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_tel_math_func_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamTelTypeAttr
convert_sai_tam_tel_type_attr_t_to_proto(const sai_int32_t val);
sai_tam_tel_type_attr_t convert_sai_tam_tel_type_attr_t_to_sai(
    lemming::dataplane::sai::TamTelTypeAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_tel_type_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_tel_type_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamTelemetryAttr
convert_sai_tam_telemetry_attr_t_to_proto(const sai_int32_t val);
sai_tam_telemetry_attr_t convert_sai_tam_telemetry_attr_t_to_sai(
    lemming::dataplane::sai::TamTelemetryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_telemetry_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_telemetry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamTelemetryType
convert_sai_tam_telemetry_type_t_to_proto(const sai_int32_t val);
sai_tam_telemetry_type_t convert_sai_tam_telemetry_type_t_to_sai(
    lemming::dataplane::sai::TamTelemetryType val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_telemetry_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_telemetry_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamTransportAttr
convert_sai_tam_transport_attr_t_to_proto(const sai_int32_t val);
sai_tam_transport_attr_t convert_sai_tam_transport_attr_t_to_sai(
    lemming::dataplane::sai::TamTransportAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_transport_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_transport_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamTransportAuthType
convert_sai_tam_transport_auth_type_t_to_proto(const sai_int32_t val);
sai_tam_transport_auth_type_t convert_sai_tam_transport_auth_type_t_to_sai(
    lemming::dataplane::sai::TamTransportAuthType val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_transport_auth_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_transport_auth_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TamTransportType
convert_sai_tam_transport_type_t_to_proto(const sai_int32_t val);
sai_tam_transport_type_t convert_sai_tam_transport_type_t_to_sai(
    lemming::dataplane::sai::TamTransportType val);
google::protobuf::RepeatedField<int>
convert_list_sai_tam_transport_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tam_transport_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TlvType convert_sai_tlv_type_t_to_proto(
    const sai_int32_t val);
sai_tlv_type_t convert_sai_tlv_type_t_to_sai(
    lemming::dataplane::sai::TlvType val);
google::protobuf::RepeatedField<int> convert_list_sai_tlv_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tlv_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelAttr convert_sai_tunnel_attr_t_to_proto(
    const sai_int32_t val);
sai_tunnel_attr_t convert_sai_tunnel_attr_t_to_sai(
    lemming::dataplane::sai::TunnelAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_tunnel_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tunnel_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelDecapEcnMode
convert_sai_tunnel_decap_ecn_mode_t_to_proto(const sai_int32_t val);
sai_tunnel_decap_ecn_mode_t convert_sai_tunnel_decap_ecn_mode_t_to_sai(
    lemming::dataplane::sai::TunnelDecapEcnMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_decap_ecn_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tunnel_decap_ecn_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelDscpMode convert_sai_tunnel_dscp_mode_t_to_proto(
    const sai_int32_t val);
sai_tunnel_dscp_mode_t convert_sai_tunnel_dscp_mode_t_to_sai(
    lemming::dataplane::sai::TunnelDscpMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_dscp_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tunnel_dscp_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelEncapEcnMode
convert_sai_tunnel_encap_ecn_mode_t_to_proto(const sai_int32_t val);
sai_tunnel_encap_ecn_mode_t convert_sai_tunnel_encap_ecn_mode_t_to_sai(
    lemming::dataplane::sai::TunnelEncapEcnMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_encap_ecn_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tunnel_encap_ecn_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelMapAttr convert_sai_tunnel_map_attr_t_to_proto(
    const sai_int32_t val);
sai_tunnel_map_attr_t convert_sai_tunnel_map_attr_t_to_sai(
    lemming::dataplane::sai::TunnelMapAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_map_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tunnel_map_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelMapEntryAttr
convert_sai_tunnel_map_entry_attr_t_to_proto(const sai_int32_t val);
sai_tunnel_map_entry_attr_t convert_sai_tunnel_map_entry_attr_t_to_sai(
    lemming::dataplane::sai::TunnelMapEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_map_entry_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tunnel_map_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelMapType convert_sai_tunnel_map_type_t_to_proto(
    const sai_int32_t val);
sai_tunnel_map_type_t convert_sai_tunnel_map_type_t_to_sai(
    lemming::dataplane::sai::TunnelMapType val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_map_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tunnel_map_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelPeerMode convert_sai_tunnel_peer_mode_t_to_proto(
    const sai_int32_t val);
sai_tunnel_peer_mode_t convert_sai_tunnel_peer_mode_t_to_sai(
    lemming::dataplane::sai::TunnelPeerMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_peer_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tunnel_peer_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelStat convert_sai_tunnel_stat_t_to_proto(
    const sai_int32_t val);
sai_tunnel_stat_t convert_sai_tunnel_stat_t_to_sai(
    lemming::dataplane::sai::TunnelStat val);
google::protobuf::RepeatedField<int> convert_list_sai_tunnel_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tunnel_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelTermTableEntryAttr
convert_sai_tunnel_term_table_entry_attr_t_to_proto(const sai_int32_t val);
sai_tunnel_term_table_entry_attr_t
convert_sai_tunnel_term_table_entry_attr_t_to_sai(
    lemming::dataplane::sai::TunnelTermTableEntryAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_term_table_entry_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tunnel_term_table_entry_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelTermTableEntryType
convert_sai_tunnel_term_table_entry_type_t_to_proto(const sai_int32_t val);
sai_tunnel_term_table_entry_type_t
convert_sai_tunnel_term_table_entry_type_t_to_sai(
    lemming::dataplane::sai::TunnelTermTableEntryType val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_term_table_entry_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tunnel_term_table_entry_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelTtlMode convert_sai_tunnel_ttl_mode_t_to_proto(
    const sai_int32_t val);
sai_tunnel_ttl_mode_t convert_sai_tunnel_ttl_mode_t_to_sai(
    lemming::dataplane::sai::TunnelTtlMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_ttl_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_tunnel_ttl_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelType convert_sai_tunnel_type_t_to_proto(
    const sai_int32_t val);
sai_tunnel_type_t convert_sai_tunnel_type_t_to_sai(
    lemming::dataplane::sai::TunnelType val);
google::protobuf::RepeatedField<int> convert_list_sai_tunnel_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tunnel_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::TunnelVxlanUdpSportMode
convert_sai_tunnel_vxlan_udp_sport_mode_t_to_proto(const sai_int32_t val);
sai_tunnel_vxlan_udp_sport_mode_t
convert_sai_tunnel_vxlan_udp_sport_mode_t_to_sai(
    lemming::dataplane::sai::TunnelVxlanUdpSportMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_vxlan_udp_sport_mode_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_tunnel_vxlan_udp_sport_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::UdfAttr convert_sai_udf_attr_t_to_proto(
    const sai_int32_t val);
sai_udf_attr_t convert_sai_udf_attr_t_to_sai(
    lemming::dataplane::sai::UdfAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_udf_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_udf_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::UdfBase convert_sai_udf_base_t_to_proto(
    const sai_int32_t val);
sai_udf_base_t convert_sai_udf_base_t_to_sai(
    lemming::dataplane::sai::UdfBase val);
google::protobuf::RepeatedField<int> convert_list_sai_udf_base_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_udf_base_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::UdfGroupAttr convert_sai_udf_group_attr_t_to_proto(
    const sai_int32_t val);
sai_udf_group_attr_t convert_sai_udf_group_attr_t_to_sai(
    lemming::dataplane::sai::UdfGroupAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_udf_group_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_udf_group_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::UdfGroupType convert_sai_udf_group_type_t_to_proto(
    const sai_int32_t val);
sai_udf_group_type_t convert_sai_udf_group_type_t_to_sai(
    lemming::dataplane::sai::UdfGroupType val);
google::protobuf::RepeatedField<int> convert_list_sai_udf_group_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_udf_group_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::UdfMatchAttr convert_sai_udf_match_attr_t_to_proto(
    const sai_int32_t val);
sai_udf_match_attr_t convert_sai_udf_match_attr_t_to_sai(
    lemming::dataplane::sai::UdfMatchAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_udf_match_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_udf_match_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::VirtualRouterAttr
convert_sai_virtual_router_attr_t_to_proto(const sai_int32_t val);
sai_virtual_router_attr_t convert_sai_virtual_router_attr_t_to_sai(
    lemming::dataplane::sai::VirtualRouterAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_virtual_router_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_virtual_router_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::VlanAttr convert_sai_vlan_attr_t_to_proto(
    const sai_int32_t val);
sai_vlan_attr_t convert_sai_vlan_attr_t_to_sai(
    lemming::dataplane::sai::VlanAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_vlan_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_vlan_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::VlanFloodControlType
convert_sai_vlan_flood_control_type_t_to_proto(const sai_int32_t val);
sai_vlan_flood_control_type_t convert_sai_vlan_flood_control_type_t_to_sai(
    lemming::dataplane::sai::VlanFloodControlType val);
google::protobuf::RepeatedField<int>
convert_list_sai_vlan_flood_control_type_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_vlan_flood_control_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::VlanMcastLookupKeyType
convert_sai_vlan_mcast_lookup_key_type_t_to_proto(const sai_int32_t val);
sai_vlan_mcast_lookup_key_type_t
convert_sai_vlan_mcast_lookup_key_type_t_to_sai(
    lemming::dataplane::sai::VlanMcastLookupKeyType val);
google::protobuf::RepeatedField<int>
convert_list_sai_vlan_mcast_lookup_key_type_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_vlan_mcast_lookup_key_type_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::VlanMemberAttr convert_sai_vlan_member_attr_t_to_proto(
    const sai_int32_t val);
sai_vlan_member_attr_t convert_sai_vlan_member_attr_t_to_sai(
    lemming::dataplane::sai::VlanMemberAttr val);
google::protobuf::RepeatedField<int>
convert_list_sai_vlan_member_attr_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_vlan_member_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::VlanStat convert_sai_vlan_stat_t_to_proto(
    const sai_int32_t val);
sai_vlan_stat_t convert_sai_vlan_stat_t_to_sai(
    lemming::dataplane::sai::VlanStat val);
google::protobuf::RepeatedField<int> convert_list_sai_vlan_stat_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_vlan_stat_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::VlanTaggingMode
convert_sai_vlan_tagging_mode_t_to_proto(const sai_int32_t val);
sai_vlan_tagging_mode_t convert_sai_vlan_tagging_mode_t_to_sai(
    lemming::dataplane::sai::VlanTaggingMode val);
google::protobuf::RepeatedField<int>
convert_list_sai_vlan_tagging_mode_t_to_proto(const sai_s32_list_t& list);
void convert_list_sai_vlan_tagging_mode_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

lemming::dataplane::sai::WredAttr convert_sai_wred_attr_t_to_proto(
    const sai_int32_t val);
sai_wred_attr_t convert_sai_wred_attr_t_to_sai(
    lemming::dataplane::sai::WredAttr val);
google::protobuf::RepeatedField<int> convert_list_sai_wred_attr_t_to_proto(
    const sai_s32_list_t& list);
void convert_list_sai_wred_attr_t_to_sai(
    int32_t* list, const google::protobuf::RepeatedField<int>& proto_list,
    uint32_t* count);

#endif  // DATAPLANE_STANDALONE_SAI_ENUM_H_
