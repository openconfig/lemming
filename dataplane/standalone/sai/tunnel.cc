

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

#include "dataplane/standalone/sai/tunnel.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/tunnel.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_tunnel_api_t l_tunnel = {
    .create_tunnel_map = l_create_tunnel_map,
    .remove_tunnel_map = l_remove_tunnel_map,
    .set_tunnel_map_attribute = l_set_tunnel_map_attribute,
    .get_tunnel_map_attribute = l_get_tunnel_map_attribute,
    .create_tunnel = l_create_tunnel,
    .remove_tunnel = l_remove_tunnel,
    .set_tunnel_attribute = l_set_tunnel_attribute,
    .get_tunnel_attribute = l_get_tunnel_attribute,
    .get_tunnel_stats = l_get_tunnel_stats,
    .get_tunnel_stats_ext = l_get_tunnel_stats_ext,
    .clear_tunnel_stats = l_clear_tunnel_stats,
    .create_tunnel_term_table_entry = l_create_tunnel_term_table_entry,
    .remove_tunnel_term_table_entry = l_remove_tunnel_term_table_entry,
    .set_tunnel_term_table_entry_attribute =
        l_set_tunnel_term_table_entry_attribute,
    .get_tunnel_term_table_entry_attribute =
        l_get_tunnel_term_table_entry_attribute,
    .create_tunnel_map_entry = l_create_tunnel_map_entry,
    .remove_tunnel_map_entry = l_remove_tunnel_map_entry,
    .set_tunnel_map_entry_attribute = l_set_tunnel_map_entry_attribute,
    .get_tunnel_map_entry_attribute = l_get_tunnel_map_entry_attribute,
};

