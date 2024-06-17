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

#include "dataplane/standalone/sai/macsec.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/macsec.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_macsec_api_t l_macsec = {
    .create_macsec = l_create_macsec,
    .remove_macsec = l_remove_macsec,
    .set_macsec_attribute = l_set_macsec_attribute,
    .get_macsec_attribute = l_get_macsec_attribute,
    .create_macsec_port = l_create_macsec_port,
    .remove_macsec_port = l_remove_macsec_port,
    .set_macsec_port_attribute = l_set_macsec_port_attribute,
    .get_macsec_port_attribute = l_get_macsec_port_attribute,
    .get_macsec_port_stats = l_get_macsec_port_stats,
    .get_macsec_port_stats_ext = l_get_macsec_port_stats_ext,
    .clear_macsec_port_stats = l_clear_macsec_port_stats,
    .create_macsec_flow = l_create_macsec_flow,
    .remove_macsec_flow = l_remove_macsec_flow,
    .set_macsec_flow_attribute = l_set_macsec_flow_attribute,
    .get_macsec_flow_attribute = l_get_macsec_flow_attribute,
    .get_macsec_flow_stats = l_get_macsec_flow_stats,
    .get_macsec_flow_stats_ext = l_get_macsec_flow_stats_ext,
    .clear_macsec_flow_stats = l_clear_macsec_flow_stats,
    .create_macsec_sc = l_create_macsec_sc,
    .remove_macsec_sc = l_remove_macsec_sc,
    .set_macsec_sc_attribute = l_set_macsec_sc_attribute,
    .get_macsec_sc_attribute = l_get_macsec_sc_attribute,
    .get_macsec_sc_stats = l_get_macsec_sc_stats,
    .get_macsec_sc_stats_ext = l_get_macsec_sc_stats_ext,
    .clear_macsec_sc_stats = l_clear_macsec_sc_stats,
    .create_macsec_sa = l_create_macsec_sa,
    .remove_macsec_sa = l_remove_macsec_sa,
    .set_macsec_sa_attribute = l_set_macsec_sa_attribute,
    .get_macsec_sa_attribute = l_get_macsec_sa_attribute,
    .get_macsec_sa_stats = l_get_macsec_sa_stats,
    .get_macsec_sa_stats_ext = l_get_macsec_sa_stats_ext,
    .clear_macsec_sa_stats = l_clear_macsec_sa_stats,
};

