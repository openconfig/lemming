




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

#include "dataplane/standalone/sai/debug_counter.h"
#include <glog/logging.h>
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"
#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/debug_counter.pb.h"

const sai_debug_counter_api_t l_debug_counter = {
	.create_debug_counter = l_create_debug_counter,
	.remove_debug_counter = l_remove_debug_counter,
	.set_debug_counter_attribute = l_set_debug_counter_attribute,
	.get_debug_counter_attribute = l_get_debug_counter_attribute,
};


sai_status_t l_create_debug_counter(sai_object_id_t *debug_counter_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateDebugCounterRequest req;
	lemming::dataplane::sai::CreateDebugCounterResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
 	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_DEBUG_COUNTER_ATTR_TYPE:
	req.set_type(static_cast<lemming::dataplane::sai::DebugCounterType>(attr_list[i].value.s32 + 1));
	break;
  case SAI_DEBUG_COUNTER_ATTR_BIND_METHOD:
	req.set_bind_method(static_cast<lemming::dataplane::sai::DebugCounterBindMethod>(attr_list[i].value.s32 + 1));
	break;
}

	}
	grpc::Status status = debug_counter->CreateDebugCounter(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	*debug_counter_id = resp.oid(); 

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_debug_counter(sai_object_id_t debug_counter_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveDebugCounterRequest req;
	lemming::dataplane::sai::RemoveDebugCounterResponse resp;
	grpc::ClientContext context;
	req.set_oid(debug_counter_id); 
	
	grpc::Status status = debug_counter->RemoveDebugCounter(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_debug_counter_attribute(sai_object_id_t debug_counter_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_debug_counter_attribute(sai_object_id_t debug_counter_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetDebugCounterAttributeRequest req;
	lemming::dataplane::sai::GetDebugCounterAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(debug_counter_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(static_cast<lemming::dataplane::sai::DebugCounterAttr>(attr_list[i].id + 1));
	}
	grpc::Status status = debug_counter->GetDebugCounterAttribute(&context, req, &resp);
	if (!status.ok()) {
		LOG(ERROR) << status.error_message();
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		

switch (attr_list[i].id) {
  
  case SAI_DEBUG_COUNTER_ATTR_INDEX:
	 attr_list[i].value.u32 =   resp.attr().index();
	break;
  case SAI_DEBUG_COUNTER_ATTR_TYPE:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().type() - 1);
	break;
  case SAI_DEBUG_COUNTER_ATTR_BIND_METHOD:
	 attr_list[i].value.s32 =  static_cast<int>(resp.attr().bind_method() - 1);
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

