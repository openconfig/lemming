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

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/ipsec.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

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

lemming::dataplane::sai::CreateIpsecRequest convert_create_ipsec(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t* attr_list) {
  lemming::dataplane::sai::CreateIpsecRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_IPSEC_ATTR_WARM_BOOT_ENABLE:
        msg.set_warm_boot_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_ATTR_EXTERNAL_SA_INDEX_ENABLE:
        msg.set_external_sa_index_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_ATTR_CTAG_TPID:
        msg.set_ctag_tpid(attr_list[i].value.u16);
        break;
      case SAI_IPSEC_ATTR_STAG_TPID:
        msg.set_stag_tpid(attr_list[i].value.u16);
        break;
      case SAI_IPSEC_ATTR_MAX_VLAN_TAGS_PARSED:
        msg.set_max_vlan_tags_parsed(attr_list[i].value.u8);
        break;
      case SAI_IPSEC_ATTR_OCTET_COUNT_HIGH_WATERMARK:
        msg.set_octet_count_high_watermark(attr_list[i].value.u64);
        break;
      case SAI_IPSEC_ATTR_OCTET_COUNT_LOW_WATERMARK:
        msg.set_octet_count_low_watermark(attr_list[i].value.u64);
        break;
      case SAI_IPSEC_ATTR_STATS_MODE:
        msg.set_stats_mode(
            convert_sai_stats_mode_t_to_proto(attr_list[i].value.s32));
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateIpsecPortRequest convert_create_ipsec_port(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t* attr_list) {
  lemming::dataplane::sai::CreateIpsecPortRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_IPSEC_PORT_ATTR_PORT_ID:
        msg.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_IPSEC_PORT_ATTR_CTAG_ENABLE:
        msg.set_ctag_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_PORT_ATTR_STAG_ENABLE:
        msg.set_stag_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_PORT_ATTR_NATIVE_VLAN_ID:
        msg.set_native_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_IPSEC_PORT_ATTR_VRF_FROM_PACKET_VLAN_ENABLE:
        msg.set_vrf_from_packet_vlan_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
        msg.set_switch_switching_mode(
            convert_sai_switch_switching_mode_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_IPSEC_PORT_ATTR_STATS_COUNT_MODE:
        msg.set_stats_count_mode(
            convert_sai_stats_count_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_IPSEC_PORT_ATTR_SELECTIVE_COUNTER_LIST:
        msg.mutable_selective_counter_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateIpsecSaRequest convert_create_ipsec_sa(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t* attr_list) {
  lemming::dataplane::sai::CreateIpsecSaRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_IPSEC_SA_ATTR_IPSEC_DIRECTION:
        msg.set_ipsec_direction(
            convert_sai_ipsec_direction_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_ID:
        msg.set_ipsec_id(attr_list[i].value.oid);
        break;
      case SAI_IPSEC_SA_ATTR_EXTERNAL_SA_INDEX:
        msg.set_external_sa_index(attr_list[i].value.u32);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_PORT_LIST:
        msg.mutable_ipsec_port_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_SPI:
        msg.set_ipsec_spi(attr_list[i].value.u32);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_ESN_ENABLE:
        msg.set_ipsec_esn_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_CIPHER:
        msg.set_ipsec_cipher(
            convert_sai_ipsec_cipher_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_IPSEC_SA_ATTR_SALT:
        msg.set_salt(attr_list[i].value.u32);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_ENABLE:
        msg.set_ipsec_replay_protection_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_WINDOW:
        msg.set_ipsec_replay_protection_window(attr_list[i].value.u32);
        break;
      case SAI_IPSEC_SA_ATTR_TERM_DST_IP:
        msg.set_term_dst_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_IPSEC_SA_ATTR_TERM_VLAN_ID_ENABLE:
        msg.set_term_vlan_id_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_SA_ATTR_TERM_VLAN_ID:
        msg.set_term_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_IPSEC_SA_ATTR_TERM_SRC_IP_ENABLE:
        msg.set_term_src_ip_enable(attr_list[i].value.booldata);
        break;
      case SAI_IPSEC_SA_ATTR_TERM_SRC_IP:
        msg.set_term_src_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_IPSEC_SA_ATTR_EGRESS_ESN:
        msg.set_egress_esn(attr_list[i].value.u64);
        break;
      case SAI_IPSEC_SA_ATTR_MINIMUM_INGRESS_ESN:
        msg.set_minimum_ingress_esn(attr_list[i].value.u64);
        break;
      case SAI_IPSEC_SA_ATTR_STATS_COUNT_MODE:
        msg.set_stats_count_mode(
            convert_sai_stats_count_mode_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_IPSEC_SA_ATTR_SELECTIVE_COUNTER_LIST:
        msg.mutable_selective_counter_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_ipsec(sai_object_id_t* ipsec_id,
                            sai_object_id_t switch_id, uint32_t attr_count,
                            const sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIpsecRequest req =
      convert_create_ipsec(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateIpsecResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = ipsec->CreateIpsec(&context, req, &resp);
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
  if (ipsec_id) {
    *ipsec_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_ipsec(sai_object_id_t ipsec_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveIpsecRequest req;
  lemming::dataplane::sai::RemoveIpsecResponse resp;
  grpc::ClientContext context;
  req.set_oid(ipsec_id);

  grpc::Status status = ipsec->RemoveIpsec(&context, req, &resp);
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

sai_status_t l_set_ipsec_attribute(sai_object_id_t ipsec_id,
                                   const sai_attribute_t* attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetIpsecAttributeRequest req;
  lemming::dataplane::sai::SetIpsecAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(ipsec_id);

  switch (attr->id) {
    case SAI_IPSEC_ATTR_WARM_BOOT_ENABLE:
      req.set_warm_boot_enable(attr->value.booldata);
      break;
    case SAI_IPSEC_ATTR_CTAG_TPID:
      req.set_ctag_tpid(attr->value.u16);
      break;
    case SAI_IPSEC_ATTR_STAG_TPID:
      req.set_stag_tpid(attr->value.u16);
      break;
    case SAI_IPSEC_ATTR_MAX_VLAN_TAGS_PARSED:
      req.set_max_vlan_tags_parsed(attr->value.u8);
      break;
    case SAI_IPSEC_ATTR_OCTET_COUNT_HIGH_WATERMARK:
      req.set_octet_count_high_watermark(attr->value.u64);
      break;
    case SAI_IPSEC_ATTR_OCTET_COUNT_LOW_WATERMARK:
      req.set_octet_count_low_watermark(attr->value.u64);
      break;
    case SAI_IPSEC_ATTR_STATS_MODE:
      req.set_stats_mode(convert_sai_stats_mode_t_to_proto(attr->value.s32));
      break;
  }

  grpc::Status status = ipsec->SetIpsecAttribute(&context, req, &resp);
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

sai_status_t l_get_ipsec_attribute(sai_object_id_t ipsec_id,
                                   uint32_t attr_count,
                                   sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIpsecAttributeRequest req;
  lemming::dataplane::sai::GetIpsecAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(ipsec_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_ipsec_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = ipsec->GetIpsecAttribute(&context, req, &resp);
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
      case SAI_IPSEC_ATTR_TERM_REMOTE_IP_MATCH_SUPPORTED:
        attr_list[i].value.booldata =
            resp.attr().term_remote_ip_match_supported();
        break;
      case SAI_IPSEC_ATTR_SWITCHING_MODE_CUT_THROUGH_SUPPORTED:
        attr_list[i].value.booldata =
            resp.attr().switching_mode_cut_through_supported();
        break;
      case SAI_IPSEC_ATTR_SWITCHING_MODE_STORE_AND_FORWARD_SUPPORTED:
        attr_list[i].value.booldata =
            resp.attr().switching_mode_store_and_forward_supported();
        break;
      case SAI_IPSEC_ATTR_STATS_MODE_READ_SUPPORTED:
        attr_list[i].value.booldata = resp.attr().stats_mode_read_supported();
        break;
      case SAI_IPSEC_ATTR_STATS_MODE_READ_CLEAR_SUPPORTED:
        attr_list[i].value.booldata =
            resp.attr().stats_mode_read_clear_supported();
        break;
      case SAI_IPSEC_ATTR_SN_32BIT_SUPPORTED:
        attr_list[i].value.booldata = resp.attr().sn_32bit_supported();
        break;
      case SAI_IPSEC_ATTR_ESN_64BIT_SUPPORTED:
        attr_list[i].value.booldata = resp.attr().esn_64bit_supported();
        break;
      case SAI_IPSEC_ATTR_SUPPORTED_CIPHER_LIST:
        convert_list_sai_ipsec_cipher_t_to_sai(
            attr_list[i].value.s32list.list,
            resp.attr().supported_cipher_list(),
            &attr_list[i].value.s32list.count);
        break;
      case SAI_IPSEC_ATTR_SYSTEM_SIDE_MTU:
        attr_list[i].value.u16 = resp.attr().system_side_mtu();
        break;
      case SAI_IPSEC_ATTR_WARM_BOOT_SUPPORTED:
        attr_list[i].value.booldata = resp.attr().warm_boot_supported();
        break;
      case SAI_IPSEC_ATTR_WARM_BOOT_ENABLE:
        attr_list[i].value.booldata = resp.attr().warm_boot_enable();
        break;
      case SAI_IPSEC_ATTR_EXTERNAL_SA_INDEX_ENABLE:
        attr_list[i].value.booldata = resp.attr().external_sa_index_enable();
        break;
      case SAI_IPSEC_ATTR_CTAG_TPID:
        attr_list[i].value.u16 = resp.attr().ctag_tpid();
        break;
      case SAI_IPSEC_ATTR_STAG_TPID:
        attr_list[i].value.u16 = resp.attr().stag_tpid();
        break;
      case SAI_IPSEC_ATTR_MAX_VLAN_TAGS_PARSED:
        attr_list[i].value.u8 = resp.attr().max_vlan_tags_parsed();
        break;
      case SAI_IPSEC_ATTR_OCTET_COUNT_HIGH_WATERMARK:
        attr_list[i].value.u64 = resp.attr().octet_count_high_watermark();
        break;
      case SAI_IPSEC_ATTR_OCTET_COUNT_LOW_WATERMARK:
        attr_list[i].value.u64 = resp.attr().octet_count_low_watermark();
        break;
      case SAI_IPSEC_ATTR_STATS_MODE:
        attr_list[i].value.s32 =
            convert_sai_stats_mode_t_to_sai(resp.attr().stats_mode());
        break;
      case SAI_IPSEC_ATTR_AVAILABLE_IPSEC_SA:
        attr_list[i].value.u32 = resp.attr().available_ipsec_sa();
        break;
      case SAI_IPSEC_ATTR_SA_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().sa_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_ipsec_port(sai_object_id_t* ipsec_port_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIpsecPortRequest req =
      convert_create_ipsec_port(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateIpsecPortResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = ipsec->CreateIpsecPort(&context, req, &resp);
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
  if (ipsec_port_id) {
    *ipsec_port_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_ipsec_port(sai_object_id_t ipsec_port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveIpsecPortRequest req;
  lemming::dataplane::sai::RemoveIpsecPortResponse resp;
  grpc::ClientContext context;
  req.set_oid(ipsec_port_id);

  grpc::Status status = ipsec->RemoveIpsecPort(&context, req, &resp);
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

sai_status_t l_set_ipsec_port_attribute(sai_object_id_t ipsec_port_id,
                                        const sai_attribute_t* attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetIpsecPortAttributeRequest req;
  lemming::dataplane::sai::SetIpsecPortAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(ipsec_port_id);

  switch (attr->id) {
    case SAI_IPSEC_PORT_ATTR_CTAG_ENABLE:
      req.set_ctag_enable(attr->value.booldata);
      break;
    case SAI_IPSEC_PORT_ATTR_STAG_ENABLE:
      req.set_stag_enable(attr->value.booldata);
      break;
    case SAI_IPSEC_PORT_ATTR_VRF_FROM_PACKET_VLAN_ENABLE:
      req.set_vrf_from_packet_vlan_enable(attr->value.booldata);
      break;
    case SAI_IPSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
      req.set_switch_switching_mode(
          convert_sai_switch_switching_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_IPSEC_PORT_ATTR_STATS_COUNT_MODE:
      req.set_stats_count_mode(
          convert_sai_stats_count_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_IPSEC_PORT_ATTR_SELECTIVE_COUNTER_LIST:
      req.mutable_selective_counter_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
  }

  grpc::Status status = ipsec->SetIpsecPortAttribute(&context, req, &resp);
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

sai_status_t l_get_ipsec_port_attribute(sai_object_id_t ipsec_port_id,
                                        uint32_t attr_count,
                                        sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIpsecPortAttributeRequest req;
  lemming::dataplane::sai::GetIpsecPortAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(ipsec_port_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_ipsec_port_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = ipsec->GetIpsecPortAttribute(&context, req, &resp);
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
      case SAI_IPSEC_PORT_ATTR_PORT_ID:
        attr_list[i].value.oid = resp.attr().port_id();
        break;
      case SAI_IPSEC_PORT_ATTR_CTAG_ENABLE:
        attr_list[i].value.booldata = resp.attr().ctag_enable();
        break;
      case SAI_IPSEC_PORT_ATTR_STAG_ENABLE:
        attr_list[i].value.booldata = resp.attr().stag_enable();
        break;
      case SAI_IPSEC_PORT_ATTR_NATIVE_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().native_vlan_id();
        break;
      case SAI_IPSEC_PORT_ATTR_VRF_FROM_PACKET_VLAN_ENABLE:
        attr_list[i].value.booldata = resp.attr().vrf_from_packet_vlan_enable();
        break;
      case SAI_IPSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
        attr_list[i].value.s32 = convert_sai_switch_switching_mode_t_to_sai(
            resp.attr().switch_switching_mode());
        break;
      case SAI_IPSEC_PORT_ATTR_STATS_COUNT_MODE:
        attr_list[i].value.s32 = convert_sai_stats_count_mode_t_to_sai(
            resp.attr().stats_count_mode());
        break;
      case SAI_IPSEC_PORT_ATTR_SELECTIVE_COUNTER_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().selective_counter_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_ipsec_port_stats(sai_object_id_t ipsec_port_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t* counter_ids,
                                    uint64_t* counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIpsecPortStatsRequest req;
  lemming::dataplane::sai::GetIpsecPortStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(ipsec_port_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(convert_sai_ipsec_port_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status = ipsec->GetIpsecPortStats(&context, req, &resp);
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

sai_status_t l_get_ipsec_port_stats_ext(sai_object_id_t ipsec_port_id,
                                        uint32_t number_of_counters,
                                        const sai_stat_id_t* counter_ids,
                                        sai_stats_mode_t mode,
                                        uint64_t* counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_ipsec_port_stats(sai_object_id_t ipsec_port_id,
                                      uint32_t number_of_counters,
                                      const sai_stat_id_t* counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_ipsec_sa(sai_object_id_t* ipsec_sa_id,
                               sai_object_id_t switch_id, uint32_t attr_count,
                               const sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIpsecSaRequest req =
      convert_create_ipsec_sa(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateIpsecSaResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = ipsec->CreateIpsecSa(&context, req, &resp);
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
  if (ipsec_sa_id) {
    *ipsec_sa_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_ipsec_sa(sai_object_id_t ipsec_sa_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveIpsecSaRequest req;
  lemming::dataplane::sai::RemoveIpsecSaResponse resp;
  grpc::ClientContext context;
  req.set_oid(ipsec_sa_id);

  grpc::Status status = ipsec->RemoveIpsecSa(&context, req, &resp);
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

sai_status_t l_set_ipsec_sa_attribute(sai_object_id_t ipsec_sa_id,
                                      const sai_attribute_t* attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetIpsecSaAttributeRequest req;
  lemming::dataplane::sai::SetIpsecSaAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(ipsec_sa_id);

  switch (attr->id) {
    case SAI_IPSEC_SA_ATTR_EXTERNAL_SA_INDEX:
      req.set_external_sa_index(attr->value.u32);
      break;
    case SAI_IPSEC_SA_ATTR_IPSEC_PORT_LIST:
      req.mutable_ipsec_port_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_ENABLE:
      req.set_ipsec_replay_protection_enable(attr->value.booldata);
      break;
    case SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_WINDOW:
      req.set_ipsec_replay_protection_window(attr->value.u32);
      break;
    case SAI_IPSEC_SA_ATTR_EGRESS_ESN:
      req.set_egress_esn(attr->value.u64);
      break;
    case SAI_IPSEC_SA_ATTR_MINIMUM_INGRESS_ESN:
      req.set_minimum_ingress_esn(attr->value.u64);
      break;
    case SAI_IPSEC_SA_ATTR_STATS_COUNT_MODE:
      req.set_stats_count_mode(
          convert_sai_stats_count_mode_t_to_proto(attr->value.s32));
      break;
    case SAI_IPSEC_SA_ATTR_SELECTIVE_COUNTER_LIST:
      req.mutable_selective_counter_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
  }

  grpc::Status status = ipsec->SetIpsecSaAttribute(&context, req, &resp);
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

sai_status_t l_get_ipsec_sa_attribute(sai_object_id_t ipsec_sa_id,
                                      uint32_t attr_count,
                                      sai_attribute_t* attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIpsecSaAttributeRequest req;
  lemming::dataplane::sai::GetIpsecSaAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(ipsec_sa_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_ipsec_sa_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = ipsec->GetIpsecSaAttribute(&context, req, &resp);
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
      case SAI_IPSEC_SA_ATTR_IPSEC_DIRECTION:
        attr_list[i].value.s32 =
            convert_sai_ipsec_direction_t_to_sai(resp.attr().ipsec_direction());
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_ID:
        attr_list[i].value.oid = resp.attr().ipsec_id();
        break;
      case SAI_IPSEC_SA_ATTR_OCTET_COUNT_STATUS:
        attr_list[i].value.s32 =
            convert_sai_ipsec_sa_octet_count_status_t_to_sai(
                resp.attr().octet_count_status());
        break;
      case SAI_IPSEC_SA_ATTR_EXTERNAL_SA_INDEX:
        attr_list[i].value.u32 = resp.attr().external_sa_index();
        break;
      case SAI_IPSEC_SA_ATTR_SA_INDEX:
        attr_list[i].value.u32 = resp.attr().sa_index();
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_PORT_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().ipsec_port_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_SPI:
        attr_list[i].value.u32 = resp.attr().ipsec_spi();
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_ESN_ENABLE:
        attr_list[i].value.booldata = resp.attr().ipsec_esn_enable();
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_CIPHER:
        attr_list[i].value.s32 =
            convert_sai_ipsec_cipher_t_to_sai(resp.attr().ipsec_cipher());
        break;
      case SAI_IPSEC_SA_ATTR_SALT:
        attr_list[i].value.u32 = resp.attr().salt();
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_ENABLE:
        attr_list[i].value.booldata =
            resp.attr().ipsec_replay_protection_enable();
        break;
      case SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_WINDOW:
        attr_list[i].value.u32 = resp.attr().ipsec_replay_protection_window();
        break;
      case SAI_IPSEC_SA_ATTR_TERM_DST_IP:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().term_dst_ip());
        break;
      case SAI_IPSEC_SA_ATTR_TERM_VLAN_ID_ENABLE:
        attr_list[i].value.booldata = resp.attr().term_vlan_id_enable();
        break;
      case SAI_IPSEC_SA_ATTR_TERM_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().term_vlan_id();
        break;
      case SAI_IPSEC_SA_ATTR_TERM_SRC_IP_ENABLE:
        attr_list[i].value.booldata = resp.attr().term_src_ip_enable();
        break;
      case SAI_IPSEC_SA_ATTR_TERM_SRC_IP:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().term_src_ip());
        break;
      case SAI_IPSEC_SA_ATTR_EGRESS_ESN:
        attr_list[i].value.u64 = resp.attr().egress_esn();
        break;
      case SAI_IPSEC_SA_ATTR_MINIMUM_INGRESS_ESN:
        attr_list[i].value.u64 = resp.attr().minimum_ingress_esn();
        break;
      case SAI_IPSEC_SA_ATTR_STATS_COUNT_MODE:
        attr_list[i].value.s32 = convert_sai_stats_count_mode_t_to_sai(
            resp.attr().stats_count_mode());
        break;
      case SAI_IPSEC_SA_ATTR_SELECTIVE_COUNTER_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().selective_counter_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_ipsec_sa_stats(sai_object_id_t ipsec_sa_id,
                                  uint32_t number_of_counters,
                                  const sai_stat_id_t* counter_ids,
                                  uint64_t* counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIpsecSaStatsRequest req;
  lemming::dataplane::sai::GetIpsecSaStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(ipsec_sa_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(convert_sai_ipsec_sa_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status = ipsec->GetIpsecSaStats(&context, req, &resp);
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

sai_status_t l_get_ipsec_sa_stats_ext(sai_object_id_t ipsec_sa_id,
                                      uint32_t number_of_counters,
                                      const sai_stat_id_t* counter_ids,
                                      sai_stats_mode_t mode,
                                      uint64_t* counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_ipsec_sa_stats(sai_object_id_t ipsec_sa_id,
                                    uint32_t number_of_counters,
                                    const sai_stat_id_t* counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
