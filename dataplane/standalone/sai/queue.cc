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

#include "dataplane/standalone/sai/queue.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/queue.pb.h"
#include <glog/logging.h>

const sai_queue_api_t l_queue = {
	.create_queue = l_create_queue,
	.remove_queue = l_remove_queue,
	.set_queue_attribute = l_set_queue_attribute,
	.get_queue_attribute = l_get_queue_attribute,
	.get_queue_stats = l_get_queue_stats,
	.get_queue_stats_ext = l_get_queue_stats_ext,
	.clear_queue_stats = l_clear_queue_stats,
};


lemming::dataplane::sai::CreateQueueRequest convert_create_queue(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateQueueRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_QUEUE_ATTR_TYPE:
	msg.set_type(convert_sai_queue_type_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_QUEUE_ATTR_PORT:
	msg.set_port(attr_list[i].value.oid);
	break;
  case SAI_QUEUE_ATTR_INDEX:
	msg.set_index(attr_list[i].value.u8);
	break;
  case SAI_QUEUE_ATTR_PARENT_SCHEDULER_NODE:
	msg.set_parent_scheduler_node(attr_list[i].value.oid);
	break;
  case SAI_QUEUE_ATTR_WRED_PROFILE_ID:
	msg.set_wred_profile_id(attr_list[i].value.oid);
	break;
  case SAI_QUEUE_ATTR_BUFFER_PROFILE_ID:
	msg.set_buffer_profile_id(attr_list[i].value.oid);
	break;
  case SAI_QUEUE_ATTR_SCHEDULER_PROFILE_ID:
	msg.set_scheduler_profile_id(attr_list[i].value.oid);
	break;
  case SAI_QUEUE_ATTR_ENABLE_PFC_DLDR:
	msg.set_enable_pfc_dldr(attr_list[i].value.booldata);
	break;
  case SAI_QUEUE_ATTR_PFC_DLR_INIT:
	msg.set_pfc_dlr_init(attr_list[i].value.booldata);
	break;
  case SAI_QUEUE_ATTR_TAM_OBJECT:
	msg.mutable_tam_object()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_QUEUE_ATTR_PFC_DLR_PACKET_ACTION:
	msg.set_pfc_dlr_packet_action(convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
	break;
}

}
return msg;
}

sai_status_t l_create_queue(sai_object_id_t *queue_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateQueueRequest req = convert_create_queue(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateQueueResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = queue->CreateQueue(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (queue_id) {
	*queue_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_queue(sai_object_id_t queue_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveQueueRequest req;
	lemming::dataplane::sai::RemoveQueueResponse resp;
	grpc::ClientContext context;
	req.set_oid(queue_id); 
	
	grpc::Status status = queue->RemoveQueue(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_queue_attribute(sai_object_id_t queue_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetQueueAttributeRequest req;
	lemming::dataplane::sai::SetQueueAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(queue_id); 
	
	
	

switch (attr->id) {
  
  case SAI_QUEUE_ATTR_PARENT_SCHEDULER_NODE:
	req.set_parent_scheduler_node(attr->value.oid);
	break;
  case SAI_QUEUE_ATTR_WRED_PROFILE_ID:
	req.set_wred_profile_id(attr->value.oid);
	break;
  case SAI_QUEUE_ATTR_BUFFER_PROFILE_ID:
	req.set_buffer_profile_id(attr->value.oid);
	break;
  case SAI_QUEUE_ATTR_SCHEDULER_PROFILE_ID:
	req.set_scheduler_profile_id(attr->value.oid);
	break;
  case SAI_QUEUE_ATTR_ENABLE_PFC_DLDR:
	req.set_enable_pfc_dldr(attr->value.booldata);
	break;
  case SAI_QUEUE_ATTR_PFC_DLR_INIT:
	req.set_pfc_dlr_init(attr->value.booldata);
	break;
  case SAI_QUEUE_ATTR_TAM_OBJECT:
	req.mutable_tam_object()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
  case SAI_QUEUE_ATTR_PFC_DLR_PACKET_ACTION:
	req.set_pfc_dlr_packet_action(convert_sai_packet_action_t_to_proto(attr->value.s32));
	break;
}

	grpc::Status status = queue->SetQueueAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_queue_attribute(sai_object_id_t queue_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetQueueAttributeRequest req;
	lemming::dataplane::sai::GetQueueAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(queue_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_queue_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = queue->GetQueueAttribute(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < attr_count; i++ ) {
		
		

switch (attr_list[i].id) {
  
  case SAI_QUEUE_ATTR_TYPE:
	 attr_list[i].value.s32 =  convert_sai_queue_type_t_to_sai(resp.attr().type());
	break;
  case SAI_QUEUE_ATTR_PORT:
	 attr_list[i].value.oid =   resp.attr().port();
	break;
  case SAI_QUEUE_ATTR_INDEX:
	 attr_list[i].value.u8 =   resp.attr().index();
	break;
  case SAI_QUEUE_ATTR_PARENT_SCHEDULER_NODE:
	 attr_list[i].value.oid =   resp.attr().parent_scheduler_node();
	break;
  case SAI_QUEUE_ATTR_WRED_PROFILE_ID:
	 attr_list[i].value.oid =   resp.attr().wred_profile_id();
	break;
  case SAI_QUEUE_ATTR_BUFFER_PROFILE_ID:
	 attr_list[i].value.oid =   resp.attr().buffer_profile_id();
	break;
  case SAI_QUEUE_ATTR_SCHEDULER_PROFILE_ID:
	 attr_list[i].value.oid =   resp.attr().scheduler_profile_id();
	break;
  case SAI_QUEUE_ATTR_PAUSE_STATUS:
	 attr_list[i].value.booldata =   resp.attr().pause_status();
	break;
  case SAI_QUEUE_ATTR_ENABLE_PFC_DLDR:
	 attr_list[i].value.booldata =   resp.attr().enable_pfc_dldr();
	break;
  case SAI_QUEUE_ATTR_PFC_DLR_INIT:
	 attr_list[i].value.booldata =   resp.attr().pfc_dlr_init();
	break;
  case SAI_QUEUE_ATTR_TAM_OBJECT:
	copy_list(attr_list[i].value.objlist.list, resp.attr().tam_object(), &attr_list[i].value.objlist.count);
	break;
  case SAI_QUEUE_ATTR_PFC_DLR_PACKET_ACTION:
	 attr_list[i].value.s32 =  convert_sai_packet_action_t_to_sai(resp.attr().pfc_dlr_packet_action());
	break;
  case SAI_QUEUE_ATTR_PFC_CONTINUOUS_DEADLOCK_STATE:
	 attr_list[i].value.s32 =  convert_sai_queue_pfc_continuous_deadlock_state_t_to_sai(resp.attr().pfc_continuous_deadlock_state());
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_queue_stats(sai_object_id_t queue_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetQueueStatsRequest req;
	lemming::dataplane::sai::GetQueueStatsResponse resp;
	grpc::ClientContext context;
	req.set_oid(queue_id); 
	
	for (uint32_t i = 0; i < number_of_counters; i++) {
		req.add_counter_ids(convert_sai_queue_stat_t_to_proto(counter_ids[i]));
	}
	grpc::Status status = queue->GetQueueStats(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	for(uint32_t i = 0; i < number_of_counters && i < uint32_t(resp.values_size()); i++ ) {
		counters[i] = resp.values(i);
	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_queue_stats_ext(sai_object_id_t queue_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids, sai_stats_mode_t mode, uint64_t *counters) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_queue_stats(sai_object_id_t queue_id, uint32_t number_of_counters, const sai_stat_id_t *counter_ids) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

