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

#include "dataplane/standalone/sai/udf.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/udf.pb.h"
#include <glog/logging.h>

const sai_udf_api_t l_udf = {
	.create_udf = l_create_udf,
	.remove_udf = l_remove_udf,
	.set_udf_attribute = l_set_udf_attribute,
	.get_udf_attribute = l_get_udf_attribute,
	.create_udf_match = l_create_udf_match,
	.remove_udf_match = l_remove_udf_match,
	.set_udf_match_attribute = l_set_udf_match_attribute,
	.get_udf_match_attribute = l_get_udf_match_attribute,
	.create_udf_group = l_create_udf_group,
	.remove_udf_group = l_remove_udf_group,
	.set_udf_group_attribute = l_set_udf_group_attribute,
	.get_udf_group_attribute = l_get_udf_group_attribute,
};


lemming::dataplane::sai::CreateUdfRequest convert_create_udf(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateUdfRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_UDF_ATTR_MATCH_ID:
	msg.set_match_id(attr_list[i].value.oid);
	break;
  case SAI_UDF_ATTR_GROUP_ID:
	msg.set_group_id(attr_list[i].value.oid);
	break;
  case SAI_UDF_ATTR_BASE:
	msg.set_base(convert_sai_udf_base_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_UDF_ATTR_OFFSET:
	msg.set_offset(attr_list[i].value.u16);
	break;
  case SAI_UDF_ATTR_HASH_MASK:
	msg.mutable_hash_mask()->Add(attr_list[i].value.u8list.list, attr_list[i].value.u8list.list + attr_list[i].value.u8list.count);
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateUdfMatchRequest convert_create_udf_match(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateUdfMatchRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_UDF_MATCH_ATTR_L2_TYPE:
	*msg.mutable_l2_type() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_UDF_MATCH_ATTR_L3_TYPE:
	*msg.mutable_l3_type() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u8, attr_list[i].value.aclfield.mask.u8 );
	break;
  case SAI_UDF_MATCH_ATTR_GRE_TYPE:
	*msg.mutable_gre_type() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
  case SAI_UDF_MATCH_ATTR_PRIORITY:
	msg.set_priority(attr_list[i].value.u8);
	break;
  case SAI_UDF_MATCH_ATTR_L4_DST_PORT_TYPE:
	*msg.mutable_l4_dst_port_type() = convert_from_acl_field_data( attr_list[i].value.aclfield, attr_list[i].value.aclfield.data.u16, attr_list[i].value.aclfield.mask.u16 );
	break;
}

}
return msg;
}

