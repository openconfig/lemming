

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

#include "dataplane/standalone/sai/wred.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/wred.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_wred_api_t l_wred = {
    .create_wred = l_create_wred,
    .remove_wred = l_remove_wred,
    .set_wred_attribute = l_set_wred_attribute,
    .get_wred_attribute = l_get_wred_attribute,
};

sai_status_t l_create_wred(sai_object_id_t *wred_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateWredRequest req;
  lemming::dataplane::sai::CreateWredResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_WRED_ATTR_GREEN_ENABLE:
        req.set_green_enable(attr_list[i].value.booldata);
        break;
      case SAI_WRED_ATTR_GREEN_MIN_THRESHOLD:
        req.set_green_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_GREEN_MAX_THRESHOLD:
        req.set_green_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_GREEN_DROP_PROBABILITY:
        req.set_green_drop_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_YELLOW_ENABLE:
        req.set_yellow_enable(attr_list[i].value.booldata);
        break;
      case SAI_WRED_ATTR_YELLOW_MIN_THRESHOLD:
        req.set_yellow_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_YELLOW_MAX_THRESHOLD:
        req.set_yellow_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_YELLOW_DROP_PROBABILITY:
        req.set_yellow_drop_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_RED_ENABLE:
        req.set_red_enable(attr_list[i].value.booldata);
        break;
      case SAI_WRED_ATTR_RED_MIN_THRESHOLD:
        req.set_red_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_RED_MAX_THRESHOLD:
        req.set_red_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_RED_DROP_PROBABILITY:
        req.set_red_drop_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_WEIGHT:
        req.set_weight(attr_list[i].value.u8);
        break;
      case SAI_WRED_ATTR_ECN_MARK_MODE:
        req.set_ecn_mark_mode(static_cast<lemming::dataplane::sai::EcnMarkMode>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_WRED_ATTR_ECN_GREEN_MIN_THRESHOLD:
        req.set_ecn_green_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_GREEN_MAX_THRESHOLD:
        req.set_ecn_green_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_GREEN_MARK_PROBABILITY:
        req.set_ecn_green_mark_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_YELLOW_MIN_THRESHOLD:
        req.set_ecn_yellow_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_YELLOW_MAX_THRESHOLD:
        req.set_ecn_yellow_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_YELLOW_MARK_PROBABILITY:
        req.set_ecn_yellow_mark_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_RED_MIN_THRESHOLD:
        req.set_ecn_red_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_RED_MAX_THRESHOLD:
        req.set_ecn_red_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_RED_MARK_PROBABILITY:
        req.set_ecn_red_mark_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MIN_THRESHOLD:
        req.set_ecn_color_unaware_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MAX_THRESHOLD:
        req.set_ecn_color_unaware_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MARK_PROBABILITY:
        req.set_ecn_color_unaware_mark_probability(attr_list[i].value.u32);
        break;
    }
  }
  grpc::Status status = wred->CreateWred(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *wred_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_WRED, wred_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_wred(sai_object_id_t wred_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_WRED, wred_id);
}

sai_status_t l_set_wred_attribute(sai_object_id_t wred_id,
                                  const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_WRED, wred_id, attr);
}

sai_status_t l_get_wred_attribute(sai_object_id_t wred_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_WRED, wred_id, attr_count,
                                   attr_list);
}
