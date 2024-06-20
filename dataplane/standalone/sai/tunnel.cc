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

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/tunnel.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

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
    .create_tunnels = l_create_tunnels,
    .remove_tunnels = l_remove_tunnels,
    .set_tunnels_attribute = l_set_tunnels_attribute,
    .get_tunnels_attribute = l_get_tunnels_attribute,
};

lemming::dataplane::sai::CreateTunnelMapRequest convert_create_tunnel_map(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateTunnelMapRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_MAP_ATTR_TYPE:
        msg.set_type(
            convert_sai_tunnel_map_type_t_to_proto(attr_list[i].value.s32));
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateTunnelRequest convert_create_tunnel(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateTunnelRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_ATTR_TYPE:
        msg.set_type(
            convert_sai_tunnel_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_ATTR_UNDERLAY_INTERFACE:
        msg.set_underlay_interface(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_OVERLAY_INTERFACE:
        msg.set_overlay_interface(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_PEER_MODE:
        msg.set_peer_mode(
            convert_sai_tunnel_peer_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_SRC_IP:
        msg.set_encap_src_ip(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_DST_IP:
        msg.set_encap_dst_ip(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_TTL_MODE:
        msg.set_encap_ttl_mode(
            convert_sai_tunnel_ttl_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_TTL_VAL:
        msg.set_encap_ttl_val(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_DSCP_MODE:
        msg.set_encap_dscp_mode(
            convert_sai_tunnel_dscp_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_DSCP_VAL:
        msg.set_encap_dscp_val(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_GRE_KEY_VALID:
        msg.set_encap_gre_key_valid(attr_list[i].value.booldata);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_GRE_KEY:
        msg.set_encap_gre_key(attr_list[i].value.u32);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_ECN_MODE:
        msg.set_encap_ecn_mode(convert_sai_tunnel_encap_ecn_mode_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_ATTR_ENCAP_MAPPERS:
        msg.mutable_encap_mappers()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_DECAP_ECN_MODE:
        msg.set_decap_ecn_mode(convert_sai_tunnel_decap_ecn_mode_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_ATTR_DECAP_MAPPERS:
        msg.mutable_decap_mappers()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_DECAP_TTL_MODE:
        msg.set_decap_ttl_mode(
            convert_sai_tunnel_ttl_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_ATTR_DECAP_DSCP_MODE:
        msg.set_decap_dscp_mode(
            convert_sai_tunnel_dscp_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
        msg.set_loopback_packet_action(
            convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MODE:
        msg.set_vxlan_udp_sport_mode(
            convert_sai_tunnel_vxlan_udp_sport_mode_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT:
        msg.set_vxlan_udp_sport(attr_list[i].value.u16);
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
        msg.set_vxlan_udp_sport_mask(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_ATTR_SA_INDEX:
        msg.set_sa_index(attr_list[i].value.u32);
        break;
      case SAI_TUNNEL_ATTR_IPSEC_SA_PORT_LIST:
        msg.mutable_ipsec_sa_port_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
        msg.set_encap_qos_tc_and_color_to_dscp_map(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
        msg.set_encap_qos_tc_to_queue_map(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
        msg.set_decap_qos_dscp_to_tc_map(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
        msg.set_decap_qos_tc_to_priority_group_map(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_SECURITY:
        msg.set_vxlan_udp_sport_security(attr_list[i].value.booldata);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateTunnelTermTableEntryRequest
convert_create_tunnel_term_table_entry(sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateTunnelTermTableEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_VR_ID:
        msg.set_vr_id(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_TYPE:
        msg.set_type(convert_sai_tunnel_term_table_entry_type_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP:
        msg.set_dst_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP_MASK:
        msg.set_dst_ip_mask(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP:
        msg.set_src_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP_MASK:
        msg.set_src_ip_mask(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_TUNNEL_TYPE:
        msg.set_tunnel_type(
            convert_sai_tunnel_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_ACTION_TUNNEL_ID:
        msg.set_action_tunnel_id(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_IPSEC_VERIFIED:
        msg.set_ipsec_verified(attr_list[i].value.booldata);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateTunnelMapEntryRequest
convert_create_tunnel_map_entry(sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateTunnelMapEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP_TYPE:
        msg.set_tunnel_map_type(
            convert_sai_tunnel_map_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP:
        msg.set_tunnel_map(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_OECN_KEY:
        msg.set_oecn_key(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_OECN_VALUE:
        msg.set_oecn_value(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_UECN_KEY:
        msg.set_uecn_key(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_UECN_VALUE:
        msg.set_uecn_value(attr_list[i].value.u8);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_KEY:
        msg.set_vlan_id_key(attr_list[i].value.u16);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_VALUE:
        msg.set_vlan_id_value(attr_list[i].value.u16);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VNI_ID_KEY:
        msg.set_vni_id_key(attr_list[i].value.u32);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VNI_ID_VALUE:
        msg.set_vni_id_value(attr_list[i].value.u32);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_KEY:
        msg.set_bridge_id_key(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_VALUE:
        msg.set_bridge_id_value(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_KEY:
        msg.set_virtual_router_id_key(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_VALUE:
        msg.set_virtual_router_id_value(attr_list[i].value.oid);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VSID_ID_KEY:
        msg.set_vsid_id_key(attr_list[i].value.u32);
        break;
      case SAI_TUNNEL_MAP_ENTRY_ATTR_VSID_ID_VALUE:
        msg.set_vsid_id_value(attr_list[i].value.u32);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_tunnel_map(sai_object_id_t *tunnel_map_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTunnelMapRequest req =
      convert_create_tunnel_map(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateTunnelMapResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

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
    req.add_attr_type(convert_sai_tunnel_map_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = tunnel->GetTunnelMapAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_MAP_ATTR_TYPE:
        attr_list[i].value.s32 =
            convert_sai_tunnel_map_type_t_to_sai(resp.attr().type());
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

  lemming::dataplane::sai::CreateTunnelRequest req =
      convert_create_tunnel(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateTunnelResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

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
          convert_sai_tunnel_ttl_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_TUNNEL_ATTR_ENCAP_TTL_VAL:
      req.set_encap_ttl_val(attr->value.u8);
      break;
    case SAI_TUNNEL_ATTR_ENCAP_DSCP_MODE:
      req.set_encap_dscp_mode(
          convert_sai_tunnel_dscp_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_TUNNEL_ATTR_ENCAP_DSCP_VAL:
      req.set_encap_dscp_val(attr->value.u8);
      break;
    case SAI_TUNNEL_ATTR_ENCAP_GRE_KEY:
      req.set_encap_gre_key(attr->value.u32);
      break;
    case SAI_TUNNEL_ATTR_DECAP_TTL_MODE:
      req.set_decap_ttl_mode(
          convert_sai_tunnel_ttl_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_TUNNEL_ATTR_DECAP_DSCP_MODE:
      req.set_decap_dscp_mode(
          convert_sai_tunnel_dscp_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
      req.set_loopback_packet_action(
          convert_sai_packet_action_t_to_proto(attr->value.s32));
      break;
    case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MODE:
      req.set_vxlan_udp_sport_mode(
          convert_sai_tunnel_vxlan_udp_sport_mode_t_to_proto(attr->value.s32));
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
    req.add_attr_type(convert_sai_tunnel_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = tunnel->GetTunnelAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TUNNEL_ATTR_TYPE:
        attr_list[i].value.s32 =
            convert_sai_tunnel_type_t_to_sai(resp.attr().type());
        break;
      case SAI_TUNNEL_ATTR_UNDERLAY_INTERFACE:
        attr_list[i].value.oid = resp.attr().underlay_interface();
        break;
      case SAI_TUNNEL_ATTR_OVERLAY_INTERFACE:
        attr_list[i].value.oid = resp.attr().overlay_interface();
        break;
      case SAI_TUNNEL_ATTR_PEER_MODE:
        attr_list[i].value.s32 =
            convert_sai_tunnel_peer_mode_t_to_sai(resp.attr().peer_mode());
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
            convert_sai_tunnel_ttl_mode_t_to_sai(resp.attr().encap_ttl_mode());
        break;
      case SAI_TUNNEL_ATTR_ENCAP_TTL_VAL:
        attr_list[i].value.u8 = resp.attr().encap_ttl_val();
        break;
      case SAI_TUNNEL_ATTR_ENCAP_DSCP_MODE:
        attr_list[i].value.s32 = convert_sai_tunnel_dscp_mode_t_to_sai(
            resp.attr().encap_dscp_mode());
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
        attr_list[i].value.s32 = convert_sai_tunnel_encap_ecn_mode_t_to_sai(
            resp.attr().encap_ecn_mode());
        break;
      case SAI_TUNNEL_ATTR_ENCAP_MAPPERS:
        copy_list(attr_list[i].value.objlist.list, resp.attr().encap_mappers(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_DECAP_ECN_MODE:
        attr_list[i].value.s32 = convert_sai_tunnel_decap_ecn_mode_t_to_sai(
            resp.attr().decap_ecn_mode());
        break;
      case SAI_TUNNEL_ATTR_DECAP_MAPPERS:
        copy_list(attr_list[i].value.objlist.list, resp.attr().decap_mappers(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_DECAP_TTL_MODE:
        attr_list[i].value.s32 =
            convert_sai_tunnel_ttl_mode_t_to_sai(resp.attr().decap_ttl_mode());
        break;
      case SAI_TUNNEL_ATTR_DECAP_DSCP_MODE:
        attr_list[i].value.s32 = convert_sai_tunnel_dscp_mode_t_to_sai(
            resp.attr().decap_dscp_mode());
        break;
      case SAI_TUNNEL_ATTR_TERM_TABLE_ENTRY_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().term_table_entry_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
        attr_list[i].value.s32 = convert_sai_packet_action_t_to_sai(
            resp.attr().loopback_packet_action());
        break;
      case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MODE:
        attr_list[i].value.s32 =
            convert_sai_tunnel_vxlan_udp_sport_mode_t_to_sai(
                resp.attr().vxlan_udp_sport_mode());
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

  lemming::dataplane::sai::GetTunnelStatsRequest req;
  lemming::dataplane::sai::GetTunnelStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(tunnel_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(convert_sai_tunnel_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status = tunnel->GetTunnelStats(&context, req, &resp);
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

  lemming::dataplane::sai::CreateTunnelTermTableEntryRequest req =
      convert_create_tunnel_term_table_entry(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateTunnelTermTableEntryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

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
        convert_sai_tunnel_term_table_entry_attr_t_to_proto(attr_list[i].id));
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
        attr_list[i].value.s32 =
            convert_sai_tunnel_term_table_entry_type_t_to_sai(
                resp.attr().type());
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
            convert_sai_tunnel_type_t_to_sai(resp.attr().tunnel_type());
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_ACTION_TUNNEL_ID:
        attr_list[i].value.oid = resp.attr().action_tunnel_id();
        break;
      case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_IP_ADDR_FAMILY:
        attr_list[i].value.s32 =
            convert_sai_ip_addr_family_t_to_sai(resp.attr().ip_addr_family());
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

  lemming::dataplane::sai::CreateTunnelMapEntryRequest req =
      convert_create_tunnel_map_entry(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateTunnelMapEntryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

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
    req.add_attr_type(
        convert_sai_tunnel_map_entry_attr_t_to_proto(attr_list[i].id));
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
            convert_sai_tunnel_map_type_t_to_sai(resp.attr().tunnel_map_type());
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

sai_status_t l_create_tunnels(sai_object_id_t switch_id, uint32_t object_count,
                              const uint32_t *attr_count,
                              const sai_attribute_t **attr_list,
                              sai_bulk_op_error_mode_t mode,
                              sai_object_id_t *object_id,
                              sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTunnelsRequest req;
  lemming::dataplane::sai::CreateTunnelsResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    auto r = convert_create_tunnel(switch_id, attr_count[i], attr_list[i]);
    *req.add_reqs() = r;
  }

  grpc::Status status = tunnel->CreateTunnels(&context, req, &resp);
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

sai_status_t l_remove_tunnels(uint32_t object_count,
                              const sai_object_id_t *object_id,
                              sai_bulk_op_error_mode_t mode,
                              sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveTunnelsRequest req;
  lemming::dataplane::sai::RemoveTunnelsResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    req.add_reqs()->set_oid(object_id[i]);
  }

  grpc::Status status = tunnel->RemoveTunnels(&context, req, &resp);
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

sai_status_t l_set_tunnels_attribute(uint32_t object_count,
                                     const sai_object_id_t *object_id,
                                     const sai_attribute_t *attr_list,
                                     sai_bulk_op_error_mode_t mode,
                                     sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_tunnels_attribute(uint32_t object_count,
                                     const sai_object_id_t *object_id,
                                     const uint32_t *attr_count,
                                     sai_attribute_t **attr_list,
                                     sai_bulk_op_error_mode_t mode,
                                     sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}
