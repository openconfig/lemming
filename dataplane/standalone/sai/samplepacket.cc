

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

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/samplepacket.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_samplepacket_api_t l_samplepacket = {
    .create_samplepacket = l_create_samplepacket,
    .remove_samplepacket = l_remove_samplepacket,
    .set_samplepacket_attribute = l_set_samplepacket_attribute,
    .get_samplepacket_attribute = l_get_samplepacket_attribute,
};

sai_status_t l_create_samplepacket(sai_object_id_t *samplepacket_id,
                                   sai_object_id_t switch_id,
                                   uint32_t attr_count,
                                   const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateSamplepacketRequest req;
  lemming::dataplane::sai::CreateSamplepacketResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SAMPLEPACKET_ATTR_SAMPLE_RATE:
        req.set_sample_rate(attr_list[i].value.u32);
        break;
      case SAI_SAMPLEPACKET_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::SamplepacketType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_SAMPLEPACKET_ATTR_MODE:
        req.set_mode(static_cast<lemming::dataplane::sai::SamplepacketMode>(
            attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = samplepacket->CreateSamplepacket(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *samplepacket_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_SAMPLEPACKET, samplepacket_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_samplepacket(sai_object_id_t samplepacket_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_SAMPLEPACKET, samplepacket_id);
}

sai_status_t l_set_samplepacket_attribute(sai_object_id_t samplepacket_id,
                                          const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_SAMPLEPACKET,
                                   samplepacket_id, attr);
}

sai_status_t l_get_samplepacket_attribute(sai_object_id_t samplepacket_id,
                                          uint32_t attr_count,
                                          sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_SAMPLEPACKET,
                                   samplepacket_id, attr_count, attr_list);
}
