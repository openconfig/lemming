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

#include "dataplane/standalone/switch.h"

#include <glog/logging.h>

#include <string>
#include <thread>  // NOLINT
#include <vector>

#include "dataplane/standalone/acl.h"
#include "dataplane/standalone/bridge.h"
#include "dataplane/standalone/buffer.h"
#include "dataplane/standalone/dtel.h"
#include "dataplane/standalone/hostif.h"
#include "dataplane/standalone/lucius/lucius_clib.h"
#include "dataplane/standalone/port.h"
#include "dataplane/standalone/route.h"
#include "dataplane/standalone/router_interface.h"
#include "dataplane/standalone/translator.h"
#include "dataplane/standalone/vlan.h"

extern "C" {
#include "inc/sai.h"
}

sai_status_t Switch::create(_In_ uint32_t attr_count,
                            _In_ const sai_attribute_t *attr_list) {
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      // TODO(dgrau): handle this attributes specially.
      case SAI_SWITCH_ATTR_INIT_SWITCH: {
        attr_list[i].value.booldata;
      }
      case SAI_SWITCH_ATTR_FDB_EVENT_NOTIFY: {
        reinterpret_cast<sai_fdb_event_notification_fn>(attr_list[i].value.ptr);
      }
      case SAI_SWITCH_ATTR_PORT_STATE_CHANGE_NOTIFY: {
        this->port_callback_fn =
            reinterpret_cast<sai_port_state_change_notification_fn>(
                attr_list[i].value.ptr);
      }
      case SAI_SWITCH_ATTR_SWITCH_SHUTDOWN_REQUEST_NOTIFY: {
        reinterpret_cast<sai_switch_shutdown_request_notification_fn>(
            attr_list[i].value.ptr);
      }
      case SAI_SWITCH_ATTR_SRC_MAC_ADDRESS: {
        attr_list[i].value.mac;
      }
    }
  }

  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_NUMBER_OF_ACTIVE_PORTS,
      .value = {.u32 = 0},
  });

  attrs.push_back({
      .id = SAI_SWITCH_ATTR_PORT_LIST,
      .value = {.objlist = {.count = 0}},
  });

  auto portOid = this->attrMgr->create(SAI_OBJECT_TYPE_PORT, this->id);
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_CPU_PORT,
      .value = {.oid = portOid},
  });

  attrs.push_back({
      .id = SAI_SWITCH_ATTR_ACL_ENTRY_MINIMUM_PRIORITY,
      .value = {.u32 = 1},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_ACL_ENTRY_MAXIMUM_PRIORITY,
      .value = {.u32 = 16000},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_MAX_ACL_ACTION_COUNT,
      .value = {.u32 = 53},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_ACL_STAGE_INGRESS,
      .value = {.aclcapability = {}},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_ACL_STAGE_EGRESS,
      .value = {.aclcapability = {}},
  });

  auto vlanOid = this->attrMgr->create(SAI_OBJECT_TYPE_VLAN, this->id);
  this->create_child(SAI_OBJECT_TYPE_VLAN, vlanOid, 0, nullptr);

  attrs.push_back({
      .id = SAI_SWITCH_ATTR_DEFAULT_VLAN_ID,
      .value = {.oid = vlanOid},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_NUMBER_OF_ECMP_GROUPS,
      .value = {.u32 = 1024},
  });

  auto stpOid = this->attrMgr->create(SAI_OBJECT_TYPE_STP, this->id);
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_DEFAULT_STP_INST_ID,
      .value = {.oid = stpOid},
  });

  auto vrOid = this->attrMgr->create(SAI_OBJECT_TYPE_VIRTUAL_ROUTER, this->id);
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_DEFAULT_VIRTUAL_ROUTER_ID,
      .value = {.oid = vrOid},
  });

  attrs.push_back({
      .id = SAI_SWITCH_ATTR_DEFAULT_OVERRIDE_VIRTUAL_ROUTER_ID,
      .value = {.oid = vrOid},
  });

  auto brOid = this->attrMgr->create(SAI_OBJECT_TYPE_BRIDGE, this->id);
  this->create_child(SAI_OBJECT_TYPE_BRIDGE, brOid, 0, nullptr);

  attrs.push_back({
      .id = SAI_SWITCH_ATTR_DEFAULT_1Q_BRIDGE_ID,
      .value = {.oid = brOid},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_INGRESS_ACL,
      .value = {.oid = SAI_NULL_OBJECT_ID},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_EGRESS_ACL,
      .value = {.oid = SAI_NULL_OBJECT_ID},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_TRAFFIC_CLASSES,
      .value = {.u8 = 0},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_TOTAL_BUFFER_SIZE,
      .value = {.u64 = 1024 * 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_IPV4_ROUTE_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_IPV6_ROUTE_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_IPV4_NEXTHOP_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_IPV6_NEXTHOP_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_IPV4_NEIGHBOR_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_IPV6_NEIGHBOR_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_MEMBER_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_FDB_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_L2MC_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_IPMC_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_SNAT_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_DNAT_ENTRY,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_ACL_TABLE,
      .value = {.u32 = 1024},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_AVAILABLE_ACL_TABLE_GROUP,
      .value = {.u32 = 1024},
  });

  auto trGrpOid =
      this->attrMgr->create(SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP, this->id);
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_DEFAULT_TRAP_GROUP,
      .value = {.oid = trGrpOid},
  });
  auto ecmpOid = this->attrMgr->create(SAI_OBJECT_TYPE_HASH, this->id);
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_ECMP_HASH,
      .value = {.oid = ecmpOid},
  });

  auto hashOid = this->attrMgr->create(SAI_OBJECT_TYPE_HASH, this->id);
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_LAG_HASH,
      .value = {.oid = hashOid},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_RESTART_WARM,
      .value = {.booldata = false},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_WARM_RECOVER,
      .value = {.booldata = false},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM,
      .value = {.s32 = SAI_HASH_ALGORITHM_CRC},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_SEED,
      .value = {.u32 = 0},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH,
      .value = {.booldata = false},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_QOS_DEFAULT_TC,
      .value = {.u8 = 0},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP,
      .value = {.oid = SAI_NULL_OBJECT_ID},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP,
      .value = {.oid = SAI_NULL_OBJECT_ID},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP,
      .value = {.oid = SAI_NULL_OBJECT_ID},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP,
      .value = {.oid = SAI_NULL_OBJECT_ID},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP,
      .value = {.oid = SAI_NULL_OBJECT_ID},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_SWITCH_SHELL_ENABLE,
      .value = {.booldata = false},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_SWITCH_PROFILE_ID,
      .value = {.u32 = 0},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_SWITCH_HARDWARE_INFO,
      .value = {.s8list = {.count = 0}},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_FIRMWARE_PATH_NAME,
      .value = {.s8list = {.count = 0}},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID,
      .value = {.oid = SAI_NULL_OBJECT_ID},
  });
  attrs.push_back({
      .id = SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP,
      .value = {.oid = SAI_NULL_OBJECT_ID},
  });

  LOG(INFO) << "Starting notif thread";
  std::thread thread(&Switch::handle_notification, this);
  thread.detach();

  APIBase::create(attrs.size(), attrs.data());
  LOG(INFO) << "Switch created successfuly";
  return SAI_STATUS_SUCCESS;
}

