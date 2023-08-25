

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

#include "dataplane/standalone/proto/bfd.pb.h"
#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_bfd_api_t l_bfd = {
    .create_bfd_session = l_create_bfd_session,
    .remove_bfd_session = l_remove_bfd_session,
    .set_bfd_session_attribute = l_set_bfd_session_attribute,
    .get_bfd_session_attribute = l_get_bfd_session_attribute,
    .get_bfd_session_stats = l_get_bfd_session_stats,
    .get_bfd_session_stats_ext = l_get_bfd_session_stats_ext,
    .clear_bfd_session_stats = l_clear_bfd_session_stats,
};

sai_status_t l_create_bfd_session(sai_object_id_t *bfd_session_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateBfdSessionRequest req;
  lemming::dataplane::sai::CreateBfdSessionResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BFD_SESSION_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::BfdSessionType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_BFD_SESSION_ATTR_HW_LOOKUP_VALID:
        req.set_hw_lookup_valid(attr_list[i].value.booldata);
        break;
      case SAI_BFD_SESSION_ATTR_VIRTUAL_ROUTER:
        req.set_virtual_router(attr_list[i].value.oid);
        break;
      case SAI_BFD_SESSION_ATTR_PORT:
        req.set_port(attr_list[i].value.oid);
        break;
      case SAI_BFD_SESSION_ATTR_LOCAL_DISCRIMINATOR:
        req.set_local_discriminator(attr_list[i].value.u32);
        break;
      case SAI_BFD_SESSION_ATTR_REMOTE_DISCRIMINATOR:
        req.set_remote_discriminator(attr_list[i].value.u32);
        break;
      case SAI_BFD_SESSION_ATTR_UDP_SRC_PORT:
        req.set_udp_src_port(attr_list[i].value.u32);
        break;
      case SAI_BFD_SESSION_ATTR_TC:
        req.set_tc(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_TPID:
        req.set_vlan_tpid(attr_list[i].value.u16);
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_ID:
        req.set_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_PRI:
        req.set_vlan_pri(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_CFI:
        req.set_vlan_cfi(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_VLAN_HEADER_VALID:
        req.set_vlan_header_valid(attr_list[i].value.booldata);
        break;
      case SAI_BFD_SESSION_ATTR_BFD_ENCAPSULATION_TYPE:
        req.set_bfd_encapsulation_type(
            static_cast<lemming::dataplane::sai::BfdEncapsulationType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BFD_SESSION_ATTR_IPHDR_VERSION:
        req.set_iphdr_version(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_TOS:
        req.set_tos(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_TTL:
        req.set_ttl(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_SRC_IP_ADDRESS:
        req.set_src_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_BFD_SESSION_ATTR_DST_IP_ADDRESS:
        req.set_dst_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_TOS:
        req.set_tunnel_tos(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_TTL:
        req.set_tunnel_ttl(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_SRC_IP_ADDRESS:
        req.set_tunnel_src_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_BFD_SESSION_ATTR_TUNNEL_DST_IP_ADDRESS:
        req.set_tunnel_dst_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_BFD_SESSION_ATTR_SRC_MAC_ADDRESS:
        req.set_src_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_BFD_SESSION_ATTR_DST_MAC_ADDRESS:
        req.set_dst_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_BFD_SESSION_ATTR_ECHO_ENABLE:
        req.set_echo_enable(attr_list[i].value.booldata);
        break;
      case SAI_BFD_SESSION_ATTR_MULTIHOP:
        req.set_multihop(attr_list[i].value.booldata);
        break;
      case SAI_BFD_SESSION_ATTR_CBIT:
        req.set_cbit(attr_list[i].value.booldata);
        break;
      case SAI_BFD_SESSION_ATTR_MIN_TX:
        req.set_min_tx(attr_list[i].value.u32);
        break;
      case SAI_BFD_SESSION_ATTR_MIN_RX:
        req.set_min_rx(attr_list[i].value.u32);
        break;
      case SAI_BFD_SESSION_ATTR_MULTIPLIER:
        req.set_multiplier(attr_list[i].value.u8);
        break;
      case SAI_BFD_SESSION_ATTR_OFFLOAD_TYPE:
        req.set_offload_type(
            static_cast<lemming::dataplane::sai::BfdSessionOffloadType>(
                attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = bfd->CreateBfdSession(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *bfd_session_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_BFD_SESSION, bfd_session_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_bfd_session(sai_object_id_t bfd_session_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_BFD_SESSION, bfd_session_id);
}

sai_status_t l_set_bfd_session_attribute(sai_object_id_t bfd_session_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_BFD_SESSION, bfd_session_id,
                                   attr);
}

sai_status_t l_get_bfd_session_attribute(sai_object_id_t bfd_session_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_BFD_SESSION, bfd_session_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_bfd_session_stats(sai_object_id_t bfd_session_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_BFD_SESSION, bfd_session_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_bfd_session_stats_ext(sai_object_id_t bfd_session_id,
                                         uint32_t number_of_counters,
                                         const sai_stat_id_t *counter_ids,
                                         sai_stats_mode_t mode,
                                         uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_BFD_SESSION, bfd_session_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_bfd_session_stats(sai_object_id_t bfd_session_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_BFD_SESSION, bfd_session_id,
                                 number_of_counters, counter_ids);
}
