

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

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_mirror_session(sai_object_id_t mirror_session_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveMirrorSessionRequest req;
  lemming::dataplane::sai::RemoveMirrorSessionResponse resp;
  grpc::ClientContext context;
  req.set_oid(mirror_session_id);

  grpc::Status status = mirror->RemoveMirrorSession(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_mirror_session_attribute(sai_object_id_t mirror_session_id,
                                            const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetMirrorSessionAttributeRequest req;
  lemming::dataplane::sai::SetMirrorSessionAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(mirror_session_id);

  switch (attr->id) {
    case SAI_MIRROR_SESSION_ATTR_MONITOR_PORT:
      req.set_monitor_port(attr->value.oid);
      break;
    case SAI_MIRROR_SESSION_ATTR_TRUNCATE_SIZE:
      req.set_truncate_size(attr->value.u16);
      break;
    case SAI_MIRROR_SESSION_ATTR_SAMPLE_RATE:
      req.set_sample_rate(attr->value.u32);
      break;
    case SAI_MIRROR_SESSION_ATTR_CONGESTION_MODE:
      req.set_congestion_mode(
          static_cast<lemming::dataplane::sai::MirrorSessionCongestionMode>(
              attr->value.s32 + 1));
      break;
    case SAI_MIRROR_SESSION_ATTR_TC:
      req.set_tc(attr->value.u8);
      break;
    case SAI_MIRROR_SESSION_ATTR_VLAN_TPID:
      req.set_vlan_tpid(attr->value.u16);
      break;
    case SAI_MIRROR_SESSION_ATTR_VLAN_ID:
      req.set_vlan_id(attr->value.u16);
      break;
    case SAI_MIRROR_SESSION_ATTR_VLAN_PRI:
      req.set_vlan_pri(attr->value.u8);
      break;
    case SAI_MIRROR_SESSION_ATTR_VLAN_CFI:
      req.set_vlan_cfi(attr->value.u8);
      break;
    case SAI_MIRROR_SESSION_ATTR_VLAN_HEADER_VALID:
      req.set_vlan_header_valid(attr->value.booldata);
      break;
    case SAI_MIRROR_SESSION_ATTR_IPHDR_VERSION:
      req.set_iphdr_version(attr->value.u8);
      break;
    case SAI_MIRROR_SESSION_ATTR_TOS:
      req.set_tos(attr->value.u8);
      break;
    case SAI_MIRROR_SESSION_ATTR_TTL:
      req.set_ttl(attr->value.u8);
      break;
    case SAI_MIRROR_SESSION_ATTR_SRC_IP_ADDRESS:
      req.set_src_ip_address(convert_from_ip_address(attr->value.ipaddr));
      break;
    case SAI_MIRROR_SESSION_ATTR_DST_IP_ADDRESS:
      req.set_dst_ip_address(convert_from_ip_address(attr->value.ipaddr));
      break;
    case SAI_MIRROR_SESSION_ATTR_SRC_MAC_ADDRESS:
      req.set_src_mac_address(attr->value.mac, sizeof(attr->value.mac));
      break;
    case SAI_MIRROR_SESSION_ATTR_DST_MAC_ADDRESS:
      req.set_dst_mac_address(attr->value.mac, sizeof(attr->value.mac));
      break;
    case SAI_MIRROR_SESSION_ATTR_GRE_PROTOCOL_TYPE:
      req.set_gre_protocol_type(attr->value.u16);
      break;
    case SAI_MIRROR_SESSION_ATTR_MONITOR_PORTLIST:
      req.mutable_monitor_portlist()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_MIRROR_SESSION_ATTR_POLICER:
      req.set_policer(attr->value.oid);
      break;
    case SAI_MIRROR_SESSION_ATTR_UDP_SRC_PORT:
      req.set_udp_src_port(attr->value.u16);
      break;
    case SAI_MIRROR_SESSION_ATTR_UDP_DST_PORT:
      req.set_udp_dst_port(attr->value.u16);
      break;
  }

  grpc::Status status = mirror->SetMirrorSessionAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_mirror_session_attribute(sai_object_id_t mirror_session_id,
                                            uint32_t attr_count,
                                            sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMirrorSessionAttributeRequest req;
  lemming::dataplane::sai::GetMirrorSessionAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(mirror_session_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::MirrorSessionAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = mirror->GetMirrorSessionAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MIRROR_SESSION_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_MIRROR_SESSION_ATTR_MONITOR_PORT:
        attr_list[i].value.oid = resp.attr().monitor_port();
        break;
      case SAI_MIRROR_SESSION_ATTR_TRUNCATE_SIZE:
        attr_list[i].value.u16 = resp.attr().truncate_size();
        break;
      case SAI_MIRROR_SESSION_ATTR_SAMPLE_RATE:
        attr_list[i].value.u32 = resp.attr().sample_rate();
        break;
      case SAI_MIRROR_SESSION_ATTR_CONGESTION_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().congestion_mode() - 1);
        break;
      case SAI_MIRROR_SESSION_ATTR_TC:
        attr_list[i].value.u8 = resp.attr().tc();
        break;
      case SAI_MIRROR_SESSION_ATTR_VLAN_TPID:
        attr_list[i].value.u16 = resp.attr().vlan_tpid();
        break;
      case SAI_MIRROR_SESSION_ATTR_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().vlan_id();
        break;
      case SAI_MIRROR_SESSION_ATTR_VLAN_PRI:
        attr_list[i].value.u8 = resp.attr().vlan_pri();
        break;
      case SAI_MIRROR_SESSION_ATTR_VLAN_CFI:
        attr_list[i].value.u8 = resp.attr().vlan_cfi();
        break;
      case SAI_MIRROR_SESSION_ATTR_VLAN_HEADER_VALID:
        attr_list[i].value.booldata = resp.attr().vlan_header_valid();
        break;
      case SAI_MIRROR_SESSION_ATTR_ERSPAN_ENCAPSULATION_TYPE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().erspan_encapsulation_type() - 1);
        break;
      case SAI_MIRROR_SESSION_ATTR_IPHDR_VERSION:
        attr_list[i].value.u8 = resp.attr().iphdr_version();
        break;
      case SAI_MIRROR_SESSION_ATTR_TOS:
        attr_list[i].value.u8 = resp.attr().tos();
        break;
      case SAI_MIRROR_SESSION_ATTR_TTL:
        attr_list[i].value.u8 = resp.attr().ttl();
        break;
      case SAI_MIRROR_SESSION_ATTR_SRC_IP_ADDRESS:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().src_ip_address());
        break;
      case SAI_MIRROR_SESSION_ATTR_DST_IP_ADDRESS:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().dst_ip_address());
        break;
      case SAI_MIRROR_SESSION_ATTR_SRC_MAC_ADDRESS:
        memcpy(attr_list[i].value.mac, resp.attr().src_mac_address().data(),
               sizeof(sai_mac_t));
        break;
      case SAI_MIRROR_SESSION_ATTR_DST_MAC_ADDRESS:
        memcpy(attr_list[i].value.mac, resp.attr().dst_mac_address().data(),
               sizeof(sai_mac_t));
        break;
      case SAI_MIRROR_SESSION_ATTR_GRE_PROTOCOL_TYPE:
        attr_list[i].value.u16 = resp.attr().gre_protocol_type();
        break;
      case SAI_MIRROR_SESSION_ATTR_MONITOR_PORTLIST_VALID:
        attr_list[i].value.booldata = resp.attr().monitor_portlist_valid();
        break;
      case SAI_MIRROR_SESSION_ATTR_MONITOR_PORTLIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().monitor_portlist(),
                  attr_list[i].value.objlist.count);
        break;
      case SAI_MIRROR_SESSION_ATTR_POLICER:
        attr_list[i].value.oid = resp.attr().policer();
        break;
      case SAI_MIRROR_SESSION_ATTR_UDP_SRC_PORT:
        attr_list[i].value.u16 = resp.attr().udp_src_port();
        break;
      case SAI_MIRROR_SESSION_ATTR_UDP_DST_PORT:
        attr_list[i].value.u16 = resp.attr().udp_dst_port();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
