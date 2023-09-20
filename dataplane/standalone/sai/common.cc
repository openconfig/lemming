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

#include "dataplane/standalone/proto/common.pb.h"

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
