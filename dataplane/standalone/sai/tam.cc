




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
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"
#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/tam.pb.h"

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


sai_status_t l_create_tam(sai_object_id_t *tam_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamRequest req;
	lemming::dataplane::sai::CreateTamResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_ATTR_TELEMETRY_OBJECTS_LIST:
	req.mutable_telemetry_objects_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_ATTR_EVENT_OBJECTS_LIST:
	req.mutable_event_objects_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_ATTR_INT_OBJECTS_LIST:
	req.mutable_int_objects_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
}

	}
	grpc::Status status = tam->CreateTam(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	*tam_id = resp.oid(); 

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam(sai_object_id_t tam_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamRequest req;
	lemming::dataplane::sai::RemoveTamResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_id); 
	
	grpc::Status status = tam->RemoveTam(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_attribute(sai_object_id_t tam_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamAttributeRequest req;
	lemming::dataplane::sai::SetTamAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_ATTR_TELEMETRY_OBJECTS_LIST:
	req.mutable_telemetry_objects_list()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
  case SAI_TAM_ATTR_EVENT_OBJECTS_LIST:
	req.mutable_event_objects_list()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
  case SAI_TAM_ATTR_INT_OBJECTS_LIST:
	req.mutable_int_objects_list()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
}

	grpc::Status status = tam->SetTamAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_attribute(sai_object_id_t tam_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamAttributeRequest req;
	lemming::dataplane::sai::GetTamAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_ATTR_TELEMETRY_OBJECTS_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().telemetry_objects_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_ATTR_EVENT_OBJECTS_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().event_objects_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_ATTR_INT_OBJECTS_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().int_objects_list(), &attr_list[i].value.objlist.count);
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tam_math_func(sai_object_id_t *tam_math_func_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamMathFuncRequest req;
	lemming::dataplane::sai::CreateTamMathFuncResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_MATH_FUNC_ATTR_TAM_TEL_MATH_FUNC_TYPE:
	req.set_tam_tel_math_func_type(static_cast<lemming::dataplane::sai::TamTelMathFuncType>(attr_list[i].value.s32 + 1));
	break;
}

	}
	grpc::Status status = tam->CreateTamMathFunc(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	*tam_math_func_id = resp.oid(); 

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam_math_func(sai_object_id_t tam_math_func_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamMathFuncRequest req;
	lemming::dataplane::sai::RemoveTamMathFuncResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_math_func_id); 
	
	grpc::Status status = tam->RemoveTamMathFunc(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_math_func_attribute(sai_object_id_t tam_math_func_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamMathFuncAttributeRequest req;
	lemming::dataplane::sai::SetTamMathFuncAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_math_func_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_MATH_FUNC_ATTR_TAM_TEL_MATH_FUNC_TYPE:
	req.set_tam_tel_math_func_type(static_cast<lemming::dataplane::sai::TamTelMathFuncType>(attr->value.s32 + 1));
	break;
}

	grpc::Status status = tam->SetTamMathFuncAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_math_func_attribute(sai_object_id_t tam_math_func_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamMathFuncAttributeRequest req;
	lemming::dataplane::sai::GetTamMathFuncAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_math_func_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamMathFuncAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamMathFuncAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_MATH_FUNC_ATTR_TAM_TEL_MATH_FUNC_TYPE:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().tam_tel_math_func_type() - 1);
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tam_report(sai_object_id_t *tam_report_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamReportRequest req;
	lemming::dataplane::sai::CreateTamReportResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_REPORT_ATTR_TYPE:
	req.set_type(static_cast<lemming::dataplane::sai::TamReportType>(attr_list[i].value.s32 + 1));
	break;
  case SAI_TAM_REPORT_ATTR_HISTOGRAM_NUMBER_OF_BINS:
	req.set_histogram_number_of_bins(attr_list[i].value.u32);
	break;
  case SAI_TAM_REPORT_ATTR_HISTOGRAM_BIN_BOUNDARY:
	req.mutable_histogram_bin_boundary()->Add(attr_list[i].value.u32list.list, attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
	break;
  case SAI_TAM_REPORT_ATTR_QUOTA:
	req.set_quota(attr_list[i].value.u32);
	break;
  case SAI_TAM_REPORT_ATTR_REPORT_MODE:
	req.set_report_mode(static_cast<lemming::dataplane::sai::TamReportMode>(attr_list[i].value.s32 + 1));
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

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam_report(sai_object_id_t tam_report_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamReportRequest req;
	lemming::dataplane::sai::RemoveTamReportResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_report_id); 
	
	grpc::Status status = tam->RemoveTamReport(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_report_attribute(sai_object_id_t tam_report_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamReportAttributeRequest req;
	lemming::dataplane::sai::SetTamReportAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_report_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_REPORT_ATTR_TYPE:
	req.set_type(static_cast<lemming::dataplane::sai::TamReportType>(attr->value.s32 + 1));
	break;
  case SAI_TAM_REPORT_ATTR_QUOTA:
	req.set_quota(attr->value.u32);
	break;
  case SAI_TAM_REPORT_ATTR_REPORT_INTERVAL:
	req.set_report_interval(attr->value.u32);
	break;
  case SAI_TAM_REPORT_ATTR_ENTERPRISE_NUMBER:
	req.set_enterprise_number(attr->value.u32);
	break;
}

	grpc::Status status = tam->SetTamReportAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_report_attribute(sai_object_id_t tam_report_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamReportAttributeRequest req;
	lemming::dataplane::sai::GetTamReportAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_report_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamReportAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamReportAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_REPORT_ATTR_TYPE:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().type() - 1);
	break;
  case SAI_TAM_REPORT_ATTR_HISTOGRAM_NUMBER_OF_BINS:
	 attr_list[i].value.u32 =   resp.attr().histogram_number_of_bins();
	break;
  case SAI_TAM_REPORT_ATTR_HISTOGRAM_BIN_BOUNDARY:
	copy_list(attr_list[i].value.u32list.list, resp.attr().histogram_bin_boundary(), &attr_list[i].value.u32list.count);
	break;
  case SAI_TAM_REPORT_ATTR_QUOTA:
	 attr_list[i].value.u32 =   resp.attr().quota();
	break;
  case SAI_TAM_REPORT_ATTR_REPORT_MODE:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().report_mode() - 1);
	break;
  case SAI_TAM_REPORT_ATTR_REPORT_INTERVAL:
	 attr_list[i].value.u32 =   resp.attr().report_interval();
	break;
  case SAI_TAM_REPORT_ATTR_ENTERPRISE_NUMBER:
	 attr_list[i].value.u32 =   resp.attr().enterprise_number();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tam_event_threshold(sai_object_id_t *tam_event_threshold_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamEventThresholdRequest req;
	lemming::dataplane::sai::CreateTamEventThresholdResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

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
	req.set_unit(static_cast<lemming::dataplane::sai::TamEventThresholdUnit>(attr_list[i].value.s32 + 1));
	break;
}

	}
	grpc::Status status = tam->CreateTamEventThreshold(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	*tam_event_threshold_id = resp.oid(); 

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam_event_threshold(sai_object_id_t tam_event_threshold_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamEventThresholdRequest req;
	lemming::dataplane::sai::RemoveTamEventThresholdResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_event_threshold_id); 
	
	grpc::Status status = tam->RemoveTamEventThreshold(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_event_threshold_attribute(sai_object_id_t tam_event_threshold_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamEventThresholdAttributeRequest req;
	lemming::dataplane::sai::SetTamEventThresholdAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_event_threshold_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_EVENT_THRESHOLD_ATTR_HIGH_WATERMARK:
	req.set_high_watermark(attr->value.u32);
	break;
  case SAI_TAM_EVENT_THRESHOLD_ATTR_LOW_WATERMARK:
	req.set_low_watermark(attr->value.u32);
	break;
  case SAI_TAM_EVENT_THRESHOLD_ATTR_LATENCY:
	req.set_latency(attr->value.u32);
	break;
  case SAI_TAM_EVENT_THRESHOLD_ATTR_RATE:
	req.set_rate(attr->value.u32);
	break;
  case SAI_TAM_EVENT_THRESHOLD_ATTR_ABS_VALUE:
	req.set_abs_value(attr->value.u32);
	break;
  case SAI_TAM_EVENT_THRESHOLD_ATTR_UNIT:
	req.set_unit(static_cast<lemming::dataplane::sai::TamEventThresholdUnit>(attr->value.s32 + 1));
	break;
}

	grpc::Status status = tam->SetTamEventThresholdAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_event_threshold_attribute(sai_object_id_t tam_event_threshold_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamEventThresholdAttributeRequest req;
	lemming::dataplane::sai::GetTamEventThresholdAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_event_threshold_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamEventThresholdAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamEventThresholdAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_EVENT_THRESHOLD_ATTR_HIGH_WATERMARK:
	 attr_list[i].value.u32 =   resp.attr().high_watermark();
	break;
  case SAI_TAM_EVENT_THRESHOLD_ATTR_LOW_WATERMARK:
	 attr_list[i].value.u32 =   resp.attr().low_watermark();
	break;
  case SAI_TAM_EVENT_THRESHOLD_ATTR_LATENCY:
	 attr_list[i].value.u32 =   resp.attr().latency();
	break;
  case SAI_TAM_EVENT_THRESHOLD_ATTR_RATE:
	 attr_list[i].value.u32 =   resp.attr().rate();
	break;
  case SAI_TAM_EVENT_THRESHOLD_ATTR_ABS_VALUE:
	 attr_list[i].value.u32 =   resp.attr().abs_value();
	break;
  case SAI_TAM_EVENT_THRESHOLD_ATTR_UNIT:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().unit() - 1);
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tam_int(sai_object_id_t *tam_int_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamIntRequest req;
	lemming::dataplane::sai::CreateTamIntResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_INT_ATTR_TYPE:
	req.set_type(static_cast<lemming::dataplane::sai::TamIntType>(attr_list[i].value.s32 + 1));
	break;
  case SAI_TAM_INT_ATTR_DEVICE_ID:
	req.set_device_id(attr_list[i].value.u32);
	break;
  case SAI_TAM_INT_ATTR_IOAM_TRACE_TYPE:
	req.set_ioam_trace_type(attr_list[i].value.u32);
	break;
  case SAI_TAM_INT_ATTR_INT_PRESENCE_TYPE:
	req.set_int_presence_type(static_cast<lemming::dataplane::sai::TamIntPresenceType>(attr_list[i].value.s32 + 1));
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
	req.mutable_collector_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
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

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam_int(sai_object_id_t tam_int_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamIntRequest req;
	lemming::dataplane::sai::RemoveTamIntResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_int_id); 
	
	grpc::Status status = tam->RemoveTamInt(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_int_attribute(sai_object_id_t tam_int_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamIntAttributeRequest req;
	lemming::dataplane::sai::SetTamIntAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_int_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_INT_ATTR_IOAM_TRACE_TYPE:
	req.set_ioam_trace_type(attr->value.u32);
	break;
  case SAI_TAM_INT_ATTR_TRACE_VECTOR:
	req.set_trace_vector(attr->value.u16);
	break;
  case SAI_TAM_INT_ATTR_ACTION_VECTOR:
	req.set_action_vector(attr->value.u16);
	break;
  case SAI_TAM_INT_ATTR_P4_INT_INSTRUCTION_BITMAP:
	req.set_p4_int_instruction_bitmap(attr->value.u16);
	break;
  case SAI_TAM_INT_ATTR_METADATA_FRAGMENT_ENABLE:
	req.set_metadata_fragment_enable(attr->value.booldata);
	break;
  case SAI_TAM_INT_ATTR_METADATA_CHECKSUM_ENABLE:
	req.set_metadata_checksum_enable(attr->value.booldata);
	break;
  case SAI_TAM_INT_ATTR_REPORT_ALL_PACKETS:
	req.set_report_all_packets(attr->value.booldata);
	break;
  case SAI_TAM_INT_ATTR_FLOW_LIVENESS_PERIOD:
	req.set_flow_liveness_period(attr->value.u16);
	break;
  case SAI_TAM_INT_ATTR_LATENCY_SENSITIVITY:
	req.set_latency_sensitivity(attr->value.u8);
	break;
  case SAI_TAM_INT_ATTR_ACL_GROUP:
	req.set_acl_group(attr->value.oid);
	break;
  case SAI_TAM_INT_ATTR_MAX_HOP_COUNT:
	req.set_max_hop_count(attr->value.u8);
	break;
  case SAI_TAM_INT_ATTR_MAX_LENGTH:
	req.set_max_length(attr->value.u8);
	break;
  case SAI_TAM_INT_ATTR_NAME_SPACE_ID:
	req.set_name_space_id(attr->value.u8);
	break;
  case SAI_TAM_INT_ATTR_NAME_SPACE_ID_GLOBAL:
	req.set_name_space_id_global(attr->value.booldata);
	break;
  case SAI_TAM_INT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
	req.set_ingress_samplepacket_enable(attr->value.oid);
	break;
  case SAI_TAM_INT_ATTR_COLLECTOR_LIST:
	req.mutable_collector_list()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
  case SAI_TAM_INT_ATTR_MATH_FUNC:
	req.set_math_func(attr->value.oid);
	break;
}

	grpc::Status status = tam->SetTamIntAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_int_attribute(sai_object_id_t tam_int_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamIntAttributeRequest req;
	lemming::dataplane::sai::GetTamIntAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_int_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamIntAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamIntAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_INT_ATTR_TYPE:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().type() - 1);
	break;
  case SAI_TAM_INT_ATTR_DEVICE_ID:
	 attr_list[i].value.u32 =   resp.attr().device_id();
	break;
  case SAI_TAM_INT_ATTR_IOAM_TRACE_TYPE:
	 attr_list[i].value.u32 =   resp.attr().ioam_trace_type();
	break;
  case SAI_TAM_INT_ATTR_INT_PRESENCE_TYPE:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().int_presence_type() - 1);
	break;
  case SAI_TAM_INT_ATTR_INT_PRESENCE_PB1:
	 attr_list[i].value.u32 =   resp.attr().int_presence_pb1();
	break;
  case SAI_TAM_INT_ATTR_INT_PRESENCE_PB2:
	 attr_list[i].value.u32 =   resp.attr().int_presence_pb2();
	break;
  case SAI_TAM_INT_ATTR_INT_PRESENCE_DSCP_VALUE:
	 attr_list[i].value.u8 =   resp.attr().int_presence_dscp_value();
	break;
  case SAI_TAM_INT_ATTR_INLINE:
	 attr_list[i].value.booldata =   resp.attr().inline_();
	break;
  case SAI_TAM_INT_ATTR_INT_PRESENCE_L3_PROTOCOL:
	 attr_list[i].value.u8 =   resp.attr().int_presence_l3_protocol();
	break;
  case SAI_TAM_INT_ATTR_TRACE_VECTOR:
	 attr_list[i].value.u16 =   resp.attr().trace_vector();
	break;
  case SAI_TAM_INT_ATTR_ACTION_VECTOR:
	 attr_list[i].value.u16 =   resp.attr().action_vector();
	break;
  case SAI_TAM_INT_ATTR_P4_INT_INSTRUCTION_BITMAP:
	 attr_list[i].value.u16 =   resp.attr().p4_int_instruction_bitmap();
	break;
  case SAI_TAM_INT_ATTR_METADATA_FRAGMENT_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().metadata_fragment_enable();
	break;
  case SAI_TAM_INT_ATTR_METADATA_CHECKSUM_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().metadata_checksum_enable();
	break;
  case SAI_TAM_INT_ATTR_REPORT_ALL_PACKETS:
	 attr_list[i].value.booldata =   resp.attr().report_all_packets();
	break;
  case SAI_TAM_INT_ATTR_FLOW_LIVENESS_PERIOD:
	 attr_list[i].value.u16 =   resp.attr().flow_liveness_period();
	break;
  case SAI_TAM_INT_ATTR_LATENCY_SENSITIVITY:
	 attr_list[i].value.u8 =   resp.attr().latency_sensitivity();
	break;
  case SAI_TAM_INT_ATTR_ACL_GROUP:
	 attr_list[i].value.oid =   resp.attr().acl_group();
	break;
  case SAI_TAM_INT_ATTR_MAX_HOP_COUNT:
	 attr_list[i].value.u8 =   resp.attr().max_hop_count();
	break;
  case SAI_TAM_INT_ATTR_MAX_LENGTH:
	 attr_list[i].value.u8 =   resp.attr().max_length();
	break;
  case SAI_TAM_INT_ATTR_NAME_SPACE_ID:
	 attr_list[i].value.u8 =   resp.attr().name_space_id();
	break;
  case SAI_TAM_INT_ATTR_NAME_SPACE_ID_GLOBAL:
	 attr_list[i].value.booldata =   resp.attr().name_space_id_global();
	break;
  case SAI_TAM_INT_ATTR_INGRESS_SAMPLEPACKET_ENABLE:
	 attr_list[i].value.oid =   resp.attr().ingress_samplepacket_enable();
	break;
  case SAI_TAM_INT_ATTR_COLLECTOR_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().collector_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_INT_ATTR_MATH_FUNC:
	 attr_list[i].value.oid =   resp.attr().math_func();
	break;
  case SAI_TAM_INT_ATTR_REPORT_ID:
	 attr_list[i].value.oid =   resp.attr().report_id();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tam_tel_type(sai_object_id_t *tam_tel_type_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamTelTypeRequest req;
	lemming::dataplane::sai::CreateTamTelTypeResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_TEL_TYPE_ATTR_TAM_TELEMETRY_TYPE:
	req.set_tam_telemetry_type(static_cast<lemming::dataplane::sai::TamTelemetryType>(attr_list[i].value.s32 + 1));
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
	req.set_switch_enable_resource_utilization_stats(attr_list[i].value.booldata);
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

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam_tel_type(sai_object_id_t tam_tel_type_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamTelTypeRequest req;
	lemming::dataplane::sai::RemoveTamTelTypeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_tel_type_id); 
	
	grpc::Status status = tam->RemoveTamTelType(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_tel_type_attribute(sai_object_id_t tam_tel_type_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamTelTypeAttributeRequest req;
	lemming::dataplane::sai::SetTamTelTypeAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_tel_type_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_TEL_TYPE_ATTR_INT_SWITCH_IDENTIFIER:
	req.set_int_switch_identifier(attr->value.u32);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS:
	req.set_switch_enable_port_stats(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_INGRESS:
	req.set_switch_enable_port_stats_ingress(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_EGRESS:
	req.set_switch_enable_port_stats_egress(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_VIRTUAL_QUEUE_STATS:
	req.set_switch_enable_virtual_queue_stats(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_OUTPUT_QUEUE_STATS:
	req.set_switch_enable_output_queue_stats(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_MMU_STATS:
	req.set_switch_enable_mmu_stats(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FABRIC_STATS:
	req.set_switch_enable_fabric_stats(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FILTER_STATS:
	req.set_switch_enable_filter_stats(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_RESOURCE_UTILIZATION_STATS:
	req.set_switch_enable_resource_utilization_stats(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_FABRIC_Q:
	req.set_fabric_q(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_NE_ENABLE:
	req.set_ne_enable(attr->value.booldata);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_DSCP_VALUE:
	req.set_dscp_value(attr->value.u8);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_MATH_FUNC:
	req.set_math_func(attr->value.oid);
	break;
}

	grpc::Status status = tam->SetTamTelTypeAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_tel_type_attribute(sai_object_id_t tam_tel_type_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamTelTypeAttributeRequest req;
	lemming::dataplane::sai::GetTamTelTypeAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_tel_type_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamTelTypeAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamTelTypeAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_TEL_TYPE_ATTR_TAM_TELEMETRY_TYPE:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().tam_telemetry_type() - 1);
	break;
  case SAI_TAM_TEL_TYPE_ATTR_INT_SWITCH_IDENTIFIER:
	 attr_list[i].value.u32 =   resp.attr().int_switch_identifier();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS:
	 attr_list[i].value.booldata =   resp.attr().switch_enable_port_stats();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_INGRESS:
	 attr_list[i].value.booldata =   resp.attr().switch_enable_port_stats_ingress();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_PORT_STATS_EGRESS:
	 attr_list[i].value.booldata =   resp.attr().switch_enable_port_stats_egress();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_VIRTUAL_QUEUE_STATS:
	 attr_list[i].value.booldata =   resp.attr().switch_enable_virtual_queue_stats();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_OUTPUT_QUEUE_STATS:
	 attr_list[i].value.booldata =   resp.attr().switch_enable_output_queue_stats();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_MMU_STATS:
	 attr_list[i].value.booldata =   resp.attr().switch_enable_mmu_stats();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FABRIC_STATS:
	 attr_list[i].value.booldata =   resp.attr().switch_enable_fabric_stats();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_FILTER_STATS:
	 attr_list[i].value.booldata =   resp.attr().switch_enable_filter_stats();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_SWITCH_ENABLE_RESOURCE_UTILIZATION_STATS:
	 attr_list[i].value.booldata =   resp.attr().switch_enable_resource_utilization_stats();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_FABRIC_Q:
	 attr_list[i].value.booldata =   resp.attr().fabric_q();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_NE_ENABLE:
	 attr_list[i].value.booldata =   resp.attr().ne_enable();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_DSCP_VALUE:
	 attr_list[i].value.u8 =   resp.attr().dscp_value();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_MATH_FUNC:
	 attr_list[i].value.oid =   resp.attr().math_func();
	break;
  case SAI_TAM_TEL_TYPE_ATTR_REPORT_ID:
	 attr_list[i].value.oid =   resp.attr().report_id();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tam_transport(sai_object_id_t *tam_transport_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamTransportRequest req;
	lemming::dataplane::sai::CreateTamTransportResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_TRANSPORT_ATTR_TRANSPORT_TYPE:
	req.set_transport_type(static_cast<lemming::dataplane::sai::TamTransportType>(attr_list[i].value.s32 + 1));
	break;
  case SAI_TAM_TRANSPORT_ATTR_SRC_PORT:
	req.set_src_port(attr_list[i].value.u32);
	break;
  case SAI_TAM_TRANSPORT_ATTR_DST_PORT:
	req.set_dst_port(attr_list[i].value.u32);
	break;
  case SAI_TAM_TRANSPORT_ATTR_TRANSPORT_AUTH_TYPE:
	req.set_transport_auth_type(static_cast<lemming::dataplane::sai::TamTransportAuthType>(attr_list[i].value.s32 + 1));
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

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam_transport(sai_object_id_t tam_transport_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamTransportRequest req;
	lemming::dataplane::sai::RemoveTamTransportResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_transport_id); 
	
	grpc::Status status = tam->RemoveTamTransport(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_transport_attribute(sai_object_id_t tam_transport_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamTransportAttributeRequest req;
	lemming::dataplane::sai::SetTamTransportAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_transport_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_TRANSPORT_ATTR_SRC_PORT:
	req.set_src_port(attr->value.u32);
	break;
  case SAI_TAM_TRANSPORT_ATTR_DST_PORT:
	req.set_dst_port(attr->value.u32);
	break;
  case SAI_TAM_TRANSPORT_ATTR_TRANSPORT_AUTH_TYPE:
	req.set_transport_auth_type(static_cast<lemming::dataplane::sai::TamTransportAuthType>(attr->value.s32 + 1));
	break;
  case SAI_TAM_TRANSPORT_ATTR_MTU:
	req.set_mtu(attr->value.u32);
	break;
}

	grpc::Status status = tam->SetTamTransportAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_transport_attribute(sai_object_id_t tam_transport_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamTransportAttributeRequest req;
	lemming::dataplane::sai::GetTamTransportAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_transport_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamTransportAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamTransportAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_TRANSPORT_ATTR_TRANSPORT_TYPE:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().transport_type() - 1);
	break;
  case SAI_TAM_TRANSPORT_ATTR_SRC_PORT:
	 attr_list[i].value.u32 =   resp.attr().src_port();
	break;
  case SAI_TAM_TRANSPORT_ATTR_DST_PORT:
	 attr_list[i].value.u32 =   resp.attr().dst_port();
	break;
  case SAI_TAM_TRANSPORT_ATTR_TRANSPORT_AUTH_TYPE:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().transport_auth_type() - 1);
	break;
  case SAI_TAM_TRANSPORT_ATTR_MTU:
	 attr_list[i].value.u32 =   resp.attr().mtu();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tam_telemetry(sai_object_id_t *tam_telemetry_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamTelemetryRequest req;
	lemming::dataplane::sai::CreateTamTelemetryResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_TELEMETRY_ATTR_TAM_TYPE_LIST:
	req.mutable_tam_type_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_TELEMETRY_ATTR_COLLECTOR_LIST:
	req.mutable_collector_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_TELEMETRY_ATTR_TAM_REPORTING_UNIT:
	req.set_tam_reporting_unit(static_cast<lemming::dataplane::sai::TamReportingUnit>(attr_list[i].value.s32 + 1));
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

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam_telemetry(sai_object_id_t tam_telemetry_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamTelemetryRequest req;
	lemming::dataplane::sai::RemoveTamTelemetryResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_telemetry_id); 
	
	grpc::Status status = tam->RemoveTamTelemetry(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_telemetry_attribute(sai_object_id_t tam_telemetry_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamTelemetryAttributeRequest req;
	lemming::dataplane::sai::SetTamTelemetryAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_telemetry_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_TELEMETRY_ATTR_TAM_TYPE_LIST:
	req.mutable_tam_type_list()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
  case SAI_TAM_TELEMETRY_ATTR_TAM_REPORTING_UNIT:
	req.set_tam_reporting_unit(static_cast<lemming::dataplane::sai::TamReportingUnit>(attr->value.s32 + 1));
	break;
  case SAI_TAM_TELEMETRY_ATTR_REPORTING_INTERVAL:
	req.set_reporting_interval(attr->value.u32);
	break;
}

	grpc::Status status = tam->SetTamTelemetryAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_telemetry_attribute(sai_object_id_t tam_telemetry_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamTelemetryAttributeRequest req;
	lemming::dataplane::sai::GetTamTelemetryAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_telemetry_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamTelemetryAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamTelemetryAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_TELEMETRY_ATTR_TAM_TYPE_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().tam_type_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_TELEMETRY_ATTR_COLLECTOR_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().collector_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_TELEMETRY_ATTR_TAM_REPORTING_UNIT:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().tam_reporting_unit() - 1);
	break;
  case SAI_TAM_TELEMETRY_ATTR_REPORTING_INTERVAL:
	 attr_list[i].value.u32 =   resp.attr().reporting_interval();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tam_collector(sai_object_id_t *tam_collector_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamCollectorRequest req;
	lemming::dataplane::sai::CreateTamCollectorResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

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

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam_collector(sai_object_id_t tam_collector_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamCollectorRequest req;
	lemming::dataplane::sai::RemoveTamCollectorResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_collector_id); 
	
	grpc::Status status = tam->RemoveTamCollector(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_collector_attribute(sai_object_id_t tam_collector_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamCollectorAttributeRequest req;
	lemming::dataplane::sai::SetTamCollectorAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_collector_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_COLLECTOR_ATTR_SRC_IP:
	req.set_src_ip(convert_from_ip_address(attr->value.ipaddr));
	break;
  case SAI_TAM_COLLECTOR_ATTR_DST_IP:
	req.set_dst_ip(convert_from_ip_address(attr->value.ipaddr));
	break;
  case SAI_TAM_COLLECTOR_ATTR_LOCALHOST:
	req.set_localhost(attr->value.booldata);
	break;
  case SAI_TAM_COLLECTOR_ATTR_VIRTUAL_ROUTER_ID:
	req.set_virtual_router_id(attr->value.oid);
	break;
  case SAI_TAM_COLLECTOR_ATTR_TRUNCATE_SIZE:
	req.set_truncate_size(attr->value.u16);
	break;
  case SAI_TAM_COLLECTOR_ATTR_TRANSPORT:
	req.set_transport(attr->value.oid);
	break;
  case SAI_TAM_COLLECTOR_ATTR_DSCP_VALUE:
	req.set_dscp_value(attr->value.u8);
	break;
}

	grpc::Status status = tam->SetTamCollectorAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_collector_attribute(sai_object_id_t tam_collector_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamCollectorAttributeRequest req;
	lemming::dataplane::sai::GetTamCollectorAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_collector_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamCollectorAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamCollectorAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_COLLECTOR_ATTR_SRC_IP:
	 attr_list[i].value.ipaddr =  convert_to_ip_address(resp.attr().src_ip());
	break;
  case SAI_TAM_COLLECTOR_ATTR_DST_IP:
	 attr_list[i].value.ipaddr =  convert_to_ip_address(resp.attr().dst_ip());
	break;
  case SAI_TAM_COLLECTOR_ATTR_LOCALHOST:
	 attr_list[i].value.booldata =   resp.attr().localhost();
	break;
  case SAI_TAM_COLLECTOR_ATTR_VIRTUAL_ROUTER_ID:
	 attr_list[i].value.oid =   resp.attr().virtual_router_id();
	break;
  case SAI_TAM_COLLECTOR_ATTR_TRUNCATE_SIZE:
	 attr_list[i].value.u16 =   resp.attr().truncate_size();
	break;
  case SAI_TAM_COLLECTOR_ATTR_TRANSPORT:
	 attr_list[i].value.oid =   resp.attr().transport();
	break;
  case SAI_TAM_COLLECTOR_ATTR_DSCP_VALUE:
	 attr_list[i].value.u8 =   resp.attr().dscp_value();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tam_event_action(sai_object_id_t *tam_event_action_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamEventActionRequest req;
	lemming::dataplane::sai::CreateTamEventActionResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

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

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam_event_action(sai_object_id_t tam_event_action_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamEventActionRequest req;
	lemming::dataplane::sai::RemoveTamEventActionResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_event_action_id); 
	
	grpc::Status status = tam->RemoveTamEventAction(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_event_action_attribute(sai_object_id_t tam_event_action_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamEventActionAttributeRequest req;
	lemming::dataplane::sai::SetTamEventActionAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_event_action_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_EVENT_ACTION_ATTR_REPORT_TYPE:
	req.set_report_type(attr->value.oid);
	break;
  case SAI_TAM_EVENT_ACTION_ATTR_QOS_ACTION_TYPE:
	req.set_qos_action_type(attr->value.u32);
	break;
}

	grpc::Status status = tam->SetTamEventActionAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_event_action_attribute(sai_object_id_t tam_event_action_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamEventActionAttributeRequest req;
	lemming::dataplane::sai::GetTamEventActionAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_event_action_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamEventActionAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamEventActionAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_EVENT_ACTION_ATTR_REPORT_TYPE:
	 attr_list[i].value.oid =   resp.attr().report_type();
	break;
  case SAI_TAM_EVENT_ACTION_ATTR_QOS_ACTION_TYPE:
	 attr_list[i].value.u32 =   resp.attr().qos_action_type();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_tam_event(sai_object_id_t *tam_event_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateTamEventRequest req;
	lemming::dataplane::sai::CreateTamEventResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_EVENT_ATTR_TYPE:
	req.set_type(static_cast<lemming::dataplane::sai::TamEventType>(attr_list[i].value.s32 + 1));
	break;
  case SAI_TAM_EVENT_ATTR_ACTION_LIST:
	req.mutable_action_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_EVENT_ATTR_COLLECTOR_LIST:
	req.mutable_collector_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
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

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_tam_event(sai_object_id_t tam_event_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveTamEventRequest req;
	lemming::dataplane::sai::RemoveTamEventResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_event_id); 
	
	grpc::Status status = tam->RemoveTamEvent(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_tam_event_attribute(sai_object_id_t tam_event_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetTamEventAttributeRequest req;
	lemming::dataplane::sai::SetTamEventAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(tam_event_id); 
	
	

switch (attr->id) {
  
  case SAI_TAM_EVENT_ATTR_THRESHOLD:
	req.set_threshold(attr->value.oid);
	break;
  case SAI_TAM_EVENT_ATTR_DSCP_VALUE:
	req.set_dscp_value(attr->value.u8);
	break;
}

	grpc::Status status = tam->SetTamEventAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_tam_event_attribute(sai_object_id_t tam_event_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetTamEventAttributeRequest req;
	lemming::dataplane::sai::GetTamEventAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(tam_event_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::TamEventAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = tam->GetTamEventAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_TAM_EVENT_ATTR_TYPE:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().type() - 1);
	break;
  case SAI_TAM_EVENT_ATTR_ACTION_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().action_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_EVENT_ATTR_COLLECTOR_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().collector_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_TAM_EVENT_ATTR_THRESHOLD:
	 attr_list[i].value.oid =   resp.attr().threshold();
	break;
  case SAI_TAM_EVENT_ATTR_DSCP_VALUE:
	 attr_list[i].value.u8 =   resp.attr().dscp_value();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

