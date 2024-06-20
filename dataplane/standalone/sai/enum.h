
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
#include "dataplane/proto/sai/bfd.pb.h"
#include "dataplane/proto/sai/bmtor.pb.h"
#include "dataplane/proto/sai/bridge.pb.h"
#include "dataplane/proto/sai/buffer.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/counter.pb.h"
#include "dataplane/proto/sai/debug_counter.pb.h"
#include "dataplane/proto/sai/dtel.pb.h"
#include "dataplane/proto/sai/fdb.pb.h"
#include "dataplane/proto/sai/generic_programmable.pb.h"
#include "dataplane/proto/sai/hash.pb.h"
#include "dataplane/proto/sai/hostif.pb.h"
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
#include "dataplane/proto/sai/udf.pb.h"
#include "dataplane/proto/sai/virtual_router.pb.h"
#include "dataplane/proto/sai/vlan.pb.h"
#include "dataplane/proto/sai/wred.pb.h"

extern "C" {
#include "inc/sai.h"
#include "experimental/saiextensions.h"
}

lemming::dataplane::sai::PortModuleType convert_sai_port_module_type_t_to_proto(
    const sai_int32_t val);
sai_port_module_type_t convert_sai_port_module_type_t_to_sai(
    lemming::dataplane::sai::PortModuleType val);

lemming::dataplane::sai::QueuePfcContinuousDeadlockState
convert_sai_queue_pfc_continuous_deadlock_state_t_to_proto(
    const sai_int32_t val);
sai_queue_pfc_continuous_deadlock_state_t
convert_sai_queue_pfc_continuous_deadlock_state_t_to_sai(
    lemming::dataplane::sai::QueuePfcContinuousDeadlockState val);

lemming::dataplane::sai::PacketVlan convert_sai_packet_vlan_t_to_proto(
    const sai_int32_t val);
sai_packet_vlan_t convert_sai_packet_vlan_t_to_sai(
    lemming::dataplane::sai::PacketVlan val);

lemming::dataplane::sai::SwitchSwitchingMode
convert_sai_switch_switching_mode_t_to_proto(const sai_int32_t val);
sai_switch_switching_mode_t convert_sai_switch_switching_mode_t_to_sai(
    lemming::dataplane::sai::SwitchSwitchingMode val);

lemming::dataplane::sai::BridgePortTaggingMode
convert_sai_bridge_port_tagging_mode_t_to_proto(const sai_int32_t val);
sai_bridge_port_tagging_mode_t convert_sai_bridge_port_tagging_mode_t_to_sai(
    lemming::dataplane::sai::BridgePortTaggingMode val);

lemming::dataplane::sai::BridgeFloodControlType
convert_sai_bridge_flood_control_type_t_to_proto(const sai_int32_t val);
sai_bridge_flood_control_type_t convert_sai_bridge_flood_control_type_t_to_sai(
    lemming::dataplane::sai::BridgeFloodControlType val);

lemming::dataplane::sai::IpsecSaStat convert_sai_ipsec_sa_stat_t_to_proto(
    const sai_int32_t val);
sai_ipsec_sa_stat_t convert_sai_ipsec_sa_stat_t_to_sai(
    lemming::dataplane::sai::IpsecSaStat val);

lemming::dataplane::sai::StpPortState convert_sai_stp_port_state_t_to_proto(
    const sai_int32_t val);
sai_stp_port_state_t convert_sai_stp_port_state_t_to_sai(
    lemming::dataplane::sai::StpPortState val);

lemming::dataplane::sai::SwitchOperStatus
convert_sai_switch_oper_status_t_to_proto(const sai_int32_t val);
sai_switch_oper_status_t convert_sai_switch_oper_status_t_to_sai(
    lemming::dataplane::sai::SwitchOperStatus val);

lemming::dataplane::sai::TamTransportType
convert_sai_tam_transport_type_t_to_proto(const sai_int32_t val);
sai_tam_transport_type_t convert_sai_tam_transport_type_t_to_sai(
    lemming::dataplane::sai::TamTransportType val);

lemming::dataplane::sai::BufferPoolStat convert_sai_buffer_pool_stat_t_to_proto(
    const sai_int32_t val);
sai_buffer_pool_stat_t convert_sai_buffer_pool_stat_t_to_sai(
    lemming::dataplane::sai::BufferPoolStat val);

lemming::dataplane::sai::InsegEntryPopTtlMode
convert_sai_inseg_entry_pop_ttl_mode_t_to_proto(const sai_int32_t val);
sai_inseg_entry_pop_ttl_mode_t convert_sai_inseg_entry_pop_ttl_mode_t_to_sai(
    lemming::dataplane::sai::InsegEntryPopTtlMode val);

lemming::dataplane::sai::PortMdixModeConfig
convert_sai_port_mdix_mode_config_t_to_proto(const sai_int32_t val);
sai_port_mdix_mode_config_t convert_sai_port_mdix_mode_config_t_to_sai(
    lemming::dataplane::sai::PortMdixModeConfig val);

lemming::dataplane::sai::FdbEntryType convert_sai_fdb_entry_type_t_to_proto(
    const sai_int32_t val);
sai_fdb_entry_type_t convert_sai_fdb_entry_type_t_to_sai(
    lemming::dataplane::sai::FdbEntryType val);

lemming::dataplane::sai::MacsecFlowStat convert_sai_macsec_flow_stat_t_to_proto(
    const sai_int32_t val);
sai_macsec_flow_stat_t convert_sai_macsec_flow_stat_t_to_sai(
    lemming::dataplane::sai::MacsecFlowStat val);

lemming::dataplane::sai::SystemPortType convert_sai_system_port_type_t_to_proto(
    const sai_int32_t val);
sai_system_port_type_t convert_sai_system_port_type_t_to_sai(
    lemming::dataplane::sai::SystemPortType val);

lemming::dataplane::sai::IpsecDirection convert_sai_ipsec_direction_t_to_proto(
    const sai_int32_t val);
sai_ipsec_direction_t convert_sai_ipsec_direction_t_to_sai(
    lemming::dataplane::sai::IpsecDirection val);

lemming::dataplane::sai::IpsecPortStat convert_sai_ipsec_port_stat_t_to_proto(
    const sai_int32_t val);
sai_ipsec_port_stat_t convert_sai_ipsec_port_stat_t_to_sai(
    lemming::dataplane::sai::IpsecPortStat val);

lemming::dataplane::sai::MirrorSessionType
convert_sai_mirror_session_type_t_to_proto(const sai_int32_t val);
sai_mirror_session_type_t convert_sai_mirror_session_type_t_to_sai(
    lemming::dataplane::sai::MirrorSessionType val);

lemming::dataplane::sai::PortType convert_sai_port_type_t_to_proto(
    const sai_int32_t val);
sai_port_type_t convert_sai_port_type_t_to_sai(
    lemming::dataplane::sai::PortType val);

lemming::dataplane::sai::TunnelEncapEcnMode
convert_sai_tunnel_encap_ecn_mode_t_to_proto(const sai_int32_t val);
sai_tunnel_encap_ecn_mode_t convert_sai_tunnel_encap_ecn_mode_t_to_sai(
    lemming::dataplane::sai::TunnelEncapEcnMode val);

lemming::dataplane::sai::BfdSessionOffloadType
convert_sai_bfd_session_offload_type_t_to_proto(const sai_int32_t val);
sai_bfd_session_offload_type_t convert_sai_bfd_session_offload_type_t_to_sai(
    lemming::dataplane::sai::BfdSessionOffloadType val);

lemming::dataplane::sai::TableBitmapClassificationEntryStat
convert_sai_table_bitmap_classification_entry_stat_t_to_proto(
    const sai_int32_t val);
sai_table_bitmap_classification_entry_stat_t
convert_sai_table_bitmap_classification_entry_stat_t_to_sai(
    lemming::dataplane::sai::TableBitmapClassificationEntryStat val);

lemming::dataplane::sai::HostifTableEntryChannelType
convert_sai_hostif_table_entry_channel_type_t_to_proto(const sai_int32_t val);
sai_hostif_table_entry_channel_type_t
convert_sai_hostif_table_entry_channel_type_t_to_sai(
    lemming::dataplane::sai::HostifTableEntryChannelType val);

lemming::dataplane::sai::VlanTaggingMode
convert_sai_vlan_tagging_mode_t_to_proto(const sai_int32_t val);
sai_vlan_tagging_mode_t convert_sai_vlan_tagging_mode_t_to_sai(
    lemming::dataplane::sai::VlanTaggingMode val);

lemming::dataplane::sai::PortDualMedia convert_sai_port_dual_media_t_to_proto(
    const sai_int32_t val);
sai_port_dual_media_t convert_sai_port_dual_media_t_to_sai(
    lemming::dataplane::sai::PortDualMedia val);

lemming::dataplane::sai::MacsecCipherSuite
convert_sai_macsec_cipher_suite_t_to_proto(const sai_int32_t val);
sai_macsec_cipher_suite_t convert_sai_macsec_cipher_suite_t_to_sai(
    lemming::dataplane::sai::MacsecCipherSuite val);

lemming::dataplane::sai::NatEvent convert_sai_nat_event_t_to_proto(
    const sai_int32_t val);
sai_nat_event_t convert_sai_nat_event_t_to_sai(
    lemming::dataplane::sai::NatEvent val);

lemming::dataplane::sai::PortConnectorFailoverMode
convert_sai_port_connector_failover_mode_t_to_proto(const sai_int32_t val);
sai_port_connector_failover_mode_t
convert_sai_port_connector_failover_mode_t_to_sai(
    lemming::dataplane::sai::PortConnectorFailoverMode val);

lemming::dataplane::sai::HostifTableEntryType
convert_sai_hostif_table_entry_type_t_to_proto(const sai_int32_t val);
sai_hostif_table_entry_type_t convert_sai_hostif_table_entry_type_t_to_sai(
    lemming::dataplane::sai::HostifTableEntryType val);

lemming::dataplane::sai::PortMediaType convert_sai_port_media_type_t_to_proto(
    const sai_int32_t val);
sai_port_media_type_t convert_sai_port_media_type_t_to_sai(
    lemming::dataplane::sai::PortMediaType val);

lemming::dataplane::sai::PortBreakoutModeType
convert_sai_port_breakout_mode_type_t_to_proto(const sai_int32_t val);
sai_port_breakout_mode_type_t convert_sai_port_breakout_mode_type_t_to_sai(
    lemming::dataplane::sai::PortBreakoutModeType val);

