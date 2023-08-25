

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

#include "dataplane/standalone/sai/mirror.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/mirror.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_mirror_api_t l_mirror = {
    .create_mirror_session = l_create_mirror_session,
    .remove_mirror_session = l_remove_mirror_session,
    .set_mirror_session_attribute = l_set_mirror_session_attribute,
    .get_mirror_session_attribute = l_get_mirror_session_attribute,
};

sai_status_t l_create_mirror_session(sai_object_id_t *mirror_session_id,
                                     sai_object_id_t switch_id,
                                     uint32_t attr_count,
                                     const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateMirrorSessionRequest req;
  lemming::dataplane::sai::CreateMirrorSessionResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MIRROR_SESSION_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::MirrorSessionType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_MIRROR_SESSION_ATTR_MONITOR_PORT:
        req.set_monitor_port(attr_list[i].value.oid);
        break;
      case SAI_MIRROR_SESSION_ATTR_TRUNCATE_SIZE:
        req.set_truncate_size(attr_list[i].value.u16);
        break;
      case SAI_MIRROR_SESSION_ATTR_SAMPLE_RATE:
        req.set_sample_rate(attr_list[i].value.u32);
        break;
      case SAI_MIRROR_SESSION_ATTR_CONGESTION_MODE:
        req.set_congestion_mode(
            static_cast<lemming::dataplane::sai::MirrorSessionCongestionMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_MIRROR_SESSION_ATTR_TC:
        req.set_tc(attr_list[i].value.u8);
        break;
      case SAI_MIRROR_SESSION_ATTR_VLAN_TPID:
        req.set_vlan_tpid(attr_list[i].value.u16);
        break;
      case SAI_MIRROR_SESSION_ATTR_VLAN_ID:
        req.set_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_MIRROR_SESSION_ATTR_VLAN_PRI:
        req.set_vlan_pri(attr_list[i].value.u8);
        break;
      case SAI_MIRROR_SESSION_ATTR_VLAN_CFI:
        req.set_vlan_cfi(attr_list[i].value.u8);
        break;
      case SAI_MIRROR_SESSION_ATTR_VLAN_HEADER_VALID:
        req.set_vlan_header_valid(attr_list[i].value.booldata);
        break;
      case SAI_MIRROR_SESSION_ATTR_ERSPAN_ENCAPSULATION_TYPE:
        req.set_erspan_encapsulation_type(
            static_cast<lemming::dataplane::sai::ErspanEncapsulationType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_MIRROR_SESSION_ATTR_IPHDR_VERSION:
        req.set_iphdr_version(attr_list[i].value.u8);
        break;
      case SAI_MIRROR_SESSION_ATTR_TOS:
        req.set_tos(attr_list[i].value.u8);
        break;
      case SAI_MIRROR_SESSION_ATTR_TTL:
        req.set_ttl(attr_list[i].value.u8);
        break;
      case SAI_MIRROR_SESSION_ATTR_SRC_IP_ADDRESS:
        req.set_src_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_MIRROR_SESSION_ATTR_DST_IP_ADDRESS:
        req.set_dst_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_MIRROR_SESSION_ATTR_SRC_MAC_ADDRESS:
        req.set_src_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_MIRROR_SESSION_ATTR_DST_MAC_ADDRESS:
        req.set_dst_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_MIRROR_SESSION_ATTR_GRE_PROTOCOL_TYPE:
        req.set_gre_protocol_type(attr_list[i].value.u16);
        break;
      case SAI_MIRROR_SESSION_ATTR_MONITOR_PORTLIST_VALID:
        req.set_monitor_portlist_valid(attr_list[i].value.booldata);
        break;
      case SAI_MIRROR_SESSION_ATTR_MONITOR_PORTLIST:
        req.mutable_monitor_portlist()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_MIRROR_SESSION_ATTR_POLICER:
        req.set_policer(attr_list[i].value.oid);
        break;
      case SAI_MIRROR_SESSION_ATTR_UDP_SRC_PORT:
        req.set_udp_src_port(attr_list[i].value.u16);
        break;
      case SAI_MIRROR_SESSION_ATTR_UDP_DST_PORT:
        req.set_udp_dst_port(attr_list[i].value.u16);
        break;
    }
  }
  grpc::Status status = mirror->CreateMirrorSession(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *mirror_session_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_MIRROR_SESSION, mirror_session_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_mirror_session(sai_object_id_t mirror_session_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_MIRROR_SESSION, mirror_session_id);
}

sai_status_t l_set_mirror_session_attribute(sai_object_id_t mirror_session_id,
                                            const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_MIRROR_SESSION,
                                   mirror_session_id, attr);
}

sai_status_t l_get_mirror_session_attribute(sai_object_id_t mirror_session_id,
                                            uint32_t attr_count,
                                            sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_MIRROR_SESSION,
                                   mirror_session_id, attr_count, attr_list);
}