sai_status_t Switch::set_attribute(const sai_attribute_t *attr) {
  return sai_status_t();
}

sai_status_t Switch::create_child(sai_object_type_t type, sai_object_id_t id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  switch (type) {
    case SAI_OBJECT_TYPE_PORT:
      this->apis[std::to_string(id)] = std::make_unique<Port>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_ROUTER_INTERFACE:
      this->apis[std::to_string(id)] = std::make_unique<RouterInterface>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_VLAN:
      this->apis[std::to_string(id)] = std::make_unique<VLAN>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_BRIDGE:
      this->apis[std::to_string(id)] = std::make_unique<Bridge>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_HOSTIF:
      this->apis[std::to_string(id)] = std::make_unique<HostIf>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY:
      this->apis[std::to_string(id)] = std::make_unique<HostIfTableEntry>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_HOSTIF_TRAP:
      this->apis[std::to_string(id)] = std::make_unique<HostIfTrap>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_DTEL:
      this->apis[std::to_string(id)] = std::make_unique<DTEL>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_ACL_TABLE:
      this->apis[std::to_string(id)] = std::make_unique<ACLTable>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
    case SAI_OBJECT_TYPE_BUFFER_POOL:
      this->apis[std::to_string(id)] = std::make_unique<BufferPool>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_BUFFER_PROFILE:
      this->apis[std::to_string(id)] = std::make_unique<BufferProfile>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_NEXT_HOP:
      this->apis[std::to_string(id)] = std::make_unique<NextHop>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_NEXT_HOP_GROUP:
      this->apis[std::to_string(id)] = std::make_unique<NextHopGroup>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    case SAI_OBJECT_TYPE_NEXT_HOP_GROUP_MEMBER:
      this->apis[std::to_string(id)] = std::make_unique<NextHopGroupMember>(
          std::to_string(id), this->attrMgr, this->fwd, this->dataplane);
      break;
    default:
      return SAI_STATUS_FAILURE;
  }
  return this->apis[std::to_string(id)]->create(attr_count, attr_list);
}

