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

#include "dataplane/standalone/translator.h"

extern "C" {
#include "inc/sai.h"
#include "meta/saimetadata.h"
#include "switch.h"
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

  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_NUMBER_OF_ACTIVE_PORTS,
                                   .value = {.u32 = 0},
                               });

  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_PORT_LIST,
                                   .value = {.objlist = {.count = 0}},
                               });

  auto portOid = this->attrMgr->create(SAI_OBJECT_TYPE_PORT, this->id);

  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_CPU_PORT,
                                   .value = {.oid = portOid},
                               });

  auto vlanOid = this->attrMgr->create(SAI_OBJECT_TYPE_VLAN, this->id);
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_DEFAULT_VLAN_ID,
                                   .value = {.oid = vlanOid},
                               });

  auto stpOid = this->attrMgr->create(SAI_OBJECT_TYPE_STP, this->id);
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_DEFAULT_STP_INST_ID,
                                   .value = {.oid = stpOid},
                               });

  auto vrOid = this->attrMgr->create(SAI_OBJECT_TYPE_VIRTUAL_ROUTER, this->id);
  this->attrMgr->set_attribute(
      std::to_string(this->id),
      sai_attribute_t{
          .id = SAI_SWITCH_ATTR_DEFAULT_VIRTUAL_ROUTER_ID,
          .value = {.oid = vrOid},
      });

  this->attrMgr->set_attribute(
      std::to_string(this->id),
      sai_attribute_t{
          .id = SAI_SWITCH_ATTR_DEFAULT_OVERRIDE_VIRTUAL_ROUTER_ID,
          .value = {.oid = vrOid},
      });

  auto brOid = this->attrMgr->create(SAI_OBJECT_TYPE_BRIDGE, this->id);
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_DEFAULT_1Q_BRIDGE_ID,
                                   .value = {.oid = brOid},
                               });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_INGRESS_ACL,
                                   .value = {.oid = SAI_NULL_OBJECT_ID},
                               });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_EGRESS_ACL,
                                   .value = {.oid = SAI_NULL_OBJECT_ID},
                               });
  this->attrMgr->set_attribute(
      std::to_string(this->id),
      sai_attribute_t{
          .id = SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_TRAFFIC_CLASSES,
          .value = {.u8 = 0},
      });

  auto trGrpOid =
      this->attrMgr->create(SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP, this->id);
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_DEFAULT_TRAP_GROUP,
                                   .value = {.oid = trGrpOid},
                               });
  auto ecmpOid = this->attrMgr->create(SAI_OBJECT_TYPE_HASH, this->id);
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_ECMP_HASH,
                                   .value = {.oid = ecmpOid},
                               });

  auto hashOid = this->attrMgr->create(SAI_OBJECT_TYPE_HASH, this->id);
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_LAG_HASH,
                                   .value = {.oid = hashOid},
                               });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_RESTART_WARM,
                                   .value = {.booldata = false},
                               });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_WARM_RECOVER,
                                   .value = {.booldata = false},
                               });
  this->attrMgr->set_attribute(
      std::to_string(this->id),
      sai_attribute_t{
          .id = SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM,
          .value = {.s32 = SAI_HASH_ALGORITHM_CRC},
      });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_SEED,
                                   .value = {.u32 = 0},
                               });
  this->attrMgr->set_attribute(
      std::to_string(this->id),
      sai_attribute_t{
          .id = SAI_SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH,
          .value = {.booldata = false},
      });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_QOS_DEFAULT_TC,
                                   .value = {.u8 = 0},
                               });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP,
                                   .value = {.oid = SAI_NULL_OBJECT_ID},
                               });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP,
                                   .value = {.oid = SAI_NULL_OBJECT_ID},
                               });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP,
                                   .value = {.oid = SAI_NULL_OBJECT_ID},
                               });
  this->attrMgr->set_attribute(
      std::to_string(this->id),
      sai_attribute_t{
          .id = SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP,
          .value = {.oid = SAI_NULL_OBJECT_ID},
      });
  this->attrMgr->set_attribute(
      std::to_string(this->id),
      sai_attribute_t{
          .id = SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP,
          .value = {.oid = SAI_NULL_OBJECT_ID},
      });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_SWITCH_SHELL_ENABLE,
                                   .value = {.booldata = false},
                               });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_SWITCH_PROFILE_ID,
                                   .value = {.u32 = 0},
                               });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_SWITCH_HARDWARE_INFO,
                                   .value = {.s8list = {.count = 0}},
                               });
  this->attrMgr->set_attribute(std::to_string(this->id),
                               sai_attribute_t{
                                   .id = SAI_SWITCH_ATTR_FIRMWARE_PATH_NAME,
                                   .value = {.s8list = {.count = 0}},
                               });

  this->attrMgr->set_attribute(
      std::to_string(this->id),
      sai_attribute_t{
          .id = SAI_SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID,
          .value = {.oid = SAI_NULL_OBJECT_ID},
      });

  this->attrMgr->set_attribute(
      std::to_string(this->id),
      sai_attribute_t{
          .id = SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP,
          .value = {.oid = SAI_NULL_OBJECT_ID},
      });

  return SAI_STATUS_SUCCESS;
}

sai_status_t Switch::create_child(sai_object_type_t type, sai_object_id_t id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  switch (type) {
    case SAI_OBJECT_TYPE_PORT:
      this->apis[std::to_string(id)] =
          std::make_unique<Port>(this->attrMgr, this->client);
    default:
      return SAI_STATUS_FAILURE;
  }
  return SAI_STATUS_SUCCESS;
}

sai_status_t Switch::create_child(sai_object_type_t type, common_entry_t id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  return sai_status_t();
}

sai_status_t Switch::set_child_attr(sai_object_type_t type, std::string id,
                                    const sai_attribute_t *attr) {
  return this->apis[id]->set_attribute(attr);
}
