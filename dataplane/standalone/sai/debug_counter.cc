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

#include "dataplane/standalone/sai/debug_counter.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/debug_counter.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_debug_counter_api_t l_debug_counter = {
    .create_debug_counter = l_create_debug_counter,
    .remove_debug_counter = l_remove_debug_counter,
    .set_debug_counter_attribute = l_set_debug_counter_attribute,
    .get_debug_counter_attribute = l_get_debug_counter_attribute,
};

lemming::dataplane::sai::CreateDebugCounterRequest convert_create_debug_counter(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateDebugCounterRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DEBUG_COUNTER_ATTR_TYPE:
        msg.set_type(
            convert_sai_debug_counter_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_DEBUG_COUNTER_ATTR_BIND_METHOD:
        msg.set_bind_method(convert_sai_debug_counter_bind_method_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_DEBUG_COUNTER_ATTR_IN_DROP_REASON_LIST:
        msg.mutable_in_drop_reason_list()->CopyFrom(
            convert_list_sai_in_drop_reason_t_to_proto(
                attr_list[i].value.s32list));
        break;
      case SAI_DEBUG_COUNTER_ATTR_OUT_DROP_REASON_LIST:
        msg.mutable_out_drop_reason_list()->CopyFrom(
            convert_list_sai_out_drop_reason_t_to_proto(
                attr_list[i].value.s32list));
        break;
    }
  }
  return msg;
}

sai_status_t l_create_debug_counter(sai_object_id_t *debug_counter_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateDebugCounterRequest req =
      convert_create_debug_counter(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateDebugCounterResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = debug_counter->CreateDebugCounter(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (debug_counter_id) {
    *debug_counter_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_debug_counter(sai_object_id_t debug_counter_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveDebugCounterRequest req;
  lemming::dataplane::sai::RemoveDebugCounterResponse resp;
  grpc::ClientContext context;
  req.set_oid(debug_counter_id);

  grpc::Status status = debug_counter->RemoveDebugCounter(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_debug_counter_attribute(sai_object_id_t debug_counter_id,
                                           const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetDebugCounterAttributeRequest req;
  lemming::dataplane::sai::SetDebugCounterAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(debug_counter_id);

  switch (attr->id) {
    case SAI_DEBUG_COUNTER_ATTR_IN_DROP_REASON_LIST:
      req.mutable_in_drop_reason_list()->CopyFrom(
          convert_list_sai_in_drop_reason_t_to_proto(attr->value.s32list));
      break;
    case SAI_DEBUG_COUNTER_ATTR_OUT_DROP_REASON_LIST:
      req.mutable_out_drop_reason_list()->CopyFrom(
          convert_list_sai_out_drop_reason_t_to_proto(attr->value.s32list));
      break;
  }

  grpc::Status status =
      debug_counter->SetDebugCounterAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_debug_counter_attribute(sai_object_id_t debug_counter_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetDebugCounterAttributeRequest req;
  lemming::dataplane::sai::GetDebugCounterAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(debug_counter_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_debug_counter_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status =
      debug_counter->GetDebugCounterAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DEBUG_COUNTER_ATTR_INDEX:
        attr_list[i].value.u32 = resp.attr().index();
        break;
      case SAI_DEBUG_COUNTER_ATTR_TYPE:
        attr_list[i].value.s32 =
            convert_sai_debug_counter_type_t_to_sai(resp.attr().type());
        break;
      case SAI_DEBUG_COUNTER_ATTR_BIND_METHOD:
        attr_list[i].value.s32 = convert_sai_debug_counter_bind_method_t_to_sai(
            resp.attr().bind_method());
        break;
      case SAI_DEBUG_COUNTER_ATTR_IN_DROP_REASON_LIST:
        convert_list_sai_in_drop_reason_t_to_sai(
            attr_list[i].value.s32list.list, resp.attr().in_drop_reason_list(),
            &attr_list[i].value.s32list.count);
        break;
      case SAI_DEBUG_COUNTER_ATTR_OUT_DROP_REASON_LIST:
        convert_list_sai_out_drop_reason_t_to_sai(
            attr_list[i].value.s32list.list, resp.attr().out_drop_reason_list(),
            &attr_list[i].value.s32list.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
