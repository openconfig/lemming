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

#include "dataplane/standalone/sai/rpf_group.h"
#include <glog/logging.h>
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_rpf_group_api_t l_rpf_group = {
	.create_rpf_group = l_create_rpf_group,
	.remove_rpf_group = l_remove_rpf_group,
	.set_rpf_group_attribute = l_set_rpf_group_attribute,
	.get_rpf_group_attribute = l_get_rpf_group_attribute,
	.create_rpf_group_member = l_create_rpf_group_member,
	.remove_rpf_group_member = l_remove_rpf_group_member,
	.set_rpf_group_member_attribute = l_set_rpf_group_member_attribute,
	.get_rpf_group_member_attribute = l_get_rpf_group_member_attribute,
};


sai_status_t l_create_rpf_group(sai_object_id_t *rpf_group_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_RPF_GROUP, rpf_group_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_rpf_group(sai_object_id_t rpf_group_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_RPF_GROUP, rpf_group_id);
}

sai_status_t l_set_rpf_group_attribute(sai_object_id_t rpf_group_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_RPF_GROUP, rpf_group_id, attr);
}

sai_status_t l_get_rpf_group_attribute(sai_object_id_t rpf_group_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_RPF_GROUP, rpf_group_id, attr_count, attr_list);
}

sai_status_t l_create_rpf_group_member(sai_object_id_t *rpf_group_member_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_RPF_GROUP_MEMBER, rpf_group_member_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_rpf_group_member(sai_object_id_t rpf_group_member_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_RPF_GROUP_MEMBER, rpf_group_member_id);
}

sai_status_t l_set_rpf_group_member_attribute(sai_object_id_t rpf_group_member_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_RPF_GROUP_MEMBER, rpf_group_member_id, attr);
}

sai_status_t l_get_rpf_group_member_attribute(sai_object_id_t rpf_group_member_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_RPF_GROUP_MEMBER, rpf_group_member_id, attr_count, attr_list);
}

