

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

#include "dataplane/standalone/sai/samplepacket.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/samplepacket.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_samplepacket_api_t l_samplepacket = {
    .create_samplepacket = l_create_samplepacket,
    .remove_samplepacket = l_remove_samplepacket,
    .set_samplepacket_attribute = l_set_samplepacket_attribute,
    .get_samplepacket_attribute = l_get_samplepacket_attribute,
};

lemming::dataplane::sai::CreateSamplepacketRequest convert_create_samplepacket(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateSamplepacketRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SAMPLEPACKET_ATTR_SAMPLE_RATE:
        msg.set_sample_rate(attr_list[i].value.u32);
        break;
      case SAI_SAMPLEPACKET_ATTR_TYPE:
        msg.set_type(static_cast<lemming::dataplane::sai::SamplepacketType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_SAMPLEPACKET_ATTR_MODE:
        msg.set_mode(static_cast<lemming::dataplane::sai::SamplepacketMode>(
            attr_list[i].value.s32 + 1));
        break;
    }
  }
  return msg;
}

sai_status_t l_create_samplepacket(sai_object_id_t *samplepacket_id,
                                   sai_object_id_t switch_id,
                                   uint32_t attr_count,
                                   const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateSamplepacketRequest req =
      convert_create_samplepacket(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateSamplepacketResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = samplepacket->CreateSamplepacket(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *samplepacket_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_samplepacket(sai_object_id_t samplepacket_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveSamplepacketRequest req;
  lemming::dataplane::sai::RemoveSamplepacketResponse resp;
  grpc::ClientContext context;
  req.set_oid(samplepacket_id);

  grpc::Status status = samplepacket->RemoveSamplepacket(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_samplepacket_attribute(sai_object_id_t samplepacket_id,
                                          const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetSamplepacketAttributeRequest req;
  lemming::dataplane::sai::SetSamplepacketAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(samplepacket_id);

  switch (attr->id) {
    case SAI_SAMPLEPACKET_ATTR_SAMPLE_RATE:
      req.set_sample_rate(attr->value.u32);
      break;
  }

  grpc::Status status =
      samplepacket->SetSamplepacketAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_samplepacket_attribute(sai_object_id_t samplepacket_id,
                                          uint32_t attr_count,
                                          sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetSamplepacketAttributeRequest req;
  lemming::dataplane::sai::GetSamplepacketAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(samplepacket_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::SamplepacketAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status =
      samplepacket->GetSamplepacketAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SAMPLEPACKET_ATTR_SAMPLE_RATE:
        attr_list[i].value.u32 = resp.attr().sample_rate();
        break;
      case SAI_SAMPLEPACKET_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_SAMPLEPACKET_ATTR_MODE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().mode() - 1);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
