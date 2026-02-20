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

#ifndef DATAPLANE_STANDALONE_SAI_NEXT_HOP_H_
#define DATAPLANE_STANDALONE_SAI_NEXT_HOP_H_

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

extern const sai_next_hop_api_t l_next_hop;

sai_status_t l_create_next_hop(sai_object_id_t* next_hop_id,
                               sai_object_id_t switch_id, uint32_t attr_count,
                               const sai_attribute_t* attr_list);

sai_status_t l_remove_next_hop(sai_object_id_t next_hop_id);

sai_status_t l_set_next_hop_attribute(sai_object_id_t next_hop_id,
                                      const sai_attribute_t* attr);

sai_status_t l_get_next_hop_attribute(sai_object_id_t next_hop_id,
                                      uint32_t attr_count,
                                      sai_attribute_t* attr_list);

sai_status_t l_create_next_hops(sai_object_id_t switch_id,
                                uint32_t object_count,
                                const uint32_t* attr_count,
                                const sai_attribute_t** attr_list,
                                sai_bulk_op_error_mode_t mode,
                                sai_object_id_t* object_id,
                                sai_status_t* object_statuses);

sai_status_t l_remove_next_hops(uint32_t object_count,
                                const sai_object_id_t* object_id,
                                sai_bulk_op_error_mode_t mode,
                                sai_status_t* object_statuses);

sai_status_t l_set_next_hops_attribute(uint32_t object_count,
                                       const sai_object_id_t* object_id,
                                       const sai_attribute_t* attr_list,
                                       sai_bulk_op_error_mode_t mode,
                                       sai_status_t* object_statuses);

sai_status_t l_get_next_hops_attribute(uint32_t object_count,
                                       const sai_object_id_t* object_id,
                                       const uint32_t* attr_count,
                                       sai_attribute_t** attr_list,
                                       sai_bulk_op_error_mode_t mode,
                                       sai_status_t* object_statuses);

#endif  // DATAPLANE_STANDALONE_SAI_NEXT_HOP_H_