lemming::dataplane::sai::ObjectStage convert_sai_object_stage_t_to_proto(
    const sai_int32_t val);
sai_object_stage_t convert_sai_object_stage_t_to_sai(
    lemming::dataplane::sai::ObjectStage val);

lemming::dataplane::sai::TableMetaTunnelEntryStat
convert_sai_table_meta_tunnel_entry_stat_t_to_proto(const sai_int32_t val);
sai_table_meta_tunnel_entry_stat_t
convert_sai_table_meta_tunnel_entry_stat_t_to_sai(
    lemming::dataplane::sai::TableMetaTunnelEntryStat val);

lemming::dataplane::sai::FdbFlushEntryType
convert_sai_fdb_flush_entry_type_t_to_proto(const sai_int32_t val);
sai_fdb_flush_entry_type_t convert_sai_fdb_flush_entry_type_t_to_sai(
    lemming::dataplane::sai::FdbFlushEntryType val);

lemming::dataplane::sai::HostifVlanTag convert_sai_hostif_vlan_tag_t_to_proto(
    const sai_int32_t val);
sai_hostif_vlan_tag_t convert_sai_hostif_vlan_tag_t_to_sai(
    lemming::dataplane::sai::HostifVlanTag val);

lemming::dataplane::sai::SwitchType convert_sai_switch_type_t_to_proto(
    const sai_int32_t val);
sai_switch_type_t convert_sai_switch_type_t_to_sai(
    lemming::dataplane::sai::SwitchType val);

lemming::dataplane::sai::TamBindPointType
convert_sai_tam_bind_point_type_t_to_proto(const sai_int32_t val);
sai_tam_bind_point_type_t convert_sai_tam_bind_point_type_t_to_sai(
    lemming::dataplane::sai::TamBindPointType val);

lemming::dataplane::sai::VlanFloodControlType
convert_sai_vlan_flood_control_type_t_to_proto(const sai_int32_t val);
sai_vlan_flood_control_type_t convert_sai_vlan_flood_control_type_t_to_sai(
    lemming::dataplane::sai::VlanFloodControlType val);

lemming::dataplane::sai::LogLevel convert_sai_log_level_t_to_proto(
    const sai_int32_t val);
sai_log_level_t convert_sai_log_level_t_to_sai(
    lemming::dataplane::sai::LogLevel val);

lemming::dataplane::sai::NextHopGroupMapType
convert_sai_next_hop_group_map_type_t_to_proto(const sai_int32_t val);
sai_next_hop_group_map_type_t convert_sai_next_hop_group_map_type_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMapType val);

lemming::dataplane::sai::QueueStat convert_sai_queue_stat_t_to_proto(
    const sai_int32_t val);
sai_queue_stat_t convert_sai_queue_stat_t_to_sai(
    lemming::dataplane::sai::QueueStat val);

lemming::dataplane::sai::PortInterfaceType
convert_sai_port_interface_type_t_to_proto(const sai_int32_t val);
sai_port_interface_type_t convert_sai_port_interface_type_t_to_sai(
    lemming::dataplane::sai::PortInterfaceType val);

lemming::dataplane::sai::SamplepacketType
convert_sai_samplepacket_type_t_to_proto(const sai_int32_t val);
sai_samplepacket_type_t convert_sai_samplepacket_type_t_to_sai(
    lemming::dataplane::sai::SamplepacketType val);

lemming::dataplane::sai::PacketAction convert_sai_packet_action_t_to_proto(
    const sai_int32_t val);
sai_packet_action_t convert_sai_packet_action_t_to_sai(
    lemming::dataplane::sai::PacketAction val);

lemming::dataplane::sai::TunnelPeerMode convert_sai_tunnel_peer_mode_t_to_proto(
    const sai_int32_t val);
sai_tunnel_peer_mode_t convert_sai_tunnel_peer_mode_t_to_sai(
    lemming::dataplane::sai::TunnelPeerMode val);

lemming::dataplane::sai::BfdSessionState
convert_sai_bfd_session_state_t_to_proto(const sai_int32_t val);
sai_bfd_session_state_t convert_sai_bfd_session_state_t_to_sai(
    lemming::dataplane::sai::BfdSessionState val);

lemming::dataplane::sai::OutDropReason convert_sai_out_drop_reason_t_to_proto(
    const sai_int32_t val);
sai_out_drop_reason_t convert_sai_out_drop_reason_t_to_sai(
    lemming::dataplane::sai::OutDropReason val);

lemming::dataplane::sai::NextHopGroupMemberConfiguredRole
convert_sai_next_hop_group_member_configured_role_t_to_proto(
    const sai_int32_t val);
sai_next_hop_group_member_configured_role_t
convert_sai_next_hop_group_member_configured_role_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMemberConfiguredRole val);

lemming::dataplane::sai::BulkOpErrorMode
convert_sai_bulk_op_error_mode_t_to_proto(const sai_int32_t val);
sai_bulk_op_error_mode_t convert_sai_bulk_op_error_mode_t_to_sai(
    lemming::dataplane::sai::BulkOpErrorMode val);

lemming::dataplane::sai::MacsecPortStat convert_sai_macsec_port_stat_t_to_proto(
    const sai_int32_t val);
sai_macsec_port_stat_t convert_sai_macsec_port_stat_t_to_sai(
    lemming::dataplane::sai::MacsecPortStat val);

lemming::dataplane::sai::MacsecSaStat convert_sai_macsec_sa_stat_t_to_proto(
    const sai_int32_t val);
sai_macsec_sa_stat_t convert_sai_macsec_sa_stat_t_to_sai(
    lemming::dataplane::sai::MacsecSaStat val);

lemming::dataplane::sai::Srv6SidlistType
convert_sai_srv6_sidlist_type_t_to_proto(const sai_int32_t val);
sai_srv6_sidlist_type_t convert_sai_srv6_sidlist_type_t_to_sai(
    lemming::dataplane::sai::Srv6SidlistType val);

lemming::dataplane::sai::BridgePortStat convert_sai_bridge_port_stat_t_to_proto(
    const sai_int32_t val);
sai_bridge_port_stat_t convert_sai_bridge_port_stat_t_to_sai(
    lemming::dataplane::sai::BridgePortStat val);

lemming::dataplane::sai::DebugCounterBindMethod
convert_sai_debug_counter_bind_method_t_to_proto(const sai_int32_t val);
sai_debug_counter_bind_method_t convert_sai_debug_counter_bind_method_t_to_sai(
    lemming::dataplane::sai::DebugCounterBindMethod val);

lemming::dataplane::sai::TunnelType convert_sai_tunnel_type_t_to_proto(
    const sai_int32_t val);
sai_tunnel_type_t convert_sai_tunnel_type_t_to_sai(
    lemming::dataplane::sai::TunnelType val);

lemming::dataplane::sai::PortFecMode convert_sai_port_fec_mode_t_to_proto(
    const sai_int32_t val);
sai_port_fec_mode_t convert_sai_port_fec_mode_t_to_sai(
    lemming::dataplane::sai::PortFecMode val);

lemming::dataplane::sai::PortPtpMode convert_sai_port_ptp_mode_t_to_proto(
    const sai_int32_t val);
sai_port_ptp_mode_t convert_sai_port_ptp_mode_t_to_sai(
    lemming::dataplane::sai::PortPtpMode val);

lemming::dataplane::sai::TunnelDscpMode convert_sai_tunnel_dscp_mode_t_to_proto(
    const sai_int32_t val);
sai_tunnel_dscp_mode_t convert_sai_tunnel_dscp_mode_t_to_sai(
    lemming::dataplane::sai::TunnelDscpMode val);

lemming::dataplane::sai::BfdEncapsulationType
convert_sai_bfd_encapsulation_type_t_to_proto(const sai_int32_t val);
sai_bfd_encapsulation_type_t convert_sai_bfd_encapsulation_type_t_to_sai(
    lemming::dataplane::sai::BfdEncapsulationType val);

lemming::dataplane::sai::BufferPoolType convert_sai_buffer_pool_type_t_to_proto(
    const sai_int32_t val);
sai_buffer_pool_type_t convert_sai_buffer_pool_type_t_to_sai(
    lemming::dataplane::sai::BufferPoolType val);

lemming::dataplane::sai::CounterStat convert_sai_counter_stat_t_to_proto(
    const sai_int32_t val);
sai_counter_stat_t convert_sai_counter_stat_t_to_sai(
    lemming::dataplane::sai::CounterStat val);

lemming::dataplane::sai::QosMapType convert_sai_qos_map_type_t_to_proto(
    const sai_int32_t val);
sai_qos_map_type_t convert_sai_qos_map_type_t_to_sai(
    lemming::dataplane::sai::QosMapType val);

lemming::dataplane::sai::SwitchMcastSnoopingCapability
convert_sai_switch_mcast_snooping_capability_t_to_proto(const sai_int32_t val);
sai_switch_mcast_snooping_capability_t
convert_sai_switch_mcast_snooping_capability_t_to_sai(
    lemming::dataplane::sai::SwitchMcastSnoopingCapability val);

lemming::dataplane::sai::SwitchHardwareAccessBus
convert_sai_switch_hardware_access_bus_t_to_proto(const sai_int32_t val);
sai_switch_hardware_access_bus_t
convert_sai_switch_hardware_access_bus_t_to_sai(
    lemming::dataplane::sai::SwitchHardwareAccessBus val);

lemming::dataplane::sai::SwitchFirmwareLoadType
convert_sai_switch_firmware_load_type_t_to_proto(const sai_int32_t val);
sai_switch_firmware_load_type_t convert_sai_switch_firmware_load_type_t_to_sai(
    lemming::dataplane::sai::SwitchFirmwareLoadType val);

lemming::dataplane::sai::TunnelTtlMode convert_sai_tunnel_ttl_mode_t_to_proto(
    const sai_int32_t val);
sai_tunnel_ttl_mode_t convert_sai_tunnel_ttl_mode_t_to_sai(
    lemming::dataplane::sai::TunnelTtlMode val);

lemming::dataplane::sai::DtelEventType convert_sai_dtel_event_type_t_to_proto(
    const sai_int32_t val);
sai_dtel_event_type_t convert_sai_dtel_event_type_t_to_sai(
    lemming::dataplane::sai::DtelEventType val);

