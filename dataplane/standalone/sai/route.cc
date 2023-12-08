

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

#include "dataplane/standalone/sai/route.h"

#include <glog/logging.h>

#include "dataplane/proto/common.pb.h"
#include "dataplane/proto/route.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_route_api_t l_route = {
    .create_route_entry = l_create_route_entry,
    .remove_route_entry = l_remove_route_entry,
    .set_route_entry_attribute = l_set_route_entry_attribute,
    .get_route_entry_attribute = l_get_route_entry_attribute,
    .create_route_entries = l_create_route_entries,
    .remove_route_entries = l_remove_route_entries,
    .set_route_entries_attribute = l_set_route_entries_attribute,
    .get_route_entries_attribute = l_get_route_entries_attribute,
};

lemming::dataplane::sai::CreateRouteEntryRequest convert_create_route_entry(
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateRouteEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ROUTE_ENTRY_ATTR_PACKET_ACTION:
        msg.set_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_ROUTE_ENTRY_ATTR_USER_TRAP_ID:
        msg.set_user_trap_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTE_ENTRY_ATTR_NEXT_HOP_ID:
        msg.set_next_hop_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTE_ENTRY_ATTR_META_DATA:
        msg.set_meta_data(attr_list[i].value.u32);
        break;
      case SAI_ROUTE_ENTRY_ATTR_COUNTER_ID:
        msg.set_counter_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_route_entry(const sai_route_entry_t *route_entry,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateRouteEntryRequest req =
      convert_create_route_entry(attr_count, attr_list);
  lemming::dataplane::sai::CreateRouteEntryResponse resp;
  grpc::ClientContext context;

  *req.mutable_entry() = convert_from_route_entry(*route_entry);
  grpc::Status status = route->CreateRouteEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_route_entry(const sai_route_entry_t *route_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveRouteEntryRequest req;
  lemming::dataplane::sai::RemoveRouteEntryResponse resp;
  grpc::ClientContext context;

  *req.mutable_entry() = convert_from_route_entry(*route_entry);
  grpc::Status status = route->RemoveRouteEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_route_entry_attribute(const sai_route_entry_t *route_entry,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetRouteEntryAttributeRequest req;
  lemming::dataplane::sai::SetRouteEntryAttributeResponse resp;
  grpc::ClientContext context;

  *req.mutable_entry() = convert_from_route_entry(*route_entry);

  switch (attr->id) {
    case SAI_ROUTE_ENTRY_ATTR_PACKET_ACTION:
      req.set_packet_action(static_cast<lemming::dataplane::sai::PacketAction>(
          attr->value.s32 + 1));
      break;
    case SAI_ROUTE_ENTRY_ATTR_USER_TRAP_ID:
      req.set_user_trap_id(attr->value.oid);
      break;
    case SAI_ROUTE_ENTRY_ATTR_NEXT_HOP_ID:
      req.set_next_hop_id(attr->value.oid);
      break;
    case SAI_ROUTE_ENTRY_ATTR_META_DATA:
      req.set_meta_data(attr->value.u32);
      break;
    case SAI_ROUTE_ENTRY_ATTR_COUNTER_ID:
      req.set_counter_id(attr->value.oid);
      break;
  }

  grpc::Status status = route->SetRouteEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_route_entry_attribute(const sai_route_entry_t *route_entry,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetRouteEntryAttributeRequest req;
  lemming::dataplane::sai::GetRouteEntryAttributeResponse resp;
  grpc::ClientContext context;
  *req.mutable_entry() = convert_from_route_entry(*route_entry);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::RouteEntryAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = route->GetRouteEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ROUTE_ENTRY_ATTR_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().packet_action() - 1);
        break;
      case SAI_ROUTE_ENTRY_ATTR_USER_TRAP_ID:
        attr_list[i].value.oid = resp.attr().user_trap_id();
        break;
      case SAI_ROUTE_ENTRY_ATTR_NEXT_HOP_ID:
        attr_list[i].value.oid = resp.attr().next_hop_id();
        break;
      case SAI_ROUTE_ENTRY_ATTR_META_DATA:
        attr_list[i].value.u32 = resp.attr().meta_data();
        break;
      case SAI_ROUTE_ENTRY_ATTR_IP_ADDR_FAMILY:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().ip_addr_family() - 1);
        break;
      case SAI_ROUTE_ENTRY_ATTR_COUNTER_ID:
        attr_list[i].value.oid = resp.attr().counter_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_route_entries(uint32_t object_count,
                                    const sai_route_entry_t *route_entry,
                                    const uint32_t *attr_count,
                                    const sai_attribute_t **attr_list,
                                    sai_bulk_op_error_mode_t mode,
                                    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateRouteEntriesRequest req;
  lemming::dataplane::sai::CreateRouteEntriesResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    auto r = convert_create_route_entry(attr_count[i], attr_list[i]);

    *r.mutable_entry() = convert_from_route_entry(*route_entry);
    *req.add_reqs() = r;
  }

  grpc::Status status = route->CreateRouteEntries(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < object_count; i++) {
    object_statuses[i] = SAI_STATUS_SUCCESS;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_route_entries(uint32_t object_count,
                                    const sai_route_entry_t *route_entry,
                                    sai_bulk_op_error_mode_t mode,
                                    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_route_entries_attribute(uint32_t object_count,
                                           const sai_route_entry_t *route_entry,
                                           const sai_attribute_t *attr_list,
                                           sai_bulk_op_error_mode_t mode,
                                           sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_route_entries_attribute(uint32_t object_count,
                                           const sai_route_entry_t *route_entry,
                                           const uint32_t *attr_count,
                                           sai_attribute_t **attr_list,
                                           sai_bulk_op_error_mode_t mode,
                                           sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}
