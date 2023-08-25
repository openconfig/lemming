

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

#include "dataplane/standalone/sai/tam.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/tam.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_tam_api_t l_tam = {
    .create_tam = l_create_tam,
    .remove_tam = l_remove_tam,
    .set_tam_attribute = l_set_tam_attribute,
    .get_tam_attribute = l_get_tam_attribute,
    .create_tam_math_func = l_create_tam_math_func,
    .remove_tam_math_func = l_remove_tam_math_func,
    .set_tam_math_func_attribute = l_set_tam_math_func_attribute,
    .get_tam_math_func_attribute = l_get_tam_math_func_attribute,
    .create_tam_report = l_create_tam_report,
    .remove_tam_report = l_remove_tam_report,
    .set_tam_report_attribute = l_set_tam_report_attribute,
    .get_tam_report_attribute = l_get_tam_report_attribute,
    .create_tam_event_threshold = l_create_tam_event_threshold,
    .remove_tam_event_threshold = l_remove_tam_event_threshold,
    .set_tam_event_threshold_attribute = l_set_tam_event_threshold_attribute,
    .get_tam_event_threshold_attribute = l_get_tam_event_threshold_attribute,
    .create_tam_int = l_create_tam_int,
    .remove_tam_int = l_remove_tam_int,
    .set_tam_int_attribute = l_set_tam_int_attribute,
    .get_tam_int_attribute = l_get_tam_int_attribute,
    .create_tam_tel_type = l_create_tam_tel_type,
    .remove_tam_tel_type = l_remove_tam_tel_type,
    .set_tam_tel_type_attribute = l_set_tam_tel_type_attribute,
    .get_tam_tel_type_attribute = l_get_tam_tel_type_attribute,
    .create_tam_transport = l_create_tam_transport,
    .remove_tam_transport = l_remove_tam_transport,
    .set_tam_transport_attribute = l_set_tam_transport_attribute,
    .get_tam_transport_attribute = l_get_tam_transport_attribute,
    .create_tam_telemetry = l_create_tam_telemetry,
    .remove_tam_telemetry = l_remove_tam_telemetry,
    .set_tam_telemetry_attribute = l_set_tam_telemetry_attribute,
    .get_tam_telemetry_attribute = l_get_tam_telemetry_attribute,
    .create_tam_collector = l_create_tam_collector,
    .remove_tam_collector = l_remove_tam_collector,
    .set_tam_collector_attribute = l_set_tam_collector_attribute,
    .get_tam_collector_attribute = l_get_tam_collector_attribute,
    .create_tam_event_action = l_create_tam_event_action,
    .remove_tam_event_action = l_remove_tam_event_action,
    .set_tam_event_action_attribute = l_set_tam_event_action_attribute,
    .get_tam_event_action_attribute = l_get_tam_event_action_attribute,
    .create_tam_event = l_create_tam_event,
    .remove_tam_event = l_remove_tam_event,
    .set_tam_event_attribute = l_set_tam_event_attribute,
    .get_tam_event_attribute = l_get_tam_event_attribute,
};

sai_status_t l_create_tam(sai_object_id_t *tam_id, sai_object_id_t switch_id,
                          uint32_t attr_count,
                          const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamRequest req;
  lemming::dataplane::sai::CreateTamResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_ATTR_TELEMETRY_OBJECTS_LIST:
        req.mutable_telemetry_objects_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TAM_ATTR_EVENT_OBJECTS_LIST:
        req.mutable_event_objects_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TAM_ATTR_INT_OBJECTS_LIST:
        req.mutable_int_objects_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
    }
  }
  grpc::Status status = tam->CreateTam(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM, tam_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_tam(sai_object_id_t tam_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM, tam_id);
}

sai_status_t l_set_tam_attribute(sai_object_id_t tam_id,
                                 const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM, tam_id, attr);
}

sai_status_t l_get_tam_attribute(sai_object_id_t tam_id, uint32_t attr_count,
                                 sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM, tam_id, attr_count,
                                   attr_list);
}

