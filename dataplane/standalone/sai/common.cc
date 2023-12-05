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

#include "dataplane/standalone/sai/common.h"

#include <glog/logging.h>

#include <algorithm>
#include <string>

#include "dataplane/proto/common.pb.h"

std::string convert_from_ip_addr(sai_ip_addr_family_t addr_family,
                                 const sai_ip_addr_t& addr) {
  if (addr_family == SAI_IP_ADDR_FAMILY_IPV4) {
    return std::string(
        reinterpret_cast<const char*>(&addr.ip4),
        reinterpret_cast<const char*>(&addr.ip4) + sizeof(sai_ip4_t));
  }
  return std::string(addr.ip6, addr.ip6 + 16);
}

std::string convert_from_ip_address(const sai_ip_address_t& val) {
  return convert_from_ip_addr(val.addr_family, val.addr);
}

lemming::dataplane::sai::RouteEntry convert_from_route_entry(
    const sai_route_entry_t& entry) {
  lemming::dataplane::sai::RouteEntry re;
  re.set_switch_id(entry.switch_id);
  re.set_vr_id(entry.vr_id);
  *re.mutable_destination() = convert_from_ip_prefix(entry.destination);
  return re;
}

lemming::dataplane::sai::IpPrefix convert_from_ip_prefix(
    const sai_ip_prefix_t& ip_prefix) {
  lemming::dataplane::sai::IpPrefix ip;
  ip.set_addr(convert_from_ip_addr(ip_prefix.addr_family, ip_prefix.addr));
  ip.set_mask(convert_from_ip_addr(ip_prefix.addr_family, ip_prefix.mask));
  return ip;
}

sai_ip_addr_t convert_to_ip_addr(std::string val) {
  sai_ip_addr_t addr;
  if (val.length() == 4) {
    addr.ip4 = *reinterpret_cast<sai_uint32_t*>(&val[0]);
  } else if (val.length() == 16) {
    memcpy(addr.ip6, val.data(), sizeof(sai_ip6_t));
  }
  return addr;
}

sai_ip_address_t convert_to_ip_address(std::string str) {
  sai_ip_address_t ip;
  ip.addr = convert_to_ip_addr(str);
  if (str.length() == 4) {
    ip.addr_family = SAI_IP_ADDR_FAMILY_IPV4;
  } else if (str.length() == 16) {
    ip.addr_family = SAI_IP_ADDR_FAMILY_IPV6;
  }
  return ip;
}

sai_route_entry_t convert_to_route_entry(
    const lemming::dataplane::sai::RouteEntry& entry) {
  sai_route_entry_t re;
  re.switch_id = entry.switch_id();
  re.vr_id = entry.vr_id();
  re.destination = convert_to_ip_prefix(entry.destination());
  return re;
}

sai_ip_prefix_t convert_to_ip_prefix(
    const lemming::dataplane::sai::IpPrefix& ip_prefix) {
  sai_ip_prefix_t ip;
  ip.addr = convert_to_ip_addr(ip_prefix.addr());
  ip.mask = convert_to_ip_addr(ip_prefix.mask());
  return ip;
}

std::vector<sai_port_oper_status_notification_t> convert_to_oper_status(
    const lemming::dataplane::sai::PortStateChangeNotificationResponse& resp) {
  std::vector<sai_port_oper_status_notification_t> list;
  for (auto d : resp.data()) {
    list.push_back({
        .port_id = d.port_id(),
        .port_state = static_cast<sai_port_oper_status_t>(d.port_state() - 1),
    });
  }
  return list;
}

lemming::dataplane::sai::NeighborEntry convert_from_neighbor_entry(
    const sai_neighbor_entry_t& entry) {
  lemming::dataplane::sai::NeighborEntry ne;
  ne.set_switch_id(entry.switch_id);
  ne.set_rif_id(entry.rif_id);
  ne.set_ip_address(convert_from_ip_address(entry.ip_address));

  return ne;
}

