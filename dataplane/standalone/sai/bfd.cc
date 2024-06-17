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

#include "dataplane/standalone/sai/bfd.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/bfd.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_bfd_api_t l_bfd = {
    .create_bfd_session = l_create_bfd_session,
    .remove_bfd_session = l_remove_bfd_session,
    .set_bfd_session_attribute = l_set_bfd_session_attribute,
    .get_bfd_session_attribute = l_get_bfd_session_attribute,
    .get_bfd_session_stats = l_get_bfd_session_stats,
    .get_bfd_session_stats_ext = l_get_bfd_session_stats_ext,
    .clear_bfd_session_stats = l_clear_bfd_session_stats,
};

lemming::dataplane::sai::CreateBfdSessionRequest convert_create_bfd_session(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateBfdSessionRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BFD_SESSION_ATTR_TYPE:
        msg.set_type(static_cast<lemming::dataplane::sai::BfdSessionType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_BFD_SESSION_ATTR_HW_LOOKUP_VALID:
        msg.set_hw_lookup_valid(attr_list[i].value.booldata);
        break;
      case SAI_BFD_SESSION_ATTR_VIRTUAL_ROUTER:
        msg.set_virtual_router(attr_list[i].value.oid);
        break;
      case SAI_BFD_SESSION_ATTR_PORT:
        msg.set_port(attr_list[i].value.oid);
        break;
      case SAI_BFD_SESSION_ATTR_LOCAL_DISCRIMINATOR:
        msg.set_local_discriminator(attr_list[i].value.u32);
        break;
      case SAI_BFD_SESSION_ATTR_REMOTE_DISCRIMINATOR:
        msg.set_remote_discriminator(attr_list[i].value.u32);
        break;
      case SAI_BFD_SESSION_ATTR_UDP_SRC_PORT:
        msg.set_udp_src_port(attr_list[i].value.u32);
        break;
      case SAI_BFD_SESSION_ATTR_TC:
        msg.set_tc(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_TPID:
        msg.set_vlan_tpid(attr_list[i].value.u16);
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_ID:
        msg.set_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_PRI:
        msg.set_vlan_pri(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_CFI:
        msg.set_vlan_cfi(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_HEADER_VALID:
        msg.set_vlan_header_valid(attr_list[i].value.booldata);
        break;
      case SAI_BFD_SESSION_ATTR_BFD_ENCAPSULATION_TYPE:
        msg.set_bfd_encapsulation_type(
            static_cast<lemming::dataplane::sai::BfdEncapsulationType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BFD_SESSION_ATTR_IPHDR_VERSION:
        msg.set_iphdr_version(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_TOS:
        msg.set_tos(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_TTL:
        msg.set_ttl(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_SRC_IP_ADDRESS:
        msg.set_src_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_BFD_SESSION_ATTR_DST_IP_ADDRESS:
        msg.set_dst_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_TOS:
        msg.set_tunnel_tos(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_TTL:
        msg.set_tunnel_ttl(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_SRC_IP_ADDRESS:
        msg.set_tunnel_src_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_DST_IP_ADDRESS:
        msg.set_tunnel_dst_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_BFD_SESSION_ATTR_SRC_MAC_ADDRESS:
        msg.set_src_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_BFD_SESSION_ATTR_DST_MAC_ADDRESS:
        msg.set_dst_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_BFD_SESSION_ATTR_ECHO_ENABLE:
        msg.set_echo_enable(attr_list[i].value.booldata);
        break;
      case SAI_BFD_SESSION_ATTR_MULTIHOP:
        msg.set_multihop(attr_list[i].value.booldata);
        break;
      case SAI_BFD_SESSION_ATTR_CBIT:
        msg.set_cbit(attr_list[i].value.booldata);
        break;
      case SAI_BFD_SESSION_ATTR_MIN_TX:
        msg.set_min_tx(attr_list[i].value.u32);
        break;
      case SAI_BFD_SESSION_ATTR_MIN_RX:
        msg.set_min_rx(attr_list[i].value.u32);
        break;
      case SAI_BFD_SESSION_ATTR_MULTIPLIER:
        msg.set_multiplier(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_OFFLOAD_TYPE:
        msg.set_offload_type(
            static_cast<lemming::dataplane::sai::BfdSessionOffloadType>(
                attr_list[i].value.s32 + 1));
        break;
    }
  }
  return msg;
}

sai_status_t l_create_bfd_session(sai_object_id_t *bfd_session_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateBfdSessionRequest req =
      convert_create_bfd_session(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateBfdSessionResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = bfd->CreateBfdSession(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *bfd_session_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_bfd_session(sai_object_id_t bfd_session_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveBfdSessionRequest req;
  lemming::dataplane::sai::RemoveBfdSessionResponse resp;
  grpc::ClientContext context;
  req.set_oid(bfd_session_id);

  grpc::Status status = bfd->RemoveBfdSession(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_bfd_session_attribute(sai_object_id_t bfd_session_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetBfdSessionAttributeRequest req;
  lemming::dataplane::sai::SetBfdSessionAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(bfd_session_id);

  switch (attr->id) {
    case SAI_BFD_SESSION_ATTR_VIRTUAL_ROUTER:
      req.set_virtual_router(attr->value.oid);
      break;
    case SAI_BFD_SESSION_ATTR_PORT:
      req.set_port(attr->value.oid);
      break;
    case SAI_BFD_SESSION_ATTR_TC:
      req.set_tc(attr->value.u8);
      break;
    case SAI_BFD_SESSION_ATTR_VLAN_TPID:
      req.set_vlan_tpid(attr->value.u16);
      break;
    case SAI_BFD_SESSION_ATTR_VLAN_PRI:
      req.set_vlan_pri(attr->value.u8);
      break;
    case SAI_BFD_SESSION_ATTR_VLAN_CFI:
      req.set_vlan_cfi(attr->value.u8);
      break;
    case SAI_BFD_SESSION_ATTR_IPHDR_VERSION:
      req.set_iphdr_version(attr->value.u8);
      break;
    case SAI_BFD_SESSION_ATTR_TOS:
      req.set_tos(attr->value.u8);
      break;
    case SAI_BFD_SESSION_ATTR_TTL:
      req.set_ttl(attr->value.u8);
      break;
    case SAI_BFD_SESSION_ATTR_TUNNEL_TOS:
      req.set_tunnel_tos(attr->value.u8);
      break;
    case SAI_BFD_SESSION_ATTR_TUNNEL_TTL:
      req.set_tunnel_ttl(attr->value.u8);
      break;
    case SAI_BFD_SESSION_ATTR_SRC_MAC_ADDRESS:
      req.set_src_mac_address(attr->value.mac, sizeof(attr->value.mac));
      break;
    case SAI_BFD_SESSION_ATTR_DST_MAC_ADDRESS:
      req.set_dst_mac_address(attr->value.mac, sizeof(attr->value.mac));
      break;
    case SAI_BFD_SESSION_ATTR_ECHO_ENABLE:
      req.set_echo_enable(attr->value.booldata);
      break;
    case SAI_BFD_SESSION_ATTR_MIN_TX:
      req.set_min_tx(attr->value.u32);
      break;
    case SAI_BFD_SESSION_ATTR_MIN_RX:
      req.set_min_rx(attr->value.u32);
      break;
    case SAI_BFD_SESSION_ATTR_MULTIPLIER:
      req.set_multiplier(attr->value.u8);
      break;
  }

  grpc::Status status = bfd->SetBfdSessionAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_bfd_session_attribute(sai_object_id_t bfd_session_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetBfdSessionAttributeRequest req;
  lemming::dataplane::sai::GetBfdSessionAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(bfd_session_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::BfdSessionAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = bfd->GetBfdSessionAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BFD_SESSION_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_BFD_SESSION_ATTR_HW_LOOKUP_VALID:
        attr_list[i].value.booldata = resp.attr().hw_lookup_valid();
        break;
      case SAI_BFD_SESSION_ATTR_VIRTUAL_ROUTER:
        attr_list[i].value.oid = resp.attr().virtual_router();
        break;
      case SAI_BFD_SESSION_ATTR_PORT:
        attr_list[i].value.oid = resp.attr().port();
        break;
      case SAI_BFD_SESSION_ATTR_LOCAL_DISCRIMINATOR:
        attr_list[i].value.u32 = resp.attr().local_discriminator();
        break;
      case SAI_BFD_SESSION_ATTR_REMOTE_DISCRIMINATOR:
        attr_list[i].value.u32 = resp.attr().remote_discriminator();
        break;
      case SAI_BFD_SESSION_ATTR_UDP_SRC_PORT:
        attr_list[i].value.u32 = resp.attr().udp_src_port();
        break;
      case SAI_BFD_SESSION_ATTR_TC:
        attr_list[i].value.u8 = resp.attr().tc();
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_TPID:
        attr_list[i].value.u16 = resp.attr().vlan_tpid();
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().vlan_id();
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_PRI:
        attr_list[i].value.u8 = resp.attr().vlan_pri();
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_CFI:
        attr_list[i].value.u8 = resp.attr().vlan_cfi();
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_HEADER_VALID:
        attr_list[i].value.booldata = resp.attr().vlan_header_valid();
        break;
      case SAI_BFD_SESSION_ATTR_BFD_ENCAPSULATION_TYPE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().bfd_encapsulation_type() - 1);
        break;
      case SAI_BFD_SESSION_ATTR_IPHDR_VERSION:
        attr_list[i].value.u8 = resp.attr().iphdr_version();
        break;
      case SAI_BFD_SESSION_ATTR_TOS:
        attr_list[i].value.u8 = resp.attr().tos();
        break;
      case SAI_BFD_SESSION_ATTR_TTL:
        attr_list[i].value.u8 = resp.attr().ttl();
        break;
      case SAI_BFD_SESSION_ATTR_SRC_IP_ADDRESS:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().src_ip_address());
        break;
      case SAI_BFD_SESSION_ATTR_DST_IP_ADDRESS:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().dst_ip_address());
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_TOS:
        attr_list[i].value.u8 = resp.attr().tunnel_tos();
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_TTL:
        attr_list[i].value.u8 = resp.attr().tunnel_ttl();
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_SRC_IP_ADDRESS:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().tunnel_src_ip_address());
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_DST_IP_ADDRESS:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().tunnel_dst_ip_address());
        break;
      case SAI_BFD_SESSION_ATTR_SRC_MAC_ADDRESS:
        memcpy(attr_list[i].value.mac, resp.attr().src_mac_address().data(),
               sizeof(sai_mac_t));
        break;
      case SAI_BFD_SESSION_ATTR_DST_MAC_ADDRESS:
        memcpy(attr_list[i].value.mac, resp.attr().dst_mac_address().data(),
               sizeof(sai_mac_t));
        break;
      case SAI_BFD_SESSION_ATTR_ECHO_ENABLE:
        attr_list[i].value.booldata = resp.attr().echo_enable();
        break;
      case SAI_BFD_SESSION_ATTR_MULTIHOP:
        attr_list[i].value.booldata = resp.attr().multihop();
        break;
      case SAI_BFD_SESSION_ATTR_CBIT:
        attr_list[i].value.booldata = resp.attr().cbit();
        break;
      case SAI_BFD_SESSION_ATTR_MIN_TX:
        attr_list[i].value.u32 = resp.attr().min_tx();
        break;
      case SAI_BFD_SESSION_ATTR_MIN_RX:
        attr_list[i].value.u32 = resp.attr().min_rx();
        break;
      case SAI_BFD_SESSION_ATTR_MULTIPLIER:
        attr_list[i].value.u8 = resp.attr().multiplier();
        break;
      case SAI_BFD_SESSION_ATTR_REMOTE_MIN_TX:
        attr_list[i].value.u32 = resp.attr().remote_min_tx();
        break;
      case SAI_BFD_SESSION_ATTR_REMOTE_MIN_RX:
        attr_list[i].value.u32 = resp.attr().remote_min_rx();
        break;
      case SAI_BFD_SESSION_ATTR_STATE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().state() - 1);
        break;
      case SAI_BFD_SESSION_ATTR_OFFLOAD_TYPE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().offload_type() - 1);
        break;
      case SAI_BFD_SESSION_ATTR_NEGOTIATED_TX:
        attr_list[i].value.u32 = resp.attr().negotiated_tx();
        break;
      case SAI_BFD_SESSION_ATTR_NEGOTIATED_RX:
        attr_list[i].value.u32 = resp.attr().negotiated_rx();
        break;
      case SAI_BFD_SESSION_ATTR_LOCAL_DIAG:
        attr_list[i].value.u8 = resp.attr().local_diag();
        break;
      case SAI_BFD_SESSION_ATTR_REMOTE_DIAG:
        attr_list[i].value.u8 = resp.attr().remote_diag();
        break;
      case SAI_BFD_SESSION_ATTR_REMOTE_MULTIPLIER:
        attr_list[i].value.u8 = resp.attr().remote_multiplier();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_bfd_session_stats(sai_object_id_t bfd_session_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetBfdSessionStatsRequest req;
  lemming::dataplane::sai::GetBfdSessionStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(bfd_session_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(static_cast<lemming::dataplane::sai::BfdSessionStat>(
        counter_ids[i] + 1));
  }
  grpc::Status status = bfd->GetBfdSessionStats(&context, req, &resp);
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

sai_status_t l_get_bfd_session_stats_ext(sai_object_id_t bfd_session_id,
                                         uint32_t number_of_counters,
                                         const sai_stat_id_t *counter_ids,
                                         sai_stats_mode_t mode,
                                         uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_bfd_session_stats(sai_object_id_t bfd_session_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