sai_status_t l_create_tam_math_func(sai_object_id_t *tam_math_func_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamMathFuncRequest req;
  lemming::dataplane::sai::CreateTamMathFuncResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_MATH_FUNC_ATTR_TAM_TEL_MATH_FUNC_TYPE:
        req.set_tam_tel_math_func_type(
            static_cast<lemming::dataplane::sai::TamTelMathFuncType>(
                attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = tam->CreateTamMathFunc(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_math_func_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM_MATH_FUNC, tam_math_func_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_math_func(sai_object_id_t tam_math_func_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM_MATH_FUNC, tam_math_func_id);
}

sai_status_t l_set_tam_math_func_attribute(sai_object_id_t tam_math_func_id,
                                           const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM_MATH_FUNC,
                                   tam_math_func_id, attr);
}

sai_status_t l_get_tam_math_func_attribute(sai_object_id_t tam_math_func_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM_MATH_FUNC,
                                   tam_math_func_id, attr_count, attr_list);
}

sai_status_t l_create_tam_report(sai_object_id_t *tam_report_id,
                                 sai_object_id_t switch_id, uint32_t attr_count,
                                 const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamReportRequest req;
  lemming::dataplane::sai::CreateTamReportResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_REPORT_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::TamReportType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_TAM_REPORT_ATTR_HISTOGRAM_NUMBER_OF_BINS:
        req.set_histogram_number_of_bins(attr_list[i].value.u32);
        break;
      case SAI_TAM_REPORT_ATTR_HISTOGRAM_BIN_BOUNDARY:
        req.mutable_histogram_bin_boundary()->Add(
            attr_list[i].value.u32list.list,
            attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
        break;
      case SAI_TAM_REPORT_ATTR_QUOTA:
        req.set_quota(attr_list[i].value.u32);
        break;
      case SAI_TAM_REPORT_ATTR_REPORT_MODE:
        req.set_report_mode(static_cast<lemming::dataplane::sai::TamReportMode>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_TAM_REPORT_ATTR_REPORT_INTERVAL:
        req.set_report_interval(attr_list[i].value.u32);
        break;
      case SAI_TAM_REPORT_ATTR_ENTERPRISE_NUMBER:
        req.set_enterprise_number(attr_list[i].value.u32);
        break;
    }
  }
  grpc::Status status = tam->CreateTamReport(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_report_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM_REPORT, tam_report_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_report(sai_object_id_t tam_report_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM_REPORT, tam_report_id);
}

sai_status_t l_set_tam_report_attribute(sai_object_id_t tam_report_id,
                                        const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM_REPORT, tam_report_id,
                                   attr);
}

sai_status_t l_get_tam_report_attribute(sai_object_id_t tam_report_id,
                                        uint32_t attr_count,
                                        sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM_REPORT, tam_report_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_tam_event_threshold(
    sai_object_id_t *tam_event_threshold_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamEventThresholdRequest req;
  lemming::dataplane::sai::CreateTamEventThresholdResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_EVENT_THRESHOLD_ATTR_HIGH_WATERMARK:
        req.set_high_watermark(attr_list[i].value.u32);
        break;
      case SAI_TAM_EVENT_THRESHOLD_ATTR_LOW_WATERMARK:
        req.set_low_watermark(attr_list[i].value.u32);
        break;
      case SAI_TAM_EVENT_THRESHOLD_ATTR_LATENCY:
        req.set_latency(attr_list[i].value.u32);
        break;
      case SAI_TAM_EVENT_THRESHOLD_ATTR_RATE:
        req.set_rate(attr_list[i].value.u32);
        break;
      case SAI_TAM_EVENT_THRESHOLD_ATTR_ABS_VALUE:
        req.set_abs_value(attr_list[i].value.u32);
        break;
      case SAI_TAM_EVENT_THRESHOLD_ATTR_UNIT:
        req.set_unit(
            static_cast<lemming::dataplane::sai::TamEventThresholdUnit>(
                attr_list[i].value.s32 + 1));
        break;
    }
  }
  grpc::Status status = tam->CreateTamEventThreshold(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_event_threshold_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM_EVENT_THRESHOLD,
                            tam_event_threshold_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_tam_event_threshold(
    sai_object_id_t tam_event_threshold_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM_EVENT_THRESHOLD,
                            tam_event_threshold_id);
}

sai_status_t l_set_tam_event_threshold_attribute(
    sai_object_id_t tam_event_threshold_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM_EVENT_THRESHOLD,
                                   tam_event_threshold_id, attr);
}

sai_status_t l_get_tam_event_threshold_attribute(
    sai_object_id_t tam_event_threshold_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM_EVENT_THRESHOLD,
                                   tam_event_threshold_id, attr_count,
                                   attr_list);
}