lemming::dataplane::sai::IpsecCipher convert_sai_ipsec_cipher_t_to_proto(
    const sai_int32_t val);
sai_ipsec_cipher_t convert_sai_ipsec_cipher_t_to_sai(
    lemming::dataplane::sai::IpsecCipher val);

lemming::dataplane::sai::PortInternalLoopbackMode
convert_sai_port_internal_loopback_mode_t_to_proto(const sai_int32_t val);
sai_port_internal_loopback_mode_t
convert_sai_port_internal_loopback_mode_t_to_sai(
    lemming::dataplane::sai::PortInternalLoopbackMode val);

lemming::dataplane::sai::IpAddrFamily convert_sai_ip_addr_family_t_to_proto(
    const sai_int32_t val);
sai_ip_addr_family_t convert_sai_ip_addr_family_t_to_sai(
    lemming::dataplane::sai::IpAddrFamily val);

lemming::dataplane::sai::AclBindPointType
convert_sai_acl_bind_point_type_t_to_proto(const sai_int32_t val);
sai_acl_bind_point_type_t convert_sai_acl_bind_point_type_t_to_sai(
    lemming::dataplane::sai::AclBindPointType val);

lemming::dataplane::sai::OutsegTtlMode convert_sai_outseg_ttl_mode_t_to_proto(
    const sai_int32_t val);
sai_outseg_ttl_mode_t convert_sai_outseg_ttl_mode_t_to_sai(
    lemming::dataplane::sai::OutsegTtlMode val);

lemming::dataplane::sai::SwitchStat convert_sai_switch_stat_t_to_proto(
    const sai_int32_t val);
sai_switch_stat_t convert_sai_switch_stat_t_to_sai(
    lemming::dataplane::sai::SwitchStat val);

lemming::dataplane::sai::AclTableGroupType
convert_sai_acl_table_group_type_t_to_proto(const sai_int32_t val);
sai_acl_table_group_type_t convert_sai_acl_table_group_type_t_to_sai(
    lemming::dataplane::sai::AclTableGroupType val);

lemming::dataplane::sai::MeterType convert_sai_meter_type_t_to_proto(
    const sai_int32_t val);
sai_meter_type_t convert_sai_meter_type_t_to_sai(
    lemming::dataplane::sai::MeterType val);

lemming::dataplane::sai::PolicerColorSource
convert_sai_policer_color_source_t_to_proto(const sai_int32_t val);
sai_policer_color_source_t convert_sai_policer_color_source_t_to_sai(
    lemming::dataplane::sai::PolicerColorSource val);

lemming::dataplane::sai::TableMetaTunnelEntryAction
convert_sai_table_meta_tunnel_entry_action_t_to_proto(const sai_int32_t val);
sai_table_meta_tunnel_entry_action_t
convert_sai_table_meta_tunnel_entry_action_t_to_sai(
    lemming::dataplane::sai::TableMetaTunnelEntryAction val);

lemming::dataplane::sai::InsegEntryPopQosMode
convert_sai_inseg_entry_pop_qos_mode_t_to_proto(const sai_int32_t val);
sai_inseg_entry_pop_qos_mode_t convert_sai_inseg_entry_pop_qos_mode_t_to_sai(
    lemming::dataplane::sai::InsegEntryPopQosMode val);

lemming::dataplane::sai::RouterInterfaceType
convert_sai_router_interface_type_t_to_proto(const sai_int32_t val);
sai_router_interface_type_t convert_sai_router_interface_type_t_to_sai(
    lemming::dataplane::sai::RouterInterfaceType val);

lemming::dataplane::sai::TamTransportAuthType
convert_sai_tam_transport_auth_type_t_to_proto(const sai_int32_t val);
sai_tam_transport_auth_type_t convert_sai_tam_transport_auth_type_t_to_sai(
    lemming::dataplane::sai::TamTransportAuthType val);

lemming::dataplane::sai::AclIpType convert_sai_acl_ip_type_t_to_proto(
    const sai_int32_t val);
sai_acl_ip_type_t convert_sai_acl_ip_type_t_to_sai(
    lemming::dataplane::sai::AclIpType val);

lemming::dataplane::sai::DebugCounterType
convert_sai_debug_counter_type_t_to_proto(const sai_int32_t val);
sai_debug_counter_type_t convert_sai_debug_counter_type_t_to_sai(
    lemming::dataplane::sai::DebugCounterType val);

lemming::dataplane::sai::PortPrbsConfig convert_sai_port_prbs_config_t_to_proto(
    const sai_int32_t val);
sai_port_prbs_config_t convert_sai_port_prbs_config_t_to_sai(
    lemming::dataplane::sai::PortPrbsConfig val);

lemming::dataplane::sai::BufferPoolThresholdMode
convert_sai_buffer_pool_threshold_mode_t_to_proto(const sai_int32_t val);
sai_buffer_pool_threshold_mode_t
convert_sai_buffer_pool_threshold_mode_t_to_sai(
    lemming::dataplane::sai::BufferPoolThresholdMode val);

lemming::dataplane::sai::TableBitmapRouterEntryStat
convert_sai_table_bitmap_router_entry_stat_t_to_proto(const sai_int32_t val);
sai_table_bitmap_router_entry_stat_t
convert_sai_table_bitmap_router_entry_stat_t_to_sai(
    lemming::dataplane::sai::TableBitmapRouterEntryStat val);

lemming::dataplane::sai::NativeHashField
convert_sai_native_hash_field_t_to_proto(const sai_int32_t val);
sai_native_hash_field_t convert_sai_native_hash_field_t_to_sai(
    lemming::dataplane::sai::NativeHashField val);

lemming::dataplane::sai::MacsecDirection
convert_sai_macsec_direction_t_to_proto(const sai_int32_t val);
sai_macsec_direction_t convert_sai_macsec_direction_t_to_sai(
    lemming::dataplane::sai::MacsecDirection val);

lemming::dataplane::sai::MacsecScStat convert_sai_macsec_sc_stat_t_to_proto(
    const sai_int32_t val);
sai_macsec_sc_stat_t convert_sai_macsec_sc_stat_t_to_sai(
    lemming::dataplane::sai::MacsecScStat val);

lemming::dataplane::sai::Api convert_sai_api_t_to_proto(const sai_int32_t val);
sai_api_t convert_sai_api_t_to_sai(lemming::dataplane::sai::Api val);

lemming::dataplane::sai::AclRangeType convert_sai_acl_range_type_t_to_proto(
    const sai_int32_t val);
sai_acl_range_type_t convert_sai_acl_range_type_t_to_sai(
    lemming::dataplane::sai::AclRangeType val);

lemming::dataplane::sai::BridgePortType convert_sai_bridge_port_type_t_to_proto(
    const sai_int32_t val);
sai_bridge_port_type_t convert_sai_bridge_port_type_t_to_sai(
    lemming::dataplane::sai::BridgePortType val);

lemming::dataplane::sai::PortOperStatus convert_sai_port_oper_status_t_to_proto(
    const sai_int32_t val);
sai_port_oper_status_t convert_sai_port_oper_status_t_to_sai(
    lemming::dataplane::sai::PortOperStatus val);

lemming::dataplane::sai::PortLoopbackMode
convert_sai_port_loopback_mode_t_to_proto(const sai_int32_t val);
sai_port_loopback_mode_t convert_sai_port_loopback_mode_t_to_sai(
    lemming::dataplane::sai::PortLoopbackMode val);

lemming::dataplane::sai::PortLinkTrainingFailureStatus
convert_sai_port_link_training_failure_status_t_to_proto(const sai_int32_t val);
sai_port_link_training_failure_status_t
convert_sai_port_link_training_failure_status_t_to_sai(
    lemming::dataplane::sai::PortLinkTrainingFailureStatus val);

lemming::dataplane::sai::MySidEntryEndpointBehavior
convert_sai_my_sid_entry_endpoint_behavior_t_to_proto(const sai_int32_t val);
sai_my_sid_entry_endpoint_behavior_t
convert_sai_my_sid_entry_endpoint_behavior_t_to_sai(
    lemming::dataplane::sai::MySidEntryEndpointBehavior val);

lemming::dataplane::sai::UdfGroupType convert_sai_udf_group_type_t_to_proto(
    const sai_int32_t val);
sai_udf_group_type_t convert_sai_udf_group_type_t_to_sai(
    lemming::dataplane::sai::UdfGroupType val);

lemming::dataplane::sai::BfdSessionStat convert_sai_bfd_session_stat_t_to_proto(
    const sai_int32_t val);
sai_bfd_session_stat_t convert_sai_bfd_session_stat_t_to_sai(
    lemming::dataplane::sai::BfdSessionStat val);

lemming::dataplane::sai::HostifTxType convert_sai_hostif_tx_type_t_to_proto(
    const sai_int32_t val);
sai_hostif_tx_type_t convert_sai_hostif_tx_type_t_to_sai(
    lemming::dataplane::sai::HostifTxType val);

lemming::dataplane::sai::MirrorSessionCongestionMode
convert_sai_mirror_session_congestion_mode_t_to_proto(const sai_int32_t val);
sai_mirror_session_congestion_mode_t
convert_sai_mirror_session_congestion_mode_t_to_sai(
    lemming::dataplane::sai::MirrorSessionCongestionMode val);

lemming::dataplane::sai::TamReportType convert_sai_tam_report_type_t_to_proto(
    const sai_int32_t val);
sai_tam_report_type_t convert_sai_tam_report_type_t_to_sai(
    lemming::dataplane::sai::TamReportType val);

lemming::dataplane::sai::TamEventType convert_sai_tam_event_type_t_to_proto(
    const sai_int32_t val);
sai_tam_event_type_t convert_sai_tam_event_type_t_to_sai(
    lemming::dataplane::sai::TamEventType val);

lemming::dataplane::sai::UdfBase convert_sai_udf_base_t_to_proto(
    const sai_int32_t val);
sai_udf_base_t convert_sai_udf_base_t_to_sai(
    lemming::dataplane::sai::UdfBase val);

lemming::dataplane::sai::BfdSessionType convert_sai_bfd_session_type_t_to_proto(
    const sai_int32_t val);
sai_bfd_session_type_t convert_sai_bfd_session_type_t_to_sai(
    lemming::dataplane::sai::BfdSessionType val);

