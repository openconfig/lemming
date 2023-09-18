

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

#include "dataplane/standalone/sai/neighbor.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/neighbor.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_neighbor_api_t l_neighbor = {
    .create_neighbor_entry = l_create_neighbor_entry,
    .remove_neighbor_entry = l_remove_neighbor_entry,
    .set_neighbor_entry_attribute = l_set_neighbor_entry_attribute,
    .get_neighbor_entry_attribute = l_get_neighbor_entry_attribute,
    .remove_all_neighbor_entries = l_remove_all_neighbor_entries,
};

sai_status_t l_create_neighbor_entry(const sai_neighbor_entry_t *neighbor_entry,
                                     uint32_t attr_count,
                                     const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateNeighborEntryRequest req;
  lemming::dataplane::sai::CreateNeighborEntryResponse resp;
  grpc::ClientContext context;

  *req.mutable_entry() = convert_from_neighbor_entry(*neighbor_entry);
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_NEIGHBOR_ENTRY_ATTR_DST_MAC_ADDRESS:
        req.set_dst_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_PACKET_ACTION:
        req.set_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_USER_TRAP_ID:
        req.set_user_trap_id(attr_list[i].value.oid);
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_NO_HOST_ROUTE:
        req.set_no_host_route(attr_list[i].value.booldata);
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_META_DATA:
        req.set_meta_data(attr_list[i].value.u32);
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_COUNTER_ID:
        req.set_counter_id(attr_list[i].value.oid);
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_ENCAP_INDEX:
        req.set_encap_index(attr_list[i].value.u32);
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_ENCAP_IMPOSE_INDEX:
        req.set_encap_impose_index(attr_list[i].value.booldata);
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_IS_LOCAL:
        req.set_is_local(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status = neighbor->CreateNeighborEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_neighbor_entry(
    const sai_neighbor_entry_t *neighbor_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveNeighborEntryRequest req;
  lemming::dataplane::sai::RemoveNeighborEntryResponse resp;
  grpc::ClientContext context;

  *req.mutable_entry() = convert_from_neighbor_entry(*neighbor_entry);
  grpc::Status status = neighbor->RemoveNeighborEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_neighbor_entry_attribute(
    const sai_neighbor_entry_t *neighbor_entry, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetNeighborEntryAttributeRequest req;
  lemming::dataplane::sai::SetNeighborEntryAttributeResponse resp;
  grpc::ClientContext context;

  *req.mutable_entry() = convert_from_neighbor_entry(*neighbor_entry);

  switch (attr->id) {
    case SAI_NEIGHBOR_ENTRY_ATTR_DST_MAC_ADDRESS:
      req.set_dst_mac_address(attr->value.mac, sizeof(attr->value.mac));
      break;
    case SAI_NEIGHBOR_ENTRY_ATTR_PACKET_ACTION:
      req.set_packet_action(static_cast<lemming::dataplane::sai::PacketAction>(
          attr->value.s32 + 1));
      break;
    case SAI_NEIGHBOR_ENTRY_ATTR_USER_TRAP_ID:
      req.set_user_trap_id(attr->value.oid);
      break;
    case SAI_NEIGHBOR_ENTRY_ATTR_NO_HOST_ROUTE:
      req.set_no_host_route(attr->value.booldata);
      break;
    case SAI_NEIGHBOR_ENTRY_ATTR_META_DATA:
      req.set_meta_data(attr->value.u32);
      break;
    case SAI_NEIGHBOR_ENTRY_ATTR_COUNTER_ID:
      req.set_counter_id(attr->value.oid);
      break;
    case SAI_NEIGHBOR_ENTRY_ATTR_ENCAP_INDEX:
      req.set_encap_index(attr->value.u32);
      break;
    case SAI_NEIGHBOR_ENTRY_ATTR_ENCAP_IMPOSE_INDEX:
      req.set_encap_impose_index(attr->value.booldata);
      break;
    case SAI_NEIGHBOR_ENTRY_ATTR_IS_LOCAL:
      req.set_is_local(attr->value.booldata);
      break;
  }

  grpc::Status status =
      neighbor->SetNeighborEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_neighbor_entry_attribute(
    const sai_neighbor_entry_t *neighbor_entry, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetNeighborEntryAttributeRequest req;
  lemming::dataplane::sai::GetNeighborEntryAttributeResponse resp;
  grpc::ClientContext context;
  *req.mutable_entry() = convert_from_neighbor_entry(*neighbor_entry);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::NeighborEntryAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status =
      neighbor->GetNeighborEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_NEIGHBOR_ENTRY_ATTR_DST_MAC_ADDRESS:
        memcpy(attr_list[i].value.mac, resp.attr().dst_mac_address().data(),
               sizeof(sai_mac_t));
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().packet_action() - 1);
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_USER_TRAP_ID:
        attr_list[i].value.oid = resp.attr().user_trap_id();
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_NO_HOST_ROUTE:
        attr_list[i].value.booldata = resp.attr().no_host_route();
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_META_DATA:
        attr_list[i].value.u32 = resp.attr().meta_data();
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_COUNTER_ID:
        attr_list[i].value.oid = resp.attr().counter_id();
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_ENCAP_INDEX:
        attr_list[i].value.u32 = resp.attr().encap_index();
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_ENCAP_IMPOSE_INDEX:
        attr_list[i].value.booldata = resp.attr().encap_impose_index();
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_IS_LOCAL:
        attr_list[i].value.booldata = resp.attr().is_local();
        break;
      case SAI_NEIGHBOR_ENTRY_ATTR_IP_ADDR_FAMILY:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().ip_addr_family() - 1);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_all_neighbor_entries(sai_object_id_t switch_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}
