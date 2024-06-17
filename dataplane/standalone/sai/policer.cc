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

#include "dataplane/standalone/sai/policer.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/policer.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_policer_api_t l_policer = {
    .create_policer = l_create_policer,
    .remove_policer = l_remove_policer,
    .set_policer_attribute = l_set_policer_attribute,
    .get_policer_attribute = l_get_policer_attribute,
    .get_policer_stats = l_get_policer_stats,
    .get_policer_stats_ext = l_get_policer_stats_ext,
    .clear_policer_stats = l_clear_policer_stats,
};

lemming::dataplane::sai::CreatePolicerRequest convert_create_policer(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreatePolicerRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_POLICER_ATTR_METER_TYPE:
        msg.set_meter_type(static_cast<lemming::dataplane::sai::MeterType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_MODE:
        msg.set_mode(static_cast<lemming::dataplane::sai::PolicerMode>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_COLOR_SOURCE:
        msg.set_color_source(
            static_cast<lemming::dataplane::sai::PolicerColorSource>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_CBS:
        msg.set_cbs(attr_list[i].value.u64);
        break;
      case SAI_POLICER_ATTR_CIR:
        msg.set_cir(attr_list[i].value.u64);
        break;
      case SAI_POLICER_ATTR_PBS:
        msg.set_pbs(attr_list[i].value.u64);
        break;
      case SAI_POLICER_ATTR_PIR:
        msg.set_pir(attr_list[i].value.u64);
        break;
      case SAI_POLICER_ATTR_GREEN_PACKET_ACTION:
        msg.set_green_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_YELLOW_PACKET_ACTION:
        msg.set_yellow_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_RED_PACKET_ACTION:
        msg.set_red_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_OBJECT_STAGE:
        msg.set_object_stage(static_cast<lemming::dataplane::sai::ObjectStage>(
            attr_list[i].value.s32 + 1));
        break;
    }
  }
  return msg;
}

sai_status_t l_create_policer(sai_object_id_t *policer_id,
                              sai_object_id_t switch_id, uint32_t attr_count,
                              const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePolicerRequest req =
      convert_create_policer(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreatePolicerResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = policer->CreatePolicer(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *policer_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_policer(sai_object_id_t policer_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemovePolicerRequest req;
  lemming::dataplane::sai::RemovePolicerResponse resp;
  grpc::ClientContext context;
  req.set_oid(policer_id);

  grpc::Status status = policer->RemovePolicer(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_policer_attribute(sai_object_id_t policer_id,
                                     const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetPolicerAttributeRequest req;
  lemming::dataplane::sai::SetPolicerAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(policer_id);

  switch (attr->id) {
    case SAI_POLICER_ATTR_CBS:
      req.set_cbs(attr->value.u64);
      break;
    case SAI_POLICER_ATTR_CIR:
      req.set_cir(attr->value.u64);
      break;
    case SAI_POLICER_ATTR_PBS:
      req.set_pbs(attr->value.u64);
      break;
    case SAI_POLICER_ATTR_PIR:
      req.set_pir(attr->value.u64);
      break;
    case SAI_POLICER_ATTR_GREEN_PACKET_ACTION:
      req.set_green_packet_action(
          static_cast<lemming::dataplane::sai::PacketAction>(attr->value.s32 +
                                                             1));
      break;
    case SAI_POLICER_ATTR_YELLOW_PACKET_ACTION:
      req.set_yellow_packet_action(
          static_cast<lemming::dataplane::sai::PacketAction>(attr->value.s32 +
                                                             1));
      break;
    case SAI_POLICER_ATTR_RED_PACKET_ACTION:
      req.set_red_packet_action(
          static_cast<lemming::dataplane::sai::PacketAction>(attr->value.s32 +
                                                             1));
      break;
  }

  grpc::Status status = policer->SetPolicerAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_policer_attribute(sai_object_id_t policer_id,
                                     uint32_t attr_count,
                                     sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetPolicerAttributeRequest req;
  lemming::dataplane::sai::GetPolicerAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(policer_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::PolicerAttr>(attr_list[i].id + 1));
  }
  grpc::Status status = policer->GetPolicerAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_POLICER_ATTR_METER_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().meter_type() - 1);
        break;
      case SAI_POLICER_ATTR_MODE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().mode() - 1);
        break;
      case SAI_POLICER_ATTR_COLOR_SOURCE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().color_source() - 1);
        break;
      case SAI_POLICER_ATTR_CBS:
        attr_list[i].value.u64 = resp.attr().cbs();
        break;
      case SAI_POLICER_ATTR_CIR:
        attr_list[i].value.u64 = resp.attr().cir();
        break;
      case SAI_POLICER_ATTR_PBS:
        attr_list[i].value.u64 = resp.attr().pbs();
        break;
      case SAI_POLICER_ATTR_PIR:
        attr_list[i].value.u64 = resp.attr().pir();
        break;
      case SAI_POLICER_ATTR_GREEN_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().green_packet_action() - 1);
        break;
      case SAI_POLICER_ATTR_YELLOW_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().yellow_packet_action() - 1);
        break;
      case SAI_POLICER_ATTR_RED_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().red_packet_action() - 1);
        break;
      case SAI_POLICER_ATTR_OBJECT_STAGE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().object_stage() - 1);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_policer_stats(sai_object_id_t policer_id,
                                 uint32_t number_of_counters,
                                 const sai_stat_id_t *counter_ids,
                                 uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetPolicerStatsRequest req;
  lemming::dataplane::sai::GetPolicerStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(policer_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(
        static_cast<lemming::dataplane::sai::PolicerStat>(counter_ids[i] + 1));
  }
  grpc::Status status = policer->GetPolicerStats(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_policer_stats_ext(sai_object_id_t policer_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     sai_stats_mode_t mode,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_policer_stats(sai_object_id_t policer_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
