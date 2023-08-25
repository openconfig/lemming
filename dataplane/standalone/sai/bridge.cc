

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

#include "dataplane/standalone/proto/bridge.pb.h"
#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

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

sai_status_t l_create_bridge(sai_object_id_t *bridge_id,
                             sai_object_id_t switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateBridgeRequest req;
  lemming::dataplane::sai::CreateBridgeResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BRIDGE_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::BridgeType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_BRIDGE_ATTR_MAX_LEARNED_ADDRESSES:
        req.set_max_learned_addresses(attr_list[i].value.u32);
        break;
      case SAI_BRIDGE_ATTR_LEARN_DISABLE:
        req.set_learn_disable(attr_list[i].value.booldata);
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
        req.set_unknown_unicast_flood_control_type(
            static_cast<lemming::dataplane::sai::BridgeFloodControlType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
        req.set_unknown_unicast_flood_group(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
        req.set_unknown_multicast_flood_control_type(
            static_cast<lemming::dataplane::sai::BridgeFloodControlType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
        req.set_unknown_multicast_flood_group(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
        req.set_broadcast_flood_control_type(
            static_cast<lemming::dataplane::sai::BridgeFloodControlType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BRIDGE_ATTR_BROADCAST_FLOOD_GROUP:
        req.set_broadcast_flood_group(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = bridge->CreateBridge(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *bridge_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_BRIDGE, bridge_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_bridge(sai_object_id_t bridge_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_BRIDGE, bridge_id);
}

sai_status_t l_set_bridge_attribute(sai_object_id_t bridge_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_BRIDGE, bridge_id, attr);
}

sai_status_t l_get_bridge_attribute(sai_object_id_t bridge_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_BRIDGE, bridge_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_bridge_stats(sai_object_id_t bridge_id,
                                uint32_t number_of_counters,
                                const sai_stat_id_t *counter_ids,
                                uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_BRIDGE, bridge_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_bridge_stats_ext(sai_object_id_t bridge_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t *counter_ids,
                                    sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_BRIDGE, bridge_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_bridge_stats(sai_object_id_t bridge_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_BRIDGE, bridge_id,
                                 number_of_counters, counter_ids);
}

sai_status_t l_create_bridge_port(sai_object_id_t *bridge_port_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateBridgePortRequest req;
  lemming::dataplane::sai::CreateBridgePortResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_BRIDGE_PORT_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::BridgePortType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_BRIDGE_PORT_ATTR_PORT_ID:
        req.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_PORT_ATTR_TAGGING_MODE:
        req.set_tagging_mode(
            static_cast<lemming::dataplane::sai::BridgePortTaggingMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BRIDGE_PORT_ATTR_VLAN_ID:
        req.set_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_BRIDGE_PORT_ATTR_RIF_ID:
        req.set_rif_id(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_PORT_ATTR_TUNNEL_ID:
        req.set_tunnel_id(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_PORT_ATTR_BRIDGE_ID:
        req.set_bridge_id(attr_list[i].value.oid);
        break;
      case SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_MODE:
        req.set_fdb_learning_mode(
            static_cast<lemming::dataplane::sai::BridgePortFdbLearningMode>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BRIDGE_PORT_ATTR_MAX_LEARNED_ADDRESSES:
        req.set_max_learned_addresses(attr_list[i].value.u32);
        break;
      case SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_LIMIT_VIOLATION_PACKET_ACTION:
        req.set_fdb_learning_limit_violation_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_BRIDGE_PORT_ATTR_ADMIN_STATE:
        req.set_admin_state(attr_list[i].value.booldata);
        break;
      case SAI_BRIDGE_PORT_ATTR_INGRESS_FILTERING:
        req.set_ingress_filtering(attr_list[i].value.booldata);
        break;
      case SAI_BRIDGE_PORT_ATTR_EGRESS_FILTERING:
        req.set_egress_filtering(attr_list[i].value.booldata);
        break;
      case SAI_BRIDGE_PORT_ATTR_ISOLATION_GROUP:
        req.set_isolation_group(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = bridge->CreateBridgePort(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *bridge_port_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_BRIDGE_PORT, bridge_port_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_bridge_port(sai_object_id_t bridge_port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_BRIDGE_PORT, bridge_port_id);
}

sai_status_t l_set_bridge_port_attribute(sai_object_id_t bridge_port_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_BRIDGE_PORT, bridge_port_id,
                                   attr);
}

sai_status_t l_get_bridge_port_attribute(sai_object_id_t bridge_port_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_BRIDGE_PORT, bridge_port_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_bridge_port_stats(sai_object_id_t bridge_port_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_BRIDGE_PORT, bridge_port_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_bridge_port_stats_ext(sai_object_id_t bridge_port_id,
                                         uint32_t number_of_counters,
                                         const sai_stat_id_t *counter_ids,
                                         sai_stats_mode_t mode,
                                         uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_BRIDGE_PORT, bridge_port_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_bridge_port_stats(sai_object_id_t bridge_port_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_BRIDGE_PORT, bridge_port_id,
                                 number_of_counters, counter_ids);
}
