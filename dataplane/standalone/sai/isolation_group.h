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

#ifndef DATAPLANE_STANDALONE_SAI_ISOLATION_GROUP_H_
#define DATAPLANE_STANDALONE_SAI_ISOLATION_GROUP_H_

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

extern const sai_isolation_group_api_t l_isolation_group;

sai_status_t l_create_isolation_group(sai_object_id_t *isolation_group_id,
                                      sai_object_id_t switch_id,
                                      uint32_t attr_count,
                                      const sai_attribute_t *attr_list);

sai_status_t l_remove_isolation_group(sai_object_id_t isolation_group_id);

sai_status_t l_set_isolation_group_attribute(sai_object_id_t isolation_group_id,
                                             const sai_attribute_t *attr);

sai_status_t l_get_isolation_group_attribute(sai_object_id_t isolation_group_id,
                                             uint32_t attr_count,
                                             sai_attribute_t *attr_list);

sai_status_t l_create_isolation_group_member(
    sai_object_id_t *isolation_group_member_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_isolation_group_member(
    sai_object_id_t isolation_group_member_id);

sai_status_t l_set_isolation_group_member_attribute(
    sai_object_id_t isolation_group_member_id, const sai_attribute_t *attr);

sai_status_t l_get_isolation_group_member_attribute(
    sai_object_id_t isolation_group_member_id, uint32_t attr_count,
    sai_attribute_t *attr_list);

#endif  // DATAPLANE_STANDALONE_SAI_ISOLATION_GROUP_H_
