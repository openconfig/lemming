#include "dataplane/standalone/sai/common.h"

#include <glog/logging.h>

#include "common.h"
#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/sai/entry.h"

std::string convert_from_ip_addr(sai_ip_addr_family_t addr_family,
                                 const sai_ip_addr_t& addr) {
  if (addr_family == SAI_IP_ADDR_FAMILY_IPV4) {
    sai_ip4_t ip = addr.ip4;
    return reinterpret_cast<char*>(&ip);
  }
  sai_ip6_t ip;
  std::copy(addr.ip6, addr.ip6 + sizeof(sai_ip6_t), ip);
  return reinterpret_cast<char*>(ip);
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
