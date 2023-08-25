

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

#include "dataplane/standalone/sai/queue.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/queue.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_queue_api_t l_queue = {
    .create_queue = l_create_queue,
    .remove_queue = l_remove_queue,
    .set_queue_attribute = l_set_queue_attribute,
    .get_queue_attribute = l_get_queue_attribute,
    .get_queue_stats = l_get_queue_stats,
    .get_queue_stats_ext = l_get_queue_stats_ext,
    .clear_queue_stats = l_clear_queue_stats,
};

sai_status_t l_create_queue(sai_object_id_t *queue_id,
                            sai_object_id_t switch_id, uint32_t attr_count,
                            const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateQueueRequest req;
  lemming::dataplane::sai::CreateQueueResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_QUEUE_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::QueueType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_QUEUE_ATTR_PORT:
        req.set_port(attr_list[i].value.oid);
        break;
      case SAI_QUEUE_ATTR_INDEX:
        req.set_index(attr_list[i].value.u8);
        break;
      case SAI_QUEUE_ATTR_PARENT_SCHEDULER_NODE:
        req.set_parent_scheduler_node(attr_list[i].value.oid);
        break;
      case SAI_QUEUE_ATTR_WRED_PROFILE_ID:
        req.set_wred_profile_id(attr_list[i].value.oid);
        break;
      case SAI_QUEUE_ATTR_BUFFER_PROFILE_ID:
        req.set_buffer_profile_id(attr_list[i].value.oid);
        break;
      case SAI_QUEUE_ATTR_SCHEDULER_PROFILE_ID:
        req.set_scheduler_profile_id(attr_list[i].value.oid);
        break;
      case SAI_QUEUE_ATTR_ENABLE_PFC_DLDR:
        req.set_enable_pfc_dldr(attr_list[i].value.booldata);
        break;
      case SAI_QUEUE_ATTR_PFC_DLR_INIT:
        req.set_pfc_dlr_init(attr_list[i].value.booldata);
        break;
      case SAI_QUEUE_ATTR_TAM_OBJECT:
        req.mutable_tam_object()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
    }
  }
  grpc::Status status = queue->CreateQueue(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *queue_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_QUEUE, queue_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_queue(sai_object_id_t queue_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_QUEUE, queue_id);
}

sai_status_t l_set_queue_attribute(sai_object_id_t queue_id,
                                   const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_QUEUE, queue_id, attr);
}

sai_status_t l_get_queue_attribute(sai_object_id_t queue_id,
                                   uint32_t attr_count,
                                   sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_QUEUE, queue_id, attr_count,
                                   attr_list);
}

sai_status_t l_get_queue_stats(sai_object_id_t queue_id,
                               uint32_t number_of_counters,
                               const sai_stat_id_t *counter_ids,
                               uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_QUEUE, queue_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_queue_stats_ext(sai_object_id_t queue_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids,
                                   sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_QUEUE, queue_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_queue_stats(sai_object_id_t queue_id,
                                 uint32_t number_of_counters,
                                 const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_QUEUE, queue_id,
                                 number_of_counters, counter_ids);
}
