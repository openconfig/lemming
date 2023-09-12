

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

#include "dataplane/standalone/sai/isolation_group.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/isolation_group.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_isolation_group_api_t l_isolation_group = {
    .create_isolation_group = l_create_isolation_group,
    .remove_isolation_group = l_remove_isolation_group,
    .set_isolation_group_attribute = l_set_isolation_group_attribute,
    .get_isolation_group_attribute = l_get_isolation_group_attribute,
    .create_isolation_group_member = l_create_isolation_group_member,
    .remove_isolation_group_member = l_remove_isolation_group_member,
    .set_isolation_group_member_attribute =
        l_set_isolation_group_member_attribute,
    .get_isolation_group_member_attribute =
        l_get_isolation_group_member_attribute,
};

sai_status_t l_create_isolation_group(sai_object_id_t *isolation_group_id,
                                      sai_object_id_t switch_id,
                                      uint32_t attr_count,
                                      const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIsolationGroupRequest req;
  lemming::dataplane::sai::CreateIsolationGroupResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ISOLATION_GROUP_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::IsolationGroupType>(
            attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status =
      isolation_group->CreateIsolationGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *isolation_group_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_isolation_group(sai_object_id_t isolation_group_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveIsolationGroupRequest req;
  lemming::dataplane::sai::RemoveIsolationGroupResponse resp;
  grpc::ClientContext context;
  req.set_oid(isolation_group_id);

  grpc::Status status =
      isolation_group->RemoveIsolationGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_isolation_group_attribute(sai_object_id_t isolation_group_id,
                                             const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_isolation_group_attribute(sai_object_id_t isolation_group_id,
                                             uint32_t attr_count,
                                             sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIsolationGroupAttributeRequest req;
  lemming::dataplane::sai::GetIsolationGroupAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(isolation_group_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::IsolationGroupAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status =
      isolation_group->GetIsolationGroupAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ISOLATION_GROUP_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_ISOLATION_GROUP_ATTR_ISOLATION_MEMBER_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().isolation_member_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_isolation_group_member(
    sai_object_id_t *isolation_group_member_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateIsolationGroupMemberRequest req;
  lemming::dataplane::sai::CreateIsolationGroupMemberResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_GROUP_ID:
        req.set_isolation_group_id(attr_list[i].value.oid);
        break;
      case SAI_ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_OBJECT:
        req.set_isolation_object(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status =
      isolation_group->CreateIsolationGroupMember(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *isolation_group_member_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_isolation_group_member(
    sai_object_id_t isolation_group_member_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveIsolationGroupMemberRequest req;
  lemming::dataplane::sai::RemoveIsolationGroupMemberResponse resp;
  grpc::ClientContext context;
  req.set_oid(isolation_group_member_id);

  grpc::Status status =
      isolation_group->RemoveIsolationGroupMember(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_isolation_group_member_attribute(
    sai_object_id_t isolation_group_member_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_isolation_group_member_attribute(
    sai_object_id_t isolation_group_member_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetIsolationGroupMemberAttributeRequest req;
  lemming::dataplane::sai::GetIsolationGroupMemberAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(isolation_group_member_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::IsolationGroupMemberAttr>(
            attr_list[i].id + 1));
  }
  grpc::Status status =
      isolation_group->GetIsolationGroupMemberAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_GROUP_ID:
        attr_list[i].value.oid = resp.attr().isolation_group_id();
        break;
      case SAI_ISOLATION_GROUP_MEMBER_ATTR_ISOLATION_OBJECT:
        attr_list[i].value.oid = resp.attr().isolation_object();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
