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

#include "dataplane/standalone/route.h"

#include <glog/logging.h>

#include <bitset>
#include <string>
#include <vector>

#include "dataplane/standalone/translator.h"

sai_status_t Route::create(common_entry_t entry, _In_ uint32_t attr_count,
                           _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);
  sai_packet_action_t act;
  sai_object_id_t oid;
  for (auto attr : attrs) {
    switch (attr.id) {
      case SAI_ROUTE_ENTRY_ATTR_PACKET_ACTION:
        act = static_cast<sai_packet_action_t>(attr.value.s32);
        break;
      case SAI_ROUTE_ENTRY_ATTR_NEXT_HOP_ID:
        oid = attr.value.oid;
        break;
    }
  }
  lemming::dataplane::AddIPRouteRequest req;
  switch (entry.route_entry->destination.addr_family) {
    case SAI_IP_ADDR_FAMILY_IPV4:
      req.mutable_route()->mutable_prefix()->mutable_mask()->set_addr(
          &entry.route_entry->destination.addr.ip4,
          sizeof(entry.route_entry->destination.addr.ip4));
      req.mutable_route()->mutable_prefix()->mutable_mask()->set_mask(
          &entry.route_entry->destination.mask.ip4,
          sizeof(entry.route_entry->destination.mask.ip4));
      break;
    case SAI_IP_ADDR_FAMILY_IPV6:
      req.mutable_route()->mutable_prefix()->mutable_mask()->set_addr(
          &entry.route_entry->destination.addr.ip6,
          sizeof(entry.route_entry->destination.addr.ip6));
      req.mutable_route()->mutable_prefix()->mutable_mask()->set_mask(
          &entry.route_entry->destination.mask.ip6,
          sizeof(entry.route_entry->destination.mask.ip6));
      break;
    default:
      return SAI_STATUS_INVALID_PARAMETER;
  }
  req.mutable_route()->mutable_prefix()->set_vrf_id(entry.route_entry->vr_id);

  // TODO(dgrau): Implement CPU actions.

  switch (act) {
    case SAI_PACKET_ACTION_DROP:
    case SAI_PACKET_ACTION_TRAP:  // COPY and DROP
    case SAI_PACKET_ACTION_DENY:  // COPY_CANCEL and DROP
      req.mutable_route()->set_action(lemming::dataplane::PACKET_ACTION_DROP);
      break;
    case SAI_PACKET_ACTION_FORWARD:
      req.mutable_route()->set_action(
          lemming::dataplane::PACKET_ACTION_FORWARD);
    case SAI_PACKET_ACTION_LOG:      // COPY and FORWARD
    case SAI_PACKET_ACTION_TRANSIT:  // COPY_CANCEL and FORWARD
      req.mutable_route()->set_action(
          lemming::dataplane::PACKET_ACTION_FORWARD);
      break;
    case SAI_PACKET_ACTION_COPY:
      break;
    case SAI_PACKET_ACTION_COPY_CANCEL:
      break;
  }

  std::string hop_id = std::to_string(oid);
  sai_object_type_t obj_type = this->attrMgr->get_type(hop_id);

  // If the packet action is drop, then next hop is optional.
  if (req.route().action() == lemming::dataplane::PACKET_ACTION_FORWARD) {
    switch (obj_type) {
      case SAI_OBJECT_TYPE_NEXT_HOP:
        req.mutable_route()->set_next_hop_id(oid);
        break;
      case SAI_OBJECT_TYPE_NEXT_HOP_GROUP:
        req.mutable_route()->set_next_hop_group_id(oid);
        break;
      case SAI_OBJECT_TYPE_ROUTER_INTERFACE:
        req.mutable_route()->set_interface_id(hop_id);
        break;
      case SAI_OBJECT_TYPE_PORT:
        req.mutable_route()->set_port_id(hop_id);
        break;
      default:
        return SAI_STATUS_INVALID_OBJECT_TYPE;
    }
  }

  LOG(INFO) << "create route dest " << this->id << ", action " << act
            << ", next hop id " << oid << ", type " << obj_type;

  grpc::ClientContext context;
  lemming::dataplane::AddIPRouteResponse resp;
  auto status = this->dataplane->AddIPRoute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << "Failed to create route: " << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  APIBase::create(attrs.size(), attrs.data());
  return SAI_STATUS_SUCCESS;
}

sai_status_t Route::set_attribute(_In_ const sai_attribute_t* attr) {
  return SAI_STATUS_SUCCESS;
}
