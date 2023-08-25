

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

#include "dataplane/standalone/sai/ipsec.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/ipsec.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_ipsec_api_t l_ipsec = {
    .create_ipsec = l_create_ipsec,
    .remove_ipsec = l_remove_ipsec,
    .set_ipsec_attribute = l_set_ipsec_attribute,
    .get_ipsec_attribute = l_get_ipsec_attribute,
    .create_ipsec_port = l_create_ipsec_port,
    .remove_ipsec_port = l_remove_ipsec_port,
    .set_ipsec_port_attribute = l_set_ipsec_port_attribute,
    .get_ipsec_port_attribute = l_get_ipsec_port_attribute,
    .get_ipsec_port_stats = l_get_ipsec_port_stats,
    .get_ipsec_port_stats_ext = l_get_ipsec_port_stats_ext,
    .clear_ipsec_port_stats = l_clear_ipsec_port_stats,
    .create_ipsec_sa = l_create_ipsec_sa,
    .remove_ipsec_sa = l_remove_ipsec_sa,
    .set_ipsec_sa_attribute = l_set_ipsec_sa_attribute,
    .get_ipsec_sa_attribute = l_get_ipsec_sa_attribute,
    .get_ipsec_sa_stats = l_get_ipsec_sa_stats,
    .get_ipsec_sa_stats_ext = l_get_ipsec_sa_stats_ext,
    .clear_ipsec_sa_stats = l_clear_ipsec_sa_stats,
};

