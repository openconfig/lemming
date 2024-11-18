
// Copyright 2024 Google LLC
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

#include "dataplane/standalone/saiserver/enum.h"

lemming::dataplane::sai::AclActionType convert_sai_acl_action_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_ACTION_TYPE_REDIRECT:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_REDIRECT;

    case SAI_ACL_ACTION_TYPE_ENDPOINT_IP:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_ENDPOINT_IP;

    case SAI_ACL_ACTION_TYPE_REDIRECT_LIST:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_REDIRECT_LIST;

    case SAI_ACL_ACTION_TYPE_PACKET_ACTION:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_PACKET_ACTION;

    case SAI_ACL_ACTION_TYPE_FLOOD:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_FLOOD;

    case SAI_ACL_ACTION_TYPE_COUNTER:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_COUNTER;

    case SAI_ACL_ACTION_TYPE_MIRROR_INGRESS:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_MIRROR_INGRESS;

    case SAI_ACL_ACTION_TYPE_MIRROR_EGRESS:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_MIRROR_EGRESS;

    case SAI_ACL_ACTION_TYPE_SET_POLICER:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_POLICER;

    case SAI_ACL_ACTION_TYPE_DECREMENT_TTL:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_DECREMENT_TTL;

    case SAI_ACL_ACTION_TYPE_SET_TC:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_TC;

    case SAI_ACL_ACTION_TYPE_SET_PACKET_COLOR:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_PACKET_COLOR;

    case SAI_ACL_ACTION_TYPE_SET_INNER_VLAN_ID:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_INNER_VLAN_ID;

    case SAI_ACL_ACTION_TYPE_SET_INNER_VLAN_PRI:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_INNER_VLAN_PRI;

    case SAI_ACL_ACTION_TYPE_SET_OUTER_VLAN_ID:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_OUTER_VLAN_ID;

    case SAI_ACL_ACTION_TYPE_SET_OUTER_VLAN_PRI:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_OUTER_VLAN_PRI;

    case SAI_ACL_ACTION_TYPE_ADD_VLAN_ID:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_ADD_VLAN_ID;

    case SAI_ACL_ACTION_TYPE_ADD_VLAN_PRI:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_ADD_VLAN_PRI;

    case SAI_ACL_ACTION_TYPE_SET_SRC_MAC:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_SRC_MAC;

    case SAI_ACL_ACTION_TYPE_SET_DST_MAC:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_DST_MAC;

    case SAI_ACL_ACTION_TYPE_SET_SRC_IP:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_SRC_IP;

    case SAI_ACL_ACTION_TYPE_SET_DST_IP:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_DST_IP;

    case SAI_ACL_ACTION_TYPE_SET_SRC_IPV6:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_SRC_IPV6;

    case SAI_ACL_ACTION_TYPE_SET_DST_IPV6:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_DST_IPV6;

    case SAI_ACL_ACTION_TYPE_SET_DSCP:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_DSCP;

    case SAI_ACL_ACTION_TYPE_SET_ECN:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_ECN;

    case SAI_ACL_ACTION_TYPE_SET_L4_SRC_PORT:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_L4_SRC_PORT;

    case SAI_ACL_ACTION_TYPE_SET_L4_DST_PORT:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_L4_DST_PORT;

    case SAI_ACL_ACTION_TYPE_INGRESS_SAMPLEPACKET_ENABLE:
      return lemming::dataplane::sai::
          ACL_ACTION_TYPE_INGRESS_SAMPLEPACKET_ENABLE;

    case SAI_ACL_ACTION_TYPE_EGRESS_SAMPLEPACKET_ENABLE:
      return lemming::dataplane::sai::
          ACL_ACTION_TYPE_EGRESS_SAMPLEPACKET_ENABLE;

    case SAI_ACL_ACTION_TYPE_SET_ACL_META_DATA:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_ACL_META_DATA;

    case SAI_ACL_ACTION_TYPE_EGRESS_BLOCK_PORT_LIST:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_EGRESS_BLOCK_PORT_LIST;

    case SAI_ACL_ACTION_TYPE_SET_USER_TRAP_ID:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_USER_TRAP_ID;

    case SAI_ACL_ACTION_TYPE_SET_DO_NOT_LEARN:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_DO_NOT_LEARN;

    case SAI_ACL_ACTION_TYPE_ACL_DTEL_FLOW_OP:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_ACL_DTEL_FLOW_OP;

    case SAI_ACL_ACTION_TYPE_DTEL_INT_SESSION:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_DTEL_INT_SESSION;

    case SAI_ACL_ACTION_TYPE_DTEL_DROP_REPORT_ENABLE:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_DTEL_DROP_REPORT_ENABLE;

    case SAI_ACL_ACTION_TYPE_DTEL_TAIL_DROP_REPORT_ENABLE:
      return lemming::dataplane::sai::
          ACL_ACTION_TYPE_DTEL_TAIL_DROP_REPORT_ENABLE;

    case SAI_ACL_ACTION_TYPE_DTEL_FLOW_SAMPLE_PERCENT:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_DTEL_FLOW_SAMPLE_PERCENT;

    case SAI_ACL_ACTION_TYPE_DTEL_REPORT_ALL_PACKETS:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_DTEL_REPORT_ALL_PACKETS;

    case SAI_ACL_ACTION_TYPE_NO_NAT:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_NO_NAT;

    case SAI_ACL_ACTION_TYPE_INT_INSERT:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_INT_INSERT;

    case SAI_ACL_ACTION_TYPE_INT_DELETE:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_INT_DELETE;

    case SAI_ACL_ACTION_TYPE_INT_REPORT_FLOW:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_INT_REPORT_FLOW;

    case SAI_ACL_ACTION_TYPE_INT_REPORT_DROPS:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_INT_REPORT_DROPS;

    case SAI_ACL_ACTION_TYPE_INT_REPORT_TAIL_DROPS:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_INT_REPORT_TAIL_DROPS;

    case SAI_ACL_ACTION_TYPE_TAM_INT_OBJECT:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_TAM_INT_OBJECT;

    case SAI_ACL_ACTION_TYPE_SET_ISOLATION_GROUP:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_ISOLATION_GROUP;

    case SAI_ACL_ACTION_TYPE_MACSEC_FLOW:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_MACSEC_FLOW;

    case SAI_ACL_ACTION_TYPE_SET_LAG_HASH_ID:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_LAG_HASH_ID;

    case SAI_ACL_ACTION_TYPE_SET_ECMP_HASH_ID:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_ECMP_HASH_ID;

    case SAI_ACL_ACTION_TYPE_SET_VRF:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_VRF;

    case SAI_ACL_ACTION_TYPE_SET_FORWARDING_CLASS:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_SET_FORWARDING_CLASS;

    default:
      return lemming::dataplane::sai::ACL_ACTION_TYPE_UNSPECIFIED;
  }
}
sai_acl_action_type_t convert_sai_acl_action_type_t_to_sai(
    lemming::dataplane::sai::AclActionType val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_ACTION_TYPE_REDIRECT:
      return SAI_ACL_ACTION_TYPE_REDIRECT;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_ENDPOINT_IP:
      return SAI_ACL_ACTION_TYPE_ENDPOINT_IP;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_REDIRECT_LIST:
      return SAI_ACL_ACTION_TYPE_REDIRECT_LIST;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_PACKET_ACTION:
      return SAI_ACL_ACTION_TYPE_PACKET_ACTION;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_FLOOD:
      return SAI_ACL_ACTION_TYPE_FLOOD;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_COUNTER:
      return SAI_ACL_ACTION_TYPE_COUNTER;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_MIRROR_INGRESS:
      return SAI_ACL_ACTION_TYPE_MIRROR_INGRESS;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_MIRROR_EGRESS:
      return SAI_ACL_ACTION_TYPE_MIRROR_EGRESS;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_POLICER:
      return SAI_ACL_ACTION_TYPE_SET_POLICER;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_DECREMENT_TTL:
      return SAI_ACL_ACTION_TYPE_DECREMENT_TTL;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_TC:
      return SAI_ACL_ACTION_TYPE_SET_TC;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_PACKET_COLOR:
      return SAI_ACL_ACTION_TYPE_SET_PACKET_COLOR;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_INNER_VLAN_ID:
      return SAI_ACL_ACTION_TYPE_SET_INNER_VLAN_ID;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_INNER_VLAN_PRI:
      return SAI_ACL_ACTION_TYPE_SET_INNER_VLAN_PRI;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_OUTER_VLAN_ID:
      return SAI_ACL_ACTION_TYPE_SET_OUTER_VLAN_ID;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_OUTER_VLAN_PRI:
      return SAI_ACL_ACTION_TYPE_SET_OUTER_VLAN_PRI;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_ADD_VLAN_ID:
      return SAI_ACL_ACTION_TYPE_ADD_VLAN_ID;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_ADD_VLAN_PRI:
      return SAI_ACL_ACTION_TYPE_ADD_VLAN_PRI;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_SRC_MAC:
      return SAI_ACL_ACTION_TYPE_SET_SRC_MAC;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_DST_MAC:
      return SAI_ACL_ACTION_TYPE_SET_DST_MAC;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_SRC_IP:
      return SAI_ACL_ACTION_TYPE_SET_SRC_IP;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_DST_IP:
      return SAI_ACL_ACTION_TYPE_SET_DST_IP;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_SRC_IPV6:
      return SAI_ACL_ACTION_TYPE_SET_SRC_IPV6;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_DST_IPV6:
      return SAI_ACL_ACTION_TYPE_SET_DST_IPV6;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_DSCP:
      return SAI_ACL_ACTION_TYPE_SET_DSCP;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_ECN:
      return SAI_ACL_ACTION_TYPE_SET_ECN;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_L4_SRC_PORT:
      return SAI_ACL_ACTION_TYPE_SET_L4_SRC_PORT;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_L4_DST_PORT:
      return SAI_ACL_ACTION_TYPE_SET_L4_DST_PORT;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_INGRESS_SAMPLEPACKET_ENABLE:
      return SAI_ACL_ACTION_TYPE_INGRESS_SAMPLEPACKET_ENABLE;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_EGRESS_SAMPLEPACKET_ENABLE:
      return SAI_ACL_ACTION_TYPE_EGRESS_SAMPLEPACKET_ENABLE;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_ACL_META_DATA:
      return SAI_ACL_ACTION_TYPE_SET_ACL_META_DATA;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_EGRESS_BLOCK_PORT_LIST:
      return SAI_ACL_ACTION_TYPE_EGRESS_BLOCK_PORT_LIST;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_USER_TRAP_ID:
      return SAI_ACL_ACTION_TYPE_SET_USER_TRAP_ID;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_DO_NOT_LEARN:
      return SAI_ACL_ACTION_TYPE_SET_DO_NOT_LEARN;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_ACL_DTEL_FLOW_OP:
      return SAI_ACL_ACTION_TYPE_ACL_DTEL_FLOW_OP;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_DTEL_INT_SESSION:
      return SAI_ACL_ACTION_TYPE_DTEL_INT_SESSION;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_DTEL_DROP_REPORT_ENABLE:
      return SAI_ACL_ACTION_TYPE_DTEL_DROP_REPORT_ENABLE;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_DTEL_TAIL_DROP_REPORT_ENABLE:
      return SAI_ACL_ACTION_TYPE_DTEL_TAIL_DROP_REPORT_ENABLE;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_DTEL_FLOW_SAMPLE_PERCENT:
      return SAI_ACL_ACTION_TYPE_DTEL_FLOW_SAMPLE_PERCENT;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_DTEL_REPORT_ALL_PACKETS:
      return SAI_ACL_ACTION_TYPE_DTEL_REPORT_ALL_PACKETS;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_NO_NAT:
      return SAI_ACL_ACTION_TYPE_NO_NAT;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_INT_INSERT:
      return SAI_ACL_ACTION_TYPE_INT_INSERT;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_INT_DELETE:
      return SAI_ACL_ACTION_TYPE_INT_DELETE;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_INT_REPORT_FLOW:
      return SAI_ACL_ACTION_TYPE_INT_REPORT_FLOW;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_INT_REPORT_DROPS:
      return SAI_ACL_ACTION_TYPE_INT_REPORT_DROPS;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_INT_REPORT_TAIL_DROPS:
      return SAI_ACL_ACTION_TYPE_INT_REPORT_TAIL_DROPS;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_TAM_INT_OBJECT:
      return SAI_ACL_ACTION_TYPE_TAM_INT_OBJECT;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_ISOLATION_GROUP:
      return SAI_ACL_ACTION_TYPE_SET_ISOLATION_GROUP;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_MACSEC_FLOW:
      return SAI_ACL_ACTION_TYPE_MACSEC_FLOW;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_LAG_HASH_ID:
      return SAI_ACL_ACTION_TYPE_SET_LAG_HASH_ID;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_ECMP_HASH_ID:
      return SAI_ACL_ACTION_TYPE_SET_ECMP_HASH_ID;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_VRF:
      return SAI_ACL_ACTION_TYPE_SET_VRF;

    case lemming::dataplane::sai::ACL_ACTION_TYPE_SET_FORWARDING_CLASS:
      return SAI_ACL_ACTION_TYPE_SET_FORWARDING_CLASS;

    default:
      return SAI_ACL_ACTION_TYPE_REDIRECT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_acl_action_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_action_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_action_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_action_type_t_to_sai(
        static_cast<lemming::dataplane::sai::AclActionType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclBindPointType
convert_sai_acl_bind_point_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_BIND_POINT_TYPE_PORT:
      return lemming::dataplane::sai::ACL_BIND_POINT_TYPE_PORT;

    case SAI_ACL_BIND_POINT_TYPE_LAG:
      return lemming::dataplane::sai::ACL_BIND_POINT_TYPE_LAG;

    case SAI_ACL_BIND_POINT_TYPE_VLAN:
      return lemming::dataplane::sai::ACL_BIND_POINT_TYPE_VLAN;

    case SAI_ACL_BIND_POINT_TYPE_ROUTER_INTERFACE:
      return lemming::dataplane::sai::ACL_BIND_POINT_TYPE_ROUTER_INTERFACE;

    case SAI_ACL_BIND_POINT_TYPE_SWITCH:
      return lemming::dataplane::sai::ACL_BIND_POINT_TYPE_SWITCH;

    default:
      return lemming::dataplane::sai::ACL_BIND_POINT_TYPE_UNSPECIFIED;
  }
}
sai_acl_bind_point_type_t convert_sai_acl_bind_point_type_t_to_sai(
    lemming::dataplane::sai::AclBindPointType val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_BIND_POINT_TYPE_PORT:
      return SAI_ACL_BIND_POINT_TYPE_PORT;

    case lemming::dataplane::sai::ACL_BIND_POINT_TYPE_LAG:
      return SAI_ACL_BIND_POINT_TYPE_LAG;

    case lemming::dataplane::sai::ACL_BIND_POINT_TYPE_VLAN:
      return SAI_ACL_BIND_POINT_TYPE_VLAN;

    case lemming::dataplane::sai::ACL_BIND_POINT_TYPE_ROUTER_INTERFACE:
      return SAI_ACL_BIND_POINT_TYPE_ROUTER_INTERFACE;

    case lemming::dataplane::sai::ACL_BIND_POINT_TYPE_SWITCH:
      return SAI_ACL_BIND_POINT_TYPE_SWITCH;

    default:
      return SAI_ACL_BIND_POINT_TYPE_PORT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_acl_bind_point_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_bind_point_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_bind_point_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_bind_point_type_t_to_sai(
        static_cast<lemming::dataplane::sai::AclBindPointType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclCounterAttr convert_sai_acl_counter_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_COUNTER_ATTR_TABLE_ID:
      return lemming::dataplane::sai::ACL_COUNTER_ATTR_TABLE_ID;

    case SAI_ACL_COUNTER_ATTR_ENABLE_PACKET_COUNT:
      return lemming::dataplane::sai::ACL_COUNTER_ATTR_ENABLE_PACKET_COUNT;

    case SAI_ACL_COUNTER_ATTR_ENABLE_BYTE_COUNT:
      return lemming::dataplane::sai::ACL_COUNTER_ATTR_ENABLE_BYTE_COUNT;

    case SAI_ACL_COUNTER_ATTR_PACKETS:
      return lemming::dataplane::sai::ACL_COUNTER_ATTR_PACKETS;

    case SAI_ACL_COUNTER_ATTR_BYTES:
      return lemming::dataplane::sai::ACL_COUNTER_ATTR_BYTES;

    case SAI_ACL_COUNTER_ATTR_LABEL:
      return lemming::dataplane::sai::ACL_COUNTER_ATTR_LABEL;

    default:
      return lemming::dataplane::sai::ACL_COUNTER_ATTR_UNSPECIFIED;
  }
}
sai_acl_counter_attr_t convert_sai_acl_counter_attr_t_to_sai(
    lemming::dataplane::sai::AclCounterAttr val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_COUNTER_ATTR_TABLE_ID:
      return SAI_ACL_COUNTER_ATTR_TABLE_ID;

    case lemming::dataplane::sai::ACL_COUNTER_ATTR_ENABLE_PACKET_COUNT:
      return SAI_ACL_COUNTER_ATTR_ENABLE_PACKET_COUNT;

    case lemming::dataplane::sai::ACL_COUNTER_ATTR_ENABLE_BYTE_COUNT:
      return SAI_ACL_COUNTER_ATTR_ENABLE_BYTE_COUNT;

    case lemming::dataplane::sai::ACL_COUNTER_ATTR_PACKETS:
      return SAI_ACL_COUNTER_ATTR_PACKETS;

    case lemming::dataplane::sai::ACL_COUNTER_ATTR_BYTES:
      return SAI_ACL_COUNTER_ATTR_BYTES;

    case lemming::dataplane::sai::ACL_COUNTER_ATTR_LABEL:
      return SAI_ACL_COUNTER_ATTR_LABEL;

    default:
      return SAI_ACL_COUNTER_ATTR_TABLE_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_acl_counter_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_counter_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_counter_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_counter_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::AclCounterAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclDtelFlowOp convert_sai_acl_dtel_flow_op_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_DTEL_FLOW_OP_NOP:
      return lemming::dataplane::sai::ACL_DTEL_FLOW_OP_NOP;

    case SAI_ACL_DTEL_FLOW_OP_INT:
      return lemming::dataplane::sai::ACL_DTEL_FLOW_OP_INT;

    case SAI_ACL_DTEL_FLOW_OP_IOAM:
      return lemming::dataplane::sai::ACL_DTEL_FLOW_OP_IOAM;

    case SAI_ACL_DTEL_FLOW_OP_POSTCARD:
      return lemming::dataplane::sai::ACL_DTEL_FLOW_OP_POSTCARD;

    default:
      return lemming::dataplane::sai::ACL_DTEL_FLOW_OP_UNSPECIFIED;
  }
}
sai_acl_dtel_flow_op_t convert_sai_acl_dtel_flow_op_t_to_sai(
    lemming::dataplane::sai::AclDtelFlowOp val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_DTEL_FLOW_OP_NOP:
      return SAI_ACL_DTEL_FLOW_OP_NOP;

    case lemming::dataplane::sai::ACL_DTEL_FLOW_OP_INT:
      return SAI_ACL_DTEL_FLOW_OP_INT;

    case lemming::dataplane::sai::ACL_DTEL_FLOW_OP_IOAM:
      return SAI_ACL_DTEL_FLOW_OP_IOAM;

    case lemming::dataplane::sai::ACL_DTEL_FLOW_OP_POSTCARD:
      return SAI_ACL_DTEL_FLOW_OP_POSTCARD;

    default:
      return SAI_ACL_DTEL_FLOW_OP_NOP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_acl_dtel_flow_op_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_dtel_flow_op_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_dtel_flow_op_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_dtel_flow_op_t_to_sai(
        static_cast<lemming::dataplane::sai::AclDtelFlowOp>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclEntryAttr convert_sai_acl_entry_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_ENTRY_ATTR_TABLE_ID:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_TABLE_ID;

    case SAI_ACL_ENTRY_ATTR_PRIORITY:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_PRIORITY;

    case SAI_ACL_ENTRY_ATTR_ADMIN_STATE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ADMIN_STATE;

    case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IPV6;

    case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD3:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD3;

    case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD2:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD2;

    case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD1:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD1;

    case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD0:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD0;

    case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IPV6;

    case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD3:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD3;

    case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD2:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD2;

    case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD1:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD1;

    case SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD0:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD0;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_IPV6:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_SRC_IPV6;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_IPV6:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_DST_IPV6;

    case SAI_ACL_ENTRY_ATTR_FIELD_SRC_MAC:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_MAC;

    case SAI_ACL_ENTRY_ATTR_FIELD_DST_MAC:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_MAC;

    case SAI_ACL_ENTRY_ATTR_FIELD_SRC_IP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IP;

    case SAI_ACL_ENTRY_ATTR_FIELD_DST_IP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IP;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_IP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_SRC_IP;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_IP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_DST_IP;

    case SAI_ACL_ENTRY_ATTR_FIELD_IN_PORTS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IN_PORTS;

    case SAI_ACL_ENTRY_ATTR_FIELD_OUT_PORTS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_OUT_PORTS;

    case SAI_ACL_ENTRY_ATTR_FIELD_IN_PORT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IN_PORT;

    case SAI_ACL_ENTRY_ATTR_FIELD_OUT_PORT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_OUT_PORT;

    case SAI_ACL_ENTRY_ATTR_FIELD_SRC_PORT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_PORT;

    case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_ID:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_ID;

    case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_PRI:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_PRI;

    case SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_CFI:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_CFI;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_ID:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_VLAN_ID;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_PRI:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_VLAN_PRI;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_CFI:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_VLAN_CFI;

    case SAI_ACL_ENTRY_ATTR_FIELD_L4_SRC_PORT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_L4_SRC_PORT;

    case SAI_ACL_ENTRY_ATTR_FIELD_L4_DST_PORT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_L4_DST_PORT;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_L4_SRC_PORT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_L4_SRC_PORT;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_L4_DST_PORT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_L4_DST_PORT;

    case SAI_ACL_ENTRY_ATTR_FIELD_ETHER_TYPE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ETHER_TYPE;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_ETHER_TYPE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_ETHER_TYPE;

    case SAI_ACL_ENTRY_ATTR_FIELD_IP_PROTOCOL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IP_PROTOCOL;

    case SAI_ACL_ENTRY_ATTR_FIELD_INNER_IP_PROTOCOL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_IP_PROTOCOL;

    case SAI_ACL_ENTRY_ATTR_FIELD_IP_IDENTIFICATION:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IP_IDENTIFICATION;

    case SAI_ACL_ENTRY_ATTR_FIELD_DSCP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DSCP;

    case SAI_ACL_ENTRY_ATTR_FIELD_ECN:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ECN;

    case SAI_ACL_ENTRY_ATTR_FIELD_TTL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TTL;

    case SAI_ACL_ENTRY_ATTR_FIELD_TOS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TOS;

    case SAI_ACL_ENTRY_ATTR_FIELD_IP_FLAGS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IP_FLAGS;

    case SAI_ACL_ENTRY_ATTR_FIELD_TCP_FLAGS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TCP_FLAGS;

    case SAI_ACL_ENTRY_ATTR_FIELD_ACL_IP_TYPE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ACL_IP_TYPE;

    case SAI_ACL_ENTRY_ATTR_FIELD_ACL_IP_FRAG:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ACL_IP_FRAG;

    case SAI_ACL_ENTRY_ATTR_FIELD_IPV6_FLOW_LABEL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IPV6_FLOW_LABEL;

    case SAI_ACL_ENTRY_ATTR_FIELD_TC:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TC;

    case SAI_ACL_ENTRY_ATTR_FIELD_ICMP_TYPE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ICMP_TYPE;

    case SAI_ACL_ENTRY_ATTR_FIELD_ICMP_CODE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ICMP_CODE;

    case SAI_ACL_ENTRY_ATTR_FIELD_ICMPV6_TYPE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ICMPV6_TYPE;

    case SAI_ACL_ENTRY_ATTR_FIELD_ICMPV6_CODE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ICMPV6_CODE;

    case SAI_ACL_ENTRY_ATTR_FIELD_PACKET_VLAN:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_PACKET_VLAN;

    case SAI_ACL_ENTRY_ATTR_FIELD_TUNNEL_VNI:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TUNNEL_VNI;

    case SAI_ACL_ENTRY_ATTR_FIELD_HAS_VLAN_TAG:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_HAS_VLAN_TAG;

    case SAI_ACL_ENTRY_ATTR_FIELD_MACSEC_SCI:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MACSEC_SCI;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_LABEL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_LABEL;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_TTL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_TTL;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_EXP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_EXP;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_BOS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_BOS;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_LABEL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_LABEL;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_TTL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_TTL;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_EXP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_EXP;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_BOS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_BOS;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_LABEL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_LABEL;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_TTL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_TTL;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_EXP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_EXP;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_BOS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_BOS;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_LABEL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_LABEL;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_TTL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_TTL;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_EXP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_EXP;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_BOS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_BOS;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_LABEL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_LABEL;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_TTL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_TTL;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_EXP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_EXP;

    case SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_BOS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_BOS;

    case SAI_ACL_ENTRY_ATTR_FIELD_FDB_DST_USER_META:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_FDB_DST_USER_META;

    case SAI_ACL_ENTRY_ATTR_FIELD_ROUTE_DST_USER_META:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ROUTE_DST_USER_META;

    case SAI_ACL_ENTRY_ATTR_FIELD_NEIGHBOR_DST_USER_META:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_FIELD_NEIGHBOR_DST_USER_META;

    case SAI_ACL_ENTRY_ATTR_FIELD_PORT_USER_META:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_PORT_USER_META;

    case SAI_ACL_ENTRY_ATTR_FIELD_VLAN_USER_META:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_VLAN_USER_META;

    case SAI_ACL_ENTRY_ATTR_FIELD_ACL_USER_META:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ACL_USER_META;

    case SAI_ACL_ENTRY_ATTR_FIELD_FDB_NPU_META_DST_HIT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_FDB_NPU_META_DST_HIT;

    case SAI_ACL_ENTRY_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT;

    case SAI_ACL_ENTRY_ATTR_FIELD_ROUTE_NPU_META_DST_HIT:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_FIELD_ROUTE_NPU_META_DST_HIT;

    case SAI_ACL_ENTRY_ATTR_FIELD_BTH_OPCODE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_BTH_OPCODE;

    case SAI_ACL_ENTRY_ATTR_FIELD_AETH_SYNDROME:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_AETH_SYNDROME;

    case SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN;

    case SAI_ACL_ENTRY_ATTR_FIELD_ACL_RANGE_TYPE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ACL_RANGE_TYPE;

    case SAI_ACL_ENTRY_ATTR_FIELD_IPV6_NEXT_HEADER:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IPV6_NEXT_HEADER;

    case SAI_ACL_ENTRY_ATTR_FIELD_GRE_KEY:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_GRE_KEY;

    case SAI_ACL_ENTRY_ATTR_FIELD_TAM_INT_TYPE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TAM_INT_TYPE;

    case SAI_ACL_ENTRY_ATTR_ACTION_REDIRECT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_REDIRECT;

    case SAI_ACL_ENTRY_ATTR_ACTION_ENDPOINT_IP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_ENDPOINT_IP;

    case SAI_ACL_ENTRY_ATTR_ACTION_REDIRECT_LIST:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_REDIRECT_LIST;

    case SAI_ACL_ENTRY_ATTR_ACTION_PACKET_ACTION:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_PACKET_ACTION;

    case SAI_ACL_ENTRY_ATTR_ACTION_FLOOD:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_FLOOD;

    case SAI_ACL_ENTRY_ATTR_ACTION_COUNTER:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_COUNTER;

    case SAI_ACL_ENTRY_ATTR_ACTION_MIRROR_INGRESS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_MIRROR_INGRESS;

    case SAI_ACL_ENTRY_ATTR_ACTION_MIRROR_EGRESS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_MIRROR_EGRESS;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_POLICER:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_POLICER;

    case SAI_ACL_ENTRY_ATTR_ACTION_DECREMENT_TTL:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_DECREMENT_TTL;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_TC:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_TC;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_PACKET_COLOR:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_PACKET_COLOR;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_INNER_VLAN_ID:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_INNER_VLAN_ID;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_INNER_VLAN_PRI:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_INNER_VLAN_PRI;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_OUTER_VLAN_ID:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_OUTER_VLAN_ID;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_OUTER_VLAN_PRI:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_OUTER_VLAN_PRI;

    case SAI_ACL_ENTRY_ATTR_ACTION_ADD_VLAN_ID:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_ADD_VLAN_ID;

    case SAI_ACL_ENTRY_ATTR_ACTION_ADD_VLAN_PRI:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_ADD_VLAN_PRI;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_SRC_MAC:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_SRC_MAC;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_DST_MAC:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_DST_MAC;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_SRC_IP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_SRC_IP;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_DST_IP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_DST_IP;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_SRC_IPV6:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_SRC_IPV6;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_DST_IPV6:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_DST_IPV6;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_DSCP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_DSCP;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_ECN:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_ECN;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_L4_SRC_PORT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_L4_SRC_PORT;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_L4_DST_PORT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_L4_DST_PORT;

    case SAI_ACL_ENTRY_ATTR_ACTION_INGRESS_SAMPLEPACKET_ENABLE:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_ACTION_INGRESS_SAMPLEPACKET_ENABLE;

    case SAI_ACL_ENTRY_ATTR_ACTION_EGRESS_SAMPLEPACKET_ENABLE:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_ACTION_EGRESS_SAMPLEPACKET_ENABLE;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_ACL_META_DATA:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_ACL_META_DATA;

    case SAI_ACL_ENTRY_ATTR_ACTION_EGRESS_BLOCK_PORT_LIST:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_ACTION_EGRESS_BLOCK_PORT_LIST;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_USER_TRAP_ID:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_USER_TRAP_ID;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_DO_NOT_LEARN:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_DO_NOT_LEARN;

    case SAI_ACL_ENTRY_ATTR_ACTION_ACL_DTEL_FLOW_OP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_ACL_DTEL_FLOW_OP;

    case SAI_ACL_ENTRY_ATTR_ACTION_DTEL_INT_SESSION:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_DTEL_INT_SESSION;

    case SAI_ACL_ENTRY_ATTR_ACTION_DTEL_DROP_REPORT_ENABLE:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_ACTION_DTEL_DROP_REPORT_ENABLE;

    case SAI_ACL_ENTRY_ATTR_ACTION_DTEL_TAIL_DROP_REPORT_ENABLE:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_ACTION_DTEL_TAIL_DROP_REPORT_ENABLE;

    case SAI_ACL_ENTRY_ATTR_ACTION_DTEL_FLOW_SAMPLE_PERCENT:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_ACTION_DTEL_FLOW_SAMPLE_PERCENT;

    case SAI_ACL_ENTRY_ATTR_ACTION_DTEL_REPORT_ALL_PACKETS:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_ACTION_DTEL_REPORT_ALL_PACKETS;

    case SAI_ACL_ENTRY_ATTR_ACTION_NO_NAT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_NO_NAT;

    case SAI_ACL_ENTRY_ATTR_ACTION_INT_INSERT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_INT_INSERT;

    case SAI_ACL_ENTRY_ATTR_ACTION_INT_DELETE:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_INT_DELETE;

    case SAI_ACL_ENTRY_ATTR_ACTION_INT_REPORT_FLOW:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_INT_REPORT_FLOW;

    case SAI_ACL_ENTRY_ATTR_ACTION_INT_REPORT_DROPS:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_INT_REPORT_DROPS;

    case SAI_ACL_ENTRY_ATTR_ACTION_INT_REPORT_TAIL_DROPS:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_ACTION_INT_REPORT_TAIL_DROPS;

    case SAI_ACL_ENTRY_ATTR_ACTION_TAM_INT_OBJECT:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_TAM_INT_OBJECT;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_ISOLATION_GROUP:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_ISOLATION_GROUP;

    case SAI_ACL_ENTRY_ATTR_ACTION_MACSEC_FLOW:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_MACSEC_FLOW;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_LAG_HASH_ID:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_LAG_HASH_ID;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_ECMP_HASH_ID:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_ECMP_HASH_ID;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_VRF:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_VRF;

    case SAI_ACL_ENTRY_ATTR_ACTION_SET_FORWARDING_CLASS:
      return lemming::dataplane::sai::
          ACL_ENTRY_ATTR_ACTION_SET_FORWARDING_CLASS;

    default:
      return lemming::dataplane::sai::ACL_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_acl_entry_attr_t convert_sai_acl_entry_attr_t_to_sai(
    lemming::dataplane::sai::AclEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_ENTRY_ATTR_TABLE_ID:
      return SAI_ACL_ENTRY_ATTR_TABLE_ID;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_PRIORITY:
      return SAI_ACL_ENTRY_ATTR_PRIORITY;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ADMIN_STATE:
      return SAI_ACL_ENTRY_ATTR_ADMIN_STATE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IPV6:
      return SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD3:
      return SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD3;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD2:
      return SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD2;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD1:
      return SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD1;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD0:
      return SAI_ACL_ENTRY_ATTR_FIELD_SRC_IPV6_WORD0;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IPV6:
      return SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD3:
      return SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD3;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD2:
      return SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD2;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD1:
      return SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD1;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD0:
      return SAI_ACL_ENTRY_ATTR_FIELD_DST_IPV6_WORD0;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_SRC_IPV6:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_IPV6;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_DST_IPV6:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_IPV6;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_MAC:
      return SAI_ACL_ENTRY_ATTR_FIELD_SRC_MAC;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_MAC:
      return SAI_ACL_ENTRY_ATTR_FIELD_DST_MAC;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_IP:
      return SAI_ACL_ENTRY_ATTR_FIELD_SRC_IP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DST_IP:
      return SAI_ACL_ENTRY_ATTR_FIELD_DST_IP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_SRC_IP:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_SRC_IP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_DST_IP:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_DST_IP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IN_PORTS:
      return SAI_ACL_ENTRY_ATTR_FIELD_IN_PORTS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_OUT_PORTS:
      return SAI_ACL_ENTRY_ATTR_FIELD_OUT_PORTS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IN_PORT:
      return SAI_ACL_ENTRY_ATTR_FIELD_IN_PORT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_OUT_PORT:
      return SAI_ACL_ENTRY_ATTR_FIELD_OUT_PORT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_SRC_PORT:
      return SAI_ACL_ENTRY_ATTR_FIELD_SRC_PORT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_ID:
      return SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_ID;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_PRI:
      return SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_PRI;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_CFI:
      return SAI_ACL_ENTRY_ATTR_FIELD_OUTER_VLAN_CFI;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_VLAN_ID:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_ID;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_VLAN_PRI:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_PRI;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_VLAN_CFI:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_VLAN_CFI;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_L4_SRC_PORT:
      return SAI_ACL_ENTRY_ATTR_FIELD_L4_SRC_PORT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_L4_DST_PORT:
      return SAI_ACL_ENTRY_ATTR_FIELD_L4_DST_PORT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_L4_SRC_PORT:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_L4_SRC_PORT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_L4_DST_PORT:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_L4_DST_PORT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ETHER_TYPE:
      return SAI_ACL_ENTRY_ATTR_FIELD_ETHER_TYPE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_ETHER_TYPE:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_ETHER_TYPE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IP_PROTOCOL:
      return SAI_ACL_ENTRY_ATTR_FIELD_IP_PROTOCOL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_INNER_IP_PROTOCOL:
      return SAI_ACL_ENTRY_ATTR_FIELD_INNER_IP_PROTOCOL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IP_IDENTIFICATION:
      return SAI_ACL_ENTRY_ATTR_FIELD_IP_IDENTIFICATION;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_DSCP:
      return SAI_ACL_ENTRY_ATTR_FIELD_DSCP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ECN:
      return SAI_ACL_ENTRY_ATTR_FIELD_ECN;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TTL:
      return SAI_ACL_ENTRY_ATTR_FIELD_TTL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TOS:
      return SAI_ACL_ENTRY_ATTR_FIELD_TOS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IP_FLAGS:
      return SAI_ACL_ENTRY_ATTR_FIELD_IP_FLAGS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TCP_FLAGS:
      return SAI_ACL_ENTRY_ATTR_FIELD_TCP_FLAGS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ACL_IP_TYPE:
      return SAI_ACL_ENTRY_ATTR_FIELD_ACL_IP_TYPE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ACL_IP_FRAG:
      return SAI_ACL_ENTRY_ATTR_FIELD_ACL_IP_FRAG;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IPV6_FLOW_LABEL:
      return SAI_ACL_ENTRY_ATTR_FIELD_IPV6_FLOW_LABEL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TC:
      return SAI_ACL_ENTRY_ATTR_FIELD_TC;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ICMP_TYPE:
      return SAI_ACL_ENTRY_ATTR_FIELD_ICMP_TYPE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ICMP_CODE:
      return SAI_ACL_ENTRY_ATTR_FIELD_ICMP_CODE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ICMPV6_TYPE:
      return SAI_ACL_ENTRY_ATTR_FIELD_ICMPV6_TYPE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ICMPV6_CODE:
      return SAI_ACL_ENTRY_ATTR_FIELD_ICMPV6_CODE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_PACKET_VLAN:
      return SAI_ACL_ENTRY_ATTR_FIELD_PACKET_VLAN;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TUNNEL_VNI:
      return SAI_ACL_ENTRY_ATTR_FIELD_TUNNEL_VNI;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_HAS_VLAN_TAG:
      return SAI_ACL_ENTRY_ATTR_FIELD_HAS_VLAN_TAG;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MACSEC_SCI:
      return SAI_ACL_ENTRY_ATTR_FIELD_MACSEC_SCI;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_LABEL:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_LABEL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_TTL:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_TTL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_EXP:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_EXP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_BOS:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL0_BOS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_LABEL:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_LABEL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_TTL:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_TTL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_EXP:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_EXP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_BOS:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL1_BOS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_LABEL:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_LABEL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_TTL:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_TTL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_EXP:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_EXP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_BOS:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL2_BOS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_LABEL:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_LABEL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_TTL:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_TTL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_EXP:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_EXP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_BOS:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL3_BOS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_LABEL:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_LABEL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_TTL:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_TTL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_EXP:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_EXP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_BOS:
      return SAI_ACL_ENTRY_ATTR_FIELD_MPLS_LABEL4_BOS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_FDB_DST_USER_META:
      return SAI_ACL_ENTRY_ATTR_FIELD_FDB_DST_USER_META;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ROUTE_DST_USER_META:
      return SAI_ACL_ENTRY_ATTR_FIELD_ROUTE_DST_USER_META;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_NEIGHBOR_DST_USER_META:
      return SAI_ACL_ENTRY_ATTR_FIELD_NEIGHBOR_DST_USER_META;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_PORT_USER_META:
      return SAI_ACL_ENTRY_ATTR_FIELD_PORT_USER_META;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_VLAN_USER_META:
      return SAI_ACL_ENTRY_ATTR_FIELD_VLAN_USER_META;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ACL_USER_META:
      return SAI_ACL_ENTRY_ATTR_FIELD_ACL_USER_META;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_FDB_NPU_META_DST_HIT:
      return SAI_ACL_ENTRY_ATTR_FIELD_FDB_NPU_META_DST_HIT;

    case lemming::dataplane::sai::
        ACL_ENTRY_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT:
      return SAI_ACL_ENTRY_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ROUTE_NPU_META_DST_HIT:
      return SAI_ACL_ENTRY_ATTR_FIELD_ROUTE_NPU_META_DST_HIT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_BTH_OPCODE:
      return SAI_ACL_ENTRY_ATTR_FIELD_BTH_OPCODE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_AETH_SYNDROME:
      return SAI_ACL_ENTRY_ATTR_FIELD_AETH_SYNDROME;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN:
      return SAI_ACL_ENTRY_ATTR_USER_DEFINED_FIELD_GROUP_MIN;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_ACL_RANGE_TYPE:
      return SAI_ACL_ENTRY_ATTR_FIELD_ACL_RANGE_TYPE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_IPV6_NEXT_HEADER:
      return SAI_ACL_ENTRY_ATTR_FIELD_IPV6_NEXT_HEADER;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_GRE_KEY:
      return SAI_ACL_ENTRY_ATTR_FIELD_GRE_KEY;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_FIELD_TAM_INT_TYPE:
      return SAI_ACL_ENTRY_ATTR_FIELD_TAM_INT_TYPE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_REDIRECT:
      return SAI_ACL_ENTRY_ATTR_ACTION_REDIRECT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_ENDPOINT_IP:
      return SAI_ACL_ENTRY_ATTR_ACTION_ENDPOINT_IP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_REDIRECT_LIST:
      return SAI_ACL_ENTRY_ATTR_ACTION_REDIRECT_LIST;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_PACKET_ACTION:
      return SAI_ACL_ENTRY_ATTR_ACTION_PACKET_ACTION;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_FLOOD:
      return SAI_ACL_ENTRY_ATTR_ACTION_FLOOD;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_COUNTER:
      return SAI_ACL_ENTRY_ATTR_ACTION_COUNTER;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_MIRROR_INGRESS:
      return SAI_ACL_ENTRY_ATTR_ACTION_MIRROR_INGRESS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_MIRROR_EGRESS:
      return SAI_ACL_ENTRY_ATTR_ACTION_MIRROR_EGRESS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_POLICER:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_POLICER;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_DECREMENT_TTL:
      return SAI_ACL_ENTRY_ATTR_ACTION_DECREMENT_TTL;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_TC:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_TC;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_PACKET_COLOR:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_PACKET_COLOR;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_INNER_VLAN_ID:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_INNER_VLAN_ID;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_INNER_VLAN_PRI:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_INNER_VLAN_PRI;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_OUTER_VLAN_ID:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_OUTER_VLAN_ID;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_OUTER_VLAN_PRI:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_OUTER_VLAN_PRI;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_ADD_VLAN_ID:
      return SAI_ACL_ENTRY_ATTR_ACTION_ADD_VLAN_ID;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_ADD_VLAN_PRI:
      return SAI_ACL_ENTRY_ATTR_ACTION_ADD_VLAN_PRI;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_SRC_MAC:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_SRC_MAC;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_DST_MAC:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_DST_MAC;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_SRC_IP:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_SRC_IP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_DST_IP:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_DST_IP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_SRC_IPV6:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_SRC_IPV6;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_DST_IPV6:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_DST_IPV6;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_DSCP:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_DSCP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_ECN:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_ECN;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_L4_SRC_PORT:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_L4_SRC_PORT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_L4_DST_PORT:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_L4_DST_PORT;

    case lemming::dataplane::sai::
        ACL_ENTRY_ATTR_ACTION_INGRESS_SAMPLEPACKET_ENABLE:
      return SAI_ACL_ENTRY_ATTR_ACTION_INGRESS_SAMPLEPACKET_ENABLE;

    case lemming::dataplane::sai::
        ACL_ENTRY_ATTR_ACTION_EGRESS_SAMPLEPACKET_ENABLE:
      return SAI_ACL_ENTRY_ATTR_ACTION_EGRESS_SAMPLEPACKET_ENABLE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_ACL_META_DATA:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_ACL_META_DATA;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_EGRESS_BLOCK_PORT_LIST:
      return SAI_ACL_ENTRY_ATTR_ACTION_EGRESS_BLOCK_PORT_LIST;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_USER_TRAP_ID:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_USER_TRAP_ID;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_DO_NOT_LEARN:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_DO_NOT_LEARN;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_ACL_DTEL_FLOW_OP:
      return SAI_ACL_ENTRY_ATTR_ACTION_ACL_DTEL_FLOW_OP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_DTEL_INT_SESSION:
      return SAI_ACL_ENTRY_ATTR_ACTION_DTEL_INT_SESSION;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_DTEL_DROP_REPORT_ENABLE:
      return SAI_ACL_ENTRY_ATTR_ACTION_DTEL_DROP_REPORT_ENABLE;

    case lemming::dataplane::sai::
        ACL_ENTRY_ATTR_ACTION_DTEL_TAIL_DROP_REPORT_ENABLE:
      return SAI_ACL_ENTRY_ATTR_ACTION_DTEL_TAIL_DROP_REPORT_ENABLE;

    case lemming::dataplane::sai::
        ACL_ENTRY_ATTR_ACTION_DTEL_FLOW_SAMPLE_PERCENT:
      return SAI_ACL_ENTRY_ATTR_ACTION_DTEL_FLOW_SAMPLE_PERCENT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_DTEL_REPORT_ALL_PACKETS:
      return SAI_ACL_ENTRY_ATTR_ACTION_DTEL_REPORT_ALL_PACKETS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_NO_NAT:
      return SAI_ACL_ENTRY_ATTR_ACTION_NO_NAT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_INT_INSERT:
      return SAI_ACL_ENTRY_ATTR_ACTION_INT_INSERT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_INT_DELETE:
      return SAI_ACL_ENTRY_ATTR_ACTION_INT_DELETE;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_INT_REPORT_FLOW:
      return SAI_ACL_ENTRY_ATTR_ACTION_INT_REPORT_FLOW;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_INT_REPORT_DROPS:
      return SAI_ACL_ENTRY_ATTR_ACTION_INT_REPORT_DROPS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_INT_REPORT_TAIL_DROPS:
      return SAI_ACL_ENTRY_ATTR_ACTION_INT_REPORT_TAIL_DROPS;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_TAM_INT_OBJECT:
      return SAI_ACL_ENTRY_ATTR_ACTION_TAM_INT_OBJECT;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_ISOLATION_GROUP:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_ISOLATION_GROUP;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_MACSEC_FLOW:
      return SAI_ACL_ENTRY_ATTR_ACTION_MACSEC_FLOW;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_LAG_HASH_ID:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_LAG_HASH_ID;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_ECMP_HASH_ID:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_ECMP_HASH_ID;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_VRF:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_VRF;

    case lemming::dataplane::sai::ACL_ENTRY_ATTR_ACTION_SET_FORWARDING_CLASS:
      return SAI_ACL_ENTRY_ATTR_ACTION_SET_FORWARDING_CLASS;

    default:
      return SAI_ACL_ENTRY_ATTR_TABLE_ID;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_acl_entry_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::AclEntryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclIpFrag convert_sai_acl_ip_frag_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_IP_FRAG_ANY:
      return lemming::dataplane::sai::ACL_IP_FRAG_ANY;

    case SAI_ACL_IP_FRAG_NON_FRAG:
      return lemming::dataplane::sai::ACL_IP_FRAG_NON_FRAG;

    case SAI_ACL_IP_FRAG_NON_FRAG_OR_HEAD:
      return lemming::dataplane::sai::ACL_IP_FRAG_NON_FRAG_OR_HEAD;

    case SAI_ACL_IP_FRAG_HEAD:
      return lemming::dataplane::sai::ACL_IP_FRAG_HEAD;

    case SAI_ACL_IP_FRAG_NON_HEAD:
      return lemming::dataplane::sai::ACL_IP_FRAG_NON_HEAD;

    default:
      return lemming::dataplane::sai::ACL_IP_FRAG_UNSPECIFIED;
  }
}
sai_acl_ip_frag_t convert_sai_acl_ip_frag_t_to_sai(
    lemming::dataplane::sai::AclIpFrag val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_IP_FRAG_ANY:
      return SAI_ACL_IP_FRAG_ANY;

    case lemming::dataplane::sai::ACL_IP_FRAG_NON_FRAG:
      return SAI_ACL_IP_FRAG_NON_FRAG;

    case lemming::dataplane::sai::ACL_IP_FRAG_NON_FRAG_OR_HEAD:
      return SAI_ACL_IP_FRAG_NON_FRAG_OR_HEAD;

    case lemming::dataplane::sai::ACL_IP_FRAG_HEAD:
      return SAI_ACL_IP_FRAG_HEAD;

    case lemming::dataplane::sai::ACL_IP_FRAG_NON_HEAD:
      return SAI_ACL_IP_FRAG_NON_HEAD;

    default:
      return SAI_ACL_IP_FRAG_ANY;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_acl_ip_frag_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_ip_frag_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_ip_frag_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_ip_frag_t_to_sai(
        static_cast<lemming::dataplane::sai::AclIpFrag>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclIpType convert_sai_acl_ip_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_IP_TYPE_ANY:
      return lemming::dataplane::sai::ACL_IP_TYPE_ANY;

    case SAI_ACL_IP_TYPE_IP:
      return lemming::dataplane::sai::ACL_IP_TYPE_IP;

    case SAI_ACL_IP_TYPE_NON_IP:
      return lemming::dataplane::sai::ACL_IP_TYPE_NON_IP;

    case SAI_ACL_IP_TYPE_IPV4ANY:
      return lemming::dataplane::sai::ACL_IP_TYPE_IPV4ANY;

    case SAI_ACL_IP_TYPE_NON_IPV4:
      return lemming::dataplane::sai::ACL_IP_TYPE_NON_IPV4;

    case SAI_ACL_IP_TYPE_IPV6ANY:
      return lemming::dataplane::sai::ACL_IP_TYPE_IPV6ANY;

    case SAI_ACL_IP_TYPE_NON_IPV6:
      return lemming::dataplane::sai::ACL_IP_TYPE_NON_IPV6;

    case SAI_ACL_IP_TYPE_ARP:
      return lemming::dataplane::sai::ACL_IP_TYPE_ARP;

    case SAI_ACL_IP_TYPE_ARP_REQUEST:
      return lemming::dataplane::sai::ACL_IP_TYPE_ARP_REQUEST;

    case SAI_ACL_IP_TYPE_ARP_REPLY:
      return lemming::dataplane::sai::ACL_IP_TYPE_ARP_REPLY;

    default:
      return lemming::dataplane::sai::ACL_IP_TYPE_UNSPECIFIED;
  }
}
sai_acl_ip_type_t convert_sai_acl_ip_type_t_to_sai(
    lemming::dataplane::sai::AclIpType val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_IP_TYPE_ANY:
      return SAI_ACL_IP_TYPE_ANY;

    case lemming::dataplane::sai::ACL_IP_TYPE_IP:
      return SAI_ACL_IP_TYPE_IP;

    case lemming::dataplane::sai::ACL_IP_TYPE_NON_IP:
      return SAI_ACL_IP_TYPE_NON_IP;

    case lemming::dataplane::sai::ACL_IP_TYPE_IPV4ANY:
      return SAI_ACL_IP_TYPE_IPV4ANY;

    case lemming::dataplane::sai::ACL_IP_TYPE_NON_IPV4:
      return SAI_ACL_IP_TYPE_NON_IPV4;

    case lemming::dataplane::sai::ACL_IP_TYPE_IPV6ANY:
      return SAI_ACL_IP_TYPE_IPV6ANY;

    case lemming::dataplane::sai::ACL_IP_TYPE_NON_IPV6:
      return SAI_ACL_IP_TYPE_NON_IPV6;

    case lemming::dataplane::sai::ACL_IP_TYPE_ARP:
      return SAI_ACL_IP_TYPE_ARP;

    case lemming::dataplane::sai::ACL_IP_TYPE_ARP_REQUEST:
      return SAI_ACL_IP_TYPE_ARP_REQUEST;

    case lemming::dataplane::sai::ACL_IP_TYPE_ARP_REPLY:
      return SAI_ACL_IP_TYPE_ARP_REPLY;

    default:
      return SAI_ACL_IP_TYPE_ANY;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_acl_ip_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_ip_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_ip_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_ip_type_t_to_sai(
        static_cast<lemming::dataplane::sai::AclIpType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclRangeAttr convert_sai_acl_range_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_RANGE_ATTR_TYPE:
      return lemming::dataplane::sai::ACL_RANGE_ATTR_TYPE;

    case SAI_ACL_RANGE_ATTR_LIMIT:
      return lemming::dataplane::sai::ACL_RANGE_ATTR_LIMIT;

    default:
      return lemming::dataplane::sai::ACL_RANGE_ATTR_UNSPECIFIED;
  }
}
sai_acl_range_attr_t convert_sai_acl_range_attr_t_to_sai(
    lemming::dataplane::sai::AclRangeAttr val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_RANGE_ATTR_TYPE:
      return SAI_ACL_RANGE_ATTR_TYPE;

    case lemming::dataplane::sai::ACL_RANGE_ATTR_LIMIT:
      return SAI_ACL_RANGE_ATTR_LIMIT;

    default:
      return SAI_ACL_RANGE_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_acl_range_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_range_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_range_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_range_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::AclRangeAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclRangeType convert_sai_acl_range_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_RANGE_TYPE_L4_SRC_PORT_RANGE:
      return lemming::dataplane::sai::ACL_RANGE_TYPE_L4_SRC_PORT_RANGE;

    case SAI_ACL_RANGE_TYPE_L4_DST_PORT_RANGE:
      return lemming::dataplane::sai::ACL_RANGE_TYPE_L4_DST_PORT_RANGE;

    case SAI_ACL_RANGE_TYPE_OUTER_VLAN:
      return lemming::dataplane::sai::ACL_RANGE_TYPE_OUTER_VLAN;

    case SAI_ACL_RANGE_TYPE_INNER_VLAN:
      return lemming::dataplane::sai::ACL_RANGE_TYPE_INNER_VLAN;

    case SAI_ACL_RANGE_TYPE_PACKET_LENGTH:
      return lemming::dataplane::sai::ACL_RANGE_TYPE_PACKET_LENGTH;

    default:
      return lemming::dataplane::sai::ACL_RANGE_TYPE_UNSPECIFIED;
  }
}
sai_acl_range_type_t convert_sai_acl_range_type_t_to_sai(
    lemming::dataplane::sai::AclRangeType val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_RANGE_TYPE_L4_SRC_PORT_RANGE:
      return SAI_ACL_RANGE_TYPE_L4_SRC_PORT_RANGE;

    case lemming::dataplane::sai::ACL_RANGE_TYPE_L4_DST_PORT_RANGE:
      return SAI_ACL_RANGE_TYPE_L4_DST_PORT_RANGE;

    case lemming::dataplane::sai::ACL_RANGE_TYPE_OUTER_VLAN:
      return SAI_ACL_RANGE_TYPE_OUTER_VLAN;

    case lemming::dataplane::sai::ACL_RANGE_TYPE_INNER_VLAN:
      return SAI_ACL_RANGE_TYPE_INNER_VLAN;

    case lemming::dataplane::sai::ACL_RANGE_TYPE_PACKET_LENGTH:
      return SAI_ACL_RANGE_TYPE_PACKET_LENGTH;

    default:
      return SAI_ACL_RANGE_TYPE_L4_SRC_PORT_RANGE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_acl_range_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_range_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_range_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_range_type_t_to_sai(
        static_cast<lemming::dataplane::sai::AclRangeType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclStage convert_sai_acl_stage_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_STAGE_INGRESS:
      return lemming::dataplane::sai::ACL_STAGE_INGRESS;

    case SAI_ACL_STAGE_EGRESS:
      return lemming::dataplane::sai::ACL_STAGE_EGRESS;

    case SAI_ACL_STAGE_INGRESS_MACSEC:
      return lemming::dataplane::sai::ACL_STAGE_INGRESS_MACSEC;

    case SAI_ACL_STAGE_EGRESS_MACSEC:
      return lemming::dataplane::sai::ACL_STAGE_EGRESS_MACSEC;

    case SAI_ACL_STAGE_PRE_INGRESS:
      return lemming::dataplane::sai::ACL_STAGE_PRE_INGRESS;

    default:
      return lemming::dataplane::sai::ACL_STAGE_UNSPECIFIED;
  }
}
sai_acl_stage_t convert_sai_acl_stage_t_to_sai(
    lemming::dataplane::sai::AclStage val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_STAGE_INGRESS:
      return SAI_ACL_STAGE_INGRESS;

    case lemming::dataplane::sai::ACL_STAGE_EGRESS:
      return SAI_ACL_STAGE_EGRESS;

    case lemming::dataplane::sai::ACL_STAGE_INGRESS_MACSEC:
      return SAI_ACL_STAGE_INGRESS_MACSEC;

    case lemming::dataplane::sai::ACL_STAGE_EGRESS_MACSEC:
      return SAI_ACL_STAGE_EGRESS_MACSEC;

    case lemming::dataplane::sai::ACL_STAGE_PRE_INGRESS:
      return SAI_ACL_STAGE_PRE_INGRESS;

    default:
      return SAI_ACL_STAGE_INGRESS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_acl_stage_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_stage_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_stage_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_stage_t_to_sai(
        static_cast<lemming::dataplane::sai::AclStage>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclTableAttr convert_sai_acl_table_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_TABLE_ATTR_ACL_STAGE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_ACL_STAGE;

    case SAI_ACL_TABLE_ATTR_ACL_BIND_POINT_TYPE_LIST:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_ACL_BIND_POINT_TYPE_LIST;

    case SAI_ACL_TABLE_ATTR_SIZE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_SIZE;

    case SAI_ACL_TABLE_ATTR_ACL_ACTION_TYPE_LIST:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_ACL_ACTION_TYPE_LIST;

    case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IPV6;

    case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD3:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD3;

    case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD2:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD2;

    case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD1:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD1;

    case SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD0:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD0;

    case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IPV6;

    case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD3:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD3;

    case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD2:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD2;

    case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD1:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD1;

    case SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD0:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD0;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_IPV6:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_SRC_IPV6;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_IPV6:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_DST_IPV6;

    case SAI_ACL_TABLE_ATTR_FIELD_SRC_MAC:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_MAC;

    case SAI_ACL_TABLE_ATTR_FIELD_DST_MAC:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_MAC;

    case SAI_ACL_TABLE_ATTR_FIELD_SRC_IP:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IP;

    case SAI_ACL_TABLE_ATTR_FIELD_DST_IP:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IP;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_IP:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_SRC_IP;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_IP:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_DST_IP;

    case SAI_ACL_TABLE_ATTR_FIELD_IN_PORTS:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IN_PORTS;

    case SAI_ACL_TABLE_ATTR_FIELD_OUT_PORTS:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_OUT_PORTS;

    case SAI_ACL_TABLE_ATTR_FIELD_IN_PORT:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IN_PORT;

    case SAI_ACL_TABLE_ATTR_FIELD_OUT_PORT:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_OUT_PORT;

    case SAI_ACL_TABLE_ATTR_FIELD_SRC_PORT:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_PORT;

    case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_ID:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_OUTER_VLAN_ID;

    case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_PRI:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_OUTER_VLAN_PRI;

    case SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_CFI:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_OUTER_VLAN_CFI;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_ID:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_VLAN_ID;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_PRI:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_VLAN_PRI;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_CFI:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_VLAN_CFI;

    case SAI_ACL_TABLE_ATTR_FIELD_L4_SRC_PORT:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_L4_SRC_PORT;

    case SAI_ACL_TABLE_ATTR_FIELD_L4_DST_PORT:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_L4_DST_PORT;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_L4_SRC_PORT:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_L4_SRC_PORT;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_L4_DST_PORT:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_L4_DST_PORT;

    case SAI_ACL_TABLE_ATTR_FIELD_ETHER_TYPE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ETHER_TYPE;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_ETHER_TYPE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_ETHER_TYPE;

    case SAI_ACL_TABLE_ATTR_FIELD_IP_PROTOCOL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IP_PROTOCOL;

    case SAI_ACL_TABLE_ATTR_FIELD_INNER_IP_PROTOCOL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_IP_PROTOCOL;

    case SAI_ACL_TABLE_ATTR_FIELD_IP_IDENTIFICATION:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IP_IDENTIFICATION;

    case SAI_ACL_TABLE_ATTR_FIELD_DSCP:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DSCP;

    case SAI_ACL_TABLE_ATTR_FIELD_ECN:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ECN;

    case SAI_ACL_TABLE_ATTR_FIELD_TTL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TTL;

    case SAI_ACL_TABLE_ATTR_FIELD_TOS:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TOS;

    case SAI_ACL_TABLE_ATTR_FIELD_IP_FLAGS:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IP_FLAGS;

    case SAI_ACL_TABLE_ATTR_FIELD_TCP_FLAGS:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TCP_FLAGS;

    case SAI_ACL_TABLE_ATTR_FIELD_ACL_IP_TYPE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ACL_IP_TYPE;

    case SAI_ACL_TABLE_ATTR_FIELD_ACL_IP_FRAG:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ACL_IP_FRAG;

    case SAI_ACL_TABLE_ATTR_FIELD_IPV6_FLOW_LABEL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IPV6_FLOW_LABEL;

    case SAI_ACL_TABLE_ATTR_FIELD_TC:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TC;

    case SAI_ACL_TABLE_ATTR_FIELD_ICMP_TYPE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ICMP_TYPE;

    case SAI_ACL_TABLE_ATTR_FIELD_ICMP_CODE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ICMP_CODE;

    case SAI_ACL_TABLE_ATTR_FIELD_ICMPV6_TYPE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ICMPV6_TYPE;

    case SAI_ACL_TABLE_ATTR_FIELD_ICMPV6_CODE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ICMPV6_CODE;

    case SAI_ACL_TABLE_ATTR_FIELD_PACKET_VLAN:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_PACKET_VLAN;

    case SAI_ACL_TABLE_ATTR_FIELD_TUNNEL_VNI:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TUNNEL_VNI;

    case SAI_ACL_TABLE_ATTR_FIELD_HAS_VLAN_TAG:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_HAS_VLAN_TAG;

    case SAI_ACL_TABLE_ATTR_FIELD_MACSEC_SCI:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MACSEC_SCI;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_LABEL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_LABEL;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_TTL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_TTL;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_EXP:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_EXP;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_BOS:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_BOS;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_LABEL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_LABEL;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_TTL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_TTL;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_EXP:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_EXP;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_BOS:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_BOS;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_LABEL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_LABEL;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_TTL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_TTL;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_EXP:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_EXP;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_BOS:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_BOS;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_LABEL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_LABEL;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_TTL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_TTL;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_EXP:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_EXP;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_BOS:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_BOS;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_LABEL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_LABEL;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_TTL:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_TTL;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_EXP:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_EXP;

    case SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_BOS:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_BOS;

    case SAI_ACL_TABLE_ATTR_FIELD_FDB_DST_USER_META:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_FDB_DST_USER_META;

    case SAI_ACL_TABLE_ATTR_FIELD_ROUTE_DST_USER_META:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ROUTE_DST_USER_META;

    case SAI_ACL_TABLE_ATTR_FIELD_NEIGHBOR_DST_USER_META:
      return lemming::dataplane::sai::
          ACL_TABLE_ATTR_FIELD_NEIGHBOR_DST_USER_META;

    case SAI_ACL_TABLE_ATTR_FIELD_PORT_USER_META:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_PORT_USER_META;

    case SAI_ACL_TABLE_ATTR_FIELD_VLAN_USER_META:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_VLAN_USER_META;

    case SAI_ACL_TABLE_ATTR_FIELD_ACL_USER_META:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ACL_USER_META;

    case SAI_ACL_TABLE_ATTR_FIELD_FDB_NPU_META_DST_HIT:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_FDB_NPU_META_DST_HIT;

    case SAI_ACL_TABLE_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT:
      return lemming::dataplane::sai::
          ACL_TABLE_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT;

    case SAI_ACL_TABLE_ATTR_FIELD_ROUTE_NPU_META_DST_HIT:
      return lemming::dataplane::sai::
          ACL_TABLE_ATTR_FIELD_ROUTE_NPU_META_DST_HIT;

    case SAI_ACL_TABLE_ATTR_FIELD_BTH_OPCODE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_BTH_OPCODE;

    case SAI_ACL_TABLE_ATTR_FIELD_AETH_SYNDROME:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_AETH_SYNDROME;

    case SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN:
      return lemming::dataplane::sai::
          ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN;

    case SAI_ACL_TABLE_ATTR_FIELD_ACL_RANGE_TYPE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ACL_RANGE_TYPE;

    case SAI_ACL_TABLE_ATTR_FIELD_IPV6_NEXT_HEADER:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IPV6_NEXT_HEADER;

    case SAI_ACL_TABLE_ATTR_FIELD_GRE_KEY:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_GRE_KEY;

    case SAI_ACL_TABLE_ATTR_FIELD_TAM_INT_TYPE:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TAM_INT_TYPE;

    case SAI_ACL_TABLE_ATTR_ENTRY_LIST:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_ENTRY_LIST;

    case SAI_ACL_TABLE_ATTR_AVAILABLE_ACL_ENTRY:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_AVAILABLE_ACL_ENTRY;

    case SAI_ACL_TABLE_ATTR_AVAILABLE_ACL_COUNTER:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_AVAILABLE_ACL_COUNTER;

    default:
      return lemming::dataplane::sai::ACL_TABLE_ATTR_UNSPECIFIED;
  }
}
sai_acl_table_attr_t convert_sai_acl_table_attr_t_to_sai(
    lemming::dataplane::sai::AclTableAttr val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_TABLE_ATTR_ACL_STAGE:
      return SAI_ACL_TABLE_ATTR_ACL_STAGE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_ACL_BIND_POINT_TYPE_LIST:
      return SAI_ACL_TABLE_ATTR_ACL_BIND_POINT_TYPE_LIST;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_SIZE:
      return SAI_ACL_TABLE_ATTR_SIZE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_ACL_ACTION_TYPE_LIST:
      return SAI_ACL_TABLE_ATTR_ACL_ACTION_TYPE_LIST;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IPV6:
      return SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD3:
      return SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD3;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD2:
      return SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD2;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD1:
      return SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD1;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD0:
      return SAI_ACL_TABLE_ATTR_FIELD_SRC_IPV6_WORD0;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IPV6:
      return SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD3:
      return SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD3;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD2:
      return SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD2;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD1:
      return SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD1;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD0:
      return SAI_ACL_TABLE_ATTR_FIELD_DST_IPV6_WORD0;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_SRC_IPV6:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_IPV6;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_DST_IPV6:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_IPV6;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_MAC:
      return SAI_ACL_TABLE_ATTR_FIELD_SRC_MAC;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_MAC:
      return SAI_ACL_TABLE_ATTR_FIELD_DST_MAC;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_IP:
      return SAI_ACL_TABLE_ATTR_FIELD_SRC_IP;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DST_IP:
      return SAI_ACL_TABLE_ATTR_FIELD_DST_IP;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_SRC_IP:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_SRC_IP;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_DST_IP:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_DST_IP;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IN_PORTS:
      return SAI_ACL_TABLE_ATTR_FIELD_IN_PORTS;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_OUT_PORTS:
      return SAI_ACL_TABLE_ATTR_FIELD_OUT_PORTS;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IN_PORT:
      return SAI_ACL_TABLE_ATTR_FIELD_IN_PORT;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_OUT_PORT:
      return SAI_ACL_TABLE_ATTR_FIELD_OUT_PORT;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_SRC_PORT:
      return SAI_ACL_TABLE_ATTR_FIELD_SRC_PORT;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_OUTER_VLAN_ID:
      return SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_ID;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_OUTER_VLAN_PRI:
      return SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_PRI;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_OUTER_VLAN_CFI:
      return SAI_ACL_TABLE_ATTR_FIELD_OUTER_VLAN_CFI;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_VLAN_ID:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_ID;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_VLAN_PRI:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_PRI;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_VLAN_CFI:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_VLAN_CFI;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_L4_SRC_PORT:
      return SAI_ACL_TABLE_ATTR_FIELD_L4_SRC_PORT;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_L4_DST_PORT:
      return SAI_ACL_TABLE_ATTR_FIELD_L4_DST_PORT;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_L4_SRC_PORT:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_L4_SRC_PORT;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_L4_DST_PORT:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_L4_DST_PORT;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ETHER_TYPE:
      return SAI_ACL_TABLE_ATTR_FIELD_ETHER_TYPE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_ETHER_TYPE:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_ETHER_TYPE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IP_PROTOCOL:
      return SAI_ACL_TABLE_ATTR_FIELD_IP_PROTOCOL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_INNER_IP_PROTOCOL:
      return SAI_ACL_TABLE_ATTR_FIELD_INNER_IP_PROTOCOL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IP_IDENTIFICATION:
      return SAI_ACL_TABLE_ATTR_FIELD_IP_IDENTIFICATION;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_DSCP:
      return SAI_ACL_TABLE_ATTR_FIELD_DSCP;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ECN:
      return SAI_ACL_TABLE_ATTR_FIELD_ECN;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TTL:
      return SAI_ACL_TABLE_ATTR_FIELD_TTL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TOS:
      return SAI_ACL_TABLE_ATTR_FIELD_TOS;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IP_FLAGS:
      return SAI_ACL_TABLE_ATTR_FIELD_IP_FLAGS;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TCP_FLAGS:
      return SAI_ACL_TABLE_ATTR_FIELD_TCP_FLAGS;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ACL_IP_TYPE:
      return SAI_ACL_TABLE_ATTR_FIELD_ACL_IP_TYPE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ACL_IP_FRAG:
      return SAI_ACL_TABLE_ATTR_FIELD_ACL_IP_FRAG;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IPV6_FLOW_LABEL:
      return SAI_ACL_TABLE_ATTR_FIELD_IPV6_FLOW_LABEL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TC:
      return SAI_ACL_TABLE_ATTR_FIELD_TC;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ICMP_TYPE:
      return SAI_ACL_TABLE_ATTR_FIELD_ICMP_TYPE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ICMP_CODE:
      return SAI_ACL_TABLE_ATTR_FIELD_ICMP_CODE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ICMPV6_TYPE:
      return SAI_ACL_TABLE_ATTR_FIELD_ICMPV6_TYPE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ICMPV6_CODE:
      return SAI_ACL_TABLE_ATTR_FIELD_ICMPV6_CODE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_PACKET_VLAN:
      return SAI_ACL_TABLE_ATTR_FIELD_PACKET_VLAN;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TUNNEL_VNI:
      return SAI_ACL_TABLE_ATTR_FIELD_TUNNEL_VNI;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_HAS_VLAN_TAG:
      return SAI_ACL_TABLE_ATTR_FIELD_HAS_VLAN_TAG;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MACSEC_SCI:
      return SAI_ACL_TABLE_ATTR_FIELD_MACSEC_SCI;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_LABEL:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_LABEL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_TTL:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_TTL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_EXP:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_EXP;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_BOS:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL0_BOS;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_LABEL:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_LABEL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_TTL:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_TTL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_EXP:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_EXP;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_BOS:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL1_BOS;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_LABEL:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_LABEL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_TTL:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_TTL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_EXP:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_EXP;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_BOS:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL2_BOS;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_LABEL:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_LABEL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_TTL:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_TTL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_EXP:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_EXP;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_BOS:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL3_BOS;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_LABEL:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_LABEL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_TTL:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_TTL;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_EXP:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_EXP;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_BOS:
      return SAI_ACL_TABLE_ATTR_FIELD_MPLS_LABEL4_BOS;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_FDB_DST_USER_META:
      return SAI_ACL_TABLE_ATTR_FIELD_FDB_DST_USER_META;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ROUTE_DST_USER_META:
      return SAI_ACL_TABLE_ATTR_FIELD_ROUTE_DST_USER_META;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_NEIGHBOR_DST_USER_META:
      return SAI_ACL_TABLE_ATTR_FIELD_NEIGHBOR_DST_USER_META;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_PORT_USER_META:
      return SAI_ACL_TABLE_ATTR_FIELD_PORT_USER_META;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_VLAN_USER_META:
      return SAI_ACL_TABLE_ATTR_FIELD_VLAN_USER_META;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ACL_USER_META:
      return SAI_ACL_TABLE_ATTR_FIELD_ACL_USER_META;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_FDB_NPU_META_DST_HIT:
      return SAI_ACL_TABLE_ATTR_FIELD_FDB_NPU_META_DST_HIT;

    case lemming::dataplane::sai::
        ACL_TABLE_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT:
      return SAI_ACL_TABLE_ATTR_FIELD_NEIGHBOR_NPU_META_DST_HIT;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ROUTE_NPU_META_DST_HIT:
      return SAI_ACL_TABLE_ATTR_FIELD_ROUTE_NPU_META_DST_HIT;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_BTH_OPCODE:
      return SAI_ACL_TABLE_ATTR_FIELD_BTH_OPCODE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_AETH_SYNDROME:
      return SAI_ACL_TABLE_ATTR_FIELD_AETH_SYNDROME;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN:
      return SAI_ACL_TABLE_ATTR_USER_DEFINED_FIELD_GROUP_MIN;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_ACL_RANGE_TYPE:
      return SAI_ACL_TABLE_ATTR_FIELD_ACL_RANGE_TYPE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_IPV6_NEXT_HEADER:
      return SAI_ACL_TABLE_ATTR_FIELD_IPV6_NEXT_HEADER;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_GRE_KEY:
      return SAI_ACL_TABLE_ATTR_FIELD_GRE_KEY;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_FIELD_TAM_INT_TYPE:
      return SAI_ACL_TABLE_ATTR_FIELD_TAM_INT_TYPE;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_ENTRY_LIST:
      return SAI_ACL_TABLE_ATTR_ENTRY_LIST;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_AVAILABLE_ACL_ENTRY:
      return SAI_ACL_TABLE_ATTR_AVAILABLE_ACL_ENTRY;

    case lemming::dataplane::sai::ACL_TABLE_ATTR_AVAILABLE_ACL_COUNTER:
      return SAI_ACL_TABLE_ATTR_AVAILABLE_ACL_COUNTER;

    default:
      return SAI_ACL_TABLE_ATTR_ACL_STAGE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_acl_table_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_table_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_table_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_table_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::AclTableAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclTableGroupAttr
convert_sai_acl_table_group_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_TABLE_GROUP_ATTR_ACL_STAGE:
      return lemming::dataplane::sai::ACL_TABLE_GROUP_ATTR_ACL_STAGE;

    case SAI_ACL_TABLE_GROUP_ATTR_ACL_BIND_POINT_TYPE_LIST:
      return lemming::dataplane::sai::
          ACL_TABLE_GROUP_ATTR_ACL_BIND_POINT_TYPE_LIST;

    case SAI_ACL_TABLE_GROUP_ATTR_TYPE:
      return lemming::dataplane::sai::ACL_TABLE_GROUP_ATTR_TYPE;

    case SAI_ACL_TABLE_GROUP_ATTR_MEMBER_LIST:
      return lemming::dataplane::sai::ACL_TABLE_GROUP_ATTR_MEMBER_LIST;

    default:
      return lemming::dataplane::sai::ACL_TABLE_GROUP_ATTR_UNSPECIFIED;
  }
}
sai_acl_table_group_attr_t convert_sai_acl_table_group_attr_t_to_sai(
    lemming::dataplane::sai::AclTableGroupAttr val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_TABLE_GROUP_ATTR_ACL_STAGE:
      return SAI_ACL_TABLE_GROUP_ATTR_ACL_STAGE;

    case lemming::dataplane::sai::ACL_TABLE_GROUP_ATTR_ACL_BIND_POINT_TYPE_LIST:
      return SAI_ACL_TABLE_GROUP_ATTR_ACL_BIND_POINT_TYPE_LIST;

    case lemming::dataplane::sai::ACL_TABLE_GROUP_ATTR_TYPE:
      return SAI_ACL_TABLE_GROUP_ATTR_TYPE;

    case lemming::dataplane::sai::ACL_TABLE_GROUP_ATTR_MEMBER_LIST:
      return SAI_ACL_TABLE_GROUP_ATTR_MEMBER_LIST;

    default:
      return SAI_ACL_TABLE_GROUP_ATTR_ACL_STAGE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_group_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_table_group_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_table_group_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_table_group_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::AclTableGroupAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclTableGroupMemberAttr
convert_sai_acl_table_group_member_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_GROUP_ID:
      return lemming::dataplane::sai::
          ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_GROUP_ID;

    case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_ID:
      return lemming::dataplane::sai::ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_ID;

    case SAI_ACL_TABLE_GROUP_MEMBER_ATTR_PRIORITY:
      return lemming::dataplane::sai::ACL_TABLE_GROUP_MEMBER_ATTR_PRIORITY;

    default:
      return lemming::dataplane::sai::ACL_TABLE_GROUP_MEMBER_ATTR_UNSPECIFIED;
  }
}
sai_acl_table_group_member_attr_t
convert_sai_acl_table_group_member_attr_t_to_sai(
    lemming::dataplane::sai::AclTableGroupMemberAttr val) {
  switch (val) {
    case lemming::dataplane::sai::
        ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_GROUP_ID:
      return SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_GROUP_ID;

    case lemming::dataplane::sai::ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_ID:
      return SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_ID;

    case lemming::dataplane::sai::ACL_TABLE_GROUP_MEMBER_ATTR_PRIORITY:
      return SAI_ACL_TABLE_GROUP_MEMBER_ATTR_PRIORITY;

    default:
      return SAI_ACL_TABLE_GROUP_MEMBER_ATTR_ACL_TABLE_GROUP_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_group_member_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_acl_table_group_member_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_table_group_member_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_table_group_member_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::AclTableGroupMemberAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::AclTableGroupType
convert_sai_acl_table_group_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ACL_TABLE_GROUP_TYPE_SEQUENTIAL:
      return lemming::dataplane::sai::ACL_TABLE_GROUP_TYPE_SEQUENTIAL;

    case SAI_ACL_TABLE_GROUP_TYPE_PARALLEL:
      return lemming::dataplane::sai::ACL_TABLE_GROUP_TYPE_PARALLEL;

    default:
      return lemming::dataplane::sai::ACL_TABLE_GROUP_TYPE_UNSPECIFIED;
  }
}
sai_acl_table_group_type_t convert_sai_acl_table_group_type_t_to_sai(
    lemming::dataplane::sai::AclTableGroupType val) {
  switch (val) {
    case lemming::dataplane::sai::ACL_TABLE_GROUP_TYPE_SEQUENTIAL:
      return SAI_ACL_TABLE_GROUP_TYPE_SEQUENTIAL;

    case lemming::dataplane::sai::ACL_TABLE_GROUP_TYPE_PARALLEL:
      return SAI_ACL_TABLE_GROUP_TYPE_PARALLEL;

    default:
      return SAI_ACL_TABLE_GROUP_TYPE_SEQUENTIAL;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_acl_table_group_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_acl_table_group_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_acl_table_group_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_acl_table_group_type_t_to_sai(
        static_cast<lemming::dataplane::sai::AclTableGroupType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::ApiExtensions convert_sai_api_extensions_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_API_EXTENSIONS_RANGE_START:
      return lemming::dataplane::sai::API_EXTENSIONS_RANGE_START;

    case SAI_API_EXTENSIONS_RANGE_END:
      return lemming::dataplane::sai::API_EXTENSIONS_RANGE_END;

    default:
      return lemming::dataplane::sai::API_EXTENSIONS_UNSPECIFIED;
  }
}
sai_api_extensions_t convert_sai_api_extensions_t_to_sai(
    lemming::dataplane::sai::ApiExtensions val) {
  switch (val) {
    case lemming::dataplane::sai::API_EXTENSIONS_RANGE_START:
      return SAI_API_EXTENSIONS_RANGE_START;

    case lemming::dataplane::sai::API_EXTENSIONS_RANGE_END:
      return SAI_API_EXTENSIONS_RANGE_END;

    default:
      return SAI_API_EXTENSIONS_RANGE_START;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_api_extensions_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_api_extensions_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_api_extensions_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_api_extensions_t_to_sai(
        static_cast<lemming::dataplane::sai::ApiExtensions>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::Api convert_sai_api_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_API_UNSPECIFIED:
      return lemming::dataplane::sai::API_SAI_UNSPECIFIED;

    case SAI_API_SWITCH:
      return lemming::dataplane::sai::API_SWITCH;

    case SAI_API_PORT:
      return lemming::dataplane::sai::API_PORT;

    case SAI_API_FDB:
      return lemming::dataplane::sai::API_FDB;

    case SAI_API_VLAN:
      return lemming::dataplane::sai::API_VLAN;

    case SAI_API_VIRTUAL_ROUTER:
      return lemming::dataplane::sai::API_VIRTUAL_ROUTER;

    case SAI_API_ROUTE:
      return lemming::dataplane::sai::API_ROUTE;

    case SAI_API_NEXT_HOP:
      return lemming::dataplane::sai::API_NEXT_HOP;

    case SAI_API_NEXT_HOP_GROUP:
      return lemming::dataplane::sai::API_NEXT_HOP_GROUP;

    case SAI_API_ROUTER_INTERFACE:
      return lemming::dataplane::sai::API_ROUTER_INTERFACE;

    case SAI_API_NEIGHBOR:
      return lemming::dataplane::sai::API_NEIGHBOR;

    case SAI_API_ACL:
      return lemming::dataplane::sai::API_ACL;

    case SAI_API_HOSTIF:
      return lemming::dataplane::sai::API_HOSTIF;

    case SAI_API_MIRROR:
      return lemming::dataplane::sai::API_MIRROR;

    case SAI_API_SAMPLEPACKET:
      return lemming::dataplane::sai::API_SAMPLEPACKET;

    case SAI_API_STP:
      return lemming::dataplane::sai::API_STP;

    case SAI_API_LAG:
      return lemming::dataplane::sai::API_LAG;

    case SAI_API_POLICER:
      return lemming::dataplane::sai::API_POLICER;

    case SAI_API_WRED:
      return lemming::dataplane::sai::API_WRED;

    case SAI_API_QOS_MAP:
      return lemming::dataplane::sai::API_QOS_MAP;

    case SAI_API_QUEUE:
      return lemming::dataplane::sai::API_QUEUE;

    case SAI_API_SCHEDULER:
      return lemming::dataplane::sai::API_SCHEDULER;

    case SAI_API_SCHEDULER_GROUP:
      return lemming::dataplane::sai::API_SCHEDULER_GROUP;

    case SAI_API_BUFFER:
      return lemming::dataplane::sai::API_BUFFER;

    case SAI_API_HASH:
      return lemming::dataplane::sai::API_HASH;

    case SAI_API_UDF:
      return lemming::dataplane::sai::API_UDF;

    case SAI_API_TUNNEL:
      return lemming::dataplane::sai::API_TUNNEL;

    case SAI_API_L2MC:
      return lemming::dataplane::sai::API_L2MC;

    case SAI_API_IPMC:
      return lemming::dataplane::sai::API_IPMC;

    case SAI_API_RPF_GROUP:
      return lemming::dataplane::sai::API_RPF_GROUP;

    case SAI_API_L2MC_GROUP:
      return lemming::dataplane::sai::API_L2MC_GROUP;

    case SAI_API_IPMC_GROUP:
      return lemming::dataplane::sai::API_IPMC_GROUP;

    case SAI_API_MCAST_FDB:
      return lemming::dataplane::sai::API_MCAST_FDB;

    case SAI_API_BRIDGE:
      return lemming::dataplane::sai::API_BRIDGE;

    case SAI_API_TAM:
      return lemming::dataplane::sai::API_TAM;

    case SAI_API_SRV6:
      return lemming::dataplane::sai::API_SRV6;

    case SAI_API_MPLS:
      return lemming::dataplane::sai::API_MPLS;

    case SAI_API_DTEL:
      return lemming::dataplane::sai::API_DTEL;

    case SAI_API_BFD:
      return lemming::dataplane::sai::API_BFD;

    case SAI_API_ISOLATION_GROUP:
      return lemming::dataplane::sai::API_ISOLATION_GROUP;

    case SAI_API_NAT:
      return lemming::dataplane::sai::API_NAT;

    case SAI_API_COUNTER:
      return lemming::dataplane::sai::API_COUNTER;

    case SAI_API_DEBUG_COUNTER:
      return lemming::dataplane::sai::API_DEBUG_COUNTER;

    case SAI_API_MACSEC:
      return lemming::dataplane::sai::API_MACSEC;

    case SAI_API_SYSTEM_PORT:
      return lemming::dataplane::sai::API_SYSTEM_PORT;

    case SAI_API_MY_MAC:
      return lemming::dataplane::sai::API_MY_MAC;

    case SAI_API_IPSEC:
      return lemming::dataplane::sai::API_IPSEC;

    case SAI_API_GENERIC_PROGRAMMABLE:
      return lemming::dataplane::sai::API_GENERIC_PROGRAMMABLE;

    case SAI_API_MAX:
      return lemming::dataplane::sai::API_MAX;

    default:
      return lemming::dataplane::sai::API_UNSPECIFIED;
  }
}
sai_api_t convert_sai_api_t_to_sai(lemming::dataplane::sai::Api val) {
  switch (val) {
    case lemming::dataplane::sai::API_SAI_UNSPECIFIED:
      return SAI_API_UNSPECIFIED;

    case lemming::dataplane::sai::API_SWITCH:
      return SAI_API_SWITCH;

    case lemming::dataplane::sai::API_PORT:
      return SAI_API_PORT;

    case lemming::dataplane::sai::API_FDB:
      return SAI_API_FDB;

    case lemming::dataplane::sai::API_VLAN:
      return SAI_API_VLAN;

    case lemming::dataplane::sai::API_VIRTUAL_ROUTER:
      return SAI_API_VIRTUAL_ROUTER;

    case lemming::dataplane::sai::API_ROUTE:
      return SAI_API_ROUTE;

    case lemming::dataplane::sai::API_NEXT_HOP:
      return SAI_API_NEXT_HOP;

    case lemming::dataplane::sai::API_NEXT_HOP_GROUP:
      return SAI_API_NEXT_HOP_GROUP;

    case lemming::dataplane::sai::API_ROUTER_INTERFACE:
      return SAI_API_ROUTER_INTERFACE;

    case lemming::dataplane::sai::API_NEIGHBOR:
      return SAI_API_NEIGHBOR;

    case lemming::dataplane::sai::API_ACL:
      return SAI_API_ACL;

    case lemming::dataplane::sai::API_HOSTIF:
      return SAI_API_HOSTIF;

    case lemming::dataplane::sai::API_MIRROR:
      return SAI_API_MIRROR;

    case lemming::dataplane::sai::API_SAMPLEPACKET:
      return SAI_API_SAMPLEPACKET;

    case lemming::dataplane::sai::API_STP:
      return SAI_API_STP;

    case lemming::dataplane::sai::API_LAG:
      return SAI_API_LAG;

    case lemming::dataplane::sai::API_POLICER:
      return SAI_API_POLICER;

    case lemming::dataplane::sai::API_WRED:
      return SAI_API_WRED;

    case lemming::dataplane::sai::API_QOS_MAP:
      return SAI_API_QOS_MAP;

    case lemming::dataplane::sai::API_QUEUE:
      return SAI_API_QUEUE;

    case lemming::dataplane::sai::API_SCHEDULER:
      return SAI_API_SCHEDULER;

    case lemming::dataplane::sai::API_SCHEDULER_GROUP:
      return SAI_API_SCHEDULER_GROUP;

    case lemming::dataplane::sai::API_BUFFER:
      return SAI_API_BUFFER;

    case lemming::dataplane::sai::API_HASH:
      return SAI_API_HASH;

    case lemming::dataplane::sai::API_UDF:
      return SAI_API_UDF;

    case lemming::dataplane::sai::API_TUNNEL:
      return SAI_API_TUNNEL;

    case lemming::dataplane::sai::API_L2MC:
      return SAI_API_L2MC;

    case lemming::dataplane::sai::API_IPMC:
      return SAI_API_IPMC;

    case lemming::dataplane::sai::API_RPF_GROUP:
      return SAI_API_RPF_GROUP;

    case lemming::dataplane::sai::API_L2MC_GROUP:
      return SAI_API_L2MC_GROUP;

    case lemming::dataplane::sai::API_IPMC_GROUP:
      return SAI_API_IPMC_GROUP;

    case lemming::dataplane::sai::API_MCAST_FDB:
      return SAI_API_MCAST_FDB;

    case lemming::dataplane::sai::API_BRIDGE:
      return SAI_API_BRIDGE;

    case lemming::dataplane::sai::API_TAM:
      return SAI_API_TAM;

    case lemming::dataplane::sai::API_SRV6:
      return SAI_API_SRV6;

    case lemming::dataplane::sai::API_MPLS:
      return SAI_API_MPLS;

    case lemming::dataplane::sai::API_DTEL:
      return SAI_API_DTEL;

    case lemming::dataplane::sai::API_BFD:
      return SAI_API_BFD;

    case lemming::dataplane::sai::API_ISOLATION_GROUP:
      return SAI_API_ISOLATION_GROUP;

    case lemming::dataplane::sai::API_NAT:
      return SAI_API_NAT;

    case lemming::dataplane::sai::API_COUNTER:
      return SAI_API_COUNTER;

    case lemming::dataplane::sai::API_DEBUG_COUNTER:
      return SAI_API_DEBUG_COUNTER;

    case lemming::dataplane::sai::API_MACSEC:
      return SAI_API_MACSEC;

    case lemming::dataplane::sai::API_SYSTEM_PORT:
      return SAI_API_SYSTEM_PORT;

    case lemming::dataplane::sai::API_MY_MAC:
      return SAI_API_MY_MAC;

    case lemming::dataplane::sai::API_IPSEC:
      return SAI_API_IPSEC;

    case lemming::dataplane::sai::API_GENERIC_PROGRAMMABLE:
      return SAI_API_GENERIC_PROGRAMMABLE;

    case lemming::dataplane::sai::API_MAX:
      return SAI_API_MAX;

    default:
      return SAI_API_UNSPECIFIED;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_api_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_api_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_api_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_api_t_to_sai(
        static_cast<lemming::dataplane::sai::Api>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BfdEncapsulationType
convert_sai_bfd_encapsulation_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_BFD_ENCAPSULATION_TYPE_IP_IN_IP:
      return lemming::dataplane::sai::BFD_ENCAPSULATION_TYPE_IP_IN_IP;

    case SAI_BFD_ENCAPSULATION_TYPE_L3_GRE_TUNNEL:
      return lemming::dataplane::sai::BFD_ENCAPSULATION_TYPE_L3_GRE_TUNNEL;

    case SAI_BFD_ENCAPSULATION_TYPE_NONE:
      return lemming::dataplane::sai::BFD_ENCAPSULATION_TYPE_NONE;

    default:
      return lemming::dataplane::sai::BFD_ENCAPSULATION_TYPE_UNSPECIFIED;
  }
}
sai_bfd_encapsulation_type_t convert_sai_bfd_encapsulation_type_t_to_sai(
    lemming::dataplane::sai::BfdEncapsulationType val) {
  switch (val) {
    case lemming::dataplane::sai::BFD_ENCAPSULATION_TYPE_IP_IN_IP:
      return SAI_BFD_ENCAPSULATION_TYPE_IP_IN_IP;

    case lemming::dataplane::sai::BFD_ENCAPSULATION_TYPE_L3_GRE_TUNNEL:
      return SAI_BFD_ENCAPSULATION_TYPE_L3_GRE_TUNNEL;

    case lemming::dataplane::sai::BFD_ENCAPSULATION_TYPE_NONE:
      return SAI_BFD_ENCAPSULATION_TYPE_NONE;

    default:
      return SAI_BFD_ENCAPSULATION_TYPE_IP_IN_IP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bfd_encapsulation_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bfd_encapsulation_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bfd_encapsulation_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bfd_encapsulation_type_t_to_sai(
        static_cast<lemming::dataplane::sai::BfdEncapsulationType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BfdSessionAttr convert_sai_bfd_session_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BFD_SESSION_ATTR_TYPE:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_TYPE;

    case SAI_BFD_SESSION_ATTR_HW_LOOKUP_VALID:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_HW_LOOKUP_VALID;

    case SAI_BFD_SESSION_ATTR_VIRTUAL_ROUTER:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_VIRTUAL_ROUTER;

    case SAI_BFD_SESSION_ATTR_PORT:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_PORT;

    case SAI_BFD_SESSION_ATTR_LOCAL_DISCRIMINATOR:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_LOCAL_DISCRIMINATOR;

    case SAI_BFD_SESSION_ATTR_REMOTE_DISCRIMINATOR:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_REMOTE_DISCRIMINATOR;

    case SAI_BFD_SESSION_ATTR_UDP_SRC_PORT:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_UDP_SRC_PORT;

    case SAI_BFD_SESSION_ATTR_TC:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_TC;

    case SAI_BFD_SESSION_ATTR_VLAN_TPID:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_VLAN_TPID;

    case SAI_BFD_SESSION_ATTR_VLAN_ID:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_VLAN_ID;

    case SAI_BFD_SESSION_ATTR_VLAN_PRI:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_VLAN_PRI;

    case SAI_BFD_SESSION_ATTR_VLAN_CFI:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_VLAN_CFI;

    case SAI_BFD_SESSION_ATTR_VLAN_HEADER_VALID:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_VLAN_HEADER_VALID;

    case SAI_BFD_SESSION_ATTR_BFD_ENCAPSULATION_TYPE:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_BFD_ENCAPSULATION_TYPE;

    case SAI_BFD_SESSION_ATTR_IPHDR_VERSION:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_IPHDR_VERSION;

    case SAI_BFD_SESSION_ATTR_TOS:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_TOS;

    case SAI_BFD_SESSION_ATTR_TTL:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_TTL;

    case SAI_BFD_SESSION_ATTR_SRC_IP_ADDRESS:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_SRC_IP_ADDRESS;

    case SAI_BFD_SESSION_ATTR_DST_IP_ADDRESS:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_DST_IP_ADDRESS;

    case SAI_BFD_SESSION_ATTR_TUNNEL_TOS:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_TUNNEL_TOS;

    case SAI_BFD_SESSION_ATTR_TUNNEL_TTL:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_TUNNEL_TTL;

    case SAI_BFD_SESSION_ATTR_TUNNEL_SRC_IP_ADDRESS:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_TUNNEL_SRC_IP_ADDRESS;

    case SAI_BFD_SESSION_ATTR_TUNNEL_DST_IP_ADDRESS:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_TUNNEL_DST_IP_ADDRESS;

    case SAI_BFD_SESSION_ATTR_SRC_MAC_ADDRESS:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_SRC_MAC_ADDRESS;

    case SAI_BFD_SESSION_ATTR_DST_MAC_ADDRESS:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_DST_MAC_ADDRESS;

    case SAI_BFD_SESSION_ATTR_ECHO_ENABLE:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_ECHO_ENABLE;

    case SAI_BFD_SESSION_ATTR_MULTIHOP:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_MULTIHOP;

    case SAI_BFD_SESSION_ATTR_CBIT:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_CBIT;

    case SAI_BFD_SESSION_ATTR_MIN_TX:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_MIN_TX;

    case SAI_BFD_SESSION_ATTR_MIN_RX:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_MIN_RX;

    case SAI_BFD_SESSION_ATTR_MULTIPLIER:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_MULTIPLIER;

    case SAI_BFD_SESSION_ATTR_REMOTE_MIN_TX:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_REMOTE_MIN_TX;

    case SAI_BFD_SESSION_ATTR_REMOTE_MIN_RX:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_REMOTE_MIN_RX;

    case SAI_BFD_SESSION_ATTR_STATE:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_STATE;

    case SAI_BFD_SESSION_ATTR_OFFLOAD_TYPE:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_OFFLOAD_TYPE;

    case SAI_BFD_SESSION_ATTR_NEGOTIATED_TX:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_NEGOTIATED_TX;

    case SAI_BFD_SESSION_ATTR_NEGOTIATED_RX:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_NEGOTIATED_RX;

    case SAI_BFD_SESSION_ATTR_LOCAL_DIAG:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_LOCAL_DIAG;

    case SAI_BFD_SESSION_ATTR_REMOTE_DIAG:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_REMOTE_DIAG;

    case SAI_BFD_SESSION_ATTR_REMOTE_MULTIPLIER:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_REMOTE_MULTIPLIER;

    default:
      return lemming::dataplane::sai::BFD_SESSION_ATTR_UNSPECIFIED;
  }
}
sai_bfd_session_attr_t convert_sai_bfd_session_attr_t_to_sai(
    lemming::dataplane::sai::BfdSessionAttr val) {
  switch (val) {
    case lemming::dataplane::sai::BFD_SESSION_ATTR_TYPE:
      return SAI_BFD_SESSION_ATTR_TYPE;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_HW_LOOKUP_VALID:
      return SAI_BFD_SESSION_ATTR_HW_LOOKUP_VALID;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_VIRTUAL_ROUTER:
      return SAI_BFD_SESSION_ATTR_VIRTUAL_ROUTER;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_PORT:
      return SAI_BFD_SESSION_ATTR_PORT;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_LOCAL_DISCRIMINATOR:
      return SAI_BFD_SESSION_ATTR_LOCAL_DISCRIMINATOR;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_REMOTE_DISCRIMINATOR:
      return SAI_BFD_SESSION_ATTR_REMOTE_DISCRIMINATOR;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_UDP_SRC_PORT:
      return SAI_BFD_SESSION_ATTR_UDP_SRC_PORT;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_TC:
      return SAI_BFD_SESSION_ATTR_TC;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_VLAN_TPID:
      return SAI_BFD_SESSION_ATTR_VLAN_TPID;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_VLAN_ID:
      return SAI_BFD_SESSION_ATTR_VLAN_ID;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_VLAN_PRI:
      return SAI_BFD_SESSION_ATTR_VLAN_PRI;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_VLAN_CFI:
      return SAI_BFD_SESSION_ATTR_VLAN_CFI;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_VLAN_HEADER_VALID:
      return SAI_BFD_SESSION_ATTR_VLAN_HEADER_VALID;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_BFD_ENCAPSULATION_TYPE:
      return SAI_BFD_SESSION_ATTR_BFD_ENCAPSULATION_TYPE;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_IPHDR_VERSION:
      return SAI_BFD_SESSION_ATTR_IPHDR_VERSION;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_TOS:
      return SAI_BFD_SESSION_ATTR_TOS;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_TTL:
      return SAI_BFD_SESSION_ATTR_TTL;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_SRC_IP_ADDRESS:
      return SAI_BFD_SESSION_ATTR_SRC_IP_ADDRESS;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_DST_IP_ADDRESS:
      return SAI_BFD_SESSION_ATTR_DST_IP_ADDRESS;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_TUNNEL_TOS:
      return SAI_BFD_SESSION_ATTR_TUNNEL_TOS;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_TUNNEL_TTL:
      return SAI_BFD_SESSION_ATTR_TUNNEL_TTL;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_TUNNEL_SRC_IP_ADDRESS:
      return SAI_BFD_SESSION_ATTR_TUNNEL_SRC_IP_ADDRESS;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_TUNNEL_DST_IP_ADDRESS:
      return SAI_BFD_SESSION_ATTR_TUNNEL_DST_IP_ADDRESS;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_SRC_MAC_ADDRESS:
      return SAI_BFD_SESSION_ATTR_SRC_MAC_ADDRESS;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_DST_MAC_ADDRESS:
      return SAI_BFD_SESSION_ATTR_DST_MAC_ADDRESS;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_ECHO_ENABLE:
      return SAI_BFD_SESSION_ATTR_ECHO_ENABLE;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_MULTIHOP:
      return SAI_BFD_SESSION_ATTR_MULTIHOP;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_CBIT:
      return SAI_BFD_SESSION_ATTR_CBIT;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_MIN_TX:
      return SAI_BFD_SESSION_ATTR_MIN_TX;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_MIN_RX:
      return SAI_BFD_SESSION_ATTR_MIN_RX;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_MULTIPLIER:
      return SAI_BFD_SESSION_ATTR_MULTIPLIER;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_REMOTE_MIN_TX:
      return SAI_BFD_SESSION_ATTR_REMOTE_MIN_TX;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_REMOTE_MIN_RX:
      return SAI_BFD_SESSION_ATTR_REMOTE_MIN_RX;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_STATE:
      return SAI_BFD_SESSION_ATTR_STATE;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_OFFLOAD_TYPE:
      return SAI_BFD_SESSION_ATTR_OFFLOAD_TYPE;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_NEGOTIATED_TX:
      return SAI_BFD_SESSION_ATTR_NEGOTIATED_TX;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_NEGOTIATED_RX:
      return SAI_BFD_SESSION_ATTR_NEGOTIATED_RX;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_LOCAL_DIAG:
      return SAI_BFD_SESSION_ATTR_LOCAL_DIAG;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_REMOTE_DIAG:
      return SAI_BFD_SESSION_ATTR_REMOTE_DIAG;

    case lemming::dataplane::sai::BFD_SESSION_ATTR_REMOTE_MULTIPLIER:
      return SAI_BFD_SESSION_ATTR_REMOTE_MULTIPLIER;

    default:
      return SAI_BFD_SESSION_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bfd_session_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bfd_session_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bfd_session_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bfd_session_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::BfdSessionAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BfdSessionOffloadType
convert_sai_bfd_session_offload_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_BFD_SESSION_OFFLOAD_TYPE_NONE:
      return lemming::dataplane::sai::BFD_SESSION_OFFLOAD_TYPE_NONE;

    case SAI_BFD_SESSION_OFFLOAD_TYPE_FULL:
      return lemming::dataplane::sai::BFD_SESSION_OFFLOAD_TYPE_FULL;

    case SAI_BFD_SESSION_OFFLOAD_TYPE_SUSTENANCE:
      return lemming::dataplane::sai::BFD_SESSION_OFFLOAD_TYPE_SUSTENANCE;

    default:
      return lemming::dataplane::sai::BFD_SESSION_OFFLOAD_TYPE_UNSPECIFIED;
  }
}
sai_bfd_session_offload_type_t convert_sai_bfd_session_offload_type_t_to_sai(
    lemming::dataplane::sai::BfdSessionOffloadType val) {
  switch (val) {
    case lemming::dataplane::sai::BFD_SESSION_OFFLOAD_TYPE_NONE:
      return SAI_BFD_SESSION_OFFLOAD_TYPE_NONE;

    case lemming::dataplane::sai::BFD_SESSION_OFFLOAD_TYPE_FULL:
      return SAI_BFD_SESSION_OFFLOAD_TYPE_FULL;

    case lemming::dataplane::sai::BFD_SESSION_OFFLOAD_TYPE_SUSTENANCE:
      return SAI_BFD_SESSION_OFFLOAD_TYPE_SUSTENANCE;

    default:
      return SAI_BFD_SESSION_OFFLOAD_TYPE_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bfd_session_offload_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_bfd_session_offload_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bfd_session_offload_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bfd_session_offload_type_t_to_sai(
        static_cast<lemming::dataplane::sai::BfdSessionOffloadType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BfdSessionStat convert_sai_bfd_session_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BFD_SESSION_STAT_IN_PACKETS:
      return lemming::dataplane::sai::BFD_SESSION_STAT_IN_PACKETS;

    case SAI_BFD_SESSION_STAT_OUT_PACKETS:
      return lemming::dataplane::sai::BFD_SESSION_STAT_OUT_PACKETS;

    case SAI_BFD_SESSION_STAT_DROP_PACKETS:
      return lemming::dataplane::sai::BFD_SESSION_STAT_DROP_PACKETS;

    default:
      return lemming::dataplane::sai::BFD_SESSION_STAT_UNSPECIFIED;
  }
}
sai_bfd_session_stat_t convert_sai_bfd_session_stat_t_to_sai(
    lemming::dataplane::sai::BfdSessionStat val) {
  switch (val) {
    case lemming::dataplane::sai::BFD_SESSION_STAT_IN_PACKETS:
      return SAI_BFD_SESSION_STAT_IN_PACKETS;

    case lemming::dataplane::sai::BFD_SESSION_STAT_OUT_PACKETS:
      return SAI_BFD_SESSION_STAT_OUT_PACKETS;

    case lemming::dataplane::sai::BFD_SESSION_STAT_DROP_PACKETS:
      return SAI_BFD_SESSION_STAT_DROP_PACKETS;

    default:
      return SAI_BFD_SESSION_STAT_IN_PACKETS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bfd_session_stat_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bfd_session_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bfd_session_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bfd_session_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::BfdSessionStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BfdSessionState
convert_sai_bfd_session_state_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_BFD_SESSION_STATE_ADMIN_DOWN:
      return lemming::dataplane::sai::BFD_SESSION_STATE_ADMIN_DOWN;

    case SAI_BFD_SESSION_STATE_DOWN:
      return lemming::dataplane::sai::BFD_SESSION_STATE_DOWN;

    case SAI_BFD_SESSION_STATE_INIT:
      return lemming::dataplane::sai::BFD_SESSION_STATE_INIT;

    case SAI_BFD_SESSION_STATE_UP:
      return lemming::dataplane::sai::BFD_SESSION_STATE_UP;

    default:
      return lemming::dataplane::sai::BFD_SESSION_STATE_UNSPECIFIED;
  }
}
sai_bfd_session_state_t convert_sai_bfd_session_state_t_to_sai(
    lemming::dataplane::sai::BfdSessionState val) {
  switch (val) {
    case lemming::dataplane::sai::BFD_SESSION_STATE_ADMIN_DOWN:
      return SAI_BFD_SESSION_STATE_ADMIN_DOWN;

    case lemming::dataplane::sai::BFD_SESSION_STATE_DOWN:
      return SAI_BFD_SESSION_STATE_DOWN;

    case lemming::dataplane::sai::BFD_SESSION_STATE_INIT:
      return SAI_BFD_SESSION_STATE_INIT;

    case lemming::dataplane::sai::BFD_SESSION_STATE_UP:
      return SAI_BFD_SESSION_STATE_UP;

    default:
      return SAI_BFD_SESSION_STATE_ADMIN_DOWN;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bfd_session_state_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bfd_session_state_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bfd_session_state_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bfd_session_state_t_to_sai(
        static_cast<lemming::dataplane::sai::BfdSessionState>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BfdSessionType convert_sai_bfd_session_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BFD_SESSION_TYPE_DEMAND_ACTIVE:
      return lemming::dataplane::sai::BFD_SESSION_TYPE_DEMAND_ACTIVE;

    case SAI_BFD_SESSION_TYPE_DEMAND_PASSIVE:
      return lemming::dataplane::sai::BFD_SESSION_TYPE_DEMAND_PASSIVE;

    case SAI_BFD_SESSION_TYPE_ASYNC_ACTIVE:
      return lemming::dataplane::sai::BFD_SESSION_TYPE_ASYNC_ACTIVE;

    case SAI_BFD_SESSION_TYPE_ASYNC_PASSIVE:
      return lemming::dataplane::sai::BFD_SESSION_TYPE_ASYNC_PASSIVE;

    default:
      return lemming::dataplane::sai::BFD_SESSION_TYPE_UNSPECIFIED;
  }
}
sai_bfd_session_type_t convert_sai_bfd_session_type_t_to_sai(
    lemming::dataplane::sai::BfdSessionType val) {
  switch (val) {
    case lemming::dataplane::sai::BFD_SESSION_TYPE_DEMAND_ACTIVE:
      return SAI_BFD_SESSION_TYPE_DEMAND_ACTIVE;

    case lemming::dataplane::sai::BFD_SESSION_TYPE_DEMAND_PASSIVE:
      return SAI_BFD_SESSION_TYPE_DEMAND_PASSIVE;

    case lemming::dataplane::sai::BFD_SESSION_TYPE_ASYNC_ACTIVE:
      return SAI_BFD_SESSION_TYPE_ASYNC_ACTIVE;

    case lemming::dataplane::sai::BFD_SESSION_TYPE_ASYNC_PASSIVE:
      return SAI_BFD_SESSION_TYPE_ASYNC_PASSIVE;

    default:
      return SAI_BFD_SESSION_TYPE_DEMAND_ACTIVE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bfd_session_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bfd_session_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bfd_session_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bfd_session_type_t_to_sai(
        static_cast<lemming::dataplane::sai::BfdSessionType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BridgeAttr convert_sai_bridge_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BRIDGE_ATTR_TYPE:
      return lemming::dataplane::sai::BRIDGE_ATTR_TYPE;

    case SAI_BRIDGE_ATTR_PORT_LIST:
      return lemming::dataplane::sai::BRIDGE_ATTR_PORT_LIST;

    case SAI_BRIDGE_ATTR_MAX_LEARNED_ADDRESSES:
      return lemming::dataplane::sai::BRIDGE_ATTR_MAX_LEARNED_ADDRESSES;

    case SAI_BRIDGE_ATTR_LEARN_DISABLE:
      return lemming::dataplane::sai::BRIDGE_ATTR_LEARN_DISABLE;

    case SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
      return lemming::dataplane::sai::
          BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE;

    case SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
      return lemming::dataplane::sai::BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP;

    case SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
      return lemming::dataplane::sai::
          BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE;

    case SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
      return lemming::dataplane::sai::BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP;

    case SAI_BRIDGE_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
      return lemming::dataplane::sai::BRIDGE_ATTR_BROADCAST_FLOOD_CONTROL_TYPE;

    case SAI_BRIDGE_ATTR_BROADCAST_FLOOD_GROUP:
      return lemming::dataplane::sai::BRIDGE_ATTR_BROADCAST_FLOOD_GROUP;

    default:
      return lemming::dataplane::sai::BRIDGE_ATTR_UNSPECIFIED;
  }
}
sai_bridge_attr_t convert_sai_bridge_attr_t_to_sai(
    lemming::dataplane::sai::BridgeAttr val) {
  switch (val) {
    case lemming::dataplane::sai::BRIDGE_ATTR_TYPE:
      return SAI_BRIDGE_ATTR_TYPE;

    case lemming::dataplane::sai::BRIDGE_ATTR_PORT_LIST:
      return SAI_BRIDGE_ATTR_PORT_LIST;

    case lemming::dataplane::sai::BRIDGE_ATTR_MAX_LEARNED_ADDRESSES:
      return SAI_BRIDGE_ATTR_MAX_LEARNED_ADDRESSES;

    case lemming::dataplane::sai::BRIDGE_ATTR_LEARN_DISABLE:
      return SAI_BRIDGE_ATTR_LEARN_DISABLE;

    case lemming::dataplane::sai::
        BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
      return SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE;

    case lemming::dataplane::sai::BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
      return SAI_BRIDGE_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP;

    case lemming::dataplane::sai::
        BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
      return SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE;

    case lemming::dataplane::sai::BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
      return SAI_BRIDGE_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP;

    case lemming::dataplane::sai::BRIDGE_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
      return SAI_BRIDGE_ATTR_BROADCAST_FLOOD_CONTROL_TYPE;

    case lemming::dataplane::sai::BRIDGE_ATTR_BROADCAST_FLOOD_GROUP:
      return SAI_BRIDGE_ATTR_BROADCAST_FLOOD_GROUP;

    default:
      return SAI_BRIDGE_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_bridge_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bridge_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bridge_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bridge_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::BridgeAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BridgeFloodControlType
convert_sai_bridge_flood_control_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_BRIDGE_FLOOD_CONTROL_TYPE_SUB_PORTS:
      return lemming::dataplane::sai::BRIDGE_FLOOD_CONTROL_TYPE_SUB_PORTS;

    case SAI_BRIDGE_FLOOD_CONTROL_TYPE_NONE:
      return lemming::dataplane::sai::BRIDGE_FLOOD_CONTROL_TYPE_NONE;

    case SAI_BRIDGE_FLOOD_CONTROL_TYPE_L2MC_GROUP:
      return lemming::dataplane::sai::BRIDGE_FLOOD_CONTROL_TYPE_L2MC_GROUP;

    case SAI_BRIDGE_FLOOD_CONTROL_TYPE_COMBINED:
      return lemming::dataplane::sai::BRIDGE_FLOOD_CONTROL_TYPE_COMBINED;

    default:
      return lemming::dataplane::sai::BRIDGE_FLOOD_CONTROL_TYPE_UNSPECIFIED;
  }
}
sai_bridge_flood_control_type_t convert_sai_bridge_flood_control_type_t_to_sai(
    lemming::dataplane::sai::BridgeFloodControlType val) {
  switch (val) {
    case lemming::dataplane::sai::BRIDGE_FLOOD_CONTROL_TYPE_SUB_PORTS:
      return SAI_BRIDGE_FLOOD_CONTROL_TYPE_SUB_PORTS;

    case lemming::dataplane::sai::BRIDGE_FLOOD_CONTROL_TYPE_NONE:
      return SAI_BRIDGE_FLOOD_CONTROL_TYPE_NONE;

    case lemming::dataplane::sai::BRIDGE_FLOOD_CONTROL_TYPE_L2MC_GROUP:
      return SAI_BRIDGE_FLOOD_CONTROL_TYPE_L2MC_GROUP;

    case lemming::dataplane::sai::BRIDGE_FLOOD_CONTROL_TYPE_COMBINED:
      return SAI_BRIDGE_FLOOD_CONTROL_TYPE_COMBINED;

    default:
      return SAI_BRIDGE_FLOOD_CONTROL_TYPE_SUB_PORTS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bridge_flood_control_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_bridge_flood_control_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bridge_flood_control_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bridge_flood_control_type_t_to_sai(
        static_cast<lemming::dataplane::sai::BridgeFloodControlType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BridgePortAttr convert_sai_bridge_port_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BRIDGE_PORT_ATTR_TYPE:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_TYPE;

    case SAI_BRIDGE_PORT_ATTR_PORT_ID:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_PORT_ID;

    case SAI_BRIDGE_PORT_ATTR_TAGGING_MODE:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_TAGGING_MODE;

    case SAI_BRIDGE_PORT_ATTR_VLAN_ID:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_VLAN_ID;

    case SAI_BRIDGE_PORT_ATTR_RIF_ID:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_RIF_ID;

    case SAI_BRIDGE_PORT_ATTR_TUNNEL_ID:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_TUNNEL_ID;

    case SAI_BRIDGE_PORT_ATTR_BRIDGE_ID:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_BRIDGE_ID;

    case SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_MODE:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_FDB_LEARNING_MODE;

    case SAI_BRIDGE_PORT_ATTR_MAX_LEARNED_ADDRESSES:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_MAX_LEARNED_ADDRESSES;

    case SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_LIMIT_VIOLATION_PACKET_ACTION:
      return lemming::dataplane::sai::
          BRIDGE_PORT_ATTR_FDB_LEARNING_LIMIT_VIOLATION_PACKET_ACTION;

    case SAI_BRIDGE_PORT_ATTR_ADMIN_STATE:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_ADMIN_STATE;

    case SAI_BRIDGE_PORT_ATTR_INGRESS_FILTERING:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_INGRESS_FILTERING;

    case SAI_BRIDGE_PORT_ATTR_EGRESS_FILTERING:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_EGRESS_FILTERING;

    case SAI_BRIDGE_PORT_ATTR_ISOLATION_GROUP:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_ISOLATION_GROUP;

    default:
      return lemming::dataplane::sai::BRIDGE_PORT_ATTR_UNSPECIFIED;
  }
}
sai_bridge_port_attr_t convert_sai_bridge_port_attr_t_to_sai(
    lemming::dataplane::sai::BridgePortAttr val) {
  switch (val) {
    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_TYPE:
      return SAI_BRIDGE_PORT_ATTR_TYPE;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_PORT_ID:
      return SAI_BRIDGE_PORT_ATTR_PORT_ID;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_TAGGING_MODE:
      return SAI_BRIDGE_PORT_ATTR_TAGGING_MODE;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_VLAN_ID:
      return SAI_BRIDGE_PORT_ATTR_VLAN_ID;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_RIF_ID:
      return SAI_BRIDGE_PORT_ATTR_RIF_ID;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_TUNNEL_ID:
      return SAI_BRIDGE_PORT_ATTR_TUNNEL_ID;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_BRIDGE_ID:
      return SAI_BRIDGE_PORT_ATTR_BRIDGE_ID;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_FDB_LEARNING_MODE:
      return SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_MODE;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_MAX_LEARNED_ADDRESSES:
      return SAI_BRIDGE_PORT_ATTR_MAX_LEARNED_ADDRESSES;

    case lemming::dataplane::sai::
        BRIDGE_PORT_ATTR_FDB_LEARNING_LIMIT_VIOLATION_PACKET_ACTION:
      return SAI_BRIDGE_PORT_ATTR_FDB_LEARNING_LIMIT_VIOLATION_PACKET_ACTION;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_ADMIN_STATE:
      return SAI_BRIDGE_PORT_ATTR_ADMIN_STATE;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_INGRESS_FILTERING:
      return SAI_BRIDGE_PORT_ATTR_INGRESS_FILTERING;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_EGRESS_FILTERING:
      return SAI_BRIDGE_PORT_ATTR_EGRESS_FILTERING;

    case lemming::dataplane::sai::BRIDGE_PORT_ATTR_ISOLATION_GROUP:
      return SAI_BRIDGE_PORT_ATTR_ISOLATION_GROUP;

    default:
      return SAI_BRIDGE_PORT_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bridge_port_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bridge_port_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bridge_port_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bridge_port_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::BridgePortAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BridgePortFdbLearningMode
convert_sai_bridge_port_fdb_learning_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_BRIDGE_PORT_FDB_LEARNING_MODE_DROP:
      return lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_DROP;

    case SAI_BRIDGE_PORT_FDB_LEARNING_MODE_DISABLE:
      return lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_DISABLE;

    case SAI_BRIDGE_PORT_FDB_LEARNING_MODE_HW:
      return lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_HW;

    case SAI_BRIDGE_PORT_FDB_LEARNING_MODE_CPU_TRAP:
      return lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_CPU_TRAP;

    case SAI_BRIDGE_PORT_FDB_LEARNING_MODE_CPU_LOG:
      return lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_CPU_LOG;

    case SAI_BRIDGE_PORT_FDB_LEARNING_MODE_FDB_NOTIFICATION:
      return lemming::dataplane::sai::
          BRIDGE_PORT_FDB_LEARNING_MODE_FDB_NOTIFICATION;

    default:
      return lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_UNSPECIFIED;
  }
}
sai_bridge_port_fdb_learning_mode_t
convert_sai_bridge_port_fdb_learning_mode_t_to_sai(
    lemming::dataplane::sai::BridgePortFdbLearningMode val) {
  switch (val) {
    case lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_DROP:
      return SAI_BRIDGE_PORT_FDB_LEARNING_MODE_DROP;

    case lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_DISABLE:
      return SAI_BRIDGE_PORT_FDB_LEARNING_MODE_DISABLE;

    case lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_HW:
      return SAI_BRIDGE_PORT_FDB_LEARNING_MODE_HW;

    case lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_CPU_TRAP:
      return SAI_BRIDGE_PORT_FDB_LEARNING_MODE_CPU_TRAP;

    case lemming::dataplane::sai::BRIDGE_PORT_FDB_LEARNING_MODE_CPU_LOG:
      return SAI_BRIDGE_PORT_FDB_LEARNING_MODE_CPU_LOG;

    case lemming::dataplane::sai::
        BRIDGE_PORT_FDB_LEARNING_MODE_FDB_NOTIFICATION:
      return SAI_BRIDGE_PORT_FDB_LEARNING_MODE_FDB_NOTIFICATION;

    default:
      return SAI_BRIDGE_PORT_FDB_LEARNING_MODE_DROP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bridge_port_fdb_learning_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_bridge_port_fdb_learning_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bridge_port_fdb_learning_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bridge_port_fdb_learning_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::BridgePortFdbLearningMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BridgePortStat convert_sai_bridge_port_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BRIDGE_PORT_STAT_IN_OCTETS:
      return lemming::dataplane::sai::BRIDGE_PORT_STAT_IN_OCTETS;

    case SAI_BRIDGE_PORT_STAT_IN_PACKETS:
      return lemming::dataplane::sai::BRIDGE_PORT_STAT_IN_PACKETS;

    case SAI_BRIDGE_PORT_STAT_OUT_OCTETS:
      return lemming::dataplane::sai::BRIDGE_PORT_STAT_OUT_OCTETS;

    case SAI_BRIDGE_PORT_STAT_OUT_PACKETS:
      return lemming::dataplane::sai::BRIDGE_PORT_STAT_OUT_PACKETS;

    default:
      return lemming::dataplane::sai::BRIDGE_PORT_STAT_UNSPECIFIED;
  }
}
sai_bridge_port_stat_t convert_sai_bridge_port_stat_t_to_sai(
    lemming::dataplane::sai::BridgePortStat val) {
  switch (val) {
    case lemming::dataplane::sai::BRIDGE_PORT_STAT_IN_OCTETS:
      return SAI_BRIDGE_PORT_STAT_IN_OCTETS;

    case lemming::dataplane::sai::BRIDGE_PORT_STAT_IN_PACKETS:
      return SAI_BRIDGE_PORT_STAT_IN_PACKETS;

    case lemming::dataplane::sai::BRIDGE_PORT_STAT_OUT_OCTETS:
      return SAI_BRIDGE_PORT_STAT_OUT_OCTETS;

    case lemming::dataplane::sai::BRIDGE_PORT_STAT_OUT_PACKETS:
      return SAI_BRIDGE_PORT_STAT_OUT_PACKETS;

    default:
      return SAI_BRIDGE_PORT_STAT_IN_OCTETS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bridge_port_stat_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bridge_port_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bridge_port_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bridge_port_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::BridgePortStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BridgePortTaggingMode
convert_sai_bridge_port_tagging_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_BRIDGE_PORT_TAGGING_MODE_UNTAGGED:
      return lemming::dataplane::sai::BRIDGE_PORT_TAGGING_MODE_UNTAGGED;

    case SAI_BRIDGE_PORT_TAGGING_MODE_TAGGED:
      return lemming::dataplane::sai::BRIDGE_PORT_TAGGING_MODE_TAGGED;

    default:
      return lemming::dataplane::sai::BRIDGE_PORT_TAGGING_MODE_UNSPECIFIED;
  }
}
sai_bridge_port_tagging_mode_t convert_sai_bridge_port_tagging_mode_t_to_sai(
    lemming::dataplane::sai::BridgePortTaggingMode val) {
  switch (val) {
    case lemming::dataplane::sai::BRIDGE_PORT_TAGGING_MODE_UNTAGGED:
      return SAI_BRIDGE_PORT_TAGGING_MODE_UNTAGGED;

    case lemming::dataplane::sai::BRIDGE_PORT_TAGGING_MODE_TAGGED:
      return SAI_BRIDGE_PORT_TAGGING_MODE_TAGGED;

    default:
      return SAI_BRIDGE_PORT_TAGGING_MODE_UNTAGGED;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bridge_port_tagging_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_bridge_port_tagging_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bridge_port_tagging_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bridge_port_tagging_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::BridgePortTaggingMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BridgePortType convert_sai_bridge_port_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BRIDGE_PORT_TYPE_PORT:
      return lemming::dataplane::sai::BRIDGE_PORT_TYPE_PORT;

    case SAI_BRIDGE_PORT_TYPE_SUB_PORT:
      return lemming::dataplane::sai::BRIDGE_PORT_TYPE_SUB_PORT;

    case SAI_BRIDGE_PORT_TYPE_1Q_ROUTER:
      return lemming::dataplane::sai::BRIDGE_PORT_TYPE_1Q_ROUTER;

    case SAI_BRIDGE_PORT_TYPE_1D_ROUTER:
      return lemming::dataplane::sai::BRIDGE_PORT_TYPE_1D_ROUTER;

    case SAI_BRIDGE_PORT_TYPE_TUNNEL:
      return lemming::dataplane::sai::BRIDGE_PORT_TYPE_TUNNEL;

    default:
      return lemming::dataplane::sai::BRIDGE_PORT_TYPE_UNSPECIFIED;
  }
}
sai_bridge_port_type_t convert_sai_bridge_port_type_t_to_sai(
    lemming::dataplane::sai::BridgePortType val) {
  switch (val) {
    case lemming::dataplane::sai::BRIDGE_PORT_TYPE_PORT:
      return SAI_BRIDGE_PORT_TYPE_PORT;

    case lemming::dataplane::sai::BRIDGE_PORT_TYPE_SUB_PORT:
      return SAI_BRIDGE_PORT_TYPE_SUB_PORT;

    case lemming::dataplane::sai::BRIDGE_PORT_TYPE_1Q_ROUTER:
      return SAI_BRIDGE_PORT_TYPE_1Q_ROUTER;

    case lemming::dataplane::sai::BRIDGE_PORT_TYPE_1D_ROUTER:
      return SAI_BRIDGE_PORT_TYPE_1D_ROUTER;

    case lemming::dataplane::sai::BRIDGE_PORT_TYPE_TUNNEL:
      return SAI_BRIDGE_PORT_TYPE_TUNNEL;

    default:
      return SAI_BRIDGE_PORT_TYPE_PORT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bridge_port_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bridge_port_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bridge_port_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bridge_port_type_t_to_sai(
        static_cast<lemming::dataplane::sai::BridgePortType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BridgeStat convert_sai_bridge_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BRIDGE_STAT_IN_OCTETS:
      return lemming::dataplane::sai::BRIDGE_STAT_IN_OCTETS;

    case SAI_BRIDGE_STAT_IN_PACKETS:
      return lemming::dataplane::sai::BRIDGE_STAT_IN_PACKETS;

    case SAI_BRIDGE_STAT_OUT_OCTETS:
      return lemming::dataplane::sai::BRIDGE_STAT_OUT_OCTETS;

    case SAI_BRIDGE_STAT_OUT_PACKETS:
      return lemming::dataplane::sai::BRIDGE_STAT_OUT_PACKETS;

    default:
      return lemming::dataplane::sai::BRIDGE_STAT_UNSPECIFIED;
  }
}
sai_bridge_stat_t convert_sai_bridge_stat_t_to_sai(
    lemming::dataplane::sai::BridgeStat val) {
  switch (val) {
    case lemming::dataplane::sai::BRIDGE_STAT_IN_OCTETS:
      return SAI_BRIDGE_STAT_IN_OCTETS;

    case lemming::dataplane::sai::BRIDGE_STAT_IN_PACKETS:
      return SAI_BRIDGE_STAT_IN_PACKETS;

    case lemming::dataplane::sai::BRIDGE_STAT_OUT_OCTETS:
      return SAI_BRIDGE_STAT_OUT_OCTETS;

    case lemming::dataplane::sai::BRIDGE_STAT_OUT_PACKETS:
      return SAI_BRIDGE_STAT_OUT_PACKETS;

    default:
      return SAI_BRIDGE_STAT_IN_OCTETS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_bridge_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bridge_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bridge_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bridge_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::BridgeStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BridgeType convert_sai_bridge_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BRIDGE_TYPE_1Q:
      return lemming::dataplane::sai::BRIDGE_TYPE_1Q;

    case SAI_BRIDGE_TYPE_1D:
      return lemming::dataplane::sai::BRIDGE_TYPE_1D;

    default:
      return lemming::dataplane::sai::BRIDGE_TYPE_UNSPECIFIED;
  }
}
sai_bridge_type_t convert_sai_bridge_type_t_to_sai(
    lemming::dataplane::sai::BridgeType val) {
  switch (val) {
    case lemming::dataplane::sai::BRIDGE_TYPE_1Q:
      return SAI_BRIDGE_TYPE_1Q;

    case lemming::dataplane::sai::BRIDGE_TYPE_1D:
      return SAI_BRIDGE_TYPE_1D;

    default:
      return SAI_BRIDGE_TYPE_1Q;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_bridge_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bridge_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bridge_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bridge_type_t_to_sai(
        static_cast<lemming::dataplane::sai::BridgeType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BufferPoolAttr convert_sai_buffer_pool_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BUFFER_POOL_ATTR_SHARED_SIZE:
      return lemming::dataplane::sai::BUFFER_POOL_ATTR_SHARED_SIZE;

    case SAI_BUFFER_POOL_ATTR_TYPE:
      return lemming::dataplane::sai::BUFFER_POOL_ATTR_TYPE;

    case SAI_BUFFER_POOL_ATTR_SIZE:
      return lemming::dataplane::sai::BUFFER_POOL_ATTR_SIZE;

    case SAI_BUFFER_POOL_ATTR_THRESHOLD_MODE:
      return lemming::dataplane::sai::BUFFER_POOL_ATTR_THRESHOLD_MODE;

    case SAI_BUFFER_POOL_ATTR_TAM:
      return lemming::dataplane::sai::BUFFER_POOL_ATTR_TAM;

    case SAI_BUFFER_POOL_ATTR_XOFF_SIZE:
      return lemming::dataplane::sai::BUFFER_POOL_ATTR_XOFF_SIZE;

    case SAI_BUFFER_POOL_ATTR_WRED_PROFILE_ID:
      return lemming::dataplane::sai::BUFFER_POOL_ATTR_WRED_PROFILE_ID;

    default:
      return lemming::dataplane::sai::BUFFER_POOL_ATTR_UNSPECIFIED;
  }
}
sai_buffer_pool_attr_t convert_sai_buffer_pool_attr_t_to_sai(
    lemming::dataplane::sai::BufferPoolAttr val) {
  switch (val) {
    case lemming::dataplane::sai::BUFFER_POOL_ATTR_SHARED_SIZE:
      return SAI_BUFFER_POOL_ATTR_SHARED_SIZE;

    case lemming::dataplane::sai::BUFFER_POOL_ATTR_TYPE:
      return SAI_BUFFER_POOL_ATTR_TYPE;

    case lemming::dataplane::sai::BUFFER_POOL_ATTR_SIZE:
      return SAI_BUFFER_POOL_ATTR_SIZE;

    case lemming::dataplane::sai::BUFFER_POOL_ATTR_THRESHOLD_MODE:
      return SAI_BUFFER_POOL_ATTR_THRESHOLD_MODE;

    case lemming::dataplane::sai::BUFFER_POOL_ATTR_TAM:
      return SAI_BUFFER_POOL_ATTR_TAM;

    case lemming::dataplane::sai::BUFFER_POOL_ATTR_XOFF_SIZE:
      return SAI_BUFFER_POOL_ATTR_XOFF_SIZE;

    case lemming::dataplane::sai::BUFFER_POOL_ATTR_WRED_PROFILE_ID:
      return SAI_BUFFER_POOL_ATTR_WRED_PROFILE_ID;

    default:
      return SAI_BUFFER_POOL_ATTR_SHARED_SIZE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_buffer_pool_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_buffer_pool_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_buffer_pool_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_buffer_pool_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::BufferPoolAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BufferPoolStat convert_sai_buffer_pool_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BUFFER_POOL_STAT_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_CURR_OCCUPANCY_BYTES;

    case SAI_BUFFER_POOL_STAT_WATERMARK_BYTES:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_WATERMARK_BYTES;

    case SAI_BUFFER_POOL_STAT_DROPPED_PACKETS:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_DROPPED_PACKETS;

    case SAI_BUFFER_POOL_STAT_GREEN_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_GREEN_WRED_DROPPED_PACKETS;

    case SAI_BUFFER_POOL_STAT_GREEN_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_GREEN_WRED_DROPPED_BYTES;

    case SAI_BUFFER_POOL_STAT_YELLOW_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_YELLOW_WRED_DROPPED_PACKETS;

    case SAI_BUFFER_POOL_STAT_YELLOW_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_YELLOW_WRED_DROPPED_BYTES;

    case SAI_BUFFER_POOL_STAT_RED_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_RED_WRED_DROPPED_PACKETS;

    case SAI_BUFFER_POOL_STAT_RED_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_RED_WRED_DROPPED_BYTES;

    case SAI_BUFFER_POOL_STAT_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_WRED_DROPPED_PACKETS;

    case SAI_BUFFER_POOL_STAT_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_WRED_DROPPED_BYTES;

    case SAI_BUFFER_POOL_STAT_GREEN_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_GREEN_WRED_ECN_MARKED_PACKETS;

    case SAI_BUFFER_POOL_STAT_GREEN_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_GREEN_WRED_ECN_MARKED_BYTES;

    case SAI_BUFFER_POOL_STAT_YELLOW_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_YELLOW_WRED_ECN_MARKED_PACKETS;

    case SAI_BUFFER_POOL_STAT_YELLOW_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_YELLOW_WRED_ECN_MARKED_BYTES;

    case SAI_BUFFER_POOL_STAT_RED_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_RED_WRED_ECN_MARKED_PACKETS;

    case SAI_BUFFER_POOL_STAT_RED_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_RED_WRED_ECN_MARKED_BYTES;

    case SAI_BUFFER_POOL_STAT_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_WRED_ECN_MARKED_PACKETS;

    case SAI_BUFFER_POOL_STAT_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_WRED_ECN_MARKED_BYTES;

    case SAI_BUFFER_POOL_STAT_XOFF_ROOM_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_XOFF_ROOM_CURR_OCCUPANCY_BYTES;

    case SAI_BUFFER_POOL_STAT_XOFF_ROOM_WATERMARK_BYTES:
      return lemming::dataplane::sai::
          BUFFER_POOL_STAT_XOFF_ROOM_WATERMARK_BYTES;

    case SAI_BUFFER_POOL_STAT_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::BUFFER_POOL_STAT_UNSPECIFIED;
  }
}
sai_buffer_pool_stat_t convert_sai_buffer_pool_stat_t_to_sai(
    lemming::dataplane::sai::BufferPoolStat val) {
  switch (val) {
    case lemming::dataplane::sai::BUFFER_POOL_STAT_CURR_OCCUPANCY_BYTES:
      return SAI_BUFFER_POOL_STAT_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_WATERMARK_BYTES:
      return SAI_BUFFER_POOL_STAT_WATERMARK_BYTES;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_DROPPED_PACKETS:
      return SAI_BUFFER_POOL_STAT_DROPPED_PACKETS;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_GREEN_WRED_DROPPED_PACKETS:
      return SAI_BUFFER_POOL_STAT_GREEN_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_GREEN_WRED_DROPPED_BYTES:
      return SAI_BUFFER_POOL_STAT_GREEN_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_YELLOW_WRED_DROPPED_PACKETS:
      return SAI_BUFFER_POOL_STAT_YELLOW_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_YELLOW_WRED_DROPPED_BYTES:
      return SAI_BUFFER_POOL_STAT_YELLOW_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_RED_WRED_DROPPED_PACKETS:
      return SAI_BUFFER_POOL_STAT_RED_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_RED_WRED_DROPPED_BYTES:
      return SAI_BUFFER_POOL_STAT_RED_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_WRED_DROPPED_PACKETS:
      return SAI_BUFFER_POOL_STAT_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_WRED_DROPPED_BYTES:
      return SAI_BUFFER_POOL_STAT_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::
        BUFFER_POOL_STAT_GREEN_WRED_ECN_MARKED_PACKETS:
      return SAI_BUFFER_POOL_STAT_GREEN_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_GREEN_WRED_ECN_MARKED_BYTES:
      return SAI_BUFFER_POOL_STAT_GREEN_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::
        BUFFER_POOL_STAT_YELLOW_WRED_ECN_MARKED_PACKETS:
      return SAI_BUFFER_POOL_STAT_YELLOW_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_YELLOW_WRED_ECN_MARKED_BYTES:
      return SAI_BUFFER_POOL_STAT_YELLOW_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_RED_WRED_ECN_MARKED_PACKETS:
      return SAI_BUFFER_POOL_STAT_RED_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_RED_WRED_ECN_MARKED_BYTES:
      return SAI_BUFFER_POOL_STAT_RED_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_WRED_ECN_MARKED_PACKETS:
      return SAI_BUFFER_POOL_STAT_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_WRED_ECN_MARKED_BYTES:
      return SAI_BUFFER_POOL_STAT_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::
        BUFFER_POOL_STAT_XOFF_ROOM_CURR_OCCUPANCY_BYTES:
      return SAI_BUFFER_POOL_STAT_XOFF_ROOM_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_XOFF_ROOM_WATERMARK_BYTES:
      return SAI_BUFFER_POOL_STAT_XOFF_ROOM_WATERMARK_BYTES;

    case lemming::dataplane::sai::BUFFER_POOL_STAT_CUSTOM_RANGE_BASE:
      return SAI_BUFFER_POOL_STAT_CUSTOM_RANGE_BASE;

    default:
      return SAI_BUFFER_POOL_STAT_CURR_OCCUPANCY_BYTES;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_buffer_pool_stat_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_buffer_pool_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_buffer_pool_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_buffer_pool_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::BufferPoolStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BufferPoolThresholdMode
convert_sai_buffer_pool_threshold_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_BUFFER_POOL_THRESHOLD_MODE_STATIC:
      return lemming::dataplane::sai::BUFFER_POOL_THRESHOLD_MODE_STATIC;

    case SAI_BUFFER_POOL_THRESHOLD_MODE_DYNAMIC:
      return lemming::dataplane::sai::BUFFER_POOL_THRESHOLD_MODE_DYNAMIC;

    default:
      return lemming::dataplane::sai::BUFFER_POOL_THRESHOLD_MODE_UNSPECIFIED;
  }
}
sai_buffer_pool_threshold_mode_t
convert_sai_buffer_pool_threshold_mode_t_to_sai(
    lemming::dataplane::sai::BufferPoolThresholdMode val) {
  switch (val) {
    case lemming::dataplane::sai::BUFFER_POOL_THRESHOLD_MODE_STATIC:
      return SAI_BUFFER_POOL_THRESHOLD_MODE_STATIC;

    case lemming::dataplane::sai::BUFFER_POOL_THRESHOLD_MODE_DYNAMIC:
      return SAI_BUFFER_POOL_THRESHOLD_MODE_DYNAMIC;

    default:
      return SAI_BUFFER_POOL_THRESHOLD_MODE_STATIC;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_buffer_pool_threshold_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_buffer_pool_threshold_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_buffer_pool_threshold_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_buffer_pool_threshold_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::BufferPoolThresholdMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BufferPoolType convert_sai_buffer_pool_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_BUFFER_POOL_TYPE_INGRESS:
      return lemming::dataplane::sai::BUFFER_POOL_TYPE_INGRESS;

    case SAI_BUFFER_POOL_TYPE_EGRESS:
      return lemming::dataplane::sai::BUFFER_POOL_TYPE_EGRESS;

    case SAI_BUFFER_POOL_TYPE_BOTH:
      return lemming::dataplane::sai::BUFFER_POOL_TYPE_BOTH;

    default:
      return lemming::dataplane::sai::BUFFER_POOL_TYPE_UNSPECIFIED;
  }
}
sai_buffer_pool_type_t convert_sai_buffer_pool_type_t_to_sai(
    lemming::dataplane::sai::BufferPoolType val) {
  switch (val) {
    case lemming::dataplane::sai::BUFFER_POOL_TYPE_INGRESS:
      return SAI_BUFFER_POOL_TYPE_INGRESS;

    case lemming::dataplane::sai::BUFFER_POOL_TYPE_EGRESS:
      return SAI_BUFFER_POOL_TYPE_EGRESS;

    case lemming::dataplane::sai::BUFFER_POOL_TYPE_BOTH:
      return SAI_BUFFER_POOL_TYPE_BOTH;

    default:
      return SAI_BUFFER_POOL_TYPE_INGRESS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_buffer_pool_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_buffer_pool_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_buffer_pool_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_buffer_pool_type_t_to_sai(
        static_cast<lemming::dataplane::sai::BufferPoolType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BufferProfileAttr
convert_sai_buffer_profile_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_BUFFER_PROFILE_ATTR_POOL_ID:
      return lemming::dataplane::sai::BUFFER_PROFILE_ATTR_POOL_ID;

    case SAI_BUFFER_PROFILE_ATTR_RESERVED_BUFFER_SIZE:
      return lemming::dataplane::sai::BUFFER_PROFILE_ATTR_RESERVED_BUFFER_SIZE;

    case SAI_BUFFER_PROFILE_ATTR_THRESHOLD_MODE:
      return lemming::dataplane::sai::BUFFER_PROFILE_ATTR_THRESHOLD_MODE;

    case SAI_BUFFER_PROFILE_ATTR_SHARED_DYNAMIC_TH:
      return lemming::dataplane::sai::BUFFER_PROFILE_ATTR_SHARED_DYNAMIC_TH;

    case SAI_BUFFER_PROFILE_ATTR_SHARED_STATIC_TH:
      return lemming::dataplane::sai::BUFFER_PROFILE_ATTR_SHARED_STATIC_TH;

    case SAI_BUFFER_PROFILE_ATTR_XOFF_TH:
      return lemming::dataplane::sai::BUFFER_PROFILE_ATTR_XOFF_TH;

    case SAI_BUFFER_PROFILE_ATTR_XON_TH:
      return lemming::dataplane::sai::BUFFER_PROFILE_ATTR_XON_TH;

    case SAI_BUFFER_PROFILE_ATTR_XON_OFFSET_TH:
      return lemming::dataplane::sai::BUFFER_PROFILE_ATTR_XON_OFFSET_TH;

    default:
      return lemming::dataplane::sai::BUFFER_PROFILE_ATTR_UNSPECIFIED;
  }
}
sai_buffer_profile_attr_t convert_sai_buffer_profile_attr_t_to_sai(
    lemming::dataplane::sai::BufferProfileAttr val) {
  switch (val) {
    case lemming::dataplane::sai::BUFFER_PROFILE_ATTR_POOL_ID:
      return SAI_BUFFER_PROFILE_ATTR_POOL_ID;

    case lemming::dataplane::sai::BUFFER_PROFILE_ATTR_RESERVED_BUFFER_SIZE:
      return SAI_BUFFER_PROFILE_ATTR_RESERVED_BUFFER_SIZE;

    case lemming::dataplane::sai::BUFFER_PROFILE_ATTR_THRESHOLD_MODE:
      return SAI_BUFFER_PROFILE_ATTR_THRESHOLD_MODE;

    case lemming::dataplane::sai::BUFFER_PROFILE_ATTR_SHARED_DYNAMIC_TH:
      return SAI_BUFFER_PROFILE_ATTR_SHARED_DYNAMIC_TH;

    case lemming::dataplane::sai::BUFFER_PROFILE_ATTR_SHARED_STATIC_TH:
      return SAI_BUFFER_PROFILE_ATTR_SHARED_STATIC_TH;

    case lemming::dataplane::sai::BUFFER_PROFILE_ATTR_XOFF_TH:
      return SAI_BUFFER_PROFILE_ATTR_XOFF_TH;

    case lemming::dataplane::sai::BUFFER_PROFILE_ATTR_XON_TH:
      return SAI_BUFFER_PROFILE_ATTR_XON_TH;

    case lemming::dataplane::sai::BUFFER_PROFILE_ATTR_XON_OFFSET_TH:
      return SAI_BUFFER_PROFILE_ATTR_XON_OFFSET_TH;

    default:
      return SAI_BUFFER_PROFILE_ATTR_POOL_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_buffer_profile_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_buffer_profile_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_buffer_profile_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_buffer_profile_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::BufferProfileAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BufferProfileThresholdMode
convert_sai_buffer_profile_threshold_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_BUFFER_PROFILE_THRESHOLD_MODE_STATIC:
      return lemming::dataplane::sai::BUFFER_PROFILE_THRESHOLD_MODE_STATIC;

    case SAI_BUFFER_PROFILE_THRESHOLD_MODE_DYNAMIC:
      return lemming::dataplane::sai::BUFFER_PROFILE_THRESHOLD_MODE_DYNAMIC;

    default:
      return lemming::dataplane::sai::BUFFER_PROFILE_THRESHOLD_MODE_UNSPECIFIED;
  }
}
sai_buffer_profile_threshold_mode_t
convert_sai_buffer_profile_threshold_mode_t_to_sai(
    lemming::dataplane::sai::BufferProfileThresholdMode val) {
  switch (val) {
    case lemming::dataplane::sai::BUFFER_PROFILE_THRESHOLD_MODE_STATIC:
      return SAI_BUFFER_PROFILE_THRESHOLD_MODE_STATIC;

    case lemming::dataplane::sai::BUFFER_PROFILE_THRESHOLD_MODE_DYNAMIC:
      return SAI_BUFFER_PROFILE_THRESHOLD_MODE_DYNAMIC;

    default:
      return SAI_BUFFER_PROFILE_THRESHOLD_MODE_STATIC;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_buffer_profile_threshold_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_buffer_profile_threshold_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_buffer_profile_threshold_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_buffer_profile_threshold_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::BufferProfileThresholdMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::BulkOpErrorMode
convert_sai_bulk_op_error_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_BULK_OP_ERROR_MODE_STOP_ON_ERROR:
      return lemming::dataplane::sai::BULK_OP_ERROR_MODE_STOP_ON_ERROR;

    case SAI_BULK_OP_ERROR_MODE_IGNORE_ERROR:
      return lemming::dataplane::sai::BULK_OP_ERROR_MODE_IGNORE_ERROR;

    default:
      return lemming::dataplane::sai::BULK_OP_ERROR_MODE_UNSPECIFIED;
  }
}
sai_bulk_op_error_mode_t convert_sai_bulk_op_error_mode_t_to_sai(
    lemming::dataplane::sai::BulkOpErrorMode val) {
  switch (val) {
    case lemming::dataplane::sai::BULK_OP_ERROR_MODE_STOP_ON_ERROR:
      return SAI_BULK_OP_ERROR_MODE_STOP_ON_ERROR;

    case lemming::dataplane::sai::BULK_OP_ERROR_MODE_IGNORE_ERROR:
      return SAI_BULK_OP_ERROR_MODE_IGNORE_ERROR;

    default:
      return SAI_BULK_OP_ERROR_MODE_STOP_ON_ERROR;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_bulk_op_error_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_bulk_op_error_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_bulk_op_error_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_bulk_op_error_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::BulkOpErrorMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::CommonApi convert_sai_common_api_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_COMMON_API_CREATE:
      return lemming::dataplane::sai::COMMON_API_CREATE;

    case SAI_COMMON_API_REMOVE:
      return lemming::dataplane::sai::COMMON_API_REMOVE;

    case SAI_COMMON_API_SET:
      return lemming::dataplane::sai::COMMON_API_SET;

    case SAI_COMMON_API_GET:
      return lemming::dataplane::sai::COMMON_API_GET;

    case SAI_COMMON_API_BULK_CREATE:
      return lemming::dataplane::sai::COMMON_API_BULK_CREATE;

    case SAI_COMMON_API_BULK_REMOVE:
      return lemming::dataplane::sai::COMMON_API_BULK_REMOVE;

    case SAI_COMMON_API_BULK_SET:
      return lemming::dataplane::sai::COMMON_API_BULK_SET;

    case SAI_COMMON_API_BULK_GET:
      return lemming::dataplane::sai::COMMON_API_BULK_GET;

    case SAI_COMMON_API_MAX:
      return lemming::dataplane::sai::COMMON_API_MAX;

    default:
      return lemming::dataplane::sai::COMMON_API_UNSPECIFIED;
  }
}
sai_common_api_t convert_sai_common_api_t_to_sai(
    lemming::dataplane::sai::CommonApi val) {
  switch (val) {
    case lemming::dataplane::sai::COMMON_API_CREATE:
      return SAI_COMMON_API_CREATE;

    case lemming::dataplane::sai::COMMON_API_REMOVE:
      return SAI_COMMON_API_REMOVE;

    case lemming::dataplane::sai::COMMON_API_SET:
      return SAI_COMMON_API_SET;

    case lemming::dataplane::sai::COMMON_API_GET:
      return SAI_COMMON_API_GET;

    case lemming::dataplane::sai::COMMON_API_BULK_CREATE:
      return SAI_COMMON_API_BULK_CREATE;

    case lemming::dataplane::sai::COMMON_API_BULK_REMOVE:
      return SAI_COMMON_API_BULK_REMOVE;

    case lemming::dataplane::sai::COMMON_API_BULK_SET:
      return SAI_COMMON_API_BULK_SET;

    case lemming::dataplane::sai::COMMON_API_BULK_GET:
      return SAI_COMMON_API_BULK_GET;

    case lemming::dataplane::sai::COMMON_API_MAX:
      return SAI_COMMON_API_MAX;

    default:
      return SAI_COMMON_API_CREATE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_common_api_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_common_api_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_common_api_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_common_api_t_to_sai(
        static_cast<lemming::dataplane::sai::CommonApi>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::CounterAttr convert_sai_counter_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_COUNTER_ATTR_TYPE:
      return lemming::dataplane::sai::COUNTER_ATTR_TYPE;

    case SAI_COUNTER_ATTR_LABEL:
      return lemming::dataplane::sai::COUNTER_ATTR_LABEL;

    default:
      return lemming::dataplane::sai::COUNTER_ATTR_UNSPECIFIED;
  }
}
sai_counter_attr_t convert_sai_counter_attr_t_to_sai(
    lemming::dataplane::sai::CounterAttr val) {
  switch (val) {
    case lemming::dataplane::sai::COUNTER_ATTR_TYPE:
      return SAI_COUNTER_ATTR_TYPE;

    case lemming::dataplane::sai::COUNTER_ATTR_LABEL:
      return SAI_COUNTER_ATTR_LABEL;

    default:
      return SAI_COUNTER_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_counter_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_counter_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_counter_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_counter_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::CounterAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::CounterStat convert_sai_counter_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_COUNTER_STAT_PACKETS:
      return lemming::dataplane::sai::COUNTER_STAT_PACKETS;

    case SAI_COUNTER_STAT_BYTES:
      return lemming::dataplane::sai::COUNTER_STAT_BYTES;

    case SAI_COUNTER_STAT_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::COUNTER_STAT_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::COUNTER_STAT_UNSPECIFIED;
  }
}
sai_counter_stat_t convert_sai_counter_stat_t_to_sai(
    lemming::dataplane::sai::CounterStat val) {
  switch (val) {
    case lemming::dataplane::sai::COUNTER_STAT_PACKETS:
      return SAI_COUNTER_STAT_PACKETS;

    case lemming::dataplane::sai::COUNTER_STAT_BYTES:
      return SAI_COUNTER_STAT_BYTES;

    case lemming::dataplane::sai::COUNTER_STAT_CUSTOM_RANGE_BASE:
      return SAI_COUNTER_STAT_CUSTOM_RANGE_BASE;

    default:
      return SAI_COUNTER_STAT_PACKETS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_counter_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_counter_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_counter_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_counter_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::CounterStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::CounterType convert_sai_counter_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_COUNTER_TYPE_REGULAR:
      return lemming::dataplane::sai::COUNTER_TYPE_REGULAR;

    default:
      return lemming::dataplane::sai::COUNTER_TYPE_UNSPECIFIED;
  }
}
sai_counter_type_t convert_sai_counter_type_t_to_sai(
    lemming::dataplane::sai::CounterType val) {
  switch (val) {
    case lemming::dataplane::sai::COUNTER_TYPE_REGULAR:
      return SAI_COUNTER_TYPE_REGULAR;

    default:
      return SAI_COUNTER_TYPE_REGULAR;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_counter_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_counter_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_counter_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_counter_type_t_to_sai(
        static_cast<lemming::dataplane::sai::CounterType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::DebugCounterAttr
convert_sai_debug_counter_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_DEBUG_COUNTER_ATTR_INDEX:
      return lemming::dataplane::sai::DEBUG_COUNTER_ATTR_INDEX;

    case SAI_DEBUG_COUNTER_ATTR_TYPE:
      return lemming::dataplane::sai::DEBUG_COUNTER_ATTR_TYPE;

    case SAI_DEBUG_COUNTER_ATTR_BIND_METHOD:
      return lemming::dataplane::sai::DEBUG_COUNTER_ATTR_BIND_METHOD;

    case SAI_DEBUG_COUNTER_ATTR_IN_DROP_REASON_LIST:
      return lemming::dataplane::sai::DEBUG_COUNTER_ATTR_IN_DROP_REASON_LIST;

    case SAI_DEBUG_COUNTER_ATTR_OUT_DROP_REASON_LIST:
      return lemming::dataplane::sai::DEBUG_COUNTER_ATTR_OUT_DROP_REASON_LIST;

    default:
      return lemming::dataplane::sai::DEBUG_COUNTER_ATTR_UNSPECIFIED;
  }
}
sai_debug_counter_attr_t convert_sai_debug_counter_attr_t_to_sai(
    lemming::dataplane::sai::DebugCounterAttr val) {
  switch (val) {
    case lemming::dataplane::sai::DEBUG_COUNTER_ATTR_INDEX:
      return SAI_DEBUG_COUNTER_ATTR_INDEX;

    case lemming::dataplane::sai::DEBUG_COUNTER_ATTR_TYPE:
      return SAI_DEBUG_COUNTER_ATTR_TYPE;

    case lemming::dataplane::sai::DEBUG_COUNTER_ATTR_BIND_METHOD:
      return SAI_DEBUG_COUNTER_ATTR_BIND_METHOD;

    case lemming::dataplane::sai::DEBUG_COUNTER_ATTR_IN_DROP_REASON_LIST:
      return SAI_DEBUG_COUNTER_ATTR_IN_DROP_REASON_LIST;

    case lemming::dataplane::sai::DEBUG_COUNTER_ATTR_OUT_DROP_REASON_LIST:
      return SAI_DEBUG_COUNTER_ATTR_OUT_DROP_REASON_LIST;

    default:
      return SAI_DEBUG_COUNTER_ATTR_INDEX;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_debug_counter_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_debug_counter_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_debug_counter_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_debug_counter_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::DebugCounterAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::DebugCounterBindMethod
convert_sai_debug_counter_bind_method_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_DEBUG_COUNTER_BIND_METHOD_AUTOMATIC:
      return lemming::dataplane::sai::DEBUG_COUNTER_BIND_METHOD_AUTOMATIC;

    default:
      return lemming::dataplane::sai::DEBUG_COUNTER_BIND_METHOD_UNSPECIFIED;
  }
}
sai_debug_counter_bind_method_t convert_sai_debug_counter_bind_method_t_to_sai(
    lemming::dataplane::sai::DebugCounterBindMethod val) {
  switch (val) {
    case lemming::dataplane::sai::DEBUG_COUNTER_BIND_METHOD_AUTOMATIC:
      return SAI_DEBUG_COUNTER_BIND_METHOD_AUTOMATIC;

    default:
      return SAI_DEBUG_COUNTER_BIND_METHOD_AUTOMATIC;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_debug_counter_bind_method_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_debug_counter_bind_method_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_debug_counter_bind_method_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_debug_counter_bind_method_t_to_sai(
        static_cast<lemming::dataplane::sai::DebugCounterBindMethod>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::DebugCounterType
convert_sai_debug_counter_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_DEBUG_COUNTER_TYPE_PORT_IN_DROP_REASONS:
      return lemming::dataplane::sai::DEBUG_COUNTER_TYPE_PORT_IN_DROP_REASONS;

    case SAI_DEBUG_COUNTER_TYPE_PORT_OUT_DROP_REASONS:
      return lemming::dataplane::sai::DEBUG_COUNTER_TYPE_PORT_OUT_DROP_REASONS;

    case SAI_DEBUG_COUNTER_TYPE_SWITCH_IN_DROP_REASONS:
      return lemming::dataplane::sai::DEBUG_COUNTER_TYPE_SWITCH_IN_DROP_REASONS;

    case SAI_DEBUG_COUNTER_TYPE_SWITCH_OUT_DROP_REASONS:
      return lemming::dataplane::sai::
          DEBUG_COUNTER_TYPE_SWITCH_OUT_DROP_REASONS;

    default:
      return lemming::dataplane::sai::DEBUG_COUNTER_TYPE_UNSPECIFIED;
  }
}
sai_debug_counter_type_t convert_sai_debug_counter_type_t_to_sai(
    lemming::dataplane::sai::DebugCounterType val) {
  switch (val) {
    case lemming::dataplane::sai::DEBUG_COUNTER_TYPE_PORT_IN_DROP_REASONS:
      return SAI_DEBUG_COUNTER_TYPE_PORT_IN_DROP_REASONS;

    case lemming::dataplane::sai::DEBUG_COUNTER_TYPE_PORT_OUT_DROP_REASONS:
      return SAI_DEBUG_COUNTER_TYPE_PORT_OUT_DROP_REASONS;

    case lemming::dataplane::sai::DEBUG_COUNTER_TYPE_SWITCH_IN_DROP_REASONS:
      return SAI_DEBUG_COUNTER_TYPE_SWITCH_IN_DROP_REASONS;

    case lemming::dataplane::sai::DEBUG_COUNTER_TYPE_SWITCH_OUT_DROP_REASONS:
      return SAI_DEBUG_COUNTER_TYPE_SWITCH_OUT_DROP_REASONS;

    default:
      return SAI_DEBUG_COUNTER_TYPE_PORT_IN_DROP_REASONS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_debug_counter_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_debug_counter_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_debug_counter_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_debug_counter_type_t_to_sai(
        static_cast<lemming::dataplane::sai::DebugCounterType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::DtelAttr convert_sai_dtel_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_DTEL_ATTR_INT_ENDPOINT_ENABLE:
      return lemming::dataplane::sai::DTEL_ATTR_INT_ENDPOINT_ENABLE;

    case SAI_DTEL_ATTR_INT_TRANSIT_ENABLE:
      return lemming::dataplane::sai::DTEL_ATTR_INT_TRANSIT_ENABLE;

    case SAI_DTEL_ATTR_POSTCARD_ENABLE:
      return lemming::dataplane::sai::DTEL_ATTR_POSTCARD_ENABLE;

    case SAI_DTEL_ATTR_DROP_REPORT_ENABLE:
      return lemming::dataplane::sai::DTEL_ATTR_DROP_REPORT_ENABLE;

    case SAI_DTEL_ATTR_QUEUE_REPORT_ENABLE:
      return lemming::dataplane::sai::DTEL_ATTR_QUEUE_REPORT_ENABLE;

    case SAI_DTEL_ATTR_SWITCH_ID:
      return lemming::dataplane::sai::DTEL_ATTR_SWITCH_ID;

    case SAI_DTEL_ATTR_FLOW_STATE_CLEAR_CYCLE:
      return lemming::dataplane::sai::DTEL_ATTR_FLOW_STATE_CLEAR_CYCLE;

    case SAI_DTEL_ATTR_LATENCY_SENSITIVITY:
      return lemming::dataplane::sai::DTEL_ATTR_LATENCY_SENSITIVITY;

    case SAI_DTEL_ATTR_SINK_PORT_LIST:
      return lemming::dataplane::sai::DTEL_ATTR_SINK_PORT_LIST;

    case SAI_DTEL_ATTR_INT_L4_DSCP:
      return lemming::dataplane::sai::DTEL_ATTR_INT_L4_DSCP;

    default:
      return lemming::dataplane::sai::DTEL_ATTR_UNSPECIFIED;
  }
}
sai_dtel_attr_t convert_sai_dtel_attr_t_to_sai(
    lemming::dataplane::sai::DtelAttr val) {
  switch (val) {
    case lemming::dataplane::sai::DTEL_ATTR_INT_ENDPOINT_ENABLE:
      return SAI_DTEL_ATTR_INT_ENDPOINT_ENABLE;

    case lemming::dataplane::sai::DTEL_ATTR_INT_TRANSIT_ENABLE:
      return SAI_DTEL_ATTR_INT_TRANSIT_ENABLE;

    case lemming::dataplane::sai::DTEL_ATTR_POSTCARD_ENABLE:
      return SAI_DTEL_ATTR_POSTCARD_ENABLE;

    case lemming::dataplane::sai::DTEL_ATTR_DROP_REPORT_ENABLE:
      return SAI_DTEL_ATTR_DROP_REPORT_ENABLE;

    case lemming::dataplane::sai::DTEL_ATTR_QUEUE_REPORT_ENABLE:
      return SAI_DTEL_ATTR_QUEUE_REPORT_ENABLE;

    case lemming::dataplane::sai::DTEL_ATTR_SWITCH_ID:
      return SAI_DTEL_ATTR_SWITCH_ID;

    case lemming::dataplane::sai::DTEL_ATTR_FLOW_STATE_CLEAR_CYCLE:
      return SAI_DTEL_ATTR_FLOW_STATE_CLEAR_CYCLE;

    case lemming::dataplane::sai::DTEL_ATTR_LATENCY_SENSITIVITY:
      return SAI_DTEL_ATTR_LATENCY_SENSITIVITY;

    case lemming::dataplane::sai::DTEL_ATTR_SINK_PORT_LIST:
      return SAI_DTEL_ATTR_SINK_PORT_LIST;

    case lemming::dataplane::sai::DTEL_ATTR_INT_L4_DSCP:
      return SAI_DTEL_ATTR_INT_L4_DSCP;

    default:
      return SAI_DTEL_ATTR_INT_ENDPOINT_ENABLE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_dtel_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_dtel_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_dtel_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_dtel_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::DtelAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::DtelEventAttr convert_sai_dtel_event_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_DTEL_EVENT_ATTR_TYPE:
      return lemming::dataplane::sai::DTEL_EVENT_ATTR_TYPE;

    case SAI_DTEL_EVENT_ATTR_REPORT_SESSION:
      return lemming::dataplane::sai::DTEL_EVENT_ATTR_REPORT_SESSION;

    case SAI_DTEL_EVENT_ATTR_DSCP_VALUE:
      return lemming::dataplane::sai::DTEL_EVENT_ATTR_DSCP_VALUE;

    default:
      return lemming::dataplane::sai::DTEL_EVENT_ATTR_UNSPECIFIED;
  }
}
sai_dtel_event_attr_t convert_sai_dtel_event_attr_t_to_sai(
    lemming::dataplane::sai::DtelEventAttr val) {
  switch (val) {
    case lemming::dataplane::sai::DTEL_EVENT_ATTR_TYPE:
      return SAI_DTEL_EVENT_ATTR_TYPE;

    case lemming::dataplane::sai::DTEL_EVENT_ATTR_REPORT_SESSION:
      return SAI_DTEL_EVENT_ATTR_REPORT_SESSION;

    case lemming::dataplane::sai::DTEL_EVENT_ATTR_DSCP_VALUE:
      return SAI_DTEL_EVENT_ATTR_DSCP_VALUE;

    default:
      return SAI_DTEL_EVENT_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_dtel_event_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_dtel_event_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_dtel_event_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_dtel_event_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::DtelEventAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::DtelEventType convert_sai_dtel_event_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_DTEL_EVENT_TYPE_FLOW_STATE:
      return lemming::dataplane::sai::DTEL_EVENT_TYPE_FLOW_STATE;

    case SAI_DTEL_EVENT_TYPE_FLOW_REPORT_ALL_PACKETS:
      return lemming::dataplane::sai::DTEL_EVENT_TYPE_FLOW_REPORT_ALL_PACKETS;

    case SAI_DTEL_EVENT_TYPE_FLOW_TCPFLAG:
      return lemming::dataplane::sai::DTEL_EVENT_TYPE_FLOW_TCPFLAG;

    case SAI_DTEL_EVENT_TYPE_QUEUE_REPORT_THRESHOLD_BREACH:
      return lemming::dataplane::sai::
          DTEL_EVENT_TYPE_QUEUE_REPORT_THRESHOLD_BREACH;

    case SAI_DTEL_EVENT_TYPE_QUEUE_REPORT_TAIL_DROP:
      return lemming::dataplane::sai::DTEL_EVENT_TYPE_QUEUE_REPORT_TAIL_DROP;

    case SAI_DTEL_EVENT_TYPE_DROP_REPORT:
      return lemming::dataplane::sai::DTEL_EVENT_TYPE_DROP_REPORT;

    case SAI_DTEL_EVENT_TYPE_MAX:
      return lemming::dataplane::sai::DTEL_EVENT_TYPE_MAX;

    default:
      return lemming::dataplane::sai::DTEL_EVENT_TYPE_UNSPECIFIED;
  }
}
sai_dtel_event_type_t convert_sai_dtel_event_type_t_to_sai(
    lemming::dataplane::sai::DtelEventType val) {
  switch (val) {
    case lemming::dataplane::sai::DTEL_EVENT_TYPE_FLOW_STATE:
      return SAI_DTEL_EVENT_TYPE_FLOW_STATE;

    case lemming::dataplane::sai::DTEL_EVENT_TYPE_FLOW_REPORT_ALL_PACKETS:
      return SAI_DTEL_EVENT_TYPE_FLOW_REPORT_ALL_PACKETS;

    case lemming::dataplane::sai::DTEL_EVENT_TYPE_FLOW_TCPFLAG:
      return SAI_DTEL_EVENT_TYPE_FLOW_TCPFLAG;

    case lemming::dataplane::sai::DTEL_EVENT_TYPE_QUEUE_REPORT_THRESHOLD_BREACH:
      return SAI_DTEL_EVENT_TYPE_QUEUE_REPORT_THRESHOLD_BREACH;

    case lemming::dataplane::sai::DTEL_EVENT_TYPE_QUEUE_REPORT_TAIL_DROP:
      return SAI_DTEL_EVENT_TYPE_QUEUE_REPORT_TAIL_DROP;

    case lemming::dataplane::sai::DTEL_EVENT_TYPE_DROP_REPORT:
      return SAI_DTEL_EVENT_TYPE_DROP_REPORT;

    case lemming::dataplane::sai::DTEL_EVENT_TYPE_MAX:
      return SAI_DTEL_EVENT_TYPE_MAX;

    default:
      return SAI_DTEL_EVENT_TYPE_FLOW_STATE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_dtel_event_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_dtel_event_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_dtel_event_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_dtel_event_type_t_to_sai(
        static_cast<lemming::dataplane::sai::DtelEventType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::DtelIntSessionAttr
convert_sai_dtel_int_session_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_DTEL_INT_SESSION_ATTR_MAX_HOP_COUNT:
      return lemming::dataplane::sai::DTEL_INT_SESSION_ATTR_MAX_HOP_COUNT;

    case SAI_DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_ID:
      return lemming::dataplane::sai::DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_ID;

    case SAI_DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_PORTS:
      return lemming::dataplane::sai::
          DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_PORTS;

    case SAI_DTEL_INT_SESSION_ATTR_COLLECT_INGRESS_TIMESTAMP:
      return lemming::dataplane::sai::
          DTEL_INT_SESSION_ATTR_COLLECT_INGRESS_TIMESTAMP;

    case SAI_DTEL_INT_SESSION_ATTR_COLLECT_EGRESS_TIMESTAMP:
      return lemming::dataplane::sai::
          DTEL_INT_SESSION_ATTR_COLLECT_EGRESS_TIMESTAMP;

    case SAI_DTEL_INT_SESSION_ATTR_COLLECT_QUEUE_INFO:
      return lemming::dataplane::sai::DTEL_INT_SESSION_ATTR_COLLECT_QUEUE_INFO;

    default:
      return lemming::dataplane::sai::DTEL_INT_SESSION_ATTR_UNSPECIFIED;
  }
}
sai_dtel_int_session_attr_t convert_sai_dtel_int_session_attr_t_to_sai(
    lemming::dataplane::sai::DtelIntSessionAttr val) {
  switch (val) {
    case lemming::dataplane::sai::DTEL_INT_SESSION_ATTR_MAX_HOP_COUNT:
      return SAI_DTEL_INT_SESSION_ATTR_MAX_HOP_COUNT;

    case lemming::dataplane::sai::DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_ID:
      return SAI_DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_ID;

    case lemming::dataplane::sai::DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_PORTS:
      return SAI_DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_PORTS;

    case lemming::dataplane::sai::
        DTEL_INT_SESSION_ATTR_COLLECT_INGRESS_TIMESTAMP:
      return SAI_DTEL_INT_SESSION_ATTR_COLLECT_INGRESS_TIMESTAMP;

    case lemming::dataplane::sai::
        DTEL_INT_SESSION_ATTR_COLLECT_EGRESS_TIMESTAMP:
      return SAI_DTEL_INT_SESSION_ATTR_COLLECT_EGRESS_TIMESTAMP;

    case lemming::dataplane::sai::DTEL_INT_SESSION_ATTR_COLLECT_QUEUE_INFO:
      return SAI_DTEL_INT_SESSION_ATTR_COLLECT_QUEUE_INFO;

    default:
      return SAI_DTEL_INT_SESSION_ATTR_MAX_HOP_COUNT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_dtel_int_session_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_dtel_int_session_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_dtel_int_session_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_dtel_int_session_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::DtelIntSessionAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::DtelQueueReportAttr
convert_sai_dtel_queue_report_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_DTEL_QUEUE_REPORT_ATTR_QUEUE_ID:
      return lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_QUEUE_ID;

    case SAI_DTEL_QUEUE_REPORT_ATTR_DEPTH_THRESHOLD:
      return lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_DEPTH_THRESHOLD;

    case SAI_DTEL_QUEUE_REPORT_ATTR_LATENCY_THRESHOLD:
      return lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_LATENCY_THRESHOLD;

    case SAI_DTEL_QUEUE_REPORT_ATTR_BREACH_QUOTA:
      return lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_BREACH_QUOTA;

    case SAI_DTEL_QUEUE_REPORT_ATTR_TAIL_DROP:
      return lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_TAIL_DROP;

    default:
      return lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_UNSPECIFIED;
  }
}
sai_dtel_queue_report_attr_t convert_sai_dtel_queue_report_attr_t_to_sai(
    lemming::dataplane::sai::DtelQueueReportAttr val) {
  switch (val) {
    case lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_QUEUE_ID:
      return SAI_DTEL_QUEUE_REPORT_ATTR_QUEUE_ID;

    case lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_DEPTH_THRESHOLD:
      return SAI_DTEL_QUEUE_REPORT_ATTR_DEPTH_THRESHOLD;

    case lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_LATENCY_THRESHOLD:
      return SAI_DTEL_QUEUE_REPORT_ATTR_LATENCY_THRESHOLD;

    case lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_BREACH_QUOTA:
      return SAI_DTEL_QUEUE_REPORT_ATTR_BREACH_QUOTA;

    case lemming::dataplane::sai::DTEL_QUEUE_REPORT_ATTR_TAIL_DROP:
      return SAI_DTEL_QUEUE_REPORT_ATTR_TAIL_DROP;

    default:
      return SAI_DTEL_QUEUE_REPORT_ATTR_QUEUE_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_dtel_queue_report_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_dtel_queue_report_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_dtel_queue_report_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_dtel_queue_report_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::DtelQueueReportAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::DtelReportSessionAttr
convert_sai_dtel_report_session_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_DTEL_REPORT_SESSION_ATTR_SRC_IP:
      return lemming::dataplane::sai::DTEL_REPORT_SESSION_ATTR_SRC_IP;

    case SAI_DTEL_REPORT_SESSION_ATTR_DST_IP_LIST:
      return lemming::dataplane::sai::DTEL_REPORT_SESSION_ATTR_DST_IP_LIST;

    case SAI_DTEL_REPORT_SESSION_ATTR_VIRTUAL_ROUTER_ID:
      return lemming::dataplane::sai::
          DTEL_REPORT_SESSION_ATTR_VIRTUAL_ROUTER_ID;

    case SAI_DTEL_REPORT_SESSION_ATTR_TRUNCATE_SIZE:
      return lemming::dataplane::sai::DTEL_REPORT_SESSION_ATTR_TRUNCATE_SIZE;

    case SAI_DTEL_REPORT_SESSION_ATTR_UDP_DST_PORT:
      return lemming::dataplane::sai::DTEL_REPORT_SESSION_ATTR_UDP_DST_PORT;

    default:
      return lemming::dataplane::sai::DTEL_REPORT_SESSION_ATTR_UNSPECIFIED;
  }
}
sai_dtel_report_session_attr_t convert_sai_dtel_report_session_attr_t_to_sai(
    lemming::dataplane::sai::DtelReportSessionAttr val) {
  switch (val) {
    case lemming::dataplane::sai::DTEL_REPORT_SESSION_ATTR_SRC_IP:
      return SAI_DTEL_REPORT_SESSION_ATTR_SRC_IP;

    case lemming::dataplane::sai::DTEL_REPORT_SESSION_ATTR_DST_IP_LIST:
      return SAI_DTEL_REPORT_SESSION_ATTR_DST_IP_LIST;

    case lemming::dataplane::sai::DTEL_REPORT_SESSION_ATTR_VIRTUAL_ROUTER_ID:
      return SAI_DTEL_REPORT_SESSION_ATTR_VIRTUAL_ROUTER_ID;

    case lemming::dataplane::sai::DTEL_REPORT_SESSION_ATTR_TRUNCATE_SIZE:
      return SAI_DTEL_REPORT_SESSION_ATTR_TRUNCATE_SIZE;

    case lemming::dataplane::sai::DTEL_REPORT_SESSION_ATTR_UDP_DST_PORT:
      return SAI_DTEL_REPORT_SESSION_ATTR_UDP_DST_PORT;

    default:
      return SAI_DTEL_REPORT_SESSION_ATTR_SRC_IP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_dtel_report_session_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_dtel_report_session_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_dtel_report_session_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_dtel_report_session_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::DtelReportSessionAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::EcnMarkMode convert_sai_ecn_mark_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ECN_MARK_MODE_NONE:
      return lemming::dataplane::sai::ECN_MARK_MODE_NONE;

    case SAI_ECN_MARK_MODE_GREEN:
      return lemming::dataplane::sai::ECN_MARK_MODE_GREEN;

    case SAI_ECN_MARK_MODE_YELLOW:
      return lemming::dataplane::sai::ECN_MARK_MODE_YELLOW;

    case SAI_ECN_MARK_MODE_RED:
      return lemming::dataplane::sai::ECN_MARK_MODE_RED;

    case SAI_ECN_MARK_MODE_GREEN_YELLOW:
      return lemming::dataplane::sai::ECN_MARK_MODE_GREEN_YELLOW;

    case SAI_ECN_MARK_MODE_GREEN_RED:
      return lemming::dataplane::sai::ECN_MARK_MODE_GREEN_RED;

    case SAI_ECN_MARK_MODE_YELLOW_RED:
      return lemming::dataplane::sai::ECN_MARK_MODE_YELLOW_RED;

    case SAI_ECN_MARK_MODE_ALL:
      return lemming::dataplane::sai::ECN_MARK_MODE_ALL;

    default:
      return lemming::dataplane::sai::ECN_MARK_MODE_UNSPECIFIED;
  }
}
sai_ecn_mark_mode_t convert_sai_ecn_mark_mode_t_to_sai(
    lemming::dataplane::sai::EcnMarkMode val) {
  switch (val) {
    case lemming::dataplane::sai::ECN_MARK_MODE_NONE:
      return SAI_ECN_MARK_MODE_NONE;

    case lemming::dataplane::sai::ECN_MARK_MODE_GREEN:
      return SAI_ECN_MARK_MODE_GREEN;

    case lemming::dataplane::sai::ECN_MARK_MODE_YELLOW:
      return SAI_ECN_MARK_MODE_YELLOW;

    case lemming::dataplane::sai::ECN_MARK_MODE_RED:
      return SAI_ECN_MARK_MODE_RED;

    case lemming::dataplane::sai::ECN_MARK_MODE_GREEN_YELLOW:
      return SAI_ECN_MARK_MODE_GREEN_YELLOW;

    case lemming::dataplane::sai::ECN_MARK_MODE_GREEN_RED:
      return SAI_ECN_MARK_MODE_GREEN_RED;

    case lemming::dataplane::sai::ECN_MARK_MODE_YELLOW_RED:
      return SAI_ECN_MARK_MODE_YELLOW_RED;

    case lemming::dataplane::sai::ECN_MARK_MODE_ALL:
      return SAI_ECN_MARK_MODE_ALL;

    default:
      return SAI_ECN_MARK_MODE_NONE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_ecn_mark_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ecn_mark_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ecn_mark_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ecn_mark_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::EcnMarkMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::ErspanEncapsulationType
convert_sai_erspan_encapsulation_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ERSPAN_ENCAPSULATION_TYPE_MIRROR_L3_GRE_TUNNEL:
      return lemming::dataplane::sai::
          ERSPAN_ENCAPSULATION_TYPE_MIRROR_L3_GRE_TUNNEL;

    default:
      return lemming::dataplane::sai::ERSPAN_ENCAPSULATION_TYPE_UNSPECIFIED;
  }
}
sai_erspan_encapsulation_type_t convert_sai_erspan_encapsulation_type_t_to_sai(
    lemming::dataplane::sai::ErspanEncapsulationType val) {
  switch (val) {
    case lemming::dataplane::sai::
        ERSPAN_ENCAPSULATION_TYPE_MIRROR_L3_GRE_TUNNEL:
      return SAI_ERSPAN_ENCAPSULATION_TYPE_MIRROR_L3_GRE_TUNNEL;

    default:
      return SAI_ERSPAN_ENCAPSULATION_TYPE_MIRROR_L3_GRE_TUNNEL;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_erspan_encapsulation_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_erspan_encapsulation_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_erspan_encapsulation_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_erspan_encapsulation_type_t_to_sai(
        static_cast<lemming::dataplane::sai::ErspanEncapsulationType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::FdbEntryAttr convert_sai_fdb_entry_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_FDB_ENTRY_ATTR_TYPE:
      return lemming::dataplane::sai::FDB_ENTRY_ATTR_TYPE;

    case SAI_FDB_ENTRY_ATTR_PACKET_ACTION:
      return lemming::dataplane::sai::FDB_ENTRY_ATTR_PACKET_ACTION;

    case SAI_FDB_ENTRY_ATTR_USER_TRAP_ID:
      return lemming::dataplane::sai::FDB_ENTRY_ATTR_USER_TRAP_ID;

    case SAI_FDB_ENTRY_ATTR_BRIDGE_PORT_ID:
      return lemming::dataplane::sai::FDB_ENTRY_ATTR_BRIDGE_PORT_ID;

    case SAI_FDB_ENTRY_ATTR_META_DATA:
      return lemming::dataplane::sai::FDB_ENTRY_ATTR_META_DATA;

    case SAI_FDB_ENTRY_ATTR_ENDPOINT_IP:
      return lemming::dataplane::sai::FDB_ENTRY_ATTR_ENDPOINT_IP;

    case SAI_FDB_ENTRY_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::FDB_ENTRY_ATTR_COUNTER_ID;

    case SAI_FDB_ENTRY_ATTR_ALLOW_MAC_MOVE:
      return lemming::dataplane::sai::FDB_ENTRY_ATTR_ALLOW_MAC_MOVE;

    default:
      return lemming::dataplane::sai::FDB_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_fdb_entry_attr_t convert_sai_fdb_entry_attr_t_to_sai(
    lemming::dataplane::sai::FdbEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::FDB_ENTRY_ATTR_TYPE:
      return SAI_FDB_ENTRY_ATTR_TYPE;

    case lemming::dataplane::sai::FDB_ENTRY_ATTR_PACKET_ACTION:
      return SAI_FDB_ENTRY_ATTR_PACKET_ACTION;

    case lemming::dataplane::sai::FDB_ENTRY_ATTR_USER_TRAP_ID:
      return SAI_FDB_ENTRY_ATTR_USER_TRAP_ID;

    case lemming::dataplane::sai::FDB_ENTRY_ATTR_BRIDGE_PORT_ID:
      return SAI_FDB_ENTRY_ATTR_BRIDGE_PORT_ID;

    case lemming::dataplane::sai::FDB_ENTRY_ATTR_META_DATA:
      return SAI_FDB_ENTRY_ATTR_META_DATA;

    case lemming::dataplane::sai::FDB_ENTRY_ATTR_ENDPOINT_IP:
      return SAI_FDB_ENTRY_ATTR_ENDPOINT_IP;

    case lemming::dataplane::sai::FDB_ENTRY_ATTR_COUNTER_ID:
      return SAI_FDB_ENTRY_ATTR_COUNTER_ID;

    case lemming::dataplane::sai::FDB_ENTRY_ATTR_ALLOW_MAC_MOVE:
      return SAI_FDB_ENTRY_ATTR_ALLOW_MAC_MOVE;

    default:
      return SAI_FDB_ENTRY_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_fdb_entry_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_fdb_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_fdb_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_fdb_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::FdbEntryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::FdbEntryType convert_sai_fdb_entry_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_FDB_ENTRY_TYPE_DYNAMIC:
      return lemming::dataplane::sai::FDB_ENTRY_TYPE_DYNAMIC;

    case SAI_FDB_ENTRY_TYPE_STATIC:
      return lemming::dataplane::sai::FDB_ENTRY_TYPE_STATIC;

    default:
      return lemming::dataplane::sai::FDB_ENTRY_TYPE_UNSPECIFIED;
  }
}
sai_fdb_entry_type_t convert_sai_fdb_entry_type_t_to_sai(
    lemming::dataplane::sai::FdbEntryType val) {
  switch (val) {
    case lemming::dataplane::sai::FDB_ENTRY_TYPE_DYNAMIC:
      return SAI_FDB_ENTRY_TYPE_DYNAMIC;

    case lemming::dataplane::sai::FDB_ENTRY_TYPE_STATIC:
      return SAI_FDB_ENTRY_TYPE_STATIC;

    default:
      return SAI_FDB_ENTRY_TYPE_DYNAMIC;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_fdb_entry_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_fdb_entry_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_fdb_entry_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_fdb_entry_type_t_to_sai(
        static_cast<lemming::dataplane::sai::FdbEntryType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::FdbEvent convert_sai_fdb_event_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_FDB_EVENT_LEARNED:
      return lemming::dataplane::sai::FDB_EVENT_LEARNED;

    case SAI_FDB_EVENT_AGED:
      return lemming::dataplane::sai::FDB_EVENT_AGED;

    case SAI_FDB_EVENT_MOVE:
      return lemming::dataplane::sai::FDB_EVENT_MOVE;

    case SAI_FDB_EVENT_FLUSHED:
      return lemming::dataplane::sai::FDB_EVENT_FLUSHED;

    default:
      return lemming::dataplane::sai::FDB_EVENT_UNSPECIFIED;
  }
}
sai_fdb_event_t convert_sai_fdb_event_t_to_sai(
    lemming::dataplane::sai::FdbEvent val) {
  switch (val) {
    case lemming::dataplane::sai::FDB_EVENT_LEARNED:
      return SAI_FDB_EVENT_LEARNED;

    case lemming::dataplane::sai::FDB_EVENT_AGED:
      return SAI_FDB_EVENT_AGED;

    case lemming::dataplane::sai::FDB_EVENT_MOVE:
      return SAI_FDB_EVENT_MOVE;

    case lemming::dataplane::sai::FDB_EVENT_FLUSHED:
      return SAI_FDB_EVENT_FLUSHED;

    default:
      return SAI_FDB_EVENT_LEARNED;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_fdb_event_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_fdb_event_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_fdb_event_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_fdb_event_t_to_sai(
        static_cast<lemming::dataplane::sai::FdbEvent>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::FdbFlushEntryType
convert_sai_fdb_flush_entry_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_FDB_FLUSH_ENTRY_TYPE_DYNAMIC:
      return lemming::dataplane::sai::FDB_FLUSH_ENTRY_TYPE_DYNAMIC;

    case SAI_FDB_FLUSH_ENTRY_TYPE_STATIC:
      return lemming::dataplane::sai::FDB_FLUSH_ENTRY_TYPE_STATIC;

    case SAI_FDB_FLUSH_ENTRY_TYPE_ALL:
      return lemming::dataplane::sai::FDB_FLUSH_ENTRY_TYPE_ALL;

    default:
      return lemming::dataplane::sai::FDB_FLUSH_ENTRY_TYPE_UNSPECIFIED;
  }
}
sai_fdb_flush_entry_type_t convert_sai_fdb_flush_entry_type_t_to_sai(
    lemming::dataplane::sai::FdbFlushEntryType val) {
  switch (val) {
    case lemming::dataplane::sai::FDB_FLUSH_ENTRY_TYPE_DYNAMIC:
      return SAI_FDB_FLUSH_ENTRY_TYPE_DYNAMIC;

    case lemming::dataplane::sai::FDB_FLUSH_ENTRY_TYPE_STATIC:
      return SAI_FDB_FLUSH_ENTRY_TYPE_STATIC;

    case lemming::dataplane::sai::FDB_FLUSH_ENTRY_TYPE_ALL:
      return SAI_FDB_FLUSH_ENTRY_TYPE_ALL;

    default:
      return SAI_FDB_FLUSH_ENTRY_TYPE_DYNAMIC;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_fdb_flush_entry_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_fdb_flush_entry_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_fdb_flush_entry_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_fdb_flush_entry_type_t_to_sai(
        static_cast<lemming::dataplane::sai::FdbFlushEntryType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::FineGrainedHashFieldAttr
convert_sai_fine_grained_hash_field_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_FINE_GRAINED_HASH_FIELD_ATTR_NATIVE_HASH_FIELD:
      return lemming::dataplane::sai::
          FINE_GRAINED_HASH_FIELD_ATTR_NATIVE_HASH_FIELD;

    case SAI_FINE_GRAINED_HASH_FIELD_ATTR_IPV4_MASK:
      return lemming::dataplane::sai::FINE_GRAINED_HASH_FIELD_ATTR_IPV4_MASK;

    case SAI_FINE_GRAINED_HASH_FIELD_ATTR_IPV6_MASK:
      return lemming::dataplane::sai::FINE_GRAINED_HASH_FIELD_ATTR_IPV6_MASK;

    case SAI_FINE_GRAINED_HASH_FIELD_ATTR_SEQUENCE_ID:
      return lemming::dataplane::sai::FINE_GRAINED_HASH_FIELD_ATTR_SEQUENCE_ID;

    default:
      return lemming::dataplane::sai::FINE_GRAINED_HASH_FIELD_ATTR_UNSPECIFIED;
  }
}
sai_fine_grained_hash_field_attr_t
convert_sai_fine_grained_hash_field_attr_t_to_sai(
    lemming::dataplane::sai::FineGrainedHashFieldAttr val) {
  switch (val) {
    case lemming::dataplane::sai::
        FINE_GRAINED_HASH_FIELD_ATTR_NATIVE_HASH_FIELD:
      return SAI_FINE_GRAINED_HASH_FIELD_ATTR_NATIVE_HASH_FIELD;

    case lemming::dataplane::sai::FINE_GRAINED_HASH_FIELD_ATTR_IPV4_MASK:
      return SAI_FINE_GRAINED_HASH_FIELD_ATTR_IPV4_MASK;

    case lemming::dataplane::sai::FINE_GRAINED_HASH_FIELD_ATTR_IPV6_MASK:
      return SAI_FINE_GRAINED_HASH_FIELD_ATTR_IPV6_MASK;

    case lemming::dataplane::sai::FINE_GRAINED_HASH_FIELD_ATTR_SEQUENCE_ID:
      return SAI_FINE_GRAINED_HASH_FIELD_ATTR_SEQUENCE_ID;

    default:
      return SAI_FINE_GRAINED_HASH_FIELD_ATTR_NATIVE_HASH_FIELD;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_fine_grained_hash_field_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_fine_grained_hash_field_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_fine_grained_hash_field_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_fine_grained_hash_field_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::FineGrainedHashFieldAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::GenericProgrammableAttr
convert_sai_generic_programmable_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_GENERIC_PROGRAMMABLE_ATTR_OBJECT_NAME:
      return lemming::dataplane::sai::GENERIC_PROGRAMMABLE_ATTR_OBJECT_NAME;

    case SAI_GENERIC_PROGRAMMABLE_ATTR_ENTRY:
      return lemming::dataplane::sai::GENERIC_PROGRAMMABLE_ATTR_ENTRY;

    case SAI_GENERIC_PROGRAMMABLE_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::GENERIC_PROGRAMMABLE_ATTR_COUNTER_ID;

    default:
      return lemming::dataplane::sai::GENERIC_PROGRAMMABLE_ATTR_UNSPECIFIED;
  }
}
sai_generic_programmable_attr_t convert_sai_generic_programmable_attr_t_to_sai(
    lemming::dataplane::sai::GenericProgrammableAttr val) {
  switch (val) {
    case lemming::dataplane::sai::GENERIC_PROGRAMMABLE_ATTR_OBJECT_NAME:
      return SAI_GENERIC_PROGRAMMABLE_ATTR_OBJECT_NAME;

    case lemming::dataplane::sai::GENERIC_PROGRAMMABLE_ATTR_ENTRY:
      return SAI_GENERIC_PROGRAMMABLE_ATTR_ENTRY;

    case lemming::dataplane::sai::GENERIC_PROGRAMMABLE_ATTR_COUNTER_ID:
      return SAI_GENERIC_PROGRAMMABLE_ATTR_COUNTER_ID;

    default:
      return SAI_GENERIC_PROGRAMMABLE_ATTR_OBJECT_NAME;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_generic_programmable_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_generic_programmable_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_generic_programmable_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_generic_programmable_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::GenericProgrammableAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HashAlgorithm convert_sai_hash_algorithm_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_HASH_ALGORITHM_CRC:
      return lemming::dataplane::sai::HASH_ALGORITHM_CRC;

    case SAI_HASH_ALGORITHM_XOR:
      return lemming::dataplane::sai::HASH_ALGORITHM_XOR;

    case SAI_HASH_ALGORITHM_RANDOM:
      return lemming::dataplane::sai::HASH_ALGORITHM_RANDOM;

    case SAI_HASH_ALGORITHM_CRC_32LO:
      return lemming::dataplane::sai::HASH_ALGORITHM_CRC_32LO;

    case SAI_HASH_ALGORITHM_CRC_32HI:
      return lemming::dataplane::sai::HASH_ALGORITHM_CRC_32HI;

    case SAI_HASH_ALGORITHM_CRC_CCITT:
      return lemming::dataplane::sai::HASH_ALGORITHM_CRC_CCITT;

    case SAI_HASH_ALGORITHM_CRC_XOR:
      return lemming::dataplane::sai::HASH_ALGORITHM_CRC_XOR;

    default:
      return lemming::dataplane::sai::HASH_ALGORITHM_UNSPECIFIED;
  }
}
sai_hash_algorithm_t convert_sai_hash_algorithm_t_to_sai(
    lemming::dataplane::sai::HashAlgorithm val) {
  switch (val) {
    case lemming::dataplane::sai::HASH_ALGORITHM_CRC:
      return SAI_HASH_ALGORITHM_CRC;

    case lemming::dataplane::sai::HASH_ALGORITHM_XOR:
      return SAI_HASH_ALGORITHM_XOR;

    case lemming::dataplane::sai::HASH_ALGORITHM_RANDOM:
      return SAI_HASH_ALGORITHM_RANDOM;

    case lemming::dataplane::sai::HASH_ALGORITHM_CRC_32LO:
      return SAI_HASH_ALGORITHM_CRC_32LO;

    case lemming::dataplane::sai::HASH_ALGORITHM_CRC_32HI:
      return SAI_HASH_ALGORITHM_CRC_32HI;

    case lemming::dataplane::sai::HASH_ALGORITHM_CRC_CCITT:
      return SAI_HASH_ALGORITHM_CRC_CCITT;

    case lemming::dataplane::sai::HASH_ALGORITHM_CRC_XOR:
      return SAI_HASH_ALGORITHM_CRC_XOR;

    default:
      return SAI_HASH_ALGORITHM_CRC;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_hash_algorithm_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_hash_algorithm_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hash_algorithm_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hash_algorithm_t_to_sai(
        static_cast<lemming::dataplane::sai::HashAlgorithm>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HashAttr convert_sai_hash_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_HASH_ATTR_NATIVE_HASH_FIELD_LIST:
      return lemming::dataplane::sai::HASH_ATTR_NATIVE_HASH_FIELD_LIST;

    case SAI_HASH_ATTR_UDF_GROUP_LIST:
      return lemming::dataplane::sai::HASH_ATTR_UDF_GROUP_LIST;

    case SAI_HASH_ATTR_FINE_GRAINED_HASH_FIELD_LIST:
      return lemming::dataplane::sai::HASH_ATTR_FINE_GRAINED_HASH_FIELD_LIST;

    default:
      return lemming::dataplane::sai::HASH_ATTR_UNSPECIFIED;
  }
}
sai_hash_attr_t convert_sai_hash_attr_t_to_sai(
    lemming::dataplane::sai::HashAttr val) {
  switch (val) {
    case lemming::dataplane::sai::HASH_ATTR_NATIVE_HASH_FIELD_LIST:
      return SAI_HASH_ATTR_NATIVE_HASH_FIELD_LIST;

    case lemming::dataplane::sai::HASH_ATTR_UDF_GROUP_LIST:
      return SAI_HASH_ATTR_UDF_GROUP_LIST;

    case lemming::dataplane::sai::HASH_ATTR_FINE_GRAINED_HASH_FIELD_LIST:
      return SAI_HASH_ATTR_FINE_GRAINED_HASH_FIELD_LIST;

    default:
      return SAI_HASH_ATTR_NATIVE_HASH_FIELD_LIST;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_hash_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_hash_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hash_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hash_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::HashAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifAttr convert_sai_hostif_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_ATTR_TYPE:
      return lemming::dataplane::sai::HOSTIF_ATTR_TYPE;

    case SAI_HOSTIF_ATTR_OBJ_ID:
      return lemming::dataplane::sai::HOSTIF_ATTR_OBJ_ID;

    case SAI_HOSTIF_ATTR_NAME:
      return lemming::dataplane::sai::HOSTIF_ATTR_NAME;

    case SAI_HOSTIF_ATTR_OPER_STATUS:
      return lemming::dataplane::sai::HOSTIF_ATTR_OPER_STATUS;

    case SAI_HOSTIF_ATTR_QUEUE:
      return lemming::dataplane::sai::HOSTIF_ATTR_QUEUE;

    case SAI_HOSTIF_ATTR_VLAN_TAG:
      return lemming::dataplane::sai::HOSTIF_ATTR_VLAN_TAG;

    case SAI_HOSTIF_ATTR_GENETLINK_MCGRP_NAME:
      return lemming::dataplane::sai::HOSTIF_ATTR_GENETLINK_MCGRP_NAME;

    default:
      return lemming::dataplane::sai::HOSTIF_ATTR_UNSPECIFIED;
  }
}
sai_hostif_attr_t convert_sai_hostif_attr_t_to_sai(
    lemming::dataplane::sai::HostifAttr val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_ATTR_TYPE:
      return SAI_HOSTIF_ATTR_TYPE;

    case lemming::dataplane::sai::HOSTIF_ATTR_OBJ_ID:
      return SAI_HOSTIF_ATTR_OBJ_ID;

    case lemming::dataplane::sai::HOSTIF_ATTR_NAME:
      return SAI_HOSTIF_ATTR_NAME;

    case lemming::dataplane::sai::HOSTIF_ATTR_OPER_STATUS:
      return SAI_HOSTIF_ATTR_OPER_STATUS;

    case lemming::dataplane::sai::HOSTIF_ATTR_QUEUE:
      return SAI_HOSTIF_ATTR_QUEUE;

    case lemming::dataplane::sai::HOSTIF_ATTR_VLAN_TAG:
      return SAI_HOSTIF_ATTR_VLAN_TAG;

    case lemming::dataplane::sai::HOSTIF_ATTR_GENETLINK_MCGRP_NAME:
      return SAI_HOSTIF_ATTR_GENETLINK_MCGRP_NAME;

    default:
      return SAI_HOSTIF_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_hostif_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_hostif_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifTableEntryAttr
convert_sai_hostif_table_entry_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_TABLE_ENTRY_ATTR_TYPE:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_TYPE;

    case SAI_HOSTIF_TABLE_ENTRY_ATTR_OBJ_ID:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_OBJ_ID;

    case SAI_HOSTIF_TABLE_ENTRY_ATTR_TRAP_ID:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_TRAP_ID;

    case SAI_HOSTIF_TABLE_ENTRY_ATTR_CHANNEL_TYPE:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_CHANNEL_TYPE;

    case SAI_HOSTIF_TABLE_ENTRY_ATTR_HOST_IF:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_HOST_IF;

    default:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_hostif_table_entry_attr_t convert_sai_hostif_table_entry_attr_t_to_sai(
    lemming::dataplane::sai::HostifTableEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_TYPE:
      return SAI_HOSTIF_TABLE_ENTRY_ATTR_TYPE;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_OBJ_ID:
      return SAI_HOSTIF_TABLE_ENTRY_ATTR_OBJ_ID;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_TRAP_ID:
      return SAI_HOSTIF_TABLE_ENTRY_ATTR_TRAP_ID;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_CHANNEL_TYPE:
      return SAI_HOSTIF_TABLE_ENTRY_ATTR_CHANNEL_TYPE;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_ATTR_HOST_IF:
      return SAI_HOSTIF_TABLE_ENTRY_ATTR_HOST_IF;

    default:
      return SAI_HOSTIF_TABLE_ENTRY_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_hostif_table_entry_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_hostif_table_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_table_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_table_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifTableEntryAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifTableEntryChannelType
convert_sai_hostif_table_entry_channel_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_CB:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_CB;

    case SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_FD:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_FD;

    case SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_PHYSICAL_PORT:
      return lemming::dataplane::sai::
          HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_PHYSICAL_PORT;

    case SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_LOGICAL_PORT:
      return lemming::dataplane::sai::
          HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_LOGICAL_PORT;

    case SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_L3:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_L3;

    case SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_GENETLINK:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_GENETLINK;

    default:
      return lemming::dataplane::sai::
          HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_UNSPECIFIED;
  }
}
sai_hostif_table_entry_channel_type_t
convert_sai_hostif_table_entry_channel_type_t_to_sai(
    lemming::dataplane::sai::HostifTableEntryChannelType val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_CB:
      return SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_CB;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_FD:
      return SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_FD;

    case lemming::dataplane::sai::
        HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_PHYSICAL_PORT:
      return SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_PHYSICAL_PORT;

    case lemming::dataplane::sai::
        HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_LOGICAL_PORT:
      return SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_LOGICAL_PORT;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_L3:
      return SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_NETDEV_L3;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_GENETLINK:
      return SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_GENETLINK;

    default:
      return SAI_HOSTIF_TABLE_ENTRY_CHANNEL_TYPE_CB;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_hostif_table_entry_channel_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_hostif_table_entry_channel_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_table_entry_channel_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_table_entry_channel_type_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifTableEntryChannelType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifTableEntryType
convert_sai_hostif_table_entry_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_TABLE_ENTRY_TYPE_PORT:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_PORT;

    case SAI_HOSTIF_TABLE_ENTRY_TYPE_LAG:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_LAG;

    case SAI_HOSTIF_TABLE_ENTRY_TYPE_VLAN:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_VLAN;

    case SAI_HOSTIF_TABLE_ENTRY_TYPE_TRAP_ID:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_TRAP_ID;

    case SAI_HOSTIF_TABLE_ENTRY_TYPE_WILDCARD:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_WILDCARD;

    default:
      return lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_UNSPECIFIED;
  }
}
sai_hostif_table_entry_type_t convert_sai_hostif_table_entry_type_t_to_sai(
    lemming::dataplane::sai::HostifTableEntryType val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_PORT:
      return SAI_HOSTIF_TABLE_ENTRY_TYPE_PORT;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_LAG:
      return SAI_HOSTIF_TABLE_ENTRY_TYPE_LAG;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_VLAN:
      return SAI_HOSTIF_TABLE_ENTRY_TYPE_VLAN;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_TRAP_ID:
      return SAI_HOSTIF_TABLE_ENTRY_TYPE_TRAP_ID;

    case lemming::dataplane::sai::HOSTIF_TABLE_ENTRY_TYPE_WILDCARD:
      return SAI_HOSTIF_TABLE_ENTRY_TYPE_WILDCARD;

    default:
      return SAI_HOSTIF_TABLE_ENTRY_TYPE_PORT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_hostif_table_entry_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_hostif_table_entry_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_table_entry_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_table_entry_type_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifTableEntryType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifTrapAttr convert_sai_hostif_trap_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_TRAP_ATTR_TRAP_TYPE:
      return lemming::dataplane::sai::HOSTIF_TRAP_ATTR_TRAP_TYPE;

    case SAI_HOSTIF_TRAP_ATTR_PACKET_ACTION:
      return lemming::dataplane::sai::HOSTIF_TRAP_ATTR_PACKET_ACTION;

    case SAI_HOSTIF_TRAP_ATTR_TRAP_PRIORITY:
      return lemming::dataplane::sai::HOSTIF_TRAP_ATTR_TRAP_PRIORITY;

    case SAI_HOSTIF_TRAP_ATTR_EXCLUDE_PORT_LIST:
      return lemming::dataplane::sai::HOSTIF_TRAP_ATTR_EXCLUDE_PORT_LIST;

    case SAI_HOSTIF_TRAP_ATTR_TRAP_GROUP:
      return lemming::dataplane::sai::HOSTIF_TRAP_ATTR_TRAP_GROUP;

    case SAI_HOSTIF_TRAP_ATTR_MIRROR_SESSION:
      return lemming::dataplane::sai::HOSTIF_TRAP_ATTR_MIRROR_SESSION;

    case SAI_HOSTIF_TRAP_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::HOSTIF_TRAP_ATTR_COUNTER_ID;

    default:
      return lemming::dataplane::sai::HOSTIF_TRAP_ATTR_UNSPECIFIED;
  }
}
sai_hostif_trap_attr_t convert_sai_hostif_trap_attr_t_to_sai(
    lemming::dataplane::sai::HostifTrapAttr val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_TRAP_ATTR_TRAP_TYPE:
      return SAI_HOSTIF_TRAP_ATTR_TRAP_TYPE;

    case lemming::dataplane::sai::HOSTIF_TRAP_ATTR_PACKET_ACTION:
      return SAI_HOSTIF_TRAP_ATTR_PACKET_ACTION;

    case lemming::dataplane::sai::HOSTIF_TRAP_ATTR_TRAP_PRIORITY:
      return SAI_HOSTIF_TRAP_ATTR_TRAP_PRIORITY;

    case lemming::dataplane::sai::HOSTIF_TRAP_ATTR_EXCLUDE_PORT_LIST:
      return SAI_HOSTIF_TRAP_ATTR_EXCLUDE_PORT_LIST;

    case lemming::dataplane::sai::HOSTIF_TRAP_ATTR_TRAP_GROUP:
      return SAI_HOSTIF_TRAP_ATTR_TRAP_GROUP;

    case lemming::dataplane::sai::HOSTIF_TRAP_ATTR_MIRROR_SESSION:
      return SAI_HOSTIF_TRAP_ATTR_MIRROR_SESSION;

    case lemming::dataplane::sai::HOSTIF_TRAP_ATTR_COUNTER_ID:
      return SAI_HOSTIF_TRAP_ATTR_COUNTER_ID;

    default:
      return SAI_HOSTIF_TRAP_ATTR_TRAP_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_hostif_trap_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_hostif_trap_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_trap_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_trap_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifTrapAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifTrapGroupAttr
convert_sai_hostif_trap_group_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_TRAP_GROUP_ATTR_ADMIN_STATE:
      return lemming::dataplane::sai::HOSTIF_TRAP_GROUP_ATTR_ADMIN_STATE;

    case SAI_HOSTIF_TRAP_GROUP_ATTR_QUEUE:
      return lemming::dataplane::sai::HOSTIF_TRAP_GROUP_ATTR_QUEUE;

    case SAI_HOSTIF_TRAP_GROUP_ATTR_POLICER:
      return lemming::dataplane::sai::HOSTIF_TRAP_GROUP_ATTR_POLICER;

    case SAI_HOSTIF_TRAP_GROUP_ATTR_OBJECT_STAGE:
      return lemming::dataplane::sai::HOSTIF_TRAP_GROUP_ATTR_OBJECT_STAGE;

    default:
      return lemming::dataplane::sai::HOSTIF_TRAP_GROUP_ATTR_UNSPECIFIED;
  }
}
sai_hostif_trap_group_attr_t convert_sai_hostif_trap_group_attr_t_to_sai(
    lemming::dataplane::sai::HostifTrapGroupAttr val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_TRAP_GROUP_ATTR_ADMIN_STATE:
      return SAI_HOSTIF_TRAP_GROUP_ATTR_ADMIN_STATE;

    case lemming::dataplane::sai::HOSTIF_TRAP_GROUP_ATTR_QUEUE:
      return SAI_HOSTIF_TRAP_GROUP_ATTR_QUEUE;

    case lemming::dataplane::sai::HOSTIF_TRAP_GROUP_ATTR_POLICER:
      return SAI_HOSTIF_TRAP_GROUP_ATTR_POLICER;

    case lemming::dataplane::sai::HOSTIF_TRAP_GROUP_ATTR_OBJECT_STAGE:
      return SAI_HOSTIF_TRAP_GROUP_ATTR_OBJECT_STAGE;

    default:
      return SAI_HOSTIF_TRAP_GROUP_ATTR_ADMIN_STATE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_hostif_trap_group_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_hostif_trap_group_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_trap_group_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_trap_group_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifTrapGroupAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifTrapType convert_sai_hostif_trap_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_TRAP_TYPE_START:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_START;

    case SAI_HOSTIF_TRAP_TYPE_LACP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_LACP;

    case SAI_HOSTIF_TRAP_TYPE_EAPOL:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_EAPOL;

    case SAI_HOSTIF_TRAP_TYPE_LLDP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_LLDP;

    case SAI_HOSTIF_TRAP_TYPE_PVRST:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PVRST;

    case SAI_HOSTIF_TRAP_TYPE_IGMP_TYPE_QUERY:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IGMP_TYPE_QUERY;

    case SAI_HOSTIF_TRAP_TYPE_IGMP_TYPE_LEAVE:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IGMP_TYPE_LEAVE;

    case SAI_HOSTIF_TRAP_TYPE_IGMP_TYPE_V1_REPORT:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IGMP_TYPE_V1_REPORT;

    case SAI_HOSTIF_TRAP_TYPE_IGMP_TYPE_V2_REPORT:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IGMP_TYPE_V2_REPORT;

    case SAI_HOSTIF_TRAP_TYPE_IGMP_TYPE_V3_REPORT:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IGMP_TYPE_V3_REPORT;

    case SAI_HOSTIF_TRAP_TYPE_SAMPLEPACKET:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_SAMPLEPACKET;

    case SAI_HOSTIF_TRAP_TYPE_UDLD:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_UDLD;

    case SAI_HOSTIF_TRAP_TYPE_CDP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_CDP;

    case SAI_HOSTIF_TRAP_TYPE_VTP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_VTP;

    case SAI_HOSTIF_TRAP_TYPE_DTP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DTP;

    case SAI_HOSTIF_TRAP_TYPE_PAGP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PAGP;

    case SAI_HOSTIF_TRAP_TYPE_PTP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PTP;

    case SAI_HOSTIF_TRAP_TYPE_PTP_TX_EVENT:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PTP_TX_EVENT;

    case SAI_HOSTIF_TRAP_TYPE_DHCP_L2:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DHCP_L2;

    case SAI_HOSTIF_TRAP_TYPE_DHCPV6_L2:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DHCPV6_L2;

    case SAI_HOSTIF_TRAP_TYPE_SWITCH_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_SWITCH_CUSTOM_RANGE_BASE;

    case SAI_HOSTIF_TRAP_TYPE_ARP_REQUEST:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_ARP_REQUEST;

    case SAI_HOSTIF_TRAP_TYPE_ARP_RESPONSE:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_ARP_RESPONSE;

    case SAI_HOSTIF_TRAP_TYPE_DHCP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DHCP;

    case SAI_HOSTIF_TRAP_TYPE_OSPF:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_OSPF;

    case SAI_HOSTIF_TRAP_TYPE_PIM:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PIM;

    case SAI_HOSTIF_TRAP_TYPE_VRRP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_VRRP;

    case SAI_HOSTIF_TRAP_TYPE_DHCPV6:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DHCPV6;

    case SAI_HOSTIF_TRAP_TYPE_OSPFV6:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_OSPFV6;

    case SAI_HOSTIF_TRAP_TYPE_VRRPV6:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_VRRPV6;

    case SAI_HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_DISCOVERY:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_DISCOVERY;

    case SAI_HOSTIF_TRAP_TYPE_IPV6_MLD_V1_V2:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IPV6_MLD_V1_V2;

    case SAI_HOSTIF_TRAP_TYPE_IPV6_MLD_V1_REPORT:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IPV6_MLD_V1_REPORT;

    case SAI_HOSTIF_TRAP_TYPE_IPV6_MLD_V1_DONE:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IPV6_MLD_V1_DONE;

    case SAI_HOSTIF_TRAP_TYPE_MLD_V2_REPORT:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_MLD_V2_REPORT;

    case SAI_HOSTIF_TRAP_TYPE_UNKNOWN_L3_MULTICAST:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_UNKNOWN_L3_MULTICAST;

    case SAI_HOSTIF_TRAP_TYPE_SNAT_MISS:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_SNAT_MISS;

    case SAI_HOSTIF_TRAP_TYPE_DNAT_MISS:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DNAT_MISS;

    case SAI_HOSTIF_TRAP_TYPE_NAT_HAIRPIN:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_NAT_HAIRPIN;

    case SAI_HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_SOLICITATION:
      return lemming::dataplane::sai::
          HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_SOLICITATION;

    case SAI_HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_ADVERTISEMENT:
      return lemming::dataplane::sai::
          HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_ADVERTISEMENT;

    case SAI_HOSTIF_TRAP_TYPE_ISIS:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_ISIS;

    case SAI_HOSTIF_TRAP_TYPE_ROUTER_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_ROUTER_CUSTOM_RANGE_BASE;

    case SAI_HOSTIF_TRAP_TYPE_IP2ME:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IP2ME;

    case SAI_HOSTIF_TRAP_TYPE_SSH:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_SSH;

    case SAI_HOSTIF_TRAP_TYPE_SNMP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_SNMP;

    case SAI_HOSTIF_TRAP_TYPE_BGP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BGP;

    case SAI_HOSTIF_TRAP_TYPE_BGPV6:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BGPV6;

    case SAI_HOSTIF_TRAP_TYPE_BFD:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BFD;

    case SAI_HOSTIF_TRAP_TYPE_BFDV6:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BFDV6;

    case SAI_HOSTIF_TRAP_TYPE_BFD_MICRO:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BFD_MICRO;

    case SAI_HOSTIF_TRAP_TYPE_BFDV6_MICRO:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BFDV6_MICRO;

    case SAI_HOSTIF_TRAP_TYPE_LDP:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_LDP;

    case SAI_HOSTIF_TRAP_TYPE_GNMI:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_GNMI;

    case SAI_HOSTIF_TRAP_TYPE_P4RT:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_P4RT;

    case SAI_HOSTIF_TRAP_TYPE_NTPCLIENT:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_NTPCLIENT;

    case SAI_HOSTIF_TRAP_TYPE_NTPSERVER:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_NTPSERVER;

    case SAI_HOSTIF_TRAP_TYPE_LOCAL_IP_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::
          HOSTIF_TRAP_TYPE_LOCAL_IP_CUSTOM_RANGE_BASE;

    case SAI_HOSTIF_TRAP_TYPE_L3_MTU_ERROR:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_L3_MTU_ERROR;

    case SAI_HOSTIF_TRAP_TYPE_TTL_ERROR:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_TTL_ERROR;

    case SAI_HOSTIF_TRAP_TYPE_STATIC_FDB_MOVE:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_STATIC_FDB_MOVE;

    case SAI_HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_EGRESS_BUFFER:
      return lemming::dataplane::sai::
          HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_EGRESS_BUFFER;

    case SAI_HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_WRED:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_WRED;

    case SAI_HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_ROUTER:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_ROUTER;

    case SAI_HOSTIF_TRAP_TYPE_MPLS_TTL_ERROR:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_MPLS_TTL_ERROR;

    case SAI_HOSTIF_TRAP_TYPE_MPLS_ROUTER_ALERT_LABEL:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_MPLS_ROUTER_ALERT_LABEL;

    case SAI_HOSTIF_TRAP_TYPE_MPLS_LABEL_LOOKUP_MISS:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_MPLS_LABEL_LOOKUP_MISS;

    case SAI_HOSTIF_TRAP_TYPE_CUSTOM_EXCEPTION_RANGE_BASE:
      return lemming::dataplane::sai::
          HOSTIF_TRAP_TYPE_CUSTOM_EXCEPTION_RANGE_BASE;

    case SAI_HOSTIF_TRAP_TYPE_END:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_END;

    default:
      return lemming::dataplane::sai::HOSTIF_TRAP_TYPE_UNSPECIFIED;
  }
}
sai_hostif_trap_type_t convert_sai_hostif_trap_type_t_to_sai(
    lemming::dataplane::sai::HostifTrapType val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_START:
      return SAI_HOSTIF_TRAP_TYPE_START;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_LACP:
      return SAI_HOSTIF_TRAP_TYPE_LACP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_EAPOL:
      return SAI_HOSTIF_TRAP_TYPE_EAPOL;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_LLDP:
      return SAI_HOSTIF_TRAP_TYPE_LLDP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PVRST:
      return SAI_HOSTIF_TRAP_TYPE_PVRST;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IGMP_TYPE_QUERY:
      return SAI_HOSTIF_TRAP_TYPE_IGMP_TYPE_QUERY;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IGMP_TYPE_LEAVE:
      return SAI_HOSTIF_TRAP_TYPE_IGMP_TYPE_LEAVE;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IGMP_TYPE_V1_REPORT:
      return SAI_HOSTIF_TRAP_TYPE_IGMP_TYPE_V1_REPORT;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IGMP_TYPE_V2_REPORT:
      return SAI_HOSTIF_TRAP_TYPE_IGMP_TYPE_V2_REPORT;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IGMP_TYPE_V3_REPORT:
      return SAI_HOSTIF_TRAP_TYPE_IGMP_TYPE_V3_REPORT;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_SAMPLEPACKET:
      return SAI_HOSTIF_TRAP_TYPE_SAMPLEPACKET;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_UDLD:
      return SAI_HOSTIF_TRAP_TYPE_UDLD;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_CDP:
      return SAI_HOSTIF_TRAP_TYPE_CDP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_VTP:
      return SAI_HOSTIF_TRAP_TYPE_VTP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DTP:
      return SAI_HOSTIF_TRAP_TYPE_DTP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PAGP:
      return SAI_HOSTIF_TRAP_TYPE_PAGP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PTP:
      return SAI_HOSTIF_TRAP_TYPE_PTP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PTP_TX_EVENT:
      return SAI_HOSTIF_TRAP_TYPE_PTP_TX_EVENT;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DHCP_L2:
      return SAI_HOSTIF_TRAP_TYPE_DHCP_L2;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DHCPV6_L2:
      return SAI_HOSTIF_TRAP_TYPE_DHCPV6_L2;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_SWITCH_CUSTOM_RANGE_BASE:
      return SAI_HOSTIF_TRAP_TYPE_SWITCH_CUSTOM_RANGE_BASE;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_ARP_REQUEST:
      return SAI_HOSTIF_TRAP_TYPE_ARP_REQUEST;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_ARP_RESPONSE:
      return SAI_HOSTIF_TRAP_TYPE_ARP_RESPONSE;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DHCP:
      return SAI_HOSTIF_TRAP_TYPE_DHCP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_OSPF:
      return SAI_HOSTIF_TRAP_TYPE_OSPF;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PIM:
      return SAI_HOSTIF_TRAP_TYPE_PIM;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_VRRP:
      return SAI_HOSTIF_TRAP_TYPE_VRRP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DHCPV6:
      return SAI_HOSTIF_TRAP_TYPE_DHCPV6;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_OSPFV6:
      return SAI_HOSTIF_TRAP_TYPE_OSPFV6;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_VRRPV6:
      return SAI_HOSTIF_TRAP_TYPE_VRRPV6;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_DISCOVERY:
      return SAI_HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_DISCOVERY;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IPV6_MLD_V1_V2:
      return SAI_HOSTIF_TRAP_TYPE_IPV6_MLD_V1_V2;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IPV6_MLD_V1_REPORT:
      return SAI_HOSTIF_TRAP_TYPE_IPV6_MLD_V1_REPORT;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IPV6_MLD_V1_DONE:
      return SAI_HOSTIF_TRAP_TYPE_IPV6_MLD_V1_DONE;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_MLD_V2_REPORT:
      return SAI_HOSTIF_TRAP_TYPE_MLD_V2_REPORT;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_UNKNOWN_L3_MULTICAST:
      return SAI_HOSTIF_TRAP_TYPE_UNKNOWN_L3_MULTICAST;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_SNAT_MISS:
      return SAI_HOSTIF_TRAP_TYPE_SNAT_MISS;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_DNAT_MISS:
      return SAI_HOSTIF_TRAP_TYPE_DNAT_MISS;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_NAT_HAIRPIN:
      return SAI_HOSTIF_TRAP_TYPE_NAT_HAIRPIN;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_SOLICITATION:
      return SAI_HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_SOLICITATION;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_ADVERTISEMENT:
      return SAI_HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_ADVERTISEMENT;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_ISIS:
      return SAI_HOSTIF_TRAP_TYPE_ISIS;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_ROUTER_CUSTOM_RANGE_BASE:
      return SAI_HOSTIF_TRAP_TYPE_ROUTER_CUSTOM_RANGE_BASE;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_IP2ME:
      return SAI_HOSTIF_TRAP_TYPE_IP2ME;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_SSH:
      return SAI_HOSTIF_TRAP_TYPE_SSH;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_SNMP:
      return SAI_HOSTIF_TRAP_TYPE_SNMP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BGP:
      return SAI_HOSTIF_TRAP_TYPE_BGP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BGPV6:
      return SAI_HOSTIF_TRAP_TYPE_BGPV6;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BFD:
      return SAI_HOSTIF_TRAP_TYPE_BFD;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BFDV6:
      return SAI_HOSTIF_TRAP_TYPE_BFDV6;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BFD_MICRO:
      return SAI_HOSTIF_TRAP_TYPE_BFD_MICRO;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_BFDV6_MICRO:
      return SAI_HOSTIF_TRAP_TYPE_BFDV6_MICRO;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_LDP:
      return SAI_HOSTIF_TRAP_TYPE_LDP;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_GNMI:
      return SAI_HOSTIF_TRAP_TYPE_GNMI;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_P4RT:
      return SAI_HOSTIF_TRAP_TYPE_P4RT;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_NTPCLIENT:
      return SAI_HOSTIF_TRAP_TYPE_NTPCLIENT;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_NTPSERVER:
      return SAI_HOSTIF_TRAP_TYPE_NTPSERVER;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_LOCAL_IP_CUSTOM_RANGE_BASE:
      return SAI_HOSTIF_TRAP_TYPE_LOCAL_IP_CUSTOM_RANGE_BASE;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_L3_MTU_ERROR:
      return SAI_HOSTIF_TRAP_TYPE_L3_MTU_ERROR;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_TTL_ERROR:
      return SAI_HOSTIF_TRAP_TYPE_TTL_ERROR;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_STATIC_FDB_MOVE:
      return SAI_HOSTIF_TRAP_TYPE_STATIC_FDB_MOVE;

    case lemming::dataplane::sai::
        HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_EGRESS_BUFFER:
      return SAI_HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_EGRESS_BUFFER;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_WRED:
      return SAI_HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_WRED;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_ROUTER:
      return SAI_HOSTIF_TRAP_TYPE_PIPELINE_DISCARD_ROUTER;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_MPLS_TTL_ERROR:
      return SAI_HOSTIF_TRAP_TYPE_MPLS_TTL_ERROR;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_MPLS_ROUTER_ALERT_LABEL:
      return SAI_HOSTIF_TRAP_TYPE_MPLS_ROUTER_ALERT_LABEL;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_MPLS_LABEL_LOOKUP_MISS:
      return SAI_HOSTIF_TRAP_TYPE_MPLS_LABEL_LOOKUP_MISS;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_CUSTOM_EXCEPTION_RANGE_BASE:
      return SAI_HOSTIF_TRAP_TYPE_CUSTOM_EXCEPTION_RANGE_BASE;

    case lemming::dataplane::sai::HOSTIF_TRAP_TYPE_END:
      return SAI_HOSTIF_TRAP_TYPE_END;

    default:
      return SAI_HOSTIF_TRAP_TYPE_START;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_hostif_trap_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_hostif_trap_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_trap_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_trap_type_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifTrapType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifTxType convert_sai_hostif_tx_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_TX_TYPE_PIPELINE_BYPASS:
      return lemming::dataplane::sai::HOSTIF_TX_TYPE_PIPELINE_BYPASS;

    case SAI_HOSTIF_TX_TYPE_PIPELINE_LOOKUP:
      return lemming::dataplane::sai::HOSTIF_TX_TYPE_PIPELINE_LOOKUP;

    case SAI_HOSTIF_TX_TYPE_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::HOSTIF_TX_TYPE_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::HOSTIF_TX_TYPE_UNSPECIFIED;
  }
}
sai_hostif_tx_type_t convert_sai_hostif_tx_type_t_to_sai(
    lemming::dataplane::sai::HostifTxType val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_TX_TYPE_PIPELINE_BYPASS:
      return SAI_HOSTIF_TX_TYPE_PIPELINE_BYPASS;

    case lemming::dataplane::sai::HOSTIF_TX_TYPE_PIPELINE_LOOKUP:
      return SAI_HOSTIF_TX_TYPE_PIPELINE_LOOKUP;

    case lemming::dataplane::sai::HOSTIF_TX_TYPE_CUSTOM_RANGE_BASE:
      return SAI_HOSTIF_TX_TYPE_CUSTOM_RANGE_BASE;

    default:
      return SAI_HOSTIF_TX_TYPE_PIPELINE_BYPASS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_hostif_tx_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_hostif_tx_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_tx_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_tx_type_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifTxType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifType convert_sai_hostif_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_TYPE_NETDEV:
      return lemming::dataplane::sai::HOSTIF_TYPE_NETDEV;

    case SAI_HOSTIF_TYPE_FD:
      return lemming::dataplane::sai::HOSTIF_TYPE_FD;

    case SAI_HOSTIF_TYPE_GENETLINK:
      return lemming::dataplane::sai::HOSTIF_TYPE_GENETLINK;

    default:
      return lemming::dataplane::sai::HOSTIF_TYPE_UNSPECIFIED;
  }
}
sai_hostif_type_t convert_sai_hostif_type_t_to_sai(
    lemming::dataplane::sai::HostifType val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_TYPE_NETDEV:
      return SAI_HOSTIF_TYPE_NETDEV;

    case lemming::dataplane::sai::HOSTIF_TYPE_FD:
      return SAI_HOSTIF_TYPE_FD;

    case lemming::dataplane::sai::HOSTIF_TYPE_GENETLINK:
      return SAI_HOSTIF_TYPE_GENETLINK;

    default:
      return SAI_HOSTIF_TYPE_NETDEV;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_hostif_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_hostif_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_type_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifUserDefinedTrapAttr
convert_sai_hostif_user_defined_trap_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TYPE:
      return lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_ATTR_TYPE;

    case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_PRIORITY:
      return lemming::dataplane::sai::
          HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_PRIORITY;

    case SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_GROUP:
      return lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_GROUP;

    default:
      return lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_ATTR_UNSPECIFIED;
  }
}
sai_hostif_user_defined_trap_attr_t
convert_sai_hostif_user_defined_trap_attr_t_to_sai(
    lemming::dataplane::sai::HostifUserDefinedTrapAttr val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_ATTR_TYPE:
      return SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TYPE;

    case lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_PRIORITY:
      return SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_PRIORITY;

    case lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_GROUP:
      return SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TRAP_GROUP;

    default:
      return SAI_HOSTIF_USER_DEFINED_TRAP_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_hostif_user_defined_trap_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_hostif_user_defined_trap_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_user_defined_trap_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_user_defined_trap_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifUserDefinedTrapAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifUserDefinedTrapType
convert_sai_hostif_user_defined_trap_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_START:
      return lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_START;

    case SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_NEIGHBOR:
      return lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_NEIGHBOR;

    case SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_ACL:
      return lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_ACL;

    case SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_FDB:
      return lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_FDB;

    case SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_INSEG_ENTRY:
      return lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_INSEG_ENTRY;

    case SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::
          HOSTIF_USER_DEFINED_TRAP_TYPE_CUSTOM_RANGE_BASE;

    case SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_END:
      return lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_END;

    default:
      return lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_UNSPECIFIED;
  }
}
sai_hostif_user_defined_trap_type_t
convert_sai_hostif_user_defined_trap_type_t_to_sai(
    lemming::dataplane::sai::HostifUserDefinedTrapType val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_START:
      return SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_START;

    case lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_NEIGHBOR:
      return SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_NEIGHBOR;

    case lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_ACL:
      return SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_ACL;

    case lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_FDB:
      return SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_FDB;

    case lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_INSEG_ENTRY:
      return SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_INSEG_ENTRY;

    case lemming::dataplane::sai::
        HOSTIF_USER_DEFINED_TRAP_TYPE_CUSTOM_RANGE_BASE:
      return SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_CUSTOM_RANGE_BASE;

    case lemming::dataplane::sai::HOSTIF_USER_DEFINED_TRAP_TYPE_END:
      return SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_END;

    default:
      return SAI_HOSTIF_USER_DEFINED_TRAP_TYPE_START;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_hostif_user_defined_trap_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_hostif_user_defined_trap_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_user_defined_trap_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_user_defined_trap_type_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifUserDefinedTrapType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::HostifVlanTag convert_sai_hostif_vlan_tag_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_HOSTIF_VLAN_TAG_STRIP:
      return lemming::dataplane::sai::HOSTIF_VLAN_TAG_STRIP;

    case SAI_HOSTIF_VLAN_TAG_KEEP:
      return lemming::dataplane::sai::HOSTIF_VLAN_TAG_KEEP;

    case SAI_HOSTIF_VLAN_TAG_ORIGINAL:
      return lemming::dataplane::sai::HOSTIF_VLAN_TAG_ORIGINAL;

    default:
      return lemming::dataplane::sai::HOSTIF_VLAN_TAG_UNSPECIFIED;
  }
}
sai_hostif_vlan_tag_t convert_sai_hostif_vlan_tag_t_to_sai(
    lemming::dataplane::sai::HostifVlanTag val) {
  switch (val) {
    case lemming::dataplane::sai::HOSTIF_VLAN_TAG_STRIP:
      return SAI_HOSTIF_VLAN_TAG_STRIP;

    case lemming::dataplane::sai::HOSTIF_VLAN_TAG_KEEP:
      return SAI_HOSTIF_VLAN_TAG_KEEP;

    case lemming::dataplane::sai::HOSTIF_VLAN_TAG_ORIGINAL:
      return SAI_HOSTIF_VLAN_TAG_ORIGINAL;

    default:
      return SAI_HOSTIF_VLAN_TAG_STRIP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_hostif_vlan_tag_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_hostif_vlan_tag_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_hostif_vlan_tag_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_hostif_vlan_tag_t_to_sai(
        static_cast<lemming::dataplane::sai::HostifVlanTag>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::InDropReason convert_sai_in_drop_reason_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IN_DROP_REASON_START:
      return lemming::dataplane::sai::IN_DROP_REASON_START;

    case SAI_IN_DROP_REASON_SMAC_MULTICAST:
      return lemming::dataplane::sai::IN_DROP_REASON_SMAC_MULTICAST;

    case SAI_IN_DROP_REASON_SMAC_EQUALS_DMAC:
      return lemming::dataplane::sai::IN_DROP_REASON_SMAC_EQUALS_DMAC;

    case SAI_IN_DROP_REASON_DMAC_RESERVED:
      return lemming::dataplane::sai::IN_DROP_REASON_DMAC_RESERVED;

    case SAI_IN_DROP_REASON_VLAN_TAG_NOT_ALLOWED:
      return lemming::dataplane::sai::IN_DROP_REASON_VLAN_TAG_NOT_ALLOWED;

    case SAI_IN_DROP_REASON_INGRESS_VLAN_FILTER:
      return lemming::dataplane::sai::IN_DROP_REASON_INGRESS_VLAN_FILTER;

    case SAI_IN_DROP_REASON_INGRESS_STP_FILTER:
      return lemming::dataplane::sai::IN_DROP_REASON_INGRESS_STP_FILTER;

    case SAI_IN_DROP_REASON_FDB_UC_DISCARD:
      return lemming::dataplane::sai::IN_DROP_REASON_FDB_UC_DISCARD;

    case SAI_IN_DROP_REASON_FDB_MC_DISCARD:
      return lemming::dataplane::sai::IN_DROP_REASON_FDB_MC_DISCARD;

    case SAI_IN_DROP_REASON_L2_LOOPBACK_FILTER:
      return lemming::dataplane::sai::IN_DROP_REASON_L2_LOOPBACK_FILTER;

    case SAI_IN_DROP_REASON_EXCEEDS_L2_MTU:
      return lemming::dataplane::sai::IN_DROP_REASON_EXCEEDS_L2_MTU;

    case SAI_IN_DROP_REASON_L3_ANY:
      return lemming::dataplane::sai::IN_DROP_REASON_L3_ANY;

    case SAI_IN_DROP_REASON_EXCEEDS_L3_MTU:
      return lemming::dataplane::sai::IN_DROP_REASON_EXCEEDS_L3_MTU;

    case SAI_IN_DROP_REASON_TTL:
      return lemming::dataplane::sai::IN_DROP_REASON_TTL;

    case SAI_IN_DROP_REASON_L3_LOOPBACK_FILTER:
      return lemming::dataplane::sai::IN_DROP_REASON_L3_LOOPBACK_FILTER;

    case SAI_IN_DROP_REASON_NON_ROUTABLE:
      return lemming::dataplane::sai::IN_DROP_REASON_NON_ROUTABLE;

    case SAI_IN_DROP_REASON_NO_L3_HEADER:
      return lemming::dataplane::sai::IN_DROP_REASON_NO_L3_HEADER;

    case SAI_IN_DROP_REASON_IP_HEADER_ERROR:
      return lemming::dataplane::sai::IN_DROP_REASON_IP_HEADER_ERROR;

    case SAI_IN_DROP_REASON_UC_DIP_MC_DMAC:
      return lemming::dataplane::sai::IN_DROP_REASON_UC_DIP_MC_DMAC;

    case SAI_IN_DROP_REASON_DIP_LOOPBACK:
      return lemming::dataplane::sai::IN_DROP_REASON_DIP_LOOPBACK;

    case SAI_IN_DROP_REASON_SIP_LOOPBACK:
      return lemming::dataplane::sai::IN_DROP_REASON_SIP_LOOPBACK;

    case SAI_IN_DROP_REASON_SIP_MC:
      return lemming::dataplane::sai::IN_DROP_REASON_SIP_MC;

    case SAI_IN_DROP_REASON_SIP_CLASS_E:
      return lemming::dataplane::sai::IN_DROP_REASON_SIP_CLASS_E;

    case SAI_IN_DROP_REASON_SIP_UNSPECIFIED:
      return lemming::dataplane::sai::IN_DROP_REASON_SIP_UNSPECIFIED;

    case SAI_IN_DROP_REASON_MC_DMAC_MISMATCH:
      return lemming::dataplane::sai::IN_DROP_REASON_MC_DMAC_MISMATCH;

    case SAI_IN_DROP_REASON_SIP_EQUALS_DIP:
      return lemming::dataplane::sai::IN_DROP_REASON_SIP_EQUALS_DIP;

    case SAI_IN_DROP_REASON_SIP_BC:
      return lemming::dataplane::sai::IN_DROP_REASON_SIP_BC;

    case SAI_IN_DROP_REASON_DIP_LOCAL:
      return lemming::dataplane::sai::IN_DROP_REASON_DIP_LOCAL;

    case SAI_IN_DROP_REASON_DIP_LINK_LOCAL:
      return lemming::dataplane::sai::IN_DROP_REASON_DIP_LINK_LOCAL;

    case SAI_IN_DROP_REASON_SIP_LINK_LOCAL:
      return lemming::dataplane::sai::IN_DROP_REASON_SIP_LINK_LOCAL;

    case SAI_IN_DROP_REASON_IPV6_MC_SCOPE0:
      return lemming::dataplane::sai::IN_DROP_REASON_IPV6_MC_SCOPE0;

    case SAI_IN_DROP_REASON_IPV6_MC_SCOPE1:
      return lemming::dataplane::sai::IN_DROP_REASON_IPV6_MC_SCOPE1;

    case SAI_IN_DROP_REASON_IRIF_DISABLED:
      return lemming::dataplane::sai::IN_DROP_REASON_IRIF_DISABLED;

    case SAI_IN_DROP_REASON_ERIF_DISABLED:
      return lemming::dataplane::sai::IN_DROP_REASON_ERIF_DISABLED;

    case SAI_IN_DROP_REASON_LPM4_MISS:
      return lemming::dataplane::sai::IN_DROP_REASON_LPM4_MISS;

    case SAI_IN_DROP_REASON_LPM6_MISS:
      return lemming::dataplane::sai::IN_DROP_REASON_LPM6_MISS;

    case SAI_IN_DROP_REASON_BLACKHOLE_ROUTE:
      return lemming::dataplane::sai::IN_DROP_REASON_BLACKHOLE_ROUTE;

    case SAI_IN_DROP_REASON_BLACKHOLE_ARP:
      return lemming::dataplane::sai::IN_DROP_REASON_BLACKHOLE_ARP;

    case SAI_IN_DROP_REASON_UNRESOLVED_NEXT_HOP:
      return lemming::dataplane::sai::IN_DROP_REASON_UNRESOLVED_NEXT_HOP;

    case SAI_IN_DROP_REASON_L3_EGRESS_LINK_DOWN:
      return lemming::dataplane::sai::IN_DROP_REASON_L3_EGRESS_LINK_DOWN;

    case SAI_IN_DROP_REASON_DECAP_ERROR:
      return lemming::dataplane::sai::IN_DROP_REASON_DECAP_ERROR;

    case SAI_IN_DROP_REASON_ACL_ANY:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_ANY;

    case SAI_IN_DROP_REASON_ACL_INGRESS_PORT:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_INGRESS_PORT;

    case SAI_IN_DROP_REASON_ACL_INGRESS_LAG:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_INGRESS_LAG;

    case SAI_IN_DROP_REASON_ACL_INGRESS_VLAN:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_INGRESS_VLAN;

    case SAI_IN_DROP_REASON_ACL_INGRESS_RIF:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_INGRESS_RIF;

    case SAI_IN_DROP_REASON_ACL_INGRESS_SWITCH:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_INGRESS_SWITCH;

    case SAI_IN_DROP_REASON_ACL_EGRESS_PORT:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_EGRESS_PORT;

    case SAI_IN_DROP_REASON_ACL_EGRESS_LAG:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_EGRESS_LAG;

    case SAI_IN_DROP_REASON_ACL_EGRESS_VLAN:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_EGRESS_VLAN;

    case SAI_IN_DROP_REASON_ACL_EGRESS_RIF:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_EGRESS_RIF;

    case SAI_IN_DROP_REASON_ACL_EGRESS_SWITCH:
      return lemming::dataplane::sai::IN_DROP_REASON_ACL_EGRESS_SWITCH;

    case SAI_IN_DROP_REASON_FDB_AND_BLACKHOLE_DISCARDS:
      return lemming::dataplane::sai::IN_DROP_REASON_FDB_AND_BLACKHOLE_DISCARDS;

    case SAI_IN_DROP_REASON_MPLS_MISS:
      return lemming::dataplane::sai::IN_DROP_REASON_MPLS_MISS;

    case SAI_IN_DROP_REASON_SRV6_LOCAL_SID_DROP:
      return lemming::dataplane::sai::IN_DROP_REASON_SRV6_LOCAL_SID_DROP;

    case SAI_IN_DROP_REASON_END:
      return lemming::dataplane::sai::IN_DROP_REASON_END;

    case SAI_IN_DROP_REASON_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::IN_DROP_REASON_CUSTOM_RANGE_BASE;

    case SAI_IN_DROP_REASON_CUSTOM_RANGE_END:
      return lemming::dataplane::sai::IN_DROP_REASON_CUSTOM_RANGE_END;

    default:
      return lemming::dataplane::sai::IN_DROP_REASON_UNSPECIFIED;
  }
}
sai_in_drop_reason_t convert_sai_in_drop_reason_t_to_sai(
    lemming::dataplane::sai::InDropReason val) {
  switch (val) {
    case lemming::dataplane::sai::IN_DROP_REASON_START:
      return SAI_IN_DROP_REASON_START;

    case lemming::dataplane::sai::IN_DROP_REASON_SMAC_MULTICAST:
      return SAI_IN_DROP_REASON_SMAC_MULTICAST;

    case lemming::dataplane::sai::IN_DROP_REASON_SMAC_EQUALS_DMAC:
      return SAI_IN_DROP_REASON_SMAC_EQUALS_DMAC;

    case lemming::dataplane::sai::IN_DROP_REASON_DMAC_RESERVED:
      return SAI_IN_DROP_REASON_DMAC_RESERVED;

    case lemming::dataplane::sai::IN_DROP_REASON_VLAN_TAG_NOT_ALLOWED:
      return SAI_IN_DROP_REASON_VLAN_TAG_NOT_ALLOWED;

    case lemming::dataplane::sai::IN_DROP_REASON_INGRESS_VLAN_FILTER:
      return SAI_IN_DROP_REASON_INGRESS_VLAN_FILTER;

    case lemming::dataplane::sai::IN_DROP_REASON_INGRESS_STP_FILTER:
      return SAI_IN_DROP_REASON_INGRESS_STP_FILTER;

    case lemming::dataplane::sai::IN_DROP_REASON_FDB_UC_DISCARD:
      return SAI_IN_DROP_REASON_FDB_UC_DISCARD;

    case lemming::dataplane::sai::IN_DROP_REASON_FDB_MC_DISCARD:
      return SAI_IN_DROP_REASON_FDB_MC_DISCARD;

    case lemming::dataplane::sai::IN_DROP_REASON_L2_LOOPBACK_FILTER:
      return SAI_IN_DROP_REASON_L2_LOOPBACK_FILTER;

    case lemming::dataplane::sai::IN_DROP_REASON_EXCEEDS_L2_MTU:
      return SAI_IN_DROP_REASON_EXCEEDS_L2_MTU;

    case lemming::dataplane::sai::IN_DROP_REASON_L3_ANY:
      return SAI_IN_DROP_REASON_L3_ANY;

    case lemming::dataplane::sai::IN_DROP_REASON_EXCEEDS_L3_MTU:
      return SAI_IN_DROP_REASON_EXCEEDS_L3_MTU;

    case lemming::dataplane::sai::IN_DROP_REASON_TTL:
      return SAI_IN_DROP_REASON_TTL;

    case lemming::dataplane::sai::IN_DROP_REASON_L3_LOOPBACK_FILTER:
      return SAI_IN_DROP_REASON_L3_LOOPBACK_FILTER;

    case lemming::dataplane::sai::IN_DROP_REASON_NON_ROUTABLE:
      return SAI_IN_DROP_REASON_NON_ROUTABLE;

    case lemming::dataplane::sai::IN_DROP_REASON_NO_L3_HEADER:
      return SAI_IN_DROP_REASON_NO_L3_HEADER;

    case lemming::dataplane::sai::IN_DROP_REASON_IP_HEADER_ERROR:
      return SAI_IN_DROP_REASON_IP_HEADER_ERROR;

    case lemming::dataplane::sai::IN_DROP_REASON_UC_DIP_MC_DMAC:
      return SAI_IN_DROP_REASON_UC_DIP_MC_DMAC;

    case lemming::dataplane::sai::IN_DROP_REASON_DIP_LOOPBACK:
      return SAI_IN_DROP_REASON_DIP_LOOPBACK;

    case lemming::dataplane::sai::IN_DROP_REASON_SIP_LOOPBACK:
      return SAI_IN_DROP_REASON_SIP_LOOPBACK;

    case lemming::dataplane::sai::IN_DROP_REASON_SIP_MC:
      return SAI_IN_DROP_REASON_SIP_MC;

    case lemming::dataplane::sai::IN_DROP_REASON_SIP_CLASS_E:
      return SAI_IN_DROP_REASON_SIP_CLASS_E;

    case lemming::dataplane::sai::IN_DROP_REASON_SIP_UNSPECIFIED:
      return SAI_IN_DROP_REASON_SIP_UNSPECIFIED;

    case lemming::dataplane::sai::IN_DROP_REASON_MC_DMAC_MISMATCH:
      return SAI_IN_DROP_REASON_MC_DMAC_MISMATCH;

    case lemming::dataplane::sai::IN_DROP_REASON_SIP_EQUALS_DIP:
      return SAI_IN_DROP_REASON_SIP_EQUALS_DIP;

    case lemming::dataplane::sai::IN_DROP_REASON_SIP_BC:
      return SAI_IN_DROP_REASON_SIP_BC;

    case lemming::dataplane::sai::IN_DROP_REASON_DIP_LOCAL:
      return SAI_IN_DROP_REASON_DIP_LOCAL;

    case lemming::dataplane::sai::IN_DROP_REASON_DIP_LINK_LOCAL:
      return SAI_IN_DROP_REASON_DIP_LINK_LOCAL;

    case lemming::dataplane::sai::IN_DROP_REASON_SIP_LINK_LOCAL:
      return SAI_IN_DROP_REASON_SIP_LINK_LOCAL;

    case lemming::dataplane::sai::IN_DROP_REASON_IPV6_MC_SCOPE0:
      return SAI_IN_DROP_REASON_IPV6_MC_SCOPE0;

    case lemming::dataplane::sai::IN_DROP_REASON_IPV6_MC_SCOPE1:
      return SAI_IN_DROP_REASON_IPV6_MC_SCOPE1;

    case lemming::dataplane::sai::IN_DROP_REASON_IRIF_DISABLED:
      return SAI_IN_DROP_REASON_IRIF_DISABLED;

    case lemming::dataplane::sai::IN_DROP_REASON_ERIF_DISABLED:
      return SAI_IN_DROP_REASON_ERIF_DISABLED;

    case lemming::dataplane::sai::IN_DROP_REASON_LPM4_MISS:
      return SAI_IN_DROP_REASON_LPM4_MISS;

    case lemming::dataplane::sai::IN_DROP_REASON_LPM6_MISS:
      return SAI_IN_DROP_REASON_LPM6_MISS;

    case lemming::dataplane::sai::IN_DROP_REASON_BLACKHOLE_ROUTE:
      return SAI_IN_DROP_REASON_BLACKHOLE_ROUTE;

    case lemming::dataplane::sai::IN_DROP_REASON_BLACKHOLE_ARP:
      return SAI_IN_DROP_REASON_BLACKHOLE_ARP;

    case lemming::dataplane::sai::IN_DROP_REASON_UNRESOLVED_NEXT_HOP:
      return SAI_IN_DROP_REASON_UNRESOLVED_NEXT_HOP;

    case lemming::dataplane::sai::IN_DROP_REASON_L3_EGRESS_LINK_DOWN:
      return SAI_IN_DROP_REASON_L3_EGRESS_LINK_DOWN;

    case lemming::dataplane::sai::IN_DROP_REASON_DECAP_ERROR:
      return SAI_IN_DROP_REASON_DECAP_ERROR;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_ANY:
      return SAI_IN_DROP_REASON_ACL_ANY;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_INGRESS_PORT:
      return SAI_IN_DROP_REASON_ACL_INGRESS_PORT;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_INGRESS_LAG:
      return SAI_IN_DROP_REASON_ACL_INGRESS_LAG;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_INGRESS_VLAN:
      return SAI_IN_DROP_REASON_ACL_INGRESS_VLAN;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_INGRESS_RIF:
      return SAI_IN_DROP_REASON_ACL_INGRESS_RIF;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_INGRESS_SWITCH:
      return SAI_IN_DROP_REASON_ACL_INGRESS_SWITCH;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_EGRESS_PORT:
      return SAI_IN_DROP_REASON_ACL_EGRESS_PORT;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_EGRESS_LAG:
      return SAI_IN_DROP_REASON_ACL_EGRESS_LAG;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_EGRESS_VLAN:
      return SAI_IN_DROP_REASON_ACL_EGRESS_VLAN;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_EGRESS_RIF:
      return SAI_IN_DROP_REASON_ACL_EGRESS_RIF;

    case lemming::dataplane::sai::IN_DROP_REASON_ACL_EGRESS_SWITCH:
      return SAI_IN_DROP_REASON_ACL_EGRESS_SWITCH;

    case lemming::dataplane::sai::IN_DROP_REASON_FDB_AND_BLACKHOLE_DISCARDS:
      return SAI_IN_DROP_REASON_FDB_AND_BLACKHOLE_DISCARDS;

    case lemming::dataplane::sai::IN_DROP_REASON_MPLS_MISS:
      return SAI_IN_DROP_REASON_MPLS_MISS;

    case lemming::dataplane::sai::IN_DROP_REASON_SRV6_LOCAL_SID_DROP:
      return SAI_IN_DROP_REASON_SRV6_LOCAL_SID_DROP;

    case lemming::dataplane::sai::IN_DROP_REASON_END:
      return SAI_IN_DROP_REASON_END;

    case lemming::dataplane::sai::IN_DROP_REASON_CUSTOM_RANGE_BASE:
      return SAI_IN_DROP_REASON_CUSTOM_RANGE_BASE;

    case lemming::dataplane::sai::IN_DROP_REASON_CUSTOM_RANGE_END:
      return SAI_IN_DROP_REASON_CUSTOM_RANGE_END;

    default:
      return SAI_IN_DROP_REASON_START;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_in_drop_reason_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_in_drop_reason_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_in_drop_reason_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_in_drop_reason_t_to_sai(
        static_cast<lemming::dataplane::sai::InDropReason>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IngressPriorityGroupAttr
convert_sai_ingress_priority_group_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_INGRESS_PRIORITY_GROUP_ATTR_BUFFER_PROFILE:
      return lemming::dataplane::sai::
          INGRESS_PRIORITY_GROUP_ATTR_BUFFER_PROFILE;

    case SAI_INGRESS_PRIORITY_GROUP_ATTR_PORT:
      return lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_ATTR_PORT;

    case SAI_INGRESS_PRIORITY_GROUP_ATTR_TAM:
      return lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_ATTR_TAM;

    case SAI_INGRESS_PRIORITY_GROUP_ATTR_INDEX:
      return lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_ATTR_INDEX;

    default:
      return lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_ATTR_UNSPECIFIED;
  }
}
sai_ingress_priority_group_attr_t
convert_sai_ingress_priority_group_attr_t_to_sai(
    lemming::dataplane::sai::IngressPriorityGroupAttr val) {
  switch (val) {
    case lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_ATTR_BUFFER_PROFILE:
      return SAI_INGRESS_PRIORITY_GROUP_ATTR_BUFFER_PROFILE;

    case lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_ATTR_PORT:
      return SAI_INGRESS_PRIORITY_GROUP_ATTR_PORT;

    case lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_ATTR_TAM:
      return SAI_INGRESS_PRIORITY_GROUP_ATTR_TAM;

    case lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_ATTR_INDEX:
      return SAI_INGRESS_PRIORITY_GROUP_ATTR_INDEX;

    default:
      return SAI_INGRESS_PRIORITY_GROUP_ATTR_BUFFER_PROFILE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_ingress_priority_group_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_ingress_priority_group_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ingress_priority_group_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ingress_priority_group_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::IngressPriorityGroupAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IngressPriorityGroupStat
convert_sai_ingress_priority_group_stat_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_INGRESS_PRIORITY_GROUP_STAT_PACKETS:
      return lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_STAT_PACKETS;

    case SAI_INGRESS_PRIORITY_GROUP_STAT_BYTES:
      return lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_STAT_BYTES;

    case SAI_INGRESS_PRIORITY_GROUP_STAT_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::
          INGRESS_PRIORITY_GROUP_STAT_CURR_OCCUPANCY_BYTES;

    case SAI_INGRESS_PRIORITY_GROUP_STAT_WATERMARK_BYTES:
      return lemming::dataplane::sai::
          INGRESS_PRIORITY_GROUP_STAT_WATERMARK_BYTES;

    case SAI_INGRESS_PRIORITY_GROUP_STAT_SHARED_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::
          INGRESS_PRIORITY_GROUP_STAT_SHARED_CURR_OCCUPANCY_BYTES;

    case SAI_INGRESS_PRIORITY_GROUP_STAT_SHARED_WATERMARK_BYTES:
      return lemming::dataplane::sai::
          INGRESS_PRIORITY_GROUP_STAT_SHARED_WATERMARK_BYTES;

    case SAI_INGRESS_PRIORITY_GROUP_STAT_XOFF_ROOM_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::
          INGRESS_PRIORITY_GROUP_STAT_XOFF_ROOM_CURR_OCCUPANCY_BYTES;

    case SAI_INGRESS_PRIORITY_GROUP_STAT_XOFF_ROOM_WATERMARK_BYTES:
      return lemming::dataplane::sai::
          INGRESS_PRIORITY_GROUP_STAT_XOFF_ROOM_WATERMARK_BYTES;

    case SAI_INGRESS_PRIORITY_GROUP_STAT_DROPPED_PACKETS:
      return lemming::dataplane::sai::
          INGRESS_PRIORITY_GROUP_STAT_DROPPED_PACKETS;

    case SAI_INGRESS_PRIORITY_GROUP_STAT_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::
          INGRESS_PRIORITY_GROUP_STAT_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_STAT_UNSPECIFIED;
  }
}
sai_ingress_priority_group_stat_t
convert_sai_ingress_priority_group_stat_t_to_sai(
    lemming::dataplane::sai::IngressPriorityGroupStat val) {
  switch (val) {
    case lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_STAT_PACKETS:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_PACKETS;

    case lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_STAT_BYTES:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_BYTES;

    case lemming::dataplane::sai::
        INGRESS_PRIORITY_GROUP_STAT_CURR_OCCUPANCY_BYTES:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_STAT_WATERMARK_BYTES:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_WATERMARK_BYTES;

    case lemming::dataplane::sai::
        INGRESS_PRIORITY_GROUP_STAT_SHARED_CURR_OCCUPANCY_BYTES:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_SHARED_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::
        INGRESS_PRIORITY_GROUP_STAT_SHARED_WATERMARK_BYTES:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_SHARED_WATERMARK_BYTES;

    case lemming::dataplane::sai::
        INGRESS_PRIORITY_GROUP_STAT_XOFF_ROOM_CURR_OCCUPANCY_BYTES:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_XOFF_ROOM_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::
        INGRESS_PRIORITY_GROUP_STAT_XOFF_ROOM_WATERMARK_BYTES:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_XOFF_ROOM_WATERMARK_BYTES;

    case lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_STAT_DROPPED_PACKETS:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_DROPPED_PACKETS;

    case lemming::dataplane::sai::INGRESS_PRIORITY_GROUP_STAT_CUSTOM_RANGE_BASE:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_CUSTOM_RANGE_BASE;

    default:
      return SAI_INGRESS_PRIORITY_GROUP_STAT_PACKETS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_ingress_priority_group_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_ingress_priority_group_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ingress_priority_group_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ingress_priority_group_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::IngressPriorityGroupStat>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::InsegEntryAttr convert_sai_inseg_entry_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_INSEG_ENTRY_ATTR_NUM_OF_POP:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_NUM_OF_POP;

    case SAI_INSEG_ENTRY_ATTR_PACKET_ACTION:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_PACKET_ACTION;

    case SAI_INSEG_ENTRY_ATTR_TRAP_PRIORITY:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_TRAP_PRIORITY;

    case SAI_INSEG_ENTRY_ATTR_NEXT_HOP_ID:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_NEXT_HOP_ID;

    case SAI_INSEG_ENTRY_ATTR_PSC_TYPE:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_PSC_TYPE;

    case SAI_INSEG_ENTRY_ATTR_QOS_TC:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_QOS_TC;

    case SAI_INSEG_ENTRY_ATTR_MPLS_EXP_TO_TC_MAP:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_MPLS_EXP_TO_TC_MAP;

    case SAI_INSEG_ENTRY_ATTR_MPLS_EXP_TO_COLOR_MAP:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_MPLS_EXP_TO_COLOR_MAP;

    case SAI_INSEG_ENTRY_ATTR_POP_TTL_MODE:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_POP_TTL_MODE;

    case SAI_INSEG_ENTRY_ATTR_POP_QOS_MODE:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_POP_QOS_MODE;

    case SAI_INSEG_ENTRY_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_COUNTER_ID;

    default:
      return lemming::dataplane::sai::INSEG_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_inseg_entry_attr_t convert_sai_inseg_entry_attr_t_to_sai(
    lemming::dataplane::sai::InsegEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_NUM_OF_POP:
      return SAI_INSEG_ENTRY_ATTR_NUM_OF_POP;

    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_PACKET_ACTION:
      return SAI_INSEG_ENTRY_ATTR_PACKET_ACTION;

    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_TRAP_PRIORITY:
      return SAI_INSEG_ENTRY_ATTR_TRAP_PRIORITY;

    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_NEXT_HOP_ID:
      return SAI_INSEG_ENTRY_ATTR_NEXT_HOP_ID;

    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_PSC_TYPE:
      return SAI_INSEG_ENTRY_ATTR_PSC_TYPE;

    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_QOS_TC:
      return SAI_INSEG_ENTRY_ATTR_QOS_TC;

    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_MPLS_EXP_TO_TC_MAP:
      return SAI_INSEG_ENTRY_ATTR_MPLS_EXP_TO_TC_MAP;

    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_MPLS_EXP_TO_COLOR_MAP:
      return SAI_INSEG_ENTRY_ATTR_MPLS_EXP_TO_COLOR_MAP;

    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_POP_TTL_MODE:
      return SAI_INSEG_ENTRY_ATTR_POP_TTL_MODE;

    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_POP_QOS_MODE:
      return SAI_INSEG_ENTRY_ATTR_POP_QOS_MODE;

    case lemming::dataplane::sai::INSEG_ENTRY_ATTR_COUNTER_ID:
      return SAI_INSEG_ENTRY_ATTR_COUNTER_ID;

    default:
      return SAI_INSEG_ENTRY_ATTR_NUM_OF_POP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_inseg_entry_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_inseg_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_inseg_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_inseg_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::InsegEntryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::InsegEntryPopQosMode
convert_sai_inseg_entry_pop_qos_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_INSEG_ENTRY_POP_QOS_MODE_UNIFORM:
      return lemming::dataplane::sai::INSEG_ENTRY_POP_QOS_MODE_UNIFORM;

    case SAI_INSEG_ENTRY_POP_QOS_MODE_PIPE:
      return lemming::dataplane::sai::INSEG_ENTRY_POP_QOS_MODE_PIPE;

    default:
      return lemming::dataplane::sai::INSEG_ENTRY_POP_QOS_MODE_UNSPECIFIED;
  }
}
sai_inseg_entry_pop_qos_mode_t convert_sai_inseg_entry_pop_qos_mode_t_to_sai(
    lemming::dataplane::sai::InsegEntryPopQosMode val) {
  switch (val) {
    case lemming::dataplane::sai::INSEG_ENTRY_POP_QOS_MODE_UNIFORM:
      return SAI_INSEG_ENTRY_POP_QOS_MODE_UNIFORM;

    case lemming::dataplane::sai::INSEG_ENTRY_POP_QOS_MODE_PIPE:
      return SAI_INSEG_ENTRY_POP_QOS_MODE_PIPE;

    default:
      return SAI_INSEG_ENTRY_POP_QOS_MODE_UNIFORM;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_inseg_entry_pop_qos_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_inseg_entry_pop_qos_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_inseg_entry_pop_qos_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_inseg_entry_pop_qos_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::InsegEntryPopQosMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::InsegEntryPopTtlMode
convert_sai_inseg_entry_pop_ttl_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_INSEG_ENTRY_POP_TTL_MODE_UNIFORM:
      return lemming::dataplane::sai::INSEG_ENTRY_POP_TTL_MODE_UNIFORM;

    case SAI_INSEG_ENTRY_POP_TTL_MODE_PIPE:
      return lemming::dataplane::sai::INSEG_ENTRY_POP_TTL_MODE_PIPE;

    default:
      return lemming::dataplane::sai::INSEG_ENTRY_POP_TTL_MODE_UNSPECIFIED;
  }
}
sai_inseg_entry_pop_ttl_mode_t convert_sai_inseg_entry_pop_ttl_mode_t_to_sai(
    lemming::dataplane::sai::InsegEntryPopTtlMode val) {
  switch (val) {
    case lemming::dataplane::sai::INSEG_ENTRY_POP_TTL_MODE_UNIFORM:
      return SAI_INSEG_ENTRY_POP_TTL_MODE_UNIFORM;

    case lemming::dataplane::sai::INSEG_ENTRY_POP_TTL_MODE_PIPE:
      return SAI_INSEG_ENTRY_POP_TTL_MODE_PIPE;

    default:
      return SAI_INSEG_ENTRY_POP_TTL_MODE_UNIFORM;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_inseg_entry_pop_ttl_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_inseg_entry_pop_ttl_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_inseg_entry_pop_ttl_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_inseg_entry_pop_ttl_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::InsegEntryPopTtlMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::InsegEntryPscType
convert_sai_inseg_entry_psc_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_INSEG_ENTRY_PSC_TYPE_ELSP:
      return lemming::dataplane::sai::INSEG_ENTRY_PSC_TYPE_ELSP;

    case SAI_INSEG_ENTRY_PSC_TYPE_LLSP:
      return lemming::dataplane::sai::INSEG_ENTRY_PSC_TYPE_LLSP;

    default:
      return lemming::dataplane::sai::INSEG_ENTRY_PSC_TYPE_UNSPECIFIED;
  }
}
sai_inseg_entry_psc_type_t convert_sai_inseg_entry_psc_type_t_to_sai(
    lemming::dataplane::sai::InsegEntryPscType val) {
  switch (val) {
    case lemming::dataplane::sai::INSEG_ENTRY_PSC_TYPE_ELSP:
      return SAI_INSEG_ENTRY_PSC_TYPE_ELSP;

    case lemming::dataplane::sai::INSEG_ENTRY_PSC_TYPE_LLSP:
      return SAI_INSEG_ENTRY_PSC_TYPE_LLSP;

    default:
      return SAI_INSEG_ENTRY_PSC_TYPE_ELSP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_inseg_entry_psc_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_inseg_entry_psc_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_inseg_entry_psc_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_inseg_entry_psc_type_t_to_sai(
        static_cast<lemming::dataplane::sai::InsegEntryPscType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpAddrFamily convert_sai_ip_addr_family_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IP_ADDR_FAMILY_IPV4:
      return lemming::dataplane::sai::IP_ADDR_FAMILY_IPV4;

    case SAI_IP_ADDR_FAMILY_IPV6:
      return lemming::dataplane::sai::IP_ADDR_FAMILY_IPV6;

    default:
      return lemming::dataplane::sai::IP_ADDR_FAMILY_UNSPECIFIED;
  }
}
sai_ip_addr_family_t convert_sai_ip_addr_family_t_to_sai(
    lemming::dataplane::sai::IpAddrFamily val) {
  switch (val) {
    case lemming::dataplane::sai::IP_ADDR_FAMILY_IPV4:
      return SAI_IP_ADDR_FAMILY_IPV4;

    case lemming::dataplane::sai::IP_ADDR_FAMILY_IPV6:
      return SAI_IP_ADDR_FAMILY_IPV6;

    default:
      return SAI_IP_ADDR_FAMILY_IPV4;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_ip_addr_family_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ip_addr_family_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ip_addr_family_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ip_addr_family_t_to_sai(
        static_cast<lemming::dataplane::sai::IpAddrFamily>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpmcEntryAttr convert_sai_ipmc_entry_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IPMC_ENTRY_ATTR_PACKET_ACTION:
      return lemming::dataplane::sai::IPMC_ENTRY_ATTR_PACKET_ACTION;

    case SAI_IPMC_ENTRY_ATTR_OUTPUT_GROUP_ID:
      return lemming::dataplane::sai::IPMC_ENTRY_ATTR_OUTPUT_GROUP_ID;

    case SAI_IPMC_ENTRY_ATTR_RPF_GROUP_ID:
      return lemming::dataplane::sai::IPMC_ENTRY_ATTR_RPF_GROUP_ID;

    case SAI_IPMC_ENTRY_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::IPMC_ENTRY_ATTR_COUNTER_ID;

    default:
      return lemming::dataplane::sai::IPMC_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_ipmc_entry_attr_t convert_sai_ipmc_entry_attr_t_to_sai(
    lemming::dataplane::sai::IpmcEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::IPMC_ENTRY_ATTR_PACKET_ACTION:
      return SAI_IPMC_ENTRY_ATTR_PACKET_ACTION;

    case lemming::dataplane::sai::IPMC_ENTRY_ATTR_OUTPUT_GROUP_ID:
      return SAI_IPMC_ENTRY_ATTR_OUTPUT_GROUP_ID;

    case lemming::dataplane::sai::IPMC_ENTRY_ATTR_RPF_GROUP_ID:
      return SAI_IPMC_ENTRY_ATTR_RPF_GROUP_ID;

    case lemming::dataplane::sai::IPMC_ENTRY_ATTR_COUNTER_ID:
      return SAI_IPMC_ENTRY_ATTR_COUNTER_ID;

    default:
      return SAI_IPMC_ENTRY_ATTR_PACKET_ACTION;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_ipmc_entry_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipmc_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipmc_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipmc_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::IpmcEntryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpmcEntryType convert_sai_ipmc_entry_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IPMC_ENTRY_TYPE_SG:
      return lemming::dataplane::sai::IPMC_ENTRY_TYPE_SG;

    case SAI_IPMC_ENTRY_TYPE_XG:
      return lemming::dataplane::sai::IPMC_ENTRY_TYPE_XG;

    default:
      return lemming::dataplane::sai::IPMC_ENTRY_TYPE_UNSPECIFIED;
  }
}
sai_ipmc_entry_type_t convert_sai_ipmc_entry_type_t_to_sai(
    lemming::dataplane::sai::IpmcEntryType val) {
  switch (val) {
    case lemming::dataplane::sai::IPMC_ENTRY_TYPE_SG:
      return SAI_IPMC_ENTRY_TYPE_SG;

    case lemming::dataplane::sai::IPMC_ENTRY_TYPE_XG:
      return SAI_IPMC_ENTRY_TYPE_XG;

    default:
      return SAI_IPMC_ENTRY_TYPE_SG;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_ipmc_entry_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipmc_entry_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipmc_entry_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipmc_entry_type_t_to_sai(
        static_cast<lemming::dataplane::sai::IpmcEntryType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpmcGroupAttr convert_sai_ipmc_group_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IPMC_GROUP_ATTR_IPMC_OUTPUT_COUNT:
      return lemming::dataplane::sai::IPMC_GROUP_ATTR_IPMC_OUTPUT_COUNT;

    case SAI_IPMC_GROUP_ATTR_IPMC_MEMBER_LIST:
      return lemming::dataplane::sai::IPMC_GROUP_ATTR_IPMC_MEMBER_LIST;

    default:
      return lemming::dataplane::sai::IPMC_GROUP_ATTR_UNSPECIFIED;
  }
}
sai_ipmc_group_attr_t convert_sai_ipmc_group_attr_t_to_sai(
    lemming::dataplane::sai::IpmcGroupAttr val) {
  switch (val) {
    case lemming::dataplane::sai::IPMC_GROUP_ATTR_IPMC_OUTPUT_COUNT:
      return SAI_IPMC_GROUP_ATTR_IPMC_OUTPUT_COUNT;

    case lemming::dataplane::sai::IPMC_GROUP_ATTR_IPMC_MEMBER_LIST:
      return SAI_IPMC_GROUP_ATTR_IPMC_MEMBER_LIST;

    default:
      return SAI_IPMC_GROUP_ATTR_IPMC_OUTPUT_COUNT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_ipmc_group_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipmc_group_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipmc_group_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipmc_group_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::IpmcGroupAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpmcGroupMemberAttr
convert_sai_ipmc_group_member_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_IPMC_GROUP_MEMBER_ATTR_IPMC_GROUP_ID:
      return lemming::dataplane::sai::IPMC_GROUP_MEMBER_ATTR_IPMC_GROUP_ID;

    case SAI_IPMC_GROUP_MEMBER_ATTR_IPMC_OUTPUT_ID:
      return lemming::dataplane::sai::IPMC_GROUP_MEMBER_ATTR_IPMC_OUTPUT_ID;

    default:
      return lemming::dataplane::sai::IPMC_GROUP_MEMBER_ATTR_UNSPECIFIED;
  }
}
sai_ipmc_group_member_attr_t convert_sai_ipmc_group_member_attr_t_to_sai(
    lemming::dataplane::sai::IpmcGroupMemberAttr val) {
  switch (val) {
    case lemming::dataplane::sai::IPMC_GROUP_MEMBER_ATTR_IPMC_GROUP_ID:
      return SAI_IPMC_GROUP_MEMBER_ATTR_IPMC_GROUP_ID;

    case lemming::dataplane::sai::IPMC_GROUP_MEMBER_ATTR_IPMC_OUTPUT_ID:
      return SAI_IPMC_GROUP_MEMBER_ATTR_IPMC_OUTPUT_ID;

    default:
      return SAI_IPMC_GROUP_MEMBER_ATTR_IPMC_GROUP_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_ipmc_group_member_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipmc_group_member_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipmc_group_member_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipmc_group_member_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::IpmcGroupMemberAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpsecAttr convert_sai_ipsec_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IPSEC_ATTR_TERM_REMOTE_IP_MATCH_SUPPORTED:
      return lemming::dataplane::sai::IPSEC_ATTR_TERM_REMOTE_IP_MATCH_SUPPORTED;

    case SAI_IPSEC_ATTR_SWITCHING_MODE_CUT_THROUGH_SUPPORTED:
      return lemming::dataplane::sai::
          IPSEC_ATTR_SWITCHING_MODE_CUT_THROUGH_SUPPORTED;

    case SAI_IPSEC_ATTR_SWITCHING_MODE_STORE_AND_FORWARD_SUPPORTED:
      return lemming::dataplane::sai::
          IPSEC_ATTR_SWITCHING_MODE_STORE_AND_FORWARD_SUPPORTED;

    case SAI_IPSEC_ATTR_STATS_MODE_READ_SUPPORTED:
      return lemming::dataplane::sai::IPSEC_ATTR_STATS_MODE_READ_SUPPORTED;

    case SAI_IPSEC_ATTR_STATS_MODE_READ_CLEAR_SUPPORTED:
      return lemming::dataplane::sai::
          IPSEC_ATTR_STATS_MODE_READ_CLEAR_SUPPORTED;

    case SAI_IPSEC_ATTR_SN_32BIT_SUPPORTED:
      return lemming::dataplane::sai::IPSEC_ATTR_SN_32BIT_SUPPORTED;

    case SAI_IPSEC_ATTR_ESN_64BIT_SUPPORTED:
      return lemming::dataplane::sai::IPSEC_ATTR_ESN_64BIT_SUPPORTED;

    case SAI_IPSEC_ATTR_SUPPORTED_CIPHER_LIST:
      return lemming::dataplane::sai::IPSEC_ATTR_SUPPORTED_CIPHER_LIST;

    case SAI_IPSEC_ATTR_SYSTEM_SIDE_MTU:
      return lemming::dataplane::sai::IPSEC_ATTR_SYSTEM_SIDE_MTU;

    case SAI_IPSEC_ATTR_WARM_BOOT_SUPPORTED:
      return lemming::dataplane::sai::IPSEC_ATTR_WARM_BOOT_SUPPORTED;

    case SAI_IPSEC_ATTR_WARM_BOOT_ENABLE:
      return lemming::dataplane::sai::IPSEC_ATTR_WARM_BOOT_ENABLE;

    case SAI_IPSEC_ATTR_EXTERNAL_SA_INDEX_ENABLE:
      return lemming::dataplane::sai::IPSEC_ATTR_EXTERNAL_SA_INDEX_ENABLE;

    case SAI_IPSEC_ATTR_CTAG_TPID:
      return lemming::dataplane::sai::IPSEC_ATTR_CTAG_TPID;

    case SAI_IPSEC_ATTR_STAG_TPID:
      return lemming::dataplane::sai::IPSEC_ATTR_STAG_TPID;

    case SAI_IPSEC_ATTR_MAX_VLAN_TAGS_PARSED:
      return lemming::dataplane::sai::IPSEC_ATTR_MAX_VLAN_TAGS_PARSED;

    case SAI_IPSEC_ATTR_OCTET_COUNT_HIGH_WATERMARK:
      return lemming::dataplane::sai::IPSEC_ATTR_OCTET_COUNT_HIGH_WATERMARK;

    case SAI_IPSEC_ATTR_OCTET_COUNT_LOW_WATERMARK:
      return lemming::dataplane::sai::IPSEC_ATTR_OCTET_COUNT_LOW_WATERMARK;

    case SAI_IPSEC_ATTR_STATS_MODE:
      return lemming::dataplane::sai::IPSEC_ATTR_STATS_MODE;

    case SAI_IPSEC_ATTR_AVAILABLE_IPSEC_SA:
      return lemming::dataplane::sai::IPSEC_ATTR_AVAILABLE_IPSEC_SA;

    case SAI_IPSEC_ATTR_SA_LIST:
      return lemming::dataplane::sai::IPSEC_ATTR_SA_LIST;

    default:
      return lemming::dataplane::sai::IPSEC_ATTR_UNSPECIFIED;
  }
}
sai_ipsec_attr_t convert_sai_ipsec_attr_t_to_sai(
    lemming::dataplane::sai::IpsecAttr val) {
  switch (val) {
    case lemming::dataplane::sai::IPSEC_ATTR_TERM_REMOTE_IP_MATCH_SUPPORTED:
      return SAI_IPSEC_ATTR_TERM_REMOTE_IP_MATCH_SUPPORTED;

    case lemming::dataplane::sai::
        IPSEC_ATTR_SWITCHING_MODE_CUT_THROUGH_SUPPORTED:
      return SAI_IPSEC_ATTR_SWITCHING_MODE_CUT_THROUGH_SUPPORTED;

    case lemming::dataplane::sai::
        IPSEC_ATTR_SWITCHING_MODE_STORE_AND_FORWARD_SUPPORTED:
      return SAI_IPSEC_ATTR_SWITCHING_MODE_STORE_AND_FORWARD_SUPPORTED;

    case lemming::dataplane::sai::IPSEC_ATTR_STATS_MODE_READ_SUPPORTED:
      return SAI_IPSEC_ATTR_STATS_MODE_READ_SUPPORTED;

    case lemming::dataplane::sai::IPSEC_ATTR_STATS_MODE_READ_CLEAR_SUPPORTED:
      return SAI_IPSEC_ATTR_STATS_MODE_READ_CLEAR_SUPPORTED;

    case lemming::dataplane::sai::IPSEC_ATTR_SN_32BIT_SUPPORTED:
      return SAI_IPSEC_ATTR_SN_32BIT_SUPPORTED;

    case lemming::dataplane::sai::IPSEC_ATTR_ESN_64BIT_SUPPORTED:
      return SAI_IPSEC_ATTR_ESN_64BIT_SUPPORTED;

    case lemming::dataplane::sai::IPSEC_ATTR_SUPPORTED_CIPHER_LIST:
      return SAI_IPSEC_ATTR_SUPPORTED_CIPHER_LIST;

    case lemming::dataplane::sai::IPSEC_ATTR_SYSTEM_SIDE_MTU:
      return SAI_IPSEC_ATTR_SYSTEM_SIDE_MTU;

    case lemming::dataplane::sai::IPSEC_ATTR_WARM_BOOT_SUPPORTED:
      return SAI_IPSEC_ATTR_WARM_BOOT_SUPPORTED;

    case lemming::dataplane::sai::IPSEC_ATTR_WARM_BOOT_ENABLE:
      return SAI_IPSEC_ATTR_WARM_BOOT_ENABLE;

    case lemming::dataplane::sai::IPSEC_ATTR_EXTERNAL_SA_INDEX_ENABLE:
      return SAI_IPSEC_ATTR_EXTERNAL_SA_INDEX_ENABLE;

    case lemming::dataplane::sai::IPSEC_ATTR_CTAG_TPID:
      return SAI_IPSEC_ATTR_CTAG_TPID;

    case lemming::dataplane::sai::IPSEC_ATTR_STAG_TPID:
      return SAI_IPSEC_ATTR_STAG_TPID;

    case lemming::dataplane::sai::IPSEC_ATTR_MAX_VLAN_TAGS_PARSED:
      return SAI_IPSEC_ATTR_MAX_VLAN_TAGS_PARSED;

    case lemming::dataplane::sai::IPSEC_ATTR_OCTET_COUNT_HIGH_WATERMARK:
      return SAI_IPSEC_ATTR_OCTET_COUNT_HIGH_WATERMARK;

    case lemming::dataplane::sai::IPSEC_ATTR_OCTET_COUNT_LOW_WATERMARK:
      return SAI_IPSEC_ATTR_OCTET_COUNT_LOW_WATERMARK;

    case lemming::dataplane::sai::IPSEC_ATTR_STATS_MODE:
      return SAI_IPSEC_ATTR_STATS_MODE;

    case lemming::dataplane::sai::IPSEC_ATTR_AVAILABLE_IPSEC_SA:
      return SAI_IPSEC_ATTR_AVAILABLE_IPSEC_SA;

    case lemming::dataplane::sai::IPSEC_ATTR_SA_LIST:
      return SAI_IPSEC_ATTR_SA_LIST;

    default:
      return SAI_IPSEC_ATTR_TERM_REMOTE_IP_MATCH_SUPPORTED;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_ipsec_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipsec_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipsec_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipsec_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::IpsecAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpsecCipher convert_sai_ipsec_cipher_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IPSEC_CIPHER_AES128_GCM16:
      return lemming::dataplane::sai::IPSEC_CIPHER_AES128_GCM16;

    case SAI_IPSEC_CIPHER_AES256_GCM16:
      return lemming::dataplane::sai::IPSEC_CIPHER_AES256_GCM16;

    case SAI_IPSEC_CIPHER_AES128_GMAC:
      return lemming::dataplane::sai::IPSEC_CIPHER_AES128_GMAC;

    case SAI_IPSEC_CIPHER_AES256_GMAC:
      return lemming::dataplane::sai::IPSEC_CIPHER_AES256_GMAC;

    default:
      return lemming::dataplane::sai::IPSEC_CIPHER_UNSPECIFIED;
  }
}
sai_ipsec_cipher_t convert_sai_ipsec_cipher_t_to_sai(
    lemming::dataplane::sai::IpsecCipher val) {
  switch (val) {
    case lemming::dataplane::sai::IPSEC_CIPHER_AES128_GCM16:
      return SAI_IPSEC_CIPHER_AES128_GCM16;

    case lemming::dataplane::sai::IPSEC_CIPHER_AES256_GCM16:
      return SAI_IPSEC_CIPHER_AES256_GCM16;

    case lemming::dataplane::sai::IPSEC_CIPHER_AES128_GMAC:
      return SAI_IPSEC_CIPHER_AES128_GMAC;

    case lemming::dataplane::sai::IPSEC_CIPHER_AES256_GMAC:
      return SAI_IPSEC_CIPHER_AES256_GMAC;

    default:
      return SAI_IPSEC_CIPHER_AES128_GCM16;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_ipsec_cipher_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipsec_cipher_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipsec_cipher_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipsec_cipher_t_to_sai(
        static_cast<lemming::dataplane::sai::IpsecCipher>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpsecDirection convert_sai_ipsec_direction_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IPSEC_DIRECTION_EGRESS:
      return lemming::dataplane::sai::IPSEC_DIRECTION_EGRESS;

    case SAI_IPSEC_DIRECTION_INGRESS:
      return lemming::dataplane::sai::IPSEC_DIRECTION_INGRESS;

    default:
      return lemming::dataplane::sai::IPSEC_DIRECTION_UNSPECIFIED;
  }
}
sai_ipsec_direction_t convert_sai_ipsec_direction_t_to_sai(
    lemming::dataplane::sai::IpsecDirection val) {
  switch (val) {
    case lemming::dataplane::sai::IPSEC_DIRECTION_EGRESS:
      return SAI_IPSEC_DIRECTION_EGRESS;

    case lemming::dataplane::sai::IPSEC_DIRECTION_INGRESS:
      return SAI_IPSEC_DIRECTION_INGRESS;

    default:
      return SAI_IPSEC_DIRECTION_EGRESS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_ipsec_direction_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipsec_direction_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipsec_direction_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipsec_direction_t_to_sai(
        static_cast<lemming::dataplane::sai::IpsecDirection>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpsecPortAttr convert_sai_ipsec_port_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IPSEC_PORT_ATTR_PORT_ID:
      return lemming::dataplane::sai::IPSEC_PORT_ATTR_PORT_ID;

    case SAI_IPSEC_PORT_ATTR_CTAG_ENABLE:
      return lemming::dataplane::sai::IPSEC_PORT_ATTR_CTAG_ENABLE;

    case SAI_IPSEC_PORT_ATTR_STAG_ENABLE:
      return lemming::dataplane::sai::IPSEC_PORT_ATTR_STAG_ENABLE;

    case SAI_IPSEC_PORT_ATTR_NATIVE_VLAN_ID:
      return lemming::dataplane::sai::IPSEC_PORT_ATTR_NATIVE_VLAN_ID;

    case SAI_IPSEC_PORT_ATTR_VRF_FROM_PACKET_VLAN_ENABLE:
      return lemming::dataplane::sai::
          IPSEC_PORT_ATTR_VRF_FROM_PACKET_VLAN_ENABLE;

    case SAI_IPSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
      return lemming::dataplane::sai::IPSEC_PORT_ATTR_SWITCH_SWITCHING_MODE;

    default:
      return lemming::dataplane::sai::IPSEC_PORT_ATTR_UNSPECIFIED;
  }
}
sai_ipsec_port_attr_t convert_sai_ipsec_port_attr_t_to_sai(
    lemming::dataplane::sai::IpsecPortAttr val) {
  switch (val) {
    case lemming::dataplane::sai::IPSEC_PORT_ATTR_PORT_ID:
      return SAI_IPSEC_PORT_ATTR_PORT_ID;

    case lemming::dataplane::sai::IPSEC_PORT_ATTR_CTAG_ENABLE:
      return SAI_IPSEC_PORT_ATTR_CTAG_ENABLE;

    case lemming::dataplane::sai::IPSEC_PORT_ATTR_STAG_ENABLE:
      return SAI_IPSEC_PORT_ATTR_STAG_ENABLE;

    case lemming::dataplane::sai::IPSEC_PORT_ATTR_NATIVE_VLAN_ID:
      return SAI_IPSEC_PORT_ATTR_NATIVE_VLAN_ID;

    case lemming::dataplane::sai::IPSEC_PORT_ATTR_VRF_FROM_PACKET_VLAN_ENABLE:
      return SAI_IPSEC_PORT_ATTR_VRF_FROM_PACKET_VLAN_ENABLE;

    case lemming::dataplane::sai::IPSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
      return SAI_IPSEC_PORT_ATTR_SWITCH_SWITCHING_MODE;

    default:
      return SAI_IPSEC_PORT_ATTR_PORT_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_ipsec_port_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipsec_port_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipsec_port_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipsec_port_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::IpsecPortAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpsecPortStat convert_sai_ipsec_port_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IPSEC_PORT_STAT_TX_ERROR_PKTS:
      return lemming::dataplane::sai::IPSEC_PORT_STAT_TX_ERROR_PKTS;

    case SAI_IPSEC_PORT_STAT_TX_IPSEC_PKTS:
      return lemming::dataplane::sai::IPSEC_PORT_STAT_TX_IPSEC_PKTS;

    case SAI_IPSEC_PORT_STAT_TX_NON_IPSEC_PKTS:
      return lemming::dataplane::sai::IPSEC_PORT_STAT_TX_NON_IPSEC_PKTS;

    case SAI_IPSEC_PORT_STAT_RX_ERROR_PKTS:
      return lemming::dataplane::sai::IPSEC_PORT_STAT_RX_ERROR_PKTS;

    case SAI_IPSEC_PORT_STAT_RX_IPSEC_PKTS:
      return lemming::dataplane::sai::IPSEC_PORT_STAT_RX_IPSEC_PKTS;

    case SAI_IPSEC_PORT_STAT_RX_NON_IPSEC_PKTS:
      return lemming::dataplane::sai::IPSEC_PORT_STAT_RX_NON_IPSEC_PKTS;

    default:
      return lemming::dataplane::sai::IPSEC_PORT_STAT_UNSPECIFIED;
  }
}
sai_ipsec_port_stat_t convert_sai_ipsec_port_stat_t_to_sai(
    lemming::dataplane::sai::IpsecPortStat val) {
  switch (val) {
    case lemming::dataplane::sai::IPSEC_PORT_STAT_TX_ERROR_PKTS:
      return SAI_IPSEC_PORT_STAT_TX_ERROR_PKTS;

    case lemming::dataplane::sai::IPSEC_PORT_STAT_TX_IPSEC_PKTS:
      return SAI_IPSEC_PORT_STAT_TX_IPSEC_PKTS;

    case lemming::dataplane::sai::IPSEC_PORT_STAT_TX_NON_IPSEC_PKTS:
      return SAI_IPSEC_PORT_STAT_TX_NON_IPSEC_PKTS;

    case lemming::dataplane::sai::IPSEC_PORT_STAT_RX_ERROR_PKTS:
      return SAI_IPSEC_PORT_STAT_RX_ERROR_PKTS;

    case lemming::dataplane::sai::IPSEC_PORT_STAT_RX_IPSEC_PKTS:
      return SAI_IPSEC_PORT_STAT_RX_IPSEC_PKTS;

    case lemming::dataplane::sai::IPSEC_PORT_STAT_RX_NON_IPSEC_PKTS:
      return SAI_IPSEC_PORT_STAT_RX_NON_IPSEC_PKTS;

    default:
      return SAI_IPSEC_PORT_STAT_TX_ERROR_PKTS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_ipsec_port_stat_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipsec_port_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipsec_port_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipsec_port_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::IpsecPortStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpsecSaAttr convert_sai_ipsec_sa_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IPSEC_SA_ATTR_IPSEC_DIRECTION:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_DIRECTION;

    case SAI_IPSEC_SA_ATTR_IPSEC_ID:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_ID;

    case SAI_IPSEC_SA_ATTR_OCTET_COUNT_STATUS:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_OCTET_COUNT_STATUS;

    case SAI_IPSEC_SA_ATTR_EXTERNAL_SA_INDEX:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_EXTERNAL_SA_INDEX;

    case SAI_IPSEC_SA_ATTR_SA_INDEX:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_SA_INDEX;

    case SAI_IPSEC_SA_ATTR_IPSEC_PORT_LIST:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_PORT_LIST;

    case SAI_IPSEC_SA_ATTR_IPSEC_SPI:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_SPI;

    case SAI_IPSEC_SA_ATTR_IPSEC_ESN_ENABLE:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_ESN_ENABLE;

    case SAI_IPSEC_SA_ATTR_IPSEC_CIPHER:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_CIPHER;

    case SAI_IPSEC_SA_ATTR_ENCRYPT_KEY:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_ENCRYPT_KEY;

    case SAI_IPSEC_SA_ATTR_SALT:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_SALT;

    case SAI_IPSEC_SA_ATTR_AUTH_KEY:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_AUTH_KEY;

    case SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_ENABLE:
      return lemming::dataplane::sai::
          IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_ENABLE;

    case SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_WINDOW:
      return lemming::dataplane::sai::
          IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_WINDOW;

    case SAI_IPSEC_SA_ATTR_TERM_DST_IP:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_TERM_DST_IP;

    case SAI_IPSEC_SA_ATTR_TERM_VLAN_ID_ENABLE:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_TERM_VLAN_ID_ENABLE;

    case SAI_IPSEC_SA_ATTR_TERM_VLAN_ID:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_TERM_VLAN_ID;

    case SAI_IPSEC_SA_ATTR_TERM_SRC_IP_ENABLE:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_TERM_SRC_IP_ENABLE;

    case SAI_IPSEC_SA_ATTR_TERM_SRC_IP:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_TERM_SRC_IP;

    case SAI_IPSEC_SA_ATTR_EGRESS_ESN:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_EGRESS_ESN;

    case SAI_IPSEC_SA_ATTR_MINIMUM_INGRESS_ESN:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_MINIMUM_INGRESS_ESN;

    default:
      return lemming::dataplane::sai::IPSEC_SA_ATTR_UNSPECIFIED;
  }
}
sai_ipsec_sa_attr_t convert_sai_ipsec_sa_attr_t_to_sai(
    lemming::dataplane::sai::IpsecSaAttr val) {
  switch (val) {
    case lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_DIRECTION:
      return SAI_IPSEC_SA_ATTR_IPSEC_DIRECTION;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_ID:
      return SAI_IPSEC_SA_ATTR_IPSEC_ID;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_OCTET_COUNT_STATUS:
      return SAI_IPSEC_SA_ATTR_OCTET_COUNT_STATUS;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_EXTERNAL_SA_INDEX:
      return SAI_IPSEC_SA_ATTR_EXTERNAL_SA_INDEX;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_SA_INDEX:
      return SAI_IPSEC_SA_ATTR_SA_INDEX;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_PORT_LIST:
      return SAI_IPSEC_SA_ATTR_IPSEC_PORT_LIST;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_SPI:
      return SAI_IPSEC_SA_ATTR_IPSEC_SPI;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_ESN_ENABLE:
      return SAI_IPSEC_SA_ATTR_IPSEC_ESN_ENABLE;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_CIPHER:
      return SAI_IPSEC_SA_ATTR_IPSEC_CIPHER;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_ENCRYPT_KEY:
      return SAI_IPSEC_SA_ATTR_ENCRYPT_KEY;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_SALT:
      return SAI_IPSEC_SA_ATTR_SALT;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_AUTH_KEY:
      return SAI_IPSEC_SA_ATTR_AUTH_KEY;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_ENABLE:
      return SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_ENABLE;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_WINDOW:
      return SAI_IPSEC_SA_ATTR_IPSEC_REPLAY_PROTECTION_WINDOW;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_TERM_DST_IP:
      return SAI_IPSEC_SA_ATTR_TERM_DST_IP;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_TERM_VLAN_ID_ENABLE:
      return SAI_IPSEC_SA_ATTR_TERM_VLAN_ID_ENABLE;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_TERM_VLAN_ID:
      return SAI_IPSEC_SA_ATTR_TERM_VLAN_ID;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_TERM_SRC_IP_ENABLE:
      return SAI_IPSEC_SA_ATTR_TERM_SRC_IP_ENABLE;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_TERM_SRC_IP:
      return SAI_IPSEC_SA_ATTR_TERM_SRC_IP;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_EGRESS_ESN:
      return SAI_IPSEC_SA_ATTR_EGRESS_ESN;

    case lemming::dataplane::sai::IPSEC_SA_ATTR_MINIMUM_INGRESS_ESN:
      return SAI_IPSEC_SA_ATTR_MINIMUM_INGRESS_ESN;

    default:
      return SAI_IPSEC_SA_ATTR_IPSEC_DIRECTION;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_ipsec_sa_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipsec_sa_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipsec_sa_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipsec_sa_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::IpsecSaAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpsecSaOctetCountStatus
convert_sai_ipsec_sa_octet_count_status_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_IPSEC_SA_OCTET_COUNT_STATUS_BELOW_LOW_WATERMARK:
      return lemming::dataplane::sai::
          IPSEC_SA_OCTET_COUNT_STATUS_BELOW_LOW_WATERMARK;

    case SAI_IPSEC_SA_OCTET_COUNT_STATUS_BELOW_HIGH_WATERMARK:
      return lemming::dataplane::sai::
          IPSEC_SA_OCTET_COUNT_STATUS_BELOW_HIGH_WATERMARK;

    case SAI_IPSEC_SA_OCTET_COUNT_STATUS_ABOVE_HIGH_WATERMARK:
      return lemming::dataplane::sai::
          IPSEC_SA_OCTET_COUNT_STATUS_ABOVE_HIGH_WATERMARK;

    default:
      return lemming::dataplane::sai::IPSEC_SA_OCTET_COUNT_STATUS_UNSPECIFIED;
  }
}
sai_ipsec_sa_octet_count_status_t
convert_sai_ipsec_sa_octet_count_status_t_to_sai(
    lemming::dataplane::sai::IpsecSaOctetCountStatus val) {
  switch (val) {
    case lemming::dataplane::sai::
        IPSEC_SA_OCTET_COUNT_STATUS_BELOW_LOW_WATERMARK:
      return SAI_IPSEC_SA_OCTET_COUNT_STATUS_BELOW_LOW_WATERMARK;

    case lemming::dataplane::sai::
        IPSEC_SA_OCTET_COUNT_STATUS_BELOW_HIGH_WATERMARK:
      return SAI_IPSEC_SA_OCTET_COUNT_STATUS_BELOW_HIGH_WATERMARK;

    case lemming::dataplane::sai::
        IPSEC_SA_OCTET_COUNT_STATUS_ABOVE_HIGH_WATERMARK:
      return SAI_IPSEC_SA_OCTET_COUNT_STATUS_ABOVE_HIGH_WATERMARK;

    default:
      return SAI_IPSEC_SA_OCTET_COUNT_STATUS_BELOW_LOW_WATERMARK;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_ipsec_sa_octet_count_status_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_ipsec_sa_octet_count_status_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipsec_sa_octet_count_status_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipsec_sa_octet_count_status_t_to_sai(
        static_cast<lemming::dataplane::sai::IpsecSaOctetCountStatus>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IpsecSaStat convert_sai_ipsec_sa_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_IPSEC_SA_STAT_PROTECTED_OCTETS:
      return lemming::dataplane::sai::IPSEC_SA_STAT_PROTECTED_OCTETS;

    case SAI_IPSEC_SA_STAT_PROTECTED_PKTS:
      return lemming::dataplane::sai::IPSEC_SA_STAT_PROTECTED_PKTS;

    case SAI_IPSEC_SA_STAT_GOOD_PKTS:
      return lemming::dataplane::sai::IPSEC_SA_STAT_GOOD_PKTS;

    case SAI_IPSEC_SA_STAT_BAD_HEADER_PKTS_IN:
      return lemming::dataplane::sai::IPSEC_SA_STAT_BAD_HEADER_PKTS_IN;

    case SAI_IPSEC_SA_STAT_REPLAYED_PKTS_IN:
      return lemming::dataplane::sai::IPSEC_SA_STAT_REPLAYED_PKTS_IN;

    case SAI_IPSEC_SA_STAT_LATE_PKTS_IN:
      return lemming::dataplane::sai::IPSEC_SA_STAT_LATE_PKTS_IN;

    case SAI_IPSEC_SA_STAT_BAD_TRAILER_PKTS_IN:
      return lemming::dataplane::sai::IPSEC_SA_STAT_BAD_TRAILER_PKTS_IN;

    case SAI_IPSEC_SA_STAT_AUTH_FAIL_PKTS_IN:
      return lemming::dataplane::sai::IPSEC_SA_STAT_AUTH_FAIL_PKTS_IN;

    case SAI_IPSEC_SA_STAT_DUMMY_DROPPED_PKTS_IN:
      return lemming::dataplane::sai::IPSEC_SA_STAT_DUMMY_DROPPED_PKTS_IN;

    case SAI_IPSEC_SA_STAT_OTHER_DROPPED_PKTS:
      return lemming::dataplane::sai::IPSEC_SA_STAT_OTHER_DROPPED_PKTS;

    default:
      return lemming::dataplane::sai::IPSEC_SA_STAT_UNSPECIFIED;
  }
}
sai_ipsec_sa_stat_t convert_sai_ipsec_sa_stat_t_to_sai(
    lemming::dataplane::sai::IpsecSaStat val) {
  switch (val) {
    case lemming::dataplane::sai::IPSEC_SA_STAT_PROTECTED_OCTETS:
      return SAI_IPSEC_SA_STAT_PROTECTED_OCTETS;

    case lemming::dataplane::sai::IPSEC_SA_STAT_PROTECTED_PKTS:
      return SAI_IPSEC_SA_STAT_PROTECTED_PKTS;

    case lemming::dataplane::sai::IPSEC_SA_STAT_GOOD_PKTS:
      return SAI_IPSEC_SA_STAT_GOOD_PKTS;

    case lemming::dataplane::sai::IPSEC_SA_STAT_BAD_HEADER_PKTS_IN:
      return SAI_IPSEC_SA_STAT_BAD_HEADER_PKTS_IN;

    case lemming::dataplane::sai::IPSEC_SA_STAT_REPLAYED_PKTS_IN:
      return SAI_IPSEC_SA_STAT_REPLAYED_PKTS_IN;

    case lemming::dataplane::sai::IPSEC_SA_STAT_LATE_PKTS_IN:
      return SAI_IPSEC_SA_STAT_LATE_PKTS_IN;

    case lemming::dataplane::sai::IPSEC_SA_STAT_BAD_TRAILER_PKTS_IN:
      return SAI_IPSEC_SA_STAT_BAD_TRAILER_PKTS_IN;

    case lemming::dataplane::sai::IPSEC_SA_STAT_AUTH_FAIL_PKTS_IN:
      return SAI_IPSEC_SA_STAT_AUTH_FAIL_PKTS_IN;

    case lemming::dataplane::sai::IPSEC_SA_STAT_DUMMY_DROPPED_PKTS_IN:
      return SAI_IPSEC_SA_STAT_DUMMY_DROPPED_PKTS_IN;

    case lemming::dataplane::sai::IPSEC_SA_STAT_OTHER_DROPPED_PKTS:
      return SAI_IPSEC_SA_STAT_OTHER_DROPPED_PKTS;

    default:
      return SAI_IPSEC_SA_STAT_PROTECTED_OCTETS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_ipsec_sa_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_ipsec_sa_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_ipsec_sa_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_ipsec_sa_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::IpsecSaStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IsolationGroupAttr
convert_sai_isolation_group_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ISOLATION_GROUP_ATTR_TYPE:
      return lemming::dataplane::sai::ISOLATION_GROUP_ATTR_TYPE;

    case SAI_ISOLATION_GROUP_ATTR_ISOLATION_MEMBER_LIST:
      return lemming::dataplane::sai::
          ISOLATION_GROUP_ATTR_ISOLATION_MEMBER_LIST;

    default:
      return lemming::dataplane::sai::ISOLATION_GROUP_ATTR_UNSPECIFIED;
  }
}
sai_isolation_group_attr_t convert_sai_isolation_group_attr_t_to_sai(
    lemming::dataplane::sai::IsolationGroupAttr val) {
  switch (val) {
    case lemming::dataplane::sai::ISOLATION_GROUP_ATTR_TYPE:
      return SAI_ISOLATION_GROUP_ATTR_TYPE;

    case lemming::dataplane::sai::ISOLATION_GROUP_ATTR_ISOLATION_MEMBER_LIST:
      return SAI_ISOLATION_GROUP_ATTR_ISOLATION_MEMBER_LIST;

    default:
      return SAI_ISOLATION_GROUP_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_isolation_group_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_isolation_group_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_isolation_group_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_isolation_group_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::IsolationGroupAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IsolationGroupMemberAttr
convert_sai_isolation_group_member_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_GROUP_ID:
      return lemming::dataplane::sai::
          ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_GROUP_ID;

    case SAI_ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_OBJECT:
      return lemming::dataplane::sai::
          ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_OBJECT;

    default:
      return lemming::dataplane::sai::ISOLATION_GROUP_MEMBER_ATTR_UNSPECIFIED;
  }
}
sai_isolation_group_member_attr_t
convert_sai_isolation_group_member_attr_t_to_sai(
    lemming::dataplane::sai::IsolationGroupMemberAttr val) {
  switch (val) {
    case lemming::dataplane::sai::
        ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_GROUP_ID:
      return SAI_ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_GROUP_ID;

    case lemming::dataplane::sai::ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_OBJECT:
      return SAI_ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_OBJECT;

    default:
      return SAI_ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_GROUP_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_isolation_group_member_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_isolation_group_member_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_isolation_group_member_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_isolation_group_member_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::IsolationGroupMemberAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::IsolationGroupType
convert_sai_isolation_group_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ISOLATION_GROUP_TYPE_PORT:
      return lemming::dataplane::sai::ISOLATION_GROUP_TYPE_PORT;

    case SAI_ISOLATION_GROUP_TYPE_BRIDGE_PORT:
      return lemming::dataplane::sai::ISOLATION_GROUP_TYPE_BRIDGE_PORT;

    default:
      return lemming::dataplane::sai::ISOLATION_GROUP_TYPE_UNSPECIFIED;
  }
}
sai_isolation_group_type_t convert_sai_isolation_group_type_t_to_sai(
    lemming::dataplane::sai::IsolationGroupType val) {
  switch (val) {
    case lemming::dataplane::sai::ISOLATION_GROUP_TYPE_PORT:
      return SAI_ISOLATION_GROUP_TYPE_PORT;

    case lemming::dataplane::sai::ISOLATION_GROUP_TYPE_BRIDGE_PORT:
      return SAI_ISOLATION_GROUP_TYPE_BRIDGE_PORT;

    default:
      return SAI_ISOLATION_GROUP_TYPE_PORT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_isolation_group_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_isolation_group_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_isolation_group_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_isolation_group_type_t_to_sai(
        static_cast<lemming::dataplane::sai::IsolationGroupType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::L2mcEntryAttr convert_sai_l2mc_entry_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_L2MC_ENTRY_ATTR_PACKET_ACTION:
      return lemming::dataplane::sai::L2MC_ENTRY_ATTR_PACKET_ACTION;

    case SAI_L2MC_ENTRY_ATTR_OUTPUT_GROUP_ID:
      return lemming::dataplane::sai::L2MC_ENTRY_ATTR_OUTPUT_GROUP_ID;

    default:
      return lemming::dataplane::sai::L2MC_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_l2mc_entry_attr_t convert_sai_l2mc_entry_attr_t_to_sai(
    lemming::dataplane::sai::L2mcEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::L2MC_ENTRY_ATTR_PACKET_ACTION:
      return SAI_L2MC_ENTRY_ATTR_PACKET_ACTION;

    case lemming::dataplane::sai::L2MC_ENTRY_ATTR_OUTPUT_GROUP_ID:
      return SAI_L2MC_ENTRY_ATTR_OUTPUT_GROUP_ID;

    default:
      return SAI_L2MC_ENTRY_ATTR_PACKET_ACTION;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_l2mc_entry_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_l2mc_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_l2mc_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_l2mc_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::L2mcEntryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::L2mcEntryType convert_sai_l2mc_entry_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_L2MC_ENTRY_TYPE_SG:
      return lemming::dataplane::sai::L2MC_ENTRY_TYPE_SG;

    case SAI_L2MC_ENTRY_TYPE_XG:
      return lemming::dataplane::sai::L2MC_ENTRY_TYPE_XG;

    default:
      return lemming::dataplane::sai::L2MC_ENTRY_TYPE_UNSPECIFIED;
  }
}
sai_l2mc_entry_type_t convert_sai_l2mc_entry_type_t_to_sai(
    lemming::dataplane::sai::L2mcEntryType val) {
  switch (val) {
    case lemming::dataplane::sai::L2MC_ENTRY_TYPE_SG:
      return SAI_L2MC_ENTRY_TYPE_SG;

    case lemming::dataplane::sai::L2MC_ENTRY_TYPE_XG:
      return SAI_L2MC_ENTRY_TYPE_XG;

    default:
      return SAI_L2MC_ENTRY_TYPE_SG;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_l2mc_entry_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_l2mc_entry_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_l2mc_entry_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_l2mc_entry_type_t_to_sai(
        static_cast<lemming::dataplane::sai::L2mcEntryType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::L2mcGroupAttr convert_sai_l2mc_group_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_L2MC_GROUP_ATTR_L2MC_OUTPUT_COUNT:
      return lemming::dataplane::sai::L2MC_GROUP_ATTR_L2MC_OUTPUT_COUNT;

    case SAI_L2MC_GROUP_ATTR_L2MC_MEMBER_LIST:
      return lemming::dataplane::sai::L2MC_GROUP_ATTR_L2MC_MEMBER_LIST;

    default:
      return lemming::dataplane::sai::L2MC_GROUP_ATTR_UNSPECIFIED;
  }
}
sai_l2mc_group_attr_t convert_sai_l2mc_group_attr_t_to_sai(
    lemming::dataplane::sai::L2mcGroupAttr val) {
  switch (val) {
    case lemming::dataplane::sai::L2MC_GROUP_ATTR_L2MC_OUTPUT_COUNT:
      return SAI_L2MC_GROUP_ATTR_L2MC_OUTPUT_COUNT;

    case lemming::dataplane::sai::L2MC_GROUP_ATTR_L2MC_MEMBER_LIST:
      return SAI_L2MC_GROUP_ATTR_L2MC_MEMBER_LIST;

    default:
      return SAI_L2MC_GROUP_ATTR_L2MC_OUTPUT_COUNT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_l2mc_group_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_l2mc_group_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_l2mc_group_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_l2mc_group_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::L2mcGroupAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::L2mcGroupMemberAttr
convert_sai_l2mc_group_member_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_GROUP_ID:
      return lemming::dataplane::sai::L2MC_GROUP_MEMBER_ATTR_L2MC_GROUP_ID;

    case SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_OUTPUT_ID:
      return lemming::dataplane::sai::L2MC_GROUP_MEMBER_ATTR_L2MC_OUTPUT_ID;

    case SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_ENDPOINT_IP:
      return lemming::dataplane::sai::L2MC_GROUP_MEMBER_ATTR_L2MC_ENDPOINT_IP;

    default:
      return lemming::dataplane::sai::L2MC_GROUP_MEMBER_ATTR_UNSPECIFIED;
  }
}
sai_l2mc_group_member_attr_t convert_sai_l2mc_group_member_attr_t_to_sai(
    lemming::dataplane::sai::L2mcGroupMemberAttr val) {
  switch (val) {
    case lemming::dataplane::sai::L2MC_GROUP_MEMBER_ATTR_L2MC_GROUP_ID:
      return SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_GROUP_ID;

    case lemming::dataplane::sai::L2MC_GROUP_MEMBER_ATTR_L2MC_OUTPUT_ID:
      return SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_OUTPUT_ID;

    case lemming::dataplane::sai::L2MC_GROUP_MEMBER_ATTR_L2MC_ENDPOINT_IP:
      return SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_ENDPOINT_IP;

    default:
      return SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_GROUP_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_l2mc_group_member_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_l2mc_group_member_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_l2mc_group_member_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_l2mc_group_member_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::L2mcGroupMemberAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::LagAttr convert_sai_lag_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_LAG_ATTR_PORT_LIST:
      return lemming::dataplane::sai::LAG_ATTR_PORT_LIST;

    case SAI_LAG_ATTR_INGRESS_ACL:
      return lemming::dataplane::sai::LAG_ATTR_INGRESS_ACL;

    case SAI_LAG_ATTR_EGRESS_ACL:
      return lemming::dataplane::sai::LAG_ATTR_EGRESS_ACL;

    case SAI_LAG_ATTR_PORT_VLAN_ID:
      return lemming::dataplane::sai::LAG_ATTR_PORT_VLAN_ID;

    case SAI_LAG_ATTR_DEFAULT_VLAN_PRIORITY:
      return lemming::dataplane::sai::LAG_ATTR_DEFAULT_VLAN_PRIORITY;

    case SAI_LAG_ATTR_DROP_UNTAGGED:
      return lemming::dataplane::sai::LAG_ATTR_DROP_UNTAGGED;

    case SAI_LAG_ATTR_DROP_TAGGED:
      return lemming::dataplane::sai::LAG_ATTR_DROP_TAGGED;

    case SAI_LAG_ATTR_TPID:
      return lemming::dataplane::sai::LAG_ATTR_TPID;

    case SAI_LAG_ATTR_SYSTEM_PORT_AGGREGATE_ID:
      return lemming::dataplane::sai::LAG_ATTR_SYSTEM_PORT_AGGREGATE_ID;

    case SAI_LAG_ATTR_LABEL:
      return lemming::dataplane::sai::LAG_ATTR_LABEL;

    default:
      return lemming::dataplane::sai::LAG_ATTR_UNSPECIFIED;
  }
}
sai_lag_attr_t convert_sai_lag_attr_t_to_sai(
    lemming::dataplane::sai::LagAttr val) {
  switch (val) {
    case lemming::dataplane::sai::LAG_ATTR_PORT_LIST:
      return SAI_LAG_ATTR_PORT_LIST;

    case lemming::dataplane::sai::LAG_ATTR_INGRESS_ACL:
      return SAI_LAG_ATTR_INGRESS_ACL;

    case lemming::dataplane::sai::LAG_ATTR_EGRESS_ACL:
      return SAI_LAG_ATTR_EGRESS_ACL;

    case lemming::dataplane::sai::LAG_ATTR_PORT_VLAN_ID:
      return SAI_LAG_ATTR_PORT_VLAN_ID;

    case lemming::dataplane::sai::LAG_ATTR_DEFAULT_VLAN_PRIORITY:
      return SAI_LAG_ATTR_DEFAULT_VLAN_PRIORITY;

    case lemming::dataplane::sai::LAG_ATTR_DROP_UNTAGGED:
      return SAI_LAG_ATTR_DROP_UNTAGGED;

    case lemming::dataplane::sai::LAG_ATTR_DROP_TAGGED:
      return SAI_LAG_ATTR_DROP_TAGGED;

    case lemming::dataplane::sai::LAG_ATTR_TPID:
      return SAI_LAG_ATTR_TPID;

    case lemming::dataplane::sai::LAG_ATTR_SYSTEM_PORT_AGGREGATE_ID:
      return SAI_LAG_ATTR_SYSTEM_PORT_AGGREGATE_ID;

    case lemming::dataplane::sai::LAG_ATTR_LABEL:
      return SAI_LAG_ATTR_LABEL;

    default:
      return SAI_LAG_ATTR_PORT_LIST;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_lag_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_lag_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_lag_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_lag_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::LagAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::LagMemberAttr convert_sai_lag_member_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_LAG_MEMBER_ATTR_LAG_ID:
      return lemming::dataplane::sai::LAG_MEMBER_ATTR_LAG_ID;

    case SAI_LAG_MEMBER_ATTR_PORT_ID:
      return lemming::dataplane::sai::LAG_MEMBER_ATTR_PORT_ID;

    case SAI_LAG_MEMBER_ATTR_EGRESS_DISABLE:
      return lemming::dataplane::sai::LAG_MEMBER_ATTR_EGRESS_DISABLE;

    case SAI_LAG_MEMBER_ATTR_INGRESS_DISABLE:
      return lemming::dataplane::sai::LAG_MEMBER_ATTR_INGRESS_DISABLE;

    default:
      return lemming::dataplane::sai::LAG_MEMBER_ATTR_UNSPECIFIED;
  }
}
sai_lag_member_attr_t convert_sai_lag_member_attr_t_to_sai(
    lemming::dataplane::sai::LagMemberAttr val) {
  switch (val) {
    case lemming::dataplane::sai::LAG_MEMBER_ATTR_LAG_ID:
      return SAI_LAG_MEMBER_ATTR_LAG_ID;

    case lemming::dataplane::sai::LAG_MEMBER_ATTR_PORT_ID:
      return SAI_LAG_MEMBER_ATTR_PORT_ID;

    case lemming::dataplane::sai::LAG_MEMBER_ATTR_EGRESS_DISABLE:
      return SAI_LAG_MEMBER_ATTR_EGRESS_DISABLE;

    case lemming::dataplane::sai::LAG_MEMBER_ATTR_INGRESS_DISABLE:
      return SAI_LAG_MEMBER_ATTR_INGRESS_DISABLE;

    default:
      return SAI_LAG_MEMBER_ATTR_LAG_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_lag_member_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_lag_member_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_lag_member_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_lag_member_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::LagMemberAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::LogLevel convert_sai_log_level_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_LOG_LEVEL_DEBUG:
      return lemming::dataplane::sai::LOG_LEVEL_DEBUG;

    case SAI_LOG_LEVEL_INFO:
      return lemming::dataplane::sai::LOG_LEVEL_INFO;

    case SAI_LOG_LEVEL_NOTICE:
      return lemming::dataplane::sai::LOG_LEVEL_NOTICE;

    case SAI_LOG_LEVEL_WARN:
      return lemming::dataplane::sai::LOG_LEVEL_WARN;

    case SAI_LOG_LEVEL_ERROR:
      return lemming::dataplane::sai::LOG_LEVEL_ERROR;

    case SAI_LOG_LEVEL_CRITICAL:
      return lemming::dataplane::sai::LOG_LEVEL_CRITICAL;

    default:
      return lemming::dataplane::sai::LOG_LEVEL_UNSPECIFIED;
  }
}
sai_log_level_t convert_sai_log_level_t_to_sai(
    lemming::dataplane::sai::LogLevel val) {
  switch (val) {
    case lemming::dataplane::sai::LOG_LEVEL_DEBUG:
      return SAI_LOG_LEVEL_DEBUG;

    case lemming::dataplane::sai::LOG_LEVEL_INFO:
      return SAI_LOG_LEVEL_INFO;

    case lemming::dataplane::sai::LOG_LEVEL_NOTICE:
      return SAI_LOG_LEVEL_NOTICE;

    case lemming::dataplane::sai::LOG_LEVEL_WARN:
      return SAI_LOG_LEVEL_WARN;

    case lemming::dataplane::sai::LOG_LEVEL_ERROR:
      return SAI_LOG_LEVEL_ERROR;

    case lemming::dataplane::sai::LOG_LEVEL_CRITICAL:
      return SAI_LOG_LEVEL_CRITICAL;

    default:
      return SAI_LOG_LEVEL_DEBUG;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_log_level_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_log_level_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_log_level_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_log_level_t_to_sai(
        static_cast<lemming::dataplane::sai::LogLevel>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecAttr convert_sai_macsec_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_ATTR_DIRECTION:
      return lemming::dataplane::sai::MACSEC_ATTR_DIRECTION;

    case SAI_MACSEC_ATTR_SWITCHING_MODE_CUT_THROUGH_SUPPORTED:
      return lemming::dataplane::sai::
          MACSEC_ATTR_SWITCHING_MODE_CUT_THROUGH_SUPPORTED;

    case SAI_MACSEC_ATTR_SWITCHING_MODE_STORE_AND_FORWARD_SUPPORTED:
      return lemming::dataplane::sai::
          MACSEC_ATTR_SWITCHING_MODE_STORE_AND_FORWARD_SUPPORTED;

    case SAI_MACSEC_ATTR_STATS_MODE_READ_SUPPORTED:
      return lemming::dataplane::sai::MACSEC_ATTR_STATS_MODE_READ_SUPPORTED;

    case SAI_MACSEC_ATTR_STATS_MODE_READ_CLEAR_SUPPORTED:
      return lemming::dataplane::sai::
          MACSEC_ATTR_STATS_MODE_READ_CLEAR_SUPPORTED;

    case SAI_MACSEC_ATTR_SCI_IN_INGRESS_MACSEC_ACL:
      return lemming::dataplane::sai::MACSEC_ATTR_SCI_IN_INGRESS_MACSEC_ACL;

    case SAI_MACSEC_ATTR_SUPPORTED_CIPHER_SUITE_LIST:
      return lemming::dataplane::sai::MACSEC_ATTR_SUPPORTED_CIPHER_SUITE_LIST;

    case SAI_MACSEC_ATTR_PN_32BIT_SUPPORTED:
      return lemming::dataplane::sai::MACSEC_ATTR_PN_32BIT_SUPPORTED;

    case SAI_MACSEC_ATTR_XPN_64BIT_SUPPORTED:
      return lemming::dataplane::sai::MACSEC_ATTR_XPN_64BIT_SUPPORTED;

    case SAI_MACSEC_ATTR_GCM_AES128_SUPPORTED:
      return lemming::dataplane::sai::MACSEC_ATTR_GCM_AES128_SUPPORTED;

    case SAI_MACSEC_ATTR_GCM_AES256_SUPPORTED:
      return lemming::dataplane::sai::MACSEC_ATTR_GCM_AES256_SUPPORTED;

    case SAI_MACSEC_ATTR_SECTAG_OFFSETS_SUPPORTED:
      return lemming::dataplane::sai::MACSEC_ATTR_SECTAG_OFFSETS_SUPPORTED;

    case SAI_MACSEC_ATTR_SYSTEM_SIDE_MTU:
      return lemming::dataplane::sai::MACSEC_ATTR_SYSTEM_SIDE_MTU;

    case SAI_MACSEC_ATTR_WARM_BOOT_SUPPORTED:
      return lemming::dataplane::sai::MACSEC_ATTR_WARM_BOOT_SUPPORTED;

    case SAI_MACSEC_ATTR_WARM_BOOT_ENABLE:
      return lemming::dataplane::sai::MACSEC_ATTR_WARM_BOOT_ENABLE;

    case SAI_MACSEC_ATTR_CTAG_TPID:
      return lemming::dataplane::sai::MACSEC_ATTR_CTAG_TPID;

    case SAI_MACSEC_ATTR_STAG_TPID:
      return lemming::dataplane::sai::MACSEC_ATTR_STAG_TPID;

    case SAI_MACSEC_ATTR_MAX_VLAN_TAGS_PARSED:
      return lemming::dataplane::sai::MACSEC_ATTR_MAX_VLAN_TAGS_PARSED;

    case SAI_MACSEC_ATTR_STATS_MODE:
      return lemming::dataplane::sai::MACSEC_ATTR_STATS_MODE;

    case SAI_MACSEC_ATTR_PHYSICAL_BYPASS_ENABLE:
      return lemming::dataplane::sai::MACSEC_ATTR_PHYSICAL_BYPASS_ENABLE;

    case SAI_MACSEC_ATTR_SUPPORTED_PORT_LIST:
      return lemming::dataplane::sai::MACSEC_ATTR_SUPPORTED_PORT_LIST;

    case SAI_MACSEC_ATTR_AVAILABLE_MACSEC_FLOW:
      return lemming::dataplane::sai::MACSEC_ATTR_AVAILABLE_MACSEC_FLOW;

    case SAI_MACSEC_ATTR_FLOW_LIST:
      return lemming::dataplane::sai::MACSEC_ATTR_FLOW_LIST;

    case SAI_MACSEC_ATTR_AVAILABLE_MACSEC_SC:
      return lemming::dataplane::sai::MACSEC_ATTR_AVAILABLE_MACSEC_SC;

    case SAI_MACSEC_ATTR_AVAILABLE_MACSEC_SA:
      return lemming::dataplane::sai::MACSEC_ATTR_AVAILABLE_MACSEC_SA;

    case SAI_MACSEC_ATTR_MAX_SECURE_ASSOCIATIONS_PER_SC:
      return lemming::dataplane::sai::
          MACSEC_ATTR_MAX_SECURE_ASSOCIATIONS_PER_SC;

    default:
      return lemming::dataplane::sai::MACSEC_ATTR_UNSPECIFIED;
  }
}
sai_macsec_attr_t convert_sai_macsec_attr_t_to_sai(
    lemming::dataplane::sai::MacsecAttr val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_ATTR_DIRECTION:
      return SAI_MACSEC_ATTR_DIRECTION;

    case lemming::dataplane::sai::
        MACSEC_ATTR_SWITCHING_MODE_CUT_THROUGH_SUPPORTED:
      return SAI_MACSEC_ATTR_SWITCHING_MODE_CUT_THROUGH_SUPPORTED;

    case lemming::dataplane::sai::
        MACSEC_ATTR_SWITCHING_MODE_STORE_AND_FORWARD_SUPPORTED:
      return SAI_MACSEC_ATTR_SWITCHING_MODE_STORE_AND_FORWARD_SUPPORTED;

    case lemming::dataplane::sai::MACSEC_ATTR_STATS_MODE_READ_SUPPORTED:
      return SAI_MACSEC_ATTR_STATS_MODE_READ_SUPPORTED;

    case lemming::dataplane::sai::MACSEC_ATTR_STATS_MODE_READ_CLEAR_SUPPORTED:
      return SAI_MACSEC_ATTR_STATS_MODE_READ_CLEAR_SUPPORTED;

    case lemming::dataplane::sai::MACSEC_ATTR_SCI_IN_INGRESS_MACSEC_ACL:
      return SAI_MACSEC_ATTR_SCI_IN_INGRESS_MACSEC_ACL;

    case lemming::dataplane::sai::MACSEC_ATTR_SUPPORTED_CIPHER_SUITE_LIST:
      return SAI_MACSEC_ATTR_SUPPORTED_CIPHER_SUITE_LIST;

    case lemming::dataplane::sai::MACSEC_ATTR_PN_32BIT_SUPPORTED:
      return SAI_MACSEC_ATTR_PN_32BIT_SUPPORTED;

    case lemming::dataplane::sai::MACSEC_ATTR_XPN_64BIT_SUPPORTED:
      return SAI_MACSEC_ATTR_XPN_64BIT_SUPPORTED;

    case lemming::dataplane::sai::MACSEC_ATTR_GCM_AES128_SUPPORTED:
      return SAI_MACSEC_ATTR_GCM_AES128_SUPPORTED;

    case lemming::dataplane::sai::MACSEC_ATTR_GCM_AES256_SUPPORTED:
      return SAI_MACSEC_ATTR_GCM_AES256_SUPPORTED;

    case lemming::dataplane::sai::MACSEC_ATTR_SECTAG_OFFSETS_SUPPORTED:
      return SAI_MACSEC_ATTR_SECTAG_OFFSETS_SUPPORTED;

    case lemming::dataplane::sai::MACSEC_ATTR_SYSTEM_SIDE_MTU:
      return SAI_MACSEC_ATTR_SYSTEM_SIDE_MTU;

    case lemming::dataplane::sai::MACSEC_ATTR_WARM_BOOT_SUPPORTED:
      return SAI_MACSEC_ATTR_WARM_BOOT_SUPPORTED;

    case lemming::dataplane::sai::MACSEC_ATTR_WARM_BOOT_ENABLE:
      return SAI_MACSEC_ATTR_WARM_BOOT_ENABLE;

    case lemming::dataplane::sai::MACSEC_ATTR_CTAG_TPID:
      return SAI_MACSEC_ATTR_CTAG_TPID;

    case lemming::dataplane::sai::MACSEC_ATTR_STAG_TPID:
      return SAI_MACSEC_ATTR_STAG_TPID;

    case lemming::dataplane::sai::MACSEC_ATTR_MAX_VLAN_TAGS_PARSED:
      return SAI_MACSEC_ATTR_MAX_VLAN_TAGS_PARSED;

    case lemming::dataplane::sai::MACSEC_ATTR_STATS_MODE:
      return SAI_MACSEC_ATTR_STATS_MODE;

    case lemming::dataplane::sai::MACSEC_ATTR_PHYSICAL_BYPASS_ENABLE:
      return SAI_MACSEC_ATTR_PHYSICAL_BYPASS_ENABLE;

    case lemming::dataplane::sai::MACSEC_ATTR_SUPPORTED_PORT_LIST:
      return SAI_MACSEC_ATTR_SUPPORTED_PORT_LIST;

    case lemming::dataplane::sai::MACSEC_ATTR_AVAILABLE_MACSEC_FLOW:
      return SAI_MACSEC_ATTR_AVAILABLE_MACSEC_FLOW;

    case lemming::dataplane::sai::MACSEC_ATTR_FLOW_LIST:
      return SAI_MACSEC_ATTR_FLOW_LIST;

    case lemming::dataplane::sai::MACSEC_ATTR_AVAILABLE_MACSEC_SC:
      return SAI_MACSEC_ATTR_AVAILABLE_MACSEC_SC;

    case lemming::dataplane::sai::MACSEC_ATTR_AVAILABLE_MACSEC_SA:
      return SAI_MACSEC_ATTR_AVAILABLE_MACSEC_SA;

    case lemming::dataplane::sai::MACSEC_ATTR_MAX_SECURE_ASSOCIATIONS_PER_SC:
      return SAI_MACSEC_ATTR_MAX_SECURE_ASSOCIATIONS_PER_SC;

    default:
      return SAI_MACSEC_ATTR_DIRECTION;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_macsec_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecCipherSuite
convert_sai_macsec_cipher_suite_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_CIPHER_SUITE_GCM_AES_128:
      return lemming::dataplane::sai::MACSEC_CIPHER_SUITE_GCM_AES_128;

    case SAI_MACSEC_CIPHER_SUITE_GCM_AES_256:
      return lemming::dataplane::sai::MACSEC_CIPHER_SUITE_GCM_AES_256;

    case SAI_MACSEC_CIPHER_SUITE_GCM_AES_XPN_128:
      return lemming::dataplane::sai::MACSEC_CIPHER_SUITE_GCM_AES_XPN_128;

    case SAI_MACSEC_CIPHER_SUITE_GCM_AES_XPN_256:
      return lemming::dataplane::sai::MACSEC_CIPHER_SUITE_GCM_AES_XPN_256;

    default:
      return lemming::dataplane::sai::MACSEC_CIPHER_SUITE_UNSPECIFIED;
  }
}
sai_macsec_cipher_suite_t convert_sai_macsec_cipher_suite_t_to_sai(
    lemming::dataplane::sai::MacsecCipherSuite val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_CIPHER_SUITE_GCM_AES_128:
      return SAI_MACSEC_CIPHER_SUITE_GCM_AES_128;

    case lemming::dataplane::sai::MACSEC_CIPHER_SUITE_GCM_AES_256:
      return SAI_MACSEC_CIPHER_SUITE_GCM_AES_256;

    case lemming::dataplane::sai::MACSEC_CIPHER_SUITE_GCM_AES_XPN_128:
      return SAI_MACSEC_CIPHER_SUITE_GCM_AES_XPN_128;

    case lemming::dataplane::sai::MACSEC_CIPHER_SUITE_GCM_AES_XPN_256:
      return SAI_MACSEC_CIPHER_SUITE_GCM_AES_XPN_256;

    default:
      return SAI_MACSEC_CIPHER_SUITE_GCM_AES_128;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_macsec_cipher_suite_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_cipher_suite_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_cipher_suite_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_cipher_suite_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecCipherSuite>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecDirection
convert_sai_macsec_direction_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_DIRECTION_EGRESS:
      return lemming::dataplane::sai::MACSEC_DIRECTION_EGRESS;

    case SAI_MACSEC_DIRECTION_INGRESS:
      return lemming::dataplane::sai::MACSEC_DIRECTION_INGRESS;

    default:
      return lemming::dataplane::sai::MACSEC_DIRECTION_UNSPECIFIED;
  }
}
sai_macsec_direction_t convert_sai_macsec_direction_t_to_sai(
    lemming::dataplane::sai::MacsecDirection val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_DIRECTION_EGRESS:
      return SAI_MACSEC_DIRECTION_EGRESS;

    case lemming::dataplane::sai::MACSEC_DIRECTION_INGRESS:
      return SAI_MACSEC_DIRECTION_INGRESS;

    default:
      return SAI_MACSEC_DIRECTION_EGRESS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_macsec_direction_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_direction_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_direction_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_direction_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecDirection>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecFlowAttr convert_sai_macsec_flow_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_FLOW_ATTR_MACSEC_DIRECTION:
      return lemming::dataplane::sai::MACSEC_FLOW_ATTR_MACSEC_DIRECTION;

    case SAI_MACSEC_FLOW_ATTR_ACL_ENTRY_LIST:
      return lemming::dataplane::sai::MACSEC_FLOW_ATTR_ACL_ENTRY_LIST;

    case SAI_MACSEC_FLOW_ATTR_SC_LIST:
      return lemming::dataplane::sai::MACSEC_FLOW_ATTR_SC_LIST;

    default:
      return lemming::dataplane::sai::MACSEC_FLOW_ATTR_UNSPECIFIED;
  }
}
sai_macsec_flow_attr_t convert_sai_macsec_flow_attr_t_to_sai(
    lemming::dataplane::sai::MacsecFlowAttr val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_FLOW_ATTR_MACSEC_DIRECTION:
      return SAI_MACSEC_FLOW_ATTR_MACSEC_DIRECTION;

    case lemming::dataplane::sai::MACSEC_FLOW_ATTR_ACL_ENTRY_LIST:
      return SAI_MACSEC_FLOW_ATTR_ACL_ENTRY_LIST;

    case lemming::dataplane::sai::MACSEC_FLOW_ATTR_SC_LIST:
      return SAI_MACSEC_FLOW_ATTR_SC_LIST;

    default:
      return SAI_MACSEC_FLOW_ATTR_MACSEC_DIRECTION;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_macsec_flow_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_flow_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_flow_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_flow_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecFlowAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecFlowStat convert_sai_macsec_flow_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_FLOW_STAT_OTHER_ERR:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_OTHER_ERR;

    case SAI_MACSEC_FLOW_STAT_OCTETS_UNCONTROLLED:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_OCTETS_UNCONTROLLED;

    case SAI_MACSEC_FLOW_STAT_OCTETS_CONTROLLED:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_OCTETS_CONTROLLED;

    case SAI_MACSEC_FLOW_STAT_OUT_OCTETS_COMMON:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_OUT_OCTETS_COMMON;

    case SAI_MACSEC_FLOW_STAT_UCAST_PKTS_UNCONTROLLED:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_UCAST_PKTS_UNCONTROLLED;

    case SAI_MACSEC_FLOW_STAT_UCAST_PKTS_CONTROLLED:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_UCAST_PKTS_CONTROLLED;

    case SAI_MACSEC_FLOW_STAT_MULTICAST_PKTS_UNCONTROLLED:
      return lemming::dataplane::sai::
          MACSEC_FLOW_STAT_MULTICAST_PKTS_UNCONTROLLED;

    case SAI_MACSEC_FLOW_STAT_MULTICAST_PKTS_CONTROLLED:
      return lemming::dataplane::sai::
          MACSEC_FLOW_STAT_MULTICAST_PKTS_CONTROLLED;

    case SAI_MACSEC_FLOW_STAT_BROADCAST_PKTS_UNCONTROLLED:
      return lemming::dataplane::sai::
          MACSEC_FLOW_STAT_BROADCAST_PKTS_UNCONTROLLED;

    case SAI_MACSEC_FLOW_STAT_BROADCAST_PKTS_CONTROLLED:
      return lemming::dataplane::sai::
          MACSEC_FLOW_STAT_BROADCAST_PKTS_CONTROLLED;

    case SAI_MACSEC_FLOW_STAT_CONTROL_PKTS:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_CONTROL_PKTS;

    case SAI_MACSEC_FLOW_STAT_PKTS_UNTAGGED:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_PKTS_UNTAGGED;

    case SAI_MACSEC_FLOW_STAT_IN_TAGGED_CONTROL_PKTS:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_TAGGED_CONTROL_PKTS;

    case SAI_MACSEC_FLOW_STAT_OUT_PKTS_TOO_LONG:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_OUT_PKTS_TOO_LONG;

    case SAI_MACSEC_FLOW_STAT_IN_PKTS_NO_TAG:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_PKTS_NO_TAG;

    case SAI_MACSEC_FLOW_STAT_IN_PKTS_BAD_TAG:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_PKTS_BAD_TAG;

    case SAI_MACSEC_FLOW_STAT_IN_PKTS_NO_SCI:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_PKTS_NO_SCI;

    case SAI_MACSEC_FLOW_STAT_IN_PKTS_UNKNOWN_SCI:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_PKTS_UNKNOWN_SCI;

    case SAI_MACSEC_FLOW_STAT_IN_PKTS_OVERRUN:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_PKTS_OVERRUN;

    default:
      return lemming::dataplane::sai::MACSEC_FLOW_STAT_UNSPECIFIED;
  }
}
sai_macsec_flow_stat_t convert_sai_macsec_flow_stat_t_to_sai(
    lemming::dataplane::sai::MacsecFlowStat val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_FLOW_STAT_OTHER_ERR:
      return SAI_MACSEC_FLOW_STAT_OTHER_ERR;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_OCTETS_UNCONTROLLED:
      return SAI_MACSEC_FLOW_STAT_OCTETS_UNCONTROLLED;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_OCTETS_CONTROLLED:
      return SAI_MACSEC_FLOW_STAT_OCTETS_CONTROLLED;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_OUT_OCTETS_COMMON:
      return SAI_MACSEC_FLOW_STAT_OUT_OCTETS_COMMON;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_UCAST_PKTS_UNCONTROLLED:
      return SAI_MACSEC_FLOW_STAT_UCAST_PKTS_UNCONTROLLED;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_UCAST_PKTS_CONTROLLED:
      return SAI_MACSEC_FLOW_STAT_UCAST_PKTS_CONTROLLED;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_MULTICAST_PKTS_UNCONTROLLED:
      return SAI_MACSEC_FLOW_STAT_MULTICAST_PKTS_UNCONTROLLED;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_MULTICAST_PKTS_CONTROLLED:
      return SAI_MACSEC_FLOW_STAT_MULTICAST_PKTS_CONTROLLED;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_BROADCAST_PKTS_UNCONTROLLED:
      return SAI_MACSEC_FLOW_STAT_BROADCAST_PKTS_UNCONTROLLED;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_BROADCAST_PKTS_CONTROLLED:
      return SAI_MACSEC_FLOW_STAT_BROADCAST_PKTS_CONTROLLED;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_CONTROL_PKTS:
      return SAI_MACSEC_FLOW_STAT_CONTROL_PKTS;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_PKTS_UNTAGGED:
      return SAI_MACSEC_FLOW_STAT_PKTS_UNTAGGED;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_TAGGED_CONTROL_PKTS:
      return SAI_MACSEC_FLOW_STAT_IN_TAGGED_CONTROL_PKTS;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_OUT_PKTS_TOO_LONG:
      return SAI_MACSEC_FLOW_STAT_OUT_PKTS_TOO_LONG;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_PKTS_NO_TAG:
      return SAI_MACSEC_FLOW_STAT_IN_PKTS_NO_TAG;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_PKTS_BAD_TAG:
      return SAI_MACSEC_FLOW_STAT_IN_PKTS_BAD_TAG;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_PKTS_NO_SCI:
      return SAI_MACSEC_FLOW_STAT_IN_PKTS_NO_SCI;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_PKTS_UNKNOWN_SCI:
      return SAI_MACSEC_FLOW_STAT_IN_PKTS_UNKNOWN_SCI;

    case lemming::dataplane::sai::MACSEC_FLOW_STAT_IN_PKTS_OVERRUN:
      return SAI_MACSEC_FLOW_STAT_IN_PKTS_OVERRUN;

    default:
      return SAI_MACSEC_FLOW_STAT_OTHER_ERR;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_macsec_flow_stat_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_flow_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_flow_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_flow_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecFlowStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecMaxSecureAssociationsPerSc
convert_sai_macsec_max_secure_associations_per_sc_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_MAX_SECURE_ASSOCIATIONS_PER_SC_TWO:
      return lemming::dataplane::sai::MACSEC_MAX_SECURE_ASSOCIATIONS_PER_SC_TWO;

    case SAI_MACSEC_MAX_SECURE_ASSOCIATIONS_PER_SC_FOUR:
      return lemming::dataplane::sai::
          MACSEC_MAX_SECURE_ASSOCIATIONS_PER_SC_FOUR;

    default:
      return lemming::dataplane::sai::
          MACSEC_MAX_SECURE_ASSOCIATIONS_PER_SC_UNSPECIFIED;
  }
}
sai_macsec_max_secure_associations_per_sc_t
convert_sai_macsec_max_secure_associations_per_sc_t_to_sai(
    lemming::dataplane::sai::MacsecMaxSecureAssociationsPerSc val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_MAX_SECURE_ASSOCIATIONS_PER_SC_TWO:
      return SAI_MACSEC_MAX_SECURE_ASSOCIATIONS_PER_SC_TWO;

    case lemming::dataplane::sai::MACSEC_MAX_SECURE_ASSOCIATIONS_PER_SC_FOUR:
      return SAI_MACSEC_MAX_SECURE_ASSOCIATIONS_PER_SC_FOUR;

    default:
      return SAI_MACSEC_MAX_SECURE_ASSOCIATIONS_PER_SC_TWO;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_macsec_max_secure_associations_per_sc_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_max_secure_associations_per_sc_t_to_proto(
        list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_max_secure_associations_per_sc_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_max_secure_associations_per_sc_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecMaxSecureAssociationsPerSc>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecPortAttr convert_sai_macsec_port_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_PORT_ATTR_MACSEC_DIRECTION:
      return lemming::dataplane::sai::MACSEC_PORT_ATTR_MACSEC_DIRECTION;

    case SAI_MACSEC_PORT_ATTR_PORT_ID:
      return lemming::dataplane::sai::MACSEC_PORT_ATTR_PORT_ID;

    case SAI_MACSEC_PORT_ATTR_CTAG_ENABLE:
      return lemming::dataplane::sai::MACSEC_PORT_ATTR_CTAG_ENABLE;

    case SAI_MACSEC_PORT_ATTR_STAG_ENABLE:
      return lemming::dataplane::sai::MACSEC_PORT_ATTR_STAG_ENABLE;

    case SAI_MACSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
      return lemming::dataplane::sai::MACSEC_PORT_ATTR_SWITCH_SWITCHING_MODE;

    default:
      return lemming::dataplane::sai::MACSEC_PORT_ATTR_UNSPECIFIED;
  }
}
sai_macsec_port_attr_t convert_sai_macsec_port_attr_t_to_sai(
    lemming::dataplane::sai::MacsecPortAttr val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_PORT_ATTR_MACSEC_DIRECTION:
      return SAI_MACSEC_PORT_ATTR_MACSEC_DIRECTION;

    case lemming::dataplane::sai::MACSEC_PORT_ATTR_PORT_ID:
      return SAI_MACSEC_PORT_ATTR_PORT_ID;

    case lemming::dataplane::sai::MACSEC_PORT_ATTR_CTAG_ENABLE:
      return SAI_MACSEC_PORT_ATTR_CTAG_ENABLE;

    case lemming::dataplane::sai::MACSEC_PORT_ATTR_STAG_ENABLE:
      return SAI_MACSEC_PORT_ATTR_STAG_ENABLE;

    case lemming::dataplane::sai::MACSEC_PORT_ATTR_SWITCH_SWITCHING_MODE:
      return SAI_MACSEC_PORT_ATTR_SWITCH_SWITCHING_MODE;

    default:
      return SAI_MACSEC_PORT_ATTR_MACSEC_DIRECTION;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_macsec_port_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_port_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_port_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_port_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecPortAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecPortStat convert_sai_macsec_port_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_PORT_STAT_PRE_MACSEC_DROP_PKTS:
      return lemming::dataplane::sai::MACSEC_PORT_STAT_PRE_MACSEC_DROP_PKTS;

    case SAI_MACSEC_PORT_STAT_CONTROL_PKTS:
      return lemming::dataplane::sai::MACSEC_PORT_STAT_CONTROL_PKTS;

    case SAI_MACSEC_PORT_STAT_DATA_PKTS:
      return lemming::dataplane::sai::MACSEC_PORT_STAT_DATA_PKTS;

    default:
      return lemming::dataplane::sai::MACSEC_PORT_STAT_UNSPECIFIED;
  }
}
sai_macsec_port_stat_t convert_sai_macsec_port_stat_t_to_sai(
    lemming::dataplane::sai::MacsecPortStat val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_PORT_STAT_PRE_MACSEC_DROP_PKTS:
      return SAI_MACSEC_PORT_STAT_PRE_MACSEC_DROP_PKTS;

    case lemming::dataplane::sai::MACSEC_PORT_STAT_CONTROL_PKTS:
      return SAI_MACSEC_PORT_STAT_CONTROL_PKTS;

    case lemming::dataplane::sai::MACSEC_PORT_STAT_DATA_PKTS:
      return SAI_MACSEC_PORT_STAT_DATA_PKTS;

    default:
      return SAI_MACSEC_PORT_STAT_PRE_MACSEC_DROP_PKTS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_macsec_port_stat_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_port_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_port_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_port_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecPortStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecSaAttr convert_sai_macsec_sa_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_SA_ATTR_MACSEC_DIRECTION:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_MACSEC_DIRECTION;

    case SAI_MACSEC_SA_ATTR_SC_ID:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_SC_ID;

    case SAI_MACSEC_SA_ATTR_AN:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_AN;

    case SAI_MACSEC_SA_ATTR_SAK:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_SAK;

    case SAI_MACSEC_SA_ATTR_SALT:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_SALT;

    case SAI_MACSEC_SA_ATTR_AUTH_KEY:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_AUTH_KEY;

    case SAI_MACSEC_SA_ATTR_CONFIGURED_EGRESS_XPN:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_CONFIGURED_EGRESS_XPN;

    case SAI_MACSEC_SA_ATTR_CURRENT_XPN:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_CURRENT_XPN;

    case SAI_MACSEC_SA_ATTR_MINIMUM_INGRESS_XPN:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_MINIMUM_INGRESS_XPN;

    case SAI_MACSEC_SA_ATTR_MACSEC_SSCI:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_MACSEC_SSCI;

    default:
      return lemming::dataplane::sai::MACSEC_SA_ATTR_UNSPECIFIED;
  }
}
sai_macsec_sa_attr_t convert_sai_macsec_sa_attr_t_to_sai(
    lemming::dataplane::sai::MacsecSaAttr val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_SA_ATTR_MACSEC_DIRECTION:
      return SAI_MACSEC_SA_ATTR_MACSEC_DIRECTION;

    case lemming::dataplane::sai::MACSEC_SA_ATTR_SC_ID:
      return SAI_MACSEC_SA_ATTR_SC_ID;

    case lemming::dataplane::sai::MACSEC_SA_ATTR_AN:
      return SAI_MACSEC_SA_ATTR_AN;

    case lemming::dataplane::sai::MACSEC_SA_ATTR_SAK:
      return SAI_MACSEC_SA_ATTR_SAK;

    case lemming::dataplane::sai::MACSEC_SA_ATTR_SALT:
      return SAI_MACSEC_SA_ATTR_SALT;

    case lemming::dataplane::sai::MACSEC_SA_ATTR_AUTH_KEY:
      return SAI_MACSEC_SA_ATTR_AUTH_KEY;

    case lemming::dataplane::sai::MACSEC_SA_ATTR_CONFIGURED_EGRESS_XPN:
      return SAI_MACSEC_SA_ATTR_CONFIGURED_EGRESS_XPN;

    case lemming::dataplane::sai::MACSEC_SA_ATTR_CURRENT_XPN:
      return SAI_MACSEC_SA_ATTR_CURRENT_XPN;

    case lemming::dataplane::sai::MACSEC_SA_ATTR_MINIMUM_INGRESS_XPN:
      return SAI_MACSEC_SA_ATTR_MINIMUM_INGRESS_XPN;

    case lemming::dataplane::sai::MACSEC_SA_ATTR_MACSEC_SSCI:
      return SAI_MACSEC_SA_ATTR_MACSEC_SSCI;

    default:
      return SAI_MACSEC_SA_ATTR_MACSEC_DIRECTION;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_macsec_sa_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_sa_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_sa_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_sa_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecSaAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecSaStat convert_sai_macsec_sa_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_SA_STAT_OCTETS_ENCRYPTED:
      return lemming::dataplane::sai::MACSEC_SA_STAT_OCTETS_ENCRYPTED;

    case SAI_MACSEC_SA_STAT_OCTETS_PROTECTED:
      return lemming::dataplane::sai::MACSEC_SA_STAT_OCTETS_PROTECTED;

    case SAI_MACSEC_SA_STAT_OUT_PKTS_ENCRYPTED:
      return lemming::dataplane::sai::MACSEC_SA_STAT_OUT_PKTS_ENCRYPTED;

    case SAI_MACSEC_SA_STAT_OUT_PKTS_PROTECTED:
      return lemming::dataplane::sai::MACSEC_SA_STAT_OUT_PKTS_PROTECTED;

    case SAI_MACSEC_SA_STAT_IN_PKTS_UNCHECKED:
      return lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_UNCHECKED;

    case SAI_MACSEC_SA_STAT_IN_PKTS_DELAYED:
      return lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_DELAYED;

    case SAI_MACSEC_SA_STAT_IN_PKTS_LATE:
      return lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_LATE;

    case SAI_MACSEC_SA_STAT_IN_PKTS_INVALID:
      return lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_INVALID;

    case SAI_MACSEC_SA_STAT_IN_PKTS_NOT_VALID:
      return lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_NOT_VALID;

    case SAI_MACSEC_SA_STAT_IN_PKTS_NOT_USING_SA:
      return lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_NOT_USING_SA;

    case SAI_MACSEC_SA_STAT_IN_PKTS_UNUSED_SA:
      return lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_UNUSED_SA;

    case SAI_MACSEC_SA_STAT_IN_PKTS_OK:
      return lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_OK;

    default:
      return lemming::dataplane::sai::MACSEC_SA_STAT_UNSPECIFIED;
  }
}
sai_macsec_sa_stat_t convert_sai_macsec_sa_stat_t_to_sai(
    lemming::dataplane::sai::MacsecSaStat val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_SA_STAT_OCTETS_ENCRYPTED:
      return SAI_MACSEC_SA_STAT_OCTETS_ENCRYPTED;

    case lemming::dataplane::sai::MACSEC_SA_STAT_OCTETS_PROTECTED:
      return SAI_MACSEC_SA_STAT_OCTETS_PROTECTED;

    case lemming::dataplane::sai::MACSEC_SA_STAT_OUT_PKTS_ENCRYPTED:
      return SAI_MACSEC_SA_STAT_OUT_PKTS_ENCRYPTED;

    case lemming::dataplane::sai::MACSEC_SA_STAT_OUT_PKTS_PROTECTED:
      return SAI_MACSEC_SA_STAT_OUT_PKTS_PROTECTED;

    case lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_UNCHECKED:
      return SAI_MACSEC_SA_STAT_IN_PKTS_UNCHECKED;

    case lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_DELAYED:
      return SAI_MACSEC_SA_STAT_IN_PKTS_DELAYED;

    case lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_LATE:
      return SAI_MACSEC_SA_STAT_IN_PKTS_LATE;

    case lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_INVALID:
      return SAI_MACSEC_SA_STAT_IN_PKTS_INVALID;

    case lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_NOT_VALID:
      return SAI_MACSEC_SA_STAT_IN_PKTS_NOT_VALID;

    case lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_NOT_USING_SA:
      return SAI_MACSEC_SA_STAT_IN_PKTS_NOT_USING_SA;

    case lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_UNUSED_SA:
      return SAI_MACSEC_SA_STAT_IN_PKTS_UNUSED_SA;

    case lemming::dataplane::sai::MACSEC_SA_STAT_IN_PKTS_OK:
      return SAI_MACSEC_SA_STAT_IN_PKTS_OK;

    default:
      return SAI_MACSEC_SA_STAT_OCTETS_ENCRYPTED;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_macsec_sa_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_sa_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_sa_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_sa_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecSaStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecScAttr convert_sai_macsec_sc_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_SC_ATTR_MACSEC_DIRECTION:
      return lemming::dataplane::sai::MACSEC_SC_ATTR_MACSEC_DIRECTION;

    case SAI_MACSEC_SC_ATTR_FLOW_ID:
      return lemming::dataplane::sai::MACSEC_SC_ATTR_FLOW_ID;

    case SAI_MACSEC_SC_ATTR_MACSEC_SCI:
      return lemming::dataplane::sai::MACSEC_SC_ATTR_MACSEC_SCI;

    case SAI_MACSEC_SC_ATTR_MACSEC_EXPLICIT_SCI_ENABLE:
      return lemming::dataplane::sai::MACSEC_SC_ATTR_MACSEC_EXPLICIT_SCI_ENABLE;

    case SAI_MACSEC_SC_ATTR_MACSEC_SECTAG_OFFSET:
      return lemming::dataplane::sai::MACSEC_SC_ATTR_MACSEC_SECTAG_OFFSET;

    case SAI_MACSEC_SC_ATTR_ACTIVE_EGRESS_SA_ID:
      return lemming::dataplane::sai::MACSEC_SC_ATTR_ACTIVE_EGRESS_SA_ID;

    case SAI_MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_ENABLE:
      return lemming::dataplane::sai::
          MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_ENABLE;

    case SAI_MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_WINDOW:
      return lemming::dataplane::sai::
          MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_WINDOW;

    case SAI_MACSEC_SC_ATTR_SA_LIST:
      return lemming::dataplane::sai::MACSEC_SC_ATTR_SA_LIST;

    case SAI_MACSEC_SC_ATTR_MACSEC_CIPHER_SUITE:
      return lemming::dataplane::sai::MACSEC_SC_ATTR_MACSEC_CIPHER_SUITE;

    case SAI_MACSEC_SC_ATTR_ENCRYPTION_ENABLE:
      return lemming::dataplane::sai::MACSEC_SC_ATTR_ENCRYPTION_ENABLE;

    default:
      return lemming::dataplane::sai::MACSEC_SC_ATTR_UNSPECIFIED;
  }
}
sai_macsec_sc_attr_t convert_sai_macsec_sc_attr_t_to_sai(
    lemming::dataplane::sai::MacsecScAttr val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_SC_ATTR_MACSEC_DIRECTION:
      return SAI_MACSEC_SC_ATTR_MACSEC_DIRECTION;

    case lemming::dataplane::sai::MACSEC_SC_ATTR_FLOW_ID:
      return SAI_MACSEC_SC_ATTR_FLOW_ID;

    case lemming::dataplane::sai::MACSEC_SC_ATTR_MACSEC_SCI:
      return SAI_MACSEC_SC_ATTR_MACSEC_SCI;

    case lemming::dataplane::sai::MACSEC_SC_ATTR_MACSEC_EXPLICIT_SCI_ENABLE:
      return SAI_MACSEC_SC_ATTR_MACSEC_EXPLICIT_SCI_ENABLE;

    case lemming::dataplane::sai::MACSEC_SC_ATTR_MACSEC_SECTAG_OFFSET:
      return SAI_MACSEC_SC_ATTR_MACSEC_SECTAG_OFFSET;

    case lemming::dataplane::sai::MACSEC_SC_ATTR_ACTIVE_EGRESS_SA_ID:
      return SAI_MACSEC_SC_ATTR_ACTIVE_EGRESS_SA_ID;

    case lemming::dataplane::sai::
        MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_ENABLE:
      return SAI_MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_ENABLE;

    case lemming::dataplane::sai::
        MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_WINDOW:
      return SAI_MACSEC_SC_ATTR_MACSEC_REPLAY_PROTECTION_WINDOW;

    case lemming::dataplane::sai::MACSEC_SC_ATTR_SA_LIST:
      return SAI_MACSEC_SC_ATTR_SA_LIST;

    case lemming::dataplane::sai::MACSEC_SC_ATTR_MACSEC_CIPHER_SUITE:
      return SAI_MACSEC_SC_ATTR_MACSEC_CIPHER_SUITE;

    case lemming::dataplane::sai::MACSEC_SC_ATTR_ENCRYPTION_ENABLE:
      return SAI_MACSEC_SC_ATTR_ENCRYPTION_ENABLE;

    default:
      return SAI_MACSEC_SC_ATTR_MACSEC_DIRECTION;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_macsec_sc_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_sc_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_sc_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_sc_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecScAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MacsecScStat convert_sai_macsec_sc_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MACSEC_SC_STAT_SA_NOT_IN_USE:
      return lemming::dataplane::sai::MACSEC_SC_STAT_SA_NOT_IN_USE;

    default:
      return lemming::dataplane::sai::MACSEC_SC_STAT_UNSPECIFIED;
  }
}
sai_macsec_sc_stat_t convert_sai_macsec_sc_stat_t_to_sai(
    lemming::dataplane::sai::MacsecScStat val) {
  switch (val) {
    case lemming::dataplane::sai::MACSEC_SC_STAT_SA_NOT_IN_USE:
      return SAI_MACSEC_SC_STAT_SA_NOT_IN_USE;

    default:
      return SAI_MACSEC_SC_STAT_SA_NOT_IN_USE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_macsec_sc_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_macsec_sc_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_macsec_sc_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_macsec_sc_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::MacsecScStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::McastFdbEntryAttr
convert_sai_mcast_fdb_entry_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_MCAST_FDB_ENTRY_ATTR_GROUP_ID:
      return lemming::dataplane::sai::MCAST_FDB_ENTRY_ATTR_GROUP_ID;

    case SAI_MCAST_FDB_ENTRY_ATTR_PACKET_ACTION:
      return lemming::dataplane::sai::MCAST_FDB_ENTRY_ATTR_PACKET_ACTION;

    case SAI_MCAST_FDB_ENTRY_ATTR_META_DATA:
      return lemming::dataplane::sai::MCAST_FDB_ENTRY_ATTR_META_DATA;

    default:
      return lemming::dataplane::sai::MCAST_FDB_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_mcast_fdb_entry_attr_t convert_sai_mcast_fdb_entry_attr_t_to_sai(
    lemming::dataplane::sai::McastFdbEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::MCAST_FDB_ENTRY_ATTR_GROUP_ID:
      return SAI_MCAST_FDB_ENTRY_ATTR_GROUP_ID;

    case lemming::dataplane::sai::MCAST_FDB_ENTRY_ATTR_PACKET_ACTION:
      return SAI_MCAST_FDB_ENTRY_ATTR_PACKET_ACTION;

    case lemming::dataplane::sai::MCAST_FDB_ENTRY_ATTR_META_DATA:
      return SAI_MCAST_FDB_ENTRY_ATTR_META_DATA;

    default:
      return SAI_MCAST_FDB_ENTRY_ATTR_GROUP_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_mcast_fdb_entry_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_mcast_fdb_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_mcast_fdb_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_mcast_fdb_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::McastFdbEntryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MeterType convert_sai_meter_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_METER_TYPE_PACKETS:
      return lemming::dataplane::sai::METER_TYPE_PACKETS;

    case SAI_METER_TYPE_BYTES:
      return lemming::dataplane::sai::METER_TYPE_BYTES;

    case SAI_METER_TYPE_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::METER_TYPE_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::METER_TYPE_UNSPECIFIED;
  }
}
sai_meter_type_t convert_sai_meter_type_t_to_sai(
    lemming::dataplane::sai::MeterType val) {
  switch (val) {
    case lemming::dataplane::sai::METER_TYPE_PACKETS:
      return SAI_METER_TYPE_PACKETS;

    case lemming::dataplane::sai::METER_TYPE_BYTES:
      return SAI_METER_TYPE_BYTES;

    case lemming::dataplane::sai::METER_TYPE_CUSTOM_RANGE_BASE:
      return SAI_METER_TYPE_CUSTOM_RANGE_BASE;

    default:
      return SAI_METER_TYPE_PACKETS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_meter_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_meter_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_meter_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_meter_type_t_to_sai(
        static_cast<lemming::dataplane::sai::MeterType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MirrorSessionAttr
convert_sai_mirror_session_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_MIRROR_SESSION_ATTR_TYPE:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_TYPE;

    case SAI_MIRROR_SESSION_ATTR_MONITOR_PORT:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_MONITOR_PORT;

    case SAI_MIRROR_SESSION_ATTR_TRUNCATE_SIZE:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_TRUNCATE_SIZE;

    case SAI_MIRROR_SESSION_ATTR_SAMPLE_RATE:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_SAMPLE_RATE;

    case SAI_MIRROR_SESSION_ATTR_CONGESTION_MODE:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_CONGESTION_MODE;

    case SAI_MIRROR_SESSION_ATTR_TC:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_TC;

    case SAI_MIRROR_SESSION_ATTR_VLAN_TPID:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_VLAN_TPID;

    case SAI_MIRROR_SESSION_ATTR_VLAN_ID:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_VLAN_ID;

    case SAI_MIRROR_SESSION_ATTR_VLAN_PRI:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_VLAN_PRI;

    case SAI_MIRROR_SESSION_ATTR_VLAN_CFI:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_VLAN_CFI;

    case SAI_MIRROR_SESSION_ATTR_VLAN_HEADER_VALID:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_VLAN_HEADER_VALID;

    case SAI_MIRROR_SESSION_ATTR_ERSPAN_ENCAPSULATION_TYPE:
      return lemming::dataplane::sai::
          MIRROR_SESSION_ATTR_ERSPAN_ENCAPSULATION_TYPE;

    case SAI_MIRROR_SESSION_ATTR_IPHDR_VERSION:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_IPHDR_VERSION;

    case SAI_MIRROR_SESSION_ATTR_TOS:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_TOS;

    case SAI_MIRROR_SESSION_ATTR_TTL:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_TTL;

    case SAI_MIRROR_SESSION_ATTR_SRC_IP_ADDRESS:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_SRC_IP_ADDRESS;

    case SAI_MIRROR_SESSION_ATTR_DST_IP_ADDRESS:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_DST_IP_ADDRESS;

    case SAI_MIRROR_SESSION_ATTR_SRC_MAC_ADDRESS:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_SRC_MAC_ADDRESS;

    case SAI_MIRROR_SESSION_ATTR_DST_MAC_ADDRESS:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_DST_MAC_ADDRESS;

    case SAI_MIRROR_SESSION_ATTR_GRE_PROTOCOL_TYPE:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_GRE_PROTOCOL_TYPE;

    case SAI_MIRROR_SESSION_ATTR_MONITOR_PORTLIST_VALID:
      return lemming::dataplane::sai::
          MIRROR_SESSION_ATTR_MONITOR_PORTLIST_VALID;

    case SAI_MIRROR_SESSION_ATTR_MONITOR_PORTLIST:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_MONITOR_PORTLIST;

    case SAI_MIRROR_SESSION_ATTR_POLICER:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_POLICER;

    case SAI_MIRROR_SESSION_ATTR_UDP_SRC_PORT:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_UDP_SRC_PORT;

    case SAI_MIRROR_SESSION_ATTR_UDP_DST_PORT:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_UDP_DST_PORT;

    case SAI_MIRROR_SESSION_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_COUNTER_ID;

    default:
      return lemming::dataplane::sai::MIRROR_SESSION_ATTR_UNSPECIFIED;
  }
}
sai_mirror_session_attr_t convert_sai_mirror_session_attr_t_to_sai(
    lemming::dataplane::sai::MirrorSessionAttr val) {
  switch (val) {
    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_TYPE:
      return SAI_MIRROR_SESSION_ATTR_TYPE;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_MONITOR_PORT:
      return SAI_MIRROR_SESSION_ATTR_MONITOR_PORT;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_TRUNCATE_SIZE:
      return SAI_MIRROR_SESSION_ATTR_TRUNCATE_SIZE;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_SAMPLE_RATE:
      return SAI_MIRROR_SESSION_ATTR_SAMPLE_RATE;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_CONGESTION_MODE:
      return SAI_MIRROR_SESSION_ATTR_CONGESTION_MODE;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_TC:
      return SAI_MIRROR_SESSION_ATTR_TC;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_VLAN_TPID:
      return SAI_MIRROR_SESSION_ATTR_VLAN_TPID;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_VLAN_ID:
      return SAI_MIRROR_SESSION_ATTR_VLAN_ID;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_VLAN_PRI:
      return SAI_MIRROR_SESSION_ATTR_VLAN_PRI;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_VLAN_CFI:
      return SAI_MIRROR_SESSION_ATTR_VLAN_CFI;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_VLAN_HEADER_VALID:
      return SAI_MIRROR_SESSION_ATTR_VLAN_HEADER_VALID;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_ERSPAN_ENCAPSULATION_TYPE:
      return SAI_MIRROR_SESSION_ATTR_ERSPAN_ENCAPSULATION_TYPE;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_IPHDR_VERSION:
      return SAI_MIRROR_SESSION_ATTR_IPHDR_VERSION;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_TOS:
      return SAI_MIRROR_SESSION_ATTR_TOS;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_TTL:
      return SAI_MIRROR_SESSION_ATTR_TTL;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_SRC_IP_ADDRESS:
      return SAI_MIRROR_SESSION_ATTR_SRC_IP_ADDRESS;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_DST_IP_ADDRESS:
      return SAI_MIRROR_SESSION_ATTR_DST_IP_ADDRESS;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_SRC_MAC_ADDRESS:
      return SAI_MIRROR_SESSION_ATTR_SRC_MAC_ADDRESS;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_DST_MAC_ADDRESS:
      return SAI_MIRROR_SESSION_ATTR_DST_MAC_ADDRESS;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_GRE_PROTOCOL_TYPE:
      return SAI_MIRROR_SESSION_ATTR_GRE_PROTOCOL_TYPE;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_MONITOR_PORTLIST_VALID:
      return SAI_MIRROR_SESSION_ATTR_MONITOR_PORTLIST_VALID;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_MONITOR_PORTLIST:
      return SAI_MIRROR_SESSION_ATTR_MONITOR_PORTLIST;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_POLICER:
      return SAI_MIRROR_SESSION_ATTR_POLICER;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_UDP_SRC_PORT:
      return SAI_MIRROR_SESSION_ATTR_UDP_SRC_PORT;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_UDP_DST_PORT:
      return SAI_MIRROR_SESSION_ATTR_UDP_DST_PORT;

    case lemming::dataplane::sai::MIRROR_SESSION_ATTR_COUNTER_ID:
      return SAI_MIRROR_SESSION_ATTR_COUNTER_ID;

    default:
      return SAI_MIRROR_SESSION_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_mirror_session_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_mirror_session_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_mirror_session_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_mirror_session_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::MirrorSessionAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MirrorSessionCongestionMode
convert_sai_mirror_session_congestion_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_MIRROR_SESSION_CONGESTION_MODE_INDEPENDENT:
      return lemming::dataplane::sai::
          MIRROR_SESSION_CONGESTION_MODE_INDEPENDENT;

    case SAI_MIRROR_SESSION_CONGESTION_MODE_CORRELATED:
      return lemming::dataplane::sai::MIRROR_SESSION_CONGESTION_MODE_CORRELATED;

    default:
      return lemming::dataplane::sai::
          MIRROR_SESSION_CONGESTION_MODE_UNSPECIFIED;
  }
}
sai_mirror_session_congestion_mode_t
convert_sai_mirror_session_congestion_mode_t_to_sai(
    lemming::dataplane::sai::MirrorSessionCongestionMode val) {
  switch (val) {
    case lemming::dataplane::sai::MIRROR_SESSION_CONGESTION_MODE_INDEPENDENT:
      return SAI_MIRROR_SESSION_CONGESTION_MODE_INDEPENDENT;

    case lemming::dataplane::sai::MIRROR_SESSION_CONGESTION_MODE_CORRELATED:
      return SAI_MIRROR_SESSION_CONGESTION_MODE_CORRELATED;

    default:
      return SAI_MIRROR_SESSION_CONGESTION_MODE_INDEPENDENT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_mirror_session_congestion_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_mirror_session_congestion_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_mirror_session_congestion_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_mirror_session_congestion_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::MirrorSessionCongestionMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MirrorSessionType
convert_sai_mirror_session_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_MIRROR_SESSION_TYPE_LOCAL:
      return lemming::dataplane::sai::MIRROR_SESSION_TYPE_LOCAL;

    case SAI_MIRROR_SESSION_TYPE_REMOTE:
      return lemming::dataplane::sai::MIRROR_SESSION_TYPE_REMOTE;

    case SAI_MIRROR_SESSION_TYPE_ENHANCED_REMOTE:
      return lemming::dataplane::sai::MIRROR_SESSION_TYPE_ENHANCED_REMOTE;

    case SAI_MIRROR_SESSION_TYPE_SFLOW:
      return lemming::dataplane::sai::MIRROR_SESSION_TYPE_SFLOW;

    default:
      return lemming::dataplane::sai::MIRROR_SESSION_TYPE_UNSPECIFIED;
  }
}
sai_mirror_session_type_t convert_sai_mirror_session_type_t_to_sai(
    lemming::dataplane::sai::MirrorSessionType val) {
  switch (val) {
    case lemming::dataplane::sai::MIRROR_SESSION_TYPE_LOCAL:
      return SAI_MIRROR_SESSION_TYPE_LOCAL;

    case lemming::dataplane::sai::MIRROR_SESSION_TYPE_REMOTE:
      return SAI_MIRROR_SESSION_TYPE_REMOTE;

    case lemming::dataplane::sai::MIRROR_SESSION_TYPE_ENHANCED_REMOTE:
      return SAI_MIRROR_SESSION_TYPE_ENHANCED_REMOTE;

    case lemming::dataplane::sai::MIRROR_SESSION_TYPE_SFLOW:
      return SAI_MIRROR_SESSION_TYPE_SFLOW;

    default:
      return SAI_MIRROR_SESSION_TYPE_LOCAL;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_mirror_session_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_mirror_session_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_mirror_session_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_mirror_session_type_t_to_sai(
        static_cast<lemming::dataplane::sai::MirrorSessionType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MyMacAttr convert_sai_my_mac_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MY_MAC_ATTR_PRIORITY:
      return lemming::dataplane::sai::MY_MAC_ATTR_PRIORITY;

    case SAI_MY_MAC_ATTR_PORT_ID:
      return lemming::dataplane::sai::MY_MAC_ATTR_PORT_ID;

    case SAI_MY_MAC_ATTR_VLAN_ID:
      return lemming::dataplane::sai::MY_MAC_ATTR_VLAN_ID;

    case SAI_MY_MAC_ATTR_MAC_ADDRESS:
      return lemming::dataplane::sai::MY_MAC_ATTR_MAC_ADDRESS;

    case SAI_MY_MAC_ATTR_MAC_ADDRESS_MASK:
      return lemming::dataplane::sai::MY_MAC_ATTR_MAC_ADDRESS_MASK;

    default:
      return lemming::dataplane::sai::MY_MAC_ATTR_UNSPECIFIED;
  }
}
sai_my_mac_attr_t convert_sai_my_mac_attr_t_to_sai(
    lemming::dataplane::sai::MyMacAttr val) {
  switch (val) {
    case lemming::dataplane::sai::MY_MAC_ATTR_PRIORITY:
      return SAI_MY_MAC_ATTR_PRIORITY;

    case lemming::dataplane::sai::MY_MAC_ATTR_PORT_ID:
      return SAI_MY_MAC_ATTR_PORT_ID;

    case lemming::dataplane::sai::MY_MAC_ATTR_VLAN_ID:
      return SAI_MY_MAC_ATTR_VLAN_ID;

    case lemming::dataplane::sai::MY_MAC_ATTR_MAC_ADDRESS:
      return SAI_MY_MAC_ATTR_MAC_ADDRESS;

    case lemming::dataplane::sai::MY_MAC_ATTR_MAC_ADDRESS_MASK:
      return SAI_MY_MAC_ATTR_MAC_ADDRESS_MASK;

    default:
      return SAI_MY_MAC_ATTR_PRIORITY;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_my_mac_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_my_mac_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_my_mac_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_my_mac_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::MyMacAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MySidEntryAttr
convert_sai_my_sid_entry_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR:
      return lemming::dataplane::sai::MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR;

    case SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR_FLAVOR:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR_FLAVOR;

    case SAI_MY_SID_ENTRY_ATTR_PACKET_ACTION:
      return lemming::dataplane::sai::MY_SID_ENTRY_ATTR_PACKET_ACTION;

    case SAI_MY_SID_ENTRY_ATTR_TRAP_PRIORITY:
      return lemming::dataplane::sai::MY_SID_ENTRY_ATTR_TRAP_PRIORITY;

    case SAI_MY_SID_ENTRY_ATTR_NEXT_HOP_ID:
      return lemming::dataplane::sai::MY_SID_ENTRY_ATTR_NEXT_HOP_ID;

    case SAI_MY_SID_ENTRY_ATTR_TUNNEL_ID:
      return lemming::dataplane::sai::MY_SID_ENTRY_ATTR_TUNNEL_ID;

    case SAI_MY_SID_ENTRY_ATTR_VRF:
      return lemming::dataplane::sai::MY_SID_ENTRY_ATTR_VRF;

    case SAI_MY_SID_ENTRY_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::MY_SID_ENTRY_ATTR_COUNTER_ID;

    default:
      return lemming::dataplane::sai::MY_SID_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_my_sid_entry_attr_t convert_sai_my_sid_entry_attr_t_to_sai(
    lemming::dataplane::sai::MySidEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR:
      return SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR;

    case lemming::dataplane::sai::MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR_FLAVOR:
      return SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR_FLAVOR;

    case lemming::dataplane::sai::MY_SID_ENTRY_ATTR_PACKET_ACTION:
      return SAI_MY_SID_ENTRY_ATTR_PACKET_ACTION;

    case lemming::dataplane::sai::MY_SID_ENTRY_ATTR_TRAP_PRIORITY:
      return SAI_MY_SID_ENTRY_ATTR_TRAP_PRIORITY;

    case lemming::dataplane::sai::MY_SID_ENTRY_ATTR_NEXT_HOP_ID:
      return SAI_MY_SID_ENTRY_ATTR_NEXT_HOP_ID;

    case lemming::dataplane::sai::MY_SID_ENTRY_ATTR_TUNNEL_ID:
      return SAI_MY_SID_ENTRY_ATTR_TUNNEL_ID;

    case lemming::dataplane::sai::MY_SID_ENTRY_ATTR_VRF:
      return SAI_MY_SID_ENTRY_ATTR_VRF;

    case lemming::dataplane::sai::MY_SID_ENTRY_ATTR_COUNTER_ID:
      return SAI_MY_SID_ENTRY_ATTR_COUNTER_ID;

    default:
      return SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_my_sid_entry_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_my_sid_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_my_sid_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_my_sid_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::MySidEntryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MySidEntryEndpointBehaviorFlavor
convert_sai_my_sid_entry_endpoint_behavior_flavor_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_NONE:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_NONE;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USP:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USP;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USD:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USD;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USP:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USP;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USD_AND_USP:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USD_AND_USP;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USD:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USD;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USP_AND_USD:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USP_AND_USD;

    default:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_UNSPECIFIED;
  }
}
sai_my_sid_entry_endpoint_behavior_flavor_t
convert_sai_my_sid_entry_endpoint_behavior_flavor_t_to_sai(
    lemming::dataplane::sai::MySidEntryEndpointBehaviorFlavor val) {
  switch (val) {
    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_NONE:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_NONE;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USP:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USP;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USD:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USD;

    case lemming::dataplane::sai::
        MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USP:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USP;

    case lemming::dataplane::sai::
        MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USD_AND_USP:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_USD_AND_USP;

    case lemming::dataplane::sai::
        MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USD:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USD;

    case lemming::dataplane::sai::
        MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USP_AND_USD:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_PSP_AND_USP_AND_USD;

    default:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_FLAVOR_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_my_sid_entry_endpoint_behavior_flavor_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_my_sid_entry_endpoint_behavior_flavor_t_to_proto(
        list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_my_sid_entry_endpoint_behavior_flavor_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_my_sid_entry_endpoint_behavior_flavor_t_to_sai(
        static_cast<lemming::dataplane::sai::MySidEntryEndpointBehaviorFlavor>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::MySidEntryEndpointBehavior
convert_sai_my_sid_entry_endpoint_behavior_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_E:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_E;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_X:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_X;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_T:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_T;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DX6:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DX6;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DX4:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DX4;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT6:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT6;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT4:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT4;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT46:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT46;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_ENCAPS:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_ENCAPS;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_ENCAPS_RED:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_ENCAPS_RED;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_INSERT:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_INSERT;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_INSERT_RED:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_INSERT_RED;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_UN:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_UN;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_UA:
      return lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_UA;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_CUSTOM_RANGE_START:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_CUSTOM_RANGE_START;

    case SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_CUSTOM_RANGE_END:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_CUSTOM_RANGE_END;

    default:
      return lemming::dataplane::sai::
          MY_SID_ENTRY_ENDPOINT_BEHAVIOR_UNSPECIFIED;
  }
}
sai_my_sid_entry_endpoint_behavior_t
convert_sai_my_sid_entry_endpoint_behavior_t_to_sai(
    lemming::dataplane::sai::MySidEntryEndpointBehavior val) {
  switch (val) {
    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_E:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_E;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_X:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_X;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_T:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_T;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DX6:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DX6;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DX4:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DX4;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT6:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT6;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT4:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT4;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT46:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_DT46;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_ENCAPS:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_ENCAPS;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_ENCAPS_RED:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_ENCAPS_RED;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_INSERT:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_INSERT;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_INSERT_RED:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_B6_INSERT_RED;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_UN:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_UN;

    case lemming::dataplane::sai::MY_SID_ENTRY_ENDPOINT_BEHAVIOR_UA:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_UA;

    case lemming::dataplane::sai::
        MY_SID_ENTRY_ENDPOINT_BEHAVIOR_CUSTOM_RANGE_START:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_CUSTOM_RANGE_START;

    case lemming::dataplane::sai::
        MY_SID_ENTRY_ENDPOINT_BEHAVIOR_CUSTOM_RANGE_END:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_CUSTOM_RANGE_END;

    default:
      return SAI_MY_SID_ENTRY_ENDPOINT_BEHAVIOR_E;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_my_sid_entry_endpoint_behavior_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_my_sid_entry_endpoint_behavior_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_my_sid_entry_endpoint_behavior_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_my_sid_entry_endpoint_behavior_t_to_sai(
        static_cast<lemming::dataplane::sai::MySidEntryEndpointBehavior>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NatEntryAttr convert_sai_nat_entry_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_NAT_ENTRY_ATTR_NAT_TYPE:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_NAT_TYPE;

    case SAI_NAT_ENTRY_ATTR_SRC_IP:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_SRC_IP;

    case SAI_NAT_ENTRY_ATTR_SRC_IP_MASK:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_SRC_IP_MASK;

    case SAI_NAT_ENTRY_ATTR_VR_ID:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_VR_ID;

    case SAI_NAT_ENTRY_ATTR_DST_IP:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_DST_IP;

    case SAI_NAT_ENTRY_ATTR_DST_IP_MASK:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_DST_IP_MASK;

    case SAI_NAT_ENTRY_ATTR_L4_SRC_PORT:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_L4_SRC_PORT;

    case SAI_NAT_ENTRY_ATTR_L4_DST_PORT:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_L4_DST_PORT;

    case SAI_NAT_ENTRY_ATTR_ENABLE_PACKET_COUNT:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_ENABLE_PACKET_COUNT;

    case SAI_NAT_ENTRY_ATTR_PACKET_COUNT:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_PACKET_COUNT;

    case SAI_NAT_ENTRY_ATTR_ENABLE_BYTE_COUNT:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_ENABLE_BYTE_COUNT;

    case SAI_NAT_ENTRY_ATTR_BYTE_COUNT:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_BYTE_COUNT;

    case SAI_NAT_ENTRY_ATTR_HIT_BIT_COR:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_HIT_BIT_COR;

    case SAI_NAT_ENTRY_ATTR_HIT_BIT:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_HIT_BIT;

    case SAI_NAT_ENTRY_ATTR_AGING_TIME:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_AGING_TIME;

    default:
      return lemming::dataplane::sai::NAT_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_nat_entry_attr_t convert_sai_nat_entry_attr_t_to_sai(
    lemming::dataplane::sai::NatEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::NAT_ENTRY_ATTR_NAT_TYPE:
      return SAI_NAT_ENTRY_ATTR_NAT_TYPE;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_SRC_IP:
      return SAI_NAT_ENTRY_ATTR_SRC_IP;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_SRC_IP_MASK:
      return SAI_NAT_ENTRY_ATTR_SRC_IP_MASK;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_VR_ID:
      return SAI_NAT_ENTRY_ATTR_VR_ID;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_DST_IP:
      return SAI_NAT_ENTRY_ATTR_DST_IP;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_DST_IP_MASK:
      return SAI_NAT_ENTRY_ATTR_DST_IP_MASK;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_L4_SRC_PORT:
      return SAI_NAT_ENTRY_ATTR_L4_SRC_PORT;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_L4_DST_PORT:
      return SAI_NAT_ENTRY_ATTR_L4_DST_PORT;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_ENABLE_PACKET_COUNT:
      return SAI_NAT_ENTRY_ATTR_ENABLE_PACKET_COUNT;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_PACKET_COUNT:
      return SAI_NAT_ENTRY_ATTR_PACKET_COUNT;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_ENABLE_BYTE_COUNT:
      return SAI_NAT_ENTRY_ATTR_ENABLE_BYTE_COUNT;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_BYTE_COUNT:
      return SAI_NAT_ENTRY_ATTR_BYTE_COUNT;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_HIT_BIT_COR:
      return SAI_NAT_ENTRY_ATTR_HIT_BIT_COR;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_HIT_BIT:
      return SAI_NAT_ENTRY_ATTR_HIT_BIT;

    case lemming::dataplane::sai::NAT_ENTRY_ATTR_AGING_TIME:
      return SAI_NAT_ENTRY_ATTR_AGING_TIME;

    default:
      return SAI_NAT_ENTRY_ATTR_NAT_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_nat_entry_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_nat_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_nat_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_nat_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::NatEntryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NatEvent convert_sai_nat_event_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_NAT_EVENT_NONE:
      return lemming::dataplane::sai::NAT_EVENT_NONE;

    case SAI_NAT_EVENT_AGED:
      return lemming::dataplane::sai::NAT_EVENT_AGED;

    default:
      return lemming::dataplane::sai::NAT_EVENT_UNSPECIFIED;
  }
}
sai_nat_event_t convert_sai_nat_event_t_to_sai(
    lemming::dataplane::sai::NatEvent val) {
  switch (val) {
    case lemming::dataplane::sai::NAT_EVENT_NONE:
      return SAI_NAT_EVENT_NONE;

    case lemming::dataplane::sai::NAT_EVENT_AGED:
      return SAI_NAT_EVENT_AGED;

    default:
      return SAI_NAT_EVENT_NONE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_nat_event_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_nat_event_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_nat_event_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_nat_event_t_to_sai(
        static_cast<lemming::dataplane::sai::NatEvent>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NatType convert_sai_nat_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_NAT_TYPE_NONE:
      return lemming::dataplane::sai::NAT_TYPE_NONE;

    case SAI_NAT_TYPE_SOURCE_NAT:
      return lemming::dataplane::sai::NAT_TYPE_SOURCE_NAT;

    case SAI_NAT_TYPE_DESTINATION_NAT:
      return lemming::dataplane::sai::NAT_TYPE_DESTINATION_NAT;

    case SAI_NAT_TYPE_DOUBLE_NAT:
      return lemming::dataplane::sai::NAT_TYPE_DOUBLE_NAT;

    case SAI_NAT_TYPE_DESTINATION_NAT_POOL:
      return lemming::dataplane::sai::NAT_TYPE_DESTINATION_NAT_POOL;

    default:
      return lemming::dataplane::sai::NAT_TYPE_UNSPECIFIED;
  }
}
sai_nat_type_t convert_sai_nat_type_t_to_sai(
    lemming::dataplane::sai::NatType val) {
  switch (val) {
    case lemming::dataplane::sai::NAT_TYPE_NONE:
      return SAI_NAT_TYPE_NONE;

    case lemming::dataplane::sai::NAT_TYPE_SOURCE_NAT:
      return SAI_NAT_TYPE_SOURCE_NAT;

    case lemming::dataplane::sai::NAT_TYPE_DESTINATION_NAT:
      return SAI_NAT_TYPE_DESTINATION_NAT;

    case lemming::dataplane::sai::NAT_TYPE_DOUBLE_NAT:
      return SAI_NAT_TYPE_DOUBLE_NAT;

    case lemming::dataplane::sai::NAT_TYPE_DESTINATION_NAT_POOL:
      return SAI_NAT_TYPE_DESTINATION_NAT_POOL;

    default:
      return SAI_NAT_TYPE_NONE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_nat_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_nat_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_nat_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_nat_type_t_to_sai(
        static_cast<lemming::dataplane::sai::NatType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NatZoneCounterAttr
convert_sai_nat_zone_counter_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_NAT_ZONE_COUNTER_ATTR_NAT_TYPE:
      return lemming::dataplane::sai::NAT_ZONE_COUNTER_ATTR_NAT_TYPE;

    case SAI_NAT_ZONE_COUNTER_ATTR_ZONE_ID:
      return lemming::dataplane::sai::NAT_ZONE_COUNTER_ATTR_ZONE_ID;

    case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_DISCARD:
      return lemming::dataplane::sai::NAT_ZONE_COUNTER_ATTR_ENABLE_DISCARD;

    case SAI_NAT_ZONE_COUNTER_ATTR_DISCARD_PACKET_COUNT:
      return lemming::dataplane::sai::
          NAT_ZONE_COUNTER_ATTR_DISCARD_PACKET_COUNT;

    case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATION_NEEDED:
      return lemming::dataplane::sai::
          NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATION_NEEDED;

    case SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATION_NEEDED_PACKET_COUNT:
      return lemming::dataplane::sai::
          NAT_ZONE_COUNTER_ATTR_TRANSLATION_NEEDED_PACKET_COUNT;

    case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATIONS:
      return lemming::dataplane::sai::NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATIONS;

    case SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATIONS_PACKET_COUNT:
      return lemming::dataplane::sai::
          NAT_ZONE_COUNTER_ATTR_TRANSLATIONS_PACKET_COUNT;

    default:
      return lemming::dataplane::sai::NAT_ZONE_COUNTER_ATTR_UNSPECIFIED;
  }
}
sai_nat_zone_counter_attr_t convert_sai_nat_zone_counter_attr_t_to_sai(
    lemming::dataplane::sai::NatZoneCounterAttr val) {
  switch (val) {
    case lemming::dataplane::sai::NAT_ZONE_COUNTER_ATTR_NAT_TYPE:
      return SAI_NAT_ZONE_COUNTER_ATTR_NAT_TYPE;

    case lemming::dataplane::sai::NAT_ZONE_COUNTER_ATTR_ZONE_ID:
      return SAI_NAT_ZONE_COUNTER_ATTR_ZONE_ID;

    case lemming::dataplane::sai::NAT_ZONE_COUNTER_ATTR_ENABLE_DISCARD:
      return SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_DISCARD;

    case lemming::dataplane::sai::NAT_ZONE_COUNTER_ATTR_DISCARD_PACKET_COUNT:
      return SAI_NAT_ZONE_COUNTER_ATTR_DISCARD_PACKET_COUNT;

    case lemming::dataplane::sai::
        NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATION_NEEDED:
      return SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATION_NEEDED;

    case lemming::dataplane::sai::
        NAT_ZONE_COUNTER_ATTR_TRANSLATION_NEEDED_PACKET_COUNT:
      return SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATION_NEEDED_PACKET_COUNT;

    case lemming::dataplane::sai::NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATIONS:
      return SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATIONS;

    case lemming::dataplane::sai::
        NAT_ZONE_COUNTER_ATTR_TRANSLATIONS_PACKET_COUNT:
      return SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATIONS_PACKET_COUNT;

    default:
      return SAI_NAT_ZONE_COUNTER_ATTR_NAT_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_nat_zone_counter_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_nat_zone_counter_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_nat_zone_counter_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_nat_zone_counter_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::NatZoneCounterAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NativeHashField
convert_sai_native_hash_field_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_NATIVE_HASH_FIELD_SRC_IP:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_SRC_IP;

    case SAI_NATIVE_HASH_FIELD_DST_IP:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_DST_IP;

    case SAI_NATIVE_HASH_FIELD_INNER_SRC_IP:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_SRC_IP;

    case SAI_NATIVE_HASH_FIELD_INNER_DST_IP:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_DST_IP;

    case SAI_NATIVE_HASH_FIELD_SRC_IPV4:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_SRC_IPV4;

    case SAI_NATIVE_HASH_FIELD_DST_IPV4:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_DST_IPV4;

    case SAI_NATIVE_HASH_FIELD_SRC_IPV6:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_SRC_IPV6;

    case SAI_NATIVE_HASH_FIELD_DST_IPV6:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_DST_IPV6;

    case SAI_NATIVE_HASH_FIELD_INNER_SRC_IPV4:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_SRC_IPV4;

    case SAI_NATIVE_HASH_FIELD_INNER_DST_IPV4:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_DST_IPV4;

    case SAI_NATIVE_HASH_FIELD_INNER_SRC_IPV6:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_SRC_IPV6;

    case SAI_NATIVE_HASH_FIELD_INNER_DST_IPV6:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_DST_IPV6;

    case SAI_NATIVE_HASH_FIELD_VLAN_ID:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_VLAN_ID;

    case SAI_NATIVE_HASH_FIELD_IP_PROTOCOL:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_IP_PROTOCOL;

    case SAI_NATIVE_HASH_FIELD_ETHERTYPE:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_ETHERTYPE;

    case SAI_NATIVE_HASH_FIELD_L4_SRC_PORT:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_L4_SRC_PORT;

    case SAI_NATIVE_HASH_FIELD_L4_DST_PORT:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_L4_DST_PORT;

    case SAI_NATIVE_HASH_FIELD_SRC_MAC:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_SRC_MAC;

    case SAI_NATIVE_HASH_FIELD_DST_MAC:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_DST_MAC;

    case SAI_NATIVE_HASH_FIELD_IN_PORT:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_IN_PORT;

    case SAI_NATIVE_HASH_FIELD_INNER_IP_PROTOCOL:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_IP_PROTOCOL;

    case SAI_NATIVE_HASH_FIELD_INNER_ETHERTYPE:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_ETHERTYPE;

    case SAI_NATIVE_HASH_FIELD_INNER_L4_SRC_PORT:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_L4_SRC_PORT;

    case SAI_NATIVE_HASH_FIELD_INNER_L4_DST_PORT:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_L4_DST_PORT;

    case SAI_NATIVE_HASH_FIELD_INNER_SRC_MAC:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_SRC_MAC;

    case SAI_NATIVE_HASH_FIELD_INNER_DST_MAC:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_DST_MAC;

    case SAI_NATIVE_HASH_FIELD_MPLS_LABEL_ALL:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_ALL;

    case SAI_NATIVE_HASH_FIELD_MPLS_LABEL_0:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_0;

    case SAI_NATIVE_HASH_FIELD_MPLS_LABEL_1:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_1;

    case SAI_NATIVE_HASH_FIELD_MPLS_LABEL_2:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_2;

    case SAI_NATIVE_HASH_FIELD_MPLS_LABEL_3:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_3;

    case SAI_NATIVE_HASH_FIELD_MPLS_LABEL_4:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_4;

    case SAI_NATIVE_HASH_FIELD_IPV6_FLOW_LABEL:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_IPV6_FLOW_LABEL;

    case SAI_NATIVE_HASH_FIELD_NONE:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_NONE;

    default:
      return lemming::dataplane::sai::NATIVE_HASH_FIELD_UNSPECIFIED;
  }
}
sai_native_hash_field_t convert_sai_native_hash_field_t_to_sai(
    lemming::dataplane::sai::NativeHashField val) {
  switch (val) {
    case lemming::dataplane::sai::NATIVE_HASH_FIELD_SRC_IP:
      return SAI_NATIVE_HASH_FIELD_SRC_IP;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_DST_IP:
      return SAI_NATIVE_HASH_FIELD_DST_IP;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_SRC_IP:
      return SAI_NATIVE_HASH_FIELD_INNER_SRC_IP;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_DST_IP:
      return SAI_NATIVE_HASH_FIELD_INNER_DST_IP;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_SRC_IPV4:
      return SAI_NATIVE_HASH_FIELD_SRC_IPV4;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_DST_IPV4:
      return SAI_NATIVE_HASH_FIELD_DST_IPV4;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_SRC_IPV6:
      return SAI_NATIVE_HASH_FIELD_SRC_IPV6;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_DST_IPV6:
      return SAI_NATIVE_HASH_FIELD_DST_IPV6;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_SRC_IPV4:
      return SAI_NATIVE_HASH_FIELD_INNER_SRC_IPV4;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_DST_IPV4:
      return SAI_NATIVE_HASH_FIELD_INNER_DST_IPV4;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_SRC_IPV6:
      return SAI_NATIVE_HASH_FIELD_INNER_SRC_IPV6;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_DST_IPV6:
      return SAI_NATIVE_HASH_FIELD_INNER_DST_IPV6;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_VLAN_ID:
      return SAI_NATIVE_HASH_FIELD_VLAN_ID;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_IP_PROTOCOL:
      return SAI_NATIVE_HASH_FIELD_IP_PROTOCOL;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_ETHERTYPE:
      return SAI_NATIVE_HASH_FIELD_ETHERTYPE;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_L4_SRC_PORT:
      return SAI_NATIVE_HASH_FIELD_L4_SRC_PORT;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_L4_DST_PORT:
      return SAI_NATIVE_HASH_FIELD_L4_DST_PORT;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_SRC_MAC:
      return SAI_NATIVE_HASH_FIELD_SRC_MAC;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_DST_MAC:
      return SAI_NATIVE_HASH_FIELD_DST_MAC;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_IN_PORT:
      return SAI_NATIVE_HASH_FIELD_IN_PORT;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_IP_PROTOCOL:
      return SAI_NATIVE_HASH_FIELD_INNER_IP_PROTOCOL;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_ETHERTYPE:
      return SAI_NATIVE_HASH_FIELD_INNER_ETHERTYPE;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_L4_SRC_PORT:
      return SAI_NATIVE_HASH_FIELD_INNER_L4_SRC_PORT;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_L4_DST_PORT:
      return SAI_NATIVE_HASH_FIELD_INNER_L4_DST_PORT;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_SRC_MAC:
      return SAI_NATIVE_HASH_FIELD_INNER_SRC_MAC;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_INNER_DST_MAC:
      return SAI_NATIVE_HASH_FIELD_INNER_DST_MAC;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_ALL:
      return SAI_NATIVE_HASH_FIELD_MPLS_LABEL_ALL;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_0:
      return SAI_NATIVE_HASH_FIELD_MPLS_LABEL_0;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_1:
      return SAI_NATIVE_HASH_FIELD_MPLS_LABEL_1;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_2:
      return SAI_NATIVE_HASH_FIELD_MPLS_LABEL_2;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_3:
      return SAI_NATIVE_HASH_FIELD_MPLS_LABEL_3;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_MPLS_LABEL_4:
      return SAI_NATIVE_HASH_FIELD_MPLS_LABEL_4;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_IPV6_FLOW_LABEL:
      return SAI_NATIVE_HASH_FIELD_IPV6_FLOW_LABEL;

    case lemming::dataplane::sai::NATIVE_HASH_FIELD_NONE:
      return SAI_NATIVE_HASH_FIELD_NONE;

    default:
      return SAI_NATIVE_HASH_FIELD_SRC_IP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_native_hash_field_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_native_hash_field_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_native_hash_field_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_native_hash_field_t_to_sai(
        static_cast<lemming::dataplane::sai::NativeHashField>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NeighborEntryAttr
convert_sai_neighbor_entry_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_NEIGHBOR_ENTRY_ATTR_DST_MAC_ADDRESS:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_DST_MAC_ADDRESS;

    case SAI_NEIGHBOR_ENTRY_ATTR_PACKET_ACTION:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_PACKET_ACTION;

    case SAI_NEIGHBOR_ENTRY_ATTR_USER_TRAP_ID:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_USER_TRAP_ID;

    case SAI_NEIGHBOR_ENTRY_ATTR_NO_HOST_ROUTE:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_NO_HOST_ROUTE;

    case SAI_NEIGHBOR_ENTRY_ATTR_META_DATA:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_META_DATA;

    case SAI_NEIGHBOR_ENTRY_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_COUNTER_ID;

    case SAI_NEIGHBOR_ENTRY_ATTR_ENCAP_INDEX:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_ENCAP_INDEX;

    case SAI_NEIGHBOR_ENTRY_ATTR_ENCAP_IMPOSE_INDEX:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_ENCAP_IMPOSE_INDEX;

    case SAI_NEIGHBOR_ENTRY_ATTR_IS_LOCAL:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_IS_LOCAL;

    case SAI_NEIGHBOR_ENTRY_ATTR_IP_ADDR_FAMILY:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_IP_ADDR_FAMILY;

    default:
      return lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_neighbor_entry_attr_t convert_sai_neighbor_entry_attr_t_to_sai(
    lemming::dataplane::sai::NeighborEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_DST_MAC_ADDRESS:
      return SAI_NEIGHBOR_ENTRY_ATTR_DST_MAC_ADDRESS;

    case lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_PACKET_ACTION:
      return SAI_NEIGHBOR_ENTRY_ATTR_PACKET_ACTION;

    case lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_USER_TRAP_ID:
      return SAI_NEIGHBOR_ENTRY_ATTR_USER_TRAP_ID;

    case lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_NO_HOST_ROUTE:
      return SAI_NEIGHBOR_ENTRY_ATTR_NO_HOST_ROUTE;

    case lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_META_DATA:
      return SAI_NEIGHBOR_ENTRY_ATTR_META_DATA;

    case lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_COUNTER_ID:
      return SAI_NEIGHBOR_ENTRY_ATTR_COUNTER_ID;

    case lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_ENCAP_INDEX:
      return SAI_NEIGHBOR_ENTRY_ATTR_ENCAP_INDEX;

    case lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_ENCAP_IMPOSE_INDEX:
      return SAI_NEIGHBOR_ENTRY_ATTR_ENCAP_IMPOSE_INDEX;

    case lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_IS_LOCAL:
      return SAI_NEIGHBOR_ENTRY_ATTR_IS_LOCAL;

    case lemming::dataplane::sai::NEIGHBOR_ENTRY_ATTR_IP_ADDR_FAMILY:
      return SAI_NEIGHBOR_ENTRY_ATTR_IP_ADDR_FAMILY;

    default:
      return SAI_NEIGHBOR_ENTRY_ATTR_DST_MAC_ADDRESS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_neighbor_entry_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_neighbor_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_neighbor_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_neighbor_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::NeighborEntryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NextHopAttr convert_sai_next_hop_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_NEXT_HOP_ATTR_TYPE:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_TYPE;

    case SAI_NEXT_HOP_ATTR_IP:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_IP;

    case SAI_NEXT_HOP_ATTR_ROUTER_INTERFACE_ID:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_ROUTER_INTERFACE_ID;

    case SAI_NEXT_HOP_ATTR_TUNNEL_ID:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_TUNNEL_ID;

    case SAI_NEXT_HOP_ATTR_TUNNEL_VNI:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_TUNNEL_VNI;

    case SAI_NEXT_HOP_ATTR_TUNNEL_MAC:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_TUNNEL_MAC;

    case SAI_NEXT_HOP_ATTR_SRV6_SIDLIST_ID:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_SRV6_SIDLIST_ID;

    case SAI_NEXT_HOP_ATTR_LABELSTACK:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_LABELSTACK;

    case SAI_NEXT_HOP_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_COUNTER_ID;

    case SAI_NEXT_HOP_ATTR_DISABLE_DECREMENT_TTL:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_DISABLE_DECREMENT_TTL;

    case SAI_NEXT_HOP_ATTR_OUTSEG_TYPE:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_OUTSEG_TYPE;

    case SAI_NEXT_HOP_ATTR_OUTSEG_TTL_MODE:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_OUTSEG_TTL_MODE;

    case SAI_NEXT_HOP_ATTR_OUTSEG_TTL_VALUE:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_OUTSEG_TTL_VALUE;

    case SAI_NEXT_HOP_ATTR_OUTSEG_EXP_MODE:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_OUTSEG_EXP_MODE;

    case SAI_NEXT_HOP_ATTR_OUTSEG_EXP_VALUE:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_OUTSEG_EXP_VALUE;

    case SAI_NEXT_HOP_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
      return lemming::dataplane::sai::
          NEXT_HOP_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP;

    default:
      return lemming::dataplane::sai::NEXT_HOP_ATTR_UNSPECIFIED;
  }
}
sai_next_hop_attr_t convert_sai_next_hop_attr_t_to_sai(
    lemming::dataplane::sai::NextHopAttr val) {
  switch (val) {
    case lemming::dataplane::sai::NEXT_HOP_ATTR_TYPE:
      return SAI_NEXT_HOP_ATTR_TYPE;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_IP:
      return SAI_NEXT_HOP_ATTR_IP;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_ROUTER_INTERFACE_ID:
      return SAI_NEXT_HOP_ATTR_ROUTER_INTERFACE_ID;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_TUNNEL_ID:
      return SAI_NEXT_HOP_ATTR_TUNNEL_ID;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_TUNNEL_VNI:
      return SAI_NEXT_HOP_ATTR_TUNNEL_VNI;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_TUNNEL_MAC:
      return SAI_NEXT_HOP_ATTR_TUNNEL_MAC;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_SRV6_SIDLIST_ID:
      return SAI_NEXT_HOP_ATTR_SRV6_SIDLIST_ID;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_LABELSTACK:
      return SAI_NEXT_HOP_ATTR_LABELSTACK;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_COUNTER_ID:
      return SAI_NEXT_HOP_ATTR_COUNTER_ID;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_DISABLE_DECREMENT_TTL:
      return SAI_NEXT_HOP_ATTR_DISABLE_DECREMENT_TTL;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_OUTSEG_TYPE:
      return SAI_NEXT_HOP_ATTR_OUTSEG_TYPE;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_OUTSEG_TTL_MODE:
      return SAI_NEXT_HOP_ATTR_OUTSEG_TTL_MODE;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_OUTSEG_TTL_VALUE:
      return SAI_NEXT_HOP_ATTR_OUTSEG_TTL_VALUE;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_OUTSEG_EXP_MODE:
      return SAI_NEXT_HOP_ATTR_OUTSEG_EXP_MODE;

    case lemming::dataplane::sai::NEXT_HOP_ATTR_OUTSEG_EXP_VALUE:
      return SAI_NEXT_HOP_ATTR_OUTSEG_EXP_VALUE;

    case lemming::dataplane::sai::
        NEXT_HOP_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
      return SAI_NEXT_HOP_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP;

    default:
      return SAI_NEXT_HOP_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_next_hop_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_next_hop_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_next_hop_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_next_hop_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::NextHopAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NextHopGroupAttr
convert_sai_next_hop_group_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_COUNT:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_NEXT_HOP_COUNT;

    case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_LIST:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_LIST;

    case SAI_NEXT_HOP_GROUP_ATTR_TYPE:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_TYPE;

    case SAI_NEXT_HOP_GROUP_ATTR_SET_SWITCHOVER:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_SET_SWITCHOVER;

    case SAI_NEXT_HOP_GROUP_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_COUNTER_ID;

    case SAI_NEXT_HOP_GROUP_ATTR_CONFIGURED_SIZE:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_CONFIGURED_SIZE;

    case SAI_NEXT_HOP_GROUP_ATTR_REAL_SIZE:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_REAL_SIZE;

    case SAI_NEXT_HOP_GROUP_ATTR_SELECTION_MAP:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_SELECTION_MAP;

    case SAI_NEXT_HOP_GROUP_ATTR_HIERARCHICAL_NEXTHOP:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_HIERARCHICAL_NEXTHOP;

    case SAI_NEXT_HOP_GROUP_ATTR_ARS_OBJECT_ID:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_ARS_OBJECT_ID;

    case SAI_NEXT_HOP_GROUP_ATTR_ARS_PACKET_DROPS:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_ARS_PACKET_DROPS;

    case SAI_NEXT_HOP_GROUP_ATTR_ARS_NEXT_HOP_REASSIGNMENTS:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_ATTR_ARS_NEXT_HOP_REASSIGNMENTS;

    case SAI_NEXT_HOP_GROUP_ATTR_ARS_PORT_REASSIGNMENTS:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_ATTR_ARS_PORT_REASSIGNMENTS;

    case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_LIST:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_NEXT_HOP_LIST;

    case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_WEIGHT_LIST:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_WEIGHT_LIST;

    case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_COUNTER_LIST:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_COUNTER_LIST;

    default:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_UNSPECIFIED;
  }
}
sai_next_hop_group_attr_t convert_sai_next_hop_group_attr_t_to_sai(
    lemming::dataplane::sai::NextHopGroupAttr val) {
  switch (val) {
    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_NEXT_HOP_COUNT:
      return SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_COUNT;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_LIST:
      return SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_LIST;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_TYPE:
      return SAI_NEXT_HOP_GROUP_ATTR_TYPE;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_SET_SWITCHOVER:
      return SAI_NEXT_HOP_GROUP_ATTR_SET_SWITCHOVER;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_COUNTER_ID:
      return SAI_NEXT_HOP_GROUP_ATTR_COUNTER_ID;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_CONFIGURED_SIZE:
      return SAI_NEXT_HOP_GROUP_ATTR_CONFIGURED_SIZE;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_REAL_SIZE:
      return SAI_NEXT_HOP_GROUP_ATTR_REAL_SIZE;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_SELECTION_MAP:
      return SAI_NEXT_HOP_GROUP_ATTR_SELECTION_MAP;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_HIERARCHICAL_NEXTHOP:
      return SAI_NEXT_HOP_GROUP_ATTR_HIERARCHICAL_NEXTHOP;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_ARS_OBJECT_ID:
      return SAI_NEXT_HOP_GROUP_ATTR_ARS_OBJECT_ID;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_ARS_PACKET_DROPS:
      return SAI_NEXT_HOP_GROUP_ATTR_ARS_PACKET_DROPS;

    case lemming::dataplane::sai::
        NEXT_HOP_GROUP_ATTR_ARS_NEXT_HOP_REASSIGNMENTS:
      return SAI_NEXT_HOP_GROUP_ATTR_ARS_NEXT_HOP_REASSIGNMENTS;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_ARS_PORT_REASSIGNMENTS:
      return SAI_NEXT_HOP_GROUP_ATTR_ARS_PORT_REASSIGNMENTS;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_ATTR_NEXT_HOP_LIST:
      return SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_LIST;

    case lemming::dataplane::sai::
        NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_WEIGHT_LIST:
      return SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_WEIGHT_LIST;

    case lemming::dataplane::sai::
        NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_COUNTER_LIST:
      return SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_COUNTER_LIST;

    default:
      return SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_COUNT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_next_hop_group_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_next_hop_group_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_next_hop_group_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::NextHopGroupAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NextHopGroupMapAttr
convert_sai_next_hop_group_map_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_NEXT_HOP_GROUP_MAP_ATTR_TYPE:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MAP_ATTR_TYPE;

    case SAI_NEXT_HOP_GROUP_MAP_ATTR_MAP_TO_VALUE_LIST:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MAP_ATTR_MAP_TO_VALUE_LIST;

    default:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MAP_ATTR_UNSPECIFIED;
  }
}
sai_next_hop_group_map_attr_t convert_sai_next_hop_group_map_attr_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMapAttr val) {
  switch (val) {
    case lemming::dataplane::sai::NEXT_HOP_GROUP_MAP_ATTR_TYPE:
      return SAI_NEXT_HOP_GROUP_MAP_ATTR_TYPE;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MAP_ATTR_MAP_TO_VALUE_LIST:
      return SAI_NEXT_HOP_GROUP_MAP_ATTR_MAP_TO_VALUE_LIST;

    default:
      return SAI_NEXT_HOP_GROUP_MAP_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_map_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_next_hop_group_map_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_next_hop_group_map_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_next_hop_group_map_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::NextHopGroupMapAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NextHopGroupMapType
convert_sai_next_hop_group_map_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_NEXT_HOP_GROUP_MAP_TYPE_FORWARDING_CLASS_TO_INDEX:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MAP_TYPE_FORWARDING_CLASS_TO_INDEX;

    default:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MAP_TYPE_UNSPECIFIED;
  }
}
sai_next_hop_group_map_type_t convert_sai_next_hop_group_map_type_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMapType val) {
  switch (val) {
    case lemming::dataplane::sai::
        NEXT_HOP_GROUP_MAP_TYPE_FORWARDING_CLASS_TO_INDEX:
      return SAI_NEXT_HOP_GROUP_MAP_TYPE_FORWARDING_CLASS_TO_INDEX;

    default:
      return SAI_NEXT_HOP_GROUP_MAP_TYPE_FORWARDING_CLASS_TO_INDEX;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_map_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_next_hop_group_map_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_next_hop_group_map_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_next_hop_group_map_type_t_to_sai(
        static_cast<lemming::dataplane::sai::NextHopGroupMapType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NextHopGroupMemberAttr
convert_sai_next_hop_group_member_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_GROUP_ID:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_GROUP_ID;

    case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_ID:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_ID;

    case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_WEIGHT:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_WEIGHT;

    case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_CONFIGURED_ROLE:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MEMBER_ATTR_CONFIGURED_ROLE;

    case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_OBSERVED_ROLE:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_OBSERVED_ROLE;

    case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_MONITORED_OBJECT:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MEMBER_ATTR_MONITORED_OBJECT;

    case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_INDEX:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_INDEX;

    case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_SEQUENCE_ID:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_SEQUENCE_ID;

    case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_COUNTER_ID;

    case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_ARS_ALTERNATE_PATH:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MEMBER_ATTR_ARS_ALTERNATE_PATH;

    default:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_UNSPECIFIED;
  }
}
sai_next_hop_group_member_attr_t
convert_sai_next_hop_group_member_attr_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMemberAttr val) {
  switch (val) {
    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_GROUP_ID:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_GROUP_ID;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_ID:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_ID;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_WEIGHT:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_WEIGHT;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_CONFIGURED_ROLE:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_CONFIGURED_ROLE;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_OBSERVED_ROLE:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_OBSERVED_ROLE;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_MONITORED_OBJECT:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_MONITORED_OBJECT;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_INDEX:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_INDEX;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_SEQUENCE_ID:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_SEQUENCE_ID;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_COUNTER_ID:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_COUNTER_ID;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_ATTR_ARS_ALTERNATE_PATH:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_ARS_ALTERNATE_PATH;

    default:
      return SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_GROUP_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_member_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_next_hop_group_member_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_next_hop_group_member_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_next_hop_group_member_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::NextHopGroupMemberAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NextHopGroupMemberConfiguredRole
convert_sai_next_hop_group_member_configured_role_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_NEXT_HOP_GROUP_MEMBER_CONFIGURED_ROLE_PRIMARY:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MEMBER_CONFIGURED_ROLE_PRIMARY;

    case SAI_NEXT_HOP_GROUP_MEMBER_CONFIGURED_ROLE_STANDBY:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MEMBER_CONFIGURED_ROLE_STANDBY;

    default:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MEMBER_CONFIGURED_ROLE_UNSPECIFIED;
  }
}
sai_next_hop_group_member_configured_role_t
convert_sai_next_hop_group_member_configured_role_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMemberConfiguredRole val) {
  switch (val) {
    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_CONFIGURED_ROLE_PRIMARY:
      return SAI_NEXT_HOP_GROUP_MEMBER_CONFIGURED_ROLE_PRIMARY;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_CONFIGURED_ROLE_STANDBY:
      return SAI_NEXT_HOP_GROUP_MEMBER_CONFIGURED_ROLE_STANDBY;

    default:
      return SAI_NEXT_HOP_GROUP_MEMBER_CONFIGURED_ROLE_PRIMARY;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_member_configured_role_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_next_hop_group_member_configured_role_t_to_proto(
        list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_next_hop_group_member_configured_role_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_next_hop_group_member_configured_role_t_to_sai(
        static_cast<lemming::dataplane::sai::NextHopGroupMemberConfiguredRole>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NextHopGroupMemberObservedRole
convert_sai_next_hop_group_member_observed_role_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_NEXT_HOP_GROUP_MEMBER_OBSERVED_ROLE_ACTIVE:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MEMBER_OBSERVED_ROLE_ACTIVE;

    case SAI_NEXT_HOP_GROUP_MEMBER_OBSERVED_ROLE_INACTIVE:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MEMBER_OBSERVED_ROLE_INACTIVE;

    default:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_MEMBER_OBSERVED_ROLE_UNSPECIFIED;
  }
}
sai_next_hop_group_member_observed_role_t
convert_sai_next_hop_group_member_observed_role_t_to_sai(
    lemming::dataplane::sai::NextHopGroupMemberObservedRole val) {
  switch (val) {
    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_OBSERVED_ROLE_ACTIVE:
      return SAI_NEXT_HOP_GROUP_MEMBER_OBSERVED_ROLE_ACTIVE;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_MEMBER_OBSERVED_ROLE_INACTIVE:
      return SAI_NEXT_HOP_GROUP_MEMBER_OBSERVED_ROLE_INACTIVE;

    default:
      return SAI_NEXT_HOP_GROUP_MEMBER_OBSERVED_ROLE_ACTIVE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_member_observed_role_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_next_hop_group_member_observed_role_t_to_proto(
        list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_next_hop_group_member_observed_role_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_next_hop_group_member_observed_role_t_to_sai(
        static_cast<lemming::dataplane::sai::NextHopGroupMemberObservedRole>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NextHopGroupType
convert_sai_next_hop_group_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP:
      return lemming::dataplane::sai::
          NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP;

    case SAI_NEXT_HOP_GROUP_TYPE_DYNAMIC_ORDERED_ECMP:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_DYNAMIC_ORDERED_ECMP;

    case SAI_NEXT_HOP_GROUP_TYPE_FINE_GRAIN_ECMP:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_FINE_GRAIN_ECMP;

    case SAI_NEXT_HOP_GROUP_TYPE_PROTECTION:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_PROTECTION;

    case SAI_NEXT_HOP_GROUP_TYPE_CLASS_BASED:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_CLASS_BASED;

    case SAI_NEXT_HOP_GROUP_TYPE_ECMP_WITH_MEMBERS:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_ECMP_WITH_MEMBERS;

    default:
      return lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_UNSPECIFIED;
  }
}
sai_next_hop_group_type_t convert_sai_next_hop_group_type_t_to_sai(
    lemming::dataplane::sai::NextHopGroupType val) {
  switch (val) {
    case lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP:
      return SAI_NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_DYNAMIC_ORDERED_ECMP:
      return SAI_NEXT_HOP_GROUP_TYPE_DYNAMIC_ORDERED_ECMP;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_FINE_GRAIN_ECMP:
      return SAI_NEXT_HOP_GROUP_TYPE_FINE_GRAIN_ECMP;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_PROTECTION:
      return SAI_NEXT_HOP_GROUP_TYPE_PROTECTION;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_CLASS_BASED:
      return SAI_NEXT_HOP_GROUP_TYPE_CLASS_BASED;

    case lemming::dataplane::sai::NEXT_HOP_GROUP_TYPE_ECMP_WITH_MEMBERS:
      return SAI_NEXT_HOP_GROUP_TYPE_ECMP_WITH_MEMBERS;

    default:
      return SAI_NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_next_hop_group_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_next_hop_group_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_next_hop_group_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_next_hop_group_type_t_to_sai(
        static_cast<lemming::dataplane::sai::NextHopGroupType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::NextHopType convert_sai_next_hop_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_NEXT_HOP_TYPE_IP:
      return lemming::dataplane::sai::NEXT_HOP_TYPE_IP;

    case SAI_NEXT_HOP_TYPE_MPLS:
      return lemming::dataplane::sai::NEXT_HOP_TYPE_MPLS;

    case SAI_NEXT_HOP_TYPE_TUNNEL_ENCAP:
      return lemming::dataplane::sai::NEXT_HOP_TYPE_TUNNEL_ENCAP;

    case SAI_NEXT_HOP_TYPE_SRV6_SIDLIST:
      return lemming::dataplane::sai::NEXT_HOP_TYPE_SRV6_SIDLIST;

    default:
      return lemming::dataplane::sai::NEXT_HOP_TYPE_UNSPECIFIED;
  }
}
sai_next_hop_type_t convert_sai_next_hop_type_t_to_sai(
    lemming::dataplane::sai::NextHopType val) {
  switch (val) {
    case lemming::dataplane::sai::NEXT_HOP_TYPE_IP:
      return SAI_NEXT_HOP_TYPE_IP;

    case lemming::dataplane::sai::NEXT_HOP_TYPE_MPLS:
      return SAI_NEXT_HOP_TYPE_MPLS;

    case lemming::dataplane::sai::NEXT_HOP_TYPE_TUNNEL_ENCAP:
      return SAI_NEXT_HOP_TYPE_TUNNEL_ENCAP;

    case lemming::dataplane::sai::NEXT_HOP_TYPE_SRV6_SIDLIST:
      return SAI_NEXT_HOP_TYPE_SRV6_SIDLIST;

    default:
      return SAI_NEXT_HOP_TYPE_IP;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_next_hop_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_next_hop_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_next_hop_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_next_hop_type_t_to_sai(
        static_cast<lemming::dataplane::sai::NextHopType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::ObjectStage convert_sai_object_stage_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_OBJECT_STAGE_BOTH:
      return lemming::dataplane::sai::OBJECT_STAGE_BOTH;

    case SAI_OBJECT_STAGE_INGRESS:
      return lemming::dataplane::sai::OBJECT_STAGE_INGRESS;

    case SAI_OBJECT_STAGE_EGRESS:
      return lemming::dataplane::sai::OBJECT_STAGE_EGRESS;

    default:
      return lemming::dataplane::sai::OBJECT_STAGE_UNSPECIFIED;
  }
}
sai_object_stage_t convert_sai_object_stage_t_to_sai(
    lemming::dataplane::sai::ObjectStage val) {
  switch (val) {
    case lemming::dataplane::sai::OBJECT_STAGE_BOTH:
      return SAI_OBJECT_STAGE_BOTH;

    case lemming::dataplane::sai::OBJECT_STAGE_INGRESS:
      return SAI_OBJECT_STAGE_INGRESS;

    case lemming::dataplane::sai::OBJECT_STAGE_EGRESS:
      return SAI_OBJECT_STAGE_EGRESS;

    default:
      return SAI_OBJECT_STAGE_BOTH;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_object_stage_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_object_stage_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_object_stage_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_object_stage_t_to_sai(
        static_cast<lemming::dataplane::sai::ObjectStage>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::ObjectTypeExtensions
convert_sai_object_type_extensions_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_OBJECT_TYPE_EXTENSIONS_RANGE_START:
      return lemming::dataplane::sai::OBJECT_TYPE_EXTENSIONS_RANGE_START;

    case SAI_OBJECT_TYPE_TABLE_BITMAP_ROUTER_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_TABLE_BITMAP_ROUTER_ENTRY;

    case SAI_OBJECT_TYPE_TABLE_META_TUNNEL_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_TABLE_META_TUNNEL_ENTRY;

    case SAI_OBJECT_TYPE_EXTENSIONS_RANGE_END:
      return lemming::dataplane::sai::OBJECT_TYPE_EXTENSIONS_RANGE_END;

    default:
      return lemming::dataplane::sai::OBJECT_TYPE_EXTENSIONS_UNSPECIFIED;
  }
}
sai_object_type_extensions_t convert_sai_object_type_extensions_t_to_sai(
    lemming::dataplane::sai::ObjectTypeExtensions val) {
  switch (val) {
    case lemming::dataplane::sai::OBJECT_TYPE_EXTENSIONS_RANGE_START:
      return SAI_OBJECT_TYPE_EXTENSIONS_RANGE_START;

    case lemming::dataplane::sai::OBJECT_TYPE_TABLE_BITMAP_ROUTER_ENTRY:
      return SAI_OBJECT_TYPE_TABLE_BITMAP_ROUTER_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_TABLE_META_TUNNEL_ENTRY:
      return SAI_OBJECT_TYPE_TABLE_META_TUNNEL_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_EXTENSIONS_RANGE_END:
      return SAI_OBJECT_TYPE_EXTENSIONS_RANGE_END;

    default:
      return SAI_OBJECT_TYPE_EXTENSIONS_RANGE_START;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_object_type_extensions_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_object_type_extensions_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_object_type_extensions_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_object_type_extensions_t_to_sai(
        static_cast<lemming::dataplane::sai::ObjectTypeExtensions>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::ObjectType convert_sai_object_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_OBJECT_TYPE_NULL:
      return lemming::dataplane::sai::OBJECT_TYPE_NULL;

    case SAI_OBJECT_TYPE_PORT:
      return lemming::dataplane::sai::OBJECT_TYPE_PORT;

    case SAI_OBJECT_TYPE_LAG:
      return lemming::dataplane::sai::OBJECT_TYPE_LAG;

    case SAI_OBJECT_TYPE_VIRTUAL_ROUTER:
      return lemming::dataplane::sai::OBJECT_TYPE_VIRTUAL_ROUTER;

    case SAI_OBJECT_TYPE_NEXT_HOP:
      return lemming::dataplane::sai::OBJECT_TYPE_NEXT_HOP;

    case SAI_OBJECT_TYPE_NEXT_HOP_GROUP:
      return lemming::dataplane::sai::OBJECT_TYPE_NEXT_HOP_GROUP;

    case SAI_OBJECT_TYPE_ROUTER_INTERFACE:
      return lemming::dataplane::sai::OBJECT_TYPE_ROUTER_INTERFACE;

    case SAI_OBJECT_TYPE_ACL_TABLE:
      return lemming::dataplane::sai::OBJECT_TYPE_ACL_TABLE;

    case SAI_OBJECT_TYPE_ACL_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_ACL_ENTRY;

    case SAI_OBJECT_TYPE_ACL_COUNTER:
      return lemming::dataplane::sai::OBJECT_TYPE_ACL_COUNTER;

    case SAI_OBJECT_TYPE_ACL_RANGE:
      return lemming::dataplane::sai::OBJECT_TYPE_ACL_RANGE;

    case SAI_OBJECT_TYPE_ACL_TABLE_GROUP:
      return lemming::dataplane::sai::OBJECT_TYPE_ACL_TABLE_GROUP;

    case SAI_OBJECT_TYPE_ACL_TABLE_GROUP_MEMBER:
      return lemming::dataplane::sai::OBJECT_TYPE_ACL_TABLE_GROUP_MEMBER;

    case SAI_OBJECT_TYPE_HOSTIF:
      return lemming::dataplane::sai::OBJECT_TYPE_HOSTIF;

    case SAI_OBJECT_TYPE_MIRROR_SESSION:
      return lemming::dataplane::sai::OBJECT_TYPE_MIRROR_SESSION;

    case SAI_OBJECT_TYPE_SAMPLEPACKET:
      return lemming::dataplane::sai::OBJECT_TYPE_SAMPLEPACKET;

    case SAI_OBJECT_TYPE_STP:
      return lemming::dataplane::sai::OBJECT_TYPE_STP;

    case SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP:
      return lemming::dataplane::sai::OBJECT_TYPE_HOSTIF_TRAP_GROUP;

    case SAI_OBJECT_TYPE_POLICER:
      return lemming::dataplane::sai::OBJECT_TYPE_POLICER;

    case SAI_OBJECT_TYPE_WRED:
      return lemming::dataplane::sai::OBJECT_TYPE_WRED;

    case SAI_OBJECT_TYPE_QOS_MAP:
      return lemming::dataplane::sai::OBJECT_TYPE_QOS_MAP;

    case SAI_OBJECT_TYPE_QUEUE:
      return lemming::dataplane::sai::OBJECT_TYPE_QUEUE;

    case SAI_OBJECT_TYPE_SCHEDULER:
      return lemming::dataplane::sai::OBJECT_TYPE_SCHEDULER;

    case SAI_OBJECT_TYPE_SCHEDULER_GROUP:
      return lemming::dataplane::sai::OBJECT_TYPE_SCHEDULER_GROUP;

    case SAI_OBJECT_TYPE_BUFFER_POOL:
      return lemming::dataplane::sai::OBJECT_TYPE_BUFFER_POOL;

    case SAI_OBJECT_TYPE_BUFFER_PROFILE:
      return lemming::dataplane::sai::OBJECT_TYPE_BUFFER_PROFILE;

    case SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP:
      return lemming::dataplane::sai::OBJECT_TYPE_INGRESS_PRIORITY_GROUP;

    case SAI_OBJECT_TYPE_LAG_MEMBER:
      return lemming::dataplane::sai::OBJECT_TYPE_LAG_MEMBER;

    case SAI_OBJECT_TYPE_HASH:
      return lemming::dataplane::sai::OBJECT_TYPE_HASH;

    case SAI_OBJECT_TYPE_UDF:
      return lemming::dataplane::sai::OBJECT_TYPE_UDF;

    case SAI_OBJECT_TYPE_UDF_MATCH:
      return lemming::dataplane::sai::OBJECT_TYPE_UDF_MATCH;

    case SAI_OBJECT_TYPE_UDF_GROUP:
      return lemming::dataplane::sai::OBJECT_TYPE_UDF_GROUP;

    case SAI_OBJECT_TYPE_FDB_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_FDB_ENTRY;

    case SAI_OBJECT_TYPE_SWITCH:
      return lemming::dataplane::sai::OBJECT_TYPE_SWITCH;

    case SAI_OBJECT_TYPE_HOSTIF_TRAP:
      return lemming::dataplane::sai::OBJECT_TYPE_HOSTIF_TRAP;

    case SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_HOSTIF_TABLE_ENTRY;

    case SAI_OBJECT_TYPE_NEIGHBOR_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_NEIGHBOR_ENTRY;

    case SAI_OBJECT_TYPE_ROUTE_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_ROUTE_ENTRY;

    case SAI_OBJECT_TYPE_VLAN:
      return lemming::dataplane::sai::OBJECT_TYPE_VLAN;

    case SAI_OBJECT_TYPE_VLAN_MEMBER:
      return lemming::dataplane::sai::OBJECT_TYPE_VLAN_MEMBER;

    case SAI_OBJECT_TYPE_HOSTIF_PACKET:
      return lemming::dataplane::sai::OBJECT_TYPE_HOSTIF_PACKET;

    case SAI_OBJECT_TYPE_TUNNEL_MAP:
      return lemming::dataplane::sai::OBJECT_TYPE_TUNNEL_MAP;

    case SAI_OBJECT_TYPE_TUNNEL:
      return lemming::dataplane::sai::OBJECT_TYPE_TUNNEL;

    case SAI_OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY;

    case SAI_OBJECT_TYPE_FDB_FLUSH:
      return lemming::dataplane::sai::OBJECT_TYPE_FDB_FLUSH;

    case SAI_OBJECT_TYPE_NEXT_HOP_GROUP_MEMBER:
      return lemming::dataplane::sai::OBJECT_TYPE_NEXT_HOP_GROUP_MEMBER;

    case SAI_OBJECT_TYPE_STP_PORT:
      return lemming::dataplane::sai::OBJECT_TYPE_STP_PORT;

    case SAI_OBJECT_TYPE_RPF_GROUP:
      return lemming::dataplane::sai::OBJECT_TYPE_RPF_GROUP;

    case SAI_OBJECT_TYPE_RPF_GROUP_MEMBER:
      return lemming::dataplane::sai::OBJECT_TYPE_RPF_GROUP_MEMBER;

    case SAI_OBJECT_TYPE_L2MC_GROUP:
      return lemming::dataplane::sai::OBJECT_TYPE_L2MC_GROUP;

    case SAI_OBJECT_TYPE_L2MC_GROUP_MEMBER:
      return lemming::dataplane::sai::OBJECT_TYPE_L2MC_GROUP_MEMBER;

    case SAI_OBJECT_TYPE_IPMC_GROUP:
      return lemming::dataplane::sai::OBJECT_TYPE_IPMC_GROUP;

    case SAI_OBJECT_TYPE_IPMC_GROUP_MEMBER:
      return lemming::dataplane::sai::OBJECT_TYPE_IPMC_GROUP_MEMBER;

    case SAI_OBJECT_TYPE_L2MC_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_L2MC_ENTRY;

    case SAI_OBJECT_TYPE_IPMC_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_IPMC_ENTRY;

    case SAI_OBJECT_TYPE_MCAST_FDB_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_MCAST_FDB_ENTRY;

    case SAI_OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP:
      return lemming::dataplane::sai::OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP;

    case SAI_OBJECT_TYPE_BRIDGE:
      return lemming::dataplane::sai::OBJECT_TYPE_BRIDGE;

    case SAI_OBJECT_TYPE_BRIDGE_PORT:
      return lemming::dataplane::sai::OBJECT_TYPE_BRIDGE_PORT;

    case SAI_OBJECT_TYPE_TUNNEL_MAP_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_TUNNEL_MAP_ENTRY;

    case SAI_OBJECT_TYPE_TAM:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM;

    case SAI_OBJECT_TYPE_SRV6_SIDLIST:
      return lemming::dataplane::sai::OBJECT_TYPE_SRV6_SIDLIST;

    case SAI_OBJECT_TYPE_PORT_POOL:
      return lemming::dataplane::sai::OBJECT_TYPE_PORT_POOL;

    case SAI_OBJECT_TYPE_INSEG_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_INSEG_ENTRY;

    case SAI_OBJECT_TYPE_DTEL:
      return lemming::dataplane::sai::OBJECT_TYPE_DTEL;

    case SAI_OBJECT_TYPE_DTEL_QUEUE_REPORT:
      return lemming::dataplane::sai::OBJECT_TYPE_DTEL_QUEUE_REPORT;

    case SAI_OBJECT_TYPE_DTEL_INT_SESSION:
      return lemming::dataplane::sai::OBJECT_TYPE_DTEL_INT_SESSION;

    case SAI_OBJECT_TYPE_DTEL_REPORT_SESSION:
      return lemming::dataplane::sai::OBJECT_TYPE_DTEL_REPORT_SESSION;

    case SAI_OBJECT_TYPE_DTEL_EVENT:
      return lemming::dataplane::sai::OBJECT_TYPE_DTEL_EVENT;

    case SAI_OBJECT_TYPE_BFD_SESSION:
      return lemming::dataplane::sai::OBJECT_TYPE_BFD_SESSION;

    case SAI_OBJECT_TYPE_ISOLATION_GROUP:
      return lemming::dataplane::sai::OBJECT_TYPE_ISOLATION_GROUP;

    case SAI_OBJECT_TYPE_ISOLATION_GROUP_MEMBER:
      return lemming::dataplane::sai::OBJECT_TYPE_ISOLATION_GROUP_MEMBER;

    case SAI_OBJECT_TYPE_TAM_MATH_FUNC:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM_MATH_FUNC;

    case SAI_OBJECT_TYPE_TAM_REPORT:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM_REPORT;

    case SAI_OBJECT_TYPE_TAM_EVENT_THRESHOLD:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM_EVENT_THRESHOLD;

    case SAI_OBJECT_TYPE_TAM_TEL_TYPE:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM_TEL_TYPE;

    case SAI_OBJECT_TYPE_TAM_TRANSPORT:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM_TRANSPORT;

    case SAI_OBJECT_TYPE_TAM_TELEMETRY:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM_TELEMETRY;

    case SAI_OBJECT_TYPE_TAM_COLLECTOR:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM_COLLECTOR;

    case SAI_OBJECT_TYPE_TAM_EVENT_ACTION:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM_EVENT_ACTION;

    case SAI_OBJECT_TYPE_TAM_EVENT:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM_EVENT;

    case SAI_OBJECT_TYPE_NAT_ZONE_COUNTER:
      return lemming::dataplane::sai::OBJECT_TYPE_NAT_ZONE_COUNTER;

    case SAI_OBJECT_TYPE_NAT_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_NAT_ENTRY;

    case SAI_OBJECT_TYPE_TAM_INT:
      return lemming::dataplane::sai::OBJECT_TYPE_TAM_INT;

    case SAI_OBJECT_TYPE_COUNTER:
      return lemming::dataplane::sai::OBJECT_TYPE_COUNTER;

    case SAI_OBJECT_TYPE_DEBUG_COUNTER:
      return lemming::dataplane::sai::OBJECT_TYPE_DEBUG_COUNTER;

    case SAI_OBJECT_TYPE_PORT_CONNECTOR:
      return lemming::dataplane::sai::OBJECT_TYPE_PORT_CONNECTOR;

    case SAI_OBJECT_TYPE_PORT_SERDES:
      return lemming::dataplane::sai::OBJECT_TYPE_PORT_SERDES;

    case SAI_OBJECT_TYPE_MACSEC:
      return lemming::dataplane::sai::OBJECT_TYPE_MACSEC;

    case SAI_OBJECT_TYPE_MACSEC_PORT:
      return lemming::dataplane::sai::OBJECT_TYPE_MACSEC_PORT;

    case SAI_OBJECT_TYPE_MACSEC_FLOW:
      return lemming::dataplane::sai::OBJECT_TYPE_MACSEC_FLOW;

    case SAI_OBJECT_TYPE_MACSEC_SC:
      return lemming::dataplane::sai::OBJECT_TYPE_MACSEC_SC;

    case SAI_OBJECT_TYPE_MACSEC_SA:
      return lemming::dataplane::sai::OBJECT_TYPE_MACSEC_SA;

    case SAI_OBJECT_TYPE_SYSTEM_PORT:
      return lemming::dataplane::sai::OBJECT_TYPE_SYSTEM_PORT;

    case SAI_OBJECT_TYPE_FINE_GRAINED_HASH_FIELD:
      return lemming::dataplane::sai::OBJECT_TYPE_FINE_GRAINED_HASH_FIELD;

    case SAI_OBJECT_TYPE_SWITCH_TUNNEL:
      return lemming::dataplane::sai::OBJECT_TYPE_SWITCH_TUNNEL;

    case SAI_OBJECT_TYPE_MY_SID_ENTRY:
      return lemming::dataplane::sai::OBJECT_TYPE_MY_SID_ENTRY;

    case SAI_OBJECT_TYPE_MY_MAC:
      return lemming::dataplane::sai::OBJECT_TYPE_MY_MAC;

    case SAI_OBJECT_TYPE_NEXT_HOP_GROUP_MAP:
      return lemming::dataplane::sai::OBJECT_TYPE_NEXT_HOP_GROUP_MAP;

    case SAI_OBJECT_TYPE_IPSEC:
      return lemming::dataplane::sai::OBJECT_TYPE_IPSEC;

    case SAI_OBJECT_TYPE_IPSEC_PORT:
      return lemming::dataplane::sai::OBJECT_TYPE_IPSEC_PORT;

    case SAI_OBJECT_TYPE_IPSEC_SA:
      return lemming::dataplane::sai::OBJECT_TYPE_IPSEC_SA;

    case SAI_OBJECT_TYPE_GENERIC_PROGRAMMABLE:
      return lemming::dataplane::sai::OBJECT_TYPE_GENERIC_PROGRAMMABLE;

    case SAI_OBJECT_TYPE_MAX:
      return lemming::dataplane::sai::OBJECT_TYPE_MAX;

    default:
      return lemming::dataplane::sai::OBJECT_TYPE_UNSPECIFIED;
  }
}
sai_object_type_t convert_sai_object_type_t_to_sai(
    lemming::dataplane::sai::ObjectType val) {
  switch (val) {
    case lemming::dataplane::sai::OBJECT_TYPE_NULL:
      return SAI_OBJECT_TYPE_NULL;

    case lemming::dataplane::sai::OBJECT_TYPE_PORT:
      return SAI_OBJECT_TYPE_PORT;

    case lemming::dataplane::sai::OBJECT_TYPE_LAG:
      return SAI_OBJECT_TYPE_LAG;

    case lemming::dataplane::sai::OBJECT_TYPE_VIRTUAL_ROUTER:
      return SAI_OBJECT_TYPE_VIRTUAL_ROUTER;

    case lemming::dataplane::sai::OBJECT_TYPE_NEXT_HOP:
      return SAI_OBJECT_TYPE_NEXT_HOP;

    case lemming::dataplane::sai::OBJECT_TYPE_NEXT_HOP_GROUP:
      return SAI_OBJECT_TYPE_NEXT_HOP_GROUP;

    case lemming::dataplane::sai::OBJECT_TYPE_ROUTER_INTERFACE:
      return SAI_OBJECT_TYPE_ROUTER_INTERFACE;

    case lemming::dataplane::sai::OBJECT_TYPE_ACL_TABLE:
      return SAI_OBJECT_TYPE_ACL_TABLE;

    case lemming::dataplane::sai::OBJECT_TYPE_ACL_ENTRY:
      return SAI_OBJECT_TYPE_ACL_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_ACL_COUNTER:
      return SAI_OBJECT_TYPE_ACL_COUNTER;

    case lemming::dataplane::sai::OBJECT_TYPE_ACL_RANGE:
      return SAI_OBJECT_TYPE_ACL_RANGE;

    case lemming::dataplane::sai::OBJECT_TYPE_ACL_TABLE_GROUP:
      return SAI_OBJECT_TYPE_ACL_TABLE_GROUP;

    case lemming::dataplane::sai::OBJECT_TYPE_ACL_TABLE_GROUP_MEMBER:
      return SAI_OBJECT_TYPE_ACL_TABLE_GROUP_MEMBER;

    case lemming::dataplane::sai::OBJECT_TYPE_HOSTIF:
      return SAI_OBJECT_TYPE_HOSTIF;

    case lemming::dataplane::sai::OBJECT_TYPE_MIRROR_SESSION:
      return SAI_OBJECT_TYPE_MIRROR_SESSION;

    case lemming::dataplane::sai::OBJECT_TYPE_SAMPLEPACKET:
      return SAI_OBJECT_TYPE_SAMPLEPACKET;

    case lemming::dataplane::sai::OBJECT_TYPE_STP:
      return SAI_OBJECT_TYPE_STP;

    case lemming::dataplane::sai::OBJECT_TYPE_HOSTIF_TRAP_GROUP:
      return SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP;

    case lemming::dataplane::sai::OBJECT_TYPE_POLICER:
      return SAI_OBJECT_TYPE_POLICER;

    case lemming::dataplane::sai::OBJECT_TYPE_WRED:
      return SAI_OBJECT_TYPE_WRED;

    case lemming::dataplane::sai::OBJECT_TYPE_QOS_MAP:
      return SAI_OBJECT_TYPE_QOS_MAP;

    case lemming::dataplane::sai::OBJECT_TYPE_QUEUE:
      return SAI_OBJECT_TYPE_QUEUE;

    case lemming::dataplane::sai::OBJECT_TYPE_SCHEDULER:
      return SAI_OBJECT_TYPE_SCHEDULER;

    case lemming::dataplane::sai::OBJECT_TYPE_SCHEDULER_GROUP:
      return SAI_OBJECT_TYPE_SCHEDULER_GROUP;

    case lemming::dataplane::sai::OBJECT_TYPE_BUFFER_POOL:
      return SAI_OBJECT_TYPE_BUFFER_POOL;

    case lemming::dataplane::sai::OBJECT_TYPE_BUFFER_PROFILE:
      return SAI_OBJECT_TYPE_BUFFER_PROFILE;

    case lemming::dataplane::sai::OBJECT_TYPE_INGRESS_PRIORITY_GROUP:
      return SAI_OBJECT_TYPE_INGRESS_PRIORITY_GROUP;

    case lemming::dataplane::sai::OBJECT_TYPE_LAG_MEMBER:
      return SAI_OBJECT_TYPE_LAG_MEMBER;

    case lemming::dataplane::sai::OBJECT_TYPE_HASH:
      return SAI_OBJECT_TYPE_HASH;

    case lemming::dataplane::sai::OBJECT_TYPE_UDF:
      return SAI_OBJECT_TYPE_UDF;

    case lemming::dataplane::sai::OBJECT_TYPE_UDF_MATCH:
      return SAI_OBJECT_TYPE_UDF_MATCH;

    case lemming::dataplane::sai::OBJECT_TYPE_UDF_GROUP:
      return SAI_OBJECT_TYPE_UDF_GROUP;

    case lemming::dataplane::sai::OBJECT_TYPE_FDB_ENTRY:
      return SAI_OBJECT_TYPE_FDB_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_SWITCH:
      return SAI_OBJECT_TYPE_SWITCH;

    case lemming::dataplane::sai::OBJECT_TYPE_HOSTIF_TRAP:
      return SAI_OBJECT_TYPE_HOSTIF_TRAP;

    case lemming::dataplane::sai::OBJECT_TYPE_HOSTIF_TABLE_ENTRY:
      return SAI_OBJECT_TYPE_HOSTIF_TABLE_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_NEIGHBOR_ENTRY:
      return SAI_OBJECT_TYPE_NEIGHBOR_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_ROUTE_ENTRY:
      return SAI_OBJECT_TYPE_ROUTE_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_VLAN:
      return SAI_OBJECT_TYPE_VLAN;

    case lemming::dataplane::sai::OBJECT_TYPE_VLAN_MEMBER:
      return SAI_OBJECT_TYPE_VLAN_MEMBER;

    case lemming::dataplane::sai::OBJECT_TYPE_HOSTIF_PACKET:
      return SAI_OBJECT_TYPE_HOSTIF_PACKET;

    case lemming::dataplane::sai::OBJECT_TYPE_TUNNEL_MAP:
      return SAI_OBJECT_TYPE_TUNNEL_MAP;

    case lemming::dataplane::sai::OBJECT_TYPE_TUNNEL:
      return SAI_OBJECT_TYPE_TUNNEL;

    case lemming::dataplane::sai::OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY:
      return SAI_OBJECT_TYPE_TUNNEL_TERM_TABLE_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_FDB_FLUSH:
      return SAI_OBJECT_TYPE_FDB_FLUSH;

    case lemming::dataplane::sai::OBJECT_TYPE_NEXT_HOP_GROUP_MEMBER:
      return SAI_OBJECT_TYPE_NEXT_HOP_GROUP_MEMBER;

    case lemming::dataplane::sai::OBJECT_TYPE_STP_PORT:
      return SAI_OBJECT_TYPE_STP_PORT;

    case lemming::dataplane::sai::OBJECT_TYPE_RPF_GROUP:
      return SAI_OBJECT_TYPE_RPF_GROUP;

    case lemming::dataplane::sai::OBJECT_TYPE_RPF_GROUP_MEMBER:
      return SAI_OBJECT_TYPE_RPF_GROUP_MEMBER;

    case lemming::dataplane::sai::OBJECT_TYPE_L2MC_GROUP:
      return SAI_OBJECT_TYPE_L2MC_GROUP;

    case lemming::dataplane::sai::OBJECT_TYPE_L2MC_GROUP_MEMBER:
      return SAI_OBJECT_TYPE_L2MC_GROUP_MEMBER;

    case lemming::dataplane::sai::OBJECT_TYPE_IPMC_GROUP:
      return SAI_OBJECT_TYPE_IPMC_GROUP;

    case lemming::dataplane::sai::OBJECT_TYPE_IPMC_GROUP_MEMBER:
      return SAI_OBJECT_TYPE_IPMC_GROUP_MEMBER;

    case lemming::dataplane::sai::OBJECT_TYPE_L2MC_ENTRY:
      return SAI_OBJECT_TYPE_L2MC_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_IPMC_ENTRY:
      return SAI_OBJECT_TYPE_IPMC_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_MCAST_FDB_ENTRY:
      return SAI_OBJECT_TYPE_MCAST_FDB_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP:
      return SAI_OBJECT_TYPE_HOSTIF_USER_DEFINED_TRAP;

    case lemming::dataplane::sai::OBJECT_TYPE_BRIDGE:
      return SAI_OBJECT_TYPE_BRIDGE;

    case lemming::dataplane::sai::OBJECT_TYPE_BRIDGE_PORT:
      return SAI_OBJECT_TYPE_BRIDGE_PORT;

    case lemming::dataplane::sai::OBJECT_TYPE_TUNNEL_MAP_ENTRY:
      return SAI_OBJECT_TYPE_TUNNEL_MAP_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM:
      return SAI_OBJECT_TYPE_TAM;

    case lemming::dataplane::sai::OBJECT_TYPE_SRV6_SIDLIST:
      return SAI_OBJECT_TYPE_SRV6_SIDLIST;

    case lemming::dataplane::sai::OBJECT_TYPE_PORT_POOL:
      return SAI_OBJECT_TYPE_PORT_POOL;

    case lemming::dataplane::sai::OBJECT_TYPE_INSEG_ENTRY:
      return SAI_OBJECT_TYPE_INSEG_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_DTEL:
      return SAI_OBJECT_TYPE_DTEL;

    case lemming::dataplane::sai::OBJECT_TYPE_DTEL_QUEUE_REPORT:
      return SAI_OBJECT_TYPE_DTEL_QUEUE_REPORT;

    case lemming::dataplane::sai::OBJECT_TYPE_DTEL_INT_SESSION:
      return SAI_OBJECT_TYPE_DTEL_INT_SESSION;

    case lemming::dataplane::sai::OBJECT_TYPE_DTEL_REPORT_SESSION:
      return SAI_OBJECT_TYPE_DTEL_REPORT_SESSION;

    case lemming::dataplane::sai::OBJECT_TYPE_DTEL_EVENT:
      return SAI_OBJECT_TYPE_DTEL_EVENT;

    case lemming::dataplane::sai::OBJECT_TYPE_BFD_SESSION:
      return SAI_OBJECT_TYPE_BFD_SESSION;

    case lemming::dataplane::sai::OBJECT_TYPE_ISOLATION_GROUP:
      return SAI_OBJECT_TYPE_ISOLATION_GROUP;

    case lemming::dataplane::sai::OBJECT_TYPE_ISOLATION_GROUP_MEMBER:
      return SAI_OBJECT_TYPE_ISOLATION_GROUP_MEMBER;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM_MATH_FUNC:
      return SAI_OBJECT_TYPE_TAM_MATH_FUNC;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM_REPORT:
      return SAI_OBJECT_TYPE_TAM_REPORT;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM_EVENT_THRESHOLD:
      return SAI_OBJECT_TYPE_TAM_EVENT_THRESHOLD;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM_TEL_TYPE:
      return SAI_OBJECT_TYPE_TAM_TEL_TYPE;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM_TRANSPORT:
      return SAI_OBJECT_TYPE_TAM_TRANSPORT;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM_TELEMETRY:
      return SAI_OBJECT_TYPE_TAM_TELEMETRY;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM_COLLECTOR:
      return SAI_OBJECT_TYPE_TAM_COLLECTOR;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM_EVENT_ACTION:
      return SAI_OBJECT_TYPE_TAM_EVENT_ACTION;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM_EVENT:
      return SAI_OBJECT_TYPE_TAM_EVENT;

    case lemming::dataplane::sai::OBJECT_TYPE_NAT_ZONE_COUNTER:
      return SAI_OBJECT_TYPE_NAT_ZONE_COUNTER;

    case lemming::dataplane::sai::OBJECT_TYPE_NAT_ENTRY:
      return SAI_OBJECT_TYPE_NAT_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_TAM_INT:
      return SAI_OBJECT_TYPE_TAM_INT;

    case lemming::dataplane::sai::OBJECT_TYPE_COUNTER:
      return SAI_OBJECT_TYPE_COUNTER;

    case lemming::dataplane::sai::OBJECT_TYPE_DEBUG_COUNTER:
      return SAI_OBJECT_TYPE_DEBUG_COUNTER;

    case lemming::dataplane::sai::OBJECT_TYPE_PORT_CONNECTOR:
      return SAI_OBJECT_TYPE_PORT_CONNECTOR;

    case lemming::dataplane::sai::OBJECT_TYPE_PORT_SERDES:
      return SAI_OBJECT_TYPE_PORT_SERDES;

    case lemming::dataplane::sai::OBJECT_TYPE_MACSEC:
      return SAI_OBJECT_TYPE_MACSEC;

    case lemming::dataplane::sai::OBJECT_TYPE_MACSEC_PORT:
      return SAI_OBJECT_TYPE_MACSEC_PORT;

    case lemming::dataplane::sai::OBJECT_TYPE_MACSEC_FLOW:
      return SAI_OBJECT_TYPE_MACSEC_FLOW;

    case lemming::dataplane::sai::OBJECT_TYPE_MACSEC_SC:
      return SAI_OBJECT_TYPE_MACSEC_SC;

    case lemming::dataplane::sai::OBJECT_TYPE_MACSEC_SA:
      return SAI_OBJECT_TYPE_MACSEC_SA;

    case lemming::dataplane::sai::OBJECT_TYPE_SYSTEM_PORT:
      return SAI_OBJECT_TYPE_SYSTEM_PORT;

    case lemming::dataplane::sai::OBJECT_TYPE_FINE_GRAINED_HASH_FIELD:
      return SAI_OBJECT_TYPE_FINE_GRAINED_HASH_FIELD;

    case lemming::dataplane::sai::OBJECT_TYPE_SWITCH_TUNNEL:
      return SAI_OBJECT_TYPE_SWITCH_TUNNEL;

    case lemming::dataplane::sai::OBJECT_TYPE_MY_SID_ENTRY:
      return SAI_OBJECT_TYPE_MY_SID_ENTRY;

    case lemming::dataplane::sai::OBJECT_TYPE_MY_MAC:
      return SAI_OBJECT_TYPE_MY_MAC;

    case lemming::dataplane::sai::OBJECT_TYPE_NEXT_HOP_GROUP_MAP:
      return SAI_OBJECT_TYPE_NEXT_HOP_GROUP_MAP;

    case lemming::dataplane::sai::OBJECT_TYPE_IPSEC:
      return SAI_OBJECT_TYPE_IPSEC;

    case lemming::dataplane::sai::OBJECT_TYPE_IPSEC_PORT:
      return SAI_OBJECT_TYPE_IPSEC_PORT;

    case lemming::dataplane::sai::OBJECT_TYPE_IPSEC_SA:
      return SAI_OBJECT_TYPE_IPSEC_SA;

    case lemming::dataplane::sai::OBJECT_TYPE_GENERIC_PROGRAMMABLE:
      return SAI_OBJECT_TYPE_GENERIC_PROGRAMMABLE;

    case lemming::dataplane::sai::OBJECT_TYPE_MAX:
      return SAI_OBJECT_TYPE_MAX;

    default:
      return SAI_OBJECT_TYPE_NULL;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_object_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_object_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_object_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_object_type_t_to_sai(
        static_cast<lemming::dataplane::sai::ObjectType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::OutDropReason convert_sai_out_drop_reason_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_OUT_DROP_REASON_START:
      return lemming::dataplane::sai::OUT_DROP_REASON_START;

    case SAI_OUT_DROP_REASON_EGRESS_VLAN_FILTER:
      return lemming::dataplane::sai::OUT_DROP_REASON_EGRESS_VLAN_FILTER;

    case SAI_OUT_DROP_REASON_L3_ANY:
      return lemming::dataplane::sai::OUT_DROP_REASON_L3_ANY;

    case SAI_OUT_DROP_REASON_L3_EGRESS_LINK_DOWN:
      return lemming::dataplane::sai::OUT_DROP_REASON_L3_EGRESS_LINK_DOWN;

    case SAI_OUT_DROP_REASON_TUNNEL_LOOPBACK_PACKET_DROP:
      return lemming::dataplane::sai::
          OUT_DROP_REASON_TUNNEL_LOOPBACK_PACKET_DROP;

    case SAI_OUT_DROP_REASON_END:
      return lemming::dataplane::sai::OUT_DROP_REASON_END;

    case SAI_OUT_DROP_REASON_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::OUT_DROP_REASON_CUSTOM_RANGE_BASE;

    case SAI_OUT_DROP_REASON_CUSTOM_RANGE_END:
      return lemming::dataplane::sai::OUT_DROP_REASON_CUSTOM_RANGE_END;

    default:
      return lemming::dataplane::sai::OUT_DROP_REASON_UNSPECIFIED;
  }
}
sai_out_drop_reason_t convert_sai_out_drop_reason_t_to_sai(
    lemming::dataplane::sai::OutDropReason val) {
  switch (val) {
    case lemming::dataplane::sai::OUT_DROP_REASON_START:
      return SAI_OUT_DROP_REASON_START;

    case lemming::dataplane::sai::OUT_DROP_REASON_EGRESS_VLAN_FILTER:
      return SAI_OUT_DROP_REASON_EGRESS_VLAN_FILTER;

    case lemming::dataplane::sai::OUT_DROP_REASON_L3_ANY:
      return SAI_OUT_DROP_REASON_L3_ANY;

    case lemming::dataplane::sai::OUT_DROP_REASON_L3_EGRESS_LINK_DOWN:
      return SAI_OUT_DROP_REASON_L3_EGRESS_LINK_DOWN;

    case lemming::dataplane::sai::OUT_DROP_REASON_TUNNEL_LOOPBACK_PACKET_DROP:
      return SAI_OUT_DROP_REASON_TUNNEL_LOOPBACK_PACKET_DROP;

    case lemming::dataplane::sai::OUT_DROP_REASON_END:
      return SAI_OUT_DROP_REASON_END;

    case lemming::dataplane::sai::OUT_DROP_REASON_CUSTOM_RANGE_BASE:
      return SAI_OUT_DROP_REASON_CUSTOM_RANGE_BASE;

    case lemming::dataplane::sai::OUT_DROP_REASON_CUSTOM_RANGE_END:
      return SAI_OUT_DROP_REASON_CUSTOM_RANGE_END;

    default:
      return SAI_OUT_DROP_REASON_START;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_out_drop_reason_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_out_drop_reason_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_out_drop_reason_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_out_drop_reason_t_to_sai(
        static_cast<lemming::dataplane::sai::OutDropReason>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::OutsegExpMode convert_sai_outseg_exp_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_OUTSEG_EXP_MODE_UNIFORM:
      return lemming::dataplane::sai::OUTSEG_EXP_MODE_UNIFORM;

    case SAI_OUTSEG_EXP_MODE_PIPE:
      return lemming::dataplane::sai::OUTSEG_EXP_MODE_PIPE;

    default:
      return lemming::dataplane::sai::OUTSEG_EXP_MODE_UNSPECIFIED;
  }
}
sai_outseg_exp_mode_t convert_sai_outseg_exp_mode_t_to_sai(
    lemming::dataplane::sai::OutsegExpMode val) {
  switch (val) {
    case lemming::dataplane::sai::OUTSEG_EXP_MODE_UNIFORM:
      return SAI_OUTSEG_EXP_MODE_UNIFORM;

    case lemming::dataplane::sai::OUTSEG_EXP_MODE_PIPE:
      return SAI_OUTSEG_EXP_MODE_PIPE;

    default:
      return SAI_OUTSEG_EXP_MODE_UNIFORM;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_outseg_exp_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_outseg_exp_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_outseg_exp_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_outseg_exp_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::OutsegExpMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::OutsegTtlMode convert_sai_outseg_ttl_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_OUTSEG_TTL_MODE_UNIFORM:
      return lemming::dataplane::sai::OUTSEG_TTL_MODE_UNIFORM;

    case SAI_OUTSEG_TTL_MODE_PIPE:
      return lemming::dataplane::sai::OUTSEG_TTL_MODE_PIPE;

    default:
      return lemming::dataplane::sai::OUTSEG_TTL_MODE_UNSPECIFIED;
  }
}
sai_outseg_ttl_mode_t convert_sai_outseg_ttl_mode_t_to_sai(
    lemming::dataplane::sai::OutsegTtlMode val) {
  switch (val) {
    case lemming::dataplane::sai::OUTSEG_TTL_MODE_UNIFORM:
      return SAI_OUTSEG_TTL_MODE_UNIFORM;

    case lemming::dataplane::sai::OUTSEG_TTL_MODE_PIPE:
      return SAI_OUTSEG_TTL_MODE_PIPE;

    default:
      return SAI_OUTSEG_TTL_MODE_UNIFORM;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_outseg_ttl_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_outseg_ttl_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_outseg_ttl_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_outseg_ttl_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::OutsegTtlMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::OutsegType convert_sai_outseg_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_OUTSEG_TYPE_PUSH:
      return lemming::dataplane::sai::OUTSEG_TYPE_PUSH;

    case SAI_OUTSEG_TYPE_SWAP:
      return lemming::dataplane::sai::OUTSEG_TYPE_SWAP;

    default:
      return lemming::dataplane::sai::OUTSEG_TYPE_UNSPECIFIED;
  }
}
sai_outseg_type_t convert_sai_outseg_type_t_to_sai(
    lemming::dataplane::sai::OutsegType val) {
  switch (val) {
    case lemming::dataplane::sai::OUTSEG_TYPE_PUSH:
      return SAI_OUTSEG_TYPE_PUSH;

    case lemming::dataplane::sai::OUTSEG_TYPE_SWAP:
      return SAI_OUTSEG_TYPE_SWAP;

    default:
      return SAI_OUTSEG_TYPE_PUSH;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_outseg_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_outseg_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_outseg_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_outseg_type_t_to_sai(
        static_cast<lemming::dataplane::sai::OutsegType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PacketAction convert_sai_packet_action_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PACKET_ACTION_DROP:
      return lemming::dataplane::sai::PACKET_ACTION_DROP;

    case SAI_PACKET_ACTION_FORWARD:
      return lemming::dataplane::sai::PACKET_ACTION_FORWARD;

    case SAI_PACKET_ACTION_COPY:
      return lemming::dataplane::sai::PACKET_ACTION_COPY;

    case SAI_PACKET_ACTION_COPY_CANCEL:
      return lemming::dataplane::sai::PACKET_ACTION_COPY_CANCEL;

    case SAI_PACKET_ACTION_TRAP:
      return lemming::dataplane::sai::PACKET_ACTION_TRAP;

    case SAI_PACKET_ACTION_LOG:
      return lemming::dataplane::sai::PACKET_ACTION_LOG;

    case SAI_PACKET_ACTION_DENY:
      return lemming::dataplane::sai::PACKET_ACTION_DENY;

    case SAI_PACKET_ACTION_TRANSIT:
      return lemming::dataplane::sai::PACKET_ACTION_TRANSIT;

    case SAI_PACKET_ACTION_DONOTDROP:
      return lemming::dataplane::sai::PACKET_ACTION_DONOTDROP;

    default:
      return lemming::dataplane::sai::PACKET_ACTION_UNSPECIFIED;
  }
}
sai_packet_action_t convert_sai_packet_action_t_to_sai(
    lemming::dataplane::sai::PacketAction val) {
  switch (val) {
    case lemming::dataplane::sai::PACKET_ACTION_DROP:
      return SAI_PACKET_ACTION_DROP;

    case lemming::dataplane::sai::PACKET_ACTION_FORWARD:
      return SAI_PACKET_ACTION_FORWARD;

    case lemming::dataplane::sai::PACKET_ACTION_COPY:
      return SAI_PACKET_ACTION_COPY;

    case lemming::dataplane::sai::PACKET_ACTION_COPY_CANCEL:
      return SAI_PACKET_ACTION_COPY_CANCEL;

    case lemming::dataplane::sai::PACKET_ACTION_TRAP:
      return SAI_PACKET_ACTION_TRAP;

    case lemming::dataplane::sai::PACKET_ACTION_LOG:
      return SAI_PACKET_ACTION_LOG;

    case lemming::dataplane::sai::PACKET_ACTION_DENY:
      return SAI_PACKET_ACTION_DENY;

    case lemming::dataplane::sai::PACKET_ACTION_TRANSIT:
      return SAI_PACKET_ACTION_TRANSIT;

    case lemming::dataplane::sai::PACKET_ACTION_DONOTDROP:
      return SAI_PACKET_ACTION_DONOTDROP;

    default:
      return SAI_PACKET_ACTION_DROP;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_packet_action_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_packet_action_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_packet_action_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_packet_action_t_to_sai(
        static_cast<lemming::dataplane::sai::PacketAction>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PacketColor convert_sai_packet_color_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PACKET_COLOR_GREEN:
      return lemming::dataplane::sai::PACKET_COLOR_GREEN;

    case SAI_PACKET_COLOR_YELLOW:
      return lemming::dataplane::sai::PACKET_COLOR_YELLOW;

    case SAI_PACKET_COLOR_RED:
      return lemming::dataplane::sai::PACKET_COLOR_RED;

    default:
      return lemming::dataplane::sai::PACKET_COLOR_UNSPECIFIED;
  }
}
sai_packet_color_t convert_sai_packet_color_t_to_sai(
    lemming::dataplane::sai::PacketColor val) {
  switch (val) {
    case lemming::dataplane::sai::PACKET_COLOR_GREEN:
      return SAI_PACKET_COLOR_GREEN;

    case lemming::dataplane::sai::PACKET_COLOR_YELLOW:
      return SAI_PACKET_COLOR_YELLOW;

    case lemming::dataplane::sai::PACKET_COLOR_RED:
      return SAI_PACKET_COLOR_RED;

    default:
      return SAI_PACKET_COLOR_GREEN;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_packet_color_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_packet_color_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_packet_color_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_packet_color_t_to_sai(
        static_cast<lemming::dataplane::sai::PacketColor>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PacketVlan convert_sai_packet_vlan_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PACKET_VLAN_UNTAG:
      return lemming::dataplane::sai::PACKET_VLAN_UNTAG;

    case SAI_PACKET_VLAN_SINGLE_OUTER_TAG:
      return lemming::dataplane::sai::PACKET_VLAN_SINGLE_OUTER_TAG;

    case SAI_PACKET_VLAN_DOUBLE_TAG:
      return lemming::dataplane::sai::PACKET_VLAN_DOUBLE_TAG;

    default:
      return lemming::dataplane::sai::PACKET_VLAN_UNSPECIFIED;
  }
}
sai_packet_vlan_t convert_sai_packet_vlan_t_to_sai(
    lemming::dataplane::sai::PacketVlan val) {
  switch (val) {
    case lemming::dataplane::sai::PACKET_VLAN_UNTAG:
      return SAI_PACKET_VLAN_UNTAG;

    case lemming::dataplane::sai::PACKET_VLAN_SINGLE_OUTER_TAG:
      return SAI_PACKET_VLAN_SINGLE_OUTER_TAG;

    case lemming::dataplane::sai::PACKET_VLAN_DOUBLE_TAG:
      return SAI_PACKET_VLAN_DOUBLE_TAG;

    default:
      return SAI_PACKET_VLAN_UNTAG;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_packet_vlan_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_packet_vlan_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_packet_vlan_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_packet_vlan_t_to_sai(
        static_cast<lemming::dataplane::sai::PacketVlan>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PolicerAttr convert_sai_policer_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_POLICER_ATTR_METER_TYPE:
      return lemming::dataplane::sai::POLICER_ATTR_METER_TYPE;

    case SAI_POLICER_ATTR_MODE:
      return lemming::dataplane::sai::POLICER_ATTR_MODE;

    case SAI_POLICER_ATTR_COLOR_SOURCE:
      return lemming::dataplane::sai::POLICER_ATTR_COLOR_SOURCE;

    case SAI_POLICER_ATTR_CBS:
      return lemming::dataplane::sai::POLICER_ATTR_CBS;

    case SAI_POLICER_ATTR_CIR:
      return lemming::dataplane::sai::POLICER_ATTR_CIR;

    case SAI_POLICER_ATTR_PBS:
      return lemming::dataplane::sai::POLICER_ATTR_PBS;

    case SAI_POLICER_ATTR_PIR:
      return lemming::dataplane::sai::POLICER_ATTR_PIR;

    case SAI_POLICER_ATTR_GREEN_PACKET_ACTION:
      return lemming::dataplane::sai::POLICER_ATTR_GREEN_PACKET_ACTION;

    case SAI_POLICER_ATTR_YELLOW_PACKET_ACTION:
      return lemming::dataplane::sai::POLICER_ATTR_YELLOW_PACKET_ACTION;

    case SAI_POLICER_ATTR_RED_PACKET_ACTION:
      return lemming::dataplane::sai::POLICER_ATTR_RED_PACKET_ACTION;

    case SAI_POLICER_ATTR_ENABLE_COUNTER_PACKET_ACTION_LIST:
      return lemming::dataplane::sai::
          POLICER_ATTR_ENABLE_COUNTER_PACKET_ACTION_LIST;

    case SAI_POLICER_ATTR_OBJECT_STAGE:
      return lemming::dataplane::sai::POLICER_ATTR_OBJECT_STAGE;

    default:
      return lemming::dataplane::sai::POLICER_ATTR_UNSPECIFIED;
  }
}
sai_policer_attr_t convert_sai_policer_attr_t_to_sai(
    lemming::dataplane::sai::PolicerAttr val) {
  switch (val) {
    case lemming::dataplane::sai::POLICER_ATTR_METER_TYPE:
      return SAI_POLICER_ATTR_METER_TYPE;

    case lemming::dataplane::sai::POLICER_ATTR_MODE:
      return SAI_POLICER_ATTR_MODE;

    case lemming::dataplane::sai::POLICER_ATTR_COLOR_SOURCE:
      return SAI_POLICER_ATTR_COLOR_SOURCE;

    case lemming::dataplane::sai::POLICER_ATTR_CBS:
      return SAI_POLICER_ATTR_CBS;

    case lemming::dataplane::sai::POLICER_ATTR_CIR:
      return SAI_POLICER_ATTR_CIR;

    case lemming::dataplane::sai::POLICER_ATTR_PBS:
      return SAI_POLICER_ATTR_PBS;

    case lemming::dataplane::sai::POLICER_ATTR_PIR:
      return SAI_POLICER_ATTR_PIR;

    case lemming::dataplane::sai::POLICER_ATTR_GREEN_PACKET_ACTION:
      return SAI_POLICER_ATTR_GREEN_PACKET_ACTION;

    case lemming::dataplane::sai::POLICER_ATTR_YELLOW_PACKET_ACTION:
      return SAI_POLICER_ATTR_YELLOW_PACKET_ACTION;

    case lemming::dataplane::sai::POLICER_ATTR_RED_PACKET_ACTION:
      return SAI_POLICER_ATTR_RED_PACKET_ACTION;

    case lemming::dataplane::sai::
        POLICER_ATTR_ENABLE_COUNTER_PACKET_ACTION_LIST:
      return SAI_POLICER_ATTR_ENABLE_COUNTER_PACKET_ACTION_LIST;

    case lemming::dataplane::sai::POLICER_ATTR_OBJECT_STAGE:
      return SAI_POLICER_ATTR_OBJECT_STAGE;

    default:
      return SAI_POLICER_ATTR_METER_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_policer_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_policer_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_policer_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_policer_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::PolicerAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PolicerColorSource
convert_sai_policer_color_source_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_POLICER_COLOR_SOURCE_BLIND:
      return lemming::dataplane::sai::POLICER_COLOR_SOURCE_BLIND;

    case SAI_POLICER_COLOR_SOURCE_AWARE:
      return lemming::dataplane::sai::POLICER_COLOR_SOURCE_AWARE;

    case SAI_POLICER_COLOR_SOURCE_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::POLICER_COLOR_SOURCE_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::POLICER_COLOR_SOURCE_UNSPECIFIED;
  }
}
sai_policer_color_source_t convert_sai_policer_color_source_t_to_sai(
    lemming::dataplane::sai::PolicerColorSource val) {
  switch (val) {
    case lemming::dataplane::sai::POLICER_COLOR_SOURCE_BLIND:
      return SAI_POLICER_COLOR_SOURCE_BLIND;

    case lemming::dataplane::sai::POLICER_COLOR_SOURCE_AWARE:
      return SAI_POLICER_COLOR_SOURCE_AWARE;

    case lemming::dataplane::sai::POLICER_COLOR_SOURCE_CUSTOM_RANGE_BASE:
      return SAI_POLICER_COLOR_SOURCE_CUSTOM_RANGE_BASE;

    default:
      return SAI_POLICER_COLOR_SOURCE_BLIND;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_policer_color_source_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_policer_color_source_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_policer_color_source_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_policer_color_source_t_to_sai(
        static_cast<lemming::dataplane::sai::PolicerColorSource>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PolicerMode convert_sai_policer_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_POLICER_MODE_SR_TCM:
      return lemming::dataplane::sai::POLICER_MODE_SR_TCM;

    case SAI_POLICER_MODE_TR_TCM:
      return lemming::dataplane::sai::POLICER_MODE_TR_TCM;

    case SAI_POLICER_MODE_STORM_CONTROL:
      return lemming::dataplane::sai::POLICER_MODE_STORM_CONTROL;

    case SAI_POLICER_MODE_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::POLICER_MODE_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::POLICER_MODE_UNSPECIFIED;
  }
}
sai_policer_mode_t convert_sai_policer_mode_t_to_sai(
    lemming::dataplane::sai::PolicerMode val) {
  switch (val) {
    case lemming::dataplane::sai::POLICER_MODE_SR_TCM:
      return SAI_POLICER_MODE_SR_TCM;

    case lemming::dataplane::sai::POLICER_MODE_TR_TCM:
      return SAI_POLICER_MODE_TR_TCM;

    case lemming::dataplane::sai::POLICER_MODE_STORM_CONTROL:
      return SAI_POLICER_MODE_STORM_CONTROL;

    case lemming::dataplane::sai::POLICER_MODE_CUSTOM_RANGE_BASE:
      return SAI_POLICER_MODE_CUSTOM_RANGE_BASE;

    default:
      return SAI_POLICER_MODE_SR_TCM;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_policer_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_policer_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_policer_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_policer_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::PolicerMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PolicerStat convert_sai_policer_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_POLICER_STAT_PACKETS:
      return lemming::dataplane::sai::POLICER_STAT_PACKETS;

    case SAI_POLICER_STAT_ATTR_BYTES:
      return lemming::dataplane::sai::POLICER_STAT_ATTR_BYTES;

    case SAI_POLICER_STAT_GREEN_PACKETS:
      return lemming::dataplane::sai::POLICER_STAT_GREEN_PACKETS;

    case SAI_POLICER_STAT_GREEN_BYTES:
      return lemming::dataplane::sai::POLICER_STAT_GREEN_BYTES;

    case SAI_POLICER_STAT_YELLOW_PACKETS:
      return lemming::dataplane::sai::POLICER_STAT_YELLOW_PACKETS;

    case SAI_POLICER_STAT_YELLOW_BYTES:
      return lemming::dataplane::sai::POLICER_STAT_YELLOW_BYTES;

    case SAI_POLICER_STAT_RED_PACKETS:
      return lemming::dataplane::sai::POLICER_STAT_RED_PACKETS;

    case SAI_POLICER_STAT_RED_BYTES:
      return lemming::dataplane::sai::POLICER_STAT_RED_BYTES;

    case SAI_POLICER_STAT_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::POLICER_STAT_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::POLICER_STAT_UNSPECIFIED;
  }
}
sai_policer_stat_t convert_sai_policer_stat_t_to_sai(
    lemming::dataplane::sai::PolicerStat val) {
  switch (val) {
    case lemming::dataplane::sai::POLICER_STAT_PACKETS:
      return SAI_POLICER_STAT_PACKETS;

    case lemming::dataplane::sai::POLICER_STAT_ATTR_BYTES:
      return SAI_POLICER_STAT_ATTR_BYTES;

    case lemming::dataplane::sai::POLICER_STAT_GREEN_PACKETS:
      return SAI_POLICER_STAT_GREEN_PACKETS;

    case lemming::dataplane::sai::POLICER_STAT_GREEN_BYTES:
      return SAI_POLICER_STAT_GREEN_BYTES;

    case lemming::dataplane::sai::POLICER_STAT_YELLOW_PACKETS:
      return SAI_POLICER_STAT_YELLOW_PACKETS;

    case lemming::dataplane::sai::POLICER_STAT_YELLOW_BYTES:
      return SAI_POLICER_STAT_YELLOW_BYTES;

    case lemming::dataplane::sai::POLICER_STAT_RED_PACKETS:
      return SAI_POLICER_STAT_RED_PACKETS;

    case lemming::dataplane::sai::POLICER_STAT_RED_BYTES:
      return SAI_POLICER_STAT_RED_BYTES;

    case lemming::dataplane::sai::POLICER_STAT_CUSTOM_RANGE_BASE:
      return SAI_POLICER_STAT_CUSTOM_RANGE_BASE;

    default:
      return SAI_POLICER_STAT_PACKETS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_policer_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_policer_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_policer_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_policer_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::PolicerStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortAttr convert_sai_port_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_ATTR_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_TYPE;

    case SAI_PORT_ATTR_OPER_STATUS:
      return lemming::dataplane::sai::PORT_ATTR_OPER_STATUS;

    case SAI_PORT_ATTR_SUPPORTED_BREAKOUT_MODE_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_SUPPORTED_BREAKOUT_MODE_TYPE;

    case SAI_PORT_ATTR_CURRENT_BREAKOUT_MODE_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_CURRENT_BREAKOUT_MODE_TYPE;

    case SAI_PORT_ATTR_QOS_NUMBER_OF_QUEUES:
      return lemming::dataplane::sai::PORT_ATTR_QOS_NUMBER_OF_QUEUES;

    case SAI_PORT_ATTR_QOS_QUEUE_LIST:
      return lemming::dataplane::sai::PORT_ATTR_QOS_QUEUE_LIST;

    case SAI_PORT_ATTR_QOS_NUMBER_OF_SCHEDULER_GROUPS:
      return lemming::dataplane::sai::PORT_ATTR_QOS_NUMBER_OF_SCHEDULER_GROUPS;

    case SAI_PORT_ATTR_QOS_SCHEDULER_GROUP_LIST:
      return lemming::dataplane::sai::PORT_ATTR_QOS_SCHEDULER_GROUP_LIST;

    case SAI_PORT_ATTR_QOS_MAXIMUM_HEADROOM_SIZE:
      return lemming::dataplane::sai::PORT_ATTR_QOS_MAXIMUM_HEADROOM_SIZE;

    case SAI_PORT_ATTR_SUPPORTED_SPEED:
      return lemming::dataplane::sai::PORT_ATTR_SUPPORTED_SPEED;

    case SAI_PORT_ATTR_SUPPORTED_FEC_MODE:
      return lemming::dataplane::sai::PORT_ATTR_SUPPORTED_FEC_MODE;

    case SAI_PORT_ATTR_SUPPORTED_FEC_MODE_EXTENDED:
      return lemming::dataplane::sai::PORT_ATTR_SUPPORTED_FEC_MODE_EXTENDED;

    case SAI_PORT_ATTR_SUPPORTED_HALF_DUPLEX_SPEED:
      return lemming::dataplane::sai::PORT_ATTR_SUPPORTED_HALF_DUPLEX_SPEED;

    case SAI_PORT_ATTR_SUPPORTED_AUTO_NEG_MODE:
      return lemming::dataplane::sai::PORT_ATTR_SUPPORTED_AUTO_NEG_MODE;

    case SAI_PORT_ATTR_SUPPORTED_FLOW_CONTROL_MODE:
      return lemming::dataplane::sai::PORT_ATTR_SUPPORTED_FLOW_CONTROL_MODE;

    case SAI_PORT_ATTR_SUPPORTED_ASYMMETRIC_PAUSE_MODE:
      return lemming::dataplane::sai::PORT_ATTR_SUPPORTED_ASYMMETRIC_PAUSE_MODE;

    case SAI_PORT_ATTR_SUPPORTED_MEDIA_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_SUPPORTED_MEDIA_TYPE;

    case SAI_PORT_ATTR_REMOTE_ADVERTISED_SPEED:
      return lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_SPEED;

    case SAI_PORT_ATTR_REMOTE_ADVERTISED_FEC_MODE:
      return lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_FEC_MODE;

    case SAI_PORT_ATTR_REMOTE_ADVERTISED_FEC_MODE_EXTENDED:
      return lemming::dataplane::sai::
          PORT_ATTR_REMOTE_ADVERTISED_FEC_MODE_EXTENDED;

    case SAI_PORT_ATTR_REMOTE_ADVERTISED_HALF_DUPLEX_SPEED:
      return lemming::dataplane::sai::
          PORT_ATTR_REMOTE_ADVERTISED_HALF_DUPLEX_SPEED;

    case SAI_PORT_ATTR_REMOTE_ADVERTISED_AUTO_NEG_MODE:
      return lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_AUTO_NEG_MODE;

    case SAI_PORT_ATTR_REMOTE_ADVERTISED_FLOW_CONTROL_MODE:
      return lemming::dataplane::sai::
          PORT_ATTR_REMOTE_ADVERTISED_FLOW_CONTROL_MODE;

    case SAI_PORT_ATTR_REMOTE_ADVERTISED_ASYMMETRIC_PAUSE_MODE:
      return lemming::dataplane::sai::
          PORT_ATTR_REMOTE_ADVERTISED_ASYMMETRIC_PAUSE_MODE;

    case SAI_PORT_ATTR_REMOTE_ADVERTISED_MEDIA_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_MEDIA_TYPE;

    case SAI_PORT_ATTR_REMOTE_ADVERTISED_OUI_CODE:
      return lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_OUI_CODE;

    case SAI_PORT_ATTR_NUMBER_OF_INGRESS_PRIORITY_GROUPS:
      return lemming::dataplane::sai::
          PORT_ATTR_NUMBER_OF_INGRESS_PRIORITY_GROUPS;

    case SAI_PORT_ATTR_INGRESS_PRIORITY_GROUP_LIST:
      return lemming::dataplane::sai::PORT_ATTR_INGRESS_PRIORITY_GROUP_LIST;

    case SAI_PORT_ATTR_EYE_VALUES:
      return lemming::dataplane::sai::PORT_ATTR_EYE_VALUES;

    case SAI_PORT_ATTR_OPER_SPEED:
      return lemming::dataplane::sai::PORT_ATTR_OPER_SPEED;

    case SAI_PORT_ATTR_HW_LANE_LIST:
      return lemming::dataplane::sai::PORT_ATTR_HW_LANE_LIST;

    case SAI_PORT_ATTR_SPEED:
      return lemming::dataplane::sai::PORT_ATTR_SPEED;

    case SAI_PORT_ATTR_FULL_DUPLEX_MODE:
      return lemming::dataplane::sai::PORT_ATTR_FULL_DUPLEX_MODE;

    case SAI_PORT_ATTR_AUTO_NEG_MODE:
      return lemming::dataplane::sai::PORT_ATTR_AUTO_NEG_MODE;

    case SAI_PORT_ATTR_ADMIN_STATE:
      return lemming::dataplane::sai::PORT_ATTR_ADMIN_STATE;

    case SAI_PORT_ATTR_MEDIA_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_MEDIA_TYPE;

    case SAI_PORT_ATTR_ADVERTISED_SPEED:
      return lemming::dataplane::sai::PORT_ATTR_ADVERTISED_SPEED;

    case SAI_PORT_ATTR_ADVERTISED_FEC_MODE:
      return lemming::dataplane::sai::PORT_ATTR_ADVERTISED_FEC_MODE;

    case SAI_PORT_ATTR_ADVERTISED_FEC_MODE_EXTENDED:
      return lemming::dataplane::sai::PORT_ATTR_ADVERTISED_FEC_MODE_EXTENDED;

    case SAI_PORT_ATTR_ADVERTISED_HALF_DUPLEX_SPEED:
      return lemming::dataplane::sai::PORT_ATTR_ADVERTISED_HALF_DUPLEX_SPEED;

    case SAI_PORT_ATTR_ADVERTISED_AUTO_NEG_MODE:
      return lemming::dataplane::sai::PORT_ATTR_ADVERTISED_AUTO_NEG_MODE;

    case SAI_PORT_ATTR_ADVERTISED_FLOW_CONTROL_MODE:
      return lemming::dataplane::sai::PORT_ATTR_ADVERTISED_FLOW_CONTROL_MODE;

    case SAI_PORT_ATTR_ADVERTISED_ASYMMETRIC_PAUSE_MODE:
      return lemming::dataplane::sai::
          PORT_ATTR_ADVERTISED_ASYMMETRIC_PAUSE_MODE;

    case SAI_PORT_ATTR_ADVERTISED_MEDIA_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_ADVERTISED_MEDIA_TYPE;

    case SAI_PORT_ATTR_ADVERTISED_OUI_CODE:
      return lemming::dataplane::sai::PORT_ATTR_ADVERTISED_OUI_CODE;

    case SAI_PORT_ATTR_PORT_VLAN_ID:
      return lemming::dataplane::sai::PORT_ATTR_PORT_VLAN_ID;

    case SAI_PORT_ATTR_DEFAULT_VLAN_PRIORITY:
      return lemming::dataplane::sai::PORT_ATTR_DEFAULT_VLAN_PRIORITY;

    case SAI_PORT_ATTR_DROP_UNTAGGED:
      return lemming::dataplane::sai::PORT_ATTR_DROP_UNTAGGED;

    case SAI_PORT_ATTR_DROP_TAGGED:
      return lemming::dataplane::sai::PORT_ATTR_DROP_TAGGED;

    case SAI_PORT_ATTR_INTERNAL_LOOPBACK_MODE:
      return lemming::dataplane::sai::PORT_ATTR_INTERNAL_LOOPBACK_MODE;

    case SAI_PORT_ATTR_USE_EXTENDED_FEC:
      return lemming::dataplane::sai::PORT_ATTR_USE_EXTENDED_FEC;

    case SAI_PORT_ATTR_FEC_MODE:
      return lemming::dataplane::sai::PORT_ATTR_FEC_MODE;

    case SAI_PORT_ATTR_FEC_MODE_EXTENDED:
      return lemming::dataplane::sai::PORT_ATTR_FEC_MODE_EXTENDED;

    case SAI_PORT_ATTR_UPDATE_DSCP:
      return lemming::dataplane::sai::PORT_ATTR_UPDATE_DSCP;

    case SAI_PORT_ATTR_MTU:
      return lemming::dataplane::sai::PORT_ATTR_MTU;

    case SAI_PORT_ATTR_FLOOD_STORM_CONTROL_POLICER_ID:
      return lemming::dataplane::sai::PORT_ATTR_FLOOD_STORM_CONTROL_POLICER_ID;

    case SAI_PORT_ATTR_BROADCAST_STORM_CONTROL_POLICER_ID:
      return lemming::dataplane::sai::
          PORT_ATTR_BROADCAST_STORM_CONTROL_POLICER_ID;

    case SAI_PORT_ATTR_MULTICAST_STORM_CONTROL_POLICER_ID:
      return lemming::dataplane::sai::
          PORT_ATTR_MULTICAST_STORM_CONTROL_POLICER_ID;

    case SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_MODE:
      return lemming::dataplane::sai::PORT_ATTR_GLOBAL_FLOW_CONTROL_MODE;

    case SAI_PORT_ATTR_INGRESS_ACL:
      return lemming::dataplane::sai::PORT_ATTR_INGRESS_ACL;

    case SAI_PORT_ATTR_EGRESS_ACL:
      return lemming::dataplane::sai::PORT_ATTR_EGRESS_ACL;

    case SAI_PORT_ATTR_INGRESS_MACSEC_ACL:
      return lemming::dataplane::sai::PORT_ATTR_INGRESS_MACSEC_ACL;

    case SAI_PORT_ATTR_EGRESS_MACSEC_ACL:
      return lemming::dataplane::sai::PORT_ATTR_EGRESS_MACSEC_ACL;

    case SAI_PORT_ATTR_MACSEC_PORT_LIST:
      return lemming::dataplane::sai::PORT_ATTR_MACSEC_PORT_LIST;

    case SAI_PORT_ATTR_INGRESS_MIRROR_SESSION:
      return lemming::dataplane::sai::PORT_ATTR_INGRESS_MIRROR_SESSION;

    case SAI_PORT_ATTR_EGRESS_MIRROR_SESSION:
      return lemming::dataplane::sai::PORT_ATTR_EGRESS_MIRROR_SESSION;

    case SAI_PORT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
      return lemming::dataplane::sai::PORT_ATTR_INGRESS_SAMPLEPACKET_ENABLE;

    case SAI_PORT_ATTR_EGRESS_SAMPLEPACKET_ENABLE:
      return lemming::dataplane::sai::PORT_ATTR_EGRESS_SAMPLEPACKET_ENABLE;

    case SAI_PORT_ATTR_INGRESS_SAMPLE_MIRROR_SESSION:
      return lemming::dataplane::sai::PORT_ATTR_INGRESS_SAMPLE_MIRROR_SESSION;

    case SAI_PORT_ATTR_EGRESS_SAMPLE_MIRROR_SESSION:
      return lemming::dataplane::sai::PORT_ATTR_EGRESS_SAMPLE_MIRROR_SESSION;

    case SAI_PORT_ATTR_POLICER_ID:
      return lemming::dataplane::sai::PORT_ATTR_POLICER_ID;

    case SAI_PORT_ATTR_QOS_DEFAULT_TC:
      return lemming::dataplane::sai::PORT_ATTR_QOS_DEFAULT_TC;

    case SAI_PORT_ATTR_QOS_DOT1P_TO_TC_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_DOT1P_TO_TC_MAP;

    case SAI_PORT_ATTR_QOS_DOT1P_TO_COLOR_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_DOT1P_TO_COLOR_MAP;

    case SAI_PORT_ATTR_QOS_DSCP_TO_TC_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_DSCP_TO_TC_MAP;

    case SAI_PORT_ATTR_QOS_DSCP_TO_COLOR_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_DSCP_TO_COLOR_MAP;

    case SAI_PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_TC_TO_QUEUE_MAP;

    case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP;

    case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP;

    case SAI_PORT_ATTR_QOS_TC_TO_PRIORITY_GROUP_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_TC_TO_PRIORITY_GROUP_MAP;

    case SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_PRIORITY_GROUP_MAP:
      return lemming::dataplane::sai::
          PORT_ATTR_QOS_PFC_PRIORITY_TO_PRIORITY_GROUP_MAP;

    case SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_QUEUE_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_PFC_PRIORITY_TO_QUEUE_MAP;

    case SAI_PORT_ATTR_QOS_SCHEDULER_PROFILE_ID:
      return lemming::dataplane::sai::PORT_ATTR_QOS_SCHEDULER_PROFILE_ID;

    case SAI_PORT_ATTR_QOS_INGRESS_BUFFER_PROFILE_LIST:
      return lemming::dataplane::sai::PORT_ATTR_QOS_INGRESS_BUFFER_PROFILE_LIST;

    case SAI_PORT_ATTR_QOS_EGRESS_BUFFER_PROFILE_LIST:
      return lemming::dataplane::sai::PORT_ATTR_QOS_EGRESS_BUFFER_PROFILE_LIST;

    case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_MODE:
      return lemming::dataplane::sai::PORT_ATTR_PRIORITY_FLOW_CONTROL_MODE;

    case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL:
      return lemming::dataplane::sai::PORT_ATTR_PRIORITY_FLOW_CONTROL;

    case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_RX:
      return lemming::dataplane::sai::PORT_ATTR_PRIORITY_FLOW_CONTROL_RX;

    case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_TX:
      return lemming::dataplane::sai::PORT_ATTR_PRIORITY_FLOW_CONTROL_TX;

    case SAI_PORT_ATTR_META_DATA:
      return lemming::dataplane::sai::PORT_ATTR_META_DATA;

    case SAI_PORT_ATTR_EGRESS_BLOCK_PORT_LIST:
      return lemming::dataplane::sai::PORT_ATTR_EGRESS_BLOCK_PORT_LIST;

    case SAI_PORT_ATTR_HW_PROFILE_ID:
      return lemming::dataplane::sai::PORT_ATTR_HW_PROFILE_ID;

    case SAI_PORT_ATTR_EEE_ENABLE:
      return lemming::dataplane::sai::PORT_ATTR_EEE_ENABLE;

    case SAI_PORT_ATTR_EEE_IDLE_TIME:
      return lemming::dataplane::sai::PORT_ATTR_EEE_IDLE_TIME;

    case SAI_PORT_ATTR_EEE_WAKE_TIME:
      return lemming::dataplane::sai::PORT_ATTR_EEE_WAKE_TIME;

    case SAI_PORT_ATTR_PORT_POOL_LIST:
      return lemming::dataplane::sai::PORT_ATTR_PORT_POOL_LIST;

    case SAI_PORT_ATTR_ISOLATION_GROUP:
      return lemming::dataplane::sai::PORT_ATTR_ISOLATION_GROUP;

    case SAI_PORT_ATTR_PKT_TX_ENABLE:
      return lemming::dataplane::sai::PORT_ATTR_PKT_TX_ENABLE;

    case SAI_PORT_ATTR_TAM_OBJECT:
      return lemming::dataplane::sai::PORT_ATTR_TAM_OBJECT;

    case SAI_PORT_ATTR_SERDES_PREEMPHASIS:
      return lemming::dataplane::sai::PORT_ATTR_SERDES_PREEMPHASIS;

    case SAI_PORT_ATTR_SERDES_IDRIVER:
      return lemming::dataplane::sai::PORT_ATTR_SERDES_IDRIVER;

    case SAI_PORT_ATTR_SERDES_IPREDRIVER:
      return lemming::dataplane::sai::PORT_ATTR_SERDES_IPREDRIVER;

    case SAI_PORT_ATTR_LINK_TRAINING_ENABLE:
      return lemming::dataplane::sai::PORT_ATTR_LINK_TRAINING_ENABLE;

    case SAI_PORT_ATTR_PTP_MODE:
      return lemming::dataplane::sai::PORT_ATTR_PTP_MODE;

    case SAI_PORT_ATTR_INTERFACE_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_INTERFACE_TYPE;

    case SAI_PORT_ATTR_ADVERTISED_INTERFACE_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_ADVERTISED_INTERFACE_TYPE;

    case SAI_PORT_ATTR_REFERENCE_CLOCK:
      return lemming::dataplane::sai::PORT_ATTR_REFERENCE_CLOCK;

    case SAI_PORT_ATTR_PRBS_POLYNOMIAL:
      return lemming::dataplane::sai::PORT_ATTR_PRBS_POLYNOMIAL;

    case SAI_PORT_ATTR_PORT_SERDES_ID:
      return lemming::dataplane::sai::PORT_ATTR_PORT_SERDES_ID;

    case SAI_PORT_ATTR_LINK_TRAINING_FAILURE_STATUS:
      return lemming::dataplane::sai::PORT_ATTR_LINK_TRAINING_FAILURE_STATUS;

    case SAI_PORT_ATTR_LINK_TRAINING_RX_STATUS:
      return lemming::dataplane::sai::PORT_ATTR_LINK_TRAINING_RX_STATUS;

    case SAI_PORT_ATTR_PRBS_CONFIG:
      return lemming::dataplane::sai::PORT_ATTR_PRBS_CONFIG;

    case SAI_PORT_ATTR_PRBS_LOCK_STATUS:
      return lemming::dataplane::sai::PORT_ATTR_PRBS_LOCK_STATUS;

    case SAI_PORT_ATTR_PRBS_LOCK_LOSS_STATUS:
      return lemming::dataplane::sai::PORT_ATTR_PRBS_LOCK_LOSS_STATUS;

    case SAI_PORT_ATTR_PRBS_RX_STATUS:
      return lemming::dataplane::sai::PORT_ATTR_PRBS_RX_STATUS;

    case SAI_PORT_ATTR_PRBS_RX_STATE:
      return lemming::dataplane::sai::PORT_ATTR_PRBS_RX_STATE;

    case SAI_PORT_ATTR_AUTO_NEG_STATUS:
      return lemming::dataplane::sai::PORT_ATTR_AUTO_NEG_STATUS;

    case SAI_PORT_ATTR_DISABLE_DECREMENT_TTL:
      return lemming::dataplane::sai::PORT_ATTR_DISABLE_DECREMENT_TTL;

    case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_MPLS_EXP_TO_TC_MAP;

    case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
      return lemming::dataplane::sai::PORT_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP;

    case SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
      return lemming::dataplane::sai::
          PORT_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP;

    case SAI_PORT_ATTR_TPID:
      return lemming::dataplane::sai::PORT_ATTR_TPID;

    case SAI_PORT_ATTR_ERR_STATUS_LIST:
      return lemming::dataplane::sai::PORT_ATTR_ERR_STATUS_LIST;

    case SAI_PORT_ATTR_FABRIC_ATTACHED:
      return lemming::dataplane::sai::PORT_ATTR_FABRIC_ATTACHED;

    case SAI_PORT_ATTR_FABRIC_ATTACHED_SWITCH_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_FABRIC_ATTACHED_SWITCH_TYPE;

    case SAI_PORT_ATTR_FABRIC_ATTACHED_SWITCH_ID:
      return lemming::dataplane::sai::PORT_ATTR_FABRIC_ATTACHED_SWITCH_ID;

    case SAI_PORT_ATTR_FABRIC_ATTACHED_PORT_INDEX:
      return lemming::dataplane::sai::PORT_ATTR_FABRIC_ATTACHED_PORT_INDEX;

    case SAI_PORT_ATTR_FABRIC_REACHABILITY:
      return lemming::dataplane::sai::PORT_ATTR_FABRIC_REACHABILITY;

    case SAI_PORT_ATTR_SYSTEM_PORT:
      return lemming::dataplane::sai::PORT_ATTR_SYSTEM_PORT;

    case SAI_PORT_ATTR_AUTO_NEG_FEC_MODE_OVERRIDE:
      return lemming::dataplane::sai::PORT_ATTR_AUTO_NEG_FEC_MODE_OVERRIDE;

    case SAI_PORT_ATTR_LOOPBACK_MODE:
      return lemming::dataplane::sai::PORT_ATTR_LOOPBACK_MODE;

    case SAI_PORT_ATTR_MDIX_MODE_STATUS:
      return lemming::dataplane::sai::PORT_ATTR_MDIX_MODE_STATUS;

    case SAI_PORT_ATTR_MDIX_MODE_CONFIG:
      return lemming::dataplane::sai::PORT_ATTR_MDIX_MODE_CONFIG;

    case SAI_PORT_ATTR_AUTO_NEG_CONFIG_MODE:
      return lemming::dataplane::sai::PORT_ATTR_AUTO_NEG_CONFIG_MODE;

    case SAI_PORT_ATTR_1000X_SGMII_SLAVE_AUTODETECT:
      return lemming::dataplane::sai::PORT_ATTR_1000X_SGMII_SLAVE_AUTODETECT;

    case SAI_PORT_ATTR_MODULE_TYPE:
      return lemming::dataplane::sai::PORT_ATTR_MODULE_TYPE;

    case SAI_PORT_ATTR_DUAL_MEDIA:
      return lemming::dataplane::sai::PORT_ATTR_DUAL_MEDIA;

    case SAI_PORT_ATTR_AUTO_NEG_FEC_MODE_EXTENDED:
      return lemming::dataplane::sai::PORT_ATTR_AUTO_NEG_FEC_MODE_EXTENDED;

    case SAI_PORT_ATTR_IPG:
      return lemming::dataplane::sai::PORT_ATTR_IPG;

    case SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_FORWARD:
      return lemming::dataplane::sai::PORT_ATTR_GLOBAL_FLOW_CONTROL_FORWARD;

    case SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_FORWARD:
      return lemming::dataplane::sai::PORT_ATTR_PRIORITY_FLOW_CONTROL_FORWARD;

    case SAI_PORT_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
      return lemming::dataplane::sai::
          PORT_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP;

    case SAI_PORT_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
      return lemming::dataplane::sai::
          PORT_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP;

    case SAI_PORT_ATTR_IPSEC_PORT:
      return lemming::dataplane::sai::PORT_ATTR_IPSEC_PORT;

    case SAI_PORT_ATTR_PFC_TC_DLD_INTERVAL_RANGE:
      return lemming::dataplane::sai::PORT_ATTR_PFC_TC_DLD_INTERVAL_RANGE;

    case SAI_PORT_ATTR_PFC_TC_DLD_INTERVAL:
      return lemming::dataplane::sai::PORT_ATTR_PFC_TC_DLD_INTERVAL;

    case SAI_PORT_ATTR_PFC_TC_DLR_INTERVAL_RANGE:
      return lemming::dataplane::sai::PORT_ATTR_PFC_TC_DLR_INTERVAL_RANGE;

    case SAI_PORT_ATTR_PFC_TC_DLR_INTERVAL:
      return lemming::dataplane::sai::PORT_ATTR_PFC_TC_DLR_INTERVAL;

    case SAI_PORT_ATTR_SUPPORTED_LINK_TRAINING_MODE:
      return lemming::dataplane::sai::PORT_ATTR_SUPPORTED_LINK_TRAINING_MODE;

    case SAI_PORT_ATTR_RX_SIGNAL_DETECT:
      return lemming::dataplane::sai::PORT_ATTR_RX_SIGNAL_DETECT;

    case SAI_PORT_ATTR_RX_LOCK_STATUS:
      return lemming::dataplane::sai::PORT_ATTR_RX_LOCK_STATUS;

    case SAI_PORT_ATTR_PCS_RX_LINK_STATUS:
      return lemming::dataplane::sai::PORT_ATTR_PCS_RX_LINK_STATUS;

    case SAI_PORT_ATTR_FEC_ALIGNMENT_LOCK:
      return lemming::dataplane::sai::PORT_ATTR_FEC_ALIGNMENT_LOCK;

    case SAI_PORT_ATTR_FABRIC_ISOLATE:
      return lemming::dataplane::sai::PORT_ATTR_FABRIC_ISOLATE;

    case SAI_PORT_ATTR_MAX_FEC_SYMBOL_ERRORS_DETECTABLE:
      return lemming::dataplane::sai::
          PORT_ATTR_MAX_FEC_SYMBOL_ERRORS_DETECTABLE;

    default:
      return lemming::dataplane::sai::PORT_ATTR_UNSPECIFIED;
  }
}
sai_port_attr_t convert_sai_port_attr_t_to_sai(
    lemming::dataplane::sai::PortAttr val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_ATTR_TYPE:
      return SAI_PORT_ATTR_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_OPER_STATUS:
      return SAI_PORT_ATTR_OPER_STATUS;

    case lemming::dataplane::sai::PORT_ATTR_SUPPORTED_BREAKOUT_MODE_TYPE:
      return SAI_PORT_ATTR_SUPPORTED_BREAKOUT_MODE_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_CURRENT_BREAKOUT_MODE_TYPE:
      return SAI_PORT_ATTR_CURRENT_BREAKOUT_MODE_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_QOS_NUMBER_OF_QUEUES:
      return SAI_PORT_ATTR_QOS_NUMBER_OF_QUEUES;

    case lemming::dataplane::sai::PORT_ATTR_QOS_QUEUE_LIST:
      return SAI_PORT_ATTR_QOS_QUEUE_LIST;

    case lemming::dataplane::sai::PORT_ATTR_QOS_NUMBER_OF_SCHEDULER_GROUPS:
      return SAI_PORT_ATTR_QOS_NUMBER_OF_SCHEDULER_GROUPS;

    case lemming::dataplane::sai::PORT_ATTR_QOS_SCHEDULER_GROUP_LIST:
      return SAI_PORT_ATTR_QOS_SCHEDULER_GROUP_LIST;

    case lemming::dataplane::sai::PORT_ATTR_QOS_MAXIMUM_HEADROOM_SIZE:
      return SAI_PORT_ATTR_QOS_MAXIMUM_HEADROOM_SIZE;

    case lemming::dataplane::sai::PORT_ATTR_SUPPORTED_SPEED:
      return SAI_PORT_ATTR_SUPPORTED_SPEED;

    case lemming::dataplane::sai::PORT_ATTR_SUPPORTED_FEC_MODE:
      return SAI_PORT_ATTR_SUPPORTED_FEC_MODE;

    case lemming::dataplane::sai::PORT_ATTR_SUPPORTED_FEC_MODE_EXTENDED:
      return SAI_PORT_ATTR_SUPPORTED_FEC_MODE_EXTENDED;

    case lemming::dataplane::sai::PORT_ATTR_SUPPORTED_HALF_DUPLEX_SPEED:
      return SAI_PORT_ATTR_SUPPORTED_HALF_DUPLEX_SPEED;

    case lemming::dataplane::sai::PORT_ATTR_SUPPORTED_AUTO_NEG_MODE:
      return SAI_PORT_ATTR_SUPPORTED_AUTO_NEG_MODE;

    case lemming::dataplane::sai::PORT_ATTR_SUPPORTED_FLOW_CONTROL_MODE:
      return SAI_PORT_ATTR_SUPPORTED_FLOW_CONTROL_MODE;

    case lemming::dataplane::sai::PORT_ATTR_SUPPORTED_ASYMMETRIC_PAUSE_MODE:
      return SAI_PORT_ATTR_SUPPORTED_ASYMMETRIC_PAUSE_MODE;

    case lemming::dataplane::sai::PORT_ATTR_SUPPORTED_MEDIA_TYPE:
      return SAI_PORT_ATTR_SUPPORTED_MEDIA_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_SPEED:
      return SAI_PORT_ATTR_REMOTE_ADVERTISED_SPEED;

    case lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_FEC_MODE:
      return SAI_PORT_ATTR_REMOTE_ADVERTISED_FEC_MODE;

    case lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_FEC_MODE_EXTENDED:
      return SAI_PORT_ATTR_REMOTE_ADVERTISED_FEC_MODE_EXTENDED;

    case lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_HALF_DUPLEX_SPEED:
      return SAI_PORT_ATTR_REMOTE_ADVERTISED_HALF_DUPLEX_SPEED;

    case lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_AUTO_NEG_MODE:
      return SAI_PORT_ATTR_REMOTE_ADVERTISED_AUTO_NEG_MODE;

    case lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_FLOW_CONTROL_MODE:
      return SAI_PORT_ATTR_REMOTE_ADVERTISED_FLOW_CONTROL_MODE;

    case lemming::dataplane::sai::
        PORT_ATTR_REMOTE_ADVERTISED_ASYMMETRIC_PAUSE_MODE:
      return SAI_PORT_ATTR_REMOTE_ADVERTISED_ASYMMETRIC_PAUSE_MODE;

    case lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_MEDIA_TYPE:
      return SAI_PORT_ATTR_REMOTE_ADVERTISED_MEDIA_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_REMOTE_ADVERTISED_OUI_CODE:
      return SAI_PORT_ATTR_REMOTE_ADVERTISED_OUI_CODE;

    case lemming::dataplane::sai::PORT_ATTR_NUMBER_OF_INGRESS_PRIORITY_GROUPS:
      return SAI_PORT_ATTR_NUMBER_OF_INGRESS_PRIORITY_GROUPS;

    case lemming::dataplane::sai::PORT_ATTR_INGRESS_PRIORITY_GROUP_LIST:
      return SAI_PORT_ATTR_INGRESS_PRIORITY_GROUP_LIST;

    case lemming::dataplane::sai::PORT_ATTR_EYE_VALUES:
      return SAI_PORT_ATTR_EYE_VALUES;

    case lemming::dataplane::sai::PORT_ATTR_OPER_SPEED:
      return SAI_PORT_ATTR_OPER_SPEED;

    case lemming::dataplane::sai::PORT_ATTR_HW_LANE_LIST:
      return SAI_PORT_ATTR_HW_LANE_LIST;

    case lemming::dataplane::sai::PORT_ATTR_SPEED:
      return SAI_PORT_ATTR_SPEED;

    case lemming::dataplane::sai::PORT_ATTR_FULL_DUPLEX_MODE:
      return SAI_PORT_ATTR_FULL_DUPLEX_MODE;

    case lemming::dataplane::sai::PORT_ATTR_AUTO_NEG_MODE:
      return SAI_PORT_ATTR_AUTO_NEG_MODE;

    case lemming::dataplane::sai::PORT_ATTR_ADMIN_STATE:
      return SAI_PORT_ATTR_ADMIN_STATE;

    case lemming::dataplane::sai::PORT_ATTR_MEDIA_TYPE:
      return SAI_PORT_ATTR_MEDIA_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_ADVERTISED_SPEED:
      return SAI_PORT_ATTR_ADVERTISED_SPEED;

    case lemming::dataplane::sai::PORT_ATTR_ADVERTISED_FEC_MODE:
      return SAI_PORT_ATTR_ADVERTISED_FEC_MODE;

    case lemming::dataplane::sai::PORT_ATTR_ADVERTISED_FEC_MODE_EXTENDED:
      return SAI_PORT_ATTR_ADVERTISED_FEC_MODE_EXTENDED;

    case lemming::dataplane::sai::PORT_ATTR_ADVERTISED_HALF_DUPLEX_SPEED:
      return SAI_PORT_ATTR_ADVERTISED_HALF_DUPLEX_SPEED;

    case lemming::dataplane::sai::PORT_ATTR_ADVERTISED_AUTO_NEG_MODE:
      return SAI_PORT_ATTR_ADVERTISED_AUTO_NEG_MODE;

    case lemming::dataplane::sai::PORT_ATTR_ADVERTISED_FLOW_CONTROL_MODE:
      return SAI_PORT_ATTR_ADVERTISED_FLOW_CONTROL_MODE;

    case lemming::dataplane::sai::PORT_ATTR_ADVERTISED_ASYMMETRIC_PAUSE_MODE:
      return SAI_PORT_ATTR_ADVERTISED_ASYMMETRIC_PAUSE_MODE;

    case lemming::dataplane::sai::PORT_ATTR_ADVERTISED_MEDIA_TYPE:
      return SAI_PORT_ATTR_ADVERTISED_MEDIA_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_ADVERTISED_OUI_CODE:
      return SAI_PORT_ATTR_ADVERTISED_OUI_CODE;

    case lemming::dataplane::sai::PORT_ATTR_PORT_VLAN_ID:
      return SAI_PORT_ATTR_PORT_VLAN_ID;

    case lemming::dataplane::sai::PORT_ATTR_DEFAULT_VLAN_PRIORITY:
      return SAI_PORT_ATTR_DEFAULT_VLAN_PRIORITY;

    case lemming::dataplane::sai::PORT_ATTR_DROP_UNTAGGED:
      return SAI_PORT_ATTR_DROP_UNTAGGED;

    case lemming::dataplane::sai::PORT_ATTR_DROP_TAGGED:
      return SAI_PORT_ATTR_DROP_TAGGED;

    case lemming::dataplane::sai::PORT_ATTR_INTERNAL_LOOPBACK_MODE:
      return SAI_PORT_ATTR_INTERNAL_LOOPBACK_MODE;

    case lemming::dataplane::sai::PORT_ATTR_USE_EXTENDED_FEC:
      return SAI_PORT_ATTR_USE_EXTENDED_FEC;

    case lemming::dataplane::sai::PORT_ATTR_FEC_MODE:
      return SAI_PORT_ATTR_FEC_MODE;

    case lemming::dataplane::sai::PORT_ATTR_FEC_MODE_EXTENDED:
      return SAI_PORT_ATTR_FEC_MODE_EXTENDED;

    case lemming::dataplane::sai::PORT_ATTR_UPDATE_DSCP:
      return SAI_PORT_ATTR_UPDATE_DSCP;

    case lemming::dataplane::sai::PORT_ATTR_MTU:
      return SAI_PORT_ATTR_MTU;

    case lemming::dataplane::sai::PORT_ATTR_FLOOD_STORM_CONTROL_POLICER_ID:
      return SAI_PORT_ATTR_FLOOD_STORM_CONTROL_POLICER_ID;

    case lemming::dataplane::sai::PORT_ATTR_BROADCAST_STORM_CONTROL_POLICER_ID:
      return SAI_PORT_ATTR_BROADCAST_STORM_CONTROL_POLICER_ID;

    case lemming::dataplane::sai::PORT_ATTR_MULTICAST_STORM_CONTROL_POLICER_ID:
      return SAI_PORT_ATTR_MULTICAST_STORM_CONTROL_POLICER_ID;

    case lemming::dataplane::sai::PORT_ATTR_GLOBAL_FLOW_CONTROL_MODE:
      return SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_MODE;

    case lemming::dataplane::sai::PORT_ATTR_INGRESS_ACL:
      return SAI_PORT_ATTR_INGRESS_ACL;

    case lemming::dataplane::sai::PORT_ATTR_EGRESS_ACL:
      return SAI_PORT_ATTR_EGRESS_ACL;

    case lemming::dataplane::sai::PORT_ATTR_INGRESS_MACSEC_ACL:
      return SAI_PORT_ATTR_INGRESS_MACSEC_ACL;

    case lemming::dataplane::sai::PORT_ATTR_EGRESS_MACSEC_ACL:
      return SAI_PORT_ATTR_EGRESS_MACSEC_ACL;

    case lemming::dataplane::sai::PORT_ATTR_MACSEC_PORT_LIST:
      return SAI_PORT_ATTR_MACSEC_PORT_LIST;

    case lemming::dataplane::sai::PORT_ATTR_INGRESS_MIRROR_SESSION:
      return SAI_PORT_ATTR_INGRESS_MIRROR_SESSION;

    case lemming::dataplane::sai::PORT_ATTR_EGRESS_MIRROR_SESSION:
      return SAI_PORT_ATTR_EGRESS_MIRROR_SESSION;

    case lemming::dataplane::sai::PORT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
      return SAI_PORT_ATTR_INGRESS_SAMPLEPACKET_ENABLE;

    case lemming::dataplane::sai::PORT_ATTR_EGRESS_SAMPLEPACKET_ENABLE:
      return SAI_PORT_ATTR_EGRESS_SAMPLEPACKET_ENABLE;

    case lemming::dataplane::sai::PORT_ATTR_INGRESS_SAMPLE_MIRROR_SESSION:
      return SAI_PORT_ATTR_INGRESS_SAMPLE_MIRROR_SESSION;

    case lemming::dataplane::sai::PORT_ATTR_EGRESS_SAMPLE_MIRROR_SESSION:
      return SAI_PORT_ATTR_EGRESS_SAMPLE_MIRROR_SESSION;

    case lemming::dataplane::sai::PORT_ATTR_POLICER_ID:
      return SAI_PORT_ATTR_POLICER_ID;

    case lemming::dataplane::sai::PORT_ATTR_QOS_DEFAULT_TC:
      return SAI_PORT_ATTR_QOS_DEFAULT_TC;

    case lemming::dataplane::sai::PORT_ATTR_QOS_DOT1P_TO_TC_MAP:
      return SAI_PORT_ATTR_QOS_DOT1P_TO_TC_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_DOT1P_TO_COLOR_MAP:
      return SAI_PORT_ATTR_QOS_DOT1P_TO_COLOR_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_DSCP_TO_TC_MAP:
      return SAI_PORT_ATTR_QOS_DSCP_TO_TC_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_DSCP_TO_COLOR_MAP:
      return SAI_PORT_ATTR_QOS_DSCP_TO_COLOR_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
      return SAI_PORT_ATTR_QOS_TC_TO_QUEUE_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
      return SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
      return SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_TC_TO_PRIORITY_GROUP_MAP:
      return SAI_PORT_ATTR_QOS_TC_TO_PRIORITY_GROUP_MAP;

    case lemming::dataplane::sai::
        PORT_ATTR_QOS_PFC_PRIORITY_TO_PRIORITY_GROUP_MAP:
      return SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_PRIORITY_GROUP_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_PFC_PRIORITY_TO_QUEUE_MAP:
      return SAI_PORT_ATTR_QOS_PFC_PRIORITY_TO_QUEUE_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_SCHEDULER_PROFILE_ID:
      return SAI_PORT_ATTR_QOS_SCHEDULER_PROFILE_ID;

    case lemming::dataplane::sai::PORT_ATTR_QOS_INGRESS_BUFFER_PROFILE_LIST:
      return SAI_PORT_ATTR_QOS_INGRESS_BUFFER_PROFILE_LIST;

    case lemming::dataplane::sai::PORT_ATTR_QOS_EGRESS_BUFFER_PROFILE_LIST:
      return SAI_PORT_ATTR_QOS_EGRESS_BUFFER_PROFILE_LIST;

    case lemming::dataplane::sai::PORT_ATTR_PRIORITY_FLOW_CONTROL_MODE:
      return SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_MODE;

    case lemming::dataplane::sai::PORT_ATTR_PRIORITY_FLOW_CONTROL:
      return SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL;

    case lemming::dataplane::sai::PORT_ATTR_PRIORITY_FLOW_CONTROL_RX:
      return SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_RX;

    case lemming::dataplane::sai::PORT_ATTR_PRIORITY_FLOW_CONTROL_TX:
      return SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_TX;

    case lemming::dataplane::sai::PORT_ATTR_META_DATA:
      return SAI_PORT_ATTR_META_DATA;

    case lemming::dataplane::sai::PORT_ATTR_EGRESS_BLOCK_PORT_LIST:
      return SAI_PORT_ATTR_EGRESS_BLOCK_PORT_LIST;

    case lemming::dataplane::sai::PORT_ATTR_HW_PROFILE_ID:
      return SAI_PORT_ATTR_HW_PROFILE_ID;

    case lemming::dataplane::sai::PORT_ATTR_EEE_ENABLE:
      return SAI_PORT_ATTR_EEE_ENABLE;

    case lemming::dataplane::sai::PORT_ATTR_EEE_IDLE_TIME:
      return SAI_PORT_ATTR_EEE_IDLE_TIME;

    case lemming::dataplane::sai::PORT_ATTR_EEE_WAKE_TIME:
      return SAI_PORT_ATTR_EEE_WAKE_TIME;

    case lemming::dataplane::sai::PORT_ATTR_PORT_POOL_LIST:
      return SAI_PORT_ATTR_PORT_POOL_LIST;

    case lemming::dataplane::sai::PORT_ATTR_ISOLATION_GROUP:
      return SAI_PORT_ATTR_ISOLATION_GROUP;

    case lemming::dataplane::sai::PORT_ATTR_PKT_TX_ENABLE:
      return SAI_PORT_ATTR_PKT_TX_ENABLE;

    case lemming::dataplane::sai::PORT_ATTR_TAM_OBJECT:
      return SAI_PORT_ATTR_TAM_OBJECT;

    case lemming::dataplane::sai::PORT_ATTR_SERDES_PREEMPHASIS:
      return SAI_PORT_ATTR_SERDES_PREEMPHASIS;

    case lemming::dataplane::sai::PORT_ATTR_SERDES_IDRIVER:
      return SAI_PORT_ATTR_SERDES_IDRIVER;

    case lemming::dataplane::sai::PORT_ATTR_SERDES_IPREDRIVER:
      return SAI_PORT_ATTR_SERDES_IPREDRIVER;

    case lemming::dataplane::sai::PORT_ATTR_LINK_TRAINING_ENABLE:
      return SAI_PORT_ATTR_LINK_TRAINING_ENABLE;

    case lemming::dataplane::sai::PORT_ATTR_PTP_MODE:
      return SAI_PORT_ATTR_PTP_MODE;

    case lemming::dataplane::sai::PORT_ATTR_INTERFACE_TYPE:
      return SAI_PORT_ATTR_INTERFACE_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_ADVERTISED_INTERFACE_TYPE:
      return SAI_PORT_ATTR_ADVERTISED_INTERFACE_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_REFERENCE_CLOCK:
      return SAI_PORT_ATTR_REFERENCE_CLOCK;

    case lemming::dataplane::sai::PORT_ATTR_PRBS_POLYNOMIAL:
      return SAI_PORT_ATTR_PRBS_POLYNOMIAL;

    case lemming::dataplane::sai::PORT_ATTR_PORT_SERDES_ID:
      return SAI_PORT_ATTR_PORT_SERDES_ID;

    case lemming::dataplane::sai::PORT_ATTR_LINK_TRAINING_FAILURE_STATUS:
      return SAI_PORT_ATTR_LINK_TRAINING_FAILURE_STATUS;

    case lemming::dataplane::sai::PORT_ATTR_LINK_TRAINING_RX_STATUS:
      return SAI_PORT_ATTR_LINK_TRAINING_RX_STATUS;

    case lemming::dataplane::sai::PORT_ATTR_PRBS_CONFIG:
      return SAI_PORT_ATTR_PRBS_CONFIG;

    case lemming::dataplane::sai::PORT_ATTR_PRBS_LOCK_STATUS:
      return SAI_PORT_ATTR_PRBS_LOCK_STATUS;

    case lemming::dataplane::sai::PORT_ATTR_PRBS_LOCK_LOSS_STATUS:
      return SAI_PORT_ATTR_PRBS_LOCK_LOSS_STATUS;

    case lemming::dataplane::sai::PORT_ATTR_PRBS_RX_STATUS:
      return SAI_PORT_ATTR_PRBS_RX_STATUS;

    case lemming::dataplane::sai::PORT_ATTR_PRBS_RX_STATE:
      return SAI_PORT_ATTR_PRBS_RX_STATE;

    case lemming::dataplane::sai::PORT_ATTR_AUTO_NEG_STATUS:
      return SAI_PORT_ATTR_AUTO_NEG_STATUS;

    case lemming::dataplane::sai::PORT_ATTR_DISABLE_DECREMENT_TTL:
      return SAI_PORT_ATTR_DISABLE_DECREMENT_TTL;

    case lemming::dataplane::sai::PORT_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
      return SAI_PORT_ATTR_QOS_MPLS_EXP_TO_TC_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
      return SAI_PORT_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP;

    case lemming::dataplane::sai::PORT_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
      return SAI_PORT_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP;

    case lemming::dataplane::sai::PORT_ATTR_TPID:
      return SAI_PORT_ATTR_TPID;

    case lemming::dataplane::sai::PORT_ATTR_ERR_STATUS_LIST:
      return SAI_PORT_ATTR_ERR_STATUS_LIST;

    case lemming::dataplane::sai::PORT_ATTR_FABRIC_ATTACHED:
      return SAI_PORT_ATTR_FABRIC_ATTACHED;

    case lemming::dataplane::sai::PORT_ATTR_FABRIC_ATTACHED_SWITCH_TYPE:
      return SAI_PORT_ATTR_FABRIC_ATTACHED_SWITCH_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_FABRIC_ATTACHED_SWITCH_ID:
      return SAI_PORT_ATTR_FABRIC_ATTACHED_SWITCH_ID;

    case lemming::dataplane::sai::PORT_ATTR_FABRIC_ATTACHED_PORT_INDEX:
      return SAI_PORT_ATTR_FABRIC_ATTACHED_PORT_INDEX;

    case lemming::dataplane::sai::PORT_ATTR_FABRIC_REACHABILITY:
      return SAI_PORT_ATTR_FABRIC_REACHABILITY;

    case lemming::dataplane::sai::PORT_ATTR_SYSTEM_PORT:
      return SAI_PORT_ATTR_SYSTEM_PORT;

    case lemming::dataplane::sai::PORT_ATTR_AUTO_NEG_FEC_MODE_OVERRIDE:
      return SAI_PORT_ATTR_AUTO_NEG_FEC_MODE_OVERRIDE;

    case lemming::dataplane::sai::PORT_ATTR_LOOPBACK_MODE:
      return SAI_PORT_ATTR_LOOPBACK_MODE;

    case lemming::dataplane::sai::PORT_ATTR_MDIX_MODE_STATUS:
      return SAI_PORT_ATTR_MDIX_MODE_STATUS;

    case lemming::dataplane::sai::PORT_ATTR_MDIX_MODE_CONFIG:
      return SAI_PORT_ATTR_MDIX_MODE_CONFIG;

    case lemming::dataplane::sai::PORT_ATTR_AUTO_NEG_CONFIG_MODE:
      return SAI_PORT_ATTR_AUTO_NEG_CONFIG_MODE;

    case lemming::dataplane::sai::PORT_ATTR_1000X_SGMII_SLAVE_AUTODETECT:
      return SAI_PORT_ATTR_1000X_SGMII_SLAVE_AUTODETECT;

    case lemming::dataplane::sai::PORT_ATTR_MODULE_TYPE:
      return SAI_PORT_ATTR_MODULE_TYPE;

    case lemming::dataplane::sai::PORT_ATTR_DUAL_MEDIA:
      return SAI_PORT_ATTR_DUAL_MEDIA;

    case lemming::dataplane::sai::PORT_ATTR_AUTO_NEG_FEC_MODE_EXTENDED:
      return SAI_PORT_ATTR_AUTO_NEG_FEC_MODE_EXTENDED;

    case lemming::dataplane::sai::PORT_ATTR_IPG:
      return SAI_PORT_ATTR_IPG;

    case lemming::dataplane::sai::PORT_ATTR_GLOBAL_FLOW_CONTROL_FORWARD:
      return SAI_PORT_ATTR_GLOBAL_FLOW_CONTROL_FORWARD;

    case lemming::dataplane::sai::PORT_ATTR_PRIORITY_FLOW_CONTROL_FORWARD:
      return SAI_PORT_ATTR_PRIORITY_FLOW_CONTROL_FORWARD;

    case lemming::dataplane::sai::PORT_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
      return SAI_PORT_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP;

    case lemming::dataplane::sai::
        PORT_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
      return SAI_PORT_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP;

    case lemming::dataplane::sai::PORT_ATTR_IPSEC_PORT:
      return SAI_PORT_ATTR_IPSEC_PORT;

    case lemming::dataplane::sai::PORT_ATTR_PFC_TC_DLD_INTERVAL_RANGE:
      return SAI_PORT_ATTR_PFC_TC_DLD_INTERVAL_RANGE;

    case lemming::dataplane::sai::PORT_ATTR_PFC_TC_DLD_INTERVAL:
      return SAI_PORT_ATTR_PFC_TC_DLD_INTERVAL;

    case lemming::dataplane::sai::PORT_ATTR_PFC_TC_DLR_INTERVAL_RANGE:
      return SAI_PORT_ATTR_PFC_TC_DLR_INTERVAL_RANGE;

    case lemming::dataplane::sai::PORT_ATTR_PFC_TC_DLR_INTERVAL:
      return SAI_PORT_ATTR_PFC_TC_DLR_INTERVAL;

    case lemming::dataplane::sai::PORT_ATTR_SUPPORTED_LINK_TRAINING_MODE:
      return SAI_PORT_ATTR_SUPPORTED_LINK_TRAINING_MODE;

    case lemming::dataplane::sai::PORT_ATTR_RX_SIGNAL_DETECT:
      return SAI_PORT_ATTR_RX_SIGNAL_DETECT;

    case lemming::dataplane::sai::PORT_ATTR_RX_LOCK_STATUS:
      return SAI_PORT_ATTR_RX_LOCK_STATUS;

    case lemming::dataplane::sai::PORT_ATTR_PCS_RX_LINK_STATUS:
      return SAI_PORT_ATTR_PCS_RX_LINK_STATUS;

    case lemming::dataplane::sai::PORT_ATTR_FEC_ALIGNMENT_LOCK:
      return SAI_PORT_ATTR_FEC_ALIGNMENT_LOCK;

    case lemming::dataplane::sai::PORT_ATTR_FABRIC_ISOLATE:
      return SAI_PORT_ATTR_FABRIC_ISOLATE;

    case lemming::dataplane::sai::PORT_ATTR_MAX_FEC_SYMBOL_ERRORS_DETECTABLE:
      return SAI_PORT_ATTR_MAX_FEC_SYMBOL_ERRORS_DETECTABLE;

    default:
      return SAI_PORT_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_port_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::PortAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortAutoNegConfigMode
convert_sai_port_auto_neg_config_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_AUTO_NEG_CONFIG_MODE_DISABLED:
      return lemming::dataplane::sai::PORT_AUTO_NEG_CONFIG_MODE_DISABLED;

    case SAI_PORT_AUTO_NEG_CONFIG_MODE_AUTO:
      return lemming::dataplane::sai::PORT_AUTO_NEG_CONFIG_MODE_AUTO;

    case SAI_PORT_AUTO_NEG_CONFIG_MODE_SLAVE:
      return lemming::dataplane::sai::PORT_AUTO_NEG_CONFIG_MODE_SLAVE;

    case SAI_PORT_AUTO_NEG_CONFIG_MODE_MASTER:
      return lemming::dataplane::sai::PORT_AUTO_NEG_CONFIG_MODE_MASTER;

    default:
      return lemming::dataplane::sai::PORT_AUTO_NEG_CONFIG_MODE_UNSPECIFIED;
  }
}
sai_port_auto_neg_config_mode_t convert_sai_port_auto_neg_config_mode_t_to_sai(
    lemming::dataplane::sai::PortAutoNegConfigMode val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_AUTO_NEG_CONFIG_MODE_DISABLED:
      return SAI_PORT_AUTO_NEG_CONFIG_MODE_DISABLED;

    case lemming::dataplane::sai::PORT_AUTO_NEG_CONFIG_MODE_AUTO:
      return SAI_PORT_AUTO_NEG_CONFIG_MODE_AUTO;

    case lemming::dataplane::sai::PORT_AUTO_NEG_CONFIG_MODE_SLAVE:
      return SAI_PORT_AUTO_NEG_CONFIG_MODE_SLAVE;

    case lemming::dataplane::sai::PORT_AUTO_NEG_CONFIG_MODE_MASTER:
      return SAI_PORT_AUTO_NEG_CONFIG_MODE_MASTER;

    default:
      return SAI_PORT_AUTO_NEG_CONFIG_MODE_DISABLED;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_auto_neg_config_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_port_auto_neg_config_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_auto_neg_config_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_auto_neg_config_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::PortAutoNegConfigMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortBreakoutModeType
convert_sai_port_breakout_mode_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_BREAKOUT_MODE_TYPE_1_LANE:
      return lemming::dataplane::sai::PORT_BREAKOUT_MODE_TYPE_1_LANE;

    case SAI_PORT_BREAKOUT_MODE_TYPE_2_LANE:
      return lemming::dataplane::sai::PORT_BREAKOUT_MODE_TYPE_2_LANE;

    case SAI_PORT_BREAKOUT_MODE_TYPE_4_LANE:
      return lemming::dataplane::sai::PORT_BREAKOUT_MODE_TYPE_4_LANE;

    case SAI_PORT_BREAKOUT_MODE_TYPE_MAX:
      return lemming::dataplane::sai::PORT_BREAKOUT_MODE_TYPE_MAX;

    default:
      return lemming::dataplane::sai::PORT_BREAKOUT_MODE_TYPE_UNSPECIFIED;
  }
}
sai_port_breakout_mode_type_t convert_sai_port_breakout_mode_type_t_to_sai(
    lemming::dataplane::sai::PortBreakoutModeType val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_BREAKOUT_MODE_TYPE_1_LANE:
      return SAI_PORT_BREAKOUT_MODE_TYPE_1_LANE;

    case lemming::dataplane::sai::PORT_BREAKOUT_MODE_TYPE_2_LANE:
      return SAI_PORT_BREAKOUT_MODE_TYPE_2_LANE;

    case lemming::dataplane::sai::PORT_BREAKOUT_MODE_TYPE_4_LANE:
      return SAI_PORT_BREAKOUT_MODE_TYPE_4_LANE;

    case lemming::dataplane::sai::PORT_BREAKOUT_MODE_TYPE_MAX:
      return SAI_PORT_BREAKOUT_MODE_TYPE_MAX;

    default:
      return SAI_PORT_BREAKOUT_MODE_TYPE_1_LANE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_breakout_mode_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_port_breakout_mode_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_breakout_mode_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_breakout_mode_type_t_to_sai(
        static_cast<lemming::dataplane::sai::PortBreakoutModeType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortConnectorAttr
convert_sai_port_connector_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_PORT_ID:
      return lemming::dataplane::sai::PORT_CONNECTOR_ATTR_SYSTEM_SIDE_PORT_ID;

    case SAI_PORT_CONNECTOR_ATTR_LINE_SIDE_PORT_ID:
      return lemming::dataplane::sai::PORT_CONNECTOR_ATTR_LINE_SIDE_PORT_ID;

    case SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_FAILOVER_PORT_ID:
      return lemming::dataplane::sai::
          PORT_CONNECTOR_ATTR_SYSTEM_SIDE_FAILOVER_PORT_ID;

    case SAI_PORT_CONNECTOR_ATTR_LINE_SIDE_FAILOVER_PORT_ID:
      return lemming::dataplane::sai::
          PORT_CONNECTOR_ATTR_LINE_SIDE_FAILOVER_PORT_ID;

    case SAI_PORT_CONNECTOR_ATTR_FAILOVER_MODE:
      return lemming::dataplane::sai::PORT_CONNECTOR_ATTR_FAILOVER_MODE;

    default:
      return lemming::dataplane::sai::PORT_CONNECTOR_ATTR_UNSPECIFIED;
  }
}
sai_port_connector_attr_t convert_sai_port_connector_attr_t_to_sai(
    lemming::dataplane::sai::PortConnectorAttr val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_CONNECTOR_ATTR_SYSTEM_SIDE_PORT_ID:
      return SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_PORT_ID;

    case lemming::dataplane::sai::PORT_CONNECTOR_ATTR_LINE_SIDE_PORT_ID:
      return SAI_PORT_CONNECTOR_ATTR_LINE_SIDE_PORT_ID;

    case lemming::dataplane::sai::
        PORT_CONNECTOR_ATTR_SYSTEM_SIDE_FAILOVER_PORT_ID:
      return SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_FAILOVER_PORT_ID;

    case lemming::dataplane::sai::
        PORT_CONNECTOR_ATTR_LINE_SIDE_FAILOVER_PORT_ID:
      return SAI_PORT_CONNECTOR_ATTR_LINE_SIDE_FAILOVER_PORT_ID;

    case lemming::dataplane::sai::PORT_CONNECTOR_ATTR_FAILOVER_MODE:
      return SAI_PORT_CONNECTOR_ATTR_FAILOVER_MODE;

    default:
      return SAI_PORT_CONNECTOR_ATTR_SYSTEM_SIDE_PORT_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_connector_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_connector_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_connector_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_connector_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::PortConnectorAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortConnectorFailoverMode
convert_sai_port_connector_failover_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_CONNECTOR_FAILOVER_MODE_DISABLE:
      return lemming::dataplane::sai::PORT_CONNECTOR_FAILOVER_MODE_DISABLE;

    case SAI_PORT_CONNECTOR_FAILOVER_MODE_PRIMARY:
      return lemming::dataplane::sai::PORT_CONNECTOR_FAILOVER_MODE_PRIMARY;

    case SAI_PORT_CONNECTOR_FAILOVER_MODE_SECONDARY:
      return lemming::dataplane::sai::PORT_CONNECTOR_FAILOVER_MODE_SECONDARY;

    default:
      return lemming::dataplane::sai::PORT_CONNECTOR_FAILOVER_MODE_UNSPECIFIED;
  }
}
sai_port_connector_failover_mode_t
convert_sai_port_connector_failover_mode_t_to_sai(
    lemming::dataplane::sai::PortConnectorFailoverMode val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_CONNECTOR_FAILOVER_MODE_DISABLE:
      return SAI_PORT_CONNECTOR_FAILOVER_MODE_DISABLE;

    case lemming::dataplane::sai::PORT_CONNECTOR_FAILOVER_MODE_PRIMARY:
      return SAI_PORT_CONNECTOR_FAILOVER_MODE_PRIMARY;

    case lemming::dataplane::sai::PORT_CONNECTOR_FAILOVER_MODE_SECONDARY:
      return SAI_PORT_CONNECTOR_FAILOVER_MODE_SECONDARY;

    default:
      return SAI_PORT_CONNECTOR_FAILOVER_MODE_DISABLE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_connector_failover_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_port_connector_failover_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_connector_failover_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_connector_failover_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::PortConnectorFailoverMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortDualMedia convert_sai_port_dual_media_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_DUAL_MEDIA_NONE:
      return lemming::dataplane::sai::PORT_DUAL_MEDIA_NONE;

    case SAI_PORT_DUAL_MEDIA_COPPER_ONLY:
      return lemming::dataplane::sai::PORT_DUAL_MEDIA_COPPER_ONLY;

    case SAI_PORT_DUAL_MEDIA_FIBER_ONLY:
      return lemming::dataplane::sai::PORT_DUAL_MEDIA_FIBER_ONLY;

    case SAI_PORT_DUAL_MEDIA_COPPER_PREFERRED:
      return lemming::dataplane::sai::PORT_DUAL_MEDIA_COPPER_PREFERRED;

    case SAI_PORT_DUAL_MEDIA_FIBER_PREFERRED:
      return lemming::dataplane::sai::PORT_DUAL_MEDIA_FIBER_PREFERRED;

    default:
      return lemming::dataplane::sai::PORT_DUAL_MEDIA_UNSPECIFIED;
  }
}
sai_port_dual_media_t convert_sai_port_dual_media_t_to_sai(
    lemming::dataplane::sai::PortDualMedia val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_DUAL_MEDIA_NONE:
      return SAI_PORT_DUAL_MEDIA_NONE;

    case lemming::dataplane::sai::PORT_DUAL_MEDIA_COPPER_ONLY:
      return SAI_PORT_DUAL_MEDIA_COPPER_ONLY;

    case lemming::dataplane::sai::PORT_DUAL_MEDIA_FIBER_ONLY:
      return SAI_PORT_DUAL_MEDIA_FIBER_ONLY;

    case lemming::dataplane::sai::PORT_DUAL_MEDIA_COPPER_PREFERRED:
      return SAI_PORT_DUAL_MEDIA_COPPER_PREFERRED;

    case lemming::dataplane::sai::PORT_DUAL_MEDIA_FIBER_PREFERRED:
      return SAI_PORT_DUAL_MEDIA_FIBER_PREFERRED;

    default:
      return SAI_PORT_DUAL_MEDIA_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_dual_media_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_dual_media_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_dual_media_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_dual_media_t_to_sai(
        static_cast<lemming::dataplane::sai::PortDualMedia>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortErrStatus convert_sai_port_err_status_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_ERR_STATUS_DATA_UNIT_CRC_ERROR:
      return lemming::dataplane::sai::PORT_ERR_STATUS_DATA_UNIT_CRC_ERROR;

    case SAI_PORT_ERR_STATUS_DATA_UNIT_SIZE:
      return lemming::dataplane::sai::PORT_ERR_STATUS_DATA_UNIT_SIZE;

    case SAI_PORT_ERR_STATUS_DATA_UNIT_MISALIGNMENT_ERROR:
      return lemming::dataplane::sai::
          PORT_ERR_STATUS_DATA_UNIT_MISALIGNMENT_ERROR;

    case SAI_PORT_ERR_STATUS_CODE_GROUP_ERROR:
      return lemming::dataplane::sai::PORT_ERR_STATUS_CODE_GROUP_ERROR;

    case SAI_PORT_ERR_STATUS_SIGNAL_LOCAL_ERROR:
      return lemming::dataplane::sai::PORT_ERR_STATUS_SIGNAL_LOCAL_ERROR;

    case SAI_PORT_ERR_STATUS_NO_RX_REACHABILITY:
      return lemming::dataplane::sai::PORT_ERR_STATUS_NO_RX_REACHABILITY;

    case SAI_PORT_ERR_STATUS_CRC_RATE:
      return lemming::dataplane::sai::PORT_ERR_STATUS_CRC_RATE;

    case SAI_PORT_ERR_STATUS_REMOTE_FAULT_STATUS:
      return lemming::dataplane::sai::PORT_ERR_STATUS_REMOTE_FAULT_STATUS;

    case SAI_PORT_ERR_STATUS_MAX:
      return lemming::dataplane::sai::PORT_ERR_STATUS_MAX;

    default:
      return lemming::dataplane::sai::PORT_ERR_STATUS_UNSPECIFIED;
  }
}
sai_port_err_status_t convert_sai_port_err_status_t_to_sai(
    lemming::dataplane::sai::PortErrStatus val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_ERR_STATUS_DATA_UNIT_CRC_ERROR:
      return SAI_PORT_ERR_STATUS_DATA_UNIT_CRC_ERROR;

    case lemming::dataplane::sai::PORT_ERR_STATUS_DATA_UNIT_SIZE:
      return SAI_PORT_ERR_STATUS_DATA_UNIT_SIZE;

    case lemming::dataplane::sai::PORT_ERR_STATUS_DATA_UNIT_MISALIGNMENT_ERROR:
      return SAI_PORT_ERR_STATUS_DATA_UNIT_MISALIGNMENT_ERROR;

    case lemming::dataplane::sai::PORT_ERR_STATUS_CODE_GROUP_ERROR:
      return SAI_PORT_ERR_STATUS_CODE_GROUP_ERROR;

    case lemming::dataplane::sai::PORT_ERR_STATUS_SIGNAL_LOCAL_ERROR:
      return SAI_PORT_ERR_STATUS_SIGNAL_LOCAL_ERROR;

    case lemming::dataplane::sai::PORT_ERR_STATUS_NO_RX_REACHABILITY:
      return SAI_PORT_ERR_STATUS_NO_RX_REACHABILITY;

    case lemming::dataplane::sai::PORT_ERR_STATUS_CRC_RATE:
      return SAI_PORT_ERR_STATUS_CRC_RATE;

    case lemming::dataplane::sai::PORT_ERR_STATUS_REMOTE_FAULT_STATUS:
      return SAI_PORT_ERR_STATUS_REMOTE_FAULT_STATUS;

    case lemming::dataplane::sai::PORT_ERR_STATUS_MAX:
      return SAI_PORT_ERR_STATUS_MAX;

    default:
      return SAI_PORT_ERR_STATUS_DATA_UNIT_CRC_ERROR;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_err_status_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_err_status_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_err_status_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_err_status_t_to_sai(
        static_cast<lemming::dataplane::sai::PortErrStatus>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortFecModeExtended
convert_sai_port_fec_mode_extended_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_FEC_MODE_EXTENDED_NONE:
      return lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_NONE;

    case SAI_PORT_FEC_MODE_EXTENDED_RS528:
      return lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_RS528;

    case SAI_PORT_FEC_MODE_EXTENDED_RS544:
      return lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_RS544;

    case SAI_PORT_FEC_MODE_EXTENDED_RS544_INTERLEAVED:
      return lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_RS544_INTERLEAVED;

    case SAI_PORT_FEC_MODE_EXTENDED_FC:
      return lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_FC;

    default:
      return lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_UNSPECIFIED;
  }
}
sai_port_fec_mode_extended_t convert_sai_port_fec_mode_extended_t_to_sai(
    lemming::dataplane::sai::PortFecModeExtended val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_NONE:
      return SAI_PORT_FEC_MODE_EXTENDED_NONE;

    case lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_RS528:
      return SAI_PORT_FEC_MODE_EXTENDED_RS528;

    case lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_RS544:
      return SAI_PORT_FEC_MODE_EXTENDED_RS544;

    case lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_RS544_INTERLEAVED:
      return SAI_PORT_FEC_MODE_EXTENDED_RS544_INTERLEAVED;

    case lemming::dataplane::sai::PORT_FEC_MODE_EXTENDED_FC:
      return SAI_PORT_FEC_MODE_EXTENDED_FC;

    default:
      return SAI_PORT_FEC_MODE_EXTENDED_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_fec_mode_extended_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_fec_mode_extended_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_fec_mode_extended_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_fec_mode_extended_t_to_sai(
        static_cast<lemming::dataplane::sai::PortFecModeExtended>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortFecMode convert_sai_port_fec_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_FEC_MODE_NONE:
      return lemming::dataplane::sai::PORT_FEC_MODE_NONE;

    case SAI_PORT_FEC_MODE_RS:
      return lemming::dataplane::sai::PORT_FEC_MODE_RS;

    case SAI_PORT_FEC_MODE_FC:
      return lemming::dataplane::sai::PORT_FEC_MODE_FC;

    default:
      return lemming::dataplane::sai::PORT_FEC_MODE_UNSPECIFIED;
  }
}
sai_port_fec_mode_t convert_sai_port_fec_mode_t_to_sai(
    lemming::dataplane::sai::PortFecMode val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_FEC_MODE_NONE:
      return SAI_PORT_FEC_MODE_NONE;

    case lemming::dataplane::sai::PORT_FEC_MODE_RS:
      return SAI_PORT_FEC_MODE_RS;

    case lemming::dataplane::sai::PORT_FEC_MODE_FC:
      return SAI_PORT_FEC_MODE_FC;

    default:
      return SAI_PORT_FEC_MODE_NONE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_port_fec_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_fec_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_fec_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_fec_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::PortFecMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortFlowControlMode
convert_sai_port_flow_control_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_FLOW_CONTROL_MODE_DISABLE:
      return lemming::dataplane::sai::PORT_FLOW_CONTROL_MODE_DISABLE;

    case SAI_PORT_FLOW_CONTROL_MODE_TX_ONLY:
      return lemming::dataplane::sai::PORT_FLOW_CONTROL_MODE_TX_ONLY;

    case SAI_PORT_FLOW_CONTROL_MODE_RX_ONLY:
      return lemming::dataplane::sai::PORT_FLOW_CONTROL_MODE_RX_ONLY;

    case SAI_PORT_FLOW_CONTROL_MODE_BOTH_ENABLE:
      return lemming::dataplane::sai::PORT_FLOW_CONTROL_MODE_BOTH_ENABLE;

    default:
      return lemming::dataplane::sai::PORT_FLOW_CONTROL_MODE_UNSPECIFIED;
  }
}
sai_port_flow_control_mode_t convert_sai_port_flow_control_mode_t_to_sai(
    lemming::dataplane::sai::PortFlowControlMode val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_FLOW_CONTROL_MODE_DISABLE:
      return SAI_PORT_FLOW_CONTROL_MODE_DISABLE;

    case lemming::dataplane::sai::PORT_FLOW_CONTROL_MODE_TX_ONLY:
      return SAI_PORT_FLOW_CONTROL_MODE_TX_ONLY;

    case lemming::dataplane::sai::PORT_FLOW_CONTROL_MODE_RX_ONLY:
      return SAI_PORT_FLOW_CONTROL_MODE_RX_ONLY;

    case lemming::dataplane::sai::PORT_FLOW_CONTROL_MODE_BOTH_ENABLE:
      return SAI_PORT_FLOW_CONTROL_MODE_BOTH_ENABLE;

    default:
      return SAI_PORT_FLOW_CONTROL_MODE_DISABLE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_flow_control_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_flow_control_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_flow_control_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_flow_control_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::PortFlowControlMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortInterfaceType
convert_sai_port_interface_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_INTERFACE_TYPE_NONE:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_NONE;

    case SAI_PORT_INTERFACE_TYPE_CR:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_CR;

    case SAI_PORT_INTERFACE_TYPE_CR2:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_CR2;

    case SAI_PORT_INTERFACE_TYPE_CR4:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_CR4;

    case SAI_PORT_INTERFACE_TYPE_SR:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_SR;

    case SAI_PORT_INTERFACE_TYPE_SR2:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_SR2;

    case SAI_PORT_INTERFACE_TYPE_SR4:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_SR4;

    case SAI_PORT_INTERFACE_TYPE_LR:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_LR;

    case SAI_PORT_INTERFACE_TYPE_LR4:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_LR4;

    case SAI_PORT_INTERFACE_TYPE_KR:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_KR;

    case SAI_PORT_INTERFACE_TYPE_KR4:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_KR4;

    case SAI_PORT_INTERFACE_TYPE_CAUI:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_CAUI;

    case SAI_PORT_INTERFACE_TYPE_GMII:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_GMII;

    case SAI_PORT_INTERFACE_TYPE_SFI:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_SFI;

    case SAI_PORT_INTERFACE_TYPE_XLAUI:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_XLAUI;

    case SAI_PORT_INTERFACE_TYPE_KR2:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_KR2;

    case SAI_PORT_INTERFACE_TYPE_CAUI4:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_CAUI4;

    case SAI_PORT_INTERFACE_TYPE_XAUI:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_XAUI;

    case SAI_PORT_INTERFACE_TYPE_XFI:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_XFI;

    case SAI_PORT_INTERFACE_TYPE_XGMII:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_XGMII;

    case SAI_PORT_INTERFACE_TYPE_CR8:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_CR8;

    case SAI_PORT_INTERFACE_TYPE_KR8:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_KR8;

    case SAI_PORT_INTERFACE_TYPE_SR8:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_SR8;

    case SAI_PORT_INTERFACE_TYPE_LR8:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_LR8;

    case SAI_PORT_INTERFACE_TYPE_MAX:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_MAX;

    default:
      return lemming::dataplane::sai::PORT_INTERFACE_TYPE_UNSPECIFIED;
  }
}
sai_port_interface_type_t convert_sai_port_interface_type_t_to_sai(
    lemming::dataplane::sai::PortInterfaceType val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_NONE:
      return SAI_PORT_INTERFACE_TYPE_NONE;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_CR:
      return SAI_PORT_INTERFACE_TYPE_CR;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_CR2:
      return SAI_PORT_INTERFACE_TYPE_CR2;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_CR4:
      return SAI_PORT_INTERFACE_TYPE_CR4;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_SR:
      return SAI_PORT_INTERFACE_TYPE_SR;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_SR2:
      return SAI_PORT_INTERFACE_TYPE_SR2;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_SR4:
      return SAI_PORT_INTERFACE_TYPE_SR4;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_LR:
      return SAI_PORT_INTERFACE_TYPE_LR;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_LR4:
      return SAI_PORT_INTERFACE_TYPE_LR4;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_KR:
      return SAI_PORT_INTERFACE_TYPE_KR;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_KR4:
      return SAI_PORT_INTERFACE_TYPE_KR4;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_CAUI:
      return SAI_PORT_INTERFACE_TYPE_CAUI;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_GMII:
      return SAI_PORT_INTERFACE_TYPE_GMII;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_SFI:
      return SAI_PORT_INTERFACE_TYPE_SFI;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_XLAUI:
      return SAI_PORT_INTERFACE_TYPE_XLAUI;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_KR2:
      return SAI_PORT_INTERFACE_TYPE_KR2;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_CAUI4:
      return SAI_PORT_INTERFACE_TYPE_CAUI4;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_XAUI:
      return SAI_PORT_INTERFACE_TYPE_XAUI;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_XFI:
      return SAI_PORT_INTERFACE_TYPE_XFI;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_XGMII:
      return SAI_PORT_INTERFACE_TYPE_XGMII;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_CR8:
      return SAI_PORT_INTERFACE_TYPE_CR8;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_KR8:
      return SAI_PORT_INTERFACE_TYPE_KR8;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_SR8:
      return SAI_PORT_INTERFACE_TYPE_SR8;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_LR8:
      return SAI_PORT_INTERFACE_TYPE_LR8;

    case lemming::dataplane::sai::PORT_INTERFACE_TYPE_MAX:
      return SAI_PORT_INTERFACE_TYPE_MAX;

    default:
      return SAI_PORT_INTERFACE_TYPE_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_interface_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_interface_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_interface_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_interface_type_t_to_sai(
        static_cast<lemming::dataplane::sai::PortInterfaceType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortInternalLoopbackMode
convert_sai_port_internal_loopback_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_INTERNAL_LOOPBACK_MODE_NONE:
      return lemming::dataplane::sai::PORT_INTERNAL_LOOPBACK_MODE_NONE;

    case SAI_PORT_INTERNAL_LOOPBACK_MODE_PHY:
      return lemming::dataplane::sai::PORT_INTERNAL_LOOPBACK_MODE_PHY;

    case SAI_PORT_INTERNAL_LOOPBACK_MODE_MAC:
      return lemming::dataplane::sai::PORT_INTERNAL_LOOPBACK_MODE_MAC;

    default:
      return lemming::dataplane::sai::PORT_INTERNAL_LOOPBACK_MODE_UNSPECIFIED;
  }
}
sai_port_internal_loopback_mode_t
convert_sai_port_internal_loopback_mode_t_to_sai(
    lemming::dataplane::sai::PortInternalLoopbackMode val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_INTERNAL_LOOPBACK_MODE_NONE:
      return SAI_PORT_INTERNAL_LOOPBACK_MODE_NONE;

    case lemming::dataplane::sai::PORT_INTERNAL_LOOPBACK_MODE_PHY:
      return SAI_PORT_INTERNAL_LOOPBACK_MODE_PHY;

    case lemming::dataplane::sai::PORT_INTERNAL_LOOPBACK_MODE_MAC:
      return SAI_PORT_INTERNAL_LOOPBACK_MODE_MAC;

    default:
      return SAI_PORT_INTERNAL_LOOPBACK_MODE_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_internal_loopback_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_port_internal_loopback_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_internal_loopback_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_internal_loopback_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::PortInternalLoopbackMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortLinkTrainingFailureStatus
convert_sai_port_link_training_failure_status_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_LINK_TRAINING_FAILURE_STATUS_NO_ERROR:
      return lemming::dataplane::sai::
          PORT_LINK_TRAINING_FAILURE_STATUS_NO_ERROR;

    case SAI_PORT_LINK_TRAINING_FAILURE_STATUS_FRAME_LOCK_ERROR:
      return lemming::dataplane::sai::
          PORT_LINK_TRAINING_FAILURE_STATUS_FRAME_LOCK_ERROR;

    case SAI_PORT_LINK_TRAINING_FAILURE_STATUS_SNR_LOWER_THRESHOLD:
      return lemming::dataplane::sai::
          PORT_LINK_TRAINING_FAILURE_STATUS_SNR_LOWER_THRESHOLD;

    case SAI_PORT_LINK_TRAINING_FAILURE_STATUS_TIME_OUT:
      return lemming::dataplane::sai::
          PORT_LINK_TRAINING_FAILURE_STATUS_TIME_OUT;

    default:
      return lemming::dataplane::sai::
          PORT_LINK_TRAINING_FAILURE_STATUS_UNSPECIFIED;
  }
}
sai_port_link_training_failure_status_t
convert_sai_port_link_training_failure_status_t_to_sai(
    lemming::dataplane::sai::PortLinkTrainingFailureStatus val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_LINK_TRAINING_FAILURE_STATUS_NO_ERROR:
      return SAI_PORT_LINK_TRAINING_FAILURE_STATUS_NO_ERROR;

    case lemming::dataplane::sai::
        PORT_LINK_TRAINING_FAILURE_STATUS_FRAME_LOCK_ERROR:
      return SAI_PORT_LINK_TRAINING_FAILURE_STATUS_FRAME_LOCK_ERROR;

    case lemming::dataplane::sai::
        PORT_LINK_TRAINING_FAILURE_STATUS_SNR_LOWER_THRESHOLD:
      return SAI_PORT_LINK_TRAINING_FAILURE_STATUS_SNR_LOWER_THRESHOLD;

    case lemming::dataplane::sai::PORT_LINK_TRAINING_FAILURE_STATUS_TIME_OUT:
      return SAI_PORT_LINK_TRAINING_FAILURE_STATUS_TIME_OUT;

    default:
      return SAI_PORT_LINK_TRAINING_FAILURE_STATUS_NO_ERROR;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_link_training_failure_status_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_port_link_training_failure_status_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_link_training_failure_status_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_link_training_failure_status_t_to_sai(
        static_cast<lemming::dataplane::sai::PortLinkTrainingFailureStatus>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortLinkTrainingRxStatus
convert_sai_port_link_training_rx_status_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_LINK_TRAINING_RX_STATUS_NOT_TRAINED:
      return lemming::dataplane::sai::PORT_LINK_TRAINING_RX_STATUS_NOT_TRAINED;

    case SAI_PORT_LINK_TRAINING_RX_STATUS_TRAINED:
      return lemming::dataplane::sai::PORT_LINK_TRAINING_RX_STATUS_TRAINED;

    default:
      return lemming::dataplane::sai::PORT_LINK_TRAINING_RX_STATUS_UNSPECIFIED;
  }
}
sai_port_link_training_rx_status_t
convert_sai_port_link_training_rx_status_t_to_sai(
    lemming::dataplane::sai::PortLinkTrainingRxStatus val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_LINK_TRAINING_RX_STATUS_NOT_TRAINED:
      return SAI_PORT_LINK_TRAINING_RX_STATUS_NOT_TRAINED;

    case lemming::dataplane::sai::PORT_LINK_TRAINING_RX_STATUS_TRAINED:
      return SAI_PORT_LINK_TRAINING_RX_STATUS_TRAINED;

    default:
      return SAI_PORT_LINK_TRAINING_RX_STATUS_NOT_TRAINED;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_link_training_rx_status_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_port_link_training_rx_status_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_link_training_rx_status_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_link_training_rx_status_t_to_sai(
        static_cast<lemming::dataplane::sai::PortLinkTrainingRxStatus>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortLoopbackMode
convert_sai_port_loopback_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_LOOPBACK_MODE_NONE:
      return lemming::dataplane::sai::PORT_LOOPBACK_MODE_NONE;

    case SAI_PORT_LOOPBACK_MODE_PHY:
      return lemming::dataplane::sai::PORT_LOOPBACK_MODE_PHY;

    case SAI_PORT_LOOPBACK_MODE_MAC:
      return lemming::dataplane::sai::PORT_LOOPBACK_MODE_MAC;

    case SAI_PORT_LOOPBACK_MODE_PHY_REMOTE:
      return lemming::dataplane::sai::PORT_LOOPBACK_MODE_PHY_REMOTE;

    default:
      return lemming::dataplane::sai::PORT_LOOPBACK_MODE_UNSPECIFIED;
  }
}
sai_port_loopback_mode_t convert_sai_port_loopback_mode_t_to_sai(
    lemming::dataplane::sai::PortLoopbackMode val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_LOOPBACK_MODE_NONE:
      return SAI_PORT_LOOPBACK_MODE_NONE;

    case lemming::dataplane::sai::PORT_LOOPBACK_MODE_PHY:
      return SAI_PORT_LOOPBACK_MODE_PHY;

    case lemming::dataplane::sai::PORT_LOOPBACK_MODE_MAC:
      return SAI_PORT_LOOPBACK_MODE_MAC;

    case lemming::dataplane::sai::PORT_LOOPBACK_MODE_PHY_REMOTE:
      return SAI_PORT_LOOPBACK_MODE_PHY_REMOTE;

    default:
      return SAI_PORT_LOOPBACK_MODE_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_loopback_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_loopback_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_loopback_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_loopback_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::PortLoopbackMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortMdixModeConfig
convert_sai_port_mdix_mode_config_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_MDIX_MODE_CONFIG_AUTO:
      return lemming::dataplane::sai::PORT_MDIX_MODE_CONFIG_AUTO;

    case SAI_PORT_MDIX_MODE_CONFIG_STRAIGHT:
      return lemming::dataplane::sai::PORT_MDIX_MODE_CONFIG_STRAIGHT;

    case SAI_PORT_MDIX_MODE_CONFIG_CROSSOVER:
      return lemming::dataplane::sai::PORT_MDIX_MODE_CONFIG_CROSSOVER;

    default:
      return lemming::dataplane::sai::PORT_MDIX_MODE_CONFIG_UNSPECIFIED;
  }
}
sai_port_mdix_mode_config_t convert_sai_port_mdix_mode_config_t_to_sai(
    lemming::dataplane::sai::PortMdixModeConfig val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_MDIX_MODE_CONFIG_AUTO:
      return SAI_PORT_MDIX_MODE_CONFIG_AUTO;

    case lemming::dataplane::sai::PORT_MDIX_MODE_CONFIG_STRAIGHT:
      return SAI_PORT_MDIX_MODE_CONFIG_STRAIGHT;

    case lemming::dataplane::sai::PORT_MDIX_MODE_CONFIG_CROSSOVER:
      return SAI_PORT_MDIX_MODE_CONFIG_CROSSOVER;

    default:
      return SAI_PORT_MDIX_MODE_CONFIG_AUTO;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_mdix_mode_config_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_mdix_mode_config_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_mdix_mode_config_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_mdix_mode_config_t_to_sai(
        static_cast<lemming::dataplane::sai::PortMdixModeConfig>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortMdixModeStatus
convert_sai_port_mdix_mode_status_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_MDIX_MODE_STATUS_STRAIGHT:
      return lemming::dataplane::sai::PORT_MDIX_MODE_STATUS_STRAIGHT;

    case SAI_PORT_MDIX_MODE_STATUS_CROSSOVER:
      return lemming::dataplane::sai::PORT_MDIX_MODE_STATUS_CROSSOVER;

    default:
      return lemming::dataplane::sai::PORT_MDIX_MODE_STATUS_UNSPECIFIED;
  }
}
sai_port_mdix_mode_status_t convert_sai_port_mdix_mode_status_t_to_sai(
    lemming::dataplane::sai::PortMdixModeStatus val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_MDIX_MODE_STATUS_STRAIGHT:
      return SAI_PORT_MDIX_MODE_STATUS_STRAIGHT;

    case lemming::dataplane::sai::PORT_MDIX_MODE_STATUS_CROSSOVER:
      return SAI_PORT_MDIX_MODE_STATUS_CROSSOVER;

    default:
      return SAI_PORT_MDIX_MODE_STATUS_STRAIGHT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_mdix_mode_status_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_mdix_mode_status_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_mdix_mode_status_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_mdix_mode_status_t_to_sai(
        static_cast<lemming::dataplane::sai::PortMdixModeStatus>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortMediaType convert_sai_port_media_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_MEDIA_TYPE_NOT_PRESENT:
      return lemming::dataplane::sai::PORT_MEDIA_TYPE_NOT_PRESENT;

    case SAI_PORT_MEDIA_TYPE_UNKNOWN:
      return lemming::dataplane::sai::PORT_MEDIA_TYPE_UNKNOWN;

    case SAI_PORT_MEDIA_TYPE_FIBER:
      return lemming::dataplane::sai::PORT_MEDIA_TYPE_FIBER;

    case SAI_PORT_MEDIA_TYPE_COPPER:
      return lemming::dataplane::sai::PORT_MEDIA_TYPE_COPPER;

    case SAI_PORT_MEDIA_TYPE_BACKPLANE:
      return lemming::dataplane::sai::PORT_MEDIA_TYPE_BACKPLANE;

    default:
      return lemming::dataplane::sai::PORT_MEDIA_TYPE_UNSPECIFIED;
  }
}
sai_port_media_type_t convert_sai_port_media_type_t_to_sai(
    lemming::dataplane::sai::PortMediaType val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_MEDIA_TYPE_NOT_PRESENT:
      return SAI_PORT_MEDIA_TYPE_NOT_PRESENT;

    case lemming::dataplane::sai::PORT_MEDIA_TYPE_UNKNOWN:
      return SAI_PORT_MEDIA_TYPE_UNKNOWN;

    case lemming::dataplane::sai::PORT_MEDIA_TYPE_FIBER:
      return SAI_PORT_MEDIA_TYPE_FIBER;

    case lemming::dataplane::sai::PORT_MEDIA_TYPE_COPPER:
      return SAI_PORT_MEDIA_TYPE_COPPER;

    case lemming::dataplane::sai::PORT_MEDIA_TYPE_BACKPLANE:
      return SAI_PORT_MEDIA_TYPE_BACKPLANE;

    default:
      return SAI_PORT_MEDIA_TYPE_NOT_PRESENT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_media_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_media_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_media_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_media_type_t_to_sai(
        static_cast<lemming::dataplane::sai::PortMediaType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortModuleType convert_sai_port_module_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_MODULE_TYPE_1000BASE_X:
      return lemming::dataplane::sai::PORT_MODULE_TYPE_1000BASE_X;

    case SAI_PORT_MODULE_TYPE_100FX:
      return lemming::dataplane::sai::PORT_MODULE_TYPE_100FX;

    case SAI_PORT_MODULE_TYPE_SGMII_SLAVE:
      return lemming::dataplane::sai::PORT_MODULE_TYPE_SGMII_SLAVE;

    default:
      return lemming::dataplane::sai::PORT_MODULE_TYPE_UNSPECIFIED;
  }
}
sai_port_module_type_t convert_sai_port_module_type_t_to_sai(
    lemming::dataplane::sai::PortModuleType val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_MODULE_TYPE_1000BASE_X:
      return SAI_PORT_MODULE_TYPE_1000BASE_X;

    case lemming::dataplane::sai::PORT_MODULE_TYPE_100FX:
      return SAI_PORT_MODULE_TYPE_100FX;

    case lemming::dataplane::sai::PORT_MODULE_TYPE_SGMII_SLAVE:
      return SAI_PORT_MODULE_TYPE_SGMII_SLAVE;

    default:
      return SAI_PORT_MODULE_TYPE_1000BASE_X;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_module_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_module_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_module_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_module_type_t_to_sai(
        static_cast<lemming::dataplane::sai::PortModuleType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortOperStatus convert_sai_port_oper_status_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_OPER_STATUS_UNKNOWN:
      return lemming::dataplane::sai::PORT_OPER_STATUS_UNKNOWN;

    case SAI_PORT_OPER_STATUS_UP:
      return lemming::dataplane::sai::PORT_OPER_STATUS_UP;

    case SAI_PORT_OPER_STATUS_DOWN:
      return lemming::dataplane::sai::PORT_OPER_STATUS_DOWN;

    case SAI_PORT_OPER_STATUS_TESTING:
      return lemming::dataplane::sai::PORT_OPER_STATUS_TESTING;

    case SAI_PORT_OPER_STATUS_NOT_PRESENT:
      return lemming::dataplane::sai::PORT_OPER_STATUS_NOT_PRESENT;

    default:
      return lemming::dataplane::sai::PORT_OPER_STATUS_UNSPECIFIED;
  }
}
sai_port_oper_status_t convert_sai_port_oper_status_t_to_sai(
    lemming::dataplane::sai::PortOperStatus val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_OPER_STATUS_UNKNOWN:
      return SAI_PORT_OPER_STATUS_UNKNOWN;

    case lemming::dataplane::sai::PORT_OPER_STATUS_UP:
      return SAI_PORT_OPER_STATUS_UP;

    case lemming::dataplane::sai::PORT_OPER_STATUS_DOWN:
      return SAI_PORT_OPER_STATUS_DOWN;

    case lemming::dataplane::sai::PORT_OPER_STATUS_TESTING:
      return SAI_PORT_OPER_STATUS_TESTING;

    case lemming::dataplane::sai::PORT_OPER_STATUS_NOT_PRESENT:
      return SAI_PORT_OPER_STATUS_NOT_PRESENT;

    default:
      return SAI_PORT_OPER_STATUS_UNKNOWN;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_oper_status_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_oper_status_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_oper_status_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_oper_status_t_to_sai(
        static_cast<lemming::dataplane::sai::PortOperStatus>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortPoolAttr convert_sai_port_pool_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_POOL_ATTR_PORT_ID:
      return lemming::dataplane::sai::PORT_POOL_ATTR_PORT_ID;

    case SAI_PORT_POOL_ATTR_BUFFER_POOL_ID:
      return lemming::dataplane::sai::PORT_POOL_ATTR_BUFFER_POOL_ID;

    case SAI_PORT_POOL_ATTR_QOS_WRED_PROFILE_ID:
      return lemming::dataplane::sai::PORT_POOL_ATTR_QOS_WRED_PROFILE_ID;

    default:
      return lemming::dataplane::sai::PORT_POOL_ATTR_UNSPECIFIED;
  }
}
sai_port_pool_attr_t convert_sai_port_pool_attr_t_to_sai(
    lemming::dataplane::sai::PortPoolAttr val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_POOL_ATTR_PORT_ID:
      return SAI_PORT_POOL_ATTR_PORT_ID;

    case lemming::dataplane::sai::PORT_POOL_ATTR_BUFFER_POOL_ID:
      return SAI_PORT_POOL_ATTR_BUFFER_POOL_ID;

    case lemming::dataplane::sai::PORT_POOL_ATTR_QOS_WRED_PROFILE_ID:
      return SAI_PORT_POOL_ATTR_QOS_WRED_PROFILE_ID;

    default:
      return SAI_PORT_POOL_ATTR_PORT_ID;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_port_pool_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_pool_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_pool_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_pool_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::PortPoolAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortPoolStat convert_sai_port_pool_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_POOL_STAT_IF_OCTETS:
      return lemming::dataplane::sai::PORT_POOL_STAT_IF_OCTETS;

    case SAI_PORT_POOL_STAT_GREEN_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::PORT_POOL_STAT_GREEN_WRED_DROPPED_PACKETS;

    case SAI_PORT_POOL_STAT_GREEN_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::PORT_POOL_STAT_GREEN_WRED_DROPPED_BYTES;

    case SAI_PORT_POOL_STAT_YELLOW_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::
          PORT_POOL_STAT_YELLOW_WRED_DROPPED_PACKETS;

    case SAI_PORT_POOL_STAT_YELLOW_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::PORT_POOL_STAT_YELLOW_WRED_DROPPED_BYTES;

    case SAI_PORT_POOL_STAT_RED_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::PORT_POOL_STAT_RED_WRED_DROPPED_PACKETS;

    case SAI_PORT_POOL_STAT_RED_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::PORT_POOL_STAT_RED_WRED_DROPPED_BYTES;

    case SAI_PORT_POOL_STAT_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::PORT_POOL_STAT_WRED_DROPPED_PACKETS;

    case SAI_PORT_POOL_STAT_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::PORT_POOL_STAT_WRED_DROPPED_BYTES;

    case SAI_PORT_POOL_STAT_GREEN_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::
          PORT_POOL_STAT_GREEN_WRED_ECN_MARKED_PACKETS;

    case SAI_PORT_POOL_STAT_GREEN_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::
          PORT_POOL_STAT_GREEN_WRED_ECN_MARKED_BYTES;

    case SAI_PORT_POOL_STAT_YELLOW_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::
          PORT_POOL_STAT_YELLOW_WRED_ECN_MARKED_PACKETS;

    case SAI_PORT_POOL_STAT_YELLOW_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::
          PORT_POOL_STAT_YELLOW_WRED_ECN_MARKED_BYTES;

    case SAI_PORT_POOL_STAT_RED_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::
          PORT_POOL_STAT_RED_WRED_ECN_MARKED_PACKETS;

    case SAI_PORT_POOL_STAT_RED_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::PORT_POOL_STAT_RED_WRED_ECN_MARKED_BYTES;

    case SAI_PORT_POOL_STAT_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::PORT_POOL_STAT_WRED_ECN_MARKED_PACKETS;

    case SAI_PORT_POOL_STAT_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::PORT_POOL_STAT_WRED_ECN_MARKED_BYTES;

    case SAI_PORT_POOL_STAT_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::PORT_POOL_STAT_CURR_OCCUPANCY_BYTES;

    case SAI_PORT_POOL_STAT_WATERMARK_BYTES:
      return lemming::dataplane::sai::PORT_POOL_STAT_WATERMARK_BYTES;

    case SAI_PORT_POOL_STAT_SHARED_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::
          PORT_POOL_STAT_SHARED_CURR_OCCUPANCY_BYTES;

    case SAI_PORT_POOL_STAT_SHARED_WATERMARK_BYTES:
      return lemming::dataplane::sai::PORT_POOL_STAT_SHARED_WATERMARK_BYTES;

    case SAI_PORT_POOL_STAT_DROPPED_PKTS:
      return lemming::dataplane::sai::PORT_POOL_STAT_DROPPED_PKTS;

    default:
      return lemming::dataplane::sai::PORT_POOL_STAT_UNSPECIFIED;
  }
}
sai_port_pool_stat_t convert_sai_port_pool_stat_t_to_sai(
    lemming::dataplane::sai::PortPoolStat val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_POOL_STAT_IF_OCTETS:
      return SAI_PORT_POOL_STAT_IF_OCTETS;

    case lemming::dataplane::sai::PORT_POOL_STAT_GREEN_WRED_DROPPED_PACKETS:
      return SAI_PORT_POOL_STAT_GREEN_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::PORT_POOL_STAT_GREEN_WRED_DROPPED_BYTES:
      return SAI_PORT_POOL_STAT_GREEN_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_YELLOW_WRED_DROPPED_PACKETS:
      return SAI_PORT_POOL_STAT_YELLOW_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::PORT_POOL_STAT_YELLOW_WRED_DROPPED_BYTES:
      return SAI_PORT_POOL_STAT_YELLOW_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_RED_WRED_DROPPED_PACKETS:
      return SAI_PORT_POOL_STAT_RED_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::PORT_POOL_STAT_RED_WRED_DROPPED_BYTES:
      return SAI_PORT_POOL_STAT_RED_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_WRED_DROPPED_PACKETS:
      return SAI_PORT_POOL_STAT_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::PORT_POOL_STAT_WRED_DROPPED_BYTES:
      return SAI_PORT_POOL_STAT_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_GREEN_WRED_ECN_MARKED_PACKETS:
      return SAI_PORT_POOL_STAT_GREEN_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::PORT_POOL_STAT_GREEN_WRED_ECN_MARKED_BYTES:
      return SAI_PORT_POOL_STAT_GREEN_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_YELLOW_WRED_ECN_MARKED_PACKETS:
      return SAI_PORT_POOL_STAT_YELLOW_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::PORT_POOL_STAT_YELLOW_WRED_ECN_MARKED_BYTES:
      return SAI_PORT_POOL_STAT_YELLOW_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_RED_WRED_ECN_MARKED_PACKETS:
      return SAI_PORT_POOL_STAT_RED_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::PORT_POOL_STAT_RED_WRED_ECN_MARKED_BYTES:
      return SAI_PORT_POOL_STAT_RED_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_WRED_ECN_MARKED_PACKETS:
      return SAI_PORT_POOL_STAT_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::PORT_POOL_STAT_WRED_ECN_MARKED_BYTES:
      return SAI_PORT_POOL_STAT_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_CURR_OCCUPANCY_BYTES:
      return SAI_PORT_POOL_STAT_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_WATERMARK_BYTES:
      return SAI_PORT_POOL_STAT_WATERMARK_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_SHARED_CURR_OCCUPANCY_BYTES:
      return SAI_PORT_POOL_STAT_SHARED_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_SHARED_WATERMARK_BYTES:
      return SAI_PORT_POOL_STAT_SHARED_WATERMARK_BYTES;

    case lemming::dataplane::sai::PORT_POOL_STAT_DROPPED_PKTS:
      return SAI_PORT_POOL_STAT_DROPPED_PKTS;

    default:
      return SAI_PORT_POOL_STAT_IF_OCTETS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_port_pool_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_pool_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_pool_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_pool_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::PortPoolStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortPrbsConfig convert_sai_port_prbs_config_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_PRBS_CONFIG_DISABLE:
      return lemming::dataplane::sai::PORT_PRBS_CONFIG_DISABLE;

    case SAI_PORT_PRBS_CONFIG_ENABLE_TX_RX:
      return lemming::dataplane::sai::PORT_PRBS_CONFIG_ENABLE_TX_RX;

    case SAI_PORT_PRBS_CONFIG_ENABLE_RX:
      return lemming::dataplane::sai::PORT_PRBS_CONFIG_ENABLE_RX;

    case SAI_PORT_PRBS_CONFIG_ENABLE_TX:
      return lemming::dataplane::sai::PORT_PRBS_CONFIG_ENABLE_TX;

    default:
      return lemming::dataplane::sai::PORT_PRBS_CONFIG_UNSPECIFIED;
  }
}
sai_port_prbs_config_t convert_sai_port_prbs_config_t_to_sai(
    lemming::dataplane::sai::PortPrbsConfig val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_PRBS_CONFIG_DISABLE:
      return SAI_PORT_PRBS_CONFIG_DISABLE;

    case lemming::dataplane::sai::PORT_PRBS_CONFIG_ENABLE_TX_RX:
      return SAI_PORT_PRBS_CONFIG_ENABLE_TX_RX;

    case lemming::dataplane::sai::PORT_PRBS_CONFIG_ENABLE_RX:
      return SAI_PORT_PRBS_CONFIG_ENABLE_RX;

    case lemming::dataplane::sai::PORT_PRBS_CONFIG_ENABLE_TX:
      return SAI_PORT_PRBS_CONFIG_ENABLE_TX;

    default:
      return SAI_PORT_PRBS_CONFIG_DISABLE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_prbs_config_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_prbs_config_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_prbs_config_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_prbs_config_t_to_sai(
        static_cast<lemming::dataplane::sai::PortPrbsConfig>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortPrbsRxStatus
convert_sai_port_prbs_rx_status_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_PRBS_RX_STATUS_OK:
      return lemming::dataplane::sai::PORT_PRBS_RX_STATUS_OK;

    case SAI_PORT_PRBS_RX_STATUS_LOCK_WITH_ERRORS:
      return lemming::dataplane::sai::PORT_PRBS_RX_STATUS_LOCK_WITH_ERRORS;

    case SAI_PORT_PRBS_RX_STATUS_NOT_LOCKED:
      return lemming::dataplane::sai::PORT_PRBS_RX_STATUS_NOT_LOCKED;

    case SAI_PORT_PRBS_RX_STATUS_LOST_LOCK:
      return lemming::dataplane::sai::PORT_PRBS_RX_STATUS_LOST_LOCK;

    default:
      return lemming::dataplane::sai::PORT_PRBS_RX_STATUS_UNSPECIFIED;
  }
}
sai_port_prbs_rx_status_t convert_sai_port_prbs_rx_status_t_to_sai(
    lemming::dataplane::sai::PortPrbsRxStatus val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_PRBS_RX_STATUS_OK:
      return SAI_PORT_PRBS_RX_STATUS_OK;

    case lemming::dataplane::sai::PORT_PRBS_RX_STATUS_LOCK_WITH_ERRORS:
      return SAI_PORT_PRBS_RX_STATUS_LOCK_WITH_ERRORS;

    case lemming::dataplane::sai::PORT_PRBS_RX_STATUS_NOT_LOCKED:
      return SAI_PORT_PRBS_RX_STATUS_NOT_LOCKED;

    case lemming::dataplane::sai::PORT_PRBS_RX_STATUS_LOST_LOCK:
      return SAI_PORT_PRBS_RX_STATUS_LOST_LOCK;

    default:
      return SAI_PORT_PRBS_RX_STATUS_OK;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_prbs_rx_status_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_prbs_rx_status_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_prbs_rx_status_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_prbs_rx_status_t_to_sai(
        static_cast<lemming::dataplane::sai::PortPrbsRxStatus>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortPriorityFlowControlMode
convert_sai_port_priority_flow_control_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_PRIORITY_FLOW_CONTROL_MODE_COMBINED:
      return lemming::dataplane::sai::PORT_PRIORITY_FLOW_CONTROL_MODE_COMBINED;

    case SAI_PORT_PRIORITY_FLOW_CONTROL_MODE_SEPARATE:
      return lemming::dataplane::sai::PORT_PRIORITY_FLOW_CONTROL_MODE_SEPARATE;

    default:
      return lemming::dataplane::sai::
          PORT_PRIORITY_FLOW_CONTROL_MODE_UNSPECIFIED;
  }
}
sai_port_priority_flow_control_mode_t
convert_sai_port_priority_flow_control_mode_t_to_sai(
    lemming::dataplane::sai::PortPriorityFlowControlMode val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_PRIORITY_FLOW_CONTROL_MODE_COMBINED:
      return SAI_PORT_PRIORITY_FLOW_CONTROL_MODE_COMBINED;

    case lemming::dataplane::sai::PORT_PRIORITY_FLOW_CONTROL_MODE_SEPARATE:
      return SAI_PORT_PRIORITY_FLOW_CONTROL_MODE_SEPARATE;

    default:
      return SAI_PORT_PRIORITY_FLOW_CONTROL_MODE_COMBINED;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_priority_flow_control_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_port_priority_flow_control_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_priority_flow_control_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_priority_flow_control_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::PortPriorityFlowControlMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortPtpMode convert_sai_port_ptp_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_PTP_MODE_NONE:
      return lemming::dataplane::sai::PORT_PTP_MODE_NONE;

    case SAI_PORT_PTP_MODE_SINGLE_STEP_TIMESTAMP:
      return lemming::dataplane::sai::PORT_PTP_MODE_SINGLE_STEP_TIMESTAMP;

    case SAI_PORT_PTP_MODE_TWO_STEP_TIMESTAMP:
      return lemming::dataplane::sai::PORT_PTP_MODE_TWO_STEP_TIMESTAMP;

    default:
      return lemming::dataplane::sai::PORT_PTP_MODE_UNSPECIFIED;
  }
}
sai_port_ptp_mode_t convert_sai_port_ptp_mode_t_to_sai(
    lemming::dataplane::sai::PortPtpMode val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_PTP_MODE_NONE:
      return SAI_PORT_PTP_MODE_NONE;

    case lemming::dataplane::sai::PORT_PTP_MODE_SINGLE_STEP_TIMESTAMP:
      return SAI_PORT_PTP_MODE_SINGLE_STEP_TIMESTAMP;

    case lemming::dataplane::sai::PORT_PTP_MODE_TWO_STEP_TIMESTAMP:
      return SAI_PORT_PTP_MODE_TWO_STEP_TIMESTAMP;

    default:
      return SAI_PORT_PTP_MODE_NONE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_port_ptp_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_ptp_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_ptp_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_ptp_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::PortPtpMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortSerdesAttr convert_sai_port_serdes_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_SERDES_ATTR_PORT_ID:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_PORT_ID;

    case SAI_PORT_SERDES_ATTR_PREEMPHASIS:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_PREEMPHASIS;

    case SAI_PORT_SERDES_ATTR_IDRIVER:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_IDRIVER;

    case SAI_PORT_SERDES_ATTR_IPREDRIVER:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_IPREDRIVER;

    case SAI_PORT_SERDES_ATTR_TX_FIR_PRE1:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_PRE1;

    case SAI_PORT_SERDES_ATTR_TX_FIR_PRE2:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_PRE2;

    case SAI_PORT_SERDES_ATTR_TX_FIR_PRE3:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_PRE3;

    case SAI_PORT_SERDES_ATTR_TX_FIR_MAIN:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_MAIN;

    case SAI_PORT_SERDES_ATTR_TX_FIR_POST1:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_POST1;

    case SAI_PORT_SERDES_ATTR_TX_FIR_POST2:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_POST2;

    case SAI_PORT_SERDES_ATTR_TX_FIR_POST3:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_POST3;

    case SAI_PORT_SERDES_ATTR_TX_FIR_ATTN:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_ATTN;

    default:
      return lemming::dataplane::sai::PORT_SERDES_ATTR_UNSPECIFIED;
  }
}
sai_port_serdes_attr_t convert_sai_port_serdes_attr_t_to_sai(
    lemming::dataplane::sai::PortSerdesAttr val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_SERDES_ATTR_PORT_ID:
      return SAI_PORT_SERDES_ATTR_PORT_ID;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_PREEMPHASIS:
      return SAI_PORT_SERDES_ATTR_PREEMPHASIS;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_IDRIVER:
      return SAI_PORT_SERDES_ATTR_IDRIVER;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_IPREDRIVER:
      return SAI_PORT_SERDES_ATTR_IPREDRIVER;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_PRE1:
      return SAI_PORT_SERDES_ATTR_TX_FIR_PRE1;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_PRE2:
      return SAI_PORT_SERDES_ATTR_TX_FIR_PRE2;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_PRE3:
      return SAI_PORT_SERDES_ATTR_TX_FIR_PRE3;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_MAIN:
      return SAI_PORT_SERDES_ATTR_TX_FIR_MAIN;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_POST1:
      return SAI_PORT_SERDES_ATTR_TX_FIR_POST1;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_POST2:
      return SAI_PORT_SERDES_ATTR_TX_FIR_POST2;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_POST3:
      return SAI_PORT_SERDES_ATTR_TX_FIR_POST3;

    case lemming::dataplane::sai::PORT_SERDES_ATTR_TX_FIR_ATTN:
      return SAI_PORT_SERDES_ATTR_TX_FIR_ATTN;

    default:
      return SAI_PORT_SERDES_ATTR_PORT_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_port_serdes_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_serdes_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_serdes_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_serdes_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::PortSerdesAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortStat convert_sai_port_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_STAT_IF_IN_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_OCTETS;

    case SAI_PORT_STAT_IF_IN_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_UCAST_PKTS;

    case SAI_PORT_STAT_IF_IN_NON_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_NON_UCAST_PKTS;

    case SAI_PORT_STAT_IF_IN_DISCARDS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_DISCARDS;

    case SAI_PORT_STAT_IF_IN_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_ERRORS;

    case SAI_PORT_STAT_IF_IN_UNKNOWN_PROTOS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_UNKNOWN_PROTOS;

    case SAI_PORT_STAT_IF_IN_BROADCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_BROADCAST_PKTS;

    case SAI_PORT_STAT_IF_IN_MULTICAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_MULTICAST_PKTS;

    case SAI_PORT_STAT_IF_IN_VLAN_DISCARDS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_VLAN_DISCARDS;

    case SAI_PORT_STAT_IF_OUT_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_IF_OUT_OCTETS;

    case SAI_PORT_STAT_IF_OUT_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IF_OUT_UCAST_PKTS;

    case SAI_PORT_STAT_IF_OUT_NON_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IF_OUT_NON_UCAST_PKTS;

    case SAI_PORT_STAT_IF_OUT_DISCARDS:
      return lemming::dataplane::sai::PORT_STAT_IF_OUT_DISCARDS;

    case SAI_PORT_STAT_IF_OUT_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_IF_OUT_ERRORS;

    case SAI_PORT_STAT_IF_OUT_QLEN:
      return lemming::dataplane::sai::PORT_STAT_IF_OUT_QLEN;

    case SAI_PORT_STAT_IF_OUT_BROADCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IF_OUT_BROADCAST_PKTS;

    case SAI_PORT_STAT_IF_OUT_MULTICAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IF_OUT_MULTICAST_PKTS;

    case SAI_PORT_STAT_ETHER_STATS_DROP_EVENTS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_DROP_EVENTS;

    case SAI_PORT_STAT_ETHER_STATS_MULTICAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_MULTICAST_PKTS;

    case SAI_PORT_STAT_ETHER_STATS_BROADCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_BROADCAST_PKTS;

    case SAI_PORT_STAT_ETHER_STATS_UNDERSIZE_PKTS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_UNDERSIZE_PKTS;

    case SAI_PORT_STAT_ETHER_STATS_FRAGMENTS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_FRAGMENTS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS_64_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_PKTS_64_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS_65_TO_127_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_STATS_PKTS_65_TO_127_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS_128_TO_255_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_STATS_PKTS_128_TO_255_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS_256_TO_511_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_STATS_PKTS_256_TO_511_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS_512_TO_1023_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_STATS_PKTS_512_TO_1023_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS_1024_TO_1518_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_STATS_PKTS_1024_TO_1518_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS_1519_TO_2047_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_STATS_PKTS_1519_TO_2047_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS_2048_TO_4095_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_STATS_PKTS_2048_TO_4095_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS_4096_TO_9216_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_STATS_PKTS_4096_TO_9216_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS_9217_TO_16383_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_STATS_PKTS_9217_TO_16383_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_OVERSIZE_PKTS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_OVERSIZE_PKTS;

    case SAI_PORT_STAT_ETHER_RX_OVERSIZE_PKTS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_RX_OVERSIZE_PKTS;

    case SAI_PORT_STAT_ETHER_TX_OVERSIZE_PKTS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_TX_OVERSIZE_PKTS;

    case SAI_PORT_STAT_ETHER_STATS_JABBERS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_JABBERS;

    case SAI_PORT_STAT_ETHER_STATS_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_OCTETS;

    case SAI_PORT_STAT_ETHER_STATS_PKTS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_PKTS;

    case SAI_PORT_STAT_ETHER_STATS_COLLISIONS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_COLLISIONS;

    case SAI_PORT_STAT_ETHER_STATS_CRC_ALIGN_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_CRC_ALIGN_ERRORS;

    case SAI_PORT_STAT_ETHER_STATS_TX_NO_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_TX_NO_ERRORS;

    case SAI_PORT_STAT_ETHER_STATS_RX_NO_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_STATS_RX_NO_ERRORS;

    case SAI_PORT_STAT_IP_IN_RECEIVES:
      return lemming::dataplane::sai::PORT_STAT_IP_IN_RECEIVES;

    case SAI_PORT_STAT_IP_IN_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_IP_IN_OCTETS;

    case SAI_PORT_STAT_IP_IN_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IP_IN_UCAST_PKTS;

    case SAI_PORT_STAT_IP_IN_NON_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IP_IN_NON_UCAST_PKTS;

    case SAI_PORT_STAT_IP_IN_DISCARDS:
      return lemming::dataplane::sai::PORT_STAT_IP_IN_DISCARDS;

    case SAI_PORT_STAT_IP_OUT_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_IP_OUT_OCTETS;

    case SAI_PORT_STAT_IP_OUT_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IP_OUT_UCAST_PKTS;

    case SAI_PORT_STAT_IP_OUT_NON_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IP_OUT_NON_UCAST_PKTS;

    case SAI_PORT_STAT_IP_OUT_DISCARDS:
      return lemming::dataplane::sai::PORT_STAT_IP_OUT_DISCARDS;

    case SAI_PORT_STAT_IPV6_IN_RECEIVES:
      return lemming::dataplane::sai::PORT_STAT_IPV6_IN_RECEIVES;

    case SAI_PORT_STAT_IPV6_IN_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_IPV6_IN_OCTETS;

    case SAI_PORT_STAT_IPV6_IN_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IPV6_IN_UCAST_PKTS;

    case SAI_PORT_STAT_IPV6_IN_NON_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IPV6_IN_NON_UCAST_PKTS;

    case SAI_PORT_STAT_IPV6_IN_MCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IPV6_IN_MCAST_PKTS;

    case SAI_PORT_STAT_IPV6_IN_DISCARDS:
      return lemming::dataplane::sai::PORT_STAT_IPV6_IN_DISCARDS;

    case SAI_PORT_STAT_IPV6_OUT_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_IPV6_OUT_OCTETS;

    case SAI_PORT_STAT_IPV6_OUT_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IPV6_OUT_UCAST_PKTS;

    case SAI_PORT_STAT_IPV6_OUT_NON_UCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IPV6_OUT_NON_UCAST_PKTS;

    case SAI_PORT_STAT_IPV6_OUT_MCAST_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IPV6_OUT_MCAST_PKTS;

    case SAI_PORT_STAT_IPV6_OUT_DISCARDS:
      return lemming::dataplane::sai::PORT_STAT_IPV6_OUT_DISCARDS;

    case SAI_PORT_STAT_GREEN_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::PORT_STAT_GREEN_WRED_DROPPED_PACKETS;

    case SAI_PORT_STAT_GREEN_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::PORT_STAT_GREEN_WRED_DROPPED_BYTES;

    case SAI_PORT_STAT_YELLOW_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::PORT_STAT_YELLOW_WRED_DROPPED_PACKETS;

    case SAI_PORT_STAT_YELLOW_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::PORT_STAT_YELLOW_WRED_DROPPED_BYTES;

    case SAI_PORT_STAT_RED_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::PORT_STAT_RED_WRED_DROPPED_PACKETS;

    case SAI_PORT_STAT_RED_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::PORT_STAT_RED_WRED_DROPPED_BYTES;

    case SAI_PORT_STAT_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::PORT_STAT_WRED_DROPPED_PACKETS;

    case SAI_PORT_STAT_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::PORT_STAT_WRED_DROPPED_BYTES;

    case SAI_PORT_STAT_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::PORT_STAT_ECN_MARKED_PACKETS;

    case SAI_PORT_STAT_ETHER_IN_PKTS_64_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_64_OCTETS;

    case SAI_PORT_STAT_ETHER_IN_PKTS_65_TO_127_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_65_TO_127_OCTETS;

    case SAI_PORT_STAT_ETHER_IN_PKTS_128_TO_255_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_128_TO_255_OCTETS;

    case SAI_PORT_STAT_ETHER_IN_PKTS_256_TO_511_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_256_TO_511_OCTETS;

    case SAI_PORT_STAT_ETHER_IN_PKTS_512_TO_1023_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_IN_PKTS_512_TO_1023_OCTETS;

    case SAI_PORT_STAT_ETHER_IN_PKTS_1024_TO_1518_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_IN_PKTS_1024_TO_1518_OCTETS;

    case SAI_PORT_STAT_ETHER_IN_PKTS_1519_TO_2047_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_IN_PKTS_1519_TO_2047_OCTETS;

    case SAI_PORT_STAT_ETHER_IN_PKTS_2048_TO_4095_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_IN_PKTS_2048_TO_4095_OCTETS;

    case SAI_PORT_STAT_ETHER_IN_PKTS_4096_TO_9216_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_IN_PKTS_4096_TO_9216_OCTETS;

    case SAI_PORT_STAT_ETHER_IN_PKTS_9217_TO_16383_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_IN_PKTS_9217_TO_16383_OCTETS;

    case SAI_PORT_STAT_ETHER_OUT_PKTS_64_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_64_OCTETS;

    case SAI_PORT_STAT_ETHER_OUT_PKTS_65_TO_127_OCTETS:
      return lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_65_TO_127_OCTETS;

    case SAI_PORT_STAT_ETHER_OUT_PKTS_128_TO_255_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_OUT_PKTS_128_TO_255_OCTETS;

    case SAI_PORT_STAT_ETHER_OUT_PKTS_256_TO_511_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_OUT_PKTS_256_TO_511_OCTETS;

    case SAI_PORT_STAT_ETHER_OUT_PKTS_512_TO_1023_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_OUT_PKTS_512_TO_1023_OCTETS;

    case SAI_PORT_STAT_ETHER_OUT_PKTS_1024_TO_1518_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_OUT_PKTS_1024_TO_1518_OCTETS;

    case SAI_PORT_STAT_ETHER_OUT_PKTS_1519_TO_2047_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_OUT_PKTS_1519_TO_2047_OCTETS;

    case SAI_PORT_STAT_ETHER_OUT_PKTS_2048_TO_4095_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_OUT_PKTS_2048_TO_4095_OCTETS;

    case SAI_PORT_STAT_ETHER_OUT_PKTS_4096_TO_9216_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_OUT_PKTS_4096_TO_9216_OCTETS;

    case SAI_PORT_STAT_ETHER_OUT_PKTS_9217_TO_16383_OCTETS:
      return lemming::dataplane::sai::
          PORT_STAT_ETHER_OUT_PKTS_9217_TO_16383_OCTETS;

    case SAI_PORT_STAT_IN_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::PORT_STAT_IN_CURR_OCCUPANCY_BYTES;

    case SAI_PORT_STAT_IN_WATERMARK_BYTES:
      return lemming::dataplane::sai::PORT_STAT_IN_WATERMARK_BYTES;

    case SAI_PORT_STAT_IN_SHARED_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::PORT_STAT_IN_SHARED_CURR_OCCUPANCY_BYTES;

    case SAI_PORT_STAT_IN_SHARED_WATERMARK_BYTES:
      return lemming::dataplane::sai::PORT_STAT_IN_SHARED_WATERMARK_BYTES;

    case SAI_PORT_STAT_OUT_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::PORT_STAT_OUT_CURR_OCCUPANCY_BYTES;

    case SAI_PORT_STAT_OUT_WATERMARK_BYTES:
      return lemming::dataplane::sai::PORT_STAT_OUT_WATERMARK_BYTES;

    case SAI_PORT_STAT_OUT_SHARED_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::PORT_STAT_OUT_SHARED_CURR_OCCUPANCY_BYTES;

    case SAI_PORT_STAT_OUT_SHARED_WATERMARK_BYTES:
      return lemming::dataplane::sai::PORT_STAT_OUT_SHARED_WATERMARK_BYTES;

    case SAI_PORT_STAT_IN_DROPPED_PKTS:
      return lemming::dataplane::sai::PORT_STAT_IN_DROPPED_PKTS;

    case SAI_PORT_STAT_OUT_DROPPED_PKTS:
      return lemming::dataplane::sai::PORT_STAT_OUT_DROPPED_PKTS;

    case SAI_PORT_STAT_PAUSE_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PAUSE_RX_PKTS;

    case SAI_PORT_STAT_PAUSE_TX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PAUSE_TX_PKTS;

    case SAI_PORT_STAT_PFC_0_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_0_RX_PKTS;

    case SAI_PORT_STAT_PFC_0_TX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_0_TX_PKTS;

    case SAI_PORT_STAT_PFC_1_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_1_RX_PKTS;

    case SAI_PORT_STAT_PFC_1_TX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_1_TX_PKTS;

    case SAI_PORT_STAT_PFC_2_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_2_RX_PKTS;

    case SAI_PORT_STAT_PFC_2_TX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_2_TX_PKTS;

    case SAI_PORT_STAT_PFC_3_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_3_RX_PKTS;

    case SAI_PORT_STAT_PFC_3_TX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_3_TX_PKTS;

    case SAI_PORT_STAT_PFC_4_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_4_RX_PKTS;

    case SAI_PORT_STAT_PFC_4_TX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_4_TX_PKTS;

    case SAI_PORT_STAT_PFC_5_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_5_RX_PKTS;

    case SAI_PORT_STAT_PFC_5_TX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_5_TX_PKTS;

    case SAI_PORT_STAT_PFC_6_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_6_RX_PKTS;

    case SAI_PORT_STAT_PFC_6_TX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_6_TX_PKTS;

    case SAI_PORT_STAT_PFC_7_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_7_RX_PKTS;

    case SAI_PORT_STAT_PFC_7_TX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_7_TX_PKTS;

    case SAI_PORT_STAT_PFC_0_RX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_0_RX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_0_TX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_0_TX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_1_RX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_1_RX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_1_TX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_1_TX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_2_RX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_2_RX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_2_TX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_2_TX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_3_RX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_3_RX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_3_TX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_3_TX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_4_RX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_4_RX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_4_TX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_4_TX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_5_RX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_5_RX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_5_TX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_5_TX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_6_RX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_6_RX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_6_TX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_6_TX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_7_RX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_7_RX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_7_TX_PAUSE_DURATION:
      return lemming::dataplane::sai::PORT_STAT_PFC_7_TX_PAUSE_DURATION;

    case SAI_PORT_STAT_PFC_0_RX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_0_RX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_0_TX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_0_TX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_1_RX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_1_RX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_1_TX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_1_TX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_2_RX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_2_RX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_2_TX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_2_TX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_3_RX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_3_RX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_3_TX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_3_TX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_4_RX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_4_RX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_4_TX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_4_TX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_5_RX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_5_RX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_5_TX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_5_TX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_6_RX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_6_RX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_6_TX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_6_TX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_7_RX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_7_RX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_7_TX_PAUSE_DURATION_US:
      return lemming::dataplane::sai::PORT_STAT_PFC_7_TX_PAUSE_DURATION_US;

    case SAI_PORT_STAT_PFC_0_ON2OFF_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_0_ON2OFF_RX_PKTS;

    case SAI_PORT_STAT_PFC_1_ON2OFF_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_1_ON2OFF_RX_PKTS;

    case SAI_PORT_STAT_PFC_2_ON2OFF_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_2_ON2OFF_RX_PKTS;

    case SAI_PORT_STAT_PFC_3_ON2OFF_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_3_ON2OFF_RX_PKTS;

    case SAI_PORT_STAT_PFC_4_ON2OFF_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_4_ON2OFF_RX_PKTS;

    case SAI_PORT_STAT_PFC_5_ON2OFF_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_5_ON2OFF_RX_PKTS;

    case SAI_PORT_STAT_PFC_6_ON2OFF_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_6_ON2OFF_RX_PKTS;

    case SAI_PORT_STAT_PFC_7_ON2OFF_RX_PKTS:
      return lemming::dataplane::sai::PORT_STAT_PFC_7_ON2OFF_RX_PKTS;

    case SAI_PORT_STAT_DOT3_STATS_ALIGNMENT_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_DOT3_STATS_ALIGNMENT_ERRORS;

    case SAI_PORT_STAT_DOT3_STATS_FCS_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_DOT3_STATS_FCS_ERRORS;

    case SAI_PORT_STAT_DOT3_STATS_SINGLE_COLLISION_FRAMES:
      return lemming::dataplane::sai::
          PORT_STAT_DOT3_STATS_SINGLE_COLLISION_FRAMES;

    case SAI_PORT_STAT_DOT3_STATS_MULTIPLE_COLLISION_FRAMES:
      return lemming::dataplane::sai::
          PORT_STAT_DOT3_STATS_MULTIPLE_COLLISION_FRAMES;

    case SAI_PORT_STAT_DOT3_STATS_SQE_TEST_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_DOT3_STATS_SQE_TEST_ERRORS;

    case SAI_PORT_STAT_DOT3_STATS_DEFERRED_TRANSMISSIONS:
      return lemming::dataplane::sai::
          PORT_STAT_DOT3_STATS_DEFERRED_TRANSMISSIONS;

    case SAI_PORT_STAT_DOT3_STATS_LATE_COLLISIONS:
      return lemming::dataplane::sai::PORT_STAT_DOT3_STATS_LATE_COLLISIONS;

    case SAI_PORT_STAT_DOT3_STATS_EXCESSIVE_COLLISIONS:
      return lemming::dataplane::sai::PORT_STAT_DOT3_STATS_EXCESSIVE_COLLISIONS;

    case SAI_PORT_STAT_DOT3_STATS_INTERNAL_MAC_TRANSMIT_ERRORS:
      return lemming::dataplane::sai::
          PORT_STAT_DOT3_STATS_INTERNAL_MAC_TRANSMIT_ERRORS;

    case SAI_PORT_STAT_DOT3_STATS_CARRIER_SENSE_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_DOT3_STATS_CARRIER_SENSE_ERRORS;

    case SAI_PORT_STAT_DOT3_STATS_FRAME_TOO_LONGS:
      return lemming::dataplane::sai::PORT_STAT_DOT3_STATS_FRAME_TOO_LONGS;

    case SAI_PORT_STAT_DOT3_STATS_INTERNAL_MAC_RECEIVE_ERRORS:
      return lemming::dataplane::sai::
          PORT_STAT_DOT3_STATS_INTERNAL_MAC_RECEIVE_ERRORS;

    case SAI_PORT_STAT_DOT3_STATS_SYMBOL_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_DOT3_STATS_SYMBOL_ERRORS;

    case SAI_PORT_STAT_DOT3_CONTROL_IN_UNKNOWN_OPCODES:
      return lemming::dataplane::sai::PORT_STAT_DOT3_CONTROL_IN_UNKNOWN_OPCODES;

    case SAI_PORT_STAT_EEE_TX_EVENT_COUNT:
      return lemming::dataplane::sai::PORT_STAT_EEE_TX_EVENT_COUNT;

    case SAI_PORT_STAT_EEE_RX_EVENT_COUNT:
      return lemming::dataplane::sai::PORT_STAT_EEE_RX_EVENT_COUNT;

    case SAI_PORT_STAT_EEE_TX_DURATION:
      return lemming::dataplane::sai::PORT_STAT_EEE_TX_DURATION;

    case SAI_PORT_STAT_EEE_RX_DURATION:
      return lemming::dataplane::sai::PORT_STAT_EEE_RX_DURATION;

    case SAI_PORT_STAT_PRBS_ERROR_COUNT:
      return lemming::dataplane::sai::PORT_STAT_PRBS_ERROR_COUNT;

    case SAI_PORT_STAT_IF_IN_FEC_CORRECTABLE_FRAMES:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CORRECTABLE_FRAMES;

    case SAI_PORT_STAT_IF_IN_FEC_NOT_CORRECTABLE_FRAMES:
      return lemming::dataplane::sai::
          PORT_STAT_IF_IN_FEC_NOT_CORRECTABLE_FRAMES;

    case SAI_PORT_STAT_IF_IN_FEC_SYMBOL_ERRORS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_SYMBOL_ERRORS;

    case SAI_PORT_STAT_IF_IN_FABRIC_DATA_UNITS:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FABRIC_DATA_UNITS;

    case SAI_PORT_STAT_IF_OUT_FABRIC_DATA_UNITS:
      return lemming::dataplane::sai::PORT_STAT_IF_OUT_FABRIC_DATA_UNITS;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S0:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S0;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S1:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S1;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S2:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S2;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S3:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S3;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S4:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S4;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S5:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S5;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S6:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S6;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S7:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S7;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S8:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S8;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S9:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S9;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S10:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S10;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S11:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S11;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S12:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S12;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S13:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S13;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S14:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S14;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S15:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S15;

    case SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S16:
      return lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S16;

    case SAI_PORT_STAT_IN_DROP_REASON_RANGE_BASE:
      return lemming::dataplane::sai::PORT_STAT_IN_DROP_REASON_RANGE_BASE;

    case SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_IN_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS;

    case SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_IN_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS;

    case SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_IN_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS;

    case SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_IN_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS;

    case SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_IN_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS;

    case SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_IN_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS;

    case SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_IN_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS;

    case SAI_PORT_STAT_IN_DROP_REASON_RANGE_END:
      return lemming::dataplane::sai::PORT_STAT_IN_DROP_REASON_RANGE_END;

    case SAI_PORT_STAT_OUT_DROP_REASON_RANGE_BASE:
      return lemming::dataplane::sai::PORT_STAT_OUT_DROP_REASON_RANGE_BASE;

    case SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_OUT_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS;

    case SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_OUT_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS;

    case SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_OUT_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS;

    case SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_OUT_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS;

    case SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_OUT_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS;

    case SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_OUT_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS;

    case SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS:
      return lemming::dataplane::sai::
          PORT_STAT_OUT_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS;

    case SAI_PORT_STAT_OUT_DROP_REASON_RANGE_END:
      return lemming::dataplane::sai::PORT_STAT_OUT_DROP_REASON_RANGE_END;

    default:
      return lemming::dataplane::sai::PORT_STAT_UNSPECIFIED;
  }
}
sai_port_stat_t convert_sai_port_stat_t_to_sai(
    lemming::dataplane::sai::PortStat val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_STAT_IF_IN_OCTETS:
      return SAI_PORT_STAT_IF_IN_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_UCAST_PKTS:
      return SAI_PORT_STAT_IF_IN_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_NON_UCAST_PKTS:
      return SAI_PORT_STAT_IF_IN_NON_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_DISCARDS:
      return SAI_PORT_STAT_IF_IN_DISCARDS;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_ERRORS:
      return SAI_PORT_STAT_IF_IN_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_UNKNOWN_PROTOS:
      return SAI_PORT_STAT_IF_IN_UNKNOWN_PROTOS;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_BROADCAST_PKTS:
      return SAI_PORT_STAT_IF_IN_BROADCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_MULTICAST_PKTS:
      return SAI_PORT_STAT_IF_IN_MULTICAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_VLAN_DISCARDS:
      return SAI_PORT_STAT_IF_IN_VLAN_DISCARDS;

    case lemming::dataplane::sai::PORT_STAT_IF_OUT_OCTETS:
      return SAI_PORT_STAT_IF_OUT_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_IF_OUT_UCAST_PKTS:
      return SAI_PORT_STAT_IF_OUT_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IF_OUT_NON_UCAST_PKTS:
      return SAI_PORT_STAT_IF_OUT_NON_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IF_OUT_DISCARDS:
      return SAI_PORT_STAT_IF_OUT_DISCARDS;

    case lemming::dataplane::sai::PORT_STAT_IF_OUT_ERRORS:
      return SAI_PORT_STAT_IF_OUT_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_IF_OUT_QLEN:
      return SAI_PORT_STAT_IF_OUT_QLEN;

    case lemming::dataplane::sai::PORT_STAT_IF_OUT_BROADCAST_PKTS:
      return SAI_PORT_STAT_IF_OUT_BROADCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IF_OUT_MULTICAST_PKTS:
      return SAI_PORT_STAT_IF_OUT_MULTICAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_DROP_EVENTS:
      return SAI_PORT_STAT_ETHER_STATS_DROP_EVENTS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_MULTICAST_PKTS:
      return SAI_PORT_STAT_ETHER_STATS_MULTICAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_BROADCAST_PKTS:
      return SAI_PORT_STAT_ETHER_STATS_BROADCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_UNDERSIZE_PKTS:
      return SAI_PORT_STAT_ETHER_STATS_UNDERSIZE_PKTS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_FRAGMENTS:
      return SAI_PORT_STAT_ETHER_STATS_FRAGMENTS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_PKTS_64_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS_64_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_PKTS_65_TO_127_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS_65_TO_127_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_PKTS_128_TO_255_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS_128_TO_255_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_PKTS_256_TO_511_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS_256_TO_511_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_PKTS_512_TO_1023_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS_512_TO_1023_OCTETS;

    case lemming::dataplane::sai::
        PORT_STAT_ETHER_STATS_PKTS_1024_TO_1518_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS_1024_TO_1518_OCTETS;

    case lemming::dataplane::sai::
        PORT_STAT_ETHER_STATS_PKTS_1519_TO_2047_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS_1519_TO_2047_OCTETS;

    case lemming::dataplane::sai::
        PORT_STAT_ETHER_STATS_PKTS_2048_TO_4095_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS_2048_TO_4095_OCTETS;

    case lemming::dataplane::sai::
        PORT_STAT_ETHER_STATS_PKTS_4096_TO_9216_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS_4096_TO_9216_OCTETS;

    case lemming::dataplane::sai::
        PORT_STAT_ETHER_STATS_PKTS_9217_TO_16383_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS_9217_TO_16383_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_OVERSIZE_PKTS:
      return SAI_PORT_STAT_ETHER_STATS_OVERSIZE_PKTS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_RX_OVERSIZE_PKTS:
      return SAI_PORT_STAT_ETHER_RX_OVERSIZE_PKTS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_TX_OVERSIZE_PKTS:
      return SAI_PORT_STAT_ETHER_TX_OVERSIZE_PKTS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_JABBERS:
      return SAI_PORT_STAT_ETHER_STATS_JABBERS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_OCTETS:
      return SAI_PORT_STAT_ETHER_STATS_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_PKTS:
      return SAI_PORT_STAT_ETHER_STATS_PKTS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_COLLISIONS:
      return SAI_PORT_STAT_ETHER_STATS_COLLISIONS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_CRC_ALIGN_ERRORS:
      return SAI_PORT_STAT_ETHER_STATS_CRC_ALIGN_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_TX_NO_ERRORS:
      return SAI_PORT_STAT_ETHER_STATS_TX_NO_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_STATS_RX_NO_ERRORS:
      return SAI_PORT_STAT_ETHER_STATS_RX_NO_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_IP_IN_RECEIVES:
      return SAI_PORT_STAT_IP_IN_RECEIVES;

    case lemming::dataplane::sai::PORT_STAT_IP_IN_OCTETS:
      return SAI_PORT_STAT_IP_IN_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_IP_IN_UCAST_PKTS:
      return SAI_PORT_STAT_IP_IN_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IP_IN_NON_UCAST_PKTS:
      return SAI_PORT_STAT_IP_IN_NON_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IP_IN_DISCARDS:
      return SAI_PORT_STAT_IP_IN_DISCARDS;

    case lemming::dataplane::sai::PORT_STAT_IP_OUT_OCTETS:
      return SAI_PORT_STAT_IP_OUT_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_IP_OUT_UCAST_PKTS:
      return SAI_PORT_STAT_IP_OUT_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IP_OUT_NON_UCAST_PKTS:
      return SAI_PORT_STAT_IP_OUT_NON_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IP_OUT_DISCARDS:
      return SAI_PORT_STAT_IP_OUT_DISCARDS;

    case lemming::dataplane::sai::PORT_STAT_IPV6_IN_RECEIVES:
      return SAI_PORT_STAT_IPV6_IN_RECEIVES;

    case lemming::dataplane::sai::PORT_STAT_IPV6_IN_OCTETS:
      return SAI_PORT_STAT_IPV6_IN_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_IPV6_IN_UCAST_PKTS:
      return SAI_PORT_STAT_IPV6_IN_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IPV6_IN_NON_UCAST_PKTS:
      return SAI_PORT_STAT_IPV6_IN_NON_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IPV6_IN_MCAST_PKTS:
      return SAI_PORT_STAT_IPV6_IN_MCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IPV6_IN_DISCARDS:
      return SAI_PORT_STAT_IPV6_IN_DISCARDS;

    case lemming::dataplane::sai::PORT_STAT_IPV6_OUT_OCTETS:
      return SAI_PORT_STAT_IPV6_OUT_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_IPV6_OUT_UCAST_PKTS:
      return SAI_PORT_STAT_IPV6_OUT_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IPV6_OUT_NON_UCAST_PKTS:
      return SAI_PORT_STAT_IPV6_OUT_NON_UCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IPV6_OUT_MCAST_PKTS:
      return SAI_PORT_STAT_IPV6_OUT_MCAST_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IPV6_OUT_DISCARDS:
      return SAI_PORT_STAT_IPV6_OUT_DISCARDS;

    case lemming::dataplane::sai::PORT_STAT_GREEN_WRED_DROPPED_PACKETS:
      return SAI_PORT_STAT_GREEN_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::PORT_STAT_GREEN_WRED_DROPPED_BYTES:
      return SAI_PORT_STAT_GREEN_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::PORT_STAT_YELLOW_WRED_DROPPED_PACKETS:
      return SAI_PORT_STAT_YELLOW_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::PORT_STAT_YELLOW_WRED_DROPPED_BYTES:
      return SAI_PORT_STAT_YELLOW_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::PORT_STAT_RED_WRED_DROPPED_PACKETS:
      return SAI_PORT_STAT_RED_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::PORT_STAT_RED_WRED_DROPPED_BYTES:
      return SAI_PORT_STAT_RED_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::PORT_STAT_WRED_DROPPED_PACKETS:
      return SAI_PORT_STAT_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::PORT_STAT_WRED_DROPPED_BYTES:
      return SAI_PORT_STAT_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::PORT_STAT_ECN_MARKED_PACKETS:
      return SAI_PORT_STAT_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_64_OCTETS:
      return SAI_PORT_STAT_ETHER_IN_PKTS_64_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_65_TO_127_OCTETS:
      return SAI_PORT_STAT_ETHER_IN_PKTS_65_TO_127_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_128_TO_255_OCTETS:
      return SAI_PORT_STAT_ETHER_IN_PKTS_128_TO_255_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_256_TO_511_OCTETS:
      return SAI_PORT_STAT_ETHER_IN_PKTS_256_TO_511_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_512_TO_1023_OCTETS:
      return SAI_PORT_STAT_ETHER_IN_PKTS_512_TO_1023_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_1024_TO_1518_OCTETS:
      return SAI_PORT_STAT_ETHER_IN_PKTS_1024_TO_1518_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_1519_TO_2047_OCTETS:
      return SAI_PORT_STAT_ETHER_IN_PKTS_1519_TO_2047_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_2048_TO_4095_OCTETS:
      return SAI_PORT_STAT_ETHER_IN_PKTS_2048_TO_4095_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_4096_TO_9216_OCTETS:
      return SAI_PORT_STAT_ETHER_IN_PKTS_4096_TO_9216_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_IN_PKTS_9217_TO_16383_OCTETS:
      return SAI_PORT_STAT_ETHER_IN_PKTS_9217_TO_16383_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_64_OCTETS:
      return SAI_PORT_STAT_ETHER_OUT_PKTS_64_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_65_TO_127_OCTETS:
      return SAI_PORT_STAT_ETHER_OUT_PKTS_65_TO_127_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_128_TO_255_OCTETS:
      return SAI_PORT_STAT_ETHER_OUT_PKTS_128_TO_255_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_256_TO_511_OCTETS:
      return SAI_PORT_STAT_ETHER_OUT_PKTS_256_TO_511_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_512_TO_1023_OCTETS:
      return SAI_PORT_STAT_ETHER_OUT_PKTS_512_TO_1023_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_1024_TO_1518_OCTETS:
      return SAI_PORT_STAT_ETHER_OUT_PKTS_1024_TO_1518_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_1519_TO_2047_OCTETS:
      return SAI_PORT_STAT_ETHER_OUT_PKTS_1519_TO_2047_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_2048_TO_4095_OCTETS:
      return SAI_PORT_STAT_ETHER_OUT_PKTS_2048_TO_4095_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_4096_TO_9216_OCTETS:
      return SAI_PORT_STAT_ETHER_OUT_PKTS_4096_TO_9216_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_ETHER_OUT_PKTS_9217_TO_16383_OCTETS:
      return SAI_PORT_STAT_ETHER_OUT_PKTS_9217_TO_16383_OCTETS;

    case lemming::dataplane::sai::PORT_STAT_IN_CURR_OCCUPANCY_BYTES:
      return SAI_PORT_STAT_IN_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::PORT_STAT_IN_WATERMARK_BYTES:
      return SAI_PORT_STAT_IN_WATERMARK_BYTES;

    case lemming::dataplane::sai::PORT_STAT_IN_SHARED_CURR_OCCUPANCY_BYTES:
      return SAI_PORT_STAT_IN_SHARED_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::PORT_STAT_IN_SHARED_WATERMARK_BYTES:
      return SAI_PORT_STAT_IN_SHARED_WATERMARK_BYTES;

    case lemming::dataplane::sai::PORT_STAT_OUT_CURR_OCCUPANCY_BYTES:
      return SAI_PORT_STAT_OUT_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::PORT_STAT_OUT_WATERMARK_BYTES:
      return SAI_PORT_STAT_OUT_WATERMARK_BYTES;

    case lemming::dataplane::sai::PORT_STAT_OUT_SHARED_CURR_OCCUPANCY_BYTES:
      return SAI_PORT_STAT_OUT_SHARED_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::PORT_STAT_OUT_SHARED_WATERMARK_BYTES:
      return SAI_PORT_STAT_OUT_SHARED_WATERMARK_BYTES;

    case lemming::dataplane::sai::PORT_STAT_IN_DROPPED_PKTS:
      return SAI_PORT_STAT_IN_DROPPED_PKTS;

    case lemming::dataplane::sai::PORT_STAT_OUT_DROPPED_PKTS:
      return SAI_PORT_STAT_OUT_DROPPED_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PAUSE_RX_PKTS:
      return SAI_PORT_STAT_PAUSE_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PAUSE_TX_PKTS:
      return SAI_PORT_STAT_PAUSE_TX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_0_RX_PKTS:
      return SAI_PORT_STAT_PFC_0_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_0_TX_PKTS:
      return SAI_PORT_STAT_PFC_0_TX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_1_RX_PKTS:
      return SAI_PORT_STAT_PFC_1_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_1_TX_PKTS:
      return SAI_PORT_STAT_PFC_1_TX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_2_RX_PKTS:
      return SAI_PORT_STAT_PFC_2_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_2_TX_PKTS:
      return SAI_PORT_STAT_PFC_2_TX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_3_RX_PKTS:
      return SAI_PORT_STAT_PFC_3_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_3_TX_PKTS:
      return SAI_PORT_STAT_PFC_3_TX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_4_RX_PKTS:
      return SAI_PORT_STAT_PFC_4_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_4_TX_PKTS:
      return SAI_PORT_STAT_PFC_4_TX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_5_RX_PKTS:
      return SAI_PORT_STAT_PFC_5_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_5_TX_PKTS:
      return SAI_PORT_STAT_PFC_5_TX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_6_RX_PKTS:
      return SAI_PORT_STAT_PFC_6_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_6_TX_PKTS:
      return SAI_PORT_STAT_PFC_6_TX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_7_RX_PKTS:
      return SAI_PORT_STAT_PFC_7_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_7_TX_PKTS:
      return SAI_PORT_STAT_PFC_7_TX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_0_RX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_0_RX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_0_TX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_0_TX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_1_RX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_1_RX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_1_TX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_1_TX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_2_RX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_2_RX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_2_TX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_2_TX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_3_RX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_3_RX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_3_TX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_3_TX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_4_RX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_4_RX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_4_TX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_4_TX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_5_RX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_5_RX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_5_TX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_5_TX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_6_RX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_6_RX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_6_TX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_6_TX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_7_RX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_7_RX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_7_TX_PAUSE_DURATION:
      return SAI_PORT_STAT_PFC_7_TX_PAUSE_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PFC_0_RX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_0_RX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_0_TX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_0_TX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_1_RX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_1_RX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_1_TX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_1_TX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_2_RX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_2_RX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_2_TX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_2_TX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_3_RX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_3_RX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_3_TX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_3_TX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_4_RX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_4_RX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_4_TX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_4_TX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_5_RX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_5_RX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_5_TX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_5_TX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_6_RX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_6_RX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_6_TX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_6_TX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_7_RX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_7_RX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_7_TX_PAUSE_DURATION_US:
      return SAI_PORT_STAT_PFC_7_TX_PAUSE_DURATION_US;

    case lemming::dataplane::sai::PORT_STAT_PFC_0_ON2OFF_RX_PKTS:
      return SAI_PORT_STAT_PFC_0_ON2OFF_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_1_ON2OFF_RX_PKTS:
      return SAI_PORT_STAT_PFC_1_ON2OFF_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_2_ON2OFF_RX_PKTS:
      return SAI_PORT_STAT_PFC_2_ON2OFF_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_3_ON2OFF_RX_PKTS:
      return SAI_PORT_STAT_PFC_3_ON2OFF_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_4_ON2OFF_RX_PKTS:
      return SAI_PORT_STAT_PFC_4_ON2OFF_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_5_ON2OFF_RX_PKTS:
      return SAI_PORT_STAT_PFC_5_ON2OFF_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_6_ON2OFF_RX_PKTS:
      return SAI_PORT_STAT_PFC_6_ON2OFF_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_PFC_7_ON2OFF_RX_PKTS:
      return SAI_PORT_STAT_PFC_7_ON2OFF_RX_PKTS;

    case lemming::dataplane::sai::PORT_STAT_DOT3_STATS_ALIGNMENT_ERRORS:
      return SAI_PORT_STAT_DOT3_STATS_ALIGNMENT_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_DOT3_STATS_FCS_ERRORS:
      return SAI_PORT_STAT_DOT3_STATS_FCS_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_DOT3_STATS_SINGLE_COLLISION_FRAMES:
      return SAI_PORT_STAT_DOT3_STATS_SINGLE_COLLISION_FRAMES;

    case lemming::dataplane::sai::
        PORT_STAT_DOT3_STATS_MULTIPLE_COLLISION_FRAMES:
      return SAI_PORT_STAT_DOT3_STATS_MULTIPLE_COLLISION_FRAMES;

    case lemming::dataplane::sai::PORT_STAT_DOT3_STATS_SQE_TEST_ERRORS:
      return SAI_PORT_STAT_DOT3_STATS_SQE_TEST_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_DOT3_STATS_DEFERRED_TRANSMISSIONS:
      return SAI_PORT_STAT_DOT3_STATS_DEFERRED_TRANSMISSIONS;

    case lemming::dataplane::sai::PORT_STAT_DOT3_STATS_LATE_COLLISIONS:
      return SAI_PORT_STAT_DOT3_STATS_LATE_COLLISIONS;

    case lemming::dataplane::sai::PORT_STAT_DOT3_STATS_EXCESSIVE_COLLISIONS:
      return SAI_PORT_STAT_DOT3_STATS_EXCESSIVE_COLLISIONS;

    case lemming::dataplane::sai::
        PORT_STAT_DOT3_STATS_INTERNAL_MAC_TRANSMIT_ERRORS:
      return SAI_PORT_STAT_DOT3_STATS_INTERNAL_MAC_TRANSMIT_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_DOT3_STATS_CARRIER_SENSE_ERRORS:
      return SAI_PORT_STAT_DOT3_STATS_CARRIER_SENSE_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_DOT3_STATS_FRAME_TOO_LONGS:
      return SAI_PORT_STAT_DOT3_STATS_FRAME_TOO_LONGS;

    case lemming::dataplane::sai::
        PORT_STAT_DOT3_STATS_INTERNAL_MAC_RECEIVE_ERRORS:
      return SAI_PORT_STAT_DOT3_STATS_INTERNAL_MAC_RECEIVE_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_DOT3_STATS_SYMBOL_ERRORS:
      return SAI_PORT_STAT_DOT3_STATS_SYMBOL_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_DOT3_CONTROL_IN_UNKNOWN_OPCODES:
      return SAI_PORT_STAT_DOT3_CONTROL_IN_UNKNOWN_OPCODES;

    case lemming::dataplane::sai::PORT_STAT_EEE_TX_EVENT_COUNT:
      return SAI_PORT_STAT_EEE_TX_EVENT_COUNT;

    case lemming::dataplane::sai::PORT_STAT_EEE_RX_EVENT_COUNT:
      return SAI_PORT_STAT_EEE_RX_EVENT_COUNT;

    case lemming::dataplane::sai::PORT_STAT_EEE_TX_DURATION:
      return SAI_PORT_STAT_EEE_TX_DURATION;

    case lemming::dataplane::sai::PORT_STAT_EEE_RX_DURATION:
      return SAI_PORT_STAT_EEE_RX_DURATION;

    case lemming::dataplane::sai::PORT_STAT_PRBS_ERROR_COUNT:
      return SAI_PORT_STAT_PRBS_ERROR_COUNT;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CORRECTABLE_FRAMES:
      return SAI_PORT_STAT_IF_IN_FEC_CORRECTABLE_FRAMES;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_NOT_CORRECTABLE_FRAMES:
      return SAI_PORT_STAT_IF_IN_FEC_NOT_CORRECTABLE_FRAMES;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_SYMBOL_ERRORS:
      return SAI_PORT_STAT_IF_IN_FEC_SYMBOL_ERRORS;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FABRIC_DATA_UNITS:
      return SAI_PORT_STAT_IF_IN_FABRIC_DATA_UNITS;

    case lemming::dataplane::sai::PORT_STAT_IF_OUT_FABRIC_DATA_UNITS:
      return SAI_PORT_STAT_IF_OUT_FABRIC_DATA_UNITS;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S0:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S0;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S1:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S1;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S2:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S2;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S3:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S3;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S4:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S4;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S5:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S5;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S6:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S6;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S7:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S7;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S8:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S8;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S9:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S9;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S10:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S10;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S11:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S11;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S12:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S12;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S13:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S13;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S14:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S14;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S15:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S15;

    case lemming::dataplane::sai::PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S16:
      return SAI_PORT_STAT_IF_IN_FEC_CODEWORD_ERRORS_S16;

    case lemming::dataplane::sai::PORT_STAT_IN_DROP_REASON_RANGE_BASE:
      return SAI_PORT_STAT_IN_DROP_REASON_RANGE_BASE;

    case lemming::dataplane::sai::
        PORT_STAT_IN_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS:
      return SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_IN_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS:
      return SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_IN_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS:
      return SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_IN_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS:
      return SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_IN_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS:
      return SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_IN_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS:
      return SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_IN_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS:
      return SAI_PORT_STAT_IN_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS;

    case lemming::dataplane::sai::PORT_STAT_IN_DROP_REASON_RANGE_END:
      return SAI_PORT_STAT_IN_DROP_REASON_RANGE_END;

    case lemming::dataplane::sai::PORT_STAT_OUT_DROP_REASON_RANGE_BASE:
      return SAI_PORT_STAT_OUT_DROP_REASON_RANGE_BASE;

    case lemming::dataplane::sai::
        PORT_STAT_OUT_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS:
      return SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_OUT_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS:
      return SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_OUT_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS:
      return SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_OUT_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS:
      return SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_OUT_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS:
      return SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_OUT_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS:
      return SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS;

    case lemming::dataplane::sai::
        PORT_STAT_OUT_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS:
      return SAI_PORT_STAT_OUT_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS;

    case lemming::dataplane::sai::PORT_STAT_OUT_DROP_REASON_RANGE_END:
      return SAI_PORT_STAT_OUT_DROP_REASON_RANGE_END;

    default:
      return SAI_PORT_STAT_IF_IN_OCTETS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_port_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::PortStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::PortType convert_sai_port_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_PORT_TYPE_LOGICAL:
      return lemming::dataplane::sai::PORT_TYPE_LOGICAL;

    case SAI_PORT_TYPE_CPU:
      return lemming::dataplane::sai::PORT_TYPE_CPU;

    case SAI_PORT_TYPE_FABRIC:
      return lemming::dataplane::sai::PORT_TYPE_FABRIC;

    case SAI_PORT_TYPE_RECYCLE:
      return lemming::dataplane::sai::PORT_TYPE_RECYCLE;

    default:
      return lemming::dataplane::sai::PORT_TYPE_UNSPECIFIED;
  }
}
sai_port_type_t convert_sai_port_type_t_to_sai(
    lemming::dataplane::sai::PortType val) {
  switch (val) {
    case lemming::dataplane::sai::PORT_TYPE_LOGICAL:
      return SAI_PORT_TYPE_LOGICAL;

    case lemming::dataplane::sai::PORT_TYPE_CPU:
      return SAI_PORT_TYPE_CPU;

    case lemming::dataplane::sai::PORT_TYPE_FABRIC:
      return SAI_PORT_TYPE_FABRIC;

    case lemming::dataplane::sai::PORT_TYPE_RECYCLE:
      return SAI_PORT_TYPE_RECYCLE;

    default:
      return SAI_PORT_TYPE_LOGICAL;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_port_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_port_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_port_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_port_type_t_to_sai(
        static_cast<lemming::dataplane::sai::PortType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::QosMapAttr convert_sai_qos_map_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_QOS_MAP_ATTR_TYPE:
      return lemming::dataplane::sai::QOS_MAP_ATTR_TYPE;

    case SAI_QOS_MAP_ATTR_MAP_TO_VALUE_LIST:
      return lemming::dataplane::sai::QOS_MAP_ATTR_MAP_TO_VALUE_LIST;

    default:
      return lemming::dataplane::sai::QOS_MAP_ATTR_UNSPECIFIED;
  }
}
sai_qos_map_attr_t convert_sai_qos_map_attr_t_to_sai(
    lemming::dataplane::sai::QosMapAttr val) {
  switch (val) {
    case lemming::dataplane::sai::QOS_MAP_ATTR_TYPE:
      return SAI_QOS_MAP_ATTR_TYPE;

    case lemming::dataplane::sai::QOS_MAP_ATTR_MAP_TO_VALUE_LIST:
      return SAI_QOS_MAP_ATTR_MAP_TO_VALUE_LIST;

    default:
      return SAI_QOS_MAP_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_qos_map_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_qos_map_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_qos_map_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_qos_map_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::QosMapAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::QosMapType convert_sai_qos_map_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_QOS_MAP_TYPE_DOT1P_TO_TC:
      return lemming::dataplane::sai::QOS_MAP_TYPE_DOT1P_TO_TC;

    case SAI_QOS_MAP_TYPE_DOT1P_TO_COLOR:
      return lemming::dataplane::sai::QOS_MAP_TYPE_DOT1P_TO_COLOR;

    case SAI_QOS_MAP_TYPE_DSCP_TO_TC:
      return lemming::dataplane::sai::QOS_MAP_TYPE_DSCP_TO_TC;

    case SAI_QOS_MAP_TYPE_DSCP_TO_COLOR:
      return lemming::dataplane::sai::QOS_MAP_TYPE_DSCP_TO_COLOR;

    case SAI_QOS_MAP_TYPE_TC_TO_QUEUE:
      return lemming::dataplane::sai::QOS_MAP_TYPE_TC_TO_QUEUE;

    case SAI_QOS_MAP_TYPE_TC_AND_COLOR_TO_DSCP:
      return lemming::dataplane::sai::QOS_MAP_TYPE_TC_AND_COLOR_TO_DSCP;

    case SAI_QOS_MAP_TYPE_TC_AND_COLOR_TO_DOT1P:
      return lemming::dataplane::sai::QOS_MAP_TYPE_TC_AND_COLOR_TO_DOT1P;

    case SAI_QOS_MAP_TYPE_TC_TO_PRIORITY_GROUP:
      return lemming::dataplane::sai::QOS_MAP_TYPE_TC_TO_PRIORITY_GROUP;

    case SAI_QOS_MAP_TYPE_PFC_PRIORITY_TO_PRIORITY_GROUP:
      return lemming::dataplane::sai::
          QOS_MAP_TYPE_PFC_PRIORITY_TO_PRIORITY_GROUP;

    case SAI_QOS_MAP_TYPE_PFC_PRIORITY_TO_QUEUE:
      return lemming::dataplane::sai::QOS_MAP_TYPE_PFC_PRIORITY_TO_QUEUE;

    case SAI_QOS_MAP_TYPE_MPLS_EXP_TO_TC:
      return lemming::dataplane::sai::QOS_MAP_TYPE_MPLS_EXP_TO_TC;

    case SAI_QOS_MAP_TYPE_MPLS_EXP_TO_COLOR:
      return lemming::dataplane::sai::QOS_MAP_TYPE_MPLS_EXP_TO_COLOR;

    case SAI_QOS_MAP_TYPE_TC_AND_COLOR_TO_MPLS_EXP:
      return lemming::dataplane::sai::QOS_MAP_TYPE_TC_AND_COLOR_TO_MPLS_EXP;

    case SAI_QOS_MAP_TYPE_DSCP_TO_FORWARDING_CLASS:
      return lemming::dataplane::sai::QOS_MAP_TYPE_DSCP_TO_FORWARDING_CLASS;

    case SAI_QOS_MAP_TYPE_MPLS_EXP_TO_FORWARDING_CLASS:
      return lemming::dataplane::sai::QOS_MAP_TYPE_MPLS_EXP_TO_FORWARDING_CLASS;

    case SAI_QOS_MAP_TYPE_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::QOS_MAP_TYPE_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::QOS_MAP_TYPE_UNSPECIFIED;
  }
}
sai_qos_map_type_t convert_sai_qos_map_type_t_to_sai(
    lemming::dataplane::sai::QosMapType val) {
  switch (val) {
    case lemming::dataplane::sai::QOS_MAP_TYPE_DOT1P_TO_TC:
      return SAI_QOS_MAP_TYPE_DOT1P_TO_TC;

    case lemming::dataplane::sai::QOS_MAP_TYPE_DOT1P_TO_COLOR:
      return SAI_QOS_MAP_TYPE_DOT1P_TO_COLOR;

    case lemming::dataplane::sai::QOS_MAP_TYPE_DSCP_TO_TC:
      return SAI_QOS_MAP_TYPE_DSCP_TO_TC;

    case lemming::dataplane::sai::QOS_MAP_TYPE_DSCP_TO_COLOR:
      return SAI_QOS_MAP_TYPE_DSCP_TO_COLOR;

    case lemming::dataplane::sai::QOS_MAP_TYPE_TC_TO_QUEUE:
      return SAI_QOS_MAP_TYPE_TC_TO_QUEUE;

    case lemming::dataplane::sai::QOS_MAP_TYPE_TC_AND_COLOR_TO_DSCP:
      return SAI_QOS_MAP_TYPE_TC_AND_COLOR_TO_DSCP;

    case lemming::dataplane::sai::QOS_MAP_TYPE_TC_AND_COLOR_TO_DOT1P:
      return SAI_QOS_MAP_TYPE_TC_AND_COLOR_TO_DOT1P;

    case lemming::dataplane::sai::QOS_MAP_TYPE_TC_TO_PRIORITY_GROUP:
      return SAI_QOS_MAP_TYPE_TC_TO_PRIORITY_GROUP;

    case lemming::dataplane::sai::QOS_MAP_TYPE_PFC_PRIORITY_TO_PRIORITY_GROUP:
      return SAI_QOS_MAP_TYPE_PFC_PRIORITY_TO_PRIORITY_GROUP;

    case lemming::dataplane::sai::QOS_MAP_TYPE_PFC_PRIORITY_TO_QUEUE:
      return SAI_QOS_MAP_TYPE_PFC_PRIORITY_TO_QUEUE;

    case lemming::dataplane::sai::QOS_MAP_TYPE_MPLS_EXP_TO_TC:
      return SAI_QOS_MAP_TYPE_MPLS_EXP_TO_TC;

    case lemming::dataplane::sai::QOS_MAP_TYPE_MPLS_EXP_TO_COLOR:
      return SAI_QOS_MAP_TYPE_MPLS_EXP_TO_COLOR;

    case lemming::dataplane::sai::QOS_MAP_TYPE_TC_AND_COLOR_TO_MPLS_EXP:
      return SAI_QOS_MAP_TYPE_TC_AND_COLOR_TO_MPLS_EXP;

    case lemming::dataplane::sai::QOS_MAP_TYPE_DSCP_TO_FORWARDING_CLASS:
      return SAI_QOS_MAP_TYPE_DSCP_TO_FORWARDING_CLASS;

    case lemming::dataplane::sai::QOS_MAP_TYPE_MPLS_EXP_TO_FORWARDING_CLASS:
      return SAI_QOS_MAP_TYPE_MPLS_EXP_TO_FORWARDING_CLASS;

    case lemming::dataplane::sai::QOS_MAP_TYPE_CUSTOM_RANGE_BASE:
      return SAI_QOS_MAP_TYPE_CUSTOM_RANGE_BASE;

    default:
      return SAI_QOS_MAP_TYPE_DOT1P_TO_TC;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_qos_map_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_qos_map_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_qos_map_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_qos_map_type_t_to_sai(
        static_cast<lemming::dataplane::sai::QosMapType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::QueueAttr convert_sai_queue_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_QUEUE_ATTR_TYPE:
      return lemming::dataplane::sai::QUEUE_ATTR_TYPE;

    case SAI_QUEUE_ATTR_PORT:
      return lemming::dataplane::sai::QUEUE_ATTR_PORT;

    case SAI_QUEUE_ATTR_INDEX:
      return lemming::dataplane::sai::QUEUE_ATTR_INDEX;

    case SAI_QUEUE_ATTR_PARENT_SCHEDULER_NODE:
      return lemming::dataplane::sai::QUEUE_ATTR_PARENT_SCHEDULER_NODE;

    case SAI_QUEUE_ATTR_WRED_PROFILE_ID:
      return lemming::dataplane::sai::QUEUE_ATTR_WRED_PROFILE_ID;

    case SAI_QUEUE_ATTR_BUFFER_PROFILE_ID:
      return lemming::dataplane::sai::QUEUE_ATTR_BUFFER_PROFILE_ID;

    case SAI_QUEUE_ATTR_SCHEDULER_PROFILE_ID:
      return lemming::dataplane::sai::QUEUE_ATTR_SCHEDULER_PROFILE_ID;

    case SAI_QUEUE_ATTR_PAUSE_STATUS:
      return lemming::dataplane::sai::QUEUE_ATTR_PAUSE_STATUS;

    case SAI_QUEUE_ATTR_ENABLE_PFC_DLDR:
      return lemming::dataplane::sai::QUEUE_ATTR_ENABLE_PFC_DLDR;

    case SAI_QUEUE_ATTR_PFC_DLR_INIT:
      return lemming::dataplane::sai::QUEUE_ATTR_PFC_DLR_INIT;

    case SAI_QUEUE_ATTR_TAM_OBJECT:
      return lemming::dataplane::sai::QUEUE_ATTR_TAM_OBJECT;

    case SAI_QUEUE_ATTR_PFC_DLR_PACKET_ACTION:
      return lemming::dataplane::sai::QUEUE_ATTR_PFC_DLR_PACKET_ACTION;

    case SAI_QUEUE_ATTR_PFC_CONTINUOUS_DEADLOCK_STATE:
      return lemming::dataplane::sai::QUEUE_ATTR_PFC_CONTINUOUS_DEADLOCK_STATE;

    default:
      return lemming::dataplane::sai::QUEUE_ATTR_UNSPECIFIED;
  }
}
sai_queue_attr_t convert_sai_queue_attr_t_to_sai(
    lemming::dataplane::sai::QueueAttr val) {
  switch (val) {
    case lemming::dataplane::sai::QUEUE_ATTR_TYPE:
      return SAI_QUEUE_ATTR_TYPE;

    case lemming::dataplane::sai::QUEUE_ATTR_PORT:
      return SAI_QUEUE_ATTR_PORT;

    case lemming::dataplane::sai::QUEUE_ATTR_INDEX:
      return SAI_QUEUE_ATTR_INDEX;

    case lemming::dataplane::sai::QUEUE_ATTR_PARENT_SCHEDULER_NODE:
      return SAI_QUEUE_ATTR_PARENT_SCHEDULER_NODE;

    case lemming::dataplane::sai::QUEUE_ATTR_WRED_PROFILE_ID:
      return SAI_QUEUE_ATTR_WRED_PROFILE_ID;

    case lemming::dataplane::sai::QUEUE_ATTR_BUFFER_PROFILE_ID:
      return SAI_QUEUE_ATTR_BUFFER_PROFILE_ID;

    case lemming::dataplane::sai::QUEUE_ATTR_SCHEDULER_PROFILE_ID:
      return SAI_QUEUE_ATTR_SCHEDULER_PROFILE_ID;

    case lemming::dataplane::sai::QUEUE_ATTR_PAUSE_STATUS:
      return SAI_QUEUE_ATTR_PAUSE_STATUS;

    case lemming::dataplane::sai::QUEUE_ATTR_ENABLE_PFC_DLDR:
      return SAI_QUEUE_ATTR_ENABLE_PFC_DLDR;

    case lemming::dataplane::sai::QUEUE_ATTR_PFC_DLR_INIT:
      return SAI_QUEUE_ATTR_PFC_DLR_INIT;

    case lemming::dataplane::sai::QUEUE_ATTR_TAM_OBJECT:
      return SAI_QUEUE_ATTR_TAM_OBJECT;

    case lemming::dataplane::sai::QUEUE_ATTR_PFC_DLR_PACKET_ACTION:
      return SAI_QUEUE_ATTR_PFC_DLR_PACKET_ACTION;

    case lemming::dataplane::sai::QUEUE_ATTR_PFC_CONTINUOUS_DEADLOCK_STATE:
      return SAI_QUEUE_ATTR_PFC_CONTINUOUS_DEADLOCK_STATE;

    default:
      return SAI_QUEUE_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_queue_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_queue_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_queue_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_queue_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::QueueAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::QueuePfcContinuousDeadlockState
convert_sai_queue_pfc_continuous_deadlock_state_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_NOT_PAUSED:
      return lemming::dataplane::sai::
          QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_NOT_PAUSED;

    case SAI_QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_PAUSED:
      return lemming::dataplane::sai::
          QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_PAUSED;

    case SAI_QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_PAUSED_NOT_CONTINUOUS:
      return lemming::dataplane::sai::
          QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_PAUSED_NOT_CONTINUOUS;

    default:
      return lemming::dataplane::sai::
          QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_UNSPECIFIED;
  }
}
sai_queue_pfc_continuous_deadlock_state_t
convert_sai_queue_pfc_continuous_deadlock_state_t_to_sai(
    lemming::dataplane::sai::QueuePfcContinuousDeadlockState val) {
  switch (val) {
    case lemming::dataplane::sai::
        QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_NOT_PAUSED:
      return SAI_QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_NOT_PAUSED;

    case lemming::dataplane::sai::QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_PAUSED:
      return SAI_QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_PAUSED;

    case lemming::dataplane::sai::
        QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_PAUSED_NOT_CONTINUOUS:
      return SAI_QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_PAUSED_NOT_CONTINUOUS;

    default:
      return SAI_QUEUE_PFC_CONTINUOUS_DEADLOCK_STATE_NOT_PAUSED;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_queue_pfc_continuous_deadlock_state_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_queue_pfc_continuous_deadlock_state_t_to_proto(
        list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_queue_pfc_continuous_deadlock_state_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_queue_pfc_continuous_deadlock_state_t_to_sai(
        static_cast<lemming::dataplane::sai::QueuePfcContinuousDeadlockState>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::QueuePfcDeadlockEventType
convert_sai_queue_pfc_deadlock_event_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_QUEUE_PFC_DEADLOCK_EVENT_TYPE_DETECTED:
      return lemming::dataplane::sai::QUEUE_PFC_DEADLOCK_EVENT_TYPE_DETECTED;

    case SAI_QUEUE_PFC_DEADLOCK_EVENT_TYPE_RECOVERED:
      return lemming::dataplane::sai::QUEUE_PFC_DEADLOCK_EVENT_TYPE_RECOVERED;

    default:
      return lemming::dataplane::sai::QUEUE_PFC_DEADLOCK_EVENT_TYPE_UNSPECIFIED;
  }
}
sai_queue_pfc_deadlock_event_type_t
convert_sai_queue_pfc_deadlock_event_type_t_to_sai(
    lemming::dataplane::sai::QueuePfcDeadlockEventType val) {
  switch (val) {
    case lemming::dataplane::sai::QUEUE_PFC_DEADLOCK_EVENT_TYPE_DETECTED:
      return SAI_QUEUE_PFC_DEADLOCK_EVENT_TYPE_DETECTED;

    case lemming::dataplane::sai::QUEUE_PFC_DEADLOCK_EVENT_TYPE_RECOVERED:
      return SAI_QUEUE_PFC_DEADLOCK_EVENT_TYPE_RECOVERED;

    default:
      return SAI_QUEUE_PFC_DEADLOCK_EVENT_TYPE_DETECTED;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_queue_pfc_deadlock_event_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_queue_pfc_deadlock_event_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_queue_pfc_deadlock_event_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_queue_pfc_deadlock_event_type_t_to_sai(
        static_cast<lemming::dataplane::sai::QueuePfcDeadlockEventType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::QueueStat convert_sai_queue_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_QUEUE_STAT_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_PACKETS;

    case SAI_QUEUE_STAT_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_BYTES;

    case SAI_QUEUE_STAT_DROPPED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_DROPPED_PACKETS;

    case SAI_QUEUE_STAT_DROPPED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_DROPPED_BYTES;

    case SAI_QUEUE_STAT_GREEN_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_GREEN_PACKETS;

    case SAI_QUEUE_STAT_GREEN_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_GREEN_BYTES;

    case SAI_QUEUE_STAT_GREEN_DROPPED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_GREEN_DROPPED_PACKETS;

    case SAI_QUEUE_STAT_GREEN_DROPPED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_GREEN_DROPPED_BYTES;

    case SAI_QUEUE_STAT_YELLOW_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_YELLOW_PACKETS;

    case SAI_QUEUE_STAT_YELLOW_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_YELLOW_BYTES;

    case SAI_QUEUE_STAT_YELLOW_DROPPED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_YELLOW_DROPPED_PACKETS;

    case SAI_QUEUE_STAT_YELLOW_DROPPED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_YELLOW_DROPPED_BYTES;

    case SAI_QUEUE_STAT_RED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_RED_PACKETS;

    case SAI_QUEUE_STAT_RED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_RED_BYTES;

    case SAI_QUEUE_STAT_RED_DROPPED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_RED_DROPPED_PACKETS;

    case SAI_QUEUE_STAT_RED_DROPPED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_RED_DROPPED_BYTES;

    case SAI_QUEUE_STAT_GREEN_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_GREEN_WRED_DROPPED_PACKETS;

    case SAI_QUEUE_STAT_GREEN_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_GREEN_WRED_DROPPED_BYTES;

    case SAI_QUEUE_STAT_YELLOW_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_YELLOW_WRED_DROPPED_PACKETS;

    case SAI_QUEUE_STAT_YELLOW_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_YELLOW_WRED_DROPPED_BYTES;

    case SAI_QUEUE_STAT_RED_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_RED_WRED_DROPPED_PACKETS;

    case SAI_QUEUE_STAT_RED_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_RED_WRED_DROPPED_BYTES;

    case SAI_QUEUE_STAT_WRED_DROPPED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_WRED_DROPPED_PACKETS;

    case SAI_QUEUE_STAT_WRED_DROPPED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_WRED_DROPPED_BYTES;

    case SAI_QUEUE_STAT_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_CURR_OCCUPANCY_BYTES;

    case SAI_QUEUE_STAT_WATERMARK_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_WATERMARK_BYTES;

    case SAI_QUEUE_STAT_SHARED_CURR_OCCUPANCY_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_SHARED_CURR_OCCUPANCY_BYTES;

    case SAI_QUEUE_STAT_SHARED_WATERMARK_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_SHARED_WATERMARK_BYTES;

    case SAI_QUEUE_STAT_GREEN_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_GREEN_WRED_ECN_MARKED_PACKETS;

    case SAI_QUEUE_STAT_GREEN_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_GREEN_WRED_ECN_MARKED_BYTES;

    case SAI_QUEUE_STAT_YELLOW_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_YELLOW_WRED_ECN_MARKED_PACKETS;

    case SAI_QUEUE_STAT_YELLOW_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_YELLOW_WRED_ECN_MARKED_BYTES;

    case SAI_QUEUE_STAT_RED_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_RED_WRED_ECN_MARKED_PACKETS;

    case SAI_QUEUE_STAT_RED_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_RED_WRED_ECN_MARKED_BYTES;

    case SAI_QUEUE_STAT_WRED_ECN_MARKED_PACKETS:
      return lemming::dataplane::sai::QUEUE_STAT_WRED_ECN_MARKED_PACKETS;

    case SAI_QUEUE_STAT_WRED_ECN_MARKED_BYTES:
      return lemming::dataplane::sai::QUEUE_STAT_WRED_ECN_MARKED_BYTES;

    case SAI_QUEUE_STAT_CURR_OCCUPANCY_LEVEL:
      return lemming::dataplane::sai::QUEUE_STAT_CURR_OCCUPANCY_LEVEL;

    case SAI_QUEUE_STAT_WATERMARK_LEVEL:
      return lemming::dataplane::sai::QUEUE_STAT_WATERMARK_LEVEL;

    case SAI_QUEUE_STAT_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::QUEUE_STAT_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::QUEUE_STAT_UNSPECIFIED;
  }
}
sai_queue_stat_t convert_sai_queue_stat_t_to_sai(
    lemming::dataplane::sai::QueueStat val) {
  switch (val) {
    case lemming::dataplane::sai::QUEUE_STAT_PACKETS:
      return SAI_QUEUE_STAT_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_BYTES:
      return SAI_QUEUE_STAT_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_DROPPED_PACKETS:
      return SAI_QUEUE_STAT_DROPPED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_DROPPED_BYTES:
      return SAI_QUEUE_STAT_DROPPED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_GREEN_PACKETS:
      return SAI_QUEUE_STAT_GREEN_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_GREEN_BYTES:
      return SAI_QUEUE_STAT_GREEN_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_GREEN_DROPPED_PACKETS:
      return SAI_QUEUE_STAT_GREEN_DROPPED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_GREEN_DROPPED_BYTES:
      return SAI_QUEUE_STAT_GREEN_DROPPED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_YELLOW_PACKETS:
      return SAI_QUEUE_STAT_YELLOW_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_YELLOW_BYTES:
      return SAI_QUEUE_STAT_YELLOW_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_YELLOW_DROPPED_PACKETS:
      return SAI_QUEUE_STAT_YELLOW_DROPPED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_YELLOW_DROPPED_BYTES:
      return SAI_QUEUE_STAT_YELLOW_DROPPED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_RED_PACKETS:
      return SAI_QUEUE_STAT_RED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_RED_BYTES:
      return SAI_QUEUE_STAT_RED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_RED_DROPPED_PACKETS:
      return SAI_QUEUE_STAT_RED_DROPPED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_RED_DROPPED_BYTES:
      return SAI_QUEUE_STAT_RED_DROPPED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_GREEN_WRED_DROPPED_PACKETS:
      return SAI_QUEUE_STAT_GREEN_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_GREEN_WRED_DROPPED_BYTES:
      return SAI_QUEUE_STAT_GREEN_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_YELLOW_WRED_DROPPED_PACKETS:
      return SAI_QUEUE_STAT_YELLOW_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_YELLOW_WRED_DROPPED_BYTES:
      return SAI_QUEUE_STAT_YELLOW_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_RED_WRED_DROPPED_PACKETS:
      return SAI_QUEUE_STAT_RED_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_RED_WRED_DROPPED_BYTES:
      return SAI_QUEUE_STAT_RED_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_WRED_DROPPED_PACKETS:
      return SAI_QUEUE_STAT_WRED_DROPPED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_WRED_DROPPED_BYTES:
      return SAI_QUEUE_STAT_WRED_DROPPED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_CURR_OCCUPANCY_BYTES:
      return SAI_QUEUE_STAT_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_WATERMARK_BYTES:
      return SAI_QUEUE_STAT_WATERMARK_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_SHARED_CURR_OCCUPANCY_BYTES:
      return SAI_QUEUE_STAT_SHARED_CURR_OCCUPANCY_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_SHARED_WATERMARK_BYTES:
      return SAI_QUEUE_STAT_SHARED_WATERMARK_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_GREEN_WRED_ECN_MARKED_PACKETS:
      return SAI_QUEUE_STAT_GREEN_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_GREEN_WRED_ECN_MARKED_BYTES:
      return SAI_QUEUE_STAT_GREEN_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_YELLOW_WRED_ECN_MARKED_PACKETS:
      return SAI_QUEUE_STAT_YELLOW_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_YELLOW_WRED_ECN_MARKED_BYTES:
      return SAI_QUEUE_STAT_YELLOW_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_RED_WRED_ECN_MARKED_PACKETS:
      return SAI_QUEUE_STAT_RED_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_RED_WRED_ECN_MARKED_BYTES:
      return SAI_QUEUE_STAT_RED_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_WRED_ECN_MARKED_PACKETS:
      return SAI_QUEUE_STAT_WRED_ECN_MARKED_PACKETS;

    case lemming::dataplane::sai::QUEUE_STAT_WRED_ECN_MARKED_BYTES:
      return SAI_QUEUE_STAT_WRED_ECN_MARKED_BYTES;

    case lemming::dataplane::sai::QUEUE_STAT_CURR_OCCUPANCY_LEVEL:
      return SAI_QUEUE_STAT_CURR_OCCUPANCY_LEVEL;

    case lemming::dataplane::sai::QUEUE_STAT_WATERMARK_LEVEL:
      return SAI_QUEUE_STAT_WATERMARK_LEVEL;

    case lemming::dataplane::sai::QUEUE_STAT_CUSTOM_RANGE_BASE:
      return SAI_QUEUE_STAT_CUSTOM_RANGE_BASE;

    default:
      return SAI_QUEUE_STAT_PACKETS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_queue_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_queue_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_queue_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_queue_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::QueueStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::QueueType convert_sai_queue_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_QUEUE_TYPE_ALL:
      return lemming::dataplane::sai::QUEUE_TYPE_ALL;

    case SAI_QUEUE_TYPE_UNICAST:
      return lemming::dataplane::sai::QUEUE_TYPE_UNICAST;

    case SAI_QUEUE_TYPE_MULTICAST:
      return lemming::dataplane::sai::QUEUE_TYPE_MULTICAST;

    case SAI_QUEUE_TYPE_UNICAST_VOQ:
      return lemming::dataplane::sai::QUEUE_TYPE_UNICAST_VOQ;

    case SAI_QUEUE_TYPE_MULTICAST_VOQ:
      return lemming::dataplane::sai::QUEUE_TYPE_MULTICAST_VOQ;

    case SAI_QUEUE_TYPE_FABRIC_TX:
      return lemming::dataplane::sai::QUEUE_TYPE_FABRIC_TX;

    case SAI_QUEUE_TYPE_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::QUEUE_TYPE_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::QUEUE_TYPE_UNSPECIFIED;
  }
}
sai_queue_type_t convert_sai_queue_type_t_to_sai(
    lemming::dataplane::sai::QueueType val) {
  switch (val) {
    case lemming::dataplane::sai::QUEUE_TYPE_ALL:
      return SAI_QUEUE_TYPE_ALL;

    case lemming::dataplane::sai::QUEUE_TYPE_UNICAST:
      return SAI_QUEUE_TYPE_UNICAST;

    case lemming::dataplane::sai::QUEUE_TYPE_MULTICAST:
      return SAI_QUEUE_TYPE_MULTICAST;

    case lemming::dataplane::sai::QUEUE_TYPE_UNICAST_VOQ:
      return SAI_QUEUE_TYPE_UNICAST_VOQ;

    case lemming::dataplane::sai::QUEUE_TYPE_MULTICAST_VOQ:
      return SAI_QUEUE_TYPE_MULTICAST_VOQ;

    case lemming::dataplane::sai::QUEUE_TYPE_FABRIC_TX:
      return SAI_QUEUE_TYPE_FABRIC_TX;

    case lemming::dataplane::sai::QUEUE_TYPE_CUSTOM_RANGE_BASE:
      return SAI_QUEUE_TYPE_CUSTOM_RANGE_BASE;

    default:
      return SAI_QUEUE_TYPE_ALL;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_queue_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_queue_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_queue_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_queue_type_t_to_sai(
        static_cast<lemming::dataplane::sai::QueueType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::RouteEntryAttr convert_sai_route_entry_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_ROUTE_ENTRY_ATTR_PACKET_ACTION:
      return lemming::dataplane::sai::ROUTE_ENTRY_ATTR_PACKET_ACTION;

    case SAI_ROUTE_ENTRY_ATTR_USER_TRAP_ID:
      return lemming::dataplane::sai::ROUTE_ENTRY_ATTR_USER_TRAP_ID;

    case SAI_ROUTE_ENTRY_ATTR_NEXT_HOP_ID:
      return lemming::dataplane::sai::ROUTE_ENTRY_ATTR_NEXT_HOP_ID;

    case SAI_ROUTE_ENTRY_ATTR_META_DATA:
      return lemming::dataplane::sai::ROUTE_ENTRY_ATTR_META_DATA;

    case SAI_ROUTE_ENTRY_ATTR_IP_ADDR_FAMILY:
      return lemming::dataplane::sai::ROUTE_ENTRY_ATTR_IP_ADDR_FAMILY;

    case SAI_ROUTE_ENTRY_ATTR_COUNTER_ID:
      return lemming::dataplane::sai::ROUTE_ENTRY_ATTR_COUNTER_ID;

    default:
      return lemming::dataplane::sai::ROUTE_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_route_entry_attr_t convert_sai_route_entry_attr_t_to_sai(
    lemming::dataplane::sai::RouteEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::ROUTE_ENTRY_ATTR_PACKET_ACTION:
      return SAI_ROUTE_ENTRY_ATTR_PACKET_ACTION;

    case lemming::dataplane::sai::ROUTE_ENTRY_ATTR_USER_TRAP_ID:
      return SAI_ROUTE_ENTRY_ATTR_USER_TRAP_ID;

    case lemming::dataplane::sai::ROUTE_ENTRY_ATTR_NEXT_HOP_ID:
      return SAI_ROUTE_ENTRY_ATTR_NEXT_HOP_ID;

    case lemming::dataplane::sai::ROUTE_ENTRY_ATTR_META_DATA:
      return SAI_ROUTE_ENTRY_ATTR_META_DATA;

    case lemming::dataplane::sai::ROUTE_ENTRY_ATTR_IP_ADDR_FAMILY:
      return SAI_ROUTE_ENTRY_ATTR_IP_ADDR_FAMILY;

    case lemming::dataplane::sai::ROUTE_ENTRY_ATTR_COUNTER_ID:
      return SAI_ROUTE_ENTRY_ATTR_COUNTER_ID;

    default:
      return SAI_ROUTE_ENTRY_ATTR_PACKET_ACTION;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_route_entry_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_route_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_route_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_route_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::RouteEntryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::RouterInterfaceAttr
convert_sai_router_interface_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ROUTER_INTERFACE_ATTR_VIRTUAL_ROUTER_ID:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_VIRTUAL_ROUTER_ID;

    case SAI_ROUTER_INTERFACE_ATTR_TYPE:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_TYPE;

    case SAI_ROUTER_INTERFACE_ATTR_PORT_ID:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_PORT_ID;

    case SAI_ROUTER_INTERFACE_ATTR_VLAN_ID:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_VLAN_ID;

    case SAI_ROUTER_INTERFACE_ATTR_OUTER_VLAN_ID:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_OUTER_VLAN_ID;

    case SAI_ROUTER_INTERFACE_ATTR_INNER_VLAN_ID:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_INNER_VLAN_ID;

    case SAI_ROUTER_INTERFACE_ATTR_BRIDGE_ID:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_BRIDGE_ID;

    case SAI_ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS;

    case SAI_ROUTER_INTERFACE_ATTR_ADMIN_V4_STATE:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_ADMIN_V4_STATE;

    case SAI_ROUTER_INTERFACE_ATTR_ADMIN_V6_STATE:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_ADMIN_V6_STATE;

    case SAI_ROUTER_INTERFACE_ATTR_MTU:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_MTU;

    case SAI_ROUTER_INTERFACE_ATTR_INGRESS_ACL:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_INGRESS_ACL;

    case SAI_ROUTER_INTERFACE_ATTR_EGRESS_ACL:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_EGRESS_ACL;

    case SAI_ROUTER_INTERFACE_ATTR_NEIGHBOR_MISS_PACKET_ACTION:
      return lemming::dataplane::sai::
          ROUTER_INTERFACE_ATTR_NEIGHBOR_MISS_PACKET_ACTION;

    case SAI_ROUTER_INTERFACE_ATTR_V4_MCAST_ENABLE:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_V4_MCAST_ENABLE;

    case SAI_ROUTER_INTERFACE_ATTR_V6_MCAST_ENABLE:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_V6_MCAST_ENABLE;

    case SAI_ROUTER_INTERFACE_ATTR_LOOPBACK_PACKET_ACTION:
      return lemming::dataplane::sai::
          ROUTER_INTERFACE_ATTR_LOOPBACK_PACKET_ACTION;

    case SAI_ROUTER_INTERFACE_ATTR_IS_VIRTUAL:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_IS_VIRTUAL;

    case SAI_ROUTER_INTERFACE_ATTR_NAT_ZONE_ID:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_NAT_ZONE_ID;

    case SAI_ROUTER_INTERFACE_ATTR_DISABLE_DECREMENT_TTL:
      return lemming::dataplane::sai::
          ROUTER_INTERFACE_ATTR_DISABLE_DECREMENT_TTL;

    case SAI_ROUTER_INTERFACE_ATTR_ADMIN_MPLS_STATE:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_ADMIN_MPLS_STATE;

    default:
      return lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_UNSPECIFIED;
  }
}
sai_router_interface_attr_t convert_sai_router_interface_attr_t_to_sai(
    lemming::dataplane::sai::RouterInterfaceAttr val) {
  switch (val) {
    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_VIRTUAL_ROUTER_ID:
      return SAI_ROUTER_INTERFACE_ATTR_VIRTUAL_ROUTER_ID;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_TYPE:
      return SAI_ROUTER_INTERFACE_ATTR_TYPE;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_PORT_ID:
      return SAI_ROUTER_INTERFACE_ATTR_PORT_ID;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_VLAN_ID:
      return SAI_ROUTER_INTERFACE_ATTR_VLAN_ID;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_OUTER_VLAN_ID:
      return SAI_ROUTER_INTERFACE_ATTR_OUTER_VLAN_ID;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_INNER_VLAN_ID:
      return SAI_ROUTER_INTERFACE_ATTR_INNER_VLAN_ID;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_BRIDGE_ID:
      return SAI_ROUTER_INTERFACE_ATTR_BRIDGE_ID;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS:
      return SAI_ROUTER_INTERFACE_ATTR_SRC_MAC_ADDRESS;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_ADMIN_V4_STATE:
      return SAI_ROUTER_INTERFACE_ATTR_ADMIN_V4_STATE;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_ADMIN_V6_STATE:
      return SAI_ROUTER_INTERFACE_ATTR_ADMIN_V6_STATE;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_MTU:
      return SAI_ROUTER_INTERFACE_ATTR_MTU;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_INGRESS_ACL:
      return SAI_ROUTER_INTERFACE_ATTR_INGRESS_ACL;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_EGRESS_ACL:
      return SAI_ROUTER_INTERFACE_ATTR_EGRESS_ACL;

    case lemming::dataplane::sai::
        ROUTER_INTERFACE_ATTR_NEIGHBOR_MISS_PACKET_ACTION:
      return SAI_ROUTER_INTERFACE_ATTR_NEIGHBOR_MISS_PACKET_ACTION;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_V4_MCAST_ENABLE:
      return SAI_ROUTER_INTERFACE_ATTR_V4_MCAST_ENABLE;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_V6_MCAST_ENABLE:
      return SAI_ROUTER_INTERFACE_ATTR_V6_MCAST_ENABLE;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_LOOPBACK_PACKET_ACTION:
      return SAI_ROUTER_INTERFACE_ATTR_LOOPBACK_PACKET_ACTION;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_IS_VIRTUAL:
      return SAI_ROUTER_INTERFACE_ATTR_IS_VIRTUAL;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_NAT_ZONE_ID:
      return SAI_ROUTER_INTERFACE_ATTR_NAT_ZONE_ID;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_DISABLE_DECREMENT_TTL:
      return SAI_ROUTER_INTERFACE_ATTR_DISABLE_DECREMENT_TTL;

    case lemming::dataplane::sai::ROUTER_INTERFACE_ATTR_ADMIN_MPLS_STATE:
      return SAI_ROUTER_INTERFACE_ATTR_ADMIN_MPLS_STATE;

    default:
      return SAI_ROUTER_INTERFACE_ATTR_VIRTUAL_ROUTER_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_router_interface_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_router_interface_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_router_interface_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_router_interface_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::RouterInterfaceAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::RouterInterfaceStat
convert_sai_router_interface_stat_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ROUTER_INTERFACE_STAT_IN_OCTETS:
      return lemming::dataplane::sai::ROUTER_INTERFACE_STAT_IN_OCTETS;

    case SAI_ROUTER_INTERFACE_STAT_IN_PACKETS:
      return lemming::dataplane::sai::ROUTER_INTERFACE_STAT_IN_PACKETS;

    case SAI_ROUTER_INTERFACE_STAT_OUT_OCTETS:
      return lemming::dataplane::sai::ROUTER_INTERFACE_STAT_OUT_OCTETS;

    case SAI_ROUTER_INTERFACE_STAT_OUT_PACKETS:
      return lemming::dataplane::sai::ROUTER_INTERFACE_STAT_OUT_PACKETS;

    case SAI_ROUTER_INTERFACE_STAT_IN_ERROR_OCTETS:
      return lemming::dataplane::sai::ROUTER_INTERFACE_STAT_IN_ERROR_OCTETS;

    case SAI_ROUTER_INTERFACE_STAT_IN_ERROR_PACKETS:
      return lemming::dataplane::sai::ROUTER_INTERFACE_STAT_IN_ERROR_PACKETS;

    case SAI_ROUTER_INTERFACE_STAT_OUT_ERROR_OCTETS:
      return lemming::dataplane::sai::ROUTER_INTERFACE_STAT_OUT_ERROR_OCTETS;

    case SAI_ROUTER_INTERFACE_STAT_OUT_ERROR_PACKETS:
      return lemming::dataplane::sai::ROUTER_INTERFACE_STAT_OUT_ERROR_PACKETS;

    default:
      return lemming::dataplane::sai::ROUTER_INTERFACE_STAT_UNSPECIFIED;
  }
}
sai_router_interface_stat_t convert_sai_router_interface_stat_t_to_sai(
    lemming::dataplane::sai::RouterInterfaceStat val) {
  switch (val) {
    case lemming::dataplane::sai::ROUTER_INTERFACE_STAT_IN_OCTETS:
      return SAI_ROUTER_INTERFACE_STAT_IN_OCTETS;

    case lemming::dataplane::sai::ROUTER_INTERFACE_STAT_IN_PACKETS:
      return SAI_ROUTER_INTERFACE_STAT_IN_PACKETS;

    case lemming::dataplane::sai::ROUTER_INTERFACE_STAT_OUT_OCTETS:
      return SAI_ROUTER_INTERFACE_STAT_OUT_OCTETS;

    case lemming::dataplane::sai::ROUTER_INTERFACE_STAT_OUT_PACKETS:
      return SAI_ROUTER_INTERFACE_STAT_OUT_PACKETS;

    case lemming::dataplane::sai::ROUTER_INTERFACE_STAT_IN_ERROR_OCTETS:
      return SAI_ROUTER_INTERFACE_STAT_IN_ERROR_OCTETS;

    case lemming::dataplane::sai::ROUTER_INTERFACE_STAT_IN_ERROR_PACKETS:
      return SAI_ROUTER_INTERFACE_STAT_IN_ERROR_PACKETS;

    case lemming::dataplane::sai::ROUTER_INTERFACE_STAT_OUT_ERROR_OCTETS:
      return SAI_ROUTER_INTERFACE_STAT_OUT_ERROR_OCTETS;

    case lemming::dataplane::sai::ROUTER_INTERFACE_STAT_OUT_ERROR_PACKETS:
      return SAI_ROUTER_INTERFACE_STAT_OUT_ERROR_PACKETS;

    default:
      return SAI_ROUTER_INTERFACE_STAT_IN_OCTETS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_router_interface_stat_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_router_interface_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_router_interface_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_router_interface_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::RouterInterfaceStat>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::RouterInterfaceType
convert_sai_router_interface_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_ROUTER_INTERFACE_TYPE_PORT:
      return lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_PORT;

    case SAI_ROUTER_INTERFACE_TYPE_VLAN:
      return lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_VLAN;

    case SAI_ROUTER_INTERFACE_TYPE_LOOPBACK:
      return lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_LOOPBACK;

    case SAI_ROUTER_INTERFACE_TYPE_MPLS_ROUTER:
      return lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_MPLS_ROUTER;

    case SAI_ROUTER_INTERFACE_TYPE_SUB_PORT:
      return lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_SUB_PORT;

    case SAI_ROUTER_INTERFACE_TYPE_BRIDGE:
      return lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_BRIDGE;

    case SAI_ROUTER_INTERFACE_TYPE_QINQ_PORT:
      return lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_QINQ_PORT;

    default:
      return lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_UNSPECIFIED;
  }
}
sai_router_interface_type_t convert_sai_router_interface_type_t_to_sai(
    lemming::dataplane::sai::RouterInterfaceType val) {
  switch (val) {
    case lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_PORT:
      return SAI_ROUTER_INTERFACE_TYPE_PORT;

    case lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_VLAN:
      return SAI_ROUTER_INTERFACE_TYPE_VLAN;

    case lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_LOOPBACK:
      return SAI_ROUTER_INTERFACE_TYPE_LOOPBACK;

    case lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_MPLS_ROUTER:
      return SAI_ROUTER_INTERFACE_TYPE_MPLS_ROUTER;

    case lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_SUB_PORT:
      return SAI_ROUTER_INTERFACE_TYPE_SUB_PORT;

    case lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_BRIDGE:
      return SAI_ROUTER_INTERFACE_TYPE_BRIDGE;

    case lemming::dataplane::sai::ROUTER_INTERFACE_TYPE_QINQ_PORT:
      return SAI_ROUTER_INTERFACE_TYPE_QINQ_PORT;

    default:
      return SAI_ROUTER_INTERFACE_TYPE_PORT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_router_interface_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_router_interface_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_router_interface_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_router_interface_type_t_to_sai(
        static_cast<lemming::dataplane::sai::RouterInterfaceType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::RpfGroupAttr convert_sai_rpf_group_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_RPF_GROUP_ATTR_RPF_INTERFACE_COUNT:
      return lemming::dataplane::sai::RPF_GROUP_ATTR_RPF_INTERFACE_COUNT;

    case SAI_RPF_GROUP_ATTR_RPF_MEMBER_LIST:
      return lemming::dataplane::sai::RPF_GROUP_ATTR_RPF_MEMBER_LIST;

    default:
      return lemming::dataplane::sai::RPF_GROUP_ATTR_UNSPECIFIED;
  }
}
sai_rpf_group_attr_t convert_sai_rpf_group_attr_t_to_sai(
    lemming::dataplane::sai::RpfGroupAttr val) {
  switch (val) {
    case lemming::dataplane::sai::RPF_GROUP_ATTR_RPF_INTERFACE_COUNT:
      return SAI_RPF_GROUP_ATTR_RPF_INTERFACE_COUNT;

    case lemming::dataplane::sai::RPF_GROUP_ATTR_RPF_MEMBER_LIST:
      return SAI_RPF_GROUP_ATTR_RPF_MEMBER_LIST;

    default:
      return SAI_RPF_GROUP_ATTR_RPF_INTERFACE_COUNT;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_rpf_group_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_rpf_group_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_rpf_group_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_rpf_group_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::RpfGroupAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::RpfGroupMemberAttr
convert_sai_rpf_group_member_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_RPF_GROUP_MEMBER_ATTR_RPF_GROUP_ID:
      return lemming::dataplane::sai::RPF_GROUP_MEMBER_ATTR_RPF_GROUP_ID;

    case SAI_RPF_GROUP_MEMBER_ATTR_RPF_INTERFACE_ID:
      return lemming::dataplane::sai::RPF_GROUP_MEMBER_ATTR_RPF_INTERFACE_ID;

    default:
      return lemming::dataplane::sai::RPF_GROUP_MEMBER_ATTR_UNSPECIFIED;
  }
}
sai_rpf_group_member_attr_t convert_sai_rpf_group_member_attr_t_to_sai(
    lemming::dataplane::sai::RpfGroupMemberAttr val) {
  switch (val) {
    case lemming::dataplane::sai::RPF_GROUP_MEMBER_ATTR_RPF_GROUP_ID:
      return SAI_RPF_GROUP_MEMBER_ATTR_RPF_GROUP_ID;

    case lemming::dataplane::sai::RPF_GROUP_MEMBER_ATTR_RPF_INTERFACE_ID:
      return SAI_RPF_GROUP_MEMBER_ATTR_RPF_INTERFACE_ID;

    default:
      return SAI_RPF_GROUP_MEMBER_ATTR_RPF_GROUP_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_rpf_group_member_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_rpf_group_member_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_rpf_group_member_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_rpf_group_member_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::RpfGroupMemberAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SamplepacketAttr
convert_sai_samplepacket_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SAMPLEPACKET_ATTR_SAMPLE_RATE:
      return lemming::dataplane::sai::SAMPLEPACKET_ATTR_SAMPLE_RATE;

    case SAI_SAMPLEPACKET_ATTR_TYPE:
      return lemming::dataplane::sai::SAMPLEPACKET_ATTR_TYPE;

    case SAI_SAMPLEPACKET_ATTR_MODE:
      return lemming::dataplane::sai::SAMPLEPACKET_ATTR_MODE;

    default:
      return lemming::dataplane::sai::SAMPLEPACKET_ATTR_UNSPECIFIED;
  }
}
sai_samplepacket_attr_t convert_sai_samplepacket_attr_t_to_sai(
    lemming::dataplane::sai::SamplepacketAttr val) {
  switch (val) {
    case lemming::dataplane::sai::SAMPLEPACKET_ATTR_SAMPLE_RATE:
      return SAI_SAMPLEPACKET_ATTR_SAMPLE_RATE;

    case lemming::dataplane::sai::SAMPLEPACKET_ATTR_TYPE:
      return SAI_SAMPLEPACKET_ATTR_TYPE;

    case lemming::dataplane::sai::SAMPLEPACKET_ATTR_MODE:
      return SAI_SAMPLEPACKET_ATTR_MODE;

    default:
      return SAI_SAMPLEPACKET_ATTR_SAMPLE_RATE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_samplepacket_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_samplepacket_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_samplepacket_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_samplepacket_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::SamplepacketAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SamplepacketMode
convert_sai_samplepacket_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SAMPLEPACKET_MODE_EXCLUSIVE:
      return lemming::dataplane::sai::SAMPLEPACKET_MODE_EXCLUSIVE;

    case SAI_SAMPLEPACKET_MODE_SHARED:
      return lemming::dataplane::sai::SAMPLEPACKET_MODE_SHARED;

    default:
      return lemming::dataplane::sai::SAMPLEPACKET_MODE_UNSPECIFIED;
  }
}
sai_samplepacket_mode_t convert_sai_samplepacket_mode_t_to_sai(
    lemming::dataplane::sai::SamplepacketMode val) {
  switch (val) {
    case lemming::dataplane::sai::SAMPLEPACKET_MODE_EXCLUSIVE:
      return SAI_SAMPLEPACKET_MODE_EXCLUSIVE;

    case lemming::dataplane::sai::SAMPLEPACKET_MODE_SHARED:
      return SAI_SAMPLEPACKET_MODE_SHARED;

    default:
      return SAI_SAMPLEPACKET_MODE_EXCLUSIVE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_samplepacket_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_samplepacket_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_samplepacket_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_samplepacket_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::SamplepacketMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SamplepacketType
convert_sai_samplepacket_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SAMPLEPACKET_TYPE_SLOW_PATH:
      return lemming::dataplane::sai::SAMPLEPACKET_TYPE_SLOW_PATH;

    case SAI_SAMPLEPACKET_TYPE_MIRROR_SESSION:
      return lemming::dataplane::sai::SAMPLEPACKET_TYPE_MIRROR_SESSION;

    default:
      return lemming::dataplane::sai::SAMPLEPACKET_TYPE_UNSPECIFIED;
  }
}
sai_samplepacket_type_t convert_sai_samplepacket_type_t_to_sai(
    lemming::dataplane::sai::SamplepacketType val) {
  switch (val) {
    case lemming::dataplane::sai::SAMPLEPACKET_TYPE_SLOW_PATH:
      return SAI_SAMPLEPACKET_TYPE_SLOW_PATH;

    case lemming::dataplane::sai::SAMPLEPACKET_TYPE_MIRROR_SESSION:
      return SAI_SAMPLEPACKET_TYPE_MIRROR_SESSION;

    default:
      return SAI_SAMPLEPACKET_TYPE_SLOW_PATH;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_samplepacket_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_samplepacket_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_samplepacket_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_samplepacket_type_t_to_sai(
        static_cast<lemming::dataplane::sai::SamplepacketType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SchedulerAttr convert_sai_scheduler_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_SCHEDULER_ATTR_SCHEDULING_TYPE:
      return lemming::dataplane::sai::SCHEDULER_ATTR_SCHEDULING_TYPE;

    case SAI_SCHEDULER_ATTR_SCHEDULING_WEIGHT:
      return lemming::dataplane::sai::SCHEDULER_ATTR_SCHEDULING_WEIGHT;

    case SAI_SCHEDULER_ATTR_METER_TYPE:
      return lemming::dataplane::sai::SCHEDULER_ATTR_METER_TYPE;

    case SAI_SCHEDULER_ATTR_MIN_BANDWIDTH_RATE:
      return lemming::dataplane::sai::SCHEDULER_ATTR_MIN_BANDWIDTH_RATE;

    case SAI_SCHEDULER_ATTR_MIN_BANDWIDTH_BURST_RATE:
      return lemming::dataplane::sai::SCHEDULER_ATTR_MIN_BANDWIDTH_BURST_RATE;

    case SAI_SCHEDULER_ATTR_MAX_BANDWIDTH_RATE:
      return lemming::dataplane::sai::SCHEDULER_ATTR_MAX_BANDWIDTH_RATE;

    case SAI_SCHEDULER_ATTR_MAX_BANDWIDTH_BURST_RATE:
      return lemming::dataplane::sai::SCHEDULER_ATTR_MAX_BANDWIDTH_BURST_RATE;

    default:
      return lemming::dataplane::sai::SCHEDULER_ATTR_UNSPECIFIED;
  }
}
sai_scheduler_attr_t convert_sai_scheduler_attr_t_to_sai(
    lemming::dataplane::sai::SchedulerAttr val) {
  switch (val) {
    case lemming::dataplane::sai::SCHEDULER_ATTR_SCHEDULING_TYPE:
      return SAI_SCHEDULER_ATTR_SCHEDULING_TYPE;

    case lemming::dataplane::sai::SCHEDULER_ATTR_SCHEDULING_WEIGHT:
      return SAI_SCHEDULER_ATTR_SCHEDULING_WEIGHT;

    case lemming::dataplane::sai::SCHEDULER_ATTR_METER_TYPE:
      return SAI_SCHEDULER_ATTR_METER_TYPE;

    case lemming::dataplane::sai::SCHEDULER_ATTR_MIN_BANDWIDTH_RATE:
      return SAI_SCHEDULER_ATTR_MIN_BANDWIDTH_RATE;

    case lemming::dataplane::sai::SCHEDULER_ATTR_MIN_BANDWIDTH_BURST_RATE:
      return SAI_SCHEDULER_ATTR_MIN_BANDWIDTH_BURST_RATE;

    case lemming::dataplane::sai::SCHEDULER_ATTR_MAX_BANDWIDTH_RATE:
      return SAI_SCHEDULER_ATTR_MAX_BANDWIDTH_RATE;

    case lemming::dataplane::sai::SCHEDULER_ATTR_MAX_BANDWIDTH_BURST_RATE:
      return SAI_SCHEDULER_ATTR_MAX_BANDWIDTH_BURST_RATE;

    default:
      return SAI_SCHEDULER_ATTR_SCHEDULING_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_scheduler_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_scheduler_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_scheduler_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_scheduler_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::SchedulerAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SchedulerGroupAttr
convert_sai_scheduler_group_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SCHEDULER_GROUP_ATTR_CHILD_COUNT:
      return lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_CHILD_COUNT;

    case SAI_SCHEDULER_GROUP_ATTR_CHILD_LIST:
      return lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_CHILD_LIST;

    case SAI_SCHEDULER_GROUP_ATTR_PORT_ID:
      return lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_PORT_ID;

    case SAI_SCHEDULER_GROUP_ATTR_LEVEL:
      return lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_LEVEL;

    case SAI_SCHEDULER_GROUP_ATTR_MAX_CHILDS:
      return lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_MAX_CHILDS;

    case SAI_SCHEDULER_GROUP_ATTR_SCHEDULER_PROFILE_ID:
      return lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_SCHEDULER_PROFILE_ID;

    case SAI_SCHEDULER_GROUP_ATTR_PARENT_NODE:
      return lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_PARENT_NODE;

    default:
      return lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_UNSPECIFIED;
  }
}
sai_scheduler_group_attr_t convert_sai_scheduler_group_attr_t_to_sai(
    lemming::dataplane::sai::SchedulerGroupAttr val) {
  switch (val) {
    case lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_CHILD_COUNT:
      return SAI_SCHEDULER_GROUP_ATTR_CHILD_COUNT;

    case lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_CHILD_LIST:
      return SAI_SCHEDULER_GROUP_ATTR_CHILD_LIST;

    case lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_PORT_ID:
      return SAI_SCHEDULER_GROUP_ATTR_PORT_ID;

    case lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_LEVEL:
      return SAI_SCHEDULER_GROUP_ATTR_LEVEL;

    case lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_MAX_CHILDS:
      return SAI_SCHEDULER_GROUP_ATTR_MAX_CHILDS;

    case lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_SCHEDULER_PROFILE_ID:
      return SAI_SCHEDULER_GROUP_ATTR_SCHEDULER_PROFILE_ID;

    case lemming::dataplane::sai::SCHEDULER_GROUP_ATTR_PARENT_NODE:
      return SAI_SCHEDULER_GROUP_ATTR_PARENT_NODE;

    default:
      return SAI_SCHEDULER_GROUP_ATTR_CHILD_COUNT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_scheduler_group_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_scheduler_group_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_scheduler_group_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_scheduler_group_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::SchedulerGroupAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SchedulingType convert_sai_scheduling_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_SCHEDULING_TYPE_STRICT:
      return lemming::dataplane::sai::SCHEDULING_TYPE_STRICT;

    case SAI_SCHEDULING_TYPE_WRR:
      return lemming::dataplane::sai::SCHEDULING_TYPE_WRR;

    case SAI_SCHEDULING_TYPE_DWRR:
      return lemming::dataplane::sai::SCHEDULING_TYPE_DWRR;

    default:
      return lemming::dataplane::sai::SCHEDULING_TYPE_UNSPECIFIED;
  }
}
sai_scheduling_type_t convert_sai_scheduling_type_t_to_sai(
    lemming::dataplane::sai::SchedulingType val) {
  switch (val) {
    case lemming::dataplane::sai::SCHEDULING_TYPE_STRICT:
      return SAI_SCHEDULING_TYPE_STRICT;

    case lemming::dataplane::sai::SCHEDULING_TYPE_WRR:
      return SAI_SCHEDULING_TYPE_WRR;

    case lemming::dataplane::sai::SCHEDULING_TYPE_DWRR:
      return SAI_SCHEDULING_TYPE_DWRR;

    default:
      return SAI_SCHEDULING_TYPE_STRICT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_scheduling_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_scheduling_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_scheduling_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_scheduling_type_t_to_sai(
        static_cast<lemming::dataplane::sai::SchedulingType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::Srv6SidlistAttr
convert_sai_srv6_sidlist_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SRV6_SIDLIST_ATTR_TYPE:
      return lemming::dataplane::sai::SRV6_SIDLIST_ATTR_TYPE;

    case SAI_SRV6_SIDLIST_ATTR_TLV_LIST:
      return lemming::dataplane::sai::SRV6_SIDLIST_ATTR_TLV_LIST;

    case SAI_SRV6_SIDLIST_ATTR_SEGMENT_LIST:
      return lemming::dataplane::sai::SRV6_SIDLIST_ATTR_SEGMENT_LIST;

    default:
      return lemming::dataplane::sai::SRV6_SIDLIST_ATTR_UNSPECIFIED;
  }
}
sai_srv6_sidlist_attr_t convert_sai_srv6_sidlist_attr_t_to_sai(
    lemming::dataplane::sai::Srv6SidlistAttr val) {
  switch (val) {
    case lemming::dataplane::sai::SRV6_SIDLIST_ATTR_TYPE:
      return SAI_SRV6_SIDLIST_ATTR_TYPE;

    case lemming::dataplane::sai::SRV6_SIDLIST_ATTR_TLV_LIST:
      return SAI_SRV6_SIDLIST_ATTR_TLV_LIST;

    case lemming::dataplane::sai::SRV6_SIDLIST_ATTR_SEGMENT_LIST:
      return SAI_SRV6_SIDLIST_ATTR_SEGMENT_LIST;

    default:
      return SAI_SRV6_SIDLIST_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_srv6_sidlist_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_srv6_sidlist_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_srv6_sidlist_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_srv6_sidlist_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::Srv6SidlistAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::Srv6SidlistType
convert_sai_srv6_sidlist_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SRV6_SIDLIST_TYPE_INSERT:
      return lemming::dataplane::sai::SRV6_SIDLIST_TYPE_INSERT;

    case SAI_SRV6_SIDLIST_TYPE_INSERT_RED:
      return lemming::dataplane::sai::SRV6_SIDLIST_TYPE_INSERT_RED;

    case SAI_SRV6_SIDLIST_TYPE_ENCAPS:
      return lemming::dataplane::sai::SRV6_SIDLIST_TYPE_ENCAPS;

    case SAI_SRV6_SIDLIST_TYPE_ENCAPS_RED:
      return lemming::dataplane::sai::SRV6_SIDLIST_TYPE_ENCAPS_RED;

    case SAI_SRV6_SIDLIST_TYPE_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::SRV6_SIDLIST_TYPE_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::SRV6_SIDLIST_TYPE_UNSPECIFIED;
  }
}
sai_srv6_sidlist_type_t convert_sai_srv6_sidlist_type_t_to_sai(
    lemming::dataplane::sai::Srv6SidlistType val) {
  switch (val) {
    case lemming::dataplane::sai::SRV6_SIDLIST_TYPE_INSERT:
      return SAI_SRV6_SIDLIST_TYPE_INSERT;

    case lemming::dataplane::sai::SRV6_SIDLIST_TYPE_INSERT_RED:
      return SAI_SRV6_SIDLIST_TYPE_INSERT_RED;

    case lemming::dataplane::sai::SRV6_SIDLIST_TYPE_ENCAPS:
      return SAI_SRV6_SIDLIST_TYPE_ENCAPS;

    case lemming::dataplane::sai::SRV6_SIDLIST_TYPE_ENCAPS_RED:
      return SAI_SRV6_SIDLIST_TYPE_ENCAPS_RED;

    case lemming::dataplane::sai::SRV6_SIDLIST_TYPE_CUSTOM_RANGE_BASE:
      return SAI_SRV6_SIDLIST_TYPE_CUSTOM_RANGE_BASE;

    default:
      return SAI_SRV6_SIDLIST_TYPE_INSERT;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_srv6_sidlist_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_srv6_sidlist_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_srv6_sidlist_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_srv6_sidlist_type_t_to_sai(
        static_cast<lemming::dataplane::sai::Srv6SidlistType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::StatsMode convert_sai_stats_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_STATS_MODE_READ:
      return lemming::dataplane::sai::STATS_MODE_READ;

    case SAI_STATS_MODE_READ_AND_CLEAR:
      return lemming::dataplane::sai::STATS_MODE_READ_AND_CLEAR;

    case SAI_STATS_MODE_BULK_READ:
      return lemming::dataplane::sai::STATS_MODE_BULK_READ;

    case SAI_STATS_MODE_BULK_CLEAR:
      return lemming::dataplane::sai::STATS_MODE_BULK_CLEAR;

    case SAI_STATS_MODE_BULK_READ_AND_CLEAR:
      return lemming::dataplane::sai::STATS_MODE_BULK_READ_AND_CLEAR;

    default:
      return lemming::dataplane::sai::STATS_MODE_UNSPECIFIED;
  }
}
sai_stats_mode_t convert_sai_stats_mode_t_to_sai(
    lemming::dataplane::sai::StatsMode val) {
  switch (val) {
    case lemming::dataplane::sai::STATS_MODE_READ:
      return SAI_STATS_MODE_READ;

    case lemming::dataplane::sai::STATS_MODE_READ_AND_CLEAR:
      return SAI_STATS_MODE_READ_AND_CLEAR;

    case lemming::dataplane::sai::STATS_MODE_BULK_READ:
      return SAI_STATS_MODE_BULK_READ;

    case lemming::dataplane::sai::STATS_MODE_BULK_CLEAR:
      return SAI_STATS_MODE_BULK_CLEAR;

    case lemming::dataplane::sai::STATS_MODE_BULK_READ_AND_CLEAR:
      return SAI_STATS_MODE_BULK_READ_AND_CLEAR;

    default:
      return SAI_STATS_MODE_READ;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_stats_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_stats_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_stats_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_stats_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::StatsMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::StpAttr convert_sai_stp_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_STP_ATTR_VLAN_LIST:
      return lemming::dataplane::sai::STP_ATTR_VLAN_LIST;

    case SAI_STP_ATTR_BRIDGE_ID:
      return lemming::dataplane::sai::STP_ATTR_BRIDGE_ID;

    case SAI_STP_ATTR_PORT_LIST:
      return lemming::dataplane::sai::STP_ATTR_PORT_LIST;

    default:
      return lemming::dataplane::sai::STP_ATTR_UNSPECIFIED;
  }
}
sai_stp_attr_t convert_sai_stp_attr_t_to_sai(
    lemming::dataplane::sai::StpAttr val) {
  switch (val) {
    case lemming::dataplane::sai::STP_ATTR_VLAN_LIST:
      return SAI_STP_ATTR_VLAN_LIST;

    case lemming::dataplane::sai::STP_ATTR_BRIDGE_ID:
      return SAI_STP_ATTR_BRIDGE_ID;

    case lemming::dataplane::sai::STP_ATTR_PORT_LIST:
      return SAI_STP_ATTR_PORT_LIST;

    default:
      return SAI_STP_ATTR_VLAN_LIST;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_stp_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_stp_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_stp_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_stp_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::StpAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::StpPortAttr convert_sai_stp_port_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_STP_PORT_ATTR_STP:
      return lemming::dataplane::sai::STP_PORT_ATTR_STP;

    case SAI_STP_PORT_ATTR_BRIDGE_PORT:
      return lemming::dataplane::sai::STP_PORT_ATTR_BRIDGE_PORT;

    case SAI_STP_PORT_ATTR_STATE:
      return lemming::dataplane::sai::STP_PORT_ATTR_STATE;

    default:
      return lemming::dataplane::sai::STP_PORT_ATTR_UNSPECIFIED;
  }
}
sai_stp_port_attr_t convert_sai_stp_port_attr_t_to_sai(
    lemming::dataplane::sai::StpPortAttr val) {
  switch (val) {
    case lemming::dataplane::sai::STP_PORT_ATTR_STP:
      return SAI_STP_PORT_ATTR_STP;

    case lemming::dataplane::sai::STP_PORT_ATTR_BRIDGE_PORT:
      return SAI_STP_PORT_ATTR_BRIDGE_PORT;

    case lemming::dataplane::sai::STP_PORT_ATTR_STATE:
      return SAI_STP_PORT_ATTR_STATE;

    default:
      return SAI_STP_PORT_ATTR_STP;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_stp_port_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_stp_port_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_stp_port_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_stp_port_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::StpPortAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::StpPortState convert_sai_stp_port_state_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_STP_PORT_STATE_LEARNING:
      return lemming::dataplane::sai::STP_PORT_STATE_LEARNING;

    case SAI_STP_PORT_STATE_FORWARDING:
      return lemming::dataplane::sai::STP_PORT_STATE_FORWARDING;

    case SAI_STP_PORT_STATE_BLOCKING:
      return lemming::dataplane::sai::STP_PORT_STATE_BLOCKING;

    default:
      return lemming::dataplane::sai::STP_PORT_STATE_UNSPECIFIED;
  }
}
sai_stp_port_state_t convert_sai_stp_port_state_t_to_sai(
    lemming::dataplane::sai::StpPortState val) {
  switch (val) {
    case lemming::dataplane::sai::STP_PORT_STATE_LEARNING:
      return SAI_STP_PORT_STATE_LEARNING;

    case lemming::dataplane::sai::STP_PORT_STATE_FORWARDING:
      return SAI_STP_PORT_STATE_FORWARDING;

    case lemming::dataplane::sai::STP_PORT_STATE_BLOCKING:
      return SAI_STP_PORT_STATE_BLOCKING;

    default:
      return SAI_STP_PORT_STATE_LEARNING;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_stp_port_state_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_stp_port_state_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_stp_port_state_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_stp_port_state_t_to_sai(
        static_cast<lemming::dataplane::sai::StpPortState>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchAttrExtensions
convert_sai_switch_attr_extensions_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_ATTR_EXTENSIONS_RANGE_START:
      return lemming::dataplane::sai::SWITCH_ATTR_EXTENSIONS_RANGE_START;

    case SAI_SWITCH_ATTR_EXTENSIONS_RANGE_END:
      return lemming::dataplane::sai::SWITCH_ATTR_EXTENSIONS_RANGE_END;

    default:
      return lemming::dataplane::sai::SWITCH_ATTR_EXTENSIONS_UNSPECIFIED;
  }
}
sai_switch_attr_extensions_t convert_sai_switch_attr_extensions_t_to_sai(
    lemming::dataplane::sai::SwitchAttrExtensions val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_ATTR_EXTENSIONS_RANGE_START:
      return SAI_SWITCH_ATTR_EXTENSIONS_RANGE_START;

    case lemming::dataplane::sai::SWITCH_ATTR_EXTENSIONS_RANGE_END:
      return SAI_SWITCH_ATTR_EXTENSIONS_RANGE_END;

    default:
      return SAI_SWITCH_ATTR_EXTENSIONS_RANGE_START;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_switch_attr_extensions_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_switch_attr_extensions_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_attr_extensions_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_attr_extensions_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchAttrExtensions>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchAttr convert_sai_switch_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_ATTR_NUMBER_OF_ACTIVE_PORTS:
      return lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_ACTIVE_PORTS;

    case SAI_SWITCH_ATTR_MAX_NUMBER_OF_SUPPORTED_PORTS:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_NUMBER_OF_SUPPORTED_PORTS;

    case SAI_SWITCH_ATTR_PORT_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_PORT_LIST;

    case SAI_SWITCH_ATTR_PORT_MAX_MTU:
      return lemming::dataplane::sai::SWITCH_ATTR_PORT_MAX_MTU;

    case SAI_SWITCH_ATTR_CPU_PORT:
      return lemming::dataplane::sai::SWITCH_ATTR_CPU_PORT;

    case SAI_SWITCH_ATTR_MAX_VIRTUAL_ROUTERS:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_VIRTUAL_ROUTERS;

    case SAI_SWITCH_ATTR_FDB_TABLE_SIZE:
      return lemming::dataplane::sai::SWITCH_ATTR_FDB_TABLE_SIZE;

    case SAI_SWITCH_ATTR_L3_NEIGHBOR_TABLE_SIZE:
      return lemming::dataplane::sai::SWITCH_ATTR_L3_NEIGHBOR_TABLE_SIZE;

    case SAI_SWITCH_ATTR_L3_ROUTE_TABLE_SIZE:
      return lemming::dataplane::sai::SWITCH_ATTR_L3_ROUTE_TABLE_SIZE;

    case SAI_SWITCH_ATTR_LAG_MEMBERS:
      return lemming::dataplane::sai::SWITCH_ATTR_LAG_MEMBERS;

    case SAI_SWITCH_ATTR_NUMBER_OF_LAGS:
      return lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_LAGS;

    case SAI_SWITCH_ATTR_ECMP_MEMBERS:
      return lemming::dataplane::sai::SWITCH_ATTR_ECMP_MEMBERS;

    case SAI_SWITCH_ATTR_NUMBER_OF_ECMP_GROUPS:
      return lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_ECMP_GROUPS;

    case SAI_SWITCH_ATTR_NUMBER_OF_UNICAST_QUEUES:
      return lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_UNICAST_QUEUES;

    case SAI_SWITCH_ATTR_NUMBER_OF_MULTICAST_QUEUES:
      return lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_MULTICAST_QUEUES;

    case SAI_SWITCH_ATTR_NUMBER_OF_QUEUES:
      return lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_QUEUES;

    case SAI_SWITCH_ATTR_NUMBER_OF_CPU_QUEUES:
      return lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_CPU_QUEUES;

    case SAI_SWITCH_ATTR_ON_LINK_ROUTE_SUPPORTED:
      return lemming::dataplane::sai::SWITCH_ATTR_ON_LINK_ROUTE_SUPPORTED;

    case SAI_SWITCH_ATTR_OPER_STATUS:
      return lemming::dataplane::sai::SWITCH_ATTR_OPER_STATUS;

    case SAI_SWITCH_ATTR_MAX_NUMBER_OF_TEMP_SENSORS:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_NUMBER_OF_TEMP_SENSORS;

    case SAI_SWITCH_ATTR_TEMP_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_TEMP_LIST;

    case SAI_SWITCH_ATTR_MAX_TEMP:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_TEMP;

    case SAI_SWITCH_ATTR_AVERAGE_TEMP:
      return lemming::dataplane::sai::SWITCH_ATTR_AVERAGE_TEMP;

    case SAI_SWITCH_ATTR_ACL_TABLE_MINIMUM_PRIORITY:
      return lemming::dataplane::sai::SWITCH_ATTR_ACL_TABLE_MINIMUM_PRIORITY;

    case SAI_SWITCH_ATTR_ACL_TABLE_MAXIMUM_PRIORITY:
      return lemming::dataplane::sai::SWITCH_ATTR_ACL_TABLE_MAXIMUM_PRIORITY;

    case SAI_SWITCH_ATTR_ACL_ENTRY_MINIMUM_PRIORITY:
      return lemming::dataplane::sai::SWITCH_ATTR_ACL_ENTRY_MINIMUM_PRIORITY;

    case SAI_SWITCH_ATTR_ACL_ENTRY_MAXIMUM_PRIORITY:
      return lemming::dataplane::sai::SWITCH_ATTR_ACL_ENTRY_MAXIMUM_PRIORITY;

    case SAI_SWITCH_ATTR_ACL_TABLE_GROUP_MINIMUM_PRIORITY:
      return lemming::dataplane::sai::
          SWITCH_ATTR_ACL_TABLE_GROUP_MINIMUM_PRIORITY;

    case SAI_SWITCH_ATTR_ACL_TABLE_GROUP_MAXIMUM_PRIORITY:
      return lemming::dataplane::sai::
          SWITCH_ATTR_ACL_TABLE_GROUP_MAXIMUM_PRIORITY;

    case SAI_SWITCH_ATTR_FDB_DST_USER_META_DATA_RANGE:
      return lemming::dataplane::sai::SWITCH_ATTR_FDB_DST_USER_META_DATA_RANGE;

    case SAI_SWITCH_ATTR_ROUTE_DST_USER_META_DATA_RANGE:
      return lemming::dataplane::sai::
          SWITCH_ATTR_ROUTE_DST_USER_META_DATA_RANGE;

    case SAI_SWITCH_ATTR_NEIGHBOR_DST_USER_META_DATA_RANGE:
      return lemming::dataplane::sai::
          SWITCH_ATTR_NEIGHBOR_DST_USER_META_DATA_RANGE;

    case SAI_SWITCH_ATTR_PORT_USER_META_DATA_RANGE:
      return lemming::dataplane::sai::SWITCH_ATTR_PORT_USER_META_DATA_RANGE;

    case SAI_SWITCH_ATTR_VLAN_USER_META_DATA_RANGE:
      return lemming::dataplane::sai::SWITCH_ATTR_VLAN_USER_META_DATA_RANGE;

    case SAI_SWITCH_ATTR_ACL_USER_META_DATA_RANGE:
      return lemming::dataplane::sai::SWITCH_ATTR_ACL_USER_META_DATA_RANGE;

    case SAI_SWITCH_ATTR_ACL_USER_TRAP_ID_RANGE:
      return lemming::dataplane::sai::SWITCH_ATTR_ACL_USER_TRAP_ID_RANGE;

    case SAI_SWITCH_ATTR_DEFAULT_VLAN_ID:
      return lemming::dataplane::sai::SWITCH_ATTR_DEFAULT_VLAN_ID;

    case SAI_SWITCH_ATTR_DEFAULT_STP_INST_ID:
      return lemming::dataplane::sai::SWITCH_ATTR_DEFAULT_STP_INST_ID;

    case SAI_SWITCH_ATTR_MAX_STP_INSTANCE:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_STP_INSTANCE;

    case SAI_SWITCH_ATTR_DEFAULT_VIRTUAL_ROUTER_ID:
      return lemming::dataplane::sai::SWITCH_ATTR_DEFAULT_VIRTUAL_ROUTER_ID;

    case SAI_SWITCH_ATTR_DEFAULT_OVERRIDE_VIRTUAL_ROUTER_ID:
      return lemming::dataplane::sai::
          SWITCH_ATTR_DEFAULT_OVERRIDE_VIRTUAL_ROUTER_ID;

    case SAI_SWITCH_ATTR_DEFAULT_1Q_BRIDGE_ID:
      return lemming::dataplane::sai::SWITCH_ATTR_DEFAULT_1Q_BRIDGE_ID;

    case SAI_SWITCH_ATTR_INGRESS_ACL:
      return lemming::dataplane::sai::SWITCH_ATTR_INGRESS_ACL;

    case SAI_SWITCH_ATTR_EGRESS_ACL:
      return lemming::dataplane::sai::SWITCH_ATTR_EGRESS_ACL;

    case SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_TRAFFIC_CLASSES:
      return lemming::dataplane::sai::
          SWITCH_ATTR_QOS_MAX_NUMBER_OF_TRAFFIC_CLASSES;

    case SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_SCHEDULER_GROUP_HIERARCHY_LEVELS:
      return lemming::dataplane::sai::
          SWITCH_ATTR_QOS_MAX_NUMBER_OF_SCHEDULER_GROUP_HIERARCHY_LEVELS;

    case SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_SCHEDULER_GROUPS_PER_HIERARCHY_LEVEL:
      return lemming::dataplane::sai::
          SWITCH_ATTR_QOS_MAX_NUMBER_OF_SCHEDULER_GROUPS_PER_HIERARCHY_LEVEL;

    case SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_CHILDS_PER_SCHEDULER_GROUP:
      return lemming::dataplane::sai::
          SWITCH_ATTR_QOS_MAX_NUMBER_OF_CHILDS_PER_SCHEDULER_GROUP;

    case SAI_SWITCH_ATTR_TOTAL_BUFFER_SIZE:
      return lemming::dataplane::sai::SWITCH_ATTR_TOTAL_BUFFER_SIZE;

    case SAI_SWITCH_ATTR_INGRESS_BUFFER_POOL_NUM:
      return lemming::dataplane::sai::SWITCH_ATTR_INGRESS_BUFFER_POOL_NUM;

    case SAI_SWITCH_ATTR_EGRESS_BUFFER_POOL_NUM:
      return lemming::dataplane::sai::SWITCH_ATTR_EGRESS_BUFFER_POOL_NUM;

    case SAI_SWITCH_ATTR_AVAILABLE_IPV4_ROUTE_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV4_ROUTE_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_IPV6_ROUTE_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV6_ROUTE_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_IPV4_NEXTHOP_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV4_NEXTHOP_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_IPV6_NEXTHOP_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV6_NEXTHOP_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_IPV4_NEIGHBOR_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV4_NEIGHBOR_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_IPV6_NEIGHBOR_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV6_NEIGHBOR_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_ENTRY:
      return lemming::dataplane::sai::
          SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_MEMBER_ENTRY:
      return lemming::dataplane::sai::
          SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_MEMBER_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_FDB_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_FDB_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_L2MC_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_L2MC_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_IPMC_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPMC_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_SNAT_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_SNAT_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_DNAT_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_DNAT_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_DOUBLE_NAT_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_DOUBLE_NAT_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_ACL_TABLE:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_ACL_TABLE;

    case SAI_SWITCH_ATTR_AVAILABLE_ACL_TABLE_GROUP:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_ACL_TABLE_GROUP;

    case SAI_SWITCH_ATTR_AVAILABLE_MY_SID_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_MY_SID_ENTRY;

    case SAI_SWITCH_ATTR_DEFAULT_TRAP_GROUP:
      return lemming::dataplane::sai::SWITCH_ATTR_DEFAULT_TRAP_GROUP;

    case SAI_SWITCH_ATTR_ECMP_HASH:
      return lemming::dataplane::sai::SWITCH_ATTR_ECMP_HASH;

    case SAI_SWITCH_ATTR_LAG_HASH:
      return lemming::dataplane::sai::SWITCH_ATTR_LAG_HASH;

    case SAI_SWITCH_ATTR_RESTART_WARM:
      return lemming::dataplane::sai::SWITCH_ATTR_RESTART_WARM;

    case SAI_SWITCH_ATTR_WARM_RECOVER:
      return lemming::dataplane::sai::SWITCH_ATTR_WARM_RECOVER;

    case SAI_SWITCH_ATTR_RESTART_TYPE:
      return lemming::dataplane::sai::SWITCH_ATTR_RESTART_TYPE;

    case SAI_SWITCH_ATTR_MIN_PLANNED_RESTART_INTERVAL:
      return lemming::dataplane::sai::SWITCH_ATTR_MIN_PLANNED_RESTART_INTERVAL;

    case SAI_SWITCH_ATTR_NV_STORAGE_SIZE:
      return lemming::dataplane::sai::SWITCH_ATTR_NV_STORAGE_SIZE;

    case SAI_SWITCH_ATTR_MAX_ACL_ACTION_COUNT:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_ACL_ACTION_COUNT;

    case SAI_SWITCH_ATTR_MAX_ACL_RANGE_COUNT:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_ACL_RANGE_COUNT;

    case SAI_SWITCH_ATTR_ACL_CAPABILITY:
      return lemming::dataplane::sai::SWITCH_ATTR_ACL_CAPABILITY;

    case SAI_SWITCH_ATTR_MCAST_SNOOPING_CAPABILITY:
      return lemming::dataplane::sai::SWITCH_ATTR_MCAST_SNOOPING_CAPABILITY;

    case SAI_SWITCH_ATTR_SWITCHING_MODE:
      return lemming::dataplane::sai::SWITCH_ATTR_SWITCHING_MODE;

    case SAI_SWITCH_ATTR_BCAST_CPU_FLOOD_ENABLE:
      return lemming::dataplane::sai::SWITCH_ATTR_BCAST_CPU_FLOOD_ENABLE;

    case SAI_SWITCH_ATTR_MCAST_CPU_FLOOD_ENABLE:
      return lemming::dataplane::sai::SWITCH_ATTR_MCAST_CPU_FLOOD_ENABLE;

    case SAI_SWITCH_ATTR_SRC_MAC_ADDRESS:
      return lemming::dataplane::sai::SWITCH_ATTR_SRC_MAC_ADDRESS;

    case SAI_SWITCH_ATTR_MAX_LEARNED_ADDRESSES:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_LEARNED_ADDRESSES;

    case SAI_SWITCH_ATTR_FDB_AGING_TIME:
      return lemming::dataplane::sai::SWITCH_ATTR_FDB_AGING_TIME;

    case SAI_SWITCH_ATTR_FDB_UNICAST_MISS_PACKET_ACTION:
      return lemming::dataplane::sai::
          SWITCH_ATTR_FDB_UNICAST_MISS_PACKET_ACTION;

    case SAI_SWITCH_ATTR_FDB_BROADCAST_MISS_PACKET_ACTION:
      return lemming::dataplane::sai::
          SWITCH_ATTR_FDB_BROADCAST_MISS_PACKET_ACTION;

    case SAI_SWITCH_ATTR_FDB_MULTICAST_MISS_PACKET_ACTION:
      return lemming::dataplane::sai::
          SWITCH_ATTR_FDB_MULTICAST_MISS_PACKET_ACTION;

    case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_ALGORITHM:
      return lemming::dataplane::sai::SWITCH_ATTR_ECMP_DEFAULT_HASH_ALGORITHM;

    case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_SEED:
      return lemming::dataplane::sai::SWITCH_ATTR_ECMP_DEFAULT_HASH_SEED;

    case SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_OFFSET:
      return lemming::dataplane::sai::SWITCH_ATTR_ECMP_DEFAULT_HASH_OFFSET;

    case SAI_SWITCH_ATTR_ECMP_DEFAULT_SYMMETRIC_HASH:
      return lemming::dataplane::sai::SWITCH_ATTR_ECMP_DEFAULT_SYMMETRIC_HASH;

    case SAI_SWITCH_ATTR_ECMP_HASH_IPV4:
      return lemming::dataplane::sai::SWITCH_ATTR_ECMP_HASH_IPV4;

    case SAI_SWITCH_ATTR_ECMP_HASH_IPV4_IN_IPV4:
      return lemming::dataplane::sai::SWITCH_ATTR_ECMP_HASH_IPV4_IN_IPV4;

    case SAI_SWITCH_ATTR_ECMP_HASH_IPV6:
      return lemming::dataplane::sai::SWITCH_ATTR_ECMP_HASH_IPV6;

    case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM:
      return lemming::dataplane::sai::SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM;

    case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_SEED:
      return lemming::dataplane::sai::SWITCH_ATTR_LAG_DEFAULT_HASH_SEED;

    case SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_OFFSET:
      return lemming::dataplane::sai::SWITCH_ATTR_LAG_DEFAULT_HASH_OFFSET;

    case SAI_SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH:
      return lemming::dataplane::sai::SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH;

    case SAI_SWITCH_ATTR_LAG_HASH_IPV4:
      return lemming::dataplane::sai::SWITCH_ATTR_LAG_HASH_IPV4;

    case SAI_SWITCH_ATTR_LAG_HASH_IPV4_IN_IPV4:
      return lemming::dataplane::sai::SWITCH_ATTR_LAG_HASH_IPV4_IN_IPV4;

    case SAI_SWITCH_ATTR_LAG_HASH_IPV6:
      return lemming::dataplane::sai::SWITCH_ATTR_LAG_HASH_IPV6;

    case SAI_SWITCH_ATTR_COUNTER_REFRESH_INTERVAL:
      return lemming::dataplane::sai::SWITCH_ATTR_COUNTER_REFRESH_INTERVAL;

    case SAI_SWITCH_ATTR_QOS_DEFAULT_TC:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_DEFAULT_TC;

    case SAI_SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP;

    case SAI_SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP;

    case SAI_SWITCH_ATTR_QOS_DSCP_TO_TC_MAP:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_DSCP_TO_TC_MAP;

    case SAI_SWITCH_ATTR_QOS_DSCP_TO_COLOR_MAP:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_DSCP_TO_COLOR_MAP;

    case SAI_SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP;

    case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP;

    case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP;

    case SAI_SWITCH_ATTR_SWITCH_SHELL_ENABLE:
      return lemming::dataplane::sai::SWITCH_ATTR_SWITCH_SHELL_ENABLE;

    case SAI_SWITCH_ATTR_SWITCH_PROFILE_ID:
      return lemming::dataplane::sai::SWITCH_ATTR_SWITCH_PROFILE_ID;

    case SAI_SWITCH_ATTR_SWITCH_HARDWARE_INFO:
      return lemming::dataplane::sai::SWITCH_ATTR_SWITCH_HARDWARE_INFO;

    case SAI_SWITCH_ATTR_FIRMWARE_PATH_NAME:
      return lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_PATH_NAME;

    case SAI_SWITCH_ATTR_INIT_SWITCH:
      return lemming::dataplane::sai::SWITCH_ATTR_INIT_SWITCH;

    case SAI_SWITCH_ATTR_SWITCH_STATE_CHANGE_NOTIFY:
      return lemming::dataplane::sai::SWITCH_ATTR_SWITCH_STATE_CHANGE_NOTIFY;

    case SAI_SWITCH_ATTR_SWITCH_SHUTDOWN_REQUEST_NOTIFY:
      return lemming::dataplane::sai::
          SWITCH_ATTR_SWITCH_SHUTDOWN_REQUEST_NOTIFY;

    case SAI_SWITCH_ATTR_FDB_EVENT_NOTIFY:
      return lemming::dataplane::sai::SWITCH_ATTR_FDB_EVENT_NOTIFY;

    case SAI_SWITCH_ATTR_PORT_STATE_CHANGE_NOTIFY:
      return lemming::dataplane::sai::SWITCH_ATTR_PORT_STATE_CHANGE_NOTIFY;

    case SAI_SWITCH_ATTR_PACKET_EVENT_NOTIFY:
      return lemming::dataplane::sai::SWITCH_ATTR_PACKET_EVENT_NOTIFY;

    case SAI_SWITCH_ATTR_FAST_API_ENABLE:
      return lemming::dataplane::sai::SWITCH_ATTR_FAST_API_ENABLE;

    case SAI_SWITCH_ATTR_MIRROR_TC:
      return lemming::dataplane::sai::SWITCH_ATTR_MIRROR_TC;

    case SAI_SWITCH_ATTR_ACL_STAGE_INGRESS:
      return lemming::dataplane::sai::SWITCH_ATTR_ACL_STAGE_INGRESS;

    case SAI_SWITCH_ATTR_ACL_STAGE_EGRESS:
      return lemming::dataplane::sai::SWITCH_ATTR_ACL_STAGE_EGRESS;

    case SAI_SWITCH_ATTR_SRV6_MAX_SID_DEPTH:
      return lemming::dataplane::sai::SWITCH_ATTR_SRV6_MAX_SID_DEPTH;

    case SAI_SWITCH_ATTR_SRV6_TLV_TYPE:
      return lemming::dataplane::sai::SWITCH_ATTR_SRV6_TLV_TYPE;

    case SAI_SWITCH_ATTR_QOS_NUM_LOSSLESS_QUEUES:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_NUM_LOSSLESS_QUEUES;

    case SAI_SWITCH_ATTR_QUEUE_PFC_DEADLOCK_NOTIFY:
      return lemming::dataplane::sai::SWITCH_ATTR_QUEUE_PFC_DEADLOCK_NOTIFY;

    case SAI_SWITCH_ATTR_PFC_DLR_PACKET_ACTION:
      return lemming::dataplane::sai::SWITCH_ATTR_PFC_DLR_PACKET_ACTION;

    case SAI_SWITCH_ATTR_PFC_TC_DLD_INTERVAL_RANGE:
      return lemming::dataplane::sai::SWITCH_ATTR_PFC_TC_DLD_INTERVAL_RANGE;

    case SAI_SWITCH_ATTR_PFC_TC_DLD_INTERVAL:
      return lemming::dataplane::sai::SWITCH_ATTR_PFC_TC_DLD_INTERVAL;

    case SAI_SWITCH_ATTR_PFC_TC_DLR_INTERVAL_RANGE:
      return lemming::dataplane::sai::SWITCH_ATTR_PFC_TC_DLR_INTERVAL_RANGE;

    case SAI_SWITCH_ATTR_PFC_TC_DLR_INTERVAL:
      return lemming::dataplane::sai::SWITCH_ATTR_PFC_TC_DLR_INTERVAL;

    case SAI_SWITCH_ATTR_SUPPORTED_PROTECTED_OBJECT_TYPE:
      return lemming::dataplane::sai::
          SWITCH_ATTR_SUPPORTED_PROTECTED_OBJECT_TYPE;

    case SAI_SWITCH_ATTR_TPID_OUTER_VLAN:
      return lemming::dataplane::sai::SWITCH_ATTR_TPID_OUTER_VLAN;

    case SAI_SWITCH_ATTR_TPID_INNER_VLAN:
      return lemming::dataplane::sai::SWITCH_ATTR_TPID_INNER_VLAN;

    case SAI_SWITCH_ATTR_CRC_CHECK_ENABLE:
      return lemming::dataplane::sai::SWITCH_ATTR_CRC_CHECK_ENABLE;

    case SAI_SWITCH_ATTR_CRC_RECALCULATION_ENABLE:
      return lemming::dataplane::sai::SWITCH_ATTR_CRC_RECALCULATION_ENABLE;

    case SAI_SWITCH_ATTR_BFD_SESSION_STATE_CHANGE_NOTIFY:
      return lemming::dataplane::sai::
          SWITCH_ATTR_BFD_SESSION_STATE_CHANGE_NOTIFY;

    case SAI_SWITCH_ATTR_NUMBER_OF_BFD_SESSION:
      return lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_BFD_SESSION;

    case SAI_SWITCH_ATTR_MAX_BFD_SESSION:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_BFD_SESSION;

    case SAI_SWITCH_ATTR_SUPPORTED_IPV4_BFD_SESSION_OFFLOAD_TYPE:
      return lemming::dataplane::sai::
          SWITCH_ATTR_SUPPORTED_IPV4_BFD_SESSION_OFFLOAD_TYPE;

    case SAI_SWITCH_ATTR_SUPPORTED_IPV6_BFD_SESSION_OFFLOAD_TYPE:
      return lemming::dataplane::sai::
          SWITCH_ATTR_SUPPORTED_IPV6_BFD_SESSION_OFFLOAD_TYPE;

    case SAI_SWITCH_ATTR_MIN_BFD_RX:
      return lemming::dataplane::sai::SWITCH_ATTR_MIN_BFD_RX;

    case SAI_SWITCH_ATTR_MIN_BFD_TX:
      return lemming::dataplane::sai::SWITCH_ATTR_MIN_BFD_TX;

    case SAI_SWITCH_ATTR_ECN_ECT_THRESHOLD_ENABLE:
      return lemming::dataplane::sai::SWITCH_ATTR_ECN_ECT_THRESHOLD_ENABLE;

    case SAI_SWITCH_ATTR_VXLAN_DEFAULT_ROUTER_MAC:
      return lemming::dataplane::sai::SWITCH_ATTR_VXLAN_DEFAULT_ROUTER_MAC;

    case SAI_SWITCH_ATTR_VXLAN_DEFAULT_PORT:
      return lemming::dataplane::sai::SWITCH_ATTR_VXLAN_DEFAULT_PORT;

    case SAI_SWITCH_ATTR_MAX_MIRROR_SESSION:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_MIRROR_SESSION;

    case SAI_SWITCH_ATTR_MAX_SAMPLED_MIRROR_SESSION:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_SAMPLED_MIRROR_SESSION;

    case SAI_SWITCH_ATTR_SUPPORTED_EXTENDED_STATS_MODE:
      return lemming::dataplane::sai::SWITCH_ATTR_SUPPORTED_EXTENDED_STATS_MODE;

    case SAI_SWITCH_ATTR_UNINIT_DATA_PLANE_ON_REMOVAL:
      return lemming::dataplane::sai::SWITCH_ATTR_UNINIT_DATA_PLANE_ON_REMOVAL;

    case SAI_SWITCH_ATTR_TAM_OBJECT_ID:
      return lemming::dataplane::sai::SWITCH_ATTR_TAM_OBJECT_ID;

    case SAI_SWITCH_ATTR_TAM_EVENT_NOTIFY:
      return lemming::dataplane::sai::SWITCH_ATTR_TAM_EVENT_NOTIFY;

    case SAI_SWITCH_ATTR_SUPPORTED_OBJECT_TYPE_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_SUPPORTED_OBJECT_TYPE_LIST;

    case SAI_SWITCH_ATTR_PRE_SHUTDOWN:
      return lemming::dataplane::sai::SWITCH_ATTR_PRE_SHUTDOWN;

    case SAI_SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID:
      return lemming::dataplane::sai::SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID;

    case SAI_SWITCH_ATTR_NAT_ENABLE:
      return lemming::dataplane::sai::SWITCH_ATTR_NAT_ENABLE;

    case SAI_SWITCH_ATTR_HARDWARE_ACCESS_BUS:
      return lemming::dataplane::sai::SWITCH_ATTR_HARDWARE_ACCESS_BUS;

    case SAI_SWITCH_ATTR_PLATFROM_CONTEXT:
      return lemming::dataplane::sai::SWITCH_ATTR_PLATFROM_CONTEXT;

    case SAI_SWITCH_ATTR_REGISTER_READ:
      return lemming::dataplane::sai::SWITCH_ATTR_REGISTER_READ;

    case SAI_SWITCH_ATTR_REGISTER_WRITE:
      return lemming::dataplane::sai::SWITCH_ATTR_REGISTER_WRITE;

    case SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_BROADCAST:
      return lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_DOWNLOAD_BROADCAST;

    case SAI_SWITCH_ATTR_FIRMWARE_LOAD_METHOD:
      return lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_LOAD_METHOD;

    case SAI_SWITCH_ATTR_FIRMWARE_LOAD_TYPE:
      return lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_LOAD_TYPE;

    case SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_EXECUTE:
      return lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_DOWNLOAD_EXECUTE;

    case SAI_SWITCH_ATTR_FIRMWARE_BROADCAST_STOP:
      return lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_BROADCAST_STOP;

    case SAI_SWITCH_ATTR_FIRMWARE_VERIFY_AND_INIT_SWITCH:
      return lemming::dataplane::sai::
          SWITCH_ATTR_FIRMWARE_VERIFY_AND_INIT_SWITCH;

    case SAI_SWITCH_ATTR_FIRMWARE_STATUS:
      return lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_STATUS;

    case SAI_SWITCH_ATTR_FIRMWARE_MAJOR_VERSION:
      return lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_MAJOR_VERSION;

    case SAI_SWITCH_ATTR_FIRMWARE_MINOR_VERSION:
      return lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_MINOR_VERSION;

    case SAI_SWITCH_ATTR_PORT_CONNECTOR_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_PORT_CONNECTOR_LIST;

    case SAI_SWITCH_ATTR_PROPOGATE_PORT_STATE_FROM_LINE_TO_SYSTEM_PORT_SUPPORT:
      return lemming::dataplane::sai::
          SWITCH_ATTR_PROPOGATE_PORT_STATE_FROM_LINE_TO_SYSTEM_PORT_SUPPORT;

    case SAI_SWITCH_ATTR_TYPE:
      return lemming::dataplane::sai::SWITCH_ATTR_TYPE;

    case SAI_SWITCH_ATTR_MACSEC_OBJECT_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_MACSEC_OBJECT_LIST;

    case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_MPLS_EXP_TO_TC_MAP;

    case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
      return lemming::dataplane::sai::SWITCH_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP;

    case SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
      return lemming::dataplane::sai::
          SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP;

    case SAI_SWITCH_ATTR_SWITCH_ID:
      return lemming::dataplane::sai::SWITCH_ATTR_SWITCH_ID;

    case SAI_SWITCH_ATTR_MAX_SYSTEM_CORES:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_SYSTEM_CORES;

    case SAI_SWITCH_ATTR_SYSTEM_PORT_CONFIG_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_SYSTEM_PORT_CONFIG_LIST;

    case SAI_SWITCH_ATTR_NUMBER_OF_SYSTEM_PORTS:
      return lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_SYSTEM_PORTS;

    case SAI_SWITCH_ATTR_SYSTEM_PORT_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_SYSTEM_PORT_LIST;

    case SAI_SWITCH_ATTR_NUMBER_OF_FABRIC_PORTS:
      return lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_FABRIC_PORTS;

    case SAI_SWITCH_ATTR_FABRIC_PORT_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_FABRIC_PORT_LIST;

    case SAI_SWITCH_ATTR_PACKET_DMA_MEMORY_POOL_SIZE:
      return lemming::dataplane::sai::SWITCH_ATTR_PACKET_DMA_MEMORY_POOL_SIZE;

    case SAI_SWITCH_ATTR_FAILOVER_CONFIG_MODE:
      return lemming::dataplane::sai::SWITCH_ATTR_FAILOVER_CONFIG_MODE;

    case SAI_SWITCH_ATTR_SUPPORTED_FAILOVER_MODE:
      return lemming::dataplane::sai::SWITCH_ATTR_SUPPORTED_FAILOVER_MODE;

    case SAI_SWITCH_ATTR_TUNNEL_OBJECTS_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_TUNNEL_OBJECTS_LIST;

    case SAI_SWITCH_ATTR_PACKET_AVAILABLE_DMA_MEMORY_POOL_SIZE:
      return lemming::dataplane::sai::
          SWITCH_ATTR_PACKET_AVAILABLE_DMA_MEMORY_POOL_SIZE;

    case SAI_SWITCH_ATTR_PRE_INGRESS_ACL:
      return lemming::dataplane::sai::SWITCH_ATTR_PRE_INGRESS_ACL;

    case SAI_SWITCH_ATTR_AVAILABLE_SNAPT_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_SNAPT_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_DNAPT_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_DNAPT_ENTRY;

    case SAI_SWITCH_ATTR_AVAILABLE_DOUBLE_NAPT_ENTRY:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_DOUBLE_NAPT_ENTRY;

    case SAI_SWITCH_ATTR_SLAVE_MDIO_ADDR_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_SLAVE_MDIO_ADDR_LIST;

    case SAI_SWITCH_ATTR_MY_MAC_TABLE_MINIMUM_PRIORITY:
      return lemming::dataplane::sai::SWITCH_ATTR_MY_MAC_TABLE_MINIMUM_PRIORITY;

    case SAI_SWITCH_ATTR_MY_MAC_TABLE_MAXIMUM_PRIORITY:
      return lemming::dataplane::sai::SWITCH_ATTR_MY_MAC_TABLE_MAXIMUM_PRIORITY;

    case SAI_SWITCH_ATTR_MY_MAC_LIST:
      return lemming::dataplane::sai::SWITCH_ATTR_MY_MAC_LIST;

    case SAI_SWITCH_ATTR_INSTALLED_MY_MAC_ENTRIES:
      return lemming::dataplane::sai::SWITCH_ATTR_INSTALLED_MY_MAC_ENTRIES;

    case SAI_SWITCH_ATTR_AVAILABLE_MY_MAC_ENTRIES:
      return lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_MY_MAC_ENTRIES;

    case SAI_SWITCH_ATTR_MAX_NUMBER_OF_FORWARDING_CLASSES:
      return lemming::dataplane::sai::
          SWITCH_ATTR_MAX_NUMBER_OF_FORWARDING_CLASSES;

    case SAI_SWITCH_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
      return lemming::dataplane::sai::
          SWITCH_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP;

    case SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
      return lemming::dataplane::sai::
          SWITCH_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP;

    case SAI_SWITCH_ATTR_IPSEC_OBJECT_ID:
      return lemming::dataplane::sai::SWITCH_ATTR_IPSEC_OBJECT_ID;

    case SAI_SWITCH_ATTR_IPSEC_SA_TAG_TPID:
      return lemming::dataplane::sai::SWITCH_ATTR_IPSEC_SA_TAG_TPID;

    case SAI_SWITCH_ATTR_IPSEC_SA_STATUS_CHANGE_NOTIFY:
      return lemming::dataplane::sai::SWITCH_ATTR_IPSEC_SA_STATUS_CHANGE_NOTIFY;

    case SAI_SWITCH_ATTR_NAT_EVENT_NOTIFY:
      return lemming::dataplane::sai::SWITCH_ATTR_NAT_EVENT_NOTIFY;

    case SAI_SWITCH_ATTR_MAX_ECMP_MEMBER_COUNT:
      return lemming::dataplane::sai::SWITCH_ATTR_MAX_ECMP_MEMBER_COUNT;

    case SAI_SWITCH_ATTR_ECMP_MEMBER_COUNT:
      return lemming::dataplane::sai::SWITCH_ATTR_ECMP_MEMBER_COUNT;

    default:
      return lemming::dataplane::sai::SWITCH_ATTR_UNSPECIFIED;
  }
}
sai_switch_attr_t convert_sai_switch_attr_t_to_sai(
    lemming::dataplane::sai::SwitchAttr val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_ACTIVE_PORTS:
      return SAI_SWITCH_ATTR_NUMBER_OF_ACTIVE_PORTS;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_NUMBER_OF_SUPPORTED_PORTS:
      return SAI_SWITCH_ATTR_MAX_NUMBER_OF_SUPPORTED_PORTS;

    case lemming::dataplane::sai::SWITCH_ATTR_PORT_LIST:
      return SAI_SWITCH_ATTR_PORT_LIST;

    case lemming::dataplane::sai::SWITCH_ATTR_PORT_MAX_MTU:
      return SAI_SWITCH_ATTR_PORT_MAX_MTU;

    case lemming::dataplane::sai::SWITCH_ATTR_CPU_PORT:
      return SAI_SWITCH_ATTR_CPU_PORT;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_VIRTUAL_ROUTERS:
      return SAI_SWITCH_ATTR_MAX_VIRTUAL_ROUTERS;

    case lemming::dataplane::sai::SWITCH_ATTR_FDB_TABLE_SIZE:
      return SAI_SWITCH_ATTR_FDB_TABLE_SIZE;

    case lemming::dataplane::sai::SWITCH_ATTR_L3_NEIGHBOR_TABLE_SIZE:
      return SAI_SWITCH_ATTR_L3_NEIGHBOR_TABLE_SIZE;

    case lemming::dataplane::sai::SWITCH_ATTR_L3_ROUTE_TABLE_SIZE:
      return SAI_SWITCH_ATTR_L3_ROUTE_TABLE_SIZE;

    case lemming::dataplane::sai::SWITCH_ATTR_LAG_MEMBERS:
      return SAI_SWITCH_ATTR_LAG_MEMBERS;

    case lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_LAGS:
      return SAI_SWITCH_ATTR_NUMBER_OF_LAGS;

    case lemming::dataplane::sai::SWITCH_ATTR_ECMP_MEMBERS:
      return SAI_SWITCH_ATTR_ECMP_MEMBERS;

    case lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_ECMP_GROUPS:
      return SAI_SWITCH_ATTR_NUMBER_OF_ECMP_GROUPS;

    case lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_UNICAST_QUEUES:
      return SAI_SWITCH_ATTR_NUMBER_OF_UNICAST_QUEUES;

    case lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_MULTICAST_QUEUES:
      return SAI_SWITCH_ATTR_NUMBER_OF_MULTICAST_QUEUES;

    case lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_QUEUES:
      return SAI_SWITCH_ATTR_NUMBER_OF_QUEUES;

    case lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_CPU_QUEUES:
      return SAI_SWITCH_ATTR_NUMBER_OF_CPU_QUEUES;

    case lemming::dataplane::sai::SWITCH_ATTR_ON_LINK_ROUTE_SUPPORTED:
      return SAI_SWITCH_ATTR_ON_LINK_ROUTE_SUPPORTED;

    case lemming::dataplane::sai::SWITCH_ATTR_OPER_STATUS:
      return SAI_SWITCH_ATTR_OPER_STATUS;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_NUMBER_OF_TEMP_SENSORS:
      return SAI_SWITCH_ATTR_MAX_NUMBER_OF_TEMP_SENSORS;

    case lemming::dataplane::sai::SWITCH_ATTR_TEMP_LIST:
      return SAI_SWITCH_ATTR_TEMP_LIST;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_TEMP:
      return SAI_SWITCH_ATTR_MAX_TEMP;

    case lemming::dataplane::sai::SWITCH_ATTR_AVERAGE_TEMP:
      return SAI_SWITCH_ATTR_AVERAGE_TEMP;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_TABLE_MINIMUM_PRIORITY:
      return SAI_SWITCH_ATTR_ACL_TABLE_MINIMUM_PRIORITY;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_TABLE_MAXIMUM_PRIORITY:
      return SAI_SWITCH_ATTR_ACL_TABLE_MAXIMUM_PRIORITY;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_ENTRY_MINIMUM_PRIORITY:
      return SAI_SWITCH_ATTR_ACL_ENTRY_MINIMUM_PRIORITY;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_ENTRY_MAXIMUM_PRIORITY:
      return SAI_SWITCH_ATTR_ACL_ENTRY_MAXIMUM_PRIORITY;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_TABLE_GROUP_MINIMUM_PRIORITY:
      return SAI_SWITCH_ATTR_ACL_TABLE_GROUP_MINIMUM_PRIORITY;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_TABLE_GROUP_MAXIMUM_PRIORITY:
      return SAI_SWITCH_ATTR_ACL_TABLE_GROUP_MAXIMUM_PRIORITY;

    case lemming::dataplane::sai::SWITCH_ATTR_FDB_DST_USER_META_DATA_RANGE:
      return SAI_SWITCH_ATTR_FDB_DST_USER_META_DATA_RANGE;

    case lemming::dataplane::sai::SWITCH_ATTR_ROUTE_DST_USER_META_DATA_RANGE:
      return SAI_SWITCH_ATTR_ROUTE_DST_USER_META_DATA_RANGE;

    case lemming::dataplane::sai::SWITCH_ATTR_NEIGHBOR_DST_USER_META_DATA_RANGE:
      return SAI_SWITCH_ATTR_NEIGHBOR_DST_USER_META_DATA_RANGE;

    case lemming::dataplane::sai::SWITCH_ATTR_PORT_USER_META_DATA_RANGE:
      return SAI_SWITCH_ATTR_PORT_USER_META_DATA_RANGE;

    case lemming::dataplane::sai::SWITCH_ATTR_VLAN_USER_META_DATA_RANGE:
      return SAI_SWITCH_ATTR_VLAN_USER_META_DATA_RANGE;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_USER_META_DATA_RANGE:
      return SAI_SWITCH_ATTR_ACL_USER_META_DATA_RANGE;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_USER_TRAP_ID_RANGE:
      return SAI_SWITCH_ATTR_ACL_USER_TRAP_ID_RANGE;

    case lemming::dataplane::sai::SWITCH_ATTR_DEFAULT_VLAN_ID:
      return SAI_SWITCH_ATTR_DEFAULT_VLAN_ID;

    case lemming::dataplane::sai::SWITCH_ATTR_DEFAULT_STP_INST_ID:
      return SAI_SWITCH_ATTR_DEFAULT_STP_INST_ID;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_STP_INSTANCE:
      return SAI_SWITCH_ATTR_MAX_STP_INSTANCE;

    case lemming::dataplane::sai::SWITCH_ATTR_DEFAULT_VIRTUAL_ROUTER_ID:
      return SAI_SWITCH_ATTR_DEFAULT_VIRTUAL_ROUTER_ID;

    case lemming::dataplane::sai::
        SWITCH_ATTR_DEFAULT_OVERRIDE_VIRTUAL_ROUTER_ID:
      return SAI_SWITCH_ATTR_DEFAULT_OVERRIDE_VIRTUAL_ROUTER_ID;

    case lemming::dataplane::sai::SWITCH_ATTR_DEFAULT_1Q_BRIDGE_ID:
      return SAI_SWITCH_ATTR_DEFAULT_1Q_BRIDGE_ID;

    case lemming::dataplane::sai::SWITCH_ATTR_INGRESS_ACL:
      return SAI_SWITCH_ATTR_INGRESS_ACL;

    case lemming::dataplane::sai::SWITCH_ATTR_EGRESS_ACL:
      return SAI_SWITCH_ATTR_EGRESS_ACL;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_MAX_NUMBER_OF_TRAFFIC_CLASSES:
      return SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_TRAFFIC_CLASSES;

    case lemming::dataplane::sai::
        SWITCH_ATTR_QOS_MAX_NUMBER_OF_SCHEDULER_GROUP_HIERARCHY_LEVELS:
      return SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_SCHEDULER_GROUP_HIERARCHY_LEVELS;

    case lemming::dataplane::sai::
        SWITCH_ATTR_QOS_MAX_NUMBER_OF_SCHEDULER_GROUPS_PER_HIERARCHY_LEVEL:
      return SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_SCHEDULER_GROUPS_PER_HIERARCHY_LEVEL;

    case lemming::dataplane::sai::
        SWITCH_ATTR_QOS_MAX_NUMBER_OF_CHILDS_PER_SCHEDULER_GROUP:
      return SAI_SWITCH_ATTR_QOS_MAX_NUMBER_OF_CHILDS_PER_SCHEDULER_GROUP;

    case lemming::dataplane::sai::SWITCH_ATTR_TOTAL_BUFFER_SIZE:
      return SAI_SWITCH_ATTR_TOTAL_BUFFER_SIZE;

    case lemming::dataplane::sai::SWITCH_ATTR_INGRESS_BUFFER_POOL_NUM:
      return SAI_SWITCH_ATTR_INGRESS_BUFFER_POOL_NUM;

    case lemming::dataplane::sai::SWITCH_ATTR_EGRESS_BUFFER_POOL_NUM:
      return SAI_SWITCH_ATTR_EGRESS_BUFFER_POOL_NUM;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV4_ROUTE_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_IPV4_ROUTE_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV6_ROUTE_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_IPV6_ROUTE_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV4_NEXTHOP_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_IPV4_NEXTHOP_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV6_NEXTHOP_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_IPV6_NEXTHOP_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV4_NEIGHBOR_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_IPV4_NEIGHBOR_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPV6_NEIGHBOR_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_IPV6_NEIGHBOR_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_ENTRY;

    case lemming::dataplane::sai::
        SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_MEMBER_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_NEXT_HOP_GROUP_MEMBER_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_FDB_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_FDB_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_L2MC_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_L2MC_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_IPMC_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_IPMC_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_SNAT_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_SNAT_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_DNAT_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_DNAT_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_DOUBLE_NAT_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_DOUBLE_NAT_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_ACL_TABLE:
      return SAI_SWITCH_ATTR_AVAILABLE_ACL_TABLE;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_ACL_TABLE_GROUP:
      return SAI_SWITCH_ATTR_AVAILABLE_ACL_TABLE_GROUP;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_MY_SID_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_MY_SID_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_DEFAULT_TRAP_GROUP:
      return SAI_SWITCH_ATTR_DEFAULT_TRAP_GROUP;

    case lemming::dataplane::sai::SWITCH_ATTR_ECMP_HASH:
      return SAI_SWITCH_ATTR_ECMP_HASH;

    case lemming::dataplane::sai::SWITCH_ATTR_LAG_HASH:
      return SAI_SWITCH_ATTR_LAG_HASH;

    case lemming::dataplane::sai::SWITCH_ATTR_RESTART_WARM:
      return SAI_SWITCH_ATTR_RESTART_WARM;

    case lemming::dataplane::sai::SWITCH_ATTR_WARM_RECOVER:
      return SAI_SWITCH_ATTR_WARM_RECOVER;

    case lemming::dataplane::sai::SWITCH_ATTR_RESTART_TYPE:
      return SAI_SWITCH_ATTR_RESTART_TYPE;

    case lemming::dataplane::sai::SWITCH_ATTR_MIN_PLANNED_RESTART_INTERVAL:
      return SAI_SWITCH_ATTR_MIN_PLANNED_RESTART_INTERVAL;

    case lemming::dataplane::sai::SWITCH_ATTR_NV_STORAGE_SIZE:
      return SAI_SWITCH_ATTR_NV_STORAGE_SIZE;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_ACL_ACTION_COUNT:
      return SAI_SWITCH_ATTR_MAX_ACL_ACTION_COUNT;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_ACL_RANGE_COUNT:
      return SAI_SWITCH_ATTR_MAX_ACL_RANGE_COUNT;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_CAPABILITY:
      return SAI_SWITCH_ATTR_ACL_CAPABILITY;

    case lemming::dataplane::sai::SWITCH_ATTR_MCAST_SNOOPING_CAPABILITY:
      return SAI_SWITCH_ATTR_MCAST_SNOOPING_CAPABILITY;

    case lemming::dataplane::sai::SWITCH_ATTR_SWITCHING_MODE:
      return SAI_SWITCH_ATTR_SWITCHING_MODE;

    case lemming::dataplane::sai::SWITCH_ATTR_BCAST_CPU_FLOOD_ENABLE:
      return SAI_SWITCH_ATTR_BCAST_CPU_FLOOD_ENABLE;

    case lemming::dataplane::sai::SWITCH_ATTR_MCAST_CPU_FLOOD_ENABLE:
      return SAI_SWITCH_ATTR_MCAST_CPU_FLOOD_ENABLE;

    case lemming::dataplane::sai::SWITCH_ATTR_SRC_MAC_ADDRESS:
      return SAI_SWITCH_ATTR_SRC_MAC_ADDRESS;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_LEARNED_ADDRESSES:
      return SAI_SWITCH_ATTR_MAX_LEARNED_ADDRESSES;

    case lemming::dataplane::sai::SWITCH_ATTR_FDB_AGING_TIME:
      return SAI_SWITCH_ATTR_FDB_AGING_TIME;

    case lemming::dataplane::sai::SWITCH_ATTR_FDB_UNICAST_MISS_PACKET_ACTION:
      return SAI_SWITCH_ATTR_FDB_UNICAST_MISS_PACKET_ACTION;

    case lemming::dataplane::sai::SWITCH_ATTR_FDB_BROADCAST_MISS_PACKET_ACTION:
      return SAI_SWITCH_ATTR_FDB_BROADCAST_MISS_PACKET_ACTION;

    case lemming::dataplane::sai::SWITCH_ATTR_FDB_MULTICAST_MISS_PACKET_ACTION:
      return SAI_SWITCH_ATTR_FDB_MULTICAST_MISS_PACKET_ACTION;

    case lemming::dataplane::sai::SWITCH_ATTR_ECMP_DEFAULT_HASH_ALGORITHM:
      return SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_ALGORITHM;

    case lemming::dataplane::sai::SWITCH_ATTR_ECMP_DEFAULT_HASH_SEED:
      return SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_SEED;

    case lemming::dataplane::sai::SWITCH_ATTR_ECMP_DEFAULT_HASH_OFFSET:
      return SAI_SWITCH_ATTR_ECMP_DEFAULT_HASH_OFFSET;

    case lemming::dataplane::sai::SWITCH_ATTR_ECMP_DEFAULT_SYMMETRIC_HASH:
      return SAI_SWITCH_ATTR_ECMP_DEFAULT_SYMMETRIC_HASH;

    case lemming::dataplane::sai::SWITCH_ATTR_ECMP_HASH_IPV4:
      return SAI_SWITCH_ATTR_ECMP_HASH_IPV4;

    case lemming::dataplane::sai::SWITCH_ATTR_ECMP_HASH_IPV4_IN_IPV4:
      return SAI_SWITCH_ATTR_ECMP_HASH_IPV4_IN_IPV4;

    case lemming::dataplane::sai::SWITCH_ATTR_ECMP_HASH_IPV6:
      return SAI_SWITCH_ATTR_ECMP_HASH_IPV6;

    case lemming::dataplane::sai::SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM:
      return SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_ALGORITHM;

    case lemming::dataplane::sai::SWITCH_ATTR_LAG_DEFAULT_HASH_SEED:
      return SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_SEED;

    case lemming::dataplane::sai::SWITCH_ATTR_LAG_DEFAULT_HASH_OFFSET:
      return SAI_SWITCH_ATTR_LAG_DEFAULT_HASH_OFFSET;

    case lemming::dataplane::sai::SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH:
      return SAI_SWITCH_ATTR_LAG_DEFAULT_SYMMETRIC_HASH;

    case lemming::dataplane::sai::SWITCH_ATTR_LAG_HASH_IPV4:
      return SAI_SWITCH_ATTR_LAG_HASH_IPV4;

    case lemming::dataplane::sai::SWITCH_ATTR_LAG_HASH_IPV4_IN_IPV4:
      return SAI_SWITCH_ATTR_LAG_HASH_IPV4_IN_IPV4;

    case lemming::dataplane::sai::SWITCH_ATTR_LAG_HASH_IPV6:
      return SAI_SWITCH_ATTR_LAG_HASH_IPV6;

    case lemming::dataplane::sai::SWITCH_ATTR_COUNTER_REFRESH_INTERVAL:
      return SAI_SWITCH_ATTR_COUNTER_REFRESH_INTERVAL;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_DEFAULT_TC:
      return SAI_SWITCH_ATTR_QOS_DEFAULT_TC;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP:
      return SAI_SWITCH_ATTR_QOS_DOT1P_TO_TC_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP:
      return SAI_SWITCH_ATTR_QOS_DOT1P_TO_COLOR_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_DSCP_TO_TC_MAP:
      return SAI_SWITCH_ATTR_QOS_DSCP_TO_TC_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_DSCP_TO_COLOR_MAP:
      return SAI_SWITCH_ATTR_QOS_DSCP_TO_COLOR_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP:
      return SAI_SWITCH_ATTR_QOS_TC_TO_QUEUE_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP:
      return SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DOT1P_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP:
      return SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_DSCP_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_SWITCH_SHELL_ENABLE:
      return SAI_SWITCH_ATTR_SWITCH_SHELL_ENABLE;

    case lemming::dataplane::sai::SWITCH_ATTR_SWITCH_PROFILE_ID:
      return SAI_SWITCH_ATTR_SWITCH_PROFILE_ID;

    case lemming::dataplane::sai::SWITCH_ATTR_SWITCH_HARDWARE_INFO:
      return SAI_SWITCH_ATTR_SWITCH_HARDWARE_INFO;

    case lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_PATH_NAME:
      return SAI_SWITCH_ATTR_FIRMWARE_PATH_NAME;

    case lemming::dataplane::sai::SWITCH_ATTR_INIT_SWITCH:
      return SAI_SWITCH_ATTR_INIT_SWITCH;

    case lemming::dataplane::sai::SWITCH_ATTR_SWITCH_STATE_CHANGE_NOTIFY:
      return SAI_SWITCH_ATTR_SWITCH_STATE_CHANGE_NOTIFY;

    case lemming::dataplane::sai::SWITCH_ATTR_SWITCH_SHUTDOWN_REQUEST_NOTIFY:
      return SAI_SWITCH_ATTR_SWITCH_SHUTDOWN_REQUEST_NOTIFY;

    case lemming::dataplane::sai::SWITCH_ATTR_FDB_EVENT_NOTIFY:
      return SAI_SWITCH_ATTR_FDB_EVENT_NOTIFY;

    case lemming::dataplane::sai::SWITCH_ATTR_PORT_STATE_CHANGE_NOTIFY:
      return SAI_SWITCH_ATTR_PORT_STATE_CHANGE_NOTIFY;

    case lemming::dataplane::sai::SWITCH_ATTR_PACKET_EVENT_NOTIFY:
      return SAI_SWITCH_ATTR_PACKET_EVENT_NOTIFY;

    case lemming::dataplane::sai::SWITCH_ATTR_FAST_API_ENABLE:
      return SAI_SWITCH_ATTR_FAST_API_ENABLE;

    case lemming::dataplane::sai::SWITCH_ATTR_MIRROR_TC:
      return SAI_SWITCH_ATTR_MIRROR_TC;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_STAGE_INGRESS:
      return SAI_SWITCH_ATTR_ACL_STAGE_INGRESS;

    case lemming::dataplane::sai::SWITCH_ATTR_ACL_STAGE_EGRESS:
      return SAI_SWITCH_ATTR_ACL_STAGE_EGRESS;

    case lemming::dataplane::sai::SWITCH_ATTR_SRV6_MAX_SID_DEPTH:
      return SAI_SWITCH_ATTR_SRV6_MAX_SID_DEPTH;

    case lemming::dataplane::sai::SWITCH_ATTR_SRV6_TLV_TYPE:
      return SAI_SWITCH_ATTR_SRV6_TLV_TYPE;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_NUM_LOSSLESS_QUEUES:
      return SAI_SWITCH_ATTR_QOS_NUM_LOSSLESS_QUEUES;

    case lemming::dataplane::sai::SWITCH_ATTR_QUEUE_PFC_DEADLOCK_NOTIFY:
      return SAI_SWITCH_ATTR_QUEUE_PFC_DEADLOCK_NOTIFY;

    case lemming::dataplane::sai::SWITCH_ATTR_PFC_DLR_PACKET_ACTION:
      return SAI_SWITCH_ATTR_PFC_DLR_PACKET_ACTION;

    case lemming::dataplane::sai::SWITCH_ATTR_PFC_TC_DLD_INTERVAL_RANGE:
      return SAI_SWITCH_ATTR_PFC_TC_DLD_INTERVAL_RANGE;

    case lemming::dataplane::sai::SWITCH_ATTR_PFC_TC_DLD_INTERVAL:
      return SAI_SWITCH_ATTR_PFC_TC_DLD_INTERVAL;

    case lemming::dataplane::sai::SWITCH_ATTR_PFC_TC_DLR_INTERVAL_RANGE:
      return SAI_SWITCH_ATTR_PFC_TC_DLR_INTERVAL_RANGE;

    case lemming::dataplane::sai::SWITCH_ATTR_PFC_TC_DLR_INTERVAL:
      return SAI_SWITCH_ATTR_PFC_TC_DLR_INTERVAL;

    case lemming::dataplane::sai::SWITCH_ATTR_SUPPORTED_PROTECTED_OBJECT_TYPE:
      return SAI_SWITCH_ATTR_SUPPORTED_PROTECTED_OBJECT_TYPE;

    case lemming::dataplane::sai::SWITCH_ATTR_TPID_OUTER_VLAN:
      return SAI_SWITCH_ATTR_TPID_OUTER_VLAN;

    case lemming::dataplane::sai::SWITCH_ATTR_TPID_INNER_VLAN:
      return SAI_SWITCH_ATTR_TPID_INNER_VLAN;

    case lemming::dataplane::sai::SWITCH_ATTR_CRC_CHECK_ENABLE:
      return SAI_SWITCH_ATTR_CRC_CHECK_ENABLE;

    case lemming::dataplane::sai::SWITCH_ATTR_CRC_RECALCULATION_ENABLE:
      return SAI_SWITCH_ATTR_CRC_RECALCULATION_ENABLE;

    case lemming::dataplane::sai::SWITCH_ATTR_BFD_SESSION_STATE_CHANGE_NOTIFY:
      return SAI_SWITCH_ATTR_BFD_SESSION_STATE_CHANGE_NOTIFY;

    case lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_BFD_SESSION:
      return SAI_SWITCH_ATTR_NUMBER_OF_BFD_SESSION;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_BFD_SESSION:
      return SAI_SWITCH_ATTR_MAX_BFD_SESSION;

    case lemming::dataplane::sai::
        SWITCH_ATTR_SUPPORTED_IPV4_BFD_SESSION_OFFLOAD_TYPE:
      return SAI_SWITCH_ATTR_SUPPORTED_IPV4_BFD_SESSION_OFFLOAD_TYPE;

    case lemming::dataplane::sai::
        SWITCH_ATTR_SUPPORTED_IPV6_BFD_SESSION_OFFLOAD_TYPE:
      return SAI_SWITCH_ATTR_SUPPORTED_IPV6_BFD_SESSION_OFFLOAD_TYPE;

    case lemming::dataplane::sai::SWITCH_ATTR_MIN_BFD_RX:
      return SAI_SWITCH_ATTR_MIN_BFD_RX;

    case lemming::dataplane::sai::SWITCH_ATTR_MIN_BFD_TX:
      return SAI_SWITCH_ATTR_MIN_BFD_TX;

    case lemming::dataplane::sai::SWITCH_ATTR_ECN_ECT_THRESHOLD_ENABLE:
      return SAI_SWITCH_ATTR_ECN_ECT_THRESHOLD_ENABLE;

    case lemming::dataplane::sai::SWITCH_ATTR_VXLAN_DEFAULT_ROUTER_MAC:
      return SAI_SWITCH_ATTR_VXLAN_DEFAULT_ROUTER_MAC;

    case lemming::dataplane::sai::SWITCH_ATTR_VXLAN_DEFAULT_PORT:
      return SAI_SWITCH_ATTR_VXLAN_DEFAULT_PORT;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_MIRROR_SESSION:
      return SAI_SWITCH_ATTR_MAX_MIRROR_SESSION;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_SAMPLED_MIRROR_SESSION:
      return SAI_SWITCH_ATTR_MAX_SAMPLED_MIRROR_SESSION;

    case lemming::dataplane::sai::SWITCH_ATTR_SUPPORTED_EXTENDED_STATS_MODE:
      return SAI_SWITCH_ATTR_SUPPORTED_EXTENDED_STATS_MODE;

    case lemming::dataplane::sai::SWITCH_ATTR_UNINIT_DATA_PLANE_ON_REMOVAL:
      return SAI_SWITCH_ATTR_UNINIT_DATA_PLANE_ON_REMOVAL;

    case lemming::dataplane::sai::SWITCH_ATTR_TAM_OBJECT_ID:
      return SAI_SWITCH_ATTR_TAM_OBJECT_ID;

    case lemming::dataplane::sai::SWITCH_ATTR_TAM_EVENT_NOTIFY:
      return SAI_SWITCH_ATTR_TAM_EVENT_NOTIFY;

    case lemming::dataplane::sai::SWITCH_ATTR_SUPPORTED_OBJECT_TYPE_LIST:
      return SAI_SWITCH_ATTR_SUPPORTED_OBJECT_TYPE_LIST;

    case lemming::dataplane::sai::SWITCH_ATTR_PRE_SHUTDOWN:
      return SAI_SWITCH_ATTR_PRE_SHUTDOWN;

    case lemming::dataplane::sai::SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID:
      return SAI_SWITCH_ATTR_NAT_ZONE_COUNTER_OBJECT_ID;

    case lemming::dataplane::sai::SWITCH_ATTR_NAT_ENABLE:
      return SAI_SWITCH_ATTR_NAT_ENABLE;

    case lemming::dataplane::sai::SWITCH_ATTR_HARDWARE_ACCESS_BUS:
      return SAI_SWITCH_ATTR_HARDWARE_ACCESS_BUS;

    case lemming::dataplane::sai::SWITCH_ATTR_PLATFROM_CONTEXT:
      return SAI_SWITCH_ATTR_PLATFROM_CONTEXT;

    case lemming::dataplane::sai::SWITCH_ATTR_REGISTER_READ:
      return SAI_SWITCH_ATTR_REGISTER_READ;

    case lemming::dataplane::sai::SWITCH_ATTR_REGISTER_WRITE:
      return SAI_SWITCH_ATTR_REGISTER_WRITE;

    case lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_DOWNLOAD_BROADCAST:
      return SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_BROADCAST;

    case lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_LOAD_METHOD:
      return SAI_SWITCH_ATTR_FIRMWARE_LOAD_METHOD;

    case lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_LOAD_TYPE:
      return SAI_SWITCH_ATTR_FIRMWARE_LOAD_TYPE;

    case lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_DOWNLOAD_EXECUTE:
      return SAI_SWITCH_ATTR_FIRMWARE_DOWNLOAD_EXECUTE;

    case lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_BROADCAST_STOP:
      return SAI_SWITCH_ATTR_FIRMWARE_BROADCAST_STOP;

    case lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_VERIFY_AND_INIT_SWITCH:
      return SAI_SWITCH_ATTR_FIRMWARE_VERIFY_AND_INIT_SWITCH;

    case lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_STATUS:
      return SAI_SWITCH_ATTR_FIRMWARE_STATUS;

    case lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_MAJOR_VERSION:
      return SAI_SWITCH_ATTR_FIRMWARE_MAJOR_VERSION;

    case lemming::dataplane::sai::SWITCH_ATTR_FIRMWARE_MINOR_VERSION:
      return SAI_SWITCH_ATTR_FIRMWARE_MINOR_VERSION;

    case lemming::dataplane::sai::SWITCH_ATTR_PORT_CONNECTOR_LIST:
      return SAI_SWITCH_ATTR_PORT_CONNECTOR_LIST;

    case lemming::dataplane::sai::
        SWITCH_ATTR_PROPOGATE_PORT_STATE_FROM_LINE_TO_SYSTEM_PORT_SUPPORT:
      return SAI_SWITCH_ATTR_PROPOGATE_PORT_STATE_FROM_LINE_TO_SYSTEM_PORT_SUPPORT;

    case lemming::dataplane::sai::SWITCH_ATTR_TYPE:
      return SAI_SWITCH_ATTR_TYPE;

    case lemming::dataplane::sai::SWITCH_ATTR_MACSEC_OBJECT_LIST:
      return SAI_SWITCH_ATTR_MACSEC_OBJECT_LIST;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_MPLS_EXP_TO_TC_MAP:
      return SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_TC_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP:
      return SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_COLOR_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP:
      return SAI_SWITCH_ATTR_QOS_TC_AND_COLOR_TO_MPLS_EXP_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_SWITCH_ID:
      return SAI_SWITCH_ATTR_SWITCH_ID;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_SYSTEM_CORES:
      return SAI_SWITCH_ATTR_MAX_SYSTEM_CORES;

    case lemming::dataplane::sai::SWITCH_ATTR_SYSTEM_PORT_CONFIG_LIST:
      return SAI_SWITCH_ATTR_SYSTEM_PORT_CONFIG_LIST;

    case lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_SYSTEM_PORTS:
      return SAI_SWITCH_ATTR_NUMBER_OF_SYSTEM_PORTS;

    case lemming::dataplane::sai::SWITCH_ATTR_SYSTEM_PORT_LIST:
      return SAI_SWITCH_ATTR_SYSTEM_PORT_LIST;

    case lemming::dataplane::sai::SWITCH_ATTR_NUMBER_OF_FABRIC_PORTS:
      return SAI_SWITCH_ATTR_NUMBER_OF_FABRIC_PORTS;

    case lemming::dataplane::sai::SWITCH_ATTR_FABRIC_PORT_LIST:
      return SAI_SWITCH_ATTR_FABRIC_PORT_LIST;

    case lemming::dataplane::sai::SWITCH_ATTR_PACKET_DMA_MEMORY_POOL_SIZE:
      return SAI_SWITCH_ATTR_PACKET_DMA_MEMORY_POOL_SIZE;

    case lemming::dataplane::sai::SWITCH_ATTR_FAILOVER_CONFIG_MODE:
      return SAI_SWITCH_ATTR_FAILOVER_CONFIG_MODE;

    case lemming::dataplane::sai::SWITCH_ATTR_SUPPORTED_FAILOVER_MODE:
      return SAI_SWITCH_ATTR_SUPPORTED_FAILOVER_MODE;

    case lemming::dataplane::sai::SWITCH_ATTR_TUNNEL_OBJECTS_LIST:
      return SAI_SWITCH_ATTR_TUNNEL_OBJECTS_LIST;

    case lemming::dataplane::sai::
        SWITCH_ATTR_PACKET_AVAILABLE_DMA_MEMORY_POOL_SIZE:
      return SAI_SWITCH_ATTR_PACKET_AVAILABLE_DMA_MEMORY_POOL_SIZE;

    case lemming::dataplane::sai::SWITCH_ATTR_PRE_INGRESS_ACL:
      return SAI_SWITCH_ATTR_PRE_INGRESS_ACL;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_SNAPT_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_SNAPT_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_DNAPT_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_DNAPT_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_DOUBLE_NAPT_ENTRY:
      return SAI_SWITCH_ATTR_AVAILABLE_DOUBLE_NAPT_ENTRY;

    case lemming::dataplane::sai::SWITCH_ATTR_SLAVE_MDIO_ADDR_LIST:
      return SAI_SWITCH_ATTR_SLAVE_MDIO_ADDR_LIST;

    case lemming::dataplane::sai::SWITCH_ATTR_MY_MAC_TABLE_MINIMUM_PRIORITY:
      return SAI_SWITCH_ATTR_MY_MAC_TABLE_MINIMUM_PRIORITY;

    case lemming::dataplane::sai::SWITCH_ATTR_MY_MAC_TABLE_MAXIMUM_PRIORITY:
      return SAI_SWITCH_ATTR_MY_MAC_TABLE_MAXIMUM_PRIORITY;

    case lemming::dataplane::sai::SWITCH_ATTR_MY_MAC_LIST:
      return SAI_SWITCH_ATTR_MY_MAC_LIST;

    case lemming::dataplane::sai::SWITCH_ATTR_INSTALLED_MY_MAC_ENTRIES:
      return SAI_SWITCH_ATTR_INSTALLED_MY_MAC_ENTRIES;

    case lemming::dataplane::sai::SWITCH_ATTR_AVAILABLE_MY_MAC_ENTRIES:
      return SAI_SWITCH_ATTR_AVAILABLE_MY_MAC_ENTRIES;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_NUMBER_OF_FORWARDING_CLASSES:
      return SAI_SWITCH_ATTR_MAX_NUMBER_OF_FORWARDING_CLASSES;

    case lemming::dataplane::sai::SWITCH_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP:
      return SAI_SWITCH_ATTR_QOS_DSCP_TO_FORWARDING_CLASS_MAP;

    case lemming::dataplane::sai::
        SWITCH_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP:
      return SAI_SWITCH_ATTR_QOS_MPLS_EXP_TO_FORWARDING_CLASS_MAP;

    case lemming::dataplane::sai::SWITCH_ATTR_IPSEC_OBJECT_ID:
      return SAI_SWITCH_ATTR_IPSEC_OBJECT_ID;

    case lemming::dataplane::sai::SWITCH_ATTR_IPSEC_SA_TAG_TPID:
      return SAI_SWITCH_ATTR_IPSEC_SA_TAG_TPID;

    case lemming::dataplane::sai::SWITCH_ATTR_IPSEC_SA_STATUS_CHANGE_NOTIFY:
      return SAI_SWITCH_ATTR_IPSEC_SA_STATUS_CHANGE_NOTIFY;

    case lemming::dataplane::sai::SWITCH_ATTR_NAT_EVENT_NOTIFY:
      return SAI_SWITCH_ATTR_NAT_EVENT_NOTIFY;

    case lemming::dataplane::sai::SWITCH_ATTR_MAX_ECMP_MEMBER_COUNT:
      return SAI_SWITCH_ATTR_MAX_ECMP_MEMBER_COUNT;

    case lemming::dataplane::sai::SWITCH_ATTR_ECMP_MEMBER_COUNT:
      return SAI_SWITCH_ATTR_ECMP_MEMBER_COUNT;

    default:
      return SAI_SWITCH_ATTR_NUMBER_OF_ACTIVE_PORTS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_switch_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_switch_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchFailoverConfigMode
convert_sai_switch_failover_config_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_FAILOVER_CONFIG_MODE_NO_HITLESS:
      return lemming::dataplane::sai::SWITCH_FAILOVER_CONFIG_MODE_NO_HITLESS;

    case SAI_SWITCH_FAILOVER_CONFIG_MODE_HITLESS:
      return lemming::dataplane::sai::SWITCH_FAILOVER_CONFIG_MODE_HITLESS;

    default:
      return lemming::dataplane::sai::SWITCH_FAILOVER_CONFIG_MODE_UNSPECIFIED;
  }
}
sai_switch_failover_config_mode_t
convert_sai_switch_failover_config_mode_t_to_sai(
    lemming::dataplane::sai::SwitchFailoverConfigMode val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_FAILOVER_CONFIG_MODE_NO_HITLESS:
      return SAI_SWITCH_FAILOVER_CONFIG_MODE_NO_HITLESS;

    case lemming::dataplane::sai::SWITCH_FAILOVER_CONFIG_MODE_HITLESS:
      return SAI_SWITCH_FAILOVER_CONFIG_MODE_HITLESS;

    default:
      return SAI_SWITCH_FAILOVER_CONFIG_MODE_NO_HITLESS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_switch_failover_config_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_switch_failover_config_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_failover_config_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_failover_config_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchFailoverConfigMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchFirmwareLoadMethod
convert_sai_switch_firmware_load_method_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_FIRMWARE_LOAD_METHOD_NONE:
      return lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_METHOD_NONE;

    case SAI_SWITCH_FIRMWARE_LOAD_METHOD_INTERNAL:
      return lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_METHOD_INTERNAL;

    case SAI_SWITCH_FIRMWARE_LOAD_METHOD_EEPROM:
      return lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_METHOD_EEPROM;

    default:
      return lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_METHOD_UNSPECIFIED;
  }
}
sai_switch_firmware_load_method_t
convert_sai_switch_firmware_load_method_t_to_sai(
    lemming::dataplane::sai::SwitchFirmwareLoadMethod val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_METHOD_NONE:
      return SAI_SWITCH_FIRMWARE_LOAD_METHOD_NONE;

    case lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_METHOD_INTERNAL:
      return SAI_SWITCH_FIRMWARE_LOAD_METHOD_INTERNAL;

    case lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_METHOD_EEPROM:
      return SAI_SWITCH_FIRMWARE_LOAD_METHOD_EEPROM;

    default:
      return SAI_SWITCH_FIRMWARE_LOAD_METHOD_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_switch_firmware_load_method_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_switch_firmware_load_method_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_firmware_load_method_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_firmware_load_method_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchFirmwareLoadMethod>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchFirmwareLoadType
convert_sai_switch_firmware_load_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_FIRMWARE_LOAD_TYPE_SKIP:
      return lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_TYPE_SKIP;

    case SAI_SWITCH_FIRMWARE_LOAD_TYPE_FORCE:
      return lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_TYPE_FORCE;

    case SAI_SWITCH_FIRMWARE_LOAD_TYPE_AUTO:
      return lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_TYPE_AUTO;

    default:
      return lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_TYPE_UNSPECIFIED;
  }
}
sai_switch_firmware_load_type_t convert_sai_switch_firmware_load_type_t_to_sai(
    lemming::dataplane::sai::SwitchFirmwareLoadType val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_TYPE_SKIP:
      return SAI_SWITCH_FIRMWARE_LOAD_TYPE_SKIP;

    case lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_TYPE_FORCE:
      return SAI_SWITCH_FIRMWARE_LOAD_TYPE_FORCE;

    case lemming::dataplane::sai::SWITCH_FIRMWARE_LOAD_TYPE_AUTO:
      return SAI_SWITCH_FIRMWARE_LOAD_TYPE_AUTO;

    default:
      return SAI_SWITCH_FIRMWARE_LOAD_TYPE_SKIP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_switch_firmware_load_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_switch_firmware_load_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_firmware_load_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_firmware_load_type_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchFirmwareLoadType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchHardwareAccessBus
convert_sai_switch_hardware_access_bus_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_HARDWARE_ACCESS_BUS_MDIO:
      return lemming::dataplane::sai::SWITCH_HARDWARE_ACCESS_BUS_MDIO;

    case SAI_SWITCH_HARDWARE_ACCESS_BUS_I2C:
      return lemming::dataplane::sai::SWITCH_HARDWARE_ACCESS_BUS_I2C;

    case SAI_SWITCH_HARDWARE_ACCESS_BUS_CPLD:
      return lemming::dataplane::sai::SWITCH_HARDWARE_ACCESS_BUS_CPLD;

    default:
      return lemming::dataplane::sai::SWITCH_HARDWARE_ACCESS_BUS_UNSPECIFIED;
  }
}
sai_switch_hardware_access_bus_t
convert_sai_switch_hardware_access_bus_t_to_sai(
    lemming::dataplane::sai::SwitchHardwareAccessBus val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_HARDWARE_ACCESS_BUS_MDIO:
      return SAI_SWITCH_HARDWARE_ACCESS_BUS_MDIO;

    case lemming::dataplane::sai::SWITCH_HARDWARE_ACCESS_BUS_I2C:
      return SAI_SWITCH_HARDWARE_ACCESS_BUS_I2C;

    case lemming::dataplane::sai::SWITCH_HARDWARE_ACCESS_BUS_CPLD:
      return SAI_SWITCH_HARDWARE_ACCESS_BUS_CPLD;

    default:
      return SAI_SWITCH_HARDWARE_ACCESS_BUS_MDIO;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_switch_hardware_access_bus_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_switch_hardware_access_bus_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_hardware_access_bus_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_hardware_access_bus_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchHardwareAccessBus>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchMcastSnoopingCapability
convert_sai_switch_mcast_snooping_capability_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_MCAST_SNOOPING_CAPABILITY_NONE:
      return lemming::dataplane::sai::SWITCH_MCAST_SNOOPING_CAPABILITY_NONE;

    case SAI_SWITCH_MCAST_SNOOPING_CAPABILITY_XG:
      return lemming::dataplane::sai::SWITCH_MCAST_SNOOPING_CAPABILITY_XG;

    case SAI_SWITCH_MCAST_SNOOPING_CAPABILITY_SG:
      return lemming::dataplane::sai::SWITCH_MCAST_SNOOPING_CAPABILITY_SG;

    case SAI_SWITCH_MCAST_SNOOPING_CAPABILITY_XG_AND_SG:
      return lemming::dataplane::sai::
          SWITCH_MCAST_SNOOPING_CAPABILITY_XG_AND_SG;

    default:
      return lemming::dataplane::sai::
          SWITCH_MCAST_SNOOPING_CAPABILITY_UNSPECIFIED;
  }
}
sai_switch_mcast_snooping_capability_t
convert_sai_switch_mcast_snooping_capability_t_to_sai(
    lemming::dataplane::sai::SwitchMcastSnoopingCapability val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_MCAST_SNOOPING_CAPABILITY_NONE:
      return SAI_SWITCH_MCAST_SNOOPING_CAPABILITY_NONE;

    case lemming::dataplane::sai::SWITCH_MCAST_SNOOPING_CAPABILITY_XG:
      return SAI_SWITCH_MCAST_SNOOPING_CAPABILITY_XG;

    case lemming::dataplane::sai::SWITCH_MCAST_SNOOPING_CAPABILITY_SG:
      return SAI_SWITCH_MCAST_SNOOPING_CAPABILITY_SG;

    case lemming::dataplane::sai::SWITCH_MCAST_SNOOPING_CAPABILITY_XG_AND_SG:
      return SAI_SWITCH_MCAST_SNOOPING_CAPABILITY_XG_AND_SG;

    default:
      return SAI_SWITCH_MCAST_SNOOPING_CAPABILITY_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_switch_mcast_snooping_capability_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_switch_mcast_snooping_capability_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_mcast_snooping_capability_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_mcast_snooping_capability_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchMcastSnoopingCapability>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchOperStatus
convert_sai_switch_oper_status_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_OPER_STATUS_UNKNOWN:
      return lemming::dataplane::sai::SWITCH_OPER_STATUS_UNKNOWN;

    case SAI_SWITCH_OPER_STATUS_UP:
      return lemming::dataplane::sai::SWITCH_OPER_STATUS_UP;

    case SAI_SWITCH_OPER_STATUS_DOWN:
      return lemming::dataplane::sai::SWITCH_OPER_STATUS_DOWN;

    case SAI_SWITCH_OPER_STATUS_FAILED:
      return lemming::dataplane::sai::SWITCH_OPER_STATUS_FAILED;

    default:
      return lemming::dataplane::sai::SWITCH_OPER_STATUS_UNSPECIFIED;
  }
}
sai_switch_oper_status_t convert_sai_switch_oper_status_t_to_sai(
    lemming::dataplane::sai::SwitchOperStatus val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_OPER_STATUS_UNKNOWN:
      return SAI_SWITCH_OPER_STATUS_UNKNOWN;

    case lemming::dataplane::sai::SWITCH_OPER_STATUS_UP:
      return SAI_SWITCH_OPER_STATUS_UP;

    case lemming::dataplane::sai::SWITCH_OPER_STATUS_DOWN:
      return SAI_SWITCH_OPER_STATUS_DOWN;

    case lemming::dataplane::sai::SWITCH_OPER_STATUS_FAILED:
      return SAI_SWITCH_OPER_STATUS_FAILED;

    default:
      return SAI_SWITCH_OPER_STATUS_UNKNOWN;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_switch_oper_status_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_switch_oper_status_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_oper_status_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_oper_status_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchOperStatus>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchRestartType
convert_sai_switch_restart_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_RESTART_TYPE_NONE:
      return lemming::dataplane::sai::SWITCH_RESTART_TYPE_NONE;

    case SAI_SWITCH_RESTART_TYPE_PLANNED:
      return lemming::dataplane::sai::SWITCH_RESTART_TYPE_PLANNED;

    case SAI_SWITCH_RESTART_TYPE_ANY:
      return lemming::dataplane::sai::SWITCH_RESTART_TYPE_ANY;

    default:
      return lemming::dataplane::sai::SWITCH_RESTART_TYPE_UNSPECIFIED;
  }
}
sai_switch_restart_type_t convert_sai_switch_restart_type_t_to_sai(
    lemming::dataplane::sai::SwitchRestartType val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_RESTART_TYPE_NONE:
      return SAI_SWITCH_RESTART_TYPE_NONE;

    case lemming::dataplane::sai::SWITCH_RESTART_TYPE_PLANNED:
      return SAI_SWITCH_RESTART_TYPE_PLANNED;

    case lemming::dataplane::sai::SWITCH_RESTART_TYPE_ANY:
      return SAI_SWITCH_RESTART_TYPE_ANY;

    default:
      return SAI_SWITCH_RESTART_TYPE_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_switch_restart_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_switch_restart_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_restart_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_restart_type_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchRestartType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchStat convert_sai_switch_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_STAT_IN_DROP_REASON_RANGE_BASE:
      return lemming::dataplane::sai::SWITCH_STAT_IN_DROP_REASON_RANGE_BASE;

    case SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS;

    case SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS;

    case SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS;

    case SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS;

    case SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS;

    case SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS;

    case SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS;

    case SAI_SWITCH_STAT_IN_DROP_REASON_RANGE_END:
      return lemming::dataplane::sai::SWITCH_STAT_IN_DROP_REASON_RANGE_END;

    case SAI_SWITCH_STAT_OUT_DROP_REASON_RANGE_BASE:
      return lemming::dataplane::sai::SWITCH_STAT_OUT_DROP_REASON_RANGE_BASE;

    case SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS;

    case SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS;

    case SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS;

    case SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS;

    case SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS;

    case SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS;

    case SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS:
      return lemming::dataplane::sai::
          SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS;

    case SAI_SWITCH_STAT_OUT_DROP_REASON_RANGE_END:
      return lemming::dataplane::sai::SWITCH_STAT_OUT_DROP_REASON_RANGE_END;

    case SAI_SWITCH_STAT_FABRIC_DROP_REASON_RANGE_BASE:
      return lemming::dataplane::sai::SWITCH_STAT_FABRIC_DROP_REASON_RANGE_BASE;

    case SAI_SWITCH_STAT_REACHABILITY_DROP:
      return lemming::dataplane::sai::SWITCH_STAT_REACHABILITY_DROP;

    case SAI_SWITCH_STAT_HIGHEST_QUEUE_CONGESTION_LEVEL:
      return lemming::dataplane::sai::
          SWITCH_STAT_HIGHEST_QUEUE_CONGESTION_LEVEL;

    case SAI_SWITCH_STAT_GLOBAL_DROP:
      return lemming::dataplane::sai::SWITCH_STAT_GLOBAL_DROP;

    case SAI_SWITCH_STAT_FABRIC_DROP_REASON_RANGE_END:
      return lemming::dataplane::sai::SWITCH_STAT_FABRIC_DROP_REASON_RANGE_END;

    default:
      return lemming::dataplane::sai::SWITCH_STAT_UNSPECIFIED;
  }
}
sai_switch_stat_t convert_sai_switch_stat_t_to_sai(
    lemming::dataplane::sai::SwitchStat val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_STAT_IN_DROP_REASON_RANGE_BASE:
      return SAI_SWITCH_STAT_IN_DROP_REASON_RANGE_BASE;

    case lemming::dataplane::sai::
        SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS:
      return SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS:
      return SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS:
      return SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS:
      return SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS:
      return SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS:
      return SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS:
      return SAI_SWITCH_STAT_IN_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS;

    case lemming::dataplane::sai::SWITCH_STAT_IN_DROP_REASON_RANGE_END:
      return SAI_SWITCH_STAT_IN_DROP_REASON_RANGE_END;

    case lemming::dataplane::sai::SWITCH_STAT_OUT_DROP_REASON_RANGE_BASE:
      return SAI_SWITCH_STAT_OUT_DROP_REASON_RANGE_BASE;

    case lemming::dataplane::sai::
        SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS:
      return SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_1_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS:
      return SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_2_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS:
      return SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_3_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS:
      return SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_4_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS:
      return SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_5_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS:
      return SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_6_DROPPED_PKTS;

    case lemming::dataplane::sai::
        SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS:
      return SAI_SWITCH_STAT_OUT_CONFIGURED_DROP_REASONS_7_DROPPED_PKTS;

    case lemming::dataplane::sai::SWITCH_STAT_OUT_DROP_REASON_RANGE_END:
      return SAI_SWITCH_STAT_OUT_DROP_REASON_RANGE_END;

    case lemming::dataplane::sai::SWITCH_STAT_FABRIC_DROP_REASON_RANGE_BASE:
      return SAI_SWITCH_STAT_FABRIC_DROP_REASON_RANGE_BASE;

    case lemming::dataplane::sai::SWITCH_STAT_REACHABILITY_DROP:
      return SAI_SWITCH_STAT_REACHABILITY_DROP;

    case lemming::dataplane::sai::SWITCH_STAT_HIGHEST_QUEUE_CONGESTION_LEVEL:
      return SAI_SWITCH_STAT_HIGHEST_QUEUE_CONGESTION_LEVEL;

    case lemming::dataplane::sai::SWITCH_STAT_GLOBAL_DROP:
      return SAI_SWITCH_STAT_GLOBAL_DROP;

    case lemming::dataplane::sai::SWITCH_STAT_FABRIC_DROP_REASON_RANGE_END:
      return SAI_SWITCH_STAT_FABRIC_DROP_REASON_RANGE_END;

    default:
      return SAI_SWITCH_STAT_IN_DROP_REASON_RANGE_BASE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_switch_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_switch_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchSwitchingMode
convert_sai_switch_switching_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_SWITCHING_MODE_CUT_THROUGH:
      return lemming::dataplane::sai::SWITCH_SWITCHING_MODE_CUT_THROUGH;

    case SAI_SWITCH_SWITCHING_MODE_STORE_AND_FORWARD:
      return lemming::dataplane::sai::SWITCH_SWITCHING_MODE_STORE_AND_FORWARD;

    default:
      return lemming::dataplane::sai::SWITCH_SWITCHING_MODE_UNSPECIFIED;
  }
}
sai_switch_switching_mode_t convert_sai_switch_switching_mode_t_to_sai(
    lemming::dataplane::sai::SwitchSwitchingMode val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_SWITCHING_MODE_CUT_THROUGH:
      return SAI_SWITCH_SWITCHING_MODE_CUT_THROUGH;

    case lemming::dataplane::sai::SWITCH_SWITCHING_MODE_STORE_AND_FORWARD:
      return SAI_SWITCH_SWITCHING_MODE_STORE_AND_FORWARD;

    default:
      return SAI_SWITCH_SWITCHING_MODE_CUT_THROUGH;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_switch_switching_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_switch_switching_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_switching_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_switching_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchSwitchingMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchTunnelAttr
convert_sai_switch_tunnel_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_TYPE:
      return lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_TUNNEL_TYPE;

    case SAI_SWITCH_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
      return lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION;

    case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_ENCAP_ECN_MODE:
      return lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_TUNNEL_ENCAP_ECN_MODE;

    case SAI_SWITCH_TUNNEL_ATTR_ENCAP_MAPPERS:
      return lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_ENCAP_MAPPERS;

    case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_DECAP_ECN_MODE:
      return lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_TUNNEL_DECAP_ECN_MODE;

    case SAI_SWITCH_TUNNEL_ATTR_DECAP_MAPPERS:
      return lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_DECAP_MAPPERS;

    case SAI_SWITCH_TUNNEL_ATTR_TUNNEL_VXLAN_UDP_SPORT_MODE:
      return lemming::dataplane::sai::
          SWITCH_TUNNEL_ATTR_TUNNEL_VXLAN_UDP_SPORT_MODE;

    case SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT:
      return lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT;

    case SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
      return lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK;

    case SAI_SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
      return lemming::dataplane::sai::
          SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP;

    case SAI_SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
      return lemming::dataplane::sai::
          SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP;

    case SAI_SWITCH_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
      return lemming::dataplane::sai::
          SWITCH_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP;

    case SAI_SWITCH_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
      return lemming::dataplane::sai::
          SWITCH_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP;

    default:
      return lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_UNSPECIFIED;
  }
}
sai_switch_tunnel_attr_t convert_sai_switch_tunnel_attr_t_to_sai(
    lemming::dataplane::sai::SwitchTunnelAttr val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_TUNNEL_TYPE:
      return SAI_SWITCH_TUNNEL_ATTR_TUNNEL_TYPE;

    case lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
      return SAI_SWITCH_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION;

    case lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_TUNNEL_ENCAP_ECN_MODE:
      return SAI_SWITCH_TUNNEL_ATTR_TUNNEL_ENCAP_ECN_MODE;

    case lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_ENCAP_MAPPERS:
      return SAI_SWITCH_TUNNEL_ATTR_ENCAP_MAPPERS;

    case lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_TUNNEL_DECAP_ECN_MODE:
      return SAI_SWITCH_TUNNEL_ATTR_TUNNEL_DECAP_ECN_MODE;

    case lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_DECAP_MAPPERS:
      return SAI_SWITCH_TUNNEL_ATTR_DECAP_MAPPERS;

    case lemming::dataplane::sai::
        SWITCH_TUNNEL_ATTR_TUNNEL_VXLAN_UDP_SPORT_MODE:
      return SAI_SWITCH_TUNNEL_ATTR_TUNNEL_VXLAN_UDP_SPORT_MODE;

    case lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT:
      return SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT;

    case lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
      return SAI_SWITCH_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK;

    case lemming::dataplane::sai::
        SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
      return SAI_SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP;

    case lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
      return SAI_SWITCH_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP;

    case lemming::dataplane::sai::SWITCH_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
      return SAI_SWITCH_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP;

    case lemming::dataplane::sai::
        SWITCH_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
      return SAI_SWITCH_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP;

    default:
      return SAI_SWITCH_TUNNEL_ATTR_TUNNEL_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_switch_tunnel_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_switch_tunnel_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_tunnel_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_tunnel_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchTunnelAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SwitchType convert_sai_switch_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_SWITCH_TYPE_NPU:
      return lemming::dataplane::sai::SWITCH_TYPE_NPU;

    case SAI_SWITCH_TYPE_PHY:
      return lemming::dataplane::sai::SWITCH_TYPE_PHY;

    case SAI_SWITCH_TYPE_VOQ:
      return lemming::dataplane::sai::SWITCH_TYPE_VOQ;

    case SAI_SWITCH_TYPE_FABRIC:
      return lemming::dataplane::sai::SWITCH_TYPE_FABRIC;

    default:
      return lemming::dataplane::sai::SWITCH_TYPE_UNSPECIFIED;
  }
}
sai_switch_type_t convert_sai_switch_type_t_to_sai(
    lemming::dataplane::sai::SwitchType val) {
  switch (val) {
    case lemming::dataplane::sai::SWITCH_TYPE_NPU:
      return SAI_SWITCH_TYPE_NPU;

    case lemming::dataplane::sai::SWITCH_TYPE_PHY:
      return SAI_SWITCH_TYPE_PHY;

    case lemming::dataplane::sai::SWITCH_TYPE_VOQ:
      return SAI_SWITCH_TYPE_VOQ;

    case lemming::dataplane::sai::SWITCH_TYPE_FABRIC:
      return SAI_SWITCH_TYPE_FABRIC;

    default:
      return SAI_SWITCH_TYPE_NPU;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_switch_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_switch_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_switch_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_switch_type_t_to_sai(
        static_cast<lemming::dataplane::sai::SwitchType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SystemPortAttr convert_sai_system_port_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_SYSTEM_PORT_ATTR_TYPE:
      return lemming::dataplane::sai::SYSTEM_PORT_ATTR_TYPE;

    case SAI_SYSTEM_PORT_ATTR_QOS_NUMBER_OF_VOQS:
      return lemming::dataplane::sai::SYSTEM_PORT_ATTR_QOS_NUMBER_OF_VOQS;

    case SAI_SYSTEM_PORT_ATTR_QOS_VOQ_LIST:
      return lemming::dataplane::sai::SYSTEM_PORT_ATTR_QOS_VOQ_LIST;

    case SAI_SYSTEM_PORT_ATTR_PORT:
      return lemming::dataplane::sai::SYSTEM_PORT_ATTR_PORT;

    case SAI_SYSTEM_PORT_ATTR_ADMIN_STATE:
      return lemming::dataplane::sai::SYSTEM_PORT_ATTR_ADMIN_STATE;

    case SAI_SYSTEM_PORT_ATTR_CONFIG_INFO:
      return lemming::dataplane::sai::SYSTEM_PORT_ATTR_CONFIG_INFO;

    case SAI_SYSTEM_PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
      return lemming::dataplane::sai::SYSTEM_PORT_ATTR_QOS_TC_TO_QUEUE_MAP;

    default:
      return lemming::dataplane::sai::SYSTEM_PORT_ATTR_UNSPECIFIED;
  }
}
sai_system_port_attr_t convert_sai_system_port_attr_t_to_sai(
    lemming::dataplane::sai::SystemPortAttr val) {
  switch (val) {
    case lemming::dataplane::sai::SYSTEM_PORT_ATTR_TYPE:
      return SAI_SYSTEM_PORT_ATTR_TYPE;

    case lemming::dataplane::sai::SYSTEM_PORT_ATTR_QOS_NUMBER_OF_VOQS:
      return SAI_SYSTEM_PORT_ATTR_QOS_NUMBER_OF_VOQS;

    case lemming::dataplane::sai::SYSTEM_PORT_ATTR_QOS_VOQ_LIST:
      return SAI_SYSTEM_PORT_ATTR_QOS_VOQ_LIST;

    case lemming::dataplane::sai::SYSTEM_PORT_ATTR_PORT:
      return SAI_SYSTEM_PORT_ATTR_PORT;

    case lemming::dataplane::sai::SYSTEM_PORT_ATTR_ADMIN_STATE:
      return SAI_SYSTEM_PORT_ATTR_ADMIN_STATE;

    case lemming::dataplane::sai::SYSTEM_PORT_ATTR_CONFIG_INFO:
      return SAI_SYSTEM_PORT_ATTR_CONFIG_INFO;

    case lemming::dataplane::sai::SYSTEM_PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
      return SAI_SYSTEM_PORT_ATTR_QOS_TC_TO_QUEUE_MAP;

    default:
      return SAI_SYSTEM_PORT_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_system_port_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_system_port_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_system_port_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_system_port_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::SystemPortAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::SystemPortType convert_sai_system_port_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_SYSTEM_PORT_TYPE_LOCAL:
      return lemming::dataplane::sai::SYSTEM_PORT_TYPE_LOCAL;

    case SAI_SYSTEM_PORT_TYPE_REMOTE:
      return lemming::dataplane::sai::SYSTEM_PORT_TYPE_REMOTE;

    default:
      return lemming::dataplane::sai::SYSTEM_PORT_TYPE_UNSPECIFIED;
  }
}
sai_system_port_type_t convert_sai_system_port_type_t_to_sai(
    lemming::dataplane::sai::SystemPortType val) {
  switch (val) {
    case lemming::dataplane::sai::SYSTEM_PORT_TYPE_LOCAL:
      return SAI_SYSTEM_PORT_TYPE_LOCAL;

    case lemming::dataplane::sai::SYSTEM_PORT_TYPE_REMOTE:
      return SAI_SYSTEM_PORT_TYPE_REMOTE;

    default:
      return SAI_SYSTEM_PORT_TYPE_LOCAL;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_system_port_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_system_port_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_system_port_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_system_port_type_t_to_sai(
        static_cast<lemming::dataplane::sai::SystemPortType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TableBitmapClassificationEntryAction
convert_sai_table_bitmap_classification_entry_action_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ACTION_SET_METADATA:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_ACTION_SET_METADATA;

    case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ACTION_NOACTION:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_ACTION_NOACTION;

    default:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_ACTION_UNSPECIFIED;
  }
}
sai_table_bitmap_classification_entry_action_t
convert_sai_table_bitmap_classification_entry_action_t_to_sai(
    lemming::dataplane::sai::TableBitmapClassificationEntryAction val) {
  switch (val) {
    case lemming::dataplane::sai::
        TABLE_BITMAP_CLASSIFICATION_ENTRY_ACTION_SET_METADATA:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ACTION_SET_METADATA;

    case lemming::dataplane::sai::
        TABLE_BITMAP_CLASSIFICATION_ENTRY_ACTION_NOACTION:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ACTION_NOACTION;

    default:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ACTION_SET_METADATA;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_classification_entry_action_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_table_bitmap_classification_entry_action_t_to_proto(
            list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_table_bitmap_classification_entry_action_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_table_bitmap_classification_entry_action_t_to_sai(
        static_cast<
            lemming::dataplane::sai::TableBitmapClassificationEntryAction>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TableBitmapClassificationEntryAttr
convert_sai_table_bitmap_classification_entry_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ACTION:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ACTION;

    case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ROUTER_INTERFACE_KEY:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ROUTER_INTERFACE_KEY;

    case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IS_DEFAULT:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IS_DEFAULT;

    case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IN_RIF_METADATA:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IN_RIF_METADATA;

    default:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_table_bitmap_classification_entry_attr_t
convert_sai_table_bitmap_classification_entry_attr_t_to_sai(
    lemming::dataplane::sai::TableBitmapClassificationEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ACTION:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ACTION;

    case lemming::dataplane::sai::
        TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ROUTER_INTERFACE_KEY:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ROUTER_INTERFACE_KEY;

    case lemming::dataplane::sai::
        TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IS_DEFAULT:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IS_DEFAULT;

    case lemming::dataplane::sai::
        TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IN_RIF_METADATA:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IN_RIF_METADATA;

    default:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ACTION;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_classification_entry_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_table_bitmap_classification_entry_attr_t_to_proto(
            list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_table_bitmap_classification_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_table_bitmap_classification_entry_attr_t_to_sai(
        static_cast<
            lemming::dataplane::sai::TableBitmapClassificationEntryAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TableBitmapClassificationEntryStat
convert_sai_table_bitmap_classification_entry_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_STAT_HIT_PACKETS:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_STAT_HIT_PACKETS;

    case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_STAT_HIT_OCTETS:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_STAT_HIT_OCTETS;

    default:
      return lemming::dataplane::sai::
          TABLE_BITMAP_CLASSIFICATION_ENTRY_STAT_UNSPECIFIED;
  }
}
sai_table_bitmap_classification_entry_stat_t
convert_sai_table_bitmap_classification_entry_stat_t_to_sai(
    lemming::dataplane::sai::TableBitmapClassificationEntryStat val) {
  switch (val) {
    case lemming::dataplane::sai::
        TABLE_BITMAP_CLASSIFICATION_ENTRY_STAT_HIT_PACKETS:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_STAT_HIT_PACKETS;

    case lemming::dataplane::sai::
        TABLE_BITMAP_CLASSIFICATION_ENTRY_STAT_HIT_OCTETS:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_STAT_HIT_OCTETS;

    default:
      return SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_STAT_HIT_PACKETS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_classification_entry_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_table_bitmap_classification_entry_stat_t_to_proto(
            list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_table_bitmap_classification_entry_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_table_bitmap_classification_entry_stat_t_to_sai(
        static_cast<
            lemming::dataplane::sai::TableBitmapClassificationEntryStat>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TableBitmapRouterEntryAction
convert_sai_table_bitmap_router_entry_action_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_NEXTHOP:
      return lemming::dataplane::sai::
          TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_NEXTHOP;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_LOCAL:
      return lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_LOCAL;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_CPU:
      return lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_CPU;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_DROP:
      return lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ACTION_DROP;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_NOACTION:
      return lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ACTION_NOACTION;

    default:
      return lemming::dataplane::sai::
          TABLE_BITMAP_ROUTER_ENTRY_ACTION_UNSPECIFIED;
  }
}
sai_table_bitmap_router_entry_action_t
convert_sai_table_bitmap_router_entry_action_t_to_sai(
    lemming::dataplane::sai::TableBitmapRouterEntryAction val) {
  switch (val) {
    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_NEXTHOP:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_NEXTHOP;

    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_LOCAL:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_LOCAL;

    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_CPU:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_CPU;

    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ACTION_DROP:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_DROP;

    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ACTION_NOACTION:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_NOACTION;

    default:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ACTION_TO_NEXTHOP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_router_entry_action_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_table_bitmap_router_entry_action_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_table_bitmap_router_entry_action_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_table_bitmap_router_entry_action_t_to_sai(
        static_cast<lemming::dataplane::sai::TableBitmapRouterEntryAction>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TableBitmapRouterEntryAttr
convert_sai_table_bitmap_router_entry_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_ACTION:
      return lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_ACTION;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_PRIORITY:
      return lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_PRIORITY;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_KEY:
      return lemming::dataplane::sai::
          TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_KEY;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_MASK:
      return lemming::dataplane::sai::
          TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_MASK;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_DST_IP_KEY:
      return lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_DST_IP_KEY;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_TUNNEL_INDEX:
      return lemming::dataplane::sai::
          TABLE_BITMAP_ROUTER_ENTRY_ATTR_TUNNEL_INDEX;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_NEXT_HOP:
      return lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_NEXT_HOP;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_ROUTER_INTERFACE:
      return lemming::dataplane::sai::
          TABLE_BITMAP_ROUTER_ENTRY_ATTR_ROUTER_INTERFACE;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_TRAP_ID:
      return lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_TRAP_ID;

    default:
      return lemming::dataplane::sai::
          TABLE_BITMAP_ROUTER_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_table_bitmap_router_entry_attr_t
convert_sai_table_bitmap_router_entry_attr_t_to_sai(
    lemming::dataplane::sai::TableBitmapRouterEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_ACTION:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_ACTION;

    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_PRIORITY:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_PRIORITY;

    case lemming::dataplane::sai::
        TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_KEY:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_KEY;

    case lemming::dataplane::sai::
        TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_MASK:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_MASK;

    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_DST_IP_KEY:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_DST_IP_KEY;

    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_TUNNEL_INDEX:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_TUNNEL_INDEX;

    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_NEXT_HOP:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_NEXT_HOP;

    case lemming::dataplane::sai::
        TABLE_BITMAP_ROUTER_ENTRY_ATTR_ROUTER_INTERFACE:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_ROUTER_INTERFACE;

    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_ATTR_TRAP_ID:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_TRAP_ID;

    default:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_ACTION;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_router_entry_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_table_bitmap_router_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_table_bitmap_router_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_table_bitmap_router_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TableBitmapRouterEntryAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TableBitmapRouterEntryStat
convert_sai_table_bitmap_router_entry_stat_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TABLE_BITMAP_ROUTER_ENTRY_STAT_HIT_PACKETS:
      return lemming::dataplane::sai::
          TABLE_BITMAP_ROUTER_ENTRY_STAT_HIT_PACKETS;

    case SAI_TABLE_BITMAP_ROUTER_ENTRY_STAT_HIT_OCTETS:
      return lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_STAT_HIT_OCTETS;

    default:
      return lemming::dataplane::sai::
          TABLE_BITMAP_ROUTER_ENTRY_STAT_UNSPECIFIED;
  }
}
sai_table_bitmap_router_entry_stat_t
convert_sai_table_bitmap_router_entry_stat_t_to_sai(
    lemming::dataplane::sai::TableBitmapRouterEntryStat val) {
  switch (val) {
    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_STAT_HIT_PACKETS:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_STAT_HIT_PACKETS;

    case lemming::dataplane::sai::TABLE_BITMAP_ROUTER_ENTRY_STAT_HIT_OCTETS:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_STAT_HIT_OCTETS;

    default:
      return SAI_TABLE_BITMAP_ROUTER_ENTRY_STAT_HIT_PACKETS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_table_bitmap_router_entry_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_table_bitmap_router_entry_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_table_bitmap_router_entry_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_table_bitmap_router_entry_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::TableBitmapRouterEntryStat>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TableMetaTunnelEntryAction
convert_sai_table_meta_tunnel_entry_action_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TABLE_META_TUNNEL_ENTRY_ACTION_TUNNEL_ENCAP:
      return lemming::dataplane::sai::
          TABLE_META_TUNNEL_ENTRY_ACTION_TUNNEL_ENCAP;

    case SAI_TABLE_META_TUNNEL_ENTRY_ACTION_NOACTION:
      return lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ACTION_NOACTION;

    default:
      return lemming::dataplane::sai::
          TABLE_META_TUNNEL_ENTRY_ACTION_UNSPECIFIED;
  }
}
sai_table_meta_tunnel_entry_action_t
convert_sai_table_meta_tunnel_entry_action_t_to_sai(
    lemming::dataplane::sai::TableMetaTunnelEntryAction val) {
  switch (val) {
    case lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ACTION_TUNNEL_ENCAP:
      return SAI_TABLE_META_TUNNEL_ENTRY_ACTION_TUNNEL_ENCAP;

    case lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ACTION_NOACTION:
      return SAI_TABLE_META_TUNNEL_ENTRY_ACTION_NOACTION;

    default:
      return SAI_TABLE_META_TUNNEL_ENTRY_ACTION_TUNNEL_ENCAP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_table_meta_tunnel_entry_action_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_table_meta_tunnel_entry_action_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_table_meta_tunnel_entry_action_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_table_meta_tunnel_entry_action_t_to_sai(
        static_cast<lemming::dataplane::sai::TableMetaTunnelEntryAction>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TableMetaTunnelEntryAttr
convert_sai_table_meta_tunnel_entry_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_ACTION:
      return lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_ACTION;

    case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_METADATA_KEY:
      return lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_METADATA_KEY;

    case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_IS_DEFAULT:
      return lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_IS_DEFAULT;

    case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_TUNNEL_ID:
      return lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_TUNNEL_ID;

    case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_UNDERLAY_DIP:
      return lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_UNDERLAY_DIP;

    default:
      return lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_table_meta_tunnel_entry_attr_t
convert_sai_table_meta_tunnel_entry_attr_t_to_sai(
    lemming::dataplane::sai::TableMetaTunnelEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_ACTION:
      return SAI_TABLE_META_TUNNEL_ENTRY_ATTR_ACTION;

    case lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_METADATA_KEY:
      return SAI_TABLE_META_TUNNEL_ENTRY_ATTR_METADATA_KEY;

    case lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_IS_DEFAULT:
      return SAI_TABLE_META_TUNNEL_ENTRY_ATTR_IS_DEFAULT;

    case lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_TUNNEL_ID:
      return SAI_TABLE_META_TUNNEL_ENTRY_ATTR_TUNNEL_ID;

    case lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_ATTR_UNDERLAY_DIP:
      return SAI_TABLE_META_TUNNEL_ENTRY_ATTR_UNDERLAY_DIP;

    default:
      return SAI_TABLE_META_TUNNEL_ENTRY_ATTR_ACTION;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_table_meta_tunnel_entry_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_table_meta_tunnel_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_table_meta_tunnel_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_table_meta_tunnel_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TableMetaTunnelEntryAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TableMetaTunnelEntryStat
convert_sai_table_meta_tunnel_entry_stat_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TABLE_META_TUNNEL_ENTRY_STAT_HIT_PACKETS:
      return lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_STAT_HIT_PACKETS;

    case SAI_TABLE_META_TUNNEL_ENTRY_STAT_HIT_OCTETS:
      return lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_STAT_HIT_OCTETS;

    default:
      return lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_STAT_UNSPECIFIED;
  }
}
sai_table_meta_tunnel_entry_stat_t
convert_sai_table_meta_tunnel_entry_stat_t_to_sai(
    lemming::dataplane::sai::TableMetaTunnelEntryStat val) {
  switch (val) {
    case lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_STAT_HIT_PACKETS:
      return SAI_TABLE_META_TUNNEL_ENTRY_STAT_HIT_PACKETS;

    case lemming::dataplane::sai::TABLE_META_TUNNEL_ENTRY_STAT_HIT_OCTETS:
      return SAI_TABLE_META_TUNNEL_ENTRY_STAT_HIT_OCTETS;

    default:
      return SAI_TABLE_META_TUNNEL_ENTRY_STAT_HIT_PACKETS;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_table_meta_tunnel_entry_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_table_meta_tunnel_entry_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_table_meta_tunnel_entry_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_table_meta_tunnel_entry_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::TableMetaTunnelEntryStat>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamAttr convert_sai_tam_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_ATTR_TELEMETRY_OBJECTS_LIST:
      return lemming::dataplane::sai::TAM_ATTR_TELEMETRY_OBJECTS_LIST;

    case SAI_TAM_ATTR_EVENT_OBJECTS_LIST:
      return lemming::dataplane::sai::TAM_ATTR_EVENT_OBJECTS_LIST;

    case SAI_TAM_ATTR_INT_OBJECTS_LIST:
      return lemming::dataplane::sai::TAM_ATTR_INT_OBJECTS_LIST;

    case SAI_TAM_ATTR_TAM_BIND_POINT_TYPE_LIST:
      return lemming::dataplane::sai::TAM_ATTR_TAM_BIND_POINT_TYPE_LIST;

    default:
      return lemming::dataplane::sai::TAM_ATTR_UNSPECIFIED;
  }
}
sai_tam_attr_t convert_sai_tam_attr_t_to_sai(
    lemming::dataplane::sai::TamAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_ATTR_TELEMETRY_OBJECTS_LIST:
      return SAI_TAM_ATTR_TELEMETRY_OBJECTS_LIST;

    case lemming::dataplane::sai::TAM_ATTR_EVENT_OBJECTS_LIST:
      return SAI_TAM_ATTR_EVENT_OBJECTS_LIST;

    case lemming::dataplane::sai::TAM_ATTR_INT_OBJECTS_LIST:
      return SAI_TAM_ATTR_INT_OBJECTS_LIST;

    case lemming::dataplane::sai::TAM_ATTR_TAM_BIND_POINT_TYPE_LIST:
      return SAI_TAM_ATTR_TAM_BIND_POINT_TYPE_LIST;

    default:
      return SAI_TAM_ATTR_TELEMETRY_OBJECTS_LIST;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_tam_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamBindPointType
convert_sai_tam_bind_point_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_BIND_POINT_TYPE_QUEUE:
      return lemming::dataplane::sai::TAM_BIND_POINT_TYPE_QUEUE;

    case SAI_TAM_BIND_POINT_TYPE_PORT:
      return lemming::dataplane::sai::TAM_BIND_POINT_TYPE_PORT;

    case SAI_TAM_BIND_POINT_TYPE_LAG:
      return lemming::dataplane::sai::TAM_BIND_POINT_TYPE_LAG;

    case SAI_TAM_BIND_POINT_TYPE_VLAN:
      return lemming::dataplane::sai::TAM_BIND_POINT_TYPE_VLAN;

    case SAI_TAM_BIND_POINT_TYPE_SWITCH:
      return lemming::dataplane::sai::TAM_BIND_POINT_TYPE_SWITCH;

    case SAI_TAM_BIND_POINT_TYPE_IPG:
      return lemming::dataplane::sai::TAM_BIND_POINT_TYPE_IPG;

    case SAI_TAM_BIND_POINT_TYPE_BSP:
      return lemming::dataplane::sai::TAM_BIND_POINT_TYPE_BSP;

    default:
      return lemming::dataplane::sai::TAM_BIND_POINT_TYPE_UNSPECIFIED;
  }
}
sai_tam_bind_point_type_t convert_sai_tam_bind_point_type_t_to_sai(
    lemming::dataplane::sai::TamBindPointType val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_BIND_POINT_TYPE_QUEUE:
      return SAI_TAM_BIND_POINT_TYPE_QUEUE;

    case lemming::dataplane::sai::TAM_BIND_POINT_TYPE_PORT:
      return SAI_TAM_BIND_POINT_TYPE_PORT;

    case lemming::dataplane::sai::TAM_BIND_POINT_TYPE_LAG:
      return SAI_TAM_BIND_POINT_TYPE_LAG;

    case lemming::dataplane::sai::TAM_BIND_POINT_TYPE_VLAN:
      return SAI_TAM_BIND_POINT_TYPE_VLAN;

    case lemming::dataplane::sai::TAM_BIND_POINT_TYPE_SWITCH:
      return SAI_TAM_BIND_POINT_TYPE_SWITCH;

    case lemming::dataplane::sai::TAM_BIND_POINT_TYPE_IPG:
      return SAI_TAM_BIND_POINT_TYPE_IPG;

    case lemming::dataplane::sai::TAM_BIND_POINT_TYPE_BSP:
      return SAI_TAM_BIND_POINT_TYPE_BSP;

    default:
      return SAI_TAM_BIND_POINT_TYPE_QUEUE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_bind_point_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_bind_point_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_bind_point_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_bind_point_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TamBindPointType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamCollectorAttr
convert_sai_tam_collector_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_COLLECTOR_ATTR_SRC_IP:
      return lemming::dataplane::sai::TAM_COLLECTOR_ATTR_SRC_IP;

    case SAI_TAM_COLLECTOR_ATTR_DST_IP:
      return lemming::dataplane::sai::TAM_COLLECTOR_ATTR_DST_IP;

    case SAI_TAM_COLLECTOR_ATTR_LOCALHOST:
      return lemming::dataplane::sai::TAM_COLLECTOR_ATTR_LOCALHOST;

    case SAI_TAM_COLLECTOR_ATTR_VIRTUAL_ROUTER_ID:
      return lemming::dataplane::sai::TAM_COLLECTOR_ATTR_VIRTUAL_ROUTER_ID;

    case SAI_TAM_COLLECTOR_ATTR_TRUNCATE_SIZE:
      return lemming::dataplane::sai::TAM_COLLECTOR_ATTR_TRUNCATE_SIZE;

    case SAI_TAM_COLLECTOR_ATTR_TRANSPORT:
      return lemming::dataplane::sai::TAM_COLLECTOR_ATTR_TRANSPORT;

    case SAI_TAM_COLLECTOR_ATTR_DSCP_VALUE:
      return lemming::dataplane::sai::TAM_COLLECTOR_ATTR_DSCP_VALUE;

    default:
      return lemming::dataplane::sai::TAM_COLLECTOR_ATTR_UNSPECIFIED;
  }
}
sai_tam_collector_attr_t convert_sai_tam_collector_attr_t_to_sai(
    lemming::dataplane::sai::TamCollectorAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_COLLECTOR_ATTR_SRC_IP:
      return SAI_TAM_COLLECTOR_ATTR_SRC_IP;

    case lemming::dataplane::sai::TAM_COLLECTOR_ATTR_DST_IP:
      return SAI_TAM_COLLECTOR_ATTR_DST_IP;

    case lemming::dataplane::sai::TAM_COLLECTOR_ATTR_LOCALHOST:
      return SAI_TAM_COLLECTOR_ATTR_LOCALHOST;

    case lemming::dataplane::sai::TAM_COLLECTOR_ATTR_VIRTUAL_ROUTER_ID:
      return SAI_TAM_COLLECTOR_ATTR_VIRTUAL_ROUTER_ID;

    case lemming::dataplane::sai::TAM_COLLECTOR_ATTR_TRUNCATE_SIZE:
      return SAI_TAM_COLLECTOR_ATTR_TRUNCATE_SIZE;

    case lemming::dataplane::sai::TAM_COLLECTOR_ATTR_TRANSPORT:
      return SAI_TAM_COLLECTOR_ATTR_TRANSPORT;

    case lemming::dataplane::sai::TAM_COLLECTOR_ATTR_DSCP_VALUE:
      return SAI_TAM_COLLECTOR_ATTR_DSCP_VALUE;

    default:
      return SAI_TAM_COLLECTOR_ATTR_SRC_IP;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_collector_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_collector_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_collector_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_collector_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamCollectorAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamEventActionAttr
convert_sai_tam_event_action_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_EVENT_ACTION_ATTR_REPORT_TYPE:
      return lemming::dataplane::sai::TAM_EVENT_ACTION_ATTR_REPORT_TYPE;

    case SAI_TAM_EVENT_ACTION_ATTR_QOS_ACTION_TYPE:
      return lemming::dataplane::sai::TAM_EVENT_ACTION_ATTR_QOS_ACTION_TYPE;

    default:
      return lemming::dataplane::sai::TAM_EVENT_ACTION_ATTR_UNSPECIFIED;
  }
}
sai_tam_event_action_attr_t convert_sai_tam_event_action_attr_t_to_sai(
    lemming::dataplane::sai::TamEventActionAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_EVENT_ACTION_ATTR_REPORT_TYPE:
      return SAI_TAM_EVENT_ACTION_ATTR_REPORT_TYPE;

    case lemming::dataplane::sai::TAM_EVENT_ACTION_ATTR_QOS_ACTION_TYPE:
      return SAI_TAM_EVENT_ACTION_ATTR_QOS_ACTION_TYPE;

    default:
      return SAI_TAM_EVENT_ACTION_ATTR_REPORT_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_event_action_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_event_action_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_event_action_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_event_action_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamEventActionAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamEventAttr convert_sai_tam_event_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_EVENT_ATTR_TYPE:
      return lemming::dataplane::sai::TAM_EVENT_ATTR_TYPE;

    case SAI_TAM_EVENT_ATTR_ACTION_LIST:
      return lemming::dataplane::sai::TAM_EVENT_ATTR_ACTION_LIST;

    case SAI_TAM_EVENT_ATTR_COLLECTOR_LIST:
      return lemming::dataplane::sai::TAM_EVENT_ATTR_COLLECTOR_LIST;

    case SAI_TAM_EVENT_ATTR_THRESHOLD:
      return lemming::dataplane::sai::TAM_EVENT_ATTR_THRESHOLD;

    case SAI_TAM_EVENT_ATTR_DSCP_VALUE:
      return lemming::dataplane::sai::TAM_EVENT_ATTR_DSCP_VALUE;

    default:
      return lemming::dataplane::sai::TAM_EVENT_ATTR_UNSPECIFIED;
  }
}
sai_tam_event_attr_t convert_sai_tam_event_attr_t_to_sai(
    lemming::dataplane::sai::TamEventAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_EVENT_ATTR_TYPE:
      return SAI_TAM_EVENT_ATTR_TYPE;

    case lemming::dataplane::sai::TAM_EVENT_ATTR_ACTION_LIST:
      return SAI_TAM_EVENT_ATTR_ACTION_LIST;

    case lemming::dataplane::sai::TAM_EVENT_ATTR_COLLECTOR_LIST:
      return SAI_TAM_EVENT_ATTR_COLLECTOR_LIST;

    case lemming::dataplane::sai::TAM_EVENT_ATTR_THRESHOLD:
      return SAI_TAM_EVENT_ATTR_THRESHOLD;

    case lemming::dataplane::sai::TAM_EVENT_ATTR_DSCP_VALUE:
      return SAI_TAM_EVENT_ATTR_DSCP_VALUE;

    default:
      return SAI_TAM_EVENT_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_tam_event_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_event_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_event_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_event_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamEventAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamEventThresholdAttr
convert_sai_tam_event_threshold_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_EVENT_THRESHOLD_ATTR_HIGH_WATERMARK:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_HIGH_WATERMARK;

    case SAI_TAM_EVENT_THRESHOLD_ATTR_LOW_WATERMARK:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_LOW_WATERMARK;

    case SAI_TAM_EVENT_THRESHOLD_ATTR_LATENCY:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_LATENCY;

    case SAI_TAM_EVENT_THRESHOLD_ATTR_RATE:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_RATE;

    case SAI_TAM_EVENT_THRESHOLD_ATTR_ABS_VALUE:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_ABS_VALUE;

    case SAI_TAM_EVENT_THRESHOLD_ATTR_UNIT:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_UNIT;

    default:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_UNSPECIFIED;
  }
}
sai_tam_event_threshold_attr_t convert_sai_tam_event_threshold_attr_t_to_sai(
    lemming::dataplane::sai::TamEventThresholdAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_HIGH_WATERMARK:
      return SAI_TAM_EVENT_THRESHOLD_ATTR_HIGH_WATERMARK;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_LOW_WATERMARK:
      return SAI_TAM_EVENT_THRESHOLD_ATTR_LOW_WATERMARK;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_LATENCY:
      return SAI_TAM_EVENT_THRESHOLD_ATTR_LATENCY;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_RATE:
      return SAI_TAM_EVENT_THRESHOLD_ATTR_RATE;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_ABS_VALUE:
      return SAI_TAM_EVENT_THRESHOLD_ATTR_ABS_VALUE;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_ATTR_UNIT:
      return SAI_TAM_EVENT_THRESHOLD_ATTR_UNIT;

    default:
      return SAI_TAM_EVENT_THRESHOLD_ATTR_HIGH_WATERMARK;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_event_threshold_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_tam_event_threshold_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_event_threshold_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_event_threshold_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamEventThresholdAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamEventThresholdUnit
convert_sai_tam_event_threshold_unit_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_EVENT_THRESHOLD_UNIT_NANOSEC:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_NANOSEC;

    case SAI_TAM_EVENT_THRESHOLD_UNIT_USEC:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_USEC;

    case SAI_TAM_EVENT_THRESHOLD_UNIT_MSEC:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_MSEC;

    case SAI_TAM_EVENT_THRESHOLD_UNIT_PERCENT:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_PERCENT;

    case SAI_TAM_EVENT_THRESHOLD_UNIT_BYTES:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_BYTES;

    case SAI_TAM_EVENT_THRESHOLD_UNIT_PACKETS:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_PACKETS;

    case SAI_TAM_EVENT_THRESHOLD_UNIT_CELLS:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_CELLS;

    default:
      return lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_UNSPECIFIED;
  }
}
sai_tam_event_threshold_unit_t convert_sai_tam_event_threshold_unit_t_to_sai(
    lemming::dataplane::sai::TamEventThresholdUnit val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_NANOSEC:
      return SAI_TAM_EVENT_THRESHOLD_UNIT_NANOSEC;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_USEC:
      return SAI_TAM_EVENT_THRESHOLD_UNIT_USEC;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_MSEC:
      return SAI_TAM_EVENT_THRESHOLD_UNIT_MSEC;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_PERCENT:
      return SAI_TAM_EVENT_THRESHOLD_UNIT_PERCENT;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_BYTES:
      return SAI_TAM_EVENT_THRESHOLD_UNIT_BYTES;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_PACKETS:
      return SAI_TAM_EVENT_THRESHOLD_UNIT_PACKETS;

    case lemming::dataplane::sai::TAM_EVENT_THRESHOLD_UNIT_CELLS:
      return SAI_TAM_EVENT_THRESHOLD_UNIT_CELLS;

    default:
      return SAI_TAM_EVENT_THRESHOLD_UNIT_NANOSEC;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_event_threshold_unit_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_tam_event_threshold_unit_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_event_threshold_unit_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_event_threshold_unit_t_to_sai(
        static_cast<lemming::dataplane::sai::TamEventThresholdUnit>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamEventType convert_sai_tam_event_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_EVENT_TYPE_FLOW_STATE:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_FLOW_STATE;

    case SAI_TAM_EVENT_TYPE_FLOW_WATCHLIST:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_FLOW_WATCHLIST;

    case SAI_TAM_EVENT_TYPE_FLOW_TCPFLAG:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_FLOW_TCPFLAG;

    case SAI_TAM_EVENT_TYPE_QUEUE_THRESHOLD:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_QUEUE_THRESHOLD;

    case SAI_TAM_EVENT_TYPE_QUEUE_TAIL_DROP:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_QUEUE_TAIL_DROP;

    case SAI_TAM_EVENT_TYPE_PACKET_DROP:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_PACKET_DROP;

    case SAI_TAM_EVENT_TYPE_RESOURCE_UTILIZATION:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_RESOURCE_UTILIZATION;

    case SAI_TAM_EVENT_TYPE_IPG_SHARED:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_IPG_SHARED;

    case SAI_TAM_EVENT_TYPE_IPG_XOFF_ROOM:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_IPG_XOFF_ROOM;

    case SAI_TAM_EVENT_TYPE_BSP:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_BSP;

    default:
      return lemming::dataplane::sai::TAM_EVENT_TYPE_UNSPECIFIED;
  }
}
sai_tam_event_type_t convert_sai_tam_event_type_t_to_sai(
    lemming::dataplane::sai::TamEventType val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_EVENT_TYPE_FLOW_STATE:
      return SAI_TAM_EVENT_TYPE_FLOW_STATE;

    case lemming::dataplane::sai::TAM_EVENT_TYPE_FLOW_WATCHLIST:
      return SAI_TAM_EVENT_TYPE_FLOW_WATCHLIST;

    case lemming::dataplane::sai::TAM_EVENT_TYPE_FLOW_TCPFLAG:
      return SAI_TAM_EVENT_TYPE_FLOW_TCPFLAG;

    case lemming::dataplane::sai::TAM_EVENT_TYPE_QUEUE_THRESHOLD:
      return SAI_TAM_EVENT_TYPE_QUEUE_THRESHOLD;

    case lemming::dataplane::sai::TAM_EVENT_TYPE_QUEUE_TAIL_DROP:
      return SAI_TAM_EVENT_TYPE_QUEUE_TAIL_DROP;

    case lemming::dataplane::sai::TAM_EVENT_TYPE_PACKET_DROP:
      return SAI_TAM_EVENT_TYPE_PACKET_DROP;

    case lemming::dataplane::sai::TAM_EVENT_TYPE_RESOURCE_UTILIZATION:
      return SAI_TAM_EVENT_TYPE_RESOURCE_UTILIZATION;

    case lemming::dataplane::sai::TAM_EVENT_TYPE_IPG_SHARED:
      return SAI_TAM_EVENT_TYPE_IPG_SHARED;

    case lemming::dataplane::sai::TAM_EVENT_TYPE_IPG_XOFF_ROOM:
      return SAI_TAM_EVENT_TYPE_IPG_XOFF_ROOM;

    case lemming::dataplane::sai::TAM_EVENT_TYPE_BSP:
      return SAI_TAM_EVENT_TYPE_BSP;

    default:
      return SAI_TAM_EVENT_TYPE_FLOW_STATE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_tam_event_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_event_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_event_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_event_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TamEventType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamIntAttr convert_sai_tam_int_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_INT_ATTR_TYPE:
      return lemming::dataplane::sai::TAM_INT_ATTR_TYPE;

    case SAI_TAM_INT_ATTR_DEVICE_ID:
      return lemming::dataplane::sai::TAM_INT_ATTR_DEVICE_ID;

    case SAI_TAM_INT_ATTR_IOAM_TRACE_TYPE:
      return lemming::dataplane::sai::TAM_INT_ATTR_IOAM_TRACE_TYPE;

    case SAI_TAM_INT_ATTR_INT_PRESENCE_TYPE:
      return lemming::dataplane::sai::TAM_INT_ATTR_INT_PRESENCE_TYPE;

    case SAI_TAM_INT_ATTR_INT_PRESENCE_PB1:
      return lemming::dataplane::sai::TAM_INT_ATTR_INT_PRESENCE_PB1;

    case SAI_TAM_INT_ATTR_INT_PRESENCE_PB2:
      return lemming::dataplane::sai::TAM_INT_ATTR_INT_PRESENCE_PB2;

    case SAI_TAM_INT_ATTR_INT_PRESENCE_DSCP_VALUE:
      return lemming::dataplane::sai::TAM_INT_ATTR_INT_PRESENCE_DSCP_VALUE;

    case SAI_TAM_INT_ATTR_INLINE:
      return lemming::dataplane::sai::TAM_INT_ATTR_INLINE;

    case SAI_TAM_INT_ATTR_INT_PRESENCE_L3_PROTOCOL:
      return lemming::dataplane::sai::TAM_INT_ATTR_INT_PRESENCE_L3_PROTOCOL;

    case SAI_TAM_INT_ATTR_TRACE_VECTOR:
      return lemming::dataplane::sai::TAM_INT_ATTR_TRACE_VECTOR;

    case SAI_TAM_INT_ATTR_ACTION_VECTOR:
      return lemming::dataplane::sai::TAM_INT_ATTR_ACTION_VECTOR;

    case SAI_TAM_INT_ATTR_P4_INT_INSTRUCTION_BITMAP:
      return lemming::dataplane::sai::TAM_INT_ATTR_P4_INT_INSTRUCTION_BITMAP;

    case SAI_TAM_INT_ATTR_METADATA_FRAGMENT_ENABLE:
      return lemming::dataplane::sai::TAM_INT_ATTR_METADATA_FRAGMENT_ENABLE;

    case SAI_TAM_INT_ATTR_METADATA_CHECKSUM_ENABLE:
      return lemming::dataplane::sai::TAM_INT_ATTR_METADATA_CHECKSUM_ENABLE;

    case SAI_TAM_INT_ATTR_REPORT_ALL_PACKETS:
      return lemming::dataplane::sai::TAM_INT_ATTR_REPORT_ALL_PACKETS;

    case SAI_TAM_INT_ATTR_FLOW_LIVENESS_PERIOD:
      return lemming::dataplane::sai::TAM_INT_ATTR_FLOW_LIVENESS_PERIOD;

    case SAI_TAM_INT_ATTR_LATENCY_SENSITIVITY:
      return lemming::dataplane::sai::TAM_INT_ATTR_LATENCY_SENSITIVITY;

    case SAI_TAM_INT_ATTR_ACL_GROUP:
      return lemming::dataplane::sai::TAM_INT_ATTR_ACL_GROUP;

    case SAI_TAM_INT_ATTR_MAX_HOP_COUNT:
      return lemming::dataplane::sai::TAM_INT_ATTR_MAX_HOP_COUNT;

    case SAI_TAM_INT_ATTR_MAX_LENGTH:
      return lemming::dataplane::sai::TAM_INT_ATTR_MAX_LENGTH;

    case SAI_TAM_INT_ATTR_NAME_SPACE_ID:
      return lemming::dataplane::sai::TAM_INT_ATTR_NAME_SPACE_ID;

    case SAI_TAM_INT_ATTR_NAME_SPACE_ID_GLOBAL:
      return lemming::dataplane::sai::TAM_INT_ATTR_NAME_SPACE_ID_GLOBAL;

    case SAI_TAM_INT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
      return lemming::dataplane::sai::TAM_INT_ATTR_INGRESS_SAMPLEPACKET_ENABLE;

    case SAI_TAM_INT_ATTR_COLLECTOR_LIST:
      return lemming::dataplane::sai::TAM_INT_ATTR_COLLECTOR_LIST;

    case SAI_TAM_INT_ATTR_MATH_FUNC:
      return lemming::dataplane::sai::TAM_INT_ATTR_MATH_FUNC;

    case SAI_TAM_INT_ATTR_REPORT_ID:
      return lemming::dataplane::sai::TAM_INT_ATTR_REPORT_ID;

    default:
      return lemming::dataplane::sai::TAM_INT_ATTR_UNSPECIFIED;
  }
}
sai_tam_int_attr_t convert_sai_tam_int_attr_t_to_sai(
    lemming::dataplane::sai::TamIntAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_INT_ATTR_TYPE:
      return SAI_TAM_INT_ATTR_TYPE;

    case lemming::dataplane::sai::TAM_INT_ATTR_DEVICE_ID:
      return SAI_TAM_INT_ATTR_DEVICE_ID;

    case lemming::dataplane::sai::TAM_INT_ATTR_IOAM_TRACE_TYPE:
      return SAI_TAM_INT_ATTR_IOAM_TRACE_TYPE;

    case lemming::dataplane::sai::TAM_INT_ATTR_INT_PRESENCE_TYPE:
      return SAI_TAM_INT_ATTR_INT_PRESENCE_TYPE;

    case lemming::dataplane::sai::TAM_INT_ATTR_INT_PRESENCE_PB1:
      return SAI_TAM_INT_ATTR_INT_PRESENCE_PB1;

    case lemming::dataplane::sai::TAM_INT_ATTR_INT_PRESENCE_PB2:
      return SAI_TAM_INT_ATTR_INT_PRESENCE_PB2;

    case lemming::dataplane::sai::TAM_INT_ATTR_INT_PRESENCE_DSCP_VALUE:
      return SAI_TAM_INT_ATTR_INT_PRESENCE_DSCP_VALUE;

    case lemming::dataplane::sai::TAM_INT_ATTR_INLINE:
      return SAI_TAM_INT_ATTR_INLINE;

    case lemming::dataplane::sai::TAM_INT_ATTR_INT_PRESENCE_L3_PROTOCOL:
      return SAI_TAM_INT_ATTR_INT_PRESENCE_L3_PROTOCOL;

    case lemming::dataplane::sai::TAM_INT_ATTR_TRACE_VECTOR:
      return SAI_TAM_INT_ATTR_TRACE_VECTOR;

    case lemming::dataplane::sai::TAM_INT_ATTR_ACTION_VECTOR:
      return SAI_TAM_INT_ATTR_ACTION_VECTOR;

    case lemming::dataplane::sai::TAM_INT_ATTR_P4_INT_INSTRUCTION_BITMAP:
      return SAI_TAM_INT_ATTR_P4_INT_INSTRUCTION_BITMAP;

    case lemming::dataplane::sai::TAM_INT_ATTR_METADATA_FRAGMENT_ENABLE:
      return SAI_TAM_INT_ATTR_METADATA_FRAGMENT_ENABLE;

    case lemming::dataplane::sai::TAM_INT_ATTR_METADATA_CHECKSUM_ENABLE:
      return SAI_TAM_INT_ATTR_METADATA_CHECKSUM_ENABLE;

    case lemming::dataplane::sai::TAM_INT_ATTR_REPORT_ALL_PACKETS:
      return SAI_TAM_INT_ATTR_REPORT_ALL_PACKETS;

    case lemming::dataplane::sai::TAM_INT_ATTR_FLOW_LIVENESS_PERIOD:
      return SAI_TAM_INT_ATTR_FLOW_LIVENESS_PERIOD;

    case lemming::dataplane::sai::TAM_INT_ATTR_LATENCY_SENSITIVITY:
      return SAI_TAM_INT_ATTR_LATENCY_SENSITIVITY;

    case lemming::dataplane::sai::TAM_INT_ATTR_ACL_GROUP:
      return SAI_TAM_INT_ATTR_ACL_GROUP;

    case lemming::dataplane::sai::TAM_INT_ATTR_MAX_HOP_COUNT:
      return SAI_TAM_INT_ATTR_MAX_HOP_COUNT;

    case lemming::dataplane::sai::TAM_INT_ATTR_MAX_LENGTH:
      return SAI_TAM_INT_ATTR_MAX_LENGTH;

    case lemming::dataplane::sai::TAM_INT_ATTR_NAME_SPACE_ID:
      return SAI_TAM_INT_ATTR_NAME_SPACE_ID;

    case lemming::dataplane::sai::TAM_INT_ATTR_NAME_SPACE_ID_GLOBAL:
      return SAI_TAM_INT_ATTR_NAME_SPACE_ID_GLOBAL;

    case lemming::dataplane::sai::TAM_INT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
      return SAI_TAM_INT_ATTR_INGRESS_SAMPLEPACKET_ENABLE;

    case lemming::dataplane::sai::TAM_INT_ATTR_COLLECTOR_LIST:
      return SAI_TAM_INT_ATTR_COLLECTOR_LIST;

    case lemming::dataplane::sai::TAM_INT_ATTR_MATH_FUNC:
      return SAI_TAM_INT_ATTR_MATH_FUNC;

    case lemming::dataplane::sai::TAM_INT_ATTR_REPORT_ID:
      return SAI_TAM_INT_ATTR_REPORT_ID;

    default:
      return SAI_TAM_INT_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_tam_int_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_int_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_int_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_int_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamIntAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamIntPresenceType
convert_sai_tam_int_presence_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_INT_PRESENCE_TYPE_UNDEFINED:
      return lemming::dataplane::sai::TAM_INT_PRESENCE_TYPE_UNDEFINED;

    case SAI_TAM_INT_PRESENCE_TYPE_PB:
      return lemming::dataplane::sai::TAM_INT_PRESENCE_TYPE_PB;

    case SAI_TAM_INT_PRESENCE_TYPE_L3_PROTOCOL:
      return lemming::dataplane::sai::TAM_INT_PRESENCE_TYPE_L3_PROTOCOL;

    case SAI_TAM_INT_PRESENCE_TYPE_DSCP:
      return lemming::dataplane::sai::TAM_INT_PRESENCE_TYPE_DSCP;

    default:
      return lemming::dataplane::sai::TAM_INT_PRESENCE_TYPE_UNSPECIFIED;
  }
}
sai_tam_int_presence_type_t convert_sai_tam_int_presence_type_t_to_sai(
    lemming::dataplane::sai::TamIntPresenceType val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_INT_PRESENCE_TYPE_UNDEFINED:
      return SAI_TAM_INT_PRESENCE_TYPE_UNDEFINED;

    case lemming::dataplane::sai::TAM_INT_PRESENCE_TYPE_PB:
      return SAI_TAM_INT_PRESENCE_TYPE_PB;

    case lemming::dataplane::sai::TAM_INT_PRESENCE_TYPE_L3_PROTOCOL:
      return SAI_TAM_INT_PRESENCE_TYPE_L3_PROTOCOL;

    case lemming::dataplane::sai::TAM_INT_PRESENCE_TYPE_DSCP:
      return SAI_TAM_INT_PRESENCE_TYPE_DSCP;

    default:
      return SAI_TAM_INT_PRESENCE_TYPE_UNDEFINED;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_int_presence_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_int_presence_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_int_presence_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_int_presence_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TamIntPresenceType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamIntType convert_sai_tam_int_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_INT_TYPE_IOAM:
      return lemming::dataplane::sai::TAM_INT_TYPE_IOAM;

    case SAI_TAM_INT_TYPE_IFA1:
      return lemming::dataplane::sai::TAM_INT_TYPE_IFA1;

    case SAI_TAM_INT_TYPE_IFA2:
      return lemming::dataplane::sai::TAM_INT_TYPE_IFA2;

    case SAI_TAM_INT_TYPE_P4_INT_1:
      return lemming::dataplane::sai::TAM_INT_TYPE_P4_INT_1;

    case SAI_TAM_INT_TYPE_P4_INT_2:
      return lemming::dataplane::sai::TAM_INT_TYPE_P4_INT_2;

    case SAI_TAM_INT_TYPE_DIRECT_EXPORT:
      return lemming::dataplane::sai::TAM_INT_TYPE_DIRECT_EXPORT;

    case SAI_TAM_INT_TYPE_IFA1_TAILSTAMP:
      return lemming::dataplane::sai::TAM_INT_TYPE_IFA1_TAILSTAMP;

    default:
      return lemming::dataplane::sai::TAM_INT_TYPE_UNSPECIFIED;
  }
}
sai_tam_int_type_t convert_sai_tam_int_type_t_to_sai(
    lemming::dataplane::sai::TamIntType val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_INT_TYPE_IOAM:
      return SAI_TAM_INT_TYPE_IOAM;

    case lemming::dataplane::sai::TAM_INT_TYPE_IFA1:
      return SAI_TAM_INT_TYPE_IFA1;

    case lemming::dataplane::sai::TAM_INT_TYPE_IFA2:
      return SAI_TAM_INT_TYPE_IFA2;

    case lemming::dataplane::sai::TAM_INT_TYPE_P4_INT_1:
      return SAI_TAM_INT_TYPE_P4_INT_1;

    case lemming::dataplane::sai::TAM_INT_TYPE_P4_INT_2:
      return SAI_TAM_INT_TYPE_P4_INT_2;

    case lemming::dataplane::sai::TAM_INT_TYPE_DIRECT_EXPORT:
      return SAI_TAM_INT_TYPE_DIRECT_EXPORT;

    case lemming::dataplane::sai::TAM_INT_TYPE_IFA1_TAILSTAMP:
      return SAI_TAM_INT_TYPE_IFA1_TAILSTAMP;

    default:
      return SAI_TAM_INT_TYPE_IOAM;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_tam_int_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_int_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_int_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_int_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TamIntType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamMathFuncAttr
convert_sai_tam_math_func_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_MATH_FUNC_ATTR_TAM_TEL_MATH_FUNC_TYPE:
      return lemming::dataplane::sai::TAM_MATH_FUNC_ATTR_TAM_TEL_MATH_FUNC_TYPE;

    default:
      return lemming::dataplane::sai::TAM_MATH_FUNC_ATTR_UNSPECIFIED;
  }
}
sai_tam_math_func_attr_t convert_sai_tam_math_func_attr_t_to_sai(
    lemming::dataplane::sai::TamMathFuncAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_MATH_FUNC_ATTR_TAM_TEL_MATH_FUNC_TYPE:
      return SAI_TAM_MATH_FUNC_ATTR_TAM_TEL_MATH_FUNC_TYPE;

    default:
      return SAI_TAM_MATH_FUNC_ATTR_TAM_TEL_MATH_FUNC_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_math_func_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_math_func_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_math_func_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_math_func_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamMathFuncAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamReportAttr convert_sai_tam_report_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_REPORT_ATTR_TYPE:
      return lemming::dataplane::sai::TAM_REPORT_ATTR_TYPE;

    case SAI_TAM_REPORT_ATTR_HISTOGRAM_NUMBER_OF_BINS:
      return lemming::dataplane::sai::TAM_REPORT_ATTR_HISTOGRAM_NUMBER_OF_BINS;

    case SAI_TAM_REPORT_ATTR_HISTOGRAM_BIN_BOUNDARY:
      return lemming::dataplane::sai::TAM_REPORT_ATTR_HISTOGRAM_BIN_BOUNDARY;

    case SAI_TAM_REPORT_ATTR_QUOTA:
      return lemming::dataplane::sai::TAM_REPORT_ATTR_QUOTA;

    case SAI_TAM_REPORT_ATTR_REPORT_MODE:
      return lemming::dataplane::sai::TAM_REPORT_ATTR_REPORT_MODE;

    case SAI_TAM_REPORT_ATTR_REPORT_INTERVAL:
      return lemming::dataplane::sai::TAM_REPORT_ATTR_REPORT_INTERVAL;

    case SAI_TAM_REPORT_ATTR_ENTERPRISE_NUMBER:
      return lemming::dataplane::sai::TAM_REPORT_ATTR_ENTERPRISE_NUMBER;

    case SAI_TAM_REPORT_ATTR_TEMPLATE_REPORT_INTERVAL:
      return lemming::dataplane::sai::TAM_REPORT_ATTR_TEMPLATE_REPORT_INTERVAL;

    default:
      return lemming::dataplane::sai::TAM_REPORT_ATTR_UNSPECIFIED;
  }
}
sai_tam_report_attr_t convert_sai_tam_report_attr_t_to_sai(
    lemming::dataplane::sai::TamReportAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_REPORT_ATTR_TYPE:
      return SAI_TAM_REPORT_ATTR_TYPE;

    case lemming::dataplane::sai::TAM_REPORT_ATTR_HISTOGRAM_NUMBER_OF_BINS:
      return SAI_TAM_REPORT_ATTR_HISTOGRAM_NUMBER_OF_BINS;

    case lemming::dataplane::sai::TAM_REPORT_ATTR_HISTOGRAM_BIN_BOUNDARY:
      return SAI_TAM_REPORT_ATTR_HISTOGRAM_BIN_BOUNDARY;

    case lemming::dataplane::sai::TAM_REPORT_ATTR_QUOTA:
      return SAI_TAM_REPORT_ATTR_QUOTA;

    case lemming::dataplane::sai::TAM_REPORT_ATTR_REPORT_MODE:
      return SAI_TAM_REPORT_ATTR_REPORT_MODE;

    case lemming::dataplane::sai::TAM_REPORT_ATTR_REPORT_INTERVAL:
      return SAI_TAM_REPORT_ATTR_REPORT_INTERVAL;

    case lemming::dataplane::sai::TAM_REPORT_ATTR_ENTERPRISE_NUMBER:
      return SAI_TAM_REPORT_ATTR_ENTERPRISE_NUMBER;

    case lemming::dataplane::sai::TAM_REPORT_ATTR_TEMPLATE_REPORT_INTERVAL:
      return SAI_TAM_REPORT_ATTR_TEMPLATE_REPORT_INTERVAL;

    default:
      return SAI_TAM_REPORT_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_report_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_report_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_report_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_report_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamReportAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamReportMode convert_sai_tam_report_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_REPORT_MODE_ALL:
      return lemming::dataplane::sai::TAM_REPORT_MODE_ALL;

    case SAI_TAM_REPORT_MODE_BULK:
      return lemming::dataplane::sai::TAM_REPORT_MODE_BULK;

    default:
      return lemming::dataplane::sai::TAM_REPORT_MODE_UNSPECIFIED;
  }
}
sai_tam_report_mode_t convert_sai_tam_report_mode_t_to_sai(
    lemming::dataplane::sai::TamReportMode val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_REPORT_MODE_ALL:
      return SAI_TAM_REPORT_MODE_ALL;

    case lemming::dataplane::sai::TAM_REPORT_MODE_BULK:
      return SAI_TAM_REPORT_MODE_BULK;

    default:
      return SAI_TAM_REPORT_MODE_ALL;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_report_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_report_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_report_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_report_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::TamReportMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamReportType convert_sai_tam_report_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_REPORT_TYPE_SFLOW:
      return lemming::dataplane::sai::TAM_REPORT_TYPE_SFLOW;

    case SAI_TAM_REPORT_TYPE_IPFIX:
      return lemming::dataplane::sai::TAM_REPORT_TYPE_IPFIX;

    case SAI_TAM_REPORT_TYPE_PROTO:
      return lemming::dataplane::sai::TAM_REPORT_TYPE_PROTO;

    case SAI_TAM_REPORT_TYPE_THRIFT:
      return lemming::dataplane::sai::TAM_REPORT_TYPE_THRIFT;

    case SAI_TAM_REPORT_TYPE_JSON:
      return lemming::dataplane::sai::TAM_REPORT_TYPE_JSON;

    case SAI_TAM_REPORT_TYPE_P4_EXTN:
      return lemming::dataplane::sai::TAM_REPORT_TYPE_P4_EXTN;

    case SAI_TAM_REPORT_TYPE_HISTOGRAM:
      return lemming::dataplane::sai::TAM_REPORT_TYPE_HISTOGRAM;

    case SAI_TAM_REPORT_TYPE_VENDOR_EXTN:
      return lemming::dataplane::sai::TAM_REPORT_TYPE_VENDOR_EXTN;

    default:
      return lemming::dataplane::sai::TAM_REPORT_TYPE_UNSPECIFIED;
  }
}
sai_tam_report_type_t convert_sai_tam_report_type_t_to_sai(
    lemming::dataplane::sai::TamReportType val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_REPORT_TYPE_SFLOW:
      return SAI_TAM_REPORT_TYPE_SFLOW;

    case lemming::dataplane::sai::TAM_REPORT_TYPE_IPFIX:
      return SAI_TAM_REPORT_TYPE_IPFIX;

    case lemming::dataplane::sai::TAM_REPORT_TYPE_PROTO:
      return SAI_TAM_REPORT_TYPE_PROTO;

    case lemming::dataplane::sai::TAM_REPORT_TYPE_THRIFT:
      return SAI_TAM_REPORT_TYPE_THRIFT;

    case lemming::dataplane::sai::TAM_REPORT_TYPE_JSON:
      return SAI_TAM_REPORT_TYPE_JSON;

    case lemming::dataplane::sai::TAM_REPORT_TYPE_P4_EXTN:
      return SAI_TAM_REPORT_TYPE_P4_EXTN;

    case lemming::dataplane::sai::TAM_REPORT_TYPE_HISTOGRAM:
      return SAI_TAM_REPORT_TYPE_HISTOGRAM;

    case lemming::dataplane::sai::TAM_REPORT_TYPE_VENDOR_EXTN:
      return SAI_TAM_REPORT_TYPE_VENDOR_EXTN;

    default:
      return SAI_TAM_REPORT_TYPE_SFLOW;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_report_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_report_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_report_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_report_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TamReportType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamReportingUnit
convert_sai_tam_reporting_unit_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_REPORTING_UNIT_SEC:
      return lemming::dataplane::sai::TAM_REPORTING_UNIT_SEC;

    case SAI_TAM_REPORTING_UNIT_MINUTE:
      return lemming::dataplane::sai::TAM_REPORTING_UNIT_MINUTE;

    case SAI_TAM_REPORTING_UNIT_HOUR:
      return lemming::dataplane::sai::TAM_REPORTING_UNIT_HOUR;

    case SAI_TAM_REPORTING_UNIT_DAY:
      return lemming::dataplane::sai::TAM_REPORTING_UNIT_DAY;

    default:
      return lemming::dataplane::sai::TAM_REPORTING_UNIT_UNSPECIFIED;
  }
}
sai_tam_reporting_unit_t convert_sai_tam_reporting_unit_t_to_sai(
    lemming::dataplane::sai::TamReportingUnit val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_REPORTING_UNIT_SEC:
      return SAI_TAM_REPORTING_UNIT_SEC;

    case lemming::dataplane::sai::TAM_REPORTING_UNIT_MINUTE:
      return SAI_TAM_REPORTING_UNIT_MINUTE;

    case lemming::dataplane::sai::TAM_REPORTING_UNIT_HOUR:
      return SAI_TAM_REPORTING_UNIT_HOUR;

    case lemming::dataplane::sai::TAM_REPORTING_UNIT_DAY:
      return SAI_TAM_REPORTING_UNIT_DAY;

    default:
      return SAI_TAM_REPORTING_UNIT_SEC;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_reporting_unit_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_reporting_unit_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_reporting_unit_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_reporting_unit_t_to_sai(
        static_cast<lemming::dataplane::sai::TamReportingUnit>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamTelMathFuncType
convert_sai_tam_tel_math_func_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_TEL_MATH_FUNC_TYPE_NONE:
      return lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_NONE;

    case SAI_TAM_TEL_MATH_FUNC_TYPE_GEO_MEAN:
      return lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_GEO_MEAN;

    case SAI_TAM_TEL_MATH_FUNC_TYPE_ALGEBRAIC_MEAN:
      return lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_ALGEBRAIC_MEAN;

    case SAI_TAM_TEL_MATH_FUNC_TYPE_AVERAGE:
      return lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_AVERAGE;

    case SAI_TAM_TEL_MATH_FUNC_TYPE_MODE:
      return lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_MODE;

    case SAI_TAM_TEL_MATH_FUNC_TYPE_RATE:
      return lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_RATE;

    default:
      return lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_UNSPECIFIED;
  }
}
sai_tam_tel_math_func_type_t convert_sai_tam_tel_math_func_type_t_to_sai(
    lemming::dataplane::sai::TamTelMathFuncType val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_NONE:
      return SAI_TAM_TEL_MATH_FUNC_TYPE_NONE;

    case lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_GEO_MEAN:
      return SAI_TAM_TEL_MATH_FUNC_TYPE_GEO_MEAN;

    case lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_ALGEBRAIC_MEAN:
      return SAI_TAM_TEL_MATH_FUNC_TYPE_ALGEBRAIC_MEAN;

    case lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_AVERAGE:
      return SAI_TAM_TEL_MATH_FUNC_TYPE_AVERAGE;

    case lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_MODE:
      return SAI_TAM_TEL_MATH_FUNC_TYPE_MODE;

    case lemming::dataplane::sai::TAM_TEL_MATH_FUNC_TYPE_RATE:
      return SAI_TAM_TEL_MATH_FUNC_TYPE_RATE;

    default:
      return SAI_TAM_TEL_MATH_FUNC_TYPE_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_tel_math_func_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_tel_math_func_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_tel_math_func_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_tel_math_func_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TamTelMathFuncType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamTelTypeAttr
convert_sai_tam_tel_type_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_TEL_TYPE_ATTR_TAM_TELEMETRY_TYPE:
      return lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_TAM_TELEMETRY_TYPE;

    case SAI_TAM_TEL_TYPE_ATTR_INT_SWITCH_IDENTIFIER:
      return lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_INT_SWITCH_IDENTIFIER;

    case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS:
      return lemming::dataplane::sai::
          TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS;

    case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_INGRESS:
      return lemming::dataplane::sai::
          TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_INGRESS;

    case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_EGRESS:
      return lemming::dataplane::sai::
          TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_EGRESS;

    case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_VIRTUAL_QUEUE_STATS:
      return lemming::dataplane::sai::
          TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_VIRTUAL_QUEUE_STATS;

    case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_OUTPUT_QUEUE_STATS:
      return lemming::dataplane::sai::
          TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_OUTPUT_QUEUE_STATS;

    case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_MMU_STATS:
      return lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_MMU_STATS;

    case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FABRIC_STATS:
      return lemming::dataplane::sai::
          TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FABRIC_STATS;

    case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FILTER_STATS:
      return lemming::dataplane::sai::
          TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FILTER_STATS;

    case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_RESOURCE_UTILIZATION_STATS:
      return lemming::dataplane::sai::
          TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_RESOURCE_UTILIZATION_STATS;

    case SAI_TAM_TEL_TYPE_ATTR_FABRIC_Q:
      return lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_FABRIC_Q;

    case SAI_TAM_TEL_TYPE_ATTR_NE_ENABLE:
      return lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_NE_ENABLE;

    case SAI_TAM_TEL_TYPE_ATTR_DSCP_VALUE:
      return lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_DSCP_VALUE;

    case SAI_TAM_TEL_TYPE_ATTR_MATH_FUNC:
      return lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_MATH_FUNC;

    case SAI_TAM_TEL_TYPE_ATTR_REPORT_ID:
      return lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_REPORT_ID;

    default:
      return lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_UNSPECIFIED;
  }
}
sai_tam_tel_type_attr_t convert_sai_tam_tel_type_attr_t_to_sai(
    lemming::dataplane::sai::TamTelTypeAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_TAM_TELEMETRY_TYPE:
      return SAI_TAM_TEL_TYPE_ATTR_TAM_TELEMETRY_TYPE;

    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_INT_SWITCH_IDENTIFIER:
      return SAI_TAM_TEL_TYPE_ATTR_INT_SWITCH_IDENTIFIER;

    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS:
      return SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS;

    case lemming::dataplane::sai::
        TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_INGRESS:
      return SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_INGRESS;

    case lemming::dataplane::sai::
        TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_EGRESS:
      return SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_EGRESS;

    case lemming::dataplane::sai::
        TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_VIRTUAL_QUEUE_STATS:
      return SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_VIRTUAL_QUEUE_STATS;

    case lemming::dataplane::sai::
        TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_OUTPUT_QUEUE_STATS:
      return SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_OUTPUT_QUEUE_STATS;

    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_MMU_STATS:
      return SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_MMU_STATS;

    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FABRIC_STATS:
      return SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FABRIC_STATS;

    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FILTER_STATS:
      return SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FILTER_STATS;

    case lemming::dataplane::sai::
        TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_RESOURCE_UTILIZATION_STATS:
      return SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_RESOURCE_UTILIZATION_STATS;

    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_FABRIC_Q:
      return SAI_TAM_TEL_TYPE_ATTR_FABRIC_Q;

    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_NE_ENABLE:
      return SAI_TAM_TEL_TYPE_ATTR_NE_ENABLE;

    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_DSCP_VALUE:
      return SAI_TAM_TEL_TYPE_ATTR_DSCP_VALUE;

    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_MATH_FUNC:
      return SAI_TAM_TEL_TYPE_ATTR_MATH_FUNC;

    case lemming::dataplane::sai::TAM_TEL_TYPE_ATTR_REPORT_ID:
      return SAI_TAM_TEL_TYPE_ATTR_REPORT_ID;

    default:
      return SAI_TAM_TEL_TYPE_ATTR_TAM_TELEMETRY_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_tel_type_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_tel_type_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_tel_type_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_tel_type_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamTelTypeAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamTelemetryAttr
convert_sai_tam_telemetry_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_TELEMETRY_ATTR_TAM_TYPE_LIST:
      return lemming::dataplane::sai::TAM_TELEMETRY_ATTR_TAM_TYPE_LIST;

    case SAI_TAM_TELEMETRY_ATTR_COLLECTOR_LIST:
      return lemming::dataplane::sai::TAM_TELEMETRY_ATTR_COLLECTOR_LIST;

    case SAI_TAM_TELEMETRY_ATTR_TAM_REPORTING_UNIT:
      return lemming::dataplane::sai::TAM_TELEMETRY_ATTR_TAM_REPORTING_UNIT;

    case SAI_TAM_TELEMETRY_ATTR_REPORTING_INTERVAL:
      return lemming::dataplane::sai::TAM_TELEMETRY_ATTR_REPORTING_INTERVAL;

    default:
      return lemming::dataplane::sai::TAM_TELEMETRY_ATTR_UNSPECIFIED;
  }
}
sai_tam_telemetry_attr_t convert_sai_tam_telemetry_attr_t_to_sai(
    lemming::dataplane::sai::TamTelemetryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_TELEMETRY_ATTR_TAM_TYPE_LIST:
      return SAI_TAM_TELEMETRY_ATTR_TAM_TYPE_LIST;

    case lemming::dataplane::sai::TAM_TELEMETRY_ATTR_COLLECTOR_LIST:
      return SAI_TAM_TELEMETRY_ATTR_COLLECTOR_LIST;

    case lemming::dataplane::sai::TAM_TELEMETRY_ATTR_TAM_REPORTING_UNIT:
      return SAI_TAM_TELEMETRY_ATTR_TAM_REPORTING_UNIT;

    case lemming::dataplane::sai::TAM_TELEMETRY_ATTR_REPORTING_INTERVAL:
      return SAI_TAM_TELEMETRY_ATTR_REPORTING_INTERVAL;

    default:
      return SAI_TAM_TELEMETRY_ATTR_TAM_TYPE_LIST;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_telemetry_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_telemetry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_telemetry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_telemetry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamTelemetryAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamTelemetryType
convert_sai_tam_telemetry_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_TELEMETRY_TYPE_NE:
      return lemming::dataplane::sai::TAM_TELEMETRY_TYPE_NE;

    case SAI_TAM_TELEMETRY_TYPE_SWITCH:
      return lemming::dataplane::sai::TAM_TELEMETRY_TYPE_SWITCH;

    case SAI_TAM_TELEMETRY_TYPE_FABRIC:
      return lemming::dataplane::sai::TAM_TELEMETRY_TYPE_FABRIC;

    case SAI_TAM_TELEMETRY_TYPE_FLOW:
      return lemming::dataplane::sai::TAM_TELEMETRY_TYPE_FLOW;

    case SAI_TAM_TELEMETRY_TYPE_INT:
      return lemming::dataplane::sai::TAM_TELEMETRY_TYPE_INT;

    default:
      return lemming::dataplane::sai::TAM_TELEMETRY_TYPE_UNSPECIFIED;
  }
}
sai_tam_telemetry_type_t convert_sai_tam_telemetry_type_t_to_sai(
    lemming::dataplane::sai::TamTelemetryType val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_TELEMETRY_TYPE_NE:
      return SAI_TAM_TELEMETRY_TYPE_NE;

    case lemming::dataplane::sai::TAM_TELEMETRY_TYPE_SWITCH:
      return SAI_TAM_TELEMETRY_TYPE_SWITCH;

    case lemming::dataplane::sai::TAM_TELEMETRY_TYPE_FABRIC:
      return SAI_TAM_TELEMETRY_TYPE_FABRIC;

    case lemming::dataplane::sai::TAM_TELEMETRY_TYPE_FLOW:
      return SAI_TAM_TELEMETRY_TYPE_FLOW;

    case lemming::dataplane::sai::TAM_TELEMETRY_TYPE_INT:
      return SAI_TAM_TELEMETRY_TYPE_INT;

    default:
      return SAI_TAM_TELEMETRY_TYPE_NE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_telemetry_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_telemetry_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_telemetry_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_telemetry_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TamTelemetryType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamTransportAttr
convert_sai_tam_transport_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_TRANSPORT_ATTR_TRANSPORT_TYPE:
      return lemming::dataplane::sai::TAM_TRANSPORT_ATTR_TRANSPORT_TYPE;

    case SAI_TAM_TRANSPORT_ATTR_SRC_PORT:
      return lemming::dataplane::sai::TAM_TRANSPORT_ATTR_SRC_PORT;

    case SAI_TAM_TRANSPORT_ATTR_DST_PORT:
      return lemming::dataplane::sai::TAM_TRANSPORT_ATTR_DST_PORT;

    case SAI_TAM_TRANSPORT_ATTR_TRANSPORT_AUTH_TYPE:
      return lemming::dataplane::sai::TAM_TRANSPORT_ATTR_TRANSPORT_AUTH_TYPE;

    case SAI_TAM_TRANSPORT_ATTR_MTU:
      return lemming::dataplane::sai::TAM_TRANSPORT_ATTR_MTU;

    default:
      return lemming::dataplane::sai::TAM_TRANSPORT_ATTR_UNSPECIFIED;
  }
}
sai_tam_transport_attr_t convert_sai_tam_transport_attr_t_to_sai(
    lemming::dataplane::sai::TamTransportAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_TRANSPORT_ATTR_TRANSPORT_TYPE:
      return SAI_TAM_TRANSPORT_ATTR_TRANSPORT_TYPE;

    case lemming::dataplane::sai::TAM_TRANSPORT_ATTR_SRC_PORT:
      return SAI_TAM_TRANSPORT_ATTR_SRC_PORT;

    case lemming::dataplane::sai::TAM_TRANSPORT_ATTR_DST_PORT:
      return SAI_TAM_TRANSPORT_ATTR_DST_PORT;

    case lemming::dataplane::sai::TAM_TRANSPORT_ATTR_TRANSPORT_AUTH_TYPE:
      return SAI_TAM_TRANSPORT_ATTR_TRANSPORT_AUTH_TYPE;

    case lemming::dataplane::sai::TAM_TRANSPORT_ATTR_MTU:
      return SAI_TAM_TRANSPORT_ATTR_MTU;

    default:
      return SAI_TAM_TRANSPORT_ATTR_TRANSPORT_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_transport_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_transport_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_transport_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_transport_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TamTransportAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamTransportAuthType
convert_sai_tam_transport_auth_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_TRANSPORT_AUTH_TYPE_NONE:
      return lemming::dataplane::sai::TAM_TRANSPORT_AUTH_TYPE_NONE;

    case SAI_TAM_TRANSPORT_AUTH_TYPE_SSL:
      return lemming::dataplane::sai::TAM_TRANSPORT_AUTH_TYPE_SSL;

    case SAI_TAM_TRANSPORT_AUTH_TYPE_TLS:
      return lemming::dataplane::sai::TAM_TRANSPORT_AUTH_TYPE_TLS;

    default:
      return lemming::dataplane::sai::TAM_TRANSPORT_AUTH_TYPE_UNSPECIFIED;
  }
}
sai_tam_transport_auth_type_t convert_sai_tam_transport_auth_type_t_to_sai(
    lemming::dataplane::sai::TamTransportAuthType val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_TRANSPORT_AUTH_TYPE_NONE:
      return SAI_TAM_TRANSPORT_AUTH_TYPE_NONE;

    case lemming::dataplane::sai::TAM_TRANSPORT_AUTH_TYPE_SSL:
      return SAI_TAM_TRANSPORT_AUTH_TYPE_SSL;

    case lemming::dataplane::sai::TAM_TRANSPORT_AUTH_TYPE_TLS:
      return SAI_TAM_TRANSPORT_AUTH_TYPE_TLS;

    default:
      return SAI_TAM_TRANSPORT_AUTH_TYPE_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_transport_auth_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_tam_transport_auth_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_transport_auth_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_transport_auth_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TamTransportAuthType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TamTransportType
convert_sai_tam_transport_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TAM_TRANSPORT_TYPE_NONE:
      return lemming::dataplane::sai::TAM_TRANSPORT_TYPE_NONE;

    case SAI_TAM_TRANSPORT_TYPE_TCP:
      return lemming::dataplane::sai::TAM_TRANSPORT_TYPE_TCP;

    case SAI_TAM_TRANSPORT_TYPE_UDP:
      return lemming::dataplane::sai::TAM_TRANSPORT_TYPE_UDP;

    case SAI_TAM_TRANSPORT_TYPE_GRPC:
      return lemming::dataplane::sai::TAM_TRANSPORT_TYPE_GRPC;

    case SAI_TAM_TRANSPORT_TYPE_MIRROR:
      return lemming::dataplane::sai::TAM_TRANSPORT_TYPE_MIRROR;

    default:
      return lemming::dataplane::sai::TAM_TRANSPORT_TYPE_UNSPECIFIED;
  }
}
sai_tam_transport_type_t convert_sai_tam_transport_type_t_to_sai(
    lemming::dataplane::sai::TamTransportType val) {
  switch (val) {
    case lemming::dataplane::sai::TAM_TRANSPORT_TYPE_NONE:
      return SAI_TAM_TRANSPORT_TYPE_NONE;

    case lemming::dataplane::sai::TAM_TRANSPORT_TYPE_TCP:
      return SAI_TAM_TRANSPORT_TYPE_TCP;

    case lemming::dataplane::sai::TAM_TRANSPORT_TYPE_UDP:
      return SAI_TAM_TRANSPORT_TYPE_UDP;

    case lemming::dataplane::sai::TAM_TRANSPORT_TYPE_GRPC:
      return SAI_TAM_TRANSPORT_TYPE_GRPC;

    case lemming::dataplane::sai::TAM_TRANSPORT_TYPE_MIRROR:
      return SAI_TAM_TRANSPORT_TYPE_MIRROR;

    default:
      return SAI_TAM_TRANSPORT_TYPE_NONE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tam_transport_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tam_transport_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tam_transport_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tam_transport_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TamTransportType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TlvType convert_sai_tlv_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TLV_TYPE_INGRESS:
      return lemming::dataplane::sai::TLV_TYPE_INGRESS;

    case SAI_TLV_TYPE_EGRESS:
      return lemming::dataplane::sai::TLV_TYPE_EGRESS;

    case SAI_TLV_TYPE_OPAQUE:
      return lemming::dataplane::sai::TLV_TYPE_OPAQUE;

    case SAI_TLV_TYPE_HMAC:
      return lemming::dataplane::sai::TLV_TYPE_HMAC;

    default:
      return lemming::dataplane::sai::TLV_TYPE_UNSPECIFIED;
  }
}
sai_tlv_type_t convert_sai_tlv_type_t_to_sai(
    lemming::dataplane::sai::TlvType val) {
  switch (val) {
    case lemming::dataplane::sai::TLV_TYPE_INGRESS:
      return SAI_TLV_TYPE_INGRESS;

    case lemming::dataplane::sai::TLV_TYPE_EGRESS:
      return SAI_TLV_TYPE_EGRESS;

    case lemming::dataplane::sai::TLV_TYPE_OPAQUE:
      return SAI_TLV_TYPE_OPAQUE;

    case lemming::dataplane::sai::TLV_TYPE_HMAC:
      return SAI_TLV_TYPE_HMAC;

    default:
      return SAI_TLV_TYPE_INGRESS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_tlv_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tlv_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tlv_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tlv_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TlvType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelAttr convert_sai_tunnel_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_ATTR_TYPE:
      return lemming::dataplane::sai::TUNNEL_ATTR_TYPE;

    case SAI_TUNNEL_ATTR_UNDERLAY_INTERFACE:
      return lemming::dataplane::sai::TUNNEL_ATTR_UNDERLAY_INTERFACE;

    case SAI_TUNNEL_ATTR_OVERLAY_INTERFACE:
      return lemming::dataplane::sai::TUNNEL_ATTR_OVERLAY_INTERFACE;

    case SAI_TUNNEL_ATTR_PEER_MODE:
      return lemming::dataplane::sai::TUNNEL_ATTR_PEER_MODE;

    case SAI_TUNNEL_ATTR_ENCAP_SRC_IP:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_SRC_IP;

    case SAI_TUNNEL_ATTR_ENCAP_DST_IP:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_DST_IP;

    case SAI_TUNNEL_ATTR_ENCAP_TTL_MODE:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_TTL_MODE;

    case SAI_TUNNEL_ATTR_ENCAP_TTL_VAL:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_TTL_VAL;

    case SAI_TUNNEL_ATTR_ENCAP_DSCP_MODE:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_DSCP_MODE;

    case SAI_TUNNEL_ATTR_ENCAP_DSCP_VAL:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_DSCP_VAL;

    case SAI_TUNNEL_ATTR_ENCAP_GRE_KEY_VALID:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_GRE_KEY_VALID;

    case SAI_TUNNEL_ATTR_ENCAP_GRE_KEY:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_GRE_KEY;

    case SAI_TUNNEL_ATTR_ENCAP_ECN_MODE:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_ECN_MODE;

    case SAI_TUNNEL_ATTR_ENCAP_MAPPERS:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_MAPPERS;

    case SAI_TUNNEL_ATTR_DECAP_ECN_MODE:
      return lemming::dataplane::sai::TUNNEL_ATTR_DECAP_ECN_MODE;

    case SAI_TUNNEL_ATTR_DECAP_MAPPERS:
      return lemming::dataplane::sai::TUNNEL_ATTR_DECAP_MAPPERS;

    case SAI_TUNNEL_ATTR_DECAP_TTL_MODE:
      return lemming::dataplane::sai::TUNNEL_ATTR_DECAP_TTL_MODE;

    case SAI_TUNNEL_ATTR_DECAP_DSCP_MODE:
      return lemming::dataplane::sai::TUNNEL_ATTR_DECAP_DSCP_MODE;

    case SAI_TUNNEL_ATTR_TERM_TABLE_ENTRY_LIST:
      return lemming::dataplane::sai::TUNNEL_ATTR_TERM_TABLE_ENTRY_LIST;

    case SAI_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
      return lemming::dataplane::sai::TUNNEL_ATTR_LOOPBACK_PACKET_ACTION;

    case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MODE:
      return lemming::dataplane::sai::TUNNEL_ATTR_VXLAN_UDP_SPORT_MODE;

    case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT:
      return lemming::dataplane::sai::TUNNEL_ATTR_VXLAN_UDP_SPORT;

    case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
      return lemming::dataplane::sai::TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK;

    case SAI_TUNNEL_ATTR_SA_INDEX:
      return lemming::dataplane::sai::TUNNEL_ATTR_SA_INDEX;

    case SAI_TUNNEL_ATTR_IPSEC_SA_PORT_LIST:
      return lemming::dataplane::sai::TUNNEL_ATTR_IPSEC_SA_PORT_LIST;

    case SAI_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
      return lemming::dataplane::sai::
          TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP;

    case SAI_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
      return lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP;

    case SAI_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
      return lemming::dataplane::sai::TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP;

    case SAI_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
      return lemming::dataplane::sai::
          TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP;

    case SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_SECURITY:
      return lemming::dataplane::sai::TUNNEL_ATTR_VXLAN_UDP_SPORT_SECURITY;

    default:
      return lemming::dataplane::sai::TUNNEL_ATTR_UNSPECIFIED;
  }
}
sai_tunnel_attr_t convert_sai_tunnel_attr_t_to_sai(
    lemming::dataplane::sai::TunnelAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_ATTR_TYPE:
      return SAI_TUNNEL_ATTR_TYPE;

    case lemming::dataplane::sai::TUNNEL_ATTR_UNDERLAY_INTERFACE:
      return SAI_TUNNEL_ATTR_UNDERLAY_INTERFACE;

    case lemming::dataplane::sai::TUNNEL_ATTR_OVERLAY_INTERFACE:
      return SAI_TUNNEL_ATTR_OVERLAY_INTERFACE;

    case lemming::dataplane::sai::TUNNEL_ATTR_PEER_MODE:
      return SAI_TUNNEL_ATTR_PEER_MODE;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_SRC_IP:
      return SAI_TUNNEL_ATTR_ENCAP_SRC_IP;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_DST_IP:
      return SAI_TUNNEL_ATTR_ENCAP_DST_IP;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_TTL_MODE:
      return SAI_TUNNEL_ATTR_ENCAP_TTL_MODE;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_TTL_VAL:
      return SAI_TUNNEL_ATTR_ENCAP_TTL_VAL;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_DSCP_MODE:
      return SAI_TUNNEL_ATTR_ENCAP_DSCP_MODE;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_DSCP_VAL:
      return SAI_TUNNEL_ATTR_ENCAP_DSCP_VAL;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_GRE_KEY_VALID:
      return SAI_TUNNEL_ATTR_ENCAP_GRE_KEY_VALID;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_GRE_KEY:
      return SAI_TUNNEL_ATTR_ENCAP_GRE_KEY;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_ECN_MODE:
      return SAI_TUNNEL_ATTR_ENCAP_ECN_MODE;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_MAPPERS:
      return SAI_TUNNEL_ATTR_ENCAP_MAPPERS;

    case lemming::dataplane::sai::TUNNEL_ATTR_DECAP_ECN_MODE:
      return SAI_TUNNEL_ATTR_DECAP_ECN_MODE;

    case lemming::dataplane::sai::TUNNEL_ATTR_DECAP_MAPPERS:
      return SAI_TUNNEL_ATTR_DECAP_MAPPERS;

    case lemming::dataplane::sai::TUNNEL_ATTR_DECAP_TTL_MODE:
      return SAI_TUNNEL_ATTR_DECAP_TTL_MODE;

    case lemming::dataplane::sai::TUNNEL_ATTR_DECAP_DSCP_MODE:
      return SAI_TUNNEL_ATTR_DECAP_DSCP_MODE;

    case lemming::dataplane::sai::TUNNEL_ATTR_TERM_TABLE_ENTRY_LIST:
      return SAI_TUNNEL_ATTR_TERM_TABLE_ENTRY_LIST;

    case lemming::dataplane::sai::TUNNEL_ATTR_LOOPBACK_PACKET_ACTION:
      return SAI_TUNNEL_ATTR_LOOPBACK_PACKET_ACTION;

    case lemming::dataplane::sai::TUNNEL_ATTR_VXLAN_UDP_SPORT_MODE:
      return SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MODE;

    case lemming::dataplane::sai::TUNNEL_ATTR_VXLAN_UDP_SPORT:
      return SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT;

    case lemming::dataplane::sai::TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK:
      return SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_MASK;

    case lemming::dataplane::sai::TUNNEL_ATTR_SA_INDEX:
      return SAI_TUNNEL_ATTR_SA_INDEX;

    case lemming::dataplane::sai::TUNNEL_ATTR_IPSEC_SA_PORT_LIST:
      return SAI_TUNNEL_ATTR_IPSEC_SA_PORT_LIST;

    case lemming::dataplane::sai::
        TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP:
      return SAI_TUNNEL_ATTR_ENCAP_QOS_TC_AND_COLOR_TO_DSCP_MAP;

    case lemming::dataplane::sai::TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP:
      return SAI_TUNNEL_ATTR_ENCAP_QOS_TC_TO_QUEUE_MAP;

    case lemming::dataplane::sai::TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP:
      return SAI_TUNNEL_ATTR_DECAP_QOS_DSCP_TO_TC_MAP;

    case lemming::dataplane::sai::
        TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP:
      return SAI_TUNNEL_ATTR_DECAP_QOS_TC_TO_PRIORITY_GROUP_MAP;

    case lemming::dataplane::sai::TUNNEL_ATTR_VXLAN_UDP_SPORT_SECURITY:
      return SAI_TUNNEL_ATTR_VXLAN_UDP_SPORT_SECURITY;

    default:
      return SAI_TUNNEL_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_tunnel_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelDecapEcnMode
convert_sai_tunnel_decap_ecn_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_DECAP_ECN_MODE_STANDARD:
      return lemming::dataplane::sai::TUNNEL_DECAP_ECN_MODE_STANDARD;

    case SAI_TUNNEL_DECAP_ECN_MODE_COPY_FROM_OUTER:
      return lemming::dataplane::sai::TUNNEL_DECAP_ECN_MODE_COPY_FROM_OUTER;

    case SAI_TUNNEL_DECAP_ECN_MODE_USER_DEFINED:
      return lemming::dataplane::sai::TUNNEL_DECAP_ECN_MODE_USER_DEFINED;

    default:
      return lemming::dataplane::sai::TUNNEL_DECAP_ECN_MODE_UNSPECIFIED;
  }
}
sai_tunnel_decap_ecn_mode_t convert_sai_tunnel_decap_ecn_mode_t_to_sai(
    lemming::dataplane::sai::TunnelDecapEcnMode val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_DECAP_ECN_MODE_STANDARD:
      return SAI_TUNNEL_DECAP_ECN_MODE_STANDARD;

    case lemming::dataplane::sai::TUNNEL_DECAP_ECN_MODE_COPY_FROM_OUTER:
      return SAI_TUNNEL_DECAP_ECN_MODE_COPY_FROM_OUTER;

    case lemming::dataplane::sai::TUNNEL_DECAP_ECN_MODE_USER_DEFINED:
      return SAI_TUNNEL_DECAP_ECN_MODE_USER_DEFINED;

    default:
      return SAI_TUNNEL_DECAP_ECN_MODE_STANDARD;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_decap_ecn_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_decap_ecn_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_decap_ecn_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_decap_ecn_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelDecapEcnMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelDscpMode convert_sai_tunnel_dscp_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_DSCP_MODE_UNIFORM_MODEL:
      return lemming::dataplane::sai::TUNNEL_DSCP_MODE_UNIFORM_MODEL;

    case SAI_TUNNEL_DSCP_MODE_PIPE_MODEL:
      return lemming::dataplane::sai::TUNNEL_DSCP_MODE_PIPE_MODEL;

    default:
      return lemming::dataplane::sai::TUNNEL_DSCP_MODE_UNSPECIFIED;
  }
}
sai_tunnel_dscp_mode_t convert_sai_tunnel_dscp_mode_t_to_sai(
    lemming::dataplane::sai::TunnelDscpMode val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_DSCP_MODE_UNIFORM_MODEL:
      return SAI_TUNNEL_DSCP_MODE_UNIFORM_MODEL;

    case lemming::dataplane::sai::TUNNEL_DSCP_MODE_PIPE_MODEL:
      return SAI_TUNNEL_DSCP_MODE_PIPE_MODEL;

    default:
      return SAI_TUNNEL_DSCP_MODE_UNIFORM_MODEL;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_dscp_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_dscp_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_dscp_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_dscp_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelDscpMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelEncapEcnMode
convert_sai_tunnel_encap_ecn_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_ENCAP_ECN_MODE_STANDARD:
      return lemming::dataplane::sai::TUNNEL_ENCAP_ECN_MODE_STANDARD;

    case SAI_TUNNEL_ENCAP_ECN_MODE_USER_DEFINED:
      return lemming::dataplane::sai::TUNNEL_ENCAP_ECN_MODE_USER_DEFINED;

    default:
      return lemming::dataplane::sai::TUNNEL_ENCAP_ECN_MODE_UNSPECIFIED;
  }
}
sai_tunnel_encap_ecn_mode_t convert_sai_tunnel_encap_ecn_mode_t_to_sai(
    lemming::dataplane::sai::TunnelEncapEcnMode val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_ENCAP_ECN_MODE_STANDARD:
      return SAI_TUNNEL_ENCAP_ECN_MODE_STANDARD;

    case lemming::dataplane::sai::TUNNEL_ENCAP_ECN_MODE_USER_DEFINED:
      return SAI_TUNNEL_ENCAP_ECN_MODE_USER_DEFINED;

    default:
      return SAI_TUNNEL_ENCAP_ECN_MODE_STANDARD;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_encap_ecn_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_encap_ecn_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_encap_ecn_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_encap_ecn_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelEncapEcnMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelMapAttr convert_sai_tunnel_map_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_MAP_ATTR_TYPE:
      return lemming::dataplane::sai::TUNNEL_MAP_ATTR_TYPE;

    case SAI_TUNNEL_MAP_ATTR_ENTRY_LIST:
      return lemming::dataplane::sai::TUNNEL_MAP_ATTR_ENTRY_LIST;

    default:
      return lemming::dataplane::sai::TUNNEL_MAP_ATTR_UNSPECIFIED;
  }
}
sai_tunnel_map_attr_t convert_sai_tunnel_map_attr_t_to_sai(
    lemming::dataplane::sai::TunnelMapAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_MAP_ATTR_TYPE:
      return SAI_TUNNEL_MAP_ATTR_TYPE;

    case lemming::dataplane::sai::TUNNEL_MAP_ATTR_ENTRY_LIST:
      return SAI_TUNNEL_MAP_ATTR_ENTRY_LIST;

    default:
      return SAI_TUNNEL_MAP_ATTR_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_map_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_map_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_map_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_map_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelMapAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelMapEntryAttr
convert_sai_tunnel_map_entry_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP_TYPE:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP_TYPE;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_OECN_KEY:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_OECN_KEY;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_OECN_VALUE:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_OECN_VALUE;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_UECN_KEY:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_UECN_KEY;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_UECN_VALUE:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_UECN_VALUE;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_KEY:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_KEY;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_VALUE:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_VALUE;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_VNI_ID_KEY:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VNI_ID_KEY;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_VNI_ID_VALUE:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VNI_ID_VALUE;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_KEY:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_KEY;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_VALUE:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_VALUE;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_KEY:
      return lemming::dataplane::sai::
          TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_KEY;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_VALUE:
      return lemming::dataplane::sai::
          TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_VALUE;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_VSID_ID_KEY:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VSID_ID_KEY;

    case SAI_TUNNEL_MAP_ENTRY_ATTR_VSID_ID_VALUE:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VSID_ID_VALUE;

    default:
      return lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_tunnel_map_entry_attr_t convert_sai_tunnel_map_entry_attr_t_to_sai(
    lemming::dataplane::sai::TunnelMapEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP_TYPE:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP_TYPE;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_OECN_KEY:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_OECN_KEY;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_OECN_VALUE:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_OECN_VALUE;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_UECN_KEY:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_UECN_KEY;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_UECN_VALUE:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_UECN_VALUE;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_KEY:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_KEY;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_VALUE:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_VLAN_ID_VALUE;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VNI_ID_KEY:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_VNI_ID_KEY;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VNI_ID_VALUE:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_VNI_ID_VALUE;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_KEY:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_KEY;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_VALUE:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_BRIDGE_ID_VALUE;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_KEY:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_KEY;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_VALUE:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_VIRTUAL_ROUTER_ID_VALUE;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VSID_ID_KEY:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_VSID_ID_KEY;

    case lemming::dataplane::sai::TUNNEL_MAP_ENTRY_ATTR_VSID_ID_VALUE:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_VSID_ID_VALUE;

    default:
      return SAI_TUNNEL_MAP_ENTRY_ATTR_TUNNEL_MAP_TYPE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_map_entry_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_map_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_map_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_map_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelMapEntryAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelMapType convert_sai_tunnel_map_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_MAP_TYPE_OECN_TO_UECN:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_OECN_TO_UECN;

    case SAI_TUNNEL_MAP_TYPE_UECN_OECN_TO_OECN:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_UECN_OECN_TO_OECN;

    case SAI_TUNNEL_MAP_TYPE_VNI_TO_VLAN_ID:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_VNI_TO_VLAN_ID;

    case SAI_TUNNEL_MAP_TYPE_VLAN_ID_TO_VNI:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_VLAN_ID_TO_VNI;

    case SAI_TUNNEL_MAP_TYPE_VNI_TO_BRIDGE_IF:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_VNI_TO_BRIDGE_IF;

    case SAI_TUNNEL_MAP_TYPE_BRIDGE_IF_TO_VNI:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_BRIDGE_IF_TO_VNI;

    case SAI_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID;

    case SAI_TUNNEL_MAP_TYPE_VIRTUAL_ROUTER_ID_TO_VNI:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_VIRTUAL_ROUTER_ID_TO_VNI;

    case SAI_TUNNEL_MAP_TYPE_VSID_TO_VLAN_ID:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_VSID_TO_VLAN_ID;

    case SAI_TUNNEL_MAP_TYPE_VLAN_ID_TO_VSID:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_VLAN_ID_TO_VSID;

    case SAI_TUNNEL_MAP_TYPE_VSID_TO_BRIDGE_IF:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_VSID_TO_BRIDGE_IF;

    case SAI_TUNNEL_MAP_TYPE_BRIDGE_IF_TO_VSID:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_BRIDGE_IF_TO_VSID;

    case SAI_TUNNEL_MAP_TYPE_CUSTOM_RANGE_BASE:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_CUSTOM_RANGE_BASE;

    default:
      return lemming::dataplane::sai::TUNNEL_MAP_TYPE_UNSPECIFIED;
  }
}
sai_tunnel_map_type_t convert_sai_tunnel_map_type_t_to_sai(
    lemming::dataplane::sai::TunnelMapType val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_OECN_TO_UECN:
      return SAI_TUNNEL_MAP_TYPE_OECN_TO_UECN;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_UECN_OECN_TO_OECN:
      return SAI_TUNNEL_MAP_TYPE_UECN_OECN_TO_OECN;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_VNI_TO_VLAN_ID:
      return SAI_TUNNEL_MAP_TYPE_VNI_TO_VLAN_ID;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_VLAN_ID_TO_VNI:
      return SAI_TUNNEL_MAP_TYPE_VLAN_ID_TO_VNI;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_VNI_TO_BRIDGE_IF:
      return SAI_TUNNEL_MAP_TYPE_VNI_TO_BRIDGE_IF;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_BRIDGE_IF_TO_VNI:
      return SAI_TUNNEL_MAP_TYPE_BRIDGE_IF_TO_VNI;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID:
      return SAI_TUNNEL_MAP_TYPE_VNI_TO_VIRTUAL_ROUTER_ID;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_VIRTUAL_ROUTER_ID_TO_VNI:
      return SAI_TUNNEL_MAP_TYPE_VIRTUAL_ROUTER_ID_TO_VNI;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_VSID_TO_VLAN_ID:
      return SAI_TUNNEL_MAP_TYPE_VSID_TO_VLAN_ID;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_VLAN_ID_TO_VSID:
      return SAI_TUNNEL_MAP_TYPE_VLAN_ID_TO_VSID;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_VSID_TO_BRIDGE_IF:
      return SAI_TUNNEL_MAP_TYPE_VSID_TO_BRIDGE_IF;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_BRIDGE_IF_TO_VSID:
      return SAI_TUNNEL_MAP_TYPE_BRIDGE_IF_TO_VSID;

    case lemming::dataplane::sai::TUNNEL_MAP_TYPE_CUSTOM_RANGE_BASE:
      return SAI_TUNNEL_MAP_TYPE_CUSTOM_RANGE_BASE;

    default:
      return SAI_TUNNEL_MAP_TYPE_OECN_TO_UECN;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_map_type_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_map_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_map_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_map_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelMapType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelPeerMode convert_sai_tunnel_peer_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_PEER_MODE_P2P:
      return lemming::dataplane::sai::TUNNEL_PEER_MODE_P2P;

    case SAI_TUNNEL_PEER_MODE_P2MP:
      return lemming::dataplane::sai::TUNNEL_PEER_MODE_P2MP;

    default:
      return lemming::dataplane::sai::TUNNEL_PEER_MODE_UNSPECIFIED;
  }
}
sai_tunnel_peer_mode_t convert_sai_tunnel_peer_mode_t_to_sai(
    lemming::dataplane::sai::TunnelPeerMode val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_PEER_MODE_P2P:
      return SAI_TUNNEL_PEER_MODE_P2P;

    case lemming::dataplane::sai::TUNNEL_PEER_MODE_P2MP:
      return SAI_TUNNEL_PEER_MODE_P2MP;

    default:
      return SAI_TUNNEL_PEER_MODE_P2P;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_peer_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_peer_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_peer_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_peer_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelPeerMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelStat convert_sai_tunnel_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_STAT_IN_OCTETS:
      return lemming::dataplane::sai::TUNNEL_STAT_IN_OCTETS;

    case SAI_TUNNEL_STAT_IN_PACKETS:
      return lemming::dataplane::sai::TUNNEL_STAT_IN_PACKETS;

    case SAI_TUNNEL_STAT_OUT_OCTETS:
      return lemming::dataplane::sai::TUNNEL_STAT_OUT_OCTETS;

    case SAI_TUNNEL_STAT_OUT_PACKETS:
      return lemming::dataplane::sai::TUNNEL_STAT_OUT_PACKETS;

    default:
      return lemming::dataplane::sai::TUNNEL_STAT_UNSPECIFIED;
  }
}
sai_tunnel_stat_t convert_sai_tunnel_stat_t_to_sai(
    lemming::dataplane::sai::TunnelStat val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_STAT_IN_OCTETS:
      return SAI_TUNNEL_STAT_IN_OCTETS;

    case lemming::dataplane::sai::TUNNEL_STAT_IN_PACKETS:
      return SAI_TUNNEL_STAT_IN_PACKETS;

    case lemming::dataplane::sai::TUNNEL_STAT_OUT_OCTETS:
      return SAI_TUNNEL_STAT_OUT_OCTETS;

    case lemming::dataplane::sai::TUNNEL_STAT_OUT_PACKETS:
      return SAI_TUNNEL_STAT_OUT_PACKETS;

    default:
      return SAI_TUNNEL_STAT_IN_OCTETS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_tunnel_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelTermTableEntryAttr
convert_sai_tunnel_term_table_entry_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_VR_ID:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_VR_ID;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_TYPE:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_TYPE;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP_MASK:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP_MASK;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP_MASK:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP_MASK;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_TUNNEL_TYPE:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_TUNNEL_TYPE;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_ACTION_TUNNEL_ID:
      return lemming::dataplane::sai::
          TUNNEL_TERM_TABLE_ENTRY_ATTR_ACTION_TUNNEL_ID;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_IP_ADDR_FAMILY:
      return lemming::dataplane::sai::
          TUNNEL_TERM_TABLE_ENTRY_ATTR_IP_ADDR_FAMILY;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_IPSEC_VERIFIED:
      return lemming::dataplane::sai::
          TUNNEL_TERM_TABLE_ENTRY_ATTR_IPSEC_VERIFIED;

    default:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_UNSPECIFIED;
  }
}
sai_tunnel_term_table_entry_attr_t
convert_sai_tunnel_term_table_entry_attr_t_to_sai(
    lemming::dataplane::sai::TunnelTermTableEntryAttr val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_VR_ID:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_VR_ID;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_TYPE:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_TYPE;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP_MASK:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_DST_IP_MASK;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP_MASK:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_SRC_IP_MASK;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_TUNNEL_TYPE:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_TUNNEL_TYPE;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_ACTION_TUNNEL_ID:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_ACTION_TUNNEL_ID;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_IP_ADDR_FAMILY:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_IP_ADDR_FAMILY;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_ATTR_IPSEC_VERIFIED:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_IPSEC_VERIFIED;

    default:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_ATTR_VR_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_term_table_entry_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_tunnel_term_table_entry_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_term_table_entry_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_term_table_entry_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelTermTableEntryAttr>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelTermTableEntryType
convert_sai_tunnel_term_table_entry_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2P:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2P;

    case SAI_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2MP:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2MP;

    default:
      return lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_TYPE_UNSPECIFIED;
  }
}
sai_tunnel_term_table_entry_type_t
convert_sai_tunnel_term_table_entry_type_t_to_sai(
    lemming::dataplane::sai::TunnelTermTableEntryType val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2MP;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2P:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2P;

    case lemming::dataplane::sai::TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2MP:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_TYPE_MP2MP;

    default:
      return SAI_TUNNEL_TERM_TABLE_ENTRY_TYPE_P2P;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_term_table_entry_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_tunnel_term_table_entry_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_term_table_entry_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_term_table_entry_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelTermTableEntryType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelTtlMode convert_sai_tunnel_ttl_mode_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_TTL_MODE_UNIFORM_MODEL:
      return lemming::dataplane::sai::TUNNEL_TTL_MODE_UNIFORM_MODEL;

    case SAI_TUNNEL_TTL_MODE_PIPE_MODEL:
      return lemming::dataplane::sai::TUNNEL_TTL_MODE_PIPE_MODEL;

    default:
      return lemming::dataplane::sai::TUNNEL_TTL_MODE_UNSPECIFIED;
  }
}
sai_tunnel_ttl_mode_t convert_sai_tunnel_ttl_mode_t_to_sai(
    lemming::dataplane::sai::TunnelTtlMode val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_TTL_MODE_UNIFORM_MODEL:
      return SAI_TUNNEL_TTL_MODE_UNIFORM_MODEL;

    case lemming::dataplane::sai::TUNNEL_TTL_MODE_PIPE_MODEL:
      return SAI_TUNNEL_TTL_MODE_PIPE_MODEL;

    default:
      return SAI_TUNNEL_TTL_MODE_UNIFORM_MODEL;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_ttl_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_ttl_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_ttl_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_ttl_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelTtlMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelType convert_sai_tunnel_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_TYPE_IPINIP:
      return lemming::dataplane::sai::TUNNEL_TYPE_IPINIP;

    case SAI_TUNNEL_TYPE_IPINIP_GRE:
      return lemming::dataplane::sai::TUNNEL_TYPE_IPINIP_GRE;

    case SAI_TUNNEL_TYPE_VXLAN:
      return lemming::dataplane::sai::TUNNEL_TYPE_VXLAN;

    case SAI_TUNNEL_TYPE_MPLS:
      return lemming::dataplane::sai::TUNNEL_TYPE_MPLS;

    case SAI_TUNNEL_TYPE_SRV6:
      return lemming::dataplane::sai::TUNNEL_TYPE_SRV6;

    case SAI_TUNNEL_TYPE_NVGRE:
      return lemming::dataplane::sai::TUNNEL_TYPE_NVGRE;

    case SAI_TUNNEL_TYPE_IPINIP_ESP:
      return lemming::dataplane::sai::TUNNEL_TYPE_IPINIP_ESP;

    case SAI_TUNNEL_TYPE_IPINIP_UDP_ESP:
      return lemming::dataplane::sai::TUNNEL_TYPE_IPINIP_UDP_ESP;

    case SAI_TUNNEL_TYPE_VXLAN_UDP_ESP:
      return lemming::dataplane::sai::TUNNEL_TYPE_VXLAN_UDP_ESP;

    default:
      return lemming::dataplane::sai::TUNNEL_TYPE_UNSPECIFIED;
  }
}
sai_tunnel_type_t convert_sai_tunnel_type_t_to_sai(
    lemming::dataplane::sai::TunnelType val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_TYPE_IPINIP:
      return SAI_TUNNEL_TYPE_IPINIP;

    case lemming::dataplane::sai::TUNNEL_TYPE_IPINIP_GRE:
      return SAI_TUNNEL_TYPE_IPINIP_GRE;

    case lemming::dataplane::sai::TUNNEL_TYPE_VXLAN:
      return SAI_TUNNEL_TYPE_VXLAN;

    case lemming::dataplane::sai::TUNNEL_TYPE_MPLS:
      return SAI_TUNNEL_TYPE_MPLS;

    case lemming::dataplane::sai::TUNNEL_TYPE_SRV6:
      return SAI_TUNNEL_TYPE_SRV6;

    case lemming::dataplane::sai::TUNNEL_TYPE_NVGRE:
      return SAI_TUNNEL_TYPE_NVGRE;

    case lemming::dataplane::sai::TUNNEL_TYPE_IPINIP_ESP:
      return SAI_TUNNEL_TYPE_IPINIP_ESP;

    case lemming::dataplane::sai::TUNNEL_TYPE_IPINIP_UDP_ESP:
      return SAI_TUNNEL_TYPE_IPINIP_UDP_ESP;

    case lemming::dataplane::sai::TUNNEL_TYPE_VXLAN_UDP_ESP:
      return SAI_TUNNEL_TYPE_VXLAN_UDP_ESP;

    default:
      return SAI_TUNNEL_TYPE_IPINIP;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_tunnel_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_tunnel_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_type_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::TunnelVxlanUdpSportMode
convert_sai_tunnel_vxlan_udp_sport_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_TUNNEL_VXLAN_UDP_SPORT_MODE_USER_DEFINED:
      return lemming::dataplane::sai::TUNNEL_VXLAN_UDP_SPORT_MODE_USER_DEFINED;

    case SAI_TUNNEL_VXLAN_UDP_SPORT_MODE_EPHEMERAL:
      return lemming::dataplane::sai::TUNNEL_VXLAN_UDP_SPORT_MODE_EPHEMERAL;

    default:
      return lemming::dataplane::sai::TUNNEL_VXLAN_UDP_SPORT_MODE_UNSPECIFIED;
  }
}
sai_tunnel_vxlan_udp_sport_mode_t
convert_sai_tunnel_vxlan_udp_sport_mode_t_to_sai(
    lemming::dataplane::sai::TunnelVxlanUdpSportMode val) {
  switch (val) {
    case lemming::dataplane::sai::TUNNEL_VXLAN_UDP_SPORT_MODE_USER_DEFINED:
      return SAI_TUNNEL_VXLAN_UDP_SPORT_MODE_USER_DEFINED;

    case lemming::dataplane::sai::TUNNEL_VXLAN_UDP_SPORT_MODE_EPHEMERAL:
      return SAI_TUNNEL_VXLAN_UDP_SPORT_MODE_EPHEMERAL;

    default:
      return SAI_TUNNEL_VXLAN_UDP_SPORT_MODE_USER_DEFINED;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_tunnel_vxlan_udp_sport_mode_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_tunnel_vxlan_udp_sport_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_tunnel_vxlan_udp_sport_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_tunnel_vxlan_udp_sport_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::TunnelVxlanUdpSportMode>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::UdfAttr convert_sai_udf_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_UDF_ATTR_MATCH_ID:
      return lemming::dataplane::sai::UDF_ATTR_MATCH_ID;

    case SAI_UDF_ATTR_GROUP_ID:
      return lemming::dataplane::sai::UDF_ATTR_GROUP_ID;

    case SAI_UDF_ATTR_BASE:
      return lemming::dataplane::sai::UDF_ATTR_BASE;

    case SAI_UDF_ATTR_OFFSET:
      return lemming::dataplane::sai::UDF_ATTR_OFFSET;

    case SAI_UDF_ATTR_HASH_MASK:
      return lemming::dataplane::sai::UDF_ATTR_HASH_MASK;

    default:
      return lemming::dataplane::sai::UDF_ATTR_UNSPECIFIED;
  }
}
sai_udf_attr_t convert_sai_udf_attr_t_to_sai(
    lemming::dataplane::sai::UdfAttr val) {
  switch (val) {
    case lemming::dataplane::sai::UDF_ATTR_MATCH_ID:
      return SAI_UDF_ATTR_MATCH_ID;

    case lemming::dataplane::sai::UDF_ATTR_GROUP_ID:
      return SAI_UDF_ATTR_GROUP_ID;

    case lemming::dataplane::sai::UDF_ATTR_BASE:
      return SAI_UDF_ATTR_BASE;

    case lemming::dataplane::sai::UDF_ATTR_OFFSET:
      return SAI_UDF_ATTR_OFFSET;

    case lemming::dataplane::sai::UDF_ATTR_HASH_MASK:
      return SAI_UDF_ATTR_HASH_MASK;

    default:
      return SAI_UDF_ATTR_MATCH_ID;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_udf_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_udf_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_udf_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_udf_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::UdfAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::UdfBase convert_sai_udf_base_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_UDF_BASE_L2:
      return lemming::dataplane::sai::UDF_BASE_L2;

    case SAI_UDF_BASE_L3:
      return lemming::dataplane::sai::UDF_BASE_L3;

    case SAI_UDF_BASE_L4:
      return lemming::dataplane::sai::UDF_BASE_L4;

    default:
      return lemming::dataplane::sai::UDF_BASE_UNSPECIFIED;
  }
}
sai_udf_base_t convert_sai_udf_base_t_to_sai(
    lemming::dataplane::sai::UdfBase val) {
  switch (val) {
    case lemming::dataplane::sai::UDF_BASE_L2:
      return SAI_UDF_BASE_L2;

    case lemming::dataplane::sai::UDF_BASE_L3:
      return SAI_UDF_BASE_L3;

    case lemming::dataplane::sai::UDF_BASE_L4:
      return SAI_UDF_BASE_L4;

    default:
      return SAI_UDF_BASE_L2;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_udf_base_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_udf_base_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_udf_base_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_udf_base_t_to_sai(
        static_cast<lemming::dataplane::sai::UdfBase>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::UdfGroupAttr convert_sai_udf_group_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_UDF_GROUP_ATTR_UDF_LIST:
      return lemming::dataplane::sai::UDF_GROUP_ATTR_UDF_LIST;

    case SAI_UDF_GROUP_ATTR_TYPE:
      return lemming::dataplane::sai::UDF_GROUP_ATTR_TYPE;

    case SAI_UDF_GROUP_ATTR_LENGTH:
      return lemming::dataplane::sai::UDF_GROUP_ATTR_LENGTH;

    default:
      return lemming::dataplane::sai::UDF_GROUP_ATTR_UNSPECIFIED;
  }
}
sai_udf_group_attr_t convert_sai_udf_group_attr_t_to_sai(
    lemming::dataplane::sai::UdfGroupAttr val) {
  switch (val) {
    case lemming::dataplane::sai::UDF_GROUP_ATTR_UDF_LIST:
      return SAI_UDF_GROUP_ATTR_UDF_LIST;

    case lemming::dataplane::sai::UDF_GROUP_ATTR_TYPE:
      return SAI_UDF_GROUP_ATTR_TYPE;

    case lemming::dataplane::sai::UDF_GROUP_ATTR_LENGTH:
      return SAI_UDF_GROUP_ATTR_LENGTH;

    default:
      return SAI_UDF_GROUP_ATTR_UDF_LIST;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_udf_group_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_udf_group_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_udf_group_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_udf_group_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::UdfGroupAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::UdfGroupType convert_sai_udf_group_type_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_UDF_GROUP_TYPE_START:
      return lemming::dataplane::sai::UDF_GROUP_TYPE_START;

    case SAI_UDF_GROUP_TYPE_HASH:
      return lemming::dataplane::sai::UDF_GROUP_TYPE_HASH;

    case SAI_UDF_GROUP_TYPE_END:
      return lemming::dataplane::sai::UDF_GROUP_TYPE_END;

    default:
      return lemming::dataplane::sai::UDF_GROUP_TYPE_UNSPECIFIED;
  }
}
sai_udf_group_type_t convert_sai_udf_group_type_t_to_sai(
    lemming::dataplane::sai::UdfGroupType val) {
  switch (val) {
    case lemming::dataplane::sai::UDF_GROUP_TYPE_START:
      return SAI_UDF_GROUP_TYPE_START;

    case lemming::dataplane::sai::UDF_GROUP_TYPE_HASH:
      return SAI_UDF_GROUP_TYPE_HASH;

    case lemming::dataplane::sai::UDF_GROUP_TYPE_END:
      return SAI_UDF_GROUP_TYPE_END;

    default:
      return SAI_UDF_GROUP_TYPE_START;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_udf_group_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_udf_group_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_udf_group_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_udf_group_type_t_to_sai(
        static_cast<lemming::dataplane::sai::UdfGroupType>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::UdfMatchAttr convert_sai_udf_match_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_UDF_MATCH_ATTR_L2_TYPE:
      return lemming::dataplane::sai::UDF_MATCH_ATTR_L2_TYPE;

    case SAI_UDF_MATCH_ATTR_L3_TYPE:
      return lemming::dataplane::sai::UDF_MATCH_ATTR_L3_TYPE;

    case SAI_UDF_MATCH_ATTR_GRE_TYPE:
      return lemming::dataplane::sai::UDF_MATCH_ATTR_GRE_TYPE;

    case SAI_UDF_MATCH_ATTR_PRIORITY:
      return lemming::dataplane::sai::UDF_MATCH_ATTR_PRIORITY;

    default:
      return lemming::dataplane::sai::UDF_MATCH_ATTR_UNSPECIFIED;
  }
}
sai_udf_match_attr_t convert_sai_udf_match_attr_t_to_sai(
    lemming::dataplane::sai::UdfMatchAttr val) {
  switch (val) {
    case lemming::dataplane::sai::UDF_MATCH_ATTR_L2_TYPE:
      return SAI_UDF_MATCH_ATTR_L2_TYPE;

    case lemming::dataplane::sai::UDF_MATCH_ATTR_L3_TYPE:
      return SAI_UDF_MATCH_ATTR_L3_TYPE;

    case lemming::dataplane::sai::UDF_MATCH_ATTR_GRE_TYPE:
      return SAI_UDF_MATCH_ATTR_GRE_TYPE;

    case lemming::dataplane::sai::UDF_MATCH_ATTR_PRIORITY:
      return SAI_UDF_MATCH_ATTR_PRIORITY;

    default:
      return SAI_UDF_MATCH_ATTR_L2_TYPE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_udf_match_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_udf_match_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_udf_match_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_udf_match_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::UdfMatchAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::VirtualRouterAttr
convert_sai_virtual_router_attr_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_VIRTUAL_ROUTER_ATTR_ADMIN_V4_STATE:
      return lemming::dataplane::sai::VIRTUAL_ROUTER_ATTR_ADMIN_V4_STATE;

    case SAI_VIRTUAL_ROUTER_ATTR_ADMIN_V6_STATE:
      return lemming::dataplane::sai::VIRTUAL_ROUTER_ATTR_ADMIN_V6_STATE;

    case SAI_VIRTUAL_ROUTER_ATTR_SRC_MAC_ADDRESS:
      return lemming::dataplane::sai::VIRTUAL_ROUTER_ATTR_SRC_MAC_ADDRESS;

    case SAI_VIRTUAL_ROUTER_ATTR_VIOLATION_TTL1_PACKET_ACTION:
      return lemming::dataplane::sai::
          VIRTUAL_ROUTER_ATTR_VIOLATION_TTL1_PACKET_ACTION;

    case SAI_VIRTUAL_ROUTER_ATTR_VIOLATION_IP_OPTIONS_PACKET_ACTION:
      return lemming::dataplane::sai::
          VIRTUAL_ROUTER_ATTR_VIOLATION_IP_OPTIONS_PACKET_ACTION;

    case SAI_VIRTUAL_ROUTER_ATTR_UNKNOWN_L3_MULTICAST_PACKET_ACTION:
      return lemming::dataplane::sai::
          VIRTUAL_ROUTER_ATTR_UNKNOWN_L3_MULTICAST_PACKET_ACTION;

    case SAI_VIRTUAL_ROUTER_ATTR_LABEL:
      return lemming::dataplane::sai::VIRTUAL_ROUTER_ATTR_LABEL;

    default:
      return lemming::dataplane::sai::VIRTUAL_ROUTER_ATTR_UNSPECIFIED;
  }
}
sai_virtual_router_attr_t convert_sai_virtual_router_attr_t_to_sai(
    lemming::dataplane::sai::VirtualRouterAttr val) {
  switch (val) {
    case lemming::dataplane::sai::VIRTUAL_ROUTER_ATTR_ADMIN_V4_STATE:
      return SAI_VIRTUAL_ROUTER_ATTR_ADMIN_V4_STATE;

    case lemming::dataplane::sai::VIRTUAL_ROUTER_ATTR_ADMIN_V6_STATE:
      return SAI_VIRTUAL_ROUTER_ATTR_ADMIN_V6_STATE;

    case lemming::dataplane::sai::VIRTUAL_ROUTER_ATTR_SRC_MAC_ADDRESS:
      return SAI_VIRTUAL_ROUTER_ATTR_SRC_MAC_ADDRESS;

    case lemming::dataplane::sai::
        VIRTUAL_ROUTER_ATTR_VIOLATION_TTL1_PACKET_ACTION:
      return SAI_VIRTUAL_ROUTER_ATTR_VIOLATION_TTL1_PACKET_ACTION;

    case lemming::dataplane::sai::
        VIRTUAL_ROUTER_ATTR_VIOLATION_IP_OPTIONS_PACKET_ACTION:
      return SAI_VIRTUAL_ROUTER_ATTR_VIOLATION_IP_OPTIONS_PACKET_ACTION;

    case lemming::dataplane::sai::
        VIRTUAL_ROUTER_ATTR_UNKNOWN_L3_MULTICAST_PACKET_ACTION:
      return SAI_VIRTUAL_ROUTER_ATTR_UNKNOWN_L3_MULTICAST_PACKET_ACTION;

    case lemming::dataplane::sai::VIRTUAL_ROUTER_ATTR_LABEL:
      return SAI_VIRTUAL_ROUTER_ATTR_LABEL;

    default:
      return SAI_VIRTUAL_ROUTER_ATTR_ADMIN_V4_STATE;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_virtual_router_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_virtual_router_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_virtual_router_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_virtual_router_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::VirtualRouterAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::VlanAttr convert_sai_vlan_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_VLAN_ATTR_VLAN_ID:
      return lemming::dataplane::sai::VLAN_ATTR_VLAN_ID;

    case SAI_VLAN_ATTR_MEMBER_LIST:
      return lemming::dataplane::sai::VLAN_ATTR_MEMBER_LIST;

    case SAI_VLAN_ATTR_MAX_LEARNED_ADDRESSES:
      return lemming::dataplane::sai::VLAN_ATTR_MAX_LEARNED_ADDRESSES;

    case SAI_VLAN_ATTR_STP_INSTANCE:
      return lemming::dataplane::sai::VLAN_ATTR_STP_INSTANCE;

    case SAI_VLAN_ATTR_LEARN_DISABLE:
      return lemming::dataplane::sai::VLAN_ATTR_LEARN_DISABLE;

    case SAI_VLAN_ATTR_IPV4_MCAST_LOOKUP_KEY_TYPE:
      return lemming::dataplane::sai::VLAN_ATTR_IPV4_MCAST_LOOKUP_KEY_TYPE;

    case SAI_VLAN_ATTR_IPV6_MCAST_LOOKUP_KEY_TYPE:
      return lemming::dataplane::sai::VLAN_ATTR_IPV6_MCAST_LOOKUP_KEY_TYPE;

    case SAI_VLAN_ATTR_UNKNOWN_NON_IP_MCAST_OUTPUT_GROUP_ID:
      return lemming::dataplane::sai::
          VLAN_ATTR_UNKNOWN_NON_IP_MCAST_OUTPUT_GROUP_ID;

    case SAI_VLAN_ATTR_UNKNOWN_IPV4_MCAST_OUTPUT_GROUP_ID:
      return lemming::dataplane::sai::
          VLAN_ATTR_UNKNOWN_IPV4_MCAST_OUTPUT_GROUP_ID;

    case SAI_VLAN_ATTR_UNKNOWN_IPV6_MCAST_OUTPUT_GROUP_ID:
      return lemming::dataplane::sai::
          VLAN_ATTR_UNKNOWN_IPV6_MCAST_OUTPUT_GROUP_ID;

    case SAI_VLAN_ATTR_UNKNOWN_LINKLOCAL_MCAST_OUTPUT_GROUP_ID:
      return lemming::dataplane::sai::
          VLAN_ATTR_UNKNOWN_LINKLOCAL_MCAST_OUTPUT_GROUP_ID;

    case SAI_VLAN_ATTR_INGRESS_ACL:
      return lemming::dataplane::sai::VLAN_ATTR_INGRESS_ACL;

    case SAI_VLAN_ATTR_EGRESS_ACL:
      return lemming::dataplane::sai::VLAN_ATTR_EGRESS_ACL;

    case SAI_VLAN_ATTR_META_DATA:
      return lemming::dataplane::sai::VLAN_ATTR_META_DATA;

    case SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
      return lemming::dataplane::sai::
          VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE;

    case SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
      return lemming::dataplane::sai::VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP;

    case SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
      return lemming::dataplane::sai::
          VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE;

    case SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
      return lemming::dataplane::sai::VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP;

    case SAI_VLAN_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
      return lemming::dataplane::sai::VLAN_ATTR_BROADCAST_FLOOD_CONTROL_TYPE;

    case SAI_VLAN_ATTR_BROADCAST_FLOOD_GROUP:
      return lemming::dataplane::sai::VLAN_ATTR_BROADCAST_FLOOD_GROUP;

    case SAI_VLAN_ATTR_CUSTOM_IGMP_SNOOPING_ENABLE:
      return lemming::dataplane::sai::VLAN_ATTR_CUSTOM_IGMP_SNOOPING_ENABLE;

    case SAI_VLAN_ATTR_TAM_OBJECT:
      return lemming::dataplane::sai::VLAN_ATTR_TAM_OBJECT;

    default:
      return lemming::dataplane::sai::VLAN_ATTR_UNSPECIFIED;
  }
}
sai_vlan_attr_t convert_sai_vlan_attr_t_to_sai(
    lemming::dataplane::sai::VlanAttr val) {
  switch (val) {
    case lemming::dataplane::sai::VLAN_ATTR_VLAN_ID:
      return SAI_VLAN_ATTR_VLAN_ID;

    case lemming::dataplane::sai::VLAN_ATTR_MEMBER_LIST:
      return SAI_VLAN_ATTR_MEMBER_LIST;

    case lemming::dataplane::sai::VLAN_ATTR_MAX_LEARNED_ADDRESSES:
      return SAI_VLAN_ATTR_MAX_LEARNED_ADDRESSES;

    case lemming::dataplane::sai::VLAN_ATTR_STP_INSTANCE:
      return SAI_VLAN_ATTR_STP_INSTANCE;

    case lemming::dataplane::sai::VLAN_ATTR_LEARN_DISABLE:
      return SAI_VLAN_ATTR_LEARN_DISABLE;

    case lemming::dataplane::sai::VLAN_ATTR_IPV4_MCAST_LOOKUP_KEY_TYPE:
      return SAI_VLAN_ATTR_IPV4_MCAST_LOOKUP_KEY_TYPE;

    case lemming::dataplane::sai::VLAN_ATTR_IPV6_MCAST_LOOKUP_KEY_TYPE:
      return SAI_VLAN_ATTR_IPV6_MCAST_LOOKUP_KEY_TYPE;

    case lemming::dataplane::sai::
        VLAN_ATTR_UNKNOWN_NON_IP_MCAST_OUTPUT_GROUP_ID:
      return SAI_VLAN_ATTR_UNKNOWN_NON_IP_MCAST_OUTPUT_GROUP_ID;

    case lemming::dataplane::sai::VLAN_ATTR_UNKNOWN_IPV4_MCAST_OUTPUT_GROUP_ID:
      return SAI_VLAN_ATTR_UNKNOWN_IPV4_MCAST_OUTPUT_GROUP_ID;

    case lemming::dataplane::sai::VLAN_ATTR_UNKNOWN_IPV6_MCAST_OUTPUT_GROUP_ID:
      return SAI_VLAN_ATTR_UNKNOWN_IPV6_MCAST_OUTPUT_GROUP_ID;

    case lemming::dataplane::sai::
        VLAN_ATTR_UNKNOWN_LINKLOCAL_MCAST_OUTPUT_GROUP_ID:
      return SAI_VLAN_ATTR_UNKNOWN_LINKLOCAL_MCAST_OUTPUT_GROUP_ID;

    case lemming::dataplane::sai::VLAN_ATTR_INGRESS_ACL:
      return SAI_VLAN_ATTR_INGRESS_ACL;

    case lemming::dataplane::sai::VLAN_ATTR_EGRESS_ACL:
      return SAI_VLAN_ATTR_EGRESS_ACL;

    case lemming::dataplane::sai::VLAN_ATTR_META_DATA:
      return SAI_VLAN_ATTR_META_DATA;

    case lemming::dataplane::sai::VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE:
      return SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_CONTROL_TYPE;

    case lemming::dataplane::sai::VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP:
      return SAI_VLAN_ATTR_UNKNOWN_UNICAST_FLOOD_GROUP;

    case lemming::dataplane::sai::
        VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE:
      return SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_CONTROL_TYPE;

    case lemming::dataplane::sai::VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP:
      return SAI_VLAN_ATTR_UNKNOWN_MULTICAST_FLOOD_GROUP;

    case lemming::dataplane::sai::VLAN_ATTR_BROADCAST_FLOOD_CONTROL_TYPE:
      return SAI_VLAN_ATTR_BROADCAST_FLOOD_CONTROL_TYPE;

    case lemming::dataplane::sai::VLAN_ATTR_BROADCAST_FLOOD_GROUP:
      return SAI_VLAN_ATTR_BROADCAST_FLOOD_GROUP;

    case lemming::dataplane::sai::VLAN_ATTR_CUSTOM_IGMP_SNOOPING_ENABLE:
      return SAI_VLAN_ATTR_CUSTOM_IGMP_SNOOPING_ENABLE;

    case lemming::dataplane::sai::VLAN_ATTR_TAM_OBJECT:
      return SAI_VLAN_ATTR_TAM_OBJECT;

    default:
      return SAI_VLAN_ATTR_VLAN_ID;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_vlan_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_vlan_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_vlan_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_vlan_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::VlanAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::VlanFloodControlType
convert_sai_vlan_flood_control_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_VLAN_FLOOD_CONTROL_TYPE_ALL:
      return lemming::dataplane::sai::VLAN_FLOOD_CONTROL_TYPE_ALL;

    case SAI_VLAN_FLOOD_CONTROL_TYPE_NONE:
      return lemming::dataplane::sai::VLAN_FLOOD_CONTROL_TYPE_NONE;

    case SAI_VLAN_FLOOD_CONTROL_TYPE_L2MC_GROUP:
      return lemming::dataplane::sai::VLAN_FLOOD_CONTROL_TYPE_L2MC_GROUP;

    case SAI_VLAN_FLOOD_CONTROL_TYPE_COMBINED:
      return lemming::dataplane::sai::VLAN_FLOOD_CONTROL_TYPE_COMBINED;

    default:
      return lemming::dataplane::sai::VLAN_FLOOD_CONTROL_TYPE_UNSPECIFIED;
  }
}
sai_vlan_flood_control_type_t convert_sai_vlan_flood_control_type_t_to_sai(
    lemming::dataplane::sai::VlanFloodControlType val) {
  switch (val) {
    case lemming::dataplane::sai::VLAN_FLOOD_CONTROL_TYPE_ALL:
      return SAI_VLAN_FLOOD_CONTROL_TYPE_ALL;

    case lemming::dataplane::sai::VLAN_FLOOD_CONTROL_TYPE_NONE:
      return SAI_VLAN_FLOOD_CONTROL_TYPE_NONE;

    case lemming::dataplane::sai::VLAN_FLOOD_CONTROL_TYPE_L2MC_GROUP:
      return SAI_VLAN_FLOOD_CONTROL_TYPE_L2MC_GROUP;

    case lemming::dataplane::sai::VLAN_FLOOD_CONTROL_TYPE_COMBINED:
      return SAI_VLAN_FLOOD_CONTROL_TYPE_COMBINED;

    default:
      return SAI_VLAN_FLOOD_CONTROL_TYPE_ALL;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_vlan_flood_control_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_vlan_flood_control_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_vlan_flood_control_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_vlan_flood_control_type_t_to_sai(
        static_cast<lemming::dataplane::sai::VlanFloodControlType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::VlanMcastLookupKeyType
convert_sai_vlan_mcast_lookup_key_type_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_VLAN_MCAST_LOOKUP_KEY_TYPE_MAC_DA:
      return lemming::dataplane::sai::VLAN_MCAST_LOOKUP_KEY_TYPE_MAC_DA;

    case SAI_VLAN_MCAST_LOOKUP_KEY_TYPE_XG:
      return lemming::dataplane::sai::VLAN_MCAST_LOOKUP_KEY_TYPE_XG;

    case SAI_VLAN_MCAST_LOOKUP_KEY_TYPE_SG:
      return lemming::dataplane::sai::VLAN_MCAST_LOOKUP_KEY_TYPE_SG;

    case SAI_VLAN_MCAST_LOOKUP_KEY_TYPE_XG_AND_SG:
      return lemming::dataplane::sai::VLAN_MCAST_LOOKUP_KEY_TYPE_XG_AND_SG;

    default:
      return lemming::dataplane::sai::VLAN_MCAST_LOOKUP_KEY_TYPE_UNSPECIFIED;
  }
}
sai_vlan_mcast_lookup_key_type_t
convert_sai_vlan_mcast_lookup_key_type_t_to_sai(
    lemming::dataplane::sai::VlanMcastLookupKeyType val) {
  switch (val) {
    case lemming::dataplane::sai::VLAN_MCAST_LOOKUP_KEY_TYPE_MAC_DA:
      return SAI_VLAN_MCAST_LOOKUP_KEY_TYPE_MAC_DA;

    case lemming::dataplane::sai::VLAN_MCAST_LOOKUP_KEY_TYPE_XG:
      return SAI_VLAN_MCAST_LOOKUP_KEY_TYPE_XG;

    case lemming::dataplane::sai::VLAN_MCAST_LOOKUP_KEY_TYPE_SG:
      return SAI_VLAN_MCAST_LOOKUP_KEY_TYPE_SG;

    case lemming::dataplane::sai::VLAN_MCAST_LOOKUP_KEY_TYPE_XG_AND_SG:
      return SAI_VLAN_MCAST_LOOKUP_KEY_TYPE_XG_AND_SG;

    default:
      return SAI_VLAN_MCAST_LOOKUP_KEY_TYPE_MAC_DA;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_vlan_mcast_lookup_key_type_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(
        convert_sai_vlan_mcast_lookup_key_type_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_vlan_mcast_lookup_key_type_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_vlan_mcast_lookup_key_type_t_to_sai(
        static_cast<lemming::dataplane::sai::VlanMcastLookupKeyType>(
            proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::VlanMemberAttr convert_sai_vlan_member_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_VLAN_MEMBER_ATTR_VLAN_ID:
      return lemming::dataplane::sai::VLAN_MEMBER_ATTR_VLAN_ID;

    case SAI_VLAN_MEMBER_ATTR_BRIDGE_PORT_ID:
      return lemming::dataplane::sai::VLAN_MEMBER_ATTR_BRIDGE_PORT_ID;

    case SAI_VLAN_MEMBER_ATTR_VLAN_TAGGING_MODE:
      return lemming::dataplane::sai::VLAN_MEMBER_ATTR_VLAN_TAGGING_MODE;

    default:
      return lemming::dataplane::sai::VLAN_MEMBER_ATTR_UNSPECIFIED;
  }
}
sai_vlan_member_attr_t convert_sai_vlan_member_attr_t_to_sai(
    lemming::dataplane::sai::VlanMemberAttr val) {
  switch (val) {
    case lemming::dataplane::sai::VLAN_MEMBER_ATTR_VLAN_ID:
      return SAI_VLAN_MEMBER_ATTR_VLAN_ID;

    case lemming::dataplane::sai::VLAN_MEMBER_ATTR_BRIDGE_PORT_ID:
      return SAI_VLAN_MEMBER_ATTR_BRIDGE_PORT_ID;

    case lemming::dataplane::sai::VLAN_MEMBER_ATTR_VLAN_TAGGING_MODE:
      return SAI_VLAN_MEMBER_ATTR_VLAN_TAGGING_MODE;

    default:
      return SAI_VLAN_MEMBER_ATTR_VLAN_ID;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_vlan_member_attr_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_vlan_member_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_vlan_member_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_vlan_member_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::VlanMemberAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::VlanStat convert_sai_vlan_stat_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_VLAN_STAT_IN_OCTETS:
      return lemming::dataplane::sai::VLAN_STAT_IN_OCTETS;

    case SAI_VLAN_STAT_IN_PACKETS:
      return lemming::dataplane::sai::VLAN_STAT_IN_PACKETS;

    case SAI_VLAN_STAT_IN_UCAST_PKTS:
      return lemming::dataplane::sai::VLAN_STAT_IN_UCAST_PKTS;

    case SAI_VLAN_STAT_IN_NON_UCAST_PKTS:
      return lemming::dataplane::sai::VLAN_STAT_IN_NON_UCAST_PKTS;

    case SAI_VLAN_STAT_IN_DISCARDS:
      return lemming::dataplane::sai::VLAN_STAT_IN_DISCARDS;

    case SAI_VLAN_STAT_IN_ERRORS:
      return lemming::dataplane::sai::VLAN_STAT_IN_ERRORS;

    case SAI_VLAN_STAT_IN_UNKNOWN_PROTOS:
      return lemming::dataplane::sai::VLAN_STAT_IN_UNKNOWN_PROTOS;

    case SAI_VLAN_STAT_OUT_OCTETS:
      return lemming::dataplane::sai::VLAN_STAT_OUT_OCTETS;

    case SAI_VLAN_STAT_OUT_PACKETS:
      return lemming::dataplane::sai::VLAN_STAT_OUT_PACKETS;

    case SAI_VLAN_STAT_OUT_UCAST_PKTS:
      return lemming::dataplane::sai::VLAN_STAT_OUT_UCAST_PKTS;

    case SAI_VLAN_STAT_OUT_NON_UCAST_PKTS:
      return lemming::dataplane::sai::VLAN_STAT_OUT_NON_UCAST_PKTS;

    case SAI_VLAN_STAT_OUT_DISCARDS:
      return lemming::dataplane::sai::VLAN_STAT_OUT_DISCARDS;

    case SAI_VLAN_STAT_OUT_ERRORS:
      return lemming::dataplane::sai::VLAN_STAT_OUT_ERRORS;

    case SAI_VLAN_STAT_OUT_QLEN:
      return lemming::dataplane::sai::VLAN_STAT_OUT_QLEN;

    default:
      return lemming::dataplane::sai::VLAN_STAT_UNSPECIFIED;
  }
}
sai_vlan_stat_t convert_sai_vlan_stat_t_to_sai(
    lemming::dataplane::sai::VlanStat val) {
  switch (val) {
    case lemming::dataplane::sai::VLAN_STAT_IN_OCTETS:
      return SAI_VLAN_STAT_IN_OCTETS;

    case lemming::dataplane::sai::VLAN_STAT_IN_PACKETS:
      return SAI_VLAN_STAT_IN_PACKETS;

    case lemming::dataplane::sai::VLAN_STAT_IN_UCAST_PKTS:
      return SAI_VLAN_STAT_IN_UCAST_PKTS;

    case lemming::dataplane::sai::VLAN_STAT_IN_NON_UCAST_PKTS:
      return SAI_VLAN_STAT_IN_NON_UCAST_PKTS;

    case lemming::dataplane::sai::VLAN_STAT_IN_DISCARDS:
      return SAI_VLAN_STAT_IN_DISCARDS;

    case lemming::dataplane::sai::VLAN_STAT_IN_ERRORS:
      return SAI_VLAN_STAT_IN_ERRORS;

    case lemming::dataplane::sai::VLAN_STAT_IN_UNKNOWN_PROTOS:
      return SAI_VLAN_STAT_IN_UNKNOWN_PROTOS;

    case lemming::dataplane::sai::VLAN_STAT_OUT_OCTETS:
      return SAI_VLAN_STAT_OUT_OCTETS;

    case lemming::dataplane::sai::VLAN_STAT_OUT_PACKETS:
      return SAI_VLAN_STAT_OUT_PACKETS;

    case lemming::dataplane::sai::VLAN_STAT_OUT_UCAST_PKTS:
      return SAI_VLAN_STAT_OUT_UCAST_PKTS;

    case lemming::dataplane::sai::VLAN_STAT_OUT_NON_UCAST_PKTS:
      return SAI_VLAN_STAT_OUT_NON_UCAST_PKTS;

    case lemming::dataplane::sai::VLAN_STAT_OUT_DISCARDS:
      return SAI_VLAN_STAT_OUT_DISCARDS;

    case lemming::dataplane::sai::VLAN_STAT_OUT_ERRORS:
      return SAI_VLAN_STAT_OUT_ERRORS;

    case lemming::dataplane::sai::VLAN_STAT_OUT_QLEN:
      return SAI_VLAN_STAT_OUT_QLEN;

    default:
      return SAI_VLAN_STAT_IN_OCTETS;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_vlan_stat_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_vlan_stat_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_vlan_stat_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_vlan_stat_t_to_sai(
        static_cast<lemming::dataplane::sai::VlanStat>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::VlanTaggingMode
convert_sai_vlan_tagging_mode_t_to_proto(const sai_int32_t val) {
  switch (val) {
    case SAI_VLAN_TAGGING_MODE_UNTAGGED:
      return lemming::dataplane::sai::VLAN_TAGGING_MODE_UNTAGGED;

    case SAI_VLAN_TAGGING_MODE_TAGGED:
      return lemming::dataplane::sai::VLAN_TAGGING_MODE_TAGGED;

    case SAI_VLAN_TAGGING_MODE_PRIORITY_TAGGED:
      return lemming::dataplane::sai::VLAN_TAGGING_MODE_PRIORITY_TAGGED;

    default:
      return lemming::dataplane::sai::VLAN_TAGGING_MODE_UNSPECIFIED;
  }
}
sai_vlan_tagging_mode_t convert_sai_vlan_tagging_mode_t_to_sai(
    lemming::dataplane::sai::VlanTaggingMode val) {
  switch (val) {
    case lemming::dataplane::sai::VLAN_TAGGING_MODE_UNTAGGED:
      return SAI_VLAN_TAGGING_MODE_UNTAGGED;

    case lemming::dataplane::sai::VLAN_TAGGING_MODE_TAGGED:
      return SAI_VLAN_TAGGING_MODE_TAGGED;

    case lemming::dataplane::sai::VLAN_TAGGING_MODE_PRIORITY_TAGGED:
      return SAI_VLAN_TAGGING_MODE_PRIORITY_TAGGED;

    default:
      return SAI_VLAN_TAGGING_MODE_UNTAGGED;
  }
}

google::protobuf::RepeatedField<int>
convert_list_sai_vlan_tagging_mode_t_to_proto(const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_vlan_tagging_mode_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_vlan_tagging_mode_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_vlan_tagging_mode_t_to_sai(
        static_cast<lemming::dataplane::sai::VlanTaggingMode>(proto_list[i]));
  }
  *count = proto_list.size();
}

lemming::dataplane::sai::WredAttr convert_sai_wred_attr_t_to_proto(
    const sai_int32_t val) {
  switch (val) {
    case SAI_WRED_ATTR_GREEN_ENABLE:
      return lemming::dataplane::sai::WRED_ATTR_GREEN_ENABLE;

    case SAI_WRED_ATTR_GREEN_MIN_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_GREEN_MIN_THRESHOLD;

    case SAI_WRED_ATTR_GREEN_MAX_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_GREEN_MAX_THRESHOLD;

    case SAI_WRED_ATTR_GREEN_DROP_PROBABILITY:
      return lemming::dataplane::sai::WRED_ATTR_GREEN_DROP_PROBABILITY;

    case SAI_WRED_ATTR_YELLOW_ENABLE:
      return lemming::dataplane::sai::WRED_ATTR_YELLOW_ENABLE;

    case SAI_WRED_ATTR_YELLOW_MIN_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_YELLOW_MIN_THRESHOLD;

    case SAI_WRED_ATTR_YELLOW_MAX_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_YELLOW_MAX_THRESHOLD;

    case SAI_WRED_ATTR_YELLOW_DROP_PROBABILITY:
      return lemming::dataplane::sai::WRED_ATTR_YELLOW_DROP_PROBABILITY;

    case SAI_WRED_ATTR_RED_ENABLE:
      return lemming::dataplane::sai::WRED_ATTR_RED_ENABLE;

    case SAI_WRED_ATTR_RED_MIN_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_RED_MIN_THRESHOLD;

    case SAI_WRED_ATTR_RED_MAX_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_RED_MAX_THRESHOLD;

    case SAI_WRED_ATTR_RED_DROP_PROBABILITY:
      return lemming::dataplane::sai::WRED_ATTR_RED_DROP_PROBABILITY;

    case SAI_WRED_ATTR_WEIGHT:
      return lemming::dataplane::sai::WRED_ATTR_WEIGHT;

    case SAI_WRED_ATTR_ECN_MARK_MODE:
      return lemming::dataplane::sai::WRED_ATTR_ECN_MARK_MODE;

    case SAI_WRED_ATTR_ECN_GREEN_MIN_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_ECN_GREEN_MIN_THRESHOLD;

    case SAI_WRED_ATTR_ECN_GREEN_MAX_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_ECN_GREEN_MAX_THRESHOLD;

    case SAI_WRED_ATTR_ECN_GREEN_MARK_PROBABILITY:
      return lemming::dataplane::sai::WRED_ATTR_ECN_GREEN_MARK_PROBABILITY;

    case SAI_WRED_ATTR_ECN_YELLOW_MIN_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_ECN_YELLOW_MIN_THRESHOLD;

    case SAI_WRED_ATTR_ECN_YELLOW_MAX_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_ECN_YELLOW_MAX_THRESHOLD;

    case SAI_WRED_ATTR_ECN_YELLOW_MARK_PROBABILITY:
      return lemming::dataplane::sai::WRED_ATTR_ECN_YELLOW_MARK_PROBABILITY;

    case SAI_WRED_ATTR_ECN_RED_MIN_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_ECN_RED_MIN_THRESHOLD;

    case SAI_WRED_ATTR_ECN_RED_MAX_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_ECN_RED_MAX_THRESHOLD;

    case SAI_WRED_ATTR_ECN_RED_MARK_PROBABILITY:
      return lemming::dataplane::sai::WRED_ATTR_ECN_RED_MARK_PROBABILITY;

    case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MIN_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_ECN_COLOR_UNAWARE_MIN_THRESHOLD;

    case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MAX_THRESHOLD:
      return lemming::dataplane::sai::WRED_ATTR_ECN_COLOR_UNAWARE_MAX_THRESHOLD;

    case SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MARK_PROBABILITY:
      return lemming::dataplane::sai::
          WRED_ATTR_ECN_COLOR_UNAWARE_MARK_PROBABILITY;

    default:
      return lemming::dataplane::sai::WRED_ATTR_UNSPECIFIED;
  }
}
sai_wred_attr_t convert_sai_wred_attr_t_to_sai(
    lemming::dataplane::sai::WredAttr val) {
  switch (val) {
    case lemming::dataplane::sai::WRED_ATTR_GREEN_ENABLE:
      return SAI_WRED_ATTR_GREEN_ENABLE;

    case lemming::dataplane::sai::WRED_ATTR_GREEN_MIN_THRESHOLD:
      return SAI_WRED_ATTR_GREEN_MIN_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_GREEN_MAX_THRESHOLD:
      return SAI_WRED_ATTR_GREEN_MAX_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_GREEN_DROP_PROBABILITY:
      return SAI_WRED_ATTR_GREEN_DROP_PROBABILITY;

    case lemming::dataplane::sai::WRED_ATTR_YELLOW_ENABLE:
      return SAI_WRED_ATTR_YELLOW_ENABLE;

    case lemming::dataplane::sai::WRED_ATTR_YELLOW_MIN_THRESHOLD:
      return SAI_WRED_ATTR_YELLOW_MIN_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_YELLOW_MAX_THRESHOLD:
      return SAI_WRED_ATTR_YELLOW_MAX_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_YELLOW_DROP_PROBABILITY:
      return SAI_WRED_ATTR_YELLOW_DROP_PROBABILITY;

    case lemming::dataplane::sai::WRED_ATTR_RED_ENABLE:
      return SAI_WRED_ATTR_RED_ENABLE;

    case lemming::dataplane::sai::WRED_ATTR_RED_MIN_THRESHOLD:
      return SAI_WRED_ATTR_RED_MIN_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_RED_MAX_THRESHOLD:
      return SAI_WRED_ATTR_RED_MAX_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_RED_DROP_PROBABILITY:
      return SAI_WRED_ATTR_RED_DROP_PROBABILITY;

    case lemming::dataplane::sai::WRED_ATTR_WEIGHT:
      return SAI_WRED_ATTR_WEIGHT;

    case lemming::dataplane::sai::WRED_ATTR_ECN_MARK_MODE:
      return SAI_WRED_ATTR_ECN_MARK_MODE;

    case lemming::dataplane::sai::WRED_ATTR_ECN_GREEN_MIN_THRESHOLD:
      return SAI_WRED_ATTR_ECN_GREEN_MIN_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_ECN_GREEN_MAX_THRESHOLD:
      return SAI_WRED_ATTR_ECN_GREEN_MAX_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_ECN_GREEN_MARK_PROBABILITY:
      return SAI_WRED_ATTR_ECN_GREEN_MARK_PROBABILITY;

    case lemming::dataplane::sai::WRED_ATTR_ECN_YELLOW_MIN_THRESHOLD:
      return SAI_WRED_ATTR_ECN_YELLOW_MIN_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_ECN_YELLOW_MAX_THRESHOLD:
      return SAI_WRED_ATTR_ECN_YELLOW_MAX_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_ECN_YELLOW_MARK_PROBABILITY:
      return SAI_WRED_ATTR_ECN_YELLOW_MARK_PROBABILITY;

    case lemming::dataplane::sai::WRED_ATTR_ECN_RED_MIN_THRESHOLD:
      return SAI_WRED_ATTR_ECN_RED_MIN_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_ECN_RED_MAX_THRESHOLD:
      return SAI_WRED_ATTR_ECN_RED_MAX_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_ECN_RED_MARK_PROBABILITY:
      return SAI_WRED_ATTR_ECN_RED_MARK_PROBABILITY;

    case lemming::dataplane::sai::WRED_ATTR_ECN_COLOR_UNAWARE_MIN_THRESHOLD:
      return SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MIN_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_ECN_COLOR_UNAWARE_MAX_THRESHOLD:
      return SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MAX_THRESHOLD;

    case lemming::dataplane::sai::WRED_ATTR_ECN_COLOR_UNAWARE_MARK_PROBABILITY:
      return SAI_WRED_ATTR_ECN_COLOR_UNAWARE_MARK_PROBABILITY;

    default:
      return SAI_WRED_ATTR_GREEN_ENABLE;
  }
}

google::protobuf::RepeatedField<int> convert_list_sai_wred_attr_t_to_proto(
    const sai_s32_list_t &list) {
  google::protobuf::RepeatedField<int> proto_list;
  for (int i = 0; i < list.count; i++) {
    proto_list.Add(convert_sai_wred_attr_t_to_proto(list.list[i]));
  }
  return proto_list;
}
void convert_list_sai_wred_attr_t_to_sai(
    int32_t *list, const google::protobuf::RepeatedField<int> &proto_list,
    uint32_t *count) {
  for (int i = 0; i < proto_list.size(); i++) {
    list[i] = convert_sai_wred_attr_t_to_sai(
        static_cast<lemming::dataplane::sai::WredAttr>(proto_list[i]));
  }
  *count = proto_list.size();
}
