

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

#include "dataplane/standalone/sai/system_port.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/system_port.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_system_port_api_t l_system_port = {
    .create_system_port = l_create_system_port,
    .remove_system_port = l_remove_system_port,
    .set_system_port_attribute = l_set_system_port_attribute,
    .get_system_port_attribute = l_get_system_port_attribute,
};

lemming::dataplane::sai::CreateSystemPortRequest convert_create_system_port(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateSystemPortRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SYSTEM_PORT_ATTR_ADMIN_STATE:
        msg.set_admin_state(attr_list[i].value.booldata);
        break;
      case SAI_SYSTEM_PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
        msg.set_qos_tc_to_queue_map(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_system_port(sai_object_id_t *system_port_id,
                                  sai_object_id_t switch_id,
                                  uint32_t attr_count,
                                  const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateSystemPortRequest req =
      convert_create_system_port(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateSystemPortResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = system_port->CreateSystemPort(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *system_port_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_system_port(sai_object_id_t system_port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveSystemPortRequest req;
  lemming::dataplane::sai::RemoveSystemPortResponse resp;
  grpc::ClientContext context;
  req.set_oid(system_port_id);

  grpc::Status status = system_port->RemoveSystemPort(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_system_port_attribute(sai_object_id_t system_port_id,
                                         const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetSystemPortAttributeRequest req;
  lemming::dataplane::sai::SetSystemPortAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(system_port_id);

  switch (attr->id) {
    case SAI_SYSTEM_PORT_ATTR_ADMIN_STATE:
      req.set_admin_state(attr->value.booldata);
      break;
    case SAI_SYSTEM_PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
      req.set_qos_tc_to_queue_map(attr->value.oid);
      break;
  }

  grpc::Status status =
      system_port->SetSystemPortAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_system_port_attribute(sai_object_id_t system_port_id,
                                         uint32_t attr_count,
                                         sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetSystemPortAttributeRequest req;
  lemming::dataplane::sai::GetSystemPortAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(system_port_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::SystemPortAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status =
      system_port->GetSystemPortAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SYSTEM_PORT_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_SYSTEM_PORT_ATTR_QOS_NUMBER_OF_VOQS:
        attr_list[i].value.u32 = resp.attr().qos_number_of_voqs();
        break;
      case SAI_SYSTEM_PORT_ATTR_QOS_VOQ_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().qos_voq_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_SYSTEM_PORT_ATTR_PORT:
        attr_list[i].value.oid = resp.attr().port();
        break;
      case SAI_SYSTEM_PORT_ATTR_ADMIN_STATE:
        attr_list[i].value.booldata = resp.attr().admin_state();
        break;
      case SAI_SYSTEM_PORT_ATTR_QOS_TC_TO_QUEUE_MAP:
        attr_list[i].value.oid = resp.attr().qos_tc_to_queue_map();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
