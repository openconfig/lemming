

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
#include "dataplane/standalone/sai/entry.h"

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

  return translator->create(SAI_OBJECT_TYPE_TUNNEL_MAP, tunnel_map_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tunnel_map(sai_object_id_t tunnel_map_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TUNNEL_MAP, tunnel_map_id);
}

sai_status_t l_set_tunnel_map_attribute(sai_object_id_t tunnel_map_id,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TUNNEL_MAP, tunnel_map_id,
                                   attr);
}

sai_status_t l_get_tunnel_map_attribute(sai_object_id_t tunnel_map_id,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TUNNEL_MAP, tunnel_map_id,
                                   attr_count, attr_list);
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
    }
  }
  grpc::Status status = tunnel->CreateTunnel(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tunnel_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TUNNEL, tunnel_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_tunnel(sai_object_id_t tunnel_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TUNNEL, tunnel_id);
}

sai_status_t l_set_tunnel_attribute(sai_object_id_t tunnel_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TUNNEL, tunnel_id, attr);
}

sai_status_t l_get_tunnel_attribute(sai_object_id_t tunnel_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TUNNEL, tunnel_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_tunnel_stats(sai_object_id_t tunnel_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids,
                                uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_TUNNEL, tunnel_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_tunnel_stats_ext(sai_object_id_t tunnel_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t *counter_ids,
                                    sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_TUNNEL, tunnel_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_tunnel_stats(sai_object_id_t tunnel_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_TUNNEL, tunnel_id,
                                 number_of_counters, counter_ids);
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

  return translator->create(SAI_OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY,
                            tunnel_term_table_entry_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_tunnel_term_table_entry(
    sai_object_id_t tunnel_term_table_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY,
                            tunnel_term_table_entry_id);
}

sai_status_t l_set_tunnel_term_table_entry_attribute(
    sai_object_id_t tunnel_term_table_entry_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY,
                                   tunnel_term_table_entry_id, attr);
}

sai_status_t l_get_tunnel_term_table_entry_attribute(
    sai_object_id_t tunnel_term_table_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY,
                                   tunnel_term_table_entry_id, attr_count,
                                   attr_list);
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

  return translator->create(SAI_OBJECT_TYPE_TUNNEL_MAP_ENTRY,
                            tunnel_map_entry_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_tunnel_map_entry(sai_object_id_t tunnel_map_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TUNNEL_MAP_ENTRY,
                            tunnel_map_entry_id);
}

sai_status_t l_set_tunnel_map_entry_attribute(
    sai_object_id_t tunnel_map_entry_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TUNNEL_MAP_ENTRY,
                                   tunnel_map_entry_id, attr);
}

sai_status_t l_get_tunnel_map_entry_attribute(
    sai_object_id_t tunnel_map_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TUNNEL_MAP_ENTRY,
                                   tunnel_map_entry_id, attr_count, attr_list);
}