lemming::dataplane::sai::CreateUdfGroupRequest convert_create_udf_group(sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {

lemming::dataplane::sai::CreateUdfGroupRequest msg;


 for(uint32_t i = 0; i < attr_count; i++ ) {
	
	

switch (attr_list[i].id) {
  
  case SAI_UDF_GROUP_ATTR_TYPE:
	msg.set_type(convert_sai_udf_group_type_t_to_proto(attr_list[i].value.s32));
	break;
  case SAI_UDF_GROUP_ATTR_LENGTH:
	msg.set_length(attr_list[i].value.u16);
	break;
}

}
return msg;
}

sai_status_t l_create_udf(sai_object_id_t *udf_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateUdfRequest req = convert_create_udf(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateUdfResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = udf->CreateUdf(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (udf_id) {
	*udf_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_udf(sai_object_id_t udf_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveUdfRequest req;
	lemming::dataplane::sai::RemoveUdfResponse resp;
	grpc::ClientContext context;
	req.set_oid(udf_id); 
	
	grpc::Status status = udf->RemoveUdf(&context, req, &resp);
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

sai_status_t l_set_udf_attribute(sai_object_id_t udf_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::SetUdfAttributeRequest req;
	lemming::dataplane::sai::SetUdfAttributeResponse resp;
	grpc::ClientContext context;
	req.set_oid(udf_id); 
	
	
	

switch (attr->id) {
  
  case SAI_UDF_ATTR_BASE:
	req.set_base(convert_sai_udf_base_t_to_proto(attr->value.s32));
	break;
  case SAI_UDF_ATTR_HASH_MASK:
	req.mutable_hash_mask()->Add(attr->value.u8list.list, attr->value.u8list.list + attr->value.u8list.count);
	break;
}

	grpc::Status status = udf->SetUdfAttribute(&context, req, &resp);
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

sai_status_t l_get_udf_attribute(sai_object_id_t udf_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetUdfAttributeRequest req;
	lemming::dataplane::sai::GetUdfAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(udf_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_udf_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = udf->GetUdfAttribute(&context, req, &resp);
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
  
  case SAI_UDF_ATTR_MATCH_ID:
	 attr_list[i].value.oid =   resp.attr().match_id();
	break;
  case SAI_UDF_ATTR_GROUP_ID:
	 attr_list[i].value.oid =   resp.attr().group_id();
	break;
  case SAI_UDF_ATTR_BASE:
	 attr_list[i].value.s32 =  convert_sai_udf_base_t_to_sai(resp.attr().base());
	break;
  case SAI_UDF_ATTR_OFFSET:
	 attr_list[i].value.u16 =   resp.attr().offset();
	break;
  case SAI_UDF_ATTR_HASH_MASK:
	copy_list(attr_list[i].value.u8list.list, resp.attr().hash_mask(), &attr_list[i].value.u8list.count);
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_udf_match(sai_object_id_t *udf_match_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateUdfMatchRequest req = convert_create_udf_match(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateUdfMatchResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = udf->CreateUdfMatch(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (udf_match_id) {
	*udf_match_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_udf_match(sai_object_id_t udf_match_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveUdfMatchRequest req;
	lemming::dataplane::sai::RemoveUdfMatchResponse resp;
	grpc::ClientContext context;
	req.set_oid(udf_match_id); 
	
	grpc::Status status = udf->RemoveUdfMatch(&context, req, &resp);
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

sai_status_t l_set_udf_match_attribute(sai_object_id_t udf_match_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_udf_match_attribute(sai_object_id_t udf_match_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetUdfMatchAttributeRequest req;
	lemming::dataplane::sai::GetUdfMatchAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(udf_match_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_udf_match_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = udf->GetUdfMatchAttribute(&context, req, &resp);
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
  
  case SAI_UDF_MATCH_ATTR_L2_TYPE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().l2_type(), resp.attr().l2_type().data_uint(), resp.attr().l2_type().mask_uint());
	break;
  case SAI_UDF_MATCH_ATTR_L3_TYPE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u8(resp.attr().l3_type(), resp.attr().l3_type().data_uint(), resp.attr().l3_type().mask_uint());
	break;
  case SAI_UDF_MATCH_ATTR_GRE_TYPE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().gre_type(), resp.attr().gre_type().data_uint(), resp.attr().gre_type().mask_uint());
	break;
  case SAI_UDF_MATCH_ATTR_PRIORITY:
	 attr_list[i].value.u8 =   resp.attr().priority();
	break;
  case SAI_UDF_MATCH_ATTR_L4_DST_PORT_TYPE:
	 attr_list[i].value.aclfield =  convert_to_acl_field_data_u16(resp.attr().l4_dst_port_type(), resp.attr().l4_dst_port_type().data_uint(), resp.attr().l4_dst_port_type().mask_uint());
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_udf_group(sai_object_id_t *udf_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::CreateUdfGroupRequest req = convert_create_udf_group(switch_id, attr_count, attr_list);
	lemming::dataplane::sai::CreateUdfGroupResponse resp;
	grpc::ClientContext context;
	 req.set_switch_(switch_id); 
	
	grpc::Status status = udf->CreateUdfGroup(&context, req, &resp);
	if (!status.ok()) {
		auto it = context.GetServerTrailingMetadata().find("traceparent");
		if (it != context.GetServerTrailingMetadata().end()) {
			LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second << " msg: " << status.error_message(); 
		} else {
			LOG(ERROR) << "Lucius RPC error: " << status.error_message(); 
		}
		return SAI_STATUS_FAILURE;
	}
	if (udf_group_id) {
	*udf_group_id = resp.oid(); 
  	}
	
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_udf_group(sai_object_id_t udf_group_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::RemoveUdfGroupRequest req;
	lemming::dataplane::sai::RemoveUdfGroupResponse resp;
	grpc::ClientContext context;
	req.set_oid(udf_group_id); 
	
	grpc::Status status = udf->RemoveUdfGroup(&context, req, &resp);
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

sai_status_t l_set_udf_group_attribute(sai_object_id_t udf_group_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_udf_group_attribute(sai_object_id_t udf_group_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	
	lemming::dataplane::sai::GetUdfGroupAttributeRequest req;
	lemming::dataplane::sai::GetUdfGroupAttributeResponse resp;
	grpc::ClientContext context;
	
	req.set_oid(udf_group_id); 

	for (uint32_t i = 0; i < attr_count; i++) {
		req.add_attr_type(convert_sai_udf_group_attr_t_to_proto(attr_list[i].id));
	}
	grpc::Status status = udf->GetUdfGroupAttribute(&context, req, &resp);
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
  
  case SAI_UDF_GROUP_ATTR_UDF_LIST:
	copy_list(attr_list[i].value.objlist.list, resp.attr().udf_list(), &attr_list[i].value.objlist.count);
	break;
  case SAI_UDF_GROUP_ATTR_TYPE:
	 attr_list[i].value.s32 =  convert_sai_udf_group_type_t_to_sai(resp.attr().type());
	break;
  case SAI_UDF_GROUP_ATTR_LENGTH:
	 attr_list[i].value.u16 =   resp.attr().length();
	break;
}

	}
	
	return SAI_STATUS_SUCCESS;
}

