

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

#include "dataplane/standalone/sai/nat.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/nat.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_nat_api_t l_nat = {
    .create_nat_entry = l_create_nat_entry,
    .remove_nat_entry = l_remove_nat_entry,
    .set_nat_entry_attribute = l_set_nat_entry_attribute,
    .get_nat_entry_attribute = l_get_nat_entry_attribute,
    .create_nat_entries = l_create_nat_entries,
    .remove_nat_entries = l_remove_nat_entries,
    .set_nat_entries_attribute = l_set_nat_entries_attribute,
    .get_nat_entries_attribute = l_get_nat_entries_attribute,
    .create_nat_zone_counter = l_create_nat_zone_counter,
    .remove_nat_zone_counter = l_remove_nat_zone_counter,
    .set_nat_zone_counter_attribute = l_set_nat_zone_counter_attribute,
    .get_nat_zone_counter_attribute = l_get_nat_zone_counter_attribute,
};

sai_status_t l_create_nat_entry(const sai_nat_entry_t *nat_entry,
                                uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateNatEntryRequest req;
  lemming::dataplane::sai::CreateNatEntryResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_NAT_ENTRY_ATTR_NAT_TYPE:
        req.set_nat_type(static_cast<lemming::dataplane::sai::NatType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_NAT_ENTRY_ATTR_SRC_IP:
        req.set_src_ip(&attr_list[i].value.ip4, sizeof(attr_list[i].value.ip4));
        break;
      case SAI_NAT_ENTRY_ATTR_SRC_IP_MASK:
        req.set_src_ip_mask(&attr_list[i].value.ip4,
                            sizeof(attr_list[i].value.ip4));
        break;
      case SAI_NAT_ENTRY_ATTR_VR_ID:
        req.set_vr_id(attr_list[i].value.oid);
        break;
      case SAI_NAT_ENTRY_ATTR_DST_IP:
        req.set_dst_ip(&attr_list[i].value.ip4, sizeof(attr_list[i].value.ip4));
        break;
      case SAI_NAT_ENTRY_ATTR_DST_IP_MASK:
        req.set_dst_ip_mask(&attr_list[i].value.ip4,
                            sizeof(attr_list[i].value.ip4));
        break;
      case SAI_NAT_ENTRY_ATTR_L4_SRC_PORT:
        req.set_l4_src_port(attr_list[i].value.u16);
        break;
      case SAI_NAT_ENTRY_ATTR_L4_DST_PORT:
        req.set_l4_dst_port(attr_list[i].value.u16);
        break;
      case SAI_NAT_ENTRY_ATTR_ENABLE_PACKET_COUNT:
        req.set_enable_packet_count(attr_list[i].value.booldata);
        break;
      case SAI_NAT_ENTRY_ATTR_PACKET_COUNT:
        req.set_packet_count(attr_list[i].value.u64);
        break;
      case SAI_NAT_ENTRY_ATTR_ENABLE_BYTE_COUNT:
        req.set_enable_byte_count(attr_list[i].value.booldata);
        break;
      case SAI_NAT_ENTRY_ATTR_BYTE_COUNT:
        req.set_byte_count(attr_list[i].value.u64);
        break;
      case SAI_NAT_ENTRY_ATTR_HIT_BIT_COR:
        req.set_hit_bit_cor(attr_list[i].value.booldata);
        break;
      case SAI_NAT_ENTRY_ATTR_HIT_BIT:
        req.set_hit_bit(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status = nat->CreateNatEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  common_entry_t entry = {.nat_entry = nat_entry};
  return translator->create(SAI_OBJECT_TYPE_NAT_ENTRY, entry, attr_count,
                            attr_list);
}

sai_status_t l_remove_nat_entry(const sai_nat_entry_t *nat_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.nat_entry = nat_entry};
  return translator->remove(SAI_OBJECT_TYPE_NAT_ENTRY, entry);
}

sai_status_t l_set_nat_entry_attribute(const sai_nat_entry_t *nat_entry,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.nat_entry = nat_entry};
  return translator->set_attribute(SAI_OBJECT_TYPE_NAT_ENTRY, entry, attr);
}

