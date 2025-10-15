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

#include "dataplane/standalone/sai/next_hop_group.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/next_hop_group.pb.h"
#include <glog/logging.h>

const sai_next_hop_group_api_t l_next_hop_group = {
	.create_next_hop_group = l_create_next_hop_group,
	.remove_next_hop_group = l_remove_next_hop_group,
	.set_next_hop_group_attribute = l_set_next_hop_group_attribute,
	.get_next_hop_group_attribute = l_get_next_hop_group_attribute,
	.create_next_hop_group_member = l_create_next_hop_group_member,
	.remove_next_hop_group_member = l_remove_next_hop_group_member,
	.set_next_hop_group_member_attribute = l_set_next_hop_group_member_attribute,
	.get_next_hop_group_member_attribute = l_get_next_hop_group_member_attribute,
	.create_next_hop_group_members = l_create_next_hop_group_members,
	.remove_next_hop_group_members = l_remove_next_hop_group_members,
	.create_next_hop_group_map = l_create_next_hop_group_map,
	.remove_next_hop_group_map = l_remove_next_hop_group_map,
	.set_next_hop_group_map_attribute = l_set_next_hop_group_map_attribute,
	.get_next_hop_group_map_attribute = l_get_next_hop_group_map_attribute,
	.set_next_hop_group_members_attribute = l_set_next_hop_group_members_attribute,
	.get_next_hop_group_members_attribute = l_get_next_hop_group_members_attribute,
	.create_next_hop_groups = l_create_next_hop_groups,
	.remove_next_hop_groups = l_remove_next_hop_groups,
	.set_next_hop_groups_attribute = l_set_next_hop_groups_attribute,
	.get_next_hop_groups_attribute = l_get_next_hop_groups_attribute,
};


