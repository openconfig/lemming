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
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/acl.pb.h"
#include <glog/logging.h>

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
	.set_acl_table_group_member_attribute = l_set_acl_table_group_member_attribute,
	.get_acl_table_group_member_attribute = l_get_acl_table_group_member_attribute,
	.create_acl_table_chain_group = l_create_acl_table_chain_group,
	.remove_acl_table_chain_group = l_remove_acl_table_chain_group,
	.set_acl_table_chain_group_attribute = l_set_acl_table_chain_group_attribute,
	.get_acl_table_chain_group_attribute = l_get_acl_table_chain_group_attribute,
};


lemming::dataplane::sai::CreateAclTableRequest convert_create_acl_table(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateAclTableRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
if (attr_list[i].id >= SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
  (*msg.mutable_user_defined_field_group_min())[attr_list[i].id - SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN] = attr_list[i].value.oid;
}
	

switch (attr_list[i].id) {
  
  case SAI_ACL_TABLE_ATTR_ACL_STAGE:
	msg.set_acl_stage(convert_sai_acl_stage_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_ACL_TABLE_ATTR_ACL_BIND_POINT_TYPE_LIST:
	msg.mutable_acl_bind_point_type_list()->CopyFrom(convert_list_sai_acl_bind_point_type_t_to_proto(attr_list[i].value.s32list));
	break;
  case SAI_ACL_TABLE_ATTR_SIZE:
	msg.set_size(attr_list[i].value.u32);
	break;
  case SAI_ACL_TABLE_ATTR_ACL_ACTION_TYPE_LIST:
	msg.mutable_acl_action_type_list()->CopyFrom(convert_list_sai_acl_action_type_t_to_proto(attr_list[i].value.s32list));
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6:
	msg.set_field_src_ipv6(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD3:
	msg.set_field_src_ipv6_word3(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD2:
	msg.set_field_src_ipv6_word2(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD1:
	msg.set_field_src_ipv6_word1(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD0:
	msg.set_field_src_ipv6_word0(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6:
	msg.set_field_dst_ipv6(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD3:
	msg.set_field_dst_ipv6_word3(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD2:
	msg.set_field_dst_ipv6_word2(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD1:
	msg.set_field_dst_ipv6_word1(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD0:
	msg.set_field_dst_ipv6_word0(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_IPV6:
	msg.set_field_inner_src_ipv6(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_IPV6:
	msg.set_field_inner_dst_ipv6(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_MAC:
	msg.set_field_src_mac(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_MAC:
	msg.set_field_dst_mac(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IP:
	msg.set_field_src_ip(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IP:
	msg.set_field_dst_ip(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_IP:
	msg.set_field_inner_src_ip(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_IP:
	msg.set_field_inner_dst_ip(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IN_PORTS:
	msg.set_field_in_ports(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_OUT_PORTS:
	msg.set_field_out_ports(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IN_PORT:
	msg.set_field_in_port(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_OUT_PORT:
	msg.set_field_out_port(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_PORT:
	msg.set_field_src_port(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_ID:
	msg.set_field_outer_vlan_id(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_PRI:
	msg.set_field_outer_vlan_pri(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_CFI:
	msg.set_field_outer_vlan_cfi(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_ID:
	msg.set_field_inner_vlan_id(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_PRI:
	msg.set_field_inner_vlan_pri(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_CFI:
	msg.set_field_inner_vlan_cfi(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_L4_SRC_PORT:
	msg.set_field_l4_src_port(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_L4_DST_PORT:
	msg.set_field_l4_dst_port(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_L4_SRC_PORT:
	msg.set_field_inner_l4_src_port(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_L4_DST_PORT:
	msg.set_field_inner_l4_dst_port(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ETHER_TYPE:
	msg.set_field_ether_type(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_ETHER_TYPE:
	msg.set_field_inner_ether_type(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IP_PROTOCOL:
	msg.set_field_ip_protocol(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_IP_PROTOCOL:
	msg.set_field_inner_ip_protocol(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IP_IDENTIFICATION:
	msg.set_field_ip_identification(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DSCP:
	msg.set_field_dscp(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ECN:
	msg.set_field_ecn(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TTL:
	msg.set_field_ttl(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TOS:
	msg.set_field_tos(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IP_FLAGS:
	msg.set_field_ip_flags(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TCP_FLAGS:
	msg.set_field_tcp_flags(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ACL_IP_TYPE:
	msg.set_field_acl_ip_type(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ACL_IP_FRAG:
	msg.set_field_acl_ip_frag(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IPV6_FLOW_LABEL:
	msg.set_field_ipv6_flow_label(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TC:
	msg.set_field_tc(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ICMP_TYPE:
	msg.set_field_icmp_type(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ICMP_CODE:
	msg.set_field_icmp_code(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ICMPV6_TYPE:
	msg.set_field_icmpv6_type(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ICMPV6_CODE:
	msg.set_field_icmpv6_code(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_PACKET_VLAN:
	msg.set_field_packet_vlan(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TUNNEL_VNI:
	msg.set_field_tunnel_vni(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_HAS_VLAN_TAG:
	msg.set_field_has_vlan_tag(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MACSEC_SCI:
	msg.set_field_macsec_sci(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_LABEL:
	msg.set_field_mpls_label0_label(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_TTL:
	msg.set_field_mpls_label0_ttl(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_EXP:
	msg.set_field_mpls_label0_exp(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_BOS:
	msg.set_field_mpls_label0_bos(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_LABEL:
	msg.set_field_mpls_label1_label(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_TTL:
	msg.set_field_mpls_label1_ttl(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_EXP:
	msg.set_field_mpls_label1_exp(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_BOS:
	msg.set_field_mpls_label1_bos(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_LABEL:
	msg.set_field_mpls_label2_label(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_TTL:
	msg.set_field_mpls_label2_ttl(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_EXP:
	msg.set_field_mpls_label2_exp(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_BOS:
	msg.set_field_mpls_label2_bos(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_LABEL:
	msg.set_field_mpls_label3_label(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_TTL:
	msg.set_field_mpls_label3_ttl(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_EXP:
	msg.set_field_mpls_label3_exp(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_BOS:
	msg.set_field_mpls_label3_bos(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_LABEL:
	msg.set_field_mpls_label4_label(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_TTL:
	msg.set_field_mpls_label4_ttl(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_EXP:
	msg.set_field_mpls_label4_exp(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_BOS:
	msg.set_field_mpls_label4_bos(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_FDB_DST_USER_META:
	msg.set_field_fdb_dst_user_meta(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ROUTE_DST_USER_META:
	msg.set_field_route_dst_user_meta(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_NEIGHBOR_DST_USER_META:
	msg.set_field_neighbor_dst_user_meta(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_PORT_USER_META:
	msg.set_field_port_user_meta(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_VLAN_USER_META:
	msg.set_field_vlan_user_meta(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ACL_USER_META:
	msg.set_field_acl_user_meta(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_FDB_NPU_META_DST_HIT:
	msg.set_field_fdb_npu_meta_dst_hit(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT:
	msg.set_field_neighbor_npu_meta_dst_hit(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ROUTE_NPU_META_DST_HIT:
	msg.set_field_route_npu_meta_dst_hit(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_BTH_OPCODE:
	msg.set_field_bth_opcode(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_AETH_SYNDROME:
	msg.set_field_aeth_syndrome(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ACL_RANGE_TYPE:
	msg.mutable_field_acl_range_type()->CopyFrom(convert_list_sai_acl_range_type_t_to_proto(attr_list[i].value.s32list));
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IPV6_NEXT_HEADER:
	msg.set_field_ipv6_next_header(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_GRE_KEY:
	msg.set_field_gre_key(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TAM_INT_TYPE:
	msg.set_field_tam_int_type(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_MAC:
	msg.set_field_inner_src_mac(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_MAC:
	msg.set_field_inner_dst_mac(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_VRF_ID:
	msg.set_field_vrf_id(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IPMC_NPU_META_DST_HIT:
	msg.set_field_ipmc_npu_meta_dst_hit(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TUNNEL_TERMINATED:
	msg.set_field_tunnel_terminated(attr_list[i].value.booldata);
	break;
  case SAI_ACL_TABLE_ATTR_ACL_TABLE_MATCH_TYPE:
	msg.set_acl_table_match_type(convert_sai_acl_table_match_type_t_to_proto(attr_list[i].value.s32));
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateAclEntryRequest convert_create_acl_entry(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateAclEntryRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
if (attr_list[i].id >= SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
    *(*msg.mutable_user_defined_field_group_min())[attr_list[i].id  - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN].mutable_data_u8list() = std::string(attr_list[i].value.aclfield.data.u8list.list, attr_list[i].value.aclfield.data.u8list.list + attr_list[i].value.aclfield.data.u8list.count);
    *(*msg.mutable_user_defined_field_group_min())[attr_list[i].id  - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN].mutable_mask_u8list() = std::string(attr_list[i].value.aclfield.mask.u8list.list, attr_list[i].value.aclfield.mask.u8list.list + attr_list[i].value.aclfield.mask.u8list.count);
}
	

switch (attr_list[i].id) {
  
  case SAI_ACL_ENTRY_ATTR_TABLE_ID:
	msg.set_table_id(attr_list[i].value.oid);
	break;
  case SAI_ACL_ENTRY_ATTR_PRIORITY:
	msg.set_priority(attr_list[i].value.u32);
	break;
  case SAI_ACL_ENTRY_ATTR_ADMIN_STATE:
	msg.set_admin_state(attr_list[i].value.booldata);
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6:
	*msg.mutable_field_src_ipv6() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD3:
	*msg.mutable_field_src_ipv6_word3() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD2:
	*msg.mutable_field_src_ipv6_word2() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD1:
	*msg.mutable_field_src_ipv6_word1() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD0:
	*msg.mutable_field_src_ipv6_word0() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6:
	*msg.mutable_field_dst_ipv6() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD3:
	*msg.mutable_field_dst_ipv6_word3() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD2:
	*msg.mutable_field_dst_ipv6_word2() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD1:
	*msg.mutable_field_dst_ipv6_word1() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD0:
	*msg.mutable_field_dst_ipv6_word0() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_IPV6:
	*msg.mutable_field_inner_src_ipv6() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_IPV6:
	*msg.mutable_field_inner_dst_ipv6() = convert_from_acl_field_data_ip6( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip6, attr_list[i].value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_MAC:
	*msg.mutable_field_src_mac() = convert_from_acl_field_data_mac( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.mac, attr_list[i].value.aclfield.mask.mac );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_MAC:
	*msg.mutable_field_dst_mac() = convert_from_acl_field_data_mac( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.mac, attr_list[i].value.aclfield.mask.mac );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IP:
	*msg.mutable_field_src_ip() = convert_from_acl_field_data_ip4( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip4, attr_list[i].value.aclfield.mask.ip4 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IP:
	*msg.mutable_field_dst_ip() = convert_from_acl_field_data_ip4( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip4, attr_list[i].value.aclfield.mask.ip4 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_IP:
	*msg.mutable_field_inner_src_ip() = convert_from_acl_field_data_ip4( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip4, attr_list[i].value.aclfield.mask.ip4 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_IP:
	*msg.mutable_field_inner_dst_ip() = convert_from_acl_field_data_ip4( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.ip4, attr_list[i].value.aclfield.mask.ip4 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IN_PORT:
	*msg.mutable_field_in_port() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUT_PORT:
	*msg.mutable_field_out_port() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_PORT:
	*msg.mutable_field_src_port() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_ID:
	*msg.mutable_field_outer_vlan_id() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_PRI:
	*msg.mutable_field_outer_vlan_pri() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_CFI:
	*msg.mutable_field_outer_vlan_cfi() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_ID:
	*msg.mutable_field_inner_vlan_id() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_PRI:
	*msg.mutable_field_inner_vlan_pri() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_CFI:
	*msg.mutable_field_inner_vlan_cfi() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_L4_SRC_PORT:
	*msg.mutable_field_l4_src_port() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_L4_DST_PORT:
	*msg.mutable_field_l4_dst_port() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_L4_SRC_PORT:
	*msg.mutable_field_inner_l4_src_port() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_L4_DST_PORT:
	*msg.mutable_field_inner_l4_dst_port() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ETHER_TYPE:
	*msg.mutable_field_ether_type() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_ETHER_TYPE:
	*msg.mutable_field_inner_ether_type() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IP_PROTOCOL:
	*msg.mutable_field_ip_protocol() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_IP_PROTOCOL:
	*msg.mutable_field_inner_ip_protocol() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IP_IDENTIFICATION:
	*msg.mutable_field_ip_identification() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DSCP:
	*msg.mutable_field_dscp() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ECN:
	*msg.mutable_field_ecn() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TTL:
	*msg.mutable_field_ttl() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TOS:
	*msg.mutable_field_tos() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IP_FLAGS:
	*msg.mutable_field_ip_flags() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TCP_FLAGS:
	*msg.mutable_field_tcp_flags() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ACL_IP_TYPE:
	*msg.mutable_field_acl_ip_type() = convert_from_acl_field_data_ip_type( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.s32, attr_list[i].value.aclfield.mask.s32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IPV6_FLOW_LABEL:
	*msg.mutable_field_ipv6_flow_label() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TC:
	*msg.mutable_field_tc() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMP_TYPE:
	*msg.mutable_field_icmp_type() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMP_CODE:
	*msg.mutable_field_icmp_code() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMPV6_TYPE:
	*msg.mutable_field_icmpv6_type() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMPV6_CODE:
	*msg.mutable_field_icmpv6_code() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TUNNEL_VNI:
	*msg.mutable_field_tunnel_vni() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_LABEL:
	*msg.mutable_field_mpls_label0_label() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_TTL:
	*msg.mutable_field_mpls_label0_ttl() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_EXP:
	*msg.mutable_field_mpls_label0_exp() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_LABEL:
	*msg.mutable_field_mpls_label1_label() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_TTL:
	*msg.mutable_field_mpls_label1_ttl() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_EXP:
	*msg.mutable_field_mpls_label1_exp() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_LABEL:
	*msg.mutable_field_mpls_label2_label() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_TTL:
	*msg.mutable_field_mpls_label2_ttl() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_EXP:
	*msg.mutable_field_mpls_label2_exp() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_LABEL:
	*msg.mutable_field_mpls_label3_label() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_TTL:
	*msg.mutable_field_mpls_label3_ttl() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_EXP:
	*msg.mutable_field_mpls_label3_exp() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_LABEL:
	*msg.mutable_field_mpls_label4_label() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_TTL:
	*msg.mutable_field_mpls_label4_ttl() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_EXP:
	*msg.mutable_field_mpls_label4_exp() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_FDB_DST_USER_META:
	*msg.mutable_field_fdb_dst_user_meta() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ROUTE_DST_USER_META:
	*msg.mutable_field_route_dst_user_meta() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_NEIGHBOR_DST_USER_META:
	*msg.mutable_field_neighbor_dst_user_meta() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_PORT_USER_META:
	*msg.mutable_field_port_user_meta() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_VLAN_USER_META:
	*msg.mutable_field_vlan_user_meta() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ACL_USER_META:
	*msg.mutable_field_acl_user_meta() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_BTH_OPCODE:
	*msg.mutable_field_bth_opcode() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_AETH_SYNDROME:
	*msg.mutable_field_aeth_syndrome() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IPV6_NEXT_HEADER:
	*msg.mutable_field_ipv6_next_header() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_GRE_KEY:
	*msg.mutable_field_gre_key() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u32, attr_list[i].value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_MAC:
	*msg.mutable_field_inner_src_mac() = convert_from_acl_field_data_mac( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.mac, attr_list[i].value.aclfield.mask.mac );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_MAC:
	*msg.mutable_field_inner_dst_mac() = convert_from_acl_field_data_mac( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.mac, attr_list[i].value.aclfield.mask.mac );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_VRF_ID:
	*msg.mutable_field_vrf_id() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_REDIRECT:
	*msg.mutable_action_redirect() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_PACKET_ACTION:
	*msg.mutable_action_packet_action() = convert_from_acl_action_data_action( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.s32 );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_COUNTER:
	*msg.mutable_action_counter() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_POLICER:
	*msg.mutable_action_set_policer() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_INGRESS_SAMPLEPACKET_ENABLE:
	*msg.mutable_action_ingress_samplepacket_enable() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_EGRESS_SAMPLEPACKET_ENABLE:
	*msg.mutable_action_egress_samplepacket_enable() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_USER_TRAP_ID:
	*msg.mutable_action_set_user_trap_id() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_DTEL_INT_SESSION:
	*msg.mutable_action_dtel_int_session() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_TAM_INT_OBJECT:
	*msg.mutable_action_tam_int_object() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ISOLATION_GROUP:
	*msg.mutable_action_set_isolation_group() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_MACSEC_FLOW:
	*msg.mutable_action_macsec_flow() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_LAG_HASH_ID:
	*msg.mutable_action_set_lag_hash_id() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ECMP_HASH_ID:
	*msg.mutable_action_set_ecmp_hash_id() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_VRF:
	*msg.mutable_action_set_vrf() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ARS_MONITORING:
	*msg.mutable_action_set_ars_monitoring() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ARS_OBJECT:
	*msg.mutable_action_set_ars_object() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_CHAIN_REDIRECT:
	*msg.mutable_action_chain_redirect() = convert_from_acl_action_data( attr_list[i].value.aclaction, attr_list[i].value.aclaction.parameter.oid );
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateAclCounterRequest convert_create_acl_counter(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateAclCounterRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_ACL_COUNTER_ATTR_TABLE_ID:
	msg.set_table_id(attr_list[i].value.oid);
	break;
  case SAI_ACL_COUNTER_ATTR_ENABLE_PACKET_COUNT:
	msg.set_enable_packet_count(attr_list[i].value.booldata);
	break;
  case SAI_ACL_COUNTER_ATTR_ENABLE_BYTE_COUNT:
	msg.set_enable_byte_count(attr_list[i].value.booldata);
	break;
  case SAI_ACL_COUNTER_ATTR_PACKETS:
	msg.set_packets(attr_list[i].value.u64);
	break;
  case SAI_ACL_COUNTER_ATTR_BYTES:
	msg.set_bytes(attr_list[i].value.u64);
	break;
  case SAI_ACL_COUNTER_ATTR_LABEL:
	msg.set_label(attr_list[i].value.chardata);
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateAclRangeRequest convert_create_acl_range(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateAclRangeRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_ACL_RANGE_ATTR_TYPE:
	msg.set_type(convert_sai_acl_range_type_t_to_proto(attr_list[i].value.s32));
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateAclTableGroupRequest convert_create_acl_table_group(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateAclTableGroupRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_ACL_TABLE_GROUP_ATTR_ACL_STAGE:
	msg.set_acl_stage(convert_sai_acl_stage_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_ACL_TABLE_GROUP_ATTR_ACL_BIND_POINT_TYPE_LIST:
	msg.mutable_acl_bind_point_type_list()->CopyFrom(convert_list_sai_acl_bind_point_type_t_to_proto(attr_list[i].value.s32list));
	break;
  case SAI_ACL_TABLE_GROUP_ATTR_TYPE:
	msg.set_type(convert_sai_acl_table_group_type_t_to_proto(attr_list[i].value.s32));
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateAclTableGroupMemberRequest convert_create_acl_table_group_member(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateAclTableGroupMemberRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_GROUP_ID:
	msg.set_acl_table_group_id(attr_list[i].value.oid);
	break;
  case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_ID:
	msg.set_acl_table_id(attr_list[i].value.oid);
	break;
  case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_PRIORITY:
	msg.set_priority(attr_list[i].value.u32);
	break;
  case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_CHAIN_GROUP_ID:
	msg.set_acl_table_chain_group_id(attr_list[i].value.oid);
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateAclTableChainGroupRequest convert_create_acl_table_chain_group(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateAclTableChainGroupRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_ACL_TABLE_CHAIN_GROUP_ATTR_TYPE:
	msg.set_type(convert_sai_acl_table_chain_group_type_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_ACL_TABLE_CHAIN_GROUP_ATTR_STAGE:
	msg.set_stage(convert_sai_acl_table_chain_group_stage_t_to_proto(attr_list[i].value.s32));
	break;
}

}
return msg;
}

sai_status_t l_create_acl_table(sai_object_id_t *acl_table_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateAclTableRequest req = convert_create_acl_table(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateAclTableResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = acl->CreateAclTable(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (acl_table_id) {
	*acl_table_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_acl_table(sai_object_id_t acl_table_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveAclTableRequest req;
	lemming::dataplane::sai::RemoveAclTableResponse resp;
	grpc::ClientContext context;
	req.set_oid(acl_table_id); 
	
	grpc::Status status = acl->RemoveAclTable(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_acl_table_attribute(sai_object_id_t acl_table_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_acl_table_attribute(sai_object_id_t acl_table_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetAclTableAttributeRequest req;
	lemming::dataplane::sai::GetAclTableAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(acl_table_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_acl_table_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = acl->GetAclTableAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		
if (attr_list[i].id >= SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
  attr_list[i].value.oid = resp.attr().user_defined_field_group_min().at(attr_list[i].id - SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN);
}
		

switch (attr_list[i].id) {
  
  case SAI_ACL_TABLE_ATTR_ACL_STAGE:
	 attr_list[i].value.s32 =  convert_sai_acl_stage_t_to_sai(resp.attr().acl_stage());
	break;
  case SAI_ACL_TABLE_ATTR_ACL_BIND_POINT_TYPE_LIST:
	convert_list_sai_acl_bind_point_type_t_to_sai(attr_list[i].value.s32list.list, resp.attr().acl_bind_point_type_list(), &attr_list[i].value.s32list.count);
	break;
  case SAI_ACL_TABLE_ATTR_SIZE:
	 attr_list[i].value.u32 =   resp.attr().size();
	break;
  case SAI_ACL_TABLE_ATTR_ACL_ACTION_TYPE_LIST:
	convert_list_sai_acl_action_type_t_to_sai(attr_list[i].value.s32list.list, resp.attr().acl_action_type_list(), &attr_list[i].value.s32list.count);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6:
	 attr_list[i].value.booldata =   resp.attr().field_src_ipv6();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD3:
	 attr_list[i].value.booldata =   resp.attr().field_src_ipv6_word3();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD2:
	 attr_list[i].value.booldata =   resp.attr().field_src_ipv6_word2();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD1:
	 attr_list[i].value.booldata =   resp.attr().field_src_ipv6_word1();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD0:
	 attr_list[i].value.booldata =   resp.attr().field_src_ipv6_word0();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6:
	 attr_list[i].value.booldata =   resp.attr().field_dst_ipv6();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD3:
	 attr_list[i].value.booldata =   resp.attr().field_dst_ipv6_word3();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD2:
	 attr_list[i].value.booldata =   resp.attr().field_dst_ipv6_word2();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD1:
	 attr_list[i].value.booldata =   resp.attr().field_dst_ipv6_word1();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD0:
	 attr_list[i].value.booldata =   resp.attr().field_dst_ipv6_word0();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_IPV6:
	 attr_list[i].value.booldata =   resp.attr().field_inner_src_ipv6();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_IPV6:
	 attr_list[i].value.booldata =   resp.attr().field_inner_dst_ipv6();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_MAC:
	 attr_list[i].value.booldata =   resp.attr().field_src_mac();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_MAC:
	 attr_list[i].value.booldata =   resp.attr().field_dst_mac();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_IP:
	 attr_list[i].value.booldata =   resp.attr().field_src_ip();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DST_IP:
	 attr_list[i].value.booldata =   resp.attr().field_dst_ip();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_IP:
	 attr_list[i].value.booldata =   resp.attr().field_inner_src_ip();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_IP:
	 attr_list[i].value.booldata =   resp.attr().field_inner_dst_ip();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IN_PORTS:
	 attr_list[i].value.booldata =   resp.attr().field_in_ports();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_OUT_PORTS:
	 attr_list[i].value.booldata =   resp.attr().field_out_ports();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IN_PORT:
	 attr_list[i].value.booldata =   resp.attr().field_in_port();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_OUT_PORT:
	 attr_list[i].value.booldata =   resp.attr().field_out_port();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_SRC_PORT:
	 attr_list[i].value.booldata =   resp.attr().field_src_port();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_ID:
	 attr_list[i].value.booldata =   resp.attr().field_outer_vlan_id();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_PRI:
	 attr_list[i].value.booldata =   resp.attr().field_outer_vlan_pri();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_CFI:
	 attr_list[i].value.booldata =   resp.attr().field_outer_vlan_cfi();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_ID:
	 attr_list[i].value.booldata =   resp.attr().field_inner_vlan_id();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_PRI:
	 attr_list[i].value.booldata =   resp.attr().field_inner_vlan_pri();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_CFI:
	 attr_list[i].value.booldata =   resp.attr().field_inner_vlan_cfi();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_L4_SRC_PORT:
	 attr_list[i].value.booldata =   resp.attr().field_l4_src_port();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_L4_DST_PORT:
	 attr_list[i].value.booldata =   resp.attr().field_l4_dst_port();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_L4_SRC_PORT:
	 attr_list[i].value.booldata =   resp.attr().field_inner_l4_src_port();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_L4_DST_PORT:
	 attr_list[i].value.booldata =   resp.attr().field_inner_l4_dst_port();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ETHER_TYPE:
	 attr_list[i].value.booldata =   resp.attr().field_ether_type();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_ETHER_TYPE:
	 attr_list[i].value.booldata =   resp.attr().field_inner_ether_type();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IP_PROTOCOL:
	 attr_list[i].value.booldata =   resp.attr().field_ip_protocol();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_IP_PROTOCOL:
	 attr_list[i].value.booldata =   resp.attr().field_inner_ip_protocol();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IP_IDENTIFICATION:
	 attr_list[i].value.booldata =   resp.attr().field_ip_identification();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_DSCP:
	 attr_list[i].value.booldata =   resp.attr().field_dscp();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ECN:
	 attr_list[i].value.booldata =   resp.attr().field_ecn();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TTL:
	 attr_list[i].value.booldata =   resp.attr().field_ttl();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TOS:
	 attr_list[i].value.booldata =   resp.attr().field_tos();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IP_FLAGS:
	 attr_list[i].value.booldata =   resp.attr().field_ip_flags();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TCP_FLAGS:
	 attr_list[i].value.booldata =   resp.attr().field_tcp_flags();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ACL_IP_TYPE:
	 attr_list[i].value.booldata =   resp.attr().field_acl_ip_type();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ACL_IP_FRAG:
	 attr_list[i].value.booldata =   resp.attr().field_acl_ip_frag();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IPV6_FLOW_LABEL:
	 attr_list[i].value.booldata =   resp.attr().field_ipv6_flow_label();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TC:
	 attr_list[i].value.booldata =   resp.attr().field_tc();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ICMP_TYPE:
	 attr_list[i].value.booldata =   resp.attr().field_icmp_type();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ICMP_CODE:
	 attr_list[i].value.booldata =   resp.attr().field_icmp_code();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ICMPV6_TYPE:
	 attr_list[i].value.booldata =   resp.attr().field_icmpv6_type();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ICMPV6_CODE:
	 attr_list[i].value.booldata =   resp.attr().field_icmpv6_code();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_PACKET_VLAN:
	 attr_list[i].value.booldata =   resp.attr().field_packet_vlan();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TUNNEL_VNI:
	 attr_list[i].value.booldata =   resp.attr().field_tunnel_vni();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_HAS_VLAN_TAG:
	 attr_list[i].value.booldata =   resp.attr().field_has_vlan_tag();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MACSEC_SCI:
	 attr_list[i].value.booldata =   resp.attr().field_macsec_sci();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_LABEL:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label0_label();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_TTL:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label0_ttl();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_EXP:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label0_exp();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_BOS:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label0_bos();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_LABEL:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label1_label();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_TTL:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label1_ttl();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_EXP:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label1_exp();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_BOS:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label1_bos();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_LABEL:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label2_label();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_TTL:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label2_ttl();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_EXP:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label2_exp();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_BOS:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label2_bos();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_LABEL:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label3_label();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_TTL:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label3_ttl();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_EXP:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label3_exp();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_BOS:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label3_bos();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_LABEL:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label4_label();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_TTL:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label4_ttl();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_EXP:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label4_exp();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_BOS:
	 attr_list[i].value.booldata =   resp.attr().field_mpls_label4_bos();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_FDB_DST_USER_META:
	 attr_list[i].value.booldata =   resp.attr().field_fdb_dst_user_meta();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ROUTE_DST_USER_META:
	 attr_list[i].value.booldata =   resp.attr().field_route_dst_user_meta();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_NEIGHBOR_DST_USER_META:
	 attr_list[i].value.booldata =   resp.attr().field_neighbor_dst_user_meta();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_PORT_USER_META:
	 attr_list[i].value.booldata =   resp.attr().field_port_user_meta();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_VLAN_USER_META:
	 attr_list[i].value.booldata =   resp.attr().field_vlan_user_meta();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ACL_USER_META:
	 attr_list[i].value.booldata =   resp.attr().field_acl_user_meta();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_FDB_NPU_META_DST_HIT:
	 attr_list[i].value.booldata =   resp.attr().field_fdb_npu_meta_dst_hit();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT:
	 attr_list[i].value.booldata =   resp.attr().field_neighbor_npu_meta_dst_hit();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ROUTE_NPU_META_DST_HIT:
	 attr_list[i].value.booldata =   resp.attr().field_route_npu_meta_dst_hit();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_BTH_OPCODE:
	 attr_list[i].value.booldata =   resp.attr().field_bth_opcode();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_AETH_SYNDROME:
	 attr_list[i].value.booldata =   resp.attr().field_aeth_syndrome();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_ACL_RANGE_TYPE:
	convert_list_sai_acl_range_type_t_to_sai(attr_list[i].value.s32list.list, resp.attr().field_acl_range_type(), &attr_list[i].value.s32list.count);
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IPV6_NEXT_HEADER:
	 attr_list[i].value.booldata =   resp.attr().field_ipv6_next_header();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_GRE_KEY:
	 attr_list[i].value.booldata =   resp.attr().field_gre_key();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TAM_INT_TYPE:
	 attr_list[i].value.booldata =   resp.attr().field_tam_int_type();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_MAC:
	 attr_list[i].value.booldata =   resp.attr().field_inner_src_mac();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_MAC:
	 attr_list[i].value.booldata =   resp.attr().field_inner_dst_mac();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_VRF_ID:
	 attr_list[i].value.booldata =   resp.attr().field_vrf_id();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_IPMC_NPU_META_DST_HIT:
	 attr_list[i].value.booldata =   resp.attr().field_ipmc_npu_meta_dst_hit();
	break;
  case SAI_ACL_TABLE_ATTR_FIELD_TUNNEL_TERMINATED:
	 attr_list[i].value.booldata =   resp.attr().field_tunnel_terminated();
	break;
  case SAI_ACL_TABLE_ATTR_ENTRY_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().entry_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_ACL_TABLE_ATTR_AVAILABLE_ACL_ENTRY:
	 attr_list[i].value.u32 =   resp.attr().available_acl_entry();
	break;
  case SAI_ACL_TABLE_ATTR_AVAILABLE_ACL_COUNTER:
	 attr_list[i].value.u32 =   resp.attr().available_acl_counter();
	break;
  case SAI_ACL_TABLE_ATTR_ACL_TABLE_MATCH_TYPE:
	 attr_list[i].value.s32 =  convert_sai_acl_table_match_type_t_to_sai(resp.attr().acl_table_match_type());
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_acl_entry(sai_object_id_t *acl_entry_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateAclEntryRequest req = convert_create_acl_entry(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateAclEntryResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = acl->CreateAclEntry(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (acl_entry_id) {
	*acl_entry_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_acl_entry(sai_object_id_t acl_entry_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveAclEntryRequest req;
	lemming::dataplane::sai::RemoveAclEntryResponse resp;
	grpc::ClientContext context;
	req.set_oid(acl_entry_id); 
	
	grpc::Status status = acl->RemoveAclEntry(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_acl_entry_attribute(sai_object_id_t acl_entry_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetAclEntryAttributeRequest req;
	lemming::dataplane::sai::SetAclEntryAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(acl_entry_id); 
	
	
if (attr->id >= SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr->id < SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
*(*req.mutable_user_defined_field_group_min())[attr->id  - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN].mutable_data_u8list() = std::string(attr->value.aclfield.data.u8list.list, attr->value.aclfield.data.u8list.list + attr->value.aclfield.data.u8list.count);
*(*req.mutable_user_defined_field_group_min())[attr->id  - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN].mutable_mask_u8list() = std::string(attr->value.aclfield.mask.u8list.list, attr->value.aclfield.mask.u8list.list + attr->value.aclfield.mask.u8list.count);
}
	

switch (attr->id) {
  
  case SAI_ACL_ENTRY_ATTR_PRIORITY:
	req.set_priority(attr->value.u32);
	break;
  case SAI_ACL_ENTRY_ATTR_ADMIN_STATE:
	req.set_admin_state(attr->value.booldata);
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6:
	*req.mutable_field_src_ipv6() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD3:
	*req.mutable_field_src_ipv6_word3() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD2:
	*req.mutable_field_src_ipv6_word2() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD1:
	*req.mutable_field_src_ipv6_word1() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD0:
	*req.mutable_field_src_ipv6_word0() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6:
	*req.mutable_field_dst_ipv6() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD3:
	*req.mutable_field_dst_ipv6_word3() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD2:
	*req.mutable_field_dst_ipv6_word2() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD1:
	*req.mutable_field_dst_ipv6_word1() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD0:
	*req.mutable_field_dst_ipv6_word0() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_IPV6:
	*req.mutable_field_inner_src_ipv6() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_IPV6:
	*req.mutable_field_inner_dst_ipv6() = convert_from_acl_field_data_ip6( attr->value.aclfield, attr->value.aclfield.data.ip6, attr->value.aclfield.mask.ip6 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_MAC:
	*req.mutable_field_src_mac() = convert_from_acl_field_data_mac( attr->value.aclfield, attr->value.aclfield.data.mac, attr->value.aclfield.mask.mac );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_MAC:
	*req.mutable_field_dst_mac() = convert_from_acl_field_data_mac( attr->value.aclfield, attr->value.aclfield.data.mac, attr->value.aclfield.mask.mac );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IP:
	*req.mutable_field_src_ip() = convert_from_acl_field_data_ip4( attr->value.aclfield, attr->value.aclfield.data.ip4, attr->value.aclfield.mask.ip4 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IP:
	*req.mutable_field_dst_ip() = convert_from_acl_field_data_ip4( attr->value.aclfield, attr->value.aclfield.data.ip4, attr->value.aclfield.mask.ip4 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_IP:
	*req.mutable_field_inner_src_ip() = convert_from_acl_field_data_ip4( attr->value.aclfield, attr->value.aclfield.data.ip4, attr->value.aclfield.mask.ip4 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_IP:
	*req.mutable_field_inner_dst_ip() = convert_from_acl_field_data_ip4( attr->value.aclfield, attr->value.aclfield.data.ip4, attr->value.aclfield.mask.ip4 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IN_PORT:
	*req.mutable_field_in_port() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUT_PORT:
	*req.mutable_field_out_port() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_PORT:
	*req.mutable_field_src_port() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_ID:
	*req.mutable_field_outer_vlan_id() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u16, attr->value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_PRI:
	*req.mutable_field_outer_vlan_pri() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_CFI:
	*req.mutable_field_outer_vlan_cfi() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_ID:
	*req.mutable_field_inner_vlan_id() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u16, attr->value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_PRI:
	*req.mutable_field_inner_vlan_pri() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_CFI:
	*req.mutable_field_inner_vlan_cfi() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_L4_SRC_PORT:
	*req.mutable_field_l4_src_port() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u16, attr->value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_L4_DST_PORT:
	*req.mutable_field_l4_dst_port() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u16, attr->value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_L4_SRC_PORT:
	*req.mutable_field_inner_l4_src_port() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u16, attr->value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_L4_DST_PORT:
	*req.mutable_field_inner_l4_dst_port() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u16, attr->value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ETHER_TYPE:
	*req.mutable_field_ether_type() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u16, attr->value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_ETHER_TYPE:
	*req.mutable_field_inner_ether_type() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u16, attr->value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IP_PROTOCOL:
	*req.mutable_field_ip_protocol() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_IP_PROTOCOL:
	*req.mutable_field_inner_ip_protocol() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IP_IDENTIFICATION:
	*req.mutable_field_ip_identification() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u16, attr->value.aclfield.mask.u16 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DSCP:
	*req.mutable_field_dscp() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ECN:
	*req.mutable_field_ecn() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TTL:
	*req.mutable_field_ttl() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TOS:
	*req.mutable_field_tos() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IP_FLAGS:
	*req.mutable_field_ip_flags() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TCP_FLAGS:
	*req.mutable_field_tcp_flags() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ACL_IP_TYPE:
	*req.mutable_field_acl_ip_type() = convert_from_acl_field_data_ip_type( attr->value.aclfield, attr->value.aclfield.data.s32, attr->value.aclfield.mask.s32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IPV6_FLOW_LABEL:
	*req.mutable_field_ipv6_flow_label() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TC:
	*req.mutable_field_tc() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMP_TYPE:
	*req.mutable_field_icmp_type() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMP_CODE:
	*req.mutable_field_icmp_code() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMPV6_TYPE:
	*req.mutable_field_icmpv6_type() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMPV6_CODE:
	*req.mutable_field_icmpv6_code() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TUNNEL_VNI:
	*req.mutable_field_tunnel_vni() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_LABEL:
	*req.mutable_field_mpls_label0_label() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_TTL:
	*req.mutable_field_mpls_label0_ttl() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_EXP:
	*req.mutable_field_mpls_label0_exp() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_LABEL:
	*req.mutable_field_mpls_label1_label() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_TTL:
	*req.mutable_field_mpls_label1_ttl() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_EXP:
	*req.mutable_field_mpls_label1_exp() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_LABEL:
	*req.mutable_field_mpls_label2_label() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_TTL:
	*req.mutable_field_mpls_label2_ttl() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_EXP:
	*req.mutable_field_mpls_label2_exp() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_LABEL:
	*req.mutable_field_mpls_label3_label() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_TTL:
	*req.mutable_field_mpls_label3_ttl() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_EXP:
	*req.mutable_field_mpls_label3_exp() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_LABEL:
	*req.mutable_field_mpls_label4_label() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_TTL:
	*req.mutable_field_mpls_label4_ttl() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_EXP:
	*req.mutable_field_mpls_label4_exp() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_FDB_DST_USER_META:
	*req.mutable_field_fdb_dst_user_meta() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ROUTE_DST_USER_META:
	*req.mutable_field_route_dst_user_meta() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_NEIGHBOR_DST_USER_META:
	*req.mutable_field_neighbor_dst_user_meta() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_PORT_USER_META:
	*req.mutable_field_port_user_meta() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_VLAN_USER_META:
	*req.mutable_field_vlan_user_meta() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ACL_USER_META:
	*req.mutable_field_acl_user_meta() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_BTH_OPCODE:
	*req.mutable_field_bth_opcode() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_AETH_SYNDROME:
	*req.mutable_field_aeth_syndrome() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IPV6_NEXT_HEADER:
	*req.mutable_field_ipv6_next_header() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u8, attr->value.aclfield.mask.u8 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_GRE_KEY:
	*req.mutable_field_gre_key() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.u32, attr->value.aclfield.mask.u32 );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_MAC:
	*req.mutable_field_inner_src_mac() = convert_from_acl_field_data_mac( attr->value.aclfield, attr->value.aclfield.data.mac, attr->value.aclfield.mask.mac );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_MAC:
	*req.mutable_field_inner_dst_mac() = convert_from_acl_field_data_mac( attr->value.aclfield, attr->value.aclfield.data.mac, attr->value.aclfield.mask.mac );
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_VRF_ID:
	*req.mutable_field_vrf_id() = convert_from_acl_field_data( attr->value.aclfield, attr->value.aclfield.data.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_REDIRECT:
	*req.mutable_action_redirect() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_PACKET_ACTION:
	*req.mutable_action_packet_action() = convert_from_acl_action_data_action( attr->value.aclaction, attr->value.aclaction.parameter.s32 );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_COUNTER:
	*req.mutable_action_counter() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_POLICER:
	*req.mutable_action_set_policer() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_INGRESS_SAMPLEPACKET_ENABLE:
	*req.mutable_action_ingress_samplepacket_enable() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_EGRESS_SAMPLEPACKET_ENABLE:
	*req.mutable_action_egress_samplepacket_enable() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_USER_TRAP_ID:
	*req.mutable_action_set_user_trap_id() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_DTEL_INT_SESSION:
	*req.mutable_action_dtel_int_session() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_TAM_INT_OBJECT:
	*req.mutable_action_tam_int_object() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ISOLATION_GROUP:
	*req.mutable_action_set_isolation_group() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_MACSEC_FLOW:
	*req.mutable_action_macsec_flow() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_LAG_HASH_ID:
	*req.mutable_action_set_lag_hash_id() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ECMP_HASH_ID:
	*req.mutable_action_set_ecmp_hash_id() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_VRF:
	*req.mutable_action_set_vrf() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ARS_MONITORING:
	*req.mutable_action_set_ars_monitoring() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ARS_OBJECT:
	*req.mutable_action_set_ars_object() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_CHAIN_REDIRECT:
	*req.mutable_action_chain_redirect() = convert_from_acl_action_data( attr->value.aclaction, attr->value.aclaction.parameter.oid );
	break;
}

	grpc::Status status = acl->SetAclEntryAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_acl_entry_attribute(sai_object_id_t acl_entry_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetAclEntryAttributeRequest req;
	lemming::dataplane::sai::GetAclEntryAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(acl_entry_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_acl_entry_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = acl->GetAclEntryAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		
if (attr_list[i].id >= SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN && attr_list[i].id < SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MAX) {
	auto acl_attr = resp.attr().user_defined_field_group_min().at(attr_list[i].id - SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN);
	memcpy(attr_list[i].value.aclfield.data.u8list.list, acl_attr.data_u8list().data(), acl_attr.data_u8list().size());
	memcpy(attr_list[i].value.aclfield.mask.u8list.list, acl_attr.mask_u8list().data(), acl_attr.mask_u8list().size());
    attr_list[i].value.aclfield.data.u8list.count = acl_attr.data_u8list().size();
	attr_list[i].value.aclfield.mask.u8list.count = acl_attr.mask_u8list().size();
}
		

switch (attr_list[i].id) {
  
  case SAI_ACL_ENTRY_ATTR_TABLE_ID:
	 attr_list[i].value.oid =   resp.attr().table_id();
	break;
  case SAI_ACL_ENTRY_ATTR_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().priority();
	break;
  case SAI_ACL_ENTRY_ATTR_ADMIN_STATE:
	 attr_list[i].value.booldata =   resp.attr().admin_state();
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_src_ipv6(), resp.attr().field_src_ipv6().data_ip(), resp.attr().field_src_ipv6().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD3:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_src_ipv6_word3(), resp.attr().field_src_ipv6_word3().data_ip(), resp.attr().field_src_ipv6_word3().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD2:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_src_ipv6_word2(), resp.attr().field_src_ipv6_word2().data_ip(), resp.attr().field_src_ipv6_word2().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD1:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_src_ipv6_word1(), resp.attr().field_src_ipv6_word1().data_ip(), resp.attr().field_src_ipv6_word1().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD0:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_src_ipv6_word0(), resp.attr().field_src_ipv6_word0().data_ip(), resp.attr().field_src_ipv6_word0().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_dst_ipv6(), resp.attr().field_dst_ipv6().data_ip(), resp.attr().field_dst_ipv6().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD3:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_dst_ipv6_word3(), resp.attr().field_dst_ipv6_word3().data_ip(), resp.attr().field_dst_ipv6_word3().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD2:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_dst_ipv6_word2(), resp.attr().field_dst_ipv6_word2().data_ip(), resp.attr().field_dst_ipv6_word2().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD1:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_dst_ipv6_word1(), resp.attr().field_dst_ipv6_word1().data_ip(), resp.attr().field_dst_ipv6_word1().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD0:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_dst_ipv6_word0(), resp.attr().field_dst_ipv6_word0().data_ip(), resp.attr().field_dst_ipv6_word0().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_IPV6:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_inner_src_ipv6(), resp.attr().field_inner_src_ipv6().data_ip(), resp.attr().field_inner_src_ipv6().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_IPV6:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip6(resp.attr().field_inner_dst_ipv6(), resp.attr().field_inner_dst_ipv6().data_ip(), resp.attr().field_inner_dst_ipv6().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_MAC:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_mac(resp.attr().field_src_mac(), resp.attr().field_src_mac().data_mac(), resp.attr().field_src_mac().mask_mac());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_MAC:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_mac(resp.attr().field_dst_mac(), resp.attr().field_dst_mac().data_mac(), resp.attr().field_dst_mac().mask_mac());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IP:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data(resp.attr().field_src_ip(), resp.attr().field_src_ip().data_ip(), resp.attr().field_src_ip().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DST_IP:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data(resp.attr().field_dst_ip(), resp.attr().field_dst_ip().data_ip(), resp.attr().field_dst_ip().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_IP:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data(resp.attr().field_inner_src_ip(), resp.attr().field_inner_src_ip().data_ip(), resp.attr().field_inner_src_ip().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_IP:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data(resp.attr().field_inner_dst_ip(), resp.attr().field_inner_dst_ip().data_ip(), resp.attr().field_inner_dst_ip().mask_ip());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IN_PORT:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data(resp.attr().field_in_port(), resp.attr().field_in_port().data_oid());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUT_PORT:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data(resp.attr().field_out_port(), resp.attr().field_out_port().data_oid());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_SRC_PORT:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data(resp.attr().field_src_port(), resp.attr().field_src_port().data_oid());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_ID:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().field_outer_vlan_id(), resp.attr().field_outer_vlan_id().data_uint(), resp.attr().field_outer_vlan_id().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_PRI:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_outer_vlan_pri(), resp.attr().field_outer_vlan_pri().data_uint(), resp.attr().field_outer_vlan_pri().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_CFI:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_outer_vlan_cfi(), resp.attr().field_outer_vlan_cfi().data_uint(), resp.attr().field_outer_vlan_cfi().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_ID:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().field_inner_vlan_id(), resp.attr().field_inner_vlan_id().data_uint(), resp.attr().field_inner_vlan_id().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_PRI:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_inner_vlan_pri(), resp.attr().field_inner_vlan_pri().data_uint(), resp.attr().field_inner_vlan_pri().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_CFI:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_inner_vlan_cfi(), resp.attr().field_inner_vlan_cfi().data_uint(), resp.attr().field_inner_vlan_cfi().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_L4_SRC_PORT:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().field_l4_src_port(), resp.attr().field_l4_src_port().data_uint(), resp.attr().field_l4_src_port().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_L4_DST_PORT:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().field_l4_dst_port(), resp.attr().field_l4_dst_port().data_uint(), resp.attr().field_l4_dst_port().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_L4_SRC_PORT:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().field_inner_l4_src_port(), resp.attr().field_inner_l4_src_port().data_uint(), resp.attr().field_inner_l4_src_port().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_L4_DST_PORT:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().field_inner_l4_dst_port(), resp.attr().field_inner_l4_dst_port().data_uint(), resp.attr().field_inner_l4_dst_port().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ETHER_TYPE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().field_ether_type(), resp.attr().field_ether_type().data_uint(), resp.attr().field_ether_type().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_ETHER_TYPE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().field_inner_ether_type(), resp.attr().field_inner_ether_type().data_uint(), resp.attr().field_inner_ether_type().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IP_PROTOCOL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_ip_protocol(), resp.attr().field_ip_protocol().data_uint(), resp.attr().field_ip_protocol().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_IP_PROTOCOL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_inner_ip_protocol(), resp.attr().field_inner_ip_protocol().data_uint(), resp.attr().field_inner_ip_protocol().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IP_IDENTIFICATION:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().field_ip_identification(), resp.attr().field_ip_identification().data_uint(), resp.attr().field_ip_identification().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_DSCP:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_dscp(), resp.attr().field_dscp().data_uint(), resp.attr().field_dscp().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ECN:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_ecn(), resp.attr().field_ecn().data_uint(), resp.attr().field_ecn().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TTL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_ttl(), resp.attr().field_ttl().data_uint(), resp.attr().field_ttl().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TOS:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_tos(), resp.attr().field_tos().data_uint(), resp.attr().field_tos().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IP_FLAGS:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_ip_flags(), resp.attr().field_ip_flags().data_uint(), resp.attr().field_ip_flags().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TCP_FLAGS:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_tcp_flags(), resp.attr().field_tcp_flags().data_uint(), resp.attr().field_tcp_flags().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ACL_IP_TYPE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_ip_type(resp.attr().field_acl_ip_type(), resp.attr().field_acl_ip_type().data_ip_type());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IPV6_FLOW_LABEL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_ipv6_flow_label(), resp.attr().field_ipv6_flow_label().data_uint(), resp.attr().field_ipv6_flow_label().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TC:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_tc(), resp.attr().field_tc().data_uint(), resp.attr().field_tc().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMP_TYPE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_icmp_type(), resp.attr().field_icmp_type().data_uint(), resp.attr().field_icmp_type().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMP_CODE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_icmp_code(), resp.attr().field_icmp_code().data_uint(), resp.attr().field_icmp_code().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMPV6_TYPE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_icmpv6_type(), resp.attr().field_icmpv6_type().data_uint(), resp.attr().field_icmpv6_type().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ICMPV6_CODE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_icmpv6_code(), resp.attr().field_icmpv6_code().data_uint(), resp.attr().field_icmpv6_code().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_TUNNEL_VNI:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_tunnel_vni(), resp.attr().field_tunnel_vni().data_uint(), resp.attr().field_tunnel_vni().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_LABEL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_mpls_label0_label(), resp.attr().field_mpls_label0_label().data_uint(), resp.attr().field_mpls_label0_label().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_TTL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_mpls_label0_ttl(), resp.attr().field_mpls_label0_ttl().data_uint(), resp.attr().field_mpls_label0_ttl().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_EXP:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_mpls_label0_exp(), resp.attr().field_mpls_label0_exp().data_uint(), resp.attr().field_mpls_label0_exp().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_LABEL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_mpls_label1_label(), resp.attr().field_mpls_label1_label().data_uint(), resp.attr().field_mpls_label1_label().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_TTL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_mpls_label1_ttl(), resp.attr().field_mpls_label1_ttl().data_uint(), resp.attr().field_mpls_label1_ttl().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_EXP:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_mpls_label1_exp(), resp.attr().field_mpls_label1_exp().data_uint(), resp.attr().field_mpls_label1_exp().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_LABEL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_mpls_label2_label(), resp.attr().field_mpls_label2_label().data_uint(), resp.attr().field_mpls_label2_label().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_TTL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_mpls_label2_ttl(), resp.attr().field_mpls_label2_ttl().data_uint(), resp.attr().field_mpls_label2_ttl().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_EXP:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_mpls_label2_exp(), resp.attr().field_mpls_label2_exp().data_uint(), resp.attr().field_mpls_label2_exp().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_LABEL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_mpls_label3_label(), resp.attr().field_mpls_label3_label().data_uint(), resp.attr().field_mpls_label3_label().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_TTL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_mpls_label3_ttl(), resp.attr().field_mpls_label3_ttl().data_uint(), resp.attr().field_mpls_label3_ttl().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_EXP:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_mpls_label3_exp(), resp.attr().field_mpls_label3_exp().data_uint(), resp.attr().field_mpls_label3_exp().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_LABEL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_mpls_label4_label(), resp.attr().field_mpls_label4_label().data_uint(), resp.attr().field_mpls_label4_label().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_TTL:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_mpls_label4_ttl(), resp.attr().field_mpls_label4_ttl().data_uint(), resp.attr().field_mpls_label4_ttl().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_EXP:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_mpls_label4_exp(), resp.attr().field_mpls_label4_exp().data_uint(), resp.attr().field_mpls_label4_exp().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_FDB_DST_USER_META:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_fdb_dst_user_meta(), resp.attr().field_fdb_dst_user_meta().data_uint(), resp.attr().field_fdb_dst_user_meta().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ROUTE_DST_USER_META:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_route_dst_user_meta(), resp.attr().field_route_dst_user_meta().data_uint(), resp.attr().field_route_dst_user_meta().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_NEIGHBOR_DST_USER_META:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_neighbor_dst_user_meta(), resp.attr().field_neighbor_dst_user_meta().data_uint(), resp.attr().field_neighbor_dst_user_meta().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_PORT_USER_META:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_port_user_meta(), resp.attr().field_port_user_meta().data_uint(), resp.attr().field_port_user_meta().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_VLAN_USER_META:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_vlan_user_meta(), resp.attr().field_vlan_user_meta().data_uint(), resp.attr().field_vlan_user_meta().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_ACL_USER_META:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_acl_user_meta(), resp.attr().field_acl_user_meta().data_uint(), resp.attr().field_acl_user_meta().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_BTH_OPCODE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_bth_opcode(), resp.attr().field_bth_opcode().data_uint(), resp.attr().field_bth_opcode().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_AETH_SYNDROME:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_aeth_syndrome(), resp.attr().field_aeth_syndrome().data_uint(), resp.attr().field_aeth_syndrome().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_IPV6_NEXT_HEADER:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().field_ipv6_next_header(), resp.attr().field_ipv6_next_header().data_uint(), resp.attr().field_ipv6_next_header().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_GRE_KEY:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u32(resp.attr().field_gre_key(), resp.attr().field_gre_key().data_uint(), resp.attr().field_gre_key().mask_uint());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_MAC:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_mac(resp.attr().field_inner_src_mac(), resp.attr().field_inner_src_mac().data_mac(), resp.attr().field_inner_src_mac().mask_mac());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_MAC:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_mac(resp.attr().field_inner_dst_mac(), resp.attr().field_inner_dst_mac().data_mac(), resp.attr().field_inner_dst_mac().mask_mac());
	break;
  case SAI_ACL_ENTRY_ATTR_FIELD_VRF_ID:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data(resp.attr().field_vrf_id(), resp.attr().field_vrf_id().data_oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_REDIRECT:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_redirect(), resp.attr().action_redirect().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_PACKET_ACTION:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data_action(resp.attr().action_packet_action(), resp.attr().action_packet_action().packet_action());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_COUNTER:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_counter(), resp.attr().action_counter().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_POLICER:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_set_policer(), resp.attr().action_set_policer().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_INGRESS_SAMPLEPACKET_ENABLE:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_ingress_samplepacket_enable(), resp.attr().action_ingress_samplepacket_enable().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_EGRESS_SAMPLEPACKET_ENABLE:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_egress_samplepacket_enable(), resp.attr().action_egress_samplepacket_enable().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_USER_TRAP_ID:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_set_user_trap_id(), resp.attr().action_set_user_trap_id().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_DTEL_INT_SESSION:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_dtel_int_session(), resp.attr().action_dtel_int_session().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_TAM_INT_OBJECT:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_tam_int_object(), resp.attr().action_tam_int_object().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ISOLATION_GROUP:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_set_isolation_group(), resp.attr().action_set_isolation_group().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_MACSEC_FLOW:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_macsec_flow(), resp.attr().action_macsec_flow().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_LAG_HASH_ID:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_set_lag_hash_id(), resp.attr().action_set_lag_hash_id().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ECMP_HASH_ID:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_set_ecmp_hash_id(), resp.attr().action_set_ecmp_hash_id().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_VRF:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_set_vrf(), resp.attr().action_set_vrf().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ARS_MONITORING:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_set_ars_monitoring(), resp.attr().action_set_ars_monitoring().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_SET_ARS_OBJECT:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_set_ars_object(), resp.attr().action_set_ars_object().oid());
	break;
  case SAI_ACL_ENTRY_ATTR_ACTION_CHAIN_REDIRECT:
	 attr_list[i].value.aclaction =  convert_to_acl_action_data(resp.attr().action_chain_redirect(), resp.attr().action_chain_redirect().oid());
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_acl_counter(sai_object_id_t *acl_counter_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateAclCounterRequest req = convert_create_acl_counter(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateAclCounterResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = acl->CreateAclCounter(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (acl_counter_id) {
	*acl_counter_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_acl_counter(sai_object_id_t acl_counter_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveAclCounterRequest req;
	lemming::dataplane::sai::RemoveAclCounterResponse resp;
	grpc::ClientContext context;
	req.set_oid(acl_counter_id); 
	
	grpc::Status status = acl->RemoveAclCounter(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_acl_counter_attribute(sai_object_id_t acl_counter_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetAclCounterAttributeRequest req;
	lemming::dataplane::sai::SetAclCounterAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(acl_counter_id); 
	
	
	

switch (attr->id) {
  
  case SAI_ACL_COUNTER_ATTR_PACKETS:
	req.set_packets(attr->value.u64);
	break;
  case SAI_ACL_COUNTER_ATTR_BYTES:
	req.set_bytes(attr->value.u64);
	break;
  case SAI_ACL_COUNTER_ATTR_LABEL:
	req.set_label(attr->value.chardata);
	break;
}

	grpc::Status status = acl->SetAclCounterAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_acl_counter_attribute(sai_object_id_t acl_counter_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetAclCounterAttributeRequest req;
	lemming::dataplane::sai::GetAclCounterAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(acl_counter_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_acl_counter_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = acl->GetAclCounterAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		
		

switch (attr_list[i].id) {
  
  case SAI_ACL_COUNTER_ATTR_TABLE_ID:
	 attr_list[i].value.oid =   resp.attr().table_id();
	break;
  case SAI_ACL_COUNTER_ATTR_ENABLE_PACKET_COUNT:
	 attr_list[i].value.booldata =   resp.attr().enable_packet_count();
	break;
  case SAI_ACL_COUNTER_ATTR_ENABLE_BYTE_COUNT:
	 attr_list[i].value.booldata =   resp.attr().enable_byte_count();
	break;
  case SAI_ACL_COUNTER_ATTR_PACKETS:
	 attr_list[i].value.u64 =   resp.attr().packets();
	break;
  case SAI_ACL_COUNTER_ATTR_BYTES:
	 attr_list[i].value.u64 =   resp.attr().bytes();
	break;
  case SAI_ACL_COUNTER_ATTR_LABEL:
	strncpy(attr_list[i].value.chardata, resp.attr().label().data(), 32);
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_acl_range(sai_object_id_t *acl_range_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateAclRangeRequest req = convert_create_acl_range(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateAclRangeResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = acl->CreateAclRange(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (acl_range_id) {
	*acl_range_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_acl_range(sai_object_id_t acl_range_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveAclRangeRequest req;
	lemming::dataplane::sai::RemoveAclRangeResponse resp;
	grpc::ClientContext context;
	req.set_oid(acl_range_id); 
	
	grpc::Status status = acl->RemoveAclRange(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_acl_range_attribute(sai_object_id_t acl_range_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_acl_range_attribute(sai_object_id_t acl_range_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetAclRangeAttributeRequest req;
	lemming::dataplane::sai::GetAclRangeAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(acl_range_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_acl_range_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = acl->GetAclRangeAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		
		

switch (attr_list[i].id) {
  
  case SAI_ACL_RANGE_ATTR_TYPE:
	 attr_list[i].value.s32 =  convert_sai_acl_range_type_t_to_sai(resp.attr().type());
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_acl_table_group(sai_object_id_t *acl_table_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateAclTableGroupRequest req = convert_create_acl_table_group(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateAclTableGroupResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = acl->CreateAclTableGroup(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (acl_table_group_id) {
	*acl_table_group_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_acl_table_group(sai_object_id_t acl_table_group_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveAclTableGroupRequest req;
	lemming::dataplane::sai::RemoveAclTableGroupResponse resp;
	grpc::ClientContext context;
	req.set_oid(acl_table_group_id); 
	
	grpc::Status status = acl->RemoveAclTableGroup(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_acl_table_group_attribute(sai_object_id_t acl_table_group_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_acl_table_group_attribute(sai_object_id_t acl_table_group_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetAclTableGroupAttributeRequest req;
	lemming::dataplane::sai::GetAclTableGroupAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(acl_table_group_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_acl_table_group_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = acl->GetAclTableGroupAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		
		

switch (attr_list[i].id) {
  
  case SAI_ACL_TABLE_GROUP_ATTR_ACL_STAGE:
	 attr_list[i].value.s32 =  convert_sai_acl_stage_t_to_sai(resp.attr().acl_stage());
	break;
  case SAI_ACL_TABLE_GROUP_ATTR_ACL_BIND_POINT_TYPE_LIST:
	convert_list_sai_acl_bind_point_type_t_to_sai(attr_list[i].value.s32list.list, resp.attr().acl_bind_point_type_list(), &attr_list[i].value.s32list.count);
	break;
  case SAI_ACL_TABLE_GROUP_ATTR_TYPE:
	 attr_list[i].value.s32 =  convert_sai_acl_table_group_type_t_to_sai(resp.attr().type());
	break;
  case SAI_ACL_TABLE_GROUP_ATTR_MEMBER_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().member_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_ACL_TABLE_GROUP_ATTR_CHAIN_GROUP_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().chain_group_list(), &attr_list[i].value.objlist.count);
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_acl_table_group_member(sai_object_id_t *acl_table_group_member_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateAclTableGroupMemberRequest req = convert_create_acl_table_group_member(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateAclTableGroupMemberResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = acl->CreateAclTableGroupMember(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (acl_table_group_member_id) {
	*acl_table_group_member_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_acl_table_group_member(sai_object_id_t acl_table_group_member_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveAclTableGroupMemberRequest req;
	lemming::dataplane::sai::RemoveAclTableGroupMemberResponse resp;
	grpc::ClientContext context;
	req.set_oid(acl_table_group_member_id); 
	
	grpc::Status status = acl->RemoveAclTableGroupMember(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_acl_table_group_member_attribute(sai_object_id_t acl_table_group_member_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetAclTableGroupMemberAttributeRequest req;
	lemming::dataplane::sai::SetAclTableGroupMemberAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(acl_table_group_member_id); 
	
	
	

switch (attr->id) {
  
  case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_CHAIN_GROUP_ID:
	req.set_acl_table_chain_group_id(attr->value.oid);
	break;
}

	grpc::Status status = acl->SetAclTableGroupMemberAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_acl_table_group_member_attribute(sai_object_id_t acl_table_group_member_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetAclTableGroupMemberAttributeRequest req;
	lemming::dataplane::sai::GetAclTableGroupMemberAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(acl_table_group_member_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_acl_table_group_member_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = acl->GetAclTableGroupMemberAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		
		

switch (attr_list[i].id) {
  
  case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_GROUP_ID:
	 attr_list[i].value.oid =   resp.attr().acl_table_group_id();
	break;
  case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_ID:
	 attr_list[i].value.oid =   resp.attr().acl_table_id();
	break;
  case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_PRIORITY:
	 attr_list[i].value.u32 =   resp.attr().priority();
	break;
  case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_CHAIN_GROUP_ID:
	 attr_list[i].value.oid =   resp.attr().acl_table_chain_group_id();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_acl_table_chain_group(sai_object_id_t *acl_table_chain_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateAclTableChainGroupRequest req = convert_create_acl_table_chain_group(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateAclTableChainGroupResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = acl->CreateAclTableChainGroup(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (acl_table_chain_group_id) {
	*acl_table_chain_group_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_acl_table_chain_group(sai_object_id_t acl_table_chain_group_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveAclTableChainGroupRequest req;
	lemming::dataplane::sai::RemoveAclTableChainGroupResponse resp;
	grpc::ClientContext context;
	req.set_oid(acl_table_chain_group_id); 
	
	grpc::Status status = acl->RemoveAclTableChainGroup(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_acl_table_chain_group_attribute(sai_object_id_t acl_table_chain_group_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_acl_table_chain_group_attribute(sai_object_id_t acl_table_chain_group_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetAclTableChainGroupAttributeRequest req;
	lemming::dataplane::sai::GetAclTableChainGroupAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(acl_table_chain_group_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_acl_table_chain_group_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = acl->GetAclTableChainGroupAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		
		

switch (attr_list[i].id) {
  
  case SAI_ACL_TABLE_CHAIN_GROUP_ATTR_TYPE:
	 attr_list[i].value.s32 =  convert_sai_acl_table_chain_group_type_t_to_sai(resp.attr().type());
	break;
  case SAI_ACL_TABLE_CHAIN_GROUP_ATTR_STAGE:
	 attr_list[i].value.s32 =  convert_sai_acl_table_chain_group_stage_t_to_sai(resp.attr().stage());
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

