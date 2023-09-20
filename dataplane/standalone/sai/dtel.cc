

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

#include "dataplane/standalone/sai/dtel.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/dtel.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_dtel_api_t l_dtel = {
    .create_dtel = l_create_dtel,
    .remove_dtel = l_remove_dtel,
    .set_dtel_attribute = l_set_dtel_attribute,
    .get_dtel_attribute = l_get_dtel_attribute,
    .create_dtel_queue_report = l_create_dtel_queue_report,
    .remove_dtel_queue_report = l_remove_dtel_queue_report,
    .set_dtel_queue_report_attribute = l_set_dtel_queue_report_attribute,
    .get_dtel_queue_report_attribute = l_get_dtel_queue_report_attribute,
    .create_dtel_int_session = l_create_dtel_int_session,
    .remove_dtel_int_session = l_remove_dtel_int_session,
    .set_dtel_int_session_attribute = l_set_dtel_int_session_attribute,
    .get_dtel_int_session_attribute = l_get_dtel_int_session_attribute,
    .create_dtel_report_session = l_create_dtel_report_session,
    .remove_dtel_report_session = l_remove_dtel_report_session,
    .set_dtel_report_session_attribute = l_set_dtel_report_session_attribute,
    .get_dtel_report_session_attribute = l_get_dtel_report_session_attribute,
    .create_dtel_event = l_create_dtel_event,
    .remove_dtel_event = l_remove_dtel_event,
    .set_dtel_event_attribute = l_set_dtel_event_attribute,
    .get_dtel_event_attribute = l_get_dtel_event_attribute,
};