sai_status_t l_get_nat_entry_attribute(const sai_nat_entry_t *nat_entry,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.nat_entry = nat_entry};
  return translator->get_attribute(SAI_OBJECT_TYPE_NAT_ENTRY, entry, attr_count,
                                   attr_list);
}

sai_status_t l_create_nat_entries(uint32_t object_count,
                                  const sai_nat_entry_t *nat_entry,
                                  const uint32_t *attr_count,
                                  const sai_attribute_t **attr_list,
                                  sai_bulk_op_error_mode_t mode,
                                  sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.nat_entry = nat_entry};
  return translator->create_bulk(SAI_OBJECT_TYPE_NAT_ENTRY, object_count, entry,
                                 attr_count, attr_list, mode, object_statuses);
}

sai_status_t l_remove_nat_entries(uint32_t object_count,
                                  const sai_nat_entry_t *nat_entry,
                                  sai_bulk_op_error_mode_t mode,
                                  sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.nat_entry = nat_entry};
  return translator->remove_bulk(SAI_OBJECT_TYPE_NAT_ENTRY, object_count, entry,
                                 mode, object_statuses);
}

sai_status_t l_set_nat_entries_attribute(uint32_t object_count,
                                         const sai_nat_entry_t *nat_entry,
                                         const sai_attribute_t *attr_list,
                                         sai_bulk_op_error_mode_t mode,
                                         sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.nat_entry = nat_entry};
  return translator->set_attribute_bulk(SAI_OBJECT_TYPE_NAT_ENTRY, object_count,
                                        entry, attr_list, mode,
                                        object_statuses);
}

sai_status_t l_get_nat_entries_attribute(uint32_t object_count,
                                         const sai_nat_entry_t *nat_entry,
                                         const uint32_t *attr_count,
                                         sai_attribute_t **attr_list,
                                         sai_bulk_op_error_mode_t mode,
                                         sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  common_entry_t entry = {.nat_entry = nat_entry};
  return translator->get_attribute_bulk(SAI_OBJECT_TYPE_NAT_ENTRY, object_count,
                                        entry, attr_count, attr_list, mode,
                                        object_statuses);
}

sai_status_t l_create_nat_zone_counter(sai_object_id_t *nat_zone_counter_id,
                                       sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateNatZoneCounterRequest req;
  lemming::dataplane::sai::CreateNatZoneCounterResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_NAT_ZONE_COUNTER_ATTR_NAT_TYPE:
        req.set_nat_type(static_cast<lemming::dataplane::sai::NatType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_NAT_ZONE_COUNTER_ATTR_ZONE_ID:
        req.set_zone_id(attr_list[i].value.u8);
        break;
      case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_DISCARD:
        req.set_enable_discard(attr_list[i].value.booldata);
        break;
      case SAI_NAT_ZONE_COUNTER_ATTR_DISCARD_PACKET_COUNT:
        req.set_discard_packet_count(attr_list[i].value.u64);
        break;
      case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATION_NEEDED:
        req.set_enable_translation_needed(attr_list[i].value.booldata);
        break;
      case SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATION_NEEDED_PACKET_COUNT:
        req.set_translation_needed_packet_count(attr_list[i].value.u64);
        break;
      case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATIONS:
        req.set_enable_translations(attr_list[i].value.booldata);
        break;
      case SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATIONS_PACKET_COUNT:
        req.set_translations_packet_count(attr_list[i].value.u64);
        break;
    }
  }
  grpc::Status status = nat->CreateNatZoneCounter(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *nat_zone_counter_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_NAT_ZONE_COUNTER,
                            nat_zone_counter_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_nat_zone_counter(sai_object_id_t nat_zone_counter_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_NAT_ZONE_COUNTER,
                            nat_zone_counter_id);
}

sai_status_t l_set_nat_zone_counter_attribute(
    sai_object_id_t nat_zone_counter_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_NAT_ZONE_COUNTER,
                                   nat_zone_counter_id, attr);
}

sai_status_t l_get_nat_zone_counter_attribute(
    sai_object_id_t nat_zone_counter_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_NAT_ZONE_COUNTER,
                                   nat_zone_counter_id, attr_count, attr_list);
}