lemming::dataplane::sai::BridgePortFdbLearningMode
convert_sai_bridge_port_fdb_learning_mode_t_to_proto(const sai_int32_t val);
sai_bridge_port_fdb_learning_mode_t
convert_sai_bridge_port_fdb_learning_mode_t_to_sai(
    lemming::dataplane::sai::BridgePortFdbLearningMode val);

lemming::dataplane::sai::PortStat convert_sai_port_stat_t_to_proto(
    const sai_int32_t val);
sai_port_stat_t convert_sai_port_stat_t_to_sai(
    lemming::dataplane::sai::PortStat val);

lemming::dataplane::sai::TamTelemetryType
convert_sai_tam_telemetry_type_t_to_proto(const sai_int32_t val);
sai_tam_telemetry_type_t convert_sai_tam_telemetry_type_t_to_sai(
    lemming::dataplane::sai::TamTelemetryType val);

lemming::dataplane::sai::OutsegExpMode convert_sai_outseg_exp_mode_t_to_proto(
    const sai_int32_t val);
sai_outseg_exp_mode_t convert_sai_outseg_exp_mode_t_to_sai(
    lemming::dataplane::sai::OutsegExpMode val);

lemming::dataplane::sai::TamEventThresholdUnit
convert_sai_tam_event_threshold_unit_t_to_proto(const sai_int32_t val);
sai_tam_event_threshold_unit_t convert_sai_tam_event_threshold_unit_t_to_sai(
    lemming::dataplane::sai::TamEventThresholdUnit val);

lemming::dataplane::sai::AclDtelFlowOp convert_sai_acl_dtel_flow_op_t_to_proto(
    const sai_int32_t val);
sai_acl_dtel_flow_op_t convert_sai_acl_dtel_flow_op_t_to_sai(
    lemming::dataplane::sai::AclDtelFlowOp val);

lemming::dataplane::sai::L2mcEntryType convert_sai_l2mc_entry_type_t_to_proto(
    const sai_int32_t val);
sai_l2mc_entry_type_t convert_sai_l2mc_entry_type_t_to_sai(
    lemming::dataplane::sai::L2mcEntryType val);

lemming::dataplane::sai::PortAutoNegConfigMode
convert_sai_port_auto_neg_config_mode_t_to_proto(const sai_int32_t val);
sai_port_auto_neg_config_mode_t convert_sai_port_auto_neg_config_mode_t_to_sai(
    lemming::dataplane::sai::PortAutoNegConfigMode val);

lemming::dataplane::sai::SamplepacketMode
convert_sai_samplepacket_mode_t_to_proto(const sai_int32_t val);
sai_samplepacket_mode_t convert_sai_samplepacket_mode_t_to_sai(
    lemming::dataplane::sai::SamplepacketMode val);

lemming::dataplane::sai::MySidEntryEndpointBehaviorFlavor
convert_sai_my_sid_entry_endpoint_behavior_flavor_t_to_proto(
    const sai_int32_t val);
sai_my_sid_entry_endpoint_behavior_flavor_t
convert_sai_my_sid_entry_endpoint_behavior_flavor_t_to_sai(
    lemming::dataplane::sai::MySidEntryEndpointBehaviorFlavor val);

lemming::dataplane::sai::TamReportMode convert_sai_tam_report_mode_t_to_proto(
    const sai_int32_t val);
sai_tam_report_mode_t convert_sai_tam_report_mode_t_to_sai(
    lemming::dataplane::sai::TamReportMode val);

lemming::dataplane::sai::TunnelTermTableEntryType
convert_sai_tunnel_term_table_entry_type_t_to_proto(const sai_int32_t val);
sai_tunnel_term_table_entry_type_t
convert_sai_tunnel_term_table_entry_type_t_to_sai(
    lemming::dataplane::sai::TunnelTermTableEntryType val);

lemming::dataplane::sai::OutsegType convert_sai_outseg_type_t_to_proto(
    const sai_int32_t val);
sai_outseg_type_t convert_sai_outseg_type_t_to_sai(
    lemming::dataplane::sai::OutsegType val);

lemming::dataplane::sai::IngressPriorityGroupStat
convert_sai_ingress_priority_group_stat_t_to_proto(const sai_int32_t val);
sai_ingress_priority_group_stat_t
convert_sai_ingress_priority_group_stat_t_to_sai(
    lemming::dataplane::sai::IngressPriorityGroupStat val);

lemming::dataplane::sai::NextHopGroupType
convert_sai_next_hop_group_type_t_to_proto(const sai_int32_t val);
sai_next_hop_group_type_t convert_sai_next_hop_group_type_t_to_sai(
    lemming::dataplane::sai::NextHopGroupType val);

lemming::dataplane::sai::PortPriorityFlowControlMode
convert_sai_port_priority_flow_control_mode_t_to_proto(const sai_int32_t val);
sai_port_priority_flow_control_mode_t
convert_sai_port_priority_flow_control_mode_t_to_sai(
    lemming::dataplane::sai::PortPriorityFlowControlMode val);

lemming::dataplane::sai::PortLinkTrainingRxStatus
convert_sai_port_link_training_rx_status_t_to_proto(const sai_int32_t val);
sai_port_link_training_rx_status_t
convert_sai_port_link_training_rx_status_t_to_sai(
    lemming::dataplane::sai::PortLinkTrainingRxStatus val);

lemming::dataplane::sai::TamIntPresenceType
convert_sai_tam_int_presence_type_t_to_proto(const sai_int32_t val);
sai_tam_int_presence_type_t convert_sai_tam_int_presence_type_t_to_sai(
    lemming::dataplane::sai::TamIntPresenceType val);

lemming::dataplane::sai::TunnelMapType convert_sai_tunnel_map_type_t_to_proto(
    const sai_int32_t val);
sai_tunnel_map_type_t convert_sai_tunnel_map_type_t_to_sai(
    lemming::dataplane::sai::TunnelMapType val);

lemming::dataplane::sai::CommonApi convert_sai_common_api_t_to_proto(
    const sai_int32_t val);
sai_common_api_t convert_sai_common_api_t_to_sai(
    lemming::dataplane::sai::CommonApi val);

lemming::dataplane::sai::PacketColor convert_sai_packet_color_t_to_proto(
    const sai_int32_t val);
sai_packet_color_t convert_sai_packet_color_t_to_sai(
    lemming::dataplane::sai::PacketColor val);

lemming::dataplane::sai::FdbEvent convert_sai_fdb_event_t_to_proto(
    const sai_int32_t val);
sai_fdb_event_t convert_sai_fdb_event_t_to_sai(
    lemming::dataplane::sai::FdbEvent val);

lemming::dataplane::sai::IpsecSaOctetCountStatus
convert_sai_ipsec_sa_octet_count_status_t_to_proto(const sai_int32_t val);
sai_ipsec_sa_octet_count_status_t
convert_sai_ipsec_sa_octet_count_status_t_to_sai(
    lemming::dataplane::sai::IpsecSaOctetCountStatus val);

lemming::dataplane::sai::PolicerMode convert_sai_policer_mode_t_to_proto(
    const sai_int32_t val);
sai_policer_mode_t convert_sai_policer_mode_t_to_sai(
    lemming::dataplane::sai::PolicerMode val);

lemming::dataplane::sai::ObjectTypeExtensions
convert_sai_object_type_extensions_t_to_proto(const sai_int32_t val);
sai_object_type_extensions_t convert_sai_object_type_extensions_t_to_sai(
    lemming::dataplane::sai::ObjectTypeExtensions val);

lemming::dataplane::sai::QueueType convert_sai_queue_type_t_to_proto(
    const sai_int32_t val);
sai_queue_type_t convert_sai_queue_type_t_to_sai(
    lemming::dataplane::sai::QueueType val);

lemming::dataplane::sai::SwitchAttrExtensions
convert_sai_switch_attr_extensions_t_to_proto(const sai_int32_t val);
sai_switch_attr_extensions_t convert_sai_switch_attr_extensions_t_to_sai(
    lemming::dataplane::sai::SwitchAttrExtensions val);

lemming::dataplane::sai::TamIntType convert_sai_tam_int_type_t_to_proto(
    const sai_int32_t val);
sai_tam_int_type_t convert_sai_tam_int_type_t_to_sai(
    lemming::dataplane::sai::TamIntType val);

lemming::dataplane::sai::TableBitmapClassificationEntryAction
convert_sai_table_bitmap_classification_entry_action_t_to_proto(
    const sai_int32_t val);
sai_table_bitmap_classification_entry_action_t
convert_sai_table_bitmap_classification_entry_action_t_to_sai(
    lemming::dataplane::sai::TableBitmapClassificationEntryAction val);

lemming::dataplane::sai::HostifUserDefinedTrapType
convert_sai_hostif_user_defined_trap_type_t_to_proto(const sai_int32_t val);
sai_hostif_user_defined_trap_type_t
convert_sai_hostif_user_defined_trap_type_t_to_sai(
    lemming::dataplane::sai::HostifUserDefinedTrapType val);

lemming::dataplane::sai::PortMdixModeStatus
convert_sai_port_mdix_mode_status_t_to_proto(const sai_int32_t val);
sai_port_mdix_mode_status_t convert_sai_port_mdix_mode_status_t_to_sai(
    lemming::dataplane::sai::PortMdixModeStatus val);

lemming::dataplane::sai::TableBitmapRouterEntryAction
convert_sai_table_bitmap_router_entry_action_t_to_proto(const sai_int32_t val);
sai_table_bitmap_router_entry_action_t
convert_sai_table_bitmap_router_entry_action_t_to_sai(
    lemming::dataplane::sai::TableBitmapRouterEntryAction val);

lemming::dataplane::sai::MacsecMaxSecureAssociationsPerSc
convert_sai_macsec_max_secure_associations_per_sc_t_to_proto(
    const sai_int32_t val);
sai_macsec_max_secure_associations_per_sc_t
convert_sai_macsec_max_secure_associations_per_sc_t_to_sai(
    lemming::dataplane::sai::MacsecMaxSecureAssociationsPerSc val);

lemming::dataplane::sai::InsegEntryPscType
convert_sai_inseg_entry_psc_type_t_to_proto(const sai_int32_t val);
sai_inseg_entry_psc_type_t convert_sai_inseg_entry_psc_type_t_to_sai(
    lemming::dataplane::sai::InsegEntryPscType val);