sai_status_t l_create_dtel(sai_object_id_t *dtel_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateDtelRequest req;
  lemming::dataplane::sai::CreateDtelResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DTEL_ATTR_INT_ENDPOINT_ENABLE:
        req.set_int_endpoint_enable(attr_list[i].value.booldata);
        break;
      case SAI_DTEL_ATTR_INT_TRANSIT_ENABLE:
        req.set_int_transit_enable(attr_list[i].value.booldata);
        break;
      case SAI_DTEL_ATTR_POSTCARD_ENABLE:
        req.set_postcard_enable(attr_list[i].value.booldata);
        break;
      case SAI_DTEL_ATTR_DROP_REPORT_ENABLE:
        req.set_drop_report_enable(attr_list[i].value.booldata);
        break;
      case SAI_DTEL_ATTR_QUEUE_REPORT_ENABLE:
        req.set_queue_report_enable(attr_list[i].value.booldata);
        break;
      case SAI_DTEL_ATTR_SWITCH_ID:
        req.set_switch_id(attr_list[i].value.u32);
        break;
      case SAI_DTEL_ATTR_FLOW_STATE_CLEAR_CYCLE:
        req.set_flow_state_clear_cycle(attr_list[i].value.u16);
        break;
      case SAI_DTEL_ATTR_LATENCY_SENSITIVITY:
        req.set_latency_sensitivity(attr_list[i].value.u8);
        break;
      case SAI_DTEL_ATTR_SINK_PORT_LIST:
        req.mutable_sink_port_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
    }
  }
  grpc::Status status = dtel->CreateDtel(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *dtel_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_dtel(sai_object_id_t dtel_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveDtelRequest req;
  lemming::dataplane::sai::RemoveDtelResponse resp;
  grpc::ClientContext context;
  req.set_oid(dtel_id);

  grpc::Status status = dtel->RemoveDtel(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_dtel_attribute(sai_object_id_t dtel_id,
                                  const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetDtelAttributeRequest req;
  lemming::dataplane::sai::SetDtelAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(dtel_id);

  switch (attr->id) {
    case SAI_DTEL_ATTR_INT_ENDPOINT_ENABLE:
      req.set_int_endpoint_enable(attr->value.booldata);
      break;
    case SAI_DTEL_ATTR_INT_TRANSIT_ENABLE:
      req.set_int_transit_enable(attr->value.booldata);
      break;
    case SAI_DTEL_ATTR_POSTCARD_ENABLE:
      req.set_postcard_enable(attr->value.booldata);
      break;
    case SAI_DTEL_ATTR_DROP_REPORT_ENABLE:
      req.set_drop_report_enable(attr->value.booldata);
      break;
    case SAI_DTEL_ATTR_QUEUE_REPORT_ENABLE:
      req.set_queue_report_enable(attr->value.booldata);
      break;
    case SAI_DTEL_ATTR_SWITCH_ID:
      req.set_switch_id(attr->value.u32);
      break;
    case SAI_DTEL_ATTR_FLOW_STATE_CLEAR_CYCLE:
      req.set_flow_state_clear_cycle(attr->value.u16);
      break;
    case SAI_DTEL_ATTR_LATENCY_SENSITIVITY:
      req.set_latency_sensitivity(attr->value.u8);
      break;
    case SAI_DTEL_ATTR_SINK_PORT_LIST:
      req.mutable_sink_port_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
  }

  grpc::Status status = dtel->SetDtelAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_dtel_attribute(sai_object_id_t dtel_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetDtelAttributeRequest req;
  lemming::dataplane::sai::GetDtelAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(dtel_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::DtelAttr>(attr_list[i].id + 1));
  }
  grpc::Status status = dtel->GetDtelAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DTEL_ATTR_INT_ENDPOINT_ENABLE:
        attr_list[i].value.booldata = resp.attr().int_endpoint_enable();
        break;
      case SAI_DTEL_ATTR_INT_TRANSIT_ENABLE:
        attr_list[i].value.booldata = resp.attr().int_transit_enable();
        break;
      case SAI_DTEL_ATTR_POSTCARD_ENABLE:
        attr_list[i].value.booldata = resp.attr().postcard_enable();
        break;
      case SAI_DTEL_ATTR_DROP_REPORT_ENABLE:
        attr_list[i].value.booldata = resp.attr().drop_report_enable();
        break;
      case SAI_DTEL_ATTR_QUEUE_REPORT_ENABLE:
        attr_list[i].value.booldata = resp.attr().queue_report_enable();
        break;
      case SAI_DTEL_ATTR_SWITCH_ID:
        attr_list[i].value.u32 = resp.attr().switch_id();
        break;
      case SAI_DTEL_ATTR_FLOW_STATE_CLEAR_CYCLE:
        attr_list[i].value.u16 = resp.attr().flow_state_clear_cycle();
        break;
      case SAI_DTEL_ATTR_LATENCY_SENSITIVITY:
        attr_list[i].value.u8 = resp.attr().latency_sensitivity();
        break;
      case SAI_DTEL_ATTR_SINK_PORT_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().sink_port_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_dtel_queue_report(sai_object_id_t *dtel_queue_report_id,
                                        sai_object_id_t switch_id,
                                        uint32_t attr_count,
                                        const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateDtelQueueReportRequest req;
  lemming::dataplane::sai::CreateDtelQueueReportResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DTEL_QUEUE_REPORT_ATTR_QUEUE_ID:
        req.set_queue_id(attr_list[i].value.oid);
        break;
      case SAI_DTEL_QUEUE_REPORT_ATTR_DEPTH_THRESHOLD:
        req.set_depth_threshold(attr_list[i].value.u32);
        break;
      case SAI_DTEL_QUEUE_REPORT_ATTR_LATENCY_THRESHOLD:
        req.set_latency_threshold(attr_list[i].value.u32);
        break;
      case SAI_DTEL_QUEUE_REPORT_ATTR_BREACH_QUOTA:
        req.set_breach_quota(attr_list[i].value.u32);
        break;
      case SAI_DTEL_QUEUE_REPORT_ATTR_TAIL_DROP:
        req.set_tail_drop(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status = dtel->CreateDtelQueueReport(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *dtel_queue_report_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_dtel_queue_report(sai_object_id_t dtel_queue_report_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveDtelQueueReportRequest req;
  lemming::dataplane::sai::RemoveDtelQueueReportResponse resp;
  grpc::ClientContext context;
  req.set_oid(dtel_queue_report_id);

  grpc::Status status = dtel->RemoveDtelQueueReport(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_dtel_queue_report_attribute(
    sai_object_id_t dtel_queue_report_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetDtelQueueReportAttributeRequest req;
  lemming::dataplane::sai::SetDtelQueueReportAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(dtel_queue_report_id);

  switch (attr->id) {
    case SAI_DTEL_QUEUE_REPORT_ATTR_DEPTH_THRESHOLD:
      req.set_depth_threshold(attr->value.u32);
      break;
    case SAI_DTEL_QUEUE_REPORT_ATTR_LATENCY_THRESHOLD:
      req.set_latency_threshold(attr->value.u32);
      break;
    case SAI_DTEL_QUEUE_REPORT_ATTR_BREACH_QUOTA:
      req.set_breach_quota(attr->value.u32);
      break;
    case SAI_DTEL_QUEUE_REPORT_ATTR_TAIL_DROP:
      req.set_tail_drop(attr->value.booldata);
      break;
  }

  grpc::Status status = dtel->SetDtelQueueReportAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_dtel_queue_report_attribute(
    sai_object_id_t dtel_queue_report_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetDtelQueueReportAttributeRequest req;
  lemming::dataplane::sai::GetDtelQueueReportAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(dtel_queue_report_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::DtelQueueReportAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = dtel->GetDtelQueueReportAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DTEL_QUEUE_REPORT_ATTR_QUEUE_ID:
        attr_list[i].value.oid = resp.attr().queue_id();
        break;
      case SAI_DTEL_QUEUE_REPORT_ATTR_DEPTH_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().depth_threshold();
        break;
      case SAI_DTEL_QUEUE_REPORT_ATTR_LATENCY_THRESHOLD:
        attr_list[i].value.u32 = resp.attr().latency_threshold();
        break;
      case SAI_DTEL_QUEUE_REPORT_ATTR_BREACH_QUOTA:
        attr_list[i].value.u32 = resp.attr().breach_quota();
        break;
      case SAI_DTEL_QUEUE_REPORT_ATTR_TAIL_DROP:
        attr_list[i].value.booldata = resp.attr().tail_drop();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_dtel_int_session(sai_object_id_t *dtel_int_session_id,
                                       sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateDtelIntSessionRequest req;
  lemming::dataplane::sai::CreateDtelIntSessionResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DTEL_INT_SESSION_ATTR_MAX_HOP_COUNT:
        req.set_max_hop_count(attr_list[i].value.u8);
        break;
      case SAI_DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_ID:
        req.set_collect_switch_id(attr_list[i].value.booldata);
        break;
      case SAI_DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_PORTS:
        req.set_collect_switch_ports(attr_list[i].value.booldata);
        break;
      case SAI_DTEL_INT_SESSION_ATTR_COLLECT_INGRESS_TIMESTAMP:
        req.set_collect_ingress_timestamp(attr_list[i].value.booldata);
        break;
      case SAI_DTEL_INT_SESSION_ATTR_COLLECT_EGRESS_TIMESTAMP:
        req.set_collect_egress_timestamp(attr_list[i].value.booldata);
        break;
      case SAI_DTEL_INT_SESSION_ATTR_COLLECT_QUEUE_INFO:
        req.set_collect_queue_info(attr_list[i].value.booldata);
        break;
    }
  }
  grpc::Status status = dtel->CreateDtelIntSession(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *dtel_int_session_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_dtel_int_session(sai_object_id_t dtel_int_session_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveDtelIntSessionRequest req;
  lemming::dataplane::sai::RemoveDtelIntSessionResponse resp;
  grpc::ClientContext context;
  req.set_oid(dtel_int_session_id);

  grpc::Status status = dtel->RemoveDtelIntSession(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_dtel_int_session_attribute(
    sai_object_id_t dtel_int_session_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetDtelIntSessionAttributeRequest req;
  lemming::dataplane::sai::SetDtelIntSessionAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(dtel_int_session_id);

  switch (attr->id) {
    case SAI_DTEL_INT_SESSION_ATTR_MAX_HOP_COUNT:
      req.set_max_hop_count(attr->value.u8);
      break;
    case SAI_DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_ID:
      req.set_collect_switch_id(attr->value.booldata);
      break;
    case SAI_DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_PORTS:
      req.set_collect_switch_ports(attr->value.booldata);
      break;
    case SAI_DTEL_INT_SESSION_ATTR_COLLECT_INGRESS_TIMESTAMP:
      req.set_collect_ingress_timestamp(attr->value.booldata);
      break;
    case SAI_DTEL_INT_SESSION_ATTR_COLLECT_EGRESS_TIMESTAMP:
      req.set_collect_egress_timestamp(attr->value.booldata);
      break;
    case SAI_DTEL_INT_SESSION_ATTR_COLLECT_QUEUE_INFO:
      req.set_collect_queue_info(attr->value.booldata);
      break;
  }

  grpc::Status status = dtel->SetDtelIntSessionAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_dtel_int_session_attribute(
    sai_object_id_t dtel_int_session_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetDtelIntSessionAttributeRequest req;
  lemming::dataplane::sai::GetDtelIntSessionAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(dtel_int_session_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::DtelIntSessionAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = dtel->GetDtelIntSessionAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DTEL_INT_SESSION_ATTR_MAX_HOP_COUNT:
        attr_list[i].value.u8 = resp.attr().max_hop_count();
        break;
      case SAI_DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_ID:
        attr_list[i].value.booldata = resp.attr().collect_switch_id();
        break;
      case SAI_DTEL_INT_SESSION_ATTR_COLLECT_SWITCH_PORTS:
        attr_list[i].value.booldata = resp.attr().collect_switch_ports();
        break;
      case SAI_DTEL_INT_SESSION_ATTR_COLLECT_INGRESS_TIMESTAMP:
        attr_list[i].value.booldata = resp.attr().collect_ingress_timestamp();
        break;
      case SAI_DTEL_INT_SESSION_ATTR_COLLECT_EGRESS_TIMESTAMP:
        attr_list[i].value.booldata = resp.attr().collect_egress_timestamp();
        break;
      case SAI_DTEL_INT_SESSION_ATTR_COLLECT_QUEUE_INFO:
        attr_list[i].value.booldata = resp.attr().collect_queue_info();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_dtel_report_session(
    sai_object_id_t *dtel_report_session_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateDtelReportSessionRequest req;
  lemming::dataplane::sai::CreateDtelReportSessionResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DTEL_REPORT_SESSION_ATTR_SRC_IP:
        req.set_src_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_DTEL_REPORT_SESSION_ATTR_VIRTUAL_ROUTER_ID:
        req.set_virtual_router_id(attr_list[i].value.oid);
        break;
      case SAI_DTEL_REPORT_SESSION_ATTR_TRUNCATE_SIZE:
        req.set_truncate_size(attr_list[i].value.u16);
        break;
      case SAI_DTEL_REPORT_SESSION_ATTR_UDP_DST_PORT:
        req.set_udp_dst_port(attr_list[i].value.u16);
        break;
    }
  }
  grpc::Status status = dtel->CreateDtelReportSession(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *dtel_report_session_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_dtel_report_session(
    sai_object_id_t dtel_report_session_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveDtelReportSessionRequest req;
  lemming::dataplane::sai::RemoveDtelReportSessionResponse resp;
  grpc::ClientContext context;
  req.set_oid(dtel_report_session_id);

  grpc::Status status = dtel->RemoveDtelReportSession(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_dtel_report_session_attribute(
    sai_object_id_t dtel_report_session_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetDtelReportSessionAttributeRequest req;
  lemming::dataplane::sai::SetDtelReportSessionAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(dtel_report_session_id);

  switch (attr->id) {
    case SAI_DTEL_REPORT_SESSION_ATTR_SRC_IP:
      req.set_src_ip(convert_from_ip_address(attr->value.ipaddr));
      break;
    case SAI_DTEL_REPORT_SESSION_ATTR_VIRTUAL_ROUTER_ID:
      req.set_virtual_router_id(attr->value.oid);
      break;
    case SAI_DTEL_REPORT_SESSION_ATTR_TRUNCATE_SIZE:
      req.set_truncate_size(attr->value.u16);
      break;
    case SAI_DTEL_REPORT_SESSION_ATTR_UDP_DST_PORT:
      req.set_udp_dst_port(attr->value.u16);
      break;
  }

  grpc::Status status =
      dtel->SetDtelReportSessionAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_dtel_report_session_attribute(
    sai_object_id_t dtel_report_session_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetDtelReportSessionAttributeRequest req;
  lemming::dataplane::sai::GetDtelReportSessionAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(dtel_report_session_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        static_cast<lemming::dataplane::sai::DtelReportSessionAttr>(
            attr_list[i].id + 1));
  }
  grpc::Status status =
      dtel->GetDtelReportSessionAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DTEL_REPORT_SESSION_ATTR_SRC_IP:
        attr_list[i].value.ipaddr = convert_to_ip_address(resp.attr().src_ip());
        break;
      case SAI_DTEL_REPORT_SESSION_ATTR_VIRTUAL_ROUTER_ID:
        attr_list[i].value.oid = resp.attr().virtual_router_id();
        break;
      case SAI_DTEL_REPORT_SESSION_ATTR_TRUNCATE_SIZE:
        attr_list[i].value.u16 = resp.attr().truncate_size();
        break;
      case SAI_DTEL_REPORT_SESSION_ATTR_UDP_DST_PORT:
        attr_list[i].value.u16 = resp.attr().udp_dst_port();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_dtel_event(sai_object_id_t *dtel_event_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateDtelEventRequest req;
  lemming::dataplane::sai::CreateDtelEventResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DTEL_EVENT_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::DtelEventType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_DTEL_EVENT_ATTR_REPORT_SESSION:
        req.set_report_session(attr_list[i].value.oid);
        break;
      case SAI_DTEL_EVENT_ATTR_DSCP_VALUE:
        req.set_dscp_value(attr_list[i].value.u8);
        break;
    }
  }
  grpc::Status status = dtel->CreateDtelEvent(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *dtel_event_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_dtel_event(sai_object_id_t dtel_event_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveDtelEventRequest req;
  lemming::dataplane::sai::RemoveDtelEventResponse resp;
  grpc::ClientContext context;
  req.set_oid(dtel_event_id);

  grpc::Status status = dtel->RemoveDtelEvent(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_dtel_event_attribute(sai_object_id_t dtel_event_id,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetDtelEventAttributeRequest req;
  lemming::dataplane::sai::SetDtelEventAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(dtel_event_id);

  switch (attr->id) {
    case SAI_DTEL_EVENT_ATTR_REPORT_SESSION:
      req.set_report_session(attr->value.oid);
      break;
    case SAI_DTEL_EVENT_ATTR_DSCP_VALUE:
      req.set_dscp_value(attr->value.u8);
      break;
  }

  grpc::Status status = dtel->SetDtelEventAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_dtel_event_attribute(sai_object_id_t dtel_event_id,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetDtelEventAttributeRequest req;
  lemming::dataplane::sai::GetDtelEventAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(dtel_event_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::DtelEventAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = dtel->GetDtelEventAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_DTEL_EVENT_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
      case SAI_DTEL_EVENT_ATTR_REPORT_SESSION:
        attr_list[i].value.oid = resp.attr().report_session();
        break;
      case SAI_DTEL_EVENT_ATTR_DSCP_VALUE:
        attr_list[i].value.u8 = resp.attr().dscp_value();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
