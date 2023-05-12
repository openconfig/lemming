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

#include "dataplane/standalone/sai/l2mc.h"

#include <glog/logging.h>

#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_l2mc_api_t l_l2mc = {
    .create_l2mc_entry = l_create_l2mc_entry,
    .remove_l2mc_entry = l_remove_l2mc_entry,
    .set_l2mc_entry_attribute = l_set_l2mc_entry_attribute,
    .get_l2mc_entry_attribute = l_get_l2mc_entry_attribute,
};

sai_status_t l_create_l2mc_entry(const sai_l2mc_entry_t *l2mc_entry,
                                 uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.l2mc_entry = l2mc_entry};
  return translator->create(SAI_OBJECT_TYPE_L2MC_ENTRY, entry, attr_count,
                            attr_list);
}

sai_status_t l_remove_l2mc_entry(const sai_l2mc_entry_t *l2mc_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.l2mc_entry = l2mc_entry};
  return translator->remove(SAI_OBJECT_TYPE_L2MC_ENTRY, entry);
}

sai_status_t l_set_l2mc_entry_attribute(const sai_l2mc_entry_t *l2mc_entry,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.l2mc_entry = l2mc_entry};
  return translator->set_attribute(SAI_OBJECT_TYPE_L2MC_ENTRY, entry, attr);
}

sai_status_t l_get_l2mc_entry_attribute(const sai_l2mc_entry_t *l2mc_entry,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.l2mc_entry = l2mc_entry};
  return translator->get_attribute(SAI_OBJECT_TYPE_L2MC_ENTRY, entry,
                                   attr_count, attr_list);
}
