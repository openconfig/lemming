

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

#include "dataplane/standalone/sai/next_hop.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/next_hop.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_next_hop_api_t l_next_hop = {
    .create_next_hop = l_create_next_hop,
    .remove_next_hop = l_remove_next_hop,
    .set_next_hop_attribute = l_set_next_hop_attribute,
    .get_next_hop_attribute = l_get_next_hop_attribute,
};

sai_status_t l_create_next_hop(sai_object_id_t *next_hop_id,
                               sai_object_id_t switch_id, uint32_t attr_count,
                               const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateNextHopRequest req;
  lemming::dataplane::sai::CreateNextHopResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_NEXT_HOP_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::NextHopType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_NEXT_HOP_ATTR_IP:
        req.set_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_NEXT_HOP_ATTR_ROUTER_INTERFACE_ID:
        req.set_router_interface_id(attr_list[i].value.oid);
        break;
      case SAI_NEXT_HOP_ATTR_TUNNEL_ID:
        req.set_tunnel_id(attr_list[i].value.oid);
        break;
      case SAI_NEXT_HOP_ATTR_TUNNEL_VNI:
        req.set_tunnel_vni(attr_list[i].value.u32);
        break;
      case SAI_NEXT_HOP_ATTR_TUNNEL_MAC:
        req.set_tunnel_mac(attr_list[i].value.mac,
                           sizeof(attr_list[i].value.mac));
        break;
      case SAI_NEXT_HOP_ATTR_SRV6_SIDLIST_ID:
        req.set_srv6_sidlist_id(attr_list[i].value.oid);
        break;
      case SAI_NEXT_HOP_ATTR_LABELSTACK:
        req.mutable_labelstack()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_NEXT_HOP_ATTR_COUNTER_ID:
        req.set_counter_id(attr_list[i].value.oid);
        break;
      case SAI_NEXT_HOP_ATTR_DISABLE_DECREMENT_TTL:
        req.set_disable_decrement_ttl(attr_list[i].value.booldata);
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_TYPE:
        req.set_outseg_type(static_cast<lemming::dataplane::sai::OutsegType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_TTL_MODE:
        req.set_outseg_ttl_mode(
            static_cast<lemming::dataplane::sai::OutsegTtlMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_TTL_VALUE:
        req.set_outseg_ttl_value(attr_list[i].value.u8);
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_EXP_MODE:
        req.set_outseg_exp_mode(
            static_cast<lemming::dataplane::sai::OutsegExpMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_EXP_VALUE:
        req.set_outseg_exp_value(attr_list[i].value.u8);
        break;
      case SAI_NEXT_HOP_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
        req.set_qos_tc_and_color_to_mpls_exp_map(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = next_hop->CreateNextHop(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *next_hop_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_NEXT_HOP, next_hop_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_next_hop(sai_object_id_t next_hop_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_NEXT_HOP, next_hop_id);
}

sai_status_t l_set_next_hop_attribute(sai_object_id_t next_hop_id,
                                      const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_NEXT_HOP, next_hop_id, attr);
}

sai_status_t l_get_next_hop_attribute(sai_object_id_t next_hop_id,
                                      uint32_t attr_count,
                                      sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_NEXT_HOP, next_hop_id,
                                   attr_count, attr_list);
}
