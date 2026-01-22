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

#include "dataplane/standalone/sai/fdb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/fdb.pb.h"
#include <glog/logging.h>

const sai_fdb_api_t l_fdb = {
	.create_fdb_entry = l_create_fdb_entry,
	.remove_fdb_entry = l_remove_fdb_entry,
	.set_fdb_entry_attribute = l_set_fdb_entry_attribute,
	.get_fdb_entry_attribute = l_get_fdb_entry_attribute,
	.flush_fdb_entries = l_flush_fdb_entries,
	.create_fdb_entries = l_create_fdb_entries,
	.remove_fdb_entries = l_remove_fdb_entries,
	.set_fdb_entries_attribute = l_set_fdb_entries_attribute,
	.get_fdb_entries_attribute = l_get_fdb_entries_attribute,
};


lemming::dataplane::sai::CreateFdbEntryRequest convert_create_fdb_entry(uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateFdbEntryRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_FDB_ENTRY_ATTR_TYPE:
	msg.set_type(convert_sai_fdb_entry_type_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_FDB_ENTRY_ATTR_PACKET_ACTION:
	msg.set_packet_action(convert_sai_packet_action_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_FDB_ENTRY_ATTR_USER_TRAP_ID:
	msg.set_user_trap_id(attr_list[i].value.oid);
	break;
  case SAI_FDB_ENTRY_ATTR_BRIDGE_PORT_ID:
	msg.set_bridge_port_id(attr_list[i].value.oid);
	break;
  case SAI_FDB_ENTRY_ATTR_META_DATA:
	msg.set_meta_data(attr_list[i].value.u32);
	break;
  case SAI_FDB_ENTRY_ATTR_ENDPOINT_IP:
	msg.set_endpoint_ip(convert_from_ip_address(attr_list[i].value.ipaddr));
	break;
  case SAI_FDB_ENTRY_ATTR_COUNTER_ID:
	msg.set_counter_id(attr_list[i].value.oid);
	break;
  case SAI_FDB_ENTRY_ATTR_ALLOW_MAC_MOVE:
	msg.set_allow_mac_move(attr_list[i].value.booldata);
	break;
}

}
return msg;
}

lemming::dataplane::sai::FlushFdbEntriesRequest convert_flush_fdb_entries(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::FlushFdbEntriesRequest msg;
msg.set_switch_(switch_id);
 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_FDB_FLUSH_ATTR_BRIDGE_PORT_ID:
	msg.set_bridge_port_id(attr_list[i].value.oid);
	break;
  case SAI_FDB_FLUSH_ATTR_BV_ID:
	msg.set_bv_id(attr_list[i].value.oid);
	break;
  case SAI_FDB_FLUSH_ATTR_ENTRY_TYPE:
	msg.set_entry_type(convert_sai_fdb_flush_entry_type_t_to_proto(attr_list[i].value.s32));
	break;
}

}
return msg;
}

