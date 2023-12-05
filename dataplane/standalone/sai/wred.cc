

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

#include "dataplane/proto/common.pb.h"
#include "dataplane/proto/wred.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_wred_api_t l_wred = {
    .create_wred = l_create_wred,
    .remove_wred = l_remove_wred,
    .set_wred_attribute = l_set_wred_attribute,
    .get_wred_attribute = l_get_wred_attribute,
};

lemming::dataplane::sai::CreateWredRequest convert_create_wred(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateWredRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_WRED_ATTR_GREEN_ENABLE:
        msg.set_green_enable(attr_list[i].value.booldata);
        break;
      case SAI_WRED_ATTR_GREEN_MIN_THRESHOLD:
        msg.set_green_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_GREEN_MAX_THRESHOLD:
        msg.set_green_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_GREEN_DROP_PROBABILITY:
        msg.set_green_drop_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_YELLOW_ENABLE:
        msg.set_yellow_enable(attr_list[i].value.booldata);
        break;
      case SAI_WRED_ATTR_YELLOW_MIN_THRESHOLD:
        msg.set_yellow_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_YELLOW_MAX_THRESHOLD:
        msg.set_yellow_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_YELLOW_DROP_PROBABILITY:
        msg.set_yellow_drop_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_RED_ENABLE:
        msg.set_red_enable(attr_list[i].value.booldata);
        break;
      case SAI_WRED_ATTR_RED_MIN_THRESHOLD:
        msg.set_red_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_RED_MAX_THRESHOLD:
        msg.set_red_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_RED_DROP_PROBABILITY:
        msg.set_red_drop_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_WEIGHT:
        msg.set_weight(attr_list[i].value.u8);
        break;
      case SAI_WRED_ATTR_ECN_MARK_MODE:
        msg.set_ecn_mark_mode(static_cast<lemming::dataplane::sai::EcnMarkMode>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_WRED_ATTR_ECN_GREEN_MIN_THRESHOLD:
        msg.set_ecn_green_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_GREEN_MAX_THRESHOLD:
        msg.set_ecn_green_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_GREEN_MARK_PROBABILITY:
        msg.set_ecn_green_mark_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_YELLOW_MIN_THRESHOLD:
        msg.set_ecn_yellow_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_YELLOW_MAX_THRESHOLD:
        msg.set_ecn_yellow_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_YELLOW_MARK_PROBABILITY:
        msg.set_ecn_yellow_mark_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_RED_MIN_THRESHOLD:
        msg.set_ecn_red_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_RED_MAX_THRESHOLD:
        msg.set_ecn_red_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_RED_MARK_PROBABILITY:
        msg.set_ecn_red_mark_probability(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MIN_THRESHOLD:
        msg.set_ecn_color_unaware_min_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MAX_THRESHOLD:
        msg.set_ecn_color_unaware_max_threshold(attr_list[i].value.u32);
        break;
      case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MARK_PROBABILITY:
        msg.set_ecn_color_unaware_mark_probability(attr_list[i].value.u32);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_wred(sai_object_id_t *wred_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateWredRequest req =
      convert_create_wred(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateWredResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = wred->CreateWred(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *wred_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_wred(sai_object_id_t wred_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveWredRequest req;
  lemming::dataplane::sai::RemoveWredResponse resp;
  grpc::ClientContext context;
  req.set_oid(wred_id);

  grpc::Status status = wred->RemoveWred(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_wred_attribute(sai_object_id_t wred_id,
                                  const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetWredAttributeRequest req;
  lemming::dataplane::sai::SetWredAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(wred_id);

  switch (attr->id) {
    case SAI_WRED_ATTR_GREEN_ENABLE:
      req.set_green_enable(attr->value.booldata);
      break;
    case SAI_WRED_ATTR_GREEN_MIN_THRESHOLD:
      req.set_green_min_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_GREEN_MAX_THRESHOLD:
      req.set_green_max_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_GREEN_DROP_PROBABILITY:
      req.set_green_drop_probability(attr->value.u32);
      break;
    case SAI_WRED_ATTR_YELLOW_ENABLE:
      req.set_yellow_enable(attr->value.booldata);
      break;
    case SAI_WRED_ATTR_YELLOW_MIN_THRESHOLD:
      req.set_yellow_min_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_YELLOW_MAX_THRESHOLD:
      req.set_yellow_max_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_YELLOW_DROP_PROBABILITY:
      req.set_yellow_drop_probability(attr->value.u32);
      break;
    case SAI_WRED_ATTR_RED_ENABLE:
      req.set_red_enable(attr->value.booldata);
      break;
    case SAI_WRED_ATTR_RED_MIN_THRESHOLD:
      req.set_red_min_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_RED_MAX_THRESHOLD:
      req.set_red_max_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_RED_DROP_PROBABILITY:
      req.set_red_drop_probability(attr->value.u32);
      break;
    case SAI_WRED_ATTR_WEIGHT:
      req.set_weight(attr->value.u8);
      break;
    case SAI_WRED_ATTR_ECN_MARK_MODE:
      req.set_ecn_mark_mode(static_cast<lemming::dataplane::sai::EcnMarkMode>(
          attr->value.s32 + 1));
      break;
    case SAI_WRED_ATTR_ECN_GREEN_MIN_THRESHOLD:
      req.set_ecn_green_min_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_GREEN_MAX_THRESHOLD:
      req.set_ecn_green_max_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_GREEN_MARK_PROBABILITY:
      req.set_ecn_green_mark_probability(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_YELLOW_MIN_THRESHOLD:
      req.set_ecn_yellow_min_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_YELLOW_MAX_THRESHOLD:
      req.set_ecn_yellow_max_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_YELLOW_MARK_PROBABILITY:
      req.set_ecn_yellow_mark_probability(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_RED_MIN_THRESHOLD:
      req.set_ecn_red_min_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_RED_MAX_THRESHOLD:
      req.set_ecn_red_max_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_RED_MARK_PROBABILITY:
      req.set_ecn_red_mark_probability(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MIN_THRESHOLD:
      req.set_ecn_color_unaware_min_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MAX_THRESHOLD:
      req.set_ecn_color_unaware_max_threshold(attr->value.u32);
      break;
    case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MARK_PROBABILITY:
      req.set_ecn_color_unaware_mark_probability(attr->value.u32);
      break;
  }

  grpc::Status status = wred->SetWredAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_wred_attribute(sai_object_id_t wred_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetWredAttributeRequest req;
  lemming::dataplane::sai::GetWredAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(wred_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::WredAttr>(attr_list[i].id + 1));
  }
  grpc::Status status = wred->GetWredAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_WRED_ATTR_GREEN_ENABLE:
        attr_list[i].value.booldata = resp.attr().green_enable();
        break;
      case SAI_WRED_ATTR_GREEN_MIN_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().green_min_threshold();
        break;
      case SAI_WRED_ATTR_GREEN_MAX_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().green_max_threshold();
        break;
      case SAI_WRED_ATTR_GREEN_DROP_PROBABILITY:
        attr_list[i].value.u32 = resp.attr().green_drop_probability();
        break;
      case SAI_WRED_ATTR_YELLOW_ENABLE:
        attr_list[i].value.booldata = resp.attr().yellow_enable();
        break;
      case SAI_WRED_ATTR_YELLOW_MIN_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().yellow_min_threshold();
        break;
      case SAI_WRED_ATTR_YELLOW_MAX_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().yellow_max_threshold();
        break;
      case SAI_WRED_ATTR_YELLOW_DROP_PROBABILITY:
        attr_list[i].value.u32 = resp.attr().yellow_drop_probability();
        break;
      case SAI_WRED_ATTR_RED_ENABLE:
        attr_list[i].value.booldata = resp.attr().red_enable();
        break;
      case SAI_WRED_ATTR_RED_MIN_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().red_min_threshold();
        break;
      case SAI_WRED_ATTR_RED_MAX_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().red_max_threshold();
        break;
      case SAI_WRED_ATTR_RED_DROP_PROBABILITY:
        attr_list[i].value.u32 = resp.attr().red_drop_probability();
        break;
      case SAI_WRED_ATTR_WEIGHT:
        attr_list[i].value.u8 = resp.attr().weight();
        break;
      case SAI_WRED_ATTR_ECN_MARK_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().ecn_mark_mode() - 1);
        break;
      case SAI_WRED_ATTR_ECN_GREEN_MIN_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().ecn_green_min_threshold();
        break;
      case SAI_WRED_ATTR_ECN_GREEN_MAX_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().ecn_green_max_threshold();
        break;
      case SAI_WRED_ATTR_ECN_GREEN_MARK_PROBABILITY:
        attr_list[i].value.u32 = resp.attr().ecn_green_mark_probability();
        break;
      case SAI_WRED_ATTR_ECN_YELLOW_MIN_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().ecn_yellow_min_threshold();
        break;
      case SAI_WRED_ATTR_ECN_YELLOW_MAX_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().ecn_yellow_max_threshold();
        break;
      case SAI_WRED_ATTR_ECN_YELLOW_MARK_PROBABILITY:
        attr_list[i].value.u32 = resp.attr().ecn_yellow_mark_probability();
        break;
      case SAI_WRED_ATTR_ECN_RED_MIN_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().ecn_red_min_threshold();
        break;
      case SAI_WRED_ATTR_ECN_RED_MAX_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().ecn_red_max_threshold();
        break;
      case SAI_WRED_ATTR_ECN_RED_MARK_PROBABILITY:
        attr_list[i].value.u32 = resp.attr().ecn_red_mark_probability();
        break;
      case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MIN_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().ecn_color_unaware_min_threshold();
        break;
      case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MAX_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().ecn_color_unaware_max_threshold();
        break;
      case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MARK_PROBABILITY:
        attr_list[i].value.u32 =
            resp.attr().ecn_color_unaware_mark_probability();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