sai_status_t l_create_tunnel_map(sai_object_id_t *tunnel_map_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTunnelMapRequest req;
  lemming::dataplane::sai::CreateTunnelMapResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_MAP_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::TunnelMapType>(
            attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = tunnel->CreateTunnelMap(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tunnel_map_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tunnel_map(sai_object_id_t tunnel_map_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveTunnelMapRequest req;
  lemming::dataplane::sai::RemoveTunnelMapResponse resp;
  grpc::ClientContext context;
  req.set_oid(tunnel_map_id);

  grpc::Status status = tunnel->RemoveTunnelMap(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tunnel_map_attribute(sai_object_id_t tunnel_map_id,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tunnel_map_attribute(sai_object_id_t tunnel_map_id,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetTunnelMapAttributeRequest req;
  lemming::dataplane::sai::GetTunnelMapAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(tunnel_map_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::TunnelMapAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = tunnel->GetTunnelMapAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_MAP_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_TUNNEL_MAP_ATTR_ENTRY_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().entry_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tunnel(sai_object_id_t *tunnel_id,
                             sai_object_id_t switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTunnelRequest req;
  lemming::dataplane::sai::CreateTunnelResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::TunnelType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_ATTR_UNDERLAY_INTERFACE:
        req.set_underlay_interface(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_OVERLAY_INTERFACE:
        req.set_overlay_interface(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_PEER_MODE:
        req.set_peer_mode(static_cast<lemming::dataplane::sai::TunnelPeerMode>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_SRC_IP:
        req.set_encap_src_ip(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_DST_IP:
        req.set_encap_dst_ip(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_TTL_MODE:
        req.set_encap_ttl_mode(
            static_cast<lemming::dataplane::sai::TunnelTtlMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_TTL_VAL:
        req.set_encap_ttl_val(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_DSCP_MODE:
        req.set_encap_dscp_mode(
            static_cast<lemming::dataplane::sai::TunnelDscpMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_DSCP_VAL:
        req.set_encap_dscp_val(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_GRE_KEY_VALID:
        req.set_encap_gre_key_valid(attr_list[i].value.booldata);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_GRE_KEY:
        req.set_encap_gre_key(attr_list[i].value.u32);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_ECN_MODE:
        req.set_encap_ecn_mode(
            static_cast<lemming::dataplane::sai::TunnelEncapEcnMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_MAPPERS:
        req.mutable_encap_mappers()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_DECAP_ECN_MODE:
        req.set_decap_ecn_mode(
            static_cast<lemming::dataplane::sai::TunnelDecapEcnMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_ATTR_DECAP_MAPPERS:
        req.mutable_decap_mappers()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_DECAP_TTL_MODE:
        req.set_decap_ttl_mode(
            static_cast<lemming::dataplane::sai::TunnelTtlMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_ATTR_DECAP_DSCP_MODE:
        req.set_decap_dscp_mode(
            static_cast<lemming::dataplane::sai::TunnelDscpMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
        req.set_loopback_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MODE:
        req.set_vxlan_udp_sport_mode(
            static_cast<lemming::dataplane::sai::TunnelVxlanUdpSportMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT:
        req.set_vxlan_udp_sport(attr_list[i].value.u16);
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
        req.set_vxlan_udp_sport_mask(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_ATTR_SA_INDEX:
        req.set_sa_index(attr_list[i].value.u32);
        break;
      case SAI_TUNNEL_ATTR_IPSEC_SA_PORT_LIST:
        req.mutable_ipsec_sa_port_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
        req.set_encap_qos_tc_and_color_to_dscp_map(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
        req.set_encap_qos_tc_to_queue_map(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
        req.set_decap_qos_dscp_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
        req.set_decap_qos_tc_to_priority_group_map(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_SECURITY:
        req.set_vxlan_udp_sport_security(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status = tunnel->CreateTunnel(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tunnel_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tunnel(sai_object_id_t tunnel_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveTunnelRequest req;
  lemming::dataplane::sai::RemoveTunnelResponse resp;
  grpc::ClientContext context;
  req.set_oid(tunnel_id);

  grpc::Status status = tunnel->RemoveTunnel(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tunnel_attribute(sai_object_id_t tunnel_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetTunnelAttributeRequest req;
  lemming::dataplane::sai::SetTunnelAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(tunnel_id);

  switch (attr->id) {
    case SAI_TUNNEL_ATTR_ENCAP_TTL_MODE:
      req.set_encap_ttl_mode(
          static_cast<lemming::dataplane::sai::TunnelTtlMode>(attr->value.s32 +
                                                              1));
      break;
    case SAI_TUNNEL_ATTR_ENCAP_TTL_VAL:
      req.set_encap_ttl_val(attr->value.u8);
      break;
    case SAI_TUNNEL_ATTR_ENCAP_DSCP_MODE:
      req.set_encap_dscp_mode(
          static_cast<lemming::dataplane::sai::TunnelDscpMode>(attr->value.s32 +
                                                               1));
      break;
    case SAI_TUNNEL_ATTR_ENCAP_DSCP_VAL:
      req.set_encap_dscp_val(attr->value.u8);
      break;
    case SAI_TUNNEL_ATTR_ENCAP_GRE_KEY:
      req.set_encap_gre_key(attr->value.u32);
      break;
    case SAI_TUNNEL_ATTR_DECAP_TTL_MODE:
      req.set_decap_ttl_mode(
          static_cast<lemming::dataplane::sai::TunnelTtlMode>(attr->value.s32 +
                                                              1));
      break;
    case SAI_TUNNEL_ATTR_DECAP_DSCP_MODE:
      req.set_decap_dscp_mode(
          static_cast<lemming::dataplane::sai::TunnelDscpMode>(attr->value.s32 +
                                                               1));
      break;
    case SAI_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
      req.set_loopback_packet_action(
          static_cast<lemming::dataplane::sai::PacketAction>(attr->value.s32 +
                                                             1));
      break;
    case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MODE:
      req.set_vxlan_udp_sport_mode(
          static_cast<lemming::dataplane::sai::TunnelVxlanUdpSportMode>(
              attr->value.s32 + 1));
      break;
    case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT:
      req.set_vxlan_udp_sport(attr->value.u16);
      break;
    case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
      req.set_vxlan_udp_sport_mask(attr->value.u8);
      break;
    case SAI_TUNNEL_ATTR_SA_INDEX:
      req.set_sa_index(attr->value.u32);
      break;
    case SAI_TUNNEL_ATTR_IPSEC_SA_PORT_LIST:
      req.mutable_ipsec_sa_port_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
      req.set_encap_qos_tc_and_color_to_dscp_map(attr->value.oid);
      break;
    case SAI_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
      req.set_encap_qos_tc_to_queue_map(attr->value.oid);
      break;
    case SAI_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
      req.set_decap_qos_dscp_to_tc_map(attr->value.oid);
      break;
    case SAI_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
      req.set_decap_qos_tc_to_priority_group_map(attr->value.oid);
      break;
    case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_SECURITY:
      req.set_vxlan_udp_sport_security(attr->value.booldata);
      break;
  }

  grpc::Status status = tunnel->SetTunnelAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tunnel_attribute(sai_object_id_t tunnel_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetTunnelAttributeRequest req;
  lemming::dataplane::sai::GetTunnelAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(tunnel_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::TunnelAttr>(attr_list[i].id + 1));
  }
  grpc::Status status = tunnel->GetTunnelAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_TUNNEL_ATTR_UNDERLAY_INTERFACE:
        attr_list[i].value.oid = resp.attr().underlay_interface();
        break;
      case SAI_TUNNEL_ATTR_OVERLAY_INTERFACE:
        attr_list[i].value.oid = resp.attr().overlay_interface();
        break;
      case SAI_TUNNEL_ATTR_PEER_MODE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().peer_mode() - 1);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_SRC_IP:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().encap_src_ip());
        break;
      case SAI_TUNNEL_ATTR_ENCAP_DST_IP:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().encap_dst_ip());
        break;
      case SAI_TUNNEL_ATTR_ENCAP_TTL_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().encap_ttl_mode() - 1);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_TTL_VAL:
        attr_list[i].value.u8 = resp.attr().encap_ttl_val();
        break;
      case SAI_TUNNEL_ATTR_ENCAP_DSCP_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().encap_dscp_mode() - 1);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_DSCP_VAL:
        attr_list[i].value.u8 = resp.attr().encap_dscp_val();
        break;
      case SAI_TUNNEL_ATTR_ENCAP_GRE_KEY_VALID:
        attr_list[i].value.booldata = resp.attr().encap_gre_key_valid();
        break;
      case SAI_TUNNEL_ATTR_ENCAP_GRE_KEY:
        attr_list[i].value.u32 = resp.attr().encap_gre_key();
        break;
      case SAI_TUNNEL_ATTR_ENCAP_ECN_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().encap_ecn_mode() - 1);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_MAPPERS:
        copy_list(attr_list[i].value.objlist.list, resp.attr().encap_mappers(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_DECAP_ECN_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().decap_ecn_mode() - 1);
        break;
      case SAI_TUNNEL_ATTR_DECAP_MAPPERS:
        copy_list(attr_list[i].value.objlist.list, resp.attr().decap_mappers(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_DECAP_TTL_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().decap_ttl_mode() - 1);
        break;
      case SAI_TUNNEL_ATTR_DECAP_DSCP_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().decap_dscp_mode() - 1);
        break;
      case SAI_TUNNEL_ATTR_TERM_TABLE_ENTRY_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().term_table_entry_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().loopback_packet_action() - 1);
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().vxlan_udp_sport_mode() - 1);
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT:
        attr_list[i].value.u16 = resp.attr().vxlan_udp_sport();
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
        attr_list[i].value.u8 = resp.attr().vxlan_udp_sport_mask();
        break;
      case SAI_TUNNEL_ATTR_SA_INDEX:
        attr_list[i].value.u32 = resp.attr().sa_index();
        break;
      case SAI_TUNNEL_ATTR_IPSEC_SA_PORT_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().ipsec_sa_port_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
        attr_list[i].value.oid =
            resp.attr().encap_qos_tc_and_color_to_dscp_map();
        break;
      case SAI_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
        attr_list[i].value.oid = resp.attr().encap_qos_tc_to_queue_map();
        break;
      case SAI_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
        attr_list[i].value.oid = resp.attr().decap_qos_dscp_to_tc_map();
        break;
      case SAI_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
        attr_list[i].value.oid =
            resp.attr().decap_qos_tc_to_priority_group_map();
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_SECURITY:
        attr_list[i].value.booldata = resp.attr().vxlan_udp_sport_security();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tunnel_stats(sai_object_id_t tunnel_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids,
                                uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tunnel_stats_ext(sai_object_id_t tunnel_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t *counter_ids,
                                    sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_tunnel_stats(sai_object_id_t tunnel_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tunnel_term_table_entry(
    sai_object_id_t *tunnel_term_table_entry_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTunnelTermTableEntryRequest req;
  lemming::dataplane::sai::CreateTunnelTermTableEntryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_VR_ID:
        req.set_vr_id(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_TYPE:
        req.set_type(
            static_cast<lemming::dataplane::sai::TunnelTermTableEntryType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP:
        req.set_dst_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP_MASK:
        req.set_dst_ip_mask(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP:
        req.set_src_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP_MASK:
        req.set_src_ip_mask(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_TUNNEL_TYPE:
        req.set_tunnel_type(static_cast<lemming::dataplane::sai::TunnelType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_ACTION_TUNNEL_ID:
        req.set_action_tunnel_id(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_IPSEC_VERIFIED:
        req.set_ipsec_verified(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status =
      tunnel->CreateTunnelTermTableEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tunnel_term_table_entry_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tunnel_term_table_entry(
    sai_object_id_t tunnel_term_table_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveTunnelTermTableEntryRequest req;
  lemming::dataplane::sai::RemoveTunnelTermTableEntryResponse resp;
  grpc::ClientContext context;
  req.set_oid(tunnel_term_table_entry_id);

  grpc::Status status =
      tunnel->RemoveTunnelTermTableEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tunnel_term_table_entry_attribute(
    sai_object_id_t tunnel_term_table_entry_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetTunnelTermTableEntryAttributeRequest req;
  lemming::dataplane::sai::SetTunnelTermTableEntryAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(tunnel_term_table_entry_id);

  switch (attr->id) {
    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_IPSEC_VERIFIED:
      req.set_ipsec_verified(attr->value.booldata);
      break;
  }

  grpc::Status status =
      tunnel->SetTunnelTermTableEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tunnel_term_table_entry_attribute(
    sai_object_id_t tunnel_term_table_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetTunnelTermTableEntryAttributeRequest req;
  lemming::dataplane::sai::GetTunnelTermTableEntryAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(tunnel_term_table_entry_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::TunnelTermTableEntryAttr>(
            attr_list[i].id + 1));
  }
  grpc::Status status =
      tunnel->GetTunnelTermTableEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_VR_ID:
        attr_list[i].value.oid = resp.attr().vr_id();
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP:
        attr_list[i].value.ipaddr = convert_to_ip_address(resp.attr().dst_ip());
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP_MASK:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().dst_ip_mask());
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP:
        attr_list[i].value.ipaddr = convert_to_ip_address(resp.attr().src_ip());
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP_MASK:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().src_ip_mask());
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_TUNNEL_TYPE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().tunnel_type() - 1);
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_ACTION_TUNNEL_ID:
        attr_list[i].value.oid = resp.attr().action_tunnel_id();
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_IP_ADDR_FAMILY:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().ip_addr_family() - 1);
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_IPSEC_VERIFIED:
        attr_list[i].value.booldata = resp.attr().ipsec_verified();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tunnel_map_entry(sai_object_id_t *tunnel_map_entry_id,
                                       sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTunnelMapEntryRequest req;
  lemming::dataplane::sai::CreateTunnelMapEntryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP_TYPE:
        req.set_tunnel_map_type(
            static_cast<lemming::dataplane::sai::TunnelMapType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP:
        req.set_tunnel_map(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_OECN_KEY:
        req.set_oecn_key(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_OECN_VALUE:
        req.set_oecn_value(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_UECN_KEY:
        req.set_uecn_key(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_UECN_VALUE:
        req.set_uecn_value(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_KEY:
        req.set_vlan_id_key(attr_list[i].value.u16);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_VALUE:
        req.set_vlan_id_value(attr_list[i].value.u16);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VNI_ID_KEY:
        req.set_vni_id_key(attr_list[i].value.u32);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VNI_ID_VALUE:
        req.set_vni_id_value(attr_list[i].value.u32);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_KEY:
        req.set_bridge_id_key(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_VALUE:
        req.set_bridge_id_value(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_KEY:
        req.set_virtual_router_id_key(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_VALUE:
        req.set_virtual_router_id_value(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VSID_ID_KEY:
        req.set_vsid_id_key(attr_list[i].value.u32);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VSID_ID_VALUE:
        req.set_vsid_id_value(attr_list[i].value.u32);
        break;
    }
  }
  grpc::Status status = tunnel->CreateTunnelMapEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tunnel_map_entry_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tunnel_map_entry(sai_object_id_t tunnel_map_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveTunnelMapEntryRequest req;
  lemming::dataplane::sai::RemoveTunnelMapEntryResponse resp;
  grpc::ClientContext context;
  req.set_oid(tunnel_map_entry_id);

  grpc::Status status = tunnel->RemoveTunnelMapEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tunnel_map_entry_attribute(
    sai_object_id_t tunnel_map_entry_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tunnel_map_entry_attribute(
    sai_object_id_t tunnel_map_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetTunnelMapEntryAttributeRequest req;
  lemming::dataplane::sai::GetTunnelMapEntryAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(tunnel_map_entry_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::TunnelMapEntryAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status =
      tunnel->GetTunnelMapEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP_TYPE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().tunnel_map_type() - 1);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP:
        attr_list[i].value.oid = resp.attr().tunnel_map();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_OECN_KEY:
        attr_list[i].value.u8 = resp.attr().oecn_key();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_OECN_VALUE:
        attr_list[i].value.u8 = resp.attr().oecn_value();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_UECN_KEY:
        attr_list[i].value.u8 = resp.attr().uecn_key();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_UECN_VALUE:
        attr_list[i].value.u8 = resp.attr().uecn_value();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_KEY:
        attr_list[i].value.u16 = resp.attr().vlan_id_key();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_VALUE:
        attr_list[i].value.u16 = resp.attr().vlan_id_value();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VNI_ID_KEY:
        attr_list[i].value.u32 = resp.attr().vni_id_key();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VNI_ID_VALUE:
        attr_list[i].value.u32 = resp.attr().vni_id_value();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_KEY:
        attr_list[i].value.oid = resp.attr().bridge_id_key();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_VALUE:
        attr_list[i].value.oid = resp.attr().bridge_id_value();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_KEY:
        attr_list[i].value.oid = resp.attr().virtual_router_id_key();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_VALUE:
        attr_list[i].value.oid = resp.attr().virtual_router_id_value();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VSID_ID_KEY:
        attr_list[i].value.u32 = resp.attr().vsid_id_key();
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VSID_ID_VALUE:
        attr_list[i].value.u32 = resp.attr().vsid_id_value();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
