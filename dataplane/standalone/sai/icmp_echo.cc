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

#include "dataplane/standalone/sai/icmp_echo.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/icmp_echo.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_icmp_echo_api_t l_icmp_echo = {
    .create_icmp_echo_session = l_create_icmp_echo_session,
    .remove_icmp_echo_session = l_remove_icmp_echo_session,
    .set_icmp_echo_session_attribute = l_set_icmp_echo_session_attribute,
    .get_icmp_echo_session_attribute = l_get_icmp_echo_session_attribute,
    .get_icmp_echo_session_stats = l_get_icmp_echo_session_stats,
    .get_icmp_echo_session_stats_ext = l_get_icmp_echo_session_stats_ext,
    .clear_icmp_echo_session_stats = l_clear_icmp_echo_session_stats,
};

lemming::dataplane::sai::CreateIcmpEchoSessionRequest
convert_create_icmp_echo_session(sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t* attr_list) {
  lemming::dataplane::sai::CreateIcmpEchoSessionRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ICMP_ECHO_SESSION_ATTR_HW_LOOKUP_VALID:
        msg.set_hw_lookup_valid(attr_list[i].value.booldata);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_VIRTUAL_ROUTER:
        msg.set_virtual_router(attr_list[i].value.oid);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_PORT:
        msg.set_port(attr_list[i].value.oid);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_RX_PORT:
        msg.set_rx_port(attr_list[i].value.oid);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_GUID:
        msg.set_guid(attr_list[i].value.u64);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_COOKIE:
        msg.set_cookie(attr_list[i].value.u32);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_IPHDR_VERSION:
        msg.set_iphdr_version(attr_list[i].value.u8);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_TOS:
        msg.set_tos(attr_list[i].value.u8);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_TTL:
        msg.set_ttl(attr_list[i].value.u8);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_SRC_IP_ADDRESS:
        msg.set_src_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_DST_IP_ADDRESS:
        msg.set_dst_ip_address(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_SRC_MAC_ADDRESS:
        msg.set_src_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_DST_MAC_ADDRESS:
        msg.set_dst_mac_address(attr_list[i].value.mac,
                                sizeof(attr_list[i].value.mac));
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_TX_INTERVAL:
        msg.set_tx_interval(attr_list[i].value.u32);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_RX_INTERVAL:
        msg.set_rx_interval(attr_list[i].value.u32);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_SET_NEXT_HOP_GROUP_SWITCHOVER:
        msg.set_set_next_hop_group_switchover(attr_list[i].value.booldata);
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_STATS_COUNT_MODE:
        msg.set_stats_count_mode(
            convert_sai_stats_count_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_SELECTIVE_COUNTER_LIST:
        msg.mutable_selective_counter_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_icmp_echo_session(sai_object_id_t* icmp_echo_session_id,
                                        sai_object_id_t switch_id,
                                        uint32_t attr_count,
                                        const sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIcmpEchoSessionRequest req =
      convert_create_icmp_echo_session(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateIcmpEchoSessionResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = icmp_echo->CreateIcmpEchoSession(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (icmp_echo_session_id) {
    *icmp_echo_session_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_icmp_echo_session(sai_object_id_t icmp_echo_session_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveIcmpEchoSessionRequest req;
  lemming::dataplane::sai::RemoveIcmpEchoSessionResponse resp;
  grpc::ClientContext context;
  req.set_oid(icmp_echo_session_id);

  grpc::Status status = icmp_echo->RemoveIcmpEchoSession(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_icmp_echo_session_attribute(
    sai_object_id_t icmp_echo_session_id, const sai_attribute_t* attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetIcmpEchoSessionAttributeRequest req;
  lemming::dataplane::sai::SetIcmpEchoSessionAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(icmp_echo_session_id);

  switch (attr->id) {
    case SAI_ICMP_ECHO_SESSION_ATTR_VIRTUAL_ROUTER:
      req.set_virtual_router(attr->value.oid);
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_PORT:
      req.set_port(attr->value.oid);
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_RX_PORT:
      req.set_rx_port(attr->value.oid);
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_IPHDR_VERSION:
      req.set_iphdr_version(attr->value.u8);
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_TOS:
      req.set_tos(attr->value.u8);
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_TTL:
      req.set_ttl(attr->value.u8);
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_SRC_MAC_ADDRESS:
      req.set_src_mac_address(attr->value.mac, sizeof(attr->value.mac));
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_DST_MAC_ADDRESS:
      req.set_dst_mac_address(attr->value.mac, sizeof(attr->value.mac));
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_TX_INTERVAL:
      req.set_tx_interval(attr->value.u32);
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_RX_INTERVAL:
      req.set_rx_interval(attr->value.u32);
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_SET_NEXT_HOP_GROUP_SWITCHOVER:
      req.set_set_next_hop_group_switchover(attr->value.booldata);
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_STATS_COUNT_MODE:
      req.set_stats_count_mode(
          convert_sai_stats_count_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_ICMP_ECHO_SESSION_ATTR_SELECTIVE_COUNTER_LIST:
      req.mutable_selective_counter_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
  }

  grpc::Status status =
      icmp_echo->SetIcmpEchoSessionAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_icmp_echo_session_attribute(
    sai_object_id_t icmp_echo_session_id, uint32_t attr_count,
    sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIcmpEchoSessionAttributeRequest req;
  lemming::dataplane::sai::GetIcmpEchoSessionAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(icmp_echo_session_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_icmp_echo_session_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status =
      icmp_echo->GetIcmpEchoSessionAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ICMP_ECHO_SESSION_ATTR_HW_LOOKUP_VALID:
        attr_list[i].value.booldata = resp.attr().hw_lookup_valid();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_VIRTUAL_ROUTER:
        attr_list[i].value.oid = resp.attr().virtual_router();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_PORT:
        attr_list[i].value.oid = resp.attr().port();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_RX_PORT:
        attr_list[i].value.oid = resp.attr().rx_port();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_GUID:
        attr_list[i].value.u64 = resp.attr().guid();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_COOKIE:
        attr_list[i].value.u32 = resp.attr().cookie();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_IPHDR_VERSION:
        attr_list[i].value.u8 = resp.attr().iphdr_version();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_TOS:
        attr_list[i].value.u8 = resp.attr().tos();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_TTL:
        attr_list[i].value.u8 = resp.attr().ttl();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_SRC_IP_ADDRESS:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().src_ip_address());
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_DST_IP_ADDRESS:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().dst_ip_address());
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_SRC_MAC_ADDRESS:
        memcpy(attr_list[i].value.mac, resp.attr().src_mac_address().data(),
               sizeof(sai_mac_t));
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_DST_MAC_ADDRESS:
        memcpy(attr_list[i].value.mac, resp.attr().dst_mac_address().data(),
               sizeof(sai_mac_t));
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_TX_INTERVAL:
        attr_list[i].value.u32 = resp.attr().tx_interval();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_RX_INTERVAL:
        attr_list[i].value.u32 = resp.attr().rx_interval();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_SET_NEXT_HOP_GROUP_SWITCHOVER:
        attr_list[i].value.booldata =
            resp.attr().set_next_hop_group_switchover();
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_STATE:
        attr_list[i].value.s32 =
            convert_sai_icmp_echo_session_state_t_to_sai(resp.attr().state());
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_STATS_COUNT_MODE:
        attr_list[i].value.s32 = convert_sai_stats_count_mode_t_to_sai(
            resp.attr().stats_count_mode());
        break;
      case SAI_ICMP_ECHO_SESSION_ATTR_SELECTIVE_COUNTER_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().selective_counter_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_icmp_echo_session_stats(sai_object_id_t icmp_echo_session_id,
                                           uint32_t number_of_counters,
                                           const sai_stat_id_t* counter_ids,
                                           uint64_t* counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIcmpEchoSessionStatsRequest req;
  lemming::dataplane::sai::GetIcmpEchoSessionStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(icmp_echo_session_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(
        convert_sai_icmp_echo_session_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status =
      icmp_echo->GetIcmpEchoSessionStats(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_icmp_echo_session_stats_ext(
    sai_object_id_t icmp_echo_session_id, uint32_t number_of_counters,
    const sai_stat_id_t* counter_ids, sai_stats_mode_t mode,
    uint64_t* counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_icmp_echo_session_stats(
    sai_object_id_t icmp_echo_session_id, uint32_t number_of_counters,
    const sai_stat_id_t* counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