sai_status_t Switch::create_child(sai_object_type_t type, common_entry_t entry,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  std::string id = this->attrMgr->serialize_entry(type, entry);
  switch (type) {
    case SAI_OBJECT_TYPE_ROUTE_ENTRY:
      this->apis[id] = std::make_unique<Route>(id, this->attrMgr, this->fwd,
                                               this->dataplane);
      break;
    default:
      return SAI_STATUS_FAILURE;
  }
  return this->apis[id]->create(entry, attr_count, attr_list);
}

sai_status_t Switch::set_child_attr(sai_object_type_t type, std::string id,
                                    const sai_attribute_t *attr) {
  return this->apis[id]->set_attribute(attr);
}

void Switch::handle_notification() {
  grpc::ClientContext ctx;
  forwarding::NotifySubscribeRequest req;
  char *id = getForwardCtxID();
  req.mutable_context()->set_id(id);
  auto reader = this->fwd->NotifySubscribe(&ctx, req);
  free(id);
  forwarding::EventDesc ed;
  while (reader->Read(&ed)) {
    if (!ed.has_port()) {
      continue;
    }
    auto type = attrMgr->get_type(ed.port().port_id().object_id().id());
    if (type !=
        SAI_OBJECT_TYPE_PORT) {  // Ignore notification for host if ports.
      continue;
    }
    sai_port_oper_status_notification_t sai_notif{
        .port_id = std::stoul(ed.port().port_id().object_id().id()),
    };
    switch (ed.port().port_info().oper_status()) {
      case forwarding::PORT_STATE_ENABLED_UP:
        sai_notif.port_state = SAI_PORT_OPER_STATUS_UP;
        break;
      case forwarding::PORT_STATE_DISABLED_DOWN:
        sai_notif.port_state = SAI_PORT_OPER_STATUS_DOWN;
        break;
      default:
        sai_notif.port_state = SAI_PORT_OPER_STATUS_UNKNOWN;
    }
    LOG(INFO) << "Sending port callback for port id " << sai_notif.port_id;
    this->port_callback_fn(1, &sai_notif);
  }
  grpc::Status st = reader->Finish();
  if (!st.ok()) {
    LOG(ERROR) << st.error_message();
  }
}