lemming::dataplane::sai::NextHopType convert_sai_next_hop_type_t_to_proto(
    const sai_int32_t val);
sai_next_hop_type_t convert_sai_next_hop_type_t_to_sai(
    lemming::dataplane::sai::NextHopType val);

lemming::dataplane::sai::NextHopGroupMemberObservedRole
convert_sai_next_hop_group_member_observed_role_t_to_proto(
    const sai_int32_t val);
sai_next_hop_group_member_observed_role_t
convert_sai_next_hop_group_member_observed_role_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMemberObservedRole val);

lemming::dataplane::sai::BridgeStat convert_sai_bridge_stat_t_to_proto(
    const sai_int32_t val);
sai_bridge_stat_t convert_sai_bridge_stat_t_to_sai(
    lemming::dataplane::sai::BridgeStat val);

lemming::dataplane::sai::BufferProfileThresholdMode
convert_sai_buffer_profile_threshold_mode_t_to_proto(const sai_int32_t val);
sai_buffer_profile_threshold_mode_t
convert_sai_buffer_profile_threshold_mode_t_to_sai(
    lemming::dataplane::sai::BufferProfileThresholdMode val);

lemming::dataplane::sai::CounterType convert_sai_counter_type_t_to_proto(
    const sai_int32_t val);
sai_counter_type_t convert_sai_counter_type_t_to_sai(
    lemming::dataplane::sai::CounterType val);

lemming::dataplane::sai::AclStage convert_sai_acl_stage_t_to_proto(
    const sai_int32_t val);
sai_acl_stage_t convert_sai_acl_stage_t_to_sai(
    lemming::dataplane::sai::AclStage val);

lemming::dataplane::sai::SwitchRestartType
convert_sai_switch_restart_type_t_to_proto(const sai_int32_t val);
sai_switch_restart_type_t convert_sai_switch_restart_type_t_to_sai(
    lemming::dataplane::sai::SwitchRestartType val);

lemming::dataplane::sai::SwitchFailoverConfigMode
convert_sai_switch_failover_config_mode_t_to_proto(const sai_int32_t val);
sai_switch_failover_config_mode_t
convert_sai_switch_failover_config_mode_t_to_sai(
    lemming::dataplane::sai::SwitchFailoverConfigMode val);

lemming::dataplane::sai::TamTelMathFuncType
convert_sai_tam_tel_math_func_type_t_to_proto(const sai_int32_t val);
sai_tam_tel_math_func_type_t convert_sai_tam_tel_math_func_type_t_to_sai(
    lemming::dataplane::sai::TamTelMathFuncType val);

lemming::dataplane::sai::TamReportingUnit
convert_sai_tam_reporting_unit_t_to_proto(const sai_int32_t val);
sai_tam_reporting_unit_t convert_sai_tam_reporting_unit_t_to_sai(
    lemming::dataplane::sai::TamReportingUnit val);

lemming::dataplane::sai::PortPrbsRxStatus
convert_sai_port_prbs_rx_status_t_to_proto(const sai_int32_t val);
sai_port_prbs_rx_status_t convert_sai_port_prbs_rx_status_t_to_sai(
    lemming::dataplane::sai::PortPrbsRxStatus val);

lemming::dataplane::sai::TlvType convert_sai_tlv_type_t_to_proto(
    const sai_int32_t val);
sai_tlv_type_t convert_sai_tlv_type_t_to_sai(
    lemming::dataplane::sai::TlvType val);

lemming::dataplane::sai::StatsMode convert_sai_stats_mode_t_to_proto(
    const sai_int32_t val);
sai_stats_mode_t convert_sai_stats_mode_t_to_sai(
    lemming::dataplane::sai::StatsMode val);

lemming::dataplane::sai::PortPoolStat convert_sai_port_pool_stat_t_to_proto(
    const sai_int32_t val);
sai_port_pool_stat_t convert_sai_port_pool_stat_t_to_sai(
    lemming::dataplane::sai::PortPoolStat val);

lemming::dataplane::sai::QueuePfcDeadlockEventType
convert_sai_queue_pfc_deadlock_event_type_t_to_proto(const sai_int32_t val);
sai_queue_pfc_deadlock_event_type_t
convert_sai_queue_pfc_deadlock_event_type_t_to_sai(
    lemming::dataplane::sai::QueuePfcDeadlockEventType val);

lemming::dataplane::sai::TunnelDecapEcnMode
convert_sai_tunnel_decap_ecn_mode_t_to_proto(const sai_int32_t val);
sai_tunnel_decap_ecn_mode_t convert_sai_tunnel_decap_ecn_mode_t_to_sai(
    lemming::dataplane::sai::TunnelDecapEcnMode val);

lemming::dataplane::sai::EcnMarkMode convert_sai_ecn_mark_mode_t_to_proto(
    const sai_int32_t val);
sai_ecn_mark_mode_t convert_sai_ecn_mark_mode_t_to_sai(
    lemming::dataplane::sai::EcnMarkMode val);

lemming::dataplane::sai::SwitchFirmwareLoadMethod
convert_sai_switch_firmware_load_method_t_to_proto(const sai_int32_t val);
sai_switch_firmware_load_method_t
convert_sai_switch_firmware_load_method_t_to_sai(
    lemming::dataplane::sai::SwitchFirmwareLoadMethod val);

lemming::dataplane::sai::TunnelVxlanUdpSportMode
convert_sai_tunnel_vxlan_udp_sport_mode_t_to_proto(const sai_int32_t val);
sai_tunnel_vxlan_udp_sport_mode_t
convert_sai_tunnel_vxlan_udp_sport_mode_t_to_sai(
    lemming::dataplane::sai::TunnelVxlanUdpSportMode val);

lemming::dataplane::sai::PortErrStatus convert_sai_port_err_status_t_to_proto(
    const sai_int32_t val);
sai_port_err_status_t convert_sai_port_err_status_t_to_sai(
    lemming::dataplane::sai::PortErrStatus val);

lemming::dataplane::sai::InDropReason convert_sai_in_drop_reason_t_to_proto(
    const sai_int32_t val);
sai_in_drop_reason_t convert_sai_in_drop_reason_t_to_sai(
    lemming::dataplane::sai::InDropReason val);

lemming::dataplane::sai::ApiExtensions convert_sai_api_extensions_t_to_proto(
    const sai_int32_t val);
sai_api_extensions_t convert_sai_api_extensions_t_to_sai(
    lemming::dataplane::sai::ApiExtensions val);

lemming::dataplane::sai::ObjectType convert_sai_object_type_t_to_proto(
    const sai_int32_t val);
sai_object_type_t convert_sai_object_type_t_to_sai(
    lemming::dataplane::sai::ObjectType val);

lemming::dataplane::sai::IpmcEntryType convert_sai_ipmc_entry_type_t_to_proto(
    const sai_int32_t val);
sai_ipmc_entry_type_t convert_sai_ipmc_entry_type_t_to_sai(
    lemming::dataplane::sai::IpmcEntryType val);

lemming::dataplane::sai::IsolationGroupType
convert_sai_isolation_group_type_t_to_proto(const sai_int32_t val);
sai_isolation_group_type_t convert_sai_isolation_group_type_t_to_sai(
    lemming::dataplane::sai::IsolationGroupType val);

lemming::dataplane::sai::SchedulingType convert_sai_scheduling_type_t_to_proto(
    const sai_int32_t val);
sai_scheduling_type_t convert_sai_scheduling_type_t_to_sai(
    lemming::dataplane::sai::SchedulingType val);

lemming::dataplane::sai::PortFecModeExtended
convert_sai_port_fec_mode_extended_t_to_proto(const sai_int32_t val);
sai_port_fec_mode_extended_t convert_sai_port_fec_mode_extended_t_to_sai(
    lemming::dataplane::sai::PortFecModeExtended val);

lemming::dataplane::sai::RouterInterfaceStat
convert_sai_router_interface_stat_t_to_proto(const sai_int32_t val);
sai_router_interface_stat_t convert_sai_router_interface_stat_t_to_sai(
    lemming::dataplane::sai::RouterInterfaceStat val);

lemming::dataplane::sai::HashAlgorithm convert_sai_hash_algorithm_t_to_proto(
    const sai_int32_t val);
sai_hash_algorithm_t convert_sai_hash_algorithm_t_to_sai(
    lemming::dataplane::sai::HashAlgorithm val);

lemming::dataplane::sai::AclIpFrag convert_sai_acl_ip_frag_t_to_proto(
    const sai_int32_t val);
sai_acl_ip_frag_t convert_sai_acl_ip_frag_t_to_sai(
    lemming::dataplane::sai::AclIpFrag val);

lemming::dataplane::sai::AclActionType convert_sai_acl_action_type_t_to_proto(
    const sai_int32_t val);
sai_acl_action_type_t convert_sai_acl_action_type_t_to_sai(
    lemming::dataplane::sai::AclActionType val);

lemming::dataplane::sai::HostifType convert_sai_hostif_type_t_to_proto(
    const sai_int32_t val);
sai_hostif_type_t convert_sai_hostif_type_t_to_sai(
    lemming::dataplane::sai::HostifType val);

lemming::dataplane::sai::NatType convert_sai_nat_type_t_to_proto(
    const sai_int32_t val);
sai_nat_type_t convert_sai_nat_type_t_to_sai(
    lemming::dataplane::sai::NatType val);

lemming::dataplane::sai::PolicerStat convert_sai_policer_stat_t_to_proto(
    const sai_int32_t val);
sai_policer_stat_t convert_sai_policer_stat_t_to_sai(
    lemming::dataplane::sai::PolicerStat val);

lemming::dataplane::sai::PortFlowControlMode
convert_sai_port_flow_control_mode_t_to_proto(const sai_int32_t val);
sai_port_flow_control_mode_t convert_sai_port_flow_control_mode_t_to_sai(
    lemming::dataplane::sai::PortFlowControlMode val);

lemming::dataplane::sai::TunnelStat convert_sai_tunnel_stat_t_to_proto(
    const sai_int32_t val);
sai_tunnel_stat_t convert_sai_tunnel_stat_t_to_sai(
    lemming::dataplane::sai::TunnelStat val);

lemming::dataplane::sai::VlanMcastLookupKeyType
convert_sai_vlan_mcast_lookup_key_type_t_to_proto(const sai_int32_t val);
sai_vlan_mcast_lookup_key_type_t
convert_sai_vlan_mcast_lookup_key_type_t_to_sai(
    lemming::dataplane::sai::VlanMcastLookupKeyType val);

