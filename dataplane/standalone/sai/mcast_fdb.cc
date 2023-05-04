
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

#include "dataplane/standalone/log/log.h"

const sai_mcast_fdb_api_t l_mcast_fdb = {
    .create_mcast_fdb_entry = l_create_mcast_fdb_entry,
    .remove_mcast_fdb_entry = l_remove_mcast_fdb_entry,
    .set_mcast_fdb_entry_attribute = l_set_mcast_fdb_entry_attribute,
    .get_mcast_fdb_entry_attribute = l_get_mcast_fdb_entry_attribute,
};

sai_status_t l_create_mcast_fdb_entry(
    const sai_mcast_fdb_entry_t *mcast_fdb_entry, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_remove_mcast_fdb_entry(
    const sai_mcast_fdb_entry_t *mcast_fdb_entry) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_mcast_fdb_entry_attribute(
    const sai_mcast_fdb_entry_t *mcast_fdb_entry, const sai_attribute_t *attr) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_mcast_fdb_entry_attribute(
    const sai_mcast_fdb_entry_t *mcast_fdb_entry, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LUCIUS_LOG_FUNC();
  return SAI_STATUS_NOT_IMPLEMENTED;
}
