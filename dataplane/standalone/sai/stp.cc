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

#include "dataplane/standalone/sai/stp.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/stp.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_stp_api_t l_stp = {
    .create_stp = l_create_stp,
    .remove_stp = l_remove_stp,
    .set_stp_attribute = l_set_stp_attribute,
    .get_stp_attribute = l_get_stp_attribute,
    .create_stp_port = l_create_stp_port,
    .remove_stp_port = l_remove_stp_port,
    .set_stp_port_attribute = l_set_stp_port_attribute,
    .get_stp_port_attribute = l_get_stp_port_attribute,
    .create_stp_ports = l_create_stp_ports,
    .remove_stp_ports = l_remove_stp_ports,
};

lemming::dataplane::sai::CreateStpRequest convert_create_stp(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateStpRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {}
  }
  return msg;
}

lemming::dataplane::sai::CreateStpPortRequest convert_create_stp_port(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateStpPortRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_STP_PORT_ATTR_STP:
        msg.set_stp(attr_list[i].value.oid);
        break;
      case SAI_STP_PORT_ATTR_BRIDGE_PORT:
        msg.set_bridge_port(attr_list[i].value.oid);
        break;
      case SAI_STP_PORT_ATTR_STATE:
        msg.set_state(
            convert_sai_stp_port_state_t_to_proto(attr_list[i].value.s32));
        break;
    }
  }
  return msg;
}

sai_status_t l_create_stp(sai_object_id_t *stp_id, sai_object_id_t switch_id,
                          uint32_t attr_count,
                          const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateStpRequest req =
      convert_create_stp(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateStpResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = stp->CreateStp(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (stp_id) {
    *stp_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_stp(sai_object_id_t stp_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveStpRequest req;
  lemming::dataplane::sai::RemoveStpResponse resp;
  grpc::ClientContext context;
  req.set_oid(stp_id);

  grpc::Status status = stp->RemoveStp(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_stp_attribute(sai_object_id_t stp_id,
                                 const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_stp_attribute(sai_object_id_t stp_id, uint32_t attr_count,
                                 sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetStpAttributeRequest req;
  lemming::dataplane::sai::GetStpAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(stp_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_stp_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = stp->GetStpAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_STP_ATTR_BRIDGE_ID:
        attr_list[i].value.oid = resp.attr().bridge_id();
        break;
      case SAI_STP_ATTR_PORT_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().port_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_stp_port(sai_object_id_t *stp_port_id,
                               sai_object_id_t switch_id, uint32_t attr_count,
                               const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateStpPortRequest req =
      convert_create_stp_port(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateStpPortResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = stp->CreateStpPort(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (stp_port_id) {
    *stp_port_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_stp_port(sai_object_id_t stp_port_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveStpPortRequest req;
  lemming::dataplane::sai::RemoveStpPortResponse resp;
  grpc::ClientContext context;
  req.set_oid(stp_port_id);

  grpc::Status status = stp->RemoveStpPort(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_stp_port_attribute(sai_object_id_t stp_port_id,
                                      const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetStpPortAttributeRequest req;
  lemming::dataplane::sai::SetStpPortAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(stp_port_id);

  switch (attr->id) {
    case SAI_STP_PORT_ATTR_STATE:
      req.set_state(convert_sai_stp_port_state_t_to_proto(attr->value.s32));
      break;
  }

  grpc::Status status = stp->SetStpPortAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_stp_port_attribute(sai_object_id_t stp_port_id,
                                      uint32_t attr_count,
                                      sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetStpPortAttributeRequest req;
  lemming::dataplane::sai::GetStpPortAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(stp_port_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_stp_port_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = stp->GetStpPortAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_STP_PORT_ATTR_STP:
        attr_list[i].value.oid = resp.attr().stp();
        break;
      case SAI_STP_PORT_ATTR_BRIDGE_PORT:
        attr_list[i].value.oid = resp.attr().bridge_port();
        break;
      case SAI_STP_PORT_ATTR_STATE:
        attr_list[i].value.s32 =
            convert_sai_stp_port_state_t_to_sai(resp.attr().state());
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_stp_ports(sai_object_id_t switch_id,
                                uint32_t object_count,
                                const uint32_t *attr_count,
                                const sai_attribute_t **attr_list,
                                sai_bulk_op_error_mode_t mode,
                                sai_object_id_t *object_id,
                                sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateStpPortsRequest req;
  lemming::dataplane::sai::CreateStpPortsResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    auto r = convert_create_stp_port(switch_id, attr_count[i], attr_list[i]);
    *req.add_reqs() = r;
  }

  grpc::Status status = stp->CreateStpPorts(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (object_count != resp.resps().size()) {
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < object_count; i++) {
    object_id[i] = resp.resps(i).oid();
    object_statuses[i] = SAI_STATUS_SUCCESS;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_stp_ports(uint32_t object_count,
                                const sai_object_id_t *object_id,
                                sai_bulk_op_error_mode_t mode,
                                sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveStpPortsRequest req;
  lemming::dataplane::sai::RemoveStpPortsResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    req.add_reqs()->set_oid(object_id[i]);
  }

  grpc::Status status = stp->RemoveStpPorts(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (object_count != resp.resps().size()) {
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < object_count; i++) {
    object_statuses[i] = SAI_STATUS_SUCCESS;
  }

  return SAI_STATUS_SUCCESS;
}
