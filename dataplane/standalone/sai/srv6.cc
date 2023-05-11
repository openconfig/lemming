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

#include "dataplane/standalone/sai/srv6.h"

#include <glog/logging.h>

#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_srv6_api_t l_srv6 = {
    .create_srv6_sidlist = l_create_srv6_sidlist,
    .remove_srv6_sidlist = l_remove_srv6_sidlist,
    .set_srv6_sidlist_attribute = l_set_srv6_sidlist_attribute,
    .get_srv6_sidlist_attribute = l_get_srv6_sidlist_attribute,
    .create_srv6_sidlists = l_create_srv6_sidlists,
    .remove_srv6_sidlists = l_remove_srv6_sidlists,
    .create_my_sid_entry = l_create_my_sid_entry,
    .remove_my_sid_entry = l_remove_my_sid_entry,
    .set_my_sid_entry_attribute = l_set_my_sid_entry_attribute,
    .get_my_sid_entry_attribute = l_get_my_sid_entry_attribute,
    .create_my_sid_entries = l_create_my_sid_entries,
    .remove_my_sid_entries = l_remove_my_sid_entries,
    .set_my_sid_entries_attribute = l_set_my_sid_entries_attribute,
    .get_my_sid_entries_attribute = l_get_my_sid_entries_attribute,
};

sai_status_t l_create_srv6_sidlist(sai_object_id_t *srv6_sidlist_id,
                                   sai_object_id_t switch_id,
                                   uint32_t attr_count,
                                   const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_SRV6_SIDLIST, srv6_sidlist_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_srv6_sidlist(sai_object_id_t srv6_sidlist_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_SRV6_SIDLIST, srv6_sidlist_id);
}

sai_status_t l_set_srv6_sidlist_attribute(sai_object_id_t srv6_sidlist_id,
                                          const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_SRV6_SIDLIST,
                                   srv6_sidlist_id, attr);
}

sai_status_t l_get_srv6_sidlist_attribute(sai_object_id_t srv6_sidlist_id,
                                          uint32_t attr_count,
                                          sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_SRV6_SIDLIST,
                                   srv6_sidlist_id, attr_count, attr_list);
}

sai_status_t l_create_srv6_sidlists(sai_object_id_t switch_id,
                                    uint32_t object_count,
                                    const uint32_t *attr_count,
                                    const sai_attribute_t **attr_list,
                                    sai_bulk_op_error_mode_t mode,
                                    sai_object_id_t *object_id,
                                    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create_bulk(SAI_OBJECT_TYPE_SRV6_SIDLIST, switch_id,
                                 object_count, attr_count, attr_list, mode,
                                 object_id, object_statuses);
}

sai_status_t l_remove_srv6_sidlists(uint32_t object_count,
                                    const sai_object_id_t *object_id,
                                    sai_bulk_op_error_mode_t mode,
                                    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove_bulk(SAI_OBJECT_TYPE_SRV6_SIDLIST, object_count,
                                 object_id, mode, object_statuses);
}

sai_status_t l_create_my_sid_entry(const sai_my_sid_entry_t *my_sid_entry,
                                   uint32_t attr_count,
                                   const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.my_sid_entry = my_sid_entry};
  return translator->create(SAI_OBJECT_TYPE_MY_SID_ENTRY, entry, attr_count,
                            attr_list);
}

sai_status_t l_remove_my_sid_entry(const sai_my_sid_entry_t *my_sid_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.my_sid_entry = my_sid_entry};
  return translator->remove(SAI_OBJECT_TYPE_MY_SID_ENTRY, entry);
}

sai_status_t l_set_my_sid_entry_attribute(
    const sai_my_sid_entry_t *my_sid_entry, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.my_sid_entry = my_sid_entry};
  return translator->set_attribute(SAI_OBJECT_TYPE_MY_SID_ENTRY, entry, attr);
}

sai_status_t l_get_my_sid_entry_attribute(
    const sai_my_sid_entry_t *my_sid_entry, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.my_sid_entry = my_sid_entry};
  return translator->get_attribute(SAI_OBJECT_TYPE_MY_SID_ENTRY, entry,
                                   attr_count, attr_list);
}

sai_status_t l_create_my_sid_entries(uint32_t object_count,
                                     const sai_my_sid_entry_t *my_sid_entry,
                                     const uint32_t *attr_count,
                                     const sai_attribute_t **attr_list,
                                     sai_bulk_op_error_mode_t mode,
                                     sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.my_sid_entry = my_sid_entry};
  return translator->create_bulk(SAI_OBJECT_TYPE_MY_SID_ENTRY, object_count,
                                 entry, attr_count, attr_list, mode,
                                 object_statuses);
}

sai_status_t l_remove_my_sid_entries(uint32_t object_count,
                                     const sai_my_sid_entry_t *my_sid_entry,
                                     sai_bulk_op_error_mode_t mode,
                                     sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.my_sid_entry = my_sid_entry};
  return translator->remove_bulk(SAI_OBJECT_TYPE_MY_SID_ENTRY, object_count,
                                 entry, mode, object_statuses);
}

sai_status_t l_set_my_sid_entries_attribute(
    uint32_t object_count, const sai_my_sid_entry_t *my_sid_entry,
    const sai_attribute_t *attr_list, sai_bulk_op_error_mode_t mode,
    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.my_sid_entry = my_sid_entry};
  return translator->set_attribute_bulk(SAI_OBJECT_TYPE_MY_SID_ENTRY,
                                        object_count, entry, attr_list, mode,
                                        object_statuses);
}

sai_status_t l_get_my_sid_entries_attribute(
    uint32_t object_count, const sai_my_sid_entry_t *my_sid_entry,
    const uint32_t *attr_count, sai_attribute_t **attr_list,
    sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.my_sid_entry = my_sid_entry};
  return translator->get_attribute_bulk(SAI_OBJECT_TYPE_MY_SID_ENTRY,
                                        object_count, entry, attr_count,
                                        attr_list, mode, object_statuses);
}
