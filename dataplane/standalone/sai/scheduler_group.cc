

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

#include "dataplane/standalone/sai/scheduler_group.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/scheduler_group.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_scheduler_group_api_t l_scheduler_group = {
    .create_scheduler_group = l_create_scheduler_group,
    .remove_scheduler_group = l_remove_scheduler_group,
    .set_scheduler_group_attribute = l_set_scheduler_group_attribute,
    .get_scheduler_group_attribute = l_get_scheduler_group_attribute,
};

sai_status_t l_create_scheduler_group(sai_object_id_t *scheduler_group_id,
                                      sai_object_id_t switch_id,
                                      uint32_t attr_count,
                                      const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateSchedulerGroupRequest req;
  lemming::dataplane::sai::CreateSchedulerGroupResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SCHEDULER_GROUP_ATTR_PORT_ID:
        req.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_SCHEDULER_GROUP_ATTR_LEVEL:
        req.set_level(attr_list[i].value.u8);
        break;
      case SAI_SCHEDULER_GROUP_ATTR_MAX_CHILDS:
        req.set_max_childs(attr_list[i].value.u8);
        break;
      case SAI_SCHEDULER_GROUP_ATTR_SCHEDULER_PROFILE_ID:
        req.set_scheduler_profile_id(attr_list[i].value.oid);
        break;
      case SAI_SCHEDULER_GROUP_ATTR_PARENT_NODE:
        req.set_parent_node(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status =
      scheduler_group->CreateSchedulerGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *scheduler_group_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_scheduler_group(sai_object_id_t scheduler_group_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveSchedulerGroupRequest req;
  lemming::dataplane::sai::RemoveSchedulerGroupResponse resp;
  grpc::ClientContext context;
  req.set_oid(scheduler_group_id);

  grpc::Status status =
      scheduler_group->RemoveSchedulerGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_scheduler_group_attribute(sai_object_id_t scheduler_group_id,
                                             const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetSchedulerGroupAttributeRequest req;
  lemming::dataplane::sai::SetSchedulerGroupAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(scheduler_group_id);

  switch (attr->id) {
    case SAI_SCHEDULER_GROUP_ATTR_SCHEDULER_PROFILE_ID:
      req.set_scheduler_profile_id(attr->value.oid);
      break;
    case SAI_SCHEDULER_GROUP_ATTR_PARENT_NODE:
      req.set_parent_node(attr->value.oid);
      break;
  }

  grpc::Status status =
      scheduler_group->SetSchedulerGroupAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_scheduler_group_attribute(sai_object_id_t scheduler_group_id,
                                             uint32_t attr_count,
                                             sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetSchedulerGroupAttributeRequest req;
  lemming::dataplane::sai::GetSchedulerGroupAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(scheduler_group_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::SchedulerGroupAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status =
      scheduler_group->GetSchedulerGroupAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SCHEDULER_GROUP_ATTR_CHILD_COUNT:
        attr_list[i].value.u32 = resp.attr().child_count();
        break;
      case SAI_SCHEDULER_GROUP_ATTR_CHILD_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().child_list(),
                  attr_list[i].value.objlist.count);
        break;
      case SAI_SCHEDULER_GROUP_ATTR_PORT_ID:
        attr_list[i].value.oid = resp.attr().port_id();
        break;
      case SAI_SCHEDULER_GROUP_ATTR_LEVEL:
        attr_list[i].value.u8 = resp.attr().level();
        break;
      case SAI_SCHEDULER_GROUP_ATTR_MAX_CHILDS:
        attr_list[i].value.u8 = resp.attr().max_childs();
        break;
      case SAI_SCHEDULER_GROUP_ATTR_SCHEDULER_PROFILE_ID:
        attr_list[i].value.oid = resp.attr().scheduler_profile_id();
        break;
      case SAI_SCHEDULER_GROUP_ATTR_PARENT_NODE:
        attr_list[i].value.oid = resp.attr().parent_node();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
