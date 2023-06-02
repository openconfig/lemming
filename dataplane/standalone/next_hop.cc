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

#include <glog/logging.h>

#include <bitset>
#include <string>
#include <vector>

#include "dataplane/standalone/next_hop.h"
#include "dataplane/standalone/translator.h"

sai_status_t NextHop::create(_In_ uint32_t attr_count,
                             _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);
  sai_next_hop_type_t type;
  sai_ip_address_t ip;
  sai_object_id_t oid;
  for (auto attr : attrs) {
    switch (attr.id) {
      case SAI_NEXT_HOP_ATTR_TYPE:
        type = static_cast<sai_next_hop_type_t>(attr.value.s32);
        break;
      case SAI_NEXT_HOP_ATTR_IP:
        ip = attr.value.ipaddr;
        break;
      case SAI_NEXT_HOP_ATTR_ROUTER_INTERFACE_ID:
        oid = attr.value.oid;
        break;
    }
  }
  if (type != SAI_NEXT_HOP_TYPE_IP) {
    return SAI_STATUS_NOT_SUPPORTED;
  }
  grpc::ClientContext context;
  lemming::dataplane::AddNextHopRequest req;
  lemming::dataplane::AddNextHopResponse resp;
  req.set_id(std::stoul(this->id));
  req.mutable_next_hop()->set_port(std::to_string(oid));

  switch (ip.addr_family) {
    case SAI_IP_ADDR_FAMILY_IPV4:
      req.mutable_next_hop()->set_ip_bytes(&ip.addr.ip4, sizeof(ip.addr.ip4));
      break;
    case SAI_IP_ADDR_FAMILY_IPV6:
      req.mutable_next_hop()->set_ip_bytes(ip.addr.ip6, sizeof(ip.addr.ip6));
      break;
    default:
      return SAI_STATUS_INVALID_PARAMETER;
  }

  auto status = this->dataplane->AddNextHop(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << "Failed to create route: " << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  APIBase::create(attrs.size(), attrs.data());
  return SAI_STATUS_SUCCESS;
}

sai_status_t NextHop::set_attribute(_In_ const sai_attribute_t* attr) {
  return SAI_STATUS_SUCCESS;
}