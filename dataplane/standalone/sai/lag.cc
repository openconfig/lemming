

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

#include "dataplane/standalone/sai/lag.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/lag.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_lag_api_t l_lag = {
    .create_lag = l_create_lag,
    .remove_lag = l_remove_lag,
    .set_lag_attribute = l_set_lag_attribute,
    .get_lag_attribute = l_get_lag_attribute,
    .create_lag_member = l_create_lag_member,
    .remove_lag_member = l_remove_lag_member,
    .set_lag_member_attribute = l_set_lag_member_attribute,
    .get_lag_member_attribute = l_get_lag_member_attribute,
    .create_lag_members = l_create_lag_members,
    .remove_lag_members = l_remove_lag_members,
};

sai_status_t l_create_lag(sai_object_id_t *lag_id, sai_object_id_t switch_id,
                          uint32_t attr_count,
                          const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateLagRequest req;
  lemming::dataplane::sai::CreateLagResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_LAG_ATTR_INGRESS_ACL:
        req.set_ingress_acl(attr_list[i].value.oid);
        break;
      case SAI_LAG_ATTR_EGRESS_ACL:
        req.set_egress_acl(attr_list[i].value.oid);
        break;
      case SAI_LAG_ATTR_PORT_VLAN_ID:
        req.set_port_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_LAG_ATTR_DEFAULT_VLAN_PRIORITY:
        req.set_default_vlan_priority(attr_list[i].value.u8);
        break;
      case SAI_LAG_ATTR_DROP_UNTAGGED:
        req.set_drop_untagged(attr_list[i].value.booldata);
        break;
      case SAI_LAG_ATTR_DROP_TAGGED:
        req.set_drop_tagged(attr_list[i].value.booldata);
        break;
      case SAI_LAG_ATTR_TPID:
        req.set_tpid(attr_list[i].value.u16);
        break;
      case SAI_LAG_ATTR_SYSTEM_PORT_AGGREGATE_ID:
        req.set_system_port_aggregate_id(attr_list[i].value.u32);
        break;
      case SAI_LAG_ATTR_LABEL:
        req.set_label(attr_list[i].value.chardata);
        break;
    }
  }
  grpc::Status status = lag->CreateLag(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *lag_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_lag(sai_object_id_t lag_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveLagRequest req;
  lemming::dataplane::sai::RemoveLagResponse resp;
  grpc::ClientContext context;
  req.set_oid(lag_id);

  grpc::Status status = lag->RemoveLag(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_lag_attribute(sai_object_id_t lag_id,
                                 const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetLagAttributeRequest req;
  lemming::dataplane::sai::SetLagAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(lag_id);

  switch (attr->id) {
    case SAI_LAG_ATTR_INGRESS_ACL:
      req.set_ingress_acl(attr->value.oid);
      break;
    case SAI_LAG_ATTR_EGRESS_ACL:
      req.set_egress_acl(attr->value.oid);
      break;
    case SAI_LAG_ATTR_PORT_VLAN_ID:
      req.set_port_vlan_id(attr->value.u16);
      break;
    case SAI_LAG_ATTR_DEFAULT_VLAN_PRIORITY:
      req.set_default_vlan_priority(attr->value.u8);
      break;
    case SAI_LAG_ATTR_DROP_UNTAGGED:
      req.set_drop_untagged(attr->value.booldata);
      break;
    case SAI_LAG_ATTR_DROP_TAGGED:
      req.set_drop_tagged(attr->value.booldata);
      break;
    case SAI_LAG_ATTR_TPID:
      req.set_tpid(attr->value.u16);
      break;
    case SAI_LAG_ATTR_LABEL:
      req.set_label(attr->value.chardata);
      break;
  }

  grpc::Status status = lag->SetLagAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_lag_attribute(sai_object_id_t lag_id, uint32_t attr_count,
                                 sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetLagAttributeRequest req;
  lemming::dataplane::sai::GetLagAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(lag_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::LagAttr>(attr_list[i].id + 1));
  }
  grpc::Status status = lag->GetLagAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_LAG_ATTR_PORT_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().port_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_LAG_ATTR_INGRESS_ACL:
        attr_list[i].value.oid = resp.attr().ingress_acl();
        break;
      case SAI_LAG_ATTR_EGRESS_ACL:
        attr_list[i].value.oid = resp.attr().egress_acl();
        break;
      case SAI_LAG_ATTR_PORT_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().port_vlan_id();
        break;
      case SAI_LAG_ATTR_DEFAULT_VLAN_PRIORITY:
        attr_list[i].value.u8 = resp.attr().default_vlan_priority();
        break;
      case SAI_LAG_ATTR_DROP_UNTAGGED:
        attr_list[i].value.booldata = resp.attr().drop_untagged();
        break;
      case SAI_LAG_ATTR_DROP_TAGGED:
        attr_list[i].value.booldata = resp.attr().drop_tagged();
        break;
      case SAI_LAG_ATTR_TPID:
        attr_list[i].value.u16 = resp.attr().tpid();
        break;
      case SAI_LAG_ATTR_SYSTEM_PORT_AGGREGATE_ID:
        attr_list[i].value.u32 = resp.attr().system_port_aggregate_id();
        break;
      case SAI_LAG_ATTR_LABEL:
        strncpy(attr_list[i].value.chardata, resp.attr().label().data(), 32);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_lag_member(sai_object_id_t *lag_member_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateLagMemberRequest req;
  lemming::dataplane::sai::CreateLagMemberResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_LAG_MEMBER_ATTR_LAG_ID:
        req.set_lag_id(attr_list[i].value.oid);
        break;
      case SAI_LAG_MEMBER_ATTR_PORT_ID:
        req.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_LAG_MEMBER_ATTR_EGRESS_DISABLE:
        req.set_egress_disable(attr_list[i].value.booldata);
        break;
      case SAI_LAG_MEMBER_ATTR_INGRESS_DISABLE:
        req.set_ingress_disable(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status = lag->CreateLagMember(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *lag_member_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_lag_member(sai_object_id_t lag_member_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveLagMemberRequest req;
  lemming::dataplane::sai::RemoveLagMemberResponse resp;
  grpc::ClientContext context;
  req.set_oid(lag_member_id);

  grpc::Status status = lag->RemoveLagMember(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_lag_member_attribute(sai_object_id_t lag_member_id,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetLagMemberAttributeRequest req;
  lemming::dataplane::sai::SetLagMemberAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(lag_member_id);

  switch (attr->id) {
    case SAI_LAG_MEMBER_ATTR_EGRESS_DISABLE:
      req.set_egress_disable(attr->value.booldata);
      break;
    case SAI_LAG_MEMBER_ATTR_INGRESS_DISABLE:
      req.set_ingress_disable(attr->value.booldata);
      break;
  }

  grpc::Status status = lag->SetLagMemberAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_lag_member_attribute(sai_object_id_t lag_member_id,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetLagMemberAttributeRequest req;
  lemming::dataplane::sai::GetLagMemberAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(lag_member_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::LagMemberAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = lag->GetLagMemberAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_LAG_MEMBER_ATTR_LAG_ID:
        attr_list[i].value.oid = resp.attr().lag_id();
        break;
      case SAI_LAG_MEMBER_ATTR_PORT_ID:
        attr_list[i].value.oid = resp.attr().port_id();
        break;
      case SAI_LAG_MEMBER_ATTR_EGRESS_DISABLE:
        attr_list[i].value.booldata = resp.attr().egress_disable();
        break;
      case SAI_LAG_MEMBER_ATTR_INGRESS_DISABLE:
        attr_list[i].value.booldata = resp.attr().ingress_disable();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_lag_members(sai_object_id_t switch_id,
                                  uint32_t object_count,
                                  const uint32_t *attr_count,
                                  const sai_attribute_t **attr_list,
                                  sai_bulk_op_error_mode_t mode,
                                  sai_object_id_t *object_id,
                                  sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_lag_members(uint32_t object_count,
                                  const sai_object_id_t *object_id,
                                  sai_bulk_op_error_mode_t mode,
                                  sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
