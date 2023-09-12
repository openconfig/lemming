

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

#include "dataplane/standalone/sai/l2mc_group.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/l2mc_group.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_l2mc_group_api_t l_l2mc_group = {
    .create_l2mc_group = l_create_l2mc_group,
    .remove_l2mc_group = l_remove_l2mc_group,
    .set_l2mc_group_attribute = l_set_l2mc_group_attribute,
    .get_l2mc_group_attribute = l_get_l2mc_group_attribute,
    .create_l2mc_group_member = l_create_l2mc_group_member,
    .remove_l2mc_group_member = l_remove_l2mc_group_member,
    .set_l2mc_group_member_attribute = l_set_l2mc_group_member_attribute,
    .get_l2mc_group_member_attribute = l_get_l2mc_group_member_attribute,
};

sai_status_t l_create_l2mc_group(sai_object_id_t *l2mc_group_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateL2mcGroupRequest req;
  lemming::dataplane::sai::CreateL2mcGroupResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {}
  }
  grpc::Status status = l2mc_group->CreateL2mcGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *l2mc_group_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_l2mc_group(sai_object_id_t l2mc_group_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveL2mcGroupRequest req;
  lemming::dataplane::sai::RemoveL2mcGroupResponse resp;
  grpc::ClientContext context;
  req.set_oid(l2mc_group_id);

  grpc::Status status = l2mc_group->RemoveL2mcGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_l2mc_group_attribute(sai_object_id_t l2mc_group_id,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_l2mc_group_attribute(sai_object_id_t l2mc_group_id,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetL2mcGroupAttributeRequest req;
  lemming::dataplane::sai::GetL2mcGroupAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(l2mc_group_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::L2mcGroupAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = l2mc_group->GetL2mcGroupAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_L2MC_GROUP_ATTR_L2MC_OUTPUT_COUNT:
        attr_list[i].value.u32 = resp.attr().l2mc_output_count();
        break;
      case SAI_L2MC_GROUP_ATTR_L2MC_MEMBER_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().l2mc_member_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_l2mc_group_member(sai_object_id_t *l2mc_group_member_id,
                                        sai_object_id_t switch_id,
                                        uint32_t attr_count,
                                        const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateL2mcGroupMemberRequest req;
  lemming::dataplane::sai::CreateL2mcGroupMemberResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_GROUP_ID:
        req.set_l2mc_group_id(attr_list[i].value.oid);
        break;
      case SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_OUTPUT_ID:
        req.set_l2mc_output_id(attr_list[i].value.oid);
        break;
      case SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_ENDPOINT_IP:
        req.set_l2mc_endpoint_ip(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
    }
  }
  grpc::Status status = l2mc_group->CreateL2mcGroupMember(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *l2mc_group_member_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_l2mc_group_member(sai_object_id_t l2mc_group_member_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveL2mcGroupMemberRequest req;
  lemming::dataplane::sai::RemoveL2mcGroupMemberResponse resp;
  grpc::ClientContext context;
  req.set_oid(l2mc_group_member_id);

  grpc::Status status = l2mc_group->RemoveL2mcGroupMember(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_l2mc_group_member_attribute(
    sai_object_id_t l2mc_group_member_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_l2mc_group_member_attribute(
    sai_object_id_t l2mc_group_member_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetL2mcGroupMemberAttributeRequest req;
  lemming::dataplane::sai::GetL2mcGroupMemberAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(l2mc_group_member_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::L2mcGroupMemberAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status =
      l2mc_group->GetL2mcGroupMemberAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_GROUP_ID:
        attr_list[i].value.oid = resp.attr().l2mc_group_id();
        break;
      case SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_OUTPUT_ID:
        attr_list[i].value.oid = resp.attr().l2mc_output_id();
        break;
      case SAI_L2MC_GROUP_MEMBER_ATTR_L2MC_ENDPOINT_IP:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().l2mc_endpoint_ip());
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
