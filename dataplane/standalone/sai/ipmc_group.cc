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

#include "dataplane/standalone/sai/ipmc_group.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/ipmc_group.pb.h"
#include <glog/logging.h>

const sai_ipmc_group_api_t l_ipmc_group = {
	.create_ipmc_group = l_create_ipmc_group,
	.remove_ipmc_group = l_remove_ipmc_group,
	.set_ipmc_group_attribute = l_set_ipmc_group_attribute,
	.get_ipmc_group_attribute = l_get_ipmc_group_attribute,
	.create_ipmc_group_member = l_create_ipmc_group_member,
	.remove_ipmc_group_member = l_remove_ipmc_group_member,
	.set_ipmc_group_member_attribute = l_set_ipmc_group_member_attribute,
	.get_ipmc_group_member_attribute = l_get_ipmc_group_member_attribute,
};


lemming::dataplane::sai::CreateIpmcGroupRequest convert_create_ipmc_group(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateIpmcGroupRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
}

}
return msg;
}

lemming::dataplane::sai::CreateIpmcGroupMemberRequest convert_create_ipmc_group_member(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateIpmcGroupMemberRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_IPMC_GROUP_MEMBER_ATTR_IPMC_GROUP_ID:
	msg.set_ipmc_group_id(attr_list[i].value.oid);
	break;
  case SAI_IPMC_GROUP_MEMBER_ATTR_IPMC_OUTPUT_ID:
	msg.set_ipmc_output_id(attr_list[i].value.oid);
	break;
}

}
return msg;
}

sai_status_t l_create_ipmc_group(sai_object_id_t *ipmc_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateIpmcGroupRequest req = convert_create_ipmc_group(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateIpmcGroupResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = ipmc_group->CreateIpmcGroup(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (ipmc_group_id) {
	*ipmc_group_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_ipmc_group(sai_object_id_t ipmc_group_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveIpmcGroupRequest req;
	lemming::dataplane::sai::RemoveIpmcGroupResponse resp;
	grpc::ClientContext context;
	req.set_oid(ipmc_group_id); 
	
	grpc::Status status = ipmc_group->RemoveIpmcGroup(&context, req, &resp);
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

sai_status_t l_set_ipmc_group_attribute(sai_object_id_t ipmc_group_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_ipmc_group_attribute(sai_object_id_t ipmc_group_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetIpmcGroupAttributeRequest req;
	lemming::dataplane::sai::GetIpmcGroupAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(ipmc_group_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_ipmc_group_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = ipmc_group->GetIpmcGroupAttribute(&context, req, &resp);
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
  
  case SAI_IPMC_GROUP_ATTR_IPMC_OUTPUT_COUNT:
	 attr_list[i].value.u32 =   resp.attr().ipmc_output_count();
	break;
  case SAI_IPMC_GROUP_ATTR_IPMC_MEMBER_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().ipmc_member_list(), &attr_list[i].value.objlist.count);
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_ipmc_group_member(sai_object_id_t *ipmc_group_member_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateIpmcGroupMemberRequest req = convert_create_ipmc_group_member(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateIpmcGroupMemberResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = ipmc_group->CreateIpmcGroupMember(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (ipmc_group_member_id) {
	*ipmc_group_member_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_ipmc_group_member(sai_object_id_t ipmc_group_member_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveIpmcGroupMemberRequest req;
	lemming::dataplane::sai::RemoveIpmcGroupMemberResponse resp;
	grpc::ClientContext context;
	req.set_oid(ipmc_group_member_id); 
	
	grpc::Status status = ipmc_group->RemoveIpmcGroupMember(&context, req, &resp);
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

sai_status_t l_set_ipmc_group_member_attribute(sai_object_id_t ipmc_group_member_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_ipmc_group_member_attribute(sai_object_id_t ipmc_group_member_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetIpmcGroupMemberAttributeRequest req;
	lemming::dataplane::sai::GetIpmcGroupMemberAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(ipmc_group_member_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_ipmc_group_member_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = ipmc_group->GetIpmcGroupMemberAttribute(&context, req, &resp);
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
  
  case SAI_IPMC_GROUP_MEMBER_ATTR_IPMC_GROUP_ID:
	 attr_list[i].value.oid =   resp.attr().ipmc_group_id();
	break;
  case SAI_IPMC_GROUP_MEMBER_ATTR_IPMC_OUTPUT_ID:
	 attr_list[i].value.oid =   resp.attr().ipmc_output_id();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