lemming::dataplane::sai::BridgeType convert_sai_bridge_type_t_to_proto(
    const sai_int32_t val);
sai_bridge_type_t convert_sai_bridge_type_t_to_sai(
    lemming::dataplane::sai::BridgeType val);

lemming::dataplane::sai::HostifTrapType convert_sai_hostif_trap_type_t_to_proto(
    const sai_int32_t val);
sai_hostif_trap_type_t convert_sai_hostif_trap_type_t_to_sai(
    lemming::dataplane::sai::HostifTrapType val);

lemming::dataplane::sai::ErspanEncapsulationType
convert_sai_erspan_encapsulation_type_t_to_proto(const sai_int32_t val);
sai_erspan_encapsulation_type_t convert_sai_erspan_encapsulation_type_t_to_sai(
    lemming::dataplane::sai::ErspanEncapsulationType val);

lemming::dataplane::sai::VlanStat convert_sai_vlan_stat_t_to_proto(
    const sai_int32_t val);
sai_vlan_stat_t convert_sai_vlan_stat_t_to_sai(
    lemming::dataplane::sai::VlanStat val);

lemming::dataplane::sai::GenericProgrammableAttr
convert_sai_generic_programmable_attr_t_to_proto(const sai_int32_t val);
sai_generic_programmable_attr_t convert_sai_generic_programmable_attr_t_to_sai(
    lemming::dataplane::sai::GenericProgrammableAttr val);

lemming::dataplane::sai::TamEventAttr convert_sai_tam_event_attr_t_to_proto(
    const sai_int32_t val);
sai_tam_event_attr_t convert_sai_tam_event_attr_t_to_sai(
    lemming::dataplane::sai::TamEventAttr val);

lemming::dataplane::sai::BufferPoolAttr convert_sai_buffer_pool_attr_t_to_proto(
    const sai_int32_t val);
sai_buffer_pool_attr_t convert_sai_buffer_pool_attr_t_to_sai(
    lemming::dataplane::sai::BufferPoolAttr val);

lemming::dataplane::sai::IpmcEntryAttr convert_sai_ipmc_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_ipmc_entry_attr_t convert_sai_ipmc_entry_attr_t_to_sai(
    lemming::dataplane::sai::IpmcEntryAttr val);

lemming::dataplane::sai::TamCollectorAttr
convert_sai_tam_collector_attr_t_to_proto(const sai_int32_t val);
sai_tam_collector_attr_t convert_sai_tam_collector_attr_t_to_sai(
    lemming::dataplane::sai::TamCollectorAttr val);

lemming::dataplane::sai::BridgeAttr convert_sai_bridge_attr_t_to_proto(
    const sai_int32_t val);
sai_bridge_attr_t convert_sai_bridge_attr_t_to_sai(
    lemming::dataplane::sai::BridgeAttr val);

lemming::dataplane::sai::MacsecFlowAttr convert_sai_macsec_flow_attr_t_to_proto(
    const sai_int32_t val);
sai_macsec_flow_attr_t convert_sai_macsec_flow_attr_t_to_sai(
    lemming::dataplane::sai::MacsecFlowAttr val);

lemming::dataplane::sai::RouteEntryAttr convert_sai_route_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_route_entry_attr_t convert_sai_route_entry_attr_t_to_sai(
    lemming::dataplane::sai::RouteEntryAttr val);

lemming::dataplane::sai::UdfAttr convert_sai_udf_attr_t_to_proto(
    const sai_int32_t val);
sai_udf_attr_t convert_sai_udf_attr_t_to_sai(
    lemming::dataplane::sai::UdfAttr val);

lemming::dataplane::sai::TableBitmapClassificationEntryAttr
convert_sai_table_bitmap_classification_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_table_bitmap_classification_entry_attr_t
convert_sai_table_bitmap_classification_entry_attr_t_to_sai(
    lemming::dataplane::sai::TableBitmapClassificationEntryAttr val);

lemming::dataplane::sai::HostifAttr convert_sai_hostif_attr_t_to_proto(
    const sai_int32_t val);
sai_hostif_attr_t convert_sai_hostif_attr_t_to_sai(
    lemming::dataplane::sai::HostifAttr val);

lemming::dataplane::sai::InsegEntryAttr convert_sai_inseg_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_inseg_entry_attr_t convert_sai_inseg_entry_attr_t_to_sai(
    lemming::dataplane::sai::InsegEntryAttr val);

lemming::dataplane::sai::RouterInterfaceAttr
convert_sai_router_interface_attr_t_to_proto(const sai_int32_t val);
sai_router_interface_attr_t convert_sai_router_interface_attr_t_to_sai(
    lemming::dataplane::sai::RouterInterfaceAttr val);

lemming::dataplane::sai::HostifTrapAttr convert_sai_hostif_trap_attr_t_to_proto(
    const sai_int32_t val);
sai_hostif_trap_attr_t convert_sai_hostif_trap_attr_t_to_sai(
    lemming::dataplane::sai::HostifTrapAttr val);

lemming::dataplane::sai::IpsecSaAttr convert_sai_ipsec_sa_attr_t_to_proto(
    const sai_int32_t val);
sai_ipsec_sa_attr_t convert_sai_ipsec_sa_attr_t_to_sai(
    lemming::dataplane::sai::IpsecSaAttr val);

lemming::dataplane::sai::SystemPortAttr convert_sai_system_port_attr_t_to_proto(
    const sai_int32_t val);
sai_system_port_attr_t convert_sai_system_port_attr_t_to_sai(
    lemming::dataplane::sai::SystemPortAttr val);

lemming::dataplane::sai::HostifTableEntryAttr
convert_sai_hostif_table_entry_attr_t_to_proto(const sai_int32_t val);
sai_hostif_table_entry_attr_t convert_sai_hostif_table_entry_attr_t_to_sai(
    lemming::dataplane::sai::HostifTableEntryAttr val);

lemming::dataplane::sai::PortAttr convert_sai_port_attr_t_to_proto(
    const sai_int32_t val);
sai_port_attr_t convert_sai_port_attr_t_to_sai(
    lemming::dataplane::sai::PortAttr val);

lemming::dataplane::sai::DtelAttr convert_sai_dtel_attr_t_to_proto(
    const sai_int32_t val);
sai_dtel_attr_t convert_sai_dtel_attr_t_to_sai(
    lemming::dataplane::sai::DtelAttr val);

lemming::dataplane::sai::IsolationGroupAttr
convert_sai_isolation_group_attr_t_to_proto(const sai_int32_t val);
sai_isolation_group_attr_t convert_sai_isolation_group_attr_t_to_sai(
    lemming::dataplane::sai::IsolationGroupAttr val);

lemming::dataplane::sai::MyMacAttr convert_sai_my_mac_attr_t_to_proto(
    const sai_int32_t val);
sai_my_mac_attr_t convert_sai_my_mac_attr_t_to_sai(
    lemming::dataplane::sai::MyMacAttr val);

lemming::dataplane::sai::NextHopAttr convert_sai_next_hop_attr_t_to_proto(
    const sai_int32_t val);
sai_next_hop_attr_t convert_sai_next_hop_attr_t_to_sai(
    lemming::dataplane::sai::NextHopAttr val);

lemming::dataplane::sai::TunnelAttr convert_sai_tunnel_attr_t_to_proto(
    const sai_int32_t val);
sai_tunnel_attr_t convert_sai_tunnel_attr_t_to_sai(
    lemming::dataplane::sai::TunnelAttr val);

lemming::dataplane::sai::AclEntryAttr convert_sai_acl_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_acl_entry_attr_t convert_sai_acl_entry_attr_t_to_sai(
    lemming::dataplane::sai::AclEntryAttr val);

lemming::dataplane::sai::DebugCounterAttr
convert_sai_debug_counter_attr_t_to_proto(const sai_int32_t val);
sai_debug_counter_attr_t convert_sai_debug_counter_attr_t_to_sai(
    lemming::dataplane::sai::DebugCounterAttr val);

lemming::dataplane::sai::DtelIntSessionAttr
convert_sai_dtel_int_session_attr_t_to_proto(const sai_int32_t val);
sai_dtel_int_session_attr_t convert_sai_dtel_int_session_attr_t_to_sai(
    lemming::dataplane::sai::DtelIntSessionAttr val);

lemming::dataplane::sai::FdbEntryAttr convert_sai_fdb_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_fdb_entry_attr_t convert_sai_fdb_entry_attr_t_to_sai(
    lemming::dataplane::sai::FdbEntryAttr val);

lemming::dataplane::sai::L2mcGroupAttr convert_sai_l2mc_group_attr_t_to_proto(
    const sai_int32_t val);
sai_l2mc_group_attr_t convert_sai_l2mc_group_attr_t_to_sai(
    lemming::dataplane::sai::L2mcGroupAttr val);

lemming::dataplane::sai::TamMathFuncAttr
convert_sai_tam_math_func_attr_t_to_proto(const sai_int32_t val);
sai_tam_math_func_attr_t convert_sai_tam_math_func_attr_t_to_sai(
    lemming::dataplane::sai::TamMathFuncAttr val);

lemming::dataplane::sai::TunnelTermTableEntryAttr
convert_sai_tunnel_term_table_entry_attr_t_to_proto(const sai_int32_t val);
sai_tunnel_term_table_entry_attr_t
convert_sai_tunnel_term_table_entry_attr_t_to_sai(
    lemming::dataplane::sai::TunnelTermTableEntryAttr val);

lemming::dataplane::sai::UdfMatchAttr convert_sai_udf_match_attr_t_to_proto(
    const sai_int32_t val);
sai_udf_match_attr_t convert_sai_udf_match_attr_t_to_sai(
    lemming::dataplane::sai::UdfMatchAttr val);

lemming::dataplane::sai::NextHopGroupMapAttr
convert_sai_next_hop_group_map_attr_t_to_proto(const sai_int32_t val);
sai_next_hop_group_map_attr_t convert_sai_next_hop_group_map_attr_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMapAttr val);

lemming::dataplane::sai::SchedulerAttr convert_sai_scheduler_attr_t_to_proto(
    const sai_int32_t val);
sai_scheduler_attr_t convert_sai_scheduler_attr_t_to_sai(
    lemming::dataplane::sai::SchedulerAttr val);

