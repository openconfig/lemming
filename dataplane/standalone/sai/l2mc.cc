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

#include "dataplane/standalone/sai/l2mc.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/l2mc.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_l2mc_api_t l_l2mc = {
    .create_l2mc_entry = l_create_l2mc_entry,
    .remove_l2mc_entry = l_remove_l2mc_entry,
    .set_l2mc_entry_attribute = l_set_l2mc_entry_attribute,
    .get_l2mc_entry_attribute = l_get_l2mc_entry_attribute,
};

lemming::dataplane::sai::CreateL2mcEntryRequest convert_create_l2mc_entry(
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateL2mcEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_L2MC_ENTRY_ATTR_PACKET_ACTION:
        msg.set_packet_action(
            convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_L2MC_ENTRY_ATTR_OUTPUT_GROUP_ID:
        msg.set_output_group_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_l2mc_entry(const sai_l2mc_entry_t *l2mc_entry,
                                 uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateL2mcEntryRequest req =
      convert_create_l2mc_entry(attr_count, attr_list);
  lemming::dataplane::sai::CreateL2mcEntryResponse resp;
  grpc::ClientContext context;

  grpc::Status status = l2mc->CreateL2mcEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_l2mc_entry(const sai_l2mc_entry_t *l2mc_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveL2mcEntryRequest req;
  lemming::dataplane::sai::RemoveL2mcEntryResponse resp;
  grpc::ClientContext context;

  grpc::Status status = l2mc->RemoveL2mcEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_l2mc_entry_attribute(const sai_l2mc_entry_t *l2mc_entry,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetL2mcEntryAttributeRequest req;
  lemming::dataplane::sai::SetL2mcEntryAttributeResponse resp;
  grpc::ClientContext context;

  switch (attr->id) {
    case SAI_L2MC_ENTRY_ATTR_PACKET_ACTION:
      req.set_packet_action(
          convert_sai_packet_action_t_to_proto(attr->value.s32));
      break;
    case SAI_L2MC_ENTRY_ATTR_OUTPUT_GROUP_ID:
      req.set_output_group_id(attr->value.oid);
      break;
  }

  grpc::Status status = l2mc->SetL2mcEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_l2mc_entry_attribute(const sai_l2mc_entry_t *l2mc_entry,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetL2mcEntryAttributeRequest req;
  lemming::dataplane::sai::GetL2mcEntryAttributeResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_l2mc_entry_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = l2mc->GetL2mcEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_L2MC_ENTRY_ATTR_PACKET_ACTION:
        attr_list[i].value.s32 =
            convert_sai_packet_action_t_to_sai(resp.attr().packet_action());
        break;
      case SAI_L2MC_ENTRY_ATTR_OUTPUT_GROUP_ID:
        attr_list[i].value.oid = resp.attr().output_group_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
