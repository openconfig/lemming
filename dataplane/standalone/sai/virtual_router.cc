

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

#include "dataplane/standalone/sai/virtual_router.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/virtual_router.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_virtual_router_api_t l_virtual_router = {
    .create_virtual_router = l_create_virtual_router,
    .remove_virtual_router = l_remove_virtual_router,
    .set_virtual_router_attribute = l_set_virtual_router_attribute,
    .get_virtual_router_attribute = l_get_virtual_router_attribute,
};

sai_status_t l_create_virtual_router(sai_object_id_t *virtual_router_id,
                                     sai_object_id_t switch_id,
                                     uint32_t attr_count,
                                     const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateVirtualRouterRequest req;
  lemming::dataplane::sai::CreateVirtualRouterResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_VIRTUAL_ROUTER_ATTR_ADMIN_V4_STATE:
        req.set_admin_v4_state(attr_list[i].value.booldata);
        break;
      case SAI_VIRTUAL_ROUTER_ATTR_ADMIN_V6_STATE:
        req.set_admin_v6_state(attr_list[i].value.booldata);
        break;
      case SAI_VIRTUAL_ROUTER_ATTR_SRC_MAC_ADDRESS:
        req.set_src_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_VIRTUAL_ROUTER_ATTR_VIOLATION_TTL1_PACKET_ACTION:
        req.set_violation_ttl1_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_VIRTUAL_ROUTER_ATTR_VIOLATION_IP_OPTIONS_PACKET_ACTION:
        req.set_violation_ip_options_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_VIRTUAL_ROUTER_ATTR_UNKNOWN_L3_MULTICAST_PACKET_ACTION:
        req.set_unknown_l3_multicast_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_VIRTUAL_ROUTER_ATTR_LABEL:
        req.set_label(attr_list[i].value.chardata);
        break;
    }
  }
  grpc::Status status =
      virtual_router->CreateVirtualRouter(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *virtual_router_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_VIRTUAL_ROUTER, virtual_router_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_virtual_router(sai_object_id_t virtual_router_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_VIRTUAL_ROUTER, virtual_router_id);
}

sai_status_t l_set_virtual_router_attribute(sai_object_id_t virtual_router_id,
                                            const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_VIRTUAL_ROUTER,
                                   virtual_router_id, attr);
}

sai_status_t l_get_virtual_router_attribute(sai_object_id_t virtual_router_id,
                                            uint32_t attr_count,
                                            sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_VIRTUAL_ROUTER,
                                   virtual_router_id, attr_count, attr_list);
}
