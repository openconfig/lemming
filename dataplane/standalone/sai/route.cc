

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

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/route.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

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

sai_status_t l_create_route_entry(const sai_route_entry_t *route_entry,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateRouteEntryRequest req;
  lemming::dataplane::sai::CreateRouteEntryResponse resp;
  grpc::ClientContext context;

  *req.mutable_entry() = convert_from_route_entry(*route_entry);
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ROUTE_ENTRY_ATTR_PACKET_ACTION:
        req.set_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_ROUTE_ENTRY_ATTR_USER_TRAP_ID:
        req.set_user_trap_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTE_ENTRY_ATTR_NEXT_HOP_ID:
        req.set_next_hop_id(attr_list[i].value.oid);
        break;
      case SAI_ROUTE_ENTRY_ATTR_META_DATA:
        req.set_meta_data(attr_list[i].value.u32);
        break;
      case SAI_ROUTE_ENTRY_ATTR_COUNTER_ID:
        req.set_counter_id(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = route->CreateRouteEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  common_entry_t entry = {.route_entry = route_entry};
  return translator->create(SAI_OBJECT_TYPE_ROUTE_ENTRY, entry, attr_count,
                            attr_list);
}

sai_status_t l_remove_route_entry(const sai_route_entry_t *route_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.route_entry = route_entry};
  return translator->remove(SAI_OBJECT_TYPE_ROUTE_ENTRY, entry);
}

sai_status_t l_set_route_entry_attribute(const sai_route_entry_t *route_entry,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.route_entry = route_entry};
  return translator->set_attribute(SAI_OBJECT_TYPE_ROUTE_ENTRY, entry, attr);
}

sai_status_t l_get_route_entry_attribute(const sai_route_entry_t *route_entry,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.route_entry = route_entry};
  return translator->get_attribute(SAI_OBJECT_TYPE_ROUTE_ENTRY, entry,
                                   attr_count, attr_list);
}

sai_status_t l_create_route_entries(uint32_t object_count,
                                    const sai_route_entry_t *route_entry,
                                    const uint32_t *attr_count,
                                    const sai_attribute_t **attr_list,
                                    sai_bulk_op_error_mode_t mode,
                                    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.route_entry = route_entry};
  return translator->create_bulk(SAI_OBJECT_TYPE_ROUTE_ENTRY, object_count,
                                 entry, attr_count, attr_list, mode,
                                 object_statuses);
}

sai_status_t l_remove_route_entries(uint32_t object_count,
                                    const sai_route_entry_t *route_entry,
                                    sai_bulk_op_error_mode_t mode,
                                    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.route_entry = route_entry};
  return translator->remove_bulk(SAI_OBJECT_TYPE_ROUTE_ENTRY, object_count,
                                 entry, mode, object_statuses);
}

sai_status_t l_set_route_entries_attribute(uint32_t object_count,
                                           const sai_route_entry_t *route_entry,
                                           const sai_attribute_t *attr_list,
                                           sai_bulk_op_error_mode_t mode,
                                           sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.route_entry = route_entry};
  return translator->set_attribute_bulk(SAI_OBJECT_TYPE_ROUTE_ENTRY,
                                        object_count, entry, attr_list, mode,
                                        object_statuses);
}

sai_status_t l_get_route_entries_attribute(uint32_t object_count,
                                           const sai_route_entry_t *route_entry,
                                           const uint32_t *attr_count,
                                           sai_attribute_t **attr_list,
                                           sai_bulk_op_error_mode_t mode,
                                           sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.route_entry = route_entry};
  return translator->get_attribute_bulk(SAI_OBJECT_TYPE_ROUTE_ENTRY,
                                        object_count, entry, attr_count,
                                        attr_list, mode, object_statuses);
}
