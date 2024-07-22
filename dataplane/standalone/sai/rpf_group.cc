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

#include "dataplane/standalone/sai/rpf_group.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/rpf_group.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_rpf_group_api_t l_rpf_group = {
    .create_rpf_group = l_create_rpf_group,
    .remove_rpf_group = l_remove_rpf_group,
    .set_rpf_group_attribute = l_set_rpf_group_attribute,
    .get_rpf_group_attribute = l_get_rpf_group_attribute,
    .create_rpf_group_member = l_create_rpf_group_member,
    .remove_rpf_group_member = l_remove_rpf_group_member,
    .set_rpf_group_member_attribute = l_set_rpf_group_member_attribute,
    .get_rpf_group_member_attribute = l_get_rpf_group_member_attribute,
};

lemming::dataplane::sai::CreateRpfGroupRequest convert_create_rpf_group(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateRpfGroupRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {}
  }
  return msg;
}

lemming::dataplane::sai::CreateRpfGroupMemberRequest
convert_create_rpf_group_member(sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateRpfGroupMemberRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_RPF_GROUP_MEMBER_ATTR_RPF_GROUP_ID:
        msg.set_rpf_group_id(attr_list[i].value.oid);
        break;
      case SAI_RPF_GROUP_MEMBER_ATTR_RPF_INTERFACE_ID:
        msg.set_rpf_interface_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_rpf_group(sai_object_id_t *rpf_group_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateRpfGroupRequest req =
      convert_create_rpf_group(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateRpfGroupResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = rpf_group->CreateRpfGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  if (rpf_group_id) {
    *rpf_group_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_rpf_group(sai_object_id_t rpf_group_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveRpfGroupRequest req;
  lemming::dataplane::sai::RemoveRpfGroupResponse resp;
  grpc::ClientContext context;
  req.set_oid(rpf_group_id);

  grpc::Status status = rpf_group->RemoveRpfGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_rpf_group_attribute(sai_object_id_t rpf_group_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_rpf_group_attribute(sai_object_id_t rpf_group_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetRpfGroupAttributeRequest req;
  lemming::dataplane::sai::GetRpfGroupAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(rpf_group_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_rpf_group_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = rpf_group->GetRpfGroupAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_RPF_GROUP_ATTR_RPF_INTERFACE_COUNT:
        attr_list[i].value.u32 = resp.attr().rpf_interface_count();
        break;
      case SAI_RPF_GROUP_ATTR_RPF_MEMBER_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().rpf_member_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_rpf_group_member(sai_object_id_t *rpf_group_member_id,
                                       sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateRpfGroupMemberRequest req =
      convert_create_rpf_group_member(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateRpfGroupMemberResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = rpf_group->CreateRpfGroupMember(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  if (rpf_group_member_id) {
    *rpf_group_member_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_rpf_group_member(sai_object_id_t rpf_group_member_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveRpfGroupMemberRequest req;
  lemming::dataplane::sai::RemoveRpfGroupMemberResponse resp;
  grpc::ClientContext context;
  req.set_oid(rpf_group_member_id);

  grpc::Status status = rpf_group->RemoveRpfGroupMember(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_rpf_group_member_attribute(
    sai_object_id_t rpf_group_member_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_rpf_group_member_attribute(
    sai_object_id_t rpf_group_member_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetRpfGroupMemberAttributeRequest req;
  lemming::dataplane::sai::GetRpfGroupMemberAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(rpf_group_member_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_rpf_group_member_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status =
      rpf_group->GetRpfGroupMemberAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_RPF_GROUP_MEMBER_ATTR_RPF_GROUP_ID:
        attr_list[i].value.oid = resp.attr().rpf_group_id();
        break;
      case SAI_RPF_GROUP_MEMBER_ATTR_RPF_INTERFACE_ID:
        attr_list[i].value.oid = resp.attr().rpf_interface_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