sai_status_t l_create_tam_int(sai_object_id_t *tam_int_id,
                              sai_object_id_t switch_id, uint32_t attr_count,
                              const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamIntRequest req;
  lemming::dataplane::sai::CreateTamIntResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_INT_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::TamIntType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_TAM_INT_ATTR_DEVICE_ID:
        req.set_device_id(attr_list[i].value.u32);
        break;
      case SAI_TAM_INT_ATTR_IOAM_TRACE_TYPE:
        req.set_ioam_trace_type(attr_list[i].value.u32);
        break;
      case SAI_TAM_INT_ATTR_INT_PRESENCE_TYPE:
        req.set_int_presence_type(
            static_cast<lemming::dataplane::sai::TamIntPresenceType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TAM_INT_ATTR_INT_PRESENCE_PB1:
        req.set_int_presence_pb1(attr_list[i].value.u32);
        break;
      case SAI_TAM_INT_ATTR_INT_PRESENCE_PB2:
        req.set_int_presence_pb2(attr_list[i].value.u32);
        break;
      case SAI_TAM_INT_ATTR_INT_PRESENCE_DSCP_VALUE:
        req.set_int_presence_dscp_value(attr_list[i].value.u8);
        break;
      case SAI_TAM_INT_ATTR_INLINE:
        req.set_inline_(attr_list[i].value.booldata);
        break;
      case SAI_TAM_INT_ATTR_INT_PRESENCE_L3_PROTOCOL:
        req.set_int_presence_l3_protocol(attr_list[i].value.u8);
        break;
      case SAI_TAM_INT_ATTR_TRACE_VECTOR:
        req.set_trace_vector(attr_list[i].value.u16);
        break;
      case SAI_TAM_INT_ATTR_ACTION_VECTOR:
        req.set_action_vector(attr_list[i].value.u16);
        break;
      case SAI_TAM_INT_ATTR_P4_INT_INSTRUCTION_BITMAP:
        req.set_p4_int_instruction_bitmap(attr_list[i].value.u16);
        break;
      case SAI_TAM_INT_ATTR_METADATA_FRAGMENT_ENABLE:
        req.set_metadata_fragment_enable(attr_list[i].value.booldata);
        break;
      case SAI_TAM_INT_ATTR_METADATA_CHECKSUM_ENABLE:
        req.set_metadata_checksum_enable(attr_list[i].value.booldata);
        break;
      case SAI_TAM_INT_ATTR_REPORT_ALL_PACKETS:
        req.set_report_all_packets(attr_list[i].value.booldata);
        break;
      case SAI_TAM_INT_ATTR_FLOW_LIVENESS_PERIOD:
        req.set_flow_liveness_period(attr_list[i].value.u16);
        break;
      case SAI_TAM_INT_ATTR_LATENCY_SENSITIVITY:
        req.set_latency_sensitivity(attr_list[i].value.u8);
        break;
      case SAI_TAM_INT_ATTR_ACL_GROUP:
        req.set_acl_group(attr_list[i].value.oid);
        break;
      case SAI_TAM_INT_ATTR_MAX_HOP_COUNT:
        req.set_max_hop_count(attr_list[i].value.u8);
        break;
      case SAI_TAM_INT_ATTR_MAX_LENGTH:
        req.set_max_length(attr_list[i].value.u8);
        break;
      case SAI_TAM_INT_ATTR_NAME_SPACE_ID:
        req.set_name_space_id(attr_list[i].value.u8);
        break;
      case SAI_TAM_INT_ATTR_NAME_SPACE_ID_GLOBAL:
        req.set_name_space_id_global(attr_list[i].value.booldata);
        break;
      case SAI_TAM_INT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
        req.set_ingress_samplepacket_enable(attr_list[i].value.oid);
        break;
      case SAI_TAM_INT_ATTR_COLLECTOR_LIST:
        req.mutable_collector_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TAM_INT_ATTR_MATH_FUNC:
        req.set_math_func(attr_list[i].value.oid);
        break;
      case SAI_TAM_INT_ATTR_REPORT_ID:
        req.set_report_id(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = tam->CreateTamInt(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_int_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM_INT, tam_int_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_tam_int(sai_object_id_t tam_int_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM_INT, tam_int_id);
}

sai_status_t l_set_tam_int_attribute(sai_object_id_t tam_int_id,
                                     const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM_INT, tam_int_id, attr);
}

sai_status_t l_get_tam_int_attribute(sai_object_id_t tam_int_id,
                                     uint32_t attr_count,
                                     sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM_INT, tam_int_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_tam_tel_type(sai_object_id_t *tam_tel_type_id,
                                   sai_object_id_t switch_id,
                                   uint32_t attr_count,
                                   const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamTelTypeRequest req;
  lemming::dataplane::sai::CreateTamTelTypeResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_TEL_TYPE_ATTR_TAM_TELEMETRY_TYPE:
        req.set_tam_telemetry_type(
            static_cast<lemming::dataplane::sai::TamTelemetryType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TAM_TEL_TYPE_ATTR_INT_SWITCH_IDENTIFIER:
        req.set_int_switch_identifier(attr_list[i].value.u32);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS:
        req.set_switch_enable_port_stats(attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_INGRESS:
        req.set_switch_enable_port_stats_ingress(attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_EGRESS:
        req.set_switch_enable_port_stats_egress(attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_VIRTUAL_QUEUE_STATS:
        req.set_switch_enable_virtual_queue_stats(attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_OUTPUT_QUEUE_STATS:
        req.set_switch_enable_output_queue_stats(attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_MMU_STATS:
        req.set_switch_enable_mmu_stats(attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FABRIC_STATS:
        req.set_switch_enable_fabric_stats(attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FILTER_STATS:
        req.set_switch_enable_filter_stats(attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_RESOURCE_UTILIZATION_STATS:
        req.set_switch_enable_resource_utilization_stats(
            attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_FABRIC_Q:
        req.set_fabric_q(attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_NE_ENABLE:
        req.set_ne_enable(attr_list[i].value.booldata);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_DSCP_VALUE:
        req.set_dscp_value(attr_list[i].value.u8);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_MATH_FUNC:
        req.set_math_func(attr_list[i].value.oid);
        break;
      case SAI_TAM_TEL_TYPE_ATTR_REPORT_ID:
        req.set_report_id(attr_list[i].value.oid);
        break;
    }
  }
  grpc::Status status = tam->CreateTamTelType(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_tel_type_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM_TEL_TYPE, tam_tel_type_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_tel_type(sai_object_id_t tam_tel_type_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM_TEL_TYPE, tam_tel_type_id);
}

sai_status_t l_set_tam_tel_type_attribute(sai_object_id_t tam_tel_type_id,
                                          const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM_TEL_TYPE,
                                   tam_tel_type_id, attr);
}

sai_status_t l_get_tam_tel_type_attribute(sai_object_id_t tam_tel_type_id,
                                          uint32_t attr_count,
                                          sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM_TEL_TYPE,
                                   tam_tel_type_id, attr_count, attr_list);
}

sai_status_t l_create_tam_transport(sai_object_id_t *tam_transport_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamTransportRequest req;
  lemming::dataplane::sai::CreateTamTransportResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_TRANSPORT_ATTR_TRANSPORT_TYPE:
        req.set_transport_type(
            static_cast<lemming::dataplane::sai::TamTransportType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TAM_TRANSPORT_ATTR_SRC_PORT:
        req.set_src_port(attr_list[i].value.u32);
        break;
      case SAI_TAM_TRANSPORT_ATTR_DST_PORT:
        req.set_dst_port(attr_list[i].value.u32);
        break;
      case SAI_TAM_TRANSPORT_ATTR_TRANSPORT_AUTH_TYPE:
        req.set_transport_auth_type(
            static_cast<lemming::dataplane::sai::TamTransportAuthType>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TAM_TRANSPORT_ATTR_MTU:
        req.set_mtu(attr_list[i].value.u32);
        break;
    }
  }
  grpc::Status status = tam->CreateTamTransport(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_transport_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM_TRANSPORT, tam_transport_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_transport(sai_object_id_t tam_transport_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM_TRANSPORT, tam_transport_id);
}

sai_status_t l_set_tam_transport_attribute(sai_object_id_t tam_transport_id,
                                           const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM_TRANSPORT,
                                   tam_transport_id, attr);
}

sai_status_t l_get_tam_transport_attribute(sai_object_id_t tam_transport_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM_TRANSPORT,
                                   tam_transport_id, attr_count, attr_list);
}

sai_status_t l_create_tam_telemetry(sai_object_id_t *tam_telemetry_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamTelemetryRequest req;
  lemming::dataplane::sai::CreateTamTelemetryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_TELEMETRY_ATTR_TAM_TYPE_LIST:
        req.mutable_tam_type_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TAM_TELEMETRY_ATTR_COLLECTOR_LIST:
        req.mutable_collector_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TAM_TELEMETRY_ATTR_TAM_REPORTING_UNIT:
        req.set_tam_reporting_unit(
            static_cast<lemming::dataplane::sai::TamReportingUnit>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_TAM_TELEMETRY_ATTR_REPORTING_INTERVAL:
        req.set_reporting_interval(attr_list[i].value.u32);
        break;
    }
  }
  grpc::Status status = tam->CreateTamTelemetry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_telemetry_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM_TELEMETRY, tam_telemetry_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_telemetry(sai_object_id_t tam_telemetry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM_TELEMETRY, tam_telemetry_id);
}

sai_status_t l_set_tam_telemetry_attribute(sai_object_id_t tam_telemetry_id,
                                           const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM_TELEMETRY,
                                   tam_telemetry_id, attr);
}

sai_status_t l_get_tam_telemetry_attribute(sai_object_id_t tam_telemetry_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM_TELEMETRY,
                                   tam_telemetry_id, attr_count, attr_list);
}

sai_status_t l_create_tam_collector(sai_object_id_t *tam_collector_id,
                                    sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamCollectorRequest req;
  lemming::dataplane::sai::CreateTamCollectorResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_COLLECTOR_ATTR_SRC_IP:
        req.set_src_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TAM_COLLECTOR_ATTR_DST_IP:
        req.set_dst_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
      case SAI_TAM_COLLECTOR_ATTR_LOCALHOST:
        req.set_localhost(attr_list[i].value.booldata);
        break;
      case SAI_TAM_COLLECTOR_ATTR_VIRTUAL_ROUTER_ID:
        req.set_virtual_router_id(attr_list[i].value.oid);
        break;
      case SAI_TAM_COLLECTOR_ATTR_TRUNCATE_SIZE:
        req.set_truncate_size(attr_list[i].value.u16);
        break;
      case SAI_TAM_COLLECTOR_ATTR_TRANSPORT:
        req.set_transport(attr_list[i].value.oid);
        break;
      case SAI_TAM_COLLECTOR_ATTR_DSCP_VALUE:
        req.set_dscp_value(attr_list[i].value.u8);
        break;
    }
  }
  grpc::Status status = tam->CreateTamCollector(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_collector_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM_COLLECTOR, tam_collector_id,
                            switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_collector(sai_object_id_t tam_collector_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM_COLLECTOR, tam_collector_id);
}

sai_status_t l_set_tam_collector_attribute(sai_object_id_t tam_collector_id,
                                           const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM_COLLECTOR,
                                   tam_collector_id, attr);
}

sai_status_t l_get_tam_collector_attribute(sai_object_id_t tam_collector_id,
                                           uint32_t attr_count,
                                           sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM_COLLECTOR,
                                   tam_collector_id, attr_count, attr_list);
}

sai_status_t l_create_tam_event_action(sai_object_id_t *tam_event_action_id,
                                       sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamEventActionRequest req;
  lemming::dataplane::sai::CreateTamEventActionResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_EVENT_ACTION_ATTR_REPORT_TYPE:
        req.set_report_type(attr_list[i].value.oid);
        break;
      case SAI_TAM_EVENT_ACTION_ATTR_QOS_ACTION_TYPE:
        req.set_qos_action_type(attr_list[i].value.u32);
        break;
    }
  }
  grpc::Status status = tam->CreateTamEventAction(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_event_action_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM_EVENT_ACTION,
                            tam_event_action_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_tam_event_action(sai_object_id_t tam_event_action_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM_EVENT_ACTION,
                            tam_event_action_id);
}

sai_status_t l_set_tam_event_action_attribute(
    sai_object_id_t tam_event_action_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM_EVENT_ACTION,
                                   tam_event_action_id, attr);
}

sai_status_t l_get_tam_event_action_attribute(
    sai_object_id_t tam_event_action_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM_EVENT_ACTION,
                                   tam_event_action_id, attr_count, attr_list);
}

sai_status_t l_create_tam_event(sai_object_id_t *tam_event_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTamEventRequest req;
  lemming::dataplane::sai::CreateTamEventResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TAM_EVENT_ATTR_TYPE:
        req.set_type(static_cast<lemming::dataplane::sai::TamEventType>(
            attr_list[i].value.s32 + 1));
        break;
      case SAI_TAM_EVENT_ATTR_ACTION_LIST:
        req.mutable_action_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TAM_EVENT_ATTR_COLLECTOR_LIST:
        req.mutable_collector_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_TAM_EVENT_ATTR_THRESHOLD:
        req.set_threshold(attr_list[i].value.oid);
        break;
      case SAI_TAM_EVENT_ATTR_DSCP_VALUE:
        req.set_dscp_value(attr_list[i].value.u8);
        break;
    }
  }
  grpc::Status status = tam->CreateTamEvent(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *tam_event_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_TAM_EVENT, tam_event_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_tam_event(sai_object_id_t tam_event_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_TAM_EVENT, tam_event_id);
}

sai_status_t l_set_tam_event_attribute(sai_object_id_t tam_event_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_TAM_EVENT, tam_event_id,
                                   attr);
}

sai_status_t l_get_tam_event_attribute(sai_object_id_t tam_event_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_TAM_EVENT, tam_event_id,
                                   attr_count, attr_list);
}
