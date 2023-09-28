

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

#include "dataplane/standalone/sai/counter.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/counter.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_counter_api_t l_counter = {
    .create_counter = l_create_counter,
    .remove_counter = l_remove_counter,
    .set_counter_attribute = l_set_counter_attribute,
    .get_counter_attribute = l_get_counter_attribute,
    .get_counter_stats = l_get_counter_stats,
    .get_counter_stats_ext = l_get_counter_stats_ext,
    .clear_counter_stats = l_clear_counter_stats,
};

lemming::dataplane::sai::CreateCounterRequest convert_create_counter(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateCounterRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_COUNTER_ATTR_TYPE:
        msg.set_type(static_cast<lemming::dataplane::sai::CounterType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_COUNTER_ATTR_LABEL:
        msg.set_label(attr_list[i].value.chardata);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_counter(sai_object_id_t *counter_id,
                              sai_object_id_t switch_id, uint32_t attr_count,
                              const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateCounterRequest req =
      convert_create_counter(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateCounterResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = counter->CreateCounter(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *counter_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_counter(sai_object_id_t counter_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveCounterRequest req;
  lemming::dataplane::sai::RemoveCounterResponse resp;
  grpc::ClientContext context;
  req.set_oid(counter_id);

  grpc::Status status = counter->RemoveCounter(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_counter_attribute(sai_object_id_t counter_id,
                                     const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetCounterAttributeRequest req;
  lemming::dataplane::sai::SetCounterAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(counter_id);

  switch (attr->id) {
    case SAI_COUNTER_ATTR_LABEL:
      req.set_label(attr->value.chardata);
      break;
  }

  grpc::Status status = counter->SetCounterAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_counter_attribute(sai_object_id_t counter_id,
                                     uint32_t attr_count,
                                     sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetCounterAttributeRequest req;
  lemming::dataplane::sai::GetCounterAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(counter_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::CounterAttr>(attr_list[i].id + 1));
  }
  grpc::Status status = counter->GetCounterAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_COUNTER_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_COUNTER_ATTR_LABEL:
        strncpy(attr_list[i].value.chardata, resp.attr().label().data(), 32);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_counter_stats(sai_object_id_t counter_id,
                                 uint32_t number_of_counters,
                                 const sai_stat_id_t *counter_ids,
                                 uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_counter_stats_ext(sai_object_id_t counter_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     sai_stats_mode_t mode,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_counter_stats(sai_object_id_t counter_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
