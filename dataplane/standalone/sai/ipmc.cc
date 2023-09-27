

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

#include "dataplane/standalone/sai/ipmc.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/ipmc.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_ipmc_api_t l_ipmc = {
    .create_ipmc_entry = l_create_ipmc_entry,
    .remove_ipmc_entry = l_remove_ipmc_entry,
    .set_ipmc_entry_attribute = l_set_ipmc_entry_attribute,
    .get_ipmc_entry_attribute = l_get_ipmc_entry_attribute,
};

lemming::dataplane::sai::CreateIpmcEntryRequest convert_create_ipmc_entry(
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateIpmcEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_IPMC_ENTRY_ATTR_PACKET_ACTION:
        msg.set_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_IPMC_ENTRY_ATTR_OUTPUT_GROUP_ID:
        msg.set_output_group_id(attr_list[i].value.oid);
        break;
      case SAI_IPMC_ENTRY_ATTR_RPF_GROUP_ID:
        msg.set_rpf_group_id(attr_list[i].value.oid);
        break;
      case SAI_IPMC_ENTRY_ATTR_COUNTER_ID:
        msg.set_counter_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_ipmc_entry(const sai_ipmc_entry_t *ipmc_entry,
                                 uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIpmcEntryRequest req =
      convert_create_ipmc_entry(attr_count, attr_list);
  lemming::dataplane::sai::CreateIpmcEntryResponse resp;
  grpc::ClientContext context;

  grpc::Status status = ipmc->CreateIpmcEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_ipmc_entry(const sai_ipmc_entry_t *ipmc_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveIpmcEntryRequest req;
  lemming::dataplane::sai::RemoveIpmcEntryResponse resp;
  grpc::ClientContext context;

  grpc::Status status = ipmc->RemoveIpmcEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_ipmc_entry_attribute(const sai_ipmc_entry_t *ipmc_entry,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetIpmcEntryAttributeRequest req;
  lemming::dataplane::sai::SetIpmcEntryAttributeResponse resp;
  grpc::ClientContext context;

  switch (attr->id) {
    case SAI_IPMC_ENTRY_ATTR_PACKET_ACTION:
      req.set_packet_action(static_cast<lemming::dataplane::sai::PacketAction>(
          attr->value.s32 + 1));
      break;
    case SAI_IPMC_ENTRY_ATTR_OUTPUT_GROUP_ID:
      req.set_output_group_id(attr->value.oid);
      break;
    case SAI_IPMC_ENTRY_ATTR_RPF_GROUP_ID:
      req.set_rpf_group_id(attr->value.oid);
      break;
    case SAI_IPMC_ENTRY_ATTR_COUNTER_ID:
      req.set_counter_id(attr->value.oid);
      break;
  }

  grpc::Status status = ipmc->SetIpmcEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_ipmc_entry_attribute(const sai_ipmc_entry_t *ipmc_entry,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIpmcEntryAttributeRequest req;
  lemming::dataplane::sai::GetIpmcEntryAttributeResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::IpmcEntryAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = ipmc->GetIpmcEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_IPMC_ENTRY_ATTR_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().packet_action() - 1);
        break;
      case SAI_IPMC_ENTRY_ATTR_OUTPUT_GROUP_ID:
        attr_list[i].value.oid = resp.attr().output_group_id();
        break;
      case SAI_IPMC_ENTRY_ATTR_RPF_GROUP_ID:
        attr_list[i].value.oid = resp.attr().rpf_group_id();
        break;
      case SAI_IPMC_ENTRY_ATTR_COUNTER_ID:
        attr_list[i].value.oid = resp.attr().counter_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
