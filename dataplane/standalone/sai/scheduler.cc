

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

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/scheduler.pb.h"
#include "dataplane/standalone/sai/common.h"

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

  lemming::dataplane::sai::CreateSchedulerRequest req;
  lemming::dataplane::sai::CreateSchedulerResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SCHEDULER_ATTR_SCHEDULING_TYPE:
        req.set_scheduling_type(
            static_cast<lemming::dataplane::sai::SchedulingType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_SCHEDULER_ATTR_SCHEDULING_WEIGHT:
        req.set_scheduling_weight(attr_list[i].value.u8);
        break;
      case SAI_SCHEDULER_ATTR_METER_TYPE:
        req.set_meter_type(static_cast<lemming::dataplane::sai::MeterType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_SCHEDULER_ATTR_MIN_BANDWIDTH_RATE:
        req.set_min_bandwidth_rate(attr_list[i].value.u64);
        break;
      case SAI_SCHEDULER_ATTR_MIN_BANDWIDTH_BURST_RATE:
        req.set_min_bandwidth_burst_rate(attr_list[i].value.u64);
        break;
      case SAI_SCHEDULER_ATTR_MAX_BANDWIDTH_RATE:
        req.set_max_bandwidth_rate(attr_list[i].value.u64);
        break;
      case SAI_SCHEDULER_ATTR_MAX_BANDWIDTH_BURST_RATE:
        req.set_max_bandwidth_burst_rate(attr_list[i].value.u64);
        break;
    }
  }
  grpc::Status status = scheduler->CreateScheduler(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *scheduler_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_scheduler(sai_object_id_t scheduler_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveSchedulerRequest req;
  lemming::dataplane::sai::RemoveSchedulerResponse resp;
  grpc::ClientContext context;
  req.set_oid(scheduler_id);

  grpc::Status status = scheduler->RemoveScheduler(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_scheduler_attribute(sai_object_id_t scheduler_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetSchedulerAttributeRequest req;
  lemming::dataplane::sai::SetSchedulerAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(scheduler_id);

  switch (attr->id) {
    case SAI_SCHEDULER_ATTR_SCHEDULING_TYPE:
      req.set_scheduling_type(
          static_cast<lemming::dataplane::sai::SchedulingType>(attr->value.s32 +
                                                               1));
      break;
    case SAI_SCHEDULER_ATTR_SCHEDULING_WEIGHT:
      req.set_scheduling_weight(attr->value.u8);
      break;
    case SAI_SCHEDULER_ATTR_METER_TYPE:
      req.set_meter_type(
          static_cast<lemming::dataplane::sai::MeterType>(attr->value.s32 + 1));
      break;
    case SAI_SCHEDULER_ATTR_MIN_BANDWIDTH_RATE:
      req.set_min_bandwidth_rate(attr->value.u64);
      break;
    case SAI_SCHEDULER_ATTR_MIN_BANDWIDTH_BURST_RATE:
      req.set_min_bandwidth_burst_rate(attr->value.u64);
      break;
    case SAI_SCHEDULER_ATTR_MAX_BANDWIDTH_RATE:
      req.set_max_bandwidth_rate(attr->value.u64);
      break;
    case SAI_SCHEDULER_ATTR_MAX_BANDWIDTH_BURST_RATE:
      req.set_max_bandwidth_burst_rate(attr->value.u64);
      break;
  }

  grpc::Status status = scheduler->SetSchedulerAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_scheduler_attribute(sai_object_id_t scheduler_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetSchedulerAttributeRequest req;
  lemming::dataplane::sai::GetSchedulerAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(scheduler_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::SchedulerAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = scheduler->GetSchedulerAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SCHEDULER_ATTR_SCHEDULING_TYPE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().scheduling_type() - 1);
        break;
      case SAI_SCHEDULER_ATTR_SCHEDULING_WEIGHT:
        attr_list[i].value.u8 = resp.attr().scheduling_weight();
        break;
      case SAI_SCHEDULER_ATTR_METER_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().meter_type() - 1);
        break;
      case SAI_SCHEDULER_ATTR_MIN_BANDWIDTH_RATE:
        attr_list[i].value.u64 = resp.attr().min_bandwidth_rate();
        break;
      case SAI_SCHEDULER_ATTR_MIN_BANDWIDTH_BURST_RATE:
        attr_list[i].value.u64 = resp.attr().min_bandwidth_burst_rate();
        break;
      case SAI_SCHEDULER_ATTR_MAX_BANDWIDTH_RATE:
        attr_list[i].value.u64 = resp.attr().max_bandwidth_rate();
        break;
      case SAI_SCHEDULER_ATTR_MAX_BANDWIDTH_BURST_RATE:
        attr_list[i].value.u64 = resp.attr().max_bandwidth_burst_rate();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
