

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

#include "dataplane/standalone/sai/mcast_fdb.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/mcast_fdb.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_mcast_fdb_api_t l_mcast_fdb = {
    .create_mcast_fdb_entry = l_create_mcast_fdb_entry,
    .remove_mcast_fdb_entry = l_remove_mcast_fdb_entry,
    .set_mcast_fdb_entry_attribute = l_set_mcast_fdb_entry_attribute,
    .get_mcast_fdb_entry_attribute = l_get_mcast_fdb_entry_attribute,
};

sai_status_t l_create_mcast_fdb_entry(
    const sai_mcast_fdb_entry_t *mcast_fdb_entry, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateMcastFdbEntryRequest req;
  lemming::dataplane::sai::CreateMcastFdbEntryResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MCAST_FDB_ENTRY_ATTR_GROUP_ID:
        req.set_group_id(attr_list[i].value.oid);
        break;
      case SAI_MCAST_FDB_ENTRY_ATTR_PACKET_ACTION:
        req.set_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_MCAST_FDB_ENTRY_ATTR_META_DATA:
        req.set_meta_data(attr_list[i].value.u32);
        break;
    }
  }
  grpc::Status status = mcast_fdb->CreateMcastFdbEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  common_entry_t entry = {.mcast_fdb_entry = mcast_fdb_entry};
  return translator->create(SAI_OBJECT_TYPE_MCAST_FDB_ENTRY, entry, attr_count,
                            attr_list);
}

sai_status_t l_remove_mcast_fdb_entry(
    const sai_mcast_fdb_entry_t *mcast_fdb_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.mcast_fdb_entry = mcast_fdb_entry};
  return translator->remove(SAI_OBJECT_TYPE_MCAST_FDB_ENTRY, entry);
}

sai_status_t l_set_mcast_fdb_entry_attribute(
    const sai_mcast_fdb_entry_t *mcast_fdb_entry, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.mcast_fdb_entry = mcast_fdb_entry};
  return translator->set_attribute(SAI_OBJECT_TYPE_MCAST_FDB_ENTRY, entry,
                                   attr);
}

sai_status_t l_get_mcast_fdb_entry_attribute(
    const sai_mcast_fdb_entry_t *mcast_fdb_entry, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.mcast_fdb_entry = mcast_fdb_entry};
  return translator->get_attribute(SAI_OBJECT_TYPE_MCAST_FDB_ENTRY, entry,
                                   attr_count, attr_list);
}