lemming::dataplane::sai::CreateNextHopGroupRequest convert_create_next_hop_group(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateNextHopGroupRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_NEXT_HOP_GROUP_ATTR_TYPE:
	msg.set_type(convert_sai_next_hop_group_type_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_SET_SWITCHOVER:
	msg.set_set_switchover(attr_list[i].value.booldata);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_COUNTER_ID:
	msg.set_counter_id(attr_list[i].value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_CONFIGURED_SIZE:
	msg.set_configured_size(attr_list[i].value.u32);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_SELECTION_MAP:
	msg.set_selection_map(attr_list[i].value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_HIERARCHICAL_NEXTHOP:
	msg.set_hierarchical_nexthop(attr_list[i].value.booldata);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_ARS_OBJECT_ID:
	msg.set_ars_object_id(attr_list[i].value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_LIST:
	msg.mutable_next_hop_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_WEIGHT_LIST:
	msg.mutable_next_hop_member_weight_list()->Add(attr_list[i].value.u32list.list, attr_list[i].value.u32list.list + attr_list[i].value.u32list.count);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_COUNTER_LIST:
	msg.mutable_next_hop_member_counter_list()->Add(attr_list[i].value.objlist.list, attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateNextHopGroupMemberRequest convert_create_next_hop_group_member(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateNextHopGroupMemberRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_GROUP_ID:
	msg.set_next_hop_group_id(attr_list[i].value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_ID:
	msg.set_next_hop_id(attr_list[i].value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_WEIGHT:
	msg.set_weight(attr_list[i].value.u32);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_CONFIGURED_ROLE:
	msg.set_configured_role(convert_sai_next_hop_group_member_configured_role_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_MONITORED_OBJECT:
	msg.set_monitored_object(attr_list[i].value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_INDEX:
	msg.set_index(attr_list[i].value.u32);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_SEQUENCE_ID:
	msg.set_sequence_id(attr_list[i].value.u32);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_COUNTER_ID:
	msg.set_counter_id(attr_list[i].value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_ARS_ALTERNATE_PATH:
	msg.set_ars_alternate_path(attr_list[i].value.booldata);
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateNextHopGroupMapRequest convert_create_next_hop_group_map(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateNextHopGroupMapRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_NEXT_HOP_GROUP_MAP_ATTR_TYPE:
	msg.set_type(convert_sai_next_hop_group_map_type_t_to_proto(attr_list[i].value.s32));
	break;
}

}
return msg;
}

sai_status_t l_create_next_hop_group(sai_object_id_t *next_hop_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateNextHopGroupRequest req = convert_create_next_hop_group(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateNextHopGroupResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = next_hop_group->CreateNextHopGroup(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (next_hop_group_id) {
	*next_hop_group_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_next_hop_group(sai_object_id_t next_hop_group_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveNextHopGroupRequest req;
	lemming::dataplane::sai::RemoveNextHopGroupResponse resp;
	grpc::ClientContext context;
	req.set_oid(next_hop_group_id); 
	
	grpc::Status status = next_hop_group->RemoveNextHopGroup(&context, req, &resp);
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

sai_status_t l_set_next_hop_group_attribute(sai_object_id_t next_hop_group_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetNextHopGroupAttributeRequest req;
	lemming::dataplane::sai::SetNextHopGroupAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(next_hop_group_id); 
	
	
	

switch (attr->id) {
  
  case SAI_NEXT_HOP_GROUP_ATTR_SET_SWITCHOVER:
	req.set_set_switchover(attr->value.booldata);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_COUNTER_ID:
	req.set_counter_id(attr->value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_SELECTION_MAP:
	req.set_selection_map(attr->value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_ARS_OBJECT_ID:
	req.set_ars_object_id(attr->value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_LIST:
	req.mutable_next_hop_list()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_WEIGHT_LIST:
	req.mutable_next_hop_member_weight_list()->Add(attr->value.u32list.list, attr->value.u32list.list + attr->value.u32list.count);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_COUNTER_LIST:
	req.mutable_next_hop_member_counter_list()->Add(attr->value.objlist.list, attr->value.objlist.list + attr->value.objlist.count);
	break;
}

	grpc::Status status = next_hop_group->SetNextHopGroupAttribute(&context, req, &resp);
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

sai_status_t l_get_next_hop_group_attribute(sai_object_id_t next_hop_group_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetNextHopGroupAttributeRequest req;
	lemming::dataplane::sai::GetNextHopGroupAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(next_hop_group_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_next_hop_group_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = next_hop_group->GetNextHopGroupAttribute(&context, req, &resp);
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
  
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_COUNT:
	 attr_list[i].value.u32 =   resp.attr().next_hop_count();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().next_hop_member_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_TYPE:
	 attr_list[i].value.s32 =  convert_sai_next_hop_group_type_t_to_sai(resp.attr().type());
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_SET_SWITCHOVER:
	 attr_list[i].value.booldata =   resp.attr().set_switchover();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_COUNTER_ID:
	 attr_list[i].value.oid =   resp.attr().counter_id();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_CONFIGURED_SIZE:
	 attr_list[i].value.u32 =   resp.attr().configured_size();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_REAL_SIZE:
	 attr_list[i].value.u32 =   resp.attr().real_size();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_SELECTION_MAP:
	 attr_list[i].value.oid =   resp.attr().selection_map();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_HIERARCHICAL_NEXTHOP:
	 attr_list[i].value.booldata =   resp.attr().hierarchical_nexthop();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_ARS_OBJECT_ID:
	 attr_list[i].value.oid =   resp.attr().ars_object_id();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_ARS_PACKET_DROPS:
	 attr_list[i].value.u32 =   resp.attr().ars_packet_drops();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_ARS_NEXT_HOP_REASSIGNMENTS:
	 attr_list[i].value.u32 =   resp.attr().ars_next_hop_reassignments();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_ARS_PORT_REASSIGNMENTS:
	 attr_list[i].value.u32 =   resp.attr().ars_port_reassignments();
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().next_hop_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_WEIGHT_LIST:
	copy_list(attr_list[i].value.u32list.list, resp.attr().next_hop_member_weight_list(), &attr_list[i].value.u32list.count);
	break;
  case SAI_NEXT_HOP_GROUP_ATTR_NEXT_HOP_MEMBER_COUNTER_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().next_hop_member_counter_list(), &attr_list[i].value.objlist.count);
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_next_hop_group_member(sai_object_id_t *next_hop_group_member_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateNextHopGroupMemberRequest req = convert_create_next_hop_group_member(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateNextHopGroupMemberResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = next_hop_group->CreateNextHopGroupMember(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (next_hop_group_member_id) {
	*next_hop_group_member_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_next_hop_group_member(sai_object_id_t next_hop_group_member_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveNextHopGroupMemberRequest req;
	lemming::dataplane::sai::RemoveNextHopGroupMemberResponse resp;
	grpc::ClientContext context;
	req.set_oid(next_hop_group_member_id); 
	
	grpc::Status status = next_hop_group->RemoveNextHopGroupMember(&context, req, &resp);
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

sai_status_t l_set_next_hop_group_member_attribute(sai_object_id_t next_hop_group_member_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetNextHopGroupMemberAttributeRequest req;
	lemming::dataplane::sai::SetNextHopGroupMemberAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(next_hop_group_member_id); 
	
	
	

switch (attr->id) {
  
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_ID:
	req.set_next_hop_id(attr->value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_WEIGHT:
	req.set_weight(attr->value.u32);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_MONITORED_OBJECT:
	req.set_monitored_object(attr->value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_SEQUENCE_ID:
	req.set_sequence_id(attr->value.u32);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_COUNTER_ID:
	req.set_counter_id(attr->value.oid);
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_ARS_ALTERNATE_PATH:
	req.set_ars_alternate_path(attr->value.booldata);
	break;
}

	grpc::Status status = next_hop_group->SetNextHopGroupMemberAttribute(&context, req, &resp);
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

sai_status_t l_get_next_hop_group_member_attribute(sai_object_id_t next_hop_group_member_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetNextHopGroupMemberAttributeRequest req;
	lemming::dataplane::sai::GetNextHopGroupMemberAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(next_hop_group_member_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_next_hop_group_member_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = next_hop_group->GetNextHopGroupMemberAttribute(&context, req, &resp);
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
  
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_GROUP_ID:
	 attr_list[i].value.oid =   resp.attr().next_hop_group_id();
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_ID:
	 attr_list[i].value.oid =   resp.attr().next_hop_id();
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_WEIGHT:
	 attr_list[i].value.u32 =   resp.attr().weight();
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_CONFIGURED_ROLE:
	 attr_list[i].value.s32 =  convert_sai_next_hop_group_member_configured_role_t_to_sai(resp.attr().configured_role());
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_OBSERVED_ROLE:
	 attr_list[i].value.s32 =  convert_sai_next_hop_group_member_observed_role_t_to_sai(resp.attr().observed_role());
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_MONITORED_OBJECT:
	 attr_list[i].value.oid =   resp.attr().monitored_object();
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_INDEX:
	 attr_list[i].value.u32 =   resp.attr().index();
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_SEQUENCE_ID:
	 attr_list[i].value.u32 =   resp.attr().sequence_id();
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_COUNTER_ID:
	 attr_list[i].value.oid =   resp.attr().counter_id();
	break;
  case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_ARS_ALTERNATE_PATH:
	 attr_list[i].value.booldata =   resp.attr().ars_alternate_path();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_next_hop_group_members(sai_object_id_t switch_id, uint32_t object_count, const uint32_t *attr_count, const sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_object_id_t *object_id, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateNextHopGroupMembersRequest req;
	lemming::dataplane::sai::CreateNextHopGroupMembersResponse resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) {
		auto r = convert_create_next_hop_group_member(switch_id, attr_count[i],attr_list[i]);
		*req.add_reqs() = r;
	}

	grpc::Status status = next_hop_group->CreateNextHopGroupMembers(&context, req, &resp);
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
		object_id[i] = resp.resps(i).oid(); 
		object_statuses[i] = SAI_STATUS_SUCCESS;
	}

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_next_hop_group_members(uint32_t object_count, const sai_object_id_t *object_id, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveNextHopGroupMembersRequest req;
	lemming::dataplane::sai::RemoveNextHopGroupMembersResponse resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) {
		req.add_reqs()->set_oid(object_id[i]); 
		
	}

	grpc::Status status = next_hop_group->RemoveNextHopGroupMembers(&context, req, &resp);
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

sai_status_t l_create_next_hop_group_map(sai_object_id_t *next_hop_group_map_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateNextHopGroupMapRequest req = convert_create_next_hop_group_map(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateNextHopGroupMapResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = next_hop_group->CreateNextHopGroupMap(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (next_hop_group_map_id) {
	*next_hop_group_map_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_next_hop_group_map(sai_object_id_t next_hop_group_map_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveNextHopGroupMapRequest req;
	lemming::dataplane::sai::RemoveNextHopGroupMapResponse resp;
	grpc::ClientContext context;
	req.set_oid(next_hop_group_map_id); 
	
	grpc::Status status = next_hop_group->RemoveNextHopGroupMap(&context, req, &resp);
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

sai_status_t l_set_next_hop_group_map_attribute(sai_object_id_t next_hop_group_map_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_next_hop_group_map_attribute(sai_object_id_t next_hop_group_map_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetNextHopGroupMapAttributeRequest req;
	lemming::dataplane::sai::GetNextHopGroupMapAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(next_hop_group_map_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_next_hop_group_map_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = next_hop_group->GetNextHopGroupMapAttribute(&context, req, &resp);
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
  
  case SAI_NEXT_HOP_GROUP_MAP_ATTR_TYPE:
	 attr_list[i].value.s32 =  convert_sai_next_hop_group_map_type_t_to_sai(resp.attr().type());
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_next_hop_group_members_attribute(uint32_t object_count, const sai_object_id_t *object_id, const sai_attribute_t *attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_next_hop_group_members_attribute(uint32_t object_count, const sai_object_id_t *object_id, const uint32_t *attr_count, sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_next_hop_groups(sai_object_id_t switch_id, uint32_t object_count, const uint32_t *attr_count, const sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_object_id_t *object_id, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateNextHopGroupsRequest req;
	lemming::dataplane::sai::CreateNextHopGroupsResponse resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) {
		auto r = convert_create_next_hop_group(switch_id, attr_count[i],attr_list[i]);
		*req.add_reqs() = r;
	}

	grpc::Status status = next_hop_group->CreateNextHopGroups(&context, req, &resp);
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
		object_id[i] = resp.resps(i).oid(); 
		object_statuses[i] = SAI_STATUS_SUCCESS;
	}

	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_next_hop_groups(uint32_t object_count, const sai_object_id_t *object_id, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveNextHopGroupsRequest req;
	lemming::dataplane::sai::RemoveNextHopGroupsResponse resp;
	grpc::ClientContext context;

	for (uint32_t i = 0; i < object_count; i++) {
		req.add_reqs()->set_oid(object_id[i]); 
		
	}

	grpc::Status status = next_hop_group->RemoveNextHopGroups(&context, req, &resp);
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

sai_status_t l_set_next_hop_groups_attribute(uint32_t object_count, const sai_object_id_t *object_id, const sai_attribute_t *attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__ << " is not implemented but by-passing check";
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_next_hop_groups_attribute(uint32_t object_count, const sai_object_id_t *object_id, const uint32_t *attr_count, sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__ << " is not implemented but by-passing check";
	return SAI_STATUS_SUCCESS;
}

