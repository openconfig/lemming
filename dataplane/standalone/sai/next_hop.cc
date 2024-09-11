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

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/next_hop.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_next_hop_api_t l_next_hop = {
    .create_next_hop = l_create_next_hop,
    .remove_next_hop = l_remove_next_hop,
    .set_next_hop_attribute = l_set_next_hop_attribute,
    .get_next_hop_attribute = l_get_next_hop_attribute,
    .create_next_hops = l_create_next_hops,
    .remove_next_hops = l_remove_next_hops,
    .set_next_hops_attribute = l_set_next_hops_attribute,
    .get_next_hops_attribute = l_get_next_hops_attribute,
};

lemming::dataplane::sai::CreateNextHopRequest convert_create_next_hop(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateNextHopRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_NEXT_HOP_ATTR_TYPE:
        msg.set_type(
            convert_sai_next_hop_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_NEXT_HOP_ATTR_IP:
        msg.set_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_NEXT_HOP_ATTR_ROUTER_INTERFACE_ID:
        msg.set_router_interface_id(attr_list[i].value.oid);
        break;
      case SAI_NEXT_HOP_ATTR_TUNNEL_ID:
        msg.set_tunnel_id(attr_list[i].value.oid);
        break;
      case SAI_NEXT_HOP_ATTR_TUNNEL_VNI:
        msg.set_tunnel_vni(attr_list[i].value.u32);
        break;
      case SAI_NEXT_HOP_ATTR_TUNNEL_MAC:
        msg.set_tunnel_mac(attr_list[i].value.mac,
                           sizeof(attr_list[i].value.mac));
        break;
      case SAI_NEXT_HOP_ATTR_SRV6_SIDLIST_ID:
        msg.set_srv6_sidlist_id(attr_list[i].value.oid);
        break;
      case SAI_NEXT_HOP_ATTR_LABELSTACK:
        msg.mutable_labelstack()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_NEXT_HOP_ATTR_COUNTER_ID:
        msg.set_counter_id(attr_list[i].value.oid);
        break;
      case SAI_NEXT_HOP_ATTR_DISABLE_DECREMENT_TTL:
        msg.set_disable_decrement_ttl(attr_list[i].value.booldata);
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_TYPE:
        msg.set_outseg_type(
            convert_sai_outseg_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_TTL_MODE:
        msg.set_outseg_ttl_mode(
            convert_sai_outseg_ttl_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_TTL_VALUE:
        msg.set_outseg_ttl_value(attr_list[i].value.u8);
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_EXP_MODE:
        msg.set_outseg_exp_mode(
            convert_sai_outseg_exp_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_EXP_VALUE:
        msg.set_outseg_exp_value(attr_list[i].value.u8);
        break;
      case SAI_NEXT_HOP_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
        msg.set_qos_tc_and_color_to_mpls_exp_map(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_next_hop(sai_object_id_t *next_hop_id,
                               sai_object_id_t switch_id, uint32_t attr_count,
                               const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateNextHopRequest req =
      convert_create_next_hop(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateNextHopResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = next_hop->CreateNextHop(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (next_hop_id) {
    *next_hop_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_next_hop(sai_object_id_t next_hop_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveNextHopRequest req;
  lemming::dataplane::sai::RemoveNextHopResponse resp;
  grpc::ClientContext context;
  req.set_oid(next_hop_id);

  grpc::Status status = next_hop->RemoveNextHop(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_next_hop_attribute(sai_object_id_t next_hop_id,
                                      const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetNextHopAttributeRequest req;
  lemming::dataplane::sai::SetNextHopAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(next_hop_id);

  switch (attr->id) {
    case SAI_NEXT_HOP_ATTR_TUNNEL_VNI:
      req.set_tunnel_vni(attr->value.u32);
      break;
    case SAI_NEXT_HOP_ATTR_TUNNEL_MAC:
      req.set_tunnel_mac(attr->value.mac, sizeof(attr->value.mac));
      break;
    case SAI_NEXT_HOP_ATTR_COUNTER_ID:
      req.set_counter_id(attr->value.oid);
      break;
    case SAI_NEXT_HOP_ATTR_DISABLE_DECREMENT_TTL:
      req.set_disable_decrement_ttl(attr->value.booldata);
      break;
    case SAI_NEXT_HOP_ATTR_OUTSEG_TYPE:
      req.set_outseg_type(convert_sai_outseg_type_t_to_proto(attr->value.s32));
      break;
    case SAI_NEXT_HOP_ATTR_OUTSEG_TTL_MODE:
      req.set_outseg_ttl_mode(
          convert_sai_outseg_ttl_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_NEXT_HOP_ATTR_OUTSEG_TTL_VALUE:
      req.set_outseg_ttl_value(attr->value.u8);
      break;
    case SAI_NEXT_HOP_ATTR_OUTSEG_EXP_MODE:
      req.set_outseg_exp_mode(
          convert_sai_outseg_exp_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_NEXT_HOP_ATTR_OUTSEG_EXP_VALUE:
      req.set_outseg_exp_value(attr->value.u8);
      break;
    case SAI_NEXT_HOP_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
      req.set_qos_tc_and_color_to_mpls_exp_map(attr->value.oid);
      break;
  }

  grpc::Status status = next_hop->SetNextHopAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_next_hop_attribute(sai_object_id_t next_hop_id,
                                      uint32_t attr_count,
                                      sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetNextHopAttributeRequest req;
  lemming::dataplane::sai::GetNextHopAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(next_hop_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_next_hop_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = next_hop->GetNextHopAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_NEXT_HOP_ATTR_TYPE:
        attr_list[i].value.s32 =
            convert_sai_next_hop_type_t_to_sai(resp.attr().type());
        break;
      case SAI_NEXT_HOP_ATTR_IP:
        attr_list[i].value.ipaddr = convert_to_ip_address(resp.attr().ip());
        break;
      case SAI_NEXT_HOP_ATTR_ROUTER_INTERFACE_ID:
        attr_list[i].value.oid = resp.attr().router_interface_id();
        break;
      case SAI_NEXT_HOP_ATTR_TUNNEL_ID:
        attr_list[i].value.oid = resp.attr().tunnel_id();
        break;
      case SAI_NEXT_HOP_ATTR_TUNNEL_VNI:
        attr_list[i].value.u32 = resp.attr().tunnel_vni();
        break;
      case SAI_NEXT_HOP_ATTR_TUNNEL_MAC:
        memcpy(attr_list[i].value.mac, resp.attr().tunnel_mac().data(),
               sizeof(sai_mac_t));
        break;
      case SAI_NEXT_HOP_ATTR_SRV6_SIDLIST_ID:
        attr_list[i].value.oid = resp.attr().srv6_sidlist_id();
        break;
      case SAI_NEXT_HOP_ATTR_LABELSTACK:
        copy_list(attr_list[i].value.u32list.list, resp.attr().labelstack(),
                  &attr_list[i].value.u32list.count);
        break;
      case SAI_NEXT_HOP_ATTR_COUNTER_ID:
        attr_list[i].value.oid = resp.attr().counter_id();
        break;
      case SAI_NEXT_HOP_ATTR_DISABLE_DECREMENT_TTL:
        attr_list[i].value.booldata = resp.attr().disable_decrement_ttl();
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_TYPE:
        attr_list[i].value.s32 =
            convert_sai_outseg_type_t_to_sai(resp.attr().outseg_type());
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_TTL_MODE:
        attr_list[i].value.s32 =
            convert_sai_outseg_ttl_mode_t_to_sai(resp.attr().outseg_ttl_mode());
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_TTL_VALUE:
        attr_list[i].value.u8 = resp.attr().outseg_ttl_value();
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_EXP_MODE:
        attr_list[i].value.s32 =
            convert_sai_outseg_exp_mode_t_to_sai(resp.attr().outseg_exp_mode());
        break;
      case SAI_NEXT_HOP_ATTR_OUTSEG_EXP_VALUE:
        attr_list[i].value.u8 = resp.attr().outseg_exp_value();
        break;
      case SAI_NEXT_HOP_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
        attr_list[i].value.oid = resp.attr().qos_tc_and_color_to_mpls_exp_map();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_next_hops(sai_object_id_t switch_id,
                                uint32_t object_count,
                                const uint32_t *attr_count,
                                const sai_attribute_t **attr_list,
                                sai_bulk_op_error_mode_t mode,
                                sai_object_id_t *object_id,
                                sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateNextHopsRequest req;
  lemming::dataplane::sai::CreateNextHopsResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    auto r = convert_create_next_hop(switch_id, attr_count[i], attr_list[i]);
    *req.add_reqs() = r;
  }

  grpc::Status status = next_hop->CreateNextHops(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (object_count != resp.resps().size()) {
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < object_count; i++) {
    object_id[i] = resp.resps(i).oid();
    object_statuses[i] = SAI_STATUS_SUCCESS;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_next_hops(uint32_t object_count,
                                const sai_object_id_t *object_id,
                                sai_bulk_op_error_mode_t mode,
                                sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveNextHopsRequest req;
  lemming::dataplane::sai::RemoveNextHopsResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    req.add_reqs()->set_oid(object_id[i]);
  }

  grpc::Status status = next_hop->RemoveNextHops(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (object_count != resp.resps().size()) {
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < object_count; i++) {
    object_statuses[i] = SAI_STATUS_SUCCESS;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_next_hops_attribute(uint32_t object_count,
                                       const sai_object_id_t *object_id,
                                       const sai_attribute_t *attr_list,
                                       sai_bulk_op_error_mode_t mode,
                                       sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_next_hops_attribute(uint32_t object_count,
                                       const sai_object_id_t *object_id,
                                       const uint32_t *attr_count,
                                       sai_attribute_t **attr_list,
                                       sai_bulk_op_error_mode_t mode,
                                       sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}