sai_status_t l_create_fdb_entry(const sai_fdb_entry_t *fdb_entry, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateFdbEntryRequest req = convert_create_fdb_entry(attr_count, attr_list);
	lemming::dataplane::sai::CreateFdbEntryResponse resp;
	grpc::ClientContext context;
	
	
	grpc::Status status = fdb->CreateFdbEntry(&context, req, &resp);
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

sai_status_t l_remove_fdb_entry(const sai_fdb_entry_t *fdb_entry) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveFdbEntryRequest req;
	lemming::dataplane::sai::RemoveFdbEntryResponse resp;
	grpc::ClientContext context;
	
	
	grpc::Status status = fdb->RemoveFdbEntry(&context, req, &resp);
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

sai_status_t l_set_fdb_entry_attribute(const sai_fdb_entry_t *fdb_entry, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetFdbEntryAttributeRequest req;
	lemming::dataplane::sai::SetFdbEntryAttributeResponse resp;
	grpc::ClientContext context;
	
	
	
	

switch (attr->id) {
  
  case SAI_FDB_ENTRY_ATTR_TYPE:
	req.set_type(convert_sai_fdb_entry_type_t_to_proto(attr->value.s32));
	break;
  case SAI_FDB_ENTRY_ATTR_PACKET_ACTION:
	req.set_packet_action(convert_sai_packet_action_t_to_proto(attr->value.s32));
	break;
  case SAI_FDB_ENTRY_ATTR_USER_TRAP_ID:
	req.set_user_trap_id(attr->value.oid);
	break;
  case SAI_FDB_ENTRY_ATTR_BRIDGE_PORT_ID:
	req.set_bridge_port_id(attr->value.oid);
	break;
  case SAI_FDB_ENTRY_ATTR_META_DATA:
	req.set_meta_data(attr->value.u32);
	break;
  case SAI_FDB_ENTRY_ATTR_ENDPOINT_IP:
	req.set_endpoint_ip(convert_from_ip_address(attr->value.ipaddr));
	break;
  case SAI_FDB_ENTRY_ATTR_COUNTER_ID:
	req.set_counter_id(attr->value.oid);
	break;
  case SAI_FDB_ENTRY_ATTR_ALLOW_MAC_MOVE:
	req.set_allow_mac_move(attr->value.booldata);
	break;
}

	grpc::Status status = fdb->SetFdbEntryAttribute(&context, req, &resp);
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

sai_status_t l_get_fdb_entry_attribute(const sai_fdb_entry_t *fdb_entry, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetFdbEntryAttributeRequest req;
	lemming::dataplane::sai::GetFdbEntryAttributeResponse resp;
	grpc::ClientContext context;
	
	

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_fdb_entry_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = fdb->GetFdbEntryAttribute(&context, req, &resp);
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
  
  case SAI_FDB_ENTRY_ATTR_TYPE:
	 attr_list[i].value.s32 =  convert_sai_fdb_entry_type_t_to_sai(resp.attr().type());
	break;
  case SAI_FDB_ENTRY_ATTR_PACKET_ACTION:
	 attr_list[i].value.s32 =  convert_sai_packet_action_t_to_sai(resp.attr().packet_action());
	break;
  case SAI_FDB_ENTRY_ATTR_USER_TRAP_ID:
	 attr_list[i].value.oid =   resp.attr().user_trap_id();
	break;
  case SAI_FDB_ENTRY_ATTR_BRIDGE_PORT_ID:
	 attr_list[i].value.oid =   resp.attr().bridge_port_id();
	break;
  case SAI_FDB_ENTRY_ATTR_META_DATA:
	 attr_list[i].value.u32 =   resp.attr().meta_data();
	break;
  case SAI_FDB_ENTRY_ATTR_ENDPOINT_IP:
	 attr_list[i].value.ipaddr =  convert_to_ip_address(resp.attr().endpoint_ip());
	break;
  case SAI_FDB_ENTRY_ATTR_COUNTER_ID:
	 attr_list[i].value.oid =   resp.attr().counter_id();
	break;
  case SAI_FDB_ENTRY_ATTR_ALLOW_MAC_MOVE:
	 attr_list[i].value.booldata =   resp.attr().allow_mac_move();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_flush_fdb_entries(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	lemming::dataplane::sai::FlushFdbEntriesRequest req = convert_flush_fdb_entries(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::FlushFdbEntriesResponse resp;
	grpc::ClientContext context;

	grpc::Status status = fdb->FlushFdbEntries(&context, req, &resp);
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

sai_status_t l_create_fdb_entries(uint32_t object_count, const sai_fdb_entry_t *fdb_entry, const uint32_t *attr_count, const sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateFdbEntriesRequest req;
	lemming::dataplane::sai::CreateFdbEntriesResponse resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) { 
		auto r = convert_create_fdb_entry(attr_count[i], attr_list[i]);
		*req.add_reqs() = r;
	}

	grpc::Status status = fdb->CreateFdbEntries(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
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

sai_status_t l_remove_fdb_entries(uint32_t object_count, const sai_fdb_entry_t *fdb_entry, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveFdbEntriesRequest req;
	lemming::dataplane::sai::RemoveFdbEntriesResponse resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) {
		
		
	}

	grpc::Status status = fdb->RemoveFdbEntries(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
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

sai_status_t l_set_fdb_entries_attribute(uint32_t object_count, const sai_fdb_entry_t *fdb_entry, const sai_attribute_t *attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_fdb_entries_attribute(uint32_t object_count, const sai_fdb_entry_t *fdb_entry, const uint32_t *attr_count, sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return SAI_STATUS_NOT_IMPLEMENTED;
}

