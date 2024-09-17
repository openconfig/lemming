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

#include "dataplane/standalone/sai/bridge.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/bridge.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_bridge_api_t l_bridge = {
    .create_bridge = l_create_bridge,
    .remove_bridge = l_remove_bridge,
    .set_bridge_attribute = l_set_bridge_attribute,
    .get_bridge_attribute = l_get_bridge_attribute,
    .get_bridge_stats = l_get_bridge_stats,
    .get_bridge_stats_ext = l_get_bridge_stats_ext,
    .clear_bridge_stats = l_clear_bridge_stats,
    .create_bridge_port = l_create_bridge_port,
    .remove_bridge_port = l_remove_bridge_port,
    .set_bridge_port_attribute = l_set_bridge_port_attribute,
    .get_bridge_port_attribute = l_get_bridge_port_attribute,
    .get_bridge_port_stats = l_get_bridge_port_stats,
    .get_bridge_port_stats_ext = l_get_bridge_port_stats_ext,
    .clear_bridge_port_stats = l_clear_bridge_port_stats,
};

lemming::dataplane::sai::CreateBridgeRequest convert_create_bridge(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateBridgeRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BRIDGE_ATTR_TYPE:
        msg.set_type(
            convert_sai_bridge_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_BRIDGE_ATTR_MAX_LEARNED_ADDRESSES:
        msg.set_max_learned_addresses(attr_list[i].value.u32);
        break;
      case SAI_BRIDGE_ATTR_LEARN_DISABLE:
        msg.set_learn_disable(attr_list[i].value.booldata);
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
        msg.set_unknown_unicast_flood_control_type(
            convert_sai_bridge_flood_control_type_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
        msg.set_unknown_unicast_flood_group(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
        msg.set_unknown_multicast_flood_control_type(
            convert_sai_bridge_flood_control_type_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
        msg.set_unknown_multicast_flood_group(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
        msg.set_broadcast_flood_control_type(
            convert_sai_bridge_flood_control_type_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_BRIDGE_ATTR_BROADCAST_FLOOD_GROUP:
        msg.set_broadcast_flood_group(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateBridgePortRequest convert_create_bridge_port(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateBridgePortRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BRIDGE_PORT_ATTR_TYPE:
        msg.set_type(
            convert_sai_bridge_port_type_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_BRIDGE_PORT_ATTR_PORT_ID:
        msg.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_PORT_ATTR_TAGGING_MODE:
        msg.set_tagging_mode(convert_sai_bridge_port_tagging_mode_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_BRIDGE_PORT_ATTR_VLAN_ID:
        msg.set_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_BRIDGE_PORT_ATTR_RIF_ID:
        msg.set_rif_id(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_PORT_ATTR_TUNNEL_ID:
        msg.set_tunnel_id(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_PORT_ATTR_BRIDGE_ID:
        msg.set_bridge_id(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_MODE:
        msg.set_fdb_learning_mode(
            convert_sai_bridge_port_fdb_learning_mode_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_BRIDGE_PORT_ATTR_MAX_LEARNED_ADDRESSES:
        msg.set_max_learned_addresses(attr_list[i].value.u32);
        break;
      case SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_LIMIT_VIOLATION_PACKET_ACTION:
        msg.set_fdb_learning_limit_violation_packet_action(
            convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_BRIDGE_PORT_ATTR_ADMIN_STATE:
        msg.set_admin_state(attr_list[i].value.booldata);
        break;
      case SAI_BRIDGE_PORT_ATTR_INGRESS_FILTERING:
        msg.set_ingress_filtering(attr_list[i].value.booldata);
        break;
      case SAI_BRIDGE_PORT_ATTR_EGRESS_FILTERING:
        msg.set_egress_filtering(attr_list[i].value.booldata);
        break;
      case SAI_BRIDGE_PORT_ATTR_ISOLATION_GROUP:
        msg.set_isolation_group(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_bridge(sai_object_id_t *bridge_id,
                             sai_object_id_t switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateBridgeRequest req =
      convert_create_bridge(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateBridgeResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = bridge->CreateBridge(&context, req, &resp);
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
  if (bridge_id) {
    *bridge_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_bridge(sai_object_id_t bridge_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveBridgeRequest req;
  lemming::dataplane::sai::RemoveBridgeResponse resp;
  grpc::ClientContext context;
  req.set_oid(bridge_id);

  grpc::Status status = bridge->RemoveBridge(&context, req, &resp);
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

sai_status_t l_set_bridge_attribute(sai_object_id_t bridge_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetBridgeAttributeRequest req;
  lemming::dataplane::sai::SetBridgeAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(bridge_id);

  switch (attr->id) {
    case SAI_BRIDGE_ATTR_MAX_LEARNED_ADDRESSES:
      req.set_max_learned_addresses(attr->value.u32);
      break;
    case SAI_BRIDGE_ATTR_LEARN_DISABLE:
      req.set_learn_disable(attr->value.booldata);
      break;
    case SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
      req.set_unknown_unicast_flood_control_type(
          convert_sai_bridge_flood_control_type_t_to_proto(attr->value.s32));
      break;
    case SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
      req.set_unknown_unicast_flood_group(attr->value.oid);
      break;
    case SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
      req.set_unknown_multicast_flood_control_type(
          convert_sai_bridge_flood_control_type_t_to_proto(attr->value.s32));
      break;
    case SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
      req.set_unknown_multicast_flood_group(attr->value.oid);
      break;
    case SAI_BRIDGE_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
      req.set_broadcast_flood_control_type(
          convert_sai_bridge_flood_control_type_t_to_proto(attr->value.s32));
      break;
    case SAI_BRIDGE_ATTR_BROADCAST_FLOOD_GROUP:
      req.set_broadcast_flood_group(attr->value.oid);
      break;
  }

  grpc::Status status = bridge->SetBridgeAttribute(&context, req, &resp);
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

sai_status_t l_get_bridge_attribute(sai_object_id_t bridge_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetBridgeAttributeRequest req;
  lemming::dataplane::sai::GetBridgeAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(bridge_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_bridge_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = bridge->GetBridgeAttribute(&context, req, &resp);
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
      case SAI_BRIDGE_ATTR_TYPE:
        attr_list[i].value.s32 =
            convert_sai_bridge_type_t_to_sai(resp.attr().type());
        break;
      case SAI_BRIDGE_ATTR_PORT_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().port_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_BRIDGE_ATTR_MAX_LEARNED_ADDRESSES:
        attr_list[i].value.u32 = resp.attr().max_learned_addresses();
        break;
      case SAI_BRIDGE_ATTR_LEARN_DISABLE:
        attr_list[i].value.booldata = resp.attr().learn_disable();
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
        attr_list[i].value.s32 = convert_sai_bridge_flood_control_type_t_to_sai(
            resp.attr().unknown_unicast_flood_control_type());
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
        attr_list[i].value.oid = resp.attr().unknown_unicast_flood_group();
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
        attr_list[i].value.s32 = convert_sai_bridge_flood_control_type_t_to_sai(
            resp.attr().unknown_multicast_flood_control_type());
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
        attr_list[i].value.oid = resp.attr().unknown_multicast_flood_group();
        break;
      case SAI_BRIDGE_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
        attr_list[i].value.s32 = convert_sai_bridge_flood_control_type_t_to_sai(
            resp.attr().broadcast_flood_control_type());
        break;
      case SAI_BRIDGE_ATTR_BROADCAST_FLOOD_GROUP:
        attr_list[i].value.oid = resp.attr().broadcast_flood_group();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_bridge_stats(sai_object_id_t bridge_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids,
                                uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetBridgeStatsRequest req;
  lemming::dataplane::sai::GetBridgeStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(bridge_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(convert_sai_bridge_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status = bridge->GetBridgeStats(&context, req, &resp);
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

sai_status_t l_get_bridge_stats_ext(sai_object_id_t bridge_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t *counter_ids,
                                    sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_bridge_stats(sai_object_id_t bridge_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_bridge_port(sai_object_id_t *bridge_port_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateBridgePortRequest req =
      convert_create_bridge_port(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateBridgePortResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = bridge->CreateBridgePort(&context, req, &resp);
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
  if (bridge_port_id) {
    *bridge_port_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_bridge_port(sai_object_id_t bridge_port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveBridgePortRequest req;
  lemming::dataplane::sai::RemoveBridgePortResponse resp;
  grpc::ClientContext context;
  req.set_oid(bridge_port_id);

  grpc::Status status = bridge->RemoveBridgePort(&context, req, &resp);
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

sai_status_t l_set_bridge_port_attribute(sai_object_id_t bridge_port_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetBridgePortAttributeRequest req;
  lemming::dataplane::sai::SetBridgePortAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(bridge_port_id);

  switch (attr->id) {
    case SAI_BRIDGE_PORT_ATTR_TAGGING_MODE:
      req.set_tagging_mode(
          convert_sai_bridge_port_tagging_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_BRIDGE_PORT_ATTR_BRIDGE_ID:
      req.set_bridge_id(attr->value.oid);
      break;
    case SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_MODE:
      req.set_fdb_learning_mode(
          convert_sai_bridge_port_fdb_learning_mode_t_to_proto(
              attr->value.s32));
      break;
    case SAI_BRIDGE_PORT_ATTR_MAX_LEARNED_ADDRESSES:
      req.set_max_learned_addresses(attr->value.u32);
      break;
    case SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_LIMIT_VIOLATION_PACKET_ACTION:
      req.set_fdb_learning_limit_violation_packet_action(
          convert_sai_packet_action_t_to_proto(attr->value.s32));
      break;
    case SAI_BRIDGE_PORT_ATTR_ADMIN_STATE:
      req.set_admin_state(attr->value.booldata);
      break;
    case SAI_BRIDGE_PORT_ATTR_INGRESS_FILTERING:
      req.set_ingress_filtering(attr->value.booldata);
      break;
    case SAI_BRIDGE_PORT_ATTR_EGRESS_FILTERING:
      req.set_egress_filtering(attr->value.booldata);
      break;
    case SAI_BRIDGE_PORT_ATTR_ISOLATION_GROUP:
      req.set_isolation_group(attr->value.oid);
      break;
  }

  grpc::Status status = bridge->SetBridgePortAttribute(&context, req, &resp);
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

sai_status_t l_get_bridge_port_attribute(sai_object_id_t bridge_port_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetBridgePortAttributeRequest req;
  lemming::dataplane::sai::GetBridgePortAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(bridge_port_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_bridge_port_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = bridge->GetBridgePortAttribute(&context, req, &resp);
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
      case SAI_BRIDGE_PORT_ATTR_TYPE:
        attr_list[i].value.s32 =
            convert_sai_bridge_port_type_t_to_sai(resp.attr().type());
        break;
      case SAI_BRIDGE_PORT_ATTR_PORT_ID:
        attr_list[i].value.oid = resp.attr().port_id();
        break;
      case SAI_BRIDGE_PORT_ATTR_TAGGING_MODE:
        attr_list[i].value.s32 = convert_sai_bridge_port_tagging_mode_t_to_sai(
            resp.attr().tagging_mode());
        break;
      case SAI_BRIDGE_PORT_ATTR_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().vlan_id();
        break;
      case SAI_BRIDGE_PORT_ATTR_RIF_ID:
        attr_list[i].value.oid = resp.attr().rif_id();
        break;
      case SAI_BRIDGE_PORT_ATTR_TUNNEL_ID:
        attr_list[i].value.oid = resp.attr().tunnel_id();
        break;
      case SAI_BRIDGE_PORT_ATTR_BRIDGE_ID:
        attr_list[i].value.oid = resp.attr().bridge_id();
        break;
      case SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_MODE:
        attr_list[i].value.s32 =
            convert_sai_bridge_port_fdb_learning_mode_t_to_sai(
                resp.attr().fdb_learning_mode());
        break;
      case SAI_BRIDGE_PORT_ATTR_MAX_LEARNED_ADDRESSES:
        attr_list[i].value.u32 = resp.attr().max_learned_addresses();
        break;
      case SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_LIMIT_VIOLATION_PACKET_ACTION:
        attr_list[i].value.s32 = convert_sai_packet_action_t_to_sai(
            resp.attr().fdb_learning_limit_violation_packet_action());
        break;
      case SAI_BRIDGE_PORT_ATTR_ADMIN_STATE:
        attr_list[i].value.booldata = resp.attr().admin_state();
        break;
      case SAI_BRIDGE_PORT_ATTR_INGRESS_FILTERING:
        attr_list[i].value.booldata = resp.attr().ingress_filtering();
        break;
      case SAI_BRIDGE_PORT_ATTR_EGRESS_FILTERING:
        attr_list[i].value.booldata = resp.attr().egress_filtering();
        break;
      case SAI_BRIDGE_PORT_ATTR_ISOLATION_GROUP:
        attr_list[i].value.oid = resp.attr().isolation_group();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_bridge_port_stats(sai_object_id_t bridge_port_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetBridgePortStatsRequest req;
  lemming::dataplane::sai::GetBridgePortStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(bridge_port_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(
        convert_sai_bridge_port_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status = bridge->GetBridgePortStats(&context, req, &resp);
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

sai_status_t l_get_bridge_port_stats_ext(sai_object_id_t bridge_port_id,
                                         uint32_t number_of_counters,
                                         const sai_stat_id_t *counter_ids,
                                         sai_stats_mode_t mode,
                                         uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_bridge_port_stats(sai_object_id_t bridge_port_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