lemming::dataplane::sai::TunnelMapEntryAttr
convert_sai_tunnel_map_entry_attr_t_to_proto(const sai_int32_t val);
sai_tunnel_map_entry_attr_t convert_sai_tunnel_map_entry_attr_t_to_sai(
    lemming::dataplane::sai::TunnelMapEntryAttr val);

lemming::dataplane::sai::VirtualRouterAttr
convert_sai_virtual_router_attr_t_to_proto(const sai_int32_t val);
sai_virtual_router_attr_t convert_sai_virtual_router_attr_t_to_sai(
    lemming::dataplane::sai::VirtualRouterAttr val);

lemming::dataplane::sai::IsolationGroupMemberAttr
convert_sai_isolation_group_member_attr_t_to_proto(const sai_int32_t val);
sai_isolation_group_member_attr_t
convert_sai_isolation_group_member_attr_t_to_sai(
    lemming::dataplane::sai::IsolationGroupMemberAttr val);

lemming::dataplane::sai::MirrorSessionAttr
convert_sai_mirror_session_attr_t_to_proto(const sai_int32_t val);
sai_mirror_session_attr_t convert_sai_mirror_session_attr_t_to_sai(
    lemming::dataplane::sai::MirrorSessionAttr val);

lemming::dataplane::sai::Srv6SidlistAttr
convert_sai_srv6_sidlist_attr_t_to_proto(const sai_int32_t val);
sai_srv6_sidlist_attr_t convert_sai_srv6_sidlist_attr_t_to_sai(
    lemming::dataplane::sai::Srv6SidlistAttr val);

lemming::dataplane::sai::TamIntAttr convert_sai_tam_int_attr_t_to_proto(
    const sai_int32_t val);
sai_tam_int_attr_t convert_sai_tam_int_attr_t_to_sai(
    lemming::dataplane::sai::TamIntAttr val);

lemming::dataplane::sai::AclCounterAttr convert_sai_acl_counter_attr_t_to_proto(
    const sai_int32_t val);
sai_acl_counter_attr_t convert_sai_acl_counter_attr_t_to_sai(
    lemming::dataplane::sai::AclCounterAttr val);

lemming::dataplane::sai::BfdSessionAttr convert_sai_bfd_session_attr_t_to_proto(
    const sai_int32_t val);
sai_bfd_session_attr_t convert_sai_bfd_session_attr_t_to_sai(
    lemming::dataplane::sai::BfdSessionAttr val);

lemming::dataplane::sai::HashAttr convert_sai_hash_attr_t_to_proto(
    const sai_int32_t val);
sai_hash_attr_t convert_sai_hash_attr_t_to_sai(
    lemming::dataplane::sai::HashAttr val);

lemming::dataplane::sai::AclTableGroupAttr
convert_sai_acl_table_group_attr_t_to_proto(const sai_int32_t val);
sai_acl_table_group_attr_t convert_sai_acl_table_group_attr_t_to_sai(
    lemming::dataplane::sai::AclTableGroupAttr val);

lemming::dataplane::sai::AclTableAttr convert_sai_acl_table_attr_t_to_proto(
    const sai_int32_t val);
sai_acl_table_attr_t convert_sai_acl_table_attr_t_to_sai(
    lemming::dataplane::sai::AclTableAttr val);

lemming::dataplane::sai::TamAttr convert_sai_tam_attr_t_to_proto(
    const sai_int32_t val);
sai_tam_attr_t convert_sai_tam_attr_t_to_sai(
    lemming::dataplane::sai::TamAttr val);

lemming::dataplane::sai::McastFdbEntryAttr
convert_sai_mcast_fdb_entry_attr_t_to_proto(const sai_int32_t val);
sai_mcast_fdb_entry_attr_t convert_sai_mcast_fdb_entry_attr_t_to_sai(
    lemming::dataplane::sai::McastFdbEntryAttr val);

lemming::dataplane::sai::QosMapAttr convert_sai_qos_map_attr_t_to_proto(
    const sai_int32_t val);
sai_qos_map_attr_t convert_sai_qos_map_attr_t_to_sai(
    lemming::dataplane::sai::QosMapAttr val);

lemming::dataplane::sai::SwitchAttr convert_sai_switch_attr_t_to_proto(
    const sai_int32_t val);
sai_switch_attr_t convert_sai_switch_attr_t_to_sai(
    lemming::dataplane::sai::SwitchAttr val);

lemming::dataplane::sai::PolicerAttr convert_sai_policer_attr_t_to_proto(
    const sai_int32_t val);
sai_policer_attr_t convert_sai_policer_attr_t_to_sai(
    lemming::dataplane::sai::PolicerAttr val);

lemming::dataplane::sai::RpfGroupAttr convert_sai_rpf_group_attr_t_to_proto(
    const sai_int32_t val);
sai_rpf_group_attr_t convert_sai_rpf_group_attr_t_to_sai(
    lemming::dataplane::sai::RpfGroupAttr val);

lemming::dataplane::sai::CounterAttr convert_sai_counter_attr_t_to_proto(
    const sai_int32_t val);
sai_counter_attr_t convert_sai_counter_attr_t_to_sai(
    lemming::dataplane::sai::CounterAttr val);

lemming::dataplane::sai::PortPoolAttr convert_sai_port_pool_attr_t_to_proto(
    const sai_int32_t val);
sai_port_pool_attr_t convert_sai_port_pool_attr_t_to_sai(
    lemming::dataplane::sai::PortPoolAttr val);

lemming::dataplane::sai::PortConnectorAttr
convert_sai_port_connector_attr_t_to_proto(const sai_int32_t val);
sai_port_connector_attr_t convert_sai_port_connector_attr_t_to_sai(
    lemming::dataplane::sai::PortConnectorAttr val);

lemming::dataplane::sai::MySidEntryAttr
convert_sai_my_sid_entry_attr_t_to_proto(const sai_int32_t val);
sai_my_sid_entry_attr_t convert_sai_my_sid_entry_attr_t_to_sai(
    lemming::dataplane::sai::MySidEntryAttr val);

lemming::dataplane::sai::NatEntryAttr convert_sai_nat_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_nat_entry_attr_t convert_sai_nat_entry_attr_t_to_sai(
    lemming::dataplane::sai::NatEntryAttr val);

lemming::dataplane::sai::NatZoneCounterAttr
convert_sai_nat_zone_counter_attr_t_to_proto(const sai_int32_t val);
sai_nat_zone_counter_attr_t convert_sai_nat_zone_counter_attr_t_to_sai(
    lemming::dataplane::sai::NatZoneCounterAttr val);

lemming::dataplane::sai::HostifUserDefinedTrapAttr
convert_sai_hostif_user_defined_trap_attr_t_to_proto(const sai_int32_t val);
sai_hostif_user_defined_trap_attr_t
convert_sai_hostif_user_defined_trap_attr_t_to_sai(
    lemming::dataplane::sai::HostifUserDefinedTrapAttr val);

lemming::dataplane::sai::L2mcEntryAttr convert_sai_l2mc_entry_attr_t_to_proto(
    const sai_int32_t val);
sai_l2mc_entry_attr_t convert_sai_l2mc_entry_attr_t_to_sai(
    lemming::dataplane::sai::L2mcEntryAttr val);

lemming::dataplane::sai::UdfGroupAttr convert_sai_udf_group_attr_t_to_proto(
    const sai_int32_t val);
sai_udf_group_attr_t convert_sai_udf_group_attr_t_to_sai(
    lemming::dataplane::sai::UdfGroupAttr val);

lemming::dataplane::sai::VlanMemberAttr convert_sai_vlan_member_attr_t_to_proto(
    const sai_int32_t val);
sai_vlan_member_attr_t convert_sai_vlan_member_attr_t_to_sai(
    lemming::dataplane::sai::VlanMemberAttr val);

lemming::dataplane::sai::FineGrainedHashFieldAttr
convert_sai_fine_grained_hash_field_attr_t_to_proto(const sai_int32_t val);
sai_fine_grained_hash_field_attr_t
convert_sai_fine_grained_hash_field_attr_t_to_sai(
    lemming::dataplane::sai::FineGrainedHashFieldAttr val);

lemming::dataplane::sai::IpmcGroupMemberAttr
convert_sai_ipmc_group_member_attr_t_to_proto(const sai_int32_t val);
sai_ipmc_group_member_attr_t convert_sai_ipmc_group_member_attr_t_to_sai(
    lemming::dataplane::sai::IpmcGroupMemberAttr val);

lemming::dataplane::sai::TunnelMapAttr convert_sai_tunnel_map_attr_t_to_proto(
    const sai_int32_t val);
sai_tunnel_map_attr_t convert_sai_tunnel_map_attr_t_to_sai(
    lemming::dataplane::sai::TunnelMapAttr val);

lemming::dataplane::sai::AclRangeAttr convert_sai_acl_range_attr_t_to_proto(
    const sai_int32_t val);
sai_acl_range_attr_t convert_sai_acl_range_attr_t_to_sai(
    lemming::dataplane::sai::AclRangeAttr val);

lemming::dataplane::sai::DtelQueueReportAttr
convert_sai_dtel_queue_report_attr_t_to_proto(const sai_int32_t val);
sai_dtel_queue_report_attr_t convert_sai_dtel_queue_report_attr_t_to_sai(
    lemming::dataplane::sai::DtelQueueReportAttr val);

lemming::dataplane::sai::HostifTrapGroupAttr
convert_sai_hostif_trap_group_attr_t_to_proto(const sai_int32_t val);
sai_hostif_trap_group_attr_t convert_sai_hostif_trap_group_attr_t_to_sai(
    lemming::dataplane::sai::HostifTrapGroupAttr val);

lemming::dataplane::sai::LagAttr convert_sai_lag_attr_t_to_proto(
    const sai_int32_t val);
sai_lag_attr_t convert_sai_lag_attr_t_to_sai(
    lemming::dataplane::sai::LagAttr val);

lemming::dataplane::sai::RpfGroupMemberAttr
convert_sai_rpf_group_member_attr_t_to_proto(const sai_int32_t val);
sai_rpf_group_member_attr_t convert_sai_rpf_group_member_attr_t_to_sai(
    lemming::dataplane::sai::RpfGroupMemberAttr val);

