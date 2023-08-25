

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

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/policer.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_policer_api_t l_policer = {
    .create_policer = l_create_policer,
    .remove_policer = l_remove_policer,
    .set_policer_attribute = l_set_policer_attribute,
    .get_policer_attribute = l_get_policer_attribute,
    .get_policer_stats = l_get_policer_stats,
    .get_policer_stats_ext = l_get_policer_stats_ext,
    .clear_policer_stats = l_clear_policer_stats,
};

sai_status_t l_create_policer(sai_object_id_t *policer_id,
                              sai_object_id_t switch_id, uint32_t attr_count,
                              const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreatePolicerRequest req;
  lemming::dataplane::sai::CreatePolicerResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_POLICER_ATTR_METER_TYPE:
        req.set_meter_type(static_cast<lemming::dataplane::sai::MeterType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_MODE:
        req.set_mode(static_cast<lemming::dataplane::sai::PolicerMode>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_COLOR_SOURCE:
        req.set_color_source(
            static_cast<lemming::dataplane::sai::PolicerColorSource>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_CBS:
        req.set_cbs(attr_list[i].value.u64);
        break;
      case SAI_POLICER_ATTR_CIR:
        req.set_cir(attr_list[i].value.u64);
        break;
      case SAI_POLICER_ATTR_PBS:
        req.set_pbs(attr_list[i].value.u64);
        break;
      case SAI_POLICER_ATTR_PIR:
        req.set_pir(attr_list[i].value.u64);
        break;
      case SAI_POLICER_ATTR_GREEN_PACKET_ACTION:
        req.set_green_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_YELLOW_PACKET_ACTION:
        req.set_yellow_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_POLICER_ATTR_RED_PACKET_ACTION:
        req.set_red_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = policer->CreatePolicer(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *policer_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_POLICER, policer_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_policer(sai_object_id_t policer_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_POLICER, policer_id);
}

sai_status_t l_set_policer_attribute(sai_object_id_t policer_id,
                                     const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_POLICER, policer_id, attr);
}

sai_status_t l_get_policer_attribute(sai_object_id_t policer_id,
                                     uint32_t attr_count,
                                     sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_POLICER, policer_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_policer_stats(sai_object_id_t policer_id,
                                 uint32_t number_of_counters,
                                 const sai_stat_id_t *counter_ids,
                                 uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_POLICER, policer_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_policer_stats_ext(sai_object_id_t policer_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     sai_stats_mode_t mode,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_POLICER, policer_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_policer_stats(sai_object_id_t policer_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_POLICER, policer_id,
                                 number_of_counters, counter_ids);
}