lemming::dataplane::sai::CreateMacsecRequest convert_create_macsec(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateMacsecRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MACSEC_ATTR_DIRECTION:
        msg.set_direction(static_cast<lemming::dataplane::sai::MacsecDirection>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_MACSEC_ATTR_WARM_BOOT_ENABLE:
        msg.set_warm_boot_enable(attr_list[i].value.booldata);
        break;
      case SAI_MACSEC_ATTR_CTAG_TPID:
        msg.set_ctag_tpid(attr_list[i].value.u16);
        break;
      case SAI_MACSEC_ATTR_STAG_TPID:
        msg.set_stag_tpid(attr_list[i].value.u16);
        break;
      case SAI_MACSEC_ATTR_MAX_VLAN_TAGS_PARSED:
        msg.set_max_vlan_tags_parsed(attr_list[i].value.u8);
        break;
      case SAI_MACSEC_ATTR_STATS_MODE:
        msg.set_stats_mode(static_cast<lemming::dataplane::sai::StatsMode>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_MACSEC_ATTR_PHYSICAL_BYPASS_ENABLE:
        msg.set_physical_bypass_enable(attr_list[i].value.booldata);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateMacsecPortRequest convert_create_macsec_port(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateMacsecPortRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MACSEC_PORT_ATTR_MACSEC_DIRECTION:
        msg.set_macsec_direction(
            static_cast<lemming::dataplane::sai::MacsecDirection>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_MACSEC_PORT_ATTR_PORT_ID:
        msg.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_MACSEC_PORT_ATTR_CTAG_ENABLE:
        msg.set_ctag_enable(attr_list[i].value.booldata);
        break;
      case SAI_MACSEC_PORT_ATTR_STAG_ENABLE:
        msg.set_stag_enable(attr_list[i].value.booldata);
        break;
      case SAI_MACSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
        msg.set_switch_switching_mode(
            static_cast<lemming::dataplane::sai::SwitchSwitchingMode>(
                attr_list[i].value.s32 + 1));
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateMacsecFlowRequest convert_create_macsec_flow(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateMacsecFlowRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MACSEC_FLOW_ATTR_MACSEC_DIRECTION:
        msg.set_macsec_direction(
            static_cast<lemming::dataplane::sai::MacsecDirection>(
                attr_list[i].value.s32 + 1));
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateMacsecScRequest convert_create_macsec_sc(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateMacsecScRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MACSEC_SC_ATTR_MACSEC_DIRECTION:
        msg.set_macsec_direction(
            static_cast<lemming::dataplane::sai::MacsecDirection>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_MACSEC_SC_ATTR_FLOW_ID:
        msg.set_flow_id(attr_list[i].value.oid);
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_SCI:
        msg.set_macsec_sci(attr_list[i].value.u64);
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_EXPLICIT_SCI_ENABLE:
        msg.set_macsec_explicit_sci_enable(attr_list[i].value.booldata);
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_SECTAG_OFFSET:
        msg.set_macsec_sectag_offset(attr_list[i].value.u8);
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_ENABLE:
        msg.set_macsec_replay_protection_enable(attr_list[i].value.booldata);
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_WINDOW:
        msg.set_macsec_replay_protection_window(attr_list[i].value.u32);
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_CIPHER_SUITE:
        msg.set_macsec_cipher_suite(
            static_cast<lemming::dataplane::sai::MacsecCipherSuite>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_MACSEC_SC_ATTR_ENCRYPTION_ENABLE:
        msg.set_encryption_enable(attr_list[i].value.booldata);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateMacsecSaRequest convert_create_macsec_sa(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateMacsecSaRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MACSEC_SA_ATTR_MACSEC_DIRECTION:
        msg.set_macsec_direction(
            static_cast<lemming::dataplane::sai::MacsecDirection>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_MACSEC_SA_ATTR_SC_ID:
        msg.set_sc_id(attr_list[i].value.oid);
        break;
      case SAI_MACSEC_SA_ATTR_AN:
        msg.set_an(attr_list[i].value.u8);
        break;
      case SAI_MACSEC_SA_ATTR_CONFIGURED_EGRESS_XPN:
        msg.set_configured_egress_xpn(attr_list[i].value.u64);
        break;
      case SAI_MACSEC_SA_ATTR_MINIMUM_INGRESS_XPN:
        msg.set_minimum_ingress_xpn(attr_list[i].value.u64);
        break;
      case SAI_MACSEC_SA_ATTR_MACSEC_SSCI:
        msg.set_macsec_ssci(attr_list[i].value.u32);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_macsec(sai_object_id_t *macsec_id,
                             sai_object_id_t switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateMacsecRequest req =
      convert_create_macsec(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateMacsecResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = macsec->CreateMacsec(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *macsec_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_macsec(sai_object_id_t macsec_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveMacsecRequest req;
  lemming::dataplane::sai::RemoveMacsecResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_id);

  grpc::Status status = macsec->RemoveMacsec(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_macsec_attribute(sai_object_id_t macsec_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetMacsecAttributeRequest req;
  lemming::dataplane::sai::SetMacsecAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_id);

  switch (attr->id) {
    case SAI_MACSEC_ATTR_WARM_BOOT_ENABLE:
      req.set_warm_boot_enable(attr->value.booldata);
      break;
    case SAI_MACSEC_ATTR_CTAG_TPID:
      req.set_ctag_tpid(attr->value.u16);
      break;
    case SAI_MACSEC_ATTR_STAG_TPID:
      req.set_stag_tpid(attr->value.u16);
      break;
    case SAI_MACSEC_ATTR_MAX_VLAN_TAGS_PARSED:
      req.set_max_vlan_tags_parsed(attr->value.u8);
      break;
    case SAI_MACSEC_ATTR_STATS_MODE:
      req.set_stats_mode(
          static_cast<lemming::dataplane::sai::StatsMode>(attr->value.s32 + 1));
      break;
    case SAI_MACSEC_ATTR_PHYSICAL_BYPASS_ENABLE:
      req.set_physical_bypass_enable(attr->value.booldata);
      break;
  }

  grpc::Status status = macsec->SetMacsecAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_attribute(sai_object_id_t macsec_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMacsecAttributeRequest req;
  lemming::dataplane::sai::GetMacsecAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(macsec_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::MacsecAttr>(attr_list[i].id + 1));
  }
  grpc::Status status = macsec->GetMacsecAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MACSEC_ATTR_DIRECTION:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().direction() - 1);
        break;
      case SAI_MACSEC_ATTR_SWITCHING_MODE_CUT_THROUGH_SUPPORTED:
        attr_list[i].value.booldata =
            resp.attr().switching_mode_cut_through_supported();
        break;
      case SAI_MACSEC_ATTR_SWITCHING_MODE_STORE_AND_FORWARD_SUPPORTED:
        attr_list[i].value.booldata =
            resp.attr().switching_mode_store_and_forward_supported();
        break;
      case SAI_MACSEC_ATTR_STATS_MODE_READ_SUPPORTED:
        attr_list[i].value.booldata = resp.attr().stats_mode_read_supported();
        break;
      case SAI_MACSEC_ATTR_STATS_MODE_READ_CLEAR_SUPPORTED:
        attr_list[i].value.booldata =
            resp.attr().stats_mode_read_clear_supported();
        break;
      case SAI_MACSEC_ATTR_SCI_IN_INGRESS_MACSEC_ACL:
        attr_list[i].value.booldata = resp.attr().sci_in_ingress_macsec_acl();
        break;
      case SAI_MACSEC_ATTR_PN_32BIT_SUPPORTED:
        attr_list[i].value.booldata = resp.attr().pn_32bit_supported();
        break;
      case SAI_MACSEC_ATTR_XPN_64BIT_SUPPORTED:
        attr_list[i].value.booldata = resp.attr().xpn_64bit_supported();
        break;
      case SAI_MACSEC_ATTR_GCM_AES128_SUPPORTED:
        attr_list[i].value.booldata = resp.attr().gcm_aes128_supported();
        break;
      case SAI_MACSEC_ATTR_GCM_AES256_SUPPORTED:
        attr_list[i].value.booldata = resp.attr().gcm_aes256_supported();
        break;
      case SAI_MACSEC_ATTR_SECTAG_OFFSETS_SUPPORTED:
        copy_list(attr_list[i].value.u8list.list,
                  resp.attr().sectag_offsets_supported(),
                  &attr_list[i].value.u8list.count);
        break;
      case SAI_MACSEC_ATTR_SYSTEM_SIDE_MTU:
        attr_list[i].value.u16 = resp.attr().system_side_mtu();
        break;
      case SAI_MACSEC_ATTR_WARM_BOOT_SUPPORTED:
        attr_list[i].value.booldata = resp.attr().warm_boot_supported();
        break;
      case SAI_MACSEC_ATTR_WARM_BOOT_ENABLE:
        attr_list[i].value.booldata = resp.attr().warm_boot_enable();
        break;
      case SAI_MACSEC_ATTR_CTAG_TPID:
        attr_list[i].value.u16 = resp.attr().ctag_tpid();
        break;
      case SAI_MACSEC_ATTR_STAG_TPID:
        attr_list[i].value.u16 = resp.attr().stag_tpid();
        break;
      case SAI_MACSEC_ATTR_MAX_VLAN_TAGS_PARSED:
        attr_list[i].value.u8 = resp.attr().max_vlan_tags_parsed();
        break;
      case SAI_MACSEC_ATTR_STATS_MODE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().stats_mode() - 1);
        break;
      case SAI_MACSEC_ATTR_PHYSICAL_BYPASS_ENABLE:
        attr_list[i].value.booldata = resp.attr().physical_bypass_enable();
        break;
      case SAI_MACSEC_ATTR_SUPPORTED_PORT_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().supported_port_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_MACSEC_ATTR_AVAILABLE_MACSEC_FLOW:
        attr_list[i].value.u32 = resp.attr().available_macsec_flow();
        break;
      case SAI_MACSEC_ATTR_FLOW_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().flow_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_MACSEC_ATTR_AVAILABLE_MACSEC_SC:
        attr_list[i].value.u32 = resp.attr().available_macsec_sc();
        break;
      case SAI_MACSEC_ATTR_AVAILABLE_MACSEC_SA:
        attr_list[i].value.u32 = resp.attr().available_macsec_sa();
        break;
      case SAI_MACSEC_ATTR_MAX_SECURE_ASSOCIATIONS_PER_SC:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().max_secure_associations_per_sc() - 1);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_macsec_port(sai_object_id_t *macsec_port_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateMacsecPortRequest req =
      convert_create_macsec_port(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateMacsecPortResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = macsec->CreateMacsecPort(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *macsec_port_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_macsec_port(sai_object_id_t macsec_port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveMacsecPortRequest req;
  lemming::dataplane::sai::RemoveMacsecPortResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_port_id);

  grpc::Status status = macsec->RemoveMacsecPort(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_macsec_port_attribute(sai_object_id_t macsec_port_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetMacsecPortAttributeRequest req;
  lemming::dataplane::sai::SetMacsecPortAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_port_id);

  switch (attr->id) {
    case SAI_MACSEC_PORT_ATTR_CTAG_ENABLE:
      req.set_ctag_enable(attr->value.booldata);
      break;
    case SAI_MACSEC_PORT_ATTR_STAG_ENABLE:
      req.set_stag_enable(attr->value.booldata);
      break;
    case SAI_MACSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
      req.set_switch_switching_mode(
          static_cast<lemming::dataplane::sai::SwitchSwitchingMode>(
              attr->value.s32 + 1));
      break;
  }

  grpc::Status status = macsec->SetMacsecPortAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_port_attribute(sai_object_id_t macsec_port_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMacsecPortAttributeRequest req;
  lemming::dataplane::sai::GetMacsecPortAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(macsec_port_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::MacsecPortAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = macsec->GetMacsecPortAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MACSEC_PORT_ATTR_MACSEC_DIRECTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().macsec_direction() - 1);
        break;
      case SAI_MACSEC_PORT_ATTR_PORT_ID:
        attr_list[i].value.oid = resp.attr().port_id();
        break;
      case SAI_MACSEC_PORT_ATTR_CTAG_ENABLE:
        attr_list[i].value.booldata = resp.attr().ctag_enable();
        break;
      case SAI_MACSEC_PORT_ATTR_STAG_ENABLE:
        attr_list[i].value.booldata = resp.attr().stag_enable();
        break;
      case SAI_MACSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().switch_switching_mode() - 1);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_port_stats(sai_object_id_t macsec_port_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMacsecPortStatsRequest req;
  lemming::dataplane::sai::GetMacsecPortStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_port_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(static_cast<lemming::dataplane::sai::MacsecPortStat>(
        counter_ids[i] + 1));
  }
  grpc::Status status = macsec->GetMacsecPortStats(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_port_stats_ext(sai_object_id_t macsec_port_id,
                                         uint32_t number_of_counters,
                                         const sai_stat_id_t *counter_ids,
                                         sai_stats_mode_t mode,
                                         uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_macsec_port_stats(sai_object_id_t macsec_port_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_macsec_flow(sai_object_id_t *macsec_flow_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateMacsecFlowRequest req =
      convert_create_macsec_flow(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateMacsecFlowResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = macsec->CreateMacsecFlow(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *macsec_flow_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_macsec_flow(sai_object_id_t macsec_flow_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveMacsecFlowRequest req;
  lemming::dataplane::sai::RemoveMacsecFlowResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_flow_id);

  grpc::Status status = macsec->RemoveMacsecFlow(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_macsec_flow_attribute(sai_object_id_t macsec_flow_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_flow_attribute(sai_object_id_t macsec_flow_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMacsecFlowAttributeRequest req;
  lemming::dataplane::sai::GetMacsecFlowAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(macsec_flow_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::MacsecFlowAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = macsec->GetMacsecFlowAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MACSEC_FLOW_ATTR_MACSEC_DIRECTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().macsec_direction() - 1);
        break;
      case SAI_MACSEC_FLOW_ATTR_ACL_ENTRY_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().acl_entry_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_MACSEC_FLOW_ATTR_SC_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().sc_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_flow_stats(sai_object_id_t macsec_flow_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids,
                                     uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMacsecFlowStatsRequest req;
  lemming::dataplane::sai::GetMacsecFlowStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_flow_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(static_cast<lemming::dataplane::sai::MacsecFlowStat>(
        counter_ids[i] + 1));
  }
  grpc::Status status = macsec->GetMacsecFlowStats(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_flow_stats_ext(sai_object_id_t macsec_flow_id,
                                         uint32_t number_of_counters,
                                         const sai_stat_id_t *counter_ids,
                                         sai_stats_mode_t mode,
                                         uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_macsec_flow_stats(sai_object_id_t macsec_flow_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_macsec_sc(sai_object_id_t *macsec_sc_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateMacsecScRequest req =
      convert_create_macsec_sc(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateMacsecScResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = macsec->CreateMacsecSc(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *macsec_sc_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_macsec_sc(sai_object_id_t macsec_sc_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveMacsecScRequest req;
  lemming::dataplane::sai::RemoveMacsecScResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_sc_id);

  grpc::Status status = macsec->RemoveMacsecSc(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_macsec_sc_attribute(sai_object_id_t macsec_sc_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetMacsecScAttributeRequest req;
  lemming::dataplane::sai::SetMacsecScAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_sc_id);

  switch (attr->id) {
    case SAI_MACSEC_SC_ATTR_MACSEC_EXPLICIT_SCI_ENABLE:
      req.set_macsec_explicit_sci_enable(attr->value.booldata);
      break;
    case SAI_MACSEC_SC_ATTR_MACSEC_SECTAG_OFFSET:
      req.set_macsec_sectag_offset(attr->value.u8);
      break;
    case SAI_MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_ENABLE:
      req.set_macsec_replay_protection_enable(attr->value.booldata);
      break;
    case SAI_MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_WINDOW:
      req.set_macsec_replay_protection_window(attr->value.u32);
      break;
    case SAI_MACSEC_SC_ATTR_MACSEC_CIPHER_SUITE:
      req.set_macsec_cipher_suite(
          static_cast<lemming::dataplane::sai::MacsecCipherSuite>(
              attr->value.s32 + 1));
      break;
    case SAI_MACSEC_SC_ATTR_ENCRYPTION_ENABLE:
      req.set_encryption_enable(attr->value.booldata);
      break;
  }

  grpc::Status status = macsec->SetMacsecScAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_sc_attribute(sai_object_id_t macsec_sc_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMacsecScAttributeRequest req;
  lemming::dataplane::sai::GetMacsecScAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(macsec_sc_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::MacsecScAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = macsec->GetMacsecScAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MACSEC_SC_ATTR_MACSEC_DIRECTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().macsec_direction() - 1);
        break;
      case SAI_MACSEC_SC_ATTR_FLOW_ID:
        attr_list[i].value.oid = resp.attr().flow_id();
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_SCI:
        attr_list[i].value.u64 = resp.attr().macsec_sci();
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_EXPLICIT_SCI_ENABLE:
        attr_list[i].value.booldata = resp.attr().macsec_explicit_sci_enable();
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_SECTAG_OFFSET:
        attr_list[i].value.u8 = resp.attr().macsec_sectag_offset();
        break;
      case SAI_MACSEC_SC_ATTR_ACTIVE_EGRESS_SA_ID:
        attr_list[i].value.oid = resp.attr().active_egress_sa_id();
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_ENABLE:
        attr_list[i].value.booldata =
            resp.attr().macsec_replay_protection_enable();
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_WINDOW:
        attr_list[i].value.u32 = resp.attr().macsec_replay_protection_window();
        break;
      case SAI_MACSEC_SC_ATTR_SA_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().sa_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_MACSEC_SC_ATTR_MACSEC_CIPHER_SUITE:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().macsec_cipher_suite() - 1);
        break;
      case SAI_MACSEC_SC_ATTR_ENCRYPTION_ENABLE:
        attr_list[i].value.booldata = resp.attr().encryption_enable();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_sc_stats(sai_object_id_t macsec_sc_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids,
                                   uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMacsecScStatsRequest req;
  lemming::dataplane::sai::GetMacsecScStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_sc_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(
        static_cast<lemming::dataplane::sai::MacsecScStat>(counter_ids[i] + 1));
  }
  grpc::Status status = macsec->GetMacsecScStats(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_sc_stats_ext(sai_object_id_t macsec_sc_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids,
                                       sai_stats_mode_t mode,
                                       uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_macsec_sc_stats(sai_object_id_t macsec_sc_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_macsec_sa(sai_object_id_t *macsec_sa_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateMacsecSaRequest req =
      convert_create_macsec_sa(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateMacsecSaResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = macsec->CreateMacsecSa(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *macsec_sa_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_macsec_sa(sai_object_id_t macsec_sa_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveMacsecSaRequest req;
  lemming::dataplane::sai::RemoveMacsecSaResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_sa_id);

  grpc::Status status = macsec->RemoveMacsecSa(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_macsec_sa_attribute(sai_object_id_t macsec_sa_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetMacsecSaAttributeRequest req;
  lemming::dataplane::sai::SetMacsecSaAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_sa_id);

  switch (attr->id) {
    case SAI_MACSEC_SA_ATTR_CONFIGURED_EGRESS_XPN:
      req.set_configured_egress_xpn(attr->value.u64);
      break;
    case SAI_MACSEC_SA_ATTR_MINIMUM_INGRESS_XPN:
      req.set_minimum_ingress_xpn(attr->value.u64);
      break;
  }

  grpc::Status status = macsec->SetMacsecSaAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_sa_attribute(sai_object_id_t macsec_sa_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMacsecSaAttributeRequest req;
  lemming::dataplane::sai::GetMacsecSaAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(macsec_sa_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::MacsecSaAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = macsec->GetMacsecSaAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MACSEC_SA_ATTR_MACSEC_DIRECTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().macsec_direction() - 1);
        break;
      case SAI_MACSEC_SA_ATTR_SC_ID:
        attr_list[i].value.oid = resp.attr().sc_id();
        break;
      case SAI_MACSEC_SA_ATTR_AN:
        attr_list[i].value.u8 = resp.attr().an();
        break;
      case SAI_MACSEC_SA_ATTR_CONFIGURED_EGRESS_XPN:
        attr_list[i].value.u64 = resp.attr().configured_egress_xpn();
        break;
      case SAI_MACSEC_SA_ATTR_CURRENT_XPN:
        attr_list[i].value.u64 = resp.attr().current_xpn();
        break;
      case SAI_MACSEC_SA_ATTR_MINIMUM_INGRESS_XPN:
        attr_list[i].value.u64 = resp.attr().minimum_ingress_xpn();
        break;
      case SAI_MACSEC_SA_ATTR_MACSEC_SSCI:
        attr_list[i].value.u32 = resp.attr().macsec_ssci();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_sa_stats(sai_object_id_t macsec_sa_id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t *counter_ids,
                                   uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMacsecSaStatsRequest req;
  lemming::dataplane::sai::GetMacsecSaStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(macsec_sa_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(
        static_cast<lemming::dataplane::sai::MacsecSaStat>(counter_ids[i] + 1));
  }
  grpc::Status status = macsec->GetMacsecSaStats(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_macsec_sa_stats_ext(sai_object_id_t macsec_sa_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t *counter_ids,
                                       sai_stats_mode_t mode,
                                       uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_macsec_sa_stats(sai_object_id_t macsec_sa_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
