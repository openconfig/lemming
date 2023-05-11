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

#include "dataplane/standalone/sai/scheduler.h"

#include <glog/logging.h>

#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_scheduler_api_t l_scheduler = {
    .create_scheduler = l_create_scheduler,
    .remove_scheduler = l_remove_scheduler,
    .set_scheduler_attribute = l_set_scheduler_attribute,
    .get_scheduler_attribute = l_get_scheduler_attribute,
};

sai_status_t l_create_scheduler(sai_object_id_t *scheduler_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_SCHEDULER, scheduler_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_scheduler(sai_object_id_t scheduler_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_SCHEDULER, scheduler_id);
}

sai_status_t l_set_scheduler_attribute(sai_object_id_t scheduler_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_SCHEDULER, scheduler_id,
                                   attr);
}

sai_status_t l_get_scheduler_attribute(sai_object_id_t scheduler_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_SCHEDULER, scheduler_id,
                                   attr_count, attr_list);
}
