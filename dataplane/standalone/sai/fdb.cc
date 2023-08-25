

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

#include "dataplane/standalone/sai/fdb.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/fdb.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_fdb_api_t l_fdb = {
    .create_fdb_entry = l_create_fdb_entry,
    .remove_fdb_entry = l_remove_fdb_entry,
    .set_fdb_entry_attribute = l_set_fdb_entry_attribute,
    .get_fdb_entry_attribute = l_get_fdb_entry_attribute,
    .flush_fdb_entries = l_flush_fdb_entries,
    .create_fdb_entries = l_create_fdb_entries,
    .remove_fdb_entries = l_remove_fdb_entries,
    .set_fdb_entries_attribute = l_set_fdb_entries_attribute,
    .get_fdb_entries_attribute = l_get_fdb_entries_attribute,
};

sai_status_t l_create_fdb_entry(const sai_fdb_entry_t *fdb_entry,
                                uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateFdbEntryRequest req;
  lemming::dataplane::sai::CreateFdbEntryResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_FDB_ENTRY_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::FdbEntryType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_FDB_ENTRY_ATTR_PACKET_ACTION:
        req.set_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_FDB_ENTRY_ATTR_USER_TRAP_ID:
        req.set_user_trap_id(attr_list[i].value.oid);
        break;
      case SAI_FDB_ENTRY_ATTR_BRIDGE_PORT_ID:
        req.set_bridge_port_id(attr_list[i].value.oid);
        break;
      case SAI_FDB_ENTRY_ATTR_META_DATA:
        req.set_meta_data(attr_list[i].value.u32);
        break;
      case SAI_FDB_ENTRY_ATTR_ENDPOINT_IP:
        req.set_endpoint_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_FDB_ENTRY_ATTR_COUNTER_ID:
        req.set_counter_id(attr_list[i].value.oid);
        break;
      case SAI_FDB_ENTRY_ATTR_ALLOW_MAC_MOVE:
        req.set_allow_mac_move(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status = fdb->CreateFdbEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  common_entry_t entry = {.fdb_entry = fdb_entry};
  return translator->create(SAI_OBJECT_TYPE_FDB_ENTRY, entry, attr_count,
                            attr_list);
}

sai_status_t l_remove_fdb_entry(const sai_fdb_entry_t *fdb_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.fdb_entry = fdb_entry};
  return translator->remove(SAI_OBJECT_TYPE_FDB_ENTRY, entry);
}

sai_status_t l_set_fdb_entry_attribute(const sai_fdb_entry_t *fdb_entry,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.fdb_entry = fdb_entry};
  return translator->set_attribute(SAI_OBJECT_TYPE_FDB_ENTRY, entry, attr);
}

sai_status_t l_get_fdb_entry_attribute(const sai_fdb_entry_t *fdb_entry,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.fdb_entry = fdb_entry};
  return translator->get_attribute(SAI_OBJECT_TYPE_FDB_ENTRY, entry, attr_count,
                                   attr_list);
}

sai_status_t l_flush_fdb_entries(sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_fdb_entries(uint32_t object_count,
                                  const sai_fdb_entry_t *fdb_entry,
                                  const uint32_t *attr_count,
                                  const sai_attribute_t **attr_list,
                                  sai_bulk_op_error_mode_t mode,
                                  sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.fdb_entry = fdb_entry};
  return translator->create_bulk(SAI_OBJECT_TYPE_FDB_ENTRY, object_count, entry,
                                 attr_count, attr_list, mode, object_statuses);
}

sai_status_t l_remove_fdb_entries(uint32_t object_count,
                                  const sai_fdb_entry_t *fdb_entry,
                                  sai_bulk_op_error_mode_t mode,
                                  sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.fdb_entry = fdb_entry};
  return translator->remove_bulk(SAI_OBJECT_TYPE_FDB_ENTRY, object_count, entry,
                                 mode, object_statuses);
}

sai_status_t l_set_fdb_entries_attribute(uint32_t object_count,
                                         const sai_fdb_entry_t *fdb_entry,
                                         const sai_attribute_t *attr_list,
                                         sai_bulk_op_error_mode_t mode,
                                         sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.fdb_entry = fdb_entry};
  return translator->set_attribute_bulk(SAI_OBJECT_TYPE_FDB_ENTRY, object_count,
                                        entry, attr_list, mode,
                                        object_statuses);
}

sai_status_t l_get_fdb_entries_attribute(uint32_t object_count,
                                         const sai_fdb_entry_t *fdb_entry,
                                         const uint32_t *attr_count,
                                         sai_attribute_t **attr_list,
                                         sai_bulk_op_error_mode_t mode,
                                         sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.fdb_entry = fdb_entry};
  return translator->get_attribute_bulk(SAI_OBJECT_TYPE_FDB_ENTRY, object_count,
                                        entry, attr_count, attr_list, mode,
                                        object_statuses);
}
