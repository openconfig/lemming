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

#include "dataplane/standalone/neighbor.h"

#include <glog/logging.h>

#include <bitset>
#include <string>
#include <vector>

#include "dataplane/standalone/translator.h"

sai_status_t Neighbor::create(common_entry_t entry, _In_ uint32_t attr_count,
                              _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);
  sai_mac_t* mac;
  for (auto attr : attrs) {
    switch (attr.id) {
      case SAI_NEIGHBOR_ENTRY_ATTR_DST_MAC_ADDRESS:
        mac = &attr.value.mac;
        break;
    }
  }
  lemming::dataplane::AddNeighborRequest req;
  req.set_mac(mac, sizeof(sai_mac_t));
  req.set_port_id(std::to_string(entry.neighbor_entry->rif_id));

  switch (entry.neighbor_entry->ip_address.addr_family) {
    case SAI_IP_ADDR_FAMILY_IPV4:
      req.set_ip_bytes(&entry.neighbor_entry->ip_address.addr.ip4,
                       sizeof(entry.neighbor_entry->ip_address.addr.ip4));
      break;
    case SAI_IP_ADDR_FAMILY_IPV6:
      req.set_ip_bytes(&entry.neighbor_entry->ip_address.addr.ip6,
                       sizeof(entry.neighbor_entry->ip_address.addr.ip6));
      break;
    default:
      return SAI_STATUS_INVALID_PARAMETER;
  }
  grpc::ClientContext context;
  lemming::dataplane::AddNeighborResponse resp;
  auto status = this->dataplane->AddNeighbor(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << "Failed to add neighbor: " << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  APIBase::create(attrs.size(), attrs.data());
  return SAI_STATUS_SUCCESS;
}

sai_status_t Neighbor::set_attribute(_In_ const sai_attribute_t* attr) {
  return SAI_STATUS_SUCCESS;
}
