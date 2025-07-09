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

#include "dataplane/standalone/sai/nat.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/nat.pb.h"
#include <glog/logging.h>

const sai_nat_api_t l_nat = {
	.create_nat_entry = l_create_nat_entry,
	.remove_nat_entry = l_remove_nat_entry,
	.set_nat_entry_attribute = l_set_nat_entry_attribute,
	.get_nat_entry_attribute = l_get_nat_entry_attribute,
	.create_nat_entries = l_create_nat_entries,
	.remove_nat_entries = l_remove_nat_entries,
	.set_nat_entries_attribute = l_set_nat_entries_attribute,
	.get_nat_entries_attribute = l_get_nat_entries_attribute,
	.create_nat_zone_counter = l_create_nat_zone_counter,
	.remove_nat_zone_counter = l_remove_nat_zone_counter,
	.set_nat_zone_counter_attribute = l_set_nat_zone_counter_attribute,
	.get_nat_zone_counter_attribute = l_get_nat_zone_counter_attribute,
};


lemming::dataplane::sai::CreateNatEntryRequest convert_create_nat_entry(uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateNatEntryRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_NAT_ENTRY_ATTR_NAT_TYPE:
	msg.set_nat_type(convert_sai_nat_type_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_NAT_ENTRY_ATTR_SRC_IP:
	msg.set_src_ip(&attr_list[i].value.ip4, sizeof(attr_list[i].value.ip4));
	break;
  case SAI_NAT_ENTRY_ATTR_SRC_IP_MASK:
	msg.set_src_ip_mask(&attr_list[i].value.ip4, sizeof(attr_list[i].value.ip4));
	break;
  case SAI_NAT_ENTRY_ATTR_VR_ID:
	msg.set_vr_id(attr_list[i].value.oid);
	break;
  case SAI_NAT_ENTRY_ATTR_DST_IP:
	msg.set_dst_ip(&attr_list[i].value.ip4, sizeof(attr_list[i].value.ip4));
	break;
  case SAI_NAT_ENTRY_ATTR_DST_IP_MASK:
	msg.set_dst_ip_mask(&attr_list[i].value.ip4, sizeof(attr_list[i].value.ip4));
	break;
  case SAI_NAT_ENTRY_ATTR_L4_SRC_PORT:
	msg.set_l4_src_port(attr_list[i].value.u16);
	break;
  case SAI_NAT_ENTRY_ATTR_L4_DST_PORT:
	msg.set_l4_dst_port(attr_list[i].value.u16);
	break;
  case SAI_NAT_ENTRY_ATTR_ENABLE_PACKET_COUNT:
	msg.set_enable_packet_count(attr_list[i].value.booldata);
	break;
  case SAI_NAT_ENTRY_ATTR_PACKET_COUNT:
	msg.set_packet_count(attr_list[i].value.u64);
	break;
  case SAI_NAT_ENTRY_ATTR_ENABLE_BYTE_COUNT:
	msg.set_enable_byte_count(attr_list[i].value.booldata);
	break;
  case SAI_NAT_ENTRY_ATTR_BYTE_COUNT:
	msg.set_byte_count(attr_list[i].value.u64);
	break;
  case SAI_NAT_ENTRY_ATTR_HIT_BIT_COR:
	msg.set_hit_bit_cor(attr_list[i].value.booldata);
	break;
  case SAI_NAT_ENTRY_ATTR_HIT_BIT:
	msg.set_hit_bit(attr_list[i].value.booldata);
	break;
  case SAI_NAT_ENTRY_ATTR_AGING_TIME:
	msg.set_aging_time(attr_list[i].value.u32);
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateNatZoneCounterRequest convert_create_nat_zone_counter(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateNatZoneCounterRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_NAT_ZONE_COUNTER_ATTR_NAT_TYPE:
	msg.set_nat_type(convert_sai_nat_type_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_ZONE_ID:
	msg.set_zone_id(attr_list[i].value.u8);
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_DISCARD:
	msg.set_enable_discard(attr_list[i].value.booldata);
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_DISCARD_PACKET_COUNT:
	msg.set_discard_packet_count(attr_list[i].value.u64);
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATION_NEEDED:
	msg.set_enable_translation_needed(attr_list[i].value.booldata);
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATION_NEEDED_PACKET_COUNT:
	msg.set_translation_needed_packet_count(attr_list[i].value.u64);
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATIONS:
	msg.set_enable_translations(attr_list[i].value.booldata);
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATIONS_PACKET_COUNT:
	msg.set_translations_packet_count(attr_list[i].value.u64);
	break;
}

}
return msg;
}

sai_status_t l_create_nat_entry(const sai_nat_entry_t *nat_entry, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateNatEntryRequest req = convert_create_nat_entry(attr_count, attr_list);
	lemming::dataplane::sai::CreateNatEntryResponse resp;
	grpc::ClientContext context;
	
	
	grpc::Status status = nat->CreateNatEntry(&context, req, &resp);
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

sai_status_t l_remove_nat_entry(const sai_nat_entry_t *nat_entry) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveNatEntryRequest req;
	lemming::dataplane::sai::RemoveNatEntryResponse resp;
	grpc::ClientContext context;
	
	
	grpc::Status status = nat->RemoveNatEntry(&context, req, &resp);
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

sai_status_t l_set_nat_entry_attribute(const sai_nat_entry_t *nat_entry, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetNatEntryAttributeRequest req;
	lemming::dataplane::sai::SetNatEntryAttributeResponse resp;
	grpc::ClientContext context;
	
	
	
	

switch (attr->id) {
  
  case SAI_NAT_ENTRY_ATTR_NAT_TYPE:
	req.set_nat_type(convert_sai_nat_type_t_to_proto(attr->value.s32));
	break;
  case SAI_NAT_ENTRY_ATTR_SRC_IP:
	req.set_src_ip(&attr->value.ip4, sizeof(attr->value.ip4));
	break;
  case SAI_NAT_ENTRY_ATTR_SRC_IP_MASK:
	req.set_src_ip_mask(&attr->value.ip4, sizeof(attr->value.ip4));
	break;
  case SAI_NAT_ENTRY_ATTR_VR_ID:
	req.set_vr_id(attr->value.oid);
	break;
  case SAI_NAT_ENTRY_ATTR_DST_IP:
	req.set_dst_ip(&attr->value.ip4, sizeof(attr->value.ip4));
	break;
  case SAI_NAT_ENTRY_ATTR_DST_IP_MASK:
	req.set_dst_ip_mask(&attr->value.ip4, sizeof(attr->value.ip4));
	break;
  case SAI_NAT_ENTRY_ATTR_L4_SRC_PORT:
	req.set_l4_src_port(attr->value.u16);
	break;
  case SAI_NAT_ENTRY_ATTR_L4_DST_PORT:
	req.set_l4_dst_port(attr->value.u16);
	break;
  case SAI_NAT_ENTRY_ATTR_ENABLE_PACKET_COUNT:
	req.set_enable_packet_count(attr->value.booldata);
	break;
  case SAI_NAT_ENTRY_ATTR_PACKET_COUNT:
	req.set_packet_count(attr->value.u64);
	break;
  case SAI_NAT_ENTRY_ATTR_ENABLE_BYTE_COUNT:
	req.set_enable_byte_count(attr->value.booldata);
	break;
  case SAI_NAT_ENTRY_ATTR_BYTE_COUNT:
	req.set_byte_count(attr->value.u64);
	break;
  case SAI_NAT_ENTRY_ATTR_HIT_BIT_COR:
	req.set_hit_bit_cor(attr->value.booldata);
	break;
  case SAI_NAT_ENTRY_ATTR_HIT_BIT:
	req.set_hit_bit(attr->value.booldata);
	break;
  case SAI_NAT_ENTRY_ATTR_AGING_TIME:
	req.set_aging_time(attr->value.u32);
	break;
}

	grpc::Status status = nat->SetNatEntryAttribute(&context, req, &resp);
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

sai_status_t l_get_nat_entry_attribute(const sai_nat_entry_t *nat_entry, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetNatEntryAttributeRequest req;
	lemming::dataplane::sai::GetNatEntryAttributeResponse resp;
	grpc::ClientContext context;
	
	

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_nat_entry_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = nat->GetNatEntryAttribute(&context, req, &resp);
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
  
  case SAI_NAT_ENTRY_ATTR_NAT_TYPE:
	 attr_list[i].value.s32 =  convert_sai_nat_type_t_to_sai(resp.attr().nat_type());
	break;
  case SAI_NAT_ENTRY_ATTR_SRC_IP:
	memcpy(&attr_list[i].value.ip4, resp.attr().src_ip().data(), sizeof(sai_ip4_t));
	break;
  case SAI_NAT_ENTRY_ATTR_SRC_IP_MASK:
	memcpy(&attr_list[i].value.ip4, resp.attr().src_ip_mask().data(), sizeof(sai_ip4_t));
	break;
  case SAI_NAT_ENTRY_ATTR_VR_ID:
	 attr_list[i].value.oid =   resp.attr().vr_id();
	break;
  case SAI_NAT_ENTRY_ATTR_DST_IP:
	memcpy(&attr_list[i].value.ip4, resp.attr().dst_ip().data(), sizeof(sai_ip4_t));
	break;
  case SAI_NAT_ENTRY_ATTR_DST_IP_MASK:
	memcpy(&attr_list[i].value.ip4, resp.attr().dst_ip_mask().data(), sizeof(sai_ip4_t));
	break;
  case SAI_NAT_ENTRY_ATTR_L4_SRC_PORT:
	 attr_list[i].value.u16 =   resp.attr().l4_src_port();
	break;
  case SAI_NAT_ENTRY_ATTR_L4_DST_PORT:
	 attr_list[i].value.u16 =   resp.attr().l4_dst_port();
	break;
  case SAI_NAT_ENTRY_ATTR_ENABLE_PACKET_COUNT:
	 attr_list[i].value.booldata =   resp.attr().enable_packet_count();
	break;
  case SAI_NAT_ENTRY_ATTR_PACKET_COUNT:
	 attr_list[i].value.u64 =   resp.attr().packet_count();
	break;
  case SAI_NAT_ENTRY_ATTR_ENABLE_BYTE_COUNT:
	 attr_list[i].value.booldata =   resp.attr().enable_byte_count();
	break;
  case SAI_NAT_ENTRY_ATTR_BYTE_COUNT:
	 attr_list[i].value.u64 =   resp.attr().byte_count();
	break;
  case SAI_NAT_ENTRY_ATTR_HIT_BIT_COR:
	 attr_list[i].value.booldata =   resp.attr().hit_bit_cor();
	break;
  case SAI_NAT_ENTRY_ATTR_HIT_BIT:
	 attr_list[i].value.booldata =   resp.attr().hit_bit();
	break;
  case SAI_NAT_ENTRY_ATTR_AGING_TIME:
	 attr_list[i].value.u32 =   resp.attr().aging_time();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_nat_entries(uint32_t object_count, const sai_nat_entry_t *nat_entry, const uint32_t *attr_count, const sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateNatEntriesRequest req;
	lemming::dataplane::sai::CreateNatEntriesResponse resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) { 
		auto r = convert_create_nat_entry(attr_count[i], attr_list[i]);
		*req.add_reqs() = r;
	}

	grpc::Status status = nat->CreateNatEntries(&context, req, &resp);
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

sai_status_t l_remove_nat_entries(uint32_t object_count, const sai_nat_entry_t *nat_entry, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveNatEntriesRequest req;
	lemming::dataplane::sai::RemoveNatEntriesResponse resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) {
		
		
	}

	grpc::Status status = nat->RemoveNatEntries(&context, req, &resp);
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

sai_status_t l_set_nat_entries_attribute(uint32_t object_count, const sai_nat_entry_t *nat_entry, const sai_attribute_t *attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_nat_entries_attribute(uint32_t object_count, const sai_nat_entry_t *nat_entry, const uint32_t *attr_count, sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_nat_zone_counter(sai_object_id_t *nat_zone_counter_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateNatZoneCounterRequest req = convert_create_nat_zone_counter(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateNatZoneCounterResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = nat->CreateNatZoneCounter(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (nat_zone_counter_id) {
	*nat_zone_counter_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_nat_zone_counter(sai_object_id_t nat_zone_counter_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveNatZoneCounterRequest req;
	lemming::dataplane::sai::RemoveNatZoneCounterResponse resp;
	grpc::ClientContext context;
	req.set_oid(nat_zone_counter_id); 
	
	grpc::Status status = nat->RemoveNatZoneCounter(&context, req, &resp);
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

sai_status_t l_set_nat_zone_counter_attribute(sai_object_id_t nat_zone_counter_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetNatZoneCounterAttributeRequest req;
	lemming::dataplane::sai::SetNatZoneCounterAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(nat_zone_counter_id); 
	
	
	

switch (attr->id) {
  
  case SAI_NAT_ZONE_COUNTER_ATTR_NAT_TYPE:
	req.set_nat_type(convert_sai_nat_type_t_to_proto(attr->value.s32));
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_ZONE_ID:
	req.set_zone_id(attr->value.u8);
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_DISCARD_PACKET_COUNT:
	req.set_discard_packet_count(attr->value.u64);
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATION_NEEDED_PACKET_COUNT:
	req.set_translation_needed_packet_count(attr->value.u64);
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATIONS_PACKET_COUNT:
	req.set_translations_packet_count(attr->value.u64);
	break;
}

	grpc::Status status = nat->SetNatZoneCounterAttribute(&context, req, &resp);
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

sai_status_t l_get_nat_zone_counter_attribute(sai_object_id_t nat_zone_counter_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetNatZoneCounterAttributeRequest req;
	lemming::dataplane::sai::GetNatZoneCounterAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(nat_zone_counter_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_nat_zone_counter_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = nat->GetNatZoneCounterAttribute(&context, req, &resp);
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
  
  case SAI_NAT_ZONE_COUNTER_ATTR_NAT_TYPE:
	 attr_list[i].value.s32 =  convert_sai_nat_type_t_to_sai(resp.attr().nat_type());
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_ZONE_ID:
	 attr_list[i].value.u8 =   resp.attr().zone_id();
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_DISCARD:
	 attr_list[i].value.booldata =   resp.attr().enable_discard();
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_DISCARD_PACKET_COUNT:
	 attr_list[i].value.u64 =   resp.attr().discard_packet_count();
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATION_NEEDED:
	 attr_list[i].value.booldata =   resp.attr().enable_translation_needed();
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATION_NEEDED_PACKET_COUNT:
	 attr_list[i].value.u64 =   resp.attr().translation_needed_packet_count();
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_ENABLE_TRANSLATIONS:
	 attr_list[i].value.booldata =   resp.attr().enable_translations();
	break;
  case SAI_NAT_ZONE_COUNTER_ATTR_TRANSLATIONS_PACKET_COUNT:
	 attr_list[i].value.u64 =   resp.attr().translations_packet_count();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

