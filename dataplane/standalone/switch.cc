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

#include "dataplane/standalone/log/log.h"
#include "dataplane/standalone/translator.h"

sai_status_t Switch::create_switch(_Out_ sai_object_id_t *switch_id,
                                   _In_ uint32_t attr_count,
                                   _In_ const sai_attribute_t *attr_list) {
  *switch_id = this->translator->createObject(SAI_OBJECT_TYPE_SWITCH);
  for (uint32_t i = 0; i < attr_count; i++) {
    LOG(attr_list[i].id);
    this->translator->setAttribute(*switch_id, sai_attribute_t{
                                                   .id = attr_list[i].id,
                                                   .value = attr_list[i].value,
                                               });
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
  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_NUMBER_OF_ACTIVE_PORTS,
                      .value = {.u32 = 0},
                  });  // 0

  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_PORT_LIST,
                                     .value = {.objlist = {.count = 0}},
                                 });  // 2

  auto portOid = this->translator->createObject(SAI_OBJECT_TYPE_PORT);

  this->translator->setAttribute(*switch_id, sai_attribute_t{
                                                 .id = SAI_SWITCH_ATTR_CPU_PORT,
                                                 .value = {.oid = portOid},
                                             });  // 4

  auto vlanOid = this->translator->createObject(SAI_OBJECT_TYPE_VLAN);
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_DEFAULT_VLAN_ID,
                                     .value = {.oid = vlanOid},
                                 });  // 36

  auto stpOid = this->translator->createObject(SAI_OBJECT_TYPE_STP);
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_DEFAULT_STP_INST_ID,
                                     .value = {.oid = stpOid},
                                 });  // 37

  auto vrOid = this->translator->createObject(SAI_OBJECT_TYPE_VIRTUAL_ROUTER);
  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_DEFAULT_VIRTUAL_ROUTER_ID,
                      .value = {.oid = vrOid},
                  });  // 39

  auto brOid = this->translator->createObject(SAI_OBJECT_TYPE_BRIDGE);
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_DEFAULT_1Q_BRIDGE_ID,
                                     .value = {.oid = brOid},
                                 });  // 34
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_INGRESS_ACL,
                                     .value = {.oid = SAI_NULL_OBJECT_ID},
                                 });  // 41
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_EGRESS_ACL,
                                     .value = {.oid = SAI_NULL_OBJECT_ID},
                                 });  // 42
  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_TRAFFIC_CLASSES,
                      .value = {.u8 = 0},
                  });  // 43

  auto hashOid = this->translator->createObject(SAI_OBJECT_TYPE_HASH);
  this->translator->setAttribute(*switch_id, sai_attribute_t{
                                                 .id = SAI_SWITCH_ATTR_LAG_HASH,
                                                 .value = {.oid = hashOid},
                                             });  // 68
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_RESTART_WARM,
                                     .value = {.booldata = false},
                                 });  // 69
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_WARM_RECOVER,
                                     .value = {.booldata = false},
                                 });  // 70
  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM,
                      .value = {.s32 = SAI_HASH_ALGORITHM_CRC},
                  });  // 93
  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_SEED,
                      .value = {.u32 = 0},
                  });  // 94
  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH,
                      .value = {.booldata = false},
                  });  // 95
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_QOS_DEFAULT_TC,
                                     .value = {.u8 = 0},
                                 });  // 100
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP,
                                     .value = {.oid = SAI_NULL_OBJECT_ID},
                                 });  // 101
  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP,
                      .value = {.oid = SAI_NULL_OBJECT_ID},
                  });  // 102
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP,
                                     .value = {.oid = SAI_NULL_OBJECT_ID},
                                 });  // 105
  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP,
                      .value = {.oid = SAI_NULL_OBJECT_ID},
                  });  // 106
  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP,
                      .value = {.oid = SAI_NULL_OBJECT_ID},
                  });  // 107
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_SWITCH_SHELL_ENABLE,
                                     .value = {.booldata = false},
                                 });  // 108
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_SWITCH_PROFILE_ID,
                                     .value = {.u32 = 0},
                                 });  // 109
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_SWITCH_HARDWARE_INFO,
                                     .value = {.s8list = {.count = 0}},
                                 });  // 110
  this->translator->setAttribute(*switch_id,
                                 sai_attribute_t{
                                     .id = SAI_SWITCH_ATTR_FIRMWARE_PATH_NAME,
                                     .value = {.s8list = {.count = 0}},
                                 });  // 111
  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID,
                      .value = {.oid = SAI_NULL_OBJECT_ID},
                  });  // 154

  this->translator->setAttribute(
      *switch_id, sai_attribute_t{
                      .id = SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP,
                      .value = {.oid = SAI_NULL_OBJECT_ID},
                  });  // 175

  LOG("created switch");
  return SAI_STATUS_SUCCESS;
}

sai_status_t Switch::set_switch_attribute(_In_ sai_object_id_t switch_id,
                                          _In_ const sai_attribute_t *attr) {
  this->translator->setAttribute(switch_id, *attr);
  return SAI_STATUS_SUCCESS;
}

sai_status_t Switch::get_switch_attribute(_In_ sai_object_id_t switch_id,
                                          _In_ uint32_t attr_count,
                                          _Inout_ sai_attribute_t *attr_list) {
  for (uint32_t i = 0; i < attr_count; i++) {
    LOG(attr_list[i].id);
    if (auto ret = this->translator->getAttribute(switch_id, &attr_list[i]);
        ret != SAI_STATUS_SUCCESS) {
      return ret;
    }
  }
  return SAI_STATUS_SUCCESS;
}
