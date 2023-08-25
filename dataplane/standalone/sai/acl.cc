

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

#include "dataplane/standalone/sai/acl.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/acl.pb.h"
#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_acl_api_t l_acl = {
    .create_acl_table = l_create_acl_table,
    .remove_acl_table = l_remove_acl_table,
    .set_acl_table_attribute = l_set_acl_table_attribute,
    .get_acl_table_attribute = l_get_acl_table_attribute,
    .create_acl_entry = l_create_acl_entry,
    .remove_acl_entry = l_remove_acl_entry,
    .set_acl_entry_attribute = l_set_acl_entry_attribute,
    .get_acl_entry_attribute = l_get_acl_entry_attribute,
    .create_acl_counter = l_create_acl_counter,
    .remove_acl_counter = l_remove_acl_counter,
    .set_acl_counter_attribute = l_set_acl_counter_attribute,
    .get_acl_counter_attribute = l_get_acl_counter_attribute,
    .create_acl_range = l_create_acl_range,
    .remove_acl_range = l_remove_acl_range,
    .set_acl_range_attribute = l_set_acl_range_attribute,
    .get_acl_range_attribute = l_get_acl_range_attribute,
    .create_acl_table_group = l_create_acl_table_group,
    .remove_acl_table_group = l_remove_acl_table_group,
    .set_acl_table_group_attribute = l_set_acl_table_group_attribute,
    .get_acl_table_group_attribute = l_get_acl_table_group_attribute,
    .create_acl_table_group_member = l_create_acl_table_group_member,
    .remove_acl_table_group_member = l_remove_acl_table_group_member,
    .set_acl_table_group_member_attribute =
        l_set_acl_table_group_member_attribute,
    .get_acl_table_group_member_attribute =
        l_get_acl_table_group_member_attribute,
};