sai_neighbor_entry_t convert_to_neighbor_entry(
    const lemming::dataplane::sai::NeighborEntry& entry) {
  sai_neighbor_entry_t ne;
  ne.switch_id = entry.switch_id();
  ne.rif_id = entry.rif_id();
  ne.ip_address = convert_to_ip_address(entry.ip_address());

  return ne;
}

void convert_to_acl_capability(
    sai_acl_capability_t& out,
    const lemming::dataplane::sai::ACLCapability& in) {
  out.is_action_list_mandatory = in.is_action_list_mandatory();
  for (int i = 0; i < in.action_list().size(); i++) {
    out.action_list.list[0] = in.action_list(i) - 1;
  }
  out.action_list.count = in.action_list().size();
}

lemming::dataplane::sai::AclActionData convert_from_acl_action_data(
    const sai_acl_action_data_t& in, sai_object_id_t id) {
  lemming::dataplane::sai::AclActionData out;
  out.set_enable(in.enable);
  out.set_oid(id);
  return out;
}

lemming::dataplane::sai::AclActionData convert_from_acl_action_data_action(
    const sai_acl_action_data_t& in, sai_int32_t val) {
  lemming::dataplane::sai::AclActionData out;
  out.set_enable(in.enable);
  out.set_packet_action(
      static_cast<lemming::dataplane::sai::PacketAction>(val + 1));
  return out;
}

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data(
    const sai_acl_field_data_t& in, sai_ip4_t data, sai_ip4_t mask) {
  lemming::dataplane::sai::AclFieldData out;
  out.set_enable(in.enable);
  *out.mutable_data_ip() =
      std::string(reinterpret_cast<const char*>(&data), sizeof(sai_ip4_t));
  *out.mutable_mask_ip() =
      std::string(reinterpret_cast<const char*>(&mask), sizeof(sai_ip4_t));
  return out;
}

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data(
    const sai_acl_field_data_t& in, sai_uint8_t data, sai_uint8_t mask) {
  lemming::dataplane::sai::AclFieldData out;
  out.set_enable(in.enable);
  out.set_data_uint(data);
  out.set_mask_uint(mask);
  return out;
}

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data(
    const sai_acl_field_data_t& in, sai_uint16_t data, sai_uint16_t mask) {
  lemming::dataplane::sai::AclFieldData out;
  out.set_enable(in.enable);
  out.set_data_uint(data);
  out.set_mask_uint(mask);
  return out;
}

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data(
    const sai_acl_field_data_t& in, sai_object_id_t data) {
  lemming::dataplane::sai::AclFieldData out;
  out.set_enable(in.enable);
  out.set_data_oid(data);
  return out;
}

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data_ip6(
    const sai_acl_field_data_t& in, const sai_ip6_t data,
    const sai_ip6_t mask) {
  lemming::dataplane::sai::AclFieldData out;
  out.set_enable(in.enable);
  *out.mutable_data_ip() =
      std::string(reinterpret_cast<const char*>(data), sizeof(sai_ip6_t));
  *out.mutable_mask_ip() =
      std::string(reinterpret_cast<const char*>(mask), sizeof(sai_ip6_t));
  return out;
}

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data_mac(
    const sai_acl_field_data_t& in, const sai_mac_t data,
    const sai_mac_t mask) {
  lemming::dataplane::sai::AclFieldData out;
  out.set_enable(in.enable);
  *out.mutable_data_mac() =
      std::string(reinterpret_cast<const char*>(data), sizeof(sai_mac_t));
  *out.mutable_mask_mac() =
      std::string(reinterpret_cast<const char*>(mask), sizeof(sai_mac_t));
  return out;
}

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data_ip_type(
    const sai_acl_field_data_t& in, sai_int32_t type, sai_int32_t mask) {
  lemming::dataplane::sai::AclFieldData out;
  out.set_enable(in.enable);
  out.set_data_ip_type(
      static_cast<lemming::dataplane::sai::AclIpType>(type + 1));
  out.set_mask_int(mask);
  return out;
}