lemming::dataplane::sai::StpPortAttr convert_sai_stp_port_attr_t_to_proto(
    const sai_int32_t val);
sai_stp_port_attr_t convert_sai_stp_port_attr_t_to_sai(
    lemming::dataplane::sai::StpPortAttr val);

lemming::dataplane::sai::TableBitmapRouterEntryAttr
convert_sai_table_bitmap_router_entry_attr_t_to_proto(const sai_int32_t val);
sai_table_bitmap_router_entry_attr_t
convert_sai_table_bitmap_router_entry_attr_t_to_sai(
    lemming::dataplane::sai::TableBitmapRouterEntryAttr val);

lemming::dataplane::sai::LagMemberAttr convert_sai_lag_member_attr_t_to_proto(
    const sai_int32_t val);
sai_lag_member_attr_t convert_sai_lag_member_attr_t_to_sai(
    lemming::dataplane::sai::LagMemberAttr val);

lemming::dataplane::sai::PortSerdesAttr convert_sai_port_serdes_attr_t_to_proto(
    const sai_int32_t val);
sai_port_serdes_attr_t convert_sai_port_serdes_attr_t_to_sai(
    lemming::dataplane::sai::PortSerdesAttr val);

lemming::dataplane::sai::TamReportAttr convert_sai_tam_report_attr_t_to_proto(
    const sai_int32_t val);
sai_tam_report_attr_t convert_sai_tam_report_attr_t_to_sai(
    lemming::dataplane::sai::TamReportAttr val);

lemming::dataplane::sai::TamTransportAttr
convert_sai_tam_transport_attr_t_to_proto(const sai_int32_t val);
sai_tam_transport_attr_t convert_sai_tam_transport_attr_t_to_sai(
    lemming::dataplane::sai::TamTransportAttr val);

lemming::dataplane::sai::VlanAttr convert_sai_vlan_attr_t_to_proto(
    const sai_int32_t val);
sai_vlan_attr_t convert_sai_vlan_attr_t_to_sai(
    lemming::dataplane::sai::VlanAttr val);

lemming::dataplane::sai::IngressPriorityGroupAttr
convert_sai_ingress_priority_group_attr_t_to_proto(const sai_int32_t val);
sai_ingress_priority_group_attr_t
convert_sai_ingress_priority_group_attr_t_to_sai(
    lemming::dataplane::sai::IngressPriorityGroupAttr val);

lemming::dataplane::sai::MacsecAttr convert_sai_macsec_attr_t_to_proto(
    const sai_int32_t val);
sai_macsec_attr_t convert_sai_macsec_attr_t_to_sai(
    lemming::dataplane::sai::MacsecAttr val);

lemming::dataplane::sai::SwitchTunnelAttr
convert_sai_switch_tunnel_attr_t_to_proto(const sai_int32_t val);
sai_switch_tunnel_attr_t convert_sai_switch_tunnel_attr_t_to_sai(
    lemming::dataplane::sai::SwitchTunnelAttr val);

lemming::dataplane::sai::TamEventActionAttr
convert_sai_tam_event_action_attr_t_to_proto(const sai_int32_t val);
sai_tam_event_action_attr_t convert_sai_tam_event_action_attr_t_to_sai(
    lemming::dataplane::sai::TamEventActionAttr val);

lemming::dataplane::sai::DtelEventAttr convert_sai_dtel_event_attr_t_to_proto(
    const sai_int32_t val);
sai_dtel_event_attr_t convert_sai_dtel_event_attr_t_to_sai(
    lemming::dataplane::sai::DtelEventAttr val);

lemming::dataplane::sai::TableMetaTunnelEntryAttr
convert_sai_table_meta_tunnel_entry_attr_t_to_proto(const sai_int32_t val);
sai_table_meta_tunnel_entry_attr_t
convert_sai_table_meta_tunnel_entry_attr_t_to_sai(
    lemming::dataplane::sai::TableMetaTunnelEntryAttr val);

lemming::dataplane::sai::IpsecPortAttr convert_sai_ipsec_port_attr_t_to_proto(
    const sai_int32_t val);
sai_ipsec_port_attr_t convert_sai_ipsec_port_attr_t_to_sai(
    lemming::dataplane::sai::IpsecPortAttr val);

lemming::dataplane::sai::StpAttr convert_sai_stp_attr_t_to_proto(
    const sai_int32_t val);
sai_stp_attr_t convert_sai_stp_attr_t_to_sai(
    lemming::dataplane::sai::StpAttr val);

lemming::dataplane::sai::TamTelemetryAttr
convert_sai_tam_telemetry_attr_t_to_proto(const sai_int32_t val);
sai_tam_telemetry_attr_t convert_sai_tam_telemetry_attr_t_to_sai(
    lemming::dataplane::sai::TamTelemetryAttr val);

lemming::dataplane::sai::QueueAttr convert_sai_queue_attr_t_to_proto(
    const sai_int32_t val);
sai_queue_attr_t convert_sai_queue_attr_t_to_sai(
    lemming::dataplane::sai::QueueAttr val);

lemming::dataplane::sai::TamTelTypeAttr
convert_sai_tam_tel_type_attr_t_to_proto(const sai_int32_t val);
sai_tam_tel_type_attr_t convert_sai_tam_tel_type_attr_t_to_sai(
    lemming::dataplane::sai::TamTelTypeAttr val);

lemming::dataplane::sai::BufferProfileAttr
convert_sai_buffer_profile_attr_t_to_proto(const sai_int32_t val);
sai_buffer_profile_attr_t convert_sai_buffer_profile_attr_t_to_sai(
    lemming::dataplane::sai::BufferProfileAttr val);

lemming::dataplane::sai::L2mcGroupMemberAttr
convert_sai_l2mc_group_member_attr_t_to_proto(const sai_int32_t val);
sai_l2mc_group_member_attr_t convert_sai_l2mc_group_member_attr_t_to_sai(
    lemming::dataplane::sai::L2mcGroupMemberAttr val);

lemming::dataplane::sai::SamplepacketAttr
convert_sai_samplepacket_attr_t_to_proto(const sai_int32_t val);
sai_samplepacket_attr_t convert_sai_samplepacket_attr_t_to_sai(
    lemming::dataplane::sai::SamplepacketAttr val);

lemming::dataplane::sai::NextHopGroupMemberAttr
convert_sai_next_hop_group_member_attr_t_to_proto(const sai_int32_t val);
sai_next_hop_group_member_attr_t
convert_sai_next_hop_group_member_attr_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMemberAttr val);

lemming::dataplane::sai::IpsecAttr convert_sai_ipsec_attr_t_to_proto(
    const sai_int32_t val);
sai_ipsec_attr_t convert_sai_ipsec_attr_t_to_sai(
    lemming::dataplane::sai::IpsecAttr val);

lemming::dataplane::sai::MacsecScAttr convert_sai_macsec_sc_attr_t_to_proto(
    const sai_int32_t val);
sai_macsec_sc_attr_t convert_sai_macsec_sc_attr_t_to_sai(
    lemming::dataplane::sai::MacsecScAttr val);

lemming::dataplane::sai::NextHopGroupAttr
convert_sai_next_hop_group_attr_t_to_proto(const sai_int32_t val);
sai_next_hop_group_attr_t convert_sai_next_hop_group_attr_t_to_sai(
    lemming::dataplane::sai::NextHopGroupAttr val);

lemming::dataplane::sai::BridgePortAttr convert_sai_bridge_port_attr_t_to_proto(
    const sai_int32_t val);
sai_bridge_port_attr_t convert_sai_bridge_port_attr_t_to_sai(
    lemming::dataplane::sai::BridgePortAttr val);

lemming::dataplane::sai::AclTableGroupMemberAttr
convert_sai_acl_table_group_member_attr_t_to_proto(const sai_int32_t val);
sai_acl_table_group_member_attr_t
convert_sai_acl_table_group_member_attr_t_to_sai(
    lemming::dataplane::sai::AclTableGroupMemberAttr val);

lemming::dataplane::sai::IpmcGroupAttr convert_sai_ipmc_group_attr_t_to_proto(
    const sai_int32_t val);
sai_ipmc_group_attr_t convert_sai_ipmc_group_attr_t_to_sai(
    lemming::dataplane::sai::IpmcGroupAttr val);

lemming::dataplane::sai::WredAttr convert_sai_wred_attr_t_to_proto(
    const sai_int32_t val);
sai_wred_attr_t convert_sai_wred_attr_t_to_sai(
    lemming::dataplane::sai::WredAttr val);

lemming::dataplane::sai::DtelReportSessionAttr
convert_sai_dtel_report_session_attr_t_to_proto(const sai_int32_t val);
sai_dtel_report_session_attr_t convert_sai_dtel_report_session_attr_t_to_sai(
    lemming::dataplane::sai::DtelReportSessionAttr val);

lemming::dataplane::sai::MacsecPortAttr convert_sai_macsec_port_attr_t_to_proto(
    const sai_int32_t val);
sai_macsec_port_attr_t convert_sai_macsec_port_attr_t_to_sai(
    lemming::dataplane::sai::MacsecPortAttr val);

lemming::dataplane::sai::NeighborEntryAttr
convert_sai_neighbor_entry_attr_t_to_proto(const sai_int32_t val);
sai_neighbor_entry_attr_t convert_sai_neighbor_entry_attr_t_to_sai(
    lemming::dataplane::sai::NeighborEntryAttr val);

lemming::dataplane::sai::MacsecSaAttr convert_sai_macsec_sa_attr_t_to_proto(
    const sai_int32_t val);
sai_macsec_sa_attr_t convert_sai_macsec_sa_attr_t_to_sai(
    lemming::dataplane::sai::MacsecSaAttr val);

lemming::dataplane::sai::SchedulerGroupAttr
convert_sai_scheduler_group_attr_t_to_proto(const sai_int32_t val);
sai_scheduler_group_attr_t convert_sai_scheduler_group_attr_t_to_sai(
    lemming::dataplane::sai::SchedulerGroupAttr val);

lemming::dataplane::sai::TamEventThresholdAttr
convert_sai_tam_event_threshold_attr_t_to_proto(const sai_int32_t val);
sai_tam_event_threshold_attr_t convert_sai_tam_event_threshold_attr_t_to_sai(
    lemming::dataplane::sai::TamEventThresholdAttr val);

#endif  // DATAPLANE_STANDALONE_SAI_ENUM_H_