sai_status_t l_create_acl_table(sai_object_id_t *acl_table_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateAclTableRequest req;
  lemming::dataplane::sai::CreateAclTableResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ACL_TABLE_ATTR_ACL_STAGE:
        req.set_acl_stage(static_cast<lemming::dataplane::sai::AclStage>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_ACL_TABLE_ATTR_SIZE:
        req.set_size(attr_list[i].value.u32);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6:
        req.set_field_src_ipv6(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD3:
        req.set_field_src_ipv6_word3(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD2:
        req.set_field_src_ipv6_word2(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD1:
        req.set_field_src_ipv6_word1(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD0:
        req.set_field_src_ipv6_word0(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6:
        req.set_field_dst_ipv6(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD3:
        req.set_field_dst_ipv6_word3(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD2:
        req.set_field_dst_ipv6_word2(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD1:
        req.set_field_dst_ipv6_word1(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD0:
        req.set_field_dst_ipv6_word0(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_IPV6:
        req.set_field_inner_src_ipv6(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_IPV6:
        req.set_field_inner_dst_ipv6(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_SRC_MAC:
        req.set_field_src_mac(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_DST_MAC:
        req.set_field_dst_mac(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_SRC_IP:
        req.set_field_src_ip(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_DST_IP:
        req.set_field_dst_ip(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_IP:
        req.set_field_inner_src_ip(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_IP:
        req.set_field_inner_dst_ip(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_IN_PORTS:
        req.set_field_in_ports(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_OUT_PORTS:
        req.set_field_out_ports(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_IN_PORT:
        req.set_field_in_port(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_OUT_PORT:
        req.set_field_out_port(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_SRC_PORT:
        req.set_field_src_port(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_ID:
        req.set_field_outer_vlan_id(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_PRI:
        req.set_field_outer_vlan_pri(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_CFI:
        req.set_field_outer_vlan_cfi(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_ID:
        req.set_field_inner_vlan_id(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_PRI:
        req.set_field_inner_vlan_pri(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_CFI:
        req.set_field_inner_vlan_cfi(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_L4_SRC_PORT:
        req.set_field_l4_src_port(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_L4_DST_PORT:
        req.set_field_l4_dst_port(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_L4_SRC_PORT:
        req.set_field_inner_l4_src_port(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_L4_DST_PORT:
        req.set_field_inner_l4_dst_port(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ETHER_TYPE:
        req.set_field_ether_type(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_ETHER_TYPE:
        req.set_field_inner_ether_type(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_IP_PROTOCOL:
        req.set_field_ip_protocol(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_INNER_IP_PROTOCOL:
        req.set_field_inner_ip_protocol(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_IP_IDENTIFICATION:
        req.set_field_ip_identification(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_DSCP:
        req.set_field_dscp(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ECN:
        req.set_field_ecn(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_TTL:
        req.set_field_ttl(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_TOS:
        req.set_field_tos(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_IP_FLAGS:
        req.set_field_ip_flags(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_TCP_FLAGS:
        req.set_field_tcp_flags(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ACL_IP_TYPE:
        req.set_field_acl_ip_type(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ACL_IP_FRAG:
        req.set_field_acl_ip_frag(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_IPV6_FLOW_LABEL:
        req.set_field_ipv6_flow_label(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_TC:
        req.set_field_tc(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ICMP_TYPE:
        req.set_field_icmp_type(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ICMP_CODE:
        req.set_field_icmp_code(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ICMPV6_TYPE:
        req.set_field_icmpv6_type(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ICMPV6_CODE:
        req.set_field_icmpv6_code(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_PACKET_VLAN:
        req.set_field_packet_vlan(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_TUNNEL_VNI:
        req.set_field_tunnel_vni(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_HAS_VLAN_TAG:
        req.set_field_has_vlan_tag(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MACSEC_SCI:
        req.set_field_macsec_sci(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_LABEL:
        req.set_field_mpls_label0_label(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_TTL:
        req.set_field_mpls_label0_ttl(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_EXP:
        req.set_field_mpls_label0_exp(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_BOS:
        req.set_field_mpls_label0_bos(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_LABEL:
        req.set_field_mpls_label1_label(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_TTL:
        req.set_field_mpls_label1_ttl(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_EXP:
        req.set_field_mpls_label1_exp(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_BOS:
        req.set_field_mpls_label1_bos(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_LABEL:
        req.set_field_mpls_label2_label(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_TTL:
        req.set_field_mpls_label2_ttl(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_EXP:
        req.set_field_mpls_label2_exp(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_BOS:
        req.set_field_mpls_label2_bos(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_LABEL:
        req.set_field_mpls_label3_label(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_TTL:
        req.set_field_mpls_label3_ttl(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_EXP:
        req.set_field_mpls_label3_exp(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_BOS:
        req.set_field_mpls_label3_bos(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_LABEL:
        req.set_field_mpls_label4_label(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_TTL:
        req.set_field_mpls_label4_ttl(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_EXP:
        req.set_field_mpls_label4_exp(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_BOS:
        req.set_field_mpls_label4_bos(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_FDB_DST_USER_META:
        req.set_field_fdb_dst_user_meta(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ROUTE_DST_USER_META:
        req.set_field_route_dst_user_meta(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_NEIGHBOR_DST_USER_META:
        req.set_field_neighbor_dst_user_meta(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_PORT_USER_META:
        req.set_field_port_user_meta(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_VLAN_USER_META:
        req.set_field_vlan_user_meta(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ACL_USER_META:
        req.set_field_acl_user_meta(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_FDB_NPU_META_DST_HIT:
        req.set_field_fdb_npu_meta_dst_hit(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT:
        req.set_field_neighbor_npu_meta_dst_hit(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_ROUTE_NPU_META_DST_HIT:
        req.set_field_route_npu_meta_dst_hit(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_BTH_OPCODE:
        req.set_field_bth_opcode(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_AETH_SYNDROME:
        req.set_field_aeth_syndrome(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN:
        req.set_user_defined_field_group_min(attr_list[i].value.oid);
        break;
      case SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MAX:
        req.set_user_defined_field_group_max(attr_list[i].value.oid);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_IPV6_NEXT_HEADER:
        req.set_field_ipv6_next_header(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_GRE_KEY:
        req.set_field_gre_key(attr_list[i].value.booldata);
        break;
      case SAI_ACL_TABLE_ATTR_FIELD_TAM_INT_TYPE:
        req.set_field_tam_int_type(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status = acl->CreateAclTable(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *acl_table_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_ACL_TABLE, acl_table_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_acl_table(sai_object_id_t acl_table_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_ACL_TABLE, acl_table_id);
}

sai_status_t l_set_acl_table_attribute(sai_object_id_t acl_table_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_ACL_TABLE, acl_table_id,
                                   attr);
}

sai_status_t l_get_acl_table_attribute(sai_object_id_t acl_table_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_ACL_TABLE, acl_table_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_acl_entry(sai_object_id_t *acl_entry_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateAclEntryRequest req;
  lemming::dataplane::sai::CreateAclEntryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ACL_ENTRY_ATTR_TABLE_ID:
        req.set_table_id(attr_list[i].value.oid);
        break;
      case SAI_ACL_ENTRY_ATTR_PRIORITY:
        req.set_priority(attr_list[i].value.u32);
        break;
      case SAI_ACL_ENTRY_ATTR_ADMIN_STATE:
        req.set_admin_state(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status = acl->CreateAclEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *acl_entry_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_ACL_ENTRY, acl_entry_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_acl_entry(sai_object_id_t acl_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_ACL_ENTRY, acl_entry_id);
}

sai_status_t l_set_acl_entry_attribute(sai_object_id_t acl_entry_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_ACL_ENTRY, acl_entry_id,
                                   attr);
}

sai_status_t l_get_acl_entry_attribute(sai_object_id_t acl_entry_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_ACL_ENTRY, acl_entry_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_acl_counter(sai_object_id_t *acl_counter_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateAclCounterRequest req;
  lemming::dataplane::sai::CreateAclCounterResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ACL_COUNTER_ATTR_TABLE_ID:
        req.set_table_id(attr_list[i].value.oid);
        break;
      case SAI_ACL_COUNTER_ATTR_ENABLE_PACKET_COUNT:
        req.set_enable_packet_count(attr_list[i].value.booldata);
        break;
      case SAI_ACL_COUNTER_ATTR_ENABLE_BYTE_COUNT:
        req.set_enable_byte_count(attr_list[i].value.booldata);
        break;
      case SAI_ACL_COUNTER_ATTR_PACKETS:
        req.set_packets(attr_list[i].value.u64);
        break;
      case SAI_ACL_COUNTER_ATTR_BYTES:
        req.set_bytes(attr_list[i].value.u64);
        break;
    }
  }
  grpc::Status status = acl->CreateAclCounter(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *acl_counter_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_ACL_COUNTER, acl_counter_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_acl_counter(sai_object_id_t acl_counter_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_ACL_COUNTER, acl_counter_id);
}

sai_status_t l_set_acl_counter_attribute(sai_object_id_t acl_counter_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_ACL_COUNTER, acl_counter_id,
                                   attr);
}

sai_status_t l_get_acl_counter_attribute(sai_object_id_t acl_counter_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_ACL_COUNTER, acl_counter_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_acl_range(sai_object_id_t *acl_range_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateAclRangeRequest req;
  lemming::dataplane::sai::CreateAclRangeResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ACL_RANGE_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::AclRangeType>(
            attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = acl->CreateAclRange(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *acl_range_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_ACL_RANGE, acl_range_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_acl_range(sai_object_id_t acl_range_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_ACL_RANGE, acl_range_id);
}

sai_status_t l_set_acl_range_attribute(sai_object_id_t acl_range_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_ACL_RANGE, acl_range_id,
                                   attr);
}

sai_status_t l_get_acl_range_attribute(sai_object_id_t acl_range_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_ACL_RANGE, acl_range_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_acl_table_group(sai_object_id_t *acl_table_group_id,
                                      sai_object_id_t switch_id,
                                      uint32_t attr_count,
                                      const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateAclTableGroupRequest req;
  lemming::dataplane::sai::CreateAclTableGroupResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ACL_TABLE_GROUP_ATTR_ACL_STAGE:
        req.set_acl_stage(static_cast<lemming::dataplane::sai::AclStage>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_ACL_TABLE_GROUP_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::AclTableGroupType>(
            attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = acl->CreateAclTableGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *acl_table_group_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_ACL_TABLE_GROUP, acl_table_group_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_acl_table_group(sai_object_id_t acl_table_group_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_ACL_TABLE_GROUP,
                            acl_table_group_id);
}

sai_status_t l_set_acl_table_group_attribute(sai_object_id_t acl_table_group_id,
                                             const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_ACL_TABLE_GROUP,
                                   acl_table_group_id, attr);
}

sai_status_t l_get_acl_table_group_attribute(sai_object_id_t acl_table_group_id,
                                             uint32_t attr_count,
                                             sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_ACL_TABLE_GROUP,
                                   acl_table_group_id, attr_count, attr_list);
}

sai_status_t l_create_acl_table_group_member(
    sai_object_id_t *acl_table_group_member_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateAclTableGroupMemberRequest req;
  lemming::dataplane::sai::CreateAclTableGroupMemberResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_GROUP_ID:
        req.set_acl_table_group_id(attr_list[i].value.oid);
        break;
      case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_ID:
        req.set_acl_table_id(attr_list[i].value.oid);
        break;
      case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_PRIORITY:
        req.set_priority(attr_list[i].value.u32);
        break;
    }
  }
  grpc::Status status = acl->CreateAclTableGroupMember(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *acl_table_group_member_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_ACL_TABLE_GROUP_MEMBER,
                            acl_table_group_member_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_acl_table_group_member(
    sai_object_id_t acl_table_group_member_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_ACL_TABLE_GROUP_MEMBER,
                            acl_table_group_member_id);
}

sai_status_t l_set_acl_table_group_member_attribute(
    sai_object_id_t acl_table_group_member_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_ACL_TABLE_GROUP_MEMBER,
                                   acl_table_group_member_id, attr);
}

sai_status_t l_get_acl_table_group_member_attribute(
    sai_object_id_t acl_table_group_member_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_ACL_TABLE_GROUP_MEMBER,
                                   acl_table_group_member_id, attr_count,
                                   attr_list);
}