sai_status_t l_create_ipsec(sai_object_id_t *ipsec_id,
                            sai_object_id_t switch_id, uint32_t attr_count,
                            const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIpsecRequest req;
  lemming::dataplane::sai::CreateIpsecResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_IPSEC_ATTR_WARM_BOOT_ENABLE:
        req.set_warm_boot_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_ATTR_EXTERNAL_SA_INDEX_ENABLE:
        req.set_external_sa_index_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_ATTR_CTAG_TPID:
        req.set_ctag_tpid(attr_list[i].value.u16);
        break;
      case SAI_IPSEC_ATTR_STAG_TPID:
        req.set_stag_tpid(attr_list[i].value.u16);
        break;
      case SAI_IPSEC_ATTR_MAX_VLAN_TAGS_PARSED:
        req.set_max_vlan_tags_parsed(attr_list[i].value.u8);
        break;
      case SAI_IPSEC_ATTR_OCTET_COUNT_HIGH_WATERMARK:
        req.set_octet_count_high_watermark(attr_list[i].value.u64);
        break;
      case SAI_IPSEC_ATTR_OCTET_COUNT_LOW_WATERMARK:
        req.set_octet_count_low_watermark(attr_list[i].value.u64);
        break;
      case SAI_IPSEC_ATTR_STATS_MODE:
        req.set_stats_mode(static_cast<lemming::dataplane::sai::StatsMode>(
            attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = ipsec->CreateIpsec(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *ipsec_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_IPSEC, ipsec_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_ipsec(sai_object_id_t ipsec_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_IPSEC, ipsec_id);
}

sai_status_t l_set_ipsec_attribute(sai_object_id_t ipsec_id,
                                   const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_IPSEC, ipsec_id, attr);
}

sai_status_t l_get_ipsec_attribute(sai_object_id_t ipsec_id,
                                   uint32_t attr_count,
                                   sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_IPSEC, ipsec_id, attr_count,
                                   attr_list);
}

sai_status_t l_create_ipsec_port(sai_object_id_t *ipsec_port_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIpsecPortRequest req;
  lemming::dataplane::sai::CreateIpsecPortResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_IPSEC_PORT_ATTR_PORT_ID:
        req.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_IPSEC_PORT_ATTR_CTAG_ENABLE:
        req.set_ctag_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_PORT_ATTR_STAG_ENABLE:
        req.set_stag_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_PORT_ATTR_NATIVE_VLAN_ID:
        req.set_native_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_IPSEC_PORT_ATTR_VRF_FROM_PACKET_VLAN_ENABLE:
        req.set_vrf_from_packet_vlan_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
        req.set_switch_switching_mode(
            static_cast<lemming::dataplane::sai::SwitchSwitchingMode>(
                attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = ipsec->CreateIpsecPort(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *ipsec_port_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_IPSEC_PORT, ipsec_port_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_ipsec_port(sai_object_id_t ipsec_port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_IPSEC_PORT, ipsec_port_id);
}

sai_status_t l_set_ipsec_port_attribute(sai_object_id_t ipsec_port_id,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_IPSEC_PORT, ipsec_port_id,
                                   attr);
}

sai_status_t l_get_ipsec_port_attribute(sai_object_id_t ipsec_port_id,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_IPSEC_PORT, ipsec_port_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_ipsec_port_stats(sai_object_id_t ipsec_port_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t *counter_ids,
                                    uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_IPSEC_PORT, ipsec_port_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_ipsec_port_stats_ext(sai_object_id_t ipsec_port_id,
                                        uint32_t number_of_counters,
                                        const sai_stat_id_t *counter_ids,
                                        sai_stats_mode_t mode,
                                        uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_IPSEC_PORT, ipsec_port_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_ipsec_port_stats(sai_object_id_t ipsec_port_id,
                                      uint32_t number_of_counters,
                                      const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_IPSEC_PORT, ipsec_port_id,
                                 number_of_counters, counter_ids);
}

sai_status_t l_create_ipsec_sa(sai_object_id_t *ipsec_sa_id,
                               sai_object_id_t switch_id, uint32_t attr_count,
                               const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIpsecSaRequest req;
  lemming::dataplane::sai::CreateIpsecSaResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_IPSEC_SA_ATTR_IPSEC_DIRECTION:
        req.set_ipsec_direction(
            static_cast<lemming::dataplane::sai::IpsecDirection>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_ID:
        req.set_ipsec_id(attr_list[i].value.oid);
        break;
      case SAI_IPSEC_SA_ATTR_EXTERNAL_SA_INDEX:
        req.set_external_sa_index(attr_list[i].value.u32);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_PORT_LIST:
        req.mutable_ipsec_port_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_SPI:
        req.set_ipsec_spi(attr_list[i].value.u32);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_ESN_ENABLE:
        req.set_ipsec_esn_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_CIPHER:
        req.set_ipsec_cipher(static_cast<lemming::dataplane::sai::IpsecCipher>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_IPSEC_SA_ATTR_SALT:
        req.set_salt(attr_list[i].value.u32);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_ENABLE:
        req.set_ipsec_replay_protection_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_WINDOW:
        req.set_ipsec_replay_protection_window(attr_list[i].value.u32);
        break;
      case SAI_IPSEC_SA_ATTR_TERM_DST_IP:
        req.set_term_dst_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_IPSEC_SA_ATTR_TERM_VLAN_ID_ENABLE:
        req.set_term_vlan_id_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_SA_ATTR_TERM_VLAN_ID:
        req.set_term_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_IPSEC_SA_ATTR_TERM_SRC_IP_ENABLE:
        req.set_term_src_ip_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_SA_ATTR_TERM_SRC_IP:
        req.set_term_src_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_IPSEC_SA_ATTR_EGRESS_ESN:
        req.set_egress_esn(attr_list[i].value.u64);
        break;
      case SAI_IPSEC_SA_ATTR_MINIMUM_INGRESS_ESN:
        req.set_minimum_ingress_esn(attr_list[i].value.u64);
        break;
    }
  }
  grpc::Status status = ipsec->CreateIpsecSa(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *ipsec_sa_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_IPSEC_SA, ipsec_sa_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_ipsec_sa(sai_object_id_t ipsec_sa_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_IPSEC_SA, ipsec_sa_id);
}

sai_status_t l_set_ipsec_sa_attribute(sai_object_id_t ipsec_sa_id,
                                      const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_IPSEC_SA, ipsec_sa_id, attr);
}

sai_status_t l_get_ipsec_sa_attribute(sai_object_id_t ipsec_sa_id,
                                      uint32_t attr_count,
                                      sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_IPSEC_SA, ipsec_sa_id,
                                   attr_count, attr_list);
}

sai_status_t l_get_ipsec_sa_stats(sai_object_id_t ipsec_sa_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t *counter_ids,
                                  uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats(SAI_OBJECT_TYPE_IPSEC_SA, ipsec_sa_id,
                               number_of_counters, counter_ids, counters);
}

sai_status_t l_get_ipsec_sa_stats_ext(sai_object_id_t ipsec_sa_id,
                                      uint32_t number_of_counters,
                                      const sai_stat_id_t *counter_ids,
                                      sai_stats_mode_t mode,
                                      uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_stats_ext(SAI_OBJECT_TYPE_IPSEC_SA, ipsec_sa_id,
                                   number_of_counters, counter_ids, mode,
                                   counters);
}

sai_status_t l_clear_ipsec_sa_stats(sai_object_id_t ipsec_sa_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->clear_stats(SAI_OBJECT_TYPE_IPSEC_SA, ipsec_sa_id,
                                 number_of_counters, counter_ids);
}
